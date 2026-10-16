package jscomp

// Differential test against tools/oracle/java/InfraDump.java (Go-only:
// test-helper): the same sweep of AstValidator (on parsed and on
// deliberately corrupted SCRIPTs), MinimizedCondition, NodeIterators,
// NodeNameExtractor, TemplateAstMatcher, DataFlowAnalysis and
// DataFlowAnalysis.computeEscaped over parsed files, printed in the same
// format and compared file by file.
//
//   - TestInfraDumpGolden replays the checked-in golden dump
//     testdata/infra/golden.txt.gz: the hand-written inputs
//     testdata/infra/js/*.js, then (when /home/user/jsgo/.cache/corpus
//     exists) the corpus files listed in testdata/infra/golden-files.txt
//     (every 400th file of corpus-files.txt). Regenerate with
//       CP=.cache/closure-compiler.jar:.cache/oracle-build
//       java -Xss1g -cp $CP com.google.javascript.jscomp.InfraDump a.gz \
//         internal/jscomp/testdata/infra/js <the js/*.js names>
//       java -Xss1g -cp $CP com.google.javascript.jscomp.InfraDump b.gz \
//         .cache/corpus --list=internal/jscomp/testdata/infra/golden-files.txt
//       cat a.gz b.gz > internal/jscomp/testdata/infra/golden.txt.gz
//   - TestInfraDumpOracle compares any other dump: INFRA_ORACLE=<dump.gz>
//     [INFRA_ROOT=<root>] (default root: the corpus under
//     /home/user/jsgo/.cache/corpus) [INFRA_DIFF_DIR=<dir>: writes the Go
//     dump of every differing file there]. The held-out corpus
//     (.cache/corpus-holdout, corpus-holdout-files.txt) is for a final
//     check only.

import (
	"bufio"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const infraCorpusRoot = "/home/user/jsgo/.cache/corpus"

// infraEscape is InfraDump.escape (ParsingDump's escaping).
func infraEscape(s string) string {
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

// infraExc is InfraDump.exc: toString() for IllegalState/IllegalArgument
// exceptions, the class name otherwise.
func infraExc(t jutil.Throwable) string {
	switch t.(type) {
	case jutil.IIllegalStateException, jutil.IIllegalArgumentException:
		return t.Error()
	}
	full := t.Error()
	if msg := t.GetMessage(); msg != "" && strings.HasSuffix(full, ": "+msg) {
		return strings.TrimSuffix(full, ": "+msg)
	}
	if i := strings.Index(full, ": "); i >= 0 {
		return full[:i]
	}
	return full
}

// infraShow is InfraDump.show.
func infraShow(n *rhino.Node) string {
	if n == nil {
		return "null"
	}
	return n.ToStringWithFlags(true, false, false)
}

// infraSexpr is InfraDump.sexpr.
func infraSexpr(n *rhino.Node) string {
	var b strings.Builder
	infraSexprTo(&b, n)
	return b.String()
}

func infraSexprTo(b *strings.Builder, n *rhino.Node) {
	if n == nil {
		b.WriteString("null")
		return
	}
	b.WriteByte('(')
	b.WriteString(n.ToStringWithFlags(false, false, false))
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		b.WriteByte(' ')
		infraSexprTo(b, c)
	}
	b.WriteByte(')')
}

func infraPreorder(n *rhino.Node, out *[]*rhino.Node) {
	*out = append(*out, n)
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		infraPreorder(c, out)
	}
}

func infraSize(n *rhino.Node, cap_ int) int {
	s := 1
	for c := n.GetFirstChild(); c != nil && s <= cap_; c = c.GetNext() {
		s += infraSize(c, cap_)
	}
	return s
}

// infraDumper is InfraDump.
type infraDumper struct {
	sb       strings.Builder
	compiler *Compiler
	ids      map[*rhino.Node]int
}

