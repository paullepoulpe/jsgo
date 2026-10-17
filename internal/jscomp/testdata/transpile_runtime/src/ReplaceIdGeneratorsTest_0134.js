/** @const */
var $jscomp = $jscomp || {};
/** @const */
$jscomp.scope = {};
/**
 * @nosideeffects
 * @noinline
 * @param {!ITemplateArray} arrayStrings
 * @return {!ITemplateArray}
 */
$jscomp.createTemplateTagFirstArg = function(arrayStrings) {
  return $jscomp.createTemplateTagFirstArgWithRaw(arrayStrings, arrayStrings)
};
/**
 * @nosideeffects
 * @noinline
 * @param {!ITemplateArray} arrayStrings
 * @param {!ITemplateArray} rawArrayStrings
 * @return {!ITemplateArray}
 */
$jscomp.createTemplateTagFirstArgWithRaw = function(arrayStrings, rawArrayStrings) {
  arrayStrings.raw = rawArrayStrings;
  Object.freeze && (Object.freeze(arrayStrings), Object.freeze(rawArrayStrings));
  return (/** @type {!ITemplateArray} */ (arrayStrings));
};
/** @idGenerator {stable} */ function tagId(strings) {return strings[0];};
function Foo() { 'AAGMxg'; }
