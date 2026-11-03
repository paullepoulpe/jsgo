package jscomp

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func runPassWithThreads(
	t *testing.T,
	inputs []*SourceFile,
	numThreads int,
	configure func(opts *CompilerOptions),
	makePass func(c *Compiler) CompilerPass,
) (errs []string, warns []string, report string) {
	t.Helper()

	opts := NewCompilerOptions()
	opts.SetChecksOnly(true)
	opts.SetClosurePass(true)
	opts.SetNumParallelThreads(numThreads)
	if configure != nil {
		configure(opts)
	}

	var reportBuf bytes.Buffer
	comp := NewCompilerWithOutStream(jutil.NewPrintStream(&reportBuf))
	comp.DisableThreads()
	comp.Init([]*SourceFile{}, inputs, opts)
	root := comp.ParseInputs()
	if root.IsNil() || comp.HasErrors() {
		for _, e := range comp.GetErrors() {
			errs = append(errs, e.String())
		}
		return errs, nil, reportBuf.String()
	}
	externsRoot := root.GetFirstChild()
	jsRoot := root.GetSecondChild()
	makePass(comp).Process(externsRoot, jsRoot)

	for _, e := range comp.GetErrors() {
		errs = append(errs, e.String())
	}
	for _, w := range comp.GetWarnings() {
		warns = append(warns, w.String())
	}
	return errs, warns, reportBuf.String()
}

func TestCheckGoogModuleTypeScriptNameParallel_MultiFileBugFix(t *testing.T) {
	inputs := []*SourceFile{
		SourceFile_FromCode("gws/js/first.js", `goog.module('google3.gws.js.first');
exports = 1;
`),
		SourceFile_FromCode("gws/js/second.js", `goog.module('wrong.second');
exports = 2;
`),
		SourceFile_FromCode("javascript/apps/other/ignored.js", `goog.module('custom.name');
exports = 3;
`),
		SourceFile_FromCode("gws/js/fourth.js", `goog.module('wrong.fourth');
exports = 4;
`),
	}

	configure := func(opts *CompilerOptions) {
		opts.SetWarningLevel(DiagnosticGroups_LINT_CHECKS, CheckLevel_WARNING)
	}

	// Standalone CheckGoogModuleTypeScriptName:
	seqErrs, seqWarns, _ := runPassWithThreads(t, inputs, 1, configure, func(c *Compiler) CompilerPass {
		return NewCheckGoogModuleTypeScriptName(c)
	})
	parErrs, parWarns, _ := runPassWithThreads(t, inputs, 4, configure, func(c *Compiler) CompilerPass {
		return NewCheckGoogModuleTypeScriptName(c)
	})

	if len(seqErrs) != 0 || len(parErrs) != 0 {
		t.Fatalf("unexpected errors: seq=%v par=%v", seqErrs, parErrs)
	}
	if strings.Join(seqWarns, "\n") != strings.Join(parWarns, "\n") {
		t.Fatalf("standalone warnings mismatch:\nseq:\n%s\npar:\n%s",
			strings.Join(seqWarns, "\n"), strings.Join(parWarns, "\n"))
	}
	if len(seqWarns) != 2 {
		t.Fatalf("expected 2 warnings across multi-file input (second.js and fourth.js), got %d:\n%s",
			len(seqWarns), strings.Join(seqWarns, "\n"))
	}

	// CombinedCompilerPass execution (exercises cloneForParallel + ShouldTraverse per script):
	combSeqErrs, combSeqWarns, _ := runPassWithThreads(t, inputs, 1, configure, func(c *Compiler) CompilerPass {
		return NewCombinedCompilerPass(c, NewCheckGoogModuleTypeScriptName(c))
	})
	combParErrs, combParWarns, _ := runPassWithThreads(t, inputs, 4, configure, func(c *Compiler) CompilerPass {
		return NewCombinedCompilerPass(c, NewCheckGoogModuleTypeScriptName(c))
	})
	if len(combSeqErrs) != 0 || len(combParErrs) != 0 {
		t.Fatalf("unexpected combined errors: seq=%v par=%v", combSeqErrs, combParErrs)
	}
	if strings.Join(combSeqWarns, "\n") != strings.Join(combParWarns, "\n") ||
		strings.Join(combSeqWarns, "\n") != strings.Join(seqWarns, "\n") {
		t.Fatalf("combined warnings mismatch:\nseq:\n%s\ncombSeq:\n%s\ncombPar:\n%s",
			strings.Join(seqWarns, "\n"), strings.Join(combSeqWarns, "\n"), strings.Join(combParWarns, "\n"))
	}
}

