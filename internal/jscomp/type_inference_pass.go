// Port of com/google/javascript/jscomp/TypeInferencePass.java @ closure-compiler v20261028

/*
 * Copyright 2009 The Closure Compiler Authors.
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
	"sort"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TypeInferencePass is a compiler pass to run the type inference analysis.
type TypeInferencePass struct {
	compiler                AbstractCompiler
	registry                *rhino.JSTypeRegistry
	reverseInterpreter      ReverseAbstractInterpreter
	topScope                *TypedScope
	scopeCreator            *TypedScopeCreator
	assertionFunctionLookup *CodingConvention_AssertionFunctionLookup
	reusableTypeInference   TypeInference
	reusableCfg             ReusableCfg

	// (stepCount, Token) -> populationCount
	stepCountHistogram map[int]*[256]int
}

// Java: TypeInferencePass(AbstractCompiler, ReverseAbstractInterpreter, TypedScopeCreator)
func NewTypeInferencePass(
	compiler AbstractCompiler,
	reverseInterpreter ReverseAbstractInterpreter,
	scopeCreator *TypedScopeCreator) *TypeInferencePass {
	p := &TypeInferencePass{}
	p.compiler = compiler
	p.registry = compiler.GetTypeRegistry()
	p.reverseInterpreter = reverseInterpreter
	p.scopeCreator = scopeCreator
	p.assertionFunctionLookup =
		CodingConvention_AssertionFunctionLookup_Of(compiler.GetCodingConvention().GetAssertionFunctions())
	if compiler.IsDebugLoggingEnabled() {
		p.stepCountHistogram = make(map[int]*[256]int)
	}
	return p
}

// InferAllScopes executes type inference running over part of the scope
// tree.
//
// Returns the top scope, either newly created, or patched by this
// inference.
//
// Java: TypeInferencePass.inferAllScopes(Node)
func (p *TypeInferencePass) InferAllScopes(inferenceRoot rhino.Node) *TypedScope {
	// Type analysis happens in two major phases.
	// 1) Finding all the symbols.
	// 2) Propagating all the inferred types.
	//
	// The order of this analysis is non-obvious. In a complete inference
	// system, we may need to backtrack arbitrarily far. But the compile-time
	// costs would be unacceptable.
	//
	// We do one pass where we do typed scope creation for all scopes
	// in pre-order.
	//
	// Then we do a second pass where we do all type inference
	// (type propagation) in pre-order.
	//
	// We use a memoized scope creator so that we never create a scope
	// more than once.
	//
	// This will allow us to handle cases like:
	// var ns = {};
	// (function() { /** JSDoc */ ns.method = function() {}; })();
	// ns.method();
	// In this code, we need to build the symbol table for the inner scope in
	// order to propagate the type of ns.method in the outer scope.
	var entries []typedScopeOrderEntry
	var scopes []*TypedScope
	func() {
		closer := p.registry.GetResolver().OpenForDefinition()
		defer closer.Close()
		jutil.CheckState(inferenceRoot.IsRoot())
		jutil.CheckState(inferenceRoot.GetParent().IsNil())
		jutil.CheckState(p.topScope == nil)
		p.scopeCreator.recordScopeOrder = inferenceRoot.HasTwoChildren()
		defer func() {
			p.scopeCreator.recordScopeOrder = false
			p.scopeCreator.scopeOrder = nil
		}()
		p.topScope = p.scopeCreator.CreateScope(inferenceRoot, nil).(*TypedScope)
		entries = p.scopeCreator.takeScopeOrder()

		if len(entries) > 0 {
			scopes = make([]*TypedScope, len(entries))
			scopes[0] = p.topScope
			p.buildAllScopesFromOrder(inferenceRoot, entries, scopes)
		} else {
			var firstCb TypeInferencePass_FirstScopeBuildingCallback
			b1 := NodeTraversal_Builder{
				compiler:     p.compiler,
				callback:     &firstCb,
				scopeCreator: p.scopeCreator,
			}
			b1.TraverseWithScope(inferenceRoot, p.topScope)
		}
		p.scopeCreator.ResolveWeakImportsPreResolution()
	}()
	p.scopeCreator.FinishAndFreeze()

	if len(entries) > 0 {
		p.inferAllScopesFromOrder(inferenceRoot, entries, scopes)
		entries = nil
		scopes = nil
	} else {
		secondCb := TypeInferencePass_SecondScopeBuildingCallback{outer: p}
		b2 := NodeTraversal_Builder{
			compiler:     p.compiler,
			callback:     &secondCb,
			scopeCreator: p.scopeCreator,
		}
		b2.TraverseWithScope(inferenceRoot, p.topScope)
	}

	// Normalize TypedVars to have the '?' type instead of null after inference is complete. This
	// currently cannot be done any earlier because it breaks inference of variables assigned in
	// local scopes.
	// TODO(b/149843534): this should be a crash instead.
	unknownType := p.registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	for var_ := range p.scopeCreator.GetAllSymbols() {
		if var_.GetType() == nil {
			var_.SetType(unknownType)
		}
	}

	if p.stepCountHistogram != nil {
		func() {
			histogram := p.compiler.CreateOrReopenLog("TypeInferencePass", "step_histogram.log")
			defer histogram.Close()
			histogram.LogString("step_count token population")

			totals := []int{0, 0}
			stepCounts := make([]int, 0, len(p.stepCountHistogram))
			for stepCount := range p.stepCountHistogram { // order-independent: stepCounts is sorted immediately below
				stepCounts = append(stepCounts, stepCount)
			}
			// Comparator.<Integer>naturalOrder().reversed()
			sort.SliceStable(stepCounts, func(i, j int) bool { return stepCounts[i] > stepCounts[j] })
			for _, stepCount := range stepCounts {
				type entry struct {
					element rhino.Token
					count   int
				}
				var entries []entry
				counts := p.stepCountHistogram[stepCount]
				for tok, count := range counts {
					if count > 0 {
						entries = append(entries, entry{rhino.Token(tok), count})
					}
				}
				// comparingInt(Multiset.Entry::getCount)
				sort.SliceStable(entries, func(i, j int) bool { return entries[i].count < entries[j].count })
				for _, e := range entries {
					totals[0] += stepCount * e.count
					totals[1] += e.count
					histogram.LogFormat("%s %s %s", stepCount, e.element, e.count)
				}
			}
			histogram.LogFormat("%s TOTAL %s", totals[0], totals[1])
		}()
	}

	return p.topScope
}

