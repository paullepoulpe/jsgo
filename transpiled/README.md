# Mechanical Closure Compiler Java-to-Go Pipeline (`transpiled/`)

This directory contains the ground-truth-validated, mechanical migration of upstream Google Closure Compiler (`v20261104`, `813` production `.java` files + `2` pure-Java collection helpers = **`815` `.java` files -> `815` `.go` files**, 1:1 file correspondence) to idiomatic, `unsafe`-free Go.

## Pipeline Architecture

1. **`transpiled/prep/` — Non-Algorithmic Java-Prep Patches (`01`–`91`)**:
   - Deterministic, reviewable `.patch` series applied cleanly via `./transpiled/prep/apply_prep.sh` onto pristine upstream `closure-src.jar` (`v20261104`).
   - Strictly non-algorithmic: reflection removal (`args4j`, `Gson`), stream/optional/switch normalization, deterministic collection ordering, primitive/value unboxing, inline small-collection fast paths, inductive bitflag memoization, and reusable per-pass/per-scope state resets.
2. **`transpiled/translator/java/j2go/` — Mechanical Java-to-Go Source Transpiler**:
   - Built strictly on OpenJDK 21's standard `javax.tools.JavaCompiler` and `com.sun.source.*` (`JavacTask`, `Trees`, `Types`, `Elements`) APIs (`100%` copyleft-clean; zero internal compiler imports).
   - Emits `815` `.go` files across 17 acyclic SCC Go packages + `j2gort` in ~6.8 seconds (`./transpiled/translator/run_j2go.sh`).
   - Lowers `rhino.Node` to a 32-bit `nodestore.Node` handle (`0` = `null`), lowers Java enums without constant-specific overrides to unboxed `uint8`/`uint16` integers, lowers small immutable value classes (`SourcePosition`, `SourceRange`, `parser.Token`, `Property.Key`, `ScopedName`) to stack Go structs, and emits strongly-typed external bindings in `gen/j2gort/`.
3. **`transpiled/translator/passes/` — 11 Idempotent Go-to-Go AST + `go/types` Passes (`gopasses`)**:
   - Modular, independently toggleable passes verified for **100% second-run idempotency** (`Run 2 == 0 sites`, `15,004` sites lowered on Run 1 across 647 files) and post-pass `go/types` type-checking:
     1. `native_slice_lowering`: Lowers non-escaping `ArrayList` locals, fields, and slice-accumulator parameters to native `[]T` slices.
     2. `native_set_lowering`: Lowers non-escaping `HashSet`/`LinkedHashSet` locals, fields (including nullable lazy sets), and `Set`-returning methods to `j2gort.OrderedSet[T]` / `PlainSet[T]`.
     3. `native_map_lowering`: Lowers non-escaping `HashMap`/`LinkedHashMap`/`IdentityHashMap` locals and fields to `j2gort.OrderedMap[K, V]`, `CustomOrderedMap[K, V]`, or `PlainMap[K, V]`.
     4. `immutable_list_lowering`: Lowers `ImmutableList` fields/locals/getters and `ImmutableList.Builder` locals to native `[]T` slices.
     5. `native_deque_lowering`: Lowers `ArrayDeque`/`Deque`/`Queue` locals and fields to ring-buffer `j2gort.Deque[T]`.
     6. `string_builder_lowering`: Lowers non-escaping `StringBuilder` locals to stack-buffered `j2gort.LocalStringBuilder`.
     7. `slab_alloc_lowering`: Lowers compilation-scoped/pass-scoped AST, type, scope, `Var`, CFG/DFA, and pass state allocations onto contiguous `j2gort.Slab[T]` / `PooledSlab[T]` arenas with deterministic `ReleaseSlabs_()` lifecycle hooks, plus Cache-Line-0 hot-field packing.
     8. `node_traversal_lowering`: Pools `jscomp.NodeTraversal`, inlines `[12]`-deep scope stacks, reuses `SyntacticScopeCreator`, and emits block-cached `traverseBranchBlk_` avoiding redundant `nodestore` block loads.
     9. `dead_cast_cleanup`: Strips redundant self-conversions (`T(x)`, `AsIface[T](x)`).
     10. `abstract_base_devirt`: Devirtualizes non-overridden abstract-class interface calls (`JSType`, `ObjectType`, `AbstractScope`, `AbstractVar`, `CompilerInput`, etc.) directly to `*<I>Base` struct methods/fields.
     11. `j2gort_peephole`: Eliminates redundant `Cast*`/`AsIface`/`Box*`/`StrVal`/`CheckNotNull*`/`IfaceEquals` wrappers, unwraps zero-arg statement/call-arg IIFEs, and lowers `Node` accessors to direct `nodestore` intrinsics.
