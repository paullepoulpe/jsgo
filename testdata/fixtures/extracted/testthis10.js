/** @constructor */function A(){};
A.prototype.foo = 3;
/** @this {A}
@return {string} */
A.bar = function() { return this.foo; };
