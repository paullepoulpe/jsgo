/** @constructor */ function Foo() {}
(function(ctor) {
  /** @type {boolean} */ ctor.prop = true;
})(Foo);
/** @param {number} x */ function f(x) {}
f(Foo.prop);
