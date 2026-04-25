/** @interface */
class Foo {}
/** @interface */
class Bar {}
/** @typedef {!Foo|!Bar} */
let Baz;
/**
 * @interface
 * @extends {Baz}
 */
class Blah {}
