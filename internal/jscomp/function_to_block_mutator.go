// Port of com/google/javascript/jscomp/FunctionToBlockMutator.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: function, loop, declaration, and return shapes are read through
// package ast while keeping FunctionToBlockMutator's rhino.Node API intact.

// FunctionToBlockMutator is a class to transform the body of a function into
// a generic block suitable for inlining.
type FunctionToBlockMutator struct {
	compiler                 AbstractCompiler
	safeNameIdSupplier       func() string
	functionArgumentInjector *FunctionArgumentInjector
	localScopeCreator        SyntacticScopeCreator
	localMakeUnique          MakeDeclaredNamesUnique
	localRenamerStackBuf     [16]MakeDeclaredNamesUnique_Renamer
	allowMoveCallArgs        bool
	currentCallNode          rhino.Node
	onMovedArg               func(rhino.Node)
	errorNodes               map[rhino.Node]struct{}
	lastDiagCount            int
	paramNamesScratch        inlineIDSet
	modifiedParamsScratch    inlineIDSet
	argsScratch              []callParamBinding
	allNamesToAliasScratch   inlineIDSet
	paramReplacementsScratch inlineIDNodeMap
	newAliasesToAddScratch   []rhino.Node
}

// Java: FunctionToBlockMutator(AbstractCompiler, Supplier<String>)
func NewFunctionToBlockMutator(compiler AbstractCompiler, safeNameIdSupplier func() string) *FunctionToBlockMutator {
	m := &FunctionToBlockMutator{}
	m.compiler = jutil.CheckNotNull(compiler)
	m.safeNameIdSupplier = safeNameIdSupplier
	m.functionArgumentInjector = NewFunctionArgumentInjector(compiler.GetAstAnalyzer())
	return m
}

// Mutate returns a clone of the function body mutated to be suitable for
// injection as a statement into another code block.
//
// fnName is the name to use when preparing human readable names; fnNode the
// function to prepare; callNode the call node that will be replaced;
// resultName the name function results should be assigned to (empty: Java
// null); needsDefaultResult whether the result value must be set;
// isCallInLoop whether the function body must be prepared to be injected
// into the body of a loop.
//
// Java: FunctionToBlockMutator.mutate(String, Node, Node, String, boolean, boolean)
func (m *FunctionToBlockMutator) Mutate(
	fnName string,
	fnNode rhino.Node,
	callNode rhino.Node,
	resultName string,
	needsDefaultResult bool,
	isCallInLoop bool) rhino.Node {
	return m.mutateWithSummary(
		fnName, fnNode, callNode, resultName, needsDefaultResult, isCallInLoop, nil)
}

func (m *FunctionToBlockMutator) mutateWithSummary(
	fnName string,
	fnNode rhino.Node,
	callNode rhino.Node,
	resultName string,
	needsDefaultResult bool,
	isCallInLoop bool,
	summary *calleeSummary,
) rhino.Node {
	if m.allowMoveCallArgs {
		m.currentCallNode = callNode
	}
	res := m.mutateInternalWithSummary(
		fnName,
		fnNode,
		callNode,
		resultName,
		needsDefaultResult,
		isCallInLoop,
		/* renameLocals= */ true,
		summary)
	m.currentCallNode = rhino.NilNode
	return res
}

// MutateWithoutRenaming is used where the inlining occurs into an isolated
// scope such as a module. Renaming is avoided since the associated JSDoc
// annotations are not updated.
//
// fnName is the name to use when preparing human readable names; fnNode the
// function to prepare; callNode the call node that will be replaced;
// resultName the name function results should be assigned to;
// needsDefaultResult whether the result value must be set; isCallInLoop
// whether the function body must be prepared to be injected into the body
// of a loop. Returns a clone of the function body mutated to be suitable for
// injection as a statement into another code block.
//
// Java: FunctionToBlockMutator.mutateWithoutRenaming(String, Node, Node, String, boolean, boolean)
func (m *FunctionToBlockMutator) MutateWithoutRenaming(
	fnName string,
	fnNode rhino.Node,
	callNode rhino.Node,
	resultName string,
	needsDefaultResult bool,
	isCallInLoop bool) rhino.Node {
	return m.MutateInternal(
		fnName,
		fnNode,
		callNode,
		resultName,
		needsDefaultResult,
		isCallInLoop,
		/* renameLocals= */ false)
}

