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
	"iter"
	"slices"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// SymbolTable is a symbol table for people that want to use Closure Compiler
// as an indexer.
//
// Contains an index of all the symbols in the code within a compilation job.
// The API is designed for people who want to visit all the symbols, rather
// than people who want to lookup a specific symbol by a certain key.
//
// We can use this to combine different types of symbol tables. For example,
// one class might have a StaticSymbolTable of all variable references, and
// another class might have a StaticSymbolTable of all type names in JSDoc
// comments. This class allows you to combine them into a unified index.
//
// Most passes build their own "partial" symbol table that implements the same
// interface (StaticSymbolTable, StaticSlot, and friends). Individual compiler
// passes usually need more or less metadata about the certainty of symbol
// information. Building a complete symbol table with all the necessary
// metadata for all passes would be too slow. However, as long as these
// "partial" symbol tables implement the proper interfaces, we should be able
// to add them to this symbol table to make it more complete.
//
// If clients want fast lookup, they should build their own wrapper around this
// symbol table that indexes symbols or references by the desired lookup key.
//
// By design, when this symbol table creates symbols for types, it tries to
// mimic the symbol table you would get in an OO language. For example, the
// "type Foo" and "the constructor that creates objects of type Foo" are the
// same symbol. The types of "Foo.prototype" and "new Foo()" also have the same
// symbol. Although JSCompiler internally treats these as distinct symbols, we
// assume that most clients will not care about the distinction.
//
// See SymbolTable_AddSymbolsFrom for more information on how to write plugins
// for this symbol table.
//
// Go-only: language-adapter — Guava's Ordering fields (sourceNameOrdering,
// nodeOrdering, lexicalScopeOrdering, symbolOrdering) are comparator
// functions; getNaturalSymbolOrdering returns one. The generic methods
// addScopes and addSymbolsFrom are the package-level functions
// SymbolTable_AddScopes and SymbolTable_AddSymbolsFrom.
type symbolTableKey struct {
	declNode rhino.Node
	name     string
}

// Java: public final class SymbolTable
type SymbolTable struct {
	// All symbols in the program, uniquely identified by the node where
	// they're declared and their name.
	symbols symbolOrderedMap[symbolTableKey, *SymbolTable_Symbol]

	// All syntactic scopes in the program, uniquely identified by the node
	// where they're declared.
	scopes symbolOrderedMap[rhino.Node, *SymbolTable_SymbolScope]

	// All Nodes with JSDocInfo in the program.
	docInfos []rhino.Node

	globalScope *SymbolTable_SymbolScope

	compiler AbstractCompiler

	registry *rhino.JSTypeRegistry

	// Map from a JSType to the symbol that declares it.
	symbolThatDeclaresType map[rhino.JSType]*SymbolTable_Symbol

	// Comparators
	sourceNameOrdering   func(a, b *string) int
	nodeOrdering         func(a, b rhino.Node) int
	lexicalScopeOrdering func(a, b *SymbolTable_SymbolScope) int
	symbolOrdering       func(a, b *SymbolTable_Symbol) int
}

// SymbolTable_GLOBAL_THIS is the name we use for the JavaScript built-in
// Global object. It's anonymous in JavaScript, so we have to give it an
// invalid identifier to avoid conflicts with user-defined property names.
//
// Java: SymbolTable.GLOBAL_THIS
const SymbolTable_GLOBAL_THIS = "*global*"

// NewSymbolTable creates a symbol table. Clients should get a symbol table by
// asking the compiler at the end of a compilation job.
//
// Java: SymbolTable(AbstractCompiler, JSTypeRegistry)
func NewSymbolTable(compiler AbstractCompiler, registry *rhino.JSTypeRegistry) *SymbolTable {
	t := &SymbolTable{
		symbolThatDeclaresType: make(map[rhino.JSType]*SymbolTable_Symbol),
		compiler:               compiler,
		registry:               registry,
	}
	t.initOrderings()
	return t
}

// Java: SymbolTable.getReferences(Symbol)
func (t *SymbolTable) GetReferences(symbol *SymbolTable_Symbol) iter.Seq[*SymbolTable_Reference] {
	return symbol.references.Values()
}

// Java: SymbolTable.getReferenceList(Symbol)
func (t *SymbolTable) GetReferenceList(symbol *SymbolTable_Symbol) []*SymbolTable_Reference {
	return slices.Collect(symbol.references.Values())
}

// Java: SymbolTable.getAllSymbols()
func (t *SymbolTable) GetAllSymbols() []*SymbolTable_Symbol {
	return t.symbols.ValueSlice()
}

// GetAllSymbolsSorted gets the symbols in their natural ordering. Always
// returns a mutable list.
//
// Java: SymbolTable.getAllSymbolsSorted()
func (t *SymbolTable) GetAllSymbolsSorted() []*SymbolTable_Symbol {
	result := t.symbols.ValueSlice()
	slices.SortStableFunc(result, t.GetNaturalSymbolOrdering())
	return result
}

// GetNaturalSymbolOrdering gets the 'natural' ordering of symbols.
//
// Right now, we only guarantee that symbols in the global scope will come
// before symbols in local scopes. After that, the order is deterministic but
// undefined.
//
// Java: SymbolTable.getNaturalSymbolOrdering()
func (t *SymbolTable) GetNaturalSymbolOrdering() func(a, b *SymbolTable_Symbol) int {
	return t.symbolOrdering
}

