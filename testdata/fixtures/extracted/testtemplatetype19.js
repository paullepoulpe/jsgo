/**
 * @param {T} t
 * @param {U} u
 * @return {{t:T, u:U}}\s
 * @template T,U
 */
function fn(t, u) { return {t:t, u:u}; }
/** @type {null} */ var x = fn(1, 'str');
