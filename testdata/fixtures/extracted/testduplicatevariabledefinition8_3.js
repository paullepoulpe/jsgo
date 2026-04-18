var ns = {}
/** @record */ function rec(){}
/** @type {string} */ rec.prototype.prop;

/** @type {!rec} */ ns.x;
/** @type {{prop:string}} */ ns.x;

/** @type {{prop:number}} */ ns.y;
/** @type {!rec} */ ns.y;
