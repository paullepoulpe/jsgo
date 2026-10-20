#!/usr/bin/env python3
"""Decode actual JS-produced API payloads with independent native protobuf."""
import argparse
import hashlib
import json
from pathlib import Path
import subprocess

p = argparse.ArgumentParser(description=__doc__)
for name in ('source', 'js-proof', 'native-oracle', 'protoc', 'out'):
    p.add_argument('--'+name, required=True)
a = p.parse_args(); source = Path(a.source).resolve(); oracle = Path(a.native_oracle).resolve(); out = Path(a.out).resolve()
sha = lambda path: hashlib.sha256(Path(path).read_bytes()).hexdigest()
cap = Path('/sys/fs/cgroup'); quota, period = (cap/'cpu.max').read_text().split()
if ((cap/'memory.max').read_text().strip() != '4294967296' or (cap/'memory.swap.max').read_text().strip() != '0'
    or quota == 'max' or int(quota) != 2*int(period) or (cap/'cpuset.cpus.effective').read_text().strip() not in ('0-1', '0,1')):
    raise ValueError('actual4GiB/swap0/CPU2 reference check required')
if out.exists() or sha(a.protoc) != '7bf21b9e0a5390dbebe784a4aca1d64683af157b5a22847369737509e8e6cf30':
    raise ValueError('new output and exact pinned native protoc required')
expected = json.loads((oracle/'wire-oracle.json').read_text()); actual = json.loads(Path(a.js_proof).read_text())
if expected['status'] != 'PASS' or actual['status'] != 'PASS' or expected['source_manifest_sha256'] != sha(source/'source-manifest.json'):
    raise ValueError('need actual native and JS payload proof from the same schemas')
manifest = json.loads((source/'source-manifest.json').read_text())
for row in manifest['records']:
    if sha(source/'schemas'/row['path']) != row['sha256']:
        raise ValueError('upstream schema changed')
out.mkdir(parents=True); results = {}
for name in ('table', 'job'):
    values = actual['actualWires'][name]
    if not values or any(type(value) is not int or not 0 <= value <= 255 for value in values):
        raise ValueError('invalid actual JS binary payload')
    raw = bytes(values); (out/(name+'.js.pb')).write_bytes(raw)
    command = [a.protoc, '--proto_path='+str(source/'schemas'), '--decode=google.cloud.bigquery.v2.'+name.title(), 'google/cloud/bigquery/v2/'+name+'.proto']
    decoded = subprocess.run(command, input=raw, capture_output=True)
    (out/(name+'.decoded.textproto')).write_bytes(decoded.stdout); (out/(name+'.stderr')).write_bytes(decoded.stderr)
    reference = oracle/(name+'.decoded.textproto')
    if decoded.returncode or sha(reference) != expected['payloads'][name]['decoded_sha256'] or decoded.stdout != reference.read_bytes():
        raise RuntimeError('actual JS wire does not have the same complete native-decoded API payload: '+name)
    results[name] = {'status': 'PASS', 'argv': command, 'actual_js_wire_sha256': hashlib.sha256(raw).hexdigest(), 'native_decoded_sha256': sha(reference)}
(out/'result.json').write_text(json.dumps({'status': 'PASS', 'scope': 'Independent native protobuf semantics; no compiler/performance claim', 'js_proof_sha256': sha(a.js_proof), 'schemas_sha256': sha(source/'source-manifest.json'), 'payloads': results}, indent=2)+'\n')
print(json.dumps({'status': 'PASS', 'complete_native_payloads_match': list(results)}))
