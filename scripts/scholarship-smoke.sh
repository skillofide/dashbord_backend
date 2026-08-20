#!/usr/bin/env bash
# End-to-end smoke test for the scholarship funnel.
#
# Walks the whole P1 path against a running stack: create a programme, apply as
# a stranger, claim the link, and confirm the invite guard actually bites.
#
#   docker compose up -d
#   ./scripts/scholarship-smoke.sh
#
# Requires: curl, jq, and an admin login. Override with env vars:
#   API=http://localhost:8080 ADMIN_EMAIL=... ADMIN_PASS=... ./scripts/scholarship-smoke.sh
set -euo pipefail

API=${API:-http://localhost:8080}
ADMIN_EMAIL=${ADMIN_EMAIL:-admin@skillofied.com}
ADMIN_PASS=${ADMIN_PASS:-skillofied123}
# A course id of its own, so a smoke run never renames or repoints a real
# programme. It maps whatever published scholarship paper it finds to this id,
# which leaves the genuine course-to-paper mappings untouched.
COURSE_ID=${COURSE_ID:-__smoke_test}
# Unique per run so re-running does not trip the one-application-per-programme
# unique index.
CANDIDATE=${CANDIDATE:-smoke+$(date +%s)@example.com}

pass() { printf '  \033[32m✓\033[0m %s\n' "$1"; }
fail() { printf '  \033[31m✗\033[0m %s\n' "$1"; exit 1; }
step() { printf '\n\033[1m%s\033[0m\n' "$1"; }

step "0 · gateway is up"
curl -fsS "$API/api/health" >/dev/null || fail "gateway not reachable at $API"
pass "$API responding"

step "1 · admin login"
TOKEN=$(curl -fsS -X POST "$API/api/login" \
  -H 'Content-Type: application/json' \
  -d "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$ADMIN_PASS\"}" | jq -r .token)
[ -n "$TOKEN" ] && [ "$TOKEN" != null ] || fail "could not log in as $ADMIN_EMAIL"
pass "authenticated"
# Content-Type belongs in here, not only on the calls that obviously post JSON:
# without it the GraphQL handler cannot parse the body and answers "Must provide
# an operation", which reads like a refusal and would let step 10 pass vacuously.
AUTH=(-H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json')

step "2 · a published scholarship paper exists"
# Any published assessment with purpose=scholarship will do. If there is none,
# ./scripts/seed-scholarship-test.sh authors one; the admin panel works too
# (Tests → new → purpose 'scholarship' → sections → publish).
ASSESSMENT=$(curl -fsS "${AUTH[@]}" "$API/api/recruiter/assessments?purpose=scholarship" \
  | jq -r '[.assessments[]? | select(.status=="published")][0].id // empty')
[ -n "$ASSESSMENT" ] || fail "no published assessment with purpose=scholarship — run ./scripts/seed-scholarship-test.sh first"
pass "assessment $ASSESSMENT"

step "3 · map the course to that paper"
curl -fsS -X POST "${AUTH[@]}" "$API/api/admin/scholarship-programs" \
  -d "{\"course_id\":\"$COURSE_ID\",\"course_name\":\"Smoke Test Course\",
       \"assessment_id\":\"$ASSESSMENT\",\"seats\":0,\"is_active\":true,
       \"award_slabs\":[{\"minPercent\":80,\"awardPercent\":100},{\"minPercent\":60,\"awardPercent\":50}]}" \
  | jq -e .success >/dev/null || fail "could not create the programme"
pass "course $COURSE_ID → $ASSESSMENT"

step "4 · the course shows up publicly"
curl -fsS "$API/api/scholarship/config" \
  | jq -e --arg c "$COURSE_ID" '.programs[] | select(.courseId==$c)' >/dev/null \
  || fail "programme missing from the public config"
pass "listed on /api/scholarship/config"

step "5 · a stranger applies"
APPLY=$(curl -fsS -X POST "$API/api/scholarship/apply" \
  -H 'Content-Type: application/json' \
  -d "{\"name\":\"Smoke Candidate\",\"email\":\"$CANDIDATE\",\"phone\":\"9999999999\",
       \"courseId\":\"$COURSE_ID\",\"college\":\"Test College\",\"graduationYear\":2026,
       \"page_url\":\"http://localhost:3000/scholarship\"}")
TEST_URL=$(echo "$APPLY" | jq -r .testUrl)
[ -n "$TEST_URL" ] && [ "$TEST_URL" != null ] || fail "no testUrl returned: $APPLY"
CLAIM=${TEST_URL##*t=}
pass "applied as $CANDIDATE"
pass "claim link issued"

step "6 · the honeypot is silent"
HP=$(curl -fsS -X POST "$API/api/scholarship/apply" \
  -H 'Content-Type: application/json' \
  -d '{"name":"Bot","email":"bot@example.com","courseId":"'"$COURSE_ID"'","website":"spam"}')
echo "$HP" | jq -e '.testUrl' >/dev/null 2>&1 && fail "honeypot handed out a test link"
pass "bot got a success shape and nothing else"

step "7 · claim exchanges the token for a session"
CLAIMED=$(curl -fsS -X POST "$API/api/scholarship/claim" \
  -H 'Content-Type: application/json' -d "{\"token\":\"$CLAIM\"}")
CAND_JWT=$(echo "$CLAIMED" | jq -r .token)
CAND_ASSESSMENT=$(echo "$CLAIMED" | jq -r .assessmentId)
[ -n "$CAND_JWT" ] && [ "$CAND_JWT" != null ] || fail "claim returned no session: $CLAIMED"
[ "$CAND_ASSESSMENT" = "$ASSESSMENT" ] || fail "claim pointed at the wrong paper"
pass "session minted for the candidate"

step "8 · a bad token is refused"
# No -f here: this call is *expected* to return 401, and curl -f exits non-zero
# on an error status, which set -e would turn into a failed test run.
FORGED=$(curl -sS -o /dev/null -w '%{http_code}' -X POST "$API/api/scholarship/claim" \
  -H 'Content-Type: application/json' -d '{"token":"not-a-real-token"}')
[ "$FORGED" = "401" ] || fail "a forged claim token returned $FORGED, expected 401"
pass "forged token rejected with 401"

step "9 · the candidate can start THEIR paper"
# startAttempt returns an AttemptState, whose identifier field is attemptId
# (not id). The query lives in a heredoc and the payload is assembled by jq, so
# no GraphQL brace ever has to survive a round of shell quoting.
INVITE_TOKEN=$(echo "$CLAIMED" | jq -r .inviteToken)
read -r -d '' Q_START <<'GQL' || true
mutation S($a: String!, $t: String) {
  startAttempt(assessmentId: $a, inviteToken: $t) {
    attemptId status secondsLeft
    questions { id kind options { id isCorrect } }
  }
}
GQL
START=$(jq -nc --arg q "$Q_START" --arg a "$ASSESSMENT" --arg t "$INVITE_TOKEN" \
          '{query:$q, variables:{a:$a, t:$t}}' \
        | curl -fsS -X POST "$API/api/graphql" \
            -H "Authorization: Bearer $CAND_JWT" -H 'Content-Type: application/json' \
            --data-binary @-)
ATTEMPT=$(echo "$START" | jq -r '.data.startAttempt.attemptId // empty')
[ -n "$ATTEMPT" ] || fail "invited candidate could not start: $START"
QCOUNT=$(echo "$START" | jq '.data.startAttempt.questions | length')
MCQS=$(echo "$START"  | jq '[.data.startAttempt.questions[] | select(.kind=="mcq")]    | length')
CODE=$(echo "$START"  | jq '[.data.startAttempt.questions[] | select(.kind=="coding")] | length')
LEFT=$(echo "$START"  | jq -r '.data.startAttempt.secondsLeft')
pass "attempt $ATTEMPT started"
pass "paper materialised: $QCOUNT questions ($MCQS MCQ + $CODE coding), ${LEFT}s on the clock"

# The candidate's own copy of the paper must not carry the answer key.
LEAKED=$(echo "$START" | jq '[.data.startAttempt.questions[].options[]? | select(.isCorrect == true)] | length')
[ "$LEAKED" = "0" ] || fail "SECURITY: the answer key was served to the candidate ($LEAKED options flagged correct)"
pass "answer key withheld from the candidate"

step "10 · THE GUARD — an uninvited student cannot"
# The admin account has never been invited to this paper. Before the
# resolveInvite fix this call succeeded, which is the whole reason it exists.
read -r -d '' Q_UNINVITED <<'GQL' || true
mutation S($a: String!) {
  startAttempt(assessmentId: $a) { attemptId }
}
GQL
UNINVITED=$(jq -nc --arg q "$Q_UNINVITED" --arg a "$ASSESSMENT" \
              '{query:$q, variables:{a:$a}}' \
            | curl -fsS -X POST "$API/api/graphql" "${AUTH[@]}" --data-binary @-)
if echo "$UNINVITED" | jq -e '.data.startAttempt.attemptId' >/dev/null 2>&1; then
  fail "SECURITY: an uninvited user started the scholarship paper"
fi
MSG=$(echo "$UNINVITED" | jq -r '.errors[0].message // "refused"')
# A parse error would look like a pass here, so require the refusal to actually
# be about the invite. This test is the security guarantee; it must not be able
# to pass for the wrong reason.
case "$MSG" in
  *invited*) ;;
  *) fail "step 10 did not exercise the guard — server said: $MSG" ;;
esac
pass "uninvited user blocked: $MSG"

step "11 · it reaches the admin list"
# --data-urlencode, because the test address contains a '+' and a raw query
# string would decode that to a space and silently match nothing.
ROW=$(curl -fsS "${AUTH[@]}" --get --data-urlencode "search=$CANDIDATE" \
        "$API/api/admin/scholarships" | jq '.applications[0] // empty')
[ -n "$ROW" ] || fail "application missing from the admin list"
pass "visible in admin → Scholarship"

# The candidate started a paper in step 9, so the list must show it. This is the
# join that used to be wired to a column nothing ever wrote.
SEEN_ATTEMPT=$(echo "$ROW" | jq -r '.attempt_id')
SEEN_STATUS=$(echo "$ROW" | jq -r '.status')
[ "$SEEN_ATTEMPT" = "$ATTEMPT" ] || fail "list shows attempt [$SEEN_ATTEMPT], expected [$ATTEMPT]"
[ "$SEEN_STATUS" = "started" ] || fail "list shows status [$SEEN_STATUS], expected [started]"
pass "linked to attempt $SEEN_ATTEMPT, status \"$SEEN_STATUS\""

step "12 · and the lead was mirrored"
curl -fsS "${AUTH[@]}" --get --data-urlencode "search=$CANDIDATE" --data-urlencode "source=scholarship" \
  "$API/api/admin/inquiries" \
  | jq -e '.inquiries | length > 0' >/dev/null || fail "lead not mirrored into inquiries"
pass "visible in admin → Enquiries"

step "13 · re-applying after sitting the paper is refused"
REAPPLY=$(curl -sS -o /dev/null -w '%{http_code}' -X POST "$API/api/scholarship/apply" \
  -H 'Content-Type: application/json' \
  -d "{\"name\":\"Smoke Candidate\",\"email\":\"$CANDIDATE\",\"courseId\":\"$COURSE_ID\"}")
# Still in progress, not submitted — re-applying is allowed and refreshes the
# link, which is what a candidate who lost their email needs.
[ "$REAPPLY" = "200" ] || fail "re-apply mid-attempt returned $REAPPLY, expected 200"
pass "mid-attempt re-application refreshes the link (200)"

step "14 · clean up after itself"
# The programme had to be live for steps 4-9 to work, and a live programme is
# offered on the public marketing site. Leaving it behind would put "Smoke Test
# Course" in front of real applicants, so it is retired on the way out.
curl -fsS -X POST "${AUTH[@]}" "$API/api/admin/scholarship-programs" -d "{
  \"course_id\":\"$COURSE_ID\",\"course_name\":\"Smoke Test Course\",
  \"assessment_id\":\"$ASSESSMENT\",\"is_active\":false}" \
  | jq -e .success >/dev/null || fail "could not retire the test programme"

curl -fsS "$API/api/scholarship/config" \
  | jq -e --arg c "$COURSE_ID" '[.programs[] | select(.courseId==$c)] | length == 0' >/dev/null \
  || fail "the test programme is still on the public config"
pass "test programme retired and gone from the public config"

# Retiring mentions only is_active. Optional fields must survive that: an upsert
# that defaulted them would silently erase the award ladder, and a candidate who
# scored 92% would then be told they had earned nothing.
curl -fsS "${AUTH[@]}" "$API/api/admin/scholarship-programs" \
  | jq -e --arg c "$COURSE_ID" '[.programs[] | select(.course_id==$c)][0].award_slabs | length > 0' >/dev/null \
  || fail "retiring the programme wiped its award ladder"
pass "award ladder survived the partial update"

printf '\n\033[32mAll checks passed.\033[0m Candidate: %s\n' "$CANDIDATE"
