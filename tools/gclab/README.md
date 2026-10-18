# GC lab scripts (2026-10-18)

Reference material for the experiment in
[docs/gclab-2026-10-18.md](../../docs/gclab-2026-10-18.md); not part of the
build. The Python scripts drove the GC-setting sweep, the heap profiles and
the summaries, and carry the paths of the session that ran them. The opt-in
profiling hook (`hook-*.go.txt`) was compiled in only through
`go build -overlay`, together with one-line calls added to temporary copies
of `cmd/jsgo/main.go`, the tracer, `PhaseOptimizer` and `System` at the
points the report names (end of parsing, checks, type checking, and
mid-optimization). To re-run: restore the hook as `.go` files in an overlay,
add those calls, and follow the report's method section.
