package passes_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

// setupWave17J2gortOverlayWorkspace builds a temporary workspace overlaid with the Wave 17
// runtime helpers AND the updated JIterCursor (kind=10 jTypedViewCursorSource) extracted
// directly from RuntimeCollectionsBindings.java.
func setupWave17J2gortOverlayWorkspace(t *testing.T) (tmpDir, repoRoot string) {
	t.Helper()
	tmpDir, repoRoot = setupWave16J2gortOverlayWorkspace(t)

	rtPath := filepath.Join(tmpDir, "j2gort", "ext_collections_rt.go")
	rtBytes, err := os.ReadFile(rtPath)
	if err != nil {
		t.Fatal(err)
	}
	rtContent := string(rtBytes)
	if strings.Contains(rtContent, "jTypedViewCursorSource") {
		return tmpDir, repoRoot
	}

	javaPath := filepath.Join(repoRoot, "transpiled", "translator", "java", "j2go", "RuntimeCollectionsBindings.java")
	rawJava, err := os.ReadFile(javaPath)
	if err != nil {
		t.Fatalf("read RuntimeCollectionsBindings.java: %v", err)
	}
	javaSrc := string(rawJava)

	icStart := "        type JIterCursor struct {"
	icEnd := "        func newJavaSet2(e0, e1 any) Ext_com_google_common_collect_ImmutableSet {"
	sIdx := strings.Index(javaSrc, icStart)
	eIdx := strings.Index(javaSrc, icEnd)
	if sIdx < 0 || eIdx <= sIdx {
		t.Fatalf("could not locate JIterCursor block in RuntimeCollectionsBindings.java (sIdx=%d, eIdx=%d)", sIdx, eIdx)
	}
	var icLines []string
	for _, line := range strings.Split(javaSrc[sIdx:eIdx], "\n") {
		lt := strings.TrimSpace(line)
		if lt == `""";` || lt == `return """` || strings.HasPrefix(lt, "private static String buildColRt") || (lt == "}" && strings.HasPrefix(line, "  }")) {
			continue
		}
		trimmed := strings.TrimPrefix(line, "        ")
		trimmed = strings.ReplaceAll(trimmed, `\t`, "\t")
		trimmed = strings.ReplaceAll(trimmed, `\"`, `"`)
		icLines = append(icLines, trimmed)
	}
	newIterCursorGo := strings.Join(icLines, "\n") + "\n"

	oldStart := strings.Index(rtContent, "type JIterCursor struct {")
	oldEnd := strings.Index(rtContent, "func newJavaSet2(e0, e1 any)")
	if oldStart < 0 || oldEnd <= oldStart {
		t.Fatalf("could not locate JIterCursor block in ext_collections_rt.go (s=%d, e=%d)", oldStart, oldEnd)
	}
	rtContent = rtContent[:oldStart] + newIterCursorGo + rtContent[oldEnd:]
	if err := os.WriteFile(rtPath, []byte(rtContent), 0o644); err != nil {
		t.Fatal(err)
	}
	return tmpDir, repoRoot
}

