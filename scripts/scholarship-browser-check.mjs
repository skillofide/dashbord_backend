#!/usr/bin/env node
// Browser-level check of the scholarship hand-off, end to end.
//
// The shell smoke test proves the API contract; this proves the part only a
// browser can: that the link in a candidate's email actually lands them in a
// running paper. It drives real Chrome over CDP using the WebSocket built into
// Node 22+, so there is nothing to install.
//
//   node scripts/scholarship-browser-check.mjs
//
// Needs: the stack up, knovate-web on :3000, skillofied-app on :5173, and a
// scholarship programme mapped (./scripts/seed-scholarship-test.sh).
//
// Env: WEB=http://localhost:3000  PORT=9333  COURSE_ID=5  CHROME=<path>

const WEB = process.env.WEB || 'http://localhost:3000';
const PORT = Number(process.env.PORT || 9333);
const COURSE_ID = process.env.COURSE_ID || '5';
const CHROME = process.env.CHROME
  || '/Applications/Google Chrome.app/Contents/MacOS/Google Chrome';

const sleep = (ms) => new Promise((r) => setTimeout(r, ms));
const ok = (m) => console.log(`  \x1b[32m✓\x1b[0m ${m}`);
const step = (m) => console.log(`\n\x1b[1m${m}\x1b[0m`);
const die = (m) => { console.error(`  \x1b[31m✗\x1b[0m ${m}`); process.exit(1); };

// ── an application, straight through the marketing site's own proxy ──────────
step('1 · apply on the marketing site');
const email = `browser+${Date.now()}@example.com`;
const applyResp = await fetch(`${WEB}/api/scholarship/apply`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ name: 'Browser Check', email, courseId: COURSE_ID }),
});
const applied = await applyResp.json();
if (!applied.testUrl) die(`no testUrl returned: ${JSON.stringify(applied)}`);
ok(`applied as ${email}`);

// ── drive it ─────────────────────────────────────────────────────────────────
const { spawn } = await import('node:child_process');
const chrome = spawn(CHROME, [
  '--headless=new', '--disable-gpu', '--no-sandbox',
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

let id = 0;
const pending = new Map();
ws.onmessage = (ev) => {
  const msg = JSON.parse(ev.data);
  if (msg.id && pending.has(msg.id)) { pending.get(msg.id)(msg); pending.delete(msg.id); }
};
const send = (method, params = {}) => new Promise((res) => {
  const i = ++id; pending.set(i, res);
  ws.send(JSON.stringify({ id: i, method, params }));
});
const evaluate = async (expression) => {
  const r = await send('Runtime.evaluate', { expression, returnByValue: true, awaitPromise: true });
  return r.result?.result?.value;
};
async function waitFor(expr, label, tries = 80) {
  for (let i = 0; i < tries; i++) {
    if (await evaluate(expr)) return;
    await sleep(500);
  }
  die(`timed out waiting for ${label}`);
}

await send('Page.enable');
await send('Runtime.enable');

step('2 · follow the link from the email');
await send('Page.navigate', { url: applied.testUrl });
await waitFor(`location.pathname.startsWith('/scholarship/instructions')`, 'the claim redirect');
ok('one-time token exchanged for a session');
if (!(await evaluate(`!!localStorage.getItem('token') && localStorage.getItem('isLoggedIn')==='true'`))) {
  die('no session was stored — later authenticated calls would fail');
}
ok('session stored under the same keys the login form uses');

step('3 · the instructions gate');
await waitFor(`!!document.querySelector('input[type=checkbox]')`, 'the instructions to render');
ok(`paper: ${await evaluate(`document.querySelector('h1')?.textContent`)}`);

const startBtn = `[...document.querySelectorAll('button')].find(b=>/start my test/i.test(b.textContent))`;
if ((await evaluate(`${startBtn}?.disabled`)) !== true) {
  die('Start was enabled before consent — a candidate could begin without reading the rules');
}
ok('Start is disabled until the candidate consents');

await evaluate(`document.querySelector('input[type=checkbox]').click(); true`);
await sleep(400);
if ((await evaluate(`${startBtn}?.disabled`)) !== false) die('Start stayed disabled after consent');
ok('Start enables on consent');

step('4 · start the paper');
await evaluate(`${startBtn}.click(); true`);
await waitFor(`location.pathname.startsWith('/placement/tests/attempt/')`, 'the attempt to start');
ok(`in the player: ${await evaluate('location.pathname')}`);

// The palette lists every question, so its length is the paper size.
await waitFor(`document.body.innerText.includes('Not answered')`, 'the player to render');
const sections = await evaluate(
  `JSON.stringify([...document.querySelectorAll('button')].map(b=>b.textContent.trim()).filter(t=>/\\d+\\/\\d+$/.test(t)))`);
const parsed = JSON.parse(sections);
if (parsed.length < 2) die(`expected an MCQ and a coding section, saw ${sections}`);
parsed.forEach((s) => ok(`section rendered — ${s}`));

if (!(await evaluate(`Object.keys(sessionStorage).filter(k=>k.startsWith('scholarship.invite')).length === 0`))) {
  die('the invite token was left behind in sessionStorage');
}
ok('invite token consumed');

console.log(`\n\x1b[32mThe hand-off works end to end.\x1b[0m Candidate: ${email}`);
ws.close();
chrome.kill();
process.exit(0);
