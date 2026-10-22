// Port of com/google/javascript/jscomp/InlineAndCollapseProperties.java @ closure-compiler v20261015

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
	"reflect"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/diagnostic"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// InlineAndCollapseProperties performs inlining of aliases and collapsing of
// qualified names in order to improve later optimizations, such as
// RemoveUnusedCode.
type InlineAndCollapseProperties struct {
	compiler                 AbstractCompiler
	propertyCollapseLevel    CompilerOptions_PropertyCollapseLevel
	chunkOutputType          CompilerOptions_ChunkOutputType
	haveModulesBeenRewritten bool
	moduleResolutionMode     ModuleLoader_ResolutionMode

	// Used by `AggressiveInlineAliasesTest` to enable execution of the
	// aggressive inlining logic without doing any collapsing.
	testAggressiveInliningOnly bool

	// Supplied by `AggressiveInlineAliasesTest`.
	//
	// The `GlobalNamespace` created by `AggressiveInlineAliases` will be
	// passed to this `Consumer` for examination.
	//
	// Go-only: language-adapter — Optional<Consumer<GlobalNamespace>>; nil is
	// Optional.empty().
	optionalGlobalNamespaceTester func(*GlobalNamespace)

	// Records decisions made by this class and related logic.
	//
	// This field is allocated and cleaned up by process(). It's a class field
	// to avoid having to pass it as an extra argument through a lot of
	// methods.
	decisionsLog diagnostic.LogFile

	// A `GlobalNamespace` that is shared by alias inlining and property
	// collapsing code.
	namespace *GlobalNamespace
}

// Warnings

// Java: InlineAndCollapseProperties.PARTIAL_NAMESPACE_WARNING
var InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING = DiagnosticType_Warning(
	"JSC_PARTIAL_NAMESPACE",
	"Partial alias created for namespace {0}, possibly due to await/yield transpilation.\n"+
		"This may prevent optimization of anything nested under this namespace.\n"+
		"See https://github.com/google/closure-compiler/wiki/FAQ#i-got-a-partial-alias-created-for-namespace-error--what-do-i-do"+
		" for more details.")

// Java: InlineAndCollapseProperties.NAMESPACE_REDEFINED_WARNING
var InlineAndCollapseProperties_NAMESPACE_REDEFINED_WARNING = DiagnosticType_Warning("JSC_NAMESPACE_REDEFINED", "namespace {0} should not be redefined")

// Java: InlineAndCollapseProperties.RECEIVER_AFFECTED_BY_COLLAPSE
var InlineAndCollapseProperties_RECEIVER_AFFECTED_BY_COLLAPSE = DiagnosticType_Warning(
	"JSC_RECEIVER_AFFECTED_BY_COLLAPSE",
	"Receiver reference in function {0} changes meaning when namespace is collapsed.\n"+
		" Consider annotating @nocollapse; however, other properties on the receiver may"+
		" still be collapsed.")

// Java: InlineAndCollapseProperties.UNSAFE_CTOR_ALIASING
var InlineAndCollapseProperties_UNSAFE_CTOR_ALIASING = DiagnosticType_Warning(
	"JSC_UNSAFE_CTOR_ALIASING",
	"Variable {0} aliases a constructor, so it cannot be assigned multiple times")

// Java: InlineAndCollapseProperties.ALIAS_CYCLE
var InlineAndCollapseProperties_ALIAS_CYCLE = DiagnosticType_Error("JSC_ALIAS_CYCLE", "Alias path contains a cycle: {0} to {1}")

// Java: InlineAndCollapseProperties(Builder)
func NewInlineAndCollapseProperties(builder *InlineAndCollapseProperties_Builder) *InlineAndCollapseProperties {
	p := &InlineAndCollapseProperties{decisionsLog: nil}
	p.compiler = builder.compiler
	p.propertyCollapseLevel = builder.propertyCollapseLevel
	p.chunkOutputType = builder.chunkOutputType
	p.haveModulesBeenRewritten = builder.haveModulesBeenRewritten
	p.moduleResolutionMode = builder.moduleResolutionMode
	p.testAggressiveInliningOnly = builder.testAggressiveInliningOnly
	p.optionalGlobalNamespaceTester = builder.optionalGlobalNamespaceTester
	return p
}

// InlineAndCollapseProperties_Builder is InlineAndCollapseProperties.Builder.
type InlineAndCollapseProperties_Builder struct {
	compiler                      AbstractCompiler
	propertyCollapseLevel         CompilerOptions_PropertyCollapseLevel
	chunkOutputType               CompilerOptions_ChunkOutputType
	haveModulesBeenRewritten      bool
	moduleResolutionMode          ModuleLoader_ResolutionMode
	testAggressiveInliningOnly    bool
	optionalGlobalNamespaceTester func(*GlobalNamespace)
}

// Java: InlineAndCollapseProperties.Builder(AbstractCompiler)
func NewInlineAndCollapseProperties_Builder(compiler AbstractCompiler) *InlineAndCollapseProperties_Builder {
	return &InlineAndCollapseProperties_Builder{
		compiler:                      compiler,
		testAggressiveInliningOnly:    false,
		optionalGlobalNamespaceTester: nil,
	}
}

// Java: InlineAndCollapseProperties.Builder.setPropertyCollapseLevel(PropertyCollapseLevel)
func (b *InlineAndCollapseProperties_Builder) SetPropertyCollapseLevel(propertyCollapseLevel CompilerOptions_PropertyCollapseLevel) *InlineAndCollapseProperties_Builder {
	b.propertyCollapseLevel = propertyCollapseLevel
	return b
}

// Java: InlineAndCollapseProperties.Builder.setChunkOutputType(ChunkOutputType)
func (b *InlineAndCollapseProperties_Builder) SetChunkOutputType(chunkOutputType CompilerOptions_ChunkOutputType) *InlineAndCollapseProperties_Builder {
	b.chunkOutputType = chunkOutputType
	return b
}

// Java: InlineAndCollapseProperties.Builder.setHaveModulesBeenRewritten(boolean)
func (b *InlineAndCollapseProperties_Builder) SetHaveModulesBeenRewritten(haveModulesBeenRewritten bool) *InlineAndCollapseProperties_Builder {
	b.haveModulesBeenRewritten = haveModulesBeenRewritten
	return b
}

// Java: InlineAndCollapseProperties.Builder.setModuleResolutionMode(ResolutionMode)
func (b *InlineAndCollapseProperties_Builder) SetModuleResolutionMode(moduleResolutionMode ModuleLoader_ResolutionMode) *InlineAndCollapseProperties_Builder {
	b.moduleResolutionMode = moduleResolutionMode
	return b
}

// Java: InlineAndCollapseProperties.Builder.testAggressiveInliningOnly(Consumer<GlobalNamespace>)
func (b *InlineAndCollapseProperties_Builder) TestAggressiveInliningOnly(globalNamespaceTester func(*GlobalNamespace)) *InlineAndCollapseProperties_Builder {
	b.testAggressiveInliningOnly = true
	b.optionalGlobalNamespaceTester = jutil.CheckNotNull(globalNamespaceTester)
	return b
}

// Java: InlineAndCollapseProperties.Builder.build()
func (b *InlineAndCollapseProperties_Builder) Build() *InlineAndCollapseProperties {
	return NewInlineAndCollapseProperties(b)
}

// InlineAndCollapseProperties_Builder_ is the static builder(compiler).
//
// Java: InlineAndCollapseProperties.builder(AbstractCompiler)
func InlineAndCollapseProperties_Builder_(compiler AbstractCompiler) *InlineAndCollapseProperties_Builder {
	return NewInlineAndCollapseProperties_Builder(compiler)
}

// Java: InlineAndCollapseProperties.process(Node, Node)
func (p *InlineAndCollapseProperties) Process(externs *rhino.Node, root *rhino.Node) {
	jutil.CheckStateMsg1(
		!p.testAggressiveInliningOnly || p.propertyCollapseLevel == CompilerOptions_PropertyCollapseLevel_ALL,
		"testAggressiveInlining is invalid for: %s",
		p.propertyCollapseLevel)
	defer func() {
		p.decisionsLog = nil
		if p.namespace != nil && p.optionalGlobalNamespaceTester == nil {
			p.namespace.Release()
			p.namespace = nil
		}
		evictCachedGlobalNamespaceArenas()
	}()
	logFile := p.compiler.CreateOrReopenIndexedLog(reflect.TypeOf(p).Elem(), "decisions.log")
	defer logFile.Close()
	// NOTE: decisionsLog will be a do-nothing proxy object unless the compiler
	// was given an option telling it to generate log files and where to put them.
	p.decisionsLog = logFile
	switch p.propertyCollapseLevel {
	case CompilerOptions_PropertyCollapseLevel_NONE:
		p.PerformMinimalInliningAndNoCollapsing(externs, root)
	case CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT:
		p.PerformMinimalInliningAndModuleExportCollapsing(externs, root)
	case CompilerOptions_PropertyCollapseLevel_ALL:
		if p.testAggressiveInliningOnly {
			p.PerformAggressiveInliningForTest(externs, root)
		} else {
			p.PerformAggressiveInliningAndCollapsing(externs, root)
		}
	}
}

// Java: InlineAndCollapseProperties.performMinimalInliningAndNoCollapsing(Node, Node)
func (p *InlineAndCollapseProperties) PerformMinimalInliningAndNoCollapsing(externs *rhino.Node, root *rhino.Node) {
	// TODO(b/124915436): Remove InlineAliases completely after cleaning up the codebase.
	NewInlineAndCollapseProperties_InlineAliases(p).Process(externs, root)
}

// Java: InlineAndCollapseProperties.performMinimalInliningAndModuleExportCollapsing(Node, Node)
func (p *InlineAndCollapseProperties) PerformMinimalInliningAndModuleExportCollapsing(externs *rhino.Node, root *rhino.Node) {
	// TODO(b/124915436): Remove InlineAliases completely after cleaning up the codebase.
	NewInlineAndCollapseProperties_InlineAliases(p).Process(externs, root)
	// CollapseProperties needs this namespace.
	// TODO(bradfordcsmith): Have `InlineAliases` update the namespace it already created
	// and reuse that one instead.
	if p.namespace != nil {
		p.namespace.Release()
	}
	p.namespace = NewGlobalNamespaceWithLog(p.decisionsLog, p.compiler, root)
	NewInlineAndCollapseProperties_CollapseProperties(p).Process(externs, root)
}

// Java: InlineAndCollapseProperties.performAggressiveInliningAndCollapsing(Node, Node)
func (p *InlineAndCollapseProperties) PerformAggressiveInliningAndCollapsing(externs *rhino.Node, root *rhino.Node) {
	NewConcretizeStaticInheritanceForInlining(p.compiler).Process(externs, root)
	NewInlineAndCollapseProperties_AggressiveInlineAliases(p).Process(externs, root)

	NewInlineAndCollapseProperties_CollapseProperties(p).Process(externs, root)
}

// Java: InlineAndCollapseProperties.performAggressiveInliningForTest(Node, Node)
func (p *InlineAndCollapseProperties) PerformAggressiveInliningForTest(externs *rhino.Node, root *rhino.Node) {
	aggressiveInlineAliases := NewInlineAndCollapseProperties_AggressiveInlineAliases(p)
	aggressiveInlineAliases.Process(externs, root)
	if p.optionalGlobalNamespaceTester == nil {
		// Optional.get() on an empty Optional
		panic(jutil.NewNoSuchElementExceptionMsg("No value present"))
	}
	p.optionalGlobalNamespaceTester(aggressiveInlineAliases.GetLastUsedGlobalNamespace())
}

// InlineAndCollapseProperties_AggressiveInlineAliases inlines type aliases
// if they are explicitly or effectively const. Also inlines inherited static
// property accesses for ES6 classes.
//
// This frees subsequent optimization passes from the responsibility of
// having to reason about alias chains and is a requirement for correct
// behavior in at least CollapseProperties and J2clPropertyInlinerPass.
//
// This is designed to be no more unsafe than CollapseProperties. It will in
// some cases inline properties, possibly past places that change the
// property value. However, it will only do so in cases where
// CollapseProperties would unsafely collapse the property anyway.
type InlineAndCollapseProperties_AggressiveInlineAliases struct {
	outer *InlineAndCollapseProperties

	codeChanged bool
}

// Java: InlineAndCollapseProperties.AggressiveInlineAliases()
func NewInlineAndCollapseProperties_AggressiveInlineAliases(outer *InlineAndCollapseProperties) *InlineAndCollapseProperties_AggressiveInlineAliases {
	a := &InlineAndCollapseProperties_AggressiveInlineAliases{outer: outer}
	a.codeChanged = true
	return a
}

// Java: InlineAndCollapseProperties.AggressiveInlineAliases.getLastUsedGlobalNamespace()
func (a *InlineAndCollapseProperties_AggressiveInlineAliases) GetLastUsedGlobalNamespace() *GlobalNamespace {
	return a.outer.namespace
}

// Java: InlineAndCollapseProperties.AggressiveInlineAliases.process(Node, Node)
func (a *InlineAndCollapseProperties_AggressiveInlineAliases) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(a.outer.compiler, root, NewInlineAndCollapseProperties_RewriteSimpleDestructuringAliases())

	// Building the `GlobalNamespace` dominates the cost of this pass, so it is built once and
	// updated as changes are made so it can be reused for the next iteration.
	if a.outer.namespace != nil {
		a.outer.namespace.Release()
	}
	a.outer.namespace = NewGlobalNamespaceWithLog(a.outer.decisionsLog, a.outer.compiler, root)
	for a.codeChanged {
		a.codeChanged = false
		a.InlineAliases(a.outer.namespace)
	}
}

