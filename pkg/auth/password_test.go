package auth

import (
	"strings"
	"testing"
)

func TestHashPasswordProducesVerifiableHash(t *testing.T) {
	hashed, err := HashPassword("knovate123")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if !IsHashed(hashed) {
		t.Fatalf("HashPassword produced something IsHashed rejects: %q", hashed)
	}
	if strings.Contains(hashed, "knovate123") {
		t.Fatal("hash contains the plaintext password")
	}

	ok, legacy := CheckPassword(hashed, "knovate123")
	if !ok {
		t.Fatal("correct password rejected against its own hash")
	}
	if legacy {
		t.Fatal("a freshly created hash was reported as legacy plaintext")
	}
}

func TestHashPasswordIsSalted(t *testing.T) {
	a, err := HashPassword("same-password")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	b, err := HashPassword("same-password")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if a == b {
		t.Fatal("identical passwords produced identical hashes — salt is missing")
	}
}

func TestCheckPasswordRejectsWrongPassword(t *testing.T) {
	hashed, err := HashPassword("correct-horse")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if ok, _ := CheckPassword(hashed, "battery-staple"); ok {
		t.Fatal("wrong password accepted")
	}
}

// The whole point of the legacy branch: accounts created before hashing must
// still log in, and must be flagged so the caller can upgrade the row.
func TestCheckPasswordAcceptsLegacyPlaintextAndFlagsIt(t *testing.T) {
	ok, legacy := CheckPassword("knovate123", "knovate123")
	if !ok {
		t.Fatal("legacy plaintext password rejected — existing users would be locked out")
	}
	if !legacy {
		t.Fatal("legacy plaintext row not flagged for re-hashing")
	}
}

func TestCheckPasswordRejectsWrongLegacyPassword(t *testing.T) {
	if ok, _ := CheckPassword("knovate123", "wrong"); ok {
		t.Fatal("wrong password accepted against a legacy plaintext row")
	}
}

// A stored plaintext value that happens to start with a bcrypt prefix must not
// be treated as a valid hash it cannot possibly be.
func TestCheckPasswordRejectsMalformedHash(t *testing.T) {
	if ok, _ := CheckPassword("$2a$not-a-real-hash", "anything"); ok {
		t.Fatal("malformed bcrypt hash accepted")
	}
}

// bcrypt silently ignores input past 72 bytes, which would make two different
// long passwords interchangeable. HashPassword must refuse instead.
func TestHashPasswordRejectsOverlongInput(t *testing.T) {
	if _, err := HashPassword(strings.Repeat("a", 73)); err == nil {
		t.Fatal("password longer than bcrypt's 72-byte limit was accepted")
	}
}
