// Port of com/google/javascript/jscomp/InlineAndCollapseProperties.java @ closure-compiler v20261028

/*
 * Copyright 2021 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// aggressiveInlineAliases inlines type aliases if they are explicitly or
// effectively const. Also inlines inherited static property accesses for ES6
// classes.
//
// This frees subsequent optimization passes from the responsibility of
// having to reason about alias chains and is a requirement for correct
// behavior in at least CollapseProperties and J2clPropertyInlinerPass.
//
// This is designed to be no more unsafe than CollapseProperties. It will in
// some cases inline properties, possibly past places that change the
// property value. However, it will only do so in cases where
// CollapseProperties would unsafely collapse the property anyway.
type aggressiveInlineAliases struct {
	outer *InlineAndCollapseProperties

	scopeCreator    *SyntacticScopeCreator
	scopeCollectors map[rhino.Node]*ReferenceCollector
	cachedVarNames  map[uint32]struct{}
	newNodes        astChangeSet

	codeChanged bool
}

// Java: InlineAndCollapseProperties.AggressiveInlineAliases()
func NewInlineAndCollapseProperties_AggressiveInlineAliases(outer *InlineAndCollapseProperties) *aggressiveInlineAliases {
	return &aggressiveInlineAliases{outer: outer, codeChanged: true}
}

// Java: InlineAndCollapseProperties.AggressiveInlineAliases.getLastUsedGlobalNamespace()
func (a *aggressiveInlineAliases) GetLastUsedGlobalNamespace() *GlobalNamespace {
	return a.outer.namespace
}

func (a *aggressiveInlineAliases) getScopeCreator() *SyntacticScopeCreator {
	if a.scopeCreator == nil {
		a.scopeCreator = NewSyntacticScopeCreator(a.outer.compiler)
	}
	return a.scopeCreator
}

func (a *aggressiveInlineAliases) invalidateScopeCollectors() {
	if len(a.scopeCollectors) > 0 {
		clear(a.scopeCollectors)
		clear(a.cachedVarNames)
	}
}

func leftmostQNameID(n rhino.Node) (uint32, bool) {
	if n.IsNil() {
		return 0, false
	}
	for n.IsGetProp() {
		n = n.GetFirstChild()
	}
	if n.IsName() {
		return n.GetStringID(), true
	}
	if n.IsThis() || n.IsSuper() || n.IsImportMeta() {
		return 0, true
	}
	return 0, false
}

func (a *aggressiveInlineAliases) noteMutatedQName(n rhino.Node) {
	if len(a.scopeCollectors) == 0 || n.IsNil() || n.IsNull() {
		return
	}
	id, ok := leftmostQNameID(n)
	if !ok {
		a.invalidateScopeCollectors()
		return
	}
	if id != 0 {
		if _, cached := a.cachedVarNames[id]; cached {
			a.invalidateScopeCollectors()
		}
	}
}

func (a *aggressiveInlineAliases) getReferencesInScope(
	aliasScope *Scope, aliasVar *Var,
) (*ReferenceCollector, *ReferenceCollection) {
	scopeRoot := aliasScope.GetRootNode()
	if c := a.scopeCollectors[scopeRoot]; c != nil {
		if refs := c.GetReferences(aliasVar); refs != nil {
			return c, refs
		}
	}
	collector := NewReferenceCollectorWithVarFilter(
		a.outer.compiler,
		ReferenceCollector_DO_NOTHING_BEHAVIOR,
		a.getScopeCreator(),
		func(v *Var) bool { return v.scope.GetRootNode() == scopeRoot },
	)
	collector.ProcessScope(aliasScope)
	if a.scopeCollectors == nil {
		a.scopeCollectors = make(map[rhino.Node]*ReferenceCollector, 8)
		a.cachedVarNames = make(map[uint32]struct{}, len(aliasScope.varList)+2)
	}
	a.scopeCollectors[scopeRoot] = collector
	for _, v := range aliasScope.varList {
		if v != nil {
			a.cachedVarNames[v.nameID] = struct{}{}
		}
	}
	if aliasScope.hasImplicitVars {
		a.cachedVarNames[implicitArgumentsID] = struct{}{}
	}
	return collector, collector.GetReferences(aliasVar)
}

// Java: InlineAndCollapseProperties.AggressiveInlineAliases.process(Node, Node)
func (a *aggressiveInlineAliases) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(a.outer.compiler, root, &rewriteSimpleDestructuringAliases{})

	// Building the `GlobalNamespace` dominates the cost of this pass, so it is built once and
	// updated as changes are made so it can be reused for the next iteration.
	a.outer.namespace = NewGlobalNamespaceWithLog(a.outer.decisionsLog, a.outer.compiler, root)
	for a.codeChanged {
		a.codeChanged = false
		a.inlineAliases(a.outer.namespace)
	}
	a.scopeCreator = nil
	a.scopeCollectors = nil
	a.cachedVarNames = nil
	a.newNodes = astChangeSet{}
}

// inlineAliases: for each qualified name N in the global scope, we check if:
// (a) No ancestor of N is ever aliased or assigned an unknown value type.
// (If N = "a.b.c", "a" and "a.b" are never aliased). (b) N has exactly one
// write, and it lives in the global scope. (c) N is aliased in a local
// scope. (d) N is aliased in global scope
//
// If (a) is true, then GlobalNamespace must know all the writes to N. If (a)
// and (b) are true, then N cannot change during the execution of a local
// scope. If (a) and (b) and (c) are true, then the alias can be inlined if
// the alias obeys the usual rules for how we decide whether a variable is
// inlineable. If (a) and (b) and (d) are true, then inline the alias if
// possible (if it is assigned exactly once unconditionally).
//
// For (a), (b), and (c) are true and the alias is of a constructor, we may
// also partially inline the alias - i.e. replace some references with the
// constructor but not all - since constructor properties are always
// collapsed, so we want to be more aggressive about removing aliases. This
// is similar to what FlowSensitiveInlineVariables does.
//
// If (a) is not true, but the property is a 'declared type' (which
// CollapseProperties will unsafely collapse), we also inline any properties
// without @nocollapse. This is unsafe but no more unsafe than what
// CollapseProperties does. This pass and CollapseProperties share the logic
// to determine when a name is unsafely collapsible in Name#canCollapse()
//
// See InlineVariables.
//
// Java: InlineAndCollapseProperties.AggressiveInlineAliases.inlineAliases(GlobalNamespace)
func (a *aggressiveInlineAliases) inlineAliases(namespace *GlobalNamespace) {
	// Invariant: All the names in the worklist meet condition (a).
	// adds all top-level names to the worklist, but not any properties on those names.
	//
	// Go-only: language-adapter — the ArrayDeque is a slice; pop() takes the
	// first element and add() appends.
	workList := namespace.NameForestSlice()

	for i := 0; i < len(workList); i++ {
		name := workList[i]

		// Don't attempt to inline a getter or setter property as a variable.
		if name.IsGetOrSetDefinition() {
			continue
		}

		if !name.InExterns() && // not an externs definition
			name.GetGlobalSets() == 1 && // set exactly once and only set in the global scope
			name.GetLocalSets() == 0 {
			// `name` meets condition (b). Find all of its aliases
			// and try to inline them.
			a.maybeInlineInnerName(name)
			if name.GetAliasingGets() > 0 || name.GetSubclassingGets() > 0 {
				// condition (c) and/or condition (d) are true
				a.inlineAliasesForName(name, namespace)
			}
		}
		InlineAndCollapseProperties_MaybeAddPropertiesToWorklist(name, &workList)
	}
}

// inlineAliasesForName inlines a global name into all the places where
// references for aliases to it currently exist.
//
// e.g.
//
//	const globalName = { method() {} };
//	const aliasForGlobalName = globalName;
//	aliasForGlobalName.method(); // replace this with globalName.method();
//
// This method only handles aliases created by assignment. In particular, it
// doesn't handle aliases created by inner names on class or function
// expressions. See (maybeInlineInnerName() for that).
//
// name is the global name whose aliases will be replaced; namespace is used
// to find references to the global name that create aliases e.g. `const
// aliasName = globalName`.
//
// Java: InlineAndCollapseProperties.AggressiveInlineAliases.inlineAliasesForName(Name, GlobalNamespace)
func (a *aggressiveInlineAliases) inlineAliasesForName(name *GlobalNamespace_Name, namespace *GlobalNamespace) {
	var refsBuf [8]*GlobalNamespace_Ref
	refs := name.snapshotRefs(&refsBuf)
	for _, ref := range refs {
		if ref.IsAliasingGet() && !InlineAndCollapseProperties_MayBeGlobalAlias(ref) && !ref.IsTwin() {
			// `name` meets condition (c). Try to inline it.
			// TODO(johnlenz): consider picking up new aliases at the end
			// of the pass instead of immediately like we do for global
			// inlines.
			a.inlineAliasIfPossible(name, ref, namespace)
		} else if ref.IsAliasingGet() &&
			ref.scope.GetClosestHoistScope().IsGlobal() &&
			!ref.IsTwin() {
			a.inlineGlobalAliasIfPossible(name, ref, namespace)
		} else if name.IsClass() && ref.IsSubclassingGet() && name.HasProps() {
			for prop := range name.Props() {
				a.rewriteAllSubclassInheritedAccesses(name, ref, prop, namespace)
			}
		}
	}
}

// maybeInlineInnerName: if the global name is a function with an
// inner-scope name, inline references to that name with the global name.
// Class inner name rewrites are handled in Es6NormalizeClasses.
//
// e.g.
//
//	var globalFunction = function innerName() {
//	  // change this to globalFunction.someProp
//	  use(innerName.someProp);
//	}
//
// Java: InlineAndCollapseProperties.AggressiveInlineAliases.maybeInlineInnerName(Name)
func (a *aggressiveInlineAliases) maybeInlineInnerName(globalName *GlobalNamespace_Name) {
	globalNameDeclaration := jutil.CheckNotNullMsg1(globalName.getDeclarationRef(), "%s", globalName)
	globalDeclarationNode := rhino.CheckNotNullMsg1(globalNameDeclaration.GetNode(), "%s", globalNameDeclaration)
	valueNode := NodeUtil_GetRValueOfLValue(globalDeclarationNode)
	if valueNode.IsNil() {
		// no function or class expression, so no inner name
		return
	}
	innerNameNode := InlineAndCollapseProperties_MaybeGetInnerNameNode(valueNode)
	if innerNameNode.IsNil() {
		// no inner name to require inlining
		return
	}
	innerScope := a.getScopeCreator().CreateScope(valueNode, globalNameDeclaration.scope).(*Scope)
	innerNameVar := jutil.CheckNotNull(innerScope.GetVarByNode(innerNameNode))
	collector := NewReferenceCollectorWithVarFilter(
		a.outer.compiler,
		ReferenceCollector_DO_NOTHING_BEHAVIOR,
		a.getScopeCreator(),
		func(v *Var) bool { return jutil.ObjectsEquals(innerNameVar, v) })
	collector.ProcessScope(innerScope)
	innerNameRefs := collector.GetReferences(innerNameVar)

	newNodes := &a.newNodes
	newNodes.reset()

	for innerNameRef := range innerNameRefs.All() {
		// replace all references to the inner name other than its declaration
		innerNameRefNode := innerNameRef.GetNode()
		if NodeUtil_IsNormalGet(innerNameRefNode.GetParent()) {
			// Replace `innerName` with `globalName` for `innerName.prop` and `innerName[expr]`
			//
			// TODO(b/148237949): We are intentionally ignoring cases where the inner name
			// escapes to other scopes where properties may be accessed on it (e.g. `use(InnerName)`).
			// This is unsafe, but currently necessary to avoid large code size regressions.
			//
			// NOTE: We also don't want to introduce a global reference for cases like
			// `x instanceof innerName`. It would be safe to inline these, but it also isn't
			// necessary,
			// and the introduction of a reference to a global in a local scope can cause other
			// optimizations to back off.
			a.noteMutatedQName(innerNameRefNode)
			newNodes.add(a.replaceAliasReference(globalNameDeclaration, innerNameRef))
		}
	}
	a.outer.namespace.ScanNewNodes(newNodes)
}

// rewriteAllSubclassInheritedAccesses inlines all references to inherited
// static superclass properties from the subclass or any descendant of the
// given subclass. Avoids inlining references to inherited methods when
// possible, since they may use this or super().
//
// superclassNameObj is the Name of the superclass; superclassRef the
// SUBCLASSING_REF; prop the property on the superclass to rewrite, if any
// descendant accesses it; namespace the GlobalNamespace containing
// superclassNameObj.
//
// Java: InlineAndCollapseProperties.AggressiveInlineAliases.rewriteAllSubclassInheritedAccesses(Name, Ref, Name, GlobalNamespace)
func (a *aggressiveInlineAliases) rewriteAllSubclassInheritedAccesses(
	superclassNameObj *GlobalNamespace_Name, superclassRef *GlobalNamespace_Ref, prop *GlobalNamespace_Name, namespace *GlobalNamespace) bool {
	if !prop.CanCollapse() {
		return false // inlining is a) unnecessary if there is @nocollapse and b) might break
		// usages of `this` in the method
	}
	subclass := InlineAndCollapseProperties_GetSubclassForEs6Superclass(superclassRef.GetNode())
	if subclass.IsNil() {
		return false
	}
	subclassNameObj, ok := namespace.lookupGlobalNode(subclass)
	if !ok {
		return false
	}
	subclassPropNameObj := namespace.lookupGlobalChildByID(subclassNameObj, prop.baseNameID)
	// Don't rewrite if the subclass ever shadows the parent static property.
	// This may also back off on cases where the subclass first accesses the parent property, then
	// shadows it.
	if subclassPropNameObj != nil &&
		(subclassPropNameObj.GetLocalSets() > 0 || subclassPropNameObj.GetGlobalSets() > 0) {
		return false
	}

	// Recurse to find potential sub-subclass accesses of the superclass property.
	if subclassNameObj != nil && subclassNameObj.SubclassingGetCount() > 0 {
		var subRefsBuf [8]*GlobalNamespace_Ref
		for _, ref := range subclassNameObj.snapshotRefs(&subRefsBuf) {
			if ref.IsSubclassingGet() {
				a.rewriteAllSubclassInheritedAccesses(superclassNameObj, ref, prop, namespace)
			}
		}
	}

	if subclassPropNameObj != nil {
		// Use this node as a template for rewriteNestedAliasReference.
		superclassNameNode := superclassNameObj.getDeclarationRef().GetNode()
		if superclassNameNode.IsName() {
			superclassNameNode = superclassNameNode.CloneNode()
		} else if superclassNameNode.IsGetProp() {
			superclassNameNode = superclassNameNode.CloneTree()
		} else {
			return false
		}

		newNodes := &a.newNodes
		newNodes.reset()
		a.rewriteNestedAliasReference(superclassNameNode, 0, newNodes, subclassPropNameObj)
		namespace.ScanNewNodes(newNodes)
	}
	return true
}

// isInlineableModuleExportsAlias recognizes aliases for the special global
// variables representing the `exports` values for goog modules which are
// safe to inline.
//
// The compiler will enforce that references to these module objects only
// occur after they have been initialized and that they are only ever
// assigned one time. So, as long as a variable that aliases one of them is
// only ever assigned that one module exports value and only used for
// getting properties off of it with destructuring or dot-property access,
// we know it is safe to inline those aliases.
//
// This inlining is important to support `await goog.requireDynamic()` when
// transpiling down to es5. In that case the compiler generates one of these
// aliases which risks producing the PARTIAL_NAMESPACE_WARNING.
//
// Java: InlineAndCollapseProperties.AggressiveInlineAliases.isInlineableModuleExportsAlias(Name, ReferenceCollection)
func (a *aggressiveInlineAliases) isInlineableModuleExportsAlias(name *GlobalNamespace_Name, aliasRefs *ReferenceCollection) bool {
	if !strings.HasPrefix(name.RootBaseName(), ClosureRewriteModule_MODULE_EXPORTS_PREFIX) {
		return false
	}
	aliasedNameStr := name.GetFullName()
	// the first rhs value Node we find assigned to the name
	var firstRhs rhino.Node
	size := len(aliasRefs.references)
	for i := 0; i < size; i++ {
		ref := aliasRefs.references[i]
		if ref.IsVarDeclaration() || ref.GetParent().IsAssign() {
			var rhs rhino.Node
			if assign, ok := ast.AsAssign(ref.GetParent()); ok {
				rhs = assign.Value().Node()
			} else if nameNode, ok := ast.AsName(ref.GetNode()); ok {
				if initExpr, ok := nameNode.Init(); ok {
					rhs = initExpr.Node()
				}
			}

			if !rhs.IsNil() {
				// Make sure that every time a value is assigned to the alias, it is
				// the aliased name we expect. If not, we cannot inline this alias.
				if firstRhs.IsNil() {
					// NOTE: comparing with a string is slower, but necessary the first time.
					if !rhs.MatchesQualifiedName(aliasedNameStr) {
						return false
					}
					firstRhs = rhs
				} else {
					// comparing nodes is faster than comparing with the string name
					if !rhs.MatchesQualifiedNameNode(firstRhs) {
						return false
					}
				}
			}
			continue
		}
		if !ref.IsDotPropertyAccess() && !ref.IsAssignedToObjectDestructuringPattern() {
			return false
		}
	}
	return true
}

// inlineAliasIfPossible attempts to inline a non-global alias of a global
// name.
//
// It is assumed that the name for which it is an alias meets conditions (a)
// and (b).
//
// The non-global alias is only inlinable if it is well-defined and assigned
// once, according to the definitions in ReferenceCollection
//
// If the aliasing name is completely removed, also deletes the aliasing Ref.
//
// name is the global name being aliased; alias the aliasing reference to the
// name to remove.
//
// Java: InlineAndCollapseProperties.AggressiveInlineAliases.inlineAliasIfPossible(Name, Ref, GlobalNamespace)
func (a *aggressiveInlineAliases) inlineAliasIfPossible(name *GlobalNamespace_Name, alias *GlobalNamespace_Ref, namespace *GlobalNamespace) {
	// Ensure that the alias is assigned to a local variable at that
	// variable's declaration. If the alias's parent is a NAME,
	// then the NAME must be the child of a VAR, LET, or CONST node, and we must
	// be in a VAR, LET, or CONST assignment.
	// Otherwise if the parent is an assign, we are in a "a = alias" case.
	aliasParent := alias.GetNode().GetParent()
	if lhsNode, ok := aliasLhsNode(aliasParent); ok {
		aliasVar := alias.scope.GetVarByNode(lhsNode)
		jutil.CheckStateMsg1(aliasVar != nil, "Expected variable to be defined in scope (%s)", lhsNode.GetString())
		aliasScope := aliasVar.scope
		collector, aliasRefs := a.getReferencesInScope(aliasScope, aliasVar)

		if aliasRefs.IsWellDefined() &&
			(aliasRefs.IsAssignedOnceInLifetime() ||
				a.isInlineableModuleExportsAlias(name, aliasRefs)) {
			// The alias is well-formed, so do the inlining now.
			newNodes := &a.newNodes
			newNodes.reset()
			size := len(aliasRefs.references)
			// It's initialized on either the first or second reference.
			firstRead := 2
			if aliasRefs.references[0].IsInitializingDeclaration() {
				firstRead = 1
			}
			collector.forgetVar(aliasVar)
			for i := firstRead; i < size; i++ {
				aliasRef := aliasRefs.references[i]
				newNodes.add(a.replaceAliasReference(alias, aliasRef))
			}

			// just set the original alias to null.
			a.tryReplacingAliasingAssignment(alias, name, lhsNode)

			// Inlining the variable may have introduced new references
			// to descendants of `name`. So those need to be collected now.
			namespace.ScanNewNodes(newNodes)
			return
		}

		if name.IsConstructor() {
			// TODO(lharker): the main reason this was added is because method decomposition inside
			// generators introduces some constructor aliases that weren't getting inlined.
			// If we find another (safer) way to avoid aliasing in method decomposition, consider
			// removing this.
			if !a.partiallyInlineAlias(alias, name, namespace, collector, aliasVar, aliasRefs, lhsNode) {
				// If we can't inline all alias references, make sure there are no unsafe property
				// accesses.
				if InlineAndCollapseProperties_ReferencesCollapsibleProperty(aliasRefs, name) {
					a.outer.compiler.Report(JSError_MakeWithNode(aliasParent, InlineAndCollapseProperties_UNSAFE_CTOR_ALIASING, lhsNode.GetString()))
				}
			}
		}
	}
}

// partiallyInlineAlias inlines some references to an alias with its value.
// This handles cases where the alias is not declared at initialization. It
// does nothing if the alias is reassigned after being initialized, unless
// the reassignment occurs because of an enclosing function or a loop.
//
// alias is an alias of some variable, which may not be well-defined;
// namespace the GlobalNamespace, which will be updated with all new nodes
// created; aliasRefs all references to the alias in its scope; aliasLhsNode
// the lhs name of the alias when it is first initialized. Returns whether
// all references to the alias were inlined.
//
// Java: InlineAndCollapseProperties.AggressiveInlineAliases.partiallyInlineAlias(Ref, Name, GlobalNamespace, ReferenceCollection, Node)
func (a *aggressiveInlineAliases) partiallyInlineAlias(
	alias *GlobalNamespace_Ref,
	aliasingName *GlobalNamespace_Name,
	namespace *GlobalNamespace,
	collector *ReferenceCollector,
	aliasVar *Var,
	aliasRefs *ReferenceCollection,
	aliasLhsNode rhino.Node) bool {
	var aliasBlock *BasicBlock
	// This initial iteration through all the alias references does two things:
	// a) Find the control flow block in which the alias is assigned.
	// b) See if the alias var is assigned to in multiple places, and return if that's the case.
	//    NOTE: we still may inline if the alias is assigned in a loop or inner function and that
	//    assignment statement is potentially executed multiple times.
	//    This is more aggressive than what "inlineAliasIfPossible" does.
	for aliasRef := range aliasRefs.All() {
		aliasRefNode := aliasRef.GetNode()
		if aliasRefNode == aliasLhsNode {
			aliasBlock = aliasRef.GetBasicBlock()
			continue
		} else if aliasRef.IsLvalue() {
			// Don't replace any references if the alias is reassigned
			return false
		}
	}

	newNodes := &a.newNodes
	newNodes.reset()
	alreadySeenInitialAlias := false
	foundNonReplaceableAlias := false
	replacedAny := false
	// Do a second iteration through all the alias references, and replace any inlinable
	// references.
	for aliasRef := range aliasRefs.All() {
		aliasRefNode := aliasRef.GetNode()
		if aliasRefNode == aliasLhsNode {
			alreadySeenInitialAlias = true
			continue
		} else if aliasRef.IsDeclaration() {
			// Ignore any alias declarations, e.g. "var alias;", since there's nothing to inline.
			continue
		}

		refBlock := aliasRef.GetBasicBlock()
		if (refBlock != aliasBlock && aliasBlock.ProvablyExecutesBefore(refBlock)) ||
			(refBlock == aliasBlock && alreadySeenInitialAlias) {
			// We replace the alias only if the alias and reference are in the same BasicBlock,
			// the aliasing assignment takes place before the reference, and the alias is
			// never reassigned.
			a.codeChanged = true
			replacedAny = true
			newNodes.add(a.replaceAliasReference(alias, aliasRef))
		} else {
			foundNonReplaceableAlias = true
		}
	}

	if replacedAny || !foundNonReplaceableAlias {
		collector.forgetVar(aliasVar)
	}

	// We removed all references to the alias, so remove the original aliasing assignment.
	if !foundNonReplaceableAlias {
		a.tryReplacingAliasingAssignment(alias, aliasingName, aliasLhsNode)
	}

	if a.codeChanged {
		// Inlining the variable may have introduced new references
		// to descendants of `name`. So those need to be collected now.
		namespace.ScanNewNodes(newNodes)
	}
	return !foundNonReplaceableAlias
}

// tryReplacingAliasingAssignment replaces the rhs of an aliasing assignment
// with null, unless the assignment result is used in a complex expression.
//
// Java: InlineAndCollapseProperties.AggressiveInlineAliases.tryReplacingAliasingAssignment(Ref, Name, Node)
func (a *aggressiveInlineAliases) tryReplacingAliasingAssignment(
	alias *GlobalNamespace_Ref, aliasName *GlobalNamespace_Name, aliasLhsNode rhino.Node) bool {
	// either VAR/CONST/LET or ASSIGN.
	assignment := aliasLhsNode.GetParent()
	if !NodeUtil_IsNameDeclaration(assignment) && NodeUtil_IsExpressionResultUsed(assignment) {
		// e.g. don't change "if (alias = someVariable)" to "if (alias = null)"
		// TODO(lharker): instead replace the entire assignment with the RHS - "alias = x" becomes
		// "x"
		return false
	}
	aliasParent := alias.GetNode().GetParent()
	a.noteMutatedQName(alias.GetNode())
	alias.GetNode().ReplaceWith(rhino.IR_NullNode())
	aliasName.RemoveRef(alias)
	a.codeChanged = true
	a.outer.compiler.ReportChangeToEnclosingScope(aliasParent)
	return true
}

// replaceAliasReference: alias is a GlobalNamespace.Ref of the variable
// being aliased; aliasRef one particular usage of an alias that we want to
// replace with the aliased var. Returns an AstChange representing the new
// node(s) added to the AST.
//
// Java: InlineAndCollapseProperties.AggressiveInlineAliases.replaceAliasReference(Ref, Reference)
func (a *aggressiveInlineAliases) replaceAliasReference(alias *GlobalNamespace_Ref, aliasRef *Reference) GlobalNamespace_AstChange {
	originalRefNode := alias.GetNode()
	nodeToReplace := aliasRef.GetNode()
	jutil.CheckStateMsg1(nodeToReplace.IsQualifiedName(), "%s", nodeToReplace)
	if nodeToReplace.HasChildren() {
		a.invalidateScopeCollectors()
	}
	// If the reference node is a NAME it could be
	// const origName = value;
	// If we use cloneTree() for that we'll clone the value, which we don't want.
	// Otherwise, we do want to clone the tree of GETPROP nodes.
	var newNode rhino.Node
	if originalRefNode.IsName() {
		newNode = originalRefNode.CloneNode()
	} else {
		newNode = originalRefNode.CloneTree()
	}
	newNode.SrcrefTree(nodeToReplace)
	nodeToReplace.ReplaceWith(newNode)
	a.noteMutatedQName(newNode)
	a.outer.compiler.ReportChangeToEnclosingScope(newNode)
	return NewGlobalNamespace_SimpleAstChange(newNode, alias.GetChunk(), alias.scope)
}

// inlineGlobalAliasIfPossible attempts to inline an global alias of a
// global name. This requires that the name is well defined: assigned
// unconditionally, assigned exactly once. It is assumed that, the name for
// which it is an alias must already meet these same requirements.
//
// If the alias is completely removed, also deletes the aliasing Ref.
//
// name is the global name being aliased; alias the alias to inline.
//
// Java: InlineAndCollapseProperties.AggressiveInlineAliases.inlineGlobalAliasIfPossible(Name, Ref, GlobalNamespace)
func (a *aggressiveInlineAliases) inlineGlobalAliasIfPossible(name *GlobalNamespace_Name, alias *GlobalNamespace_Ref, namespace *GlobalNamespace) {
	// Ensure that the alias is assigned to global name at that the
	// declaration.
	aliasParent := alias.GetNode().GetParent()
	if lvalue, ok := aliasLhsNode(aliasParent); ok &&
		(NodeUtil_IsExecutedExactlyOnce(aliasParent) ||
			// We special-case for constructors here, to inline constructor aliases
			// more aggressively in global scope.
			// We do this because constructor properties are always collapsed,
			// so we want to inline the aliases also to avoid breakages.
			(aliasParent.IsName() && name.IsConstructor())) {
		if lvalue.IsNil() || !lvalue.IsQualifiedName() {
			return
		}
		if lvalue.IsName() &&
			a.outer.compiler.GetCodingConvention().IsExportedWithLocal(lvalue.GetString() /* local= */, false) {
			return
		}
		aliasingName, _ := namespace.lookupGlobalNode(lvalue)

		if aliasingName == nil {
			// this is true for names in externs or properties on extern names
			return
		}

		if name == aliasingName && aliasParent.IsAssign() {
			// Ignore `a.b.c = a.b.c;` with `a.b.c;`.
			return
		}

		aliasInlinability := aliasingName.CalculateInlinability()
		if !aliasInlinability.ShouldInlineUsages() {
			// nothing to do here
			return
		}

		// Rewrite all references to the aliasing name, except for the initialization & any cross-
		// chunk references if they exist.
		newNodes := a.rewriteAliasReferences(aliasingName, alias, aliasInlinability)
		a.rewriteAliasProps(aliasingName, alias.GetNode(), 0, newNodes)

		if aliasInlinability.ShouldRemoveDeclaration() {
			// Rewrite the initialization of the alias.
			aliasInitialization := aliasingName.GetInitialization()
			if aliasInitialization.IsTwin() {
				// This is in a nested assign.
				// Replace
				//   a.b = aliasing.name = aliased.name
				// with
				//   a.b = aliased.name
				assign, ok := ast.AsAssign(aliasParent)
				jutil.CheckStateMsg1(ok, "%s", aliasParent)
				aliasGrandparent := aliasParent.GetParent()
				a.noteMutatedQName(assign.Target().Node())
				a.noteMutatedQName(alias.GetNode())
				aliasParent.ReplaceWith(alias.GetNode().Detach())
				// Remove the ref to 'aliasing.name' entirely
				aliasingName.RemoveRef(aliasInitialization)
				// Force GlobalNamespace to revisit the new reference to 'aliased.name' and update its
				// internal state.
				newNodes.add(NewGlobalNamespace_RefBasedAstChange(alias, alias.GetNode()))
				a.outer.compiler.ReportChangeToEnclosingScope(aliasGrandparent)
			} else {
				// Replace
				//  aliasing.name = aliased.name
				// with
				//  aliasing.name = null;
				a.noteMutatedQName(alias.GetNode())
				alias.GetNode().ReplaceWith(rhino.IR_NullNode())
				a.outer.compiler.ReportChangeToEnclosingScope(aliasParent)
			}
			a.codeChanged = true
			// Update the original aliased name to say that it has one less ALIASING_REF.
			name.RemoveRef(alias)
		}

		// Inlining the variable may have introduced new references
		// to descendants of `name`. So those need to be collected now.
		namespace.ScanNewNodes(newNodes)
	}
}

