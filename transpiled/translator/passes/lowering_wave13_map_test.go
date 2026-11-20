package passes_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

// extractLoweringHelpersFromJava extracts the exact Go source of the post-pass runtime helpers
// (ImmListFromTypedSlice, SliceIterator, OrderedSet, LocalOrderedSet, OrderedMap, CustomOrderedMap)
// from RuntimeCollectionsBindings.java so tests compile and execute the real emitted runtime code
// against transpiled/gen/j2gort.
func extractLoweringHelpersFromJava(t *testing.T, repoRoot string) string {
	t.Helper()
	javaPath := filepath.Join(repoRoot, "transpiled", "translator", "java", "j2go", "RuntimeCollectionsBindings.java")
	rawBytes, err := os.ReadFile(javaPath)
	if err != nil {
		t.Fatalf("read RuntimeCollectionsBindings.java: %v", err)
	}
	src := string(rawBytes)
	startMarker := "// --- Post-Pass Native Slice, OrderedSet, and LocalStringBuilder Runtime Helpers ---"
	sIdx := strings.Index(src, startMarker)
	if sIdx < 0 {
		t.Fatalf("could not locate lowering helpers section in RuntimeCollectionsBindings.java (sIdx=%d)", sIdx)
	}
	block := src[sIdx:]
	var lines []string
	for _, line := range strings.Split(block, "\n") {
		lt := strings.TrimSpace(line)
		if lt == `""";` || lt == `return """` || strings.HasPrefix(lt, "private static String buildColRt") ||
			(lt == "}" && (line == "}" || strings.HasPrefix(line, "  }"))) {
			continue
		}
		trimmed := strings.TrimPrefix(line, "        ")
		trimmed = strings.ReplaceAll(trimmed, `\t`, "\t")
		trimmed = strings.ReplaceAll(trimmed, `\"`, `"`)
		lines = append(lines, trimmed)
	}
	return "package j2gort\n\nimport (\n\t\"fmt\"\n\t\"math/bits\"\n\t\"reflect\"\n\t\"slices\"\n\t\"sort\"\n\t\"strconv\"\n\t\"strings\"\n\t\"sync\"\n\t\"sync/atomic\"\n\t\"unicode/utf16\"\n\n\t\"github.com/paullepoulpe/jsgo/transpiled/shims/javalang\"\n\t\"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore\"\n)\n\nvar (\n\t_ = bits.TrailingZeros64\n\t_ = reflect.TypeOf\n\t_ = slices.Sort[[]int]\n\t_ = sort.Slice\n\t_ = strconv.AppendInt\n\t_ = utf16.DecodeRune\n\t_ = sync.Pool{}\n\t_ = atomic.Int32{}\n)\n\n" + strings.Join(lines, "\n")
}

func setupRealJ2gortOverlayWorkspace(t *testing.T) (tmpDir, repoRoot string) {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	repoRoot = filepath.Clean(filepath.Join(wd, "..", "..", ".."))
	realJ2gortDir := filepath.Join(repoRoot, "transpiled", "gen", "j2gort")

	tmpDir = t.TempDir()
	goMod := "module example.com/w13maptest\n\ngo 1.24.0\n\nrequire (\n\tgithub.com/paullepoulpe/jsgo v0.0.0\n\tgolang.org/x/text v0.34.0\n)\n\nreplace github.com/paullepoulpe/jsgo => " + filepath.ToSlash(repoRoot) + "\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}
	if sumBytes, err := os.ReadFile(filepath.Join(repoRoot, "go.sum")); err == nil {
		_ = os.WriteFile(filepath.Join(tmpDir, "go.sum"), sumBytes, 0o644)
	}

	dstJ2 := filepath.Join(tmpDir, "j2gort")
	if err := os.MkdirAll(dstJ2, 0o755); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(realJ2gortDir)
	if err != nil {
		t.Fatalf("read real j2gort dir: %v", err)
	}
	startMarker := "// --- Post-Pass Native Slice, OrderedSet, and LocalStringBuilder Runtime Helpers ---"
	needCompPoolResetStub := true
	needRuntimePoolPurgerStub := true
	hasPurgeAllInBase := false
	for _, ent := range entries {
		if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".go") || strings.HasSuffix(ent.Name(), "_test.go") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(realJ2gortDir, ent.Name()))
		if err != nil {
			t.Fatal(err)
		}
		content := string(data)
		if ent.Name() == "ext_collections_rt.go" {
			sIdx := strings.Index(content, startMarker)
			if sIdx >= 0 {
				content = content[:sIdx]
			}
		}
		if strings.Contains(content, "func registerCompPoolReset(") {
			needCompPoolResetStub = false
		}
		if strings.Contains(content, "func registerRuntimePoolPurger(") {
			needRuntimePoolPurgerStub = false
		}
		if strings.Contains(content, "func PurgeAllRuntimePools(") {
			hasPurgeAllInBase = true
		}
		if err := os.WriteFile(filepath.Join(dstJ2, ent.Name()), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	helpersGo := extractLoweringHelpersFromJava(t, repoRoot)
	if needCompPoolResetStub && strings.Contains(helpersGo, "registerCompPoolReset(") && !strings.Contains(helpersGo, "func registerCompPoolReset(") {
		helpersGo += "\nvar (\n\tcompPoolResetMu sync.Mutex\n\tcompPoolResets  []func()\n)\n\nfunc registerCompPoolReset(fn func()) {\n\tcompPoolResetMu.Lock()\n\tcompPoolResets = append(compPoolResets, fn)\n\tcompPoolResetMu.Unlock()\n}\n"
	}
	if needRuntimePoolPurgerStub && strings.Contains(helpersGo, "registerRuntimePoolPurger(") && !strings.Contains(helpersGo, "func registerRuntimePoolPurger(") {
		helpersGo += "\nvar (\n\truntimePoolPurgerMu sync.Mutex\n\truntimePoolPurgers  []func()\n)\n\nfunc registerRuntimePoolPurger(fn func()) {\n\truntimePoolPurgerMu.Lock()\n\truntimePoolPurgers = append(runtimePoolPurgers, fn)\n\truntimePoolPurgerMu.Unlock()\n}\n"
		if !hasPurgeAllInBase && !strings.Contains(helpersGo, "func PurgeAllRuntimePools(") {
			helpersGo += "\nfunc PurgeAllRuntimePools() {\n\tTrimStaleSlabPools(0)\n\truntimePoolPurgerMu.Lock()\n\tfns := runtimePoolPurgers\n\truntimePoolPurgerMu.Unlock()\n\tfor _, fn := range fns {\n\t\tfn()\n\t}\n}\n"
		}
	}
	if err := os.WriteFile(filepath.Join(dstJ2, "lowering_helpers_overlay.go"), []byte(helpersGo), 0o644); err != nil {
		t.Fatal(err)
	}
	return tmpDir, repoRoot
}

