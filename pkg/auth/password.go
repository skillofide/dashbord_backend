package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
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

// unusablePrefix marks a stored password that can never authenticate anybody.
//
// Two things make this value safe, and both are deliberate.
//
// It is checked explicitly below rather than relied on to "just not match",
// because CheckPassword falls back to a plaintext comparison for anything that
// is not a bcrypt hash — an arbitrary placeholder would quietly become a live
// credential for whoever could read the column.
//
// And it begins with "$2a$" so that IsHashed reports true, which means even a
// build of this package from *before* the explicit check takes the bcrypt
// branch, fails to parse the value as a hash, and rejects. That matters because
// several services verify passwords with their own compiled copy of this
// package: without the prefix, the safety of every provisioned account would
// depend on all of them being redeployed together, and the one left behind
// would treat the placeholder as a plaintext password. A security property
// should not rest on deploy order.
const unusablePrefix = "$2a$unusable$"

// UnusablePassword returns a value for an account that has no password yet.
//
// Accounts provisioned on someone's behalf — a scholarship applicant who has
// only ever filled in a form — need a password column that nothing can satisfy
// until they set one through the reset flow. Generating a random password and
// bcrypting it would do the same job, but it costs a deliberate ~100ms of CPU
// per account to protect a secret nobody is ever told: at fifty applications a
// minute that is the slowest thing in the request.
//
// The random suffix keeps rows distinct so the column cannot be used to spot
// which accounts were provisioned together.
func UnusablePassword() (string, error) {
	var b [18]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return unusablePrefix + base64.RawURLEncoding.EncodeToString(b[:]), nil
}

// IsUnusable reports whether an account has no password set.
func IsUnusable(stored string) bool { return strings.HasPrefix(stored, unusablePrefix) }

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
	// Checked first, and before the plaintext fallback below could ever see it.
	if IsUnusable(stored) {
		return false, false
	}
	if IsHashed(stored) {
		return bcrypt.CompareHashAndPassword([]byte(stored), []byte(plain)) == nil, false
	}
	// Legacy plaintext row. Compared in constant time so a wrong password
	// cannot be narrowed down by how quickly it is rejected.
	return subtle.ConstantTimeCompare([]byte(stored), []byte(plain)) == 1, true
}
