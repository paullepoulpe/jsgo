#!/usr/bin/env python3
"""4-Layer Fast-Mode Equivalence & Optimization Quality Gate for jsgo.

Validates that `jsgo --fast_mode` (and `JSGO_FAST_MODE=1`) satisfies:
  - Layer 1 (Go-to-Go Multi-Thread Determinism & Race Freedom):
      Exit code == compat exit code, no data race / panic in stderr, and
      SHA256(fast_out[t=1]) == SHA256(fast_out[t=k]) across all thread counts.
  - Layer 2 (Lockstep AST & Lexical-Scope Alpha-Isomorphism):
      Parses compat and fast outputs with jsgo's own ECMAScript parser, builds
      exact lexical scopes via SyntacticScopeCreator, and walks both ASTs in
      lockstep to verify 100% identical AST tokens, child counts, flags,
      literals, property accesses, extern identifiers, and a strict 1-to-1
      lexical-scope variable bijection (varA <-> varB).
  - Layer 3 (Raw & Gzipped Output Size Non-Regression):
      Verifies raw_bytes(fast_out) <= raw_bytes(compat_out) * 1.001 and
      gzip_bytes(fast_out) <= gzip_bytes(compat_out) * 1.001.
  - Layer 4 (AST / Token Structural Census):
      Counts `function`, `return`, `throw`, control-flow keywords (`if`, `for`,
      `while`, `switch`, `case`, `try`, `catch`), string literals, number
      literals, and total tokens, asserting structural non-regression.

Usage:
  python3 tools/bench/fast_mode_gate.py --jsgo /path/to/jsgo
  python3 tools/bench/fast_mode_gate.py --jsgo /path/to/jsgo --tasks closure-math
  python3 tools/bench/fast_mode_gate.py --jsgo /path/to/jsgo --all-tasks --threads 1,2,4
"""
import argparse
from collections import Counter
import gzip
import hashlib
import json
import math
import os
from pathlib import Path
import re
import shutil
import subprocess
import sys
import tempfile
import time

sys.dont_write_bytecode = True
HERE = Path(__file__).resolve().parent
sys.path.insert(0, str(HERE))
import bench  # noqa: E402
import parallel_gate  # noqa: E402

try:
    import tomllib
except ImportError:
    tomllib = None

DEFAULT_FAST_TASKS = (
    "jquery-adv",
    "d3-adv",
    "three-adv",
    "closure-math",
    "svelte-bundle",
)

CLOSURE_MATH_ENTRY = """goog.provide("closure_math_entry");
goog.require("goog.math");
goog.require("goog.math.AffineTransform");
goog.require("goog.math.Bezier");
goog.require("goog.math.Box");
goog.require("goog.math.Coordinate");
goog.require("goog.math.Coordinate3");
goog.require("goog.math.Integer");
goog.require("goog.math.Line");
goog.require("goog.math.Long");
goog.require("goog.math.Path");
goog.require("goog.math.paths");
goog.require("goog.math.Range");
goog.require("goog.math.Rect");
goog.require("goog.math.Size");
goog.require("goog.math.Vec2");
goog.require("goog.math.Vec3");
window["googMath"] = [
  goog.math, goog.math.AffineTransform, goog.math.Bezier, goog.math.Box,
  goog.math.Coordinate, goog.math.Coordinate3, goog.math.Integer, goog.math.Line,
  goog.math.Long, goog.math.Path, goog.math.paths, goog.math.Range,
  goog.math.Rect, goog.math.Size, goog.math.Vec2, goog.math.Vec3
];
"""

# ---------------------------------------------------------------------------------------------
# Pure-Python JS Lexer & Alpha-Normalizer (Layers 2 & 4)

JS_KEYWORDS = frozenset({
    "async", "await", "break", "case", "catch", "class", "const", "continue",
    "debugger", "default", "delete", "do", "else", "enum", "export", "extends",
    "false", "finally", "for", "function", "if", "import", "in", "instanceof",
    "let", "new", "null", "of", "return", "static", "super", "switch", "this",
    "throw", "true", "try", "typeof", "var", "void", "while", "with", "yield",
})

CONTROL_FLOW_KEYWORDS = frozenset({
    "if", "for", "while", "switch", "case", "try", "catch",
})

REGEX_PRECEDING_KEYWORDS = frozenset({
    "return", "throw", "case", "delete", "void", "typeof",
    "instanceof", "in", "of", "new", "yield", "await", "else", "do",
})

