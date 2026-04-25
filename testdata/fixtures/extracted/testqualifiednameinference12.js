/** @param {function(this:Object)} x */ function f(x) {}
/** @constructor */ function Foo() {
  /** @type {number} */ this.bar = 3;
  f(function() { this.bar = true; });
}
