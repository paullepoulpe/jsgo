// Node runner of the behavioural oracle (docs/behavioural-oracle.md).
//
//   node runner.js [--env script|browser] [--freeze-date] [--onload]
//                  [--pre <script.js>] [--max-timers N] --result <out.json>
//                  <program.js>
//
// Runs one program (an uncompiled test or a compiled build of it) as a classic
// script in this process's main realm (vm.runInThisContext: top-level `var`
// and function declarations become global properties, `this` is the global
// object, as in a browser <script> or a test262 host), after a deterministic
// prelude:
//
//   - Math.random is a fixed-seed generator;
//   - timers (setTimeout, setInterval, setImmediate and their clear
//     functions) run on a virtual clock: after the script, the runner fires
//     them in (due time, creation order), jumping the clock to each due time,
//     and drains the microtask queue (promise jobs) between two timers;
//   - with --freeze-date, Date reads the virtual clock (a fixed epoch plus the
//     virtual time), so Date.now() and `new Date()` are deterministic too.
//     It is off for test262, whose Date tests check the constructor's own
//     shape;
//   - `print` writes its arguments to stdout, and `$262` is the test262 host
//     object (global, evalScript, detachArrayBuffer);
//   - with --env browser, a minimal DOM stub (window, document, navigator,
//     location) for code written for a page, such as Closure Library's
//     goog.testing; with --onload, window.onload runs after the script;
//   - --pre runs a script before the program, in the same realm (Closure's
//     force_polyfills.js, which deletes built-ins so that a build's polyfills
//     run).
//
// The program's own output goes to stdout and stderr. The runner writes its
// record of what happened to --result as JSON:
//
//   completion   "normal", or "throw" when the script, a timer callback or a
//                microtask threw an exception that nothing caught
//   phase        where it was thrown: "parse" (the program is not valid
//                JavaScript for this engine), "script", "timer", "microtask"
//   error        {name, message}: name is the error's `name` for an Error, or
//                "Test262Error" for test262's error class (recognised by its
//                toString, so that renaming does not hide it), else the
//                constructor's name or the value's type; message is
//                String(e.message), or String(e) for a non-Error
//   rejections   the reasons of the promise rejections still unhandled at the
//                end, in the order they were rejected
//   timers       {fired, pending, capped}: timer callbacks run, timers left,
//                and whether --max-timers stopped the clock
//   suite        for --env browser, Closure's G_testRunner results if it ran:
//                {finished, success, tests: {name: [failure messages]}}
//
// Stack traces are not recorded. A program that never ends is killed by the
// caller (tools/behave/behave.py), which then has no record.
'use strict';

const fs = require('fs');
const vm = require('vm');

const args = process.argv.slice(2);
const opts = {env: 'script', freezeDate: false, onload: false, pre: null, maxTimers: 100000, result: null, file: null};
for (let i = 0; i < args.length; i++) {
  const a = args[i];
  if (a === '--env') opts.env = args[++i];
  else if (a === '--freeze-date') opts.freezeDate = true;
  else if (a === '--onload') opts.onload = true;
  else if (a === '--pre') opts.pre = args[++i];
  else if (a === '--max-timers') opts.maxTimers = Number(args[++i]);
  else if (a === '--result') opts.result = args[++i];
  else opts.file = a;
}
if (!opts.file || !opts.result) {
  process.stderr.write('usage: node runner.js [--env script|browser] [--freeze-date] [--onload] --result <out.json> <program.js>\n');
  process.exit(2);
}

// Host functions and built-ins the runner itself uses, taken before the
// prelude, --pre or the program can replace or delete them.
const G = globalThis;
const RealPromise = Promise;
const construct = Reflect.construct;
const objectKeys = Object.keys;
const imul = Math.imul;
const realSetImmediate = setImmediate;
const realWrite = (fd, s) => fs.writeSync(fd, s);
const RealDate = Date;
const EPOCH = RealDate.UTC(2026, 0, 1, 12, 0, 0);

const record = {completion: 'normal', phase: null, error: null, rejections: [], timers: null};

// describe(e): {name, message} of a thrown value, without running user code
// that could throw again.
function describe(e) {
  let name, message;
  try {
    if (e instanceof Error) {
      name = String(e.name);
      message = String(e.message);
    } else if (e !== null && typeof e === 'object') {
      let s = '';
      try { s = String(e); } catch (_) { /* uninspectable */ }
      if (s.slice(0, 14) === 'Test262Error: ') {
        name = 'Test262Error';
        message = 'message' in e ? String(e.message) : s.slice(14);
      } else {
        const c = e.constructor;
        name = c && typeof c.name === 'string' && c.name ? c.name : 'object';
        message = 'message' in e ? String(e.message) : s;
      }
    } else {
      name = typeof e;
      message = String(e);
    }
  } catch (_) {
    name = name || '<uninspectable>';
    message = message === undefined ? '<uninspectable>' : message;
  }
  return {name, message};
}

