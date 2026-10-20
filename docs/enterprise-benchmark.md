# Enterprise benchmark protocol

The comparator is the repository's pinned, unmodified Java jar. Both engines
receive the same complete argv, ordered source files, explicit externs, default
extern source bytes, compilation level, warning/type-check settings, chunks and
output settings. A failure or memory-limit breach is a failed task, never an
infinite speedup. No persistent Go worker, parsed-extern cache or Java class
overlay is used.

The existing tasks already cover substantial single-output compilation:
`typed-full` and `shaka-blockly` use ADVANCED with type checking, while
`bundle-advanced` and `bundle-advanced-ts` use ADVANCED with checkVars. Reuse these
tasks rather than inventing smaller substitutes. Their documentation counts are
approximate: the pinned inputs currently expand to 1,820 files / 25,919,443 bytes
for typed-full and 848 files / 8,523,185 bytes for shaka-blockly.

`enterprise_specs.py` additionally prepares `typed-full-chunks`, using the same
actual inputs in four chunks: core (971), Shaka (522), Blockly (324) and app (3).
The latter three depend on core. This is real multi-chunk ADVANCED work, without
copying or duplicating source. Every chunk's JS, map and input manifest, plus the
chunk dependency graph, is compared. Single-output tasks that emit code also emit a common
input manifest. Checks-only retains the same manifest flag, but stock Java sets
`OutputJs.NONE` and skips manifest emission. Both runtimes must leave that file
absent; stdout, stderr and exit still match exactly. These additional manifest outputs change the measured job from
the historical task; both runtimes perform the extra work.

## Resource budget

Every compiler runs in its own Docker cgroup with kernel `memory.max=4294967296`,
`memory.swap.max=0`, two CPUs of quota and CPU affinity 0,1. Both Docker metadata
and the actual kernel cgroup files are checked. Java uses `-Xmx3g` so a 4 GiB
process budget leaves space for native memory, JIT code, stacks and metadata.
Go uses its default GC/memory/CPU settings; inherited runtime tuning is removed.
The heap flag is disclosed as Java runtime configuration, not claimed as the
process-memory limit. `prlimit --rss` on Linux and Java heap limits alone are
insufficient enforcement.

The pinned Python container image supplies the same userspace for both engines.
The installed JDK is mounted read-only, preserving its actual version and binary
identity. External runtime configuration/data reached through JDK symlinks
(such as Debian security settings and trust store) are mounted individually
read-only, SHA-256 recorded and revalidated between cohorts. The Go binary,
input directories, driver and jar are also read-only.
Network access is disabled during compilation. `validate_budget.py` verifies the
actual 4 GiB kernel settings and a separate 64 MiB allocation canary that must be
OOM-killed. It does not allocate 4 GiB merely to test the limit.

Java is frozen with `docker pause` during each Go invocation, including all JVM
JIT/GC threads. Its resident pages remain charged to its own cgroup, so it can
still contribute to host memory pressure. The host must accommodate both capped
processes plus the coordinator. Outer cgroup memory state is recorded, and an
outer OOM invalidates the cohort. Run benchmarks without builds or other heavy
jobs. The coordinator and Docker transport are outside the compiler budget;
container memory.peak includes the small launcher and charged file cache.

## Paired measurements

`CliPaired` constructs a fresh stock CommandLineRunner and Compiler for every
request, including normal default-extern parsing. It retains the JVM for JIT
warmup but retains no compiler or result. Compiler stdout/stderr are written to
real files inside the timed interval on both engines. Java protocol transport,
artifact hashing and Docker startup/pause overhead are outside the compiler
timings. Go timing includes launching and waiting for each fresh process.

Start with 12 Java warmups, extending to at most 30 according to predeclared
rolling median, spread and CompilationMXBean-time thresholds. Retain every warmup
row. Lack of settling is reported as inconclusive rather than assuming warmed
Java is settled. Three additional fresh Go processes warm executable/filesystem
pages; these are unmeasured and have the same artifact gates.

