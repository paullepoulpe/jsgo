package jscomp

// Go-only: differential test — the Go side of tools/oracle/java/CoreG2Dump.java
// (core group g2: AstAnalyzer, ExpressionDecomposer, OptionalChainRewriter,
// MakeDeclaredNamesUnique, AstFactory). coreG2Dump reproduces, byte for byte,
// what the Java tool prints (see its header for the format of each mode).
//
//   - TestCoreG2OracleGolden compares against the checked-in Java dumps in
//     testdata/core_g2/: <mode>.txt.gz for the sources under
//     testdata/core_g2/src (listed in testdata/core_g2/files.txt: every
//     JavaScript string of AstAnalyzerTest, ExpressionDecomposerTest,
//     MakeDeclaredNamesUniqueTest, RewriteOptionalChainingOperatorTest and
//     AstFactoryTest), fixtures-<mode>.txt.gz for the repository's
//     testdata/fixtures, and factory.txt.gz. Regenerate with
//       java -Xss1g -cp closure-compiler.jar:<build> \
//         com.google.javascript.jscomp.CoreG2Dump <mode> \
//         internal/jscomp/testdata/core_g2/<mode>.txt.gz \
//         internal/jscomp/testdata/core_g2/src @internal/jscomp/testdata/core_g2/files.txt
//     (testdata/fixtures with @.cache/fixtures-files.txt for fixtures-<mode>,
//     and `... CoreG2Dump factory internal/jscomp/testdata/core_g2/factory.txt.gz`).
//   - TestCoreG2OracleCorpus compares a whole corpus when CORE_G2_ROOT (the
//     root of the sources), CORE_G2_FILES (the list of relative paths) and
//     CORE_G2_ORACLE_DIR (holding corpus-<mode>.txt.gz written by the Java
//     tool for that list) are set; CORE_G2_MODES (comma-separated) restricts
//     the modes, CORE_G2_EVERY=N compares only every Nth file of the dump
//     (a deterministic sample), and CORE_G2_DIFF_DIR, if set, receives the
//     Go and Java dumps of every mismatching file. The Java dumps of /home/user/jsgo/.cache/corpus are
//     produced with the same command on @.cache/corpus-files.txt. The
//     held-out corpus (.cache/corpus-holdout, corpus-holdout-files.txt) is
//     for a final check only.

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ---------------------------------------------------------------------------
// Shared helpers

// coreG2Escape is CoreG2Dump.escape (ParsingDump.escape).
func coreG2Escape(s string) string {
	var b strings.Builder
	for _, c := range jutil.StringToUTF16(s) {
		switch {
		case c == '\\':
			b.WriteString("\\\\")
		case c >= 0x20 && c < 0x7f:
			b.WriteByte(byte(c))
		default:
			fmt.Fprintf(&b, "\\u%04x", c)
		}
	}
	return b.String()
}

type coreG2Out struct {
	sb strings.Builder
}

func (o *coreG2Out) line(s string) {
	o.sb.WriteString(coreG2Escape(s))
	o.sb.WriteByte('\n')
}

// coreG2Exc is CoreG2Dump.exc: the root cause's toString(), a
// NullPointerException's class name only.
func coreG2Exc(e jutil.Throwable) string {
	for e.GetCause() != nil {
		e = e.GetCause()
	}
	if _, ok := e.(*jutil.NullPointerException); ok {
		return "java.lang.NullPointerException"
	}
	return e.Error()
}

// coreG2Try runs body; a thrown Java exception is returned (nil if none).
func coreG2Try(body func()) (thrown jutil.Throwable) {
	jutil.Try(body, func(e jutil.IRuntimeException) { thrown = e })
	return thrown
}

// coreG2ParserConfig is the parsing.Config CompilerInput uses for
// --language_in=ECMASCRIPT_NEXT (CoreG2Dump.compilerParserConfig).
func coreG2ParserConfig() *parsing.Config {
	return parsing.ParserRunner_CreateConfigWithModes(
		parsing.Config_LanguageMode_ES_NEXT,
		parsing.Config_JsDocParsing_TYPES_ONLY,
		parsing.Config_RunMode_STOP_AFTER_ERROR,
		nil,
		true,
		parsing.Config_StrictMode_STRICT)
}

type coreG2ErrorFlag struct{ hadError bool }

func (r *coreG2ErrorFlag) Warning(message string, sourceName string, line int, lineOffset int) {}
func (r *coreG2ErrorFlag) Error(message string, sourceName string, line int, lineOffset int) {
	r.hadError = true
}

// coreG2Parse is CoreG2Dump.parse: nil when the parse reports an error or
// throws.
func coreG2Parse(config *parsing.Config, name string, src string) rhino.Node {
	reporter := &coreG2ErrorFlag{}
	var result *parsing.ParserRunner_ParseResult
	if coreG2Try(func() {
		result = parsing.ParserRunner_Parse(SourceFile_FromCode(name, src), src, config, reporter)
	}) != nil {
		return rhino.NilNode
	}
	if reporter.hadError || result.Ast.IsNil() {
		return rhino.NilNode
	}
	return result.Ast
}

// coreG2Preorder is CoreG2Dump.preorder.
func coreG2Preorder(root rhino.Node) []rhino.Node {
	var result []rhino.Node
	stack := []rhino.Node{root}
	for len(stack) != 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		result = append(result, n)
		for c := n.GetLastChild(); !c.IsNil(); c = c.GetPrevious() {
			stack = append(stack, c)
		}
	}
	return result
}

// coreG2NewCompiler is `new Compiler()` + initOptions(new CompilerOptions()).
func coreG2NewCompiler() AbstractCompiler {
	return newTestCompiler(NewCompilerOptions())
}

func (o *coreG2Out) printTree(n rhino.Node) {
	for _, l := range strings.Split(n.ToStringTree(), "\n") {
		if l != "" {
			o.line(l)
		}
	}
}

// coreG2Counter is CoreG2Dump.counter: "0", "1", ...
func coreG2Counter() func() string {
	next := 0
	return func() string {
		s := strconv.Itoa(next)
		next++
		return s
	}
}

// ---------------------------------------------------------------------------
// analysis

