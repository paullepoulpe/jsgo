(function() {
  const c = goog.module.get('a.b.c');
  c;
  {
    const c = goog.module.get('other.a.b.c');
  }
})
