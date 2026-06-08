package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	userv1 "github.com/skillofide/proto/user/v1"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

// VerifyUser verifies credentials against the database.
func (r *UserRepository) VerifyUser(ctx context.Context, email, password string) (*userv1.VerifyUserResponse, error) {
	var id, name, dbPassword, role string
	err := r.pool.QueryRow(ctx, `
		SELECT id::text, name, password, role
		FROM   users
		WHERE  email = $1
	`, email).Scan(&id, &name, &dbPassword, &role)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, fmt.Errorf("query user: %w", err)
	}

	// Plain-text check as currently done by api-gateway
	if password != dbPassword {
		return nil, fmt.Errorf("invalid password")
	}

	return &userv1.VerifyUserResponse{
		Id:    id,
		Email: email,
		Name:  name,
		Role:  role,
	}, nil
}

// CreateOrUpdateUser inserts or updates a user record.
func (r *UserRepository) CreateOrUpdateUser(ctx context.Context, email, name, password, role string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO users (email, name, password, role)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (email) 
		DO UPDATE SET name = EXCLUDED.name, password = EXCLUDED.password, role = EXCLUDED.role, updated_at = now();
	`, email, name, password, role)
	if err != nil {
		return fmt.Errorf("upsert user: %w", err)
	}
	return nil
}

// EnsureUsersTable creates the users table if missing and seeds the default admin user.
func (r *UserRepository) EnsureUsersTable(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email      TEXT NOT NULL UNIQUE,
			name       TEXT NOT NULL,
			password   TEXT NOT NULL,
			role       TEXT NOT NULL DEFAULT 'student',
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
	`)
	if err != nil {
		return fmt.Errorf("create users table: %w", err)
	}

	// Seed default user
	_, err = r.pool.Exec(ctx, `
		INSERT INTO users (email, name, password, role)
		VALUES ('admin@skillofied.com', 'Admin User', 'skillofied123', 'admin')
		ON CONFLICT (email) DO NOTHING;
	`)
	if err != nil {
		return fmt.Errorf("seed default user: %w", err)
	}

	return nil
}

