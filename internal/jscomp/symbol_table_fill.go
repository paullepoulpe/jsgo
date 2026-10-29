// Port of com/google/javascript/jscomp/SymbolTable.java @ closure-compiler v20261021

/*
 * Copyright 2011 The Closure Compiler Authors.
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
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var symbolTable_logger = jutil.Logger_GetLogger("com.google.javascript.jscomp.SymbolTable")

type externSymbolPosKey struct {
	sourceFile string
	position   int
}

// FillNamespaceReferences fills in namespace references. Not all symbol
// tables record references to "namespace" objects. For example, if you have:
// goog.dom.DomHelper = function() {}; The symbol table may not record that as
// a reference to "goog.dom", because that would be redundant.
//
// Java: SymbolTable.fillNamespaceReferences()
func (t *SymbolTable) FillNamespaceReferences() {
	for _, symbol := range t.GetAllSymbols() {
		qName := symbol.GetName()
		// Go-only: byte indices of the ASCII '.' (only used to cut the string).
		rootIndex := strings.IndexByte(qName, '.')
		if rootIndex == -1 {
			continue
		}

		root := symbol.scope.GetQualifiedSlot(qName[:rootIndex])
		if root == nil {
			// In theory, this should never happen, but we fail quietly anyway
			// just to be safe.
			continue
		}

		for ref := range t.GetReferences(symbol) {
			currentNode := ref.GetNode()
			if !currentNode.IsQualifiedName() {
				continue
			}

			for gp, ok := ast.AsGetProp(currentNode); ok; gp, ok = ast.AsGetProp(currentNode) {
				currentNode = gp.Object().Node()

				name, hasName := currentNode.GetQualifiedNameString()
				if !hasName {
					continue
				}
				namespace := t.isAnySymbolDeclared(name, currentNode, root.scope)
				if namespace == nil {
					namespace = root.scope.GetQualifiedSlot(name)
				}

				if namespace == nil && root.scope.IsGlobalScope() {
					// Originally UNKNOWN_TYPE has been always used for namespace symbols even though
					// compiler does have type information attached to a node. Unclear why. Changing code to
					// property type mostly works. It only fails on Foo.prototype cases for some reason.
					// It's pretty rare case when Foo.prototype defined in global scope though so for now
					// just carve it out.
					nodeType := currentNode.GetJSType()
					var symbolType rhino.JSType
					if nodeType != nil && !nodeType.IsFunctionPrototypeType() {
						symbolType = nodeType
					} else {
						symbolType = t.registry.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
					}
					namespace =
						t.declareSymbol(
							name, symbolType, true, root.scope, currentNode, nil /* JsDoc info */)
				}

				if namespace != nil {
					namespace.DefineReferenceAt(currentNode)
				}
			}
		}
	}
}

