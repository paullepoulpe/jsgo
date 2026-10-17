class SuperClass {
  /**
   * @param {number} num
   */
   constructor(num) {
   }
}
goog.exportSymbol('SuperClass', SuperClass);
class SubClass extends SuperClass {
}
goog.exportSymbol('SubClass', SubClass);
