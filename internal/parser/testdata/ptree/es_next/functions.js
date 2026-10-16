function f() {}
function g(a, b = 1, ...rest) { return a + b; }
function h({a, b: [c, d = 2]}, [e, , f] = [], {...g}) {}
function* gen() { yield 1; }
async function af() { await 1; return await af(); }
async function* agf() { yield await 1; }
var fe = function () {};
var nfe = function named() {};
var ge = function* () {};
var afe = async function () {};
var anfe = async function named2() {};
var a1 = x => x;
var a2 = (x) => x * 2;
var a3 = (x, y) => { return x + y; };
var a4 = () => {};
var a5 = (...args) => args;
var a6 = (a, ...b) => b;
var a7 = ({a}) => a;
var a8 = ([a]) => a;
var a9 = (a = 1, b = a) => b;
var a10 = (a,) => a;
var a11 = async x => x;
var a12 = async (x) => x;
var a13 = async () => {};
var a14 = async (a, b) => { await a; };
var a15 = x => y => z => x + y + z;
var a16 = () => ({});
f(a => a, b);
f(...args, ...more,);
new f(1, 2,);
(function () { 'use strict'; })();
(() => 1)();