// Java: SymbolTable.fillPropertyScopes()
func (t *SymbolTable) FillPropertyScopes() {
	// Collect all object symbols.
	// All symbols that came from goog.module are collected separately because they will have to
	// be processed first. See explanation below.
	types := []*SymbolTable_Symbol{}
	exports := []*SymbolTable_Symbol{}

	// Create a property scope for each named type and each anonymous object,
	// and populate it with that object's properties.
	//
	// We notably don't want to create a property scope for 'x' in
	// var x = new Foo();
	// where x is just an instance of another type.
	for _, sym := range t.GetAllSymbols() {
		if t.needsPropertyScope(sym) {
			name := sym.GetName()
			if name == "exports" {
				exports = append(exports, sym)
			} else {
				types = append(types, sym)
			}
		}
	}

	// The order of operations here is significant.
	//
	// When we add properties to Foo, we'll remove Foo.prototype from
	// the symbol table and replace it with a fresh symbol in Foo's
	// property scope. So the symbol for Foo.prototype in
	// {@code instances} will be stale.
	//
	// If we order them in reverse lexicographical order symbols x.y and x will be processed before
	// foo. This is wrong as foo is in fact property of x.y namespace. So we must process all
	// module$exports$ symbols first. That's why we collected them in a separate list.
	natural := t.GetNaturalSymbolOrdering()
	slices.SortStableFunc(types, func(a, b *SymbolTable_Symbol) int { return natural(b, a) })
	allTypes := append(slices.Clip(types), exports...)

	// If you thought we are done with tricky case - you were wrong. There is another one!
	// The problem with the same property scope appearing several times. For example when using
	// aliases:
	//
	// const OBJ = {one: 1};
	// function() {
	//   const alias = OBJ;
	//   console.log(alias.one);
	// }
	//
	// In this case both 'OBJ' and 'alias' are considered property scopes and are candidates for
	// processing even though they share the same "type" which is "{one: 1}". As they share the same
	// type we need to process only one of them. To do that we build a "type => root symbol" map.
	// In this case the map will be {one: 1} => OBJ. Using this map will skip 'alias' when creating
	// property scopes.
	//
	// NOTE: we are using LinkedIdentityHashMap to compare types using == because we need to find
	// symbols
	// that point to the exact same type instance.
	for _, s := range allTypes {
		// Symbols are sorted in reverse order so that those with more outer scope will come later in
		// the list, and therefore override those set by aliases in more inner scope. The sorting
		// happens few lines above.
		type_ := t.getType(s)

		symbolForType := t.GetSymbolForTypeHelper(type_ /* linkToCtor= */, false)
		if (type_.IsNominalConstructorOrInterface() ||
			type_.IsFunctionPrototypeType() ||
			type_.IsEnumType()) &&
			s != symbolForType {
			// Some cases can't be handled by sorting. For example
			//
			// goog.module('some.Foo');
			// class Foo {}
			// exports = Foo;
			//
			// goog.module('some.bar');
			// const Foo = goog.require('some.Foo');
			//
			// Both Foo symbols (declaration and alias) have the same type and equal candidates to
			// become "symbol that declares Foo type". But in reality we should use the Foo in some.Foo
			// module. So if symbolForType exists and not equal to current symbol - it means it's an
			// alias.
			continue
		}
		t.symbolThatDeclaresType[type_] = s
	}

	for _, s := range allTypes {
		type_ := t.getType(s)

		// Create property scopes only based on "root" symbols for each type to handle aliases.
		if type_ == nil || s == t.symbolThatDeclaresType[type_] {
			t.createPropertyScopeFor(s)
		}
	}

	// Now we need to set the new property scope symbol to all aliases.
	for _, s := range allTypes {
		type_ := t.getType(s)

		if type_ != nil {
			if decl := t.symbolThatDeclaresType[type_]; decl != nil {
				s.propertyScope = decl.GetPropertyScope()
			}
		}
	}

	t.PruneOrphanedNames()
}

// Java: SymbolTable.needsPropertyScope(Symbol)
func (t *SymbolTable) needsPropertyScope(sym *SymbolTable_Symbol) bool {
	jsType := t.getType(sym)
	if jsType != nil && jsType.IsRecordType() {
		return true
	}
	type_ := rhino.ObjectType_Cast(jsType)
	if type_ == nil {
		return false
	}

	// Anonymous objects
	if type_.GetReferenceName() == nil {
		return true
	}

	// Constructors/prototypes
	if type_.IsNominalConstructorOrInterface() || type_.IsFunctionPrototypeType() {
		return true
	}
	if type_.IsEnumType() {
		return true
	}
	return false
}

// PruneOrphanedNames removes symbols where the namespace they're on has been
// removed.
//
// After filling property scopes, we may have two symbols represented in
// different ways. For example, "A.superClass_.foo" and B.prototype.foo".
//
// This resolves that ambiguity by pruning the duplicates. If we have a
// lexical symbol with a constructor in its property chain, then we assume
// there's also a property path to this symbol. In other words, we can remove
// "A.superClass_.foo" because it's rooted at "A", and we built a property
// scope for "A" above.
//
// Java: SymbolTable.pruneOrphanedNames()
func (t *SymbolTable) PruneOrphanedNames() {
nextSymbol:
	for _, s := range t.GetAllSymbols() {
		currentName := s.GetName()
		if s.IsProperty() || t.compiler.GetModuleMap().GetClosureModule(currentName) != nil {
			continue
		}

		// Go-only: byte indices of the ASCII '.' (only used to cut the string).
		dot := -1
		for {
			dot = strings.LastIndexByte(currentName, '.')
			if dot == -1 {
				break
			}
			currentName = currentName[:dot]

			owner := s.scope.GetQualifiedSlot(currentName)
			if owner != nil {
				if ownerType := t.getType(owner); ownerType != nil &&
					(ownerType.IsNominalConstructorOrInterface() ||
						ownerType.IsFunctionPrototypeType() ||
						ownerType.IsEnumType()) {
					t.removeSymbol(s)
					continue nextSymbol
				}
			}
		}
	}
}

