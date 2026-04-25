/** @interface */
function Foo() {}
/** @return {number} */
Foo.prototype.method = function() {};
/**
 * @constructor
 * @implements {Foo}
 */
function Bar() {}
Bar.prototype.method;
var /** null */ n = (new Bar).method();