// MutateInternal: fnName is the name to use when preparing human readable
// names; fnNode the function to prepare; callNode the call node that will be
// replaced; resultName the name function results should be assigned to;
// needsDefaultResult whether the result value must be set; isCallInLoop
// whether the function body must be prepared to be injected into the body
// of a loop; renameLocals false if the inlining is part of module rewriting
// and doesn't require making local names unique. Returns a clone of the
// function body mutated to be suitable for injection as a statement into
// another code block.
//
// Java: FunctionToBlockMutator.mutateInternal(String, Node, Node, String, boolean, boolean, boolean)
func (m *FunctionToBlockMutator) MutateInternal(
	fnName string,
	fnNode rhino.Node,
	callNode rhino.Node,
	resultName string,
	needsDefaultResult bool,
	isCallInLoop bool,
	renameLocals bool) rhino.Node {
	return m.mutateInternalWithSummary(
		fnName, fnNode, callNode, resultName, needsDefaultResult, isCallInLoop, renameLocals, nil)
}

func (m *FunctionToBlockMutator) mutateInternalWithSummary(
	fnName string,
	fnNode rhino.Node,
	callNode rhino.Node,
	resultName string,
	needsDefaultResult bool,
	isCallInLoop bool,
	renameLocals bool,
	summary *calleeSummary,
) rhino.Node {
	var localSummary calleeSummary
	if summary == nil || summary.fnNode != fnNode {
		if scanCalleeSummary(fnName, fnNode, &localSummary) {
			summary = &localSummary
		} else {
			summary = nil
		}
	}

	newFn := ast.MustFunction(fnNode.CloneTree())
	newFnNode := newFn.Node()
	// Wrap the clone in a script, in a root, so that makeLocalNamesUnique sees a coherent tree.
	rhino.IR_Root(rhino.IR_ScriptOfNodes(newFnNode))

	if renameLocals {
		// Now that parameter names have been replaced, make sure all the local
		// names are unique, to allow functions to be inlined multiple times
		// without causing conflicts.
		m.makeLocalNamesUnique(newFnNode, isCallInLoop, summary == nil || summary.hasLabels)

		// Function declarations must be rewritten as function expressions as
		// they will be within a block and normalization prevents function
		// declarations within block as browser implementations vary.
		if summary == nil || summary.hasInnerFunctionDecls {
			FunctionToBlockMutator_RewriteFunctionDeclarations(newFn.Body().Node())
		}
	}

	m.argsScratch = m.functionArgumentInjector.fillCallParamBindings(
		m.argsScratch[:0], newFnNode, callNode, m.safeNameIdSupplier)

	// TODO(johnlenz): Mark NAME nodes constant for parameters that are not modified.
	modifiedParameters := &m.modifiedParamsScratch
	if summary != nil && summary.canReuseAfterRename(renameLocals) {
		modifiedParameters.clear()
		if summary.numModified > 0 {
			bindings := m.argsScratch
			for _, paramIdx := range summary.modifiedParamIndices() {
				if int(paramIdx)+1 < len(bindings) {
					modifiedParameters.add(bindings[int(paramIdx)+1].nameID)
				}
			}
		}
	} else {
		m.functionArgumentInjector.findModifiedParametersByIDInto(newFnNode, &m.paramNamesScratch, modifiedParameters)
	}

	hasArgs := len(m.argsScratch) > 0
	allNamesToAlias := &m.allNamesToAliasScratch
	allNamesToAlias.clear()
	if hasArgs {
		m.functionArgumentInjector.gatherCallArgumentsNeedingTempsByID(
			m.compiler, newFnNode, m.argsScratch, modifiedParameters, m.compiler.GetCodingConvention(), summary, renameLocals, allNamesToAlias)
	}

	mayHaveInjectedDecls := false
	if isCallInLoop && summary != nil && !summary.hasVarOrLetDecls && !(renameLocals && summary.hasInnerFunctionDecls) {
		for _, b := range m.argsScratch {
			if b.argNode.HasChildren() && functionToBlockMutator_hasVarOrLetOutsideFn(b.argNode) {
				mayHaveInjectedDecls = true
				break
			}
		}
	}

	newBlock := newFn.Body().MustBlock().Detach().Node()

	if hasArgs {
		inlineResult := m.aliasAndInlineArgumentsByID(
			newBlock, m.argsScratch, allNamesToAlias, summary != nil, summary != nil && summary.referencesThis)
		jutil.CheckState(newBlock == inlineResult)
	}

	//
	// For calls inlined into loops, VAR declarations are not reinitialized to
	// undefined as they would have been if the function were called, so ensure
	// that they are properly initialized.
	//
	if isCallInLoop && (summary == nil || summary.hasVarOrLetDecls || (renameLocals && summary.hasInnerFunctionDecls) || mayHaveInjectedDecls) {
		FunctionToBlockMutator_FixUninitializedVarDeclarations(newBlock, newBlock)
	}

	labelName := m.GetLabelNameForFunction(fnName)
	var injectableBlock rhino.Node
	if summary != nil {
		injectableBlock = functionToBlockMutator_replaceReturnsWithCount(
			newBlock, resultName, labelName, needsDefaultResult, int(summary.shallowReturnCount))
	} else {
		injectableBlock = FunctionToBlockMutator_ReplaceReturns(newBlock, resultName, labelName, needsDefaultResult)
	}
	jutil.CheckState(!injectableBlock.IsNil())

	return injectableBlock
}

