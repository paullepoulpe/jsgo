package jscomp

// Go side of tools/oracle/java/NodeUtilSweep.java: an API sweep of NodeUtil
// over real ASTs. Every NodeUtil predicate/getter whose arguments are just a
// node (plus fixed constants, the node's parent or the test compiler) is
// called on every node of every AST, and the results are printed one line
// per node in the order NodeUtilSweep prints them. The two outputs must be
// byte-identical.
//
//   - TestNodeUtilSweepGolden compares against the checked-in Java output
//     (testdata/node_util/sweep_golden.txt.gz, made with
//     NodeUtilSweep --gzip --snippets=testdata/node_util/sweep_snippets.txt)
//     for the snippets of testdata/node_util/sweep_snippets.txt: every code
//     string of NodeUtilTest plus hand-written edge cases.
//   - TestNodeUtilSweepOracle (env-gated) compares any Java output:
//     NODEUTIL_SWEEP_ROOT=<corpus root> NODEUTIL_SWEEP_FILES=<list file>
//     NODEUTIL_SWEEP_ORACLE=<NodeUtilSweep output, .gz or plain>
//     [NODEUTIL_SWEEP_MAX_DIFFS=<n>]. The Java output for the corpus is made by
//     java -Xss1g -cp closure-compiler.jar:build com.google.javascript.jscomp.NodeUtilSweep
//     --gzip <out.gz> .cache/corpus @<list file>.

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
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var (
	nodeUtilSweep_GOOG_MODULE_GET   = rhino.QualifiedName_Of("goog.module.get")
	nodeUtilSweep_GOOG_PROVIDE_NODE = rhino.IR_Getprop(rhino.IR_Name("goog"), "provide")
)

// nodeUtilSweepConfig is NodeUtilSweep.compilerParserConfig(): the config
// Compiler.getParserConfig(ConfigContext.DEFAULT) builds for ECMASCRIPT_NEXT
// with default options.
func nodeUtilSweepConfig() *parsing.Config {
	return parsing.ParserRunner_CreateConfigWithModes(parsing.Config_LanguageMode_ES_NEXT,
		parsing.Config_JsDocParsing_TYPES_ONLY, parsing.Config_RunMode_STOP_AFTER_ERROR, nil, true,
		parsing.Config_StrictMode_STRICT)
}

type nodeUtilSweeper struct {
	index    map[*rhino.Node]int
	compiler AbstractCompiler
	script   *rhino.Node
	mutate   bool

	// hasConvention is whether compiler.GetCodingConvention() works; the
	// positions of the probes that need it are masked when it does not.
	hasConvention bool
	maskBools     map[int]bool
	maskFields    map[int]bool
}

// newNodeUtilSweeper is the sweeper of NodeUtilSweep: `new Compiler()` +
// `initOptions(new CompilerOptions())` is the test compiler.
func newNodeUtilSweeper() *nodeUtilSweeper {
	s := &nodeUtilSweeper{
		compiler:   newTestCompiler(NewCompilerOptions()),
		maskBools:  map[int]bool{},
		maskFields: map[int]bool{},
	}
	s.hasConvention = true
	func() {
		defer func() {
			if r := recover(); r != nil {
				s.hasConvention = false
			}
		}()
		s.compiler.GetCodingConvention()
	}()
	return s
}

// nodeUtilSweepMask masks, in a Java node line, the probes that the Go sweep
// could not run (written '~'), so the rest of the line is still compared.
func (s *nodeUtilSweeper) mask(chunk string) string {
	if s.hasConvention {
		return chunk
	}
	lines := strings.Split(chunk, "\n")
	for i, l := range lines {
		if l == "" || l[0] < '0' || l[0] > '9' {
			continue
		}
		fields := strings.Split(l, "\t")
		head := []byte(fields[0])
		start := strings.IndexByte(fields[0], ' ')
		start += strings.IndexByte(fields[0][start+1:], ' ') + 2
		for k := range s.maskBools {
			if start+k < len(head) {
				head[start+k] = '~'
			}
		}
		fields[0] = string(head)
		for k := range s.maskFields {
			if k+1 < len(fields) {
				fields[k+1] = "~"
			}
		}
		lines[i] = strings.Join(fields, "\t")
	}
	return strings.Join(lines, "\n")
}

type nodeUtilSilentReporter struct{}

func (nodeUtilSilentReporter) Warning(message string, sourceName string, line int, lineOffset int) {}
func (nodeUtilSilentReporter) Error(message string, sourceName string, line int, lineOffset int)   {}

// sweepFile is NodeUtilSweep.sweepFile.
func (s *nodeUtilSweeper) sweepFile(sb *strings.Builder, config *parsing.Config, name, src string) {
	sb.WriteString("# file " + nodeUtilSweepEscape(name) + "\n")
	var result *parsing.ParserRunner_ParseResult
	var exception jutil.Throwable
	jutil.Try(func() {
		result = parsing.ParserRunner_Parse(SourceFile_FromCode(name, src), src, config, nodeUtilSilentReporter{})
	}, func(e jutil.IRuntimeException) { exception = e })
	if exception != nil {
		sb.WriteString("# exception " + nodeUtilSweepClassName(exception) + "\n")
		return
	}
	if result.Ast == nil {
		sb.WriteString("# tree null\n")
		return
	}
	s.index = map[*rhino.Node]int{}
	var nodes []*rhino.Node
	nodes = s.collect(result.Ast, nodes)
	s.script = result.Ast
	s.mutate = len(nodes) <= nodeUtilSweepMaxNodesForMutations
	for i, n := range nodes {
		s.sweepNode(sb, i, n)
	}
	s.sweepScript(sb, result.Ast)
}

func (s *nodeUtilSweeper) collect(n *rhino.Node, nodes []*rhino.Node) []*rhino.Node {
	s.index[n] = len(nodes)
	nodes = append(nodes, n)
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		nodes = s.collect(c, nodes)
	}
	return nodes
}

// nodeUtilSweepClassName is Throwable.getClass().getName().
func nodeUtilSweepClassName(t jutil.Throwable) string {
	full := t.Error()
	if msg := t.GetMessage(); strings.HasSuffix(full, ": "+msg) {
		return strings.TrimSuffix(full, ": "+msg)
	}
	return full
}

// nodeUtilSweepCode is NodeUtilSweep.code(Throwable).
func nodeUtilSweepCode(r any) byte {
	t := jutil.AsThrowable(r)
	switch t.(type) {
	case jutil.IIllegalStateException:
		return 'S'
	case jutil.IIllegalArgumentException:
		return 'A'
	case jutil.INullPointerException:
		return 'N'
	case jutil.IUnsupportedOperationException:
		return 'U'
	case jutil.IClassCastException:
		return 'C'
	case jutil.IIndexOutOfBoundsException:
		return 'I'
	case jutil.IError:
		return 'E'
	}
	return 'X'
}