WELL_KNOWN_GLOBALS = frozenset({
    "window", "document", "globalThis", "self", "Math", "Object", "Array",
    "String", "Number", "Boolean", "Function", "Symbol", "BigInt", "Error",
    "TypeError", "RangeError", "ReferenceError", "SyntaxError", "URIError",
    "EvalError", "RegExp", "Date", "Promise", "Map", "Set", "WeakMap",
    "WeakSet", "ArrayBuffer", "SharedArrayBuffer", "DataView", "Int8Array",
    "Uint8Array", "Uint8ClampedArray", "Int16Array", "Uint16Array",
    "Int32Array", "Uint32Array", "Float32Array", "Float64Array",
    "BigInt64Array", "BigUint64Array", "JSON", "Reflect", "Proxy", "Intl",
    "console", "undefined", "NaN", "Infinity", "arguments", "isNaN",
    "isFinite", "parseInt", "parseFloat", "encodeURI", "decodeURI",
    "encodeURIComponent", "decodeURIComponent", "setTimeout", "clearTimeout",
    "setInterval", "clearInterval", "requestAnimationFrame",
    "cancelAnimationFrame", "navigator", "location", "history", "performance",
})

WELL_KNOWN_PROPS = frozenset({
    "prototype", "constructor", "call", "apply", "bind", "length", "toString",
    "valueOf", "hasOwnProperty", "isPrototypeOf", "propertyIsEnumerable",
    "toLocaleString", "push", "pop", "shift", "unshift", "slice", "splice",
    "concat", "join", "indexOf", "lastIndexOf", "forEach", "map", "filter",
    "reduce", "reduceRight", "every", "some", "find", "findIndex", "includes",
    "sort", "reverse", "keys", "values", "entries", "from", "isArray",
    "assign", "create", "defineProperty", "defineProperties",
    "getOwnPropertyDescriptor", "getOwnPropertyDescriptors",
    "getOwnPropertyNames", "getOwnPropertySymbols", "getPrototypeOf",
    "setPrototypeOf", "freeze", "seal", "preventExtensions", "is", "charAt",
    "charCodeAt", "codePointAt", "substring", "substr", "toLowerCase",
    "toUpperCase", "trim", "split", "replace", "match", "search", "test",
    "exec", "then", "catch", "finally", "resolve", "reject", "all", "race",
    "get", "set", "has", "add", "delete", "clear", "size", "next", "done",
    "value", "message", "name", "stack", "abs", "acos", "asin", "atan",
    "atan2", "ceil", "cos", "exp", "floor", "imul", "log", "max", "min",
    "pow", "random", "round", "sin", "sqrt", "tan", "clz32", "fround",
    "hypot", "sign", "trunc", "PI", "E", "LN2", "LN10", "LOG2E", "LOG10E",
    "SQRT1_2", "SQRT2", "MAX_SAFE_INTEGER", "MIN_SAFE_INTEGER", "MAX_VALUE",
    "MIN_VALUE", "EPSILON",
})

TOKEN_MASTER_RE = re.compile(
    r"""
    \s+ |
    //[^\r\n]* |
    /\*[\s\S]*?\*/ |
    (?P<STR>'[^'\\]*(?:\\[\s\S][^'\\]*)*'|"[^"\\]*(?:\\[\s\S][^"\\]*)*"|`[^`\\]*(?:\\[\s\S][^`\\]*)*`) |
    (?P<NUM>(?:0[xX][0-9a-fA-F_]+|0[bB][01_]+|0[oO][0-7_]+|\d+(?:_\d+)*(?:\.\d+(?:_\d+)*)?|\.\d+(?:_\d+)*)(?:[eE][+-]?\d+(?:_\d+)*)?n?) |
    (?P<ID>[A-Za-z_$][A-Za-z0-9_$]*) |
    (?P<SLASH>/=?) |
    (?P<PUNCT>>>>=|===|!==|>>>|<<=|>>=|\*\*=|\+\+|--|<<|>>|<=|>=|==|!=|&&|\|\||\?\?|\+=|-=|\*=|%=|&=|\|=|\^=|=>|\.\.\.|[(){}\[\].,;:?~+\-*%&|^!<>=])
    """,
    re.VERBOSE,
)

REGEX_LITERAL_RE = re.compile(
    r"/(?![/*])(?:[^/\\\[\r\n]|\\[\s\S]|\[(?:[^\]\\]|\\[\s\S])*\])+/[a-zA-Z]*"
)


