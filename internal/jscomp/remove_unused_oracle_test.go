package jscomp

// Differential test against tools/oracle/java/RemoveUnusedDump.java: the
// w5-remove-unused passes (RemoveUnusedCode, PureFunctionIdentifier,
// DeadAssignmentsElimination, CoalesceVariableNames, LiveVariablesAnalysis,
// OptimizeParameters, OptimizeReturns, OptimizeConstructors,
// DevirtualizeMethods, FunctionRewriter and the cross-chunk passes) run on
// parsed and normalized files, printed in the same format and compared pass by
// pass.
//
//   - TestRemoveUnusedGolden replays the checked-in golden dump
//     testdata/removeunused/golden.txt.gz (with --full trees): the hand-written
//     inputs testdata/removeunused/js/*.js, followed by the corpus files of
//     testdata/removeunused/golden-files.txt (skipped when the corpus is not
//     present). Regenerate with
//       CP=.cache/closure-compiler.jar:<RemoveUnusedDump build dir>
//       J="java -Xss1g -XX:-ShowCodeDetailsInExceptionMessages -cp $CP"
//       $J com.google.javascript.jscomp.RemoveUnusedDump --full a.gz \
//         internal/jscomp/testdata/removeunused/js <the *.js names>
//       $J com.google.javascript.jscomp.RemoveUnusedDump --full b.gz \
//         .cache/corpus --list=internal/jscomp/testdata/removeunused/golden-files.txt
//       cat a.gz b.gz > internal/jscomp/testdata/removeunused/golden.txt.gz
//   - TestRemoveUnusedOracle compares any other dump (hashes, or --full):
//     REMOVE_UNUSED_ORACLE=<dump.gz> [REMOVE_UNUSED_ROOT=<root>] (default
//     root: the corpus under /home/user/jsgo/.cache/corpus).

