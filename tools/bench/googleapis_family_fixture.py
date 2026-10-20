#!/usr/bin/env python3
"""Prepare an entire genuine Google APIs family and exact transitive schemas.

No Java/Go Closure compilation or timing is performed. Upstream files remain
byte-identical; incompatible schemas/generator versions fail with diagnostics.
"""
import argparse
import base64
import hashlib
import json
from pathlib import Path
import re
import shutil
import subprocess
import sys
import tarfile
from urllib.request import Request, urlopen

sys.dont_write_bytecode = True
HERE = Path(__file__).resolve().parent
PIN = HERE/'generated/googleapis_bigquery_pin.json'
IMPORT = re.compile(r'^\s*import\s+(?:(?:public|weak)\s+)?"([^"]+)"\s*;', re.M)


def sha(path):
    return hashlib.sha256(Path(path).read_bytes()).hexdigest()


def save(path, value):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(value, indent=2)+'\n')


def safe(name):
    if Path(name).is_absolute() or '..' in Path(name).parts or '\\' in name:
        raise ValueError('unsafe upstream import path: '+name)
    return name


def fresh(path):
    if path.exists():
        raise ValueError('refusing to overwrite upstream/evidence files')
    path.mkdir(parents=True)


def inventory(root):
    return [{'path': p.relative_to(root).as_posix(), 'bytes': p.stat().st_size, 'sha256': sha(p)}
            for p in sorted(root.rglob('*')) if p.is_file()]


def fetch(a):
    pin = json.loads(PIN.read_text()); out = Path(a.out).resolve(); archive = Path(a.protobuf_archive).resolve()
    if sha(archive) != pin['protobuf_archive_sha256']:
        raise ValueError('matching genuine protobuf3.20.3 source archive required')
    fresh(out); schemas = out/'schemas'; schemas.mkdir()
    roots = {x['path']: x for x in pin['files']}; pending = list(sorted(roots)); records = {}; total = 0
    try:
        with tarfile.open(archive) as protobuf:
            while pending:
                name = safe(pending.pop(0))
                if name in records:
                    continue
                if len(records) >= 256:
                    raise ValueError('unexpectedly broad schema closure; inspect before expanding')
                if name.startswith('google/protobuf/'):
                    member = protobuf.getmember('protobuf-3.20.3/src/'+name)
                    if not member.isfile():
                        raise ValueError('not an ordinary upstream protobuf source')
                    data = protobuf.extractfile(member).read()
                    origin = 'protobuf3.20.3 source archive'
                else:
                    origin = 'https://raw.githubusercontent.com/googleapis/googleapis/'+pin['revision']+'/'+name
                    with urlopen(Request(origin, headers={'User-Agent': 'jsgo-real-api-family-fixture'}), timeout=30) as response:
                        data = response.read((4 << 20)+1)
                total += len(data)
                if len(data) > 4 << 20 or total > 20 << 20:
                    raise ValueError('schema source budget exceeded')
                blob = hashlib.sha1(b'blob '+str(len(data)).encode()+b'\0'+data).hexdigest()
                if name in roots and (len(data), blob) != (roots[name]['bytes'], roots[name]['git_blob_sha1']):
                    raise ValueError('family source differs from pinned published git object')
                path = schemas/name; path.parent.mkdir(parents=True, exist_ok=True); path.write_bytes(data)
                imports = IMPORT.findall(data.decode('utf-8'))
                records[name] = {'path': name, 'bytes': len(data), 'sha256': sha(path), 'git_blob_sha1': blob,
                                 'origin': origin, 'family_root': name in roots, 'imports': imports}
                pending.extend(safe(item) for item in imports if item not in records)
        with urlopen('https://raw.githubusercontent.com/googleapis/googleapis/'+pin['revision']+'/LICENSE', timeout=30) as response:
            (out/'GOOGLEAPIS-LICENSE').write_bytes(response.read(65536))
        with tarfile.open(archive) as protobuf:
            (out/'PROTOBUF-LICENSE').write_bytes(protobuf.extractfile('protobuf-3.20.3/LICENSE').read())
        save(out/'source-manifest.json', {'schema': 'jsgo-googleapis-schema-closure-v1', 'status': 'FETCHED; generator compatibility not yet established',
             'pin': pin, 'pin_sha256': sha(PIN), 'protobuf_archive': str(archive), 'roots': sorted(roots),
             'records': sorted(records.values(), key=lambda x: x['path']), 'provided_schema_count': len(records),
             'family_schema_count': len(roots), 'source_bytes': total, 'files_unchanged': True})
    except BaseException as error:
        save(out/'failure.json', {'type': type(error).__name__, 'message': str(error), 'partial_records': records})
        raise
    print(json.dumps({'family_schemas': len(roots), 'with_transitive_schemas': len(records), 'source_bytes': total}))