// Java: SymbolTable.getScope(Symbol)
func (t *SymbolTable) GetScope(slot *SymbolTable_Symbol) *SymbolTable_SymbolScope {
	return slot.scope
}

// Java: SymbolTable.getAllJSDocInfoNodes()
func (t *SymbolTable) GetAllJSDocInfoNodes() []rhino.Node {
	return t.docInfos
}

// GetEnclosingScope gets the scope that contains the given node. If n is a
// function name, we return the scope that contains the function, not the
// function itself.
//
// Java: SymbolTable.getEnclosingScope(Node)
func (t *SymbolTable) GetEnclosingScope(n rhino.Node) *SymbolTable_SymbolScope {
	current := n.GetParent()
	if n.IsName() && n.GetParent().IsFunction() {
		current = current.GetParent()
	}

	for ; !current.IsNil(); current = current.GetParent() {
		if scope := t.scopes.Get(current); scope != nil {
			return scope
		}
	}
	return nil
}

// GetEnclosingFunctionScope gets the scope that contains the given node. If n
// is a function name, we return the scope that contains the function, not the
// function itself. The returned scope is either function or global scope.
//
// Java: SymbolTable.getEnclosingFunctionScope(Node)
func (t *SymbolTable) GetEnclosingFunctionScope(n rhino.Node) *SymbolTable_SymbolScope {
	current := n.GetParent()
	if n.IsName() && !current.IsNil() && current.IsFunction() {
		current = current.GetParent()
	}

	for ; !current.IsNil(); current = current.GetParent() {
		scope := t.scopes.Get(current)
		if scope != nil && !scope.IsBlockScope() {
			return scope
		}
	}
	return t.globalScope
}

// GetParameterInFunction tries to find, if sym is a function, a Symbol for a
// parameter with the given name.
//
// Returns null if we couldn't find one.
//
// Notice that this just makes a best effort, and may not be able to find
// parameters for non-conventional function definitions. For example, we would
// not be able to find "y" in this code:
// var x = x() ? function(y) {} : function(y) {};
//
// Java: SymbolTable.getParameterInFunction(Symbol, String)
func (t *SymbolTable) GetParameterInFunction(sym *SymbolTable_Symbol, paramName string) *SymbolTable_Symbol {
	scope := t.GetScopeInFunction(sym)
	if scope != nil {
		param := scope.GetSlot(paramName)
		if param != nil && param.scope == scope {
			return param
		}
	}
	return nil
}

// Java: SymbolTable.getScopeInFunction(Symbol)
func (t *SymbolTable) GetScopeInFunction(sym *SymbolTable_Symbol) *SymbolTable_SymbolScope {
	type_ := sym.GetFunctionType()
	if type_ == nil {
		return nil
	}

	functionNode := type_.GetSource()
	if functionNode.IsNil() {
		return nil
	}

	return t.scopes.Get(functionNode)
}

// GetSymbolForScope returns the symbol associated with the given scope. All
// local scopes are associated with a function, and some functions are
// associated with a symbol.
//
// Java: SymbolTable.getSymbolForScope(SymbolScope)
func (t *SymbolTable) GetSymbolForScope(scope *SymbolTable_SymbolScope) *SymbolTable_Symbol {
	if scope.GetSymbolForScope() == nil {
		scope.SetSymbolForScope(t.FindSymbolForScope(scope))
	}
	return scope.GetSymbolForScope()
}

// FindSymbolForScope finds the symbol associated with the given scope. Notice
// that we won't always be able to figure out this association dynamically, so
// sometimes we'll just create the association when we create the scope.
//
// Java: SymbolTable.findSymbolForScope(SymbolScope)
func (t *SymbolTable) FindSymbolForScope(scope *SymbolTable_SymbolScope) *SymbolTable_Symbol {
	rootNode := scope.GetRootNode()
	if rootNode.GetParent().IsNil() {
		return t.globalScope.GetSlot(SymbolTable_GLOBAL_THIS)
	}

	if !rootNode.IsFunction() {
		return nil
	}

	name := NodeUtil_GetBestLValueName(NodeUtil_GetBestLValue(rootNode))
	if name == nil {
		return nil
	}
	return scope.GetParentScope().GetQualifiedSlot(*name)
}

// GetAllSymbolsForTypeOf gets all symbols associated with the type of the
// given symbol.
//
// For example, given a variable x declared as /* @type {Array|Date} / var x =
// f(); this will return the constructors for Array and Date.
//
// Java: SymbolTable.getAllSymbolsForTypeOf(Symbol)
func (t *SymbolTable) GetAllSymbolsForTypeOf(sym *SymbolTable_Symbol) []*SymbolTable_Symbol {
	return t.GetAllSymbolsForType(t.getType(sym))
}

// GetGlobalScope returns the global scope.
//
// Java: SymbolTable.getGlobalScope()
func (t *SymbolTable) GetGlobalScope() *SymbolTable_SymbolScope {
	return t.globalScope
}