func TestWave17Task1_SlabAllocLowering_UnitAndRealPackages(t *testing.T) {
	tmpDir, repoRoot := setupWave17J2gortOverlayWorkspace(t)

	// 1. Unit + differential + negative soundness test for slab_alloc_lowering + j2gort.Slab[T].
	origDir := filepath.Join(tmpDir, "origslab")
	optDir := filepath.Join(tmpDir, "optslab")
	if err := os.MkdirAll(origDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(optDir, 0o755); err != nil {
		t.Fatal(err)
	}

	makeSlabPkgSrc := func(pkgName string) map[string]string {
		return map[string]string{
			"trees.go": `package ` + pkgName + `

import "example.com/w13maptest/j2gort"

type ParseTree interface {
	Implements_` + pkgName + `_ParseTree()
	IsJavaNil_() bool
	GetParseTreeBase() *ParseTreeBase
}

type ParseTreeBase struct {
	self_ ParseTree
	loc   int32
}

func (b *ParseTreeBase) BindSelf_(s any)                 { b.self_ = s.(ParseTree) }
func (b *ParseTreeBase) IsJavaNil_() bool                { return b == nil }
func (b *ParseTreeBase) GetParseTreeBase() *ParseTreeBase { return b }
func (*ParseTreeBase) Implements_` + pkgName + `_ParseTree() {}

func InitParseTree(b *ParseTreeBase, loc int32) { b.loc = loc }

type IdentTree struct {
	ParseTreeBase
	name string
}

func (t *IdentTree) BindSelf_(s any)                 { t.ParseTreeBase.BindSelf_(s) }
func (t *IdentTree) IsJavaNil_() bool                { return t == nil }
func (*IdentTree) Implements_` + pkgName + `_ParseTree() {}
func (t *IdentTree) initFields_IdentTree()           { _ = t }

func NewIdentTree(loc int32, name string) *IdentTree {
	this_ := &IdentTree{}
	this_.BindSelf_(this_)
	InitIdentTree(this_, loc, name)
	return this_
}

func InitIdentTree(this_ *IdentTree, loc int32, name string) {
	InitParseTree(&this_.ParseTreeBase, loc)
	this_.initFields_IdentTree()
	this_.name = name
}

type ParamListTree struct {
	ParseTreeBase
	count int32
}

func (t *ParamListTree) BindSelf_(s any)                 { t.ParseTreeBase.BindSelf_(s) }
func (t *ParamListTree) IsJavaNil_() bool                { return t == nil }
func (*ParamListTree) Implements_` + pkgName + `_ParseTree() {}
func (t *ParamListTree) initFields_ParamListTree()       { _ = t }

func NewParamListTree(loc int32, count int32) *ParamListTree {
	this_ := &ParamListTree{}
	this_.BindSelf_(this_)
	InitParamListTree(this_, loc, count)
	return this_
}

func InitParamListTree(this_ *ParamListTree, loc int32, count int32) {
	InitParseTree(&this_.ParseTreeBase, loc)
	this_.initFields_ParamListTree()
	this_.count = count
}

type FuncDeclTree struct {
	ParseTreeBase
	name   string
	params *ParamListTree
}

func (t *FuncDeclTree) BindSelf_(s any)                 { t.ParseTreeBase.BindSelf_(s) }
func (t *FuncDeclTree) IsJavaNil_() bool                { return t == nil }
func (*FuncDeclTree) Implements_` + pkgName + `_ParseTree() {}
func (t *FuncDeclTree) initFields_FuncDeclTree()        { _ = t }

func NewFuncDeclTree(b *FuncDeclTree_Builder) *FuncDeclTree {
	this_ := &FuncDeclTree{}
	this_.BindSelf_(this_)
	InitFuncDeclTree(this_, b)
	return this_
}

func InitFuncDeclTree(this_ *FuncDeclTree, b *FuncDeclTree_Builder) {
	InitParseTree(&this_.ParseTreeBase, b.loc)
	this_.initFields_FuncDeclTree()
	this_.name = b.name
	this_.params = j2gort.CheckNotNullPtr(b.params)
}

func FuncDeclTree_S_Builder(name string) *FuncDeclTree_Builder {
	_ = name
	return NewFuncDeclTree_Builder(name)
}

type FuncDeclTree_Builder struct {
	name   string
	params *ParamListTree
	loc    int32
}

func (b *FuncDeclTree_Builder) initFields_FuncDeclTree_Builder() { _ = b }

func NewFuncDeclTree_Builder(name string) *FuncDeclTree_Builder {
	this_ := &FuncDeclTree_Builder{}
	InitFuncDeclTree_Builder(this_, name)
	return this_
}

func InitFuncDeclTree_Builder(this_ *FuncDeclTree_Builder, name string) {
	this_.initFields_FuncDeclTree_Builder()
	this_.name = name
}

func (this_ *FuncDeclTree_Builder) SetParams(p *ParamListTree) *FuncDeclTree_Builder {
	this_.params = p
	return this_
}

func (this_ *FuncDeclTree_Builder) Build(loc int32) *FuncDeclTree {
	_ = loc
	this_.loc = loc
	return NewFuncDeclTree(this_)
}

type Comment struct {
	value string
}

func NewComment(v string) *Comment {
	this_ := &Comment{}
	InitComment(this_, v)
	return this_
}

func InitComment(this_ *Comment, v string) { this_.value = v }
`,
			"parser_and_guards.go": `package ` + pkgName + `

import "example.com/w13maptest/j2gort"

type Parser struct {
	comments []*Comment
}

func NewParser() *Parser { return &Parser{} }

func Parser_NewEmptyParams(loc int32) *ParamListTree {
	_ = loc
	return NewParamListTree(loc, 0)
}

func (this_ *Parser) ParseSample(n int) (*FuncDeclTree, []*IdentTree) {
	this_.comments = append(this_.comments, NewComment("header"))
	var idents []*IdentTree
	for i := 0; i < n; i++ {
		idents = append(idents, NewIdentTree(int32(i), "id"))
	}
	b := FuncDeclTree_S_Builder("fn").SetParams(Parser_NewEmptyParams(99))
	return b.Build(100), idents
}

func (this_ *Parser) Summary(n int) int32 {
	fn, ids := this_.ParseSample(n)
	return fn.loc + int32(len(ids)) + fn.params.loc + int32(len(this_.comments))
}

// Negative guard 1: ReusableParser has a Clear() method so its slabs would leak across resets — MUST NOT lower!
type ReusableParser struct {
	count int32
}

func (this_ *ReusableParser) Clear() { this_.count = 0 }
func (this_ *ReusableParser) ParseOne() ParseTree {
	return NewIdentTree(1, "reused")
}

// Negative guard 2: StatelessPass has no retaining fields for StatelessPass_MeasuredNode — MUST NOT lower!
type StatelessPass struct{}

type StatelessPass_MeasuredNode struct {
	cost int32
}

func NewStatelessPass_MeasuredNode(c int32) *StatelessPass_MeasuredNode {
	this_ := &StatelessPass_MeasuredNode{}
	InitStatelessPass_MeasuredNode(this_, c)
	return this_
}

func InitStatelessPass_MeasuredNode(this_ *StatelessPass_MeasuredNode, c int32) { this_.cost = c }

func (this_ *StatelessPass) Measure(c int32) int32 {
	m := NewStatelessPass_MeasuredNode(c)
	return m.cost
}

var _ = j2gort.NullString
`,
		}
	}

	for name, content := range makeSlabPkgSrc("origslab") {
		if err := os.WriteFile(filepath.Join(origDir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for name, content := range makeSlabPkgSrc("optslab") {
		if err := os.WriteFile(filepath.Join(optDir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	allPasses := passes.AllPasses()
	stats1, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 on optslab failed: %v", err)
	}
	if stats1.ByPass["slab_alloc_lowering"] != 4 {
		t.Fatalf("expected 4 slab_alloc_lowering sites (NewIdentTree, FuncDeclTree_S_Builder, Parser_NewEmptyParams, b.Build), got %d (%+v)", stats1.ByPass["slab_alloc_lowering"], stats1.ByPass)
	}
	stats2, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 on optslab failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites on optslab Run 2 (idempotency), got %+v", stats2)
	}

	optParserBytes, err := os.ReadFile(filepath.Join(optDir, "parser_and_guards.go"))
	if err != nil {
		t.Fatal(err)
	}
	optParserCode := string(optParserBytes)
	for _, want := range []string{
		"slab_IdentTree",
		"slab_FuncDeclTree_Builder",
		"slab_ParamListTree",
		"slab_FuncDeclTree",
		"allocPooled_IdentTree(&this_.slab_IdentTree",
		"allocPooledVia_FuncDeclTree_S_Builder(&this_.slab_FuncDeclTree_Builder",
		"allocPooledVia_Parser_NewEmptyParams(&this_.slab_ParamListTree",
		"b.buildWithPooledSlab_FuncDeclTree(&this_.slab_FuncDeclTree",
	} {
		if !strings.Contains(optParserCode, want) {
			t.Fatalf("expected %q in lowered parser_and_guards.go:\n%s", want, optParserCode)
		}
	}
	for _, forbidden := range []string{
		"slab_Comment",
		"slab_StatelessPass_MeasuredNode",
	} {
		if strings.Contains(optParserCode, forbidden) {
			t.Fatalf("forbidden unsafe slab field %q found in:\n%s", forbidden, optParserCode)
		}
	}
	// Verify ReusableParser (has Clear()) did NOT get any slab_* field!
	reusableIdx := strings.Index(optParserCode, "type ReusableParser struct")
	reusableEnd := strings.Index(optParserCode, "// Negative guard 2:")
	if reusableIdx < 0 || reusableEnd <= reusableIdx || strings.Contains(optParserCode[reusableIdx:reusableEnd], "slab_") {
		t.Fatalf("ReusableParser with Clear() must not receive slab fields:\n%s", optParserCode[reusableIdx:])
	}

	// Verify Slab[T] runtime chunk sizing (0 -> 16 -> 32 -> 64) and differential parity!
	slabRuntimeTestSrc := `package w17_slab_test

import (
	"testing"

	"example.com/w13maptest/j2gort"
	"example.com/w13maptest/optslab"
	"example.com/w13maptest/origslab"
)

type pairItem struct{ a, b int64 }

func TestSlabGrowthAndParity(t *testing.T) {
	// Unused Slab allocates 0 times.
	if allocs := testing.AllocsPerRun(100, func() {
		var s j2gort.Slab[pairItem]
		_ = s
	}); allocs != 0 {
		t.Fatalf("unused Slab allocs = %v, want 0", allocs)
	}
	// 1..16 items allocate exactly 1 backing chunk (n=16).
	if allocs := testing.AllocsPerRun(100, func() {
		var s j2gort.Slab[pairItem]
		for i := 0; i < 16; i++ {
			p := s.Alloc()
			p.a = int64(i)
		}
	}); allocs != 1 {
		t.Fatalf("16-item Slab allocs = %v, want 1", allocs)
	}
	// 17..48 items allocate exactly 2 backing chunks (16 + 32 = 48).
	if allocs := testing.AllocsPerRun(100, func() {
		var s j2gort.Slab[pairItem]
		for i := 0; i < 48; i++ {
			p := s.Alloc()
			p.a = int64(i)
		}
	}); allocs != 2 {
		t.Fatalf("48-item Slab allocs = %v, want 2", allocs)
	}
	// 49..112 items allocate exactly 3 backing chunks (16 + 32 + 64 = 112).
	if allocs := testing.AllocsPerRun(100, func() {
		var s j2gort.Slab[pairItem]
		for i := 0; i < 112; i++ {
			p := s.Alloc()
			p.a = int64(i)
		}
	}); allocs != 3 {
		t.Fatalf("112-item Slab allocs = %v, want 3", allocs)
	}

	// Verify distinct addresses and differential parity across 100 parsed IdentTrees!
	pOrig := origslab.NewParser()
	pOpt := optslab.NewParser()
	if got, want := pOpt.Summary(100), pOrig.Summary(100); got != want {
		t.Fatalf("Parser.Summary(100) = %d, want %d", got, want)
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "w17_slab_test.go"), []byte(slabRuntimeTestSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "test", "-v", "-run", "TestSlabGrowthAndParity", ".")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test TestSlabGrowthAndParity failed: %v\n%s", err, string(out))
	}

	// 2. Real copied `parser` and `graph` packages: verify slab_alloc_lowering sites, Run-2 idempotency, and `go build`!
	copyGenPkg := func(pkgName string) string {
		srcDir := filepath.Join(repoRoot, "transpiled", "gen", pkgName)
		dstDir := filepath.Join(tmpDir, pkgName)
		if err := os.MkdirAll(dstDir, 0o755); err != nil {
			t.Fatal(err)
		}
		ents, err := os.ReadDir(srcDir)
		if err != nil {
			t.Fatal(err)
		}
		for _, ent := range ents {
			if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".go") || strings.HasSuffix(ent.Name(), "_test.go") {
				continue
			}
			raw, err := os.ReadFile(filepath.Join(srcDir, ent.Name()))
			if err != nil {
				t.Fatal(err)
			}
			rewritten := strings.ReplaceAll(string(raw),
				`"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"`,
				`"example.com/w13maptest/j2gort"`)
			if err := os.WriteFile(filepath.Join(dstDir, ent.Name()), []byte(rewritten), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return dstDir
	}

	parserDir := copyGenPkg("parser")
	graphDir := copyGenPkg("graph")
	origParserGo, _ := os.ReadFile(filepath.Join(parserDir, "parser.go"))
	parserPreLowered := strings.Contains(string(origParserGo), "slab_MemberExpressionTree")
	origGraphGo, _ := os.ReadFile(filepath.Join(graphDir, "linked_directed_graph.go"))
	graphPreLowered := strings.Contains(string(origGraphGo), "slab_")

	pStats1, err := passes.RunOnPackageDir(parserDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 on real parser failed: %v", err)
	}
	if !parserPreLowered && pStats1.ByPass["slab_alloc_lowering"] < 140 {
		t.Fatalf("expected >= 140 slab_alloc_lowering sites in real parser, got %d (%+v)", pStats1.ByPass["slab_alloc_lowering"], pStats1.ByPass)
	}
	pStats2, err := passes.RunOnPackageDir(parserDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 on real parser failed: %v", err)
	}
	if pStats2.TotalSites != 0 || pStats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites on real parser Run 2, got %+v", pStats2)
	}
	realParserGo, err := os.ReadFile(filepath.Join(parserDir, "parser.go"))
	if err != nil {
		t.Fatal(err)
	}
	for _, wantSlab := range []string{
		"slab_MemberExpressionTree",
		"slab_IdentifierExpressionTree",
		"slab_CallExpressionTree",
		"slab_BinaryOperatorTree",
		"slab_LiteralExpressionTree",
		"slab_FunctionDeclarationTree",
		"slab_FunctionDeclarationTree_Builder",
		"slab_FormalParameterListTree",
	} {
		if !strings.Contains(string(realParserGo), wantSlab) {
			t.Fatalf("missing %s on real Parser struct", wantSlab)
		}
	}
	if strings.Contains(string(realParserGo), "slab_Comment") {
		t.Fatalf("Parser.comments escapes via ParseResult and must NOT be slab-allocated")
	}

	gStats1, err := passes.RunOnPackageDir(graphDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 on real graph failed: %v", err)
	}
	if !graphPreLowered && gStats1.ByPass["slab_alloc_lowering"] < 6 {
		t.Fatalf("expected >= 6 slab_alloc_lowering sites in real graph (3 LinkedDirectedGraph + 2 LinkedUndirectedGraph + 1 StandardUnionFind), got %d (%+v)", gStats1.ByPass["slab_alloc_lowering"], gStats1.ByPass)
	}
	gStats2, err := passes.RunOnPackageDir(graphDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 on real graph failed: %v", err)
	}
	if gStats2.TotalSites != 0 || gStats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites on real graph Run 2, got %+v", gStats2)
	}

	buildCmd := exec.Command("go", "build", "./parser", "./graph")
	buildCmd.Dir = tmpDir
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("go build ./parser ./graph failed: %v\n%s", err, string(out))
	}
}

func TestWave17Task2_OrderedMapCap8AndLazyValueMap(t *testing.T) {
	tmpDir, _ := setupWave17J2gortOverlayWorkspace(t)

	origDir := filepath.Join(tmpDir, "origpm")
	optDir := filepath.Join(tmpDir, "optpm")
	_ = os.MkdirAll(origDir, 0o755)
	_ = os.MkdirAll(optDir, 0o755)

	makePropertyMapSrc := func(pkgName string) string {
		return `package ` + pkgName + `

import "example.com/w13maptest/j2gort"

type Property struct {
	Name string
	Slot int32
}

type SymKey struct {
	Id int32
}

func (s *SymKey) IsJavaNil_() bool     { return s == nil }
func (s *SymKey) HashCode() int32      { return s.Id * 31 }
func (s *SymKey) Equals(other any) bool {
	o, ok := other.(*SymKey)
	return ok && o != nil && s.Id == o.Id
}

type PropertyMap struct {
	properties   j2gort.Ext_java_util_Map
	knownSymbols j2gort.Ext_java_util_Map
	removableMap j2gort.Ext_java_util_Map
}

func (this_ *PropertyMap) initFields_PropertyMap() {
	this_.properties = nil
	this_.knownSymbols = nil
	this_.removableMap = nil
}

func NewPropertyMap() *PropertyMap {
	this_ := &PropertyMap{}
	InitPropertyMap(this_)
	return this_
}

func InitPropertyMap(this_ *PropertyMap) {
	this_.initFields_PropertyMap()
}

func (this_ *PropertyMap) PutStringProp(name string, p *Property) {
	if this_.properties == nil {
		this_.properties = j2gort.NewExt_java_util_LinkedHashMap()
		j2gort.MapPutStr(this_.properties, name, j2gort.PtrToAny(p))
		return
	}
	j2gort.MapPutStr(this_.properties, name, j2gort.PtrToAny(p))
}

func (this_ *PropertyMap) GetStringProp(name string) *Property {
	if this_.properties != nil {
		return j2gort.CastPtr[Property](j2gort.MapGetStr(this_.properties, name))
	}
	return nil
}

func (this_ *PropertyMap) CountStringProps() int32 {
	if this_.properties == nil {
		return 0
	}
	return this_.properties.Size()
}

func (this_ *PropertyMap) PutSymbolProp(sym *SymKey, p *Property) {
	if this_.knownSymbols == nil {
		this_.knownSymbols = j2gort.NewExt_java_util_LinkedHashMap()
	}
	this_.knownSymbols.Put_2(j2gort.PtrToAny(sym), j2gort.PtrToAny(p))
}

func (this_ *PropertyMap) GetSymbolProp(sym *SymKey) *Property {
	if this_.knownSymbols != nil {
		return j2gort.CastPtr[Property](this_.knownSymbols.Get(j2gort.PtrToAny(sym)))
	}
	return nil
}

// Negative guard: removableMap has Remove, so becoming empty after Remove would conflate empty with nil — MUST stay pointer!
func (this_ *PropertyMap) PutRemovable(name string, p *Property) {
	if this_.removableMap == nil {
		this_.removableMap = j2gort.NewExt_java_util_LinkedHashMap()
		j2gort.MapPutStr(this_.removableMap, name, j2gort.PtrToAny(p))
		return
	}
	j2gort.MapPutStr(this_.removableMap, name, j2gort.PtrToAny(p))
}

func (this_ *PropertyMap) RemoveRemovable(name string) bool {
	if this_.removableMap == nil {
		return false
	}
	_ = j2gort.CastPtr[Property](j2gort.MapRemoveStr(this_.removableMap, name))
	return this_.removableMap != nil
}
`
	}

	if err := os.WriteFile(filepath.Join(origDir, "pm.go"), []byte(makePropertyMapSrc("origpm")), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(optDir, "pm.go"), []byte(makePropertyMapSrc("optpm")), 0o644); err != nil {
		t.Fatal(err)
	}

	allPasses := passes.AllPasses()
	if _, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{}); err != nil {
		t.Fatalf("Run 1 on optpm failed: %v", err)
	}
	stats2, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 on optpm failed: %v", err)
	}
	if stats2.TotalSites != 0 {
		t.Fatalf("expected 0 sites on optpm Run 2, got %+v", stats2)
	}

	optCodeBytes, err := os.ReadFile(filepath.Join(optDir, "pm.go"))
	if err != nil {
		t.Fatal(err)
	}
	optCode := string(optCodeBytes)
	if !strings.Contains(optCode, "properties   j2gort.OrderedMap[string, *Property]") &&
		!strings.Contains(optCode, "properties j2gort.OrderedMap[string, *Property]") {
		t.Fatalf("expected PropertyMap.properties to lower to value j2gort.OrderedMap:\n%s", optCode)
	}
	if !strings.Contains(optCode, "*j2gort.CustomOrderedMap[*SymKey, *Property]") {
		t.Fatalf("expected PropertyMap.knownSymbols to stay pointer *j2gort.CustomOrderedMap:\n%s", optCode)
	}
	if !strings.Contains(optCode, "removableMap *j2gort.OrderedMap[string, *Property]") {
		t.Fatalf("expected removableMap with Remove() to stay pointer *j2gort.OrderedMap:\n%s", optCode)
	}
	if !strings.Contains(optCode, "this_.properties.IsEmpty()") {
		t.Fatalf("expected nil check on lazy value OrderedMap to become .IsEmpty():\n%s", optCode)
	}

	testSrc := `package w17_cap8_test

import (
	"fmt"
	"testing"

	"example.com/w13maptest/j2gort"
	"example.com/w13maptest/optpm"
	"example.com/w13maptest/origpm"
)

func TestCap8SingleAllocAndPropertyMapParity(t *testing.T) {
	keys := [8]string{"v0", "v1", "v2", "v3", "v4", "v5", "v6", "v7"}
	// 1..8 entries in OrderedMap do 2 slice allocations (cap=4 -> cap=8) and 0 hmap allocations!
	if allocs := testing.AllocsPerRun(200, func() {
		var m j2gort.OrderedMap[string, int32]
		for i := 0; i < 8; i++ {
			m.Put(keys[i], int32(i))
		}
		if m.Get("v7") != 7 {
			t.Fatal("lookup failed")
		}
	}); allocs != 2 {
		t.Fatalf("8-entry OrderedMap allocs = %v, want 2 (cap-4 -> cap-8 backing slice, 0 hmap)", allocs)
	}

	// CustomOrderedMap also starts at cap=4 (2 slice allocations for 8 entries, 0 buckets).
	symKeys := [8]*optpm.SymKey{{Id: 1}, {Id: 2}, {Id: 3}, {Id: 4}, {Id: 5}, {Id: 6}, {Id: 7}, {Id: 8}}
	if allocs := testing.AllocsPerRun(200, func() {
		var cm j2gort.CustomOrderedMap[*optpm.SymKey, int32]
		for i := 0; i < 8; i++ {
			cm.Put(symKeys[i], int32(i+10))
		}
		if cm.Get(symKeys[7]) != 17 {
			t.Fatal("custom lookup failed")
		}
	}); allocs != 2 {
		t.Fatalf("8-entry CustomOrderedMap allocs = %v, want 2", allocs)
	}

	// Differential parity on empty and populated PropertyMap!
	a := origpm.NewPropertyMap()
	b := optpm.NewPropertyMap()
	if a.CountStringProps() != b.CountStringProps() || a.GetStringProp("x") != nil || b.GetStringProp("x") != nil {
		t.Fatal("empty PropertyMap mismatch")
	}
	for i := 0; i < 10; i++ {
		k := fmt.Sprintf("p%d", i)
		a.PutStringProp(k, &origpm.Property{Name: k, Slot: int32(i)})
		b.PutStringProp(k, &optpm.Property{Name: k, Slot: int32(i)})
	}
	if a.CountStringProps() != b.CountStringProps() || b.GetStringProp("p9").Slot != 9 {
		t.Fatal("populated PropertyMap mismatch")
	}
	a.PutRemovable("r", &origpm.Property{Name: "r", Slot: 1})
	b.PutRemovable("r", &optpm.Property{Name: "r", Slot: 1})
	if a.RemoveRemovable("r") != b.RemoveRemovable("r") {
		t.Fatal("removableMap non-nil parity mismatch after removing sole entry")
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "w17_cap8_test.go"), []byte(testSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "test", "-v", "-run", "TestCap8SingleAllocAndPropertyMapParity", ".")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test TestCap8SingleAllocAndPropertyMapParity failed: %v\n%s", err, string(out))
	}
}

func TestWave17Task3_PeepholeRemoveDeadCodeArrayDeque_ZeroAlloc(t *testing.T) {
	tmpDir, _ := setupWave17J2gortOverlayWorkspace(t)

	origDir := filepath.Join(tmpDir, "origpeep")
	optDir := filepath.Join(tmpDir, "optpeep")
	_ = os.MkdirAll(origDir, 0o755)
	_ = os.MkdirAll(optDir, 0o755)

	makePeepholeSrc := func(pkgName string) string {
		return `package ` + pkgName + `

import (
	"example.com/w13maptest/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
)

type PeepholeRemoveDeadCode struct {
	deleted int32
}

func (this_ *PeepholeRemoveDeadCode) TrySimplifyUnusedResult(expr nodestore.Node, childCount int, forceBadInvariant bool) nodestore.Node {
	var sideEffectRoots *j2gort.Ext_java_util_ArrayDeque = j2gort.NewExt_java_util_ArrayDeque()
	_ = sideEffectRoots
	var atFixedPoint bool = this_.priv_PeepholeRemoveDeadCode_trySimplifyUnusedResultInternal(expr, childCount, sideEffectRoots)
	_ = atFixedPoint
	if atFixedPoint {
		return expr
	} else if sideEffectRoots.IsEmpty() {
		this_.deleted++
		return nodestore.Node(0)
	} else if nodestore.Node(j2gort.DequePeekFirstNode(sideEffectRoots)) == expr {
		if !(sideEffectRoots.Size() == int32(1)) || forceBadInvariant {
			j2gort.CheckStateMsg(false, sideEffectRoots)
		}
		return expr
	} else {
		var acc uint32
		for !(sideEffectRoots.IsEmpty()) {
			next := nodestore.Node(j2gort.DequePollFirstNode(sideEffectRoots))
			acc += uint32(next)
		}
		return nodestore.Node(acc)
	}
}

func (this_ *PeepholeRemoveDeadCode) priv_PeepholeRemoveDeadCode_trySimplifyUnusedResultInternal(tree nodestore.Node, childCount int, sideEffectRoots *j2gort.Ext_java_util_ArrayDeque) bool {
	if childCount < 0 {
		j2gort.DequeAddLastNode(sideEffectRoots, uint32(tree))
		return true
	}
	if childCount == 0 {
		return false
	}
	if childCount == 1 {
		j2gort.DequeAddLastNode(sideEffectRoots, uint32(tree))
		return false
	}
	for i := 0; i < childCount; i++ {
		this_.priv_PeepholeRemoveDeadCode_trySimplifyUnusedResultInternal(nodestore.Node(uint32(tree)+uint32(i)+1), 1, sideEffectRoots)
	}
	return false
}

func JoinBreadcrumbQueue(parts []string) string {
	var q *j2gort.Ext_java_util_ArrayDeque = j2gort.NewExt_java_util_ArrayDeque()
	for _, p := range parts {
		q.AddLast(j2gort.BoxString(p))
	}
	return j2gort.Ext_java_lang_String_Join_2_any_Ext_java_lang_Iterable(".", j2gort.AsIface[j2gort.Ext_java_lang_Iterable](q))
}
`
	}

	if err := os.WriteFile(filepath.Join(origDir, "peep.go"), []byte(makePeepholeSrc("origpeep")), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(optDir, "peep.go"), []byte(makePeepholeSrc("optpeep")), 0o644); err != nil {
		t.Fatal(err)
	}

	allPasses := passes.AllPasses()
	stats1, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 on optpeep failed: %v", err)
	}
	if stats1.ByPass["native_deque_lowering"] == 0 {
		t.Fatalf("expected native_deque_lowering to fire on PeepholeRemoveDeadCode, stats=%+v", stats1.ByPass)
	}
	stats2, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 on optpeep failed: %v", err)
	}
	if stats2.TotalSites != 0 {
		t.Fatalf("expected 0 sites on optpeep Run 2, got %+v", stats2)
	}

	optCodeBytes, err := os.ReadFile(filepath.Join(optDir, "peep.go"))
	if err != nil {
		t.Fatal(err)
	}
	optCode := string(optCodeBytes)
	if strings.Contains(optCode, "Ext_java_util_ArrayDeque") {
		t.Fatalf("expected zero remaining Ext_java_util_ArrayDeque in optpeep:\n%s", optCode)
	}
	if !strings.Contains(optCode, "var sideEffectRoots j2gort.Deque[nodestore.Node]") {
		t.Fatalf("expected stack value Deque[nodestore.Node] in TrySimplifyUnusedResult:\n%s", optCode)
	}
	if !strings.Contains(optCode, "sideEffectRoots *j2gort.Deque[nodestore.Node]") {
		t.Fatalf("expected pointer parameter *Deque[nodestore.Node] in trySimplifyUnusedResultInternal:\n%s", optCode)
	}
	if !strings.Contains(optCode, "j2gort.CheckStateStr(false, sideEffectRoots.String())") {
		t.Fatalf("expected CheckStateMsg(false, sideEffectRoots) to lower to CheckStateStr(false, sideEffectRoots.String()):\n%s", optCode)
	}

	testSrc := `package w17_peep_test

import (
	"fmt"
	"testing"

	"example.com/w13maptest/optpeep"
	"example.com/w13maptest/origpeep"
	"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
)

func TestPeepholeZeroAllocAndParity(t *testing.T) {
	pOrig := &origpeep.PeepholeRemoveDeadCode{}
	pOpt := &optpeep.PeepholeRemoveDeadCode{}

	for _, count := range []int{-1, 0, 1, 2, 3, 4, 10} {
		got := pOpt.TrySimplifyUnusedResult(nodestore.Node(10), count, false)
		want := pOrig.TrySimplifyUnusedResult(nodestore.Node(10), count, false)
		if got != want {
			t.Fatalf("count=%d: got %v, want %v", count, got, want)
		}
	}

	// Verify 0 heap allocations for 0..4 side-effect roots (both stack non-escape AND [4]Node inline buffer)!
	for _, count := range []int{-1, 0, 1, 3, 4} {
		allocs := testing.AllocsPerRun(500, func() {
			_ = pOpt.TrySimplifyUnusedResult(nodestore.Node(42), count, false)
		})
		if allocs != 0 {
			t.Fatalf("count=%d: expected 0 heap allocs/op on lowered PeepholeRemoveDeadCode ArrayDeque, got %v", count, allocs)
		}
	}

	// Verify CheckStateMsg panic message parity!
	catchMsg := func(fn func()) (msg string) {
		defer func() {
			if r := recover(); r != nil {
				msg = fmt.Sprint(r)
			}
		}()
		fn()
		return ""
	}
	origMsg := catchMsg(func() { pOrig.TrySimplifyUnusedResult(nodestore.Node(7), 1, true) })
	optMsg := catchMsg(func() { pOpt.TrySimplifyUnusedResult(nodestore.Node(7), 1, true) })
	if origMsg == "" || origMsg != optMsg {
		t.Fatalf("CheckStateMsg panic parity mismatch: got %q, want %q", optMsg, origMsg)
	}

	if got, want := optpeep.JoinBreadcrumbQueue([]string{"a", "b", "c"}), origpeep.JoinBreadcrumbQueue([]string{"a", "b", "c"}); got != want {
		t.Fatalf("JoinBreadcrumbQueue = %q, want %q", got, want)
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "w17_peep_test.go"), []byte(testSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "test", "-v", "-run", "TestPeepholeZeroAllocAndParity", ".")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test TestPeepholeZeroAllocAndParity failed: %v\n%s", err, string(out))
	}
}

func TestWave17Task4_HotMapAndSliceUnlocks(t *testing.T) {
	tmpDir, _ := setupWave17J2gortOverlayWorkspace(t)

	origDir := filepath.Join(tmpDir, "orighot")
	optDir := filepath.Join(tmpDir, "opthot")
	_ = os.MkdirAll(origDir, 0o755)
	_ = os.MkdirAll(optDir, 0o755)

	makeHotSitesSrc := func(pkgName string) string {
		return `package ` + pkgName + `

import (
	"example.com/w13maptest/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
)

// 1. TypedScopeCreator.memoized (method-ref IIFE GetNodeToScopeMapper)
type TypedScope struct {
	Root nodestore.Node
	Tag  string
}

type TypedScopeCreator struct {
	memoized *j2gort.Ext_java_util_LinkedHashMap
}

func (this_ *TypedScopeCreator) initFields_TypedScopeCreator() {
	this_.memoized = j2gort.NewExt_java_util_LinkedHashMap()
}

func NewTypedScopeCreator() *TypedScopeCreator {
	this_ := &TypedScopeCreator{}
	this_.initFields_TypedScopeCreator()
	return this_
}

func (this_ *TypedScopeCreator) CreateScope(root nodestore.Node, tag string) *TypedScope {
	s := &TypedScope{Root: root, Tag: tag}
	j2gort.MapPutNode(this_.memoized, root, j2gort.PtrToAny(s))
	return s
}

func (this_ *TypedScopeCreator) LookupScope(root nodestore.Node) *TypedScope {
	return j2gort.CastPtr[TypedScope](j2gort.MapGetNode(this_.memoized, root))
}

func (this_ *TypedScopeCreator) GetNodeToScopeMapper() j2gort.Ext_java_util_function_Function {
	return func() j2gort.Ext_java_util_function_Function {
		bound_5 := this_.memoized
		return j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return bound_5.Get(mrp0_) })
	}()
}

// 2. GlobalNamespace_Name.refsForNode (singleton-or-map union + priv_Name_castRefsForNodeMap + FirstValue)
type GlobalNamespace_Ref struct {
	node  nodestore.Node
	label string
}

type GlobalNamespace_Name struct {
	refsForNode any
}

func NewGlobalNamespace_Name() *GlobalNamespace_Name {
	return &GlobalNamespace_Name{refsForNode: nil}
}

func (this_ *GlobalNamespace_Name) priv_Name_castRefsForNodeMap_3054() j2gort.Ext_java_util_Map {
	return j2gort.AsIface[j2gort.Ext_java_util_Map](this_.refsForNode)
}

func (this_ *GlobalNamespace_Name) PutRef(node nodestore.Node, label string) *GlobalNamespace_Ref {
	ref := &GlobalNamespace_Ref{node: node, label: label}
	if this_.refsForNode == nil {
		this_.refsForNode = j2gort.PtrToAny(ref)
		return ref
	}
	if j2gort.InstanceOf[*GlobalNamespace_Ref](this_.refsForNode) {
		var refsForNodeMap j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
		_ = refsForNodeMap
		var existingRef *GlobalNamespace_Ref = j2gort.CastPtr[GlobalNamespace_Ref](this_.refsForNode)
		_ = existingRef
		j2gort.MapPutNode(refsForNodeMap, existingRef.node, j2gort.PtrToAny(existingRef))
		this_.refsForNode = refsForNodeMap
	}
	j2gort.MapPutNode(this_.priv_Name_castRefsForNodeMap_3054(), node, j2gort.PtrToAny(ref))
	return ref
}

func (this_ *GlobalNamespace_Name) UpdateRefNode(ref *GlobalNamespace_Ref, newNode nodestore.Node) {
	oldNode := ref.node
	ref.node = newNode
	if this_.refsForNode == nil {
		this_.refsForNode = j2gort.PtrToAny(ref)
	} else if j2gort.InstanceOf[*GlobalNamespace_Ref](this_.refsForNode) {
		if !(j2gort.RefEquals(this_.refsForNode, ref)) {
			j2gort.CheckState(false)
		}
	} else {
		var refsForNodeMap j2gort.Ext_java_util_Map = this_.priv_Name_castRefsForNodeMap_3054()
		j2gort.MapRemoveNode(refsForNodeMap, oldNode)
		if newNode != 0 {
			j2gort.MapPutNode(refsForNodeMap, newNode, j2gort.PtrToAny(ref))
		}
	}
}

func (this_ *GlobalNamespace_Name) GetRefForNode(node nodestore.Node) *GlobalNamespace_Ref {
	if this_.refsForNode == nil {
		return nil
	}
	if ref, ok := j2gort.TryCast[*GlobalNamespace_Ref](this_.refsForNode); ok {
		if ref.node == node {
			return ref
		}
		return nil
	}
	return j2gort.CastPtr[GlobalNamespace_Ref](j2gort.MapGetNode(this_.priv_Name_castRefsForNodeMap_3054(), node))
}

func (this_ *GlobalNamespace_Name) GetFirstRef() *GlobalNamespace_Ref {
	j2gort.CheckNotNull(this_.refsForNode)
	if ref, ok := j2gort.TryCast[*GlobalNamespace_Ref](this_.refsForNode); ok {
		return ref
	}
	return j2gort.CastPtr[GlobalNamespace_Ref](j2gort.AsIface[j2gort.Ext_java_util_Iterator](j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.priv_Name_castRefsForNodeMap_3054().Values()).Iterator()).Next())
}

func (this_ *GlobalNamespace_Name) GetRefs() j2gort.Ext_java_util_Collection {
	if this_.refsForNode == nil {
		return j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0())
	}
	if ref, ok := j2gort.TryCast[*GlobalNamespace_Ref](this_.refsForNode); ok {
		return j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_1(j2gort.PtrToAny(ref)))
	}
	return j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.priv_Name_castRefsForNodeMap_3054().Values())
}