// InlineAliases: for each qualified name N in the global scope, we check if:
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
func (a *InlineAndCollapseProperties_AggressiveInlineAliases) InlineAliases(namespace *GlobalNamespace) {
	// Invariant: All the names in the worklist meet condition (a).
	// adds all top-level names to the worklist, but not any properties on those names.
	//
	// Go-only: language-adapter — the ArrayDeque is a slice; pop() takes the
	// first element and add() appends.
	workList := namespace.GetNameForest().ToSlice()

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
			a.MaybeInlineInnerName(name)
			if name.GetAliasingGets() > 0 || name.GetSubclassingGets() > 0 {
				// condition (c) and/or condition (d) are true
				a.InlineAliasesForName(name, namespace)
			}
		}
		InlineAndCollapseProperties_MaybeAddPropertiesToWorklist(name, &workList)
	}
}

// InlineAliasesForName inlines a global name into all the places where
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
func (a *InlineAndCollapseProperties_AggressiveInlineAliases) InlineAliasesForName(name *GlobalNamespace_Name, namespace *GlobalNamespace) {
	refs := name.GetRefs()
	for _, ref := range refs {
		if ref.IsAliasingGet() && !InlineAndCollapseProperties_MayBeGlobalAlias(ref) && !ref.IsTwin() {
			// `name` meets condition (c). Try to inline it.
			// TODO(johnlenz): consider picking up new aliases at the end
			// of the pass instead of immediately like we do for global
			// inlines.
			a.InlineAliasIfPossible(name, ref, namespace)
		} else if ref.IsAliasingGet() &&
			ref.scope.GetClosestHoistScope().IsGlobal() &&
			!ref.IsTwin() {
			a.InlineGlobalAliasIfPossible(name, ref, namespace)
		} else if name.IsClass() && ref.IsSubclassingGet() && name.props != nil {
			for prop := range name.props.All() {
				a.RewriteAllSubclassInheritedAccesses(name, ref, prop, namespace)
			}
		}
	}
}

// MaybeInlineInnerName: if the global name is a function with an
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
func (a *InlineAndCollapseProperties_AggressiveInlineAliases) MaybeInlineInnerName(globalName *GlobalNamespace_Name) {
	globalNameDeclaration := jutil.CheckNotNullMsg1(globalName.getDeclarationRef(), "%s", globalName)
	globalDeclarationNode := jutil.CheckNotNullMsg1(globalNameDeclaration.GetNode(), "%s", globalNameDeclaration)
	valueNode := NodeUtil_GetRValueOfLValue(globalDeclarationNode)
	if valueNode == nil {
		// no function or class expression, so no inner name
		return
	}
	innerNameNode := InlineAndCollapseProperties_MaybeGetInnerNameNode(valueNode)
	if innerNameNode == nil {
		// no inner name to require inlining
		return
	}
	innerName := innerNameNode.GetString()
	syntacticScopeCreator := NewSyntacticScopeCreator(a.outer.compiler)
	syntacticScopeCreator.EnableArenaPool()
	innerScope := syntacticScopeCreator.CreateScope(valueNode, globalNameDeclaration.scope).(*Scope)
	innerNameVar := jutil.CheckNotNull(innerScope.GetVar(innerName))
	collector := NewReferenceCollectorWithVarFilter(
		a.outer.compiler,
		ReferenceCollector_DO_NOTHING_BEHAVIOR,
		syntacticScopeCreator,
		func(v *Var) bool { return jutil.ObjectsEquals(innerNameVar, v) })
	collector.enableArenaPool()
	defer collector.release()
	collector.ProcessScope(innerScope)
	innerNameRefs := collector.GetReferences(innerNameVar)

	newNodes := jutil.NewLinkedHashSet[GlobalNamespace_AstChange]()

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
			newNodes.Add(a.ReplaceAliasReference(globalNameDeclaration, innerNameRef))
		}
	}
	a.outer.namespace.ScanNewNodes(newNodes)
}

// RewriteAllSubclassInheritedAccesses inlines all references to inherited
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
func (a *InlineAndCollapseProperties_AggressiveInlineAliases) RewriteAllSubclassInheritedAccesses(
	superclassNameObj *GlobalNamespace_Name, superclassRef *GlobalNamespace_Ref, prop *GlobalNamespace_Name, namespace *GlobalNamespace) bool {
	if !prop.CanCollapse() {
		return false // inlining is a) unnecessary if there is @nocollapse and b) might break
		// usages of `this` in the method
	}
	subclass := InlineAndCollapseProperties_GetSubclassForEs6Superclass(superclassRef.GetNode())
	if subclass == nil {
		return false
	}
	subclassName, ok := namespace.qnameStringCached(subclass)
	if !ok {
		return false
	}
	subclassQualifiedPropName := namespace.concatDotCached(subclassName, prop.GetBaseName())
	subclassPropNameObj, _ := namespace.GetOwnSlot(subclassQualifiedPropName).(*GlobalNamespace_Name)
	// Don't rewrite if the subclass ever shadows the parent static property.
	// This may also back off on cases where the subclass first accesses the parent property, then
	// shadows it.
	if subclassPropNameObj != nil &&
		(subclassPropNameObj.GetLocalSets() > 0 || subclassPropNameObj.GetGlobalSets() > 0) {
		return false
	}

	// Recurse to find potential sub-subclass accesses of the superclass property.
	subclassNameObj, _ := namespace.GetOwnSlot(subclassName).(*GlobalNamespace_Name)
	if subclassNameObj != nil && subclassNameObj.SubclassingGetCount() > 0 {
		for _, ref := range subclassNameObj.GetRefs() {
			if ref.IsSubclassingGet() {
				a.RewriteAllSubclassInheritedAccesses(superclassNameObj, ref, prop, namespace)
			}
		}
	}

	if subclassPropNameObj != nil {
		newNodes := jutil.NewLinkedHashSet[GlobalNamespace_AstChange]()

		// Use this node as a template for rewriteNestedAliasReference.
		superclassNameNode := superclassNameObj.getDeclarationRef().GetNode()
		if superclassNameNode.IsName() {
			superclassNameNode = superclassNameNode.CloneNode()
		} else if superclassNameNode.IsGetProp() {
			superclassNameNode = superclassNameNode.CloneTree()
		} else {
			return false
		}

		a.RewriteNestedAliasReference(superclassNameNode, 0, newNodes, subclassPropNameObj)
		namespace.ScanNewNodes(newNodes)
	}
	return true
}

