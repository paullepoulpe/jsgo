/** @interface */
class I1 { /** @return {number} */ m() {} }
/** @interface @extends {I1} */
class I2 extends I1 {}
/** @interface @extends {I1} */
function I3() {}
/** @record */
class R1 { constructor() { /** @type {string} */ this.s; } }
/** @implements {I1} */
class Impl { m() { return 1; } }
/** @implements {I2} @implements {R1} */
class Impl2 { m() { return 2; } }
/** @interface */
const IExpr = class {};
