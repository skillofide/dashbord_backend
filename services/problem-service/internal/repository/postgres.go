// Package repository provides Postgres-backed data access for the problem service.
package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	problemv1 "github.com/skillofide/proto/problem/v1"
)

// ProblemRepository wraps a pgxpool and implements all problem data queries.
type ProblemRepository struct {
	pool *pgxpool.Pool
}

// New constructs a ProblemRepository.
func New(pool *pgxpool.Pool) *ProblemRepository {
	return &ProblemRepository{pool: pool}
}

// ListProblems returns a filtered, paginated list of problems (summary fields only).
func (r *ProblemRepository) ListProblems(ctx context.Context, req *problemv1.ListProblemsRequest) (*problemv1.ListProblemsResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}
	offset := (page - 1) * pageSize

	// Build dynamic WHERE clause
	var filterClauses []string
	var filterArgs []interface{}

	if req.SetId != "" {
		filterArgs = append(filterArgs, req.SetId)
		filterClauses = append(filterClauses, fmt.Sprintf("p.set_id::text = $%d", len(filterArgs)))
	}
	if req.Topic != "" {
		filterArgs = append(filterArgs, req.Topic)
		filterClauses = append(filterClauses, fmt.Sprintf("p.topic = $%d", len(filterArgs)))
	}
	if req.Difficulty != "" {
		filterArgs = append(filterArgs, req.Difficulty)
		filterClauses = append(filterClauses, fmt.Sprintf("p.difficulty = $%d", len(filterArgs)))
	}

	whereSQL := ""
	if len(filterClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(filterClauses, " AND ")
	}

	// Build pagination args
	mainArgs := make([]interface{}, len(filterArgs))
	copy(mainArgs, filterArgs)

	var userIdParamIndex int
	if req.UserId != "" {
		mainArgs = append(mainArgs, req.UserId)
		userIdParamIndex = len(mainArgs)
	}

	limitOffset := fmt.Sprintf("LIMIT $%d OFFSET $%d", len(mainArgs)+1, len(mainArgs)+2)
	paginationArgs := append(mainArgs, pageSize, offset)

	var query string
	if req.UserId != "" {
		query = fmt.Sprintf(`
			SELECT p.id, p.slug, p.title, p.difficulty, p.topic, p.xp,
			       COALESCE(p.set_id::text, '') AS set_id,
			       COALESCE(pus.status, 'Unsolved') AS user_status
			FROM   problems p
			LEFT JOIN problem_user_status pus ON pus.problem_id = p.id AND pus.user_id = $%d
			%s
			ORDER  BY p.created_at DESC
			%s
		`, userIdParamIndex, whereSQL, limitOffset)
	} else {
		query = fmt.Sprintf(`
			SELECT p.id, p.slug, p.title, p.difficulty, p.topic, p.xp,
			       COALESCE(p.set_id::text, '') AS set_id,
			       'Unsolved' AS user_status
			FROM   problems p
			%s
			ORDER  BY p.created_at DESC
			%s
		`, whereSQL, limitOffset)
	}

	rows, err := r.pool.Query(ctx, query, paginationArgs...)
	if err != nil {
		return nil, fmt.Errorf("list problems query: %w", err)
	}
	defer rows.Close()

	var problems []*problemv1.Problem
	for rows.Next() {
		p := &problemv1.Problem{}
		if err := rows.Scan(&p.Id, &p.Slug, &p.Title, &p.Difficulty, &p.Topic, &p.Xp, &p.SetId, &p.UserStatus); err != nil {
			return nil, fmt.Errorf("scan problem row: %w", err)
		}
		problems = append(problems, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate problem rows: %w", err)
	}

	// Count total matching rows
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM problems p %s", whereSQL)
	var total int32
	_ = r.pool.QueryRow(ctx, countQuery, filterArgs...).Scan(&total)

	return &problemv1.ListProblemsResponse{
		Problems: problems,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetProblem returns the full detail of a single problem by UUID or slug.
func (r *ProblemRepository) GetProblem(ctx context.Context, req *problemv1.GetProblemRequest) (*problemv1.Problem, error) {
	p := &problemv1.Problem{}

	err := r.pool.QueryRow(ctx, `
		SELECT id, slug, title, difficulty, topic, xp, statement,
		       COALESCE(set_id::text, '') AS set_id
		FROM   problems
		WHERE  id::text = $1 OR slug = $1
	`, req.Id).Scan(&p.Id, &p.Slug, &p.Title, &p.Difficulty, &p.Topic, &p.Xp, &p.Statement, &p.SetId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("problem not found: %s", req.Id)
		}
		return nil, fmt.Errorf("get problem: %w", err)
	}

	// Constraints
	cRows, _ := r.pool.Query(ctx,
		"SELECT constraint_text FROM problem_constraints WHERE problem_id = $1 ORDER BY order_index", p.Id)
	if cRows != nil {
		for cRows.Next() {
			var c string
			cRows.Scan(&c) //nolint:errcheck
			p.Constraints = append(p.Constraints, c)
		}
		cRows.Close()
	}

	// Tags
	tRows, _ := r.pool.Query(ctx, "SELECT tag FROM problem_tags WHERE problem_id = $1", p.Id)
	if tRows != nil {
		for tRows.Next() {
			var t string
			tRows.Scan(&t) //nolint:errcheck
			p.Tags = append(p.Tags, t)
		}
		tRows.Close()
	}

	// Examples
	eRows, _ := r.pool.Query(ctx,
		"SELECT input, output, explanation FROM examples WHERE problem_id = $1 ORDER BY order_index", p.Id)
	if eRows != nil {
		for eRows.Next() {
			ex := &problemv1.Example{}
			eRows.Scan(&ex.Input, &ex.Output, &ex.Explanation) //nolint:errcheck
			p.Examples = append(p.Examples, ex)
		}
		eRows.Close()
	}

	// Hints
	hRows, _ := r.pool.Query(ctx,
		"SELECT order_index, title, body FROM hints WHERE problem_id = $1 ORDER BY order_index", p.Id)
	if hRows != nil {
		for hRows.Next() {
			h := &problemv1.Hint{}
			hRows.Scan(&h.Order, &h.Title, &h.Body) //nolint:errcheck
			p.Hints = append(p.Hints, h)
		}
		hRows.Close()
	}

	// Starter codes
	sc := &problemv1.StarterCodes{}
	r.pool.QueryRow(ctx, //nolint:errcheck
		"SELECT javascript, python, java, cpp, go FROM starter_codes WHERE problem_id = $1", p.Id).
		Scan(&sc.Javascript, &sc.Python, &sc.Java, &sc.Cpp, &sc.Go)
	p.StarterCodes = sc

	return p, nil
}

// GetTestCases returns the test cases for a problem.
// When includeHidden is false (Run button), only visible test cases are returned.
func (r *ProblemRepository) GetTestCases(ctx context.Context, req *problemv1.GetTestCasesRequest) (*problemv1.GetTestCasesResponse, error) {
	query := `
		SELECT id, input, expected_output, is_hidden, time_limit_ms, memory_limit_mb, order_index
		FROM   test_cases
		WHERE  problem_id = $1
	`
	args := []interface{}{req.ProblemId}

	if !req.IncludeHidden {
		query += " AND is_hidden = false"
	}
	query += " ORDER BY order_index"

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get test cases: %w", err)
	}
	defer rows.Close()

	var testCases []*problemv1.TestCase
	for rows.Next() {
		tc := &problemv1.TestCase{}
		if err := rows.Scan(&tc.Id, &tc.Input, &tc.ExpectedOutput, &tc.IsHidden,
			&tc.TimeLimitMs, &tc.MemoryLimitMb, &tc.OrderIndex); err != nil {
			return nil, fmt.Errorf("scan test case: %w", err)
		}
		testCases = append(testCases, tc)
	}

	return &problemv1.GetTestCasesResponse{TestCases: testCases}, nil
}

