// Port of com/google/javascript/jscomp/CoalesceVariableNames.java @ closure-compiler v20261015

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"cmp"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CoalesceVariableNames reuses variable names if possible.
//
// For example, from
//
//	var x = 1; print(x); var y = 2; print(y);
//
// to
//
//	var x = 1; print(x); x = 2; print(x)
//
// The benefits are slightly shorter code because of the removed var
// declaration, less unique variables in hope for better renaming, and
// finally better gzip compression.
//
// The pass operates similar to a typical register allocator found in an
// optimizing compiler by first computing live ranges with
// LiveVariablesAnalysis and a variable interference graph. Then it uses
// graph coloring in GraphColoring to determine which two variables can be
// merge together safely.
type CoalesceVariableNames struct {
	NodeTraversal_AbstractCfgCallbackBase

	compiler     AbstractCompiler
	scopeCreator *MemoizedScopeCreator
	// Go-only: language-adapter — the ArrayDeques used as stacks
	// (push/pop/peek at the front); the front is the last element of the
	// slice.
	colorings      []GraphColoring[*Var, any]
	liveAnalyses   []*LiveVariablesAnalysis
	usePseudoNames bool
	astFactory     *AstFactory
	liveness       *LiveVariablesAnalysis

	coloringTieBreaker func(a, b *Var) int

	// A stack of shouldOptimizeScope results.
	//
	// Go-only: language-adapter — see colorings.
	shouldOptimizeScopeStack []bool
}

// NewCoalesceVariableNames: usePseudoNames is for debug purposes, when
// merging variable foo and bar to foo, rename both variable to foo_bar.
//
// Java: CoalesceVariableNames(AbstractCompiler, boolean)
func NewCoalesceVariableNames(compiler AbstractCompiler, usePseudoNames bool) *CoalesceVariableNames {
	c := &CoalesceVariableNames{}
	c.InitNodeTraversal_AbstractCfgCallbackBase(c)
	c.coloringTieBreaker = func(a, b *Var) int {
		// comparingInt((Var arg) -> liveness.getVarIndex(arg.getName()))
		return cmp.Compare(c.liveness.GetVarIndex(a.GetName()), c.liveness.GetVarIndex(b.GetName()))
	}
	c.shouldOptimizeScopeStack = nil

	// The code is normalized at this point in the compilation process. This allows us to use the
	// fact that all variables have been given unique names. We can hoist coalesced variables to
	// VARS because we know that shadowing can't occur.
	jutil.CheckState(compiler.GetLifeCycleStage().IsNormalized())

	c.compiler = compiler
	c.colorings = nil
	c.liveAnalyses = nil
	c.usePseudoNames = usePseudoNames
	c.astFactory = compiler.CreateAstFactory()
	c.scopeCreator = NewMemoizedScopeCreator(NewSyntacticScopeCreator(compiler))
	return c
}

// Java: CoalesceVariableNames.process(Node, Node)
func (c *CoalesceVariableNames) Process(externs *rhino.Node, root *rhino.Node) {
	jutil.CheckNotNull(externs)
	jutil.CheckNotNull(root)
	NodeTraversal_Builder_().
		SetCompiler(c.compiler).
		SetCallback(c).
		SetScopeCreator(c.scopeCreator).
		Traverse(root)
	c.compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_RAW)
}

// CoalesceVariableNames_ShouldOptimizeScope returns populated
// AllVarsDeclaredInFunction object iff shouldOptimizeScope is true.
//
// Java: CoalesceVariableNames.shouldOptimizeScope(NodeTraversal)
func CoalesceVariableNames_ShouldOptimizeScope(t *NodeTraversal) *NodeUtil_AllVarsDeclaredInFunction {
	// TODO(user): We CAN do this in the global scope, just need to be
	// careful when something is exported. Liveness uses bit-vector for live
	// sets so I don't see compilation time will be a problem for running this
	// pass in the global scope.

	if t.GetScopeRoot().IsFunction() {
		allVarsDeclaredInFunction :=
			NodeUtil_GetAllVarsDeclaredInFunction(t.GetCompiler(), t.GetScopeCreator(), t.GetScope())
		if LiveVariablesAnalysis_MAX_VARIABLES_TO_ANALYZE >
			len(allVarsDeclaredInFunction.GetAllVariablesInOrder()) {
			return allVarsDeclaredInFunction
		}
	}

	return nil
}