func (d *infraDumper) line(s string) {
	d.sb.WriteString(infraEscape(s))
	d.sb.WriteByte('\n')
}

func (d *infraDumper) id(n *rhino.Node) string {
	if n == nil {
		return "R"
	}
	if i, ok := d.ids[n]; ok {
		return strconv.Itoa(i)
	}
	return "?"
}

// infraTry runs body, returning the Java exception it threw (or nil).
func infraTry(body func()) (thrown jutil.Throwable) {
	jutil.Try(body, func(e jutil.IRuntimeException) { thrown = e })
	return thrown
}

// dumpOne is InfraDump.dumpOne.
func (d *infraDumper) dumpOne(name, src string) string {
	d.sb.Reset()
	d.line("## file " + name)
	d.compiler = NewCompiler()
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	d.compiler.InitOptions(options)
	var script *rhino.Node
	if e := infraTry(func() {
		script = NewCompilerInput(SourceFile_FromCode(name, src)).GetAstRoot(d.compiler)
	}); e != nil {
		d.line("# parse-exception " + infraExc(e))
		return d.sb.String()
	}
	d.line("# errors " + strconv.Itoa(d.compiler.GetErrorCount()))
	var nodes []*rhino.Node
	infraPreorder(script, &nodes)
	d.ids = make(map[*rhino.Node]int, len(nodes))
	for i, n := range nodes {
		d.ids[n] = i
	}
	d.line("# nodes " + strconv.Itoa(len(nodes)))

	d.validate(script, "")
	d.corrupt(script, len(nodes))
	d.minimize(nodes)
	d.iterators(nodes)
	d.names(nodes)
	d.templates(nodes)
	d.dataflow(script, nodes)
	d.escaped(script)
	return d.sb.String()
}

// ---------------------------------------------------------------------------
// AstValidator

type infraViolations struct {
	indent string
	lines  []string
}

func (v *infraViolations) HandleViolation(message string, n *rhino.Node) {
	v.lines = append(v.lines, v.indent+"V "+message+"\t"+infraShow(n))
}

func (d *infraDumper) validate(script *rhino.Node, indent string) {
	violations := &infraViolations{indent: indent}
	validator := NewAstValidatorWithHandler(
		d.compiler, violations /* validateScriptFeatures= */, true /* shouldValidateRequiredInlinings= */, false)
	validator.SetTypeValidationMode(AstValidator_TypeInfoValidation_NONE)
	e := infraTry(func() { validator.ValidateScript(script) })
	for _, v := range violations.lines {
		d.line(v)
	}
	if e != nil {
		d.line(indent + "X " + infraExc(e))
	}
}

const infraCorruptions = 6

func (d *infraDumper) corrupt(script *rhino.Node, count int) {
	if count < 2 {
		return
	}
	step := max(1, (count-1)/12)
	kind := 0
	for k := 1; k < count; k, kind = k+step, kind+1 {
		clone := script.CloneTree()
		var cloneNodes []*rhino.Node
		infraPreorder(clone, &cloneNodes)
		n := cloneNodes[k]
		var what string
		switch kind % infraCorruptions {
		case 0:
			if n.HasChildren() {
				n.GetFirstChild().Detach()
				what = "remove-first-child"
			} else {
				n.AddChildToBack(rhino.IR_Number(0))
				what = "add-child"
			}
		case 1:
			n.AddChildToBack(rhino.IR_Name("extra"))
			what = "add-name-child"
		case 2:
			if n.IsName() {
				n.SetToken(rhino.Token_STRINGLIT)
			} else {
				n.SetToken(rhino.Token_NAME)
			}
			what = "retoken"
		case 3:
			if next := n.GetNext(); next != nil {
				next.Detach()
				next.InsertBefore(n)
				what = "swap-with-next"
			} else {
				n.GetParent().AddChildToFront(rhino.IR_Empty())
				what = "add-empty-sibling"
			}
		case 4:
			n.ReplaceWith(rhino.IR_Empty())
			what = "replace-with-empty"
		default:
			block := rhino.NewNode(rhino.Token_BLOCK)
			n.ReplaceWith(block)
			block.AddChildToBack(n)
			what = "wrap-in-block"
		}
		d.line("C " + strconv.Itoa(k) + " " + what + " " + n.GetToken().String())
		d.validate(clone, "  ")
	}
}

