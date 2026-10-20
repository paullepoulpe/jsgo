# Exploratory Go factor screen

This tool screens Go configurations before expensive live Java confirmation.
It reports **Go variant / Go control** elapsed and peak-RSS ratios, paired by
round. A value below one means the variant takes less time or memory. It reports
no Java speed or memory comparison. Positive candidates require a subsequent
live paired-JVM run; a screen result is not a claimed victory over Java.

Provide the same variant JSON format as `paired_factors.py`: objects containing
`name`, absolute `binary`, string `go_env` experiment flags, and optional expected
`sha256`. Stock main, factor off, each of seven audited factors, and compound can
share one immutable task. Variant configuration cannot replace argv, inputs,
externs, chunks, output requirements, locale or CPU/GC settings.

```sh
python3 tools/bench/go_factor_screen.py --spec /abs/task.json \
  --variants /abs/ten-variants.json --control factor-off \
  --oracle /abs/completed-stock-Java/stock-java-oracle-digests.json \
  --image 'image@sha256:exact-digest' --out /abs/new-screen --rounds 3 \
  --provenance /abs/frozen-build-manifest.json
```

The oracle is explicit. Its accompanying `metadata.json` (or explicit
`--oracle-metadata`) must identify the pinned stock Java jar, the same complete
argv, ordered inputs and hashes, builtin extern identities, output artifacts and
required absences. Every stdout, stderr, zero exit, artifact and absence is gated
on every row. A digest alone with no matching work provenance is insufficient.
Downloads, images, builds and fixture generation must already be complete.

Every row is a fresh process in a new kernel-enforced 4 GiB, zero-swap, two-CPU,
affinity 0,1 container using the frozen helper API. Runtime defaults remain
unchanged. No JVM is launched by the screen. There are no compiler warmup rows;
all binaries are hashed before round one and source inputs are frozen. Seeded
random Latin rotations balance order: full V-round blocks visit every variant
in every position, and partial-block position counts differ by at most one.
Round pairing limits drift but does not create independent JVM cohorts.

The output contains the complete row plan, raw elapsed/CPU/RSS/cgroup metrics,
every strict gate, median/range/sample variance, and per-round control ratios.
Failures, negative results and partial evidence are retained. OOM is failure;
an outer OOM invalidates the screen. No confidence interval or Java performance
ratio is inferred from these exploratory repetitions.

Helper source hashes are pinned and executed from a new read-only snapshot;
imports compile exact checked bytes without stale `.pyc` caches. Binary hashes,
full input manifest, stock-oracle provenance and source snapshots are retained.
Use repeatable `--provenance` to copy/hash actual build commands, source commits
and build manifests. A flag value alone is not proof that a factor dispatched;
retain independent actual CLI activation evidence before attribution.
