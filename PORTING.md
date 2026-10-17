# Porting conventions

This is the contract for translating Closure Compiler (Java, pinned at
**v20261015**, see `tools/oracle/pin.env`) into Go. Every ported file follows it, so that files ported
independently (by different people or agents, in parallel) fit together
without negotiation, and so that a Java reader can find any Go symbol
mechanically.

The one-line summary: **translate whole Java files, top to bottom,
mechanically; adapt only where Go cannot express the Java; never add
behaviour that the Java does not have.**

**No GPL code, anywhere.** What we translate is Closure Compiler itself
(Apache 2.0), plus the parts of Guava (Apache 2.0), Gson (Apache 2.0),
args4j (MIT) and protobuf (BSD) that it uses. The JDK is GPLv2 with the
Classpath Exception, and **no JDK code may enter this repository** in any
form:
- no translation of JDK source files;
- no translation of `javap` disassembly of JDK classes;
- no tables or data extracted from JDK internals (reflection on `sun.*`,
  `jdk.internal.*` and the like).

Where Go must reproduce Java library behaviour (`internal/jutil`), write it
from the public specification:
- the Javadoc;
- the JLS;
- Unicode and IETF standards;
- published algorithms.

Then verify it black-box against outputs that the JDK's public API prints
(vector generators under `tools/oracle/java`). Don't open the JDK's source
while writing such code. If exact behaviour cannot be specified that way,
stop and raise it; never fall back to the JDK's code.

**Licence notices.** Every ported file keeps the licence comment of the
Java file it ports, right after its `// Port of` line; a file that ports
several keeps each one's. Closure ships the Rhino-derived files under
`com/google/javascript/rhino` under a dual licence: MPL 1.1, or
alternatively GPL 2.0. This repository uses them under the MPL 1.1 only.
Their ports carry the MPL 1.1 notice with the GPL alternative removed.
`go run ./tools/licenses` writes those notices, for main and test sources
alike. The few Java files that have no licence comment (the AutoValue and
AutoBuilder output) are listed, with the reason, in `noNotice` in
`tools/licenses/closure.go`.

The same rule applies to ports of Guava, Gson, protobuf-java or args4j
code, at the versions the pinned jar ships (Guava 33.4.6-jre, Gson 2.9.1,
protobuf-java 4.30.2, args4j 2.33):
- A translation starts with `// Port of <path> (<project> <version>)`,
  followed by the upstream file's own copyright and licence header.
- A file that only borrows a function starts instead with a short
  `// Derived from <path> (<project> <version>)` block. The block says
  what is derived, and gives the copyright line and the licence, with a
  pointer to its text under `LICENSES/`.
- Code written from the documentation needs neither.

List each such file in `NOTICE`. `go run ./tools/licenses -check`, run by
`make vet`, fails in four cases:
- a tracked file mentions the GPL;
- a file names a Closure source without a Closure licence comment, unless
  `noNotice` lists the source;
- a file names a third-party source without that project's notice;
- `NOTICE` or a licence text under `LICENSES/` is missing.

## 1. Unit of work

- One Java top-level class (with all of its nested classes) → one Go file.
  `com/google/javascript/rhino/Node.java` → `internal/rhino/node.go`
  (snake_case of the class name). Nested classes stay in the same file.
- A file is ported **completely**: every field, constructor, method and
  nested type, in Java declaration order. No stubs, no "subset" ports, no
  `TODO: port body`. If a member depends on a class that is not ported yet,
  port that class first (we port bottom-up in dependency order) — or, if it
  is far up the stack (e.g. a debugging hook), port the member with a call
  into a clearly named, not-yet-ported symbol and let the build fail until
  that dependency lands. Never paper over it with invented logic.
- **Forward declarations are retired.** While the port was in progress,
  `Compiler` and `DefaultPassConfig` (one SCC that references every pass)
  could build before every pass existed because `tools/pending` generated
  `internal/jscomp/pending_gen.go`: declarations of the not-yet-ported
  classes, whose bodies panicked with `pending port: <Java signature>` and
  whose tests skipped. The port is complete: the file, its generator and
  the `PendingPort` skip plumbing are gone, and nothing may be
  forward-declared.
- Each file starts with:
  ```go
  // Port of com/google/javascript/rhino/Node.java @ closure-compiler v20261015
  ```
- Every Go function carries a `// Java: <Class>.<method>(<param types>)`
  line in its doc comment, which makes overload mapping greppable. Java doc
  comments are kept (reflowed as Go comments).
