#!/usr/bin/env python3
"""Prepare pinned, actual generated JS inputs; compiler trials are separate.

fetch downloads exact external assets. stage extracts them and, for Soy only,
explicitly runs the pinned upstream generator outside benchmark measurements.
specs freezes complete compiler work for paired_cli.py. verify executes an
independent semantic check of one compiler's output outside measured time.
No source, generated artifact, downloaded jar or runtime stub is committed.
"""
import argparse
import base64
import hashlib
import json
import os
from pathlib import Path
import shutil
import subprocess
import sys
import tarfile
import tempfile
from urllib.request import Request, urlopen
import zipfile

sys.dont_write_bytecode = True
HERE = Path(__file__).resolve().parent
sys.path.insert(0, str(HERE))
import bench
from enterprise_specs import embedded_externs, replace_threads

ASSETS = HERE/'generated/assets.json'
FAMILIES = ('j2cl', 'protobuf', 'soy')
SEMANTICS = {
    'j2cl': ('jsgoEnterpriseJ2cl', 'J2CL:13,11;remaining=1'),
    'protobuf': ('jsgoEnterpriseProtobuf', 'enterprise.proto|Invoice|total|1700000000|123456789|Ada'),
    'soy': ('jsgoEnterpriseSoy', 'Hello &lt;Ada&gt;!|Hello world!'),
}


def sha(path):
    return hashlib.sha256(Path(path).read_bytes()).hexdigest()


def save(path, value):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(value, indent=2)+'\n')


def selected(text):
    values = text.split(',')
    if not values or set(values)-set(FAMILIES) or len(set(values)) != len(values):
        raise ValueError('families must be distinct j2cl,protobuf,soy selections')
    return values


def pins(families):
    return [x for x in json.loads(ASSETS.read_text())['assets'] if x['family'] in families]


def validate(path, asset):
    data = path.read_bytes()
    if 'integrity' in asset:
        algorithm, expected = asset['integrity'].split('-', 1)
        actual = base64.b64encode(hashlib.new(algorithm, data).digest()).decode()
    else:
        algorithm = 'sha256' if 'sha256' in asset else 'sha1'
        expected, actual = asset[algorithm], hashlib.new(algorithm, data).hexdigest()
    if actual != expected:
        raise ValueError('pinned artifact mismatch: '+str(path))
    return {'name': asset['name'], 'url': asset['url'], 'bytes': len(data),
            'sha256': hashlib.sha256(data).hexdigest(), 'upstream_pin': algorithm+'-'+expected,
            'license': asset['license']}


def fetch(args):
    cache = Path(args.cache).resolve()
    cache.mkdir(parents=True, exist_ok=True)
    records = []
    for asset in pins(selected(args.families)):
        path = cache/asset['name']
        if not path.exists():
            # Unique partial files and atomic non-overwriting publication permit
            # a shared download cache without replacing an existing pinned file.
            with tempfile.NamedTemporaryFile(dir=cache, prefix=asset['name']+'.', delete=False) as stream:
                partial = Path(stream.name)
                try:
                    with urlopen(Request(asset['url'], headers={'User-Agent': 'jsgo-generated-fixtures'}), timeout=60) as response:
                        total = 0
                        while block := response.read(1 << 20):
                            total += len(block)
                            if total > (64 << 20):
                                raise ValueError('asset exceeds declared setup safety bound')
                            stream.write(block)
                    stream.flush()
                    validate(partial, asset)
                    try:
                        os.link(partial, path)
                    except FileExistsError:
                        pass  # Validate the winner rather than overwrite it.
                finally:
                    partial.unlink(missing_ok=True)
        records.append(validate(path, asset))
    save(cache/'download-manifest.json', {'assets_pin_sha256': sha(ASSETS), 'records': records})
    print(json.dumps({'downloaded_or_verified': len(records), 'cache': str(cache)}))


def member_target(root, name):
    rel = Path(name)
    if rel.is_absolute() or '..' in rel.parts or '\\' in name:
        raise ValueError('unsafe archive member: '+name)
    target = root/rel
    target.parent.mkdir(parents=True, exist_ok=True)
    if target.exists():
        raise ValueError('duplicate archive member: '+name)
    return target


def unzip(archive, root):
    with zipfile.ZipFile(archive) as source:
        for item in source.infolist():
            if item.is_dir():
                continue
            if (item.external_attr >> 16) & 0o170000 == 0o120000:
                raise ValueError('archive symlink is not an input file')
            member_target(root, item.filename).write_bytes(source.read(item))


def unpack_package(archive, root):
    with tarfile.open(archive, 'r:gz') as source:
        for item in source:
            if item.isdir():
                continue
            if not item.isfile() or not item.name.startswith('package/'):
                raise ValueError('unexpected npm archive entry: '+item.name)
            target = member_target(root, item.name[len('package/'):])
            with source.extractfile(item) as stream:
                target.write_bytes(stream.read())