// ---------------------------------------------------------------------------
// MinimizedCondition

func infraIsCondition(n *rhino.Node) bool {
	parent := n.GetParent()
	if parent == nil {
		return false
	}
	switch parent.GetToken() {
	case rhino.Token_IF, rhino.Token_WHILE, rhino.Token_DO, rhino.Token_FOR:
		return NodeUtil_GetConditionExpression(parent) == n && !n.IsEmpty()
	default:
		return false
	}
}

func (d *infraDumper) minimize(nodes []*rhino.Node) {
	for _, n := range nodes {
		var candidate bool
		switch n.GetToken() {
		case rhino.Token_NOT, rhino.Token_AND, rhino.Token_OR, rhino.Token_HOOK, rhino.Token_COMMA:
			candidate = true
		default:
			candidate = infraIsCondition(n)
		}
		if !candidate {
			continue
		}
		for _, style := range MinimizedCondition_MinimizationStyle_values() {
			clone := n.CloneTree()
			holder := rhino.NewNodeWithChild(rhino.Token_EXPR_RESULT, clone)
			var result string
			if e := infraTry(func() {
				mc := MinimizedCondition_FromConditionNode(clone)
				m := mc.GetMinimized(style)
				will := m.WillChange(clone)
				result = "len=" + strconv.Itoa(m.length) +
					" changed=" + strconv.FormatBool(m.changed) +
					" will=" + strconv.FormatBool(will)
				if will {
					m.ApplyTo(clone)
					result += " " + infraSexpr(holder.GetFirstChild())
				}
			}); e != nil {
				result = "X " + infraExc(e)
			}
			d.line("M " + d.id(n) + " " + n.GetToken().String() + " " + style.String() + ": " + result)
		}
	}
}

// ---------------------------------------------------------------------------
// NodeIterators

const infraMaxMotion = 40

func (d *infraDumper) motion(header string, it *NodeIterators_LocalVarMotion) {
	var b strings.Builder
	b.WriteString(header)
	b.WriteByte(':')
	count := 0
	for it.HasNext() {
		next := it.Next()
		if count == infraMaxMotion {
			b.WriteString(" ...")
			break
		}
		count++
		b.WriteString(" " + next.GetToken().String() + "@" + d.id(next))
	}
	d.line(b.String())
}

func (d *infraDumper) iterators(nodes []*rhino.Node) {
	for _, n := range nodes {
		if NodeUtil_IsNameDeclaration(n) && NodeUtil_IsStatement(n) {
			for name := n.GetFirstChild(); name != nil; name = name.GetNext() {
				if !name.IsName() {
					continue
				}
				header := "I var " + d.id(name)
				// motion prints its line only once the iteration is done.
				if e := infraTry(func() {
					d.motion(header, NodeIterators_LocalVarMotion_ForVar(d.compiler, name, n, n.GetParent()))
				}); e != nil {
					d.line(header + ": X " + infraExc(e))
				}
			}
		} else if n.IsExprResult() && n.GetFirstChild().IsAssign() && n.GetFirstFirstChild().IsName() {
			assign := n.GetFirstChild()
			header := "I assign " + d.id(assign.GetFirstChild())
			if e := infraTry(func() {
				d.motion(header, NodeIterators_LocalVarMotion_ForAssign(
					d.compiler, assign.GetFirstChild(), assign, n, n.GetParent()))
			}); e != nil {
				d.line(header + ": X " + infraExc(e))
			}
		}
	}
}

// ---------------------------------------------------------------------------
// NodeNameExtractor

