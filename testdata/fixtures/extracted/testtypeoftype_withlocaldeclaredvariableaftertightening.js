function f() {
  /** @type {string|number} */
  var x = 'str';
  if (typeof x == 'string') {
    /** @type {typeof x} */
    let g = null
  }
}