func TestWave13OrderedAndCustomOrderedMap_DirectRuntimeSemantics(t *testing.T) {
	tmpDir, _ := setupRealJ2gortOverlayWorkspace(t)

	unitTestSrc := `package w13_rt_test

import (
	"fmt"
	"testing"

	"example.com/w13maptest/j2gort"
)

type ScopedKey struct {
	Name      string
	ScopeRoot uint32
	ForceHash int32
}

func (k *ScopedKey) IsJavaNil_() bool { return k == nil }
func (k *ScopedKey) HashCode() int32 {
	if k.ForceHash != 0 {
		return k.ForceHash
	}
	return int32(len(k.Name))*31 + int32(k.ScopeRoot)
}
func (k *ScopedKey) Equals(other any) bool {
	o, ok := other.(*ScopedKey)
	return ok && o != nil && o.ScopeRoot == k.ScopeRoot && o.Name == k.Name
}
func (k *ScopedKey) String() string {
	return fmt.Sprintf("%s@%d", k.Name, k.ScopeRoot)
}

func TestOrderedMapTransitionsAndClone(t *testing.T) {
	var uninit *j2gort.OrderedMap[string, int32]
	if !uninit.IsNil() {
		t.Fatalf("expected nil *OrderedMap to report IsNil")
	}

	m := j2gort.NewOrderedMap[string, int32]()
	if m.IsNil() || !m.IsEmpty() || m.Size() != 0 {
		t.Fatalf("unexpected initial state")
	}

	// Lazy alloc (<=12) -> indexed (>12)
	for i := 0; i < 24; i++ {
		k := fmt.Sprintf("k%02d", i)
		if prev, ok := m.PutOk(k, int32(i*10)); ok || prev != 0 {
			t.Fatalf("unexpected prev on fresh insert %d", i)
		}
	}
	if m.Size() != 24 {
		t.Fatalf("want size 24, got %d", m.Size())
	}
	// Update in place preserves order
	if prev := m.Put("k03", 999); prev != 30 {
		t.Fatalf("want prev 30, got %d", prev)
	}
	// Remove > half to trigger compaction, then reinsert k01 (must move to end!)
	for i := 0; i < 18; i++ {
		if i == 3 {
			continue
		}
		k := fmt.Sprintf("k%02d", i)
		if _, ok := m.RemoveOk(k); !ok {
			t.Fatalf("expected remove ok for %s", k)
		}
	}
	m.Put("k01", 111)

	cloned := j2gort.OrderedMapFrom(m)
	if !j2gort.OrderedMapEquals(m, cloned, false) || m.HashCodeMap(false) != cloned.HashCodeMap(false) {
		t.Fatalf("cloned OrderedMap not equal to original")
	}

	var keys []string
	for cur := cloned.IterKeys(); cur.HasNext(); {
		keys = append(keys, cur.Next())
	}
	wantKeys := "[k03 k18 k19 k20 k21 k22 k23 k01]"
	if fmt.Sprint(keys) != wantKeys {
		t.Fatalf("got keys %v, want %s", keys, wantKeys)
	}
}

func TestCustomOrderedMapCollisionsNullKeyAndClone(t *testing.T) {
	m := j2gort.NewCustomOrderedMap[*ScopedKey, string]()

	// Null key support
	if prev := j2gort.CustomOrderedMapPutStr(&m, nil, "nullVal"); prev != j2gort.NullString {
		t.Fatalf("expected NullString on first null-key put, got %q", prev)
	}
	if got := j2gort.CustomOrderedMapGetStr(m, nil); got != "nullVal" {
		t.Fatalf("expected nullVal for nil key, got %q", got)
	}

	// Insert 22 keys with deliberate hash collisions AND structural-equality lookup via distinct pointers!
	for i := 0; i < 22; i++ {
		k := &ScopedKey{Name: fmt.Sprintf("v%d", i), ScopeRoot: uint32(i % 3), ForceHash: int32(i % 4)}
		j2gort.CustomOrderedMapPutStr(&m, k, fmt.Sprintf("val%d", i))
	}
	if m.Size() != 23 {
		t.Fatalf("want size 23, got %d", m.Size())
	}

	// Lookup with freshly allocated ScopedKey pointers (tests value equality, not pointer identity!)
	for i := 0; i < 22; i++ {
		probe := &ScopedKey{Name: fmt.Sprintf("v%d", i), ScopeRoot: uint32(i % 3), ForceHash: int32(i % 4)}
		if got := j2gort.CustomOrderedMapGetStr(m, probe); got != fmt.Sprintf("val%d", i) {
			t.Fatalf("lookup %d got %q", i, got)
		}
	}
	missing := &ScopedKey{Name: "missing", ScopeRoot: 1, ForceHash: 1}
	if got := j2gort.CustomOrderedMapGetStr(m, missing); got != j2gort.NullString {
		t.Fatalf("expected NullString for missing key, got %q", got)
	}

	// Clone before compaction (holes == 0 fast memcpy path)
	cloneFast := j2gort.CustomOrderedMapFrom(m)
	if !j2gort.CustomOrderedMapEquals(m, cloneFast) || m.HashCode() != cloneFast.HashCode() {
		t.Fatalf("fast clone mismatch")
	}

	// Remove most entries including null key to exercise bucket unlink + compaction, then clone again
	if old := j2gort.CustomOrderedMapRemoveStr(&m, nil); old != "nullVal" {
		t.Fatalf("remove nil key got %q", old)
	}
	for i := 0; i < 18; i++ {
		probe := &ScopedKey{Name: fmt.Sprintf("v%d", i), ScopeRoot: uint32(i % 3), ForceHash: int32(i % 4)}
		if old := j2gort.CustomOrderedMapRemoveStr(&m, probe); old != fmt.Sprintf("val%d", i) {
			t.Fatalf("remove %d got %q", i, old)
		}
	}
	if m.Size() != 4 {
		t.Fatalf("want size 4 after removals, got %d", m.Size())
	}
	cloneInline := j2gort.CustomOrderedMapFrom(m)
	if !j2gort.CustomOrderedMapEquals(m, cloneInline) {
		t.Fatalf("inline clone after compaction mismatch")
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "w13_rt_test.go"), []byte(unitTestSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "test", "-count=1", "-race", "./...")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("runtime unit test failed: %v\n%s", err, string(out))
	}
}

func TestWave13NativeMapLowering_HotspotsDifferentialParityAndIdempotence(t *testing.T) {
	tmpDir, _ := setupRealJ2gortOverlayWorkspace(t)

	workloadSrc := `package orig

import "example.com/w13maptest/j2gort"

type Node uint32

func (n Node) IsJavaNil_() bool  { return n == 0 }
func (n Node) RawNodeId_() uint32 { return uint32(n) }

// --- 1. AbstractScope-style nullable lazy LinkedHashMap<String, AbstractVar> ---
type AbstractVar interface {
	IsJavaNil_() bool
	GetName() string
	GetScopeRoot() Node
	Equals(any) bool
	HashCode() int32
}

type Var struct {
	name      string
	scopeRoot Node
}

func strHash(s string) int32 {
	var h int32
	for i := 0; i < len(s); i++ {
		h = 31*h + int32(s[i])
	}
	return h
}

func NewVar(name string, root Node) *Var { return &Var{name: name, scopeRoot: root} }
func (v *Var) IsJavaNil_() bool           { return v == nil }
func (v *Var) GetName() string            { return v.name }
func (v *Var) GetScopeRoot() Node         { return v.scopeRoot }
func (v *Var) HashCode() int32            { return strHash(v.name)*31 + int32(v.scopeRoot) }
func (v *Var) Equals(other any) bool {
	o, ok := other.(AbstractVar)
	return ok && o != nil && !o.IsJavaNil_() && o.GetScopeRoot() == v.scopeRoot && o.GetName() == v.name
}
func (v *Var) String() string { return v.name }

type Scope struct {
	parent *Scope
	var0   AbstractVar
	vars   *j2gort.Ext_java_util_LinkedHashMap
}

func NewScope(parent *Scope) *Scope {
	s := &Scope{parent: parent}
	s.initFields_Scope()
	return s
}

func (this_ *Scope) initFields_Scope() {
	this_.var0 = nil
	this_.vars = nil
}

func (this_ *Scope) GetScopeBase() *Scope { return this_ }

func (this_ *Scope) Declare(name string, v AbstractVar) {
	if this_.vars != nil {
		j2gort.MapPutStr(this_.vars, name, v)
		return
	}
	if this_.var0 == nil {
		this_.var0 = v
		return
	}
	this_.vars = j2gort.NewExt_java_util_LinkedHashMap_1_int32(int32(8))
	j2gort.MapPutStr(this_.vars, this_.var0.GetName(), this_.var0)
	j2gort.MapPutStr(this_.vars, name, v)
	this_.var0 = nil
}

func (this_ *Scope) GetOwnVar(name string) AbstractVar {
	if this_.vars != nil {
		return j2gort.AsIface[AbstractVar](j2gort.MapGetStr(this_.vars, name))
	}
	if this_.var0 != nil && this_.var0.GetName() == name {
		return this_.var0
	}
	return nil
}

func (this_ *Scope) Undeclare(v AbstractVar) {
	if this_.vars != nil {
		j2gort.MapRemoveStr(this_.vars, v.GetName())
	} else if this_.var0 != nil && this_.var0.GetName() == v.GetName() {
		this_.var0 = nil
	}
}

func (this_ *Scope) VarCount() int32 {
	if this_.vars != nil {
		return this_.vars.Size()
	}
	if this_.var0 == nil {
		return 0
	}
	return 1
}

func (this_ *Scope) ClearVars() {
	this_.var0 = nil
	if (this_.vars != nil) && !(this_.vars.IsEmpty()) {
		this_.vars.Clear()
	}
}

func (this_ *Scope) GetVarIterable() j2gort.Ext_java_lang_Iterable {
	if this_.vars != nil {
		return j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.vars.Values())
	}
	if this_.var0 == nil {
		return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_java_util_Collections_EmptyList())
	}
	return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_java_util_Collections_SingletonList(this_.var0))
}

