var [a, b] = arr;
var [c, , d] = arr;
var [e = 1, [f, g] = [], ...h] = arr;
var {i, j: k, l = 1, m: {n}, ['o' + 1]: p, ...q} = obj;
let {r: [s, ...t]} = obj;
const {u = 1} = obj;
[a, b] = [b, a];
({a, b} = obj);
({a: x.y, b: z[0], c = 3} = obj);
[x.y, z[0] = 1, ...w.v] = arr;
for (const [k2, v2] of map) {}
for ({a, b} of list) {}
for ([a, b] of list) {}
try {} catch ({message}) {}
try {} catch ([first]) {}
function fd({a = 1, b: {c} = {}} = {}) {}
var {if: kw} = obj;
({if: kw2} = obj);