// IsInlineableModuleExportsAlias recognizes aliases for the special global
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
func (a *InlineAndCollapseProperties_AggressiveInlineAliases) IsInlineableModuleExportsAlias(name *GlobalNamespace_Name, aliasRefs *ReferenceCollection) bool {
	aliasedNameStr := name.GetFullName()
	if !strings.HasPrefix(aliasedNameStr, ClosureRewriteModule_MODULE_EXPORTS_PREFIX) {
		return false
	}
	// the first rhs value Node we find assigned to the name
	var firstRhs *rhino.Node
	size := aliasRefs.references.Size()
	for i := 0; i < size; i++ {
		ref := aliasRefs.references.Get(i)
		if ref.IsVarDeclaration() || ref.GetParent().IsAssign() {
			var rhs *rhino.Node
			if ref.GetParent().IsAssign() {
				rhs = ref.GetParent().GetSecondChild()
			} else {
				rhs = ref.GetNode().GetFirstChild()
			}

			if rhs != nil {
				// Make sure that every time a value is assigned to the alias, it is
				// the aliased name we expect. If not, we cannot inline this alias.
				if firstRhs == nil {
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

// InlineAliasIfPossible attempts to inline a non-global alias of a global
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
func (a *InlineAndCollapseProperties_AggressiveInlineAliases) InlineAliasIfPossible(name *GlobalNamespace_Name, alias *GlobalNamespace_Ref, namespace *GlobalNamespace) {
	// Ensure that the alias is assigned to a local variable at that
	// variable's declaration. If the alias's parent is a NAME,
	// then the NAME must be the child of a VAR, LET, or CONST node, and we must
	// be in a VAR, LET, or CONST assignment.
	// Otherwise if the parent is an assign, we are in a "a = alias" case.
	aliasParent := alias.GetNode().GetParent()
	if aliasParent.IsName() || aliasParent.IsAssign() {
		var aliasLhsNode *rhino.Node
		if aliasParent.IsName() {
			aliasLhsNode = aliasParent
		} else {
			aliasLhsNode = aliasParent.GetFirstChild()
		}
		aliasVarName := aliasLhsNode.GetString()

		aliasVar := alias.scope.GetVar(aliasVarName)
		jutil.CheckStateMsg1(aliasVar != nil, "Expected variable to be defined in scope (%s)", aliasVarName)
		collector := NewReferenceCollectorWithVarFilter(
			a.outer.compiler,
			ReferenceCollector_DO_NOTHING_BEHAVIOR,
			NewSyntacticScopeCreator(a.outer.compiler),
			func(v *Var) bool { return jutil.ObjectsEquals(aliasVar, v) })
		collector.enableArenaPool()
		defer collector.release()
		aliasScope := aliasVar.GetScope().(*Scope)
		collector.ProcessScope(aliasScope)

		aliasRefs := collector.GetReferences(aliasVar)
		newNodes := jutil.NewLinkedHashSet[GlobalNamespace_AstChange]()

		if aliasRefs.IsWellDefined() &&
			(aliasRefs.IsAssignedOnceInLifetime() ||
				a.IsInlineableModuleExportsAlias(name, aliasRefs)) {
			// The alias is well-formed, so do the inlining now.
			size := aliasRefs.references.Size()
			// It's initialized on either the first or second reference.
			firstRead := 2
			if aliasRefs.references.Get(0).IsInitializingDeclaration() {
				firstRead = 1
			}
			for i := firstRead; i < size; i++ {
				aliasRef := aliasRefs.references.Get(i)
				newNodes.Add(a.ReplaceAliasReference(alias, aliasRef))
			}

			// just set the original alias to null.
			a.TryReplacingAliasingAssignment(alias, name, aliasLhsNode)

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
			if !a.PartiallyInlineAlias(alias, name, namespace, aliasRefs, aliasLhsNode) {
				// If we can't inline all alias references, make sure there are no unsafe property
				// accesses.
				if InlineAndCollapseProperties_ReferencesCollapsibleProperty(aliasRefs, name, namespace) {
					a.outer.compiler.Report(JSError_MakeWithNode(aliasParent, InlineAndCollapseProperties_UNSAFE_CTOR_ALIASING, aliasVarName))
				}
			}
		}
	}
}

// PartiallyInlineAlias inlines some references to an alias with its value.
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
func (a *InlineAndCollapseProperties_AggressiveInlineAliases) PartiallyInlineAlias(
	alias *GlobalNamespace_Ref,
	aliasingName *GlobalNamespace_Name,
	namespace *GlobalNamespace,
	aliasRefs *ReferenceCollection,
	aliasLhsNode *rhino.Node) bool {
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

	newNodes := jutil.NewLinkedHashSet[GlobalNamespace_AstChange]()
	alreadySeenInitialAlias := false
	foundNonReplaceableAlias := false
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
			newNodes.Add(a.ReplaceAliasReference(alias, aliasRef))
		} else {
			foundNonReplaceableAlias = true
		}
	}

	// We removed all references to the alias, so remove the original aliasing assignment.
	if !foundNonReplaceableAlias {
		a.TryReplacingAliasingAssignment(alias, aliasingName, aliasLhsNode)
	}

	if a.codeChanged {
		// Inlining the variable may have introduced new references
		// to descendants of `name`. So those need to be collected now.
		namespace.ScanNewNodes(newNodes)
	}
	return !foundNonReplaceableAlias
}

// TryReplacingAliasingAssignment replaces the rhs of an aliasing assignment
// with null, unless the assignment result is used in a complex expression.
//
// Java: InlineAndCollapseProperties.AggressiveInlineAliases.tryReplacingAliasingAssignment(Ref, Name, Node)
func (a *InlineAndCollapseProperties_AggressiveInlineAliases) TryReplacingAliasingAssignment(
	alias *GlobalNamespace_Ref, aliasName *GlobalNamespace_Name, aliasLhsNode *rhino.Node) bool {
	// either VAR/CONST/LET or ASSIGN.
	assignment := aliasLhsNode.GetParent()
	if !NodeUtil_IsNameDeclaration(assignment) && NodeUtil_IsExpressionResultUsed(assignment) {
		// e.g. don't change "if (alias = someVariable)" to "if (alias = null)"
		// TODO(lharker): instead replace the entire assignment with the RHS - "alias = x" becomes
		// "x"
		return false
	}
	aliasParent := alias.GetNode().GetParent()
	alias.GetNode().ReplaceWith(rhino.IR_NullNode())
	aliasName.RemoveRef(alias)
	a.codeChanged = true
	a.outer.compiler.ReportChangeToEnclosingScope(aliasParent)
	return true
}

// ReplaceAliasReference: alias is a GlobalNamespace.Ref of the variable
// being aliased; aliasRef one particular usage of an alias that we want to
// replace with the aliased var. Returns an AstChange representing the new
// node(s) added to the AST.
//
// Java: InlineAndCollapseProperties.AggressiveInlineAliases.replaceAliasReference(Ref, Reference)
func (a *InlineAndCollapseProperties_AggressiveInlineAliases) ReplaceAliasReference(alias *GlobalNamespace_Ref, aliasRef *Reference) GlobalNamespace_AstChange {
	originalRefNode := alias.GetNode()
	nodeToReplace := aliasRef.GetNode()
	jutil.CheckStateMsg1(nodeToReplace.IsQualifiedName(), "%s", nodeToReplace)
	// If the reference node is a NAME it could be
	// const origName = value;
	// If we use cloneTree() for that we'll clone the value, which we don't want.
	// Otherwise, we do want to clone the tree of GETPROP nodes.
	var newNode *rhino.Node
	if originalRefNode.IsName() {
		newNode = originalRefNode.CloneNode()
	} else {
		newNode = originalRefNode.CloneTree()
	}
	newNode.SrcrefTree(nodeToReplace)
	nodeToReplace.ReplaceWith(newNode)
	a.outer.compiler.ReportChangeToEnclosingScope(newNode)
	return NewGlobalNamespace_SimpleAstChange(newNode, alias.GetChunk(), alias.scope)
}

// InlineGlobalAliasIfPossible attempts to inline an global alias of a
// global name. This requires that the name is well defined: assigned
// unconditionally, assigned exactly once. It is assumed that, the name for
// which it is an alias must already meet these same requirements.
//
// If the alias is completely removed, also deletes the aliasing Ref.
//
// name is the global name being aliased; alias the alias to inline.
//
// Java: InlineAndCollapseProperties.AggressiveInlineAliases.inlineGlobalAliasIfPossible(Name, Ref, GlobalNamespace)
func (a *InlineAndCollapseProperties_AggressiveInlineAliases) InlineGlobalAliasIfPossible(name *GlobalNamespace_Name, alias *GlobalNamespace_Ref, namespace *GlobalNamespace) {
	// Ensure that the alias is assigned to global name at that the
	// declaration.
	aliasParent := alias.GetNode().GetParent()
	if ((aliasParent.IsAssign() || aliasParent.IsName()) &&
		NodeUtil_IsExecutedExactlyOnce(aliasParent)) ||
		// We special-case for constructors here, to inline constructor aliases
		// more aggressively in global scope.
		// We do this because constructor properties are always collapsed,
		// so we want to inline the aliases also to avoid breakages.
		(aliasParent.IsName() && name.IsConstructor()) {
		var lvalue *rhino.Node
		if aliasParent.IsName() {
			lvalue = aliasParent
		} else {
			lvalue = aliasParent.GetFirstChild()
		}
		lvalueName, ok := namespace.qnameStringCached(lvalue)
		if !ok {
			return
		}
		if lvalue.IsName() &&
			a.outer.compiler.GetCodingConvention().IsExportedWithLocal(lvalueName /* local= */, false) {
			return
		}
		aliasingName, _ := namespace.GetSlot(lvalueName).(*GlobalNamespace_Name)

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
		newNodes := a.RewriteAliasReferences(aliasingName, alias, aliasInlinability)
		a.RewriteAliasProps(aliasingName, alias.GetNode(), 0, newNodes)

		if aliasInlinability.ShouldRemoveDeclaration() {
			// Rewrite the initialization of the alias.
			aliasInitialization := aliasingName.GetInitialization()
			if aliasInitialization.IsTwin() {
				// This is in a nested assign.
				// Replace
				//   a.b = aliasing.name = aliased.name
				// with
				//   a.b = aliased.name
				jutil.CheckStateMsg1(aliasParent.IsAssign(), "%s", aliasParent)
				aliasGrandparent := aliasParent.GetParent()
				aliasParent.ReplaceWith(alias.GetNode().Detach())
				// Remove the ref to 'aliasing.name' entirely
				aliasingName.RemoveRef(aliasInitialization)
				// Force GlobalNamespace to revisit the new reference to 'aliased.name' and update its
				// internal state.
				newNodes.Add(NewGlobalNamespace_RefBasedAstChange(alias, alias.GetNode()))
				a.outer.compiler.ReportChangeToEnclosingScope(aliasGrandparent)
			} else {
				// Replace
				//  aliasing.name = aliased.name
				// with
				//  aliasing.name = null;
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

// RewriteAliasReferences replaces reads of a name with the name it aliases
//
// aliasingName is the alias - e.g. 'alsoFoo' in 'const alsoFoo = foo;';
// aliasInlinability indicates whether all references to the alias should be
// inlined, or only those that are in a valid cross-chunk dependency.
//
// Java: InlineAndCollapseProperties.AggressiveInlineAliases.rewriteAliasReferences(Name, Ref, Inlinability)
func (a *InlineAndCollapseProperties_AggressiveInlineAliases) RewriteAliasReferences(
	aliasingName *GlobalNamespace_Name, aliasingRef *GlobalNamespace_Ref, aliasInlinability GlobalNamespace_Inlinability) *jutil.LinkedHashSet[GlobalNamespace_AstChange] {
	newNodes := jutil.NewLinkedHashSet[GlobalNamespace_AstChange]()
	refs := aliasingName.GetRefs()
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
			DestructuringGlobalNameExtractor_ReassignDestructringLvalue(
				ref.GetNode(), aliasingRef.GetNode().CloneTree(), newNodes, ref, a.outer.compiler)
		} else {
			// e.g. `x.y`
			jutil.CheckState(ref.GetNode().IsGetProp() || ref.GetNode().IsName())
			newNode := aliasingRef.GetNode().CloneTree()
			node := ref.GetNode()
			newNode.Srcref(node)
			node.ReplaceWith(newNode)
			a.outer.compiler.ReportChangeToEnclosingScope(newNode)
			newNodes.Add(NewGlobalNamespace_RefBasedAstChange(ref, newNode))
		}
		aliasingName.RemoveRef(ref)
	}
	return newNodes
}

// RewriteAliasProps: name is the Name whose properties references should be
// updated; value the value to use when rewriting; depth the chain depth;
// newNodes expression nodes that have been updated.
//
// Java: InlineAndCollapseProperties.AggressiveInlineAliases.rewriteAliasProps(Name, Node, int, Set<AstChange>)
func (a *InlineAndCollapseProperties_AggressiveInlineAliases) RewriteAliasProps(
	name *GlobalNamespace_Name, value *rhino.Node, depth int, newNodes *jutil.LinkedHashSet[GlobalNamespace_AstChange]) {
	if name.props == nil {
		return
	}
	jutil.CheckStateMsg2(
		!value.MatchesQualifiedName(name.GetFullName()),
		"%s should not match name %s",
		value,
		name.GetFullName())
	for prop := range name.props.All() {
		a.RewriteNestedAliasReference(value, depth, newNodes, prop)
	}
}

// RewriteNestedAliasReference replaces references to an alias that are
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
func (a *InlineAndCollapseProperties_AggressiveInlineAliases) RewriteNestedAliasReference(
	value *rhino.Node, depth int, newNodes *jutil.LinkedHashSet[GlobalNamespace_AstChange], prop *GlobalNamespace_Name) {
	a.RewriteAliasProps(prop, value, depth+1, newNodes)
	refs := prop.GetRefs()
	for _, ref := range refs {
		target := ref.GetNode()
		if target.IsStringKey() && target.GetParent().IsDestructuringPattern() {
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
			jutil.CheckStateMsg1(
				target.GetGrandparent().IsAssign() || target.GetGrandparent().IsDestructuringLhs(),
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
			if target.IsGetProp() {
				target = target.GetFirstChild()
			} else if NodeUtil_IsObjectLitKey(target) {
				// Object literal key definitions are a little trickier, as we
				// need to find the assignment target
				gparent := target.GetGrandparent()
				if gparent.IsAssign() {
					target = gparent.GetFirstChild()
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
		target.ReplaceWith(newValue)
		a.outer.compiler.ReportChangeToEnclosingScope(newValue)
		prop.RemoveRef(ref)
		// Rescan the expression root.
		newNodes.Add(NewGlobalNamespace_RefBasedAstChange(ref, ref.GetNode()))
		a.codeChanged = true
	}
}

// InlineAndCollapseProperties_RewriteSimpleDestructuringAliases rewrites
// "simple" destructuring aliases to a format that is more amenable to
// inlining.
//
// To be specific, this rewrites aliases of the form: const {x} =
// qualified.name; to: const x = qualified.name.x;
type InlineAndCollapseProperties_RewriteSimpleDestructuringAliases struct {
	NodeTraversal_AbstractPostOrderCallbackBase
}

// Java: InlineAndCollapseProperties.RewriteSimpleDestructuringAliases()
func NewInlineAndCollapseProperties_RewriteSimpleDestructuringAliases() *InlineAndCollapseProperties_RewriteSimpleDestructuringAliases {
	return &InlineAndCollapseProperties_RewriteSimpleDestructuringAliases{}
}

// Java: InlineAndCollapseProperties.RewriteSimpleDestructuringAliases.isSimpleDestructuringAlias(Node)
func InlineAndCollapseProperties_RewriteSimpleDestructuringAliases_IsSimpleDestructuringAlias(n *rhino.Node) bool {
	if !NodeUtil_IsStatement(n) || !n.IsConst() {
		return false
	}
	jutil.CheckState(n.HasOneChild())
	destructuringLhs := n.GetFirstChild()
	if !destructuringLhs.IsDestructuringLhs() {
		return false
	}
	objectPattern := destructuringLhs.GetFirstChild()
	if !objectPattern.IsObjectPattern() {
		return false
	}
	rhs := destructuringLhs.GetLastChild()
	if !rhs.IsQualifiedName() {
		return false
	}
	return InlineAndCollapseProperties_RewriteSimpleDestructuringAliases_IsSimpleDestructuringPattern(objectPattern)
}

// Java: InlineAndCollapseProperties.RewriteSimpleDestructuringAliases.isSimpleDestructuringPattern(Node)
func InlineAndCollapseProperties_RewriteSimpleDestructuringAliases_IsSimpleDestructuringPattern(objectPattern *rhino.Node) bool {
	jutil.CheckArgument(objectPattern.IsObjectPattern())
	for key := objectPattern.GetFirstChild(); key != nil; key = key.GetNext() {
		if !key.IsStringKey() || key.IsQuotedStringKey() {
			return false
		}
		jutil.CheckState(key.HasOneChild())
		rhs := key.GetFirstChild()
		if !rhs.IsObjectPattern() && !rhs.IsName() {
			return false
		}
		if rhs.IsObjectPattern() && !InlineAndCollapseProperties_RewriteSimpleDestructuringAliases_IsSimpleDestructuringPattern(rhs) {
			return false
		}
	}
	return true
}

// Java: InlineAndCollapseProperties.RewriteSimpleDestructuringAliases.visit(NodeTraversal, Node, Node)
func (r *InlineAndCollapseProperties_RewriteSimpleDestructuringAliases) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if !InlineAndCollapseProperties_RewriteSimpleDestructuringAliases_IsSimpleDestructuringAlias(n) {
		return
	}

	insertionPoint := n

	destructuringLhs := n.GetFirstChild()
	objectPattern := destructuringLhs.GetFirstChild()
	rhs := destructuringLhs.GetLastChild()
	unusedNewInsertionPoint := r.ExpandObjectPattern(t, insertionPoint, objectPattern, rhs)
	_ = unusedNewInsertionPoint
	n.Detach()
	t.ReportCodeChange()
}

// Java: InlineAndCollapseProperties.RewriteSimpleDestructuringAliases.expandObjectPattern(NodeTraversal, Node, Node, Node)
func (r *InlineAndCollapseProperties_RewriteSimpleDestructuringAliases) ExpandObjectPattern(
	t *NodeTraversal, insertionPoint *rhino.Node, objectPattern *rhino.Node, rhs *rhino.Node) *rhino.Node {
	for key := objectPattern.GetFirstChild(); key != nil; key = key.GetNext() {
		keyChild := key.GetFirstChild()
		var nameNode *rhino.Node
		if keyChild.IsName() {
			nameNode = keyChild.Detach()
		} else {
			jutil.CheckState(keyChild.IsObjectPattern())
			uniqueId := t.GetCompiler().GetUniqueIdSupplier().GetUniqueId(t.GetInput())
			nameNode = rhino.IR_Name("destructuring$" + uniqueId).Srcref(keyChild)
			nameNode.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
		}
		newRhs := rhino.IR_Getprop(rhs.CloneTree(), key.GetString()).Srcref(keyChild)
		newConstNode := rhino.IR_ConstNode(nameNode, newRhs).Srcref(objectPattern)
		newConstNode.InsertAfter(insertionPoint)
		insertionPoint = newConstNode
		if keyChild.IsObjectPattern() {
			insertionPoint = r.ExpandObjectPattern(t, insertionPoint, keyChild, nameNode.CloneNode())
		}
	}
	return insertionPoint
}

// Java: InlineAndCollapseProperties.maybeGetInnerNameNode(Node)
func InlineAndCollapseProperties_MaybeGetInnerNameNode(maybeFunctionOrClassNode *rhino.Node) *rhino.Node {
	if NodeUtil_IsFunctionExpression(maybeFunctionOrClassNode) {
		nameNode := maybeFunctionOrClassNode.GetFirstChild()
		jutil.CheckStateMsg1(nameNode.IsName(), "%s", nameNode)
		// functions with no name have a NAME node with an empty string
		if nameNode.GetString() == "" {
			return nil
		}
		return nameNode
	} else if NodeUtil_IsClassExpression(maybeFunctionOrClassNode) {
		nameNode := maybeFunctionOrClassNode.GetFirstChild()
		// classes with no name have an EMPTY node first child
		if nameNode.IsName() {
			return nameNode
		}
		return nil
	} else {
		return nil // not a function or class expression
	}
}

// InlineAndCollapseProperties_MaybeAddPropertiesToWorklist adds properties
// of `name` to the worklist if the following conditions hold:
//
//  1. The given property of `name` either meets condition (a) or is
//     unsafely collapsible (as defined by Name#canCollapse()
//  2. `name` meets condition (b)
//
// This only adds direct properties of a name, not all its descendants. For
// example, this adds `a.b` given `a`, but not `a.b.c`.
//
// Go-only: language-adapter — the Deque<Name> is a pointer to the slice.
//
// Java: InlineAndCollapseProperties.maybeAddPropertiesToWorklist(Name, Deque<Name>)
func InlineAndCollapseProperties_MaybeAddPropertiesToWorklist(name *GlobalNamespace_Name, workList *[]*GlobalNamespace_Name) {
	if !(name.IsObjectLiteral() || name.IsFunction() || name.IsClass()) {
		// Don't add properties for things like `Foo` in
		//   const Foo = someMysteriousFunctionCall();
		// Since `Foo` is not declared as an object, class, or function literal, assume its value
		// may be aliased somewhere and its properties do not meet condition (a).
		return
	}
	if InlineAndCollapseProperties_IsUnsafelyReassigned(name) {
		// Don't add properties if this was assigned multiple times, except for 'safe'
		// reassignments:
		//    var ns = ns || {};
		// This is equivalent to condition (b)
		return
	}
	if name.props == nil {
		return
	}

	if name.GetAliasingGets() == 0 {
		// All of `name`'s children meet condition (a), so they can be
		// added to the worklist.
		*workList = append(*workList, name.props.Slice()...)
	} else {
		// The children do NOT meet condition (a) but we may try to add them anyway.
		// This is because CollapseProperties will unsafely collapse properties on constructors and
		// enums, so we want to be more aggressive about inlining references to their children.
		for property := range name.props.All() {
			// Only add properties that would be unsafely collapsed by CollapseProperties
			if property.CanCollapse() {
				*workList = append(*workList, property)
			}
		}
	}
}

func inlineAndCollapseProperties_getProperty(
	parent *GlobalNamespace_Name,
	propName string,
	nameMap *jutil.LinkedHashMap[string, *GlobalNamespace_Name],
) *GlobalNamespace_Name {
	if parent.props == nil {
		return nil
	}
	props := parent.props.Slice()
	if len(props) <= 16 {
		for _, p := range props {
			if p.baseName == propName {
				return p
			}
		}
		return nil
	}
	if parent.outer != nil {
		return nameMap.Get(parent.outer.concatDotCached(parent.GetFullName(), propName))
	}
	return nameMap.Get(parent.GetFullName() + "." + propName)
}

// InlineAndCollapseProperties_MayBeGlobalAlias returns true if the alias is
// possibly defined in the global scope, which we handle with more caution
// than with locally scoped variables. May return false positives.
//
// alias is an aliasing get. Returns if the alias is possibly defined in the
// global scope.
//
// Java: InlineAndCollapseProperties.mayBeGlobalAlias(Ref)
func InlineAndCollapseProperties_MayBeGlobalAlias(alias *GlobalNamespace_Ref) bool {
	// Note: alias.scope is the closest scope in which the aliasing assignment occurred.
	// So for "if (true) { var alias = aliasedVar; }", the alias.scope would be the IF block
	// scope.
	if alias.scope.IsGlobal() {
		return true
	}
	// If the scope in which the alias is assigned is not global, look up the LHS of the
	// assignment.
	aliasParent := alias.GetNode().GetParent()
	if !aliasParent.IsAssign() && !aliasParent.IsName() {
		// Only handle variable assignments and initializing declarations.
		return true
	}
	var aliasLhsNode *rhino.Node
	if aliasParent.IsName() {
		aliasLhsNode = aliasParent
	} else {
		aliasLhsNode = aliasParent.GetFirstChild()
	}
	if !aliasLhsNode.IsName() {
		// Only handle assignments to simple names, not qualified names or GETPROPs.
		return true
	}
	aliasVarName := aliasLhsNode.GetString()
	aliasVar := alias.scope.GetVar(aliasVarName)
	if aliasVar != nil {
		return aliasVar.IsGlobal()
	}
	return true
}

// InlineAndCollapseProperties_ReferencesCollapsibleProperty returns whether
// a ReferenceCollection for some aliasing variable references a property on
// the original aliased variable that may be collapsed in
// CollapseProperties.
//
// See Name#canCollapse for what can/cannot be collapsed.
//
// Java: InlineAndCollapseProperties.referencesCollapsibleProperty(ReferenceCollection, Name, GlobalNamespace)
func InlineAndCollapseProperties_ReferencesCollapsibleProperty(
	aliasRefs *ReferenceCollection, aliasedName *GlobalNamespace_Name, namespace *GlobalNamespace) bool {
	for ref := range aliasRefs.references.All() {
		if ref.GetParent() == nil {
			continue
		}
		if NodeUtil_IsNormalOrOptChainGetProp(ref.GetParent()) {
			// e.g. if the reference is "alias.b.someProp", this will be "b".
			propertyName := ref.GetParent().GetString()
			// e.g. if the aliased name is "originalName", this will be "originalName.b".
			originalProperty := inlineAndCollapseProperties_getProperty(aliasedName, propertyName, namespace.GetNameIndex())
			// If the original property isn't in the namespace or can't be collapsed, keep going.
			if originalProperty == nil || !originalProperty.CanCollapse() {
				continue
			}
			return true
		}
	}
	return false
}

// InlineAndCollapseProperties_IsUnsafelyReassigned checks if the name has
// multiple sets that are not of the form "a = a || {}".
//
// Java: InlineAndCollapseProperties.isUnsafelyReassigned(Name)
func InlineAndCollapseProperties_IsUnsafelyReassigned(name *GlobalNamespace_Name) bool {
	if name.GetTotalSets() <= 1 {
		return false
	}
	foundOriginalDefinition := false
	for _, ref := range name.GetRefs() {
		if !ref.IsSet() {
			continue
		}
		if InlineAndCollapseProperties_IsSafeNamespaceReinit(ref) {
			continue
		}
		if !foundOriginalDefinition {
			foundOriginalDefinition = true
		} else {
			return true
		}
	}
	return false
}

// InlineAndCollapseProperties_GetSubclassForEs6Superclass tries to find an
// lvalue for the subclass given the superclass node in an `class ... extends
// ` clause
//
// Only handles cases where we have either a class declaration or a class
// expression in an assignment or name declaration. Otherwise returns null.
//
// Java: InlineAndCollapseProperties.getSubclassForEs6Superclass(Node)
func InlineAndCollapseProperties_GetSubclassForEs6Superclass(superclass *rhino.Node) *rhino.Node {
	classNode := superclass.GetParent()
	jutil.CheckArgumentMsg1(classNode.IsClass(), "%s", classNode)
	if NodeUtil_IsNameDeclaration(classNode.GetGrandparent()) {
		// const Clazz = class extends Super {
		return classNode.GetParent()
	} else if superclass.GetGrandparent().IsAssign() {
		// ns.foo.Clazz = class extends Super {
		return classNode.GetPrevious()
	} else if NodeUtil_IsClassDeclaration(classNode) {
		// class Clazz extends Super {
		return classNode.GetFirstChild()
	}
	return nil
}

// InlineAndCollapseProperties_CollapseProperties flattens global
// objects/namespaces by replacing each '.' with '$' in their names.
//
// This reduces the number of property lookups the browser has to do and
// allows the RenameVars pass to shorten namespaced names. For example,
// goog.events.handleEvent() -> goog$events$handleEvent() -> Za().
//
// If a global object's name is assigned to more than once, or if a property
// is added to the global object in a complex expression, then none of its
// properties will be collapsed (for safety/correctness).
//
// If, after a global object is declared, it is never referenced except when
// its properties are read or set, then the object will be removed after its
// properties have been collapsed.
//
// Uninitialized variable stubs are created at a global object's declaration
// site for any of its properties that are added late in a local scope.
//
// Static properties of constructors are always collapsed, unsafely! For
// other objects: if, after an object is declared, it is referenced directly
// in a way that might create an alias for it, then none of its properties
// will be collapsed. This behavior is a safeguard to prevent the values
// associated with the flattened names from getting out of sync with the
// object's actual property values. For example, in the following case, an
// alias a$b, if created, could easily keep the value 0 even after a.b
// became 5: `a = {b: 0}; c = a; c.b = 5;`.
//
// This pass may break code, but relies on AggressiveInlineAliases running
// before this pass to make some common patterns safer.
//
// This pass doesn't flatten property accesses of the form: a[b].
//
// For lots of examples, see the unit test.
type InlineAndCollapseProperties_CollapseProperties struct {
	outer *InlineAndCollapseProperties

	// Maps names (e.g. "a.b.c") to nodes in the global namespace tree
	nameMap *jutil.LinkedHashMap[string, *GlobalNamespace_Name]

	dynamicallyImportedModules *jutil.LinkedHashSet[string]
}

// Java: InlineAndCollapseProperties.CollapseProperties()
func NewInlineAndCollapseProperties_CollapseProperties(outer *InlineAndCollapseProperties) *InlineAndCollapseProperties_CollapseProperties {
	return &InlineAndCollapseProperties_CollapseProperties{
		outer:                      outer,
		dynamicallyImportedModules: jutil.NewLinkedHashSet[string](),
	}
}

// Java: InlineAndCollapseProperties.CollapseProperties.process(Node, Node)
func (c *InlineAndCollapseProperties_CollapseProperties) Process(externs *rhino.Node, root *rhino.Node) {
	if c.outer.propertyCollapseLevel == CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT ||
		c.outer.chunkOutputType == CompilerOptions_ChunkOutputType_ES_MODULES {
		NodeTraversal_Traverse(
			c.outer.compiler,
			root,
			NewInlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules(c, c.outer.haveModulesBeenRewritten, c.outer.moduleResolutionMode))
	}

	c.nameMap = jutil.CheckNotNullMsg(c.outer.namespace, "namespace was not initialized").GetNameIndex()
	globalNames := c.outer.namespace.GetNameForest()
	escaped := c.CheckNamespaces()
	for name := range globalNames.All() {
		c.FlattenReferencesToCollapsibleDescendantNames(name, name.GetBaseName(), escaped)
		// We collapse property definitions after collapsing property references
		// because this step can alter the parse tree above property references,
		// invalidating the node ancestry stored with each reference.
		c.CollapseDeclarationOfNameAndDescendants(name, name.GetBaseName(), escaped)
	}
}

// Java: InlineAndCollapseProperties.CollapseProperties.canCollapse(Name)
func (c *InlineAndCollapseProperties_CollapseProperties) CanCollapse(name *GlobalNamespace_Name) bool {
	inlinability := name.CanCollapseOrInline()
	if !inlinability.CanCollapse() {
		c.LogDecisionForNameWithInlinability(name, inlinability, "canCollapse() returns false")
		return false
	}

	if c.outer.propertyCollapseLevel == CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT {
		if !name.IsModuleExport() {
			c.LogDecisionForNameWithInlinability(name, inlinability, "module export: canCollapse() returns false")
			return false
		} else if c.dynamicallyImportedModules.Contains(name.GetBaseName()) {
			c.LogDecisionForNameWithInlinability(
				name, inlinability, "dynamic module export: canCollapse() returns false")
			return false
		}
	}

	c.LogDecisionForNameWithInlinability(name, inlinability, "canCollapse() returns true")
	return true
}

// CanEliminate returns true if it is safe to completely remove the rvalue
// assigned to the name from the AST.
//
// Java: InlineAndCollapseProperties.CollapseProperties.canEliminate(Name, Node)
func (c *InlineAndCollapseProperties_CollapseProperties) CanEliminate(name *GlobalNamespace_Name, rvalue *rhino.Node) bool {
	if !name.CanEliminate() {
		return false
	}

	if rvalue.IsObjectLit() {
		for child := rvalue.GetFirstChild(); child != nil; child = child.GetNext() {
			switch child.GetToken() {
			case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_COMPUTED_PROP, rhino.Token_OBJECT_SPREAD:
				// We cannot completely eliminate the object literal, because it contains
				// at least one feature that we cannot extract as a collapsed variable name.
				return false
			case rhino.Token_STRING_KEY, rhino.Token_MEMBER_FUNCTION_DEF:
			default:
				panic(jutil.NewIllegalStateExceptionMsg(
					"Unexpected child of OBJECTLIT: " + child.ToStringTree()))
			}
		}
	}
	if name.props == nil ||
		name.props.IsEmpty() ||
		c.outer.propertyCollapseLevel != CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT {
		return true
	}

	return false
}

// CheckNamespaces runs through all namespaces (prefixes of classes and
// enums), and checks if any of them have been used in an unsafe way.
//
// Java: InlineAndCollapseProperties.CollapseProperties.checkNamespaces()
func (c *InlineAndCollapseProperties_CollapseProperties) CheckNamespaces() *jutil.LinkedHashSet[*GlobalNamespace_Name] {
	escaped := jutil.NewLinkedHashSet[*GlobalNamespace_Name]()
	dynamicallyImportedModuleRefs := jutil.NewLinkedHashSet(c.dynamicallyImportedModules.ToSlice()...)
	if !c.dynamicallyImportedModules.IsEmpty() {
		// When the output chunk type is ES_MODULES, properties of the module namespace
		// must not be collapsed as they are referenced off the namespace object. The namespace
		// objects escape via the dynamic import expression.
		for name := range c.nameMap.Values() {
			// Test if the name is a rewritten module namespace variable and if so mark it as escaped
			// to prevent any property collapsing.
			//
			// example:
			// /** @const */ var module$foo = {};
			if c.dynamicallyImportedModules.Contains(name.GetFullName()) {
				c.LogDecisionForName(name, "escapes - dynamically imported module namespace")
				escaped.Add(name)
				if name.props == nil {
					continue
				}
				for prop := range name.props.All() {
					propDeclaration := prop.getDeclarationRef()
					if propDeclaration == nil {
						continue
					}
					if propDeclaration.GetNode() != nil {
						// ES Module rewriting creates aliases on the module namespace object. These aliased
						// names also escape and their properties may not be collapsed.
						//
						// example:
						// class Foo$$module$foo {
						//   static bar() { return 'bar'; }
						// }
						// /** @const */ var module$foo = {};
						// /** @const */ module$foo.Foo = Foo$$module$foo;
						//
						// While module$foo.Foo cannot be collapsed because we marked the module namespace
						// as escaped, we also need to prevent any property collapsing on the
						// Foo$$module$foo
						// class itself.
						rValue := NodeUtil_GetRValueOfLValue(propDeclaration.GetNode())
						if rValue.IsName() {
							c.LogDecisionForName(
								name, "escapes - dynamically imported module namespace property alias")
							dynamicallyImportedModuleRefs.Add(*rValue.GetQualifiedName())
						}
					}
				}
			}
		}
	}

	for name := range c.nameMap.Values() {
		if !dynamicallyImportedModuleRefs.IsEmpty() &&
			dynamicallyImportedModuleRefs.Contains(name.GetFullName()) {
			escaped.Add(name)
		}
		if !name.IsNamespaceObjectLit() {
			continue
		}
		if name.GetAliasingGets() == 0 &&
			name.GetLocalSets()+name.GetGlobalSets() <= 1 &&
			name.GetDeleteProps() == 0 {
			continue
		}
		initialized := name.GetDeclaration() != nil
		for _, ref := range name.GetRefs() {
			if ref.IsDeleteProp() {
				if initialized {
					c.WarnAboutNamespaceRedefinition(name, ref)
				}
			} else if ref.IsSet() && ref != name.getDeclarationRef() {
				if initialized && !InlineAndCollapseProperties_IsSafeNamespaceReinit(ref) {
					c.WarnAboutNamespaceRedefinition(name, ref)
				}

				initialized = true
			} else if ref.IsAliasingGet() {
				c.WarnAboutNamespaceAliasing(name, ref)
				c.LogDecisionForName(name, "escapes")
				escaped.Add(name)
				break
			}
		}
	}
	return escaped
}

// WarnAboutNamespaceAliasing reports a warning because a namespace was
// aliased.
//
// nameObj is a namespace that is being aliased; ref the reference that
// forced the alias.
//
// Java: InlineAndCollapseProperties.CollapseProperties.warnAboutNamespaceAliasing(Name, Ref)
func (c *InlineAndCollapseProperties_CollapseProperties) WarnAboutNamespaceAliasing(nameObj *GlobalNamespace_Name, ref *GlobalNamespace_Ref) {
	c.outer.compiler.Report(
		JSError_MakeWithNode(ref.GetNode(), InlineAndCollapseProperties_PARTIAL_NAMESPACE_WARNING, nameObj.GetFullName()))
}

// WarnAboutNamespaceRedefinition reports a warning because a namespace was
// redefined.
//
// nameObj is a namespace that is being redefined; ref the reference that set
// the namespace.
//
// Java: InlineAndCollapseProperties.CollapseProperties.warnAboutNamespaceRedefinition(Name, Ref)
func (c *InlineAndCollapseProperties_CollapseProperties) WarnAboutNamespaceRedefinition(nameObj *GlobalNamespace_Name, ref *GlobalNamespace_Ref) {
	c.outer.compiler.Report(
		JSError_MakeWithNode(ref.GetNode(), InlineAndCollapseProperties_NAMESPACE_REDEFINED_WARNING, nameObj.GetFullName()))
}

// FlattenReferencesToCollapsibleDescendantNames flattens all references to
// collapsible properties of a global name except their initial definitions.
// Recurs on subnames.
//
// n is an object representing a global name; alias the flattened name for
// n.
//
// Java: InlineAndCollapseProperties.CollapseProperties.flattenReferencesToCollapsibleDescendantNames(Name, String, Set<Name>)
func (c *InlineAndCollapseProperties_CollapseProperties) FlattenReferencesToCollapsibleDescendantNames(
	n *GlobalNamespace_Name, alias string, escaped *jutil.LinkedHashSet[*GlobalNamespace_Name]) {
	if n.props == nil {
		return
	}
	if n.IsCollapsingExplicitlyDenied() {
		c.LogDecisionForName(n, "@nocollapse: will not flatten descendant name references")
		return
	}
	if escaped.Contains(n) {
		c.LogDecisionForName(n, "escapes: will not flatten descendant name references")
		return
	}

	for p := range n.props.All() {
		propAlias := c.AppendPropForAlias(alias, p.GetBaseName())
		inlinability := p.CanCollapseOrInline()

		isAllowedToCollapse :=
			c.outer.propertyCollapseLevel != CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT || p.IsModuleExport()
		if isAllowedToCollapse {
			if inlinability.CanCollapse() {
				c.LogDecisionForNameWithInlinability(p, inlinability, "will flatten references")
				c.flattenReferencesTo(p, propAlias, escaped)
			} else if p.IsCollapsingExplicitlyDenied() {
				c.LogDecisionForName(p, "@nocollapse: will not flatten references")
			} else if p.IsSimpleStubDeclaration() {
				c.LogDecisionForName(p, "simple stub declaration: will flatten references")
				c.FlattenSimpleStubDeclaration(p, propAlias)
			} else {
				c.LogDecisionForNameWithInlinability(p, inlinability, "will not flatten references")
			}
		}

		c.FlattenReferencesToCollapsibleDescendantNames(p, propAlias, escaped)
	}
}

// Java: InlineAndCollapseProperties.CollapseProperties.logDecisionForName(Name, Inlinability, String)
func (c *InlineAndCollapseProperties_CollapseProperties) LogDecisionForNameWithInlinability(
	name *GlobalNamespace_Name, inlinability GlobalNamespace_Inlinability, message string) {
	if c.outer.decisionsLog.IsLogging() {
		c.LogDecisionForNameSupplier(name, func() string { return "inlinability " + inlinability.String() + ": " + message })
	}
}

// Java: InlineAndCollapseProperties.CollapseProperties.logDecisionForName(Name, String)
func (c *InlineAndCollapseProperties_CollapseProperties) LogDecisionForName(name *GlobalNamespace_Name, message string) {
	if c.outer.decisionsLog.IsLogging() {
		c.outer.decisionsLog.LogSupplier(func() string { return name.GetFullName() + ": " + message })
	}
}

// Java: InlineAndCollapseProperties.CollapseProperties.logDecisionForName(Name, Supplier<String>)
func (c *InlineAndCollapseProperties_CollapseProperties) LogDecisionForNameSupplier(name *GlobalNamespace_Name, messageSupplier func() string) {
	if c.outer.decisionsLog.IsLogging() {
		c.outer.decisionsLog.LogSupplier(func() string { return name.GetFullName() + ": " + messageSupplier() })
	}
}

// FlattenSimpleStubDeclaration flattens a stub declaration. This is mostly a
// hack to support legacy users.
//
// Java: InlineAndCollapseProperties.CollapseProperties.flattenSimpleStubDeclaration(Name, String)
func (c *InlineAndCollapseProperties_CollapseProperties) FlattenSimpleStubDeclaration(name *GlobalNamespace_Name, alias string) {
	refs := name.GetRefs()
	// Iterables.getOnlyElement(name.getRefs())
	if len(refs) == 0 {
		panic(jutil.NewNoSuchElementException())
	} else if len(refs) > 1 {
		panic(jutil.NewIllegalArgumentExceptionMsg(inlineAndCollapseProperties_onlyElementMessage(refs)))
	}
	ref := refs[0]
	nameNode := NodeUtil_NewNameWithOriginalName(c.outer.compiler, alias, ref.GetNode(), name.GetFullName())
	varNode := rhino.IR_Var(nameNode).SrcrefIfMissing(nameNode)

	jutil.CheckState(ref.GetNode().GetParent().IsExprResult())
	parent := ref.GetNode().GetParent()
	parent.ReplaceWith(varNode)
	c.outer.compiler.ReportChangeToEnclosingScope(varNode)
}

// inlineAndCollapseProperties_onlyElementMessage is the message of Guava's
// Iterables.getOnlyElement for more than one element: "expected one element
// but was: <e1, e2, ...>" (at most five elements, then ", ...").
//
// Go-only: language-adapter — Guava's message, computed only on failure.
func inlineAndCollapseProperties_onlyElementMessage(refs []*GlobalNamespace_Ref) string {
	var sb strings.Builder
	sb.WriteString("expected one element but was: <")
	for i, ref := range refs {
		if i == 5 {
			sb.WriteString(", ...")
			break
		}
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(ref.String())
	}
	sb.WriteByte('>')
	return sb.String()
}

// ShouldAddConstantName returns whether a NAME node with the given
// collapsedName, for a name originally declared at declarationNode, should
// be annotated with Node#IS_CONSTANT_NAME
//
// Validity checks in unit tests will otherwise throw an exception.
//
// Java: InlineAndCollapseProperties.CollapseProperties.shouldAddConstantName(Node, String)
func (c *InlineAndCollapseProperties_CollapseProperties) ShouldAddConstantName(declarationNode *rhino.Node, collapsedName string) bool {
	if declarationNode == nil {
		return false
	}
	info := NodeUtil_GetBestJSDocInfo(declarationNode.GetParent())
	return (info != nil && info.IsConstant()) ||
		declarationNode.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME) ||
		c.outer.compiler.GetCodingConvention().IsConstant(collapsedName)
}

func (c *InlineAndCollapseProperties_CollapseProperties) willFlattenReferences(p *GlobalNamespace_Name) bool {
	isAllowedToCollapse :=
		c.outer.propertyCollapseLevel != CompilerOptions_PropertyCollapseLevel_MODULE_EXPORT || p.IsModuleExport()
	return isAllowedToCollapse && p.CanCollapseOrInline().CanCollapse()
}

// FlattenReferencesTo flattens all references to a collapsible property of a
// global name except its initial definition.
//
// n is a global property name (e.g. "a.b" or "a.b.c.d"); alias the
// flattened name (e.g. "a$b" or "a$b$c$d").
//
// Java: InlineAndCollapseProperties.CollapseProperties.flattenReferencesTo(Name, String)
func (c *InlineAndCollapseProperties_CollapseProperties) FlattenReferencesTo(n *GlobalNamespace_Name, alias string) {
	c.flattenReferencesTo(n, alias, nil)
}

func (c *InlineAndCollapseProperties_CollapseProperties) flattenReferencesTo(
	n *GlobalNamespace_Name, alias string, escaped *jutil.LinkedHashSet[*GlobalNamespace_Name]) {
	originalName := n.GetFullName()
	var declarationNode *rhino.Node
	if n.GetDeclaration() != nil {
		declarationNode = n.getDeclarationRef().GetNode()
	}
	isConstantName := c.ShouldAddConstantName(declarationNode, alias)
	for _, r := range n.GetRefs() {
		if r == n.getDeclarationRef() {
			// Declarations are handled separately.
			continue
		}
		rParent := r.GetNode().GetParent()
		// We shouldn't flatten a reference that's an object literal key, because duplicate keys
		// show up as refs.
		if !NodeUtil_MayBeObjectLitKey(r.GetNode()) {
			c.FlattenNameRef(alias, r.GetNode(), rParent, originalName, isConstantName)
		} else if r.GetNode().IsStringKey() && r.GetNode().GetParent().IsObjectPattern() {
			newNode := rhino.IR_Name(alias).Srcref(r.GetNode())
			NodeUtil_CopyNameAnnotations(r.GetNode(), newNode)
			DestructuringGlobalNameExtractor_ReassignDestructringLvalue(
				r.GetNode(), newNode, nil, r, c.outer.compiler)
		}
	}

	// Flatten all occurrences of a name as a prefix of its subnames. For
	// example, if `n` corresponds to the name "a.b", then "a.b" will be
	// replaced with "a$b" in all occurrences of "a.b.c", "a.b.c.d", etc.
	if n.props != nil {
		descendantsVisited := escaped != nil && !n.IsCollapsingExplicitlyDenied() && !escaped.Contains(n)
		for p := range n.props.All() {
			if flattenPrefixesHasWork(p) {
				c.flattenPrefixes(alias, originalName+p.GetBaseName(), p, 1, isConstantName, escaped, descendantsVisited)
			}
		}
	}
}

func flattenPrefixesHasWork(n *GlobalNamespace_Name) bool {
	if n.props != nil {
		return true
	}
	decl := n.getDeclarationRef()
	if decl != nil && decl.GetNode() != nil && decl.GetNode().IsGetProp() {
		return true
	}
	refs := n.GetRefs()
	return len(refs) > 1 || (len(refs) == 1 && refs[0] != decl)
}

// FlattenPrefixes flattens all occurrences of a name as a prefix of subnames
// beginning with a particular subname.
//
// n is a global property name (e.g. "a.b.c.d"); alias a flattened prefix
// name (e.g. "a$b"); originalName the full original name of the global
// property (e.g. "a.b.c.d") equivalent to n.getFullName(), but pre-computed
// to save on intermediate string allocation; depth the difference in depth
// between the property name and the prefix name (e.g. 2).
//
// Java: InlineAndCollapseProperties.CollapseProperties.flattenPrefixes(String, String, Name, int, boolean)
func (c *InlineAndCollapseProperties_CollapseProperties) FlattenPrefixes(
	alias string, originalName string, n *GlobalNamespace_Name, depth int, isConstantName bool) {
	c.flattenPrefixes(alias, originalName, n, depth, isConstantName, nil, false)
}

func (c *InlineAndCollapseProperties_CollapseProperties) flattenPrefixes(
	alias string,
	originalName string,
	n *GlobalNamespace_Name,
	depth int,
	isConstantName bool,
	escaped *jutil.LinkedHashSet[*GlobalNamespace_Name],
	descendantsVisited bool) {
	// Only flatten the prefix of a name declaration if the name being
	// initialized is fully qualified (i.e. not an object literal key).
	decl := n.getDeclarationRef()
	if decl != nil && decl.GetNode() != nil && decl.GetNode().IsGetProp() {
		c.FlattenNameRefAtDepth(alias, decl.GetNode(), depth, originalName, isConstantName)
	}

	// If FlattenReferencesToCollapsibleDescendantNames will visit n and call
	// flattenReferencesTo(n, ...), then every non-declaration GETPROP ref in
	// n.GetRefs() and every prefix in n.props will be replaced with n's own
	// collapsed alias, completely overwriting any intermediate prefix node
	// created here for n's ancestor.
	if descendantsVisited && c.willFlattenReferences(n) {
		return
	}

	for _, r := range n.GetRefs() {
		if r == decl {
			// Declarations are handled separately.
			continue
		}

		c.FlattenNameRefAtDepth(alias, r.GetNode(), depth, originalName, isConstantName)
	}

	if n.props != nil {
		nextDescendantsVisited := descendantsVisited && !n.IsCollapsingExplicitlyDenied() && !escaped.Contains(n)
		for p := range n.props.All() {
			if flattenPrefixesHasWork(p) {
				c.flattenPrefixes(alias, originalName+p.GetBaseName(), p, depth+1, isConstantName, escaped, nextDescendantsVisited)
			}
		}
	}
}

// FlattenNameRefAtDepth flattens a particular prefix of a single name
// reference.
//
// alias is a flattened prefix name (e.g. "a$b"); n the node corresponding to
// a subproperty name (e.g. "a.b.c.d"); depth the difference in depth between
// the property name and the prefix name (e.g. 2); originalName the string
// version of the property name.
//
// Java: InlineAndCollapseProperties.CollapseProperties.flattenNameRefAtDepth(String, Node, int, String, boolean)
func (c *InlineAndCollapseProperties_CollapseProperties) FlattenNameRefAtDepth(
	alias string, n *rhino.Node, depth int, originalName string, isConstantName bool) {
	// This method has to work for both GETPROP chains and, in rare cases,
	// OBJLIT keys, possibly nested. That's why we check for children before
	// proceeding. In the OBJLIT case, we don't need to do anything.
	nType := n.GetToken()
	isQName := nType == rhino.Token_NAME || nType == rhino.Token_GETPROP
	isObjKey := NodeUtil_MayBeObjectLitKey(n)
	jutil.CheckState(isObjKey || isQName)
	if isQName {
		for i := 1; i < depth && n.HasChildren(); i++ {
			n = n.GetFirstChild()
		}
		if n.IsGetProp() && n.GetFirstChild().IsGetProp() {
			c.FlattenNameRef(alias, n.GetFirstChild(), n, originalName, isConstantName)
		}
	}
}

// FlattenNameRef replaces a GETPROP a.b.c with a NAME a$b$c.
//
// alias is a flattened prefix name (e.g. "a$b"); n the GETPROP node
// corresponding to the original name (e.g. "a.b"); parent n's parent;
// originalName the string version of the property name; isConstantName
// whether to annotate with Node#IS_CONSTANT_NAME.
//
// Java: InlineAndCollapseProperties.CollapseProperties.flattenNameRef(String, Node, Node, String, boolean)
func (c *InlineAndCollapseProperties_CollapseProperties) FlattenNameRef(
	alias string, n *rhino.Node, parent *rhino.Node, originalName string, isConstantName bool) {
	jutil.CheckArgumentMsg2(
		n.IsGetProp(), "Expected GETPROP, found %s. Node: %s", n.GetToken(), n)

	// BEFORE:
	//   getprop
	//     getprop
	//       name a
	//       string b
	//     string c
	// AFTER:
	//   name a$b$c
	ref := NodeUtil_NewNameWithOriginalName(c.outer.compiler, alias, n, originalName).CopyTypeFrom(n)
	NodeUtil_CopyNameAnnotations(n, ref)
	if isConstantName {
		ref.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, isConstantName)
	}
	if NodeUtil_IsNormalOrOptChainCall(parent) && n.IsFirstChildOf(parent) {
		// The node was a call target. We are deliberately flattening these as
		// the "this" isn't provided by the namespace. Mark it as such:
		parent.PutBooleanProp(rhino.Node_FREE_CALL, true)
	}

	n.ReplaceWith(ref)
	c.outer.compiler.ReportChangeToEnclosingScope(ref)
}

// CollapseDeclarationOfNameAndDescendants collapses definitions of the
// collapsible properties of a global name. Recurs on subnames that also
// represent JavaScript objects with collapsible properties.
//
// n is a node representing a global name; alias the flattened name for n.
//
// Java: InlineAndCollapseProperties.CollapseProperties.collapseDeclarationOfNameAndDescendants(Name, String, Set<Name>)
func (c *InlineAndCollapseProperties_CollapseProperties) CollapseDeclarationOfNameAndDescendants(
	n *GlobalNamespace_Name, alias string, escaped *jutil.LinkedHashSet[*GlobalNamespace_Name]) {
	childNameInlinability := n.CanCollapseOrInlineChildNames()
	var canCollapseChildNames bool
	if !childNameInlinability.CanCollapse() {
		c.LogDecisionForNameSupplier(
			n,
			func() string {
				return "child name inlinability: " + childNameInlinability.String() + ": will not collapse child names"
			})
		canCollapseChildNames = false
	} else if escaped.Contains(n) {
		c.LogDecisionForName(n, "escapes: will not collapse child names")
		canCollapseChildNames = false
	} else {
		canCollapseChildNames = true
	}

	// Handle this name first so that nested object literals get unrolled.
	if c.CanCollapse(n) {
		c.LogDecisionForName(n, "collapsing")
		c.UpdateGlobalNameDeclaration(n, alias, canCollapseChildNames)
	}

	if n.props == nil || escaped.Contains(n) {
		return
	}
	c.LogDecisionForName(n, "collapsing descendants")
	for p := range n.props.All() {
		c.CollapseDeclarationOfNameAndDescendants(
			p, c.AppendPropForAlias(alias, p.GetBaseName()), escaped)
	}
}

// UpdateTwinnedDeclaration updates the initial assignment to a collapsible
// property at global scope by adding a VAR stub and collapsing the
// property. e.g. c = a.b = 1; => var a$b; c = a$b = 1; This specifically
// handles "twinned" assignments, which are those where the assignment is
// also used as a reference and which need special handling.
//
// alias is the flattened property name (e.g. "a$b"); refName the name for
// the reference being updated; ref an object containing information about
// the assignment getting updated.
//
// Java: InlineAndCollapseProperties.CollapseProperties.updateTwinnedDeclaration(String, Name, Ref)
func (c *InlineAndCollapseProperties_CollapseProperties) UpdateTwinnedDeclaration(alias string, refName *GlobalNamespace_Name, ref *GlobalNamespace_Ref) {
	jutil.CheckStateMsg1(ref.IsTwin(), "%s", ref)
	// Don't handle declarations of an already flat name, just qualified names.
	if !ref.GetNode().IsGetProp() {
		return
	}
	rvalue := ref.GetNode().GetNext()
	parent := ref.GetNode().GetParent()
	grandparent := parent.GetParent()

	if rvalue != nil && rvalue.IsFunction() {
		c.CheckForReceiverAffectedByCollapse(rvalue, refName.GetJSDocInfo(), refName)
	}

	// Create the new alias node.
	nameNode := NodeUtil_NewNameWithOriginalName(c.outer.compiler, alias, grandparent.GetFirstChild(), refName.GetFullName())
	NodeUtil_CopyNameAnnotations(ref.GetNode(), nameNode)

	// BEFORE:
	// ... (x.y = 3);
	//
	// AFTER:
	// var x$y;
	// ... (x$y = 3);

	current := grandparent
	currentParent := grandparent.GetParent()
	for ; !currentParent.IsScript() && !currentParent.IsBlock(); current, currentParent = currentParent, currentParent.GetParent() {
	}

	// Create a stub variable declaration right
	// before the current statement.
	stubVar := rhino.IR_Var(nameNode.CloneTree()).SrcrefIfMissing(nameNode)
	stubVar.InsertBefore(current)

	ref.GetNode().ReplaceWith(nameNode)
	c.outer.compiler.ReportChangeToEnclosingScope(nameNode)
}

// UpdateGlobalNameDeclaration updates the first initialization (a.k.a
// "declaration") of a global name. This involves flattening the global name
// (if it's not just a global variable name already), collapsing object
// literal keys into global variables, declaring stub global variables for
// properties added later in a local scope.
//
// It may seem odd that this function also takes care of declaring stubs for
// direct children. The ultimate goal of this function is to eliminate the
// global name entirely (when possible), so that "middlemen" namespaces
// disappear, and to do that we need to make sure that all the direct
// children will be collapsed as well.
//
// n is an object representing a global name (e.g. "a", "a.b.c"); alias the
// flattened name for n (e.g. "a", "a$b$c"); canCollapseChildNames whether
// it's possible to collapse children of this name. (This is mostly passed
// for convenience; it's equivalent to n.canCollapseChildNames()).
//
// Java: InlineAndCollapseProperties.CollapseProperties.updateGlobalNameDeclaration(Name, String, boolean)
func (c *InlineAndCollapseProperties_CollapseProperties) UpdateGlobalNameDeclaration(n *GlobalNamespace_Name, alias string, canCollapseChildNames bool) {
	decl := n.getDeclarationRef()
	if decl == nil {
		// Some names do not have declarations, because they
		// are only defined in local scopes.
		c.LogDecisionForName(n, "no global declaration found")
		return
	}

	isSimpleName := n.IsSimpleName()
	declNode := decl.GetNode()
	switch declNode.GetParent().GetToken() {
	case rhino.Token_ASSIGN:
		c.LogDeclarationAction(n, declNode, "updating assignment")
		c.UpdateGlobalNameDeclarationAtAssignNode(n, alias, canCollapseChildNames)
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		c.LogDeclarationAction(n, declNode, "updating variable declaration")
		c.UpdateGlobalNameDeclarationAtVariableNode(n, canCollapseChildNames)
	case rhino.Token_FUNCTION:
		c.LogDeclarationAction(n, declNode, "updating function declaration")
		c.UpdateGlobalNameDeclarationAtFunctionNode(n, canCollapseChildNames)
	case rhino.Token_CLASS:
		c.LogDeclarationAction(n, declNode, "updating class declaration")
		c.UpdateGlobalNameDeclarationAtClassNode(n, canCollapseChildNames)
	case rhino.Token_CLASS_MEMBERS:
		c.LogDeclarationAction(n, declNode, "updating static member declaration")
		c.UpdateGlobalNameDeclarationAtStaticMemberNode(n, alias, canCollapseChildNames)
	default:
		c.LogDeclarationAction(
			n, declNode, "not updating an unsupported type of declaration node")
	}

	if !isSimpleName {
		commonAncestor := n.GetDeepestCommonAncestorChunk(c.outer.compiler.GetChunkGraph())
		if commonAncestor != decl.GetChunk() {
			c.MoveDeclarationIntoHigherChunk(decl, commonAncestor)
		}
	}
}

// Java: InlineAndCollapseProperties.CollapseProperties.logDeclarationAction(Name, Node, String)
func (c *InlineAndCollapseProperties_CollapseProperties) LogDeclarationAction(name *GlobalNamespace_Name, declarationNode *rhino.Node, message string) {
	c.LogDecisionForNameSupplier(name, func() string { return declarationNode.String() + ": " + message })
}

// Java: InlineAndCollapseProperties.CollapseProperties.moveDeclarationIntoHigherChunk(Ref, JSChunk)
func (c *InlineAndCollapseProperties_CollapseProperties) MoveDeclarationIntoHigherChunk(decl *GlobalNamespace_Ref, commonAncestor *JSChunk) {
	// Var hoisting means that early references are generally ok. E.g. `a$b; var a$b;` is valid.
	// However, if `var a$b;` were in a different, conditionally loaded chunk, then it might not
	// be parsed before `a$b;`. So we need to avoid a ReferenceError when referencing `a$b`.
	newNode := decl.GetNode()
	varNode := newNode.GetParent()
	initializer := newNode.GetFirstChild()
	jutil.CheckStateMsg1(varNode.IsVar(), "Unexpected declaration %s", varNode)

	newRoot := commonAncestor.GetInputs()[0].GetAstRoot(c.outer.compiler)
	if initializer == nil {
		// Already a stub declaration - "var a$b;". Move it to the ancestor chunk.
		newRoot.AddChildToFront(varNode.Detach())
	} else {
		// Define `var a$b;` in the ancestor chunk.
		newVarNode := rhino.IR_Var(newNode.CloneNode()).Srcref(newNode)
		newRoot.AddChildToFront(newVarNode)
		// Replace `var a$b = () => {};` with `a$b = () = {};`. Having multiple `var` declarations
		// of the same name is valid JavaScript but violates Normalize constraints.
		newAssign := rhino.IR_ExprResult(rhino.IR_Assign(newNode.Detach(), initializer.Detach())).
			SrcrefTreeIfMissing(newNode)
		varNode.ReplaceWith(newAssign)
	}
	c.outer.compiler.ReportChangeToChangeScope(newRoot)
}

// UpdateGlobalNameDeclarationAtAssignNode updates the first initialization
// (a.k.a "declaration") of a global name that occurs at an ASSIGN node. See
// comment for updateGlobalNameDeclaration.
//
// n is an object representing a global name (e.g. "a", "a.b.c"); alias the
// flattened name for n (e.g. "a", "a$b$c").
//
// Java: InlineAndCollapseProperties.CollapseProperties.updateGlobalNameDeclarationAtAssignNode(Name, String, boolean)
func (c *InlineAndCollapseProperties_CollapseProperties) UpdateGlobalNameDeclarationAtAssignNode(
	n *GlobalNamespace_Name, alias string, canCollapseChildNames bool) {
	// NOTE: It's important that we don't add additional nodes
	// (e.g. a var node before the exprstmt) because the exprstmt might be
	// the child of an if statement that's not inside a block).

	// All qualified names - even for variables that are initially declared as LETS and CONSTS -
	// are being declared as VAR statements, but this is not incorrect because
	// we are only collapsing for global names.
	ref := n.getDeclarationRef()
	rvalue := ref.GetNode().GetNext()
	if ref.IsTwin() {
		c.UpdateTwinnedDeclaration(alias, n, ref)
		return
	}
	varNode := rhino.NewNode(rhino.Token_VAR)
	varParent := ref.GetNode().GetAncestor(3)
	grandparent := ref.GetNode().GetAncestor(2)
	isObjLit := rvalue.IsObjectLit()
	insertedVarNode := false

	if isObjLit && c.CanEliminate(n, rvalue) {
		// Eliminate the object literal altogether.
		grandparent.ReplaceWith(varNode)
		n.UpdateRefNode(ref, nil)
		insertedVarNode = true
		c.outer.compiler.ReportChangeToEnclosingScope(varNode)
	} else if !n.IsSimpleName() {
		// Create a VAR node to declare the name.
		if rvalue.IsFunction() {
			c.CheckForReceiverAffectedByCollapse(rvalue, n.GetJSDocInfo(), n)
		}

		c.outer.compiler.ReportChangeToEnclosingScope(rvalue)
		rvalue.Detach()

		nameNode := NodeUtil_NewNameWithOriginalName(c.outer.compiler, alias, ref.GetNode().GetAncestor(2), n.GetFullName())

		nameNode.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, c.ShouldAddConstantName(ref.GetNode(), alias))

		info := NodeUtil_GetBestJSDocInfo(ref.GetNode().GetParent())
		if info != nil {
			varNode.SetJSDocInfo(info)
		}
		varNode.AddChildToBack(nameNode)
		nameNode.AddChildToFront(rvalue)
		grandparent.ReplaceWith(varNode)

		// Update the node ancestry stored in the reference.
		n.UpdateRefNode(ref, nameNode)
		insertedVarNode = true
		c.outer.compiler.ReportChangeToEnclosingScope(varNode)
	}

	if canCollapseChildNames {
		if isObjLit {
			c.DeclareVariablesForObjLitValues(n, alias, rvalue, varNode, varNode.GetPrevious())
		}

		c.AddStubsForUndeclaredProperties(n, alias, varParent, varNode)
	}

	if insertedVarNode {
		if !varNode.HasChildren() {
			varNode.Detach()
		}
	}
}

// CheckForReceiverAffectedByCollapse warns about any references to "this" in
// the given FUNCTION. The function is getting collapsed, so the references
// will change.
//
// Java: InlineAndCollapseProperties.CollapseProperties.checkForReceiverAffectedByCollapse(Node, JSDocInfo, Name)
func (c *InlineAndCollapseProperties_CollapseProperties) CheckForReceiverAffectedByCollapse(function *rhino.Node, docInfo *rhino.JSDocInfo, name *GlobalNamespace_Name) {
	jutil.CheckState(function.IsFunction())

	if docInfo != nil {
		// Don't rely on type inference for this check.

		if docInfo.IsConstructorOrInterface() {
			return // Ctors and interfaces need to be able to reference `this`
		}
		if docInfo.HasThisType() {
			/*
			 * Use `@this` as a signal that the reference is intentional.
			 *
			 * <p>TODO(b/156823102): This signal also silences the check on all transpiled static
			 * methods.
			 */
			return
		}
	}

	// Use the NodeUtil method so we don't forget to update this logic.
	if NodeUtil_ReferencesOwnReceiver(function) {
		c.outer.compiler.Report(JSError_MakeWithNode(function, InlineAndCollapseProperties_RECEIVER_AFFECTED_BY_COLLAPSE, name.GetFullName()))
	}
}

// UpdateGlobalNameDeclarationAtVariableNode updates the first
// initialization (a.k.a "declaration") of a global name that occurs at a VAR
// node. See comment for updateGlobalNameDeclaration.
//
// n is an object representing a global name (e.g. "a").
//
// Java: InlineAndCollapseProperties.CollapseProperties.updateGlobalNameDeclarationAtVariableNode(Name, boolean)
func (c *InlineAndCollapseProperties_CollapseProperties) UpdateGlobalNameDeclarationAtVariableNode(n *GlobalNamespace_Name, canCollapseChildNames bool) {
	if !canCollapseChildNames {
		c.LogDecisionForName(n, "cannot collapse child names: skipping")
		return
	}

	ref := n.getDeclarationRef()
	name := ref.GetNode().GetString()
	rvalue := ref.GetNode().GetFirstChild()
	variableNode := ref.GetNode().GetParent()
	grandparent := variableNode.GetParent()

	isObjLit := rvalue.IsObjectLit()

	if isObjLit {
		c.DeclareVariablesForObjLitValues(n, name, rvalue, variableNode, variableNode.GetPrevious())
	}

	c.AddStubsForUndeclaredProperties(n, name, grandparent, variableNode)

	if isObjLit && c.CanEliminate(n, rvalue) {
		ref.GetNode().Detach()
		c.outer.compiler.ReportChangeToEnclosingScope(variableNode)
		if !variableNode.HasChildren() {
			variableNode.Detach()
		}

		// Clear out the object reference, since we've eliminated it from the
		// parse tree.
		n.UpdateRefNode(ref, nil)
	}
}

// UpdateGlobalNameDeclarationAtFunctionNode updates the first
// initialization (a.k.a "declaration") of a global name that occurs at a
// FUNCTION node. See comment for updateGlobalNameDeclaration.
//
// n is an object representing a global name (e.g. "a").
//
// Java: InlineAndCollapseProperties.CollapseProperties.updateGlobalNameDeclarationAtFunctionNode(Name, boolean)
func (c *InlineAndCollapseProperties_CollapseProperties) UpdateGlobalNameDeclarationAtFunctionNode(n *GlobalNamespace_Name, canCollapseChildNames bool) {
	if !canCollapseChildNames || !c.CanCollapse(n) {
		return
	}

	ref := n.getDeclarationRef()
	fnName := ref.GetNode().GetString()
	c.AddStubsForUndeclaredProperties(
		n, fnName, ref.GetNode().GetAncestor(2), ref.GetNode().GetParent())
}

// UpdateGlobalNameDeclarationAtClassNode updates the first initialization
// (a.k.a "declaration") of a global name that occurs at a CLASS node. See
// comment for updateGlobalNameDeclaration.
//
// n is an object representing a global name (e.g. "a").
//
// Java: InlineAndCollapseProperties.CollapseProperties.updateGlobalNameDeclarationAtClassNode(Name, boolean)
func (c *InlineAndCollapseProperties_CollapseProperties) UpdateGlobalNameDeclarationAtClassNode(n *GlobalNamespace_Name, canCollapseChildNames bool) {
	if !canCollapseChildNames || !c.CanCollapse(n) {
		return
	}

	ref := n.getDeclarationRef()
	className := ref.GetNode().GetString()
	c.AddStubsForUndeclaredProperties(
		n, className, ref.GetNode().GetAncestor(2), ref.GetNode().GetParent())
}

// UpdateGlobalNameDeclarationAtStaticMemberNode updates the first
// initialization (a.k.a "declaration") of a global name that occurs in a
// static MEMBER_FUNCTION_DEF or MEMBER_FIELD_DEF in a class. See comment for
// updateGlobalNameDeclaration.
//
// n is a static MEMBER_FUNCTION_DEF or MEMBER_FIELD_DEF in a class assigned
// to a global name (e.g. `a.b`); alias the new flattened name for `n` (e.g.
// "a$b"); canCollapseChildNames whether properties of `n` are also
// collapsible, meaning that any properties only assigned locally need stub
// declarations.
//
// Java: InlineAndCollapseProperties.CollapseProperties.updateGlobalNameDeclarationAtStaticMemberNode(Name, String, boolean)
func (c *InlineAndCollapseProperties_CollapseProperties) UpdateGlobalNameDeclarationAtStaticMemberNode(
	n *GlobalNamespace_Name, alias string, canCollapseChildNames bool) {

	declaration := n.getDeclarationRef()
	classNode := declaration.GetNode().GetGrandparent()
	jutil.CheckStateMsg1(classNode.IsClass(), "%s", classNode)
	enclosingStatement := NodeUtil_GetEnclosingStatement(classNode)

	if canCollapseChildNames {
		c.AddStubsForUndeclaredProperties(n, alias, enclosingStatement.GetParent(), classNode)
	}

	// detach `static m() {}` from `class Foo { static m() {} }`
	member := declaration.GetNode().Detach()
	var varDecl *rhino.Node
	if member.HasChildren() {
		memberVal := member.GetOnlyChild()
		c.CheckForReceiverAffectedByCollapse(memberVal, member.GetJSDocInfo(), n)
		memberVal = memberVal.Detach()

		// add a var declaration, creating `var Foo$m = function() {}; class Foo {}`
		varDecl = rhino.IR_VarWithValue(NodeUtil_NewNameWithSrcref(c.outer.compiler, alias, member), memberVal).Srcref(member)
	} else {
		// Member field defs may not have initializers.
		varDecl = rhino.IR_Var(NodeUtil_NewNameWithSrcref(c.outer.compiler, alias, member)).Srcref(member)
	}
	varDecl.InsertBefore(enclosingStatement)
	// We would lose optimization-relevant jsdoc tags here because they are stored on the class
	// member node, not the function node. Copy them over to the new declaration statement so
	// later passes can make use of them.
	varDecl.SetJSDocInfo(member.GetJSDocInfo())
	c.outer.compiler.ReportChangeToEnclosingScope(varDecl)

	// collapsing this name's properties requires updating this Ref
	n.UpdateRefNode(declaration, varDecl.GetFirstChild())
}

// DeclareVariablesForObjLitValues declares global variables to serve as
// aliases for the values in an object literal, optionally removing all of
// the object literal's keys and values.
//
// alias is the object literal's flattened name (e.g. "a$b$c"); objlit the
// OBJLIT node; varNode the VAR node to which new global variables should be
// added as children; nameToAddAfter the child of varNode after which new
// variables should be added (may be null).
//
// Java: InlineAndCollapseProperties.CollapseProperties.declareVariablesForObjLitValues(Name, String, Node, Node, Node)
func (c *InlineAndCollapseProperties_CollapseProperties) DeclareVariablesForObjLitValues(
	objlitName *GlobalNamespace_Name, alias string, objlit *rhino.Node, varNode *rhino.Node, nameToAddAfter *rhino.Node) {
	arbitraryNameCounter := 0
	discardKeys := !objlitName.ShouldKeepKeys()

	var nextKey *rhino.Node
	for key := objlit.GetFirstChild(); key != nil; key = nextKey {
		value := key.GetFirstChild()
		nextKey = key.GetNext()

		// A computed property, or a get or a set can not be rewritten as a VAR. We don't know what
		// properties will be generated by a spread.
		switch key.GetToken() {
		case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_COMPUTED_PROP, rhino.Token_OBJECT_SPREAD:
			continue
		case rhino.Token_STRING_KEY, rhino.Token_MEMBER_FUNCTION_DEF:
		default:
			panic(jutil.NewIllegalStateExceptionMsg(
				"Unexpected child of OBJECTLIT: " + key.ToStringTree()))
		}

		// We generate arbitrary names for keys that aren't valid JavaScript
		// identifiers, since those keys are never referenced. (If they were,
		// this object literal's child names wouldn't be collapsible.) The only
		// reason that we don't eliminate them entirely is the off chance that
		// their values are expressions that have side effects.
		isJsIdentifier := !key.IsNumber() && rhino.TokenStream_IsJSIdentifier(key.GetString())
		var propName string
		if isJsIdentifier {
			propName = key.GetString()
		} else {
			arbitraryNameCounter++
			propName = strconv.Itoa(arbitraryNameCounter)
		}

		// If the name cannot be collapsed, skip it.
		p := inlineAndCollapseProperties_getProperty(objlitName, propName, c.nameMap)
		if p != nil && !c.CanCollapse(p) {
			continue
		}

		propAlias := c.AppendPropForAlias(alias, propName)
		var refNode *rhino.Node
		isConstantName := c.ShouldAddConstantName(key, propAlias)
		if discardKeys {
			key.Detach()
			value.Detach()
			// Don't report a change here because the objlit has already been removed from the tree.
		} else {
			// Substitute a reference for the value.
			refNode = rhino.IR_Name(propAlias)
			if isConstantName {
				refNode.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
			}

			value.ReplaceWith(refNode)
			c.outer.compiler.ReportChangeToEnclosingScope(refNode)
		}

		// Declare the collapsed name as a variable with the original value.
		nameNode := rhino.IR_Name(propAlias)
		nameNode.AddChildToFront(value)
		if isConstantName {
			nameNode.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
		}
		newVar := rhino.IR_Var(nameNode).SrcrefTreeIfMissing(key)
		if nameToAddAfter != nil {
			newVar.InsertAfter(nameToAddAfter)
		} else {
			newVar.InsertBefore(varNode)
		}
		c.outer.compiler.ReportChangeToEnclosingScope(newVar)
		nameToAddAfter = newVar

		// Update the global name's node ancestry if it hasn't already been
		// done. (Duplicate keys in an object literal can bring us here twice
		// for the same global name.)
		if isJsIdentifier && p != nil {
			if !discardKeys {
				p.AddAliasingGetClonedFromDeclaration(refNode)
			}

			p.UpdateRefNode(p.getDeclarationRef(), nameNode)

			if value.IsFunction() {
				c.CheckForReceiverAffectedByCollapse(value, key.GetJSDocInfo(), p)
			}
		}
	}
}

// AddStubsForUndeclaredProperties adds global variable "stubs" for any
// properties of a global name that are only set in a local scope or read
// but never set.
//
// n is an object representing a global name (e.g. "a", "a.b.c"); alias the
// flattened name of the object whose properties we are adding stubs for
// (e.g. "a$b$c"); parent the node to which new global variables should be
// added as children; addAfter the child of after which new variables should
// be added.
//
// Java: InlineAndCollapseProperties.CollapseProperties.addStubsForUndeclaredProperties(Name, String, Node, Node)
func (c *InlineAndCollapseProperties_CollapseProperties) AddStubsForUndeclaredProperties(
	n *GlobalNamespace_Name, alias string, parent *rhino.Node, addAfter *rhino.Node) {
	jutil.CheckStateMsg1(n.CanCollapseUnannotatedChildNames(), "%s", n)
	jutil.CheckArgumentMsg1(NodeUtil_IsStatementBlock(parent), "%s", parent)
	jutil.CheckNotNull(addAfter)
	if n.props == nil {
		return
	}
	for p := range n.props.All() {
		if !p.NeedsToBeStubbed() {
			continue
		}

		propAlias := c.AppendPropForAlias(alias, p.GetBaseName())
		nameNode := rhino.IR_Name(propAlias)
		newVar := rhino.IR_Var(nameNode).SrcrefTreeIfMissing(addAfter)
		newVar.InsertAfter(addAfter)

		// Determine if this is a constant var by checking the first
		// reference to it. Don't check the declaration, as it might be null.
		constPropNode := p.GetFirstRef().GetNode()
		isConstantName :=
			// Don't call shouldAddConstantName because it expects to see a declaration node for
			// a name, and this name is, by definition, undeclared.
			constPropNode.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME) ||
				c.outer.compiler.GetCodingConvention().IsConstant(propAlias)
		if isConstantName {
			nameNode.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
		}

		c.outer.compiler.ReportChangeToEnclosingScope(newVar)
		addAfter = newVar
	}
}

// Java: InlineAndCollapseProperties.CollapseProperties.appendPropForAlias(String, String)
func (c *InlineAndCollapseProperties_CollapseProperties) AppendPropForAlias(root string, prop string) string {
	if strings.IndexByte(prop, '$') != -1 {
		// Encode '$' in a property as '$0'. Because '0' cannot be the
		// start of an identifier, this will never conflict with our
		// encoding from '.' -> '$'.
		prop = strings.ReplaceAll(prop, "$", "$0")
	}
	result := root + "$" + prop
	id := 1
	for c.nameMap.ContainsKey(result) {
		result = root + "$" + prop + "$" + strconv.Itoa(id)
		id++
	}
	return result
}

// InlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules
// finds all the module namespace objects which are referenced by a dynamic
// import.
type InlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *InlineAndCollapseProperties_CollapseProperties

	processCommonJSModules bool
	moduleResolutionMode   ModuleLoader_ResolutionMode
}

// Java: InlineAndCollapseProperties.CollapseProperties.FindDynamicallyImportedModules(boolean, ResolutionMode)
func NewInlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules(
	outer *InlineAndCollapseProperties_CollapseProperties,
	processCommonJSModules bool,
	resolutionMode ModuleLoader_ResolutionMode) *InlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules {
	return &InlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules{
		outer:                  outer,
		processCommonJSModules: processCommonJSModules,
		moduleResolutionMode:   resolutionMode,
	}
}

// Java: InlineAndCollapseProperties.CollapseProperties.FindDynamicallyImportedModules.visit(NodeTraversal, Node, Node)
func (f *InlineAndCollapseProperties_CollapseProperties_FindDynamicallyImportedModules) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	// After rewriting, CommonJS dynamic imports are of the form
	// __webpack_require__.e(2).then(function() { return module$mod1.default; })
	//
	// Mark the base module name as being dynamically imported.
	if f.processCommonJSModules &&
		n.IsGetProp() &&
		n.IsQualifiedName() &&
		n.GetParent() != nil &&
		n.GetParent().IsReturn() &&
		n.GetGrandparent().IsBlock() &&
		n.GetGrandparent().HasOneChild() &&
		n.GetGrandparent().GetParent().IsFunction() {
		potentialCallback := NodeUtil_GetEnclosingFunction(n)
		if potentialCallback != nil &&
			ProcessCommonJSModules_IsCommonJsDynamicImportCallback(
				NodeUtil_GetEnclosingFunction(potentialCallback), f.moduleResolutionMode) {
			f.outer.dynamicallyImportedModules.Add(NodeUtil_GetRootOfQualifiedNameString(*n.GetQualifiedName()))
		}
	} else if ConvertChunksToESModules_IsDynamicImportCallback(n) {
		moduleNamespace := ConvertChunksToESModules_GetDynamicImportCallbackModuleNamespace(f.outer.outer.compiler, n)
		if moduleNamespace != nil {
			f.outer.dynamicallyImportedModules.Add(*moduleNamespace.GetQualifiedName())
		}
	}
}

// Java: InlineAndCollapseProperties.isSafeNamespaceReinit(Ref)
func InlineAndCollapseProperties_IsSafeNamespaceReinit(ref *GlobalNamespace_Ref) bool {
	// allow "a = a || {}" or "var a = a || {}" or "var a;"
	valParent := InlineAndCollapseProperties_GetValueParent(ref)
	val := valParent.GetLastChild()
	if val != nil && val.IsOr() {
		maybeName := val.GetFirstChild()
		if ref.GetNode().MatchesQualifiedNameNode(maybeName) {
			return true
		}
	}

	return false
}

// InlineAndCollapseProperties_GetValueParent gets the parent node of the
// value for any assignment to a Name. For example, in the assignment `var x
// = 3;` the parent would be the NAME node.
//
// Java: InlineAndCollapseProperties.getValueParent(Ref)
func InlineAndCollapseProperties_GetValueParent(ref *GlobalNamespace_Ref) *rhino.Node {
	// there are four types of declarations: VARs, LETs, CONSTs, and ASSIGNs
	n := ref.GetNode().GetParent()
	if n != nil && NodeUtil_IsNameDeclaration(n) {
		return ref.GetNode()
	}
	return ref.GetNode().GetParent()
}

// InlineAndCollapseProperties_InlineAliases inlines constant aliases
//
// This pass was originally necessary because typechecking did not handle
// type aliases well. Now typechecking understands type aliases. In theory,
// this pass can be deleted, but in practice this pass affects some check
// passes that run post-typechecking.
//
// This alias inliner is not very aggressive. It will only inline explicitly
// const aliases but not effectively const ones (for example ones that are
// only ever assigned a value once). This is done to be conservative since
// it's not a good idea to be making dramatic AST changes during checks (or
// really, any AST changes at all). There is a more aggressive alias inliner
// that runs at the start of optimization.
//
// TODO(b/124915436): Delete this pass.
type InlineAndCollapseProperties_InlineAliases struct {
	outer *InlineAndCollapseProperties

	aliases    *jutil.LinkedHashMap[string, string]
	namespace  *GlobalNamespace
	astFactory *AstFactory
}

// Java: InlineAndCollapseProperties.InlineAliases()
func NewInlineAndCollapseProperties_InlineAliases(outer *InlineAndCollapseProperties) *InlineAndCollapseProperties_InlineAliases {
	a := &InlineAndCollapseProperties_InlineAliases{
		outer:   outer,
		aliases: jutil.NewLinkedHashMap[string, string](),
	}
	a.astFactory = outer.compiler.CreateAstFactory()
	return a
}

// Java: InlineAndCollapseProperties.InlineAliases.process(Node, Node)
func (a *InlineAndCollapseProperties_InlineAliases) Process(externs *rhino.Node, root *rhino.Node) {
	a.namespace = NewGlobalNamespaceWithRoots(a.outer.compiler, externs, root)
	defer func() {
		a.namespace.Release()
		a.namespace = nil
	}()
	NodeTraversal_TraverseRoots(a.outer.compiler, NewInlineAndCollapseProperties_InlineAliases_AliasesCollector(a), externs, root)
	if !a.aliases.IsEmpty() {
		NodeTraversal_TraverseRoots(a.outer.compiler, NewInlineAndCollapseProperties_InlineAliases_AliasesInliner(a), externs, root)
	}
}

// InlineAndCollapseProperties_InlineAliases_AliasesCollector is
// InlineAliases.AliasesCollector.
type InlineAndCollapseProperties_InlineAliases_AliasesCollector struct {
	NodeTraversal_ExternsSkippingCallbackBase
	outer *InlineAndCollapseProperties_InlineAliases
}

// Java: InlineAndCollapseProperties.InlineAliases.AliasesCollector()
func NewInlineAndCollapseProperties_InlineAliases_AliasesCollector(outer *InlineAndCollapseProperties_InlineAliases) *InlineAndCollapseProperties_InlineAliases_AliasesCollector {
	return &InlineAndCollapseProperties_InlineAliases_AliasesCollector{outer: outer}
}

// Java: InlineAndCollapseProperties.InlineAliases.AliasesCollector.visit(NodeTraversal, Node, Node)
func (c *InlineAndCollapseProperties_InlineAliases_AliasesCollector) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_VAR, rhino.Token_CONST, rhino.Token_LET:
		if n.HasOneChild() && t.InGlobalScope() {
			c.VisitAliasDefinition(n.GetFirstChild(), NodeUtil_GetBestJSDocInfo(n.GetFirstChild()))
		}
	case rhino.Token_ASSIGN:
		if parent != nil && parent.IsExprResult() && t.InGlobalScope() {
			c.VisitAliasDefinition(n.GetFirstChild(), n.GetJSDocInfo())
		}
	default:
	}
}

