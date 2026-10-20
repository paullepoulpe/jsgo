#!/usr/bin/env python3
"""Planner-model tests only; no extraction, compiler or fabricated benchmark."""
import importlib.util
from pathlib import Path
import sys
import unittest
sys.dont_write_bytecode = True
spec = importlib.util.spec_from_file_location('planner', Path(__file__).with_name('many_chunk_specs.py'))
planner = importlib.util.module_from_spec(spec); spec.loader.exec_module(planner)


class ChunkPlannerTest(unittest.TestCase):
    def fixture(self, count):
        items, records, argv = [], [], ['--compilation_level=ADVANCED', '--dependency_mode=SORT_ONLY']
        for n in range(count):
            path, digest = f'f{n}.js', f'{n:064x}'
            items.append({'flag': '--js', 'link': 'inputs', 'relative': path, 'bytes': n+1, 'sha256': digest})
            records.append({'path': 'inputs/'+path, 'bytes': n+1, 'sha256': digest, 'provides': ['n'+str(n)],
                            'requires': [] if n == 0 else ['n'+str(n-1)], 'type_requires': []})
            argv += ['--js', 'inputs/'+path]
        argv += ['--js_output_file=out.js', '--create_source_map=out.js.map', '--output_manifest=inputs.manifest']
        base = {'name': 'model-only', 'argv': argv, 'inputs': items, 'links': {}, 'artifacts': ['out.js'], 'same_work': {}}
        inventory = {'stock_jar_sha256': planner.PINNED_JAR, 'errors': [], 'files': records}
        return base, inventory

    def test_all_real_paths_kept_once_dependencies_precede_32_nonempty_chunks(self):
        base, deps = self.fixture(50)
        result, proof = planner.plan(base, deps, 32)
        self.assertEqual(32, len(proof['chunks']))
        self.assertEqual(set(x['relative'] for x in base['inputs']), set(x['relative'] for x in result['inputs']))
        self.assertEqual(50, len(result['inputs']))
        self.assertTrue(all(c['input_count'] for c in proof['chunks']))
        self.assertEqual(98, len(result['artifacts']))
        self.assertIn('chunks/$weak$.js.inputs', result['artifacts'])
        self.assertEqual(['chunks/$weak$.js', 'chunks/$weak$.js.map'], result['absent_artifacts'])
        self.assertIn('--compilation_level=ADVANCED', result['argv'])
        self.assertIn('--dependency_mode=SORT_ONLY', result['argv'])

    def test_type_cycle_scc_is_never_split(self):
        base, deps = self.fixture(40)
        deps['files'][10]['type_requires'] = ['n12']
        result, proof = planner.plan(base, deps, 32)
        positions = {node: i for i, c in enumerate(proof['chunks']) for node in c['original_input_indices']}
        self.assertEqual(positions[10], positions[11]); self.assertEqual(positions[11], positions[12])

    def test_unresolved_duplicate_changed_sources_or_wrong_stock_fail(self):
        for kind in ('unresolved', 'duplicate', 'changed', 'jar', 'order'):
            base, deps = self.fixture(40)
            if kind == 'unresolved': deps['files'][0]['requires'] = ['missing']
            if kind == 'duplicate': deps['files'][1]['provides'] = ['n0']
            if kind == 'changed': deps['files'][1]['sha256'] = 'wrong'
            if kind == 'jar': deps['stock_jar_sha256'] = 'wrong'
            if kind == 'order': deps['files'].reverse()
            with self.assertRaises(ValueError): planner.plan(base, deps, 32)

    def test_zero_byte_tail_and_iterative_long_chain(self):
        buckets = planner.partition([[n] for n in range(40)], [1]*35+[0]*5, 32)
        self.assertEqual(list(range(40)), [n for b in buckets for n in b])
        edges = [set()] + [{n-1} for n in range(1, 2500)]
        owner, groups = planner.strong_components(edges)
        ordered = planner.topological_components(edges, owner, groups)
        self.assertEqual(list(range(2500)), [n for g in ordered for n in g])


if __name__ == '__main__': unittest.main()
