#!/usr/bin/env python3
"""Normalize five frozen candidates and verify exact forward/reverse application.

Source-only: no Java, build, benchmark, network, archive, or class-file access.
The output directory must be fresh apart from this script and README.txt.
"""

from datetime import datetime, timezone
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
    ('traversal-leaf-name', 'com/google/javascript/jscomp/NodeTraversal.java',
     'MEASURED_SCREEN_NEGATIVE_NO_PROMOTION'),
    ('reference-function-first-child', 'com/google/javascript/jscomp/Reference.java',
     'NOT_MEASURED'),
    ('graph-edge-copy', 'com/google/javascript/jscomp/graph/LinkedDirectedGraph.java',
     'NOT_MEASURED'),
    ('source-position-last', 'com/google/javascript/jscomp/parsing/parser/LineNumberScanner.java',
     'NOT_MEASURED'),
    ('equality-prototype-scalar', 'com/google/javascript/rhino/jstype/JSType.java',
     'TIMING_PENDING'),
]


def sha(data):
    return hashlib.sha256(data).hexdigest()


def command(argv, cwd):
    result = subprocess.run(argv, cwd=cwd, text=True, capture_output=True, check=False)
    receipt = {
        'argv': argv, 'cwd': str(cwd), 'exit_status': result.returncode,
        'stdout': result.stdout, 'stderr': result.stderr,
    }
    if result.returncode:
        raise RuntimeError(json.dumps(receipt, indent=2))
    return receipt


def frozen_hash(manifest, relative):
    if 'files_sha256' in manifest:
        return manifest['files_sha256'][relative]
    return manifest['files'][relative]['sha256']


