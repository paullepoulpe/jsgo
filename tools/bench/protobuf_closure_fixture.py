#!/usr/bin/env python3
"""Prepare real protoc Closure-style inputs; full compiler gates are separate.

Run generate ONLY through the existing capped setup launcher. Native protoc
version and generation then share its real 4 GiB / CPU2 process-tree budget.
No upstream source/generated file is edited or committed. The only owned JS is
a semantic harness. There are no synthetic schemas, fake externs or stubs.
"""
import argparse
import hashlib
import json
import os
from pathlib import Path
import shutil
import subprocess
import sys
import tarfile
from urllib.request import Request, urlopen

sys.dont_write_bytecode = True
HERE = Path(__file__).resolve().parent
PINS = HERE/'generated/protobuf_closure_assets.json'
SCHEMAS = ['google/protobuf/'+name+'.proto' for name in ('descriptor', 'struct', 'timestamp')]
RUNTIME = ['js/'+name+'.js' for name in ('debug', 'map', 'message')]+[
    'js/binary/'+name+'.js' for name in ('arith', 'constants', 'decoder', 'encoder', 'reader', 'utils', 'writer')]


def sha(path):
    return hashlib.sha256(Path(path).read_bytes()).hexdigest()


def save(path, value):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(value, indent=2)+'\n')


def inventory(root):
    return [{'relative': p.relative_to(root).as_posix(), 'bytes': p.stat().st_size, 'sha256': sha(p)}
            for p in sorted(root.rglob('*')) if p.is_file() and p.name != 'fixture-manifest.json']


def validate_assets(cache):
    records = json.loads(PINS.read_text())['assets']
    for row in records:
        path = cache/row['name']
        if path.stat().st_size != row['bytes'] or sha(path) != row['sha256']:
            raise ValueError('pinned upstream asset mismatch: '+str(path))
        if row.get('upstream_sha1') and hashlib.sha1(path.read_bytes()).hexdigest() != row['upstream_sha1']:
            raise ValueError('Maven published SHA1 mismatch')
    return records


def fresh(path):
    if path.exists():
        raise ValueError('refusing to overwrite fixture evidence: '+str(path))
    path.mkdir(parents=True)


def fetch(args):
    cache = Path(args.cache).resolve(); cache.mkdir(parents=True, exist_ok=True)
    for row in json.loads(PINS.read_text())['assets']:
        path = cache/row['name']
        if not path.exists():
            partial = path.with_name(path.name+'.partial-'+str(os.getpid()))
            try:
                with urlopen(Request(row['url'], headers={'User-Agent': 'jsgo-protobuf-closure-fixture'}), timeout=60) as response:
                    data = response.read(row['bytes']+1)
                partial.write_bytes(data)
                if len(data) != row['bytes'] or sha(partial) != row['sha256']:
                    raise ValueError('download differs from exact upstream pin')
                os.link(partial, path)
            finally:
                partial.unlink(missing_ok=True)
    records = validate_assets(cache)
    (cache/'protoc-3.20.3-linux-x86_64.exe').chmod(0o755)
    save(cache/'pinned-download-manifest.json', {'pins_sha256': sha(PINS), 'assets': records})


def stage(args):
    cache, out = Path(args.cache).resolve(), Path(args.out).resolve()
    records = validate_assets(cache); fresh(out)
    selected = RUNTIME+['src/'+name for name in SCHEMAS]+['LICENSE', 'js/package.json']
    with tarfile.open(cache/'protobuf-v3.20.3.tar.gz') as archive:
        for name in selected:
            member = archive.getmember('protobuf-3.20.3/'+name)
            if not member.isfile():
                raise ValueError('expected an ordinary upstream source file: '+name)
            target = out/'upstream'/name
            target.parent.mkdir(parents=True, exist_ok=True)
            target.write_bytes(archive.extractfile(member).read())
    shutil.copyfile(HERE/'generated/protobuf_closure_entry.js', out/'protobuf_closure_entry.js')
    save(out/'fixture-manifest.json', {'schema': 'jsgo-protobuf-closure-stage-v1', 'status': 'STAGED; not compiled',
         'assets': records, 'pins_sha256': sha(PINS), 'upstream': json.loads(PINS.read_text()),
         'synthetic_schemas': False, 'files': inventory(out)})


