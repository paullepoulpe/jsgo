/**
 * @constructor
 * @dict
 */
function Foo() {}
function Bar() {}
Bar.prototype = new Foo();
Bar.prototype.someprop = 123;
