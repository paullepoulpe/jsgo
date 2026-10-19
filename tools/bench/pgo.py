#!/usr/bin/env python3
"""Regenerates cmd/jsgo/default.pgo, the CPU profile `go build` uses for
profile-guided optimisation of jsgo (see tools/bench/README.md, "PGO profile").

  python3 -I tools/bench/pgo.py [--profile cmd/jsgo/default.pgo] [--cache DIR]

It sets up the benchmark's task workspaces and invocations the way bench.py
does (the pinned inputs must have been fetched: tools/bench/fetch-inputs.sh),
builds jsgo with a CPU-profiling hook (through -overlay; the sources are not
changed), runs the first invocation of each whole-program task and an evenly
spaced sample of the single-file task, and merges the profiles into one. The
build ignores any existing default.pgo, so the profile describes the
unoptimised code paths.
"""
import argparse
import json
import os
import shutil
import subprocess
import sys
import tempfile

HERE = os.path.dirname(os.path.abspath(__file__))
REPO = os.path.dirname(os.path.dirname(HERE))
sys.path.insert(0, HERE)  # python3 -I leaves the script's directory out
sys.dont_write_bytecode = True
import bench  # noqa: E402
GO_ENV_DROP = ("GOGC", "GOMEMLIMIT", "GOMAXPROCS", "GODEBUG")

HOOK = '''package jutil

// ProfExitHook runs before System_Exit ends the process (tools/bench/pgo.py).
var ProfExitHook func()
'''

MAIN_PROF = '''package main

import (
	"os"
	"runtime/pprof"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func init() {
	f, err := os.Create(os.Getenv("JSGO_PGO_PROFILE"))
	if err != nil {
		panic(err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	jutil.ProfExitHook = func() {
		pprof.StopCPUProfile()
		f.Close()
	}
}
'''


def overlay(tmp):
    """Writes the overlay files and returns the -overlay JSON path."""
    system = os.path.join(REPO, "internal", "jutil", "system.go")
    src = open(system).read()
    head = "func System_Exit(status int) {\n"
    if src.count(head) != 1:
        sys.exit("internal/jutil/system.go: System_Exit not found")
    src = src.replace(head, head + "\tif ProfExitHook != nil {\n\t\tProfExitHook()\n\t}\n")
    main = os.path.join(REPO, "cmd", "jsgo", "main.go")
    msrc = open(main).read()
    call = "\tjscomp.CommandLineRunner_Main(os.Args[1:])\n"
    if msrc.count(call) != 1:
        sys.exit("cmd/jsgo/main.go: the CommandLineRunner_Main call not found")
    msrc = msrc.replace(call, call + "\tif jutil.ProfExitHook != nil {\n\t\tjutil.ProfExitHook()\n\t}\n")
    files = {
        system: src,
        os.path.join(REPO, "internal", "jutil", "zz_pgo_hook.go"): HOOK,
        main: msrc,
        os.path.join(REPO, "cmd", "jsgo", "zz_pgo_main.go"): MAIN_PROF,
    }
    replace = {}
    for i, (path, text) in enumerate(files.items()):
        p = os.path.join(tmp, f"{i}.go")
        open(p, "w").write(text)
        replace[path] = p
    o = os.path.join(tmp, "overlay.json")
    json.dump({"Replace": replace}, open(o, "w"))
    return o


def main():
    ap = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    ap.add_argument("--cache", default=bench.default_cache(), help="the .cache directory (default: the main checkout's)")
    ap.add_argument("--profile", default=os.path.join(REPO, "cmd", "jsgo", "default.pgo"))
    ap.add_argument("--tasks", default="typed-app,bundle,downlevel-whitespace,downlevel-simple,checks-only")
    ap.add_argument("--single-files", type=int, default=8, help="single-file invocations to profile")
    a = ap.parse_args()

    env = {k: v for k, v in os.environ.items() if k not in GO_ENV_DROP}
    env["LC_ALL"] = "C.UTF-8"
    with tempfile.TemporaryDirectory() as tmp:
        binary = os.path.join(tmp, "jsgo-prof")
        subprocess.run(["go", "build", "-pgo=off", "-overlay", overlay(tmp), "-o", binary, "./cmd/jsgo"],
                       cwd=REPO, check=True)
        ws = os.path.join(tmp, "ws")
        bench.make_workspace(ws, a.cache)
        tasks = {t["name"]: t for t in bench.load_tasks()}
        runs = []
        for task in a.tasks.split(","):
            runs.append((task, ws, bench.task_jobs(tasks[task], ws)[0], 1800))
        if a.single_files:
            jobs = bench.task_jobs(tasks["single-files"], ws)
            step = max(1, len(jobs) // a.single_files)
            for i in range(0, len(jobs), step)[:a.single_files]:
                # A compile that deadlocks hangs while the profiler runs.
                runs.append((f"single-files-{i}", ws, jobs[i], 60))
        profiles = []
        for name, ws, args, timeout in runs:
            outdir = os.path.join(tmp, name)
            os.makedirs(outdir)
            # Outputs go to the scratch directory, not the workspace.
            args = [x.replace("--js_output_file=out.js", f"--js_output_file={outdir}/out.js")
                     .replace("--create_source_map=out.js.map", f"--create_source_map={outdir}/out.js.map")
                    for x in args]
            prof = os.path.join(tmp, name + ".pprof")
            try:
                subprocess.run([binary] + args, cwd=ws, env=dict(env, JSGO_PGO_PROFILE=prof),
                               stdin=subprocess.DEVNULL, stdout=subprocess.DEVNULL,
                               stderr=subprocess.DEVNULL, timeout=timeout)
            except subprocess.TimeoutExpired:
                print(f"{name}: timed out, skipped")
                continue
            if os.path.exists(prof) and os.path.getsize(prof) > 0:
                profiles.append(prof)
                print(f"{name}: profiled")
            else:
                print(f"{name}: no profile, skipped")
        merged = os.path.join(tmp, "merged.pgo")
        with open(merged, "wb") as f:
            subprocess.run(["go", "tool", "pprof", "-proto"] + profiles, stdout=f, check=True)
        shutil.copyfile(merged, a.profile)
    print(f"wrote {a.profile} from {len(profiles)} profiles")


if __name__ == "__main__":
    main()