import (
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// removeUnusedOracleExterns is RemoveUnusedDump.EXTERNS.
var removeUnusedOracleExterns = strings.Join([]string{
	"var window;",
	"var document;",
	"var console;",
	"var Arguments;",
	"/** @constructor @param {*=} opt_value @return {!Object} */",
	"function Object(opt_value) {}",
	"/** @type {Function} */ Object.defineProperties = function() {};",
	"/** @type {Function} */ Object.defineProperty = function() {};",
	"/** @constructor @param {...*} var_args @return {!Array} */",
	"function Array(var_args) {}",
	"/** @constructor @param {*=} message */ function Error(message) {}",
	"/** @nosideeffects */ function externPure() {}",
	"function alert(a) {}",
	"",
}, "\n")

type removeUnusedOraclePass struct {
	name string
	run  func(c *Compiler, externs *rhino.Node, js *rhino.Node, extra *[]string)
}

// removeUnusedOracleRucSimple is RemoveUnusedDump.rucSimple.
func removeUnusedOracleRucSimple(c *Compiler) *RemoveUnusedCode_Builder {
	return NewRemoveUnusedCode_Builder(c).
		RemoveUnusedPrototypeProperties(false).
		RemoveLocalVars(true).
		RemoveGlobals(false).
		PreserveFunctionExpressionNames(false).
		RemoveUnusedThisProperties(false).
		RemoveUnusedObjectDefinePropertiesDefinitions(false).
		AssumeGettersArePure(true).
		RemoveUnusedPolyfills(false)
}

// removeUnusedOracleRucAdvanced is RemoveUnusedDump.rucAdvanced.
func removeUnusedOracleRucAdvanced(c *Compiler) *RemoveUnusedCode_Builder {
	return NewRemoveUnusedCode_Builder(c).
		RemoveUnusedPrototypeProperties(true).
		RemoveLocalVars(true).
		RemoveGlobals(true).
		PreserveFunctionExpressionNames(false).
		RemoveUnusedThisProperties(true).
		RemoveUnusedObjectDefinePropertiesDefinitions(true).
		AssumeGettersArePure(true).
		RemoveUnusedPolyfills(false)
}

// removeUnusedOraclePasses is RemoveUnusedDump.passes.
func removeUnusedOraclePasses() []removeUnusedOraclePass {
	return []removeUnusedOraclePass{
		{"ruc-simple", func(c *Compiler, ext, js *rhino.Node, extra *[]string) {
			removeUnusedOracleRucSimple(c).Build().Process(ext, js)
		}},
		{"ruc-advanced", func(c *Compiler, ext, js *rhino.Node, extra *[]string) {
			removeUnusedOracleRucAdvanced(c).Build().Process(ext, js)
		}},
		{"pfi-ruc-advanced", func(c *Compiler, ext, js *rhino.Node, extra *[]string) {
			NewPureFunctionIdentifier_Driver(c).Process(ext, js)
			removeUnusedOracleRucAdvanced(c).Build().Process(ext, js)
		}},
		{"pfi", func(c *Compiler, ext, js *rhino.Node, extra *[]string) {
			NewPureFunctionIdentifier_Driver(c).Process(ext, js)
			removeUnusedOracleCollectFlags(ext, extra)
			removeUnusedOracleCollectFlags(js, extra)
		}},
		{"dae", func(c *Compiler, ext, js *rhino.Node, extra *[]string) {
			NewDeadAssignmentsElimination(c).Process(ext, js)
		}},
		{"cvn", func(c *Compiler, ext, js *rhino.Node, extra *[]string) {
			NewCoalesceVariableNames(c, false).Process(ext, js)
		}},
		{"cvn-pseudo", func(c *Compiler, ext, js *rhino.Node, extra *[]string) {
			NewCoalesceVariableNames(c, true).Process(ext, js)
		}},
		{"lva", func(c *Compiler, ext, js *rhino.Node, extra *[]string) {
			NodeTraversal_Traverse(c, js, newRemoveUnusedOracleLva(c, extra))
		}},
		{"optparams", func(c *Compiler, ext, js *rhino.Node, extra *[]string) {
			NewOptimizeParameters(c).Process(ext, js)
		}},
		{"optreturns", func(c *Compiler, ext, js *rhino.Node, extra *[]string) {
			NewOptimizeReturns(c).Process(ext, js)
		}},
		{"optconstructors", func(c *Compiler, ext, js *rhino.Node, extra *[]string) {
			NewOptimizeConstructors(c).Process(ext, js)
		}},
		{"devirtualize", func(c *Compiler, ext, js *rhino.Node, extra *[]string) {
			OptimizeCalls_Builder_().
				SetCompiler(c).
				SetConsiderExterns(false).
				AddPass(NewDevirtualizeMethods(c)).
				Build().
				Process(ext, js)
		}},
		{"rewrite-functions", func(c *Compiler, ext, js *rhino.Node, extra *[]string) {
			NewFunctionRewriter(c).Process(ext, js)
		}},
	}
}

// removeUnusedOracleChunkPasses is RemoveUnusedDump.chunkPasses.
func removeUnusedOracleChunkPasses() []removeUnusedOraclePass {
	return []removeUnusedOraclePass{
		{"crosschunk-refs", func(c *Compiler, ext, js *rhino.Node, extra *[]string) {
			collector := NewCrossChunkReferenceCollector(c, NewSyntacticScopeCreator(c))
			collector.ProcessRoot(js)
			for _, s := range collector.GetTopLevelStatements() {
				var l strings.Builder
				l.WriteString("stmt " + strconv.Itoa(s.GetOriginalOrder()) + " " + s.GetChunk().GetName() +
					" " + removeUnusedOraclePos(s.GetStatementNode()))
				if s.IsDeclarationStatement() {
					value := "-"
					if s.GetDeclaredValueNode() != nil {
						value = removeUnusedOraclePos(s.GetDeclaredValueNode())
					}
					l.WriteString(" decl " + removeUnusedOraclePos(s.GetDeclaredNameReference().GetNode()) +
						" value " + value + " movable " + strconv.FormatBool(s.IsMovableDeclaration()))
				}
				l.WriteString(" refs")
				for _, r := range s.GetNonDeclarationReferences() {
					l.WriteString(" " + removeUnusedOraclePos(r.GetNode()))
				}
				*extra = append(*extra, l.String())
			}
			for name, v := range collector.GetGlobalVariableNamesMap().All() {
				refs := collector.GetReferences(v)
				var l strings.Builder
				l.WriteString("var " + name)
				if refs == nil {
					l.WriteString(" none")
				} else {
					l.WriteString(" wellDefined " + strconv.FormatBool(refs.IsWellDefined()) +
						" assignedOnce " + strconv.FormatBool(refs.IsAssignedOnceInLifetime()) + " refs")
					for r := range refs.All() {
						l.WriteString(" " + removeUnusedOraclePos(r.GetNode()))
					}
				}
				*extra = append(*extra, l.String())
			}
		}},
		{"crosschunk-code", func(c *Compiler, ext, js *rhino.Node, extra *[]string) {
			NewCrossChunkCodeMotion(c, c.GetChunkGraph(), false).Process(ext, js)
		}},
		{"crosschunk-method", func(c *Compiler, ext, js *rhino.Node, extra *[]string) {
			NewCrossChunkMethodMotion(c, NewIdGenerator(), false, false).Process(ext, js)
		}},
	}
}

// removeUnusedOracleLva is RemoveUnusedDump.Lva: LiveVariablesAnalysis for
// every function, as DeadAssignmentsElimination computes it.
type removeUnusedOracleLva struct {
	NodeTraversal_AbstractCfgCallbackBase
	compiler *Compiler
	extra    *[]string
}

func newRemoveUnusedOracleLva(compiler *Compiler, extra *[]string) *removeUnusedOracleLva {
	l := &removeUnusedOracleLva{compiler: compiler, extra: extra}
	l.InitNodeTraversal_AbstractCfgCallbackBase(l)
	return l
}

func (l *removeUnusedOracleLva) ExitScopeWithCfg(t *NodeTraversal) {
	if !t.InFunctionBlockScope() {
		return
	}
	blockScope := t.GetScope()
	functionScope := blockScope.GetParent()
	function := t.GetScopeRoot().GetParent()
	cfg := l.GetControlFlowGraph(l.compiler)
	scopeCreator := NewSyntacticScopeCreator(l.compiler)
	allVarsDeclaredInFunction := NodeUtil_GetAllVarsDeclaredInFunction(l.compiler, scopeCreator, functionScope)
	liveness := NewLiveVariablesAnalysis(
		cfg, functionScope, blockScope, l.compiler, scopeCreator, allVarsDeclaredInFunction)
	liveness.Analyze()
	var escaped []string
	for v := range liveness.GetEscapedLocals().All() {
		escaped = append(escaped, v.GetName())
	}
	slices.SortFunc(escaped, jutil.StringCompareTo) // Collections.sort: String.compareTo
	*l.extra = append(*l.extra, "function "+removeUnusedOraclePos(function)+" "+
		function.GetFirstChild().GetString()+" vars "+strconv.Itoa(len(liveness.GetAllVariablesInOrder()))+
		" escaped "+strings.Join(escaped, ","))
	for _, v := range liveness.GetAllVariablesInOrder() {
		*l.extra = append(*l.extra, "var "+strconv.Itoa(liveness.GetVarIndex(v.GetName()))+" "+v.GetName())
	}
	for _, cfgNode := range cfg.GetDiGraphNodes() {
		state := dataFlowAnalysis_linearFlowStateOf[*LiveVariablesAnalysis_LiveVariableLattice](cfgNode)
		n := cfgNode.GetValue()
		where := "IMPLICIT_RETURN"
		if n != nil {
			where = removeUnusedOraclePos(n)
		}
		*l.extra = append(*l.extra, "node "+where+" in "+removeUnusedOracleBits(state.GetIn())+
			" out "+removeUnusedOracleBits(state.GetOut()))
	}
}

// removeUnusedOracleBits is RemoveUnusedDump.bits.
func removeUnusedOracleBits(l *LiveVariablesAnalysis_LiveVariableLattice) string {
	var sb strings.Builder
	sb.WriteByte('{')
	first := true
	for i := l.NextSetBit(0); i >= 0; i = l.NextSetBit(i + 1) {
		if !first {
			sb.WriteByte(',')
		}
		first = false
		sb.WriteString(strconv.Itoa(i))
	}
	sb.WriteByte('}')
	return sb.String()
}

// removeUnusedOracleCollectFlags is RemoveUnusedDump.collectFlags.
func removeUnusedOracleCollectFlags(n *rhino.Node, extra *[]string) {
	if n.IsCall() || n.IsNew() || n.IsTaggedTemplateLit() || n.IsOptChainCall() {
		*extra = append(*extra, "flags "+removeUnusedOraclePos(n)+" "+strconv.Itoa(int(n.GetSideEffectFlags())))
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		removeUnusedOracleCollectFlags(c, extra)
	}
}

func removeUnusedOraclePos(n *rhino.Node) string {
	return n.GetToken().String() + "@" + strconv.Itoa(n.GetLineno()) + ":" + strconv.Itoa(n.GetCharno())
}

func removeUnusedOracleOptions() *CompilerOptions {
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	options.SetLanguageOut(CompilerOptions_LanguageMode_NO_TRANSPILE)
	// VarCheck declares the undeclared names (as the checks phase does before these passes run);
	// its "undefined variable" errors would only repeat in every run.
	options.SetWarningLevel(DiagnosticGroups_UNDEFINED_VARIABLES, CheckLevel_OFF)
	return options
}

// removeUnusedOracleNewCompiler is RemoveUnusedDump.newCompiler.
func removeUnusedOracleNewCompiler(rel string, src string) *Compiler {
	compiler := NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
	compiler.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", removeUnusedOracleExterns)},
		[]*SourceFile{SourceFile_FromCode(rel, src)},
		removeUnusedOracleOptions())
	compiler.Parse()
	return compiler
}

