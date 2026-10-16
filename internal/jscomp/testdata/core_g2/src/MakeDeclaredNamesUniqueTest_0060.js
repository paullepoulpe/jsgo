var loopObjName = {};
for(;;) {
   var fn = (function f(loopObjName){ return function() {}; })(loopObjName);
}