// rewriteAliasReferences replaces reads of a name with the name it aliases
//
// aliasingName is the alias - e.g. 'alsoFoo' in 'const alsoFoo = foo;';
// aliasInlinability indicates whether all references to the alias should be
// inlined, or only those that are in a valid cross-chunk dependency.
//
// Java: InlineAndCollapseProperties.AggressiveInlineAliases.rewriteAliasReferences(Name, Ref, Inlinability)
func (a *aggressiveInlineAliases) rewriteAliasReferences(
	aliasingName *GlobalNamespace_Name, aliasingRef *GlobalNamespace_Ref, aliasInlinability GlobalNamespace_Inlinability) *astChangeSet {
	newNodes := &a.newNodes
	newNodes.reset()
	var refsBuf [8]*GlobalNamespace_Ref
	refs := aliasingName.snapshotRefs(&refsBuf)
	var initializationChunk *JSChunk
	if aliasingName.GetInitialization() != nil {
		initializationChunk = aliasingName.GetInitialization().GetChunk()
	} else {
		initializationChunk = aliasingName.getDeclarationRef().GetChunk()
	}
	for _, ref := range refs {
		if ref.IsSetFromGlobal() {
			// Handled elsewhere
			// This should either be the initialization or an irrelevant uninitialized declaration.
			jutil.CheckStateMsg1(
				ref.GetChunk() == initializationChunk || ref.IsUninitializedDeclaration(), "%s", ref)
			continue
		}

		jutil.CheckStateMsg1(ref.IsGet(), "%s", ref) // names with local sets should not be rewritten
		// Twin refs that are both gets and sets are handled later, with the other sets.
		jutil.CheckStateMsg1(!ref.IsTwin(), "%s", ref)

		if aliasInlinability == GlobalNamespace_Inlinability_INLINE_UNLESS_INVALID_CROSS_CHUNK_DEPENDENCY &&
			ref.GetChunk() != initializationChunk &&
			!a.outer.compiler.GetChunkGraph().DependsOn(ref.GetChunk(), initializationChunk) {
			continue
		}

		if ref.GetNode().IsStringKey() {
			// e.g. `y` in `const {y} = x;`
			if ref.scope != nil && ref.scope.IsLocal() {
				a.invalidateScopeCollectors()
			}
			DestructuringGlobalNameExtractor_ReassignDestructringLvalue(
				ref.GetNode(), aliasingRef.GetNode().CloneTree(), newNodes, ref, a.outer.compiler)
		} else {
			// e.g. `x.y`
			jutil.CheckState(ref.GetNode().IsGetProp() || ref.GetNode().IsName())
			newNode := aliasingRef.GetNode().CloneTree()
			node := ref.GetNode()
			a.noteMutatedQName(node)
			a.noteMutatedQName(newNode)
			newNode.Srcref(node)
			node.ReplaceWith(newNode)
			a.outer.compiler.ReportChangeToEnclosingScope(newNode)
			newNodes.add(NewGlobalNamespace_RefBasedAstChange(ref, newNode))
		}
		aliasingName.RemoveRef(ref)
	}
	return newNodes
}

