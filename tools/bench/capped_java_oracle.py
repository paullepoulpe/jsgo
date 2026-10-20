#!/usr/bin/env python3
"""One fresh stock-Java artifact oracle under the paired runner's exact cap.

This is diagnostic correctness evidence only. One protocol ``warmup`` request
is followed by STOP; there is no settled warm JVM or performance ratio.
No compiler, driver, image or input is built or downloaded by this command.
"""
import argparse
import json
from pathlib import Path
import re
import subprocess
import sys

sys.dont_write_bytecode = True
import paired_factors as factors

HERE = Path(__file__).resolve().parent
MODE = 'diagnostic artifact oracle only; fresh JVM, one compile; no performance or warm ratio'


class Transcript:
    def __init__(self, stream, evidence):
        self.stream, self.evidence = stream, evidence

    def readline(self):
        line = self.stream.readline()
        self.evidence.write(line)
        self.evidence.flush()
        return line

    def drain(self):
        self.evidence.write(self.stream.read())
        self.evidence.flush()


def invoke(helper, engine, root, spec, args):
    """Run the unchanged driver and retain raw protocol, even on cap failure."""
    work = root/'java'
    helper.workspace(work, spec)
    (work/'argv.txt').write_text('\n'.join(spec['argv'])+'\n')
    command = ['python3', str(helper.HERE/'cgroup_exec.py'), str(Path(args.java_home)/'bin/java'),
               '-Xmx3g', '-cp', args.jar+':'+args.driver, 'CliPaired', str(work/'argv.txt')]
    helper.save(work/'request.json', {'argv': command, 'compiler_argv': spec['argv'],
                'mode': MODE, 'protocol_requests': ['1\twarmup', 'STOP'],
                'environment': {'LC_ALL': 'C.UTF-8'}, 'protocol_kind': 'warmup',
                'semantic_kind': 'diagnostic artifact gate only'})
    before = helper.outer_state()
    helper.save(root/'outer-before.json', before)
    process, cid = None, None
    errors, transcript = (work/'driver.stderr').open('w'), (work/'driver.stdout').open('w')
    try:
        cid = engine.create(work, command, java=True)
        process = subprocess.Popen(helper.DOCKER+['start', '-a', '-i', cid], stdin=subprocess.PIPE,
                                   stdout=subprocess.PIPE, stderr=errors, text=True, env=engine.env)
        process.stdout = Transcript(process.stdout, transcript)
        ready_line = process.stdout.readline()
        helper.save(work/'ready.json', {'raw': ready_line})
        helper.save(work/'ready.json', {'raw': ready_line, 'fields': factors.validate_ready(ready_line)})
        # The frozen driver accepts warmup|measure, not "diagnostic".
        # Retain that actual protocol kind and label its purpose separately.
        row = helper.java_request(process, work, spec, 1, 'warmup')
        row.update({'semantic_kind': 'diagnostic artifact gate only', 'mode': MODE})
        helper.save(root/'java-row.json', row)
        process.stdin.write('STOP\n'); process.stdin.flush(); process.stdin.close()
        status = process.wait(timeout=30)
        process.stdout.drain()
        helper.save(work/'driver-exit.json', {'exit': status})
        if status != 0:
            raise RuntimeError('stock Java driver failed after its single diagnostic compile')
        after = helper.outer_state()
        helper.save(root/'outer-after.json', after)
        if any(helper.event_count(after.get('memory.events', ''), name) !=
               helper.event_count(before.get('memory.events', ''), name) for name in ('oom', 'oom_kill')):
            raise RuntimeError('outer cgroup OOM invalidates stock Java oracle')
        factors.check_frozen(args.source_records)
        if helper.input_identity(spec) != args.input_manifest:
            raise ValueError('source input changed during stock Java oracle')
        if helper.java_runtime_dependencies(args.java_home) != args.java_dependencies:
            raise ValueError('JDK external runtime dependencies changed during stock Java oracle')
        for path, digest in args.executable_digests.items():
            if helper.sha(path) != digest:
                raise ValueError('executable changed during stock Java oracle: '+path)
        helper.save(root/'stock-java-oracle-digests.json', row['digests'])
        result = {'status': 'COMPLETE', 'mode': MODE, 'protocol_kind': 'warmup',
                  'compilations': 1, 'warm_target_settled': False, 'row': row}
        helper.save(root/'result.json', result)
        return result
    except BaseException as error:
        failure = {'type': type(error).__name__, 'message': str(error), 'mode': MODE,
                   'outer_state': helper.outer_state()}
        if cid is not None:
            try:
                failure['container_inspect'] = json.loads(engine.call('inspect', cid, capture_output=True, text=True).stdout)
            except Exception as inspect_error:
                failure['inspect_error'] = str(inspect_error)
        helper.save(root/'failure.json', failure)
        raise
    finally:
        # Remove the capped container before reaping docker attach, including
        # startup/READY failures and an OOM that never produced a RESULT row.
        try:
            engine.close()
        finally:
            try:
                if process is not None:
                    if process.poll() is None:
                        try:
                            process.wait(timeout=10)
                        except subprocess.TimeoutExpired:
                            process.kill(); process.wait()
                    process.stdout.drain()
            finally:
                errors.close(); transcript.close()

