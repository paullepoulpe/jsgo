// Port of com/google/javascript/jscomp/Es6RewriteBlockScopedDeclaration.java @ closure-compiler v20261015

/*
 * Copyright 2014 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Es6RewriteBlockScopedDeclaration rewrites "let"s and "const"s as "var"s.
// Renames block-scoped declarations and their references when necessary.
//
// Note that this must run after Es6RewriteDestructuring, since it does not
// process destructuring let/const declarations at all.
type Es6RewriteBlockScopedDeclaration struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler         AbstractCompiler
	astFactory       *AstFactory
	letConsts        *jutil.LinkedHashSet[*rhino.Node]
	uniqueIdSupplier *UniqueIdSupplier
}

// Java: Es6RewriteBlockScopedDeclaration.transpiledFeatures
var es6RewriteBlockScopedDeclaration_transpiledFeatures = parser.FeatureSet_BARE_MINIMUM.WithFeatures(parser.FeatureSet_Feature_LET_DECLARATIONS, parser.FeatureSet_Feature_CONST_DECLARATIONS)

// Java: Es6RewriteBlockScopedDeclaration.LOOP_PARAM_NAME_PREFIX
const es6RewriteBlockScopedDeclaration_LOOP_PARAM_NAME_PREFIX = "$jscomp$loop_param$"

// Java: Es6RewriteBlockScopedDeclaration(AbstractCompiler)
func NewEs6RewriteBlockScopedDeclaration(compiler AbstractCompiler) *Es6RewriteBlockScopedDeclaration {
	e := &Es6RewriteBlockScopedDeclaration{}
	e.letConsts = jutil.NewLinkedHashSet[*rhino.Node]()
	e.compiler = compiler
	e.uniqueIdSupplier = compiler.GetUniqueIdSupplier()
	e.astFactory = compiler.CreateAstFactory()
	return e
}

// Java: Es6RewriteBlockScopedDeclaration.visit(NodeTraversal, Node, Node)
func (e *Es6RewriteBlockScopedDeclaration) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	// NOTE: This pass depends on for-of being transpiled away before it runs.
	jutil.CheckStateMsg1(parent == nil || !parent.IsForOf(), "%s", parent)

	if n.IsLet() || n.IsConst() {
		// nothing to do for vars here
		e.letConsts.Add(n)
		for nameNode := n.GetFirstChild(); nameNode != nil; nameNode = nameNode.GetNext() {
			e.VisitBlockScopedNameDeclaration(n, nameNode)
		}
	}
}

// Java: Es6RewriteBlockScopedDeclaration.process(Node, Node)
func (e *Es6RewriteBlockScopedDeclaration) Process(externs *rhino.Node, root *rhino.Node) {
	// All declared names are already unique post normalize, so we can safely just switch 'let' or
	// 'const' to 'var' in all non-loop cases.
	// - Gather a list of let & const variables
	// - Also add `= void 0` to any that are not initialized.
	NodeTraversal_Traverse(e.compiler, root, e)
	if !e.letConsts.IsEmpty() {
		sc := NewSyntacticScopeCreator(e.compiler)
		sc.EnableArenaPool()
		defer sc.ReleaseArena()
		transformer := newEs6RewriteBlockScopedDeclaration_LoopClosureTransformer(e)
		NodeTraversal_Builder_().SetCompiler(e.compiler).SetScopeCreator(sc).SetCallback(transformer).Traverse(root)
		transformer.TransformLoopClosure()
		e.RewriteDeclsToVars()
	}
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(e.compiler, root, es6RewriteBlockScopedDeclaration_transpiledFeatures)
}

// VisitBlockScopedNameDeclaration normalizes block-scoped let/const
// declarations with no initializer in the loop body to be initialized to
// undefined.
//
// Java: Es6RewriteBlockScopedDeclaration.visitBlockScopedNameDeclaration(Node, Node)
func (e *Es6RewriteBlockScopedDeclaration) VisitBlockScopedNameDeclaration(decl *rhino.Node, nameNode *rhino.Node) {
	parent := decl.GetParent()
	// Normalize "let x;" to "let x = undefined;" if in a loop, since we later convert x
	// to be $jscomp$loop$0.x and want to reset the property to undefined every loop iteration.
	if !nameNode.HasChildren() && (parent == nil || !parent.IsForIn()) && e.InLoop(decl) {
		undefined := e.astFactory.CreateUndefinedValue().SrcrefTree(nameNode)
		nameNode.AddChildToFront(undefined)
		e.compiler.ReportChangeToEnclosingScope(undefined)
	}
}

// InLoop reports whether n is inside a loop. If n is inside a function which
// is inside a loop, we do not consider it to be inside a loop.
//
// Java: Es6RewriteBlockScopedDeclaration.inLoop(Node)
func (e *Es6RewriteBlockScopedDeclaration) InLoop(n *rhino.Node) bool {
	enclosingNode :=
		NodeUtil_GetEnclosingNode(n, Es6RewriteBlockScopedDeclaration_IsLoopOrFunction)
	return enclosingNode != nil && !enclosingNode.IsFunction()
}

// Java: Es6RewriteBlockScopedDeclaration.isLoopOrFunction(Node)
func Es6RewriteBlockScopedDeclaration_IsLoopOrFunction(n *rhino.Node) bool {
	return n.IsFunction() || NodeUtil_IsLoopStructure(n)
}

// Java: Es6RewriteBlockScopedDeclaration.extractInlineJSDoc(Node, Node, Node)
func Es6RewriteBlockScopedDeclaration_ExtractInlineJSDoc(srcDeclaration *rhino.Node, srcName *rhino.Node, destDeclaration *rhino.Node) {
	existingInfo := srcDeclaration.GetJSDocInfo()
	if existingInfo == nil {
		// Extract inline JSDoc from "src" and add it to the "dest" node.
		existingInfo = srcName.GetJSDocInfo()
		srcName.SetJSDocInfo(nil)
	}
	builder := rhino.JSDocInfo_Builder_MaybeCopyFrom(existingInfo)
	destDeclaration.SetJSDocInfo(builder.Build())
}

// Java: Es6RewriteBlockScopedDeclaration.maybeAddConstJSDoc(Node, Node, Node)
func Es6RewriteBlockScopedDeclaration_MaybeAddConstJSDoc(srcDeclaration *rhino.Node, srcName *rhino.Node, destDeclaration *rhino.Node) {
	if srcDeclaration.IsConst() {
		Es6RewriteBlockScopedDeclaration_ExtractInlineJSDoc(srcDeclaration, srcName, destDeclaration)
		builder := rhino.JSDocInfo_Builder_MaybeCopyFrom(destDeclaration.GetJSDocInfo())
		builder.RecordConstancy()
		destDeclaration.SetJSDocInfo(builder.Build())
	}
}

// HandleLetConstDeclarationList converts all declared names of the given
// declarationList of let/const declarations to their own, separate
// (normalized) var declarations.
//
// "const i = 0, j = 0;" becomes "/** @const */ var i = 0; /** @const */ var j = 0;"
//
// If the declarationList of let/const declarations is in a FOR intializer,
// moves those into separate declarations outside the FOR loop (for
// maintaining normalization).
//
// We can only have let/const declaration lists (not var) here as this pass
// runs post normalize.
//
// Java: Es6RewriteBlockScopedDeclaration.handleLetConstDeclarationList(Node, Node)
func (e *Es6RewriteBlockScopedDeclaration) HandleLetConstDeclarationList(declarationList *rhino.Node, parent *rhino.Node) {
	if declarationList.IsVar() {
		// This declaration list is already handled and we can safely return. This happens because
		// this method is also called by {@code replaceDeclarationWithProperty} where it gets
		// repeatedly called for each name in the declarationList. After the first call
		// to this (for the first name), this declarationList would no longer be {@code Token.CONST}
		// (i.e. it would've become a separate var with an {@code /** const */} annotation).
		jutil.CheckStateMsg(
			!declarationList.HasChildren() || declarationList.HasOneChild(),
			"var declaration can have maximum one child post normalization.")
		return
	}

	if parent.IsVanillaFor() {
		/*
		 * This is needed to handle the case where we get:
		 *
		 * <ol>
		 *   <li>`for (let x = ...)` or
		 *   <li>`for (const x = ...`. <
		 * </ol>
		 *
		 * <p>Normalize only moves "var" outside the for initializer, it allows let/const within the
		 * initializer. So both #1 and #2 can exist regardless if this pass runs before and after
		 * normalize.
		 */
		e.HandleLetConstDeclarationListInVanillaForInitializer(declarationList, parent)
		return
	}

	if parent.IsForIn() {
		/*
		 * This is needed to handle the case where we get:
		 *
		 * <ol>
		 *   <li>`for (let x in ...)` or
		 *   <li>`for (const x in ...`. <
		 * </ol>
		 *
		 *
		 * <p>Normalize only moves "var" outside the for initializer, it allows let/const within the
		 * initializer. So both #1 and #2 can exist regardless if this pass runs before and after
		 * normalize.
		 */
		e.HandleLetConstDeclarationListInForInInitializer(declarationList, parent)
		return
	}

	// convert all names to their own, separate (normalized) declarations
	for declarationList.HasChildren() {
		name := declarationList.GetLastChild()
		newDeclaration := rhino.IR_Var(name.Detach()).Srcref(declarationList)
		/*
		 * This method gets called multiple times by {@code replaceDeclarationWithProperty}. In the
		 * first call, it splits the declarationList into individual /\** @const *\/ var declarations.
		 * i.e. `const a,b` --> `/\** @const *\/ var a; /\** @const *\/ var b;`
		 *
		 * <p>Then it gets invoked for each individual var declaration. once for each name in the
		 * original declarationList). However, after the first call to this (for the first name), this
		 * declarationList would no longer be {@code Token.CONST} (i.e. it would've become a separate
		 * var with an {@code /*\* const *\/} annotation).
		 *
		 * <p>Previously, this method rewrote those individual var declarations again into new var
		 * declarations. Without copying over the JSDoc, those newly created var declarations did not
		 * contain `@const`, and produce this:
		 *
		 * <p>`/\** @const *\/ var a; /\** @const *\/ var b;` -->`var a; var b;`
		 *
		 * <p>Now, even though the "redundant rewriting" is handled by early-returning from this
		 * method whenever this method gets called with a non-declaration var list, it's still
		 * important that the individual var declarations created by the first rewriting contain the
		 * annotations from original declaration list. Hence we propagate JSDoc from declarationList
		 * into the individual var declarations.
		 */
		Es6RewriteBlockScopedDeclaration_ExtractInlineJSDoc(declarationList, name, newDeclaration)
		Es6RewriteBlockScopedDeclaration_MaybeAddConstJSDoc(declarationList, name, newDeclaration)
		newDeclaration.InsertAfter(declarationList)
		e.compiler.ReportChangeToEnclosingScope(parent)
	}

	// declarationList has no children left. Remove.
	declarationList.Detach()
	e.compiler.ReportChangeToEnclosingScope(parent)
}