// Java: CoalesceVariableNames.enterScopeWithCfg(NodeTraversal)
func (c *CoalesceVariableNames) EnterScopeWithCfg(t *NodeTraversal) {
	allVarsDeclaredInFunction := CoalesceVariableNames_ShouldOptimizeScope(t)
	if allVarsDeclaredInFunction == nil {
		c.shouldOptimizeScopeStack = append(c.shouldOptimizeScopeStack, false)
		return
	}
	c.shouldOptimizeScopeStack = append(c.shouldOptimizeScopeStack, true)

	scope := t.GetScope()
	jutil.CheckStateMsg1(scope.IsFunctionScope(), "%s", scope)

	// live variables analysis is based off of the control flow graph
	cfg := c.GetControlFlowGraph(c.compiler)

	c.liveness =
		NewLiveVariablesAnalysis(
			cfg, scope, nil, c.compiler, c.scopeCreator, allVarsDeclaredInFunction)

	if parser.FeatureSet_ES3.Contains(c.compiler.GetOptions().GetOutputFeatureSet()) {
		// If the function has exactly 2 params, mark them as escaped. This is a work-around for a
		// bug in IE 8 and below, where it throws an exception if you write to the parameters of the
		// callback in a sort(). See http://blickly.github.io/closure-compiler-issues/#58 and
		// https://www.zachleat.com/web/array-sort/
		enclosingFunction := scope.GetRootNode()
		if NodeUtil_GetFunctionParameters(enclosingFunction).HasTwoChildren() {
			c.liveness.MarkAllParametersEscaped()
		}
	}

	c.liveness.Analyze()
	c.liveAnalyses = append(c.liveAnalyses, c.liveness)

	// The interference graph has the function's variables as its nodes and any interference
	// between the variables as the edges. Interference between two variables means that they are
	// alive at overlapping times, which means that their variable names cannot be coalesced.
	interferenceGraph :=
		c.ComputeVariableNamesInterferenceGraph(cfg, c.liveness.GetEscapedLocals())

	// Color any interfering variables with different colors and any variables that can be safely
	// coalesced wih the same color.
	var coloring GraphColoring[*Var, any] = NewGraphColoring_GreedyGraphColoringWithTieBreaker[*Var, any](interferenceGraph, c.coloringTieBreaker)
	coloring.Color()
	c.colorings = append(c.colorings, coloring)
}

// Java: CoalesceVariableNames.exitScopeWithCfg(NodeTraversal)
func (c *CoalesceVariableNames) ExitScopeWithCfg(t *NodeTraversal) {
	if len(c.shouldOptimizeScopeStack) == 0 {
		panic(jutil.NewNoSuchElementException()) // ArrayDeque.pop()
	}
	shouldOptimize := c.shouldOptimizeScopeStack[len(c.shouldOptimizeScopeStack)-1]
	c.shouldOptimizeScopeStack = c.shouldOptimizeScopeStack[:len(c.shouldOptimizeScopeStack)-1]
	if !shouldOptimize {
		return
	}
	if len(c.colorings) == 0 || len(c.liveAnalyses) == 0 {
		panic(jutil.NewNoSuchElementException()) // ArrayDeque.pop()
	}
	c.colorings[len(c.colorings)-1] = nil
	c.colorings = c.colorings[:len(c.colorings)-1]
	c.liveAnalyses[len(c.liveAnalyses)-1] = nil
	c.liveAnalyses = c.liveAnalyses[:len(c.liveAnalyses)-1]
	// ArrayDeque.peek(): null when empty.
	if len(c.liveAnalyses) == 0 {
		c.liveness = nil
	} else {
		c.liveness = c.liveAnalyses[len(c.liveAnalyses)-1]
	}
}