func coreG2Analyzers(compiler AbstractCompiler) []*AstAnalyzer {
	accessors := jutil.NewLinkedHashMap[string, AccessorSummary_PropertyAccessKind]()
	accessors.Put("getter", AccessorSummary_PropertyAccessKind_GETTER_ONLY)
	accessors.Put("setter", AccessorSummary_PropertyAccessKind_SETTER_ONLY)
	accessors.Put("length", AccessorSummary_PropertyAccessKind_GETTER_ONLY)
	accessors.Put("value", AccessorSummary_PropertyAccessKind_SETTER_ONLY)
	accessors.Put("x", AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER)
	accessors.Put("prototype", AccessorSummary_PropertyAccessKind_GETTER_AND_SETTER)
	accessors.Put("a", AccessorSummary_PropertyAccessKind_GETTER_ONLY)
	summary := AccessorSummary_Create(accessors)
	return []*AstAnalyzer{
		compiler.GetAstAnalyzer(),
		NewAstAnalyzer(
			AstAnalyzer_Options_Builder_().
				SetUseTypesForLocalOptimization(false).
				SetAssumeGettersArePure(false).
				SetHasRegexpGlobalReferences(false).
				SetAssumeKnownBuiltinsArePure(true).
				Build(),
			compiler.GetTypeRegistry(),
			summary),
		NewAstAnalyzer(
			AstAnalyzer_Options_Builder_().
				SetUseTypesForLocalOptimization(false).
				SetAssumeGettersArePure(true).
				SetHasRegexpGlobalReferences(true).
				SetAssumeKnownBuiltinsArePure(false).
				Build(),
			compiler.GetTypeRegistry(),
			nil),
		NewAstAnalyzer(
			AstAnalyzer_Options_Builder_().
				SetUseTypesForLocalOptimization(true).
				SetAssumeGettersArePure(true).
				SetHasRegexpGlobalReferences(false).
				SetAssumeKnownBuiltinsArePure(true).
				Build(),
			compiler.GetTypeRegistry(),
			nil),
	}
}

func coreG2Decomposers(compiler AbstractCompiler) []*ExpressionDecomposer {
	return []*ExpressionDecomposer{
		compiler.CreateDefaultExpressionDecomposer(),
		NewExpressionDecomposer(
			compiler,
			coreG2Counter(),
			map[string]bool{"foo": true, "goog": true, "fn": true, "require": true, "Object": true},
			Scope_CreateGlobalScope(rhino.NewNode(rhino.Token_SCRIPT)),
			true),
	}
}

var coreG2QueryNames = []string{"MHSE", "MEMS", "NTMHSE", "FCHSE", "CCHSE"}

var coreG2Queries = []func(*AstAnalyzer, rhino.Node) bool{
	(*AstAnalyzer).MayHaveSideEffects,
	(*AstAnalyzer).MayEffectMutableState,
	(*AstAnalyzer).NodeTypeMayHaveSideEffects,
	(*AstAnalyzer).FunctionCallHasSideEffects,
	(*AstAnalyzer).ConstructorCallHasSideEffects,
}

func (o *coreG2Out) analysis(script rhino.Node) {
	compiler := coreG2NewCompiler()
	analyzers := coreG2Analyzers(compiler)
	decomposers := coreG2Decomposers(compiler)
	for _, n := range coreG2Preorder(script) {
		var sb strings.Builder
		sb.WriteString(n.GetToken().String())
		var exceptions []string
		isCall := n.IsCall() || n.IsOptChainCall() || n.IsTaggedTemplateLit()
		for a, analyzer := range analyzers {
			sb.WriteByte(' ')
			for q, query := range coreG2Queries {
				if (q == 3 && !isCall) || (q == 4 && !n.IsNew()) {
					sb.WriteByte('-')
					continue
				}
				var result bool
				if e := coreG2Try(func() { result = query(analyzer, n) }); e != nil {
					sb.WriteByte('!')
					exceptions = append(exceptions, "! a"+strconv.Itoa(a)+" "+coreG2QueryNames[q]+" "+coreG2Exc(e))
				} else if result {
					sb.WriteByte('1')
				} else {
					sb.WriteByte('0')
				}
			}
		}
		for d, decomposer := range decomposers {
			sb.WriteByte(' ')
			if !rhino.IR_MayBeExpression(n) {
				sb.WriteByte('-')
				continue
			}
			var result ExpressionDecomposer_DecompositionType
			if e := coreG2Try(func() { result = decomposer.CanExposeExpression(n) }); e != nil {
				sb.WriteByte('!')
				exceptions = append(exceptions, "! d"+strconv.Itoa(d)+" canExposeExpression "+coreG2Exc(e))
				continue
			}
			switch result {
			case ExpressionDecomposer_DecompositionType_UNDECOMPOSABLE:
				sb.WriteByte('U')
			case ExpressionDecomposer_DecompositionType_MOVABLE:
				sb.WriteByte('M')
			case ExpressionDecomposer_DecompositionType_DECOMPOSABLE:
				sb.WriteByte('D')
			}
		}
		o.line(sb.String())
		for _, e := range exceptions {
			o.line(e)
		}
	}
}

// ---------------------------------------------------------------------------
// mutate

const coreG2RegionLines = 60

// printRegion is CoreG2Dump.printRegion.
func (o *coreG2Out) printRegion(block rhino.Node, from int, count int) {
	var lines []string
	c := block.GetChildAtIndex(from)
	for i := 0; i < count && !c.IsNil(); i, c = i+1, c.GetNext() {
		for _, l := range strings.Split(c.ToStringTree(), "\n") {
			if l != "" {
				lines = append(lines, coreG2Escape(l))
			}
		}
	}
	crc := crc32.NewIEEE()
	for i, l := range lines {
		if i < coreG2RegionLines {
			o.sb.WriteString(l)
			o.sb.WriteByte('\n')
		}
		crc.Write([]byte(l + "\n"))
	}
	if len(lines) > coreG2RegionLines {
		o.line("... " + strconv.Itoa(len(lines)) + " lines crc32 " + strconv.FormatUint(uint64(crc.Sum32()), 16))
	}
}

