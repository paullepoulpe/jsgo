#!/usr/bin/env python3
"""Normalize two closed sources, checking individual exact forward/reverse bytes."""
from datetime import datetime, timezone
import difflib
import hashlib
import json
from pathlib import Path
import subprocess
import tempfile

HERE = Path(__file__).resolve().parent
BASE = HERE.parent
PINNED = Path('/workspace/jsgo/.cache/closure-src')
INPUTS = {}
SPECS = [
    ('empty-parameters', 'com/google/javascript/jscomp/FunctionArgumentInjector.java',
     BASE/'candidates/scope-maps/control-src',
     BASE/'candidates/scope-maps/empty-parameters-only/src',
     BASE/'audits/empty-parameters-callsite-v1/source-review.json'),
    ('scope-declare-once', 'com/google/javascript/jscomp/AbstractScope.java',
     BASE/'candidates/scope-declare-once/control-src',
     BASE/'candidates/scope-declare-once/src',
     BASE/'candidates/scope-declare-once/frozen-source-manifest.json'),
]


def require(condition, detail):
    if not condition:
        raise ValueError(detail)


def digest(value):
    return hashlib.sha256(value).hexdigest()


def read(path):
    path = Path(path)
    require(path.suffix not in ('.class', '.jar', '.zip', '.jfr'), str(path))
    value = path.read_bytes()
    observed = digest(value)
    require(INPUTS.setdefault(str(path), observed) == observed, 'Changed input '+str(path))
    return value


def obj(path):
    return json.loads(read(path))


def run(argv, directory):
    result = subprocess.run(argv, cwd=directory, capture_output=True, text=True, check=False)
    receipt = {'argv':argv, 'cwd':str(directory), 'exit_status':result.returncode,
               'stdout':result.stdout, 'stderr':result.stderr}
    require(result.returncode == 0, json.dumps(receipt))
    return receipt


