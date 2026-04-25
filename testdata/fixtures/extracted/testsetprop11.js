/**
 * @constructor
 * @struct
 */
function Foo() {}
/** @constructor */
function Bar() {}
Bar.prototype = new Foo();
Bar.prototype.someprop = 123;
