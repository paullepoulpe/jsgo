// Port of com/google/javascript/jscomp/InlineAndCollapseProperties.java @ closure-compiler v20261021

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
	p := &InlineAndCollapseProperties{decisionsLog: diagnostic.LogFile{}}
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
func (p *InlineAndCollapseProperties) Process(externs rhino.Node, root rhino.Node) {
	jutil.CheckStateMsg1(
		!p.testAggressiveInliningOnly || p.propertyCollapseLevel == CompilerOptions_PropertyCollapseLevel_ALL,
		"testAggressiveInlining is invalid for: %s",
		p.propertyCollapseLevel)
	defer func() {
		p.decisionsLog = diagnostic.LogFile{}
		if p.optionalGlobalNamespaceTester == nil {
			p.namespace = nil
		}
	}()
	logFile := p.compiler.CreateOrReopenIndexedLog("InlineAndCollapseProperties", "decisions.log")
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
func (p *InlineAndCollapseProperties) PerformMinimalInliningAndNoCollapsing(externs rhino.Node, root rhino.Node) {
	// TODO(b/124915436): Remove InlineAliases completely after cleaning up the codebase.
	NewInlineAndCollapseProperties_InlineAliases(p).Process(externs, root)
}

// Java: InlineAndCollapseProperties.performMinimalInliningAndModuleExportCollapsing(Node, Node)
func (p *InlineAndCollapseProperties) PerformMinimalInliningAndModuleExportCollapsing(externs rhino.Node, root rhino.Node) {
	// TODO(b/124915436): Remove InlineAliases completely after cleaning up the codebase.
	NewInlineAndCollapseProperties_InlineAliases(p).Process(externs, root)
	// CollapseProperties needs this namespace.
	// TODO(bradfordcsmith): Have `InlineAliases` update the namespace it already created
	// and reuse that one instead.
	p.namespace = NewGlobalNamespaceWithLog(p.decisionsLog, p.compiler, root)
	NewInlineAndCollapseProperties_CollapseProperties(p).Process(externs, root)
}

// Java: InlineAndCollapseProperties.performAggressiveInliningAndCollapsing(Node, Node)
func (p *InlineAndCollapseProperties) PerformAggressiveInliningAndCollapsing(externs rhino.Node, root rhino.Node) {
	NewConcretizeStaticInheritanceForInlining(p.compiler).Process(externs, root)
	NewInlineAndCollapseProperties_AggressiveInlineAliases(p).Process(externs, root)

	NewInlineAndCollapseProperties_CollapseProperties(p).Process(externs, root)
}

// Java: InlineAndCollapseProperties.performAggressiveInliningForTest(Node, Node)
func (p *InlineAndCollapseProperties) PerformAggressiveInliningForTest(externs rhino.Node, root rhino.Node) {
	aggressiveInlineAliases := NewInlineAndCollapseProperties_AggressiveInlineAliases(p)
	aggressiveInlineAliases.Process(externs, root)
	if p.optionalGlobalNamespaceTester == nil {
		// Optional.get() on an empty Optional
		panic(jutil.NewNoSuchElementExceptionMsg("No value present"))
	}
	p.optionalGlobalNamespaceTester(aggressiveInlineAliases.GetLastUsedGlobalNamespace())
}

// Java: InlineAndCollapseProperties.maybeGetInnerNameNode(Node)
func InlineAndCollapseProperties_MaybeGetInnerNameNode(maybeFunctionOrClassNode rhino.Node) rhino.Node {
	if NodeUtil_IsFunctionExpression(maybeFunctionOrClassNode) {
		nameNode := maybeFunctionOrClassNode.GetFirstChild()
		jutil.CheckStateMsg1(nameNode.IsName(), "%s", nameNode)
		// functions with no name have a NAME node with an empty string
		if nameNode.GetString() == "" {
			return rhino.NilNode
		}
		return nameNode
	} else if NodeUtil_IsClassExpression(maybeFunctionOrClassNode) {
		nameNode := maybeFunctionOrClassNode.GetFirstChild()
		// classes with no name have an EMPTY node first child
		if nameNode.IsName() {
			return nameNode
		}
		return rhino.NilNode
	} else {
		return rhino.NilNode // not a function or class expression
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
	if !name.HasProps() {
		return
	}

	if name.GetAliasingGets() == 0 {
		// All of `name`'s children meet condition (a), so they can be
		// added to the worklist.
		name.AppendPropsTo(workList)
	} else {
		// The children do NOT meet condition (a) but we may try to add them anyway.
		// This is because CollapseProperties will unsafely collapse properties on constructors and
		// enums, so we want to be more aggressive about inlining references to their children.
		for property := range name.Props() {
			// Only add properties that would be unsafely collapsed by CollapseProperties
			if property.CanCollapse() {
				*workList = append(*workList, property)
			}
		}
	}
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
	var aliasLhsNode rhino.Node
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
	aliasRefs *ReferenceCollection, aliasedName *GlobalNamespace_Name) bool {
	for _, ref := range aliasRefs.references {
		if ref.GetParent().IsNil() {
			continue
		}
		if NodeUtil_IsNormalOrOptChainGetProp(ref.GetParent()) {
			// e.g. if the reference is "alias.b.someProp", this will be "b".
			propertyName := ref.GetParent().GetString()
			// e.g. if the aliased name is "originalName", this will be "originalName.b".
			originalProperty := aliasedName.LookupChild(propertyName)
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
	for ref := range name.Refs() {
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
func InlineAndCollapseProperties_GetSubclassForEs6Superclass(superclass rhino.Node) rhino.Node {
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
	return rhino.NilNode
}
