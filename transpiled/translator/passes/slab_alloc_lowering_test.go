package passes_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

func setupWave18SynthWorkspace(t *testing.T) (tmpDir, repoRoot string) {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	repoRoot = filepath.Clean(filepath.Join(wd, "..", "..", ".."))
	tmpDir = t.TempDir()

	goMod := "module example.com/w18slabtest\n\ngo 1.24.0\n\nrequire (\n\tgithub.com/paullepoulpe/jsgo v0.0.0\n)\n\nreplace github.com/paullepoulpe/jsgo => " + filepath.ToSlash(repoRoot) + "\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}
	if sumBytes, err := os.ReadFile(filepath.Join(repoRoot, "go.sum")); err == nil {
		_ = os.WriteFile(filepath.Join(tmpDir, "go.sum"), sumBytes, 0o644)
	}

	j2Dir := filepath.Join(tmpDir, "j2gort")
	if err := os.MkdirAll(j2Dir, 0o755); err != nil {
		t.Fatal(err)
	}
	j2Stub := `package j2gort

import (
	"fmt"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
)

var _ = nodestore.NilNode

const NullString = "\x00"

var SlabReleaseCalls atomic.Int64
var PooledSlabReleaseCalls atomic.Int64

func ThrowNPE() {
	panic(javalang.NewNullPointerException())
}

func CheckStateStr(b bool, msg string) {
	if !b {
		panic(javalang.NewIllegalStateException(msg))
	}
}

func CheckStateMsg(b bool, msg any) {
	if !b {
		panic(javalang.NewIllegalStateException(fmt.Sprint(msg)))
	}
}

type BoxDouble struct {
	Val float64
}

func ValueOfDouble(v float64) *BoxDouble { return &BoxDouble{Val: v} }
func UnboxDouble(b *BoxDouble) float64 {
	if b == nil {
		ThrowNPE()
	}
	return b.Val
}

type Slab[T any] struct {
	buf []T
}

func (s *Slab[T]) Alloc() *T {
	i := len(s.buf)
	if i == cap(s.buf) {
		n := cap(s.buf) * 2
		if n == 0 {
			n = 4
		} else if n > 64 {
			n = 64
		}
		c := make([]T, 1, n)
		s.buf = c
		return &c[0]
	}
	b := s.buf[:i+1]
	s.buf = b
	return &b[i]
}

func (s *Slab[T]) Release() {
	if s != nil {
		s.buf = nil
		SlabReleaseCalls.Add(1)
	}
}

type PooledSlab[T any] struct {
	buf []T
}

func (s *PooledSlab[T]) Alloc() *T {
	i := len(s.buf)
	if i == cap(s.buf) {
		c := make([]T, 1, 64)
		s.buf = c
		return &c[0]
	}
	b := s.buf[:i+1]
	s.buf = b
	return &b[i]
}

func (s *PooledSlab[T]) Release() {
	if s != nil {
		s.buf = nil
		SlabReleaseCalls.Add(1)
		PooledSlabReleaseCalls.Add(1)
	}
}

type Deque[T any] struct {
	items []T
}

func (d *Deque[T]) Push(v T) {
	d.items = append(d.items, v)
}

func (d *Deque[T]) Pop() T {
	last := len(d.items) - 1
	v := d.items[last]
	var zero T
	d.items[last] = zero
	d.items = d.items[:last]
	return v
}
`
	if err := os.WriteFile(filepath.Join(j2Dir, "runtime.go"), []byte(j2Stub), 0o644); err != nil {
		t.Fatal(err)
	}

	rhinoDir := filepath.Join(tmpDir, "rhino")
	if err := os.MkdirAll(rhinoDir, 0o755); err != nil {
		t.Fatal(err)
	}
	rhinoStub := `package rhino

import (
	"fmt"
	"math"
	"strings"
	"sync/atomic"

	"example.com/w18slabtest/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
)

type Node nodestore.Node
type Token uint8

func StaticGet_Token_POS() Token       { return 22 }
func StaticGet_Token_ADD() Token       { return 14 }
func StaticGet_Token_NAME() Token      { return 33 }
func StaticGet_Token_NUMBER() Token    { return 34 }
func StaticGet_Token_STRINGLIT() Token { return 36 }
func StaticGet_Token_HOOK() Token      { return 68 }
func StaticGet_Token_EMPTY() Token     { return 89 }
func StaticGet_Token_BLOCK() Token     { return 91 }
func StaticGet_Token_SCRIPT() Token    { return 95 }

type StaticSourceFile interface {
	IsJavaNil_() bool
	GetName() string
}

type SimpleSourceFile struct {
	name   string
	handle atomic.Uint32
}

func NewSimpleSourceFile(name string) *SimpleSourceFile      { return &SimpleSourceFile{name: name} }
func (s *SimpleSourceFile) IsJavaNil_() bool                 { return s == nil }
func (s *SimpleSourceFile) GetName() string                  { return s.name }
func (s *SimpleSourceFile) SourceFileHandlePtr_() *atomic.Uint32 { return &s.handle }

func checkDetached(n Node) *nodestore.NodeData {
	if n == 0 {
		panic(javalang.NewNullPointerException())
	}
	rec := nodestore.Node(n).Rec()
	if (rec.Parent | rec.Next | rec.Previous) != 0 {
		panic(javalang.NewIllegalStateException("child not detached"))
	}
	return rec
}

func NewNode_1(tok Token) Node {
	raw, rec := nodestore.AllocNodeFastRec()
	rec.Token = uint8(tok) - 1
	rec.LinenoCharno = -1
	return Node(raw)
}

func NewNode_2(tok Token, c1 Node) Node {
	raw, rec := nodestore.AllocNodeFastRec()
	rec.Token = uint8(tok) - 1
	rec.LinenoCharno = -1
	rec.First = nodestore.Node(c1)
	r1 := checkDetached(c1)
	r1.Previous = nodestore.Node(c1)
	r1.Parent = raw
	return Node(raw)
}

func NewNode_3(tok Token, c1, c2 Node) Node {
	raw, rec := nodestore.AllocNodeFastRec()
	rec.Token = uint8(tok) - 1
	rec.LinenoCharno = -1
	rec.First = nodestore.Node(c1)
	r1 := checkDetached(c1)
	r1.Next = nodestore.Node(c2)
	r1.Previous = nodestore.Node(c2)
	r1.Parent = raw
	r2 := checkDetached(c2)
	r2.Previous = nodestore.Node(c1)
	r2.Parent = raw
	return Node(raw)
}

func NewNode_4(tok Token, c1, c2, c3 Node) Node {
	raw, rec := nodestore.AllocNodeFastRec()
	rec.Token = uint8(tok) - 1
	rec.LinenoCharno = -1
	rec.First = nodestore.Node(c1)
	r1 := checkDetached(c1)
	r1.Next = nodestore.Node(c2)
	r1.Previous = nodestore.Node(c3)
	r1.Parent = raw
	r2 := checkDetached(c2)
	r2.Next = nodestore.Node(c3)
	r2.Previous = nodestore.Node(c1)
	r2.Parent = raw
	r3 := checkDetached(c3)
	r3.Previous = nodestore.Node(c2)
	r3.Parent = raw
	return Node(raw)
}

func Node_NewString_2(tok Token, str string) Node {
	raw, rec := nodestore.AllocNodeFastRec()
	rec.Token = uint8(tok) - 1
	rec.LinenoCharno = -1
	rec.KindFlags = 3
	if str == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	rec.Payload = uint32(nodestore.Intern(str))
	return Node(raw)
}

func IR_String(s string) Node {
	return Node_NewString_2(StaticGet_Token_STRINGLIT(), s)
}

func IR_Number(d float64) Node {
	if math.IsNaN(d) || d < 0 || (d == 0 && 1/d < 0) {
		j2gort.CheckStateMsg(false, d)
	}
	raw, rec := nodestore.AllocNodeFastRec()
	rec.Token = uint8(StaticGet_Token_NUMBER()) - 1
	rec.LinenoCharno = -1
	rec.KindFlags = 1
	nodestore.Node(raw).SetNumberVal(d)
	return Node(raw)
}

func (n Node) ClonePropsFrom(other Node) Node {
	if !nodestore.Node(n).ClonePropsFrom(nodestore.Node(other)) {
		j2gort.CheckStateStr(false, "Node has existing properties.")
	}
	return n
}

func (n Node) CloneNode_0() Node {
	c := NewNode_1(Token(nodestore.Node(n).Token() + 1))
	return c.ClonePropsFrom(n)
}

func (n Node) SetStaticSourceFile(sf StaticSourceFile) {
	nodestore.Node(n).SetSourceFileObj(sf)
}

func (n Node) GetSourceFileName() string {
	if sf, ok := nodestore.Node(n).SourceFileObj().(StaticSourceFile); ok && sf != nil {
		return sf.GetName()
	}
	return ""
}

func (n Node) SetIsInClosureUnawareSubtree(v bool) {
	pf := nodestore.Node(n).PropFlags()
	if v {
		pf |= 16777216
	} else {
		pf &^= 16777216
	}
	nodestore.Node(n).SetPropFlags(pf)
}

func (n Node) GetIsInClosureUnawareSubtree() bool {
	return nodestore.Node(n).PropFlags()&16777216 != 0
}

func (n Node) GetFirstChild() Node { return Node(nodestore.Node(n).FirstChild()) }
func (n Node) GetLastChild() Node {
	f := nodestore.Node(n).FirstChild()
	if f == 0 {
		return 0
	}
	return Node(f.Previous())
}

func (n Node) Summary() string {
	var b strings.Builder
	var walk func(nodestore.Node)
	walk = func(cur nodestore.Node) {
		fmt.Fprintf(&b, "(tok=%d kind=%d sf=%d pf=%d", cur.Token(), cur.KindFlags(), cur.SourceFile(), cur.PropFlags())
		switch cur.KindFlags() {
		case 1:
			fmt.Fprintf(&b, " num=%v", cur.NumberVal())
		case 3:
			fmt.Fprintf(&b, " str=%s", nodestore.String(nodestore.StringID(cur.Payload())))
		}
		for c := cur.FirstChild(); c != 0; c = c.Next() {
			b.WriteByte(' ')
			walk(c)
		}
		b.WriteByte(')')
	}
	walk(nodestore.Node(n))
	return b.String()
}
`
	if err := os.WriteFile(filepath.Join(rhinoDir, "node.go"), []byte(rhinoStub), 0o644); err != nil {
		t.Fatal(err)
	}
	return tmpDir, repoRoot
}

func realGenJ2gortReady(repoRoot string) bool {
	cmd := exec.Command("go", "build", "./gen/j2gort", "./gen/rhino")
	cmd.Dir = filepath.Join(repoRoot, "transpiled")
	return cmd.Run() == nil
}

