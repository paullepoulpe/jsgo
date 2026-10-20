# Performance experiments, October 2026

The retained CLI improves every workload in the existing five-task suite.
Warm Java remains faster on those full builds. The persistent worker beats
warm Java on the held-out small WHITESPACE workload; that is a specific
workflow result, not a claim that Go generally outperforms HotSpot.

## Full compilations

Five interleaved runs per Go binary; medians below. Original is commit
`c9872522` built with Go 1.24.7. Candidate uses Go 1.27.2 and the repository's
existing PGO profile. Both use default runtime GC settings. Each invocation
recompiles the program and starts a new Go process.

| Workload | Original Go | Candidate Go | Less elapsed time | Fresh warm Java |
|---|---:|---:|---:|---:|
| typed-app | 6.094s | 5.290s | 13.2% | 3.687s |
| bundle | 23.102s | 18.531s | 19.8% | 16.568s |
| downlevel-whitespace | 1.554s | 1.418s | 8.8% | 0.933s |
| downlevel-simple | 4.762s | 4.038s | 15.2% | 2.660s |
| checks-only | 8.872s | 8.028s | 9.5% | 5.930s |

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

## Repeated small applications

These small DOM-counter workloads were held out of PGO training. Java uses
100 warm-up and 30 measured iterations. Go uses the production `jsgo-worker`
JSON protocol, 10 warm-up and 30 measured requests, including request/response
transport and client decoding. All 40 responses and generated artifacts are
checked against Java for each workload.

| Workload | Go worker | Warm Java | Java time / Go time |
|---|---:|---:|---:|
| small-advanced | 948.45ms | 546.45ms | 0.58x |
| small-simple | 286.59ms | 232.90ms | 0.81x |
| small-whitespace | 45.64ms | 57.90ms | 1.27x |
| small-checks | 516.47ms | 357.35ms | 0.69x |

Only the WHITESPACE row is a win over Java. User code is reread and compiled on
every request; this is not compiled-output replay. The worker clones immutable
built-in extern parse templates and creates a new compiler for each request.
See [worker usage](worker.md).

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
- Add a serialized persistent CLI worker with fresh compiler state and an
  immutable, bounded built-in extern parse cache. Clone JSDoc type expressions
  and transformations, preserve inline-type bits, and rebind source metadata.
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
axios, ajv and Monaco inputs. Thirteen mixed worker requests also matched Java
and the ordinary CLI: source changes, output maps, warnings, parse errors, bad
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
records a limited worker win and measured CLI improvements without redefining
the JVM baseline to make the other rows pass.

Raw medians, iteration data, binary hashes, build settings and machine metadata
are in [performance-results.json](performance-results.json).

## Reproduce

Build the baseline and candidate separately, then use prebuilt binaries during
measurement. The experiment harness refuses any run that changes Java's output.
Use `tools/bench/bench.py` to establish a fresh Java reference, and
`tools/bench/experiments.py` for interleaved Go variants. The worker benchmark
uses the production JSON request/response protocol. To reproduce the small-app
comparison separately, first prepare its warm Java reference and then run the
protocol harness below. Pin both commands to the same CPU set (for example
`taskset -c 0,1` on Linux) and keep heavy jobs off those CPUs.

```sh
go build -o bin/jsgo ./cmd/jsgo
go build -pgo=cmd/jsgo/default.pgo -o bin/jsgo-worker ./cmd/jsgo-worker
python3 -I tools/bench/bench.py run --jsgo bin/jsgo --out /tmp/jsgo-reference \
  --tasks typed-app,bundle,downlevel-whitespace,downlevel-simple,checks-only \
  --runs 5 --warmup 12 --warm-iters 10 --skip-trace

python3 -I tools/bench/bench.py run --skip-jsgo --out /tmp/jsgo-small-java \
  --tasks small-advanced,small-simple,small-whitespace,small-checks \
  --runs 1 --warmup 100 --warm-iters 30 --skip-trace
python3 -I tools/bench/worker_protocol.py --binary bin/jsgo-worker \
  --reference /tmp/jsgo-small-java --out /tmp/jsgo-small-worker
```
