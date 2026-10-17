goog.scope(function() {
  var Foo = ns.Foo;
  var x = {
    /** @param {Foo} foo */ y: function(foo) {}
  };
  x.y('');
});