func (this_ *Scope) AllAccessibleNames() string {
	var acc j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	for s := this_; s != nil; s = s.parent {
		if s.GetScopeBase().vars != nil {
			cur := j2gort.NewMapValCursor(s.GetScopeBase().vars)
			for cur.HasNext() {
				raw := cur.Next()
				v := j2gort.AsIface[AbstractVar](raw)
				acc.PutIfAbsent(j2gort.BoxString(v.GetName()), j2gort.BoxString(v.GetName()))
			}
		} else if s.var0 != nil {
			acc.PutIfAbsent(j2gort.BoxString(s.var0.GetName()), j2gort.BoxString(s.var0.GetName()))
		}
	}
	var out string
	curK := j2gort.NewMapKeyCursor(acc)
	for curK.HasNext() {
		rawK := curK.Next()
		if out != "" {
			out += ","
		}
		out += j2gort.CastString(rawK)
	}
	return out
}

// --- 2. MustBeReachingVariableDef_MustDef + MustDefJoin (PutAll clone, resultMap alias, ForEach_1, EntryCursor SetValue(nil), Equals/HashCode) ---
type Definition struct {
	node           Node
	invalidatedBy  *Var
}

func (d *Definition) IsJavaNil_() bool { return d == nil }
func (d *Definition) HashCode() int32  { return int32(d.node) }
func (d *Definition) Equals(other any) bool {
	o, ok := other.(*Definition)
	return ok && o != nil && o.node == d.node
}

