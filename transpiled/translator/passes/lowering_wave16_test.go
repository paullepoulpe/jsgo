package passes_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

// setupWave16J2gortOverlayWorkspace builds a temporary workspace overlaid with BOTH the Wave 16
// lowering runtime helpers AND the inline [4]uint64 Ext_java_util_BitSet implementation extracted
// directly from RuntimeCollectionsBindings.java.
func setupWave16J2gortOverlayWorkspace(t *testing.T) (tmpDir, repoRoot string) {
	t.Helper()
	tmpDir, repoRoot = setupRealJ2gortOverlayWorkspace(t)

	javaPath := filepath.Join(repoRoot, "transpiled", "translator", "java", "j2go", "RuntimeCollectionsBindings.java")
	rawJava, err := os.ReadFile(javaPath)
	if err != nil {
		t.Fatalf("read RuntimeCollectionsBindings.java: %v", err)
	}
	javaSrc := string(rawJava)

	bsStart := "// --- BitSet ("
	bsEnd := "// --- Iterables, Iterators, AbstractIterator, PeekingIterator, Ordering, ComparisonChain, Range ---"
	sIdx := strings.Index(javaSrc, bsStart)
	eIdx := strings.Index(javaSrc, bsEnd)
	if sIdx < 0 || eIdx <= sIdx {
		t.Fatalf("could not locate inline BitSet block in RuntimeCollectionsBindings.java (sIdx=%d, eIdx=%d)", sIdx, eIdx)
	}
	var bsLines []string
	for _, line := range strings.Split(javaSrc[sIdx:eIdx], "\n") {
		lt := strings.TrimSpace(line)
		if lt == `""";` || lt == `return """` || strings.HasPrefix(lt, "private static String buildColRt") || (lt == "}" && strings.HasPrefix(line, "  }")) {
			continue
		}
		trimmed := strings.TrimPrefix(line, "        ")
		trimmed = strings.ReplaceAll(trimmed, `\t`, "\t")
		trimmed = strings.ReplaceAll(trimmed, `\"`, `"`)
		bsLines = append(bsLines, trimmed)
	}
	inlineBitSetGo := strings.Join(bsLines, "\n") + "\n"

	rtPath := filepath.Join(tmpDir, "j2gort", "ext_collections_rt.go")
	rtBytes, err := os.ReadFile(rtPath)
	if err != nil {
		t.Fatal(err)
	}
	rtContent := string(rtBytes)
	oldBsStart := strings.Index(rtContent, "// --- BitSet ---")
	if oldBsStart < 0 {
		oldBsStart = strings.Index(rtContent, bsStart)
	}
	oldBsEnd := strings.Index(rtContent, bsEnd)
	if oldBsStart < 0 || oldBsEnd <= oldBsStart {
		t.Fatalf("could not locate BitSet section in ext_collections_rt.go (s=%d, e=%d)", oldBsStart, oldBsEnd)
	}
	rtContent = rtContent[:oldBsStart] + inlineBitSetGo + rtContent[oldBsEnd:]
	if err := os.WriteFile(rtPath, []byte(rtContent), 0o644); err != nil {
		t.Fatal(err)
	}

	extPath := filepath.Join(tmpDir, "j2gort", "ext_bindings.go")
	extBytes, err := os.ReadFile(extPath)
	if err != nil {
		t.Fatal(err)
	}
	extContent := string(extBytes)
	wantStructFields := "words []uint64 `gson:\"-\"`\n\tbuf   [4]uint64 `gson:\"-\"`"
	if strings.Contains(inlineBitSetGo, "bsInlineWords = 16") {
		wantStructFields = "words []uint64   `gson:\"-\"`\n\tbuf   [16]uint64 `gson:\"-\"`\n\thi    [8]int32   `gson:\"-\"`\n\tinitW int32      `gson:\"-\"`\n\thiLen int32      `gson:\"-\"`"
	} else if strings.Contains(inlineBitSetGo, "bsInlineWords = 2") {
		wantStructFields = "words []uint64  `gson:\"-\"`\n\tbuf   [2]uint64 `gson:\"-\"`\n\thi    [8]int32  `gson:\"-\"`\n\tinitW int32     `gson:\"-\"`\n\thiLen int32     `gson:\"-\"`"
	} else if strings.Contains(inlineBitSetGo, "bsSparseCap") {
		wantStructFields = "words []uint64  `gson:\"-\"`\n\tbuf   [8]uint64 `gson:\"-\"`\n\thi    [8]int32  `gson:\"-\"`\n\tinitW int32     `gson:\"-\"`\n\thiLen int32     `gson:\"-\"`"
	} else if strings.Contains(inlineBitSetGo, "hiBit") {
		wantStructFields = "words []uint64  `gson:\"-\"`\n\tbuf   [4]uint64 `gson:\"-\"`\n\tinitW int32     `gson:\"-\"`\n\thiBit int32     `gson:\"-\"`"
	}
	bsTypeIdx := strings.Index(extContent, "type Ext_java_util_BitSet struct {")
	if bsTypeIdx < 0 {
		t.Fatalf("could not locate Ext_java_util_BitSet struct block in ext_bindings.go")
	}
	bsEndBrace := strings.Index(extContent[bsTypeIdx:], "\n}\n")
	if bsEndBrace < 0 {
		t.Fatalf("could not locate closing brace of Ext_java_util_BitSet in ext_bindings.go")
	}
	extContent = extContent[:bsTypeIdx] + "type Ext_java_util_BitSet struct {\n\t" + wantStructFields + "\n}" + extContent[bsTypeIdx+bsEndBrace+2:]
	extContent = strings.Replace(extContent, "obj.Core_ = javautil.NewBitSet()", "InitExt_java_util_BitSet_0(obj)", 1)
	extContent = strings.Replace(extContent, "obj.Core_ = javautil.NewBitSetWithCapacity(int(arg0))", "InitExt_java_util_BitSet_1(obj, arg0)", 1)
	// Also add Flip_1, Set_2_int32_bool, Clear, Clear_1, Clear_2, Get_2, NextClearBit, PreviousClearBit, Xor wrappers if not yet emitted in ext_bindings.go.
	extraWrappers := `
func (this_ *Ext_java_util_BitSet) Flip_1(arg0 int32) { bitSetOp_Flip_1(this_, arg0) }
func (this_ *Ext_java_util_BitSet) Set_2_int32_bool(arg0 int32, arg1 bool) { bitSetOp_Set_2_int32_bool(this_, arg0, arg1) }
func (this_ *Ext_java_util_BitSet) Set_3(arg0 int32, arg1 int32, arg2 bool) { bitSetOp_Set_3(this_, arg0, arg1, arg2) }
func (this_ *Ext_java_util_BitSet) Clear_0()          { bitSetOp_Clear_0(this_) }
func (this_ *Ext_java_util_BitSet) Clear_1(arg0 int32) { bitSetOp_Clear_1(this_, arg0) }
func (this_ *Ext_java_util_BitSet) Clear_2(arg0 int32, arg1 int32) { bitSetOp_Clear_2(this_, arg0, arg1) }
func (this_ *Ext_java_util_BitSet) Get_2(arg0 int32, arg1 int32) *Ext_java_util_BitSet { return bitSetOp_Get_2(this_, arg0, arg1) }
func (this_ *Ext_java_util_BitSet) NextClearBit(arg0 int32) int32 { return bitSetOp_NextClearBit(this_, arg0) }
func (this_ *Ext_java_util_BitSet) PreviousSetBit(arg0 int32) int32 { return bitSetOp_PreviousSetBit(this_, arg0) }
func (this_ *Ext_java_util_BitSet) PreviousClearBit(arg0 int32) int32 { return bitSetOp_PreviousClearBit(this_, arg0) }
func (this_ *Ext_java_util_BitSet) And(arg0 *Ext_java_util_BitSet) { bitSetOp_And(this_, arg0) }
func (this_ *Ext_java_util_BitSet) Xor(arg0 *Ext_java_util_BitSet) { bitSetOp_Xor(this_, arg0) }
`
	for _, fnSig := range []string{"Flip_1(", "Set_2_int32_bool(", "Set_3(", "Clear_0(", "Clear_1(", "Clear_2(", "Get_2(", "NextClearBit(", "PreviousSetBit(", "PreviousClearBit(", "And(", "Xor("} {
		if strings.Contains(extContent, "func (this_ *Ext_java_util_BitSet) "+fnSig) {
			// Remove duplicate line from extraWrappers if already present in ext_bindings.go.
			var kept []string
			for _, l := range strings.Split(extraWrappers, "\n") {
				if !strings.Contains(l, " "+fnSig) {
					kept = append(kept, l)
				}
			}
			extraWrappers = strings.Join(kept, "\n")
		}
	}
	extContent += "\n" + extraWrappers
	if !strings.Contains(extContent, "func Ext_com_google_common_collect_Sets_NewIdentityHashSet(") {
		extContent += "\nfunc Ext_com_google_common_collect_Sets_NewIdentityHashSet() Ext_java_util_Set {\n\tobj := &Ext_java_util_HashSetBase{}\n\tobj.BindExtSelf_(obj)\n\tobj.inlineSet.identity = true\n\treturn obj\n}\n"
	}
	if err := os.WriteFile(extPath, []byte(extContent), 0o644); err != nil {
		t.Fatal(err)
	}
	return tmpDir, repoRoot
}

