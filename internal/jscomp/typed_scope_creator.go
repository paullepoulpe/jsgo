// Port of com/google/javascript/jscomp/TypedScopeCreator.java @ closure-compiler v20261015

/*
 * Copyright 2004 The Closure Compiler Authors.
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
	"iter"
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TypedScopeCreator creates the symbol table of variables available in the
// current scope and their types.
//
// Scopes created by this class are very different from scopes created by
// the syntactic scope creator. These scopes have type information, and
// include some qualified names in addition to variables (like
// Class.staticMethod).
//
// When building scope information, also declares relevant information about
// types in the type registry.
type TypedScopeCreator struct {
	compiler                 AbstractCompiler
	typeParsingErrorReporter rhino.ErrorReporter
	validator                *TypeValidator
	codingConvention         CodingConvention
	typeRegistry             *rhino.JSTypeRegistry
	moduleMap                *ModuleMap
	metadataMap              *ModuleMetadataMap
	moduleImportResolver     *ModuleImportResolver
	processClosurePrimitives bool
	memoized                 *jutil.LinkedHashMap[*rhino.Node, *TypedScope]

	// Maps from scope root to declared variable names. Populated by FirstOrderFunctionAnalyzer to
	// reserve names before the TypedScope is populated.
	reservedNamesForScope *jutil.ListMultimap[*rhino.Node, string]

	// Set of functions with non-empty returns, for passing to FunctionTypeBuilder.
	functionsWithNonEmptyReturns *jutil.LinkedHashSet[*rhino.Node]
	// Includes both simple and qualified names.
	escapedVarNames *jutil.ObjectLinkedHashSet[ScopedName]
	// Count of how many times each variable is assigned, for marking effectively final.
	assignedVarNames *jutil.ObjectHashMultiset[ScopedName]

	// For convenience
	unknownType rhino.ObjectType

	// All names imported through goog.requireType. Resolve these after all scopes are created.
	weakImports []*TypedScopeCreator_WeakModuleImport

	unresolvedNodes []*rhino.Node

	// Set of NAME, GETPROP, and STRING_KEY lvalues which should be treated as const declarations when
	// assigned. Treat simple names in this list as if they were declared `const`. E.g. treat `exports
	// = class {};` as `const exports = class {};`. Treat GETPROP and STRING_KEY nodes as if they were
	// annotated @const.
	undeclaredNamesForClosure *jutil.LinkedHashSet[*rhino.Node]

	// Maps EXPR_RESULT nodes from goog.provides to all implicitly provided names from the call
	providedNamesFromCall *jutil.LinkedHashMultimap[*rhino.Node, *ProcessClosureProvidesAndRequires_ProvidedName]

	stage TypedScopeCreator_Stage
}

var TypedScopeCreator_MALFORMED_TYPEDEF = DiagnosticType_Warning(
	"JSC_MALFORMED_TYPEDEF", "Typedef for {0} does not have any type information")

var TypedScopeCreator_ENUM_INITIALIZER = DiagnosticType_Warning(
	"JSC_ENUM_INITIALIZER_NOT_ENUM", "enum initializer must be an object literal or an enum")

var TypedScopeCreator_INVALID_ENUM_KEY = DiagnosticType_Warning(
	"JSC_INVALID_ENUM_KEY", "enum key must be a string or numeric literal")

var TypedScopeCreator_CTOR_INITIALIZER = DiagnosticType_Warning(
	"JSC_CTOR_INITIALIZER_NOT_CTOR", "Constructor {0} must be initialized at declaration")

var TypedScopeCreator_IFACE_INITIALIZER = DiagnosticType_Warning(
	"JSC_IFACE_INITIALIZER_NOT_IFACE", "Interface {0} must be initialized at declaration")

var TypedScopeCreator_CONSTRUCTOR_EXPECTED = DiagnosticType_Warning(
	"JSC_REFLECT_CONSTRUCTOR_EXPECTED", "Constructor expected as first argument")

var TypedScopeCreator_UNKNOWN_LENDS = DiagnosticType_Warning(
	"JSC_UNKNOWN_LENDS", "Variable {0} not declared before @lends annotation.")

var TypedScopeCreator_LENDS_ON_NON_OBJECT = DiagnosticType_Warning(
	"JSC_LENDS_ON_NON_OBJECT", "May only lend properties to object types. {0} has type {1}.")

var TypedScopeCreator_INCOMPATIBLE_ALIAS_ANNOTATION = DiagnosticType_Warning(
	"JSC_INCOMPATIBLE_ALIAS_ANNOTATION",
	"Annotation {0} on {1} incompatible with aliased type.")

var TypedScopeCreator_DYNAMIC_EXTENDS_WITHOUT_JSDOC = DiagnosticType_Warning(
	"JSC_DYNAMIC_EXTENDS_WITHOUT_JSDOC",
	"The right-hand side of an extends clause must be a qualified name, or else @extends must"+
		" be specified in JSDoc")

// TypedScopeCreator_Stage is the stage of a TypedScopeCreator.
type TypedScopeCreator_Stage int32

const (
	TypedScopeCreator_Stage_BUILDING TypedScopeCreator_Stage = iota
	TypedScopeCreator_Stage_FROZEN
)

var typedScopeCreator_Stage_names = [...]string{"BUILDING", "FROZEN"}

// Java: TypedScopeCreator.Stage.values()
func TypedScopeCreator_Stage_values() []TypedScopeCreator_Stage {
	return []TypedScopeCreator_Stage{TypedScopeCreator_Stage_BUILDING, TypedScopeCreator_Stage_FROZEN}
}

// String is Enum.name().
func (s TypedScopeCreator_Stage) String() string { return typedScopeCreator_Stage_names[s] }

// Ordinal is Enum.ordinal().
func (s TypedScopeCreator_Stage) Ordinal() int { return int(s) }

// TypedScopeCreator_WeakModuleImport tracks information about a weakly
// ordered import, such as goog.requireType or goog.forwardDeclare.
//
// TypedScopeCreator traverses through all files once in AST order, so when
// visiting a goog.requireType() call, the required module may not have its
// scope populated yet. Thus we create a 'WeakModuleImport' object and revisit
// the import after TypedScopeCreator has traversed all files.
type TypedScopeCreator_WeakModuleImport struct {
	// Go-only: language-adapter — the enclosing TypedScopeCreator instance of
	// this inner class.
	outer *TypedScopeCreator

	moduleLocalNode  *rhino.Node // the NAME or IMPORT_STAR representing the assignee
	scopedImport     ScopedName  // the name (possibly qualified) & scope being imported
	localModuleScope *TypedScope // the scope containing the import
}

// Java: TypedScopeCreator.WeakModuleImport(Node, ScopedName, TypedScope)
func NewTypedScopeCreator_WeakModuleImport(outer *TypedScopeCreator, moduleLocalNode *rhino.Node, scopedImport ScopedName, localModuleScope *TypedScope) *TypedScopeCreator_WeakModuleImport {
	w := &TypedScopeCreator_WeakModuleImport{outer: outer}
	jutil.CheckStateMsg1(moduleLocalNode.IsName() || moduleLocalNode.IsImportStar(), "%s", moduleLocalNode)
	w.moduleLocalNode = moduleLocalNode
	w.scopedImport = scopedImport
	w.localModuleScope = localModuleScope
	return w
}

// Java: TypedScopeCreator.WeakModuleImport.resolve()
func (w *TypedScopeCreator_WeakModuleImport) Resolve() {
	// resolvedScope may be null if this is importing a nonexistent module
	resolvedScope := w.outer.memoized.Get(w.scopedImport.GetScopeRoot())
	var requiredVar *TypedVar
	if resolvedScope != nil {
		requiredVar, _ = resolvedScope.GetSlot(w.scopedImport.GetName()).(*TypedVar)
	}

	var type_ rhino.JSType
	var isInferred bool
	if requiredVar != nil {
		type_ = requiredVar.GetType()
		isInferred = requiredVar.IsTypeInferred()
	} else if resolvedScope != nil {
		// Some imports might not exist as fully qualified names in the given scope, but are still
		// resolvable via properties. This is always true for "provideAlreadyProvided" names.
		nsType := resolvedScope.GetTypeThroughNamespace(w.scopedImport.GetName())
		if nsType != nil {
			type_ = nsType
		} else {
			type_ = w.outer.unknownType
		}
		isInferred = nsType == nil
	} else {
		// if resolvedScope is null, this code is importing a module that doesn't actually exist
		// and we'll error elsewhere.
		type_ = w.outer.unknownType
		isInferred = true
	}
	w.localModuleScope.Declare(
		w.moduleLocalNode.GetString(),
		w.moduleLocalNode,
		type_,
		w.outer.compiler.GetInput(NodeUtil_GetInputId(w.moduleLocalNode)),
		isInferred)
	if requiredVar != nil && requiredVar.GetNameNode().GetTypedefTypeProp() != nil {
		// Propagate the 'typedef type' from the module export to this variable. Otherwise
		// NamedTypes pointing to the imported name fail to resolve.
		typedefType := requiredVar.GetNameNode().GetTypedefTypeProp()
		w.moduleLocalNode.SetTypedefTypeProp(typedefType)
		w.outer.typeRegistry.DeclareType(w.localModuleScope, w.moduleLocalNode.GetString(), typedefType)
	}
}

// TypedScopeCreator_RValueInfo stores the type and qualified name for a
// destructuring rvalue.
type TypedScopeCreator_RValueInfo struct {
	type_         rhino.JSType
	qualifiedName rhino.QualifiedName
}

// Java: TypedScopeCreator.RValueInfo(JSType, QualifiedName)
func NewTypedScopeCreator_RValueInfo(type_ rhino.JSType, qualifiedName rhino.QualifiedName) *TypedScopeCreator_RValueInfo {
	r := &TypedScopeCreator_RValueInfo{}
	r.type_ = type_
	r.qualifiedName = qualifiedName
	return r
}

// Java: TypedScopeCreator.RValueInfo.empty()
func TypedScopeCreator_RValueInfo_Empty() *TypedScopeCreator_RValueInfo {
	return NewTypedScopeCreator_RValueInfo(nil, nil)
}

// Java: TypedScopeCreator(AbstractCompiler)
func NewTypedScopeCreator(compiler AbstractCompiler) *TypedScopeCreator {
	return NewTypedScopeCreatorWithCodingConvention(compiler, compiler.GetCodingConvention())
}

// Java: TypedScopeCreator(AbstractCompiler, CodingConvention)
func NewTypedScopeCreatorWithCodingConvention(compiler AbstractCompiler, codingConvention CodingConvention) *TypedScopeCreator {
	c := &TypedScopeCreator{}
	c.memoized = jutil.NewLinkedHashMap[*rhino.Node, *TypedScope]()
	c.stage = TypedScopeCreator_Stage_BUILDING
	c.compiler = compiler
	c.validator = compiler.GetTypeValidator()
	c.codingConvention = codingConvention
	c.typeRegistry = compiler.GetTypeRegistry()
	c.typeParsingErrorReporter = c.typeRegistry.GetErrorReporter()
	c.unknownType = c.typeRegistry.GetNativeObjectType(rhino.JSTypeNative_UNKNOWN_TYPE)
	if compiler.GetModuleMetadataMap() != nil {
		c.metadataMap = compiler.GetModuleMetadataMap()
	} else {
		c.metadataMap = NewModuleMetadataMap(
			jutil.NewLinkedHashMap[string, *ModuleMetadataMap_ModuleMetadata](),
			jutil.NewLinkedHashMap[string, *ModuleMetadataMap_ModuleMetadata]())
	}
	c.moduleMap = compiler.GetModuleMap()
	c.moduleImportResolver = NewModuleImportResolver(c.moduleMap, c.GetNodeToScopeMapper(), c.typeRegistry)
	c.processClosurePrimitives = !c.metadataMap.GetModulesByGoogNamespace().IsEmpty()

	// Reset state to empty collections.
	c.ClearCommonState()
	return c
}

// Java: TypedScopeCreator.report(JSError)
func (c *TypedScopeCreator) Report(error_ *JSError) {
	c.compiler.Report(error_)
}

// Java: TypedScopeCreator.getReferences(TypedVar)
func (c *TypedScopeCreator) GetReferences(var_ *TypedVar) iter.Seq[*TypedVar] {
	return slices.Values([]*TypedVar{var_})
}

// Java: TypedScopeCreator.getScope(TypedVar)
func (c *TypedScopeCreator) GetScope(var_ *TypedVar) rhino.StaticScope {
	return var_.GetScope()
}

// Java: TypedScopeCreator.getAllSymbols()
func (c *TypedScopeCreator) GetAllSymbols() iter.Seq[*TypedVar] {
	vars := []*TypedVar{}
	for s := range c.memoized.Values() {
		for v := range s.GetAllSymbols() {
			vars = append(vars, v)
		}
	}
	return slices.Values(vars)
}

// GetNodeToScopeMapper returns a function mapping a scope root node to a
// TypedScope.
//
// This method mostly exists in lieu of an interface representing root node
// -> scope.
//
// Java: TypedScopeCreator.getNodeToScopeMapper()
func (c *TypedScopeCreator) GetNodeToScopeMapper() func(*rhino.Node) *TypedScope {
	return c.memoized.Get
}

// Java: TypedScopeCreator.getAllMemoizedScopes()
func (c *TypedScopeCreator) GetAllMemoizedScopes() []*TypedScope {
	// Return scopes in reverse order of creation so that IIFEs will
	// come before the global scope.
	scopes := slices.Collect(c.memoized.Values())
	slices.Reverse(scopes)
	return scopes
}

// CreateScopeOfNode creates a scope if it doesn't already exist, looking up
// in the map for the parent scope.
//
// Go-only: language-adapter — the overload createScope(Node) is named
// apart from the ScopeCreator method CreateScope(Node, AbstractScope).
//
// Java: TypedScopeCreator.createScope(Node)
func (c *TypedScopeCreator) CreateScopeOfNode(n *rhino.Node) *TypedScope {
	s := c.memoized.Get(n)
	if s != nil {
		return s
	}
	return c.CreateScope(n, c.CreateScopeOfNode(NodeUtil_GetEnclosingScopeRoot(n.GetParent()))).(*TypedScope)
}

// CreateScope creates a scope with all types declared. Declares newly
// discovered types and type properties in the type registry.
//
// Java: TypedScopeCreator.createScope(Node, AbstractScope)
func (c *TypedScopeCreator) CreateScope(root *rhino.Node, parent AbstractScope) AbstractScope {
	_, parentIsTyped := parent.(*TypedScope)
	jutil.CheckArgument(parent == nil || parentIsTyped)
	var typedParent *TypedScope
	if parent != nil {
		typedParent = parent.(*TypedScope)
	}

	scope := c.memoized.Get(root)
	if scope != nil {
		jutil.CheckState(typedParent == scope.GetParent())
	} else {
		scope = c.CreateScopeInternal(root, typedParent)
		c.memoized.Put(root, scope)
	}
	return scope
}

// Java: TypedScopeCreator.createScopeInternal(Node, TypedScope)
func (c *TypedScopeCreator) CreateScopeInternal(root *rhino.Node, typedParent *TypedScope) *TypedScope {
	jutil.CheckStateMsg1(
		c.stage == TypedScopeCreator_Stage_BUILDING,
		"Cannot create scope for %s after TypedScopeCreator is frozen",
		root)
	// Constructing the global scope is very different than constructing
	// inner scopes, because only global scopes can contain named classes that
	// show up in the type registry.
	var newScope *TypedScope

	var scopeBuilder TypedScopeCreator_AbstractScopeBuilder

	module := ModuleImportResolver_GetModuleFromScopeRoot(c.moduleMap, c.compiler, root)
	if typedParent == nil {
		jutil.CheckStateMsg1(root.IsRoot(), "%s", root)
		externsRoot := root.GetFirstChild()
		jsRoot := root.GetSecondChild()
		jutil.CheckStateMsg1(externsRoot.IsRoot(), "%s", externsRoot)
		jutil.CheckStateMsg1(jsRoot.IsRoot(), "%s", jsRoot)
		globalThis := c.typeRegistry.GetNativeObjectType(rhino.JSTypeNative_GLOBAL_THIS)

		// Mark the main root, the externs root, and the src root
		// with the global this type.
		root.SetJSType(globalThis)
		externsRoot.SetJSType(globalThis)
		jsRoot.SetJSType(globalThis)

		// Find all the classes in the global scope.
		newScope = c.CreateInitialScope(root)
	} else {
		// Because JSTypeRegistry#getType looks up the scope in which a root of a qualified name is
		// declared, pre-populate this TypedScope with all qualified name roots. This prevents
		// type resolution from accidentally returning a type from an outer scope that is shadowed.
		reservedNames := jutil.NewLinkedHashSet[string]()
		reservedNames.AddAll(c.reservedNamesForScope.RemoveAll(root)...)
		if module != nil && module.Metadata().IsGoogModule() {
			// TypedScopeCreator treats default export assignments, like `exports = class {};`, as
			// declarations. However, the untyped scope only contains an implicit slot for `exports`.
			reservedNames.Add("exports")
		} else if root.IsFunction() && NodeUtil_IsBundledGoogModuleCall(root.GetParent()) {
			// Pretend that 'exports' is declared in the block of goog.loadModule
			// functions, not the function scope. See the above comment for why.
			reservedNames.Remove("exports")
		}

		newScope = NewTypedScopeWithReservedNames(typedParent, root, reservedNames, module)
	}

	if root.IsFunction() {
		scopeBuilder = NewTypedScopeCreator_FunctionScopeBuilder(c, newScope)
	} else if root.IsClass() || root.IsMemberFieldDef() || root.IsComputedFieldDef() {
		scopeBuilder = NewTypedScopeCreator_ClassScopeBuilder(c, newScope)
	} else {
		scopeBuilder = NewTypedScopeCreator_NormalScopeBuilder(c, newScope)
		if root.IsModuleBody() {
			// Store the module scope on the CompilerInput object so it can be found and used for
			// program analysis purposes.
			// DefaultPassConfig is responsible for removing it when it is no longer needed.
			inputId := jutil.CheckNotNull(NodeUtil_GetInputId(root))
			compilerInput := jutil.CheckNotNull(c.compiler.GetInput(inputId))
			compilerInput.SetTypedScope(newScope)
		}
	}
	scopeBuilder.Build()
	if module != nil && module.Metadata().IsEs6Module() {
		// Declare an implicit variable representing the namespace of this module, then add a property
		// for each exported name to that variable's type.
		namespaceType := c.typeRegistry.CreateAnonymousObjectType(nil)
		newScope.Declare(
			Export_NAMESPACE,
			root, // Use the given MODULE_BODY as the 'declaration node' for lack of a better option.
			namespaceType,
			c.compiler.GetInput(NodeUtil_GetInputId(root)),
			/* inferred= */ false)

		// Store the module object type on the MODULE_BODY.
		// The Es6RewriteModules will retrieve it from there for use when it creates a global for
		// the module object.
		root.SetJSType(namespaceType)
		c.moduleImportResolver.UpdateEsModuleNamespaceType(namespaceType, module, newScope)
	}

	return newScope
}

// MarkGoogModuleExportsAsConst adds nodes representing goog.module exports
// to a list, to treat them as @const.
//
// This method handles the following styles of exports:
//
//   - `exports = class {}` adds the NAME node `exports`
//   - `exports = {Foo};` adds the NAME node `exports` and the STRING_KEY
//     node `Foo`
//   - `exports.Foo = Foo;` adds the GETPROP node `exports.Foo`
//
// Java: TypedScopeCreator.markGoogModuleExportsAsConst(Node)
func (c *TypedScopeCreator) MarkGoogModuleExportsAsConst(moduleBody *rhino.Node) {
	// TODO(lharker): Use the source nodes from the Bindings once we no longer rewrite before
	// typechecking. This is not feasible currently because a few places will rewrite exports = ...
	for statement := moduleBody.GetFirstChild(); statement != nil; statement = statement.GetNext() {
		if !NodeUtil_IsExprAssign(statement) {
			continue
		}
		lhs := statement.GetFirstFirstChild()
		if lhs.MatchesName("exports") {
			c.undeclaredNamesForClosure.Add(lhs)
			// If this is full of named exports, add all the string key nodes.
			if NodeUtil_IsNamedExportsLiteral(lhs.GetNext()) {
				for key := lhs.GetNext().GetFirstChild(); key != nil; key = key.GetNext() {
					c.undeclaredNamesForClosure.Add(key)
				}
			}
		} else if lhs.IsGetProp() && lhs.GetFirstChild().MatchesName("exports") {
			c.undeclaredNamesForClosure.Add(lhs)
		}
	}
}

// GatherAllProvides gathers all namespaces created by goog.provide and any
// definitions in code.
//
// This method does not actually declare anything in the scope. In order to
// accurately report redefinition warnings, wait to declare implicit names
// until the actual goog.provide call.
//
// root is the global ROOT or a SCRIPT.
//
// Java: TypedScopeCreator.gatherAllProvides(Node)
func (c *TypedScopeCreator) GatherAllProvides(root *rhino.Node) {
	if !c.processClosurePrimitives {
		return
	}

	externs := root.GetFirstChild()
	js := root.GetSecondChild()
	providedNames :=
		NewProcessClosureProvidesAndRequires(c.compiler /* preserveGoogProvidesAndRequires= */, true).
			CollectProvidedNames(externs, js)

	for name := range providedNames.Values() {
		if name.HasImplicitInitialization() {
			continue
		}
		if name.GetCandidateDefinition() != nil {
			// This name will be defined eventually in the source code.
			firstDefinitionNode := name.GetCandidateDefinition()
			if NodeUtil_IsExprAssign(firstDefinitionNode) &&
				firstDefinitionNode.GetFirstFirstChild().IsName() {
				// Treat assignments of provided names as declarations.
				c.undeclaredNamesForClosure.Add(firstDefinitionNode.GetFirstFirstChild())
			}
		} else if name.GetFirstProvideCall() != nil &&
			NodeUtil_IsExprCall(name.GetFirstProvideCall()) &&
			!name.IsFromLegacyModule() {
			// This name is implicitly created by a goog.provide call; declare it in the scope once
			// reaching the provide call. The exception is legacy goog.modules, which are declared
			// once leaving the module.
			c.providedNamesFromCall.Put(name.GetFirstProvideCall(), name)
		}

		if name.IsExplicitlyProvided() && !name.IsFromLegacyModule() {
			c.typeRegistry.RegisterLegacyClosureNamespace(name.GetNamespace())
		}
	}
}

