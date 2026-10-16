cr.define('namespace', function() {
  /** I'm function's JSDoc */
  function internalStaticMethod() {
    alert(42);
  }
  return {
    externalStaticMethod: internalStaticMethod
  };
});