// 3. JSDocInfo.GetTypeNodes() (local ArrayList returned as Ext_java_util_Collection)
type JSTypeExpr struct {
	root nodestore.Node
}

type JSDocInfo struct {
	exprs []*JSTypeExpr
}

func NewJSDocInfo(roots ...nodestore.Node) *JSDocInfo {
	es := make([]*JSTypeExpr, len(roots))
	for i, r := range roots {
		es[i] = &JSTypeExpr{root: r}
	}
	return &JSDocInfo{exprs: es}
}

func (this_ *JSDocInfo) GetTypeExpressions() j2gort.Ext_java_util_Collection {
	return j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.ImmListFromTypedSlice(this_.exprs))
}

func (this_ *JSDocInfo) GetTypeNodes() j2gort.Ext_java_util_Collection {
	var exprs j2gort.Ext_java_util_Collection = this_.GetTypeExpressions()
	if exprs.IsEmpty() {
		return j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	}
	var nodes j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList_1_int32(exprs.Size())
	{
		cur := j2gort.NewIterCursor(j2gort.CheckNotNull(exprs))
		for cur.HasNext() {
			raw := cur.Next()
			te := j2gort.CastPtr[JSTypeExpr](raw)
			nodes.Add_1(j2gort.BoxNode(te.root))
		}
	}
	return nodes
}
`
	}

	if err := os.WriteFile(filepath.Join(origDir, "hot.go"), []byte(makeHotSitesSrc("orighot")), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(optDir, "hot.go"), []byte(makeHotSitesSrc("opthot")), 0o644); err != nil {
		t.Fatal(err)
	}

	allPasses := passes.AllPasses()
	if _, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{}); err != nil {
		t.Fatalf("Run 1 on opthot failed: %v", err)
	}
	stats2, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 on opthot failed: %v", err)
	}
	if stats2.TotalSites != 0 {
		t.Fatalf("expected 0 sites on opthot Run 2, got %+v", stats2)
	}

	optBytes, err := os.ReadFile(filepath.Join(optDir, "hot.go"))
	if err != nil {
		t.Fatal(err)
	}
	optCode := string(optBytes)
	for _, want := range []string{
		"memoized j2gort.OrderedMap[nodestore.Node, *TypedScope]",
		"j2gort.OrderedMapGetNodeAny(bound_5, mrp0_)",
		"priv_Name_castRefsForNodeMap_3054() *j2gort.OrderedMap[nodestore.Node, *GlobalNamespace_Ref]",
		".FirstValue()",
		"j2gort.WrapTypedArrayList(nodes)",
	} {
		if !strings.Contains(optCode, want) {
			t.Fatalf("expected %q in lowered opthot/hot.go:\n%s", want, optCode)
		}
	}
	for _, forbidden := range []string{
		"MapPutNode(",
		"MapGetNode(",
		"NewExt_java_util_ArrayList",
	} {
		if strings.Contains(optCode, forbidden) {
			t.Fatalf("expected %q to be completely eliminated from opthot/hot.go:\n%s", forbidden, optCode)
		}
	}

	testSrc := `package w17_hot_test