// CreateInitialScope creates the outermost scope. This scope contains
// native binding such as Object, Date, etc.
//
// root is the global ROOT node.
//
// Java: TypedScopeCreator.createInitialScope(Node)
func (c *TypedScopeCreator) CreateInitialScope(root *rhino.Node) *TypedScope {
	jutil.CheckArgumentMsg1(root.IsRoot(), "%s", root)

	// Gather global information used in typed scope creation. Use a memoized scope creator because
	// scope-building takes a nontrivial amount of time.
	scopeCreator := NewMemoizedScopeCreator(NewSyntacticScopeCreator(c.compiler))

	NodeTraversal_Builder_().
		SetCompiler(c.compiler).
		SetCallback(NewTypedScopeCreator_FirstOrderFunctionAnalyzer(c)).
		SetScopeCreator(scopeCreator).
		TraverseRoots(root.GetFirstChild(), root.GetLastChild())

	NodeTraversal_Builder_().
		SetCompiler(c.compiler).
		SetCallback(NewTypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable(c.typeRegistry)).
		SetScopeCreator(scopeCreator).
		Traverse(root)

	s := TypedScope_CreateGlobalScope(root)
	c.DeclareNativeFunctionType(s, rhino.JSTypeNative_ARRAY_FUNCTION_TYPE)
	c.DeclareNativeFunctionType(s, rhino.JSTypeNative_BIGINT_OBJECT_FUNCTION_TYPE)
	c.DeclareNativeFunctionType(s, rhino.JSTypeNative_BOOLEAN_OBJECT_FUNCTION_TYPE)
	c.DeclareNativeFunctionType(s, rhino.JSTypeNative_DATE_FUNCTION_TYPE)
	c.DeclareNativeFunctionType(s, rhino.JSTypeNative_FUNCTION_FUNCTION_TYPE)
	c.DeclareNativeFunctionType(s, rhino.JSTypeNative_GENERATOR_FUNCTION_TYPE)
	c.DeclareNativeFunctionType(s, rhino.JSTypeNative_ITERABLE_FUNCTION_TYPE)
	c.DeclareNativeFunctionType(s, rhino.JSTypeNative_ITERATOR_FUNCTION_TYPE)
	c.DeclareNativeFunctionType(s, rhino.JSTypeNative_NUMBER_OBJECT_FUNCTION_TYPE)
	c.DeclareNativeFunctionType(s, rhino.JSTypeNative_OBJECT_FUNCTION_TYPE)
	c.DeclareNativeFunctionType(s, rhino.JSTypeNative_REGEXP_FUNCTION_TYPE)
	c.DeclareNativeFunctionType(s, rhino.JSTypeNative_STRING_OBJECT_FUNCTION_TYPE)
	c.DeclareNativeFunctionType(s, rhino.JSTypeNative_SYMBOL_OBJECT_FUNCTION_TYPE)
	c.DeclareNativeValueType(s, "undefined", rhino.JSTypeNative_VOID_TYPE)
	c.AddWellKnownSymbols(s)

	c.GatherAllProvides(root)

	// Memoize the global scope so that module scope creation can access it. See
	// AbstractScopeBuilder#shouldTraverse - modules are traversed early, as if they were always
	// executed when control flow reaches the module body.
	c.memoized.Put(root, s)

	return s
}

// typedScopeCreator_WELL_KNOWN_SYMBOLS are the symbols that are defined in
// the language spec.
//
// NOTE: adding a symbol to this list is necessary only to make it available
// in Clutz output. By default, JSCompiler will note any Symbol.* references
// in externs whether or not they are in this list, but Clutz runs without
// including all externs so relies on this list to know what is / is not a
// well-known symbol.
var typedScopeCreator_WELL_KNOWN_SYMBOLS = []string{
	// go/keep-sorted start
	"Symbol.asyncDispose",
	"Symbol.asyncIterator",
	"Symbol.dispose",
	"Symbol.hasInstance",
	"Symbol.isConcatSpreadable",
	"Symbol.iterator",
	"Symbol.match",
	"Symbol.replace",
	"Symbol.species",
	"Symbol.toPrimitive",
	"Symbol.toStringTag",
	"Symbol.unscopables",
	// go/keep-sorted end
}

// Java: TypedScopeCreator.addWellKnownSymbols(TypedScope)
func (c *TypedScopeCreator) AddWellKnownSymbols(scope *TypedScope) {
	for _, symbol := range typedScopeCreator_WELL_KNOWN_SYMBOLS {
		type_ := rhino.NewKnownSymbolType(c.typeRegistry, symbol)
		TypedScopeCreator_DeclareNativeType(scope, symbol, type_)
	}
}

// Java: TypedScopeCreator.declareNativeFunctionType(TypedScope, JSTypeNative)
func (c *TypedScopeCreator) DeclareNativeFunctionType(scope *TypedScope, tId rhino.JSTypeNative) {
	t := c.typeRegistry.GetNativeFunctionType(tId)
	TypedScopeCreator_DeclareNativeType(scope, *t.GetInstanceType().GetReferenceName(), t)
	TypedScopeCreator_DeclareNativeType(scope, *t.GetPrototype().GetReferenceName(), t.GetPrototype())
}

// Java: TypedScopeCreator.declareNativeValueType(TypedScope, String, JSTypeNative)
func (c *TypedScopeCreator) DeclareNativeValueType(scope *TypedScope, name string, tId rhino.JSTypeNative) {
	TypedScopeCreator_DeclareNativeType(scope, name, c.typeRegistry.GetNativeType(tId))
}

// Java: TypedScopeCreator.declareNativeType(TypedScope, String, JSType)
func TypedScopeCreator_DeclareNativeType(scope *TypedScope, name string, t rhino.JSType) {
	scope.Declare(name, nil, t, nil, false)
}

// SetDeferredType sets the type for a node now.
//
// If the type is unresolved, enqueue it to be updated with a resolved type
// later.
//
// Java: TypedScopeCreator.setDeferredType(Node, JSType)
func (c *TypedScopeCreator) SetDeferredType(node *rhino.Node, type_ rhino.JSType) {
	// Other parts of this pass may read the not-yet-resolved type off the node.
	// (like when we set the LHS of an assign with a typed RHS function.)
	node.SetJSType(type_)
	if !type_.IsResolved() {
		c.unresolvedNodes = append(c.unresolvedNodes, node)
	}
}

// ResolveWeakImportsPreResolution needs to run pre-type-resolution to
// handle weak module imports.
//
// Java: TypedScopeCreator.resolveWeakImportsPreResolution()
func (c *TypedScopeCreator) ResolveWeakImportsPreResolution() {
	// Declare goog.module type requires in scope.
	for _, weakImport := range c.weakImports {
		weakImport.Resolve()
	}
}

// FinishAndFreeze performs some final work to resolve remaining types.
//
// After this call, calling this.createScope will crash if passed a scope
// root that hasn't already been visited.
//
// Java: TypedScopeCreator.finishAndFreeze()
func (c *TypedScopeCreator) FinishAndFreeze() {
	// Resolve types and attach them to nodes.
	for _, node := range c.unresolvedNodes {
		resolved := node.GetJSType().Resolve(c.typeParsingErrorReporter)
		node.SetJSType(resolved)
	}

	// Resolve types and attach them to scope slots.
	for _, scope := range c.GetAllMemoizedScopes() {
		for var_ := range scope.GetVarIterable() {
			var_.ResolveType(c.typeParsingErrorReporter)
		}
		scope.ValidateCompletelyBuilt()
	}

	// free up memory. we no longer need this state now that all scopes have been visited
	c.ClearCommonState()
	c.stage = TypedScopeCreator_Stage_FROZEN

	// we must keep the 'memoized' set of TypedScopes around, since later passes will use it.
}

// Java: TypedScopeCreator.clearCommonState()
func (c *TypedScopeCreator) ClearCommonState() {
	c.reservedNamesForScope = jutil.ArrayListMultimap_Create[*rhino.Node, string]()
	c.functionsWithNonEmptyReturns = jutil.NewLinkedHashSet[*rhino.Node]()
	c.escapedVarNames = jutil.NewObjectLinkedHashSet[ScopedName]()
	c.assignedVarNames = jutil.ObjectHashMultiset_Create[ScopedName]()
	c.weakImports = []*TypedScopeCreator_WeakModuleImport{}
	c.unresolvedNodes = []*rhino.Node{}
	c.undeclaredNamesForClosure = jutil.NewLinkedHashSet[*rhino.Node]()
	c.providedNamesFromCall = jutil.NewLinkedHashMultimap[*rhino.Node, *ProcessClosureProvidesAndRequires_ProvidedName]()
}

// TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable adds all enums and
// typedefs to the registry's list of non-nullable types.
type TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	registry *rhino.JSTypeRegistry
}

// Java: TypedScopeCreator.IdentifyEnumsAndTypedefsAsNonNullable(JSTypeRegistry)
func NewTypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable(registry *rhino.JSTypeRegistry) *TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable {
	v := &TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable{}
	v.registry = registry
	return v
}

// Java: TypedScopeCreator.IdentifyEnumsAndTypedefsAsNonNullable.visit(NodeTraversal, Node, Node)
func (v *TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable) Visit(t *NodeTraversal, node *rhino.Node, parent *rhino.Node) {
	switch node.GetToken() {
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_VAR:
		for child := node.GetFirstChild(); child != nil; child = child.GetNext() {
			// TODO(b/116853368): make this work for destructuring aliases as well.
			v.IdentifyEnumOrTypedefDeclaration(
				t, child, child.GetFirstChild(), NodeUtil_GetBestJSDocInfo(child))
		}
	case rhino.Token_EXPR_RESULT:
		firstChild := node.GetFirstChild()
		if firstChild.IsAssign() {
			assign := firstChild
			v.IdentifyEnumOrTypedefDeclaration(
				t, assign.GetFirstChild(), assign.GetSecondChild(), assign.GetJSDocInfo())
		} else if firstChild.IsGetProp() {
			v.IdentifyEnumOrTypedefDeclaration(
				t, firstChild /* rvalue= */, nil, firstChild.GetJSDocInfo())
		}
	default:
	}
}

// Java: TypedScopeCreator.IdentifyEnumsAndTypedefsAsNonNullable.identifyEnumOrTypedefDeclaration(NodeTraversal, Node, Node, JSDocInfo)
func (v *TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable) IdentifyEnumOrTypedefDeclaration(
	t *NodeTraversal, nameNode *rhino.Node, rvalue *rhino.Node, info *rhino.JSDocInfo) {
	if !nameNode.IsQualifiedName() {
		return
	}
	if info != nil && info.HasEnumParameterType() {
		v.registry.IdentifyNonNullableName(t.GetScope(), *nameNode.GetQualifiedName())
	} else if info != nil && info.HasTypedefType() {
		v.registry.IdentifyNonNullableName(t.GetScope(), *nameNode.GetQualifiedName())
	} else if rvalue != nil &&
		rvalue.IsQualifiedName() &&
		v.registry.IsNonNullableName(t.GetScope(), *rvalue.GetQualifiedName()) &&
		NodeUtil_IsConstantDeclaration(info, nameNode) {
		v.registry.IdentifyNonNullableName(t.GetScope(), *nameNode.GetQualifiedName())
	}
}

// Java: TypedScopeCreator.getNativeType(JSTypeNative)
func (c *TypedScopeCreator) GetNativeType(nativeType rhino.JSTypeNative) rhino.JSType {
	return c.typeRegistry.GetNativeType(nativeType)
}

// TypedScopeCreator_AbstractScopeBuilder is the abstract base of the scope
// builders.
type TypedScopeCreator_AbstractScopeBuilder interface {
	NodeTraversal_Callback

	// Build traverses the scope root and builds it.
	//
	// Java: TypedScopeCreator.AbstractScopeBuilder.build()
	Build()

	// VisitPreorder is called by shouldTraverse on nodes after ensuring the
	// inputId is set.
	//
	// Java: TypedScopeCreator.AbstractScopeBuilder.visitPreorder(NodeTraversal, Node, Node)
	VisitPreorder(t *NodeTraversal, n *rhino.Node, parent *rhino.Node)

	// VisitPostorder is called by visit on nodes after updating the inputId.
	//
	// Java: TypedScopeCreator.AbstractScopeBuilder.visitPostorder(NodeTraversal, Node, Node)
	VisitPostorder(t *NodeTraversal, n *rhino.Node, parent *rhino.Node)

	// GetTypedScopeCreator_AbstractScopeBuilderBase gives access to the
	// fields of the base class.
	//
	// Go-only: language-adapter — field access through the interface
	// (PORTING.md §4).
	GetTypedScopeCreator_AbstractScopeBuilderBase() *TypedScopeCreator_AbstractScopeBuilderBase
}

// TypedScopeCreator_AbstractScopeBuilderBase holds the fields and method
// bodies of TypedScopeCreator.AbstractScopeBuilder.
type TypedScopeCreator_AbstractScopeBuilderBase struct {
	self TypedScopeCreator_AbstractScopeBuilder

	// Go-only: language-adapter — the enclosing TypedScopeCreator instance of
	// this inner class.
	outer *TypedScopeCreator

	// The scope that we're building.
	currentScope *TypedScope

	// The current hoist scope.
	currentHoistScope *TypedScope

	// The current source file that we're in.
	sourceName *string

	// The InputId of the current node.
	inputId *rhino.InputId

	// Some actions need to be deferred, such as analyzing object literals with lends annotations,
	// or resolving type-less stubs. These actions are added to this map, keyed by the node that
	// should be waited for before running.
	//
	// Go-only: language-adapter — a Runnable is a *func(): every Java lambda
	// is a distinct object, which the pointer to a fresh func value keeps.
	deferredActions *jutil.LinkedHashMultimap[*rhino.Node, *func()]
}

// InitTypedScopeCreator_AbstractScopeBuilderBase is the AbstractScopeBuilder
// constructor; self is the concrete builder.
//
// Java: TypedScopeCreator.AbstractScopeBuilder(TypedScope)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) InitTypedScopeCreator_AbstractScopeBuilderBase(
	self TypedScopeCreator_AbstractScopeBuilder, outer *TypedScopeCreator, scope *TypedScope) {
	b.self = self
	b.outer = outer
	b.deferredActions = jutil.NewLinkedHashMultimap[*rhino.Node, *func()]()
	b.currentScope = scope
	b.currentHoistScope = scope.GetClosestHoistScope()
}

// GetTypedScopeCreator_AbstractScopeBuilderBase gives access to the fields
// of the base class.
func (b *TypedScopeCreator_AbstractScopeBuilderBase) GetTypedScopeCreator_AbstractScopeBuilderBase() *TypedScopeCreator_AbstractScopeBuilderBase {
	return b
}

// GetCompilerInput returns the current compiler input.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.getCompilerInput()
func (b *TypedScopeCreator_AbstractScopeBuilderBase) GetCompilerInput() *CompilerInput {
	return b.outer.compiler.GetInput(b.inputId)
}

// Java: TypedScopeCreator.AbstractScopeBuilder.getModule()
func (b *TypedScopeCreator_AbstractScopeBuilderBase) GetModule() *Module {
	return b.currentScope.GetModule()
}

// Build traverses the scope root and builds it.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.build()
func (b *TypedScopeCreator_AbstractScopeBuilderBase) Build() {
	b.InitializeModuleScope(b.currentScope.GetRootNode())

	NodeTraversal_Builder_().
		SetCompiler(b.outer.compiler).
		SetCallback(b.self).
		SetScopeCreator(ScopeCreator_ASSERT_NO_SCOPES_CREATED).
		TraverseAtScope(b.currentScope)

	b.FinishDeclaringGoogModule()
}

// InitializeModuleScope builds the beginning of a module-scope. This can be
// an ES module or a goog.module.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.initializeModuleScope(Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) InitializeModuleScope(moduleBody *rhino.Node) {
	if b.GetModule() == nil {
		return
	}
	moduleType := b.GetModule().Metadata().ModuleType()
	switch moduleType {
	case ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE, ModuleMetadataMap_ModuleType_GOOG_MODULE:
		b.DeclareExportsInGoogModuleScope(b.GetModule(), moduleBody)
		b.outer.MarkGoogModuleExportsAsConst(moduleBody)
	case ModuleMetadataMap_ModuleType_ES6_MODULE:
		unresolvedImports := b.outer.moduleImportResolver.DeclareEsModuleImports(
			b.GetModule(),
			b.currentScope,
			b.outer.compiler.GetInput(NodeUtil_GetInputId(moduleBody)))
		for key, value := range unresolvedImports.All() {
			b.outer.weakImports = append(b.outer.weakImports,
				NewTypedScopeCreator_WeakModuleImport(b.outer, key, value, b.currentScope))
		}
	default:
		// This should have been a build-breaking DUPLICATE_NAMESPACE_AND_MODULE error in
		// GatherModuleMetadata
		panic(jutil.NewIllegalStateExceptionMsg(
			jutil.StringFormat(
				"Unexpected module type %s in module %s", moduleType, b.GetModule())))
	}
}

// DeclareExportsInGoogModuleScope ensures that the name `exports` is
// declared in goog.module scope.
//
// If a goog.module explicitly assigns to exports, we want to treat that
// assignment inside the scope as if it were a declaration:
// `const exports = ...`. This method only handles cases where we want to
// treat exports as implicitly declared.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.declareExportsInGoogModuleScope(Module, Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) DeclareExportsInGoogModuleScope(googModule *Module, moduleBody *rhino.Node) {
	for statement := moduleBody.GetFirstChild(); statement != nil; statement = statement.GetNext() {
		if !NodeUtil_IsExprAssign(statement) {
			continue
		}
		lhs := statement.GetFirstFirstChild()
		if lhs.MatchesName("exports") {
			return // found a direct assignment `exports = [...]`
		}
	}
	root := googModule.Metadata().RootNode()
	// Synthesize an object literal namespace 'exports'
	NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(b).
		ForDeclarationNode(root.GetFirstChild()).
		ForVariableName("exports").
		WithType(b.outer.typeRegistry.CreateAnonymousObjectType(nil)).
		AllowLaterTypeInference(false).
		InScope(b.currentScope).
		DefineSlot()
}

// FinishDeclaringGoogModule declares goog.module.declareLegacyNamespace()
// names in the global scope and annotates the AST with type information
// about module exports.
//
// For example, given
//
//	goog.module('a.b.Foo');
//	goog.module.declareLegacyNamespace();
//
//	exports = class Foo {};
//
// This function is responsible for declaring a global name `a.b.Foo` and
// giving it the type from `exports = class Foo {};`
//
// Java: TypedScopeCreator.AbstractScopeBuilder.finishDeclaringGoogModule()
func (b *TypedScopeCreator_AbstractScopeBuilderBase) FinishDeclaringGoogModule() {
	if b.GetModule() == nil || !b.GetModule().Metadata().IsGoogModule() {
		return
	}
	exportsVar := jutil.CheckNotNull(b.currentScope.GetSlot("exports")).(*TypedVar)

	if b.GetModule().Metadata().IsLegacyGoogModule() {
		b.outer.typeRegistry.RegisterLegacyClosureNamespace(*b.GetModule().ClosureNamespace())
		moduleNamespace := rhino.QualifiedName_Of(*b.GetModule().ClosureNamespace())
		NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(b).
			InScope(b.currentScope.GetGlobalScope()).
			ForDeclarationNode(exportsVar.GetNameNode()).
			ForVariableName(moduleNamespace.Join()).
			WithType(exportsVar.GetType()).
			AllowLaterTypeInference(exportsVar.IsTypeInferred()).
			ForGoogProvidedName().
			DefineSlot()

		if !moduleNamespace.IsSimple() && !exportsVar.IsTypeInferred() {
			parentType :=
				b.currentScope.GetGlobalScope().LookupQualifiedName(moduleNamespace.GetOwner())
			if parentType != nil && parentType.ToMaybeObjectType() != nil {
				// Declare the namespace on the parent name, propagating any JSDoc annotations
				// on "exports = Foo".
				typeDoc := NodeUtil_GetBestJSDocInfo(exportsVar.GetNameNode())
				b.DeclarePropertyIfNamespaceType(
					parentType.ToMaybeObjectType(),
					exportsVar.GetNameNode(),
					rhino.NewProperty_StringKey(moduleNamespace.GetComponent()),
					exportsVar.GetType(),
					exportsVar.GetNameNode(),
					typeDoc)
			}
		}
		b.DeclareAliasTypeIfRvalueIsAliasableWithLValueName(
			b.GetModule().ClosureNamespace(),
			exportsVar.GetNameNode(), // Pretend that 'exports = '... is the lvalue node.
			rhino.QualifiedName_Of("exports"),
			exportsVar.GetType(),
			b.currentScope,
			b.currentScope.GetGlobalScope())
	} else {
		b.outer.typeRegistry.RegisterNonLegacyClosureNamespace(
			*b.GetModule().ClosureNamespace(), exportsVar.GetNameNode(), exportsVar.GetType())
	}
	// Store the type of the namespace on the AST for the convenience of later passes that want
	// to access it.
	rootNode := b.currentScope.GetRootNode()
	if rootNode.IsModuleBody() {
		rootNode.SetJSType(exportsVar.GetType())
	} else {
		// For goog.loadModule, give the `exports` parameter the correct type.
		jutil.CheckStateMsg1(rootNode.IsBlock(), "%s", rootNode)
		fn := rootNode.GetParent()
		paramList := NodeUtil_GetFunctionParameters(fn)
		paramList.GetOnlyChild().SetJSType(exportsVar.GetType())
	}
}