def inventory(root):
    return [{'relative': p.relative_to(root).as_posix(), 'bytes': p.stat().st_size, 'sha256': sha(p)}
            for p in sorted(root.rglob('*')) if p.is_file() and p.name != 'fixture-manifest.json']


def stage(args):
    families = selected(args.families)
    cache, out = Path(args.cache).resolve(), Path(args.out).resolve()
    if out.exists():
        raise ValueError('refusing to replace an existing fixture')
    records = [validate(cache/x['name'], x) for x in pins(families)]
    out.mkdir(parents=True)
    manifest = {'schema': 'jsgo-generated-fixture-v1', 'families': families,
                'asset_manifest_sha256': sha(ASSETS), 'assets': records,
                'generator': None, 'status': 'PREPARING; no compiler or semantic gates yet'}
    try:
        if 'j2cl' in families:
            unzip(cache/'j2cl-jre.zip', out/'j2cl/jre')
            unzip(cache/'j2cl-bootstrap.zip', out/'j2cl/bootstrap')
        if 'protobuf' in families:
            unpack_package(cache/'google-protobuf.tgz', out/'node_modules/google-protobuf')
        if 'soy' in families:
            if not args.java:
                raise ValueError('Soy staging requires explicit --java; generation is setup, never a benchmark trial')
            source, generated = out/'soy/source', out/'soy/generated'
            source.mkdir(parents=True); generated.mkdir()
            for asset in pins(['soy']):
                if asset['name'].endswith('.soy'):
                    shutil.copyfile(cache/asset['name'], source/asset['name'])
            shutil.copyfile(cache/'soyutils_usegoog.js', out/'soy/soyutils_usegoog.js')
            # The pinned compiler's default namespace output matches its SAME
            # release usegoog runtime. Retain actual invocation/generated bytes.
            command = [str(Path(args.java).resolve()), '-Xmx512m', '-jar', str(cache/'soy-compiler.jar'),
                       '--outputPathFormat', str(generated/'{INPUT_FILE_NAME_NO_EXT}.js')]
            command += ['--srcs', ','.join(str(p) for p in sorted(source.glob('*.soy')))]
            with (out/'soy/generator.stdout').open('wb') as stdout, (out/'soy/generator.stderr').open('wb') as stderr:
                result = subprocess.run(command, stdout=stdout, stderr=stderr)
            manifest['generator'] = {'argv': command, 'exit': result.returncode,
                                     'java_binary_sha256': sha(Path(args.java).resolve())}
            if result.returncode or len(list(generated.glob('*.js'))) != len(list(source.glob('*.soy'))):
                raise ValueError('pinned Soy generator failed; preserve fixture diagnostics, do not hand-edit output')
        for family in families:
            shutil.copyfile(HERE/'generated'/f'{family}_entry.js', out/f'{family}_entry.js')
        manifest.update({'status': 'STAGED; compiler and semantic compatibility not established', 'files': inventory(out)})
        save(out/'fixture-manifest.json', manifest)
    except BaseException as error:
        manifest['failure'] = str(error)
        save(out/'fixture-manifest.json', manifest)
        raise
    print(json.dumps({'families': families, 'out': str(out), 'files': len(manifest['files'])}))


