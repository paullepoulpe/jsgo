/** @constructor */
function Foo() {}
/** @type {number} */
Foo.prototype.myprop = 2;

/** @constructor @extends {Foo} */
function Bar() {}
/** @override @type {string} */
Bar.prototype.myprop;