// VisitAliasDefinition maybe records that given lvalue is an alias of the
// qualified name on its rhs. Note that since we are doing a post-order
// traversal, any previous aliases contained in the rhs will have already
// been substituted by the time we record the new alias.
//
// Java: InlineAndCollapseProperties.InlineAliases.AliasesCollector.visitAliasDefinition(Node, JSDocInfo)
func (c *InlineAndCollapseProperties_InlineAliases_AliasesCollector) VisitAliasDefinition(lhs *rhino.Node, info *rhino.JSDocInfo) {
	if c.IsDeclaredConst(lhs, info) &&
		(info == nil || !info.HasTypeInformation()) {
		if lhsQName, ok := c.outer.namespace.qnameStringCached(lhs); ok {
			rhs := NodeUtil_GetRValueOfLValue(lhs)
			if rhsQName, ok := c.outer.namespace.qnameStringCached(rhs); ok {
				lhsName, _ := c.outer.namespace.GetOwnSlot(lhsQName).(*GlobalNamespace_Name)
				rhsName, _ := c.outer.namespace.GetOwnSlot(rhsQName).(*GlobalNamespace_Name)
				if lhsName != nil &&
					lhsName.CalculateInlinability().ShouldInlineUsages() &&
					rhsName != nil &&
					rhsName.CalculateInlinability().ShouldInlineUsages() {
					c.outer.aliases.Put(lhsQName, rhsQName)
				}
			}
		}
	}
}

