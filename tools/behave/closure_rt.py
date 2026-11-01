"""Closure Compiler's runtime tests for the behavioural oracle (docs/behavioural-oracle.md).

`test/com/google/javascript/jscomp/runtime_tests` is a set of Closure Library
jsunit tests written to run compiled (language features, ES modules,
goog.module, polyfills). GitHub's mirror of Closure Compiler stopped carrying
it after v20231112, so the oracle uses that tag (tools/behave/pins.json), with
the Closure Library release of the same period from the corpus
(google-closure-library@20230802.0.0) for goog.testing.

Each `*_test.js` is compiled the way the directory's utils/build.sh compiles
it (--dependency_mode PRUNE with the test as entry point,
--process_common_js_modules, --module_resolution NODE), with the library and
test files it needs found by a namespace and import scan (`deps`), and run
with the runner's browser stub: the test runner of goog.testing.jsunit starts
from window.onload, and its results are the verdict. The polyfill tests also
run in a "#forced" variant (ES5 output only): the directory's
force_polyfills.js deletes the engine's built-ins before the build runs, so
that the build's polyfills are the code under test.

goog.testing finds the tests by name (global functions for jsunit-style
tests, the keys of the object given to testSuite otherwise), which ADVANCED
renames. Closure's own test builds use CompilerOptions.exportTestFunctions,
which the command line does not expose, so each case gets two generated
inputs (under .cache/behave/closure-gen/, the same for every configuration):
an externs file whose object literal lists the test names (so that property
renaming keeps the testSuite keys), and for a jsunit-style test a second entry
point that calls goog.exportSymbol on each test function.

There is no uncompiled form: goog.module and ES modules need the compiler (or
a debug loader) to run at all.
"""
import os
import re

PROVIDE_RE = re.compile(r"""goog\.(?:provide|module|declareModuleId)\(\s*['"]([\w.$]+)['"]\s*\)""")
REQUIRE_RE = re.compile(r"""goog\.(?:require|requireType)\(\s*['"]([\w.$]+)['"]\s*\)""")
IMPORT_RE = re.compile(r"""(?:^|[;\s])(?:import|export)\b[^'";]*?['"]([^'"\n]+)['"]""", re.M)

TESTS_REL = "test/com/google/javascript/jscomp/runtime_tests"

# Tests that cannot run under the runner's DOM stub, with the reason.
SKIP = {
    "polyfill_tests/polyfill_isolation_chrome_test.js": "needs a real DOM (its _dom.html page)",
    "polyfill_tests/polyfill_isolation_ie_test.js": "needs a real DOM (its _dom.html page)",
}

FLAGS = ["--process_common_js_modules", "--module_resolution", "NODE", "--dependency_mode", "PRUNE"]
RUN_OPTS = ["--env", "browser", "--onload", "--freeze-date"]


def strip_comments(src):
    """The source without // and /* */ comments (a rough cut, string-unaware:
    enough to ignore commented-out requires)."""
    src = re.sub(r"/\*.*?\*/", " ", src, flags=re.S)
    return re.sub(r"(?m)(^|[^:\\'\"])//.*$", r"\1", src)


class Index:
    """Namespace -> file over the library and the test tree; file -> deps."""

    def __init__(self, cache, lib_rel, tests_rel):
        self.cache = cache
        self.provides = {}
        self.deps = {}
        self.base = os.path.join(lib_rel, "closure/goog/base.js")
        for top in (os.path.join(lib_rel, "closure/goog"), tests_rel):
            if top is None:
                continue
            for dirpath, dirnames, filenames in os.walk(os.path.join(cache, top)):
                dirnames.sort()
                for fn in sorted(filenames):
                    if not fn.endswith(".js"):
                        continue
                    rel = os.path.relpath(os.path.join(dirpath, fn), cache)
                    with open(os.path.join(dirpath, fn), encoding="utf-8", errors="replace") as f:
                        src = strip_comments(f.read())
                    for ns in PROVIDE_RE.findall(src):
                        self.provides.setdefault(ns, rel)
                    reqs = REQUIRE_RE.findall(src)
                    imps = []
                    for spec in IMPORT_RE.findall(src):
                        if spec.startswith("goog:"):
                            reqs.append(spec[5:])
                        elif spec.startswith("."):
                            imps.append(os.path.normpath(os.path.join(os.path.dirname(rel), spec)))
                    self.deps[rel] = (reqs, imps)

    def closure(self, entry):
        """The files entry needs, transitively (base.js included), sorted."""
        seen = {self.base}
        stack = [entry]
        while stack:
            f = stack.pop()
            if f in seen or f not in self.deps:
                continue
            seen.add(f)
            reqs, imps = self.deps[f]
            for ns in reqs:
                if ns in self.provides:
                    stack.append(self.provides[ns])
            for p in imps:
                for cand in (p, p + ".js"):
                    if cand in self.deps:
                        stack.append(cand)
                        break
        return [self.base] + sorted(seen - {self.base})


