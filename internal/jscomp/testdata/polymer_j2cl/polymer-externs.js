/** @constructor */
var Element = function() {};
/** @constructor @extends {Element} */
var HTMLElement = function() {};
/** @constructor @extends {HTMLElement} */
var HTMLInputElement = function() {};
/** @constructor @extends {HTMLElement} */
var PolymerElement = function() {};
/**
 * @typedef {{
 *   type: !Function,
 *   value: (* | undefined),
 *   readOnly: (boolean | undefined),
 *   computed: (string | undefined),
 *   reflectToAttribute: (boolean | undefined),
 *   notify: (boolean | undefined),
 *   observer: (string | function(this:?, ?, ?) | undefined)
 * }}
 */
let PolymerElementPropertiesMeta;

/**
 * @typedef {Object<string, !Function|!PolymerElementPropertiesMeta>}
 */
let PolymerElementProperties;
/** @type {!Object} */
PolymerElement.prototype.$;
PolymerElement.prototype.created = function() {};
PolymerElement.prototype.ready = function() {};
PolymerElement.prototype.attached = function() {};
PolymerElement.prototype.domReady = function() {};
PolymerElement.prototype.detached = function() {};
/**
 * Call the callback after a timeout. Calling job again with the same name
 * resets the timer but will not result in additional calls to callback.
 *
 * @param {string} name
 * @param {Function} callback
 * @param {number} timeoutMillis The minimum delay in milliseconds before
 *     calling the callback.
 */
PolymerElement.prototype.job = function(name, callback, timeoutMillis) {};
/**
 * @param a {!Object}
 * @return {!Function}
 */
var Polymer = function(a) {};
// In the actual code Polymer.Element comes from a mixin. This definition is simplified.
/**
 * @polymer
 */
Polymer.Element = class {};

/** @type {!Object<string, !Element>} */
Polymer.Element.prototype.$;
var alert = function(msg) {};
/** @const */ var $jscomp = $jscomp || {};
/** @const */ $jscomp.scope = {};
/**
 * @param {?Object} type
 * @param {T} object
 * @return {T}
 * @template T
 */
$jscomp.reflectObject = function (type, object) {};
