// Port of com/google/javascript/jscomp/ReachingUseDefTester.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/ReachingUseDefTester.java)

/*
 * Copyright 2020 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ReachingUseDefTester is the utility for testing MaybeReachingVariableUse
// and MustBeReachingVariableDef.
type ReachingUseDefTester struct {
	compiler     *Compiler
	scopeCreator *SyntacticScopeCreator
	labelFinder  *ReachingUseDefTester_LabelFinder
	root         *rhino.Node
	reachingUse  *MaybeReachingVariableUse
	reachingDef  *MustBeReachingVariableDef
}

// Java: ReachingUseDefTester(Compiler, SyntacticScopeCreator, LabelFinder)
func newReachingUseDefTester(
	compiler *Compiler, scopeCreator *SyntacticScopeCreator, labelFinder *ReachingUseDefTester_LabelFinder) *ReachingUseDefTester {
	return &ReachingUseDefTester{
		compiler:     compiler,
		scopeCreator: scopeCreator,
		labelFinder:  labelFinder,
	}
}

// Java: ReachingUseDefTester.create()
func ReachingUseDefTester_Create() *ReachingUseDefTester {
	compiler := reachingUseDefTester_createCompiler()
	return newReachingUseDefTester(
		compiler, NewSyntacticScopeCreator(compiler), NewReachingUseDefTester_LabelFinder())
}

// Java: ReachingUseDefTester.createCompiler()
func reachingUseDefTester_createCompiler() *Compiler {
	compiler := NewCompiler()
	compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
	options := NewCompilerOptions()
	options.SetCodingConvention(NewGoogleCodingConvention())
	options.SetLanguage(CompilerOptions_LanguageMode_UNSUPPORTED)
	compiler.InitOptions(options)
	return compiler
}

// ComputeReachingUses runs `MaybeReachingVariableUse` pass to compute and
// store the may-be-reaching uses for all definitions of each variable in the
// test source.
//
// Java: ReachingUseDefTester.computeReachingUses(String, boolean)
func (r *ReachingUseDefTester) ComputeReachingUses(src string, async bool) {
	script := r.parseScript(src, async)
	r.root = script.GetFirstChild()
	funcBlockScope := r.computeFunctionBlockScope(script, r.root)
	cfg := r.computeCfg(r.root)
	escaped := jutil.NewObjectLinkedHashSet[*Var]()

	allVarsDeclaredInFunction :=
		NodeUtil_GetAllVarsDeclaredInFunction(r.compiler, r.scopeCreator, funcBlockScope.GetParent())
	allVarsInFn := allVarsDeclaredInFunction.GetAllVariables()
	DataFlowAnalysis_ComputeEscaped(
		funcBlockScope.GetParent(), escaped, r.compiler, r.scopeCreator, allVarsInFn)
	r.reachingUse = NewMaybeReachingVariableUse(cfg, escaped, allVarsInFn)
	r.reachingUse.Analyze()
}

// ComputeReachingDef runs `MustBeReachingVariableDef` pass to compute and
// store the must-be-reaching definition for all uses of each variable in the
// test source.
//
// Java: ReachingUseDefTester.computeReachingDef(String)
func (r *ReachingUseDefTester) ComputeReachingDef(src string) {
	script := r.parseScript(src /* async= */, false)
	r.root = script.GetFirstChild()
	funcBlockScope := r.computeFunctionBlockScope(script, r.root)
	cfg := r.computeCfg(r.root)
	escaped := jutil.NewObjectLinkedHashSet[*Var]()

	allVarsDeclaredInFunction :=
		NodeUtil_GetAllVarsDeclaredInFunction(r.compiler, r.scopeCreator, funcBlockScope.GetParent())
	allVarsInFn := allVarsDeclaredInFunction.GetAllVariables()
	DataFlowAnalysis_ComputeEscaped(
		funcBlockScope.GetParent(), escaped, r.compiler, r.scopeCreator, allVarsInFn)
	r.reachingDef = NewMustBeReachingVariableDef(cfg, r.compiler, escaped, allVarsInFn)
	r.reachingDef.Analyze()
}