// FillPropertySymbols creates symbols and references for all properties of
// types in this symbol table.
//
// This gets a little bit tricky, because of the way this symbol table
// conflates "type Foo" and "the constructor of type Foo". So if you have:
// SymbolTable symbolTable = for("var x = new Foo();");
// Symbol x = symbolTable.getGlobalScope().getSlot("x");
// Symbol type = symbolTable.getAllSymbolsForType(getType(x)).get(0);
// Then type.getPropertyScope() will have the properties of the constructor
// "Foo". To get the properties of instances of "Foo", you will need to call:
// Symbol instance = symbolTable.getSymbolForInstancesOf(type);
// As described at the top of this file, notice that "new Foo()" and
// "Foo.prototype" are represented by the same symbol.
//
// Java: SymbolTable.fillPropertySymbols(Node, Node)
func (t *SymbolTable) FillPropertySymbols(externs rhino.Node, root rhino.Node) {
	NewSymbolTable_PropertyRefCollector(t).Process(externs, root)
}

// FillJSDocInfo indexes JSDocInfo.
//
// Java: SymbolTable.fillJSDocInfo(Node, Node)
func (t *SymbolTable) FillJSDocInfo(externs rhino.Node, root rhino.Node) {
	NodeTraversal_TraverseRoots(
		t.compiler, NewSymbolTable_JSDocInfoCollector(t, t.compiler.GetTypeRegistry()), externs, root)

	// Create references to parameters in the JSDoc.
	for _, sym := range t.GetAllSymbols() {
		info := sym.GetJSDocInfo()
		if info == nil {
			continue
		}

		markers := info.GetMarkers()
		for i := range markers {
			marker := &markers[i]
			pos := marker.GetNameNode()
			if pos == nil {
				continue
			}

			paramNode := pos.GetItem()
			name := paramNode.GetString()
			param := t.GetParameterInFunction(sym, name)
			if param == nil {
				// There is no reference to this parameter in the actual JavaScript
				// code, so we'll try to create a special JsDoc-only symbol in
				// a JsDoc-only scope.
				typePos := marker.GetType()
				var type_ rhino.JSType
				if typePos != nil {
					type_ = typePos.GetItem().GetJSType()
				}

				if sym.docScope == nil {
					sym.docScope =
						NewSymbolTable_SymbolScope(
							rhino.NilNode /* root */, nil /* parent scope */, nil /* type of this */, sym)
				}

				// Check to make sure there's no existing symbol. In theory, this
				// should never happen, but we check anyway and fail silently
				// if our assumptions are wrong. (We do not want to put the symbol
				// table into an invalid state).
				existingSymbol := t.isAnySymbolDeclared(name, paramNode, sym.docScope)
				if existingSymbol == nil {
					t.declareSymbol(name, type_, type_ == nil, sym.docScope, paramNode /* info= */, nil)
				}
			} else {
				param.DefineReferenceAt(paramNode)
			}
		}
	}
}

// FillSymbolVisibility records the visibility of each symbol.
//
// Java: SymbolTable.fillSymbolVisibility(Node, Node)
func (t *SymbolTable) FillSymbolVisibility(externs rhino.Node, root rhino.Node) {
	collectPass := NewCollectFileOverviewVisibility(t.compiler)
	collectPass.Process(externs, root)
	if collectPass.hasDuplicate {
		panic(jutil.NewIllegalArgumentExceptionMsg("Multiple entries with same key"))
	}

	NodeTraversal_TraverseRoots(t.compiler, NewSymbolTable_VisibilityCollector(t, collectPass.visibilities), externs, root)
}

