"""The coverage rule of the behavioural oracle (PLAN.md §8; docs/behavioural-oracle.md).

A fast build that is not byte-identical to its compat build passes only where
the oracle can vouch for it:

  1. locate: tools/bench/ast_isomorphism -diff aligns the two builds' ASTs and
     reports every place where they differ (a node with other flags or
     literal, a node inserted or deleted, a reference whose names do not
     correspond under one consistent renaming), with the function (or
     top-level statement) around it. A function with no change in it is
     still alpha-isomorphic to its counterpart.
  2. run: both builds run under V8 block coverage (NODE_V8_COVERAGE). A
     change ran when its range in the fast build has a non-zero count; a
     deletion, when its range in the compat build does.
  3. judge, per change:
       vouched  the change ran, the test checks something (its compat build
                passes, or passes at least one of its tests) and the fast build
                behaved the same (core.compare: same or same-renamed)
       blind    the change never ran, or ran in a build that checks nothing
                or whose behaviour differs only where nothing is checked:
                reported, never passed
       fail     the fast build behaved differently on something the test
                checks (a regression or behaviour difference), and the change
                ran (changes that did not run stay blind)
     A build is identical, isomorphic, vouched (every change vouched), blind
     (some change blind, none failed) or fail.

jsgo itself logs nothing: the structural diff of the outputs locates every
difference, including the knock-on effects of a fast-only decision on later
passes, without trusting the code under test to report its own changes.
"""
import bisect
import collections
import hashlib
import json
import os
import re
import subprocess

import core

ISO_PKG = "./tools/bench/ast_isomorphism"
IDENT_RE = re.compile(r"[A-Za-z_$][A-Za-z0-9_$]*")
COMMENT_RE = re.compile(r"/\*.*?\*/|//[^\n]*", re.S)


def tool_digest():
    d = os.path.join(core.REPO, "tools", "bench", "ast_isomorphism")
    rels = [f for f in os.listdir(d) if f.endswith(".go") and not f.endswith("_test.go")]
    return core.tree_digest(d, rels)[0][:16]


def astdiff_binary(cache):
    """Builds tools/bench/ast_isomorphism into the cache (once per source
    digest) and returns its path."""
    out = os.path.join(cache, "behave", "bin", f"astdiff-{tool_digest()}")
    if not os.path.exists(out):
        os.makedirs(os.path.dirname(out), exist_ok=True)
        tmp = out + f".tmp{os.getpid()}"
        subprocess.run(["go", "build", "-o", tmp, ISO_PKG], cwd=core.REPO, check=True)
        os.replace(tmp, out)
    return out


def extern_names(astdiff, cache):
    """The property names of the default externs (astdiff -extern-names),
    cached beside the binary."""
    path = astdiff + ".extern-names.txt"
    if not os.path.exists(path):
        p = subprocess.run([astdiff, "-extern-names"], capture_output=True, text=True, check=True)
        tmp = path + f".tmp{os.getpid()}"
        with open(tmp, "w") as f:
            f.write(p.stdout)
        os.replace(tmp, path)
    return path


def case_externs(case):
    """The --externs files of a case (paths relative to its root)."""
    out = []
    for i, f in enumerate(case.flags):
        if f == "--externs" and i + 1 < len(case.flags):
            out.append(case.flags[i + 1])
        elif f.startswith("--externs="):
            out.append(f.split("=", 1)[1])
    return out


def fixed_names_file(case, base, work):
    """The names that keep their spelling in an ADVANCED build of case: the
    default externs' property names and every identifier of the case's own
    externs files. Written into work; returns its path."""
    extra = set()
    for rel in case_externs(case):
        try:
            with open(os.path.join(case.root, rel), encoding="utf-8", errors="replace") as f:
                src = f.read()
        except OSError:
            continue
        extra |= set(IDENT_RE.findall(COMMENT_RE.sub(" ", src)))
    if not extra:
        return base
    path = os.path.join(work, "fixed-names.txt")
    with open(base) as f:
        names = f.read()
    with open(path, "w") as f:
        f.write(names)
        f.write("\n".join(sorted(extra)) + "\n")
    return path


def renames_properties(config, case):
    return config.startswith("advanced") and "--property_renaming=OFF" not in case.flags


def locate(astdiff, compat_js, fast_js, renamed_props, fixed):
    """Runs the diff. Returns its JSON, or {"error": message}."""
    argv = [astdiff, "-diff", "-props", "renamed" if renamed_props else "exact"]
    if renamed_props and fixed:
        argv += ["-fixed", fixed]
    p = subprocess.run(argv + [compat_js, fast_js], capture_output=True, text=True, timeout=600)
    if p.returncode != 0:
        return {"error": (p.stderr or p.stdout).strip()[:2000]}
    return json.loads(p.stdout)


class Coverage:
    """V8 block coverage of one script: the count at a point is that of the
    smallest range around it."""

    def __init__(self, ranges):
        # V8's ranges nest; sorted by (start, -end), a stack sweep finds the
        # innermost range around each point.
        self.ranges = sorted(ranges, key=lambda r: (r[0], -r[1]))
        self.starts = [r[0] for r in self.ranges]

    def counts(self, points):
        """{point: count or None} for the points (any order)."""
        out = {}
        stack = []
        k = 0
        for p in sorted(set(points)):
            while k < len(self.ranges) and self.ranges[k][0] <= p:
                stack.append(self.ranges[k])
                k += 1
            # Ranges that ended before p cannot contain a later point either.
            while stack and stack[-1][1] <= p:
                stack.pop()
            c = None
            for r in reversed(stack):
                if r[0] <= p < r[1]:
                    c = r[2]
                    break
            out[p] = c
        return out


