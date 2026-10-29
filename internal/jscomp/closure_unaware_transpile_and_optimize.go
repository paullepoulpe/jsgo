// Port of com/google/javascript/jscomp/TranspileAndOptimizeClosureUnaware.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/ManageClosureUnawareCode.java @ closure-compiler v20261021

/*
 * Copyright 2024 The Closure Compiler Authors.
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

/*
 * Copyright 2025 The Closure Compiler Authors.
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
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TranspileAndOptimizeClosureUnaware runs a new Compiler instance over all
// @closureUnaware code (i.e. "shadow roots") in the AST.
//
// This is a separate Compiler instance to avoid leaking compiler state
// between the shadow AST and main AST. Otherwise, we risk doing unsafe
// optimizations on @closureUnaware code.
type TranspileAndOptimizeClosureUnaware struct {
	original AbstractCompiler
	mode     NestedCompilerRunner_Mode
}

// Java: TranspileAndOptimizeClosureUnaware(AbstractCompiler, NestedCompilerRunner.Mode)
func NewTranspileAndOptimizeClosureUnaware(original AbstractCompiler, mode NestedCompilerRunner_Mode) *TranspileAndOptimizeClosureUnaware {
	return &TranspileAndOptimizeClosureUnaware{original: original, mode: mode}
}

// Java: TranspileAndOptimizeClosureUnaware.process(Node, Node)
func (p *TranspileAndOptimizeClosureUnaware) Process(externs rhino.Node, root rhino.Node) {
	collector := newTranspileAndOptimizeClosureUnaware_CollectShadowAsts()
	NodeTraversal_Traverse(p.original, root, collector)
	if collector.isEmpty() {
		return
	}

	allValues := collector.values()
	p.setScriptFeaturesets(allValues)
	shadowOptions := ClosureUnawareOptions_Convert(p.original.GetOptions())
	shadowPassConfig :=
		newTranspileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig(
			shadowOptions, p.createInjectRuntimeLibrariesPass(allValues))
	shadowCompiler :=
		NestedCompilerRunner_Create(p.original, shadowOptions, p.mode, shadowPassConfig)

	p.initShadowInputs(shadowCompiler, collector.groups)
	shadowCompiler.Compile()
	p.reattachShadowNodes(allValues)
}

// transpileAndOptimizeClosureUnaware_ShadowAst is the record
// ShadowAst(Node originalRoot, Node callNode, Node script).
type transpileAndOptimizeClosureUnaware_ShadowAst struct {
	originalRoot rhino.Node
	callNode     rhino.Node
	script       rhino.Node
}

// Java: TranspileAndOptimizeClosureUnaware.ShadowAst(Node, Node, Node)
func newTranspileAndOptimizeClosureUnaware_ShadowAst(originalRoot rhino.Node, callNode rhino.Node, script rhino.Node) *transpileAndOptimizeClosureUnaware_ShadowAst {
	return &transpileAndOptimizeClosureUnaware_ShadowAst{originalRoot: originalRoot, callNode: callNode, script: script}
}

// Java: TranspileAndOptimizeClosureUnaware.ShadowAst.originalRoot()
func (a *transpileAndOptimizeClosureUnaware_ShadowAst) OriginalRoot() rhino.Node {
	return a.originalRoot
}

// Java: TranspileAndOptimizeClosureUnaware.ShadowAst.callNode()
func (a *transpileAndOptimizeClosureUnaware_ShadowAst) CallNode() rhino.Node { return a.callNode }

// Java: TranspileAndOptimizeClosureUnaware.ShadowAst.script()
func (a *transpileAndOptimizeClosureUnaware_ShadowAst) Script() rhino.Node { return a.script }

// Equals is the record's equals: the components are equal.
//
// Java: TranspileAndOptimizeClosureUnaware.ShadowAst.equals(Object)
func (a *transpileAndOptimizeClosureUnaware_ShadowAst) Equals(o any) bool {
	other, ok := o.(*transpileAndOptimizeClosureUnaware_ShadowAst)
	if !ok || other == nil {
		return false
	}
	return jutil.ObjectsEquals(a.originalRoot, other.originalRoot) &&
		jutil.ObjectsEquals(a.callNode, other.callNode) &&
		jutil.ObjectsEquals(a.script, other.script)
}

// HashCode is the record's hashCode: 31-combined component hash codes.
//
// Java: TranspileAndOptimizeClosureUnaware.ShadowAst.hashCode()
func (a *transpileAndOptimizeClosureUnaware_ShadowAst) HashCode() int32 {
	h := jutil.ObjectHashCode(a.originalRoot)
	h = 31*h + jutil.ObjectHashCode(a.callNode)
	h = 31*h + jutil.ObjectHashCode(a.script)
	return h
}

// String is the record's toString.
//
// Java: TranspileAndOptimizeClosureUnaware.ShadowAst.toString()
func (a *transpileAndOptimizeClosureUnaware_ShadowAst) String() string {
	return "ShadowAst[originalRoot=" + jutil.ToJavaString(a.originalRoot) +
		", callNode=" + jutil.ToJavaString(a.callNode) +
		", script=" + jutil.ToJavaString(a.script) + "]"
}

type shadowFileGroup struct {
	file *SourceFile
	asts []*transpileAndOptimizeClosureUnaware_ShadowAst
}

// transpileAndOptimizeClosureUnaware_CollectShadowAsts traverses all
// closure unaware SCRIPTs to collect shadow ASTs.
type transpileAndOptimizeClosureUnaware_CollectShadowAsts struct {
	groups    []shadowFileGroup
	fileIndex map[*SourceFile]int
}

// Java: TranspileAndOptimizeClosureUnaware.CollectShadowAsts()
func newTranspileAndOptimizeClosureUnaware_CollectShadowAsts() *transpileAndOptimizeClosureUnaware_CollectShadowAsts {
	return &transpileAndOptimizeClosureUnaware_CollectShadowAsts{
		fileIndex: make(map[*SourceFile]int),
	}
}

func (c *transpileAndOptimizeClosureUnaware_CollectShadowAsts) add(
	file *SourceFile, ast *transpileAndOptimizeClosureUnaware_ShadowAst) {
	if idx, ok := c.fileIndex[file]; ok {
		for _, existing := range c.groups[idx].asts {
			if jutil.ObjectsEquals(existing, ast) {
				return
			}
		}
		c.groups[idx].asts = append(c.groups[idx].asts, ast)
		return
	}
	c.fileIndex[file] = len(c.groups)
	c.groups = append(c.groups, shadowFileGroup{
		file: file,
		asts: []*transpileAndOptimizeClosureUnaware_ShadowAst{ast},
	})
}

func (c *transpileAndOptimizeClosureUnaware_CollectShadowAsts) isEmpty() bool {
	return len(c.groups) == 0
}

func (c *transpileAndOptimizeClosureUnaware_CollectShadowAsts) values() []*transpileAndOptimizeClosureUnaware_ShadowAst {
	var out []*transpileAndOptimizeClosureUnaware_ShadowAst
	for _, g := range c.groups {
		out = append(out, g.asts...)
	}
	return out
}

// Java: TranspileAndOptimizeClosureUnaware.CollectShadowAsts.shouldTraverse(NodeTraversal, Node, Node)
func (c *transpileAndOptimizeClosureUnaware_CollectShadowAsts) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if !n.IsScript() {
		return true
	}
	if !n.IsClosureUnawareCode() {
		return false
	}
	fileJsDoc := n.GetJSDocInfo()
	if fileJsDoc == nil || fileJsDoc.GetPerFileClosureUnawareMode() == nil {
		// TODO: b/475302905 - remove this null check, once the JSC release to always serialize a
		// ClosureUnawareMode lands.
		return true
	}
	minificationLevel := *fileJsDoc.GetPerFileClosureUnawareMode()
	switch minificationLevel {
	case rhino.JSDocInfo_PerFileClosureUnawareMode_UNSPECIFIED, rhino.JSDocInfo_PerFileClosureUnawareMode_SIMPLE:
		return true
	case rhino.JSDocInfo_PerFileClosureUnawareMode_WHITESPACE:
		return false
	}
	// Go-only: language-adapter — unreachable for a valid enum value: the
	// Java switch expression is exhaustive.
	panic(jutil.NewIllegalStateException())
}

// Java: TranspileAndOptimizeClosureUnaware.CollectShadowAsts.visit(NodeTraversal, Node, Node)
func (c *transpileAndOptimizeClosureUnaware_CollectShadowAsts) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	shadow := n.GetClosureUnawareShadow()
	if shadow.IsNil() {
		return
	}
	// CALL
	//   GETPROP .call
	//     NAME $jscomp_wrap_closure_unaware_code
	// NAME globalThis
	//   or
	// CALL
	//   GETPROP .call
	//     NAME $jscomp_wrap_closure_unaware_code
	// NAME undefined
	//   or
	// CALL
	//   NAME $jscomp_wrap_closure_unaware_code
	var callNode rhino.Node
	if n.GetParent().IsCall() {
		callNode = n.GetParent()
	} else {
		callNode = n.GetGrandparent()
	}
	jutil.CheckStateMsg1(callNode.IsCall(), "%s", callNode)
	c.add(
		t.GetInput().GetSourceFile(), newTranspileAndOptimizeClosureUnaware_ShadowAst(shadow, callNode, shadow.GetOnlyChild()))
}

// Java: TranspileAndOptimizeClosureUnaware.initShadowInputs(NestedCompilerRunner, LinkedHashMultimap<SourceFile, ShadowAst>)
func (p *TranspileAndOptimizeClosureUnaware) initShadowInputs(
	shadowCompiler *NestedCompilerRunner,
	groups []shadowFileGroup) {
	for _, g := range groups {
		for indexInFile, shadowAst := range g.asts {
			uniqueName := g.file.GetName() + ".shadow" + strconv.Itoa(indexInFile)
			shadowCompiler.AddScript(shadowAst.Script().Detach(), uniqueName)
		}
	}
}

// reattachShadowNodes re-attaches the shadow AST scripts, that were
// previously detached for compilation, to their associated ROOT nodes from
// the main AST.
//
// Java: TranspileAndOptimizeClosureUnaware.reattachShadowNodes(LinkedHashMultimap<SourceFile, ShadowAst>)
func (p *TranspileAndOptimizeClosureUnaware) reattachShadowNodes(
	inputs []*transpileAndOptimizeClosureUnaware_ShadowAst) {
	for _, shadowAst := range inputs {
		script := shadowAst.Script()
		jutil.CheckStateMsg1(script.HasOneChild(), "%s", script.ToStringTree())
		shadowAst.OriginalRoot().AddChildToFront(script.Detach())
		// Remove traces of the shadow compilation.
		script.SetInputId(nil)
		script.SetStaticSourceFile(nil)
	}
}

// Java: TranspileAndOptimizeClosureUnaware.setScriptFeaturesets(Collection<ShadowAst>)
func (p *TranspileAndOptimizeClosureUnaware) setScriptFeaturesets(inputs []*transpileAndOptimizeClosureUnaware_ShadowAst) {
	for _, shadowAst := range inputs {
		shadowScript := shadowAst.Script()
		collector := parsing.NewFeatureCollector()
		collector.VisitScript(shadowScript)
		shadowScript.PutProp(rhino.Node_FEATURE_SET, collector.AllFeatures())
	}
}

// transpileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig is the
// DefaultPassConfig with the closure-unaware runtime library injection
// after markUnnormalized.
type transpileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig struct {
	PassConfigDelegate

	injectRuntimeLibs *PassFactory
}

// Java: TranspileAndOptimizeClosureUnaware.RuntimeLibInjectionPassConfig(CompilerOptions, PassFactory)
func newTranspileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig(options *CompilerOptions, injectRuntimeLibs *PassFactory) *transpileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig {
	return &transpileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig{
		PassConfigDelegate: NewPassConfigDelegate(NewDefaultPassConfig(options)),
		injectRuntimeLibs:  injectRuntimeLibs,
	}
}

// Java: TranspileAndOptimizeClosureUnaware.RuntimeLibInjectionPassConfig.getTranspileOnlyPasses()
func (c *transpileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig) GetTranspileOnlyPasses() *PassListBuilder {
	passes := c.PassConfigDelegate.GetTranspileOnlyPasses()
	passes.AddAfter(c.injectRuntimeLibs, PassNames_MARK_UNNORMALIZED)
	return passes
}

// Java: TranspileAndOptimizeClosureUnaware.RuntimeLibInjectionPassConfig.getFinalizations()
func (c *transpileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig) GetFinalizations() *PassListBuilder {
	passes := c.PassConfigDelegate.GetFinalizations()
	passes.AddAfter(c.injectRuntimeLibs, PassNames_MARK_UNNORMALIZED)
	return passes
}

// Java: TranspileAndOptimizeClosureUnaware.createInjectRuntimeLibrariesPass(Collection<ShadowAst>)
func (p *TranspileAndOptimizeClosureUnaware) createInjectRuntimeLibrariesPass(shadowAsts []*transpileAndOptimizeClosureUnaware_ShadowAst) *PassFactory {
	var callSites []*InjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite
	for _, shadowAst := range shadowAsts {
		closureUnawareFunction :=
			NodeUtil_FindPreorder(shadowAst.Script(), (rhino.Node).IsFunction, func(childToVisit rhino.Node) bool { return true })
		callSites = append(callSites,
			NewInjectClosureUnawareRuntimeLibraries_ClosureUnawareCallSite(
				shadowAst.CallNode(), closureUnawareFunction))
	}
	return makePassFactory("injectClosureUnawareRuntimeLibs", func(compiler AbstractCompiler) CompilerPass {
		return NewInjectClosureUnawareRuntimeLibraries(p.original, compiler, callSites)
	})
}

// ManageClosureUnawareCode is the pass that "unwraps" closure-unaware code
// into a regular AST.
//
// ManageClosureUnawareCode.unwrap() should be run after all other passes have
// run, to unwrap the code and re-expose it to the code-printing stage of the
// compiler.
type ManageClosureUnawareCode struct {
	compiler AbstractCompiler
}

// Java: ManageClosureUnawareCode(AbstractCompiler)
func NewManageClosureUnawareCode(compiler AbstractCompiler) *ManageClosureUnawareCode {
	return &ManageClosureUnawareCode{compiler: compiler}
}

// Java: ManageClosureUnawareCode.unwrap(AbstractCompiler)
func ManageClosureUnawareCode_Unwrap(compiler AbstractCompiler) *ManageClosureUnawareCode {
	return NewManageClosureUnawareCode(compiler)
}

// Java: ManageClosureUnawareCode.process(Node, Node)
func (m *ManageClosureUnawareCode) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(m.compiler, root, NewManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode())
}

// TODO jameswr: Given that the CodePrinter supports printing out the shadow instead of the shadow
// host node, do we even need to revert the AST back to the original form at the end of
// compilation?

// ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode is the
// traversal that replaces each closure-unaware shadow host with the original
// code.
type ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode struct{}

// Java: UnwrapConcealedClosureUnawareCode()
func NewManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode() *ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode {
	return &ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode{}
}

// Java: UnwrapConcealedClosureUnawareCode.shouldTraverse(NodeTraversal, Node, Node)
func (u *ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if parent.IsNil() || !parent.IsScript() {
		return true // keep going
	}

	if parent.IsScript() && !parent.IsClosureUnawareCode() {
		return false
	}

	// Once inside a closureUnaware script, we want to traverse the entire thing to make sure we
	// find all the nodes marked as closure-unaware shadows.
	return true
}

// Java: UnwrapConcealedClosureUnawareCode.visit(NodeTraversal, Node, Node)
func (u *ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	u.TryUnwrapClosureUnawareShadowedCode(n)
}

// Java: UnwrapConcealedClosureUnawareCode.tryUnwrapClosureUnawareShadowedCode(Node)
func (u *ManageClosureUnawareCode_UnwrapConcealedClosureUnawareCode) TryUnwrapClosureUnawareShadowedCode(n rhino.Node) {
	shadowAstRoot := n.GetClosureUnawareShadow()
	if shadowAstRoot.IsNil() {
		return
	}

	// ROOT -> SCRIPT -> EXPR_RESULT -> CALL -> FUNCTION
	shadowScript := shadowAstRoot.GetOnlyChild()
	jutil.CheckStateMsg1(shadowScript.IsScript(), "%s", shadowScript)
	jutil.CheckStateMsg1(shadowScript.HasOneChild(), "%s", shadowScript)
	exprResult := shadowScript.GetOnlyChild()
	jutil.CheckStateMsg1(exprResult.IsExprResult(), "%s", exprResult)
	callNode := exprResult.GetOnlyChild()
	jutil.CheckStateMsg1(callNode.IsCall(), "%s", callNode)
	originalCodeFunction := callNode.GetLastChild()
	jutil.CheckStateMsg1(originalCodeFunction.IsFunction(), "%s", originalCodeFunction)
	originalCodeFunction.Detach()
	n.ReplaceWith(originalCodeFunction)
}
