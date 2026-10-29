package jscomp

// Differential test against tools/oracle/java/ScopeDump.java: the same sweep
// of the NodeTraversal, SyntacticScopeCreator / Scope / Var,
// MemoizedScopeCreator and ReferenceCollector / ReferenceCollection /
// Reference API over parsed files, printed in the same format (see the
// header of ScopeDump.java), compared file by file.
//
//   - TestScopeDumpGolden replays the checked-in dump
//     testdata/scope/golden.txt.gz: the hand-written inputs in
//     testdata/scope/js, then the corpus files of
//     testdata/scope/golden-files.txt (every 150th file of corpus-files.txt
//     of at most 16 kB; read from SCOPE_CORPUS, default
//     /home/user/jsgo/.cache/corpus, and skipped if it is missing).
//     Regenerate with
//       CP=.cache/$CLOSURE_JAR:<build>
//       java -Xss1g -cp $CP com.google.javascript.jscomp.ScopeDump a.gz \
//         internal/jscomp/testdata/scope/js <the js names, sorted>
//       java -Xss1g -cp $CP com.google.javascript.jscomp.ScopeDump b.gz \
//         .cache/corpus @internal/jscomp/testdata/scope/golden-files.txt
//       cat a.gz b.gz > internal/jscomp/testdata/scope/golden.txt.gz
//   - TestScopeDumpOracle compares any other dump(s): SCOPE_ORACLE=<a.gz>[,<b.gz>...]
//     [SCOPE_ROOT=<root>] (default: the corpus), SCOPE_DIFF_DIR=<dir> receives
//     the Go dump of each differing file. The held-out corpus
//     (.cache/corpus-holdout, corpus-holdout-files.txt) is for a final check
//     only.
//
// The Java side parses with ParserRunner.parse and the parsing.Config of a
// Compiler whose options have language_in ECMASCRIPT_NEXT, sets the SCRIPT's
// InputId, puts it under a ROOT and registers a CompilerInput for it with the
// Compiler's private putCompilerInput; so does the Go side, with the test
// compiler.