// ShouldTraverse is final in Java.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.shouldTraverse(NodeTraversal, Node, Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	b.inputId = t.GetInputId()

	if n.IsFunction() || n.IsScript() || (parent == nil && b.inputId != nil) {
		jutil.CheckNotNull(b.inputId)
		b.sourceName = NodeUtil_GetSourceName(n)
	}
	if parent == nil || b.InCurrentScope(t) {
		b.self.VisitPreorder(t, n, parent)
		return true
	} else if n.IsModuleBody() {
		// Visit modules pre-order. While this doesn't exactly match execution semantics, it
		// does match how the compiler rewrites modules into the global scope.
		b.outer.CreateScope(n, b.currentScope)
	} else if NodeUtil_IsBundledGoogModuleScopeRoot(n) {
		functionScope := b.outer.CreateScope(parent, b.currentScope).(*TypedScope)
		b.outer.CreateScope(n, functionScope)
	}
	return false
}

// Java: TypedScopeCreator.AbstractScopeBuilder.inCurrentScope(NodeTraversal)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) InCurrentScope(t *NodeTraversal) bool {
	return t.GetScopeRoot() == b.currentScope.GetRootNode()
}

// Visit is final in Java.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.visit(NodeTraversal, Node, Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	isScopeRoot := parent == nil
	b.inputId = t.GetInputId()
	if !isScopeRoot || !n.IsClass() {
		// Do not type literals or visit class scope roots
		b.AttachLiteralTypes(n)
		b.self.VisitPostorder(t, n, parent)
	}
	if !isScopeRoot && b.deferredActions.ContainsKey(n) { // streams are expensive, only make if needed
		for _, action := range b.deferredActions.RemoveAll(n) {
			(*action)()
		}
	} else if isScopeRoot && !b.deferredActions.IsEmpty() {
		// Run *all* remaining deferred actions, in case any were missed.
		for action := range b.deferredActions.Values() {
			(*action)()
		}
	}
}

// VisitPreorder is called by shouldTraverse on nodes after ensuring the
// inputId is set.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.visitPreorder(NodeTraversal, Node, Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) VisitPreorder(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
}

// VisitPostorder is called by visit on nodes after updating the inputId.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.visitPostorder(NodeTraversal, Node, Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) VisitPostorder(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
}

// Java: TypedScopeCreator.AbstractScopeBuilder.attachLiteralTypes(Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) AttachLiteralTypes(n *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_NULL:
		n.SetJSType(b.outer.GetNativeType(rhino.JSTypeNative_NULL_TYPE))
	case rhino.Token_VOID:
		n.SetJSType(b.outer.GetNativeType(rhino.JSTypeNative_VOID_TYPE))
	case rhino.Token_STRINGLIT, rhino.Token_TEMPLATELIT_STRING:
		n.SetJSType(b.outer.GetNativeType(rhino.JSTypeNative_STRING_TYPE))
	case rhino.Token_NUMBER:
		n.SetJSType(b.outer.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE))
	case rhino.Token_BIGINT:
		n.SetJSType(b.outer.GetNativeType(rhino.JSTypeNative_BIGINT_TYPE))
	case rhino.Token_TRUE, rhino.Token_FALSE:
		n.SetJSType(b.outer.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE))
	case rhino.Token_REGEXP:
		n.SetJSType(b.outer.GetNativeType(rhino.JSTypeNative_REGEXP_TYPE))
	case rhino.Token_OBJECTLIT:
		info := n.GetJSDocInfo()
		if info != nil && info.HasLendsName() {
			// Defer analyzing object literals with a @lends annotation until we
			// reach the root of the statement they're defined in.
			//
			// This ensures that if there are any @lends annotations on the object
			// literals, the type on the @lends annotation resolves correctly.
			//
			// For more information, see http://blickly.github.io/closure-compiler-issues/#314
			action := func() { b.DefineObjectLiteral(n) }
			b.deferredActions.Put(NodeUtil_GetEnclosingStatement(n), &action)
		} else {
			b.DefineObjectLiteral(n)
		}
	case rhino.Token_CLASS:
		// NOTE(sdh): We can't handle function nodes here because they need special behavior to
		// deal with hoisting.  But since classes aren't hoisted, and may need to be handled in
		// such places as default method initializers (i.e. in a FunctionScope) or class extends
		// clauses (technically part of the ClassScope, but visited instead by the NormalScope),
		// they can be handled consistently in all scopes.
		b.DefineClassLiteral(n)
	// NOTE(johnlenz): If we ever support Array tuples,
	// we will need to handle them here as we do object literals
	// above.
	case rhino.Token_ARRAYLIT:
		n.SetJSType(b.outer.GetNativeType(rhino.JSTypeNative_ARRAY_TYPE))
	default:
	}
}

// Java: TypedScopeCreator.AbstractScopeBuilder.defineObjectLiteral(Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) DefineObjectLiteral(objectLit *rhino.Node) {
	// Handle the @lends annotation.
	var type_ rhino.JSType
	info := objectLit.GetJSDocInfo()
	if info != nil && info.HasLendsName() {
		lendsName := info.GetLendsName().GetRoot().GetString()
		lendsVar := b.currentScope.GetVar(lendsName)
		if lendsVar == nil {
			b.outer.Report(JSError_MakeWithNode(objectLit, TypedScopeCreator_UNKNOWN_LENDS, lendsName))
		} else {
			type_ = lendsVar.GetType()
			if type_ == nil {
				type_ = b.outer.unknownType
			}
			if !type_.IsSubtypeOf(b.outer.typeRegistry.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE)) {
				b.outer.Report(JSError_MakeWithNode(objectLit, TypedScopeCreator_LENDS_ON_NON_OBJECT, lendsName, type_.String()))
				type_ = nil
			} else {
				objectLit.SetJSType(type_)
			}
		}
	}

	info = NodeUtil_GetBestJSDocInfo(objectLit)
	createEnumType := info != nil && info.HasEnumParameterType()
	if createEnumType {
		lValue := NodeUtil_GetBestLValue(objectLit)
		lValueName := NodeUtil_GetBestLValueName(lValue)
		type_ = b.CreateEnumTypeFromNodes(objectLit, lValueName, lValue, info)
	}

	if type_ == nil {
		type_ = b.outer.typeRegistry.CreateAnonymousObjectType(info)
	}

	b.outer.SetDeferredType(objectLit, type_)

	// If this is an enum, the properties were already taken care of above.
	b.ProcessObjectLitProperties(
		objectLit, rhino.ObjectType_Cast(objectLit.GetJSType()), !createEnumType)
}

// ProcessObjectLitProperties processes an object literal and all the types
// on it.
//
// objLit is the OBJECTLIT node; objLitType the type of the OBJECTLIT node,
// which might be a named type, because of the lends annotation;
// declareOnOwner, if true, declares properties on the objLitType as well.
// If false, the caller should take care of this.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.processObjectLitProperties(Node, ObjectType, boolean)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) ProcessObjectLitProperties(objLit *rhino.Node, objLitType rhino.ObjectType, declareOnOwner bool) {
	for keyNode := objLit.GetFirstChild(); keyNode != nil; keyNode = keyNode.GetNext() {
		if keyNode.IsSpread() {
			// Don't try defining computed or spread properties on an object. Note that for spread
			// type inference will try to determine the properties and types. We cannot do it here as
			// we don't have all the type information of the spread object.
			continue
		}
		var memberName rhino.Property_Key
		var value *rhino.Node
		var qualifiedName *string
		if keyNode.IsComputedProp() {
			symbol := TypedScopeCreator_ExtractKnownSymbolKey(b.currentScope, keyNode.GetFirstChild())
			if symbol == nil {
				continue
			}
			qualifiedName = nil
			memberName = rhino.NewProperty_SymbolKey(symbol)
			value = keyNode.GetSecondChild()
		} else {
			value = keyNode.GetFirstChild()
			memberName = rhino.NewProperty_StringKey(NodeUtil_GetObjectOrClassLitKeyName(keyNode))
			qualifiedName = NodeUtil_GetBestLValueName(keyNode)
		}
		info := keyNode.GetJSDocInfo()
		valueType := b.GetDeclaredType(info, keyNode, value, nil)
		var keyType rhino.JSType
		if objLitType.IsEnumType() {
			keyType = objLitType.ToMaybeEnumType().GetElementsType()
		} else {
			keyType = TypeCheck_GetObjectLitKeyTypeFromValueType(keyNode, valueType)
		}

		// Try to declare this property in the current scope if it
		// has an authoritative name.
		if qualifiedName != nil {
			NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(b).
				ForDeclarationNode(keyNode).
				ForVariableName(*qualifiedName).
				InScope(b.GetLValueRootScope(keyNode)).
				WithType(keyType).
				AllowLaterTypeInference(keyType == nil).
				DefineSlot()
		} else if keyType != nil {
			b.outer.SetDeferredType(keyNode, keyType)
		}

		if keyType != nil && objLitType != nil && declareOnOwner {
			// Declare this property on its object literal.
			objLitType.DefineDeclaredPropertyByKey(memberName, keyType, keyNode)
		}
	}
}

// GetDeclaredTypeInAnnotation returns the type specified in a JSDoc
// annotation near a GETPROP, NAME, object literal member, or class field.
//
// Extracts type information from the @type tag.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.getDeclaredTypeInAnnotation(Node, JSDocInfo)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) GetDeclaredTypeInAnnotation(node *rhino.Node, info *rhino.JSDocInfo) rhino.JSType {
	jutil.CheckArgumentMsg1(info.HasType(), "%s", info)
	ownerTypeKeys := b.FindOwnerTypeKeys(node)

	var templateScope rhino.StaticTypedScope
	if len(ownerTypeKeys) != 0 {
		templateScope = b.outer.typeRegistry.CreateScopeWithTemplates(b.currentScope, ownerTypeKeys)
	} else {
		templateScope = b.currentScope
	}
	return info.GetType().Evaluate(templateScope, b.outer.typeRegistry)
}

// Java: TypedScopeCreator.AbstractScopeBuilder.findOwnerTypeKeys(Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) FindOwnerTypeKeys(node *rhino.Node) []*rhino.TemplateType {
	if node != nil && !node.IsStaticMember() && node.GetParent().IsClassMembers() {
		classType := node.GetGrandparent().GetJSType()
		if classType != nil {
			return classType.ToMaybeFunctionType().GetTemplateTypeMap().GetTemplateKeys()
		}
		return []*rhino.TemplateType{}
	}
	ownerNode := NodeUtil_GetBestLValueOwner(node)
	ownerName := NodeUtil_GetBestLValueName(ownerNode)
	if ownerName != nil {
		ownerVar := b.currentScope.GetVar(*ownerName)
		if ownerVar != nil {
			ownerType := b.GetPrototypeOwnerType(rhino.ObjectType_Cast(ownerVar.GetType()))
			if ownerType != nil {
				return ownerType.GetTemplateTypeMap().GetTemplateKeys()
			}
		}
	}
	return []*rhino.TemplateType{}
}

// AssertDefinitionNode asserts that it's OK to define this node's name. The
// node should have a source name and be of the specified type.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.assertDefinitionNode(Node, Token)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) AssertDefinitionNode(n *rhino.Node, type_ rhino.Token) {
	jutil.CheckState(b.sourceName != nil)
	jutil.CheckStateMsg1(n.GetToken() == type_, "%s", n)
}

// DefineCatch defines a catch parameter.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.defineCatch(Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) DefineCatch(n *rhino.Node) {
	b.AssertDefinitionNode(n, rhino.Token_CATCH)
	// Though almost certainly a terrible idea, it is possible to do destructuring in
	// the catch declaration.
	// e.g. `} catch ({message, errno}) {`
	NodeUtil_VisitLhsNodesInNode(
		n,
		func(catchName *rhino.Node) {
			type_ := b.GetDeclaredType(catchName.GetJSDocInfo(), catchName, nil, nil)
			NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(b).
				ForDeclarationNode(catchName).
				ForVariableName(catchName.GetString()).
				InScope(b.currentScope).
				WithType(type_).
				AllowLaterTypeInference(type_ == nil).
				DefineSlot()
		})
}

// DefineAssignAsIfVarDeclaration defines an assignment to a name as if it
// were an actual `var` declaration.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.defineAssignAsIfVarDeclaration(Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) DefineAssignAsIfVarDeclaration(assignment *rhino.Node) {
	info := assignment.GetJSDocInfo()
	name := assignment.GetFirstChild()
	jutil.CheckArgumentMsg1(name.IsName(), "%s", name)
	rvalue := assignment.GetSecondChild()
	b.DefineName(name, rvalue, b.currentScope.GetClosestHoistScope(), info)
}

// DefineVars defines a variable declared with `var`, `let`, or `const`.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.defineVars(Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) DefineVars(n *rhino.Node) {
	jutil.CheckState(b.sourceName != nil)
	jutil.CheckState(NodeUtil_IsNameDeclaration(n))
	info := n.GetJSDocInfo()
	// `var` declarations are hoisted, but `let` and `const` are not.
	scope := b.currentScope
	if n.IsVar() {
		scope = b.currentHoistScope
	}

	if n.HasMoreThanOneChild() && info != nil {
		b.outer.Report(JSError_MakeWithNode(n, TypeCheck_MULTIPLE_VAR_DEF))
	}

	for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
		b.DefineVarChild(info, child, scope)
	}
	if n.HasOneChild() && b.IsValidTypedefDeclaration(n.GetOnlyChild(), n.GetJSDocInfo()) {
		b.DeclareTypedefType(n.GetOnlyChild(), n.GetJSDocInfo())
	}
}

// DefineVarChild defines a variable declared with `var`, `let`, or `const`.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.defineVarChild(JSDocInfo, Node, TypedScope)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) DefineVarChild(declarationInfo *rhino.JSDocInfo, child *rhino.Node, scope *TypedScope) {
	if child.IsName() {
		if declarationInfo == nil {
			declarationInfo = child.GetJSDocInfo()
			// TODO(bradfordcsmith): Report an error if both the declaration node and the name itself
			//     have JSDoc.
		}
		b.DefineName(child, child.GetFirstChild(), scope, declarationInfo)
	} else {
		jutil.CheckStateMsg1(child.IsDestructuringLhs(), "%s", child)
		pattern := child.GetFirstChild()
		value := child.GetSecondChild()

		if ModuleImportResolver_IsGoogModuleDependencyCall(value) {
			// Define destructuring names here, since goog.require destructuring patterns can only
			// have one level and require some special handling.
			defaultImport := b.outer.moduleImportResolver.GetClosureNamespaceTypeFromCall(value)
			for key := pattern.GetFirstChild(); key != nil; key = key.GetNext() {
				b.DefineModuleImport(key.GetFirstChild(), defaultImport, jutil.StringPtr(key.GetString()), scope)
			}
			return
		}

		b.DefineDestructuringPatternInVarDeclaration(
			pattern,
			scope,
			func() *TypedScopeCreator_RValueInfo {
				// Note that value will be null if we are in an enhanced for loop
				//   for (const {x, y} of data) {
				if value != nil {
					return NewTypedScopeCreator_RValueInfo(
						b.GetDeclaredRValueType( /* lValue= */ nil, value),
						value.GetQualifiedNameObject())
				}
				return NewTypedScopeCreator_RValueInfo(b.outer.unknownType /* qualifiedName= */, nil)
			})
	}
}

// InferTypeForDestructuredTarget returns information about the qualified
// name and type of the target, if it exists.
//
// Never returns null, but will return an RValueInfo with null `type` and
// `qualifiedName` slots.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.inferTypeForDestructuredTarget(DestructuredTarget, Supplier)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) InferTypeForDestructuredTarget(
	target *DestructuredTarget, patternTypeSupplier func() *TypedScopeCreator_RValueInfo) *TypedScopeCreator_RValueInfo {
	// Currently we only do type inference for string key nodes in object patterns here, to
	// handle aliasing types. e.g
	//   const {Foo} = ns;
	// TypeInference takes care of the rest.
	// Note that although DestructuredTarget includes logic for inferring types, we don't use
	// it here because we only do some very limited type inference during TypedScopeCreation,
	// and only return a non-null type here if we are accessing a declared property on a known
	// type.
	if !target.HasStringKey() || target.HasDefaultValue() {
		return TypedScopeCreator_RValueInfo_Empty()
	}
	rvalue := patternTypeSupplier()
	patternType := rvalue.type_
	propertyName := target.GetStringKey().GetString()
	var qname rhino.QualifiedName
	if rvalue.qualifiedName != nil {
		qname = rvalue.qualifiedName.Getprop(propertyName)
	}
	if patternType == nil || patternType.IsUnknownType() {
		return NewTypedScopeCreator_RValueInfo(nil, qname)
	}
	if patternType.HasProperty(propertyName) {
		type_ := patternType.FindPropertyType(propertyName)
		return NewTypedScopeCreator_RValueInfo(type_, qname)
	}
	return NewTypedScopeCreator_RValueInfo(nil, qname)
}

// Java: TypedScopeCreator.AbstractScopeBuilder.defineDestructuringPatternInVarDeclaration(Node, TypedScope, Supplier)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) DefineDestructuringPatternInVarDeclaration(
	pattern *rhino.Node, scope *TypedScope, patternTypeSupplier func() *TypedScopeCreator_RValueInfo) {
	for _, target := range DestructuredTarget_CreateAllNonEmptyTargetsInPatternWithSupplier(
		b.outer.typeRegistry, func() rhino.JSType { return patternTypeSupplier().type_ }, pattern) {

		typeSupplier := func() *TypedScopeCreator_RValueInfo {
			return b.InferTypeForDestructuredTarget(target, patternTypeSupplier)
		}

		if target.GetNode().IsDestructuringPattern() {
			b.DefineDestructuringPatternInVarDeclaration(target.GetNode(), scope, typeSupplier)
		} else {
			name := target.GetNode()
			jutil.CheckStateMsg1(name.IsName(), "This method is only for declaring variables: %s", name)

			// variable's type
			type_ := b.GetDeclaredType(name.GetJSDocInfo(), name /* rValue= */, nil, typeSupplier)
			if type_ == nil {
				// The variable's type will be inferred.
				if name.IsFromExterns() {
					type_ = b.outer.unknownType
				} else {
					type_ = nil
				}
			}
			NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(b).
				ForDeclarationNode(name).
				ForVariableName(name.GetString()).
				InScope(scope).
				WithType(type_).
				AllowLaterTypeInference(type_ == nil).
				DefineSlot()
		}
	}
}

// DefineClassLiteral defines a class literal. Handles any of the following
// cases:
//
//   - Class declarations: `class Foo { ... }`
//   - Class assignments: `foo.Bar = class { ... }`
//   - Bleeding names: `foo.Bar = class Baz { ... }`
//   - Properties: `{foo: class { ... }}`
//   - Callbacks: `foo(class { ... })`
//
// Java: TypedScopeCreator.AbstractScopeBuilder.defineClassLiteral(Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) DefineClassLiteral(n *rhino.Node) {
	b.AssertDefinitionNode(n, rhino.Token_CLASS)

	// Determine the name and JSDocInfo and l-value for the class.
	// Any of these may be null.
	lValue := NodeUtil_GetBestLValue(n)
	info := NodeUtil_GetBestJSDocInfo(n)
	className := NodeUtil_GetBestLValueName(lValue)
	classTypeIdentifier := b.GetBestTypeName(lValue, className)

	// Create the type and assign it on the CLASS node.
	classType := b.CreateClassTypeFromNodes(n, classTypeIdentifier, className, info, lValue)
	b.outer.SetDeferredType(n, classType)

	// Declare this symbol in the current scope iff it's a class
	// declaration. Otherwise, the declaration will happen in other
	// code paths.
	if NodeUtil_IsClassDeclaration(n) {
		jutil.CheckNotNull(className)
		NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(b).
			ForDeclarationNode(n.GetFirstChild()).
			ForVariableName(*className).
			InScope(b.currentScope).
			WithType(classType).
			AllowLaterTypeInference(classType == nil).
			DefineSlot()
	}
}

// Java: TypedScopeCreator.AbstractScopeBuilder.getBestTypeName(Node, String)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) GetBestTypeName(lvalue *rhino.Node, syntacticLvalueName *string) *string {
	if syntacticLvalueName == nil {
		return nil
	}
	if b.IsGoogModuleExports(lvalue) {
		return jutil.StringPtr(jutil.StringReplace(*syntacticLvalueName, "exports", *b.GetModule().ClosureNamespace()))
	}
	return syntacticLvalueName
}

// DefineFunctionLiteral defines a function literal.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.defineFunctionLiteral(Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) DefineFunctionLiteral(n *rhino.Node) {
	b.AssertDefinitionNode(n, rhino.Token_FUNCTION)

	// Determine the name and JSDocInfo and l-value for the function.
	// Any of these may be null.
	lValue := NodeUtil_GetBestLValue(n)
	info := NodeUtil_GetBestJSDocInfo(n)
	functionName := NodeUtil_GetBestLValueName(lValue)
	functionType := b.CreateFunctionTypeFromNodes(n, functionName, info, lValue)

	// Assigning the function type to the function node
	b.outer.SetDeferredType(n, functionType)

	// Declare this symbol in the current scope iff it's a function
	// declaration. Otherwise, the declaration will happen in other
	// code paths.
	if NodeUtil_IsFunctionDeclaration(n) {
		NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(b).
			ForDeclarationNode(n.GetFirstChild()).
			ForVariableName(*functionName).
			InScope(b.currentScope).
			WithType(functionType).
			AllowLaterTypeInference(functionType == nil).
			DefineSlot()
	}
}

