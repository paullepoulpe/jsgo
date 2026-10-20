/** Matching-release Soy examples and runtime, with observable HTML escaping. */
goog.provide('enterprise.soy');
goog.require('soy.examples.simple');
goog.require('soy.examples.features');

// Retain the genuine generated feature examples as a public template library.
globalThis['jsgoEnterpriseSoyTemplates'] = soy.examples.features;
globalThis['jsgoEnterpriseSoy'] = function() {
  return String(soy.examples.simple.helloName({'name': '<Ada>'})) + '|' +
      String(soy.examples.simple.helloWorld());
};
