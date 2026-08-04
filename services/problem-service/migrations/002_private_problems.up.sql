-- Company-owned coding questions used in hiring drives must never surface in
-- the public practice listings. ListProblems filters on is_private = false.
ALTER TABLE problems ADD COLUMN IF NOT EXISTS is_private       BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE problems ADD COLUMN IF NOT EXISTS owner_company_id UUID;

CREATE INDEX IF NOT EXISTS idx_problems_private ON problems(is_private) WHERE is_private = false;
