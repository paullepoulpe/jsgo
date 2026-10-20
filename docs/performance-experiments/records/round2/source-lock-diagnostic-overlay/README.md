# Source-lock activation diagnostic overlay

This is diagnostic-only instrumentation for the frozen
`/workspace/jsgo-structural-next` source. It changes no checked-out file and must
not be included in elapsed-time/RSS benchmark binaries. `manifest.json` records
the exact source commit, source/replacement hashes and asserted hook counts.

Run `python3 prepare.py` to regenerate the overlay and gofmt its replacement
files. `--repo PATH` selects another clean frozen worktree. `--merge-overlay PATH`
combines a profile/main overlay only when it targets that same worktree and
touches disjoint files; collisions are rejected rather than silently overwritten.

Build with Go's `-overlay=/workspace/jsgo-results/round2/source-lock-diagnostic-overlay/overlay.json`
using the normal actual import path, compiler flags and balanced PGO source.
The only added runtime variable is
`JSGO_SOURCE_LOCK_DIAGNOSTIC_JSON=/absolute/path/report.json`.
Use otherwise identical real task CLI arguments with the source-lock experiment
off and on. Reports write directly to that file, never stdout/stderr or compiler
reporters. No option, tracer mode or eligibility fence is changed.

The JSON records each ownership attempt, exact guard failure reason/details,
loader-type distribution, already-loaded files and wrapped count. It separately
counts actual owned cold loads and original monitor cold loads, including success
and source UTF-8 bytes, owned atomic-cache hits, reader/streaming paths, and Set/
Clear operations. Attempts increase immediately before original loader calls;
successes increase only after successful original bookkeeping. These byte counts
are UTF-8 storage length, not Java string/code-unit metrics. No source contents or
SourceFile pointers are retained by the diagnostic.

The report flush is registered on DoRun entry; the later owner-restoration defer
runs first, after all compiler jobs have joined. It therefore captures ordinary
success, early returns and internal failure unwinding. Counter instrumentation
uses atomic counters; ownership metadata is recorded serially under a private
diagnostic mutex. The normal timed source has no diagnostic fields/globals/hooks.

Source audit before capture found that SourceFile constructors retain text in a
Preloaded loader but leave `SourceFile.code` nil. Builtin extern creation therefore
does not itself initialize bookkeeping. Ownership installs after stock source
creation and before Init. The diagnostic is needed to determine whether real task
observer/loader eligibility rejects the entire batch, or whether loads are already
completed by an earlier path not evident from the ordinary source lifecycle.

Preparation validation only has run: unique hook matching, clean-source check,
source hashes, Python parsing and gofmt. Root owns build/capture slots.