def main():
    outputs = [HERE/(s[0]+'.patch') for s in SPECS]+[HERE/'proof.json']
    require(not any(path.exists() for path in outputs), 'Preserve prior output; use fresh directory')
    review_path = BASE/'reviews/empty-scope-probes-v1/review.json'
    review = obj(review_path)
    require(review['status'] == 'PASS_ALL_THREE_INDEPENDENT_RECEIPT_REVIEWS', 'Probe review')
    proof = {'status':'RUNNING', 'pinned_version':'v20261015',
             'prepared_at_utc':datetime.now(timezone.utc).isoformat(),
             'promotion_status':'NONE_PROMOTED', 'combined_patches_validated':False,
             'performance_status':'UNMEASURED_AT_TASK_ASSIGNMENT_NO_ACTIVE_SCREEN_READ',
             'scope':'Source-only exact patches and existing build/probe receipt links; no new runtime.',
             'prepare_sha256':digest(read(Path(__file__))),
             'readme_sha256':digest(read(HERE/'README.txt')),
             'independent_probe_review':str(review_path),
             'independent_probe_review_sha256':digest(read(review_path)), 'candidates':[]}
    for name, relative, control_root, candidate_root, closed_path in SPECS:
        control_path, candidate_path = control_root/relative, candidate_root/relative
        original, candidate = read(control_path), read(candidate_path)
        original_sha, candidate_sha = digest(original), digest(candidate)
        require(original == read(PINNED/relative), name+' original differs from pinned source')
        require(original != candidate, name+' no change')
        require(original.endswith(b'\n') and candidate.endswith(b'\n'), name)
        closed = obj(closed_path)
        if name == 'empty-parameters':
            closed_hashes = {entry['path']:entry['sha256'] for entry in closed['files']}
            require(closed_hashes[str(control_path)] == original_sha, name)
            require(closed_hashes[str(candidate_path)] == candidate_sha, name)
        else:
            require(closed['files_sha256']['control-src/'+relative] == original_sha, name)
            require(closed['files_sha256']['src/'+relative] == candidate_sha, name)

        manifest_path = BASE/'build'/name/'manifest.json'
        manifest = obj(manifest_path)
        require(manifest['status'] == 'PASS_BUILD_VALIDATION', name)
        require(manifest['compiler_source_families'] == [relative[:-5]], name)
        build_proofs = {}
        for variant, effective_sha in [('CONTROL', original_sha), (name, candidate_sha)]:
            metadata = manifest['variants'][variant]
            build_path = Path(metadata['build_proof'])
            build_bytes = read(build_path)
            require(digest(build_bytes) == metadata['build_proof_sha256'], name+' '+variant)
            receipt = json.loads(build_bytes)
            require(receipt['status'] == 'PASS' and receipt['exit'] == 0, name+' '+variant)
            require(receipt['original_sources'] == {relative:original_sha}, name+' '+variant)
            require(receipt['effective_sources'] == {relative:effective_sha}, name+' '+variant)
            build_proofs[variant] = {'path':str(build_path), 'sha256':digest(build_bytes),
                'original_source_sha256':original_sha, 'effective_source_sha256':effective_sha,
                'status':receipt['status'], 'exit_status':receipt['exit']}

        repository_path = 'src/'+relative
        old_label, new_label = 'a/'+repository_path, 'b/'+repository_path
        patch = ('diff --git '+old_label+' '+new_label+'\n'+''.join(difflib.unified_diff(
            original.decode().splitlines(keepends=True), candidate.decode().splitlines(keepends=True),
            fromfile=old_label, tofile=new_label, n=3))).encode()
        patch_path = HERE/(name+'.patch')
        patch_path.write_bytes(patch)
        with tempfile.TemporaryDirectory(prefix='closure-patch-v3-'+name+'-') as temporary:
            fixture = Path(temporary)
            target = fixture/repository_path
            target.parent.mkdir(parents=True)
            target.write_bytes(original)
            check = run(['git', 'apply', '--check', '--whitespace=error', str(patch_path)], fixture)
            require(target.read_bytes() == original, name)
            apply = run(['git', 'apply', '--whitespace=error', str(patch_path)], fixture)
            applied = target.read_bytes()
            require(applied == candidate, name+' forward bytes')
            forward_files = sorted(str(p.relative_to(fixture)) for p in fixture.rglob('*') if p.is_file())
            require(forward_files == [repository_path], name+' unexpected fixture file')
            reverse_check = run(['git', 'apply', '--reverse', '--check', '--whitespace=error', str(patch_path)], fixture)
            require(target.read_bytes() == candidate, name)
            reverse_apply = run(['git', 'apply', '--reverse', '--whitespace=error', str(patch_path)], fixture)
            restored = target.read_bytes()
            require(restored == original, name+' reverse bytes')
            reverse_files = sorted(str(p.relative_to(fixture)) for p in fixture.rglob('*') if p.is_file())
            require(reverse_files == [repository_path], name)
        proof['candidates'].append({'name':name, 'status':'PASS_EXACT_FORWARD_AND_REVERSE',
            'performance_status':'UNMEASURED_AT_TASK_ASSIGNMENT', 'promotion_status':'NOT_PROMOTED',
            'repository_path':repository_path, 'patch':patch_path.name, 'patch_sha256':digest(patch),
            'original_source':str(control_path), 'original_sha256':original_sha,
            'candidate_source':str(candidate_path), 'candidate_sha256':candidate_sha,
            'closed_source_record':str(closed_path), 'closed_source_record_sha256':digest(read(closed_path)),
            'original_bytes_equal_pinned':True, 'applied_bytes_equal_candidate':True,
            'restored_bytes_equal_original':True, 'applied_sha256':digest(applied),
            'restored_sha256':digest(restored), 'build_manifest':str(manifest_path),
            'build_manifest_sha256':digest(read(manifest_path)), 'existing_build_proofs':build_proofs,
            'check':check, 'apply':apply, 'reverse_check':reverse_check, 'reverse_apply':reverse_apply,
            'fixture_files_after_apply':forward_files, 'fixture_files_after_reverse':reverse_files,
            'temporary_fixture_removed':not fixture.exists()})
        print(name+': PASS exact forward/reverse bytes, pinned/closed/build source hashes')
    for path, expected in list(INPUTS.items()):
        require(digest(read(path)) == expected, 'Changed input '+path)
    for entry in proof['candidates']:
        require(digest((HERE/entry['patch']).read_bytes()) == entry['patch_sha256'], entry['name'])
    proof['inputs_unchanged_after_generation'] = True
    proof['status'] = 'PASS_BOTH_EXACT_SOURCE_PATCHES'
    (HERE/'proof.json').write_text(json.dumps(proof,indent=2)+'\n')


if __name__ == '__main__':
    main()