TEST_NAME = r"(?:test\w*|setUp|tearDown|setUpPage|tearDownPage|shouldRunTests)"
METHOD_RE = re.compile(r"(?m)^\s*(?:async\s+)?(" + TEST_NAME + r")\s*(?:\(|:)")
GLOBAL_FN_RE = re.compile(r"(?m)^(?:async\s+)?function\s*\*?\s*(" + TEST_NAME + r")\s*\(")


def test_names(src):
    """The test and lifecycle names a test file defines: (all names, the
    top-level functions of a jsunit-style test)."""
    src = strip_comments(src)
    globals_ = sorted(set(GLOBAL_FN_RE.findall(src)))
    names = sorted(set(METHOD_RE.findall(src)) | set(globals_))
    return names, globals_


def write_if_changed(path, text):
    os.makedirs(os.path.dirname(path), exist_ok=True)
    if os.path.exists(path):
        with open(path) as f:
            if f.read() == text:
                return
    with open(path, "w") as f:
        f.write(text)


def generated_inputs(cache, cid, src):
    """The generated externs file and, for a jsunit-style test, the export
    entry point (paths relative to cache)."""
    names, globals_ = test_names(src)
    gen = os.path.join("behave", "closure-gen", re.sub(r"[^A-Za-z0-9._-]", "_", cid))
    externs = os.path.join(gen, "test_names_externs.js")
    write_if_changed(os.path.join(cache, externs),
                     "/**\n * @fileoverview Generated by tools/behave/closure_rt.py: the test names, kept\n"
                     " * by ADVANCED's property renaming.\n * @externs\n */\n"
                     "/** @const */\nvar behaveTestNames = {\n" + "".join(f"  {n}: 0,\n" for n in names) + "};\n")
    exports = None
    if globals_ and "goog.testing.jsunit" in src:
        exports = os.path.join(gen, "export_tests.js")
        write_if_changed(os.path.join(cache, exports),
                         "// Generated by tools/behave/closure_rt.py: exports the test functions under their\n"
                         "// names, so that goog.testing finds them in ADVANCED builds.\n"
                         "goog.require('goog.testing.jsunit');\n" +
                         "".join(f"goog.exportSymbol('{n}', {n});\n" for n in globals_))
    return externs, exports


def verdict(rec, stdout):
    if rec is None:
        return False, "no-record"
    if rec.get("phase") == "parse":
        return False, "parse-error"
    if rec.get("completion") == "throw":
        return False, f"threw-{(rec.get('error') or {}).get('name')}"
    s = rec.get("suite")
    if not s:
        return False, "no-test-runner"
    if not s.get("finished"):
        return False, "not-finished"
    tests = s.get("tests") or {}
    failed = sum(1 for v in tests.values() if v)
    if not s.get("success") or failed:
        return False, f"test-failures-{failed}"
    if not tests:
        return False, "no-tests"
    return True, "ok"


def cases(root, cache):
    """The cases of Closure's runtime tests. root is the closure-compiler
    checkout (relative to cache, pins.json "closure"."dir")."""
    import core
    pins = core.read_pins()
    lib_rel = pins["closure_library"]["dir"]
    tests_rel = os.path.join(os.path.relpath(root, cache), TESTS_REL)
    index = Index(cache, lib_rel, tests_rel)
    out = []
    for dirpath, dirnames, filenames in os.walk(os.path.join(cache, tests_rel)):
        dirnames.sort()
        for fn in sorted(filenames):
            if not fn.endswith("_test.js"):
                continue
            rel = os.path.relpath(os.path.join(dirpath, fn), cache)
            cid = os.path.relpath(rel, tests_rel)
            if cid in SKIP:
                continue
            inputs = index.closure(rel)
            flags = FLAGS + ["--entry_point", rel]
            with open(os.path.join(dirpath, fn), encoding="utf-8", errors="replace") as f:
                externs, exports = generated_inputs(cache, cid, f.read())
            flags += ["--externs", externs]
            if exports:
                inputs = inputs + [exports]
                flags += ["--entry_point", exports]
            if any(i.endswith("goog/testing/cspviolationobserver.js") for i in inputs):
                flags += ["--externs", os.path.join(lib_rel, "closure/goog/testing/cspviolationobserver_externs.js")]
            out.append(core.Case("closure", cid, cache, inputs, flags, None, RUN_OPTS, verdict,
                                 stratum=os.path.dirname(cid) or "."))
            if cid.startswith("polyfill_tests/"):
                force = os.path.join(cache, os.path.dirname(rel), "force_polyfills.js")
                if os.path.exists(force):
                    c = core.Case("closure", cid + "#forced", cache, inputs, flags, None,
                                  RUN_OPTS + ["--pre", force], verdict, stratum=os.path.dirname(cid) + "#forced")
                    c.configs = [k for k in core.CONFIGS if k.endswith("-es5")]
                    out.append(c)
    return out
