/**
 * @param {number} a
 * @param {string=} b
 * @param {...boolean} c
 */
function f(a, b, ...c) {}
/** @param {{x: number, y: string}} p */
function g({x, y: [z] = []}) {}
function h(a = 1, {b} = {}, [c, d]) {}
/** @param {number=} a @param {number} b */
function badOrder(a, b) {}
/** @param {...number} a @param {number} b */
function badVar(a, b) {}
/** @param {number} nothere */
function inexistent(x) {}
(function(x, y) { var q = x; })(3, 'a');
var outerNum = 3;
(function(n, ...r) {})(outerNum, 1, 2);
(function(m) {})(...[1]);
const arrow = (a, b) => a + b;
const arrow2 = () => { return 1; };
function opt_args(opt_x, var_args) {}
