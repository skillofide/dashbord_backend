package judge

import "testing"

func TestMatchesTyped(t *testing.T) {
	exact := CompareSpec{Mode: CompareExact}

	tests := []struct {
		name     string
		expected string
		actual   string
		spec     CompareSpec
		want     bool
	}{
		// The formatting differences that used to produce false Wrong Answers.
		{"python list spacing", "[12,8,20,5]", "[12, 8, 20, 5]", exact, true},
		{"quoted string result", `"fox brown quick the"`, `"fox brown quick the"`, exact, true},
		{"bool", "true", "true", exact, true},
		{"trailing newline", "[0,1]", "[0,1]\n", exact, true},

		// Genuine failures must stay failures.
		{"wrong values", "[0,1]", "[1,2]", exact, false},
		{"wrong length", "[0,1]", "[0,1,2]", exact, false},
		{"empty vs answer", "[0,1]", "[]", exact, false},
		{"bool flipped", "true", "false", exact, false},
		{"string differs", `"abc"`, `"abd"`, exact, false},
		{"null vs value", "null", "0", exact, false},

		// "Return the answer in any order."
		{"unordered pair", "[0,1]", "[1,0]", CompareSpec{Mode: CompareUnordered}, true},
		{"unordered still checks content", "[0,1]", "[1,2]", CompareSpec{Mode: CompareUnordered}, false},
		{"exact rejects reorder", "[0,1]", "[1,0]", exact, false},
		{"unordered nested", "[[1,2],[3,4]]", "[[4,3],[2,1]]", CompareSpec{Mode: CompareUnordered}, true},

		// Sets ignore duplicates as well.
		{"set dedupes", "[1,2,3]", "[3,2,1,1]", CompareSpec{Mode: CompareSet}, true},
		{"unordered does not dedupe", "[1,2,3]", "[3,2,1,1]", CompareSpec{Mode: CompareUnordered}, false},

		// Floats.
		{"float tolerance", "0.3", "0.30000000000000004", CompareSpec{Mode: CompareExact, Eps: 1e-6}, true},
		{"float outside tolerance", "0.3", "0.4", CompareSpec{Mode: CompareExact, Eps: 1e-6}, false},

		// Non-JSON expectations fall back to text comparison rather than
		// silently passing or silently failing.
		{"non-json both sides", "Hello World", "Hello World", exact, true},
		{"non-json differs", "Hello World", "Goodbye", exact, false},
		{"stack trace is not a pass", "[0,1]", "Traceback (most recent call last)", exact, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := MatchesTyped(tc.expected, tc.actual, tc.spec); got != tc.want {
				t.Errorf("MatchesTyped(%q, %q, %v) = %v, want %v",
					tc.expected, tc.actual, tc.spec.Mode, got, tc.want)
			}
		})
	}
}
