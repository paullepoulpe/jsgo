# Enabled parser profile: output and remaining work

> Attribution correction: retained stack-family totals are lower bounds, not hard performance ceilings. CPU stacks are limited to64 PCs; owning pass frames can be lost during generic recursion. See `profile-stack-depth-audit.md`.

Source: frozen 11534351, all six experiment flags enabled. These are the same byte-gated non-tracing profiles described in `parser-profile-evidence.md`; CPU is sampled CPU, not wall elapsed, and allocation is sampled `alloc_space`, not peak RSS. Inclusive categories overlap. A single short whitespace profile has coarse 10ms samples.

| Task | CPU samples total | Emitter + writer/output union | Sampled allocation total | Emitter + writer/output |
|---|---:|---:|---:|---:|
| typed-app | 6.87s | 0.08s (1.2%) | 1146.5MiB | 2.0MiB |
| bundle | 26.59s | 0.34s (1.3%) | 4713.4MiB | 23.5MiB |
| downlevel-whitespace | 1.91s | 0.11s (5.8%) | 332.2MiB | 2.0MiB |
| downlevel-simple | 5.25s | 0.08s (1.5%) | 1044.5MiB | no positive samples |
| checks-only | 10.20s | 0.01s (0.1%) | 1430.5MiB | no positive samples |

The emitter cannot credibly explain a ~35% whole-program gap on these inputs. `StringBuilder.Append` already copies ASCII directly into its UTF-16 buffer without a temporary conversion. `MappedCodePrinter.Append` already counts UTF-16 columns from the builder-length delta. Pretty-printer numeric-source extraction already caches source conversion and has an ASCII byte-substring path. Removing these protections risks source-map columns, unpaired UTF-16 behavior, or output bytes for a small ceiling.

One actual double conversion remains in the output adapter: `BufferedWriter.Write` converts to UTF-16, `FlushBuffer` converts the buffer to WTF-8, then `charsetOutputStreamWriter.Write` converts that back to UTF-16. A private UTF-16 bulk-writer capability could avoid the middle conversion only for audited concrete charset writers, retaining the public string-writer call schedule for custom writers. Surrogate pairing across buffer boundaries, pending-high-surrogate behavior, exact charset replacement, close/flush and partial-write failure ordering must remain unchanged. This is a small bounded target, not the primary next redesign.

## More material allocation findings

Whitespace has 81.0MiB sampled allocation and 0.25s CPU under linked-map/set operations (24.4% allocation / 13.1% CPU, including their callers). Scope scanning accounts for 41.3MiB / 0.15s. Normalization accounts for 48.1MiB / 0.16s, including contextual rename inversion 33.5MiB / 0.06s. Module resolution and name conversion account for 44.0MiB / 0.08s. These buckets overlap and cannot be added.

For bundle, reference collection is 924MiB / 3.83s, scope scanning 517MiB / 2.36s, and contextual rename inversion 326.7MiB / 0.78s. The latter includes the existing duplicate-conflict-pair suppression and is not an untouched quadratic algorithm.

Two concrete stateless module-path opportunities merit a small independent screen:

1. `ModuleNames_CanonicalizePath` always splits, allocates a second segment array, and rejoins even when the result equals the input. A segment scan can return the original immutable string when no segment is `.` or `..` and no trailing slash needs Java split-limit-zero normalization. Interior empty segments must remain; `"/"`, `"///"`, empty input, leading `..`, absolute `..`, and repeated separators need exact oracle cases. This removes work rather than retaining a result cache.
2. `ModuleNames_ToJSIdentifier` performs nine sequential `ReplaceAll` calls, repeatedly allocating strings for long absolute module paths. A single byte-oriented pass can apply the same ASCII substitutions and `%20` sequence rule after the existing suffix/prefix stripping. Preserve arbitrary UTF-8/WTF-8 bytes, `.js.i.js` stripping, one-leading-slash behavior in `ToModuleName`, and exact replacement ordering. All substitutions introduce no new percent sequence, so a fused pass has a straightforward equivalence argument.

`NodeModuleResolver.ResolveJsModuleFile` also computes a canonical candidate and then `Locate` canonicalizes it again when package-json remapping is absent. Avoiding that duplicate needs a concrete-owned, stock-path-escaper guard: the public resolver has custom escaper callbacks whose invocation count can be observable. This is a less attractive first patch than the pure helper fast paths.

Raw category data: `parser-emission-profile-analysis.json` and `parser-emission-pass-analysis.json`. No new compiler build or timed benchmark was run for this audit.