// DefineName defines a variable based on the NAME node passed.
//
// name is the NAME node; value optionally the value assigned to the name
// node; info the JSDocInfo information relating to this name node.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.defineName(Node, Node, TypedScope, JSDocInfo)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) DefineName(name *rhino.Node, value *rhino.Node, scope *TypedScope, info *rhino.JSDocInfo) {
	if ModuleImportResolver_IsGoogModuleDependencyCall(value) {
		b.DefineModuleImport(
			name, b.outer.moduleImportResolver.GetClosureNamespaceTypeFromCall(value), nil, scope)
		return
	}
	type_ := b.GetDeclaredType(info, name, value /* declaredRValueTypeSupplier= */, nil)
	if type_ == nil {
		// The variable's type will be inferred.
		if name.IsFromExterns() {
			type_ = b.outer.unknownType
		} else {
			type_ = nil
		}
	}
	NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(b).
		ForDeclarationNode(name).
		ForVariableName(name.GetString()).
		InScope(scope).
		WithType(type_).
		AllowLaterTypeInference(type_ == nil).
		DefineSlot()
}

// DefineModuleImport defines the name of an import.
//
// localNameNode is the name node of the LHS of the import being defined;
// importedModuleObject the root node of the scope in which the type being
// imported was defined, along with the local name of the overall module
// object inside the scope where it was defined, or null if no module with
// that name exists; optionalProperty the property name of the locally
// imported type on the module object, if destructuring-style importing was
// used, or null if this is a namespace import; scopeToDeclareIn the scope in
// which localNameNode is defined.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.defineModuleImport(Node, ScopedName, String, TypedScope)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) DefineModuleImport(
	localNameNode *rhino.Node,
	importedModuleObject ScopedName,
	optionalProperty *string,
	scopeToDeclareIn *TypedScope) {
	if importedModuleObject == nil {
		// We could not find the module defining this import. If the compiler is assuming forward
		// declared types (or if explicitly forward declared), create a NamedType so it resolves to
		// NoResolvedType. Otherwise, fall back to unknownType to avoid spurious type check errors.
		var fallbackType rhino.JSType
		if b.outer.compiler.GetOptions().AssumeForwardDeclaredForMissingTypes() {
			fallbackType = b.outer.typeRegistry.CreateNamedType(
				b.currentScope,
				localNameNode.GetString(),
				localNameNode.GetSourceFileName(),
				localNameNode.GetLineno(),
				localNameNode.GetCharno())
		} else {
			fallbackType = b.outer.unknownType
		}
		NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(b).
			ForDeclarationNode(localNameNode).
			ForVariableName(localNameNode.GetString()).
			InScope(scopeToDeclareIn).
			WithType(fallbackType).
			AllowLaterTypeInference(
				fallbackType.IsUnknownType() &&
					!b.outer.compiler.GetOptions().AssumeForwardDeclaredForMissingTypes()).
			DefineSlot()
		return
	}

	var exportedName ScopedName
	if optionalProperty == nil {
		exportedName = importedModuleObject
	} else {
		// If this is a destucuring-style import, find its type.
		exportedName = ScopedName_Of(
			importedModuleObject.GetName()+"."+*optionalProperty,
			importedModuleObject.GetScopeRoot())
	}
	// Try getting the actual scope. The scope will be null in the following cases:
	//   - Someone has required a module that does not exist at all.
	//   - Someone has requireType'd or forwardDeclare'd a module that exists, but does not have
	//     an associated scope yet.
	var exportScope *TypedScope
	if exportedName.GetScopeRoot() != nil {
		exportScope = b.outer.memoized.Get(exportedName.GetScopeRoot())
	}

	// The scope is null for modules that were not visited yet.
	if exportScope != nil {
		type_ := exportScope.LookupQualifiedName(rhino.QualifiedName_Of(exportedName.GetName()))

		// The type is null if either the name is inferred or this is an early ref.
		if type_ != nil {
			b.DeclareAliasTypeIfRvalueIsAliasable(
				localNameNode, rhino.QualifiedName_Of(exportedName.GetName()), type_, exportScope)

			NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(b).
				ForDeclarationNode(localNameNode).
				ForVariableName(localNameNode.GetString()).
				InScope(scopeToDeclareIn).
				WithType(type_).
				AllowLaterTypeInference(type_ == nil).
				DefineSlot()
			return
		}
	}
	// Defer defining this name until after we have visited the entire AST.
	b.outer.weakImports = append(b.outer.weakImports,
		NewTypedScopeCreator_WeakModuleImport(b.outer, localNameNode, exportedName, scopeToDeclareIn))
}

// ShouldUseFunctionLiteralType: if a variable is assigned a function
// literal in the global scope, make that a declared type (even if there's
// no doc info). There's only one exception to this rule: if the return type
// is inferred, and we're in a local scope, we should assume the whole
// function is inferred.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.shouldUseFunctionLiteralType(FunctionType, JSDocInfo, Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) ShouldUseFunctionLiteralType(type_ rhino.FunctionType, info *rhino.JSDocInfo, lValue *rhino.Node) bool {
	if info != nil {
		return true
	}
	if lValue != nil && NodeUtil_MayBeObjectLitKey(lValue) {
		return false
	}
	// TODO(johnlenz): consider unifying global and local behavior
	return b.IsLValueRootedInGlobalScope(lValue) || !type_.IsReturnTypeInferred()
}

// CreateClassTypeFromNodes creates a new class type from the given class
// literal. This function does not need to worry about stubs and aliases
// because they are handled by createFunctionTypeFromNodes instead.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.createClassTypeFromNodes(Node, String, String, JSDocInfo, Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) CreateClassTypeFromNodes(
	clazz *rhino.Node,
	name *string,
	syntacticName *string,
	info *rhino.JSDocInfo,
	lvalueNode *rhino.Node) rhino.FunctionType {
	jutil.CheckArgumentMsg1(clazz.IsClass(), "%s", clazz)

	declarationScope := b.currentScope
	if lvalueNode != nil {
		declarationScope = b.GetLValueRootScope(lvalueNode)
	}
	builder := NewFunctionTypeBuilder(name, b.outer.compiler, clazz, b.currentScope).
		UsingClassSyntax().
		SetSyntacticFunctionName(syntacticName).
		SetContents(NewFunctionTypeBuilder_AstFunctionContents(clazz)).
		SetDeclarationScope(declarationScope).
		InferKind(info).
		InferTemplateTypeName(info, nil)

	extendsClause := clazz.GetSecondChild()

	// Look at the extends clause and/or JSDoc info to find a super class.  Use generics from the
	// JSDoc to supplement the extends type when available.
	baseType := b.FindSuperClassFromNodes(extendsClause, info)
	builder.InferInheritance(info, baseType)

	// Look for an explicit constructor.
	constructor := NodeUtil_GetEs6ClassConstructorMemberFunctionDef(clazz)
	if constructor != nil {
		constructor = constructor.GetOnlyChild() // We want the FUNCTION, not the member.
	}

	if constructor != nil {
		// Note: constructor should have the following structure:
		//   MEMBER_FUNCTION_DEF [jsdoc_info]
		//     FUNCTION
		//       NAME
		//       PARAM_LIST ...
		//       BLOCK ...
		// NodeUtil.getFirstPropMatchingKey returns the FUNCTION node.
		ctorInfo := NodeUtil_GetBestJSDocInfo(constructor)
		builder.InferConstructorParameters(constructor.GetSecondChild(), ctorInfo)
	} else if extendsClause.IsEmpty() {
		// No explicit constructor and no superclass: constructor is no-args.
		builder.InferImplicitConstructorParameters([]*rhino.FunctionType_Parameter{})
	} else {
		// No explicit constructor, but we have a superclass.  If we know its type, then copy its
		// constructor arguments (and templates).  If not, make the constructor arguments unknown.
		// TODO(sdh): consider allowing attaching constructor @param annotations somewhere else?
		var extendsCtor rhino.FunctionType
		if baseType != nil {
			extendsCtor = baseType.GetConstructor()
		}
		if extendsCtor != nil {
			// Known superclass: copy the parameters node.
			builder.InferImplicitConstructorParameters(extendsCtor.GetParameters())
		} else {
			// Unresolveable extends clause: suppress typechecking.
			builder.InferImplicitConstructorParameters(
				b.outer.typeRegistry.CreateParametersWithVarArgs(
					b.outer.typeRegistry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)))
		}
	}

	// TODO(sdh): Handle template parameters.  The constructor should store all parameters,
	// while the instance type should only have the class-level parameters?

	// Add the type for the "constructor" property.
	classType := builder.BuildAndRegister()
	if classType.IsConstructor() {
		// NOTE: This logic is similar to the goog.inherits handling in
		// ClosureCodingConvention#applySubclassRelationship. If this logic is modified
		// it is likely that code needs to be modified as well.

		// Notice that constructor functions do not need to be covariant on the superclass.
		// So if G extends F, new G() and new F() can accept completely different argument
		// types, but G.prototype.constructor needs to be covariant on F.prototype.constructor.
		// To get around this, we just turn off type-checking on arguments and return types
		// of G.prototype.constructor.

		// NOTE: For final classes we could do better here and retain the parameter types.

		qmarkCtor := classType.ForgetParameterAndReturnTypes()
		classPrototypeType := classType.GetPrototypeProperty()
		classPrototypeType.DefineDeclaredProperty("constructor", qmarkCtor, constructor)
	}
	if classType.HasInstanceType() {
		classPrototype := classType.GetSlot("prototype")
		// SymbolTable users expect the class prototype and actual class to have the same
		// declaration node.
		if lvalueNode != nil {
			classPrototype.SetNode(lvalueNode)
		} else {
			classPrototype.SetNode(classPrototype.GetNode())
		}
	}

	return classType
}

// FindSuperClassFromNodes looks at the extends clause to find the instance
// type being extended. Returns null if there is no such clause, and unknown
// if the type cannot be determined.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.findSuperClassFromNodes(Node, JSDocInfo)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) FindSuperClassFromNodes(extendsNode *rhino.Node, info *rhino.JSDocInfo) rhino.ObjectType {
	if extendsNode.IsEmpty() {
		// No extends clause: return null.
		return nil
	}
	ctorType := extendsNode.GetJSType()
	if ctorType == nil {
		if extendsNode.IsQualifiedName() {
			superclass := *extendsNode.GetQualifiedName()
			// Look up qualified names in the scope (types won't be set on the AST until inference).
			ctorType = b.currentScope.LookupQualifiedName(rhino.QualifiedName_Of(superclass))
			if ctorType == nil {
				var_ := b.currentScope.GetVar(superclass)
				if var_ != nil {
					ctorType = var_.GetType()
				} else {
					ctorType = nil
				}
			}
			// If that doesn't work, then fall back on the registry
			if ctorType == nil {
				return rhino.ObjectType_Cast(
					b.outer.typeRegistry.GetTypeWithLocation(
						b.currentScope,
						superclass,
						extendsNode.GetSourceFileName(),
						extendsNode.GetLineno(),
						extendsNode.GetCharno()))
			}
		} else if b.IsGoogModuleGetProperty(extendsNode) {
			ctorType = b.GetCtorForGoogModuleGet(extendsNode)
		} else {
			// Anything TypedScopeCreator can infer has already been read off the AST.  This is likely
			// a CALL or GETELEM, which are unknown until TypeInference.  Instead, ignore it for now,
			// require an @extends tag in the JSDoc, and verify correctness in TypeCheck.
			if info == nil || !info.HasBaseType() {
				b.outer.Report(JSError_MakeWithNode(extendsNode, TypedScopeCreator_DYNAMIC_EXTENDS_WITHOUT_JSDOC))
			}
		}
	}

	if ctorType != nil {
		if ctorType.IsConstructor() || ctorType.IsInterface() {
			return ctorType.ToMaybeFunctionType().GetInstanceType()
		} else if ctorType.IsUnknownType() || ctorType.IsNoResolvedType() {
			// The constructor could have an unknown type for cases where it is dynamically
			// created or passed in from elsewhere.
			// e.g. with a mixin pattern
			// function mixinSomething(ctor) {
			//   return class extends ctor { ... };
			// }
			// In that case consider the super class instance type to be unknown.
			// On the other hand, the constructor can be a NoResolvedType if running in Clutz's
			// partial compilation mode.
			return ctorType.ToMaybeObjectType()
		}
	}

	// We couldn't determine the type, so for TypedScope creation purposes we will treat it as if
	// there were no extends clause.  TypeCheck will give a more precise error later.
	return nil
}

// Java: TypedScopeCreator.AbstractScopeBuilder.isGoogModuleGetProperty(Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) IsGoogModuleGetProperty(extendsClause *rhino.Node) bool {
	if extendsClause.IsCall() {
		return ModuleImportResolver_IsGoogModuleDependencyCall(extendsClause)
	} else if extendsClause.IsGetProp() {
		return b.IsGoogModuleGetProperty(extendsClause.GetFirstChild())
	} else {
		return false
	}
}

// GetCtorForGoogModuleGet looks up the type of a goog.module.get property
// access.
//
// extendsNode is `goog.module.get('x');` or a getprop
// `goog.module.get('x').y.z`.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.getCtorForGoogModuleGet(Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) GetCtorForGoogModuleGet(extendsNode *rhino.Node) rhino.JSType {
	call := extendsNode
	properties := []string{}
	for !call.IsCall() {
		properties = slices.Insert(properties, 0, call.GetString())
		call = call.GetFirstChild()
	}
	extendsCall := b.outer.moduleImportResolver.GetClosureNamespaceTypeFromCall(call)
	if extendsCall == nil || !b.outer.memoized.ContainsKey(extendsCall.GetScopeRoot()) {
		// Handle invalid goog.module.get calls
		return nil
	}
	moduleScope := b.outer.memoized.Get(extendsCall.GetScopeRoot())
	properties = slices.Insert(properties, 0, extendsCall.GetName())
	superclassName := rhino.QualifiedName_Of(strings.Join(properties, "."))
	return moduleScope.LookupQualifiedName(superclassName)
}

// CreateFunctionTypeFromNodes creates a new function type, based on the
// given nodes.
//
// This handles two cases that are semantically very different, but are not
// mutually exclusive: - A function literal that needs a type attached to it
// (called from defineClassLiteral with a non-null FUNCTION node for rValue).
// - An assignment expression with function-type info in the JsDoc (called
// from getDeclaredType on a stub (rValue == null) or alias (rValue is a
// qualified name).
//
// All parameters are optional, and we will do the best we can to create a
// function type.
//
// This function will always create a function type, so only call it if
// you're sure that's what you want.
//
// rValue is the function node; name the function's name; info the JSDocInfo
// attached to the function definition; lvalueNode the node where this
// function is being assigned. For example, `A.prototype.foo = ...` would be
// used to determine that this function is a method of A.prototype. May be
// null to indicate that this is not being assigned to a qualified name.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.createFunctionTypeFromNodes(Node, String, JSDocInfo, Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) CreateFunctionTypeFromNodes(
	rValue *rhino.Node,
	name *string,
	info *rhino.JSDocInfo,
	lvalueNode *rhino.Node) rhino.FunctionType {
	// Check for an alias.
	if rValue != nil && rValue.IsQualifiedName() && lvalueNode != nil {
		var_ := b.currentScope.GetVar(*rValue.GetQualifiedName())
		if var_ != nil && var_.GetType() != nil && var_.GetType().IsFunctionType() {
			aliasedType := var_.GetType().ToMaybeFunctionType()
			if aliasedType.IsConstructor() || aliasedType.IsInterface() {
				// TODO(nick): Remove this. This should already be handled by normal type resolution.
				if name != nil {
					b.outer.typeRegistry.DeclareType(b.currentScope, *name, aliasedType.GetInstanceType())
				}
				b.CheckFunctionAliasAnnotations(lvalueNode, aliasedType, info)
				return aliasedType
			}
		}
	}

	// No alias: look for an explicit @type in JSDocInfo.
	if info != nil && info.HasType() {
		type_ := info.GetType().Evaluate(b.currentScope, b.outer.typeRegistry)

		// Known to be not null since we have the FUNCTION token there.
		type_ = type_.RestrictByNotNullOrUndefined()
		if type_.IsFunctionType() {
			functionType := type_.ToMaybeFunctionType()
			functionType.SetJSDocInfo(info)
			return functionType
		}
	}

	// No alias or explicit @type, so look for a function literal, or @param/@return.
	errorRoot := rValue
	if rValue == nil {
		errorRoot = lvalueNode
	}
	isFnLiteral := rValue != nil && rValue.IsFunction()
	var fnRoot *rhino.Node
	var parametersNode *rhino.Node
	if isFnLiteral {
		fnRoot = rValue
		parametersNode = rValue.GetSecondChild()
	}

	// If this function is being assigned as a property on a type, try finding the owner type
	// and the property name.
	// This is easy to do for class members because the owner type is on the CLASS node and
	// the property name is the MEMBER_FUNCTION_DEF/GETTER_DEF/SETTER_DEF string.
	// For other functions, we rely on NodeUtil.getBestLValueOwner.
	var classRoot *rhino.Node
	if lvalueNode != nil && lvalueNode.GetParent().IsClassMembers() {
		classRoot = lvalueNode.GetGrandparent()
	}
	ownerNode := NodeUtil_GetBestLValueOwner(lvalueNode)

	var ownerType rhino.ObjectType
	var propName *string
	if classRoot != nil {
		// Static members are owned by the constructor, non-statics are owned by the prototype.
		if ownerFn := rhino.JSType_ToMaybeFunctionType(classRoot.GetJSType()); ownerFn != nil {
			ownerType = ownerFn
		}
		if !lvalueNode.IsStaticMember() && ownerType != nil {
			ownerType = ownerType.(rhino.FunctionType).GetPrototype()
		}
		if !lvalueNode.IsComputedProp() {
			propName = jutil.StringPtr(lvalueNode.GetString())
		}
	} else {
		ownerName := NodeUtil_GetBestLValueName(ownerNode)
		if ownerName != nil {
			ownerType =
				rhino.ObjectType_Cast(b.currentScope.LookupQualifiedName(rhino.QualifiedName_Of(*ownerName)))
		}

		if ownerName != nil && name != nil {
			// TODO(b/111621092): Use the AST rather than manipulating strings here.
			jutil.CheckStateMsg2(
				jutil.StringStartsWith(*name, *ownerName), "Expected \"%s\" to start with \"%s\"", *name, *ownerName)
			propName = jutil.StringPtr(jutil.StringSubstringFrom(*name, jutil.StringLength(*ownerName)+1))
		}
	}

	prototypeOwner := b.GetPrototypeOwnerType(ownerType)
	var prototypeOwnerTypeMap *rhino.TemplateTypeMap
	if prototypeOwner != nil && prototypeOwner.GetTypeOfThis() != nil {
		prototypeOwnerTypeMap = prototypeOwner.GetTypeOfThis().GetTemplateTypeMap()
	}

	// Find the type of any overridden function.
	var overriddenType rhino.FunctionType
	if ownerType != nil && propName != nil {
		// the type of the property this overrides, not necessarily a function.
		overriddenPropType :=
			b.FindOverriddenProperty(ownerType, *propName, prototypeOwnerTypeMap)
		if overriddenPropType != nil {
			// Overridden getters and setters need special handling because we declare
			// getters/setters as simple properties with their respective return/parameter type. This
			// causes a split during inference where left and right sides of a getter/setter
			// declaration will be inferred to have different types; if the left side has type `T`,
			// the right side will be some function type involving `T`.
			if lvalueNode.IsGetterDef() {
				// Convert `number` to `function(): number`
				overriddenType = b.outer.typeRegistry.CreateFunctionType(overriddenPropType)
			} else if lvalueNode.IsSetterDef() {
				// Convert `number` to `function(number): undefined`
				overriddenType =
					b.outer.typeRegistry.CreateFunctionType(b.outer.GetNativeType(rhino.JSTypeNative_VOID_TYPE), overriddenPropType)
			} else if overriddenPropType.IsFunctionType() {
				// for cases where we override a non-method (e.g. a number) with a method, don't put the
				// non-method type (e.g. number) on the function.
				// Instead do some basic inference to create a function type.
				// we will warn during typechecking for an invalid override, but we don't want to put a
				// non-function type on this function because that will interfere with type inference
				// inside the function.
				overriddenType = overriddenPropType.ToMaybeFunctionType()
			}
		}
	}

	var contents *FunctionTypeBuilder_AstFunctionContents
	if fnRoot != nil {
		contents = NewFunctionTypeBuilder_AstFunctionContents(fnRoot)
	}
	if b.outer.functionsWithNonEmptyReturns.Contains(fnRoot) {
		contents.RecordNonEmptyReturn()
	}
	// Go-only: language-adapter — a nil *AstFunctionContents must be passed
	// as a nil FunctionContents interface (Java: setContents(null)).
	var contentsOrNil FunctionTypeBuilder_FunctionContents
	if contents != nil {
		contentsOrNil = contents
	}

	bestTypeName := b.GetBestTypeName(lvalueNode, name)
	declarationScope := b.currentScope
	if lvalueNode != nil {
		declarationScope = b.GetLValueRootScope(lvalueNode)
	}
	// Go-only: language-adapter — a nil ObjectType must be passed as a nil
	// JSType interface.
	var prototypeOwnerOrNil rhino.JSType
	if prototypeOwner != nil {
		prototypeOwnerOrNil = prototypeOwner
	}
	builder := NewFunctionTypeBuilder(bestTypeName, b.outer.compiler, errorRoot, b.currentScope).
		SetSyntacticFunctionName(name).
		SetContents(contentsOrNil).
		SetDeclarationScope(declarationScope).
		InferFromOverriddenFunction(overriddenType, parametersNode).
		InferKind(info).
		InferClosurePrimitive(info).
		InferTemplateTypeName(info, prototypeOwnerOrNil).
		InferInheritance(info, nil)

	if info == nil || !info.HasReturnType() {
		// when there is no @return annotation, look for inline return type declaration
		if rValue != nil && rValue.IsFunction() && rValue.HasChildren() {
			nameDocInfo := rValue.GetFirstChild().GetJSDocInfo()
			builder.InferReturnType(nameDocInfo, true)
		}
	} else {
		builder.InferReturnType(info, false)
	}

	// Infer the context type.
	var fallbackReceiverType rhino.JSType
	if ownerType != nil &&
		ownerType.IsFunctionPrototypeType() &&
		ownerType.GetOwnerFunction().HasInstanceType() {
		fallbackReceiverType = ownerType.GetOwnerFunction().GetInstanceType()
	} else if ownerType != nil &&
		ownerType.IsFunctionType() &&
		ownerType.ToMaybeFunctionType().HasInstanceType() &&
		lvalueNode != nil &&
		lvalueNode.IsStaticMember() {
		// Limit this case to members of ctors and interfaces decalared using `static`. Most
		// namespaces, like object literals, are assumed to declare free functions, so we exclude
		// them. Additionally, methods *assigned* to a ctor, especially an ES5 ctor, were never
		// designed with static polymorphism in mind, so excluding them preserves their assumptions.
		fallbackReceiverType = ownerType
	} else if ownerNode != nil && ownerNode.IsThis() {
		fallbackReceiverType = b.currentScope.GetTypeOfThis()
	}

	fnType := builder.
		InferThisTypeWithFallback(info, fallbackReceiverType).
		InferParameterTypesWithParamsParent(parametersNode, info).
		BuildAndRegister()

	// Do some additional validation for constructors and interfaces.
	if fnType.HasInstanceType() && lvalueNode != nil {
		prototypeSlot := fnType.GetSlot("prototype")

		// We want to make sure that the function and its prototype are declared at the same node.
		// This consistency is helpful to users of SymbolTable, because everything gets declared at
		// the same place.
		prototypeSlot.SetNode(lvalueNode)
	}
	return fnType
}

