#!/usr/bin/env python3
"""One fresh Go CLI in the identical kernel-enforced benchmark budget.

Useful for correctness/profiling before a paired screen. Without --oracle it is
explicitly ungated diagnostic evidence, never a measured compatibility win.
"""
import argparse
import json
from pathlib import Path
import sys
sys.dont_write_bytecode = True
from paired_cli import Engine, go_request, input_identity, save, sha

p = argparse.ArgumentParser(description=__doc__)
p.add_argument('--spec', required=True)
p.add_argument('--binary', required=True)
p.add_argument('--out', required=True)
p.add_argument('--image', required=True)
p.add_argument('--go-env', default='{}')
p.add_argument('--oracle', help='pinned stock Java artifact digest JSON')
p.add_argument('--readonly-bind', action='append', default=[])
a = p.parse_args()
a.binary = str(Path(a.binary).resolve())
spec = json.loads(Path(a.spec).read_text())
root = Path(a.out).resolve()
if root.exists():
    p.error('refusing to overwrite diagnostic evidence')
root.mkdir(parents=True)
identity = input_identity(spec)
a.readonly_bind = sorted(set(a.readonly_bind+list(spec['links'].values())))
save(root/'metadata.json', vars(a) | {'input_manifest': identity, 'spec': spec,
     'binary_sha256': sha(a.binary), 'mode': 'untimed diagnostic/correctness invocation; raw clock retained, not performance evidence'})
engine = Engine(a)
try:
    row = go_request(engine, root/'run', spec, a)
    if a.oracle:
        if row['digests'] != json.loads(Path(a.oracle).read_text()):
            raise ValueError('artifact mismatch against stock Java')
        row['oracle_gate'] = 'MATCH'
    else:
        row['oracle_gate'] = 'UNGATED diagnostic; compare complete artifacts before using as evidence'
    save(root/'result.json', row)
    print(json.dumps(row, indent=2))
finally:
    engine.close()
