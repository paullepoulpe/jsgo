#!/usr/bin/env python3
"""Summaries of suite.py results: means and Java/jsgo ratios with 95% CIs.

  python3 -I analyze.py RESULTS.jsonl [--base java] [--vs final] [--json]

Cold comparisons are paired by block (paired Student-t on log ratios, df = n-1).
Comparisons with 'java-warm' (separate JVM iterations, not paired) use Welch's t on
log values. Ratios are geometric means of base/other: >1 means 'other' is faster/smaller.
"""
import json
import math
import statistics
import sys

# Two-sided 95% Student-t quantiles.
T975 = {1: 12.706, 2: 4.303, 3: 3.182, 4: 2.776, 5: 2.571, 6: 2.447, 7: 2.365, 8: 2.306, 9: 2.262,
        10: 2.228, 11: 2.201, 12: 2.179, 13: 2.160, 14: 2.145, 15: 2.131, 16: 2.120, 17: 2.110,
        18: 2.101, 19: 2.093, 20: 2.086, 25: 2.060, 30: 2.042, 40: 2.021, 60: 2.000, 120: 1.980}


def tq(df):
    df = max(1, int(math.floor(df)))
    keys = sorted(T975)
    for k in keys:
        if df <= k:
            return T975[k]
    return 1.96


def paired(a, b):
    logs = [math.log(x / y) for x, y in zip(a, b)]
    m = statistics.mean(logs)
    if len(logs) < 2:
        return math.exp(m), None, None
    h = tq(len(logs) - 1) * statistics.stdev(logs) / math.sqrt(len(logs))
    return math.exp(m), math.exp(m - h), math.exp(m + h)


def welch(a, b):
    la, lb = [math.log(x) for x in a], [math.log(x) for x in b]
    m = statistics.mean(la) - statistics.mean(lb)
    if len(la) < 2 or len(lb) < 2:
        return math.exp(m), None, None
    va, vb = statistics.variance(la) / len(la), statistics.variance(lb) / len(lb)
    se = math.sqrt(va + vb)
    if se == 0:
        return math.exp(m), math.exp(m), math.exp(m)
    df = (va + vb) ** 2 / (va ** 2 / (len(la) - 1) + vb ** 2 / (len(lb) - 1))
    h = tq(df) * se
    return math.exp(m), math.exp(m - h), math.exp(m + h)


def load(path):
    recs = [json.loads(l) for l in open(path) if l.strip()]
    out = {}
    for r in recs:
        out.setdefault(r["task"], {}).setdefault(r["who"], []).append(r)
    for t in out.values():
        for w in t:
            t[w].sort(key=lambda r: r["block"])
    return out


def compare(by, base, other, key):
    a, b = by[base], by[other]
    if base == "java-warm" or other == "java-warm":
        return welch([r[key] for r in a], [r[key] for r in b])
    n = min(len(a), len(b))
    return paired([r[key] for r in a[:n]], [r[key] for r in b[:n]])


def stage_compare(by, base, other, stage, key):
    a = [r["stages"][stage][key] for r in by[base]]
    b = [r["stages"][stage][key] for r in by[other]]
    n = min(len(a), len(b))
    return paired(a[:n], b[:n])


def main():
    path = sys.argv[1]
    data = load(path)
    out = {}
    for task, by in data.items():
        t = {"who": {}}
        for w, rs in by.items():
            t["who"][w] = {
                "n": len(rs),
                "wall_s": statistics.mean(r["wall_s"] for r in rs),
                "cpu_s": statistics.mean(r["cpu_s"] for r in rs),
                "rss_mb": statistics.mean(r["rss_mb"] for r in rs) if w != "java-warm" else rs[0]["rss_mb"],
                "gates": sorted({r["gate"].split(" ")[0] for r in rs}),
            }
            if rs and "stages" in rs[0]:
                t["who"][w]["stages"] = {
                    st: {k: statistics.mean(r["stages"][st][k] for r in rs) for k in ("wall_s", "cpu_s", "rss_mb")}
                    for st in rs[0]["stages"]
                }
        t["ratios"] = {}
        ws = list(by)
        for base in ("java", "java-warm", "main"):
            if base not in by:
                continue
            for other in ws:
                if other == base or other.startswith("java"):
                    continue
                if base == "main" and other == "main":
                    continue
                k = f"{base}/{other}"
                t["ratios"][k] = {key: compare(by, base, other, key) for key in ("wall_s", "cpu_s")}
                if base != "java-warm":
                    t["ratios"][k]["rss_mb"] = compare(by, base, other, "rss_mb")
                if "stages" in by[base][0] and base != "java-warm":
                    t["ratios"][k]["stages"] = {
                        st: {key: stage_compare(by, base, other, st, key) for key in ("wall_s", "rss_mb")}
                        for st in by[base][0]["stages"]
                    }
        out[task] = t
    if "--json" in sys.argv:
        print(json.dumps(out, indent=1))
        return
    for task, t in out.items():
        print(f"== {task}")
        for w, s in t["who"].items():
            print(f"   {w:10} n={s['n']} wall {s['wall_s']:8.2f}  cpu {s['cpu_s']:8.2f}  rss {s['rss_mb']:7.0f}  {s['gates']}")
        for k, v in t["ratios"].items():
            line = []
            for key in ("wall_s", "cpu_s", "rss_mb"):
                if key in v:
                    g, lo, hi = v[key]
                    line.append(f"{key} {g:.3f}" + (f" [{lo:.3f}, {hi:.3f}]" if lo else ""))
            print(f"   {k:18} " + "  ".join(line))


if __name__ == "__main__":
    main()
