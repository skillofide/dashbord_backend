-- ─── Companies ────────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS companies (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT        NOT NULL,
    slug       TEXT        NOT NULL UNIQUE,
    logo_url   TEXT        NOT NULL DEFAULT '',
    website    TEXT        NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Which recruiter users may administer a company's drives. Every recruiter
-- query joins this table; a client-supplied company_id is never trusted.
CREATE TABLE IF NOT EXISTS company_members (
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    user_id    UUID NOT NULL,
    role       TEXT NOT NULL DEFAULT 'recruiter' CHECK (role IN ('recruiter', 'owner')),
    PRIMARY KEY (company_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_company_members_user ON company_members(user_id);

-- ─── MCQ / descriptive question bank ──────────────────────────────────────────
-- Coding questions are NOT duplicated here: section_questions references
-- problem-service problems by id.

CREATE TABLE IF NOT EXISTS mcq_questions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id  UUID        REFERENCES companies(id) ON DELETE CASCADE, -- NULL = platform bank
    topic       TEXT        NOT NULL DEFAULT 'General',
    difficulty  TEXT        NOT NULL DEFAULT 'Medium' CHECK (difficulty IN ('Easy', 'Medium', 'Hard')),
    body        TEXT        NOT NULL,
    kind        TEXT        NOT NULL DEFAULT 'single' CHECK (kind IN ('single', 'multiple', 'numeric')),
    explanation TEXT        NOT NULL DEFAULT '',
    is_active   BOOLEAN     NOT NULL DEFAULT true,
    created_by  UUID,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_mcq_company    ON mcq_questions(company_id);
CREATE INDEX IF NOT EXISTS idx_mcq_topic      ON mcq_questions(topic);
CREATE INDEX IF NOT EXISTS idx_mcq_difficulty ON mcq_questions(difficulty);

-- Options live server-side and the answer key never crosses the service
-- boundary toward a student-facing response.
CREATE TABLE IF NOT EXISTS mcq_options (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_id UUID    NOT NULL REFERENCES mcq_questions(id) ON DELETE CASCADE,
    body        TEXT    NOT NULL,
    is_correct  BOOLEAN NOT NULL DEFAULT false,
    order_index INT     NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_mcq_options_question ON mcq_options(question_id);

-- ─── Test definition ──────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS assessments (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id        UUID         REFERENCES companies(id) ON DELETE CASCADE, -- NULL = practice test
    title             TEXT         NOT NULL,
    description       TEXT         NOT NULL DEFAULT '',
    purpose           TEXT         NOT NULL DEFAULT 'practice' CHECK (purpose IN ('practice', 'hiring')),
    duration_minutes  INT          NOT NULL DEFAULT 60,
    total_marks       INT          NOT NULL DEFAULT 0,   -- derived on publish, cached
    passing_marks     INT          NOT NULL DEFAULT 0,
    negative_marking  NUMERIC(4,2) NOT NULL DEFAULT 0,
    shuffle_questions BOOLEAN      NOT NULL DEFAULT true,
    shuffle_options   BOOLEAN      NOT NULL DEFAULT true,
    allow_backtrack   BOOLEAN      NOT NULL DEFAULT true,
    reveal_results    BOOLEAN      NOT NULL DEFAULT true,
    proctoring        JSONB        NOT NULL DEFAULT '{}',
    status            TEXT         NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    opens_at          TIMESTAMPTZ,
    closes_at         TIMESTAMPTZ,
    max_attempts      INT          NOT NULL DEFAULT 1,
    created_by        UUID         NOT NULL,
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ  NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_assessments_company ON assessments(company_id);
CREATE INDEX IF NOT EXISTS idx_assessments_status  ON assessments(status, purpose);

-- A test is an ordered list of sections; a mixed-pattern test is simply an
-- 'mcq' section plus a 'coding' section.
CREATE TABLE IF NOT EXISTS assessment_sections (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_id    UUID    NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    title            TEXT    NOT NULL,
    kind             TEXT    NOT NULL CHECK (kind IN ('mcq', 'coding', 'descriptive')),
    order_index      INT     NOT NULL DEFAULT 0,
    duration_minutes INT,                     -- NULL = shares the global test timer
    cutoff_marks     INT     NOT NULL DEFAULT 0,
    -- Random draw from the bank instead of a fixed question list.
    pick_count       INT,                     -- NULL = use the explicit list
    pick_topic       TEXT    NOT NULL DEFAULT '',
    pick_difficulty  TEXT    NOT NULL DEFAULT '',
    pick_marks       INT     NOT NULL DEFAULT 1,  -- marks per bank-drawn question
    partial_credit   BOOLEAN NOT NULL DEFAULT false
);
CREATE INDEX IF NOT EXISTS idx_sections_assessment ON assessment_sections(assessment_id, order_index);

CREATE TABLE IF NOT EXISTS section_questions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    section_id      UUID NOT NULL REFERENCES assessment_sections(id) ON DELETE CASCADE,
    mcq_question_id UUID REFERENCES mcq_questions(id) ON DELETE RESTRICT,
    problem_id      UUID,   -- FK-by-convention into problem-service
    marks           INT  NOT NULL DEFAULT 1,
    order_index     INT  NOT NULL DEFAULT 0,
    CHECK (num_nonnulls(mcq_question_id, problem_id) = 1)
);
CREATE INDEX IF NOT EXISTS idx_section_questions ON section_questions(section_id, order_index);

-- ─── Invitations (hiring drives) ──────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS assessment_invites (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_id UUID NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    email         TEXT NOT NULL,
    user_id       UUID,                    -- resolved on first claim
    token         TEXT NOT NULL UNIQUE,
    status        TEXT NOT NULL DEFAULT 'invited'
                  CHECK (status IN ('invited', 'opened', 'started', 'submitted', 'expired')),
    expires_at    TIMESTAMPTZ,
    sent_at       TIMESTAMPTZ,
    UNIQUE (assessment_id, email)
);
CREATE INDEX IF NOT EXISTS idx_invites_email ON assessment_invites(lower(email));

-- ─── Attempts ─────────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS attempts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_id   UUID         NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    user_id         UUID         NOT NULL,
    invite_id       UUID         REFERENCES assessment_invites(id) ON DELETE SET NULL,
    attempt_no      INT          NOT NULL DEFAULT 1,
    status          TEXT         NOT NULL DEFAULT 'in_progress'
                    CHECK (status IN ('in_progress', 'submitted', 'evaluating', 'evaluated', 'disqualified', 'expired')),
    seed            BIGINT       NOT NULL,
    started_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
    expires_at      TIMESTAMPTZ  NOT NULL,  -- server-authoritative deadline
    submitted_at    TIMESTAMPTZ,
    evaluated_at    TIMESTAMPTZ,
    score           NUMERIC(8,2) NOT NULL DEFAULT 0,
    max_score       NUMERIC(8,2) NOT NULL DEFAULT 0,
    section_scores  JSONB        NOT NULL DEFAULT '{}',
    integrity_score NUMERIC(5,2) NOT NULL DEFAULT 100,
    UNIQUE (assessment_id, user_id, attempt_no)
);
CREATE INDEX IF NOT EXISTS idx_attempts_assessment ON attempts(assessment_id, status);
CREATE INDEX IF NOT EXISTS idx_attempts_user       ON attempts(user_id);
-- Drives the expiry sweeper.
CREATE INDEX IF NOT EXISTS idx_attempts_live       ON attempts(expires_at) WHERE status = 'in_progress';

-- The materialized paper. Frozen at StartAttempt so a bank edit mid-drive
-- cannot change what a candidate already saw, and so the same seed reproduces
-- the paper for dispute review.
CREATE TABLE IF NOT EXISTS attempt_questions (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attempt_id       UUID         NOT NULL REFERENCES attempts(id) ON DELETE CASCADE,
    section_id       UUID         NOT NULL,
    kind             TEXT         NOT NULL CHECK (kind IN ('mcq', 'coding', 'descriptive')),
    mcq_question_id  UUID,
    problem_id       UUID,
    order_index      INT          NOT NULL,
    marks            NUMERIC(6,2) NOT NULL DEFAULT 1,
    option_order     UUID[]       NOT NULL DEFAULT '{}',
    -- answer state
    selected_options UUID[]       NOT NULL DEFAULT '{}',
    text_answer      TEXT,
    submission_id    UUID,
    language         TEXT         NOT NULL DEFAULT '',
    code             TEXT         NOT NULL DEFAULT '',
    awarded_marks    NUMERIC(6,2),
    grading_status   TEXT         NOT NULL DEFAULT 'ungraded'
                     CHECK (grading_status IN ('ungraded', 'pending', 'graded', 'manual_review')),
    time_spent_ms    BIGINT       NOT NULL DEFAULT 0,
    visited          BOOLEAN      NOT NULL DEFAULT false,
    marked_review    BOOLEAN      NOT NULL DEFAULT false,
    UNIQUE (attempt_id, order_index)
);
CREATE INDEX IF NOT EXISTS idx_attempt_questions_attempt    ON attempt_questions(attempt_id);
CREATE INDEX IF NOT EXISTS idx_attempt_questions_submission ON attempt_questions(submission_id);

-- Every graded code submission inside a test. Scoring takes the BEST result
-- per question, not the last, so this history is load-bearing.
CREATE TABLE IF NOT EXISTS attempt_submissions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attempt_question_id UUID        NOT NULL REFERENCES attempt_questions(id) ON DELETE CASCADE,
    submission_id       UUID        NOT NULL UNIQUE,
    language            TEXT        NOT NULL,
    code                TEXT        NOT NULL,
    passed_count        INT,
    total_count         INT,
    status              TEXT        NOT NULL DEFAULT 'Pending',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_attempt_submissions_question ON attempt_submissions(attempt_question_id);

-- ─── Proctoring ───────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS proctor_events (
    id          BIGSERIAL PRIMARY KEY,
    attempt_id  UUID        NOT NULL REFERENCES attempts(id) ON DELETE CASCADE,
    kind        TEXT        NOT NULL,
    detail      TEXT        NOT NULL DEFAULT '',
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_proctor_events_attempt ON proctor_events(attempt_id, occurred_at);

-- ─── Shortlisting ─────────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS shortlists (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assessment_id UUID        NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    name          TEXT        NOT NULL,
    criteria      JSONB       NOT NULL DEFAULT '{}',
    created_by    UUID        NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_shortlists_assessment ON shortlists(assessment_id);

CREATE TABLE IF NOT EXISTS shortlist_entries (
    shortlist_id UUID NOT NULL REFERENCES shortlists(id) ON DELETE CASCADE,
    attempt_id   UUID NOT NULL REFERENCES attempts(id) ON DELETE CASCADE,
    user_id      UUID NOT NULL,
    rank         INT  NOT NULL,
    decision     TEXT NOT NULL DEFAULT 'shortlisted'
                 CHECK (decision IN ('shortlisted', 'rejected', 'on_hold', 'hired')),
    notes        TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (shortlist_id, attempt_id)
);
