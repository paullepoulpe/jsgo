goog.module('modA');

class Foo {}
// This is not a named exports object because of the value literal
exports = {Foo, Bar: [1,2,3]};
