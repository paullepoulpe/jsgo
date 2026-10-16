goog.module('FooWrapper');

const {Foo: FooFromBaseModule} = goog.require('base');
exports = FooFromBaseModule;