4. **`transpiled/shims/` — Clean-Room Go Runtime Shims (`0` `"unsafe"` Imports)**:
   - 12 minimal Go packages (`nodestore`, `javalang`, `javaio`, `javanio`, `javautil`, `javazip`, `javamath`, `javaregex`, `guava`, `gson`, `protobuf`, `resources`) documented in [`shims/API_SURFACE.md`](shims/API_SURFACE.md).

## Reproducing Full Regeneration & Verification

```bash
export GOFLAGS="-p=2"

# 1. Verify and apply all 91 Java-prep patches onto .cache/closure-src:
./transpiled/prep/apply_prep.sh --check
./transpiled/prep/apply_prep.sh

# 2. Run mechanical j2go transpiler (815 .java -> 815 .go):
./transpiled/translator/run_j2go.sh \
  --src .cache/closure-src \
  --out transpiled/gen \
  --shaded-jar .cache/closure-compiler-v20261104-maven.jar \
  --unshaded-jar .cache/closure-unshaded-deps.jar

# 3. Run idempotent Go-to-Go AST passes (and verify Run 2 == 0 sites):
(cd transpiled && go run ./translator/passes/cmd/gopasses -stats ./gen/...)
(cd transpiled && go run ./translator/passes/cmd/gopasses -stats ./gen/...)

# 4. Static, determinism, unsafe, and license gates:
go run ./tools/maprangecheck
go run ./tools/licenses -check

# 5. Build stock Go binary (-pgo=off, default GOGC/GOMEMLIMIT):
(cd transpiled && go build -p 2 -pgo=off -o bin/jsgo_post ./cmd/jsgo)

# 6. Run single-threaded googleapis-bench (7,460 chunks, 3 stages) and verify SHA-256 parity:
python3 tools/bench/run_googleapis_bench.py --compiler transpiled/bin/jsgo_post --threads 1 --keep-states
```

## Differential Parity & Performance (`Java ÷ Go`)

All milestones maintain **100% byte-for-byte output and state parity** against upstream Closure Compiler Java `v20261104`:
- `googleapis-bench` (`7,460` output chunks, 3 stages):
  - `checks.state` SHA-256: `8e2b94eb327ca0545d398573fd654cc8dc2c735db39301674adb43686c29f2fc`
  - `optimizations.state` SHA-256: `de254875c355e1db7d3da016d59c43f3c4b3c3d154f8f74aff7da4d3ca48b2e1`
  - `7,460` output chunks SHA-256: `7968e63b5cd96980ce477dea033cb4f3fbaef3268e622595dca3198a16fc1467`
- **9-Task Standard Suite (`bench.py`, `--runs 3`, OpenJDK 27)**: **`9/9 MATCH`** on both `jsgo_post` (**`68.82s`** vs **`84.16s`** Cold Java, **`Java ÷ Go = 1.22×`**, beats Cold Java on 8/9 tasks, `1.63×–2.89×` less RSS) and `jsgo_raw` (`84.66s`, `Raw ÷ Post = 1.23×`).
- **`clie2e` Differential Parity**: **`2,974 / 2,974`** non-`StackOverflowError` fixtures identical; **`110`** `projects -set dev` jobs (`0` panics, `0` hangs, 100% JS/sourcemap/diagnostic parity).
- Constraints: **Zero algorithmic compiler changes**, **`-pgo=off`**, **default Go GC (`GOGC`/`GOMEMLIMIT` unset)**, **`0` `"unsafe"` imports**, `--num_parallel_threads=1`.

### `googleapis-bench` Progression (`--threads=1`)