def tokenize_js(code: str):
    """Tokenizes JS source into a list of (kind, value) tuples, stripping whitespace/comments."""
    tokens = []
    pos = 0
    n = len(code)
    match = TOKEN_MASTER_RE.match
    regex_match = REGEX_LITERAL_RE.match

    while pos < n:
        m = match(code, pos)
        if m is None:
            pos += 1
            continue
        kind = m.lastgroup
        if kind is None:
            pos = m.end()
            continue
        val = m.group()
        if kind == "SLASH":
            prev_kind, prev_val = tokens[-1] if tokens else (None, None)
            can_regex = (
                prev_kind is None
                or (prev_kind == "PUNCT" and prev_val not in (")", "]", "}", "++", "--"))
                or (prev_kind == "KW" and prev_val in REGEX_PRECEDING_KEYWORDS)
            )
            if can_regex:
                rm = regex_match(code, pos)
                if rm is not None:
                    tokens.append(("REGEX", rm.group()))
                    pos = rm.end()
                    continue
            tokens.append(("PUNCT", val))
        elif kind == "ID":
            if val in JS_KEYWORDS:
                tokens.append(("KW", val))
            else:
                tokens.append(("ID", val))
        else:
            tokens.append((kind, val))
        pos = m.end()
    return tokens


def alpha_normalize_and_census(tokens):
    """Computes the canonical alpha-normalized token stream and Layer 4 structural census.

    Each function scope gets its own first-appearance alpha-renaming map so that
    local/global variable renaming permutations, property renaming permutations,
    and `$jscomp$*` / `JSCompiler_*` unique-ID numbering differences normalize
    to identical canonical tokens while preserving intra-function dataflow,
    keywords, literals, operators, and well-known globals/properties.
    """
    census = {
        "function": 0,
        "return": 0,
        "throw": 0,
        "control_flow": 0,
        "string_literals": 0,
        "number_literals": 0,
        "total_tokens": len(tokens),
    }

    canonical = []
    # Stack of (brace_depth_at_fn_entry, id_map_dict)
    scope_stack = [(0, {})]
    brace_depth = 0
    pending_fn_scope = False

    for i, (kind, val) in enumerate(tokens):
        if kind == "KW":
            if val == "function":
                census["function"] += 1
                pending_fn_scope = True
            elif val == "return":
                census["return"] += 1
            elif val == "throw":
                census["throw"] += 1
            elif val in CONTROL_FLOW_KEYWORDS:
                census["control_flow"] += 1
            canonical.append(val)
        elif kind == "PUNCT":
            if val == "(":
                if pending_fn_scope:
                    # Start function scope at parameter list so parameters are alpha-normalized
                    # with the function body.
                    scope_stack.append((brace_depth, {}))
                    pending_fn_scope = False
            elif val == "{":
                brace_depth += 1
            elif val == "}":
                if len(scope_stack) > 1 and scope_stack[-1][0] == brace_depth - 1:
                    scope_stack.pop()
                brace_depth = max(0, brace_depth - 1)
            canonical.append(val)
        elif kind == "STR":
            census["string_literals"] += 1
            canonical.append(val)
        elif kind == "NUM":
            census["number_literals"] += 1
            canonical.append(val)
        elif kind == "REGEX":
            canonical.append(val)
        elif kind == "ID":
            prev_kind, prev_val = tokens[i - 1] if i > 0 else (None, None)
            is_prop = prev_kind == "PUNCT" and prev_val == "."
            if is_prop:
                if val in WELL_KNOWN_PROPS:
                    canonical.append(val)
                else:
                    cur_map = scope_stack[-1][1]
                    key = "." + val
                    cid = cur_map.get(key)
                    if cid is None:
                        cid = f"@P{len(cur_map)}"
                        cur_map[key] = cid
                    canonical.append(cid)
            else:
                if val in WELL_KNOWN_GLOBALS:
                    canonical.append(val)
                else:
                    cur_map = scope_stack[-1][1]
                    cid = cur_map.get(val)
                    if cid is None:
                        cid = f"@I{len(cur_map)}"
                        cur_map[val] = cid
                    canonical.append(cid)

    return canonical, census


def token_stream_delta(canon_a, canon_b):
    """Returns (exact_match: bool, delta_ratio: float) between two canonical token streams."""
    if canon_a == canon_b:
        return True, 0.0
    max_len = max(len(canon_a), len(canon_b), 1)
    ca = Counter(canon_a)
    cb = Counter(canon_b)
    diff_tokens = sum((ca - cb).values()) + sum((cb - ca).values())
    # L1 multiset distance normalized by 2 * max_len gives the fraction of changed tokens.
    delta_ratio = diff_tokens / (2.0 * max_len)
    return False, delta_ratio


# ---------------------------------------------------------------------------------------------
# Task & Fixture Discovery

