package grading

import "testing"

func TestScoreMcqSingle(t *testing.T) {
	base := McqInput{Kind: "single", Marks: 4, CorrectIDs: []string{"a"}, NegativeMarking: 0.25}

	tests := []struct {
		name     string
		selected []string
		want     float64
	}{
		{"correct", []string{"a"}, 4},
		{"wrong", []string{"b"}, -1},
		{"unanswered scores zero, never negative", nil, 0},
		{"multi-select on a single-answer question is wrong", []string{"a", "b"}, -1},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			in := base
			in.SelectedIDs = tc.selected
			if got := ScoreMcq(in); got != tc.want {
				t.Errorf("ScoreMcq() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestScoreMcqMultiple(t *testing.T) {
	base := McqInput{Kind: "multiple", Marks: 6, CorrectIDs: []string{"a", "b", "c"}}

	tests := []struct {
		name string
		in   McqInput
		want float64
	}{
		{
			"all-or-nothing awards nothing for a partial set",
			McqInput{Kind: "multiple", Marks: 6, CorrectIDs: base.CorrectIDs, SelectedIDs: []string{"a", "b"}},
			0,
		},
		{
			"all-or-nothing awards full marks for the exact set",
			McqInput{Kind: "multiple", Marks: 6, CorrectIDs: base.CorrectIDs, SelectedIDs: []string{"c", "a", "b"}},
			6,
		},
		{
			"partial credit awards the found fraction",
			McqInput{Kind: "multiple", Marks: 6, CorrectIDs: base.CorrectIDs, SelectedIDs: []string{"a", "b"}, PartialCredit: true},
			4,
		},
		{
			"a single wrong option voids partial credit",
			McqInput{Kind: "multiple", Marks: 6, CorrectIDs: base.CorrectIDs, SelectedIDs: []string{"a", "b", "d"}, PartialCredit: true, NegativeMarking: 0.5},
			-3,
		},
		{
			"duplicate selections do not inflate partial credit",
			McqInput{Kind: "multiple", Marks: 6, CorrectIDs: base.CorrectIDs, SelectedIDs: []string{"a", "a", "a"}, PartialCredit: true},
			2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := ScoreMcq(tc.in); got != tc.want {
				t.Errorf("ScoreMcq() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestScoreNumeric(t *testing.T) {
	in := McqInput{Kind: "numeric", Marks: 5, AcceptedAnswers: []string{"2.5"}, NegativeMarking: 0.2}

	in.TextAnswer = "2.50"
	if got := ScoreMcq(in); got != 5 {
		t.Errorf("numerically equal answer = %v, want 5", got)
	}

	in.TextAnswer = "3"
	if got := ScoreMcq(in); got != -1 {
		t.Errorf("wrong numeric answer = %v, want -1", got)
	}

	in.TextAnswer = "  "
	if got := ScoreMcq(in); got != 0 {
		t.Errorf("blank numeric answer = %v, want 0", got)
	}
}

func TestScoreCoding(t *testing.T) {
	tests := []struct {
		name          string
		passed, total int
		marks         float64
		want          float64
	}{
		{"all hidden cases pass", 10, 10, 20, 20},
		{"half pass", 5, 10, 20, 10},
		{"none pass", 0, 10, 20, 0},
		{"no test cases cannot award marks", 0, 0, 20, 0},
		{"passed above total is clamped to full marks", 12, 10, 20, 20},
		{"fractional marks are rounded to two places", 1, 3, 10, 3.33},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := ScoreCoding(tc.passed, tc.total, tc.marks); got != tc.want {
				t.Errorf("ScoreCoding(%d, %d, %v) = %v, want %v", tc.passed, tc.total, tc.marks, got, tc.want)
			}
		})
	}
}

func TestIntegrityPenaltyIsNeverFree(t *testing.T) {
	if got := IntegrityPenalty("some_future_signal"); got <= 0 {
		t.Errorf("unknown proctor event penalty = %v, want > 0", got)
	}
	if IntegrityPenalty("devtools") <= IntegrityPenalty("copy") {
		t.Error("a devtools event should cost more than a copy event")
	}
}
