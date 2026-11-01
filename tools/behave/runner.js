// Node runner of the behavioural oracle (docs/behavioural-oracle.md).
//
//   node runner.js [--env script|browser|qunit] [--freeze-date] [--date-step N]
//                  [--onload] [--pre <script.js>] [--max-timers N]
//                  --result <out.json> <program.js>
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
//     performance.now() reads that clock;
//   - with --freeze-date, Date reads the virtual clock (a fixed epoch plus the
//     virtual time), so Date.now() and `new Date()` are deterministic too.
//     It is off for test262, whose Date tests check the constructor's own
//     shape. With --date-step N, each such read also moves the virtual clock
//     on by N ms, so that code which waits for the clock in a loop (lodash's
//     throttle tests) ends, and timers fall due as they would;
//   - `print` writes its arguments to stdout, and `$262` is the test262 host
//     object (global, evalScript, detachArrayBuffer);
//   - with --env browser, a minimal DOM stub (window, document, navigator,
//     location) for code written for a page, such as Closure Library's
//     goog.testing; with --onload, window.onload runs after the script;
//   - with --env qunit, a QUnit subset (QUnit.module/test/skip/config/load/
//     start/done and the assertions of QUnit 1 and 2) for library test suites
//     written for QUnit (lodash). It lives here, uncompiled, so that the
//     checking code is not part of the build under test; externs declare its
//     names. Tests run one after the other on the virtual clock; an
//     exception in a timer callback while a test runs fails that test.
//     Failure messages name the assertion and its message argument, not the
//     values compared (renamed property names would show in those);
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
//   suite        for --env browser, Closure's G_testRunner results if it ran,
//                for --env qunit the QUnit subset's results:
//                {finished, success, tests: {name: [failure messages]}}
//
// Stack traces are not recorded. A program that never ends is killed by the
// caller (tools/behave/behave.py), which then has no record.
'use strict';

const fs = require('fs');
const vm = require('vm');

const args = process.argv.slice(2);
const opts = {env: 'script', freezeDate: false, dateStep: 0, onload: false, pre: null, maxTimers: 100000, result: null,
  file: null};