// removeUnusedOracleNewChunkCompiler is RemoveUnusedDump.newChunkCompiler.
func removeUnusedOracleNewChunkCompiler(parts []string) *Compiler {
	compiler := NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
	var chunks []*JSChunk
	for k, part := range parts {
		chunk := NewJSChunk("m" + strconv.Itoa(k+1))
		chunk.AddSourceFile(SourceFile_FromCode("m"+strconv.Itoa(k+1)+".js", part))
		if k > 0 {
			chunk.AddDependency(chunks[k-1])
		}
		chunks = append(chunks, chunk)
	}
	compiler.InitChunks(
		[]*SourceFile{SourceFile_FromCode("externs.js", removeUnusedOracleExterns)}, chunks,
		removeUnusedOracleOptions())
	compiler.Parse()
	return compiler
}

// removeUnusedOraclePrepare is RemoveUnusedDump.prepare.
func removeUnusedOraclePrepare(c *Compiler) {
	externs := c.GetExternsRoot()
	js := c.GetJsRoot()
	NewVarCheck(c).Process(externs, js)
	Normalize_CreateNormalizeForOptimizations(c).Process(externs, js)
	NewGatherGetterAndSetterProperties(c).Process(externs, js)
	NewGatherExternProperties(c, GatherExternProperties_Mode_OPTIMIZE).Process(externs, js)
}