// bool is NodeUtilSweep.bool.
func nodeUtilSweepBool(sb *strings.Builder, applicable bool, p func() bool) {
	if !applicable {
		sb.WriteByte('.')
		return
	}
	var c byte
	func() {
		defer func() {
			if r := recover(); r != nil {
				c = nodeUtilSweepCode(r)
			}
		}()
		if p() {
			c = '1'
		} else {
			c = '0'
		}
	}()
	sb.WriteByte(c)
}

// field is NodeUtilSweep.field.
func nodeUtilSweepField(sb *strings.Builder, applicable bool, p func() string) {
	sb.WriteByte('\t')
	if !applicable {
		sb.WriteByte('.')
		return
	}
	var s string
	func() {
		defer func() {
			if r := recover(); r != nil {
				s = "!" + string(nodeUtilSweepCode(r))
			}
		}()
		s = p()
	}()
	sb.WriteString(nodeUtilSweepEscape(s))
}

func (s *nodeUtilSweeper) idx(n *rhino.Node) string {
	if n == nil {
		return "-"
	}
	i, ok := s.index[n]
	if !ok {
		return "?"
	}
	return strconv.Itoa(i)
}

func (s *nodeUtilSweeper) list(nodes []*rhino.Node) string {
	var b strings.Builder
	b.WriteByte('[')
	for i, n := range nodes {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(s.idx(n))
	}
	b.WriteByte(']')
	return b.String()
}

func nodeUtilSweepStr(s *string) string {
	if s == nil {
		return "-"
	}
	return nodeUtilSweepQuote(*s)
}

func nodeUtilSweepB(v bool) string {
	if v {
		return "1"
	}
	return "0"
}

func nodeUtilSweepDbl(d *float64) string {
	if d == nil {
		return "-"
	}
	return jutil.DoubleToString(*d)
}

func nodeUtilSweepTok(t rhino.Token) string {
	return t.Name()
}

func nodeUtilSweepQuote(s string) string {
	return "\"" + strings.ReplaceAll(s, "\"", "\\\"") + "\""
}

// nodeUtilSweepEscape is NodeUtilSweep.escape (ParsingDump.escape).
func nodeUtilSweepEscape(s string) string {
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

// nodeUtilSweepCompact is NodeUtilSweep.compact(Node).
func nodeUtilSweepCompact(n *rhino.Node) string {
	var b strings.Builder
	b.WriteString(n.GetToken().Name())
	if n.IsName() || n.IsGetProp() || n.IsStringLit() || n.IsStringKey() {
		b.WriteString("[" + nodeUtilSweepQuote(n.GetString()) + "]")
	}
	b.WriteString("{" + strconv.Itoa(n.GetLength()) + "@" + strconv.Itoa(n.GetLineno()) + ":" +
		strconv.Itoa(n.GetCharno()) + "}")
	if n.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME) {
		b.WriteByte('c')
	}
	if o := n.GetOriginalName(); o != nil {
		b.WriteString("o=" + nodeUtilSweepQuote(*o))
	}
	if n.HasChildren() {
		b.WriteByte('(')
		for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
			if c != n.GetFirstChild() {
				b.WriteByte(',')
			}
			b.WriteString(nodeUtilSweepCompact(c))
		}
		b.WriteByte(')')
	}
	return b.String()
}

// nodeUtilSweepMaxNodesForMutations is NodeUtilSweep.MAX_NODES_FOR_MUTATIONS.
const nodeUtilSweepMaxNodesForMutations = 1000

// mutated is NodeUtilSweep.mutated: the mutation runs on the node with the
// same pre-order index in a fresh clone of the script.
func (s *nodeUtilSweeper) mutated(i int, m func(clone, c *rhino.Node) string) string {
	clone := s.script.CloneTree()
	var cloneNodes []*rhino.Node
	cloneNodes = nodeUtilSweepCollectInto(clone, cloneNodes)
	r := m(clone, cloneNodes[i])
	return r + ":" + strconv.Itoa(int(jutil.StringHashCode(nodeUtilSweepCompact(clone))))
}

func nodeUtilSweepCollectInto(n *rhino.Node, nodes []*rhino.Node) []*rhino.Node {
	nodes = append(nodes, n)
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		nodes = nodeUtilSweepCollectInto(c, nodes)
	}
	return nodes
}

func nodeUtilSweepHasString(n *rhino.Node) bool {
	return n.IsName() || n.IsStringKey() || n.IsGetProp() || n.IsStringLit()
}

func nodeUtilSweepIsObjectOrClassLitKeyToken(n *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_STRING_KEY, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_MEMBER_FUNCTION_DEF,
		rhino.Token_MEMBER_FIELD_DEF, rhino.Token_COMPUTED_PROP, rhino.Token_COMPUTED_FIELD_DEF:
		return true
	}
	return false
}

func nodeUtilSweepFirstKey(n *rhino.Node) string {
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		if c.IsStringKey() || c.IsMemberFunctionDef() || c.IsGetterDef() || c.IsSetterDef() {
			return c.GetString()
		}
	}
	return "x"
}

func nodeUtilSweepIsAssigningParent(n *rhino.Node) bool {
	return NodeUtil_IsNameDeclaration(n) ||
		n.IsParamList() ||
		NodeUtil_IsAssignmentOp(n) ||
		n.IsCatch() ||
		n.IsDestructuringLhs() ||
		n.IsDefaultValue() ||
		n.IsImport() ||
		NodeUtil_IsEnhancedFor(n)
}

