goog.provide('user.scope');
goog.require('a.b');
goog.scope(function() {
  var b = a.b;
  new b.C();
  new b.Sub();
  var m = goog.module.get('a.mod');
  var mm = goog.module.get('a.legacy');
  use(m.X, mm);
  /** @type {b.Sub} */ var s;
});
