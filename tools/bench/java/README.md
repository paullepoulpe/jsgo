# Equivalent persistent workers

`CliWorker` uses the same JSON/base64 request/response protocol as `jsgo-worker`.
Each request creates a fresh command-line runner and compiler; compilation
requests remain serialized. Both workers may reuse immutable parsed builtin
extern templates across requests. User sources and custom externs are reparsed.

Build the Java comparison worker with the pinned jar and source archive in the
benchmark cache:

```sh
export JSGO_CACHE="$(pwd)/.cache"
python3 tools/bench/build_java_worker.py --out /tmp/java-worker
java -Djsgo.bench.externCache=true -Djsgo.bench.verifyCache=true \
  -cp /tmp/java-worker/classes:$JSGO_CACHE/closure-compiler-v20261015-maven.jar \
  com.google.javascript.jscomp.parsing.BuiltinExternCacheTest
```

The build verifies both upstream SHA-256 pins, compiles into a fresh staging
directory, and writes `build.json` only after compilation succeeds. It leaves
the pinned jar unchanged. Two API overlays expose uncached parsing and a
JSDoc clone preserving inline-type flags and deep-copying type/TTL trees.
They are generated from the pinned Apache-licensed Closure sources; no JDK
implementation source is used. The manifest records source and class hashes.

The Java cache follows Go's eligibility, config equality, full source/name
keys, diagnostic fallback, source rebinding, and limits (four configuration
groups, 512 files per group). Hits deep-clone AST/JSDoc trees. The cache test
compares clones to uncached parsing, mutates every inline/type/TTL tree, checks
fresh source bindings, changes source/config, and verifies diagnostics recur.
`verifyCache` counters are disabled during performance measurements.

Use identical warm-up/measurement counts, affinity, and compiler thread flags
for the two production protocol workers:

```sh
taskset -c 0,1 python3 tools/bench/worker_compare.py \
  --reference /tmp/small-stock-reference --out /tmp/matched-workers \
  --go-binary /tmp/jsgo-worker --java-build /tmp/java-worker \
  --threads 1,2,4 --warmup 100 --iters 30
```

For the five full-size tasks, select them with `--tasks` and use an independently
stabilized window such as `--warmup 12 --iters 10`. References come from
`bench.py run` with the unmodified pinned Java jar. They provide bytes only;
previous timings or memory values are not used. Every request and output file
must match those stock bytes, including during warm-up. The timer covers JSON
encoding, pipe transport, compiler work, and response decoding. File checking is
outside the timer. Reported RSS is the process high-water mark over the same
complete request count for each worker, in MiB. Both ratios are Java divided by
Go: values above one favor Go. Inspect the measured-half drift before treating
a JVM as warmed up.

`--num_parallel_threads` controls compiler phase workers. `--gomaxprocs`, when
specified separately, controls Go's runtime CPU parallelism. The harness records
both settings and inherited CPU affinity rather than conflating them.

A mixed-job correctness gate also covers CLI help/version, changed source,
source maps, diagnostics, failed jobs followed by successful jobs, stdin, JSON
streams, custom/default extern options, charset, and language changes. It runs
one worker through all requested thread counts:

```sh
python3 tools/bench/worker_gate.py --command-json /tmp/java-worker/build.json \
  --out /tmp/java-worker-gate --threads 1,2,4 --require-cache-hits
```

For Go, pass a JSON file containing `["/absolute/path/to/jsgo-worker"]` and omit
`--require-cache-hits`.
