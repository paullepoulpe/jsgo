#!/usr/bin/env python3
"""Oracle lifecycle checks against real frozen helpers; no Docker/JVM launches."""
import importlib.util
import io
import json
from pathlib import Path
import sys
import subprocess
import tempfile
from types import SimpleNamespace
import unittest
from unittest.mock import patch

sys.dont_write_bytecode = True
sys.path.insert(0, str(Path(__file__).resolve().parent))
source = importlib.util.spec_from_file_location('capped_java_oracle', Path(__file__).with_name('capped_java_oracle.py'))
oracle = importlib.util.module_from_spec(source); source.loader.exec_module(oracle)


class Process:
    def __init__(self, work, ready, exit_code=0, oom=0, required=True):
        self.commands, self.stopped = [], False
        self.stdin = self
        self.stdout = io.StringIO(ready+'\n'+'\t'.join(map(str, [
            'RESULT', 1, 'warmup', 10, 10, 0, 0, 0, 0, 1024, 2048, exit_code, 100,
            f'oom {oom};oom_kill {oom}']))+'\n'+'protocol-tail\n')
        self.work, self.exit_code, self.required = work, exit_code, required

    def write(self, text):
        self.commands.append(text)
        if text == '1\twarmup\n':
            directory = self.work/'protocol/1'; directory.mkdir(parents=True)
            for name, data in [('stdout', b''), ('stderr', b'diagnostic\n'), ('exit', f'{self.exit_code}\n'.encode())]:
                (directory/name).write_bytes(data)
            if self.required:
                (self.work/'output.js').write_bytes(b'compiled;\n')

    def flush(self): pass
    def close(self): pass
    def wait(self, timeout=None): self.stopped = True; return 0
    def poll(self): return 0 if self.stopped else None
    def kill(self): self.stopped = True


