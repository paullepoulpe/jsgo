"""The suite-independent engine of the behavioural oracle (docs/behavioural-oracle.md).

A suite (t262.py, closure_rt.py; later phases add theirs) turns its tests into
Case objects. The engine compiles a case with jsgo in compat mode and in fast
mode under one configuration, runs the uncompiled program (when the suite has
one) and both builds with tools/behave/runner.js, and compares the builds'
behaviour:

  behaviour   what a build did: whether it compiled (and its error codes if
              not), whether the run timed out, the exit status, the runner's
              record (completion, phase, error name and message, unhandled
              rejections, timer counts, test-suite results) and the program's
              stdout and stderr
  verdict     the suite's pass/fail rule applied to one run (Case.verdict)

Compat is the reference. compare() returns "same" when the two behaviours are
equal, "same-renamed" when they are equal once renamed identifiers are
normalised (normaliser()), and otherwise one of:

  regression      compat passes, fast fails
  behaviour       both pass, but the behaviours differ
  fast-only-pass  compat fails, fast passes
  weak            both fail, differently

Only the first two are differences on tests that check something (a test
that fails on the compat build carries no signal, PLAN.md §8); the others
are reported but do not make a run fail.
"""
import gzip
import hashlib
import json
import os
import re
import subprocess
import time

HERE = os.path.dirname(os.path.abspath(__file__))
REPO = os.path.dirname(os.path.dirname(HERE))
RUNNER = os.path.join(HERE, "runner.js")
PINS = os.path.join(HERE, "pins.json")

# The configurations of PLAN.md §8: SIMPLE and ADVANCED, ES5 and ES2015 output.
CONFIGS = {
    "simple-es5": ["-O", "SIMPLE", "--language_in", "ECMASCRIPT_NEXT", "--language_out", "ECMASCRIPT5"],
    "simple-es2015": ["-O", "SIMPLE", "--language_in", "ECMASCRIPT_NEXT", "--language_out", "ECMASCRIPT_2015"],
    "advanced-es5": ["-O", "ADVANCED", "--language_in", "ECMASCRIPT_NEXT", "--language_out", "ECMASCRIPT5"],
    "advanced-es2015": ["-O", "ADVANCED", "--language_in", "ECMASCRIPT_NEXT", "--language_out", "ECMASCRIPT_2015"],
}

# Variables that would change jsgo's mode or the Go runtime's defaults
# (tools/bench/bench.py drops the same), or node's.
ENV_DROP = ("GOGC", "GOMEMLIMIT", "GOMAXPROCS", "GODEBUG", "JSGO_FAST_MODE", "JSGO_FAST_FAULT",
            "NODE_OPTIONS", "NODE_V8_COVERAGE")
JAVA_ENV_DROP = ("JAVA_TOOL_OPTIONS", "_JAVA_OPTIONS", "JDK_JAVA_OPTIONS")

COMPILE_TIMEOUT = 300
RUN_TIMEOUT = 30
OUTPUT_LIMIT = 16 * 1024

ERROR_CODE_RE = re.compile(r"ERROR - \[(\w+)\]")
CRASH_RE = re.compile(r"^(?:panic: |fatal error: )", re.M)
IDENT_RE = re.compile(r"[A-Za-z_$][A-Za-z0-9_$]*")

# Words of engine and harness messages that are never treated as renamed
# identifiers, even when a build happens to use them as names.
ENGLISH = frozenset("""
a an and are as at be been but by can cannot could did do does for from has have if in into is it
its may must no not of on or should the then this to too was were will with without
""".split())


def read_pins():
    with open(PINS) as f:
        return json.load(f)


def default_cache():
    """The main checkout's .cache/ (shared by all worktrees), or $JSGO_CACHE."""
    if os.environ.get("JSGO_CACHE"):
        return os.environ["JSGO_CACHE"]
    common = subprocess.run(["git", "-C", REPO, "rev-parse", "--path-format=absolute", "--git-common-dir"],
                            capture_output=True, text=True, check=True).stdout.strip()
    return os.path.join(os.path.dirname(common), ".cache")


