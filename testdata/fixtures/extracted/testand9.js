function f(/** (number|string|null) */ x) {
  if (x && typeof x === 'number') {
    takesNum(x);
  }
}
function takesNum(/** number */ x) {}