- Java tests are ported alongside (`NodeTest.java` → `node_test.go`). A test
  helper class that Java tests of two Go packages share, and that imports a
  higher package, lives in an importable `…testing` package (like
  `internal/rhino/rhinotesting`, or `internal/sourcemap/sourcemaptesting`
  for `SourceMapTestCase`), and the tests that use it are external
  (`package foo_test`) test files.

## 2. Packages

Go forbids import cycles; Java allows them. A Go package is therefore one
**strongly connected component** of the Java package import graph
(generated protobuf classes count as a leaf package):

| Go package (`internal/…`) | Java packages |
|---|---|
| `jutil` | — (Java runtime + Guava semantics; see §6) |
| `base` | `jscomp.base` |
| `sourcemap` | `debugging.sourcemap` |
| `parser` | `jscomp.parsing.parser`, `.parser.trees`, `.parser.util` |
| `rhino` | `rhino`, `rhino.jstype` |
| `dtoa` | `rhino.dtoa` |
| `parsing` | `jscomp.parsing` |
| `regex` | `jscomp.regex` |
| `colors`, `diagnostic`, `resources` | `jscomp.colors`, `jscomp.diagnostic`, `jscomp.resources` |
| `jscomp` | `jscomp` and every package in its cycle (`deps`, `graph`, `lint`, `modules`, `transpile`, `type`, `disambiguate`, `ijs`, `instrumentation`, `serialization`, `bundle`, `refactoring`) |

Regenerate the table with `tools/oracle/pkgscc.py` if the pin changes.

## 3. Names

The rules are mechanical so that code can be written against a class that
another agent is porting at the same time.

| Java | Go |
|---|---|
| top-level class/interface/enum `Foo` | type `Foo` |
| nested type `Outer.Inner` | type `Outer_Inner` (e.g. `Node_Prop`, `Config_LanguageMode`) |
| instance method `getFirstChild()` | method `GetFirstChild()` — **all** methods are exported, whatever their Java visibility |
| field `firstChild` | field `firstChild` — fields keep their Java name, unexported |
| field read from another Go package | exported accessor method if Java has one; otherwise field renamed `FirstChild` |
| static method `IR.name(s)` | func `IR_Name(s)` |
| static field / constant `Node.MAX_X` | var/const `Node_MAX_X` |
| enum constant `Token.NAME` | const `Token_NAME` |
| constructor `new Foo(a, b)` | func `NewFoo(a, b) *Foo` |

**Overloads.** Go has no overloading. The overload with the fewest
parameters keeps the plain name; each other overload appends a suffix naming
what distinguishes it: `NewString(s)` / `NewStringWithToken(tok, s)`,
`Block()` / `BlockOfNode(n)` / `BlockOfList(list)`. The `// Java:`
signature line on each function records the mapping. When calling a class
that is not ported yet, guess by this rule; the compiler flags any mismatch
at integration.

Go keywords and predeclared names that collide with Java identifiers get a
trailing underscore (`type_`, `func_`, `range_`, `len_`).

**`Iterable` parameters.** A parameter declared `Iterable<T>` is an
`iter.Seq[T]`. Callers pass `slices.Values(list)` or a collection's `All()`.

## 4. Classes, interfaces and inheritance

- **Concrete class without subclasses** → `struct` with pointer-receiver
  methods. Java `this` → receiver named after the class's initial (`n`
  for `Node`).
- **Class with subclasses (abstract or not)** → an `interface` named after
  the class, declaring every overridable/public method, plus a struct
  `FooBase` holding the fields and the base implementations. Subclasses
  embed `FooBase`. Variables declared with the Java type `Foo` use the
  interface.
  - Virtual calls from base methods go through a `self Foo` field set by the
    constructor: `b.self.Method()`. Without this, Go would call the base
    implementation and silently diverge from Java.
  - `super.method()` → `x.FooBase.Method()`.
  - Field access through the interface: the interface declares
    `GetFooBase() *FooBase`.
  - `x instanceof Bar` → `_, ok := x.(*Bar)`; `(Bar) x` → `x.(*Bar)`.
  - Constructors: `NewBar(...)` allocates, sets `self`, then calls
    `b.InitFooBase(...)` (the Go form of `super(...)`), then runs Bar's
    own constructor body.
- **Private subclasses that only override accessors** (e.g. `Node`'s
  `StringNode`, `NumberNode`) may be folded into the parent struct with a
  payload field, as long as every method keeps the exact Java behaviour,
  including the exceptions the base implementation throws. Mark the fold
  with `// Go-only: language-adapter — <why>`.
- **Java interface** → Go interface of the same name. Default methods →
  package function `Foo_Method(f Foo, ...)` or a `FooDefaults` struct to
  embed.
