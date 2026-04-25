/** @param {number} x */ function f(x) {}
function g() {
  var y = null;
  for (var i = 0; i < 10; i++) {
    f(y);
    if (y != null) {
      // y is None the first time it goes through this branch
    } else {
      y = {};
    }
  }
};
