function fn(dir) {
  var node = dir ? 1 : 2;
  fn(/** @type {number} */ (node));
}
