var [a, , b = 2, ...rest] = [1, 2, 3, 4];
var {c, d: {e}, f = 5, ['g' + 1]: g, ...others} = {c: 1, d: {e: 2}};
function params({x, y = 1}, [z, w] = [], ...[q, r]) { return x + y + z + w + q + r; }
function defaults(a = 1, b = a + 1, c) { return a + b + c; }
let v1, v2;
[v1, v2] = [v2, v1];
({v1, v2} = {v1: 1, v2: 2});
for (const [k, val] of Object.entries({})) { console.log(k, val); }
for (var {length} in {abc: 1}) {}
try { throw {code: 1}; } catch ({code}) { console.log(code); }
const fn = ({a: {b: [c1, c2]}}) => c1 + c2;
var x = (([p1, p2] = foo()), p1);
