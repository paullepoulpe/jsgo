#!/usr/bin/env bash
# verify_prep.sh: Compiles all 813+ non-Ant patched Java files under OpenJDK 27
# and verifies byte-for-byte CLI parity against upstream closure-compiler-v20261104-maven.jar.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
CLOSURE_SRC="$REPO_ROOT/.cache/closure-src"
PATCHED_CLASSES="$REPO_ROOT/.cache/closure-patched-classes"
CLOSURE_JAR="$REPO_ROOT/.cache/closure-compiler-v20261104-maven.jar"
STAMP_FILE="$CLOSURE_SRC/.prep_applied"

if [[ -d "/opt/homebrew/opt/openjdk/bin" ]]; then
  export PATH="/opt/homebrew/opt/openjdk/bin:$PATH"
elif [[ -d "/opt/homebrew/Cellar/openjdk/27.0.1.0.3/bin" ]]; then
  export PATH="/opt/homebrew/Cellar/openjdk/27.0.1.0.3/bin:$PATH"
fi

echo "==> Checking OpenJDK toolchain..."
javac -version
java -version

if [[ ! -f "$STAMP_FILE" ]]; then
  echo "ERROR: $STAMP_FILE not found. Run ./transpiled/prep/apply_prep.sh first." >&2
  exit 1
fi

if [[ ! -f "$CLOSURE_JAR" ]]; then
  echo "ERROR: Upstream Closure jar not found at $CLOSURE_JAR" >&2
  exit 1
fi

for helper in SimpleMultimap.java SimpleTable.java; do
  if [[ ! -f "$CLOSURE_SRC/com/google/javascript/jscomp/$helper" ]]; then
    echo "ERROR: Missing patched helper $helper in $CLOSURE_SRC" >&2
    exit 1
  fi
done

WORK_TMP="$(mktemp -d "${TMPDIR:-/tmp}/verify-prep.XXXXXX")"
trap 'rm -rf "$WORK_TMP"' EXIT

SHADED_SRC="$WORK_TMP/shaded-src"
mkdir -p "$SHADED_SRC"

python3 - "$CLOSURE_SRC" "$SHADED_SRC" << 'PYEOF'
import re
import sys
from pathlib import Path

src_root = Path(sys.argv[1])
dst_root = Path(sys.argv[2])

# Third-party packages relocated under com.google.javascript.jscomp.jarjar in the Maven uber-jar.
SHADED_PKGS = (
    "com.google.auto.",
    "com.google.common.",
    "com.google.errorprone.",
    "com.google.gson.",
    "com.google.j2objc.",
    "com.google.protobuf.",
    "com.google.re2j.",
    "com.google.thirdparty.",
    "org.kohsuke.args4j.",
)
pattern = re.compile(
    r"(?<![\w$])(" + "|".join(re.escape(p) for p in SHADED_PKGS) + r")"
)

count = 0
for p in sorted(src_root.rglob("*.java")):
    rel = p.relative_to(src_root)
    if rel.parts[:5] == ("com", "google", "javascript", "jscomp", "ant"):
        continue
    text = p.read_text(encoding="utf-8")
    shaded = pattern.sub(r"com.google.javascript.jscomp.jarjar.\1", text)
    out = dst_root / rel
    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text(shaded, encoding="utf-8")
    count += 1

print(f"    Prepared {count} non-Ant Java files for shaded-jar compilation.")
PYEOF

JAVA_LIST="$WORK_TMP/sources.txt"
find "$SHADED_SRC" -name "*.java" | sort > "$JAVA_LIST"

SRC_COUNT="$(wc -l < "$JAVA_LIST" | tr -d ' ')"
echo "==> Compiling all $SRC_COUNT non-Ant Java source files into .cache/closure-patched-classes/..."
if [[ "$SRC_COUNT" -lt 813 ]]; then
  echo "ERROR: Expected at least 813 non-Ant Java files, found $SRC_COUNT" >&2
  exit 1
fi

rm -rf "$PATCHED_CLASSES"
mkdir -p "$PATCHED_CLASSES"

javac -encoding UTF-8 -nowarn \
  -cp "$CLOSURE_JAR" \
  -d "$PATCHED_CLASSES" \
  "@$JAVA_LIST"