func (o *coreG2Out) mutate(script rhino.Node) {
	nodes := coreG2Preorder(script)
	// ExpressionDecomposer.maybeExposeExpression / moveExpression
	{
		compiler := coreG2NewCompiler()
		probe := compiler.CreateDefaultExpressionDecomposer()
		exposed := 0
		moved := 0
		for i := 0; i < len(nodes) && (exposed < 8 || moved < 4); i++ {
			n := nodes[i]
			if !rhino.IR_MayBeExpression(n) {
				continue
			}
			var type_ ExpressionDecomposer_DecompositionType
			if coreG2Try(func() { type_ = probe.CanExposeExpression(n) }) != nil {
				continue
			}
			expose := type_ == ExpressionDecomposer_DecompositionType_DECOMPOSABLE && exposed < 8
			move := type_ == ExpressionDecomposer_DecompositionType_MOVABLE && n.IsCall() && moved < 4
			if !expose && !move {
				continue
			}
			if expose {
				exposed++
			} else {
				moved++
			}
			op := "## move "
			if expose {
				op = "## expose "
			}
			o.line(op + strconv.Itoa(i) + " " + n.GetToken().String())
			clone := script.CloneTree()
			target := coreG2Preorder(clone)[i]
			decomposer := compiler.CreateExpressionDecomposer(
				coreG2Counter(), nil, Scope_CreateGlobalScope(rhino.NewNode(rhino.Token_SCRIPT)))
			var block rhino.Node
			var from, before int
			if e := coreG2Try(func() {
				injectionPoint := ExpressionDecomposer_FindInjectionPoint(target)
				block = injectionPoint.GetParent()
				from = block.GetIndexOfChild(injectionPoint)
				before = block.GetChildCount()
				if expose {
					decomposer.MaybeExposeExpression(target)
				} else {
					decomposer.MoveExpression(target)
				}
			}); e != nil {
				o.line("! " + coreG2Exc(e))
				continue
			}
			o.printRegion(block, from, block.GetChildCount()-before+1)
		}
	}
	// OptionalChainRewriter
	{
		compiler := coreG2NewCompiler()
		rewritten := 0
		for i := 0; i < len(nodes) && rewritten < 8; i++ {
			n := nodes[i]
			if !NodeUtil_IsEndOfFullOptChain(n) {
				continue
			}
			rewritten++
			o.line("## optchain " + strconv.Itoa(i) + " " + n.GetToken().String())
			clone := script.CloneTree()
			target := coreG2Preorder(clone)[i]
			next := 0
			tmpVarNameCreator := func() string {
				s := "tmp$" + strconv.Itoa(next)
				next++
				return s
			}
			var scope *Scope
			if rewritten%2 == 0 {
				scope = Scope_CreateGlobalScope(rhino.NewNode(rhino.Token_SCRIPT))
			}
			var block rhino.Node
			var from, before int
			if e := coreG2Try(func() {
				rewriter := NewOptionalChainRewriter(compiler, compiler.CreateAstFactory(), tmpVarNameCreator, scope, target)
				block = rewriter.insertionPoint.GetParent()
				from = block.GetIndexOfChild(rewriter.insertionPoint)
				before = block.GetChildCount()
				rewriter.Rewrite()
			}); e != nil {
				o.line("! " + coreG2Exc(e))
				continue
			}
			o.printRegion(block, from, block.GetChildCount()-before+1)
		}
	}
}

// ---------------------------------------------------------------------------
// mdnu

// printRenames is CoreG2Dump.printRenames.
func (o *coreG2Out) printRenames(orig []rhino.Node, tree rhino.Node) {
	now := coreG2Preorder(tree)
	for i, n := range orig {
		if n.IsName() || n.IsImportStar() {
			before := n.GetString()
			after := now[i].GetString()
			if before != after {
				o.line(strconv.Itoa(i) + " " + before + " " + after)
			}
		}
	}
}

func (o *coreG2Out) mdnu(script rhino.Node) {
	orig := coreG2Preorder(script)
	var convention CodingConvention = CodingConventions_GetDefault()
	// ContextualRenamer, then the inverter, then the assertOnChange check.
	{
		compiler := coreG2NewCompiler()
		clone := script.CloneTree()
		o.line("## contextual")
		if e := coreG2Try(func() {
			NodeTraversal_Traverse(compiler, clone, NewMakeDeclaredNamesUnique())
			o.printRenames(orig, clone)
			o.line("## inverted")
			MakeDeclaredNamesUnique_GetContextualRenameInverter(compiler).Process(rhino.NilNode, clone)
			o.printRenames(orig, clone)
			o.line("## assert")
			again := clone.CloneTree()
			NodeTraversal_Traverse(compiler, again, NewMakeDeclaredNamesUniqueWith(nil, true, true))
			o.line("ok")
		}); e != nil {
			o.line("! " + coreG2Exc(e))
		}
	}
	// InlineRenamer on the first three functions.
	{
		compiler := coreG2NewCompiler()
		done := 0
		for i := 0; i < len(orig) && done < 3; i++ {
			if !orig[i].IsFunction() {
				continue
			}
			done++
			o.line("## inline " + strconv.Itoa(i))
			// As FunctionToBlockMutator.makeLocalNamesUnique does, on a detached clone.
			fn := orig[i].CloneTree()
			if e := coreG2Try(func() {
				NodeTraversal_TraverseScopeRoots(
					compiler,
					[]rhino.Node{fn},
					NewMakeDeclaredNamesUniqueWith(
						NewMakeDeclaredNamesUnique_InlineRenamer(
							convention, coreG2Counter(), "inline_", done%2 == 0, true, nil),
						false, false),
					true)
				o.printRenames(coreG2Preorder(orig[i]), fn)
			}); e != nil {
				o.line("! " + coreG2Exc(e))
			}
		}
	}
	// BoilerplateRenamer
	{
		compiler := coreG2NewCompiler()
		clone := script.CloneTree()
		o.line("## boilerplate")
		if e := coreG2Try(func() {
			NodeTraversal_Traverse(
				compiler,
				clone,
				NewMakeDeclaredNamesUniqueWith(
					NewMakeDeclaredNamesUnique_BoilerplateRenamer(convention, coreG2Counter(), "boiler_"),
					false, false))
			o.printRenames(orig, clone)
		}); e != nil {
			o.line("! " + coreG2Exc(e))
		}
	}
	// TargettedRenamer
	{
		compiler := coreG2NewCompiler()
		clone := script.CloneTree()
		o.line("## targetted")
		if e := coreG2Try(func() {
			NodeTraversal_Traverse(
				compiler,
				clone,
				NewMakeDeclaredNamesUniqueWith(
					NewMakeDeclaredNamesUnique_TargettedRenamer(
						NewMakeDeclaredNamesUnique_ContextualRenamer(),
						jutil.NewLinkedHashSet("i", "e", "a", "x", "n", "t")),
					true, false))
			o.printRenames(orig, clone)
		}); e != nil {
			o.line("! " + coreG2Exc(e))
		}
	}
}

// ---------------------------------------------------------------------------
// Driver

// coreG2Dump is CoreG2Dump.main for a corpus mode.
func coreG2Dump(mode string, root string, files []string) (map[string]string, error) {
	config := coreG2ParserConfig()
	perFile := map[string]string{}
	header := &coreG2Out{}
	header.line("# config languageMode=" + config.LanguageMode().String() +
		" jsDocParsingMode=" + config.JsDocParsingMode().String() +
		" runMode=" + config.RunMode().String() +
		" strictMode=" + config.StrictMode().String())
	perFile[""] = header.sb.String()
	for _, rel := range files {
		b, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			return nil, err
		}
		o := &coreG2Out{}
		o.line("# file " + rel)
		script := coreG2Parse(config, rel, coreG2DecodeUTF8(b))
		if script.IsNil() {
			o.line("# parse-error")
		} else {
			switch mode {
			case "analysis":
				o.analysis(script)
			case "mutate":
				o.mutate(script)
			case "mdnu":
				o.mdnu(script)
			default:
				return nil, fmt.Errorf("unknown mode %q", mode)
			}
		}
		perFile[rel] = o.sb.String()
	}
	return perFile, nil
}

