/**
 * @param {!IArrayLike<T>} x
 * @return {T}
 * @template T
*/
function f(x) { return x[0]; }
function g(/** !Array<string> */ x) {
  var /** null */ y = f(x);
}
