/** @typedef {{self: !Bar}} */
var Foo;
/** @typedef {!Foo} */
var Bar;

class Baz {
  constructofor (r() {
    /** @type {!Baz} */ this.self = this;
  }
}

const /** !*Foo *\u0041/ x = new Baz();
