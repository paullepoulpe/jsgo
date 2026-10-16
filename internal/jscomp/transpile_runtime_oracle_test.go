package jscomp

// Go-only: differential test — the Go side of
// tools/oracle/java/TranspileRuntimeDump.java (the level-11 cluster
// l11-transpile-runtime: its passes and helper classes run on real files).
// trDumpFile reproduces, byte for byte, what the Java tool prints for one
// file (see the tool's header for the format of each section).
//
//   - TestTranspileRuntimeOracleGolden compares against the checked-in Java
//     dumps of the sources under testdata/transpile_runtime/src, listed in
//     testdata/transpile_runtime/files.txt (every JavaScript string of the
//     Java unit tests of the cluster's classes, and of
//     RewriteLogicalAssignmentOperatorsTest, that parses):
//     golden.txt.gz (with --full) for every section but force and
//     inject-tree, golden-inject.txt.gz (without --full: their trees hold
//     the injected runtime libraries) for those two. Regenerate with
//       cd internal/jscomp/testdata/transpile_runtime && \
//       java -Xss1g -cp closure-compiler.jar:<build> \
//         com.google.javascript.jscomp.TranspileRuntimeDump --full \
//         --sections=<all but force,inject-tree> golden.txt.gz src @files.txt
//     and --sections=force,inject-tree golden-inject.txt.gz without --full.
//   - TestTranspileRuntimeOracleCorpus compares the Java dumps named by
//     TRANSPILE_RUNTIME_ORACLE (comma-separated .txt.gz files, written
//     without --full) for the sources under TRANSPILE_RUNTIME_ROOT (default
//     /home/user/jsgo/.cache/corpus). The files, and each file's sections,
//     are the ones the dump holds. TRANSPILE_RUNTIME_DIFF_DIR, if set,
//     receives the Go dump of every mismatching file. The dumps of the
//     corpus are produced on deterministic samples of
//     .cache/corpus-files.txt (every 20th file for all sections but force
//     and inject-tree, every 100th file for those two), e.g.
//       awk 'NR%20==1' .cache/corpus-files.txt > sample20.txt
//       nice -n 10 java -Xss1g -cp closure-compiler.jar:<build> \
//         com.google.javascript.jscomp.TranspileRuntimeDump \
//         --sections=<all but force,inject-tree> corpus-sample20.txt.gz \
//         .cache/corpus @sample20.txt
//     The held-out corpus (.cache/corpus-holdout, corpus-holdout-files.txt)
//     is for a final check only.

