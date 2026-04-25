/** @record */ function I1() {};
/** @type {number|undefined} */ I1.prototype.prop;

/** @record @extends {I1} */ function I2() {}
/** @record @extends {I2} */ function I3() {}
/** @type {string} */ I3.prototype.prop2;
/** @constructor @implements {I3} */ function C() {
    /** @type {number} */ this.prop = 1;
    /** @type {string} */ this.prop2 = 'str';
}

/** @param {I3} a */ function fn(a) {}
