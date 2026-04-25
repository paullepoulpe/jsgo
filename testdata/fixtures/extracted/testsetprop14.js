/**
 * @constructor
 * @struct
 */
function Top() {}
/**
 * @constructor
 * @extends {Top}
 */
function Mid() {}
/** blah blah */
Mid.prototype.foo = function() { return 1; };
/**
 * @constructor
 * @struct
 * @extends {Mid}
 */
function Bottom() {}
/** @override */
Bottom.prototype.foo = function() { return 3; };
