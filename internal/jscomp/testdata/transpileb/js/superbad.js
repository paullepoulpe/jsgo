class Q2 { set q(v) { super.x = v; } }
var obj = { __proto__: Object.prototype, m() { return super.toString(); } };