// removeUnusedOracleDropChangeTime is RemoveUnusedDump.dropChangeTime.
func removeUnusedOracleDropChangeTime(l string) string {
	const key = " [change_time: "
	i := strings.Index(l, key)
	if i < 0 {
		return l
	}
	end := i + strings.IndexByte(l[i:], ']')
	return l[:i] + " [change_time]" + l[end+1:]
}

// removeUnusedOracleCollectChanges is RemoveUnusedDump.collectChanges.
func removeUnusedOracleCollectChanges(n *rhino.Node, out *[]string) {
	if (n.IsScript() || n.IsFunction()) && n.GetChangeTime() != 0 {
		*out = append(*out, removeUnusedOraclePos(n))
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		removeUnusedOracleCollectChanges(c, out)
	}
}

// removeUnusedOracleChunkOffsets is RemoveUnusedDump.chunkOffsets.
func removeUnusedOracleChunkOffsets(parsed *Compiler) (offsets []int, skip string) {
	script := parsed.GetJsRoot().GetFirstChild()
	if script.HasChildren() && script.GetFirstChild().IsModuleBody() {
		return nil, "module"
	}
	var statements []*rhino.Node
	for c := script.GetFirstChild(); c != nil; c = c.GetNext() {
		statements = append(statements, c)
	}
	n := len(statements)
	if n < 3 {
		return nil, "small"
	}
	for k := 1; k <= 2; k++ {
		last := statements[k*n/3-1]
		offsets = append(offsets, last.GetSourceOffset()+last.GetLength())
	}
	return offsets, ""
}

type removeUnusedOracleDumper struct {
	full bool
	sb   strings.Builder
}

func (d *removeUnusedOracleDumper) line(s string) {
	d.sb.WriteString(normalizeOracleEscape(s))
	d.sb.WriteByte('\n')
}

// emitBlock is RemoveUnusedDump.emitBlock.
func (d *removeUnusedOracleDumper) emitBlock(kind string, raw []string) {
	var escaped []string
	for _, l := range raw {
		escaped = append(escaped, normalizeOracleEscape(l))
	}
	d.line(kind + " " + strconv.Itoa(len(raw)) + " " + normalizeOracleSha16(escaped))
	if d.full {
		for _, l := range escaped {
			d.sb.WriteString(l + "\n")
		}
	}
}