- **Anonymous classes and lambdas** → Go closures, or a small unexported
  struct when the interface has several methods. A single-method interface
  `Foo` that Java code implements with lambdas gets one adapter next to the
  interface: `type FooFunc func(...)` whose method calls it (e.g.
  `CompilerPassFunc`, `NodeUtil_VisitorFunc`), marked
  `// Go-only: language-adapter`. Whoever ports the interface writes it; until
  then, call it by that name.
- **Enums** → `type Foo int32` + `Foo_CONSTANT` consts in Java ordinal
  order, plus `Foo_values()`, `String()` (Java `name()`), and switch-based
  methods for enum bodies. Enum fields → lookup tables indexed by ordinal.
  `EnumSet<Foo>` → `jutil` bitset or `*jutil.LinkedHashSet[Foo]` (iteration
  in ordinal order must be preserved where observable).
- **AutoValue / records** → plain struct + constructor + `Equals`/`HashCode`
  when the Java relies on value equality.
- **`equals`/`hashCode`/`toString`** → `Equals(o any) bool`,
  `HashCode() int32`, `String() string`.

## 5. Statements and expressions

- `null` → `nil`. Never store a typed nil pointer in an interface variable
  (Go's `iface != nil` would then be true where Java's `x != null` is false).
- Object identity: Java gives every `new` a distinct identity. In Go,
  pointers to *zero-size* structs may compare equal. A class with no fields
  (common for stateless strategy objects, guards and passes) that is ever
  compared by identity, or used as a map/set key, must have a non-zero
  size. Give its base struct a one-byte field, marked
  `// Go-only: language-adapter — distinct identity`. Wave 2 found this in
  `ComposeWarningsGuard`, where two guards collapsed into one.
- Exceptions: `throw new X(..)` → `panic(NewX(..))`; `try/catch` →
  `jutil.Try(func() {...}, func(e jutil.IRuntimeException) {...})`, where
  the marker interface gives Java's catch-by-subclass semantics;
  `finally` → `defer`. A `return` inside a Java `try` body becomes a
  captured result variable plus a flag. Exception classes embed their Java
  superclass struct (ultimately `jutil.ThrowableBase`) and declare their own
  marker interface.
- A null dereference in Java is a `NullPointerException` whose message the
  JVM derives from the bytecode (JEP 358: `Cannot invoke "…Node.getSecondChild()"
  because "clazz" is null`). A Go nil dereference becomes one without a
  message. Where the message reaches Java's output (an INTERNAL COMPILER
  ERROR report prints it), check for nil and throw
  `jutil.NewHelpfulNullPointerException` with the JVM's text, marked as a
  language adapter. Oracle dumps made with
  `-XX:-ShowCodeDetailsInExceptionMessages` compare with
  `jutil.SetShowCodeDetailsInExceptionMessages(false)`.
- Guava `Preconditions.checkState/checkArgument/checkNotNull` →
  `jutil.CheckState`, `jutil.CheckArgument`, `jutil.CheckNotNull`. These are
  behaviour, not decoration: keep every one.
  - With a message template, use the fixed-arity form matching the number
    of template arguments:
    - `jutil.CheckStateMsg1/2/3`;
    - `jutil.CheckArgumentMsg1/2/3`;
    - `jutil.CheckNotNullMsg1/2`.
  - Use the variadic `...Msg` form only for 4+ arguments. The variadic form
    boxes every argument on every call, even when the check passes. It cost
    8% of parse time before it was fixed. These forms mirror Guava's own
    non-varargs overloads.
- Allocation is the main cost of a line-for-line port in Go. The JVM's
  escape analysis removes allocations that Go keeps. So, in code that runs
  per node, per token or per character:
  - do not pass `...any`;
  - do not convert to `any`;
  - do not create closures for `jutil.Try`;
  - do not build strings eagerly for messages that are only needed on
    failure.

  Keep such costs on the failure path, as `jutil` does. Measure with
  `tools/bench` before optimising anything else; behaviour must not change
  (the oracles stay at zero diffs).
- `switch` with fallthrough → Go `switch` with explicit `fallthrough`;
  Java switch expressions → a Go `switch` assigning a variable.
- Labeled `break`/`continue` map directly.
- Streams → explicit loops, preserving encounter order and short-circuiting.
- `assert` → drop (Java runs without `-ea` in production).
- Threads → sequential code. Concurrency is out of scope until parity;
  Closure is deterministic regardless of threads.
