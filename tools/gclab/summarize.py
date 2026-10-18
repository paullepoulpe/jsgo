#!/usr/bin/env python3
"""Tables from sweep/sweep.json: the GC settings sweep and the GC tax per task.

  python3 -I summarize.py sweep/sweep.json > sweep/summary.md
"""
import json, os, statistics as st, sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from common import REF  # noqa

LABEL = {"gogc100": "GOGC=100", "gogc200": "GOGC=200", "gogc400": "GOGC=400",
         "off10g": "off, 10 GiB", "offjava": "off, Java RSS"}


def med(rs, k):
    return st.median(r[k] for r in rs)


def rng(rs, k):
    v = [r[k] for r in rs]
    return min(v), max(v)


def stored(task):
    t = "single-files" if task == "single-100" else task
    j = json.load(open(os.path.join(REF, t, "jsgo.json")))["runs"]
    w = json.load(open(os.path.join(REF, t, "java-warm.json")))
    c = json.load(open(os.path.join(REF, t, "java-cold.json")))["runs"]
    wm = st.median(i["wall_ms"] for i in w["iterations"] if i["kind"] == "measure") / 1000
    return st.median(r["wall_s"] for r in j), wm, st.median(r["wall_s"] for r in c)


def main():
    res = json.load(open(sys.argv[1]))
    print("| Task | Setting | Wall s (min–max) | CPU s (user / sys) | Peak RSS MB | GC cycles | GC CPU s | GC share of CPU | Alloc GiB | Δ wall vs GOGC=100 | Outputs identical |")
    print("|---|---|---|---|---|---|---|---|---|---|---|")
    tax = []
    for task, v in res.items():
        runs = v["runs"]
        base = med(runs["gogc100"], "wall_s")
        jr = v["java_rss_mb"]
        jr = f"{st.median(jr):.0f} (median per file)" if isinstance(jr, list) else f"{jr:.0f}"
        for s, rs in runs.items():
            lo, hi = rng(rs, "wall_s")
            w = med(rs, "wall_s")
            same = sum(r["same_as_jsgo_ref"] for r in rs)
            samej = sum(r["same_as_java_ref"] for r in rs)
            tot = sum(r["jobs"] for r in rs)
            lab = LABEL[s] + (f" ({jr} MiB)" if s == "offjava" else "")
            print(f"| {task} | {lab} | {w:.2f} ({lo:.2f}–{hi:.2f}) | {med(rs, 'cpu_s'):.2f} ({med(rs, 'user_s'):.2f} / {med(rs, 'sys_s'):.2f}) | "
                  f"{med(rs, 'rss_mb'):.0f} | {med(rs, 'gc_cycles'):.0f} | {med(rs, 'gc_cpu_s'):.2f} | "
                  f"{100 * med(rs, 'gc_cpu_s') / med(rs, 'cpu_s'):.0f}% | {med(rs, 'alloc_bytes') / 2**30:.2f} | "
                  f"{100 * (w - base) / base:+.1f}% | {same}/{tot} (Java: {samej}/{tot}) |")
        sj, wj, cj = stored(task)
        if task == "single-100":
            # The stored numbers are for 1,000 files; scale them to this subset by jsgo's own time.
            scale = base / sj
            sj, wj, cj = base, wj * scale, cj * scale
        off = med(runs["off10g"], "wall_s") if "off10g" in runs else None
        offj = med(runs["offjava"], "wall_s") if "offjava" in runs else None
        tax.append((task, base, off, offj, sj, wj))
    print()
    print("| Task | Wall GOGC=100 | Wall off/10 GiB | GC tax (s, % of wall) | Wall off/Java RSS | Tax with Java-RSS budget | Stored jsgo / warm Java | Gap to warm Java | Tax share of the gap (10 GiB / Java RSS) |")
    print("|---|---|---|---|---|---|---|---|---|")
    for task, base, off, offj, sj, wj in tax:
        t1 = base - off
        t2 = base - offj
        f1, f2 = t1 / base, t2 / base
        gap = sj - wj
        print(f"| {task} | {base:.2f} | {off:.2f} | {t1:.2f} ({100 * f1:.0f}%) | {offj:.2f} | {t2:.2f} ({100 * f2:.0f}%) | "
              f"{sj:.2f} / {wj:.2f} | {gap:.2f} ({100 * gap / sj:.0f}% of jsgo) | {100 * f1 * sj / gap:.0f}% / {100 * f2 * sj / gap:.0f}% |")




def paired(path, a="gogc100", b="offjava"):
    """Per-run paired differences (the two settings ran back to back): wall, CPU, user."""
    res = json.load(open(path))
    print(f"\n| Task | N | Wall {a} (median) | Wall {b} | Paired Δ wall s: median (min–max) | Paired Δ wall % | Δ CPU s | Δ user s | Δ sys s |")
    print("|---|---|---|---|---|---|---|---|---|")
    for task, v in res.items():
        A, B = v["runs"][a], v["runs"][b]
        d = [x["wall_s"] - y["wall_s"] for x, y in zip(A, B)]
        dp = [100 * (x["wall_s"] - y["wall_s"]) / x["wall_s"] for x, y in zip(A, B)]
        dc = [x["cpu_s"] - y["cpu_s"] for x, y in zip(A, B)]
        du = [x["user_s"] - y["user_s"] for x, y in zip(A, B)]
        ds = [x["sys_s"] - y["sys_s"] for x, y in zip(A, B)]
        print(f"| {task} | {len(d)} | {med(A, 'wall_s'):.2f} | {med(B, 'wall_s'):.2f} | {st.median(d):.2f} ({min(d):.2f}–{max(d):.2f}) | "
              f"{st.median(dp):.1f}% | {st.median(dc):.2f} | {st.median(du):.2f} | {st.median(ds):.2f} |")


if __name__ == "__main__":
    if len(sys.argv) > 2 and sys.argv[2] == "--paired":
        paired(sys.argv[1])
    else:
        main()
        paired(sys.argv[1], "gogc100", "offjava")
        paired(sys.argv[1], "gogc100", "off10g")