def clean_env(drop=ENV_DROP, extra=None):
    env = {k: v for k, v in os.environ.items() if k not in drop}
    env["TZ"] = "UTC"
    env["LC_ALL"] = "C.UTF-8"
    if extra:
        env.update(extra)
    return env


def sha256_file(path):
    h = hashlib.sha256()
    with open(path, "rb") as f:
        for chunk in iter(lambda: f.read(1 << 20), b""):
            h.update(chunk)
    return h.hexdigest()


def tree_digest(root, rels):
    """sha256 over the sorted lines "<path>\\0<sha256 of the file>\\n" of the
    files rels (relative to root). Returns (digest, number of files)."""
    h = hashlib.sha256()
    rels = sorted(set(r.replace(os.sep, "/") for r in rels))
    for rel in rels:
        h.update(f"{rel}\0{sha256_file(os.path.join(root, rel))}\n".encode())
    return h.hexdigest(), len(rels)


class Case:
    """One test of a suite, ready to compile and run.

    root        the directory compilations run in; inputs are relative to it
    inputs      the --js inputs, in order
    flags       extra compiler flags for this case (strictness, dependency mode,
                externs), added after the configuration's
    uncompiled  the uncompiled program's source text, or None when the suite
                has no uncompiled form
    run_opts    extra runner.js options (--env browser, --onload, --freeze-date)
    verdict     function(record, stdout) -> (passed, why): the suite's rule
    """

    def __init__(self, suite, cid, root, inputs, flags, uncompiled, run_opts, verdict, stratum=""):
        self.suite = suite
        self.id = cid
        self.root = root
        self.inputs = inputs
        self.flags = flags
        self.uncompiled = uncompiled
        self.run_opts = run_opts
        self.verdict = verdict
        self.stratum = stratum
        self.configs = None  # the configurations this case runs in (None: all)

    def runs_in(self, configs):
        return [c for c in configs if self.configs is None or c in self.configs]

    def compile_argv(self, compiler, config, mode, out_js):
        """The compiler command line. compiler is the jsgo binary, or a list
        (the Java command) for the Java check; Java gets no --fast_mode."""
        argv = list(compiler) if isinstance(compiler, (list, tuple)) else [compiler]
        if mode == "fast":
            argv.append("--fast_mode")
        if mode in ("compat", "fast"):
            argv += ["--num_parallel_threads", "1"]
        argv += CONFIGS[config] + self.flags
        for inp in self.inputs:
            argv += ["--js", inp]
        argv += ["--js_output_file", out_js]
        return argv

    def safe_id(self):
        return re.sub(r"[^A-Za-z0-9._-]", "_", self.id)


def _clip(s):
    if len(s) > OUTPUT_LIMIT:
        return s[:OUTPUT_LIMIT] + f"\n[... {len(s) - OUTPUT_LIMIT} more characters]"
    return s


def compile_case(jsgo, case, config, mode, out_js, fault=None):
    """Compiles case with jsgo. Returns {exit, errors, crash, secs, stderr}."""
    extra = {"JSGO_FAST_FAULT": fault} if fault else None
    argv = case.compile_argv(jsgo, config, mode, out_js)
    t0 = time.monotonic()
    try:
        p = subprocess.run(argv, cwd=case.root, env=clean_env(extra=extra), capture_output=True,
                           timeout=COMPILE_TIMEOUT)
        code = p.returncode
        err = p.stderr.decode("utf-8", "replace")
    except subprocess.TimeoutExpired as e:
        code = "timeout"
        err = (e.stderr or b"").decode("utf-8", "replace")
    secs = time.monotonic() - t0
    errors = sorted(set(ERROR_CODE_RE.findall(err)))
    crash = bool(CRASH_RE.search(err)) or code == "timeout"
    ok = code == 0 and os.path.exists(out_js)
    return {"ok": ok, "exit": code, "errors": errors, "crash": crash, "secs": round(secs, 3),
            "stderr": _clip(err) if not ok else ""}


