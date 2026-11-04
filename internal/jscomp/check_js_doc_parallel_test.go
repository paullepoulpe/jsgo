package jscomp

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func runCheckJSDocWithThreads(
	t *testing.T,
	externs []*SourceFile,
	inputs []*SourceFile,
	numThreads int,
	preserveDetailedSourceInfo bool,
) (errs []string, warns []string, report string, jsdocTypeStrings []string) {
	t.Helper()

	opts := NewCompilerOptions()
	opts.SetChecksOnly(true)
	opts.SetClosurePass(true)
	opts.SetNumParallelThreads(numThreads)
	opts.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_INCLUDE_DESCRIPTIONS_WITH_WHITESPACE)
	opts.SetPreserveDetailedSourceInfo(preserveDetailedSourceInfo)
	opts.SetWarningLevel(DiagnosticGroups_NON_STANDARD_JSDOC, CheckLevel_WARNING)
	opts.SetWarningLevel(DiagnosticGroups_MISPLACED_SUPPRESS, CheckLevel_WARNING)

	var reportBuf bytes.Buffer
	comp := NewCompilerWithOutStream(jutil.NewPrintStream(&reportBuf))
	comp.DisableThreads()
	comp.Init(externs, inputs, opts)
	root := comp.ParseInputs()
	if root.IsNil() || comp.HasErrors() {
		for _, e := range comp.GetErrors() {
			errs = append(errs, e.String())
		}
		return errs, nil, reportBuf.String(), nil
	}

	externsRoot := root.GetFirstChild()
	jsRoot := root.GetSecondChild()
	NewCheckJSDoc(comp).Process(externsRoot, jsRoot)

	for _, e := range comp.GetErrors() {
		errs = append(errs, e.String())
	}
	for _, w := range comp.GetWarnings() {
		warns = append(warns, w.String())
	}
	comp.GetErrorManager().GenerateReport()

	NodeUtil_VisitPreOrder(root, func(n rhino.Node) {
		info := n.GetJSDocInfo()
		if info == nil {
			return
		}
		for it := info.TypeExprIter(); it.Next(); {
			NodeUtil_VisitPreOrder(it.Node(), func(tn rhino.Node) {
				if str, ok := ast.AsStringLit(tn); ok {
					srcFile, _ := n.SourceFileName()
					jsdocTypeStrings = append(jsdocTypeStrings, fmt.Sprintf("%s:%d:%s", srcFile, n.GetLineno(), str.Value()))
				}
			})
		}
	})

	return errs, warns, reportBuf.String(), jsdocTypeStrings
}

