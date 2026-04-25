/** @enum */ var YourEnum = {FOO: 3};
/** @enum */ var MyEnum = YourEnum;
/** @param {string} x */ function f(x) {} f(MyEnum.FOO);