// sweepNode is NodeUtilSweep.sweepNode: the probes in the same order.
func (s *nodeUtilSweeper) sweepNode(sb *strings.Builder, i int, n *rhino.Node) {
	parent := n.GetParent()
	t := n.GetToken()
	var line strings.Builder
	line.WriteString(strconv.Itoa(i) + " " + t.Name() + " ")

	isFn := n.IsFunction()
	isFnOrClass := n.IsFunction() || n.IsClass()
	hasParent := parent != nil
	isName := n.IsName()
	isGet := NodeUtil_IsNormalOrOptChainGet(n)
	isQName := n.IsQualifiedName()

	nBools, nFields := 0, 0
	bool_ := func(applicable bool, p func() bool) { nodeUtilSweepBool(&line, applicable, p); nBools++ }
	field := func(applicable bool, p func() string) { nodeUtilSweepField(&line, applicable, p); nFields++ }
	// Probes that need the compiler's coding convention: without one (the
	// test compiler needs ClosureCodingConvention for the default), they are
	// written as '~' and masked in the Java line (see nodeUtilSweepMask).
	boolC := func(applicable bool, p func() bool) {
		if s.hasConvention {
			bool_(applicable, p)
			return
		}
		s.maskBools[nBools] = true
		line.WriteByte('~')
		nBools++
	}
	fieldC := func(applicable bool, p func() string) {
		if s.hasConvention {
			field(applicable, p)
			return
		}
		s.maskFields[nFields] = true
		line.WriteString("\t~")
		nFields++
	}

	bool_(true, func() bool { return NodeUtil_IsImmutableValue(n) })
	bool_(true, func() bool { return NodeUtil_IsSymmetricOperation(n) })
	bool_(true, func() bool { return NodeUtil_IsRelationalOperation(n) })
	bool_(true, func() bool { return NodeUtil_IsLiteralValue(n, false) })
	bool_(true, func() bool { return NodeUtil_IsLiteralValue(n, true) })
	bool_(true, func() bool { return NodeUtil_IsSomeCompileTimeConstStringValue(n) })
	bool_(true, func() bool { return NodeUtil_IsEmptyBlock(n) })
	bool_(true, func() bool { return NodeUtil_IsBinaryOperator(n) })
	bool_(true, func() bool { return NodeUtil_IsUnaryOperator(n) })
	bool_(true, func() bool { return NodeUtil_IsUpdateOperator(n) })
	bool_(true, func() bool { return NodeUtil_IsSimpleOperator(n) })
	bool_(hasParent, func() bool { return NodeUtil_IsNamespaceDecl(n) })
	bool_(true, func() bool { return NodeUtil_IsPureIterable(n) })
	bool_(true, func() bool { return NodeUtil_CanBeSideEffected(n) })
	bool_(true, func() bool { return NodeUtil_IsUndefined(n) })
	bool_(true, func() bool { return NodeUtil_IsNullOrUndefined(n) })
	bool_(true, func() bool { return NodeUtil_IsNumericResult(n) })
	bool_(true, func() bool { return NodeUtil_IsBigIntResult(n) })
	bool_(true, func() bool { return NodeUtil_IsBooleanResult(n) })
	bool_(true, func() bool { return NodeUtil_IsStringResult(n) })
	bool_(true, func() bool { return NodeUtil_IsObjectResult(n) })
	bool_(true, func() bool { return NodeUtil_MayBeString(n) })
	bool_(true, func() bool { return NodeUtil_MayBeStringWithUseType(n, true) })
	bool_(true, func() bool { return NodeUtil_MayBeObject(n) })
	bool_(true, func() bool { return NodeUtil_IsAssociative(t) })
	bool_(true, func() bool { return NodeUtil_IsCommutative(t) })
	bool_(true, func() bool { return NodeUtil_IsAssignmentOp(n) })
	bool_(true, func() bool { return NodeUtil_IsLogicalAssignmentOp(n) })
	bool_(true, func() bool { return NodeUtil_IsCompoundAssignmentOp(n) })
	bool_(true, func() bool { return NodeUtil_IsHoistScopeRoot(n) })
	bool_(true, func() bool { return NodeUtil_IsInFunction(n) })
	bool_(true, func() bool { return NodeUtil_ReferencesEnclosingReceiver(n) })
	bool_(true, func() bool { return NodeUtil_ReferencesSuper(n) })
	bool_(true, func() bool { return NodeUtil_IsNormalOrOptChainGet(n) })
	bool_(true, func() bool { return NodeUtil_IsNormalOrOptChainGetProp(n) })
	bool_(true, func() bool { return NodeUtil_IsNormalOrOptChainCall(n) })
	bool_(true, func() bool { return NodeUtil_IsNormalGet(n) })
	bool_(true, func() bool { return NodeUtil_IsOptChainGet(n) })
	bool_(true, func() bool { return NodeUtil_IsOptChainNode(n) })
	bool_(true, func() bool { return NodeUtil_IsEndOfFullOptChain(n) })
	bool_(true, func() bool { return NodeUtil_IsEndOfOptChainSegment(n) })
	bool_(true, func() bool { return NodeUtil_IsBlockScopedDeclaration(n) })
	bool_(true, func() bool { return NodeUtil_IsNameDeclaration(n) })
	bool_(true, func() bool { return NodeUtil_IsDestructuringDeclaration(n) })
	bool_(true, func() bool { return NodeUtil_IsExprAssign(n) })
	bool_(true, func() bool { return NodeUtil_IsExprCall(n) })
	bool_(true, func() bool { return NodeUtil_IsNonArrowFunction(n) })
	bool_(true, func() bool { return NodeUtil_IsEnhancedFor(n) })
	bool_(true, func() bool { return NodeUtil_IsAnyFor(n) })
	bool_(true, func() bool { return NodeUtil_IsLoopStructure(n) })
	bool_(true, func() bool { return NodeUtil_IsWithinLoop(n) })
	bool_(true, func() bool { return NodeUtil_IsControlStructure(n) })
	bool_(true, func() bool { return NodeUtil_IsStatementBlock(n) })
	bool_(true, func() bool { return NodeUtil_CreatesBlockScope(n) })
	bool_(true, func() bool { return NodeUtil_CreatesScope(n) })
	bool_(true, func() bool { return NodeUtil_IsValidCfgRoot(n) })
	bool_(true, func() bool { return NodeUtil_IsStatement(n) })
	bool_(true, func() bool { return NodeUtil_IsStatementParent(n) })
	bool_(true, func() bool { return NodeUtil_IsDeclarationParent(n) })
	bool_(true, func() bool { return NodeUtil_IsSwitchCase(n) })
	bool_(true, func() bool { return NodeUtil_IsReferenceName(n) })
	bool_(true, func() bool { return NodeUtil_IsInSyntheticScript(n) })
	bool_(true, func() bool { return NodeUtil_IsCallOrNew(n) })
	bool_(true, func() bool { return NodeUtil_IsDeclaration(n) })
	bool_(true, func() bool { return NodeUtil_IsNamedExportsLiteral(n) })
	bool_(true, func() bool { return NodeUtil_IsFunctionDeclaration(n) })
	bool_(true, func() bool { return NodeUtil_IsMethodDeclaration(n) })
	bool_(true, func() bool { return NodeUtil_IsClassDeclaration(n) })
	bool_(true, func() bool { return NodeUtil_IsHoistedFunctionDeclaration(n) })
	bool_(true, func() bool { return NodeUtil_IsBlockScopedFunctionDeclaration(n) })
	bool_(true, func() bool { return NodeUtil_IsFunctionBlock(n) })
	bool_(true, func() bool { return NodeUtil_IsClassStaticBlock(n) })
	bool_(true, func() bool { return NodeUtil_IsFunctionExpression(n) })
	bool_(true, func() bool { return NodeUtil_IsNamedFunctionExpression(n) })
	bool_(true, func() bool { return NodeUtil_IsClassExpression(n) })
	bool_(true, func() bool { return NodeUtil_IsNamedClassExpression(n) })
	bool_(true, func() bool { return NodeUtil_IsNamedFunction(n) })
	bool_(true, func() bool { return NodeUtil_IsNamedClass(n) })
	bool_(true, func() bool { return NodeUtil_IsBleedingFunctionName(n) })
	bool_(true, func() bool { return NodeUtil_IsEmptyFunctionExpression(n) })
	bool_(true, func() bool { return NodeUtil_IsFunctionObjectCall(n) })
	bool_(true, func() bool { return NodeUtil_IsFunctionObjectApply(n) })
	bool_(true, func() bool { return NodeUtil_IsLValue(n) })
	bool_(true, func() bool { return NodeUtil_IsDeclarationLValue(n) })
	bool_(true, func() bool { return NodeUtil_IsLhsOfAssign(n) })
	bool_(hasParent, func() bool { return NodeUtil_IsImportedName(n) })
	bool_(true, func() bool { return NodeUtil_IsLhsByDestructuring(n) })
	bool_(true, func() bool { return NodeUtil_MayBeObjectLitKey(n) })
	bool_(hasParent, func() bool { return NodeUtil_IsObjectLitKey(n) })
	bool_(true, func() bool { return NodeUtil_IsGetOrSetKey(n) })
	bool_(true, func() bool { return NodeUtil_CanContainHoistedVarsDecls(n) })
	bool_(true, func() bool { return NodeUtil_IsObjectDefinePropertiesDefinition(n) })
	bool_(true, func() bool { return NodeUtil_IsObjectDefinePropertyDefinition(n) })
	bool_(true, func() bool { return NodeUtil_IsPrototypePropertyDeclaration(n) })
	bool_(true, func() bool { return NodeUtil_IsPrototypeProperty(n) })
	bool_(true, func() bool { return NodeUtil_IsPrototypeMethod(n) })
	bool_(true, func() bool { return NodeUtil_IsPrototypeAssignment(n) })
	bool_(hasParent, func() bool { return NodeUtil_IsPropertyAbsenceTest(n) })
	bool_(true, func() bool { return NodeUtil_IsShallowStatementTree(n) })
	bool_(true, func() bool { return NewNodeUtil_MatchShallowStatement().Apply(n) })
	bool_(true, func() bool { return NewNodeUtil_MatchDeclaration().Apply(n) })
	bool_(true, func() bool { return NodeUtil_IsConstantVar(n, nil) })
	bool_(true, func() bool { return NodeUtil_IsConstantName(n) })
	bool_(true, func() bool { return NodeUtil_FunctionHasInlineJsdocs(n) })
	bool_(true, func() bool { return NodeUtil_EvaluatesToLocalValue(n) })
	bool_(true, func() bool { return NodeUtil_MayBeUndefined(n) })
	bool_(true, func() bool { return NodeUtil_IsDefinedValue(n) })
	bool_(true, func() bool { return NodeUtil_IsInvocationTarget(n) })
	bool_(true, func() bool { return NodeUtil_IsInvocation(n) })
	bool_(true, func() bool { return NodeUtil_IsCallOrNewArgument(n) })
	bool_(hasParent, func() bool { return NodeUtil_IsExpressionResultUsed(n) })
	bool_(hasParent, func() bool { return NodeUtil_IsExecutedExactlyOnce(n) })
	bool_(true, func() bool { return NodeUtil_IsNaN(n) })
	bool_(true, func() bool { return NodeUtil_IsGoogProvideCall(n) })
	bool_(true, func() bool { return NodeUtil_IsGoogModuleCall(n) })
	bool_(true, func() bool { return NodeUtil_IsGoogModuleGetCall(n) })
	bool_(true, func() bool { return NodeUtil_IsGoogRequireCall(n) })
	bool_(true, func() bool { return NodeUtil_IsGoogRequireTypeCall(n) })
	bool_(true, func() bool { return NodeUtil_IsGoogForwardDeclareCall(n) })
	bool_(true, func() bool { return NodeUtil_IsGoogRequireDynamicCall(n) })
	bool_(true, func() bool { return NodeUtil_IsGoogWeakUsageCall(n) })
	bool_(true, func() bool { return NodeUtil_IsModuleScopeRoot(n) })
	bool_(true, func() bool { return NodeUtil_IsBundledGoogModuleCall(n) })
	bool_(true, func() bool { return NodeUtil_IsBundledGoogModuleScopeRoot(n) })
	bool_(true, func() bool { return NodeUtil_IsGoogModuleDeclareLegacyNamespaceCall(n) })
	bool_(true, func() bool { return NodeUtil_IsGoogSetTestOnlyCall(n) })
	bool_(true, func() bool { return NodeUtil_IsTopLevel(n) })
	bool_(true, func() bool { return NodeUtil_IsGoogModuleFile(n) })
	bool_(true, func() bool { return NodeUtil_IsLegacyGoogModuleFile(n) })
	bool_(true, func() bool { return NodeUtil_IsConstructor(n) })
	bool_(true, func() bool { return NodeUtil_IsEs6ConstructorMemberFunctionDef(n) })
	bool_(true, func() bool { return NodeUtil_IsEs6Constructor(n) })
	bool_(true, func() bool { return NodeUtil_IsGetterOrSetter(n) })
	bool_(true, func() bool { return NodeUtil_IsCallTo(n, "goog.require") })
	bool_(true, func() bool { return NodeUtil_IsCallToQualifiedName(n, nodeUtilSweep_GOOG_MODULE_GET) })
	bool_(true, func() bool { return NodeUtil_IsCallToNode(n, nodeUtilSweep_GOOG_PROVIDE_NODE) })
	bool_(true, func() bool { return NodeUtil_IsObjLitProperty(n) })
	bool_(true, func() bool { return NodeUtil_IsBlocklessArrowFunctionResult(n) })
	bool_(hasParent, func() bool { return NodeUtil_IsTryCatchNodeContainer(n) })
	bool_(hasParent, func() bool { return NodeUtil_IsTryFinallyNode(parent, n) })
	bool_(hasParent, func() bool { return NodeUtil_IsNameDeclOrSimpleAssignLhs(n, parent) })
	bool_(isName, func() bool { return NodeUtil_IsNonlocalModuleExportName(n) })
	boolC(isGet, func() bool { return NodeUtil_IsPropertyTest(s.compiler, n) })
	bool_(true, func() bool { return NodeUtil_IsKnownGlobalObjectReference(n) })
	bool_(true, func() bool { return NodeUtil_IsRhsClass(n) })
	bool_(true, func() bool { return NodeUtil_IsRhsFunction(n) })
	bool_(n.IsCall() || n.IsOptChainCall(), func() bool { return NodeUtil_IsToStringMethodCall(n) })
	bool_(isFn, func() bool { return NodeUtil_DoesFunctionReferenceOwnArgumentsObject(n) })
	bool_(isFn, func() bool { return NodeUtil_ReferencesOwnReceiver(n) })
	bool_(n.IsScript(), func() bool { return NodeUtil_IsFromTypeSummary(n) })
	bool_(true, func() bool { return NodeUtil_CanMergeBlock(n) })
	bool_(n.IsExport(), func() bool { return NodeUtil_IsExportFrom(n) })
	bool_(n.IsTry(), func() bool { return NodeUtil_HasFinally(n) })
	bool_(n.IsBlock(), func() bool { return NodeUtil_HasCatchHandler(n) })
	bool_(n.IsSpread() || n.IsYield() || n.IsForOf() || n.IsForAwaitOf() || n.IsRest(),
		func() bool { return NodeUtil_IteratesImpureIterable(n) })
	bool_(n.IsNew(), func() bool { return NodeUtil_NewHasLocalResult(n) })
	bool_(NodeUtil_IsCallOrNew(n), func() bool { return NodeUtil_AllArgsUnescapedLocal(n) })
	bool_(isFn, func() bool { return NodeUtil_IsNameReferenced(n, "arguments") })
	bool_(hasParent &&
		(NodeUtil_IsObjectLitKey(n) ||
			NodeUtil_IsLhsOfAssign(n) ||
			(isName && NodeUtil_IsDeclarationLValue(n)) ||
			n.IsMemberFieldDef() ||
			n.IsComputedFieldDef()),
		func() bool { return NodeUtil_IsConstantDeclaration(NodeUtil_GetBestJSDocInfo(n), n) })
	bool_(hasParent && NodeUtil_IsControlStructure(parent),
		func() bool { return NodeUtil_IsControlStructureCodeBlock(parent, n) })

	field(true, func() string { return NodeUtil_GetBooleanValue(n).Name() })
	field(true, func() string { return nodeUtilSweepStr(NodeUtil_GetStringValue(n)) })
	field(true, func() string { return nodeUtilSweepStr(NodeUtil_GetArrayElementStringValue(n)) })
	field(n.IsArrayLit(), func() string { return nodeUtilSweepStr(NodeUtil_ArrayToString(n)) })
	field(true, func() string { return nodeUtilSweepDbl(NodeUtil_GetNumberValue(n)) })
	field(true, func() string { return nodeUtilSweepDbl(NodeUtil_GetNumberValueNoConversions(n)) })
	field(true, func() string {
		x := NodeUtil_GetBigIntValue(n)
		if x == nil {
			return "-"
		}
		return x.String()
	})
	field(isFnOrClass, func() string { return nodeUtilSweepStr(NodeUtil_GetName(n)) })
	field(isFnOrClass, func() string { return s.idx(NodeUtil_GetNameNode(n)) })
	field(true, func() string { return nodeUtilSweepStr(NodeUtil_GetNearestFunctionName(n)) })
	field(n.IsClass(), func() string { return s.idx(NodeUtil_GetClassMembers(n)) })
	field(n.IsClass(), func() string { return s.idx(NodeUtil_GetEs6ClassConstructorMemberFunctionDef(n)) })
	field(NodeUtil_IsRelationalOperation(n), func() string { return nodeUtilSweepTok(NodeUtil_GetInverseOperator(t)) })
	field(n.IsBlock(), func() string { return s.idx(NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(n)) })
	field(true, func() string { return strconv.Itoa(NodeUtil_Precedence(t)) })
	field(true, func() string { return NodeUtil_GetKnownValueType(n).Name() })
	field(NodeUtil_IsCompoundAssignmentOp(n), func() string { return nodeUtilSweepTok(NodeUtil_GetOpFromAssignmentOp(n)) })
	field(true, func() string { return s.idx(NodeUtil_GetEnclosingNonArrowFunction(n)) })
	field(true, func() string { return s.idx(NodeUtil_GetEnclosingClass(n)) })
	field(true, func() string { return s.idx(NodeUtil_GetEnclosingModuleIfPresent(n)) })
	field(true, func() string { return s.idx(NodeUtil_GetEnclosingFunction(n)) })
	field(true, func() string { return s.idx(NodeUtil_GetEnclosingScript(n)) })
	field(true, func() string { return s.idx(NodeUtil_GetEnclosingBlock(n)) })
	field(true, func() string { return s.idx(NodeUtil_GetEnclosingBlockScopeRoot(n)) })
	field(true, func() string { return s.idx(NodeUtil_GetEnclosingScopeRoot(n)) })
	field(true, func() string { return s.idx(NodeUtil_GetEnclosingHoistScopeRoot(n)) })
	field(true, func() string { return s.idx(NodeUtil_GetEnclosingStatement(n)) })
	field(true, func() string { return s.idx(NodeUtil_GetEnclosingType(n, rhino.Token_CALL)) })
	field(n.IsObjectLit() || n.IsClassMembers(),
		func() string { return s.idx(NodeUtil_GetFirstPropMatchingKey(n, nodeUtilSweepFirstKey(n))) })
	field(n.IsObjectLit() || n.IsClassMembers(),
		func() string { return s.idx(NodeUtil_GetFirstGetterMatchingKey(n, nodeUtilSweepFirstKey(n))) })
	field(NodeUtil_IsOptChainNode(n), func() string { return s.idx(NodeUtil_GetStartOfOptChainSegment(n)) })
	field(NodeUtil_IsOptChainNode(n), func() string { return s.idx(NodeUtil_GetEndOfOptChainSegment(n)) })
	field(NodeUtil_IsOptChainNode(n), func() string { return nodeUtilSweepTok(NodeUtil_GetNonOptChainToken(t)) })
	field(isName || n.IsGetProp(), func() string { return s.idx(NodeUtil_GetAssignedValue(n)) })
	field(true, func() string { return s.idx(NodeUtil_GetLoopCodeBlock(n)) })
	field(NodeUtil_IsControlStructure(n), func() string { return s.idx(NodeUtil_GetConditionExpression(n)) })
	field(isFn, func() string { return s.idx(NodeUtil_GetFunctionBody(n)) })
	field(n.IsCall() || n.IsNew() || n.IsTaggedTemplateLit(),
		func() string { return s.idx(NodeUtil_GetCallTargetResolvingIndirectCalls(n)) })
	isTarget := hasParent && n.IsValidAssignmentTarget()
	field(isTarget, func() string { return s.idx(NodeUtil_GetEnclosingTarget(n)) })
	field(isTarget, func() string { return s.idx(NodeUtil_GetRootTarget(n)) })
	field(isTarget, func() string { return s.idx(NodeUtil_GetDeclaringParent(n)) })
	field(nodeUtilSweepIsObjectOrClassLitKeyToken(n),
		func() string { return nodeUtilSweepQuote(NodeUtil_GetObjectOrClassLitKeyName(n)) })
	field(nodeUtilSweepIsObjectOrClassLitKeyToken(n), func() string { return s.idx(NodeUtil_GetObjectOrClassLitKeyNode(n)) })
	field(true, func() string { return nodeUtilSweepStr(NodeUtil_OpToStr(t)) })
	field(true, func() string { return NodeUtil_OpToStrNoFail(t) })
	field(isQName, func() string { return s.idx(NodeUtil_GetRootOfQualifiedName(n)) })
	field(isQName, func() string {
		return nodeUtilSweepQuote(NodeUtil_GetRootOfQualifiedNameString(*n.GetQualifiedName()))
	})
	field(isQName, func() string {
		return nodeUtilSweepB(NodeUtil_IsValidQualifiedNameFeatureSet(parser.FeatureSet_ES3, *n.GetQualifiedName()))
	})
	field(isQName, func() string {
		return nodeUtilSweepB(NodeUtil_IsValidQualifiedNameFeatureSet(parser.FeatureSet_ES_NEXT, *n.GetQualifiedName()))
	})
	field(nodeUtilSweepHasString(n), func() string {
		return nodeUtilSweepB(NodeUtil_IsValidPropertyName(parser.FeatureSet_ES3, n.GetString())) +
			nodeUtilSweepB(NodeUtil_IsValidPropertyName(parser.FeatureSet_ES_NEXT, n.GetString())) +
			nodeUtilSweepB(NodeUtil_IsValidSimpleName(n.GetString())) +
			nodeUtilSweepB(NodeUtil_IsLatin(n.GetString()))
	})
	field(n.IsStringLit(), func() string { return nodeUtilSweepDbl(NodeUtil_GetStringNumberValue(n.GetString())) })
	field(n.IsStringLit(), func() string {
		x := NodeUtil_GetStringBigIntValue(n.GetString())
		if x == nil {
			return "-"
		}
		return x.String()
	})
	field(n.IsStringLit(), func() string { return nodeUtilSweepQuote(NodeUtil_TrimJsWhiteSpace(n.GetString())) })
	field(nodeUtilSweepIsAssigningParent(n), func() string {
		var lhs []*rhino.Node
		NodeUtil_VisitLhsNodesInNode(n, func(x *rhino.Node) { lhs = append(lhs, x) })
		return s.list(lhs)
	})
	field(n.IsDestructuringPattern(), func() string {
		var lhs []*rhino.Node
		NodeUtil_VisitLhsNodesInDestructuringPattern(n, func(x *rhino.Node) { lhs = append(lhs, x) })
		return s.list(lhs)
	})
	field(NodeUtil_IsObjectDefinePropertiesDefinition(n),
		func() string { return s.list(NodeUtil_GetObjectDefinedPropertiesKeys(n)) })
	field(true, func() string { return s.idx(NodeUtil_GetPrototypeClassName(n)) })
	field(isQName && strings.Contains(*n.GetQualifiedName(), ".prototype."),
		func() string { return nodeUtilSweepQuote(NodeUtil_GetPrototypePropertyName(n)) })
	field(true, func() string { return strconv.Itoa(NodeUtil_CountAstSizeUpToLimit(n, 10)) })
	field(isFn || n.IsScript(), func() string { return strconv.Itoa(NodeUtil_CountAstSize(n)) })
	field(n.IsTry(), func() string { return s.idx(NodeUtil_GetCatchBlock(n)) })
	field(isFn, func() string { return s.idx(NodeUtil_GetFunctionParameters(n)) })
	field(true, func() string { return nodeUtilSweepStr(NodeUtil_GetSourceName(n)) })
	field(true, func() string {
		f := NodeUtil_GetSourceFile(n)
		if f == nil {
			return "-"
		}
		return nodeUtilSweepQuote(f.GetName())
	})
	field(true, func() string {
		id := NodeUtil_GetInputId(n)
		if id == nil {
			return "-"
		}
		return nodeUtilSweepQuote(id.GetIdName())
	})
	field(isFn, func() string { return s.idx(NodeUtil_GetArgumentForFunction(n, 0)) })
	field(isFn, func() string { return s.idx(NodeUtil_GetArgumentForFunction(n, 1)) })
	field(NodeUtil_IsCallOrNew(n), func() string { return s.idx(NodeUtil_GetArgumentForCallOrNew(n, 0)) })
	field(isName || n.IsStringKey(), func() string {
		e := NodeUtil_GetDeclaredTypeExpression(n)
		if e == nil {
			return "-"
		}
		return nodeUtilSweepQuote(e.GetRoot().String())
	})
	field(true, func() string { return s.idx(NodeUtil_GetBestJSDocInfoNode(n)) })
	field(true, func() string {
		if NodeUtil_GetBestJSDocInfo(n) == nil {
			return "-"
		}
		return "J"
	})
	field(true, func() string { return s.idx(NodeUtil_GetBestJsDocInfoNodeStrict(n)) })
	field(true, func() string { return s.idx(NodeUtil_GetBestLValue(n)) })
	field(hasParent, func() string { return s.idx(NodeUtil_GetRValueOfLValue(n)) })
	field(true, func() string { return s.idx(NodeUtil_GetBestLValueOwner(n)) })
	field(true, func() string { return nodeUtilSweepStr(NodeUtil_GetBestLValueName(n)) })
	field(true, func() string { return s.idx(NodeUtil_GetBestLValueRoot(n)) })
	field(NodeUtil_IsInvocation(n), func() string {
		var args []*rhino.Node
		for a := range NodeUtil_GetInvocationArgsAsIterable(n) {
			args = append(args, a)
		}
		return s.list(args)
	})
	field(NodeUtil_IsInvocation(n), func() string { return strconv.Itoa(NodeUtil_GetInvocationArgsCount(n)) })
	field(n.IsScript(), func() string {
		fs := NodeUtil_GetFeatureSetOfScript(n)
		if fs == nil {
			return "-"
		}
		return nodeUtilSweepQuote(fs.String())
	})
	field(isName || n.IsImportStar(), func() string {
		r := NodeUtil_GetGoogRequireInfoNode(n)
		if r == nil {
			return "-"
		}
		return nodeUtilSweepQuote(r.String())
	})
	field(n.IsScript(), func() string { return strconv.Itoa(NodeUtil_EstimateNumLines(n)) })
	field(n.IsParamList() || n.IsDestructuringPattern(), func() string {
		var names []*rhino.Node
		NodeUtil_GetParamOrPatternNames(n, func(x *rhino.Node) { names = append(names, x) })
		return s.list(names)
	})
	field(n.IsScript() || n.IsBlock(), func() string { return s.list(NodeUtil_GetVarsDeclaredInBranch(n)) })
	fieldC(isQName && !(hasParent && parent.IsGetProp() && parent.GetFirstChild() == n),
		func() string { return nodeUtilSweepCompact(NodeUtil_NewQName(s.compiler, *n.GetQualifiedName())) })
	field(isFn || n.IsScript(), func() string {
		return strconv.Itoa(NodeUtil_GetNodeTypeReferenceCount(n, rhino.Token_NAME, NodeUtil_MATCH_NOT_FUNCTION)) +
			"," + strconv.Itoa(NodeUtil_GetNameReferenceCount(n, "a"))
	})
	field(isFn || n.IsScript(), func() string {
		var pre []*rhino.Node
		for x := range NodeUtil_PreOrderIterableWithPred(n, NodeUtil_MATCH_NOT_FUNCTION) {
			pre = append(pre, x)
		}
		var post []*rhino.Node
		NodeUtil_VisitPostOrderWithPred(n, NodeUtil_VisitorFunc(func(x *rhino.Node) { post = append(post, x) }),
			NodeUtil_MATCH_NOT_FUNCTION)
		var pre2 []*rhino.Node
		NodeUtil_VisitPreOrder(n, NodeUtil_VisitorFunc(func(x *rhino.Node) { pre2 = append(pre2, x) }))
		var h int64
		for _, x := range pre {
			h = h*31 + int64(s.indexOrMinusOne(x))
		}
		for _, x := range post {
			h = h*37 + int64(s.indexOrMinusOne(x))
		}
		return strconv.Itoa(len(pre)) + "," + strconv.Itoa(len(post)) + "," + strconv.Itoa(len(pre2)) + "," +
			strconv.FormatInt(h, 10)
	})
	field(isFn || n.IsScript(),
		func() string {
			return s.idx(NodeUtil_FindPreorder(n, (*rhino.Node).IsName, NodeUtil_MATCH_NOT_FUNCTION))
		})

	// Mutations, each on a fresh clone of the script.
	field(s.mutate && hasParent, func() string {
		return s.mutated(i, func(clone, c *rhino.Node) string { NodeUtil_RemoveChild(c.GetParent(), c); return "" })
	})
	field(s.mutate && n.IsBlock() && hasParent, func() string {
		return s.mutated(i, func(clone, c *rhino.Node) string { return nodeUtilSweepB(NodeUtil_TryMergeBlock(c, false)) })
	})
	field(s.mutate && n.IsBlock() && hasParent, func() string {
		return s.mutated(i, func(clone, c *rhino.Node) string { return nodeUtilSweepB(NodeUtil_TryMergeBlock(c, true)) })
	})
	field(s.mutate && hasParent, func() string {
		return s.mutated(i, func(clone, c *rhino.Node) string { NodeUtil_RedeclareVarsInsideBranch(c); return "" })
	})
	field(s.mutate && NodeUtil_IsEndOfOptChainSegment(n), func() string {
		return s.mutated(i, func(clone, c *rhino.Node) string { NodeUtil_ConvertToNonOptionalChainSegment(c); return "" })
	})
	field(s.mutate && hasParent && NodeUtil_IsNameDeclaration(parent), func() string {
		return s.mutated(i, func(clone, c *rhino.Node) string {
			NodeUtil_ReplaceDeclarationChild(c, rhino.IR_ExprResult(rhino.IR_Number(1)))
			return ""
		})
	})
	field(s.mutate && isFnOrClass, func() string {
		return s.mutated(i, func(clone, c *rhino.Node) string { NodeUtil_RemoveName(c); return "" })
	})
	field(s.mutate && n.IsTry(), func() string {
		return s.mutated(i, func(clone, c *rhino.Node) string { NodeUtil_MaybeAddFinally(c); return "" })
	})

	// Constructions (the AST is not changed).
	field(n.IsNumber() || n.IsName() || n.IsNeg() || n.IsStringLit(), func() string {
		v := NodeUtil_GetNumberValue(n)
		if v == nil {
			return "-"
		}
		return nodeUtilSweepCompact(NodeUtil_NumberNode(*v, n))
	})
	field(n.IsNumber() || n.IsName(), func() string { return nodeUtilSweepCompact(NodeUtil_NewUndefinedNode(n)) })
	isOuterQName := isQName && !(hasParent && parent.IsGetProp() && parent.GetFirstChild() == n)
	fieldC(isOuterQName, func() string {
		return nodeUtilSweepCompact(NodeUtil_NewQNameDeclaration(s.compiler, *n.GetQualifiedName(), nil, nil))
	})
	fieldC(isOuterQName, func() string {
		return nodeUtilSweepCompact(NodeUtil_NewQNameDeclarationWithType(
			s.compiler, *n.GetQualifiedName(), rhino.IR_Number(0), nil, rhino.Token_LET))
	})
	fieldC(isOuterQName, func() string {
		return nodeUtilSweepCompact(NodeUtil_NewQNameWithOriginalName(s.compiler, *n.GetQualifiedName(), n, "orig"))
	})
	fieldC(isOuterQName, func() string {
		return nodeUtilSweepCompact(
			NodeUtil_NewQNameWithOriginalName(s.compiler, *n.GetQualifiedName(), n, *n.GetQualifiedName()))
	})
	field(isName, func() string { return nodeUtilSweepCompact(NodeUtil_NewVarNode(n.GetString(), rhino.IR_Number(1))) })
	fieldC(isName, func() string {
		return nodeUtilSweepCompact(NodeUtil_NewNameWithOriginalName(s.compiler, n.GetString(), n, "orig"))
	})
	fieldC(n.IsGetProp(), func() string {
		return nodeUtilSweepCompact(NodeUtil_NewPropertyAccess(s.compiler, rhino.IR_Name("x"), n.GetString()))
	})

	sb.WriteString(line.String())
	sb.WriteByte('\n')
}