// runPass is RemoveUnusedDump.runPass.
func (d *removeUnusedOracleDumper) runPass(pass removeUnusedOraclePass, c *Compiler, parseDiags []string) {
	d.line("## pass " + pass.name)
	var extra []string
	before := ""
	exception := ""
	func() {
		defer func() {
			if r := recover(); r != nil {
				exception = "exception " + normalizeOracleException(r)
			}
		}()
		removeUnusedOraclePrepare(c)
		before = c.GetRoot().ToStringTree()
		pass.run(c, c.GetExternsRoot(), c.GetJsRoot(), &extra)
	}()
	diags := normalizeOracleDiagnostics(c)
	for _, diag := range parseDiags {
		diags = normalizeOracleRemoveFirst(diags, diag)
	}
	for _, diag := range diags {
		d.line(diag)
	}
	if exception != "" {
		d.line(exception)
		return
	}
	after := c.GetRoot().ToStringTree()
	if after == before {
		d.line("tree unchanged")
	} else {
		var lines []string
		for _, l := range strings.Split(after, "\n") {
			lines = append(lines, normalizeOracleEscape(removeUnusedOracleDropChangeTime(l)))
		}
		if lines[len(lines)-1] == "" {
			lines = lines[:len(lines)-1]
		}
		d.line("tree changed " + normalizeOracleSha16(lines))
		if d.full {
			for _, l := range lines {
				d.sb.WriteString(l + "\n")
			}
		}
	}
	var changes []string
	removeUnusedOracleCollectChanges(c.GetRoot(), &changes)
	d.emitBlock("changes", changes)
	d.emitBlock("extra", extra)
}

// dumpFile is RemoveUnusedDump.dumpFile.
func (d *removeUnusedOracleDumper) dumpFile(rel string, src string) {
	d.line("## file " + rel)
	compiler := removeUnusedOracleNewCompiler(rel, src)
	parseDiags := normalizeOracleDiagnostics(compiler)
	if compiler.HasErrors() {
		d.line("# parse errors")
	} else {
		d.line("# parse ok")
	}
	for _, diag := range parseDiags {
		d.line(diag)
	}
	if compiler.HasErrors() {
		return
	}
	for _, pass := range removeUnusedOraclePasses() {
		d.runPass(pass, removeUnusedOracleNewCompiler(rel, src), parseDiags)
	}
	offsets, skip := removeUnusedOracleChunkOffsets(compiler)
	if offsets == nil {
		d.line("## chunks skipped " + skip)
		return
	}
	units := jutil.StringToUTF16(src)
	if !(0 <= offsets[0] && offsets[0] <= offsets[1] && offsets[1] <= len(units)) {
		d.line("## chunks skipped offsets " + strconv.Itoa(offsets[0]) + " " + strconv.Itoa(offsets[1]))
		return
	}
	d.line("## chunks " + strconv.Itoa(offsets[0]) + " " + strconv.Itoa(offsets[1]))
	parts := []string{
		jutil.StringFromUTF16(units[:offsets[0]]),
		jutil.StringFromUTF16(units[offsets[0]:offsets[1]]),
		jutil.StringFromUTF16(units[offsets[1]:]),
	}
	probe := removeUnusedOracleNewChunkCompiler(parts)
	chunkParseDiags := normalizeOracleDiagnostics(probe)
	if probe.HasErrors() {
		d.line("# chunks parse errors")
	} else {
		d.line("# chunks parse ok")
	}
	for _, diag := range chunkParseDiags {
		d.line(diag)
	}
	if probe.HasErrors() {
		return
	}
	for _, pass := range removeUnusedOracleChunkPasses() {
		d.runPass(pass, removeUnusedOracleNewChunkCompiler(parts), chunkParseDiags)
	}
}

// removeUnusedOracleSections splits one file's dump into its "## pass" and
// "## chunks" sections (the header lines before the first one included in
// the first section).
func removeUnusedOracleSections(dump string) []string {
	var out []string
	var cur strings.Builder
	for _, l := range strings.SplitAfter(dump, "\n") {
		if (strings.HasPrefix(l, "## pass ") || strings.HasPrefix(l, "## chunks")) && cur.Len() > 0 {
			out = append(out, cur.String())
			cur.Reset()
		}
		cur.WriteString(l)
	}
	if cur.Len() > 0 {
		out = append(out, cur.String())
	}
	return out
}

