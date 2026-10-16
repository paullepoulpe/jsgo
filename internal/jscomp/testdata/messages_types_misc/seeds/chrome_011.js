cr.define('namespace', function() {
  function internalStaticMethod() {
    alert(42);
  }
  function letsUseIt() {
    internalStaticMethod();
  }
  return {
    externalStaticMethod: internalStaticMethod
  };
});