import (
	"bufio"
	"compress/gzip"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ---------------------------------------------------------------------------
// Shared helpers (TranspileRuntimeDump's)

// trEscape is TranspileRuntimeDump.escape (ParsingDump.escape).
func trEscape(s string) string {
	var b strings.Builder
	for _, c := range jutil.StringToUTF16(s) {
		switch {
		case c == '\\':
			b.WriteString("\\\\")
		case c >= 0x20 && c < 0x7f:
			b.WriteByte(byte(c))
		default:
			const hex = "0123456789abcdef"
			b.WriteString("\\u")
			b.WriteByte(hex[c>>12&0xf])
			b.WriteByte(hex[c>>8&0xf])
			b.WriteByte(hex[c>>4&0xf])
			b.WriteByte(hex[c&0xf])
		}
	}
	return b.String()
}

// trExc is TranspileRuntimeDump.exc: the root cause's toString(), a
// NullPointerException's class name only.
func trExc(e jutil.Throwable) string {
	for e.GetCause() != nil {
		e = e.GetCause()
	}
	if _, ok := e.(*jutil.NullPointerException); ok {
		return "java.lang.NullPointerException"
	}
	return e.Error()
}

// trTreeLines is TranspileRuntimeDump.treeLines.
func trTreeLines(lines []string, n *rhino.Node, level int) []string {
	indent := strings.Repeat("    ", level)
	lines = append(lines, trEscape(indent+n.ToStringWithFlags(true, true, false)))
	if info := n.GetJSDocInfo(); info != nil {
		lines = append(lines, trEscape(indent+"  @jsdoc "+info.ToStringVerbose()))
	}
	if shadow := n.GetClosureUnawareShadow(); shadow != nil {
		lines = append(lines, trEscape(indent+"  @shadow"))
		lines = trTreeLines(lines, shadow, level+1)
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		lines = trTreeLines(lines, c, level+1)
	}
	return lines
}

func trTree(n *rhino.Node) []string { return trTreeLines(nil, n, 0) }

func trCrc(lines []string) uint32 {
	var c uint32
	for _, l := range lines {
		c = crc32.Update(c, crc32.IEEETable, []byte(l))
		c = crc32.Update(c, crc32.IEEETable, []byte{'\n'})
	}
	return c
}

func trSameLines(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

type trDumper struct {
	sb   strings.Builder
	full bool
}

func (d *trDumper) line(s string) {
	d.sb.WriteString(trEscape(s))
	d.sb.WriteByte('\n')
}

func (d *trDumper) raw(s string) {
	d.sb.WriteString(s)
	d.sb.WriteByte('\n')
}

// printTree is TranspileRuntimeDump.printTree.
func (d *trDumper) printTree(before []string, n *rhino.Node) {
	after := trTree(n)
	same := trSameLines(after, before)
	state := "changed"
	if same {
		state = "same"
	}
	d.raw("tree " + state + " " + strconv.Itoa(len(after)) + " " + strconv.FormatUint(uint64(trCrc(after)), 16))
	if d.full && !same {
		for _, l := range after {
			d.raw(l)
		}
	}
}

// trNewCompiler is TranspileRuntimeDump.newCompiler.
func trNewCompiler(rel string, src string, asExterns bool, settings func(*CompilerOptions)) *Compiler {
	compiler := NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	settings(options)
	file := SourceFile_FromCode(rel, src)
	if asExterns {
		compiler.Init([]*SourceFile{file}, []*SourceFile{}, options)
	} else {
		compiler.Init([]*SourceFile{}, []*SourceFile{file}, options)
	}
	compiler.Parse()
	return compiler
}

// trHasParseError is TranspileRuntimeDump.hasParseError.
func trHasParseError(compiler *Compiler) bool {
	for _, e := range compiler.GetErrors() {
		if e.GetType().key != "JSC_JS_MODULE_LOAD_WARNING" {
			return true
		}
	}
	return false
}

// printDiagnostics is TranspileRuntimeDump.printDiagnostics.
func (d *trDumper) printDiagnostics(compiler *Compiler, skipErrors int, skipWarnings int) {
	errors := compiler.GetErrors()
	for i := skipErrors; i < len(errors); i++ {
		e := errors[i]
		d.line("ERROR " + e.GetType().key + " " + strconv.Itoa(e.GetLineNumber()) + ":" +
			strconv.Itoa(e.GetCharno()) + " " + e.GetDescription())
	}
	warnings := compiler.GetWarnings()
	for i := skipWarnings; i < len(warnings); i++ {
		e := warnings[i]
		d.line("WARNING " + e.GetType().key + " " + strconv.Itoa(e.GetLineNumber()) + ":" +
			strconv.Itoa(e.GetCharno()) + " " + e.GetDescription())
	}
}

// trSection runs a section: it appends its result lines and returns the tree
// to print (nil for none).
type trSection func(c *Compiler, results *[]string) *rhino.Node

// trDumpFile is TranspileRuntimeDump.dumpFile.
func trDumpFile(rel string, src string, sections []string, full bool) string {
	d := &trDumper{full: full}
	d.line("# file " + rel)
	probe := trNewCompiler(rel, src, false, func(*CompilerOptions) {})
	if trHasParseError(probe) {
		d.line("# parse-error")
		return d.sb.String()
	}
	jsBefore := trTree(probe.GetJsRoot())
	var externsBefore []string
	if full {
		d.line("# before")
		for _, l := range jsBefore {
			d.raw(l)
		}
	}
	for _, section := range sections {
		d.line("## " + section)
		asExterns := section == "window" || section == "externprops"
		compiler := trNewCompiler(rel, src, asExterns, trSettings(section))
		if trHasParseError(compiler) {
			d.line("# parse-error")
			continue
		}
		errors0 := len(compiler.GetErrors())
		warnings0 := len(compiler.GetWarnings())
		if asExterns && externsBefore == nil {
			externsBefore = trTree(compiler.GetExternsRoot())
		}
		before := jsBefore
		if asExterns {
			before = externsBefore
		}
		var results []string
		var printed *rhino.Node
		var thrown jutil.Throwable
		run := trSectionNamed(section)
		jutil.Try(func() { printed = run(compiler, &results) }, func(e jutil.Throwable) { thrown = e })
		d.printDiagnostics(compiler, errors0, warnings0)
		if thrown != nil {
			d.line("! " + trExc(thrown))
			continue
		}
		for _, r := range results {
			d.line(r)
		}
		if printed != nil {
			d.printTree(before, printed)
		}
	}
	return d.sb.String()
}

// trSettings is TranspileRuntimeDump.settings.
func trSettings(section string) func(*CompilerOptions) {
	switch section {
	case "inject-es5":
		return func(o *CompilerOptions) { trInjectSettings(o, CompilerOptions_LanguageMode_ECMASCRIPT5, false) }
	case "inject-es2015":
		return func(o *CompilerOptions) { trInjectSettings(o, CompilerOptions_LanguageMode_ECMASCRIPT_2015, false) }
	case "inject-es2017":
		return func(o *CompilerOptions) { trInjectSettings(o, CompilerOptions_LanguageMode_ECMASCRIPT_2017, false) }
	case "inject-es2017-async":
		return func(o *CompilerOptions) { trInjectSettings(o, CompilerOptions_LanguageMode_ECMASCRIPT_2017, true) }
	case "inject-es2019":
		return func(o *CompilerOptions) { trInjectSettings(o, CompilerOptions_LanguageMode_ECMASCRIPT_2019, false) }
	case "inject-tree":
		return func(o *CompilerOptions) {
			o.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT5)
			o.SetRuntimeLibraryMode(RuntimeJsLibManager_RuntimeLibraryMode_INJECT)
		}
	case "force":
		return func(o *CompilerOptions) { o.SetRuntimeLibraryMode(RuntimeJsLibManager_RuntimeLibraryMode_INJECT) }
	case "template":
		return func(o *CompilerOptions) { o.SetRuntimeLibraryMode(RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY) }
	default:
		return func(*CompilerOptions) {}
	}
}

func trInjectSettings(o *CompilerOptions, out CompilerOptions_LanguageMode, async bool) {
	o.SetLanguageOut(out)
	o.SetInstrumentAsyncContext(async)
	o.SetRuntimeLibraryMode(RuntimeJsLibManager_RuntimeLibraryMode_RECORD_ONLY)
}

func trAddInjected(c *Compiler, results *[]string) {
	for _, lib := range c.GetRuntimeJsLibManager().GetInjectedLibraries() {
		*results = append(*results, "injected "+lib)
	}
}

func trVisibilityMap(c *Compiler) *jutil.LinkedHashMap[rhino.StaticSourceFile, rhino.JSDocInfo_Visibility] {
	pass := NewCollectFileOverviewVisibility(c)
	pass.Process(c.GetExternsRoot(), c.GetJsRoot())
	return pass.GetFileOverviewVisibilityMap()
}

// trAccessCallback is the access section's traversal.
type trAccessCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	visibilities *jutil.LinkedHashMap[rhino.StaticSourceFile, rhino.JSDocInfo_Visibility]
	names        int
	getprops     int
	crc          uint32
}

