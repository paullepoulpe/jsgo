// Extra externs of the w5-checks oracle (tools/oracle/java/W5ChecksDump.java,
// internal/jscomp/w5_checks_oracle_test.go), after
// tools/oracle/typecheck-externs.js: what CheckConformanceTest adds to the
// default externs, plus a minimal DOM for the types the conformance configs
// name (testdata/w5_checks/conformance-*.textproto).

/** @type {Function} */ Arguments.prototype.callee;
/** @type {Function} */ Arguments.prototype.caller;
/** @type {Arguments} */ var arguments;

/**
 * @constructor
 * @param {*=} opt_message
 * @param {*=} opt_file
 * @param {*=} opt_line
 * @return {!Error}
 */
function Error(opt_message, opt_file, opt_line) {}
/** @type {string} */ Error.prototype.message;
/** @type {string} */ Error.prototype.stack;

/** @constructor */ function ObjectWithNoProps() {}

/** @param {*} message @return {undefined} */
function alert(message) {}
/** @param {string} code @return {*} */
function eval(code) {}
/** @param {*} s @param {number=} radix @return {number} */
function parseInt(s, radix) {}
/**
 * @param {!Function|string} f
 * @param {number=} delay
 * @param {...*} var_args
 * @return {number}
 */
function setTimeout(f, delay, var_args) {}

/** @constructor */ function Location() {}
/** @type {string} */ Location.prototype.href;

/** @constructor */ function Window() {}
/** @type {!Location} */ Window.prototype.location;
/** @type {!Window} */ var window;

/** @constructor */ function Node() {}
/** @type {?Node} */ Node.prototype.parentNode;

/** @constructor @extends {Node} */ function Element() {}
/** @type {string} */ Element.prototype.innerHTML;
/** @type {string} */ Element.prototype.outerHTML;
/** @type {string} */ Element.prototype.href;
/** @type {string} */ Element.prototype.attr;
/** @param {string} name @param {?} value @return {undefined} */
Element.prototype.setAttribute = function(name, value) {};

/** @constructor @extends {Element} */ function HTMLScriptElement() {}
/** @type {string} */ HTMLScriptElement.prototype.src;

/** @constructor @extends {Node} */ function Document() {}
/** @type {string} */ Document.prototype.domain;
/** @param {...*} var_args @return {undefined} */
Document.prototype.write = function(var_args) {};
/** @param {string} tagName @return {!Element} */
Document.prototype.createElement = function(tagName) {};
/**
 * @param {string} command
 * @param {boolean=} showUi
 * @param {*=} value
 * @return {boolean}
 */
Document.prototype.execCommand = function(command, showUi, value) {};
/** @type {!Document} */ var document;

/** @constructor */ function StyleSheet() {}
/** @type {string} */ StyleSheet.prototype.href;
/** @constructor */ function CSSImportRule() {}
/** @type {string} */ CSSImportRule.prototype.href;