func TestWave18_IRFactoryNodeAllocatorAndParserReleaseSlabs(t *testing.T) {
	tmpDir, repoRoot := setupWave18SynthWorkspace(t)
	allPasses := passes.AllPasses()

	makeSynthParser := func(pkgName string) string {
		return `package ` + pkgName + `

import "example.com/w18slabtest/j2gort"

type ParseTree interface {
	Implements_` + pkgName + `_ParseTree()
	IsJavaNil_() bool
}

type IdentTree struct {
	Name string
}

func (*IdentTree) Implements_` + pkgName + `_ParseTree() {}
func (t *IdentTree) IsJavaNil_() bool                    { return t == nil }
func (t *IdentTree) initFields_IdentTree()               { _ = t }

func NewIdentTree(name string) *IdentTree {
	this_ := &IdentTree{}
	InitIdentTree(this_, name)
	return this_
}

func InitIdentTree(this_ *IdentTree, name string) {
	this_.initFields_IdentTree()
	this_.Name = name
}

type ProgramTree struct {
	Items []*IdentTree
}

func (*ProgramTree) Implements_` + pkgName + `_ParseTree() {}
func (t *ProgramTree) IsJavaNil_() bool                     { return t == nil }

type Parser struct {
	src string
}

func NewParser(src string) *Parser {
	return &Parser{src: src}
}

func (this_ *Parser) ParseProgram() *ProgramTree {
	a := NewIdentTree(this_.src + "_a")
	if this_.src == "PANIC" {
		panic("parse_error")
	}
	b := NewIdentTree(this_.src + "_b")
	return &ProgramTree{Items: []*IdentTree{a, b}}
}

var _ = j2gort.NullString
`
	}

	makeSynthParsing := func(pkgName, parserImport string) map[string]string {
		return map[string]string{
			"ir_factory.go": `package ` + pkgName + `

import (
	"example.com/w18slabtest/j2gort"
	pkg_parser "` + parserImport + `"
	pkg_rhino "example.com/w18slabtest/rhino"
)

type IRFactory struct {
	templateNode               pkg_rhino.Node
	defaultTemplateNode        pkg_rhino.Node
	closureUnawareTemplateNode pkg_rhino.Node
	resultNode                 pkg_rhino.Node
}

func (this_ *IRFactory) initFields_IRFactory() {
	_ = this_
}

func NewIRFactory(sf pkg_rhino.StaticSourceFile) *IRFactory {
	this_ := &IRFactory{}
	InitIRFactory(this_, sf)
	return this_
}

func InitIRFactory(this_ *IRFactory, sf pkg_rhino.StaticSourceFile) {
	_ = sf
	this_.initFields_IRFactory()
	tmpl := pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_SCRIPT())
	tmpl.SetStaticSourceFile(sf)
	this_.defaultTemplateNode = tmpl
	this_.templateNode = this_.defaultTemplateNode
	cu := tmpl.CloneNode_0()
	if sf != nil && !sf.IsJavaNil_() {
		cu.SetIsInClosureUnawareSubtree(true)
	}
	this_.closureUnawareTemplateNode = cu
}

func IRFactory_TransformTree(tree *pkg_parser.ProgramTree, sf pkg_rhino.StaticSourceFile, closureUnaware bool) *IRFactory {
	var val_1 IRFactory
	InitIRFactory(&val_1, sf)
	var irFactory *IRFactory = &val_1
	_ = irFactory
	irFactory.resultNode = irFactory.Transform(tree, closureUnaware)
	return irFactory
}

func (this_ *IRFactory) Transform(tree *pkg_parser.ProgramTree, closureUnaware bool) pkg_rhino.Node {
	if closureUnaware {
		this_.templateNode = this_.closureUnawareTemplateNode
	}
	s1 := this_.NewStringNode_1(tree.Items[0].Name)
	s2 := this_.NewStringNode_2(pkg_rhino.StaticGet_Token_NAME(), tree.Items[1].Name)
	n0 := this_.NewNumberNode(j2gort.ValueOfDouble(0.0))
	n1 := this_.NewNumberNode(j2gort.ValueOfDouble(42.5))
	u := this_.NewNode_2(pkg_rhino.StaticGet_Token_POS(), n0)
	b := this_.NewNode_3(pkg_rhino.StaticGet_Token_ADD(), s1, s2)
	c := this_.NewNode_4(pkg_rhino.StaticGet_Token_HOOK(), u, b, n1)
	if closureUnaware {
		this_.templateNode = this_.defaultTemplateNode
	}
	empty := this_.NewNode_1(pkg_rhino.StaticGet_Token_EMPTY())
	return this_.NewNode_3(pkg_rhino.StaticGet_Token_BLOCK(), c, empty)
}

func (this_ *IRFactory) NewNode_1(type_ pkg_rhino.Token) pkg_rhino.Node {
	_ = type_
	return pkg_rhino.NewNode_1(type_).ClonePropsFrom(this_.templateNode)
}

func (this_ *IRFactory) NewNode_2(type_ pkg_rhino.Token, child1 pkg_rhino.Node) pkg_rhino.Node {
	_ = type_
	_ = child1
	return pkg_rhino.NewNode_2(type_, child1).ClonePropsFrom(this_.templateNode)
}

func (this_ *IRFactory) NewNode_3(type_ pkg_rhino.Token, child1 pkg_rhino.Node, child2 pkg_rhino.Node) pkg_rhino.Node {
	_ = type_
	_ = child1
	_ = child2
	return pkg_rhino.NewNode_3(type_, child1, child2).ClonePropsFrom(this_.templateNode)
}

func (this_ *IRFactory) NewNode_4(type_ pkg_rhino.Token, child1 pkg_rhino.Node, child2 pkg_rhino.Node, child3 pkg_rhino.Node) pkg_rhino.Node {
	_ = type_
	_ = child1
	_ = child2
	_ = child3
	return pkg_rhino.NewNode_4(type_, child1, child2, child3).ClonePropsFrom(this_.templateNode)
}

func (this_ *IRFactory) NewStringNode_1(value string) pkg_rhino.Node {
	_ = value
	return pkg_rhino.IR_String(value).ClonePropsFrom(this_.templateNode)
}

func (this_ *IRFactory) NewStringNode_2(type_ pkg_rhino.Token, value string) pkg_rhino.Node {
	_ = type_
	_ = value
	return pkg_rhino.Node_NewString_2(type_, value).ClonePropsFrom(this_.templateNode)
}

func (this_ *IRFactory) NewNumberNode(value *j2gort.BoxDouble) pkg_rhino.Node {
	_ = value
	return pkg_rhino.IR_Number(j2gort.UnboxDouble(value)).ClonePropsFrom(this_.templateNode)
}
`,
			"parser_runner.go": `package ` + pkgName + `

import (
	pkg_parser "` + parserImport + `"
	pkg_rhino "example.com/w18slabtest/rhino"
)

func ParserRunner_Parse(src string, sf pkg_rhino.StaticSourceFile, closureUnaware bool) pkg_rhino.Node {
	var out pkg_rhino.Node
	func() {
		var p *pkg_parser.Parser = pkg_parser.NewParser(src)
		_ = p
		tree := p.ParseProgram()
		f := IRFactory_TransformTree(tree, sf, closureUnaware)
		out = f.resultNode
	}()
	return out
}
`,
		}
	}

	origParserDir := filepath.Join(tmpDir, "origparser")
	origParsingDir := filepath.Join(tmpDir, "origparsing")
	optParserDir := filepath.Join(tmpDir, "optparser")
	optParsingDir := filepath.Join(tmpDir, "optparsing")
	for _, d := range []string{origParserDir, origParsingDir, optParserDir, optParsingDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(origParserDir, "parser.go"), []byte(makeSynthParser("origparser")), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(optParserDir, "parser.go"), []byte(makeSynthParser("optparser")), 0o644); err != nil {
		t.Fatal(err)
	}
	for fn, src := range makeSynthParsing("origparsing", "example.com/w18slabtest/origparser") {
		if err := os.WriteFile(filepath.Join(origParsingDir, fn), []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for fn, src := range makeSynthParsing("optparsing", "example.com/w18slabtest/optparser") {
		if err := os.WriteFile(filepath.Join(optParsingDir, fn), []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	stats1, err := passes.Run([]string{optParserDir, optParsingDir}, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 on optparser+optparsing failed: %v", err)
	}
	if stats1.ByPass["slab_alloc_lowering"] < 10 {
		t.Fatalf("expected >= 10 slab_alloc_lowering sites across optparser+optparsing, got %+v", stats1.ByPass)
	}
	// Run 2 MUST be 100% idempotent (0 sites, 0 files modified)!
	stats2, err := passes.Run([]string{optParserDir, optParsingDir}, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 on optparser+optparsing failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites & 0 modified files on Run 2, got %+v", stats2)
	}

	optParserSrc, err := os.ReadFile(filepath.Join(optParserDir, "parser.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(optParserSrc), "func (this_ *Parser) ReleaseSlabs_()") ||
		!strings.Contains(string(optParserSrc), "this_.slab_IdentTree.Release()") ||
		!strings.Contains(string(optParserSrc), "j2gort.PooledSlab[IdentTree]") ||
		!strings.Contains(string(optParserSrc), "allocPooled_IdentTree") {
		t.Fatalf("expected PooledSlab + ReleaseSlabs_() on optparser.Parser:\n%s", string(optParserSrc))
	}

	optRunnerSrc, err := os.ReadFile(filepath.Join(optParsingDir, "parser_runner.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(optRunnerSrc), "defer p.ReleaseSlabs_()") {
		t.Fatalf("expected defer p.ReleaseSlabs_() in optparsing/parser_runner.go:\n%s", string(optRunnerSrc))
	}

	optIRSrc, err := os.ReadFile(filepath.Join(optParsingDir, "ir_factory.go"))
	if err != nil {
		t.Fatal(err)
	}
	nsCheckSrc, _ := os.ReadFile(filepath.Join(repoRoot, "transpiled", "shims", "nodestore", "node_store.go"))
	wantNode1Call := "this_.nodeAlloc_.AllocNode1("
	if strings.Contains(string(nsCheckSrc), "func (a *NodeAllocator) AllocNode1Cached(") {
		wantNode1Call = "this_.nodeAlloc_.AllocNode1Cached("
	}
	wantStrNodeCall := "this_.nodeAlloc_.AllocStringNode("
	if strings.Contains(string(nsCheckSrc), "func (a *NodeAllocator) AllocStringNodeCached(") {
		wantStrNodeCall = "this_.nodeAlloc_.AllocStringNodeCached("
	}
	for _, want := range []string{
		"nodeAlloc_",
		"nodestore.NodeAllocator",
		"this_.nodeAlloc_.InitLocal()",
		"this_.nodeAlloc_.SetTemplate(nodestore.Node(this_.templateNode))",
		"defer irFactory.nodeAlloc_.Flush()",
		wantNode1Call,
		"this_.nodeAlloc_.AllocNode2(",
		"this_.nodeAlloc_.AllocNode3(",
		"this_.nodeAlloc_.AllocNode4(",
		wantStrNodeCall,
		"this_.nodeAlloc_.AllocNumberNode(",
	} {
		if !strings.Contains(string(optIRSrc), want) {
			t.Fatalf("missing %q in lowered ir_factory.go:\n%s", want, string(optIRSrc))
		}
	}
	if strings.Count(string(optIRSrc), "SetTemplate(nodestore.Node(this_.templateNode))") != 3 {
		t.Fatalf("expected 3 SetTemplate calls (InitIRFactory + 2 in Transform), got %d:\n%s",
			strings.Count(string(optIRSrc), "SetTemplate(nodestore.Node(this_.templateNode))"), string(optIRSrc))
	}

	// Execute differential runtime parity test between origparsing and optparsing!
	diffTestSrc := `package w18_ir_test

import (
	"math"
	"testing"

	"example.com/w18slabtest/j2gort"
	"example.com/w18slabtest/optparsing"
	"example.com/w18slabtest/origparsing"
	pkg_rhino "example.com/w18slabtest/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
)

func TestIRFactoryDifferentialParity(t *testing.T) {
	store := nodestore.NewNodeStore()
	prev := nodestore.SetCurrentNodeStore(store)
	defer nodestore.SetCurrentNodeStore(prev)

	sf := pkg_rhino.NewSimpleSourceFile("test_input.js")
	j2gort.SlabReleaseCalls.Store(0)

	for _, cu := range []bool{false, true} {
		before := store.NodeCount()
		nOpt := optparsing.ParserRunner_Parse("alpha", sf, cu)
		optUsed := store.NodeCount() - before

		nOrig := origparsing.ParserRunner_Parse("alpha", sf, cu)
		if nOpt.Summary() != nOrig.Summary() {
			t.Fatalf("AST mismatch (cu=%v):\n opt=%s\norig=%s", cu, nOpt.Summary(), nOrig.Summary())
		}
		if nOpt.GetSourceFileName() != "test_input.js" {
			t.Fatalf("sourceFileName = %q, want test_input.js", nOpt.GetSourceFileName())
		}
		innerHook := nOpt.GetFirstChild()
		outerEmpty := nOpt.GetLastChild()
		if innerHook.GetIsInClosureUnawareSubtree() != cu {
			t.Fatalf("innerHook closureUnaware = %v, want %v", innerHook.GetIsInClosureUnawareSubtree(), cu)
		}
		if outerEmpty.GetIsInClosureUnawareSubtree() {
			t.Fatalf("outerEmpty should have restored defaultTemplateNode (closureUnaware=false)")
		}
		// Flush() must rewind unused trailing IDs (2 templates + 9 AST nodes = 11 total, << 256!).
		if optUsed != 11 {
			t.Fatalf("Flush() did not rewind unused local batch IDs: optUsed=%d, want 11", optUsed)
		}
	}

	if got := j2gort.SlabReleaseCalls.Load(); got != 2 {
		t.Fatalf("expected 2 Slab.Release() calls after 2 parses, got %d", got)
	}

	// Verify defer p.ReleaseSlabs_() ALSO executes when ParseProgram panics!
	func() {
		defer func() { _ = recover() }()
		_ = optparsing.ParserRunner_Parse("PANIC", sf, false)
	}()
	if got := j2gort.SlabReleaseCalls.Load(); got != 3 {
		t.Fatalf("expected defer p.ReleaseSlabs_() to fire on panic (count=3), got %d", got)
	}

	// NullString in NewStringNode_1 / NewStringNode_2 must throw NPE identically.
	f := optparsing.NewIRFactory(sf)
	func() {
		defer func() {
			r := recover()
			if _, ok := r.(*javalang.NullPointerException); !ok {
				t.Fatalf("expected *javalang.NullPointerException on NullString, got %T (%v)", r, r)
			}
		}()
		f.NewStringNode_1(j2gort.NullString)
	}()

	// -0.0 in NewNumberNode must delegate to IR_Number and throw IllegalStateException identically.
	func() {
		defer func() {
			r := recover()
			if _, ok := r.(*javalang.IllegalStateException); !ok {
				t.Fatalf("expected *javalang.IllegalStateException on -0.0, got %T (%v)", r, r)
			}
		}()
		f.NewNumberNode(j2gort.ValueOfDouble(math.Copysign(0, -1)))
	}()
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "w18_ir_test.go"), []byte(diffTestSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "test", "-v", "-run", "TestIRFactoryDifferentialParity", ".")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test TestIRFactoryDifferentialParity failed: %v\n%s", err, string(out))
	}

	// Also verify DryRun on real gen/parsing when gen/j2gort is not mid-edit!
	if realGenJ2gortReady(repoRoot) {
		realParsingDir := filepath.Join(repoRoot, "transpiled", "gen", "parsing")
		irSrc, _ := os.ReadFile(filepath.Join(realParsingDir, "ir_factory.go"))
		jsdocSrc, _ := os.ReadFile(filepath.Join(realParsingDir, "js_doc_info_parser.go"))
		prSrc, _ := os.ReadFile(filepath.Join(realParsingDir, "parser_runner.go"))
		realParserGo, _ := os.ReadFile(filepath.Join(repoRoot, "transpiled", "gen", "parser", "parser.go"))
		nsSrc, _ := os.ReadFile(filepath.Join(repoRoot, "transpiled", "shims", "nodestore", "node_store.go"))
		hasAllocIntern := strings.Contains(string(nsSrc), "func (a *NodeAllocator) Intern(")
		hasAllocNode1Cached := strings.Contains(string(nsSrc), "func (a *NodeAllocator) AllocNode1Cached(")
		hasAllocStrCached := strings.Contains(string(nsSrc), "func (a *NodeAllocator) AllocStringNodeCached(")
		hasParserResetForReuse := strings.Contains(string(realParserGo), "ResetForReuse_(")
		parsingPreLowered := strings.Contains(string(irSrc), "nodeAlloc_") &&
			strings.Contains(string(jsdocSrc), "nodeAlloc_") &&
			(!hasAllocIntern || strings.Contains(string(irSrc), "nodeAlloc_.Intern(")) &&
			(!hasAllocNode1Cached || strings.Contains(string(irSrc), "AllocNode1Cached(")) &&
			(!hasAllocStrCached || strings.Contains(string(irSrc), "AllocStringNodeCached(")) &&
			(!hasParserResetForReuse || strings.Contains(string(prSrc), "acquirePooledParser_("))
		dryStats, err := passes.RunOnPackageDir(realParsingDir, []passes.Pass{passes.NewSlabAllocLoweringPass()}, passes.Options{DryRun: true})
		if err != nil {
			t.Fatalf("DryRun slab_alloc_lowering on real gen/parsing failed: %v", err)
		}
		if !parsingPreLowered && dryStats.ByPass["slab_alloc_lowering"] < 1 {
			t.Fatalf("expected >= 1 slab_alloc_lowering sites on real gen/parsing, got %+v", dryStats.ByPass)
		}
		if parsingPreLowered && dryStats.ByPass["slab_alloc_lowering"] != 0 {
			t.Fatalf("expected 0 slab_alloc_lowering sites on already-lowered gen/parsing, got %+v", dryStats.ByPass)
		}
	}
}

func TestWave18_ReferenceCollectorAndSyntacticScopeVarSlabs(t *testing.T) {
	tmpDir, repoRoot := setupWave18SynthWorkspace(t)
	allPasses := passes.AllPasses()

	makeScopeAndRefPkg := func(pkgName string) map[string]string {
		return map[string]string{
			"scope_and_var.go": `package ` + pkgName + `

import "example.com/w18slabtest/j2gort"

type Var struct {
	name  string
	scope *Scope
	index int32
}

func (this_ *Var) initFields_Var() { _ = this_ }

func NewVar(name string, scope *Scope, index int32) *Var {
	this_ := &Var{}
	InitVar(this_, name, scope, index)
	return this_
}

func InitVar(this_ *Var, name string, scope *Scope, index int32) {
	this_.initFields_Var()
	this_.name = name
	this_.scope = scope
	this_.index = index
}

func Var_CreateImplicitGoogNamespace(name string, scope *Scope) *Var {
	_ = name
	_ = scope
	return NewVar(name, scope, -1)
}

type Scope struct {
	parent *Scope
	depth  int32
	vars   map[string]*Var
	rootID int32
}

func (this_ *Scope) BindSelf_(s any) { _ = s }
func (this_ *Scope) initFields_Scope() {
	this_.vars = make(map[string]*Var)
}

func NewScope_1(rootID int32) *Scope {
	this_ := &Scope{}
	this_.BindSelf_(this_)
	InitScope_1(this_, rootID)
	return this_
}

func InitScope_1(this_ *Scope, rootID int32) {
	this_.initFields_Scope()
	this_.parent = nil
	this_.depth = 0
	this_.rootID = rootID
}

func NewScope_2(parent *Scope, rootID int32) *Scope {
	this_ := &Scope{}
	this_.BindSelf_(this_)
	InitScope_2(this_, parent, rootID)
	return this_
}

func InitScope_2(this_ *Scope, parent *Scope, rootID int32) {
	this_.initFields_Scope()
	this_.parent = parent
	this_.depth = parent.depth + 1
	this_.rootID = rootID
}

func Scope_CreateGlobalScope(rootID int32) *Scope {
	_ = rootID
	return NewScope_1(rootID)
}

func Scope_CreateChildScope(parent *Scope, rootID int32) *Scope {
	_ = parent
	_ = rootID
	return NewScope_2(parent, rootID)
}

func (this_ *Scope) Declare(name string) *Var {
	v := NewVar(name, this_, int32(len(this_.vars)))
	this_.vars[name] = v
	return v
}

func (this_ *Scope) DeclareImplicitGoogNamespaceIfAbsent(name string) *Var {
	v := this_.vars[name]
	if v == nil {
		v = Var_CreateImplicitGoogNamespace(name, this_)
		this_.vars[name] = v
	}
	return v
}

func (this_ *Scope) MakeImplicitVar(name string) *Var {
	return NewVar(name, this_, -1)
}

type SyntacticScopeCreator struct {
	scannerInUse bool
}

func NewSyntacticScopeCreator() *SyntacticScopeCreator {
	return &SyntacticScopeCreator{}
}

func (this_ *SyntacticScopeCreator) CreateScope_2_Node_Scope(n int32, parent *Scope) *Scope {
	var scope *Scope
	if parent == nil {
		scope = Scope_CreateGlobalScope(n)
	} else {
		scope = Scope_CreateChildScope(parent, n)
	}
	return scope
}

var _ = j2gort.NullString
`,
			"reference_collector.go": `package ` + pkgName + `

import "example.com/w18slabtest/j2gort"

type Reference struct {
	nodeID int32
	scope  *Scope
}

func (this_ *Reference) initFields_Reference() { _ = this_ }

func NewReference_2(nodeID int32, scope *Scope) *Reference {
	this_ := &Reference{}
	InitReference_2(this_, nodeID, scope)
	return this_
}

func InitReference_2(this_ *Reference, nodeID int32, scope *Scope) {
	this_.initFields_Reference()
	this_.nodeID = nodeID
	this_.scope = scope
}

type ReferenceCollector struct {
	referenceMap map[*Var][]*Reference
}

func (*ReferenceCollector) Implements_` + pkgName + `_CompilerPass() {}

func NewReferenceCollector() *ReferenceCollector {
	return &ReferenceCollector{referenceMap: make(map[*Var][]*Reference)}
}

func (this_ *ReferenceCollector) AddRefStackInit(v *Var, nodeID int32, s *Scope) *Reference {
	var val_10 Reference
	InitReference_2(&val_10, nodeID, s)
	var reference *Reference = &val_10
	this_.referenceMap[v] = append(this_.referenceMap[v], reference)
	return reference
}

func (this_ *ReferenceCollector) AddRefDirectCtor(v *Var, nodeID int32, s *Scope) *Reference {
	ref := NewReference_2(nodeID, s)
	this_.referenceMap[v] = append(this_.referenceMap[v], ref)
	return ref
}

var _ = j2gort.NullString
`,
		}
	}

	origDir := filepath.Join(tmpDir, "origjscomp")
	optDir := filepath.Join(tmpDir, "optjscomp")
	for _, d := range []string{origDir, optDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	for fn, src := range makeScopeAndRefPkg("origjscomp") {
		if err := os.WriteFile(filepath.Join(origDir, fn), []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for fn, src := range makeScopeAndRefPkg("optjscomp") {
		if err := os.WriteFile(filepath.Join(optDir, fn), []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	stats1, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 on optjscomp failed: %v", err)
	}
	if stats1.ByPass["slab_alloc_lowering"] != 7 {
		t.Fatalf("expected 7 slab_alloc_lowering sites on optjscomp, got %d (%+v)", stats1.ByPass["slab_alloc_lowering"], stats1.ByPass)
	}
	stats2, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 on optjscomp failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites & 0 modified files on optjscomp Run 2, got %+v", stats2)
	}

	scopeCode, err := os.ReadFile(filepath.Join(optDir, "scope_and_var.go"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"slab_Scope",
		"j2gort.Slab[Scope]",
		"slab_Var",
		"j2gort.Slab[Var]",
		"varSlab_",
		"*j2gort.Slab[Var]",
		"scope.varSlab_ = &this_.slab_Var",
		"this_.varSlab_ = parent.varSlab_",
		"allocOpt_Var(this_.varSlab_",
		"allocOptVia_Var_CreateImplicitGoogNamespace(this_.varSlab_",
	} {
		if !strings.Contains(string(scopeCode), want) {
			t.Fatalf("missing %q in lowered scope_and_var.go:\n%s", want, string(scopeCode))
		}
	}

	refCode, err := os.ReadFile(filepath.Join(optDir, "reference_collector.go"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"slab_Reference j2gort.Slab[Reference]",
		"val_10 := this_.slab_Reference.Alloc()",
		"alloc_Reference_2(&this_.slab_Reference",
	} {
		if !strings.Contains(string(refCode), want) {
			t.Fatalf("missing %q in lowered reference_collector.go:\n%s", want, string(refCode))
		}
	}

	rtTest := `package w18_jscomp_test

import (
	"testing"

	"example.com/w18slabtest/optjscomp"
	"example.com/w18slabtest/origjscomp"
)

func TestScopeVarAndRefParity(t *testing.T) {
	sscOpt := optjscomp.NewSyntacticScopeCreator()
	gOpt := sscOpt.CreateScope_2_Node_Scope(1, nil)
	vGlobal := gOpt.Declare("window")
	vGoog := gOpt.DeclareImplicitGoogNamespaceIfAbsent("goog.dom")
	vImpl := gOpt.MakeImplicitVar("arguments")
	cOpt := sscOpt.CreateScope_2_Node_Scope(2, gOpt)
	vLocal := cOpt.Declare("x")
	cDirect := optjscomp.Scope_CreateChildScope(gOpt, 3)
	vDirect := cDirect.Declare("y")

	// Standalone Scope_CreateGlobalScope (varSlab_ == nil) must still work cleanly without nil deref!
	standalone := optjscomp.Scope_CreateGlobalScope(99)
	vStandalone := standalone.Declare("solo")
	if vStandalone == nil {
		t.Fatal("standalone Declare returned nil")
	}

	rcOpt := optjscomp.NewReferenceCollector()
	r1 := rcOpt.AddRefStackInit(vGlobal, 10, gOpt)
	r2 := rcOpt.AddRefDirectCtor(vLocal, 20, cOpt)
	if r1 == r2 || vGlobal == vGoog || vGoog == vImpl || vLocal == vDirect {
		t.Fatal("expected distinct slab-allocated pointers")
	}

	sscOrig := origjscomp.NewSyntacticScopeCreator()
	gOrig := sscOrig.CreateScope_2_Node_Scope(1, nil)
	_ = gOrig.Declare("window")
	_ = gOrig.DeclareImplicitGoogNamespaceIfAbsent("goog.dom")
	_ = gOrig.MakeImplicitVar("arguments")
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "w18_jscomp_test.go"), []byte(rtTest), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "test", "-v", "-run", "TestScopeVarAndRefParity", ".")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test TestScopeVarAndRefParity failed: %v\n%s", err, string(out))
	}

	if realGenJ2gortReady(repoRoot) {
		realJscompDir := filepath.Join(repoRoot, "transpiled", "gen", "jscomp")
		rcSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "reference_collector.go"))
		colRtSrc, _ := os.ReadFile(filepath.Join(repoRoot, "transpiled", "gen", "j2gort", "collections_runtime.go"))
		extRtSrc, _ := os.ReadFile(filepath.Join(repoRoot, "transpiled", "gen", "j2gort", "ext_collections_rt.go"))
		ovRtSrc, _ := os.ReadFile(filepath.Join(repoRoot, "transpiled", "gen", "j2gort", "lowering_helpers_overlay.go"))
		j2AllSrc := string(colRtSrc) + "\n" + string(extRtSrc) + "\n" + string(ovRtSrc)
		tipSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "type_inference_pass.go"))
		rucSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "remove_unused_code.go"))
		sscSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "syntactic_scope_creator.go"))
		ivSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "inline_variables.go"))
		iolSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "inline_object_literals.go"))
		fsivSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "flow_sensitive_inline_variables.go"))
		daeSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "dead_assignments_elimination.go"))
		gcSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "guarded_callback.go"))
		gnSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "global_namespace.go"))
		rvSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "rename_vars.go"))
		compSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "compiler.go"))
		tasSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "typed_ast_serializer.go"))
		cttcSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "convert_types_to_colors.go"))
		ccrcSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "cross_chunk_reference_collector.go"))
		nsDir := filepath.Join(repoRoot, "transpiled", "shims", "nodestore")
		nsEntries, _ := os.ReadDir(nsDir)
		var nsAll strings.Builder
		for _, ne := range nsEntries {
			if !ne.IsDir() && strings.HasSuffix(ne.Name(), ".go") && !strings.HasSuffix(ne.Name(), "_test.go") {
				b, _ := os.ReadFile(filepath.Join(nsDir, ne.Name()))
				nsAll.Write(b)
				nsAll.WriteByte('\n')
			}
		}
		hasPooledInGen := strings.Contains(j2AllSrc, "type PooledSlab[")
		hasReleaseWithInGen := strings.Contains(j2AllSrc, "ReleaseWith(")
		hasSlabResetInGen := strings.Contains(j2AllSrc, "func (s *PooledSlab[T]) Reset()")
		hasAcquireCustomCapInGen := strings.Contains(j2AllSrc, "AcquireCustomOrderedMapCap[")
		hasReleaseOrderedMapInGen := strings.Contains(j2AllSrc, "func ReleaseOrderedMap[")
		hasPooledSliceInGen := strings.Contains(j2AllSrc, "func ReleasePooledSlice[")
		hasTrimStaleInGen := strings.Contains(j2AllSrc, "func TrimStaleSlabPools(")
		hasPurgeInGen := strings.Contains(j2AllSrc, "func PurgeAllRuntimePools(")
		hasRelJSTypesInNS := strings.Contains(nsAll.String(), "func ReleaseCurrentStoreJSTypes(")
		hasInternOwnedInNS := strings.Contains(nsAll.String(), "func InternOwned(")
		_ = hasTrimStaleInGen
		_ = hasPurgeInGen
		mcRealSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "minimized_condition.go"))
		spRealSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "string_pool.go"))
		pmcSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "peephole_minimize_conditions.go"))
		ocSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "optimize_calls.go"))
		mdnuRealSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "make_declared_names_unique.go"))
		ifRealSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "inline_functions.go"))
		pdRealSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "process_defines.go"))
		dpaeRealSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "dead_property_assignment_elimination.go"))
		tscRealSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "typed_scope_creator.go"))
		ftbRealSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "function_type_builder.go"))
		tiRealSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "type_inference.go"))
		poRealSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "phase_optimizer.go"))
		cfgRealSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "control_flow_graph.go"))
		lvaRealSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "live_variables_analysis.go"))
		pfiRealSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "pure_function_identifier.go"))
		bbRealSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "basic_block.go"))
		scopeRealSrc, _ := os.ReadFile(filepath.Join(realJscompDir, "scope.go"))
		_ = ftbRealSrc
		jscompPreLowered := strings.Contains(string(rcSrc), "slab_Reference") &&
			strings.Contains(string(tscRealSrc), "acquireNormalScopeBuilder_") &&
			strings.Contains(string(tiRealSrc), ".InitBuilder_(&fb_") &&
			(!hasPooledInGen || ((strings.Contains(string(tipSrc), "ReleaseSlabs_") || strings.Contains(string(tipSrc), "ControlFlowAnalysis_ReleaseCfg")) &&
				strings.Contains(string(rucSrc), "ReleaseSlabs_") &&
				strings.Contains(string(rucSrc), "slab_RemoveUnusedCode_VarInfo") &&
				strings.Contains(string(ivSrc), "ReleaseSlabs_") &&
				(strings.Contains(string(iolSrc), "ReleaseSlabs_") || strings.Contains(string(iolSrc), "slab_")) &&
				(strings.Contains(string(fsivSrc), "ReleaseSlabs_") || strings.Contains(string(fsivSrc), "ControlFlowAnalysis_ReleaseCfg")) &&
				strings.Contains(string(daeSrc), "ReleaseSlabs_") &&
				strings.Contains(string(rvSrc), "ReleaseSlabs_()") &&
				strings.Contains(string(ccrcSrc), "func (this_ *CrossChunkReferenceCollector) ReleaseSlabs_()") &&
				strings.Contains(string(gcSrc), "func (this_ *GuardedCallbackBase) ReleaseSlabs_()") &&
				strings.Contains(string(gcSrc), "slabCtxCount_") &&
				strings.Contains(string(pmcSrc), "slab_MinimizedCondition_MeasuredNode") &&
				strings.Contains(string(ocSrc), "slab_OptimizeCalls_NodeList") &&
				strings.Contains(string(ocSrc), "NewSyntacticScopeCreator_1") &&
				strings.Contains(string(mdnuRealSrc), "slab_ContextualRenamer") &&
				strings.Contains(string(ifRealSrc), "slab_InlineFunctions_Reference") &&
				strings.Contains(string(dpaeRealSrc), "slab_DeadPropertyAssignmentElimination_Property") &&
				strings.Contains(string(cfgRealSrc), "slab_LinearFlowState") &&
				strings.Contains(string(lvaRealSrc), "slab_LiveVariableLattice") &&
				strings.Contains(string(pfiRealSrc), "slab_SideEffectPropagation"))) &&
			(!hasSlabResetInGen || strings.Contains(string(gcSrc), ".slab_GuardedCallback_Context.Reset()")) &&
			(!hasReleaseWithInGen || (strings.Contains(string(sscSrc), "ReleaseSlabs_") &&
				strings.Contains(string(rcSrc), "ReleaseSlabs_") &&
				strings.Contains(string(gnSrc), "func (this_ *GlobalNamespace) ReleaseSlabs_()") &&
				strings.Contains(string(gnSrc), "slab_GlobalNamespace_Ref"))) &&
			(!hasReleaseOrderedMapInGen || (strings.Contains(string(sscSrc), "ReleaseOrderedMap(") &&
				strings.Contains(string(gnSrc), "ReleaseOrderedMap("))) &&
			(!hasPooledSliceInGen || (strings.Contains(string(rcSrc), "ReleasePooledSlice(") &&
				strings.Contains(string(rucSrc), "ReleasePooledSlice("))) &&
			(!hasPurgeInGen || (strings.Count(string(compSrc), "PurgeAllRuntimePools(") == 1 &&
				strings.Count(string(tipSrc), "PurgeAllRuntimePools(") == 1 &&
				!strings.Contains(string(pdRealSrc), "PurgeAllRuntimePools(") &&
				strings.Contains(string(tasSrc), "PurgeAllRuntimePools("))) &&
			(!hasRelJSTypesInNS || strings.Contains(string(cttcSrc), "ReleaseCurrentStoreJSTypes(")) &&
			(!hasInternOwnedInNS || strings.Contains(string(spRealSrc), "InternOwned(")) &&
			strings.Contains(string(ivSrc), "scratchVarExpert_") &&
			strings.Contains(string(mcRealSrc), "lazyTok_") &&
			strings.Contains(string(bbRealSrc), "p_.GetToken()") &&
			!strings.Contains(string(rucSrc), "InitRemoveUnusedCode_VarInfo_1(val_") &&
			strings.Index(string(scopeRealSrc), "varSlab_") < strings.Index(string(scopeRealSrc), "depth ") &&
			(!hasAcquireCustomCapInGen || strings.Contains(string(rcSrc), "AcquireCustomOrderedMapCap") || strings.Contains(string(rcSrc), "AcquireOrderedMapCap"))
		_ = hasTrimStaleInGen
		_ = poRealSrc
		dryStats, err := passes.RunOnPackageDir(realJscompDir, []passes.Pass{passes.NewSlabAllocLoweringPass()}, passes.Options{DryRun: true})
		if err != nil {
			t.Fatalf("DryRun slab_alloc_lowering on real gen/jscomp failed: %v", err)
		}
		if !jscompPreLowered && dryStats.ByPass["slab_alloc_lowering"] < 1 {
			t.Fatalf("expected >= 1 slab_alloc_lowering sites on real gen/jscomp, got %+v", dryStats.ByPass)
		}
		if jscompPreLowered && dryStats.ByPass["slab_alloc_lowering"] != 0 {
			t.Fatalf("expected 0 slab_alloc_lowering sites on already-lowered gen/jscomp, got %+v", dryStats.ByPass)
		}
	}
}