func (a *trAccessCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	var v rhino.JSDocInfo_Visibility
	if n.IsName() && n.GetString() != "" {
		var_ := t.GetScope().GetVar(n.GetString())
		if var_ == nil {
			return
		}
		v = AccessControlUtils_GetEffectiveNameVisibility(n, var_, a.visibilities)
		a.names++
	} else if n.IsGetProp() {
		v = AccessControlUtils_GetEffectivePropertyVisibility(n, nil, a.visibilities)
		a.getprops++
	} else {
		return
	}
	a.crc = crc32.Update(a.crc, crc32.IEEETable, []byte(v.String()+"\n"))
}

// trTemplateCallback is the template section's traversal.
type trTemplateCallback struct {
	converter      *Es6TemplateLiterals
	insertionPoint *rhino.Node
}

func (cb *trTemplateCallback) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if n.IsScript() {
		cb.insertionPoint = n.GetFirstChild()
	}
	return true
}

func (cb *trTemplateCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsTaggedTemplateLit() {
		cb.converter.VisitTaggedTemplateLiteral(t, n, cb.insertionPoint)
	} else if n.IsTemplateLit() && !parent.IsTaggedTemplateLit() {
		cb.converter.VisitTemplateLiteral(t, n)
	}
}

// trLogicalCallback is the logical section's traversal.
type trLogicalCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	helper *RewriteLogicalAssignmentOperatorsHelper
}

