# Host memory-pressure investigation

The sole v4 accumulator/bundle confirmation stopped in cohort 2 during candidate
request 68 because the measured parent cgroup recorded 387 additional `max`
events and 14,879 pages of direct scan/reclaim. This is a failed, incomplete
confirmation, not an accepted five-cohort result. Its original records and
independent failure review remain unchanged.

After all JVMs exited, the 8 GiB parent still held roughly 5.64 GB: about 2.38 GB
of file cache and 3.18 GB of kernel memory, almost all reported as reclaimable
slab, with only about 80 MB of anonymous memory. These post-run observations
show limited headroom; they do not prove the exact object or file responsible
for every pressure event. No new OOM event occurred during the failed run.

Root asked the independent reviewer to withhold partial performance ratios while
checking whether a clean restart was feasible. Three bounded, read-only
`POSIX_FADV_DONTNEED` attempts requested eviction of clean files from owned task
directories. Open descriptors used `O_RDONLY | O_NOFOLLOW`; each file's device,
inode, size and modification time were checked before and after the call. No
files were deleted or truncated. JVM flags, GC settings, memory caps and benchmark
rules were unchanged. No global `drop_caches`, synthetic memory pressure or
privileged cgroup write was used.

| Attempt | Scope | Files advised | Errors | Parent bytes before / immediately after |
| --- | --- | ---: | ---: | --- |
| 1 | Five task/cache roots, files at least 1 MiB | 2,640 | 0 | 5,627,949,056 / 4,951,568,384 |
| 2 | The same roots, all regular files | 104,714 | 0 | 4,951,216,128 / 4,752,760,832 |
| 3 | Task worktrees, toolchains and task cache, all regular files | 1,176,889 | 0 | 4,729,708,544 / 4,952,178,688 |

A later snapshot after attempt 3 was 4,696,010,752 bytes. Traversing many files can
itself populate filesystem metadata caches; a successful advice call does not
guarantee that the kernel evicted pages. None of the observed final totals met
the proposed initial-headroom target of at most 4 GiB parent usage before starting
the two JVMs. That target was a feasibility check for a possible new run, not a
retroactive admission criterion applied to old results.

The restart decision was closed before root saw any partial performance ratio:
**no restart in this batch**. The original failed run remains incomplete, its
completed cohort is not promoted to a smaller confirmation, and no fresh cohort
is substituted. A future full confirmation needs adequate host headroom in
addition to the unchanged per-JVM 4 GiB cap.

The three original JSON receipts, full file inventories and the executed script
versions are included in the evidence archive. They are environment-preparation
records, not compiler benchmarks. This report is the compact public view; the
large per-file inventories are not duplicated as standalone report files.