// coreG2DecodeUTF8 is Files.readString(path, UTF_8): malformed input is
// replaced by U+FFFD as the JDK decoder does.
func coreG2DecodeUTF8(b []byte) string {
	return strings.ToValidUTF8(string(b), "\uFFFD")
}

// coreG2ReadJavaDump splits a Java dump into its header ("") and per-file
// sections keyed by relative path.
func coreG2ReadJavaDump(path string) (map[string]string, []string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		return nil, nil, err
	}
	sections := map[string]string{}
	var order []string
	current := ""
	var sb strings.Builder
	flush := func() {
		sections[current] = sb.String()
		sb.Reset()
	}
	sc := bufio.NewScanner(zr)
	sc.Buffer(make([]byte, 1<<20), 1<<30)
	for sc.Scan() {
		l := sc.Text()
		if strings.HasPrefix(l, "# file ") {
			flush()
			current = strings.TrimPrefix(l, "# file ")
			order = append(order, current)
		}
		sb.WriteString(l)
		sb.WriteByte('\n')
	}
	if err := sc.Err(); err != nil {
		return nil, nil, err
	}
	flush()
	return sections, order, nil
}

func coreG2Compare(t *testing.T, mode string, root string, javaDump string, diffDir string) {
	want, files, err := coreG2ReadJavaDump(javaDump)
	if err != nil {
		t.Fatal(err)
	}
	if every, _ := strconv.Atoi(os.Getenv("CORE_G2_EVERY")); every > 1 {
		var sample []string
		for i := 0; i < len(files); i += every {
			sample = append(sample, files[i])
		}
		files = sample
	}
	got, err := coreG2Dump(mode, root, files)
	if err != nil {
		t.Fatal(err)
	}
	if got[""] != want[""] {
		t.Errorf("%s: header: got %q, want %q", mode, got[""], want[""])
	}
	mismatches := 0
	for _, rel := range files {
		if got[rel] == want[rel] {
			continue
		}
		mismatches++
		if mismatches <= 10 {
			t.Errorf("%s: %s: %s", mode, rel, coreG2FirstDiff(got[rel], want[rel]))
		}
		if diffDir != "" {
			p := filepath.Join(diffDir, mode, rel+".go.txt")
			_ = os.MkdirAll(filepath.Dir(p), 0o755)
			_ = os.WriteFile(p, []byte(got[rel]), 0o644)
			_ = os.WriteFile(filepath.Join(diffDir, mode, rel+".java.txt"), []byte(want[rel]), 0o644)
		}
	}
	if mismatches != 0 {
		t.Errorf("%s: %d of %d files differ", mode, mismatches, len(files))
	}
}

// coreG2FirstDiff describes the first differing line.
func coreG2FirstDiff(got, want string) string {
	g := strings.Split(got, "\n")
	w := strings.Split(want, "\n")
	for i := 0; i < len(g) || i < len(w); i++ {
		var gl, wl string
		if i < len(g) {
			gl = g[i]
		}
		if i < len(w) {
			wl = w[i]
		}
		if gl != wl {
			return fmt.Sprintf("line %d:\n  go:   %s\n  java: %s", i+1, gl, wl)
		}
	}
	return "identical?"
}

func TestCoreG2OracleGolden(t *testing.T) {
	dir := filepath.Join("testdata", "core_g2")
	for _, mode := range []string{"analysis", "mutate", "mdnu"} {
		t.Run(mode, func(t *testing.T) {
			coreG2Compare(t, mode, filepath.Join(dir, "src"), filepath.Join(dir, mode+".txt.gz"), "")
		})
		t.Run("fixtures-"+mode, func(t *testing.T) {
			coreG2Compare(t, mode, filepath.Join("..", "..", "testdata", "fixtures"),
				filepath.Join(dir, "fixtures-"+mode+".txt.gz"), "")
		})
	}
	t.Run("factory", func(t *testing.T) {
		want, _, err := coreG2ReadJavaDump(filepath.Join(dir, "factory.txt.gz"))
		if err != nil {
			t.Fatal(err)
		}
		if got := coreG2FactorySweep(); got != want[""] {
			t.Errorf("factory: %s", coreG2FirstDiff(got, want[""]))
		}
	})
}

func TestCoreG2OracleCorpus(t *testing.T) {
	root, list, oracleDir := os.Getenv("CORE_G2_ROOT"), os.Getenv("CORE_G2_FILES"), os.Getenv("CORE_G2_ORACLE_DIR")
	if root == "" || list == "" || oracleDir == "" {
		t.Skip("set CORE_G2_ROOT, CORE_G2_FILES and CORE_G2_ORACLE_DIR")
	}
	modes := []string{"analysis", "mutate", "mdnu"}
	if m := os.Getenv("CORE_G2_MODES"); m != "" {
		modes = strings.Split(m, ",")
	}
	for _, mode := range modes {
		t.Run(mode, func(t *testing.T) {
			coreG2Compare(t, mode, root, filepath.Join(oracleDir, "corpus-"+mode+".txt.gz"), os.Getenv("CORE_G2_DIFF_DIR"))
		})
	}
}

// ---------------------------------------------------------------------------
// factory

// coreG2TestScope is CoreG2Dump.TestScope: a StaticScope whose slots are
// declared by NAME nodes.
type coreG2TestScope struct {
	declarations *jutil.LinkedHashMap[string, rhino.Node]
}

func newCoreG2TestScope() *coreG2TestScope {
	return &coreG2TestScope{declarations: jutil.NewLinkedHashMap[string, rhino.Node]()}
}

func (s *coreG2TestScope) declare(name string, declaration rhino.Node) *coreG2TestScope {
	s.declarations.Put(name, declaration)
	return s
}

func (s *coreG2TestScope) GetRootNode() rhino.Node              { return rhino.NilNode }
func (s *coreG2TestScope) GetParentScope() rhino.StaticScope    { return nil }
func (s *coreG2TestScope) GetSlot(name string) rhino.StaticSlot { return s.GetOwnSlot(name) }
func (s *coreG2TestScope) GetOwnSlot(name string) rhino.StaticSlot {
	declaration := s.declarations.Get(name)
	if declaration.IsNil() {
		return nil
	}
	return &coreG2TestSlot{name: name, node: declaration, scope: s}
}
func (s *coreG2TestScope) GetTopmostScopeOfEventualDeclaration(name string) rhino.StaticScope {
	return rhino.StaticScope_GetTopmostScopeOfEventualDeclaration(s, name)
}
func (s *coreG2TestScope) String() string {
	return "TestScope" + jutil.CollectionToString(s.declarations.KeySlice())
}

