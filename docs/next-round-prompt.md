# Prompt: performance round 4 and benchmark extensions for jsgo

You are working on `jsgo` (`github.com/paullepoulpe/jsgo`), a faithful Go port of
Google Closure Compiler v20261015. Its output is byte-identical with the Java
compiler, and after three optimisation rounds it is roughly level with a warm
JVM single-threaded, and ahead of it with `--num_parallel_threads=4`. Your job:
(A) a fourth optimisation round that makes jsgo reliably faster than warm Java,
and (B) three benchmark extensions. Work on a branch, commit per logical change,
and keep `main` green.

## Read first
- `PORTING.md` (all): the porting rules. Code is a file-for-file, 1:1 port of
  the Java; Go-only optimisations live inside function bodies or in the shared
  substrate (`internal/rhino`, `internal/jutil`, `NodeTraversal`) and are marked
  `// Go-only: performance`.
- `PLAN.md` §5–7, especially §7, "Decided by the repository owner".
- `docs/performance-roadmap.md` (Rounds 1–3), `docs/gclab-2026-10-18.md` (GC
  and heap measurements), `docs/benchmarks.md`, `tools/bench/README.md`.

## Hard rules
1. **Byte-identical output with Java v20261015** (stdout, stderr, exit status,
   every output file and source map) at `--num_parallel_threads=1` and `=4`.
   Never edit a golden by hand to make a test pass. The only accepted
   differences are listed in PLAN.md §7 (JVM stack frames, Java's catchable
   StackOverflowError, non-deterministic Java ordering, gzip bytes compared
   decompressed).
2. **Equivalent work to Java.** No caching or precomputation that Java does not
   do (for example, no pre-parsed externs embedded in the binary, and no
   results carried between invocations). Make the same work cheaper.
3. No GPL code, no JDK source (behaviour comes from Javadoc or specs). Keep the
   licence headers; `make vet` runs `tools/licenses -check`.
4. Never `go test ./...` on a small machine; test packages individually.
   `internal/jscomp` is one long package run (`-timeout 3h`).

## Correctness gate (before and after each optimisation)
Fast, after each change:
- `go build ./... && make vet && gofmt -l .`
- `go test` on the changed packages, plus `go test -race` on those that run
  goroutines.
- The bench gate (`python3 -I tools/bench/bench.py run --runs 1 --skip-warm
  --skip-trace`, or `--skip-java` against stored Java outputs): every task MATCH
  (single-files: the 11 java-frames diffs are accepted).
- `tools/bench/parallel_gate.py` at threads 1, 2 and 4.

Full, once before you start (the baseline: Round 3 was never checked this way)
and once at the end, at threads 1 and 4. Commands are in
`docs/performance-roadmap.md` and `tools/clie2e/main.go`:
- `go run ./tools/clie2e -suite fixtures` (2,974 of 2,976; the other 2 are
  StackOverflowError);
- the held-out sweep: 1,011 files × 7 configs, all identical;
- the flag sweeps: development `@flags` 16,200, held-out 3,340, test262
  7,947, all identical;
- `-suite projects -set all`: 120 of 120;
- `-suite matrix`;
- `go test -count=1 ./internal/jscomp` (and with `-race` once).

If the baseline run finds differences, fix them first, as separate commits with
regression tests.

## A. Optimisations (do all; measure each)
Measure each change in isolation: interleaved A/B, N≥5, same session, PGO on
(see the A/B script pattern in `docs/gclab-2026-10-18.md`). Keep it only if it
is a measured win and the gate passes. Record each result in
`docs/performance-roadmap.md` under "Round 4".

1. **Allocation blocks reused instead of returned to the OS.** In the current
   profile (macOS), `runtime.madvise` is 14% of CPU, plus `mmap`,
   `mspan.init` and `memclr`: the runtime scavenger returns the pages of
   dropped slabs and the next pass faults them back in. Reuse slabs (AST, CFG
   and graph nodes, data-flow lattices, source-map mappings) through per-type
   free lists or `sync.Pool`, reset rather than reallocate, and size them to
   avoid the >32 KB large-object path. Profile on Linux too; the scavenger
   behaves differently there.
2. **Thread handoffs and system calls.** `pthread_cond_wait`/`signal`,
   `usleep` and `pthread_kill` take about 4%, `syscall.rawsyscalln` 3.4%.
   Make parallel work coarser (whole files or large groups per goroutine),
   find the source of the syscalls (unbuffered writes, repeated stat calls)
   and fix it.
