#!/usr/bin/env python3
"""Per-pass and per-phase roll-ups of a gclab pass log and GC log.

  python3 -I passana.py PREFIX   (PREFIX-passes.tsv, PREFIX-gclog.tsv)
"""
import collections, csv, os, sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from common import WT  # noqa


def phase_map():
    pm = {}
    for line in open(os.path.join(WT, "tools", "bench", "phases.tsv")):
        if line.startswith("#") or not line.strip():
            continue
        k, v = line.rstrip("\n").split("\t")
        pm[k] = v
    return pm


def rows(path):
    with open(path) as f:
        return list(csv.DictReader(f, delimiter="\t"))


def main():
    pre = sys.argv[1]
    pm = phase_map()
    P = rows(pre + "-passes.tsv")
    total = [r for r in P if r["name"] == "#total"][0]
    P = [r for r in P if r["name"] != "#total" and r["depth"] == "0"]
    tb, to = int(total["alloc_bytes"]), int(total["alloc_objs"])
    byp = collections.defaultdict(lambda: [0, 0, 0.0, 0, 0])
    byph = collections.defaultdict(lambda: [0, 0, 0.0, 0])
    for r in P:
        n = r["name"]
        x = byp[n]
        x[0] += int(r["alloc_bytes"]); x[1] += int(r["alloc_objs"]); x[2] += float(r["dur_ms"]); x[3] += 1
        x[4] += int(r["gc_cycles"])
        ph = pm.get(n, "other passes")
        y = byph[ph]
        y[0] += int(r["alloc_bytes"]); y[1] += int(r["alloc_objs"]); y[2] += float(r["dur_ms"]); y[3] += int(r["gc_cycles"])
    sb = sum(v[0] for v in byp.values())
    so = sum(v[1] for v in byp.values())
    print(f"total alloc {tb / 2**20:.0f} MiB, {to / 1e6:.2f} M objects; inside tracers {100 * sb / tb:.0f}% / {100 * so / to:.0f}%; "
          f"GC cycles {total['gc_cycles']}, GC CPU {float(total['gc_cpu_s']):.2f}s, wall in process {float(total['dur_ms']) / 1000:.2f}s")
    print("\n## by phase\nphase\tMiB\t%bytes\tMobj\t%objs\tms\tMiB/s\tgc_cycles")
    for ph, y in sorted(byph.items(), key=lambda kv: -kv[1][0]):
        print(f"{ph}\t{y[0] / 2**20:.0f}\t{100 * y[0] / tb:.1f}\t{y[1] / 1e6:.2f}\t{100 * y[1] / to:.1f}\t{y[2]:.0f}\t{y[0] / 2**20 / max(y[2], 1) * 1000:.0f}\t{y[3]}")
    print(f"outside tracers (startup, printing, untraced)\t{(tb - sb) / 2**20:.0f}\t{100 * (tb - sb) / tb:.1f}\t{(to - so) / 1e6:.2f}\t{100 * (to - so) / to:.1f}")
    print("\n## top passes by bytes\npass\truns\tMiB\t%bytes\tMobj\t%objs\tms\tMiB/s\tB/obj\tgc_cycles")
    for n, x in sorted(byp.items(), key=lambda kv: -kv[1][0])[:25]:
        print(f"{n}\t{x[3]}\t{x[0] / 2**20:.0f}\t{100 * x[0] / tb:.1f}\t{x[1] / 1e6:.2f}\t{100 * x[1] / to:.1f}\t{x[2]:.0f}\t"
              f"{x[0] / 2**20 / max(x[2], 1) * 1000:.0f}\t{x[0] / max(x[1], 1):.0f}\t{x[4]}")
    live = [(int(r["idx"]), r["name"], int(r["live_after"]), int(r["scan_after"])) for r in P if int(r["live_after"])]
    if live:
        pk = max(live, key=lambda t: t[2])
        print(f"\n## live heap after passes (forced GC): peak {pk[2] / 2**20:.0f} MiB after #{pk[0]} {pk[1]} (scannable {pk[3] / 2**20:.0f} MiB)")
        step = max(1, len(live) // 40)
        for t in live[::step]:
            print(f"#{t[0]}\t{t[1]}\t{pm.get(t[1], 'other')}\tlive {t[2] / 2**20:.0f} MiB\tscan {t[3] / 2**20:.0f} MiB ({100 * t[3] / max(t[2], 1):.0f}%)")
    g = pre + "-gclog.tsv"
    if os.path.exists(g):
        G = rows(g)
        if G:
            scan = sum(int(r["scan_heap"]) for r in G)
            liv = sum(int(r["live"]) for r in G)
            print(f"\n## GC log: {len(G)} cycles, Σ live at cycle end {liv / 2**30:.2f} GiB, Σ scannable {scan / 2**30:.2f} GiB "
                  f"({100 * scan / max(liv, 1):.0f}%), max live {max(int(r['live']) for r in G) / 2**20:.0f} MiB")
            byc = collections.defaultdict(lambda: [0, 0])
            for r in G:
                ph = pm.get(r["pass"], "startup" if r["pass"] == "startup" else "other passes")
                byc[ph][0] += 1
                byc[ph][1] += int(r["scan_heap"])
            print("phase\tcycles\tscanned GiB")
            for ph, (c, s) in sorted(byc.items(), key=lambda kv: -kv[1][1]):
                print(f"{ph}\t{c}\t{s / 2**30:.2f}")


if __name__ == "__main__":
    main()