func TestCheckJSDocParallelParity(t *testing.T) {
	externs := []*SourceFile{
		SourceFile_FromCode("externs_0.js", strings.Join([]string{
			"/** @fileoverview @externs */",
			"/* Ordinary extern block comment without annotations */",
			"/** @constructor */ function ExtElem() {}",
			"/** @type {string} @implicitCast */ ExtElem.prototype.innerHTML;",
			"/** @param {number} x @modifies {this} @nosideeffects */ ExtElem.prototype. mut = function(x) {};",
		}, "\n")),
		SourceFile_FromCode("externs_1.js", strings.Join([]string{
			"/** @fileoverview @externs */",
			"/* @type {number} */",
			"var extBlockCommentWarn;",
			"/** @export @type {number} */ var extMisplacedExport;",
		}, "\n")),
		SourceFile_FromCode("externs_2.js", strings.Join([]string{
			"/** @fileoverview @externs */",
			"/* Contact maintainer@example.com for details */",
			"/** @nocollapse */ ExtElem.prototype.badNoCollapse = function() {};",
		}, "\n")),
	}

	var inputs []*SourceFile
	for i := range 16 {
		name := fmt.Sprintf("script_%02d.js", i)
		var lines []string
		lines = append(lines,
			"/* Plain block comment with no at-sign at all, exercising fast-path skip. */",
			"/* Multi-line plain block comment\n * without any JSDoc tags.\n */",
			"/* Email reference dev@example.com with @ but not a JSDoc tag */",
		)
		switch i % 6 {
		case 0:
			lines = append(lines,
				"/* @type {string} */",
				fmt.Sprintf("var blockWarn_%d = 'hi';", i),
				fmt.Sprintf("/** @constructor */ var Ctor_%d = () => {};", i),
				fmt.Sprintf("/** @type {module$exports$pkg$mod_%d.Widget} */ var ref_%d = null;", i, i),
			)
		case 1:
			lines = append(lines,
				"/*\n * @param {number} x\n */",
				fmt.Sprintf("function fnDefault_%d(/** number */ x = 1) { return x; }", i),
				fmt.Sprintf("function fnRest_%d(/** number */ ...rest) { return rest.length; }", i),
			)
		case 2:
			lines = append(lines,
				fmt.Sprintf("/** @constructor */ function Cls_%d() {}", i),
				fmt.Sprintf("/** @type {string} @implicitCast */ Cls_%d.prototype.badCast = '';", i),
				fmt.Sprintf("/** @modifies {this} */ function badMod_%d() {}", i),
				fmt.Sprintf("/** @nosideeffects */ function badNoSE_%d() {}", i),
			)
		case 3:
			lines = append(lines,
				fmt.Sprintf("/** @define {boolean} */ let BAD_DEF_%d = true;", i),
				fmt.Sprintf("function retDoc_%d() { /** @return {number} */ return 1; }", i),
				fmt.Sprintf("/** @type {module$contents$pkg$mod_%d_Internal} */ var cref_%d = null;", i, i),
			)
		case 4:
			lines = append(lines,
				fmt.Sprintf("class Es6_%d {", i),
				"  /** @interface */ method() {}",
				"}",
				fmt.Sprintf("/** @usedViaDotConstructor */ function notACtor_%d() {}", i),
				fmt.Sprintf("var obj_%d = { /** @suppress {checkTypes} */ a: 1 };", i),
			)
		case 5:
			lines = append(lines,
				fmt.Sprintf("export class ModClass_%d { constructor() {} }", i),
				"/* @return {boolean} */",
				fmt.Sprintf("export function modFn_%d() { return true; }", i),
			)
		}
		inputs = append(inputs, SourceFile_FromCode(name, strings.Join(lines, "\n")))
	}

	for _, preserveComments := range []bool{true, false} {
		t.Run(fmt.Sprintf("preserveComments=%v", preserveComments), func(t *testing.T) {
			errs1, warns1, report1, types1 := runCheckJSDocWithThreads(t, externs, inputs, 1, preserveComments)
			errs4, warns4, report4, types4 := runCheckJSDocWithThreads(t, externs, inputs, 4, preserveComments)

			if len(errs1) == 0 || len(warns1) == 0 {
				t.Fatalf("expected non-empty errors and warnings, got %d errors and %d warnings", len(errs1), len(warns1))
			}
			if !slices.Equal(errs1, errs4) {
				t.Fatalf("errors mismatch:\nthreads=1:\n%s\nthreads=4:\n%s",
					strings.Join(errs1, "\n"), strings.Join(errs4, "\n"))
			}
			if !slices.Equal(warns1, warns4) {
				t.Fatalf("warnings mismatch:\nthreads=1:\n%s\nthreads=4:\n%s",
					strings.Join(warns1, "\n"), strings.Join(warns4, "\n"))
			}
			if report1 != report4 {
				t.Fatalf("report mismatch:\nthreads=1:\n%s\nthreads=4:\n%s", report1, report4)
			}
			if !slices.Equal(types1, types4) {
				t.Fatalf("rewritten JSDoc type strings mismatch:\nthreads=1:\n%s\nthreads=4:\n%s",
					strings.Join(types1, "\n"), strings.Join(types4, "\n"))
			}

			hasUnrecognized := false
			for _, ts := range types4 {
				if strings.Contains(ts, "UnrecognizedType_module$") {
					hasUnrecognized = true
					break
				}
			}
			if !hasUnrecognized {
				t.Fatalf("expected UnrecognizedType_module$ rewrite in JSDoc type strings, got %v", types4)
			}

			hasBlockCommentWarn := false
			for _, w := range warns4 {
				if strings.Contains(w, "JSC_JSDOC_IN_BLOCK_COMMENT") {
					hasBlockCommentWarn = true
					break
				}
			}
			if hasBlockCommentWarn != preserveComments {
				t.Fatalf("preserveComments=%v: hasBlockCommentWarn=%v, warnings=%v", preserveComments, hasBlockCommentWarn, warns4)
			}
		})
	}
}
