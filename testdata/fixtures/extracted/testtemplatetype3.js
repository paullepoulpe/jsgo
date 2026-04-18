/**
 * @param {T} v
 * @param {function(T)} f
 * @template T
 */
function call(v, f) { f.call(null, v); }
/** @type {string} */ var s;
call(3, function(x) {
 x = true;
 s = x;
});
