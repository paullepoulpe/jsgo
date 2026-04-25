/** @return {(string|null)} */ function f() { return null; }
/** @type {(string|null)} */ var a = f();
/** @type {string} */ var b = 'foo';
/** @type {null} */ var c = null;
if (a) {
  b = a;
} else {
  c = a;
}
