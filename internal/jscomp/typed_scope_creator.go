// Port of com/google/javascript/jscomp/TypedScopeCreator.java @ closure-compiler v20261021

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

	"github.com/paullepoulpe/jsgo/internal/ast"
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
type typedScopeCreator_ScopedNameKey struct {
	name      string
	scopeRoot rhino.Node
}

const (
	typedScopeCreator_VarFlagEscaped           uint8 = 1 << 0
	typedScopeCreator_VarFlagNotSingleAssigned uint8 = 1 << 1
)

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
	memoized                 map[rhino.Node]*TypedScope
	memoizedOrder            []*TypedScope

	// Maps from scope root to declared variable names. Populated by FirstOrderFunctionAnalyzer to
	// reserve names before the TypedScope is populated.
	reservedNamesForScope map[rhino.Node][]string

	scopeCount int

	// Set of functions with non-empty returns, for passing to FunctionTypeBuilder.
	functionsWithNonEmptyReturns map[rhino.Node]struct{}
	// Tracks escaped variables (simple and qualified) and local variables not assigned exactly once.
	varFlags map[typedScopeCreator_ScopedNameKey]uint8

	// For convenience
	unknownType rhino.ObjectType

	// All names imported through goog.requireType. Resolve these after all scopes are created.
	weakImports []*TypedScopeCreator_WeakModuleImport

	unresolvedNodes []rhino.Node

	// Set of NAME, GETPROP, and STRING_KEY lvalues which should be treated as const declarations when
	// assigned. Treat simple names in this list as if they were declared `const`. E.g. treat `exports
	// = class {};` as `const exports = class {};`. Treat GETPROP and STRING_KEY nodes as if they were
	// annotated @const.
	undeclaredNamesForClosure map[rhino.Node]struct{}

	// Maps EXPR_RESULT nodes from goog.provides to all implicitly provided names from the call
	providedNamesFromCall map[rhino.Node][]*ProcessClosureProvidesAndRequires_ProvidedName
	providedNamesCount    int

	builderPool []*TypedScopeCreator_scopeBuilder

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

	moduleLocalNode  rhino.Node  // the NAME or IMPORT_STAR representing the assignee
	scopedImport     ScopedName  // the name (possibly qualified) & scope being imported
	localModuleScope *TypedScope // the scope containing the import
}

