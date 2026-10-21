#!/usr/bin/env python3
"""Set up the production-scale googleapis multi-stage, multi-chunk benchmark.

Architecture:
1. Code generation (~35M AST nodes):
   - Runs `protoc --js_out=import_style=closure,binary:<out_dir>` once per
     `.proto` directory in `googleapis` (plus the 11 well-known types in
     `google/protobuf/*.proto`), each into its own output directory.
   - Adds the Closure protobuf runtime (`protobuf-javascript` v3.21.4) and
     Closure Library (`closure-library@v20230802`).
2. Chunks (~7.5k chunks):
   - Splits Closure Library + protobuf runtime along its `goog.require` /
     `goog.requireType` graph (53 chunks, depth 19, `base.js` at the root).
   - Creates one chunk per `.proto` file (7,407 `.proto` chunks; 7,460 total
     chunks, total DAG depth 30), with its `.proto` imports and protobuf
     runtime chunks as parent chunks.
3. Keep code alive:
   - Generates a small keep-alive module per leaf `.proto` chunk (1,492 leaf
     chunks) that instantiates, `serializeBinary()`s, and `deserializeBinary()`s
     messages from that chunk and exports the result via `goog.exportSymbol`.
4. Outputs:
   - `.cache/bench-inputs/googleapis-bench/flags.txt`: common flags, `--js`
     inputs, and `--chunk` specs in strict topological order for 3-stage
     `CHECKS -> OPTIMIZATIONS -> FINALIZATIONS` compilation.
   - `.cache/bench-inputs/googleapis-bench/manifest.json`: summary statistics.
"""

from __future__ import annotations

import concurrent.futures
import glob
import json
import os
import platform
import re
import shutil
import subprocess
import sys
import time
import urllib.request

ROOT_DIR = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
CACHE_DIR = os.path.join(ROOT_DIR, ".cache", "bench-inputs")
BENCH_DIR = os.path.join(CACHE_DIR, "googleapis-bench")

GOOGLEAPIS_DIR = os.path.join(CACHE_DIR, "googleapis")
CLOSURE_LIB_DIR = os.path.join(CACHE_DIR, "closure-library@v20230802")
PB_BIN_DIR = os.path.join(CACHE_DIR, "protobuf-javascript")
PB_SRC_DIR = os.path.join(CACHE_DIR, "protobuf-javascript-src")

PB_VERSION = "3.21.4"

WELL_KNOWN_PROTOS = [
    "google/protobuf/any.proto",
    "google/protobuf/api.proto",
    "google/protobuf/descriptor.proto",
    "google/protobuf/duration.proto",
    "google/protobuf/empty.proto",
    "google/protobuf/field_mask.proto",
    "google/protobuf/source_context.proto",
    "google/protobuf/struct.proto",
    "google/protobuf/timestamp.proto",
    "google/protobuf/type.proto",
    "google/protobuf/wrappers.proto",
]

PB_RUNTIME_REL_FILES = [
    "asserts.js",
    "debug.js",
    "map.js",
    "message.js",
    "binary/arith.js",
    "binary/constants.js",
    "binary/decoder.js",
    "binary/encoder.js",
    "binary/reader.js",
    "binary/utils.js",
    "binary/writer.js",
]

KEEPALIVE_MESSAGES_PER_LEAF = 10


