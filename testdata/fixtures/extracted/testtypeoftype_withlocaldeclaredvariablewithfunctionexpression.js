function f() {
  /** @type {string|number} */
  var x = 1;
  var g = function (/** typeof x */ a) {}
  x = 'str';
  g(null);
}