def sample_points(start, end):
    step = max(1, (end - start) // 16)
    return sorted(set(range(start, max(start + 1, end), step)) | {max(start, end - 1)})


def ran_flags(spans, cov):
    """For each [start, end) span (or None), whether a point of it has a
    non-zero count."""
    pts = []
    for s in spans:
        if s is not None:
            pts += sample_points(*s)
    counts = cov.counts(pts) if cov is not None and pts else {}
    out = []
    for s in spans:
        if s is None or cov is None:
            out.append(False)
            continue
        out.append(any((counts.get(p) or 0) > 0 for p in sample_points(*s)))
    return out


def signal(o):
    """Whether an outcome checks something: it passes, or (for suites with
    per-test results) at least one of its tests passes."""
    return bool(o["pass"]) or core.passing_tests(o) > 0


def judge(diff, cmp, compat_o, fast_cov, compat_cov):
    """The coverage rule's verdicts for one differing build (see the module
    comment). diff is locate()'s result; fast_cov and compat_cov are Coverage
    objects (compat_cov may be None when nothing was deleted). Returns (the
    build's summary, [(change, verdict, ran)])."""
    if "error" in diff:
        return {"verdict": "blind", "reason": "diff-error", "error": diff["error"][:500]}, []
    if diff["isomorphic"]:
        return {"verdict": "isomorphic", "functions": diff["functions_b"]}, []
    changes = diff["changes"]
    fast_spans = [tuple(c["b"]) if c["kind"] != "delete" and c.get("b") else None for c in changes]
    compat_spans = [tuple(c["a"]) if c["kind"] == "delete" and c.get("a") else None for c in changes]
    ran_f = ran_flags(fast_spans, fast_cov)
    ran_c = ran_flags(compat_spans, compat_cov)
    ran = [rf if c["kind"] != "delete" else rc for c, rf, rc in zip(changes, ran_f, ran_c)]
    sig = signal(compat_o)
    if cmp in core.SIGNAL_DIFFS:
        reason = "behaved-differently"
        verdicts = ["fail" if r else "blind" for r in ran]
    elif not sig:
        reason = "no-signal"
        verdicts = ["blind"] * len(changes)
    elif cmp in ("same", "same-renamed"):
        reason = None
        verdicts = ["vouched" if r else "blind" for r in ran]
    else:  # weak, fast-only-pass: differs where the test checks nothing
        reason = "differs-unchecked"
        verdicts = ["blind"] * len(changes)
    counts = collections.Counter(verdicts)
    # Per changed function (unit of the fast build; of the compat build for a
    # deletion): fail > blind > vouched.
    rank = {"vouched": 0, "blind": 1, "fail": 2}
    units = {}
    for c, v in zip(changes, verdicts):
        key = ("b", c["unit_b"]) if c["kind"] != "delete" or c["unit_b"] >= 0 else ("a", c["unit_a"])
        if key not in units or rank[v] > rank[units[key]]:
            units[key] = v
    ucounts = collections.Counter(units.values())
    if counts.get("fail") or cmp in core.SIGNAL_DIFFS:
        # A behaviour difference fails the build even when no change shows as
        # run (a run that timed out writes no coverage).
        verdict = "fail"
    elif counts.get("blind"):
        verdict = "blind"
    else:
        verdict = "vouched"
    blind = []
    for c, v, r in zip(changes, verdicts, ran):
        if v == "blind" and len(blind) < 5:
            u = diff["units_b"][c["unit_b"]] if c["unit_b"] >= 0 else diff["units_a"][c["unit_a"]]
            blind.append(f"{c['kind']} {c['what'][:60]} in {u['name']} (line {u['line']}){'' if not r else ' ran'}")
    out = {"verdict": verdict, "changes": len(changes), "vouched": counts.get("vouched", 0),
           "blind": counts.get("blind", 0), "fail": counts.get("fail", 0), "functions": diff["functions_b"],
           "units": len(units), "units_vouched": ucounts.get("vouched", 0), "units_blind": ucounts.get("blind", 0),
           "units_fail": ucounts.get("fail", 0), "ran": sum(ran), "exact_positions": diff["exact_positions"],
           "kinds": dict(collections.Counter(c["kind"] for c in changes)), "blind_examples": blind}
    if reason:
        out["reason"] = reason
    return out, list(zip(changes, verdicts, ran))


def _norm(text):
    return IDENT_RE.sub(lambda m: m.group(0) if len(m.group(0)) > 3 else "\u03b9", text)


def signatures(diff, ctext, ftext):
    """A key per change that is the same wherever the same code changed the
    same way (in another build of the same library or harness): the
    normalised text of its unit in each build (identifiers of up to three
    characters, which renaming makes, as one placeholder), its kind, its own
    normalised text, and its rank among the changes of that unit with the same
    text."""
    def unit_text(units, idx, text):
        if idx < 0:
            return ""
        u = units[idx]
        return _norm(text[u["start"]:u["end"]])
    seen = collections.Counter()
    out = []
    for c in diff["changes"]:
        span = c.get("b") if c["kind"] != "delete" else c.get("a")
        own = _norm((ftext if c["kind"] != "delete" else ctext)[span[0]:span[1]]) if span else ""
        base = (unit_text(diff["units_a"], c["unit_a"], ctext), unit_text(diff["units_b"], c["unit_b"], ftext),
                c["kind"], own)
        seen[base] += 1
        h = hashlib.sha256()
        for part in base + (str(seen[base]),):
            h.update(part.encode("utf-8", "surrogateescape") + b"\0")
        out.append(h.hexdigest()[:20])
    return out