func functionToBlockMutator_hasVarOrLetOutsideFn(n rhino.Node) bool {
	if n.IsFunction() {
		return false
	}
	if n.IsVar() || n.IsLet() {
		return true
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if c.HasChildren() && functionToBlockMutator_hasVarOrLetOutsideFn(c) {
			return true
		}
	}
	return false
}

// FunctionToBlockMutator_RewriteFunctionDeclarations: n is the node to
// inspect.
//
// Java: FunctionToBlockMutator.rewriteFunctionDeclarations(Node)
func FunctionToBlockMutator_RewriteFunctionDeclarations(n rhino.Node) rhino.Node {
	if fn, ok := ast.AsFunction(n); ok {
		if NodeUtil_IsFunctionDeclaration(n) {
			// Rewrite: function f() {} ==> var f = function() {}
			fnName := fn.NameNode()

			name := rhino.IR_Name(fnName.Ident()).Srcref(fnName.Node())
			var_ := rhino.IR_Var(name).Srcref(n)

			fn.RemoveName()
			// Add the VAR, remove the FUNCTION
			n.ReplaceWith(var_)
			// readd the function as a function expression
			name.AddChildToFront(n)

			return var_
		}
		return rhino.NilNode
	}

	// Keep track of any rewritten functions and hoist them to the top
	// of the block they are defined in. This isn't fully compliant hoisting
	// but it does address a large set of use cases.
	var hoistBuf [4]rhino.Node
	functionsToHoist := hoistBuf[:0]
	var next rhino.Node
	for c := n.GetFirstChild(); !c.IsNil(); c = next {
		next = c.GetNext() // We may rewrite "c"
		fnVar := FunctionToBlockMutator_RewriteFunctionDeclarations(c)
		if !fnVar.IsNil() {
			functionsToHoist = append(functionsToHoist, fnVar)
		}
	}

	for i := len(functionsToHoist) - 1; i >= 0; i-- {
		fnVar := functionsToHoist[i]
		if n.GetFirstChild() != fnVar {
			n.AddChildToFront(fnVar.Detach())
		}
	}

	return rhino.NilNode
}

// FunctionToBlockMutator_FixUninitializedVarDeclarations: for all VAR
// declarations (and LET declarations outside loops) that are uninitialized,
// set their initial values to "undefined" and hoist them to the containing
// block.
//
// When a function is called within a loop, local variables are recreated on
// each call in normal execution. When inlined into the loop, uninitialized
// variables would otherwise retain their values across loop iterations.
//
// Java: FunctionToBlockMutator.fixUninitializedVarDeclarations(Node, Node)
func FunctionToBlockMutator_FixUninitializedVarDeclarations(n rhino.Node, containingBlock rhino.Node) {
	FunctionToBlockMutator_FixUninitializedVarDeclarationsWithInLoop(n, containingBlock /* inLoop= */, false)
}

