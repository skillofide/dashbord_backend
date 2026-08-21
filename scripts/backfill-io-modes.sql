-- Backfill problems.io_mode from the structure of the problem itself.
--
-- The three shapes were always present in the data; nothing recorded which was
-- which, so the judge tried to treat all of them as function-mode and generated
-- a driver for programs that already had a main().

BEGIN;

-- 1. SQL problems are identified by their tag and handed to the SQL runner
--    verbatim. Every one of them carries the 'SQL' tag.
UPDATE problems p SET io_mode = 'sql'
WHERE EXISTS (
    SELECT 1 FROM problem_tags t WHERE t.problem_id = p.id AND t.tag = 'SQL'
);

-- 2. A Go starter that declares func main() is a complete stdio program. These
--    read stdin (or print a fixed answer) and must run untouched — generating a
--    main() for them produced "main redeclared" and they could not compile in
--    any language.
UPDATE problems p SET io_mode = 'stdio'
WHERE p.io_mode <> 'sql'
  AND EXISTS (
    SELECT 1 FROM starter_codes s
    WHERE s.problem_id = p.id
      AND s.go ~ '(?n)^[[:space:]]*func[[:space:]]+main[[:space:]]*\([[:space:]]*\)'
  );

-- 3. Everything else keeps the default, 'function'.

COMMIT;
