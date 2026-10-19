from pathlib import Path
import json,math,statistics,hashlib,csv,io,datetime
O=Path('/workspace/jsgo-publish-records/docs/performance-experiments');R=Path('/workspace/jsgo-results/round2')
D=json.loads((O/'datasets.json').read_text()); manifest=json.loads((O/'measurement-manifest.json').read_text());sources=json.loads((O/'source-manifest.json').read_text())
# Preserve a browseable metric ledger wherever the saved schema supplies it.
metrics=[]
for item in D:
 obj=json.loads((O/item['browse']).read_text())
 for task,row in obj.items() if isinstance(obj,dict) else []:
  if not isinstance(row,dict):continue
  if isinstance(row.get('variants'),dict):
   for name,v in row['variants'].items():
    if not isinstance(v,dict) or 'median_s' not in v:continue
    metrics.append({'dataset':item['path'],'task':task,'variant':name,'go_s':v.get('median_s'),'go_peak_rss_mib':v.get('median_peak_rss_mb'),'java_warm_s':row.get('java_warm_s'),'java_peak_rss_mib':row.get('java_peak_rss_mb'),'java_over_go_time':v.get('speed_vs_java_warm'),'java_over_go_rss':v.get('memory_ratio_java_over_go'),'all_outputs_match':v.get('all_outputs_match')})
  else:
   for name,v in row.items():
    if isinstance(v,dict) and 'go_median_s' in v:
     metrics.append({'dataset':item['path'],'task':task,'variant':'threads='+str(name),'go_s':v.get('go_median_s'),'go_peak_rss_mib':v.get('go_peak_rss_mb'),'java_warm_s':v.get('java_warm_median_s'),'java_peak_rss_mib':v.get('java_peak_rss_mb'),'java_over_go_time':v.get('time_ratio_java_over_go'),'java_over_go_rss':v.get('memory_ratio_java_over_go'),'all_outputs_match':v.get('all_outputs_match')})
(O/'metrics.json').write_text(json.dumps(metrics,indent=2)+'\n')
with (O/'metrics.csv').open('w') as f:
 w=csv.DictWriter(f,fieldnames=list(metrics[0]),lineterminator='\n');w.writeheader();w.writerows(metrics)
S=json.loads((R/'inline-bitset-template-screen/summary.json').read_text()); rows=[];reductions=[]
for task,row in S.items():
 old=json.loads((R/'baseline'/task/'jsgo.json').read_text())['runs']; ot=statistics.median(v['wall_s'] for v in old);om=statistics.median(v['rss_mb'] for v in old);n=row['variants']['inline'];delta=(1-n['median_s']/ot)*100;reductions.append(n['median_s']/ot)
 rows.append(f"| {task} | {ot:.3f} → {n['median_s']:.3f} | {delta:.1f}% | {om:.1f} → {n['median_peak_rss_mb']:.1f} | {n['speed_vs_java_warm']:.3f} | {n['memory_ratio_java_over_go']:.2f} |")