func (d *infraDumper) names(nodes []*rhino.Node) {
	extractor := NewNodeNameExtractor('$')
	for i := 1; i < len(nodes); i += 3 {
		n := nodes[i]
		if infraSize(n, 25) > 24 {
			continue
		}
		var result string
		if e := infraTry(func() { result = jutil.ToJavaString(extractor.GetName(n)) }); e != nil {
			result = "X " + infraExc(e)
		}
		d.line("N " + strconv.Itoa(i) + " " + n.GetToken().String() + ": " + result)
	}
}

// ---------------------------------------------------------------------------
// TemplateAstMatcher

// infraTemplates is InfraDump.TEMPLATES.
var infraTemplates = []string{
	"function template() { 3; }",
	"/** @param {?} a */ function template(a) { a.length; }",
	"/** @param {?} a\n * @param {?} b */ function template(a, b) { a.indexOf(b) != -1; }",
	"function template() { var x = 1; x++; }",
	"/** @param {string} string_literal_s */\n" +
		"function template(string_literal_s) { goog.require(string_literal_s); }",
	"/** @param {number} n */ function template(n) { Math.floor(n); }",
	"/** @param {?} a\n * @param {?} b */ function template(a, b) { a = b; }",
}

func (d *infraDumper) templates(nodes []*rhino.Node) {
	for t, template := range infraTemplates {
		var matcher, exact *TemplateAstMatcher
		if e := infraTry(func() {
			templateScript := NewCompilerInput(SourceFile_FromCode("template"+strconv.Itoa(t), template)).
				GetAstRoot(d.compiler)
			fn := templateScript.GetFirstChild()
			matcher = NewTemplateAstMatcher(d.compiler, fn, TypeMatchingStrategy_LOOSE)
			if t == 1 {
				exact = NewTemplateAstMatcher(d.compiler, fn, TypeMatchingStrategy_EXACT)
			}
		}); e != nil {
			d.line("T " + strconv.Itoa(t) + " X " + infraExc(e))
			continue
		}
		for _, n := range nodes {
			d.match("T "+strconv.Itoa(t), matcher, n)
			if exact != nil {
				d.match("T "+strconv.Itoa(t)+"e", exact, n)
			}
		}
	}
}

func (d *infraDumper) match(tag string, matcher *TemplateAstMatcher, n *rhino.Node) {
	var matches bool
	if e := infraTry(func() { matches = matcher.Matches(n) }); e != nil {
		d.line(tag + " " + d.id(n) + ": X " + infraExc(e))
		return
	}
	if !matches {
		return
	}
	var b strings.Builder
	b.WriteString(tag + " " + d.id(n) + ": loose=" + strconv.FormatBool(matcher.IsLooseMatch()) + " {")
	first := true
	for k, v := range matcher.GetTemplateNodeToMatchMap().All() {
		if !first {
			b.WriteString(", ")
		}
		first = false
		b.WriteString(k + "=" + infraSexpr(v))
	}
	b.WriteByte('}')
	d.line(b.String())
}

// ---------------------------------------------------------------------------
// DataFlowAnalysis

// infraNameSet is InfraDump.NameSet: a set of names, compared by value.
type infraNameSet struct {
	names *jutil.TreeSet[string]
}

func newInfraNameSet() *infraNameSet {
	return &infraNameSet{names: jutil.NewTreeSet(jutil.StringCompareTo)}
}

func newInfraNameSetCopy(other *infraNameSet) *infraNameSet {
	s := newInfraNameSet()
	for n := range other.names.All() {
		s.names.Add(n)
	}
	return s
}

func (s *infraNameSet) Equals(o any) bool {
	other, ok := o.(*infraNameSet)
	if !ok || other == nil || other.names.Size() != s.names.Size() {
		return false
	}
	for n := range s.names.All() {
		if !other.names.Contains(n) {
			return false
		}
	}
	return true
}

