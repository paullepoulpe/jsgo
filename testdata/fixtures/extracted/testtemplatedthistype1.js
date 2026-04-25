/** @constructor */
function Foo() {}
/**
 * @this {T}
 * @return {T}
 * @template T
 */
Foo.prototype.method = function() {};
/**
 * @constructor
 * @extends {Foo}
 */
function Bar() {}
var g = new Bar().method();
/**
 * @param {number} a
 */
function compute(a) {};
compute(g);
