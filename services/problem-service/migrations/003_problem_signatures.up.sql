-- Problem signatures: the typed contract between a problem, the starter code a
-- learner sees, and the driver the judge wraps around their submission.
--
-- Before this table the execution-service guessed both: it regex-scanned the
-- submitted source for the first function definition and hoped the test-case
-- input could be parsed as a JSON argument list. That silently mis-graded every
-- problem whose input spanned more than one line, whose entry point was not the
-- first function in the file, or whose return value was a string. A declared
-- signature removes the guessing from both ends.
CREATE TABLE IF NOT EXISTS problem_signatures (
    problem_id  UUID PRIMARY KEY REFERENCES problems(id) ON DELETE CASCADE,

    -- Name the driver calls, e.g. 'twoSum'. Same across every language; the
    -- code generator applies each language's casing convention.
    entry_point TEXT NOT NULL,

    -- Ordered [{"name":"nums","type":"int[]"}, {"name":"target","type":"int"}].
    -- Order is the argument order and the order of the input lines.
    params      JSONB NOT NULL DEFAULT '[]'::jsonb,

    -- One of the supported type names, or 'void'.
    return_type TEXT NOT NULL,

    -- How the judge compares actual to expected:
    --   exact     deep equality after JSON decode
    --   unordered element order is not significant (top-level sort first)
    --   set       duplicates are not significant either
    --   float     numeric comparison within float_eps
    compare     TEXT NOT NULL DEFAULT 'exact',
    float_eps   DOUBLE PRECISION NOT NULL DEFAULT 1e-6,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT problem_signatures_compare_check
        CHECK (compare IN ('exact', 'unordered', 'set', 'float')),
    CONSTRAINT problem_signatures_params_is_array
        CHECK (jsonb_typeof(params) = 'array')
);

-- How a problem is executed. 'function' is the LeetCode shape: the learner
-- writes one function and a generated driver feeds it arguments. 'stdio' is a
-- complete program reading stdin, which the Go basics problems are written as.
-- 'sql' is handed to the SQL runner verbatim.
ALTER TABLE problems ADD COLUMN IF NOT EXISTS io_mode TEXT NOT NULL DEFAULT 'function';

ALTER TABLE problems DROP CONSTRAINT IF EXISTS problems_io_mode_check;
ALTER TABLE problems ADD CONSTRAINT problems_io_mode_check
    CHECK (io_mode IN ('function', 'stdio', 'sql'));

-- Reference solutions exist to be run in CI, never to be served to a learner.
-- Every problem must have one per offered language, it must pass all the
-- problem's test cases, and the generated starter for that problem must NOT
-- pass. That pair of assertions is what catches both a broken driver and an
-- answer that has leaked into a starter template.
CREATE TABLE IF NOT EXISTS reference_solutions (
    problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    language   TEXT NOT NULL,
    code       TEXT NOT NULL,
    PRIMARY KEY (problem_id, language)
);
