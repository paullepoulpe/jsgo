function rest(a, ...others) { return others.length + a; }
var arr = [1, ...[2, 3], ...rest, 4];
foo(...arr);
foo(1, ...arr, 2);
obj.method(...arr);
getObj().method(...arr);
new Klass(...arr);
var o1 = {a: 1, ...o2, b: 2, ...o3};
const {p, ...q} = o1;
function args() { return f(...arguments); }
