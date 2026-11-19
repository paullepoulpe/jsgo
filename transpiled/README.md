# Mechanical Closure Compiler Java-to-Go Pipeline (`transpiled/`)

This directory contains the ground-truth-validated, mechanical migration of upstream Google Closure Compiler (`v20261104`, `813` production `.java` files + `2` pure-Java collection helpers = **`815` `.java` files -> `815` `.go` files**, 1:1 file correspondence) to idiomatic, `unsafe`-free Go.

## Pipeline Architecture

1. **`transpiled/prep/` — Non-Algorithmic Java-Prep Patches (`01`–`86`)**:
   - Deterministic, reviewable `.patch` series applied cleanly via `./transpiled/prep/apply_prep.sh` onto pristine upstream `closure-src.jar` (`v20261104`).
   - Strictly non-algorithmic: reflection removal (`args4j`, `Gson`), stream/optional/switch normalization, deterministic collection ordering, primitive/value unboxing, inline small-collection fast paths, inductive bitflag memoization, and reusable per-pass/per-scope state resets.
2. **`transpiled/translator/java/j2go/` — Mechanical Java-to-Go Source Transpiler**:
   - Built strictly on OpenJDK 21's standard `javax.tools.JavaCompiler` and `com.sun.source.*` (`JavacTask`, `Trees`, `Types`, `Elements`) APIs (`100%` copyleft-clean; zero internal compiler imports).
   - Emits `815` `.go` files across 17 acyclic SCC Go packages + `j2gort` in ~7 seconds (`./transpiled/translator/run_j2go.sh`).
   - Lowers `rhino.Node` to a 32-bit `nodestore.Node` handle (`0` = `null`), lowers Java enums without constant-specific overrides to unboxed `uint8`/`uint16` integers, lowers small immutable value classes (`SourcePosition`, `SourceRange`, `parser.Token`, `Property.Key`, `ScopedName`) to stack Go structs, and emits strongly-typed external bindings in `gen/j2gort/`.
3. **`transpiled/translator/passes/` — 11 Idempotent Go-to-Go AST + `go/types` Passes (`gopasses`)**:
   - Modular, independently toggleable passes verified for **100% second-run idempotency** (`Run 2 == 0 sites`) and post-pass `go/types` type-checking:
     1. `native_slice_lowering`: Lowers non-escaping `ArrayList` locals, fields, and slice-accumulator parameters to native `[]T` slices.
     2. `native_set_lowering`: Lowers non-escaping `HashSet`/`LinkedHashSet` locals, fields, and `Set`-returning methods to `j2gort.OrderedSet[T]` / `PlainSet[T]`.
     3. `native_map_lowering`: Lowers non-escaping `HashMap`/`LinkedHashMap`/`IdentityHashMap` locals and fields to `j2gort.OrderedMap[K, V]`, `CustomOrderedMap[K, V]`, or `PlainMap[K, V]`.
     4. `immutable_list_lowering`: Lowers `ImmutableList` fields/locals/getters and `ImmutableList.Builder` locals to native `[]T` slices.
     5. `native_deque_lowering`: Lowers `ArrayDeque`/`Deque`/`Queue` locals and fields to ring-buffer `j2gort.Deque[T]`.
     6. `string_builder_lowering`: Lowers non-escaping `StringBuilder` locals to stack-buffered `j2gort.LocalStringBuilder`.
     7. `slab_alloc_lowering`: Lowers compilation-scoped/pass-scoped AST, type, scope, `Var`, CFG/DFA, and pass state allocations onto contiguous `j2gort.Slab[T]` / `PooledSlab[T]` arenas with deterministic `ReleaseSlabs_()` lifecycle hooks, plus hot field cache-line packing.
     8. `node_traversal_lowering`: Pools `jscomp.NodeTraversal`, inlines `[12]`-deep scope stacks, and reuses `SyntacticScopeCreator`.
     9. `dead_cast_cleanup`: Strips redundant self-conversions (`T(x)`, `AsIface[T](x)`).
     10. `abstract_base_devirt`: Devirtualizes non-overridden abstract-class interface calls (`JSType`, `ObjectType`, `AbstractScope`, `AbstractVar`, `CompilerInput`, etc.) directly to `*<I>Base` struct methods/fields.
     11. `j2gort_peephole`: Eliminates redundant `Cast*`/`AsIface`/`Box*`/`StrVal`/`CheckNotNull*`/`IfaceEquals` wrappers, unwraps zero-arg statement/call-arg IIFEs, and lowers `Node` accessors to direct `nodestore` intrinsics.