import (
	"bufio"
	"compress/gzip"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// newScopeDumpCompiler is ScopeDump.newCompiler.
func newScopeDumpCompiler() *Compiler {
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	return newTestCompiler(options)
}

// scopeDumpParserConfig is Compiler.getParserConfig(ConfigContext.DEFAULT)
// for those options (see internal/parsing/parsingdump_test.go).
func scopeDumpParserConfig() *parsing.Config {
	return parsing.ParserRunner_CreateConfigWithModes(
		parsing.Config_LanguageMode_ES_NEXT, parsing.Config_JsDocParsing_TYPES_ONLY,
		parsing.Config_RunMode_STOP_AFTER_ERROR, nil, true, parsing.Config_StrictMode_STRICT)
}

type scopeDumpReporter struct{}

func (scopeDumpReporter) Warning(message string, sourceName string, line int, lineOffset int) {}
func (scopeDumpReporter) Error(message string, sourceName string, line int, lineOffset int)   {}

// scopeDumper accumulates one file's dump.
type scopeDumper struct {
	sb strings.Builder
}

// line is ScopeDump.line: every line escaped like ParsingDump.
func (d *scopeDumper) line(s string) {
	for _, c := range jutil.StringToUTF16(s) {
		switch {
		case c == '\\':
			d.sb.WriteString("\\\\")
		case c >= 0x20 && c < 0x7f:
			d.sb.WriteByte(byte(c))
		default:
			const hex = "0123456789abcdef"
			d.sb.WriteString("\\u")
			d.sb.WriteByte(hex[c>>12&0xf])
			d.sb.WriteByte(hex[c>>8&0xf])
			d.sb.WriteByte(hex[c>>4&0xf])
			d.sb.WriteByte(hex[c&0xf])
		}
	}
	d.sb.WriteByte('\n')
}

// scopeDumpNode is ScopeDump.node.
func scopeDumpNode(n rhino.Node) string {
	if n.IsNil() {
		return "null"
	}
	return n.GetToken().String() + "@" + strconv.Itoa(n.GetLineno()) + ":" + strconv.Itoa(n.GetCharno())
}

// scopeDumpJavaClass is Throwable.getClass().getName().
func scopeDumpJavaClass(t jutil.Throwable) string {
	full := t.Error()
	if msg := jutil.GetMessageOrNull(t); msg != nil {
		return strings.TrimSuffix(full, ": "+*msg)
	}
	return full
}

// scopeDumpExc is ScopeDump.exc.
func scopeDumpExc(e jutil.Throwable) string {
	var b strings.Builder
	b.WriteString(scopeDumpJavaClass(e))
	for c := e.GetCause(); c != nil && c != e; c = c.GetCause() {
		b.WriteByte('<')
		b.WriteString(scopeDumpJavaClass(c))
		if c.GetCause() == c {
			break
		}
	}
	return b.String()
}

// scopeDumpStr is ScopeDump.str.
func scopeDumpStr(s func() string) (result string) {
	jutil.Try(func() { result = s() }, func(e jutil.Throwable) { result = "!" + scopeDumpExc(e) })
	return result
}

// scopeDumpFlag is ScopeDump.flag.
func scopeDumpFlag(f *strings.Builder, c byte, b func() bool) {
	var set, failed bool
	jutil.Try(func() { set = b() }, func(e jutil.Throwable) { failed = true })
	if failed {
		f.WriteByte('!')
	} else if set {
		f.WriteByte(c)
	}
}

func scopeDumpBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// dumpOne is ScopeDump.dumpOne.
func (d *scopeDumper) dumpOne(name, src string) string {
	d.sb.Reset()
	compiler := newScopeDumpCompiler()
	config := scopeDumpParserConfig()
	file := SourceFile_FromCode(name, src)
	d.line("# file " + name)
	var result *parsing.ParserRunner_ParseResult
	var exception jutil.Throwable
	jutil.Try(func() {
		result = parsing.ParserRunner_Parse(file, src, config, scopeDumpReporter{})
	}, func(e jutil.IRuntimeException) { exception = e })
	if exception != nil {
		d.line("# exception " + scopeDumpExc(exception))
		return d.sb.String()
	}
	if result.Ast.IsNil() {
		d.line("# no ast")
		return d.sb.String()
	}
	script := result.Ast
	script.SetInputId(rhino.NewInputId(name))
	root := rhino.IR_Root(script)
	compiler.putCompilerInput(NewCompilerInput(file))

	d.section("# traversal", func() { d.traversal(compiler, root, false, true) })
	d.section("# obey", func() { d.traversal(compiler, root, true, false) })
	d.section("# scoperoots", func() { d.scopeRoots(compiler, root) })
	d.section("# atscope", func() { d.atScope(compiler, root) })
	d.section("# scopes", func() { d.scopes(compiler, root) })
	d.section("# references", func() { d.references(compiler, root) })
	d.section("# memoized", func() { d.memoized(compiler, root) })
	return d.sb.String()
}

// section is ScopeDump.section.
func (d *scopeDumper) section(title string, r func()) {
	d.line(title)
	jutil.Try(r, func(e jutil.Throwable) { d.line("# exception " + scopeDumpExc(e)) })
}

// scopeDumpRecorder is ScopeDump.Recorder.
type scopeDumpRecorder struct {
	d       *scopeDumper
	verbose bool
	scoped  bool
}

func (r *scopeDumpRecorder) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if r.verbose {
		r.d.line("S " + scopeDumpNode(n) +
			" p=" + scopeDumpNode(parent) +
			" d=" + scopeDumpStr(func() string { return strconv.Itoa(t.GetScopeDepth()) }) +
			" r=" + scopeDumpNode(t.GetScopeRoot()) +
			" h=" + scopeDumpNode(t.GetClosestHoistScopeRoot()) +
			" f=" + scopeDumpNode(t.GetEnclosingFunction()) +
			" b=" + scopeDumpNode(t.GetClosestScopeRootNodeBindingThisOrSuper()) +
			" m=" + scopeDumpStr(func() string { return scopeDumpTraversalFlags(t) }))
	} else {
		r.d.line("S " + scopeDumpNode(n))
	}
	return true
}

func (r *scopeDumpRecorder) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if r.verbose {
		r.d.line("V " + scopeDumpNode(n) + " p=" + scopeDumpNode(parent))
	} else {
		r.d.line("V " + scopeDumpNode(n))
	}
}

