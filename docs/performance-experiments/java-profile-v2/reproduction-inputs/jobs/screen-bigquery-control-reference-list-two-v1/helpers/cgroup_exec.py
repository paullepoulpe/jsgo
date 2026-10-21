#!/usr/bin/env python3
"""Verify the kernel budget, remove injected JVM options, then exec the JVM."""
import os
from pathlib import Path
import sys

root = Path('/sys/fs/cgroup')
if (root/'memory.max').read_text().strip() != '4294967296' or (root/'memory.swap.max').read_text().strip() != '0':
    raise SystemExit('real 4GiB memory / zero swap cgroup required')
quota, period = (root/'cpu.max').read_text().split()
if quota == 'max' or int(quota) != 2*int(period) or (root/'cpuset.cpus.effective').read_text().strip() not in ('0-1', '0,1'):
    raise SystemExit('kernel CPU quota2 and affinity0,1 required')
for key in ('JAVA_TOOL_OPTIONS', '_JAVA_OPTIONS', 'JDK_JAVA_OPTIONS', 'CLASSPATH'):
    os.environ.pop(key, None)
os.execv(sys.argv[1], sys.argv[1:])