function fail(phase, e) {
  if (record.completion === 'throw') return;
  record.completion = 'throw';
  record.phase = phase;
  record.error = describe(e);
}

// ---------------------------------------------------------------------------
// Prelude

// Math.random: mulberry32 with a fixed seed. Assignment keeps the property's
// attributes; the function keeps the name and length of the original.
{
  let s = 0x9e3779b9 >>> 0;
  Math.random = {random() {
    s = (s + 0x6d2b79f5) >>> 0;
    let t = s;
    t = imul(t ^ (t >>> 15), t | 1);
    t ^= t + imul(t ^ (t >>> 7), t | 61);
    return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
  }}.random;
}

// Virtual timers.
let now = 0;
let nextId = 1;
let seq = 0;
const timers = new Map(); // id -> {due, seq, fn, args, interval}
function addTimer(fn, delay, rest, interval) {
  const id = nextId++;
  delay = Number(delay);
  if (!(delay > 0)) delay = 0;
  if (interval && delay < 1) delay = 1;
  timers.set(id, {due: now + delay, seq: seq++, fn, args: rest, interval: interval ? delay : 0});
  return id;
}
function clearTimer(id) {
  timers.delete(typeof id === 'object' && id !== null ? Number(id) : id);
}
globalThis.setTimeout = {setTimeout(fn, delay, ...rest) { return addTimer(fn, delay, rest, false); }}.setTimeout;
globalThis.setInterval = {setInterval(fn, delay, ...rest) { return addTimer(fn, delay, rest, true); }}.setInterval;
globalThis.setImmediate = {setImmediate(fn, ...rest) { return addTimer(fn, 0, rest, false); }}.setImmediate;
globalThis.clearTimeout = {clearTimeout(id) { clearTimer(id); }}.clearTimeout;
globalThis.clearInterval = {clearInterval(id) { clearTimer(id); }}.clearInterval;
globalThis.clearImmediate = {clearImmediate(id) { clearTimer(id); }}.clearImmediate;

if (opts.freezeDate) {
  const FakeDate = function Date(...a) {
    if (new.target === undefined) return RealDate(EPOCH + now);
    if (a.length === 0) return construct(RealDate, [EPOCH + now], new.target);
    return construct(RealDate, a, new.target);
  };
  Object.defineProperty(FakeDate, 'length', {value: 7});
  FakeDate.prototype = RealDate.prototype;
  Object.defineProperty(FakeDate, 'prototype', {writable: false, enumerable: false, configurable: false});
  FakeDate.now = {now() { return EPOCH + now; }}.now;
  FakeDate.parse = RealDate.parse;
  FakeDate.UTC = RealDate.UTC;
  Object.defineProperty(RealDate.prototype, 'constructor', {value: FakeDate, writable: true, enumerable: false, configurable: true});
  globalThis.Date = FakeDate;
}

globalThis.print = {print(...a) { realWrite(1, a.map(String).join(' ') + '\n'); }}.print;

globalThis.$262 = {
  global: globalThis,
  evalScript(src) { return vm.runInThisContext(String(src)); },
  detachArrayBuffer(buf) { structuredClone(buf, {transfer: [buf]}); return null; },
};

if (opts.env === 'browser') {
  installDomStub();
}