// Java: CoalesceVariableNames.visit(NodeTraversal, Node, Node)
func (c *CoalesceVariableNames) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if len(c.colorings) == 0 || !n.IsName() || parent.IsFunction() {
		// Don't rename named functions.
		return
	}

	coloringsPeek := c.colorings[len(c.colorings)-1]
	var_ := c.liveness.GetAllVariables().Get(n.GetString())
	vNode := coloringsPeek.GetGraph().GetNode(var_)
	if vNode == nil {
		// This is not a local.
		return
	}
	coalescedVar := coloringsPeek.GetPartitionSuperNode(var_)

	if !c.usePseudoNames {
		if vNode.GetValue().Equals(coalescedVar) {
			// The coalesced name is itself, nothing to do.
			return
		}

		// Rename.
		n.SetString(coalescedVar.GetName())
		c.compiler.ReportChangeToEnclosingScope(n)
		c.UpdateDeclarationsPostCoalescing(n, coalescedVar, var_, parent)
	} else {
		// This code block is slow but since usePseudoName is for debugging,
		// we should not sacrifice performance for non-debugging compilation to
		// make this fast.
		allMergedNames := jutil.NewTreeSet[string](jutil.StringCompareTo)
		for _, iVar := range c.liveness.GetAllVariablesInOrder() {
			// Look for all the variables that can be merged (in the graph by now)
			// and it is merged with the current coalescedVar.
			if coloringsPeek.GetGraph().GetNode(iVar) != nil &&
				coloringsPeek.HaveSameColor(iVar, coalescedVar) {
				allMergedNames.Add(iVar.GetName())
			}
		}

		// Keep its original name.
		if allMergedNames.Size() == 1 {
			return
		}

		pseudoName := strings.Join(allMergedNames.ToSlice(), "_")

		for t.GetScope().HasSlot(pseudoName) {
			pseudoName += "$"
		}

		// Rename.
		n.SetString(pseudoName)
		c.compiler.ReportChangeToEnclosingScope(n)

		if vNode.GetValue().Equals(coalescedVar) {
			return
		}
		c.UpdateDeclarationsPostCoalescing(n, coalescedVar, var_, parent)
	}
}

// UpdateDeclarationsPostCoalescing updates declarations of the given name
// node and the coalesced variable. Firstly, it converts the coalesced
// variable's declaration into a `var` if it is a `const` or a `let`.
// Secondly, it removes the declaration of the given name node if it is under
// a parent which is a name declaration or if the given name is a
// destructuring LHS name under a destructuring LHS declaration.
//
// For example, when coalescing the names `a` and `b` below, it changes:
//
//	const a = 1;
//	const [{prop: b} = {prop: undefined}] = [{prop: 10}]; // declares `b`
//
// to
//
//	var a = 1;
//	[{prop: a} = {prop: undefined}] = [{prop: 10}]; // coalescing renames `b` to `a` and this
//	                                                // method undeclares it by removing the
//	                                                // `const`
//
// n is the name node which is getting coalesced with the coalescedVar;
// coalescedVar the first declared variable of a group that is being
// coalesced; originalVar the original variable being coalesced; parent the
// parent node of the name node n.
//
// Java: CoalesceVariableNames.updateDeclarationsPostCoalescing(Node, Var, Var, Node)
func (c *CoalesceVariableNames) UpdateDeclarationsPostCoalescing(
	n *rhino.Node, coalescedVar *Var, originalVar *Var, parent *rhino.Node) {
	jutil.CheckStateMsg(n.IsName(), "trying to update the declaration of a non-name node")
	if NodeUtil_IsNameDeclaration(parent) ||
		(NodeUtil_GetEnclosingType(n, rhino.Token_DESTRUCTURING_LHS) != nil &&
			NodeUtil_IsLhsByDestructuring(n) &&
			// Consider this to be a destructuring LHS declaration only if we don't hit an
			// assignment like `[a,b] = [3,4]` first.
			!c.IsNameInsideDestructuringAssignment(n)) {
		// convert the coalesced variable's declaration into a `var` if it is a `const` or a `let`
		c.MakeDeclarationVar(coalescedVar)
		// remove the declaration of the given name node as it has been coalesced with coalescedVar
		c.RemoveVarDeclaration(n, originalVar)
	}
}