func (r *scopeDumpRecorder) EnterScope(t *NodeTraversal) {
	if r.scoped {
		r.d.line("E " + scopeDumpNode(t.GetScopeRoot()) +
			" d=" + scopeDumpStr(func() string { return strconv.Itoa(t.GetScopeDepth()) }))
	}
}

func (r *scopeDumpRecorder) ExitScope(t *NodeTraversal) {
	if r.scoped {
		r.d.line("X " + scopeDumpNode(t.GetScopeRoot()) +
			" d=" + scopeDumpStr(func() string { return strconv.Itoa(t.GetScopeDepth()) }))
	}
}

// scopeDumpTraversalFlags is ScopeDump.traversalFlags.
func scopeDumpTraversalFlags(t *NodeTraversal) string {
	var f strings.Builder
	if t.InGlobalHoistScope() {
		f.WriteByte('G')
	}
	if t.InModuleHoistScope() {
		f.WriteByte('M')
	}
	if t.IsHoistScope() {
		f.WriteByte('H')
	}
	if t.InFunctionBlockScope() {
		f.WriteByte('K')
	}
	if t.InModuleScope() {
		f.WriteByte('O')
	}
	if t.InGlobalScope() {
		f.WriteByte('g')
	}
	return f.String()
}

func (d *scopeDumper) traversal(compiler AbstractCompiler, root rhino.Node, obey, verbose bool) {
	NodeTraversal_Builder_().
		SetCompiler(compiler).
		SetCallback(&scopeDumpRecorder{d: d, verbose: verbose, scoped: verbose}).
		SetObeyDestructuringAndDefaultValueExecutionOrder(obey).
		Traverse(root)
}

// scopeDumpFunctions is ScopeDump.functions.
func scopeDumpFunctions(root rhino.Node) []rhino.Node {
	var fns []rhino.Node
	var collect func(n rhino.Node)
	collect = func(n rhino.Node) {
		if n.IsFunction() {
			fns = append(fns, n)
		}
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			collect(c)
		}
	}
	collect(root)
	return fns
}

// scopeDumpSize is ScopeDump.size.
func scopeDumpSize(n rhino.Node) int {
	s := 1
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		s += scopeDumpSize(c)
	}
	return s
}

func (d *scopeDumper) scopeRoots(compiler AbstractCompiler, root rhino.Node) {
	for _, fn := range scopeDumpFunctions(root) {
		d.line("R " + scopeDumpNode(fn))
		jutil.Try(func() {
			NodeTraversal_TraverseScopeRoots(
				compiler, []rhino.Node{fn}, &scopeDumpRecorder{d: d, verbose: false, scoped: true}, false)
		}, func(e jutil.Throwable) { d.line("# exception " + scopeDumpExc(e)) })
	}
}

// scopeDumpScopeCollector is the AbstractScopedCallback of ScopeDump.atScope.
type scopeDumpScopeCollector struct {
	NodeTraversal_AbstractScopedCallbackBase
	fnScopes []*Scope
}

func (c *scopeDumpScopeCollector) EnterScope(t *NodeTraversal) {
	if t.GetScopeRoot().IsFunction() {
		c.fnScopes = append(c.fnScopes, t.GetScope())
	}
}

func (c *scopeDumpScopeCollector) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {}

func (d *scopeDumper) atScope(compiler AbstractCompiler, root rhino.Node) {
	creator := NewMemoizedScopeCreator(NewSyntacticScopeCreator(compiler))
	collector := &scopeDumpScopeCollector{}
	NodeTraversal_Builder_().
		SetCompiler(compiler).
		SetScopeCreator(creator).
		SetCallback(collector).
		Traverse(root)
	for _, s := range collector.fnScopes {
		if scopeDumpSize(s.GetRootNode()) > 64 {
			continue
		}
		d.line("A " + scopeDumpNode(s.GetRootNode()))
		jutil.Try(func() {
			NodeTraversal_Builder_().
				SetCompiler(compiler).
				SetCallback(&scopeDumpRecorder{d: d, verbose: false, scoped: true}).
				SetScopeCreator(creator).
				TraverseAtScope(s)
		}, func(e jutil.Throwable) { d.line("# exception " + scopeDumpExc(e)) })
	}
}

