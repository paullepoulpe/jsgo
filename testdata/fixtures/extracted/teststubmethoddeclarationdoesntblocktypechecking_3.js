/** @interface */
var Foo = function() {};
/** @type {number} */
Foo.prototype.num;
/**
 * @constructor
 * @implements {Foo}
 */
var Bar = function() {};
/** @type {?} */
Bar.prototype.num;
var /** string */ x = (new Bar).num;
