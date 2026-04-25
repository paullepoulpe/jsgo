/** @interface */
function Peelable() {};
/** @return {undefined} */
Peelable.prototype.peel;
/**
 * @constructor
 * @struct
 */
function Fruit() {};
/**
 * @constructor
 * @extends {Fruit}
 * @implements {Peelable}
 */
function Banana() { };
function f() {};
/** @override */
Banana.prototype.peel = f;
