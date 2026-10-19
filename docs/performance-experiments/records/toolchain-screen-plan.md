# Go 1.27.2 toolchain screen (source audit; no builds performed)

Use one frozen source checkpoint, initially `bd306adb70fb964ca189b88f4d3833480049dd60`,
and the same candidate runtime experiment environment, task arguments and core
budget for every row. Preserve the pinned unmodified warmed-Java target.
Do not mix the later parser integration into only one row.

## Recommended first screen: four binaries

This is a 2×2 comparison: checked-in versus fresh balanced PGO, each with default
inlining versus `-l=4` on jscomp/jutil/rhino. These packages contain the concrete
passes, generic collections and AST/type accessors implicated in current profiles.
Leave the standard library at its normal compiler settings.

Commands below are for the root's next build reservation. `GOGC`/`GOMEMLIMIT` here
limit compiler resource use; they must be removed by the benchmark runner as
usual. They do not bake a Go application GC configuration into the binary.

```bash
git -C /workspace/jsgo-cli-next worktree add --detach /workspace/jsgo-toolchain-screen bd306adb70fb964ca189b88f4d3833480049dd60
source /workspace/toolchains/env.sh
export PATH=/workspace/toolchains/go1.27.2/bin:$PATH
unset GOFLAGS GOEXPERIMENT
export GOAMD64=v1 CGO_ENABLED=1 GOMAXPROCS=2 GOGC=30 GOMEMLIMIT=3GiB
cd /workspace/jsgo-toolchain-screen
mkdir -p /workspace/jsgo-results/toolchain-screen/bin
sha256sum cmd/jsgo/default.pgo
go version
go env GOEXPERIMENT GOAMD64 CGO_ENABLED
```

Existing PGO control:

```bash
go build -buildvcs=false -p 1 -pgo=/workspace/jsgo-toolchain-screen/cmd/jsgo/default.pgo -o /workspace/jsgo-results/toolchain-screen/bin/old-pgo ./cmd/jsgo
```

Existing PGO with reduced non-leaf call penalties:

```bash
go build -buildvcs=false -p 1 -pgo=/workspace/jsgo-toolchain-screen/cmd/jsgo/default.pgo -gcflags='github.com/paullepoulpe/jsgo/internal/jscomp=-l=4' -gcflags='github.com/paullepoulpe/jsgo/internal/jutil=-l=4' -gcflags='github.com/paullepoulpe/jsgo/internal/rhino=-l=4' -o /workspace/jsgo-results/toolchain-screen/bin/old-pgo-l4 ./cmd/jsgo
```

After producing the fresh balanced profile below:

```bash
go build -buildvcs=false -p 1 -pgo=/workspace/jsgo-results/toolchain-screen/balanced.pgo -o /workspace/jsgo-results/toolchain-screen/bin/balanced-pgo ./cmd/jsgo
go build -buildvcs=false -p 1 -pgo=/workspace/jsgo-results/toolchain-screen/balanced.pgo -gcflags='github.com/paullepoulpe/jsgo/internal/jscomp=-l=4' -gcflags='github.com/paullepoulpe/jsgo/internal/jutil=-l=4' -gcflags='github.com/paullepoulpe/jsgo/internal/rhino=-l=4' -o /workspace/jsgo-results/toolchain-screen/bin/balanced-pgo-l4 ./cmd/jsgo
```

Each repeated `-gcflags` pattern applies only to its named package. Do not append
a later matching `-gcflags` entry expecting flags to accumulate: the last match
wins. Use one combined value per package if adding diagnostics or a budget flag.

## Fresh profile construction

Train at the exact frozen source checkpoint, with the same library candidates
enabled as the measured candidate. Build the existing non-tracing main overlay
with explicit `-pgo=off`; this avoids recycling a profile whose callsite offsets
and dominant implementations predate the library redesigns. Regenerate the
overlay's source paths for this worktree before building:

```bash
go build -buildvcs=false -p 1 -pgo=off -overlay=/workspace/jsgo-results/toolchain-screen/profile-overlay/overlay.json -o /workspace/jsgo-results/toolchain-screen/bin/profile ./cmd/jsgo
```

Run ordinary fresh CLI compilations of all five original tasks, retaining CPU
profiles and byte gates. Use enough repetitions to obtain at least roughly
1,000 CPU samples per task, especially for the small whitespace task. Do not
enable `--tracer_mode`, per-node tracing or worker/caching modes. Merge repetitions
within each task first. The following assumes one combined CPU profile per task
at `profiles/<task>/cpu.pprof`.

Plain merging weights each task by its CPU duration. That lets bundle dominate
the PGO decisions. For a balanced task objective, normalize each task to equal
sample/count mass before merging. Go's PGO reader chooses the **first** supported
sample column, normally `samples/count`, rather than automatically preferring
CPU nanoseconds. `pprof -proto -divide_by` scales all columns and truncates integer
values: upscale to one million counts/task to avoid dropping rare samples.

```python
import json, os, re, subprocess
from pathlib import Path

go = '/workspace/toolchains/go1.27.2/bin/go'
base = Path('/workspace/jsgo-results/toolchain-screen')
tasks = ['typed-app', 'bundle', 'downlevel-whitespace', 'downlevel-simple', 'checks-only']
scaled, manifest = [], []
for task in tasks:
    source = base/'profiles'/task/'cpu.pprof'
    raw = subprocess.check_output([go, 'tool', 'pprof', '-raw', str(source)], text=True)
    samples = raw.split('Samples:\n', 1)[1].split('\nLocations', 1)[0]
    assert samples.splitlines()[0].strip() == 'samples/count cpu/nanoseconds'
    count = sum(int(m.group(1)) for m in re.finditer(r'^\s*(\d+)\s+\d+:', samples, re.M))
    assert count > 0
    output = base/(task+'.normalized.pprof')
    subprocess.run([go, 'tool', 'pprof', '-proto', '-divide_by='+str(count/1000000),
                    '-output='+str(output), str(source)], check=True)
    scaled.append(str(output))
    manifest.append({'task': task, 'raw_samples': count, 'target_samples': 1000000})
subprocess.run([go, 'tool', 'pprof', '-proto', '-output='+str(base/'balanced.pgo')]+scaled, check=True)
(base/'profile-weights.json').write_text(json.dumps(manifest, indent=2)+'\n')
```