func (cb *trLogicalCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_ASSIGN_OR, rhino.Token_ASSIGN_AND, rhino.Token_ASSIGN_COALESCE:
		cb.helper.VisitLogicalAssignmentOperator(t, n)
	default:
	}
}

// trCtorsCallback is the ctors section's traversal.
type trCtorsCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	synthesizer *SynthesizeExplicitConstructors
}

func (cb *trCtorsCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsClass() && (n.GetSecondChild().IsEmpty() || n.GetSecondChild().IsQualifiedName()) {
		cb.synthesizer.SynthesizeClassConstructorIfMissing(t, n)
	}
}

// trArrowsCallback is the arrows section's traversal (Es6RewriteArrowFunction's
// use of ThisAndArgumentsReferenceUpdater).
type trArrowsCallback struct {
	c            *Compiler
	astFactory   *AstFactory
	contextStack []*ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext // top is last
}

func (cb *trArrowsCallback) peek() *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext {
	if len(cb.contextStack) == 0 {
		return nil
	}
	return cb.contextStack[len(cb.contextStack)-1]
}

func (cb *trArrowsCallback) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		cb.contextStack = append(cb.contextStack, NewThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext(
			n, false, cb.c.GetUniqueIdSupplier().GetUniqueId(t.GetInput())))
	case rhino.Token_FUNCTION:
		if !n.IsArrowFunction() {
			cb.contextStack = append(cb.contextStack, NewThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext(
				n.GetLastChild(), NodeUtil_IsEs6Constructor(n), cb.c.GetUniqueIdSupplier().GetUniqueId(t.GetInput())))
		}
	case rhino.Token_SUPER:
		context := cb.peek()
		if context != nil && context.isConstructor && parent.IsCall() && parent.GetFirstChild() == n {
			s := parent
			for s.GetParent() != context.scopeBody {
				s = NodeUtil_GetEnclosingStatement(s.GetParent())
			}
			context.lastSuperStatement = s
		}
	default:
	}
	return true
}

func (cb *trArrowsCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	context := cb.peek()
	if n.IsArrowFunction() {
		if n.GetLastChild().IsBlock() {
			n.SetIsArrowFunction(false)
			NodeTraversal_Traverse(cb.c, n.GetLastChild(), NewThisAndArgumentsReferenceUpdater(cb.c, context, cb.astFactory))
			t.ReportCodeChange()
		}
	} else if context != nil && context.scopeBody == n {
		cb.contextStack = cb.contextStack[:len(cb.contextStack)-1]
		context.AddVarDeclarations(cb.c, cb.astFactory, t.GetCurrentScript())
	}
}

