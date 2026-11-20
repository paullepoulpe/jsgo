# `j2go` Mechanical Java-to-Go Translator Design

## 1. Overview & Clean-Room Architecture

`j2go` (`transpiled/translator/java/j2go/`) is a deterministic, clean-room, mechanical Java-to-Go source-to-source translator purpose-built for Google Closure Compiler (`v20261104`, 815 prepped non-Ant `.java` files).

### Clean-Room Compliance (`PORTING.md`)
1. **Public Compiler APIs Only**: `j2go` runs on OpenJDK 27 (`javac --release 21`) using **exclusively** standard exported `java.compiler` / `jdk.compiler` public APIs:
   - `javax.tools.JavaCompiler`, `StandardJavaFileManager`, `DiagnosticCollector`
   - `com.sun.source.tree.*` & `com.sun.source.util.{JavacTask, Trees, TreePath, TreePathScanner}`
   - `javax.lang.model.element.*`, `javax.lang.model.type.*`, `javax.lang.model.util.{Elements, Types}`
   - Zero internal `com.sun.tools.javac.*`, `sun.*`, or `jdk.internal.*` APIs are imported or called.
2. **Zero JDK Inspection & Zero Copyleft**: No OpenJDK source code is read, `javap` is never invoked on JDK classes, and `GoNames.extractSanitizedLicenseHeader` automatically strips MPL 1.1 copyleft-alternative clauses from dual-licensed Rhino headers so generated Go files carry only clean Apache-2.0 / BSD / MIT / pure-MPL headers.
3. **1:1 File Mapping**: Every input `<Name>.java` file produces exactly one `<snake_name>.go` file in its target SCC Go package, preserving declaration order and emitting `// Java class: ...`, `// Java method: ...`, and `// Java constructor: ...` provenance comments.

---

## 2. Translator Pipeline & Module Structure

`j2go` consists of 11 self-contained Java classes in `transpiled/translator/java/j2go/` (~10.5K LOC total) invoked via `transpiled/translator/run_j2go.sh`:

