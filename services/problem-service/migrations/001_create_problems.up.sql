-- Create practice_sets table
CREATE TABLE IF NOT EXISTS practice_sets (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title      TEXT        NOT NULL,
    level      TEXT        NOT NULL, -- Beginner | Intermediate | Advanced
    level_color TEXT       NOT NULL DEFAULT '#6366f1',
    bg_color   TEXT        NOT NULL DEFAULT '#1e1b4b',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Create problems table
CREATE TABLE IF NOT EXISTS problems (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug        TEXT        NOT NULL UNIQUE,
    title       TEXT        NOT NULL,
    difficulty  TEXT        NOT NULL CHECK (difficulty IN ('Easy', 'Medium', 'Hard')),
    topic       TEXT        NOT NULL,
    xp          INT         NOT NULL DEFAULT 0,
    statement   TEXT        NOT NULL,
    set_id      UUID        REFERENCES practice_sets(id) ON DELETE SET NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Constraints list (ordered)
CREATE TABLE IF NOT EXISTS problem_constraints (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    problem_id      UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    constraint_text TEXT NOT NULL,
    order_index     INT  NOT NULL DEFAULT 0
);

-- Tags
CREATE TABLE IF NOT EXISTS problem_tags (
    problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    tag        TEXT NOT NULL,
    PRIMARY KEY (problem_id, tag)
);

-- Examples (input/output/explanation)
CREATE TABLE IF NOT EXISTS examples (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    problem_id  UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    input       TEXT NOT NULL,
    output      TEXT NOT NULL,
    explanation TEXT NOT NULL DEFAULT '',
    order_index INT  NOT NULL DEFAULT 0
);

-- Hints
CREATE TABLE IF NOT EXISTS hints (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    problem_id  UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    order_index INT  NOT NULL DEFAULT 0,
    title       TEXT NOT NULL,
    body        TEXT NOT NULL
);

-- Starter codes (one row per problem)
CREATE TABLE IF NOT EXISTS starter_codes (
    problem_id  UUID PRIMARY KEY REFERENCES problems(id) ON DELETE CASCADE,
    javascript  TEXT NOT NULL DEFAULT '',
    python      TEXT NOT NULL DEFAULT '',
    java        TEXT NOT NULL DEFAULT '',
    cpp         TEXT NOT NULL DEFAULT '',
    go          TEXT NOT NULL DEFAULT ''
);

-- Test cases (hidden flag controls Run vs Submit visibility)
CREATE TABLE IF NOT EXISTS test_cases (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    problem_id      UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    input           TEXT NOT NULL,
    expected_output TEXT NOT NULL,
    is_hidden       BOOLEAN NOT NULL DEFAULT false,
    time_limit_ms   INT     NOT NULL DEFAULT 2000,
    memory_limit_mb INT     NOT NULL DEFAULT 256,
    order_index     INT     NOT NULL DEFAULT 0
);

-- Per-user problem status (denormalized from progress-service for fast reads)
CREATE TABLE IF NOT EXISTS problem_user_status (
    user_id    TEXT NOT NULL,
    problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    status     TEXT NOT NULL DEFAULT 'Unsolved' CHECK (status IN ('Solved', 'InProgress', 'Unsolved')),
    solved_at  TIMESTAMPTZ,
    attempts   INT  NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, problem_id)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_problems_set_id     ON problems(set_id);
CREATE INDEX IF NOT EXISTS idx_problems_difficulty  ON problems(difficulty);
CREATE INDEX IF NOT EXISTS idx_problems_topic       ON problems(topic);
CREATE INDEX IF NOT EXISTS idx_test_cases_problem   ON test_cases(problem_id);
CREATE INDEX IF NOT EXISTS idx_pus_user_id          ON problem_user_status(user_id);
