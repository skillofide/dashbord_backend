# Scholarship — email setup and how to check the flow

Two things: getting email working, and proving the funnel end to end.

---

## 1. Email

### How it is wired

Four emails come out of this stack, all from the same five environment variables:

| Email | Sent when | Goes to |
|---|---|---|
| Your test is ready | an application is submitted | the candidate |
| New scholarship application | an application is submitted | `SCHOLARSHIP_NOTIFY_TO` |
| You've been awarded N% | an admin presses **Award** | the candidate |
| New enquiry | the contact form is submitted | `INQUIRY_NOTIFY_TO` |

```
SMTP_HOST   the relay's hostname          ← unset = email disabled, links only logged
SMTP_PORT   587 (STARTTLS) or 465/1025
SMTP_USER   username, or blank for a relay that needs no auth
SMTP_PASS   password / app password / API key
SMTP_FROM   the From address
```

`SMTP_HOST` unset is a **supported state**, not a broken one: the funnel keeps working and the
candidate's link is written to the gateway log instead. That is what local development ran on
until now.

### Local development — a real inbox, no account needed

`docker compose up -d` now starts **Mailpit**, and the gateway points at it by default. Every
email the stack sends lands there and nothing can reach a real address.

```bash
docker compose up -d
open http://localhost:8025          # the inbox
```

Nothing to configure. To check it is on:

```bash
docker compose logs api-gateway | grep "email enabled"
#  enquiry email enabled       to=["admissions@knovate.local"]
#  scholarship email enabled   staff=["admissions@knovate.local"]
#  password reset email enabled
```

### Production

Put the real values in `.env` beside `docker-compose.yml`; the compose file reads them and they
override the Mailpit defaults.

**Gmail / Google Workspace** — fine for low volume, needs an *App Password*, not the account
password (2-Step Verification must be on: Google Account → Security → App passwords):

```ini
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=admissions@yourdomain.com
SMTP_PASS=<16-character app password>
SMTP_FROM=admissions@yourdomain.com
```

**Amazon SES** — the right choice at campaign volume, and you are already on AWS:

```ini
SMTP_HOST=email-smtp.ap-south-1.amazonaws.com
SMTP_PORT=587
SMTP_USER=<SES SMTP username>      # from SES → SMTP settings, NOT your AWS key
SMTP_PASS=<SES SMTP password>
SMTP_FROM=admissions@yourdomain.com
```

Two SES gotchas that bite everyone: a new account is **sandboxed** and can only send to verified
addresses — request production access before a drive; and `SMTP_FROM` must be a verified identity.

**Brevo / Mailgun / Postmark** all work the same way — host, port 587, username, key.

Then:

```bash
docker compose up -d --no-deps api-gateway
docker compose logs api-gateway | grep "email enabled"
```

### Before a real drive

- [ ] **SPF and DKIM** on your sending domain, or the test links land in spam and candidates never
      see them. This is the single most common reason a campaign underperforms.
- [ ] `SMTP_FROM` on a domain you control — never `@gmail.com` with a custom domain's branding.
- [ ] `SCHOLARSHIP_NOTIFY_TO` set to a monitored inbox.
- [ ] `APP_BASE_URL` set to the **public** portal origin. The test link is built from it; if it
      still says `localhost:5173`, every candidate gets a dead link.
- [ ] Send one application to a real address you own and open the link from a phone.

### Two things worth knowing

Mail is **fire-and-forget**. The application is committed before any email is attempted, so a slow
or broken relay can never fail or delay a submission — a delivery failure is logged loudly and the
candidate still has their link on screen.

Credentials are only offered when `SMTP_USER` is set. Relays that need no auth — an internal
smarthost, SES on an allow-listed IP, Mailpit — work without pretending to log in.

---

## 2. Checking the flow

### The fast way — three commands

```bash
./scripts/seed-scholarship-test.sh          # a paper + a course mapped to it   (once)
./scripts/scholarship-smoke.sh              # 14 API checks, cleans up after itself
node scripts/scholarship-browser-check.mjs  # drives real Chrome through the hand-off
```

Green on all three means: applying provisions an account and an invite, the link mints a session,
an **uninvited** user is refused, the paper materialises, the lead reaches both admin screens, and
a partial programme update does not wipe the award ladder.

Under load:

```bash
JUDGE=8 node scripts/scholarship-load-check.mjs     # 50 candidates at once + 8 code runs
```

### The honest way — do it yourself, once

Scripts prove contracts. Only a person notices that a sentence is confusing or a button is in the
wrong place. Twenty minutes, in this order:

1. **Apply.** http://localhost:3000/scholarship → *Apply and start the test*. Use a real-looking
   name and an address you can watch in Mailpit.
2. **Read the email** at http://localhost:8025. Is the link right? Does the wording make sense to
   someone who has never heard of you?
3. **Follow the link** (from the email, not the redirect — that is the path most candidates take).
   You should land on the instructions with the paper's real duration and section counts.
4. **Read the rules screen as a nervous 21-year-old.** Then start.
5. **In the paper:** answer some MCQs. On a coding question, drag the panes, run the sample cases,
   put something in *Custom input*, submit once, and check the *Submissions* tab.
6. **Leave and come back.** Move to another question and return — your code must still be there.
   Then hard-refresh the browser. It must still be there.
7. **Submit** and read the result page. The award banner is the payoff: does it tell you plainly
   what you earned?
8. **Admin** → http://localhost:5174 → Scholarship. Your candidate should be listed with a live
   score. Open them, press **Award**, and check the award email in Mailpit.

### If something looks wrong

```bash
docker compose logs -f api-gateway | grep -i scholarship   # applications, claims, rate limits
docker compose logs -f execution-service                   # the judge
docker compose ps                                          # anything not Up is the answer
```

```sql
-- what the funnel actually recorded
SELECT name, email, course_id, status, created_at
  FROM scholarship_applications ORDER BY created_at DESC LIMIT 10;

-- which courses are live on the marketing site right now
SELECT course_id, course_name, is_active, seats, award_slabs FROM scholarship_programs;
```

Common ones:

| Symptom | Cause |
|---|---|
| "that course is not open for scholarship applications" | no active programme for that `course_id`, or its paper is not published |
| Candidate's link 404s | `APP_BASE_URL` points somewhere the portal is not |
| Coding answers all score 0 | runner images missing — `./scripts/build-runners.sh` |
| Correct code returns TimeLimitExceeded | judge oversubscribed — lower `EXEC_MAX_CONCURRENT`, raise `EXEC_STARTUP_GRACE_MS` |
| Applicants get 429 during a campaign | raise `SCHOLARSHIP_IP_LIMIT`; a whole campus shares one IP |
| No email at all | `SMTP_HOST` unset — the link is in the gateway log |