// scopeDumpRedeclarationHandler is the RedeclarationHandler lambda of
// ScopeDump.scopes.
type scopeDumpRedeclarationHandler struct{ d *scopeDumper }

func (h *scopeDumpRedeclarationHandler) OnRedeclaration(s *Scope, name string, n rhino.Node, input *CompilerInput) {
	inputName := "null"
	if input != nil {
		inputName = input.GetName()
	}
	h.d.line("redeclare " + scopeDumpNode(s.GetRootNode()) + " " + name + " " + scopeDumpNode(n) +
		" input=" + inputName)
}

// scopeDumpScopePrinter is the AbstractScopedCallback of ScopeDump.scopes.
type scopeDumpScopePrinter struct {
	NodeTraversal_AbstractScopedCallbackBase
	d *scopeDumper
}

func (p *scopeDumpScopePrinter) EnterScope(t *NodeTraversal) {
	p.d.printScope(t.GetScope())
}

func (p *scopeDumpScopePrinter) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {}

func (d *scopeDumper) scopes(compiler AbstractCompiler, root rhino.Node) {
	creator := NewSyntacticScopeCreatorWithTreatProvidesAsRedeclarations(
		compiler, &scopeDumpRedeclarationHandler{d: d}, true)
	NodeTraversal_Builder_().
		SetCompiler(compiler).
		SetScopeCreator(creator).
		SetCallback(&scopeDumpScopePrinter{d: d}).
		Traverse(root)
}

// scopeDumpScopeRoot is ScopeDump.scopeRoot.
func scopeDumpScopeRoot(s *Scope) string {
	if s == nil {
		return "null"
	}
	return scopeDumpNode(s.GetRootNode())
}

func (d *scopeDumper) printScope(s *Scope) {
	d.line("scope " + scopeDumpNode(s.GetRootNode()) +
		" d=" + strconv.Itoa(s.GetDepth()) +
		" p=" + scopeDumpScopeRoot(s.GetParent()) +
		" f=" + scopeDumpScopeFlags(s) +
		" h=" + scopeDumpStr(func() string { return scopeDumpScopeRoot(s.GetClosestHoistScope()) }) +
		" c=" + scopeDumpStr(func() string { return scopeDumpScopeRoot(s.GetClosestContainerScope()) }) +
		" t=" + scopeDumpStr(func() string { return scopeDumpScopeRoot(s.GetScopeOfThis()) }) +
		" cfg=" + scopeDumpStr(func() string { return scopeDumpScopeRoot(s.GetClosestCfgRootScope()) }) +
		" n=" + strconv.Itoa(s.GetVarCount()))
	implicit := "implicit"
	for _, name := range []string{"arguments", "this", "super", "exports"} {
		implicit += " " + name + "=" +
			scopeDumpStr(func() string {
				if s.HasOwnSlot(name) {
					return "Y"
				}
				return "N"
			}) +
			scopeDumpStr(func() string {
				if s.GetOwnSlot(name) != nil {
					return "V"
				}
				return "-"
			})
	}
	d.line(implicit)
	d.line("args=" + scopeDumpStr(func() string {
		v := s.GetArgumentsVar()
		if v == nil {
			return "null"
		}
		return scopeDumpNode(v.GetScopeRoot()) + "/" + strconv.Itoa(v.GetIndex())
	}))
	for v := range s.GetVarIterable() {
		d.printVar(v)
	}
}

// scopeDumpScopeFlags is ScopeDump.scopeFlags.
func scopeDumpScopeFlags(s *Scope) string {
	var f strings.Builder
	scopeDumpFlag(&f, 'B', s.IsBlockScope)
	scopeDumpFlag(&f, 'S', s.IsStaticBlockScope)
	scopeDumpFlag(&f, 'K', s.IsFunctionBlockScope)
	scopeDumpFlag(&f, 'F', s.IsFunctionScope)
	scopeDumpFlag(&f, 'M', s.IsModuleScope)
	scopeDumpFlag(&f, 'D', s.IsMemberFieldDefScope)
	scopeDumpFlag(&f, 'C', s.IsComputedFieldDefRhsScope)
	scopeDumpFlag(&f, 'T', s.IsCatchScope)
	scopeDumpFlag(&f, 'R', s.IsCfgRootScope)
	scopeDumpFlag(&f, 'H', s.IsHoistScope)
	scopeDumpFlag(&f, 'G', s.IsGlobal)
	scopeDumpFlag(&f, 'L', s.IsLocal)
	return f.String()
}