// createPropertyScopeFor builds a property scope for the given symbol. Any
// properties of the symbol will be added to the property scope.
//
// It is important that property scopes are created in order from the leaves
// up to the root, so this should only be called from #fillPropertyScopes. If
// you try to create a property scope for a parent before its leaf, then the
// leaf will get cut and re-added to the parent property scope, and weird
// things will happen.
//
// This function uses == to compare types to be exact same instances.
//
// Java: SymbolTable.createPropertyScopeFor(Symbol)
func (t *SymbolTable) createPropertyScopeFor(s *SymbolTable_Symbol) {
	// In order to build a property scope for s, we will need to build
	// a property scope for all its implicit prototypes first. This means
	// that sometimes we will already have built its property scope
	// for a previous symbol.
	if s.propertyScope != nil {
		return
	}

	jsType := t.getType(s)
	var type_ rhino.ObjectType
	if jsType != nil {
		type_ = jsType.ToObjectType()
	}
	if type_ == nil {
		return
	}

	// For cases like
	// const foo = goog.require('some.foo');
	// where createPropertyScopeFor() is called for 'foo' - instead of creating new PropertyScope we
	// reuse ModuleScope corresponding to 'some.foo' module.
	if s.GetName() == "exports" && s.GetDeclarationNode().IsModuleBody() {
		s.propertyScope = t.scopes.Get(s.GetDeclarationNode())
		return
	}
	// Create an empty property scope for the given symbol, maybe with a parent scope if it has
	// an implicit prototype.
	parentPropertyScope := t.maybeGetParentPropertyScope(type_)
	s.SetPropertyScope(NewSymbolTable_SymbolScope(rhino.NilNode, parentPropertyScope, type_, s))

	// If this symbol represents some 'a.b.c.prototype', add any instance properties of a.b.c
	// into the symbol scope.
	instanceType := type_
	propNames := type_.GetOwnPropertyNames()
	if instanceType.IsFunctionPrototypeType() {
		// Guard against modifying foo.prototype when foo is a regular (non-constructor) function.
		if instanceType.GetOwnerFunction().HasInstanceType() {
			// Merge the properties of "Foo.prototype" and "new Foo()" together.
			instanceType = instanceType.GetOwnerFunction().GetInstanceType()
			propNames = append(slices.Clip(propNames), instanceType.GetOwnPropertyNames()...)
		}
	}

	// Add all declared properties in propNames into the property scope
	for _, propName := range propNames {
		newProp := instanceType.GetSlot(propName)
		decl := newProp.GetDeclaration()
		if isNilStaticRef(decl) {
			// Skip properties without declarations. We won't know how to index
			// them, because we index things by node.
			continue
		}
		declNode := decl.GetNode()

		// We have symbol tables that do not do type analysis. They just try
		// to build a complete index of all objects in the program. So we might
		// already have symbols for things like "Foo.bar". If this happens,
		// throw out the old symbol and use the type-based symbol.
		oldProp :=
			t.symbols.Get(symbolTableKey{declNode: declNode, name: s.GetName() + "." + propName})

		if oldProp != nil && t.compiler.GetModuleMap().GetClosureModule(oldProp.GetName()) != nil {
			// This handles cases like:
			// goog.provide('a.b.c');
			// goog.provide('a.b.c.Foo');
			// Here were are creating scope for module 'a.b.c' and 'Foo' looks like a property of that
			// module. But it's a module by itself - so we don't move it to the scope of 'a.b.c' and
			// instead keeping 'a.b.c.Foo'  in global namespace.
			continue
		}

		// If we've already have an entry in the table for this symbol,
		// then skip it. This should only happen if we screwed up,
		// and declared multiple distinct properties with the same name
		// at the same node. We bail out here to be safe.
		if existing := t.symbols.Get(symbolTableKey{declNode: declNode, name: newProp.GetName()}); existing != nil {
			if symbolTable_logger.IsLoggable(jutil.Level_FINE) {
				symbolTable_logger.Fine("Found duplicate symbol " + jutil.ToJavaString(newProp))
			}
			continue
		}

		newSym := t.copySymbolTo(newProp, s.propertyScope)
		if oldProp != nil {
			if newSym.GetJSDocInfo() == nil {
				newSym.SetJSDocInfo(oldProp.GetJSDocInfo())
			}
			newSym.SetPropertyScope(oldProp.propertyScope)
			for ref := range oldProp.references.Values() {
				newSym.DefineReferenceAt(ref.GetNode())
			}
			// All references/scopes from oldProp were updated to use the newProp. Time to remove
			// oldProp.
			t.removeSymbol(oldProp)
		}
	}
}

