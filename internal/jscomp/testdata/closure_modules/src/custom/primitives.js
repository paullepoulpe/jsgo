goog.provide('custom.primitives');
goog.provide('custom.primitives.Foo');
goog.provide('custom.primitives.Object');
/** @type {Object} */ custom.primitives.Object = {};
goog.require('goog.object');
goog.forwardDeclare('custom.Fwd');
var fwd = goog.forwardDeclare('custom.Fwd2');
goog.addDependency('a.js', ['a'], []);
var x = goog.addDependency('b.js', ['b'], []);
function f() { goog.addDependency('c.js', ['c'], []); goog.forwardDeclare('custom.Inner'); }
goog.weakUsage(custom.primitives);
goog.weakUsage();
goog.weakUsage(1, 2);
goog.weakUsage('str');
var a = goog.object.create('a', 1, 'b', 2);
var b = goog.object.create(1, 2, 3);
var c = goog.object.createSet('a', 'b', 3);
var d = goog.object.createSet('a', 'a');
var e = goog.object.createSet(4, '4');
var g = goog.object.create(x, 1);
var h = goog.object.createSet([1, 2]);
var i = goog.reflect.objectProperty('foo', obj);
var j = JSCompiler_renameProperty('foo.bar', obj);
var k = JSCompiler_renameProperty(foo, obj, 1, 2);
/** @constructor */ custom.primitives.Foo = function() { custom.primitives.Foo.base(this, 'constructor', 1); };
goog.inherits(custom.primitives.Foo, Base);
custom.primitives.Foo.prototype.m = function() { return custom.primitives.Foo.base(this, 'm', 2); };
custom.primitives.Foo.prototype.n = function() { return custom.primitives.Foo.base(this, 'x'); };
custom.primitives.Foo.prototype.o = function() { return custom.primitives.Foo.base(that, 'o'); };
custom.primitives.Foo.prototype.p = function() { return Other.base(this, 'p'); };
function Ctor() { Ctor.base(this, 'constructor'); }
goog.inherits(Ctor, BaseCtor);
function Ctor2() { Ctor2.base(this, 'notconstructor'); }
goog.inherits(Ctor2, BaseCtor);
class K extends Base { m() { K.base(this, 'm'); } }
var y = Array.isArray(z);
var bound = (function() {}).bind(this, 1, 2);
var bound2 = Function.prototype.bind.call(f, null, 3);
$jscomp.inherits(Sub, Super);
goog.asserts.fail('x');
