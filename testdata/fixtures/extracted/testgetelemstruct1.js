/**
 * @constructor
 * @struct
 */
function Struct1(){ this.prop = 123; }
/** @param{Struct1} x */
function takesStruct(x) {
  var z = x;
  return z['prop'];
}