// maybeGetParentPropertyScope returns, if this type has an implicit prototype
// set, the SymbolScope corresponding to the properties of the implicit
// prototype. Otherwise returns null.
//
// Note that currently we only handle cases where the implicit prototype is a)
// a class or b) is an instance object.
//
// Java: SymbolTable.maybeGetParentPropertyScope(ObjectType)
func (t *SymbolTable) maybeGetParentPropertyScope(symbolObjectType rhino.ObjectType) *SymbolTable_SymbolScope {
	proto := symbolObjectType.GetImplicitPrototype()
	if proto == nil || proto == symbolObjectType {
		return nil
	}
	var parentSymbol *SymbolTable_Symbol
	if t.isEs6ClassConstructor(proto) {
		// given `class Foo {} class Bar extends Foo {}`, `Foo` is the implicit prototype of `Bar`.
		parentSymbol = t.GetSymbolDeclaredBy(proto.ToMaybeFunctionType())
	} else if ctor := proto.GetConstructor(); ctor != nil {
		// given
		//   /** @constructor */ function Foo() {}
		//   /** @constructor */ function Bar() {}
		//   goog.inherits(Bar, Foo);
		// the implicit prototype of Bar.prototype is the instance of Foo.
		parentSymbol = t.GetSymbolForInstancesOfFunction(ctor)
	} else {
		return nil
	}
	if parentSymbol == nil {
		return nil
	}
	t.createPropertyScopeFor(parentSymbol)
	return parentSymbol.GetPropertyScope()
}

// Java: SymbolTable.isEs6ClassConstructor(JSType)
func (t *SymbolTable) isEs6ClassConstructor(type_ rhino.JSType) bool {
	fn := type_.ToMaybeFunctionType()
	return fn != nil && !fn.GetSource().IsNil() && fn.GetSource().IsClass()
}

// FillThisReferences fills in references to "this" variables.
//
// Java: SymbolTable.fillThisReferences(Node, Node)
func (t *SymbolTable) FillThisReferences(externs rhino.Node, root rhino.Node) {
	NewSymbolTable_ThisRefCollector(t).Process(externs, root)
}

// FillSuperReferences fills in references to "super" variables.
//
// Java: SymbolTable.fillSuperReferences(Node, Node)
func (t *SymbolTable) FillSuperReferences(externs rhino.Node, root rhino.Node) {
	collectSuper := &symbolTable_superCollector{outer: t}
	NodeTraversal_TraverseRoots(t.compiler, collectSuper, externs, root)
}

var (
	symbolTableSuperVisitTokens = func() ast.TokenSet {
		var s ast.TokenSet
		s.Add(rhino.Token_SUPER)
		return s
	}()
	symbolTableCallVisitTokens = func() ast.TokenSet {
		var s ast.TokenSet
		s.Add(rhino.Token_CALL)
		return s
	}()
)

// symbolTable_superCollector is the anonymous AbstractPostOrderCallback of
// fillSuperReferences.
//
// Go-only: language-adapter — an anonymous class.
type symbolTable_superCollector struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *SymbolTable
}

func (*symbolTable_superCollector) visitTokens() *ast.TokenSet {
	return &symbolTableSuperVisitTokens
}

// Java: SymbolTable.fillSuperReferences(...).new AbstractPostOrderCallback() {...}.visit(NodeTraversal, Node, Node)
func (c *symbolTable_superCollector) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	// Process only 'super' nodes with types.
	if !n.IsSuper() || n.GetJSType() == nil {
		return
	}
	classSymbol := c.outer.GetSymbolForTypeHelper(n.GetJSType() /* linkToCtor= */, false)
	if classSymbol != nil {
		classSymbol.DefineReferenceAt(n)
	}
}

