package jscomp

// Differential test against tools/oracle/java/L11OptimizeDump.java: the
// l11-optimize passes (CollapseAnonymousFunctions, CollapseVariableDeclarations,
// Denormalize, DeadPropertyAssignmentElimination, InlineObjectLiterals,
// OptimizeCalls, FunctionArgumentInjector, ExtractPrototypeMemberDeclarations,
// AnalyzePrototypeProperties, AliasStrings, RenameLabels, RenameVars,
// ConcretizeStaticInheritanceForInlining, InferConsts, RemoveCastNodes,
// RemovePropertyRenamingCalls, RemoveUnnecessarySyntheticExterns) run on every
// file after the prerequisites DefaultPassConfig runs before them, and the
// resulting AST, diagnostics and pass-specific results (the VariableMap of
// RenameVars, the ReferenceMap of OptimizeCalls, ...) printed in the same
// format as the Java tool (see its header for the format).
//
//   - TestL11OptimizeGolden replays the checked-in golden dump
//     testdata/l11opt/golden.txt.gz (with --full trees): the hand-written
//     inputs testdata/l11opt/js/l11opt-*.js followed by the corpus files of
//     testdata/l11opt/golden-files.txt (every 150th file of corpus-files.txt
//     between 200 bytes and 12 KB). Regenerate with
//       CP=.cache/$CLOSURE_JAR:.cache/oracle-build
//       J="java -Xss1g -XX:-ShowCodeDetailsInExceptionMessages -cp $CP"
//       $J com.google.javascript.jscomp.L11OptimizeDump --full a.gz \
//         internal/jscomp/testdata/l11opt/js --list=<the l11opt-*.js names>
//       $J com.google.javascript.jscomp.L11OptimizeDump --full b.gz \
//         .cache/corpus --list=internal/jscomp/testdata/l11opt/golden-files.txt
//       cat a.gz b.gz > internal/jscomp/testdata/l11opt/golden.txt.gz
//   - TestL11OptimizeOracle compares any other dump (trees hashed, not
//     printed): L11OPT_ORACLE=<dump.gz> [L11OPT_ROOT=<root>] (default root:
//     the corpus under /home/user/jsgo/.cache/corpus). The held-out corpus
//     (.cache/corpus-holdout, corpus-holdout-files.txt) is for a final check
//     only.

