/** @constructor */
function Foo() { this.used = 1; this.unused = 2; }
Foo.prototype.method = function() { return this.used; };
Foo.prototype.unusedMethod = function() { alert(1); };
Foo.staticUnused = function() {};
Object.defineProperties(Foo.prototype, {
  definedUsed: { get: function() { return 1; } },
  definedUnused: { value: 3 }
});
class Bar extends Foo {
  constructor() { super(); this.field = 1; }
  barMethod() { return this.field; }
  static staticMethod() {}
  get getter() { return 2; }
}
var foo = new Foo();
alert(foo.method() + foo.definedUsed);
alert(new Bar().barMethod());
var obj = {a: 1, b: function() {}};
obj.a;