type MustDef struct {
	reachingDef *j2gort.Ext_java_util_LinkedHashMap
}

func NewMustDef() *MustDef {
	m := &MustDef{}
	m.initFields_MustDef()
	return m
}

func NewMustDefCopy(other *MustDef) *MustDef {
	m := &MustDef{}
	InitMustDefCopy(m, other)
	return m
}

func (this_ *MustDef) initFields_MustDef() {
	this_.reachingDef = j2gort.NewExt_java_util_LinkedHashMap()
}

func InitMustDefCopy(this_ *MustDef, other *MustDef) {
	this_.initFields_MustDef()
	this_.reachingDef.PutAll(j2gort.AsIface[j2gort.Ext_java_util_Map](other.reachingDef))
}

func (this_ *MustDef) RecordDef(v *Var, node Node, dep *Var) {
	{
		cur := j2gort.NewMapEntryCursor(this_.reachingDef)
		for cur.NextEntry() {
			otherDef := j2gort.CastPtr[Definition](cur.GetValue())
			if otherDef == nil {
				continue
			}
			if otherDef.invalidatedBy != nil && otherDef.invalidatedBy.Equals(v) {
				cur.SetValue(nil)
			}
		}
	}
	if node == 0 {
		this_.reachingDef.Put_2(j2gort.PtrToAny(v), nil)
	} else {
		this_.reachingDef.Put_2(j2gort.PtrToAny(v), &Definition{node: node, invalidatedBy: dep})
	}
}

func (this_ *MustDef) GetDefNode(v *Var) Node {
	d := j2gort.CastPtr[Definition](this_.reachingDef.Get(j2gort.PtrToAny(v)))
	if d == nil {
		return 0
	}
	return d.node
}

func (this_ *MustDef) EqualsMustDef(other *MustDef) bool {
	return other.reachingDef.Equals(j2gort.PtrToAny(this_.reachingDef))
}

func (this_ *MustDef) MustDefHash() int32 {
	return this_.reachingDef.HashCode()
}

