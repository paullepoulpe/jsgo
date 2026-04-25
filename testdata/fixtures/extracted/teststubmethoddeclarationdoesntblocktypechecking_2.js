/** @constructor */
function Foo() {}
/** @return {number} */
Foo.prototype.method = function() {};
/**
 * @constructor
 * @extends {Foo}
 */
function Bar() {}
Bar.prototype.method;
var /** null */ n = (new Bar).method();