// IsNameInsideDestructuringAssignment: is the given name node inside the LHS
// of a destructuring pattern assignment like `[a,b,c] = [1,2,3]`.
//
// Importantly, for this pass, it means that the name is not inside a
// declaration like `const [a,b,c] = [1,2,3]` and must not be undeclared.
//
// Java: CoalesceVariableNames.isNameInsideDestructuringAssignment(Node)
func (c *CoalesceVariableNames) IsNameInsideDestructuringAssignment(n *rhino.Node) bool {
	enclosingTarget := NodeUtil_GetRootTarget(n)
	if enclosingTarget.IsDestructuringPattern() && enclosingTarget.GetParent().IsAssign() {
		// This is a destructuring pattern assignment.
		return true
	}
	return false
}

// ComputeVariableNamesInterferenceGraph: in order to determine when it is
// appropriate to coalesce two variables, we use a live variables analysis to
// make sure they are not alive at the same time. We take every CFG node and
// determine which pairs of variables are alive at the same time. These pairs
// are set to true in a bit map. We take every pairing of variables and use
// the bit map to check if the two variables are alive at the same time. If
// two variables are alive at the same time, we create an edge between them
// in the interference graph. The interference graph is the input to a graph
// coloring algorithm that ensures any interfering variables are marked in
// different color groups, while variables that can safely be coalesced are
// assigned the same color group.
//
// escaped: we don't want to coalesce any escaped variables. Returns graph
// with variable nodes and edges representing variable interference.
//
// Java: CoalesceVariableNames.computeVariableNamesInterferenceGraph(ControlFlowGraph<Node>, Set<? extends Var>)
func (c *CoalesceVariableNames) ComputeVariableNamesInterferenceGraph(
	cfg ControlFlowGraph[*rhino.Node], escaped *jutil.ObjectLinkedHashSet[*Var]) UndiGraph[*Var, any] {
	var interferenceGraph UndiGraph[*Var, any] = LinkedUndirectedGraph_Create[*Var, any]()

	// First create a node for each non-escaped variable. We add these nodes in the order in which
	// they appear in the code because we want the names that appear earlier in the code to be used
	// when coalescing to variables that appear later in the code.
	orderedVariables := append([]*Var(nil), c.liveness.GetAllVariablesInOrder()...)

	// index i in interferenceGraphNodes is set to true when interferenceGraph
	// has node orderedVariables[i]
	interferenceGraphNodes := jutil.NewBitSet()

	// interferenceBitSet[i] = indices of all variables that should have an edge with
	// orderedVariables[i]
	interferenceBitSet := make([]*jutil.BitSet, len(orderedVariables))
	for i := range interferenceBitSet {
		interferenceBitSet[i] = jutil.NewBitSet()
	}

	vIndex := -1
	for _, v := range orderedVariables {
		vIndex++
		if escaped.Contains(v) {
			continue
		}

		// NOTE(user): In theory, we CAN coalesce function names just like any variables. Our
		// Liveness analysis captures this just like it as described in the specification. However, we
		// saw some zipped and unzipped size increase after this. We are not totally sure why
		// that is but, for now, we will respect the dead functions and not play around with it
		if v.GetParentNode().IsFunction() {
			continue
		}

		// NOTE: we skip class declarations for a combination of two reasons:
		// 1. they are block-scoped, so we would need to rewrite them as class expressions
		//      e.g. `class C {}` -> `var C = class {}` to avoid incorrect semantics
		//      (see testDontCoalesceClassDeclarationsWithDestructuringDeclaration).
		//    This is possible but increases pre-gzip code size and complexity.
		// 2. since function declaration coalescing seems to cause a size regression (as discussed
		//    above) we assume that coalescing class names may cause a similar size regression.
		if v.GetParentNode().IsClass() {
			continue
		}

		// Skip lets and consts that have multiple variables declared in them, otherwise this produces
		// incorrect semantics. See test case "testCapture".
		// Skipping vars technically isn't needed for correct semantics, but works around a Safari
		// bug for var redeclarations (https://github.com/google/closure-compiler/issues/3164)
		if c.IsInMultipleLvalueDecl(v) {
			continue
		}

		interferenceGraph.CreateNode(v)
		interferenceGraphNodes.Set(vIndex)
	}

	// Go through every CFG node in the program and look at variables that are live.
	// Set the pair of live variables in interferenceBitSet so we can add an edge between them.
	for _, cfgNode := range cfg.GetDiGraphNodes() {
		if cfg.IsImplicitReturn(cfgNode) {
			continue
		}

		state := dataFlowAnalysis_linearFlowStateOf[*LiveVariablesAnalysis_LiveVariableLattice](cfgNode)

		// Check the live states and add edge when possible. An edge between two variables
		// means that they are alive at overlapping times, which means that their
		// variable names cannot be coalesced.
		livein := state.GetIn()
		for i := livein.NextSetBit(0); i >= 0; i = livein.NextSetBit(i + 1) {
			for j := livein.NextSetBit(i); j >= 0; j = livein.NextSetBit(j + 1) {
				interferenceBitSet[i].Set(j)
			}
		}
		liveout := state.GetOut()
		for i := liveout.NextSetBit(0); i >= 0; i = liveout.NextSetBit(i + 1) {
			for j := liveout.NextSetBit(i); j >= 0; j = liveout.NextSetBit(j + 1) {
				interferenceBitSet[i].Set(j)
			}
		}

		liveRangeChecker :=
			NewCoalesceVariableNames_LiveRangeChecker(cfgNode.GetValue(), orderedVariables, state)
		liveRangeChecker.Check(cfgNode.GetValue())
		liveRangeChecker.SetCrossingVariables(interferenceBitSet)
	}

	// Go through each variable and try to connect them.
	v1Index := -1
	for _, v1 := range orderedVariables {
		v1Index++

		v2Index := -1
		for _, v2 := range orderedVariables {
			v2Index++
			// Skip duplicate pairs. Also avoid merging a variable with itself.
			if v1Index > v2Index {
				continue
			}

			if !interferenceGraphNodes.Get(v1Index) || !interferenceGraphNodes.Get(v2Index) {
				// Skip nodes that were not added. They are globals and escaped locals.
				continue
			}

			if (v1.IsParam() && v2.IsParam()) || interferenceBitSet[v1Index].Get(v2Index) {
				// Add an edge between variable pairs that are both parameters
				// because we don't want parameters to share a name.
				interferenceGraph.ConnectIfNotFound(v1, nil, v2)
			}
		}
	}
	return interferenceGraph
}

