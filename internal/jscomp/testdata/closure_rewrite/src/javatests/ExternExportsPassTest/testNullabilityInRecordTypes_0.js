/** @typedef {{ nonNullable: !Object, nullable: Object }} */
var foo;
/** @param {foo} record */
function x(record) {}
goog.exportSymbol('x', x);