type MustDefJoin struct {
	result    *MustDef
	resultMap *j2gort.Ext_java_util_LinkedHashMap
}

func NewMustDefJoin() *MustDefJoin {
	j := &MustDefJoin{}
	j.initFields_MustDefJoin()
	return j
}

func (this_ *MustDefJoin) initFields_MustDefJoin() {
	this_.result = NewMustDef()
	this_.resultMap = this_.result.reachingDef
}

func (this_ *MustDefJoin) JoinFlow(input *MustDef) {
	input.reachingDef.ForEach_1(func() j2gort.Ext_java_util_function_BiConsumer {
		bound := this_
		return j2gort.NewLambda_Ext_java_util_function_BiConsumer(func(mrp0_ any, mrp1_ any) {
			bound.priv_mergeVarDef(j2gort.CastPtr[Var](mrp0_), j2gort.CastPtr[Definition](mrp1_))
		})
	}())
}

func (this_ *MustDefJoin) priv_mergeVarDef(v *Var, def *Definition) {
	var resultDef *Definition = nil
	if def == nil {
		resultDef = nil
	} else if !(this_.resultMap.ContainsKey(j2gort.PtrToAny(v))) {
		resultDef = def
	} else if def.Equals(j2gort.PtrToAny(j2gort.CastPtr[Definition](this_.resultMap.Get(j2gort.PtrToAny(v))))) {
		return
	} else {
		resultDef = nil
	}
	this_.resultMap.Put_2(j2gort.PtrToAny(v), j2gort.PtrToAny(resultDef))
}

func (this_ *MustDefJoin) Finish() *MustDef { return this_.result }

// --- 3. ReferenceCollector-style ComputeIfAbsent IIFE + lambdaSingleton + LiveKeySet ---
type RefList struct {
	Count int32
}

func (*RefList) HasIdentityEquals_() {}
func (r *RefList) IsJavaNil_() bool  { return r == nil }

var lambdaSingleton_ref j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_ref = j2gort.NewLambda_Ext_java_util_function_Function(func(lp any) any {
		_ = j2gort.CastPtr[Var](lp)
		return &RefList{}
	})
}

type RefCollector struct {
	referenceMap j2gort.Ext_java_util_Map
	renames      j2gort.Ext_java_util_Map
	prefix       string
}

func NewRefCollector(prefix string) *RefCollector {
	rc := &RefCollector{prefix: prefix}
	rc.initFields_RefCollector()
	return rc
}

func (this_ *RefCollector) initFields_RefCollector() {
	this_.referenceMap = j2gort.NewExt_java_util_LinkedHashMap()
	this_.renames = j2gort.NewExt_java_util_LinkedHashMap()
}

func (this_ *RefCollector) AddRef(v *Var) int32 {
	var col *RefList = func() *RefList {
		m_1 := this_.referenceMap
		k_2 := j2gort.PtrToAny(v)
		if v_3, ok_4 := j2gort.MapTryGetComputeAny(m_1, k_2); ok_4 {
			return j2gort.CastPtr[RefList](v_3)
		}
		return j2gort.CastPtr[RefList](m_1.ComputeIfAbsent(k_2, lambdaSingleton_ref))
	}()
	col.Count++
	j2gort.MapComputeIfAbsentStr(this_.renames, v.GetName(), func() j2gort.Ext_java_util_function_Function {
		bound := this_
		return j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
			return j2gort.BoxString(bound.prefix + "_" + j2gort.CastString(mrp0_))
		})
	}())
	return col.Count
}

func (this_ *RefCollector) GetRename(name string) string {
	return j2gort.CastString(j2gort.MapGetStr(this_.renames, name))
}

func (this_ *RefCollector) GetRenameOrDefault(name, def string) string {
	return j2gort.CastString(j2gort.MapGetOrDefaultStr(this_.renames, name, j2gort.BoxString(def)))
}

func (this_ *RefCollector) AllSymbolNames() string {
	ks := j2gort.AsIface[j2gort.Ext_java_util_Set](this_.referenceMap.KeySet_0())
	cur := j2gort.NewIterCursor(j2gort.CheckNotNull(ks))
	var out string
	for cur.HasNext() {
		v := j2gort.CastPtr[Var](cur.Next())
		out += v.GetName() + ":" + this_.GetRename(v.GetName()) + ";"
	}
	return out
}

// --- 4. IdentityHashMap vs value-equality & Negative soundness guards ---
type IdentBox struct {
	Tag string
}

func (b *IdentBox) IsJavaNil_() bool  { return b == nil }
func (b *IdentBox) HashCode() int32   { return strHash(b.Tag) }
func (b *IdentBox) Equals(o any) bool { x, ok := o.(*IdentBox); return ok && x != nil && x.Tag == b.Tag }

