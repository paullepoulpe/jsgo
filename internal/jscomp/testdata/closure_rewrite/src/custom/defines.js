/** @define {boolean} */
var COMPILED = false;
/** @define {boolean} */
goog.DEBUG = true;
/** @define {string} */
goog.LOCALE = goog.define('goog.LOCALE', 'en');
/** @define {number} */
const FEATURESET = goog.define('goog.FEATURESET_YEAR', 2012);
/** @define {boolean} */
const DERIVED = goog.define('my.DERIVED', !COMPILED && goog.DEBUG);
/** @define {number} */
const BAD_VALUE = goog.define('my.BAD', Math.random());
/** @define {Object} */
var BAD_TYPE = {};
/** @define {boolean} */
var TWICE = true;
TWICE = false;
var CLOSURE_DEFINES = {'my.DERIVED': true, 'goog.LOCALE': 'fr', 'bad': {}};
CLOSURE_DEFINES['my.extra'] = 3;
const notDefine = goog.define('my.notannotated', true);
goog.define('my.bare', 1);
function f() {
  /** @define {boolean} */ var LOCAL = true;
  return goog.define('my.local', false);
}
/** @const */ var CONST_VAL = 1;
/** @define {number} */ var FROM_CONST = CONST_VAL + 1;
/** @define {string} */ var BAD_NAME = goog.define('not a name', '');
