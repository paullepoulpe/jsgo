/** @constructor */ var MainWidgetCreator = function() {};
/** @param {Function} ctor */
function createMainWidget(ctor) {
  /** @constructor */ function tempCtor() {};
  tempCtor.prototype = ctor.prototype;
  MainWidgetCreator.superClass_ = ctor.prototype;
  MainWidgetCreator.prototype = new tempCtor();
}