def main(argv=None):
    parser = argparse.ArgumentParser(description=__doc__)
    for name in ('spec', 'binary', 'jar', 'driver', 'java-home', 'image', 'out'):
        parser.add_argument('--'+name, required=True)
    parser.add_argument('--jar-sha256', help='optional expected hash, in addition to the mandatory repository pin')
    parser.add_argument('--readonly-bind', action='append', default=[])
    args = parser.parse_args(argv)
    root = Path(args.out).resolve()
    if root.exists():
        parser.error('refusing to overwrite prior oracle evidence')
    root.mkdir(parents=True)
    try:
        if not re.fullmatch(r'.+@sha256:[0-9a-f]{64}', args.image):
            raise ValueError('require an already installed image pinned by digest')
        for name in ('binary', 'jar', 'driver', 'java_home'):
            path = Path(getattr(args, name))
            if not path.is_absolute() or not path.exists():
                raise ValueError(name+' must be an existing absolute path')
        jar_hash = factors.sha(args.jar)
        pin = HERE.parent.parent/'tools/oracle/pin.env'
        config = dict(line.strip().split('=', 1) for line in pin.read_text().splitlines()
                      if line.strip() and not line.lstrip().startswith('#'))
        if jar_hash != factors.PINNED_JAR or jar_hash not in config['CLOSURE_JAR_SHA256'].strip('"').split():
            raise ValueError('Java oracle must use the unchanged repository-pinned stock jar')
        if args.jar_sha256 is not None and args.jar_sha256 != jar_hash:
            raise ValueError('supplied Java jar digest differs from pinned stock jar')
        spec_bytes = Path(args.spec).read_bytes()
        (root/'spec.json').write_bytes(spec_bytes)
        spec = json.loads(spec_bytes)
        if any(not isinstance(value, str) or any(c in value for c in '\r\n\x00') for value in spec['argv']):
            raise ValueError('compiler argv must contain strings without newline or NUL')
        helper, records, frozen_driver = factors.freeze_helpers(root, args.driver)
        own_snapshot = root/'runner-sources'/Path(__file__).name
        own_snapshot.write_bytes(Path(__file__).read_bytes())
        records[own_snapshot.name] = {'source': str(Path(__file__).resolve()),
                                      'snapshot': str(own_snapshot), 'sha256': helper.sha(own_snapshot)}
        helper.save(root/'source-manifest.json', records)
        args.driver = frozen_driver
        identity = helper.input_identity(spec)
        dependencies = helper.java_runtime_dependencies(args.java_home)
        args.readonly_bind = sorted(set(args.readonly_bind+list(spec['links'].values())))
        executable_digests = {path: helper.sha(path) for path in
                             (args.binary, args.jar, str(Path(args.java_home)/'bin/java'),
                              str(Path(args.driver)/'CliPaired.class'))}
        helper.save(root/'metadata.json', vars(args) | {'schema': 'jsgo-capped-stock-java-oracle-v1',
                    'mode': MODE, 'spec': spec, 'input_manifest': identity, 'spec_sha256': helper.sha(args.spec),
                    'jar_sha256': jar_hash, 'java_flags': ['-Xmx3g'], 'protocol_kind': 'warmup',
                    'semantic_kind': 'diagnostic artifact gate only', 'java_runtime_dependencies': dependencies,
                    'executable_digests': executable_digests, 'go_binary_role': 'Engine read-only mount only; never executed',
                    'warm_target_settled': False, 'compilations': 1})
        args.source_records, args.input_manifest = records, identity
        args.java_dependencies, args.executable_digests = dependencies, executable_digests
        engine = helper.Engine(args)
        result = invoke(helper, engine, root, spec, args)
        print(json.dumps(result, indent=2))
    except BaseException as error:
        if not (root/'failure.json').exists():
            factors.save(root/'failure.json', {'type': type(error).__name__, 'message': str(error), 'mode': MODE})
        raise


if __name__ == '__main__':
    main()