func (s *nodeUtilSweeper) indexOrMinusOne(n *rhino.Node) int {
	if i, ok := s.index[n]; ok {
		return i
	}
	return -1
}

// sweepScript is NodeUtilSweep.sweepScript.
func (s *nodeUtilSweeper) sweepScript(sb *strings.Builder, script *rhino.Node) {
	var line strings.Builder
	line.WriteString("# file-summary")
	field := func(applicable bool, p func() string) { nodeUtilSweepField(&line, applicable, p) }
	field(true, func() string {
		count := 0
		for range NodeUtil_PreOrderIterable(script) {
			count++
		}
		return strconv.Itoa(count)
	})
	field(true, func() string {
		count := 0
		NodeUtil_VisitPostOrder(script, NodeUtil_VisitorFunc(func(*rhino.Node) { count++ }))
		return strconv.Itoa(count)
	})
	field(true, func() string { return strconv.Itoa(NodeUtil_CountAstSize(script)) })
	field(true, func() string {
		return strconv.Itoa(NodeUtil_GetCount(script, (*rhino.Node).IsName, NodeUtil_MATCH_ANYTHING_BUT_NON_ARROW_FUNCTION))
	})
	field(true, func() string {
		return nodeUtilSweepB(NodeUtil_Has(script, (*rhino.Node).IsThis, NodeUtil_MATCH_ANYTHING_BUT_NON_ARROW_FUNCTION))
	})
	sb.WriteString(line.String())
	sb.WriteByte('\n')
}