// EnsureProfileTable creates the user_profiles table if missing.
func (r *UserRepository) EnsureProfileTable(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS user_profiles (
			user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,

			-- Personal Info
			gender     TEXT NOT NULL DEFAULT '',
			dob        TEXT NOT NULL DEFAULT '',
			whatsapp   TEXT NOT NULL DEFAULT '',
			phone      TEXT NOT NULL DEFAULT '',
			experience TEXT NOT NULL DEFAULT '',

			-- Generic Details
			work_experience          TEXT    NOT NULL DEFAULT '',
			career_gap               TEXT    NOT NULL DEFAULT '',
			current_state            TEXT    NOT NULL DEFAULT '',
			current_city             TEXT    NOT NULL DEFAULT '',
			preferred_locations      TEXT[]  NOT NULL DEFAULT '{}',
			github_link              TEXT    NOT NULL DEFAULT '',
			linkedin_link            TEXT    NOT NULL DEFAULT '',
			is_working_professional  BOOLEAN NOT NULL DEFAULT FALSE,
			resume_name              TEXT    NOT NULL DEFAULT '',

			-- 10th Grade
			edu10_school_name     TEXT NOT NULL DEFAULT '',
			edu10_year_of_passout TEXT NOT NULL DEFAULT '',
			edu10_marks_percent   TEXT NOT NULL DEFAULT '',

			-- 12th / PUC / Intermediate / Diploma
			edu12_school_name     TEXT NOT NULL DEFAULT '',
			edu12_year_of_passout TEXT NOT NULL DEFAULT '',
			edu12_marks_percent   TEXT NOT NULL DEFAULT '',

			-- UG Detail
			ug_university_roll_no TEXT NOT NULL DEFAULT '',
			ug_college_name       TEXT NOT NULL DEFAULT '',
			ug_course_name        TEXT NOT NULL DEFAULT '',
			ug_branch             TEXT NOT NULL DEFAULT '',
			ug_year_of_passout    TEXT NOT NULL DEFAULT '',
			ug_marks_percent      TEXT NOT NULL DEFAULT '',
			ug_cgpa               TEXT NOT NULL DEFAULT '',
			ug_active_backlogs    TEXT NOT NULL DEFAULT '',

			-- PG Detail
			pg_has_certificate BOOLEAN NOT NULL DEFAULT FALSE,

			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
	`)
	if err != nil {
		return fmt.Errorf("create user_profiles table: %w", err)
	}
	return nil
}

// GetProfile retrieves the profile for the given user ID.
// Returns nil, nil when no profile row exists yet.
func (r *UserRepository) GetProfile(ctx context.Context, userID string) (*userv1.UserProfile, error) {
	p := &userv1.UserProfile{UserID: userID}
	err := r.pool.QueryRow(ctx, `
		SELECT
			gender, dob, whatsapp, phone, experience,
			work_experience, career_gap, current_state, current_city,
			preferred_locations, github_link, linkedin_link,
			is_working_professional, resume_name,
			edu10_school_name, edu10_year_of_passout, edu10_marks_percent,
			edu12_school_name, edu12_year_of_passout, edu12_marks_percent,
			ug_university_roll_no, ug_college_name, ug_course_name, ug_branch,
			ug_year_of_passout, ug_marks_percent, ug_cgpa, ug_active_backlogs,
			pg_has_certificate
		FROM user_profiles
		WHERE user_id = $1
	`, userID).Scan(
		&p.Gender, &p.Dob, &p.Whatsapp, &p.Phone, &p.Experience,
		&p.WorkExperience, &p.CareerGap, &p.CurrentState, &p.CurrentCity,
		&p.PreferredLocations, &p.GithubLink, &p.LinkedinLink,
		&p.IsWorkingProfessional, &p.ResumeName,
		&p.Edu10SchoolName, &p.Edu10YearOfPassout, &p.Edu10MarksPercent,
		&p.Edu12SchoolName, &p.Edu12YearOfPassout, &p.Edu12MarksPercent,
		&p.UGUniversityRollNo, &p.UGCollegeName, &p.UGCourseName, &p.UGBranch,
		&p.UGYearOfPassout, &p.UGMarksPercent, &p.UGCGPA, &p.UGActiveBacklogs,
		&p.PGHasCertificate,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // no profile yet – not an error
		}
		return nil, fmt.Errorf("get profile: %w", err)
	}
	return p, nil
}

// UpsertProfile inserts or updates a user's profile row.
func (r *UserRepository) UpsertProfile(ctx context.Context, p *userv1.UserProfile) error {
	// Guard: nil slice → empty slice so the NOT NULL TEXT[] column is satisfied
	if p.PreferredLocations == nil {
		p.PreferredLocations = []string{}
	}

	_, err := r.pool.Exec(ctx, `
		INSERT INTO user_profiles (
			user_id,
			gender, dob, whatsapp, phone, experience,
			work_experience, career_gap, current_state, current_city,
			preferred_locations, github_link, linkedin_link,
			is_working_professional, resume_name,
			edu10_school_name, edu10_year_of_passout, edu10_marks_percent,
			edu12_school_name, edu12_year_of_passout, edu12_marks_percent,
			ug_university_roll_no, ug_college_name, ug_course_name, ug_branch,
			ug_year_of_passout, ug_marks_percent, ug_cgpa, ug_active_backlogs,
			pg_has_certificate
		) VALUES (
			$1,
			$2, $3, $4, $5, $6,
			$7, $8, $9, $10,
			$11, $12, $13,
			$14, $15,
			$16, $17, $18,
			$19, $20, $21,
			$22, $23, $24, $25,
			$26, $27, $28, $29,
			$30
		)
		ON CONFLICT (user_id) DO UPDATE SET
			gender                  = EXCLUDED.gender,
			dob                     = EXCLUDED.dob,
			whatsapp                = EXCLUDED.whatsapp,
			phone                   = EXCLUDED.phone,
			experience              = EXCLUDED.experience,
			work_experience         = EXCLUDED.work_experience,
			career_gap              = EXCLUDED.career_gap,
			current_state           = EXCLUDED.current_state,
			current_city            = EXCLUDED.current_city,
			preferred_locations     = EXCLUDED.preferred_locations,
			github_link             = EXCLUDED.github_link,
			linkedin_link           = EXCLUDED.linkedin_link,
			is_working_professional = EXCLUDED.is_working_professional,
			resume_name             = EXCLUDED.resume_name,
			edu10_school_name       = EXCLUDED.edu10_school_name,
			edu10_year_of_passout   = EXCLUDED.edu10_year_of_passout,
			edu10_marks_percent     = EXCLUDED.edu10_marks_percent,
			edu12_school_name       = EXCLUDED.edu12_school_name,
			edu12_year_of_passout   = EXCLUDED.edu12_year_of_passout,
			edu12_marks_percent     = EXCLUDED.edu12_marks_percent,
			ug_university_roll_no   = EXCLUDED.ug_university_roll_no,
			ug_college_name         = EXCLUDED.ug_college_name,
			ug_course_name          = EXCLUDED.ug_course_name,
			ug_branch               = EXCLUDED.ug_branch,
			ug_year_of_passout      = EXCLUDED.ug_year_of_passout,
			ug_marks_percent        = EXCLUDED.ug_marks_percent,
			ug_cgpa                 = EXCLUDED.ug_cgpa,
			ug_active_backlogs      = EXCLUDED.ug_active_backlogs,
			pg_has_certificate      = EXCLUDED.pg_has_certificate,
			updated_at              = now();
	`,
		p.UserID,
		p.Gender, p.Dob, p.Whatsapp, p.Phone, p.Experience,
		p.WorkExperience, p.CareerGap, p.CurrentState, p.CurrentCity,
		p.PreferredLocations, p.GithubLink, p.LinkedinLink,
		p.IsWorkingProfessional, p.ResumeName,
		p.Edu10SchoolName, p.Edu10YearOfPassout, p.Edu10MarksPercent,
		p.Edu12SchoolName, p.Edu12YearOfPassout, p.Edu12MarksPercent,
		p.UGUniversityRollNo, p.UGCollegeName, p.UGCourseName, p.UGBranch,
		p.UGYearOfPassout, p.UGMarksPercent, p.UGCGPA, p.UGActiveBacklogs,
		p.PGHasCertificate,
	)
	if err != nil {
		return fmt.Errorf("upsert profile: %w", err)
	}
	return nil
}

