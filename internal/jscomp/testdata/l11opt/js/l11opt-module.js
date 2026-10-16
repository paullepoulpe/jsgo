var dep = function(a, b) { return a + b; };
export const value = 'module string', other = 'module string';
export function helper(a) { var o = {p: 1}; return dep(o.p, a); }