import (
	"testing"

	"example.com/w13maptest/j2gort"
	"example.com/w13maptest/opthot"
	"example.com/w13maptest/orighot"
	"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
)

func TestHotUnlocksDifferentialAndZeroAllocIter(t *testing.T) {
	// 1. TypedScopeCreator.memoized + GetNodeToScopeMapper
	tscOrig := orighot.NewTypedScopeCreator()
	tscOpt := opthot.NewTypedScopeCreator()
	tscOrig.CreateScope(nodestore.Node(10), "global")
	tscOrig.CreateScope(nodestore.Node(20), "fn")
	tscOpt.CreateScope(nodestore.Node(10), "global")
	tscOpt.CreateScope(nodestore.Node(20), "fn")

	fnOrig := tscOrig.GetNodeToScopeMapper()
	fnOpt := tscOpt.GetNodeToScopeMapper()
	for _, probe := range []any{j2gort.BoxNode(nodestore.Node(10)), j2gort.BoxNode(nodestore.Node(20)), j2gort.BoxNode(nodestore.Node(99)), nil, "not-a-node"} {
		rOrig := fnOrig.Apply(probe)
		rOpt := fnOpt.Apply(probe)
		if (rOrig == nil) != (rOpt == nil) {
			t.Fatalf("probe=%v: nil mismatch orig=%v opt=%v", probe, rOrig, rOpt)
		}
		if rOrig != nil && rOrig.(*orighot.TypedScope).Tag != rOpt.(*opthot.TypedScope).Tag {
			t.Fatalf("probe=%v: tag mismatch", probe)
		}
	}

	// 2. GlobalNamespace_Name singleton -> OrderedMap upgrade + UpdateRefNode + GetFirstRef + GetRefs
	nOrig := orighot.NewGlobalNamespace_Name()
	nOpt := opthot.NewGlobalNamespace_Name()

	r1Orig := nOrig.PutRef(nodestore.Node(1), "a")
	r1Opt := nOpt.PutRef(nodestore.Node(1), "a")
	if nOrig.GetFirstRef().GetLabel_() != nOpt.GetFirstRef().GetLabel_() {
		t.Fatal("singleton GetFirstRef mismatch")
	}
	// Upgrade to multi-ref map
	nOrig.PutRef(nodestore.Node(2), "b")
	nOrig.PutRef(nodestore.Node(3), "c")
	nOpt.PutRef(nodestore.Node(2), "b")
	nOpt.PutRef(nodestore.Node(3), "c")

	if nOrig.GetFirstRef().GetLabel_() != nOpt.GetFirstRef().GetLabel_() {
		t.Fatal("map GetFirstRef mismatch")
	}
	nOrig.UpdateRefNode(r1Orig, nodestore.Node(11))
	nOpt.UpdateRefNode(r1Opt, nodestore.Node(11))

	for _, nd := range []nodestore.Node{1, 2, 3, 11} {
		g1 := nOrig.GetRefForNode(nd)
		g2 := nOpt.GetRefForNode(nd)
		if (g1 == nil) != (g2 == nil) || (g1 != nil && g1.GetLabel_() != g2.GetLabel_()) {
			t.Fatalf("GetRefForNode(%d) mismatch: orig=%v opt=%v", nd, g1, g2)
		}
	}
	if nOrig.GetRefs().Size() != nOpt.GetRefs().Size() {
		t.Fatalf("GetRefs().Size() mismatch")
	}

	// 3. JSDocInfo.GetTypeNodes() + zero-alloc NewIterCursor NextNode() iteration!
	docOrig := orighot.NewJSDocInfo(nodestore.Node(101), nodestore.Node(202), nodestore.Node(303))
	docOpt := opthot.NewJSDocInfo(nodestore.Node(101), nodestore.Node(202), nodestore.Node(303))
	colOrig := docOrig.GetTypeNodes()
	colOpt := docOpt.GetTypeNodes()

	var sumOrig, sumOpt uint32
	c1 := j2gort.NewIterCursor(colOrig)
	for c1.HasNext() {
		sumOrig += c1.NextNode()
	}
	iterAllocs := testing.AllocsPerRun(200, func() {
		sumOpt = 0
		c2 := j2gort.NewIterCursor(colOpt)
		for c2.HasNext() {
			sumOpt += c2.NextNode()
		}
	})
	if sumOrig != sumOpt || sumOpt != 606 {
		t.Fatalf("GetTypeNodes sum mismatch: got %d, want %d", sumOpt, sumOrig)
	}
	if iterAllocs != 0 {
		t.Fatalf("expected 0 allocs iterating WrapTypedArrayList via NewIterCursor.NextNode(), got %v", iterAllocs)
	}
}
`
	// Add GetLabel_() helper on GlobalNamespace_Ref in both packages so external test package can inspect unexported label.
	appendLabelGetter := func(dir, pkg string) {
		extra := "package " + pkg + "\nfunc (r *GlobalNamespace_Ref) GetLabel_() string { return r.label }\n"
		if err := os.WriteFile(filepath.Join(dir, "ref_helper.go"), []byte(extra), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	appendLabelGetter(origDir, "orighot")
	appendLabelGetter(optDir, "opthot")

	if err := os.WriteFile(filepath.Join(tmpDir, "w17_hot_test.go"), []byte(testSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "test", "-v", "-run", "TestHotUnlocksDifferentialAndZeroAllocIter", ".")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test TestHotUnlocksDifferentialAndZeroAllocIter failed: %v\n%s", err, string(out))
	}
}

func TestWave19A_NodeTraversalDevirtAndLeafUnrolling(t *testing.T) {
	tmpDir, _ := setupWave16J2gortOverlayWorkspace(t)

	// Create paired original (unlowered) and lowered packages with internal package name `jscomp`
	// under distinct import paths `orig19/jscomp` and `opt19/jscomp` for exact differential testing.
	origDir := filepath.Join(tmpDir, "orig19", "jscomp")
	optDir := filepath.Join(tmpDir, "opt19", "jscomp")
	if err := os.MkdirAll(origDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(optDir, 0o755); err != nil {
		t.Fatal(err)
	}

	files := map[string]string{
		"rhino_mock.go": `package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
)

