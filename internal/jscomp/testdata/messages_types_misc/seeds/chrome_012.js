cr.define('foo', function() {
  /**
   * @enum {string}
   */
  const DangerType = {
    NOT_DANGEROUS: 'NOT_DANGEROUS',
    DANGEROUS: 'DANGEROUS',
  };

  return {
    DangerType: DangerType,
  };
});
