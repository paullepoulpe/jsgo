goog.loadModule(function(exports) {
   goog.module('a.b.c');

   class Foo {}
   /** @const {*} */ Foo.prototype.x;
   exports.Foo = Foo;

   return exports;
});
