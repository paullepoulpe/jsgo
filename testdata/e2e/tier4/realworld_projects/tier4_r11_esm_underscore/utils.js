function map(arr, fn) { var res = []; for (var i = 0; i < arr.length; i++) res.push(fn(arr[i])); return res; }
function square(n) { return n * n; }
window['squared'] = map([1, 2, 3, 4], square);