// IsInMultipleLvalueDecl returns whether this variable's declaration also
// declares other names.
//
// For example, this would return true for `x` in `let [x, y, z] = []`;
//
// Java: CoalesceVariableNames.isInMultipleLvalueDecl(Var)
func (c *CoalesceVariableNames) IsInMultipleLvalueDecl(v *Var) bool {
	declarationType := v.DeclarationType()
	if declarationType == nil {
		// switch on a null enum throws.
		panic(jutil.NewNullPointerException())
	}
	switch *declarationType {
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_VAR:
		{
			nameDecl := NodeUtil_GetEnclosingNode(v.GetNode(), NodeUtil_IsNameDeclaration)

			count := 0 // for lambda access
			NodeUtil_VisitLhsNodesInNode(nameDecl, func(lhs *rhino.Node) { count++ })
			return count > 1
		}
	default:
		return false
	}
}

// CoalesceVariableNames_RemoveVarDeclaration removes variable declaration if
// the variable has been coalesced with another variable that has already
// been declared.
//
// A precondition is that if the variable has already been declared, it must
// be the only lvalue in said declaration. For example, this method will not
// accept `var x = 1, y = 2`. In theory we could leave in the `var`
// declaration, but var shadowing of params triggers a Safari bug:
// https://bugs.webkit.org/show_bug.cgi?id=182414 Another
//
// name is the name node of the variable being coalesced; originalVar the
// original variable being coalesced.
//
// Java: CoalesceVariableNames.removeVarDeclaration(Node, Var)
func (c *CoalesceVariableNames) RemoveVarDeclaration(name *rhino.Node, originalVar *Var) {
	var_ := NodeUtil_GetEnclosingNode(name, NodeUtil_IsNameDeclaration)
	parent := var_.GetParent()

	if var_.GetFirstChild().IsDestructuringLhs() {
		// convert `const [x] = arr` to `[x] = arr`
		// a precondition for this method is that `x` is the only lvalue in the destructuring pattern
		destructuringLhs := var_.GetFirstChild()
		pattern := destructuringLhs.RemoveFirstChild()
		if NodeUtil_IsEnhancedFor(parent) {
			var_.ReplaceWith(pattern)
		} else {
			rvalue := var_.GetFirstFirstChild().Detach()
			var_.ReplaceWith(NodeUtil_NewExpr(rhino.IR_Assign(pattern, rvalue).Srcref(var_)))
		}
	} else if NodeUtil_IsEnhancedFor(parent) {
		// convert `for (let x of ...` to `for (x of ...`
		var_.ReplaceWith(name.Detach())
	} else {
		// either `var x = 0;` or `var x;`
		jutil.CheckStateMsg1(var_.HasOneChild() && var_.GetFirstChild() == name, "%s", var_)
		if name.HasChildren() {
			// convert `let x = 0;` to `x = 0;`
			value := name.RemoveFirstChild()
			name.Detach()
			assign := rhino.IR_Assign(name, value).Srcref(name)

			// We don't need to wrapped it with EXPR node if it is within a FOR.
			if !parent.IsVanillaFor() {
				assign = NodeUtil_NewExpr(assign)
			}
			var_.ReplaceWith(assign)

		} else if c.IsUninitializedLetNameInLoopBodyAndLive(name, originalVar) {
			name.Detach()
			undefinedValue := c.astFactory.CreateUndefinedValue().SrcrefTree(name)
			assign := rhino.IR_Assign(name, undefinedValue).Srcref(name)
			if !parent.IsVanillaFor() {
				assign = NodeUtil_NewExpr(assign)
			}
			var_.ReplaceWith(assign)
		} else {
			// convert `let x;` to ``
			// and `for (let x;;) {}` to `for (;;) {}`
			// We can expect uninitialized declarations at this point and it's okay to remove them
			// and they've been coalesced with another declaration.
			NodeUtil_RemoveChild(parent, var_)
		}
	}
}