CLASS_COUNT="$(find "$PATCHED_CLASSES" -name "*.class" | wc -l | tr -d ' ')"
echo "==> Successfully compiled $SRC_COUNT .java files -> $CLASS_COUNT .class files."

for required_cls in \
  "com/google/javascript/jscomp/CommandLineRunner.class" \
  "com/google/javascript/jscomp/AbstractCommandLineRunner.class" \
  "com/google/javascript/jscomp/SimpleMultimap.class" \
  "com/google/javascript/jscomp/SimpleTable.class" \
  "com/google/javascript/jscomp/instrumentation/reporter/ProductionInstrumentationReporter.class"; do
  if [[ ! -f "$PATCHED_CLASSES/$required_cls" ]]; then
    echo "ERROR: Expected compiled class missing: $required_cls" >&2
    exit 1
  fi
done

echo "==> Verifying JVM classloading resolves CommandLineRunner from .cache/closure-patched-classes/..."
CLASS_ORIGIN="$(java -verbose:class -cp "$PATCHED_CLASSES:$CLOSURE_JAR" \
  com.google.javascript.jscomp.CommandLineRunner --version 2>&1 \
  | grep "com.google.javascript.jscomp.CommandLineRunner " | head -n 1 || true)"
echo "    $CLASS_ORIGIN"
if [[ "$CLASS_ORIGIN" != *"closure-patched-classes"* ]]; then
  echo "ERROR: JVM did not load CommandLineRunner from closure-patched-classes!" >&2
  exit 1
fi

