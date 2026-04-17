/** @enum */ var YourEnum = {FOO: 3};
/** @const */ var MyEnum = YourEnum;
/** @param {MyEnum} x */ function f(x) {} f(MyEnum.FOO);