def builtin_fast_tasks(ws: Path):
    """Returns built-in fast-mode tasks defined over the benchmark workspace."""
    entry_file = ws / "closure_math_entry.js"
    if not entry_file.exists():
        entry_file.write_text(CLOSURE_MATH_ENTRY, encoding="utf-8")

    closure_math_files = [
        "inputs/closure-library@v20230802/closure/goog/base.js",
        "inputs/closure-library@v20230802/closure/goog/debug/error.js",
        "inputs/closure-library@v20230802/closure/goog/dom/nodetype.js",
        "inputs/closure-library@v20230802/closure/goog/asserts/asserts.js",
        "inputs/closure-library@v20230802/closure/goog/reflect/reflect.js",
        "inputs/closure-library@v20230802/closure/goog/array/array.js",
        "inputs/closure-library@v20230802/closure/goog/object/object.js",
        "inputs/closure-library@v20230802/closure/goog/functions/functions.js",
        "inputs/closure-library@v20230802/closure/goog/math/math.js",
        "inputs/closure-library@v20230802/closure/goog/math/coordinate.js",
        "inputs/closure-library@v20230802/closure/goog/math/coordinate3.js",
        "inputs/closure-library@v20230802/closure/goog/math/size.js",
        "inputs/closure-library@v20230802/closure/goog/math/box.js",
        "inputs/closure-library@v20230802/closure/goog/math/irect.js",
        "inputs/closure-library@v20230802/closure/goog/math/rect.js",
        "inputs/closure-library@v20230802/closure/goog/math/affinetransform.js",
        "inputs/closure-library@v20230802/closure/goog/math/bezier.js",
        "inputs/closure-library@v20230802/closure/goog/math/integer.js",
        "inputs/closure-library@v20230802/closure/goog/math/line.js",
        "inputs/closure-library@v20230802/closure/goog/math/long.js",
        "inputs/closure-library@v20230802/closure/goog/math/path.js",
        "inputs/closure-library@v20230802/closure/goog/math/paths.js",
        "inputs/closure-library@v20230802/closure/goog/math/range.js",
        "inputs/closure-library@v20230802/closure/goog/math/tdma.js",
        "inputs/closure-library@v20230802/closure/goog/math/vec2.js",
        "inputs/closure-library@v20230802/closure/goog/math/vec3.js",
        "closure_math_entry.js",
    ]
    closure_math_args = [
        "--compilation_level=ADVANCED",
        "--dependency_mode=PRUNE",
        "--entry_point=goog:closure_math_entry",
        "--jscomp_warning=checkTypes",
    ]
    for f in closure_math_files:
        closure_math_args += ["--js", f]
    closure_math_args += ["--js_output_file=out.js"]

    return {
        "jquery-adv": (
            ws,
            [
                "--compilation_level=ADVANCED",
                "--jscomp_warning=checkVars",
                "--js",
                "corpus/jquery@3.7.1/package/dist/jquery.js",
                "--js_output_file=out.js",
            ],
        ),
        "d3-adv": (
            ws,
            [
                "--compilation_level=ADVANCED",
                "--jscomp_warning=checkVars",
                "--js",
                "corpus/d3@7.9.0/package/dist/d3.js",
                "--js_output_file=out.js",
            ],
        ),
        "three-adv": (
            ws,
            [
                "--compilation_level=ADVANCED",
                "--jscomp_warning=checkVars",
                "--js",
                "inputs/three@0.159.0/build/three.js",
                "--js_output_file=out.js",
            ],
        ),
        "closure-math": (ws, closure_math_args),
        "svelte-bundle": (
            ws,
            [
                "--compilation_level=ADVANCED",
                "--jscomp_warning=checkVars",
                "--js",
                "corpus/preact@10.25.4/package/dist/preact.umd.js",
                "--js",
                "corpus/vue@2.7.16/package/dist/vue.js",
                "--js_output_file=out.js",
            ],
        ),
    }


def load_toml_tasks(ws: Path):
    """Loads any tasks from tools/bench/tasks.toml if present."""
    toml_path = HERE / "tasks.toml"
    if not toml_path.exists() or tomllib is None:
        return {}
    data = tomllib.loads(toml_path.read_text(encoding="utf-8"))
    out = {}
    for name, cfg in data.get("tasks", data).items():
        if isinstance(cfg, dict) and "args" in cfg:
            args = list(cfg["args"])
            if not any(a.startswith("--js_output_file") for a in args) and "--checks_only" not in args:
                args.append("--js_output_file=out.js")
            out[name] = (ws, args)
    return out