// trSectionNamed is TranspileRuntimeDump.section.
func trSectionNamed(name string) trSection {
	switch name {
	case "substitute":
		return func(c *Compiler, r *[]string) *rhino.Node {
			NewSubstituteEs6Syntax(c).Process(c.GetExternsRoot(), c.GetJsRoot())
			return c.GetJsRoot()
		}
	case "shorthand":
		return func(c *Compiler, r *[]string) *rhino.Node {
			p := NewEs6NormalizeShorthandProperties(c)
			NodeUtil_VisitPostOrder(c.GetJsRoot(), NodeUtil_VisitorFunc(func(n *rhino.Node) { p.TranspileSubtree(n) }))
			return c.GetJsRoot()
		}
	case "catch":
		return func(c *Compiler, r *[]string) *rhino.Node {
			p := NewRewriteCatchWithNoBinding(c)
			NodeUtil_VisitPostOrder(c.GetJsRoot(), NodeUtil_VisitorFunc(func(n *rhino.Node) { p.TranspileSubtree(n) }))
			return c.GetJsRoot()
		}
	case "trycatch":
		return func(c *Compiler, r *[]string) *rhino.Node {
			NewRewriteGlobalDeclarationsForTryCatchWrapping(c).Process(c.GetExternsRoot(), c.GetJsRoot())
			return c.GetJsRoot()
		}
	case "inject-es5", "inject-es2015", "inject-es2017", "inject-es2017-async", "inject-es2019":
		return func(c *Compiler, r *[]string) *rhino.Node {
			NewInjectTranspilationRuntimeLibraries(c).Process(c.GetExternsRoot(), c.GetJsRoot())
			trAddInjected(c, r)
			return nil
		}
	case "inject-tree":
		return func(c *Compiler, r *[]string) *rhino.Node {
			NewInjectTranspilationRuntimeLibraries(c).Process(c.GetExternsRoot(), c.GetJsRoot())
			trAddInjected(c, r)
			return c.GetJsRoot()
		}
	case "force":
		return func(c *Compiler, r *[]string) *rhino.Node {
			NewInjectRuntimeLibraries(c, jutil.NewLinkedHashSet("es6/util/inherits", "es6/symbol")).
				Process(c.GetExternsRoot(), c.GetJsRoot())
			trAddInjected(c, r)
			return c.GetJsRoot()
		}
	case "unwrap":
		return func(c *Compiler, r *[]string) *rhino.Node {
			ManageClosureUnawareCode_Unwrap(c).Process(c.GetExternsRoot(), c.GetJsRoot())
			return c.GetJsRoot()
		}
	case "caller":
		return func(c *Compiler, r *[]string) *rhino.Node {
			NewRewriteCallerCodeLocation(c).Process(c.GetExternsRoot(), c.GetJsRoot())
			return c.GetJsRoot()
		}
	case "strip":
		return func(c *Compiler, r *[]string) *rhino.Node {
			c.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
			NewStripCode(
				c,
				jutil.NewLinkedHashSet(
					"goog.debug.DebugWindow",
					"goog.debug.FancyWindow",
					"goog.debug.Formatter",
					"goog.debug.HtmlFormatter",
					"goog.debug.TextFormatter",
					"goog.debug.Logger",
					"goog.debug.LogManager",
					"goog.debug.LogRecord",
					"goog.net.BrowserChannel.LogSaver",
					"goog.log",
					"GA_GoogleDebugger"),
				jutil.NewLinkedHashSet(
					"logger",
					"logger_",
					"debugWindow",
					"debugWindow_",
					"logFormatter_",
					"logBuffer_"),
				jutil.NewLinkedHashSet("trace"),
				true).
				Process(c.GetExternsRoot(), c.GetJsRoot())
			return c.GetJsRoot()
		}
	case "ids":
		return func(c *Compiler, r *[]string) *rhino.Node {
			idTestMap := &replaceIdGeneratorsTest_idTestMap{
				map_: jutil.NewImmutableMapBuilder[string, string]().
					Put("foo", ":foo:").
					Put("bar", ":bar:").
					BuildOrThrow(),
			}
			pass := NewReplaceIdGenerators(
				c,
				/* templateLiteralsAreTranspiled= */ false,
				jutil.NewImmutableMapBuilder[string, RenamingMap]().
					Put("goog.events.getUniqueId", RenamingToken_INCONSISTENT).
					Put("goog.place.getUniqueId", RenamingToken_INCONSISTENT).
					Put("id", idTestMap).
					Put("get.id", idTestMap).
					BuildOrThrow(),
				/* generatePseudoNames= */ false,
				/* previousMapSerialized= */ nil,
				/* xidHashFunction= */ nil)
			pass.Process(c.GetExternsRoot(), c.GetJsRoot())
			for _, l := range strings.Split(pass.GetSerializedIdMappings(), "\n") {
				*r = append(*r, "map "+l)
			}
			return c.GetJsRoot()
		}
	case "toggles":
		return func(c *Compiler, r *[]string) *rhino.Node {
			NewReplaceToggles(c).Process(c.GetExternsRoot(), c.GetJsRoot())
			return c.GetJsRoot()
		}
	case "window":
		return func(c *Compiler, r *[]string) *rhino.Node {
			NewDeclaredGlobalExternsOnWindow(c).Process(c.GetExternsRoot(), c.GetJsRoot())
			return c.GetExternsRoot()
		}
	case "externprops":
		return func(c *Compiler, r *[]string) *rhino.Node {
			NewGatherExternProperties(c, GatherExternProperties_Mode_CHECK_AND_OPTIMIZE).
				Process(c.GetExternsRoot(), c.GetJsRoot())
			for p := range c.GetExternProperties().All() {
				*r = append(*r, "prop "+p)
			}
			return nil
		}
	case "accessors":
		return func(c *Compiler, r *[]string) *rhino.Node {
			for k, v := range GatherGetterAndSetterProperties_Gather(c, c.GetRoot()).All() {
				*r = append(*r, "accessor "+k+" "+v.String())
			}
			return nil
		}
	case "visibility":
		return func(c *Compiler, r *[]string) *rhino.Node {
			for k, v := range trVisibilityMap(c).All() {
				*r = append(*r, "visibility "+k.GetName()+" "+v.String())
			}
			return nil
		}
	case "access":
		return func(c *Compiler, r *[]string) *rhino.Node {
			cb := &trAccessCallback{visibilities: jutil.NewLinkedHashMapFrom(trVisibilityMap(c))}
			NodeTraversal_Traverse(c, c.GetJsRoot(), cb)
			*r = append(*r, "access "+strconv.Itoa(cb.names)+" "+strconv.Itoa(cb.getprops)+" "+
				strconv.FormatUint(uint64(cb.crc), 16))
			return nil
		}
	case "template":
		return func(c *Compiler, r *[]string) *rhino.Node {
			libs := c.GetRuntimeJsLibManager()
			libs.InjectLibForField("$jscomp.createTemplateTagFirstArg")
			libs.InjectLibForField("$jscomp.createTemplateTagFirstArgWithRaw")
			NodeTraversal_Traverse(c, c.GetJsRoot(), &trTemplateCallback{converter: NewEs6TemplateLiterals(c)})
			return c.GetJsRoot()
		}
	case "logical":
		return func(c *Compiler, r *[]string) *rhino.Node {
			helper := NewRewriteLogicalAssignmentOperatorsHelper(c, c.CreateAstFactory(), c.GetUniqueIdSupplier())
			NodeTraversal_Traverse(c, c.GetJsRoot(), &trLogicalCallback{helper: helper})
			return c.GetJsRoot()
		}
	case "ctors":
		return func(c *Compiler, r *[]string) *rhino.Node {
			NodeTraversal_Traverse(c, c.GetJsRoot(), &trCtorsCallback{synthesizer: NewSynthesizeExplicitConstructors(c)})
			return c.GetJsRoot()
		}
	case "arrows":
		return func(c *Compiler, r *[]string) *rhino.Node {
			NodeTraversal_Traverse(c, c.GetJsRoot(), &trArrowsCallback{c: c, astFactory: c.CreateAstFactory()})
			return c.GetJsRoot()
		}
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg("unknown section " + name))
	}
}

