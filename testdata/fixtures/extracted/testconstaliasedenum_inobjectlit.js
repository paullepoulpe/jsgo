/** @enum */ var YourEnum = {FOO: 3};
const ns = {/** @const */ MyEnum: YourEnum}
/** @param {ns.MyEnum} x */ function f(x) {} f(ns.MyEnum.FOO);