func TestWave16Task1_IndexThreshold8_And_Task4_InlineBitSet(t *testing.T) {
	tmpDir, _ := setupWave16J2gortOverlayWorkspace(t)

	testSrc := `package w16_rt_test

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"example.com/w13maptest/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javautil"
)

type IdentNode struct {
	ID int
}

func (n *IdentNode) IsJavaNil_() bool   { return n == nil }
func (n *IdentNode) HasIdentityEquals_() {}

type CustomKey struct {
	Name string
	Slot int32
}

func (k *CustomKey) IsJavaNil_() bool { return k == nil }
func (k *CustomKey) HashCode() int32  { return int32(len(k.Name))*31 + k.Slot }
func (k *CustomKey) Equals(other any) bool {
	o, ok := other.(*CustomKey)
	return ok && o != nil && k.Name == o.Name && k.Slot == o.Slot
}

func fieldIsNil(v any, path ...string) bool {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	for _, f := range path {
		rv = rv.FieldByName(f)
	}
	return rv.IsNil()
}

func fieldIsCleared(v any, path ...string) bool {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	for _, f := range path {
		rv = rv.FieldByName(f)
	}
	if rv.IsNil() || rv.Len() == 0 {
		return true
	}
	if rv.Kind() == reflect.Slice {
		for i := 0; i < rv.Len(); i++ {
			if !rv.Index(i).IsZero() {
				return false
			}
		}
		return true
	}
	return false
}

func TestIndexThreshold8AcrossAllOrderedMapsAndSets(t *testing.T) {
	// 1. OrderedMap[string, int32] (orderedMapIndexThreshold = 12)
	var om j2gort.OrderedMap[string, int32]
	for i := 1; i <= 12; i++ {
		k := fmt.Sprintf("k%d", i)
		om.Put(k, int32(i*10))
		if !fieldIsNil(&om, "idx") {
			t.Fatalf("OrderedMap allocated idx at len=%d (want nil up to 12)", i)
		}
	}
	// Updating existing keys and removing+reinserting via linear scan while <= 12 must keep idx == nil!
	om.Put("k4", 400)
	if om.Get("k4") != 400 || !om.ContainsKey("k12") || om.ContainsKey("missing") {
		t.Fatalf("OrderedMap linear lookup mismatch")
	}
	om.Remove("k2")
	if om.ContainsKey("k2") || om.Size() != 11 || !fieldIsNil(&om, "idx") {
		t.Fatalf("OrderedMap linear remove mismatch")
	}
	// Because linear Remove compacts in-place at <=12 without tombstones, inserting k13 brings len(entries) back to 12 (idx still nil)!
	om.Put("k13", 130)
	if !fieldIsNil(&om, "idx") || om.Size() != 12 {
		t.Fatalf("OrderedMap expected idx == nil at 12 live entries after linear remove+insert")
	}
	// Pushing len(entries) past 12 allocates idx!
	om.Put("k14", 140)
	if fieldIsNil(&om, "idx") {
		t.Fatalf("OrderedMap expected idx != nil once len(entries)=13 > 12")
	}
	if om.Get("k1") != 10 || om.Get("k4") != 400 || om.Get("k13") != 130 || om.Get("k14") != 140 || om.ContainsKey("k2") {
		t.Fatalf("OrderedMap post-promotion lookup mismatch")
	}
	om.Clear()
	if !fieldIsCleared(&om, "idx") || om.Size() != 0 || om.ContainsKey("k1") || om.ContainsKey("k14") {
		t.Fatalf("OrderedMap.Clear() must clear idx and entries")
	}

	// 2. IdentityOrderedMap[*IdentNode, int32] (orderedMapIndexThreshold = 12)
	nodes := make([]*IdentNode, 18)
	for i := range nodes {
		nodes[i] = &IdentNode{ID: 7} // identical fields, distinct pointers
	}
	var im j2gort.IdentityOrderedMap[*IdentNode, int32]
	for i := 0; i < 12; i++ {
		im.Put(nodes[i], int32(i+1))
		if !fieldIsNil(&im, "raw", "idx") {
			t.Fatalf("IdentityOrderedMap allocated idx at %d", i+1)
		}
	}
	im.Put(nodes[12], 13)
	if fieldIsNil(&im, "raw", "idx") || im.Size() != 13 || im.Get(nodes[0]) != 1 || im.Get(nodes[12]) != 13 || im.ContainsKey(nodes[13]) {
		t.Fatalf("IdentityOrderedMap promotion/lookup mismatch")
	}
	im.Clear()
	if !fieldIsCleared(&im, "raw", "idx") || im.Size() != 0 || im.ContainsKey(nodes[0]) || im.ContainsKey(nodes[12]) {
		t.Fatalf("IdentityOrderedMap.Clear() must clear idx and entries")
	}

	// 3. CustomOrderedMap[*CustomKey, int32] (orderedMapIndexThreshold = 12)
	var cm j2gort.CustomOrderedMap[*CustomKey, int32]
	for i := 1; i <= 12; i++ {
		cm.Put(&CustomKey{Name: "v", Slot: int32(i)}, int32(i))
		if !fieldIsNil(&cm, "buckets") {
			t.Fatalf("CustomOrderedMap allocated buckets at %d (want nil up to 12)", i)
		}
	}
	// Equal custom key via fresh pointer updates in-place without growing entries past 12!
	if prev := cm.Put(&CustomKey{Name: "v", Slot: 5}, 55); prev != 5 || !fieldIsNil(&cm, "buckets") || cm.Size() != 12 {
		t.Fatalf("CustomOrderedMap <=12 custom Equals update failed: prev=%d size=%d", prev, cm.Size())
	}
	cm.Put(&CustomKey{Name: "v", Slot: 13}, 133)
	if fieldIsNil(&cm, "buckets") || cm.Get(&CustomKey{Name: "v", Slot: 5}) != 55 || cm.Get(&CustomKey{Name: "v", Slot: 13}) != 133 {
		t.Fatalf("CustomOrderedMap >12 bucket promotion failed")
	}
	cm.Clear()
	if !fieldIsCleared(&cm, "buckets") || cm.Size() != 0 || cm.ContainsKey(&CustomKey{Name: "v", Slot: 5}) {
		t.Fatalf("CustomOrderedMap.Clear() must clear buckets and entries")
	}

	// 4. OrderedSet, IdentityOrderedSet, LocalOrderedSet (threshold = 16), CustomOrderedSet (backed by CustomOrderedMap, threshold = 12)
	var os j2gort.OrderedSet[int32]
	var is j2gort.IdentityOrderedSet[*IdentNode]
	var cs j2gort.CustomOrderedSet[*CustomKey]
	var ls j2gort.LocalOrderedSet[int32]
	for i := 0; i < 16; i++ {
		os.Add(int32(i))
		is.Add(nodes[i])
		ls.Add(int32(i))
	}
	for i := 0; i < 12; i++ {
		cs.Add(&CustomKey{Name: "s", Slot: int32(i)})
	}
	if !fieldIsNil(&os, "idx") || !fieldIsNil(&is, "raw", "idx") || !fieldIsNil(&cs, "m", "buckets") || !fieldIsNil(&ls, "idx") {
		t.Fatalf("Expected all Set types to keep idx/buckets == nil at threshold")
	}
	// CustomOrderedSet structural duplicate at size 12 does not allocate buckets
	if cs.Add(&CustomKey{Name: "s", Slot: 3}) || !fieldIsNil(&cs, "m", "buckets") {
		t.Fatalf("CustomOrderedSet duplicate add failed")
	}
	os.Add(99)
	is.Add(nodes[16])
	cs.Add(&CustomKey{Name: "s", Slot: 99})
	ls.Add(99)
	if fieldIsNil(&os, "idx") || fieldIsNil(&is, "raw", "idx") || fieldIsNil(&cs, "m", "buckets") || fieldIsNil(&ls, "idx") {
		t.Fatalf("Expected all Set types to allocate index past threshold")
	}
	os.Clear()
	is.Clear()
	cs.Clear()
	ls.Clear()
	if !fieldIsCleared(&os, "idx") || !fieldIsCleared(&is, "raw", "idx") || !fieldIsCleared(&cs, "m", "buckets") || !fieldIsCleared(&ls, "idx") {
		t.Fatalf("Expected Clear() on all Set types to clear idx/buckets")
	}
}

func TestInlineBitSetZeroSecondaryAllocAndParity(t *testing.T) {
	// Zero heap allocations when mutating an inline Ext_java_util_BitSet up to 256 bits!
	var stackBS j2gort.Ext_java_util_BitSet
	allocs := testing.AllocsPerRun(100, func() {
		stackBS.Clear_0()
		stackBS.Set_1(0)
		stackBS.Set_1(63)
		stackBS.Set_1(128)
		stackBS.Set_1(255)
		stackBS.Flip_2(10, 200)
		if stackBS.Cardinality() == 0 || !stackBS.Get_1(255) {
			t.Fatal("unexpected BitSet state")
		}
	})
	if allocs != 0 {
		t.Fatalf("expected 0 heap allocations for <=256-bit inline BitSet operations, got %v", allocs)
	}

	// NewExt_java_util_BitSet() + <=256-bit operations must perform at most 1 alloc (the struct itself, 0 secondaryjavautil.BitSet/slice allocs).
	var sink *j2gort.Ext_java_util_BitSet
	ctorAllocs := testing.AllocsPerRun(100, func() {
		b := j2gort.NewExt_java_util_BitSet_1(192)
		b.Set_1(5)
		b.Set_1(191)
		b.Set_1(255)
		sink = b
	})
	_ = sink
	if ctorAllocs > 1 {
		t.Fatalf("expected <=1 allocation for NewExt_java_util_BitSet_1(192) + <=256-bit ops, got %v", ctorAllocs)
	}

	// Differential parity against javautil.BitSet + Java BitSet bit-oracle across inline (<=256) and heap (>256) bit indices!
	oraclePrevClearBit := func(ref *javautil.BitSet, from int) int {
		for i := from; i >= 0; i-- {
			if !ref.Get(i) {
				return i
			}
		}
		return -1
	}
	oracleToByteArray := func(ref *javautil.BitSet) []int8 {
		ln := ref.Length()
		if ln == 0 {
			return []int8{}
		}
		out := make([]int8, (ln+7)/8)
		for i := 0; i < ln; i++ {
			if ref.Get(i) {
				out[i>>3] |= int8(uint8(1) << (uint(i) & 7))
			}
		}
		return out
	}
	checkParity := func(label string, ext *j2gort.Ext_java_util_BitSet, ref *javautil.BitSet) {
		t.Helper()
		if int(ext.Length()) != ref.Length() {
			t.Fatalf("%s: Length got %d, want %d", label, ext.Length(), ref.Length())
		}
		if ext.IsEmpty() != ref.IsEmpty() {
			t.Fatalf("%s: IsEmpty got %v, want %v", label, ext.IsEmpty(), ref.IsEmpty())
		}
		if int(ext.Cardinality()) != ref.Cardinality() {
			t.Fatalf("%s: Cardinality got %d, want %d", label, ext.Cardinality(), ref.Cardinality())
		}
		if ext.HashCode() != ref.HashCode() {
			t.Fatalf("%s: HashCode got %d, want %d", label, ext.HashCode(), ref.HashCode())
		}
		if ext.String() != ref.String() {
			t.Fatalf("%s: String got %q, want %q", label, ext.String(), ref.String())
		}
		gotBytes := ext.ToByteArray()
		refBytes := oracleToByteArray(ref)
		if len(gotBytes) != len(refBytes) {
			t.Fatalf("%s: ToByteArray len got %d, want %d", label, len(gotBytes), len(refBytes))
		}
		for i := range refBytes {
			if gotBytes[i] != refBytes[i] {
				t.Fatalf("%s: ToByteArray[%d] got %d, want %d", label, i, gotBytes[i], refBytes[i])
			}
		}
		for probe := 0; probe <= 400; probe += 17 {
			if ext.Get_1(int32(probe)) != ref.Get(probe) {
				t.Fatalf("%s: Get(%d) mismatch", label, probe)
			}
			if int(ext.NextSetBit(int32(probe))) != ref.NextSetBit(probe) {
				t.Fatalf("%s: NextSetBit(%d) got %d, want %d", label, probe, ext.NextSetBit(int32(probe)), ref.NextSetBit(probe))
			}
			if int(ext.NextClearBit(int32(probe))) != ref.NextClearBit(probe) {
				t.Fatalf("%s: NextClearBit(%d) got %d, want %d", label, probe, ext.NextClearBit(int32(probe)), ref.NextClearBit(probe))
			}
			if int(ext.PreviousSetBit(int32(probe))) != ref.PreviousSetBit(probe) {
				t.Fatalf("%s: PreviousSetBit(%d) got %d, want %d", label, probe, ext.PreviousSetBit(int32(probe)), ref.PreviousSetBit(probe))
			}
			if int(ext.PreviousClearBit(int32(probe))) != oraclePrevClearBit(ref, probe) {
				t.Fatalf("%s: PreviousClearBit(%d) got %d, want %d", label, probe, ext.PreviousClearBit(int32(probe)), oraclePrevClearBit(ref, probe))
			}
		}
		if int(ext.PreviousSetBit(-1)) != -1 || int(ext.PreviousClearBit(-1)) != -1 {
			t.Fatalf("%s: Previous*(-1) expected -1", label)
		}
	}

	aExt, aRef := j2gort.NewExt_java_util_BitSet(), javautil.NewBitSet()
	bExt, bRef := j2gort.NewExt_java_util_BitSet_1(64), javautil.NewBitSetWithCapacity(64)
	checkParity("empty", aExt, aRef)

	for _, idx := range []int{0, 1, 63, 64, 127, 191, 255} {
		aExt.Set_1(int32(idx))
		aRef.Set(idx)
	}
	checkParity("inline-set", aExt, aRef)

	// Cross 256-bit inline threshold into heap slice!
	aExt.Set_2_int32_int32(250, 330)
	aRef.SetRange(250, 330)
	checkParity("promoted-range", aExt, aRef)

	aExt.Flip_2(60, 270)
	aRef.FlipRange(60, 270)
	aExt.Flip_1(300)
	aRef.Flip(300)
	checkParity("flipped", aExt, aRef)

	bExt.Set_3(50, 310, true)
	bRef.SetRange(50, 310)
	bExt.Set_2_int32_bool(100, false)
	bRef.SetVal(100, false)
	if aExt.Intersects(bExt) != aRef.Intersects(bRef) {
		t.Fatalf("Intersects mismatch")
	}

	subExt := aExt.Get_2(62, 265)
	subRef := javautil.NewBitSet()
	for i := 62; i < 265; i++ {
		if aRef.Get(i) {
			subRef.Set(i - 62)
		}
	}
	checkParity("subrange-get", subExt, subRef)

	cloneExt := aExt.Clone().(*j2gort.Ext_java_util_BitSet)
	if !cloneExt.Equals(aExt) || cloneExt == aExt {
		t.Fatalf("Clone/Equals mismatch")
	}
	cloneExt.Set_1(7)
	if cloneExt.Equals(aExt) == aRef.Equals(func() *javautil.BitSet { c := aRef.Clone(); c.Set(7); return c }()) {
		// both false
	} else {
		t.Fatalf("Mutated clone Equals mismatch")
	}

	aExt.Or(bExt)
	aRef.Or(bRef)
	checkParity("or", aExt, aRef)

	aExt.Xor(subExt)
	aRef.Xor(subRef)
	checkParity("xor", aExt, aRef)

	aExt.AndNot(bExt)
	aRef.AndNot(bRef)
	checkParity("andnot", aExt, aRef)

	bExt.And(aExt)
	bRef.And(aRef)
	checkParity("and", bExt, bRef)

	aExt.Clear_2(10, 350)
	for i := 10; i < 350; i++ {
		aRef.Clear(i)
	}
	checkParity("clear-range", aExt, aRef)

	aExt.Clear_1(0)
	aRef.Clear(0)
	checkParity("clear-bit-0", aExt, aRef)
	_ = bytes.Equal
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "w16_rt_test.go"), []byte(testSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "test", "-v", ".")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Wave 16 Task 1 & 4 runtime test failed: %v\n%s", err, string(out))
	}
}

func TestWave16Task2And3_IdentityMapSet_CustomOrderedSet_And_ComparableTypeParams(t *testing.T) {
	tmpDir, _ := setupWave16J2gortOverlayWorkspace(t)

	buildWorkload := func(pkgName string) string {
		return "package " + pkgName + `

import (
	"fmt"
	"example.com/w13maptest/j2gort"
)

type Node uint32

// --- JSType / ObjectType / *TemplateType hierarchy (mirrors rhino/template_type_replacer.go) ---
type JSType interface {
	IsJavaNil_() bool
	HashCode() int32
	Equals(any) bool
	TypeLabel() string
}

type ObjectType interface {
	JSType
	ObjTag() string
}

type ObjTypeImpl struct {
	Label string
	Kind  int32
}

func (o *ObjTypeImpl) IsJavaNil_() bool { return o == nil }
func (o *ObjTypeImpl) TypeLabel() string { return o.Label }
func (o *ObjTypeImpl) ObjTag() string    { return o.Label }
func (o *ObjTypeImpl) HashCode() int32   { return int32(len(o.Label))*31 + o.Kind }
func (o *ObjTypeImpl) Equals(other any) bool {
	x, ok := other.(*ObjTypeImpl)
	return ok && x != nil && o.Label == x.Label && o.Kind == x.Kind
}

type TemplateType struct {
	Label string
}

func (t *TemplateType) IsJavaNil_() bool { return t == nil }
func (t *TemplateType) TypeLabel() string { return t.Label }
func (t *TemplateType) HashCode() int32   { return int32(len(t.Label)) * 17 }
func (t *TemplateType) Equals(other any) bool {
	x, ok := other.(*TemplateType)
	return ok && x != nil && t.Label == x.Label
}

// 1. TemplateTypeReplacer: seenTypes (Sets_NewIdentityHashSet with *TemplateType + JSType)
//    and visitedObjectTypes (NewExt_java_util_IdentityHashMap with ObjectType + JSType)!
type TemplateTypeReplacer struct {
	seenTypes          j2gort.Ext_java_util_Set
	visitedObjectTypes *j2gort.Ext_java_util_IdentityHashMap
}

func NewTemplateTypeReplacer() *TemplateTypeReplacer {
	this_ := &TemplateTypeReplacer{}
	this_.seenTypes = nil
	this_.visitedObjectTypes = nil
	return this_
}

func (this_ *TemplateTypeReplacer) priv_TemplateTypeReplacer_initSeenTypes_1() {
	if this_.seenTypes == nil {
		this_.seenTypes = j2gort.Ext_com_google_common_collect_Sets_NewIdentityHashSet()
	}
}

func (this_ *TemplateTypeReplacer) VisitTemplate(tt *TemplateType, other JSType) bool {
	this_.priv_TemplateTypeReplacer_initSeenTypes_1()
	if this_.seenTypes.Contains(j2gort.PtrToAny(tt)) || this_.seenTypes.Contains(other) {
		return false
	}
	this_.seenTypes.Add_1(j2gort.PtrToAny(tt))
	added := this_.seenTypes.Add_1(other)
	this_.seenTypes.Remove_1_any(j2gort.PtrToAny(tt))
	return added
}

func (this_ *TemplateTypeReplacer) RecordAndLookupObj(obj ObjectType, mapped JSType, probe JSType) JSType {
	if this_.visitedObjectTypes == nil {
		this_.visitedObjectTypes = j2gort.NewExt_java_util_IdentityHashMap()
	}
	this_.visitedObjectTypes.Put_2(obj, mapped)
	return j2gort.AsIface[JSType](this_.visitedObjectTypes.Get(probe))
}

// 2. LinkedIdentityHashMap[K comparable, V any] & IdentitySetHolder[E comparable]
type LinkedIdentityEntry[K comparable, V any] struct {
	Key K
	Val V
}

func (e *LinkedIdentityEntry[K, V]) IsJavaNil_() bool   { return e == nil }
func (*LinkedIdentityEntry[K, V]) HasIdentityEquals_() {}

type LinkedIdentityHashMap[K comparable, V any] struct {
	innerMap *j2gort.Ext_java_util_IdentityHashMap
	seenKeys j2gort.Ext_java_util_Set
}

func (this_ *LinkedIdentityHashMap[K, V]) initFields_LinkedIdentityHashMap() {
	this_.innerMap = j2gort.NewExt_java_util_IdentityHashMap()
	this_.seenKeys = j2gort.Ext_com_google_common_collect_Sets_NewIdentityHashSet()
}

func NewLinkedIdentityHashMap[K comparable, V any]() *LinkedIdentityHashMap[K, V] {
	this_ := &LinkedIdentityHashMap[K, V]{}
	InitLinkedIdentityHashMap[K, V](this_)
	return this_
}

func InitLinkedIdentityHashMap[K comparable, V any](this_ *LinkedIdentityHashMap[K, V]) {
	this_.initFields_LinkedIdentityHashMap()
}

func (this_ *LinkedIdentityHashMap[K, V]) Put(key K, val V) bool {
	this_.innerMap.Put_2(j2gort.ToAny(key), j2gort.PtrToAny(&LinkedIdentityEntry[K, V]{Key: key, Val: val}))
	return this_.seenKeys.Add_1(j2gort.ToAny(key))
}

func (this_ *LinkedIdentityHashMap[K, V]) Get(key K) (V, bool) {
	if !this_.seenKeys.Contains(j2gort.ToAny(key)) {
		var zv V
		return zv, false
	}
	ent := j2gort.CastPtr[LinkedIdentityEntry[K, V]](this_.innerMap.Get(j2gort.ToAny(key)))
	return ent.Val, true
}

// 3. CustomOrderedSet[JSType] & *CustomOrderedSet[*ObjTypeImpl]
type TypeSetCollector struct {
	uniqueTypes j2gort.Ext_java_util_Set
	lazyObjs    *j2gort.Ext_java_util_LinkedHashSet
}

func NewTypeSetCollector() *TypeSetCollector {
	this_ := &TypeSetCollector{}
	InitTypeSetCollector(this_)
	return this_
}

func InitTypeSetCollector(this_ *TypeSetCollector) {
	this_.uniqueTypes = j2gort.NewExt_java_util_LinkedHashSet()
	this_.lazyObjs = nil
}

func (this_ *TypeSetCollector) AddTypes(items []JSType, objs []*ObjTypeImpl) string {
	for _, it := range items {
		this_.uniqueTypes.Add_1(it)
	}
	if this_.lazyObjs == nil {
		this_.lazyObjs = j2gort.NewExt_java_util_LinkedHashSet()
	}
	for _, ob := range objs {
		this_.lazyObjs.Add_1(j2gort.PtrToAny(ob))
	}
	var out string
	{
		cur := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.uniqueTypes))
		for cur.HasNext() {
			raw := cur.Next()
			jt := j2gort.AsIface[JSType](raw)
			out += jt.TypeLabel() + ","
		}
	}
	return fmt.Sprintf("%s|%d|%d", out, this_.uniqueTypes.Size(), this_.lazyObjs.Size())
}

