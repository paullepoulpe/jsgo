$jscomp_sink_closure_unaware_impl(
  function() {
    // Note that there isn't a JSDoc closureUnaware annotation on the inner function, as
    // it was never created during parsing of the AST.
    (function() {
      window['foo'] = 10;
    }).call(globalThis);
  });