- `synchronized` methods and blocks, and `ReentrantLock` → `jutil.Monitor`,
  never `sync.Mutex`. Java's monitors are re-entrant and Go's mutexes are
  not, so a translated call chain that comes back to a monitor it holds
  deadlocks on a mutex. This is easy to miss: the synchronized
  `ThreadSafeDelegatingErrorManager.generateReport` formats a diagnostic,
  which loads an input source map, which reports a warning through the same
  manager.
  - An instance `synchronized` method or `synchronized (this)` locks a
    `monitor jutil.Monitor` field of the object. A subclass uses the field
    of the superclass struct it embeds: it is the same Java object.
  - A `static synchronized` method locks a package-level
    `<class>_classMonitor` (the `Foo.class` monitor).
  - A method locks with `m.Lock()` then `defer m.Unlock()`; a block becomes
    a closure that does the same.
  - The owner is the goroutine (`jutil.GoroutineID`). Taking a monitor costs
    a stack walk, so keep Java's double-checked fast paths outside it.
  - `sync.Mutex` is only for Go-only state whose critical section runs no
    translated code, such as an emulated JDK cache.

## 6. Java semantics that Go does not share

These are the places where a "natural" Go translation silently changes
behaviour. Use the `internal/jutil` helpers.

- **Strings are UTF-16 in Java.** Go strings in this port hold WTF-8 (UTF-8
  that may also encode lone surrogates, which JS string literals allow).
  Any Java `String` API that takes or returns an index, a length or a char
  must use `jutil.StringLength`, `StringCharAt`, `StringSubstring`,
  `StringIndexOf`, ... — not `len(s)` or `s[i]` — unless the string is
  provably ASCII (identifiers generated by the compiler, fixed keywords).
  Ordering uses `jutil.StringCompareTo` (UTF-16 order differs from byte
  order), and hashing uses `jutil.StringHashCode`. The scanner works on
  `[]uint16`, exactly like Java's `char` arrays, so source positions are in
  UTF-16 code units as in Java.
- **`char`** → `uint16`; **`int`** → `int` for indices and counts, `int32`
  wherever 32-bit wraparound or bit operations can matter (hash codes,
  bit flags, `>>>`, JS ToInt32 folding); **`long`** → `int64`;
  **`double`** → `float64`.
- Narrowing casts `(int) d` / `(long) d` saturate in Java →
  `jutil.D2I` / `jutil.D2L`. Shifts mask the distance in Java →
  `jutil.Shl32/Shr32/Ushr32/Shl64/...`.
- `Double.toString` / `"" + d` / `String.valueOf(d)` →
  `jutil.DoubleToString` (verified against JDK 21 on 27k vectors).
  String concatenation of arbitrary values → `jutil.ToJavaString`.
- **Collections and iteration order.** Go map iteration is random; Java's is
  not. Never `range` over a Go map where the order can reach output,
  diagnostics or any decision. Mapping:
  - `ArrayList`/`List` that is mutated after being shared →
    `*jutil.ArrayList[T]`; `ImmutableList`, `List.of`, and lists never
    mutated after construction → `[]T`.
  - `LinkedHashMap`, `ImmutableMap` → `*jutil.LinkedHashMap[K,V]`;
    `LinkedHashSet`, `ImmutableSet` → `*jutil.LinkedHashSet[T]`.
  - `HashMap`, `HashSet`, Guava `HashMultimap`, `MultimapBuilder.hashKeys()
    .hashSetValues()`, and the Guava hash factories (`Maps.newHashMap*`,
    `Sets.newHashSet*`) → `*jutil.JavaHashMap[K,V]`, `*jutil.JavaHashSet[T]`,
    `*jutil.HashMultimap[K,V]` (`Maps_NewHashMap*`, `Sets_NewHashSet*`,
    `HashMultimap_Create*`, `MultimapBuilder_HashKeys_HashSetValues`), when
    both hold:
    - the collection is iterated: directly, through a copy
      (`new ArrayList<>(set)`, `ImmutableList.copyOf`), a stream, `toString`,
      `Sets.intersection`/`difference` over it, or a multimap's `keySet`,
      `values`, `entries`, `asMap` or `get(k)`;
    - its keys (for a multimap: its keys or its values) hash by value:
      `String`, boxed primitives, value classes and AutoValue types, lists of
      those.

    These reproduce the JDK 21 / shaded Guava iteration order bit for bit
    (table growth, bin order, tree bins), verified against vectors from
    `tools/oracle/java/JutilJavaHashVectors.java`. Use the constructor that
    matches the Java one: the initial capacity and the factory change the
    order. Keys hash with `ObjectHashCode` and compare with `ObjectsEquals`,
    so value classes need `Equals`/`HashCode` (and `CompareTo(T) int` if
    the Java class is `Comparable` of itself). Structural changes during
    iteration throw `ConcurrentModificationException` as in Java; remove
    through `Iterator()`.
  - A `HashMap`/`HashSet`/hash multimap that is only looked up (`get`,
    `contains`, `put`, `remove`, never iterated), or whose keys are
    identity-hashed (`Node`, `Var`, enums, plain objects: Java's own order
    is then nondeterministic, so Closure cannot depend on it; `JSType` is
    mixed, with identity hashes for constructors, `AllType`, `NoObjectType`
    and others, so treat it as identity-hashed)
    stays `*jutil.LinkedHashMap`/`LinkedHashSet` or an insertion-ordered
    multimap. A hash multimap with identity-hashed keys but value-hashed
    values whose value order is observable needs `HashMultimap`. When it is
    unclear whether the order is observable, flag the site with
    `// Go-only: language-adapter — Java HashMap order` and raise it rather
    than guessing.
  - Guava `HashBasedTable` is `LinkedHashMap`-backed (rows and columns in
    insertion order) in the pinned Guava → `*jutil.HashBasedTable`.
  - `TreeMap`/`TreeSet` → sorted `jutil` structures with the same
    comparator. `IdentityHashMap` → `jutil.LinkedHashMap` keyed by pointer
    (Java's own iteration order there is not deterministic, so Closure never
    depends on it).
  - Keys with custom `equals`/`hashCode` need the `jutil` hashed variants;
    a Go map would compare pointers instead.
