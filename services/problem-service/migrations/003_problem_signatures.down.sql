DROP TABLE IF EXISTS reference_solutions;
ALTER TABLE problems DROP CONSTRAINT IF EXISTS problems_io_mode_check;
ALTER TABLE problems DROP COLUMN IF EXISTS io_mode;
DROP TABLE IF EXISTS problem_signatures;
