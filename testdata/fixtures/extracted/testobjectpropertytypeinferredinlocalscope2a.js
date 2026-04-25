/**
  @param {!Object} o
  @param {number?} x
  @return {string}
 */
function f(o, x) {
  o.x = 'a';
  if (x) {o.x = x;}
  return o.x;
}