// Java: FunctionToBlockMutator.fixUninitializedVarDeclarations(Node, Node, boolean)
func FunctionToBlockMutator_FixUninitializedVarDeclarationsWithInLoop(
	n rhino.Node, containingBlock rhino.Node, inLoop bool) {
	// Nested functions introduce there own var hoisting scope
	// thier `var` should not be hoisted out to the containing block.
	if _, ok := ast.AsFunction(n); ok {
		return
	}

	// Inner loop structures: var declarations inside loops are still function-scoped and
	// need to be initialized at the containing block level to avoid carrying values across
	// outer loop iterations. When hoisting, any assignment at the declaration site in the loop
	// must remain in the loop. Let declarations inside loops are block-scoped to that loop and
	// do not need hoisting.
	if loop, ok := ast.AsLoop(n); ok {
		body := loop.Body().Node()
		if !body.IsNil() {
			FunctionToBlockMutator_FixUninitializedVarDeclarationsWithInLoop(body, containingBlock /* inLoop= */, true)
		}
		if vanillaFor, isFor := loop.AsFor(); isFor {
			if init, hasInit := vanillaFor.Init(); hasInit {
				FunctionToBlockMutator_FixUninitializedVarDeclarationsWithInLoop(init.Node(), containingBlock /* inLoop= */, true)
			}
		}
		return
	}

	if decl, ok := ast.AsNameDeclaration(n); ok {
		if onlyDecl, hasOne := decl.OnlyDeclarator(); hasOne {
			if inLoop && decl.Token() == rhino.Token_VAR {
				name := onlyDecl.MustName()
				nameNode := name.Node()
				inForInit := ast.RoleOf(n) == ast.RoleForInit
				// Hoist per-function declaration initialized to undefined
				hoistedName := nameNode.CloneNode().SetJSDocInfo(nil)
				hoistedName.AddChildToBack(NodeUtil_NewUndefinedNode(nameNode))
				containingBlock.AddChildToFront(rhino.IR_Var(hoistedName).Srcref(n))

				if initVal, hasInit := name.Init(); hasInit {
					// If the var in the loop had an initializer, keep the assignment at the declaration site in
					// the loop.
					value := initVal.Detach().Node()
					assign := rhino.IR_Assign(nameNode.Detach(), value).Srcref(n)
					if inForInit {
						n.ReplaceWith(assign)
					} else {
						n.ReplaceWith(NodeUtil_NewExpr(assign))
					}
				} else {
					// If uninitialized in the loop, remove the var declaration from the loop.
					if inForInit {
						n.ReplaceWith(rhino.IR_Empty())
					} else {
						n.Detach()
					}
				}
				return
			}

			if !inLoop && (decl.Token() == rhino.Token_VAR || decl.Token() == rhino.Token_LET) {
				if name, isName := onlyDecl.AsName(); isName {
					// It isn't initialized.
					if _, hasInit := name.Init(); !hasInit {
						nameNode := name.Node()
						nameNode.AddChildToBack(NodeUtil_NewUndefinedNode(nameNode))
						if ast.RoleOf(n) == ast.RoleForInit {
							n.ReplaceWith(rhino.IR_Empty())
						} else {
							n.Detach()
						}
						containingBlock.AddChildToFront(n)
					}
				}
				return
			}
		}
	}

	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		FunctionToBlockMutator_FixUninitializedVarDeclarationsWithInLoop(c, containingBlock, inLoop)
	}
}

// MakeLocalNamesUnique fixes up all local names to be unique for this
// subtree.
//
// fnNode is a mutable instance of the function to be inlined.
//
// Java: FunctionToBlockMutator.makeLocalNamesUnique(Node, boolean)
func (m *FunctionToBlockMutator) MakeLocalNamesUnique(fnNode rhino.Node, isCallInLoop bool) {
	m.makeLocalNamesUnique(fnNode, isCallInLoop, true)
}