// addRefsInGoogRequireStatement connects imported symbol back to its
// declaration. For example:
//
//	goog.module('some.Foo');
//	class Foo {}
//	exports = Foo;
//
//	goog.module('some.bar');
//	const Foo = goog.require('some.Foo');
//
// This method will add a reference from the `Foo` symbol in `some.Bar` module
// pointing to the `class Foo` in `some.Foo` module.
//
// Java: SymbolTable.addRefsInGoogRequireStatement(Node, Map<String, SymbolScope>)
func (t *SymbolTable) addRefsInGoogRequireStatement(n rhino.Node, googModuleScopes map[string]*SymbolTable_SymbolScope) {
	if t.compiler.GetOptions().ShouldRewriteModules() {
		// When modules are rewritten - there is no need to connect imported names back to definitions
		// as that was done with the help of Preprocessor tables in ClosureRewriteModule pass.
		return
	}
	require := ast.MustCall(n.GetLastChild())
	firstArg, _ := require.FirstArg()
	moduleName := firstArg.Node().GetString()
	module := t.compiler.GetModuleMap().GetClosureModule(moduleName)
	// Go-only: language-adapter — a null ModuleType (no module) makes the
	// Java switch statements below throw a NullPointerException.
	var moduleType ModuleMetadataMap_ModuleType
	hasModuleType := module != nil
	if hasModuleType {
		moduleType = module.Metadata().ModuleType()
	}
	moduleScope := googModuleScopes[moduleName]
	if n.IsName() {
		// This is `const Foo = goog.require('some.Foo');` case.
		// The type of Foo is either a class/interface in which case we can find its symbol that
		// declared class/interface.
		declaration := t.GetSymbolForTypeHelper(n.GetJSType(), false)
		if declaration != nil {
			declaration.DefineReferenceAt(n)
		}
		// This import might be for a constant like `const constant = goog.require('some.constant');`.
		// Try to find it in the corresponding definition in module.
		if !hasModuleType {
			panic(jutil.NewNullPointerException())
		}
		switch moduleType {
		case ModuleMetadataMap_ModuleType_GOOG_MODULE, ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE:
			// For goog.module search for 'exports' node and use that it as declaration. It would be
			// better to use right-handside of exports: `exports = constant;` but it's difficult to
			// get symbol for the right-handside. So do a more limited approach.
			if moduleScope == nil {
				panic(jutil.NewNullPointerException())
			}
			declaration = moduleScope.GetOwnSlot("exports")
		case ModuleMetadataMap_ModuleType_GOOG_PROVIDE:
			// For goog.provide `some.constant` should be defined in global namespace.
			declaration = t.getSymbolForName(rhino.NilNode, moduleName)
		default:
			// skip
		}
		if declaration != nil {
			declaration.DefineReferenceAt(n)
		}
	} else if dlhs, ok := ast.AsDestructuringLhs(n); ok {
		// This is `const {one} = goog.require('some.foo')` case.
		// For each imported symbol in destructuring we need to find corresponding declaration symbol
		// in `some.foo` namespace and connect them.
		objPat, ok := dlhs.Pattern().AsObjectPattern()
		jutil.CheckState(ok)
		// Iterate through all symbols in destructuring.
		for prop := range objPat.Props() {
			sk, ok := prop.AsStringKey()
			jutil.CheckState(ok)
			varName := sk.Key()
			var varDeclaration *SymbolTable_Symbol
			// AST of goog.module and goog.provide differs significantly so we need to lookup variables
			// differently.
			if !hasModuleType {
				panic(jutil.NewNullPointerException())
			}
			switch moduleType {
			case ModuleMetadataMap_ModuleType_GOOG_MODULE, ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE:
				if moduleScope == nil {
					panic(jutil.NewNullPointerException())
				}
				varDeclaration = moduleScope.GetOwnSlot(varName)
			case ModuleMetadataMap_ModuleType_GOOG_PROVIDE:
				varDeclaration = t.getSymbolForName(rhino.NilNode, moduleName+"."+varName)
			default:
				// skip
			}
			if varDeclaration != nil {
				varDeclaration.DefineReferenceAt(sk.Node())
			}
		}
	}
}

