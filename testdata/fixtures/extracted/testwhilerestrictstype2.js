/** @param {number?} x
@return {number}*/
function f(x) {
/** @type {number} */var y = 0;
while (x) {
y = x;
x = x-1;
}
return y;
}
