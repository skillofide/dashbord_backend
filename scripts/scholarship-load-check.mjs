#!/usr/bin/env node
// Concurrency dry run for the scholarship funnel.
//
// Answers one question: what happens when a college lab of candidates applies
// and starts at the same minute, from one public IP? That is the realistic
// shape of a campus drive, and it is the shape most likely to trip a limiter
// that was sized for a contact form.
//
//   node scripts/scholarship-load-check.mjs            # 50 candidates
//   N=120 node scripts/scholarship-load-check.mjs      # a bigger cohort
//   JUDGE=10 node scripts/scholarship-load-check.mjs   # + concurrent code runs
//
// Env: API, COURSE_ID, ADMIN_EMAIL, ADMIN_PASS
const API = process.env.API || 'http://localhost:8080';
const N = Number(process.env.N || 50);
const JUDGE = Number(process.env.JUDGE || 0);
const COURSE_ID = process.env.COURSE_ID || '5';
const ADMIN_EMAIL = process.env.ADMIN_EMAIL || 'admin@skillofied.com';
const ADMIN_PASS = process.env.ADMIN_PASS || 'skillofied123';

const ok = (m) => console.log(`  \x1b[32m✓\x1b[0m ${m}`);
const warn = (m) => console.log(`  \x1b[33m!\x1b[0m ${m}`);
const bad = (m) => console.log(`  \x1b[31m✗\x1b[0m ${m}`);
const step = (m) => console.log(`\n\x1b[1m${m}\x1b[0m`);

const pct = (xs, p) => {
  if (!xs.length) return 0;
  const s = [...xs].sort((a, b) => a - b);
  return Math.round(s[Math.min(s.length - 1, Math.floor((p / 100) * s.length))]);
};
const summarise = (label, ms) =>
  `${label}: p50 ${pct(ms, 50)}ms · p95 ${pct(ms, 95)}ms · max ${Math.round(Math.max(...ms, 0))}ms`;

async function timed(fn) {
  const t = Date.now();
  const value = await fn();
  return { value, ms: Date.now() - t };
}

const gql = (token, query, variables) =>
  fetch(`${API}/api/graphql`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
    body: JSON.stringify({ query, variables }),
  }).then((r) => r.json());

// ── one candidate's whole journey ────────────────────────────────────────────
async function candidate(i, assessmentId) {
  const email = `load+${Date.now()}-${i}@example.com`;
  const out = { email, stage: 'apply', rateLimited: false };

  const a = await timed(() =>
    fetch(`${API}/api/scholarship/apply`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ name: `Load ${i}`, email, courseId: COURSE_ID }),
    }));
  out.applyMs = a.ms;
  if (a.value.status === 429) { out.rateLimited = true; out.error = 'rate limited'; return out; }
  const applied = await a.value.json().catch(() => ({}));
  if (!applied.testUrl) { out.error = applied.error || `apply ${a.value.status}`; return out; }

  out.stage = 'claim';
  const token = applied.testUrl.split('t=')[1];
  const c = await timed(() =>
    fetch(`${API}/api/scholarship/claim`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ token }),
    }));
  out.claimMs = c.ms;
  if (c.value.status === 429) { out.rateLimited = true; out.error = 'rate limited'; return out; }
  const claimed = await c.value.json().catch(() => ({}));
  if (!claimed.token) { out.error = claimed.error || `claim ${c.value.status}`; return out; }

  out.stage = 'start';
  const s = await timed(() => gql(claimed.token,
    'mutation S($a:String!,$t:String){startAttempt(assessmentId:$a,inviteToken:$t){attemptId questions{id kind}}}',
    { a: assessmentId, t: claimed.inviteToken }));
  out.startMs = s.ms;
  const started = s.value?.data?.startAttempt;
  if (!started?.attemptId) { out.error = s.value?.errors?.[0]?.message || 'start failed'; return out; }

  out.stage = 'done';
  out.attemptId = started.attemptId;
  out.jwt = claimed.token;
  out.questionCount = started.questions.length;
  out.codingId = started.questions.find((q) => q.kind === 'coding')?.id;
  return out;
}

// ── go ───────────────────────────────────────────────────────────────────────
step('0 · warm up');
const login = await (await fetch(`${API}/api/login`, {
  method: 'POST', headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ email: ADMIN_EMAIL, password: ADMIN_PASS }),
})).json();
const assessments = await (await fetch(`${API}/api/recruiter/assessments?purpose=scholarship`, {
  headers: { Authorization: `Bearer ${login.token}` },
})).json();
const assessmentId = (assessments.assessments || []).find((a) => a.status === 'published')?.id;
if (!assessmentId) { bad('no published scholarship paper — run seed-scholarship-test.sh'); process.exit(1); }
ok(`paper ${assessmentId}`);

