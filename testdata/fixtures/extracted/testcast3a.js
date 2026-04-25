/** @constructor */function Base() {}
/** @constructor @extends {Base} */function Derived() {}
var baseInstance = new Base();
/** @type {!Derived} */ var baz = baseInstance;