// CheckFunctionAliasAnnotations checks that the annotations in info are
// compatible with the aliased type. Any errors will be reported at n, which
// should be the qualified name node.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.checkFunctionAliasAnnotations(Node, FunctionType, JSDocInfo)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) CheckFunctionAliasAnnotations(n *rhino.Node, type_ rhino.FunctionType, info *rhino.JSDocInfo) {
	if info == nil {
		return
	}
	var annotation *string
	if info.UsesImplicitMatch() {
		if !type_.IsStructuralInterface() {
			annotation = jutil.StringPtr("@record")
		}
	} else if info.IsInterface() {
		if !type_.IsInterface() {
			annotation = jutil.StringPtr("@interface")
		}
	} else if info.IsConstructor() && !type_.IsConstructor() {
		annotation = jutil.StringPtr("@constructor")
	}
	// TODO(sdh): consider checking @template, @param, @return, and/or @this.
	if annotation != nil &&
		// TODO(sdh): Remove this extra check once TypeScript stops passing us duplicate
		// conflicting externs.  In particular, TS considers everything an interface, but Closure
		// externs mark most things as @constructor.  The load order is not always the same, so
		// the error can show up in either the generated TS externs file or in our own extern.
		(!n.IsFromExterns() || *annotation == "@record") {
		b.outer.Report(JSError_MakeWithNode(n, TypedScopeCreator_INCOMPATIBLE_ALIAS_ANNOTATION, *annotation, jutil.ToJavaString(n.GetQualifiedName())))
	}
}

// Java: TypedScopeCreator.AbstractScopeBuilder.getPrototypeOwnerType(ObjectType)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) GetPrototypeOwnerType(ownerType rhino.ObjectType) rhino.ObjectType {
	if ownerType != nil && ownerType.IsFunctionPrototypeType() {
		if ownerFunction := ownerType.GetOwnerFunction(); ownerFunction != nil {
			return ownerFunction
		}
		return nil
	}
	return nil
}

// FindOverriddenProperty finds the property that's being overridden on this
// type, if any.
//
// Said property could be a method, field, getter, or setter. We don't
// distinguish between these when looking up a property type.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.findOverriddenProperty(ObjectType, String, TemplateTypeMap)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) FindOverriddenProperty(
	ownerType rhino.ObjectType, propName string, typeMap *rhino.TemplateTypeMap) rhino.JSType {
	var result rhino.JSType

	// First, check to see if the property is implemented
	// on a superclass.
	propType := ownerType.GetPropertyType(propName)
	if propType != nil && !propType.IsUnknownType() {
		result = propType
	} else {
		// If it's not, then check to see if it's implemented
		// on an implemented interface.
		for _, iface := range ownerType.GetCtorImplementedInterfaces() {
			propType = iface.GetPropertyType(propName)
			if propType != nil && !propType.IsUnknownType() {
				result = propType
				break
			}
		}
	}

	if result != nil && typeMap != nil && !typeMap.IsEmpty() {
		result = result.Visit(rhino.TemplateTypeReplacer_ForPartialReplacement(b.outer.typeRegistry, typeMap)).(rhino.JSType)
	}

	return result
}

// CreateEnumTypeFromNodes creates a new enum type, based on the given nodes.
//
// This handles two cases that are semantically very different, but are not
// mutually exclusive: (1) An object literal that needs an enum type attached
// to it. (2) An assignment expression with an enum tag in the JsDoc.
//
// This function will always create an enum type, so only call it if you're
// sure that's what you want.
//
// rValue is the right-hand side of the enum, or null if none; lValue the
// left-hand side of the enum; name the qualified name of the enum; info the
// JSDocInfo attached to the enum definition.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.createEnumTypeFromNodes(Node, String, Node, JSDocInfo)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) CreateEnumTypeFromNodes(
	rValue *rhino.Node, name *string, lValue *rhino.Node, info *rhino.JSDocInfo) *rhino.EnumType {
	jutil.CheckNotNull(info)
	jutil.CheckState(info.HasEnumParameterType())
	jutil.CheckStateMsg(
		lValue != nil || rValue != nil,
		"An enum initializer should come from either an lvalue or rvalue")

	var enumType *rhino.EnumType
	if rValue != nil && rValue.IsQualifiedName() {
		// Handle an aliased enum. Note that putting @enum on an enum alias is optional. If the
		// rValue is not an enum, then this assignment errors during TypeCheck.
		var_ := b.currentScope.GetVar(*rValue.GetQualifiedName())
		if var_ != nil && var_.GetType() != nil && var_.GetType().IsEnumType() {
			enumType = var_.GetType().ToMaybeEnumType()
		}
	}

	if enumType == nil {
		enumType = rhino.EnumType_Builder_(b.outer.typeRegistry).
			SetName(b.GetBestTypeName(lValue, name)).
			SetGoogModuleId(TypedScopeCreator_ContainingGoogModuleIdOf(b.currentScope)).
			SetSource(rValue).
			SetElementType(info.GetEnumParameterType().Evaluate(b.currentScope, b.outer.typeRegistry)).
			Build()

		if rValue != nil && rValue.IsObjectLit() {
			// collect enum elements
			key := rValue.GetFirstChild()
			for key != nil {
				if key.IsComputedProp() {
					b.outer.Report(JSError_MakeWithNode(key, TypedScopeCreator_INVALID_ENUM_KEY))
					key = key.GetNext()
					continue
				}
				keyName := key.GetString()
				jutil.CheckNotNullMsg1(keyName, "Invalid enum key: %s", key)
				enumType.DefineElement(keyName, key)
				key = key.GetNext()
			}
		}
	}

	if name != nil {
		b.outer.typeRegistry.DeclareType(b.currentScope, *name, enumType.GetElementsType())
	}

	if rValue == nil || !(rValue.IsObjectLit() || rValue.IsQualifiedName()) {
		errorNode := rValue
		if lValue != nil {
			errorNode = lValue
		}
		b.outer.Report(JSError_MakeWithNode(errorNode, TypedScopeCreator_ENUM_INITIALIZER))
	}
	return enumType
}

// TypedScopeCreator_AbstractScopeBuilder_SlotDefiner is responsible for
// defining typed variable "slots".
type TypedScopeCreator_AbstractScopeBuilder_SlotDefiner struct {
	// Go-only: language-adapter — the enclosing AbstractScopeBuilder instance
	// of this inner class.
	outer *TypedScopeCreator_AbstractScopeBuilderBase

	declarationNode *rhino.Node
	variableName    *string
	scope           *TypedScope
	// default is no type and a type may be inferred later
	type_                   rhino.JSType
	allowLaterTypeInference bool
	forGoogProvidedName     bool

	// TODO(bradfordcsmith): Once all the logic needed for ES_2017 features has been added,
	//     make the API to this class more restrictive to avoid accidental misuse.
	//     e.g. There will probably always be a declarationNode, so make it a constructor
	//     parameter.
}

// NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner is `new SlotDefiner()`
// inside the AbstractScopeBuilder outer.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.SlotDefiner()
func NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(outer *TypedScopeCreator_AbstractScopeBuilderBase) *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner {
	d := &TypedScopeCreator_AbstractScopeBuilder_SlotDefiner{outer: outer}
	d.allowLaterTypeInference = true
	d.forGoogProvidedName = false
	return d
}

// ForDeclarationNode sets the defining NAME or GETPROP or object literal key
// node.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.SlotDefiner.forDeclarationNode(Node)
func (d *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner) ForDeclarationNode(declarationNode *rhino.Node) *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner {
	d.declarationNode = declarationNode
	return d
}

// Java: TypedScopeCreator.AbstractScopeBuilder.SlotDefiner.readVariableNameFromDeclarationNode()
func (d *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner) ReadVariableNameFromDeclarationNode() *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner {
	// Only qualified name nodes can use this method to get the variable name
	// Object literal keys will have to compute their names themselves.
	// TODO(bradfordcsmith): Clean up these checks of the parent.
	parent := d.declarationNode.GetParent()
	if d.declarationNode.IsName() {
		jutil.CheckArgument(
			parent.IsFunction() ||
				parent.IsClass() ||
				NodeUtil_IsNameDeclaration(parent) ||
				parent.IsParamList() ||
				(parent.IsRest() && parent.GetParent().IsParamList()) ||
				parent.IsCatch())
	} else {
		jutil.CheckArgument(
			d.declarationNode.IsGetProp() && (parent.IsAssign() || parent.IsExprResult()))
	}
	d.variableName = d.declarationNode.GetQualifiedName()
	return d
}

// ForVariableName sets the name of the variable.
//
// TODO(bradfordcsmith): maybe change to withVariableName(). Need to make
// these names more consistent.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.SlotDefiner.forVariableName(String)
func (d *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner) ForVariableName(variableName string) *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner {
	d.variableName = jutil.StringPtr(variableName)
	return d
}

// InScope sets the scope in which the variable should be declared.
//
// If the given name is a qualified name, this scope should be the scope in
// which the root of the name is (or will later be) declared.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.SlotDefiner.inScope(TypedScope)
func (d *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner) InScope(scope *TypedScope) *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner {
	d.scope = jutil.CheckNotNull(scope)
	return d
}

// Java: TypedScopeCreator.AbstractScopeBuilder.SlotDefiner.withType(JSType)
func (d *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner) WithType(type_ rhino.JSType) *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner {
	d.type_ = type_
	return d
}

// Java: TypedScopeCreator.AbstractScopeBuilder.SlotDefiner.allowLaterTypeInference(boolean)
func (d *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner) AllowLaterTypeInference(allowLaterTypeInference bool) *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner {
	d.allowLaterTypeInference = allowLaterTypeInference
	return d
}

// Java: TypedScopeCreator.AbstractScopeBuilder.SlotDefiner.forGoogProvidedName()
func (d *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner) ForGoogProvidedName() *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner {
	d.forGoogProvidedName = true
	return d
}

// DefineSlot defines the slot and does related work.
//
// At minimum the declaration node and variable name must have been set.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.SlotDefiner.defineSlot()
func (d *TypedScopeCreator_AbstractScopeBuilder_SlotDefiner) DefineSlot() {
	b := d.outer
	jutil.CheckNotNullMsg(d.declarationNode, "declarationNode not set")
	jutil.CheckStateMsg1(
		d.declarationNode.IsName() ||
			d.declarationNode.IsGetProp() ||
			NodeUtil_MayBeObjectLitKey(d.declarationNode) ||
			d.declarationNode.IsModuleBody() ||
			d.declarationNode.IsExport() ||
			(d.forGoogProvidedName && NodeUtil_IsExprCall(d.declarationNode)),
		"declaration node must be an lvalue or goog.provide call, found %s",
		d.declarationNode)
	jutil.CheckNotNullMsg(d.variableName, "variableName not set")
	jutil.CheckStateMsg(d.allowLaterTypeInference || d.type_ != nil, "null type but inference not allowed")
	jutil.CheckState(*d.variableName != "")
	jutil.CheckNotNull(d.scope)
	variableName := *d.variableName

	parent := d.declarationNode.GetParent()

	scopeToDeclareIn := d.scope

	isGlobalVar := d.declarationNode.IsName() && scopeToDeclareIn.IsGlobal()
	shouldDeclareOnGlobalThis :=
		(isGlobalVar && (parent.IsVar() || parent.IsFunction())) ||
			(d.forGoogProvidedName && !strings.Contains(variableName, "."))

	// TODO(sdh): Remove this special case.  It is required to reproduce the original
	// non-block-scoped behavior, which is depended on in several places including
	// https://github.com/angular/tsickle/issues/761.  But it's more correct to always
	// declare on the owner scope.  Once all the bugs are fixed, this should be removed.
	// We may be able to get by with checking a "declared" function's source for jsdoc.
	if scopeToDeclareIn != b.currentHoistScope &&
		scopeToDeclareIn.IsGlobal() &&
		scopeToDeclareIn.HasOwnSlot(variableName) &&
		!d.forGoogProvidedName {
		scopeToDeclareIn = b.currentHoistScope
	}

	// The input may be null if we are working with a AST snippet. So read
	// the extern info from the node.

	// declared in closest scope?
	input := b.outer.compiler.GetInput(b.inputId)
	if !scopeToDeclareIn.CanDeclare(variableName) {
		oldVar := scopeToDeclareIn.GetVar(variableName)
		b.outer.validator.ExpectUndeclaredVariable(
			typedScopeCreator_sourceNameArg(b.sourceName), input, d.declarationNode, parent, oldVar, variableName, d.type_)
	} else {
		if d.type_ != nil {
			b.outer.SetDeferredType(d.declarationNode, d.type_)
		}

		b.Declare(
			scopeToDeclareIn,
			variableName,
			d.declarationNode,
			d.type_,
			input,
			d.allowLaterTypeInference)
	}

	// We need to do some additional work for constructors and interfaces.
	fnType := rhino.JSType_ToMaybeFunctionType(d.type_)
	if fnType != nil &&
		// We don't want to look at empty function types.
		!d.type_.IsEmptyType() {

		// We want to make sure that when we declare a new instance type
		// (with @constructor) that there's actually a ctor for it.
		// This doesn't apply to structural constructors (like
		// function(new:Array). Checking the constructed type against
		// the variable name is a sufficient check for this.
		if fnType.IsConstructor() || fnType.IsInterface() {
			b.FinishConstructorDefinition(
				d.declarationNode, variableName, fnType, scopeToDeclareIn, input)
		}
	}

	if shouldDeclareOnGlobalThis {
		globalThis := b.outer.typeRegistry.GetNativeObjectType(rhino.JSTypeNative_GLOBAL_THIS)
		if d.allowLaterTypeInference {
			inferredType := d.type_
			if inferredType == nil {
				inferredType = b.outer.GetNativeType(rhino.JSTypeNative_NO_TYPE)
			}
			globalThis.DefineInferredProperty(
				variableName,
				inferredType,
				d.declarationNode)
		} else {
			globalThis.DefineDeclaredProperty(variableName, d.type_, d.declarationNode)
		}
	}

	if isGlobalVar &&
		"Window" == variableName &&
		d.type_ != nil &&
		d.type_.IsFunctionType() &&
		d.type_.IsConstructor() {
		globalThisCtor :=
			b.outer.typeRegistry.GetNativeObjectType(rhino.JSTypeNative_GLOBAL_THIS).GetConstructor()
		globalThisCtor.GetInstanceType().ClearCachedValues()
		b.outer.typeRegistry.ResetImplicitPrototype(
			globalThisCtor.GetPrototype(), d.type_.ToMaybeFunctionType().GetInstanceType())
	}
}

// typedScopeCreator_sourceNameArg passes the AbstractScopeBuilder's
// sourceName (a Java String that may be null) to
// TypeValidator.expectUndeclaredVariable, whose Go port takes it as a Go
// string (and, like the Java, never reads it).
//
// Go-only: language-adapter — Java null strings are nil *string.
func typedScopeCreator_sourceNameArg(sourceName *string) string {
	if sourceName == nil {
		return ""
	}
	return *sourceName
}

// Declare declares a variable with the given name and type on the given
// scope, returning the newly-declared TypedVar. Additionally checks the
// escapedVarNames and assignedVarNames maps (which were populated during the
// FirstOrderFunctionAnalyzer and marks the result as escaped or assigned
// exactly once if appropriate.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.declare(TypedScope, String, Node, JSType, CompilerInput, boolean)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) Declare(
	scope *TypedScope, name string, n *rhino.Node, type_ rhino.JSType, input *CompilerInput, inferred bool) *TypedVar {
	var_ := scope.Declare(name, n, type_, input, inferred)
	scopedName := ScopedName_Of(name, scope.GetRootNode())
	if b.outer.escapedVarNames.Contains(scopedName) {
		var_.MarkEscaped()
	}
	if b.outer.assignedVarNames.Count(scopedName) == 1 {
		var_.MarkAssignedExactlyOnce()
	}
	b.outer.assignedVarNames.Remove(scopedName) // free up memory
	return var_
}

// Java: TypedScopeCreator.AbstractScopeBuilder.finishConstructorDefinition(Node, String, FunctionType, TypedScope, CompilerInput)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) FinishConstructorDefinition(
	declarationNode *rhino.Node,
	variableName string,
	fnType rhino.FunctionType,
	scopeToDeclareIn *TypedScope,
	input *CompilerInput) {
	// Declare var.prototype in the scope chain.
	superClassCtor := fnType.GetSuperClassConstructor()
	prototypeSlot := fnType.GetSlot("prototype")

	prototypeName := variableName + ".prototype"

	// There are some rare cases where the prototype will already
	// be declared. See TypedScopeCreatorTest#testBogusPrototypeInit.
	// Fortunately, other warnings will complain if this happens.
	prototypeVar := scopeToDeclareIn.GetVar(prototypeName)
	if prototypeVar != nil && prototypeVar.GetScope() == rhino.StaticScope(scopeToDeclareIn) {
		scopeToDeclareIn.Undeclare(prototypeVar)
	}

	scopeToDeclareIn.Declare(
		prototypeName,
		declarationNode,
		prototypeSlot.GetType(),
		input,
		// declared iff there's an explicit supertype
		superClassCtor == nil ||
			superClassCtor.GetInstanceType().Equals(b.outer.GetNativeType(rhino.JSTypeNative_OBJECT_TYPE)))
}

// IsLValueRootedInGlobalScope checks if the given node is a property of a
// name in the global scope.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.isLValueRootedInGlobalScope(Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) IsLValueRootedInGlobalScope(n *rhino.Node) bool {
	return b.GetLValueRootScope(n).IsGlobal()
}

// GetLValueRootScope returns the scope for the name of the given node.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.getLValueRootScope(Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) GetLValueRootScope(n *rhino.Node) *TypedScope {
	root := NodeUtil_GetBestLValueRoot(n)
	if root != nil {
		if root.IsName() {
			nameParent := root.GetParent()
			switch nameParent.GetToken() {
			case rhino.Token_VAR:
				return b.currentHoistScope
			case rhino.Token_LET, rhino.Token_CONST, rhino.Token_CLASS, rhino.Token_FUNCTION, rhino.Token_PARAM_LIST, rhino.Token_CATCH:
				return b.currentScope
			case rhino.Token_ITER_REST, rhino.Token_OBJECT_REST:
				// TODO(bradfordcsmith): Handle array destructuring REST
				jutil.CheckStateMsg1(nameParent.GetParent().IsParamList(), "%s", nameParent)
				return b.currentScope
			default:
				if b.IsGoogModuleExports(root) {
					// Ensure that 'exports = class {}' in a goog.module returns the module scope.
					return b.currentScope
				}
				var_ := b.currentScope.GetVar(root.GetString())
				if var_ != nil {
					return var_.GetScope().(*TypedScope)
				}
			}
		} else if root.IsThis() || root.IsSuper() {
			// We want the enclosing function scope, or the global scope if not in a function.
			return b.currentHoistScope.GetScopeOfThis()
		}
	}
	return b.currentHoistScope.GetGlobalScope()
}

