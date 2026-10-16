cr.define('namespace', function() {
  /** I'm function's JSDoc */
  var internalStaticMethod = function() {
    alert(42);
  }
  return {
    externalStaticMethod: internalStaticMethod
  };
});
