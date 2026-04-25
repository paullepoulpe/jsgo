/**
 * @param {function(this:T, ...)?} fn
 * @param {?T} opt_obj
 * @template T
 */
function baz(fn, opt_obj) {}
baz(function() {}, null);
