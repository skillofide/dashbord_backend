CREATE TABLE IF NOT EXISTS submissions (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       TEXT        NOT NULL,
    problem_id    UUID        NOT NULL,
    language      TEXT        NOT NULL,
    code          TEXT        NOT NULL,
    -- Status: Pending | Running | Accepted | WrongAnswer | TimeLimitExceeded | MemoryLimitExceeded | RuntimeError | CompileError
    status        TEXT        NOT NULL DEFAULT 'Pending',
    runtime_ms    BIGINT      NOT NULL DEFAULT 0,
    memory_kb     BIGINT      NOT NULL DEFAULT 0,
    compile_error TEXT,
    test_results  JSONB       NOT NULL DEFAULT '[]',
    submitted_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_submissions_user_id    ON submissions(user_id);
CREATE INDEX IF NOT EXISTS idx_submissions_problem_id ON submissions(problem_id);
CREATE INDEX IF NOT EXISTS idx_submissions_status     ON submissions(status);
CREATE INDEX IF NOT EXISTS idx_submissions_submitted  ON submissions(submitted_at DESC);