// GetSymbolDeclaredBy gets the symbol for the given constructor or interface.
//
// Java: SymbolTable.getSymbolDeclaredBy(FunctionType)
func (t *SymbolTable) GetSymbolDeclaredBy(fn rhino.FunctionType) *SymbolTable_Symbol {
	jutil.CheckState(fn.IsConstructor() || fn.IsInterface())
	instanceType := fn.GetInstanceType()
	return t.GetSymbolForName(fn.GetSource(), instanceType.GetReferenceName())
}

// GetSymbolDeclaredByEnum gets the symbol for the given enum.
//
// Java: SymbolTable.getSymbolDeclaredBy(EnumType)
func (t *SymbolTable) GetSymbolDeclaredByEnum(enumType *rhino.EnumType) *SymbolTable_Symbol {
	return t.GetSymbolForName(enumType.GetSource(), enumType.GetElementsType().GetReferenceName())
}

// GetSymbolForInstancesOf gets the symbol for the prototype if this is the
// symbol for a constructor or interface.
//
// Java: SymbolTable.getSymbolForInstancesOf(Symbol)
func (t *SymbolTable) GetSymbolForInstancesOf(sym *SymbolTable_Symbol) *SymbolTable_Symbol {
	fn := sym.GetFunctionType()
	if fn != nil && fn.IsNominalConstructorOrInterface() {
		return t.GetSymbolForInstancesOfFunction(fn)
	}
	return nil
}

// GetSymbolForInstancesOfFunction gets the symbol for the prototype of the
// given constructor or interface.
//
// Java: SymbolTable.getSymbolForInstancesOf(FunctionType)
func (t *SymbolTable) GetSymbolForInstancesOfFunction(fn rhino.FunctionType) *SymbolTable_Symbol {
	jutil.CheckState(fn.IsConstructor() || fn.IsInterface())
	pType := fn.GetPrototype()
	return t.GetSymbolForName(fn.GetSource(), pType.GetReferenceName())
}

// Java: SymbolTable.getSymbolForName(Node, String)
func (t *SymbolTable) GetSymbolForName(source rhino.Node, name *string) *SymbolTable_Symbol {
	if name == nil {
		return nil
	}
	return t.getSymbolForName(source, *name)
}

func (t *SymbolTable) getSymbolForName(source rhino.Node, nameStr string) *SymbolTable_Symbol {
	if t.globalScope == nil {
		return nil
	}

	var scope *SymbolTable_SymbolScope
	if source.IsNil() {
		scope = t.globalScope
	} else {
		scope = t.GetEnclosingScope(source)
	}
	if scope == nil {
		return nil
	}
	if scope.IsModuleScope() {
		// In the case of elements like
		//
		// goog.module('foo.bar');
		// exports.Pinky = {};
		//
		// we might get fully-qualified name `foo.bar.Pinky` but when looking it up within the module
		// - the symbol is just `Pinky`. So we need to remove module name to be able to find the
		// symbol in module scope.
		var metadata *ModuleMetadataMap_ModuleMetadata
		if sourceFileName := scope.GetRootNode().GetSourceFileName(); sourceFileName != nil {
			metadata = t.compiler.GetModuleMetadataMap().ModuleByPath(*sourceFileName)
		}
		if metadata == nil {
			panic(jutil.NewNullPointerException())
		}
		first := "null"
		if entries := metadata.GoogNamespaces().Entries(); len(entries) > 0 {
			first = entries[0].Key
		}
		moduleNameWithDot := first + "."
		if strings.HasPrefix(nameStr, moduleNameWithDot) {
			nameStr = nameStr[len(moduleNameWithDot):]
		}
	}

	return scope.GetQualifiedSlot(nameStr)
}

// GetAllSymbolsForType gets all symbols associated with the given type. For
// union types, this may be multiple symbols. For instance types, this will
// return the constructor of that instance.
//
// Java: SymbolTable.getAllSymbolsForType(JSType)
func (t *SymbolTable) GetAllSymbolsForType(type_ rhino.JSType) []*SymbolTable_Symbol {
	if type_ == nil {
		return []*SymbolTable_Symbol{}
	}

	if type_.IsUnionType() {
		result := []*SymbolTable_Symbol{}
		for _, alt := range type_.ToMaybeUnionType().GetAlternates() {
			// Our type system never has nested unions.
			altSym := t.GetSymbolForTypeHelper(alt, true)
			if altSym != nil {
				result = append(result, altSym)
			}
		}
		return result
	}
	result := t.GetSymbolForTypeHelper(type_, true)
	if result == nil {
		return []*SymbolTable_Symbol{}
	}
	return []*SymbolTable_Symbol{result}
}

