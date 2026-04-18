/**
 * @constructor
 * @struct
 */
function Struct1() { this.prop = 123; }
/** @constructor */
function NonStruct() { this.prop = 321; }
/** @param{(NonStruct|Struct1)} x */
function takesStruct(x) { return x['prop']; }
