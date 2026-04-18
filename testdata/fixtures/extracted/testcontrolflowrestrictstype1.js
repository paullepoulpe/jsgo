/** @return {String?} */ function f() { return null; }
/** @type {String?} */ var a = f();
/** @type {String} */ var b = new String('foo');
/** @type {null} */ var c = null;
if (a) {
  b = a;
} else {
  c = a;
}