func ExerciseIdentityAndNodeMaps(a, b *IdentBox, nodes []Node) string {
	var idMap *j2gort.Ext_java_util_IdentityHashMap = j2gort.NewExt_java_util_IdentityHashMap()
	idMap.Put_2(j2gort.PtrToAny(a), j2gort.BoxString("first"))
	idMap.Put_2(j2gort.PtrToAny(b), j2gort.BoxString("second"))

	var nodeMap j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	for i, n := range nodes {
		j2gort.MapPutNode(nodeMap, n, j2gort.ValueOfInt(int32(i+1)))
	}
	imm := j2gort.Ext_com_google_common_collect_ImmutableMap_CopyOf_1_Ext_java_util_Map(nodeMap)

	var sum int32
	cur := j2gort.NewMapEntryCursor(nodeMap)
	for cur.NextEntry() {
		kn := Node(cur.KeyNode())
		vi := j2gort.UnboxInt(j2gort.CastPtr[j2gort.BoxInt](cur.GetValue()))
		sum += int32(kn) * vi
	}

	// Negative Guard 1: erased any key (LinkedDirectedGraph.nodes pattern) MUST NOT lower!
	var erasedMap j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	var anyKey any = j2gort.PtrToAny(a)
	erasedMap.Put_2(anyKey, j2gort.BoxString("erased"))

	// Negative Guard 2: mixed IdentityHashMap + LinkedHashMap on same var MUST NOT lower!
	var mixedMap j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_IdentityHashMap()
	if len(nodes) == 999 {
		mixedMap = j2gort.NewExt_java_util_LinkedHashMap()
	}
	mixedMap.Put_2(j2gort.PtrToAny(a), j2gort.BoxString("m"))

	return j2gort.CastString(idMap.Get(j2gort.PtrToAny(a))) + "," +
		j2gort.CastString(idMap.Get(j2gort.PtrToAny(b))) + "," +
		j2gort.StrVal(j2gort.ValueOfInt(idMap.Size())) + "," +
		j2gort.StrVal(j2gort.ValueOfInt(imm.Size())) + "," +
		j2gort.StrVal(j2gort.ValueOfInt(sum)) + "," +
		j2gort.CastString(erasedMap.Get(anyKey)) + "," +
		j2gort.CastString(mixedMap.Get(j2gort.PtrToAny(a)))
}
`

	origDir := filepath.Join(tmpDir, "orig")
	optDir := filepath.Join(tmpDir, "opt")
	if err := os.MkdirAll(origDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(optDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(origDir, "workload.go"), []byte(workloadSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	optWorkload := strings.Replace(workloadSrc, "package orig", "package opt", 1)
	if err := os.WriteFile(filepath.Join(optDir, "workload.go"), []byte(optWorkload), 0o644); err != nil {
		t.Fatal(err)
	}

	allPasses := passes.AllPasses()
	stats1, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("RunOnPackageDir(opt) failed: %v", err)
	}
	if stats1.ByPass["native_map_lowering"] < 30 {
		t.Fatalf("expected >= 30 native_map_lowering sites, got %+v", stats1.ByPass)
	}

	optBytes, err := os.ReadFile(filepath.Join(optDir, "workload.go"))
	if err != nil {
		t.Fatal(err)
	}
	gotOpt := string(optBytes)
	for _, want := range []string{
		"vars   *j2gort.OrderedMap[string, AbstractVar]",
		"this_.vars = j2gort.NewOrderedMapPtrCap[string, AbstractVar](int32(8))",
		"if this_.vars != nil {",
		"this_.vars.LiveValues()",
		"var acc j2gort.OrderedMap[string, string] = j2gort.NewOrderedMap[string, string]()",
		"reachingDef j2gort.CustomOrderedMap[*Var, *Definition]",
		"resultMap *j2gort.CustomOrderedMap[*Var, *Definition]",
		"this_.resultMap = &this_.result.reachingDef",
		"this_.reachingDef.PutAllMap(other.reachingDef)",
		"input.reachingDef.ForEach(func(mrp0_ *Var, mrp1_ *Definition)",
		"j2gort.CustomOrderedMapEquals(other.reachingDef, this_.reachingDef)",
		"referenceMap j2gort.CustomOrderedMap[*Var, *RefList]",
		"j2gort.CustomOrderedMapComputeIfAbsentFn(&this_.referenceMap, v, nil, lambdaSingleton_ref)",
		"renames      j2gort.OrderedMap[string, string]",
		"j2gort.OrderedMapComputeIfAbsent(&this_.renames, v.GetName(), j2gort.NullString",
		"var idMap j2gort.IdentityOrderedMap[*IdentBox, string] = j2gort.NewIdentityOrderedMap[*IdentBox, string]()",
		"var nodeMap j2gort.OrderedMap[Node, *j2gort.BoxInt] = j2gort.NewOrderedMap[Node, *j2gort.BoxInt]()",
		"nodeMap.ToImmutableMap()",
		// Negative soundness guards preserved untouched:
		"var erasedMap j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()",
		"var mixedMap j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_IdentityHashMap()",
	} {
		if !strings.Contains(gotOpt, want) {
			t.Errorf("missing expected lowered snippet %q in opt/workload.go:\n%s", want, gotOpt)
		}
	}

	// Run 2 MUST produce strictly 0 sites and 0 modified files!
	stats2, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites and 0 modified files on Run 2, got %+v", stats2)
	}

	// Differential runtime parity check between unlowered orig and lowered opt!
	diffTestSrc := `package w13_diff_test

import (
	"fmt"
	"testing"

	"example.com/w13maptest/j2gort"
	"example.com/w13maptest/opt"
	"example.com/w13maptest/orig"
)