// HandleLetConstDeclarationListInForInInitializer: we can only have
// let/const declaration lists (not var) here as this pass runs post
// normalize.
//
// Java: Es6RewriteBlockScopedDeclaration.handleLetConstDeclarationListInForInInitializer(Node, Node)
func (e *Es6RewriteBlockScopedDeclaration) HandleLetConstDeclarationListInForInInitializer(declarationList *rhino.Node, parent *rhino.Node) {
	jutil.CheckState(parent.IsForIn())
	first := declarationList
	lhs := first.GetFirstChild()
	loopNode := parent
	if lhs.IsDestructuringLhs() {
		// This pass relies on destructuring syntax being already removed. Hence we must not enter
		// this case.
		// TODO: b/279640656 Enable this code path once this pass runs unconditionally in stage3.
		jutil.CheckStateMsg(
			false, "Destructuring syntax is unsupported in ES6RewriteBlockScopedDeclarations pass")
		// Transform:
		//    `for (let [a, b = 3] in c) {}` or `for (const [a, b = 3] in c) {}`
		// to:
		// <pre>
		//    `var a; var b; for ([a, b = 3] in c) {}`
		// </pre>
		// respectively
		NodeUtil_VisitLhsNodesInNode(
			lhs,
			func(name *rhino.Node) {
				// Add a declaration outside the for loop for the given name.
				jutil.CheckStateMsg1(
					name.IsName(),
					"lhs in destructuring declaration should be a simple name. (%s)",
					name)
				newName := rhino.IR_Name(name.GetString()).Srcref(name)
				newVar := rhino.IR_Var(newName).Srcref(name)
				Es6RewriteBlockScopedDeclaration_ExtractInlineJSDoc(declarationList, name, newVar)
				// if the initializer name was a const, the newName must no longer be a const.
				e.AddNodeBeforeLoop(newVar, loopNode)
			})

		// Transform `for (var [a, b]... )` to `for ([a, b]...`
		destructuringPattern := lhs.RemoveFirstChild()
		first.ReplaceWith(destructuringPattern)
	} else {
		// Transform:
		//    for (let a in b) {}
		// to:
		//    var a; for (a in b) {};
		// and:
		//    for (const a in b) {}
		// to:
		//    var a; for (a in b) {};
		newStatement := first.CloneTree()
		name := newStatement.GetFirstChild().CloneNode()
		// cloning also copies over any properties
		if name.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME) {
			// if the initializer name was a const, it must no longer be a marked const. However,
			// marking it as non-const requires us to also change all references to the initializer name
			// in the for-in's body. That is expensive to do. Hence, we leave the initializer name and
			// its references within the loop to be const (trade off optimizability for compile time).
		}
		newStatement.SetToken(rhino.Token_VAR)
		Es6RewriteBlockScopedDeclaration_ExtractInlineJSDoc(first, first.GetFirstChild(), newStatement)
		first.ReplaceWith(name)
		e.AddNodeBeforeLoop(newStatement, loopNode)
	}
}

