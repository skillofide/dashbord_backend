-- Class-design problems: LRUCache, MinStack, Trie, TicTacToe.
--
-- These are not single-entry-point problems. The learner builds a type with a
-- constructor and several methods, and it is graded by driving a sequence of
-- calls against one instance — state carried between them is the whole point.
-- A `put` followed by a `get` cannot be expressed as one function call, so
-- forcing them into a function signature would change what they teach.
--
-- The test input is LeetCode's own two-line form:
--
--   ["LRUCache","put","put","get","put","get"]     the call sequence
--   [[2],[1,1],[2,2],[1],[3,3],[2]]                the arguments for each
--
-- and the expectation is the list of return values, with null where a method
-- returns nothing:
--
--   [null,null,null,1,null,-1]

ALTER TABLE problem_signatures
    ADD COLUMN IF NOT EXISTS kind TEXT NOT NULL DEFAULT 'function';

ALTER TABLE problem_signatures DROP CONSTRAINT IF EXISTS problem_signatures_kind_check;
ALTER TABLE problem_signatures ADD CONSTRAINT problem_signatures_kind_check
    CHECK (kind IN ('function', 'class'));

-- For kind = 'class': entry_point is the class name, params are the
-- constructor's, and methods lists everything the call sequence may name.
--   [{"name":"put","params":[{"name":"key","type":"int"},
--                            {"name":"value","type":"int"}],
--     "returnType":"void"}, ...]
ALTER TABLE problem_signatures
    ADD COLUMN IF NOT EXISTS methods JSONB NOT NULL DEFAULT '[]'::jsonb;

ALTER TABLE problem_signatures DROP CONSTRAINT IF EXISTS problem_signatures_methods_is_array;
ALTER TABLE problem_signatures ADD CONSTRAINT problem_signatures_methods_is_array
    CHECK (jsonb_typeof(methods) = 'array');

-- A class problem must list at least one method; a function problem must not
-- list any. Getting this wrong produces a driver that compiles and grades
-- nothing.
ALTER TABLE problem_signatures DROP CONSTRAINT IF EXISTS problem_signatures_methods_match_kind;
ALTER TABLE problem_signatures ADD CONSTRAINT problem_signatures_methods_match_kind
    CHECK (
      (kind = 'class'    AND jsonb_array_length(methods) > 0) OR
      (kind = 'function' AND jsonb_array_length(methods) = 0)
    );
