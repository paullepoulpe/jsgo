//// file: w5/varref/script.js
function earlyVar() { a; var a = 1; }
function redeclared() { var b; var b; }
function tdz() { x2 = 1; let x2; }
function reassignConst() { const c = 1; c = 2; }
function unusedAssign() { let u = 1; u = 2; }
function unusedVar() { var w = 1; }
function usedLater() { let z; z = 1; return z; }
function defaults(p = q, q = 1) {}
function closureRef() { function inner() { return later; } let later = 1; return inner; }
for (let i = 0; i < 3; i++) {}
switch (1) {
  case 1:
    let s = 1;
    break;
  case 2:
    s;
}
function blockFn() {
  {
    function hoisted() {}
  }
  hoisted();
}
function redeclParam(p) { var p; }
function letShadow() { let y = 1; { let y = 2; y; } return y; }
function catchRedecl() { try {} catch (e) { var e; } }
var top1 = 1;
var top1 = 2;
function typeofTdz() { typeof t; let t; }
class Early { m() { return Late; } }
class Late {}
new Before();
class Before {}
goog.scope(function() {
  var alias = goog.foo;
  alias;
});
label: for (;;) { break label; }

//// file: w5/varref/redecl_error.js
function letThenVar() { let dup; var dup; }
function constRedecl() { const k = 1; let k = 2; }

//// file: w5/varref/mod.js
goog.module('w5.mod');
exports.x;
exports = {};

//// file: w5/varref/es.js
import {imp} from './dep.js';
imp;
export let exported = 1;
exported = 2;
export {exported as alias};