// HandleLetConstDeclarationListInVanillaForInitializer: we can only have
// let/const declaration lists (not var) here as this pass runs post
// normalize.
//
// Java: Es6RewriteBlockScopedDeclaration.handleLetConstDeclarationListInVanillaForInitializer(Node, Node)
func (e *Es6RewriteBlockScopedDeclaration) HandleLetConstDeclarationListInVanillaForInitializer(
	declarationList *rhino.Node, parent *rhino.Node) {
	jutil.CheckState(parent.IsVanillaFor())
	// if the declarationList is in a FOR initializer, move it outside
	insertSpot := e.GetInsertSpotBeforeLoop(parent)
	// convert all names to their own, separate (normalized) declarations
	for declarationList.HasChildren() {
		name := declarationList.GetLastChild()
		newDeclaration := rhino.IR_Var(name.Detach()).Srcref(declarationList)
		Es6RewriteBlockScopedDeclaration_ExtractInlineJSDoc(declarationList, name, newDeclaration)
		if name.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME) {
			// if the initializer name was a const, it must no longer be a marked const. However,
			// marking it as non-const requires us to also change all references to the initializer name
			// in the for's body. That is expensive to do. Hence, we leave the initializer name and
			// its references within the loop to be const (trade off optimizability for compile time).
		}
		// generate normalized var initializer (i.e. outside FOR)
		newDeclaration.InsertBefore(insertSpot)
		insertSpot = newDeclaration
		e.compiler.ReportChangeToEnclosingScope(parent)
	}
	// make FOR initializer empty `for (; cond; incr)`
	empty := e.astFactory.CreateEmpty().Srcref(declarationList)
	declarationList.ReplaceWith(empty)
	e.compiler.ReportChangeToEnclosingScope(empty)
}

// Java: Es6RewriteBlockScopedDeclaration.addNodeBeforeLoop(Node, Node)
func (e *Es6RewriteBlockScopedDeclaration) AddNodeBeforeLoop(newNode *rhino.Node, loopNode *rhino.Node) {
	insertSpot := e.GetInsertSpotBeforeLoop(loopNode)
	newNode.InsertBefore(insertSpot)
	e.compiler.ReportChangeToEnclosingScope(newNode)
}

