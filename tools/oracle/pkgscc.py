#!/usr/bin/env python3
# Prints the strongly connected components of the Java package import graph
# of the pinned Closure Compiler sources, dependencies first, with line counts.
# Each SCC becomes one Go package (PORTING.md section 2).
#
# usage: python3 -I tools/oracle/pkgscc.py .cache/closure-src
import os, re, sys, collections
root = sys.argv[1]
GEN = {"NodeProperty","TypedAst","AstNode","NodeKind","TypeProto","Type","ObjectTypeProto","SubtypingEdge","StringPoolProto","SourceFileProto","LazyAst","ScriptNode","OptimizationJsdoc","ColorPoolProto"}

pkg_files = collections.defaultdict(list)
for dp, dn, fn in os.walk(root):
    for f in fn:
        if f.endswith('.java'):
            p = os.path.join(dp, f)
            src = open(p, encoding='utf-8', errors='replace').read()
            m = re.search(r'^package\s+([\w.]+);', src, re.M)
            if m: pkg_files[m.group(1)].append((p, src))
pkgs = set(pkg_files)
edges = collections.defaultdict(set)
for pk, files in pkg_files.items():
    for p, src in files:
        for m in re.finditer(r"^import\s+(?:static\s+)?([\w.]+)\s*;", src, re.M):
            if ".serialization." in m.group(1) and m.group(1).rsplit(".",1)[-1] in GEN: continue
            name = m.group(1)
            # find longest package prefix
            parts = name.split('.')
            for i in range(len(parts)-1, 0, -1):
                cand = '.'.join(parts[:i])
                if cand in pkgs:
                    if cand != pk: edges[pk].add(cand)
                    break
        # fully-qualified references without import
        for m in re.finditer(r'\b(com\.google\.(?:javascript|debugging)\.[\w.]+?)\.[A-Z]\w*', src):
            cand = m.group(1)
            if cand in pkgs and cand != pk and not src[max(0,m.start()-8):m.start()].strip().endswith(('import','package')):
                edges[pk].add(cand)
# Tarjan SCC
idx = {}; low = {}; st = []; on = set(); out = []; c = [0]
sys.setrecursionlimit(10000)
def sc(v):
    idx[v] = low[v] = c[0]; c[0] += 1; st.append(v); on.add(v)
    for w in edges[v]:
        if w not in idx: sc(w); low[v] = min(low[v], low[w])
        elif w in on: low[v] = min(low[v], idx[w])
    if low[v] == idx[v]:
        comp = []
        while True:
            w = st.pop(); on.discard(w); comp.append(w)
            if w == v: break
        out.append(comp)
for v in sorted(pkgs):
    if v not in idx: sc(v)
def lines(pk): return sum(s.count('\n') for _, s in pkg_files[pk])
# out is in reverse topological order (dependencies first)
for comp in out:
    tot = sum(lines(p) for p in comp)
    short = [p.replace('com.google.javascript.','').replace('com.google.','') for p in sorted(comp)]
    print(f"{tot:7d}  {'  +  '.join(short)}")
