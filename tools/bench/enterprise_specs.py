#!/usr/bin/env python3
"""Freeze enterprise CLI argv, input bytes, chunks and extern identities.

Existing main tasks are reused verbatim except for an explicit thread flag and
the same output manifest on both engines. typed-full-chunks partitions that
actual 1,818-input workload into four dependent ADVANCED chunks; no replication.
This tool prepares files only; it does not run a compiler or fetch dependencies.
"""
import argparse
import hashlib
import io
import json
from pathlib import Path
import sys
import zipfile

sys.dont_write_bytecode = True
sys.path.insert(0, str(Path(__file__).resolve().parent))
import bench


def sha(data):
    return hashlib.sha256(data).hexdigest()


def replace_threads(argv, number):
    result, index = [], 0
    while index < len(argv):
        value = argv[index]
        if value == '--':
            return result+[f'--num_parallel_threads={number}']+argv[index:]
        if value == '--num_parallel_threads':
            if index+1 >= len(argv):
                raise ValueError('missing existing thread value')
            index += 2
        elif value.startswith('--num_parallel_threads='):
            index += 1
        else:
            result.append(value); index += 1
    return result+[f'--num_parallel_threads={number}']


def embedded_externs(jar, repo):
    with zipfile.ZipFile(jar) as archive:
        candidates = [name for name in archive.namelist() if name.endswith('externs.zip')]
        if len(candidates) != 1:
            raise ValueError('cannot identify pinned stock extern resource')
        raw = archive.read(candidates[0])
    embedded = (repo/'internal/resources/externs.zip').read_bytes()
    # ZIP metadata may differ. Compare every named source's actual bytes.
    def sources(data):
        with zipfile.ZipFile(io.BytesIO(data)) as archive:
            return {name: sha(archive.read(name)) for name in sorted(archive.namelist()) if not name.endswith('/')}
    java, go = sources(raw), sources(embedded)
    if java != go:
        raise ValueError('Go/Java builtin extern corpus differs')
    return {'java_resource': candidates[0], 'java_archive_sha256': sha(raw),
            'go_archive_sha256': sha(embedded), 'source_files': java,
            'assertion': 'all builtin extern source names and bytes match; parsed fresh on every compilation'}


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument('--tasks', default='typed-full,typed-full-chunks,bundle-advanced,bundle-advanced-ts,shaka-blockly')
    p.add_argument('--cache', default=bench.default_cache())
    p.add_argument('--jar', required=True)
    p.add_argument('--out', required=True)
    p.add_argument('--threads', type=int, default=1)
    a = p.parse_args()
    if a.threads < 1:
        p.error('threads must be positive')
    out, repo = Path(a.out).resolve(), Path(bench.REPO).resolve()
    if out.exists():
        p.error('refusing to overwrite existing input identity evidence')
    wanted = a.tasks.split(',')
    names = set(wanted)-{'typed-full-chunks'}
    if 'typed-full-chunks' in wanted:
        names.add('typed-full')
    tasks = {task['name']: task for task in bench.load_tasks(list(names))}
    externs = embedded_externs(a.jar, repo)
    out.mkdir(parents=True)
    setup = out/'workspace'
    bench.make_workspace(str(setup), a.cache)
    for name in wanted:
        task = tasks['typed-full' if name == 'typed-full-chunks' else name]
        jobs = bench.task_jobs(task, str(setup))
        if len(jobs) != 1:
            raise ValueError('enterprise paired unit must contain one compilation')
        argv = jobs[0]
        chunks = []
        if name == 'typed-full-chunks':
            inputs = bench.expand_inputs(task['inputs'], str(setup))
            groups = {key: [] for key in ('core', 'shaka', 'blockly', 'app')}
            for index in range(0, len(inputs), 2):
                flag, path = inputs[index:index+2]
                group = ('core' if '/closure-library@' in '/'+path else
                         'shaka' if '/shaka-player@' in '/'+path else
                         'blockly' if '/blockly@' in '/'+path else 'app' if path.startswith('app/') else None)
                if group is None or flag != '--js':
                    raise ValueError('unexpected chunk input: '+path)
                groups[group] += [flag, path]
            argv = list(task['flags'])
            for group, values in groups.items():
                if not values:
                    raise ValueError('empty input chunk: '+group)
                chunks.append({'name': group, 'input_count': len(values)//2,
                               'dependencies': [] if group == 'core' else ['core']})
                argv += ['--chunk', group+':'+str(len(values)//2)+('' if group == 'core' else ':core')]+values
            argv += ['--chunk_output_path_prefix=chunks/', '--create_source_map=%outname%.map',
                     '--output_manifest=%outname%.inputs', '--output_chunk_dependencies=chunk-dependencies.json']
            artifacts = [f'chunks/{chunk["name"]}.js{suffix}' for chunk in chunks for suffix in ('', '.map', '.inputs')]
            artifacts.append('chunk-dependencies.json')
        else:
            argv += ['--output_manifest=inputs.manifest']
            artifacts = list(task.get('compare', []))+['inputs.manifest']
        argv = replace_threads(argv, a.threads)
        links = {'inputs': str(Path(a.cache).resolve()/'bench-inputs'),
                 'corpus': str(Path(a.cache).resolve()/'corpus'), 'app': str(Path(bench.HERE)/'app')}
        source_inputs = []
        for index, value in enumerate(argv[:-1]):
            if value not in ('--js', '--externs'):
                continue
            path = Path(argv[index+1]); link, *relative = path.parts
            if link not in links or '..' in path.parts:
                raise ValueError('source path is not in the frozen link namespace')
            physical = Path(links[link]).joinpath(*relative)
            data = physical.read_bytes()
            source_inputs.append({'flag': value, 'link': link, 'relative': '/'.join(relative),
                                  'bytes': len(data), 'sha256': sha(data)})
        spec = {'schema': 'jsgo-enterprise-cli-v1', 'name': name, 'source_task': task['name'],
                'argv': argv, 'links': links, 'inputs': source_inputs, 'artifacts': artifacts,
                'same_work': {'mode': task['mode'], 'compiler_threads': a.threads, 'chunks': chunks,
                              'source_count': len(source_inputs), 'source_bytes': sum(x['bytes'] for x in source_inputs),
                              'builtin_externs': externs, 'explicit_externs': [x for x in source_inputs if x['flag'] == '--externs'],
                              'validation': ['identical complete argv on both engines', 'identical input and builtin extern bytes',
                                             'strict stdout/stderr/exit and ALL output artifacts for every warmup/measurement',
                                             'actual emitted input manifests and chunk graph match, not only final JS',
                                             'failure/OOM is a failed task, never a speed ratio']}}
        bench.save(str(out/(name+'.json')), spec)
        print(name, len(source_inputs), 'files', sum(x['bytes'] for x in source_inputs), 'bytes', 'chunks', chunks)


if __name__ == '__main__':
    main()
