/**
 * @constructor
 * @struct
 */
function Parent() {}
/**
 * @constructor
 * @extends {Parent}
 */
function Kid() {}
Kid.prototype.foo = 123;
var x = (new Kid()).foo;