def ensure_prerequisites() -> tuple[str, str, str]:
    os.makedirs(CACHE_DIR, exist_ok=True)

    protoc = shutil.which("protoc") or "/opt/homebrew/bin/protoc"
    if not os.path.exists(protoc):
        raise RuntimeError("protoc not found; please install protobuf (e.g. brew install protobuf)")

    # Locate protobuf include directory containing google/protobuf/any.proto
    protoc_inc = None
    for cand in [
        "/opt/homebrew/include",
        "/usr/local/include",
        "/usr/include",
        os.path.join(os.path.dirname(os.path.dirname(protoc)), "include"),
    ]:
        if os.path.exists(os.path.join(cand, "google/protobuf/any.proto")):
            protoc_inc = cand
            break
    if not protoc_inc:
        raise RuntimeError("Could not find google/protobuf/*.proto include directory")

    # Ensure protoc-gen-js binary
    protoc_gen_js = os.path.join(PB_BIN_DIR, "bin", "protoc-gen-js")
    if not os.path.exists(protoc_gen_js):
        sys_name = platform.system().lower()
        machine = platform.machine().lower()
        if sys_name == "darwin":
            arch = "aarch_64" if machine in ("arm64", "aarch64") else "x86_64"
            asset = f"protobuf-javascript-{PB_VERSION}-osx-{arch}.tar.gz"
        elif sys_name == "linux":
            arch = "aarch_64" if machine in ("arm64", "aarch64") else "x86_64"
            asset = f"protobuf-javascript-{PB_VERSION}-linux-{arch}.tar.gz"
        else:
            raise RuntimeError(f"Unsupported platform: {sys_name} {machine}")
        url = f"https://github.com/protocolbuffers/protobuf-javascript/releases/download/v{PB_VERSION}/{asset}"
        print(f"[setup] Downloading {url} ...")
        os.makedirs(PB_BIN_DIR, exist_ok=True)
        tar_path = os.path.join(CACHE_DIR, asset)
        urllib.request.urlretrieve(url, tar_path)
        subprocess.run(["tar", "-xzf", tar_path, "-C", PB_BIN_DIR], check=True)
        os.remove(tar_path)

    # Ensure protobuf-javascript Closure source files
    if not os.path.exists(os.path.join(PB_SRC_DIR, "message.js")):
        url = f"https://github.com/protocolbuffers/protobuf-javascript/archive/refs/tags/v{PB_VERSION}.tar.gz"
        print(f"[setup] Downloading {url} ...")
        os.makedirs(PB_SRC_DIR, exist_ok=True)
        tar_path = os.path.join(CACHE_DIR, f"protobuf-javascript-src-{PB_VERSION}.tar.gz")
        urllib.request.urlretrieve(url, tar_path)
        subprocess.run(
            ["tar", "-xzf", tar_path, "--strip-components=1", "-C", PB_SRC_DIR],
            check=True,
        )
        os.remove(tar_path)

    # Ensure googleapis repo
    if not os.path.exists(GOOGLEAPIS_DIR):
        print("[setup] Cloning https://github.com/googleapis/googleapis.git ...")
        subprocess.run(
            [
                "git",
                "clone",
                "--depth=1",
                "https://github.com/googleapis/googleapis.git",
                GOOGLEAPIS_DIR,
            ],
            check=True,
        )

    if not os.path.exists(CLOSURE_LIB_DIR):
        raise RuntimeError(
            f"Missing {CLOSURE_LIB_DIR}; run tools/bench/fetch-inputs.sh first"
        )

    return protoc, protoc_gen_js, protoc_inc


