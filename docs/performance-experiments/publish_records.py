import collections, datetime, gzip, hashlib, io, json, os, re, shutil, subprocess, tarfile
from pathlib import Path
ROOT=Path('/workspace/jsgo-results'); REPO=Path('/workspace/jsgo-publish-records'); OUT=REPO/'docs/performance-experiments'; OUT.mkdir(parents=True,exist_ok=True)
BASE=subprocess.check_output(['git','-C',str(REPO),'rev-parse','origin/main'],text=True).strip()
def sha(b):return hashlib.sha256(b).hexdigest()
def save(p,v):p.parent.mkdir(parents=True,exist_ok=True);p.write_text(json.dumps(v,indent=2)+'\n')
secret=re.compile(rb'(?:gh[pousr]_[A-Za-z0-9]{30,}|github_pat_[A-Za-z0-9_]{40,}|AKIA[0-9A-Z]{16}|-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----)')
def check(b,label):
 if secret.search(b):raise RuntimeError('Possible secret pattern in '+label+'; inspect without printing value')
def member(t,name,b):
 i=tarfile.TarInfo(name);i.size=len(b);i.mtime=0;i.mode=0o644;t.addfile(i,io.BytesIO(b))
keep_ext={'.json','.md','.log','.txt','.tsv','.csv','.py','.sh','.go','.java','.pprof','.pgo','.stdout','.stderr'}
records=[];excluded=collections.Counter();datasets=[]
with tarfile.open(OUT/'measurement-records.tar.gz','w:gz',compresslevel=6) as t:
 for p in sorted(ROOT.rglob('*')):
  if p.is_symlink():excluded['symlink']+=1;continue
  if not p.is_file():continue
  rel=p.relative_to(ROOT);s=p.stat();ext=p.suffix
  include=ext in keep_ext or (ext=='' and s.st_size<2000000 and p.name in {'stdout','stderr','exit','command','metadata','status','time','timing','rss'})
  if not include:excluded[ext or 'other/no-extension']+=1;continue
  b=p.read_bytes()
  if b.startswith(b'\x7fELF'):excluded['ELF binary']+=1;continue
  if ext not in {'.pprof','.pgo'}:check(b,str(p))
  member(t,'records/'+str(rel),b);records.append({'path':str(rel),'bytes':len(b),'sha256':sha(b)})
  browse=ext=='.md' or p.name in {'summary.json','machine.json'} or (ext=='.json' and len(rel.parts)<=2 and s.st_size<2000000)
  if browse:
   dest=OUT/'records'/rel;dest.parent.mkdir(parents=True,exist_ok=True);dest.write_bytes(b)
  if p.name=='summary.json':
   data=json.loads(b);datasets.append({'path':str(rel),'browse':'records/'+str(rel),'keys':list(data) if isinstance(data,dict) else [],'sha256':sha(b)})
member_manifest={'root':'jsgo-results','files':records,'excluded_counts':dict(excluded),'policy':'Includes saved metadata, raw timing/RSS JSON, scripts, validation/build logs, textual diagnostics, CPU/heap profiles and reports. Excludes compiled executables/test binaries/classes, downloaded benchmark corpora, generated JS/source maps, symlink targets and operational PID files. Artifact digests remain in run JSON. Credentials are not included.'}
save(OUT/'measurement-manifest.json',member_manifest);save(OUT/'datasets.json',datasets)
# Snapshot every existing experiment worktree, without integrating its code.
wtout=subprocess.check_output(['git','-C',str(REPO),'worktree','list','--porcelain'],text=True)
worktrees=[];notes={};sources=[]
with tarfile.open(OUT/'prototype-sources.tar.gz','w:gz',compresslevel=6) as t:
 for block in wtout.strip().split('\n\n'):
  fields=dict(line.split(' ',1) for line in block.splitlines() if ' ' in line)
  if 'worktree' not in fields:continue
  w=Path(fields['worktree'])
  if w==REPO or not w.name.startswith('jsgo'):continue
  name=w.name;head=fields.get('HEAD');branch=fields.get('branch','detached')
  b=subprocess.check_output(['git','-C',str(w),'diff','--binary','--full-index',BASE,'--'])
  check(b,name+' patch');patchpath=name+'/tracked-vs-main.patch';member(t,patchpath,b)
  extra=[]
  untracked=subprocess.check_output(['git','-C',str(w),'ls-files','--others','--exclude-standard','-z']).decode().split('\0')
  for s in untracked:
   if not s:continue
   p=w/s
   if p.is_symlink() or not p.is_file():continue
   if s.split('/')[0] not in {'internal','cmd','tools','docs','LICENSES'}:continue
   if p.suffix not in {'.go','.py','.sh','.md','.json','.txt','.js','.java','.mod','.sum'}:continue
   raw=p.read_bytes();check(raw,name+'/'+s);member(t,name+'/untracked/'+s,raw);extra.append({'path':s,'bytes':len(raw),'sha256':sha(raw)})
  for p in sorted((w/'docs').rglob('*.md')) if (w/'docs').exists() else []:
   rel=p.relative_to(w)
   if not ('experiments' in rel.parts or p.name.startswith('performance')):continue
   raw=p.read_bytes();h=sha(raw);check(raw,str(rel));dest='design-notes/'+h[:16]+'-'+p.name
   if h not in notes:
    d=OUT/dest;d.parent.mkdir(parents=True,exist_ok=True);d.write_bytes(raw);notes[h]={'path':dest,'sha256':h,'sources':[]}
   notes[h]['sources'].append({'worktree':name,'head':head,'path':str(rel)})
  sources.append({'worktree':name,'head':head,'branch':branch,'base_main':BASE,'tracked_patch':patchpath,'tracked_patch_bytes':len(b),'tracked_patch_sha256':sha(b),'untracked':extra})
save(OUT/'source-manifest.json',{'base_main':BASE,'scope':'Historical and current prototype source snapshots as patches, not compiler changes applied to this PR. Apply each patch to a clean checkout of base_main; restore its listed untracked files. Runtime experiment flags and toolchains come from measurement metadata. Snapshot records working trees, including unfinished/rejected source; availability is not a correctness or speed endorsement.','worktrees':sources})
save(OUT/'design-notes.json',list(notes.values()))
archives=[]
for p in sorted(OUT.glob('*.tar.gz')):
 b=p.read_bytes();archives.append({'path':p.name,'bytes':len(b),'sha256':sha(b)})
save(OUT/'archives.json',archives)
print(json.dumps({'records':len(records),'datasets':len(datasets),'source_worktrees':len(sources),'unique_design_notes':len(notes),'archives':archives},indent=2))