def generate(a):
    pin = json.loads(PIN.read_text()); source = Path(a.source).resolve(); out = Path(a.out).resolve(); executable = Path(a.protoc).resolve()
    if sha(executable) != pin['protoc_sha256']:
        raise ValueError('exact pinned genuine native protoc required')
    cap = Path('/sys/fs/cgroup'); quota, period = (cap/'cpu.max').read_text().split()
    if ((cap/'memory.max').read_text().strip() != '4294967296' or (cap/'memory.swap.max').read_text().strip() != '0'
        or quota == 'max' or int(quota) != 2*int(period) or (cap/'cpuset.cpus.effective').read_text().strip() not in ('0-1', '0,1')):
        raise ValueError('generation requires actual4GiB/swap0/CPU2 setup container')
    manifest = json.loads((source/'source-manifest.json').read_text())
    for row in manifest['records']:
        if sha(source/'schemas'/row['path']) != row['sha256']:
            raise ValueError('upstream schema changed')
    fresh(out); messages = out/'messages'; messages.mkdir()
    version = subprocess.run([str(executable), '--version'], capture_output=True, text=True)
    if version.returncode or version.stdout.strip() != 'libprotoc 3.20.3':
        raise ValueError('native protoc version mismatch')
    argv = [str(executable), '--proto_path='+str(source/'schemas'), '--js_out=import_style=closure,binary:'+str(messages),
            '--include_imports', '--include_source_info', '--descriptor_set_out='+str(out/'api-descriptors.pb')]
    argv += [row['path'] for row in manifest['records']]
    with (out/'stdout').open('wb') as stdout, (out/'stderr').open('wb') as stderr:
        completed = subprocess.run(argv, stdout=stdout, stderr=stderr)
    save(out/'generation.json', {'schema': 'jsgo-real-api-generation-v1', 'scope': 'Correctness preparation only; no performance claim',
         'exit': completed.returncode, 'argv': argv, 'protoc_version': version.stdout.strip(), 'protoc_sha256': sha(executable),
         'source_manifest_sha256': sha(source/'source-manifest.json'), 'generator_source_sha256': sha(__file__),
         'provided_schemas': len(manifest['records']), 'family_root_schemas': len(manifest['roots']), 'generated_files': inventory(messages)})
    if completed.returncode:
        raise RuntimeError('genuine generation failed; inspect retained incompatible-version diagnostics without editing upstream files')
    print(json.dumps({'provided_schemas': len(manifest['records']), 'generated_files': len(inventory(messages))}))


