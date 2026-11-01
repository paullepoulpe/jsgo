/**
 * @fileoverview The behavioural oracle's TypeScript workload driver
 * (tools/behave/workloads.py), compiled with TypeScript's lib/typescript.js and
 * the generated file data (BEHAVE_TSC_FILES: [path, text], the lib .d.ts files
 * first). It type-checks and emits a small program in memory, twice (ES2015
 * with declarations and source maps; ES5 CommonJS), and prints every
 * diagnostic and a digest of every emitted file (the runner keeps 16 KB of
 * output), so that the builds' outputs are compared in full.
 */
(function() {
  function digest(s) {
    var h = 0x811c9dc5;
    for (var i = 0; i < s.length; i++) {
      h ^= s.charCodeAt(i);
      h = Math.imul(h, 0x01000193) >>> 0;
    }
    return s.length + ':' + ('0000000' + h.toString(16)).slice(-8);
  }
  var files = BEHAVE_TSC_FILES;
  var text = {};
  var roots = [];
  files.forEach(function(f) {
    text[f[0]] = f[1];
    roots.push(f[0]);
  });
  function run(label, options) {
    var written = [];
    var host = {
      getSourceFile: function(name, lang) {
        return Object.prototype.hasOwnProperty.call(text, name) ? ts.createSourceFile(name, text[name], lang, true) : undefined;
      },
      getDefaultLibFileName: function() { return '/lib/lib.d.ts'; },
      writeFile: function(name, data) { written.push([name, data]); },
      getCurrentDirectory: function() { return '/'; },
      getCanonicalFileName: function(f) { return f; },
      useCaseSensitiveFileNames: function() { return true; },
      getNewLine: function() { return '\n'; },
      fileExists: function(f) { return Object.prototype.hasOwnProperty.call(text, f); },
      readFile: function(f) { return text[f]; },
      directoryExists: function() { return true; },
      getDirectories: function() { return []; }
    };
    var program = ts.createProgram(roots, options, host);
    var diags = ts.getPreEmitDiagnostics(program);
    var result = program.emit();
    diags = diags.concat(result.diagnostics);
    print('== ' + label + ': ' + diags.length + ' diagnostic(s), ' + written.length + ' file(s)');
    diags.forEach(function(d) {
      var where = '';
      if (d.file) {
        var lc = d.file.getLineAndCharacterOfPosition(d.start);
        where = d.file.fileName + ':' + (lc.line + 1) + ':' + (lc.character + 1) + ' ';
      }
      print(where + 'TS' + d.code + ' ' + ts.flattenDiagnosticMessageText(d.messageText, '\n'));
    });
    written.forEach(function(w) {
      print('-- ' + w[0] + ' ' + digest(w[1]));
    });
  }
  run('es2015', {target: ts.ScriptTarget.ES2015, module: ts.ModuleKind.ES2015, strict: true, noLib: true,
                 declaration: true, sourceMap: true, newLine: ts.NewLineKind.LineFeed, outDir: '/out'});
  run('es5', {target: ts.ScriptTarget.ES5, module: ts.ModuleKind.CommonJS, strict: true, noLib: true,
              downlevelIteration: true, importHelpers: false, newLine: ts.NewLineKind.LineFeed, outDir: '/out5'});
  print('TSC-DONE');
})();
