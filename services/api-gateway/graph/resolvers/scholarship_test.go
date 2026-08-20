package resolvers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
)

// newTestScholarshipHandler builds a handler with a nil pool.
//
// That is the whole trick: any code path that reaches Postgres panics, so these
// tests prove a rejection happened *before* the database was touched. On an
// unauthenticated endpoint that ordering is the difference between a cheap 400
// and free work for anyone who can send a request.
func newTestScholarshipHandler() *ScholarshipHandler {
	return &ScholarshipHandler{
		Log:          zap.NewNop(),
		jwtSecret:    "test-secret",
		appBase:      "http://portal.test",
		ipLimiter:    newRateLimiter(1000, time.Minute),
		emailLimiter: newRateLimiter(1000, time.Minute),
		mailer:       &scholarshipMailer{log: zap.NewNop()}, // host empty → disabled
	}
}

func postJSON(t *testing.T, h http.Handler, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestApplyRejectsBeforeTouchingTheDatabase(t *testing.T) {
	h := newTestScholarshipHandler()

	tests := []struct {
		name string
		body string
		want int
		why  string
	}{
		{
			"honeypot is answered like a success",
			`{"name":"Bot","email":"bot@example.com","courseId":"5","website":"http://spam"}`,
			http.StatusOK,
			"a bot that fills the hidden field must learn nothing from the response",
		},
		{
			"missing name",
			`{"email":"a@b.com","courseId":"5"}`,
			http.StatusBadRequest,
			"",
		},
		{
			"malformed email",
			`{"name":"A","email":"not-an-email","courseId":"5"}`,
			http.StatusBadRequest,
			"",
		},
		{
			"missing course",
			`{"name":"A","email":"a@b.com"}`,
			http.StatusBadRequest,
			"the course decides which paper they sit; there is no default",
		},
		{
			"malformed body",
			`{"name":`,
			http.StatusBadRequest,
			"",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// A panic here means the handler reached the nil pool, i.e. it did
			// database work for a request it should have rejected outright.
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("handler hit the database before rejecting: %v", r)
				}
			}()
			rec := postJSON(t, h, "/api/scholarship/apply", tc.body)
			if rec.Code != tc.want {
				t.Errorf("status = %d, want %d (%s)", rec.Code, tc.want, tc.why)
			}
		})
	}
}

func TestHoneypotResponseIsIndistinguishableFromSuccess(t *testing.T) {
	h := newTestScholarshipHandler()
	rec := postJSON(t, h, "/api/scholarship/apply",
		`{"name":"Bot","email":"bot@example.com","courseId":"5","website":"x"}`)

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("honeypot response is not JSON: %v", err)
	}
	if body["success"] != true {
		t.Errorf("honeypot response = %v, want success:true so the bot cannot adapt", body)
	}
	// It must not leak a real credential either.
	if _, ok := body["testUrl"]; ok {
		t.Error("honeypot response handed out a test link")
	}
}

func TestClaimRejectsEmptyToken(t *testing.T) {
	h := newTestScholarshipHandler()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("claim hit the database for an empty token: %v", r)
		}
	}()
	if rec := postJSON(t, h, "/api/scholarship/claim", `{"token":""}`); rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestScholarshipRouting(t *testing.T) {
	h := newTestScholarshipHandler()

	tests := []struct {
		method, path string
		want         int
	}{
		{http.MethodGet, "/api/scholarship/apply", http.StatusMethodNotAllowed},
		{http.MethodGet, "/api/scholarship/claim", http.StatusMethodNotAllowed},
		{http.MethodPost, "/api/scholarship/config", http.StatusMethodNotAllowed},
		{http.MethodPost, "/api/scholarship/nope", http.StatusNotFound},
		{http.MethodOptions, "/api/scholarship/apply", http.StatusNoContent},
		// Trailing slashes must not open a second, unrouted spelling of a path
		// that the auth middleware allow-lists by exact match.
		{http.MethodGet, "/api/scholarship/apply/", http.StatusMethodNotAllowed},
	}

	for _, tc := range tests {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != tc.want {
				t.Errorf("status = %d, want %d", rec.Code, tc.want)
			}
		})
	}
}

func TestClaimTokensAreRandomAndStoredHashed(t *testing.T) {
	a, err := randomToken(32)
	if err != nil {
		t.Fatalf("randomToken: %v", err)
	}
	b, err := randomToken(32)
	if err != nil {
		t.Fatalf("randomToken: %v", err)
	}
	if a == b {
		t.Fatal("two claim tokens came out identical")
	}
	if len(a) < 40 {
		t.Errorf("claim token is only %d chars — too short to resist guessing", len(a))
	}
	// URL-safe: the token travels as a query parameter.
	if strings.ContainsAny(a, "+/=") {
		t.Errorf("claim token %q is not URL-safe", a)
	}

	// The digest is what gets written down, and it must not be reversible to
	// the raw token by simple inspection.
	h := sha256Hex(a)
	if h == a || strings.Contains(h, a) {
		t.Error("stored digest reveals the raw token")
	}
	if len(h) != 64 {
		t.Errorf("sha256Hex produced %d chars, want 64", len(h))
	}
	if sha256Hex(a) != h {
		t.Error("sha256Hex is not stable — a stored digest would never match again")
	}
}

func TestParseOptionalTime(t *testing.T) {
	if parseOptionalTime("") != nil {
		t.Error("blank should clear the window, not fail the save")
	}
	if parseOptionalTime("not a date") != nil {
		t.Error("unparseable should clear the window, not fail the save")
	}
	got := parseOptionalTime("2026-09-01T00:00:00Z")
	if got == nil || !got.Equal(time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("parseOptionalTime = %v, want 2026-09-01T00:00:00Z", got)
	}
}