Normalization gives equal task weight, not a claim about real production traffic.
Record the source/binary/profile hashes, actual task argv/environment, profiler
settings and raw weights. Keep raw profiles so the weighting can be reviewed.
Round2 profiles from `da20b44f` are historical context, not fresh training data
for the new compiler shape.

## What the flags actually change

- `-l=4` changes the non-leaf/interface-call cost from 57 to 1 in the Go inliner.
  It does not remove safety checks or force every function inline. The ordinary
  body/callsite budget is 80, and very large callers still impose a cost-20 cap.
  This can expose wrappers, closure bodies and temporary objects to escape
  analysis while preserving their semantics.
- Standard PGO already enables interface and function-value devirtualization,
  considers the hot edges covering 99% cumulative weight, and raises hot-body
  budget to 2,000. Profile-guided interface specialization retains a guarded
  fallback. It is not a guarantee that previously escaping parameters can move
  to the stack.
- Static devirtualization explicitly rejects shaped generic receivers because
  it cannot supply their dictionaries. Inlining/PGO can help some surrounding
  calls, but compiler flags do not promise to erase every generic dictionary
  dispatch. Concrete library paths remain useful if the screen leaves those hot.
- `-d=inlbudgetslack=N` only expands candidate discovery for `GOEXPERIMENT=newinliner`.
  That experiment is off by default in this toolchain; the flag alone does not
  raise final callsite thresholds. It is a poor first blind tuning choice.
- `-d=pgoinlinebudget=4000` is supported, but test it only if compiler diagnostics
  show important hot callsites rejected at the current 2,000 budget. It still
  does not bypass the big-caller restriction. Avoid a large speculative budget
  sweep before seeing those rejections.

Source references in `/workspace/toolchains/go1.27.2/src`:
`cmd/compile/internal/inline/inl.go`,
`cmd/compile/internal/inline/inlheur/scoring.go`,
`cmd/compile/internal/base/flag.go`,
`cmd/compile/internal/devirtualize/{devirtualize,pgo}.go`,
`cmd/internal/pgo/pprof.go`,
`cmd/vendor/github.com/google/pprof/internal/report/report.go`.

## Decision and validation

Byte-gate all five workloads before interpreting any speed result. Interleave
the four variants with the current experiment runner, then independently confirm
only promising settings against the same warmed Java target. Report executable
bytes, ELF text size, build time/peak RSS, application elapsed/CPU/peak RSS,
allocated bytes and live heap from separate diagnostic runs. Preserve Go/Java
ratios and their memory-window caveat.

Large inlining can increase compiler memory use, machine code, instruction-cache
pressure and per-goroutine stack frames even while reducing allocations. These
are performance tradeoffs, not permission to relax output or Java comparisons.
Use the same `-m=2 -d=pgodebug=2` package flags in a separate diagnostic build if
needed to explain hot-site decisions. Preserve ordinary nil/bounds/write-barrier
checks, IEEE behavior and race-test coverage; do not use unsafe optimization flags.

## Followup: new-inliner activation

The inspected Go1.27.2 source has no `-d=newinliner` option. Its enabled debug
keys are the fields of `base.DebugFlags`; unknown keys fail. Full activation is
`buildcfg.Experiment.NewInliner`, selected by `GOEXPERIMENT=newinliner`.

This experiment changes exported inline metadata: noder/linker writes an extra
function-properties string and noder/reader reads it only when the experiment is
enabled. Consequently, enabling it only for selected compiler subprocesses while
reusing normally compiled dependencies is not a safe format-compatible shortcut.
The go command incorporates GOEXPERIMENT in package cache identities, giving a
matching rebuild of dependencies, including the standard library. Reuse that
cache for subsequent new-inliner candidates; do not manually splice old archives.

For the next reservation, the strongest clean alternative to `-l=4` is a new
inliner build with its default scores, the same source and the same PGO profile:

```bash
GOEXPERIMENT=newinliner GOMAXPROCS=2 GOGC=30 GOMEMLIMIT=3GiB go build -buildvcs=false -p 1 -pgo=/workspace/jsgo-toolchain-screen/cmd/jsgo/default.pgo -o /workspace/jsgo-results/toolchain-screen/bin/newinliner-pgo ./cmd/jsgo
```

Its score adjustments particularly favor calls exposing concrete interface
receivers, inlinable function arguments and branch constants. These are relevant
to this compiler's callback, collection and token-helper layers. Test default
scores first; combining a new inliner, `-l=4`, arbitrary budgets and an application
GOGC change would conceal which change caused the result. Generic dictionaries
and big-caller limits still constrain the potential benefit.

There is a package-local **diagnostic** heuristic path:
`-d=dumpinlcallsitescores=1` or `-d=dumpinlfuncprops=...` makes
`inlheur.UnitTesting()` true, which activates local heuristics. This does not
activate the experiment's exported/imported function-properties metadata and is
not equivalent to a full new-inliner build. Do not label that shortcut as the
new-inliner candidate. `inlbudgetslack`/`inlscoreadj` alone do not enable it.
