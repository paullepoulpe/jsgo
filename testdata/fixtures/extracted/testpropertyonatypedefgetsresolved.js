const a = {};
/** @typedef {{prop: string}} */
  a.b;
/** @typedef {string} */
  a.b.c;
const abAlias = a.b;
var /** !abAlias.whatIsThis */ y = 'foo';
