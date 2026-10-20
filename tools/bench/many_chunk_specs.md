# Many-chunk preparation

The existing four-chunk case remains useful but is not the many-chunk target.
The new planner targets 32 nonempty chunks from the **same 1,820 real inputs** of
`typed-full`; it does not multiply files, prune sources, replace runtimes or
weaken ADVANCED/type-check settings.

`CliDependencyInventory.java` is an original setup helper in the Closure package,
compiled against the unmodified pinned jar. It invokes stock `CompilerInput` AST
`DepsFinder` before rewrites, with SORT_ONLY and the original ordered input paths.
It uses the actual module loader, forces the AST path rather than a regex scan,
and records provides, value/type requires, input byte hashes, parser errors and
actual Compiler-class jar identity. It is not the measured Java compiler driver
and does not overlay any stock class.

Preparation sequence (not executed during measurements):

1. Freeze `typed-full` single-output spec and write its exact `--js` paths one
   per line in `input-paths.txt`; reject other/missing inputs.
2. Compile the helper in a new class directory, then run
   `com.google.javascript.jscomp.CliDependencyInventory input-paths.txt` with the
   pinned jar first on its classpath, normal input symlinks and the same real
   4 GiB container budget. Save stdout JSON, stderr and exit. This extracts
   metadata only; it establishes no performance or compiler compatibility claim.
3. Run `many_chunk_specs.py --base-spec <typed-full.json> --dependencies <json>
   --chunks 32 --out <new-spec.json>`.
4. Run the **actual new stock Java ADVANCED job under 4 GiB** before any Go job
   or timing. Check success, every emitted chunk/map/manifest and dependency JSON.
   Then perform Go byte/diagnostic/exit parity and only then measurements.

The planner constructs SCCs from value **and type** edges, topologically orders
whole SCCs with deterministic original-input ties, and balances contiguous SCC
buckets by bytes. It emits conservative chain dependencies so every prior
provider remains visible. This intentionally is not a minimally constrained or
maximally parallel DAG. Every input is retained exactly once; SCCs are never
split. Duplicate/unresolved providers, mismatched hashes, extraction errors,
fewer than 32 SCCs or empty chunks fail rather than silently inventing a layout.
If real data triggers one of those cases, preserve the failure and audit it.

All 32 JS files, 32 source maps, 32 input manifests, the synthetic `$weak$`
manifest and chunk graph are required. Stock skips weak JS/map output; their
absence is also gated. The generated spec is explicitly marked **prepared only**
until its actual stock-Java memory and correctness gates pass. The source-only
model tests are not evidence that the real 1,820-input workload succeeds.
