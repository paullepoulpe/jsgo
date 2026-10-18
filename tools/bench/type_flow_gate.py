#!/usr/bin/env python3
"""Exact CLI oracle gate for indexed inferred-type slots; no timing/building.

Runs the same typed programs through pinned Java and both Go storage backends.
Checks exit status, stdout, ordered stderr, JavaScript and source maps byte for
byte. Requires prebuilt binaries. The experimental backend remains opt-in.
"""
import argparse
import json
import os
from pathlib import Path
import shutil
import subprocess
import sys

sys.dont_write_bytecode = True
sys.path.insert(0, str(Path(__file__).resolve().parent))
import bench


PROGRAMS = {
    'forward-alias': '''/** @typedef {(Later|string)} */ var Early;
/** @param {Early} x @param {boolean} c */
function f(x,c) { if(c) {x=1;} else if(typeof x==='string') {x=x.length;} return x; }
/** @typedef {number} */ var Later;
window['result']=f(window['value'],window['condition']);
''',
    'recursive-object': '''/** @constructor */ function Link() { /** @type {?Link} */ this.next=null; }
/** @param {?Link} node @param {boolean} c */
function walk(node,c) { while(node && c) {if(node.next) {node=node.next;} else {break;}} return node; }
window['result']=walk(new Link(),window['condition']);
''',
    'unknown-branches': '''/** @param {?} x @param {boolean} c */
function unknown(x,c) { if(c) {x=3;} else {x='x';} if(typeof x==='string') {return x.length;} return x; }
window['result']=unknown(window['value'],window['condition']);
''',
    'hash-collisions': '''/** @param {(number|string)} Aa @param {(number|string)} BB @param {boolean} c */
function collisions(Aa,BB,c) {if(c) {Aa=1;BB='x';} else {Aa='x';BB=2;} return Aa+BB;}
window['result']=collisions(window['a'],window['b'],window['condition']);
''',
    'shadow-catch': '''/** @param {{value:(number|string)}} x @param {boolean} c */
function shadow(x,c) { if(c) {let x={value:1};x.value=2;} else {x.value='x';}
try {if(c) {throw x;}} catch(e) {e.value=3;} return x.value; }
window['result']=shadow({value:window['value']},window['condition']);
''',
}
PROGRAMS['large-scope-fallback'] = (
    '/** @param {boolean} c */ function large(c) {\n'
    + ''.join(f'var x{i}={i};\n' for i in range(280))
    + 'if(c) {\n' + ''.join(f'x{i}="x";\n' for i in range(280))
    + '} return x0+x279;}\nwindow["result"]=large(window["condition"]);\n'
)


def execute(command, flags, workspace, destination, env, timeout):
    output = workspace / '.type-flow-output'
    if output.exists():
        shutil.rmtree(output)
    output.mkdir()
    result = subprocess.run(command + flags, cwd=workspace, env=env,
                            capture_output=True, timeout=timeout)
    destination.mkdir(parents=True, exist_ok=True)
    (destination / 'stdout').write_bytes(result.stdout)
    (destination / 'stderr').write_bytes(result.stderr)
    (destination / 'exit').write_text(str(result.returncode) + '\n')
    artifacts = {str(p.relative_to(output)): p.read_bytes()
                 for p in sorted(output.rglob('*')) if p.is_file()}
    for name, content in artifacts.items():
        path = destination / 'artifacts' / name
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_bytes(content)
    return result.returncode, result.stdout, result.stderr, artifacts


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--jsgo', required=True, type=Path)
    parser.add_argument('--out', required=True, type=Path)
    parser.add_argument('--java', default='java')
    parser.add_argument('--jar', type=Path)
    parser.add_argument('--threads', default='1')
    parser.add_argument('--cases', help='comma-separated case names; default all')
    parser.add_argument('--timeout', type=int, default=90)
    args = parser.parse_args()
    binary = args.jsgo.resolve()
    jar = (args.jar or Path(bench.default_cache()) / bench.read_pin()['CLOSURE_JAR']).resolve()
    root = args.out.resolve()
    workspace = root / 'fixtures'
    workspace.mkdir(parents=True, exist_ok=True)
    programs = PROGRAMS
    if args.cases:
        selected = args.cases.split(',')
        if set(selected) - programs.keys():
            parser.error('unknown cases')
        programs = {name: programs[name] for name in selected}
    threads = [int(value) for value in args.threads.split(',')]
    if any(value < 1 for value in threads):
        parser.error('thread counts must be positive')
    env = {key: value for key, value in os.environ.items() if key not in
           ('GOGC', 'GODEBUG', 'GOMEMLIMIT', 'GOMAXPROCS', 'JAVA_TOOL_OPTIONS',
            '_JAVA_OPTIONS', 'JDK_JAVA_OPTIONS', 'JSGO_EXPERIMENT_INDEXED_FLOW_SLOTS')}
    env['LC_ALL'] = 'C.UTF-8'
    summary = {'jsgo': str(binary), 'jsgo_sha256': bench.sha256_file(binary),
               'jar': str(jar), 'jar_sha256': bench.sha256_file(jar), 'cases': {}}
    failures = []
    for name, source in programs.items():
        (workspace / 'input.js').write_text(source, encoding='utf-8')
        for mode in ('checks', 'advanced'):
            for count in threads:
                key = f'{name}/{mode}/threads-{count}'
                flags = [f'--num_parallel_threads={count}', '--compilation_level=ADVANCED',
                         '--jscomp_warning=checkTypes', '--language_in=ECMASCRIPT_NEXT',
                         '--js_output_file=.type-flow-output/out.js',
                         '--create_source_map=.type-flow-output/out.js.map', '--js=input.js']
                if mode == 'checks':
                    flags += ['--checks_only']
                destination = root / name / mode / f'threads-{count}'
                oracle = execute([args.java, '-jar', str(jar)], flags, workspace,
                                 destination / 'java', env, args.timeout)
                control = execute([str(binary)], flags, workspace, destination / 'go-hamt',
                                  dict(env, JSGO_EXPERIMENT_INDEXED_FLOW_SLOTS='0'), args.timeout)
                candidate = execute([str(binary)], flags, workspace, destination / 'go-indexed',
                                    dict(env, JSGO_EXPERIMENT_INDEXED_FLOW_SLOTS='1'), args.timeout)
                checks = {'java_exit': oracle[0], 'go_hamt_matches_java': control == oracle,
                          'go_indexed_matches_java': candidate == oracle,
                          'go_backends_match': candidate == control, 'args': flags}
                checks['match'] = oracle[0] == 0 and all(checks[k] for k in
                                    ('go_hamt_matches_java', 'go_indexed_matches_java', 'go_backends_match'))
                summary['cases'][key] = checks
                print(f'{key}: {"MATCH" if checks["match"] else "FAIL"}', flush=True)
                if not checks['match']:
                    failures.append(key)
                (root / 'summary.json').write_text(json.dumps(summary, indent=2) + '\n')
    if failures:
        raise SystemExit('type-flow CLI gate failed: ' + ', '.join(failures))
    print(f'{len(summary["cases"])} cases passed; results: {root / "summary.json"}', flush=True)


if __name__ == '__main__':
    main()