step(`1 · ${N} candidates apply, claim and start at once`);
const t0 = Date.now();
const results = await Promise.all(Array.from({ length: N }, (_, i) => candidate(i, assessmentId)));
const wall = Date.now() - t0;

const done = results.filter((r) => r.stage === 'done');
const limited = results.filter((r) => r.rateLimited);
const failed = results.filter((r) => r.stage !== 'done' && !r.rateLimited);

console.log(`  wall clock: ${wall}ms for ${N} candidates`);
console.log(`  ${summarise('apply', results.map((r) => r.applyMs).filter(Boolean))}`);
console.log(`  ${summarise('claim', results.map((r) => r.claimMs).filter(Boolean))}`);
console.log(`  ${summarise('start', results.map((r) => r.startMs).filter(Boolean))}`);

if (done.length === N) ok(`all ${N} reached a live paper`);
else console.log(`  ${done.length}/${N} reached a live paper`);

if (limited.length) {
  warn(`${limited.length}/${N} were RATE LIMITED — from a single IP, which is exactly what a`);
  console.log(`    college lab behind one NAT address looks like.`);
}
if (failed.length) {
  bad(`${failed.length} failed for other reasons:`);
  const byErr = {};
  failed.forEach((f) => { byErr[`${f.stage}: ${f.error}`] = (byErr[`${f.stage}: ${f.error}`] || 0) + 1; });
  Object.entries(byErr).forEach(([e, n]) => console.log(`      ${n}× ${e}`));
}

const sizes = new Set(done.map((r) => r.questionCount));
if (sizes.size === 1) ok(`every paper materialised with ${[...sizes][0]} questions`);
else bad(`papers came out different sizes: ${[...sizes].join(', ')}`);

const ids = new Set(done.map((r) => r.attemptId));
if (ids.size === done.length) ok('every candidate got their own attempt (no id collisions)');
else bad(`${done.length - ids.size} duplicate attempt ids`);

// ── the judge ────────────────────────────────────────────────────────────────
if (JUDGE > 0 && done.length) {
  step(`2 · ${JUDGE} concurrent code runs against the judge`);
  const pool = done.filter((r) => r.codingId).slice(0, JUDGE);
  const runs = await Promise.all(pool.map((r) => timed(() => gql(r.jwt,
    'mutation R($a:String!,$q:String!,$l:String!,$c:String!){runAttemptCode(attemptId:$a,questionId:$q,language:$l,code:$c){overallStatus runtimeMs testResults{status error}}}',
    { a: r.attemptId, q: r.codingId, l: 'python', c: 'def solveChallenge(input_val):\n    return input_val\n' },
  ))));
  console.log(`  ${summarise('run', runs.map((x) => x.ms))}`);

  // "A response came back" is not the same as "the code ran". A missing runner
  // image answers instantly with RuntimeError and no output, which reads as a
  // fast, healthy judge unless the verdict itself is checked.
  const noResponse = runs.filter((x) => !x.value?.data?.runAttemptCode);
  const executed = runs.filter((x) => {
    const r = x.value?.data?.runAttemptCode;
    return r && (r.testResults || []).some((t) => !/No such image/i.test(t.error || ''));
  });
  const missingImage = runs.filter((x) =>
    (x.value?.data?.runAttemptCode?.testResults || []).some((t) => /No such image/i.test(t.error || '')));

  if (noResponse.length) {
    bad(`${noResponse.length}/${pool.length} runs returned no result`);
    console.log(`      e.g. ${noResponse[0].value?.errors?.[0]?.message}`);
  }
  if (missingImage.length) {
    bad(`${missingImage.length}/${pool.length} never executed — the language runner image is missing.`);
    console.log('      Build them: ./scripts/build-runners.sh');
    console.log('      Until then every coding answer scores zero, however correct it is.');
  } else if (executed.length === pool.length) {
    ok(`all ${pool.length} runs actually executed`);
    const statuses = {};
    runs.forEach((x) => {
      const st = x.value.data.runAttemptCode.overallStatus;
      statuses[st] = (statuses[st] || 0) + 1;
    });
    console.log(`      verdicts: ${Object.entries(statuses).map(([k, v]) => `${v}× ${k}`).join(', ')}`);
  }
}

step('summary');
console.log(`  reached a paper: ${done.length}/${N}   rate limited: ${limited.length}   other failures: ${failed.length}`);
process.exit(failed.length > 0 ? 1 : 0);
