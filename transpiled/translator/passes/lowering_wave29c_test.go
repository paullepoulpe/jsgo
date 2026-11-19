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

func TestWave29C_GetAncestors_LazyPreconditions_LazyLogFile_AndRucVarInfoPools(t *testing.T) {
	tmpDir := t.TempDir()
	mustWrite := func(rel, content string) {
		t.Helper()
		full := filepath.Join(tmpDir, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	mustWrite("go.mod", "module example.com/w29ctest\n\ngo 1.24\n")

	// 1. Minimal j2gort runtime package matching production contracts + pool counters.
	mustWrite("gen/j2gort/runtime.go", `package j2gort

import (
	"fmt"
	"math/bits"
	"reflect"
	"strconv"
	"sync"
	"sync/atomic"
)

const NullString = "\xff\x00j2go_null"

type NPE struct{}

func (NPE) Error() string { return "java.lang.NullPointerException" }

type IllegalStateException struct{ Msg string }

func (e IllegalStateException) Error() string { return "IllegalStateException: " + e.Msg }

type IllegalArgumentException struct{ Msg string }

func (e IllegalArgumentException) Error() string { return "IllegalArgumentException: " + e.Msg }

func ThrowNPE() {
	panic(NPE{})
}

func CheckNotNullPtr[T any](p *T) *T {
	if p == nil {
		ThrowNPE()
	}
	return p
}

func CheckNotNullSlice[T any](s []T) []T {
	if s == nil {
		ThrowNPE()
	}
	return s
}

func CheckNotNullNode[N ~uint32](n N) N {
	if n == 0 {
		ThrowNPE()
	}
	return n
}

func CheckNotNullEnum[E ~uint8](v E) E {
	if v == 0 {
		ThrowNPE()
	}
	return v
}

func CheckNotNullString(s string) string {
	if s == NullString {
		ThrowNPE()
	}
	return s
}

func StrValString(s string) string {
	if s == NullString {
		return "null"
	}
	return s
}

func StrValInt32(v int32) string {
	return strconv.FormatInt(int64(v), 10)
}

func StrVal(v any) string {
	if v == nil {
		return "null"
	}
	switch x := v.(type) {
	case string:
		if x == NullString {
			return "null"
		}
		return x
	case *BoxInt:
		if x == nil {
			return "null"
		}
		return strconv.FormatInt(int64(x.Val), 10)
	case fmt.Stringer:
		return x.String()
	default:
		return fmt.Sprint(v)
	}
}

func BoxString(s string) any {
	if s == NullString {
		return nil
	}
	// Copy to force a measurable heap allocation when boxed into any if non-empty & > 1 byte.
	b := []byte(s)
	return string(b)
}

func BoxNode[N ~uint32](n N) any {
	if n == 0 {
		return nil
	}
	return n
}

func UnboxNode[N ~uint32](v any) N {
	if v == nil {
		ThrowNPE()
	}
	return v.(N)
}

func CastNode[N ~uint32](v any) N {
	if v == nil {
		return 0
	}
	return v.(N)
}

func PtrToAny[T any](p *T) any {
	if p == nil {
		return nil
	}
	return p
}

func ToAny(v any) any {
	return v
}

func BoxNullableVal[T comparable](v T) any {
	var z T
	if v == z {
		return nil
	}
	return v
}

type BoxInt struct{ Val int32 }

func (b *BoxInt) String() string { return strconv.FormatInt(int64(b.Val), 10) }

func ValueOfInt(v int32) *BoxInt {
	return &BoxInt{Val: v}
}

type BoxLong struct{ Val int64 }

func ValueOfLong(v int64) *BoxLong {
	return &BoxLong{Val: v}
}

type BoxBoolean struct{ Val bool }

func ValueOfBoolean(v bool) *BoxBoolean {
	return &BoxBoolean{Val: v}
}

type BoxDouble struct{ Val float64 }

func ValueOfDouble(v float64) *BoxDouble {
	return &BoxDouble{Val: v}
}

func AsIface[I any](v any) I {
	if v == nil {
		var z I
		return z
	}
	return v.(I)
}

func LenientFormat(tmpl string, args ...any) string {
	out := tmpl
	for _, a := range args {
		idx := strings.Index(out, "%s")
		if idx >= 0 {
			out = out[:idx] + StrVal(a) + out[idx+2:]
		} else {
			out += " [" + StrVal(a) + "]"
		}
	}
	return out
}

func CheckState(b bool) {
	if !b {
		panic(IllegalStateException{Msg: ""})
	}
}

func CheckStateMsg(b bool, msg any) {
	if !b {
		panic(IllegalStateException{Msg: StrVal(msg)})
	}
}

func CheckStateFmt(b bool, tmpl string, args ...any) {
	if !b {
		panic(IllegalStateException{Msg: LenientFormat(tmpl, args...)})
	}
}

func CheckArgumentMsg(b bool, msg any) {
	if !b {
		panic(IllegalArgumentException{Msg: StrVal(msg)})
	}
}

func CheckArgumentFmt(b bool, tmpl string, args ...any) {
	if !b {
		panic(IllegalArgumentException{Msg: LenientFormat(tmpl, args...)})
	}
}

func Ext_java_lang_String_Format_2(tmpl string, args []any) string {
	return LenientFormat(tmpl, args...)
}

type Ext_java_util_Iterator interface {
	HasNext() bool
	Next() any
}

type Ext_java_lang_Iterable interface {
	Iterator() any
}

type JIterCursor struct {
	it Ext_java_util_Iterator
}

func NewIterCursor(src any) JIterCursor {
	if src == nil {
		ThrowNPE()
	}
	ib := src.(Ext_java_lang_Iterable)
	return JIterCursor{it: ib.Iterator().(Ext_java_util_Iterator)}
}

func (c *JIterCursor) HasNext() bool {
	return c.it.HasNext()
}

func (c *JIterCursor) Next() any {
	return c.it.Next()
}

func (c *JIterCursor) NextNode() uint32 {
	v := c.it.Next()
	if v == nil {
		ThrowNPE()
	}
	rv := reflect.ValueOf(v)
	return uint32(rv.Uint())
}

type Ext_java_util_function_Supplier interface {
	Implements_j2gort_Ext_java_util_function_Supplier()
	IsJavaNil_() bool
	Get() any
}

func ToIface_Ext_java_util_function_Supplier[T interface {
	Ext_java_util_function_Supplier
	comparable
}](v T) Ext_java_util_function_Supplier {
	var z T
	if v == z {
		return nil
	}
	return v
}

type lambdaSupplier struct {
	fn func() any
}

func (*lambdaSupplier) Implements_j2gort_Ext_java_util_function_Supplier() {}
func (l *lambdaSupplier) IsJavaNil_() bool                                 { return l == nil }
func (l *lambdaSupplier) Get() any                                         { return l.fn() }

func NewLambda_Ext_java_util_function_Supplier(fn func() any) Ext_java_util_function_Supplier {
	return &lambdaSupplier{fn: fn}
}

// Slab + PooledSlab + PooledSlice infrastructure with live leak counters.
type Slab[T any] struct {
	chunks [][]T
	idx    int
}

func (s *Slab[T]) Alloc() *T {
	if len(s.chunks) == 0 || s.idx >= len(s.chunks[len(s.chunks)-1]) {
		s.chunks = append(s.chunks, make([]T, 64))
		s.idx = 0
	}
	c := s.chunks[len(s.chunks)-1]
	p := &c[s.idx]
	s.idx++
	return p
}

func (s *Slab[T]) Release() { s.chunks = nil; s.idx = 0 }

type PooledSlab[T any] struct {
	Slab[T]
}

func (s *PooledSlab[T]) Release() { s.Slab.Release() }

const (
	pooledSliceMinShift = 3
	pooledSliceMinCap   = 1 << pooledSliceMinShift // 8
	pooledSliceMaxShift = 18
	pooledSliceMaxCap   = 1 << pooledSliceMaxShift // 262144
	pooledSliceTiers    = pooledSliceMaxShift - pooledSliceMinShift + 1
)

type pooledSliceTier[T any] struct {
	mu   sync.Mutex
	free [][]T
}

type pooledSliceTypePool[T any] struct {
	tiers [pooledSliceTiers]pooledSliceTier[T]
}

var (
	pooledSlicePools      sync.Map
	PooledSliceAcquireCnt atomic.Int64
	PooledSliceReleaseCnt atomic.Int64
	PooledSliceLiveCnt    atomic.Int64
	PooledSliceReuseCnt   atomic.Int64
)

func ResetSlicePoolCounters() {
	PooledSliceAcquireCnt.Store(0)
	PooledSliceReleaseCnt.Store(0)
	PooledSliceLiveCnt.Store(0)
	PooledSliceReuseCnt.Store(0)
}

func getPooledSlicePool[T any]() *pooledSliceTypePool[T] {
	rt := reflect.TypeFor[T]()
	if v, ok := pooledSlicePools.Load(rt); ok {
		return v.(*pooledSliceTypePool[T])
	}
	p := &pooledSliceTypePool[T]{}
	actual, _ := pooledSlicePools.LoadOrStore(rt, p)
	return actual.(*pooledSliceTypePool[T])
}

func AcquirePooledSlice[T any](capHint int) []T {
	if capHint <= 0 {
		return nil
	}
	PooledSliceAcquireCnt.Add(1)
	if capHint > pooledSliceMaxCap {
		return make([]T, 0, capHint)
	}
	PooledSliceLiveCnt.Add(1)
	tierIdx := 0
	allocCap := pooledSliceMinCap
	if capHint > pooledSliceMinCap {
		shift := bits.Len(uint(capHint - 1))
		tierIdx = shift - pooledSliceMinShift
		allocCap = 1 << shift
	}
	t := &getPooledSlicePool[T]().tiers[tierIdx]
	t.mu.Lock()
	if n := len(t.free); n > 0 {
		s := t.free[n-1]
		t.free[n-1] = nil
		t.free = t.free[:n-1]
		t.mu.Unlock()
		PooledSliceReuseCnt.Add(1)
		return s[:0]
	}
	t.mu.Unlock()
	return make([]T, 0, allocCap)
}

func AppendPooledSlice[T any](s *[]T, v T) {
	cur := *s
	if len(cur) < cap(cur) {
		*s = append(cur, v)
		return
	}
	oldLen := len(cur)
	newCap := max(pooledSliceMinCap, cap(cur)*2)
	next := AcquirePooledSlice[T](newCap)[:oldLen+1]
	copy(next, cur)
	next[oldLen] = v
	if cap(cur) >= pooledSliceMinCap {
		ReleasePooledSlice(cur)
	}
	*s = next
}

func ReleasePooledSlice[T any](s []T) {
	c := cap(s)
	if c < pooledSliceMinCap || c > pooledSliceMaxCap {
		return
	}
	PooledSliceReleaseCnt.Add(1)
	PooledSliceLiveCnt.Add(-1)
	clear(s[:c])
	tierIdx := bits.Len(uint(c)) - 1 - pooledSliceMinShift
	t := &getPooledSlicePool[T]().tiers[tierIdx]
	t.mu.Lock()
	if len(t.free) < 64 {
		t.free = append(t.free, s[:0])
	}
	t.mu.Unlock()
}
`)
	// Fix missing strings import in gen/j2gort/runtime.go
	j2File := filepath.Join(tmpDir, "gen/j2gort/runtime.go")
	j2Bytes, _ := os.ReadFile(j2File)
	mustWrite("gen/j2gort/runtime.go", strings.Replace(string(j2Bytes), `"strconv"`, "\"strconv\"\n\t\"strings\"", 1))

	// 2. Minimal gen/diagnostic package defining LogFile + NoOpLogFile + WritingLogFile.
	mustWrite("gen/diagnostic/log_file.go", `package diagnostic

import (
	"strings"
	"sync"

	"example.com/w29ctest/gen/j2gort"
)

type LogFile_StreamedJsonProducer interface {
	WriteJson() string
}

type LogFile interface {
	Implements_diagnostic_LogFile()
	IsJavaNil_() bool
	Log_1_any(any) LogFile
	Log_1_string(string) LogFile
	Log_1_Ext_java_util_function_Supplier(j2gort.Ext_java_util_function_Supplier) LogFile
	Log_2(string, []any) LogFile
	LogJson_1_any(any) LogFile
	LogJson_1_Ext_java_util_function_Supplier(j2gort.Ext_java_util_function_Supplier) LogFile
	LogJson_1_LogFile_StreamedJsonProducer(LogFile_StreamedJsonProducer) LogFile
	IsLogging() bool
}

type NoOpLogFile struct{}

func (*NoOpLogFile) Implements_diagnostic_LogFile()                                                 {}
func (n *NoOpLogFile) IsJavaNil_() bool                                                             { return n == nil }
func (n *NoOpLogFile) IsLogging() bool                                                              { return false }
func (n *NoOpLogFile) Log_1_any(any) LogFile                                                        { return n }
func (n *NoOpLogFile) Log_1_string(string) LogFile                                                  { return n }
func (n *NoOpLogFile) Log_1_Ext_java_util_function_Supplier(j2gort.Ext_java_util_function_Supplier) LogFile { return n }
func (n *NoOpLogFile) Log_2(string, []any) LogFile                                                  { return n }
func (n *NoOpLogFile) LogJson_1_any(any) LogFile                                                    { return n }
func (n *NoOpLogFile) LogJson_1_Ext_java_util_function_Supplier(j2gort.Ext_java_util_function_Supplier) LogFile {
	return n
}
func (n *NoOpLogFile) LogJson_1_LogFile_StreamedJsonProducer(LogFile_StreamedJsonProducer) LogFile { return n }

type RecordingLogFile struct {
	mu      sync.Mutex
	Entries []string
}

func (*RecordingLogFile) Implements_diagnostic_LogFile() {}
func (r *RecordingLogFile) IsJavaNil_() bool             { return r == nil }
func (r *RecordingLogFile) IsLogging() bool              { return true }
func (r *RecordingLogFile) appendEntry(s string) LogFile {
	r.mu.Lock()
	r.Entries = append(r.Entries, s)
	r.mu.Unlock()
	return r
}
func (r *RecordingLogFile) Log_1_any(v any) LogFile {
	return r.appendEntry(j2gort.StrVal(v))
}
func (r *RecordingLogFile) Log_1_string(s string) LogFile {
	return r.appendEntry(s)
}
func (r *RecordingLogFile) Log_1_Ext_java_util_function_Supplier(sup j2gort.Ext_java_util_function_Supplier) LogFile {
	return r.appendEntry(j2gort.StrVal(sup.Get()))
}
func (r *RecordingLogFile) Log_2(tmpl string, vals []any) LogFile {
	return r.appendEntry(j2gort.Ext_java_lang_String_Format_2(tmpl, vals))
}
func (r *RecordingLogFile) LogJson_1_any(v any) LogFile {
	return r.appendEntry("JSON:" + j2gort.StrVal(v))
}
func (r *RecordingLogFile) LogJson_1_Ext_java_util_function_Supplier(sup j2gort.Ext_java_util_function_Supplier) LogFile {
	return r.appendEntry("JSON:" + j2gort.StrVal(sup.Get()))
}
func (r *RecordingLogFile) LogJson_1_LogFile_StreamedJsonProducer(p LogFile_StreamedJsonProducer) LogFile {
	return r.appendEntry("JSON:" + strings.TrimSpace(p.WriteJson()))
}
`)

	// 3. Minimal gen/rhino package defining Node (uint32), Token (uint8), AncestorIterable, GetAncestors(), Detach(), etc.
	mustWrite("gen/rhino/node.go", `package rhino

import (
	"fmt"
	"sync"

	"example.com/w29ctest/gen/j2gort"
)

type Token uint8

const (
	TokNAME     Token = 1
	TokFUNCTION Token = 2
	TokFOR      Token = 3
	TokBLOCK    Token = 4
	TokSCRIPT   Token = 5
	TokCALL     Token = 6
)

var tokNames = [...]string{
	TokNAME:     "NAME",
	TokFUNCTION: "FUNCTION",
	TokFOR:      "FOR",
	TokBLOCK:    "BLOCK",
	TokSCRIPT:   "SCRIPT",
	TokCALL:     "CALL",
}

func (t Token) String() string {
	if t == 0 {
		return "null"
	}
	return tokNames[t]
}

func (t Token) ToString() string {
	if t == 0 {
		j2gort.ThrowNPE()
	}
	return tokNames[t]
}

func (t Token) Name() string {
	if t == 0 {
		j2gort.ThrowNPE()
	}
	return tokNames[t]
}

type Node uint32

type nodeRec struct {
	tok    Token
	label  string
	parent Node
	first  Node
	next   Node
}

var (
	arenaMu sync.RWMutex
	arena   = []nodeRec{{}} // index 0 = null Node
)

func NewNode(tok Token, label string, children ...Node) Node {
	arenaMu.Lock()
	id := Node(len(arena))
	arena = append(arena, nodeRec{tok: tok, label: label})
	for i, ch := range children {
		arena[ch].parent = id
		if i == 0 {
			arena[id].first = ch
		} else {
			arena[children[i-1]].next = ch
		}
	}
	arenaMu.Unlock()
	return id
}

func (n Node) checkNonZero() {
	if n == 0 {
		j2gort.ThrowNPE()
	}
}

func (n Node) GetParent() Node {
	n.checkNonZero()
	arenaMu.RLock()
	p := arena[n].parent
	arenaMu.RUnlock()
	return p
}

func (n Node) Detach() Node {
	n.checkNonZero()
	arenaMu.Lock()
	p := arena[n].parent
	if p != 0 {
		if arena[p].first == n {
			arena[p].first = arena[n].next
		} else {
			for prev := arena[p].first; prev != 0; prev = arena[prev].next {
				if arena[prev].next == n {
					arena[prev].next = arena[n].next
					break
				}
			}
		}
		arena[n].parent = 0
		arena[n].next = 0
	}
	arenaMu.Unlock()
	return n
}

func (n Node) GetToken() Token {
	n.checkNonZero()
	arenaMu.RLock()
	t := arena[n].tok
	arenaMu.RUnlock()
	return t
}

func (n Node) Label() string {
	n.checkNonZero()
	arenaMu.RLock()
	s := arena[n].label
	arenaMu.RUnlock()
	return s
}

func (n Node) IsName() bool {
	return n.GetToken() == TokNAME
}

func (n Node) IsFunction() bool {
	return n.GetToken() == TokFUNCTION
}

func (n Node) IsFor() bool {
	return n.GetToken() == TokFOR
}

func (n Node) HasChildren() bool {
	n.checkNonZero()
	arenaMu.RLock()
	f := arena[n].first
	arenaMu.RUnlock()
	return f != 0
}

func (n Node) HasOneChild() bool {
	n.checkNonZero()
	arenaMu.RLock()
	f := arena[n].first
	one := f != 0 && arena[f].next == 0
	arenaMu.RUnlock()
	return one
}

func (n Node) ToString() string {
	n.checkNonZero()
	return fmt.Sprintf("%s(%s)", n.GetToken().String(), n.Label())
}

func (n Node) String() string {
	if n == 0 {
		return "null"
	}
	return n.ToString()
}

func (n Node) ToString_3(printSource, printType, printAnnotations bool) string {
	n.checkNonZero()
	return fmt.Sprintf("%s(%s)[%v,%v,%v]", n.GetToken().String(), n.Label(), printSource, printType, printAnnotations)
}

func (n Node) ToStringTree() string {
	n.checkNonZero()
	return "TREE<" + n.ToString() + ">"
}

func (n Node) GetLocation() string {
	n.checkNonZero()
	return "input.js:1"
}

type Node_AncestorIterable struct {
	cur Node
}

func NewNode_AncestorIterable(cur Node) *Node_AncestorIterable {
	return &Node_AncestorIterable{cur: cur}
}

func (a *Node_AncestorIterable) Iterator() any {
	return &nodeAncestorIter{outer_: a}
}

type nodeAncestorIter struct {
	outer_ *Node_AncestorIterable
}

func (it *nodeAncestorIter) HasNext() bool {
	return it.outer_.cur != 0
}

func (it *nodeAncestorIter) Next() any {
	if !it.HasNext() {
		panic("NoSuchElementException")
	}
	res := it.outer_.cur
	it.outer_.cur = it.outer_.cur.GetParent()
	return j2gort.BoxNode(res)
}

// Matches gen/rhino/node.go:3511: Node.getAncestors() starts at this_.GetParent()!
func (this_ Node) GetAncestors() *Node_AncestorIterable {
	return NewNode_AncestorIterable(this_.GetParent())
}
`)

	// 4. gen/jscomp covering Tasks 1, 2, 3, and 4 + negative soundness controls!
	mustWrite("gen/jscomp/ancestors_and_preconditions.go", `package jscomp

import (
	"sync/atomic"

	"example.com/w29ctest/gen/j2gort"
	pkg_rhino "example.com/w29ctest/gen/rhino"
)

// Task 1a: Hot NodeUtil.isWithinLoop pattern (NextNode).
func NodeUtil_IsWithinLoop(n pkg_rhino.Node) bool {
	_ = n
	{
		cur_28 := j2gort.NewIterCursor(j2gort.CheckNotNullPtr(n.GetAncestors()))
		for cur_28.HasNext() {
			parent := pkg_rhino.Node(cur_28.NextNode())
			_ = parent
			if parent.IsFor() {
				return true
			}
			if parent.IsFunction() {
				break
			}
		}
	}
	return false
}

// Task 1b: UnboxNode(cur.Next()) where body executes continue AND detaches anc mid-iteration!
func DetachForAncestorsAndCollectLabels(start pkg_rhino.Node) []string {
	var visited []string
	{
		cur_1 := j2gort.NewIterCursor(start.GetAncestors())
		for cur_1.HasNext() {
			anc := j2gort.UnboxNode[pkg_rhino.Node](cur_1.Next())
			_ = anc
			if anc.IsFor() {
				// Detach anc (sets anc.parent = 0) and continue!
				// Because curAnc_ already advanced to anc.GetParent() BEFORE body,
				// iteration MUST still continue to anc's former parent!
				anc.Detach()
				visited = append(visited, "detached:"+anc.Label())
				continue
			}
			visited = append(visited, anc.Label())
		}
	}
	return visited
}

// Task 1c: ReplaceCssNames local var ancestors = ... pattern + nested GetAncestors loop!
type CssHelper struct {
	callNode pkg_rhino.Node
}

func (this_ *CssHelper) CountAncestorPairs() int {
	var ancestors *pkg_rhino.Node_AncestorIterable = this_.callNode.GetAncestors()
	_ = ancestors
	count := 0
	{
		cur_13 := j2gort.NewIterCursor(j2gort.CheckNotNullPtr(ancestors))
		for cur_13.HasNext() {
			anc := pkg_rhino.Node(cur_13.NextNode())
			_ = anc
			count++
			for cur_14 := j2gort.NewIterCursor(j2gort.CheckNotNullPtr(anc.GetAncestors())); cur_14.HasNext(); {
				raw := cur_14.Next()
				outerAnc := j2gort.UnboxNode[pkg_rhino.Node](raw)
				_ = outerAnc
				count += 10
			}
		}
	}
	return count
}

// Task 1 Negative Control: escaping ancestors variable used after the loop MUST NOT be dropped!
func EscapingAncestorsIterable(n pkg_rhino.Node) *pkg_rhino.Node_AncestorIterable {
	var ancestors *pkg_rhino.Node_AncestorIterable = n.GetAncestors()
	{
		cur := j2gort.NewIterCursor(j2gort.CheckNotNullPtr(ancestors))
		for cur.HasNext() {
			_ = pkg_rhino.Node(cur.NextNode())
			break
		}
	}
	return ancestors
}

// Task 2: Normalize.annotateFunctionExpressionNameAsConstant & hot Preconditions calls!
type MetaHolder struct {
	ok  bool
	tag string
}

func AnnotateFunctionName(n pkg_rhino.Node) {
	j2gort.CheckStateFmt(n.IsName(), "Expected NAME node but got %s", j2gort.BoxString(j2gort.CheckNotNullEnum(n.GetToken()).String()))
}

func CheckTokenToStringAndTree(n pkg_rhino.Node, meta *MetaHolder, idx int32) {
	j2gort.CheckStateFmt(n.IsName(), "Expected NAME but got %s (%s)", j2gort.BoxString(n.GetToken().ToString()), j2gort.BoxString(n.ToString()))
	j2gort.CheckStateFmt(!(n.HasChildren()), "Expected leaf node, found %s", j2gort.BoxString(n.ToStringTree()))
	j2gort.CheckArgumentFmt(meta != nil && meta.ok, "Bad meta tag=%s idx=%s node=%s", j2gort.BoxString(meta.tag), j2gort.ValueOfInt(idx), j2gort.BoxNode(n))
	j2gort.CheckStateMsg(n != 0, j2gort.ValueOfInt(idx))
	j2gort.CheckArgumentMsg(idx >= 0, j2gort.ValueOfInt(idx))
}

var PrecondSideEffectCalls atomic.Int64

func sideEffectArg() int32 {
	return int32(PrecondSideEffectCalls.Add(1))
}

// Task 2 Negative Control: impure argument MUST NOT be guarded by if !(cond)!
func CheckStateImpureArg(cond bool) {
	j2gort.CheckStateFmt(cond, "count=%s", j2gort.ValueOfInt(sideEffectArg()))
}
`)

	mustWrite("gen/jscomp/logging_passes.go", `package jscomp

import (
	"sync/atomic"

	pkg_diagnostic "example.com/w29ctest/gen/diagnostic"
	"example.com/w29ctest/gen/j2gort"
	pkg_rhino "example.com/w29ctest/gen/rhino"
)

var DotFormatCalls atomic.Int64
var ImpureLogArgCalls atomic.Int64

type GraphvizGraph interface {
	GraphName() string
}

type ChunkGraph struct {
	name string
}

func (c *ChunkGraph) GraphName() string            { return c.name }
func (c *ChunkGraph) ToGraphvizGraph() *ChunkGraph { return c }

func DotFormatter_ToDot_1_GraphvizGraph(g GraphvizGraph) string {
	DotFormatCalls.Add(1)
	return "digraph " + g.GraphName() + " {}"
}

type StripConfigRecord struct {
	pass *LogDemoPass
}

func (*StripConfigRecord) Implements_j2gort_Ext_java_util_function_Supplier() {}
func (s *StripConfigRecord) IsJavaNil_() bool                                 { return s == nil }
func (s *StripConfigRecord) Get() any                                         { return j2gort.BoxString("config:" + s.pass.label) }

func NewStripCode_StripCodeConfigRecord(p *LogDemoPass) j2gort.Ext_java_util_function_Supplier {
	return &StripConfigRecord{pass: p}
}

func RemoveUnusedCode_RemovalLogRecord_ForVar(name string) *StripConfigRecord {
	return &StripConfigRecord{pass: &LogDemoPass{label: "var:" + name}}
}

type LogDemoPass struct {
	label        string
	chunkGraph   *ChunkGraph
	decisionsLog pkg_diagnostic.LogFile
	removalLog   pkg_diagnostic.LogFile
	typeLogFile  pkg_diagnostic.LogFile
}

func NewLogDemoPass(label string, log pkg_diagnostic.LogFile) *LogDemoPass {
	return &LogDemoPass{
		label:        label,
		chunkGraph:   &ChunkGraph{name: label},
		decisionsLog: log,
		removalLog:   log,
		typeLogFile:  log,
	}
}

func impureLogText() string {
	ImpureLogArgCalls.Add(1)
	return "impure"
}

func (this_ *LogDemoPass) RunLoggingHotPath(n pkg_rhino.Node, key string, count int32) {
	var chunkGraphLog pkg_diagnostic.LogFile = this_.decisionsLog
	_ = chunkGraphLog
	chunkGraphLog.Log_1_string(DotFormatter_ToDot_1_GraphvizGraph(j2gort.AsIface[GraphvizGraph](this_.chunkGraph.ToGraphvizGraph())))
	this_.typeLogFile.Log_1_string(n.ToString_3(true, false, true))
	this_.decisionsLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
		return j2gort.BoxString(j2gort.StrValString(key) + ": removing declaration")
	}))
	this_.decisionsLog.Log_1_Ext_java_util_function_Supplier(NewStripCode_StripCodeConfigRecord(this_))
	this_.removalLog.Log_1_Ext_java_util_function_Supplier(j2gort.ToIface_Ext_java_util_function_Supplier(RemoveUnusedCode_RemovalLogRecord_ForVar(key)))
	this_.decisionsLog.Log_2("opt %s %s", []any{j2gort.BoxString(key), j2gort.ValueOfInt(count)})

	// Already guarded in Java source -> MUST NOT be double-wrapped!
	if this_.removalLog.IsLogging() {
		this_.removalLog.Log_1_string("already-guarded")
	}
}

func (this_ *LogDemoPass) RunImpureLog() {
	// Negative control: impure argument must still execute even when !IsLogging()!
	this_.decisionsLog.Log_1_string(impureLogText())
}
`)

	mustWrite("gen/jscomp/remove_unused_code.go", `package jscomp

import (
	"example.com/w29ctest/gen/j2gort"
	pkg_rhino "example.com/w29ctest/gen/rhino"
)

type Var struct {
	name string
	hash int32
}

func NewVar(name string, hash int32) *Var {
	return &Var{name: name, hash: hash}
}

func (v *Var) GetName() string { return v.name }
func (v *Var) HashCode() int32 { return v.hash }

type RemoveUnusedCode_VarInfo struct {
	outer_    *RemoveUnusedCode
	varName   string
	removable bool
}

func InitRemoveUnusedCode_VarInfo_1(vi *RemoveUnusedCode_VarInfo, name string) {
	vi.varName = name
	vi.removable = true
}

type RemoveUnusedCode_Continuation struct {
	node pkg_rhino.Node
}

func NewRemoveUnusedCode_Continuation(outer *RemoveUnusedCode, n pkg_rhino.Node) *RemoveUnusedCode_Continuation {
	return &RemoveUnusedCode_Continuation{node: n}
}

type RemoveUnusedCode struct {
	varInfoKeys  []*Var
	varInfoVals  []*RemoveUnusedCode_VarInfo
	varInfoTable []int32
	varInfoSize  int32
	RemovedNames []string
}

func (this_ *RemoveUnusedCode) initFields_RemoveUnusedCode() {
	_ = this_
	this_.varInfoKeys = make([]*Var, int32(32))
	this_.varInfoVals = make([]*RemoveUnusedCode_VarInfo, int32(32))
	this_.varInfoTable = make([]int32, int32(64))
}

func NewRemoveUnusedCode() *RemoveUnusedCode {
	this_ := &RemoveUnusedCode{}
	this_.initFields_RemoveUnusedCode()
	return this_
}

func (this_ *RemoveUnusedCode) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node, vars []*Var) {
	_ = externs
	_ = root
	_ = NewRemoveUnusedCode_Continuation(this_, root)
	for _, v := range vars {
		vi := this_.priv_RemoveUnusedCode_getVarInfo_4622(v)
		if v.hash%2 == 0 {
			vi.removable = false
		}
	}
	for i := int32(0); i < this_.varInfoSize; i++ {
		if this_.varInfoVals[i].removable {
			this_.RemovedNames = append(this_.RemovedNames, this_.varInfoKeys[i].GetName())
		}
	}
}

func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_getVarInfo_4622(var_ *Var) *RemoveUnusedCode_VarInfo {
	j2gort.CheckNotNullPtr(var_)
	var h int32 = var_.HashCode()
	var table []int32 = this_.varInfoTable
	var mask int32 = int32(len(table)) - int32(1)
	var slot int32 = (h ^ int32(uint32(h)>>16)) & mask
	for {
		var entry int32 = table[slot]
		if entry == int32(0) {
			break
		}
		var idx int32 = entry - int32(1)
		var k *Var = this_.varInfoKeys[idx]
		if k == var_ || k.name == var_.name {
			return this_.varInfoVals[idx]
		}
		slot = (slot + int32(1)) & mask
	}
	var val_11 RemoveUnusedCode_VarInfo
	val_11.outer_ = this_
	InitRemoveUnusedCode_VarInfo_1(&val_11, var_.GetName())
	var varInfo *RemoveUnusedCode_VarInfo = &RemoveUnusedCode_VarInfo{
		outer_:    val_11.outer_,
		varName:   val_11.varName,
		removable: val_11.removable,
	}
	var newIdx int32 = this_.varInfoSize
	if newIdx == int32(len(this_.varInfoKeys)) {
		var newCap int32 = newIdx * int32(2)
		this_.varInfoKeys = func() []*Var {
			arrSrc_12 := this_.varInfoKeys
			arrLen_13 := newCap
			j2gort.CheckNotNullSlice(arrSrc_12)
			arrOut_14 := make([]*Var, arrLen_13)
			copy(arrOut_14, arrSrc_12)
			return arrOut_14
		}()
		this_.varInfoVals = func() []*RemoveUnusedCode_VarInfo {
			arrSrc_15 := this_.varInfoVals
			arrLen_16 := newCap
			j2gort.CheckNotNullSlice(arrSrc_15)
			arrOut_17 := make([]*RemoveUnusedCode_VarInfo, arrLen_16)
			copy(arrOut_17, arrSrc_15)
			return arrOut_17
		}()
	}
	this_.varInfoKeys[newIdx] = var_
	this_.varInfoVals[newIdx] = varInfo
	this_.varInfoSize = newIdx + int32(1)
	{
		a_22 := table
		i_23 := slot
		a_22[i_23] = newIdx + int32(1)
	}
	if (this_.varInfoSize * int32(4)) >= (int32(len(table)) * int32(3)) {
		this_.priv_RemoveUnusedCode_rehashVarInfoTable_4623(int32(len(table)) * int32(2))
	}
	return varInfo
}

func (this_ *RemoveUnusedCode) priv_RemoveUnusedCode_rehashVarInfoTable_4623(newTableLen int32) {
	var newTable []int32 = make([]int32, newTableLen)
	var mask int32 = newTableLen - int32(1)
	for i := int32(0); i < this_.varInfoSize; i++ {
		var h int32 = this_.varInfoKeys[i].HashCode()
		var slot int32 = (h ^ int32(uint32(h)>>16)) & mask
		for newTable[slot] != int32(0) {
			slot = (slot + int32(1)) & mask
		}
		newTable[slot] = i + int32(1)
	}
	this_.varInfoTable = newTable
}
`)

	// Run ALL passes on gen/jscomp: Run 1 + Run 2 idempotency check!
	allPasses := passes.AllPasses()
	jscompDir := filepath.Join(tmpDir, "gen/jscomp")

	res1, err := passes.RunOnPackageDir(jscompDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if res1.ByPass["j2gort_peephole"] < 15 {
		t.Fatalf("expected >= 15 j2gort_peephole sites on Run 1, got %+v", res1.ByPass)
	}
	if res1.ByPass["slab_alloc_lowering"] < 5 {
		t.Fatalf("expected >= 5 slab_alloc_lowering sites on Run 1, got %+v", res1.ByPass)
	}

	// Run 2 MUST produce 0 sites and 0 modified files!
	res2, err := passes.RunOnPackageDir(jscompDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if res2.TotalSites != 0 || res2.FilesModified != 0 {
		t.Fatalf("Run 2 NOT IDEMPOTENT: expected 0 sites and 0 modified files, got %+v", res2)
	}

	readLowered := func(rel string) string {
		t.Helper()
		b, rErr := os.ReadFile(filepath.Join(tmpDir, rel))
		if rErr != nil {
			t.Fatal(rErr)
		}
		return string(b)
	}

	// Verify AST shapes in ancestors_and_preconditions.go
	ancCode := readLowered("gen/jscomp/ancestors_and_preconditions.go")
	if strings.Count(ancCode, "GetAncestors()") != 1 {
		t.Fatalf("expected only the 1 escaping negative-control GetAncestors() call to remain, got:\n%s", ancCode)
	}
	for _, want := range []string{
		"curAnc_ = curAnc_.GetParent()",
		"if !(n.IsName()) {",
		`j2gort.CheckStateFmt(false, "Expected NAME node but got %s"`,
		"if !(!(n.HasChildren())) {",
		"if !(meta != nil && meta.ok) {",
		"if !(n != 0) {",
		"if !(idx >= 0) {",
		// Impure negative control MUST remain unguarded!
		`j2gort.CheckStateFmt(cond, "count=%s", j2gort.ValueOfInt(sideEffectArg()))`,
	} {
		if !strings.Contains(ancCode, want) {
			t.Fatalf("expected ancestors_and_preconditions.go to contain %q:\n%s", want, ancCode)
		}
	}

	// Verify AST shapes in logging_passes.go
	logCode := readLowered("gen/jscomp/logging_passes.go")
	for _, want := range []string{
		"if chunkGraphLog.IsLogging() {",
		"if this_.typeLogFile.IsLogging() {",
		"if this_.decisionsLog.IsLogging() {",
		"if this_.removalLog.IsLogging() {",
	} {
		if !strings.Contains(logCode, want) {
			t.Fatalf("expected logging_passes.go to contain %q:\n%s", want, logCode)
		}
	}
	if strings.Contains(logCode, "if this_.removalLog.IsLogging() {\n\t\tif this_.removalLog.IsLogging()") {
		t.Fatalf("double-wrapped pre-existing IsLogging() block in logging_passes.go:\n%s", logCode)
	}
	if strings.Contains(logCode, "if this_.decisionsLog.IsLogging() {\n\t\tthis_.decisionsLog.Log_1_string(impureLogText())") {
		t.Fatalf("impure Log_1_string argument must NOT be wrapped in IsLogging():\n%s", logCode)
	}

	// Verify AST shapes in remove_unused_code.go
	rucCode := readLowered("gen/jscomp/remove_unused_code.go")
	for _, want := range []string{
		"j2gort.AcquirePooledSlice[*Var](32)[:32]",
		"j2gort.AcquirePooledSlice[*RemoveUnusedCode_VarInfo](32)[:32]",
		"j2gort.AcquirePooledSlice[int32](64)[:64]",
		"j2gort.AcquirePooledSlice[*Var](int(arrLen_13))[:arrLen_13]",
		"j2gort.AcquirePooledSlice[*RemoveUnusedCode_VarInfo](int(arrLen_16))[:arrLen_16]",
		"j2gort.AcquirePooledSlice[int32](int(newTableLen))[:newTableLen]",
		"j2gort.ReleasePooledSlice(arrSrc_12)",
		"j2gort.ReleasePooledSlice(arrSrc_15)",
		"j2gort.ReleasePooledSlice(this_.varInfoTable)",
		"j2gort.ReleasePooledSlice(this_.varInfoKeys)",
		"j2gort.ReleasePooledSlice(this_.varInfoVals)",
	} {
		if !strings.Contains(rucCode, want) {
			t.Fatalf("expected remove_unused_code.go to contain %q:\n%s", want, rucCode)
		}
	}

	// 5. Runtime behavioral, zero-allocation, NPE/Exception parity, slice-pool, and -race test!
	mustWrite("gen/jscomp/wave29c_rt_test.go", `package jscomp

import (
	"fmt"
	"reflect"
	"sync"
	"testing"

	pkg_diagnostic "example.com/w29ctest/gen/diagnostic"
	"example.com/w29ctest/gen/j2gort"
	pkg_rhino "example.com/w29ctest/gen/rhino"
)

func TestWave29CRuntimeSemanticsZeroAllocsAndRace(t *testing.T) {
	// --- 1. GetAncestors(): root visits 0, normal chain, function barrier, NPE on 0, detach + continue, 0 allocs! ---
	leafInLoop := pkg_rhino.NewNode(pkg_rhino.TokNAME, "x")
	forNode := pkg_rhino.NewNode(pkg_rhino.TokFOR, "for1", leafInLoop)
	fnNode := pkg_rhino.NewNode(pkg_rhino.TokFUNCTION, "fn1", forNode)
	outerFor := pkg_rhino.NewNode(pkg_rhino.TokFOR, "outerFor", fnNode)
	scriptRoot := pkg_rhino.NewNode(pkg_rhino.TokSCRIPT, "root", outerFor)

	if NodeUtil_IsWithinLoop(scriptRoot) {
		t.Fatalf("root node has 0 ancestors and must NOT visit itself")
	}
	if NodeUtil_IsWithinLoop(outerFor) {
		t.Fatalf("outerFor has only scriptRoot ancestor and must NOT consider itself its own ancestor")
	}
	if !NodeUtil_IsWithinLoop(fnNode) {
		t.Fatalf("fnNode's parent is outerFor, so IsWithinLoop(fnNode) must be true")
	}
	leafDirectInFn := pkg_rhino.NewNode(pkg_rhino.TokNAME, "y")
	fnStop := pkg_rhino.NewNode(pkg_rhino.TokFUNCTION, "fnStop", leafDirectInFn)
	_ = pkg_rhino.NewNode(pkg_rhino.TokFOR, "enclosingFor", fnStop)
	if NodeUtil_IsWithinLoop(leafDirectInFn) {
		t.Fatalf("expected IsWithinLoop(leafDirectInFn) == false because fnStop breaks before enclosingFor")
	}
	if !NodeUtil_IsWithinLoop(leafInLoop) {
		t.Fatalf("expected IsWithinLoop(leafInLoop) == true")
	}

	// Zero heap allocations for NodeUtil_IsWithinLoop!
	ancAllocs := testing.AllocsPerRun(1000, func() {
		if !NodeUtil_IsWithinLoop(leafInLoop) || NodeUtil_IsWithinLoop(leafDirectInFn) {
			panic("unexpected")
		}
	})
	if ancAllocs != 0 {
		t.Fatalf("expected 0 allocs in lowered GetAncestors() loop, got %v", ancAllocs)
	}

	// Null Node(0) MUST throw NPE!
	func() {
		defer func() {
			r := recover()
			if _, ok := r.(j2gort.NPE); !ok {
				t.Fatalf("expected NPE on Node(0).GetAncestors(), got %#v", r)
			}
		}()
		_ = NodeUtil_IsWithinLoop(0)
	}()

	// Detach + continue mid-loop: leaf -> blockA -> forMid (detached!) -> blockB -> topScript
	nLeaf := pkg_rhino.NewNode(pkg_rhino.TokNAME, "leaf")
	nBlkA := pkg_rhino.NewNode(pkg_rhino.TokBLOCK, "blkA", nLeaf)
	nForMid := pkg_rhino.NewNode(pkg_rhino.TokFOR, "forMid", nBlkA)
	nBlkB := pkg_rhino.NewNode(pkg_rhino.TokBLOCK, "blkB", nForMid)
	_ = pkg_rhino.NewNode(pkg_rhino.TokSCRIPT, "topScript", nBlkB)

	gotLabels := DetachForAncestorsAndCollectLabels(nLeaf)
	wantLabels := []string{"blkA", "detached:forMid", "blkB", "topScript"}
	if !reflect.DeepEqual(gotLabels, wantLabels) {
		t.Fatalf("Detach + continue ancestor order mismatch: got %v, want %v", gotLabels, wantLabels)
	}
	if nForMid.GetParent() != 0 {
		t.Fatalf("expected nForMid to be detached")
	}

	// Nested GetAncestors + local var ancestors in CssHelper (0 allocs!):
	ch := &CssHelper{callNode: leafInLoop}
	// leafInLoop has 4 ancestors (for1, fn1, outerFor, root); outer counts: 3 + 2 + 1 + 0 = 6 -> 4 + 60 = 64
	if got := ch.CountAncestorPairs(); got != 64 {
		t.Fatalf("CountAncestorPairs = %d, want 64", got)
	}
	if allocs := testing.AllocsPerRun(500, func() { _ = ch.CountAncestorPairs() }); allocs != 0 {
		t.Fatalf("expected 0 allocs in nested CountAncestorPairs, got %v", allocs)
	}

	// --- 2. Lazy Preconditions: 0 allocs when cond == true, exact formatted message when cond == false! ---
	pureLeaf := pkg_rhino.NewNode(pkg_rhino.TokNAME, "constFn")
	meta := &MetaHolder{ok: true, tag: "good"}
	precondAllocs := testing.AllocsPerRun(1000, func() {
		AnnotateFunctionName(pureLeaf)
		CheckTokenToStringAndTree(pureLeaf, meta, 50000)
	})
	if precondAllocs != 0 {
		t.Fatalf("expected 0 allocs in lazy Preconditions when cond == true, got %v", precondAllocs)
	}

	// Verify exact formatted IllegalStateException when cond == false:
	func() {
		defer func() {
			r := recover()
			ise, ok := r.(j2gort.IllegalStateException)
			if !ok || ise.Msg != "Expected NAME node but got CALL" {
				t.Fatalf("expected IllegalStateException('Expected NAME node but got CALL'), got %#v", r)
			}
		}()
		callNode := pkg_rhino.NewNode(pkg_rhino.TokCALL, "foo")
		AnnotateFunctionName(callNode)
	}()

	// Verify exact formatted IllegalArgumentException when cond == false:
	func() {
		defer func() {
			r := recover()
			iae, ok := r.(j2gort.IllegalArgumentException)
			if !ok || iae.Msg != "Bad meta tag=null idx=77 node=NAME(constFn)" {
				t.Fatalf("unexpected IllegalArgumentException: %#v", r)
			}
		}()
		var nilMeta *MetaHolder
		_ = nilMeta
		CheckTokenToStringAndTree(pureLeaf, &MetaHolder{ok: false, tag: j2gort.NullString}, 77)
	}()

	// Impure negative control MUST execute sideEffectArg() even when cond == true!
	PrecondSideEffectCalls.Store(0)
	CheckStateImpureArg(true)
	if PrecondSideEffectCalls.Load() != 1 {
		t.Fatalf("impure Preconditions arg side effect was skipped: got %d", PrecondSideEffectCalls.Load())
	}

	// --- 3. Lazy LogFile IsLogging() Guard: 0 allocs & 0 DotFormatCalls when disabled, exact logs when enabled! ---
	DotFormatCalls.Store(0)
	noOp := &pkg_diagnostic.NoOpLogFile{}
	disabledPass := NewLogDemoPass("disabledGraph", noOp)
	logAllocs := testing.AllocsPerRun(1000, func() {
		disabledPass.RunLoggingHotPath(pureLeaf, "deadSym", 99999)
	})
	if logAllocs != 0 {
		t.Fatalf("expected 0 allocs when LogFile.IsLogging() == false, got %v", logAllocs)
	}
	if DotFormatCalls.Load() != 0 {
		t.Fatalf("DotFormatter ran %d times on disabled NoOpLogFile!", DotFormatCalls.Load())
	}

	// Impure negative control on disabled NoOpLogFile MUST still execute!
	ImpureLogArgCalls.Store(0)
	disabledPass.RunImpureLog()
	if ImpureLogArgCalls.Load() != 1 {
		t.Fatalf("impure LogFile arg side effect was skipped: got %d", ImpureLogArgCalls.Load())
	}

	// Enabled RecordingLogFile logs exact formatted entries!
	recLog := &pkg_diagnostic.RecordingLogFile{}
	enabledPass := NewLogDemoPass("liveGraph", recLog)
	enabledPass.RunLoggingHotPath(pureLeaf, "unusedX", 42)
	wantEntries := []string{
		"digraph liveGraph {}",
		"NAME(constFn)[true,false,true]",
		"unusedX: removing declaration",
		"config:liveGraph",
		"config:var:unusedX",
		"opt unusedX 42",
		"already-guarded",
	}
	if !reflect.DeepEqual(recLog.Entries, wantEntries) {
		t.Fatalf("enabled LogFile entries mismatch:\ngot:  %#v\nwant: %#v", recLog.Entries, wantEntries)
	}

	// --- 4. RemoveUnusedCode pooled VarInfo arrays across multiple doublings + zero live leaks! ---
	vars := make([]*Var, 300)
	var wantRemoved []string
	for i := 0; i < len(vars); i++ {
		nm := fmt.Sprintf("v_%d", i)
		h := int32((i + 1) * 2654435761 & 0x7fffffff)
		if i%2 == 1 {
			h |= 1 // odd -> removable
			wantRemoved = append(wantRemoved, nm)
		} else {
			h &^= 1 // even -> not removable
		}
		vars[i] = NewVar(nm, h)
	}

	j2gort.ResetSlicePoolCounters()
	ruc1 := NewRemoveUnusedCode()
	ruc1.Process_2(0, scriptRoot, vars)
	if !reflect.DeepEqual(ruc1.RemovedNames, wantRemoved) {
		t.Fatalf("ruc1 RemovedNames mismatch: got %d, want %d", len(ruc1.RemovedNames), len(wantRemoved))
	}
	if j2gort.PooledSliceAcquireCnt.Load() < 12 || j2gort.PooledSliceLiveCnt.Load() != 0 {
		t.Fatalf("expected >= 12 pooled slice acquires (init + 4 doublings x 3 arrays) and 0 live after ReleaseSlabs_, got acq=%d rel=%d live=%d",
			j2gort.PooledSliceAcquireCnt.Load(), j2gort.PooledSliceReleaseCnt.Load(), j2gort.PooledSliceLiveCnt.Load())
	}

	// Run 2 of RemoveUnusedCode MUST hit pooled tiers across initial + all intermediate doubling capacities!
	j2gort.ResetSlicePoolCounters()
	ruc2 := NewRemoveUnusedCode()
	ruc2.Process_2(0, scriptRoot, vars)
	if !reflect.DeepEqual(ruc2.RemovedNames, wantRemoved) {
		t.Fatalf("ruc2 RemovedNames mismatch after pool reuse!")
	}
	if j2gort.PooledSliceReuseCnt.Load() != j2gort.PooledSliceAcquireCnt.Load() || j2gort.PooledSliceLiveCnt.Load() != 0 {
		t.Fatalf("expected 100%% pooled slice tier reuse on second RemoveUnusedCode pass: acq=%d reuse=%d live=%d",
			j2gort.PooledSliceAcquireCnt.Load(), j2gort.PooledSliceReuseCnt.Load(), j2gort.PooledSliceLiveCnt.Load())
	}

	// --- 5. Concurrent -race stress test across all Wave 29C paths! ---
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for iter := 0; iter < 20; iter++ {
				if !NodeUtil_IsWithinLoop(leafInLoop) {
					panic("ancestors race")
				}
				AnnotateFunctionName(pureLeaf)
				CheckTokenToStringAndTree(pureLeaf, meta, int32(iter))
				disabledPass.RunLoggingHotPath(pureLeaf, "sym", int32(iter))
				r := NewRemoveUnusedCode()
				r.Process_2(0, scriptRoot, vars)
				if len(r.RemovedNames) != len(wantRemoved) {
					panic("ruc race")
				}
			}
		}()
	}
	wg.Wait()
}
`)

	cmd := exec.Command("go", "test", "-race", "-v", "./...")
	cmd.Dir = tmpDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("lowered Wave 29C runtime test with -race failed: %v\n%s", err, string(out))
	}
}

func TestWave29C_RealGenJscompCopyIdempotencyAndBuild(t *testing.T) {
	t.Parallel()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	repoRoot := filepath.Clean(filepath.Join(wd, "..", "..", ".."))
	if !realGenJ2gortReady(repoRoot) {
		t.Skip("gen/rhino or gen/j2gort is mid-edit; skipping real gen/jscomp copy check")
	}
	// Also verify real gen/rhino and gen/diagnostic currently build before running on a copy.
	checkCmd := exec.Command("go", "build", "./transpiled/gen/rhino", "./transpiled/gen/diagnostic")
	checkCmd.Dir = repoRoot
	if err := checkCmd.Run(); err != nil {
		t.Skipf("gen/rhino or gen/diagnostic is mid-edit (%v); skipping real gen/jscomp copy check", err)
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
	goMod := strings.Replace(string(origModBytes), "module github.com/paullepoulpe/jsgo", "module example.com/realw29c", 1) +
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

	w29Passes := []passes.Pass{
		passes.NewJ2gortPeepholePass(),
		passes.NewSlabAllocLoweringPass(),
	}
	res1, err := passes.RunOnPackageDir(copyDir, w29Passes, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 on copied real gen/jscomp failed: %v", err)
	}
	t.Logf("Wave 29C Real gen/jscomp Run 1 sites: total=%d files=%d byPass=%+v", res1.TotalSites, res1.FilesModified, res1.ByPass)

	// Verify real gen/jscomp lowered files have all 4 Wave 29C transformations!
	nuBytes, _ := os.ReadFile(filepath.Join(copyDir, "node_util.go"))
	if strings.Contains(string(nuBytes), "GetAncestors()") {
		t.Fatalf("expected NodeUtil_IsWithinLoop GetAncestors() to be lowered in real gen/jscomp/node_util.go")
	}
	normBytes, _ := os.ReadFile(filepath.Join(copyDir, "normalize.go"))
	if strings.Contains(string(normBytes), `j2gort.CheckStateFmt(n.IsName(), "Expected NAME node`) {
		t.Fatalf("expected Normalize.annotateFunctionExpressionNameAsConstant CheckStateFmt to be lazy-guarded in real gen/jscomp/normalize.go")
	}
	compBytes, _ := os.ReadFile(filepath.Join(copyDir, "compiler.go"))
	if !strings.Contains(string(compBytes), "if chunkGraphLog.IsLogging()") {
		t.Fatalf("expected chunkGraphLog.IsLogging() guard in real gen/jscomp/compiler.go")
	}
	rucBytes, _ := os.ReadFile(filepath.Join(copyDir, "remove_unused_code.go"))
	if !strings.Contains(string(rucBytes), "ReleasePooledSlice(this_.varInfoTable)") {
		t.Fatalf("expected RemoveUnusedCode varInfoTable pooling in real gen/jscomp/remove_unused_code.go")
	}

	// Verify lowered real gen/jscomp compiles cleanly!
	buildCmd := exec.Command("go", "build", "./jscomp")
	buildCmd.Dir = tmpDir
	if out, bErr := buildCmd.CombinedOutput(); bErr != nil {
		t.Fatalf("go build ./jscomp on lowered real gen/jscomp failed: %v\n%s", bErr, string(out))
	}

	// Run 2 on lowered real gen/jscomp MUST produce 0 sites and 0 modified files!
	res2, err := passes.RunOnPackageDir(copyDir, w29Passes, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 on lowered real gen/jscomp failed: %v", err)
	}
	if res2.TotalSites != 0 || res2.FilesModified != 0 {
		t.Fatalf("Run 2 on real gen/jscomp NOT IDEMPOTENT: got %+v", res2)
	}
}