// Java: InlineAndCollapseProperties.InlineAliases.AliasesCollector.isDeclaredConst(Node, JSDocInfo)
func (c *InlineAndCollapseProperties_InlineAliases_AliasesCollector) IsDeclaredConst(lhs *rhino.Node, info *rhino.JSDocInfo) bool {
	if info != nil && info.HasConstAnnotation() {
		return true
	}
	return lhs.GetParent().IsConst()
}

// InlineAndCollapseProperties_InlineAliases_AliasesInliner is
// InlineAliases.AliasesInliner.
type InlineAndCollapseProperties_InlineAliases_AliasesInliner struct {
	NodeTraversal_ExternsSkippingCallbackBase
	outer *InlineAndCollapseProperties_InlineAliases
}

// Java: InlineAndCollapseProperties.InlineAliases.AliasesInliner()
func NewInlineAndCollapseProperties_InlineAliases_AliasesInliner(outer *InlineAndCollapseProperties_InlineAliases) *InlineAndCollapseProperties_InlineAliases_AliasesInliner {
	return &InlineAndCollapseProperties_InlineAliases_AliasesInliner{outer: outer}
}

// Java: InlineAndCollapseProperties.InlineAliases.AliasesInliner.visit(NodeTraversal, Node, Node)
func (i *InlineAndCollapseProperties_InlineAliases_AliasesInliner) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_NAME, rhino.Token_GETPROP:
		if i.outer.aliases.IsEmpty() {
			return
		}
		qname, ok := i.outer.namespace.qnameStringCached(n)
		if ok && i.outer.aliases.ContainsKey(qname) {
			if i.IsLeftmostNameLocal(t, n) {
				// The alias is shadowed by a local variable. Don't rewrite.
				return
			}
			if NodeUtil_IsNameDeclOrSimpleAssignLhs(n, parent) {
				// The node defines an alias. Don't rewrite.
				return
			}

			newNode := i.outer.astFactory.CreateQName(i.outer.namespace, i.ResolveAlias(qname, n))
			if i.IsLeftmostNameLocal(t, newNode) {
				// The aliased name is shadowed by a local variable. Don't rewrite.
				return
			}

			// If n is get_prop like "obj.foo" then newNode should use only location of foo, not
			// obj.foo.
			newNode.SrcrefTree(n)
			// Similarly if n is get_prop like "obj.foo" we should index only foo. obj should not
			// be indexed as it's invisible to users.
			if newNode.IsGetProp() {
				newNode.GetFirstChild().MakeNonIndexableRecursive()
			}
			n.ReplaceWith(newNode)
			t.ReportCodeChange()
		}
	default:
	}
}