// nodeUtilSweep runs the sweep over files (relative to root) with a fresh
// test compiler, as NodeUtilSweep.main does.
func nodeUtilSweep(t *testing.T, root string, files []string, emit func(s *nodeUtilSweeper, name, out string)) {
	config := nodeUtilSweepConfig()
	s := newNodeUtilSweeper()
	if !s.hasConvention {
		t.Log("the probes that need the default CodingConvention are masked: needs ClosureCodingConvention")
	}
	for _, rel := range files {
		src, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Fatal(err)
		}
		var sb strings.Builder
		s.sweepFile(&sb, config, rel, string(src))
		emit(s, rel, sb.String())
	}
}

// nodeUtilSweepSplit splits a NodeUtilSweep output into per-file chunks.
func nodeUtilSweepSplit(r io.Reader) (map[string]string, []string, error) {
	chunks := map[string]string{}
	var order []string
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 1<<20), 1<<30)
	var cur string
	var sb strings.Builder
	flush := func() {
		if cur != "" {
			chunks[cur] = sb.String()
		}
		sb.Reset()
	}
	for sc.Scan() {
		l := sc.Text()
		if strings.HasPrefix(l, "# file ") {
			flush()
			cur = strings.TrimPrefix(l, "# file ")
			order = append(order, cur)
		}
		sb.WriteString(l)
		sb.WriteByte('\n')
	}
	flush()
	return chunks, order, sc.Err()
}

