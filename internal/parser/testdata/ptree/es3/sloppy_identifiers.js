var let = 1, yield = 2, static = 3;
var implements, interface, package, private, protected, public;
let = 4;
yield = 5;
yield: for (;;) break yield;
function f(let, yield, static) { return let + yield + static; }
function* g() { yield 1; var x = yield; }
var o = {let: 1, yield: 2, static: 3};
o.let; o.yield;
let [a] = [1];
let
b = 2;
class C { static static() {} }
var async = 1, await = 2, of = 3, get = 4, set = 5;
import(x);
for (let in o) {}
