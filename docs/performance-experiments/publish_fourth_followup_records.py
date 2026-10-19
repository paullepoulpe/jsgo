#!/usr/bin/env python3
"""Publish a fourth new/changed completed-evidence delta; never rewrite the earlier archives.

The delta tar retains records/<original-relative-path>. Browse copies have the
separate fourth-followup-records/ prefix. This script does not snapshot source worktrees,
edit README/archives.json, commit, push, or execute any archived experiment.
"""
import argparse
import collections
import csv
import gzip
import hashlib
import io
import json
import math
from pathlib import Path
import re
import statistics
import tarfile


KEEP_EXT = {'.json', '.md', '.log', '.txt', '.tsv', '.csv', '.py', '.sh', '.go',
            '.java', '.pprof', '.pgo', '.stdout', '.stderr'}
KEEP_NAMES = {'stdout', 'stderr', 'exit', 'command', 'metadata', 'status', 'time', 'timing', 'rss'}
# Active paths are explicitly supplied by the coordinator after timing ends.

SECRET = re.compile(rb'(?:gh[pousr]_[A-Za-z0-9]{30,}|github_pat_[A-Za-z0-9_]{40,}|'
                    rb'AKIA[0-9A-Z]{16}|-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----)')
FIELDS = ['dataset', 'browse', 'task', 'variant', 'go_s', 'go_peak_rss_mib',
          'java_warm_s', 'java_peak_rss_mib', 'java_over_go_time', 'java_over_go_rss',
          'all_outputs_match', 'completed_run_count']


def sha(data):
    return hashlib.sha256(data).hexdigest()


def save(path, obj):
    path.write_text(json.dumps(obj, indent=2, allow_nan=False)+'\n')


def add_member(archive, name, data):
    info = tarfile.TarInfo(name)
    info.size, info.mtime, info.mode = len(data), 0, 0o644
    info.uid = info.gid = 0
    info.uname = info.gname = ''
    archive.addfile(info, io.BytesIO(data))


