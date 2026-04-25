/** @enum */ var YourEnum = {FOO: 3};
const ns = {};
/** @const */
ns.MyEnum = YourEnum;
/** @param {ns.MyEnum} x */ function f(x) {} f(ns.MyEnum.FOO);