func (s *infraNameSet) HashCode() int32 {
	var h int32
	for n := range s.names.All() {
		h += jutil.StringHashCode(n)
	}
	return h
}

func (s *infraNameSet) String() string { return s.names.String() }

// infraUnionJoiner is InfraDump.UnionJoiner.
type infraUnionJoiner struct{ result *infraNameSet }

func (j *infraUnionJoiner) JoinFlow(input *infraNameSet) {
	for n := range input.names.All() {
		j.result.names.Add(n)
	}
}

func (j *infraUnionJoiner) Finish() *infraNameSet { return j.result }

// infraNames is InfraDump.names(Node, boolean, Set<String>).
func infraNames(n *rhino.Node, writes bool, out *jutil.TreeSet[string]) {
	if n.IsName() {
		parent := n.GetParent()
		isWrite := parent != nil &&
			(NodeUtil_IsNameDeclaration(parent) || (parent.IsAssign() && parent.GetFirstChild() == n))
		if isWrite == writes {
			out.Add(n.GetString())
		}
	}
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		switch c.GetToken() {
		case rhino.Token_BLOCK, rhino.Token_FUNCTION, rhino.Token_CLASS, rhino.Token_SCRIPT:
		default:
			infraNames(c, writes, out)
		}
	}
}

// infraWrittenNames is InfraDump.WrittenNames.
type infraWrittenNames struct {
	DataFlowAnalysisBase[*rhino.Node, *infraNameSet]
	branched bool
}

func newInfraWrittenNames(cfg ControlFlowGraph[*rhino.Node], branched bool) *infraWrittenNames {
	a := &infraWrittenNames{branched: branched}
	a.InitDataFlowAnalysisBase(a, cfg)
	return a
}

func (a *infraWrittenNames) IsForward() bool  { return true }
func (a *infraWrittenNames) IsBranched() bool { return a.branched }
func (a *infraWrittenNames) CreateFlowJoiner() DataFlowAnalysis_FlowJoiner[*infraNameSet] {
	return &infraUnionJoiner{result: newInfraNameSet()}
}
func (a *infraWrittenNames) FlowThrough(node *rhino.Node, input *infraNameSet) *infraNameSet {
	out := newInfraNameSetCopy(input)
	infraNames(node, true, out.names)
	return out
}
func (a *infraWrittenNames) CreateFlowBrancher(node *rhino.Node, output *infraNameSet) DataFlowAnalysis_FlowBrancher[*infraNameSet] {
	return func(branch ControlFlowGraph_Branch) *infraNameSet {
		if !branch.IsConditional() {
			return output
		}
		result := newInfraNameSetCopy(output)
		result.names.Add("@" + branch.String())
		return result
	}
}
func (a *infraWrittenNames) CreateInitialEstimateLattice() *infraNameSet { return newInfraNameSet() }
func (a *infraWrittenNames) CreateEntryLattice() *infraNameSet           { return newInfraNameSet() }

// infraReadNames is InfraDump.ReadNames.
type infraReadNames struct {
	DataFlowAnalysisBase[*rhino.Node, *infraNameSet]
}

func newInfraReadNames(cfg ControlFlowGraph[*rhino.Node]) *infraReadNames {
	a := &infraReadNames{}
	a.InitDataFlowAnalysisBase(a, cfg)
	return a
}

func (a *infraReadNames) IsForward() bool { return false }
func (a *infraReadNames) CreateFlowJoiner() DataFlowAnalysis_FlowJoiner[*infraNameSet] {
	return &infraUnionJoiner{result: newInfraNameSet()}
}
func (a *infraReadNames) FlowThrough(node *rhino.Node, input *infraNameSet) *infraNameSet {
	out := newInfraNameSetCopy(input)
	infraNames(node, false, out.names)
	return out
}
func (a *infraReadNames) CreateInitialEstimateLattice() *infraNameSet { return newInfraNameSet() }
func (a *infraReadNames) CreateEntryLattice() *infraNameSet           { return newInfraNameSet() }