def verify(archive_path, records):
    expected = {'records/'+r['path']: r for r in records}
    seen = set()
    with tarfile.open(archive_path, 'r:gz') as archive:
        for member in archive:
            path = Path(member.name)
            if not member.isfile() or path.is_absolute() or '..' in path.parts or member.name in seen:
                raise ValueError('unsafe/duplicate archive member: '+member.name)
            item = expected[member.name]
            data = archive.extractfile(member).read()
            if len(data) != item['bytes'] or sha(data) != item['sha256']:
                raise ValueError('archive digest mismatch: '+member.name)
            seen.add(member.name)
    if seen != set(expected):
        raise ValueError('archive member inventory differs from manifest')


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--results', type=Path, default=Path('/workspace/jsgo-results'))
    parser.add_argument('--out', type=Path, default=Path('/workspace/jsgo-publish-records/docs/performance-experiments'))
    parser.add_argument('--exclude-prefix', action='append', default=[],
                        help='active result path prefix to exclude; repeat if needed')
    args = parser.parse_args()
    root, out = args.results.resolve(), args.out.resolve()
    live_prefixes = tuple(args.exclude_prefix)
    original_manifest_path = out/'measurement-manifest.json'
    original_manifest_bytes = original_manifest_path.read_bytes()
    followup_manifest_bytes = (out/'followup-measurement-manifest.json').read_bytes()
    second_manifest_bytes = (out/'second-followup-measurement-manifest.json').read_bytes()
    third_manifest_bytes = (out/'third-followup-measurement-manifest.json').read_bytes()
    original = {r['path']: r for r in json.loads(original_manifest_bytes)['files']}
    original.update({r['path']: r for r in json.loads(followup_manifest_bytes)['files']})
    original.update({r['path']: r for r in json.loads(second_manifest_bytes)['files']})
    original.update({r['path']: r for r in json.loads(third_manifest_bytes)['files']})
    original_archive_path = out/'measurement-records.tar.gz'
    original_archive_hash = sha(original_archive_path.read_bytes())
    followup_archive_hash = sha((out/'followup-measurement-records.tar.gz').read_bytes())
    second_archive_hash = sha((out/'second-followup-measurement-records.tar.gz').read_bytes())
    third_archive_hash = sha((out/'third-followup-measurement-records.tar.gz').read_bytes())
    original_protected = {p.name: sha(p.read_bytes()) for p in out.iterdir()
                          if p.is_file() and not p.name.startswith(('fourth-followup-', 'publish_fourth_followup_'))}
    archive_path = out/'fourth-followup-measurement-records.tar.gz'
    if archive_path.exists() or (out/'fourth-followup-records').exists():
        parser.error('refusing to overwrite an existing followup archive or browse tree')
    records, datasets = [], []
    excluded = collections.Counter()
    snapshot_json = {}
    text_scanned = 0
    # Explicit gzip mtime/filename and fixed tar metadata make the archive
    # deterministic for the same ordered input bytes, independent of wall time.
    with archive_path.open('wb') as raw:
        with gzip.GzipFile(filename='', mode='wb', fileobj=raw, mtime=0, compresslevel=6) as compressed:
            with tarfile.open(fileobj=compressed, mode='w', format=tarfile.GNU_FORMAT) as archive:
                for path in sorted(root.rglob('*')):
                    rel = path.relative_to(root).as_posix()
                    if any(rel.startswith(prefix) for prefix in live_prefixes) or path.name == 'CURRENT.json':
                        excluded['explicit live prefix / CURRENT.json'] += 1
                        continue
                    if path.is_symlink():
                        excluded['symlink'] += 1
                        continue
                    if not path.is_file():
                        continue
                    before = path.stat()
                    ext = path.suffix
                    include = ext in KEEP_EXT or (ext == '' and before.st_size < 2000000 and path.name in KEEP_NAMES)
                    if not include:
                        excluded[ext or 'other/no-extension'] += 1
                        continue
                    data = path.read_bytes()
                    after = path.stat()
                    if (before.st_ino, before.st_size, before.st_mtime_ns) != (after.st_ino, after.st_size, after.st_mtime_ns):
                        excluded['changed during snapshot; not published'] += 1
                        continue
                    if data.startswith(b'\x7fELF'):
                        excluded['ELF binary'] += 1
                        continue
                    digest = sha(data)
                    prior = original.get(rel)
                    if prior and prior['sha256'] == digest:
                        excluded['unchanged from prior archives'] += 1
                        continue
                    if ext not in {'.pprof', '.pgo'}:
                        if SECRET.search(data):
                            raise ValueError('possible credential pattern; inspect privately: '+rel)
                        text_scanned += 1
                    item = {'path': rel, 'bytes': len(data), 'sha256': digest,
                            'change': 'changed' if prior else 'new',
                            'original_sha256': prior['sha256'] if prior else None}
                    records.append(item)
                    add_member(archive, 'records/'+rel, data)
                    browse = ext == '.md' or path.name in {'summary.json', 'machine.json'} or (
                        ext == '.json' and len(Path(rel).parts) <= 2 and len(data) < 2000000)
                    if browse:
                        dest = out/'fourth-followup-records'/rel
                        dest.parent.mkdir(parents=True, exist_ok=True)
                        dest.write_bytes(data)
                    if ext == '.json':
                        try:
                            snapshot_json[rel] = json.loads(data)
                        except (UnicodeDecodeError, json.JSONDecodeError):
                            # Invalid/partial JSON remains useful failure evidence.
                            pass
                    if path.name == 'summary.json':
                        obj = snapshot_json.get(rel)
                        datasets.append({'path': rel, 'browse': 'fourth-followup-records/'+rel,
                                         'keys': list(obj) if isinstance(obj, dict) else [],
                                         'sha256': digest, 'change': item['change'],
                                         'parseable_json': rel in snapshot_json})
    by_path = {r['path']: r for r in records}

    def frozen_json(path):
        rel = path.relative_to(root).as_posix()
        if rel in snapshot_json:
            return snapshot_json[rel]
        expected = by_path.get(rel) or original.get(rel)
        if not expected or not path.is_file() or path.is_symlink():
            return None
        data = path.read_bytes()
        if sha(data) != expected['sha256']:
            return None
        try:
            return json.loads(data)
        except (UnicodeDecodeError, json.JSONDecodeError):
            return None

    # Only the ordinary experiments.py completed-variant schema is normalized.
    # Diagnostics, protocol summaries and unusual schemas remain unmodified.
    metrics, skipped_metrics = [], []
    for dataset in datasets:
        obj = snapshot_json.get(dataset['path'])
        directory = root/Path(dataset['path']).parent
        machine = frozen_json(directory/'machine.json') or {}
        expected_count = machine.get('runs')
        if not isinstance(obj, dict):
            continue
        for task, row in obj.items():
            if not isinstance(row, dict) or not isinstance(row.get('variants'), dict):
                continue
            for variant, values in row['variants'].items():
                reason = None
                raw = frozen_json(directory/task/variant/'jsgo.json') or {}
                runs = raw.get('runs', [])
                if not isinstance(values, dict) or 'median_s' not in values:
                    reason = 'no ordinary completed-variant metrics'
                elif not isinstance(expected_count, int) or expected_count < 1 or len(runs) != expected_count:
                    reason = 'cannot verify configured completed raw run count'
                elif any(not isinstance(r, dict) or not isinstance(r.get('wall_s'), (float, int)) or
                         not math.isfinite(r['wall_s']) or r['wall_s'] <= 0 for r in runs):
                    reason = 'invalid/missing raw wall time'
                elif not math.isclose(statistics.median(r['wall_s'] for r in runs), values['median_s'], rel_tol=1e-12):
                    reason = 'saved median does not match completed raw observations'
                elif values.get('all_outputs_match') is True:
                    gate = frozen_json(directory/task/variant/'gate.json') or {}
                    if gate.get('status') != 'MATCH' or any(
                            not r.get('exits') or any(code != 0 for code in r['exits']) or
                            r.get('digests') != runs[0].get('digests') for r in runs):
                        reason = 'claimed output match lacks zero-exit stable-digest raw runs and MATCH gate'
                if reason:
                    skipped_metrics.append({'dataset': dataset['path'], 'task': task, 'variant': variant, 'reason': reason})
                    continue
                metrics.append({'dataset': dataset['path'], 'browse': dataset['browse'], 'task': task,
                                'variant': variant, 'go_s': values.get('median_s'),
                                'go_peak_rss_mib': values.get('median_peak_rss_mb'),
                                'java_warm_s': row.get('java_warm_s'),
                                'java_peak_rss_mib': row.get('java_peak_rss_mb'),
                                'java_over_go_time': values.get('speed_vs_java_warm'),
                                'java_over_go_rss': values.get('memory_ratio_java_over_go'),
                                'all_outputs_match': values.get('all_outputs_match'),
                                'completed_run_count': len(runs)})
    save(out/'fourth-followup-datasets.json', datasets)
    save(out/'fourth-followup-metrics.json', metrics)
    with (out/'fourth-followup-metrics.csv').open('w', newline='') as stream:
        writer = csv.DictWriter(stream, fieldnames=FIELDS, lineterminator='\n')
        writer.writeheader()
        writer.writerows(metrics)
    verify(archive_path, records)
    for path in (out/'fourth-followup-records').rglob('*'):
        if path.is_file():
            rel = path.relative_to(out/'fourth-followup-records').as_posix()
            if sha(path.read_bytes()) != by_path[rel]['sha256']:
                raise ValueError('browse copy differs from tar snapshot: '+rel)
    protected_after = {name: sha((out/name).read_bytes()) for name in original_protected}
    if protected_after != original_protected:
        raise ValueError('a protected original publication file changed during this run')
    manifest = {
        'schema': 'jsgo-measurement-fourth-followup-v1', 'root': 'jsgo-results',
        'base_manifest': {'path': 'measurement-manifest.json', 'sha256': sha(original_manifest_bytes)},
        'base_archive': {'path': 'measurement-records.tar.gz', 'sha256': original_archive_hash},
        'first_followup_manifest': {'path': 'followup-measurement-manifest.json', 'sha256': sha(followup_manifest_bytes)},
        'first_followup_archive': {'path': 'followup-measurement-records.tar.gz', 'sha256': followup_archive_hash},
        'second_followup_manifest': {'path': 'second-followup-measurement-manifest.json', 'sha256': sha(second_manifest_bytes)},
        'second_followup_archive': {'path': 'second-followup-measurement-records.tar.gz', 'sha256': second_archive_hash},
        'third_followup_manifest': {'path': 'third-followup-measurement-manifest.json', 'sha256': sha(third_manifest_bytes)},
        'third_followup_archive': {'path': 'third-followup-measurement-records.tar.gz', 'sha256': third_archive_hash},
        'archive': {'path': archive_path.name, 'bytes': archive_path.stat().st_size,
                    'sha256': sha(archive_path.read_bytes())},
        'member_prefix': 'records/', 'browse_prefix': 'fourth-followup-records/', 'files': records,
        'excluded_counts': dict(sorted(excluded.items())), 'excluded_live_prefixes': list(live_prefixes),
        'excluded_names': ['CURRENT.json'],
        'policy': 'Same extension/name policy as publish_records.py. Only new or changed hashes versus the union of original, first-followup, second-followup and third-followup manifests are included. Retains failures, negative, inconclusive and source-only evidence. Excludes compiled binaries/classes, downloaded corpora, generated JS/maps, symlink targets and operational PID files. Original records are never overwritten; changed paths supplement the historical earlier archives.',
        'normalization': 'Only completed ordinary experiments.py variants with a verified configured raw run count and matching median. Existing supplied fields are copied; absent RSS/ratios stay null. Diagnostic/nonstandard/incomplete summary schemas remain available in full.',
        'skipped_metric_variants': skipped_metrics,
        'validation': {'archive_members_and_hashes': 'PASS', 'browse_copy_hashes': 'PASS',
                       'original_publication_files_unchanged': 'PASS',
                       'text_secret_pattern_scan': 'PASS; matching credential patterns: 0',
                       'text_files_scanned': text_scanned, 'dataset_count': len(datasets),
                       'normalized_metric_rows': len(metrics),
                       'normalized_match_rows_raw_exits_digests_gates': 'PASS',
                       'archive_gzip_mtime': 0, 'archive_member_metadata': 'fixed, sorted, regular files only'},
    }
    save(out/'fourth-followup-measurement-manifest.json', manifest)
    print(json.dumps({'files': len(records), 'new': sum(r['change'] == 'new' for r in records),
                      'changed': sum(r['change'] == 'changed' for r in records),
                      'datasets': len(datasets), 'metric_rows': len(metrics),
                      'archive': manifest['archive'], 'validation': manifest['validation']}, indent=2))


if __name__ == '__main__':
    main()
