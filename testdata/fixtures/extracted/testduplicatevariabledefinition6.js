var ns = {};
/** @record */ function rec3(){}
/** @type {number} */ rec3.prototype.prop;
/** @record */ function rec4(){}
/** @type {!rec3} */ ns.x;
/** @type {!rec4} */ ns.x;
