var x = {};
/**
 * @type {function(boolean): undefined}
 */ x.abc = function(x) {};
x = {
  /**
   * @param {boolean} x
   * @suppress {duplicate}
   */
  abc: function(x) {}
};