// GetSymbolForTypeHelper gets all symbols associated with the given type. If
// there is more that one symbol associated with the given type, return null.
//
// type_ is the type. If linkToCtor is true, we should link instance types back
// to their constructor function. If false, we should link instance types back
// to their prototype. See the comments at the top of this file for more
// information on how our internal type system is more granular than Symbols.
//
// Java: SymbolTable.getSymbolForTypeHelper(JSType, boolean)
func (t *SymbolTable) GetSymbolForTypeHelper(type_ rhino.JSType, linkToCtor bool) *SymbolTable_Symbol {
	if type_ == nil {
		return nil
	}

	if type_.IsGlobalThisType() {
		return t.globalScope.GetSlot(SymbolTable_GLOBAL_THIS)
	} else if type_.IsNominalConstructorOrInterface() {
		if linkToCtor {
			return t.globalScope.GetSlot("Function")
		}
		return t.GetSymbolDeclaredBy(type_.ToMaybeFunctionType())
	} else if type_.IsFunctionPrototypeType() {
		ownerFn := type_.(rhino.ObjectType).GetOwnerFunction()
		if !ownerFn.IsConstructor() && !ownerFn.IsInterface() {
			return nil
		}
		if linkToCtor {
			return t.GetSymbolDeclaredBy(ownerFn)
		}
		return t.GetSymbolForInstancesOfFunction(ownerFn)
	} else if type_.IsInstanceType() {
		ownerFn := type_.(rhino.ObjectType).GetConstructor()
		if linkToCtor {
			return t.GetSymbolDeclaredBy(ownerFn)
		}
		return t.GetSymbolForInstancesOfFunction(ownerFn)
	} else if type_.IsFunctionType() {
		if linkToCtor {
			return t.globalScope.GetSlot("Function")
		}
		return t.globalScope.GetQualifiedSlot("Function.prototype")
	} else if autobox := type_.AutoboxesTo(); autobox != nil {
		return t.GetSymbolForTypeHelper(autobox, linkToCtor)
	} else if type_.IsEnumType() {
		return t.GetSymbolDeclaredByEnum(type_.(*rhino.EnumType))
	} else if type_.IsRecordType() {
		return t.symbolThatDeclaresType[type_]
	}
	return nil
}

// ToDebugStringTree returns a debug representation of the SymbolTable. It
// starts with the global scope and all its symbols recursively going to
// children scopes and printing their symbols.
//
// Java: SymbolTable.toDebugStringTree()
func (t *SymbolTable) ToDebugStringTree() string {
	// We are recursively going through scopes from global to nested and need to build parent =>
	// children map as it doesn't exist at the moment.
	childrenScopes := make(map[*SymbolTable_SymbolScope][]*SymbolTable_SymbolScope)
	for scope := range t.GetAllScopes() {
		if parent := scope.GetParentScope(); parent != nil {
			childrenScopes[parent] = append(childrenScopes[parent], scope)
		}
	}
	var builder strings.Builder
	// Go through all scopes that don't have parent. It's likely only Global scope.
	for scope := range t.GetAllScopes() {
		if scope.GetParentScope() == nil {
			t.toDebugStringTree(&builder, "", scope, childrenScopes)
		}
	}

	// All objects/classes have property scopes. They usually don't have parents. These scopes
	// confusingly are not included in "getAllScopes()" pass above.
	visitedScopes := make(map[*SymbolTable_SymbolScope]struct{})
	for scope := range t.GetAllScopes() {
		visitedScopes[scope] = struct{}{}
	}
	for _, symbol := range t.GetAllSymbols() {
		if symbol.propertyScope == nil {
			continue
		}
		if _, visited := visitedScopes[symbol.propertyScope]; visited {
			continue
		}
		visitedScopes[symbol.GetPropertyScope()] = struct{}{}
		t.toDebugStringTree(&builder, "", symbol.propertyScope, childrenScopes)
	}
	return builder.String()
}

// Java: SymbolTable.toDebugStringTree(StringBuilder, String, SymbolScope, ImmutableListMultimap<SymbolScope, SymbolScope>)
func (t *SymbolTable) toDebugStringTree(
	builder *strings.Builder,
	prefix string,
	scope *SymbolTable_SymbolScope,
	childrenScopes map[*SymbolTable_SymbolScope][]*SymbolTable_SymbolScope) {
	builder.WriteString(prefix)
	builder.WriteString(scope.String())
	childrenPrefix := prefix + "    "
	printedSomething := false
	if !scope.ownSymbols.IsEmpty() {
		printedSomething = true
		builder.WriteString("\n" + prefix + "  Symbols:\n")
		for symbol := range scope.ownSymbols.Values() {
			t.toDebugString(builder, childrenPrefix, symbol /* printScope= */, false)
		}
	}
	if children := childrenScopes[scope]; len(children) > 0 {
		printedSomething = true
		builder.WriteString("\n" + prefix + "  Scopes:\n")
		for _, childScope := range children {
			t.toDebugStringTree(builder, childrenPrefix, childScope, childrenScopes)
		}
	}
	if !printedSomething {
		builder.WriteByte('\n')
	}
	builder.WriteByte('\n')
}

// Java: SymbolTable.toDebugString()
func (t *SymbolTable) ToDebugString() string {
	var builder strings.Builder
	for _, symbol := range t.GetAllSymbols() {
		t.toDebugString(&builder, "", symbol /* printScope= */, true)
	}
	return builder.String()
}

// Java: SymbolTable.toDebugString(StringBuilder, String, Symbol, boolean)
func (t *SymbolTable) toDebugString(
	builder *strings.Builder, prefix string, symbol *SymbolTable_Symbol, printScope bool) {
	builder.WriteString(prefix + "'" + symbol.GetName() + "'")
	if printScope {
		builder.WriteString(" in " + symbol.scope.String())
	}
	builder.WriteByte('\n')

	refCount := 0
	for ref := range t.GetReferences(symbol) {
		node := ref.GetNode()
		indexable := "non indexable"
		if node.IsIndexable() {
			indexable = ""
		}
		builder.WriteString(prefix)
		builder.WriteString("  Ref ")
		builder.WriteString(strconv.Itoa(refCount))
		builder.WriteString(": ")
		builder.WriteString(symbolTable_nullableString(node.GetSourceFileName()))
		builder.WriteString(" line: ")
		builder.WriteString(strconv.Itoa(node.GetLineno()))
		builder.WriteString(" col: ")
		builder.WriteString(strconv.Itoa(node.GetCharno()))
		builder.WriteString(" len: ")
		builder.WriteString(strconv.Itoa(node.GetLength()))
		builder.WriteByte(' ')
		builder.WriteString(indexable)
		builder.WriteByte('\n')
		refCount++
	}
}

