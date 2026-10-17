goog.scope(function() {
  var Foo = ns.Foo;
  var x = {
    y: /** @param {Foo} foo */ function(foo) {}
  };
  x.y('');
});
