-- Authors the ten Foundational Basics problems that shipped as placeholders.
--
-- All ten had the same shape: a generic statement ("In this challenge, you are
-- asked to solve the algorithmic problem: **Title**"), a starter of
-- `solveChallenge(input) { return input; }`, and two test cases asserting
-- "1" -> "1" and "2" -> "2". The untouched starter passed both, so every one of
-- them could be cleared, and its XP collected, without typing anything.
--
-- Each is rewritten to test what its title claims. Signatures use only types the
-- code generator can emit in all five languages, which is why the two "object"
-- problems are expressed as values in and values out rather than as a struct —
-- an object literal has no common shape across JavaScript, Go and C++.
--
-- Reference solutions here cover JavaScript and Python. Java, C++ and Go still
-- need theirs; validate-problems reports those as "no-reference" rather than
-- letting them pass unverified.

BEGIN;

-- ─── Statements ──────────────────────────────────────────────────────────────

UPDATE problems SET statement = v.statement
FROM (VALUES
  ('op2', 'Given two integers `a` and `b`, return the results of four comparisons as an array of booleans, in this order:

1. `a > b`
2. `a < b`
3. `a == b`
4. `a != b`

Use the relational operators directly — no `if` statements are needed.'),

  ('op4', 'Given an integer `n`, return `"Positive"` when it is greater than zero, `"Negative"` when it is less than zero, and `"Zero"` otherwise.

Write it as a single expression using the ternary operator `? :` rather than an `if` chain.'),

  ('cond2', 'Given an integer `score` between 0 and 100, return the letter grade:

| Score | Grade |
|-------|-------|
| 90–100 | `"A"` |
| 80–89 | `"B"` |
| 70–79 | `"C"` |
| 60–69 | `"D"` |
| below 60 | `"F"` |

A `switch` on `Math.floor(score / 10)` handles every band without listing all 101 scores.'),

  ('cond5', 'Given three integers `a`, `b` and `c`, return the largest of them.

If two or more are equally large, return that value.'),

  ('loop6', 'A menu-driven program reads a choice, acts on it, and repeats until the user chooses to exit.

Given an array of menu `choices`, return the label produced by each one, in order:

| Choice | Label |
|--------|-------|
| 1 | `"Add"` |
| 2 | `"Subtract"` |
| 3 | `"Multiply"` |
| 0 | `"Exit"` |
| anything else | `"Invalid"` |

Stop as soon as you handle a `0` — `"Exit"` is included in the result, and any choices after it are ignored.'),

  ('func1', 'Write a function that takes a `name` and returns a greeting.

For the name `Ada`, return exactly `Hello, Ada!`.'),

  ('func2', 'Given the `radius` of a circle, return its area.

Use π = `3.14159`. The result is compared with a tolerance, so small floating-point differences are fine.'),

  ('str1', 'Given a string `s`, return an array of two integers: the number of vowels followed by the number of consonants.

Count only English letters — digits, spaces and punctuation are neither. Treat upper and lower case alike; `a`, `e`, `i`, `o` and `u` are the vowels.'),

  ('obj1', 'Build a record for a student from their `name`, `age` and `grade`, and return it as a formatted string:

```
Name: Ada | Age: 20 | Grade: A
```

Mind the spacing around the `|` separators — the output is compared exactly.'),

  ('obj2', 'A counter object exposes a getter and a setter over a single value.

Starting from `initial`, apply each number in `updates` with the setter, reading the value back with the getter after each one. Return every value you read, in order.

Setting a value that is already stored still counts as a read, so the result always has exactly as many entries as `updates`.')
) AS v(slug, statement)
WHERE problems.slug = v.slug;

-- ─── Test cases ──────────────────────────────────────────────────────────────

DELETE FROM test_cases WHERE problem_id IN (
  SELECT id FROM problems WHERE slug IN
    ('op2','op4','cond2','cond5','loop6','func1','func2','str1','obj1','obj2'));

INSERT INTO test_cases (problem_id, input, expected_output, is_hidden, order_index)
SELECT p.id, v.input, v.expected, v.hidden, v.idx
FROM   problems p
JOIN (VALUES
  ('op2',   E'7\n3',        '[true,false,false,true]',  false, 0),
  ('op2',   E'3\n7',        '[false,true,false,true]',  false, 1),
  ('op2',   E'5\n5',        '[false,false,true,false]', true,  2),

  ('op4',   '42',           '"Positive"',               false, 0),
  ('op4',   '-7',           '"Negative"',               false, 1),
  ('op4',   '0',            '"Zero"',                   true,  2),

  ('cond2', '95',           '"A"',                      false, 0),
  ('cond2', '72',           '"C"',                      false, 1),
  ('cond2', '100',          '"A"',                      true,  2),
  ('cond2', '59',           '"F"',                      true,  3),
  ('cond2', '60',           '"D"',                      true,  4),

  ('cond5', E'3\n9\n5',     '9',                        false, 0),
  ('cond5', E'-1\n-8\n-4',  '-1',                       false, 1),
  ('cond5', E'4\n4\n2',     '4',                        true,  2),

  ('loop6', '[1,2,0]',      '["Add","Subtract","Exit"]', false, 0),
  ('loop6', '[3,9,1]',      '["Multiply","Invalid","Add"]', false, 1),
  ('loop6', '[0,1,2]',      '["Exit"]',                 true,  2),

  ('func1', '"Ada"',        '"Hello, Ada!"',            false, 0),
  ('func1', '"Grace"',      '"Hello, Grace!"',          false, 1),
  ('func1', '""',           '"Hello, !"',               true,  2),

  ('func2', '2',            '12.56636',                 false, 0),
  ('func2', '1',            '3.14159',                  false, 1),
  ('func2', '0',            '0',                        true,  2),

  ('str1',  '"Hello World"', '[3,7]',                   false, 0),
  ('str1',  '"xyz"',        '[0,3]',                    false, 1),
  ('str1',  '"a1 e2 i3!"',  '[3,0]',                    true,  2),

  ('obj1',  E'"Ada"\n20\n"A"',    '"Name: Ada | Age: 20 | Grade: A"',   false, 0),
  ('obj1',  E'"Linus"\n33\n"B"',  '"Name: Linus | Age: 33 | Grade: B"', false, 1),

  ('obj2',  E'0\n[5,5,9]',  '[5,5,9]',                  false, 0),
  ('obj2',  E'3\n[1,2]',    '[1,2]',                    false, 1),
  ('obj2',  E'7\n[]',       '[]',                       true,  2)
) AS v(slug, input, expected, hidden, idx) ON v.slug = p.slug;

-- ─── Examples ────────────────────────────────────────────────────────────────
-- Only the explanations are stored; GetProblem fills input and output in from
-- the visible test cases so the two can never disagree.

DELETE FROM examples WHERE problem_id IN (
  SELECT id FROM problems WHERE slug IN
    ('op2','op4','cond2','cond5','loop6','func1','func2','str1','obj1','obj2'));

INSERT INTO examples (problem_id, input, output, explanation, order_index)
SELECT p.id, '', '', v.explanation, v.idx
FROM   problems p
JOIN (VALUES
  ('op2',   '7 > 3 is true, 7 < 3 is false, 7 == 3 is false, 7 != 3 is true.', 0),
  ('op4',   '42 is greater than zero.', 0),
  ('cond2', '95 falls in the 90–100 band.', 0),
  ('cond5', '9 is larger than both 3 and 5.', 0),
  ('loop6', 'Choices 1 and 2 map to Add and Subtract; 0 ends the run after adding Exit.', 0),
  ('func1', 'The name is placed between "Hello, " and "!".', 0),
  ('func2', '3.14159 x 2 x 2 = 12.56636.', 0),
  ('str1',  '"Hello World" has 3 vowels (e, o, o) and 7 consonants; the space counts as neither.', 0),
  ('obj1',  'Each field is labelled and separated by " | ".', 0),
  ('obj2',  'The value is set to 5, then to 5 again, then to 9 — read back after each.', 0)
) AS v(slug, explanation, idx) ON v.slug = p.slug;

-- ─── Signatures ──────────────────────────────────────────────────────────────

INSERT INTO problem_signatures (problem_id, entry_point, params, return_type, compare, float_eps)
SELECT p.id, v.entry_point, v.params::jsonb, v.return_type, v.compare, v.eps
FROM   problems p
JOIN (VALUES
  ('op2',   'relationalChecks', '[{"name":"a","type":"int"},{"name":"b","type":"int"}]', 'bool[]',   'exact', 1e-6),
  ('op4',   'classify',         '[{"name":"n","type":"int"}]',                          'string',   'exact', 1e-6),
  ('cond2', 'grade',            '[{"name":"score","type":"int"}]',                      'string',   'exact', 1e-6),
  ('cond5', 'maxOfThree',       '[{"name":"a","type":"int"},{"name":"b","type":"int"},{"name":"c","type":"int"}]', 'int', 'exact', 1e-6),
  ('loop6', 'runMenu',          '[{"name":"choices","type":"int[]"}]',                  'string[]', 'exact', 1e-6),
  ('func1', 'greet',            '[{"name":"name","type":"string"}]',                    'string',   'exact', 1e-6),
  -- Area is floating point; comparing it exactly would fail on the last bit.
  ('func2', 'circleArea',       '[{"name":"radius","type":"double"}]',                  'double',   'float', 1e-6),
  ('str1',  'countVowelsAndConsonants', '[{"name":"s","type":"string"}]',               'int[]',    'exact', 1e-6),
  ('obj1',  'studentSummary',   '[{"name":"name","type":"string"},{"name":"age","type":"int"},{"name":"grade","type":"string"}]', 'string', 'exact', 1e-6),
  ('obj2',  'applyUpdates',     '[{"name":"initial","type":"int"},{"name":"updates","type":"int[]"}]', 'int[]', 'exact', 1e-6)
) AS v(slug, entry_point, params, return_type, compare, eps) ON v.slug = p.slug
ON CONFLICT (problem_id) DO UPDATE
SET entry_point = EXCLUDED.entry_point,
    params      = EXCLUDED.params,
    return_type = EXCLUDED.return_type,
    compare     = EXCLUDED.compare,
    float_eps   = EXCLUDED.float_eps,
    updated_at  = now();

-- ─── Reference solutions (JavaScript, Python) ────────────────────────────────
-- Never served to a learner. They exist so validate-problems can assert the
-- problem is solvable and that the generated starter is not already a solution.

DELETE FROM reference_solutions WHERE problem_id IN (
  SELECT id FROM problems WHERE slug IN
    ('op2','op4','cond2','cond5','loop6','func1','func2','str1','obj1','obj2'));

INSERT INTO reference_solutions (problem_id, language, code)
SELECT p.id, v.lang, v.code
FROM   problems p
JOIN (VALUES
  ('op2', 'javascript', $c$function relationalChecks(a, b) {
    return [a > b, a < b, a === b, a !== b];
}$c$),
  ('op2', 'python', $c$def relationalChecks(a: int, b: int) -> list[bool]:
    return [a > b, a < b, a == b, a != b]$c$),

  ('op4', 'javascript', $c$function classify(n) {
    return n > 0 ? "Positive" : n < 0 ? "Negative" : "Zero";
}$c$),
  ('op4', 'python', $c$def classify(n: int) -> str:
    return "Positive" if n > 0 else "Negative" if n < 0 else "Zero"$c$),

  ('cond2', 'javascript', $c$function grade(score) {
    switch (Math.floor(score / 10)) {
        case 10:
        case 9: return "A";
        case 8: return "B";
        case 7: return "C";
        case 6: return "D";
        default: return "F";
    }
}$c$),
  ('cond2', 'python', $c$def grade(score: int) -> str:
    band = score // 10
    if band >= 9:
        return "A"
    if band == 8:
        return "B"
    if band == 7:
        return "C"
    if band == 6:
        return "D"
    return "F"$c$),

  ('cond5', 'javascript', $c$function maxOfThree(a, b, c) {
    return Math.max(a, Math.max(b, c));
}$c$),
  ('cond5', 'python', $c$def maxOfThree(a: int, b: int, c: int) -> int:
    return max(a, b, c)$c$),

  ('loop6', 'javascript', $c$function runMenu(choices) {
    const labels = { 1: "Add", 2: "Subtract", 3: "Multiply", 0: "Exit" };
    const out = [];
    for (const choice of choices) {
        out.push(labels[choice] !== undefined ? labels[choice] : "Invalid");
        if (choice === 0) break;
    }
    return out;
}$c$),
  ('loop6', 'python', $c$def runMenu(choices: list[int]) -> list[str]:
    labels = {1: "Add", 2: "Subtract", 3: "Multiply", 0: "Exit"}
    out = []
    for choice in choices:
        out.append(labels.get(choice, "Invalid"))
        if choice == 0:
            break
    return out$c$),

  ('func1', 'javascript', $c$function greet(name) {
    return "Hello, " + name + "!";
}$c$),
  ('func1', 'python', $c$def greet(name: str) -> str:
    return "Hello, " + name + "!"$c$),

  ('func2', 'javascript', $c$function circleArea(radius) {
    return 3.14159 * radius * radius;
}$c$),
  ('func2', 'python', $c$def circleArea(radius: float) -> float:
    return 3.14159 * radius * radius$c$),

  ('str1', 'javascript', $c$function countVowelsAndConsonants(s) {
    let vowels = 0, consonants = 0;
    for (const ch of s.toLowerCase()) {
        if (ch < "a" || ch > "z") continue;
        if ("aeiou".includes(ch)) vowels++;
        else consonants++;
    }
    return [vowels, consonants];
}$c$),
  ('str1', 'python', $c$def countVowelsAndConsonants(s: str) -> list[int]:
    vowels = consonants = 0
    for ch in s.lower():
        if not ("a" <= ch <= "z"):
            continue
        if ch in "aeiou":
            vowels += 1
        else:
            consonants += 1
    return [vowels, consonants]$c$),

  ('obj1', 'javascript', $c$function studentSummary(name, age, grade) {
    const student = { name: name, age: age, grade: grade };
    return "Name: " + student.name + " | Age: " + student.age + " | Grade: " + student.grade;
}$c$),
  ('obj1', 'python', $c$def studentSummary(name: str, age: int, grade: str) -> str:
    student = {"name": name, "age": age, "grade": grade}
    return f"Name: {student['name']} | Age: {student['age']} | Grade: {student['grade']}"$c$),

  ('obj2', 'javascript', $c$function applyUpdates(initial, updates) {
    let value = initial;
    const seen = [];
    for (const next of updates) {
        value = next;
        seen.push(value);
    }
    return seen;
}$c$),
  ('obj2', 'python', $c$def applyUpdates(initial: int, updates: list[int]) -> list[int]:
    value = initial
    seen = []
    for nxt in updates:
        value = nxt
        seen.append(value)
    return seen$c$)
) AS v(slug, lang, code) ON v.slug = p.slug;

COMMIT;
