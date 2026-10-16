function outer() {
  var a = 1;
  var b = 2;
  var once = 3;
  let twice = 1;
  twice = 2;
  var ns = {};
  function inner() { a = 2; ns.fn = function() {}; }
  b = 4;
  if (a) { ns.cond = function() {}; }
  ns.plain = function() {};
  ns.plain2 = class {};
  return function() { return b; };
}
var g1 = 1;
g1 = 2;
function noReturn() {}
function withReturn() { return 1; }
function throwOnly() { throw new Error('x'); }
async function asyncFn() {}
async function asyncRet() { return 1; }
function* gen() { yield 1; }
async function* agen() {}