// Java: TypeInferencePass.inferScope(Node, TypedScope)
func (p *TypeInferencePass) InferScope(n rhino.Node, scope *TypedScope) {
	if p.tryInferEmptyFunctionScope(n, scope) {
		return
	}
	cfg := p.ComputeCfg(n)
	typeInference := &p.reusableTypeInference
	typeInference.initTypeInference(
		p.compiler, cfg, p.reverseInterpreter, scope, p.scopeCreator, p.assertionFunctionLookup)
	typeInference.Analyze()

	if p.stepCountHistogram != nil {
		for _, node := range cfg.NodesSlice() {
			if node == cfg.GetImplicitReturn() {
				continue
			}

			stepCount := typeInference.stepCountAt(node.Index)
			counts := p.stepCountHistogram[stepCount]
			if counts == nil {
				counts = new([256]int)
				p.stepCountHistogram[stepCount] = counts
			}
			counts[node.GetValue().GetToken()]++
		}
	}
	if n.IsRoot() {
		typeInference.ReleaseLargeBuffers()
	}
	p.reusableCfg.release()
}

func (p *TypeInferencePass) tryInferEmptyFunctionScope(n rhino.Node, scope *TypedScope) bool {
	if p.stepCountHistogram != nil || p.compiler == nil || scope == nil || scope.GetRootNode() != n {
		return false
	}
	fn, ok := ast.AsFunction(n)
	if !ok {
		return false
	}
	fnBody := fn.Body().Node()
	if !fnBody.IsBlock() || fnBody.HasChildren() {
		return false
	}
	if moduleMap := p.compiler.GetModuleMap(); moduleMap != nil {
		if ModuleImportResolver_GetModuleFromScopeRoot(moduleMap, p.compiler, n) != nil ||
			ModuleImportResolver_GetModuleFromScopeRoot(moduleMap, p.compiler, fnBody) != nil {
			return false
		}
	}
	if NodeUtil_IsBundledGoogModuleCall(n.GetParent()) {
		return true
	}
	params := fn.Params()
	paramsNode := params.Node()
	if !paramsNode.HasChildren() {
		return true
	}
	for c := paramsNode.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		switch c.GetToken() {
		case rhino.Token_NAME:
		case rhino.Token_ITER_REST:
			if first := c.GetFirstChild(); first.IsNil() || !first.IsName() {
				return false
			}
		default:
			return false
		}
	}

	unknownType := p.registry.GetNativeObjectType(rhino.JSTypeNative_UNKNOWN_TYPE)
	functionType := rhino.JSType_ToMaybeFunctionType(n.GetJSType())
	for zp := range zipFunctionParams(params, functionType.GetParameters(), firstIIFEArg(fn)) {
		inferredType := zp.param.Node().GetJSType()
		if inferredType == nil {
			inferredType = unknownType
		}
		if !zp.callArg.IsNil() {
			if argType := zp.callArg.GetJSType(); argType != nil {
				inferredType = argType
			}
		} else if paramType := zp.jsDocParam.GetJSType(); paramType != nil {
			inferredType = paramType
		}
		if zp.isRest {
			inferredType = p.registry.CreateTemplatizedType(
				p.registry.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE),
				inferredType,
			)
		}
		astParam := zp.target.Node()
		var_ := scope.GetVarByID(astParam.GetStringID())
		if var_ == nil {
			panic(jutil.NewNullPointerExceptionMsg(jutil.StringFormat("Missing var for parameter %s", astParam)))
		}
		astParam.SetJSType(inferredType)
		if var_.IsTypeInferred() {
			var_.SetType(inferredType)
		}
	}
	return true
}