func (m *FunctionToBlockMutator) makeLocalNamesUnique(fnNode rhino.Node, isCallInLoop bool, checkLabels bool) {
	idSupplier := m.compiler.GetUniqueNameIdSupplier()
	rootRenamer := NewMakeDeclaredNamesUnique_InlineRenamer(
		m.compiler.GetCodingConvention(),
		idSupplier,
		"inline_",
		isCallInLoop,
		true,
		nil,
	)
	cb := &m.localMakeUnique
	cb.renamerStack = m.localRenamerStackBuf[:0]
	cb.rootRenamer = rootRenamer
	cb.markChanges = false
	cb.assertOnChange = false

	if m.localScopeCreator.compiler == nil {
		m.localScopeCreator = SyntacticScopeCreator{
			compiler:             m.compiler,
			redeclarationHandler: SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER,
		}
	} else {
		m.localScopeCreator.lastInputId = nil
		m.localScopeCreator.lastInput = nil
	}
	b := NodeTraversal_Builder{
		compiler:     m.compiler,
		callback:     cb,
		scopeCreator: &m.localScopeCreator,
	}
	b.traverseScopedRoot(fnNode, true)
	clear(cb.renamerStack)
	cb.renamerStack = nil
	cb.rootRenamer = nil

	// Make label names unique to this instance.
	if checkLabels && hasLabelNode(fnNode) {
		NewRenameLabelsWithSupplier(m.compiler, NewFunctionToBlockMutator_LabelNameSupplier(idSupplier).Get, false, false).
			Process(rhino.NilNode, fnNode)
	}
}

func hasLabelNode(n rhino.Node) bool {
	if n.IsLabel() {
		return true
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if hasLabelNode(c) {
			return true
		}
	}
	return false
}

// FunctionToBlockMutator_LabelNameSupplier is a Supplier<String>.
type FunctionToBlockMutator_LabelNameSupplier struct {
	idSupplier func() string
}

// Java: FunctionToBlockMutator.LabelNameSupplier(Supplier<String>)
func NewFunctionToBlockMutator_LabelNameSupplier(idSupplier func() string) *FunctionToBlockMutator_LabelNameSupplier {
	return &FunctionToBlockMutator_LabelNameSupplier{idSupplier: idSupplier}
}

// Java: FunctionToBlockMutator.LabelNameSupplier.get()
func (s *FunctionToBlockMutator_LabelNameSupplier) Get() string {
	return "JSCompiler_inline_label_" + s.idSupplier()
}

// GetLabelNameForFunction creates a unique label name.
//
// Java: FunctionToBlockMutator.getLabelNameForFunction(String)
func (m *FunctionToBlockMutator) GetLabelNameForFunction(fnName string) string {
	name := fnName
	if fnName == "" {
		name = "anon"
	}
	return "JSCompiler_inline_label_" + name + "_" + m.safeNameIdSupplier()
}

// GetUniqueThisName creates a unique "this" name.
//
// Java: FunctionToBlockMutator.getUniqueThisName()
func (m *FunctionToBlockMutator) GetUniqueThisName() string {
	return "JSCompiler_inline_this_" + m.safeNameIdSupplier()
}

func (m *FunctionToBlockMutator) rebindMovedArgFunctions(n rhino.Node, deletedFnsHolder rhino.Node) rhino.Node {
	if n.IsFunction() {
		newFn := n.CloneNode()
		newFn.AddChildrenToBack(n.RemoveChildren())
		if !n.GetParent().IsNil() {
			n.ReplaceWith(newFn)
		}
		if m.onMovedArg != nil {
			m.onMovedArg(n)
		}
		deletedFnsHolder.AddChildToBack(n)
		n = newFn
	} else if len(m.errorNodes) > 0 {
		if _, ok := m.errorNodes[n]; ok {
			cloned := n.CloneNode()
			cloned.AddChildrenToBack(n.RemoveChildren())
			if !n.GetParent().IsNil() {
				n.ReplaceWith(cloned)
			}
			n = cloned
		}
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		c = m.rebindMovedArgFunctions(c, deletedFnsHolder)
	}
	return n
}

func (m *FunctionToBlockMutator) takeOrCloneCallArg(origArg rhino.Node) rhino.Node {
	if !m.currentCallNode.IsNil() && origArg.GetParent() == m.currentCallNode {
		if m.compiler != nil {
			if em := m.compiler.GetErrorManager(); em != nil {
				if total := em.GetErrorCount() + em.GetWarningCount(); total > 0 && total != m.lastDiagCount {
					m.lastDiagCount = total
					if m.errorNodes == nil {
						m.errorNodes = make(map[rhino.Node]struct{}, total)
					}
					for _, e := range em.GetErrors() {
						if !e.Node().IsNil() {
							m.errorNodes[e.Node()] = struct{}{}
						}
					}
					for _, w := range em.GetWarnings() {
						if !w.Node().IsNil() {
							m.errorNodes[w.Node()] = struct{}{}
						}
					}
				}
			}
		}
		placeholder := rhino.IR_Empty()
		origArg.ReplaceWith(placeholder)
		origArg = m.rebindMovedArgFunctions(origArg, placeholder)
		if origArg.HasChildren() && m.onMovedArg != nil {
			m.onMovedArg(origArg)
		}
		return origArg
	}
	return origArg.CloneTree()
}