// Java: InlineAndCollapseProperties.InlineAliases.AliasesInliner.isLeftmostNameLocal(NodeTraversal, Node)
func (i *InlineAndCollapseProperties_InlineAliases_AliasesInliner) IsLeftmostNameLocal(t *NodeTraversal, n *rhino.Node) bool {
	jutil.CheckState(n.IsQualifiedName())
	leftmostName := NodeUtil_GetRootOfQualifiedName(n).GetString()
	v := t.LookupVar(leftmostName)
	return v != nil && v.IsLocal()
}

// ResolveAlias uses the alias table to look up the resolved name of the
// given alias. If the result is also an alias repeat until the real name is
// resolved.
//
// Java: InlineAndCollapseProperties.InlineAliases.AliasesInliner.resolveAlias(String, Node)
func (i *InlineAndCollapseProperties_InlineAliases_AliasesInliner) ResolveAlias(name string, n *rhino.Node) string {
	aliasPath := jutil.NewLinkedHashSet[string]()
	for i.outer.aliases.ContainsKey(name) {
		if !aliasPath.Add(name) {
			i.outer.outer.compiler.Report(JSError_MakeWithNode(n, InlineAndCollapseProperties_ALIAS_CYCLE, aliasPath.String(), name))

			// Cut the cycle so that it doesn't get reported more than once.
			i.outer.aliases.Remove(name)
			break
		}

		name = i.outer.aliases.Get(name)
	}
	return name
}