def check_stage(stage):
    record = json.loads((stage/'fixture-manifest.json').read_text())
    if record['files'] != inventory(stage):
        raise ValueError('staged upstream sources changed')
    return record


def generate(args):
    # This is the setup child, not a compiler benchmark or Java invocation.
    cap = Path('/sys/fs/cgroup')
    quota, period = (cap/'cpu.max').read_text().split()
    if ((cap/'memory.max').read_text().strip() != '4294967296' or
        (cap/'memory.swap.max').read_text().strip() != '0' or quota == 'max' or int(quota) != 2*int(period) or
        (cap/'cpuset.cpus.effective').read_text().strip() not in ('0-1', '0,1')):
        raise ValueError('generation requires the real 4GiB/zero-swap/CPU2 capped setup launcher')
    cache, stage, out = Path(args.cache).resolve(), Path(args.stage).resolve(), Path(args.out).resolve()
    assets = validate_assets(cache); staged = check_stage(stage); fresh(out)
    executable = cache/'protoc-3.20.3-linux-x86_64.exe'
    version_argv = [str(executable), '--version']
    version = subprocess.run(version_argv, capture_output=True, text=True)
    (out/'version.stdout').write_text(version.stdout); (out/'version.stderr').write_text(version.stderr)
    if version.returncode != 0 or version.stdout.strip() != 'libprotoc 3.20.3':
        raise ValueError('native protoc version differs from pinned 3.20.3')
    messages = out/'messages'; messages.mkdir()
    argv = [str(executable), '--proto_path='+str(stage/'upstream/src'),
            '--js_out=import_style=closure,binary:'+str(messages)]+SCHEMAS
    with (out/'generator.stdout').open('wb') as stdout, (out/'generator.stderr').open('wb') as stderr:
        result = subprocess.run(argv, stdout=stdout, stderr=stderr)
    record = {'schema': 'jsgo-protobuf-closure-generation-v1', 'argv': argv, 'version_argv': version_argv,
              'version': version.stdout.strip(), 'exit': result.returncode, 'protoc_sha256': sha(executable),
              'assets': assets, 'stage_manifest_sha256': sha(stage/'fixture-manifest.json'),
              'stage_files': staged['files'], 'generator_script_sha256': sha(__file__), 'files': inventory(messages)}
    save(out/'generation.json', record)
    if result.returncode != 0 or len(record['files']) != 3:
        raise ValueError('native upstream JS generation failed; preserve diagnostics')


def finalize(args):
    stage, generation, out = Path(args.stage).resolve(), Path(args.generation).resolve(), Path(args.out).resolve()
    staged = check_stage(stage); generated = json.loads((generation/'generation.json').read_text())
    if (generated['exit'] != 0 or generated['version'] != 'libprotoc 3.20.3' or
        generated['stage_manifest_sha256'] != sha(stage/'fixture-manifest.json') or
        generated['files'] != inventory(generation/'messages')):
        raise ValueError('native generated artifacts/stage identity changed')
    fresh(out)
    shutil.copytree(stage/'upstream', out/'upstream')
    shutil.copytree(generation/'messages', out/'messages')
    shutil.copyfile(stage/'protobuf_closure_entry.js', out/'protobuf_closure_entry.js')
    save(out/'fixture-manifest.json', {'schema': 'jsgo-protobuf-closure-fixture-v1',
         'status': 'GENERATED; Java/Go artifact and independent runtime gates pending',
         'assets': staged['assets'], 'upstream': staged['upstream'], 'generator': generated,
         'generation_evidence': str(generation), 'synthetic_schemas': False, 'files': inventory(out)})