func (d *scopeDumper) printVar(v *Var) {
	var f strings.Builder
	scopeDumpFlag(&f, 'v', v.IsVar)
	scopeDumpFlag(&f, 'l', v.IsLet)
	scopeDumpFlag(&f, 'c', v.IsConst)
	scopeDumpFlag(&f, 'k', v.IsClass)
	scopeDumpFlag(&f, 'p', v.IsParam)
	scopeDumpFlag(&f, 'h', v.IsCatch)
	scopeDumpFlag(&f, 'i', v.IsImport)
	scopeDumpFlag(&f, 'b', v.IsBleedingFunction)
	scopeDumpFlag(&f, 'a', v.IsArguments)
	scopeDumpFlag(&f, 't', v.IsThis)
	scopeDumpFlag(&f, 'e', v.IsGoogModuleExports)
	scopeDumpFlag(&f, 'n', v.IsImplicitGoogNamespace)
	scopeDumpFlag(&f, 'C', v.IsDeclaredOrInferredConst)
	scopeDumpFlag(&f, 'D', v.IsDefine)
	scopeDumpFlag(&f, 'x', v.IsExtern)
	scopeDumpFlag(&f, 'g', v.IsGlobal)
	scopeDumpFlag(&f, 'd', v.IsDefaultParam)
	doc := "0"
	if v.GetJSDocInfo() != nil {
		doc = "1"
	}
	sym := "0"
	if s, ok := v.GetSymbol().(*Var); ok && s == v {
		sym = "1"
	}
	decl := "0"
	if s, ok := v.GetDeclaration().(*Var); ok && s == v {
		decl = "1"
	}
	d.line("var " + v.GetName() +
		" n=" + scopeDumpNode(v.GetNameNode()) +
		" pn=" + scopeDumpNode(v.GetParentNode()) +
		" i=" + strconv.Itoa(v.GetIndex()) +
		" dt=" + scopeDumpStr(func() string {
		if t, ok := v.DeclarationToken(); ok {
			return t.String()
		}
		return "null"
	}) +
		" f=" + f.String() +
		" init=" + scopeDumpStr(func() string { return scopeDumpNode(v.GetInitialValue()) }) +
		" doc=" + doc +
		" sym=" + sym +
		" decl=" + decl +
		" str=" + scopeDumpStr(func() string { return v.GetImplicitGoogNamespaceStrength().String() }))
}

// scopeDumpBehavior is the Behavior lambda of ScopeDump.references.
type scopeDumpBehavior struct{ d *scopeDumper }

func (b *scopeDumpBehavior) AfterExitScope(t *NodeTraversal, referenceMap ReferenceMap) {
	var l strings.Builder
	l.WriteString("exit ")
	l.WriteString(scopeDumpNode(t.GetScopeRoot()))
	for v := range t.GetScope().GetVarIterable() {
		refs := referenceMap.GetReferences(v)
		l.WriteByte(' ')
		l.WriteString(v.GetName())
		l.WriteByte('=')
		if refs == nil {
			l.WriteString("-")
		} else {
			l.WriteString(strconv.Itoa(len(refs.references)))
		}
	}
	b.d.line(l.String())
}

func (d *scopeDumper) references(compiler AbstractCompiler, root rhino.Node) {
	collector := NewReferenceCollector(compiler, &scopeDumpBehavior{d: d}, NewSyntacticScopeCreator(compiler))
	collector.ProcessRoot(root)
	for v := range collector.GetAllSymbols() {
		refs := collector.GetReferences(v)
		d.line("sym " + v.GetName() +
			" s=" + scopeDumpNode(v.GetScopeRoot()) +
			" n=" + scopeDumpNode(v.GetNameNode()) +
			" g=" + scopeDumpNode(collector.GetScope(v).GetRootNode()))
		d.line("coll wd=" + scopeDumpStr(func() string { return scopeDumpBool(refs.IsWellDefined()) }) +
			" esc=" + scopeDumpStr(func() string { return scopeDumpBool(refs.IsEscaped()) }) +
			" once=" + scopeDumpStr(func() string { return scopeDumpBool(refs.IsAssignedOnceInLifetime()) }) +
			" never=" + scopeDumpStr(func() string { return scopeDumpBool(refs.IsNeverAssigned()) }) +
			" first=" + scopeDumpStr(func() string { return scopeDumpBool(refs.FirstReferenceIsAssigningDeclaration()) }) +
			" weak=" + scopeDumpStr(func() string { return scopeDumpBool(refs.IsReferencedWeakly()) }) +
			" init=" + scopeDumpStr(func() string {
			return strconv.Itoa(slices.Index(refs.references, refs.GetInitializingReference()))
		}) +
			" initc=" + scopeDumpStr(func() string {
			return strconv.Itoa(slices.Index(refs.references, refs.GetInitializingReferenceForConstants()))
		}) +
			" one=" + scopeDumpStr(func() string {
			return strconv.Itoa(slices.Index(refs.references, refs.GetOneAndOnlyAssignment()))
		}))
		for ref := range refs.All() {
			d.printReference(ref)
		}
	}
}