// SymbolTable_AddScopes makes sure all the given scopes in otherSymbolTable
// are in this symbol table.
//
// Java: SymbolTable.addScopes(Collection<S>)
func SymbolTable_AddScopes[S rhino.StaticScope](t *SymbolTable, scopes []S) {
	for _, scope := range scopes {
		t.CreateScopeFrom(scope)
	}
}

// FindScopes finds all the scopes and adds them to this symbol table.
//
// Java: SymbolTable.findScopes(Node, Node)
func (t *SymbolTable) FindScopes(externs rhino.Node, root rhino.Node) {
	NodeTraversal_TraverseRoots(
		t.compiler,
		&symbolTable_scopeFinder{outer: t},
		externs,
		root)
}

// symbolTable_scopeFinder is the anonymous AbstractScopedCallback of
// findScopes.
//
// Go-only: language-adapter — an anonymous class.
type symbolTable_scopeFinder struct {
	NodeTraversal_AbstractScopedCallbackBase
	outer *SymbolTable
}

// Java: SymbolTable.findScopes(...).new AbstractScopedCallback() {...}.enterScope(NodeTraversal)
func (f *symbolTable_scopeFinder) EnterScope(t *NodeTraversal) {
	f.outer.CreateScopeFrom(t.GetScope())
}

// Java: SymbolTable.findScopes(...).new AbstractScopedCallback() {...}.visit(NodeTraversal, Node, Node)
func (f *symbolTable_scopeFinder) Visit(t *NodeTraversal, n rhino.Node, p rhino.Node) {}

func (*symbolTable_scopeFinder) nodeTraversalNoVisit() {}

// GetAllScopes gets all the scopes in this symbol table.
//
// Java: SymbolTable.getAllScopes()
func (t *SymbolTable) GetAllScopes() iter.Seq[*SymbolTable_SymbolScope] {
	return t.scopes.Values()
}

// CreateScopeFrom returns, given a scope from another symbol table, the
// SymbolScope rooted at the same node. Creates one if it doesn't exist yet.
//
// Java: SymbolTable.createScopeFrom(StaticScope)
func (t *SymbolTable) CreateScopeFrom(otherScope rhino.StaticScope) *SymbolTable_SymbolScope {
	otherScopeRoot := otherScope.GetRootNode()

	myScope := t.scopes.Get(otherScopeRoot)
	if myScope == nil {
		otherScopeParent := otherScope.GetParentScope()

		// If otherScope is a global scope, and we already have a global scope,
		// then something has gone seriously wrong.
		//
		// Not all symbol tables are rooted at the same global node, and
		// we do not want to mix and match symbol tables that are rooted
		// differently.

		if isNilStaticScope(otherScopeParent) {
			// The global scope must be created before any local scopes.
			jutil.CheckStateMsg(t.globalScope == nil, "Global scopes found at different roots")
		}

		var parent *SymbolTable_SymbolScope
		if !isNilStaticScope(otherScopeParent) {
			parent = t.CreateScopeFrom(otherScopeParent)
		}
		myScope =
			NewSymbolTable_SymbolScope(
				otherScopeRoot,
				parent,
				t.getTypeOfThis(otherScope),
				nil)
		t.scopes.Put(otherScopeRoot, myScope)
		if myScope.IsGlobalScope() {
			t.globalScope = myScope
		}
	}
	return myScope
}

// AddAnonymousFunctions finds anonymous functions in local scopes, and gives
// them names and symbols. They will show up as local variables with names
// "function%0", "function%1", etc.
//
// Java: SymbolTable.addAnonymousFunctions()
func (t *SymbolTable) AddAnonymousFunctions() {
	var scopes []*SymbolTable_SymbolScope
	for scope := range t.GetAllScopes() {
		if scope.IsLexicalScope() {
			scopes = append(scopes, scope)
		}
	}
	slices.SortFunc(scopes, t.lexicalScopeOrdering)
	scopes = slices.CompactFunc(scopes, func(a, b *SymbolTable_SymbolScope) bool {
		return t.lexicalScopeOrdering(a, b) == 0
	})

	for _, scope := range scopes {
		t.addAnonymousFunctionsInScope(scope)
	}
}