def discover_all_cases(bench_ws: Path, fixture_ws: Path):
    """Builds the ordered registry of all available (workspace, args) cases."""
    registry = {}

    # 1. Synthetic fixtures from parallel_gate.py
    fixture_cases = parallel_gate.fixtures(fixture_ws)
    for name, flags in fixture_cases.items():
        args = ["--js_output_file=out.js", "--create_source_map=out.js.map"] + list(flags)
        registry[name] = (fixture_ws, args)

    # 2. Built-in fast tasks
    registry.update(builtin_fast_tasks(bench_ws))

    # 3. Optional tasks.toml tasks
    registry.update(load_toml_tasks(bench_ws))

    # 4. Whole-program tasks from tools/bench/tasks/*.json
    for task in bench.load_tasks():
        if "per_file" in task:
            continue
        tname = task["name"]
        if tname not in registry:
            try:
                jobs = bench.task_jobs(task, str(bench_ws))
                if jobs:
                    registry[tname] = (bench_ws, jobs[0])
            except SystemExit:
                pass

    return registry, list(fixture_cases.keys())


# ---------------------------------------------------------------------------------------------
# Execution & 4-Layer Gate Evaluation

def binary_supports_fast_mode_flag(jsgo_bin: Path, env: dict) -> bool:
    """Checks whether the jsgo binary recognizes the --fast_mode CLI flag."""
    res = subprocess.run(
        [str(jsgo_bin), "--fast_mode", "--version"],
        capture_output=True,
        text=True,
        env=env,
    )
    return '"--fast_mode" is not a valid option' not in res.stderr


def run_jsgo_once(jsgo_bin: Path, base_args: list, cwd: Path, out_dir: Path, env: dict,
                  threads: int, fast_mode: bool, has_fast_flag: bool, timeout: int):
    """Runs jsgo once using a consistent relative .gate-output path inside cwd."""
    if out_dir.exists():
        shutil.rmtree(out_dir)
    out_dir.mkdir(parents=True, exist_ok=True)

    gate_out = cwd / ".gate-output"
    if gate_out.exists():
        shutil.rmtree(gate_out)
    gate_out.mkdir(parents=True, exist_ok=True)

    rel_js = ".gate-output/out.js"
    rel_map = ".gate-output/out.js.map"
    gate_js_path = gate_out / "out.js"
    gate_map_path = gate_out / "out.js.map"

    args = []
    for a in base_args:
        if a == "--js_output_file=out.js" or a.startswith("--js_output_file="):
            args.append(f"--js_output_file={rel_js}")
        elif a == "--create_source_map=out.js.map" or a.startswith("--create_source_map="):
            args.append(f"--create_source_map={rel_map}")
        else:
            args.append(a)

    args.append(f"--num_parallel_threads={threads}")
    run_env = dict(env)
    if fast_mode:
        run_env["JSGO_FAST_MODE"] = "1"
        if has_fast_flag:
            args.append("--fast_mode")
    else:
        run_env.pop("JSGO_FAST_MODE", None)

    t0 = time.perf_counter()
    proc = subprocess.run(
        [str(jsgo_bin)] + args,
        cwd=str(cwd),
        env=run_env,
        stdin=subprocess.DEVNULL,
        capture_output=True,
        timeout=timeout,
    )
    wall_s = time.perf_counter() - t0

    js_bytes = gate_js_path.read_bytes() if gate_js_path.exists() else proc.stdout
    map_bytes = gate_map_path.read_bytes() if gate_map_path.exists() else b""
    if gate_js_path.exists():
        shutil.copy2(gate_js_path, out_dir / "out.js")
    if gate_map_path.exists():
        shutil.copy2(gate_map_path, out_dir / "out.js.map")
    (out_dir / "stdout").write_bytes(proc.stdout)
    (out_dir / "stderr").write_bytes(proc.stderr)
    (out_dir / "exit").write_text(f"{proc.returncode}\n", encoding="utf-8")

    stderr_text = proc.stderr.decode("utf-8", "replace")
    race_or_crash = (
        b"WARNING: DATA RACE" in proc.stderr
        or b"ThreadSanitizer" in proc.stderr
        or bool(bench.CRASH_RE.search(stderr_text))
    )

    return {
        "exit": proc.returncode,
        "wall_s": round(wall_s, 3),
        "js_bytes": js_bytes,
        "js_sha256": hashlib.sha256(js_bytes).hexdigest(),
        "map_sha256": hashlib.sha256(map_bytes).hexdigest() if map_bytes else None,
        "stderr_sha256": hashlib.sha256(proc.stderr).hexdigest(),
        "stderr": stderr_text,
        "race_or_crash": race_or_crash,
    }