// Java: Es6RewriteBlockScopedDeclaration.getInsertSpotBeforeLoop(Node)
func (e *Es6RewriteBlockScopedDeclaration) GetInsertSpotBeforeLoop(loopNode *rhino.Node) *rhino.Node {
	insertSpot := loopNode
	for insertSpot.GetParent().IsLabel() {
		insertSpot = insertSpot.GetParent()
	}
	return insertSpot
}

// Java: Es6RewriteBlockScopedDeclaration.rewriteDeclsToVars()
func (e *Es6RewriteBlockScopedDeclaration) RewriteDeclsToVars() {
	if !e.letConsts.IsEmpty() {
		for n := range e.letConsts.All() {
			// The Normalize pass splits simple let/const/var declarations already. But it doesn't split
			// let/const declaration list if it's in a for-loop initializer. Normalize only moves "var"
			// outside the for initializer, it allows let/const within the initializer. But when we're
			// rewriting those initializer let/const delcaration lists into simple "vars" lists here, we
			// split the var declaration list to maintain normalization.
			e.HandleLetConstDeclarationList(n, n.GetParent())
		}
	}
}

// es6RewriteBlockScopedDeclaration_LoopClosureTransformer transforms
// let/const declarations captured by loop closures.
type es6RewriteBlockScopedDeclaration_LoopClosureTransformer struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	outer *Es6RewriteBlockScopedDeclaration

	loopObjectMap *jutil.LinkedHashMap[*rhino.Node, *es6RewriteBlockScopedDeclaration_LoopObject]

	nodesRequiringLoopObjectsClosureMap *jutil.SetMultimap[*rhino.Node, *es6RewriteBlockScopedDeclaration_LoopObject]
	// Go-only: language-adapter — a Java HashMultimap that is only queried
	// (put, containsEntry), never iterated: an insertion-ordered multimap.
	nodesHandledForLoopObjectClosure *jutil.SetMultimap[*rhino.Node, string]
	referenceMap                     *jutil.SetMultimap[*Var, *rhino.Node]
}

// Java: Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.LOOP_OBJECT_NAME
const es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LOOP_OBJECT_NAME = "$jscomp$loop"

// Java: Es6RewriteBlockScopedDeclaration.LoopClosureTransformer()
func newEs6RewriteBlockScopedDeclaration_LoopClosureTransformer(outer *Es6RewriteBlockScopedDeclaration) *es6RewriteBlockScopedDeclaration_LoopClosureTransformer {
	l := &es6RewriteBlockScopedDeclaration_LoopClosureTransformer{outer: outer}
	l.loopObjectMap = jutil.NewLinkedHashMap[*rhino.Node, *es6RewriteBlockScopedDeclaration_LoopObject]()
	l.nodesRequiringLoopObjectsClosureMap = jutil.NewSetMultimap[*rhino.Node, *es6RewriteBlockScopedDeclaration_LoopObject]()
	l.nodesHandledForLoopObjectClosure = jutil.NewSetMultimap[*rhino.Node, string]()
	l.referenceMap = jutil.NewSetMultimap[*Var, *rhino.Node]()
	return l
}

// Java: Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.shouldTraverse(NodeTraversal, Node, Node)
func (l *es6RewriteBlockScopedDeclaration_LoopClosureTransformer) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if !NodeUtil_IsReferenceName(n) {
		return true
	}

	name := n.GetString()
	referencedIn := t.GetScope()
	var_ := referencedIn.GetVar(name)
	if var_ == nil {
		return true
	}

	if !var_.IsLet() && !var_.IsConst() {
		return true
	}

	// Traverse nodes up from let/const declaration:
	// If we hit a function or the root before a loop - Not a loop closure.
	// if we hit a loop first - maybe loop closure.
	declaredIn := var_.GetScope().(*Scope)
	var loopNode *rhino.Node
	for s := declaredIn; ; s = s.GetParent() {
		scopeRoot := s.GetRootNode()
		if NodeUtil_IsLoopStructure(scopeRoot) {
			loopNode = scopeRoot
			break
		} else if scopeRoot.HasParent() && NodeUtil_IsLoopStructure(scopeRoot.GetParent()) {
			loopNode = scopeRoot.GetParent()
			break
		} else if s.IsFunctionBlockScope() || s.IsGlobal() {
			return true
		}
	}

	l.referenceMap.Put(var_, n)

	// Traverse scopes from reference scope to declaration scope.
	// If we hit a function - loop closure detected.
	var outerMostFunctionScope *Scope
	for s := referencedIn; s != declaredIn && s.GetRootNode() != loopNode; s = s.GetParent() {
		if s.IsFunctionScope() {
			outerMostFunctionScope = s
		}
	}

	if outerMostFunctionScope != nil {
		enclosingFunction := outerMostFunctionScope.GetRootNode()

		// There are two categories of functions we might find here:
		//  1. a getter or setter in an object literal. We will wrap the entire object literal in
		//     a closure to capture the value of the let/const.
		//  2. a function declaration or expression. We will wrap the function in a closure.
		// (At this point, class methods/getters/setters and object literal member functions are
		// transpiled away.)
		var nodeToWrapInClosure *rhino.Node
		if enclosingFunction.GetParent().IsGetterDef() ||
			enclosingFunction.GetParent().IsSetterDef() {
			nodeToWrapInClosure = enclosingFunction.GetGrandparent()
			jutil.CheckState(nodeToWrapInClosure.IsObjectLit())
		} else {
			nodeToWrapInClosure = enclosingFunction
		}
		if l.nodesHandledForLoopObjectClosure.ContainsEntry(nodeToWrapInClosure, name) {
			return true
		}
		l.nodesHandledForLoopObjectClosure.Put(nodeToWrapInClosure, name)

		object := l.loopObjectMap.ComputeIfAbsent(
			loopNode, func(k *rhino.Node) *es6RewriteBlockScopedDeclaration_LoopObject {
				return newEs6RewriteBlockScopedDeclaration_LoopObject(l.CreateUniqueObjectName(t.GetInput()))
			})
		object.vars.Add(var_)
		l.nodesRequiringLoopObjectsClosureMap.Put(nodeToWrapInClosure, object)
	}
	return true
}

