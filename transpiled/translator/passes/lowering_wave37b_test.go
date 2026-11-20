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

func TestWave37B_SlabsOptimizeCallsReachingUsesFlowSensitiveInlineFunctionsAndRace(t *testing.T) {
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
	goMod := strings.Replace(string(origModBytes), "module github.com/paullepoulpe/jsgo", "module example.com/w37bslabs", 1) +
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

	mustWrite("rhino/node.go", `package rhino

type Node uint32

func (n Node) IsJavaNil_() bool { return n == 0 }
func (n Node) IsFunction() bool { return n >= 100 }
`)

	mustWrite("jscomp/common.go", `package jscomp

import "example.com/w37bslabs/rhino"

type AbstractCompiler interface {
	ReportCodeChange()
}

type ControlFlowGraph[N any] struct {
	entry N
}

func ControlFlowAnalysis_ReleaseCfg(cfg *ControlFlowGraph[rhino.Node]) {
	_ = cfg
}

type DataFlowAnalysisBase[N any, L any] struct {
	cfg *ControlFlowGraph[N]
}

type FlowJoiner[L any] interface {
	JoinFlow(x L)
}

type Scope struct {
	fnBlock bool
}

func (s *Scope) IsFunctionBlockScope() bool { return s != nil && s.fnBlock }

type NodeTraversal struct {
	scopeRoot rhino.Node
	scope     *Scope
	inGlobal  bool
}

func (t *NodeTraversal) InGlobalScope() bool   { return t.inGlobal }
func (t *NodeTraversal) GetScope() *Scope      { return t.scope }
func (t *NodeTraversal) GetScopeRoot() rhino.Node { return t.scopeRoot }

type Var struct {
	name string
}
`)

	// Task 1: OptimizeCalls_NodeList + OptimizeCalls_ReferenceMap + DevirtualizeMethods + OptimizeReturns + OptimizeParameters
	mustWrite("jscomp/opt_calls.go", `package jscomp

import (
	"example.com/w37bslabs/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

type OptimizeCalls_CallGraphCompilerPass interface {
	Process_3(externs rhino.Node, root rhino.Node, references *OptimizeCalls_ReferenceMap)
}

type OptimizeCalls_NodeList struct {
	nodes []rhino.Node
}

var static_OptimizeCalls_NodeList_EMPTY *OptimizeCalls_NodeList

func clinitSlow_OptimizeCalls_NodeList() {
	static_OptimizeCalls_NodeList_EMPTY = NewOptimizeCalls_NodeList()
}

func NewOptimizeCalls_NodeList() *OptimizeCalls_NodeList {
	this_ := &OptimizeCalls_NodeList{}
	InitOptimizeCalls_NodeList(this_)
	return this_
}

func InitOptimizeCalls_NodeList(this_ *OptimizeCalls_NodeList) {
	this_.nodes = nil
}

func (this_ *OptimizeCalls_NodeList) Add(n rhino.Node) {
	this_.nodes = append(this_.nodes, n)
}

type OptimizeCalls_ReferenceMap struct {
	names j2gort.OrderedMap[string, *OptimizeCalls_NodeList]
	props j2gort.OrderedMap[string, *OptimizeCalls_NodeList]
}

func NewOptimizeCalls_ReferenceMap() *OptimizeCalls_ReferenceMap {
	return &OptimizeCalls_ReferenceMap{
		names: j2gort.NewOrderedMap[string, *OptimizeCalls_NodeList](),
		props: j2gort.NewOrderedMap[string, *OptimizeCalls_NodeList](),
	}
}

func (this_ *OptimizeCalls_ReferenceMap) AddNameReference(name string, n rhino.Node) {
	refs := this_.names.Get(name)
	if refs == nil {
		refs = NewOptimizeCalls_NodeList()
		this_.names.Put(name, refs)
	}
	refs.Add(n)
}

func (this_ *OptimizeCalls_ReferenceMap) AddPropReference(prop string, n rhino.Node) {
	refs := this_.props.Get(prop)
	if refs == nil {
		refs = NewOptimizeCalls_NodeList()
		this_.props.Put(prop, refs)
	}
	refs.Add(n)
}

func OptimizeCalls_ReferenceMap_AddValueFunctionNodes(fnNodes *OptimizeCalls_NodeList, n rhino.Node) *OptimizeCalls_NodeList {
	if n.IsFunction() {
		if fnNodes == nil {
			fnNodes = NewOptimizeCalls_NodeList()
		}
		fnNodes.Add(n)
	}
	return fnNodes
}

func OptimizeCalls_ReferenceMap_AppendDefinitionFunctionNodes(fnNodes *OptimizeCalls_NodeList, n rhino.Node) *OptimizeCalls_NodeList {
	return OptimizeCalls_ReferenceMap_AddValueFunctionNodes(fnNodes, n)
}

func OptimizeCalls_ReferenceMap_DefinitionFunctionNodeListFor(n rhino.Node) *OptimizeCalls_NodeList {
	return OptimizeCalls_ReferenceMap_AppendDefinitionFunctionNodes(nil, n)
}

func OptimizeCalls_ReferenceMap_GetFunctionNodeList(refs *OptimizeCalls_NodeList) *OptimizeCalls_NodeList {
	var res *OptimizeCalls_NodeList
	for _, n := range refs.nodes {
		res = OptimizeCalls_ReferenceMap_AddValueFunctionNodes(res, n)
	}
	if res == nil {
		return static_OptimizeCalls_NodeList_EMPTY
	}
	return res
}

func (this_ *OptimizeCalls_ReferenceMap) GetNameFunctions(name string) *OptimizeCalls_NodeList {
	refs := this_.names.Get(name)
	if refs == nil {
		return static_OptimizeCalls_NodeList_EMPTY
	}
	return OptimizeCalls_ReferenceMap_GetFunctionNodeList(refs)
}

type OptimizeCalls struct {
	passes []OptimizeCalls_CallGraphCompilerPass
}

func OptimizeCalls_BuildPropCreation(externs rhino.Node, root rhino.Node) *OptimizeCalls_ReferenceMap {
	_ = externs
	rm := NewOptimizeCalls_ReferenceMap()
	rm.AddNameReference("fn", root)
	rm.AddPropReference("m", root)
	return rm
}

func (this_ *OptimizeCalls) Process_2(externs rhino.Node, root rhino.Node) {
	if len(this_.passes) == 0 {
		return
	}
	var refMap *OptimizeCalls_ReferenceMap = OptimizeCalls_BuildPropCreation(externs, root)
	for _, p := range this_.passes {
		p.Process_3(externs, root, refMap)
	}
}

type DevirtualizeMethods struct {
	refMap    *OptimizeCalls_ReferenceMap
	rewritten int32
}

func (this_ *DevirtualizeMethods) Process_3(externs rhino.Node, root rhino.Node, refMap *OptimizeCalls_ReferenceMap) {
	_ = externs
	_ = root
	this_.refMap = refMap
	cur := refMap.props.IterEntries()
	for cur.NextEntry() {
		this_.priv_DevirtualizeMethods_processReferenceList_2400(cur.Key(), cur.Value())
	}
}

func (this_ *DevirtualizeMethods) priv_DevirtualizeMethods_processReferenceList_2400(name string, callSiteNodes *OptimizeCalls_NodeList) {
	_ = name
	functions := NewOptimizeCalls_NodeList()
	for _, n := range callSiteNodes.nodes {
		if n.IsFunction() {
			functions.Add(n)
		}
	}
	if len(functions.nodes) > 0 {
		this_.rewritten++
	}
}

type OptimizeReturns struct {
	unused int32
}

func (this_ *OptimizeReturns) Process_3(externs rhino.Node, root rhino.Node, definitions *OptimizeCalls_ReferenceMap) {
	_ = externs
	_ = root
	var toOptimize []*OptimizeCalls_NodeList
	cur := definitions.names.IterEntries()
	for cur.NextEntry() {
		fns := OptimizeCalls_ReferenceMap_GetFunctionNodeList(cur.Value())
		if fns != nil && len(fns.nodes) > 0 {
			toOptimize = append(toOptimize, cur.Value())
		}
	}
	this_.unused = int32(len(toOptimize))
}

type OptimizeParameters struct {
	removed int32
}

func OptimizeParameters_CollectParamListNames(fn rhino.Node) *OptimizeCalls_NodeList {
	out := NewOptimizeCalls_NodeList()
	out.Add(fn)
	return out
}

func OptimizeParameters_MayReferenceParamBeforeBody(fn rhino.Node) bool {
	var namesByNames j2gort.OrderedMap[string, *OptimizeCalls_NodeList] = j2gort.NewOrderedMap[string, *OptimizeCalls_NodeList]()
	params := OptimizeParameters_CollectParamListNames(fn)
	namesByNames.Put("p", params)
	return namesByNames.Size() > 0
}

func (this_ *OptimizeParameters) Process_3(externs rhino.Node, root rhino.Node, refMap *OptimizeCalls_ReferenceMap) {
	_ = externs
	fns := OptimizeCalls_ReferenceMap_DefinitionFunctionNodeListFor(root)
	if fns != nil && len(fns.nodes) > 0 && OptimizeParameters_MayReferenceParamBeforeBody(fns.nodes[0]) {
		this_.removed++
	}
	_ = refMap
}
`)

	// Task 2: MaybeReachingVariableUse + FlowSensitiveInlineVariables + InlineFunctions
	mustWrite("jscomp/dfa_inline.go", `package jscomp

import (
	"example.com/w37bslabs/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

type MaybeReachingVariableUse_ReachingUses struct {
	mayUse []int32
}

func NewMaybeReachingVariableUse_ReachingUses_1_int32(numVars int32) *MaybeReachingVariableUse_ReachingUses {
	this_ := &MaybeReachingVariableUse_ReachingUses{}
	InitMaybeReachingVariableUse_ReachingUses_1_int32(this_, numVars)
	return this_
}

func NewMaybeReachingVariableUse_ReachingUses_1(other *MaybeReachingVariableUse_ReachingUses) *MaybeReachingVariableUse_ReachingUses {
	this_ := &MaybeReachingVariableUse_ReachingUses{}
	InitMaybeReachingVariableUse_ReachingUses_1(this_, other)
	return this_
}

func InitMaybeReachingVariableUse_ReachingUses_1_int32(this_ *MaybeReachingVariableUse_ReachingUses, numVars int32) {
	this_.mayUse = make([]int32, numVars)
}

func InitMaybeReachingVariableUse_ReachingUses_1(this_ *MaybeReachingVariableUse_ReachingUses, other *MaybeReachingVariableUse_ReachingUses) {
	this_.mayUse = append([]int32(nil), other.mayUse...)
}

type MaybeReachingVariableUse_ReachingUsesJoinOp struct {
	numVars int32
	result  *MaybeReachingVariableUse_ReachingUses
}

func (this_ *MaybeReachingVariableUse_ReachingUsesJoinOp) initFields_MaybeReachingVariableUse_ReachingUsesJoinOp() {
	_ = this_
}

func NewMaybeReachingVariableUse_ReachingUsesJoinOp(numVars int32) *MaybeReachingVariableUse_ReachingUsesJoinOp {
	this_ := &MaybeReachingVariableUse_ReachingUsesJoinOp{}
	InitMaybeReachingVariableUse_ReachingUsesJoinOp(this_, numVars)
	return this_
}

func InitMaybeReachingVariableUse_ReachingUsesJoinOp(this_ *MaybeReachingVariableUse_ReachingUsesJoinOp, numVars int32) {
	this_.initFields_MaybeReachingVariableUse_ReachingUsesJoinOp()
	this_.numVars = numVars
	this_.result = NewMaybeReachingVariableUse_ReachingUses_1_int32(numVars)
}

func (this_ *MaybeReachingVariableUse_ReachingUsesJoinOp) JoinFlow(x *MaybeReachingVariableUse_ReachingUses) {
	this_.result = x
}

type MaybeReachingVariableUse struct {
	DataFlowAnalysisBase[rhino.Node, *MaybeReachingVariableUse_ReachingUses]
	numVars      int32
	emptyLattice *MaybeReachingVariableUse_ReachingUses
}

func NewMaybeReachingVariableUse(cfg *ControlFlowGraph[rhino.Node], numVars int32) *MaybeReachingVariableUse {
	this_ := &MaybeReachingVariableUse{}
	InitMaybeReachingVariableUse(this_, cfg, numVars)
	return this_
}

func InitMaybeReachingVariableUse(this_ *MaybeReachingVariableUse, cfg *ControlFlowGraph[rhino.Node], numVars int32) {
	this_.cfg = cfg
	this_.numVars = numVars
	this_.emptyLattice = NewMaybeReachingVariableUse_ReachingUses_1_int32(numVars)
}

func (this_ *MaybeReachingVariableUse) CreateInitialEstimateLattice() *MaybeReachingVariableUse_ReachingUses {
	return NewMaybeReachingVariableUse_ReachingUses_1_int32(this_.numVars)
}

func (this_ *MaybeReachingVariableUse) FlowThrough(node rhino.Node, input *MaybeReachingVariableUse_ReachingUses) *MaybeReachingVariableUse_ReachingUses {
	_ = node
	out := NewMaybeReachingVariableUse_ReachingUses_1(input)
	if len(out.mayUse) > 0 {
		out.mayUse[0]++
	}
	return out
}

func (this_ *MaybeReachingVariableUse) CreateFlowJoiner() FlowJoiner[*MaybeReachingVariableUse_ReachingUses] {
	return NewMaybeReachingVariableUse_ReachingUsesJoinOp(this_.numVars)
}

type FlowSensitiveInlineVariables_Candidate struct {
	outer_  *FlowSensitiveInlineVariables
	varName string
	useNode rhino.Node
}

func NewFlowSensitiveInlineVariables_Candidate(outer_ *FlowSensitiveInlineVariables, varName string, useNode rhino.Node) *FlowSensitiveInlineVariables_Candidate {
	this_ := &FlowSensitiveInlineVariables_Candidate{}
	this_.outer_ = outer_
	InitFlowSensitiveInlineVariables_Candidate(this_, varName, useNode)
	return this_
}

func InitFlowSensitiveInlineVariables_Candidate(this_ *FlowSensitiveInlineVariables_Candidate, varName string, useNode rhino.Node) {
	this_.varName = varName
	this_.useNode = useNode
}

type FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback struct {
	outer_ *FlowSensitiveInlineVariables
}

func (this_ *FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback) Visit(n rhino.Node) *FlowSensitiveInlineVariables_Candidate {
	return NewFlowSensitiveInlineVariables_Candidate(this_.outer_, "x", n)
}

type FlowSensitiveInlineVariables struct {
	cfg          *ControlFlowGraph[rhino.Node]
	reachingUses *MaybeReachingVariableUse
	inlined      int32
}

func (this_ *FlowSensitiveInlineVariables) EnterScope(t *NodeTraversal) {
	if t.InGlobalScope() || !t.GetScope().IsFunctionBlockScope() {
		return
	}
	localCfg := &ControlFlowGraph[rhino.Node]{entry: t.GetScopeRoot()}
	defer func() {
		this_.cfg = nil
		this_.reachingUses = nil
		ControlFlowAnalysis_ReleaseCfg(localCfg)
	}()
	this_.cfg = localCfg
	this_.reachingUses = NewMaybeReachingVariableUse(this_.cfg, 4)
	est := this_.reachingUses.CreateInitialEstimateLattice()
	next := this_.reachingUses.FlowThrough(t.GetScopeRoot(), est)
	joiner := this_.reachingUses.CreateFlowJoiner()
	joiner.JoinFlow(next)
	cb := &FlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback{outer_: this_}
	cand := cb.Visit(t.GetScopeRoot())
	if cand.varName == "x" && len(next.mayUse) == 4 && next.mayUse[0] == 1 {
		this_.inlined++
	}
}

func (this_ *FlowSensitiveInlineVariables) Process_2(externs rhino.Node, root rhino.Node) {
	_ = externs
	t := &NodeTraversal{scopeRoot: root, scope: &Scope{fnBlock: true}}
	this_.EnterScope(t)
	this_.EnterScope(t)
}

type InlineFunctions_Function interface {
	GetName() string
}

type InlineFunctions_NamedFunction struct {
	outer_ *InlineFunctions
	fn     rhino.Node
}

func NewInlineFunctions_NamedFunction(outer_ *InlineFunctions, fn rhino.Node) *InlineFunctions_NamedFunction {
	this_ := &InlineFunctions_NamedFunction{}
	this_.outer_ = outer_
	InitInlineFunctions_NamedFunction(this_, fn)
	return this_
}
func InitInlineFunctions_NamedFunction(this_ *InlineFunctions_NamedFunction, fn rhino.Node) {
	this_.fn = fn
}
func (this_ *InlineFunctions_NamedFunction) GetName() string { return "named" }

type InlineFunctions_FunctionVar struct {
	outer_ *InlineFunctions
	fn     rhino.Node
}

func NewInlineFunctions_FunctionVar(outer_ *InlineFunctions, fn rhino.Node) *InlineFunctions_FunctionVar {
	this_ := &InlineFunctions_FunctionVar{}
	this_.outer_ = outer_
	InitInlineFunctions_FunctionVar(this_, fn)
	return this_
}
func InitInlineFunctions_FunctionVar(this_ *InlineFunctions_FunctionVar, fn rhino.Node) {
	this_.fn = fn
}
func (this_ *InlineFunctions_FunctionVar) GetName() string { return "var" }

type InlineFunctions_FunctionExpression struct {
	fn   rhino.Node
	name string
}

func NewInlineFunctions_FunctionExpression(fn rhino.Node, idx int32) *InlineFunctions_FunctionExpression {
	this_ := &InlineFunctions_FunctionExpression{}
	InitInlineFunctions_FunctionExpression(this_, fn, idx)
	return this_
}
func InitInlineFunctions_FunctionExpression(this_ *InlineFunctions_FunctionExpression, fn rhino.Node, idx int32) {
	this_.fn = fn
	_ = idx
	this_.name = "expr"
}
func (this_ *InlineFunctions_FunctionExpression) GetName() string { return this_.name }

type InlineFunctions_FindCandidateFunctions struct {
	outer_ *InlineFunctions
}

func (this_ *InlineFunctions_FindCandidateFunctions) Collect(fn rhino.Node) {
	for _, f := range []InlineFunctions_Function{
		NewInlineFunctions_NamedFunction(this_.outer_, fn),
		NewInlineFunctions_FunctionVar(this_.outer_, fn),
		NewInlineFunctions_FunctionExpression(fn, 1),
	} {
		this_.outer_.fns.Put(f.GetName(), 1)
	}
}

type InlineFunctions struct {
	fns          j2gort.OrderedMap[string, int32]
	collectedCnt int32
}

func (this_ *InlineFunctions) Process_2(externs rhino.Node, root rhino.Node) {
	_ = externs
	this_.fns = j2gort.NewOrderedMap[string, int32]()
	defer func() {
		this_.fns.Clear()
	}()
	finder := &InlineFunctions_FindCandidateFunctions{outer_: this_}
	finder.Collect(root)
	this_.collectedCnt = this_.fns.Size()
}
`)

	slabOnly := []passes.Pass{passes.NewSlabAllocLoweringPass()}
	jscompDir := filepath.Join(tmpDir, "jscomp")

	s1, err := passes.RunOnPackageDir(jscompDir, slabOnly, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if s1.TotalSites == 0 {
		t.Fatalf("expected non-zero Wave 37B slab sites on Run 1")
	}

	s2, err := passes.RunOnPackageDir(jscompDir, slabOnly, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if s2.TotalSites != 0 || s2.FilesModified != 0 {
		t.Fatalf("Run 2 NOT IDEMPOTENT: TotalSites=%d FilesModified=%d byPass=%+v", s2.TotalSites, s2.FilesModified, s2.ByPass)
	}

	optBytes, err := os.ReadFile(filepath.Join(jscompDir, "opt_calls.go"))
	if err != nil {
		t.Fatal(err)
	}
	optGot := string(optBytes)
	for _, want := range []string{
		"slab_OptimizeCalls_NodeList",
		"defer refMap.ReleaseSlabs_()",
		"allocPooledVia_OptimizeCalls_ReferenceMap_AddValueFunctionNodes",
		"allocPooledVia_OptimizeCalls_ReferenceMap_GetFunctionNodeList",
		"allocPooledVia_OptimizeParameters_CollectParamListNames",
		"this_.slab_OptimizeCalls_NodeList.Reset()",
		"defer this_.ReleaseSlabs_()",
	} {
		if !strings.Contains(optGot, want) {
			t.Fatalf("missing %q in lowered opt_calls.go:\n%s", want, optGot)
		}
	}

	dfaBytes, err := os.ReadFile(filepath.Join(jscompDir, "dfa_inline.go"))
	if err != nil {
		t.Fatal(err)
	}
	dfaGot := string(dfaBytes)
	for _, want := range []string{
		"slab_ReachingUses",
		"slab_ReachingUsesJoinOp",
		"ru_.ReleaseSlabs_()",
		"slab_Candidate",
		"this_.slab_Candidate.Reset()",
		"slab_InlineFunctions_FunctionVar",
		"slab_InlineFunctions_NamedFunction",
		"slab_InlineFunctions_FunctionExpression",
	} {
		if !strings.Contains(dfaGot, want) {
			t.Fatalf("missing %q in lowered dfa_inline.go:\n%s", want, dfaGot)
		}
	}

	mustWrite("jscomp/w37b_slab_race_test.go", `package jscomp

import (
	"sync"
	"testing"

	"example.com/w37bslabs/rhino"
)

func TestWave37BSlabRuntimeRace(t *testing.T) {
	clinitSlow_OptimizeCalls_NodeList()
	var wg sync.WaitGroup
	for g := 0; g < 4; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 30; i++ {
				dm := &DevirtualizeMethods{}
				or := &OptimizeReturns{}
				op := &OptimizeParameters{}
				oc := &OptimizeCalls{passes: []OptimizeCalls_CallGraphCompilerPass{dm, or, op}}
				oc.Process_2(1, rhino.Node(101))
				if dm.rewritten != 1 || or.unused != 1 || op.removed != 1 {
					t.Errorf("unexpected OptimizeCalls counts: dm=%d or=%d op=%d", dm.rewritten, or.unused, op.removed)
				}

				fsiv := &FlowSensitiveInlineVariables{}
				fsiv.Process_2(1, rhino.Node(102))
				if fsiv.inlined != 2 {
					t.Errorf("unexpected FlowSensitiveInlineVariables count: %d", fsiv.inlined)
				}

				inlf := &InlineFunctions{}
				inlf.Process_2(1, rhino.Node(103))
				if inlf.collectedCnt != 3 {
					t.Errorf("unexpected InlineFunctions collectedCnt: %d", inlf.collectedCnt)
				}
			}
		}()
	}
	wg.Wait()
}
`)

	cmd := exec.Command("go", "test", "-timeout=60s", "-race", "./jscomp")
	cmd.Dir = tmpDir
	if out, cErr := cmd.CombinedOutput(); cErr != nil {
		t.Fatalf("go test -race on Wave 37B slabs failed: %v\n%s", cErr, string(out))
	}
}

func TestWave37B_OrderedMapUnboxedIntLoweringAndNegativeGuards(t *testing.T) {
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
	goMod := strings.Replace(string(origModBytes), "module github.com/paullepoulpe/jsgo", "module example.com/w37bmapint", 1) +
		fmt.Sprintf("\nrequire github.com/paullepoulpe/jsgo v0.0.0\nreplace github.com/paullepoulpe/jsgo => %s\n", filepath.ToSlash(repoRoot))
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}
	if sumBytes, sErr := os.ReadFile(filepath.Join(repoRoot, "go.sum")); sErr == nil {
		_ = os.WriteFile(filepath.Join(tmpDir, "go.sum"), sumBytes, 0o644)
	}

	pkgDir := filepath.Join(tmpDir, "ser")
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatal(err)
	}

	src := `package ser

import "github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"

type Serializer struct {
	sourceFilePointers j2gort.OrderedMap[string, *j2gort.BoxInt]
	nullableBadMap     j2gort.OrderedMap[string, *j2gort.BoxInt]
	nextId             int32
}

func NewSerializer() *Serializer {
	return &Serializer{
		sourceFilePointers: j2gort.NewOrderedMap[string, *j2gort.BoxInt](),
		nullableBadMap:     j2gort.NewOrderedMap[string, *j2gort.BoxInt](),
	}
}

func (this_ *Serializer) GetPointer(file string) int32 {
	var ptr *j2gort.BoxInt = this_.sourceFilePointers.Get(file)
	_ = ptr
	if ptr == nil {
		id := this_.nextId
		this_.nextId++
		ptr = j2gort.ValueOfInt(id)
		this_.sourceFilePointers.Put(file, ptr)
	}
	return j2gort.UnboxInt(ptr)
}

func (this_ *Serializer) LookupOrNegOne(file string) (int32, any) {
	var p *j2gort.BoxInt = this_.sourceFilePointers.Get(file)
	if p != nil {
		return j2gort.UnboxInt(p), j2gort.PtrToAny(p)
	}
	return -1, nil
}

// Negative guard: explicit nil value stored in map MUST keep *j2gort.BoxInt!
func (this_ *Serializer) PutNull(k string) *j2gort.BoxInt {
	this_.nullableBadMap.Put(k, nil)
	return this_.nullableBadMap.Get(k)
}
`
	if err := os.WriteFile(filepath.Join(pkgDir, "ser.go"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}

	mapOnly := []passes.Pass{passes.NewNativeMapLoweringPass()}
	s1, err := passes.RunOnPackageDir(pkgDir, mapOnly, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if s1.TotalSites == 0 {
		t.Fatalf("expected non-zero OrderedMap unboxed int sites on Run 1")
	}

	s2, err := passes.RunOnPackageDir(pkgDir, mapOnly, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if s2.TotalSites != 0 || s2.FilesModified != 0 {
		t.Fatalf("Run 2 NOT IDEMPOTENT: TotalSites=%d FilesModified=%d", s2.TotalSites, s2.FilesModified)
	}

	gotBytes, err := os.ReadFile(filepath.Join(pkgDir, "ser.go"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(gotBytes)
	if !strings.Contains(got, "sourceFilePointers j2gort.OrderedMap[string, int32]") {
		t.Fatalf("expected sourceFilePointers to lower to OrderedMap[string, int32]:\n%s", got)
	}
	if !strings.Contains(got, ".GetOk(") {
		t.Fatalf("expected GetOk presence check in lowered ser.go:\n%s", got)
	}
	if !strings.Contains(got, "nullableBadMap     j2gort.OrderedMap[string, *j2gort.BoxInt]") &&
		!strings.Contains(got, "nullableBadMap j2gort.OrderedMap[string, *j2gort.BoxInt]") {
		t.Fatalf("expected nullableBadMap to stay *j2gort.BoxInt:\n%s", got)
	}

	testSrc := `package ser

import "testing"

func TestZeroValueVsMissingKey(t *testing.T) {
	s := NewSerializer()
	if v, box := s.LookupOrNegOne("a.js"); v != -1 || box != nil {
		t.Fatalf("expected missing key before insert, got %d %v", v, box)
	}
	// First insert stores pointer 0 — must NOT be confused with missing key on second call!
	if got0 := s.GetPointer("a.js"); got0 != 0 {
		t.Fatalf("expected first pointer 0, got %d", got0)
	}
	if got0Again := s.GetPointer("a.js"); got0Again != 0 {
		t.Fatalf("expected cached pointer 0 (not 1!), got %d", got0Again)
	}
	if got1 := s.GetPointer("b.js"); got1 != 1 {
		t.Fatalf("expected second file pointer 1, got %d", got1)
	}
	if v, box := s.LookupOrNegOne("a.js"); v != 0 || box == nil {
		t.Fatalf("expected key 'a.js' present with value 0, got %d %v", v, box)
	}
}
`
	if err := os.WriteFile(filepath.Join(pkgDir, "ser_test.go"), []byte(testSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "test", "-timeout=60s", "-race", "./ser")
	cmd.Dir = tmpDir
	if out, cErr := cmd.CombinedOutput(); cErr != nil {
		t.Fatalf("go test -race on lowered unboxed int map failed: %v\n%s", cErr, string(out))
	}
}

func TestWave37B_SetSliceAndOrderedMapValueSliceLoweringAndRace(t *testing.T) {
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
	goMod := strings.Replace(string(origModBytes), "module github.com/paullepoulpe/jsgo", "module example.com/w37bcol", 1) +
		fmt.Sprintf("\nrequire github.com/paullepoulpe/jsgo v0.0.0\nreplace github.com/paullepoulpe/jsgo => %s\n", filepath.ToSlash(repoRoot))
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}
	if sumBytes, sErr := os.ReadFile(filepath.Join(repoRoot, "go.sum")); sErr == nil {
		_ = os.WriteFile(filepath.Join(tmpDir, "go.sum"), sumBytes, 0o644)
	}

	pkgDir := filepath.Join(tmpDir, "col")
	if err := os.MkdirAll(pkgDir, 0o755); err != nil {
		t.Fatal(err)
	}

	src := `package col

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

type Removable interface {
	PropName() string
	MarkRemoved()
}

type RemovableItem struct {
	prop    string
	removed bool
}

func (r *RemovableItem) PropName() string { return r.prop }
func (r *RemovableItem) MarkRemoved()     { r.removed = true }

type PropCleaner struct {
	removablesForPropertyNames j2gort.OrderedMap[string, j2gort.Ext_java_util_ArrayList]
	validProperties            j2gort.Ext_java_util_HashSet
	candidates                 *j2gort.Ext_java_util_LinkedHashSet
	calls                      j2gort.Ext_java_util_ArrayList
	visitedSum                 int32
}

func NewPropCleaner() *PropCleaner {
	p := &PropCleaner{}
	p.initFields_PropCleaner()
	return p
}

func (this_ *PropCleaner) initFields_PropCleaner() {
	this_.removablesForPropertyNames = j2gort.NewOrderedMap[string, j2gort.Ext_java_util_ArrayList]()
	this_.validProperties = j2gort.NewExt_java_util_HashSet()
	this_.calls = j2gort.NewExt_java_util_ArrayList()
}

func (this_ *PropCleaner) Consider(r Removable) {
	var propertyName string = r.PropName()
	var list j2gort.Ext_java_util_ArrayList = this_.removablesForPropertyNames.Get(propertyName)
	_ = list
	if list == nil {
		list = j2gort.NewExt_java_util_ArrayList_1_int32(int32(4))
		this_.removablesForPropertyNames.Put(propertyName, list)
	}
	list.Add_1(r)
}

func (this_ *PropCleaner) PinProperty(propertyName string) int32 {
	var removed j2gort.Ext_java_util_ArrayList = this_.removablesForPropertyNames.Remove(propertyName)
	_ = removed
	if removed != nil {
		return removed.Size()
	}
	return 0
}

func (this_ *PropCleaner) RemoveAllRemaining() {
	cur := this_.removablesForPropertyNames.IterEntries()
	for cur.NextEntry() {
		var list j2gort.Ext_java_util_ArrayList = cur.Value()
		_ = list
		var n int32 = list.Size()
		for i := int32(0); i < n; i++ {
			j2gort.AsIface[Removable](list.Get(i)).MarkRemoved()
		}
	}
}

func (this_ *PropCleaner) CheckPropsViaLocalAlias(names []string) bool {
	var validProps j2gort.Ext_java_util_HashSet = this_.validProperties
	_ = validProps
	validProps.Clear()
	for _, n := range names {
		validProps.Add_1(n)
	}
	return validProps.Contains("ok")
}

type Cand struct {
	id   int32
	drop bool
}

func (c *Cand) HasIdentityEquals_() {}

func (this_ *PropCleaner) DrainCandidates(items []*Cand) int32 {
	this_.candidates = nil
	for _, it := range items {
		if this_.candidates == nil {
			this_.candidates = j2gort.NewExt_java_util_LinkedHashSet()
		}
		this_.candidates.Add_1(it)
	}
	kept := int32(0)
	for (this_.candidates != nil) && !(this_.candidates.IsEmpty()) {
		var first *Cand = j2gort.CastPtr[Cand](j2gort.AsIface[j2gort.Ext_java_util_Iterator](this_.candidates.Iterator()).Next())
		this_.candidates.Remove_1_any(j2gort.PtrToAny(first))
		kept += first.id
		{
			var it j2gort.Ext_java_util_Iterator = j2gort.AsIface[j2gort.Ext_java_util_Iterator](this_.candidates.Iterator())
			_ = it
			for it.HasNext() {
				var other *Cand = j2gort.CastPtr[Cand](it.Next())
				_ = other
				if other.drop {
					it.Remove()
				}
			}
		}
	}
	return kept
}

func (this_ *PropCleaner) priv_PropCleaner_onCall(v int32) {
	this_.visitedSum += v
}

func (this_ *PropCleaner) RunCallsForEach(vals []int32) int32 {
	this_.calls.Clear()
	for _, v := range vals {
		this_.calls.Add_1(j2gort.ValueOfInt(v))
	}
	this_.visitedSum = 0
	this_.calls.ForEach(func() j2gort.Ext_java_util_function_Consumer {
		bound_1 := this_
		return j2gort.NewLambda_Ext_java_util_function_Consumer(func(mrp0_ any) {
			bound_1.priv_PropCleaner_onCall(j2gort.UnboxInt(j2gort.CastPtr[j2gort.BoxInt](mrp0_)))
		})
	}())
	return this_.visitedSum
}

func FormatAliasCycle(items []string) string {
	var aliasPath j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet()
	for _, s := range items {
		j2gort.SetAddStr(aliasPath, s)
	}
	return javalang.ToJavaString(aliasPath)
}

func (this_ *PropCleaner) priv_PropCleaner_sumReadOnlySlice(stmts j2gort.Ext_java_util_ArrayList) int32 {
	total := int32(0)
	for i := int32(0); i < stmts.Size(); i++ {
		total += j2gort.UnboxInt(j2gort.CastPtr[j2gort.BoxInt](stmts.Get(i)))
	}
	return total
}

func (this_ *PropCleaner) ComputeViaReadOnlyHelper() int32 {
	stmts := j2gort.NewExt_java_util_ArrayList()
	stmts.Add_1(j2gort.ValueOfInt(10))
	stmts.Add_1(j2gort.ValueOfInt(20))
	return this_.priv_PropCleaner_sumReadOnlySlice(stmts)
}

type Def struct {
	dep string
}

func (d *Def) AddDependencyNamesTo(out j2gort.Ext_java_util_Set) {
	j2gort.SetAddStr(out, d.dep)
}

func CollectDeps(d *Def) bool {
	s := j2gort.NewExt_java_util_LinkedHashSet_1_int32(4)
	d.AddDependencyNamesTo(s)
	return j2gort.SetContainsStr(s, "x")
}

func CollectModifiedRec(items []string, idx int, unsafe *j2gort.Ext_java_util_LinkedHashSet) *j2gort.Ext_java_util_LinkedHashSet {
	if idx >= len(items) {
		return unsafe
	}
	if items[idx] != "" {
		if unsafe == nil {
			unsafe = j2gort.NewExt_java_util_LinkedHashSet()
		}
		j2gort.SetAddStr(unsafe, items[idx])
	}
	unsafe = CollectModifiedRec(items, idx+1, unsafe)
	return unsafe
}

func FindModified(items []string) j2gort.Ext_com_google_common_collect_ImmutableSet {
	var unsafe *j2gort.Ext_java_util_LinkedHashSet = CollectModifiedRec(items, 0, nil)
	if unsafe != nil {
		return j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](unsafe))
	}
	return j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
}
`
	if err := os.WriteFile(filepath.Join(pkgDir, "col.go"), []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}

	colPasses := []passes.Pass{
		passes.NewNativeSetLoweringPass(),
		passes.NewNativeSliceLoweringPass(),
		passes.NewNativeMapLoweringPass(),
	}
	s1, err := passes.RunOnPackageDir(pkgDir, colPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if s1.TotalSites == 0 {
		t.Fatalf("expected non-zero Wave 37B set/slice sites on Run 1")
	}

	s2, err := passes.RunOnPackageDir(pkgDir, colPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if s2.TotalSites != 0 || s2.FilesModified != 0 {
		t.Fatalf("Run 2 NOT IDEMPOTENT: TotalSites=%d FilesModified=%d byPass=%+v", s2.TotalSites, s2.FilesModified, s2.ByPass)
	}

	gotBytes, err := os.ReadFile(filepath.Join(pkgDir, "col.go"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(gotBytes)
	for _, want := range []string{
		"removablesForPropertyNames j2gort.OrderedMap[string, []Removable]",
		"validProperties            j2gort.OrderedSet[string]",
		"var validProps *j2gort.OrderedSet[string] = &this_.validProperties",
		"this_.candidates.GetFirst()",
		"ordRmBuf_",
		"for _, mrp0_ := range this_.calls",
		"(&aliasPath).String()",
		"stmts *[]int32",
		"ordSet_AddDependencyNamesTo",
		"ordSet_CollectModifiedRec",
		"unsafe.ToImmutableSet()",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("missing %q in lowered col.go:\n%s", want, got)
		}
	}

	testSrc := `package col

import "testing"

func TestWave37BColSemantics(t *testing.T) {
	pc := NewPropCleaner()
	var items []*RemovableItem
	// Add 10 items under "foo" (> initial cap 4!) to test slice growth write-back into OrderedMap!
	for i := 0; i < 10; i++ {
		it := &RemovableItem{prop: "foo"}
		items = append(items, it)
		pc.Consider(it)
	}
	barItem := &RemovableItem{prop: "bar"}
	pc.Consider(barItem)

	if n := pc.PinProperty("bar"); n != 1 || barItem.removed {
		t.Fatalf("expected PinProperty('bar') == 1 and unremoved, got %d %v", n, barItem.removed)
	}
	if n := pc.PinProperty("missing"); n != 0 {
		t.Fatalf("expected PinProperty('missing') == 0, got %d", n)
	}
	pc.RemoveAllRemaining()
	for i, it := range items {
		if !it.removed {
			t.Fatalf("item %d was not marked removed (slice header write-back bug!)", i)
		}
	}

	if !pc.CheckPropsViaLocalAlias([]string{"a", "ok", "b"}) {
		t.Fatalf("expected local alias CheckPropsViaLocalAlias to see 'ok'")
	}

	// Test >8 removals in buffered Iterator.Remove() loop!
	var cands []*Cand
	cands = append(cands, &Cand{id: 100, drop: false})
	for i := 0; i < 12; i++ {
		cands = append(cands, &Cand{id: int32(i + 1), drop: true})
	}
	cands = append(cands, &Cand{id: 200, drop: false})
	if sum := pc.DrainCandidates(cands); sum != 300 {
		t.Fatalf("expected DrainCandidates sum == 300, got %d", sum)
	}

	if sum := pc.RunCallsForEach([]int32{3, 4, 5}); sum != 12 {
		t.Fatalf("expected RunCallsForEach == 12, got %d", sum)
	}
	if s := FormatAliasCycle([]string{"a", "b"}); s != "[a, b]" {
		t.Fatalf("expected '[a, b]', got %q", s)
	}
	if sum := pc.ComputeViaReadOnlyHelper(); sum != 30 {
		t.Fatalf("expected ComputeViaReadOnlyHelper == 30, got %d", sum)
	}
	if !CollectDeps(&Def{dep: "x"}) {
		t.Fatalf("expected CollectDeps == true")
	}
	imm := FindModified([]string{"", "p1", "p2", "p1"})
	if imm.Size() != 2 || !imm.Contains("p1") || !imm.Contains("p2") {
		t.Fatalf("unexpected FindModified set: size=%d", imm.Size())
	}
	if emptyImm := FindModified([]string{"", ""}); !emptyImm.IsEmpty() {
		t.Fatalf("expected empty ImmutableSet")
	}
}
`
	if err := os.WriteFile(filepath.Join(pkgDir, "col_test.go"), []byte(testSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "test", "-timeout=60s", "-race", "./col")
	cmd.Dir = tmpDir
	if out, cErr := cmd.CombinedOutput(); cErr != nil {
		t.Fatalf("go test -race on Wave 37B col failed: %v\n%s", cErr, string(out))
	}
}

func TestWave37B_RealGenJscompIdempotencyAndBuild(t *testing.T) {
	t.Parallel()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	repoRoot := filepath.Clean(filepath.Join(wd, "..", "..", ".."))
	if !realGenJ2gortReady(repoRoot) {
		t.Skip("gen/rhino or gen/j2gort is mid-edit; skipping real gen/jscomp copy check")
	}
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
	goMod := strings.Replace(string(origModBytes), "module github.com/paullepoulpe/jsgo", "module example.com/realw37b", 1) +
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

	w37bPasses := []passes.Pass{
		passes.NewNativeMapLoweringPass(),
		passes.NewNativeSetLoweringPass(),
		passes.NewNativeSliceLoweringPass(),
		passes.NewSlabAllocLoweringPass(),
	}
	res1, err := passes.RunOnPackageDir(copyDir, w37bPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 on copied real gen/jscomp failed: %v", err)
	}
	t.Logf("Wave 37B Real gen/jscomp Run 1 sites: total=%d files=%d byPass=%+v", res1.TotalSites, res1.FilesModified, res1.ByPass)

	res2, err := passes.RunOnPackageDir(copyDir, w37bPasses, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 on copied real gen/jscomp failed: %v", err)
	}
	if res2.TotalSites != 0 || res2.FilesModified != 0 {
		t.Fatalf("Wave 37B real gen/jscomp Run 2 NOT IDEMPOTENT: total=%d files=%d byPass=%+v", res2.TotalSites, res2.FilesModified, res2.ByPass)
	}

	checkContains := func(file, sub string) {
		t.Helper()
		b, rErr := os.ReadFile(filepath.Join(copyDir, file))
		if rErr != nil {
			t.Fatal(rErr)
		}
		if !strings.Contains(string(b), sub) {
			t.Fatalf("expected %s to contain %q", file, sub)
		}
	}

	// Task 1: OptimizeCalls_NodeList slabs
	checkContains("optimize_calls.go", "slab_OptimizeCalls_NodeList")
	checkContains("optimize_calls.go", "allocPooledVia_OptimizeCalls_ReferenceMap_AddValueFunctionNodes")
	checkContains("devirtualize_methods.go", "slab_OptimizeCalls_NodeList")
	checkContains("optimize_returns.go", "allocPooledVia_OptimizeCalls_ReferenceMap_GetFunctionNodeList")
	checkContains("optimize_parameters.go", "allocPooledVia_OptimizeParameters_CollectParamListNames")

	// Task 2: MaybeReachingVariableUse + FlowSensitiveInlineVariables + InlineFunctions slabs
	checkContains("maybe_reaching_variable_use.go", "slab_ReachingUses")
	checkContains("flow_sensitive_inline_variables.go", "slab_Candidate")
	checkContains("flow_sensitive_inline_variables.go", "ru_.ReleaseSlabs_()")
	checkContains("inline_functions.go", "slab_InlineFunctions_FunctionVar")

	// Task 3: Unboxed int32 OrderedMaps
	checkContains("typed_ast_serializer.go", "sourceFilePointers.GetOk(")
	checkContains("color_serializer.go", "colorIdToTypePointer.GetOk(")
	if lvaBytes, rErr := os.ReadFile(filepath.Join(copyDir, "live_variables_analysis.go")); rErr == nil && strings.Contains(string(lvaBytes), "scopeVariables") {
		checkContains("live_variables_analysis.go", "scopeVariables.GetOk(")
	}

	// Task 4: Slice & Set lowerings
	checkContains("remove_unused_code.go", "j2gort.OrderedMap[string, []RemoveUnusedCode_Removable]")
	checkContains("inline_object_literals.go", "validProperties j2gort.OrderedSet[string]")
	checkContains("inline_and_collapse_properties.go", "(&aliasPath).String()")
	checkContains("flow_sensitive_inline_variables.go", "this_.candidates.GetFirst()")
	checkContains("closure_rewrite_module.go", "range this_.googModuleGetCalls")
	checkContains("function_argument_injector.go", "ordSet_FunctionArgumentInjector_CollectModifiedParameters")

	buildCmd := exec.Command("go", "build", "./jscomp")
	buildCmd.Dir = tmpDir
	if out, bErr := buildCmd.CombinedOutput(); bErr != nil {
		t.Fatalf("go build on lowered real gen/jscomp failed: %v\n%s", bErr, string(out))
	}
}