func TestCheckExtraRequiresParallel_ParityAndJSDocUsages(t *testing.T) {
	inputs := []*SourceFile{
		SourceFile_FromCode("defs.js", `
goog.provide('ns.Alpha');
goog.provide('ns.Beta');
goog.provide('ns.Gamma');
goog.provide('ns.Unused');
`),
		SourceFile_FromCode("mod_a.js", `goog.module('app.modA');
const Alpha = goog.require('ns.Alpha');
const Beta = goog.requireType('ns.Beta');
const Unused = goog.require('ns.Unused');
/** @type {?function(!Beta): !Array<!Beta>} */
let fn = null;
exports = {val: Alpha, fn};
`),
		SourceFile_FromCode("es_mod.js", `
import UsedDefault, {usedNamed, unusedNamed} from './other.js';
export const res = UsedDefault(usedNamed);
`),
	}
	for i := 0; i < 6; i++ {
		inputs = append(inputs, SourceFile_FromCode(
			fmt.Sprintf("user_%d.js", i),
			fmt.Sprintf(`goog.module('app.user%[1]d');
const Gamma = goog.require('ns.Gamma');
const Unused%[1]d = goog.require('ns.Unused');
const Fwd = goog.forwardDeclare('ns.Alpha');
/** @param {!Gamma} g */
exports = function(g) { return g; };
`, i)))
	}

	configure := func(opts *CompilerOptions) {
		opts.SetWarningLevel(DiagnosticGroups_EXTRA_REQUIRE, CheckLevel_WARNING)
		opts.SetWarningLevel(DiagnosticGroups_MODULE_LOAD, CheckLevel_OFF)
	}

	seqErrs, seqWarns, _ := runPassWithThreads(t, inputs, 1, configure, func(c *Compiler) CompilerPass {
		return NewCheckExtraRequires(c, nil)
	})
	parErrs, parWarns, _ := runPassWithThreads(t, inputs, 4, configure, func(c *Compiler) CompilerPass {
		return NewCheckExtraRequires(c, nil)
	})

	if strings.Join(seqErrs, "\n") != strings.Join(parErrs, "\n") {
		t.Fatalf("errors mismatch:\nseq:\n%s\npar:\n%s", strings.Join(seqErrs, "\n"), strings.Join(parErrs, "\n"))
	}
	if strings.Join(seqWarns, "\n") != strings.Join(parWarns, "\n") {
		t.Fatalf("warnings mismatch:\nseq:\n%s\npar:\n%s", strings.Join(seqWarns, "\n"), strings.Join(parWarns, "\n"))
	}
	if len(seqWarns) < 8 {
		t.Fatalf("expected extra-require warnings across files, got %d:\n%s", len(seqWarns), strings.Join(seqWarns, "\n"))
	}
}

func TestCheckMissingRequiresParallel_ModulesLoadModuleAndGoogScope(t *testing.T) {
	inputs := []*SourceFile{
		SourceFile_FromCode("base_defs.js", `
goog.provide('lib.Widget');
/** @constructor */
lib.Widget = function() {};

goog.provide('lib.Helper');
/** @constructor */
lib.Helper = function() {};
`),
		SourceFile_FromCode("mod_def.js", `goog.module('lib.ModType');
/** @constructor */
exports = function() {};
`),
		SourceFile_FromCode("mod_consumer.js", `goog.module('app.consumer');
/** @template T */
class Box {
  /** @param {T} v */
  constructor(v) { this.v = v; }
}
/** @type {!Box<function(!lib.ModType): !lib.Widget>} */
let box;
const w = new lib.Widget();
exports = {box, w};
`),
		SourceFile_FromCode("bundled.js", `
goog.loadModule(function(exports) {
  'use strict';
  goog.module('app.bundled');
  const h = new lib.Helper();
  exports = h;
  return exports;
});
`),
		SourceFile_FromCode("scope_file.js", `
goog.provide('app.Scoped');
goog.scope(function() {
  const w = new lib.Widget();
  app.Scoped = w;
});
`),
	}
	for i := 0; i < 5; i++ {
		inputs = append(inputs, SourceFile_FromCode(
			fmt.Sprintf("extra_mod_%d.js", i),
			fmt.Sprintf(`goog.module('app.extra%[1]d');
const u = new lib.Helper();
exports = u;
`, i)))
	}

	configure := func(opts *CompilerOptions) {
		opts.SetWarningLevel(DiagnosticGroups_MISSING_REQUIRE, CheckLevel_WARNING)
		opts.SetWarningLevel(
			DiagnosticGroup_ForType(CheckMissingRequires_MISSING_REQUIRE_TYPE),
			CheckLevel_WARNING)
	}

	makePass := func(c *Compiler) CompilerPass {
		return CompilerPassFunc(func(externs, root rhino.Node) {
			NewGatherModuleMetadata(c, false, ModuleLoader_ResolutionMode_BROWSER).Process(externs, root)
			NewCheckMissingRequires(c, c.GetModuleMetadataMap()).Process(externs, root)
		})
	}

	seqErrs, seqWarns, _ := runPassWithThreads(t, inputs, 1, configure, makePass)
	parErrs, parWarns, _ := runPassWithThreads(t, inputs, 4, configure, makePass)

	if strings.Join(seqErrs, "\n") != strings.Join(parErrs, "\n") {
		t.Fatalf("errors mismatch:\nseq:\n%s\npar:\n%s", strings.Join(seqErrs, "\n"), strings.Join(parErrs, "\n"))
	}
	if strings.Join(seqWarns, "\n") != strings.Join(parWarns, "\n") {
		t.Fatalf("warnings mismatch:\nseq:\n%s\npar:\n%s", strings.Join(seqWarns, "\n"), strings.Join(parWarns, "\n"))
	}
	all := strings.Join(append(seqErrs, seqWarns...), "\n")
	for _, want := range []string{
		"JSC_MISSING_REQUIRE",
		"JSC_MISSING_REQUIRE_TYPE",
		"JSC_MISSING_REQUIRE_IN_PROVIDES_FILE",
	} {
		if !strings.Contains(all, want) {
			t.Fatalf("expected diagnostic %s in output, got:\n%s", want, all)
		}
	}
}
