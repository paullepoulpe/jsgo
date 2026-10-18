#!/usr/bin/env python3
"""Pooled GC tax: GOGC=100 minus GOGC=off with the Java-RSS budget, over the sweep (N=5) and the
paired series (N=10), as per-round paired differences; with a bootstrap 90% interval of the median."""
import json, random, statistics as st, sys, os
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from summarize import stored  # noqa

A, B = json.load(open("sweep/sweep.json")), json.load(open("taxpair/sweep.json"))
rnd = random.Random(1)
print("| Task | N | Wall GOGC=100 (median) | Paired Δ wall s, median [90% CI] | Δ wall % | Δ CPU s (% of CPU) | Stored jsgo / warm Java | Tax share of gap [CI] |")
print("|---|---|---|---|---|---|---|---|")
for task in ["typed-app", "bundle", "downlevel-simple", "checks-only", "single-100"]:
    pairs = []
    for R in (A, B):
        if task in R:
            r = R[task]["runs"]
            pairs += list(zip(r["gogc100"], r["offjava"]))
    d = [x["wall_s"] - y["wall_s"] for x, y in pairs]
    dp = [(x["wall_s"] - y["wall_s"]) / x["wall_s"] for x, y in pairs]
    dc = [x["cpu_s"] - y["cpu_s"] for x, y in pairs]
    cpu = st.median(x["cpu_s"] for x, _ in pairs)
    base = st.median(x["wall_s"] for x, _ in pairs)
    boots = sorted(st.median(rnd.choices(dp, k=len(dp))) for _ in range(2000))
    lo, hi = boots[100], boots[1899]
    sj, wj, _ = stored(task)
    if task == "single-100":
        sc = st.median(x["wall_s"] for x, _ in pairs) / sj
        sj, wj = sj * sc, wj * sc
    gap = sj - wj
    f = st.median(dp)
    print(f"| {task} | {len(d)} | {base:.2f} | {st.median(d):.2f} [{lo * base:.2f}, {hi * base:.2f}] | {100 * f:.1f}% [{100 * lo:.1f}, {100 * hi:.1f}] | "
          f"{st.median(dc):.2f} ({100 * st.median(dc) / cpu:.0f}%) | {sj:.2f} / {wj:.2f} | {100 * f * sj / gap:.0f}% [{100 * lo * sj / gap:.0f}, {100 * hi * sj / gap:.0f}] |")
