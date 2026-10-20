#!/usr/bin/env python3
"""Pure runner-model checks; no Docker, JVM, compiler or downloaded inputs."""
import importlib.util
import json
from pathlib import Path
import statistics
import sys
import tempfile
import unittest
from types import SimpleNamespace
from unittest.mock import patch

sys.dont_write_bytecode = True
spec = importlib.util.spec_from_file_location('factors', Path(__file__).with_name('paired_factors.py'))
factors = importlib.util.module_from_spec(spec)
spec.loader.exec_module(factors)


class FactorProtocolModelTest(unittest.TestCase):
    def test_every_go_has_one_java_and_variant_position_balance(self):
        names = ['stock', 'off']+[f'factor-{n}' for n in range(6)]+['compound']
        schedule = factors.balanced_schedule(names, 3, 5, 12)
        counts = {n: [0]*len(names) for n in names}
        directions = {n: [] for n in names}
        for cohort in schedule:
            self.assertEqual(5, len(cohort['repetitions']))
            for repetition in cohort['repetitions']:
                pairs = repetition['pairs']
                self.assertEqual(set(names), {p['variant'] for p in pairs})
                self.assertEqual(len(names), len(pairs))
                for pair in pairs:
                    self.assertEqual({'go', 'java'}, set(pair['order']))
                    counts[pair['variant']][pair['position']-1] += 1
                    directions[pair['variant']].append(pair['order'][0])
        for values in counts.values():
            self.assertLessEqual(max(values)-min(values), 1)
        for values in directions.values():
            self.assertTrue(all(a != b for a, b in zip(values, values[1:])))
            self.assertLessEqual(abs(values.count('java')-values.count('go')), 1)
        self.assertEqual(schedule, factors.balanced_schedule(names, 3, 5, 12))
        self.assertNotEqual(schedule, factors.balanced_schedule(names, 3, 5, 13))

    def test_exact_helper_snapshot_and_mutation_detection(self):
        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            driver = root/'original-driver'; driver.mkdir()
            (driver/'CliPaired.class').write_bytes(b'model-only-class-bytes-never-executed')
            out = root/'evidence'; out.mkdir()
            helper, records, frozen_driver = factors.freeze_helpers(out, driver)
            self.assertEqual(out/'runner-sources', helper.HERE)
            self.assertEqual(b'model-only-class-bytes-never-executed', (Path(frozen_driver)/'CliPaired.class').read_bytes())
            factors.check_frozen(records)
            (out/'runner-sources/limited_cli.py').write_text('changed')
            with self.assertRaisesRegex(ValueError, 'frozen run helper changed'):
                factors.check_frozen(records)

    def test_current_source_helper_must_match_before_import(self):
        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            (root/'paired_cli.py').write_text('raise Exception("must never execute")')
            with patch.object(factors, 'HERE', root):
                with self.assertRaisesRegex(ValueError, 'helper source changed'):
                    factors.freeze_helpers(root/'out', root/'driver')

    def test_variants_cannot_change_work_or_runtime_budget(self):
        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            binary = root/'binary'; binary.write_bytes(b'not-executed')
            config = root/'variants.json'
            base = {'name': 'stock', 'binary': str(binary), 'go_env': {}}
            config.write_text(json.dumps([base]))
            self.assertEqual(factors.sha(binary), factors.load_variants(config)[0]['binary_sha256'])
            for change in ({'argv': ['--checks_only']}, {'go_env': {'GOGC': '200'}},
                           {'go_env': {'JSGO_MAIN_CPU_PROFILE': 'profile'}}, {'sha256': 'wrong'}, {'name': '../escape'}):
                config.write_text(json.dumps([base | change]))
                with self.assertRaises(ValueError):
                    factors.load_variants(config)

    def test_common_java_lifetime_peak_for_all_memory_ratios(self):
        helper = SimpleNamespace(stats=lambda xs: {'median': statistics.median(xs)}, warm_stability=lambda rows: {'stable': True})
        pairs = []
        for name, java_hwm, go_hwm, java_time, go_time in [('early', 32, 4, 10, 5), ('late', 64, 8, 12, 8)]:
            pairs.append({'variant': name, 'rows': {
                'java': {'peak_rss_mib': java_hwm, 'wall_s': java_time, 'cpu_s': java_time},
                'go': {'peak_rss_mib': go_hwm, 'wall_s': go_time, 'cpu_s': go_time}}})
        summary = factors.summarize_cohort(helper, 1, [{'name': 'early'}, {'name': 'late'}], pairs, [], {}, {})
        self.assertEqual(16, summary['variants']['early']['java_over_go_rss'])
        self.assertEqual(8, summary['variants']['late']['java_over_go_rss'])
        self.assertEqual(2, summary['variants']['early']['java_over_go_time'])
        self.assertEqual(1.5, summary['variants']['late']['java_over_go_time'])
        self.assertEqual(32, summary['variants']['early']['java']['peak_rss_mib']['median'])

    def test_actual_kernel_ready_budget_required(self):
        factors.validate_ready('READY\t4294967296\t0\t200000 100000\t0-1\n')
        for line in ('READY\t4294967296\t1\t200000 100000\t0-1',
                     'READY\t4294967296\t0\t100000 100000\t0-1',
                     'READY\t4294967296\t0\t200000 100000\t0-2', ''):
            with self.assertRaises(ValueError):
                factors.validate_ready(line)


if __name__ == '__main__':
    unittest.main()
