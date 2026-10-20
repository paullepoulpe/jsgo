#!/usr/bin/env python3
"""Plan a 32-chunk version of unchanged real inputs from stock AST dependencies.

Setup only. Extract dependencies with CliDependencyInventory from the pinned jar
before invoking this tool. Stock Java must then pass the actual 4GiB compilation
and complete artifact gate before Go correctness or timing. No regex dependency
scan, source replication, stubs, pruning or weakened compiler mode is introduced.
"""
import argparse
import hashlib
import heapq
import json
from pathlib import Path
import sys

sys.dont_write_bytecode = True
HERE = Path(__file__).resolve().parent
PINNED_JAR = 'd6864ebf70bca5ba6e237898de6611a389aaf8b933c2549036ced4a684b664b6'


def sha(path):
    return hashlib.sha256(Path(path).read_bytes()).hexdigest()


def save(path, value):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(value, indent=2)+'\n')


def strong_components(edges):
    # Iterative Kosaraju avoids recursion-limit failures on deep real chains.
    seen, finish = set(), []
    reverse = [set() for unused in edges]
    for node, parents in enumerate(edges):
        for parent in parents:
            reverse[parent].add(node)
    for start in range(len(edges)):
        if start in seen:
            continue
        seen.add(start); stack = [(start, iter(sorted(edges[start])))]
        while stack:
            node, children = stack[-1]
            child = next(children, None)
            if child is None:
                finish.append(node); stack.pop()
            elif child not in seen:
                seen.add(child); stack.append((child, iter(sorted(edges[child]))))
    owner, groups = {}, []
    for start in reversed(finish):
        if start in owner:
            continue
        group, stack = [], [start]
        owner[start] = len(groups)
        while stack:
            node = stack.pop(); group.append(node)
            for child in sorted(reverse[node], reverse=True):
                if child not in owner:
                    owner[child] = len(groups); stack.append(child)
        groups.append(sorted(group))
    return owner, groups


def topological_components(edges, owner, groups):
    dependencies, consumers = [set() for unused in groups], [set() for unused in groups]
    for node, parents in enumerate(edges):
        for parent in parents:
            a, b = owner[node], owner[parent]
            if a != b:
                dependencies[a].add(b); consumers[b].add(a)
    pending = [len(x) for x in dependencies]
    ready = [(min(groups[index]), index) for index, size in enumerate(pending) if not size]
    heapq.heapify(ready); result = []
    while ready:
        unused, group = heapq.heappop(ready); result.append(groups[group])
        for child in sorted(consumers[group]):
            pending[child] -= 1
            if not pending[child]:
                heapq.heappush(ready, (min(groups[child]), child))
    if len(result) != len(groups):
        raise AssertionError('SCC condensation must be acyclic')
    return result


def partition(groups, sizes, count):
    if len(groups) < count:
        raise ValueError('fewer SCCs than requested chunks; refusing empty/fake chunks')
    remaining = sum(sizes[node] for group in groups for node in group)
    buckets, offset = [], 0
    for index in range(count):
        left = count-index
        target, weight, bucket = remaining/left, 0, []
        # Keep at least one complete SCC for every later bucket.
        while offset < len(groups)-(left-1):
            group = groups[offset]
            group_weight = sum(sizes[node] for node in group)
            if left > 1 and bucket and abs(weight-target) <= abs(weight+group_weight-target):
                break
            bucket.extend(group); weight += group_weight; offset += 1
        if not bucket:
            raise AssertionError('empty chunk')
        buckets.append(bucket); remaining -= weight
    if offset != len(groups):
        raise AssertionError('partition dropped an SCC')
    return buckets


