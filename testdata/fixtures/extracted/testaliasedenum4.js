/** @enum */ var YourEnum = {FOO: 3};
/** @enum */ var MyEnum = YourEnum;
/** @param {YourEnum} x */ function f(x) {} f(YourEnum.FOO);
