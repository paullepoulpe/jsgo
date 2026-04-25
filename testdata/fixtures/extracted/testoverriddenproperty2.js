/** @constructor */ function Foo() {
  /** @type {Object} */
  this.bar = {};
}
/**
 * @constructor
 * @extends {Foo}
 */ function SubFoo() {}
/**
 * @type {Array}
 * @override
 */
SubFoo.prototype.bar = [];
