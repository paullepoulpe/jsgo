"""test262's runtime suite for the behavioural oracle (docs/behavioural-oracle.md).

The checkout is test262 at the commit pinned in tools/behave/pins.json
(tools/behave/fetch.sh). This module reads each test's front matter, decides
which tests the oracle can use (`static_exclusion`, then the probe-derived
feature exclusions of tools/behave/test262/exclusions.json), draws the pinned
stratified sample (tools/behave/test262/sample.txt, `stratified_sample`), and
turns a sampled test into a Case:

  - The harness files (assert.js, sta.js, doneprintHandle.js for async tests,
    then the test's `includes`) and the test are compiled together, as
    separate --js inputs, so that ADVANCED renames the harness and the test
    consistently. A `raw` test is compiled alone.
  - Strictness: a test without flags runs in strict mode (test262 runs it in
    both modes; strict is the one Closure assumes, --strict_mode_input). The
    uncompiled program is then "use strict"; followed by the inputs, and the
    builds get --emit_use_strict. A noStrict or raw test runs sloppy: the
    builds get --strict_mode_input=false.
  - The uncompiled program is the inputs concatenated, run as one script.

The verdict follows test262's INTERPRETING.md: a negative runtime test passes
when the script throws the named error; an async test passes when it prints
Test262:AsyncTestComplete (and no failure); any other test passes when the
script completes without an uncaught exception.
"""
import hashlib
import os
import re

FRONT_RE = re.compile(r"/\*---(.*?)---\*/", re.S)

# The harness every non-raw test gets, in this order (INTERPRETING.md).
DEFAULT_INCLUDES = ("assert.js", "sta.js")
ASYNC_INCLUDE = "doneprintHandle.js"

# Static exclusions: decided from the path and front matter alone, before any
# run. (category, reason); the first matching rule wins.
STATIC_RULES = [
    ("dir:intl402",
     "ECMA-402 (Intl): a library outside the compiler's scope whose results "
     "depend on the engine's ICU data"),
    ("dir:staging",
     "staging/: proposals and engine-contributed tests not yet in the main suite"),
    ("dir:harness",
     "harness/: tests of test262's own include files, not of the language"),
    ("negative:parse",
     "the expected outcome is an early (parse-time) error: a parser test, with no "
     "runtime behaviour to compare (the parser is checked byte for byte against "
     "Java, and test262-parser-tests is in the corpus)"),
    ("negative:resolution",
     "the expected outcome is a module resolution error: a loader test"),
    ("flag:module",
     "ES module tests need a module loader and their _FIXTURE imports; Closure "
     "rewrites modules into one script, so import.meta, live bindings across "
     "files, top-level await and namespace objects do not carry over"),
    ("flag:CanBlock",
     "CanBlockIsTrue/CanBlockIsFalse tests need $262.agent worker threads "
     "(Atomics.wait)"),
    ("feature:host",
     "needs a host hook the runner does not provide: cross-realm "
     "($262.createRealm), IsHTMLDDA, host-gc-required ($262.gc), or $262.agent "
     "(atomicsHelper.js)"),
    ("fixture",
     "a _FIXTURE file imported by module tests, not a test"),
]
STATIC_REASONS = dict(STATIC_RULES)
HOST_FEATURES = {"cross-realm", "IsHTMLDDA", "host-gc-required"}
HOST_INCLUDES = {"atomicsHelper.js"}