// Java: SymbolTable.addAnonymousFunctionsInScope(SymbolScope)
func (t *SymbolTable) addAnonymousFunctionsInScope(scope *SymbolTable_SymbolScope) {
	sym := t.GetSymbolForScope(scope)
	if sym == nil {
		// JSCompiler has no symbol for this scope. Check to see if it's a
		// local function. If it is, give it a name.
		rootNode := scope.GetRootNode()
		if scope.IsLexicalScope() &&
			!scope.IsGlobalScope() &&
			!rootNode.IsNil() &&
			!rootNode.IsFromExterns() &&
			scope.GetParentScope() != nil &&
			rootNode.IsFunction() {
			parent := scope.GetParentScope()

			innerName := "function%" + strconv.Itoa(scope.GetIndexInParent())
			type_ := rootNode.GetJSType()

			// Functions defined on anonymous objects are considered anonymous as well:
			// doFoo({bar() {}});
			// bar is not technically anonymous, but it's a method on an anonymous object literal so
			// effectively it's anonymous/inaccessible. In this case, slightly correct rootNode to
			// be a MEMBER_FUNCTION_DEF node instead of a FUNCTION node.
			if rootNode.GetParent().IsMemberFunctionDef() {
				rootNode = rootNode.GetParent()
			}

			anonymousFunctionSymbol :=
				t.declareSymbol(
					innerName, type_ /* inferred= */, true, parent, rootNode /* info= */, nil)
			scope.SetSymbolForScope(anonymousFunctionSymbol)
		}
	}
}

// symbolTable_nullableString is String.valueOf of a nullable String.
//
// Go-only: language-adapter — Java string concatenation of a null String.
func symbolTable_nullableString(s *string) string {
	if s == nil {
		return "null"
	}
	return *s
}

// SymbolTable_AddSymbolsFrom makes sure all the symbols and references in
// otherSymbolTable are in this symbol table.
//
// Uniqueness of symbols and references is determined by the associated node.
//
// If multiple symbol tables are mixed in, we do not check for consistency
// between symbol tables. The first symbol we see dictates the type
// information for that symbol.
//
// Java: SymbolTable.addSymbolsFrom(StaticSymbolTable<S, R>)
func SymbolTable_AddSymbolsFrom[S rhino.StaticSlot, R rhino.StaticRef](
	t *SymbolTable, otherSymbolTable rhino.StaticSymbolTable[S, R]) {
	for otherSymbol := range otherSymbolTable.GetAllSymbols() {
		name := otherSymbol.GetName()
		myScope := t.CreateScopeFrom(otherSymbolTable.GetScope(otherSymbol))

		decl := symbolTable_findBestDeclToAdd(t, otherSymbolTable, otherSymbol)
		var mySymbol *SymbolTable_Symbol
		if decl != nil {
			declNode := decl.GetNode()

			// If we have a declaration node, we can ensure the symbol is declared.
			mySymbol = t.isAnySymbolDeclared(name, declNode, myScope)
			if mySymbol == nil {
				mySymbol = t.copySymbolToWithDeclNode(otherSymbol, declNode, myScope)
			}
		} else {
			// If we don't have a declaration node, we won't be able to declare
			// a symbol in this symbol table. But we may be able to salvage the
			// references if we already have a symbol.
			mySymbol = myScope.GetOwnSlot(name)
		}

		if mySymbol != nil {
			for otherRef := range otherSymbolTable.GetReferences(otherSymbol) {
				if t.isGoodRefToAdd(otherRef) {
					mySymbol.DefineReferenceAt(otherRef.GetNode())
				}
			}
		}
	}
}

// isAnySymbolDeclared checks if any symbol is already declared at the given
// node and scope for the given name. If so, returns it.
//
// Java: SymbolTable.isAnySymbolDeclared(String, Node, SymbolScope)
func (t *SymbolTable) isAnySymbolDeclared(name string, declNode rhino.Node, scope *SymbolTable_SymbolScope) *SymbolTable_Symbol {
	sym := t.symbols.Get(symbolTableKey{declNode: declNode, name: name})
	if sym == nil {
		// Sometimes, our symbol tables will disagree on where the
		// declaration node should be. In the rare case where this happens,
		// trust the existing symbol.
		// See SymbolTableTest#testDeclarationDisagreement.
		return scope.ownSymbols.Get(name)
	}
	return sym
}

// symbolTable_findBestDeclToAdd is a helper for addSymbolsFrom, to determine
// the best declaration spot.
//
// Java: SymbolTable.findBestDeclToAdd(StaticSymbolTable<S, R>, S)
func symbolTable_findBestDeclToAdd[S rhino.StaticSlot, R rhino.StaticRef](
	t *SymbolTable, otherSymbolTable rhino.StaticSymbolTable[S, R], slot S) rhino.StaticRef {
	decl := slot.GetDeclaration()
	if t.isGoodRefToAdd(decl) {
		return decl
	}

	for ref := range otherSymbolTable.GetReferences(slot) {
		if t.isGoodRefToAdd(ref) {
			return ref
		}
	}

	return nil
}

// isGoodRefToAdd is a helper for addSymbolsFrom, to determine whether a
// reference is acceptable. A reference must be in the normal source tree.
//
// Java: SymbolTable.isGoodRefToAdd(StaticRef)
func (t *SymbolTable) isGoodRefToAdd(ref rhino.StaticRef) bool {
	if isNilStaticRef(ref) {
		return false
	}
	node := ref.GetNode()
	return !node.IsNil() &&
		!isNilSourceFile(node.GetStaticSourceFile()) &&
		!NodeUtil_IsInSyntheticScript(node) &&
		// Typechecking assigns implicitly goog.provided names a declaration node of the expr:
		// For example, 'some' and 'some.name' in 'goog.provide('some.name.child');' have their
		// declaration node set to the entire call node. Use the actual declaration instead.
		!NodeUtil_IsGoogProvideCall(node)
}

