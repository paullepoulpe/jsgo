// Independent untimed check of actual compiler output and its API wire payloads.
'use strict';
const fs = require('node:fs');
const vm = require('node:vm');
if (process.argv.length !== 5) throw new Error('expected compiled file, export, result');
const context = vm.createContext({console, TextEncoder, TextDecoder});
vm.runInContext(fs.readFileSync(process.argv[2], 'utf8'), context,
    {filename: process.argv[2], timeout: 30000});
if (typeof context[process.argv[3]] !== 'function') throw new Error('compiled API check export missing');
const result = vm.runInContext(process.argv[3]+'()', context, {timeout: 30000});
if (result !== process.argv[4]) throw new Error('actual generated API payload differs: '+result);
const wires = context['jsgoBigQueryV2ActualWires'];
if (!wires || !Array.isArray(wires['table']) || !Array.isArray(wires['job'])) throw new Error('actual binary payload evidence missing');
process.stdout.write(JSON.stringify({status: 'PASS', result, actualWires: wires})+'\n');