// ---------------------------------------------------------------------------
// Reading the Java dumps

// trJavaFile is one file's block of a Java dump.
type trJavaFile struct {
	rel      string
	text     string
	sections []string
	full     bool
}

func trReadJavaDump(t *testing.T, path string) []*trJavaFile {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	var files []*trJavaFile
	var cur *trJavaFile
	var sb strings.Builder
	flush := func() {
		if cur != nil {
			cur.text = sb.String()
			files = append(files, cur)
		}
		sb.Reset()
	}
	sc := bufio.NewScanner(zr)
	sc.Buffer(make([]byte, 1<<20), 1<<28)
	for sc.Scan() {
		l := sc.Text()
		if strings.HasPrefix(l, "# file ") {
			flush()
			// Every char of a path is printable ASCII or escaped; decode the escapes.
			cur = &trJavaFile{rel: trUnescape(strings.TrimPrefix(l, "# file "))}
		} else if strings.HasPrefix(l, "## ") && cur != nil {
			cur.sections = append(cur.sections, strings.TrimPrefix(l, "## "))
		} else if l == "# before" && cur != nil {
			cur.full = true
		}
		sb.WriteString(l)
		sb.WriteByte('\n')
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	flush()
	return files
}

// trUnescape decodes trEscape.
func trUnescape(s string) string {
	var units []uint16
	for i := 0; i < len(s); {
		if s[i] == '\\' && i+1 < len(s) && s[i+1] == '\\' {
			units = append(units, '\\')
			i += 2
			continue
		}
		if s[i] == '\\' && i+5 < len(s) && s[i+1] == 'u' {
			v, err := strconv.ParseUint(s[i+2:i+6], 16, 16)
			if err == nil {
				units = append(units, uint16(v))
				i += 6
				continue
			}
		}
		units = append(units, uint16(s[i]))
		i++
	}
	return jutil.StringFromUTF16(units)
}

// trCompare runs the Go dump of every file of the Java dumps and reports the
// mismatches.
func trCompare(t *testing.T, root string, files []*trJavaFile, diffDir string) {
	t.Helper()
	mismatches := 0
	for _, jf := range files {
		src, err := os.ReadFile(filepath.Join(root, jf.rel))
		if err != nil {
			t.Fatal(err)
		}
		var got string
		var crashed jutil.Throwable
		jutil.Try(func() { got = trDumpFile(jf.rel, string(src), jf.sections, jf.full) },
			func(e jutil.Throwable) { crashed = e })
		if crashed != nil {
			got = "# go crash " + crashed.Error() + "\n"
		}
		if got == jf.text {
			continue
		}
		mismatches++
		if mismatches <= 20 {
			wantLines := strings.Split(jf.text, "\n")
			gotLines := strings.Split(got, "\n")
			section := ""
			for i := 0; i < len(wantLines) || i < len(gotLines); i++ {
				var w, g string
				if i < len(wantLines) {
					w = wantLines[i]
				}
				if i < len(gotLines) {
					g = gotLines[i]
				}
				if strings.HasPrefix(w, "## ") {
					section = w
				}
				if w != g {
					t.Errorf("%s: %s: line %d:\n  java: %s\n  go:   %s", jf.rel, section, i+1, w, g)
					break
				}
			}
		}
		if diffDir != "" {
			out := filepath.Join(diffDir, jf.rel+".go.txt")
			if err := os.MkdirAll(filepath.Dir(out), 0o755); err == nil {
				_ = os.WriteFile(out, []byte(got), 0o644)
			}
		}
	}
	if mismatches > 0 {
		t.Errorf("%d of %d files differ", mismatches, len(files))
	}
}

func TestTranspileRuntimeOracleGolden(t *testing.T) {
	defer SkipIfPendingPort(t) // the Compiler reaches passes not ported yet
	dir := filepath.Join("testdata", "transpile_runtime")
	for _, dump := range []string{"golden.txt.gz", "golden-inject.txt.gz"} {
		files := trReadJavaDump(t, filepath.Join(dir, dump))
		if len(files) < 500 {
			t.Fatalf("only %d files in %s", len(files), dump)
		}
		trCompare(t, filepath.Join(dir, "src"), files, os.Getenv("TRANSPILE_RUNTIME_DIFF_DIR"))
	}
}

func TestTranspileRuntimeOracleCorpus(t *testing.T) {
	defer SkipIfPendingPort(t) // the Compiler reaches passes not ported yet
	dumps := os.Getenv("TRANSPILE_RUNTIME_ORACLE")
	if dumps == "" {
		t.Skip("TRANSPILE_RUNTIME_ORACLE not set")
	}
	root := os.Getenv("TRANSPILE_RUNTIME_ROOT")
	if root == "" {
		root = "/home/user/jsgo/.cache/corpus"
	}
	for _, dump := range strings.Split(dumps, ",") {
		trCompare(t, root, trReadJavaDump(t, dump), os.Getenv("TRANSPILE_RUNTIME_DIFF_DIR"))
	}
}