func TestWave19B_PooledSlabLinkedDirectedGraphAndCfgRelease(t *testing.T) {
	tmpDir, _ := setupWave18SynthWorkspace(t)
	allPasses := passes.AllPasses()

	graphDir := filepath.Join(tmpDir, "graph")
	jscompDir := filepath.Join(tmpDir, "jscompcfg")
	for _, d := range []string{graphDir, jscompDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	graphSrc := `package graph

import "example.com/w18slabtest/j2gort"

type Graph[N any, E any] interface {
	Implements_graph_Graph()
	IsJavaNil_() bool
}

type LinkedDirectedGraph_LinkedDiGraphNode[N any, E any] struct {
	val N
}

func (n *LinkedDirectedGraph_LinkedDiGraphNode[N, E]) initFields_Node() { _ = n }

func NewLinkedDirectedGraph_LinkedDiGraphNode[N any, E any](v N) *LinkedDirectedGraph_LinkedDiGraphNode[N, E] {
	this_ := &LinkedDirectedGraph_LinkedDiGraphNode[N, E]{}
	InitLinkedDirectedGraph_LinkedDiGraphNode(this_, v)
	return this_
}

func InitLinkedDirectedGraph_LinkedDiGraphNode[N any, E any](this_ *LinkedDirectedGraph_LinkedDiGraphNode[N, E], v N) {
	this_.initFields_Node()
	this_.val = v
}

type LinkedDirectedGraph_AnnotatedLinkedDiGraphNode[N any, E any] struct {
	val N
	ann any
}

func (n *LinkedDirectedGraph_AnnotatedLinkedDiGraphNode[N, E]) initFields_AnnNode() { _ = n }

func NewLinkedDirectedGraph_AnnotatedLinkedDiGraphNode[N any, E any](v N) *LinkedDirectedGraph_AnnotatedLinkedDiGraphNode[N, E] {
	this_ := &LinkedDirectedGraph_AnnotatedLinkedDiGraphNode[N, E]{}
	InitLinkedDirectedGraph_AnnotatedLinkedDiGraphNode(this_, v)
	return this_
}

func InitLinkedDirectedGraph_AnnotatedLinkedDiGraphNode[N any, E any](this_ *LinkedDirectedGraph_AnnotatedLinkedDiGraphNode[N, E], v N) {
	this_.initFields_AnnNode()
	this_.val = v
}

type LinkedDirectedGraph_LinkedDiGraphEdge[N any, E any] struct {
	val E
}

func (e *LinkedDirectedGraph_LinkedDiGraphEdge[N, E]) initFields_Edge() { _ = e }

func NewLinkedDirectedGraph_LinkedDiGraphEdge[N any, E any](v E) *LinkedDirectedGraph_LinkedDiGraphEdge[N, E] {
	this_ := &LinkedDirectedGraph_LinkedDiGraphEdge[N, E]{}
	InitLinkedDirectedGraph_LinkedDiGraphEdge(this_, v)
	return this_
}

func InitLinkedDirectedGraph_LinkedDiGraphEdge[N any, E any](this_ *LinkedDirectedGraph_LinkedDiGraphEdge[N, E], v E) {
	this_.initFields_Edge()
	this_.val = v
}

type LinkedDirectedGraph_AnnotatedLinkedDiGraphEdge[N any, E any] struct {
	val E
	ann any
}

func (e *LinkedDirectedGraph_AnnotatedLinkedDiGraphEdge[N, E]) initFields_AnnEdge() { _ = e }

func NewLinkedDirectedGraph_AnnotatedLinkedDiGraphEdge[N any, E any](v E) *LinkedDirectedGraph_AnnotatedLinkedDiGraphEdge[N, E] {
	this_ := &LinkedDirectedGraph_AnnotatedLinkedDiGraphEdge[N, E]{}
	InitLinkedDirectedGraph_AnnotatedLinkedDiGraphEdge(this_, v)
	return this_
}

func InitLinkedDirectedGraph_AnnotatedLinkedDiGraphEdge[N any, E any](this_ *LinkedDirectedGraph_AnnotatedLinkedDiGraphEdge[N, E], v E) {
	this_.initFields_AnnEdge()
	this_.val = v
}

type LinkedDirectedGraphBase[N any, E any] struct {
	nodes []*LinkedDirectedGraph_LinkedDiGraphNode[N, E]
	edges []*LinkedDirectedGraph_LinkedDiGraphEdge[N, E]
}

func (*LinkedDirectedGraphBase[N, E]) Implements_graph_Graph() {}
func (this_ *LinkedDirectedGraphBase[N, E]) IsJavaNil_() bool  { return this_ == nil }

func (this_ *LinkedDirectedGraphBase[N, E]) CreateNode(v N, annotated bool) any {
	if annotated {
		return NewLinkedDirectedGraph_AnnotatedLinkedDiGraphNode[N, E](v)
	}
	n := NewLinkedDirectedGraph_LinkedDiGraphNode[N, E](v)
	this_.nodes = append(this_.nodes, n)
	return n
}

func (this_ *LinkedDirectedGraphBase[N, E]) Connect(v E, annotated bool) any {
	if annotated {
		return NewLinkedDirectedGraph_AnnotatedLinkedDiGraphEdge[N, E](v)
	}
	e := NewLinkedDirectedGraph_LinkedDiGraphEdge[N, E](v)
	this_.edges = append(this_.edges, e)
	return e
}

var _ = j2gort.NullString
`

	jscompSrc := `package jscompcfg

import (
	"example.com/w18slabtest/graph"
	"example.com/w18slabtest/j2gort"
	pkg_rhino "example.com/w18slabtest/rhino"
)

type ControlFlowGraph_Branch int32

type ControlFlowGraph[N any] interface {
	IsJavaNil_() bool
	EntryVal() N
}

type ControlFlowGraphBase[N any] struct {
	graph.LinkedDirectedGraphBase[N, ControlFlowGraph_Branch]
	entry N
}

type ControlFlowGraphImpl[N any] struct {
	ControlFlowGraphBase[N]
}

func NewControlFlowGraphImpl[N any](entry N) *ControlFlowGraphImpl[N] {
	g := &ControlFlowGraphImpl[N]{}
	g.entry = entry
	g.CreateNode(entry, false)
	g.CreateNode(entry, true)
	g.Connect(1, false)
	g.Connect(2, true)
	return g
}

func (g *ControlFlowGraphImpl[N]) IsJavaNil_() bool { return g == nil }
func (g *ControlFlowGraphImpl[N]) EntryVal() N      { return g.entry }

func NodeUtil_IsValidCfgRoot(n pkg_rhino.Node) bool {
	return n != 0
}

type TypeInferencePass struct {
	panicMode bool
}

func (this_ *TypeInferencePass) computeCfg(n pkg_rhino.Node) ControlFlowGraph[pkg_rhino.Node] {
	if n == 0 {
		var nilCfg *ControlFlowGraphImpl[pkg_rhino.Node]
		return nilCfg
	}
	return NewControlFlowGraphImpl[pkg_rhino.Node](n)
}

func (this_ *TypeInferencePass) priv_TypeInferencePass_inferScope_1_Node(n pkg_rhino.Node) {
	var cfg ControlFlowGraph[pkg_rhino.Node] = this_.computeCfg(n)
	_ = cfg
	if this_.panicMode {
		panic("infer_panic")
	}
}

type NodeTraversal struct {
	root pkg_rhino.Node
}

func (t *NodeTraversal) GetScopeRoot() pkg_rhino.Node { return t.root }

type NodeTraversal_AbstractCfgCallbackBase struct {
	cfgs j2gort.Deque[ControlFlowGraph[pkg_rhino.Node]]
}

func (this_ *NodeTraversal_AbstractCfgCallbackBase) EnterScope(t *NodeTraversal) {
	if NodeUtil_IsValidCfgRoot(t.GetScopeRoot()) {
		this_.cfgs.Push(NewControlFlowGraphImpl[pkg_rhino.Node](t.GetScopeRoot()))
	}
}

func (this_ *NodeTraversal_AbstractCfgCallbackBase) ExitScope(t *NodeTraversal) {
	_ = t
	if NodeUtil_IsValidCfgRoot(t.GetScopeRoot()) {
		this_.cfgs.Pop()
	}
}

var _ = j2gort.NullString
`

	if err := os.WriteFile(filepath.Join(graphDir, "linked_directed_graph.go"), []byte(graphSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(jscompDir, "cfg_users.go"), []byte(jscompSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	stats1, err := passes.Run([]string{graphDir, jscompDir}, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 on graph+jscompcfg failed: %v", err)
	}
	if stats1.ByPass["slab_alloc_lowering"] != 6 {
		t.Fatalf("expected 6 slab_alloc_lowering sites (4 in graph + 2 in jscompcfg), got %+v", stats1.ByPass)
	}

	// Run 2 MUST be 100% idempotent (0 sites, 0 modified files)!
	stats2, err := passes.Run([]string{graphDir, jscompDir}, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 on graph+jscompcfg failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites & 0 modified files on Run 2, got %+v", stats2)
	}

	loweredGraphBytes, err := os.ReadFile(filepath.Join(graphDir, "linked_directed_graph.go"))
	if err != nil {
		t.Fatal(err)
	}
	loweredGraph := string(loweredGraphBytes)
	for _, want := range []string{
		"slab_LinkedDirectedGraph_LinkedDiGraphNode",
		"slab_LinkedDirectedGraph_AnnotatedLinkedDiGraphNode",
		"slab_LinkedDirectedGraph_LinkedDiGraphEdge",
		"slab_LinkedDirectedGraph_AnnotatedLinkedDiGraphEdge",
		"j2gort.PooledSlab[LinkedDirectedGraph_LinkedDiGraphNode[N, E]]",
		"j2gort.PooledSlab[LinkedDirectedGraph_AnnotatedLinkedDiGraphNode[N, E]]",
		"j2gort.PooledSlab[LinkedDirectedGraph_LinkedDiGraphEdge[N, E]]",
		"j2gort.PooledSlab[LinkedDirectedGraph_AnnotatedLinkedDiGraphEdge[N, E]]",
		"func (this_ *LinkedDirectedGraphBase[N, E]) ReleaseSlabs_()",
	} {
		if !strings.Contains(loweredGraph, want) {
			t.Fatalf("missing %q in lowered linked_directed_graph.go:\n%s", want, loweredGraph)
		}
	}
	if strings.Contains(loweredGraph, "j2gort.Slab[") {
		t.Fatalf("LinkedDirectedGraphBase must use PooledSlab, not Slab:\n%s", loweredGraph)
	}

	loweredJscompBytes, err := os.ReadFile(filepath.Join(jscompDir, "cfg_users.go"))
	if err != nil {
		t.Fatal(err)
	}
	loweredJscomp := string(loweredJscompBytes)
	for _, want := range []string{
		"if cfg != nil && !cfg.IsJavaNil_()",
		"defer r_.ReleaseSlabs_()",
		"poppedCfg_ := this_.cfgs.Pop()",
		"r_.ReleaseSlabs_()",
	} {
		if !strings.Contains(loweredJscomp, want) {
			t.Fatalf("missing %q in lowered cfg_users.go:\n%s", want, loweredJscomp)
		}
	}

	rtTestSrc := `package w19b_cfg_test

import (
	"testing"

	"example.com/w18slabtest/j2gort"
	"example.com/w18slabtest/jscompcfg"
	pkg_rhino "example.com/w18slabtest/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
)

func TestCfgPooledSlabReleaseRuntime(t *testing.T) {
	store := nodestore.NewNodeStore()
	prev := nodestore.SetCurrentNodeStore(store)
	defer nodestore.SetCurrentNodeStore(prev)

	root := pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_SCRIPT())
	j2gort.PooledSlabReleaseCalls.Store(0)

	tip := &jscompcfg.TypeInferencePass{}
	// Normal inferScope releases 4 PooledSlab fields on LinkedDirectedGraphBase.
	tip.ExportedInfer(root)
	if got := j2gort.PooledSlabReleaseCalls.Load(); got != 4 {
		t.Fatalf("after normal inferScope: PooledSlabReleaseCalls = %d, want 4", got)
	}

	// Typed-nil CFG (root == 0) must not panic or call ReleaseSlabs_().
	tip.ExportedInfer(0)
	if got := j2gort.PooledSlabReleaseCalls.Load(); got != 4 {
		t.Fatalf("after nil-CFG inferScope: PooledSlabReleaseCalls = %d, want 4", got)
	}

	// Panic inside inferScope still fires defer r_.ReleaseSlabs_()!
	tipPanic := &jscompcfg.TypeInferencePass{}
	tipPanic.SetPanic(true)
	func() {
		defer func() { _ = recover() }()
		tipPanic.ExportedInfer(root)
	}()
	if got := j2gort.PooledSlabReleaseCalls.Load(); got != 8 {
		t.Fatalf("after panicking inferScope: PooledSlabReleaseCalls = %d, want 8", got)
	}

	// AbstractCfgCallbackBase EnterScope + ExitScope pops CFG and releases all 4 PooledSlabs.
	cb := &jscompcfg.NodeTraversal_AbstractCfgCallbackBase{}
	nt := jscompcfg.NewNodeTraversal(root)
	cb.EnterScope(nt)
	cb.ExitScope(nt)
	if got := j2gort.PooledSlabReleaseCalls.Load(); got != 12 {
		t.Fatalf("after AbstractCfgCallbackBase.ExitScope: PooledSlabReleaseCalls = %d, want 12", got)
	}
}
`
	helperExportSrc := `package jscompcfg

import pkg_rhino "example.com/w18slabtest/rhino"

func (this_ *TypeInferencePass) SetPanic(v bool)             { this_.panicMode = v }
func (this_ *TypeInferencePass) ExportedInfer(n pkg_rhino.Node) { this_.priv_TypeInferencePass_inferScope_1_Node(n) }
func NewNodeTraversal(r pkg_rhino.Node) *NodeTraversal       { return &NodeTraversal{root: r} }
`
	if err := os.WriteFile(filepath.Join(jscompDir, "export_test_bridge.go"), []byte(helperExportSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "w19b_cfg_test.go"), []byte(rtTestSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "test", "-v", "-run", "TestCfgPooledSlabReleaseRuntime", ".")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test TestCfgPooledSlabReleaseRuntime failed: %v\n%s", err, string(out))
	}
}

func TestWave30B_SlabLeaks_VarInfo_GlobalNamespaceRef_And_PoolPurgeStripping(t *testing.T) {
	tmpDir := t.TempDir()
	mustWrite := func(rel, content string) {
		t.Helper()
		full := filepath.Join(tmpDir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	mustWrite("go.mod", "module example.com/w30bslabtest\n\ngo 1.24\n")
	mustWrite("gen/rhino/node.go", `package rhino

type Node uint32

func (n Node) IsJavaNil_() bool { return n == 0 }
`)
	mustWrite("gen/j2gort/j2gort.go", `package j2gort

import "sync/atomic"

var (
	SlabAllocCnt   atomic.Int64
	SlabReleaseCnt atomic.Int64
	PoolPurgeCnt   atomic.Int64
)

type Slab[T any] struct{}
func (s *Slab[T]) Alloc() *T { return new(T) }

type PooledSlab[T any] struct {
	items []*T
}

func (s *PooledSlab[T]) Alloc() *T {
	SlabAllocCnt.Add(1)
	v := new(T)
	s.items = append(s.items, v)
	return v
}

func (s *PooledSlab[T]) Release() {
	SlabReleaseCnt.Add(1)
	for _, it := range s.items {
		var zero T
		*it = zero
	}
	s.items = s.items[:0]
}

func (s *PooledSlab[T]) ReleaseWith(fn func(*T)) {
	for _, it := range s.items {
		fn(it)
		var zero T
		*it = zero
	}
	s.items = s.items[:0]
	SlabReleaseCnt.Add(1)
}

type OrderedMap[K comparable, V any] struct {
	m map[K]V
}

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V]    { return &OrderedMap[K, V]{m: make(map[K]V)} }
func AcquireOrderedMap[K comparable, V any]() OrderedMap[K, V] { return OrderedMap[K, V]{m: make(map[K]V)} }
func AcquireOrderedMapPtr[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{m: make(map[K]V)}
}
func (o *OrderedMap[K, V]) Put(k K, v V) {
	if o.m == nil {
		o.m = make(map[K]V)
	}
	o.m[k] = v
}
func (o *OrderedMap[K, V]) Get(k K) V    { return o.m[k] }
func (o *OrderedMap[K, V]) Clear()       { clear(o.m) }
func (o *OrderedMap[K, V]) IsEmpty() bool { return len(o.m) == 0 }
func (o *OrderedMap[K, V]) ReleaseToPool() {
	clear(o.m)
}

type MapCursor[V any] struct {
	vals []V
	idx  int
}

func (c *MapCursor[V]) HasNext() bool { return c.idx < len(c.vals) }
func (c *MapCursor[V]) Next() V {
	v := c.vals[c.idx]
	c.idx++
	return v
}
func (o *OrderedMap[K, V]) IterValues() *MapCursor[V] {
	vs := make([]V, 0, len(o.m))
	for _, v := range o.m {
		vs = append(vs, v)
	}
	return &MapCursor[V]{vals: vs}
}

func ReleaseOrderedMap[K comparable, V any](o *OrderedMap[K, V]) {
	if o != nil {
		clear(o.m)
	}
}

func AcquirePooledSlice[T any](capHint int) []T {
	if capHint < 0 {
		capHint = 0
	}
	return make([]T, 0, capHint)
}

func ReleasePooledSlice[T any](s []T) {
	clear(s)
}

func PurgeAllRuntimePools()         { PoolPurgeCnt.Add(1) }
func TrimStaleSlabPools(keep int)   { _ = keep; PoolPurgeCnt.Add(1) }
func CastPtr[T any](v any) *T       { p, _ := v.(*T); return p }
func CheckNotNullPtr[T any](p *T) *T { return p }
`)

	mustWrite("gen/jscomp/core.go", `package jscomp

import (
	"sync/atomic"

	"example.com/w30bslabtest/gen/j2gort"
	pkg_rhino "example.com/w30bslabtest/gen/rhino"
)

var ScopeCreatorReleaseCnt atomic.Int64

type AbstractCompiler interface {
	IsJavaNil_() bool
	GetRoot() pkg_rhino.Node
}

type Compiler struct{}

func (c *Compiler) IsJavaNil_() bool        { return c == nil }
func (c *Compiler) GetRoot() pkg_rhino.Node { return pkg_rhino.Node(1) }

func (this_ *Compiler) ParseInputs() pkg_rhino.Node {
	defer j2gort.PurgeAllRuntimePools()
	return pkg_rhino.Node(1)
}

type TypeInferencePass struct{}

func (this_ *TypeInferencePass) InferAllScopes(root pkg_rhino.Node) {
	_ = root
	j2gort.PurgeAllRuntimePools()
	defer j2gort.TrimStaleSlabPools(0)
}

type VariableReferenceCheck struct{}

func (this_ *VariableReferenceCheck) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	defer j2gort.PurgeAllRuntimePools()
}

type CompilerPass interface {
	Implements_jscomp_CompilerPass()
	Process_2(externs pkg_rhino.Node, root pkg_rhino.Node)
}

type ScopeCreator interface {
	IsJavaNil_() bool
	CreateScope_2_Node_Scope(n pkg_rhino.Node, parent *Scope) *Scope
}

type Var struct {
	name       string
	parentNode pkg_rhino.Node
}

func (v *Var) GetName() string               { return v.name }
func (v *Var) GetParentNode() pkg_rhino.Node { return v.parentNode }
func (v *Var) GetNameNode() pkg_rhino.Node   { return v.parentNode }

type Scope struct {
	root pkg_rhino.Node
	vars map[string]*Var
}

func (s *Scope) GetVar(name string) any {
	if s == nil || s.vars == nil {
		return (*Var)(nil)
	}
	return s.vars[name]
}

func (s *Scope) HasSlot(name string) bool {
	return s != nil && s.vars != nil && s.vars[name] != nil
}

type SyntacticScopeCreator struct {
	compiler AbstractCompiler
	released bool
}

func (s *SyntacticScopeCreator) IsJavaNil_() bool { return s == nil }

func NewSyntacticScopeCreator_1(c AbstractCompiler) *SyntacticScopeCreator {
	s := &SyntacticScopeCreator{}
	InitSyntacticScopeCreator_1(s, c)
	return s
}

func InitSyntacticScopeCreator_1(s *SyntacticScopeCreator, c AbstractCompiler) {
	s.compiler = c
}

func (s *SyntacticScopeCreator) CreateScope_2_Node_Scope(n pkg_rhino.Node, parent *Scope) *Scope {
	_ = parent
	return &Scope{
		root: n,
		vars: map[string]*Var{
			"fn": {name: "fn", parentNode: n},
			"x":  {name: "x", parentNode: n},
		},
	}
}

func (s *SyntacticScopeCreator) ReleaseSlabs_() {
	if s == nil {
		return
	}
	s.released = true
	ScopeCreatorReleaseCnt.Add(1)
}

type MemoizedScopeCreator struct {
	delegate *SyntacticScopeCreator
}

func (m *MemoizedScopeCreator) IsJavaNil_() bool { return m == nil }
func (m *MemoizedScopeCreator) CreateScope_2_Node_Scope(n pkg_rhino.Node, p *Scope) *Scope {
	return m.delegate.CreateScope_2_Node_Scope(n, p)
}
func (m *MemoizedScopeCreator) ReleaseSlabs_() {
	if m != nil && m.delegate != nil {
		m.delegate.ReleaseSlabs_()
	}
}

type NodeTraversal_Callback interface {
	Implements_jscomp_NodeTraversal_Callback()
	IsJavaNil_() bool
	ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool
	Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node)
}

type NodeTraversal_AbstractPostOrderCallbackBase struct{}

func (b *NodeTraversal_AbstractPostOrderCallbackBase) IsJavaNil_() bool { return b == nil }
func (*NodeTraversal_AbstractPostOrderCallbackBase) Implements_jscomp_NodeTraversal_Callback() {}
func (*NodeTraversal_AbstractPostOrderCallbackBase) ShouldTraverse_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node) bool {
	return true
}

type NodeTraversal struct {
	compiler     AbstractCompiler
	callback     NodeTraversal_Callback
	scopeCreator ScopeCreator
	curScope     *Scope
}

func (t *NodeTraversal) GetScope() *Scope {
	if t.curScope == nil {
		if t.scopeCreator == nil {
			t.scopeCreator = NewSyntacticScopeCreator_1(t.compiler)
		}
		t.curScope = t.scopeCreator.CreateScope_2_Node_Scope(pkg_rhino.Node(1), nil)
	}
	return t.curScope
}

type NodeTraversal_Builder struct {
	compiler     AbstractCompiler
	callback     NodeTraversal_Callback
	scopeCreator ScopeCreator
}

func NodeTraversal_S_Builder() *NodeTraversal_Builder { return &NodeTraversal_Builder{} }
func (b *NodeTraversal_Builder) SetCompiler(c AbstractCompiler) *NodeTraversal_Builder {
	b.compiler = c
	return b
}
func (b *NodeTraversal_Builder) SetCallback_1_NodeTraversal_Callback(cb NodeTraversal_Callback) *NodeTraversal_Builder {
	b.callback = cb
	return b
}
func (b *NodeTraversal_Builder) SetScopeCreator(sc ScopeCreator) *NodeTraversal_Builder {
	b.scopeCreator = sc
	return b
}
func (b *NodeTraversal_Builder) Traverse(root pkg_rhino.Node) {
	nt := &NodeTraversal{compiler: b.compiler, callback: b.callback, scopeCreator: b.scopeCreator}
	if nt.callback.ShouldTraverse_3(nt, root, 0) {
		nt.callback.Visit_3(nt, root, 0)
	}
}
func (b *NodeTraversal_Builder) TraverseRoots(externs, root pkg_rhino.Node) {
	b.Traverse(externs)
	b.Traverse(root)
}

func NodeTraversal_Traverse(c AbstractCompiler, root pkg_rhino.Node, cb NodeTraversal_Callback) {
	sc := NewSyntacticScopeCreator_1(c)
	if _, ok := any(cb).(interface{ CanReleaseSyntacticScopes_() }); ok {
		defer sc.ReleaseSlabs_()
	}
	NodeTraversal_S_Builder().SetCompiler(c).SetCallback_1_NodeTraversal_Callback(cb).SetScopeCreator(sc).Traverse(root)
}

func NodeTraversal_TraverseRoots(c AbstractCompiler, cb NodeTraversal_Callback, externs, root pkg_rhino.Node) {
	sc := NewSyntacticScopeCreator_1(c)
	if _, ok := any(cb).(interface{ CanReleaseSyntacticScopes_() }); ok {
		defer sc.ReleaseSlabs_()
	}
	NodeTraversal_S_Builder().SetCompiler(c).SetCallback_1_NodeTraversal_Callback(cb).SetScopeCreator(sc).TraverseRoots(externs, root)
}
`)

	mustWrite("gen/jscomp/scope_passes.go", `package jscomp

import (
	"example.com/w30bslabtest/gen/j2gort"
	pkg_rhino "example.com/w30bslabtest/gen/rhino"
)

// 1. CheckSuper + ConvertTypesToColors callbacks (scope-free -> MUST receive CanReleaseSyntacticScopes_())
type CheckSuper struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler     AbstractCompiler
	contextStack []any
}

func (this_ *CheckSuper) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
}

type ConvertTypesToColors_RemoveTypesBase struct {
	NodeTraversal_AbstractPostOrderCallbackBase
}

func (this_ *ConvertTypesToColors_RemoveTypesBase) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
}

type ConvertTypesToColors_RemoveTypesAndApplyColors struct {
	ConvertTypesToColors_RemoveTypesBase
	count int
}

func (this_ *ConvertTypesToColors_RemoveTypesAndApplyColors) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	this_.ConvertTypesToColors_RemoveTypesBase.Visit_3(t, n, parent)
	this_.count++
}

// Negative guard: callback calling GetScope() MUST NOT get CanReleaseSyntacticScopes_()!
type ScopeReadingCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
}

func (this_ *ScopeReadingCallback) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = n
	_ = parent
	_ = t.GetScope()
}

// 2. CheckMissingRequires.Process_2
type NodeTraversal_AbstractModuleCallbackBase struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler
}

func (b *NodeTraversal_AbstractModuleCallbackBase) GetNodeTraversal_AbstractModuleCallbackBase() *NodeTraversal_AbstractModuleCallbackBase {
	return b
}

type CheckMissingRequires struct {
	NodeTraversal_AbstractModuleCallbackBase
	templateParamNames *j2gort.OrderedMap[string, bool]
}

func (*CheckMissingRequires) Implements_jscomp_CompilerPass() {}
func (this_ *CheckMissingRequires) Visit_3(t *NodeTraversal, n pkg_rhino.Node, p pkg_rhino.Node) {
	_ = n
	_ = p
	v := j2gort.CastPtr[Var](t.GetScope().GetVar("x"))
	this_.templateParamNames.Put(v.GetName(), true)
}
func (this_ *CheckMissingRequires) Visit_4(t *NodeTraversal, n pkg_rhino.Node, p pkg_rhino.Node, m any) {
	_ = m
	this_.Visit_3(t, n, p)
}
func (this_ *CheckMissingRequires) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	NodeTraversal_Traverse(this_.GetNodeTraversal_AbstractModuleCallbackBase().compiler, root, this_)
}

// 3. ScopedAliases.Process_2
type ScopedAliases struct {
	compiler AbstractCompiler
}

func (*ScopedAliases) Implements_jscomp_CompilerPass() {}

type ScopedAliases_Traversal struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *ScopedAliases
	seen   []*Var
}

func NewScopedAliases_Traversal(outer_ *ScopedAliases) *ScopedAliases_Traversal {
	return &ScopedAliases_Traversal{outer_: outer_}
}
func (this_ *ScopedAliases_Traversal) Visit_3(t *NodeTraversal, n pkg_rhino.Node, p pkg_rhino.Node) {
	_ = n
	_ = p
	this_.seen = append(this_.seen, j2gort.CastPtr[Var](t.GetScope().GetVar("x")))
}
func (this_ *ScopedAliases) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	var traversal *ScopedAliases_Traversal = NewScopedAliases_Traversal(this_)
	NodeTraversal_Traverse(this_.compiler, root, traversal)
	_ = len(traversal.seen)
}

// 4. RewriteCallerCodeLocation.Process_2
type RewriteCallerCodeLocation struct {
	compiler                    AbstractCompiler
	callerLocationFunctionNames *j2gort.OrderedMap[string, *Var]
}

func (*RewriteCallerCodeLocation) Implements_jscomp_CompilerPass() {}

type RewriteCallerCodeLocation_FindCallerLocationFunctions struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *RewriteCallerCodeLocation
}

func NewRewriteCallerCodeLocation_FindCallerLocationFunctions(o *RewriteCallerCodeLocation) *RewriteCallerCodeLocation_FindCallerLocationFunctions {
	return &RewriteCallerCodeLocation_FindCallerLocationFunctions{outer_: o}
}
func (this_ *RewriteCallerCodeLocation_FindCallerLocationFunctions) Visit_3(t *NodeTraversal, n pkg_rhino.Node, p pkg_rhino.Node) {
	_ = n
	_ = p
	this_.outer_.callerLocationFunctionNames.Put("fn", j2gort.CastPtr[Var](t.GetScope().GetVar("fn")))
}

type RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *RewriteCallerCodeLocation
}

func NewRewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls(o *RewriteCallerCodeLocation) *RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls {
	return &RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls{outer_: o}
}
func (this_ *RewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls) Visit_3(t *NodeTraversal, n pkg_rhino.Node, p pkg_rhino.Node) {
	_ = n
	_ = p
	v := j2gort.CastPtr[Var](t.GetScope().GetVar("fn"))
	_ = v.GetNameNode() == this_.outer_.callerLocationFunctionNames.Get("fn").GetNameNode()
}

func (this_ *RewriteCallerCodeLocation) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, this_.compiler.GetRoot(), NewRewriteCallerCodeLocation_FindCallerLocationFunctions(this_))
	if !(this_.callerLocationFunctionNames.IsEmpty()) {
		NodeTraversal_Traverse(this_.compiler, this_.compiler.GetRoot(), NewRewriteCallerCodeLocation_RewriteCallerLocationFunctionCalls(this_))
	}
}

// 5. Es6NormalizeClasses.Process_2
type Es6NormalizeClasses struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler   AbstractCompiler
	classStack *j2gort.OrderedMap[string, *Var]
}

func (*Es6NormalizeClasses) Implements_jscomp_CompilerPass() {}
func (this_ *Es6NormalizeClasses) Visit_3(t *NodeTraversal, n pkg_rhino.Node, p pkg_rhino.Node) {
	_ = n
	_ = p
	this_.classStack.Put("C", j2gort.CastPtr[Var](t.GetScope().GetVar("x")))
}
func (this_ *Es6NormalizeClasses) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	NodeTraversal_TraverseRoots(this_.compiler, this_, externs, root)
}

// 6. Inline CreateScope_2_Node_Scope in InferConsts & IsolatePolyfills
type InferConsts struct {
	compiler AbstractCompiler
	found    bool
}

func (this_ *InferConsts) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = root
	var globalExternsScope *Scope = NewSyntacticScopeCreator_1(this_.compiler).CreateScope_2_Node_Scope(externs, nil)
	_ = globalExternsScope
	var v *Var = j2gort.CheckNotNullPtr(j2gort.CastPtr[Var](globalExternsScope.GetVar("x")))
	_ = v
	if v != nil && v.GetParentNode() != 0 {
		this_.found = globalExternsScope.HasSlot("x")
	}
}
`)

	mustWrite("gen/jscomp/remove_unused_code.go", `package jscomp

import (
	"example.com/w30bslabtest/gen/j2gort"
	pkg_rhino "example.com/w30bslabtest/gen/rhino"
)

type RemoveUnusedCode struct {
	canonicalUnremovableVarInfo *RemoveUnusedCode_VarInfo
	scopeCreator                ScopeCreator
	allFunctionParamScopes      []*Scope
	varInfoKeys                 []*Var
	varInfoVals                 []*RemoveUnusedCode_VarInfo
	varInfoTable                []int32
	varInfoSize                 int32
}

func (this_ *RemoveUnusedCode) initFields_RemoveUnusedCode() {
	this_.varInfoKeys = make([]*Var, int32(16))
	this_.varInfoVals = make([]*RemoveUnusedCode_VarInfo, int32(16))
	this_.varInfoTable = make([]int32, int32(32))
}

func NewRemoveUnusedCode(c AbstractCompiler) *RemoveUnusedCode {
	r := &RemoveUnusedCode{}
	InitRemoveUnusedCode(r, c)
	return r
}

func InitRemoveUnusedCode(this_ *RemoveUnusedCode, c AbstractCompiler) {
	this_.initFields_RemoveUnusedCode()
	this_.canonicalUnremovableVarInfo = NewRemoveUnusedCode_VarInfo(this_)
	this_.scopeCreator = NewSyntacticScopeCreator_1(c)
}

func (this_ *RemoveUnusedCode) GetVarInfo(v *Var) *RemoveUnusedCode_VarInfo {
	if v == nil {
		return this_.canonicalUnremovableVarInfo
	}
	var val_11 RemoveUnusedCode_VarInfo
	val_11.outer_ = this_
	InitRemoveUnusedCode_VarInfo_1(&val_11, v.GetName())
	var varInfo *RemoveUnusedCode_VarInfo = &val_11
	this_.varInfoKeys[0] = v
	this_.varInfoVals[0] = varInfo
	this_.varInfoSize = 1
	return varInfo
}

func (this_ *RemoveUnusedCode) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	sc := this_.scopeCreator.CreateScope_2_Node_Scope(root, nil)
	v := j2gort.CastPtr[Var](sc.GetVar("x"))
	_ = this_.GetVarInfo(v)
}

type RemoveUnusedCode_VarInfo struct {
	outer_               *RemoveUnusedCode
	varName              string
	canonicalUnremovable bool
	isEntirelyRemovable  bool
}

func (this_ *RemoveUnusedCode_VarInfo) initFields_RemoveUnusedCode_VarInfo() {
	_ = this_
}

func NewRemoveUnusedCode_VarInfo(outer_ *RemoveUnusedCode) *RemoveUnusedCode_VarInfo {
	this_ := &RemoveUnusedCode_VarInfo{}
	this_.outer_ = outer_
	InitRemoveUnusedCode_VarInfo(this_)
	return this_
}

func InitRemoveUnusedCode_VarInfo(this_ *RemoveUnusedCode_VarInfo) {
	this_.initFields_RemoveUnusedCode_VarInfo()
	this_.varName = ""
	this_.canonicalUnremovable = true
}

func NewRemoveUnusedCode_VarInfo_1(outer_ *RemoveUnusedCode, varName string) *RemoveUnusedCode_VarInfo {
	this_ := &RemoveUnusedCode_VarInfo{}
	this_.outer_ = outer_
	InitRemoveUnusedCode_VarInfo_1(this_, varName)
	return this_
}

func InitRemoveUnusedCode_VarInfo_1(this_ *RemoveUnusedCode_VarInfo, varName string) {
	this_.initFields_RemoveUnusedCode_VarInfo()
	this_.varName = varName
	this_.isEntirelyRemovable = true
}
`)

	mustWrite("gen/jscomp/global_namespace.go", `package jscomp

import (
	"example.com/w30bslabtest/gen/j2gort"
	pkg_rhino "example.com/w30bslabtest/gen/rhino"
)

type JSChunk struct{}
type GlobalNamespace_Ref_Type int32

type GlobalNamespace struct {
	compiler    AbstractCompiler
	globalNames []*GlobalNamespace_Name
	nameMap     j2gort.OrderedMap[string, *GlobalNamespace_Name]
}

func NewGlobalNamespace_3_AbstractCompiler_Node_Node(c AbstractCompiler, externs, root pkg_rhino.Node) *GlobalNamespace {
	gn := &GlobalNamespace{compiler: c, nameMap: j2gort.AcquireOrderedMap[string, *GlobalNamespace_Name]()}
	gn.priv_GlobalNamespace_process_1(externs, root)
	return gn
}

func (this_ *GlobalNamespace) priv_GlobalNamespace_process_1(externs, root pkg_rhino.Node) {
	var builder *NodeTraversal_Builder = NodeTraversal_S_Builder().SetCompiler(this_.compiler)
	_ = builder
	nm := NewGlobalNamespace_Name(this_, "foo")
	this_.nameMap.Put("foo", nm)
	sc := NewSyntacticScopeCreator_1(this_.compiler).CreateScope_2_Node_Scope(root, nil)
	_ = externs
	nm.AddRef(nil, sc, root, 1)
}

type GlobalNamespace_Name struct {
	outer_         *GlobalNamespace
	baseName       string
	ref0           *GlobalNamespace_Ref
	refsForNodeMap *j2gort.OrderedMap[pkg_rhino.Node, *GlobalNamespace_Ref]
}

func NewGlobalNamespace_Name(outer_ *GlobalNamespace, baseName string) *GlobalNamespace_Name {
	return &GlobalNamespace_Name{outer_: outer_, baseName: baseName}
}

func (this_ *GlobalNamespace_Name) priv_Name_createNewRef_3077(chunk *JSChunk, scope *Scope, node pkg_rhino.Node, type_ GlobalNamespace_Ref_Type) *GlobalNamespace_Ref {
	return NewGlobalNamespace_Ref(chunk, j2gort.CheckNotNullPtr(scope), node, type_)
}

func (this_ *GlobalNamespace_Name) AddRef(chunk *JSChunk, scope *Scope, node pkg_rhino.Node, type_ GlobalNamespace_Ref_Type) {
	this_.ref0 = this_.priv_Name_createNewRef_3077(chunk, scope, node, type_)
}

func (this_ *GlobalNamespace_Name) AddSingleRefForTesting(node pkg_rhino.Node, type_ GlobalNamespace_Ref_Type) *GlobalNamespace_Ref {
	var val_13 GlobalNamespace_Ref
	InitGlobalNamespace_Ref(&val_13, nil, nil, node, type_)
	var ref *GlobalNamespace_Ref = &val_13
	this_.ref0 = ref
	return ref
}

type GlobalNamespace_Ref struct {
	chunk *JSChunk
	scope *Scope
	node  pkg_rhino.Node
	type_ GlobalNamespace_Ref_Type
}

func (this_ *GlobalNamespace_Ref) initFields_GlobalNamespace_Ref() { _ = this_ }

func NewGlobalNamespace_Ref(chunk *JSChunk, scope *Scope, node pkg_rhino.Node, type_ GlobalNamespace_Ref_Type) *GlobalNamespace_Ref {
	this_ := &GlobalNamespace_Ref{}
	InitGlobalNamespace_Ref(this_, chunk, scope, node, type_)
	return this_
}

func InitGlobalNamespace_Ref(this_ *GlobalNamespace_Ref, chunk *JSChunk, scope *Scope, node pkg_rhino.Node, type_ GlobalNamespace_Ref_Type) {
	this_.initFields_GlobalNamespace_Ref()
	this_.chunk = chunk
	this_.scope = scope
	this_.node = node
	this_.type_ = type_
}

type Es6ConvertSuperConstructorCalls struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	globalNamespace *GlobalNamespace
}

func (this_ *Es6ConvertSuperConstructorCalls) SetGlobalNamespace(gn *GlobalNamespace) {
	this_.globalNamespace = gn
}
func (this_ *Es6ConvertSuperConstructorCalls) Visit_3(t *NodeTraversal, n pkg_rhino.Node, p pkg_rhino.Node) {
	_ = t
	_ = n
	_ = p
}

type Es6RewriteClass struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler                     AbstractCompiler
	convertSuperConstructorCalls *Es6ConvertSuperConstructorCalls
}

func (*Es6RewriteClass) Implements_jscomp_CompilerPass() {}
func (this_ *Es6RewriteClass) Visit_3(t *NodeTraversal, n pkg_rhino.Node, p pkg_rhino.Node) {
	_ = t
	_ = n
	_ = p
}
func (this_ *Es6RewriteClass) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	this_.convertSuperConstructorCalls.SetGlobalNamespace(NewGlobalNamespace_3_AbstractCompiler_Node_Node(this_.compiler, externs, root))
	NodeTraversal_Traverse(this_.compiler, root, this_.convertSuperConstructorCalls)
}
`)

	jscompDir := filepath.Join(tmpDir, "gen", "jscomp")
	slabPass := []passes.Pass{passes.NewSlabAllocLoweringPass()}

	// Run 1
	stats1, err := passes.RunOnPackageDir(jscompDir, slabPass, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if stats1.TotalSites < 10 {
		t.Fatalf("expected >= 10 slab_alloc_lowering sites on Run 1, got %+v", stats1)
	}

	// Run 2: 100% idempotent
	stats2, err := passes.RunOnPackageDir(jscompDir, slabPass, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("expected Run 2 idempotency (0 sites, 0 files modified), got %+v", stats2)
	}

	readFile := func(name string) string {
		t.Helper()
		b, err := os.ReadFile(filepath.Join(jscompDir, name))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}

	coreCode := readFile("core.go")
	if strings.Contains(coreCode, "TrimStaleSlabPools(") {
		t.Fatalf("expected core.go TrimStaleSlabPools calls to be stripped:\n%s", coreCode)
	}
	if got := strings.Count(coreCode, "PurgeAllRuntimePools("); got != 1 {
		t.Fatalf("expected 1 phase-boundary PurgeAllRuntimePools call (InferAllScopes entry only, with ParseInputs and VariableReferenceCheck stripped), got %d:\n%s", got, coreCode)
	}

	normWS := func(s string) string { return strings.Join(strings.Fields(s), " ") }
	spCode := readFile("scope_passes.go")
	for _, want := range []string{
		"func (*CheckSuper) CanReleaseSyntacticScopes_() { }",
		"func (*ConvertTypesToColors_RemoveTypesBase) CanReleaseSyntacticScopes_() { }",
		"func (*ConvertTypesToColors_RemoveTypesAndApplyColors) CanReleaseSyntacticScopes_() { }",
		"this_.templateParamNames.Clear()",
		"this_.callerLocationFunctionNames.Clear()",
		"this_.classStack.Clear()",
		"defer scInline_1.ReleaseSlabs_()",
	} {
		if !strings.Contains(normWS(spCode), want) {
			t.Fatalf("expected scope_passes.go to contain %q:\n%s", want, spCode)
		}
	}
	if strings.Contains(spCode, "func (*ScopeReadingCallback) CanReleaseSyntacticScopes_()") {
		t.Fatalf("ScopeReadingCallback MUST NOT receive CanReleaseSyntacticScopes_():\n%s", spCode)
	}

	rucCode := readFile("remove_unused_code.go")
	for _, want := range []string{
		"slab_RemoveUnusedCode_VarInfo j2gort.PooledSlab[RemoveUnusedCode_VarInfo]",
		"val_11 := this_.slab_RemoveUnusedCode_VarInfo.Alloc()",
		"this_.canonicalUnremovableVarInfo = NewRemoveUnusedCode_VarInfo(this_)",
		"this_.slab_RemoveUnusedCode_VarInfo.Release()",
	} {
		if !strings.Contains(rucCode, want) {
			t.Fatalf("expected remove_unused_code.go to contain %q:\n%s", want, rucCode)
		}
	}
	valsNilIdx := strings.Index(rucCode, "this_.varInfoVals = nil")
	viRelIdx := strings.Index(rucCode, "this_.slab_RemoveUnusedCode_VarInfo.Release()")
	if valsNilIdx < 0 || viRelIdx <= valsNilIdx {
		t.Fatalf("expected this_.slab_RemoveUnusedCode_VarInfo.Release() AFTER this_.varInfoVals = nil:\n%s", rucCode)
	}

	gnCode := readFile("global_namespace.go")
	for _, want := range []string{
		"slab_GlobalNamespace_Ref j2gort.PooledSlab[GlobalNamespace_Ref]",
		"func allocPooled_GlobalNamespace_Ref(outer_ *GlobalNamespace,",
		"allocPooled_GlobalNamespace_Ref(this_.outer_,",
		"this_.slab_GlobalNamespace_Ref.Release()",
		"this_.convertSuperConstructorCalls.SetGlobalNamespace(nil)",
		"gn_.ReleaseSlabs_()",
	} {
		if !strings.Contains(gnCode, want) {
			t.Fatalf("expected global_namespace.go to contain %q:\n%s", want, gnCode)
		}
	}

	// Runtime behavioral verification of lowered package
	mustWrite("gen/jscomp/w30b_runtime_test.go", `package jscomp

import (
	"testing"

	"example.com/w30bslabtest/gen/j2gort"
	pkg_rhino "example.com/w30bslabtest/gen/rhino"
)

func TestWave30BRuntimeLifecycles(t *testing.T) {
	c := &Compiler{}
	j2gort.PoolPurgeCnt.Store(0)
	_ = c.ParseInputs()
	(&TypeInferencePass{}).InferAllScopes(1)
	(&VariableReferenceCheck{}).Process_2(1, 1)
	if got := j2gort.PoolPurgeCnt.Load(); got != 1 {
		t.Fatalf("expected 1 phase-boundary pool purge (InferAllScopes entry only, ParseInputs & VariableReferenceCheck unpurged), got %d", got)
	}

	// CheckSuper + ConvertTypesToColors_RemoveTypesAndApplyColors release scopes via NodeTraversal_Traverse
	ScopeCreatorReleaseCnt.Store(0)
	NodeTraversal_Traverse(c, 1, &CheckSuper{compiler: c})
	NodeTraversal_Traverse(c, 1, &ConvertTypesToColors_RemoveTypesAndApplyColors{})
	NodeTraversal_Traverse(c, 1, &ConvertTypesToColors_RemoveTypesBase{})
	if got := ScopeCreatorReleaseCnt.Load(); got != 3 {
		t.Fatalf("expected 3 scope creator releases from scope-free callbacks, got %d", got)
	}

	// Pass-boundary releases
	ScopeCreatorReleaseCnt.Store(0)
	cmr := &CheckMissingRequires{
		NodeTraversal_AbstractModuleCallbackBase: NodeTraversal_AbstractModuleCallbackBase{compiler: c},
		templateParamNames:                       j2gort.NewOrderedMap[string, bool](),
	}
	cmr.Process_2(1, 2)
	if !cmr.templateParamNames.IsEmpty() {
		t.Fatalf("expected CheckMissingRequires.templateParamNames cleared on exit")
	}

	sa := &ScopedAliases{compiler: c}
	sa.Process_2(1, 2)

	rccl := &RewriteCallerCodeLocation{
		compiler:                    c,
		callerLocationFunctionNames: j2gort.NewOrderedMap[string, *Var](),
	}
	rccl.Process_2(1, 2)
	if !rccl.callerLocationFunctionNames.IsEmpty() {
		t.Fatalf("expected RewriteCallerCodeLocation.callerLocationFunctionNames cleared on exit")
	}

	enc := &Es6NormalizeClasses{compiler: c, classStack: j2gort.NewOrderedMap[string, *Var]()}
	enc.Process_2(1, 2)
	if !enc.classStack.IsEmpty() {
		t.Fatalf("expected Es6NormalizeClasses.classStack cleared on exit")
	}

	ic := &InferConsts{compiler: c}
	ic.Process_2(1, 2)
	if !ic.found {
		t.Fatalf("expected InferConsts.found == true")
	}
	if got := ScopeCreatorReleaseCnt.Load(); got != 5 {
		t.Fatalf("expected 5 pass-boundary/inline scope releases, got %d", got)
	}

	// RemoveUnusedCode_VarInfo slab pooling: canonicalUnremovableVarInfo MUST survive ReleaseSlabs_()!
	j2gort.SlabAllocCnt.Store(0)
	j2gort.SlabReleaseCnt.Store(0)
	ruc := NewRemoveUnusedCode(c)
	if j2gort.SlabAllocCnt.Load() != 0 {
		t.Fatalf("canonicalUnremovableVarInfo in InitRemoveUnusedCode must NOT allocate from slab")
	}
	ruc.Process_2(1, 2)
	if j2gort.SlabAllocCnt.Load() != 1 || j2gort.SlabReleaseCnt.Load() != 1 {
		t.Fatalf("expected 1 VarInfo slab alloc and 1 slab release, got alloc=%d rel=%d",
			j2gort.SlabAllocCnt.Load(), j2gort.SlabReleaseCnt.Load())
	}
	if !ruc.canonicalUnremovableVarInfo.canonicalUnremovable {
		t.Fatalf("canonicalUnremovableVarInfo was clobbered by ReleaseSlabs_()!")
	}

	// GlobalNamespace_Ref slab pooling + Es6RewriteClass release + nil-outer safety
	j2gort.SlabAllocCnt.Store(0)
	j2gort.SlabReleaseCnt.Store(0)
	erc := &Es6RewriteClass{
		compiler:                     c,
		convertSuperConstructorCalls: &Es6ConvertSuperConstructorCalls{},
	}
	erc.Process_2(1, 2)
	if erc.convertSuperConstructorCalls.globalNamespace != nil {
		t.Fatalf("expected convertSuperConstructorCalls.globalNamespace nilled on exit")
	}
	if j2gort.SlabAllocCnt.Load() < 1 || j2gort.SlabReleaseCnt.Load() < 1 {
		t.Fatalf("expected GlobalNamespace_Ref slab alloc & release during Es6RewriteClass.Process_2, got alloc=%d rel=%d",
			j2gort.SlabAllocCnt.Load(), j2gort.SlabReleaseCnt.Load())
	}
	standaloneName := NewGlobalNamespace_Name(nil, "standalone")
	if r := standaloneName.AddSingleRefForTesting(pkg_rhino.Node(3), 1); r == nil || r.node != 3 {
		t.Fatalf("expected nil-outer GlobalNamespace_Name.AddSingleRefForTesting fallback to succeed")
	}
}
`)

	cmd := exec.Command("go", "test", "-v", "./gen/jscomp")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test ./gen/jscomp failed: %v\n%s", err, string(out))
	}
}