def parse_front_matter(src):
    """The YAML front matter subset test262 uses: scalars, flow and block lists,
    and the `negative` map. Returns a dict (missing keys are absent)."""
    m = FRONT_RE.search(src)
    if not m:
        return {}
    meta = {}
    lines = m.group(1).split("\n")
    i = 0
    while i < len(lines):
        line = lines[i]
        i += 1
        if not line.strip() or line.lstrip().startswith("#") or line[:1] in (" ", "\t"):
            continue
        if ":" not in line:
            continue
        key, rest = line.split(":", 1)
        key = key.strip()
        rest = rest.strip()
        if rest.startswith("[") and rest.endswith("]"):
            meta[key] = [x.strip().strip("'\"") for x in rest[1:-1].split(",") if x.strip()]
        elif rest in ("|", ">", "|-", ">-", "|+", ">+") or (rest == "" and key != "negative"):
            # A block: a list ("- x"), or a block scalar (skipped).
            items = []
            while i < len(lines) and (lines[i][:1] in (" ", "\t") or not lines[i].strip()):
                s = lines[i].strip()
                if s.startswith("- ") and rest == "":
                    items.append(s[2:].strip().strip("'\""))
                i += 1
            if rest == "":
                meta[key] = items
        elif key == "negative":
            neg = {}
            while i < len(lines) and lines[i][:1] in (" ", "\t"):
                s = lines[i].strip()
                if ":" in s:
                    k, v = s.split(":", 1)
                    neg[k.strip()] = v.strip()
                i += 1
            meta["negative"] = neg
        else:
            meta[key] = rest.strip("'\"")
    return meta


class Test:
    """One test262 test file: its path relative to test/ and its front matter."""

    def __init__(self, rel, meta):
        self.rel = rel
        self.meta = meta
        self.flags = set(meta.get("flags", []))
        self.features = list(meta.get("features", []))
        self.includes = list(meta.get("includes", []))
        neg = meta.get("negative") or {}
        self.negative = (neg.get("phase"), neg.get("type")) if neg else None

    @property
    def strict(self):
        return "noStrict" not in self.flags and "raw" not in self.flags

    @property
    def is_async(self):
        return "async" in self.flags

    def stratum(self):
        parts = self.rel.split("/")
        return "/".join(parts[:3]) if len(parts) > 3 else "/".join(parts[:-1])

    def harness(self):
        """The harness files (relative to harness/), in load order."""
        if "raw" in self.flags:
            return []
        files = list(DEFAULT_INCLUDES)
        if self.is_async:
            files.append(ASYNC_INCLUDE)
        for inc in self.includes:
            if inc not in files:
                files.append(inc)
        return files


def static_exclusion(t):
    """The STATIC_RULES category that excludes test t, or None."""
    top = t.rel.split("/", 1)[0]
    if "_FIXTURE" in t.rel:
        return "fixture"
    if top in ("intl402", "staging", "harness"):
        return "dir:" + top
    if t.negative and t.negative[0] == "parse":
        return "negative:parse"
    if t.negative and t.negative[0] == "resolution":
        return "negative:resolution"
    if "module" in t.flags:
        return "flag:module"
    if "CanBlockIsTrue" in t.flags or "CanBlockIsFalse" in t.flags:
        return "flag:CanBlock"
    if HOST_FEATURES & set(t.features) or HOST_INCLUDES & set(t.includes):
        return "feature:host"
    return None


def load_tests(root):
    """Every test under <root>/test, sorted by path."""
    base = os.path.join(root, "test")
    tests = []
    for dirpath, dirnames, filenames in os.walk(base):
        dirnames.sort()
        for fn in sorted(filenames):
            if not fn.endswith(".js"):
                continue
            path = os.path.join(dirpath, fn)
            rel = os.path.relpath(path, base).replace(os.sep, "/")
            with open(path, encoding="utf-8", errors="surrogateescape") as f:
                src = f.read()
            tests.append(Test(rel, parse_front_matter(src)))
    tests.sort(key=lambda t: t.rel)
    return tests


def rank(seed, key):
    return hashlib.sha256(f"{seed}\0{key}".encode()).hexdigest()


