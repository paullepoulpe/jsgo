#!/usr/bin/env python3
"""Setup-only kernel-limit proof: read 4GiB limits, then a 64MiB OOM canary.

The canary never allocates 4GiB. No compiler workload or timing claim is made.
"""
import argparse
import json
import os
from pathlib import Path
import subprocess

p = argparse.ArgumentParser(description=__doc__)
p.add_argument('--image', required=True)
p.add_argument('--out', required=True)
a = p.parse_args()
if '@sha256:' not in a.image:
    p.error('pin the preinstalled image by digest')
env = {k: v for k, v in os.environ.items() if k not in
       ('DOCKER_HOST', 'DOCKER_CONTEXT', 'DOCKER_TLS', 'DOCKER_TLS_VERIFY', 'DOCKER_CERT_PATH')}
base = ['docker', '--host=unix:///var/run/docker.sock']
results = []
for limit, canary in ((4294967296, False), (67108864, True)):
    program = ("from pathlib import Path; import json; "
               "p=Path('/sys/fs/cgroup'); print(json.dumps({n:(p/n).read_text().strip() for n in "
               "['memory.max','memory.swap.max','cpu.max','cpuset.cpus.effective']}),flush=True)")
    if canary:
        program += "; b=bytearray(96<<20); exec('for i in range(0,len(b),4096): b[i]=1')"
    argv = ['create', '--memory='+str(limit), '--memory-swap='+str(limit), '--cpus=2',
            '--cpuset-cpus=0,1', '--network=none', '--read-only', a.image, 'python3', '-c', program]
    cid = subprocess.check_output(base+argv, env=env, text=True).strip()
    try:
        run = subprocess.run(base+['start', '-a', cid], env=env, capture_output=True, text=True)
        state = json.loads(subprocess.check_output(base+['inspect', cid], env=env, text=True))[0]
        kernel = json.loads(run.stdout.splitlines()[0])
        quota, period = kernel['cpu.max'].split()
        assert kernel['memory.max'] == str(limit) and kernel['memory.swap.max'] == '0'
        assert int(quota) == 2*int(period) and kernel['cpuset.cpus.effective'] in ('0-1', '0,1')
        assert state['State']['OOMKilled'] == canary
        assert state['State']['ExitCode'] == (137 if canary else 0)
        results.append({'limit_bytes': limit, 'kernel': kernel, 'oom_canary': canary,
                        'state': state['State'], 'stdout': run.stdout, 'stderr': run.stderr})
    finally:
        subprocess.run(base+['rm', '-f', cid], env=env, check=True, stdout=subprocess.DEVNULL)
out = Path(a.out)
out.parent.mkdir(parents=True, exist_ok=True)
out.write_text(json.dumps({'status': 'PASS', 'image': a.image, 'results': results,
                          'scope': 'actual kernel4GiB settings plus64MiB OOM enforcement canary; no benchmark'}, indent=2)+'\n')
print('PASS: kernel4GiB/swap0/CPU2/affinity0,1 and64MiB OOM canary')
