/** @param {boolean} x */ function f(x) {}
function g() {
  var y = null;
  if (y != null) {
    for (var i = 0; i < 10; i++) {
      f(y);
    }
  }
};
