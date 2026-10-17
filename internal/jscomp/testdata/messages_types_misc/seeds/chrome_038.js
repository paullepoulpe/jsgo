cr.define('settings', function() {
  var x = 0;
  class C {
    static create() { return new C; }
  }
  return { C: C };
});
