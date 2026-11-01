/**
 * @fileoverview The behavioural oracle's driver for marked (tools/behave/libs.py),
 * compiled with marked and the generated spec data (BEHAVE_MARKED_SPECS:
 * [suite, name, markdown, expected html, options as JSON]). Each spec is a
 * QUnit test (the runner's --env qunit) that renders the markdown and compares
 * it with the expected HTML, whitespace next to tags ignored and other runs of
 * whitespace as one (marked's own runner uses an HTML differ), or exactly with renderExact; every output is
 * also printed as a digest, so that the builds' outputs are compared in full,
 * the failing specs' included.
 */
(function() {
  var M = globalThis['marked'];
  var specs = BEHAVE_MARKED_SPECS;

  function digest(s) {
    var h = 0x811c9dc5;
    for (var i = 0; i < s.length; i++) {
      h ^= s.charCodeAt(i);
      h = Math.imul(h, 0x01000193) >>> 0;
    }
    return s.length + ':' + ('0000000' + h.toString(16)).slice(-8);
  }

  function norm(s) {
    return s.replace(/\s+/g, ' ').replace(/ ?(<[^>]*>) ?/g, '$1').trim();
  }

  specs.forEach(function(s) {
    QUnit.test(s[0] + '/' + s[1], function(assert) {
      var opts = JSON.parse(s[4]);
      var exact = !!opts['renderExact'];
      delete opts['renderExact'];
      var out;
      try {
        out = new M.Marked(opts).parse(s[2]);
      } catch (e) {
        print(s[0] + '/' + s[1] + ' threw ' + (e && e.name));
        throw e;
      }
      print(s[0] + '/' + s[1] + ' ' + digest(out));
      assert.ok(exact ? out === s[3] : norm(out) === norm(s[3]), 'html');
    });
  });
  QUnit.start();
})();
