/** @record @template T */
function rec(){}
/** @type {T} */ rec.prototype.value;

/** @type {rec<string>} */ ns.x;
/** @type {{value: string}} */ ns.x;
