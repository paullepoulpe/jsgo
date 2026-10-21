#!/usr/bin/env python3
"""Narrow source/probe receipt review and exact upstream patch; no Java/archive use."""
from datetime import datetime, timezone
import difflib
import hashlib
import json
from pathlib import Path
import subprocess
import tempfile

HERE = Path(__file__).resolve().parent
BASE = HERE.parent
NAME = 'reference-list-two'
REL = 'com/google/javascript/jscomp/ReferenceCollection.java'
FAMILY = REL[:-5]
INPUTS = {}


def require(condition, detail):
    if not condition:
        raise ValueError(detail)


def digest(value):
    return hashlib.sha256(value).hexdigest()


def read(path):
    path = Path(path)
    require(path.suffix not in ('.jar', '.zip', '.jfr'), str(path))
    value = path.read_bytes()
    require(len(value) < 2_000_000, 'Not a small receipt/source/class: '+str(path))
    observed = digest(value)
    require(INPUTS.setdefault(str(path), observed) == observed, 'Changed input '+str(path))
    return value


def obj(path):
    return json.loads(read(path))


def main():
    review_dir = BASE/'reviews'/f'{NAME}-probes-attempt-1'
    require(not review_dir.exists(), 'Preserve prior review')
    require(not (HERE/'proof.json').exists(), 'Preserve prior proof')
    candidate_root = BASE/'candidates'/NAME
    frozen_path = candidate_root/'frozen-source-manifest.json'
    frozen = obj(frozen_path)
    for relative, expected in frozen['files'].items():
        actual = read(candidate_root/relative)
        require(digest(actual) == expected['sha256'] and len(actual) == expected['bytes'], relative)
    original = read(candidate_root/'control-src'/REL)
    candidate = read(candidate_root/'src'/REL)
    require(original == read(Path('/workspace/jsgo/.cache/closure-src')/REL), 'Pinned source')
    require(candidate == original.replace(b'new ArrayList<>()', b'new ArrayList<>(2)', 1), 'One initializer')
    build_path = BASE/'build'/NAME/'manifest.json'
    build = obj(build_path)
    require(build['status'] == 'PASS_BUILD_VALIDATION' and build['compiler_source_families'] == [FAMILY], 'Build family')
    receipts = {}
    for variant, source in [('CONTROL', original), (NAME, candidate)]:
        entry = build['variants'][variant]
        receipt_path = Path(entry['build_proof'])
        receipt = obj(receipt_path)
        require(digest(read(receipt_path)) == entry['build_proof_sha256'], 'Build receipt hash')
        require(receipt['status'] == 'PASS' and receipt['exit'] == 0, 'Build success')
        require(receipt['original_sources'] == {REL:digest(original)}, 'Original source')
        require(receipt['effective_sources'] == {REL:digest(source)}, 'Effective source')
        relocated = read(BASE/'build'/NAME/variant/'src'/REL)
        expected_relocated = source.replace(b'import static com.google.common.', b'import static com.google.javascript.jscomp.jarjar.com.google.common.')
        require(relocated == expected_relocated and receipt['relocated_sources'] == {REL:digest(relocated)}, 'Relocation')
        overlay = Path(entry['overlay'])
        require({str(p.relative_to(overlay)):digest(read(p)) for p in overlay.rglob('*.class')} == entry['classes'] == receipt['classes'], 'Overlay bytes')
        receipts[variant] = {'path':str(receipt_path), 'sha256':digest(read(receipt_path)), 'effective_source_sha256':digest(source)}
    probes_path = BASE/'results'/f'{NAME}-probes-attempt-1'/'manifest.json'
    probes = obj(probes_path)
    require(probes['status'] == 'PASS' and probes['execution_complete'] and probes['expected_runs'] == 3, 'Probe completeness')
    require(probes['build_manifest_unchanged'] and probes['build_manifest_sha256'] == digest(read(build_path)), 'Probe build linkage')
    require(all(probes['source_integrity'].values()) and all(probes['overlay_integrity'].values()), 'Probe integrity')
    require(probes['compile']['status'] == 'PASS' and probes['compile']['exit'] == 0, 'Probe compile')
    require([run['variant'] for run in probes['runs']] == ['STOCK', 'CONTROL', NAME], 'Three variants')
    outputs = []
    for run in probes['runs']:
        require(run['status'] == 'PASS' and run['exit'] == 0, 'Probe run')
        stdout, stderr = read(Path(run['cwd'])/'stdout'), read(Path(run['cwd'])/'stderr')
        require(digest(stdout) == run['stdout_sha256'] and digest(stderr) == run['stderr_sha256'] and not stderr, 'Probe output bytes')
        require(b'ReferenceListTwoProbe OK serialUID=' in stdout, 'Probe completion marker')
        origin = run['class_origins'][FAMILY.replace('/', '.')]
        expected = probes['jar'] if run['variant'] == 'STOCK' else build['variants'][run['variant']]['overlay']
        require(origin['pass'] and origin['normalized_origins'] == [expected] and origin['expected_path'] == expected, 'Recorded class origin')
        log = read(Path(run['cwd'])/'class-load.log')
        require(digest(log) == run['class_load_log_sha256'], 'Class log bytes')
        require(any((FAMILY.replace('/', '.')+' source: ').encode() in line and expected.encode() in line for line in log.splitlines()), 'Actual class load line')
        outputs.append(stdout)
    require(outputs[0] == outputs[1] == outputs[2], 'Exact stdout parity')
    probe_source = candidate_root/'probe/com/google/javascript/jscomp/ReferenceListTwoProbe.java'
    require(probes['sources']['com/google/javascript/jscomp/ReferenceListTwoProbe.java']['original_sha256'] == digest(read(probe_source)), 'Frozen probe link')
    repository_path = 'src/'+REL
    old, new = 'a/'+repository_path, 'b/'+repository_path
    patch = ('diff --git '+old+' '+new+'\n'+''.join(difflib.unified_diff(original.decode().splitlines(keepends=True), candidate.decode().splitlines(keepends=True), fromfile=old, tofile=new))).encode()
    patch_path = HERE/(NAME+'.patch')
    patch_path.write_bytes(patch)
    steps = []
    with tempfile.TemporaryDirectory(prefix='closure-patch-v4-') as temporary:
        fixture = Path(temporary)
        target = fixture/repository_path
        target.parent.mkdir(parents=True)
        target.write_bytes(original)
        for args, expected in [(['--check'], original), ([], candidate), (['--reverse', '--check'], candidate), (['--reverse'], original)]:
            command = ['git', 'apply', '--whitespace=error', *args, str(patch_path)]
            done = subprocess.run(command, cwd=fixture, capture_output=True, text=True, check=False)
            require(done.returncode == 0 and target.read_bytes() == expected, 'Exact patch application '+str(args))
            require([str(p.relative_to(fixture)) for p in fixture.rglob('*') if p.is_file()] == [repository_path], 'Unexpected patch file')
            steps.append({'argv':command, 'exit':done.returncode, 'stdout':done.stdout, 'stderr':done.stderr, 'exact_expected_bytes':True})
    summary_path = BASE/'results'/f'screen-bigquery-control-{NAME}-v1'/'summary.json'
    summary = obj(summary_path)
    require(summary['complete'] and summary['candidate'] == NAME and summary['baseline'] == 'CONTROL', 'Completed screen summary')
    require(summary['source_manifest_sha256'] == digest(read(build_path)), 'Screen build link')
    review = {'status':'PASS_SOURCE_BUILD_AND_PROBE_RECEIPT_REVIEW', 'promotion_status':'NONE_PROMOTED',
        'reviewed_at_utc':datetime.now(timezone.utc).isoformat(), 'scope':'Existing small receipts/source/classes only; no new Java/build/benchmark/archive inspection.',
        'frozen_manifest':str(frozen_path), 'build_manifest':str(build_path), 'build_manifest_sha256':digest(read(build_path)),
        'probe_manifest':str(probes_path), 'probe_manifest_sha256':digest(read(probes_path)), 'runs':3,
        'all_stdout_exactly_equal':True, 'stdout_sha256':digest(outputs[0]), 'all_stderr_empty':True,
        'class_origins_and_log_hashes_verified':True, 'full_overlay_classes_verified':True,
        'original_source_sha256':digest(original), 'candidate_source_sha256':digest(candidate), 'existing_build_receipts':receipts,
        'performance_status':'NEGATIVE_SINGLE_COHORT_SCREEN_NONE_PROMOTED', 'performance_source':str(summary_path),
        'performance_summary_sha256':digest(read(summary_path)), 'performance_review_scope':'Summary transcription only; no independent full screen reanalysis.',
        'candidate_over_CONTROL':summary['cohorts'][0]['candidate_over_baseline'],
        'limits':['Focused probes are not the full upstream Java suite.', 'Current stock jar bytes were not bulk rehashed.', 'No combination with other patches validated.', 'Capacity reflection and arbitrary capacity-changing same-package callers are outside logical compiler equivalence.']}
    review_dir.mkdir(parents=True)
    (review_dir/'review.json').write_text(json.dumps(review,indent=2)+'\n')
    proof = {'status':'PASS_EXACT_FORWARD_REVERSE_SOURCE_BUILD_LINKAGE', 'promotion_status':'NONE_PROMOTED',
        'combined_patches_validated':False, 'patch':patch_path.name, 'patch_sha256':digest(patch),
        'original_sha256':digest(original), 'candidate_sha256':digest(candidate), 'steps':steps,
        'temporary_fixture_removed':not fixture.exists(), 'probe_review':str(review_dir/'review.json'),
        'probe_review_sha256':digest((review_dir/'review.json').read_bytes()), 'performance_status':review['performance_status']}
    read(HERE/'prepare.py'); read(HERE/'README.txt')
    for path, expected in list(INPUTS.items()):
        require(digest(read(path)) == expected, 'Input mutated')
    proof['inputs_unchanged'] = True
    proof['inputs_sha256'] = INPUTS
    (HERE/'proof.json').write_text(json.dumps(proof,indent=2)+'\n')
    print('PASS: frozen/build/probe receipt checks; exact forward/reverse upstream patch. No runtime started.')


if __name__ == '__main__':
    main()
