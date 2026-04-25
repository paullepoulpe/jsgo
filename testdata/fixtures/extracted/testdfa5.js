/** @return {string?} */ function u() {return 'a';}
/** @param {string?} x
@return {string} */ function f(x) {
while (!x) { x = u(); }
return x;
}
