/** @param {function(string, number)} x */
function f(x) {}
f((/** string */ str, num) => {
// TODO(b/123583824): this should be a type mismatch warning.
//    found   : number
//    expected: string
//  but the JSDoc on `str` blocks inference of `num`.
  const /** string */ newStr = num;
});
