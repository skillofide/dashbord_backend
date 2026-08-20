package resolvers

import "testing"

// TestBandFor pins down how a score turns into money off.
//
// The ordering case is the one that matters: award_slabs is operator-entered
// JSON from the admin form, and a ladder typed low-to-high would otherwise hand
// the *smallest* discount to the strongest candidate — a bug that pays out
// wrongly and looks like a pricing decision rather than a defect.
func TestBandFor(t *testing.T) {
	ladder := []awardSlab{
		{MinPercent: 80, AwardPercent: 100},
		{MinPercent: 65, AwardPercent: 50},
		{MinPercent: 50, AwardPercent: 25},
	}

	tests := []struct {
		name      string
		slabs     []awardSlab
		percent   float64
		wantAward int
		wantOK    bool
	}{
		{"top band", ladder, 92, 100, true},
		{"exactly on a boundary qualifies", ladder, 80, 100, true},
		{"just under the top band drops one", ladder, 79.9, 50, true},
		{"middle band", ladder, 70, 50, true},
		{"lowest band", ladder, 50, 25, true},
		{"below the ladder earns nothing", ladder, 49.9, 0, false},
		{"zero earns nothing", ladder, 0, 0, false},
		{"no ladder configured", nil, 95, 0, false},

		// The ordering guard.
		{
			"ascending ladder still awards the best band",
			[]awardSlab{
				{MinPercent: 50, AwardPercent: 25},
				{MinPercent: 65, AwardPercent: 50},
				{MinPercent: 80, AwardPercent: 100},
			},
			92, 100, true,
		},
		{
			"shuffled ladder still awards the best band",
			[]awardSlab{
				{MinPercent: 65, AwardPercent: 50},
				{MinPercent: 80, AwardPercent: 100},
				{MinPercent: 50, AwardPercent: 25},
			},
			85, 100, true,
		},
		{"single band, met", []awardSlab{{MinPercent: 60, AwardPercent: 40}}, 60, 40, true},
		{"single band, missed", []awardSlab{{MinPercent: 60, AwardPercent: 40}}, 59, 0, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// A copy per case: bandFor sorts in place, and a shared fixture
			// would let one case quietly fix the input for the next.
			slabs := append([]awardSlab(nil), tc.slabs...)
			got, ok := bandFor(slabs, tc.percent)
			if ok != tc.wantOK {
				t.Fatalf("qualified = %v, want %v", ok, tc.wantOK)
			}
			if ok && got.AwardPercent != tc.wantAward {
				t.Errorf("award = %d%%, want %d%% for a score of %.1f%%",
					got.AwardPercent, tc.wantAward, tc.percent)
			}
		})
	}
}