def build_ast_isomorphism_bin(out_bin: Path, env: dict) -> Path:
    """Builds the Go lockstep AST alpha-isomorphism verifier binary."""
    subprocess.run(
        ["go", "build", "-o", str(out_bin), "./tools/bench/ast_isomorphism"],
        cwd=str(HERE.parent.parent),
        env=env,
        check=True,
    )
    return out_bin


def verify_ast_isomorphism(ast_iso_bin: Path, compat_js: Path, fast_js: Path, env: dict):
    """Runs the Go lockstep AST + SyntacticScopeCreator alpha-isomorphism verifier."""
    if not compat_js.exists() and not fast_js.exists():
        return True, None, {"renamed_vars_diff": 0, "total_nodes": 0}
    if not compat_js.exists() or not fast_js.exists():
        return False, "one of compat/fast out.js is missing", {}
    proc = subprocess.run(
        [str(ast_iso_bin), str(compat_js), str(fast_js)],
        capture_output=True,
        text=True,
        env=env,
    )
    try:
        payload = json.loads(proc.stdout)
    except Exception:
        return False, f"ast_isomorphism failed: {proc.stderr or proc.stdout}", {}
    return bool(payload.get("isomorphic")), payload.get("error"), payload.get("stats") or {}


def evaluate_task(name: str, jsgo_bin: Path, ast_iso_bin: Path, base_args: list, cwd: Path, task_out_dir: Path,
                  env: dict, threads: list, has_fast_flag: bool, timeout: int):
    """Runs Step 0 (compat) and Layers 1-4 for a single task or fixture."""
    # Step 0: compat baseline (t=1, without --fast_mode)
    compat_dir = task_out_dir / "compat-t1"
    compat = run_jsgo_once(
        jsgo_bin, base_args, cwd, compat_dir,
        env, threads=1, fast_mode=False, has_fast_flag=has_fast_flag, timeout=timeout,
    )

    # Layer 1: Go-to-Go Multi-Thread Determinism & Race Freedom across threads
    fast_runs = {}
    for t in threads:
        fast_runs[t] = run_jsgo_once(
            jsgo_bin, base_args, cwd, task_out_dir / f"fast-t{t}",
            env, threads=t, fast_mode=True, has_fast_flag=has_fast_flag, timeout=timeout,
        )

    ref_t = threads[0]
    fast_ref = fast_runs[ref_t]
    fast_ref_dir = task_out_dir / f"fast-t{ref_t}"

    exits_ok = all(r["exit"] == compat["exit"] for r in fast_runs.values())
    no_race = not compat["race_or_crash"] and all(not r["race_or_crash"] for r in fast_runs.values())
    js_hashes = {r["js_sha256"] for r in fast_runs.values()}
    map_hashes = {r["map_sha256"] for r in fast_runs.values()}
    stderr_hashes = {r["stderr_sha256"] for r in fast_runs.values()}

    l1_deterministic = len(js_hashes) == 1 and len(map_hashes) == 1 and len(stderr_hashes) == 1
    l1_pass = exits_ok and no_race and l1_deterministic

    # Layer 2: Lockstep AST & Lexical-Scope Alpha-Isomorphism + Token Census
    compat_code = compat["js_bytes"].decode("utf-8", "replace")
    fast_code = fast_ref["js_bytes"].decode("utf-8", "replace")

    compat_tokens = tokenize_js(compat_code)
    fast_tokens = tokenize_js(fast_code)

    compat_canon, compat_census = alpha_normalize_and_census(compat_tokens)
    fast_canon, fast_census = alpha_normalize_and_census(fast_tokens)

    _, l2_delta = token_stream_delta(compat_canon, fast_canon)
    ast_iso_ok, ast_iso_err, ast_iso_stats = verify_ast_isomorphism(
        ast_iso_bin, compat_dir / "out.js", fast_ref_dir / "out.js", env
    )
    l2_exact = ast_iso_ok
    l2_pass = ast_iso_ok

    # Layer 3: Raw & Gzipped Output Size Non-Regression (<= +0.1%)
    compat_raw = len(compat["js_bytes"])
    fast_raw = len(fast_ref["js_bytes"])
    compat_gzip = len(gzip.compress(compat["js_bytes"], compresslevel=9)) if compat_raw > 0 else 0
    fast_gzip = len(gzip.compress(fast_ref["js_bytes"], compresslevel=9)) if fast_raw > 0 else 0

    if compat_raw == 0:
        l3_raw_ok = fast_raw == 0
        l3_gzip_ok = fast_gzip == 0
        raw_pct = 0.0 if fast_raw == 0 else 100.0
        gzip_pct = 0.0 if fast_gzip == 0 else 100.0
    else:
        l3_raw_ok = fast_raw <= math.ceil(compat_raw * 1.001)
        l3_gzip_ok = fast_gzip <= math.ceil(compat_gzip * 1.001)
        raw_pct = ((fast_raw - compat_raw) / compat_raw) * 100.0
        gzip_pct = ((fast_gzip - compat_gzip) / compat_gzip) * 100.0 if compat_gzip else 0.0
    l3_pass = l3_raw_ok and l3_gzip_ok

    # Layer 4: AST / Token Structural Census
    l4_checks = {}
    for k in ("function", "return", "throw", "control_flow"):
        c_val = compat_census[k]
        f_val = fast_census[k]
        # Allow exact match or fast_mode inlining/dead-code reduction (f_val <= c_val)
        l4_checks[k] = f_val <= c_val
    l4_pass = all(l4_checks.values())

    overall_pass = l1_pass and l2_pass and l3_pass and l4_pass

    return {
        "task": name,
        "pass": overall_pass,
        "compat_exit": compat["exit"],
        "compat_wall_s": compat["wall_s"],
        "fast_wall_s": {str(t): fast_runs[t]["wall_s"] for t in threads},
        "layer1": {
            "pass": l1_pass,
            "exits_match_compat": exits_ok,
            "no_race_or_crash": no_race,
            "deterministic_across_threads": l1_deterministic,
            "fast_js_sha256": fast_ref["js_sha256"],
            "compat_js_sha256": compat["js_sha256"],
            "byte_identical_to_compat": fast_ref["js_sha256"] == compat["js_sha256"],
        },
        "layer2": {
            "pass": l2_pass,
            "ast_isomorphic": ast_iso_ok,
            "ast_isomorphism_error": ast_iso_err,
            "ast_isomorphism_stats": ast_iso_stats,
            "exact_alpha_match": l2_exact,
            "token_delta_pct": round(l2_delta * 100.0, 4),
        },
        "layer3": {
            "pass": l3_pass,
            "compat_raw_bytes": compat_raw,
            "fast_raw_bytes": fast_raw,
            "raw_delta_pct": round(raw_pct, 3),
            "compat_gzip_bytes": compat_gzip,
            "fast_gzip_bytes": fast_gzip,
            "gzip_delta_pct": round(gzip_pct, 3),
        },
        "layer4": {
            "pass": l4_pass,
            "compat_census": compat_census,
            "fast_census": fast_census,
        },
    }


