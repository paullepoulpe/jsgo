/** @constructor */ var Foo = function() {};
/** @constructor */ var Baz = function() {};
/** @constructor
 * @extends Foo */ var Bar = function() {};
Bar.inherits = function(x){};
Bar.mixin = function(y){};
Bar.inherits(Foo);
Bar.mixin(Baz);
