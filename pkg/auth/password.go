package auth

import (
	"crypto/subtle"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Password storage.
//
// Rows created before this package existed hold the password in plain text.
// Rather than force a flag-day migration of the users table, CheckPassword
// accepts both forms and reports which one it saw, so callers can re-hash a
// legacy row the next time its owner successfully logs in. Once every active
// account has signed in at least once, the plaintext branch is dead and can be
// deleted along with the legacy return value.

// bcryptCost is the work factor. bcrypt.DefaultCost (10) is the usual
// starting point; raise it as hardware improves — existing hashes keep working
// because the cost is stored inside the hash itself.
const bcryptCost = bcrypt.DefaultCost

// HashPassword returns a bcrypt hash of plain.
//
// bcrypt ignores input beyond 72 bytes, so it rejects anything longer outright
// rather than silently truncating.
func HashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// IsHashed reports whether stored is already a bcrypt hash rather than a
// leftover plaintext password.
func IsHashed(stored string) bool {
	return strings.HasPrefix(stored, "$2a$") ||
		strings.HasPrefix(stored, "$2b$") ||
		strings.HasPrefix(stored, "$2y$")
}

// CheckPassword verifies plain against the stored value.
//
// legacy is true when the stored value was plaintext, which the caller should
// treat as a signal to re-hash and persist. It is only meaningful when ok is
// also true.
func CheckPassword(stored, plain string) (ok bool, legacy bool) {
	if IsHashed(stored) {
		return bcrypt.CompareHashAndPassword([]byte(stored), []byte(plain)) == nil, false
	}
	// Legacy plaintext row. Compared in constant time so a wrong password
	// cannot be narrowed down by how quickly it is rejected.
	return subtle.ConstantTimeCompare([]byte(stored), []byte(plain)) == 1, true
}