// rewriteAliasProps: name is the Name whose properties references should be
// updated; value the value to use when rewriting; depth the chain depth;
// newNodes expression nodes that have been updated.
//
// Java: InlineAndCollapseProperties.AggressiveInlineAliases.rewriteAliasProps(Name, Node, int, Set<AstChange>)
func (a *aggressiveInlineAliases) rewriteAliasProps(
	name *GlobalNamespace_Name, value rhino.Node, depth int, newNodes *astChangeSet) {
	if !name.HasProps() {
		return
	}
	if name.MatchesNode(value) {
		panic(jutil.NewIllegalStateExceptionMsg(
			jutil.LenientFormat("%s should not match name %s", value, name.GetFullName())))
	}
	for prop := range name.Props() {
		a.rewriteNestedAliasReference(value, depth, newNodes, prop)
	}
}

// rewriteNestedAliasReference replaces references to an alias that are
// nested inside a longer getprop chain or an object literal
//
// For example: if we have an inlined alias 'const A = B;', and reference a
// property 'A.x', then this method is responsible for replacing 'A.x' with
// 'B.x'.
//
// This is necessary because in the above example, given 'A.x', there is
// only one Ref that points to the whole name 'A.x', not a direct Ref to 'A'.
// So the only way to replace 'A.x' with 'B.x' is by looking at the property
// 'x' reference.
//
// value is the value to use when rewriting; depth the property chain depth;
// newNodes expression nodes that have been updated; prop the property to
// rewrite with value.
//
// Java: InlineAndCollapseProperties.AggressiveInlineAliases.rewriteNestedAliasReference(Node, int, Set<AstChange>, Name)
func (a *aggressiveInlineAliases) rewriteNestedAliasReference(
	value rhino.Node, depth int, newNodes *astChangeSet, prop *GlobalNamespace_Name) {
	a.rewriteAliasProps(prop, value, depth+1, newNodes)
	var refsBuf [8]*GlobalNamespace_Ref
	refs := prop.snapshotRefs(&refsBuf)
	for _, ref := range refs {
		target := ref.GetNode()
		if _, ok := ast.AsStringKey(target); ok && target.GetParent().IsDestructuringPattern() {
			// Do nothing for alias properties accessed through object destructuring. This would be
			// redundant. This method is intended for names nested inside getprop chains, because
			// GlobalNamespace only creates a single Ref for the outermost getprop. However, for
			// destructuring property accesses, GlobalNamespace creates multiple Refs, one for the
			// destructured object, and one for each string key in the pattern.
			//
			// For example, consider:
			//   const originalObj = {key: 0};
			//   const rhs = originalObj;
			//   const {key: lhs} = rhs;
			//   const otherLhs = rhs.key;
			// AggressiveInlineAliases is inlining rhs -> originalObj.
			//
			// GlobalNamespace creates two Refs for the name 'rhs': one for its declaration,
			// and one for 'const {key: lhs} = rhs;'. There is no Ref pointing directly to the 'rhs'
			// in 'const otherLhs = rhs.key', though.
			// There are also two Refs to the name 'rhs.key': one for the destructuring access and one
			// for the getprop access. This loop will visit both Refs.
			// This method is responsible for inlining "const otherLhs = originalObj.key" but not
			// "const {key: lhs} = originalObj;". We bail out at the Ref in the latter case.
			gparent := target.GetGrandparent()
			_, inAssign := ast.AsAssign(gparent)
			_, inDestructLhs := ast.AsDestructuringLhs(gparent)
			jutil.CheckStateMsg1(
				inAssign || inDestructLhs,
				// Currently GlobalNamespace doesn't create Refs for 'b' in const {a: {b}} = obj;
				// If it does start creating those Refs, we may have to update this method to handle
				// them explicitly.
				"Did not expect GlobalNamespace to create Ref for key in nested object pattern %s",
				target)
			continue
		}

		if target.IsMemberFunctionDef() {
			// In the case of a method reference, the reference to the alias is not fully qualified,
			// so we should not touch it, see b/293184904.
			continue
		}

		for i := 0; i <= depth; i++ {
			if gp, ok := ast.AsGetProp(target); ok {
				target = gp.Object().Node()
			} else if NodeUtil_IsObjectLitKey(target) {
				// Object literal key definitions are a little trickier, as we
				// need to find the assignment target
				gparent := target.GetGrandparent()
				if assign, ok := ast.AsAssign(gparent); ok {
					target = assign.Target().Node()
				} else {
					jutil.CheckState(NodeUtil_IsObjectLitKey(gparent))
					target = gparent
				}
			} else {
				panic(jutil.NewIllegalStateExceptionMsg("unexpected node: " + target.String()))
			}
		}
		jutil.CheckState(target.IsGetProp() || target.IsName())
		newValue := value.CloneTree()
		a.noteMutatedQName(target)
		a.noteMutatedQName(newValue)
		target.ReplaceWith(newValue)
		a.outer.compiler.ReportChangeToEnclosingScope(newValue)
		prop.RemoveRef(ref)
		// Rescan the expression root.
		newNodes.add(NewGlobalNamespace_RefBasedAstChange(ref, ref.GetNode()))
		a.codeChanged = true
	}
}

