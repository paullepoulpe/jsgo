/**  @constructor */ function Foo() {};
Foo.prototype.str = 'foo';

var /** {str: string, opt_num: (undefined|number)} */ x = new Foo;
