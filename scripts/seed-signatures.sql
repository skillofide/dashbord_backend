-- Declared signatures for the function-mode problems, plus the test-case
-- normalisation those signatures imply.
--
-- Test-case inputs are one JSON value per line, in parameter order. Expected
-- output is a single JSON value. Most of the data was already in that shape —
-- it was the driver that could not read it — but the three string-returning
-- problems stored their input and expectation as bare text, which cannot be
-- distinguished from a malformed JSON document.

BEGIN;

-- ─── Normalise test-case data to the JSON contract ───────────────────────────

-- reverse-words takes one string and returns one string. Both sides were stored
-- unquoted, so a correct answer ("fox brown quick the") was compared against a
-- JSON-encoded one ("\"fox brown quick the\"") and lost on the quotes alone.
UPDATE test_cases tc
SET    input           = to_json(tc.input)::text,
       expected_output = to_json(tc.expected_output)::text
FROM   problems p
WHERE  p.id = tc.problem_id
  AND  p.slug = 'reverse-words';

-- html-image-alt and css-flexbox-center take no parameters at all. Their input
-- was the literal string "None", which is not a value — it is a placeholder
-- standing in for the absence of one.
UPDATE test_cases tc
SET    input           = '',
       expected_output = to_json(tc.expected_output)::text
FROM   problems p
WHERE  p.id = tc.problem_id
  AND  p.slug IN ('html-image-alt', 'css-flexbox-center');

-- ─── Signatures ──────────────────────────────────────────────────────────────

INSERT INTO problem_signatures (problem_id, entry_point, params, return_type, compare)
SELECT p.id, v.entry_point, v.params::jsonb, v.return_type, v.compare
FROM   problems p
JOIN (VALUES
    -- LeetCode "you can return the answer in any order", so [1,0] is as correct
    -- as [0,1]. The old judge compared text and rejected the reversal.
    ('leetcode-two-sum', 'twoSum',
     '[{"name":"nums","type":"int[]"},{"name":"target","type":"int"}]', 'int[]', 'unordered'),

    ('best-time-to-buy-and-sell-stock', 'maxProfit',
     '[{"name":"prices","type":"int[]"}]', 'int', 'exact'),

    ('contains-duplicate', 'containsDuplicate',
     '[{"name":"nums","type":"int[]"}]', 'bool', 'exact'),

    ('op1', 'arithmeticOperations',
     '[{"name":"a","type":"int"},{"name":"b","type":"int"}]', 'int[]', 'exact'),

    ('cond1', 'checkEvenOdd',
     '[{"name":"n","type":"int"}]', 'string', 'exact'),

    ('loop1', 'sumOfN',
     '[{"name":"n","type":"int"}]', 'int', 'exact'),

    ('str2', 'isPalindrome',
     '[{"name":"s","type":"string"}]', 'bool', 'exact'),

    ('reverse-words', 'reverseWords',
     '[{"name":"sentence","type":"string"}]', 'string', 'exact'),

    ('html-image-alt', 'getImageTag', '[]', 'string', 'exact'),

    ('css-flexbox-center', 'getFlexCenterCSS', '[]', 'string', 'exact')
) AS v(slug, entry_point, params, return_type, compare)
  ON v.slug = p.slug
ON CONFLICT (problem_id) DO UPDATE
SET entry_point = EXCLUDED.entry_point,
    params      = EXCLUDED.params,
    return_type = EXCLUDED.return_type,
    compare     = EXCLUDED.compare,
    updated_at  = now();

COMMIT;
