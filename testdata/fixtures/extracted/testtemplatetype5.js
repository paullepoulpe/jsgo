const CGI_PARAM_RETRY_COUNT = 'rc';

/**
 * @param {...T} p
 * @return {T}
 * @template T
 */
function fn(p) { return p; }
/** @type {!Object} */ var x;

/** @return {void} */
function aScope() {
  x = fn(CGI_PARAM_RETRY_COUNT, 1);
}
