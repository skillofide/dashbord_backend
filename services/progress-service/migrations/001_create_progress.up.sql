-- User progress summary (one row per user, updated on each submission)
CREATE TABLE IF NOT EXISTS user_progress (
    user_id         TEXT    PRIMARY KEY,
    total_solved    INT     NOT NULL DEFAULT 0,
    total_attempted INT     NOT NULL DEFAULT 0,
    easy_solved     INT     NOT NULL DEFAULT 0,
    medium_solved   INT     NOT NULL DEFAULT 0,
    hard_solved     INT     NOT NULL DEFAULT 0,
    current_streak  INT     NOT NULL DEFAULT 0,
    longest_streak  INT     NOT NULL DEFAULT 0,
    total_xp        INT     NOT NULL DEFAULT 0,
    last_active_date DATE,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Per-problem-per-user status (one row per user+problem)
CREATE TABLE IF NOT EXISTS problem_progress (
    user_id    TEXT        NOT NULL,
    problem_id UUID        NOT NULL,
    set_id     UUID,
    -- Status: Solved | InProgress | Unsolved
    status     TEXT        NOT NULL DEFAULT 'Unsolved',
    solved_at  TIMESTAMPTZ,
    attempts   INT         NOT NULL DEFAULT 0,
    runtime_ms BIGINT      NOT NULL DEFAULT 0,
    memory_kb  BIGINT      NOT NULL DEFAULT 0,
    language   TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, problem_id)
);

-- Per-set progress (denormalized for fast dashboard queries)
CREATE TABLE IF NOT EXISTS set_progress (
    user_id TEXT NOT NULL,
    set_id  UUID NOT NULL,
    solved  INT  NOT NULL DEFAULT 0,
    total   INT  NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, set_id)
);

-- Activity log (one row per day per user — used for streak calculation)
CREATE TABLE IF NOT EXISTS activity_log (
    user_id       TEXT NOT NULL,
    activity_date DATE NOT NULL,
    PRIMARY KEY (user_id, activity_date)
);

CREATE INDEX IF NOT EXISTS idx_problem_progress_user   ON problem_progress(user_id);
CREATE INDEX IF NOT EXISTS idx_problem_progress_set    ON problem_progress(set_id);
CREATE INDEX IF NOT EXISTS idx_activity_log_user       ON activity_log(user_id, activity_date DESC);
