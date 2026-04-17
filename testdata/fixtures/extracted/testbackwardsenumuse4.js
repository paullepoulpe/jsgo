/** @return {number} */ function f() { return MyEnum.FOO; }
/** @enum {string} */ var YourEnum = {FOO: 'x'};
/** @enum {string} */ var MyEnum = YourEnum;
