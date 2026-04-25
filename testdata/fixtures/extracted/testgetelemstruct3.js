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
var x = (new Struct1kid())['prop'];