// Java: TypedScopeCreator.WeakModuleImport(Node, ScopedName, TypedScope)
func NewTypedScopeCreator_WeakModuleImport(outer *TypedScopeCreator, moduleLocalNode rhino.Node, scopedImport ScopedName, localModuleScope *TypedScope) *TypedScopeCreator_WeakModuleImport {
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
	resolvedScope := w.outer.memoized[w.scopedImport.GetScopeRoot()]
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

// typedScopeRValueInfo stores the type and qualified name for a destructuring
// rvalue on the stack.
type typedScopeRValueInfo struct {
	type_         rhino.JSType
	qualifiedName rhino.QualifiedName
}

// Java: TypedScopeCreator(AbstractCompiler)
func NewTypedScopeCreator(compiler AbstractCompiler) *TypedScopeCreator {
	return NewTypedScopeCreatorWithCodingConvention(compiler, compiler.GetCodingConvention())
}

// Java: TypedScopeCreator(AbstractCompiler, CodingConvention)
func NewTypedScopeCreatorWithCodingConvention(compiler AbstractCompiler, codingConvention CodingConvention) *TypedScopeCreator {
	c := &TypedScopeCreator{}
	c.memoized = make(map[rhino.Node]*TypedScope)
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
	c.processClosurePrimitives = c.metadataMap.ModulesByGoogNamespaceMap().Size() > 0

	// Reset state to empty collections.
	c.ClearCommonState()
	return c
}

func (c *TypedScopeCreator) putMemoized(root rhino.Node, scope *TypedScope) {
	if prev, exists := c.memoized[root]; !exists {
		c.memoizedOrder = append(c.memoizedOrder, scope)
	} else {
		for i, existing := range c.memoizedOrder {
			if existing == prev {
				c.memoizedOrder[i] = scope
				break
			}
		}
	}
	c.memoized[root] = scope
}

func (c *TypedScopeCreator) getMemoizedScope(n rhino.Node) *TypedScope {
	return c.memoized[n]
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
	return func(yield func(*TypedVar) bool) {
		for _, s := range c.memoizedOrder {
			for v := range s.GetAllSymbols() {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// GetNodeToScopeMapper returns a function mapping a scope root node to a
// TypedScope.
//
// This method mostly exists in lieu of an interface representing root node
// -> scope.
//
// Java: TypedScopeCreator.getNodeToScopeMapper()
func (c *TypedScopeCreator) GetNodeToScopeMapper() func(rhino.Node) *TypedScope {
	return c.getMemoizedScope
}

// Java: TypedScopeCreator.getAllMemoizedScopes()
func (c *TypedScopeCreator) GetAllMemoizedScopes() []*TypedScope {
	// Return scopes in reverse order of creation so that IIFEs will
	// come before the global scope.
	scopes := slices.Clone(c.memoizedOrder)
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
func (c *TypedScopeCreator) CreateScopeOfNode(n rhino.Node) *TypedScope {
	s := c.memoized[n]
	if s != nil {
		return s
	}
	return c.CreateScope(n, c.CreateScopeOfNode(NodeUtil_GetEnclosingScopeRoot(n.GetParent()))).(*TypedScope)
}

// CreateScope creates a scope with all types declared. Declares newly
// discovered types and type properties in the type registry.
//
// Java: TypedScopeCreator.createScope(Node, AbstractScope)
func (c *TypedScopeCreator) CreateScope(root rhino.Node, parent AbstractScope) AbstractScope {
	_, parentIsTyped := parent.(*TypedScope)
	jutil.CheckArgument(parent == nil || parentIsTyped)
	var typedParent *TypedScope
	if parent != nil {
		typedParent = parent.(*TypedScope)
	}

	scope := c.memoized[root]
	if scope != nil {
		jutil.CheckState(typedParent == scope.GetParent())
	} else {
		scope = c.CreateScopeInternal(root, typedParent)
		c.putMemoized(root, scope)
	}
	return scope
}

// Java: TypedScopeCreator.createScopeInternal(Node, TypedScope)
func (c *TypedScopeCreator) CreateScopeInternal(root rhino.Node, typedParent *TypedScope) *TypedScope {
	jutil.CheckStateMsg1(
		c.stage == TypedScopeCreator_Stage_BUILDING,
		"Cannot create scope for %s after TypedScopeCreator is frozen",
		root)
	// Constructing the global scope is very different than constructing
	// inner scopes, because only global scopes can contain named classes that
	// show up in the type registry.
	var newScope *TypedScope

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
		reservedNames := c.reservedNamesForScope[root]
		if reservedNames != nil {
			delete(c.reservedNamesForScope, root)
		}
		if module != nil && module.Metadata().IsGoogModule() {
			// TypedScopeCreator treats default export assignments, like `exports = class {};`, as
			// declarations. However, the untyped scope only contains an implicit slot for `exports`.
			if !slices.Contains(reservedNames, "exports") {
				reservedNames = append(reservedNames, "exports")
				c.varFlags[typedScopeCreator_ScopedNameKey{name: "exports", scopeRoot: root}] |= typedScopeCreator_VarFlagNotSingleAssigned
			}
		} else if len(reservedNames) > 0 && root.IsFunction() && NodeUtil_IsBundledGoogModuleCall(root.GetParent()) {
			// Pretend that 'exports' is declared in the block of goog.loadModule
			// functions, not the function scope. See the above comment for why.
			for i, rn := range reservedNames {
				if rn == "exports" {
					copy(reservedNames[i:], reservedNames[i+1:])
					reservedNames[len(reservedNames)-1] = ""
					reservedNames = reservedNames[:len(reservedNames)-1]
					break
				}
			}
		}

		newScope = newTypedScopeOwnedReservedNames(typedParent, root, reservedNames, module)
	}

	var kind scopeBuilderKind
	if root.IsFunction() {
		kind = scopeBuilderFunction
	} else if root.IsClass() || root.IsMemberFieldDef() || root.IsComputedFieldDef() {
		kind = scopeBuilderClass
	} else {
		kind = scopeBuilderNormal
		if root.IsModuleBody() {
			// Store the module scope on the CompilerInput object so it can be found and used for
			// program analysis purposes.
			// DefaultPassConfig is responsible for removing it when it is no longer needed.
			inputId := jutil.CheckNotNullPtr(NodeUtil_GetInputId(root))
			compilerInput := jutil.CheckNotNullPtr(c.compiler.GetInput(inputId))
			compilerInput.SetTypedScope(newScope)
		}
	}
	scopeBuilder := c.acquireScopeBuilder(kind, newScope)
	scopeBuilder.Build()
	c.releaseScopeBuilder(scopeBuilder)
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

	newScope.scanDone = true
	return newScope
}

func (c *TypedScopeCreator) acquireScopeBuilder(kind scopeBuilderKind, scope *TypedScope) *TypedScopeCreator_scopeBuilder {
	var b *TypedScopeCreator_scopeBuilder
	if n := len(c.builderPool); n > 0 {
		b = c.builderPool[n-1]
		c.builderPool = c.builderPool[:n-1]
	} else {
		b = &TypedScopeCreator_scopeBuilder{}
	}
	b.init(kind, c, scope)
	return b
}

func (c *TypedScopeCreator) releaseScopeBuilder(b *TypedScopeCreator_scopeBuilder) {
	b.reset()
	c.builderPool = append(c.builderPool, b)
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
func (c *TypedScopeCreator) MarkGoogModuleExportsAsConst(moduleBody rhino.Node) {
	// TODO(lharker): Use the source nodes from the Bindings once we no longer rewrite before
	// typechecking. This is not feasible currently because a few places will rewrite exports = ...
	for statement := moduleBody.GetFirstChild(); !statement.IsNil(); statement = statement.GetNext() {
		exprStmt, ok := ast.AsExprStmt(statement)
		if !ok {
			continue
		}
		assign, ok := exprStmt.Expr().AsAssign()
		if !ok {
			continue
		}
		lhs := assign.Target().Node()
		rhs := assign.Value().Node()
		if lhs.MatchesName("exports") {
			c.undeclaredNamesForClosure[lhs] = struct{}{}
			// If this is full of named exports, add all the string key nodes.
			if NodeUtil_IsNamedExportsLiteral(rhs) {
				for member := range ast.MustObjectLit(rhs).Members() {
					c.undeclaredNamesForClosure[member.Node()] = struct{}{}
				}
			}
		} else if gp, ok := ast.AsGetProp(lhs); ok && gp.Object().Node().MatchesName("exports") {
			c.undeclaredNamesForClosure[lhs] = struct{}{}
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
func (c *TypedScopeCreator) GatherAllProvides(root rhino.Node) {
	if !c.processClosurePrimitives {
		return
	}

	externs := root.GetFirstChild()
	js := root.GetSecondChild()
	providedNames :=
		NewProcessClosureProvidesAndRequires(c.compiler /* preserveGoogProvidesAndRequires= */, true).
			CollectProvidedNamesSlice(externs, js)

	for _, name := range providedNames {
		if name.HasImplicitInitialization() {
			continue
		}
		if !name.GetCandidateDefinition().IsNil() {
			// This name will be defined eventually in the source code.
			firstDefinitionNode := name.GetCandidateDefinition()
			if exprStmt, ok := ast.AsExprStmt(firstDefinitionNode); ok {
				if assign, ok := exprStmt.Expr().AsAssign(); ok && assign.Target().Node().IsName() {
					// Treat assignments of provided names as declarations.
					c.undeclaredNamesForClosure[assign.Target().Node()] = struct{}{}
				}
			}
		} else if !name.GetFirstProvideCall().IsNil() &&
			NodeUtil_IsExprCall(name.GetFirstProvideCall()) &&
			!name.IsFromLegacyModule() {
			// This name is implicitly created by a goog.provide call; declare it in the scope once
			// reaching the provide call. The exception is legacy goog.modules, which are declared
			// once leaving the module.
			if c.providedNamesFromCall == nil {
				c.providedNamesFromCall = make(map[rhino.Node][]*ProcessClosureProvidesAndRequires_ProvidedName)
			}
			call := name.GetFirstProvideCall()
			c.providedNamesFromCall[call] = append(c.providedNamesFromCall[call], name)
			c.providedNamesCount++
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
func (c *TypedScopeCreator) CreateInitialScope(root rhino.Node) *TypedScope {
	jutil.CheckArgumentMsg1(root.IsRoot(), "%s", root)

	// Gather global information used in typed scope creation.
	scopeCreator := NewSyntacticScopeCreator(c.compiler)

	analyzer := NewTypedScopeCreator_FirstOrderFunctionAnalyzer(c)
	analyzer.identifyEnumsAndTypedefs.registry = c.typeRegistry

	NodeTraversal_Builder_().
		SetCompiler(c.compiler).
		SetCallback(analyzer).
		SetScopeCreator(scopeCreator).
		TraverseRoots(root.GetFirstChild(), root.GetLastChild())

	c.GatherAllProvides(root)

	globalVarCap := 48 + analyzer.globalDeclEstimate + c.providedNamesCount
	s := newGlobalTypedScopeWithCapacity(root, globalVarCap)
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

	// Memoize the global scope so that module scope creation can access it. See
	// AbstractScopeBuilder#shouldTraverse - modules are traversed early, as if they were always
	// executed when control flow reaches the module body.
	if len(c.memoized) == 0 && c.scopeCount > 8 {
		c.memoized = make(map[rhino.Node]*TypedScope, c.scopeCount+16)
		c.memoizedOrder = make([]*TypedScope, 0, c.scopeCount+16)
	}
	c.putMemoized(root, s)

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
	scope.Declare(name, rhino.NilNode, t, nil, false)
}

// SetDeferredType sets the type for a node now.
//
// If the type is unresolved, enqueue it to be updated with a resolved type
// later.
//
// Java: TypedScopeCreator.setDeferredType(Node, JSType)
func (c *TypedScopeCreator) SetDeferredType(node rhino.Node, type_ rhino.JSType) {
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
	c.reservedNamesForScope = nil
	c.scopeCount = 0
	c.functionsWithNonEmptyReturns = nil
	c.varFlags = nil
	c.weakImports = nil
	c.unresolvedNodes = nil
	c.undeclaredNamesForClosure = nil
	c.providedNamesFromCall = nil
	c.providedNamesCount = 0
	c.builderPool = nil
	c.stage = TypedScopeCreator_Stage_FROZEN

	// we must keep the 'memoized' set of TypedScopes around, since later passes will use it.
}

// Java: TypedScopeCreator.clearCommonState()
func (c *TypedScopeCreator) ClearCommonState() {
	c.reservedNamesForScope = make(map[rhino.Node][]string)
	c.scopeCount = 0
	c.functionsWithNonEmptyReturns = make(map[rhino.Node]struct{})
	c.varFlags = make(map[typedScopeCreator_ScopedNameKey]uint8)
	c.weakImports = nil
	c.unresolvedNodes = nil
	c.undeclaredNamesForClosure = make(map[rhino.Node]struct{})
	c.providedNamesFromCall = nil
	c.providedNamesCount = 0
	c.builderPool = nil
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

var identifyEnumsAndTypedefsVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_LET)
	s.Add(rhino.Token_CONST)
	s.Add(rhino.Token_VAR)
	s.Add(rhino.Token_EXPR_RESULT)
	return s
}()

var (
	_ nodeTraversal_noVisitLeaves = (*TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable)(nil)
	_ nodeTraversal_visitTokens   = (*TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable)(nil)
)

func (*TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable) noVisitLeaves() {}

func (*TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable) visitTokens() *ast.TokenSet {
	return &identifyEnumsAndTypedefsVisitTokens
}

// Java: TypedScopeCreator.IdentifyEnumsAndTypedefsAsNonNullable.visit(NodeTraversal, Node, Node)
func (v *TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable) Visit(t *NodeTraversal, node rhino.Node, parent rhino.Node) {
	switch node.GetToken() {
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_VAR:
		for decl := range ast.MustNameDeclaration(node).Declarators() {
			// TODO(b/116853368): make this work for destructuring aliases as well.
			child := decl.Node()
			var rvalue rhino.Node
			if name, ok := decl.AsName(); ok {
				if init, ok := name.Init(); ok {
					rvalue = init.Node()
				}
			} else {
				rvalue = child.GetFirstChild()
			}
			v.IdentifyEnumOrTypedefDeclaration(
				t, child, rvalue, NodeUtil_GetBestJSDocInfo(child))
		}
	case rhino.Token_EXPR_RESULT:
		expr := ast.MustExprStmt(node).Expr()
		if assign, ok := expr.AsAssign(); ok {
			aNode := assign.Node()
			v.IdentifyEnumOrTypedefDeclaration(
				t, assign.Target().Node(), assign.Value().Node(), aNode.GetJSDocInfo())
		} else if gp, ok := expr.AsGetProp(); ok {
			gpNode := gp.Node()
			v.IdentifyEnumOrTypedefDeclaration(
				t, gpNode /* rvalue= */, rhino.NilNode, gpNode.GetJSDocInfo())
		}
	default:
	}
}

// Java: TypedScopeCreator.IdentifyEnumsAndTypedefsAsNonNullable.identifyEnumOrTypedefDeclaration(NodeTraversal, Node, Node, JSDocInfo)
func (v *TypedScopeCreator_IdentifyEnumsAndTypedefsAsNonNullable) IdentifyEnumOrTypedefDeclaration(
	t *NodeTraversal, nameNode rhino.Node, rvalue rhino.Node, info *rhino.JSDocInfo) {
	nameQname, ok := nameNode.GetQualifiedNameString()
	if !ok {
		return
	}
	if info != nil && info.HasEnumParameterType() {
		v.registry.IdentifyNonNullableName(t.resolveLookupScope(len(t.scopes)-1), nameQname)
	} else if info != nil && info.HasTypedefType() {
		v.registry.IdentifyNonNullableName(t.resolveLookupScope(len(t.scopes)-1), nameQname)
	} else if !rvalue.IsNil() && NodeUtil_IsConstantDeclaration(info, nameNode) {
		if rvalQname, ok := rvalue.GetQualifiedNameString(); ok {
			scope := t.resolveLookupScope(len(t.scopes) - 1)
			if v.registry.IsNonNullableName(scope, rvalQname) {
				v.registry.IdentifyNonNullableName(scope, nameQname)
			}
		}
	}
}

// Java: TypedScopeCreator.getNativeType(JSTypeNative)
func (c *TypedScopeCreator) GetNativeType(nativeType rhino.JSTypeNative) rhino.JSType {
	return c.typeRegistry.GetNativeType(nativeType)
}
