var internalName;
/**
 * @template T
 */
internalName = class {
  /**
   * @param {T} param1
   */
  constructor(param1) {
    /** @const */
    this.data = param1;
  }
};
goog.exportSymbol('externalName', internalName)
