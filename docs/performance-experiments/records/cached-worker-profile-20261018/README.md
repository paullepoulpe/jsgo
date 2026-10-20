Diagnostic profile, not a timing benchmark
========================================

Source: profiling helper 139b0348b6c57834d84abbc168a2f4dffe3d4eb2,
compound base 21c008cd24a7bf973665a44fff85805ebfae853b.
Binary /workspace/jsgo/bin/profileworker SHA256:
dce35142a70d0461e92c9eed94a4300a44b19bc1e68b5bfc4406ebaee75b0b34.
Go 1.27.2, existing cmd/jsgo/default.pgo SHA256:
f6aee9c5761f48f8e96beb0ee7cf5d4d8ef5f631460b2102d785d46885bdeaac.
Build GOGC=30 GOMAXPROCS=1 -p=1 limits affect compilation only. Worker
harness removes runtime tuning variables; runtime GC settings are defaults.

One process, pinned CPUs 0,1, actual production cliworker.Run executor/cache,
100 unprofiled small-whitespace warmups and 100 profiled requests. Every one
of the 200 JSON responses and generated artifacts matched the Java reference.
Source/args are the existing held-out small-whitespace benchmark. CPU sampling
starts after warmup; cumulative allocation snapshots before/after hot requests
are subtracted using pprof -base. Forced GC occurs before each snapshot outside
the CPU window. Other correctness builds ran concurrently. DO NOT reuse the
worker_protocol summary's wall time or Java ratio as a benchmark result.

Hot CPU samples: 7.02 process seconds across 7.77 seconds sampled duration.
cloneBuiltinExternResult accounts for 2.18 seconds (31.05%) cumulative CPU;
its recursive bind function accounts for 1.60 seconds (22.79%). Recursive
CloneTree contributes 0.73 seconds (10.40%), overlapping the total clone cost.
GoroutineID/runtime.Stack contributes another 1.65 seconds (23.50%), mostly
source-file reentrant monitor adapters. Background GC contributes 1.60 seconds
(22.79%). These cumulative shares overlap; do not sum them.

Allocation delta: 2844.32 MiB across 100 hot requests (28.44 MiB/request).
cloneBuiltinExternResult accounts for 2055.94 MiB (72.28%), or about 20.56
MiB/request. Node chunk allocation accounts for 1016.47 MiB (35.74% of total);
JSDocInfo.ToBuilderWithTransform accounts for 659.70 MiB (23.19%). These are
subsets of the clone allocation share and must not be added to it.

This justifies testing a frozen bulk-clone representation. Eliminating all
of the measured direct clone CPU would reduce sampled process CPU by 31%,
corresponding to 1.45x CPU throughput before any GC effects. A real clone must
still allocate and initialize independent mutable state, so the reachable gain
will be lower. This is not a strict wall-time ceiling under parallel GC and
concurrent build contention; only an exclusive unprofiled paired benchmark
can establish a speedup and RSS change.
