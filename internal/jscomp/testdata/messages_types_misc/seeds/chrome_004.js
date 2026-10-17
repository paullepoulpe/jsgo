cr.define('namespace', function() {
  var internalStaticMethod = function() {
    alert(42);
  }
  return {
    externalStaticMethod: internalStaticMethod
  };
});
