/** @constructor @param {number} x */
function Foo(x) { /** @type {number} */ this.x = x; this.y = 'a'; }
/** @return {string} */
Foo.prototype.bar = function() { return 'a'; };
Foo.prototype.baz = function(a, b) {};
/** @type {number} */ Foo.prototype.count;
Foo.prototype.untyped;
/** @const */ Foo.staticConst = 3;
Foo.staticFn = function(/** string */ s) { return s; };

/** @constructor @extends {Foo} */
function SubFoo() { Foo.call(this, 1); }
goog.inherits(SubFoo, Foo);
/** @override */
SubFoo.prototype.bar = function() { return 'b'; };

/** @interface */
function IThing() {}
/** @return {number} */ IThing.prototype.size = function() {};
/** @record */
function IRec() {}
/** @type {string} */ IRec.prototype.name;

/** @constructor @implements {IThing} */
function Thing() {}
Thing.prototype.size = function() { return 1; };

var goog = goog || {};
goog.inherits = function(childCtor, parentCtor) {};
goog.addSingletonGetter = function(ctor) {};
goog.addSingletonGetter(Thing);