// Java: Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.createUniqueObjectName(CompilerInput)
func (l *es6RewriteBlockScopedDeclaration_LoopClosureTransformer) CreateUniqueObjectName(input *CompilerInput) string {
	return es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LOOP_OBJECT_NAME + "$" + l.outer.uniqueIdSupplier.GetUniqueId(input)
}

// Java: Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.getLoopObjPropName(Var)
func (l *es6RewriteBlockScopedDeclaration_LoopClosureTransformer) GetLoopObjPropName(var_ *Var) string {
	// NOTE: var.getName() would be wrong here, because it will still contain the original
	// and possibly non-unique name for the variable. However, the name node itself will already
	// have the new and guaranteed-globally-unique name.
	return var_.GetNameNode().GetString()
}

// Java: Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.transformLoopClosure()
func (l *es6RewriteBlockScopedDeclaration_LoopClosureTransformer) TransformLoopClosure() {
	if l.loopObjectMap.IsEmpty() {
		return
	}

	wrapperFunctions := l.CreateWrapperFunctions()

	for loopNode := range l.loopObjectMap.Keys() {
		// Introduce objects to reflect the captured scope variables.
		// Fields are initially left as undefined to avoid cases like:
		//   var $jscomp$loop$0 = {$jscomp$loop$prop$i: 0, $jscomp$loop$prop$j: $jscomp$loop$0.i}
		// They are initialized lazily by changing declarations into assignments
		// later.
		loopObject := l.loopObjectMap.Get(loopNode)
		objectLitNextIteration := l.outer.astFactory.CreateObjectLit()
		l.RenameVarsToProperties(loopObject, objectLitNextIteration, loopNode)

		updateLoopObject :=
			l.outer.astFactory.CreateAssign(l.CreateLoopObjectNameNode(loopObject), objectLitNextIteration)
		objectLit :=
			// This is the only time we generate a loop object variable for a loopNode.
			rhino.IR_VarWithValue(l.CreateLoopObjectNameNode(loopObject), l.outer.astFactory.CreateObjectLit()).
				SrcrefTree(loopNode)
		l.outer.AddNodeBeforeLoop(objectLit, loopNode)
		if loopNode.IsVanillaFor() { // For
			l.ChangeVanillaForLoopHeader(loopNode, updateLoopObject)
		} else {
			loopBody := NodeUtil_GetLoopCodeBlock(loopNode)
			loopBody.AddChildToFront(rhino.IR_ExprResult(updateLoopObject).SrcrefTreeIfMissing(loopNode))
		}
		l.outer.compiler.ReportChangeToEnclosingScope(loopNode)

		l.ChangeLoopLocalVariablesToProperties(loopNode, loopObject)
	}

	// At this point, all local variables in the loop have been changed to property accesses on
	// the "loopObject" name. For the wrapper functions that we introduced in the loop, we must
	// change the name references in their body to refer to their parameter name instead of the
	// "loopObject" name.
	// TODO(bradfordcsmith): This is inefficient. We should really choose the names first, then
	// make the changes, instead of changing the same variable references twice.
	l.UpdateNamesInWrapperFunctions(wrapperFunctions)
}

// UpdateNamesInWrapperFunctions changes the loop object references in the
// introduced wrapper functions to their parameter names.
//
// Before:
//
//	var arr = [];
//	var LOOP$0 = {};
//	var i = 0;
//	for (; i < 10; LOOP$0 = {y: LOOP$0.y}, i++) {
//	  LOOP$0.y = i;
//	  arr.push((function(LOOP$0$PARAM$1) {
//	      return function() { return LOOP$0.y; }; <---- must use param name
//	  })(LOOP$0));
//	}
//
// After:
//
//	var arr = [];
//	var LOOP$0 = {};
//	var i = 0;
//	for (; i < 10; LOOP$0 = {y: LOOP$0.y}, i++) {
//	  LOOP$0.y = i;
//	  arr.push((function(LOOP$0$PARAM$1) {
//	      return function() { return LOOP$0$PARAM$1.y; }; <--- changed
//	  })(LOOP$0));
//	}
//
// Java: Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.updateNamesInWrapperFunctions(Set<Node>)
func (l *es6RewriteBlockScopedDeclaration_LoopClosureTransformer) UpdateNamesInWrapperFunctions(wrapperFunctions *jutil.LinkedHashSet[*rhino.Node]) {
	for func_ := range wrapperFunctions.All() {
		// get the param names here (and the loopObject names from it)
		paramList := func_.GetSecondChild()

		for param := paramList.GetFirstChild(); param != nil; param = param.GetNext() {
			loopObjectName :=
				jutil.StringSubstring(param.GetString(), 0, jutil.StringIndexOf(param.GetString(), es6RewriteBlockScopedDeclaration_LOOP_PARAM_NAME_PREFIX))
			l.UpdateNames( /* block */ func_.GetLastChild(), param, loopObjectName)
		}
	}
}

