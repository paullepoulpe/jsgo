/** @interface */ function Foo(){}
/**
 * @constructor
 * @struct
 * @implements {Foo}
 */
function Bar(){ this.x = 123; }
var z = /** @type {Foo} */(new Bar())['x'];