// Java: SymbolTable.copySymbolTo(StaticSlot, SymbolScope)
func (t *SymbolTable) copySymbolTo(sym rhino.StaticSlot, scope *SymbolTable_SymbolScope) *SymbolTable_Symbol {
	return t.copySymbolToWithDeclNode(sym, sym.GetDeclaration().GetNode(), scope)
}

// Java: SymbolTable.copySymbolTo(StaticSlot, Node, SymbolScope)
func (t *SymbolTable) copySymbolToWithDeclNode(sym rhino.StaticSlot, declNode rhino.Node, scope *SymbolTable_SymbolScope) *SymbolTable_Symbol {
	// All symbols must have declaration nodes.
	rhino.CheckNotNull(declNode)
	return t.declareSymbol(
		sym.GetName(), t.getType(sym), t.isTypeInferred(sym), scope, declNode, sym.GetJSDocInfo())
}

// symbolTable_sanitizeSpecialChars replaces all \ with \\ so there will be no
// \0 or \n in the string, then replaces all '\0' (NULL) with \0 and all '\n'
// (newline) with \n.
//
// Java: SymbolTable.sanitizeSpecialChars(String)
func symbolTable_sanitizeSpecialChars(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, "\\", "\\\\"), "\x00", "\\0"), "\n", "\\n")
}

// Java: SymbolTable.addSymbol(String, JSType, boolean, SymbolScope, Node)
func (t *SymbolTable) addSymbol(
	name string, type_ rhino.JSType, inferred bool, scope *SymbolTable_SymbolScope, declNode rhino.Node) *SymbolTable_Symbol {
	name = symbolTable_sanitizeSpecialChars(name)

	symbol := NewSymbolTable_Symbol(name, type_, inferred, scope)
	replacedSymbol := t.symbols.Put(symbolTableKey{declNode: declNode, name: name}, symbol)
	jutil.CheckStateMsg2(
		replacedSymbol == nil, "Found duplicate symbol %s in global index. Type %s", name, type_)

	replacedSymbol = scope.ownSymbols.Put(name, symbol)
	jutil.CheckStateMsg2(
		replacedSymbol == nil, "Found duplicate symbol %s in its scope. Type %s", name, type_)
	return symbol
}

// Java: SymbolTable.declareSymbol(String, JSType, boolean, SymbolScope, Node, JSDocInfo)
func (t *SymbolTable) declareSymbol(
	name string,
	type_ rhino.JSType,
	inferred bool,
	scope *SymbolTable_SymbolScope,
	declNode rhino.Node,
	info *rhino.JSDocInfo) *SymbolTable_Symbol {
	symbol := t.addSymbol(name, type_, inferred, scope, declNode)
	symbol.SetJSDocInfo(info)
	symbol.SetDeclaration(symbol.DefineReferenceAt(declNode))
	return symbol
}

// mergeSymbol merges 'from' symbol to 'to' symbol by moving all references to
// point to the 'to' symbol and removing 'from' symbol.
//
// Java: SymbolTable.mergeSymbol(Symbol, Symbol)
func (t *SymbolTable) mergeSymbol(from *SymbolTable_Symbol, to *SymbolTable_Symbol) {
	for nodeToMove := range from.references.Keys() {
		if nodeToMove != from.GetDeclarationNode() {
			to.DefineReferenceAt(nodeToMove)
		}
	}
	t.removeSymbol(from)
}

// Java: SymbolTable.removeSymbol(Symbol)
func (t *SymbolTable) removeSymbol(s *SymbolTable_Symbol) {
	scope := t.GetScope(s)
	if scope.ownSymbols.Remove(s.GetName()) != s {
		panic(jutil.NewIllegalStateExceptionMsg("Symbol not found in scope " + s.String()))
	}
	removed := t.symbols.Remove(symbolTableKey{declNode: s.GetDeclaration().GetNode(), name: s.GetName()})
	if removed != s {
		panic(jutil.NewIllegalStateExceptionMsg("Symbol not found in table " + s.String()))
	}
	// If s declares a property scope then all child symbols should be removed as well.
	// For example:
	// let foo = {a: 1, b: 2};
	// foo declares property scope with a and b as its children. When removing foo we should also
	// remove a and b.
	if s.propertyScope != nil && s.propertyScope.GetSymbolForScope() == s {
		// Need to iterate over copy of values list because removeSymbol() will change the map
		// and we'll get ConcurrentModificationException
		for _, childSymbol := range slices.Collect(s.propertyScope.ownSymbols.Values()) {
			t.removeSymbol(childSymbol)
		}
		t.scopes.Remove(s.GetDeclarationNode())
	}
}

// Java: SymbolTable.renameSymbol(SymbolScope, Symbol, String)
func (t *SymbolTable) renameSymbol(scope *SymbolTable_SymbolScope, sym *SymbolTable_Symbol, newName string) {
	jutil.CheckStateMsg(
		sym.propertyScope == nil,
		"Renaming supported for simple symbols that don't have property scopes.")
	existingSym := scope.GetSlot(newName)
	if existingSym == nil {
		existingSym =
			t.declareSymbol(
				newName,
				t.getType(sym),
				t.isTypeInferred(sym),
				scope,
				sym.GetDeclarationNode(),
				sym.GetJSDocInfo())
	}
	t.mergeSymbol(sym, existingSym)
}