const (
	TokSCRIPT            uint8 = 1
	TokFUNCTION          uint8 = 2
	TokMODULE_BODY       uint8 = 3
	TokCLASS             uint8 = 4
	TokCLASS_MEMBERS     uint8 = 5
	TokDEFAULT_VALUE     uint8 = 6
	TokDESTRUCTURING_LHS uint8 = 7

	TokBLOCK        uint8 = 10
	TokFOR          uint8 = 11
	TokFOR_IN       uint8 = 12
	TokFOR_OF       uint8 = 13
	TokFOR_AWAIT_OF uint8 = 14
	TokSWITCH_BODY  uint8 = 15

	TokEXPR_RESULT uint8 = 30
	TokADD         uint8 = 31
	TokNAME        uint8 = 32
	TokNUMBER      uint8 = 33
)

type Node nodestore.Node

var nodeNullPtrEx any = javalang.NewNullPointerException()

func (this_ Node) TokenSwitchOrd_() int32 {
	t := nodestore.Node(this_).Rec().Token
	if t == 0xFF {
		panic(nodeNullPtrEx)
	}
	return int32(t)
}

func (this_ Node) FirstNextTok_() (Node, Node, uint8) {
	f, nx, tok := nodestore.Node(this_).FirstNextTok()
	return Node(f), Node(nx), tok
}

func (this_ Node) HasColdProps_() bool {
	return nodestore.Node(this_).HasColdProps()
}

func (this_ Node) GetFirstChild() Node {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	return Node(nodestore.Node(this_).FirstChild())
}

func (this_ Node) GetNext() Node {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	return Node(nodestore.Node(this_).Next())
}

func (this_ Node) GetParent() Node {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	return Node(nodestore.Node(this_).Parent())
}

func (this_ Node) TokenByte_() uint8 {
	return nodestore.Node(this_).Token()
}

func (this_ Node) nodeGet_changeTime() int32 {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	return nodestore.Node(this_).ChangeTime()
}

func (this_ Node) nodeSet_changeTime(v int32) {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	nodestore.Node(this_).SetChangeTime(v)
}

type PropListItem struct {
	Key int
}

func (this_ Node) nodeSet_propListHead(v *PropListItem) {
	if this_ == 0 {
		panic(nodeNullPtrEx)
	}
	if v == nil {
		if !nodestore.Node(this_).HasColdProps() {
			return
		}
		nodestore.Node(this_).SetPropList(nil)
	} else {
		nodestore.Node(this_).SetPropList(v)
	}
}

func (this_ Node) GetChangeTime() int32             { return this_.nodeGet_changeTime() }
func (this_ Node) SetChangeTime(v int32)            { this_.nodeSet_changeTime(v) }
func (this_ Node) SetPropListHead(v *PropListItem)  { this_.nodeSet_propListHead(v) }

func NewTestNode(tok uint8, children ...Node) Node {
	n := Node(nodestore.AllocNodeToken(0, tok))
	if len(children) > 0 {
		nodestore.Node(n).SetFirstChild(nodestore.Node(children[0]))
		for i, c := range children {
			nodestore.Node(c).SetParent(nodestore.Node(n))
			if i+1 < len(children) {
				nodestore.Node(c).SetNext(nodestore.Node(children[i+1]))
			} else {
				nodestore.Node(c).SetNext(0)
			}
			if i == 0 {
				nodestore.Node(c).SetPrevious(nodestore.Node(children[len(children)-1]))
			} else {
				nodestore.Node(c).SetPrevious(nodestore.Node(children[i-1]))
			}
		}
	}
	return n
}

func AttachChild(parent, child Node) {
	first := nodestore.Node(parent).FirstChild()
	if first == 0 {
		nodestore.Node(parent).SetFirstChild(nodestore.Node(child))
		nodestore.Node(child).SetParent(nodestore.Node(parent))
		nodestore.Node(child).SetPrevious(nodestore.Node(child))
		nodestore.Node(child).SetNext(0)
		return
	}
	last := first.Previous()
	last.SetNext(nodestore.Node(child))
	nodestore.Node(child).SetPrevious(last)
	first.SetPrevious(nodestore.Node(child))
	nodestore.Node(child).SetParent(nodestore.Node(parent))
	nodestore.Node(child).SetNext(0)
}

func DetachFirstChild(parent Node) {
	first := nodestore.Node(parent).FirstChild()
	if first == 0 {
		return
	}
	nx := first.Next()
	nodestore.Node(parent).SetFirstChild(nx)
	first.SetParent(0)
	first.SetNext(0)
	first.SetPrevious(0)
}
`,
		"node_traversal.go": `package jscomp

type AbstractCompiler interface {
	IsCompiler_()
}

type MockCompiler struct{}

func (*MockCompiler) IsCompiler_() {}

type AbstractScope interface {
	RootNode() Node
}

type SimpleScope struct {
	root Node
}

func (s *SimpleScope) RootNode() Node { return s.root }

type ScopeCreator interface {
	IsJavaNil_() bool
	CreateScope(n Node, parent AbstractScope) AbstractScope
}