// coreG2TestSlot is CoreG2Dump.TestSlot.
type coreG2TestSlot struct {
	name  string
	node  rhino.Node
	scope rhino.StaticScope
}

func (s *coreG2TestSlot) GetName() string                       { return s.name }
func (s *coreG2TestSlot) GetDeclaration() rhino.StaticRef       { return s }
func (s *coreG2TestSlot) GetJSDocInfo() *rhino.JSDocInfo        { return nil }
func (s *coreG2TestSlot) GetScope() rhino.StaticScope           { return s.scope }
func (s *coreG2TestSlot) GetSymbol() rhino.StaticSlot           { return s }
func (s *coreG2TestSlot) GetNode() rhino.Node                   { return s.node }
func (s *coreG2TestSlot) GetSourceFile() rhino.StaticSourceFile { return nil }
func (s *coreG2TestSlot) String() string                        { return "TestSlot(" + s.name + ")" }

// coreG2Factory is CoreG2Dump.FactorySweep.
type coreG2Factory struct {
	o             *coreG2Out
	f             *AstFactory
	mode          string
	registry      *rhino.JSTypeRegistry
	colorRegistry *colors.ColorRegistry
	libs          *RuntimeJsLibManager
}

func coreG2FactorySweep() string {
	o := &coreG2Out{}
	for _, stage := range []AbstractCompiler_LifeCycleStage{
		AbstractCompiler_LifeCycleStage_RAW, AbstractCompiler_LifeCycleStage_NORMALIZED,
	} {
		for _, mode := range []string{"NONE", "JSTYPE", "COLOR"} {
			registry := rhino.NewJSTypeRegistry(nil)
			colorRegistry := colors.ColorRegistry_Builder_().SetDefaultNativeColorsForTesting().Build()
			libs := RuntimeJsLibManager_Create(
				RuntimeJsLibManager_RuntimeLibraryMode_RECORD_AND_VALIDATE_FIELDS,
				func(unused1 string, unused2 string) rhino.Node {
					panic(jutil.NewUnsupportedOperationException())
				},
				NewChangeTracker(),
				func() rhino.Node { panic(jutil.NewUnsupportedOperationException()) })
			var f *AstFactory
			switch mode {
			case "NONE":
				f = AstFactory_CreateFactoryWithoutTypes(stage, libs)
			case "JSTYPE":
				f = AstFactory_CreateFactoryWithTypes(stage, registry, libs)
			default:
				f = AstFactory_CreateFactoryWithColors(stage, colorRegistry, libs)
			}
			o.line("# factory " + stage.String() + " " + mode)
			(&coreG2Factory{o: o, f: f, mode: mode, registry: registry, colorRegistry: colorRegistry, libs: libs}).sweep()
		}
	}
	return o.sb.String()
}

func (s *coreG2Factory) c(name string, body func() rhino.Node) {
	s.o.line("## " + name)
	var n rhino.Node
	if e := coreG2Try(func() { n = body() }); e != nil {
		s.o.line("! " + coreG2Exc(e))
		return
	}
	if n.IsNil() {
		s.o.line("null")
	} else {
		s.o.printTree(n)
	}
}

func (s *coreG2Factory) typed(n rhino.Node, nativeType rhino.JSTypeNative, color *colors.Color) rhino.Node {
	switch s.mode {
	case "JSTYPE":
		n.SetJSType(s.registry.GetNativeType(nativeType))
	case "COLOR":
		n.SetColor(color)
	}
	return n
}

func (s *coreG2Factory) num(d float64) rhino.Node {
	return s.typed(rhino.IR_Number(d), rhino.JSTypeNative_NUMBER_TYPE, colors.StandardColors_NUMBER)
}

func (s *coreG2Factory) str(v string) rhino.Node {
	return s.typed(rhino.IR_String(v), rhino.JSTypeNative_STRING_TYPE, colors.StandardColors_STRING)
}

func (s *coreG2Factory) name(v string) rhino.Node {
	return s.typed(rhino.IR_Name(v), rhino.JSTypeNative_NUMBER_TYPE, colors.StandardColors_NUMBER)
}

func (s *coreG2Factory) untypedName(v string) rhino.Node { return rhino.IR_Name(v) }

func (s *coreG2Factory) scope() *coreG2TestScope {
	objectDecl := s.typed(rhino.IR_Name("Object"), rhino.JSTypeNative_OBJECT_FUNCTION_TYPE, colors.StandardColors_TOP_OBJECT)
	assignDecl := s.typed(rhino.IR_Name("assign"), rhino.JSTypeNative_FUNCTION_TYPE, colors.StandardColors_TOP_OBJECT)
	constDecl := s.name("k")
	constDecl.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	return newCoreG2TestScope().
		declare("a", s.name("a")).
		declare("a.b", s.str("b")).
		declare("a.b.c", s.num(3)).
		declare("k", constDecl).
		declare("untyped", s.untypedName("untyped")).
		declare("Object", objectDecl).
		declare("Object.assign", assignDecl).
		declare("$jscomp", s.untypedName("$jscomp")).
		declare("$jscomp.global", s.untypedName("global")).
		declare("$jscomp.makeIterator", s.untypedName("makeIterator")).
		declare("$jscomp.arrayFromIterator", s.untypedName("arrayFromIterator")).
		declare("$jscomp.arrayFromIterable", s.untypedName("arrayFromIterable")).
		declare("$jscomp.makeAsyncIterator", s.untypedName("makeAsyncIterator")).
		declare("$jscomp.asyncExecutePromiseGeneratorFunction", s.untypedName("asyncExecutePromiseGeneratorFunction")).
		declare("$jscomp.global.Math", s.untypedName("Math")).
		declare("$jscomp.global.Math.max", s.untypedName("max")).
		declare("Math", s.untypedName("Math")).
		declare("Math.max", s.untypedName("max"))
}

func (s *coreG2Factory) numberType() AstFactory_Type {
	return AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_NUMBER_TYPE, colors.StandardColors_NUMBER)
}

func (s *coreG2Factory) functionType() AstFactory_Type {
	return AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_FUNCTION_TYPE, colors.StandardColors_TOP_OBJECT)
}

