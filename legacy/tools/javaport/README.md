# javaport

Skeleton generator that takes a closure-compiler Java source file and
emits a Go stub ready for manual porting. Meant to accelerate the
mechanical parts of starting a new Java→Go port; it is **not** a full
transpiler.

## What it extracts

- The `// Port of com/google/javascript/…java @ v20260415 (skeleton)`
  traceability header that `scripts/ci-guards.sh` enforces.
- A Go `struct` for the top-level class, with instance fields (name,
  type) translated via a small type-mapping table.
- A method stub for every public / package / protected method,
  preserving the original javadoc as a Go doc comment.
- A constructor stub rendered as a `NewX(...)` factory.
- A commented catalogue of `JSC_*` diagnostic codes discovered in
  `DiagnosticType.{error,warning,disabled}(...)` declarations.

## What it doesn't do

- Translate method bodies — each stub carries a `// TODO: port body`
  marker and returns the zero value for its return type.
- Handle nested / anonymous classes, enums with methods, AutoValue,
  static initialisers.
- Wire the Java import graph to matching Go packages.

Hard cases fall through to a `// TODO: ...` comment so you can see
what needs manual attention.

## Usage

```sh
# Extract the upstream sources once.
curl -sSL -o /tmp/closure-src.jar \
    https://repo1.maven.org/maven2/com/google/javascript/closure-compiler/v20260415/closure-compiler-v20260415-sources.jar
mkdir -p .cache/closure-src
(cd .cache/closure-src && unzip -qo /tmp/closure-src.jar)

# Generate a skeleton.
go run ./tools/javaport \
    -in  .cache/closure-src/com/google/javascript/jscomp/CheckJSDoc.java \
    -out internal/jscomp/checkjsdoc.go \
    -pkg jscomp
```

The output compiles once you (a) hand-edit bodies, (b) add any
unmapped type names to `javaTypeMap` in `tools/javaport/main.go`,
and (c) add the discovered `JSC_*` codes to `internal/diag/diag.go`.