// 4. Comparable TypeParam Graph/Queue structures:
//    LinkedDirectedGraphBase[N comparable, E any], FixedPointGraphTraversal[N comparable, E any], UniqueQueue[T comparable]
type DiGraphNode[N any, E any] interface {
	IsJavaNil_() bool
	Implements_DiGraphNode()
	GetValue() N
}

type LinkedDiGraphNode[N any, E any] struct {
	val N
}

func (n *LinkedDiGraphNode[N, E]) IsJavaNil_() bool        { return n == nil }
func (*LinkedDiGraphNode[N, E]) HasIdentityEquals_()       {}
func (*LinkedDiGraphNode[N, E]) Implements_DiGraphNode()   {}
func (n *LinkedDiGraphNode[N, E]) GetValue() N             { return n.val }

type LinkedDirectedGraph[N comparable, E any] struct {
	nodes *j2gort.Ext_java_util_LinkedHashMap
}

func (this_ *LinkedDirectedGraph[N, E]) initFields_LinkedDirectedGraph() {
	this_.nodes = j2gort.NewExt_java_util_LinkedHashMap()
}

func NewLinkedDirectedGraph[N comparable, E any]() *LinkedDirectedGraph[N, E] {
	this_ := &LinkedDirectedGraph[N, E]{}
	InitLinkedDirectedGraph[N, E](this_)
	return this_
}

