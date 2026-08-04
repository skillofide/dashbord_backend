package judge

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

// resultSet is the shape the SQL runner emits and that SQL test cases expect:
//
//	{"rows": [ {...}, {...} ]}
//
// "ordered" only ever appears in the expected output. Most SQL problems say
// "return the result in any order", so rows are compared as a multiset by
// default; problems that specify an ORDER BY set it to true.
type resultSet struct {
	Rows    []map[string]any `json:"rows"`
	Ordered bool             `json:"ordered"`
}

// parseResultSet reports whether s is a SQL result set, and if so returns it.
// Anything that is not a JSON object carrying a "rows" array is not one, which
// is what lets the judge auto-detect SQL output without being told the language.
func parseResultSet(s string) (*resultSet, bool) {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "{") {
		return nil, false
	}

	var probe map[string]json.RawMessage
	if err := json.Unmarshal([]byte(s), &probe); err != nil {
		return nil, false
	}
	if _, present := probe["rows"]; !present {
		return nil, false
	}

	var rs resultSet
	if err := json.Unmarshal([]byte(s), &rs); err != nil {
		return nil, false
	}
	return &rs, true
}

// resultSetsEqual compares two SQL result sets semantically rather than as text.
//
// Three normalisations are applied, each for a concrete reason:
//
//   - Column names are compared case-insensitively. PostgreSQL folds unquoted
//     identifiers to lower case, so a query aliasing `AS Email` returns the key
//     "email". Grading on exact case would fail every correct answer.
//   - Numbers are compared numerically, so 4.00 and 4.0 match.
//   - Row order is ignored unless the expected output sets "ordered": true.
func resultSetsEqual(expected, actual *resultSet) bool {
	if len(expected.Rows) != len(actual.Rows) {
		return false
	}

	expectedKeys := canonicalRows(expected.Rows)
	actualKeys := canonicalRows(actual.Rows)

	if expected.Ordered {
		for i := range expectedKeys {
			if expectedKeys[i] != actualKeys[i] {
				return false
			}
		}
		return true
	}

	// Unordered: compare as multisets so duplicate rows still have to match in
	// count, which a plain set comparison would miss.
	sort.Strings(expectedKeys)
	sort.Strings(actualKeys)
	for i := range expectedKeys {
		if expectedKeys[i] != actualKeys[i] {
			return false
		}
	}
	return true
}

// canonicalRows renders each row as a stable string so rows can be compared and
// sorted without caring about key order within the row.
func canonicalRows(rows []map[string]any) []string {
	out := make([]string, 0, len(rows))
	for _, row := range rows {
		keys := make([]string, 0, len(row))
		for k := range row {
			keys = append(keys, strings.ToLower(k))
		}
		sort.Strings(keys)

		// Re-index the row by lower-cased key so lookup matches the sorted keys.
		lowered := make(map[string]any, len(row))
		for k, v := range row {
			lowered[strings.ToLower(k)] = v
		}

		var b strings.Builder
		for i, k := range keys {
			if i > 0 {
				b.WriteByte('\x1f') // unit separator: cannot occur in SQL output
			}
			b.WriteString(k)
			b.WriteByte('=')
			b.WriteString(canonicalValue(lowered[k]))
		}
		out = append(out, b.String())
	}
	return out
}

// canonicalValue renders a single cell so that values which are equal in SQL
// terms produce an identical string.
func canonicalValue(v any) string {
	switch t := v.(type) {
	case nil:
		return "<null>"

	case bool:
		return strconv.FormatBool(t)

	case float64:
		// encoding/json decodes every JSON number as float64. Render integral
		// values without a decimal point so 4 and 4.0 agree, and trim trailing
		// zeros elsewhere so 4.00 and 4.0 agree too.
		if math.Trunc(t) == t && math.Abs(t) < 1e15 {
			return strconv.FormatInt(int64(t), 10)
		}
		return strconv.FormatFloat(t, 'g', -1, 64)

	case string:
		// A numeric column can arrive as a JSON string (PostgreSQL renders
		// NUMERIC that way). Compare it as a number when it parses as one.
		if f, err := strconv.ParseFloat(strings.TrimSpace(t), 64); err == nil {
			return canonicalValue(f)
		}
		return strings.TrimSpace(t)

	default:
		// Nested arrays/objects: fall back to canonical JSON.
		b, err := json.Marshal(t)
		if err != nil {
			return fmt.Sprintf("%v", t)
		}
		return string(b)
	}
}