4. **`transpiled/shims/` — Clean-Room Go Runtime Shims (`0` `"unsafe"` Imports)**:
   - 12 minimal Go packages (`nodestore`, `javalang`, `javaio`, `javanio`, `javautil`, `javazip`, `javamath`, `javaregex`, `guava`, `gson`, `protobuf`, `resources`) documented in [`shims/API_SURFACE.md`](shims/API_SURFACE.md).

## Reproducing Full Regeneration & Verification

```bash
export GOFLAGS="-p=2"

# 1. Verify and apply all 86 Java-prep patches onto .cache/closure-src:
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

## Differential Parity & `googleapis-bench` Performance (`Java ÷ Go`)

All waves maintain **100% byte-for-byte output and state parity** against upstream Closure Compiler Java `v20261104`:
- `checks.state` SHA-256: `8e2b94eb327ca0545d398573fd654cc8dc2c735db39301674adb43686c29f2fc`
- `optimizations.state` SHA-256: `de254875c355e1db7d3da016d59c43f3c4b3c3d154f8f74aff7da4d3ca48b2e1`
- `7,460` output chunks SHA-256: `7968e63b5cd96980ce477dea033cb4f3fbaef3268e622595dca3198a16fc1467`
- Constraints: **Zero algorithmic compiler changes**, **`-pgo=off`**, **default Go GC (`GOGC`/`GOMEMLIMIT` unset)**, **`0` `"unsafe"` imports**, `--num_parallel_threads=1`.

| Milestone / Wave | `CHECKS` Wall | `OPT` Wall | `FIN` Wall | **Total Wall** | **`Java ÷ Go` Wall** | **Total CPU (`Java ÷ Go`)** | **Peak RSS** |
|---|---:|---:|---:|---:|---:|---:|---:|
| **Pre-Flip Cold Java (`v20261104`, `-t1`)** | `124.17s` | `92.81s` | `20.16s` | **`237.14s`** | `1.00×` | `397.93s` (`1.00×`) | `10,216 MiB` |
| Wave 35 (`v34c`, `35978368`) | `164.59s` | `118.42s` | `25.92s` | **`308.93s`** | `0.77×` | `314.80s` (`1.26×`) | `8,412 MiB` |
| Wave 39 (`v38`, `00dc053e`) | `144.88s` | `101.40s` | `23.15s` | **`269.43s`** | `0.88×` | `271.10s` (`1.47×`) | `7,920 MiB` |
| Wave 40 (`v39`, `69b47e35`) | `139.61s` | `100.12s` | `22.78s` | **`262.51s`** | `0.90×` | `264.15s` (`1.51×`) | `7,850 MiB` |
| Wave 41 (`v40`, `2fbb2f65`) | `137.03s` | `97.97s` | `22.54s` | **`257.54s`** | `0.92×` | `258.90s` (`1.54×`) | `7,810 MiB` |
| **Wave 42 (`v41`, `c30f865e`)** | **`128.99s`** | **`93.28s`** | **`21.67s`** | **`243.94s`** | **`0.97×`** | **`244.98s` (`1.62×`)** | **`7,775 MiB`** |
| Wave 43 (`v42b`, in verification) | *pending* | *pending* | *pending* | *pending* | *pending* | *pending* | `6,344 MiB` (`OPT`) |
