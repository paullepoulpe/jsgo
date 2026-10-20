#!/usr/bin/env python3
"""Stage a genuine 2026 generated J2CL JRE and freeze its compiler work.

No generator or compiler is run here. Both jszip artifacts are published by the
same Maven distribution/version and are extracted without modification. Failed
old Vertispan fixture evidence is separate and must be retained. Full stock-Java,
Go and independent compiled-output semantic gates establish compatibility later.
"""
import argparse
import hashlib
import json
from pathlib import Path
import shutil
import sys

sys.dont_write_bytecode = True
HERE = Path(__file__).resolve().parent
import bench
import generated_fixtures as fixtures
from enterprise_specs import embedded_externs, replace_threads

PINS = HERE/'generated/modern_j2cl_assets.json'
REQUIRED_PROPERTIES = ['jre.checks.'+name for name in ('type', 'numeric', 'bounds', 'api', 'checkLevel')]+['jre.checkedMode']


def stage(args):
    cache, out = Path(args.cache).resolve(), Path(args.out).resolve()
    pins = json.loads(PINS.read_text())
    if out.exists():
        raise ValueError('refusing to overwrite prior fixture/evidence')
    assets = []
    for row in pins['assets']:
        path = cache/row['name']; data = path.read_bytes()
        if (len(data) != row['bytes'] or hashlib.sha256(data).hexdigest() != row['sha256'] or
            hashlib.sha1(data).hexdigest() != row['published_sha1']):
            raise ValueError('published JRE/bootstrap asset pin mismatch: '+str(path))
        assets.append(row | {'verified_sha256_and_published_sha1': True})
    out.mkdir(parents=True)
    record = {'schema': 'jsgo-modern-j2cl-fixture-v1', 'distribution': pins['distribution'],
              'version': pins['version'], 'provenance': pins['provenance'], 'assets': assets,
              'pins_sha256': fixtures.sha(PINS), 'builder_sha256': fixtures.sha(__file__),
              'unzip_helper_sha256': fixtures.sha(HERE/'generated_fixtures.py'),
              'status': 'PREPARING; unvalidated compiler/runtime compatibility'}
    try:
        for name, directory in [('jre-20260527-jszip.zip', 'jre'), ('bootstrap-20260527-jszip.zip', 'bootstrap')]:
            fixtures.unzip(cache/name, out/'j2cl'/directory)
        sources = list((out/'j2cl/jre').rglob('java/lang/jre.js'))
        if len(sources) != 1:
            raise ValueError('cannot uniquely identify real JRE system-property registrations')
        source = sources[0].read_text()
        if not all("addSystemPropertyFromGoogDefine('"+name+"'," in source for name in REQUIRED_PROPERTIES):
            raise ValueError('generated JRE lacks current compiler registration ABI')
        record['abi_registration_source'] = sources[0].relative_to(out).as_posix()
        record['abi_registration_sha256'] = fixtures.sha(sources[0])
        record['required_registrations'] = REQUIRED_PROPERTIES
        shutil.copyfile(HERE/'generated/modern_j2cl_entry.js', out/'modern_j2cl_entry.js')
        record.update(status='STAGED; genuine matching JRE/bootstrap; full Java/Go/Node gates pending', files=fixtures.inventory(out))
    except BaseException as error:
        record['failure'] = {'type': type(error).__name__, 'message': str(error)}
        raise
    finally:
        fixtures.save(out/'fixture-manifest.json', record)


def specs(args):
    fixture, out = Path(args.fixture).resolve(), Path(args.out).resolve()
    record = json.loads((fixture/'fixture-manifest.json').read_text())
    if not record['status'].startswith('STAGED') or record['files'] != fixtures.inventory(fixture):
        raise ValueError('fixture incomplete or changed after staging')
    if out.exists() or args.threads < 1:
        raise ValueError('need new output directory and positive thread count')
    out.mkdir(parents=True)
    sources, auxiliary = [], []
    argv = ['--compilation_level=ADVANCED', '--jscomp_warning=checkTypes', '--jscomp_warning=checkVars',
            '--language_out=ECMASCRIPT_2017', '--dependency_mode=PRUNE', '--entry_point=goog:enterprise.modernJ2cl',
            '--j2cl_pass=AUTO']
    for path in sorted(p for p in fixture.rglob('*') if p.is_file() and p.name != 'fixture-manifest.json'):
        item = {'link': 'generated', 'relative': path.relative_to(fixture).as_posix(),
                'bytes': path.stat().st_size, 'sha256': fixtures.sha(path)}
        if path.suffix == '.js':
            argv += ['--js', 'generated/'+item['relative']]
            sources.append(item | {'flag': '--js'})
        else:
            auxiliary.append(item | {'flag': 'input-source-map-or-sidecar'})
    argv = replace_threads(argv+['--js_output_file=out.js', '--create_source_map=out.js.map', '--output_manifest=inputs.manifest'], args.threads)
    spec = {'schema': 'jsgo-enterprise-cli-v1', 'name': 'generated-modern-j2cl-jre', 'argv': argv,
            'links': {'generated': str(fixture)}, 'inputs': sources+auxiliary,
            'artifacts': ['out.js', 'out.js.map', 'inputs.manifest'], 'fixture_manifest_sha256': fixtures.sha(fixture/'fixture-manifest.json'),
            'semantic_check': {'script': str(HERE/'generated/verify_semantics.js'), 'entry': 'jsgoEnterpriseModernJ2cl',
                               'expected': 'J2CL:13,11;remaining=1',
                               'required': 'run independently on stock Java and Go ADVANCED output before performance claims'},
            'same_work': {'mode': 'ADVANCED, current J2CL AUTO, PRUNE, checkTypes/checkVars, source map, manifest',
                          'compiler_threads': args.threads, 'builtin_externs': embedded_externs(args.jar, Path(bench.REPO)),
                          'explicit_externs': [], 'source_count': len(sources), 'auxiliary_input_count': len(auxiliary),
                          'source_bytes': sum(s['bytes'] for s in sources), 'generated_family': 'modern-j2cl-jre',
                          'distribution': record['distribution'], 'version': record['version'],
                          'provenance': record['provenance'], 'fixture_status': record['status'],
                          'limitations': 'Published generated JRE and matching bootstrap with benchmark-owned JS collection harness; not a production Java application. No runtime patches/stubs or disabled J2CL pass.'}}
    fixtures.save(out/'generated-modern-j2cl-jre.json', spec)
    print(json.dumps({'spec': str(out/'generated-modern-j2cl-jre.json'), 'source_count': len(sources),
                      'auxiliary_count': len(auxiliary), 'source_bytes': spec['same_work']['source_bytes'], 'status': record['status']}))


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    sub = parser.add_subparsers(dest='operation', required=True)
    command = sub.add_parser('stage'); command.add_argument('--cache', required=True); command.add_argument('--out', required=True); command.set_defaults(run=stage)
    command = sub.add_parser('specs'); command.add_argument('--fixture', required=True); command.add_argument('--out', required=True)
    command.add_argument('--jar', required=True); command.add_argument('--threads', type=int, default=2); command.set_defaults(run=specs)
    args = parser.parse_args(); args.run(args)


if __name__ == '__main__':
    main()