func (s *coreG2Factory) sweep() {
	f := s.f
	c := s.c
	c("exprResult", func() rhino.Node { return f.ExprResult(s.num(1)) })
	c("createEmpty", func() rhino.Node { return f.CreateEmpty() })
	c("createBlock", func() rhino.Node { return f.CreateBlock(f.ExprResult(s.num(1)), f.CreateEmpty()) })
	c("createIf", func() rhino.Node { return f.CreateIf(s.name("x"), f.CreateBlock()) })
	c("createIfElse", func() rhino.Node { return f.CreateIfWithElse(s.name("x"), f.CreateBlock(), f.CreateBlock()) })
	c("createFor", func() rhino.Node {
		return f.CreateFor(f.CreateEmpty(), s.name("x"), f.CreateEmpty(), f.CreateBlock())
	})
	c("createBreak", func() rhino.Node { return f.CreateBreak() })
	c("createLabel", func() rhino.Node { return f.CreateLabel(rhino.IR_LabelName("L"), f.CreateBlock()) })
	c("createCatch", func() rhino.Node { return f.CreateCatch(s.name("e"), f.CreateBlock()) })
	c("createTryFinally", func() rhino.Node { return f.CreateTryFinally(f.CreateBlock(), f.CreateBlock()) })
	c("createTryCatchFinally", func() rhino.Node {
		return f.CreateTryCatchFinally(f.CreateBlock(), f.CreateCatch(s.name("e"), f.CreateBlock()), f.CreateBlock())
	})
	c("createTryCatchFinally/bad", func() rhino.Node {
		return f.CreateTryCatchFinally(f.CreateBlock(), f.CreateBlock(), f.CreateBlock())
	})
	c("createThrow", func() rhino.Node { return f.CreateThrow(s.str("x")) })
	c("createReturn", func() rhino.Node { return f.CreateReturn(s.num(2)) })
	c("createYield", func() rhino.Node { return f.CreateYield(s.numberType(), s.num(1)) })
	c("createAwait", func() rhino.Node { return f.CreateAwait(AstFactory_Type_(s.name("p")), s.name("p")) })
	c("createString", func() rhino.Node { return f.CreateString("hello") })
	c("createString/unicode", func() rhino.Node { return f.CreateString("h\u00e9\U0001F600") })
	c("createNumber", func() rhino.Node { return f.CreateNumber(0.5) })
	c("createNumber/negative", func() rhino.Node { return f.CreateNumber(-0.5) })
	c("createBoolean/true", func() rhino.Node { return f.CreateBoolean(true) })
	c("createBoolean/false", func() rhino.Node { return f.CreateBoolean(false) })
	c("createNull", func() rhino.Node { return f.CreateNull() })
	c("createVoid", func() rhino.Node { return f.CreateVoid(s.name("x")) })
	c("createUndefinedValue", func() rhino.Node { return f.CreateUndefinedValue() })
	c("createNot", func() rhino.Node { return f.CreateNot(s.name("x")) })
	c("createThis", func() rhino.Node { return f.CreateThis(s.numberType()) })
	c("createSuper", func() rhino.Node { return f.CreateSuper(AstFactory_Type_(s.name("s"))) })
	c("createThisForEs6Class/notClass", func() rhino.Node { return f.CreateThisForEs6Class(s.name("x")) })
	c("createSingleNameDeclaration/let", func() rhino.Node {
		return f.CreateSingleNameDeclaration(rhino.Token_LET, "v", s.num(1))
	})
	c("createSingleNameDeclaration/var", func() rhino.Node {
		return f.CreateSingleNameDeclaration(rhino.Token_VAR, "v", s.num(1))
	})
	c("createSingleNameDeclaration/const", func() rhino.Node {
		return f.CreateSingleNameDeclaration(rhino.Token_CONST, "v", s.num(1))
	})
	c("createSingleNameDeclaration/bad", func() rhino.Node {
		return f.CreateSingleNameDeclaration(rhino.Token_NAME, "v", s.num(1))
	})
	c("createSingleLetNameDeclaration", func() rhino.Node { return f.CreateSingleLetNameDeclaration("v") })
	c("createSingleLetNameDeclaration/value", func() rhino.Node {
		return f.CreateSingleLetNameDeclarationWithValue("v", s.str("s"))
	})
	c("createSingleVarNameDeclaration", func() rhino.Node { return f.CreateSingleVarNameDeclaration("v") })
	c("createSingleVarNameDeclaration/value", func() rhino.Node {
		return f.CreateSingleVarNameDeclarationWithValue("$jscomp$v", s.str("s"))
	})
	c("createSingleConstNameDeclaration", func() rhino.Node {
		return f.CreateSingleConstNameDeclaration("v", s.str("s"))
	})
	c("createSingleConstObjectPatternDeclaration", func() rhino.Node {
		return f.CreateSingleConstObjectPatternDeclaration(rhino.IR_ObjectPattern(rhino.IR_StringKey("Foo")), s.name("x"))
	})
	c("createSingleConstObjectPatternDeclaration/bad", func() rhino.Node {
		return f.CreateSingleConstObjectPatternDeclaration(s.name("y"), s.name("x"))
	})
	c("createArgumentsReference", func() rhino.Node { return f.CreateArgumentsReference() })
	c("createArgumentsAliasDeclaration", func() rhino.Node { return f.CreateArgumentsAliasDeclaration("args") })
	c("createName", func() rhino.Node { return f.CreateName("n", s.numberType()) })
	c("createName/jscompPrefix", func() rhino.Node { return f.CreateName("$jscomp$x", s.numberType()) })
	c("createName/jscomp", func() rhino.Node { return f.CreateName("$jscomp", s.numberType()) })
	c("createConstantName", func() rhino.Node { return f.CreateConstantName("n", s.numberType()) })
	c("createConstantName/jscomp", func() rhino.Node { return f.CreateConstantName("$jscomp", s.numberType()) })
	c("createName/scope", func() rhino.Node { return f.CreateNameInScope(s.scope(), "a") })
	c("createName/scope/const", func() rhino.Node { return f.CreateNameInScope(s.scope(), "k") })
	c("createName/scope/untyped", func() rhino.Node { return f.CreateNameInScope(s.scope(), "untyped") })
	c("createName/scope/missing", func() rhino.Node { return f.CreateNameInScope(s.scope(), "missing") })
	c("createName/nullScope", func() rhino.Node { return f.CreateNameInScope(nil, "a") })
	c("createName/scope/jscomp", func() rhino.Node { return f.CreateNameInScope(s.scope(), "$jscomp") })
	c("createNameWithUnknownType", func() rhino.Node { return f.CreateNameWithUnknownType("u") })
	c("createQName/dotted", func() rhino.Node { return f.CreateQName(s.scope(), "a.b.c") })
	c("createQName/missing", func() rhino.Node { return f.CreateQName(s.scope(), "a.b.missing") })
	c("createQName/iterable", func() rhino.Node { return f.CreateQNameOfNames(s.scope(), []string{"a", "b"}) })
	c("createQName/varargs", func() rhino.Node { return f.CreateQNameWithProperties(s.scope(), "a", "b", "c") })
	c("createQName/baseAndIterable", func() rhino.Node {
		return f.CreateQNameWithPropertyList(s.scope(), "a", []string{"b", "c"})
	})
	c("createQName/field/notInjected", func() rhino.Node {
		return f.CreateQNameOfField(s.scope(), s.libs.GetJsLibField("$jscomp.global"))
	})
	c("createQName/field", func() rhino.Node {
		s.libs.InjectLibForField("$jscomp.global")
		return f.CreateQNameOfField(s.scope(), s.libs.GetJsLibField("$jscomp.global"))
	})
	c("createQNameWithUnknownType", func() rhino.Node { return f.CreateQNameWithUnknownType("p.q.r") })
	c("createQNameWithUnknownType/field", func() rhino.Node {
		return f.CreateQNameWithUnknownTypeOfField(s.libs.GetJsLibField("$jscomp.global"))
	})
	c("createQNameWithUnknownType/baseAndIterable", func() rhino.Node {
		return f.CreateQNameWithUnknownTypeWithPropertyList("p", []string{"q", "r"})
	})
	c("createPrototypeAccess", func() rhino.Node { return f.CreatePrototypeAccess(s.name("Foo")) })
	c("createPrototypeAccess/untyped", func() rhino.Node { return f.CreatePrototypeAccess(s.untypedName("Foo")) })
	c("createJSCompDotGlobalAccess", func() rhino.Node { return f.CreateJSCompDotGlobalAccess(s.scope(), "Math.max") })
	c("createGetPropWithoutColor", func() rhino.Node { return f.CreateGetPropWithoutColor(s.name("o"), "p") })
	c("createGetPropWithoutColor/jscompGlobal", func() rhino.Node {
		return f.CreateGetPropWithoutColor(s.untypedName("$jscomp"), "global")
	})
	c("createGetPropWithoutColor/string", func() rhino.Node {
		return f.CreateGetPropWithoutColor(s.str("s"), "length")
	})
	c("createGetProp", func() rhino.Node { return f.CreateGetProp(s.name("o"), "p", s.numberType()) })
	c("createGetPropsWithoutColors", func() rhino.Node {
		return f.CreateGetPropsWithoutColors(s.name("o"), []string{"p", "q"})
	})
	c("createGetPropsWithUnknownType", func() rhino.Node {
		return f.CreateGetPropsWithUnknownType(s.name("o"), []string{"p", "q"})
	})
	c("createGetPropWithUnknownType", func() rhino.Node { return f.CreateGetPropWithUnknownType(s.name("o"), "p") })
	c("createStartOptChainGetprop", func() rhino.Node {
		return f.CreateStartOptChainGetprop(s.name("o"), "p", s.numberType())
	})
	c("createContinueOptChainGetprop", func() rhino.Node {
		return f.CreateContinueOptChainGetprop(
			f.CreateStartOptChainGetprop(s.name("o"), "p", s.numberType()), "q", s.numberType())
	})
	c("createStartOptChainGetelem", func() rhino.Node {
		return f.CreateStartOptChainGetelem(s.name("o"), s.num(0), s.numberType())
	})
	c("createContinueOptChainGetelem", func() rhino.Node {
		return f.CreateContinueOptChainGetelem(
			f.CreateStartOptChainGetelem(s.name("o"), s.num(0), s.numberType()), s.num(1), s.numberType())
	})
	c("createStartOptChainCall", func() rhino.Node {
		return f.CreateStartOptChainCall(s.name("o"), s.numberType(), s.num(1), s.str("x"))
	})
	c("createContinueOptChainCall", func() rhino.Node {
		return f.CreateContinueOptChainCall(
			f.CreateStartOptChainGetprop(s.name("o"), "p", s.numberType()), s.numberType())
	})
	c("createGetElem", func() rhino.Node { return f.CreateGetElem(s.name("o"), s.num(0)) })
	c("createDelProp", func() rhino.Node { return f.CreateDelProp(f.CreateGetPropWithUnknownType(s.name("o"), "p")) })
	c("createStringKey", func() rhino.Node { return f.CreateStringKey("k", s.str("v")) })
	c("createComputedProperty", func() rhino.Node { return f.CreateComputedProperty(s.str("k"), s.num(1)) })
	c("createGetterDef", func() rhino.Node { return f.CreateGetterDef("g", s.num(1)) })
	c("createIn", func() rhino.Node { return f.CreateIn(s.str("k"), s.name("o")) })
	c("createComma", func() rhino.Node { return f.CreateComma(s.name("a"), s.str("b")) })
	c("createCommas", func() rhino.Node { return f.CreateCommas(s.name("a"), s.str("b"), s.num(1), s.name("c")) })
	c("createAnd", func() rhino.Node { return f.CreateAnd(s.name("a"), s.str("b")) })
	c("createAnd/untyped", func() rhino.Node { return f.CreateAnd(s.untypedName("a"), s.str("b")) })
	c("createOr", func() rhino.Node { return f.CreateOr(s.str("a"), s.num(1)) })
	c("createAdd", func() rhino.Node { return f.CreateAdd(s.str("a"), s.num(1)) })
	c("createSub", func() rhino.Node { return f.CreateSub(s.name("a"), s.num(1)) })
	c("createInc/pre", func() rhino.Node { return f.CreateInc(s.name("a"), false) })
	c("createInc/post", func() rhino.Node { return f.CreateInc(s.name("a"), true) })
	c("createLessThan", func() rhino.Node { return f.CreateLessThan(s.name("a"), s.num(1)) })
	c("createBitwiseAnd", func() rhino.Node { return f.CreateBitwiseAnd(s.name("a"), s.num(1)) })
	c("createRightShift", func() rhino.Node { return f.CreateRightShift(s.name("a"), s.num(1)) })
	c("createCall", func() rhino.Node { return f.CreateCall(s.name("fn"), s.numberType(), s.num(1), s.str("x")) })
	c("createCall/getprop", func() rhino.Node {
		return f.CreateCall(f.CreateGetPropWithUnknownType(s.name("o"), "m"), s.numberType())
	})
	c("createCallWithUnknownType", func() rhino.Node { return f.CreateCallWithUnknownType(s.name("fn"), s.num(1)) })
	c("createObjectDotAssignCall", func() rhino.Node {
		return f.CreateObjectDotAssignCall(s.scope(), s.numberType(), f.CreateObjectLit(), s.name("a"))
	})
	c("createNewNode/unknownCtor", func() rhino.Node { return f.CreateNewNode(s.untypedName("C"), s.num(1)) })
	c("createNewNode/objectCtor", func() rhino.Node {
		return f.CreateNewNode(
			s.typed(rhino.IR_Name("Object"), rhino.JSTypeNative_OBJECT_FUNCTION_TYPE, s.colorRegistry.Get(colors.StandardColors_ARRAY_ID)),
			s.num(1))
	})
	c("createConstructorCall", func() rhino.Node {
		return f.CreateConstructorCall(
			AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_OBJECT_FUNCTION_TYPE, colors.StandardColors_TOP_OBJECT),
			rhino.IR_SuperNode(),
			s.num(1))
	})
	c("createAssignStatement", func() rhino.Node { return f.CreateAssignStatement(s.name("a"), s.str("b")) })
	c("createAssign", func() rhino.Node { return f.CreateAssign(s.name("a"), s.str("b")) })
	c("createAssign/name", func() rhino.Node { return f.CreateAssignToName("a", s.str("b")) })
	c("createObjectLit", func() rhino.Node { return f.CreateObjectLit(f.CreateStringKey("k", s.num(1))) })
	c("createObjectLit/type", func() rhino.Node {
		return f.CreateObjectLitWithType(s.numberType(), f.CreateStringKey("k", s.num(1)))
	})
	c("createQuotedStringKey", func() rhino.Node { return f.CreateQuotedStringKey("k", s.num(1)) })
	c("createEmptyFunction", func() rhino.Node { return f.CreateEmptyFunction(s.functionType()) })
	c("createEmptyGeneratorFunction", func() rhino.Node { return f.CreateEmptyGeneratorFunction(s.functionType()) })
	c("createFunction", func() rhino.Node {
		return f.CreateFunction("fn", f.CreateParamList("p", "q"), f.CreateBlock(), s.functionType())
	})
	c("createParamList", func() rhino.Node { return f.CreateParamList("p", "$jscomp$q") })
	c("createZeroArgFunction", func() rhino.Node {
		return f.CreateZeroArgFunction(
			"z", f.CreateBlock(f.CreateReturn(s.num(1))), s.registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	})
	c("createZeroArgFunction/nullReturn", func() rhino.Node {
		return f.CreateZeroArgFunction("z", f.CreateBlock(), nil)
	})
	c("createZeroArgGeneratorFunction", func() rhino.Node {
		return f.CreateZeroArgGeneratorFunction("g", f.CreateBlock(), nil)
	})
	c("createZeroArgArrowFunctionForExpression", func() rhino.Node {
		return f.CreateZeroArgArrowFunctionForExpression(s.num(1))
	})
	c("createZeroArgArrowFunctionForExpression/untyped", func() rhino.Node {
		return f.CreateZeroArgArrowFunctionForExpression(s.untypedName("u"))
	})
	c("createMemberFunctionDef", func() rhino.Node {
		return f.CreateMemberFunctionDef("m", f.CreateEmptyFunction(s.functionType()))
	})
	c("createMemberFunctionDef/named", func() rhino.Node {
		return f.CreateMemberFunctionDef("m", f.CreateFunction("n", rhino.IR_ParamList(), f.CreateBlock(), s.functionType()))
	})
	c("createSheq", func() rhino.Node { return f.CreateSheq(s.name("a"), s.num(1)) })
	c("createEq", func() rhino.Node { return f.CreateEq(s.name("a"), f.CreateNull()) })
	c("createNe", func() rhino.Node { return f.CreateNe(s.name("a"), f.CreateNull()) })
	c("createHook", func() rhino.Node { return f.CreateHook(s.name("a"), s.str("b"), s.num(1)) })
	c("createHook/untyped", func() rhino.Node { return f.CreateHook(s.name("a"), s.untypedName("b"), s.num(1)) })
	c("createArraylit", func() rhino.Node { return f.CreateArraylit(s.num(1), s.str("b")) })
	c("createArraylit/iterable", func() rhino.Node {
		elems := []rhino.Node{s.num(1), s.name("c")}
		return f.CreateArraylitOfIterable(func(yield func(rhino.Node) bool) {
			for _, e := range elems {
				if !yield(e) {
					return
				}
			}
		})
	})
	c("createJSCompMakeIteratorCall/notInjected", func() rhino.Node {
		return f.CreateJSCompMakeIteratorCall(s.name("arr"), s.scope())
	})
	c("createJSCompMakeIteratorCall", func() rhino.Node {
		s.libs.InjectLibForField("$jscomp.makeIterator")
		return f.CreateJSCompMakeIteratorCall(s.name("arr"), s.scope())
	})
	c("createJscompArrayFromIteratorCall", func() rhino.Node {
		s.libs.InjectLibForField("$jscomp.arrayFromIterator")
		return f.CreateJscompArrayFromIteratorCall(s.name("it"), s.scope())
	})
	c("createJscompArrayFromIterableCall", func() rhino.Node {
		s.libs.InjectLibForField("$jscomp.arrayFromIterable")
		return f.CreateJscompArrayFromIterableCall(s.name("it"), s.scope())
	})
	c("createJSCompMakeAsyncIteratorCall", func() rhino.Node {
		s.libs.InjectLibForField("$jscomp.makeAsyncIterator")
		return f.CreateJSCompMakeAsyncIteratorCall(s.name("it"), s.scope())
	})
	c("createEmptyAsyncGeneratorWrapperArgument/unknown", func() rhino.Node {
		return f.CreateEmptyAsyncGeneratorWrapperArgument(s.registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE))
	})
	c("createEmptyAsyncGeneratorWrapperArgument/function", func() rhino.Node {
		return f.CreateEmptyAsyncGeneratorWrapperArgument(
			s.registry.CreateFunctionType(
				s.registry.GetNativeType(rhino.JSTypeNative_VOID_TYPE),
				s.registry.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)))
	})
	c("createJscompAsyncExecutePromiseGeneratorFunctionCall", func() rhino.Node {
		s.libs.InjectLibForField("$jscomp.asyncExecutePromiseGeneratorFunction")
		return f.CreateJscompAsyncExecutePromiseGeneratorFunctionCall(
			s.scope(), f.CreateEmptyGeneratorFunction(s.functionType()))
	})
	c("type/colorId", func() rhino.Node { return f.CreateThis(AstFactory_Type_OfColorId(colors.StandardColors_ARRAY_ID)) })
	c("type/jstypeNative", func() rhino.Node {
		return f.CreateThis(AstFactory_Type_OfJSTypeNative(rhino.JSTypeNative_STRING_TYPE))
	})
	c("type/color", func() rhino.Node { return f.CreateThis(AstFactory_Type_OfColor(colors.StandardColors_STRING)) })
	c("type/jstype", func() rhino.Node {
		return f.CreateThis(AstFactory_Type_OfJSType(s.registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE)))
	})
	s.o.line("## isAddingTypes " + jutil.ToJavaString(f.IsAddingTypes()) + " isAddingColors " + jutil.ToJavaString(f.IsAddingColors()))
}
