/** @interface
 * @template T
 */
function Foo() {}

/** @constructor
 * @implements {Foo<string>}
 */
function A() {}

/** @constructor
 * @implements {Foo<number>}
 */
function B() {}

/** @type {!Foo<string>} */ var a = new A();
/** @type {!Foo<string>} */ var b = new B();
