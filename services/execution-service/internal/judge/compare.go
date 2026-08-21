package judge

import (
	"encoding/json"
	"math"
	"sort"
	"strings"
)

// CompareMode mirrors problem_signatures.compare.
type CompareMode string

const (
	CompareExact     CompareMode = "exact"
	CompareUnordered CompareMode = "unordered"
	CompareSet       CompareMode = "set"
	CompareFloat     CompareMode = "float"
)

// CompareSpec is everything the judge needs to grade one test case beyond the
// two strings themselves.
type CompareSpec struct {
	Mode CompareMode
	Eps  float64
}

// MatchesTyped decides whether actual satisfies expected for a function-mode
// problem.
//
// Both sides are decoded as JSON and compared structurally. The previous
// implementation compared whitespace-normalised text, which meant the verdict
// depended on how a language happened to format its output: Python's
// "[12, 8, 20, 5]" lost to "[12,8,20,5]", Python's "True" lost to "true", and a
// string result printed bare lost to a quoted one. None of those are wrong
// answers.
//
// If either side fails to decode — a stack trace, a partial write, a
// deliberately non-JSON expectation — this falls back to the old text
// comparison rather than declaring a match it cannot justify.
func MatchesTyped(expected, actual string, spec CompareSpec) bool {
	var exp, act interface{}
	if err := json.Unmarshal([]byte(strings.TrimSpace(expected)), &exp); err != nil {
		return normalize(actual) == normalize(expected)
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(actual)), &act); err != nil {
		return normalize(actual) == normalize(expected)
	}

	eps := spec.Eps
	if eps <= 0 {
		eps = 1e-6
	}

	switch spec.Mode {
	case CompareUnordered:
		return equalValues(canonicalSort(exp, false), canonicalSort(act, false), eps)
	case CompareSet:
		return equalValues(canonicalSort(exp, true), canonicalSort(act, true), eps)
	default:
		return equalValues(exp, act, eps)
	}
}

// equalValues is deep equality over decoded JSON, with numbers compared within
// eps so that 0.30000000000000004 and 0.3 agree.
func equalValues(a, b interface{}, eps float64) bool {
	switch av := a.(type) {
	case nil:
		return b == nil
	case bool:
		bv, ok := b.(bool)
		return ok && av == bv
	case float64:
		bv, ok := b.(float64)
		if !ok {
			return false
		}
		if math.IsNaN(av) || math.IsNaN(bv) {
			return math.IsNaN(av) && math.IsNaN(bv)
		}
		return math.Abs(av-bv) <= eps
	case string:
		bv, ok := b.(string)
		return ok && av == bv
	case []interface{}:
		bv, ok := b.([]interface{})
		if !ok || len(av) != len(bv) {
			return false
		}
		for i := range av {
			if !equalValues(av[i], bv[i], eps) {
				return false
			}
		}
		return true
	case map[string]interface{}:
		bv, ok := b.(map[string]interface{})
		if !ok || len(av) != len(bv) {
			return false
		}
		for k, v := range av {
			other, present := bv[k]
			if !present || !equalValues(v, other, eps) {
				return false
			}
		}
		return true
	}
	return false
}

// canonicalSort orders a top-level array so that "return them in any order"
// problems can be graded on content. Nested arrays are sorted too, because a
// result like [[1,2],[3,4]] is usually order-insensitive at both levels when it
// is order-insensitive at all. dedupe additionally collapses repeats for set
// comparison. Non-arrays pass through untouched.
func canonicalSort(v interface{}, dedupe bool) interface{} {
	arr, ok := v.([]interface{})
	if !ok {
		return v
	}
	out := make([]interface{}, 0, len(arr))
	for _, item := range arr {
		out = append(out, canonicalSort(item, dedupe))
	}
	sort.SliceStable(out, func(i, j int) bool {
		return sortKey(out[i]) < sortKey(out[j])
	})
	if !dedupe {
		return out
	}
	deduped := make([]interface{}, 0, len(out))
	var prev string
	for i, item := range out {
		k := sortKey(item)
		if i == 0 || k != prev {
			deduped = append(deduped, item)
		}
		prev = k
	}
	return deduped
}

// sortKey gives every decoded value a stable, total ordering. It is only used
// to impose a deterministic order before comparison, never shown to anyone, so
// the exact ordering matters far less than it being the same on both sides.
func sortKey(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}
