#!/usr/bin/env node
// Proves a scholarship candidate never sees their own score.
//
// The API check proves the server withholds it. This proves the part only a
// browser can: that the player ends on a confirmation instead of a score page,
// that the /result URL shows the same thing to anyone who types it in, and that
// no percentage, mark or pass/fail verdict is anywhere in the rendered text.
//
//   node scripts/scholarship-withheld-check.mjs
//
// Needs the stack up, knovate-web on :3000 and skillofied-app on :5173.

const WEB = process.env.WEB || 'http://localhost:3000';
const APP = process.env.APP || 'http://localhost:5173';
const PORT = Number(process.env.PORT || 9337);
const COURSE_ID = process.env.COURSE_ID || '5';
const CHROME = process.env.CHROME
  || '/Applications/Google Chrome.app/Contents/MacOS/Google Chrome';

const sleep = (ms) => new Promise((r) => setTimeout(r, ms));
const ok = (m) => console.log(`  \x1b[32m✓\x1b[0m ${m}`);
const step = (m) => console.log(`\n\x1b[1m${m}\x1b[0m`);
const die = (m) => { console.error(`  \x1b[31m✗\x1b[0m ${m}`); process.exit(1); };

step('1 · apply');
const email = `withheld+${Date.now()}@example.com`;
const applied = await (await fetch(`${WEB}/api/scholarship/apply`, {
  method: 'POST', headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ name: 'Withheld Check', email, courseId: COURSE_ID }),
})).json();
if (!applied.testUrl) die(`no testUrl returned: ${JSON.stringify(applied)}`);
ok(`applied as ${email}`);

const { spawn } = await import('node:child_process');
const chrome = spawn(CHROME, [
  '--headless=new', '--disable-gpu', '--no-sandbox',
  // The briefing asks for camera and microphone; grant them with synthetic
  // devices so the permission gate behaves as it does for a real candidate.
  '--use-fake-ui-for-media-stream', '--use-fake-device-for-media-stream',
  `--remote-debugging-port=${PORT}`, '--window-size=1280,900', 'about:blank',
], { stdio: 'ignore' });
process.on('exit', () => chrome.kill());

async function pageTarget() {
  for (let i = 0; i < 60; i++) {
    try {
      const list = await (await fetch(`http://localhost:${PORT}/json/list`)).json();
      const page = list.find((t) => t.type === 'page');
      if (page?.webSocketDebuggerUrl) return page.webSocketDebuggerUrl;
    } catch { /* not up yet */ }
    await sleep(250);
  }
  die('Chrome never exposed a debuggable page');
}

const ws = new WebSocket(await pageTarget());
await new Promise((res, rej) => { ws.onopen = res; ws.onerror = rej; });
let id = 0; const pending = new Map();
ws.onmessage = (ev) => {
  const m = JSON.parse(ev.data);
  if (m.id && pending.has(m.id)) { pending.get(m.id)(m); pending.delete(m.id); }
};
const send = (method, params = {}) => new Promise((res) => {
  const i = ++id; pending.set(i, res); ws.send(JSON.stringify({ id: i, method, params }));
});
const evaluate = async (expression) => {
  const r = await send('Runtime.evaluate', { expression, returnByValue: true, awaitPromise: true });
  return r.result?.result?.value;
};
async function waitFor(expr, label, tries = 90) {
  for (let i = 0; i < tries; i++) { if (await evaluate(expr)) return; await sleep(500); }
  die(`timed out waiting for ${label}`);
}
await send('Page.enable'); await send('Runtime.enable');

step('2 · claim the link and start the paper');
// Done over the API rather than through the briefing screen: the two-step
// briefing needs a real fullscreen grant, which headless Chrome will not give,
// and it is already covered by scholarship-browser-check.mjs. What is under
// test here is what happens *after* submit.
const API = process.env.API || 'http://localhost:8080';
const claimToken = new URL(applied.testUrl).searchParams.get('t');
const claimed = await (await fetch(`${API}/api/scholarship/claim`, {
  method: 'POST', headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ token: claimToken }),
})).json();
if (!claimed.token) die(`claim failed: ${JSON.stringify(claimed)}`);
ok('claimed a candidate session');

