#!/usr/bin/env python3
"""Byte-check mixed JSON-worker CLI jobs against the pinned stock Java jar.

This is a correctness gate, not a timing benchmark. It uses one persistent
worker across source/options/errors/thread changes and repeats builtin extern
jobs. Optional Java cache counters prove reuse outside measured runs.
"""
import argparse
import base64
import json
import hashlib
import os
from pathlib import Path
import re
import selectors
import subprocess
import sys

sys.path.insert(0, str(Path(__file__).resolve().parent))
sys.dont_write_bytecode = True
import bench
from worker_compare import thread_flags


def normalize_frames(data):
    # The Go port omits JVM stack frames; preserve messages, status, and all
    # remaining bytes, matching the existing CLI comparison convention.
    return b'\n'.join(line for line in data.split(b'\n')
                      if not (line.startswith(b'\tat ') or line.startswith(b'\t... ')))



ORACLE_SCHEMA = "jsgo-worker-stock-oracle-v1"
OUTPUT_FILES = ("out.js", "out.js.map")
STDERR_NORMALIZATION = "drop-jvm-stack-frame-lines-v1"


def object_digest(value):
    return hashlib.sha256(json.dumps(value, sort_keys=True, separators=(',', ':'),
                                     ensure_ascii=True).encode()).hexdigest()


def case_identity(flags, source, stdin):
    source_bytes = source.encode('utf-8')
    identity = {'args': flags, 'source_file': 'app.js',
                'source_utf8_base64': base64.b64encode(source_bytes).decode(),
                'source_sha256': hashlib.sha256(source_bytes).hexdigest(),
                'stdin_base64': base64.b64encode(stdin).decode(),
                'stdin_sha256': hashlib.sha256(stdin).hexdigest()}
    identity['sha256'] = object_digest(identity)
    return identity


def encode_expected(expected):
    payload = {'exit': expected['exit'],
               'stdout_base64': base64.b64encode(expected['stdout']).decode(),
               'stderr_base64': base64.b64encode(expected['stderr']).decode(),
               'files_base64': {name: None if expected[name] is None else
                                base64.b64encode(expected[name]).decode() for name in OUTPUT_FILES}}
    return {'payload': payload, 'sha256': object_digest(payload)}


def decode_expected(record):
    payload = record['payload']
    if set(payload) != {'exit', 'stdout_base64', 'stderr_base64', 'files_base64'}:
        raise ValueError('oracle expected-response schema mismatch')
    if record['sha256'] != object_digest(payload):
        raise ValueError('oracle expected-response SHA-256 mismatch')
    if type(payload['exit']) is not int or not 0 <= payload['exit'] <= 255:
        raise ValueError('oracle exit must be an actual normal process exit code')
    if set(payload['files_base64']) != set(OUTPUT_FILES):
        raise ValueError('oracle output file set mismatch')
    def decode(value):
        if not isinstance(value, str): raise ValueError('oracle byte payload must be a base64 string')
        return base64.b64decode(value, validate=True)
    expected = {'exit': payload['exit'], 'stdout': decode(payload['stdout_base64']),
                'stderr': decode(payload['stderr_base64'])}
    expected.update({name: None if value is None else decode(value)
                     for name, value in payload['files_base64'].items()})
    return expected


def load_oracle(path, jar_sha256, identities):
    path = Path(path).resolve()
    if path.is_dir(): path = path/'oracle.json'
    oracle = json.loads(path.read_text())
    if oracle.get('schema') != ORACLE_SCHEMA or oracle.get('complete') is not True:
        raise ValueError('need a complete current-schema stock oracle')
    if oracle.get('stderr_normalization') != STDERR_NORMALIZATION:
        raise ValueError('oracle stderr comparison rule mismatch')
    if oracle.get('stock_jar_sha256') != jar_sha256:
        raise ValueError('oracle pinned stock jar SHA-256 mismatch')
    if oracle.get('case_order') != list(identities) or set(oracle.get('cases', {})) != set(identities):
        raise ValueError('oracle complete case/thread sequence mismatch')
    if oracle.get('corpus_sha256') != object_digest(identities):
        raise ValueError('oracle fixture corpus SHA-256 mismatch')
    expected = {}
    for key, identity in identities.items():
        record = oracle['cases'][key]
        if record.get('input') != identity:
            raise ValueError('oracle full argv/source/stdin mismatch: '+key)
        expected[key] = decode_expected(record['expected'])
    if oracle.get('oracle_sha256') != object_digest(oracle['cases']):
        raise ValueError('oracle complete response corpus SHA-256 mismatch')
    return oracle, expected, str(path)


