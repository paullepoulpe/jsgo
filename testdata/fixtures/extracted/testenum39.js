/** @enum {Number} */ var MyEnum = {FOO: new Number(1)};
/** @param {MyEnum} x
 * @return {number} */
function f(x) { return x == MyEnum.FOO && MyEnum.FOO == x; }