def stratified_sample(items, size, seed, stratum, key, floor=2):
    """A deterministic stratified sample of about `size` items.

    Each stratum gets max(min(n, floor), round(size * n / total)) items, the
    proportional share with a small floor so that rare directories are
    represented; the shares are then scaled so they sum to `size` (or to all
    the items, if fewer). Within a stratum the items with the lowest
    sha256(seed, key) are taken."""
    groups = {}
    for it in items:
        groups.setdefault(stratum(it), []).append(it)
    total = len(items)
    if size >= total:
        return sorted(items, key=key)
    names = sorted(groups)
    alloc = {g: max(min(len(groups[g]), floor), round(size * len(groups[g]) / total)) for g in names}
    # Scale the proportional part down (or up) until the sum is `size`.
    lo, hi = 0.0, 4.0
    for _ in range(60):
        mid = (lo + hi) / 2
        s = sum(min(len(groups[g]), max(min(len(groups[g]), floor), round(mid * size * len(groups[g]) / total)))
                for g in names)
        if s > size:
            hi = mid
        else:
            lo = mid
    alloc = {g: min(len(groups[g]), max(min(len(groups[g]), floor), round(lo * size * len(groups[g]) / total)))
             for g in names}
    # Hand out what rounding left over, largest strata first.
    short = size - sum(alloc.values())
    for g in sorted(names, key=lambda g: (-len(groups[g]), g)):
        if short <= 0:
            break
        if alloc[g] < len(groups[g]):
            alloc[g] += 1
            short -= 1
    out = []
    for g in names:
        ranked = sorted(groups[g], key=lambda it: rank(seed, key(it)))
        out.extend(ranked[:alloc[g]])
    return sorted(out, key=key)


def verdict(test, rec, stdout):
    """test262's pass/fail rule for one run. rec is the runner's record (None
    when the run produced none: a timeout or a crash). Returns (passed, why)."""
    if rec is None:
        return False, "no-record"
    if rec.get("phase") == "parse":
        return False, "parse-error"
    err = rec.get("error") or {}
    if test.negative:
        want = test.negative[1]
        if rec.get("completion") != "throw" or rec.get("phase") != "script":
            return False, f"expected-{want}"
        if err.get("name") != want:
            return False, f"threw-{err.get('name')}-expected-{want}"
        return True, "threw-expected"
    if rec.get("completion") == "throw":
        return False, f"threw-{err.get('name')}"
    if test.is_async:
        if "Test262:AsyncTestFailure" in stdout:
            return False, "async-failure"
        if "Test262:AsyncTestComplete" not in stdout:
            return False, "async-incomplete"
    return True, "ok"


SUITE_DIR = os.path.join(os.path.dirname(os.path.abspath(__file__)), "test262")
EXTERNS = os.path.join(SUITE_DIR, "externs.js")
SAMPLE = os.path.join(SUITE_DIR, "sample.txt")
PER_CHANGE = os.path.join(SUITE_DIR, "per-change.txt")
EXCLUSIONS = os.path.join(SUITE_DIR, "exclusions.json")


def read_list(path):
    with open(path) as f:
        return [ln.strip() for ln in f if ln.strip() and not ln.startswith("#")]


def make_case(root, t):
    """The core.Case for test t of the checkout at root."""
    import core
    inputs = ["harness/" + h for h in t.harness()] + ["test/" + t.rel]
    flags = ["--externs", EXTERNS]
    if t.strict:
        flags.append("--emit_use_strict")
    else:
        flags.append("--strict_mode_input=false")
    parts = ['"use strict";\n'] if t.strict else []
    for inp in inputs:
        with open(os.path.join(root, inp), encoding="utf-8", errors="surrogateescape") as f:
            parts.append(f.read())
        parts.append("\n")
    return core.Case("test262", t.rel, root, inputs, flags, "".join(parts), [],
                     lambda rec, stdout, t=t: verdict(t, rec, stdout), stratum=t.stratum())


def cases(root, subset):
    """The cases of a subset: "per-change", "full" (the sample) or a file of
    test paths."""
    path = {"per-change": PER_CHANGE, "full": SAMPLE}.get(subset, subset)
    out = []
    for rel in read_list(path):
        p = os.path.join(root, "test", rel)
        with open(p, encoding="utf-8", errors="surrogateescape") as f:
            out.append(make_case(root, Test(rel, parse_front_matter(f.read()))))
    return out