func (d *infraDumper) dataflow(script *rhino.Node, nodes []*rhino.Node) {
	roots := []*rhino.Node{script}
	for _, n := range nodes {
		if n.IsFunction() {
			roots = append(roots, n)
		}
	}
	for _, root := range roots {
		var cfg ControlFlowGraph[*rhino.Node]
		if e := infraTry(func() {
			cfg = ControlFlowAnalysis_Builder_().
				SetCompiler(d.compiler).
				SetCfgRoot(root).
				SetIncludeEdgeAnnotations(true).
				ComputeCfg()
		}); e != nil {
			d.line("D " + d.id(root) + " cfg X " + infraExc(e))
			continue
		}
		d.runAnalysis("D "+d.id(root)+" fwd", cfg, newInfraWrittenNames(cfg, false), false)
		d.runAnalysis("D "+d.id(root)+" br", cfg, newInfraWrittenNames(cfg, true), true)
		d.runAnalysis("D "+d.id(root)+" bwd", cfg, newInfraReadNames(cfg), false)
	}
}

func (d *infraDumper) runAnalysis(header string, cfg ControlFlowGraph[*rhino.Node], a DataFlowAnalysis[*rhino.Node, *infraNameSet], edges bool) {
	if e := infraTry(func() { a.GetDataFlowAnalysisBase().Analyze() }); e != nil {
		d.line(header + " X " + infraExc(e))
		return
	}
	d.line(header)
	for _, node := range cfg.GetDiGraphNodes() {
		state := node.GetAnnotation().(*DataFlowAnalysis_LinearFlowState[*infraNameSet])
		d.line("  " + d.id(node.GetValue()) +
			" in=" + state.GetIn().String() +
			" out=" + state.GetOut().String() +
			" steps=" + strconv.Itoa(state.GetStepCount()))
		if edges {
			for _, e := range node.GetOutEdges() {
				d.line("    -> " + d.id(e.GetDestination().GetValue()) + " " + e.GetValue().String() +
					" " + jutil.ToJavaString(e.GetAnnotation()))
			}
		}
	}
}

// infraEscapedFinder is InfraDump.escaped's anonymous
// AbstractScopedCallback.
type infraEscapedFinder struct {
	NodeTraversal_AbstractScopedCallbackBase
	d            *infraDumper
	scopeCreator *SyntacticScopeCreator
	lines        []string
}

func (f *infraEscapedFinder) EnterScope(t *NodeTraversal) {
	scopeRoot := t.GetScopeRoot()
	if !scopeRoot.IsFunction() {
		return
	}
	scope := t.GetScope()
	var result string
	if e := infraTry(func() {
		escaped := jutil.NewObjectLinkedHashSet[*Var]()
		all := NodeUtil_GetAllVarsDeclaredInFunction(f.d.compiler, f.scopeCreator, scope)
		DataFlowAnalysis_ComputeEscaped(scope, escaped, f.d.compiler, f.scopeCreator, all.GetAllVariables())
		names := jutil.NewTreeSet(jutil.StringCompareTo)
		for v := range escaped.All() {
			names.Add(v.GetName() + "@" + f.d.id(v.GetNameNode()))
		}
		result = names.String()
	}); e != nil {
		result = "X " + infraExc(e)
	}
	f.lines = append(f.lines, "E "+f.d.id(scopeRoot)+": "+result)
}

func (f *infraEscapedFinder) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {}

func (d *infraDumper) escaped(script *rhino.Node) {
	finder := &infraEscapedFinder{d: d, scopeCreator: NewSyntacticScopeCreator(d.compiler)}
	if e := infraTry(func() {
		NodeTraversal_Builder_().
			SetCompiler(d.compiler).
			SetScopeCreator(finder.scopeCreator).
			SetCallback(finder).
			Traverse(script)
	}); e != nil {
		finder.lines = append(finder.lines, "E X "+infraExc(e))
	}
	for _, l := range finder.lines {
		d.line(l)
	}
}

