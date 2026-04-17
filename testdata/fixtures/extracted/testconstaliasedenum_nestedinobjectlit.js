/** @enum */ var YourEnum = {FOO: 3};
const ns = {/** @const */ x: {/** @const */ MyEnum: YourEnum}}
/** @param {ns.x.MyEnum} x */ function f(x) {} f(ns.x.MyEnum.FOO);
