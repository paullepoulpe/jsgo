# Performance experiments, October 2026

The retained CLI improves every workload in the existing five-task suite.
Warm Java remains faster on those full builds. Worker and parsed-extern-cache experiments have been dropped. The active
comparison uses ordinary Go CLI processes and stock warm Java, creating fresh
compilers and parsing externs on both sides.

## Full compilations

Five interleaved runs per Go binary; medians below. Original is commit
`c9872522` built with Go 1.24.7. Candidate uses Go 1.27.2 and the repository's
existing PGO profile. Both use default runtime GC settings. Each invocation
recompiles the program and starts a new Go process.

| Workload | Original Go (time / RSS) | Candidate Go (time / RSS) | Less elapsed time | Warm Java (time / RSS) | Java/Go time | Java/Go RSS |
|---|---:|---:|---:|---:|---:|---:|
| typed-app | 6.094s / 268.6MiB | 5.290s / 266.5MiB | 13.2% | 3.687s / 1263.8MiB | 0.70x | 4.74x |
| bundle | 23.102s / 675.2MiB | 18.531s / 504.2MiB | 19.8% | 16.568s / 1677.6MiB | 0.89x | 3.33x |
| downlevel-whitespace | 1.554s / 179.3MiB | 1.418s / 165.0MiB | 8.8% | 0.933s / 943.4MiB | 0.66x | 5.72x |
| downlevel-simple | 4.762s / 220.4MiB | 4.038s / 207.0MiB | 15.2% | 2.660s / 1475.5MiB | 0.66x | 7.13x |
| checks-only | 8.872s / 710.0MiB | 8.028s / 679.9MiB | 9.5% | 5.930s / 1486.7MiB | 0.74x | 2.19x |

Both ratio columns divide Java by Go: above 1x means Go takes less time or
uses less peak RSS, respectively. The time ratios compare median elapsed time;
the memory ratios compare the documented process peaks.

The candidate still needs roughly 11–34% less elapsed time, depending on the
workload, to equal these warmed JVM medians. The largest retained gain is on the
large SIMPLE bundle; the smallest is on full-build WHITESPACE.

The JVM column is a separate fresh run with 12 warm-up and 10 measured
iterations per workload, using the pinned Closure v20261015 jar and JDK
21.0.12.1. Both runtimes have the same two-core affinity on an AMD EPYC 9V74
host with a two-CPU container quota. No compilation, test execution or other
heavy task ran alongside timed compilations. CPU frequency is not locked;
small differences between experiment cohorts should not be treated as gains.

Every measured Go invocation must match Java's stdout, stderr, exit status,
JavaScript and source maps. Fresh Java reproduces its original reference.
These are end-to-end results, not sums of isolated optimization percentages.

Memory is peak resident set size (RSS), in MiB (2^20 bytes). Go CLI cells
show the median of five per-invocation peaks. Java cells show the peak of the
entire persistent JVM process, including startup, 12 warm-up iterations and
10 measured iterations. Both are measured with Linux `wait4`/`ru_maxrss`.
These values describe resident process memory, including runtime overhead;
they are not allocated bytes or live-heap size. The bundle's median Go peak
RSS falls from 675.2 to 504.2 MiB (25.3%).

## Retained changes

- Allocate AST nodes and entries of larger ordered maps in batches, reducing
  per-object allocation and scan overhead. Break allocation groups between
  independent compilations to avoid retaining earlier ASTs transitively.
- Use an indexed reaching-definition lattice for functions with at most 64
  canonical variables, preserving TOP/BOTTOM, encounter order and the sparse
  fallback for larger or noncanonical cases.
- Remove the redundant inherited self interface from Var/TypedVar, reducing
  Var from 80 to 64 bytes. Cache immutable name hashes lazily; keep Java's
  name/scope-root equality and full hash formula.
- Store AST type/color and interned string references compactly, reducing
  Node from 96 to 80 bytes. Tagged pointers remain GC-tracked Go pointers.
- Reuse repeated immutable scanner positions. Long-lived comment positions
  retain independent allocation.
- Avoid redundant rename-conflict recording in larger scopes. Small scopes
  retain the cheaper original walk; bookkeeping preserves encounter order.
- Pin Unicode 15 classification/case data from Go 1.24's BSD-licensed package
  so newer Go toolchains retain JDK 21 behavior. No JDK implementation source
  is used.