// nodeUtilSweepChunks returns an iterator over the per-file chunks of a
// NodeUtilSweep output: each call returns the next chunk's "# file" header
// and the whole chunk (header included).
func nodeUtilSweepChunks(r io.Reader) func() (string, string, bool) {
	br := bufio.NewReaderSize(r, 1<<20)
	var pending string
	done := false
	return func() (string, string, bool) {
		if done && pending == "" {
			return "", "", false
		}
		var sb strings.Builder
		header := pending
		if header != "" {
			sb.WriteString(header)
			sb.WriteByte('\n')
		}
		pending = ""
		for {
			l, err := br.ReadString('\n')
			if l == "" && err != nil {
				done = true
				break
			}
			l = strings.TrimSuffix(l, "\n")
			if strings.HasPrefix(l, "# file ") {
				if header == "" {
					header = l
					sb.WriteString(l)
					sb.WriteByte('\n')
					continue
				}
				pending = l
				break
			}
			sb.WriteString(l)
			sb.WriteByte('\n')
			if err != nil {
				done = true
				break
			}
		}
		if header == "" {
			return "", "", false
		}
		return header, sb.String(), true
	}
}

func nodeUtilSweepOpen(path string) (io.ReadCloser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(path, ".gz") {
		return f, nil
	}
	zr, err := gzip.NewReader(f)
	if err != nil {
		f.Close()
		return nil, err
	}
	return struct {
		io.Reader
		io.Closer
	}{zr, f}, nil
}

