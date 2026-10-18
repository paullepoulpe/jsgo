#!/usr/bin/env python3
"""Byte/race gate for parallel prebuilding, using prebuilt Go and pinned Java.

No compilation or timing takes place. Every case runs sequentially at each
requested thread count and compares exit status, stdout, stderr, JS and maps.
Only the benchmark suite's existing Java-exception stack-frame rule is allowed.

  python3 -I tools/bench/parallel_gate.py --jsgo /path/jsgo-race --out /tmp/parallel-gate
"""
import argparse
import base64
import json
import os
from pathlib import Path
import shutil
import subprocess
import sys

sys.dont_write_bytecode = True
sys.path.insert(0, str(Path(__file__).resolve().parent))
import bench


def fixtures(root):
    def write(name, content):
        path = root / name
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(content, encoding='utf-8')

    sources = []
    for i in range(12):
        name = f'plain/part{i:02}.js'
        write(name, f'/** @param {{number}} value @return {{number}} */\n'
                    f'function f{i}(value) {{ return value + {i}; }}\n'
                    f'var name{i}="é😀\\ud800";\n')
        sources += ['--js', name]
    for i, code in enumerate(['function syntaxA( {\n', 'let b = ;\n', 'const c = ;\n']):
        write(f'error-{i}.js', code)
    for i in range(3):
        write(f'warning-{i}.js', f'/** @type {{string}} */ var wrongType{i} = {i};\n')
    write('closure/base.js', '/** @provideGoog */ var goog = {};\n')
    write('closure/dep.js', "goog.provide('gate.dep'); gate.dep.value = 3;\n")
    write('closure/main.js', "goog.module('gate.main');\nvar dep = goog.require('gate.dep');\nwindow['result'] = dep.value;\n")
    write('closure/dead.js', "goog.provide('gate.dead'); gate.dead = 123;\n")
    write('esm/lib.js', 'export const value = 7;\n')
    write('esm/main.js', "import {value} from './lib.js'; window['result'] = value;\n")
    write('esm/missing-a.js', "import {a} from './absent-a.js'; window['missingA'] = a;\n")
    write('esm/missing-b.js', "import {b} from './absent-b.js'; window['missingB'] = b;\n")
    write('node_modules/gate/index.js', 'export const value = 9;\n')
    write('node/main.js', "import {value} from 'gate'; window['result'] = value;\n")
    write('prefix/main.js', "import {value} from '@gate/lib.js'; window['result'] = value;\n")
    for i in range(3):
        source_map = {'version': 3, 'file': f'inline-{i}.js', 'sources': [f'original-{i}.ts'],
                      'sourcesContent': [f'// original 😀 {i}\nlet original{i}: string = "é";\n'],
                      'names': [], 'mappings': 'AAAA;AACA'}
        encoded = base64.b64encode(json.dumps(source_map, separators=(',', ':')).encode()).decode()
        write(f'inline-{i}.js', f'/** @license gate license {i} */\n'
                               f'/** @type {{string}} */ var mapped{i} = "é😀\\ud800";\n'
                               f'//# sourceMappingURL=data:application/json;base64,{encoded}\n')
    write('unicode-error-a.js', 'var unicode = "😀"; let wrong = ;\n')
    write('unicode-error-b.js', 'var unicode2 = "é😀"; const wrong2 = ;\n')

    def js(*names):
        return [arg for name in names for arg in ('--js', name)]

    modules = js('esm/main.js', 'esm/lib.js')
    missing = js('esm/missing-b.js', 'esm/missing-a.js', 'esm/lib.js')
    prune = ['--dependency_mode=PRUNE', '--entry_point=esm/main.js']
    return {
        'whitespace': ['--compilation_level=WHITESPACE_ONLY'] + sources,
        'simple': ['--compilation_level=SIMPLE'] + sources,
        'advanced': ['--compilation_level=ADVANCED', '--jscomp_warning=checkTypes'] + sources,
        'parse-errors-ordered': ['--compilation_level=SIMPLE'] + js('error-2.js', 'error-0.js', 'error-1.js'),
        'type-warnings-ordered': ['--compilation_level=ADVANCED', '--jscomp_warning=checkTypes']
                                 + js('warning-2.js', 'warning-0.js', 'warning-1.js'),
        'missing-files-ordered': ['--compilation_level=SIMPLE'] + js('missing-b.js', 'missing-a.js'),
        'closure-prune-regex': ['--compilation_level=ADVANCED', '--dependency_mode=PRUNE', '--entry_point=goog:gate.main']
                               + js('closure/dead.js', 'closure/main.js', 'closure/base.js', 'closure/dep.js'),
        'esm-browser': ['--compilation_level=SIMPLE', '--module_resolution=BROWSER'] + modules,
        'esm-browser-prune': ['--compilation_level=ADVANCED', '--module_resolution=BROWSER'] + prune + modules,
        'esm-node': ['--compilation_level=SIMPLE', '--module_resolution=NODE'] + js('node/main.js', 'node_modules/gate/index.js'),
        'esm-prefix': ['--compilation_level=SIMPLE', '--module_resolution=BROWSER_WITH_TRANSFORMED_PREFIXES',
                       '--browser_resolver_prefix_replacements=@gate/=esm/'] + js('prefix/main.js', 'esm/lib.js'),
        'missing-module-warnings': ['--compilation_level=SIMPLE', '--checks_only', '--module_resolution=BROWSER', '--jscomp_warning=moduleLoad'] + missing,
        'missing-module-prune': ['--compilation_level=SIMPLE', '--checks_only', '--module_resolution=BROWSER', '--dependency_mode=PRUNE',
                                 '--entry_point=esm/missing-b.js', '--entry_point=esm/missing-a.js', '--jscomp_warning=moduleLoad'] + missing,
        'inline-maps-comments': ['--compilation_level=WHITESPACE_ONLY', '--parse_inline_source_maps', '--apply_input_source_maps',
                                 '--source_map_include_content', '--preserve_type_annotations'] + js('inline-2.js', 'inline-0.js', 'inline-1.js'),
        'utf16-diagnostic-order': ['--compilation_level=SIMPLE'] + js('unicode-error-b.js', 'unicode-error-a.js'),
    }


