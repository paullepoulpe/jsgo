# Parsing 100K+ inputs: where PrebuildAst's time goes (2026-10-26)

Builds with 100K-200K input scripts spend most of their time parsing.
This note measures the parse phase at that scale and ranks what to change.

## How it was measured

- **Workload**: N copies of typed-mega's 4,385 Closure files (OpenSphere,
  OpenLayers, Shaka, Blockly, Closure Library...; 9.7 KB per file on
  average), each copy in its own directory with every `goog.provide` /
  `goog.module` / `goog.require` namespace prefixed (`c7.a.b`), so each copy
  is a separate dependency graph. typed-mega's flags (ADVANCED, checkTypes,
  `--dependency_mode=SORT_ONLY`), inputs passed with `--flagfile`.
- **Machine**: 4 vCPUs (Xeon 2.8 GHz), 16 GB. Contention effects grow with
  the core count, so the numbers below understate them on 32-96 core build
  machines.
- **Instrumentation** (not committed): a timestamp after each step of
  `ParseInputs`, exit right after `parseInputs`, CPU profiles.
- `--num_parallel_threads` defaults to 1, as in Java: without it the parse
  is fully sequential (52,620 files: 80.8 s, against 20.9 s with 4 threads).

## Results

Parse phase, 4 threads:

| files | wall | parallel parse | serial steps | peak RSS |
|---|---|---|---|---|
| 4,385 | 1.8 s | 1.5 s | 0.3 s | |
| 13,155 | 5.6 s | 5.0 s | 0.6 s | 0.9 GB |
| 26,310 | 10.1 s | 8.9 s | 1.2 s | 1.8 GB |
| 52,620 | 20.9 s | 18.4 s | 2.5 s | 3.5 GB |
| 61,390 | 25.2 s | 22.4 s | 2.8 s | 4.3 GB |
| **100,855** | **116 s** | **108 s** | ~5 s | **9.2 GB** |

Linear up to 61K files, then three times slower than linear (about 41 s
expected at 100K) with 2 GB more memory.

### 1. A cap at 65,534 source files

Each node keeps its source file as a 16-bit id in its property word
(`nodePropSourceFileMaxID`). Past 65,534 files the id overflows, and every
node of every later file falls back to the general property list: a
heap-allocated ext box plus an entry in a locked side map, for each node.
At 100,855 files a third of the code takes that path:
`copyTemplateColdProps` is 79 s of the 323 CPU-seconds and GC marking 43 s
(4.4 s at 52K files), plus 72 s of system time for the memory.

### 2. Threads interfere on node allocation

From 2 to 4 threads (52,620 files) the parse goes from 26.6 to 18.4 s:
1.45x for twice the threads, with 10 CPU-seconds more for the same work.
The growth is in:

- node allocation and node record writes (~3.5 s): every node takes an id
  from one shared atomic counter, and consecutive ids go to different
  workers, so two workers' 32-byte records share a cache line and every
  write moves it between cores (`NodeStore.alloc`, `setProps`,
  `AddChildToBack`, the `Init*FromTemplate` constructors);
- string pool lookups (~0.7 s): each lookup read-locks one of 64 shards,
  which writes the shard's reader count; frequent identifiers (`goog`,
  `this`, `prototype`) hit the same shards from every worker;
- heap allocation (~1.5 s).

These are the costs that grow with the core count.

### 3. Where the CPU goes (52,620 files, 4 threads)

| | share |
|---|---|
| parsing (parser 23 %, IR conversion 39 % of which JSDoc 10 %) | 62 % |
| heap allocation | 9 % |
| string pool interning | 8 % |
| GC marking | 6 % |
| reading and decoding files | 5 % |

### 4. Serial steps around the parse

They grow linearly, about 4 s at 100K files: 10 % of the phase on 4
cores, most of it on 64 (the parse itself would take ~2.5 s there).

| step | 52,620 files | 100K (est.) |
|---|---|---|
| `findModulesFromEntryPoints` after the parse: a `FindModuleDependencies` traversal of every AST on the main thread | 0.67 s | 1.3 s |
| `ManageDependencies` (dependency sort; string-keyed maps) | 0.36 s | 0.75 s |
| setup before parsing (flags, creating 100K inputs, file stats) | 0.38 s | 0.7 s |
| dependency-info setup, publish loop, dependency walks, attaching scripts | 0.54 s | 1.0 s |

## Recommendations, by expected gain

