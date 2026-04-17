/** @typedef {{self: !Bar}} */
var Foo;
/** @typedef {!Foo} */
var Bar;

/** @type {!Foo} */
var foo;

class Baz {
  constructor() {
    /** @type {!Baz} */ this.self = this;
  }
}

const /** !Baz */ x = foo;