def specs(args):
    fixture, out = Path(args.fixture).resolve(), Path(args.out).resolve()
    manifest = json.loads((fixture/'fixture-manifest.json').read_text())
    if not manifest['status'].startswith('STAGED') or inventory(fixture) != manifest['files']:
        raise ValueError('fixture incomplete or changed after staging')
    if out.exists() or args.threads < 1:
        raise ValueError('need new output directory and positive thread count')
    out.mkdir(parents=True)
    cache = Path(args.cache or bench.default_cache()).resolve()
    links = {'inputs': str(cache/'bench-inputs'), 'generated': str(fixture)}
    closure = Path(links['inputs'])/'closure-library@v20230802'
    closure_js = sorted(p for prefix in ('closure/goog', 'third_party/closure/goog')
                        for p in (closure/prefix).rglob('*.js') if not p.name.endswith('_test.js'))
    externs = embedded_externs(args.jar, Path(bench.REPO))
    for family in manifest['families']:
        flags = ['--compilation_level=ADVANCED', '--jscomp_warning=checkTypes', '--jscomp_warning=checkVars',
                 '--language_out=ECMASCRIPT_2017', '--dependency_mode=PRUNE']
        if family == 'protobuf':
            files = sorted(p for p in (fixture/'node_modules/google-protobuf').rglob('*.js') if p.is_file())
            flags += ['--process_common_js_modules', '--module_resolution=NODE', '--entry_point=generated/protobuf_entry.js']
        else:
            files = sorted(p for p in (fixture/family).rglob('*.js') if p.is_file())
            flags += ['--entry_point=goog:enterprise.'+family]
            if family == 'j2cl':
                flags += ['--j2cl_pass=ON']
        sources = []
        for link, root, paths in [('inputs', Path(links['inputs']), closure_js if family == 'soy' else []),
                                  ('generated', fixture, files+[fixture/f'{family}_entry.js'])]:
            for path in paths:
                relative = path.relative_to(root).as_posix()
                flags += ['--js', link+'/'+relative]
                sources.append({'flag': '--js', 'link': link, 'relative': relative, 'bytes': path.stat().st_size, 'sha256': sha(path)})
        # Generated JRE archives carry original input maps/source sidecars.
        # Freeze these too if the compiler consults sourceMappingURL metadata.
        auxiliary = [p for p in (fixture/'j2cl').rglob('*') if p.is_file() and not p.name.endswith('.js')] if family == 'j2cl' else []
        for path in sorted(auxiliary):
            sources.append({'flag': 'input-source-map-or-sidecar', 'link': 'generated',
                            'relative': path.relative_to(fixture).as_posix(),
                            'bytes': path.stat().st_size, 'sha256': sha(path)})
        auxiliary_count = len(auxiliary)
        if family == 'protobuf':
            package = fixture/'node_modules/google-protobuf/package.json'
            sources.append({'flag': 'package-metadata', 'link': 'generated',
                            'relative': package.relative_to(fixture).as_posix(),
                            'bytes': package.stat().st_size, 'sha256': sha(package)})
            auxiliary_count += 1
        argv = replace_threads(flags+['--js_output_file=out.js', '--create_source_map=out.js.map', '--output_manifest=inputs.manifest'], args.threads)
        entry, expected = SEMANTICS[family]
        spec = {'schema': 'jsgo-enterprise-cli-v1', 'name': 'generated-'+family, 'argv': argv,
                'links': links, 'inputs': sources, 'artifacts': ['out.js', 'out.js.map', 'inputs.manifest'],
                'fixture_manifest_sha256': sha(fixture/'fixture-manifest.json'),
                'semantic_check': {'script': str(HERE/'generated/verify_semantics.js'), 'entry': entry, 'expected': expected,
                                   'required': 'run independently on stock Java and Go outputs before a valid performance claim'},
                'same_work': {'mode': 'ADVANCED, checkTypes and checkVars warnings, source map, manifest, explicit PRUNE entry',
                              'compiler_threads': args.threads, 'builtin_externs': externs, 'explicit_externs': [],
                              'source_count': len(sources)-auxiliary_count, 'auxiliary_input_count': auxiliary_count,
                              'source_bytes': sum(s['bytes'] for s in sources if s['flag'] == '--js'),
                              'generated_family': family, 'fixture_status': manifest['status'],
                              'limitations': 'real generated runtime/example fixture; not claimed to be a production enterprise application',
                              'j2cl_runtime': 'matching release JRE + 7 Closure bootstrap JS files; no duplicate full Closure tree' if family == 'j2cl' else None}}
        save(out/(spec['name']+'.json'), spec)
        print(spec['name'], spec['same_work']['source_count'], 'JS sources,', auxiliary_count, 'auxiliary inputs,', spec['same_work']['source_bytes'], 'JS bytes; UNVALIDATED')


def verify(args):
    spec = json.loads(Path(args.spec).read_text())
    semantic = spec['semantic_check']
    command = [str(Path(args.node).resolve()), semantic['script'], str(Path(args.compiled_js).resolve()), semantic['entry'], semantic['expected']]
    result = subprocess.run(command, capture_output=True, text=True, timeout=40)
    save(Path(args.out), {'command': command, 'exit': result.returncode, 'stdout': result.stdout, 'stderr': result.stderr,
                         'compiled_js_sha256': sha(args.compiled_js), 'node_binary_sha256': sha(args.node),
                         'semantic_script_sha256': sha(semantic['script']), 'status': 'PASS' if result.returncode == 0 else 'FAIL'})
    if result.returncode:
        raise SystemExit('semantic check failed; see '+args.out)
    print(result.stdout, end='')


def main():
    p = argparse.ArgumentParser(description=__doc__)
    sub = p.add_subparsers(dest='operation', required=True)
    f = sub.add_parser('fetch'); f.add_argument('--cache', required=True); f.add_argument('--families', default=','.join(FAMILIES)); f.set_defaults(run=fetch)
    s = sub.add_parser('stage'); s.add_argument('--cache', required=True); s.add_argument('--out', required=True); s.add_argument('--families', default=','.join(FAMILIES)); s.add_argument('--java'); s.set_defaults(run=stage)
    e = sub.add_parser('specs'); e.add_argument('--fixture', required=True); e.add_argument('--out', required=True); e.add_argument('--cache'); e.add_argument('--jar', required=True); e.add_argument('--threads', type=int, default=1); e.set_defaults(run=specs)
    v = sub.add_parser('verify'); v.add_argument('--spec', required=True); v.add_argument('--compiled-js', required=True); v.add_argument('--node', required=True); v.add_argument('--out', required=True); v.set_defaults(run=verify)
    args = p.parse_args(); args.run(args)


if __name__ == '__main__':
    main()
