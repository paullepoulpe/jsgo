#!/usr/bin/env python3
"""Verify publication archives and per-file manifests without executing experiments."""
import hashlib
import json
from pathlib import Path
import tarfile

root=Path(__file__).resolve().parent
sha=lambda b:hashlib.sha256(b).hexdigest()
for item in json.loads((root/'archives.json').read_text()):
 p=root/item['path'];assert p.stat().st_size==item['bytes'],p
 assert sha(p.read_bytes())==item['sha256'],p
expected={'records/'+x['path']:x for x in json.loads((root/'measurement-manifest.json').read_text())['files']}
with tarfile.open(root/'measurement-records.tar.gz','r:gz') as t:
 seen=set()
 for m in t:
  assert m.isfile() and not Path(m.name).is_absolute() and '..' not in Path(m.name).parts,m.name
  assert m.name not in seen,m.name
  b=t.extractfile(m).read();e=expected[m.name]
  assert len(b)==e['bytes'] and sha(b)==e['sha256'],m.name
  seen.add(m.name)
 assert seen==set(expected)
expected_source={}
for item in json.loads((root/'source-manifest.json').read_text())['worktrees']:
 expected_source[item['tracked_patch']]={'bytes':item['tracked_patch_bytes'],'sha256':item['tracked_patch_sha256']}
 for u in item['untracked']:expected_source[item['worktree']+'/untracked/'+u['path']]=u
with tarfile.open(root/'prototype-sources.tar.gz','r:gz') as t:
 seen=set()
 for m in t:
  assert m.isfile() and not Path(m.name).is_absolute() and '..' not in Path(m.name).parts,m.name
  assert m.name not in seen,m.name
  b=t.extractfile(m).read();e=expected_source[m.name]
  assert len(b)==e['bytes'] and sha(b)==e['sha256'],m.name
  seen.add(m.name)
 assert seen==set(expected_source)
print(f'PASS: archive hashes, {len(expected)} measurement records, {len(expected_source)} source records')