func TestHotspotsDifferentialParity(t *testing.T) {
	// 1. Scope (lazy nil -> cap(8) LinkedHashMap<String, AbstractVar>, Undeclare, GetVarIterable, AllAccessibleNames)
	oParent, xParent := orig.NewScope(nil), opt.NewScope(nil)
	oChild, xChild := orig.NewScope(oParent), opt.NewScope(xParent)
	for i := 0; i < 16; i++ {
		nm := fmt.Sprintf("var_%d", i)
		oParent.Declare(nm, orig.NewVar(nm, 1))
		xParent.Declare(nm, opt.NewVar(nm, 1))
		if i%2 == 0 {
			oChild.Declare(nm, orig.NewVar(nm, 2))
			xChild.Declare(nm, opt.NewVar(nm, 2))
		}
	}
	oParent.Undeclare(orig.NewVar("var_3", 1))
	xParent.Undeclare(opt.NewVar("var_3", 1))
	if oParent.VarCount() != xParent.VarCount() || oChild.AllAccessibleNames() != xChild.AllAccessibleNames() {
		t.Fatalf("Scope mismatch: (%d, %q) vs (%d, %q)", oParent.VarCount(), oChild.AllAccessibleNames(), xParent.VarCount(), xChild.AllAccessibleNames())
	}
	var oIterNames, xIterNames []string
	for cur := j2gort.NewIterCursor(oParent.GetVarIterable()); cur.HasNext(); {
		oIterNames = append(oIterNames, cur.Next().(orig.AbstractVar).GetName())
	}
	for cur := j2gort.NewIterCursor(xParent.GetVarIterable()); cur.HasNext(); {
		xIterNames = append(xIterNames, cur.Next().(opt.AbstractVar).GetName())
	}
	if fmt.Sprint(oIterNames) != fmt.Sprint(xIterNames) {
		t.Fatalf("GetVarIterable mismatch: %v vs %v", oIterNames, xIterNames)
	}
	oChild.ClearVars()
	xChild.ClearVars()
	if oChild.VarCount() != xChild.VarCount() {
		t.Fatalf("ClearVars mismatch")
	}

	// 2. MustDef + MustDefJoin (copy ctor, CustomOrderedMap with ScopedName equality, SetValue(nil), JoinFlow, Equals, HashCode)
	oBase, xBase := orig.NewMustDef(), opt.NewMustDef()
	for i := 0; i < 12; i++ {
		nm := fmt.Sprintf("v%d", i)
		var oDep *orig.Var
		var xDep *opt.Var
		if i > 0 {
			oDep = orig.NewVar("v0", 10)
			xDep = opt.NewVar("v0", 10)
		}
		oBase.RecordDef(orig.NewVar(nm, 10), orig.Node(100+i), oDep)
		xBase.RecordDef(opt.NewVar(nm, 10), opt.Node(100+i), xDep)
	}
	oBranchA, xBranchA := orig.NewMustDefCopy(oBase), opt.NewMustDefCopy(xBase)
	oBranchB, xBranchB := orig.NewMustDefCopy(oBase), opt.NewMustDefCopy(xBase)

	// Reassigning v0 invalidates definitions depending on v0 via EntryCursor.SetValue(nil)!
	oBranchA.RecordDef(orig.NewVar("v0", 10), orig.Node(200), nil)
	xBranchA.RecordDef(opt.NewVar("v0", 10), opt.Node(200), nil)
	oBranchB.RecordDef(orig.NewVar("v5", 10), orig.Node(999), nil)
	xBranchB.RecordDef(opt.NewVar("v5", 10), opt.Node(999), nil)

	oJoin, xJoin := orig.NewMustDefJoin(), opt.NewMustDefJoin()
	oJoin.JoinFlow(oBranchA)
	oJoin.JoinFlow(oBranchB)
	xJoin.JoinFlow(xBranchA)
	xJoin.JoinFlow(xBranchB)
	oMerged, xMerged := oJoin.Finish(), xJoin.Finish()

	if oMerged.MustDefHash() != xMerged.MustDefHash() {
		t.Fatalf("MustDefHash mismatch: %d vs %d", oMerged.MustDefHash(), xMerged.MustDefHash())
	}
	if oMerged.EqualsMustDef(orig.NewMustDefCopy(oMerged)) != xMerged.EqualsMustDef(opt.NewMustDefCopy(xMerged)) {
		t.Fatalf("EqualsMustDef copy mismatch")
	}
	for i := 0; i < 12; i++ {
		nm := fmt.Sprintf("v%d", i)
		on := uint32(oMerged.GetDefNode(orig.NewVar(nm, 10)))
		xn := uint32(xMerged.GetDefNode(opt.NewVar(nm, 10)))
		if on != xn {
			t.Fatalf("GetDefNode(%s) mismatch: orig=%d vs opt=%d", nm, on, xn)
		}
	}

	// 3. RefCollector ComputeIfAbsent + NullString vs "" semantics + LiveKeySet
	orc, xrc := orig.NewRefCollector("p"), opt.NewRefCollector("p")
	for round := 0; round < 3; round++ {
		for i := 0; i < 10; i++ {
			nm := fmt.Sprintf("sym%d", i)
			c1 := orc.AddRef(orig.NewVar(nm, 1))
			c2 := xrc.AddRef(opt.NewVar(nm, 1))
			if c1 != c2 {
				t.Fatalf("AddRef(%s) round %d: %d vs %d", nm, round, c1, c2)
			}
		}
	}
	if orc.GetRename("missing") != xrc.GetRename("missing") || xrc.GetRename("missing") != j2gort.NullString {
		t.Fatalf("missing GetRename NullString mismatch: %q vs %q", orc.GetRename("missing"), xrc.GetRename("missing"))
	}
	if orc.GetRenameOrDefault("missing", "def") != xrc.GetRenameOrDefault("missing", "def") {
		t.Fatalf("GetRenameOrDefault mismatch")
	}
	if orc.AllSymbolNames() != xrc.AllSymbolNames() {
		t.Fatalf("AllSymbolNames mismatch: %q vs %q", orc.AllSymbolNames(), xrc.AllSymbolNames())
	}

	// 4. IdentityHashMap distinct pointers with equal fields + Node map
	oa1, oa2 := &orig.IdentBox{Tag: "same"}, &orig.IdentBox{Tag: "same"}
	xa1, xa2 := &opt.IdentBox{Tag: "same"}, &opt.IdentBox{Tag: "same"}
	goOrig := orig.ExerciseIdentityAndNodeMaps(oa1, oa2, []orig.Node{1, 2, 3, 4, 5, 6})
	goOpt := opt.ExerciseIdentityAndNodeMaps(xa1, xa2, []opt.Node{1, 2, 3, 4, 5, 6})
	if goOrig != goOpt || len(goOpt) < 12 || goOpt[:12] != "first,second" {
		t.Fatalf("ExerciseIdentityAndNodeMaps mismatch: orig=%q vs opt=%q", goOrig, goOpt)
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "w13_diff_test.go"), []byte(diffTestSrc), 0o644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command("go", "test", "-count=1", "-race", "./...")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("differential runtime parity test failed: %v\n%s", err, string(out))
	}
}

