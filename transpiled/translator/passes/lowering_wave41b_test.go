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

func wave41bWriteFiles(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for rel, content := range files {
		full := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// TestWave41B_StringPoolInternAtUsesInternOwned verifies Task 1:
//   - Fresh lowering of (*StringPool).internAt_ emits nodestore.InternOwned(s_) instead of alloc_.Intern(s_).
//   - Upgrading an already-lowered StringPool.internAt_ replaces alloc_.Intern(s_) with nodestore.InternOwned(s_).
//   - NullString returns (0, false), atomic cache + Epoch invalidation work under -race, and Run 2 produces 0 sites.
func TestWave41B_StringPoolInternAtUsesInternOwned(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	nodestoreSrc := `package nodestore

import (
	"sync"
	"sync/atomic"
)

type Node uint32
type StringID uint32

const EmptyStringID StringID = 1

var (
	mu               sync.RWMutex
	epoch            atomic.Uint64
	strs             = []string{"", ""}
	strMap           = map[string]StringID{"": EmptyStringID}
	InternCalls      atomic.Int64
	InternOwnedCalls atomic.Int64
)

func init() { epoch.Store(1) }
func Epoch() uint64     { return epoch.Load() }
func BumpEpochForTest() { epoch.Add(1) }

func Intern(s string) StringID {
	InternCalls.Add(1)
	return internImpl(s)
}

func InternOwned(s string) StringID {
	InternOwnedCalls.Add(1)
	return internImpl(s)
}

func internImpl(s string) StringID {
	if len(s) == 0 {
		return EmptyStringID
	}
	mu.Lock()
	defer mu.Unlock()
	if id, ok := strMap[s]; ok {
		return id
	}
	id := StringID(len(strs))
	strs = append(strs, s)
	strMap[s] = id
	return id
}

func String(id StringID) string {
	if id == 0 {
		return "\xff\x00j2go_null"
	}
	if id == EmptyStringID {
		return ""
	}
	mu.RLock()
	defer mu.RUnlock()
	return strs[id]
}

type NodeRec struct {
	Token    uint8
	StrID    StringID
	OrigID   StringID
	Double   float64
	Tmpl     Node
}

var (
	recMu sync.RWMutex
	recs  = make([]NodeRec, 1, 256)
)

func (n Node) SetOrigNameID(id StringID) {
	recMu.Lock()
	defer recMu.Unlock()
	recs[n].OrigID = id
}

func (n Node) OrigNameID() StringID {
	recMu.RLock()
	defer recMu.RUnlock()
	return recs[n].OrigID
}

func GetRec(n Node) NodeRec {
	recMu.RLock()
	defer recMu.RUnlock()
	return recs[n]
}

type NodeAllocator struct {
	tmpl Node
}

func (a *NodeAllocator) InitLocal()            { a.tmpl = 0 }
func (a *NodeAllocator) Flush()                { a.tmpl = 0 }
func (a *NodeAllocator) SetTemplate(tmpl Node) { a.tmpl = tmpl }
func (a *NodeAllocator) Intern(s string) StringID {
	return Intern(s)
}
func (a *NodeAllocator) AllocNode1(tok uint8, tmpl Node) Node {
	recMu.Lock()
	defer recMu.Unlock()
	id := Node(len(recs))
	recs = append(recs, NodeRec{Token: tok, Tmpl: tmpl})
	return id
}
func (a *NodeAllocator) AllocStringNode(tok uint8, strID StringID, tmpl Node) Node {
	recMu.Lock()
	defer recMu.Unlock()
	id := Node(len(recs))
	recs = append(recs, NodeRec{Token: tok, StrID: strID, Tmpl: tmpl})
	return id
}
func (a *NodeAllocator) AllocNumberNode(tok uint8, val float64, tmpl Node) Node {
	recMu.Lock()
	defer recMu.Unlock()
	id := Node(len(recs))
	recs = append(recs, NodeRec{Token: tok, Double: val, Tmpl: tmpl})
	return id
}
`

	j2gortSrc := `package j2gort

const NullString = "\xff\x00j2go_null"

func ThrowNPE() { panic("NPE") }
func CheckNotNullPtr[T any](p *T) *T {
	if p == nil {
		ThrowNPE()
	}
	return p
}

type Slab[T any] struct{ items []T }
func (s *Slab[T]) Alloc() *T {
	s.items = append(s.items, *new(T))
	return &s.items[len(s.items)-1]
}
func (s *Slab[T]) Release() { s.items = nil }

type PooledSlab[T any] struct{ items []T }
func (s *PooledSlab[T]) Alloc() *T {
	s.items = append(s.items, *new(T))
	return &s.items[len(s.items)-1]
}
func (s *PooledSlab[T]) Reset()   { clear(s.items); s.items = s.items[:0] }
func (s *PooledSlab[T]) Release() { clear(s.items); s.items = nil }
func (s *PooledSlab[T]) ReleaseWith(fn func(*T)) {
	for i := range s.items {
		if fn != nil {
			fn(&s.items[i])
		}
	}
	s.Release()
}
`

	rhinoSrc := `package rhino

import "example.com/w41b_sp/shims/nodestore"

type Node uint32
type Token uint8

const (
	TokNUMBER    Token = 40
	TokSTRINGLIT Token = 41
	TokSCRIPT    Token = 90
)

func StaticGet_Token_NUMBER() Token    { return TokNUMBER }
func StaticGet_Token_STRINGLIT() Token { return TokSTRINGLIT }
func StaticGet_Token_SCRIPT() Token    { return TokSCRIPT }

func NewNode_1(tok Token) Node {
	var a nodestore.NodeAllocator
	return Node(a.AllocNode1(uint8(tok)-1, 0))
}
func IR_Number(d float64) Node {
	var a nodestore.NodeAllocator
	return Node(a.AllocNumberNode(uint8(TokNUMBER)-1, d, 0))
}
func (n Node) ClonePropsFrom(_ Node) Node { return n }
func (n Node) AddChildToBack(_ Node)      {}
func (n Node) GetString() string          { return nodestore.String(nodestore.GetRec(nodestore.Node(n)).StrID) }
func (n Node) GetOrigName() string        { return nodestore.String(nodestore.Node(n).OrigNameID()) }

type RhinoStringPool_LazyInternedStringList struct {
	items []string
}

func NewLazyList(items []string) *RhinoStringPool_LazyInternedStringList {
	return &RhinoStringPool_LazyInternedStringList{items: items}
}
func (l *RhinoStringPool_LazyInternedStringList) Size() int32      { return int32(len(l.items)) }
func (l *RhinoStringPool_LazyInternedStringList) Get(i int32) string { return l.items[i] }
`

	stringPoolSrc := `package jscomp

import (
	"example.com/w41b_sp/gen/j2gort"
	pkg_rhino "example.com/w41b_sp/gen/rhino"
)

type StringPool struct {
	pool      *pkg_rhino.RhinoStringPool_LazyInternedStringList
	maxLength int32
}

func NewStringPool(items []string) *StringPool {
	sp := &StringPool{}
	InitStringPool(sp, items)
	return sp
}

func InitStringPool(this_ *StringPool, items []string) {
	_ = j2gort.NullString
	this_.pool = pkg_rhino.NewLazyList(items)
}

func (this_ *StringPool) Get(offset int32) string {
	return this_.pool.Get(offset)
}
`

	sndSrc := `package jscomp

import (
	pkg_rhino "example.com/w41b_sp/gen/rhino"
)

type SourceFile struct{ Name string }

type AstNode struct {
	Kind     int32
	StrOff   int32
	OrigOff  int32
	NumVal   float64
	Children []*AstNode
}

func (a *AstNode) GetStringValuePointer() int32  { return a.StrOff }
func (a *AstNode) GetOriginalNamePointer() int32 { return a.OrigOff }
func (a *AstNode) GetDoubleValue() float64       { return a.NumVal }

type ScriptNodeDeserializer struct {
	stringPool *StringPool
	rootAst    *AstNode
}

func NewScriptNodeDeserializer(sp *StringPool, root *AstNode) *ScriptNodeDeserializer {
	return &ScriptNodeDeserializer{stringPool: sp, rootAst: root}
}

func (this_ *ScriptNodeDeserializer) DeserializeNew() pkg_rhino.Node {
	return NewScriptNodeDeserializer_Runner(this_).Run()
}

func (this_ *ScriptNodeDeserializer) priv_ScriptNodeDeserializer_setOriginalNameIfPresent_1(astNode *AstNode, n pkg_rhino.Node) {
	_ = astNode
	_ = n
}

func (this_ *ScriptNodeDeserializer) priv_ScriptNodeDeserializer_stringNode_2(token pkg_rhino.Token, n *AstNode, templateNode pkg_rhino.Node) pkg_rhino.Node {
	_ = this_.stringPool.Get(n.GetStringValuePointer())
	return pkg_rhino.NewNode_1(token).ClonePropsFrom(templateNode)
}

func (this_ *ScriptNodeDeserializer) priv_ScriptNodeDeserializer_newNodeFromTemplate_3(token pkg_rhino.Token, templateNode pkg_rhino.Node) pkg_rhino.Node {
	return pkg_rhino.NewNode_1(token).ClonePropsFrom(templateNode)
}

func (this_ *ScriptNodeDeserializer) priv_ScriptNodeDeserializer_deserializeSingleNode_4(n *AstNode, templateNode pkg_rhino.Node) pkg_rhino.Node {
	switch n.Kind {
	case 1:
		var numberNode pkg_rhino.Node = pkg_rhino.IR_Number(n.GetDoubleValue())
		numberNode.ClonePropsFrom(templateNode)
		return numberNode
	case 2:
		return this_.priv_ScriptNodeDeserializer_stringNode_2(pkg_rhino.StaticGet_Token_STRINGLIT(), n, templateNode)
	default:
		return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_3(pkg_rhino.StaticGet_Token_SCRIPT(), templateNode)
	}
}

type ScriptNodeDeserializer_Runner struct {
	outer_ *ScriptNodeDeserializer
}

func NewScriptNodeDeserializer_Runner(outer_ *ScriptNodeDeserializer) *ScriptNodeDeserializer_Runner {
	r := &ScriptNodeDeserializer_Runner{outer_: outer_}
	InitScriptNodeDeserializer_Runner(r)
	return r
}

func InitScriptNodeDeserializer_Runner(_ *ScriptNodeDeserializer_Runner) {}

func (this_ *ScriptNodeDeserializer_Runner) Run() pkg_rhino.Node {
	tmpl := pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_SCRIPT())
	return this_.priv_Runner_visit_5(this_.outer_.rootAst, nil, nil, tmpl)
}

func (this_ *ScriptNodeDeserializer_Runner) priv_Runner_visit_5(astNode *AstNode, c1 any, c2 any, sourceFileTemplate pkg_rhino.Node) pkg_rhino.Node {
	_ = c1
	_ = c2
	n := this_.outer_.priv_ScriptNodeDeserializer_deserializeSingleNode_4(astNode, sourceFileTemplate)
	for _, ch := range astNode.Children {
		cn := this_.priv_Runner_visit_5(ch, c1, c2, sourceFileTemplate)
		this_.outer_.priv_ScriptNodeDeserializer_setOriginalNameIfPresent_1(ch, cn)
		n.AddChildToBack(cn)
	}
	return n
}
`

	rtTestSrc := `package jscomp

import (
	"sync"
	"testing"

	"example.com/w41b_sp/gen/j2gort"
	"example.com/w41b_sp/shims/nodestore"
)

func TestStringPoolInternOwnedBehaviorAndRace(t *testing.T) {
	sp := NewStringPool([]string{"", "propAlpha", "origBeta", j2gort.NullString})
	ast := &AstNode{
		Kind: 100,
		Children: []*AstNode{
			{Kind: 2, StrOff: 1, OrigOff: 2},
			{Kind: 2, StrOff: 3, OrigOff: 3},
		},
	}
	snd := NewScriptNodeDeserializer(sp, ast)

	nodestore.InternCalls.Store(0)
	nodestore.InternOwnedCalls.Store(0)

	_ = snd.DeserializeNew()
	if got := nodestore.InternCalls.Load(); got != 0 {
		t.Fatalf("expected 0 cloning nodestore.Intern calls, got %d", got)
	}
	if got := nodestore.InternOwnedCalls.Load(); got != 2 {
		t.Fatalf("expected 2 nodestore.InternOwned calls for 2 non-null pool strings, got %d", got)
	}

	// Second DeserializeNew MUST hit atomic internedIDs_ cache (0 new InternOwned calls)!
	_ = snd.DeserializeNew()
	if got := nodestore.InternOwnedCalls.Load(); got != 2 {
		t.Fatalf("expected 0 additional InternOwned calls on cached pool, got total=%d", got)
	}

	// Epoch bump invalidates cache and still calls InternOwned (NOT Intern).
	nodestore.BumpEpochForTest()
	_ = snd.DeserializeNew()
	if got := nodestore.InternCalls.Load(); got != 0 {
		t.Fatalf("expected 0 Intern calls after Epoch bump, got %d", got)
	}
	if got := nodestore.InternOwnedCalls.Load(); got != 4 {
		t.Fatalf("expected fallback after Epoch bump to use InternOwned (total=4), got %d", got)
	}

	// Concurrent -race stress test on shared StringPool.
	var wg sync.WaitGroup
	for g := 0; g < 6; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sp2 := NewStringPool([]string{"", "k1", "k2", j2gort.NullString})
			s2 := NewScriptNodeDeserializer(sp2, ast)
			for i := 0; i < 20; i++ {
				_ = s2.DeserializeNew()
			}
		}()
	}
	wg.Wait()
	if got := nodestore.InternCalls.Load(); got != 0 {
		t.Fatalf("expected 0 Intern calls after concurrent stress, got %d", got)
	}
}
`

	wave41bWriteFiles(t, tmpDir, map[string]string{
		"go.mod":                                 "module example.com/w41b_sp\n\ngo 1.24\n",
		"shims/nodestore/nodestore.go":           nodestoreSrc,
		"gen/j2gort/j2gort.go":                   j2gortSrc,
		"gen/rhino/rhino.go":                     rhinoSrc,
		"gen/jscomp/string_pool.go":              stringPoolSrc,
		"gen/jscomp/script_node_deserializer.go": sndSrc,
	})

	jscompDir := filepath.Join(tmpDir, "gen/jscomp")
	res1, err := passes.RunOnPackageDir(jscompDir, passes.AllPasses(), passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if res1.TotalSites == 0 {
		t.Fatalf("expected non-zero sites on Run 1, got %+v", res1)
	}

	spBytes, err := os.ReadFile(filepath.Join(jscompDir, "string_pool.go"))
	if err != nil {
		t.Fatal(err)
	}
	spCode := string(spBytes)
	if !strings.Contains(spCode, "nodestore.InternOwned(s_)") {
		t.Fatalf("expected string_pool.go internAt_ to call nodestore.InternOwned(s_):\n%s", spCode)
	}
	if strings.Contains(spCode, "alloc_.Intern(") {
		t.Fatalf("expected string_pool.go internAt_ NOT to call alloc_.Intern:\n%s", spCode)
	}

	// Run 2 MUST be 100% idempotent (0 sites).
	res2, err := passes.RunOnPackageDir(jscompDir, passes.AllPasses(), passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if res2.TotalSites != 0 || res2.FilesModified != 0 {
		t.Fatalf("Run 2 NOT IDEMPOTENT: %+v", res2)
	}

	// Also verify upgrading a legacy `alloc_.Intern(s_)` string_pool.go is idempotent!
	legacySP := strings.ReplaceAll(spCode, "nodestore.InternOwned(s_)", "alloc_.Intern(s_)")
	if err := os.WriteFile(filepath.Join(jscompDir, "string_pool.go"), []byte(legacySP), 0o644); err != nil {
		t.Fatal(err)
	}
	up1, err := passes.RunOnPackageDir(jscompDir, passes.AllPasses(), passes.Options{})
	if err != nil {
		t.Fatalf("Legacy upgrade Run 1 failed: %v", err)
	}
	if up1.TotalSites != 1 {
		t.Fatalf("expected 1 upgrade site on legacy internAt_, got %+v", up1)
	}
	up2, err := passes.RunOnPackageDir(jscompDir, passes.AllPasses(), passes.Options{})
	if err != nil {
		t.Fatalf("Legacy upgrade Run 2 failed: %v", err)
	}
	if up2.TotalSites != 0 || up2.FilesModified != 0 {
		t.Fatalf("Legacy upgrade Run 2 NOT IDEMPOTENT: %+v", up2)
	}

	wave41bWriteFiles(t, tmpDir, map[string]string{
		"gen/jscomp/sp_rt_test.go": rtTestSrc,
	})

	cmd := exec.Command("go", "test", "-race", "-v", "./gen/jscomp")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test -race failed: %v\n%s", err, string(out))
	}
}

// TestWave41B_UnwrapStmtValueIIFEs verifies Task 2:
//   - Pattern A (`var ok bool = func() bool { ... }()`, `x := func() T { ... }()`, generic `R`/`V`)
//   - Pattern B (`lhsIdent = func() T { ... }()` reading lhsIdent & early-exiting loops)
//   - Pattern C (`return func() T { ... }()`)
//   - Tail-only blocks emit plain `{ ... }` with zero labels/switches; early-return / loop-return /
//     trailing-panic switch blocks emit `iifeSwitch_N: switch { default: ... break iifeSwitch_N }`.
//   - Nested non-IIFE closures, defer/go/goto/bare-return IIFEs, shadowed vars, and field/index LHS are preserved.
//   - Unlowered vs lowered behavioral & panic parity + Run-2 0-site idempotency.
func TestWave41B_UnwrapStmtValueIIFEs(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	workloadSrc := `package sample

type Box struct {
	Val int32
}

// Mirrors TypeValidator.checkWellFormedTemplatizedType: switch with inner loops,
// early returns BEFORE var declarations in cases, and trailing panic("unreachable switch expression").
func CheckWellFormedSwitch(kind int32, items []int32) (res int32) {
	var ok bool = func() bool {
		switch kind {
		case int32(128):
			var paramsOk bool = true
			for i := 0; i < len(items); i++ {
				if items[i] < 0 {
					paramsOk = false
					break
				}
			}
			if !paramsOk {
				return false
			}
			var sumAfterEarlyRet int32 = 0
			for _, v := range items {
				sumAfterEarlyRet += v
			}
			return sumAfterEarlyRet <= 100
		case int32(64):
			for _, v := range items {
				if v == 0 {
					continue
				}
				if v == 999 {
					return false
				}
			}
			var finalOk bool = len(items) > 0
			return finalOk
		case int32(256):
			return true
		}
		panic("unreachable switch expression")
	}()
	if ok {
		return 1
	}
	return 0
}

// Tail-only ternary & exhaustive switch without trailing panic -> plain { ... } with NO iifeSwitch_ label!
func TailOnlyPatterns(x int32, flag bool) int32 {
	var a int32 = func() int32 {
		if flag {
			return x + 10
		} else {
			return x - 10
		}
	}()
	b := func() int32 {
		switch x {
		case 1:
			return 100
		case 2:
			return 200
		default:
			return 300
		}
	}()
	// Pattern B reading a before returning updated value.
	a = func() int32 {
		prev := a
		delta := b
		return prev + delta
	}()
	return a
}

// Pattern C (return func() T { ... }()) + nested IIFEs + generic type parameters (like SimpleTable[R, V]).
type GenericHolder[R comparable, V any] struct {
	Data map[R]V
}

func (g *GenericHolder[R, V]) LookupOrZero(key R, fallback V, useFallback bool) V {
	resolved := func() V {
		if v, ok := g.Data[key]; ok {
			return v
		}
		if useFallback {
			return fallback
		}
		return func() V {
			var z V
			return z
		}()
	}()
	return func() V {
		if useFallback && len(g.Data) == 0 {
			return fallback
		}
		return resolved
	}()
}

// Nested non-IIFE closure inside an unwrapped IIFE MUST be preserved!
func NestedClosureInsideIIFE(base int32) int32 {
	fn := func() func(int32) int32 {
		if base > 0 {
			mult := base * 2
			return func(v int32) int32 {
				return v + mult
			}
		}
		return func(v int32) int32 {
			return v - 1
		}
	}()
	return fn(5)
}

// Negative safety guards: all of these MUST remain untouched!
var globalSink int32

func NegativeSafetyGuards(x int32, box *Box, arr []int32) int32 {
	// 1. Outer x referenced inside var x int32 = func() int32 { return x + 1 }() MUST NOT shadow outer x!
	{
		var x int32 = func() int32 {
			return x + 100
		}()
		box.Val = x
	}
	// 2. Defer inside IIFE MUST NOT unwrap.
	d := func() int32 {
		defer func() { box.Val++ }()
		return 7
	}()
	// 3. Goto inside IIFE MUST NOT unwrap.
	g := func() int32 {
		if x < 0 {
			goto neg
		}
		return 1
	neg:
		return -1
	}()
	// 4. Field & slice index LHS MUST NOT unwrap via Pattern B.
	box.Val = func() int32 { return box.Val + d }()
	arr[0] = func() int32 { return g + 1 }()
	// 5. Global variable LHS MUST NOT unwrap via Pattern B.
	globalSink = func() int32 { return x }()
	// 6. Enclosing function with defer MUST NOT unwrap Pattern C.
	defer func() {}()
	return func() int32 {
		return box.Val + arr[0] + globalSink
	}()
}
`

	origSrc := strings.Replace(workloadSrc, "package sample", "package orig", 1)
	optSrc := strings.Replace(workloadSrc, "package sample", "package opt", 1)

	parityTestSrc := `package parity_test

import (
	"testing"

	"example.com/w41b_iife/opt"
	"example.com/w41b_iife/orig"
)

func runCatch(fn func() int32) (val int32, panicked bool) {
	defer func() {
		if recover() != nil {
			panicked = true
		}
	}()
	return fn(), false
}

func TestIIFEBehavioralAndPanicParity(t *testing.T) {
	cases := []struct {
		kind  int32
		items []int32
	}{
		{128, []int32{10, 20, 30}},
		{128, []int32{10, -1, 30}},
		{128, []int32{60, 60}},
		{64, []int32{0, 1, 2}},
		{64, []int32{0, 999, 2}},
		{64, nil},
		{256, nil},
		{9999, []int32{1}}, // hits panic("unreachable switch expression")!
	}
	for _, tc := range cases {
		vOrig, pOrig := runCatch(func() int32 { return orig.CheckWellFormedSwitch(tc.kind, tc.items) })
		vOpt, pOpt := runCatch(func() int32 { return opt.CheckWellFormedSwitch(tc.kind, tc.items) })
		if pOrig != pOpt || vOrig != vOpt {
			t.Fatalf("CheckWellFormedSwitch(%d, %v): orig=(%d,%v) != opt=(%d,%v)", tc.kind, tc.items, vOrig, pOrig, vOpt, pOpt)
		}
	}

	for _, x := range []int32{0, 1, 2, 5} {
		for _, fl := range []bool{false, true} {
			if g1, g2 := orig.TailOnlyPatterns(x, fl), opt.TailOnlyPatterns(x, fl); g1 != g2 {
				t.Fatalf("TailOnlyPatterns(%d, %v): orig=%d != opt=%d", x, fl, g1, g2)
			}
		}
	}

	ghOrig := &orig.GenericHolder[string, int32]{Data: map[string]int32{"a": 42}}
	ghOpt := &opt.GenericHolder[string, int32]{Data: map[string]int32{"a": 42}}
	for _, k := range []string{"a", "missing"} {
		for _, uf := range []bool{false, true} {
			if g1, g2 := ghOrig.LookupOrZero(k, 99, uf), ghOpt.LookupOrZero(k, 99, uf); g1 != g2 {
				t.Fatalf("LookupOrZero(%q, %v): orig=%d != opt=%d", k, uf, g1, g2)
			}
		}
	}

	for _, b := range []int32{-3, 0, 4} {
		if g1, g2 := orig.NestedClosureInsideIIFE(b), opt.NestedClosureInsideIIFE(b); g1 != g2 {
			t.Fatalf("NestedClosureInsideIIFE(%d): orig=%d != opt=%d", b, g1, g2)
		}
	}

	for _, x := range []int32{-5, 0, 12} {
		b1, a1 := &orig.Box{Val: 3}, []int32{9}
		b2, a2 := &opt.Box{Val: 3}, []int32{9}
		r1 := orig.NegativeSafetyGuards(x, b1, a1)
		r2 := opt.NegativeSafetyGuards(x, b2, a2)
		if r1 != r2 || b1.Val != b2.Val || a1[0] != a2[0] {
			t.Fatalf("NegativeSafetyGuards(%d): orig=(%d,%d,%d) != opt=(%d,%d,%d)", x, r1, b1.Val, a1[0], r2, b2.Val, a2[0])
		}
	}
}
`

	wave41bWriteFiles(t, tmpDir, map[string]string{
		"go.mod":              "module example.com/w41b_iife\n\ngo 1.24\n",
		"orig/sample.go":      origSrc,
		"opt/sample.go":       optSrc,
		"parity/iife_test.go": parityTestSrc,
	})

	optDir := filepath.Join(tmpDir, "opt")
	res1, err := passes.RunOnPackageDir(optDir, passes.AllPasses(), passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if res1.TotalSites < 7 {
		t.Fatalf("expected >= 7 IIFE unwrap sites on Run 1, got %+v", res1)
	}

	// Run 2 MUST be 100% idempotent (0 sites, 0 modified files).
	res2, err := passes.RunOnPackageDir(optDir, passes.AllPasses(), passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if res2.TotalSites != 0 || res2.FilesModified != 0 {
		t.Fatalf("Run 2 NOT IDEMPOTENT: %+v", res2)
	}

	outBytes, err := os.ReadFile(filepath.Join(optDir, "sample.go"))
	if err != nil {
		t.Fatal(err)
	}
	code := string(outBytes)

	checkFn := code[:strings.Index(code, "func TailOnlyPatterns")]
	if strings.Contains(checkFn, "func()") {
		t.Fatalf("expected CheckWellFormedSwitch to have 0 remaining func() closures:\n%s", checkFn)
	}
	if !strings.Contains(checkFn, "iifeSwitch_") || !strings.Contains(checkFn, "break iifeSwitch_") {
		t.Fatalf("expected CheckWellFormedSwitch to use labeled switch break:\n%s", checkFn)
	}

	tailFn := code[strings.Index(code, "func TailOnlyPatterns"):strings.Index(code, "// Pattern C")]
	if strings.Contains(tailFn, "func()") || strings.Contains(tailFn, "iifeSwitch_") {
		t.Fatalf("expected TailOnlyPatterns to have 0 closures AND 0 iifeSwitch_ labels:\n%s", tailFn)
	}

	genFn := code[strings.Index(code, "type GenericHolder"):strings.Index(code, "func NestedClosureInsideIIFE")]
	if strings.Contains(genFn, "func()") {
		t.Fatalf("expected GenericHolder.LookupOrZero to unwrap all 3 nested/return IIFEs:\n%s", genFn)
	}

	nestedFn := code[strings.Index(code, "func NestedClosureInsideIIFE"):strings.Index(code, "var globalSink")]
	if strings.Contains(nestedFn, "func() func(") {
		t.Fatalf("expected outer IIFE in NestedClosureInsideIIFE to unwrap:\n%s", nestedFn)
	}
	if strings.Count(nestedFn, "func(v int32) int32") != 2 {
		t.Fatalf("expected both inner closures in NestedClosureInsideIIFE to be preserved:\n%s", nestedFn)
	}

	cmd := exec.Command("go", "test", "-race", "-v", "./parity")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test -race ./parity failed: %v\n%s", err, string(out))
	}
}

// TestWave41B_VarExpertScratchAndMinimizedConditionLazyNodes verifies Task 3:
//   - InlineVariables_InliningBehavior reuses scratchVarExpert_ in L1 with 0 slab_InlineVariables_VarExpert
//     allocations for non-escaping vars while promoting LazyNegative and AliasInlineRetryHandler experts.
//   - MinimizedCondition_MeasuredNode defers NewNode_1(...).Srcref(...) via lazyTok_/srcrefNode_, performing
//     ZERO NewNode_1 calls when conditions do not change and exact materialization when BuildReplacement runs.
func TestWave41B_VarExpertScratchAndMinimizedConditionLazyNodes(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()

	j2Src := `package j2gort

import "sync/atomic"

func CheckState(cond bool) {
	if !cond {
		panic("CheckState")
	}
}
func CheckNotNullNode[N ~uint32](n N) N {
	if n == 0 {
		panic("NPE")
	}
	return n
}

type Slab[T any] struct{ items []T }
func (s *Slab[T]) Alloc() *T {
	s.items = append(s.items, *new(T))
	return &s.items[len(s.items)-1]
}
func (s *Slab[T]) Release() { s.items = nil }

type PooledSlab[T any] struct {
	items      []T
	AllocCalls atomic.Int64
}

func (s *PooledSlab[T]) Alloc() *T {
	s.AllocCalls.Add(1)
	s.items = append(s.items, *new(T))
	return &s.items[len(s.items)-1]
}
func (s *PooledSlab[T]) Reset()   { clear(s.items); s.items = s.items[:0] }
func (s *PooledSlab[T]) Release() { clear(s.items); s.items = nil }
func (s *PooledSlab[T]) ReleaseWith(fn func(*T)) {
	for i := range s.items {
		if fn != nil {
			fn(&s.items[i])
		}
	}
	s.Release()
}
`

	rhinoSrc := `package rhino

import "sync/atomic"

type Node uint32
type Token uint8

const (
	TokEQ   Token = 5
	TokNE   Token = 6
	TokNOT  Token = 20
	TokSHEQ Token = 41
	TokSHNE Token = 42
	TokHOOK Token = 68
	TokOR   Token = 69
	TokAND  Token = 70
	TokIF   Token = 75
	TokNAME Token = 39
)

func StaticGet_Token_EQ() Token   { return TokEQ }
func StaticGet_Token_NE() Token   { return TokNE }
func StaticGet_Token_NOT() Token  { return TokNOT }
func StaticGet_Token_SHEQ() Token { return TokSHEQ }
func StaticGet_Token_SHNE() Token { return TokSHNE }
func StaticGet_Token_OR() Token   { return TokOR }
func StaticGet_Token_AND() Token  { return TokAND }

type nodeRec struct {
	tok      Token
	srcref   Node
	parent   Node
	children []Node
}

var (
	recs          = []nodeRec{{}}
	NewNode1Calls atomic.Int64
)

func ResetNodes() {
	recs = []nodeRec{{}}
	NewNode1Calls.Store(0)
}

func AllocTestNode(tok Token, children ...Node) Node {
	id := Node(len(recs))
	recs = append(recs, nodeRec{tok: tok})
	for _, c := range children {
		id.AddChildToBack(c)
	}
	return id
}

func NewNode_1(tok Token) Node {
	NewNode1Calls.Add(1)
	id := Node(len(recs))
	recs = append(recs, nodeRec{tok: tok})
	return id
}

func (n Node) Srcref(src Node) Node {
	if src == 0 {
		panic("Srcref(0)")
	}
	recs[n].srcref = src
	return n
}
func (n Node) SrcRefOf() Node        { return recs[n].srcref }
func (n Node) GetToken() Token       { return recs[n].tok }
func (n Node) TokenSwitchOrd_() int32 { return int32(recs[n].tok) - 1 }
func (n Node) IsNot() bool           { return recs[n].tok == TokNOT }
func (n Node) IsAnd() bool           { return recs[n].tok == TokAND }
func (n Node) IsIf() bool            { return recs[n].tok == TokIF }
func (n Node) HasParent() bool       { return recs[n].parent != 0 }
func (n Node) GetParent() Node       { return recs[n].parent }
func (n Node) GetFirstChild() Node {
	if len(recs[n].children) == 0 {
		return 0
	}
	return recs[n].children[0]
}
func (n Node) GetLastChild() Node {
	if len(recs[n].children) == 0 {
		return 0
	}
	return recs[n].children[len(recs[n].children)-1]
}
func (n Node) GetNext() Node {
	p := recs[n].parent
	if p == 0 {
		return 0
	}
	cs := recs[p].children
	for i, c := range cs {
		if c == n && i+1 < len(cs) {
			return cs[i+1]
		}
	}
	return 0
}
func (n Node) DetachChildren() {
	for _, c := range recs[n].children {
		recs[c].parent = 0
	}
	recs[n].children = nil
}
func (n Node) Detach() Node {
	p := recs[n].parent
	if p != 0 {
		cs := recs[p].children
		for i, c := range cs {
			if c == n {
				recs[p].children = append(cs[:i], cs[i+1:]...)
				break
			}
		}
		recs[n].parent = 0
	}
	return n
}
func (n Node) AddChildToBack(c Node) {
	recs[c].parent = n
	recs[n].children = append(recs[n].children, c)
}
func (n Node) ReplaceWith(repl Node) {
	p := recs[n].parent
	cs := recs[p].children
	for i, c := range cs {
		if c == n {
			cs[i] = repl
			recs[n].parent = 0
			recs[repl].parent = p
			return
		}
	}
}
`

	inlineVarsSrc := `package jscomp

import (
	"sync/atomic"

	pkg_rhino "example.com/w41b_task3/gen/rhino"
)

type AbstractCompiler interface{ IsJavaNil_() bool }
type Var struct {
	name  string
	kind  int // 0=inlineWellDefined, 1=lazyNegative, 2=waitForAlias
	alias *Var
}
type Reference struct{ node pkg_rhino.Node }
type ReferenceCollection struct{ refs []*Reference }
func (rc *ReferenceCollection) Get(i int32) *Reference { return rc.refs[i] }
func (rc *ReferenceCollection) Size() int32            { return int32(len(rc.refs)) }

type ReferenceCollector struct {
	beh  *InlineVariables_InliningBehavior
	vars []*Var
	cols map[*Var]*ReferenceCollection
}
func InitReferenceCollector(rc *ReferenceCollector, beh *InlineVariables_InliningBehavior, vars []*Var, cols map[*Var]*ReferenceCollection) {
	rc.beh = beh
	rc.vars = vars
	rc.cols = cols
}
func (rc *ReferenceCollector) Process_2(_, _ pkg_rhino.Node) {
	rc.beh.priv_InliningBehavior_doInlinesForScope_1(rc)
}

var (
	static_NO_INLINE_ANALYSIS = NewInlineVariables_InlineVarAnalysis(0, nil, nil, 0, nil, nil)
	static_NO_INLINE_EXPERT   = NewInlineVariables_VarExpert_1(static_NO_INLINE_ANALYSIS)
)
func StaticGet_NO_INLINE_ANALYSIS() *InlineVariables_InlineVarAnalysis { return static_NO_INLINE_ANALYSIS }
func StaticGet_NO_INLINE_EXPERT() *InlineVariables_VarExpert           { return static_NO_INLINE_EXPERT }

type InlineVariables struct {
	compiler     AbstractCompiler
	vars         []*Var
	cols         map[*Var]*ReferenceCollection
	inlinedNames []string
	lazySafeCnt  int
}

func (this_ *InlineVariables) initFields_InlineVariables() {}

func (this_ *InlineVariables) Process_2(externs, root pkg_rhino.Node) {
	var rc ReferenceCollector
	InitReferenceCollector(&rc, NewInlineVariables_InliningBehavior(this_), this_.vars, this_.cols)
	rc.Process_2(externs, root)
}

type InlineVariables_InlineVarAnalysis struct {
	aliasedVar          *Var
	expert              *InlineVariables_VarExpert
	initialization      *Reference
	singleReadReference *Reference
	initValue           pkg_rhino.Node
	idHash_             atomic.Int32
	kind                int8
}

func (a *InlineVariables_InlineVarAnalysis) initFields_InlineVariables_InlineVarAnalysis() {}

func NewInlineVariables_InlineVarAnalysis(kind int8, aliasedVar *Var, expert *InlineVariables_VarExpert, initValue pkg_rhino.Node, initialization, singleReadReference *Reference) *InlineVariables_InlineVarAnalysis {
	this_ := &InlineVariables_InlineVarAnalysis{}
	InitInlineVariables_InlineVarAnalysis(this_, kind, aliasedVar, expert, initValue, initialization, singleReadReference)
	return this_
}

func InitInlineVariables_InlineVarAnalysis(this_ *InlineVariables_InlineVarAnalysis, kind int8, aliasedVar *Var, expert *InlineVariables_VarExpert, initValue pkg_rhino.Node, initialization, singleReadReference *Reference) {
	this_.initFields_InlineVariables_InlineVarAnalysis()
	this_.kind = kind
	this_.aliasedVar = aliasedVar
	this_.expert = expert
	this_.initValue = initValue
	this_.initialization = initialization
	this_.singleReadReference = singleReadReference
}

func InlineVariables_InlineVarAnalysis_CreateLazyNegative(expert *InlineVariables_VarExpert) *InlineVariables_InlineVarAnalysis {
	return NewInlineVariables_InlineVarAnalysis(2, nil, expert, 0, nil, nil)
}
func InlineVariables_InlineVarAnalysis_CreateWaitForAlias(aliasedVar *Var) *InlineVariables_InlineVarAnalysis {
	return NewInlineVariables_InlineVarAnalysis(3, aliasedVar, nil, 0, nil, nil)
}
func InlineVariables_InlineVarAnalysis_CreateInlineWellDefined(expert *InlineVariables_VarExpert, initValue pkg_rhino.Node) *InlineVariables_InlineVarAnalysis {
	return NewInlineVariables_InlineVarAnalysis(5, nil, expert, initValue, nil, nil)
}

func (this_ *InlineVariables_InlineVarAnalysis) ShouldInline() bool             { return this_.kind >= 4 }
func (this_ *InlineVariables_InlineVarAnalysis) ShouldWaitForAliasedVar() bool  { return this_.kind == 3 }
func (this_ *InlineVariables_InlineVarAnalysis) GetAliasedVar() *Var            { return this_.aliasedVar }
func (this_ *InlineVariables_InlineVarAnalysis) IsSafeToInlineAliases() bool {
	if this_.kind == 2 {
		return this_.expert.IsWellDefinedAssignedOnce()
	}
	return this_.kind >= 4 || this_.kind == 1
}
func (this_ *InlineVariables_InlineVarAnalysis) PerformInline() {
	this_.expert.InlineWellDefined(this_.initValue)
}

type InlineVariables_VarExpert struct {
	fixedAnalysis                *InlineVariables_InlineVarAnalysis
	behavior                     *InlineVariables_InliningBehavior
	v                            *Var
	referenceInfo                *ReferenceCollection
	declaration                  *Reference
	initializationReference      *Reference
	idHash_                      atomic.Int32
	isDeclaredOrInferredConstant bool
	mayBeModArgs                 bool
	isWellDefinedAssignedOnce    int8
	initializationReferenceKnown bool
}

func (this_ *InlineVariables_VarExpert) initFields_InlineVariables_VarExpert() {
	this_.isWellDefinedAssignedOnce = int8(-1)
}

func NewInlineVariables_VarExpert_1(fixed *InlineVariables_InlineVarAnalysis) *InlineVariables_VarExpert {
	this_ := &InlineVariables_VarExpert{}
	InitInlineVariables_VarExpert_1(this_, fixed)
	return this_
}
func InitInlineVariables_VarExpert_1(this_ *InlineVariables_VarExpert, fixed *InlineVariables_InlineVarAnalysis) {
	this_.initFields_InlineVariables_VarExpert()
	this_.fixedAnalysis = fixed
}

func NewInlineVariables_VarExpert_5(behavior *InlineVariables_InliningBehavior, v *Var, referenceInfo *ReferenceCollection, isConst bool, modArgs bool) *InlineVariables_VarExpert {
	this_ := &InlineVariables_VarExpert{}
	InitInlineVariables_VarExpert_5(this_, behavior, v, referenceInfo, isConst, modArgs)
	return this_
}
func InitInlineVariables_VarExpert_5(this_ *InlineVariables_VarExpert, behavior *InlineVariables_InliningBehavior, v *Var, referenceInfo *ReferenceCollection, isConst bool, modArgs bool) {
	this_.initFields_InlineVariables_VarExpert()
	this_.fixedAnalysis = nil
	this_.behavior = behavior
	this_.v = v
	this_.referenceInfo = referenceInfo
	this_.declaration = referenceInfo.Get(0)
	this_.isDeclaredOrInferredConstant = isConst
	this_.mayBeModArgs = modArgs
}

func (this_ *InlineVariables_VarExpert) IsWellDefinedAssignedOnce() bool {
	if this_.isWellDefinedAssignedOnce == -1 {
		if this_.referenceInfo.Size() == 1 {
			this_.isWellDefinedAssignedOnce = 1
		} else {
			this_.isWellDefinedAssignedOnce = 0
		}
	}
	return this_.isWellDefinedAssignedOnce == 1
}

func (this_ *InlineVariables_VarExpert) InlineWellDefined(_ pkg_rhino.Node) {
	this_.behavior.outer_.inlinedNames = append(this_.behavior.outer_.inlinedNames, this_.v.name)
}

func (this_ *InlineVariables_VarExpert) Analyze() *InlineVariables_InlineVarAnalysis {
	if this_.fixedAnalysis != nil {
		return this_.fixedAnalysis
	}
	// Exercise initializationReferenceKnown to verify stale scratch state never leaks!
	if !this_.initializationReferenceKnown {
		this_.initializationReference = this_.declaration
		this_.initializationReferenceKnown = true
	} else {
		panic("stale initializationReferenceKnown leaked across scratchVarExpert_!")
	}
	switch this_.v.kind {
	case 0:
		return InlineVariables_InlineVarAnalysis_CreateInlineWellDefined(this_, this_.declaration.node)
	case 1:
		return InlineVariables_InlineVarAnalysis_CreateLazyNegative(this_)
	case 2:
		return InlineVariables_InlineVarAnalysis_CreateWaitForAlias(this_.v.alias)
	default:
		return StaticGet_NO_INLINE_ANALYSIS()
	}
}

func (this_ *InlineVariables_VarExpert) Reanalyze(aliased *Var, aliasedAna *InlineVariables_InlineVarAnalysis) *InlineVariables_InlineVarAnalysis {
	_ = aliased
	if aliasedAna.IsSafeToInlineAliases() {
		return InlineVariables_InlineVarAnalysis_CreateInlineWellDefined(this_, this_.declaration.node)
	}
	return StaticGet_NO_INLINE_ANALYSIS()
}

type InlineVariables_InliningBehavior struct {
	outer_ *InlineVariables
}

func NewInlineVariables_InliningBehavior(outer_ *InlineVariables) *InlineVariables_InliningBehavior {
	return &InlineVariables_InliningBehavior{outer_: outer_}
}

func (this_ *InlineVariables_InliningBehavior) priv_InliningBehavior_createVarExpert_2(v *Var, rc *ReferenceCollection) *InlineVariables_VarExpert {
	if rc == nil {
		return StaticGet_NO_INLINE_EXPERT()
	}
	return NewInlineVariables_VarExpert_5(this_, v, rc, true, false)
}

func (this_ *InlineVariables_InliningBehavior) priv_InliningBehavior_doInlinesForScope_1(rc *ReferenceCollector) {
	handled := make(map[*Var]*InlineVariables_InlineVarAnalysis)
	waiting := make(map[*Var][]*InlineVariables_InliningBehavior_AliasInlineRetryHandler)
	for _, v := range rc.vars {
		exp := this_.priv_InliningBehavior_createVarExpert_2(v, rc.cols[v])
		ana := exp.Analyze()
		if ana.ShouldWaitForAliasedVar() {
			rh := NewInlineVariables_InliningBehavior_AliasInlineRetryHandler(this_, v, exp)
			av := ana.GetAliasedVar()
			if prev := handled[av]; prev != nil {
				rh.Handle(av, prev)
			} else {
				waiting[av] = append(waiting[av], rh)
			}
		} else {
			if ana.ShouldInline() {
				ana.PerformInline()
			}
			handled[v] = ana
			for _, rh := range waiting[v] {
				rh.Handle(v, ana)
			}
		}
	}
	for _, ana := range handled {
		if ana.IsSafeToInlineAliases() {
			this_.outer_.lazySafeCnt++
		}
	}
}

type InlineVariables_InliningBehavior_AliasInlineRetryHandler struct {
	outer_  *InlineVariables_InliningBehavior
	v       *Var
	expert  *InlineVariables_VarExpert
	idHash_ atomic.Int32
}

func NewInlineVariables_InliningBehavior_AliasInlineRetryHandler(outer_ *InlineVariables_InliningBehavior, v *Var, expert *InlineVariables_VarExpert) *InlineVariables_InliningBehavior_AliasInlineRetryHandler {
	this_ := &InlineVariables_InliningBehavior_AliasInlineRetryHandler{}
	this_.outer_ = outer_
	InitInlineVariables_InliningBehavior_AliasInlineRetryHandler(this_, v, expert)
	return this_
}
func InitInlineVariables_InliningBehavior_AliasInlineRetryHandler(this_ *InlineVariables_InliningBehavior_AliasInlineRetryHandler, v *Var, expert *InlineVariables_VarExpert) {
	this_.v = v
	this_.expert = expert
}
func (this_ *InlineVariables_InliningBehavior_AliasInlineRetryHandler) Handle(av *Var, ana *InlineVariables_InlineVarAnalysis) {
	na := this_.expert.Reanalyze(av, ana)
	if na.ShouldInline() {
		na.PerformInline()
	}
}
`

	minCondSrc := `package jscomp

import (
	"sync/atomic"

	"example.com/w41b_task3/gen/j2gort"
	pkg_rhino "example.com/w41b_task3/gen/rhino"
)

func NodeUtil_Precedence(tok pkg_rhino.Token) int32 {
	switch tok {
	case pkg_rhino.TokOR:
		return 4
	case pkg_rhino.TokAND:
		return 5
	case pkg_rhino.TokEQ, pkg_rhino.TokNE, pkg_rhino.TokSHEQ, pkg_rhino.TokSHNE:
		return 8
	case pkg_rhino.TokNOT:
		return 13
	default:
		return 15
	}
}

type MinimizedCondition_MinimizationStyle uint8

const (
	stylePreferUnnegated MinimizedCondition_MinimizationStyle = 1
	styleAllowLeadingNot MinimizedCondition_MinimizationStyle = 2
)

func StaticGet_MinimizedCondition_MinimizationStyle_PREFER_UNNEGATED() MinimizedCondition_MinimizationStyle {
	return stylePreferUnnegated
}
func StaticGet_MinimizedCondition_MinimizationStyle_ALLOW_LEADING_NOT() MinimizedCondition_MinimizationStyle {
	return styleAllowLeadingNot
}

type MinimizedCondition struct {
	positive *MinimizedCondition_MeasuredNode
	negative *MinimizedCondition_MeasuredNode
	idHash_  atomic.Int32
}

var static_UNOPT_NEG = NewMinimizedCondition_MeasuredNode(0, -1, nil, nil, nil, 2147483647, true)

func clinit_MinimizedCondition() {}
func StaticGet_MinimizedCondition_UNOPTIMIZED_NEG() *MinimizedCondition_MeasuredNode {
	return static_UNOPT_NEG
}

func NewMinimizedCondition(p, n *MinimizedCondition_MeasuredNode) *MinimizedCondition {
	clinit_MinimizedCondition()
	this_ := &MinimizedCondition{}
	InitMinimizedCondition(this_, p, n)
	return this_
}
func InitMinimizedCondition(this_ *MinimizedCondition, p, n *MinimizedCondition_MeasuredNode) {
	this_.positive = p
	this_.negative = n.priv_MeasuredNode_change_3877()
}

func MinimizedCondition_IsMinimizable(n pkg_rhino.Node) bool {
	switch n.TokenSwitchOrd_() {
	case 19, 69, 68:
		return true
	default:
		return false
	}
}

func MinimizedCondition_FromConditionNode(n pkg_rhino.Node) *MinimizedCondition {
	j2gort.CheckState(n.HasParent())
	if MinimizedCondition_IsMinimizable(n) {
		return MinimizedCondition_ComputeMinimizedCondition(n)
	}
	return MinimizedCondition_Unoptimized(n)
}

func (this_ *MinimizedCondition) GetMinimized(style MinimizedCondition_MinimizationStyle) *MinimizedCondition_MeasuredNode {
	if style == StaticGet_MinimizedCondition_MinimizationStyle_PREFER_UNNEGATED() || this_.positive.node.IsNot() || this_.positive.length <= this_.negative.length {
		return this_.positive
	}
	return this_.negative.priv_MeasuredNode_addNot_3876()
}

func MinimizedCondition_Unoptimized(n pkg_rhino.Node) *MinimizedCondition {
	var val_1 MinimizedCondition_MeasuredNode
	InitMinimizedCondition_MeasuredNode(&val_1, n, -1, nil, nil, nil, 0, false)
	var pos *MinimizedCondition_MeasuredNode = &val_1
	return NewMinimizedCondition(pos, StaticGet_MinimizedCondition_UNOPTIMIZED_NEG())
}

func MinimizedCondition_PickBest(a, b *MinimizedCondition_MeasuredNode) *MinimizedCondition_MeasuredNode {
	if a.length == b.length {
		if b.IsChanged() {
			return a
		}
		return b
	}
	if a.length < b.length {
		return a
	}
	return b
}

func MinimizedCondition_ComputeMinimizedCondition(n pkg_rhino.Node) *MinimizedCondition {
	switch n.TokenSwitchOrd_() {
	case 19: // NOT
		sub := MinimizedCondition_ComputeMinimizedCondition(n.GetFirstChild())
		pos := MinimizedCondition_PickBest(MinimizedCondition_MeasuredNode_AddNode_2(n, sub.positive), sub.negative)
		neg := MinimizedCondition_PickBest(sub.negative.priv_MeasuredNode_negate_3874(), sub.positive)
		return NewMinimizedCondition(pos, neg)
	case 69, 68: // AND, OR
		var complementNode pkg_rhino.Node = pkg_rhino.NewNode_1(func() pkg_rhino.Token {
			if n.IsAnd() {
				return pkg_rhino.StaticGet_Token_OR()
			}
			return pkg_rhino.StaticGet_Token_AND()
		}()).Srcref(n)
		_ = complementNode
		left := MinimizedCondition_ComputeMinimizedCondition(n.GetFirstChild())
		right := MinimizedCondition_ComputeMinimizedCondition(n.GetLastChild())
		posTree := MinimizedCondition_MeasuredNode_AddNode_3(n, left.positive, right.positive)
		compTree := MinimizedCondition_MeasuredNode_AddNode_3(complementNode, left.negative, right.negative)
		pos := MinimizedCondition_PickBest(posTree, compTree.priv_MeasuredNode_negate_3874())
		neg := MinimizedCondition_PickBest(posTree.priv_MeasuredNode_negate_3874(), compTree.priv_MeasuredNode_change_3877())
		return NewMinimizedCondition(pos, neg)
	default:
		pos := MinimizedCondition_MeasuredNode_ForNode(n)
		neg := pos.priv_MeasuredNode_negate_3874()
		return NewMinimizedCondition(pos, neg)
	}
}

type MinimizedCondition_MeasuredNode struct {
	c0         *MinimizedCondition_MeasuredNode
	c1         *MinimizedCondition_MeasuredNode
	c2         *MinimizedCondition_MeasuredNode
	node       pkg_rhino.Node
	length     int32
	childCount int32
	idHash_    atomic.Int32
	changed    bool
}

func clinit_MinimizedCondition_MeasuredNode() {}
func StaticGet_MinimizedCondition_MeasuredNode_NOT_PRECEDENCE() int32 {
	return NodeUtil_Precedence(pkg_rhino.StaticGet_Token_NOT())
}

func NewMinimizedCondition_MeasuredNode(n pkg_rhino.Node, childCount int32, c0, c1, c2 *MinimizedCondition_MeasuredNode, len_ int32, ch bool) *MinimizedCondition_MeasuredNode {
	clinit_MinimizedCondition_MeasuredNode()
	this_ := &MinimizedCondition_MeasuredNode{}
	InitMinimizedCondition_MeasuredNode(this_, n, childCount, c0, c1, c2, len_, ch)
	return this_
}
func InitMinimizedCondition_MeasuredNode(this_ *MinimizedCondition_MeasuredNode, n pkg_rhino.Node, childCount int32, c0, c1, c2 *MinimizedCondition_MeasuredNode, len_ int32, ch bool) {
	this_.node = n
	this_.childCount = childCount
	this_.c0 = c0
	this_.c1 = c1
	this_.c2 = c2
	this_.length = len_
	this_.changed = ch
}

func (this_ *MinimizedCondition_MeasuredNode) IsChanged() bool { return this_.changed }
func (this_ *MinimizedCondition_MeasuredNode) IsNot() bool     { return this_.node.IsNot() }
func (this_ *MinimizedCondition_MeasuredNode) IsLowerPrecedenceThan(precedence int32) bool {
	return NodeUtil_Precedence(this_.node.GetToken()) < precedence
}

func (this_ *MinimizedCondition_MeasuredNode) WithoutNot() *MinimizedCondition_MeasuredNode {
	j2gort.CheckState(this_.IsNot())
	var first *MinimizedCondition_MeasuredNode
	if this_.childCount >= 0 {
		first = this_.c0
	} else {
		first = MinimizedCondition_MeasuredNode_ForNode(this_.node.GetFirstChild())
	}
	return first.priv_MeasuredNode_change_3877()
}

func (this_ *MinimizedCondition_MeasuredNode) priv_MeasuredNode_negate_3874() *MinimizedCondition_MeasuredNode {
	switch this_.node.TokenSwitchOrd_() {
	case 4:
		return this_.priv_MeasuredNode_updateToken_3875(pkg_rhino.StaticGet_Token_NE())
	case 5:
		return this_.priv_MeasuredNode_updateToken_3875(pkg_rhino.StaticGet_Token_EQ())
	case 19:
		return this_.WithoutNot()
	default:
		return this_.priv_MeasuredNode_addNot_3876()
	}
}

func (this_ *MinimizedCondition_MeasuredNode) priv_MeasuredNode_updateToken_3875(token pkg_rhino.Token) *MinimizedCondition_MeasuredNode {
	var newNode pkg_rhino.Node = pkg_rhino.NewNode_1(token).Srcref(this_.node)
	_ = newNode
	if this_.childCount >= 0 {
		return NewMinimizedCondition_MeasuredNode(newNode, this_.childCount, this_.c0, this_.c1, this_.c2, this_.length, true)
	}
	var first pkg_rhino.Node = this_.node.GetFirstChild()
	if first == 0 {
		return NewMinimizedCondition_MeasuredNode(newNode, -1, nil, nil, nil, this_.length, true)
	}
	var second pkg_rhino.Node = first.GetNext()
	return NewMinimizedCondition_MeasuredNode(newNode, 2, MinimizedCondition_MeasuredNode_ForNode(first), MinimizedCondition_MeasuredNode_ForNode(second), nil, this_.length, true)
}

func (this_ *MinimizedCondition_MeasuredNode) priv_MeasuredNode_addNot_3876() *MinimizedCondition_MeasuredNode {
	var cost int32 = this_.length + 1
	if this_.IsLowerPrecedenceThan(StaticGet_MinimizedCondition_MeasuredNode_NOT_PRECEDENCE()) {
		cost += 2
	}
	return NewMinimizedCondition_MeasuredNode(pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_NOT()).Srcref(this_.node), 1, this_, nil, nil, cost, true)
}

func (this_ *MinimizedCondition_MeasuredNode) priv_MeasuredNode_change_3877() *MinimizedCondition_MeasuredNode {
	if this_.IsChanged() {
		return this_
	}
	return NewMinimizedCondition_MeasuredNode(this_.node, this_.childCount, this_.c0, this_.c1, this_.c2, this_.length, true)
}

func MinimizedCondition_MeasuredNode_AddNode_2(parent pkg_rhino.Node, child0 *MinimizedCondition_MeasuredNode) *MinimizedCondition_MeasuredNode {
	cost := child0.length
	if parent.IsNot() {
		cost++
	}
	if child0.IsLowerPrecedenceThan(NodeUtil_Precedence(parent.GetToken())) {
		cost += 2
	}
	return NewMinimizedCondition_MeasuredNode(parent, 1, child0, nil, nil, cost, child0.changed)
}

func MinimizedCondition_MeasuredNode_AddNode_3(parent pkg_rhino.Node, child0, child1 *MinimizedCondition_MeasuredNode) *MinimizedCondition_MeasuredNode {
	cost := child0.length + child1.length
	if parent.IsNot() {
		cost++
	}
	pPrec := NodeUtil_Precedence(parent.GetToken())
	if child0.IsLowerPrecedenceThan(pPrec) {
		cost += 2
	}
	if child1.IsLowerPrecedenceThan(pPrec) {
		cost += 2
	}
	return NewMinimizedCondition_MeasuredNode(parent, 2, child0, child1, nil, cost, child0.changed || child1.changed)
}

func MinimizedCondition_MeasuredNode_AddNode_4(parent pkg_rhino.Node, c0, c1, c2 *MinimizedCondition_MeasuredNode) *MinimizedCondition_MeasuredNode {
	return NewMinimizedCondition_MeasuredNode(parent, 3, c0, c1, c2, c0.length+c1.length+c2.length, c0.changed || c1.changed || c2.changed)
}

func MinimizedCondition_MeasuredNode_ForNode(n pkg_rhino.Node) *MinimizedCondition_MeasuredNode {
	return NewMinimizedCondition_MeasuredNode(n, -1, nil, nil, nil, 0, false)
}

func (this_ *MinimizedCondition_MeasuredNode) WillChange(original pkg_rhino.Node) bool {
	j2gort.CheckNotNullNode(original)
	return original != this_.node || this_.IsChanged()
}

func (this_ *MinimizedCondition_MeasuredNode) BuildReplacement() pkg_rhino.Node {
	if this_.childCount >= 0 {
		this_.node.DetachChildren()
		if this_.childCount >= 1 {
			this_.node.AddChildToBack(this_.c0.BuildReplacement().Detach())
		}
		if this_.childCount >= 2 {
			this_.node.AddChildToBack(this_.c1.BuildReplacement().Detach())
		}
	}
	return this_.node
}

func (this_ *MinimizedCondition_MeasuredNode) ApplyTo(original pkg_rhino.Node) pkg_rhino.Node {
	repl := this_.BuildReplacement()
	if original != repl {
		repl.Detach()
		original.ReplaceWith(repl)
	}
	return repl
}
`

	pmcSrc := `package jscomp

import pkg_rhino "example.com/w41b_task3/gen/rhino"

type AbstractPeepholeOptimizationBase struct {
	compiler AbstractCompiler
}
func (b *AbstractPeepholeOptimizationBase) EndTraversal() { b.compiler = nil }

type PeepholeMinimizeConditions struct {
	AbstractPeepholeOptimizationBase
}

func (this_ *PeepholeMinimizeConditions) MinimizeIf(ifNode pkg_rhino.Node) pkg_rhino.Node {
	cond := ifNode.GetFirstChild()
	if !MinimizedCondition_IsMinimizable(cond) {
		return cond
	}
	minCond := MinimizedCondition_FromConditionNode(cond)
	mNode := minCond.GetMinimized(StaticGet_MinimizedCondition_MinimizationStyle_ALLOW_LEADING_NOT())
	if mNode.WillChange(cond) {
		return mNode.ApplyTo(cond)
	}
	return cond
}
`

	rtTestSrc := `package jscomp

import (
	"reflect"
	"testing"

	pkg_rhino "example.com/w41b_task3/gen/rhino"
)

func TestScratchVarExpertAndLazyMeasuredNodeRuntime(t *testing.T) {
	// 1. Verify MinimizedCondition_MeasuredNode fits in <= 48 bytes!
	if sz := reflect.TypeFor[MinimizedCondition_MeasuredNode]().Size(); sz > 48 {
		t.Fatalf("expected sizeof(MinimizedCondition_MeasuredNode) <= 48, got %d", sz)
	}

	// 2. Verify InlineVariables scratchVarExpert_ allocates 0 VarExperts on slab for non-escaping vars,
	//    and promotes ONLY LazyNegative (vLazy) + AliasInlineRetryHandler (vWait)!
	vTarget := &Var{name: "targetLazy", kind: 1}
	vWait := &Var{name: "waitAlias", kind: 2, alias: vTarget}
	vars := []*Var{vWait}
	cols := map[*Var]*ReferenceCollection{
		vWait:   {refs: []*Reference{{node: 1}}},
		vTarget: {refs: []*Reference{{node: 2}}}, // Size()==1 -> IsWellDefinedAssignedOnce()==true!
	}
	for i := 0; i < 50; i++ {
		vi := &Var{name: "localConst", kind: 0}
		vars = append(vars, vi)
		cols[vi] = &ReferenceCollection{refs: []*Reference{{node: 3}}}
	}
	vars = append(vars, vTarget)

	iv := &InlineVariables{vars: vars, cols: cols}
	iv.Process_2(1, 2)
	if len(iv.inlinedNames) != 51 {
		t.Fatalf("expected 50 localConsts + 1 retried waitAlias inlined (51 total), got %d", len(iv.inlinedNames))
	}
	if got := iv.slab_InlineVariables_VarExpert.AllocCalls.Load(); got != 2 {
		t.Fatalf("expected ONLY 2 slab_InlineVariables_VarExpert allocs (1 RetryHandler + 1 LazyNegative) across 52 vars, got %d", got)
	}

	// 3. Verify MinimizedCondition performs ZERO NewNode_1 calls when conditions are already minimal (a && b)!
	pmc := &PeepholeMinimizeConditions{}
	defer pmc.EndTraversal()

	pkg_rhino.ResetNodes()
	a := pkg_rhino.AllocTestNode(pkg_rhino.TokNAME)
	b := pkg_rhino.AllocTestNode(pkg_rhino.TokNAME)
	andCond := pkg_rhino.AllocTestNode(pkg_rhino.TokAND, a, b)
	ifAlreadyMin := pkg_rhino.AllocTestNode(pkg_rhino.TokIF, andCond)

	pkg_rhino.NewNode1Calls.Store(0)
	gotCond := pmc.MinimizeIf(ifAlreadyMin)
	if gotCond != andCond {
		t.Fatalf("expected already-minimal a && b to remain unchanged")
	}
	if calls := pkg_rhino.NewNode1Calls.Load(); calls != 0 {
		t.Fatalf("expected 0 speculative NewNode_1 calls on unchanged condition, got %d", calls)
	}

	// 4. Verify De Morgan minimization !(a != b && !c) -> a == b || c materializes
	//    ONLY the winning replacement nodes with exact Srcref preservation!
	x := pkg_rhino.AllocTestNode(pkg_rhino.TokNAME)
	y := pkg_rhino.AllocTestNode(pkg_rhino.TokNAME)
	neNode := pkg_rhino.AllocTestNode(pkg_rhino.TokNE, x, y)
	z := pkg_rhino.AllocTestNode(pkg_rhino.TokNAME)
	notZ := pkg_rhino.AllocTestNode(pkg_rhino.TokNOT, z)
	and2 := pkg_rhino.AllocTestNode(pkg_rhino.TokAND, neNode, notZ)
	outerNot := pkg_rhino.AllocTestNode(pkg_rhino.TokNOT, and2)
	ifDeMorgan := pkg_rhino.AllocTestNode(pkg_rhino.TokIF, outerNot)

	pkg_rhino.NewNode1Calls.Store(0)
	repl := pmc.MinimizeIf(ifDeMorgan)
	if repl.GetToken() != pkg_rhino.TokOR {
		t.Fatalf("expected !(x != y && !z) to minimize to TokOR (x == y || z), got tok=%v", repl.GetToken())
	}
	if repl.SrcRefOf() != and2 {
		t.Fatalf("expected materialized OR node to preserve Srcref(and2=%d), got %d", and2, repl.SrcRefOf())
	}
	eqChild := repl.GetFirstChild()
	if eqChild.GetToken() != pkg_rhino.TokEQ || eqChild.SrcRefOf() != neNode {
		t.Fatalf("expected first child to be materialized TokEQ with Srcref(neNode), got tok=%v srcref=%d", eqChild.GetToken(), eqChild.SrcRefOf())
	}
	if calls := pkg_rhino.NewNode1Calls.Load(); calls != 2 {
		t.Fatalf("expected ONLY 2 materialized NewNode_1 calls (winning OR + winning EQ), got %d", calls)
	}
}
`

	wave41bWriteFiles(t, tmpDir, map[string]string{
		"go.mod":                            "module example.com/w41b_task3\n\ngo 1.24\n",
		"gen/j2gort/j2gort.go":              j2Src,
		"gen/rhino/rhino.go":                rhinoSrc,
		"gen/jscomp/inline_variables.go":    inlineVarsSrc,
		"gen/jscomp/minimized_condition.go": minCondSrc,
		"gen/jscomp/peephole_min.go":        pmcSrc,
	})

	jscompDir := filepath.Join(tmpDir, "gen/jscomp")
	res1, err := passes.RunOnPackageDir(jscompDir, passes.AllPasses(), passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if res1.TotalSites == 0 {
		t.Fatalf("expected non-zero sites on Run 1, got %+v", res1)
	}

	// Run 2 MUST be 100% idempotent (0 sites, 0 modified files).
	res2, err := passes.RunOnPackageDir(jscompDir, passes.AllPasses(), passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if res2.TotalSites != 0 || res2.FilesModified != 0 {
		t.Fatalf("Run 2 NOT IDEMPOTENT: %+v", res2)
	}

	mcBytes, err := os.ReadFile(filepath.Join(jscompDir, "minimized_condition.go"))
	if err != nil {
		t.Fatal(err)
	}
	mcCode := string(mcBytes)
	addNotIdx := strings.Index(mcCode, "func (this_ *MinimizedCondition_MeasuredNode) addNotWithSlab_")
	changeIdx := strings.Index(mcCode, "func (this_ *MinimizedCondition_MeasuredNode) changeWithSlab_")
	if addNotIdx < 0 || changeIdx < addNotIdx {
		t.Fatalf("missing addNotWithSlab_ in lowered minimized_condition.go:\n%s", mcCode)
	}
	if strings.Contains(mcCode[addNotIdx:changeIdx], "NewNode_1") {
		t.Fatalf("addNotWithSlab_ MUST NOT call NewNode_1:\n%s", mcCode[addNotIdx:changeIdx])
	}

	wave41bWriteFiles(t, tmpDir, map[string]string{
		"gen/jscomp/task3_rt_test.go": rtTestSrc,
	})

	cmd := exec.Command("go", "test", "-race", "-v", "./gen/jscomp")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test -race ./gen/jscomp failed: %v\n%s", err, string(out))
	}
}

// TestWave41B_RealGenJscompFullPassesBuildAndIdempotency copies real gen/jscomp into an isolated
// module, runs AllPasses() Run 1, verifies all Wave 41B rewrites landed cleanly, compiles ./jscomp,
// and asserts Run 2 produces TotalSites == 0 and FilesModified == 0.
func TestWave41B_RealGenJscompFullPassesBuildAndIdempotency(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping full real gen/jscomp build + idempotency check in -short mode")
	}
	t.Parallel()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	repoRoot := filepath.Clean(filepath.Join(wd, "..", "..", ".."))
	if !realGenJ2gortReady(repoRoot) {
		t.Skip("gen/rhino or gen/j2gort is mid-edit; skipping real gen/jscomp copy check")
	}

	realJscompDir := filepath.Join(repoRoot, "transpiled", "gen", "jscomp")
	entries, err := os.ReadDir(realJscompDir)
	if err != nil {
		t.Fatal(err)
	}

	tmpDir := t.TempDir()
	origModBytes, err := os.ReadFile(filepath.Join(repoRoot, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	goMod := strings.Replace(string(origModBytes), "module github.com/paullepoulpe/jsgo", "module example.com/realw41b", 1) +
		fmt.Sprintf("\nrequire github.com/paullepoulpe/jsgo v0.0.0\nreplace github.com/paullepoulpe/jsgo => %s\n", filepath.ToSlash(repoRoot))
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}
	if sumBytes, sErr := os.ReadFile(filepath.Join(repoRoot, "go.sum")); sErr == nil {
		_ = os.WriteFile(filepath.Join(tmpDir, "go.sum"), sumBytes, 0o644)
	}

	copyDir := filepath.Join(tmpDir, "jscomp")
	if err := os.MkdirAll(copyDir, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, ent := range entries {
		if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".go") || strings.HasSuffix(ent.Name(), "_test.go") {
			continue
		}
		data, rErr := os.ReadFile(filepath.Join(realJscompDir, ent.Name()))
		if rErr != nil {
			t.Fatal(rErr)
		}
		if wErr := os.WriteFile(filepath.Join(copyDir, ent.Name()), data, 0o644); wErr != nil {
			t.Fatal(wErr)
		}
	}

	allPasses := passes.AllPasses()
	s1, err := passes.RunOnPackageDir(copyDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("AllPasses() Run 1 on copied real gen/jscomp failed: %v", err)
	}
	t.Logf("Wave 41B Real gen/jscomp Run 1 sites: total=%d files=%d byPass=%+v", s1.TotalSites, s1.FilesModified, s1.ByPass)

	s2, err := passes.RunOnPackageDir(copyDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("AllPasses() Run 2 on copied real gen/jscomp failed: %v", err)
	}
	if s2.TotalSites != 0 || s2.FilesModified != 0 {
		t.Fatalf("Wave 41B real gen/jscomp Run 2 NOT IDEMPOTENT: total=%d files=%d byPass=%+v", s2.TotalSites, s2.FilesModified, s2.ByPass)
	}

	readCopy := func(name string) string {
		b, rErr := os.ReadFile(filepath.Join(copyDir, name))
		if rErr != nil {
			t.Fatal(rErr)
		}
		return string(b)
	}

	// Verify Task 1: StringPool.internAt_ uses InternOwned.
	if sp := readCopy("string_pool.go"); !strings.Contains(sp, "InternOwned(") || strings.Contains(sp, "alloc_.Intern(") {
		t.Fatalf("expected lowered real string_pool.go to use InternOwned without alloc_.Intern")
	}
	// Verify Task 2: TypeValidator.checkWellFormedTemplatizedType IIFE unwrapped.
	tv := readCopy("type_validator.go")
	if strings.Contains(tv, "var ok bool = func() bool") {
		t.Fatalf("expected TypeValidator.checkWellFormedTemplatizedType closure IIFE to be unwrapped")
	}
	// Verify Task 3: InlineVariables scratchVarExpert_ & MinimizedCondition lazyTok_.
	if iv := readCopy("inline_variables.go"); !strings.Contains(iv, "scratchVarExpert_") || !strings.Contains(iv, "promoteIvScratchVarExpert_") {
		t.Fatalf("expected lowered real inline_variables.go to contain scratchVarExpert_ and promoteIvScratchVarExpert_")
	}
	if mc := readCopy("minimized_condition.go"); !strings.Contains(mc, "lazyTok_") || !strings.Contains(mc, "allocLazy_MinimizedCondition_MeasuredNode") {
		t.Fatalf("expected lowered real minimized_condition.go to contain lazyTok_ and allocLazy_MinimizedCondition_MeasuredNode")
	}

	buildCmd := exec.Command("go", "build", "./jscomp")
	buildCmd.Dir = tmpDir
	if out, bErr := buildCmd.CombinedOutput(); bErr != nil {
		t.Fatalf("go build ./jscomp on lowered real gen/jscomp failed: %v\n%s", bErr, string(out))
	}
}
