class Test {
 /**
  * @return {!Object}
  *
  * @suppress {checkTypes}
  */
testGoogRequireDynamicStubbedAndWithLoadedModule() {
   goog.setImportHandlerInternalDoNotCallOrElse(() => Promise.resolve(null));
   goog.setUncompiledChunkIdHandlerInternalDoNotCallOrElse(s => s);

   goog.loadModule('goog.module("a.loaded.module"); exports.foo = 12;');

   return goog.requireDynamic('a.loaded.module')
       .then(({foo}) => assertEquals(foo, 12));
 } }