def generate_proto_js(protoc: str, protoc_gen_js: str, protoc_inc: str, proto_out_dir: str) -> dict[str, list[str]]:
    """Run protoc once per proto directory into its own output directory."""
    proto_dirs: dict[str, list[str]] = {}
    for root, dirs, files in os.walk(GOOGLEAPIS_DIR):
        if ".git" in dirs:
            dirs.remove(".git")
        rel_dir = os.path.relpath(root, GOOGLEAPIS_DIR)
        # Skip preview/ which duplicates google/{api,cloud,rpc,type} with identical proto packages
        if rel_dir == "preview" or rel_dir.startswith("preview/"):
            dirs[:] = []
            continue
        protos = sorted(os.path.join(root, f) for f in files if f.endswith(".proto"))
        if protos:
            proto_dirs[rel_dir] = protos

    marker_file = os.path.join(proto_out_dir, ".generated_ok")
    if os.path.exists(marker_file):
        print(f"[setup] Reusing existing generated proto JS in {proto_out_dir}")
        return proto_dirs

    if os.path.exists(proto_out_dir):
        shutil.rmtree(proto_out_dir)
    os.makedirs(proto_out_dir, exist_ok=True)

    t0 = time.time()
    print(f"[setup] Generating Closure-style protobuf JS across {len(proto_dirs)} directories + WKT ...")

    # 1. Well-known types
    wkt_out = os.path.join(proto_out_dir, "_wkt", "google", "protobuf")
    os.makedirs(wkt_out, exist_ok=True)
    wkt_inputs = [os.path.join(protoc_inc, p) for p in WELL_KNOWN_PROTOS]
    subprocess.run(
        [
            protoc,
            f"--plugin=protoc-gen-js={protoc_gen_js}",
            "--experimental_allow_proto3_optional",
            f"-I{protoc_inc}",
            f"--js_out=import_style=closure,binary:{wkt_out}",
        ]
        + wkt_inputs,
        check=True,
    )

    # 2. Once per proto directory in googleapis
    def run_one(item: tuple[str, list[str]]) -> tuple[str, int, str]:
        rel_dir, protos = item
        out_dir = os.path.join(proto_out_dir, rel_dir)
        os.makedirs(out_dir, exist_ok=True)
        cmd = [
            protoc,
            f"--plugin=protoc-gen-js={protoc_gen_js}",
            "--experimental_allow_proto3_optional",
            f"-I{GOOGLEAPIS_DIR}",
            f"-I{protoc_inc}",
            f"--js_out=import_style=closure,binary:{out_dir}",
        ] + protos
        res = subprocess.run(cmd, capture_output=True, text=True)
        return rel_dir, res.returncode, res.stderr

    failures = []
    with concurrent.futures.ThreadPoolExecutor(max_workers=os.cpu_count() or 8) as ex:
        for rel_dir, rc, stderr in ex.map(run_one, sorted(proto_dirs.items())):
            if rc != 0:
                failures.append((rel_dir, rc, stderr))

    if failures:
        for rel_dir, rc, stderr in failures[:5]:
            print(f"ERROR in {rel_dir} (rc={rc}): {stderr}", file=sys.stderr)
        raise RuntimeError(f"protoc failed on {len(failures)} directories")

    # Remove protoc-gen-js placeholder empty_no_content_void_file.js files
    for root, _, files in os.walk(proto_out_dir):
        for f in files:
            if f == "empty_no_content_void_file.js":
                os.remove(os.path.join(root, f))

    with open(marker_file, "w") as fh:
        fh.write("ok\n")

    print(f"[setup] Generated proto JS in {time.time() - t0:.1f}s")
    return proto_dirs


def sanitize_chunk_name(prefix: str, idx: int, name: str) -> str:
    clean = re.sub(r"[^a-zA-Z0-9_]", "_", name).strip("_")
    return f"{prefix}_{idx}_{clean}"