function installDomStub() {
  class Node {
    constructor(nodeName) {
      this.nodeName = nodeName;
      this.childNodes = [];
      this.parentNode = null;
    }
    appendChild(c) { this.childNodes.push(c); c.parentNode = this; return c; }
    removeChild(c) { const i = this.childNodes.indexOf(c); if (i >= 0) this.childNodes.splice(i, 1); c.parentNode = null; return c; }
    insertBefore(c, ref) { const i = this.childNodes.indexOf(ref); if (i < 0) return this.appendChild(c); this.childNodes.splice(i, 0, c); c.parentNode = this; return c; }
    get firstChild() { return this.childNodes[0] || null; }
    get lastChild() { return this.childNodes[this.childNodes.length - 1] || null; }
    get textContent() { return this.childNodes.map((c) => c.textContent).join(''); }
    set textContent(v) { this.childNodes = [new Text(String(v))]; }
  }
  class Text extends Node {
    constructor(data) { super('#text'); this.nodeType = 3; this.data = data; }
    get textContent() { return this.data; }
    set textContent(v) { this.data = String(v); }
  }
  class Element extends Node {
    constructor(tag) {
      super(String(tag).toUpperCase());
      this.nodeType = 1;
      this.tagName = this.nodeName;
      this.style = {};
      this.attributes = {};
      this.className = '';
      this.id = '';
      this.innerHTML = '';
    }
    setAttribute(k, v) { this.attributes[k] = String(v); }
    getAttribute(k) { return Object.prototype.hasOwnProperty.call(this.attributes, k) ? this.attributes[k] : null; }
    removeAttribute(k) { delete this.attributes[k]; }
    addEventListener() {}
    removeEventListener() {}
    getElementsByTagName() { return []; }
    querySelector() { return null; }
    querySelectorAll() { return []; }
  }
  const body = new Element('body');
  const head = new Element('head');
  const documentElement = new Element('html');
  documentElement.appendChild(head);
  documentElement.appendChild(body);
  const document = {
    nodeType: 9,
    title: 'behave',
    body, head, documentElement,
    readyState: 'complete',
    location: null,
    createElement: (t) => new Element(t),
    createTextNode: (s) => new Text(String(s)),
    createComment: (s) => new Text(String(s)),
    getElementById: () => null,
    getElementsByTagName: (t) => (String(t).toLowerCase() === 'head' ? [head] : String(t).toLowerCase() === 'body' ? [body] : []),
    querySelector: () => null,
    querySelectorAll: () => [],
    addEventListener() {},
    removeEventListener() {},
    write() {},
  };
  const location = {href: 'http://localhost/behave.html', protocol: 'http:', host: 'localhost',
    hostname: 'localhost', port: '', pathname: '/behave.html', search: '', hash: '', origin: 'http://localhost'};
  document.location = location;
  globalThis.window = globalThis;
  globalThis.self = globalThis;
  globalThis.top = globalThis;
  globalThis.parent = globalThis;
  globalThis.document = document;
  globalThis.location = location;
  Object.defineProperty(globalThis, 'navigator', {value: {userAgent: 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36', platform: 'Linux x86_64', language: 'en-US'}, writable: true, configurable: true});
  globalThis.addEventListener = () => {};
  globalThis.removeEventListener = () => {};
  globalThis.onload = null;
  globalThis.onerror = null;
}

// ---------------------------------------------------------------------------
// Run

const unhandled = new Map(); // promise -> reason, in rejection order
process.on('unhandledRejection', (reason, p) => { unhandled.set(p, reason); });
process.on('rejectionHandled', (p) => { unhandled.delete(p); });
process.on('uncaughtException', (e) => { fail('microtask', e); });

function tick() {
  return new RealPromise((resolve) => realSetImmediate(resolve));
}

async function main() {
  if (opts.pre) {
    vm.runInThisContext(fs.readFileSync(opts.pre, 'utf8'), {filename: 'pre.js'});
  }
  const code = fs.readFileSync(opts.file, 'utf8');
  let script = null;
  try {
    script = new vm.Script(code, {filename: 'program.js'});
  } catch (e) {
    fail('parse', e);
  }
  if (script) {
    try {
      script.runInThisContext();
    } catch (e) {
      fail('script', e);
    }
  }
  if (script && record.completion === 'normal' && opts.onload && typeof G.onload === 'function') {
    try {
      G.onload({type: 'load'});
    } catch (e) {
      fail('script', e);
    }
  }
  let fired = 0;
  let capped = false;
  while (record.completion === 'normal') {
    await tick();
    if (record.completion !== 'normal' || timers.size === 0) break;
    if (fired >= opts.maxTimers) { capped = true; break; }
    let bestId = 0;
    let best = null;
    // Map iteration is in insertion order, and (due, seq) is a total order.
    for (const [id, t] of timers) {
      if (best === null || t.due < best.due || (t.due === best.due && t.seq < best.seq)) {
        best = t;
        bestId = id;
      }
    }
    if (best.due > now) now = best.due;
    if (best.interval) {
      best.due = now + best.interval;
      best.seq = seq++;
    } else {
      timers.delete(bestId);
    }
    fired++;
    try {
      if (typeof best.fn === 'function') best.fn(...best.args);
      else vm.runInThisContext(String(best.fn));
    } catch (e) {
      fail('timer', e);
    }
  }
  await tick();
  record.timers = {fired, pending: timers.size, capped};
  record.rejections = [];
  for (const r of unhandled.values()) record.rejections.push(describe(r));
  if (opts.env === 'browser') {
    record.suite = closureSuite();
  }
  fs.writeFileSync(opts.result, JSON.stringify(record));
}

// closureSuite(): the results of Closure Library's test runner, if one ran.
function closureSuite() {
  const tr = G.G_testRunner;
  if (!tr || typeof tr.isFinished !== 'function') return null;
  const out = {finished: false, success: false, tests: null};
  try {
    out.finished = Boolean(tr.isFinished());
    out.success = Boolean(out.finished && tr.isSuccess());
    const results = typeof tr.getTestResults === 'function' ? tr.getTestResults() : null;
    if (results) {
      const tests = {};
      for (const name of objectKeys(results).sort()) {
        tests[name] = (results[name] || []).map((r) => String(r && r.message !== undefined ? r.message : r));
      }
      out.tests = tests;
    }
  } catch (e) {
    out.error = describe(e);
  }
  return out;
}

main().catch((e) => {
  fail('runner', e);
  try { fs.writeFileSync(opts.result, JSON.stringify(record)); } catch (_) { /* nothing to report to */ }
});