def run_program(node, program, case, work_prefix, coverage_dir=None):
    """Runs one program under runner.js. Returns {exit, timeout, stdout, stderr,
    rec, secs}; rec is None when the runner wrote no record. With coverage_dir,
    node writes V8's block coverage there (NODE_V8_COVERAGE)."""
    result = work_prefix + ".rec.json"
    if os.path.exists(result):
        os.remove(result)
    argv = [node, RUNNER] + case.run_opts + ["--result", result, program]
    t0 = time.monotonic()
    timeout = False
    try:
        env = clean_env(extra={"NODE_V8_COVERAGE": coverage_dir} if coverage_dir else None)
        p = subprocess.run(argv, cwd=os.path.dirname(program), env=env, capture_output=True,
                           timeout=RUN_TIMEOUT)
        code, out, err = p.returncode, p.stdout, p.stderr
    except subprocess.TimeoutExpired as e:
        timeout = True
        code, out, err = None, e.stdout or b"", e.stderr or b""
    secs = time.monotonic() - t0
    rec = None
    if not timeout and os.path.exists(result):
        try:
            with open(result) as f:
                rec = json.load(f)
        except ValueError:
            rec = None
    return {"exit": code, "timeout": timeout, "stdout": _clip(out.decode("utf-8", "replace")),
            "stderr": _clip(err.decode("utf-8", "replace")), "rec": rec, "secs": round(secs, 3)}


def outcome(case, comp, run):
    """The outcome of one build (comp is None for the uncompiled program, run
    is None when the compilation failed), with the suite's verdict."""
    if comp is not None and not comp["ok"]:
        return {"compile": comp, "run": None, "pass": False, "why": "compile-error"}
    if run["timeout"]:
        passed, why = False, "timeout"
    else:
        passed, why = case.verdict(run["rec"], run["stdout"])
    return {"compile": comp, "run": run, "pass": passed, "why": why}


def identifiers(text):
    return set(IDENT_RE.findall(text))


def normaliser(build_text, source_ids):
    """A function that replaces, in a string, every identifier the build has
    and the sources do not (a name the compiler made up) by a placeholder
    numbered in order of first appearance. One normaliser is used for all the
    fields of a behaviour, so that the numbering is shared."""
    renamed = identifiers(build_text) - source_ids - ENGLISH
    table = {}

    def norm(s):
        if not s or not renamed:
            return s

        def sub(m):
            w = m.group(0)
            if w not in renamed:
                return w
            if w not in table:
                table[w] = f"⟨id{len(table) + 1}⟩"
            return table[w]
        return IDENT_RE.sub(sub, s)
    return norm


def behaviour(o, norm=None):
    """The comparable behaviour of one outcome, as a JSON-able list."""
    n = norm or (lambda s: s)
    comp = o.get("compile")
    if comp is not None and not comp["ok"]:
        return ["compile-error", comp["exit"] if comp["exit"] != "timeout" else "timeout", comp["errors"],
                comp["crash"]]
    r = o["run"]
    if r["timeout"]:
        return ["timeout"]
    rec = r["rec"]
    if rec is None:
        return ["no-record", r["exit"], n(r["stdout"]), n(r["stderr"])]
    err = rec.get("error") or {}
    suite = rec.get("suite")
    return ["ran", r["exit"], rec.get("completion"), rec.get("phase"), n(err.get("name")), n(err.get("message")),
            n(r["stdout"]), n(r["stderr"]),
            [[n(x.get("name")), n(x.get("message"))] for x in rec.get("rejections") or []],
            rec.get("timers"), n(json.dumps(suite, sort_keys=True)) if suite is not None else None]


def compare(compat, fast, compat_text="", fast_text="", source_ids=frozenset()):
    """Compares two outcomes (see the module comment)."""
    if behaviour(compat) == behaviour(fast):
        return "same"
    if behaviour(compat, normaliser(compat_text, source_ids)) == behaviour(fast, normaliser(fast_text, source_ids)):
        return "same-renamed"
    if compat["pass"] and not fast["pass"]:
        return "regression"
    if compat["pass"]:
        return "behaviour"
    if fast["pass"]:
        return "fast-only-pass"
    return "weak"


SIGNAL_DIFFS = ("regression", "behaviour")


def write_gz(path, data):
    os.makedirs(os.path.dirname(path), exist_ok=True)
    with gzip.open(path, "wb", compresslevel=6) as f:
        f.write(data)


