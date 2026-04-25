var ns = {};
/** @constructor */ ns.Foo = function() {
  this.bar(3, 5);
};
ns.Foo.prototype = {
  bar: function(x) {}
};
