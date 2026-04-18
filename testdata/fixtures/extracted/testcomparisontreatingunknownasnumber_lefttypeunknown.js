/** @constructor */
function Bar() {}
function f(/** ? */ x) {
  return x < (new Bar);
}