// MakeDeclarationVar converts `const` or `let` declarations to `var`
// declarations.
//
// This method should be called on the first declared variable of a group
// that are being coalesced.
//
// Because the code has already been normalized by the time this pass runs,
// we can safely redeclare any let and const coalesced variables as vars
//
// Java: CoalesceVariableNames.makeDeclarationVar(Var)
func (c *CoalesceVariableNames) MakeDeclarationVar(coalescedName *Var) {
	if coalescedName.IsConst() || coalescedName.IsLet() {
		nameNode := jutil.CheckNotNullMsg1(coalescedName.GetNameNode(), "%s", coalescedName)
		if CoalesceVariableNames_IsUninitializedLetNameInLoopBody(nameNode) {
			// We need to make sure that within a loop:
			//
			// `let x;`
			// becomes
			// `var x = void 0;`
			//
			// If we don't we won't be correctly resetting the variable to undefined on each loop
			// iteration once we turn it into a var declaration.
			//
			// Note that all other cases will already have an initializer.
			// const x = 1; // constant requires an initializer
			// let {x, y} = obj; // destructuring requires an initializer
			// let [x, y] = iterable; // destructuring requires an initializer
			undefinedValue := c.astFactory.CreateUndefinedValue().SrcrefTree(nameNode)
			nameNode.AddChildToFront(undefinedValue)
		}
		// find the declaration node in a way that works normal and destructuring declarations.
		declNode := NodeUtil_GetEnclosingNode(nameNode.GetParent(), NodeUtil_IsNameDeclaration)
		// normalization ensures that all variables in a function are uniquely named, so it's OK
		// to turn a `const` or `let` into a `var`.
		declNode.SetToken(rhino.Token_VAR)
	}
}

