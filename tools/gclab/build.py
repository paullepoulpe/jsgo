#!/usr/bin/env python3
"""Builds jsgo with the gclab profiling hook through `go build -overlay` (the sources are not
changed), like tools/bench/pgo.py does for its CPU-profiling hook. PGO stays on (default.pgo).

  python3 -I build.py SRC OUT_BINARY
"""
import json, os, subprocess, sys

HERE = os.path.dirname(os.path.abspath(__file__))
src, out = os.path.abspath(sys.argv[1]), os.path.abspath(sys.argv[2])
tmp = os.path.join(HERE, "overlay")
os.makedirs(tmp, exist_ok=True)


def patch(path, old, new):
    t = open(path).read()
    if t.count(old) != 1:
        sys.exit(f"{path}: anchor not found exactly once: {old!r}")
    return t.replace(old, new)


files = {}
system = os.path.join(src, "internal/jutil/system.go")
files[system] = patch(system, "func System_Exit(status int) {\n",
                      "func System_Exit(status int) {\n\tif GclabExitHook != nil {\n\t\tGclabExitHook()\n\t}\n")
files[os.path.join(src, "internal/jutil/zz_gclab_hook.go")] = \
    "package jutil\n\n// GclabExitHook runs before System_Exit ends the process (gclab).\nvar GclabExitHook func()\n"
main = os.path.join(src, "cmd/jsgo/main.go")
call = "\tjscomp.CommandLineRunner_Main(os.Args[1:])\n"
files[main] = patch(main, call, call + "\tif jutil.GclabExitHook != nil {\n\t\tjutil.GclabExitHook()\n\t}\n")
tracer = os.path.join(src, "internal/jscomp/tracer.go")
t = patch(tracer, "\tt.startTimeMs = Tracer_clock.CurrentTimeMillis()\n\tt.startThread = jutil.Thread_CurrentThread()\n",
          "\tt.startTimeMs = Tracer_clock.CurrentTimeMillis()\n\tt.startThread = jutil.Thread_CurrentThread()\n"
          "\tif gclabTracerStart != nil {\n\t\tgclabTracerStart(t)\n\t}\n")
files[tracer] = t
t2 = files[tracer]
old = "func (t *Tracer) StopWithThreshold(silenceThreshold int) int64 {\n"
assert t2.count(old) == 1
files[tracer] = t2.replace(old, old + "\tif gclabTracerStop != nil {\n\t\tgclabTracerStop(t)\n\t}\n")
po = os.path.join(src, "internal/jscomp/phase_optimizer.go")
files[po] = patch(po, "\tp.compiler.AfterPass(n.name)\n",
                  "\tp.compiler.AfterPass(n.name)\n\tif gclabAfterPass != nil {\n\t\tgclabAfterPass(n.name, root)\n\t}\n")
files[os.path.join(src, "internal/rhino/zz_gclab_copy.go")] = open(os.path.join(HERE, "hook", "zz_gclab_copy.go")).read()
files[os.path.join(src, "internal/jscomp/zz_gclab.go")] = open(os.path.join(HERE, "hook", "zz_gclab.go")).read()
replace = {}
for i, (path, text) in enumerate(files.items()):
    p = os.path.join(tmp, f"{i}-{os.path.basename(path)}")
    open(p, "w").write(text)
    replace[path] = p
o = os.path.join(tmp, "overlay.json")
json.dump({"Replace": replace}, open(o, "w"), indent=1)
subprocess.run(["go", "build", "-overlay", o, "-o", out, "./cmd/jsgo"], cwd=src, check=True)
print("built", out)