type RedeclarationHandler struct{}

var defaultRedeclHandler = &RedeclarationHandler{}

func StaticGet_SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER() *RedeclarationHandler {
	return defaultRedeclHandler
}

type ScopeScanner struct {
	compiler                      AbstractCompiler
	redeclarationHandler          *RedeclarationHandler
	treatProvidesAsRedeclarations bool
}

func (s *ScopeScanner) Reset() {}

type SyntacticScopeCreator struct {
	compiler                      AbstractCompiler
	redeclarationHandler          *RedeclarationHandler
	treatProvidesAsRedeclarations bool
	scannerInUse                  bool
	reusableScanner               *ScopeScanner
}

func (s *SyntacticScopeCreator) IsJavaNil_() bool { return s == nil }

func (s *SyntacticScopeCreator) CreateScope(n Node, parent AbstractScope) AbstractScope {
	_ = parent
	return &SimpleScope{root: n}
}

func InitSyntacticScopeCreator_1(this_ *SyntacticScopeCreator, compiler AbstractCompiler) {
	this_.compiler = compiler
	this_.redeclarationHandler = StaticGet_SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER()
	this_.reusableScanner = &ScopeScanner{
		compiler:             compiler,
		redeclarationHandler: this_.redeclarationHandler,
	}
}

func NewSyntacticScopeCreator_1(compiler AbstractCompiler) *SyntacticScopeCreator {
	s := &SyntacticScopeCreator{}
	InitSyntacticScopeCreator_1(s, compiler)
	return s
}

type NodeTraversal_Callback interface {
	ShouldTraverse_3(t *NodeTraversal, n Node, parent Node) bool
	Visit_3(t *NodeTraversal, n Node, parent Node)
}

type NodeTraversal_ScopedCallback interface {
	NodeTraversal_Callback
	EnterScope(t *NodeTraversal)
	ExitScope(t *NodeTraversal)
}

type NodeTraversal struct {
	compiler                                       AbstractCompiler
	callback                                       NodeTraversal_Callback
	currentNode                                    Node
	scopeCreator                                   ScopeCreator
	scopeCallback                                  NodeTraversal_ScopedCallback
	obeyDestructuringAndDefaultValueExecutionOrder bool
	scopeNodes                                     []Node
	scopeObjs                                      []AbstractScope
	scopeDepth                                     int32
	currentHoistScopeRoot                          Node
}

type NodeTraversal_Builder struct {
	compiler                                       AbstractCompiler
	callback                                       NodeTraversal_Callback
	scopeCreator                                   ScopeCreator
	obeyDestructuringAndDefaultValueExecutionOrder bool
}

func NewNodeTraversal_Builder(c AbstractCompiler, cb NodeTraversal_Callback) *NodeTraversal_Builder {
	return &NodeTraversal_Builder{compiler: c, callback: cb}
}

func (this_ *NodeTraversal_Builder) SetObeyDestructuring(v bool) *NodeTraversal_Builder {
	this_.obeyDestructuringAndDefaultValueExecutionOrder = v
	return this_
}

func clinit_NodeTraversal() {}

func InitNodeTraversal(this_ *NodeTraversal, b *NodeTraversal_Builder) {
	this_.compiler = b.compiler
	this_.callback = b.callback
	if sc, ok := b.callback.(NodeTraversal_ScopedCallback); ok {
		this_.scopeCallback = sc
	}
	if b.scopeCreator == nil || b.scopeCreator.IsJavaNil_() {
		this_.scopeCreator = NewSyntacticScopeCreator_1(this_.compiler)
	} else {
		this_.scopeCreator = b.scopeCreator
	}
	this_.obeyDestructuringAndDefaultValueExecutionOrder = b.obeyDestructuringAndDefaultValueExecutionOrder
}

func (this_ *NodeTraversal) GetCurrentNode() Node {
	return this_.currentNode
}

func (this_ *NodeTraversal) priv_NodeTraversal_ensureScopeCapacity_3801() {
	if this_.scopeDepth == int32(len(this_.scopeNodes)) {
		newLen := 8
		if this_.scopeDepth > 0 {
			newLen = int(this_.scopeDepth) << 1
		}
		ns := make([]Node, newLen)
		copy(ns, this_.scopeNodes)
		this_.scopeNodes = ns
		os := make([]AbstractScope, newLen)
		copy(os, this_.scopeObjs)
		this_.scopeObjs = os
	}
}

func (this_ *NodeTraversal) priv_NodeTraversal_pushScope_3802(n Node) {
	if this_.currentNode == 0 {
		panic("pushScope: currentNode is 0")
	}
	this_.priv_NodeTraversal_ensureScopeCapacity_3801()
	this_.scopeNodes[this_.scopeDepth] = n
	this_.scopeDepth++
	if this_.scopeCallback != nil {
		this_.scopeCallback.EnterScope(this_)
	}
}

func (this_ *NodeTraversal) priv_NodeTraversal_popScope_3803() {
	if this_.scopeCallback != nil {
		this_.scopeCallback.ExitScope(this_)
	}
	this_.scopeDepth--
	this_.scopeNodes[this_.scopeDepth] = 0
	this_.scopeObjs[this_.scopeDepth] = nil
}

func (this_ *NodeTraversal) priv_NodeTraversal_createsBlockScope_3804(n Node) bool {
	_ = this_
	switch n.TokenSwitchOrd_() {
	case int32(TokBLOCK):
		p := n.GetParent()
		if p == 0 || p.TokenByte_() == TokSCRIPT || p.TokenByte_() == TokFUNCTION {
			return false
		}
		return true
	case int32(TokFOR), int32(TokFOR_IN), int32(TokFOR_OF), int32(TokFOR_AWAIT_OF), int32(TokSWITCH_BODY), int32(TokCLASS):
		return true
	default:
		return false
	}
}

func (this_ *NodeTraversal) priv_NodeTraversal_handleScript_3810(n Node, parent Node) {
	this_.currentNode = n
	if !this_.callback.ShouldTraverse_3(this_, n, parent) {
		return
	}
	previousHoistScopeRoot := this_.currentHoistScopeRoot
	this_.currentHoistScopeRoot = n
	this_.priv_NodeTraversal_pushScope_3802(n)
	this_.priv_NodeTraversal_traverseChildren_3820(n)
	this_.priv_NodeTraversal_popScope_3803()
	this_.currentHoistScopeRoot = previousHoistScopeRoot
	this_.currentNode = n
	this_.callback.Visit_3(this_, n, parent)
}

func (this_ *NodeTraversal) priv_NodeTraversal_handleFunction_3811(n Node, parent Node) {
	this_.currentNode = n
	if !this_.callback.ShouldTraverse_3(this_, n, parent) {
		return
	}
	this_.priv_NodeTraversal_pushScope_3802(n)
	this_.priv_NodeTraversal_traverseChildren_3820(n)
	this_.priv_NodeTraversal_popScope_3803()
	this_.currentNode = n
	this_.callback.Visit_3(this_, n, parent)
}

func (this_ *NodeTraversal) priv_NodeTraversal_handleModule_3812(n Node, parent Node) {
	this_.currentNode = n
	if !this_.callback.ShouldTraverse_3(this_, n, parent) {
		return
	}
	this_.priv_NodeTraversal_pushScope_3802(n)
	this_.priv_NodeTraversal_traverseChildren_3820(n)
	this_.priv_NodeTraversal_popScope_3803()
	this_.currentNode = n
	this_.callback.Visit_3(this_, n, parent)
}

func (this_ *NodeTraversal) priv_NodeTraversal_handleClass_3813(n Node, parent Node) {
	this_.currentNode = n
	if !this_.callback.ShouldTraverse_3(this_, n, parent) {
		return
	}
	this_.priv_NodeTraversal_pushScope_3802(n)
	this_.priv_NodeTraversal_traverseChildren_3820(n)
	this_.priv_NodeTraversal_popScope_3803()
	this_.currentNode = n
	this_.callback.Visit_3(this_, n, parent)
}

func (this_ *NodeTraversal) priv_NodeTraversal_handleClassMembers_3814(n Node, parent Node) {
	this_.currentNode = n
	if !this_.callback.ShouldTraverse_3(this_, n, parent) {
		return
	}
	this_.priv_NodeTraversal_traverseChildren_3820(n)
	this_.currentNode = n
	this_.callback.Visit_3(this_, n, parent)
}

func (this_ *NodeTraversal) priv_NodeTraversal_handleDestructuringOrDefaultValue_3815(n Node, parent Node) {
	this_.currentNode = n
	if !this_.callback.ShouldTraverse_3(this_, n, parent) {
		return
	}
	for child := n.GetFirstChild(); child != 0; {
		next := child.GetNext()
		this_.priv_NodeTraversal_traverseBranch_3819(child, n)
		child = next
	}
	this_.currentNode = n
	this_.callback.Visit_3(this_, n, parent)
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverseBranch_3819(n Node, parent Node) {
	switch n.TokenSwitchOrd_() {
	case int32(TokSCRIPT):
		this_.priv_NodeTraversal_handleScript_3810(n, parent)
		return
	case int32(TokFUNCTION):
		this_.priv_NodeTraversal_handleFunction_3811(n, parent)
		return
	case int32(TokMODULE_BODY):
		this_.priv_NodeTraversal_handleModule_3812(n, parent)
		return
	case int32(TokCLASS):
		this_.priv_NodeTraversal_handleClass_3813(n, parent)
		return
	case int32(TokCLASS_MEMBERS):
		this_.priv_NodeTraversal_handleClassMembers_3814(n, parent)
		return
	case int32(TokDEFAULT_VALUE), int32(TokDESTRUCTURING_LHS):
		if this_.obeyDestructuringAndDefaultValueExecutionOrder {
			this_.priv_NodeTraversal_handleDestructuringOrDefaultValue_3815(n, parent)
			return
		}
	default:
	}
	this_.currentNode = n
	if !(this_.callback.ShouldTraverse_3(this_, n, parent)) {
		return
	}
	previousHoistScopeRoot := this_.currentHoistScopeRoot
	var createsBlockScope bool = this_.priv_NodeTraversal_createsBlockScope_3804(n)
	_ = createsBlockScope
	if createsBlockScope {
		this_.priv_NodeTraversal_pushScope_3802(n)
	}
	for child := n.GetFirstChild(); child != 0; {
		var next Node = child.GetNext()
		_ = next
		this_.priv_NodeTraversal_traverseBranch_3819(child, n)
		child = next
	}
	if createsBlockScope {
		this_.priv_NodeTraversal_popScope_3803()
		this_.currentHoistScopeRoot = previousHoistScopeRoot
	}
	this_.currentNode = n
	this_.callback.Visit_3(this_, n, parent)
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverseChildren_3820(n Node) {
	for child := n.GetFirstChild(); child != 0; {
		var next Node = child.GetNext()
		_ = next
		this_.priv_NodeTraversal_traverseBranch_3819(child, n)
		child = next
	}
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverse_3821(root Node) {
	this_.currentNode = root
	this_.priv_NodeTraversal_pushScope_3802(root)
	this_.priv_NodeTraversal_traverseBranch_3819(root, 0)
	this_.priv_NodeTraversal_popScope_3803()
}

func NodeTraversal_Traverse(compiler AbstractCompiler, root Node, cb NodeTraversal_Callback) {
	b := NewNodeTraversal_Builder(compiler, cb)
	b.Traverse(root)
}

func (this_ *NodeTraversal_Builder) Traverse(root Node) {
	t := &NodeTraversal{}
	InitNodeTraversal(t, this_)
	t.priv_NodeTraversal_traverse_3821(root)
}
`,
		"callbacks.go": `package jscomp

import "fmt"

// Base classes (end with "Base" — MUST NOT receive AlwaysShouldTraverse_ or NoOpVisit_ directly!)
type NodeTraversal_AbstractPostOrderCallbackBase struct{}

func (this_ *NodeTraversal_AbstractPostOrderCallbackBase) ShouldTraverse_3(t *NodeTraversal, n Node, parent Node) bool {
	_ = this_
	_ = t
	_ = n
	_ = parent
	return true
}

type NodeTraversal_AbstractPreOrderCallbackBase struct{}

func (this_ *NodeTraversal_AbstractPreOrderCallbackBase) Visit_3(t *NodeTraversal, n Node, parent Node) {
	_ = this_
	_ = t
	_ = n
	_ = parent
}

// Positive 1: Concrete post-order callback -> AlwaysShouldTraverse_ ONLY.
type PostOrderRecorder struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	Events []string
}

func (this_ *PostOrderRecorder) Visit_3(t *NodeTraversal, n Node, parent Node) {
	this_.Events = append(this_.Events, fmt.Sprintf("V:%d(p=%d,cur=%d)", n.TokenByte_(), parent, t.GetCurrentNode()))
}

// Positive 2: Concrete pre-order callback -> NoOpVisit_ ONLY.
type PreOrderSkipper struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	SkipTok uint8
	Events  []string
}

func (this_ *PreOrderSkipper) ShouldTraverse_3(t *NodeTraversal, n Node, parent Node) bool {
	this_.Events = append(this_.Events, fmt.Sprintf("S:%d(p=%d,cur=%d)", n.TokenByte_(), parent, t.GetCurrentNode()))
	return n.TokenByte_() != this_.SkipTok
}

// Positive 3: Scoped callback with pure return-true ShouldTraverse_3 AND pure no-op Visit_3 -> BOTH markers!
type ScopedOnlyRecorder struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	NodeTraversal_AbstractPreOrderCallbackBase
	Events []string
}