// Java: CoalesceVariableNames.isUninitializedLetNameInLoopBody(Node)
func CoalesceVariableNames_IsUninitializedLetNameInLoopBody(nameNode *rhino.Node) bool {
	jutil.CheckStateMsg1(nameNode.IsName(), "%s", nameNode)
	letNode := nameNode.GetParent()
	if !letNode.IsLet() {
		// We're looking for `let name;`
		// Note that in the case of destructuring an initializer always exists.
		// `let {name} = initializerRequiredHere;
		return false
	}
	if nameNode.HasOneChild() {
		// `let name = child;` has an initializer
		return false
	}

	letParent := letNode.GetParent()
	if NodeUtil_IsLoopStructure(letParent) {
		// `for (let x; ...`
		// `for (let x in ...`
		// `for (let x of ...`
		// `for await (let x of ...`
		// In all these cases the variable gets initialized on each loop iteration
		return false
	}
	// Inside a loop body, but not the loop control node itself
	return NodeUtil_IsWithinLoop(letParent)
}

// Java: CoalesceVariableNames.isUninitializedLetNameInLoopBodyAndLive(Node, Var)
func (c *CoalesceVariableNames) IsUninitializedLetNameInLoopBodyAndLive(nameNode *rhino.Node, originalVar *Var) bool {
	if !CoalesceVariableNames_IsUninitializedLetNameInLoopBody(nameNode) {
		return false
	}
	jutil.CheckState(c.liveness != nil)
	letNode := nameNode.GetParent()
	cfg := c.liveness.GetCfg()
	cfgNode := jutil.CheckNotNull(cfg.GetDiGraphNode(letNode))
	state := jutil.CheckNotNull(dataFlowAnalysis_linearFlowStateOf[*LiveVariablesAnalysis_LiveVariableLattice](cfgNode))
	varIndex := c.liveness.GetVarIndex(originalVar.GetName())
	return varIndex >= 0 && state.GetOut().IsLive(varIndex)
}

// CoalesceVariableNames_LiveRangeChecker is used to find written and read
// variables in the same CFG node so that the variable pairs can be marked as
// interfering in an interference bit map. Indices of written and read
// variables are put in a list. These two lists are used to mark each written
// variable as "crossing" all read variables.
type CoalesceVariableNames_LiveRangeChecker struct {
	root             *rhino.Node
	state            *DataFlowAnalysis_LinearFlowState[*LiveVariablesAnalysis_LiveVariableLattice]
	orderedVariables []*Var
	isAssignToList   []int // indices of written variables
	isReadFromList   []int // indices of read variables
}

// Java: CoalesceVariableNames.LiveRangeChecker(Node, Var[], LinearFlowState<LiveVariableLattice>)
func NewCoalesceVariableNames_LiveRangeChecker(
	root *rhino.Node, orderedVariables []*Var, state *DataFlowAnalysis_LinearFlowState[*LiveVariablesAnalysis_LiveVariableLattice]) *CoalesceVariableNames_LiveRangeChecker {
	return &CoalesceVariableNames_LiveRangeChecker{
		root:             root,
		orderedVariables: orderedVariables,
		state:            state,
	}
}

