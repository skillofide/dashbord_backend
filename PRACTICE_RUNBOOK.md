# Practice problems — deploying the execution contract

The judge used to guess how to run a submission. It scanned the source for the
first function it could find and tried to decode the whole test input as an
argument list, which fails the moment an input spans more than one line — Two
Sum's own reference solution returned `[]` on Two Sum.

Problems now **declare** a contract: entry point, typed parameters, return type,
comparison mode. Both the starter template a learner sees and the driver the
judge wraps around their submission are generated from that one declaration, so
they cannot disagree, and a starter cannot contain the answer.

Deploying this is not just new binaries. The contract lives in the database, and
a service running the new code against un-migrated, un-seeded data is worse than
one running the old code — it will report every problem as having no signature
and quietly fall back to the behaviour this replaced.

Do these in order.

---

## 1. Migrations

```bash
migrate -path=services/problem-service/migrations \
        -database "$POSTGRES_DSN&x-migrations-table=schema_migrations_problem" up
```

Adds `problem_signatures`, `reference_solutions`, `problems.io_mode`, and the
class-mode columns.

`getExecutionSpec` selects `p.io_mode` and joins `problem_signatures` on every
`GetTestCases` call — which is every Run and every Submit. Start the new
problem-service before this runs and every one of them errors on a missing
column.

**Check:**

```sql
SELECT to_regclass('problem_signatures'), to_regclass('reference_solutions');
-- both non-null
```

---

## 2. Runner images

```bash
./scripts/build-runners.sh
```

All six: `python javascript java cpp go sql`.

A missing image is not a visible failure. The judge answers
`create container: No such image: skillofide/runner-<lang>:latest`, which reaches
a learner as a RuntimeError on correct code. Locally this had silently zeroed
every SQL problem — the image had never been built, and nothing said so.

**Check:**

```bash
docker images --format '{{.Repository}}' | grep -c 'skillofide/runner-'   # 6
```

---

## 3. Classify the problems

```bash
psql "$POSTGRES_DSN" -f scripts/backfill-io-modes.sql
```

Sets `io_mode` from the structure of each problem. Three shapes were always in
the data; nothing recorded which was which:

| `io_mode` | What it is | How it runs |
|-----------|------------|-------------|
| `stdio` | a complete Go program, ~418 of them | verbatim |
| `sql` | tagged SQL | verbatim, to the SQL runner |
| `function` | one entry point | generated driver |

The stdio problems matter most here. They already declare `func main()`, and
generating another one produced `main redeclared` plus a duplicated import
block — they failed to compile in every language regardless of what the learner
wrote. Marking them `stdio` is what makes them run at all.

**Check:**

```sql
SELECT io_mode, count(*) FROM problems GROUP BY 1;
-- roughly: stdio 418, function 132, sql 17
```

---

## 4. Signatures, content, and starters

```bash
psql "$POSTGRES_DSN" -f scripts/seed-signatures.sql
psql "$POSTGRES_DSN" -f scripts/author-foundational-problems.sql
psql "$POSTGRES_DSN" -f scripts/author-remaining-placeholders.sql

go run ./services/execution-service/cmd/infer-signatures -apply
go run ./services/execution-service/cmd/infer-signatures -normalize -apply
go run ./services/execution-service/cmd/regen-starters -apply
```

`infer-signatures` derives a contract from data that already exists: parameter
names from the JavaScript starter, parameter *types* from the test cases
themselves. Types come from the data because the data is what the judge has to
handle — if every case feeds an array of integers in the first position, the
first parameter is `int[]` whatever a comment claims. It refuses rather than
guesses: a problem whose cases disagree on a type, or whose parameter is named
`root` and is therefore a serialised tree rather than the array it looks like,
is left alone.

