/** @param {{foo: !Function}} x */ function f(x) {}
f({foo: function() {}});
