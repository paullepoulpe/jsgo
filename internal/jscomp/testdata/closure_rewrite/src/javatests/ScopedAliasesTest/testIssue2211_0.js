var ns = {};
var y = 1;
goog.scope(function () {
  ns.fact = function y(n) {
    return n == 1 ? 1 : n * y(n - 1);
  };
});
