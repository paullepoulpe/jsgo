#!/usr/bin/env python3
"""Materializes per-file .i.js type summaries in .cache/bench-inputs/ijs-typed-full/.

Used by the `incremental-ijs-check` benchmark task to model incremental Blaze/Bazel
`JsChecker` type checking where a small leaf `--js` target is checked against ~1,740
pre-generated `.i.js` (`@typeSummary`) dependency files.

Idempotent: keyed on the SHA-256 of the pinned Java Closure Compiler jar in
tools/oracle/pin.env. Invoked automatically by tools/bench/fetch-inputs.sh.
"""
import hashlib
import json
import os
import re
import shutil
import subprocess
import sys
import tempfile

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, HERE)
import bench  # noqa: E402


def resolve_java():
    env_java = os.environ.get("JSGO_BENCH_JAVA")
    if env_java:
        return env_java
    brew_java = "/opt/homebrew/opt/openjdk/bin/java"
    if os.path.exists(brew_java):
        return brew_java
    return shutil.which("java") or "java"


def main():
    cache = bench.default_cache()
    out_root = sys.argv[1] if len(sys.argv) > 1 else os.path.join(cache, "bench-inputs")
    dest = os.path.join(out_root, "ijs-typed-full")
    jar = os.path.join(cache, bench.read_pin()["CLOSURE_JAR"])
    if not os.path.exists(jar):
        print(f"setup_ijs_bench: skipping ({jar} not present yet; run tools/oracle/fetch.sh)", file=sys.stderr)
        return

    jar_sha = bench.sha256_file(jar)
    marker = os.path.join(dest, ".done")
    if os.path.exists(marker):
        with open(marker) as fh:
            if fh.read().strip() == jar_sha:
                return

    with open(os.path.join(HERE, "tasks", "typed-full.json")) as fh:
        tf = json.load(fh)

    tmp_dest = dest + ".tmp"
    shutil.rmtree(tmp_dest, ignore_errors=True)
    os.makedirs(tmp_dest, exist_ok=True)

    with tempfile.TemporaryDirectory() as ws:
        os.symlink(out_root, os.path.join(ws, "inputs"))
        os.symlink(os.path.join(HERE, "app"), os.path.join(ws, "app"))
        lib_specs = [s for s in tf["inputs"] if not s.get("file", "").startswith("app/")]
        raw_args = bench.expand_inputs(lib_specs, ws)
        files = raw_args[1::2]

        js_files = []
        for f in files:
            if "/externs/" in f or f.endswith("_externs.js"):
                continue
            with open(os.path.join(ws, f), encoding="utf-8", errors="replace") as fh:
                if "@externs" in fh.read():
                    continue
            js_files.append(f)

        all_ijs = os.path.join(ws, "all.i.js")
        cmd = [
            resolve_java(), "-jar", jar,
            "--checks_only",
            "--incremental_check_mode=GENERATE_IJS",
            "--preserve_type_annotations",
            "--formatting=PRETTY_PRINT",
            "--formatting=PRINT_INPUT_DELIMITER",
            "--language_out=ECMASCRIPT_NEXT",
            "--env=CUSTOM",
            "--dependency_mode=NONE",
            "--jscomp_off=lintChecks",
            "--extra_annotation_name=listens",
            "--extra_annotation_name=exportDoc",
            "--extra_annotation_name=exportInterface",
            f"--js_output_file={all_ijs}",
        ]
        for f in js_files:
            cmd += ["--js", f]

        java_env = {k: v for k, v in os.environ.items() if k not in bench.JAVA_ENV_DROP}
        java_env["LC_ALL"] = "C.UTF-8"
        res = subprocess.run(cmd, cwd=ws, env=java_env, capture_output=True, text=True)
        if res.returncode != 0:
            sys.stderr.write(res.stderr)
            sys.exit(f"setup_ijs_bench: GENERATE_IJS failed with exit {res.returncode}")

        with open(all_ijs, encoding="utf-8") as fh:
            content = fh.read()

    matches = list(re.finditer(r"(?m)^// Input (\d+)\n", content))
    # Skip the 5 core language externs that CommandLineRunner prepends in --env=CUSTOM.
    builtin_extern_count = 5
    count = 0
    for idx in range(builtin_extern_count, len(matches)):
        start = matches[idx].end()
        end = matches[idx + 1].start() if idx + 1 < len(matches) else len(content)
        out_path = os.path.join(tmp_dest, f"lib_{count:04d}.i.js")
        with open(out_path, "w", encoding="utf-8") as out_fh:
            out_fh.write(content[start:end])
        count += 1

    with open(os.path.join(tmp_dest, ".done"), "w", encoding="utf-8") as fh:
        fh.write(jar_sha + "\n")
    shutil.rmtree(dest, ignore_errors=True)
    os.rename(tmp_dest, dest)
    print(f"ijs-typed-full: {count} .i.js files")


if __name__ == "__main__":
    main()