func InitLinkedDirectedGraph[N comparable, E any](this_ *LinkedDirectedGraph[N, E]) {
	this_.initFields_LinkedDirectedGraph()
}

func (this_ *LinkedDirectedGraph[N, E]) CreateNode(v N) DiGraphNode[N, E] {
	n := j2gort.CastPtr[LinkedDiGraphNode[N, E]](this_.nodes.Get(j2gort.ToAny(v)))
	if n == nil {
		n = &LinkedDiGraphNode[N, E]{val: v}
		this_.nodes.Put_2(j2gort.ToAny(v), j2gort.PtrToAny(n))
	}
	return n
}

type FixedPointTraversal[N comparable, E any] struct{}

func (this_ *FixedPointTraversal[N, E]) Traverse(g *LinkedDirectedGraph[N, E], seeds []N) int32 {
	var workSet *j2gort.Ext_java_util_LinkedHashSet = j2gort.NewExt_java_util_LinkedHashSet()
	for _, s := range seeds {
		workSet.Add_1(g.CreateNode(s))
	}
	var visited int32
	for !workSet.IsEmpty() {
		first := j2gort.AsIface[DiGraphNode[N, E]](workSet.GetFirst())
		this_.priv_FixedPointTraversal_visitNode_1(first, workSet)
		visited++
	}
	return visited
}

func (this_ *FixedPointTraversal[N, E]) priv_FixedPointTraversal_visitNode_1(node DiGraphNode[N, E], workSet *j2gort.Ext_java_util_LinkedHashSet) {
	workSet.Remove_1_any(node)
}

type UniqueQueue[T comparable] struct {
	seenSet *j2gort.Ext_java_util_LinkedHashSet
}

func NewUniqueQueue[T comparable]() *UniqueQueue[T] {
	this_ := &UniqueQueue[T]{}
	InitUniqueQueue[T](this_)
	return this_
}

func InitUniqueQueue[T comparable](this_ *UniqueQueue[T]) {
	this_.seenSet = j2gort.NewExt_java_util_LinkedHashSet()
}

