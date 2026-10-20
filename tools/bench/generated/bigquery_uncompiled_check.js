// Benchmark-owned setup check over unmodified upstream Closure/runtime inputs.
// This does not substitute for independent checks of Java and Go ADVANCED output.
'use strict';
const fs = require('node:fs');
const path = require('node:path');
if (process.argv.length !== 5) throw new Error('expected fixture, closure-library, expected result');
const fixture = path.resolve(process.argv[2]);
const closure = path.resolve(process.argv[3]);
const base = path.join(closure, 'closure/goog');
// The pinned source checkout omits generated deps.js. Register dependency
// metadata from the actual source declarations, without changing any runtime.
globalThis.CLOSURE_NO_DEPS = true;
require(path.join(base, 'bootstrap/nodejs.js'));

function sources(directory) {
  return fs.readdirSync(directory, {withFileTypes: true}).flatMap(item => {
    const file = path.join(directory, item.name);
    return item.isDirectory() ? sources(file) : file.endsWith('.js') ? [file] : [];
  });
}

const fixtureFiles = sources(path.join(fixture, 'upstream/js'))
    .concat(sources(path.join(fixture, 'messages')))
    .concat([path.join(fixture, 'bigquery_entry.js')]);
const files = sources(base).concat(sources(path.join(closure, 'third_party/closure/goog')))
    .filter(file => !file.endsWith('_test.js')).concat(fixtureFiles);
for (const file of files) {
  const source = fs.readFileSync(file, 'utf8');
  const provides = Array.from(source.matchAll(/^\s*goog\.provide\(['"]([^'"]+)['"]\)/gm), item => item[1]);
  const modules = Array.from(source.matchAll(/^\s*goog\.module\(['"]([^'"]+)['"]\)/gm), item => item[1]);
  const requires = Array.from(source.matchAll(/^\s*(?:(?:const|let|var)\s+[^=;]+\s*=\s*)?goog\.require\(['"]([^'"]+)['"]\)/gm), item => item[1]);
  if (provides.length + modules.length === 0) {
    if (fixtureFiles.includes(file)) throw new Error('expected genuine Closure provides: ' + file);
    continue;
  }
  goog.addDependency(path.relative(base, file).split(path.sep).join('/'), provides.concat(modules), requires,
      modules.length ? {'module': 'goog'} : {});
}
goog.require('enterprise.bigqueryV2');
const actual = globalThis['jsgoEnterpriseBigQueryV2']();
if (actual !== process.argv[4]) throw new Error('unexpected generated round trip: ' + actual);
process.stdout.write(JSON.stringify({status: 'PASS', scope: 'uncompiled real upstream runtime setup only',
  result: actual, generatedAndRuntimeInputs: fixtureFiles.length,
  actualWires: globalThis['jsgoBigQueryV2ActualWires']}) + '\n');
