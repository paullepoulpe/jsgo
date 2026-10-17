cr.define('foo', function() {
  /**
   * @enum {string}
   */
  let DangerType = {
    NOT_DANGEROUS: 'NOT_DANGEROUS',
    DANGEROUS: 'DANGEROUS',
  };

  return {
    DangerType: DangerType,
  };
});