- **Regex.** `java.util.regex` is not RE2, and the differences are not only
  backreferences and lookaround:
  - Java's `\s` includes `\x0B`;
  - Java's `.` also excludes `\r`, `\u0085`, `\u2028` and `\u2029`;
  - without MULTILINE, Java's `$` also matches before a final line
    terminator;
  - `matches()`, `lookingAt()` and `find()` differ, as do `split` and
    `replaceAll` replacement syntax.

  So every `Pattern`/`Matcher`/`String.matches`/`replaceAll`/`split(regex)`
  goes through `jutil.Pattern` (and its Matcher and String helpers), which
  translates Java syntax and semantics and panics on constructs it cannot
  express. Do not import Go's `regexp` in ported code. Where `jutil.Pattern`
  lacks a construct, extend it (jutil-only commit with JDK vectors) or write
  an equivalent hand matcher marked `// Go-only: language-adapter —
  <construct>`.

### Extending `jutil` when several ports run in parallel

`jutil` is shared by every branch, so it is where parallel work collides
(two `TreeSet`s were written independently in wave 2). To avoid that:

- **Before adding a helper,** check the main branch's `internal/jutil` and
  the other active branches (`git log <branch> -- internal/jutil`).
- **One Java/Guava class per file, named after it** (`tree_set.go`,
  `bit_set.go`, `message_format.go`), tested against vectors generated by
  the JDK (`tools/oracle/java/*Vectors.java`).
- **Commit `jutil` changes in commits that touch only `internal/jutil`**
  (plus their vector generator), so other branches can cherry-pick them.
- **Never change the signature of an existing `jutil` function.** Add a new
  one instead.
- **Duplicates are resolved on the main branch** by keeping one
  implementation, the one tested against Java. Branches merge main to pick
  that up.

## 7. What is allowed to differ from Java

Only what Go forces. Mark each such spot with
`// Go-only: language-adapter — <reason>`. Code that exists to make a test
pass, a heuristic, a special case, a "pre-pass" emulating Java, or a
shortcut is never acceptable, however good the motivation. If the port does
not match Java's behaviour, the port has a bug: find the Java logic that
produces the behaviour and translate it.

One deliberate exception, decided by the repository owner: jutil's
`java.util.zip` compression (Deflater, GZIPOutputStream) uses Go's
`compress/flate`, not a translation of zlib, so compressed bytes may differ
from Java's while the decompressed content is identical; oracles and
tools/clie2e compare gzip output decompressed.

## 8. How a port is verified

Correctness is checked against the pinned Java compiler, layer by layer
(see `PLAN.md` §4 and `tools/oracle/`):

1. Ported unit tests from the Java test suite.
2. Differential oracles on large corpora. The parser must produce a
   `rhino.Node` dump byte-identical to Java's (`tools/oracle/java/AstDump`);
   the code printer must produce byte-identical output; each pass must
   produce the same AST as Java's pass on the same input; the CLI must
   produce the same output and diagnostics.

A layer is done when its oracle diff is empty on the whole corpus, not when
a number crosses a threshold.
