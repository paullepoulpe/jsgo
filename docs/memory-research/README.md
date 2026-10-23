# Memory-layout research (2026-10-23)

Six research agents studied how other parsers, compilers, optimizers and type
checkers keep memory and GC cost low, and mapped each technique onto jsgo's
data structures. `context.md` is the brief they all received (jsgo's layouts
and heap numbers at commit 8ae12b3). Nothing here was built from the other
projects' code; the reports describe designs and link to their sources.

| Report | Projects |
|---|---|
| [go-js-toolchains.md](go-js-toolchains.md) | typescript-go, esbuild, gopls, cmd/compile |
| [rust-zig-js-toolchains.md](rust-zig-js-toolchains.md) | oxc, swc, Bun (Zig), Biome/rowan |
| [data-oriented-compilers.md](data-oriented-compilers.md) | Zig, Carbon, rustc, Go cmd/compile, LLVM idioms |
| [engines-and-managed-compilers.md](engines-and-managed-compilers.md) | V8, SpiderMonkey, Hermes, QuickJS, Closure on the JVM, Roslyn |
| [go-gc-techniques.md](go-gc-techniques.md) | Go runtime (malloc headers, Green Tea GC, maps), memory-heavy Go systems |
| [type-checkers.md](type-checkers.md) | Sorbet, TypeScript checker, Flow, Hack, Pyright/mypy, Closure Colors/TypePool |

## Findings that were verified and acted on

- **Malloc-header size-class waste.** Since Go 1.22 a pointerful object of
  513–32,760 B carries an 8-byte header, so power-of-two slab chunks land in the
  next size class: a 128-node `[]rhino.Node` chunk costs 74 B per 64-B node
  (measured). Single objects of ≤ 512 B have no header and are freed one by one.
- **Slabs pin garbage.** A live element keeps its whole chunk alive; dead AST
  nodes that still carry a JSType then pin the whole type registry (via the
  per-type `registry` pointer) after `typesToColors`.
- **`propertySymbolList`** (a package-level `[]*KnownSymbolType`) kept the
  JSTypeRegistry reachable after it was cleared.
- **Go 1.27.3 vs 1.24.7** (same code, Green Tea GC, AVX-512 machine): peak RSS
  −3 to −9 %, GC CPU −11 to −35 % on three of four tasks (+9 % on checks-only;
  single runs on a loaded machine).

## Consensus direction

The AST is already about as large as Java's (Java pays 24 B per property
item); what costs jsgo is that Go's non-generational GC re-traces every node
pointer on every cycle. Types are 2–3× Java's size. So: make the long-lived
AST pointer-free (int32 links, type/color handles) with explicit reclamation
(copying compaction at safe points), shrink the type structs, store positions
as offsets with line tables, and never let per-pass arenas outlive their pass.