func (this_ *UniqueQueue[T]) PushUnique(v T) bool {
	return this_.seenSet.Add_1(j2gort.ToAny(v))
}

// 5. Soundness Negative Guards:
//    (a) IdentityHashMap with string keys MUST NOT lower to IdentityOrderedMap!
//    (b) Generic HashSet[K comparable] instantiated with custom-Equals *ObjTypeImpl MUST NOT lower to OrderedSet[K]!
type UnsafeGenericSet[K comparable] struct {
	set *j2gort.Ext_java_util_LinkedHashSet
}

func NewUnsafeGenericSet[K comparable]() *UnsafeGenericSet[K] {
	return &UnsafeGenericSet[K]{set: j2gort.NewExt_java_util_LinkedHashSet()}
}

func (u *UnsafeGenericSet[K]) Add(k K) bool {
	return u.set.Add_1(j2gort.ToAny(k))
}

func RunAllScenarios() string {
	// Scenario 1: TemplateTypeReplacer (identity vs structural equality on JSType!)
	o1 := &ObjTypeImpl{Label: "Foo", Kind: 1}
	o2 := &ObjTypeImpl{Label: "Foo", Kind: 1} // .Equals(o1) == true, but o1 != o2!
	tt := &TemplateType{Label: "T"}
	ttr := NewTemplateTypeReplacer()
	v1 := ttr.VisitTemplate(tt, o1)
	v2 := ttr.VisitTemplate(tt, o2) // must be true because IdentityOrderedSet uses ==, not .Equals!
	v3 := ttr.VisitTemplate(tt, o1) // false (already in seenTypes)
	m1 := ttr.RecordAndLookupObj(o1, tt, o1)
	m2 := ttr.RecordAndLookupObj(o1, tt, o2) // must be nil because IdentityOrderedMap uses ==!

	// Scenario 2: LinkedIdentityHashMap[JSType, string]
	lim := NewLinkedIdentityHashMap[JSType, string]()
	p1 := lim.Put(o1, "one")
	p2 := lim.Put(o2, "two") // distinct identity!
	g1, _ := lim.Get(o1)
	g2, _ := lim.Get(o2)

	// Scenario 3: CustomOrderedSet[JSType] & *CustomOrderedSet[*ObjTypeImpl] (structural .Equals deduplication!)
	tsc := NewTypeSetCollector()
	colSummary := tsc.AddTypes([]JSType{o1, o2, tt}, []*ObjTypeImpl{o1, o2}) // o2 deduplicated!

	// Scenario 4: Comparable TypeParam graph + traversal + UniqueQueue
	g := NewLinkedDirectedGraph[Node, string]()
	tr := &FixedPointTraversal[Node, string]{}
	travCount := tr.Traverse(g, []Node{10, 20, 10, 30})
	uq := NewUniqueQueue[string]()
	uqCount := 0
	for _, s := range []string{"a", "b", "a", "c"} {
		if uq.PushUnique(s) {
			uqCount++
		}
	}

	// Scenario 5: Soundness guards
	var strIdentMap *j2gort.Ext_java_util_IdentityHashMap = j2gort.NewExt_java_util_IdentityHashMap()
	strIdentMap.Put_2(j2gort.BoxString("k"), j2gort.BoxString("v"))
	ug := NewUnsafeGenericSet[*ObjTypeImpl]()
	ug1 := ug.Add(o1)
	ug2 := ug.Add(o2) // MUST be false because HashSet uses .Equals, NOT pointer ==!

	return fmt.Sprintf("ttr=%v,%v,%v,m1=%v,m2=%v|lim=%v,%v,%s,%s|col=%s|graph=%d,uq=%d|guards=%s,%v,%v",
		v1, v2, v3, m1 != nil, m2 != nil, p1, p2, g1, g2, colSummary, travCount, uqCount,
		j2gort.CastString(strIdentMap.Get(j2gort.BoxString("k"))), ug1, ug2)
}
`
	}

	origDir := filepath.Join(tmpDir, "orig")
	optDir := filepath.Join(tmpDir, "opt")
	if err := os.MkdirAll(origDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(optDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(origDir, "workload.go"), []byte(buildWorkload("orig")), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(optDir, "workload.go"), []byte(buildWorkload("opt")), 0o644); err != nil {
		t.Fatal(err)
	}

	allPasses := passes.AllPasses()
	stats1, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if stats1.ByPass["native_map_lowering"] == 0 || stats1.ByPass["native_set_lowering"] == 0 {
		t.Fatalf("expected both map and set lowering to fire, got %+v", stats1.ByPass)
	}

	// Run 2 MUST be 100% idempotent!
	stats2, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites and 0 files modified on Run 2, got %+v", stats2)
	}

	gotBytes, err := os.ReadFile(filepath.Join(optDir, "workload.go"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(gotBytes)
	for _, want := range []string{
		"seenTypes          *j2gort.IdentityOrderedSet[JSType]",
		"visitedObjectTypes *j2gort.IdentityOrderedMap[JSType, JSType]",
		"innerMap j2gort.IdentityOrderedMap[K, *LinkedIdentityEntry[K, V]]",
		"seenKeys j2gort.IdentityOrderedSet[K]",
		"uniqueTypes j2gort.CustomOrderedSet[JSType]",
		"lazyObjs    *j2gort.CustomOrderedSet[*ObjTypeImpl]",
		"nodes j2gort.OrderedMap[N, *LinkedDiGraphNode[N, E]]",
		"var workSet j2gort.OrderedSet[DiGraphNode[N, E]] = j2gort.NewOrderedSet[DiGraphNode[N, E]]()",
		"workSet *j2gort.OrderedSet[DiGraphNode[N, E]])",
		"seenSet j2gort.OrderedSet[T]",
		// Soundness negative guards preserved unlowered:
		"var strIdentMap *j2gort.Ext_java_util_IdentityHashMap = j2gort.NewExt_java_util_IdentityHashMap()",
		"set *j2gort.Ext_java_util_LinkedHashSet",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing expected snippet %q in lowered workload.go:\n%s", want, got)
		}
	}

	diffTestSrc := `package w16_diff_test

import (
	"testing"
	"example.com/w13maptest/orig"
	"example.com/w13maptest/opt"
)

func TestDifferentialParity(t *testing.T) {
	want := orig.RunAllScenarios()
	got := opt.RunAllScenarios()
	if got != want {
		t.Fatalf("differential mismatch:\n got: %s\nwant: %s", got, want)
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "w16_diff_test.go"), []byte(diffTestSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "test", "-v", "-run", "TestDifferentialParity", ".")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Wave 16 Tasks 2 & 3 differential test failed: %v\n%s", err, string(out))
	}
}

