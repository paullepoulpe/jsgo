/**
  @param {*=} opt_o
  @return {string|null}
 */
function f(opt_o) {
  if (opt_o) {
    if (typeof opt_o !== 'string') {
      return opt_o.toString();
    }
  }
  return null;
}
