# Porting conventions

This is the contract for translating Closure Compiler (Java, pinned at
**v20260415**) into Go. Every ported file follows it, so that files ported
independently (by different people or agents, in parallel) fit together
without negotiation, and so that a Java reader can find any Go symbol
mechanically.

The one-line summary: **translate whole Java files, top to bottom,
mechanically; adapt only where Go cannot express the Java; never add
behaviour that the Java does not have.**

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
- Each file starts with:
  ```go
  // Port of com/google/javascript/rhino/Node.java @ closure-compiler v20260415
  ```
- Every Go function carries a `// Java: <Class>.<method>(<param types>)`
  line in its doc comment, which makes overload mapping greppable. Java doc
  comments are kept (reflowed as Go comments).
- Java tests are ported alongside (`NodeTest.java` → `node_test.go`) when the
  test does not depend on not-yet-ported infrastructure.

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
  struct when the interface has several methods.
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
- Exceptions: `throw new X(..)` → `panic(NewX(..))`; `try/catch` →
  `jutil.Try(func() {...}, func(e jutil.IRuntimeException) {...})`, where
  the marker interface gives Java's catch-by-subclass semantics;
  `finally` → `defer`. A `return` inside a Java `try` body becomes a
  captured result variable plus a flag. Exception classes embed their Java
  superclass struct (ultimately `jutil.ThrowableBase`) and declare their own
  marker interface.
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
- `synchronized`, threads → sequential code. Concurrency is out of scope until
  parity; Closure is deterministic regardless of threads.

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
  - `LinkedHashMap`, `ImmutableMap`, `HashMap` → `*jutil.LinkedHashMap[K,V]`;
    `LinkedHashSet`, `ImmutableSet`, `HashSet` → `*jutil.LinkedHashSet[T]`.
    A plain Java `HashMap`/`HashSet` whose iteration order is observable
    needs Java's hash order (only ~10 such sites exist in Closure); flag it
    with `// Go-only: language-adapter — Java HashMap order` and raise it
    rather than guessing.
  - `TreeMap`/`TreeSet` → sorted `jutil` structures with the same
    comparator. `IdentityHashMap` → `jutil.LinkedHashMap` keyed by pointer
    (Java's own iteration order there is not deterministic, so Closure never
    depends on it).
  - Keys with custom `equals`/`hashCode` need the `jutil` hashed variants;
    a Go map would compare pointers instead.
- **Regex.** `java.util.regex` is not RE2. Each of the ~52 uses must be
  checked for backreferences/lookaround and ported with an equivalent
  matcher where needed.

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
