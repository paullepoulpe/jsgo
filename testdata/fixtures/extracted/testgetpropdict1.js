/**
 * @constructor
 * @dict
 */
function Dict1(){ this['prop'] = 123; }
/** @param{Dict1} x */
function takesDict(x) { return x.prop; }
