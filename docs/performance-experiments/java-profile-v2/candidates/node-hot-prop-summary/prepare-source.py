#!/usr/bin/env python3
"""Source-only exact patch, head-store audit, and forward/reverse byte proof."""
from datetime import datetime, timezone
import difflib
import hashlib
import json
from pathlib import Path
import re
import subprocess
import tempfile

HERE = Path(__file__).resolve().parent
REL = 'com/google/javascript/rhino/Node.java'
PINNED = Path('/workspace/jsgo/.cache/closure-src') / REL


def sha(data):
    return hashlib.sha256(data).hexdigest()


def require(condition, message):
    if not condition:
        raise ValueError(message)


def stores(source):
    pattern = r'(?<![\w.])(?:(?:this|dest)\.)?propListHead\s*=(?!=)[^;]*;'
    return [{'line': source.count('\n', 0, m.start()) + 1,
             'statement': m.group(0)} for m in re.finditer(pattern, source)]


def run(argv, cwd):
    result = subprocess.run(argv, cwd=cwd, capture_output=True, text=True, check=False)
    receipt = {'argv': argv, 'cwd': str(cwd), 'exit_status': result.returncode,
               'stdout': result.stdout, 'stderr': result.stderr}
    require(result.returncode == 0, json.dumps(receipt))
    return receipt


def main():
    outputs = [HERE/'node-hot-prop-summary.patch', HERE/'patch-check.json',
               HERE/'profile-evidence.tsv']
    require(not any(p.exists() for p in outputs), 'Preserve prior outputs; use a new version')
    control = (HERE/'control-src'/REL).read_bytes()
    candidate = (HERE/'src'/REL).read_bytes()
    require(control == PINNED.read_bytes(), 'Control differs from pinned source')
    require(control != candidate, 'Empty candidate')
    old, new = control.decode(), candidate.decode()
    old_stores, new_stores = stores(old), stores(new)
    require(len(old_stores) == 12, 'Unexpected number of original head stores')
    require([x['statement'] for x in old_stores] == [x['statement'] for x in new_stores],
            'Head stores changed')
    enum = old.split('  enum Prop {', 1)[1].split('\n  }', 1)[0]
    enum = re.sub(r'//[^\n]*', '', enum)
    names = [x.strip() for x in enum.split(',') if x.strip()]
    require(all(re.fullmatch(r'[A-Z_0-9]+', x) for x in names), 'Unexpected enum syntax')
    require(len(names) < 128, 'Prop enum no longer fits positive byte')
    old_label, new_label = 'a/src/'+REL, 'b/src/'+REL
    patch = ('diff --git '+old_label+' '+new_label+'\n'+''.join(difflib.unified_diff(
        old.splitlines(True), new.splitlines(True), fromfile=old_label, tofile=new_label,
        n=3))).encode()
    patch_path = outputs[0]
    patch_path.write_bytes(patch)
    with tempfile.TemporaryDirectory(prefix='node-hot-prop-summary-patch-') as temporary:
        fixture = Path(temporary)
        target = fixture/'src'/REL
        target.parent.mkdir(parents=True)
        target.write_bytes(control)
        check = run(['git', 'apply', '--check', '--whitespace=error', str(patch_path)], fixture)
        require(target.read_bytes() == control, 'Check changed source')
        apply = run(['git', 'apply', '--whitespace=error', str(patch_path)], fixture)
        require(target.read_bytes() == candidate, 'Applied bytes differ')
        fixture_files = sorted(str(p.relative_to(fixture)) for p in fixture.rglob('*') if p.is_file())
        require(fixture_files == ['src/'+REL], 'Unexpected fixture file')
        reverse_check = run(['git', 'apply', '--reverse', '--check', '--whitespace=error',
                             str(patch_path)], fixture)
        require(target.read_bytes() == candidate, 'Reverse check changed source')
        reverse = run(['git', 'apply', '--reverse', '--whitespace=error', str(patch_path)], fixture)
        require(target.read_bytes() == control, 'Reversed bytes differ')
    profile = HERE.parents[1]/'results/profile-stock-bigquery-steady-v2/hot-callers.tsv'
    rows = [line for line in profile.read_text().splitlines()
            if line.startswith('CPU_CALLERS4\tcom.google.javascript.rhino.Node.lookupProperty\t')
            and ('Node.isArrowFunction' in line or 'Node.isSyntheticBlock' in line)]
    counts = {}
    for method, expected in [('isArrowFunction', 26), ('isSyntheticBlock', 21)]:
        count = sum(int(row.split('\t')[3]) for row in rows if 'Node.'+method in row)
        require(count == expected, 'Unexpected profile evidence')
        counts[method] = {'samples': count, 'all_cpu_samples': 2221, 'pct_all_cpu': count/2221*100}
    require(all(row.split('\t')[5] == '2221' for row in rows), 'Profile denominator')
    outputs[2].write_text('SOURCE\t'+str(profile)+'\n'
        'NOTE\tRecorded caller-prefix leaf samples, not call counts or expected speedup.\n'
        'COLUMNS\ttype\tleaf\trank\tsamples\tmatched_leaf_denominator\tall_cpu_denominator'
        '\tpct_leaf\tpct_all_cpu\tcaller_frames\n'+'\n'.join(rows)+'\n')
    require((HERE/'control-src'/REL).read_bytes() == control, 'Control changed')
    require((HERE/'src'/REL).read_bytes() == candidate, 'Candidate changed')
    proof = {'status': 'PASS_EXACT_SOURCE_FORWARD_REVERSE',
             'prepared_at_utc': datetime.now(timezone.utc).isoformat(),
             'pinned_version': 'v20261015', 'compile_status': 'NOT_RUN',
             'runtime_status': 'NOT_RUN', 'performance_status': 'NOT_MEASURED',
             'source_only_check': True, 'pinned_source': str(PINNED),
             'control_sha256': sha(control), 'candidate_sha256': sha(candidate),
             'prepare_script_sha256': sha(Path(__file__).read_bytes()),
             'patch_sha256': sha(patch), 'control_equals_pinned': True,
             'applied_bytes_equal_candidate': True, 'reversed_bytes_equal_control': True,
             'original_head_stores': old_stores, 'candidate_head_stores': new_stores,
             'all_head_stores_unchanged': True, 'prop_count': len(names),
             'hot_prop_ordinals': {name: names.index(name) for name in ['ARROW_FN', 'SYNTHETIC']},
             'profile_observations': counts, 'fixture_files': fixture_files,
             'check': check, 'apply': apply, 'reverse_check': reverse_check, 'reverse': reverse,
             'temporary_fixture_removed': not fixture.exists()}
    outputs[1].write_text(json.dumps(proof, indent=2)+'\n')
    print(json.dumps({key: proof[key] for key in ['status', 'control_sha256', 'candidate_sha256',
                                                'patch_sha256', 'all_head_stores_unchanged']}))


if __name__ == '__main__':
    main()
