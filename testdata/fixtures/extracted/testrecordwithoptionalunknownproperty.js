/**  @constructor */ function Foo() {};
Foo.prototype.str = 'foo';

var /** {str: string, opt_unknown: (?|undefined)} */ x = new Foo;
