package repository

import (
	"context"
	"fmt"
	"time"

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
		VALUES ('admin@knovate.com', 'Admin User', 'knovate123', 'admin'),
		       ('admin@skillofied.com', 'Admin User', 'skillofied123', 'admin')
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

// EnsureUserCoursesTable creates the user_courses table if missing.
func (r *UserRepository) EnsureUserCoursesTable(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS user_courses (
			id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			course_id  TEXT NOT NULL,
			granted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			UNIQUE (user_id, course_id)
		);
	`)
	if err != nil {
		return fmt.Errorf("create user_courses table: %w", err)
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

// CheckUserCourseAccess verifies if a user has access to a specific course.
func (r *UserRepository) CheckUserCourseAccess(ctx context.Context, userID, courseID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM user_courses
			WHERE user_id = $1 AND course_id = $2
		)
	`, userID, courseID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check course access: %w", err)
	}
	return exists, nil
}

// EnsureQuizTables creates the quiz database tables if they do not exist and seeds correct answers for Module 1.
func (r *UserRepository) EnsureQuizTables(ctx context.Context) error {
	_, err := r.pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS quiz_keys (
			module_id TEXT NOT NULL,
			question_id INT NOT NULL,
			correct_answer TEXT NOT NULL,
			PRIMARY KEY (module_id, question_id)
		);

		CREATE TABLE IF NOT EXISTS user_quiz_attempts (
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			module_id TEXT NOT NULL,
			score INT NOT NULL,
			total_questions INT NOT NULL,
			completed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			PRIMARY KEY (user_id, module_id)
		);
	`)
	if err != nil {
		return fmt.Errorf("create quiz tables: %w", err)
	}

	// Seed correct answers for Module 1 (20 questions)
	answers := []struct {
		moduleID   string
		questionID int
		correctAns string
	}{
		{"m1", 1, "D. James Gosling"},
		{"m1", 2, "C. Oak"},
		{"m1", 3, "A. Write Once, Run Anywhere (WORA)"},
		{"m1", 4, "B. JRE and development tools like 'javac'"},
		{"m1", 5, "B. Java Virtual Machine (JVM)"},
		{"m1", 6, "C. .class"},
		{"m1", 7, "C. Java Bytecode is platform-independent, but the JVM is platform-dependent."},
		{"m1", 8, "A. Garbage Collection"},
		{"m1", 9, "A. public static void main(String[] args)"},
		{"m1", 10, "C. Welcome.java"},
		{"m1", 11, "B. Providing an all-in-one text editor, build automation tool, and debugger"},
		{"m1", 12, "A. javac Test.java"},
		{"m1", 13, "B. The method does not return any value when it finishes executing."},
		{"m1", 14, "A. Strong type checking and exception handling mechanisms"},
		{"m1", 15, "A. java Demo"},
		{"m1", 16, "C. PATH"},
		{"m1", 17, "D. System.out.println(\"Hello World\");"},
		{"m1", 18, "B. To enable the JVM to call the method without creating an instance of the class first"},
		{"m1", 19, "B. The JDK is a superset that includes the complete JRE plus development tools."},
		{"m1", 20, "A. They mark the beginning of a single-line text comment."},
		{"m2", 1, "D. _variable$5"},
		{"m2", 2, "A. String"},
		{"m2", 3, "C. Widening casting happens automatically; narrowing casting must be done manually."},
		{"m2", 4, "B. /** Documentation comment */"},
		{"m2", 5, "B. next() reads input up to the next whitespace delimiter, while nextLine() reads the entire line until a newline character."},
		{"m2", 6, "A. %f"},
		{"m2", 7, "B. They skip evaluating the second condition if the overall result is already determined by the first condition."},
		{"m2", 8, "B. a=7, b=12"},
		{"m2", 9, "B. 2.0"},
		{"m2", 10, "D. 9"},
		{"m2", 11, "C. Output: 1020"},
		{"m2", 12, "C. 30 Output"},
		{"m2", 13, "A. -1"},
		{"m2", 14, "B. 10"},
		{"m2", 15, "D. -128"},
		{"m2", 16, "D. 5.68"},
		{"m2", 17, "A. n1=20, n2=10"},
		{"m2", 18, "D. 30"},
		{"m2", 19, "A. true"},
		{"m2", 20, "A. B"},
		{"m3", 1, "C. double"},
		{"m3", 2, "C. The program falls through, executing subsequent case blocks sequentially until a break or the end of the switch is encountered."},
		{"m3", 3, "A. The inner 'if' condition is evaluated only if the outer 'if' condition evaluates to true."},
		{"m3", 4, "A. 3"},
		{"m3", 5, "C. Passed"},
		{"m3", 6, "C. Block 2"},
		{"m3", 7, "B. Set Go Done"},
		{"m3", 8, "B. High"},
		{"m3", 9, "B. Divisible"},
		{"m3", 10, "D. 10"},
		{"m3", 11, "D. Compilation Error"},
		{"m3", 12, "C. Off"},
		{"m3", 13, "A. 20"},
		{"m3", 14, "C. Allowed"},
		{"m3", 15, "B. Five"},
		{"m3", 16, "D. Default One"},
		{"m3", 17, "A. Warm"},
		{"m3", 18, "B. 6"},
		{"m3", 19, "B. 8"},
		{"m3", 20, "D. Point 3"},
		{"m4", 1, "C. do-while loop"},
		{"m4", 2, "B. continue"},
		{"m4", 3, "B. Infinite loop"},
		{"m4", 4, "C. O(N^2)"},
		{"m4", 5, "C. semicolon (;)"},
		{"m5", 1, "C. void"},
		{"m5", 2, "B. StackOverflowError"},
		{"m5", 3, "B. No"},
		{"m5", 4, "B. Pass by value"},
		{"m5", 5, "C. Overriding resolution at runtime"},
		{"m6", 1, "B. 2"},
		{"m6", 2, "B. length"},
		{"m6", 3, "C. ArrayIndexOutOfBoundsException"},
		{"m6", 4, "B. Arrays.sort()"},
		{"m6", 5, "B. No"},
		{"m7", 1, "C. String Constant Pool (SCP)"},
		{"m7", 2, "B. str1.equals(str2)"},
		{"m7", 3, "C. StringBuilder"},
		{"m7", 4, "B. \"bc\""},
		{"m7", 5, "A. For security, caching, and thread safety"},
		{"m8", 1, "D. Abstraction"},
		{"m8", 2, "B. this"},
		{"m8", 3, "B. No"},
		{"m8", 4, "C. implements"},
		{"m8", 5, "B. Method Overloading"},
		{"m9", 1, "C. finally"},
		{"m9", 2, "B. throws"},
		{"m9", 3, "A. Throwable"},
		{"m9", 4, "B. Unchecked"},
		{"m9", 5, "B. Extend Exception"},
		{"m10", 1, "C. HashSet"},
		{"m10", 2, "B. HashMap"},
		{"m10", 3, "C. TreeSet"},
		{"m10", 4, "B. map.containsKey()"},
		{"m10", 5, "B. ConcurrentModificationException"},
		{"m11", 1, "B. BufferedReader"},
		{"m11", 2, "B. Try-with-resources"},
		{"m11", 3, "B. file.exists()"},
		{"m11", 4, "B. new FileWriter(\"file.txt\", true)"},
		{"m11", 5, "C. java.io"},
		{"m12", 1, "B. start()"},
		{"m12", 2, "C. synchronized"},
		{"m12", 3, "B. Java supports multiple interface implementations but only single class inheritance"},
		{"m12", 4, "B. Runnable"},
		{"m12", 5, "B. Executors"},
		{"m13", 1, "B. @FunctionalInterface"},
		{"m13", 2, "C. ::"},
		{"m13", 3, "C. reduce()"},
		{"m13", 4, "A. Optional.empty()"},
		{"m13", 5, "A. Intermediate"},
		{"m14", 1, "C. ResultSet"},
		{"m14", 2, "B. They prevent SQL Injection and cache query execution plans"},
		{"m14", 3, "B. executeQuery()"},
		{"m14", 4, "B. jdbc:mysql://..."},
		{"m14", 5, "B. 1"},
		{"m15", 1, "B. O(log N)"},
		{"m15", 2, "B. Stack"},
		{"m15", 3, "B. Breadth First Search (BFS)"},
		{"m15", 4, "C. O(N^2)"},
		{"m15", 5, "B. Queue"},
		{"m16", 1, "B. @RestController"},
		{"m16", 2, "B. Spring Initializr"},
		{"m16", 3, "B. @Autowired"},
		{"m16", 4, "C. Tomcat"},
		{"m16", 5, "B. JpaRepository"},
		{"m17", 1, "B. @Id"},
		{"m17", 2, "B. @Entity"},
		{"m17", 3, "B. @Valid"},
		{"m17", 4, "C. @RestControllerAdvice"},
		{"m17", 5, "B. spring.jpa.hibernate.ddl-auto=update"},
		{"m18", 1, "B. BCryptPasswordEncoder"},
		{"m18", 2, "B. JSON Web Token"},
		{"m18", 3, "B. Authorization Header"},
		{"m18", 4, "B. csrf().disable()"},
		{"m18", 5, "B. Authorization"},
		{"m19", 1, "B. mvn clean package"},
		{"m19", 2, "B. java -jar app.jar"},
		{"m19", 3, "A. docker build"},
		{"m19", 4, "B. target/"},
		{"m19", 5, "C. Passed via Environment Variables"},
	}

	for _, a := range answers {
		_, err = r.pool.Exec(ctx, `
			INSERT INTO quiz_keys (module_id, question_id, correct_answer)
			VALUES ($1, $2, $3)
			ON CONFLICT (module_id, question_id) 
			DO UPDATE SET correct_answer = EXCLUDED.correct_answer;
		`, a.moduleID, a.questionID, a.correctAns)
		if err != nil {
			return fmt.Errorf("seed quiz key: %w", err)
		}
	}

	return nil
}

// SubmitQuiz grades a user's quiz submission and stores the score in the database.
func (r *UserRepository) SubmitQuiz(ctx context.Context, userID, moduleID string, userAnswers []*userv1.QuizAnswer) (int, int, error) {
	// 1. Fetch correct answers for the module
	rows, err := r.pool.Query(ctx, `
		SELECT question_id, correct_answer
		FROM   quiz_keys
		WHERE  module_id = $1
	`, moduleID)
	if err != nil {
		return 0, 0, fmt.Errorf("fetch correct answers: %w", err)
	}
	defer rows.Close()

	keys := make(map[int]string)
	for rows.Next() {
		var qID int
		var ans string
		if err := rows.Scan(&qID, &ans); err != nil {
			return 0, 0, fmt.Errorf("scan correct answer: %w", err)
		}
		keys[qID] = ans
	}

	// 2. Grade
	score := 0
	total := len(keys)
	if total == 0 {
		return 0, 0, fmt.Errorf("no quiz keys found for module %s", moduleID)
	}

	for _, ua := range userAnswers {
		correctAns, ok := keys[ua.QuestionID]
		if ok && correctAns == ua.Answer {
			score++
		}
	}

	// 3. Save attempt
	_, err = r.pool.Exec(ctx, `
		INSERT INTO user_quiz_attempts (user_id, module_id, score, total_questions, completed_at)
		VALUES ($1::uuid, $2, $3, $4, now())
		ON CONFLICT (user_id, module_id)
		DO UPDATE SET score = EXCLUDED.score, total_questions = EXCLUDED.total_questions, completed_at = now();
	`, userID, moduleID, score, total)
	if err != nil {
		return 0, 0, fmt.Errorf("save quiz attempt: %w", err)
	}

	return score, total, nil
}

// GetQuizAttempts fetches all saved quiz scores for a user.
func (r *UserRepository) GetQuizAttempts(ctx context.Context, userID string) ([]*userv1.QuizAttempt, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT module_id, score, total_questions, completed_at
		FROM   user_quiz_attempts
		WHERE  user_id = $1::uuid
		ORDER  BY completed_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("get quiz attempts: %w", err)
	}
	defer rows.Close()

	var attempts []*userv1.QuizAttempt
	for rows.Next() {
		var a userv1.QuizAttempt
		var completedAt time.Time
		if err := rows.Scan(&a.ModuleID, &a.Score, &a.TotalQuestions, &completedAt); err != nil {
			return nil, fmt.Errorf("scan quiz attempt: %w", err)
		}
		a.CompletedAt = completedAt.Format(time.RFC3339)
		attempts = append(attempts, &a)
	}

	return attempts, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Admin repository methods
// ─────────────────────────────────────────────────────────────────────────────

// ImportUserRow holds a single row from the admin Excel import.
type ImportUserRow struct {
	Name      string
	Email     string
	Password  string
	Role      string
	CourseIDs []string
}

// ImportRowResult is the per-row outcome of a bulk import.
type ImportRowResult struct {
	Email   string
	Success bool
	Message string
}

// AdminUserRow is a user record returned to the admin panel.
type AdminUserRow struct {
	ID        string
	Email     string
	Name      string
	Role      string
	CourseIDs []string
	CreatedAt string
}

// BulkUpsertUsers inserts or updates multiple users and grants their course access.
func (r *UserRepository) BulkUpsertUsers(ctx context.Context, rows []ImportUserRow) []ImportRowResult {
	results := make([]ImportRowResult, 0, len(rows))
	for _, row := range rows {
		if row.Email == "" || row.Name == "" || row.Password == "" {
			results = append(results, ImportRowResult{
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

		// Upsert user and return its UUID
		var userID string
		err := r.pool.QueryRow(ctx, `
			INSERT INTO users (email, name, password, role)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (email)
			DO UPDATE SET name = EXCLUDED.name, password = EXCLUDED.password, role = EXCLUDED.role, updated_at = now()
			RETURNING id::text
		`, row.Email, row.Name, row.Password, role).Scan(&userID)
		if err != nil {
			results = append(results, ImportRowResult{
				Email:   row.Email,
				Success: false,
				Message: fmt.Sprintf("upsert user: %v", err),
			})
			continue
		}

		// Grant course access for each listed course
		for _, courseID := range row.CourseIDs {
			if courseID == "" {
				continue
			}
			_, _ = r.pool.Exec(ctx, `
				INSERT INTO user_courses (user_id, course_id)
				VALUES ($1::uuid, $2)
				ON CONFLICT (user_id, course_id) DO NOTHING
			`, userID, courseID)
		}

		results = append(results, ImportRowResult{
			Email:   row.Email,
			Success: true,
			Message: "imported successfully",
		})
	}
	return results
}

// ListAdminUsers returns a paginated list of all users with their assigned course IDs.
func (r *UserRepository) ListAdminUsers(ctx context.Context, page, pageSize int) ([]AdminUserRow, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 50
	}
	offset := (page - 1) * pageSize

	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	rows, err := r.pool.Query(ctx, `
		SELECT u.id::text, u.email, u.name, u.role, u.created_at,
		       COALESCE(array_agg(uc.course_id) FILTER (WHERE uc.course_id IS NOT NULL), '{}') AS course_ids
		FROM users u
		LEFT JOIN user_courses uc ON uc.user_id = u.id
		GROUP BY u.id, u.email, u.name, u.role, u.created_at
		ORDER BY u.created_at DESC
		LIMIT $1 OFFSET $2
	`, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []AdminUserRow
	for rows.Next() {
		var u AdminUserRow
		var createdAt time.Time
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.Role, &createdAt, &u.CourseIDs); err != nil {
			return nil, 0, fmt.Errorf("scan user row: %w", err)
		}
		u.CreatedAt = createdAt.Format(time.RFC3339)
		users = append(users, u)
	}
	return users, total, nil
}

// UpdateUserRole changes the role of an existing user.
func (r *UserRepository) UpdateUserRole(ctx context.Context, userID, role string) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE users SET role = $2, updated_at = now() WHERE id = $1::uuid
	`, userID, role)
	if err != nil {
		return fmt.Errorf("update role: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// AdminGrantCourseAccess grants a user access to a course.
func (r *UserRepository) AdminGrantCourseAccess(ctx context.Context, userID, courseID string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO user_courses (user_id, course_id)
		VALUES ($1::uuid, $2)
		ON CONFLICT (user_id, course_id) DO NOTHING
	`, userID, courseID)
	if err != nil {
		return fmt.Errorf("grant course access: %w", err)
	}
	return nil
}

// AdminRevokeCourseAccess revokes a user's access to a course.
func (r *UserRepository) AdminRevokeCourseAccess(ctx context.Context, userID, courseID string) error {
	_, err := r.pool.Exec(ctx, `
		DELETE FROM user_courses WHERE user_id = $1::uuid AND course_id = $2
	`, userID, courseID)
	if err != nil {
		return fmt.Errorf("revoke course access: %w", err)
	}
	return nil
}

// DeleteUser permanently removes a user (cascades to user_courses, user_profiles, etc.).
func (r *UserRepository) DeleteUser(ctx context.Context, userID string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1::uuid`, userID)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}
