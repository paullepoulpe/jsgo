/** @enum */ var Foo = {AAA: 1};
/** @param {number} x */ Foo.method = function(x) {};
function f() { Foo.method(true); }