def format_row(res: dict) -> str:
    l1 = res["layer1"]
    l2 = res["layer2"]
    l3 = res["layer3"]
    l4 = res["layer4"]

    l1_str = "PASS (det)" if l1["pass"] else "FAIL"
    if l2.get("ast_isomorphic"):
        renamed_diff = (l2.get("ast_isomorphism_stats") or {}).get("renamed_vars_diff", 0)
        l2_str = f"PASS (iso Δv={renamed_diff})"
    else:
        l2_str = f"FAIL ({l2.get('ast_isomorphism_error') or 'mismatch'})"

    l3_str = (
        f"{'PASS' if l3['pass'] else 'FAIL'} "
        f"({l3['raw_delta_pct']:+.2f}%/{l3['gzip_delta_pct']:+.2f}%)"
    )

    cc = l4["compat_census"]
    fc = l4["fast_census"]
    l4_str = (
        f"{'PASS' if l4['pass'] else 'FAIL'} "
        f"(fn={fc['function']}/{cc['function']} ret={fc['return']}/{cc['return']} "
        f"cf={fc['control_flow']}/{cc['control_flow']} tok={fc['total_tokens']})"
    )
    status = "PASS" if res["pass"] else "FAIL"
    return f"{res['task']:<24} | {l1_str:<12} | {l2_str:<15} | {l3_str:<21} | {l4_str:<42} | {status}"