// AliasAndInlineArguments inlines the arguments within the node tree using
// the given argument map, replaces "unsafe" names with local aliases.
//
// The aliases for unsafe require new VAR declarations, so this function can
// not be used in for direct CALL node replacement as VAR nodes can not be
// created there.
//
// Returns the node or its replacement.
//
// Java: FunctionToBlockMutator.aliasAndInlineArguments(Node, ImmutableMap<String, ParamArgPair>, Set<String>)
func (m *FunctionToBlockMutator) AliasAndInlineArguments(
	fnTemplateRoot rhino.Node,
	paramToArgMap *jutil.LinkedHashMap[string, FunctionArgumentInjector_ParamArgPair],
	namesToAlias *jutil.LinkedHashSet[string]) rhino.Node {
	m.argsScratch = bindingsFromArgMap(m.argsScratch[:0], paramToArgMap)
	m.allNamesToAliasScratch.fromLinkedHashSet(namesToAlias, true)
	return m.aliasAndInlineArgumentsByID(fnTemplateRoot, m.argsScratch, &m.allNamesToAliasScratch, false, false)
}

func (m *FunctionToBlockMutator) aliasAndInlineArgumentsByID(
	fnTemplateRoot rhino.Node,
	bindings []callParamBinding,
	namesToAlias *inlineIDSet,
	hasSummary bool,
	knownReferencesThis bool,
) rhino.Node {
	paramReplacements := &m.paramReplacementsScratch
	paramReplacements.clear()
	if namesToAlias.isEmpty() {
		// There are no names to alias. Just inline the arguments directly.
		for _, b := range bindings {
			paramReplacements.put(b.nameID, b.argNode)
		}
		result := m.functionArgumentInjector.injectByID(m.compiler, fnTemplateRoot, rhino.NilNode, paramReplacements, true)
		jutil.CheckState(result == fnTemplateRoot)
		return result
	}

	// Create local alias of names that can not be safely used directly.
	//
	// A map from function parameter to the value it should be replaced with post-inlining.
	// This is a subset of bindings: we exclude parameters for which we create
	// an explicit alias.
	// Declare the aliases in the same order as the arguments are defined.
	newAliasesToAdd := m.newAliasesToAddScratch[:0]
	for _, b := range bindings {
		if !namesToAlias.contains(b.nameID) {
			// Replace references to the parameter with the argument directly.
			paramReplacements.put(b.nameID, b.argNode)
			continue
		}
		if b.nameID == calleeThisMarkerID {
			// Update "this". We only need to create an alias if "this" is referenced
			// and the value of "this" is not Token.THIS, or the value of "this" has side effects.
			originalValue := b.argNode
			replacement := originalValue
			if !originalValue.IsThis() {
				referencesThis := knownReferencesThis
				if !hasSummary {
					referencesThis = NodeUtil_ReferencesEnclosingReceiver(fnTemplateRoot)
				}
				if referencesThis || m.compiler.GetAstAnalyzer().MayHaveSideEffects(originalValue) {
					newName := m.GetUniqueThisName()
					newValue := m.takeOrCloneCallArg(originalValue)
					newNode := NodeUtil_NewVarNode(newName, newValue).SrcrefTreeIfMissing(newValue)
					newAliasesToAdd = append(newAliasesToAdd, newNode)
					replacement = rhino.IR_Name(newName).SrcrefTree(newValue)
				}
			}
			paramReplacements.put(calleeThisMarkerID, replacement)
		} else {
			// Add an alias for a given parameter/argument. For example, if inlining this function:
			//
			//  function f(x) { use(x, x); } f(computeValue());
			//
			// This code defines
			//    var x = computeValue();
			// (as use(computeValue(), computeValue()) would compute the value twice).
			newName := b.paramNode.CloneNode()
			newValue := m.takeOrCloneCallArg(b.argNode)
			newName.Srcref(newValue) // source information should point to the argument.
			newNode := rhino.IR_VarWithValue(newName, newValue).SrcrefTreeIfMissing(newValue)
			newAliasesToAdd = append(newAliasesToAdd, newNode)
		}
	}
	m.newAliasesToAddScratch = newAliasesToAdd[:0]

	// Inline the arguments.
	result := m.functionArgumentInjector.injectByID(m.compiler, fnTemplateRoot, rhino.NilNode, paramReplacements, true)
	jutil.CheckState(result == fnTemplateRoot)

	// Now that the names have been replaced, add the new aliases for
	// the old names.
	for i := len(newAliasesToAdd) - 1; i >= 0; i-- {
		fnTemplateRoot.AddChildToFront(newAliasesToAdd[i])
	}

	return result
}

