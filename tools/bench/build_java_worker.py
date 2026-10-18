#!/usr/bin/env python3
"""Build a pinned Java worker with the same extern-template cache as jsgo.

Generates two small API overlays from the SHA-verified Apache-2.0 Closure source
archive. The original jar and stock Java benchmark remain unchanged. Put this
build directory before the jar on the worker's classpath only.
"""
import argparse
import hashlib
import json
import os
from pathlib import Path
import subprocess
import sys
import zipfile
sys.path.insert(0, str(Path(__file__).resolve().parent))
sys.dont_write_bytecode = True
import bench

ap = argparse.ArgumentParser(description=__doc__)
ap.add_argument('--cache', default=bench.default_cache())
ap.add_argument('--out', required=True)
a = ap.parse_args()
pin = bench.read_pin()
cache, out = Path(a.cache).resolve(), Path(a.out).resolve()
archive = cache/'closure-src.jar'
if bench.sha256_file(archive) != pin['CLOSURE_SRC_SHA256']:
    sys.exit('Closure source archive SHA-256 does not match the pin')
jar = cache/pin['CLOSURE_JAR']
if bench.sha256_file(jar) not in pin['CLOSURE_JAR_SHA256'].split():
    sys.exit('Closure jar SHA-256 does not match the pin')
source = out/'src'
classes = out/'classes'
classes.mkdir(parents=True, exist_ok=True)
paths = ['com/google/javascript/jscomp/parsing/ParserRunner.java',
         'com/google/javascript/rhino/JSDocInfo.java']
with zipfile.ZipFile(archive) as z:
    for name in paths:
        s = z.read(name).decode()
        if name.endswith('/ParserRunner.java'):
            anchor = '  public static ParseResult parse(\n'
            assert s.count(anchor) == 1
            wrapper = '''  public static ParseResult parse(
      StaticSourceFile file, String source, Config config, ErrorReporter reporter) {
    return BuiltinExternCache.eligible(file, config)
        ? BuiltinExternCache.parse(file, source, config, reporter)
        : parseUncached(file, source, config, reporter);
  }

  static ParseResult parseUncached(
'''
            s = s.replace(anchor, wrapper)
        else:
            anchor = '  public JSDocInfo clone(boolean cloneTypeNodes) {'
            assert s.count(anchor) == 1
            helper = '''  // Benchmark overlay: copy parser facts as well as mutable type/TTL trees.
  public JSDocInfo cloneForParseTemplate() {
    Builder builder = toBuilder(type -> new JSTypeExpression(
        type.getRoot().cloneTree(false), type.getSourceName()));
    if (isInlineType()) { builder.recordInlineType(); }
    if (!getTypeTransformations().isEmpty()) {
      LinkedHashMap<String, Node> transforms = new LinkedHashMap<>();
      for (Map.Entry<String, Node> entry : getTypeTransformations().entrySet()) {
        transforms.put(entry.getKey(), entry.getValue().cloneTree(false));
      }
      builder.setProp(TYPE_TRANSFORMATIONS, transforms);
    }
    return builder.build();
  }

'''
            s = s.replace(anchor, helper+anchor)
        # The pinned Maven jar shades these dependencies; its sources do not.
        for package in ['com.google.common', 'com.google.errorprone']:
            s = s.replace(package, 'com.google.javascript.jscomp.jarjar.'+package)
        target = source/name
        target.parent.mkdir(parents=True, exist_ok=True)
        target.write_text(s)
files = [str(source/n) for n in paths]+[
    str(Path(bench.HERE)/'java/extern-cache/BuiltinExternCache.java'),
    str(Path(bench.HERE)/'java/CliWorker.java'),
    str(Path(bench.HERE)/'java/extern-cache/BuiltinExternCacheTest.java')]
subprocess.run(['javac', '-cp', str(jar), '-d', str(classes)]+files, check=True)
metadata = {'jar_sha256': bench.sha256_file(jar),
            'source_archive_sha256': bench.sha256_file(archive),
            'source_sha256': {str(Path(f).relative_to(source)) if Path(f).is_relative_to(source)
                              else str(Path(f).relative_to(bench.HERE)): bench.sha256_file(f) for f in files},
            'class_sha256': {str(p.relative_to(classes)): bench.sha256_file(p) for p in sorted(classes.rglob('*.class'))},
            'command': ['java', '-Djsgo.bench.externCache=true', '-cp', str(classes)+':'+str(jar), 'CliWorker']}
(out/'build.json').write_text(json.dumps(metadata, indent=2)+'\n')
print(json.dumps(metadata['command']))
