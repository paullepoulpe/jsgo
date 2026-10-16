#!/usr/bin/env python3
"""Runs a command and reports its wall time, CPU time and peak RSS.

usage: python3 -I tools/bench/measure.py <report-file> <cmd> [args...]

The command's stdout/stderr pass through. The report file gets one line:
wall_s user_s sys_s maxrss_mb   (from wait4 rusage of the child)
"""
import os, subprocess, sys, time

report, cmd = sys.argv[1], sys.argv[2:]
t0 = time.monotonic()
p = subprocess.Popen(cmd)
_, status, ru = os.wait4(p.pid, 0)
wall = time.monotonic() - t0
with open(report, "w") as f:
    f.write(f"{wall:.2f} {ru.ru_utime:.2f} {ru.ru_stime:.2f} {ru.ru_maxrss / 1024:.0f}\n")
sys.exit(os.waitstatus_to_exitcode(status))
