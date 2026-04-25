/** @typedef {{self: !Bar}} */
var Foo;
/** @typedef {!Foo} */
var Bar;

class Baz {
  constructor() {
    /** @type {!Baz} */ this.self = this;
  }
}

const /** !Foo */ x = new Baz();