// CreateWrapperFunctions creates wrapper functions and calls them.
//
// Java: Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.createWrapperFunctions()
func (l *es6RewriteBlockScopedDeclaration_LoopClosureTransformer) CreateWrapperFunctions() *jutil.LinkedHashSet[*rhino.Node] {
	wrapperFunctions := jutil.NewLinkedHashSet[*rhino.Node]()
	for functionOrObjectLit := range l.nodesRequiringLoopObjectsClosureMap.KeySet() {
		// Wrapping an object literal in an IIFE would otherwise rebind `this` and `arguments` in
		// object properties/methods to the IIFE. Use ThisAndArgumentsReferenceUpdater to capture
		// and rewrite them to outer aliases before wrapping in the IIFE.
		if functionOrObjectLit.IsObjectLit() {
			enclosingFunction := NodeUtil_GetEnclosingFunction(functionOrObjectLit)
			var scopeBody *rhino.Node
			if enclosingFunction != nil {
				scopeBody = enclosingFunction.GetLastChild()
			} else {
				scopeBody = NodeUtil_GetEnclosingScript(functionOrObjectLit)
			}
			isConstructor :=
				enclosingFunction != nil && NodeUtil_IsEs6Constructor(enclosingFunction)
			script := NodeUtil_GetEnclosingScript(functionOrObjectLit)
			input := l.outer.compiler.GetInput(script.GetInputId())
			thisContext :=
				NewThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext(
					scopeBody, isConstructor, l.outer.uniqueIdSupplier.GetUniqueId(input))
			updater :=
				NewThisAndArgumentsReferenceUpdater(l.outer.compiler, thisContext, l.outer.astFactory)
			NodeTraversal_Traverse(l.outer.compiler, functionOrObjectLit, updater)
			thisContext.AddVarDeclarations(l.outer.compiler, l.outer.astFactory, script)
			for child := scopeBody.GetFirstChild(); child != nil; child = child.GetNext() {
				if child.IsLet() || child.IsConst() {
					l.outer.letConsts.Add(child)
				}
			}
		}
		returnNode := rhino.IR_ReturnNode()
		objects := l.nodesRequiringLoopObjectsClosureMap.Get(functionOrObjectLit)
		parameterNames := make([]*rhino.Node, objects.Size())
		objectNamesForCall := make([]*rhino.Node, objects.Size())
		loopObjectIndex := 0
		for object := range objects.All() {
			// This name for parameter should be unique to preserve normalization
			parameterNames[loopObjectIndex] =
				l.CreateUniqueParameterNameToUseWithinLoop(functionOrObjectLit, object)
			// this name must be the same as the loop object name
			objectNamesForCall[loopObjectIndex] = l.CreateLoopObjectNameNode(object)
			loopObjectIndex++
		}

		iife :=
			l.outer.astFactory.CreateFunction(
				"",
				rhino.IR_ParamList(parameterNames...),
				rhino.IR_BlockOfNode(returnNode),
				AstFactory_Type_OfColor(colors.StandardColors_TOP_OBJECT))
		l.outer.compiler.ReportChangeToChangeScope(iife)
		call := l.outer.astFactory.CreateCall(iife, AstFactory_Type_(functionOrObjectLit), objectNamesForCall...)
		call.PutBooleanProp(rhino.Node_FREE_CALL, true)
		var replacement *rhino.Node
		// Es6RewriteBlockScopedFunctionDeclaration must run before this pass.
		jutil.CheckStateMsg1(
			!NodeUtil_IsFunctionDeclaration(functionOrObjectLit),
			"block-scoped function declarations should not exist now: %s",
			functionOrObjectLit)
		replacement = call.SrcrefTreeIfMissing(functionOrObjectLit)
		functionOrObjectLit.ReplaceWith(replacement)
		wrapperFunctions.Add(iife)
		returnNode.AddChildToFront(functionOrObjectLit)
		l.outer.compiler.ReportChangeToEnclosingScope(replacement)
	}
	return wrapperFunctions
}

// UpdateNames renames all references within the functionOrObjectLit body
// from the old loopObjectName (e.g. "$jscomp$loop$1") to the given
// loopObjectParamName (e.g. $jscomp$loop$1$jscomp$loop_param$m123..456).
//
// Java: Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.updateNames(Node, Node, String)
func (l *es6RewriteBlockScopedDeclaration_LoopClosureTransformer) UpdateNames(node *rhino.Node, paramName *rhino.Node, loopObjectName string) {
	if node.IsName() && node.GetString() != "" /* not an anonymous function name */ {
		if node.GetString() == loopObjectName {
			node.SetString(paramName.GetString())
			return
		}
	}
	for child := node.GetFirstChild(); child != nil; child = child.GetNext() {
		l.UpdateNames(child, paramName, loopObjectName)
	}
}

