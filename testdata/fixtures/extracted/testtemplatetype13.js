/** @constructor
 * @template T
 */
function Foo() {}

/** @constructor
 * @template T
 * @extends {Foo<T>}
 */
function A() {}

var a1 = new A();
var a2 = /** @type {!A<string>} */ (new A());
var a3 = /** @type {!A<number>} */ (new A());
/** @type {!Foo<string>} */ var f1 = a1;
/** @type {!Foo<string>} */ var f2 = a2;
/** @type {!Foo<string>} */ var f3 = a3;
