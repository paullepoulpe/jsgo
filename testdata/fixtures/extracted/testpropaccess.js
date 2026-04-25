/** @param {*} x */var f = function(x) {
var o = String(x);
if (typeof o['a'] != 'undefined') { return o['a']; }
return null;
};