def specs(args):
    import bench
    from enterprise_specs import embedded_externs, replace_threads
    fixture, out = Path(args.fixture).resolve(), Path(args.out).resolve()
    record = check_stage(fixture)
    if args.threads < 1 or not record['status'].startswith('GENERATED'):
        raise ValueError('need a generated fixture and positive thread count')
    fresh(out)
    inputs = Path(args.cache or bench.default_cache()).resolve()/'bench-inputs'
    closure = inputs/'closure-library@v20230802'
    closure_js = sorted(p for prefix in ('closure/goog', 'third_party/closure/goog')
                        for p in (closure/prefix).rglob('*.js') if not p.name.endswith('_test.js'))
    if not closure_js or not (closure/'closure/goog/base.js').is_file():
        raise ValueError('existing real pinned Closure library required')
    argv = ['--compilation_level=ADVANCED', '--jscomp_warning=checkTypes', '--jscomp_warning=checkVars',
            '--language_out=ECMASCRIPT_2017', '--dependency_mode=PRUNE', '--entry_point=goog:enterprise.protobufClosure']
    sources = []
    paths = [('inputs', inputs, closure_js), ('generated', fixture,
             sorted(p for p in fixture.rglob('*.js') if p.is_file()))]
    for link, root, files in paths:
        for path in files:
            relative = path.relative_to(root).as_posix(); argv += ['--js', link+'/'+relative]
            sources.append({'flag': '--js', 'link': link, 'relative': relative,
                            'bytes': path.stat().st_size, 'sha256': sha(path)})
    for path in sorted((fixture/'upstream/src').rglob('*.proto')):
        sources.append({'flag': 'upstream-schema-provenance', 'link': 'generated',
                        'relative': path.relative_to(fixture).as_posix(), 'bytes': path.stat().st_size, 'sha256': sha(path)})
    argv = replace_threads(argv+['--js_output_file=out.js', '--create_source_map=out.js.map', '--output_manifest=inputs.manifest'], args.threads)
    spec = {'schema': 'jsgo-enterprise-cli-v1', 'name': 'generated-protobuf-closure-advanced', 'argv': argv,
            'links': {'inputs': str(inputs), 'generated': str(fixture)}, 'inputs': sources,
            'artifacts': ['out.js', 'out.js.map', 'inputs.manifest'], 'fixture_manifest_sha256': sha(fixture/'fixture-manifest.json'),
            'semantic_check': {'script': str(HERE/'generated/verify_semantics.js'), 'entry': 'jsgoEnterpriseProtobufClosure',
                               'expected': 'enterprise.proto|Invoice|total|1700000000|123456789|Ada',
                               'required': 'run independently on stock Java and Go compiled outputs before performance claims'},
            'same_work': {'mode': 'ADVANCED Closure-style generated protobuf; PRUNE, checkTypes/checkVars, maps and manifest',
                          'compiler_threads': args.threads, 'builtin_externs': embedded_externs(args.jar, Path(bench.REPO)),
                          'explicit_externs': [], 'source_count': len(sources)-3,
                          'source_bytes': sum(s['bytes'] for s in sources if s['flag'] == '--js'),
                          'generated_family': 'protobuf-closure', 'generator_version': record['generator']['version'],
                          'synthetic_schemas': False, 'fixture_status': record['status'],
                          'limitations': 'Real upstream generated well-known/descriptor messages and runtime; benchmark-owned semantic harness, not a production application.'}}
    save(out/'generated-protobuf-closure-advanced.json', spec)
    print(json.dumps({'spec': str(out/'generated-protobuf-closure-advanced.json'), 'source_count': len(sources)-3,
                      'source_bytes': spec['same_work']['source_bytes'], 'status': record['status']}))


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    sub = parser.add_subparsers(dest='operation', required=True)
    for name, function, options in [('fetch', fetch, ('cache',)), ('stage', stage, ('cache', 'out')),
        ('generate', generate, ('cache', 'stage', 'out')), ('finalize', finalize, ('stage', 'generation', 'out')),
        ('specs', specs, ('fixture', 'out', 'jar'))]:
        command = sub.add_parser(name)
        for option in options: command.add_argument('--'+option, required=True)
        if name == 'specs': command.add_argument('--cache'); command.add_argument('--threads', type=int, default=2)
        command.set_defaults(run=function)
    args = parser.parse_args(); args.run(args)


if __name__ == '__main__':
    main()
