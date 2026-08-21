-- Authors the 22 remaining placeholder problems across Foundational Basics and
-- Masters of Algorithms.
--
-- Same story as the Easy ten: a generic statement, a
-- `solveChallenge(input) { return input; }` starter, and test cases asserting
-- "1" -> "1". The untouched starter passed, so each was solvable by pressing
-- Submit on arrival.
--
-- Several titles describe a language feature rather than a computation —
-- "Function Overloading", "Pass by Value vs Reference", "Singleton Object".
-- Those are reframed so the concept is still what is being exercised but the
-- answer is a value the judge can compare across all five languages. A problem
-- that can only be expressed in one language does not belong in a set the
-- editor offers five for.
--
-- Reference solutions cover JavaScript and Python; validate-problems reports
-- the other three as unverified rather than passing them unchecked.

BEGIN;

-- ─── Statements ──────────────────────────────────────────────────────────────

UPDATE problems SET statement = v.statement
FROM (VALUES
 ('op3', 'Given an integer `n` and a shift amount `s`, return `[n << s, n >> s]`.

Left shift multiplies by two per position, right shift divides. Use the shift operators rather than arithmetic.'),
 ('op5', 'Precedence decides what binds first without any parentheses at all.

Given `a`, `b` and `c`, return the results of three expressions in this order:

1. `a + b * c`   — multiplication binds tighter than addition
2. `(a + b) * c` — parentheses override that
3. `a * b + c`'),
 ('op6', 'Compound assignment updates a variable in place.

Starting from `start`, apply `delta` four times in sequence and return the value after each step:

1. `value += delta`
2. `value -= delta`
3. `value *= delta`
4. `value %= delta` — if `delta` is `0`, use `0` for this step rather than dividing'),
 ('cond3', 'Classify an applicant using nested conditions.

Return `"Approved"` when `age` is at least 18 **and** `income` is at least 30000. When the applicant is old enough but earns less, return `"Review"`. When they are under 18, return `"Rejected"` regardless of income.'),
 ('cond4', 'Return whether `year` is a leap year.

A year is a leap year when it is divisible by 4, except that years divisible by 100 are not, unless they are also divisible by 400. So 2000 is a leap year and 1900 is not.'),
 ('cond6', 'Compute income tax from progressive brackets. Each rate applies only to the part of the income inside that band.

| Band | Rate |
|------|------|
| first 10000 | 0% |
| 10001–30000 | 10% |
| 30001–100000 | 20% |
| above 100000 | 30% |

Return the total tax. The result is compared with a tolerance.'),
 ('loop2', 'Return `n` factorial, computed with a `while` loop.

`0!` is `1`. Values up to `20!` fit in a 64-bit integer.'),
 ('loop3', 'Return the first `n` numbers of the Fibonacci series, starting `0, 1`.

For `n = 0` return an empty array.'),
 ('loop4', 'Return whether `n` is prime.

A prime has exactly two distinct divisors, 1 and itself, so `0`, `1` and negative numbers are not prime. Checking divisors up to the square root is enough.'),
 ('loop5', 'Build a centred pyramid of stars `n` rows tall and return it as an array of strings, one per row.

Row `i` (counting from 1) holds `2i - 1` stars, padded on the left with spaces so the pyramid is centred against the widest row. Do not pad the right.

For `n = 3`:

```
  *
 ***
*****
```'),
 ('func3', 'Return `n` factorial using recursion rather than a loop.

The base case is `0! = 1`; every other value calls the function again with `n - 1`.'),
 ('func4', 'Overloading picks a different behaviour for a different number of arguments. Languages that lack it reach the same end with one function that inspects what it received.

Given `values`, return:

- the single value itself, when there is one
- the sum, when there are two
- the product, when there are three or more
- `0`, when there are none'),
 ('func5', 'A scalar passed to a function is a copy; an array is a reference to the same data. Changing the copy leaves the caller''s value alone, changing the array does not.

Take `scalar` and `arr`. Inside your function add 100 to the scalar and append 100 to the array. Return `[scalarSeenByCaller, ...arrayAfterTheCall]` — that is, the scalar the *caller* still holds, followed by every element of the array after your change.'),
 ('func6', 'A closure keeps hold of the variable it was created with.

Build a counter that starts at `start`, then call it `times` times. Each call adds one and returns the new value. Return every value returned, in order.

For `start = 10` and `times = 3`, the answer is `[11, 12, 13]`.'),
 ('str3', 'Given a sentence, return it with the words in reverse order.

Words are separated by single spaces and there is no leading or trailing whitespace.'),
 ('str4', 'Return whether `a` and `b` are anagrams — the same letters in a different order.

Compare case-insensitively and ignore spaces. Everything else counts.'),
 ('str5', 'Return the index of the first occurrence of `needle` in `haystack`, or `-1` when it does not occur.

An empty needle occurs at index `0`. Knuth–Morris–Pratt reaches the answer in linear time by never re-examining a character of the haystack.'),
 ('str6', 'Return the length of the longest substring of `s` that contains no repeated character.

For `"abcabcbb"` the answer is `3`, from `"abc"`. A sliding window that moves its left edge past the previous occurrence of a repeat does this in one pass.'),
 ('obj3', 'Serialise a record to JSON.

Given a `name` and an `age`, return the JSON text for an object with those two fields, `name` first:

```
{"name":"Ada","age":36}
```

No spaces — the output is compared exactly.'),
 ('obj4', 'A subclass inherits behaviour from its parent and may override part of it.

Given an animal `kind` and a `name`, return two lines:

1. `"<name> is a <kind>"` — from the shared parent behaviour
2. the sound: `"Woof"` for `dog`, `"Meow"` for `cat`, and `"..."` for anything else'),
 ('obj5', 'A shallow copy of a nested structure shares its inner rows with the original; a deep copy does not.

Take a `matrix`, make a **deep** copy of it, then set every value in the *original* to `0`. Return your copy — if the copy was deep it still holds the values it started with.'),
 ('obj6', 'A singleton hands out the same instance every time it is asked.

Given a number of `calls`, ask for the instance that many times and return, for each call after the first, whether it was the same instance as the first one. For `calls = 3` the answer is `[true, true]`.

Return an empty array when `calls` is less than 2.')
) AS v(slug, statement) WHERE problems.slug = v.slug;

-- ─── Test cases ──────────────────────────────────────────────────────────────

DELETE FROM test_cases WHERE problem_id IN (SELECT id FROM problems WHERE slug IN
 ('op3','op5','op6','cond3','cond4','cond6','loop2','loop3','loop4','loop5',
  'func3','func4','func5','func6','str3','str4','str5','str6','obj3','obj4','obj5','obj6'));

INSERT INTO test_cases (problem_id, input, expected_output, is_hidden, order_index)
SELECT p.id, v.input, v.expected, v.hidden, v.idx FROM problems p JOIN (VALUES
 ('op3', E'8\n2', '[32,2]', false, 0), ('op3', E'1\n4', '[16,0]', false, 1), ('op3', E'255\n1', '[510,127]', true, 2),

 ('op5', E'2\n3\n4', '[14,20,10]', false, 0), ('op5', E'1\n1\n1', '[2,2,2]', false, 1), ('op5', E'5\n0\n3', '[5,15,3]', true, 2),

 ('op6', E'10\n3', '[13,10,30,0]', false, 0), ('op6', E'7\n5', '[12,7,35,0]', false, 1), ('op6', E'4\n0', '[4,4,0,0]', true, 2),

 ('cond3', E'25\n50000', '"Approved"', false, 0), ('cond3', E'25\n10000', '"Review"', false, 1),
 ('cond3', E'15\n90000', '"Rejected"', true, 2), ('cond3', E'18\n30000', '"Approved"', true, 3),

 ('cond4', '2000', 'true', false, 0), ('cond4', '1900', 'false', false, 1),
 ('cond4', '2024', 'true', true, 2), ('cond4', '2023', 'false', true, 3),

 ('cond6', '5000', '0', false, 0), ('cond6', '50000', '6000', false, 1), ('cond6', '150000', '31000', true, 2),

 ('loop2', '5', '120', false, 0), ('loop2', '0', '1', false, 1), ('loop2', '20', '2432902008176640000', true, 2),

 ('loop3', '7', '[0,1,1,2,3,5,8]', false, 0), ('loop3', '1', '[0]', false, 1), ('loop3', '0', '[]', true, 2),

 ('loop4', '17', 'true', false, 0), ('loop4', '1', 'false', false, 1),
 ('loop4', '9', 'false', true, 2), ('loop4', '2', 'true', true, 3),

 ('loop5', '3', '["  *"," ***","*****"]', false, 0), ('loop5', '1', '["*"]', false, 1),

 ('func3', '5', '120', false, 0), ('func3', '0', '1', false, 1), ('func3', '10', '3628800', true, 2),

 ('func4', '[7]', '7', false, 0), ('func4', '[3,4]', '7', false, 1),
 ('func4', '[2,3,4]', '24', true, 2), ('func4', '[]', '0', true, 3),

 ('func5', E'5\n[1,2]', '[5,1,2,100]', false, 0), ('func5', E'0\n[]', '[0,100]', false, 1),

 ('func6', E'10\n3', '[11,12,13]', false, 0), ('func6', E'0\n1', '[1]', false, 1), ('func6', E'5\n0', '[]', true, 2),

 ('str3', '"the quick brown fox"', '"fox brown quick the"', false, 0),
 ('str3', '"hello"', '"hello"', false, 1),

 ('str4', E'"listen"\n"silent"', 'true', false, 0), ('str4', E'"hello"\n"world"', 'false', false, 1),
 ('str4', E'"Dormitory"\n"dirty room"', 'true', true, 2),

 ('str5', E'"hello"\n"ll"', '2', false, 0), ('str5', E'"aaaaa"\n"bba"', '-1', false, 1),
 ('str5', E'"abc"\n""', '0', true, 2),

 ('str6', '"abcabcbb"', '3', false, 0), ('str6', '"bbbbb"', '1', false, 1),
 ('str6', '"pwwkew"', '3', true, 2), ('str6', '""', '0', true, 3),

 ('obj3', E'"Ada"\n36', '"{\\"name\\":\\"Ada\\",\\"age\\":36}"', false, 0),
 ('obj3', E'"Linus"\n54', '"{\\"name\\":\\"Linus\\",\\"age\\":54}"', false, 1),

 ('obj4', E'"dog"\n"Rex"', '["Rex is a dog","Woof"]', false, 0),
 ('obj4', E'"cat"\n"Tom"', '["Tom is a cat","Meow"]', false, 1),
 ('obj4', E'"fox"\n"Ylvis"', '["Ylvis is a fox","..."]', true, 2),

 ('obj5', '[[1,2],[3,4]]', '[[1,2],[3,4]]', false, 0), ('obj5', '[[5]]', '[[5]]', false, 1),

 ('obj6', '3', '[true,true]', false, 0), ('obj6', '1', '[]', false, 1), ('obj6', '0', '[]', true, 2)
) AS v(slug, input, expected, hidden, idx) ON v.slug = p.slug;

-- ─── Examples ────────────────────────────────────────────────────────────────
-- Input and output are filled in from the visible test cases by GetProblem.

DELETE FROM examples WHERE problem_id IN (SELECT id FROM problems WHERE slug IN
 ('op3','op5','op6','cond3','cond4','cond6','loop2','loop3','loop4','loop5',
  'func3','func4','func5','func6','str3','str4','str5','str6','obj3','obj4','obj5','obj6'));

INSERT INTO examples (problem_id, input, output, explanation, order_index)
SELECT p.id, '', '', v.explanation, 0 FROM problems p JOIN (VALUES
 ('op3',   '8 << 2 is 32; 8 >> 2 is 2.'),
 ('op5',   '2 + 3 * 4 = 14 because * binds first; (2 + 3) * 4 = 20; 2 * 3 + 4 = 10.'),
 ('op6',   '10 += 3 gives 13, then -= 3 gives 10, then *= 3 gives 30, then %= 3 gives 0.'),
 ('cond3', 'Old enough and earning at least 30000.'),
 ('cond4', '2000 is divisible by 400, so the century rule does not exclude it.'),
 ('cond6', '5000 sits entirely inside the 0% band.'),
 ('loop2', '5 x 4 x 3 x 2 x 1 = 120.'),
 ('loop3', 'Each number after the first two is the sum of the two before it.'),
 ('loop4', '17 has no divisor between 2 and its square root.'),
 ('loop5', 'Three rows of 1, 3 and 5 stars, right-aligned to width 5.'),
 ('func3', 'factorial(5) calls factorial(4), and so on down to factorial(0) = 1.'),
 ('func4', 'One value is returned unchanged.'),
 ('func5', 'The caller still sees 5 — the scalar was copied — but the array gained a 100.'),
 ('func6', 'The counter keeps its own value between calls: 11, then 12, then 13.'),
 ('str3',  'The four words are emitted last to first.'),
 ('str4',  '"listen" and "silent" use the same six letters.'),
 ('str5',  '"ll" begins at index 2 of "hello".'),
 ('str6',  '"abc" is the longest run before a character repeats.'),
 ('obj3',  'Fields appear in the order name then age, with no spaces.'),
 ('obj4',  'The first line comes from the parent, the sound from the dog override.'),
 ('obj5',  'The copy is unaffected by zeroing the original, which is what makes it deep.'),
 ('obj6',  'Calls 2 and 3 both return the same instance as call 1.')
) AS v(slug, explanation) ON v.slug = p.slug;

-- ─── Signatures ──────────────────────────────────────────────────────────────

INSERT INTO problem_signatures (problem_id, entry_point, params, return_type, compare, float_eps)
SELECT p.id, v.entry, v.params::jsonb, v.ret, v.cmp, v.eps FROM problems p JOIN (VALUES
 ('op3','shiftOps','[{"name":"n","type":"int"},{"name":"s","type":"int"}]','int[]','exact',1e-6),
 ('op5','precedence','[{"name":"a","type":"int"},{"name":"b","type":"int"},{"name":"c","type":"int"}]','int[]','exact',1e-6),
 ('op6','compound','[{"name":"start","type":"int"},{"name":"delta","type":"int"}]','int[]','exact',1e-6),
 ('cond3','classifyApplicant','[{"name":"age","type":"int"},{"name":"income","type":"int"}]','string','exact',1e-6),
 ('cond4','isLeapYear','[{"name":"year","type":"int"}]','bool','exact',1e-6),
 ('cond6','calculateTax','[{"name":"income","type":"int"}]','double','float',1e-4),
 ('loop2','factorial','[{"name":"n","type":"int"}]','long','exact',1e-6),
 ('loop3','fibonacci','[{"name":"n","type":"int"}]','int[]','exact',1e-6),
 ('loop4','isPrime','[{"name":"n","type":"int"}]','bool','exact',1e-6),
 ('loop5','pyramid','[{"name":"n","type":"int"}]','string[]','exact',1e-6),
 ('func3','factorialRecursive','[{"name":"n","type":"int"}]','long','exact',1e-6),
 ('func4','applyByArity','[{"name":"values","type":"int[]"}]','int','exact',1e-6),
 ('func5','valueVsReference','[{"name":"scalar","type":"int"},{"name":"arr","type":"int[]"}]','int[]','exact',1e-6),
 ('func6','runCounter','[{"name":"start","type":"int"},{"name":"times","type":"int"}]','int[]','exact',1e-6),
 ('str3','reverseSentence','[{"name":"sentence","type":"string"}]','string','exact',1e-6),
 ('str4','isAnagram','[{"name":"a","type":"string"},{"name":"b","type":"string"}]','bool','exact',1e-6),
 ('str5','findSubstring','[{"name":"haystack","type":"string"},{"name":"needle","type":"string"}]','int','exact',1e-6),
 ('str6','lengthOfLongestSubstring','[{"name":"s","type":"string"}]','int','exact',1e-6),
 ('obj3','serializeRecord','[{"name":"name","type":"string"},{"name":"age","type":"int"}]','string','exact',1e-6),
 ('obj4','describeAnimal','[{"name":"kind","type":"string"},{"name":"name","type":"string"}]','string[]','exact',1e-6),
 ('obj5','deepCopyMatrix','[{"name":"matrix","type":"int[][]"}]','int[][]','exact',1e-6),
 ('obj6','singletonChecks','[{"name":"calls","type":"int"}]','bool[]','exact',1e-6)
) AS v(slug, entry, params, ret, cmp, eps) ON v.slug = p.slug
ON CONFLICT (problem_id) DO UPDATE SET entry_point=EXCLUDED.entry_point, params=EXCLUDED.params,
 return_type=EXCLUDED.return_type, compare=EXCLUDED.compare, float_eps=EXCLUDED.float_eps, updated_at=now();

COMMIT;