// GetDeclaredType looks for a type declaration on a property assignment (in
// an ASSIGN or an object literal key).
//
// info is the doc info for this property; lValue the l-value node; rValue
// the node that n is being initialized to, or null if this is a stub
// declaration; declaredRValueTypeSupplier a supplier for the declared type
// of the rvalue, used for destructuring declarations where we have to do
// additional work on the rvalue.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.getDeclaredType(JSDocInfo, Node, Node, Supplier)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) GetDeclaredType(
	info *rhino.JSDocInfo,
	lValue *rhino.Node,
	rValue *rhino.Node,
	declaredRValueTypeSupplier func() *TypedScopeCreator_RValueInfo) rhino.JSType {
	if info != nil && info.HasType() {
		type_ := b.GetDeclaredTypeInAnnotation(lValue, info)

		if type_.IsSymbol() &&
			lValue.IsFromExterns() &&
			lValue.IsGetProp() &&
			lValue.GetFirstChild().MatchesName("Symbol") &&
			b.currentScope.IsGlobal() {
			// Create a unique subtype type for this name.
			// Currently we only look specifically for extern properties on the global Symbol object
			// itself - so
			//   /** @const {symbol} */
			//   Symbol.iterator;
			// but not:
			//   /** @const {symbol} */
			//   const x = Symbol();
			// We could consider relaxing that in the future & also defining known symbols from
			// source code.
			return rhino.NewKnownSymbolType(b.outer.typeRegistry, *lValue.GetQualifiedName())
		}
		return type_
	} else if rValue != nil &&
		rValue.IsFunction() &&
		b.ShouldUseFunctionLiteralType(
			rhino.JSType_ToMaybeFunctionType(rValue.GetJSType()), info, lValue) {
		return rValue.GetJSType()
	} else if rValue != nil && rValue.IsClass() {
		return rValue.GetJSType()
	} else if info != nil {
		if info.HasEnumParameterType() {
			if rValue != nil && rValue.IsObjectLit() {
				return rValue.GetJSType()
			} else {
				return b.CreateEnumTypeFromNodes(rValue, lValue.GetQualifiedName(), lValue, info)
			}
		} else if info.IsConstructorOrInterface() {
			fnType :=
				b.CreateFunctionTypeFromNodes(rValue, lValue.GetQualifiedName(), info, lValue)
			if rValue == nil && !lValue.IsFromExterns() {
				diagnostic := TypedScopeCreator_IFACE_INITIALIZER
				if fnType.IsConstructor() {
					diagnostic = TypedScopeCreator_CTOR_INITIALIZER
				}
				b.outer.Report(
					JSError_MakeWithNode(
						lValue,
						diagnostic,
						jutil.ToJavaString(lValue.GetQualifiedName())))
			}
			return fnType
		}
	}

	// Check if this is constant and if it has a known type.
	if (lValue.IsExport() && lValue.GetBooleanProp(rhino.Node_EXPORT_DEFAULT)) ||
		NodeUtil_IsConstantDeclaration(info, lValue) ||
		b.IsGoogModuleExports(lValue) {
		if rValue != nil {
			rValueType := b.GetDeclaredRValueType(lValue, rValue)
			b.DeclareAliasTypeIfRvalueIsAliasable(
				lValue, rValue.GetQualifiedNameObject(), rValueType, b.currentScope)
			if rValueType != nil {
				return rValueType
			}
		} else if declaredRValueTypeSupplier != nil {
			rvalueInfo := declaredRValueTypeSupplier()
			if rvalueInfo != nil {
				b.DeclareAliasTypeIfRvalueIsAliasable(
					lValue, rvalueInfo.qualifiedName, rvalueInfo.type_, b.currentScope)
				if rvalueInfo.type_ != nil {
					return rvalueInfo.type_
				}
			}
		}
	}

	if rValue != nil && rValue.IsAssign() {
		// Handle nested assignments. For example, TypeScript generates code like this:
		//   var Foo_1;
		//   let Foo = Foo_1 = class Foo {}
		//   Foo = Foo_1 = tslib_1.decorate(..., Foo);
		return b.GetDeclaredType(info, lValue, rValue.GetSecondChild(), nil)
	}

	if info != nil && FunctionTypeBuilder_IsFunctionTypeDeclaration(info) {
		fnName := lValue.GetQualifiedName()
		return b.CreateFunctionTypeFromNodes(nil, fnName, info, lValue)
	}

	if b.IsValidTypedefDeclaration(lValue, info) {
		return b.outer.GetNativeType(rhino.JSTypeNative_NO_TYPE)
	}

	return nil
}

// DeclareAliasTypeIfRvalueIsAliasable: for a const alias, like
// `const alias = other.name`, this may declare `alias` as a type name,
// depending on what other.name is defined to be.
//
// This method recognizes three kinds of type aliases: @typedefs,
// @constructor/@interface types, and @enums.
//
// Given any of those three types, this method redeclares the aliasing name
// in the typeRegistry. For @typedefs and global @enums, this method also
// marks the qualified name referring to the type as non-nullable by
// default.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.declareAliasTypeIfRvalueIsAliasable(Node, QualifiedName, JSType, TypedScope)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) DeclareAliasTypeIfRvalueIsAliasable(
	lValue *rhino.Node,
	rValue rhino.QualifiedName,
	rValueType rhino.JSType,
	rValueLookupScope *TypedScope) {
	b.DeclareAliasTypeIfRvalueIsAliasableWithLValueName(
		lValue.GetQualifiedName(), lValue, rValue, rValueType, rValueLookupScope, b.currentScope)
}

// DeclareAliasTypeIfRvalueIsAliasableWithLValueName: for a const alias, like
// `const alias = other.name`, this may declare `alias` as a type name,
// depending on what other.name is defined to be.
//
// NOTE: in most cases, call the version with fewer arguments. This version
// only exists to handle goog.declareLegacyNamespace, which is strange
// compared to normal type aliasing because 1) there's no GETPROP node
// representing the lvalue and 2) the type is declared in the global scope,
// not the current module-local scope.
//
// lValueName is the fully qualified lValue name, if any. If null, all this
// method will do is propagate the @typedef Node annotation to
// actualLvalueNode. aliasDeclarationScope is the scope in which to declare
// the alias name. In most cases, this should just be the currentScope.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.declareAliasTypeIfRvalueIsAliasable(String, Node, QualifiedName, JSType, TypedScope, TypedScope)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) DeclareAliasTypeIfRvalueIsAliasableWithLValueName(
	lValueName *string,
	actualLvalueNode *rhino.Node,
	rValue rhino.QualifiedName,
	rValueType rhino.JSType,
	rValueLookupScope *TypedScope,
	aliasDeclarationScope *TypedScope) {
	// NOTE: this allows some strange patterns such allowing instance properties
	// to be aliases of constructors, and then creating a local alias of that to be
	// used as a type name.  Consider restricting this.

	if rValue == nil {
		return
	}

	// Look for a @typedef annotation on the definition node
	definitionNode := b.GetDefinitionNode(rValue, rValueLookupScope)
	if definitionNode != nil {
		typedefType := definitionNode.GetTypedefTypeProp()
		if typedefType != nil {
			// Propagate typedef type to typedef aliases.
			actualLvalueNode.SetTypedefTypeProp(typedefType)
			if lValueName != nil {
				b.outer.typeRegistry.IdentifyNonNullableName(aliasDeclarationScope, *lValueName)
				b.outer.typeRegistry.DeclareType(aliasDeclarationScope, *lValueName, typedefType)
			}
			return
		}
	}

	if lValueName == nil {
		return
	}

	// Check if the provided rValueType indicates that we should declare this type
	// Note that we only look for enums and constructors/interfaces here: this step cannot work
	// for @typedefs. The 'type' of the TypedVar representing a @typedef'd name is the None type,
	// not the @typedef'd type.
	if rValueType != nil &&
		rValueType.IsFunctionType() &&
		rValueType.ToMaybeFunctionType().HasInstanceType() {
		// Look for @constructor/@interface by checking if the RHS has an instance type
		functionType := rValueType.ToMaybeFunctionType()
		b.outer.typeRegistry.DeclareType(aliasDeclarationScope, *lValueName, functionType.GetInstanceType())
		return
	}

	if rValueType != nil && rValueType.IsEnumType() {
		// Look for cases where the rValue is an Enum namespace
		b.outer.typeRegistry.DeclareType(
			aliasDeclarationScope, *lValueName, rValueType.ToMaybeEnumType().GetElementsType())
		b.outer.typeRegistry.IdentifyNonNullableName(aliasDeclarationScope, *lValueName)
	}
}

// IsGoogModuleExports reports whether this lvalue is either `exports`,
// `exports.x`, or a string key in `exports = {x}`.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.isGoogModuleExports(Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) IsGoogModuleExports(lValue *rhino.Node) bool {
	if b.GetModule() == nil || lValue == nil {
		return false
	}
	if b.outer.undeclaredNamesForClosure.Contains(lValue) {
		// includes "exports" and "exports.x"
		return true
	}
	// ASSIGN
	//   NAME "exports"
	//   OBJECT_LIT
	//     STRING_KEY "x"
	//     [value]
	return lValue.IsStringKey() &&
		lValue.GetParent().IsObjectLit() &&
		lValue.GetGrandparent().IsAssign() &&
		lValue.GetParent().GetPrevious().MatchesName("exports") &&
		b.outer.undeclaredNamesForClosure.Contains(lValue.GetParent().GetPrevious())
}

// GetDefinitionNode returns the AST node associated with the definition, if
// any.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.getDefinitionNode(QualifiedName, TypedScope)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) GetDefinitionNode(qname rhino.QualifiedName, scope *TypedScope) *rhino.Node {
	if qname.IsSimple() {
		var_ := scope.GetVar(qname.GetComponent())
		if var_ != nil {
			return var_.GetNameNode()
		}
		return nil
	}
	parent := rhino.ObjectType_Cast(scope.LookupQualifiedName(qname.GetOwner()))
	if parent != nil {
		return parent.GetPropertyDefSite(qname.GetComponent())
	}
	return nil
}

// GetDeclaredRValueType checks for common idioms of a typed R-value
// assigned to a const L-value.
//
// Normally, we would only want this sort of propagation to happen under
// type inference. But we want a declared const to be nameable in a type
// annotation, so we need to figure out the type before we try to resolve
// the annotation.
//
// lValue is the lvalue node if this is a simple assignment, null for
// destructuring.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.getDeclaredRValueType(Node, Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) GetDeclaredRValueType(lValue *rhino.Node, rValue *rhino.Node) rhino.JSType {
	// If rValue has a type-cast, we use the type in the type-cast.
	rValueInfo := rValue.GetJSDocInfo()
	if rValue.IsCast() && rValueInfo != nil && rValueInfo.HasType() {
		return rValueInfo.GetType().Evaluate(b.currentScope, b.outer.typeRegistry)
	}

	// Check if the type has already been computed during scope-creation.
	// This is mostly useful for literals like BOOLEAN, NUMBER, STRING, and
	// OBJECT_LITERAL
	type_ := rValue.GetJSType()
	if type_ != nil && !type_.IsUnknownType() {
		return type_
	}

	// If rValue is a name, try looking it up in the current scope.
	if rValue.IsQualifiedName() {
		return b.currentScope.LookupQualifiedName(rValue.GetQualifiedNameObject())
	}

	// Check for simple invariant operations, such as "!x" or "+x" or "''+x"
	if NodeUtil_IsBooleanResult(rValue) {
		return b.outer.GetNativeType(rhino.JSTypeNative_BOOLEAN_TYPE)
	}

	if NodeUtil_IsNumericResult(rValue) {
		return b.outer.GetNativeType(rhino.JSTypeNative_NUMBER_TYPE)
	}

	if NodeUtil_IsBigIntResult(rValue) {
		return b.outer.GetNativeType(rhino.JSTypeNative_BIGINT_TYPE)
	}

	if NodeUtil_IsStringResult(rValue) {
		return b.outer.GetNativeType(rhino.JSTypeNative_STRING_TYPE)
	}

	if rValue.IsNew() && rValue.GetFirstChild().IsQualifiedName() {
		targetType :=
			b.currentScope.LookupQualifiedName(rValue.GetFirstChild().GetQualifiedNameObject())
		if targetType != nil {
			fnType := targetType.RestrictByNotNullOrUndefined().ToMaybeFunctionType()
			if fnType != nil && fnType.HasInstanceType() {
				return fnType.GetInstanceType()
			}
		}
	}

	if rValue.IsCall() && rValue.GetFirstChild().MatchesName("Symbol") {
		// Type calls like `Symbol('foo')`
		symbolVar := b.currentScope.GetVar("Symbol")
		if symbolVar != nil && symbolVar.GetScope().(*TypedScope).IsGlobal() {
			return b.outer.typeRegistry.GetNativeType(rhino.JSTypeNative_SYMBOL_TYPE)
		}
	}

	// Check for a very specific JS idiom:
	// var x = x || TYPE;
	// This is used by Closure's base namespace for esoteric
	// reasons, so we only really care about that case.
	if rValue.IsOr() {
		firstClause := rValue.GetFirstChild()
		secondClause := firstClause.GetNext()
		namesMatch :=
			firstClause.IsName() &&
				lValue != nil &&
				lValue.IsName() &&
				firstClause.GetString() == lValue.GetString()
		if namesMatch {
			type_ = secondClause.GetJSType()
			if type_ != nil && !type_.IsUnknownType() {
				return type_
			}
		}
	}

	return nil
}

// CheckForClassDefiningCalls looks for class-defining calls. Because JS has
// no 'native' syntax for defining classes, this is often very
// coding-convention dependent and business-logic heavy.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.checkForClassDefiningCalls(Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) CheckForClassDefiningCalls(n *rhino.Node) {
	relationship := b.outer.codingConvention.GetClassesDefinedByCall(n)
	if relationship != nil {
		superClass :=
			TypeValidator_GetInstanceOfCtor(
				b.currentScope.LookupQualifiedName(rhino.QualifiedName_Of(relationship.superclassName)))
		subClass :=
			TypeValidator_GetInstanceOfCtor(
				b.currentScope.LookupQualifiedName(rhino.QualifiedName_Of(relationship.subclassName)))
		if superClass != nil && subClass != nil {
			// superCtor and subCtor might be structural constructors
			// (like {function(new:Object)}) so we need to resolve them back
			// to the original ctor objects.
			superCtor := superClass.GetConstructor()
			subCtor := subClass.GetConstructor()
			if superCtor != nil && subCtor != nil {
				b.outer.codingConvention.ApplySubclassRelationship(
					rhino.NewNominalTypeBuilder(superCtor, superClass),
					rhino.NewNominalTypeBuilder(subCtor, subClass),
					relationship.type_)
			}
		}
	}

	singletonGetterClassName := b.outer.codingConvention.GetSingletonGetterClassName(n)
	if singletonGetterClassName != nil {
		objectType :=
			rhino.ObjectType_Cast(b.outer.typeRegistry.GetType(b.currentScope, *singletonGetterClassName))
		if objectType != nil {
			functionType := objectType.GetConstructor()

			if functionType != nil {
				getterType := b.outer.typeRegistry.CreateFunctionType(objectType)
				b.outer.codingConvention.ApplySingletonGetter(
					rhino.NewNominalTypeBuilder(functionType, objectType), getterType)
			}
		}
	}

	objectLiteralCast := b.outer.codingConvention.GetObjectLiteralCast(n)
	if objectLiteralCast != nil {
		if objectLiteralCast.diagnosticType == nil {
			type_ :=
				rhino.ObjectType_Cast(b.outer.typeRegistry.GetType(b.currentScope, *objectLiteralCast.typeName))
			if type_ != nil && type_.GetConstructor() != nil {
				b.outer.SetDeferredType(objectLiteralCast.objectNode, type_)
				objectLiteralCast.objectNode.PutBooleanProp(rhino.Node_REFLECTED_OBJECT, true)
			} else {
				b.outer.Report(JSError_MakeWithNode(n, TypedScopeCreator_CONSTRUCTOR_EXPECTED))
			}
		} else {
			b.outer.Report(JSError_MakeWithNode(n, objectLiteralCast.diagnosticType))
		}
	}
}

// MaybeDeclareGetElem declares a computed property on its owner type.
//
// info is the doc info for this property; n a top-level GETELEM node (it
// should not be contained inside another GETPROP); rhsValue the node that n
// is being initialized to, or null if this is a stub declaration.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.maybeDeclareGetElem(JSDocInfo, Node, Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) MaybeDeclareGetElem(info *rhino.JSDocInfo, n *rhino.Node, rhsValue *rhino.Node) {
	jutil.CheckArgumentMsg1(n.IsGetElem(), "%s", n)
	ownerNode := n.GetFirstChild()
	key := n.GetSecondChild()
	keyType := TypedScopeCreator_ExtractKnownSymbolKey(b.currentScope, key)
	if keyType == nil {
		return
	}
	valueType := b.GetDeclaredType(info, n, rhsValue, nil)
	if valueType == nil {
		return
	}
	ownerName := ownerNode.GetQualifiedName()
	ownerType := b.GetObjectSlot(*ownerName)
	if ownerType == nil {
		return
	}
	b.DeclarePropertyIfNamespaceType(
		ownerType,
		ownerNode,
		rhino.NewProperty_SymbolKey(keyType),
		valueType,
		n,
		/* jsdocInfo= */ nil)
}

// MaybeDeclareQualifiedName declares the symbol for a qualified name in the
// current scope.
//
// info is the doc info for this property; n a top-level GETPROP node (it
// should not be contained inside another GETPROP); parent the parent of n;
// rhsValue the node that n is being initialized to, or null if this is a
// stub declaration.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.maybeDeclareQualifiedName(NodeTraversal, JSDocInfo, Node, Node, Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) MaybeDeclareQualifiedName(
	t *NodeTraversal, info *rhino.JSDocInfo, n *rhino.Node, parent *rhino.Node, rhsValue *rhino.Node) {
	isTypedef := b.IsValidTypedefDeclaration(n, info)
	if isTypedef {
		b.DeclareTypedefType(n, info)
	}

	ownerNode := n.GetFirstChild()
	ownerName := ownerNode.GetQualifiedName()
	qName := n.GetQualifiedName()
	propName := n.GetString()
	jutil.CheckArgument(qName != nil && ownerName != nil)

	// Precedence of type information on GETPROPs:
	// 1) @type annotation / @enum annotation
	// 2) ASSIGN to FUNCTION literal
	// 3) @param/@return annotation (with no function literal)
	// 4) ASSIGN to something marked @const
	// 5) ASSIGN to anything else
	//
	// 1, 3, and 4 are declarations, 5 is inferred, and 2 is a declaration iff
	// the function has JsDoc or has not been declared before.
	//
	// FUNCTION literals are special because TypedScopeCreator is very smart
	// about getting as much type information as possible for them.

	// Determining type for #1 + #2 + #3 + #4
	valueType := b.GetDeclaredType(info, n, rhsValue, nil)
	if valueType == nil && rhsValue != nil {
		// Determining type for #5
		valueType = rhsValue.GetJSType()
	}

	// Function prototypes are special.
	// It's a common JS idiom to do:
	// F.prototype = { ... };
	// So if F does not have an explicitly declared super type,
	// allow F.prototype to be redefined arbitrarily.
	if "prototype" == propName {
		qVar := b.currentScope.GetVar(*qName)
		if qVar != nil {
			// If the programmer has declared that F inherits from Super,
			// and they assign F.prototype to an object literal,
			// then they are responsible for making sure that the object literal's
			// implicit prototype is set up appropriately. We just obey
			// the @extends tag.
			qVarType := rhino.ObjectType_Cast(qVar.GetType())
			if qVarType != nil && rhsValue != nil && rhsValue.IsObjectLit() {
				b.outer.typeRegistry.ResetImplicitPrototype(
					rhsValue.GetJSType(), qVarType.GetImplicitPrototype())
			} else if !qVar.IsTypeInferred() {
				// If the programmer has declared that F inherits from Super,
				// and they assign F.prototype to some arbitrary expression,
				// there's not much we can do. We just ignore the expression,
				// and hope they've annotated their code in a way to tell us
				// what props are going to be on that prototype.
				return
			}

			qVar.GetScope().(*TypedScope).Undeclare(qVar)
		}
	}

	if valueType == nil {
		if parent.IsExprResult() {
			// t is mutable so make sure to capture the current state before the lambda.
			isExtern := t.GetInput() != nil && t.GetInput().IsExtern()
			action := func() { b.ResolveStubDeclaration(n, isExtern, *ownerName) }
			b.deferredActions.Put(b.currentScope.GetRootNode(), &action)
		}

		return
	}

	inferred := b.IsQualifiedNameInferred(qName, n, info, rhsValue, valueType)
	if !inferred {
		ownerType := b.GetObjectSlot(*ownerName)
		if ownerType != nil {
			b.DeclarePropertyIfNamespaceType(
				ownerType, ownerNode, rhino.NewProperty_StringKey(propName), valueType, n /* jsdocInfo= */, nil)
		}

		// this is a memory optimization: we don't need to declare .prototype props in the scope.
		// NOTE: in theory we could possibly extend this to more kinds of declarations and only
		// declare simple (non-qualified) names in the scope, but that seems to cause a lot more
		// issues with older code.
		declareInScope := ownerType == nil || !ownerType.IsFunctionPrototypeType()

		// If the property is already declared, the error will be caught when we try to declare it
		// in the current scope.
		if declareInScope {
			NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(b).
				ForDeclarationNode(n).
				ForVariableName(*qName).
				InScope(b.GetLValueRootScope(n)).
				WithType(valueType).
				AllowLaterTypeInference(inferred).
				DefineSlot()
		}
	}
}