The normal one-shot CLI, flags and default GC settings remain unchanged.

## Rejected or unretained experiments

| Experiment | Finding |
|---|---|
| Toolchain alone | Roughly 0–3% on the original code; an unmodified upgrade breaks Unicode vectors. |
| Compact parser self references | Smaller parser objects, no useful elapsed-time improvement. |
| Eager variable-name hashing | Pays work for variables never used as hashed keys; lazy hashing was retained instead. |
| Global scope/variable/reference batches | Retained too much intermediate state and increased memory/time. Per-pass boundaries did not make this worthwhile. |
| Concrete scope lookup specialization | Mixed results, no reliable suite-wide gain. |
| Allocation sizes 16/31/63/127 | Confirmed lower allocator waste for several parser records; elapsed results were too small/mixed to retain. |
| Refreshing PGO | Little overall benefit over the existing profile. |
| 64-byte packed AST | Passed pointer checks, clone/original-name/overflow tests and byte gates, but memory savings did not establish a useful elapsed-time gain. |
| More GC headroom / GC disabled | Only a small remaining wall-time gain, sometimes slower. Bundle without GC used about 5.7GB. Defaults were retained. |
| Process-lifetime arenas | Bundle slowed from about 19.1s to 21.4s and RSS grew from about 560MB to 4GB. Rejected. |
| Parallel AST parsing | Targeted cases matched Java under the race detector, but measured gains were small. Default sequential behavior was retained. |

The first worker experiment exposed an allocation-lifetime bug: a partially
filled batch could connect successive compilers through its object graph.
Compiler allocation boundaries and private storage for retained templates fix
that chain. A 22-compilation check stayed near 270–300MB RSS rather than growing
on each iteration. That early leaking run is not used as performance evidence.

## Validation and limits

All packages passed their tests. The large jscomp test binary was compiled
separately with lower *build-time* GC headroom to fit the 8GB machine, then run
with default runtime GC; the test itself passed. Vet and license checks passed.
A compression comparison test now uses an explicit flush to guarantee different
valid deflate layouts, since newer Go versions may emit identical bytes for
`BestSpeed` and `BestCompression`; its payload/header assertions remain intact.

All ten held-out project builds matched Java, covering date-fns, highlight.js,
axios, ajv and Monaco inputs. and the ordinary CLI: source changes, output maps, warnings, parse errors, bad
flags, stdin, JSON streams, custom extern settings and output encoding. JVM
stack-frame differences use the repository's existing normalization rule.

For a fraction `f` of elapsed time accelerated by `k`, the ideal whole-build
speedup is `1 / ((1 - f) + f/k)` with other costs fixed. Even making a 30%
phase twice as fast gives only 1.18x overall. The GC experiments likewise show
that eliminating GC alone cannot remove the remaining large-build gap.

There is no defensible language-level theoretical maximum. HotSpot compiles
hot code to native code too. The remaining large-build gap calls for further
algorithmic work and reducing repeated analysis, rather than assuming a newer
Go compiler, smaller structs or disabling GC guarantees a win. This report
records measured CLI improvements. The worker/extern-cache approach was dropped after the equivalent Java
comparison showed no win; it is outside the active implementation and matrix.

Raw medians, iteration data, binary hashes, build settings and machine metadata
are in [performance-results.json](performance-results.json).

## Reproduce

Build the baseline and candidate separately, then use prebuilt binaries during
measurement. The experiment harness refuses any run that changes Java's output.
Use `tools/bench/bench.py` to establish a fresh Java reference and
`tools/bench/experiments.py` for interleaved Go variants. Use
`tools/bench/thread_sweep.py` for ordinary CLI comparisons at 1, 2, and 4
compiler workers. Pin commands to the same CPU set and keep builds and tests
outside timed runs.

```sh
go build -o bin/jsgo ./cmd/jsgo
python3 -I tools/bench/bench.py run --jsgo bin/jsgo --out /tmp/jsgo-reference \
  --tasks typed-app,bundle,downlevel-whitespace,downlevel-simple,checks-only \
  --runs 5 --warmup 12 --warm-iters 10 --skip-trace
python3 -I tools/bench/thread_sweep.py --binary bin/jsgo \
  --reference /tmp/jsgo-reference --out /tmp/jsgo-threads \
  --threads 1,2,4
```
