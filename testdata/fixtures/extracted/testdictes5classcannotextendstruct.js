/** @struct @constructor*/
function StructParent() {}
/** @constructor @dict @extends {StructParent} */
function DictChild() {}
DictChild.prototype['prop'] = function() {};