// nodeUtilSweepFirstDiff describes the first differing line of two chunks.
func nodeUtilSweepFirstDiff(want, got string) string {
	wl := strings.Split(want, "\n")
	gl := strings.Split(got, "\n")
	for i := 0; i < len(wl) || i < len(gl); i++ {
		var w, g string
		if i < len(wl) {
			w = wl[i]
		}
		if i < len(gl) {
			g = gl[i]
		}
		if w != g {
			return "line " + strconv.Itoa(i+1) + ":\n  java: " + w + "\n  go:   " + g + "\n  " + nodeUtilSweepFieldDiff(w, g)
		}
	}
	return ""
}

// nodeUtilSweepFieldDiff names the first differing bool position or field.
func nodeUtilSweepFieldDiff(w, g string) string {
	wf := strings.Split(w, "\t")
	gf := strings.Split(g, "\t")
	if len(wf) > 0 && len(gf) > 0 && wf[0] != gf[0] {
		a, b := wf[0], gf[0]
		for i := 0; i < len(a) && i < len(b); i++ {
			if a[i] != b[i] {
				return "first difference at char " + strconv.Itoa(i) + " of the bools/prefix"
			}
		}
		return "prefix length differs"
	}
	for i := 1; i < len(wf) || i < len(gf); i++ {
		var a, b string
		if i < len(wf) {
			a = wf[i]
		}
		if i < len(gf) {
			b = gf[i]
		}
		if a != b {
			return "first difference in field " + strconv.Itoa(i) + ": java " + a + " go " + b
		}
	}
	return ""
}

