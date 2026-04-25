/** @param {{isVisible: boolean}} x */
function f(x){
  x.isVisible = false;
}
/** @constructor */
function Foo() {}
/**
 * @constructor
 * @extends {Foo}
 */
function SubFoo() {}
/** @type {boolean} */
SubFoo.prototype.isVisible = true;
/**
 * @param {Foo} x
 * @return {boolean}
 */
function g(x) { return x.isVisible; }