// FillGoogProvideModuleRequires connects goog.module/goog.provide calls with
// goog.require/goog.requireType. For each goog.module/goog.provide call
// creates a symbol and all goog.require/goog.requireType added as references.
//
// Java: SymbolTable.fillGoogProvideModuleRequires(Node, Node)
func (t *SymbolTable) FillGoogProvideModuleRequires(externs rhino.Node, root rhino.Node) {
	// small optimization to keep symbols created within this function for fast access instead of
	// going through SymbolTable.getSymbolForName() every time.
	declaredNamespaces := make(map[string]*SymbolTable_Symbol)
	// Map containing goog.module names to Scope object representing them. For goog.provide
	// namespaces there is no scope as elements of goog.provide are global.
	moduleScopes := make(map[string]*SymbolTable_SymbolScope)
	parseGoogCall := func(n rhino.Node) (method string, arg ast.StringLit, ok bool) {
		call, ok := ast.AsCall(n)
		if !ok {
			return "", ast.StringLit{}, false
		}
		gp, ok := call.Callee().AsGetProp()
		if !ok {
			return "", ast.StringLit{}, false
		}
		recv, ok := gp.Object().AsName()
		if !ok || recv.Ident() != "goog" {
			return "", ast.StringLit{}, false
		}
		firstArg, ok := call.FirstArg()
		if !ok {
			return "", ast.StringLit{}, false
		}
		argExpr, ok := firstArg.AsExpr()
		if !ok {
			return "", ast.StringLit{}, false
		}
		lit, ok := argExpr.AsStringLit()
		if !ok {
			return "", ast.StringLit{}, false
		}
		return gp.Property(), lit, true
	}
	collectModuleScopes :=
		&symbolTable_postOrderCallback{visit: func(nt *NodeTraversal, n rhino.Node, parent rhino.Node) {
			method, arg, ok := parseGoogCall(n)
			if !ok {
				return
			}
			modName := arg.Value()
			namespaceName := "ns$" + modName
			switch method {
			case "module", "provide":
				ns :=
					t.declareSymbol(
						namespaceName,
						/* type= */ nil,
						/* inferred= */ false,
						t.GetGlobalScope(),
						arg.Node(),
						/* info= */ nil)
				declaredNamespaces[namespaceName] = ns
				if gp := parent.GetParent(); gp.IsModuleBody() {
					moduleScopes[modName] = t.scopes.Get(gp)
				}
			default:
				// do nothing. Some other goog.xyz call.
			}
		}}
	processRequireStatements :=
		&symbolTable_postOrderCallback{visit: func(nt *NodeTraversal, n rhino.Node, parent rhino.Node) {
			method, arg, ok := parseGoogCall(n)
			if !ok {
				return
			}
			namespaceName := "ns$" + arg.Value()
			switch method {
			case "require", "requireType", "forwardDeclare":
				t.addRefsInGoogRequireStatement(parent, moduleScopes)
				symbol := declaredNamespaces[namespaceName]
				// We expect that namespace was already processed by that point, but in some broken
				// code it might be missing. In which case just skip it.
				if symbol != nil {
					symbol.DefineReferenceAt(arg.Node())
				}
			default:
				// do nothing. Some other goog.xyz call.
			}
		}}
	NodeTraversal_TraverseRoots(t.compiler, collectModuleScopes, externs, root)
	NodeTraversal_TraverseRoots(t.compiler, processRequireStatements, externs, root)
}

// symbolTable_postOrderCallback is an anonymous AbstractPostOrderCallback of
// fillGoogProvideModuleRequires.
//
// Go-only: language-adapter — an anonymous class.
type symbolTable_postOrderCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	visit func(t *NodeTraversal, n rhino.Node, parent rhino.Node)
}

func (*symbolTable_postOrderCallback) noVisitLeaves() {}
func (*symbolTable_postOrderCallback) visitTokens() *ast.TokenSet {
	return &symbolTableCallVisitTokens
}

// Java: SymbolTable.fillGoogProvideModuleRequires(...).new AbstractPostOrderCallback() {...}.visit(NodeTraversal, Node, Node)
func (c *symbolTable_postOrderCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	c.visit(t, n, parent)
}

// FlattenGoogModuleExports flattens symbols within a module scope
// (goog.module). In goog.module we often have different variations of
// declaring Foo, exports.Foo, exports = {Foo}. In some cases there are both
// Foo and exports.Foo = Foo while in others there is only exports.Foo.
// Semantically they are the same to the user so we flatten them into Foo. All
// 'exports.Foo' symbols renamed or merged into 'Foo' (depending on whether it
// already exists). That way we get more predictable table which is required
// for better index data and integration with other languages, e.g. JS <=> TS
// indexing.
//
// This is not a correct operation in 100% of times. For example it will
// produce very incorrect result for:
//
//	class Foo {}
//	exports.Foo = NotFoo;
//
// But that should be very rare case and the code is misleading and should be
// avoided.
//
// Java: SymbolTable.flattenGoogModuleExports()
func (t *SymbolTable) FlattenGoogModuleExports() {
	for moduleScope := range t.GetAllScopes() {
		if !moduleScope.IsModuleScope() {
			continue
		}
		for _, sym := range slices.Collect(moduleScope.ownSymbols.Values()) {
			if strings.HasPrefix(sym.GetName(), "exports.") {
				t.renameSymbol(moduleScope, sym, sym.GetName()[len("exports."):])
			}
		}
	}
}

// isSymbolAQuotedObjectKey checks whether symbol is a quoted object literal
// key. In the following object:
//
// var foo = {'one': 1, two: 2};
//
// 'one' is quoted key. while two is not.
//
// Java: SymbolTable.isSymbolAQuotedObjectKey(Symbol)
func (t *SymbolTable) isSymbolAQuotedObjectKey(symbol *SymbolTable_Symbol) bool {
	node := symbol.GetDeclarationNode()
	return !node.IsNil() && node.IsStringKey() && node.IsQuotedStringKey()
}

