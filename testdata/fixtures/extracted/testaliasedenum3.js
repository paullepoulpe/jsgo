/** @enum */ var YourEnum = {FOO: 3};
/** @enum */ var MyEnum = YourEnum;
/** @param {MyEnum} x */ function f(x) {} f(YourEnum.FOO);