// rewriteSimpleDestructuringAliases rewrites "simple" destructuring aliases to
// a format that is more amenable to inlining.
//
// To be specific, this rewrites aliases of the form: const {x} =
// qualified.name; to: const x = qualified.name.x;
type rewriteSimpleDestructuringAliases struct {
	NodeTraversal_AbstractPostOrderCallbackBase
}

func parseSimpleDestructuringAlias(n rhino.Node) (ast.ObjectPattern, rhino.Node, bool) {
	if !NodeUtil_IsStatement(n) {
		return ast.ObjectPattern{}, rhino.NilNode, false
	}
	decl, ok := ast.AsNameDeclarationOf(n, rhino.Token_CONST)
	if !ok {
		return ast.ObjectPattern{}, rhino.NilNode, false
	}
	declarator, ok := decl.OnlyDeclarator()
	jutil.CheckState(ok)
	destructuringLhs, ok := declarator.AsDestructuringLhs()
	if !ok {
		return ast.ObjectPattern{}, rhino.NilNode, false
	}
	objectPattern, ok := destructuringLhs.Pattern().AsObjectPattern()
	if !ok {
		return ast.ObjectPattern{}, rhino.NilNode, false
	}
	initExpr, ok := destructuringLhs.Init()
	if !ok || !initExpr.Node().IsQualifiedName() {
		return ast.ObjectPattern{}, rhino.NilNode, false
	}
	if !isSimpleDestructuringPattern(objectPattern) {
		return ast.ObjectPattern{}, rhino.NilNode, false
	}
	return objectPattern, initExpr.Node(), true
}