def build_closure_and_runtime_chunks() -> tuple[list[tuple[str, list[str], list[str]]], list[str], dict[str, int]]:
    """Split Closure Library + protobuf runtime along its goog.require graph."""
    base_js = os.path.join(CLOSURE_LIB_DIR, "closure", "goog", "base.js")
    closure_files = []
    for pat in ["closure/goog/**/*.js", "third_party/closure/goog/**/*.js"]:
        for p in sorted(glob.glob(os.path.join(CLOSURE_LIB_DIR, pat), recursive=True)):
            if p.endswith("_test.js"):
                continue
            closure_files.append(p)

    pb_files = [os.path.join(PB_SRC_DIR, f) for f in PB_RUNTIME_REL_FILES]

    prov_re = re.compile(
        r"\bgoog\.(?:provide|module|declareModuleId)\(\s*['\"]([^'\"]+)['\"]\s*\)"
    )
    req_re = re.compile(
        r"\bgoog\.(?:require|requireType)\(\s*['\"]([^'\"]+)['\"]\s*\)"
    )

    provider: dict[str, str] = {}
    file_reqs: dict[str, list[str]] = {}
    for p in closure_files + pb_files:
        with open(p, "r", encoding="utf-8", errors="ignore") as fh:
            txt = fh.read()
        txt_nc = re.sub(r"/\*.*?\*/", "", txt, flags=re.S)
        txt_nc = re.sub(r"//.*", "", txt_nc)
        for pr in prov_re.findall(txt_nc):
            provider[pr] = p
        file_reqs[p] = req_re.findall(txt_nc)

    visited: set[str] = set()
    stack = list(pb_files)
    while stack:
        curr = stack.pop()
        if curr in visited:
            continue
        visited.add(curr)
        for r in file_reqs.get(curr, []):
            if r in provider:
                stack.append(provider[r])

    deps: dict[str, set[str]] = {p: set() for p in visited}
    for p in visited:
        for r in file_reqs.get(p, []):
            if r in provider and provider[r] in visited and provider[r] != p:
                deps[p].add(provider[r])

    topo: list[str] = []
    done: set[str] = set()

    def visit(u: str) -> None:
        if u in done:
            return
        for v in sorted(deps[u]):
            visit(v)
        done.add(u)
        topo.append(u)

    for p in sorted(visited):
        visit(p)

    chunks: list[tuple[str, list[str], list[str]]] = [("base", [base_js], [])]
    file_to_chunk: dict[str, str] = {base_js: "base"}
    chunk_depth: dict[str, int] = {"base": 1}

    for i, p in enumerate(topo, start=1):
        stem = os.path.splitext(os.path.basename(p))[0]
        cname = sanitize_chunk_name("rt", i, stem)
        cdeps = sorted({file_to_chunk[v] for v in deps[p]}) if deps[p] else ["base"]
        chunks.append((cname, [p], cdeps))
        file_to_chunk[p] = cname
        chunk_depth[cname] = 1 + max(chunk_depth[d] for d in cdeps)

    rt_anchor_chunks = sorted(
        {
            file_to_chunk[provider[sym]]
            for sym in (
                "jspb.Message",
                "jspb.BinaryReader",
                "jspb.BinaryWriter",
                "jspb.Map",
            )
        }
    )
    return chunks, rt_anchor_chunks, chunk_depth


