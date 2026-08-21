ALTER TABLE problem_signatures DROP CONSTRAINT IF EXISTS problem_signatures_methods_match_kind;
ALTER TABLE problem_signatures DROP CONSTRAINT IF EXISTS problem_signatures_methods_is_array;
ALTER TABLE problem_signatures DROP COLUMN IF EXISTS methods;
ALTER TABLE problem_signatures DROP CONSTRAINT IF EXISTS problem_signatures_kind_check;
ALTER TABLE problem_signatures DROP COLUMN IF EXISTS kind;