func (this_ *ScopedOnlyRecorder) EnterScope(t *NodeTraversal) {
	this_.Events = append(this_.Events, fmt.Sprintf("ENTER(cur=%d,depth=%d)", t.GetCurrentNode(), t.scopeDepth))
}

func (this_ *ScopedOnlyRecorder) ExitScope(t *NodeTraversal) {
	this_.Events = append(this_.Events, fmt.Sprintf("EXIT(cur=%d,depth=%d)", t.GetCurrentNode(), t.scopeDepth))
}

// Positive 4: Generic concrete post-order callback -> AlwaysShouldTraverse_[T].
type GenericPostOrderCallback[T any] struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	Tag    T
	Visited int
}

func (this_ *GenericPostOrderCallback[T]) Visit_3(t *NodeTraversal, n Node, parent Node) {
	_ = t
	_ = n
	_ = parent
	this_.Visited++
}

// Negative 1: Struct with ov_ShouldTraverse_3 / ov_Visit_3 hook fields -> MUST NOT get markers!
type HookedCallbackBase struct {
	ov_ShouldTraverse_3 func(*NodeTraversal, Node, Node) bool
	ov_Visit_3          func(*NodeTraversal, Node, Node)
}

func (this_ *HookedCallbackBase) ShouldTraverse_3(t *NodeTraversal, n Node, parent Node) bool {
	_ = this_
	_ = t
	_ = n
	_ = parent
	return true
}

func (this_ *HookedCallbackBase) Visit_3(t *NodeTraversal, n Node, parent Node) {
	_ = this_
	_ = t
	_ = n
	_ = parent
}

type ConcreteHookedCallback struct {
	HookedCallbackBase
}

// Negative 2: Struct referencing self_ -> MUST NOT get markers!
type SelfRefCallback struct {
	self_ NodeTraversal_Callback
}

func (this_ *SelfRefCallback) ShouldTraverse_3(t *NodeTraversal, n Node, parent Node) bool {
	if this_.self_ != nil {
		return this_.self_.ShouldTraverse_3(t, n, parent)
	}
	return true
}

func (this_ *SelfRefCallback) Visit_3(t *NodeTraversal, n Node, parent Node) {
	if this_.self_ != nil {
		this_.self_.Visit_3(t, n, parent)
	}
}

// Negative 3: ReturnFalseCallback -> MUST NOT get AlwaysShouldTraverse_!
type ReturnFalseCallback struct {
	NodeTraversal_AbstractPreOrderCallbackBase
}

func (this_ *ReturnFalseCallback) ShouldTraverse_3(t *NodeTraversal, n Node, parent Node) bool {
	_ = this_
	_ = t
	_ = n
	_ = parent
	return false
}

// Negative 4: ConcreteEmbeddedParent embedded in OverridingChild that overrides ShouldTraverse_3 -> false!
// ConcreteEmbeddedParent MUST NOT get AlwaysShouldTraverse_ because Go would promote it onto *OverridingChild!
type ConcreteEmbeddedParent struct {
	Events []string
}

func (this_ *ConcreteEmbeddedParent) ShouldTraverse_3(t *NodeTraversal, n Node, parent Node) bool {
	_ = this_
	_ = t
	_ = n
	_ = parent
	return true
}

func (this_ *ConcreteEmbeddedParent) Visit_3(t *NodeTraversal, n Node, parent Node) {
	this_.Events = append(this_.Events, fmt.Sprintf("EP:%d", n))
}

type OverridingChild struct {
	ConcreteEmbeddedParent
}

func (this_ *OverridingChild) ShouldTraverse_3(t *NodeTraversal, n Node, parent Node) bool {
	_ = this_
	_ = t
	_ = parent
	return n.TokenByte_() != TokADD
}

// AST-Mutating ShouldTraverse_3 & Sibling-Detaching Visit_3 Callback:
// Attaches grandchildren onto a leaf NAME node inside ShouldTraverse_3, and detaches first child in Visit_3!
type MutatingLeafCallback struct {
	TargetLeaf  Node
	MutatedOnce bool
	DetachTok   uint8
	Events      []string
}

func (this_ *MutatingLeafCallback) ShouldTraverse_3(t *NodeTraversal, n Node, parent Node) bool {
	this_.Events = append(this_.Events, fmt.Sprintf("ST:%d(cur=%d)", n, t.GetCurrentNode()))
	if n == this_.TargetLeaf && !this_.MutatedOnce {
		this_.MutatedOnce = true
		gc1 := NewTestNode(TokNUMBER)
		gc2 := NewTestNode(TokBLOCK, NewTestNode(TokNUMBER))
		AttachChild(n, gc1)
		AttachChild(n, gc2)
	}
	return true
}

func (this_ *MutatingLeafCallback) Visit_3(t *NodeTraversal, n Node, parent Node) {
	this_.Events = append(this_.Events, fmt.Sprintf("VS:%d(p=%d,cur=%d)", n, parent, t.GetCurrentNode()))
	if this_.DetachTok != 0 && n.TokenByte_() == this_.DetachTok && parent != 0 {
		DetachFirstChild(parent)
	}
}
`,
	}

	// order-independent: writing independent Go source files by filename into temp directories
	for name, src := range files {
		if err := os.WriteFile(filepath.Join(origDir, name), []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(optDir, name), []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	allPasses := passes.AllPasses()
	stats1, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 on opt19/jscomp failed: %v", err)
	}
	if stats1.ByPass["node_traversal_lowering"] < 10 {
		t.Fatalf("expected >= 10 node_traversal_lowering sites, got %d (%+v)", stats1.ByPass["node_traversal_lowering"], stats1.ByPass)
	}
	stats2, err := passes.RunOnPackageDir(optDir, allPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 on opt19/jscomp failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("expected 0 sites on opt19/jscomp Run 2 (idempotency), got %+v", stats2)
	}

	ntLoweredBytes, err := os.ReadFile(filepath.Join(optDir, "node_traversal.go"))
	if err != nil {
		t.Fatal(err)
	}
	ntLowered := string(ntLoweredBytes)
	for _, want := range []string{
		"alwaysTraverse_",
		"noOpVisit_",
		"curBlk_.FirstNextTok",
		"traverseBranchBlk_",
		"nodeTraversalBlockScopeTok_",
		"nodeTraversalSpecialOrBlockTok_",
	} {
		if !strings.Contains(ntLowered, want) {
			t.Fatalf("expected %q in lowered node_traversal.go:\n%s", want, ntLowered)
		}
	}

	cbLoweredBytes, err := os.ReadFile(filepath.Join(optDir, "callbacks.go"))
	if err != nil {
		t.Fatal(err)
	}
	cbLowered := string(cbLoweredBytes)
	for _, wantMarker := range []string{
		"func (*PostOrderRecorder) AlwaysShouldTraverse_()",
		"func (*PreOrderSkipper) NoOpVisit_()",
		"func (*ScopedOnlyRecorder) AlwaysShouldTraverse_()",
		"func (*ScopedOnlyRecorder) NoOpVisit_()",
		"func (*GenericPostOrderCallback[T]) AlwaysShouldTraverse_()",
		"func (*ReturnFalseCallback) NoOpVisit_()",
	} {
		if !strings.Contains(cbLowered, wantMarker) {
			t.Fatalf("expected marker %q in lowered callbacks.go:\n%s", wantMarker, cbLowered)
		}
	}
	for _, forbiddenMarker := range []string{
		"NodeTraversal_AbstractPostOrderCallbackBase) AlwaysShouldTraverse_",
		"NodeTraversal_AbstractPreOrderCallbackBase) NoOpVisit_",
		"PostOrderRecorder) NoOpVisit_",
		"PreOrderSkipper) AlwaysShouldTraverse_",
		"ConcreteHookedCallback) AlwaysShouldTraverse_",
		"ConcreteHookedCallback) NoOpVisit_",
		"HookedCallbackBase) AlwaysShouldTraverse_",
		"SelfRefCallback) AlwaysShouldTraverse_",
		"SelfRefCallback) NoOpVisit_",
		"ReturnFalseCallback) AlwaysShouldTraverse_",
		"ConcreteEmbeddedParent) AlwaysShouldTraverse_",
		"OverridingChild) AlwaysShouldTraverse_",
	} {
		if strings.Contains(cbLowered, forbiddenMarker) {
			t.Fatalf("forbidden unsound marker %q emitted in lowered callbacks.go:\n%s", forbiddenMarker, cbLowered)
		}
	}

	diffTestSrc := `package w19a_diff_test

