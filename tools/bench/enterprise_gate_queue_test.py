#!/usr/bin/env python3
"""Pure model tests; no compiler/runtime/container work."""
import copy
import unittest

import enterprise_gate_queue as queue


def ordinary():
    return {'schema': 'jsgo-enterprise-cli-v1', 'argv': ['--js', 'generated/a.js'],
            'inputs': [{'flag': '--js', 'link': 'generated', 'relative': 'a.js'}],
            'artifacts': ['out.js', 'out.js.map', 'inputs.manifest'],
            'same_work': {'chunks': []},
            'semantic_check': {'script': '/original/semantic.js', 'entry': 'realExport', 'expected': 'real result'}}


class GateQueueTest(unittest.TestCase):
    def test_original_inputs_and_runtime_oracle_required(self):
        spec = ordinary()
        queue.validate_spec(spec, True)
        for key in ('script', 'entry', 'expected'):
            invalid = copy.deepcopy(spec); del invalid['semantic_check'][key]
            with self.assertRaises(ValueError): queue.validate_spec(invalid, True)

    def test_missing_or_duplicate_original_source_rejected(self):
        for argv in ([], ['--js', 'generated/a.js', '--js', 'generated/a.js']):
            invalid = ordinary(); invalid['argv'] = argv
            with self.assertRaises(ValueError): queue.validate_spec(invalid, True)

    def test_auxiliary_source_map_not_a_fake_js_input(self):
        spec = ordinary(); spec['inputs'].append({'flag': 'input-source-map-or-sidecar', 'link': 'generated', 'relative': 'a.js.map'})
        queue.validate_spec(spec, True)

    def test_every_chunk_artifact_and_weak_absence_required(self):
        spec = ordinary()
        spec['argv'] = ['--chunk', 'part00:1', '--js', 'generated/a.js']
        spec['same_work']['chunks'] = [{'name': 'part00', 'input_count': 1, 'dependencies': []}]
        spec['artifacts'] = ['chunks/part00.js', 'chunks/part00.js.map', 'chunks/part00.js.inputs', 'chunks/$weak$.js.inputs', 'chunk-dependencies.json']
        spec['absent_artifacts'] = ['chunks/$weak$.js', 'chunks/$weak$.js.map']
        queue.validate_spec(spec, False)
        for key in ('artifacts', 'absent_artifacts'):
            invalid = copy.deepcopy(spec); invalid[key].pop()
            with self.assertRaises(ValueError): queue.validate_spec(invalid, False)
        invalid = copy.deepcopy(spec); invalid['argv'][1] = 'part00:2'
        with self.assertRaises(ValueError): queue.validate_spec(invalid, False)

    def test_real_kernel_budget_not_just_heap_cap(self):
        queue.kernel('4294967296', '0', '200000 100000', '0-1')
        queue.kernel('4294967296', '0', '200000 100000', '0,1')
        for args in [('max', '0', '200000 100000', '0-1'),
                     ('4294967296', 'max', '200000 100000', '0-1'),
                     ('4294967296', '0', '100000 100000', '0-1'),
                     ('4294967296', '0', '200000 100000', '0-3')]:
            with self.assertRaises(ValueError): queue.kernel(*args)

    def test_outer_oom_invalidates_an_otherwise_successful_task(self):
        before = {'memory.events': 'oom 2\noom_kill 1\n'}
        queue.assert_no_outer_oom(before, before)
        for events in ['oom 3\noom_kill 1\n', 'oom 2\noom_kill 2\n']:
            with self.assertRaises(RuntimeError): queue.assert_no_outer_oom(before, {'memory.events': events})


if __name__ == '__main__':
    unittest.main()
