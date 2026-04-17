function f(/** (null | number | string) */ x) {
  (x && (typeof x === 'string')) && takesNum(x);
}
function takesNum(/** number */ n) {}