// Java: InlineAndCollapseProperties.RewriteSimpleDestructuringAliases.isSimpleDestructuringPattern(Node)
func isSimpleDestructuringPattern(objectPattern ast.ObjectPattern) bool {
	for prop := range objectPattern.Props() {
		key, ok := prop.AsStringKey()
		if !ok || key.IsQuoted() {
			return false
		}
		jutil.CheckState(key.Node().HasOneChild())
		val := key.Value()
		if _, ok := val.AsName(); ok {
			continue
		}
		nested, ok := val.AsObjectPattern()
		if !ok || !isSimpleDestructuringPattern(nested) {
			return false
		}
	}
	return true
}

var rewriteSimpleDestructuringTokens = makeVisitTokenSet(rhino.Token_CONST)

var (
	_ nodeTraversal_visitTokens   = (*rewriteSimpleDestructuringAliases)(nil)
	_ nodeTraversal_noVisitLeaves = (*rewriteSimpleDestructuringAliases)(nil)
)

func (*rewriteSimpleDestructuringAliases) visitTokens() *ast.TokenSet {
	return &rewriteSimpleDestructuringTokens
}

func (*rewriteSimpleDestructuringAliases) noVisitLeaves() {}

// Java: InlineAndCollapseProperties.RewriteSimpleDestructuringAliases.visit(NodeTraversal, Node, Node)
func (r *rewriteSimpleDestructuringAliases) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	objectPattern, rhs, ok := parseSimpleDestructuringAlias(n)
	if !ok {
		return
	}
	expandObjectPattern(t, n, objectPattern, rhs)
	n.Detach()
	t.ReportCodeChange()
}