class OracleTest(unittest.TestCase):
    def exercise(self, root, *, ready='READY\t4294967296\t0\t200000 100000\t0-1', exit_code=0, oom=0, required=True, create_error=False):
        driver = root/'driver'; driver.mkdir(); (driver/'CliPaired.class').write_bytes(b'model-class')
        helper, records, frozen_driver = oracle.factors.freeze_helpers(root, driver)
        args = SimpleNamespace(binary='/bin/true', readonly_bind=[], image='model@sha256:abc',
                               java_home='/model/jdk', jar='/model/stock.jar', driver=frozen_driver,
                               source_records=records, input_manifest=[], java_dependencies=[], executable_digests={})
        engine = helper.Engine(args)
        calls = []
        def docker(*argv, **kwargs):
            calls.append(argv)
            if argv[0] == 'create':
                if create_error:
                    raise subprocess.CalledProcessError(1, ['docker', 'create'], output='', stderr='socket permission denied')
                return SimpleNamespace(stdout='model-container\n')
            if argv[0] == 'inspect':
                return SimpleNamespace(stdout=json.dumps([{'Image': 'model-image-digest', 'State': {'OOMKilled': bool(oom)},
                    'HostConfig': {'Memory': 4294967296, 'MemorySwap': 4294967296, 'NanoCpus': 2000000000,
                                   'CpusetCpus': '0,1', 'PidsLimit': 4096, 'NetworkMode': 'none'}}]))
            if argv[0] == 'rm': return SimpleNamespace(stdout='')
            raise AssertionError(argv)
        engine.call = docker
        process = Process(root/'java', ready, exit_code, oom, required)
        self.calls, self.process = calls, process
        spec = {'argv': ['--compilation_level=ADVANCED'], 'links': {}, 'inputs': [],
                'artifacts': ['output.js'], 'absent_artifacts': ['not-emitted.map'],
                'same_work': {'builtin_externs': {'pin': 'stock'}}}
        with patch.object(helper, 'java_runtime_dependencies', return_value=[]), patch.object(oracle.subprocess, 'Popen', return_value=process):
            result = oracle.invoke(helper, engine, root, spec, args)
        return result, calls, process, spec

    def test_real_engine_and_protocol_make_one_complete_compatible_oracle(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            result, calls, process, spec = self.exercise(root)
            self.assertEqual(['1\twarmup\n', 'STOP\n'], process.commands)
            self.assertEqual(1, result['compilations'])
            self.assertFalse(result['warm_target_settled'])
            digests = json.loads((root/'stock-java-oracle-digests.json').read_text())
            self.assertEqual({'stdout', 'stderr', 'exit', 'output.js', 'not-emitted.map'}, set(digests))
            self.assertIsNone(digests['not-emitted.map'])
            import go_factor_screen
            go_factor_screen.validate_oracle(spec, [], digests,
                {'spec': spec, 'input_manifest': [], 'jar_sha256': oracle.factors.PINNED_JAR}, oracle.factors.PINNED_JAR)
            self.assertIn('--memory=4294967296', calls[0])
            self.assertIn('--memory-swap=4294967296', calls[0])
            self.assertIn('/model/stock.jar:'+str(root/'runner-sources/java-driver'), calls[0])
            self.assertEqual('rm', calls[-1][0])
            self.assertIn('protocol-tail', (root/'java/driver.stdout').read_text())

    def test_rejects_wrong_ready_before_request_and_cleans_container(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            with self.assertRaisesRegex(ValueError, 'CPU budget'):
                self.exercise(root, ready='READY\t4294967296\t0\t100000 100000\t0-1')
            self.assertTrue((root/'failure.json').is_file())
            self.assertTrue((root/'java/ready.json').is_file())
            self.assertEqual('rm', self.calls[-1][0])
            self.assertTrue(self.process.stopped)
            self.assertEqual([], self.process.commands)
            self.assertFalse((root/'stock-java-oracle-digests.json').exists())
            self.assertIn('protocol-tail', (root/'java/driver.stdout').read_text())

    def test_nonzero_oom_and_missing_artifact_preserve_evidence_without_oracle(self):
        for parameters in ({'exit_code': 1}, {'oom': 1}, {'required': False}):
            with self.subTest(parameters=parameters), tempfile.TemporaryDirectory() as tmp:
                root = Path(tmp)
                with self.assertRaises((ValueError, RuntimeError)):
                    self.exercise(root, **parameters)
                failure = json.loads((root/'failure.json').read_text())
                self.assertIn('container_inspect', failure)
                self.assertEqual('rm', self.calls[-1][0])
                self.assertTrue(self.process.stopped)
                self.assertTrue((root/'java/protocol/1/stderr').is_file())
                self.assertFalse((root/'stock-java-oracle-digests.json').exists())

    def test_existing_destination_is_never_overwritten(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp); marker = root/'marker'; marker.write_text('retain')
            with self.assertRaises(SystemExit):
                oracle.main(['--spec', 'unused', '--binary', '/bin/true', '--jar', 'unused', '--driver', 'unused',
                             '--java-home', 'unused', '--image', 'unused', '--out', str(root)])
            self.assertEqual('retain', marker.read_text())

    def test_container_create_failure_retains_complete_subprocess_diagnostics(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            with self.assertRaises(subprocess.CalledProcessError):
                self.exercise(root, create_error=True)
            details = json.loads((root/'failure.json').read_text())['subprocess']
            self.assertEqual({'argv': ['docker', 'create'], 'exit': 1, 'stdout': '',
                              'stderr': 'socket permission denied'}, details)
            self.assertFalse((root/'stock-java-oracle-digests.json').exists())

    def test_unpinned_image_preserves_preflight_failure(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)/'new'
            with self.assertRaisesRegex(ValueError, 'pinned by digest'):
                oracle.main(['--spec', 'unused', '--binary', '/bin/true', '--jar', 'unused', '--driver', 'unused',
                             '--java-home', 'unused', '--image', 'unversioned', '--out', str(root)])
            self.assertTrue((root/'failure.json').is_file())

    def test_nonstock_jar_is_rejected_before_any_container(self):
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp); jar = root/'not-stock.jar'; jar.write_bytes(b'modified compiler')
            with self.assertRaisesRegex(ValueError, 'repository-pinned stock jar'):
                oracle.main(['--spec', 'unused', '--binary', '/bin/true', '--jar', str(jar), '--driver', str(root),
                             '--java-home', str(root), '--image', 'image@sha256:'+'0'*64, '--out', str(root/'evidence')])
            self.assertTrue((root/'evidence/failure.json').is_file())
            self.assertFalse((root/'evidence/java').exists())


if __name__ == '__main__':
    unittest.main()