// FunctionToBlockMutator_ReplaceReturns converts returns to assignments and
// breaks, as needed. For example, with a labelName of 'foo': { return a; }
// becomes: foo: { a; break foo; } or foo: { resultName = a; break foo; }
//
// resultMustBeSet is whether the result must always be set to a value.
// Returns the node containing the transformed block, this may be different
// than the passed in node 'block'.
//
// Java: FunctionToBlockMutator.replaceReturns(Node, String, String, boolean)
func FunctionToBlockMutator_ReplaceReturns(
	block rhino.Node, resultName string, labelName string, resultMustBeSet bool) rhino.Node {
	rhino.CheckNotNull(block)
	returnCount := functionInjector_countShallowReturns(block)
	return functionToBlockMutator_replaceReturnsWithCount(
		block, resultName, labelName, resultMustBeSet, returnCount)
}

func functionToBlockMutator_replaceReturnsWithCount(
	block rhino.Node, resultName string, labelName string, resultMustBeSet bool, returnCount int) rhino.Node {
	rhino.CheckNotNull(block)

	root := block
	hasReturnAtExit := false
	if returnCount > 0 {
		hasReturnAtExit = FunctionToBlockMutator_HasReturnAtExit(block)
		// TODO(johnlenz): Simpler not to special case this,
		// and let it be optimized later.
		if hasReturnAtExit {
			FunctionToBlockMutator_ConvertLastReturnToStatement(block, resultName)
			returnCount--
		}

		if returnCount > 0 {
			// A label and breaks are needed.

			// Add the breaks
			FunctionToBlockMutator_ReplaceReturnWithBreak(block, rhino.NilNode, resultName, labelName)

			// Add label
			name := rhino.IR_LabelName(labelName).Srcref(block)
			label := rhino.IR_Label(name, block).Srcref(block)

			newRoot := rhino.IR_Block().Srcref(block)
			newRoot.AddChildToBack(label)

			// The label is now the root.
			root = newRoot
		}
	}

	// If there wasn't an return at the end of the function block, and we need
	// a result, add one to the block.
	if resultMustBeSet && !hasReturnAtExit && resultName != "" {
		FunctionToBlockMutator_AddDummyAssignment(block, resultName)
	}

	return root
}

/**********************************************************************
 *  Functions following here are general node transformation functions
 **********************************************************************/

// FunctionToBlockMutator_AddDummyAssignment, example: a = (void) 0;
//
// Java: FunctionToBlockMutator.addDummyAssignment(Node, String)
func FunctionToBlockMutator_AddDummyAssignment(node rhino.Node, resultName string) {
	jutil.CheckArgument(node.IsBlock())

	// A result is needed create a dummy value.
	srcLocation := node
	retVal := NodeUtil_NewUndefinedNode(srcLocation)
	resultNode := FunctionToBlockMutator_CreateAssignStatementNode(resultName, retVal)
	resultNode.SrcrefTreeIfMissing(node)

	node.AddChildToBack(resultNode)
}

// FunctionToBlockMutator_ConvertLastReturnToStatement replaces the 'return'
// statement with its child expression. "return foo()" becomes "foo()" or
// "resultName = foo()" "return" is removed or becomes "resultName = void 0".
//
// Java: FunctionToBlockMutator.convertLastReturnToStatement(Node, String)
func FunctionToBlockMutator_ConvertLastReturnToStatement(block rhino.Node, resultName string) {
	lastStmt, hasLast := ast.MustBlock(block).LastStatement()
	ret, isRet := lastStmt.AsReturn()
	jutil.CheckArgument(hasLast && isRet)
	retNode := ret.Node()
	resultNode := FunctionToBlockMutator_GetReplacementReturnStatement(retNode, resultName)

	if resultNode.IsNil() {
		retNode.Detach()
	} else {
		resultNode.SrcrefTreeIfMissing(retNode)
		retNode.ReplaceWith(resultNode)
	}
}

