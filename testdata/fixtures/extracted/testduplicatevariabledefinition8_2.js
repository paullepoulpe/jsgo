var ns = {}
/** @record */ function rec(){}
/** @type {number} */ rec.prototype.prop;

/** @type {!rec} */ ns.x;
/** @type {{prop:string}} */ ns.x;

/** @type {{prop:number}} */ ns.y;
/** @type {!rec} */ ns.y;
