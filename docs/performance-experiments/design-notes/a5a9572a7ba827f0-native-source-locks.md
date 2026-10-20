# Native source-cache locks

The worker profile showed goroutine identity discovery (`runtime.Stack` through
`jutil.Monitor`) taking substantial CPU while loading, clearing and rebuilding
source files. This candidate changes only locks whose callback graph is known
not to require Java monitor reentrancy. It does not change jutil.Monitor itself.

## SourceFile boundary

Exact concrete Preloaded, StubSourceFileCodeLoader, OnDisk and AtZip loaders use
a plain Go mutex. Every other loader, including a type embedding and overriding
a native loader, retains a reentrant jutil.Monitor. Selection happens once in
the SourceFile constructor, before publication. The loader field is not replaced
later. Known native source files use a smaller lock header; custom sources allocate
their monitor separately.

The source cache already publishes immutable text/count/offset snapshots through
an atomic pointer. Every writer still takes the selected per-file lock, so reads
and concurrent invalidation retain the same generation guarantees. Exception
paths retain deferred unlocks.

The no-reentry argument for native loaders is specific:

- GetCode calls getCodeStateLocked, which invokes the loader and the private
  setCodeAndDoBookkeepingLocked helper without acquiring the lock again.
- Offset/line paths invoke getCodeStateLocked/findLineOffsetsLocked while holding
  the same one lock; these helpers do not acquire it.
- GetCodeReader invokes OpenUncachedReader while locked and calls GetCode only
  after releasing the lock.
- SetCodeAndDoBookkeeping and restoreCachedCounts publish private state directly.
- Preloaded only returns its string; Stub throws; OnDisk calls os file I/O and
  charset decoding; AtZip opens a private ZIP-cache reader. None calls SourceFile.

Custom loader callbacks can still reenter SetCodeDeprecated, ClearCachedSource,
line-cache construction, GetCode and GetCodeReader. A focused regression exercises
that behavior for both loading and reader-opening overrides.

## ZIP cache boundary

The private global cache has a fixed loader and fixed eviction listener. Its
call graph is global lock → CachedZipFile lock → refreshIfNeededLocked →
maybeCloseLocked or archive opening. No path acquires the global cache again.
The per-archive methods already use Locked helpers and do not reacquire their
own lock. Both can therefore use sync.Mutex while retaining the current locking
scope, entry-opening atomicity, refresh semantics and exception handling.

The returned reader remains outside these locks, as before. This change makes no
new guarantee about replacing or closing an archive while an already returned
stream is being read.

## Validation

New tests cover exact-loader selection, custom embedding-loader reentry while
loading/opening, native IOException recovery, reloading and old-snapshot validity.
Existing SourceFile concurrent invalidation/shared-position tests and ZIP-cache
concurrent read/refresh tests remain applicable. Central package/race tests and
paired CLI/worker output/time/RSS gates are required before retention; no timing
claim follows from identifying a profile hotspot alone.