def read_gz(path):
    with gzip.open(path, "rb") as f:
        return f.read()


# ---------------------------------------------------------------------------------------------
# Where two builds differ, and whether a run executed it (the seeded-fault analysis; a first
# piece of PLAN.md §8's coverage rule).

SEGMENT_RE = re.compile(r"[^;{}\n]*[;{}\n]?")


def segments(text):
    """The build cut after each `;`, `{`, `}` and newline: [(start, end)]."""
    return [(m.start(), m.end()) for m in SEGMENT_RE.finditer(text) if m.end() > m.start()]


TOKEN_RE = re.compile(r"[A-Za-z_$][\w$]*|\d[\w.]*|\"(?:[^\"\\]|\\.)*\"|'(?:[^'\\]|\\.)*'|\S")


def changed_ranges(base_text, text, source_ids):
    """The ranges of `text` (offsets into it) that differ from base_text once
    the identifiers that renaming may have made (those the sources do not
    have, and all of up to three characters) are replaced by one placeholder:
    [(start, end)]. The builds are compared statement by
    statement (cut after `;`, `{`, `}`), then token by token inside each
    changed run of statements, so that a range covers the changed tokens only.
    A deletion is the empty range where it happened."""
    import difflib

    def norm_ident(w):
        # Renaming makes names of up to three characters (a, Gb, $jb, ...),
        # which sources use too (i, x, len): those are placeholders as well.
        return w if w in source_ids and len(w) > 3 else "\u03b9"

    def norm(t, segs):
        return [IDENT_RE.sub(lambda m: norm_ident(m.group(0)), t[a:b]) for a, b in segs]

    def tokens(t, a, b):
        return [(m.start(), m.end(), m.group(0)) for m in TOKEN_RE.finditer(t, a, b)]

    def tnorm(toks):
        return [norm_ident(w) if IDENT_RE.fullmatch(w) else w for _, _, w in toks]

    bs, ts = segments(base_text), segments(text)
    sm = difflib.SequenceMatcher(None, norm(base_text, bs), norm(text, ts), autojunk=False)
    out = []
    for op, i1, i2, j1, j2 in sm.get_opcodes():
        if op == "equal":
            continue
        if j2 == j1:
            at = ts[j1][0] if j1 < len(ts) else len(text)
            out.append((at, at))
            continue
        a, b = ts[j1][0], ts[j2 - 1][1]
        if i2 == i1:
            out.append((a, b))
            continue
        bt = tokens(base_text, bs[i1][0], bs[i2 - 1][1])
        tt = tokens(text, a, b)
        tm = difflib.SequenceMatcher(None, tnorm(bt), tnorm(tt), autojunk=False)
        for top, k1, k2, l1, l2 in tm.get_opcodes():
            if top == "equal":
                continue
            if l2 > l1:
                out.append((tt[l1][0], tt[l2 - 1][1]))
            else:
                at = tt[l1][0] if l1 < len(tt) else b
                out.append((at, at))
    return out


def coverage_ranges(coverage_dir, url="program.js"):
    """V8 block coverage ranges [(start, end, count)] of the script `url`."""
    import glob
    ranges = []
    for path in glob.glob(os.path.join(coverage_dir, "*.json")):
        try:
            with open(path) as f:
                data = json.load(f)
        except ValueError:
            continue
        for script in data.get("result", []):
            if script.get("url") == url:
                for fn in script.get("functions", []):
                    ranges += [(r["startOffset"], r["endOffset"], r["count"]) for r in fn.get("ranges", [])]
    return ranges


def executed(ranges, start, end):
    """Whether a point of [start, end) has a non-zero count (the count of a
    point is that of the smallest range around it); up to 16 points are
    tried."""
    step = max(1, (end - start) // 16)
    points = set(range(start, max(start + 1, end), step)) | {max(start, end - 1)}
    for p in points:
        best = None
        for a, b, c in ranges:
            if a <= p < b and (best is None or b - a < best[0]):
                best = (b - a, c)
        if best is not None and best[1] > 0:
            return True
    return False
