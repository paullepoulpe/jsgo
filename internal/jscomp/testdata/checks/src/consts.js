const a = 1;
a = 2;
a++;
/** @const */ var b = 1;
b = 2;
b += 3;
--b;
/** @const */ var c = 1;
/** @suppress {const} */ c = 2;
FOO = 1;
FOO = 2;
const {d} = {d: 1};
d = 2;
function f() { const e = 1; e = 2; }
