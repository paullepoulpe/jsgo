var x = 0;
function incX() { return ++x; }
function add(a, b) { return a + b; }
function modifies(a, b) { a = 1; return a + b; }
function inner(a) { return function() { return a; }; }
function loops(a, b) { for (var i = 0; i < 2; i++) { foo(a); } return b; }
function rest(a, ...r) { return r.length + a; }
function cond(a, b) { return a && b; }
function thisUse(a) { return this.v + a; }
add(1, 2);
add(x, incX());
add(incX(), x);
add([], {});
modifies(x, 'str');
inner(incX());
loops([], 'ab');
rest(1, 2, 3);
rest();
cond(x, incX());
thisUse.call({v: 1}, 2);
add('a', 'b', 'c', 'd');
/** @const */ var CONSTANT = 1;
let letOnce = 2, letTwice = 3;
letTwice = 4;
const constVar = 5;
var assignedLater;
assignedLater = 6;
for (var loopVar = 0; loopVar < 1; loopVar++) {}
var castMe = /** @type {number} */ (x);
var renamed = JSCompiler_renameProperty('propName', {propName: 1});
var f1 = function() { return 1; };
let f2 = function named() { return named; };
const f3 = function() {};
var a1 = 1; var a2 = 2; let b1 = 1; let b2 = 2;
var x2 = 0; for (; x2 < 1; x2++) {}
var y2; for (y2 in {}) {}
x = x + 1; x = x * 2; x = x || 3;