// ListPracticeSets returns all practice sets with optional per-user progress.
func (r *ProblemRepository) ListPracticeSets(ctx context.Context, req *problemv1.ListPracticeSetsRequest) (*problemv1.ListPracticeSetsResponse, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT ps.id, ps.title, ps.level, ps.level_color, ps.bg_color,
		       COUNT(p.id) AS total_problems
		FROM   practice_sets ps
		LEFT   JOIN problems p ON p.set_id = ps.id
		GROUP  BY ps.id
		ORDER  BY ps.created_at
	`)
	if err != nil {
		return nil, fmt.Errorf("list practice sets: %w", err)
	}
	defer rows.Close()

	var sets []*problemv1.PracticeSet
	for rows.Next() {
		s := &problemv1.PracticeSet{}
		if err := rows.Scan(&s.Id, &s.Title, &s.Level, &s.LevelColor, &s.BgColor, &s.TotalProblems); err != nil {
			return nil, fmt.Errorf("scan practice set: %w", err)
		}

		// Per-user progress if userId provided
		if req.UserId != "" {
			var solved int32
			r.pool.QueryRow(ctx, `
				SELECT COUNT(*) FROM problem_user_status pus
				JOIN   problems p ON p.id = pus.problem_id
				WHERE  pus.user_id = $1 AND p.set_id = $2 AND pus.status = 'Solved'
			`, req.UserId, s.Id).Scan(&solved) //nolint:errcheck

			if s.TotalProblems > 0 {
				s.Progress = float32(solved) / float32(s.TotalProblems) * 100
			}
		}
		sets = append(sets, s)
	}

	return &problemv1.ListPracticeSetsResponse{PracticeSets: sets}, nil
}

// GetProblemUserStatus returns the submission status for a specific user+problem.
func (r *ProblemRepository) GetProblemUserStatus(ctx context.Context, req *problemv1.GetProblemUserStatusRequest) (*problemv1.ProblemUserStatus, error) {
	pus := &problemv1.ProblemUserStatus{
		UserId:    req.UserId,
		ProblemId: req.ProblemId,
		Status:    "Unsolved",
	}

	err := r.pool.QueryRow(ctx, `
		SELECT status, COALESCE(solved_at::text, ''), attempts
		FROM   problem_user_status
		WHERE  user_id = $1 AND problem_id = $2
	`, req.UserId, req.ProblemId).Scan(&pus.Status, &pus.SolvedAt, &pus.Attempts)

	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("get problem user status: %w", err)
	}

	return pus, nil
}
