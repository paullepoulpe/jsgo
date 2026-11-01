"""Closure Library's own unit tests for the behavioural oracle (docs/behavioural-oracle.md).

Closure Library (Apache 2.0) is the library written for ADVANCED, and its
goog.testing unit tests assert on the results of almost every module. The npm
package in the corpus leaves the tests out, so tools/behave/fetch.sh checks
out GitHub's google/closure-library at the tag of the corpus version
(v20230802, pins.json "closure_library_tests"); its sources are the corpus
package's, file for file.

Each `*_test.js` is a case compiled as Closure's runtime tests are
(closure_rt.py): --dependency_mode PRUNE with the test as entry point, the
files it needs found by a namespace scan, the generated test-name externs and
exports for ADVANCED, and run with the runner's DOM stub and frozen clock;
goog.testing's per-test results are the record, so each test of a file
carries its own signal (core.compare).

Not every test file can run under the stub: the pinned list
(tools/behave/closure_lib/tests.txt) holds the files at least one of whose
tests passes on the compat build in every configuration, from a probe of the
whole library minus the directories that need a real DOM throughout
(`behave.py run --suites closure-lib --subset probe`); per-change.txt is a
subset of it.
"""
import os

import closure_rt

HERE = os.path.dirname(os.path.abspath(__file__))
LIST_DIR = os.path.join(HERE, "closure_lib")
# Directories left out of the probe: widgets, the rich-text editor and the
# like need a real DOM (layout, selection, events) in nearly every test.
PROBE_SKIP = ("ui/", "editor/", "demos/", "graphics/", "fx/", "positioning/", "style/", "a11y/", "cssom/",
              "history/", "window/")
FLAGS = ["--dependency_mode", "PRUNE"]


def read_list(path):
    with open(path) as f:
        return [ln.split("#", 1)[0].strip() for ln in f if ln.split("#", 1)[0].strip()]


def cases(root, cache, subset="per-change"):
    """The cases of Closure Library's tests. root is the checkout (pins.json
    "closure_library_tests"."dir", under cache)."""
    import core
    rel_root = os.path.relpath(root, cache)
    goog = os.path.join(rel_root, "closure", "goog")
    index = closure_rt.Index(cache, rel_root, None)
    if subset == "probe":
        ids = []
        for dirpath, dirnames, filenames in os.walk(os.path.join(cache, goog)):
            dirnames.sort()
            for fn in sorted(filenames):
                if fn.endswith("_test.js"):
                    cid = os.path.relpath(os.path.join(dirpath, fn), os.path.join(cache, goog))
                    if not cid.startswith(PROBE_SKIP):
                        ids.append(cid)
    else:
        name = {"full": "tests.txt", "per-change": "per-change.txt"}.get(subset)
        ids = read_list(os.path.join(LIST_DIR, name) if name else subset)
    out = []
    for cid in ids:
        rel = os.path.join(goog, cid)
        inputs = index.closure(rel)
        flags = FLAGS + ["--entry_point", rel]
        with open(os.path.join(cache, rel), encoding="utf-8", errors="replace") as f:
            externs, exports = closure_rt.generated_inputs(cache, "lib/" + cid, f.read())
        flags += ["--externs", externs]
        if exports:
            inputs = inputs + [exports]
            flags += ["--entry_point", exports]
        if any(i.endswith("goog/testing/cspviolationobserver.js") for i in inputs):
            flags += ["--externs", os.path.join(goog, "testing/cspviolationobserver_externs.js")]
        out.append(core.Case("closure-lib", cid, cache, inputs, flags, None, closure_rt.RUN_OPTS, closure_rt.verdict,
                             stratum=os.path.dirname(cid) or "."))
    return out
