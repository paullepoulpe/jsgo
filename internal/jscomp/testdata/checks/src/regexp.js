var r = RegExp;
RegExp.$1;
RegExp.foo;
new RegExp('a');
RegExp('b');
x instanceof RegExp;
x == RegExp;
switch (x) { case RegExp: }
var ok = /a(b/;
var bad2 = /[z-a]/;
var good = /a+b/gi;
var bad3 = /a{2,1}/;
function f(RegExp) { return RegExp.x; }