// FunctionToBlockMutator_CreateAssignStatementNode creates a valid statement
// Node containing an assignment to name of the given expression.
//
// Java: FunctionToBlockMutator.createAssignStatementNode(String, Node)
func FunctionToBlockMutator_CreateAssignStatementNode(name string, expression rhino.Node) rhino.Node {
	// Create 'name = result-expression;' statement.
	// EXPR (ASSIGN (NAME, EXPRESSION))
	nameNode := rhino.IR_Name(name)
	assign := rhino.IR_Assign(nameNode, expression)
	return NodeUtil_NewExpr(assign)
}

// FunctionToBlockMutator_GetReplacementReturnStatement replaces the 'return'
// statement with its child expression. If the result is needed (resultName
// != ""): "return foo()" becomes "resultName = foo()" "return" becomes
// "resultName = void 0". Otherwise: "return foo()" becomes "foo()" "return",
// null is returned.
//
// Java: FunctionToBlockMutator.getReplacementReturnStatement(Node, String)
func FunctionToBlockMutator_GetReplacementReturnStatement(node rhino.Node, resultName string) rhino.Node {
	var resultNode rhino.Node

	var retVal rhino.Node
	if val, hasVal := ast.MustReturn(node).Value(); hasVal {
		// Clone the child as the child hasn't been removed
		// from the node yet.
		retVal = val.Node().CloneTree()
	}

	if resultName == "" {
		if !retVal.IsNil() {
			resultNode = NodeUtil_NewExpr(retVal) // maybe null.
		}
	} else {
		if retVal.IsNil() {
			// A result is needed create a dummy value.
			srcLocation := node
			retVal = NodeUtil_NewUndefinedNode(srcLocation)
		}
		// Create a "resultName = retVal;" statement.
		resultNode = FunctionToBlockMutator_CreateAssignStatementNode(resultName, retVal)
	}

	return resultNode
}

// FunctionToBlockMutator_HasReturnAtExit returns whether the given block end
// with an return statement.
//
// Java: FunctionToBlockMutator.hasReturnAtExit(Node)
func FunctionToBlockMutator_HasReturnAtExit(block rhino.Node) bool {
	// Only inline functions that return something (empty returns
	// will be handled by ConstFolding+EmptyFunctionRemoval)
	last, hasLast := ast.MustBlock(block).LastStatement()
	if !hasLast {
		return false
	}
	_, isRet := last.AsReturn()
	return isRet
}

// FunctionToBlockMutator_ReplaceReturnWithBreak replaces the 'return'
// statement with its child expression. "return foo()" becomes "{foo();
// break;}" or "{resultName = foo(); break;}" "return" becomes {break;} or
// "{resultName = void 0;break;}".
//
// Java: FunctionToBlockMutator.replaceReturnWithBreak(Node, Node, String, String)
func FunctionToBlockMutator_ReplaceReturnWithBreak(
	current rhino.Node, parent rhino.Node, resultName string, labelName string) rhino.Node {

	if current.IsFunction() || current.IsExprResult() {
		// Don't recurse into functions definitions, and expressions can't
		// contain RETURN nodes.
		return current
	}

	if current.IsReturn() {
		jutil.CheckState(NodeUtil_IsStatementBlock(parent))

		resultNode := FunctionToBlockMutator_GetReplacementReturnStatement(current, resultName)
		breakNode := rhino.IR_BreakNodeWithLabel(rhino.IR_LabelName(labelName))

		// Replace the node in parent, and reset current to the first new child.
		breakNode.SrcrefTreeIfMissing(current)
		current.ReplaceWith(breakNode)
		if !resultNode.IsNil() {
			resultNode.SrcrefTreeIfMissing(current)
			resultNode.InsertBefore(breakNode)
		}
		current = breakNode
	} else {
		for c := current.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			// c may be replaced.
			c = FunctionToBlockMutator_ReplaceReturnWithBreak(c, current, resultName, labelName)
		}
	}

	return current
}