// ChangeVanillaForLoopHeader pulls out the initializer and places it prior
// to the loop. The increment is updated with the new loop object and
// property assignments
//
// Java: Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.changeVanillaForLoopHeader(Node, Node)
func (l *es6RewriteBlockScopedDeclaration_LoopClosureTransformer) ChangeVanillaForLoopHeader(loopNode *rhino.Node, updateLoopObject *rhino.Node) {
	initializer := loopNode.GetFirstChild()
	initializer.ReplaceWith(rhino.IR_Empty())
	if !initializer.IsEmpty() {
		if !NodeUtil_IsNameDeclaration(initializer) {
			initializer = rhino.IR_ExprResult(initializer).Srcref(initializer)
		}
		l.outer.AddNodeBeforeLoop(initializer, loopNode)
	}

	increment := loopNode.GetChildAtIndex(2)
	if increment.IsEmpty() {
		increment.ReplaceWith(updateLoopObject.SrcrefTreeIfMissing(loopNode))
	} else {
		placeHolder := rhino.IR_Empty()
		increment.ReplaceWith(placeHolder)
		placeHolder.ReplaceWith(
			l.outer.astFactory.CreateComma(updateLoopObject, increment).SrcrefTreeIfMissing(loopNode))
	}
}

// ChangeLoopLocalVariablesToProperties changes, for captured variables,
// declarations to assignments on the corresponding field of the introduced
// object. Renames all references accordingly.
//
// Java: Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.changeLoopLocalVariablesToProperties(Node, LoopObject)
func (l *es6RewriteBlockScopedDeclaration_LoopClosureTransformer) ChangeLoopLocalVariablesToProperties(loopNode *rhino.Node, loopObject *es6RewriteBlockScopedDeclaration_LoopObject) {
	for var_ := range loopObject.vars.All() {
		newPropertyName := l.GetLoopObjPropName(var_)
		for reference := range l.referenceMap.Get(var_).All() {
			// for-of loops are transpiled away before this pass runs
			jutil.CheckStateMsg1(!loopNode.IsForOf(), "%s", loopNode)
			// For-of and for-in declarations are not altered, since they are
			// used as temporary variables for assignment.
			if loopNode.IsForIn() && loopNode.GetFirstChild() == reference.GetParent() {
				l.AssignLoopVarToLoopObjectProperty(
					loopNode, loopObject, var_, newPropertyName, reference)
			} else {
				if NodeUtil_IsNameDeclaration(reference.GetParent()) {
					l.ReplaceDeclarationWithProperty(loopObject, newPropertyName, reference)
				} else {
					l.ReplaceReferenceWithProperty(loopObject, newPropertyName, reference)
				}
			}
		}
	}
}

// Java: Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.replaceDeclarationWithProperty(LoopObject, String, Node)
func (l *es6RewriteBlockScopedDeclaration_LoopClosureTransformer) ReplaceDeclarationWithProperty(
	loopObject *es6RewriteBlockScopedDeclaration_LoopObject, newPropertyName string, reference *rhino.Node) {
	declaration := reference.GetParent()
	grandParent := declaration.GetParent()
	// Record that the let / const declaration statement will get turned into one or more
	// var statements by handleDeclarationList(), so we won't try to change it again later.
	l.outer.letConsts.Remove(declaration)
	// If the declaration contains multiple declared variables, split it apart.
	// NOTE: This call could be made for each declarationList once, rather than each name in that
	// list
	l.outer.HandleLetConstDeclarationList(declaration, grandParent)

	// The variable we're working with may have been moved to a new var statement.
	declaration = reference.GetParent()
	if reference.HasChildren() {
		// Change declaration to assignment
		newReference :=
			l.CreateLoopVarReferenceReplacement(loopObject, reference, newPropertyName)
		assign := l.outer.astFactory.CreateAssign(newReference, reference.RemoveFirstChild())
		Es6RewriteBlockScopedDeclaration_ExtractInlineJSDoc(declaration, reference, declaration)
		Es6RewriteBlockScopedDeclaration_MaybeAddConstJSDoc(declaration, reference, declaration)
		assign.SetJSDocInfo(declaration.GetJSDocInfo())

		replacement := rhino.IR_ExprResult(assign).SrcrefTreeIfMissing(declaration)
		declaration.ReplaceWith(replacement)
	} else {
		// No value is assigned, so just drop the let/const statement entirely
		declaration.Detach()
	}
	l.outer.compiler.ReportChangeToEnclosingScope(grandParent)
}

// Java: Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.replaceReferenceWithProperty(LoopObject, String, Node)
func (l *es6RewriteBlockScopedDeclaration_LoopClosureTransformer) ReplaceReferenceWithProperty(
	loopObject *es6RewriteBlockScopedDeclaration_LoopObject, newPropertyName string, reference *rhino.Node) {
	referenceParent := reference.GetParent()
	reference.ReplaceWith(
		l.CreateLoopVarReferenceReplacement(loopObject, reference, newPropertyName))
	l.outer.compiler.ReportChangeToEnclosingScope(referenceParent)
}

