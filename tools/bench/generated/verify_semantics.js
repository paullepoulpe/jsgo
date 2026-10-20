// Independent, untimed runtime check. Does not affect compiler arguments/work.
// Usage: node verify_semantics.js <compiled-js> <export-name> <expected-result>
'use strict';
const fs = require('node:fs');
const vm = require('node:vm');
if (process.argv.length !== 5) throw new Error('expected output.js, export, result');
const context = vm.createContext({console, TextEncoder, TextDecoder});
vm.runInContext(fs.readFileSync(process.argv[2], 'utf8'), context,
    {filename: process.argv[2], timeout: 30000});
const entry = context[process.argv[3]];
if (typeof entry !== 'function') throw new Error('compiled public entry is missing');
const actual = vm.runInContext(process.argv[3] + '()', context, {timeout: 30000});
if (actual !== process.argv[4]) throw new Error('semantic result: ' + actual);
process.stdout.write(JSON.stringify({status: 'PASS', export: process.argv[3], result: actual}) + '\n');
