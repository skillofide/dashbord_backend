#!/usr/bin/env python3
"""
Test harness for SQL submissions.

Reads a test-case fixture from USER_INPUT, materialises it as real tables in a
throwaway PostgreSQL database, runs the learner's SQL from USER_CODE, and prints
the result as canonical JSON on stdout:

    {"rows": [ {...}, {...} ]}

The fixture format is:

    {
      "tables": { "Person": [ {"id": 1, "email": "a@b.com"}, ... ], ... },
      "verify": "SELECT id, email FROM Person ORDER BY id"   // optional
    }

"verify" exists for INSERT/UPDATE/DELETE problems, where the thing being graded
is the final state of a table rather than a result set. When present, the
learner's statement is executed for its effect and "verify" produces the output.

Only the Python standard library is used; queries are executed by shelling out
to psql, so the image needs no database driver.
"""

import json
import os
import re
import subprocess
import sys

PSQL = ["psql", "-v", "ON_ERROR_STOP=1", "-q", "-t", "-A", "-d", "judge"]

DATE_RE = re.compile(r"^\d{4}-\d{2}-\d{2}$")
TIMESTAMP_RE = re.compile(r"^\d{4}-\d{2}-\d{2}[ T]\d{2}:\d{2}")


def fail(message, code=1):
    """Report a harness//user error on stderr and exit non-zero."""
    print(message, file=sys.stderr)
    sys.exit(code)


def run_psql(sql, capture=True):
    """Execute SQL through psql. Returns stdout; raises on SQL error."""
    proc = subprocess.run(
        PSQL,
        input=sql,
        text=True,
        capture_output=True,
    )
    if proc.returncode != 0:
        # psql puts SQL errors on stderr; surface them verbatim so the learner
        # sees the real database message.
        fail(proc.stderr.strip() or "SQL execution failed")
    return proc.stdout if capture else ""


SIMPLE_IDENT_RE = re.compile(r"^[A-Za-z_][A-Za-z0-9_]*$")


def quote_ident(name):
    """
    Emit an identifier the way a learner would write it.

    Simple identifiers are left UNQUOTED so PostgreSQL folds them to lower
    case, exactly as it folds the unquoted names in the learner's query. A
    fixture table "Person" therefore becomes `person`, and `FROM Person` in the
    submission resolves to it. Quoting here would create a case-sensitive
    "Person" that `FROM Person` could never match.

    Because of that folding, result keys come back lower case, so the judge
    compares column names case-insensitively.
    """
    name = str(name)
    if SIMPLE_IDENT_RE.match(name):
        return name
    return '"' + name.replace('"', '""') + '"'


def quote_literal(value):
    if value is None:
        return "NULL"
    if isinstance(value, bool):
        return "TRUE" if value else "FALSE"
    if isinstance(value, (int, float)):
        return repr(value)
    return "'" + str(value).replace("'", "''") + "'"


def infer_type(values):
    """
    Infer a column type from every non-null value in the column.

    Checked most specific first. Booleans must be tested before integers
    because bool is a subclass of int in Python.
    """
    present = [v for v in values if v is not None]
    if not present:
        return "TEXT"

    if all(isinstance(v, bool) for v in present):
        return "BOOLEAN"
    if all(isinstance(v, int) and not isinstance(v, bool) for v in present):
        return "BIGINT"
    if all(isinstance(v, (int, float)) and not isinstance(v, bool) for v in present):
        return "NUMERIC"
    if all(isinstance(v, str) and DATE_RE.match(v) for v in present):
        return "DATE"
    if all(isinstance(v, str) and TIMESTAMP_RE.match(v) for v in present):
        return "TIMESTAMP"
    return "TEXT"


def build_schema(tables):
    """Emit CREATE TABLE + INSERT statements for the whole fixture."""
    statements = []

    for table_name, rows in tables.items():
        if not isinstance(rows, list):
            fail(f"fixture table {table_name!r} must be a list of row objects")

        # Column order follows first appearance so the schema is deterministic.
        columns = []
        for row in rows:
            for col in row:
                if col not in columns:
                    columns.append(col)

        if not columns:
            statements.append(f"CREATE TABLE {quote_ident(table_name)} ();")
            continue

        col_defs = []
        for col in columns:
            col_type = infer_type([row.get(col) for row in rows])
            col_defs.append(f"  {quote_ident(col)} {col_type}")

        statements.append(
            f"CREATE TABLE {quote_ident(table_name)} (\n" + ",\n".join(col_defs) + "\n);"
        )

        if rows:
            col_list = ", ".join(quote_ident(c) for c in columns)
            values = ",\n".join(
                "  (" + ", ".join(quote_literal(row.get(c)) for c in columns) + ")"
                for row in rows
            )
            statements.append(
                f"INSERT INTO {quote_ident(table_name)} ({col_list}) VALUES\n{values};"
            )

    return "\n".join(statements)


def strip_sql_comments(sql):
    """Remove -- line comments and /* */ block comments."""
    sql = re.sub(r"/\*.*?\*/", " ", sql, flags=re.S)
    sql = re.sub(r"--[^\n]*", " ", sql)
    return sql.strip()


def is_query(sql):
    """True when the statement produces a result set rather than mutating data."""
    head = strip_sql_comments(sql).lstrip("( \t\n").upper()
    return head.startswith("SELECT") or head.startswith("WITH") or head.startswith("TABLE")


def as_result_set(sql):
    """
    Wrap a query so PostgreSQL returns the whole result as one JSON array.

    coalesce handles the empty result, which json_agg reports as NULL. The
    trailing semicolon is stripped because the query becomes a subquery.

    Row order from the inner query is preserved into the array, so problems
    that specify an ORDER BY can be graded on order (see the "ordered" flag
    in the judge).
    """
    inner = strip_sql_comments(sql).rstrip().rstrip(";")
    return (
        "SELECT coalesce(json_agg(row_to_json(__t)), '[]'::json)::text\n"
        f"FROM (\n{inner}\n) AS __t;"
    )


def main():
    raw_input_json = os.environ.get("USER_INPUT", "").strip()
    user_code = os.environ.get("USER_CODE", "").strip()

    if not user_code:
        fail("No SQL submitted.")

    if not raw_input_json:
        fixture = {}
    else:
        try:
            fixture = json.loads(raw_input_json)
        except json.JSONDecodeError as err:
            fail(f"Malformed test fixture: {err}")

    tables = fixture.get("tables") or {}
    verify = fixture.get("verify")

    # Single-table shorthand: {"table": "Person", "rows": [...]}
    if not tables and fixture.get("table"):
        tables = {fixture["table"]: fixture.get("rows") or []}

    # 1. Build the fixture tables.
    if tables:
        run_psql(build_schema(tables), capture=False)

    # 2. Run the learner's SQL, then decide where the graded output comes from.
    if verify:
        # Modification problem: run the statement for effect, grade final state.
        run_psql(strip_sql_comments(user_code), capture=False)
        output_sql = as_result_set(verify)
    elif is_query(user_code):
        output_sql = as_result_set(user_code)
    else:
        fail(
            "Your statement does not return rows. This problem expects a SELECT query."
        )

    raw = run_psql(output_sql).strip()

    try:
        rows = json.loads(raw) if raw else []
    except json.JSONDecodeError:
        fail(f"Could not parse query result: {raw[:500]}")

    print(json.dumps({"rows": rows}, separators=(",", ":"), default=str))


if __name__ == "__main__":
    main()
