DROP INDEX IF EXISTS idx_problems_private;
ALTER TABLE problems DROP COLUMN IF EXISTS owner_company_id;
ALTER TABLE problems DROP COLUMN IF EXISTS is_private;