3. **Generic types shared across pointer types.** Go compiles one copy of a
   generic type for all pointer types and dispatches through a dictionary, so
   hot types such as `jutil.LinkedHashMap[string,*T]` cannot inline hashing or
   equality. Write hand-specialised versions for the 3–4 hottest key/value
   combinations.
4. **Double hashing and string ids.** The Java-ordered maps compute Java's
   `String.hashCode` and Go's own hash for the same key. Give names a numeric
   id once, in `RhinoStringPool`, and turn the hottest name-keyed maps into
   slices or id-keyed tables. Keep Java's iteration order (it can be
   reproduced from the stored Java hashes).
5. **Visitor interface dispatch.** `NodeTraversal.TraverseBranch` is about 45%
   cumulative; every pass's `ShouldTraverse`/`Visit` is an interface call per
   node. For the 10 hottest passes, add directly typed traversals (generic
   over the concrete callback type, or a token switch) that keep exactly the
   same visit order.
6. **Java-style iterators.** Replace `HasNext`/`Next` iterator objects on hot
   paths with plain loops or Go 1.23 range-over-func.
7. **Copies and boxing that HotSpot's escape analysis removes.** Defensive
   `ImmutableList.copyOf`-style copies where the source is never mutated
   afterwards, and ints or booleans boxed in `any` in node properties. Use
   `go build -gcflags=-m` on hot packages to find heap escapes of small
   temporaries.
8. **Exceptions emulated with panic/recover.** Find any hot path (parser error
   recovery, type-checking bailouts) that uses `panic`/`recover` for normal
   control flow and replace it with returned errors, keeping behaviour
   identical.
9. **Pointer-free AST.** Nodes are slab-allocated but still pointer-linked, so
   the GC scans them (scanning is still about 15% of CPU). Prototype nodes
   linked by int32 indices into pointer-free arenas behind the existing
   `*rhino.Node` API (substrate change only), and measure GC CPU and wall
   time. Land it only if it wins.
10. **Parallelism inside passes that run per function.** Passes that analyse
    each function independently (FlowSensitiveInlineVariables,
    CoalesceVariableNames, DeadAssignmentsElimination and similar) can process
    functions on several goroutines when `--num_parallel_threads>1`, then
    apply their AST changes in the same order as Java. Keep the output
    byte-identical; `parallel_gate.py` and `-race` must pass. Threads=1 stays
    exactly sequential.

## B. Benchmark extensions (exactly these three)
Add them to `tools/bench` (task JSON, pinned inputs in `inputs.txt` with
checksums, fetched by `fetch-inputs.sh`, never committed), with the same
correctness gate as the other tasks:

1. **Large multi-chunk ADVANCED build.** The typed-full input set (Closure
   Library + Shaka Player + Blockly + app, 28.5 MB) compiled in ADVANCED mode
   and split into 50–200 output chunks with `--chunk`, with a realistic chunk
   graph (a shared base chunk, feature chunks depending on it, some on each
   other), source maps on.
2. **Generated-code workloads.** Public stand-ins for code large companies
   mostly generate:
   - J2CL output: J2CL's samples and its JRE emulation library, from the
     j2cl repository's released or prebuilt JavaScript;
   - Closure Templates (Soy) compiled to JavaScript, for example the
     closure-templates test suites or examples;
   - protocol-buffer JavaScript (`protoc --js_out` style) for a non-trivial
     set of `.proto` files.

   Compile them typed in ADVANCED mode. Pin versions and checksums; if
   building an input needs a toolchain, script it and store the generated JS
   checksum.
3. **Memory-capped runs at 4 GB.** Every task also runs with Java at `-Xmx4g`
   and jsgo at `GOMEMLIMIT=4GiB` (GOGC left at its default), reported as its
   own columns, including failures (OOM) if any.

And for all tasks, **report variance with the medians**: min–max (or
interquartile range) and the coefficient of variation for each of Java cold,
Java warm and jsgo. Always interleave Java and jsgo runs in the same session.

## Deliverables
- Commits on your branch, each passing the fast gate; the full gate at the
  end, at threads 1 and 4, with its results table.
- `docs/performance-roadmap.md` "Round 4": each optimisation with its
  measured A/B result (kept or dropped, and why).
- `docs/benchmarks.md`: the full benchmark run (all old and new tasks, N=5,
  Java cold, Java warm and jsgo at threads 1 and 4, the 4 GB columns, with
  variance), on one quiet machine in one session, with the machine
  described.
- `default.pgo` regenerated at the end (`tools/bench/pgo.py`).
- Update PLAN.md §7 if the threading position has changed (it still says
  "no goroutine parallelism").
