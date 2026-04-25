function abstractMethod() {};
/** @interface */var base = function() {};
/** @extends {base}\s
 @interface */ var Int = function() {}
/** @type {{bar : !Function}} */ var x;\s
/** @type {!Function} */ base.prototype.bar = abstractMethod;\s
/** @type {Int} */ var foo;
foo.bar();