// isSymbolDuplicatedExternOnWindow is a heuristic method to check whether
// symbol was created by DeclaredGlobalExternsOnWindow.java pass.
//
// Java: SymbolTable.isSymbolDuplicatedExternOnWindow(Symbol)
func (t *SymbolTable) isSymbolDuplicatedExternOnWindow(symbol *SymbolTable_Symbol) bool {
	node := symbol.GetDeclarationNode()
	// Check that node is of type "window.foo";
	if node.IsIndexable() {
		return false
	}
	gp, ok := ast.AsGetProp(node)
	if !ok {
		return false
	}
	recv, ok := gp.Object().AsName()
	return ok && recv.Ident() == "window"
}

// mergeExternSymbolsDuplicatedOnWindow merges the duplicates
// DeclaredGLobalExternsOnWindow.java creates. That pass duplicates all global
// variables so that:
//
//	var foo;
//
// becomes
//
//	var foo;
//	window.foo;
//
// This function finds all such cases and merges window.foo symbol back to
// foo. It changes window.foo references to point to foo symbol.
//
// Java: SymbolTable.mergeExternSymbolsDuplicatedOnWindow()
func (t *SymbolTable) mergeExternSymbolsDuplicatedOnWindow() {
	// To find duplicated symbols we rely on the fact that duplicated symbol share the same
	// source position as original symbol and going to use filename => sourcePosition => symbol
	// table.
	//
	externSymbols := make(map[externSymbolPosKey]*SymbolTable_Symbol)
	for _, symbol := range t.symbols.ValueSlice() {
		declNode := symbol.GetDeclarationNode()
		if declNode.IsNil() {
			continue
		}
		sf := declNode.GetStaticSourceFile()
		if isNilSourceFile(sf) || !sf.IsExtern() {
			continue
		}
		sourceFilePtr := declNode.GetSourceFileName()
		if sourceFilePtr == nil {
			// Java: Table.put(null, ...) throws NullPointerException.
			panic(jutil.NewNullPointerException())
		}
		key := externSymbolPosKey{
			sourceFile: *sourceFilePtr,
			position:   declNode.GetSourcePosition(),
		}
		existingSymbol, ok := externSymbols[key]
		if !ok {
			externSymbols[key] = symbol
			continue
		}
		// Consider 2 possibilies: either symbol or existingSymbol might be the generated symbol we
		// are looking for.
		if t.isSymbolDuplicatedExternOnWindow(existingSymbol) {
			t.mergeSymbol(existingSymbol, symbol)
			externSymbols[key] = symbol
		} else if t.isSymbolDuplicatedExternOnWindow(symbol) {
			t.mergeSymbol(symbol, existingSymbol)
		}
	}
}

// RemoveGeneratedSymbols removes various generated symbols that are invisible
// to users and pollute or mess up index. Jscompiler does transpilations that
// might introduce extra nodes/symbols. Most of these symbols should not get
// into final SymbolTable because SymbolTable should contain only symbols that
// correspond to a symbol in original source code (before transpilation).
//
// Java: SymbolTable.removeGeneratedSymbols()
func (t *SymbolTable) RemoveGeneratedSymbols() {
	// Need to iterate over copy of values list because removeSymbol() will change the map
	// and we'll get ConcurrentModificationException
	for _, symbol := range t.symbols.ValueSlice() {
		if t.isSymbolAQuotedObjectKey(symbol) {
			// Quoted object keys are not considered symbols. Only unquoted keys and dot-access
			// properties are considered symbols. Remove the quoted key.
			symbolAlreadyRemoved := !t.GetScope(symbol).ownSymbols.ContainsKey(symbol.GetName())
			if !symbolAlreadyRemoved {
				t.removeSymbol(symbol)
			}
		} else if !symbol.GetDeclarationNode().IsNil() &&
			symbol.GetDeclarationNode().IsModuleBody() {
			// Symbols that represent whole module (goog.module) aren't needed for indexing. We already
			// recorded imports/references of that module in fillGoogProvideModuleRequires() function.
			// Individual exported functions/variables also have their own symbols.
			t.removeSymbol(symbol)
		}
	}
	t.mergeExternSymbolsDuplicatedOnWindow()
}
