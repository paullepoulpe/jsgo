#!/usr/bin/env python3
"""Run the 3-stage, 7.5k-chunk googleapis Closure Compiler benchmark.

Stages:
  1. CHECKS:
     --flagfile=<flags.txt>
     --segment_of_compilation_to_run=CHECKS
     --filename_to_save_to=<work_dir>/checks.state
  2. OPTIMIZATIONS:
     --flagfile=<flags.txt>
     --segment_of_compilation_to_run=OPTIMIZATIONS
     --filename_to_restore_from=<work_dir>/checks.state
     --filename_to_save_to=<work_dir>/optimizations.state
  3. FINALIZATIONS:
     --flagfile=<flags.txt>
     --segment_of_compilation_to_run=FINALIZATIONS
     --filename_to_restore_from=<work_dir>/optimizations.state
     --chunk_output_path_prefix=<work_dir>/chunks/
"""

from __future__ import annotations

import argparse
import glob
import hashlib
import json
import os
import resource
import shutil
import subprocess
import sys
import time

ROOT_DIR = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
BENCH_DIR = os.path.join(ROOT_DIR, ".cache", "bench-inputs", "googleapis-bench")


def run_stage(
    stage_name: str,
    cmd: list[str],
    env: dict[str, str],
    log_path: str,
) -> dict[str, object]:
    print(f"[{stage_name}] Running: {' '.join(cmd[:6])} ...")
    t0 = time.perf_counter()
    with open(log_path, "w", encoding="utf-8") as log_fh:
        proc = subprocess.Popen(
            cmd,
            stdout=log_fh,
            stderr=subprocess.STDOUT,
            env=env,
        )
        _, status, rusage = os.wait4(proc.pid, 0)
    wall_s = time.perf_counter() - t0
    rc = os.waitstatus_to_exitcode(status)

    # On macOS ru_maxrss is in bytes; on Linux it is in KiB.
    maxrss_bytes = (
        rusage.ru_maxrss
        if sys.platform == "darwin"
        else rusage.ru_maxrss * 1024
    )
    maxrss_mib = maxrss_bytes / (1024 * 1024)
    user_s = rusage.ru_utime
    sys_s = rusage.ru_stime

    result = {
        "stage": stage_name,
        "exit_code": rc,
        "wall_s": round(wall_s, 3),
        "user_s": round(user_s, 3),
        "sys_s": round(sys_s, 3),
        "cpu_s": round(user_s + sys_s, 3),
        "max_rss_mib": round(maxrss_mib, 1),
        "log_path": log_path,
    }
    print(
        f"[{stage_name}] exit={rc} wall={wall_s:.2f}s cpu={user_s + sys_s:.2f}s "
        f"rss={maxrss_mib:.1f} MiB"
    )
    if rc != 0:
        with open(log_path, "r", encoding="utf-8", errors="replace") as fh:
            tail = "".join(fh.readlines()[-40:])
        print(f"[{stage_name}] FAILED:\n{tail}", file=sys.stderr)
        raise RuntimeError(f"Stage {stage_name} failed with exit code {rc}")
    return result


