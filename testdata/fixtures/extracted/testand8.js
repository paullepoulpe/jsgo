function f(/** (null | number | string) */ x) {
  (x && (typeof x === 'number')) && takesNum(x);
}
function takesNum(/** number */ n) {}
