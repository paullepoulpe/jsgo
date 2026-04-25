const a = {};
/** @typedef {{prop: string}} */
  a.b;
/** @typedef {string} */
  a.b.c;
const aAlias = a;
var /** aAlias.b.c */ x;
