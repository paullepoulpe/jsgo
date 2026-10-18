#!/usr/bin/env python3
"""Lightweight oracle schema/integrity tests; no compiler or JVM is launched."""
import copy
import json
from pathlib import Path
import tempfile
import unittest
import sys
sys.path.insert(0, str(Path(__file__).resolve().parent))
sys.dont_write_bytecode = True

import worker_gate as gate


class OracleTest(unittest.TestCase):
    def setUp(self):
        self.directory = tempfile.TemporaryDirectory()
        self.path = Path(self.directory.name)/'oracle.json'
        self.inputs = {'threads-1.fixture': gate.case_identity(
            ['--js', 'app.js', '--num_parallel_threads=1'], 'var é=1;', b'\x00\xff\n')}
        self.expected = {'exit': 255, 'stdout': b'\x00\xff', 'stderr': b'warning\n\tat java.frame\n',
                         'out.js': b'compiled bytes', 'out.js.map': None}
        self.oracle = {'schema': gate.ORACLE_SCHEMA, 'complete': True, 'stock_jar_sha256': 'pinned-jar',
                       'stderr_normalization': gate.STDERR_NORMALIZATION,
                       'case_order': list(self.inputs), 'corpus_sha256': gate.object_digest(self.inputs),
                       'cases': {'threads-1.fixture': {'input': self.inputs['threads-1.fixture'],
                                                      'expected': gate.encode_expected(self.expected)}}}
        self.oracle['oracle_sha256'] = gate.object_digest(self.oracle['cases'])

    def tearDown(self):
        self.directory.cleanup()

    def load(self, oracle=None, inputs=None, jar='pinned-jar'):
        self.path.write_text(json.dumps(oracle if oracle is not None else self.oracle))
        return gate.load_oracle(self.path, jar, inputs if inputs is not None else self.inputs)

    def test_round_trip_all_raw_bytes_exit_and_missing_file(self):
        _, expected, path = self.load()
        self.assertEqual(expected['threads-1.fixture'], self.expected)
        self.assertEqual(path, str(self.path))
        self.assertEqual(gate.normalize_frames(expected['threads-1.fixture']['stderr']), b'warning\n')

    def test_reject_changed_argv_source_stdin_or_threads(self):
        for identity in [gate.case_identity(['--js','other.js'], 'var é=1;', b'\x00\xff\n'),
                         gate.case_identity(['--js','app.js','--num_parallel_threads=2'], 'var é=1;', b'\x00\xff\n'),
                         gate.case_identity(self.inputs['threads-1.fixture']['args'], 'var é=2;', b'\x00\xff\n'),
                         gate.case_identity(self.inputs['threads-1.fixture']['args'], 'var é=1;', b'different')]:
            with self.assertRaises(ValueError): self.load(inputs={'threads-1.fixture': identity})
        with self.assertRaises(ValueError): self.load(inputs={})

    def test_reject_incomplete_wrong_pin_schema_normalization_or_case_set(self):
        for key, value in [('complete',False), ('schema','old'), ('stock_jar_sha256','wrong'),
                           ('stderr_normalization','weakened'), ('case_order',[])]:
            oracle = copy.deepcopy(self.oracle); oracle[key] = value
            with self.assertRaises(ValueError): self.load(oracle)
        oracle = copy.deepcopy(self.oracle); oracle['cases']['extra'] = {}
        with self.assertRaises(ValueError): self.load(oracle)

    def test_reject_corrupted_response_or_input(self):
        oracle = copy.deepcopy(self.oracle)
        oracle['cases']['threads-1.fixture']['expected']['payload']['stdout_base64'] = 'changed'
        with self.assertRaises(ValueError): self.load(oracle)
        oracle = copy.deepcopy(self.oracle)
        oracle['cases']['threads-1.fixture']['input']['args'].append('--different')
        with self.assertRaises(ValueError): self.load(oracle)
        oracle = copy.deepcopy(self.oracle); oracle['oracle_sha256'] = 'wrong'
        with self.assertRaises(ValueError): self.load(oracle)

    def test_reject_missing_bytes_invalid_base64_and_fabricated_exit(self):
        for edit in [lambda p: p.pop('stdout_base64'),
                     lambda p: p['files_base64'].pop('out.js'),
                     lambda p: p.update(exit=True), lambda p: p.update(exit=-9),
                     lambda p: p.update(stdout_base64='%%%')]:
            record = copy.deepcopy(gate.encode_expected(self.expected))
            edit(record['payload']); record['sha256'] = gate.object_digest(record['payload'])
            with self.assertRaises((ValueError, KeyError)): gate.decode_expected(record)


if __name__ == '__main__': unittest.main()