def main():
    parser = argparse.ArgumentParser(
        description=__doc__,
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    parser.add_argument("--jsgo", required=True, type=Path, help="Path to the jsgo binary to test")
    parser.add_argument(
        "--tasks",
        help=(
            "Comma-separated list of task/fixture names to run "
            "(default: jquery-adv,d3-adv,three-adv,closure-math,svelte-bundle + synthetic fixtures)"
        ),
    )
    parser.add_argument(
        "--all-tasks",
        action="store_true",
        help="Run all tasks in tools/bench/tasks/ plus default fast tasks and synthetic fixtures",
    )
    parser.add_argument(
        "--threads",
        default="1,2,4",
        help="Comma-separated thread counts to test in fast mode (default: 1,2,4)",
    )
    parser.add_argument(
        "--cache",
        default=bench.default_cache(),
        help="Path to the .cache directory (default: main checkout .cache)",
    )
    parser.add_argument(
        "--out",
        type=Path,
        help="Optional output directory for artifacts and summary.json",
    )
    parser.add_argument("--timeout", type=int, default=180, help="Per-invocation timeout in seconds")
    a = parser.parse_args()

    jsgo_bin = a.jsgo.resolve()
    if not jsgo_bin.exists():
        parser.error(f"jsgo binary not found: {jsgo_bin}")

    threads = [int(x.strip()) for x in a.threads.split(",") if x.strip()]
    if not threads or any(t < 1 for t in threads):
        parser.error("--threads must be a comma-separated list of positive integers")

    env = {
        k: v
        for k, v in os.environ.items()
        if k not in (*bench.GO_ENV_DROP, *bench.JAVA_ENV_DROP)
    }
    env["LC_ALL"] = "C.UTF-8"

    has_fast_flag = binary_supports_fast_mode_flag(jsgo_bin, env)

    temp_ctx = tempfile.TemporaryDirectory(prefix="jsgo-fast-gate-") if a.out is None else None
    try:
        root_dir = a.out.resolve() if a.out is not None else Path(temp_ctx.name)
        root_dir.mkdir(parents=True, exist_ok=True)
        ast_iso_bin = build_ast_isomorphism_bin(root_dir / "ast_isomorphism", env)
        bench_ws = root_dir / "ws-bench"
        fixture_ws = root_dir / "ws-fixtures"
        bench.make_workspace(str(bench_ws), a.cache)
        fixture_ws.mkdir(parents=True, exist_ok=True)

        registry, fixture_names = discover_all_cases(bench_ws, fixture_ws)

        if a.all_tasks:
            selected = list(registry.keys())
        elif a.tasks:
            selected = [t.strip() for t in a.tasks.split(",") if t.strip()]
            unknown = [t for t in selected if t not in registry]
            if unknown:
                parser.error(
                    f"unknown task(s): {', '.join(unknown)} "
                    f"(available: {', '.join(sorted(registry.keys()))})"
                )
        else:
            selected = list(DEFAULT_FAST_TASKS) + fixture_names

        header = (
            f"{'Task':<24} | {'L1 (t=1==2==4)':<12} | {'L2 (α-norm)':<15} | "
            f"{'L3 (raw/gzip Δ)':<21} | {'L4 (census)':<42} | Status"
        )
        sep = "-" * len(header)
        print(
            f"Fast-Mode 4-Layer Quality Gate | binary={jsgo_bin.name} "
            f"| threads={threads} | --fast_mode_flag={'yes' if has_fast_flag else 'env-only'}",
            flush=True,
        )
        print(sep, flush=True)
        print(header, flush=True)
        print(sep, flush=True)

        results = []
        failures = []
        for name in selected:
            ws, base_args = registry[name]
            res = evaluate_task(
                name=name,
                jsgo_bin=jsgo_bin,
                ast_iso_bin=ast_iso_bin,
                base_args=base_args,
                cwd=ws,
                task_out_dir=root_dir / "runs" / name,
                env=env,
                threads=threads,
                has_fast_flag=has_fast_flag,
                timeout=a.timeout,
            )
            results.append(res)
            if not res["pass"]:
                failures.append(name)
            print(format_row(res), flush=True)

        print(sep, flush=True)
        summary = {
            "jsgo": str(jsgo_bin),
            "jsgo_sha256": bench.sha256_file(str(jsgo_bin)),
            "threads": threads,
            "supports_fast_mode_flag": has_fast_flag,
            "all_passed": len(failures) == 0,
            "passed_count": len(results) - len(failures),
            "failed_count": len(failures),
            "failures": failures,
            "results": results,
        }
        if a.out is not None:
            summary_path = root_dir / "summary.json"
            summary_path.write_text(json.dumps(summary, indent=2) + "\n", encoding="utf-8")
            print(
                f"{summary['passed_count']}/{len(results)} tasks passed "
                f"(summary: {summary_path})",
                flush=True,
            )
        else:
            print(f"{summary['passed_count']}/{len(results)} tasks passed.", flush=True)

        if failures:
            raise SystemExit(f"fast_mode_gate failed on: {', '.join(failures)}")
    finally:
        if temp_ctx is not None:
            temp_ctx.cleanup()


if __name__ == "__main__":
    main()