// TypeInferencePass_FirstScopeBuildingCallback creates every typed scope, in
// pre-order.
type TypeInferencePass_FirstScopeBuildingCallback struct {
	NodeTraversal_AbstractScopedCallbackBase
}

func (*TypeInferencePass_FirstScopeBuildingCallback) nodeTraversalNoVisit() {}

// Java: TypeInferencePass.FirstScopeBuildingCallback.enterScope(NodeTraversal)
func (c *TypeInferencePass_FirstScopeBuildingCallback) EnterScope(t *NodeTraversal) {
	t.GetTypedScope()
}

// Java: TypeInferencePass.FirstScopeBuildingCallback.visit(NodeTraversal, Node, Node)
func (c *TypeInferencePass_FirstScopeBuildingCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	// Do nothing
}

// TypeInferencePass_SecondScopeBuildingCallback runs type inference on every
// CFG root scope, in pre-order.
type TypeInferencePass_SecondScopeBuildingCallback struct {
	NodeTraversal_AbstractScopedCallbackBase
	outer *TypeInferencePass
}

func (*TypeInferencePass_SecondScopeBuildingCallback) nodeTraversalNoVisit() {}

// Java: TypeInferencePass.SecondScopeBuildingCallback.enterScope(NodeTraversal)
func (c *TypeInferencePass_SecondScopeBuildingCallback) EnterScope(t *NodeTraversal) {
	// Only infer the entry root, rather than the scope root.
	// This ensures that incremental compilation only touches the root
	// that's been swapped out.
	scope := t.GetTypedScope()
	if scope.IsCfgRootScope() && !scope.IsModuleScope() {
		// ignore scopes that don't have their own CFGs and module scopes, which are visited
		// as if they were a regular script.
		c.outer.InferScope(t.GetCurrentNode(), scope)
	}
}

// Java: TypeInferencePass.SecondScopeBuildingCallback.visit(NodeTraversal, Node, Node)
func (c *TypeInferencePass_SecondScopeBuildingCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	// Do nothing
}

// Java: TypeInferencePass.computeCfg(Node)
func (p *TypeInferencePass) ComputeCfg(n rhino.Node) *ControlFlowGraph[rhino.Node] {
	return p.reusableCfg.Compute(
		p.compiler,
		n,
		/* shouldTraverseFunctions= */ false,
		/* edgeAnnotations= */ true)
}

var dummyTypeInferenceRecoverCallback TypeInferencePass_FirstScopeBuildingCallback

func inferAllScopesRecover(compiler AbstractCompiler, curNode, curScript rhino.Node, r any) {
	t := acquireNodeTraversalDirect(compiler, &dummyTypeInferenceRecoverCallback, nil, false)
	defer releaseNodeTraversal(t)
	t.currentNode = curNode
	t.currentScript = curScript
	t.ThrowUnexpectedException(jutil.AsThrowable(r))
}

func (p *TypeInferencePass) buildAllScopesFromOrder(
	inferenceRoot rhino.Node,
	entries []typedScopeOrderEntry,
	scopes []*TypedScope,
) {
	curNode := inferenceRoot
	var curScript rhino.Node
	defer p.scopeCreator.setCurrentBuildScript(rhino.NilNode)
	defer func() {
		if r := recover(); r != nil {
			inferAllScopesRecover(p.compiler, curNode, curScript, r)
		}
	}()
	if Platform_IsThreadInterrupted() {
		panic(jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException()))
	}
	var lastScript rhino.Node
	for i := 1; i < len(entries); i++ {
		e := entries[i]
		if e.script != lastScript {
			lastScript = e.script
			curScript = e.script
			curNode = e.script
			p.scopeCreator.setCurrentBuildScript(e.script)
			if Platform_IsThreadInterrupted() {
				panic(jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException()))
			}
		}
		curNode = e.currentNode
		curScript = e.script
		scopes[i] = p.scopeCreator.CreateScope(e.node, scopes[e.parentIdx]).(*TypedScope)
	}
}

func (p *TypeInferencePass) inferAllScopesFromOrder(
	inferenceRoot rhino.Node,
	entries []typedScopeOrderEntry,
	scopes []*TypedScope,
) {
	curNode := inferenceRoot
	var curScript rhino.Node
	defer func() {
		if r := recover(); r != nil {
			inferAllScopesRecover(p.compiler, curNode, curScript, r)
		}
	}()
	if Platform_IsThreadInterrupted() {
		panic(jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException()))
	}
	var lastScript rhino.Node
	for i, e := range entries {
		if i > 0 && e.script != lastScript {
			lastScript = e.script
			curScript = e.script
			curNode = e.script
			if Platform_IsThreadInterrupted() {
				panic(jutil.NewRuntimeExceptionWithCause(jutil.NewInterruptedException()))
			}
		}
		scope := scopes[i]
		if scope.IsCfgRootScope() && !scope.IsModuleScope() {
			curNode = e.currentNode
			curScript = e.script
			p.InferScope(e.currentNode, scope)
		}
	}
}