1. **Lift the 65,534-file cap.** Keep the 16-bit id for the first 65K
   files and put larger ids in a 32-bit column allocated page by page (as
   the original-name column does): no cost below the cap, 4 bytes per node
   above, no limit. Alternative: widen the inline field by taking 2-3 bits
   from the 16-bit length field (lengths of 16K/8K and more would then use
   the ext box, as lengths of 64K do today; few nodes), which caps at 262K
   or 524K files. Expected: 100K files in ~41 s instead of 116 s on 4
   cores, ~7 GB instead of 9.2 GB. Contained change in `internal/rhino`
   (property word, template copies, clone, compaction).
2. **Give each worker its own range of node ids.** A worker claims a block
   of ids (for example one 16K-node block) with one atomic add and
   allocates from it: no shared counter, no cache lines shared between
   workers, and each file's nodes become contiguous, which also helps the
   passes that walk them later. Needs an allocator handle threaded through
   the IR factory and the JSDoc parser. Node ids already differ between
   runs with parallel parsing, and the output does not depend on them
   (`tools/bench/parallel_gate.py`; typed-app is byte-identical at 1 and 4
   threads).
3. **Reuse parser memory within a worker.** Parse trees are about three
   times the size of the AST and die as soon as their file is converted,
   but each file allocates fresh slabs. Resetting a worker's slabs for its
   next file removes most of the parse-time allocation (allocation and GC
   are 15 % of the CPU here, and they contend between threads).
4. **String pool lookups without shared writes**: a small per-worker cache
   in front of the pool, or lock-free reads (the shard table published
   atomically, rebuilt on growth).
5. **Parallelize or trim the serial steps** (Amdahl): run
   `FindModuleDependencies` per input in the parse workers, buffering
   diagnostics and publishing in input order as the dependency walks
   already do; make the dependency sort index-based instead of string-keyed;
   fold the per-input publishing into the workers.
6. Smaller: an ASCII fast path when decoding sources (reading and decoding
   is 5 %).

Measure 2-5 on the target build machines: a 4-vCPU box shows the
direction, not the size, of contention costs.

## Reproducing

`gen.py N` (below) writes N renamed copies of typed-mega's inputs under
`ws/c<k>/` next to it, and `ws/flags-N.txt`. It needs the bench inputs
(`tools/bench/fetch-inputs.sh`) and about 1.2 GB of disk for N=23. Then, in
`ws/`:

    jsgo --flagfile flags-23.txt --num_parallel_threads=4 --js_output_file=/dev/null

The step timings came from timestamps added in `ParseInputs`, and the
"exit after parsing" from a `DebugAfterPassHook` that calls `os.Exit` after
`parseInputs` (with `JSGO_CPU_PROFILE` stopped first).

```python
#!/usr/bin/env python3
"""gen.py N: N renamed copies of typed-mega's inputs in $S/many/ws/c<k>/, and a
flag file $S/many/ws/flags-N.txt (typed-mega's flags + every input)."""
import os, re, sys
sys.path.insert(0, '/home/user/jsgo/tools/bench')
import bench
N = int(sys.argv[1])
here = os.path.dirname(os.path.abspath(__file__))
ws = os.path.join(here, 'ws')
bench.make_workspace(ws, '/home/user/jsgo/.cache')
t = bench.load_tasks(['typed-mega'])[0]
args = bench.expand_inputs(t['inputs'], ws)
files = [args[i + 1] for i in range(0, len(args), 2)]
ns = re.compile(r"""\b(goog\.(?:provide|module|require|requireType|forwardDeclare|module\.get|declareModuleId))\(\s*(['"])([\w.$]+)\2""")
out = []
for k in range(N):
    pre = f'c{k}.'
    for f in files:
        dst = os.path.join(ws, f'c{k}', f)
        if not os.path.exists(dst):
            src = open(os.path.join(ws, f), encoding='utf-8', errors='surrogateescape').read()
            os.makedirs(os.path.dirname(dst), exist_ok=True)
            with open(dst, 'w', encoding='utf-8', errors='surrogateescape') as w:
                w.write(ns.sub(lambda m: f'{m.group(1)}({m.group(2)}{pre}{m.group(3)}{m.group(2)}', src))
        out += ['--js', f'c{k}/{f}']
with open(os.path.join(ws, f'flags-{N}.txt'), 'w') as w:
    w.write('\n'.join(t['flags'] + out) + '\n')
print(f'{N} copies: {len(files) * N} files')
```