// nodeUtilSweepSplitSnippets is NodeUtilSweep.splitSnippets: a line
// "#### <name>" starts the snippet <name>, whose source is every following
// line (each with its '\n') up to the next such line.
func nodeUtilSweepSplitSnippets(text string) (names, sources []string) {
	var cur *strings.Builder
	lines := strings.Split(text, "\n")
	for i, l := range lines {
		if i == len(lines)-1 && l == "" {
			break
		}
		if strings.HasPrefix(l, "#### ") {
			if cur != nil {
				sources = append(sources, cur.String())
			}
			names = append(names, strings.TrimPrefix(l, "#### "))
			cur = &strings.Builder{}
		} else if cur != nil {
			cur.WriteString(l)
			cur.WriteByte('\n')
		}
	}
	if cur != nil {
		sources = append(sources, cur.String())
	}
	return names, sources
}

func TestNodeUtilSweepGolden(t *testing.T) {
	text, err := os.ReadFile(filepath.Join("testdata", "node_util", "sweep_snippets.txt"))
	if err != nil {
		t.Fatal(err)
	}
	names, sources := nodeUtilSweepSplitSnippets(string(text))
	rc, err := nodeUtilSweepOpen(filepath.Join("testdata", "node_util", "sweep_golden.txt.gz"))
	if err != nil {
		t.Fatal(err)
	}
	defer rc.Close()
	want, order, err := nodeUtilSweepSplit(rc)
	if err != nil {
		t.Fatal(err)
	}
	if len(order) != len(names) {
		t.Fatalf("golden has %d files, the snippets %d", len(order), len(names))
	}
	config := nodeUtilSweepConfig()
	s := newNodeUtilSweeper()
	if !s.hasConvention {
		t.Log("the probes that need the default CodingConvention are masked: needs ClosureCodingConvention")
	}
	diffs := 0
	for i, name := range names {
		var sb strings.Builder
		s.sweepFile(&sb, config, name, sources[i])
		if w := s.mask(want[nodeUtilSweepEscape(name)]); sb.String() != w {
			diffs++
			if diffs <= 50 {
				t.Errorf("%s: %s", name, nodeUtilSweepFirstDiff(w, sb.String()))
			}
		}
	}
	if diffs > 0 {
		t.Logf("%d of %d snippets differ", diffs, len(names))
	}
}

func TestNodeUtilSweepOracle(t *testing.T) {
	root, list, oracle := os.Getenv("NODEUTIL_SWEEP_ROOT"), os.Getenv("NODEUTIL_SWEEP_FILES"),
		os.Getenv("NODEUTIL_SWEEP_ORACLE")
	if root == "" || list == "" || oracle == "" {
		t.Skip("set NODEUTIL_SWEEP_ROOT, NODEUTIL_SWEEP_FILES and NODEUTIL_SWEEP_ORACLE")
	}
	maxDiffs := 20
	if v := os.Getenv("NODEUTIL_SWEEP_MAX_DIFFS"); v != "" {
		maxDiffs, _ = strconv.Atoi(v)
	}
	b, err := os.ReadFile(list)
	if err != nil {
		t.Fatal(err)
	}
	var files []string
	for _, l := range strings.Split(string(b), "\n") {
		if l != "" {
			files = append(files, l)
		}
	}
	rc, err := nodeUtilSweepOpen(oracle)
	if err != nil {
		t.Fatal(err)
	}
	defer rc.Close()
	// The oracle is read one file chunk at a time, in the order of the list
	// (the corpus output is too large to hold in memory).
	chunks := nodeUtilSweepChunks(rc)
	diffs := 0
	nodeUtilSweep(t, root, files, func(s *nodeUtilSweeper, name, got string) {
		header, w, ok := chunks()
		if !ok {
			t.Fatalf("%s: the oracle ends before this file", name)
		}
		if header != "# file "+nodeUtilSweepEscape(name) {
			t.Fatalf("%s: the oracle has %q here; the list and the oracle differ", name, header)
		}
		w = s.mask(w)
		if got != w {
			diffs++
			if diffs <= maxDiffs {
				t.Errorf("%s: %s", name, nodeUtilSweepFirstDiff(w, got))
			}
		}
	})
	t.Logf("%d files, %d with differences", len(files), diffs)
}