def main() -> None:
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument(
        "--compiler",
        default=os.path.join(ROOT_DIR, "jsgo"),
        help="Path to jsgo binary (default: ./jsgo)",
    )
    ap.add_argument(
        "--java",
        action="store_true",
        help="Run Java Closure Compiler jar instead of jsgo",
    )
    ap.add_argument(
        "--java-bin",
        default="/opt/homebrew/opt/openjdk/bin/java"
        if os.path.exists("/opt/homebrew/opt/openjdk/bin/java")
        else "java",
        help="Path to Java >=21 binary",
    )
    ap.add_argument(
        "--java-jar",
        default=os.path.join(ROOT_DIR, ".cache", "closure-compiler-v20261015-maven.jar"),
        help="Path to Java Closure Compiler jar",
    )
    ap.add_argument(
        "--java-heap",
        default="20g",
        help="JVM -Xmx heap limit when --java is set (default: 20g)",
    )
    ap.add_argument(
        "--threads",
        type=int,
        default=1,
        help="--num_parallel_threads value (default: 1)",
    )
    ap.add_argument(
        "--fast-mode",
        action="store_true",
        help="Pass --fast_mode to the compiler",
    )
    ap.add_argument(
        "--tracer-mode",
        default="OFF",
        help="--tracer_mode value (e.g. OFF, AST_SIZE, TIMING_ONLY, ALL)",
    )
    ap.add_argument(
        "--stages",
        default="CHECKS,OPTIMIZATIONS,FINALIZATIONS",
        help="Comma-separated stages to run (default: CHECKS,OPTIMIZATIONS,FINALIZATIONS)",
    )
    ap.add_argument(
        "--work-dir",
        default=os.path.join(BENCH_DIR, "work"),
        help="Working directory for state files and chunk outputs",
    )
    ap.add_argument(
        "--keep-chunks",
        action="store_true",
        help="Keep generated chunk output .js files on disk after measuring",
    )
    ap.add_argument(
        "--keep-states",
        action="store_true",
        help="Keep checks.state and optimizations.state on disk after completion",
    )
    ap.add_argument(
        "--profile-dir",
        default="",
        help="Optional directory to write per-stage CPU and heap profiles",
    )
    args = ap.parse_args()

    flagfile = os.path.join(BENCH_DIR, "flags.txt")
    if not os.path.exists(flagfile):
        subprocess.run(
            [sys.executable, os.path.join(os.path.dirname(__file__), "setup_googleapis_bench.py")],
            check=True,
        )

    os.makedirs(args.work_dir, exist_ok=True)
    checks_state = os.path.join(args.work_dir, "checks.state")
    opt_state = os.path.join(args.work_dir, "optimizations.state")
    chunks_dir = os.path.join(args.work_dir, "chunks")

    if args.java:
        base_cmd = [
            args.java_bin,
            f"-Xmx{args.java_heap}",
            "-jar",
            args.java_jar,
            f"--flagfile={flagfile}",
        ]
    else:
        base_cmd = [
            args.compiler,
            f"--flagfile={flagfile}",
            f"--num_parallel_threads={args.threads}",
        ]
        if args.fast_mode:
            base_cmd.append("--fast_mode")
    if args.tracer_mode and args.tracer_mode != "OFF":
        base_cmd.append(f"--tracer_mode={args.tracer_mode}")

    env = os.environ.copy()
    stages = [s.strip().upper() for s in args.stages.split(",") if s.strip()]
    summary: dict[str, object] = {
        "compiler": args.java_jar if args.java else args.compiler,
        "java": args.java,
        "threads": args.threads if not args.java else "jvm_default",
        "fast_mode": args.fast_mode if not args.java else False,
        "tracer_mode": args.tracer_mode,
        "stages": {},
    }

    total_wall = 0.0
    peak_rss = 0.0

    if "CHECKS" in stages:
        stage_env = env.copy()
        if args.profile_dir:
            os.makedirs(args.profile_dir, exist_ok=True)
            stage_env["JSGO_CPU_PROFILE"] = os.path.join(args.profile_dir, "checks-cpu.pprof")
            stage_env["JSGO_MEM_PROFILE"] = os.path.join(args.profile_dir, "checks-heap.pprof")
        res = run_stage(
            "CHECKS",
            base_cmd
            + [
                "--segment_of_compilation_to_run=CHECKS",
                f"--filename_to_save_to={checks_state}",
            ],
            stage_env,
            os.path.join(args.work_dir, "checks.log"),
        )
        if os.path.exists(checks_state):
            res["state_size_mib"] = round(
                os.path.getsize(checks_state) / (1024 * 1024), 2
            )
            res["state_sha256"] = hashlib.sha256(open(checks_state, "rb").read()).hexdigest()
            print(f"[CHECKS] checks.state size: {res['state_size_mib']} MiB (sha256: {res['state_sha256'][:16]}...)")
        summary["stages"]["CHECKS"] = res
        total_wall += float(res["wall_s"])
        peak_rss = max(peak_rss, float(res["max_rss_mib"]))

    if "OPTIMIZATIONS" in stages:
        stage_env = env.copy()
        if args.profile_dir:
            os.makedirs(args.profile_dir, exist_ok=True)
            stage_env["JSGO_CPU_PROFILE"] = os.path.join(args.profile_dir, "optimizations-cpu.pprof")
            stage_env["JSGO_MEM_PROFILE"] = os.path.join(args.profile_dir, "optimizations-heap.pprof")
        res = run_stage(
            "OPTIMIZATIONS",
            base_cmd
            + [
                "--segment_of_compilation_to_run=OPTIMIZATIONS",
                f"--filename_to_restore_from={checks_state}",
                f"--filename_to_save_to={opt_state}",
            ],
            stage_env,
            os.path.join(args.work_dir, "optimizations.log"),
        )
        if os.path.exists(opt_state):
            res["state_size_mib"] = round(
                os.path.getsize(opt_state) / (1024 * 1024), 2
            )
            res["state_sha256"] = hashlib.sha256(open(opt_state, "rb").read()).hexdigest()
            print(f"[OPTIMIZATIONS] optimizations.state size: {res['state_size_mib']} MiB (sha256: {res['state_sha256'][:16]}...)")
        summary["stages"]["OPTIMIZATIONS"] = res
        total_wall += float(res["wall_s"])
        peak_rss = max(peak_rss, float(res["max_rss_mib"]))

    if "FINALIZATIONS" in stages:
        if os.path.exists(chunks_dir):
            shutil.rmtree(chunks_dir)
        os.makedirs(chunks_dir, exist_ok=True)
        stage_env = env.copy()
        if args.profile_dir:
            os.makedirs(args.profile_dir, exist_ok=True)
            stage_env["JSGO_CPU_PROFILE"] = os.path.join(args.profile_dir, "finalizations-cpu.pprof")
            stage_env["JSGO_MEM_PROFILE"] = os.path.join(args.profile_dir, "finalizations-heap.pprof")
        res = run_stage(
            "FINALIZATIONS",
            base_cmd
            + [
                "--segment_of_compilation_to_run=FINALIZATIONS",
                f"--filename_to_restore_from={opt_state}",
                f"--chunk_output_path_prefix={chunks_dir}/",
            ],
            stage_env,
            os.path.join(args.work_dir, "finalizations.log"),
        )
        chunk_files = sorted(glob.glob(os.path.join(chunks_dir, "*.js")))
        total_out_bytes = 0
        non_empty_chunks = 0
        composite_hash = hashlib.sha256()
        for p in chunk_files:
            data = open(p, "rb").read()
            sz = len(data)
            total_out_bytes += sz
            if sz > 1:
                non_empty_chunks += 1
            composite_hash.update(os.path.basename(p).encode("utf-8"))
            composite_hash.update(b"\0")
            composite_hash.update(hashlib.sha256(data).digest())
        res["output_chunks_count"] = len(chunk_files)
        res["non_empty_chunks_count"] = non_empty_chunks
        res["output_total_mib"] = round(total_out_bytes / (1024 * 1024), 2)
        res["chunks_sha256"] = composite_hash.hexdigest()
        print(
            f"[FINALIZATIONS] Emitted {len(chunk_files)} chunk files "
            f"({non_empty_chunks} non-empty, {res['output_total_mib']} MiB total, sha256={res['chunks_sha256'][:16]}...)"
        )
        if not args.keep_chunks:
            shutil.rmtree(chunks_dir)
        summary["stages"]["FINALIZATIONS"] = res
        total_wall += float(res["wall_s"])
        peak_rss = max(peak_rss, float(res["max_rss_mib"]))

    if not args.keep_states and stages == ["CHECKS", "OPTIMIZATIONS", "FINALIZATIONS"]:
        for sf in (checks_state, opt_state):
            if os.path.exists(sf):
                os.remove(sf)

    summary["total_wall_s"] = round(total_wall, 3)
    summary["peak_rss_mib"] = round(peak_rss, 1)

    results_path = os.path.join(args.work_dir, "results.json")
    with open(results_path, "w", encoding="utf-8") as fh:
        json.dump(summary, fh, indent=2)
    print("\n=== Summary ===")
    print(json.dumps(summary, indent=2))


if __name__ == "__main__":
    main()