// IsQualifiedNameInferred determines whether a qualified name is inferred.
// NOTE(nicksantos): Determining whether a property is declared or not is
// really really obnoxious.
//
// The problem is that there are two (equally valid) coding styles:
//
// (function() { /* The authoritative definition of goog.bar. / goog.bar =
// function() {}; })();
//
// function f() { goog.bar(); /* Reset goog.bar to a no-op. / goog.bar =
// function() {}; }
//
// In a dynamic language with first-class functions, it's very difficult to
// know which one the user intended without looking at lots of contextual
// information (the second example demonstrates a small case of this, but
// there are some really pathological cases as well).
//
// The current algorithm checks if either the declaration has JsDoc type
// information, or @const with a known type, or a function literal with a
// name we haven't seen before.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.isQualifiedNameInferred(String, Node, JSDocInfo, Node, JSType)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) IsQualifiedNameInferred(
	qName *string,
	n *rhino.Node,
	info *rhino.JSDocInfo,
	rhsValue *rhino.Node,
	valueType rhino.JSType) bool {
	// Prototypes of constructors and interfaces are always declared.
	if qName != nil && jutil.StringEndsWith(*qName, ".prototype") {
		className := jutil.StringSubstring(*qName, 0, jutil.StringLastIndexOf(*qName, ".prototype"))
		slot := b.currentScope.GetVar(className)
		var classType rhino.JSType
		if slot != nil {
			classType = slot.GetType()
		}
		if classType != nil && (classType.IsConstructor() || classType.IsInterface()) {
			return false
		}
	}
	// treat "foo = bar = <VAL>" the same as "foo = <VAL>"
	rhsValue = b.UnwrapIfAssign(rhsValue)

	// If the jsdoc or RHS specifies a concrete type, it's not inferred.
	if (info != nil &&
		(info.HasType() ||
			info.HasEnumParameterType() ||
			b.IsValidTypedefDeclaration(n, info) ||
			FunctionTypeBuilder_IsFunctionTypeDeclaration(info) ||
			(rhsValue != nil && rhsValue.IsFunction()))) ||
		b.IsTypedConstantDeclaration(info, n, valueType) {
		return false
	}

	// At this point, we're pretty sure it's inferred, since there's neither
	// useful jsdoc info, nor a useful const or doc'd function RHS.  But
	// there's still one case where it may still not be: if the RHS is a
	// class or function that is not
	//   (1) a scoped qualified name (i.e. this.b.c or super.b.c),
	//   (2) already declared in a scope,
	//   (3) assigned in a conditional block, or
	//   (4) escaped to a closure,
	// then we treat it as if it is declared, rather than inferred.
	// Stubs and other values are always considered inferred at this point.
	if rhsValue == nil || (!rhsValue.IsFunction() && !rhsValue.IsClass()) {
		return true
	}

	// "Scoped" qualified names (e.g. this.b.c or super.d) are inferred.
	if !n.IsUnscopedQualifiedName() {
		return true
	}

	// If this qname is already declared then treat this definition as inferred.
	ownerScope := b.GetLValueRootScope(n)
	if ownerScope != nil && ownerScope.HasOwnSlot(*qName) {
		return true
	}

	// Check if this is in a conditional block.
	// Functions assigned in conditional blocks are inferred.
	if b.HasControlStructureAncestor(n.GetParent()) {
		return true
	}

	// Check if this is assigned in an inner scope.
	// Functions assigned in inner scopes are inferred.
	if ownerScope != nil &&
		b.outer.escapedVarNames.Contains(ScopedName_Of(*qName, ownerScope.GetRootNode())) {
		return true
	}

	return false
}

// Java: TypedScopeCreator.AbstractScopeBuilder.unwrapIfAssign(Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) UnwrapIfAssign(maybeAssign *rhino.Node) *rhino.Node {
	if maybeAssign == nil || !maybeAssign.IsAssign() {
		return maybeAssign
	}
	return b.UnwrapIfAssign(maybeAssign.GetSecondChild())
}

// DeclareProvidedNs: given a `goog.provide()` or legacy `goog.module()` call
// and implicit ProvidedName, declares the name in the global scope.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.declareProvidedNs(Node, ProvidedName)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) DeclareProvidedNs(provideCall *rhino.Node, providedName *ProcessClosureProvidesAndRequires_ProvidedName) {
	// Redefine this name if we haven't already added a provide definition.
	// Note: in some cases, this will cause a redefinition error.
	anonymousObjectType := b.outer.typeRegistry.CreateAnonymousObjectType(nil)
	NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(b).
		InScope(b.currentScope.GetGlobalScope()).
		AllowLaterTypeInference(false).
		ForVariableName(providedName.GetNamespace()).
		ForDeclarationNode(provideCall).
		WithType(anonymousObjectType).
		ForGoogProvidedName().
		DefineSlot()

	namespace := rhino.QualifiedName_Of(providedName.GetNamespace())
	if !namespace.IsSimple() {
		ownerType := b.currentScope.LookupQualifiedName(namespace.GetOwner())
		if ownerType != nil && ownerType.IsObjectType() {
			ownerType.
				ToMaybeObjectType().
				DefineDeclaredProperty(namespace.GetComponent(), anonymousObjectType, provideCall)
		}
	}
}

// Java: TypedScopeCreator.AbstractScopeBuilder.isTypedConstantDeclaration(JSDocInfo, Node, JSType)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) IsTypedConstantDeclaration(info *rhino.JSDocInfo, n *rhino.Node, valueType rhino.JSType) bool {
	return (NodeUtil_IsConstantDeclaration(info, n) || b.IsGoogModuleExports(n)) &&
		valueType != nil
}

// Java: TypedScopeCreator.AbstractScopeBuilder.hasControlStructureAncestor(Node)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) HasControlStructureAncestor(n *rhino.Node) bool {
	for !(n.IsScript() || n.IsFunction()) {
		if NodeUtil_IsControlStructure(n) {
			return true
		}
		n = n.GetParent()
	}
	return false
}

// GetObjectSlot finds the ObjectType associated with the given slot.
//
// slotName is the name of the slot to find the type in. Returns an object
// type, or null if this slot does not contain an object.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.getObjectSlot(String)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) GetObjectSlot(slotName string) rhino.ObjectType {
	ownerVar := b.currentScope.GetVar(slotName)
	if ownerVar != nil {
		ownerVarType := ownerVar.GetType()
		if ownerVarType == nil {
			return rhino.ObjectType_Cast(nil)
		}
		return rhino.ObjectType_Cast(ownerVarType.RestrictByNotNullOrUndefined())
	}
	return nil
}

// GetInheritedInterfacePropertyType: when a class has a stub for a
// property, and the property exists on a super interface, use that type.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.getInheritedInterfacePropertyType(ObjectType, String)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) GetInheritedInterfacePropertyType(obj rhino.ObjectType, propName string) rhino.JSType {
	if obj != nil && obj.IsFunctionPrototypeType() {
		f := obj.GetOwnerFunction()
		for _, i := range f.GetImplementedInterfaces() {
			if i.HasProperty(propName) {
				return i.GetPropertyType(propName)
			}
		}
	}
	return nil
}

// ResolveStubDeclaration resolves any type-less stub declarations to
// unknown types if we could not find types for them during traversal. This
// method is only called as a deferred action after the root node is visted.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.resolveStubDeclaration(Node, boolean, String)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) ResolveStubDeclaration(n *rhino.Node, isExtern bool, ownerName string) {
	qName := n.GetQualifiedName()
	propName := n.GetString()

	// TODO(b/111216910): should this be getLValueRoot(n).hasOwnSlot(qName)?
	if b.currentScope.HasOwnSlot(*qName) {
		return
	}

	// If we see a stub property, make sure to register this property
	// in the type registry.
	ownerType := b.GetObjectSlot(ownerName)
	inheritedType := b.GetInheritedInterfacePropertyType(ownerType, propName)
	var stubType rhino.JSType
	if inheritedType == nil {
		stubType = b.outer.unknownType
	} else {
		stubType = inheritedType
	}
	NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(b).
		ForDeclarationNode(n).
		ReadVariableNameFromDeclarationNode().
		InScope(b.GetLValueRootScope(n)).
		WithType(stubType).
		AllowLaterTypeInference(true).
		DefineSlot()

	if ownerType != nil && (isExtern || ownerType.IsFunctionPrototypeType()) {
		// If this is a stub for a prototype, just declare it
		// as an unknown type. These are seen often in externs.
		ownerType.DefineInferredProperty(propName, stubType, n)
	} else {
		var registerOn rhino.JSType
		if ownerType == nil {
			registerOn = stubType
		} else {
			registerOn = ownerType
		}
		b.outer.typeRegistry.RegisterPropertyOnType(propName, registerOn)
	}
}

// IsValidTypedefDeclaration returns whether this is a valid declaration of a
// @typedef.
//
// candidate is a qualified name node; info JSDoc comments.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.isValidTypedefDeclaration(Node, JSDocInfo)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) IsValidTypedefDeclaration(candidate *rhino.Node, info *rhino.JSDocInfo) bool {
	if info == nil || !info.HasTypedefType() {
		return false
	}
	// `isUnscopedQualifiedName` excludes `this` and `super` properties.
	return candidate.IsUnscopedQualifiedName() && !NodeUtil_IsPrototypeProperty(candidate)
}

// DeclareTypedefType declares a typedef'd name in the JSTypeRegistry.
//
// Java: TypedScopeCreator.AbstractScopeBuilder.declareTypedefType(Node, JSDocInfo)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) DeclareTypedefType(candidate *rhino.Node, info *rhino.JSDocInfo) {
	typedef := *candidate.GetQualifiedName()

	// TODO(nicksantos|user): This is a terrible, terrible hack
	// to bail out on recursive typedefs. We'll eventually need
	// to handle these properly.
	b.outer.typeRegistry.DeclareType(b.currentScope, typedef, b.outer.unknownType)

	realType := info.GetTypedefType().Evaluate(b.currentScope, b.outer.typeRegistry)
	if realType == nil {
		b.outer.Report(JSError_MakeWithNode(candidate, TypedScopeCreator_MALFORMED_TYPEDEF, typedef))
	} else {
		candidate.SetTypedefTypeProp(realType)
	}

	b.outer.typeRegistry.OverwriteDeclaredType(b.currentScope, typedef, realType)
}

// Java: TypedScopeCreator.AbstractScopeBuilder.declarePropertyIfNamespaceType(ObjectType, Node, Property.Key, JSType, Node, JSDocInfo)
func (b *TypedScopeCreator_AbstractScopeBuilderBase) DeclarePropertyIfNamespaceType(
	ownerType rhino.ObjectType,
	ownerNode *rhino.Node,
	propName rhino.Property_Key,
	valueType rhino.JSType,
	declarationNode *rhino.Node,
	jsdocInfo *rhino.JSDocInfo) {
	// Only declare this as an official property if it has not been
	// declared yet.
	if ownerType.HasOwnPropertyByKey(propName) && !ownerType.IsPropertyTypeInferredByKey(propName) {
		return
	}
	// Define the property if any of the following are true:
	//   (1) it's a non-native extern type. Native types are excluded here because we don't
	//       want externs of the form "/** @type {!Object} */ var api = {}; api.foo;" to
	//       cause a property "foo" to be declared on Object.
	//   (2) it's a non-instance type. This primarily covers static properties on
	//       constructors (which are FunctionTypes, not InstanceTypes).
	//   (3) it's an assignment to 'this', which covers instance properties assigned in
	//       constructors or other methods.
	isNonNativeExtern :=
		b.GetCompilerInput() != nil &&
			b.GetCompilerInput().IsExtern() &&
			!ownerType.IsNativeObjectType()
	if isNonNativeExtern || !ownerType.IsInstanceType() || ownerNode.IsThis() {
		// If the property is undeclared or inferred, declare it now.
		ownerType.DefineDeclaredPropertyByKey(propName, valueType, declarationNode)
		ownerType.SetPropertyJSDocInfoByKey(propName, jsdocInfo)
	}
} // end AbstractScopeBuilder

// TypedScopeCreator_NormalScopeBuilder is a shallow traversal of the global
// scope to build up all classes, functions, and methods.
type TypedScopeCreator_NormalScopeBuilder struct {
	TypedScopeCreator_AbstractScopeBuilderBase
}

// Java: TypedScopeCreator.NormalScopeBuilder(TypedScope)
func NewTypedScopeCreator_NormalScopeBuilder(outer *TypedScopeCreator, scope *TypedScope) *TypedScopeCreator_NormalScopeBuilder {
	b := &TypedScopeCreator_NormalScopeBuilder{}
	b.InitTypedScopeCreator_AbstractScopeBuilderBase(b, outer, scope)
	return b
}

// Java: TypedScopeCreator.NormalScopeBuilder.visitPreorder(NodeTraversal, Node, Node)
func (b *TypedScopeCreator_NormalScopeBuilder) VisitPreorder(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	// Create any child block scopes "pre-order" as we see them.
	//
	// This is required because hoisted or qualified names defined in earlier blocks might be
	// referred to later outside the block. This isn't a big deal in most cases since a NamedType
	// will be created and resolved later, but if a NamedType is used for a superclass, we lose a
	// lot of valuable checking. Recursing into child blocks immediately prevents this from being
	// a problem.
	//
	// We don't traverse into CLASSes because we haven't yet have created the class-type on which
	// to assign members. We'll do this on the way back up (post-order) instead, after the
	// class-type has been attached to the AST.
	if parent != nil && NodeUtil_CreatesBlockScope(n) && !n.IsClass() {
		b.outer.CreateScope(n, b.currentScope)
	}

	// All other functions (and classes, etc) are handled when we see the actual function node.
	if n.IsFunction() {
		b.DefineFunctionLiteral(n)
	}
}

// Java: TypedScopeCreator.NormalScopeBuilder.visitPostorder(NodeTraversal, Node, Node)
func (b *TypedScopeCreator_NormalScopeBuilder) VisitPostorder(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_CALL:
		b.CheckForClassDefiningCalls(n)
	case rhino.Token_ASSIGN:
		// Handle initialization of properties.
		// We only allow qualified name declarations of the form
		//   /** @type {number} */ a.b.c = rhs;
		// or, for symbol-typed properties:
		//   /** @type {number} */ a.b[Symbol.something] = rhs;
		// TODO(b/77597706): Ensure that CheckJSDoc warns for JSDoc on assignments not to
		// qualified names, e.g.
		//   /** @type {number} */ [a.b.c] = someArr;
		firstChild := n.GetFirstChild()
		if firstChild.IsGetProp() && firstChild.IsQualifiedName() {
			b.MaybeDeclareQualifiedName(t, n.GetJSDocInfo(), firstChild, n, firstChild.GetNext())
		} else if firstChild.IsGetElem() && firstChild.GetFirstChild().IsQualifiedName() {
			b.MaybeDeclareGetElem(n.GetJSDocInfo(), firstChild, firstChild.GetNext())
		} else if b.outer.undeclaredNamesForClosure.Contains(firstChild) {
			b.DefineAssignAsIfVarDeclaration(n)
		}
	case rhino.Token_CATCH:
		b.DefineCatch(n)
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		b.DefineVars(n)
	case rhino.Token_GETPROP:
		// Handle stubbed properties.
		if parent.IsExprResult() && n.IsQualifiedName() {
			b.MaybeDeclareQualifiedName(t, n.GetJSDocInfo(), n, parent, nil)
		}
	case rhino.Token_GETELEM:
		// Handle stubbed properties.
		if parent.IsExprResult() && n.GetFirstChild().IsQualifiedName() {
			b.MaybeDeclareGetElem(n.GetJSDocInfo(), n, nil)
		}
	case rhino.Token_CLASS:
		// Analyse CLASS child-scopes now because later code in this scope may assign
		// properties to these class-types. We want to ensure declarations within the CLASS have
		// priority.
		b.outer.CreateScope(n, b.currentScope)
	case rhino.Token_EXPR_RESULT:
		names := b.outer.providedNamesFromCall.Get(n)
		if names != nil {
			for _, name := range names {
				b.DeclareProvidedNs(n, name)
			}
		}
	case rhino.Token_EXPORT:
		if n.GetBooleanProp(rhino.Node_EXPORT_DEFAULT) {
			// Define a dummy var for "export default <someExpr>" so that other utilities have
			// access to the type.
			// Problem: according to the debugger, if we have export default SomeName
			// SomeName does not have a type attached so we don't really export it and can't use it
			// as a type later. So we need to call getDeclaredType() instead to get the type.
			declaredType := b.GetDeclaredType(nil, n, n.GetOnlyChild(), nil)
			NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(&b.TypedScopeCreator_AbstractScopeBuilderBase).
				InScope(b.currentScope).
				ForDeclarationNode(n).
				WithType(declaredType).
				ForVariableName(Export_DEFAULT_EXPORT_NAME).
				AllowLaterTypeInference(declaredType == nil).
				DefineSlot()
		}
	default:
	}
} // end NormalScopeBuilder

// TypedScopeCreator_FunctionScopeBuilder is the scope builder subclass for
// function scopes, which only contain bleeding function names and parameter
// names. The main function body is handled by the NormalScopeBuilder on the
// function block.
type TypedScopeCreator_FunctionScopeBuilder struct {
	TypedScopeCreator_AbstractScopeBuilderBase
}

// Java: TypedScopeCreator.FunctionScopeBuilder(TypedScope)
func NewTypedScopeCreator_FunctionScopeBuilder(outer *TypedScopeCreator, scope *TypedScope) *TypedScopeCreator_FunctionScopeBuilder {
	b := &TypedScopeCreator_FunctionScopeBuilder{}
	b.InitTypedScopeCreator_AbstractScopeBuilderBase(b, outer, scope)
	return b
}

// Java: TypedScopeCreator.FunctionScopeBuilder.visitPreorder(NodeTraversal, Node, Node)
func (b *TypedScopeCreator_FunctionScopeBuilder) VisitPreorder(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if parent == nil {
		b.HandleFunctionInputs()
	} else if n.IsFunction() {
		b.DefineFunctionLiteral(n)
	}
}

// HandleFunctionInputs handles bleeding functions and function parameters.
//
// Java: TypedScopeCreator.FunctionScopeBuilder.handleFunctionInputs()
func (b *TypedScopeCreator_FunctionScopeBuilder) HandleFunctionInputs() {
	// Handle bleeding functions. These are defined as function expressions which have a non-empty
	// name, which we declare in the FUNCTION scope. Function declarations are hoisted and are
	// already declared in the containing scope; ignore those.
	fnNode := b.currentScope.GetRootNode()
	fnNameNode := fnNode.GetFirstChild()
	fnName := fnNameNode.GetString()
	if fnName != "" && NodeUtil_IsFunctionExpression(fnNode) {
		NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(&b.TypedScopeCreator_AbstractScopeBuilderBase).
			ForDeclarationNode(fnNameNode).
			ForVariableName(fnName).
			InScope(b.currentScope).
			WithType(fnNode.GetJSType()).
			AllowLaterTypeInference(false).
			DefineSlot()
	}

	b.DeclareParameters(fnNode)
}

// DeclareParameters declares all of a function's parameters inside the
// function's scope.
//
// Java: TypedScopeCreator.FunctionScopeBuilder.declareParameters(Node)
func (b *TypedScopeCreator_FunctionScopeBuilder) DeclareParameters(functionNode *rhino.Node) {
	if NodeUtil_IsBundledGoogModuleCall(functionNode.GetParent()) {
		// Skip declaring 'exports' for a goog.loadModule(function(exports) {.
		// We pretend that any assignments to 'exports' in the body are actually declarations.
		return
	}

	astParameters := functionNode.GetSecondChild()
	var iifeArgumentNode *rhino.Node

	if NodeUtil_IsInvocationTarget(functionNode) {
		iifeArgumentNode = functionNode.GetNext()
	}

	functionType := rhino.JSType_ToMaybeFunctionType(functionNode.GetJSType())
	if functionType != nil {
		// Go-only: language-adapter — the Iterator<Parameter> is an index into
		// the parameter list.
		jsdocParameters := functionType.GetParameters()
		jsdocParametersIndex := 0
		var jsDocParameter *rhino.FunctionType_Parameter
		if jsdocParametersIndex < len(jsdocParameters) {
			jsDocParameter = jsdocParameters[jsdocParametersIndex]
			jsdocParametersIndex++
		}

		for astParameter := astParameters.GetFirstChild(); astParameter != nil; astParameter = astParameter.GetNext() {
			if iifeArgumentNode != nil && iifeArgumentNode.IsSpread() {
				// don't try inferring types from spreads in iifes because we don't know how
				// many items are in the iterable.
				iifeArgumentNode = nil
			}
			var declaredType rhino.JSType
			if jsDocParameter == nil {
				declaredType = b.outer.unknownType
			} else {
				declaredType = jsDocParameter.GetJSType()
			}
			b.DeclareNamesInPositionalParameter(astParameter, declaredType, iifeArgumentNode)
			if jsDocParameter != nil {
				jsDocParameter = nil
				if jsdocParametersIndex < len(jsdocParameters) {
					jsDocParameter = jsdocParameters[jsdocParametersIndex]
					jsdocParametersIndex++
				}
			}
			if iifeArgumentNode != nil {
				iifeArgumentNode = iifeArgumentNode.GetNext()
			}
		}

		b.DeclareTemplatesInFunctionBody(functionNode, functionType)
	}
} // end declareParameters

