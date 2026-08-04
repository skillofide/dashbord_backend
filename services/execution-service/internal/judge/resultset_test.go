package judge

import "testing"

func TestMatchesExpectedSQL(t *testing.T) {
	cases := []struct {
		name     string
		expected string
		actual   string
		want     bool
	}{
		{
			name:     "identical rows",
			expected: `{"rows":[{"email":"a@b.com"}]}`,
			actual:   `{"rows":[{"email":"a@b.com"}]}`,
			want:     true,
		},
		{
			name: "column case differs (PostgreSQL folds identifiers)",
			// The problem asks for "Email"; PostgreSQL returns "email".
			expected: `{"rows":[{"Email":"a@b.com"}]}`,
			actual:   `{"rows":[{"email":"a@b.com"}]}`,
			want:     true,
		},
		{
			name:     "number formatting differs",
			expected: `{"rows":[{"score":4.00,"rank":1}]}`,
			actual:   `{"rows":[{"score":4.0,"rank":1}]}`,
			want:     true,
		},
		{
			name:     "numeric returned as string",
			expected: `{"rows":[{"average_price":6.96}]}`,
			actual:   `{"rows":[{"average_price":"6.96"}]}`,
			want:     true,
		},
		{
			name:     "row order ignored by default",
			expected: `{"rows":[{"id":1},{"id":2}]}`,
			actual:   `{"rows":[{"id":2},{"id":1}]}`,
			want:     true,
		},
		{
			name:     "row order enforced when ordered is set",
			expected: `{"rows":[{"id":1},{"id":2}],"ordered":true}`,
			actual:   `{"rows":[{"id":2},{"id":1}]}`,
			want:     false,
		},
		{
			name:     "correct order when ordered is set",
			expected: `{"rows":[{"id":1},{"id":2}],"ordered":true}`,
			actual:   `{"rows":[{"id":1},{"id":2}]}`,
			want:     true,
		},
		{
			name:     "nulls match",
			expected: `{"rows":[{"city":null,"state":null}]}`,
			actual:   `{"rows":[{"city":null,"state":null}]}`,
			want:     true,
		},
		{
			name:     "null is not the empty string",
			expected: `{"rows":[{"city":null}]}`,
			actual:   `{"rows":[{"city":""}]}`,
			want:     false,
		},
		{
			name:     "missing row fails",
			expected: `{"rows":[{"id":1},{"id":2}]}`,
			actual:   `{"rows":[{"id":1}]}`,
			want:     false,
		},
		{
			name:     "duplicate rows must match in count",
			expected: `{"rows":[{"id":1},{"id":1}]}`,
			actual:   `{"rows":[{"id":1},{"id":2}]}`,
			want:     false,
		},
		{
			name:     "wrong value fails",
			expected: `{"rows":[{"email":"a@b.com"}]}`,
			actual:   `{"rows":[{"email":"c@d.com"}]}`,
			want:     false,
		},
		{
			name:     "empty result sets match",
			expected: `{"rows":[]}`,
			actual:   `{"rows":[]}`,
			want:     true,
		},
		{
			name:     "extra column fails",
			expected: `{"rows":[{"id":1}]}`,
			actual:   `{"rows":[{"id":1,"extra":2}]}`,
			want:     false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := matchesExpected(tc.expected, tc.actual); got != tc.want {
				t.Errorf("matchesExpected() = %v, want %v\n  expected: %s\n  actual:   %s",
					got, tc.want, tc.expected, tc.actual)
			}
		})
	}
}

// Non-SQL output must keep using the original text comparison.
func TestMatchesExpectedFallsBackToText(t *testing.T) {
	cases := []struct {
		name     string
		expected string
		actual   string
		want     bool
	}{
		{"plain text equal", "42", "42", true},
		{"plain text differs", "42", "43", false},
		{"trailing whitespace ignored", "42\n", "42  \n", true},
		{"json without rows key is text", `{"value":1}`, `{"value":1}`, true},
		{"json without rows key, differing", `{"value":1}`, `{"value":2}`, false},
		{"array output stays textual", "[0,1]", "[0,1]", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := matchesExpected(tc.expected, tc.actual); got != tc.want {
				t.Errorf("matchesExpected(%q, %q) = %v, want %v", tc.expected, tc.actual, got, tc.want)
			}
		})
	}
}
