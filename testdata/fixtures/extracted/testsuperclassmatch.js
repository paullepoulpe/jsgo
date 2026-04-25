/** @constructor */ var Foo = function() {};
/** @constructor
 * @extends Foo */ var Bar = function() {};
Bar.inherits = function(x){};
Bar.inherits(Foo);