import (
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// l11optEscape is L11OptimizeDump.escape (ParsingDump's escaping).
func l11optEscape(s string) string {
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

type l11optStage int

const (
	l11optRaw l11optStage = iota
	l11optNormalized
	l11optNormalizedWithProperties
	l11optLate
	l11optPipeline
)

type l11optPass struct {
	name  string
	stage l11optStage
	run   func(c *Compiler, externs rhino.Node, js rhino.Node, extra *[]string)
}

// l11optPasses is L11OptimizeDump.passes(), in the same order.
func l11optPasses() []l11optPass {
	var p []l11optPass
	// raw
	p = append(p, l11optPass{"RemoveCastNodes", l11optRaw, func(c *Compiler, externs, js rhino.Node, extra *[]string) {
		NewRemoveCastNodes(c).Process(externs, js)
	}})
	p = append(p, l11optPass{"RemovePropertyRenamingCalls", l11optRaw, func(c *Compiler, externs, js rhino.Node, extra *[]string) {
		NewRemovePropertyRenamingCalls(c).Process(externs, js)
	}})
	p = append(p, l11optPass{"RemoveUnnecessarySyntheticExterns", l11optRaw, func(c *Compiler, externs, js rhino.Node, extra *[]string) {
		NewRemoveUnnecessarySyntheticExterns(c).Process(externs, js)
	}})
	p = append(p, l11optPass{"RenameLabels", l11optRaw, func(c *Compiler, externs, js rhino.Node, extra *[]string) {
		NewRenameLabels(c).Process(externs, js)
	}})
	p = append(p, l11optPass{"RenameVars", l11optRaw, l11optRenameVars})
	// normalized
	p = append(p, l11optPass{"CollapseAnonymousFunctions", l11optNormalized, func(c *Compiler, externs, js rhino.Node, extra *[]string) {
		NewCollapseAnonymousFunctions(c).Process(externs, js)
	}})
	for _, pattern := range ExtractPrototypeMemberDeclarations_Pattern_values() {
		p = append(p, l11optPass{"ExtractPrototypeMemberDeclarations/" + pattern.String(), l11optNormalized, func(c *Compiler, externs, js rhino.Node, extra *[]string) {
			NewExtractPrototypeMemberDeclarations(c, pattern).Process(externs, js)
		}})
	}
	for _, mode := range []CompilerOptions_AliasStringsMode{
		CompilerOptions_AliasStringsMode_ALL, CompilerOptions_AliasStringsMode_LARGE, CompilerOptions_AliasStringsMode_ALL_AGGRESSIVE,
	} {
		p = append(p, l11optPass{"AliasStrings/" + mode.String(), l11optNormalized, func(c *Compiler, externs, js rhino.Node, extra *[]string) {
			NewAliasStrings(c, c.GetChunkGraph(), false, mode).Process(externs, js)
		}})
	}
	p = append(p, l11optPass{"InlineObjectLiterals", l11optNormalized, func(c *Compiler, externs, js rhino.Node, extra *[]string) {
		NewInlineObjectLiterals(c, c.GetUniqueNameIdSupplier()).Process(externs, js)
	}})
	p = append(p, l11optPass{"DeadPropertyAssignmentElimination", l11optNormalizedWithProperties, func(c *Compiler, externs, js rhino.Node, extra *[]string) {
		NewDeadPropertyAssignmentElimination(c).Process(externs, js)
	}})
	p = append(p, l11optPass{"OptimizeCalls", l11optNormalizedWithProperties, l11optOptimizeCalls})
	p = append(p, l11optPass{"AnalyzePrototypeProperties", l11optNormalized, l11optAnalyzePrototypeProperties})
	p = append(p, l11optPass{"ConcretizeStaticInheritanceForInlining", l11optNormalized, func(c *Compiler, externs, js rhino.Node, extra *[]string) {
		NewConcretizeStaticInheritanceForInlining(c).Process(externs, js)
	}})
	p = append(p, l11optPass{"InferConsts", l11optNormalized, func(c *Compiler, externs, js rhino.Node, extra *[]string) {
		NewInferConsts(c).Process(externs, js)
		l11optDumpConsts(js, extra)
	}})
	p = append(p, l11optPass{"FunctionArgumentInjector", l11optNormalized, func(c *Compiler, externs, js rhino.Node, extra *[]string) {
		l11optFunctionArgumentInjector(c, js, extra)
	}})
	// late
	p = append(p, l11optPass{"Denormalize", l11optLate, func(c *Compiler, externs, js rhino.Node, extra *[]string) {
		NewDenormalize(c, c.GetOptions().GetOutputFeatureSet()).Process(externs, js)
	}})
	p = append(p, l11optPass{"CollapseVariableDeclarations", l11optLate, func(c *Compiler, externs, js rhino.Node, extra *[]string) {
		NewCollapseVariableDeclarations(c).Process(externs, js)
	}})
	// pipeline
	p = append(p, l11optPass{"RenameVarsAndLabels", l11optPipeline, func(c *Compiler, externs, js rhino.Node, extra *[]string) {
		l11optRenameVars(c, externs, js, extra)
		NewRenameLabels(c).Process(externs, js)
	}})
	return p
}

// l11optRenameVars is L11OptimizeDump.renameVars: RenameVars configured
// like DefaultPassConfig.runVariableRenaming with default options.
func l11optRenameVars(c *Compiler, externs rhino.Node, js rhino.Node, extra *[]string) {
	options := c.GetOptions()
	reservedChars := jutil.NewLinkedHashSet[uint16]()
	reservedNames := jutil.NewLinkedHashSet[string]()
	if options.GetRenamePrefixNamespace() != nil {
		reservedNames.Add(*options.GetRenamePrefixNamespace())
	}
	reservedNames.AddAll(c.GetExportedNames().ToSlice()...)
	reservedNames.AddAll(parsing.ParserRunner_GetReservedVars().ToSlice()...)
	rn := NewRenameVars(
		c,
		options.GetRenamePrefix(),
		options.GetVariableRenaming() == VariableRenamingPolicy_LOCAL,
		options.ShouldGeneratePseudoNames(),
		options.ShouldPreferStableNames(),
		options.GetInputVariableMap(),
		reservedChars,
		reservedNames,
		options.GetNameGenerator())
	rn.Process(externs, js)
	for k, v := range rn.GetVariableMap().GetOriginalNameToNewNameMap().All() {
		*extra = append(*extra, "map "+k+" -> "+v)
	}
}

func l11optOptimizeCalls(c *Compiler, externs rhino.Node, js rhino.Node, extra *[]string) {
	for _, considerExterns := range []bool{false, true} {
		OptimizeCalls_Builder_().
			SetCompiler(c).
			SetConsiderExterns(considerExterns).
			AddPass(OptimizeCalls_CallGraphCompilerPassFunc(func(e rhino.Node, r rhino.Node, references *OptimizeCalls_ReferenceMap) {
				*extra = append(*extra, "references considerExterns="+strconv.FormatBool(considerExterns))
				for k, refs := range references.GetNameReferences() {
					*extra = append(*extra, "name "+k+l11optRefs(refs))
					l11optFunctionNodes(refs, extra)
				}
				for k, refs := range references.GetPropReferences() {
					*extra = append(*extra, "prop "+k+l11optRefs(refs))
					l11optFunctionNodes(refs, extra)
				}
				for _, refs := range references.GetNameReferences() {
					for _, n := range refs {
						if !n.GetParent().IsNil() {
							l11optCallTarget(n, extra)
						}
					}
				}
			})).
			Build().
			Process(externs, js)
	}
}

func l11optRefs(nodes []rhino.Node) string {
	var sb strings.Builder
	for _, n := range nodes {
		sb.WriteByte(' ')
		sb.WriteString(l11optPos(n))
	}
	return sb.String()
}

func l11optFunctionNodes(nodes []rhino.Node, extra *[]string) {
	var keys []rhino.Node
	byDef := make(map[rhino.Node][]rhino.Node)
	for _, def := range nodes {
		fns := optimizeCalls_appendDefinitionFunctionNodesFor(nil, def)
		if len(fns) == 0 {
			continue
		}
		if _, ok := byDef[def]; !ok {
			keys = append(keys, def)
		}
		byDef[def] = append(byDef[def], fns...)
	}
	for _, k := range keys {
		for _, v := range byDef[k] {
			*extra = append(*extra, "  fn "+l11optPos(k)+" -> "+l11optPos(v))
		}
	}
}

func l11optCallTarget(n rhino.Node, extra *[]string) {
	var sb strings.Builder
	sb.WriteString("  target " + l11optPos(n))
	found := false
	if OptimizeCalls_ReferenceMap_IsNormalOrOptChainCallOrNewTarget(n) {
		found = true
		sb.WriteString(" call=" + l11optPos(OptimizeCalls_ReferenceMap_GetCallOrNewNodeForTarget(n)))
		arg := OptimizeCalls_ReferenceMap_GetFirstArgumentForCallOrNewOrDotCall(n)
		if arg.IsNil() {
			sb.WriteString(" arg0=null")
		} else {
			sb.WriteString(" arg0=" + l11optPos(arg))
		}
	}
	if OptimizeCalls_ReferenceMap_IsSimpleAssignmentTarget(n) {
		found = true
		sb.WriteString(" simpleAssign")
	}
	if OptimizeCalls_IsAllowedReference(n) {
		found = true
		sb.WriteString(" allowed")
	}
	if found {
		*extra = append(*extra, sb.String())
	}
}

func l11optChunkName(chunk *JSChunk) string {
	if chunk == nil {
		return "null"
	}
	return chunk.GetName()
}

func l11optAnalyzePrototypeProperties(c *Compiler, externs rhino.Node, js rhino.Node, extra *[]string) {
	app := NewAnalyzePrototypeProperties(
		c,
		c.GetChunkGraph(),
		false, /* canModifyExterns= */
		false, /* anchorUnusedVars= */
		false /* rootScopeUsesAreGlobal= */)
	app.Process(externs, js)
	for _, info := range app.GetAllNameInfo() {
		*extra = append(*extra,
			"info "+info.String()+
				" referenced="+strconv.FormatBool(info.IsReferenced())+
				" readsClosureVariables="+strconv.FormatBool(info.ReadsClosureVariables())+
				" referencesSuper="+strconv.FormatBool(info.ReferencesSuper())+
				" chunk="+l11optChunkName(info.GetDeepestCommonChunkRef()))
		for _, s := range info.GetDeclarations() {
			var kind string
			var node rhino.Node
			switch d := s.(type) {
			case *AnalyzePrototypeProperties_GlobalFunction:
				kind = "GlobalFunction"
			case *AnalyzePrototypeProperties_ClassMemberFunction:
				kind = "ClassMemberFunction"
				node = d.GetDefinitionNode()
			case AnalyzePrototypeProperties_PrototypeProperty:
				// s.getClass().getSimpleName()
				simple := "AssignmentPrototypeProperty"
				if _, ok := d.(*AnalyzePrototypeProperties_LiteralPrototypeProperty); ok {
					simple = "LiteralPrototypeProperty"
				}
				kind = simple + " prototype=" + l11optPos(d.GetPrototype())
				node = d.GetValue()
			}
			line := "  decl " + kind
			if !node.IsNil() {
				line += " node=" + l11optPos(node)
			}
			if v := s.GetRootVar(); v == nil {
				line += " rootVar=null"
			} else {
				line += " rootVar=" + v.GetName()
			}
			line += " chunk=" + l11optChunkName(s.GetChunk())
			*extra = append(*extra, line)
		}
	}
}

func l11optDumpConsts(n rhino.Node, extra *[]string) {
	if (n.IsName() || n.IsImportStar()) &&
		(n.IsDeclaredConstantVar() || n.IsInferredConstantVar()) {
		*extra = append(*extra,
			"const "+l11optPos(n)+" "+n.GetString()+
				" declared="+strconv.FormatBool(n.IsDeclaredConstantVar())+
				" inferred="+strconv.FormatBool(n.IsInferredConstantVar()))
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		l11optDumpConsts(c, extra)
	}
}

const l11optMaxInjectorPairs = 200

func l11optFunctionArgumentInjector(c *Compiler, js rhino.Node, extra *[]string) {
	injector := NewFunctionArgumentInjector(c.GetAstAnalyzer())
	var fns, calls []rhino.Node
	l11optCollect(js, &fns, &calls)
	pairs := 0
	for _, fn := range fns {
		name := fn.GetFirstChild().GetString()
		for _, call := range calls {
			if call.GetFirstChild().GetString() != name {
				continue
			}
			if pairs >= l11optMaxInjectorPairs {
				return
			}
			pairs++
			counter := 0
			ids := func() string {
				v := strconv.Itoa(counter)
				counter++
				return v
			}
			map_ := injector.GetFunctionCallParameterMap(fn, call, ids)
			var sb strings.Builder
			sb.WriteString("inject fn=" + l11optPos(fn) + " call=" + l11optPos(call) + " map")
			for k, v := range map_.All() {
				sb.WriteString(" " + k + "=" + v.Arg().GetToken().String() + "@" + l11optPos(v.Arg()))
			}
			*extra = append(*extra, sb.String())
			modified := injector.FindModifiedParameters(fn)
			*extra = append(*extra, "  modified "+modified.String())
			temps := injector.GatherCallArgumentsNeedingTemps(c, fn, map_, modified, c.GetCodingConvention())
			*extra = append(*extra, "  temps "+temps.String())
			*extra = append(*extra,
				"  conditional "+strconv.FormatBool(injector.BodyMayHaveConditionalCode(fn.GetLastChild()))+
					" "+strconv.FormatBool(injector.MayHaveConditionalCode(fn.GetLastChild())))
		}
	}
}

func l11optCollect(n rhino.Node, fns *[]rhino.Node, calls *[]rhino.Node) {
	if NodeUtil_IsFunctionDeclaration(n) && n.GetFirstChild().GetString() != "" {
		plain := true
		for p := NodeUtil_GetFunctionParameters(n).GetFirstChild(); !p.IsNil(); p = p.GetNext() {
			if !(p.IsName() || (p.IsRest() && p.GetOnlyChild().IsName())) {
				plain = false
			}
		}
		if plain {
			*fns = append(*fns, n)
		}
	}
	if n.IsCall() && n.GetFirstChild().IsName() {
		*calls = append(*calls, n)
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		l11optCollect(c, fns, calls)
	}
}

func l11optPos(n rhino.Node) string {
	return n.GetToken().String() + "@" + strconv.Itoa(n.GetLineno()) + ":" + strconv.Itoa(n.GetCharno())
}

// l11optNewCompiler is L11OptimizeDump.newCompiler.
func l11optNewCompiler(rel string, src string) *Compiler {
	compiler := NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	options.SetLanguageOut(CompilerOptions_LanguageMode_NO_TRANSPILE)
	compiler.Init(nil, []*SourceFile{SourceFile_FromCode(rel, src)}, options)
	compiler.Parse()
	return compiler
}

// l11optPrepare is L11OptimizeDump.prepare.
func l11optPrepare(c *Compiler, stage l11optStage) {
	externs := c.GetExternsRoot()
	js := c.GetJsRoot()
	if stage == l11optRaw {
		return
	}
	Normalize_CreateNormalizeForOptimizations(c).Process(externs, js)
	if stage == l11optNormalizedWithProperties {
		NewGatherGetterAndSetterProperties(c).Process(externs, js)
		NewGatherExternProperties(c, GatherExternProperties_Mode_OPTIMIZE).Process(externs, js)
	}
	if stage == l11optLate || stage == l11optPipeline {
		c.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_RAW)
	}
	if stage == l11optPipeline {
		NewDenormalize(c, c.GetOptions().GetOutputFeatureSet()).Process(externs, js)
		MakeDeclaredNamesUnique_GetContextualRenameInverter(c).Process(externs, js)
	}
}

func l11optDiag(level string, e *JSError) string {
	return "diag " + level + " " + e.GetType().key + " " + strconv.Itoa(e.GetLineno()) + ":" +
		strconv.Itoa(e.GetCharno()) + " \"" + e.GetDescription() + "\""
}

func l11optDiagnostics(c *Compiler) []string {
	var out []string
	for _, e := range c.GetErrors() {
		out = append(out, l11optDiag("ERROR", e))
	}
	for _, e := range c.GetWarnings() {
		out = append(out, l11optDiag("WARNING", e))
	}
	return out
}

// l11optException is "exception " + t.getClass().getName() + ": " +
// t.getMessage() for a recovered panic (a Go runtime fault becomes the Java
// exception the JVM throws, as in jutil.Try).
func l11optException(r any) string {
	th := jutil.AsThrowable(r)
	full := th.Error()
	msg := jutil.GetMessageOrNull(th)
	class := full
	if msg != nil && strings.HasSuffix(full, ": "+*msg) {
		class = strings.TrimSuffix(full, ": "+*msg)
	}
	return "exception " + class + ": " + jutil.ToJavaString(msg)
}

func l11optLine(sb *strings.Builder, s string) {
	sb.WriteString(l11optEscape(s))
	sb.WriteByte('\n')
}

// l11optDumpFile is L11OptimizeDump.dumpFile.
func l11optDumpFile(sb *strings.Builder, rel string, src string, full bool) {
	l11optLine(sb, "## file "+rel)
	compiler := l11optNewCompiler(rel, src)
	parseDiags := l11optDiagnostics(compiler)
	if compiler.HasErrors() {
		l11optLine(sb, "# parse errors")
	} else {
		l11optLine(sb, "# parse ok")
	}
	for _, d := range parseDiags {
		l11optLine(sb, d)
	}
	if compiler.HasErrors() {
		return
	}
	for _, pass := range l11optPasses() {
		l11optLine(sb, "## pass "+pass.name)
		c := l11optNewCompiler(rel, src)
		var extra []string
		var before, exception string
		func() {
			defer func() {
				if r := recover(); r != nil {
					exception = l11optException(r)
				}
			}()
			l11optPrepare(c, pass.stage)
			before = c.GetRoot().ToStringTree()
			pass.run(c, c.GetExternsRoot(), c.GetJsRoot(), &extra)
		}()
		diags := l11optDiagnostics(c)
		for _, d := range parseDiags {
			for i, x := range diags {
				if x == d {
					diags = append(diags[:i], diags[i+1:]...)
					break
				}
			}
		}
		for _, d := range diags {
			l11optLine(sb, d)
		}
		if exception != "" {
			l11optLine(sb, exception)
			continue
		}
		after := c.GetRoot().ToStringTree()
		if after == before {
			l11optLine(sb, "tree unchanged")
		} else {
			lines := strings.Split(after, "\n")
			for i := range lines {
				lines[i] = l11optEscape(lines[i])
			}
			if lines[len(lines)-1] == "" {
				lines = lines[:len(lines)-1]
			}
			h := sha256.New()
			for _, l := range lines {
				h.Write([]byte(l + "\n"))
			}
			l11optLine(sb, "tree changed "+hex.EncodeToString(h.Sum(nil))[:16])
			if full {
				for _, l := range lines {
					sb.WriteString(l)
					sb.WriteByte('\n')
				}
			}
		}
		for _, e := range extra {
			l11optLine(sb, e)
		}
	}
}

// l11optReadSections splits a gzipped dump into its "## file" sections.
func l11optReadSections(t *testing.T, path string) (names []string, sections map[string]string) {
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
	sections = map[string]string{}
	sc := bufio.NewScanner(zr)
	sc.Buffer(make([]byte, 1<<20), 1<<28)
	var cur string
	var sb strings.Builder
	flush := func() {
		if cur != "" {
			sections[cur] = sb.String()
		}
		sb.Reset()
	}
	for sc.Scan() {
		l := sc.Text()
		if strings.HasPrefix(l, "## file ") {
			flush()
			cur = l[len("## file "):]
			names = append(names, cur)
		}
		sb.WriteString(l)
		sb.WriteByte('\n')
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	flush()
	return names, sections
}

func l11optFirstDiff(want, got string) string {
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

// l11optCompare dumps every file of the Java dump and compares the sections.
func l11optCompare(t *testing.T, dumpPath string, full bool, rootFor func(name string) string) {
	names, sections := l11optReadSections(t, dumpPath)
	if len(names) == 0 {
		t.Fatalf("no sections in %s", dumpPath)
	}
	failures := 0
	for _, name := range names {
		src, err := os.ReadFile(filepath.Join(rootFor(name), name))
		if err != nil {
			t.Fatal(err)
		}
		var sb strings.Builder
		l11optDumpFile(&sb, name, string(src), full)
		if want, got := sections[name], sb.String(); got != want {
			failures++
			if failures <= 20 {
				t.Errorf("%s: differs from Java at %s", name, l11optFirstDiff(want, got))
			}
		}
	}
	if failures > 0 {
		t.Errorf("%d of %d files differ", failures, len(names))
	}
}

func TestL11OptimizeGolden(t *testing.T) {
	l11optCompare(t, "testdata/l11opt/golden.txt.gz", true, func(name string) string {
		if strings.HasPrefix(name, "l11opt-") {
			return "testdata/l11opt/js"
		}
		root := os.Getenv("L11OPT_CORPUS")
		if root == "" {
			root = oracleTestCorpusRoot(t)
		}
		if _, err := os.Stat(root); err != nil {
			t.Skipf("corpus not available (%v)", err)
		}
		return root
	})
}

func TestL11OptimizeOracle(t *testing.T) {
	dump := os.Getenv("L11OPT_ORACLE")
	if dump == "" {
		t.Skip("set L11OPT_ORACLE=<L11OptimizeDump output .gz> (and L11OPT_ROOT)")
	}
	root := os.Getenv("L11OPT_ROOT")
	if root == "" {
		root = oracleTestCorpusRoot(t)
	}
	l11optCompare(t, dump, false, func(string) string { return root })
}