for (let i = 0; i < args.length; i++) {
  const a = args[i];
  if (a === '--env') opts.env = args[++i];
  else if (a === '--freeze-date') opts.freezeDate = true;
  else if (a === '--date-step') opts.dateStep = Number(args[++i]);
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

// performance.now() reads the virtual clock too (goog.testing times its tests
// with it).
if (typeof performance === 'object' && performance !== null) {
  Object.defineProperty(performance, 'now', {
    value: {now() { return now; }}.now, writable: true, enumerable: false, configurable: true});
}

if (opts.freezeDate) {
  const read = () => { const t = EPOCH + now; now += opts.dateStep; return t; };
  const FakeDate = function Date(...a) {
    if (new.target === undefined) return RealDate(read());
    if (a.length === 0) return construct(RealDate, [read()], new.target);
    return construct(RealDate, a, new.target);
  };
  Object.defineProperty(FakeDate, 'length', {value: 7});
  FakeDate.prototype = RealDate.prototype;
  Object.defineProperty(FakeDate, 'prototype', {writable: false, enumerable: false, configurable: false});
  FakeDate.now = {now() { return read(); }}.now;
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
let qunit = null;
if (opts.env === 'qunit') {
  qunit = installQUnit();
}

// installQUnit(): the QUnit subset of --env qunit. Returns {onError(e),
// results()}.
function installQUnit() {
  const vSetTimeout = G.setTimeout;
  const vClearTimeout = G.clearTimeout;
  const objToString = Object.prototype.toString;
  const getProto = Object.getPrototypeOf;
  const isNaNv = (x) => x !== x;
  const tests = [];
  const names = new Map();
  const results = {};
  const doneCbs = [];
  let moduleName = '';
  let moduleHooks = null;
  let started = false;
  let finished = false;
  let index = 0;
  let cur = null;

  function equiv(a, b, seen) {
    if (a === b) return true;
    const ta = typeof a;
    if (ta !== typeof b) return false;
    if (ta === 'number') return isNaNv(a) && isNaNv(b);
    if (a === null || b === null || ta !== 'object' && ta !== 'function') return false;
    if (ta === 'function') return false;
    const tag = objToString.call(a);
    if (tag !== objToString.call(b)) return false;
    for (const [x, y] of seen) if (x === a && y === b) return true;
    seen = seen.concat([[a, b]]);
    switch (tag) {
      case '[object Number]': case '[object String]': case '[object Boolean]':
        return equiv(a.valueOf(), b.valueOf(), seen);
      case '[object Date]': return equiv(a.getTime(), b.getTime(), seen);
      case '[object RegExp]': return a.source === b.source && a.flags === b.flags;
      case '[object Array]': case '[object Arguments]': {
        if (a.length !== b.length) return false;
        for (let i = 0; i < a.length; i++) if (!equiv(a[i], b[i], seen)) return false;
        return true;
      }
      case '[object Map]': case '[object Set]': {
        if (a.size !== b.size) return false;
        const xa = [...a], xb = [...b];
        for (let i = 0; i < xa.length; i++) if (!equiv(xa[i], xb[i], seen)) return false;
        return true;
      }
    }
    // QUnit compares constructors, unless one object has no prototype or
    // inherits straight from one with none.
    const pa = getProto(a), pb = getProto(b);
    if (pa !== pb) {
      const plain = (p) => p === null || p === Object.prototype || getProto(p) === null;
      if (!(plain(pa) && plain(pb)) && a.constructor !== b.constructor) return false;
    }
    const ka = [], kb = [];
    for (const k in a) ka.push(k);
    for (const k in b) kb.push(k);
    if (ka.length !== kb.length) return false;
    ka.sort(); kb.sort();
    for (let i = 0; i < ka.length; i++) {
      if (ka[i] !== kb[i] || !equiv(a[ka[i]], b[ka[i]], seen)) return false;
    }
    return true;
  }

  function ownEquiv(a, b) {
    const strip = (v, seen) => {
      if (v === null || typeof v !== 'object') return v;
      if (seen.includes(v)) return '[cycle]';
      const out = Array.isArray(v) ? [] : {};
      for (const k of objectKeys(v)) out[k] = strip(v[k], seen.concat([v]));
      return out;
    };
    return equiv(strip(a, []), strip(b, []), []);
  }

  function push(ok, kind, message) {
    if (!cur) return;
    cur.count++;
    if (!ok) cur.failures.push(kind + ' #' + cur.count + (message === undefined ? '' : ': ' + String(message)));
  }

  function makeAssert(t) {
    const a = {
      ok(v, m) { push(!!v, 'ok', m); },
      notOk(v, m) { push(!v, 'notOk', m); },
      equal(x, y, m) { push(x == y, 'equal', m); },
      notEqual(x, y, m) { push(x != y, 'notEqual', m); },
      strictEqual(x, y, m) { push(x === y, 'strictEqual', m); },
      notStrictEqual(x, y, m) { push(x !== y, 'notStrictEqual', m); },
      deepEqual(x, y, m) { push(equiv(x, y, []), 'deepEqual', m); },
      notDeepEqual(x, y, m) { push(!equiv(x, y, []), 'notDeepEqual', m); },
      propEqual(x, y, m) { push(ownEquiv(x, y), 'propEqual', m); },
      notPropEqual(x, y, m) { push(!ownEquiv(x, y), 'notPropEqual', m); },
      throws(fn, expected, m) {
        if (typeof expected === 'string') { m = expected; expected = undefined; }
        let threw = false, e;
        try { fn.call(undefined); } catch (err) { threw = true; e = err; }
        let ok = threw;
        if (threw && expected !== undefined) {
          if (typeof expected === 'function') {
            ok = expected.prototype !== undefined && e instanceof expected ? true
              : (expected.prototype === undefined || !(e instanceof Error)) && expected.call({}, e) === true;
          } else if (expected instanceof RegExp) {
            ok = expected.test(String(e));
          } else if (expected !== null && typeof expected === 'object') {
            ok = e instanceof expected.constructor && e.name === expected.name && e.message === expected.message;
          }
        }
        push(ok, 'throws', m);
      },
      expect(n) { t.expected = n; },
      async(n) {
        let left = n === undefined ? 1 : n;
        t.pending++;
        return function done() {
          if (left <= 0) return;
          if (--left === 0 && --t.pending === 0 && cur === t) finishSoon(t);
        };
      },
      pushResult(r) { push(!!(r && r.result), 'pushResult', r && r.message); },
      step(m) { t.steps.push(String(m)); push(true, 'step', m); },
      verifySteps(want, m) { push(equiv(t.steps, want, []), 'verifySteps', m); t.steps = []; },
      timeout() {},
    };
    a.raises = a.throws;
    a.same = a.deepEqual;
    a.equals = a.equal;
    return a;
  }

  function hook(t, names2) {
    const h = t.hooks;
    if (!h) return;
    for (const n of names2) {
      if (typeof h[n] === 'function') {
        try { h[n].call(t.env, t.assert); } catch (e) { t.failures.push('hook ' + n + ' died: ' + describe(e).name); }
      }
    }
  }

  function finishSoon(t) {
    if (t.timer !== null) vClearTimeout(t.timer);
    t.timer = null;
    vSetTimeout(() => finish(t), 0);
  }

  function finish(t) {
    if (t.done) return;
    t.done = true;
    if (t.timer !== null) vClearTimeout(t.timer);
    cur = t;
    hook(t, ['afterEach', 'teardown']);
    if (t.expected !== null && t.expected !== t.count) {
      t.failures.push('expected ' + t.expected + ' assertions, ran ' + t.count);
    }
    results[t.name] = t.failures;
    cur = null;
    vSetTimeout(runNext, 0);
  }

  function runNext() {
    if (index >= tests.length) {
      if (!finished) {
        finished = true;
        const failed = objectKeys(results).filter((k) => results[k].length).length;
        for (const cb of doneCbs) {
          try { cb({failed, passed: tests.length - failed, total: tests.length}); } catch (_) { /* reporting only */ }
        }
      }
      return;
    }
    const t = tests[index++];
    t.assert = makeAssert(t);
    cur = t;
    hook(t, ['beforeEach', 'setup']);
    try {
      const r = t.fn.call(t.env, t.assert);
      if (r && typeof r.then === 'function') {
        t.pending++;
        RealPromise.resolve(r).then(() => { if (--t.pending === 0 && cur === t) finishSoon(t); },
          (e) => { t.failures.push('rejected: ' + describe(e).name + ': ' + describe(e).message); if (--t.pending === 0 && cur === t) finishSoon(t); });
      }
    } catch (e) {
      const d = describe(e);
      t.failures.push('died: ' + d.name + ': ' + d.message);
      t.pending = 0;
    }
    if (t.pending === 0) {
      finish(t);
    } else {
      t.timer = vSetTimeout(() => { t.failures.push('timeout'); t.pending = 0; finish(t); }, 60000);
    }
  }

  function add(name, fn) {
    let full = (moduleName ? moduleName + ': ' : '') + String(name);
    const k = names.get(full) || 0;
    names.set(full, k + 1);
    if (k) full += ' #' + (k + 1);
    tests.push({name: full, fn, hooks: moduleHooks, env: {}, assert: null, failures: [], count: 0, expected: null,
      pending: 0, timer: null, done: false, steps: []});
  }

  G.QUnit = {
    config: {},
    module(name, hooks) {
      moduleName = String(name);
      moduleHooks = hooks && typeof hooks === 'object' ? hooks : null;
      if (typeof hooks === 'function') hooks({beforeEach() {}, afterEach() {}});
    },
    test(name, fn) { add(name, fn); },
    only(name, fn) { add(name, fn); },
    skip() {},
    todo() {},
    load() {},
    start() {
      if (started) return;
      started = true;
      vSetTimeout(runNext, 0);
    },
    done(cb) { doneCbs.push(cb); },
    begin() {}, testStart() {}, testDone() {}, moduleStart() {}, moduleDone() {}, log() {},
    equiv(a, b) { return equiv(a, b, []); },
  };
  return {
    onError(e) {
      if (!cur || cur.done) return false;
      const d = describe(e);
      cur.failures.push('uncaught: ' + d.name + ': ' + d.message);
      return true;
    },
    results() {
      const failed = objectKeys(results).filter((k) => results[k].length).length;
      const tests2 = {};
      for (const k of objectKeys(results).sort()) tests2[k] = results[k];
      return {finished, success: finished && failed === 0, tests: tests2, registered: tests.length};
    },
  };
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
process.on('uncaughtException', (e) => { if (!(qunit && qunit.onError(e))) fail('microtask', e); });

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
      if (!(qunit && qunit.onError(e))) fail('timer', e);
    }
  }
  await tick();
  record.timers = {fired, pending: timers.size, capped};
  record.rejections = [];
  for (const r of unhandled.values()) record.rejections.push(describe(r));
  if (opts.env === 'browser') {
    record.suite = closureSuite();
  } else if (qunit) {
    record.suite = qunit.results();
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
