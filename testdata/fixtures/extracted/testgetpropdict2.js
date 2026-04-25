/**
 * @constructor
 * @dict
 */
function Dict1(){ this['prop'] = 123; }
/**
 * @constructor
 * @extends {Dict1}
 */
function Dict1kid(){ this['prop'] = 123; }
/** @param{Dict1kid} x */
function takesDict(x) { return x.prop; }