// AssignLoopVarToLoopObjectProperty transforms
// `for (const p in obj) { ... }`
//
// into `for (const p in obj) { $jscomp$loop$0.$jscomp$loop$prop$0$p = p; ... }`
//
// Java: Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.assignLoopVarToLoopObjectProperty(Node, LoopObject, Var, String, Node)
func (l *es6RewriteBlockScopedDeclaration_LoopClosureTransformer) AssignLoopVarToLoopObjectProperty(
	loopNode *rhino.Node, loopObject *es6RewriteBlockScopedDeclaration_LoopObject, var_ *Var, newPropertyName string, reference *rhino.Node) {
	// reference is the node loopVar in a for-in that looks like this:
	// `for (const loopVar in list) {`
	jutil.CheckStateMsg1(reference == var_.GetNameNode(), "%s", reference)
	referenceParent := reference.GetParent()
	jutil.CheckStateMsg1(NodeUtil_IsNameDeclaration(referenceParent), "%s", referenceParent)
	jutil.CheckStateMsg1(reference.IsName(), "%s", reference)
	// Start transpiled form of
	// `for (const p in obj) { ... }`
	// with this statement to copy the loop variable into the corresponding loop object
	// property.
	// `$jscomp$loop$0.$jscomp$loop$prop$0$p = p;`
	loopVarReference := reference.CloneNode()
	// `$jscomp$loop$0.$jscomp$loop$prop$0$p = p;`
	forInPropAssignmentStatemnt :=
		rhino.IR_ExprResult(
			l.outer.astFactory.CreateAssign(
				l.CreateLoopVarReferenceReplacement(loopObject, reference, newPropertyName),
				loopVarReference)).
			SrcrefTreeIfMissing(reference)
	// The first statement in the body should be creating a new loop object value
	// $jscomp$loop$0 = {
	//    $jscomp$loop$prop$0$p: $jscomp$loop$0.$jscomp$loop$prop$0$p,
	//    $jscomp$loop$prop$0$otherVar: $jscomp$loop$0.$jscomp$loop$prop$0$p,
	//    // other property update assignments
	// }
	// We need to update the loop variable's value to it immediately after that
	loopUpdateStatement :=
		loopNode.
			GetLastChild(). // loop body
			GetFirstChild() // first statement

	forInPropAssignmentStatemnt.InsertAfter(loopUpdateStatement)
}

// RenameVarsToProperties renames all variables in the loop object to
// properties
//
// Java: Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.renameVarsToProperties(LoopObject, Node, Node)
func (l *es6RewriteBlockScopedDeclaration_LoopClosureTransformer) RenameVarsToProperties(
	loopObject *es6RewriteBlockScopedDeclaration_LoopObject, objectLitNextIteration *rhino.Node, scopeRoot *rhino.Node) {
	for var_ := range loopObject.vars.All() {
		newPropertyName := l.GetLoopObjPropName(var_)
		var newPropertyValue *rhino.Node
		if var_.GetScopeRoot() == scopeRoot {
			newPropertyValue = l.CreateLoopVarReferenceReplacement(loopObject, var_.GetNameNode(), newPropertyName)
		} else {
			newPropertyValue = l.outer.astFactory.CreateUndefinedValue()
		}
		objectLitNextIteration.AddChildToBack(
			l.outer.astFactory.CreateStringKey(newPropertyName, newPropertyValue))
	}
}

// CreateLoopVarReferenceReplacement creates a
// `$jscomp$loop$0.$jscomp$loop$prop$varName$1` replacement for a reference
// to `varName`.
//
// Java: Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.createLoopVarReferenceReplacement(LoopObject, Node, String)
func (l *es6RewriteBlockScopedDeclaration_LoopClosureTransformer) CreateLoopVarReferenceReplacement(
	loopObject *es6RewriteBlockScopedDeclaration_LoopObject, reference *rhino.Node, propertyName string) *rhino.Node {
	replacement :=
		l.outer.astFactory.CreateGetProp(
			l.CreateLoopObjectNameNode(loopObject), propertyName, AstFactory_Type_(reference))
	replacement.SrcrefTree(reference)
	return replacement
}

// CreateUniqueParameterNameToUseWithinLoop gets a name other than the loop
// object name to use for parameter names of the new wrapper functions being
// added in the loop.
//
// For a loop object "$jscomp$loop$0", this creates a unique parameter name
// like "$jscomp$loop$0$jscomp$loop_param$1".
//
// Java: Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.createUniqueParameterNameToUseWithinLoop(Node, LoopObject)
func (l *es6RewriteBlockScopedDeclaration_LoopClosureTransformer) CreateUniqueParameterNameToUseWithinLoop(node *rhino.Node, loopObject *es6RewriteBlockScopedDeclaration_LoopObject) *rhino.Node {
	return l.outer.astFactory.CreateName(
		loopObject.name+
			es6RewriteBlockScopedDeclaration_LOOP_PARAM_NAME_PREFIX+
			l.outer.uniqueIdSupplier.GetUniqueId(l.outer.compiler.GetInput(NodeUtil_GetInputId(node))),
		AstFactory_Type_OfColor(colors.StandardColors_TOP_OBJECT))
}

// Java: Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.createLoopObjectNameNode(LoopObject)
func (l *es6RewriteBlockScopedDeclaration_LoopClosureTransformer) CreateLoopObjectNameNode(loopObject *es6RewriteBlockScopedDeclaration_LoopObject) *rhino.Node {
	return l.outer.astFactory.CreateName(loopObject.name, AstFactory_Type_OfColor(colors.StandardColors_TOP_OBJECT))
}

// Java: private static class Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.LoopObject
type es6RewriteBlockScopedDeclaration_LoopObject struct {
	// The name of the variable having the loop's internal variables as
	// properties, and the label applied to the block containing the
	// original loop body in cases where these are needed.
	name string

	vars *jutil.LinkedHashSet[*Var]
}

// Java: Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.LoopObject(String)
func newEs6RewriteBlockScopedDeclaration_LoopObject(name string) *es6RewriteBlockScopedDeclaration_LoopObject {
	o := &es6RewriteBlockScopedDeclaration_LoopObject{}
	o.vars = jutil.NewLinkedHashSet[*Var]()
	o.name = name
	return o
}
