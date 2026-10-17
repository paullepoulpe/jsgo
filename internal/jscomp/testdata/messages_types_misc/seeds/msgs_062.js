function f() {
  /**
   * @desc A message that demonstrates placeholders
   */
  a.b.MSG_SILLY = goog.getMsg(
      '{$adjective} ' + '{$someNoun}',
      {'adjective': adj, 'someNoun': noun});
}