`-normalize` rewrites test inputs to one JSON value per line. **Run it after any
signature you write by hand.** A signature declaring two parameters against an
input still stored as `l1 = [2,4,3], l2 = [5,6,4]` on a single line feeds the
whole string to the first parameter and nothing to the second. The signature and
the data have to change together, and only `-apply` on the inference path does
both automatically.

`regen-starters` replaces starters with skeletons generated from the signature
and files whatever was there into `reference_solutions`. On most problems that
previous content *was* the answer: 361 of 442 Easy problems shipped a working
solution in the editor, clearable by pressing Submit on arrival.

---

## 5. Prove it before pointing traffic at it

```bash
./scripts/validate-problems.sh -j 1
go run ./services/execution-service/cmd/driver-conformance
```

**validate-problems** asserts two things per problem, per language: the
reference solution is accepted, and the generated starter is **not**. The second
half is what catches an answer that has leaked back into a template.

**driver-conformance** round-trips every supported type through every language's
driver and asserts the value comes back byte-identical. This is what proves the
drivers — as opposed to proving whichever problems happen to have solutions
written in that language.

Serial (`-j 1`) on a shared box. Parallel runs push compiled languages past
their time limit, and a timeout is indistinguishable from a wrong answer unless
you read the flag — a correct C++ solution was misread as failing exactly this
way during development.

Both are wired into `.github/workflows/validate-problems.yml` and run on changes
to codegen, the judge, the sandbox, migrations, or seed data.

---

## 6. Services

```bash
docker compose build execution-service problem-service api-gateway
docker compose up -d execution-service problem-service api-gateway
docker exec <redis> redis-cli FLUSHDB
```

The cache flush is not optional. problem-service caches `GetProblem`, so a
correctly rebuilt service will keep serving pre-migration problems — including
the old `supportedLanguages` — until the cache turns over.

**Verify the build actually succeeded.** `docker compose up -d` will happily
restart the previous image after a failed build, and the container reports
healthy. A whole afternoon of "why is this still broken" traced to exactly that:
the Dockerfile built `./cmd/...`, which had grown to match the tools alongside
the service, and Go cannot write several binaries to one output path.

```bash
docker inspect <execution-service> --format '{{.Created}}'   # must be after the build
```

---

## Checking a real problem end to end

```bash
go run ./services/execution-service/cmd/validate-problems -slug leetcode-two-sum
```

Or in the browser: open Two Sum, leave the custom input empty, press **Run
code**. The starter returns `[]` and should report **Wrong Answer** against
`[0,1]`. That is the correct outcome — a starter that passes is a starter
containing the answer.

Then paste a real solution and it should pass both visible cases.

### If Run code says Accepted without grading

The custom-input box is non-empty. Run treats that as "run my program against
this instead", which is the ungraded path. It is no longer pre-filled, but a
stray character will do it.

### If a correct solution is marked wrong

Check what the judge actually fed it:

```sql
SELECT input, expected_output FROM test_cases tc
JOIN problems p ON p.id = tc.problem_id WHERE p.slug = '<slug>';
```

Input must be one JSON value per line, one line per parameter. If it is a single
line of `nums = [...], target = 9`, step 4's `-normalize` has not been run
against this problem.

---

## Known gaps

- **Class-design problems** (`lru-cache`, `min-stack`, `implement-trie`,
  `design-tic-tac-toe`) work in JavaScript and Python only. Java, C++ and Go
  refuse them rather than emit a driver that silently skips an unknown method
  call.
- **`linked-list-cycle-ii`** has no signature. Its expected output is the
  sentence `tail connects to node index 1`, and a cycle cannot be expressed in
  the serialised list form. It needs re-authoring, not a signature.
- **Reference solutions** cover JavaScript and Python. The other three languages
  are proved by `driver-conformance` rather than per problem, so
  `validate-problems` reports them as coverage gaps rather than failures.
- **418 stdio problems are Go-only.** Their other four starters are stock filler
  repeated across hundreds of unrelated problems, so those languages are not
  offered for them.
