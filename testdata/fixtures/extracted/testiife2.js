/** @constructor */ function Foo() {}
(function(ctor) {
  /** @type {boolean} */ ctor.prop = true;
})(Foo);
/** @return {number} */ function f() { return Foo.prop; }
