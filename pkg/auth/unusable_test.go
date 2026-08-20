package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func bcryptCompare(stored, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(stored), []byte(plain)) == nil
}

// TestUnusablePasswordCannotAuthenticate guards a sharp edge.
//
// CheckPassword falls back to a constant-time *plaintext* comparison for
// anything that is not a bcrypt hash. That makes any placeholder written into
// the password column a live credential by default — so an account provisioned
// without a password would be logged into by whoever could guess the
// placeholder. These tests exist to make sure that can never be true again.
func TestUnusablePasswordCannotAuthenticate(t *testing.T) {
	stored, err := UnusablePassword()
	if err != nil {
		t.Fatalf("UnusablePassword: %v", err)
	}
	if !IsUnusable(stored) {
		t.Fatalf("IsUnusable(%q) = false", stored)
	}

	// The obvious attack: type the stored value itself.
	if ok, _ := CheckPassword(stored, stored); ok {
		t.Error("the stored placeholder authenticated as its own password")
	}
	for _, guess := range []string{"", unusablePrefix, "password", stored[len(unusablePrefix):]} {
		if ok, _ := CheckPassword(stored, guess); ok {
			t.Errorf("guess %q authenticated against an unusable password", guess)
		}
	}
}

func TestUnusablePasswordsAreDistinct(t *testing.T) {
	a, _ := UnusablePassword()
	b, _ := UnusablePassword()
	if a == b {
		t.Error("two provisioned accounts got the same placeholder — the column would " +
			"show which accounts were created together")
	}
}

// A real password must still work, and must not be mistaken for a placeholder.
// The placeholder must also be rejected by the *old* verification path — the
// one that predates IsUnusable — because services redeploy at different times.
// Starting the value with a bcrypt prefix is what buys that: the old code takes
// the bcrypt branch, fails to parse it, and says no.
func TestUnusablePasswordFailsClosedOnOldVerification(t *testing.T) {
	stored, err := UnusablePassword()
	if err != nil {
		t.Fatalf("UnusablePassword: %v", err)
	}
	if !IsHashed(stored) {
		t.Fatal("placeholder is not bcrypt-prefixed — a build without IsUnusable would " +
			"compare it as PLAINTEXT and let anyone who read the column sign in")
	}
	// Exactly what the pre-IsUnusable CheckPassword did.
	legacyCheck := func(stored, plain string) bool {
		if IsHashed(stored) {
			return bcryptCompare(stored, plain)
		}
		return stored == plain // the plaintext fallback
	}
	for _, guess := range []string{stored, "password", ""} {
		if legacyCheck(stored, guess) {
			t.Errorf("old verification accepted %q against an unusable password", guess)
		}
	}
}

func TestRealPasswordsUnaffected(t *testing.T) {
	hashed, err := HashPassword("correct horse battery staple")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if IsUnusable(hashed) {
		t.Fatal("a real bcrypt hash was treated as unusable")
	}
	if ok, _ := CheckPassword(hashed, "correct horse battery staple"); !ok {
		t.Error("a real password stopped working")
	}
	if ok, _ := CheckPassword(hashed, "wrong"); ok {
		t.Error("a wrong password was accepted")
	}
}