// Java: InlineAndCollapseProperties.RewriteSimpleDestructuringAliases.expandObjectPattern(NodeTraversal, Node, Node, Node)
func expandObjectPattern(
	t *NodeTraversal, insertionPoint rhino.Node, objectPattern ast.ObjectPattern, rhs rhino.Node) rhino.Node {
	for prop := range objectPattern.Props() {
		key := prop.MustStringKey()
		val := key.Value()
		keyChild := val.Node()
		var nameNode rhino.Node
		var nested ast.ObjectPattern
		var isNested bool
		if name, ok := val.AsName(); ok {
			nameNode = name.Detach().Node()
		} else {
			nested, isNested = val.AsObjectPattern()
			jutil.CheckState(isNested)
			uniqueId := t.GetCompiler().GetUniqueIdSupplier().GetUniqueId(t.GetInput())
			newName := ast.NewName("destructuring$" + uniqueId)
			newName.SetConstantName(true)
			nameNode = newName.Node().Srcref(keyChild)
		}
		newRhs := rhino.IR_Getprop(rhs.CloneTree(), key.Key()).Srcref(keyChild)
		newConstNode := rhino.IR_ConstNode(nameNode, newRhs).Srcref(objectPattern.Node())
		newConstNode.InsertAfter(insertionPoint)
		insertionPoint = newConstNode
		if isNested {
			insertionPoint = expandObjectPattern(t, insertionPoint, nested, nameNode.CloneNode())
		}
	}
	return insertionPoint
}
