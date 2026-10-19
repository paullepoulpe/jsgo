# Private ASCII output buffer prototype

`JSGO_EXPERIMENT_ASCII_OUTPUT_BUFFER=1` selects byte storage for stock output
inside a fresh CLI source-ownership window. It currently requires
`JSGO_EXPERIMENT_CLI_SOURCE_LOCK=1`. Public builders/printers and custom factories,
license trackers, observers, options callbacks and compiler construction use the
original UTF-16 `jutil.StringBuilder` operations. The source-ownership window is
neither extended nor cached. Chunk-aware/custom license trackers conservatively
fall back in this first prototype.

Each per-script printer currently writes UTF-16, encodes a string, and hands that
string to the compiler join builder, which decodes/writes UTF-16 again before
encoding final output. The prototype replaces the two private buffers with ASCII
bytes until non-ASCII arrives. Lengths/positions then remain identical to UTF-16
without decoding. The only interior printer edit is insertion of a newline at
a preferred line-break position. Last-char checks, mapping offsets, suffix
comparisons, reset and line/column accounting retain the existing algorithms.

On the first non-ASCII string or char, existing bytes are copied into an original
UTF-16 builder and byte storage is released. The promotion is permanent, even
after reset. This matters for Unicode correctness: separately appending a WTF-8
high-surrogate string and low-surrogate string must combine into one code point
when the final builder is encoded. Keeping arbitrary UTF-8 strings and simply
concatenating their bytes would change that result. Promotion also preserves
inserting a newline between surrogate halves by their UTF-16 index and invalid
UTF-8 normalization. String snapshots copy the ASCII data so later edits cannot
mutate previously returned strings.

The baseline is the already recorded `structural-profiles/downlevel-whitespace`
capture (source `5b720ead`, diagnostic default-GC run), not a new timing run.
Its 1.39 CPU seconds contain 0.06 seconds (4.3%) under `Compiler.ToSource`, and
0.04 seconds (2.9%) under the printer. These are coarse 10-ms samples and an older
source checkpoint. Output cannot explain the whole remaining Java gap.

The corresponding heap512 profile has 420,373,275 sampled allocation bytes.
Output printer/join ancestry accounts for 29,779,395 bytes (7.08%), including:

| Allocation | MiB |
| --- | ---: |
| Per-script UTF-16 builder growth | 12.952 |
| Compiler join UTF-16 builder growth | 10.232 |
| Per-script UTF-16-to-string encoding | 1.667 |
| Joined UTF-16-to-string encoding | 2.547 |

ASCII bytes can roughly halve buffer growth storage, while immutable output
strings still need allocation. The private facade adds 24 bytes to each owning
printer/join header and a mode branch; original/default-mode operations still
allocate the same UTF-16 builders. No speedup, RSS improvement, allocation saving
or default-off neutrality has yet been measured. Unicode-heavy/license output
can promote early and receive less benefit.

Separate runner writer adapters account for another 9.126 MiB of
`StringToUTF16` and 4.031 MiB of `StringFromUTF16` under `WriteOutputOfCode` in
that profile. Those adapters are unchanged here. Their encoding, chunking,
flush/error behavior and custom callbacks need a separate audit.

Source-only tests compare buffer operations directly with the original builder,
including bounds exceptions before/after promotion, immutable snapshots,
surrogates, malformed UTF-8, NUL padding and repeated reset. Guarded real printer
tests record actual selected/printer/promoted counts through a nil-by-default
per-build probe, compare every mapping, and exercise pretty/compact output,
line cuts, Unicode promotion, joined source maps, suffix checks and reset.
Public/custom/copy-self fallback and actual stock WHITESPACE/ADVANCED CLI
lifecycles are included. Independent source review found no blocker. These
tests have not been built/run by the author; central validation and same-binary
off/on timing remain required.