func (d *scopeDumper) printReference(ref *Reference) {
	var f strings.Builder
	scopeDumpFlag(&f, 'D', ref.IsDeclaration)
	scopeDumpFlag(&f, 'L', ref.IsLvalue)
	scopeDumpFlag(&f, 'I', ref.IsInitializingDeclaration)
	scopeDumpFlag(&f, 'v', ref.IsVarDeclaration)
	scopeDumpFlag(&f, 'l', ref.IsLetDeclaration)
	scopeDumpFlag(&f, 'c', ref.IsConstDeclaration)
	scopeDumpFlag(&f, 'h', ref.IsHoistedFunction)
	scopeDumpFlag(&f, 's', ref.IsSimpleAssignmentToName)
	scopeDumpFlag(&f, 'o', ref.IsDotPropertyAccess)
	scopeDumpFlag(&f, 'O', ref.IsAssignedToObjectDestructuringPattern)
	bb := ref.GetBasicBlock()
	bbp := "null"
	if bb.GetParent() != nil {
		bbp = scopeDumpNode(bb.GetParent().GetRoot())
	}
	bbf := ""
	if bb.IsFunction() {
		bbf += "F"
	}
	if bb.IsLoop() {
		bbf += "L"
	}
	d.line("ref " + scopeDumpNode(ref.GetNode()) +
		" p=" + scopeDumpNode(ref.GetParent()) +
		" s=" + scopeDumpNode(ref.GetScope().GetRootNode()) +
		" bb=" + scopeDumpNode(bb.GetRoot()) +
		" bbp=" + bbp +
		" bbf=" + bbf +
		" f=" + f.String() +
		" av=" + scopeDumpStr(func() string { return scopeDumpNode(ref.GetAssignedValue()) }) +
		" sym=" + scopeDumpStr(func() string {
		s, _ := ref.GetSymbol().(*Var)
		if s == nil {
			return "null"
		}
		return scopeDumpNode(s.GetNameNode()) + "/" + scopeDumpNode(s.GetScopeRoot())
	}))
}

func (d *scopeDumper) memoized(compiler AbstractCompiler, root rhino.Node) {
	var creator ScopeCreator = NewMemoizedScopeCreator(NewSyntacticScopeCreator(compiler))
	a := creator.CreateScope(root, nil)
	b := creator.CreateScope(root, nil)
	d.line("same=" + scopeDumpBool(a == b) + " vars=" + strconv.Itoa(a.(*Scope).GetVarCount()))
	d.line("reparent=" + scopeDumpStr(func() string {
		creator.CreateScope(root, a)
		return "ok"
	}))
}

// scopeDumpReadFiles streams the "# file" sections of gzipped dumps to
// each(name, section), in order.
func scopeDumpReadFiles(t *testing.T, each func(name, section string), paths ...string) {
	t.Helper()
	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		zr, err := gzip.NewReader(f)
		if err != nil {
			t.Fatal(err)
		}
		sc := bufio.NewScanner(zr)
		sc.Buffer(make([]byte, 1<<20), 1<<28)
		var cur string
		var sb strings.Builder
		flush := func() {
			if cur != "" {
				each(cur, sb.String())
			}
			sb.Reset()
		}
		for sc.Scan() {
			l := sc.Text()
			if strings.HasPrefix(l, "# file ") {
				flush()
				cur = l[len("# file "):]
			}
			sb.WriteString(l)
			sb.WriteByte('\n')
		}
		if err := sc.Err(); err != nil {
			t.Fatal(err)
		}
		flush()
		f.Close()
	}
}