// Java: ReachingUseDefTester.parseScript(String, boolean)
func (r *ReachingUseDefTester) parseScript(src string, async bool) *rhino.Node {
	prefix := ""
	if async {
		prefix = "async "
	}
	src = prefix + "function _FUNCTION(param1, param2){" + src + "}"
	script := r.compiler.ParseTestCode(src)
	AssertThat(r.compiler.GetErrors()).IsEmpty()
	return script
}

// Java: ReachingUseDefTester.computeFunctionBlockScope(Node, Node)
func (r *ReachingUseDefTester) computeFunctionBlockScope(script *rhino.Node, fn *rhino.Node) *Scope {
	functionBlock := fn.GetLastChild()
	globalScope := r.scopeCreator.CreateScopeOfScope(script, nil)
	functionScope := r.scopeCreator.CreateScopeOfScope(r.root, globalScope)
	return r.scopeCreator.CreateScopeOfScope(functionBlock, functionScope)
}

// Java: ReachingUseDefTester.computeCfg(Node)
func (r *ReachingUseDefTester) computeCfg(fn *rhino.Node) ControlFlowGraph[*rhino.Node] {
	return ControlFlowAnalysis_Builder_().
		SetCompiler(r.compiler).
		SetCfgRoot(fn).
		SetIncludeEdgeAnnotations(true).
		ComputeCfg()
}

// GetComputedUses returns may-be-reaching uses of definition of variable `x`
// on the node extracted at label `D:`.
//
// Java: ReachingUseDefTester.getComputedUses()
func (r *ReachingUseDefTester) GetComputedUses() *jutil.LinkedHashSet[*rhino.Node] {
	uses := jutil.NewLinkedHashSet[*rhino.Node]()
	for n := range r.reachingUse.GetUses("x", r.labelFinder.extractedDef) {
		uses.Add(n)
	}
	return uses
}

// GetComputedDef returns must-be-reaching definition of variable `x` on the
// use node extracted at label `U:`.
//
// Java: ReachingUseDefTester.getComputedDef()
func (r *ReachingUseDefTester) GetComputedDef() *rhino.Node {
	return r.reachingDef.GetDefNode("x", r.labelFinder.extractedUses[0])
}

// ExtractDefAndUsesFromInputLabels runs `LabelFinder` to find the `D:` and
// `U:` labels and save the def and uses of `x`.
//
// Java: ReachingUseDefTester.extractDefAndUsesFromInputLabels()
func (r *ReachingUseDefTester) ExtractDefAndUsesFromInputLabels() {
	NodeTraversal_Builder_().
		SetCompiler(r.compiler).
		SetCallback(r.labelFinder).
		SetScopeCreator(r.scopeCreator).
		Traverse(r.root)
	AssertWithMessage("Code should have an instruction labeled D").
		That(r.labelFinder.extractedDef).
		IsNotNull()
	AssertWithMessage("Code should have an instruction label starting withing U").
		That(len(r.labelFinder.extractedUses) == 0).
		IsFalse()
}

// ReachingUseDefTester_LabelFinder finds the D: and U: label and store which
// node they point to.
type ReachingUseDefTester_LabelFinder struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	// Def and uses extracted from `D:` and `U:` labels respectively
	extractedDef  *rhino.Node
	extractedUses []*rhino.Node
}

// Java: ReachingUseDefTester.LabelFinder()
func NewReachingUseDefTester_LabelFinder() *ReachingUseDefTester_LabelFinder {
	return &ReachingUseDefTester_LabelFinder{extractedUses: []*rhino.Node{}}
}

// Java: ReachingUseDefTester.LabelFinder.visit(NodeTraversal, Node, Node)
func (f *ReachingUseDefTester_LabelFinder) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsLabel() {
		if n.GetFirstChild().GetString() == "D" {
			AssertWithMessage("Multiple D: labels in test src").That(f.extractedDef).IsNull()
			f.extractedDef = n.GetLastChild()
		} else if jutil.StringStartsWith(n.GetFirstChild().GetString(), "U") {
			f.extractedUses = append(f.extractedUses, n.GetLastChild())
		}
	}
}

// Java: ReachingUseDefTester.getExtractedUses()
func (r *ReachingUseDefTester) GetExtractedUses() []*rhino.Node {
	return r.labelFinder.extractedUses
}

// Java: ReachingUseDefTester.getExtractedDef()
func (r *ReachingUseDefTester) GetExtractedDef() *rhino.Node {
	return r.labelFinder.extractedDef
}