def plan(base, inventory, count):
    if inventory.get('stock_jar_sha256') != PINNED_JAR or inventory.get('errors'):
        raise ValueError('need successful dependency extraction from exact pinned stock jar')
    source_items = [x for x in base['inputs'] if x['flag'] == '--js']
    expected = [x['link']+'/'+x['relative'] for x in source_items]
    if len(set(expected)) != len(expected):
        raise ValueError('duplicate original input paths')
    records = inventory['files']
    if [x['path'] for x in records] != expected:
        raise ValueError('dependency inventory does not cover exact original ordered sources')
    for source, dependency in zip(source_items, records):
        if (source['sha256'], source['bytes']) != (dependency['sha256'], dependency['bytes']):
            raise ValueError('source changed between original spec and dependency extraction')
    providers = {}
    for index, record in enumerate(records):
        for symbol in record['provides']:
            if symbol in providers and providers[symbol] != index:
                raise ValueError('ambiguous namespace provider: '+symbol)
            providers[symbol] = index
    edges, missing = [], []
    for index, record in enumerate(records):
        parents = set()
        for kind, requires in [('value', record['requires']), ('type', record['type_requires'])]:
            for symbol in requires:
                if symbol not in providers:
                    missing.append({'input': record['path'], 'kind': kind, 'symbol': symbol})
                elif providers[symbol] != index:
                    parents.add(providers[symbol])
        edges.append(parents)
    if missing:
        # Never silently guess which unresolved dependencies are global externs.
        raise ValueError('unresolved provider dependencies require an explicit audited plan: '+json.dumps(missing[:12]))
    owner, groups = strong_components(edges)
    ordered = topological_components(edges, owner, groups)
    buckets = partition(ordered, [x['bytes'] for x in records], count)
    chunk_of = {node: index for index, bucket in enumerate(buckets) for node in bucket}
    if sorted(chunk_of) != list(range(len(records))) or any(chunk_of[parent] > chunk_of[node] for node, parents in enumerate(edges) for parent in parents):
        raise AssertionError('partition dropped/duplicated files or put a provider in a later chunk')
    argv, index = [], 0
    while index < len(base['argv']):
        arg = base['argv'][index]
        if arg in ('--js', '--js_output_file', '--create_source_map', '--output_manifest'):
            index += 2
        elif any(arg.startswith(flag+'=') for flag in ('--js_output_file', '--create_source_map', '--output_manifest')):
            index += 1
        elif arg.startswith('--chunk'):
            raise ValueError('start from the single-output original task, not an existing chunk layout')
        else:
            argv.append(arg); index += 1
    chunks = []
    for index, bucket in enumerate(buckets):
        name = f'part{index:02}'
        # A conservative chain preserves every prior SCC provider visibility;
        # this is not claimed as a minimally constrained or parallel DAG.
        dependencies = [] if index == 0 else [f'part{index-1:02}']
        argv += ['--chunk', name+':'+str(len(bucket))+(':'+dependencies[0] if dependencies else '')]
        for node in bucket:
            argv += ['--js', records[node]['path']]
        chunks.append({'name': name, 'input_count': len(bucket), 'source_bytes': sum(records[node]['bytes'] for node in bucket),
                       'dependencies': dependencies, 'original_input_indices': bucket})
    argv += ['--chunk_output_path_prefix=chunks/', '--create_source_map=%outname%.map', '--output_manifest=%outname%.inputs', '--output_chunk_dependencies=chunk-dependencies.json']
    artifacts = [f'chunks/{chunk["name"]}.js{suffix}' for chunk in chunks for suffix in ('', '.map', '.inputs')]
    # Stock emits the synthetic weak manifest, but explicitly skips weak JS/maps.
    artifacts += ['chunks/$weak$.js.inputs', 'chunk-dependencies.json']
    result = base | {'name': base['name']+f'-{count}-chunks', 'argv': argv, 'artifacts': artifacts,
                     'absent_artifacts': ['chunks/$weak$.js', 'chunks/$weak$.js.map']}
    result['same_work'] = base['same_work'] | {'chunks': chunks, 'chunk_partition': 'stock AST value+type SCCs, stable topological order, byte-balanced contiguous SCC buckets, conservative chain',
                                              'source_count': len(records), 'original_single_output_source_order': expected,
                                              'status': 'PREPARED ONLY; pinned stock Java4GiB correctness/artifact gate required before Go/timing'}
    result['inputs'] = [source_items[node] for bucket in buckets for node in bucket]+[x for x in base['inputs'] if x['flag'] != '--js']
    return result, {'components': len(groups), 'largest_component_inputs': max(map(len, groups)), 'chunks': chunks}


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument('--base-spec', required=True)
    p.add_argument('--dependencies', required=True)
    p.add_argument('--chunks', type=int, default=32)
    p.add_argument('--out', required=True)
    a = p.parse_args()
    out = Path(a.out).resolve()
    if out.exists() or a.chunks < 2:
        p.error('need a new output path and >=2 chunks')
    base = json.loads(Path(a.base_spec).read_text())
    inventory = json.loads(Path(a.dependencies).read_text())
    result, proof = plan(base, inventory, a.chunks)
    result['partition_provenance'] = {'base_spec_sha256': sha(a.base_spec), 'dependency_inventory_sha256': sha(a.dependencies),
                                      'planner_source_sha256': sha(__file__), 'dependency_helper_source_sha256': sha(HERE/'java/CliDependencyInventory.java')}
    save(out, result); save(out.with_suffix('.partition.json'), proof)
    print(result['name'], len(result['inputs']), 'real inputs; PREPARED, not correctness-gated')


if __name__ == '__main__':
    main()