def run(command, args, cwd, destination, env, timeout):
    outputs = cwd / '.gate-output'
    if outputs.exists():
        shutil.rmtree(outputs)
    outputs.mkdir()
    result = subprocess.run(command + args, cwd=cwd, env=env, capture_output=True, timeout=timeout)
    destination.mkdir(parents=True, exist_ok=True)
    (destination / 'stdout').write_bytes(result.stdout)
    (destination / 'stderr').write_bytes(result.stderr)
    (destination / 'exit').write_text(str(result.returncode) + '\n')
    artifacts = {str(p.relative_to(outputs)): p.read_bytes() for p in sorted(outputs.rglob('*')) if p.is_file()}
    for name, content in artifacts.items():
        path = destination / 'artifacts' / name
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_bytes(content)
    return result, artifacts


def main():
    parser = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument('--jsgo', required=True, type=Path, help='prebuilt binary; use -race to check races')
    parser.add_argument('--out', required=True, type=Path)
    parser.add_argument('--java', default='java')
    parser.add_argument('--jar', type=Path)
    parser.add_argument('--threads', default='1,2,4')
    parser.add_argument('--cases', help='comma-separated case names; default all')
    parser.add_argument('--timeout', type=int, default=90)
    a = parser.parse_args()
    a.jsgo = a.jsgo.resolve()
    a.jar = (a.jar or Path(bench.default_cache()) / bench.read_pin()['CLOSURE_JAR']).resolve()
    threads = [int(value) for value in a.threads.split(',')]
    if any(value < 1 for value in threads):
        parser.error('thread counts must be positive')
    root = a.out.resolve()
    root.mkdir(parents=True, exist_ok=True)
    workspace = root / 'fixtures'
    workspace.mkdir(exist_ok=True)
    cases = fixtures(workspace)
    if a.cases:
        chosen = a.cases.split(',')
        unknown = set(chosen) - cases.keys()
        if unknown:
            parser.error('unknown cases: ' + ','.join(sorted(unknown)))
        cases = {name: cases[name] for name in chosen}
    env = {key: value for key, value in os.environ.items() if key not in
           ('GOGC', 'GODEBUG', 'GOMEMLIMIT', 'GOMAXPROCS', 'JAVA_TOOL_OPTIONS', '_JAVA_OPTIONS', 'JDK_JAVA_OPTIONS')}
    env['LC_ALL'] = 'C.UTF-8'
    summary = {'jsgo': str(a.jsgo), 'jsgo_sha256': bench.sha256_file(a.jsgo), 'jar': str(a.jar),
               'jar_sha256': bench.sha256_file(a.jar), 'threads': threads, 'cases': {}}
    failures = []
    for name, flags in cases.items():
        for count in threads:
            args = [f'--num_parallel_threads={count}', '--js_output_file=.gate-output/out.js',
                    '--create_source_map=.gate-output/out.js.map'] + flags
            destination = root / name / f'threads-{count}'
            java, java_files = run([a.java, '-jar', str(a.jar)], args, workspace, destination / 'java', env, a.timeout)
            go, go_files = run([str(a.jsgo)], args, workspace, destination / 'go', env, a.timeout)
            race = b'WARNING: DATA RACE' in go.stderr or b'ThreadSanitizer' in go.stderr
            stderr_match = java.stderr == go.stderr
            frames_only = (not stderr_match and bench.without_java_frames(java.stderr.decode('utf-8', 'surrogateescape'))
                           == bench.without_java_frames(go.stderr.decode('utf-8', 'surrogateescape')))
            checks = {'exit_match': java.returncode == go.returncode, 'stdout_match': java.stdout == go.stdout,
                      'stderr_match': stderr_match, 'java_frames_only': frames_only,
                      'artifacts_match': java_files == go_files, 'race_detected': race,
                      'java_exit': java.returncode, 'go_exit': go.returncode, 'args': args,
                      'artifact_names': sorted(java_files)}
            checks['match'] = checks['exit_match'] and checks['stdout_match'] and (stderr_match or frames_only) and checks['artifacts_match'] and not race
            summary['cases'][f'{name}/threads-{count}'] = checks
            print(f'{name}/threads-{count}: {"MATCH" if checks["match"] else "FAIL"}', flush=True)
            if not checks['match']:
                failures.append(f'{name}/threads-{count}')
            (root / 'summary.json').write_text(json.dumps(summary, indent=2) + '\n')
    if failures:
        raise SystemExit('parallel gate failed: ' + ', '.join(failures))
    print(f'{len(summary["cases"])} cases passed; results: {root / "summary.json"}', flush=True)


if __name__ == '__main__':
    main()