func TestWave13NativeMapLowering_RealGeneratedPackages(t *testing.T) {
	tmpDir, repoRoot := setupRealJ2gortOverlayWorkspace(t)

	copyGenPkg := func(pkgName string) string {
		srcDir := filepath.Join(repoRoot, "transpiled", "gen", pkgName)
		dstDir := filepath.Join(tmpDir, pkgName)
		if err := os.MkdirAll(dstDir, 0o755); err != nil {
			t.Fatal(err)
		}
		entries, err := os.ReadDir(srcDir)
		if err != nil {
			t.Fatal(err)
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
				continue
			}
			b, err := os.ReadFile(filepath.Join(srcDir, e.Name()))
			if err != nil {
				t.Fatal(err)
			}
			rewritten := strings.ReplaceAll(string(b),
				`"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"`,
				`"example.com/w13maptest/j2gort"`,
			)
			if err := os.WriteFile(filepath.Join(dstDir, e.Name()), []byte(rewritten), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return dstDir
	}

	graphDir := copyGenPkg("graph")
	sourcemapDir := copyGenPkg("sourcemap")

	allPasses := passes.AllPasses()

	// 1. Real `graph` package: `LinkedDirectedGraphBase[N comparable, E any].nodes` lowers to value `j2gort.OrderedMap`.
	graphStats1, err := passes.RunOnPackageDir(graphDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("RunOnPackageDir(graph) failed: %v", err)
	}
	ldgBytes, err := os.ReadFile(filepath.Join(graphDir, "linked_directed_graph.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(ldgBytes), "j2gort.OrderedMap[N, *LinkedDirectedGraph_LinkedDiGraphNode[N, E]]") {
		t.Fatalf("expected LinkedDirectedGraphBase.nodes to use value j2gort.OrderedMap, stats=%+v", graphStats1.ByPass)
	}
	graphStats2, err := passes.RunOnPackageDir(graphDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 on graph failed: %v", err)
	}
	if graphStats2.TotalSites != 0 || graphStats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites on graph Run 2, got %+v", graphStats2)
	}

	// 2. Real `sourcemap` package: `SourceMapGeneratorV3` maps lower cleanly, compile, and are Run-2 idempotent!
	smStats1, err := passes.RunOnPackageDir(sourcemapDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("RunOnPackageDir(sourcemap) failed: %v", err)
	}
	smV3Bytes, err := os.ReadFile(filepath.Join(sourcemapDir, "source_map_generator_v3.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(smV3Bytes), "OrderedMap[") {
		t.Fatalf("expected OrderedMap in real sourcemap package after lowering, stats=%+v", smStats1.ByPass)
	}
	smStats2, err := passes.RunOnPackageDir(sourcemapDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 on sourcemap failed: %v", err)
	}
	if smStats2.TotalSites != 0 || smStats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites on sourcemap Run 2, got %+v", smStats2)
	}

	cmd := exec.Command("go", "build", "./graph", "./sourcemap")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build ./graph ./sourcemap failed after lowering: %v\n%s", err, string(out))
	}
}
