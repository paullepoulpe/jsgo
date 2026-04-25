/** @record @template T */ function I() {}
/** @type {T} */ I.prototype.prop;
var ns = {}
/** @record */ function rec(){}
/** @type {I<number>} */ rec.prototype.prop;

/** @type {!rec} */ ns.x;
/** @type {{prop:I<string>}} */ ns.x;
