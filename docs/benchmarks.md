# Benchmarks: Go port vs Java Closure Compiler

Per-milestone performance checks (PLAN.md §6). Every comparison runs the
same layer of both implementations on the same inputs. Both outputs are
first verified identical by the oracles.

## M1 — parser (`internal/parser` vs `parsing.parser.Parser.parseProgram`)

Harnesses:

- `tools/bench/java/ParseBench.java`;
- `tools/bench/parsebench` (Go);
- `tools/bench/measure.py` (wall time, CPU time, peak RSS via `wait4`).

**Setup:**

- **Inputs:** the full pinned corpus: 14,690 files, 84.8M UTF-16 chars
  (`tools/oracle/corpus.txt`). All files are read into memory before
  timing.
- **Runs:** 8 parses of the whole corpus per process, single-threaded, with
  the compiler's ES_NEXT parser configuration.
- **Settings:** default JVM (JDK 21, G1) and Go 1.24 (GOGC=100,
  GOMAXPROCS=4) settings.
- **Machine:** 4 vCPUs, *shared with several porting agents*. Absolute
  numbers are therefore noisy; ratios are indicative. Re-run on a quiet
  machine for the M9 report.
- **Sanity check:** both report the same 2,776 parse errors.

| Per full-corpus parse | Go (first port) | Go (after `jutil` fixes, 9dfa23e) | Java, 1st iteration (cold JIT) | Java, steady state (iters 7–8) |
|---|---:|---:|---:|---:|
| Wall time | 8.2–9.7 s | **6.6–6.8 s** | 5.9 s | **2.2 s** |
| Allocated | 3.11 GB | 2.63 GB | 2.45 GB | 2.45 GB |
| GCs / GC pause | 32 / ~5 ms | 28 / ~5 ms | 11 / 141 ms | 2 / 20 ms |
| Whole process: wall / peak RSS | 69 s / 540 MB (8 iterations) | 20.4–20.7 s / 530 MB (3 iterations) | 38 s / 2.25 GB (8 iterations) | (same process) |

**Findings:**

- **The first port was 1.4× slower than cold Java and ~4× slower than warm
  Java.** A CPU profile showed ~37% of the time in allocation.
  - The single biggest avoidable cost (8%) was `jutil.CheckArgumentMsg`
    boxing its varargs on every call, even though the check passes.
  - Two more costs were UTF-16 ↔ string conversions that allocated twice,
    and an `ArrayList.Get` that could not be inlined.
- **The fixes are now conventions (PORTING.md §5):**
  - fixed-arity Preconditions, as in Guava;
  - single-allocation string conversion;
  - out-of-line failure paths.

  Result: −22% time and −15% allocation.
- **GC policy is not the gap.** `GOGC=400` gains ~5%. `GOGC=off` is much
  slower (9.8 GB RSS, page-fault bound).
- **The remaining gap to a warm JVM (~3×) is structural.** Go's allocator
  is slower than HotSpot's thread-local bump allocation, and Go has no
  escape analysis across interface calls and inlines less. A line-for-line
  port allocates what Java allocates. Options, after parity, behind the
  same zero-diff oracles:
  - slab-allocating per-token and per-position objects;
  - reducing interface dispatch in the scanner's hot path;
  - parsing independent input files concurrently. The compiler already
    treats inputs independently at this stage, and the output stays
    deterministic.
- **Go already wins on process-level costs:**
  - no JIT warm-up: Go's first iteration is as fast as its last;
  - no JVM startup;
  - ~4× less memory.

  For single-shot CLI invocations, the common case, the comparison that
  matters is against cold Java.

To reproduce:

```sh
go build -o .cache/bench-build/parsebench ./tools/bench/parsebench
javac -cp .cache/closure-compiler.jar -d .cache/bench-build tools/bench/java/ParseBench.java
python3 -I tools/bench/measure.py go.rusage .cache/bench-build/parsebench 8 .cache/corpus @.cache/corpus-files.txt go-perfile.tsv
python3 -I tools/bench/measure.py java.rusage java -Xss1g -cp .cache/closure-compiler.jar:.cache/bench-build ParseBench 8 .cache/corpus @.cache/corpus-files.txt java-perfile.tsv
PARSEBENCH_CPUPROFILE=parse.pprof .cache/bench-build/parsebench 2 .cache/corpus @.cache/corpus-files.txt /dev/null
```

### Known hot spots to revisit after parity

- `rhino.TokenStream_IsKeyword` converts every NAME to UTF-16. That is
  ~5% of `ParserRunner.parse` allocations (found while porting IRFactory).
