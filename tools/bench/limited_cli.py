#!/usr/bin/env python3
"""Container-side fresh CLI measurement; kernel cgroup limit is mandatory."""
import json
import os
from pathlib import Path
import resource
import subprocess
import sys
import time

root = Path('/sys/fs/cgroup')
cap = {name: (root/name).read_text().strip() for name in
       ('memory.max', 'memory.swap.max', 'cpu.max', 'cpuset.cpus.effective', 'memory.events')}
if cap['memory.max'] != '4294967296' or cap['memory.swap.max'] != '0':
    raise SystemExit('requires real cgroup memory.max=4GiB and swap.max=0')
quota, period = cap['cpu.max'].split()
if quota == 'max' or int(quota) != 2*int(period) or cap['cpuset.cpus.effective'] not in ('0-1', '0,1'):
    raise SystemExit('requires actual kernel CPU quota2 and affinity0,1')
config = json.loads(Path(sys.argv[1]).read_text())
start = time.perf_counter_ns()
with open('stdout', 'wb') as out, open('stderr', 'wb') as err:
    process = subprocess.Popen(config['argv'], stdin=subprocess.DEVNULL, stdout=out, stderr=err,
                               env=config['environment'])
    pid, status, usage = os.wait4(process.pid, 0)
    code = os.waitstatus_to_exitcode(status)
    process.returncode = code
end = time.perf_counter_ns()
Path('exit').write_text(str(code)+'\n')
result = {'wall_s': (end-start)/1e9, 'cpu_s': usage.ru_utime+usage.ru_stime,
          'peak_rss_mib': usage.ru_maxrss/1024, 'exit': code, 'pid': pid,
          'cgroup': cap | {name: (root/name).read_text().strip() for name in
                           ('memory.peak', 'memory.current', 'memory.events')}}
Path('process.json').write_text(json.dumps(result, indent=2)+'\n')
raise SystemExit(code if code >= 0 else 128-code)
