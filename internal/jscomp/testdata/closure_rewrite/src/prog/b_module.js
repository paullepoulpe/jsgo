goog.module('p.B');
const A = goog.require('p.A');
const {method} = goog.require('p.A');
const Sub = goog.require('p.A.Sub');
/** @param {!A} a */
function use(a) { method(); return new Sub(); }
exports = class extends A { constructor() { super(); use(this); } };
