// Port of com/google/javascript/jscomp/TranspileAndOptimizeClosureUnaware.java @ closure-compiler v20261021

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
	if collector.shadowAsts.IsEmpty() {
		return
	}

	p.setScriptFeaturesets(collector.shadowAsts.Values())
	shadowOptions := ClosureUnawareOptions_Convert(p.original.GetOptions())
	shadowPassConfig :=
		newTranspileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig(
			shadowOptions, p.createInjectRuntimeLibrariesPass(collector.shadowAsts.Values()))
	shadowCompiler :=
		NestedCompilerRunner_Create(p.original, shadowOptions, p.mode, shadowPassConfig)

	p.initShadowInputs(shadowCompiler, collector.shadowAsts)
	shadowCompiler.Compile()
	p.reattachShadowNodes(collector.shadowAsts)
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

// transpileAndOptimizeClosureUnaware_CollectShadowAsts traverses all
// closure unaware SCRIPTs to collect shadow ASTs.
type transpileAndOptimizeClosureUnaware_CollectShadowAsts struct {
	// Go-only: language-adapter — LinkedHashMultimap<SourceFile, ShadowAst>.
	// SetMultimap iterates entries grouped by key, LinkedHashMultimap in
	// insertion order; they agree here because the traversal visits the
	// shadow ASTs of each file contiguously.
	shadowAsts *jutil.SetMultimap[*SourceFile, *transpileAndOptimizeClosureUnaware_ShadowAst]
}

// Java: TranspileAndOptimizeClosureUnaware.CollectShadowAsts()
func newTranspileAndOptimizeClosureUnaware_CollectShadowAsts() *transpileAndOptimizeClosureUnaware_CollectShadowAsts {
	return &transpileAndOptimizeClosureUnaware_CollectShadowAsts{
		shadowAsts: jutil.NewSetMultimap[*SourceFile, *transpileAndOptimizeClosureUnaware_ShadowAst](),
	}
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
	c.shadowAsts.Put(
		t.GetInput().GetSourceFile(), newTranspileAndOptimizeClosureUnaware_ShadowAst(shadow, callNode, shadow.GetOnlyChild()))
}

// Java: TranspileAndOptimizeClosureUnaware.initShadowInputs(NestedCompilerRunner, LinkedHashMultimap<SourceFile, ShadowAst>)
func (p *TranspileAndOptimizeClosureUnaware) initShadowInputs(
	shadowCompiler *NestedCompilerRunner,
	perFileInputs *jutil.SetMultimap[*SourceFile, *transpileAndOptimizeClosureUnaware_ShadowAst]) {
	for sourceFile := range perFileInputs.KeySet() {
		indexInFile := 0
		for shadowAst := range perFileInputs.Get(sourceFile).All() {
			uniqueName := sourceFile.GetName() + ".shadow" + strconv.Itoa(indexInFile)
			indexInFile++
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
	inputs *jutil.SetMultimap[*SourceFile, *transpileAndOptimizeClosureUnaware_ShadowAst]) {
	for _, shadowAst := range inputs.Values() {
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
	PassConfig_PassConfigDelegateBase

	injectRuntimeLibs *PassFactory
}

// Java: TranspileAndOptimizeClosureUnaware.RuntimeLibInjectionPassConfig(CompilerOptions, PassFactory)
func newTranspileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig(options *CompilerOptions, injectRuntimeLibs *PassFactory) *transpileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig {
	c := &transpileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig{}
	c.self = c
	c.InitPassConfig_PassConfigDelegateBase(NewDefaultPassConfig(options))
	c.injectRuntimeLibs = injectRuntimeLibs
	return c
}

// Java: TranspileAndOptimizeClosureUnaware.RuntimeLibInjectionPassConfig.getTranspileOnlyPasses()
func (c *transpileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig) GetTranspileOnlyPasses() *PassListBuilder {
	passes := c.PassConfig_PassConfigDelegateBase.GetTranspileOnlyPasses()
	passes.AddAfter(c.injectRuntimeLibs, PassNames_MARK_UNNORMALIZED)
	return passes
}

// Java: TranspileAndOptimizeClosureUnaware.RuntimeLibInjectionPassConfig.getFinalizations()
func (c *transpileAndOptimizeClosureUnaware_RuntimeLibInjectionPassConfig) GetFinalizations() *PassListBuilder {
	passes := c.PassConfig_PassConfigDelegateBase.GetFinalizations()
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
	return PassFactory_Builder_().
		SetName("injectClosureUnawareRuntimeLibs").
		SetInternalFactory(
			func(compiler AbstractCompiler) CompilerPass {
				return NewInjectClosureUnawareRuntimeLibraries(p.original, compiler, callSites)
			}).
		Build()
}