def wire_oracle(a):
    """Use independent native protobuf implementation for actual API payloads."""
    pin = json.loads(PIN.read_text()); source = Path(a.source).resolve(); out = Path(a.out).resolve(); executable = Path(a.protoc).resolve()
    cap = Path('/sys/fs/cgroup'); quota, period = (cap/'cpu.max').read_text().split()
    if ((cap/'memory.max').read_text().strip() != '4294967296' or (cap/'memory.swap.max').read_text().strip() != '0'
        or quota == 'max' or int(quota) != 2*int(period) or (cap/'cpuset.cpus.effective').read_text().strip() not in ('0-1', '0,1')):
        raise ValueError('native reference encoding requires actual capped setup process')
    if sha(executable) != pin['protoc_sha256']:
        raise ValueError('unchanged native protobuf executable required')
    manifest = json.loads((source/'source-manifest.json').read_text())
    for row in manifest['records']:
        if sha(source/'schemas'/row['path']) != row['sha256']:
            raise ValueError('schema changed before independent wire encoding')
    fresh(out); proofs = {}
    for name in ('table', 'job'):
        payload = HERE/'generated'/('bigquery_'+name+'.textproto')
        typename = 'google.cloud.bigquery.v2.'+name.title()
        common = [str(executable), '--proto_path='+str(source/'schemas')]
        schema = 'google/cloud/bigquery/v2/'+name+'.proto'
        encode = common+['--encode='+typename, schema]
        raw = subprocess.run(encode, input=payload.read_bytes(), capture_output=True)
        (out/(name+'.encode.stderr')).write_bytes(raw.stderr)
        (out/(name+'.pb')).write_bytes(raw.stdout)
        if raw.returncode:
            raise RuntimeError('native reference encoding failed: '+name)
        decode = common+['--decode='+typename, schema]
        decoded = subprocess.run(decode, input=raw.stdout, capture_output=True)
        (out/(name+'.decoded.textproto')).write_bytes(decoded.stdout)
        (out/(name+'.decode.stderr')).write_bytes(decoded.stderr)
        if decoded.returncode:
            raise RuntimeError('native reference decoding failed: '+name)
        proofs[name] = {'encode_argv': encode, 'decode_argv': decode, 'payload_path': str(payload),
                        'payload_sha256': sha(payload), 'wire_sha256': sha(out/(name+'.pb')),
                        'wire_base64': base64.b64encode(raw.stdout).decode(), 'wire_bytes': len(raw.stdout),
                        'decoded_sha256': sha(out/(name+'.decoded.textproto'))}
    save(out/'wire-oracle.json', {'schema': 'jsgo-real-api-native-wire-oracle-v1', 'status': 'PASS',
         'scope': 'Native protobuf correctness preparation; no performance claim', 'protoc_sha256': sha(executable),
         'source_manifest_sha256': sha(source/'source-manifest.json'), 'payloads': proofs})