const gql = async (query, variables) => {
  const r = await fetch(`${API}/api/graphql`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${claimed.token}` },
    body: JSON.stringify({ query, variables }),
  });
  const j = await r.json();
  if (j.errors) die(`graphql: ${JSON.stringify(j.errors)}`);
  return j.data;
};

const started = await gql(
  `mutation($a:String!,$t:String){startAttempt(assessmentId:$a,inviteToken:$t){attemptId resultsWithheld}}`,
  { a: claimed.assessmentId, t: claimed.inviteToken });
const attemptId = started.startAttempt.attemptId;
if (started.startAttempt.resultsWithheld !== true) {
  die('the server did not mark this paper as withholding results');
}
ok(`attempt ${attemptId} — server says resultsWithheld=true`);

step('3 · sit in the player and submit');
await send('Page.navigate', { url: `${APP}/` });
await waitFor(`location.origin === '${APP}'`, 'the app to load');
await evaluate(`localStorage.setItem('token', ${JSON.stringify(claimed.token)});
                localStorage.setItem('isLoggedIn','true'); true`);
await send('Page.navigate', { url: `${APP}/placement/tests/attempt/${attemptId}` });
await waitFor(`document.body.innerText.includes('Not answered')`, 'the player to render');
ok('in the player');

const btn = (re) => `[...document.querySelectorAll('button')].find(b=>${re}.test(b.textContent))`;
await evaluate(`${btn('/submit|finish|end test/i')}?.click(); true`);
await sleep(800);
// The player asks to confirm; the dialog's own button carries the same words.
await evaluate(`[...document.querySelectorAll('button')].filter(b=>/submit|yes|confirm/i.test(b.textContent)).pop()?.click(); true`);

step('4 · what the candidate sees');
await waitFor(`/your test has been submitted/i.test(document.body.innerText)`, 'the submitted panel');
ok('player ends on "Your test has been submitted"');
if (await evaluate(`location.pathname.includes('/result/')`)) {
  die('the player still routed to the score page');
}
ok(`stayed in the player (${await evaluate('location.pathname')}) — no score route`);

// The real assertion: nothing scorish is on the page, wherever it came from.
const SCORE = String.raw`/\d+\s*%|\bmarks\b|\bnot passed\b|\bpassed\b|\bscored\b|\bcorrect\b|\bincorrect\b|\bskipped\b|scholarship of|\d+% off/i`;
const leak = await evaluate(`(${SCORE}.exec(document.body.innerText)||[null])[0]`);
if (leak) die(`the submitted panel leaked "${leak}"`);
ok('no percentage, mark, verdict or award anywhere in the text');
ok(`panel says: ${JSON.stringify(await evaluate(`document.querySelector('h1')?.textContent`))}`);

step('5 · typing the result URL directly');
await send('Page.navigate', { url: `${APP}/placement/tests/result/${attemptId}` });
await waitFor(`/your test has been submitted/i.test(document.body.innerText)`, 'the result URL to show the same panel');
ok('/result/:id shows the confirmation, not a score');
const leak2 = await evaluate(`(${SCORE}.exec(document.body.innerText)||[null])[0]`);
if (leak2) die(`the result page leaked "${leak2}"`);
ok('no score there either');

// SHOT=<path> writes a screenshot of the panel, for eyeballing the copy.
if (process.env.SHOT) {
  const { writeFileSync } = await import('node:fs');
  const shot = await send('Page.captureScreenshot', { format: 'png' });
  writeFileSync(process.env.SHOT, Buffer.from(shot.result.data, 'base64'));
  ok(`screenshot: ${process.env.SHOT}`);
}

console.log('\n\x1b[32mResults are withheld from the candidate end to end.\x1b[0m\n');
process.exit(0);
