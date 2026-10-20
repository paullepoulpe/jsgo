#!/usr/bin/env python3
"""Pure exploratory-screen model checks; no compiler, JVM or Docker jobs."""
import hashlib
import importlib.util
import json
from pathlib import Path
import sys
import tempfile
from types import SimpleNamespace
import unittest

sys.dont_write_bytecode = True
spec = importlib.util.spec_from_file_location('screen', Path(__file__).with_name('go_factor_screen.py'))
screen = importlib.util.module_from_spec(spec); spec.loader.exec_module(screen)


class ScreenModelTest(unittest.TestCase):
    def test_schedule_has_every_variant_once_and_balanced_positions(self):
        names = ['stock', 'off']+[f'single-{n}' for n in range(7)]+['compound']
        for rounds in (3, 10, 13, 23):
            schedule = screen.balanced_rounds(names, rounds, 123)
            positions = {n: [0]*len(names) for n in names}
            for row in schedule:
                self.assertEqual(set(names), set(row['variants']))
                self.assertEqual(len(names), len(row['variants']))
                for index, name in enumerate(row['variants']):
                    positions[name][index] += 1
            for values in positions.values():
                self.assertLessEqual(max(values)-min(values), 1)
            self.assertEqual(schedule, screen.balanced_rounds(names, rounds, 123))

    def test_fixed_stock_oracle_requires_same_complete_work_and_absence(self):
        task = {'argv': ['--checks_only'], 'links': {'inputs': '/fixture'}, 'inputs': [],
                'artifacts': [], 'absent_artifacts': ['inputs.manifest'], 'same_work': {'builtin_externs': {'pin': 'same'}}}
        oracle = {'stdout': hashlib.sha256(b'').hexdigest(), 'stderr': hashlib.sha256(b'warning').hexdigest(),
                  'exit': hashlib.sha256(b'0\n').hexdigest(), 'inputs.manifest': None}
        metadata = {'jar_sha256': 'stock', 'spec': task, 'input_manifest': []}
        screen.validate_oracle(task, [], oracle, metadata, 'stock')
        bad_task = task | {'argv': ['--compilation_level=WHITESPACE_ONLY']}
        for data, meta in [(oracle, metadata | {'jar_sha256': 'different'}),
                           (oracle, metadata | {'spec': bad_task}),
                           (oracle | {'inputs.manifest': hashlib.sha256(b'').hexdigest()}, metadata),
                           ({k: v for k, v in oracle.items() if k != 'stderr'}, metadata),
                           (oracle | {'exit': hashlib.sha256(b'1\n').hexdigest()}, metadata)]:
            with self.assertRaises(ValueError):
                screen.validate_oracle(task, [], data, meta, 'stock')

    def test_ratios_are_paired_same_round_go_over_control(self):
        variants = [{'name': 'off'}, {'name': 'factor'}]
        rows = []
        # Deliberately vary host speed between rounds; ratio must use the
        # matching round's control rather than ratio of all pooled samples.
        for number, control, candidate in [(1, 10, 5), (2, 100, 80), (3, 20, 24)]:
            for name, seconds in [('off', control), ('factor', candidate)]:
                rows.append({'row': len(rows)+1, 'round': number, 'variant': name, 'gate': 'MATCH',
                             'metrics': {'wall_s': seconds, 'cpu_s': seconds, 'peak_rss_mib': 10 if name == 'off' else 8}})
        result = screen.summarize(rows, variants, 'off', 6)
        self.assertTrue(result['complete'])
        self.assertEqual([.5, .8, 1.2], [r['time_go_over_control'] for r in result['variants']['factor']['same_round_ratios']])
        self.assertEqual(.8, result['variants']['factor']['time_go_over_control']['median'])
        self.assertEqual(.8, result['variants']['factor']['rss_go_over_control']['median'])
        self.assertEqual(0, result['variants']['factor']['rss_go_over_control']['sample_variance'])
        self.assertNotIn('java_over_go_time', json.dumps(result))
        self.assertNotIn('java_over_go_rss', json.dumps(result))
        self.assertFalse(screen.summarize(rows[:-1], variants, 'off', 6)['complete'])

    def test_actual_engine_create_go_path_needs_no_java_attributes(self):
        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            helper, unused, records = screen.freeze(root)
            # Deliberately omit java_home, jar and driver: Go must not inspect
            # or mount Java dependencies. Exercise the real create() method.
            engine = helper.Engine(SimpleNamespace(binary='/bin/true', readonly_bind=[], image='pinned-model-image'))
            calls = []
            def docker(*argv, **kwargs):
                calls.append(argv)
                if argv[0] == 'create':
                    return SimpleNamespace(stdout='model-container\n')
                if argv[0] == 'inspect':
                    return SimpleNamespace(stdout=json.dumps([{'Image': 'pinned-model-image', 'HostConfig': {
                        'Memory': 4294967296, 'MemorySwap': 4294967296, 'NanoCpus': 2000000000,
                        'CpusetCpus': '0,1', 'PidsLimit': 4096, 'NetworkMode': 'none'}}]))
                raise AssertionError(argv)
            engine.call = docker
            work = root/'go-work'; work.mkdir()
            self.assertEqual('model-container', engine.create(work, ['model-cli']))
            self.assertIn('--memory=4294967296', calls[0])
            self.assertFalse(any('java-home' in arg or 'java-driver' in arg for arg in calls[0]))

    def test_imported_helpers_are_pinned_snapshots(self):
        with tempfile.TemporaryDirectory() as temporary:
            helper, factors, records = screen.freeze(Path(temporary))
            self.assertEqual(Path(temporary)/'runner-sources', helper.HERE)
            factors.check_frozen(records)
            (helper.HERE/'limited_cli.py').write_text('changed')
            with self.assertRaises(ValueError):
                factors.check_frozen(records)


if __name__ == '__main__':
    unittest.main()