func TestWave16Task5_HotArrayListSliceLoweringAndDifferentialParity(t *testing.T) {
	tmpDir, _ := setupWave16J2gortOverlayWorkspace(t)

	buildWorkload := func(pkgName string) string {
		return "package " + pkgName + `

import (
	"fmt"
	"strings"

	"example.com/w13maptest/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
)

type Node = nodestore.Node

type JSType interface {
	IsJavaNil_() bool
	TypeName() string
}

type SimpleType struct {
	Name string
}

func (s *SimpleType) IsJavaNil_() bool { return s == nil }
func (s *SimpleType) TypeName() string { return s.Name }

type TypeStringBuilder struct {
	sb strings.Builder
}

func (b *TypeStringBuilder) Append(s string) { b.sb.WriteString(s) }
func (b *TypeStringBuilder) AppendAll(items j2gort.Ext_java_lang_Iterable, sep string) {
	cur := j2gort.NewIterCursor(j2gort.CheckNotNull(items))
	first := true
	for cur.HasNext() {
		if !first {
			b.sb.WriteString(sep)
		}
		first = false
		b.sb.WriteString(cur.NextStr())
	}
}

// 1. UnionType_Builder.alternates (ImmListCopyToSlice) + UnionType.AppendTo (Collections_Sort_1 + TypeStringBuilder.AppendAll)
type UnionType_Builder struct {
	alternates      j2gort.Ext_java_util_List
	finalAlternates []JSType
}

func NewUnionType_Builder() *UnionType_Builder {
	this_ := &UnionType_Builder{}
	InitUnionType_Builder(this_)
	return this_
}

func InitUnionType_Builder(this_ *UnionType_Builder) {
	this_.alternates = j2gort.NewExt_java_util_ArrayList()
	this_.finalAlternates = nil
}

func (this_ *UnionType_Builder) AddAlternate(alt JSType) *UnionType_Builder {
	if this_.alternates.Size() > 0 && j2gort.AsIface[JSType](this_.alternates.Get(0)).TypeName() == alt.TypeName() {
		this_.alternates.Set(0, alt)
		return this_
	}
	this_.alternates.Add_1(alt)
	return this_
}

func (this_ *UnionType_Builder) BuildAlternates() []JSType {
	if this_.alternates.IsEmpty() {
		this_.finalAlternates = []JSType{}
	} else {
		this_.finalAlternates = j2gort.ImmListCopyToSlice[JSType](j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.alternates))
	}
	return this_.finalAlternates
}

func FormatUnionAlternates(alts []JSType) string {
	sb := &TypeStringBuilder{}
	sb.Append("(")
	var names j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	for _, a := range alts {
		names.Add_1(j2gort.BoxString(a.TypeName()))
	}
	j2gort.Ext_java_util_Collections_Sort_1(names)
	sb.AppendAll(names, "|")
	sb.Append(")")
	return sb.sb.String()
}

// 2. TypedScopeCreator: constructor chaining InitTypedScopeCreator_1 -> InitTypedScopeCreator_2 + Add_2(0, ...) + Joiner
type TypedScopeCreator struct {
	unresolvedNodes j2gort.Ext_java_util_List
}

func NewTypedScopeCreator() *TypedScopeCreator {
	this_ := &TypedScopeCreator{}
	InitTypedScopeCreator_1(this_)
	return this_
}

func InitTypedScopeCreator_1(this_ *TypedScopeCreator) {
	InitTypedScopeCreator_2(this_, true)
}

func InitTypedScopeCreator_2(this_ *TypedScopeCreator, _ bool) {
	this_.priv_TypedScopeCreator_clearCommonState_1()
}

func (this_ *TypedScopeCreator) priv_TypedScopeCreator_clearCommonState_1() {
	this_.unresolvedNodes = j2gort.NewExt_java_util_ArrayList()
}

func (this_ *TypedScopeCreator) SetDeferredType(n Node) {
	this_.unresolvedNodes.Add_1(j2gort.BoxNode(n))
}

func (this_ *TypedScopeCreator) FinishAndBuildPath(parts []string) string {
	var sum int
	{
		cur := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.unresolvedNodes))
		for cur.HasNext() {
			nd := Node(cur.NextNode())
			sum += int(nd)
		}
	}
	var properties j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	for _, p := range parts {
		properties.Add_2(int32(0), j2gort.BoxString(p))
	}
	joined := j2gort.Ext_com_google_common_base_Joiner_On_1_uint16(uint16('.')).Join_1_Ext_java_lang_Iterable(properties)
	return fmt.Sprintf("%d:%s", sum, joined)
}

// 3. RemoveUnusedCode_RemovableBuilder + RemovableBase (if/else mutually exclusive Build + chained Build overload + stack var val_8)
type Continuation struct {
	Tag string
}

func (c *Continuation) IsJavaNil_() bool { return c == nil }

type RemovableBase struct {
	continuations j2gort.Ext_java_util_List
}

func (this_ *RemovableBase) initFields_Removable() {}

func InitRemovable(this_ *RemovableBase, builder *RemovableBuilder) {
	this_.initFields_Removable()
	this_.continuations = builder.continuations
}

func (this_ *RemovableBase) Drain() string {
	var out string
	{
		cur := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.continuations))
		for cur.HasNext() {
			raw := cur.Next()
			c := j2gort.CastPtr[Continuation](raw)
			out += c.Tag + ";"
		}
	}
	this_.continuations.Clear()
	return out
}

type RemovableBuilder struct {
	continuations j2gort.Ext_java_util_List
}

func (this_ *RemovableBuilder) initFields_RemovableBuilder() {
	this_.continuations = j2gort.NewExt_java_util_ArrayList()
}

func NewRemovableBuilder() *RemovableBuilder {
	this_ := &RemovableBuilder{}
	InitRemovableBuilder(this_)
	return this_
}

func InitRemovableBuilder(this_ *RemovableBuilder) {
	this_.initFields_RemovableBuilder()
}

func (this_ *RemovableBuilder) AddContinuation(c *Continuation) *RemovableBuilder {
	this_.continuations.Add_1(j2gort.PtrToAny(c))
	return this_
}

func (this_ *RemovableBuilder) Build_1() *RemovableBase {
	return this_.Build_2(false)
}

func (this_ *RemovableBuilder) Build_2(_ bool) *RemovableBase {
	r := &RemovableBase{}
	InitRemovable(r, this_)
	return r
}

func UseRemovableBuilders(branch bool) string {
	var val_8 RemovableBuilder
	InitRemovableBuilder(&val_8)
	var bStack *RemovableBuilder = &val_8
	bStack.AddContinuation(&Continuation{Tag: "stack"})
	r0 := bStack.Build_1()

	var bBranch *RemovableBuilder = NewRemovableBuilder()
	var r1 *RemovableBase
	if branch {
		r1 = bBranch.Build_2(true)
	} else {
		bBranch.AddContinuation(&Continuation{Tag: "elseBranch"})
		r1 = bBranch.Build_2(false)
	}
	return r0.Drain() + "|" + r1.Drain()
}

// 4. FunctionToBlockMutator (functionsToHoist.Add_2(0, BoxNode(fnVar)))
func HoistFunctions(in []Node) string {
	var functionsToHoist j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	for _, fnVar := range in {
		if fnVar != 0 {
			functionsToHoist.Add_2(int32(0), j2gort.BoxNode(fnVar))
		}
	}
	var out []string
	{
		cur := j2gort.NewIterCursor(functionsToHoist)
		for cur.HasNext() {
			fn := Node(cur.NextNode())
			out = append(out, fmt.Sprintf("%d", fn))
		}
	}
	return strings.Join(out, ",")
}

// 5. PureFunctionIdentifier_CollectCallableLeaves + CollectCallableLeavesInternal + javalang.ToJavaString(status)
func PureFunctionIdentifier_CollectCallableLeaves(root int) j2gort.Ext_com_google_common_collect_ImmutableList {
	var callables j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	var allLegal bool = PureFunctionIdentifier_CollectCallableLeavesInternal(root, callables)
	return func() j2gort.Ext_com_google_common_collect_ImmutableList {
		if allLegal {
			return j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(callables)
		}
		return nil
	}()
}

func PureFunctionIdentifier_CollectCallableLeavesInternal(expr int, results j2gort.Ext_java_util_ArrayList) bool {
	if expr < 0 {
		return false
	}
	if expr <= 2 {
		results.Add_1(j2gort.BoxNode(Node(expr + 100)))
		return true
	}
	return PureFunctionIdentifier_CollectCallableLeavesInternal(expr-1, results) &&
		PureFunctionIdentifier_CollectCallableLeavesInternal(expr-2, results)
}

func SideEffectsToString(mutThis, mutGlobal, throws bool) string {
	var status j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	if mutThis {
		status.Add_1(j2gort.BoxString("this"))
	}
	if mutGlobal {
		status.Add_1(j2gort.BoxString("global"))
	}
	if throws {
		status.Add_1(j2gort.BoxString("throw"))
	}
	return javalang.ToJavaString(status)
}

// 6. General UnmodifiableList (live field view + non-return local), CopyCollection, and Lists_Reverse
type LiveListHolder struct {
	items j2gort.Ext_java_util_List
}

func NewLiveListHolder() *LiveListHolder {
	this_ := &LiveListHolder{}
	InitLiveListHolder(this_)
	return this_
}

func InitLiveListHolder(this_ *LiveListHolder) {
	this_.items = j2gort.NewExt_java_util_ArrayList()
}

func (this_ *LiveListHolder) Add(s string) {
	this_.items.Add_1(j2gort.BoxString(s))
}

func (this_ *LiveListHolder) View() j2gort.Ext_java_util_List {
	return j2gort.Ext_java_util_Collections_UnmodifiableList(this_.items)
}

func CopyReverseAndUnmod(seed j2gort.Ext_java_util_Collection) string {
	var copied j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList_1_Ext_java_util_Collection(seed)
	copied.Add_1(j2gort.BoxString("tail"))
	var revIterOut string
	{
		cur := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.Ext_com_google_common_collect_Lists_Reverse(copied)))
		for cur.HasNext() {
			s := cur.NextStr()
			revIterOut += s + "/"
		}
	}
	copied = j2gort.Ext_com_google_common_collect_Lists_Reverse(copied)
	view := j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_java_util_Collections_UnmodifiableList(copied))
	var clonedFromCand j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList_1_Ext_java_util_Collection(copied)
	return fmt.Sprintf("%s|%d|%s", revIterOut, view.Size(), j2gort.CastString(clonedFromCand.Get(0)))
}

func RunTask5Workload() string {
	ub := NewUnionType_Builder()
	ub.AddAlternate(&SimpleType{Name: "string"})
	ub.AddAlternate(&SimpleType{Name: "number"})
	ub.AddAlternate(&SimpleType{Name: "boolean"})
	alts := ub.BuildAlternates()
	uStr := FormatUnionAlternates(alts)

	tsc := NewTypedScopeCreator()
	tsc.SetDeferredType(Node(7))
	tsc.SetDeferredType(Node(11))
	tscStr := tsc.FinishAndBuildPath([]string{"Bar", "foo", "goog"})

	remStr := UseRemovableBuilders(true) + "#" + UseRemovableBuilders(false)
	hoistStr := HoistFunctions([]Node{1, 0, 2, 3})

	immLeaves := PureFunctionIdentifier_CollectCallableLeaves(4)
	extList := j2gort.NewExt_java_util_ArrayList()
	extOK := PureFunctionIdentifier_CollectCallableLeavesInternal(3, extList)
	seStr := SideEffectsToString(true, false, true)

	lh := NewLiveListHolder()
	lh.Add("x")
	liveView := lh.View()
	lh.Add("y") // live mutation AFTER View() must be visible in liveView.Size()!
	cruStr := CopyReverseAndUnmod(liveView)

	return fmt.Sprintf("u=%s|tsc=%s|rem=%s|hoist=%s|pfi=%d,%v,%d,%s|live=%d,%s",
		uStr, tscStr, remStr, hoistStr, immLeaves.Size(), extOK, extList.Size(), seStr, liveView.Size(), cruStr)
}
`
	}

	origDir := filepath.Join(tmpDir, "orig5")
	optDir := filepath.Join(tmpDir, "opt5")
	if err := os.MkdirAll(origDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(optDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(origDir, "hotlists.go"), []byte(buildWorkload("orig5")), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(optDir, "hotlists.go"), []byte(buildWorkload("opt5")), 0o644); err != nil {
		t.Fatal(err)
	}

	allPasses := passes.AllPasses()
	stats1, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if stats1.ByPass["native_slice_lowering"] == 0 {
		t.Fatalf("expected native_slice_lowering to fire, got %+v", stats1.ByPass)
	}

	// Run 2 MUST be 100% idempotent!
	stats2, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites and 0 files modified on Run 2, got %+v", stats2)
	}

	gotBytes, err := os.ReadFile(filepath.Join(optDir, "hotlists.go"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(gotBytes)
	for _, want := range []string{
		"alternates      []JSType",
		"j2gort.ImmSliceCloneNonNull(this_.alternates)",
		"j2gort.SliceSortNatural(names)",
		"sb.AppendAll(j2gort.SliceReadOnlyIterable(names), \"|\")",
		"unresolvedNodes []Node",
		"properties = j2gort.SliceInsertAt(properties, int32(0), p)",
		"Join_1_Ext_java_lang_Iterable(j2gort.SliceReadOnlyIterable(properties))",
		"continuations []*Continuation",
		"functionsToHoist = j2gort.SliceInsertAt(functionsToHoist, int32(0), fnVar)",
		"func sliceAcc_PureFunctionIdentifier_CollectCallableLeavesInternal(expr int, results *[]Node) bool",
		"sliceAcc_PureFunctionIdentifier_CollectCallableLeavesInternal(root, &callables)",
		"j2gort.SliceFormatJavaString(status)",
		"j2gort.SliceUnmodifiableListLive(&this_.items)",
		"j2gort.CopyStrColToSlice(seed)",
		"j2gort.SliceReverseInPlace(copied)",
		"j2gort.CloneNonNilSlice(copied)",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing expected lowered slice snippet %q in hotlists.go:\n%s", want, got)
		}
	}

	diffTestSrc := `package w16_t5_test

import (
	"testing"
	"example.com/w13maptest/orig5"
	"example.com/w13maptest/opt5"
)

func TestTask5DifferentialParity(t *testing.T) {
	want := orig5.RunTask5Workload()
	got := opt5.RunTask5Workload()
	if got != want {
		t.Fatalf("Task 5 differential mismatch:\n got: %s\nwant: %s", got, want)
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "w16_t5_test.go"), []byte(diffTestSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "test", "-v", "-run", "TestTask5DifferentialParity", ".")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Wave 16 Task 5 differential test failed: %v\n%s", err, string(out))
	}
}

func TestWave19C_StackUnionTypeBuilderPeepholeAndZeroAlloc(t *testing.T) {
	tmpDir, _ := setupWave16J2gortOverlayWorkspace(t)
	extPath := filepath.Join(tmpDir, "j2gort", "ext_bindings.go")
	if rawExt, err := os.ReadFile(extPath); err == nil && !strings.Contains(string(rawExt), "ImmutableList_Of_4(") {
		addOf4 := "\nfunc Ext_com_google_common_collect_ImmutableList_Of_4(arg0 any, arg1 any, arg2 any, arg3 any) Ext_com_google_common_collect_ImmutableList { return wrapImmutableListOwn([]any{CheckNotNull(ToAny(arg0)), CheckNotNull(ToAny(arg1)), CheckNotNull(ToAny(arg2)), CheckNotNull(ToAny(arg3))}) }\n"
		if err := os.WriteFile(extPath, append(rawExt, []byte(addOf4)...), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	buildRhinoUnionWorkload := func(pkgName string) string {
		return "package " + pkgName + `

import (
	"fmt"
	"strings"

	"example.com/w13maptest/j2gort"
)

type JSType interface {
	IsJavaNil_() bool
	TypeName() string
	IsResolved() bool
}

type JSTypeRegistry struct {
	NoType      JSType
	UnknownType JSType
}

func (r *JSTypeRegistry) IsJavaNil_() bool { return r == nil }

type PrimType struct {
	Name     string
	Resolved bool
}

func (p *PrimType) IsJavaNil_() bool { return p == nil }
func (p *PrimType) TypeName() string { return p.Name }
func (p *PrimType) IsResolved() bool { return p.Resolved }

type UnionType struct {
	registry   *JSTypeRegistry
	alternates j2gort.Ext_com_google_common_collect_ImmutableList
}

func (u *UnionType) IsJavaNil_() bool { return u == nil }
func (u *UnionType) IsResolved() bool { return true }
func (u *UnionType) TypeName() string {
	var parts []string
	for i := int32(0); i < u.alternates.Size(); i++ {
		parts = append(parts, j2gort.AsIface[JSType](u.alternates.Get(i)).TypeName())
	}
	return "(" + strings.Join(parts, "|") + ")"
}

func NewUnionType(b *UnionType_Builder) *UnionType {
	u := &UnionType{}
	InitUnionType(u, b)
	return u
}

func InitUnionType(this_ *UnionType, b *UnionType_Builder) {
	this_.registry = b.registry
	this_.alternates = b.finalAlternates
}

func UnionType_S_Builder(reg *JSTypeRegistry) *UnionType_Builder {
	return NewUnionType_Builder_1_JSTypeRegistry(reg)
}

type UnionType_Builder struct {
	rebuildTarget   *UnionType
	registry        *JSTypeRegistry
	alt0            JSType
	alt1            JSType
	alt2            JSType
	alt3            JSType
	alt4            JSType
	alt5            JSType
	altCount        int32
	restAlternates  j2gort.Ext_java_util_ArrayList
	finalAlternates j2gort.Ext_com_google_common_collect_ImmutableList
	hasBuilt        bool
}

func (this_ *UnionType_Builder) IsJavaNil_() bool { return this_ == nil }

func (this_ *UnionType_Builder) initFields_UnionType_Builder() {
	this_.altCount = 0
	this_.restAlternates = nil
	this_.finalAlternates = nil
	this_.hasBuilt = false
}

func NewUnionType_Builder_1_JSTypeRegistry(reg *JSTypeRegistry) *UnionType_Builder {
	this_ := &UnionType_Builder{}
	InitUnionType_Builder_1_JSTypeRegistry(this_, reg)
	return this_
}

func InitUnionType_Builder_1_JSTypeRegistry(this_ *UnionType_Builder, reg *JSTypeRegistry) {
	this_.initFields_UnionType_Builder()
	this_.rebuildTarget = nil
	this_.registry = reg
}

func NewUnionType_Builder_1_UnionType(target *UnionType) *UnionType_Builder {
	this_ := &UnionType_Builder{}
	InitUnionType_Builder_1_UnionType(this_, target)
	return this_
}

func InitUnionType_Builder_1_UnionType(this_ *UnionType_Builder, target *UnionType) {
	this_.initFields_UnionType_Builder()
	this_.rebuildTarget = target
	this_.registry = target.registry
}

func (this_ *UnionType_Builder) getAt(idx int32) JSType {
	switch idx {
	case 0:
		return this_.alt0
	case 1:
		return this_.alt1
	case 2:
		return this_.alt2
	case 3:
		return this_.alt3
	case 4:
		return this_.alt4
	case 5:
		return this_.alt5
	default:
		return j2gort.AsIface[JSType](this_.restAlternates.Get(idx - 6))
	}
}

func (this_ *UnionType_Builder) appendAlt(v JSType) {
	switch this_.altCount {
	case 0:
		this_.alt0 = v
	case 1:
		this_.alt1 = v
	case 2:
		this_.alt2 = v
	case 3:
		this_.alt3 = v
	case 4:
		this_.alt4 = v
	case 5:
		this_.alt5 = v
	default:
		if this_.restAlternates == nil {
			this_.restAlternates = j2gort.NewExt_java_util_ArrayList_1_int32(8)
		}
		this_.restAlternates.Add_1(v)
	}
	this_.altCount++
}

func (this_ *UnionType_Builder) AddAlternate(alt JSType) *UnionType_Builder {
	for i := int32(0); i < this_.altCount; i++ {
		if this_.getAt(i).TypeName() == alt.TypeName() {
			return this_
		}
	}
	this_.appendAlt(alt)
	return this_
}

func (this_ *UnionType_Builder) assignFinalAlternates() {
	if this_.altCount == 1 {
		this_.finalAlternates = j2gort.Ext_com_google_common_collect_ImmutableList_Of_1(this_.alt0)
	} else if this_.altCount == 2 {
		this_.finalAlternates = j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(this_.alt0, this_.alt1)
	} else if this_.altCount == 3 {
		this_.finalAlternates = j2gort.Ext_com_google_common_collect_ImmutableList_Of_3(this_.alt0, this_.alt1, this_.alt2)
	} else if this_.altCount == 4 {
		this_.finalAlternates = j2gort.Ext_com_google_common_collect_ImmutableList_Of_4(this_.alt0, this_.alt1, this_.alt2, this_.alt3)
	} else if this_.altCount == 5 {
		this_.finalAlternates = j2gort.Ext_com_google_common_collect_ImmutableList_Of_5(this_.alt0, this_.alt1, this_.alt2, this_.alt3, this_.alt4)
	} else if this_.altCount == 6 {
		this_.finalAlternates = j2gort.Ext_com_google_common_collect_ImmutableList_Of_6(this_.alt0, this_.alt1, this_.alt2, this_.alt3, this_.alt4, this_.alt5)
	} else {
		var all j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList_1_int32(this_.altCount)
		all.Add_1(this_.alt0)
		all.Add_1(this_.alt1)
		all.Add_1(this_.alt2)
		all.Add_1(this_.alt3)
		all.Add_1(this_.alt4)
		all.Add_1(this_.alt5)
		for i := int32(6); i < this_.altCount; i++ {
			all.Add_1(this_.restAlternates.Get(i - 6))
		}
		this_.finalAlternates = j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(all)
	}
}

func (this_ *UnionType_Builder) Build() JSType {
	this_.hasBuilt = true
	if this_.rebuildTarget != nil {
		this_.assignFinalAlternates()
		this_.rebuildTarget.alternates = this_.finalAlternates
		if this_.finalAlternates.Size() == 1 {
			return j2gort.AsIface[JSType](this_.finalAlternates.Get(0))
		}
		return this_.rebuildTarget
	}
	if this_.altCount == 0 {
		return this_.registry.NoType
	}
	if this_.altCount == 1 {
		return this_.alt0
	}
	this_.assignFinalAlternates()
	return NewUnionType(this_)
}

func CreateUnionType2(reg *JSTypeRegistry, a, b JSType) JSType {
	var builder *UnionType_Builder = UnionType_S_Builder(reg)
	_ = builder
	builder.AddAlternate(a)
	builder.AddAlternate(b)
	return builder.Build()
}

func CreateUnionChained(reg *JSTypeRegistry, a, b, c JSType) JSType {
	var builder *UnionType_Builder = UnionType_S_Builder(reg)
	return builder.AddAlternate(a).AddAlternate(b).AddAlternate(c).Build()
}

func RebuildUnion(u *UnionType, extra JSType) JSType {
	var b *UnionType_Builder = NewUnionType_Builder_1_UnionType(u)
	_ = b
	for i := int32(0); i < u.alternates.Size(); i++ {
		b.AddAlternate(j2gort.AsIface[JSType](u.alternates.Get(i)))
	}
	b.AddAlternate(extra)
	return b.Build()
}

// Soundness negative guards: NONE of these escaping builders may be stack-allocated!
func EscapingReturnDirect(reg *JSTypeRegistry) *UnionType_Builder {
	var escDirect *UnionType_Builder = UnionType_S_Builder(reg)
	return escDirect
}

func EscapingReturnChained(reg *JSTypeRegistry, a JSType) *UnionType_Builder {
	var escChained *UnionType_Builder = UnionType_S_Builder(reg)
	return escChained.AddAlternate(a)
}

func EscapingClosure(reg *JSTypeRegistry, a JSType) func() JSType {
	var escClos *UnionType_Builder = UnionType_S_Builder(reg)
	return func() JSType {
		escClos.AddAlternate(a)
		return escClos.Build()
	}
}

func sinkBuilder(b *UnionType_Builder) *UnionType_Builder { return b }

func EscapingArg(reg *JSTypeRegistry) *UnionType_Builder {
	var escArg *UnionType_Builder = UnionType_S_Builder(reg)
	return sinkBuilder(escArg)
}

func RunUnionScenarios() string {
	reg := &JSTypeRegistry{
		NoType:      &PrimType{Name: "NoType", Resolved: true},
		UnknownType: &PrimType{Name: "?", Resolved: true},
	}
	types := make([]JSType, 8)
	for i := range types {
		types[i] = &PrimType{Name: fmt.Sprintf("T%d", i), Resolved: true}
	}

	uCollapse := CreateUnionType2(reg, types[0], types[0])
	u2 := CreateUnionType2(reg, types[0], types[1])
	uChain := CreateUnionChained(reg, types[1], types[2], types[3])

	var loopNames []string
	for n := 4; n <= 8; n++ {
		var loopBuilder *UnionType_Builder = UnionType_S_Builder(reg)
		for k := 0; k < n; k++ {
			loopBuilder.AddAlternate(types[k])
		}
		loopNames = append(loopNames, loopBuilder.Build().TypeName())
	}

	rebuilt := RebuildUnion(u2.(*UnionType), types[4])
	e1 := EscapingReturnDirect(reg).AddAlternate(types[0]).Build().TypeName()
	e2 := EscapingReturnChained(reg, types[2]).Build().TypeName()
	e3 := EscapingClosure(reg, types[3])().TypeName()
	e4 := EscapingArg(reg).Build().TypeName()

	return fmt.Sprintf("%s|%s|%s|%s|%s|%s,%s,%s,%s",
		uCollapse.TypeName(), u2.TypeName(), uChain.TypeName(),
		strings.Join(loopNames, ";"), rebuilt.TypeName(), e1, e2, e3, e4)
}
`
	}

	origDir := filepath.Join(tmpDir, "orig19c")
	optDir := filepath.Join(tmpDir, "opt19c")
	if err := os.MkdirAll(origDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(optDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(origDir, "union_type.go"), []byte(buildRhinoUnionWorkload("orig19c")), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(optDir, "union_type.go"), []byte(buildRhinoUnionWorkload("opt19c")), 0o644); err != nil {
		t.Fatal(err)
	}

	allPasses := passes.AllPasses()
	stats1, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if stats1.ByPass["j2gort_peephole"] == 0 {
		t.Fatalf("expected j2gort_peephole to fire on UnionType_Builder, got %+v", stats1.ByPass)
	}

	// Run 2 MUST be 100% idempotent!
	stats2, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites and 0 files modified on Run 2, got %+v", stats2)
	}

	gotBytes, err := os.ReadFile(filepath.Join(optDir, "union_type.go"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(gotBytes)
	for _, want := range []string{
		"var builder_ubVal UnionType_Builder",
		"InitUnionType_Builder_1_JSTypeRegistry(&builder_ubVal, reg)",
		"var builder *UnionType_Builder = &builder_ubVal",
		"var b_ubVal UnionType_Builder",
		"InitUnionType_Builder_1_UnionType(&b_ubVal, u)",
		"var loopBuilder_ubVal UnionType_Builder",
		// Soundness negative guards MUST remain heap-allocated:
		"var escDirect *UnionType_Builder = UnionType_S_Builder(reg)",
		"var escChained *UnionType_Builder = UnionType_S_Builder(reg)",
		"var escClos *UnionType_Builder = UnionType_S_Builder(reg)",
		"var escArg *UnionType_Builder = UnionType_S_Builder(reg)",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing expected snippet %q in lowered union_type.go:\n%s", want, got)
		}
	}

	diffAndAllocTestSrc := `package w19c_test

import (
	"testing"
	"example.com/w13maptest/orig19c"
	"example.com/w13maptest/opt19c"
)

func TestWave19CDifferentialAndZeroAlloc(t *testing.T) {
	want := orig19c.RunUnionScenarios()
	got := opt19c.RunUnionScenarios()
	if got != want {
		t.Fatalf("Wave 19C differential mismatch:\n got: %s\nwant: %s", got, want)
	}

	reg := &opt19c.JSTypeRegistry{
		NoType:      &opt19c.PrimType{Name: "NoType", Resolved: true},
		UnknownType: &opt19c.PrimType{Name: "?", Resolved: true},
	}
	tStr := &opt19c.PrimType{Name: "string", Resolved: true}
	var sink opt19c.JSType
	allocs := testing.AllocsPerRun(200, func() {
		sink = opt19c.CreateUnionType2(reg, tStr, tStr)
	})
	_ = sink
	if allocs != 0 {
		t.Fatalf("expected 0 heap allocations for collapsing CreateUnionType2, got %v", allocs)
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "w19c_test.go"), []byte(diffAndAllocTestSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "test", "-v", "-run", "TestWave19CDifferentialAndZeroAlloc", ".")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Wave 19C differential & zero-alloc test failed: %v\n%s", err, string(out))
	}
}
