/** @constructor */ function Foo() {}
/** @constructor */ function Bar() {}
/** @return {Foo|Bar} */ function f() { return new Foo(); }
f().missing;