// removeUnusedOracleCompare dumps every file of the dump at path (sources
// under rootFor(name)) and compares section by section; it returns the
// number of files compared.
func removeUnusedOracleCompare(t *testing.T, path string, rootFor func(rel string) string, only func(rel string) bool) int {
	t.Helper()
	names, sections, full := normalizeOracleSplit(t, path)
	n, failures := 0, 0
	for _, name := range names {
		rel := normalizeOracleUnescapeName(name)
		if only != nil && !only(rel) {
			continue
		}
		src, err := os.ReadFile(filepath.Join(rootFor(rel), rel))
		if err != nil {
			t.Logf("skipping %s: %v", rel, err)
			continue
		}
		d := &removeUnusedOracleDumper{full: full}
		func() {
			defer func() {
				if r := recover(); r != nil {
					d.line("go panic " + normalizeOracleException(r))
				}
			}()
			d.dumpFile(rel, jutil.NewStringFromUTF8Bytes(src))
		}()
		n++
		got, want := removeUnusedOracleSections(d.sb.String()), removeUnusedOracleSections(sections[name])
		differs := len(got) != len(want)
		for i := 0; !differs && i < len(got); i++ {
			if got[i] == want[i] {
				continue
			}
			differs = true
			if failures < 10 {
				t.Errorf("%s: dump differs in %s:\n%s", rel, strings.SplitN(want[i], "\n", 2)[0],
					removeUnusedOracleDiff(got[i], want[i]))
			}
		}
		if differs {
			if len(got) != len(want) && failures < 10 {
				t.Errorf("%s: dump differs:\n%s", rel, normalizeOracleFirstDiff(d.sb.String(), sections[name]))
			}
			failures++
		}
	}
	if failures > 0 {
		t.Errorf("%d of %d files differ", failures, n)
	}
	return n
}

// removeUnusedOracleDiff shows the first differing line and, when that is a
// hash line ("tree changed", "changes", "extra") of a --full dump, the first
// differing line after it.
func removeUnusedOracleDiff(got, want string) string {
	g := strings.Split(got, "\n")
	w := strings.Split(want, "\n")
	var out []string
	for i := 0; (i < len(g) || i < len(w)) && len(out) < 2; i++ {
		var gl, wl string
		if i < len(g) {
			gl = g[i]
		}
		if i < len(w) {
			wl = w[i]
		}
		if gl == wl {
			continue
		}
		out = append(out, "line "+strconv.Itoa(i+1)+":\n  go:   "+gl+"\n  java: "+wl)
		if !strings.HasPrefix(wl, "tree changed ") && !strings.HasPrefix(wl, "changes ") &&
			!strings.HasPrefix(wl, "extra ") {
			break
		}
	}
	if len(out) == 0 {
		return "(no line differs)"
	}
	return strings.Join(out, "\n")
}

func TestRemoveUnusedGolden(t *testing.T) {
	golden := "testdata/removeunused/golden.txt.gz"
	jsDir := "testdata/removeunused/js"
	hand := map[string]bool{}
	entries, err := os.ReadDir(jsDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		hand[e.Name()] = true
	}
	handRoot := func(string) string { return jsDir }
	if n := removeUnusedOracleCompare(t, golden, handRoot, func(rel string) bool { return hand[rel] }); n != len(hand) {
		t.Errorf("compared %d hand-written files, want %d", n, len(hand))
	}
	corpus := "/home/user/jsgo/.cache/corpus"
	if _, err := os.Stat(corpus); err != nil {
		t.Logf("corpus not present; golden corpus files skipped")
		return
	}
	removeUnusedOracleCompare(t, golden, func(string) string { return corpus }, func(rel string) bool { return !hand[rel] })
}

func TestRemoveUnusedOracle(t *testing.T) {
	dump := os.Getenv("REMOVE_UNUSED_ORACLE")
	if dump == "" {
		t.Skip("set REMOVE_UNUSED_ORACLE=<RemoveUnusedDump output .gz> (and REMOVE_UNUSED_ROOT) to run")
	}
	root := os.Getenv("REMOVE_UNUSED_ROOT")
	if root == "" {
		root = "/home/user/jsgo/.cache/corpus"
	}
	n := removeUnusedOracleCompare(t, dump, func(string) string { return root }, nil)
	t.Logf("compared %d files", n)
}
