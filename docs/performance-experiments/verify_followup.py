#!/usr/bin/env python3
"""Verify original and incremental records without executing archived prototypes."""
import hashlib
import json
from pathlib import Path
import subprocess
import sys
import tarfile

root = Path(__file__).resolve().parent
subprocess.run([sys.executable, str(root / 'verify_archive.py')], check=True)
sha = lambda b: hashlib.sha256(b).hexdigest()
manifest = json.loads((root / 'followup-measurement-manifest.json').read_text())
for key in ('base_manifest', 'base_archive'):
    entry = manifest[key]
    assert sha((root / entry['path']).read_bytes()) == entry['sha256'], key
expected = {'records/' + e['path']: e for e in manifest['files']}
with tarfile.open(root / manifest['archive']['path'], 'r:gz') as archive:
    seen = set()
    for item in archive:
        assert item.isfile() and not Path(item.name).is_absolute()
        assert '..' not in Path(item.name).parts and item.name not in seen
        data = archive.extractfile(item).read()
        entry = expected[item.name]
        assert len(data) == entry['bytes'] and sha(data) == entry['sha256'], item.name
        browse = root / 'followup-records' / entry['path']
        if browse.exists():
            assert browse.read_bytes() == data, browse
        seen.add(item.name)
    assert seen == set(expected)
source = json.loads((root / 'followup-source-manifest.json').read_text())
bundle = source['bundle']
data = (root / bundle['path']).read_bytes()
assert len(data) == bundle['bytes'] and sha(data) == bundle['sha256']
print(f'PASS: {len(expected)} follow-up records and source bundle hash')