// Java: CoalesceVariableNames.LiveRangeChecker.check(Node)
func (l *CoalesceVariableNames_LiveRangeChecker) Check(n *rhino.Node) {
	// For most AST nodes, traverse the subtree in postorder because that's how the expressions
	// are evaluated.
	if n == l.root || !ControlFlowGraph_IsEnteringNewCfgNode(n) {
		if (n.IsDestructuringLhs() && n.HasTwoChildren()) ||
			(n.IsAssign() && n.GetFirstChild().IsDestructuringPattern()) ||
			n.IsDefaultValue() {
			// Evaluate the rhs of a destructuring assignment/declaration before the lhs.
			l.Check(n.GetSecondChild())
			l.Check(n.GetFirstChild())
		} else {
			for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
				l.Check(c)
			}
		}
		if CoalesceVariableNames_LiveRangeChecker_ShouldVisit(n) {
			l.Visit(n, n.GetParent())
		}
	}
}

// Java: CoalesceVariableNames.LiveRangeChecker.visit(Node, Node)
func (l *CoalesceVariableNames_LiveRangeChecker) Visit(n *rhino.Node, parent *rhino.Node) {
	for iVar := 0; iVar < len(l.orderedVariables); iVar++ {
		if CoalesceVariableNames_LiveRangeChecker_IsAssignTo(l.orderedVariables[iVar], n, parent) {
			l.isAssignToList = append(l.isAssignToList, iVar)
		}
	}
	if len(l.isAssignToList) != 0 {
		for iVar := 0; iVar < len(l.orderedVariables); iVar++ {
			varOutLive := l.state.GetOut().IsLive(iVar)
			if varOutLive || CoalesceVariableNames_LiveRangeChecker_IsReadFrom(l.orderedVariables[iVar], n) {
				l.isReadFromList = append(l.isReadFromList, iVar)
			}
		}
	}
}

// Java: CoalesceVariableNames.LiveRangeChecker.setCrossingVariables(BitSet[])
func (l *CoalesceVariableNames_LiveRangeChecker) SetCrossingVariables(interferenceBitSet []*jutil.BitSet) {
	for _, iWrittenVar := range l.isAssignToList {
		for _, iReadVar := range l.isReadFromList {
			interferenceBitSet[iWrittenVar].Set(iReadVar)
			interferenceBitSet[iReadVar].Set(iWrittenVar)
		}
	}
}

// CoalesceVariableNames_LiveRangeChecker_ShouldVisit returns whether any
// LiveRangeChecker would be interested in the node.
//
// Java: CoalesceVariableNames.LiveRangeChecker.shouldVisit(Node)
func CoalesceVariableNames_LiveRangeChecker_ShouldVisit(n *rhino.Node) bool {
	return (n.IsName() || (n.HasChildren() && n.GetFirstChild().IsName()))
}

// Java: CoalesceVariableNames.LiveRangeChecker.isAssignTo(Var, Node, Node)
func CoalesceVariableNames_LiveRangeChecker_IsAssignTo(var_ *Var, n *rhino.Node, parent *rhino.Node) bool {
	if n.IsName() {
		if parent.IsParamList() {
			// In a function declaration, the formal parameters are assigned.
			return var_.GetName() == n.GetString()
		} else if NodeUtil_IsNameDeclaration(parent) && n.HasChildren() {
			// If this is a VAR declaration, if the name node has a child, we are
			// assigning to that name.
			return var_.GetName() == n.GetString()
		} else if NodeUtil_IsLhsByDestructuring(n) {
			return var_.GetName() == n.GetString()
		}
	} else if NodeUtil_IsAssignmentOp(n) {
		// Lastly, any assignmentOP is also an assign.
		name := n.GetFirstChild()
		return name.IsName() && var_.GetName() == name.GetString()
	}
	return false // Definitely a read.
}

// Java: CoalesceVariableNames.LiveRangeChecker.isReadFrom(Var, Node)
func CoalesceVariableNames_LiveRangeChecker_IsReadFrom(var_ *Var, name *rhino.Node) bool {
	return name.IsName() &&
		var_.GetName() == name.GetString() &&
		!NodeUtil_IsNameDeclOrSimpleAssignLhs(name, name.GetParent())
}