run_parity_case() {
  local name="$1"
  shift
  local stdin_file="$1"
  shift
  # Remaining args are passed to CommandLineRunner, with __OUT_DIR__ substituted per run.
  local dir_ref="$WORK_TMP/${name}_ref"
  local dir_pat="$WORK_TMP/${name}_pat"
  mkdir -p "$dir_ref" "$dir_pat"

  local args_ref=()
  local args_pat=()
  for a in "$@"; do
    args_ref+=("${a//__OUT_DIR__/$dir_ref}")
    args_pat+=("${a//__OUT_DIR__/$dir_pat}")
  done

  local rc_ref=0
  local rc_pat=0

  local jvm_flags=(--enable-final-field-mutation=ALL-UNNAMED --sun-misc-unsafe-memory-access=allow)

  if [[ -n "$stdin_file" ]]; then
    java "${jvm_flags[@]}" -cp "$CLOSURE_JAR" com.google.javascript.jscomp.CommandLineRunner "${args_ref[@]}" \
      < "$stdin_file" > "$dir_ref/stdout.txt" 2> "$dir_ref/stderr.txt" || rc_ref=$?
    java "${jvm_flags[@]}" -cp "$PATCHED_CLASSES:$CLOSURE_JAR" com.google.javascript.jscomp.CommandLineRunner "${args_pat[@]}" \
      < "$stdin_file" > "$dir_pat/stdout.txt" 2> "$dir_pat/stderr.txt" || rc_pat=$?
  else
    java "${jvm_flags[@]}" -cp "$CLOSURE_JAR" com.google.javascript.jscomp.CommandLineRunner "${args_ref[@]}" \
      < /dev/null > "$dir_ref/stdout.txt" 2> "$dir_ref/stderr.txt" || rc_ref=$?
    java "${jvm_flags[@]}" -cp "$PATCHED_CLASSES:$CLOSURE_JAR" com.google.javascript.jscomp.CommandLineRunner "${args_pat[@]}" \
      < /dev/null > "$dir_pat/stdout.txt" 2> "$dir_pat/stderr.txt" || rc_pat=$?
  fi

  # Normalize __OUT_DIR__ path strings if any appear in source maps or diagnostics
  for f in "$dir_ref"/*; do
    [[ -f "$f" ]] || continue
    local base
    base="$(basename "$f")"
    if [[ ! -f "$dir_pat/$base" ]]; then
      echo "ERROR [$name]: Missing output file $base in patched run" >&2
      exit 1
    fi
    sed "s|$dir_ref|__OUT_DIR__|g" "$dir_ref/$base" > "$WORK_TMP/${name}_ref_${base}.norm"
    sed "s|$dir_pat|__OUT_DIR__|g" "$dir_pat/$base" > "$WORK_TMP/${name}_pat_${base}.norm"
    if ! diff -u "$WORK_TMP/${name}_ref_${base}.norm" "$WORK_TMP/${name}_pat_${base}.norm"; then
      echo "ERROR [$name]: Byte mismatch in $base between upstream jar and patched classes!" >&2
      exit 1
    fi
  done

  if [[ "$rc_ref" -ne "$rc_pat" ]]; then
    echo "ERROR [$name]: Exit code mismatch: upstream=$rc_ref vs patched=$rc_pat" >&2
    exit 1
  fi

  echo "    [PASS] $name (exit=$rc_pat)"
}

echo "==> Running real Java-vs-Java CLI parity suite (patched classes vs upstream maven jar)..."

# Case 1: SIMPLE optimization + ES2015+ downleveling to ES5
cat > "$WORK_TMP/case1.js" << 'EOF'
class Greeter {
  constructor(prefix = "Hello") {
    this.prefix = prefix;
  }
  greet(...names) {
    return names.map((n) => `${this.prefix}, ${n}!`).join("; ");
  }
}
const g = new Greeter("Hi");
console.log(g.greet("Closure", "Go"));
EOF
run_parity_case "simple_es6_downlevel" "" \
  --compilation_level SIMPLE \
  --language_in ECMASCRIPT_NEXT \
  --language_out ECMASCRIPT5 \
  --js "$WORK_TMP/case1.js" \
  --js_output_file "__OUT_DIR__/out.js"

# Case 2: ADVANCED optimizations + typechecking + source map generation
cat > "$WORK_TMP/case2.js" << 'EOF'
/**
 * @constructor
 * @param {number} start
 */
function StepCounter(start) {
  /** @type {number} */
  this.value = start;
}
/**
 * @param {number} delta
 * @return {number}
 */
StepCounter.prototype.add = function(delta) {
  this.value += delta;
  return this.value;
};
var c = new StepCounter(10);
window['result'] = c.add(5);
EOF
run_parity_case "advanced_typecheck_sourcemap" "" \
  --compilation_level ADVANCED \
  --warning_level VERBOSE \
  --jscomp_error=checkTypes \
  --js "$WORK_TMP/case2.js" \
  --js_output_file "__OUT_DIR__/out.js" \
  --create_source_map "__OUT_DIR__/out.js.map"

# Case 3: Type error diagnostic & non-zero exit code parity
cat > "$WORK_TMP/case3_err.js" << 'EOF'
/** @param {number} x @return {number} */
function square(x) {
  return x * x;
}
var bad = square("not_a_number");
EOF
run_parity_case "typecheck_diagnostic_error" "" \
  --compilation_level SIMPLE \
  --warning_level VERBOSE \
  --jscomp_error=checkTypes \
  --js "$WORK_TMP/case3_err.js"

# Case 4: JSON streams (--json_streams BOTH) exercising de-reflected JsonFileSpec readJson/writeJson
cat > "$WORK_TMP/case4_stream.json" << 'EOF'
[{"src":"const x = 1 + 2; window['x'] = x;\n","path":"input.js"}]
EOF
run_parity_case "json_streams_both" "$WORK_TMP/case4_stream.json" \
  --compilation_level SIMPLE \
  --json_streams BOTH

# Case 5: Closure primitives (goog.provide / goog.require) + lintChecks (CheckProvidesSorted)
cat > "$WORK_TMP/case5_b.js" << 'EOF'
goog.provide('app.zeta');
goog.provide('app.alpha');
app.alpha = 1;
app.zeta = 2;
EOF
cat > "$WORK_TMP/case5_a.js" << 'EOF'
goog.provide('app.main');
goog.require('app.alpha');
goog.require('app.zeta');
app.main = app.alpha + app.zeta;
EOF
run_parity_case "closure_deps_and_lint" "" \
  --compilation_level SIMPLE \
  --dependency_mode SORT_ONLY \
  --jscomp_warning=lintChecks \
  --js "$WORK_TMP/case5_a.js" \
  --js "$WORK_TMP/case5_b.js"

echo "==> All $SRC_COUNT Java sources compiled and passed 100% byte-identical CLI parity!"