Within each cohort, individually pair warm Java compilations and fresh Go
processes, alternating AB/BA order with a recorded seed. Restart the JVM between
at least three independent cohorts. Report raw wall/CPU/GC/JIT rows, min/max and
variance, plus elapsed and RSS ratios. Bootstrap confidence intervals resample
independent cohort median ratios, never the correlated requests in one JVM.

Java RSS is the JVM's lifetime VmHWM, including startup, warmup and JIT. Go RSS
is each fresh process's lifetime wait4 peak. These different windows are explicit;
the ratio does not measure per-compilation live heap. Container memory.peak and
memory.events are recorded separately and enforce the complete process budget.

All warmups and measurements must reproduce the fixed stock Java stdout,
stderr, exit and every artifact digest. Inputs and executable identities are
checked between cohorts. The saved spec contains every source hash, exact argv,
chunk graph, warnings/mode and every default extern source hash; checking only
the final JS would miss important differences in work.

## Setup and invocation

Compile `tools/bench/java/CliPaired.java` against the pinned jar in a new driver
directory before measurements. No Java implementation source is required.
Prepare a task with:

```sh
python3 tools/bench/enterprise_specs.py --tasks typed-full,typed-full-chunks \
  --cache /path/to/cache --jar /path/to/pinned.jar --out /path/to/new-specs
```

Then run `paired_cli.py --spec <task.json> --binary <fresh-Go-CLI> --jar <jar>
--jar-sha256 <repository-pin> --driver <compiled-driver-directory> --java-home
<installed-JDK> --image <preinstalled-image@sha256:digest> --out <new-results>`.
No image pull, compiler build or fixture generation occurs during measurements.
`capped_cli.py` uses the identical resource backend for separate correctness and
main-only CPU/allocation profiles; without a stock oracle it labels itself an
ungated diagnostic, not a performance result.

The runner source and image/kernel setup have been prepared. Full protocol,
large-task compatibility, stable warmup and speed claims require actual runs;
preparation and a successful resource canary alone establish none of those.

## Generated-code fixtures

Use actual generated artifacts with documented provenance and runtime checks,
not repeated snippets or fake runtime stubs. The planned families are genuine
J2CL JRE-generated JavaScript with collection semantics, portable Soy-generated
upstream examples with the matching runtime/compiler release, and protobuf
descriptor/well-known-type JavaScript with serialization round trips. Keep these
families separately visible: a JRE fixture is not a large production Java app,
and small generated examples do not become enterprise-scale by duplication.
They remain pending their source/runtime compatibility and stock-Java gates.

`generated_fixtures.py` separates four operations so downloads, extraction and
Soy generation cannot silently enter compiler timings. `fetch --cache <new-asset-cache>`
verifies published immutable release pins; `stage --cache <assets> --out <new-fixture>
--java <installed-java>` materializes the genuine runtimes and runs the pinned Soy
compiler. `specs --fixture <fixture> --jar <pinned-jar> --out <new-specs>` freezes
complete ADVANCED jobs for the paired runner. Use `--families` to select families
for fetch/stage. Input files, protobuf package-resolution metadata, generator
argv/exits, downloaded asset SHA-256 and generated bytes are retained.

The J2CL fixture uses the matching generated JRE and bootstrap artifacts from
Vertispan J2CL `0.11.0-9336533b6`. It exercises collection overwrite, indexing and
removal without fake runtime classes. The protobuf fixture uses actual
`google-protobuf@3.21.4` descriptor, Timestamp and Struct implementations and
checks binary round trips. Soy uses the portable `2019-10-08` compiler and
matching runtime, compiling genuine upstream examples pinned to the release-date
source commit. It checks HTML escaping and retains the generated feature library.
These are generated-code coverage fixtures, not claimed to be large production
J2CL/Soy applications.

Before any valid performance claim, independently run `generated_fixtures.py
verify --spec <spec> --compiled-js <Java-or-Go-output> --node <installed-node>
--out <new-runtime-check.json>` on both compilers' outputs. This check runs
outside the compiler timing and records the actual output/script/Node hashes.
It supplements strict compiler byte gates; it does not replace them. The fixture
builder is source-prepared; successful generation, compiler compatibility and
semantic checks must be established by retained actual runs.
