# Fresh builtin extern parser coverage

Both ordinary and owned native allocators pass every fresh original-parser
AST/JSDoc/source/property/metadata/ordered-diagnostic/exception oracle. Their
inventory JSON is identical. Resource and source/tool provenance is recorded in
`extern-inventory-provenance.json`.

| Environment | Accepted files | Accepted source bytes | Whole-file fallback bytes | Original fragments in accepted files |
| --- | ---: | ---: | ---: | ---: |
| BROWSER | 127/128 | 1,899,494 / 1,976,245 (96.12%) | 76,751 | 0 |
| CUSTOM | 4/5 | 218,053 / 294,804 (73.97%) | 76,751 | 0 |

The sole fallback is `es3.js`, at UTF-16 offset6807: `var arguments;`.
The native sink deliberately rejects that identifier. The original pipeline
handles the whole file unchanged. The browser archive selection includes its
empty directory entry, as the stock loader does.

This closes an evidence gap. It does not establish a performance improvement:
source-byte coverage is not proportional to CPU cost, and no AST is cached or
preparsed for a later compilation. Parsing test execution took0.586s with the
ordinary allocator and0.563s with the owned allocator; these are correctness run
durations, not controlled benchmark comparisons. GOGC/GOMEMLIMIT/GODEBUG were
unset, GOMAXPROCS=1 and build parallelism1.
