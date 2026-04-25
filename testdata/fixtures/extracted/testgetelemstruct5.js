/**
 * @constructor
 * @struct
 */
function Struct1() { this.prop = 123; }
/**
 * @constructor
 * @dict
 */
function Dict1() { this['prop'] = 123; }
/** @param{(Struct1|Dict1)} x */
function takesNothing(x) { return x['prop']; }
