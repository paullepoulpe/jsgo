/** @constructor
 * @template T
 */
function Foo() {}

/** @constructor
 * @extends {Foo<string>}
 */
function A() {}

/** @constructor
 * @extends {Foo<number>}
 */
function B() {}

/** @type {!Foo<string>} */ var a = new A();
/** @type {!Foo<string>} */ var b = new B();
