/**
 * @constructor
 * @struct
 */
function Struct1(){ this.prop = 123; }
/**
 * @constructor
 * @extends {Struct1}
 */
function Struct1kid(){ this.prop = 123; }
/** @param{Struct1kid} x */
function takesStruct2(x) { return x['prop']; }