def main():
    destinations = [HERE / f'{name}.patch' for name, _, _ in CANDIDATES]
    destinations.append(HERE / 'proof.json')
    if any(path.exists() for path in destinations):
        raise FileExistsError('Use a fresh output directory; existing patches/proof are preserved.')

    observed_inputs = {}

    def read_input(path):
        data = path.read_bytes()
        digest = sha(data)
        previous = observed_inputs.setdefault(str(path), digest)
        assert previous == digest, ('input changed during preparation', str(path))
        return data

    proof = {
        'status': 'RUNNING',
        'pinned_version': 'v20261015',
        'prepared_at_utc': datetime.now(timezone.utc).isoformat(),
        'promotion_status': 'NONE_PROMOTED',
        'combined_patches_validated': False,
        'scope': 'Exact individual source patches and existing build-source receipts only; '
                 'no new build/runtime/performance validation.',
        'status_snapshot_authority': 'Parent task instructions at preparation; not inferred from build success.',
        'prepare_source_sha256': sha(read_input(Path(__file__))),
        'readme_sha256': sha(read_input(HERE / 'README.txt')),
        'candidates': [],
    }

    for name, relative, performance_status in CANDIDATES:
        candidate_root = BASE / 'candidates' / name
        frozen_path = candidate_root / 'frozen-source-manifest.json'
        frozen_bytes = read_input(frozen_path)
        frozen = json.loads(frozen_bytes)
        original_path = candidate_root / 'control-src' / relative
        candidate_path = candidate_root / 'src' / relative
        original = read_input(original_path)
        candidate = read_input(candidate_path)
        original_sha, candidate_sha = sha(original), sha(candidate)
        assert original_sha == frozen_hash(frozen, 'control-src/' + relative), name
        assert candidate_sha == frozen_hash(frozen, 'src/' + relative), name
        pinned_path = PINNED_SOURCE / relative
        assert original == read_input(pinned_path), name
        assert original != candidate, name
        assert original.endswith(b'\n') and candidate.endswith(b'\n'), name

        manifest_path = BASE / 'build' / name / 'manifest.json'
        manifest_bytes = read_input(manifest_path)
        manifest = json.loads(manifest_bytes)
        assert manifest['status'] == 'PASS_BUILD_VALIDATION', name
        assert manifest['compiler_source_families'] == [relative.removesuffix('.java')], name
        build_proofs = {}
        for variant in ['CONTROL', name]:
            metadata = manifest['variants'][variant]
            receipt_path = Path(metadata['build_proof'])
            receipt_bytes = read_input(receipt_path)
            receipt_sha = sha(receipt_bytes)
            assert receipt_sha == metadata['build_proof_sha256'], (name, variant)
            receipt = json.loads(receipt_bytes)
            assert receipt['status'] == 'PASS' and receipt['exit'] == 0, (name, variant)
            assert receipt['original_sources'] == {relative: original_sha}, (name, variant)
            expected_effective = original_sha if variant == 'CONTROL' else candidate_sha
            assert receipt['effective_sources'] == {relative: expected_effective}, (name, variant)
            build_proofs[variant] = {
                'path': str(receipt_path), 'sha256': receipt_sha,
                'status': receipt['status'], 'exit_status': receipt['exit'],
                'original_source_sha256': original_sha,
                'effective_unshaded_source_sha256': expected_effective,
            }

        repository_path = 'src/' + relative
        old_label, new_label = 'a/' + repository_path, 'b/' + repository_path
        patch = (f'diff --git {old_label} {new_label}\n' + ''.join(difflib.unified_diff(
            original.decode('utf-8').splitlines(keepends=True),
            candidate.decode('utf-8').splitlines(keepends=True),
            fromfile=old_label, tofile=new_label, n=3,
        ))).encode('utf-8')
        patch_path = HERE / f'{name}.patch'
        patch_path.write_bytes(patch)

        with tempfile.TemporaryDirectory(prefix=f'closure-patch-v2-{name}-') as temporary:
            fixture = Path(temporary)
            fixture_source = fixture / repository_path
            fixture_source.parent.mkdir(parents=True)
            fixture_source.write_bytes(original)
            check = command(['git', 'apply', '--check', '--whitespace=error', str(patch_path)], fixture)
            assert fixture_source.read_bytes() == original, name
            apply = command(['git', 'apply', '--whitespace=error', str(patch_path)], fixture)
            applied = fixture_source.read_bytes()
            assert applied == candidate, name
            forward_files = sorted(str(p.relative_to(fixture)) for p in fixture.rglob('*') if p.is_file())
            assert forward_files == [repository_path], (name, forward_files)
            reverse_check = command(
                ['git', 'apply', '--reverse', '--check', '--whitespace=error', str(patch_path)], fixture)
            assert fixture_source.read_bytes() == candidate, name
            reverse_apply = command(
                ['git', 'apply', '--reverse', '--whitespace=error', str(patch_path)], fixture)
            restored = fixture_source.read_bytes()
            assert restored == original, name
            reverse_files = sorted(str(p.relative_to(fixture)) for p in fixture.rglob('*') if p.is_file())
            assert reverse_files == [repository_path], (name, reverse_files)

        entry = {
            'name': name,
            'status': 'PASS_EXACT_SOURCE_PATCH_AND_EXISTING_BUILD_SOURCE_RECEIPTS',
            'promotion_status': 'NOT_PROMOTED',
            'performance_status': performance_status,
            'repository_path': repository_path,
            'frozen_source_manifest': str(frozen_path),
            'frozen_source_manifest_sha256': sha(frozen_bytes),
            'original_source': str(original_path), 'original_sha256': original_sha,
            'candidate_source': str(candidate_path), 'candidate_sha256': candidate_sha,
            'source_hashes_match_frozen_manifest': True,
            'pinned_source': str(pinned_path), 'original_bytes_equal_pinned_source': True,
            'applied_sha256': sha(applied), 'applied_bytes_equal_frozen_candidate': True,
            'restored_sha256': sha(restored), 'restored_bytes_equal_original': True,
            'patch': patch_path.name, 'patch_sha256': sha(patch),
            'build_manifest': str(manifest_path), 'build_manifest_sha256': sha(manifest_bytes),
            'source_hashes_match_existing_build_receipts': True,
            'existing_build_proofs': build_proofs,
            'check': check, 'apply': apply,
            'reverse_check': reverse_check, 'reverse_apply': reverse_apply,
            'fixture_files_after_apply': forward_files,
            'fixture_files_after_reverse': reverse_files,
            'temporary_fixture_removed': not fixture.exists(),
        }
        if name == 'traversal-leaf-name':
            summary_path = BASE / 'results/screen-bigquery-control-traversal-leaf-name-v1/summary.json'
            summary_bytes = read_input(summary_path)
            summary = json.loads(summary_bytes)
            assert summary['complete'] is True and summary['candidate'] == name
            entry['existing_timing_evidence'] = {
                'summary': str(summary_path), 'sha256': sha(summary_bytes),
                'status': summary['status'], 'baseline': summary['baseline'],
                'cohort_count': len(summary['cohorts']),
                'candidate_over_baseline_by_cohort': [
                    c['candidate_over_baseline'] for c in summary['cohorts']],
                'interpretation': 'Measured screen negative for promotion; no useful gain established.',
            }
        proof['candidates'].append(entry)
        print(f'{name}: PASS forward/reverse exact bytes, frozen/pinned and build-source receipts; '
              f'{performance_status}')

    # Recheck only the explicit small files consumed above, not directories/archives/classes.
    for path, digest in observed_inputs.items():
        assert sha(Path(path).read_bytes()) == digest, ('input changed', path)
    for entry in proof['candidates']:
        assert sha((HERE / entry['patch']).read_bytes()) == entry['patch_sha256']
    proof['inputs_unchanged_after_generation'] = True
    proof['status'] = 'PASS_ALL_FIVE_EXACT_SOURCE_PATCHES'
    (HERE / 'proof.json').write_text(json.dumps(proof, indent=2) + '\n')


if __name__ == '__main__':
    main()