// Java: TypedScopeCreator.FunctionScopeBuilder.declareTemplatesInFunctionBody(Node, FunctionType)
func (b *TypedScopeCreator_FunctionScopeBuilder) DeclareTemplatesInFunctionBody(functionNode *rhino.Node, functionType rhino.FunctionType) {
	// Add template params to the scope.
	// This lets JSTypeRegistry resolve references to template types within the function body
	// scope. (they were already registered as types by FunctionTypeBuilder, but they need to be
	// in a TypedScope to let JSTypeRegistry correctly handle scoping).
	info := NodeUtil_GetBestJSDocInfo(functionNode)
	templateNames := jutil.NewLinkedHashSet[string]()

	isPossiblePrototypeMethod :=
		!functionType.HasInstanceType() && (info == nil || !info.HasThisType())
	if isPossiblePrototypeMethod {
		// If this is a prototype method, add any @template types from the prototype owner.
		// e.g. given
		//   /** @template T */
		//   class Foo {
		//     bar() {}
		//   }
		// We want to declare 'T' in the scope of 'bar'
		ownerNode := NodeUtil_GetBestLValue(functionNode)
		for _, templateKey := range b.FindOwnerTypeKeys(ownerNode) {
			templateNames.Add(*templateKey.GetReferenceName())
		}
	}
	if info != nil {
		// Add @template parameters from the JSDoc on this function directly. For example, given
		//   /** @template U */
		//   function foo() {}
		// This declares 'U'.
		templateNames.AddAll(info.GetTemplateTypeNames()...)
		templateNames.AddAll(info.GetTypeTransformations().KeySlice()...)
	}
	if templateNames.IsEmpty() {
		return
	}
	input := b.GetCompilerInput()
	voidType := b.outer.GetNativeType(rhino.JSTypeNative_VOID_TYPE)
	// Declare any template names in the function scope. This means that if someone shadows
	// an outer variable FOO with a @template FOO and refers to FOO inside the method, we
	// will treat it as undefined, rather than the correct type, which could lead to weird
	// errors. Ideally we'd have a "don't use me" type that gives an error at use.
	for name := range templateNames.All() {
		if !b.currentScope.CanDeclare(name) {
			b.outer.validator.ExpectUndeclaredVariable(
				typedScopeCreator_sourceNameArg(NodeUtil_GetSourceName(functionNode)),
				input,
				functionNode,
				functionNode.GetParent(),
				b.currentScope.GetVar(name),
				name,
				voidType)
		}
		b.currentScope.Declare(name, functionNode, voidType, input /* inferred= */, false)
	}
}

// DeclareNamesInPositionalParameter declares the name(s) in a positional AST
// parameter in the scope.
//
// astParameter is the positional parameter node; declaredParameterType the
// declared parameter type, or the unknown type if there is none;
// iifeArgumentNode the corresponding argument from the iife, if in an iife.
// e.g. for `(function (x) {}(3);` this would be `3`.
//
// Java: TypedScopeCreator.FunctionScopeBuilder.declareNamesInPositionalParameter(Node, JSType, Node)
func (b *TypedScopeCreator_FunctionScopeBuilder) DeclareNamesInPositionalParameter(
	astParameter *rhino.Node, declaredParameterType rhino.JSType, iifeArgumentNode *rhino.Node) {
	paramType := declaredParameterType
	isInferred := paramType.Equals(b.outer.unknownType)

	if iifeArgumentNode != nil && isInferred {
		argumentName := iifeArgumentNode.GetQualifiedName()
		var argumentVar *TypedVar
		if !(argumentName == nil || b.currentScope.GetParent() == nil) {
			argumentVar = b.currentScope.GetParent().GetVar(*argumentName)
		}
		if argumentVar != nil && !argumentVar.IsTypeInferred() {
			paramType = argumentVar.GetType()
		}
	}

	if paramType == nil {
		paramType = b.outer.unknownType
	}

	switch astParameter.GetToken() {
	case rhino.Token_NAME:
		// function f(x) {}
		b.DeclareSingleParameterName(isInferred, astParameter, paramType)
	case rhino.Token_ITER_REST:
		// function f(...x) {}
		// rest parameter is actually an array of the type specified in the JSDoc
		param := astParameter.GetFirstChild()
		arrayType := b.outer.typeRegistry.GetNativeObjectType(rhino.JSTypeNative_ARRAY_TYPE)
		restParamType := b.outer.typeRegistry.CreateTemplatizedType(arrayType, paramType)
		if param.IsName() {
			b.DeclareSingleParameterName(isInferred, astParameter.GetFirstChild(), restParamType)
		} else {
			// function f(...{length}) {}
			b.DeclareDestructuringParameter(isInferred, param, restParamType)
		}
	case rhino.Token_DEFAULT_VALUE:
		// function f(x = 3) {} or function f([x] = []) {}
		actualParam := astParameter.GetFirstChild()
		if actualParam.IsName() {
			b.DeclareSingleParameterName(isInferred, actualParam, paramType)
		} else {
			b.DeclareDestructuringParameter(isInferred, actualParam, paramType)
		}
	case rhino.Token_ARRAY_PATTERN, rhino.Token_OBJECT_PATTERN:
		// function f([x]) {}
		// function f({x}) {}
		b.DeclareDestructuringParameter(isInferred, astParameter, paramType)
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected function parameter node " + jutil.ToJavaString(astParameter)))
	}
}

// DeclareDestructuringParameter declares all names inside a destructuring
// pattern in a parameter list in the scope if we can find a non-unknown type
// for them.
//
// Unknown typed parameters are always treated as inferred, not declared.
// TypeInference may later give them a better inferred type than unknown,
// but they will never become declared.
//
// NOTE: currently, there are some less-than-ideal aspects to how we do
// this. If the pattern type is an unresolved NamedType, then we can't
// lookup properties on it (to find the individual parameter types) until
// after name resolution. The current state is to defer to TypeInference to
// type those parameters, with the drawback that they are 'inferred', not
// 'declared', and so any type can be assigned to them. In the future we
// will just enforce typing each parameter individually: relevant issue
// https://github.com/google/closure-compiler/issues/1781
//
// Java: TypedScopeCreator.FunctionScopeBuilder.declareDestructuringParameter(boolean, Node, JSType)
func (b *TypedScopeCreator_FunctionScopeBuilder) DeclareDestructuringParameter(
	isInferred bool, pattern *rhino.Node, patternType rhino.JSType) {

	for _, target := range DestructuredTarget_CreateAllNonEmptyTargetsInPattern(
		b.outer.typeRegistry, patternType, pattern) {
		parameterType := target.InferTypeWithoutUsingDefaultValue()

		if target.GetNode().IsDestructuringPattern() {
			b.DeclareDestructuringParameter(isInferred, target.GetNode(), parameterType)
		} else {
			paramName := target.GetNode()
			jutil.CheckStateMsg1(paramName.IsName(), "Expected all parameters to be names, got %s", paramName)

			if parameterType == nil || parameterType.IsUnknownType() {
				paramJSDoc := paramName.GetJSDocInfo()
				if paramJSDoc != nil && paramJSDoc.HasType() {
					// see if the parameter has its own inline JSDoc, and use that unless we already have
					// a type from @param JSDoc.
					// TODO(b/112651122): this should happen inside FunctionTypeBuilder, so that we can
					// check that calls to the function match the inline JSDoc.
					// TODO(b/111523967): we should also report a
					// warning if the inline and non-inline JSDoc conflict.
					parameterType =
						b.outer.typeRegistry.EvaluateTypeExpression(paramJSDoc.GetType(), b.currentScope)
					isInferred = false
				} else {
					// note - these parameters may get better types during TypeInference
					isInferred = true
				}
			}
			b.DeclareSingleParameterName(isInferred, paramName, parameterType)
		}
	}
}

// Java: TypedScopeCreator.FunctionScopeBuilder.declareSingleParameterName(boolean, Node, JSType)
func (b *TypedScopeCreator_FunctionScopeBuilder) DeclareSingleParameterName(isInferred bool, name *rhino.Node, type_ rhino.JSType) {
	NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(&b.TypedScopeCreator_AbstractScopeBuilderBase).
		ForDeclarationNode(name).
		ForVariableName(name.GetString()).
		InScope(b.currentScope).
		WithType(type_).
		AllowLaterTypeInference(isInferred).
		DefineSlot()
} // end FunctionScopeBuilder

// TypedScopeCreator_ClassScopeBuilder is the scope builder subclass for
// class scopes (which only contain a bleeding class name), member field def
// scopes, and RHS computed field def scopes (the latter two of which have
// `this` and `super` properties). Methods are handled by
// FunctionScopeBuilder and NormalScopeBuilder for the bodies.
type TypedScopeCreator_ClassScopeBuilder struct {
	TypedScopeCreator_AbstractScopeBuilderBase
}

// Java: TypedScopeCreator.ClassScopeBuilder(TypedScope)
func NewTypedScopeCreator_ClassScopeBuilder(outer *TypedScopeCreator, scope *TypedScope) *TypedScopeCreator_ClassScopeBuilder {
	b := &TypedScopeCreator_ClassScopeBuilder{}
	b.InitTypedScopeCreator_AbstractScopeBuilderBase(b, outer, scope)
	return b
}

// Java: TypedScopeCreator.ClassScopeBuilder.visitPreorder(NodeTraversal, Node, Node)
func (b *TypedScopeCreator_ClassScopeBuilder) VisitPreorder(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	// These are not descended into, so must be done preorder
	if !n.IsFunction() {
		return
	}

	if NodeUtil_IsEs6Constructor(n) {
		// Constructor has already been analyzed, so pull that here.
		b.outer.SetDeferredType(n, b.currentScope.GetRootNode().GetJSType())
	} else {
		b.DefineFunctionLiteral(n)
	}
}

// Java: TypedScopeCreator.ClassScopeBuilder.visitPostorder(NodeTraversal, Node, Node)
func (b *TypedScopeCreator_ClassScopeBuilder) VisitPostorder(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsName() &&
		parent == b.currentScope.GetRootNode() &&
		NodeUtil_IsClassExpression(parent) {
		// Declare bleeding class name in scope.  Pull the type off the AST.
		jutil.CheckState(n.GetString() != "") // anonymous classes have EMPTY nodes, not NAME
		NewTypedScopeCreator_AbstractScopeBuilder_SlotDefiner(&b.TypedScopeCreator_AbstractScopeBuilderBase).
			ForDeclarationNode(n).
			ReadVariableNameFromDeclarationNode().
			InScope(b.currentScope).
			WithType(parent.GetJSType()).
			AllowLaterTypeInference(false).
			DefineSlot()
	} else if NodeUtil_IsEs6ConstructorMemberFunctionDef(n) {
		// Ignore "constructor" since it has special handling in `createClassTypeFromNodes()`.
	} else if n.IsMemberFunctionDef() && parent.IsClassMembers() {
		b.DefineMemberFunction(n)
	} else if n.IsMemberFieldDef() {
		// public fields are roots of their own scope so the parent doesn't get passed into
		// visitPostorder
		b.DefineMemberField(n)
	} else if (n.IsGetterDef() || n.IsSetterDef()) && parent.IsClassMembers() {
		b.DefineGetterSetter(n)
	} else if n.IsComputedFieldDef() ||
		(n.IsComputedProp() &&
			n.GetBooleanProp(rhino.Node_COMPUTED_PROP_METHOD) &&
			parent.IsClassMembers()) {
		b.DefineComputedMemberField(n)
	}
}

// Java: TypedScopeCreator.ClassScopeBuilder.defineMemberFunction(Node)
func (b *TypedScopeCreator_ClassScopeBuilder) DefineMemberFunction(n *rhino.Node) {
	ownerType := b.DetermineOwnerTypeForClassMember(n)
	ownerType.DefineDeclaredProperty(n.GetString(), n.GetLastChild().GetJSType(), n)
}

// Java: TypedScopeCreator.ClassScopeBuilder.defineMemberField(Node)
func (b *TypedScopeCreator_ClassScopeBuilder) DefineMemberField(n *rhino.Node) {
	rhs := n.GetLastChild()
	declaredType := b.GetDeclaredType(n.GetJSDocInfo(), n, rhs, nil)
	// When there's no JSDoc type declaration for a field:
	// If we can infer a type other than UNKNOWN for an initial value, use that as the field type
	// Otherwise, set the field's type to the ALL type ('*')
	if declaredType == nil {
		var rhsType rhino.JSType
		if rhs != nil {
			rhsType = b.TypedScopeCreator_AbstractScopeBuilderBase.GetDeclaredRValueType(nil, rhs)
		}
		if rhsType != nil && !rhsType.IsUnknownType() {
			declaredType = rhsType
		} else {
			declaredType = b.outer.GetNativeType(rhino.JSTypeNative_ALL_TYPE)
		}
	}

	ownerType := b.DetermineOwnerTypeForClassMember(n)
	propName := n.GetString()
	ownerType.DefineDeclaredProperty(propName, declaredType, n)
	n.SetJSType(declaredType)
}

// Java: TypedScopeCreator.ClassScopeBuilder.defineComputedMemberField(Node)
func (b *TypedScopeCreator_ClassScopeBuilder) DefineComputedMemberField(n *rhino.Node) {
	key := n.GetFirstChild()
	if !key.IsQualifiedName() {
		return
	}
	keyType := b.currentScope.LookupQualifiedName(key.GetQualifiedNameObject())
	if keyType == nil || !keyType.IsKnownSymbolValueType() {
		return
	}
	rhs := n.GetLastChild()
	declaredType := b.GetDeclaredType(n.GetJSDocInfo(), n, rhs, nil)
	// When there's no JSDoc type declaration for a field:
	// If we can infer a type other than UNKNOWN for an initial value, use that as the field type
	// Otherwise, set the field's type to the ALL type ('*')
	if declaredType == nil {
		var rhsType rhino.JSType
		if rhs != nil {
			rhsType = b.TypedScopeCreator_AbstractScopeBuilderBase.GetDeclaredRValueType(nil, rhs)
		}
		if rhsType != nil && !rhsType.IsUnknownType() {
			declaredType = rhsType
		} else {
			declaredType = b.outer.GetNativeType(rhino.JSTypeNative_ALL_TYPE)
		}
	}

	ownerType := b.DetermineOwnerTypeForClassMember(n)
	ownerType.DefineDeclaredPropertyByKey(
		rhino.NewProperty_SymbolKey(keyType.ToMaybeKnownSymbolType()), declaredType, n)
	n.SetJSType(declaredType)
}

// Java: TypedScopeCreator.ClassScopeBuilder.defineGetterSetter(Node)
func (b *TypedScopeCreator_ClassScopeBuilder) DefineGetterSetter(n *rhino.Node) {
	name := n.GetString()
	methodType := n.GetLastChild().GetJSType().ToMaybeFunctionType()

	var propertyType rhino.JSType
	switch n.GetToken() {
	case rhino.Token_GETTER_DEF:
		// TODO(sdh): consider only falling back on unknown if the function body is empty?
		// But
		// we need to not report a conflicting type error if there's different unknowns.
		if methodType.IsReturnTypeInferred() {
			propertyType = b.outer.unknownType
		} else {
			propertyType = methodType.GetReturnType()
		}
	case rhino.Token_SETTER_DEF:
		if len(methodType.GetParameters()) == 0 {
			propertyType = b.outer.unknownType
		} else {
			propertyType = methodType.GetParameters()[0].GetJSType()
		}
	default:
		panic(jutil.NewAssertionError(n.ToStringTree()))
	}

	ownerType := b.DetermineOwnerTypeForClassMember(n)
	// TODO(b/116797078): correctly model getters/setters and stop treating this as a normal
	// property.
	ownerType.DefineDeclaredProperty(name, propertyType, n)
}

// DetermineOwnerTypeForClassMember returns the owner type for a class member
// function, getter, or setter.
//
// For a member on class C, this is either `C` for a static member or
// `C.prototype` for a nonstatic member.
//
// Java: TypedScopeCreator.ClassScopeBuilder.determineOwnerTypeForClassMember(Node)
func (b *TypedScopeCreator_ClassScopeBuilder) DetermineOwnerTypeForClassMember(member *rhino.Node) rhino.ObjectType {
	// MEMBER_FUNCTION_DEF -> CLASS_MEMBERS -> CLASS or
	// MEMBER_FIELD_DEF -> CLASS_MEMBERS -> CLASS or
	// GETTER_DEF -> CLASS_MEMBERS -> CLASS or
	// SETTER_DEF -> CLASS_MEMBERS -> CLASS
	ownerNode := member.GetGrandparent()
	jutil.CheckState(ownerNode.IsClass())
	ownerType := ownerNode.GetJSType().ToMaybeFunctionType()
	if member.IsStaticMember() {
		return ownerType
	} else if member.IsMemberFieldDef() {
		return ownerType.GetInstanceType()
	} else {
		jutil.CheckStateMsg1(
			member.IsMemberFunctionDef() ||
				member.IsGetterDef() ||
				member.IsSetterDef() ||
				member.IsComputedProp() ||
				member.IsComputedFieldDef(),
			"%s",
			member)
		return ownerType.GetPrototype()
	}
} // end ClassScopeBuilder

// TypedScopeCreator_FirstOrderFunctionAnalyzer does a first-order function
// analysis that just looks at simple things like what variables are
// escaped, and whether 'this' is used.
//
// The syntactic scopes created in this traversal are also stored for later
// use.
type TypedScopeCreator_FirstOrderFunctionAnalyzer struct {
	NodeTraversal_AbstractScopedCallbackBase

	// Go-only: language-adapter — the enclosing TypedScopeCreator instance of
	// this inner class.
	outer *TypedScopeCreator
}

// Java: TypedScopeCreator.FirstOrderFunctionAnalyzer()
func NewTypedScopeCreator_FirstOrderFunctionAnalyzer(outer *TypedScopeCreator) *TypedScopeCreator_FirstOrderFunctionAnalyzer {
	return &TypedScopeCreator_FirstOrderFunctionAnalyzer{outer: outer}
}

// Java: TypedScopeCreator.FirstOrderFunctionAnalyzer.enterScope(NodeTraversal)
func (a *TypedScopeCreator_FirstOrderFunctionAnalyzer) EnterScope(t *NodeTraversal) {
	scope := t.GetScope()
	root := scope.GetRootNode()
	for symbol := range scope.GetVarIterable() {
		a.outer.reservedNamesForScope.Put(root, symbol.GetName())
	}
}

// Java: TypedScopeCreator.FirstOrderFunctionAnalyzer.visit(NodeTraversal, Node, Node)
func (a *TypedScopeCreator_FirstOrderFunctionAnalyzer) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if t.InGlobalScope() {
		// The first-order function analyzer looks at two types of variables:
		//
		// 1) Local variables that are assigned in inner scopes ("escaped vars")
		//
		// 2) Local variables that are assigned more than once.
		//
		// We treat all global variables as escaped by default, so there's
		// no reason to do this extra computation for them.
		return
	}

	containerScope := t.GetClosestContainerScope().(*Scope)

	// Record function with returns or arrow functions without bodies
	if (n.IsReturn() && n.HasChildren()) || NodeUtil_IsBlocklessArrowFunctionResult(n) {
		a.outer.functionsWithNonEmptyReturns.Add(containerScope.GetRootNode())
	}

	// Be careful of bleeding functions, which create variables
	// in the inner scope, not the scope where the name appears.
	if n.IsName() && NodeUtil_IsLValue(n) && !NodeUtil_IsBleedingFunctionName(n) {
		name := n.GetString()
		scope := t.GetScope()
		var_ := scope.GetVar(name)
		// TODO(sdh): consider checking hasSameHoistScope instead of container scope here and
		// below. This will detect function(a) { a.foo = bar } as an escaped qualified name,
		// which seems like the right thing to do (but could possibly break things?)
		// Doing so will allow removing the warning on TypeCheckTest#testIssue1024b.
		if var_ != nil {
			ownerScope := var_.GetScope().(*Scope)
			if ownerScope.IsLocal() {
				scopedName := ScopedName_Of(name, ownerScope.GetRootNode())
				a.outer.assignedVarNames.Add(scopedName)
				if !containerScope.HasSameContainerScope(ownerScope) {
					a.outer.escapedVarNames.Add(scopedName)
				}
			}
		}
	} else if n.IsGetProp() && n.IsUnscopedQualifiedName() && NodeUtil_IsLValue(n) {
		name := NodeUtil_GetRootOfQualifiedName(n).GetString()
		scope := t.GetScope()
		var_ := scope.GetVar(name)
		if var_ != nil {
			ownerScope := var_.GetScope().(*Scope)
			if ownerScope.IsLocal() && !containerScope.HasSameContainerScope(ownerScope) {
				a.outer.escapedVarNames.Add(ScopedName_Of(*n.GetQualifiedName(), ownerScope.GetRootNode()))
			}
		}
	}
}

// Java: TypedScopeCreator.containingGoogModuleIdOf(TypedScope)
func TypedScopeCreator_ContainingGoogModuleIdOf(scope *TypedScope) *string {
	module := scope.GetModule()
	if module == nil {
		parent := scope.GetParent()
		if parent != nil {
			return TypedScopeCreator_ContainingGoogModuleIdOf(parent)
		}
		return nil
	}

	// Stop recursing once we've hit a module scope.
	metadata := module.Metadata()
	jutil.CheckStateMsg1(metadata.IsModule(), "%s", metadata)

	/*
	 * This module may not have a goog.module/goog.declareModuleId. Also don't crash if it's
	 * malformed with multiple module IDs.
	 */
	for namespace := range metadata.GoogNamespaces().All() {
		return jutil.StringPtr(namespace)
	}
	return nil
}

// Java: TypedScopeCreator.extractKnownSymbolKey(TypedScope, Node)
func TypedScopeCreator_ExtractKnownSymbolKey(currentScope *TypedScope, key *rhino.Node) *rhino.KnownSymbolType {
	if !key.IsQualifiedName() {
		return nil
	}
	type_ := currentScope.LookupQualifiedName(key.GetQualifiedNameObject())
	if type_ != nil {
		return type_.ToMaybeKnownSymbolType()
	}
	return nil
}

var (
	_ ScopeCreator                                  = (*TypedScopeCreator)(nil)
	_ rhino.StaticSymbolTable[*TypedVar, *TypedVar] = (*TypedScopeCreator)(nil)
	_ TypedScopeCreator_AbstractScopeBuilder        = (*TypedScopeCreator_NormalScopeBuilder)(nil)
	_ TypedScopeCreator_AbstractScopeBuilder        = (*TypedScopeCreator_FunctionScopeBuilder)(nil)
	_ TypedScopeCreator_AbstractScopeBuilder        = (*TypedScopeCreator_ClassScopeBuilder)(nil)
	_ NodeTraversal_Callback                        = (*TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable)(nil)
	_ NodeTraversal_ScopedCallback                  = (*TypedScopeCreator_FirstOrderFunctionAnalyzer)(nil)
)