func scopeDumpFirstDiff(want, got string) string {
	wl, gl := strings.Split(want, "\n"), strings.Split(got, "\n")
	for i := 0; i < len(wl) || i < len(gl); i++ {
		var w, g string
		if i < len(wl) {
			w = wl[i]
		}
		if i < len(gl) {
			g = gl[i]
		}
		if w != g {
			return "line " + strconv.Itoa(i+1) + ":\n  java: " + w + "\n  go:   " + g
		}
	}
	return ""
}

// scopeDumpCompare dumps every file of the Java dumps for which rootFor
// returns a root and compares (streaming: only differences are kept).
func scopeDumpCompare(t *testing.T, rootFor func(name string) string, paths ...string) {
	diffDir := os.Getenv("SCOPE_DIFF_DIR")
	type job struct{ name, want string }
	type result struct{ name, diff, got string }
	jobs := make(chan job, runtime.NumCPU())
	results := make(chan result, runtime.NumCPU())
	var wg sync.WaitGroup
	for w := 0; w < runtime.NumCPU(); w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d := &scopeDumper{}
			for j := range jobs {
				src, err := os.ReadFile(filepath.Join(rootFor(j.name), j.name))
				if err != nil {
					results <- result{name: j.name, diff: err.Error()}
					continue
				}
				got := d.dumpOne(j.name, string(src))
				if got != j.want {
					results <- result{name: j.name, diff: scopeDumpFirstDiff(j.want, got), got: got}
				} else {
					results <- result{name: j.name}
				}
			}
		}()
	}
	compared, failures := 0, 0
	done := make(chan struct{})
	go func() {
		for r := range results {
			compared++
			if r.diff == "" {
				continue
			}
			failures++
			if failures <= 20 {
				t.Errorf("%s: differs from Java at %s", r.name, r.diff)
			}
			if diffDir != "" && r.got != "" {
				out := filepath.Join(diffDir, r.name+".go.scope")
				os.MkdirAll(filepath.Dir(out), 0o755)
				os.WriteFile(out, []byte(r.got), 0o644)
			}
		}
		close(done)
	}()
	scopeDumpReadFiles(t, func(name, section string) {
		if rootFor(name) != "" {
			jobs <- job{name, section}
		}
	}, paths...)
	close(jobs)
	wg.Wait()
	close(results)
	<-done
	if compared == 0 {
		t.Fatalf("no files compared in %v", paths)
	}
	t.Logf("compared %d files: %d differ", compared, failures)
	if failures > 0 {
		t.Errorf("%d of %d files differ", failures, compared)
	}
}

func scopeDumpCorpus(t *testing.T) string {
	if root := os.Getenv("SCOPE_CORPUS"); root != "" {
		return root
	}
	return oracleTestCorpusRoot(t)
}

func TestScopeDumpGolden(t *testing.T) {
	js, err := os.ReadDir("testdata/scope/js")
	if err != nil {
		t.Fatal(err)
	}
	local := map[string]bool{}
	for _, e := range js {
		local[e.Name()] = true
	}
	corpus := scopeDumpCorpus(t)
	if _, err := os.Stat(corpus); err != nil {
		t.Logf("corpus not available (%v): comparing the hand-written inputs only", err)
		corpus = ""
	}
	scopeDumpCompare(t, func(name string) string {
		if local[name] {
			return "testdata/scope/js"
		}
		return corpus
	}, "testdata/scope/golden.txt.gz")
}

func TestScopeDumpOracle(t *testing.T) {
	oracle := os.Getenv("SCOPE_ORACLE")
	if oracle == "" {
		t.Skip("set SCOPE_ORACLE=<dump.gz>[,<dump.gz>...] (tools/oracle/java/ScopeDump.java) to run")
	}
	root := os.Getenv("SCOPE_ROOT")
	if root == "" {
		root = scopeDumpCorpus(t)
	}
	paths := strings.Split(oracle, ",")
	sort.Strings(paths)
	scopeDumpCompare(t, func(string) string { return root }, paths...)
}
