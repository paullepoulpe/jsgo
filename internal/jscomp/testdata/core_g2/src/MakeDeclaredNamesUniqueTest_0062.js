var foo = {};
for(;;) {
   var fn = (function f(bar){ return function() {}; })(foo);}