| Source File | Responsibility |
| :--- | :--- |
| `Main.java` | CLI entry point: walks `--src` directories, filters excluded paths (`jscomp/ant/`), invokes `ClasspathUnshader`, parses & type-attributes all `.java` files in one `JavacTask`, runs `PackagePlanner` + `TypeResolver`, emits 1:1 `.go` files, and emits `j2gort/`. |
| `ClasspathUnshader.java` | Pure-`java.util.jar` / `DataInputStream` bytecode constant-pool `CONSTANT_Utf8` prefix rewriter that strips `com/google/javascript/jscomp/jarjar/` shading prefixes from `closure-compiler-v20261104-maven.jar` (excluding Closure's own classes so prepped `.java` sources take precedence) for clean `javac` type resolution. |
| `PackagePlanner.java` | Discovers all `TypeElement`s, builds the cross-package symbol dependency graph across all 32 Java packages, runs Tarjan's Strongly Connected Components (SCC) algorithm in deterministic order, assigns acyclic Go package names, resolves intra-SCC simple-name collisions, and tracks per-file Go imports (`ImportTracker`). |
| `GoNames.java` | Deterministic Java -> Go identifier mangling, UpperCamelCase -> `snake_case.go` file naming, Go keyword escaping, WTF-8/UTF-16 string/char literal escaping, and license header sanitization. |
| `TypeResolver.java` | Whole-program type & symbol model: inheritance hierarchy, subclassed-class detection, Union-Find virtual method signature unification across override/interface diamonds, overload disambiguation, static/constructor naming, inner/local/anonymous class outer-instance & transitive variable capture analysis, and external library reference recording. |
| `RuntimeBindings.java` | Maps external JDK / Guava / Gson / Protobuf / Re2j types & methods to `transpiled/shims/*` (`nodestore`, `javalang`, `javaio`, `javaregex`, `resources`, etc.) and generates `<out>/j2gort/runtime.go` + `<out>/j2gort/ext_bindings.go` (pruned strictly to external method signatures actually referenced by the translated program). |
| `EmitContext.java` | Per-file and per-callable emission context: local variable symbol table, captured-variable field bindings, `try`/`catch` control-flow propagation flags (`loopDepthInTryIife`, `switchDepthInTryIife`), switch expression target stack, and deferred local/anonymous class queue. |
| `ClassEmitter.java` | Emits Go interfaces, structs (`<Class>` vs `<Class>Base` + `self_`), `com.google.javascript.rhino.Node` arena handle methods, Java enums (with fields, constructors, and per-constant method override hooks `fn_<Method>`), Java 21 `record`s (canonical/compact constructors, accessors, structural `Equals`/`HashCode`/`String`), and `<clinit>` `sync.Once` initializers. |
| `MethodEmitter.java` | Emits constructors (`New<Class>` + `init<Class>` two-phase allocation/super-chaining), instance/static methods, interface `default`/`private` helpers (`IfaceDefault_<Iface>_<Method>`), covariant/erased parameter unboxing preambles, and `Node` `KindFlags` polymorphic dispatchers. |
| `StmtEmitter.java` | Mechanical statement lowering: JLS §6.3 self-referencing local initializers, `if`/`else` with Java 21 pattern-variable hoisting, classic `for` / `while` / `do-while` / `for-each` (arrays & `Iterable`), labeled loops & labeled blocks, classic `:` switch (with fallthrough & cross-case var hoisting), arrow `->` switch statements/expressions, Java 21 type-pattern + `when`-guarded switches, `synchronized` reentrant monitors, `try`/`catch`/`finally`, and `try-with-resources`. |
| `ExprEmitter.java` | Mechanical expression lowering: JLS §5.6 numeric promotions (`char` `uint16` vs `int32`), wrapping 32/64-bit arithmetic & masked shifts (`javalang`), IEEE-754 `D2I`/`D2L` narrowing, zero-alloc `NullString` sentinel strings, generic field/method/array erasure & checked/unchecked casts, `instanceof` & binding patterns, lambdas, all 5 Java method-reference forms, and JDK/Guava intrinsics. |

---

## 3. Acyclic Go Package Layout (Tarjan SCCs)

Go forbids cyclic package imports, whereas Closure Compiler's 32 Java packages contain two multi-package dependency cycles (`rhino <-> rhino.jstype` and a 15-package `jscomp` core cycle). `PackagePlanner` computes the exact Tarjan SCC condensation DAG in post-order (leaves first), producing **17 acyclic Go packages** (plus `j2gort`) for all **815** prepped `.java` files:

| SCC # | Go Package | Files | Merged Java Packages |
| :---: | :--- | :---: | :--- |
| `00` | `sourcemap` | 19 | `com.google.debugging.sourcemap`, `...sourcemap.proto` |
| `01` | `annotations` | 1 | `com.google.javascript.jscomp.annotations` |
| `02` | `base` | 6 | `com.google.javascript.jscomp.base` |
| `03` | `diagnostic` | 4 | `com.google.javascript.jscomp.diagnostic` |
| `04` | `colors` | 6 | `com.google.javascript.jscomp.colors` |
| `05` | `parser` | 101 | `...jscomp.parsing.parser`, `...parser.trees`, `...parser.util` |
| `06` | `rhino` | 83 | `com.google.javascript.rhino`, `com.google.javascript.rhino.jstype` |
| `07` | `dtoa` | 1 | `com.google.javascript.rhino.dtoa` |
| `08` | `parsing` | 12 | `com.google.javascript.jscomp.parsing` |
| `09` | `graph` | 20 | `com.google.javascript.jscomp.graph` |
| `10` | `resources` | 2 | `com.google.javascript.jscomp.resources` |
| `11` | `regex` | 3 | `com.google.javascript.rhino.typed_ast` / `jscomp.regex` |
| `12` | `jscomptype` | 6 | `com.google.javascript.jscomp.type` |
| `13` | `jscomp` | 539 | `jscomp` + `bundle`, ` coding`, `conformance`, `crosscheck`, `deps`, `disambiguate`, `ijs`, `instrumentation`, ` isolates`, `j2clbuild`, `lint`, `modules`, `polyfill`, `prep`, `serialization`, `transpile`, `verformance` (15 cyclic Java packages) |
| `14` | `baseformat` | 8 | `com.google.javascript.jscomp.base.format` |
| `15` | `reporter` | 3 | `com.google.javascript.jscomp.deps.reporter` / error reporting leaf |
| `16` | `examples` | 1 | `com.google.javascript.jscomp.examples` |

### Intra-SCC & Cross-Package Name Disambiguation
- **File names**: When two classes in different Java packages within the same SCC share a `<snake_name>.go` basename (4 collisions across all 815 files, e.g., `deps/ModuleLoader.java` vs `modules/ModuleLoader.java`), `PackagePlanner` prefixes the relative subpackage (`deps_module_loader.go`, `bundle_Compiler.go`, `transpile_Transpiler.go`).
- **Type names**: Top-level types keep their clean `ExportedSimpleName` unless two top-level classes in the same SCC collide, in which case deterministic subpackage prefixes (`Deps_`, `Bundle_`, `Transpile_`) are prepended. Nested types flatten as `Outer_Inner` (disambiguated as `Outer__Inner` if a top-level `Outer_Inner` exists); local/anonymous classes flatten as `EnclosingTop__Local_<N>_<Name>` and `EnclosingTop__Anon_<N>`.
- **Import aliases**: Cross-package SCC imports always alias as `pkg_<gopkg>` (e.g., `pkg_colors`, `pkg_graph`, `pkg_rhino`), preventing any Java local variable or parameter named `colors`, `graph`, `parser`, or `rhino` from ever shadowing a Go package identifier.

---

## 4. Mechanical Java-to-Go Lowering Rules

### 4.1 Class Hierarchy, Open Recursion (`self_`), and Interfaces
- **Leaf concrete classes** (never subclassed in the source program, ~85% of classes):
  - Emit a single concrete struct `type Foo struct { ... }` and pointer value type `*Foo` with zero interface wrapper overhead.
- **Subclassed classes & abstract classes** (`resolver.isSubclassed(te) == true`):
  - Emit both:
    1. A Go interface `type Base interface { GetBaseBase() *BaseBase; Implements_<pkg>_Base(); <all virtual methods>... }` representing polymorphic references to `Base` or any subclass.
    2. A base struct `type BaseBase struct { self_ Base; <fields>... }` implementing `Base`.
  - Concrete subclasses embed `BaseBase` (e.g., `type Sub struct { BaseBase; <sub fields>... }`).
  - During `initBase(this_ *BaseBase, self_ Base, ...)`, `this_.self_ = self_` wires open recursion: whenever an inherited method on `*BaseBase` evaluates `this` (`ctx.selfValueExpr() -> this_.self_`) or invokes a virtual method on `this`, Go dispatches dynamically to the most-derived concrete override!
  - Explicit `super.foo(...)` calls invoke `this_.BaseBase.Foo(...)` directly on the embedded struct, bypassing virtual dispatch just like JVM `invokespecial`.
- **Interfaces**:
  - Emit Go interfaces embedding superinterfaces (`Implements_<pkg>_<Iface>()` marker methods preserve nominal `instanceof` / checked casts even on zero-method marker interfaces).
  - Interface `default` and `private` methods emit package-level functions `IfaceDefault_<Iface>_<Method>(this_ <Iface>, ...)` and thin forwarding methods on implementing structs that do not override them.

### 4.2 Virtual Method Unification, Overloads, and Generics
- **Union-Find Signature Unification (`TypeResolver`)**:
  - Java allows covariant return types and generic specialization across overrides (`T get()` -> `String get()`), whereas Go requires identical method signatures to satisfy an interface.
  - `TypeResolver` groups every override/interface-implementation family via Union-Find into a canonical `MethodSig`: if all members erase to the exact same Go parameter/return type, that concrete Go type is kept (zero `any` boxing); if two members erase to distinct Go types (covariant return or generic bridge), that slot widens to `any` (`paramWidenedToAny[i]` / `returnWidenedToAny`) and `MethodEmitter` / `ExprEmitter` automatically insert checked/unchecked conversions at entry and call sites.
- **Overload Disambiguation**:
  - Within each Union-Find root family on a type hierarchy, the primary arity-0 / most-common overload keeps clean `MethodName`; additional overloads receive deterministic arity + erased-parameter suffixes (`Foo_1_int32`, `Foo_2_Node_string`).
- **Two-Phase Constructors**:
  - Every constructor lowers to an allocator `New<Class>_*(...) *<Class>` plus an initializer `init<Class>_*(this_ *<Class>, ...)`, allowing `this(...)` and `super(...)` constructor chaining, pre-super argument evaluation, and exact JLS §12.5 instance field initialization order (`initFields_<Class>` runs right after `super(...)` returns and before the constructor body).

### 4.3 Enums & Java 21 Records
- **Enums**:
  - Embed `j2gort.EnumBase` (`Name_ string`, `Ordinal_ int32`), initialized once via `clinit_<Enum>()` (`sync.Once`), with `<Enum>_Values()` and `<Enum>_ValueOf(string)` (note: custom overloads like `JSType.MatchStatus.valueOf(boolean)` disambiguate cleanly to `JSType_MatchStatus_ValueOf_1_bool`).
  - Enum constant-specific class bodies (`CONST { @Override boolean apply(Node n) { ... } }`) compile to function pointers `fn_<Method> func(*<Enum>, ...) ...` on the enum struct, supporting both enum-declared abstract/virtual methods and inherited interface methods (`Predicate.apply`, `Function.apply`) plus `super.<method>(...)` via `baseEnum_<Method>`.
- **Java 21 `record`s**:
  - Emit immutable structs with component fields, canonical/compact constructor support, component accessor methods, and structural `Equals`, `HashCode` (matching `Float`/`Double`/`Arrays`/reference semantics), and `String`.

### 4.4 Inner, Local, and Anonymous Classes
- Non-static inner classes store an `outer_ <EnclosingGoType>` field threaded automatically as the first parameter of `New<Inner>` / `init<Inner>`. `Outer.this` and unqualified outer field/method references walk the `outer_` chain (`resolveImplicitThisValue`).
- Local and anonymous classes declared inside methods/lambdas are queued on `EmitContext.deferredLocalTypes` and emitted at top level after the enclosing callable, capturing transitive effectively-final outer local variables into `cap_<var>` struct fields initialized at the `New...` allocation site.

### 4.5 Strings (`NullString`), Primitives (`uint16` `char`), and Operators
- **Zero-Alloc Nullable Strings**:
  - Java `String` maps directly to Go value type `string` (zero `*string` heap pointer indirection!).
  - Java `(String) null` is represented by `j2gort.NullString = "\xff\x00j2go_null"` (an impossible UTF-8/WTF-8 byte sequence). `BoxString` / `CastString` / `RefEquals` / `StringEquals` / `CheckNotNull` convert seamlessly between `NullString` and `any(nil)` at generic/object boundaries.
- **Numeric & Shift Semantics**:
  - `byte -> int8`, `short -> int16`, `char -> uint16`, `int -> int32`, `long -> int64`, `float -> float32`, `double -> float64`.
  - Binary `+`, `-`, `*` on `int32`/`int64` wrap in two's complement; `/`, `%` use `javalang.Div32/Rem32/Div64/Rem64` (`MIN_VALUE / -1` overflow guard + `ArithmeticException` on `/ 0`); shifts mask shift counts by `& 0x1f` / `& 0x3f` (`javalang.Shl32/Shr32/Ushr32/...`); float-to-int casts use saturating NaN-safe `javalang.D2I` / `javalang.D2L`.

### 4.6 Control Flow, Exceptions, Switch Expressions, and Lambdas
- **Exceptions (`try`/`catch`/`finally` & `try-with-resources`)**:
  - `throw ex` lowers to `panic(j2gort.CheckNotNull(ex))`.
  - Only methods/blocks containing `try` incur a `func() { defer ... }()` IIFE wrapper; `recover()` uses `javalang.RecoverException` (re-panicking non-Java Go runtime panics unless matched, and running `finally` `defer`s in exact LIFO order after `catch`).
  - Early `return`, outer-loop `continue`, and outer-loop/switch `break` inside `try`/`catch` propagate cleanly out of the IIFE via `didRet` / `retVal` / `tryFlow` state variables.
- **Switches & Pattern Matching**:
  - Classic `:` switches group empty case headers, hoist cross-case local declarations, and emit explicit Go `fallthrough` where Java falls through.
  - Java 21 switch expressions (`yield`), arrow `->` rules, `case null`, `case Type v`, and `when` guards lower cleanly to value switches or tagless `switch { case j2gort.InstanceOfBind[T](sw, &v) && (guard): ... }` blocks.
- **Lambdas & Method References**:
  - Every functional interface `FI` gets a lightweight adapter `NewLambda_<FI>(fn)` implementing `FI`'s SAM method plus all inherited default interface methods.
  - All 5 method reference forms (`T[]::new`, `Class::new`, `Type::staticMethod`, `expr::instanceMethod` eagerly evaluating & null-checking `expr`, and unbound `Type::instanceMethod`) lower to typed `NewLambda_<FI>` closures.

---

## 5. `rhino.Node` Zero-GC Arena Lowering (`transpiled/shims/nodestore`)

`com.google.javascript.rhino.Node` and its 4 private inner subclasses (`NumberNode`, `BigIntNode`, `StringNode`, `TemplateLiteralSubstringNode`) are **never** heap-allocated as Go structs. Instead, `ClassEmitter` and `MethodEmitter` lower `Node` directly to a 32-bit value-type handle:

```go
type Node nodestore.Node // uint32 index into thread/compilation-scoped 36-byte NodeData slabs
```

1. **Zero-Pointer Arena Layout**:
   - Tree links (`parent`, `first`, `last`, `next`), `type` (`Token`), `LinenoCharno`, `Length16`, `SourceFile`, `PropFlags`, and `Payload` read/write directly on `n.Rec()` (`*nodestore.NodeData`, 36 bytes, zero GC pointers).
   - `NumberNode.number` (`float64`) packs all 64 IEEE-754 bits zero-alloc across `rec.Payload` (low 32 bits) and `rec.PropFlags` (high 32 bits).
   - `StringNode.str` and `Node.originalName` intern into `nodestore.Intern(s) -> StringID` (`0` = `NullString`, `1` = `""`, `>= 2` = pooled string), keeping `NodeStore.origNames` sparse until non-null.
   - `propListHead` and `jsdocInfo` store side-table indices or compact references via `NodeStore` side tables.
2. **Polymorphic Inner-Class Dispatch (`KindFlags`)**:
   - Allocation of `Node` sets `KindFlags = 0`; `NumberNode` = `1`, `BigIntNode` = `2`, `StringNode` = `3`, `TemplateLiteralSubstringNode` = `4`.
   - Any `Node` method overridden by any of the 4 inner classes (`getString`, `setString`, `getDouble`, `setDouble`, `getBigInt`, `isEquivalentTo`, etc.) automatically emits a `switch this_.Rec().KindFlags` dispatcher routing to `nodeSub1_*` .. `nodeSub4_*` or `nodeBase_*`.
   - `j2gort.GetObjectClass(v)` inspects `KindFlags` on `nodestore.Node` handles so `this.getClass() != node.getClass()` inside `Node.isEquivalentTo` preserves exact Java subclass differentiation!
   - Null checks on `Node` compare `n == 0` (throwing `NullPointerException` via `j2gort.ThrowNPE()` on `0.Rec()`), and generic boxing/unboxing converts `Node(0) <-> any(nil)` via `j2gort.BoxNode` / `j2gort.CastNode`.

---

## 6. Shim Binding Architecture & M2/M3 Integration Contract

`RuntimeBindings.java` bridges generated Go packages to the 12 clean-room shim packages in `transpiled/shims/` (documented in `transpiled/shims/API_SURFACE.md`):

- **Already Directly Wired for Compile & Core Execution**:
  - `transpiled/shims/nodestore`: Arena `Node`, `NodeData`, `EnsureCompilationStore`, `Intern`/`String`, `NullString` / `StringID`.
  - `transpiled/shims/javalang`: Wrapping arithmetic/shifts, IEEE-754 `D2I`/`D2L`/`DoubleToString`, `String` UTF-16 indexing & hashing, `Class`, `Throwable`/`Exception` hierarchy & `RecoverException`, `StringBuilder`.
  - `transpiled/shims/javaregex`: `String.matches`, `String.split`, `String.replaceAll`, `String.replaceFirst`.
  - `transpiled/shims/javaio`: `System.out`/`err`, `InputStream`/`Reader`/`PrintStream` core bridges.
  - `transpiled/shims/resources`: Embedded Closure externs, polyfill tables, and `.properties` lookup via `Class.getResourceAsStream`, `Class.getResource`, and `ResourceLoader`.
  - `j2gort/runtime.go`: Self-contained Guava `Preconditions` + `Strings.lenientFormat`, boxed primitives (`BoxInt`, `BoxDouble` with `IsNaN`/`IsInfinite`, etc.), `EnumBase`, `NullString`, reentrant `LockMonitor`, `Arraycopy`, and checked cast/array helpers.
- **Pruned External Bindings (`j2gort/ext_bindings.go` & `j2gort/ext_collections_rt.go`)**:
  - Every external JDK / Guava (`ImmutableList`, `ImmutableMap`, `ImmutableSet`, `Multimap`, `Sets`, `Iterables`, etc.) / Gson / Protobuf / Re2j type and method referenced by Closure's 815 files has a strongly-typed Go struct/interface + constructor/method signature emitted into `j2gort/ext_bindings.go`, pruned strictly to `TypeResolver.isNeededExternalSig` so `go build ./...` succeeds across all 18 packages with zero unused bloat.
  - Generic, zero-boxing native Go collections (`OrderedMap[K, V]`, `CustomOrderedMap[K, V]`, `IdentityOrderedMap[K, V]`, `PlainMap[K, V]`, `OrderedSet[T]`, `PlainSet[T]`, `Deque[T]`, `LocalStringBuilder`, `Slab[T]`, `PooledSlab[T]`, `AcquirePooledSlice[T]`, `CheckNotNullGeneric[T]`, `IsNullGeneric[T]`) are emitted into `j2gort/ext_collections_rt.go` with slice/map pooling and 0 `"unsafe"` imports.

## 7. Idempotent Go-to-Go AST Post-Passes (`transpiled/translator/passes/`)

After `j2go` emits raw, buildable Go (`jsgo_raw`), `cmd/gopasses` runs 11 modular Go-to-Go AST + `go/types` optimization passes (`AllPasses()` in `registry.go`) to fixed point (verified for **100% second-run idempotency**, `Run 2 == 0 sites`, with post-pass `go/types` verification after every pass):

1. **`native_slice_lowering`**: Lowers non-escaping `ArrayList` locals, unexported struct fields, and slice-accumulator parameters to native `[]T` slices.
2. **`native_set_lowering`**: Lowers non-escaping `HashSet`/`LinkedHashSet` locals, fields (including nullable lazy sets), and `Set`-returning helpers to `j2gort.OrderedSet[T]` (and order-insensitive `HashSet`s to Go `map[T]struct{}` `PlainSet[T]`).
3. **`native_map_lowering`**: Lowers non-escaping `HashMap`/`LinkedHashMap`/`IdentityHashMap` locals and fields to `j2gort.OrderedMap[K, V]`, `CustomOrderedMap[K, V]`, `IdentityOrderedMap[K, V]`, or non-iterated `PlainMap[K, V]`.
4. **`immutable_list_lowering`**: Lowers `ImmutableList` fields, locals, getters, and local `ImmutableList.Builder` pipelines directly to native `[]T` slices.
5. **`native_deque_lowering`**: Lowers `ArrayDeque`/`Deque`/`Queue` locals and fields to ring-buffer `j2gort.Deque[T]`.
6. **`string_builder_lowering`**: Lowers non-escaping `StringBuilder` locals to stack-buffered `j2gort.LocalStringBuilder`.
7. **`slab_alloc_lowering`**: Lowers high-frequency AST, `JSType`, `Scope`, `TypedScope`, `Var`, `TypedVar`, `Reference`, `GlobalNamespace_Name`, CFG/DFA graph node/edge, and pass-internal struct allocations onto contiguous `j2gort.Slab[T]` / `PooledSlab[T]` arenas (with bounded pool caps: `768 MiB` global slab cap, `64 MiB` pooled slice cap) with deterministic `ReleaseSlabs_()` cleanup and Cache-Line-0 hot-field packing.
8. **`node_traversal_lowering`**: Pools `jscomp.NodeTraversal` instances, embeds `[12]`-deep scope stacks, reuses `SyntacticScopeCreator`, and emits block-cached `traverseBranchBlk_` carrying `(n, parent, nFirst, nTok, curBlkIdx, curBlk)` in registers to eliminate redundant `nodestore` block loads.
9. **`dead_cast_cleanup`**: Strips redundant self-conversions (`T(x)`, `AsIface[T](x)`).
10. **`abstract_base_devirt`**: Devirtualizes non-overridden abstract-class interface calls (`JSType`, `ObjectType`, `AbstractScope`, `AbstractVar`, `CompilerInput`, etc.) through `*<I>Base` receivers and fields.
11. **`j2gort_peephole`**: Eliminates redundant `Cast*`/`AsIface`/`Box*`/`StrVal`/`CheckNotNull*`/`StringEquals`/`IfaceEquals` wrappers, unwraps zero-arg statement and call-argument IIFEs, optimizes `PredTestNode` / `UnwrapNodePred` / `IsAlwaysTruePred`, and lowers `rhino.Node` accessors to inlineable `nodestore` intrinsics.
