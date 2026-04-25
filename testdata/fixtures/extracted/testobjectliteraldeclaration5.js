var x = {
  /** @param {boolean} x */ abc: function(x) {}
};
/**
 * @param {boolean} x
 * @suppress {duplicate}
 */ x.abc = function(x) {};
