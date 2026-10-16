function Foo() {}; Foo.prototype.doSomething = goog.abstractMethod;
function Bar() {}; goog.inherits(Bar, Foo);
Bar.prototype.doSomething = function() {}