import (
	"fmt"
	"slices"
	"testing"

	optjs "example.com/w13maptest/opt19/jscomp"
	origjs "example.com/w13maptest/orig19/jscomp"
	"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
)

func buildSampleAST() (root uint32, mutTarget uint32) {
	// SCRIPT -> [FUNCTION -> [NAME, BLOCK(Function body -> no block scope) -> [EXPR_RESULT -> ADD -> [NAME(target), NUMBER], BLOCK(nested -> creates block scope!) -> [NAME, NUMBER]]], CLASS -> [NAME, CLASS_MEMBERS -> [NUMBER]], DEFAULT_VALUE -> [NAME, NUMBER]]
	nTarget := optjs.NewTestNode(optjs.TokNAME)
	add := optjs.NewTestNode(optjs.TokADD, nTarget, optjs.NewTestNode(optjs.TokNUMBER))
	expr := optjs.NewTestNode(optjs.TokEXPR_RESULT, add)
	innerBlk := optjs.NewTestNode(optjs.TokBLOCK, optjs.NewTestNode(optjs.TokNAME), optjs.NewTestNode(optjs.TokNUMBER))
	fnBody := optjs.NewTestNode(optjs.TokBLOCK, expr, innerBlk)
	fn := optjs.NewTestNode(optjs.TokFUNCTION, optjs.NewTestNode(optjs.TokNAME), fnBody)
	cls := optjs.NewTestNode(optjs.TokCLASS, optjs.NewTestNode(optjs.TokNAME), optjs.NewTestNode(optjs.TokCLASS_MEMBERS, optjs.NewTestNode(optjs.TokNUMBER)))
	defVal := optjs.NewTestNode(optjs.TokDEFAULT_VALUE, optjs.NewTestNode(optjs.TokNAME), optjs.NewTestNode(optjs.TokNUMBER))
	forLoop := optjs.NewTestNode(optjs.TokFOR, optjs.NewTestNode(optjs.TokNAME), optjs.NewTestNode(optjs.TokBLOCK, optjs.NewTestNode(optjs.TokNUMBER)))
	script := optjs.NewTestNode(optjs.TokSCRIPT, fn, cls, defVal, forLoop)
	return uint32(script), uint32(nTarget)
}

func TestWave19ADifferentialParity(t *testing.T) {
	s := nodestore.NewNodeStore()
	prev := nodestore.SetCurrentNodeStore(s)
	defer nodestore.SetCurrentNodeStore(prev)

	rootID, _ := buildSampleAST()
	origRoot := origjs.Node(rootID)
	optRoot := optjs.Node(rootID)

	// 1. PostOrderRecorder (AlwaysShouldTraverse_ = true, NoOpVisit_ = false)
	for _, obey := range []bool{false, true} {
		c1 := &origjs.PostOrderRecorder{}
		c2 := &optjs.PostOrderRecorder{}
		origjs.NewNodeTraversal_Builder(&origjs.MockCompiler{}, c1).SetObeyDestructuring(obey).Traverse(origRoot)
		optjs.NewNodeTraversal_Builder(&optjs.MockCompiler{}, c2).SetObeyDestructuring(obey).Traverse(optRoot)
		if !slices.Equal(c1.Events, c2.Events) || len(c1.Events) == 0 {
			t.Fatalf("PostOrderRecorder(obey=%v) mismatch:\norig=%v\nopt =%v", obey, c1.Events, c2.Events)
		}
	}

	// 2. PreOrderSkipper (AlwaysShouldTraverse_ = false, NoOpVisit_ = true)
	for _, skipTok := range []uint8{optjs.TokADD, optjs.TokBLOCK, optjs.TokNAME} {
		c1 := &origjs.PreOrderSkipper{SkipTok: skipTok}
		c2 := &optjs.PreOrderSkipper{SkipTok: skipTok}
		origjs.NodeTraversal_Traverse(&origjs.MockCompiler{}, origRoot, c1)
		optjs.NodeTraversal_Traverse(&optjs.MockCompiler{}, optRoot, c2)
		if !slices.Equal(c1.Events, c2.Events) || len(c1.Events) == 0 {
			t.Fatalf("PreOrderSkipper(skip=%d) mismatch:\norig=%v\nopt =%v", skipTok, c1.Events, c2.Events)
		}
	}

	// 3. ScopedOnlyRecorder (AlwaysShouldTraverse_ = true, NoOpVisit_ = true)
	// Verifies pushScope / EnterScope / popScope / ExitScope see 100% identical currentNode!
	{
		c1 := &origjs.ScopedOnlyRecorder{}
		c2 := &optjs.ScopedOnlyRecorder{}
		origjs.NodeTraversal_Traverse(&origjs.MockCompiler{}, origRoot, c1)
		optjs.NodeTraversal_Traverse(&optjs.MockCompiler{}, optRoot, c2)
		if !slices.Equal(c1.Events, c2.Events) || len(c1.Events) == 0 {
			t.Fatalf("ScopedOnlyRecorder mismatch:\norig=%v\nopt =%v", c1.Events, c2.Events)
		}
	}

	// 4. OverridingChild embedding ConcreteEmbeddedParent (negative soundness test!)
	{
		c1 := &origjs.OverridingChild{}
		c2 := &optjs.OverridingChild{}
		origjs.NodeTraversal_Traverse(&origjs.MockCompiler{}, origRoot, c1)
		optjs.NodeTraversal_Traverse(&optjs.MockCompiler{}, optRoot, c2)
		if !slices.Equal(c1.Events, c2.Events) {
			t.Fatalf("OverridingChild mismatch (unsound marker promotion!):\norig=%v\nopt =%v", c1.Events, c2.Events)
		}
	}

	// 5. GenericPostOrderCallback[string]
	{
		g1 := &origjs.GenericPostOrderCallback[string]{Tag: "ok"}
		g2 := &optjs.GenericPostOrderCallback[string]{Tag: "ok"}
		origjs.NodeTraversal_Traverse(&origjs.MockCompiler{}, origRoot, g1)
		optjs.NodeTraversal_Traverse(&optjs.MockCompiler{}, optRoot, g2)
		if g1.Visited == 0 || g1.Visited != g2.Visited {
			t.Fatalf("GenericPostOrderCallback Visited = %d vs %d", g2.Visited, g1.Visited)
		}
	}

	// 6. AST-mutating ShouldTraverse_3 (attaching grandchildren onto a leaf!) + sibling-detaching Visit_3!
	{
		runMutating := func(isOpt bool) []string {
			ms := nodestore.NewNodeStore()
			old := nodestore.SetCurrentNodeStore(ms)
			defer nodestore.SetCurrentNodeStore(old)
			baseNode := nodestore.AllocNodeToken(0, 1)
			rID, tID := buildSampleAST()
			var raw []string
			if !isOpt {
				m1 := &origjs.MutatingLeafCallback{TargetLeaf: origjs.Node(tID), DetachTok: origjs.TokNAME}
				origjs.NodeTraversal_Traverse(&origjs.MockCompiler{}, origjs.Node(rID), m1)
				if !m1.MutatedOnce {
					t.Fatal("orig MutatingLeafCallback did not mutate")
				}
				raw = m1.Events
			} else {
				m2 := &optjs.MutatingLeafCallback{TargetLeaf: optjs.Node(tID), DetachTok: optjs.TokNAME}
				optjs.NodeTraversal_Traverse(&optjs.MockCompiler{}, optjs.Node(rID), m2)
				if !m2.MutatedOnce {
					t.Fatal("opt MutatingLeafCallback did not mutate")
				}
				raw = m2.Events
			}
			norm := make([]string, len(raw))
			base := uint32(baseNode)
			for i, ev := range raw {
				var a, b, c uint32
				if _, err := fmt.Sscanf(ev, "ST:%d(cur=%d)", &a, &b); err == nil {
					norm[i] = fmt.Sprintf("ST:%d(cur=%d)", a-base, b-base)
				} else if _, err := fmt.Sscanf(ev, "VS:%d(p=%d,cur=%d)", &a, &b, &c); err == nil {
					pRel := uint32(0)
					if b != 0 {
						pRel = b - base
					}
					norm[i] = fmt.Sprintf("VS:%d(p=%d,cur=%d)", a-base, pRel, c-base)
				} else {
					norm[i] = ev
				}
			}
			return norm
		}
		evOrig := runMutating(false)
		evOpt := runMutating(true)
		if !slices.Equal(evOrig, evOpt) || len(evOrig) == 0 {
			t.Fatalf("MutatingLeafCallback exact event trace mismatch:\norig=%v\nopt =%v", evOrig, evOpt)
		}
	}

	// 7. Invalid 0xFF token panic parity (both root/branch and leaf child).
	catchPanic := func(fn func()) (panicked bool) {
		defer func() {
			if recover() != nil {
				panicked = true
			}
		}()
		fn()
		return false
	}
	badLeafTree := optjs.NewTestNode(optjs.TokSCRIPT, optjs.NewTestNode(0xFF))
	if !catchPanic(func() {
		origjs.NodeTraversal_Traverse(&origjs.MockCompiler{}, origjs.Node(badLeafTree), &origjs.PostOrderRecorder{})
	}) || !catchPanic(func() {
		optjs.NodeTraversal_Traverse(&optjs.MockCompiler{}, badLeafTree, &optjs.PostOrderRecorder{})
	}) {
		t.Fatalf("expected 0xFF leaf token to panic in both orig and opt")
	}

	// 8. Emitted Node helpers: ChangeTime, HasColdProps_, and nodeSet_propListHead(nil) fast path + NPE parity.
	if !catchPanic(func() { _ = optjs.Node(0).GetChangeTime() }) ||
		!catchPanic(func() { optjs.Node(0).SetChangeTime(1) }) ||
		!catchPanic(func() { optjs.Node(0).SetPropListHead(nil) }) {
		t.Fatalf("expected Node(0) accessors to throw NPE")
	}
	optRoot.SetChangeTime(12345)
	if optRoot.GetChangeTime() != 12345 {
		t.Fatalf("GetChangeTime = %d, want 12345", optRoot.GetChangeTime())
	}
	optRoot.SetChangeTime(-99)
	if optRoot.GetChangeTime() != -99 {
		t.Fatalf("GetChangeTime = %d, want -99", optRoot.GetChangeTime())
	}
	optRoot.SetChangeTime(0)
	if optRoot.GetChangeTime() != 0 || optRoot.HasColdProps_() {
		t.Fatalf("expected zero ChangeTime and false HasColdProps_")
	}
	optRoot.SetPropListHead(nil) // fast short-circuit when !HasColdProps()
	optRoot.SetPropListHead(&optjs.PropListItem{Key: 7})
	if !optRoot.HasColdProps_() {
		t.Fatalf("expected HasColdProps_() == true")
	}
	optRoot.SetPropListHead(nil)
	if optRoot.HasColdProps_() {
		t.Fatalf("expected HasColdProps_() == false after clearing propListHead")
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "w19a_diff_test.go"), []byte(diffTestSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "test", "-v", "-run", "TestWave19ADifferentialParity", ".")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test TestWave19ADifferentialParity failed: %v\n%s", err, string(out))
	}
}