geo=(1-math.exp(sum(math.log(x) for x in reductions)/len(reductions)))*100
s='''# Performance experiment record

This publishes the saved experiments through the completed template-guard V2 screen (2026-10-19). It includes positive, negative, rejected, superseded and inconclusive outcomes, validation failures, source-only audits, measurements and prototype source. Continued exploration after this snapshot is not represented as completed evidence.

**Warm Java has not been beaten across the workloads. The 3.7% figure was the latest isolated typed-app improvement, not the overall improvement. No exhaustive search or proof that further gains are impossible is claimed.** This PR publishes evidence and reproduction material; prototype compiler changes in archives are not applied to the PR's compiler/library tree.

## How to read the numbers

The latest same-binary template experiment measured typed-app 5.123 → 4.934 s (3.7% less elapsed time), RSS 275.7 → 252.6 MiB. Checks-only measured 7.504 → 7.405 s (1.3% less time), RSS 653.4 → 658.8 MiB. These five-sample ranges overlap. [Detailed round 11](records/round2/template-guard-v2-round11.md).

For a broader comparison, the following table uses ONE complete five-task candidate, a6afc95f (`inline` in round 10), against the recorded current-main baseline da20b44f (compiler tree matching main 123b46e8). Every row uses that same candidate and configuration; this does not combine different per-task winning binaries.

| Workload | Main → candidate Go seconds | Less elapsed time | Main → candidate peak MiB | Java/Go time | Java/Go RSS |
|---|---:|---:|---:|---:|---:|
'''+ '\n'.join(rows)+f'''

The equal-weight geometric mean of these five latency ratios corresponds to **{geo:.1f}% less elapsed time**. That arbitrary equal-weight aggregate is a descriptive screening statistic, not a confirmed overall improvement, workload-weighted throughput claim, or improvement introduced by this documentation PR. The baseline and candidate are separate cohorts. Candidate source includes historical experimental substrate choices, including parser-owned node batches and vendored Unicode tables; publication does not reinstate rejected production choices. Later c7998d45 was remeasured only on two typed tasks, so there is no complete five-task result for that latest binary. Earlier work already integrated into main is also outside this main-to-prototype comparison.

Both comparison cohorts use Go 1.27.2, default runtime GC, the same five inputs, shipped PGO and the recorded build options/CPU setup; candidate experiment flags and source differ and are retained in metadata. Three Go observations per cell are a screen. Java is the older stored warm JVM reference (12 warmups/10 measured compiles), with fresh compiler/extern parsing each iteration. Ratios above1 favor Go. Go RSS is median fresh CLI lifetime peak; Java RSS is whole warm JVM lifetime peak including warmup. These memory windows differ. Never treat the ratio as equivalent per-compilation live heap. [Baseline](records/round2/baseline-summary.md) · [complete round10](records/round2/inline-bitset-template-round10.md).

## Archive contents

- **{len(D)} saved summary files**: [dataset index](datasets.json) and the list below. Some are validation or protocol summaries, not performance comparisons.
- **{len(metrics)} normalized task/variant metric rows** where the schema supplies them: [CSV](metrics.csv), [JSON](metrics.json). Missing historical RSS/ratio fields remain blank/null; they are not invented. Nonstandard schemas remain available in full in each original summary.
- **{len(manifest['files'])} measurement/support files**: [measurement-records.tar.gz](measurement-records.tar.gz), [per-file SHA-256 manifest](measurement-manifest.json). Includes raw per-run timing/RSS/digests, machine/configuration/provenance, scripts, build/validation/failure logs, CPU/heap profiles, reports and diagnostic overlays.
- **{len(sources['worktrees'])} source worktree snapshots**: [prototype-sources.tar.gz](prototype-sources.tar.gz), [source manifest](source-manifest.json). Each tracked patch applies to main123b46e8; listed untracked source files are preserved too. Some snapshots are unfinished/unmeasured source, explicitly not endorsements.
- **Exact committed history**: [prototype-history.bundle](prototype-history.bundle) preserves experimental commits excluded from main, including intermediate measured implementations. It is a Git data file, not an applied merge. It complements working-tree patches so SHA references remain recoverable after the cloud workspace expires.
- **71 distinct design-note versions**: [index and source attribution](design-notes.json), browsable under `design-notes/`.
- **Archive hashes/sizes**: [archives.json](archives.json). [Verification utility](verify_archive.py).

Downloaded corpora, repeated generated JS/maps, compiled CLI/test executables/classes, symlink targets and PID files are excluded. Their pinned source references and measured artifact hashes remain in metadata; builds reproduce executables. There are no credentials in the archives. Historical `state.md`/`progress.md` and old design notes are retained verbatim as contemporaneous records, not current instructions or validated estimates. The current policies and conclusions in this index take precedence over stale expectations in those records.

## Current conclusions and remaining opportunities

- Runtime GOGC tuning is no longer an exploration target. Older GC datasets are retained only as historical evidence. New code comparisons use default runtime GC/CPU settings.
- Persistent workers and parsed-extern reuse have been dropped from the active target. Earlier unfair cached-Go/fresh-extern-Java comparisons are not valid optimization evidence; matching-Java worker followups and that correction remain archived. The target is fresh non-worker CLI compilation.
- Two-inline-word BitSets are mixed: typed/simple screens improve, bundle regresses; no default promotion. Fresh-template V1 regressed badly because per-scope admission introduced expensive monitor stack walks. V2 removes that regression and shows modest net gains against its own control. Other pass/storage/scope/parallel factors and their compound screens remain recorded, including combinations that fail to add their isolated gains.
- Compatibility tests gate the exercised paths, not all possible Java/Go programs or external API behavior. Historical broad suites and newer focused suites have distinct scopes. Skips and failed fixture/build attempts remain visible. No full-port parity claim is inferred from a focused suite.
- **Further gains remain plausible.** A consistent five-task candidate still needs roughly10–25% less wall time, depending on task, to equal its stored warm-Java reference. That does not prove the remaining gap can be eliminated. Comfortable wins need more than matching Java.
- Next evidence-driven experiments: exact BitSet capacity/Clone occupancy before a one-word layout; local parser reuse of existing global string cells with actual hit/contention counts; larger pass-specific native data structures/algorithms; safe additional per-script parallelism; semantic-type construction/storage reuse with explicit aliasing/mutation/resolution proofs. Existing native CFG/array solver state and shared string headers must not be counted again as new ideas.
- Ambitious options remain research, not forecast speedups: typed/indexed internal representations, denser pass-local work/state, ownership-aware builders and broad traversal redesign. Any must preserve compiler artifacts, diagnostic/order/failure semantics and supported public/custom fallbacks, and must be tested against fresh matched controls. No fixed theoretical maximum is supported by the profiles.
- More cores may help some phases, but this machine has a two-CPU quota. Four compiler jobs are oversubscribed. Compiler job counts bound work items and do not globally cap goroutines or change GOMAXPROCS. [Fresh default-GC whitespace thread sweep](records/round2/default-gc-white-traversal-thread-sweep/summary.json). Full performance matrices for other configurations are not invented from compatibility gates.

Useful current audits: [type/flow design](records/round2/type-flow-design-audit.md), [node/string layout](records/round2/node-string-cell-layout-audit.md), [parser front cache](records/round2/parser-string-front-cache-audit.md), [one-word BitSet](records/round2/oneword-bitset-alternative-audit.md), [template guard](records/round2/fresh-template-guard-audit.md).

## Reproduce and inspect

Run `python3 docs/performance-experiments/verify_archive.py` from a normal full clone to verify archive hashes and the measurement manifest without executing archived code. Extract the tarballs into a separate directory. Their paths preserve the original record/worktree relationships; original absolute paths are historical and must be rebased for local reproduction.

Fetch the repository's complete main history before importing the bundle (a shallow clone may lack prerequisites). `git bundle verify docs/performance-experiments/prototype-history.bundle` checks prerequisites. Fetch the specific desired archived branch into a new local experiment branch, or use the tracked patch plus listed untracked files against the exact source-manifest base. Do not overwrite an active checkout with arbitrary archived patches. Build/test commands, runtime flags, toolchains, pinned Java jar and input references are in each dataset's machine/command/provenance files. Binary identities are retained even though binaries are not.

## Dataset inventory

Families marked historical-worker or historical-GC are excluded from the current optimization claim. Validation summaries contain compatibility results only; microbenchmarks are not whole-compiler speed claims. Consult the accompanying reports/raw metadata before comparing cohorts.

| Saved dataset | Family / qualification |
|---|---|
'''
for d in D:
 path=d['path'];name=path.lower()
 if 'worker' in name or '/warm-' in '/'+name:tag='Historical worker / protocol; superseded target, inspect Java matching'
 elif any(k in name for k in ['memory-envelope','parser-gc']):tag='Historical GC settings; exploration suspended'
 elif any(k in name for k in ['gates','parallel-cli-owned','parallel-compound-final','parallel-expanded','parallel-fsiv-cli','parallel-modules-checked','parallel-zip','type-flow-cli-next']):tag='Compatibility / validation; not a timed speedup claim'
 elif any(k in name for k in ['micro','parser-tree-storage']):tag='Diagnostic / microbenchmark; not whole-compiler throughput'
 elif 'thread' in name:tag='Thread scaling; inspect environment and Java cohort'
 else:tag='CLI candidate/toolchain screen; inspect per-variant controls and sample count'
 s+=f"| [{path}]({d['browse']}) | {tag} |\n"
s+='\n## Other saved baseline/report entry points\n\n'
for p in sorted((O/'records').rglob('*.md')):
 rel=p.relative_to(O)
 if p.name in {'report.md','progress.md','state.md'}:s+=f'- [{p.relative_to(O/"records")}]({rel}) — historical snapshot; see current policies above.\n'
(O/'README.md').write_text(s)
# Include the exact archive builder and index generator, whose source is publication tooling.
shutil=None
for name in ['publish_records.py','write_index.py']:
 (O/name).write_bytes((Path('/tmp/jsgo-publication')/name).read_bytes())
archives=json.loads((O/'archives.json').read_text());p=O/'prototype-history.bundle';archives.append({'path':p.name,'bytes':p.stat().st_size,'sha256':hashlib.sha256(p.read_bytes()).hexdigest()});(O/'archives.json').write_text(json.dumps(archives,indent=2)+'\n')
print(json.dumps({'metric_rows':len(metrics),'geomean_less_time_percent':geo,'index_bytes':len(s)},indent=2))