def experiment_environment(env):
    return {key: value for key, value in sorted(env.items()) if key.startswith('JSGO_EXPERIMENT_')}


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument('--command-json', required=True, help='JSON file containing argv array or build manifest')
    ap.add_argument('--out', required=True)
    ap.add_argument('--cache', default=bench.default_cache())
    ap.add_argument('--threads', default='1,2,4')
    ap.add_argument('--timeout', type=float, default=120)
    ap.add_argument('--oracle-reference', help='explicit complete stock oracle.json (or its directory); validates the entire fixture and skips fresh Java launches')
    ap.add_argument('--require-cache-hits', action='store_true',
                    help='require Java verifyCache counters with nonzero hits and misses')
    args = ap.parse_args()
    command = json.loads(Path(args.command_json).read_text())
    if isinstance(command, dict):
        command = command['command']
    if not isinstance(command, list) or not command or not all(isinstance(s, str) for s in command):
        ap.error('command JSON must be a nonempty string array or manifest with command')
    if args.require_cache_hits:
        if Path(command[0]).name != 'java':
            ap.error('cache counter check requires a Java command')
        command.insert(1, '-Djsgo.bench.verifyCache=true')
    root = Path(args.out).resolve()
    if root.exists():
        ap.error('refusing to overwrite '+str(root))
    root.mkdir(parents=True)
    pin = bench.read_pin()
    jar = Path(args.cache).resolve()/pin['CLOSURE_JAR']
    if bench.sha256_file(jar) not in pin['CLOSURE_JAR_SHA256'].split():
        ap.error('stock Java jar does not match SHA-256 pin')
    env = {k: v for k, v in os.environ.items() if k not in bench.GO_ENV_DROP+bench.JAVA_ENV_DROP}
    env['LC_ALL'] = 'C.UTF-8'
    source = '/** @type {number} */ var count=1; window.alert("é😀"+count);\n'
    cases = [
        ('help', ['--help'], source, b''),
        ('version', ['--version'], source, b''),
        ('whitespace', ['--compilation_level=WHITESPACE_ONLY', '--js', 'app.js'], source, b''),
        ('changed-source', ['--compilation_level=WHITESPACE_ONLY', '--js', 'app.js'], source.replace('count=1','count=7'), b''),
        ('simple-map', ['--compilation_level=SIMPLE', '--js', 'app.js', '--js_output_file=out.js', '--create_source_map=out.js.map'], source, b''),
        ('advanced-warning', ['--compilation_level=ADVANCED', '--jscomp_warning=checkTypes', '--js', 'app.js'], '/** @type {string} */ var incorrect=42;window.alert(incorrect);\n', b''),
        ('parse-error', ['--js', 'app.js'], 'function broken( {', b''),
        ('bad-flag', ['--invalid-worker-test-flag'], source, b''),
        ('stdin', ['--compilation_level=WHITESPACE_ONLY', '--js', '-'], 'ignored', source.encode()),
        ('json-streams', ['--compilation_level=WHITESPACE_ONLY', '--json_streams=BOTH'], 'ignored', json.dumps([{'path':'stream.js','src':source}]).encode()),
        ('custom-externs', ['--compilation_level=SIMPLE', '--env=CUSTOM', '--js', 'app.js'], 'var counter=1;counter++;', b''),
        ('default-after-custom', ['--compilation_level=WHITESPACE_ONLY', '--js', 'app.js'], source, b''),
        ('iso-8859-1', ['--compilation_level=WHITESPACE_ONLY', '--charset=ISO-8859-1', '--js', 'app.js'], source, b''),
        ('changed-language', ['--compilation_level=WHITESPACE_ONLY', '--language_in=ECMASCRIPT5', '--js', 'app.js'], 'var es5=1;window.alert(es5);', b''),
        ('repeat-default', ['--compilation_level=WHITESPACE_ONLY', '--js', 'app.js'], source, b''),
    ]
    counts = [int(count) for count in args.threads.split(',')]
    if any(count < 1 for count in counts):
        ap.error('thread counts must be positive')
    outputs = OUTPUT_FILES
    jobs = [(f'threads-{count}.{name}', thread_flags(flags, count), code, stdin)
            for count in counts for name, flags, code, stdin in cases]
    identities = {key: case_identity(flags, code, stdin) for key, flags, code, stdin in jobs}
    if len(identities) != len(jobs):
        ap.error('thread counts must be unique')
    jar_sha256 = bench.sha256_file(jar)
    reference_expected, reference_path = {}, None
    if args.oracle_reference:
        oracle, reference_expected, reference_path = load_oracle(args.oracle_reference, jar_sha256, identities)
    else:
        oracle = {'schema': ORACLE_SCHEMA, 'complete': False, 'stock_jar_sha256': jar_sha256,
                  'stderr_normalization': STDERR_NORMALIZATION,
                  'case_order': list(identities), 'corpus_sha256': object_digest(identities), 'cases': {}}
    summary = {'command': command, 'stock_jar_sha256': jar_sha256, 'cases': {},
               'worker_experiment_environment': experiment_environment(env),
               'oracle_mode': 'verified-reference' if reference_path else 'fresh-stock-java',
               'oracle_reference': reference_path, 'oracle_corpus_sha256': oracle['corpus_sha256']}
    bench.save(str(root/'oracle.json'), oracle)
    with (root/'worker.stderr').open('wb') as errors:
        process = subprocess.Popen(command, cwd=root, env=env, stdin=subprocess.PIPE,
                                   stdout=subprocess.PIPE, stderr=errors)
        selector = selectors.DefaultSelector()
        selector.register(process.stdout, selectors.EVENT_READ)
        try:
            request_id = 0
            for key, flags, code, stdin in jobs:
                (root/'app.js').write_bytes(code.encode('utf-8'))
                for output in outputs:
                    (root/output).unlink(missing_ok=True)
                if reference_path:
                    expected = reference_expected[key]
                else:
                    stock = subprocess.run(['java', '-jar', str(jar)]+flags, input=stdin,
                                           cwd=root, env=env, capture_output=True, timeout=args.timeout)
                    expected = {'exit': stock.returncode, 'stdout': stock.stdout, 'stderr': stock.stderr}
                    expected.update({f: (root/f).read_bytes() if (root/f).exists() else None for f in outputs})
                    # Store every raw stock byte, including missing files and
                    # actual exits; never reconstruct an oracle from booleans.
                    record = {'input': identities[key], 'expected': encode_expected(expected)}
                    decode_expected(record['expected'])
                    oracle['cases'][key] = record
                    bench.save(str(root/'oracle.json'), oracle)
                for output in outputs:
                    (root/output).unlink(missing_ok=True)
                process.stdin.write(json.dumps({'id': request_id, 'args': flags,
                    'stdin_base64': base64.b64encode(stdin).decode()}).encode()+b'\n')
                process.stdin.flush()
                if not selector.select(args.timeout):
                    raise RuntimeError('worker timeout: '+key)
                response = json.loads(process.stdout.readline())
                if response['id'] != request_id:
                    raise AssertionError('worker response id mismatch: '+key)
                actual = {'exit': response['exit'],
                          'stdout': base64.b64decode(response['stdout_base64'] or '', validate=True),
                          'stderr': base64.b64decode(response['stderr_base64'] or '', validate=True)}
                actual.update({f: (root/f).read_bytes() if (root/f).exists() else None for f in outputs})
                result = {}
                for artifact, want in expected.items():
                    got = actual[artifact]
                    if artifact == 'stderr':
                        got, want = normalize_frames(got), normalize_frames(want)
                    result[artifact] = got == want
                    if artifact in ('stdout', 'stderr'):
                        (root/f'{key}.stock.{artifact}').write_bytes(expected[artifact])
                        (root/f'{key}.worker.{artifact}').write_bytes(actual[artifact])
                summary['cases'][key] = result
                bench.save(str(root/'summary.json'), summary)
                if not all(result.values()):
                    raise AssertionError(f'worker differs from stock Java: {key}: {result}')
                print(key+' MATCH', flush=True)
                request_id += 1
            process.stdin.close()
            if process.wait(timeout=args.timeout) != 0:
                raise RuntimeError('worker failed on shutdown')
        finally:
            selector.close()
            if process.poll() is None:
                process.kill()
                process.wait()
            process.stdin.close()
            process.stdout.close()
    if args.require_cache_hits:
        data = (root/'worker.stderr').read_text()
        match = re.search(r'java-worker cache: hits=(\d+) misses=(\d+)', data)
        if not match or int(match[1]) == 0 or int(match[2]) == 0:
            raise AssertionError('Java cache reuse was not demonstrated: '+data)
        summary['cache_counters'] = {'hits': int(match[1]), 'misses': int(match[2])}
    if not reference_path:
        oracle['complete'] = True
        oracle['oracle_sha256'] = object_digest(oracle['cases'])
        bench.save(str(root/'oracle.json'), oracle)
    summary['oracle_sha256'] = oracle['oracle_sha256']
    summary['all_match'] = True
    bench.save(str(root/'summary.json'), summary)
    print(f'{len(summary["cases"])} mixed requests MATCH', flush=True)


if __name__ == '__main__':
    main()