def finalize(a):
    source, generation, wire, runtime, out = [Path(getattr(a, name)).resolve() for name in ('source', 'generation', 'wire', 'runtime', 'out')]
    sm = json.loads((source/'source-manifest.json').read_text()); gm = json.loads((generation/'generation.json').read_text()); wm = json.loads((wire/'wire-oracle.json').read_text())
    if (gm['exit'] or gm['source_manifest_sha256'] != sha(source/'source-manifest.json') or
        gm['generated_files'] != inventory(generation/'messages') or wm['source_manifest_sha256'] != sha(source/'source-manifest.json') or wm['status'] != 'PASS'):
        raise ValueError('unchanged genuine generation and native wire oracle required')
    for row in sm['records']:
        if sha(source/'schemas'/row['path']) != row['sha256']:
            raise ValueError('upstream schema changed after generation')
    rm = json.loads((runtime/'fixture-manifest.json').read_text()); original = {row['relative']: row for row in rm['files']}
    expected_runtime = {name for name in original if name.startswith('upstream/js/') and name.endswith('.js')}
    if {p.relative_to(runtime).as_posix() for p in (runtime/'upstream/js').rglob('*.js')} != expected_runtime:
        raise ValueError('matching genuine runtime file set changed')
    for path in (runtime/'upstream/js').rglob('*.js'):
        if sha(path) != original[path.relative_to(runtime).as_posix()]['sha256']:
            raise ValueError('matching genuine protobuf JS runtime changed')
    fresh(out); shutil.copytree(source/'schemas', out/'schemas'); shutil.copytree(generation/'messages', out/'messages')
    shutil.copytree(runtime/'upstream/js', out/'upstream/js')
    for name in ('GOOGLEAPIS-LICENSE', 'PROTOBUF-LICENSE'):
        shutil.copyfile(source/name, out/name)
    provides = set()
    for path in (out/'messages').glob('*.js'):
        provides.update(re.findall(r"goog\.provide\('([^']+)'\)", path.read_text()))
    family = sorted(name for name in provides if name.startswith('proto.google.cloud.bigquery.v2.'))
    if not family:
        raise ValueError('missing actual generated BigQuery namespace')
    requires = sorted(set(family) | {'proto.google.protobuf.StringValue', 'proto.google.protobuf.BoolValue'})
    entry = "/** Complete real BigQueryV2 type catalog; no synthetic messages. */\ngoog.provide('enterprise.bigqueryV2');\n"
    entry += ''.join("goog.require('"+name+"');\n" for name in requires)
    entry += "globalThis['jsgoBigQueryV2Types'] = {\n"+''.join(json.dumps(name)+': '+name+',\n' for name in family)+'};\n'
    template = (HERE/'generated/bigquery_exercise.js.txt').read_text()
    for name in ('table', 'job'):
        raw = (wire/(name+'.pb')).read_bytes()
        if hashlib.sha256(raw).hexdigest() != wm['payloads'][name]['wire_sha256']:
            raise ValueError('native reference wire changed')
        template = template.replace('__'+name.upper()+'_WIRE__', json.dumps(list(raw)))
    template = template.replace('__CATALOG_COUNT__', str(len(family)))
    (out/'bigquery_entry.js').write_text(entry+template)
    save(out/'provenance.json', {'source': sm, 'generation': gm, 'native_wire_oracle': wm, 'runtime_fixture_manifest_sha256': sha(runtime/'fixture-manifest.json')})
    save(out/'fixture-manifest.json', {'schema': 'jsgo-bigquery-family-fixture-v1', 'status': 'GENERATED; Closure compiler and compiled-runtime gates pending',
         'googleapis_revision': sm['pin']['revision'], 'family_schemas': len(sm['roots']), 'all_transitive_schemas': len(sm['records']),
         'generated_js_files': len(gm['generated_files']), 'generated_js_bytes': sum(x['bytes'] for x in gm['generated_files']),
         'actual_api_type_catalog': family, 'owned_catalog_and_semantic_entry_bytes': (out/'bigquery_entry.js').stat().st_size,
         'fixture_builder_sha256': sha(__file__), 'semantic_template_sha256': sha(HERE/'generated/bigquery_exercise.js.txt'),
         'unchanged_upstream_files': True, 'synthetic_schemas': False, 'files': inventory(out)})
    print(json.dumps({'generated_messages': len(gm['generated_files']), 'actual_exported_api_types': len(family)}))