def main() -> None:
    t_start = time.time()
    protoc, protoc_gen_js, protoc_inc = ensure_prerequisites()

    proto_out_dir = os.path.join(BENCH_DIR, "proto_out")
    stubs_dir = os.path.join(BENCH_DIR, "stubs")
    keepalive_dir = os.path.join(BENCH_DIR, "keepalive")
    os.makedirs(stubs_dir, exist_ok=True)
    os.makedirs(keepalive_dir, exist_ok=True)

    proto_dirs = generate_proto_js(protoc, protoc_gen_js, protoc_inc, proto_out_dir)

    # 1. Build Closure Library + protobuf runtime chunk DAG
    rt_chunks, rt_anchor_chunks, chunk_depth = build_closure_and_runtime_chunks()
    print(
        f"[setup] Built {len(rt_chunks)} Closure Library + protobuf runtime chunks "
        f"(max depth: {max(chunk_depth.values())})"
    )

    # 2. Parse all .proto files and their imports
    import_re = re.compile(r'^\s*import\s+(?:public\s+|weak\s+)?"([^"]+)"\s*;', re.M)
    proto_imports: dict[str, list[str]] = {}

    for w in WELL_KNOWN_PROTOS:
        with open(os.path.join(protoc_inc, w), "r", encoding="utf-8") as fh:
            proto_imports[w] = import_re.findall(fh.read())

    for rel_dir, protos in sorted(proto_dirs.items()):
        for p in protos:
            rel_p = os.path.relpath(p, GOOGLEAPIS_DIR)
            with open(p, "r", encoding="utf-8") as fh:
                proto_imports[rel_p] = import_re.findall(fh.read())

    # 3. Scan generated .js files and group by // source: <proto_path>
    source_re = re.compile(r"^// source: (.+)$", re.M)
    provide_re = re.compile(r"^goog\.provide\('([^']+)'\);", re.M)
    require_re = re.compile(r"^goog\.require\('([^']+)'\);", re.M)
    msg_re = re.compile(
        r"^(proto\.[A-Za-z0-9_.]+)\.prototype\.serializeBinary\s*=", re.M
    )

    proto_js_files: dict[str, list[str]] = {p: [] for p in proto_imports}
    file_provs: dict[str, list[str]] = {}
    file_reqs: dict[str, list[str]] = {}
    proto_msgs: dict[str, list[str]] = {p: [] for p in proto_imports}

    print("[setup] Indexing generated proto JS files ...")
    for root, _, files in os.walk(proto_out_dir):
        for f in sorted(files):
            if not f.endswith(".js"):
                continue
            p = os.path.join(root, f)
            with open(p, "r", encoding="utf-8") as fh:
                content = fh.read()
            m = source_re.search(content)
            if not m:
                continue
            src_proto = m.group(1).strip()
            if src_proto not in proto_js_files:
                continue
            proto_js_files[src_proto].append(p)
            file_provs[p] = provide_re.findall(content)
            file_reqs[p] = require_re.findall(content)
            for msg in msg_re.findall(content):
                proto_msgs[src_proto].append(msg)

    # Topologically sort JS files within each .proto chunk
    for src_proto, js_list in proto_js_files.items():
        if len(js_list) <= 1:
            continue
        prov_map = {pr: p for p in js_list for pr in file_provs[p]}
        intra_deps = {
            p: {prov_map[r] for r in file_reqs[p] if r in prov_map and prov_map[r] != p}
            for p in js_list
        }
        ordered: list[str] = []
        seen: set[str] = set()

        def visit_js(u: str) -> None:
            if u in seen:
                return
            seen.add(u)
            for v in sorted(intra_deps[u]):
                visit_js(v)
            ordered.append(u)

        for p in sorted(js_list):
            visit_js(p)
        proto_js_files[src_proto] = ordered

    # Topologically sort .proto chunks
    proto_topo: list[str] = []
    proto_seen: set[str] = set()

    def visit_proto(u: str) -> None:
        if u in proto_seen:
            return
        proto_seen.add(u)
        for imp in sorted(proto_imports.get(u, [])):
            if imp in proto_imports and imp != u:
                visit_proto(imp)
        proto_topo.append(u)

    for p in sorted(proto_imports):
        visit_proto(p)

    # Identify leaf .proto files (no other .proto imports them)
    dependents: dict[str, set[str]] = {p: set() for p in proto_imports}
    for p, imps in proto_imports.items():
        for imp in imps:
            if imp in proto_imports and imp != p:
                dependents[imp].add(p)
    leaf_protos = {p for p, deps in dependents.items() if len(deps) == 0}

    # Build .proto chunk definitions and keep-alive modules for leaves
    proto_to_chunk: dict[str, str] = {}
    proto_chunks: list[tuple[str, list[str], list[str]]] = []
    keepalive_count = 0
    stub_count = 0

    for i, p in enumerate(proto_topo, start=1):
        stem = os.path.splitext(p)[0]
        cname = sanitize_chunk_name("p", i, stem)
        proto_to_chunk[p] = cname

        imp_chunks = sorted(
            {proto_to_chunk[imp] for imp in proto_imports[p] if imp in proto_to_chunk and imp != p}
        )
        cdeps = sorted(set(rt_anchor_chunks) | set(imp_chunks))
        chunk_depth[cname] = 1 + max(chunk_depth[d] for d in cdeps)

        cfiles = list(proto_js_files[p])

        if p in leaf_protos:
            # Keep-alive module for leaf chunk
            msgs = list(proto_msgs[p])
            if not msgs:
                for imp in proto_imports[p]:
                    msgs.extend(proto_msgs.get(imp, []))
            selected_msgs = msgs[:KEEPALIVE_MESSAGES_PER_LEAF]
            ka_path = os.path.join(keepalive_dir, f"{cname}.js")
            lines = [
                "/**",
                f" * @fileoverview Keep-alive module for leaf chunk {cname} ({p}).",
                " * @suppress {missingRequire,extraRequire,checkTypes}",
                " */",
                f"goog.provide('keepalive.{cname}');",
            ]
            for m in selected_msgs:
                lines.append(f"goog.require('{m}');")
            lines.append("(function() {")
            lines.append("  var kept = [")
            for idx, m in enumerate(selected_msgs):
                comma = "," if idx + 1 < len(selected_msgs) else ""
                lines.append(
                    f"    {m}.deserializeBinary((new {m}()).serializeBinary()){comma}"
                )
            lines.append("  ];")
            lines.append(f"  goog.exportSymbol('keepalive.{cname}', kept);")
            lines.append("})();")
            with open(ka_path, "w", encoding="utf-8") as fh:
                fh.write("\n".join(lines) + "\n")
            cfiles.append(ka_path)
            keepalive_count += 1
        elif not cfiles:
            # Non-leaf proto with 0 messages (e.g. google/api/annotations.proto)
            stub_path = os.path.join(stubs_dir, f"{cname}.js")
            with open(stub_path, "w", encoding="utf-8") as fh:
                fh.write(f"// stub for message-less proto: {p}\n")
            cfiles.append(stub_path)
            stub_count += 1

        proto_chunks.append((cname, cfiles, cdeps))

    all_chunks = rt_chunks + proto_chunks
    total_js_files = sum(len(cfiles) for _, cfiles, _ in all_chunks)

    # Write flags.txt
    flagfile_path = os.path.join(BENCH_DIR, "flags.txt")
    common_flags = [
        "--compilation_level=ADVANCED",
        "--language_in=ECMASCRIPT_2020",
        "--language_out=ECMASCRIPT5",
        "--warning_level=DEFAULT",
        "--use_types_for_optimization",
        "--jscomp_warning=checkTypes",
        "--jscomp_off=reportUnknownTypes",
        "--jscomp_off=strictCheckTypes",
        "--jscomp_off=lintChecks",
        "--jscomp_off=deprecated",
        "--jscomp_off=missingPolyfill",
        "--jscomp_off=visibility",
        "--jscomp_off=extraRequire",
        "--jscomp_off=missingRequire",
        "--jscomp_off=duplicate",
        "--jscomp_off=suspiciousCode",
        "--generate_exports",
        "-D",
        "COMPILED=true",
        "-D",
        "goog.DEBUG=false",
        "-D",
        "goog.asserts.ENABLE_ASSERTS=false",
    ]

    with open(flagfile_path, "w", encoding="utf-8") as fh:
        for fl in common_flags:
            fh.write(fl + "\n")
        for cname, cfiles, cdeps in all_chunks:
            for cf in cfiles:
                fh.write(f"--js={cf}\n")
        for cname, cfiles, cdeps in all_chunks:
            dep_str = (":" + ",".join(cdeps)) if cdeps else ""
            fh.write(f"--chunk={cname}:{len(cfiles)}{dep_str}\n")

    manifest = {
        "total_chunks": len(all_chunks),
        "runtime_And_closure_chunks": len(rt_chunks),
        "proto_chunks": len(proto_chunks),
        "leaf_proto_chunks": len(leaf_protos),
        "keepalive_files": keepalive_count,
        "stub_files": stub_count,
        "total_js_files": total_js_files,
        "max_chunk_depth": max(chunk_depth.values()),
        "flagfile": flagfile_path,
    }
    manifest_path = os.path.join(BENCH_DIR, "manifest.json")
    with open(manifest_path, "w", encoding="utf-8") as fh:
        json.dump(manifest, fh, indent=2)

    print(f"[setup] Done in {time.time() - t_start:.1f}s:")
    print(json.dumps(manifest, indent=2))


if __name__ == "__main__":
    main()