// Java: SymbolTable.sourceNameOrdering, nodeOrdering, lexicalScopeOrdering, symbolOrdering
func (t *SymbolTable) initOrderings() {
	t.sourceNameOrdering = func(a, b *string) int {
		// Ordering.natural().nullsFirst()
		if a == b {
			return 0
		}
		if a == nil {
			return -1
		}
		if b == nil {
			return 1
		}
		return jutil.StringCompareTo(*a, *b)
	}

	t.nodeOrdering = func(a, b rhino.Node) int {
		result := t.sourceNameOrdering(a.GetSourceFileName(), b.GetSourceFileName())
		if result != 0 {
			return result
		}

		// Source position is a bit mask of line in the top 4 bits, so this
		// is a quick way to compare order without computing absolute position.
		return int(int32(a.GetSourcePosition()) - int32(b.GetSourcePosition()))
	}

	t.lexicalScopeOrdering = func(a, b *SymbolTable_SymbolScope) int {
		jutil.CheckStateMsg(a.IsLexicalScope() && b.IsLexicalScope(), "We can only sort lexical scopes")
		result := t.nodeOrdering(a.GetRootNode(), b.GetRootNode())
		if result != 0 {
			return result
		}

		// If result = 0 it means that rootNodes either the same or that one of them was added
		// by compiler during transpilation and uses the same source info as original node.
		// In that case compare scopes by depth because one of them (probably generated one) is
		// a child of the other scope.
		// TODO(b/62349230): remove this once transpilation is disabled. No two different scopes
		// should have the same source info.
		return a.GetScopeDepth() - b.GetScopeDepth()
	}

	t.symbolOrdering = func(a, b *SymbolTable_Symbol) int {
		scopeA := t.GetScope(a)
		scopeB := t.GetScope(b)

		// More deeply nested symbols should go later.
		result := t.getLexicalScopeDepth(scopeA) - t.getLexicalScopeDepth(scopeB)
		if result != 0 {
			return result
		}

		// After than, just use lexicographic ordering.
		// This ensures "a.b" comes before "a.b.c".
		result = jutil.StringCompareTo(a.GetName(), b.GetName())
		if result != 0 {
			return result
		}

		if scopeA.IsLexicalScope() && scopeB.IsLexicalScope() {
			return t.lexicalScopeOrdering(scopeA, scopeB)
		}
		// nodeOrdering.nullsFirst()
		declA, declB := a.GetDeclarationNode(), b.GetDeclarationNode()
		if declA == declB {
			return 0
		}
		if declA.IsNil() {
			return -1
		}
		if declB.IsNil() {
			return 1
		}
		return t.nodeOrdering(declA, declB)
	}
}

// getLexicalScopeDepth returns, for a lexical scope, just the normal scope
// depth.
//
// For a property scope, returns the number of scopes we have to search to
// find the nearest lexical scope, plus that lexical scope's depth.
//
// For a doc info scope, returns 0.
//
// Java: SymbolTable.getLexicalScopeDepth(SymbolScope)
func (t *SymbolTable) getLexicalScopeDepth(scope *SymbolTable_SymbolScope) int {
	if scope.IsLexicalScope() || scope.IsDocScope() {
		return scope.GetScopeDepth()
	} else {
		jutil.CheckState(scope.IsPropertyScope())
		sym := scope.GetSymbolForScope()
		jutil.CheckNotNull(sym)
		return t.getLexicalScopeDepth(t.GetScope(sym)) + 1
	}
}

// Java: SymbolTable.getType(StaticSlot)
func (t *SymbolTable) getType(sym rhino.StaticSlot) rhino.JSType {
	if docInfo := sym.GetJSDocInfo(); docInfo != nil && docInfo.HasTypedefType() {
		type_ := t.registry.GetGlobalType(sym.GetName())
		if type_ != nil {
			return type_
		}
		return docInfo.GetTypedefType().Evaluate(symbolTable_topScope(t.compiler), t.registry)
	}
	if staticTypedSlot, ok := sym.(rhino.StaticTypedSlot); ok {
		return staticTypedSlot.GetType()
	}
	return nil
}

// symbolTable_topScope is compiler.getTopScope() as a StaticTypedScope (a
// null TypedScope is a nil interface).
//
// Go-only: language-adapter — typed nil.
func symbolTable_topScope(compiler AbstractCompiler) rhino.StaticTypedScope {
	if s := compiler.GetTopScope(); s != nil {
		return s
	}
	return nil
}

// Java: SymbolTable.getTypeOfThis(StaticScope)
func (t *SymbolTable) getTypeOfThis(s rhino.StaticScope) rhino.JSType {
	if staticTypedScope, ok := s.(rhino.StaticTypedScope); ok {
		return staticTypedScope.GetTypeOfThis()
	}
	return nil
}

// Java: SymbolTable.isTypeInferred(StaticSlot)
func (t *SymbolTable) isTypeInferred(sym rhino.StaticSlot) bool {
	if staticTypedSlot, ok := sym.(rhino.StaticTypedSlot); ok {
		return staticTypedSlot.IsTypeInferred()
	}
	return true
}