| Milestone / Wave | `CHECKS` Wall | `OPT` Wall | `FIN` Wall | **Total Wall** | **`Java ÷ Go` Wall** | **Total CPU (`Java ÷ Go`)** | **Peak RSS (`Java ÷ Go`)** |
|---|---:|---:|---:|---:|---:|---:|---:|
| **Pre-Flip Cold Java (`v20261104`, `-t1`)** | `124.17s` | `92.81s` | `20.16s` | **`237.14s`** | `1.000×` | `397.93s` (`1.000×`) | `10,216 MiB` (`1.000×`) |
| Wave 35 (`v34c`, `35978368`) | `164.59s` | `118.42s` | `25.92s` | `308.93s` | `0.768×` | `314.80s` (`1.264×`) | `8,412 MiB` (`1.214×`) |
| Wave 39 (`v38`, `00dc053e`) | `144.88s` | `101.40s` | `23.15s` | `269.43s` | `0.880×` | `271.10s` (`1.468×`) | `7,920 MiB` (`1.290×`) |
| Wave 40 (`v39`, `69b47e35`) | `139.61s` | `100.12s` | `22.78s` | `262.51s` | `0.903×` | `264.15s` (`1.506×`) | `7,850 MiB` (`1.301×`) |
| Wave 41 (`v40`, `2fbb2f65`) | `137.03s` | `97.97s` | `22.54s` | `257.54s` | `0.921×` | `258.90s` (`1.537×`) | `7,810 MiB` (`1.308×`) |
| Wave 42 (`v41`, `c30f865e`) | `128.99s` | `93.28s` | `21.67s` | `243.94s` | `0.972×` | `244.98s` (`1.624×`) | `7,775 MiB` (`1.314×`) |
| Wave 44 (`v43`, `4b75e4d0`) | `127.33s` | `93.17s` | `21.45s` | `241.95s` | `0.980×` | `242.85s` (`1.639×`) | `7,042 MiB` (`1.451×`) |
| Wave 45 (`v44`, `f20b1105`) | `126.85s` | `91.67s` | `21.09s` | `239.60s` | `0.990×` | `242.19s` (`1.643×`) | `7,883 MiB` (`1.296×`) |
| Wave 46 (`v45`, `b863c05d`) | `126.77s` | `91.86s` | `21.04s` | `239.66s` | `0.989×` | `241.10s` (`1.650×`) | `7,266 MiB` (`1.406×`) |
| Wave 47 (`v46`, `e311f3d3`) | `123.72s` | **`89.48s`** | **`20.75s`** | **`233.96s`** | **`1.014×`** | **`235.40s` (`1.690×`)** | `7,425 MiB` (`1.376×`) |
| **Wave 47 + Patch 91 (`e0999a69`, Final)** | **`123.56s`** | **`91.09s`** | `20.87s` | **`235.52s`** | **`1.007×`** | **`236.60s` (`1.682×`)** | **`7,188 MiB` (`1.421×`)** |

### 9-Task Standard Suite (`jsgo_post` vs Cold Java OpenJDK 27, 3-Run Medians)

| Task | Cold Java Wall | `jsgo_post` Wall | **`Java ÷ Go` Wall** | Cold Java RSS | `jsgo_post` RSS | **`Java ÷ Go` RSS** | Parity |
|---|---:|---:|---:|---:|---:|---:|:---:|
| `typed-app` | `6.14s` | `3.35s` | **`1.83×`** | `626 MiB` | `306 MiB` | `2.05×` | MATCH |
| `bundle` | `13.2s` | `12.0s` | **`1.09×`** | `1,734 MiB` | `803 MiB` | `2.16×` | MATCH |
| `downlevel-simple` | `4.81s` | `2.74s` | **`1.75×`** | `606 MiB` | `290 MiB` | `2.09×` | MATCH |
| `checks-only` | `6.88s` | `5.53s` | **`1.24×`** | `1,386 MiB` | `510 MiB` | `2.72×` | MATCH |
| `shaka-blockly` | `12.9s` | `8.92s` | **`1.45×`** | `1,514 MiB` | `629 MiB` | `2.41×` | MATCH |
| `typed-full` | `23.2s` | `21.3s` | **`1.09×`** | `3,488 MiB` | `1,209 MiB` | `2.89×` | MATCH |
| `generate-ijs` | `3.17s` | `4.14s` | `0.76×` | `742 MiB` | `330 MiB` | `2.25×` | MATCH |
| `incremental-ijs-check` | `4.10s` | `1.99s` | **`2.06×`** | `709 MiB` | `435 MiB` | `1.63×` | MATCH |
| `conformance-checks` | `9.76s` | `8.84s` | **`1.10×`** | `1,530 MiB` | `630 MiB` | `2.43×` | MATCH |
| **9-Task Sum / Max RSS** | **`84.16s`** | **`68.82s`** | **`1.22×`** | `3,488 MiB` | **`1,209 MiB`** | **`2.89×`** | **`9/9 MATCH`** |