// ---------------------------------------------------------------------------
// Comparison

// infraReadDump splits a (concatenated) gzip dump into per-file blocks.
func infraReadDump(t *testing.T, path string) (names []string, blocks map[string]string) {
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
	data, err := io.ReadAll(zr)
	if err != nil {
		t.Fatal(err)
	}
	blocks = map[string]string{}
	var cur string
	var b strings.Builder
	flush := func() {
		if cur != "" {
			blocks[cur] = b.String()
		}
		b.Reset()
	}
	sc := bufio.NewScanner(strings.NewReader(string(data)))
	sc.Buffer(make([]byte, 1<<20), 1<<30)
	for sc.Scan() {
		l := sc.Text()
		if name, ok := strings.CutPrefix(l, "## file "); ok {
			flush()
			cur = name
			names = append(names, name)
		}
		b.WriteString(l)
		b.WriteByte('\n')
	}
	flush()
	return names, blocks
}

// infraCompare dumps every file of the oracle and reports differences.
func infraCompare(t *testing.T, oracle string, rootFor func(name string) string, diffDir string) {
	t.Helper()
	names, blocks := infraReadDump(t, oracle)
	diffs := 0
	for _, name := range names {
		root := rootFor(name)
		if root == "" {
			continue
		}
		src, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Errorf("%s: %v", name, err)
			continue
		}
		want := blocks[name]
		var got string
		if e := infraTry(func() { got = (&infraDumper{}).dumpOne(name, string(src)) }); e != nil {
			got = "## file " + name + "\n# go-exception " + e.Error() + "\n"
		}
		if got == want {
			continue
		}
		diffs++
		if diffs <= 10 {
			gl, wl := strings.Split(got, "\n"), strings.Split(want, "\n")
			for i := 0; i < len(gl) || i < len(wl); i++ {
				var g, w string
				if i < len(gl) {
					g = gl[i]
				}
				if i < len(wl) {
					w = wl[i]
				}
				if g != w {
					t.Errorf("%s: first difference at line %d:\n  go:   %s\n  java: %s", name, i+1, g, w)
					break
				}
			}
		}
		if diffDir != "" {
			out := filepath.Join(diffDir, name)
			_ = os.MkdirAll(filepath.Dir(out), 0o755)
			_ = os.WriteFile(out+".go.txt", []byte(got), 0o644)
			_ = os.WriteFile(out+".java.txt", []byte(want), 0o644)
		}
	}
	if diffs > 0 {
		t.Errorf("%d of %d files differ", diffs, len(names))
	}
}

func TestInfraDumpGolden(t *testing.T) {
	jsDir := filepath.Join("testdata", "infra", "js")
	golden := map[string]bool{}
	if data, err := os.ReadFile(filepath.Join("testdata", "infra", "golden-files.txt")); err == nil {
		for _, l := range strings.Split(string(data), "\n") {
			if l != "" {
				golden[l] = true
			}
		}
	}
	_, corpusErr := os.Stat(infraCorpusRoot)
	infraCompare(t, filepath.Join("testdata", "infra", "golden.txt.gz"), func(name string) string {
		if golden[name] {
			if corpusErr != nil {
				return ""
			}
			return infraCorpusRoot
		}
		return jsDir
	}, "")
}

func TestInfraDumpOracle(t *testing.T) {
	oracle := os.Getenv("INFRA_ORACLE")
	if oracle == "" {
		t.Skip("set INFRA_ORACLE=<dump.gz> (tools/oracle/java/InfraDump.java)")
	}
	root := os.Getenv("INFRA_ROOT")
	if root == "" {
		root = infraCorpusRoot
	}
	infraCompare(t, oracle, func(string) string { return root }, os.Getenv("INFRA_DIFF_DIR"))
}
