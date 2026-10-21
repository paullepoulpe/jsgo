#!/usr/bin/env python3
"""Normalize frozen source diffs and verify them in tiny temporary trees.

This performs source-only checks. It never opens jars or invokes Java.
Run once in a fresh output directory containing this script and README.txt.
"""

import difflib
import hashlib
import json
from pathlib import Path
import subprocess
import tempfile


HERE = Path(__file__).resolve().parent
BASE = HERE.parent
PINNED_SOURCE = Path('/workspace/jsgo/.cache/closure-src')
CANDIDATES = [
    ('accumulator', 'scope-maps/accumulator-only/src', 'scope-maps/control-src'),
    ('identity-differs', 'identity-differs/src', 'identity-differs/control-src'),
    ('traversal-noop', 'traversal-noop/src', 'traversal-noop/control-src'),
    ('traversal-name-switch', 'traversal-name-switch/src', 'traversal-name-switch/control-src'),
    ('leaf-children-return', 'leaf-children-return/src', 'leaf-children-return/control-src'),
]


def sha(data):
    return hashlib.sha256(data).hexdigest()


def read_json(path):
    return json.loads(path.read_bytes())


def command(argv, cwd):
    result = subprocess.run(argv, cwd=cwd, text=True, capture_output=True, check=False)
    receipt = {
        'argv': argv,
        'cwd': str(cwd),
        'exit_status': result.returncode,
        'stdout': result.stdout,
        'stderr': result.stderr,
    }
    if result.returncode:
        raise RuntimeError(json.dumps(receipt, indent=2))
    return receipt


def main():
    destinations = [HERE / f'{name}.patch' for name, _, _ in CANDIDATES]
    destinations.append(HERE / 'proof.json')
    if any(path.exists() for path in destinations):
        raise FileExistsError('Use a fresh directory; existing patches/proof are preserved.')

    proof = {
        'status': 'RUNNING',
        'promotion_status': 'NONE_PROMOTED',
        'scope': 'Source normalization and exact patch application only; no new runtime validation.',
        'prepare_source_sha256': sha(Path(__file__).read_bytes()),
        'readme_sha256': sha((HERE / 'README.txt').read_bytes()),
        'candidates': [],
    }
    for name, candidate_dir, control_dir in CANDIDATES:
        manifest_path = BASE / 'build' / name / 'manifest.json'
        manifest = read_json(manifest_path)
        assert manifest['status'] == 'PASS_BUILD_VALIDATION', name
        candidate_receipt_path = Path(manifest['variants'][name]['build_proof'])
        control_receipt_path = Path(manifest['variants']['CONTROL']['build_proof'])
        candidate_receipt_bytes = candidate_receipt_path.read_bytes()
        control_receipt_bytes = control_receipt_path.read_bytes()
        assert sha(candidate_receipt_bytes) == manifest['variants'][name]['build_proof_sha256']
        assert sha(control_receipt_bytes) == manifest['variants']['CONTROL']['build_proof_sha256']
        candidate_receipt = json.loads(candidate_receipt_bytes)
        control_receipt = json.loads(control_receipt_bytes)
        assert candidate_receipt['status'] == control_receipt['status'] == 'PASS'
        assert candidate_receipt['exit'] == control_receipt['exit'] == 0
        families = manifest['compiler_source_families']
        assert len(families) == 1, (name, families)
        relative = families[0] + '.java'
        assert set(candidate_receipt['effective_sources']) == {relative}
        assert set(control_receipt['effective_sources']) == {relative}
        original_path = BASE / 'candidates' / control_dir / relative
        candidate_path = BASE / 'candidates' / candidate_dir / relative
        original = original_path.read_bytes()
        candidate = candidate_path.read_bytes()
        assert original == (PINNED_SOURCE / relative).read_bytes(), name
        assert sha(original) == candidate_receipt['original_sources'][relative]
        assert sha(original) == control_receipt['original_sources'][relative]
        assert sha(original) == control_receipt['effective_sources'][relative]
        assert sha(candidate) == candidate_receipt['effective_sources'][relative]
        assert original != candidate, name
        assert original.endswith(b'\n') and candidate.endswith(b'\n')
        repository_path = 'src/' + relative
        old_label, new_label = 'a/' + repository_path, 'b/' + repository_path
        patch = (
            f'diff --git {old_label} {new_label}\n'
            + ''.join(difflib.unified_diff(
                original.decode('utf-8').splitlines(keepends=True),
                candidate.decode('utf-8').splitlines(keepends=True),
                fromfile=old_label,
                tofile=new_label,
                n=3,
            ))
        ).encode('utf-8')
        patch_path = HERE / f'{name}.patch'
        patch_path.write_bytes(patch)

        with tempfile.TemporaryDirectory(prefix=f'closure-patch-{name}-') as temporary:
            fixture = Path(temporary)
            fixture_source = fixture / repository_path
            fixture_source.parent.mkdir(parents=True)
            fixture_source.write_bytes(original)
            check = command(['git', 'apply', '--check', '--whitespace=error', str(patch_path)], fixture)
            assert fixture_source.read_bytes() == original, name
            apply = command(['git', 'apply', '--whitespace=error', str(patch_path)], fixture)
            applied = fixture_source.read_bytes()
            assert applied == candidate, name
            fixture_files = sorted(str(path.relative_to(fixture)) for path in fixture.rglob('*') if path.is_file())
            assert fixture_files == [repository_path], fixture_files
            reverse_check = command(['git', 'apply', '--reverse', '--check', '--whitespace=error', str(patch_path)], fixture)
        proof['candidates'].append({
            'name': name,
            'status': 'PASS_EXACT_SOURCE_PATCH',
            'repository_path': repository_path,
            'original_source': str(original_path),
            'original_sha256': sha(original),
            'candidate_source': str(candidate_path),
            'candidate_sha256': sha(candidate),
            'applied_sha256': sha(applied),
            'applied_bytes_equal_candidate': True,
            'original_bytes_equal_pinned_source': True,
            'source_hashes_match_existing_build_receipts': True,
            'patch': patch_path.name,
            'patch_sha256': sha(patch),
            'build_manifest': str(manifest_path),
            'build_manifest_sha256': sha(manifest_path.read_bytes()),
            'candidate_build_proof': str(candidate_receipt_path),
            'candidate_build_proof_sha256': sha(candidate_receipt_bytes),
            'control_build_proof': str(control_receipt_path),
            'control_build_proof_sha256': sha(control_receipt_bytes),
            'check': check,
            'apply': apply,
            'reverse_check': reverse_check,
            'fixture_files': fixture_files,
            'temporary_fixture_removed': not fixture.exists(),
        })
        print(f'{name}: PASS check, apply, reverse-check, and exact candidate bytes')

    # Recheck inputs after generation; do not silently certify changed snapshots.
    for entry in proof['candidates']:
        for path_key, hash_key in [('original_source', 'original_sha256'), ('candidate_source', 'candidate_sha256')]:
            assert sha(Path(entry[path_key]).read_bytes()) == entry[hash_key]
    proof['status'] = 'PASS_ALL_FIVE_EXACT_SOURCE_PATCHES'
    (HERE / 'proof.json').write_text(json.dumps(proof, indent=2) + '\n')


if __name__ == '__main__':
    main()
