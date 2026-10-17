/** @constructor */
function SuperClass() {}
goog.exportSymbol('Foo', SuperClass);

/**
 * @constructor
 * @extends {SuperClass}
 */
function SubClass() {}
goog.exportSymbol('Bar', SubClass);
