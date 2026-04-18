/**
 * @constructor
 * @dict
 */
function Dict1() { this['prop'] = 123; }
/** @constructor */
function NonDict() { this.prop = 321; }
/** @param{(NonDict|Dict1)} x */
function takesDict(x) { return x.prop; }
