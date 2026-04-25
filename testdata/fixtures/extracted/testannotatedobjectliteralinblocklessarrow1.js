function f(/** {g: function(number): undefined} */ x) {}
() => f({/** @param {string} x */ g(x) {}});