def specs(a):
    import bench
    from enterprise_specs import embedded_externs, replace_threads
    fixture, out = Path(a.fixture).resolve(), Path(a.out).resolve(); record = json.loads((fixture/'fixture-manifest.json').read_text())
    if a.threads < 1:
        raise ValueError('positive compiler thread count required')
    actual = [x for x in inventory(fixture) if x['path'] != 'fixture-manifest.json']
    if actual != record['files']:
        raise ValueError('frozen generated fixture changed')
    cache = Path(a.cache).resolve(); inputs = cache/'bench-inputs'; closure = inputs/'closure-library@v20230802'
    closure_files = sorted(p for prefix in ('closure/goog', 'third_party/closure/goog') for p in (closure/prefix).rglob('*.js') if not p.name.endswith('_test.js'))
    fresh(out); sources = []; argv = ['--compilation_level=ADVANCED', '--jscomp_warning=checkTypes', '--jscomp_warning=checkVars', '--language_out=ECMASCRIPT_2017', '--dependency_mode=PRUNE', '--entry_point=goog:enterprise.bigqueryV2']
    for link, root, files in [('inputs', inputs, closure_files), ('generated', fixture, sorted(fixture.rglob('*.js')))]:
        for path in files:
            rel = path.relative_to(root).as_posix(); argv += ['--js', link+'/'+rel]
            sources.append({'flag': '--js', 'link': link, 'relative': rel, 'bytes': path.stat().st_size, 'sha256': sha(path)})
    source_count = len(sources)
    for path in sorted(fixture.rglob('*.proto')):
        sources.append({'flag': 'upstream-schema-provenance', 'link': 'generated', 'relative': path.relative_to(fixture).as_posix(), 'bytes': path.stat().st_size, 'sha256': sha(path)})
    argv = replace_threads(argv+['--js_output_file=out.js', '--create_source_map=out.js.map', '--output_manifest=inputs.manifest'], a.threads)
    result = {'schema': 'jsgo-enterprise-cli-v1', 'name': 'generated-bigquery-v2-advanced', 'argv': argv, 'links': {'inputs': str(inputs), 'generated': str(fixture)},
              'inputs': sources, 'artifacts': ['out.js', 'out.js.map', 'inputs.manifest'], 'fixture_manifest_sha256': sha(fixture/'fixture-manifest.json'),
              'semantic_check': {'script': str(HERE/'generated/verify_semantics.js'), 'entry': 'jsgoEnterpriseBigQueryV2', 'expected': 'benchmark-project|analytics|invoices|2|INT64|42|false|Revenue € 🧾|query-invoices-001'},
              'same_work': {'mode': 'ADVANCED/PRUNE/checkTypes/checkVars, complete genuine BigQueryV2 type catalog and native protobuf wire oracle', 'compiler_threads': a.threads,
                            'source_count': source_count, 'source_bytes': sum(x['bytes'] for x in sources if x['flag'] == '--js'), 'builtin_externs': embedded_externs(a.jar, Path(bench.REPO)),
                            'provided_family_schemas': record['family_schemas'], 'provided_transitive_schemas': record['all_transitive_schemas'], 'generated_js_files': record['generated_js_files'],
                            'actual_manifest_work': 'PENDING actual stock Java and Go correctness gate; supplied Closure candidates may be normally pruned',
                            'explicit_externs': [], 'synthetic_schemas': False, 'limitations': 'Coherent real generated API model library, not an enterprise application or a network service invocation.'}}
    save(out/(result['name']+'.json'), result); print(json.dumps({'provided_js_sources': source_count, 'provided_schema_auxiliary_inputs': len(sources)-source_count, 'js_bytes': result['same_work']['source_bytes']}))


def main():
    p = argparse.ArgumentParser(description=__doc__); sub = p.add_subparsers(dest='operation', required=True)
    f = sub.add_parser('fetch'); f.add_argument('--out', required=True); f.add_argument('--protobuf-archive', required=True); f.set_defaults(run=fetch)
    g = sub.add_parser('generate'); g.add_argument('--source', required=True); g.add_argument('--out', required=True); g.add_argument('--protoc', required=True); g.set_defaults(run=generate)
    w = sub.add_parser('wire-oracle'); w.add_argument('--source', required=True); w.add_argument('--out', required=True); w.add_argument('--protoc', required=True); w.set_defaults(run=wire_oracle)
    f = sub.add_parser('finalize')
    for name in ('source', 'generation', 'wire', 'runtime', 'out'): f.add_argument('--'+name, required=True)
    f.set_defaults(run=finalize)
    s = sub.add_parser('specs')
    for name in ('fixture', 'out', 'cache', 'jar'): s.add_argument('--'+name, required=True)
    s.add_argument('--threads', type=int, default=1); s.set_defaults(run=specs)
    a = p.parse_args(); a.run(a)


if __name__ == '__main__':
    main()
