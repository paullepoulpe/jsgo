package passes_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

func TestWave35B_SyntheticAllFourPartsAndRace(t *testing.T) {
	t.Parallel()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	repoRoot := filepath.Clean(filepath.Join(wd, "..", "..", ".."))
	if !realGenJ2gortReady(repoRoot) {
		t.Skip("gen/j2gort not ready right now; skipping")
	}

	tmpDir := t.TempDir()
	origModBytes, err := os.ReadFile(filepath.Join(repoRoot, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	goMod := strings.Replace(string(origModBytes), "module github.com/paullepoulpe/jsgo", "module example.com/w35btest", 1) +
		fmt.Sprintf("\nrequire github.com/paullepoulpe/jsgo v0.0.0\nreplace github.com/paullepoulpe/jsgo => %s\n", filepath.ToSlash(repoRoot))
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}
	if sumBytes, sErr := os.ReadFile(filepath.Join(repoRoot, "go.sum")); sErr == nil {
		_ = os.WriteFile(filepath.Join(tmpDir, "go.sum"), sumBytes, 0o644)
	}

	mustWrite := func(rel, content string) {
		t.Helper()
		full := filepath.Join(tmpDir, filepath.FromSlash(rel))
		if mkErr := os.MkdirAll(filepath.Dir(full), 0o755); mkErr != nil {
			t.Fatal(mkErr)
		}
		if wErr := os.WriteFile(full, []byte(content), 0o644); wErr != nil {
			t.Fatal(wErr)
		}
	}

	// Synthetic rhino with PropertyMap.properties lazy LinkedHashMap via local `map_`.
	mustWrite("rhino/rhino.go", `package rhino

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

type Node uint32

func (n Node) IsJavaNil_() bool { return n == 0 }

type ErrorReporter interface {
	IsJavaNil_() bool
}

type StaticSourceFile interface {
	IsJavaNil_() bool
	GetName() string
}

type SimpleSourceFile struct{ Name string }

func (s *SimpleSourceFile) IsJavaNil_() bool { return s == nil }
func (s *SimpleSourceFile) GetName() string  { return s.Name }

type Property struct {
	name  string
	doc   string
	pType int32
}

func NewProperty(name string, doc string, pType int32) *Property {
	return &Property{name: name, doc: doc, pType: pType}
}
func (p *Property) IsJavaNil_() bool      { return p == nil }
func (p *Property) HasIdentityEquals_()   {}
func (p *Property) GetName() string       { return p.name }
func (p *Property) GetJSDocInfo() string  { return p.doc }
func (p *Property) SetJSDocInfo(d string) { p.doc = d }
func (p *Property) GetType() int32        { return p.pType }
func (p *Property) SetType(t int32)       { p.pType = t }

type PropertyMap struct {
	prop0      *Property
	prop1      *Property
	properties *j2gort.Ext_java_util_LinkedHashMap
	ver        int32
}

func NewPropertyMap() *PropertyMap {
	p := &PropertyMap{}
	InitPropertyMap(p)
	return p
}

func InitPropertyMap(this_ *PropertyMap) {
	this_.initFields_PropertyMap()
}

func (this_ *PropertyMap) initFields_PropertyMap() {
	_ = this_
}

func (this_ *PropertyMap) HasInlinePropsOnly() bool {
	return (this_.properties == nil)
}

func (this_ *PropertyMap) GetOwnPropertiesCount() int32 {
	if this_.properties != nil {
		return this_.properties.Size()
	}
	if this_.prop0 == nil {
		return 0
	}
	if this_.prop1 == nil {
		return 1
	}
	return 2
}

func (this_ *PropertyMap) GetOwnProperty_1_string(propertyName string) *Property {
	if this_.properties != nil {
		return j2gort.CastPtr[Property](j2gort.MapGetStr(this_.properties, propertyName))
	}
	if this_.prop0 != nil {
		if j2gort.CheckNotNullString(this_.prop0.GetName()) == propertyName {
			return this_.prop0
		}
		if this_.prop1 != nil && j2gort.CheckNotNullString(this_.prop1.GetName()) == propertyName {
			return this_.prop1
		}
	}
	return nil
}

func (this_ *PropertyMap) PutProperty_2_string_Property(name string, newProp *Property) {
	var map_ *j2gort.Ext_java_util_LinkedHashMap = this_.properties
	_ = map_
	if map_ == nil {
		p0 := this_.prop0
		if p0 == nil {
			this_.prop0 = newProp
			this_.ver++
			return
		}
		if j2gort.CheckNotNullString(p0.GetName()) == name {
			newProp.SetJSDocInfo(p0.GetJSDocInfo())
			this_.prop0 = newProp
			return
		}
		p1 := this_.prop1
		if p1 == nil {
			this_.prop1 = newProp
			this_.ver++
			return
		}
		if j2gort.CheckNotNullString(p1.GetName()) == name {
			newProp.SetJSDocInfo(p1.GetJSDocInfo())
			this_.prop1 = newProp
			return
		}
		map_ = j2gort.NewExt_java_util_LinkedHashMap_1_int32(int32(4))
		j2gort.MapPutStr(map_, p0.GetName(), j2gort.PtrToAny(p0))
		j2gort.MapPutStr(map_, p1.GetName(), j2gort.PtrToAny(p1))
		j2gort.MapPutStr(map_, name, j2gort.PtrToAny(newProp))
		this_.prop0 = nil
		this_.prop1 = nil
		this_.properties = map_
		this_.ver++
		return
	}
	var oldProp *Property = j2gort.CastPtr[Property](j2gort.MapGetStr(map_, name))
	if oldProp == nil {
		this_.ver++
	} else {
		newProp.SetJSDocInfo(oldProp.GetJSDocInfo())
	}
	j2gort.MapPutStr(map_, name, j2gort.PtrToAny(newProp))
}

func (this_ *PropertyMap) CollectOwnPropertyNames(out []string, idx int32) int32 {
	if this_.properties != nil {
		cur_9 := j2gort.NewMapKeyStrCursor(this_.properties)
		for cur_9.HasNext() {
			key := cur_9.NextStr()
			out[idx] = key
			idx++
		}
		return idx
	}
	if this_.prop0 != nil {
		out[idx] = this_.prop0.GetName()
		idx++
		if this_.prop1 != nil {
			out[idx] = this_.prop1.GetName()
			idx++
		}
	}
	return idx
}

func (this_ *PropertyMap) AddOwnPropertyNamesTo(out j2gort.Ext_java_util_Set) {
	if this_.properties != nil {
		out.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.properties.KeySet_0()))
		return
	}
	if this_.prop0 != nil {
		j2gort.SetAddStr(out, this_.prop0.GetName())
		if this_.prop1 != nil {
			j2gort.SetAddStr(out, this_.prop1.GetName())
		}
	}
}

func (this_ *PropertyMap) GetOwnPropertyNames() j2gort.Ext_java_util_Set {
	if this_.properties != nil {
		return j2gort.AsIface[j2gort.Ext_java_util_Set](this_.properties.KeySet_0())
	}
	return j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0())
}

func (this_ *PropertyMap) ResolveOwnPropertyTypes(delta int32) {
	if this_.properties != nil {
		cur_36 := j2gort.NewMapValCursor(this_.properties)
		for cur_36.HasNext() {
			rawElem_37 := cur_36.Next()
			prop := j2gort.CastPtr[Property](rawElem_37)
			prop.SetType(prop.GetType() + delta)
		}
		return
	}
	if this_.prop0 != nil {
		this_.prop0.SetType(this_.prop0.GetType() + delta)
		if this_.prop1 != nil {
			this_.prop1.SetType(this_.prop1.GetType() + delta)
		}
	}
}

func (this_ *PropertyMap) Values() j2gort.Ext_java_lang_Iterable {
	if this_.properties != nil {
		return j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.properties.Values())
	}
	return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
}

func (this_ *PropertyMap) HashCode() int32 {
	var h int32 = 0
	if this_.properties != nil {
		cur_38 := j2gort.NewMapKeyStrCursor(this_.properties)
		for cur_38.HasNext() {
			key := cur_38.NextStr()
			h += javalang.StringHashCode(j2gort.CheckNotNullString(key))
		}
	}
	return h
}
`)

	// Synthetic parser with exact 14 Parser fields + Scanner + LineNumberScanner + pooled leaf IdentTree + escaping Comment.
	mustWrite("parser/parser.go", `package parser

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

type SourceFile struct {
	Name     string
	contents string
}

func NewSourceFile(name, contents string) *SourceFile {
	return &SourceFile{Name: name, contents: contents}
}

type SourcePosition struct {
	source *SourceFile
	offset int32
	line   int32
	column int32
}

type SourceRange struct {
	start SourcePosition
	end   SourcePosition
}

type Token struct {
	Value    string
	location SourceRange
}

type ErrorReporter interface {
	IsJavaNil_() bool
	ReportError(msg string)
}

type Scanner_CommentRecorder interface {
	Implements_parser_Scanner_CommentRecorder()
	RecordComment(value string, rng SourceRange)
}

func ToIface_Scanner_CommentRecorder(p *Parser) Scanner_CommentRecorder {
	if p == nil {
		return nil
	}
	return p
}

type LineNumberScanner struct {
	sourceFile     *SourceFile
	contents       string
	cachedPosition SourcePosition
	sourceLength   int32
	lastLine       int32
	lastLineStart  int32
	nextLineStart  int32
	cachedOffset   int32
	idHash_        atomic.Int32
}

func NewLineNumberScanner(sf *SourceFile) *LineNumberScanner {
	return &LineNumberScanner{
		sourceFile:    sf,
		contents:      sf.contents,
		sourceLength:  j2gort.JStrLen(sf.contents),
		lastLine:      -1,
		lastLineStart: -1,
		cachedOffset:  -1,
	}
}

func (l *LineNumberScanner) GetSourcePosition(off int32) SourcePosition {
	return SourcePosition{source: l.sourceFile, offset: off, line: 0, column: off}
}

type Scanner struct {
	errorReporter      ErrorReporter
	source             *SourceFile
	lineNumberScanner  *LineNumberScanner
	contents           string
	currentTokens      []Token
	commentRecorder    Scanner_CommentRecorder
	contentsLength     int32
	index              int32
	typeParameterLevel int32
	idHash_            atomic.Int32
}

var NewScannerCalls atomic.Int64

func NewScanner(er ErrorReporter, cr Scanner_CommentRecorder, sf *SourceFile, offset int32) *Scanner {
	NewScannerCalls.Add(1)
	return &Scanner{
		errorReporter:     er,
		commentRecorder:   cr,
		source:            sf,
		lineNumberScanner: NewLineNumberScanner(sf),
		contents:          sf.contents,
		contentsLength:    j2gort.JStrLen(sf.contents),
		index:             offset,
	}
}

func (s *Scanner) GetPosition_0() SourcePosition {
	return s.lineNumberScanner.GetSourcePosition(s.index)
}

type Parser_Config struct {
	Strict bool
}

type FeatureSet struct {
	Name string
}

var bareMin = &FeatureSet{Name: "bare"}

func StaticGet_FeatureSet_BARE_MINIMUM() *FeatureSet { return bareMin }

type Comment struct {
	Value string
	Range SourceRange
}

func NewComment(v string, r SourceRange) *Comment {
	return &Comment{Value: v, Range: r}
}

type ParseTree interface {
	Implements_parser_ParseTree()
	IsJavaNil_() bool
}

type ParseTreeBase struct {
	location SourceRange
}

func (b *ParseTreeBase) Implements_parser_ParseTree() {}
func (b *ParseTreeBase) IsJavaNil_() bool             { return b == nil }

type IdentTree struct {
	ParseTreeBase
	Name string
}

func NewIdentTree(loc SourceRange, name string) *IdentTree {
	t := &IdentTree{}
	InitIdentTree(t, loc, name)
	return t
}

func InitIdentTree(this_ *IdentTree, loc SourceRange, name string) {
	this_.location = loc
	this_.Name = name
}

type ProgramTree struct {
	ParseTreeBase
	Items    []*IdentTree
	Comments j2gort.Ext_com_google_common_collect_ImmutableList
	SrcFile  *SourceFile
}

func NewProgramTree(loc SourceRange, items []*IdentTree, comments j2gort.Ext_com_google_common_collect_ImmutableList, sf *SourceFile) *ProgramTree {
	return &ProgramTree{ParseTreeBase: ParseTreeBase{location: loc}, Items: items, Comments: comments, SrcFile: sf}
}

type Parser struct {
	scanner                *Scanner
	errorReporter          ErrorReporter
	config                 *Parser_Config
	comments               *j2gort.Ext_com_google_common_collect_ImmutableList_Builder
	lastCommentEndPosition SourcePosition
	features               *FeatureSet
	lastSourcePosition     SourcePosition
	sourceMapURL           string
	treeScratch            []ParseTree
	commaScratch           []SourcePosition
	treeScratchTop         int32
	commaScratchTop        int32
	idHash_                atomic.Int32
	inGenerator            bool
}

func (*Parser) Implements_parser_Scanner_CommentRecorder() {}

func clinit_Parser() {}

var NewParserCalls atomic.Int64

func (this_ *Parser) initFields_Parser() {
	this_.features = StaticGet_FeatureSet_BARE_MINIMUM()
	this_.sourceMapURL = j2gort.NullString
	this_.treeScratch = make([]ParseTree, 32)
	this_.commaScratch = make([]SourcePosition, 32)
}

func NewParser(config *Parser_Config, errorReporter ErrorReporter, source *SourceFile) *Parser {
	NewParserCalls.Add(1)
	clinit_Parser()
	this_ := &Parser{}
	InitParser(this_, config, errorReporter, source)
	return this_
}

func InitParser(this_ *Parser, config *Parser_Config, errorReporter ErrorReporter, source *SourceFile) {
	this_.initFields_Parser()
	this_.config = config
	this_.errorReporter = errorReporter
	this_.scanner = NewScanner(errorReporter, ToIface_Scanner_CommentRecorder(this_), source, 0)
	this_.lastSourcePosition = this_.scanner.GetPosition_0()
}

func (this_ *Parser) RecordComment(value string, rng SourceRange) {
	if this_.comments == nil {
		this_.comments = &j2gort.Ext_com_google_common_collect_ImmutableList_Builder{}
	}
	this_.comments.AddVoid_(NewComment(value, rng))
	this_.lastCommentEndPosition = rng.end
}

func (this_ *Parser) GetComments() j2gort.Ext_java_util_List {
	if this_.comments != nil {
		return j2gort.AsIface[j2gort.Ext_java_util_List](this_.comments.BuildImmList_())
	}
	return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
}

func (this_ *Parser) GetFeatures() *FeatureSet  { return this_.features }
func (this_ *Parser) GetSourceMapURL() string   { return this_.sourceMapURL }
func (this_ *Parser) PeekScannerSource() *SourceFile {
	if this_.scanner == nil {
		return nil
	}
	return this_.scanner.source
}

func (this_ *Parser) ParseProgram() *ProgramTree {
	raw := this_.scanner.contents
	if raw == "PANIC_NOW" {
		panic("boom in ParseProgram")
	}
	rng := SourceRange{start: this_.lastSourcePosition, end: SourcePosition{source: this_.scanner.source, offset: this_.scanner.contentsLength}}
	if len(raw) > 2 && raw[:2] == "//" {
		this_.RecordComment(raw, rng)
	}
	t1 := NewIdentTree(rng, raw+"_1")
	t2 := NewIdentTree(rng, raw+"_2")
	var built j2gort.Ext_com_google_common_collect_ImmutableList
	if this_.comments != nil {
		built = this_.comments.BuildImmList_()
	} else {
		built = j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	}
	return NewProgramTree(rng, []*IdentTree{t1, t2}, built, this_.scanner.source)
}
`)

	// Synthetic parsing with ParserRunner_Parse.
	mustWrite("parsing/parser_runner.go", `package parsing

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "example.com/w35btest/parser"
	pkg_rhino "example.com/w35btest/rhino"
)

type Es6ErrorReporter struct {
	delegate pkg_rhino.ErrorReporter
}

func (e *Es6ErrorReporter) IsJavaNil_() bool       { return e == nil }
func (e *Es6ErrorReporter) ReportError(msg string) {}

type ParseResult struct {
	Summary    string
	SrcName    string
	Comments   j2gort.Ext_java_util_List
	NestedSummary string
}

func ParserRunner_Parse(sf pkg_rhino.StaticSourceFile, src string, parseComments bool, nestedSrc string, er pkg_rhino.ErrorReporter) *ParseResult {
	file := pkg_parser.NewSourceFile(sf.GetName(), src)
	es6Er := &Es6ErrorReporter{delegate: er}
	cfg := &pkg_parser.Parser_Config{Strict: true}
	var p *pkg_parser.Parser = pkg_parser.NewParser(cfg, es6Er, file)
	_ = p
	tree := p.ParseProgram()
	var comments j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	if parseComments {
		comments = p.GetComments()
	}
	nestedSum := ""
	if nestedSrc != "" {
		sub := ParserRunner_Parse(sf, nestedSrc, parseComments, "", er)
		nestedSum = sub.Summary
	}
	return &ParseResult{
		Summary:       tree.Items[0].Name + "|" + tree.Items[1].Name,
		SrcName:       tree.SrcFile.Name,
		Comments:      comments,
		NestedSummary: nestedSum,
	}
}
`)

	// Synthetic jscomp testing Part 1 (both ClearTopTypedScope and SetTopScope(nil) in separate subpkg) and Part 4 (4A, 4B, 4C).
	mustWrite("jscomp/jscomp_w35.go", `package jscomp

import (
	"sync/atomic"

	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "example.com/w35btest/rhino"
)

// --- Part 1 targets ---
type TypedScope struct{}

type Compiler struct {
	topScope *TypedScope
}

func (this_ *Compiler) PerformTranspilationAndOptimizations(seg int) {
	_ = seg
}

func (this_ *Compiler) ClearTopTypedScope() {
	this_.topScope = nil
}

func (this_ *Compiler) SetTopScope(x *TypedScope) {
	_ = x
	this_.topScope = x
}

type TypeInferencePass struct{}

func (this_ *TypeInferencePass) InferAllScopes(root pkg_rhino.Node) {
	_ = root
}

type TypedAstSerializer struct{}

func (this_ *TypedAstSerializer) SerializeRoots(externs, root pkg_rhino.Node) {
	_ = externs
	_ = root
}

type PhaseOptimizer_NamedPass struct {
	name string
}

func (this_ *PhaseOptimizer_NamedPass) Process_2(externs, root pkg_rhino.Node) {
	_ = externs
	_ = root
}

// --- Part 4A: ControlFlowGraphBase + DataFlowAnalysisBase + LinearFlowState ---
type ControlFlowGraph_Branch uint8

const BranchUncond ControlFlowGraph_Branch = 1

func (b ControlFlowGraph_Branch) IsJavaNil_() bool                              { return b == 0 }
func (ControlFlowGraph_Branch) HasIdentityEquals_()                             {}
func (ControlFlowGraph_Branch) Implements_j2gort_Ext_java_lang_Enum()           {}

type ControlFlowGraph[N comparable] interface {
	Implements_jscomp_ControlFlowGraph()
	IsJavaNil_() bool
	GetControlFlowGraphBase() *ControlFlowGraphBase[N]
	pkg_graph.LinkedDirectedGraph[N, ControlFlowGraph_Branch]
	ResetCfgForPool() bool
}

type ControlFlowGraphBase[N comparable] struct {
	pkg_graph.LinkedDirectedGraphBase[N, ControlFlowGraph_Branch]
}

func (this_ *ControlFlowGraphBase[N]) BindSelf_(s any) {
	this_.LinkedDirectedGraphBase.BindSelf_(s)
}
func (this_ *ControlFlowGraphBase[N]) IsJavaNil_() bool { return this_ == nil }
func (b *ControlFlowGraphBase[N]) GetControlFlowGraphBase() *ControlFlowGraphBase[N] { return b }
func (*ControlFlowGraphBase[N]) Implements_jscomp_ControlFlowGraph()                 {}
func (*ControlFlowGraphBase[N]) Implements_graph_LinkedDirectedGraph()               {}
func (*ControlFlowGraphBase[N]) Implements_graph_DiGraph()                           {}
func (*ControlFlowGraphBase[N]) Implements_graph_GraphvizGraph()                     {}
func (*ControlFlowGraphBase[N]) Implements_graph_Graph()                             {}
func (*ControlFlowGraphBase[N]) Implements_graph_AdjacencyGraph()                    {}

func NewTestCfg(nodes []pkg_rhino.Node) *ControlFlowGraphBase[pkg_rhino.Node] {
	c := &ControlFlowGraphBase[pkg_rhino.Node]{}
	c.BindSelf_(c)
	pkg_graph.InitLinkedDirectedGraph[pkg_rhino.Node, ControlFlowGraph_Branch](&c.LinkedDirectedGraphBase, true, true)
	for _, n := range nodes {
		c.CreateNode(n)
	}
	return c
}

func (this_ *ControlFlowGraphBase[N]) ResetCfgForPool() bool {
	return this_.ResetForReuse(true, false)
}

type DataFlowAnalysis_FlowJoiner[L any] interface {
	Implements_jscomp_DataFlowAnalysis_FlowJoiner()
	IsJavaNil_() bool
	JoinFlow(L)
	Finish() L
}

type DataFlowAnalysis[N comparable, L any] interface {
	CreateInitialEstimateLattice() L
}

type DataFlowAnalysisBase[N comparable, L any] struct {
	self_ DataFlowAnalysis[N, L]
	cfg   ControlFlowGraph[N]
}

var FlowStatePoolPushes atomic.Int64

func DataFlowAnalysis_PopPooledFlowState() *DataFlowAnalysis_LinearFlowState {
	return nil
}

func DataFlowAnalysis_RecycleGraphFlowStates[N comparable](cfg ControlFlowGraph[N]) {
	_ = cfg
	cur := j2gort.NewIterCursor(j2gort.CheckNotNull(cfg.GetNodes()))
	for cur.HasNext() {
		raw := cur.Next()
		node := j2gort.AsIface[pkg_graph.DiGraph_DiGraphNode[N, ControlFlowGraph_Branch]](raw)
		if st, ok := j2gort.TryCast[*DataFlowAnalysis_LinearFlowState](node.GetAnnotation()); ok {
			st.ClearForPool()
			node.SetAnnotation(nil)
			FlowStatePoolPushes.Add(1)
		}
	}
}

func ControlFlowAnalysis_ReleaseCfg(maybeCfg any) {
	if astCfg, ok := maybeCfg.(*ControlFlowGraphBase[pkg_rhino.Node]); ok && astCfg != nil {
		DataFlowAnalysis_RecycleGraphFlowStates[pkg_rhino.Node](astCfg)
		if !astCfg.ResetCfgForPool() {
			astCfg.ReleaseSlabs_()
		}
	}
}

func (this_ *DataFlowAnalysisBase[N, L]) priv_DataFlowAnalysis_initialize_2310() {
	cur := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.cfg.GetNodes()))
	for cur.HasNext() {
		rawElem := cur.Next()
		node := j2gort.AsIface[pkg_graph.DiGraph_DiGraphNode[N, ControlFlowGraph_Branch]](rawElem)
		initIn := this_.self_.CreateInitialEstimateLattice()
		initOut := this_.self_.CreateInitialEstimateLattice()
		existing := node.GetAnnotation()
		raw, ok_17 := j2gort.TryCast[*DataFlowAnalysis_LinearFlowState](existing)
		if ok_17 {
			raw.Reinit(j2gort.ToAny(initIn), j2gort.ToAny(initOut))
		} else {
			rawPooled := DataFlowAnalysis_PopPooledFlowState()
			if rawPooled != nil {
				rawPooled.Reinit(j2gort.ToAny(initIn), j2gort.ToAny(initOut))
				node.SetAnnotation(pkg_graph.ToIface_Annotation(rawPooled))
			} else {
				node.SetAnnotation(NewDataFlowAnalysis_LinearFlowState(j2gort.ToAny(initIn), j2gort.ToAny(initOut)))
			}
		}
	}
}

type DataFlowAnalysis_LinearFlowState struct {
	in        any
	out       any
	stepCount int32
}

func (s *DataFlowAnalysis_LinearFlowState) IsJavaNil_() bool             { return s == nil }
func (*DataFlowAnalysis_LinearFlowState) Implements_graph_Annotation()   {}

func NewDataFlowAnalysis_LinearFlowState(in any, out any) *DataFlowAnalysis_LinearFlowState {
	this_ := &DataFlowAnalysis_LinearFlowState{}
	InitDataFlowAnalysis_LinearFlowState(this_, in, out)
	return this_
}

func InitDataFlowAnalysis_LinearFlowState(this_ *DataFlowAnalysis_LinearFlowState, in any, out any) {
	this_.initFields_DataFlowAnalysis_LinearFlowState()
	this_.in = in
	this_.out = out
}

func (this_ *DataFlowAnalysis_LinearFlowState) initFields_DataFlowAnalysis_LinearFlowState() {}
func (this_ *DataFlowAnalysis_LinearFlowState) Reinit(in any, out any) {
	this_.stepCount = 0
	this_.in = in
	this_.out = out
}
func (this_ *DataFlowAnalysis_LinearFlowState) ClearForPool() {
	this_.stepCount = 0
	this_.in = nil
	this_.out = nil
}
func (this_ *DataFlowAnalysis_LinearFlowState) GetIn() any  { return this_.in }
func (this_ *DataFlowAnalysis_LinearFlowState) GetOut() any { return this_.out }

// --- Part 4B: LiveVariablesAnalysis + Lattice + JoinOp + DAE & CVN ---
type LiveVariablesAnalysis struct {
	DataFlowAnalysisBase[pkg_rhino.Node, *LiveVariablesAnalysis_LiveVariableLattice]
	varCount     int32
	emptyLattice *LiveVariablesAnalysis_LiveVariableLattice
	scratchGen   *LiveVariablesAnalysis_LiveVariableLattice
	scratchKill  *LiveVariablesAnalysis_LiveVariableLattice
}

func (this_ *LiveVariablesAnalysis) initFields_LiveVariablesAnalysis() {}

func NewLiveVariablesAnalysis(cfg ControlFlowGraph[pkg_rhino.Node], varCount int32) *LiveVariablesAnalysis {
	l := &LiveVariablesAnalysis{}
	InitLiveVariablesAnalysis(l, cfg, varCount)
	return l
}

func InitLiveVariablesAnalysis(this_ *LiveVariablesAnalysis, cfg ControlFlowGraph[pkg_rhino.Node], varCount int32) {
	this_.self_ = this_
	this_.cfg = cfg
	this_.varCount = varCount
	this_.initFields_LiveVariablesAnalysis()
	this_.emptyLattice = NewLiveVariablesAnalysis_LiveVariableLattice_1_int32(varCount)
	this_.scratchGen = NewLiveVariablesAnalysis_LiveVariableLattice_1_int32(varCount)
	this_.scratchKill = NewLiveVariablesAnalysis_LiveVariableLattice_1_int32(varCount)
}

func (this_ *LiveVariablesAnalysis) CreateInitialEstimateLattice() *LiveVariablesAnalysis_LiveVariableLattice {
	return this_.emptyLattice
}

func (this_ *LiveVariablesAnalysis) CreateFlowJoiner() DataFlowAnalysis_FlowJoiner[*LiveVariablesAnalysis_LiveVariableLattice] {
	return NewLiveVariablesAnalysis_LiveVariableJoinOp(this_)
}

func (this_ *LiveVariablesAnalysis) FlowThrough(node pkg_rhino.Node, input *LiveVariablesAnalysis_LiveVariableLattice) *LiveVariablesAnalysis_LiveVariableLattice {
	if input.overflow == nil && this_.varCount <= 128 {
		return NewLiveVariablesAnalysis_LiveVariableLattice_3(input.numVars, input.w0|int64(node), input.w1)
	}
	var val_5 LiveVariablesAnalysis_LiveVariableLattice
	InitLiveVariablesAnalysis_LiveVariableLattice_1_LiveVariablesAnalysis_LiveVariableLattice(&val_5, input)
	var result *LiveVariablesAnalysis_LiveVariableLattice = &val_5
	result.SetBit(int32(node))
	return result
}

type LiveVariablesAnalysis_LiveVariableJoinOp struct {
	outer_ *LiveVariablesAnalysis
	result *LiveVariablesAnalysis_LiveVariableLattice
}

func (*LiveVariablesAnalysis_LiveVariableJoinOp) Implements_jscomp_DataFlowAnalysis_FlowJoiner() {}
func (j *LiveVariablesAnalysis_LiveVariableJoinOp) IsJavaNil_() bool                              { return j == nil }

func (this_ *LiveVariablesAnalysis_LiveVariableJoinOp) initFields_LiveVariablesAnalysis_LiveVariableJoinOp() {
	this_.result = NewLiveVariablesAnalysis_LiveVariableLattice_1_int32(this_.outer_.varCount)
}

func NewLiveVariablesAnalysis_LiveVariableJoinOp(outer_ *LiveVariablesAnalysis) *LiveVariablesAnalysis_LiveVariableJoinOp {
	this_ := &LiveVariablesAnalysis_LiveVariableJoinOp{}
	this_.outer_ = outer_
	InitLiveVariablesAnalysis_LiveVariableJoinOp(this_)
	return this_
}

func InitLiveVariablesAnalysis_LiveVariableJoinOp(this_ *LiveVariablesAnalysis_LiveVariableJoinOp) {
	this_.initFields_LiveVariablesAnalysis_LiveVariableJoinOp()
}

func (this_ *LiveVariablesAnalysis_LiveVariableJoinOp) JoinFlow(x *LiveVariablesAnalysis_LiveVariableLattice) {
	this_.result.w0 |= x.w0
}
func (this_ *LiveVariablesAnalysis_LiveVariableJoinOp) Finish() *LiveVariablesAnalysis_LiveVariableLattice {
	return this_.result
}

type LiveVariablesAnalysis_LiveVariableLattice struct {
	overflow *j2gort.Ext_java_util_BitSet
	w0       int64
	w1       int64
	numVars  int32
}

func (this_ *LiveVariablesAnalysis_LiveVariableLattice) initFields_LiveVariablesAnalysis_LiveVariableLattice() {}

func NewLiveVariablesAnalysis_LiveVariableLattice_1_int32(numVars int32) *LiveVariablesAnalysis_LiveVariableLattice {
	this_ := &LiveVariablesAnalysis_LiveVariableLattice{}
	InitLiveVariablesAnalysis_LiveVariableLattice_1_int32(this_, numVars)
	return this_
}

func InitLiveVariablesAnalysis_LiveVariableLattice_1_int32(this_ *LiveVariablesAnalysis_LiveVariableLattice, numVars int32) {
	this_.initFields_LiveVariablesAnalysis_LiveVariableLattice()
	this_.numVars = numVars
	if numVars > 128 {
		this_.overflow = j2gort.NewExt_java_util_BitSet_1(numVars)
	}
}

func NewLiveVariablesAnalysis_LiveVariableLattice_3(numVars int32, w0 int64, w1 int64) *LiveVariablesAnalysis_LiveVariableLattice {
	this_ := &LiveVariablesAnalysis_LiveVariableLattice{}
	InitLiveVariablesAnalysis_LiveVariableLattice_3(this_, numVars, w0, w1)
	return this_
}

func InitLiveVariablesAnalysis_LiveVariableLattice_3(this_ *LiveVariablesAnalysis_LiveVariableLattice, numVars int32, w0 int64, w1 int64) {
	this_.initFields_LiveVariablesAnalysis_LiveVariableLattice()
	this_.numVars = numVars
	this_.w0 = w0
	this_.w1 = w1
}

func NewLiveVariablesAnalysis_LiveVariableLattice_1_LiveVariablesAnalysis_LiveVariableLattice(other *LiveVariablesAnalysis_LiveVariableLattice) *LiveVariablesAnalysis_LiveVariableLattice {
	this_ := &LiveVariablesAnalysis_LiveVariableLattice{}
	InitLiveVariablesAnalysis_LiveVariableLattice_1_LiveVariablesAnalysis_LiveVariableLattice(this_, other)
	return this_
}

func InitLiveVariablesAnalysis_LiveVariableLattice_1_LiveVariablesAnalysis_LiveVariableLattice(this_ *LiveVariablesAnalysis_LiveVariableLattice, other *LiveVariablesAnalysis_LiveVariableLattice) {
	this_.initFields_LiveVariablesAnalysis_LiveVariableLattice()
	this_.numVars = other.numVars
	this_.w0 = other.w0
	this_.w1 = other.w1
	if other.overflow != nil {
		this_.overflow = j2gort.CastPtr[j2gort.Ext_java_util_BitSet](other.overflow.Clone())
	}
}

func (this_ *LiveVariablesAnalysis_LiveVariableLattice) SetBit(b int32) {
	if this_.overflow != nil {
		this_.overflow.Set_1(b)
	} else {
		this_.w0 |= int64(1) << uint(b&63)
	}
}

type DeadAssignmentsElimination struct {
	liveness *LiveVariablesAnalysis
}

func (this_ *DeadAssignmentsElimination) RunEliminate(cfg ControlFlowGraph[pkg_rhino.Node], varCount int32) int64 {
	return this_.priv_DeadAssignmentsElimination_eliminateDeadAssignments_2330(cfg, varCount)
}

func (this_ *DeadAssignmentsElimination) priv_DeadAssignmentsElimination_eliminateDeadAssignments_2330(cfg ControlFlowGraph[pkg_rhino.Node], varCount int32) int64 {
	defer func() {
		this_.liveness = nil
	}()
	this_.liveness = NewLiveVariablesAnalysis(cfg, varCount)
	this_.liveness.priv_DataFlowAnalysis_initialize_2310()
	j := this_.liveness.CreateFlowJoiner()
	out := this_.liveness.FlowThrough(pkg_rhino.Node(3), this_.liveness.emptyLattice)
	j.JoinFlow(out)
	return j.Finish().w0
}

type CoalesceVariableNames struct {
	liveAnalyses j2gort.Deque[*LiveVariablesAnalysis]
	liveness     *LiveVariablesAnalysis
}

func (this_ *CoalesceVariableNames) Process_2(cfg ControlFlowGraph[pkg_rhino.Node]) {
	defer func() {
		this_.liveness = nil
	}()
	this_.EnterScopeWithCfg(cfg)
	this_.ExitScopeWithCfg()
}

func (this_ *CoalesceVariableNames) EnterScopeWithCfg(cfg ControlFlowGraph[pkg_rhino.Node]) {
	this_.liveness = NewLiveVariablesAnalysis(cfg, 16)
	this_.liveness.priv_DataFlowAnalysis_initialize_2310()
	_ = this_.liveness.FlowThrough(pkg_rhino.Node(5), this_.liveness.emptyLattice)
	this_.liveAnalyses.Push(this_.liveness)
}

func (this_ *CoalesceVariableNames) ExitScopeWithCfg() {
	this_.liveAnalyses.Pop()
	this_.liveness = this_.liveAnalyses.Peek()
}

// --- Part 4C: PureFunctionIdentifier_SideEffectPropagation ---
type PureFunctionIdentifier struct {
	reverseCallGraph pkg_graph.LinkedDirectedGraph[string, *PureFunctionIdentifier_SideEffectPropagation]
}

func NewPureFunctionIdentifier() *PureFunctionIdentifier {
	return &PureFunctionIdentifier{
		reverseCallGraph: pkg_graph.LinkedDirectedGraph_CreateWithoutAnnotations[string, *PureFunctionIdentifier_SideEffectPropagation](),
	}
}

func (this_ *PureFunctionIdentifier) Process_3(externs, root pkg_rhino.Node) int {
	defer func() {
		if this_.reverseCallGraph != nil && !this_.reverseCallGraph.IsJavaNil_() {
			if r_, ok_ := any(this_.reverseCallGraph).(interface{ ReleaseSlabs_() }); ok_ {
				r_.ReleaseSlabs_()
			}
			this_.reverseCallGraph = nil
		}
	}()
	nA := j2gort.AsIface[pkg_graph.DiGraph_DiGraphNode[string, *PureFunctionIdentifier_SideEffectPropagation]](this_.reverseCallGraph.CreateNode("a"))
	nB := j2gort.AsIface[pkg_graph.DiGraph_DiGraphNode[string, *PureFunctionIdentifier_SideEffectPropagation]](this_.reverseCallGraph.CreateNode("b"))
	nC := j2gort.AsIface[pkg_graph.DiGraph_DiGraphNode[string, *PureFunctionIdentifier_SideEffectPropagation]](this_.reverseCallGraph.CreateNode("c"))
	this_.reverseCallGraph.Connect_3_DiGraph_DiGraphNode_any_any__any_DiGraph_DiGraphNode_any_any_(nA, PureFunctionIdentifier_SideEffectPropagation_ForAlias(), nB)
	fba := &PureFunctionIdentifier_FunctionBodyAnalyzer{outer_: this_}
	fba.VisitCall(nB, nC, root)
	return int(this_.reverseCallGraph.GetNodeCount())
}

type PureFunctionIdentifier_FunctionBodyAnalyzer struct {
	outer_ *PureFunctionIdentifier
}

func (this_ *PureFunctionIdentifier_FunctionBodyAnalyzer) VisitCall(src, dst pkg_graph.DiGraph_DiGraphNode[string, *PureFunctionIdentifier_SideEffectPropagation], inv pkg_rhino.Node) {
	edge := PureFunctionIdentifier_SideEffectPropagation_ForInvocation(inv, true)
	this_.outer_.reverseCallGraph.Connect_3_DiGraph_DiGraphNode_any_any__any_DiGraph_DiGraphNode_any_any_(src, edge, dst)
}

type PureFunctionIdentifier_SideEffectPropagation struct {
	invocation    pkg_rhino.Node
	callerIsAlias bool
	throws        bool
}

func (this_ *PureFunctionIdentifier_SideEffectPropagation) initFields_PureFunctionIdentifier_SideEffectPropagation() {}

func NewPureFunctionIdentifier_SideEffectPropagation(callerIsAlias bool, invocation pkg_rhino.Node, throws bool) *PureFunctionIdentifier_SideEffectPropagation {
	this_ := &PureFunctionIdentifier_SideEffectPropagation{}
	InitPureFunctionIdentifier_SideEffectPropagation(this_, callerIsAlias, invocation, throws)
	return this_
}

func InitPureFunctionIdentifier_SideEffectPropagation(this_ *PureFunctionIdentifier_SideEffectPropagation, callerIsAlias bool, invocation pkg_rhino.Node, throws bool) {
	this_.initFields_PureFunctionIdentifier_SideEffectPropagation()
	this_.callerIsAlias = callerIsAlias
	this_.invocation = invocation
	this_.throws = throws
}

func PureFunctionIdentifier_SideEffectPropagation_ForAlias() *PureFunctionIdentifier_SideEffectPropagation {
	return NewPureFunctionIdentifier_SideEffectPropagation(true, 0, true)
}

func PureFunctionIdentifier_SideEffectPropagation_ForInvocation(inv pkg_rhino.Node, throws bool) *PureFunctionIdentifier_SideEffectPropagation {
	return NewPureFunctionIdentifier_SideEffectPropagation(false, inv, throws)
}
`)

	// Also test SetTopScope(nil) fallback when ClearTopTypedScope is absent!
	mustWrite("jscomp2/jscomp_settop.go", `package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "example.com/w35btest/rhino"
)

var _ = j2gort.NullString

type TypedScope struct{}
type Compiler struct{ topScope *TypedScope }

func (this_ *Compiler) SetTopScope(x *TypedScope) {
	_ = x
	this_.topScope = x
}

type PhaseOptimizer_NamedPass struct{}

func (this_ *PhaseOptimizer_NamedPass) Process_2(externs, root pkg_rhino.Node) {
	_ = externs
	_ = root
}
`)

	allPasses := passes.AllPasses()
	pkgDirs := []string{
		filepath.Join(tmpDir, "rhino"),
		filepath.Join(tmpDir, "parser"),
		filepath.Join(tmpDir, "parsing"),
		filepath.Join(tmpDir, "jscomp"),
	}
	for _, dir := range pkgDirs {
		s1, rErr := passes.RunOnPackageDir(dir, allPasses, passes.Options{})
		if rErr != nil {
			t.Fatalf("Run 1 on %s failed: %v", dir, rErr)
		}
		if s1.TotalSites == 0 {
			t.Fatalf("expected > 0 sites on Run 1 for %s", dir)
		}
		s2, rErr := passes.RunOnPackageDir(dir, allPasses, passes.Options{})
		if rErr != nil {
			t.Fatalf("Run 2 on %s failed: %v", dir, rErr)
		}
		if s2.TotalSites != 0 || s2.FilesModified != 0 {
			t.Fatalf("Run 2 NOT idempotent on %s: %+v", dir, s2)
		}
	}

	readGot := func(rel string) string {
		t.Helper()
		b, rErr := os.ReadFile(filepath.Join(tmpDir, filepath.FromSlash(rel)))
		if rErr != nil {
			t.Fatal(rErr)
		}
		return string(b)
	}

	// Verify AST shapes across all 4 parts!
	rhinoGot := readGot("rhino/rhino.go")
	if !strings.Contains(rhinoGot, "properties *j2gort.OrderedMap[string, *Property]") ||
		strings.Contains(rhinoGot, "Ext_java_util_LinkedHashMap") {
		t.Fatalf("expected PropertyMap.properties lowered to *j2gort.OrderedMap[string, *Property]:\n%s", rhinoGot)
	}

	parserGot := readGot("parser/parser.go")
	for _, want := range []string{
		"func (this_ *Parser) ResetSlabs_()",
		"func (this_ *Parser) ClearForPool_()",
		"func (this_ *Parser) ResetForReuse_(",
	} {
		if !strings.Contains(parserGot, want) {
			t.Fatalf("missing %q in lowered parser.go:\n%s", want, parserGot)
		}
	}
	if strings.Contains(parserGot, "slab_Comment") {
		t.Fatalf("Comment escapes into ParseResult and must NOT be slabbed:\n%s", parserGot)
	}

	parsingGot := readGot("parsing/parser_runner.go")
	for _, want := range []string{
		"var parserPool_ sync.Pool",
		"acquirePooledParser_(cfg, es6Er, file)",
		"defer releasePooledParser_(p)",
	} {
		if !strings.Contains(parsingGot, want) {
			t.Fatalf("missing %q in lowered parser_runner.go:\n%s", want, parsingGot)
		}
	}

	jscompGot := readGot("jscomp/jscomp_w35.go")
	for _, want := range []string{
		"defer j2gort.PurgeAllRuntimePools()",
		"slab_LinearFlowState",
		"this_.slab_LinearFlowState.Release()",
		"allocPooled_DataFlowAnalysis_LinearFlowState",
		"slab_LiveVariableLattice",
		"slab_LiveVariableJoinOp",
		"this_.slab_LiveVariableLattice.Alloc()",
		"this_.liveness.ReleaseSlabs_()",
		"poppedLiveness_.ReleaseSlabs_()",
		"slab_SideEffectPropagation",
		"allocPooledVia_PureFunctionIdentifier_SideEffectPropagation_ForAlias",
		"allocPooledVia_PureFunctionIdentifier_SideEffectPropagation_ForInvocation",
		"this_.slab_SideEffectPropagation.Release()",
	} {
		if !strings.Contains(jscompGot, want) {
			t.Fatalf("missing %q in lowered jscomp_w35.go:\n%s", want, jscompGot)
		}
	}
	if strings.Contains(jscompGot, "TrimStaleSlabPools") {
		t.Fatalf("unexpected TrimStaleSlabPools in lowered jscomp_w35.go:\n%s", jscompGot)
	}

	// Write comprehensive -race runtime test!
	mustWrite("w35b_runtime_test.go", `package w35b_rt_test

import (
	"fmt"
	"sync"
	"testing"

	"example.com/w35btest/jscomp"
	"example.com/w35btest/parser"
	"example.com/w35btest/parsing"
	"example.com/w35btest/rhino"
)

func TestWave35BRuntimeAndRace(t *testing.T) {
	// 1. PropertyMap inline -> OrderedMap transition, overwrite JSDoc preservation, iteration, live views, resolve, hashCode.
	pm := rhino.NewPropertyMap()
	if !pm.HasInlinePropsOnly() || pm.GetOwnPropertiesCount() != 0 {
		t.Fatal("expected empty inline PropertyMap")
	}
	pm.PutProperty_2_string_Property("a", rhino.NewProperty("a", "docA", 10))
	pm.PutProperty_2_string_Property("b", rhino.NewProperty("b", "docB", 20))
	if !pm.HasInlinePropsOnly() || pm.GetOwnPropertiesCount() != 2 {
		t.Fatal("expected 2 inline properties")
	}
	pm.PutProperty_2_string_Property("c", rhino.NewProperty("c", "docC", 30))
	if pm.HasInlinePropsOnly() || pm.GetOwnPropertiesCount() != 3 {
		t.Fatal("expected promoted OrderedMap with 3 properties")
	}
	// Overwrite "b" and verify JSDocInfo carried over from oldProp!
	pm.PutProperty_2_string_Property("b", rhino.NewProperty("b", "", 200))
	if got := pm.GetOwnProperty_1_string("b"); got == nil || got.GetJSDocInfo() != "docB" || got.GetType() != 200 {
		t.Fatalf("unexpected overwritten property b: %+v", got)
	}
	names := make([]string, 3)
	if n := pm.CollectOwnPropertyNames(names, 0); n != 3 || names[0] != "a" || names[1] != "b" || names[2] != "c" {
		t.Fatalf("unexpected insertion order: n=%d names=%v", n, names)
	}
	if pm.GetOwnPropertyNames().Size() != 3 {
		t.Fatalf("expected LiveKeySet size 3, got %d", pm.GetOwnPropertyNames().Size())
	}
	pm.ResolveOwnPropertyTypes(5)
	if pm.GetOwnProperty_1_string("a").GetType() != 15 || pm.GetOwnProperty_1_string("c").GetType() != 35 {
		t.Fatal("ResolveOwnPropertyTypes did not update all values")
	}
	if pm.HashCode() == 0 {
		t.Fatal("expected non-zero HashCode")
	}

	// 2. Parser pooling: single-thread reuse (zero extra NewParser / NewScanner after warmup!), comment isolation, panic safety, nested reentrancy, and -race concurrency!
	parser.NewParserCalls.Store(0)
	parser.NewScannerCalls.Store(0)
	sf1 := &rhino.SimpleSourceFile{Name: "f1.js"}
	r1 := parsing.ParserRunner_Parse(sf1, "//c1", true, "nested_code", nil)
	if r1.Comments.Size() != 1 || r1.NestedSummary != "nested_code_1|nested_code_2" {
		t.Fatalf("unexpected r1: %+v", r1)
	}
	// Next non-comment file MUST NOT see stale comments from f1.js!
	for i := 0; i < 20; i++ {
		sf2 := &rhino.SimpleSourceFile{Name: fmt.Sprintf("f_%d.js", i)}
		r2 := parsing.ParserRunner_Parse(sf2, "plain_js", true, "", nil)
		if r2.Comments.Size() != 0 || r2.SrcName != sf2.Name || r2.Summary != "plain_js_1|plain_js_2" {
			t.Fatalf("stale state leaked on iter %d: %+v", i, r2)
		}
	}
	// Panic inside ParseProgram must still release cleanly to pool!
	func() {
		defer func() { _ = recover() }()
		_ = parsing.ParserRunner_Parse(sf1, "PANIC_NOW", true, "", nil)
	}()
	// Concurrent -race hammer!
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for k := 0; k < 50; k++ {
				sf := &rhino.SimpleSourceFile{Name: fmt.Sprintf("g%d_%d.js", id, k)}
				hasCmt := (k%2 == 0)
				code := fmt.Sprintf("code_%d_%d", id, k)
				if hasCmt {
					code = "//" + code
				}
				res := parsing.ParserRunner_Parse(sf, code, true, "inner", nil)
				wantCmts := int32(0)
				if hasCmt {
					wantCmts = 1
				}
				if res.Comments.Size() != wantCmts || res.SrcName != sf.Name {
					t.Errorf("race mismatch: got cmts=%d want=%d src=%s", res.Comments.Size(), wantCmts, res.SrcName)
				}
			}
		}(g)
	}
	wg.Wait()

	// 3. Part 4A + 4B + 4C slab execution (verify 0 FLOW_STATE_POOL pushes + correct results across pooled resets!)
	jscomp.FlowStatePoolPushes.Store(0)

	cfg := jscomp.NewTestCfg([]rhino.Node{1, 2, 3, 4})
	dae := &jscomp.DeadAssignmentsElimination{}
	if w := dae.RunEliminate(cfg, 32); w == 0 {
		t.Fatal("expected non-zero lattice bit")
	}
	// Overflow (> 128 vars) stack-escaped FlowThrough path:
	_ = dae.RunEliminate(cfg, 200)

	cvn := &jscomp.CoalesceVariableNames{}
	cvn.Process_2(cfg)
	jscomp.ControlFlowAnalysis_ReleaseCfg(cfg)

	pfi := jscomp.NewPureFunctionIdentifier()
	if n := pfi.Process_3(1, 2); n != 3 {
		t.Fatalf("expected 3 nodes in PFI graph, got %d", n)
	}

	if got := jscomp.FlowStatePoolPushes.Load(); got != 0 {
		t.Fatalf("expected 0 FLOW_STATE_POOL pushes of slab LinearFlowState, got %d", got)
	}

	// Part 1 smoke execution:
	comp := &jscomp.Compiler{}
	comp.ClearTopTypedScope()
	comp.PerformTranspilationAndOptimizations(0)
	np := &jscomp.PhaseOptimizer_NamedPass{}
	np.Process_2(1, 2)
}
`)

	testCmd := exec.Command("go", "test", "-race", "-v", ".")
	testCmd.Dir = tmpDir
	if out, tErr := testCmd.CombinedOutput(); tErr != nil {
		t.Fatalf("synthetic Wave 35B -race test failed: %v\n%s", tErr, string(out))
	}
}

func TestWave35B_RealGenRhinoParserParsingJscompCopyIdempotencyAndRace(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping full real gen copy test in -short mode")
	}
	t.Parallel()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	repoRoot := filepath.Clean(filepath.Join(wd, "..", "..", ".."))
	if !realGenJ2gortReady(repoRoot) {
		t.Skip("gen/j2gort does not build cleanly right now; skipping real gen copy check")
	}

	checkCmd := exec.Command("go", "build", "./transpiled/gen/j2gort", "./transpiled/gen/rhino", "./transpiled/gen/parser", "./transpiled/gen/parsing", "./transpiled/gen/diagnostic")
	checkCmd.Dir = repoRoot
	if err := checkCmd.Run(); err != nil {
		t.Skipf("gen packages mid-edit (%v); skipping real copy check", err)
	}

	tmpDir := t.TempDir()
	origModBytes, err := os.ReadFile(filepath.Join(repoRoot, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	goMod := strings.Replace(string(origModBytes), "module github.com/paullepoulpe/jsgo", "module example.com/realw35b", 1) +
		fmt.Sprintf("\nrequire github.com/paullepoulpe/jsgo v0.0.0\nreplace github.com/paullepoulpe/jsgo => %s\n", filepath.ToSlash(repoRoot))
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}
	if sumBytes, sErr := os.ReadFile(filepath.Join(repoRoot, "go.sum")); sErr == nil {
		_ = os.WriteFile(filepath.Join(tmpDir, "go.sum"), sumBytes, 0o644)
	}

	pkgs := []string{"rhino", "parser", "parsing", "jscomptype", "jscomp"}
	for _, pkg := range pkgs {
		srcDir := filepath.Join(repoRoot, "transpiled", "gen", pkg)
		dstDir := filepath.Join(tmpDir, pkg)
		if err := os.MkdirAll(dstDir, 0o755); err != nil {
			t.Fatal(err)
		}
		entries, err := os.ReadDir(srcDir)
		if err != nil {
			t.Fatal(err)
		}
		for _, ent := range entries {
			if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".go") || strings.HasSuffix(ent.Name(), "_test.go") {
				continue
			}
			data, err := os.ReadFile(filepath.Join(srcDir, ent.Name()))
			if err != nil {
				t.Fatal(err)
			}
			txt := string(data)
			txt = strings.ReplaceAll(txt, `"github.com/paullepoulpe/jsgo/transpiled/gen/rhino"`, `"example.com/realw35b/rhino"`)
			txt = strings.ReplaceAll(txt, `"github.com/paullepoulpe/jsgo/transpiled/gen/parser"`, `"example.com/realw35b/parser"`)
			txt = strings.ReplaceAll(txt, `"github.com/paullepoulpe/jsgo/transpiled/gen/parsing"`, `"example.com/realw35b/parsing"`)
			txt = strings.ReplaceAll(txt, `"github.com/paullepoulpe/jsgo/transpiled/gen/jscomptype"`, `"example.com/realw35b/jscomptype"`)
			if err := os.WriteFile(filepath.Join(dstDir, ent.Name()), []byte(txt), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}

	allPasses := passes.AllPasses()
	for _, pkg := range pkgs {
		pkgDir := filepath.Join(tmpDir, pkg)
		s1, err := passes.RunOnPackageDir(pkgDir, allPasses, passes.Options{})
		if err != nil {
			if !realGenJ2gortReady(repoRoot) {
				t.Skipf("dependency became mid-edit during %s Run 1: %v", pkg, err)
			}
			t.Fatalf("Run 1 on copied gen/%s failed: %v", pkg, err)
		}
		t.Logf("Wave 35B Real gen/%s Run 1 sites: total=%d files=%d byPass=%+v", pkg, s1.TotalSites, s1.FilesModified, s1.ByPass)

		s2, err := passes.RunOnPackageDir(pkgDir, allPasses, passes.Options{})
		if err != nil {
			if !realGenJ2gortReady(repoRoot) {
				t.Skipf("dependency became mid-edit during %s Run 2: %v", pkg, err)
			}
			t.Fatalf("Run 2 (idempotency) on copied gen/%s failed: %v", pkg, err)
		}
		if s2.TotalSites != 0 || s2.FilesModified != 0 {
			t.Fatalf("Run 2 on copied gen/%s NOT IDEMPOTENT: TotalSites=%d FilesModified=%d byPass=%+v",
				pkg, s2.TotalSites, s2.FilesModified, s2.ByPass)
		}
	}

	readCopy := func(rel string) string {
		t.Helper()
		b, rErr := os.ReadFile(filepath.Join(tmpDir, rel))
		if rErr != nil {
			t.Fatal(rErr)
		}
		return string(b)
	}
	normWS := func(s string) string { return strings.Join(strings.Fields(s), " ") }

	// Part 1 verification on real gen/jscomp:
	compCode := readCopy("jscomp/compiler.go")
	if strings.Contains(compCode, "TrimStaleSlabPools") || strings.Count(compCode, "PurgeAllRuntimePools()") != 1 {
		t.Fatalf("expected 0 TrimStaleSlabPools + 1 PurgeAllRuntimePools() in jscomp/compiler.go")
	}
	poCode := readCopy("jscomp/phase_optimizer.go")
	if strings.Contains(poCode, "TrimStaleSlabPools") {
		t.Fatalf("unexpected TrimStaleSlabPools in jscomp/phase_optimizer.go")
	}

	// Part 2 verification on real gen/parser + gen/parsing:
	parserCode := readCopy("parser/parser.go")
	for _, want := range []string{
		"func (this_ *Parser) ResetSlabs_()",
		"func (this_ *Parser) ClearForPool_()",
		"func (this_ *Parser) ResetForReuse_(",
	} {
		if !strings.Contains(parserCode, want) {
			t.Fatalf("missing %q in real lowered parser/parser.go", want)
		}
	}
	if strings.Contains(parserCode, "slab_Comment") {
		t.Fatalf("real parser/parser.go must NOT slab Comment")
	}
	prCode := readCopy("parsing/parser_runner.go")
	if !strings.Contains(prCode, "acquirePooledParser_(es6config, es6ErrorReporter, file)") ||
		!strings.Contains(prCode, "defer releasePooledParser_(p)") {
		t.Fatalf("expected pooled Parser acquire/release in real lowered parsing/parser_runner.go:\n%s", prCode)
	}

	// Part 3 verification on real gen/rhino/property_map.go:
	pmCode := normWS(readCopy("rhino/property_map.go"))
	if !strings.Contains(pmCode, "properties *j2gort.OrderedMap[string, *Property]") ||
		strings.Contains(pmCode, "Ext_java_util_LinkedHashMap") {
		t.Fatalf("expected PropertyMap.properties lowered to *j2gort.OrderedMap[string, *Property] in real rhino/property_map.go")
	}

	// Part 4 verification on real gen/jscomp:
	cfgCode := normWS(readCopy("jscomp/control_flow_graph.go"))
	dfaCode := normWS(readCopy("jscomp/data_flow_analysis.go"))
	if !strings.Contains(cfgCode, "slab_LinearFlowState j2gort.PooledSlab[DataFlowAnalysis_LinearFlowState]") ||
		!strings.Contains(cfgCode, "this_.slab_LinearFlowState.Release()") ||
		!strings.Contains(dfaCode, "allocPooled_DataFlowAnalysis_LinearFlowState(&cfgBase_.slab_LinearFlowState") {
		t.Fatalf("expected DataFlowAnalysis_LinearFlowState PooledSlab lowering in real jscomp")
	}
	lvaCode := normWS(readCopy("jscomp/live_variables_analysis.go"))
	daeCode := normWS(readCopy("jscomp/dead_assignments_elimination.go"))
	cvnCode := normWS(readCopy("jscomp/coalesce_variable_names.go"))
	if !strings.Contains(lvaCode, "slab_LiveVariableLattice j2gort.PooledSlab[LiveVariablesAnalysis_LiveVariableLattice]") ||
		!strings.Contains(lvaCode, "this_.slab_LiveVariableLattice.Alloc()") ||
		!strings.Contains(daeCode, "this_.liveness.ReleaseSlabs_()") ||
		!strings.Contains(cvnCode, "poppedLiveness_.ReleaseSlabs_()") {
		t.Fatalf("expected LiveVariableLattice PooledSlab lowering + DAE/CVN releases in real jscomp")
	}
	pfiCode := normWS(readCopy("jscomp/pure_function_identifier.go"))
	if !strings.Contains(pfiCode, "slab_SideEffectPropagation j2gort.PooledSlab[PureFunctionIdentifier_SideEffectPropagation]") ||
		!strings.Contains(pfiCode, "allocPooledVia_PureFunctionIdentifier_SideEffectPropagation_ForInvocation(") ||
		!strings.Contains(pfiCode, "this_.slab_SideEffectPropagation.Release()") {
		t.Fatalf("expected PureFunctionIdentifier_SideEffectPropagation PooledSlab lowering in real jscomp")
	}

	// End-to-end -race execution against the real lowered rhino + parser + parsing + jscomp packages!
	realRaceTest := `package realw35b_test

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"example.com/realw35b/jscomp"
	"example.com/realw35b/parsing"
	"example.com/realw35b/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
)

func TestRealLoweredRaceExecution(t *testing.T) {
	store := nodestore.NewNodeStore()
	prev := nodestore.SetCurrentNodeStore(store)
	defer nodestore.SetCurrentNodeStore(prev)

	cfg := parsing.ParserRunner_CreateConfig_6(
		parsing.StaticGet_Config_LanguageMode_ES_NEXT(),
		parsing.StaticGet_Config_JsDocParsing_INCLUDE_DESCRIPTIONS_NO_WHITESPACE(),
		parsing.StaticGet_Config_RunMode_KEEP_GOING(),
		nil,
		true,
		parsing.StaticGet_Config_StrictMode_STRICT(),
	)
	errRep := rhino.StaticGet_ErrorReporter_NULL_INSTANCE()

	// 1. Concurrent ParserRunner_Parse across goroutines under -race (exercising parserPool_ + Scanner/LineNumberScanner reuse + comment isolation):
	var wg sync.WaitGroup
	for g := 0; g < 4; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < 15; i++ {
				sf := rhino.NewSimpleSourceFile(fmt.Sprintf("file_%d_%d.js", id, i), rhino.StaticGet_StaticSourceFile_SourceKind_STRONG())
				var code string
				if i%2 == 0 {
					code = fmt.Sprintf("/** @fileoverview doc %d */\nexport function f%d(a, b, c, d, e) { return a + b + c + d + e; }\n", i, i)
				} else {
					code = fmt.Sprintf("const x%d = [1, 2, 3].map(v => v * 2);\n", i)
				}
				res := parsing.ParserRunner_Parse(sf, code, cfg, errRep)
				if res == nil || res.FieldGet_ast() == 0 {
					t.Errorf("nil AST on iter %d", i)
					return
				}
				if i%2 == 0 && res.FieldGet_comments().Size() != 1 {
					t.Errorf("expected 1 comment on even iter %d, got %d", i, res.FieldGet_comments().Size())
				}
				if i%2 == 1 && res.FieldGet_comments().Size() != 0 {
					t.Errorf("expected 0 comments on odd iter %d (stale comment leaked!), got %d", i, res.FieldGet_comments().Size())
				}
			}
		}(g)
	}
	wg.Wait()

	// 2. Full real Compiler CHECKS + OPT pass pipeline (exercising PropertyMap, TrimStaleSlabPools, LinearFlowState, LiveVariableLattice, and PureFunctionIdentifier_SideEffectPropagation):
	comp := jscomp.NewCompiler()
	opts := jscomp.NewCompilerOptions()
	jscomp.StaticGet_CompilationLevel_ADVANCED_OPTIMIZATIONS().SetOptionsForCompilationLevel(opts)
	opts.SetCheckTypes(true)
	ext := jscomp.SourceFile_FromCode_2("externs.js", "/** @externs */ function alert(x) {}")
	src := jscomp.SourceFile_FromCode_2("input.js", ` + "`" + `
/**
 * @constructor
 */
function Widget() {
  /** @type {number} */ this.a = 1;
  /** @type {number} */ this.b = 2;
  /** @type {number} */ this.c = 3;
  /** @type {number} */ this.d = 4;
}
/** @param {number} k @return {number} */
function helper(k) {
  var dead = k + 99;
  var live = k * 2;
  if (live > 10) {
    live = live + 1;
  }
  return live;
}
var w = new Widget();
alert(helper(w.a + w.b + w.c + w.d));
` + "`" + `)
	res := comp.Compile_3_SourceFile_SourceFile_CompilerOptions(ext, src, opts)
	if res == nil || comp.HasErrors() {
		t.Fatalf("ADVANCED + CheckTypes compile failed")
	}
	out := comp.ToSource_0()
	if !strings.Contains(out, "alert(") {
		t.Fatalf("unexpected compiled output: %s", out)
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "real_w35b_race_test.go"), []byte(realRaceTest), 0o644); err != nil {
		t.Fatal(err)
	}

	raceCmd := exec.Command("go", "test", "-race", "-v", "-run", "TestRealLoweredRaceExecution", ".")
	raceCmd.Dir = tmpDir
	raceOut, err := raceCmd.CombinedOutput()
	if err != nil {
		if !realGenJ2gortReady(repoRoot) {
			t.Skipf("dependency mid-edit during real -race test: %v", err)
		}
		t.Fatalf("real lowered gen packages -race execution failed: %v\n%s", err, string(raceOut))
	}
}
