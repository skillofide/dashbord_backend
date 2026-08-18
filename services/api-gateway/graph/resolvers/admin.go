package resolvers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/skillofide/api-gateway/middleware"
	pkgauth "github.com/skillofide/pkg/auth"
)

// AdminHandler handles /api/admin/* REST endpoints.
// It holds its own Postgres pool so it stays independent of the user-service module.
type AdminHandler struct {
	Pool *pgxpool.Pool
	Log  *zap.Logger
	// Inquiries serves the enquiry list; nil disables those routes.
	Inquiries *InquiryHandler
}

// ServeHTTP dispatches admin REST routes.
func (h *AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Role guard — admin only
	if middleware.RoleFromContext(r.Context()) != "admin" {
		h.jsonErr(w, http.StatusForbidden, "admin access required")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/admin")
	path = strings.TrimRight(path, "/")

	switch {
	// POST /api/admin/bulk-import
	case path == "/bulk-import" && r.Method == http.MethodPost:
		h.handleBulkImport(w, r)

	// GET /api/admin/inquiries
	case path == "/inquiries" && r.Method == http.MethodGet && h.Inquiries != nil:
		h.Inquiries.ListInquiries(w, r)

	// PATCH /api/admin/inquiries/{id}   — status / notes
	case strings.HasPrefix(path, "/inquiries/") && r.Method == http.MethodPatch && h.Inquiries != nil:
		h.Inquiries.UpdateInquiry(w, r, strings.TrimPrefix(path, "/inquiries/"))

	// GET /api/admin/users
	case path == "/users" && r.Method == http.MethodGet:
		h.handleListUsers(w, r)

	// PATCH /api/admin/users/{id}   — update role
	case strings.HasPrefix(path, "/users/") && !strings.Contains(path[len("/users/"):], "/") && r.Method == http.MethodPatch:
		userID := strings.TrimPrefix(path, "/users/")
		h.handleUpdateUser(w, r, userID)

	// DELETE /api/admin/users/{id}
	case strings.HasPrefix(path, "/users/") && !strings.Contains(path[len("/users/"):], "/") && r.Method == http.MethodDelete:
		userID := strings.TrimPrefix(path, "/users/")
		h.handleDeleteUser(w, r, userID)

	// POST /api/admin/users/{id}/courses
	case strings.HasPrefix(path, "/users/") && strings.HasSuffix(path, "/courses") && r.Method == http.MethodPost:
		inner := strings.TrimPrefix(path, "/users/")
		userID := strings.TrimSuffix(inner, "/courses")
		h.handleGrantCourse(w, r, userID)

	// DELETE /api/admin/users/{id}/courses/{courseId}
	case strings.HasPrefix(path, "/users/") && strings.Contains(path, "/courses/") && r.Method == http.MethodDelete:
		// path: /users/{id}/courses/{courseId}
		inner := strings.TrimPrefix(path, "/users/")
		parts := strings.SplitN(inner, "/courses/", 2)
		if len(parts) == 2 {
			h.handleRevokeCourse(w, r, parts[0], parts[1])
		} else {
			h.jsonErr(w, http.StatusNotFound, "not found")
		}

	default:
		h.jsonErr(w, http.StatusNotFound, "not found")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Data types
// ─────────────────────────────────────────────────────────────────────────────

type importUserRow struct {
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	Role      string   `json:"role"`
	CourseIDs []string `json:"course_ids"`
}

type importRowResult struct {
	Email   string `json:"email"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type adminUserRow struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	Name      string   `json:"name"`
	Role      string   `json:"role"`
	CourseIDs []string `json:"course_ids"`
	CreatedAt string   `json:"created_at"`
}

// ─────────────────────────────────────────────────────────────────────────────
// Handlers
// ─────────────────────────────────────────────────────────────────────────────

func (h *AdminHandler) handleBulkImport(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Users []importUserRow `json:"users"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.jsonErr(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx := context.Background()
	results := h.bulkUpsertUsers(ctx, req.Users)

	successCount := 0
	for _, res := range results {
		if res.Success {
			successCount++
		}
	}

	h.json(w, http.StatusOK, map[string]interface{}{
		"total":   len(results),
		"success": successCount,
		"failed":  len(results) - successCount,
		"results": results,
	})
}

func (h *AdminHandler) handleListUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	search := q.Get("search")
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	users, total, err := h.listUsers(r.Context(), page, pageSize, search)
	if err != nil {
		h.Log.Error("list admin users failed", zap.Error(err))
		h.jsonErr(w, http.StatusInternalServerError, "failed to list users")
		return
	}

	h.json(w, http.StatusOK, map[string]interface{}{
		"users":    users,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *AdminHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request, userID string) {
	var req struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.jsonErr(w, http.StatusBadRequest, "invalid request body")
		return
	}
	// `recruiter` grants access to the partner hiring portal; a recruiter still
	// needs a company_members row before they can see any drive.
	if req.Role != "student" && req.Role != "admin" && req.Role != "recruiter" {
		h.jsonErr(w, http.StatusBadRequest, "role must be 'student', 'recruiter' or 'admin'")
		return
	}
	tag, err := h.Pool.Exec(r.Context(),
		`UPDATE users SET role = $2, updated_at = now() WHERE id = $1::uuid`, userID, req.Role)
	if err != nil {
		h.Log.Error("update user role failed", zap.String("userID", userID), zap.Error(err))
		h.jsonErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if tag.RowsAffected() == 0 {
		h.jsonErr(w, http.StatusNotFound, "user not found")
		return
	}
	h.json(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *AdminHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request, userID string) {
	tag, err := h.Pool.Exec(r.Context(), `DELETE FROM users WHERE id = $1::uuid`, userID)
	if err != nil {
		h.Log.Error("delete user failed", zap.String("userID", userID), zap.Error(err))
		h.jsonErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if tag.RowsAffected() == 0 {
		h.jsonErr(w, http.StatusNotFound, "user not found")
		return
	}
	h.json(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *AdminHandler) handleGrantCourse(w http.ResponseWriter, r *http.Request, userID string) {
	var req struct {
		CourseID string `json:"course_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.CourseID == "" {
		h.jsonErr(w, http.StatusBadRequest, "course_id is required")
		return
	}
	_, err := h.Pool.Exec(r.Context(), `
		INSERT INTO user_courses (user_id, course_id)
		VALUES ($1::uuid, $2)
		ON CONFLICT (user_id, course_id) DO NOTHING
	`, userID, req.CourseID)
	if err != nil {
		h.Log.Error("grant course access failed", zap.Error(err))
		h.jsonErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.json(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *AdminHandler) handleRevokeCourse(w http.ResponseWriter, r *http.Request, userID, courseID string) {
	_, err := h.Pool.Exec(r.Context(),
		`DELETE FROM user_courses WHERE user_id = $1::uuid AND course_id = $2`, userID, courseID)
	if err != nil {
		h.Log.Error("revoke course access failed", zap.Error(err))
		h.jsonErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.json(w, http.StatusOK, map[string]bool{"success": true})
}

// ─────────────────────────────────────────────────────────────────────────────
// Internal DB helpers
// ─────────────────────────────────────────────────────────────────────────────

func (h *AdminHandler) bulkUpsertUsers(ctx context.Context, rows []importUserRow) []importRowResult {
	results := make([]importRowResult, 0, len(rows))
	for _, row := range rows {
		if row.Email == "" || row.Name == "" || row.Password == "" {
			results = append(results, importRowResult{
				Email:   row.Email,
				Success: false,
				Message: "email, name and password are required",
			})
			continue
		}
		role := row.Role
		if role == "" {
			role = "student"
		}

		hashed, err := pkgauth.HashPassword(row.Password)
		if err != nil {
			results = append(results, importRowResult{
				Email:   row.Email,
				Success: false,
				Message: fmt.Sprintf("hash password: %v", err),
			})
			continue
		}

		var userID string
		err = h.Pool.QueryRow(ctx, `
			INSERT INTO users (email, name, password, role)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (email)
			DO UPDATE SET name = EXCLUDED.name, password = EXCLUDED.password, role = EXCLUDED.role, updated_at = now()
			RETURNING id::text
		`, row.Email, row.Name, hashed, role).Scan(&userID)
		if err != nil {
			results = append(results, importRowResult{
				Email:   row.Email,
				Success: false,
				Message: fmt.Sprintf("upsert user: %v", err),
			})
			continue
		}

		for _, courseID := range row.CourseIDs {
			if courseID == "" {
				continue
			}
			_, _ = h.Pool.Exec(ctx, `
				INSERT INTO user_courses (user_id, course_id)
				VALUES ($1::uuid, $2)
				ON CONFLICT (user_id, course_id) DO NOTHING
			`, userID, courseID)
		}

		results = append(results, importRowResult{
			Email:   row.Email,
			Success: true,
			Message: "imported successfully",
		})
	}
	return results
}

func (h *AdminHandler) listUsers(ctx context.Context, page, pageSize int, search string) ([]adminUserRow, int, error) {
	offset := (page - 1) * pageSize

	var total int
	countQuery := `SELECT COUNT(*) FROM users`
	countArgs := []interface{}{}
	if search != "" {
		countQuery = `SELECT COUNT(*) FROM users WHERE email ILIKE $1 OR name ILIKE $1`
		countArgs = []interface{}{"%" + search + "%"}
	}
	if err := h.Pool.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	var query string
	var args []interface{}
	if search != "" {
		query = `
			SELECT u.id::text, u.email, u.name, u.role, u.created_at,
			       COALESCE(array_agg(uc.course_id) FILTER (WHERE uc.course_id IS NOT NULL), '{}') AS course_ids
			FROM users u
			LEFT JOIN user_courses uc ON uc.user_id = u.id
			WHERE u.email ILIKE $3 OR u.name ILIKE $3
			GROUP BY u.id, u.email, u.name, u.role, u.created_at
			ORDER BY u.created_at DESC
			LIMIT $1 OFFSET $2`
		args = []interface{}{pageSize, offset, "%" + search + "%"}
	} else {
		query = `
			SELECT u.id::text, u.email, u.name, u.role, u.created_at,
			       COALESCE(array_agg(uc.course_id) FILTER (WHERE uc.course_id IS NOT NULL), '{}') AS course_ids
			FROM users u
			LEFT JOIN user_courses uc ON uc.user_id = u.id
			GROUP BY u.id, u.email, u.name, u.role, u.created_at
			ORDER BY u.created_at DESC
			LIMIT $1 OFFSET $2`
		args = []interface{}{pageSize, offset}
	}

	pgRows, err := h.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	defer pgRows.Close()

	var users []adminUserRow
	for pgRows.Next() {
		var u adminUserRow
		var createdAt time.Time
		var courseIDs []string
		if err := pgRows.Scan(&u.ID, &u.Email, &u.Name, &u.Role, &createdAt, &courseIDs); err != nil {
			return nil, 0, fmt.Errorf("scan user: %w", err)
		}
		u.CreatedAt = createdAt.Format(time.RFC3339)
		if courseIDs == nil {
			courseIDs = []string{}
		}
		u.CourseIDs = courseIDs
		users = append(users, u)
	}
	if users == nil {
		users = []adminUserRow{}
	}
	return users, total, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// JSON helpers
// ─────────────────────────────────────────────────────────────────────────────

func (h *AdminHandler) json(w http.ResponseWriter, status int, v interface{}) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) //nolint:errcheck
}

func (h *AdminHandler) jsonErr(w http.ResponseWriter, status int, msg string) {
	h.json(w, status, map[string]string{"error": msg})
}

// Ensure pgx is used (avoid import cycle if pool comes from api-gateway's own init).
var _ = pgx.ErrNoRows
