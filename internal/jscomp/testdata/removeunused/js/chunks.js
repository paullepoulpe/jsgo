function usedOnlyLater() { return 1; }
var movable = function() { return 2; };
/** @constructor */
function Klass() {}
Klass.prototype.methodA = function() { return 3; };
Klass.prototype.methodB = function() { return usedOnlyLater(); };
var instance = new Klass();
var later = movable();
alert(instance.methodA());
alert(instance.methodB() + later);
var x = 1;
x = 2;
alert(x instanceof Klass);
