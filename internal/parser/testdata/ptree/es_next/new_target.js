function F() { if (new.target) {} return new.target; }
new F;
new F();
new F(1)(2);
new a.b.c();
new a[b]();
new new F()();
new (f())();
new F`t`;
new a.b`t`;
var x = new.foo;
