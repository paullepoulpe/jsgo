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
	"fmt"
	"iter"
	"slices"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/base"
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
//
// Java: public final class SymbolTable
type SymbolTable struct {
	// All symbols in the program, uniquely identified by the node where
	// they're declared and their name.
	symbols *jutil.HashBasedTable[*rhino.Node, string, *SymbolTable_Symbol]

	// All syntactic scopes in the program, uniquely identified by the node
	// where they're declared.
	scopes *jutil.LinkedHashMap[*rhino.Node, *SymbolTable_SymbolScope]

	// All Nodes with JSDocInfo in the program.
	docInfos *jutil.ArrayList[*rhino.Node]

	globalScope *SymbolTable_SymbolScope

	compiler AbstractCompiler

	registry *rhino.JSTypeRegistry

	// Map from a JSType to the symbol that declares it.
	symbolThatDeclaresType *base.LinkedIdentityHashMap[rhino.JSType, *SymbolTable_Symbol]

	// Comparators
	sourceNameOrdering   func(a, b *string) int
	nodeOrdering         func(a, b *rhino.Node) int
	lexicalScopeOrdering func(a, b *SymbolTable_SymbolScope) int
	symbolOrdering       func(a, b *SymbolTable_Symbol) int
}

var symbolTable_logger = jutil.Logger_GetLogger("com.google.javascript.jscomp.SymbolTable")

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
		symbols:                jutil.HashBasedTable_Create[*rhino.Node, string, *SymbolTable_Symbol](),
		scopes:                 jutil.NewLinkedHashMap[*rhino.Node, *SymbolTable_SymbolScope](),
		docInfos:               jutil.NewArrayList[*rhino.Node](),
		symbolThatDeclaresType: base.NewLinkedIdentityHashMap[rhino.JSType, *SymbolTable_Symbol](),
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
func (t *SymbolTable) GetAllJSDocInfoNodes() []*rhino.Node {
	return t.docInfos.Slice()
}

// GetEnclosingScope gets the scope that contains the given node. If n is a
// function name, we return the scope that contains the function, not the
// function itself.
//
// Java: SymbolTable.getEnclosingScope(Node)
func (t *SymbolTable) GetEnclosingScope(n *rhino.Node) *SymbolTable_SymbolScope {
	current := n.GetParent()
	if n.IsName() && n.GetParent().IsFunction() {
		current = current.GetParent()
	}

	for ; current != nil; current = current.GetParent() {
		if t.scopes.ContainsKey(current) {
			return t.scopes.Get(current)
		}
	}
	return nil
}

// GetEnclosingFunctionScope gets the scope that contains the given node. If n
// is a function name, we return the scope that contains the function, not the
// function itself. The returned scope is either function or global scope.
//
// Java: SymbolTable.getEnclosingFunctionScope(Node)
func (t *SymbolTable) GetEnclosingFunctionScope(n *rhino.Node) *SymbolTable_SymbolScope {
	current := n.GetParent()
	if n.IsName() && current != nil && current.IsFunction() {
		current = current.GetParent()
	}

	for ; current != nil; current = current.GetParent() {
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
	if functionNode == nil {
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
	if rootNode.GetParent() == nil {
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
func (t *SymbolTable) GetSymbolForName(source *rhino.Node, name *string) *SymbolTable_Symbol {
	if name == nil || t.globalScope == nil {
		return nil
	}
	nameStr := *name

	var scope *SymbolTable_SymbolScope
	if source == nil {
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
			metadata = t.compiler.
				GetModuleMetadataMap().
				GetModulesByPath().
				Get(*sourceFileName)
		}
		if metadata == nil {
			panic(jutil.NewNullPointerException())
		}
		first := "null"
		for ns := range metadata.GoogNamespaces().All() {
			first = ns
			break
		}
		moduleNameWithDot := first + "."
		if strings.HasPrefix(nameStr, moduleNameWithDot) {
			nameStr = jutil.StringSubstringFrom(nameStr, jutil.StringLength(moduleNameWithDot))
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
	if jutil.IsNil(type_) {
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
	if jutil.IsNil(type_) {
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
	} else if !jutil.IsNil(type_.AutoboxesTo()) {
		return t.GetSymbolForTypeHelper(type_.AutoboxesTo(), linkToCtor)
	} else if type_.IsEnumType() {
		return t.GetSymbolDeclaredByEnum(type_.(*rhino.EnumType))
	} else if type_.IsRecordType() {
		return t.symbolThatDeclaresType.Get(type_)
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
	childrenScopes := jutil.NewLinkedHashMap[*SymbolTable_SymbolScope, []*SymbolTable_SymbolScope]()
	for scope := range t.GetAllScopes() {
		if scope.GetParentScope() != nil {
			childrenScopes.Put(scope.GetParentScope(), append(childrenScopes.Get(scope.GetParentScope()), scope))
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
	visitedScopes := jutil.NewLinkedHashSet[*SymbolTable_SymbolScope]()
	for scope := range t.GetAllScopes() {
		visitedScopes.Add(scope)
	}
	for _, symbol := range t.GetAllSymbols() {
		if symbol.propertyScope == nil || visitedScopes.Contains(symbol.propertyScope) {
			continue
		}
		visitedScopes.Add(symbol.GetPropertyScope())
		t.toDebugStringTree(&builder, "", symbol.propertyScope, childrenScopes)
	}
	return builder.String()
}

// Java: SymbolTable.toDebugStringTree(StringBuilder, String, SymbolScope, ImmutableListMultimap<SymbolScope, SymbolScope>)
func (t *SymbolTable) toDebugStringTree(
	builder *strings.Builder,
	prefix string,
	scope *SymbolTable_SymbolScope,
	childrenScopes *jutil.LinkedHashMap[*SymbolTable_SymbolScope, []*SymbolTable_SymbolScope]) {
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
	if childrenScopes.ContainsKey(scope) {
		printedSomething = true
		builder.WriteString("\n" + prefix + "  Scopes:\n")
		for _, childScope := range childrenScopes.Get(scope) {
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
		builder.WriteString(fmt.Sprintf(
			"  Ref %d: %s line: %d col: %d len: %d %s\n",
			refCount,
			symbolTable_nullableString(node.GetSourceFileName()),
			node.GetLineno(),
			node.GetCharno(),
			node.GetLength(),
			indexable))
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
func (t *SymbolTable) FindScopes(externs *rhino.Node, root *rhino.Node) {
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
func (f *symbolTable_scopeFinder) Visit(t *NodeTraversal, n *rhino.Node, p *rhino.Node) {}

// GetAllScopes gets all the scopes in this symbol table.
//
// Java: SymbolTable.getAllScopes()
func (t *SymbolTable) GetAllScopes() iter.Seq[*SymbolTable_SymbolScope] {
	return t.scopes.Values()
}

// AddAnonymousFunctions finds anonymous functions in local scopes, and gives
// them names and symbols. They will show up as local variables with names
// "function%0", "function%1", etc.
//
// Java: SymbolTable.addAnonymousFunctions()
func (t *SymbolTable) AddAnonymousFunctions() {
	scopes := jutil.NewTreeSet[*SymbolTable_SymbolScope](t.lexicalScopeOrdering)
	for scope := range t.GetAllScopes() {
		if scope.IsLexicalScope() {
			scopes.Add(scope)
		}
	}

	for scope := range scopes.All() {
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
			rootNode != nil &&
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
func (t *SymbolTable) isAnySymbolDeclared(name string, declNode *rhino.Node, scope *SymbolTable_SymbolScope) *SymbolTable_Symbol {
	sym, _ := t.symbols.Get(declNode, name)
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
	return !jutil.IsNil(ref) &&
		ref.GetNode() != nil &&
		!jutil.IsNil(ref.GetNode().GetStaticSourceFile()) &&
		!NodeUtil_IsInSyntheticScript(ref.GetNode()) &&
		// Typechecking assigns implicitly goog.provided names a declaration node of the expr:
		// For example, 'some' and 'some.name' in 'goog.provide('some.name.child');' have their
		// declaration node set to the entire call node. Use the actual declaration instead.
		!NodeUtil_IsGoogProvideCall(ref.GetNode())
}

// Java: SymbolTable.copySymbolTo(StaticSlot, SymbolScope)
func (t *SymbolTable) copySymbolTo(sym rhino.StaticSlot, scope *SymbolTable_SymbolScope) *SymbolTable_Symbol {
	return t.copySymbolToWithDeclNode(sym, sym.GetDeclaration().GetNode(), scope)
}

// Java: SymbolTable.copySymbolTo(StaticSlot, Node, SymbolScope)
func (t *SymbolTable) copySymbolToWithDeclNode(sym rhino.StaticSlot, declNode *rhino.Node, scope *SymbolTable_SymbolScope) *SymbolTable_Symbol {
	// All symbols must have declaration nodes.
	jutil.CheckNotNull(declNode)
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
	name string, type_ rhino.JSType, inferred bool, scope *SymbolTable_SymbolScope, declNode *rhino.Node) *SymbolTable_Symbol {
	name = symbolTable_sanitizeSpecialChars(name)

	symbol := NewSymbolTable_Symbol(name, type_, inferred, scope)
	replacedSymbol := t.symbols.Put(declNode, name, symbol)
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
	declNode *rhino.Node,
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
	if !s.Equals(scope.ownSymbols.Remove(s.GetName())) {
		panic(jutil.NewIllegalStateExceptionMsg("Symbol not found in scope " + s.String()))
	}
	removed, _ := t.symbols.Remove(s.GetDeclaration().GetNode(), s.GetName())
	if !s.Equals(removed) {
		panic(jutil.NewIllegalStateExceptionMsg("Symbol not found in table " + s.String()))
	}
	// If s declares a property scope then all child symbols should be removed as well.
	// For example:
	// let foo = {a: 1, b: 2};
	// foo declares property scope with a and b as its children. When removing foo we should also
	// remove a and b.
	if s.propertyScope != nil && s.Equals(s.propertyScope.GetSymbolForScope()) {
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

			for currentNode.IsGetProp() {
				currentNode = currentNode.GetFirstChild()

				namePtr := currentNode.GetQualifiedName()
				if namePtr == nil {
					continue
				}
				name := *namePtr
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
					if !jutil.IsNil(nodeType) && !nodeType.IsFunctionPrototypeType() {
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
			!s.Equals(symbolForType) {
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
		t.symbolThatDeclaresType.Put(type_, s)
	}

	for _, s := range allTypes {
		type_ := t.getType(s)

		// Create property scopes only based on "root" symbols for each type to handle aliases.
		if jutil.IsNil(type_) || s.Equals(t.symbolThatDeclaresType.Get(type_)) {
			t.createPropertyScopeFor(s)
		}
	}

	// Now we need to set the new property scope symbol to all aliases.
	for _, s := range allTypes {
		type_ := t.getType(s)

		if !jutil.IsNil(type_) && t.symbolThatDeclaresType.Get(type_) != nil {
			s.propertyScope = t.symbolThatDeclaresType.Get(type_).GetPropertyScope()
		}
	}

	t.PruneOrphanedNames()
}

// Java: SymbolTable.needsPropertyScope(Symbol)
func (t *SymbolTable) needsPropertyScope(sym *SymbolTable_Symbol) bool {
	jsType := t.getType(sym)
	if !jutil.IsNil(jsType) && jsType.IsRecordType() {
		return true
	}
	type_ := rhino.ObjectType_Cast(t.getType(sym))
	if jutil.IsNil(type_) {
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
			if owner != nil &&
				!jutil.IsNil(t.getType(owner)) &&
				(t.getType(owner).IsNominalConstructorOrInterface() ||
					t.getType(owner).IsFunctionPrototypeType() ||
					t.getType(owner).IsEnumType()) {
				t.removeSymbol(s)
				continue nextSymbol
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
func (t *SymbolTable) FillPropertySymbols(externs *rhino.Node, root *rhino.Node) {
	NewSymbolTable_PropertyRefCollector(t).Process(externs, root)
}

// FillJSDocInfo indexes JSDocInfo.
//
// Java: SymbolTable.fillJSDocInfo(Node, Node)
func (t *SymbolTable) FillJSDocInfo(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_TraverseRoots(
		t.compiler, NewSymbolTable_JSDocInfoCollector(t, t.compiler.GetTypeRegistry()), externs, root)

	// Create references to parameters in the JSDoc.
	for _, sym := range t.GetAllSymbols() {
		info := sym.GetJSDocInfo()
		if info == nil {
			continue
		}

		for _, marker := range info.GetMarkers() {
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
							nil /* root */, nil /* parent scope */, nil /* type of this */, sym)
				}

				// Check to make sure there's no existing symbol. In theory, this
				// should never happen, but we check anyway and fail silently
				// if our assumptions are wrong. (We do not want to put the symbol
				// table into an invalid state).
				existingSymbol := t.isAnySymbolDeclared(name, paramNode, sym.docScope)
				if existingSymbol == nil {
					t.declareSymbol(name, type_, jutil.IsNil(type_), sym.docScope, paramNode /* info= */, nil)
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
func (t *SymbolTable) FillSymbolVisibility(externs *rhino.Node, root *rhino.Node) {
	collectPass := NewCollectFileOverviewVisibility(t.compiler)
	collectPass.Process(externs, root)
	visibilityMap := collectPass.GetFileOverviewVisibilityMap()

	NodeTraversal_TraverseRoots(t.compiler, NewSymbolTable_VisibilityCollector(t, visibilityMap), externs, root)
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
	if !jutil.IsNil(jsType) {
		type_ = jsType.ToObjectType()
	}
	if jutil.IsNil(type_) {
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
	s.SetPropertyScope(NewSymbolTable_SymbolScope(nil, parentPropertyScope, type_, s))

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
		if jutil.IsNil(newProp.GetDeclaration()) {
			// Skip properties without declarations. We won't know how to index
			// them, because we index things by node.
			continue
		}

		// We have symbol tables that do not do type analysis. They just try
		// to build a complete index of all objects in the program. So we might
		// already have symbols for things like "Foo.bar". If this happens,
		// throw out the old symbol and use the type-based symbol.
		oldProp, _ :=
			t.symbols.Get(newProp.GetDeclaration().GetNode(), s.GetName()+"."+propName)

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
		if existing, _ := t.symbols.Get(newProp.GetDeclaration().GetNode(), newProp.GetName()); existing != nil {
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
	if jutil.IsNil(proto) || proto == symbolObjectType {
		return nil
	}
	var parentSymbol *SymbolTable_Symbol
	if t.isEs6ClassConstructor(proto) {
		// given `class Foo {} class Bar extends Foo {}`, `Foo` is the implicit prototype of `Bar`.
		parentSymbol = t.GetSymbolDeclaredBy(proto.ToMaybeFunctionType())
	} else if !jutil.IsNil(proto.GetConstructor()) {
		// given
		//   /** @constructor */ function Foo() {}
		//   /** @constructor */ function Bar() {}
		//   goog.inherits(Bar, Foo);
		// the implicit prototype of Bar.prototype is the instance of Foo.
		parentSymbol = t.GetSymbolForInstancesOfFunction(proto.GetConstructor())
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
	return type_.IsFunctionType() &&
		type_.ToMaybeFunctionType().GetSource() != nil &&
		type_.ToMaybeFunctionType().GetSource().IsClass()
}

// FillThisReferences fills in references to "this" variables.
//
// Java: SymbolTable.fillThisReferences(Node, Node)
func (t *SymbolTable) FillThisReferences(externs *rhino.Node, root *rhino.Node) {
	NewSymbolTable_ThisRefCollector(t).Process(externs, root)
}

// FillSuperReferences fills in references to "super" variables.
//
// Java: SymbolTable.fillSuperReferences(Node, Node)
func (t *SymbolTable) FillSuperReferences(externs *rhino.Node, root *rhino.Node) {
	collectSuper := &symbolTable_superCollector{outer: t}
	NodeTraversal_TraverseRoots(t.compiler, collectSuper, externs, root)
}

// symbolTable_superCollector is the anonymous AbstractPostOrderCallback of
// fillSuperReferences.
//
// Go-only: language-adapter — an anonymous class.
type symbolTable_superCollector struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *SymbolTable
}

// Java: SymbolTable.fillSuperReferences(...).new AbstractPostOrderCallback() {...}.visit(NodeTraversal, Node, Node)
func (c *symbolTable_superCollector) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	// Process only 'super' nodes with types.
	if !n.IsSuper() || jutil.IsNil(n.GetJSType()) {
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
func (t *SymbolTable) addRefsInGoogRequireStatement(n *rhino.Node, googModuleScopes *jutil.LinkedHashMap[string, *SymbolTable_SymbolScope]) {
	if t.compiler.GetOptions().ShouldRewriteModules() {
		// When modules are rewritten - there is no need to connect imported names back to definitions
		// as that was done with the help of Preprocessor tables in ClosureRewriteModule pass.
		return
	}
	require := n.GetLastChild()
	moduleName := require.GetSecondChild().GetString()
	module := t.compiler.GetModuleMap().GetClosureModule(moduleName)
	// Go-only: language-adapter — a null ModuleType (no module) makes the
	// Java switch statements below throw a NullPointerException.
	var moduleType *ModuleMetadataMap_ModuleType
	if module != nil {
		mt := module.Metadata().ModuleType()
		moduleType = &mt
	}
	moduleScope := googModuleScopes.Get(moduleName)
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
		if moduleType == nil {
			panic(jutil.NewNullPointerException())
		}
		switch *moduleType {
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
			declaration = t.GetSymbolForName(nil, &moduleName)
		default:
			// skip
		}
		if declaration != nil {
			declaration.DefineReferenceAt(n)
		}
	} else if n.IsDestructuringLhs() {
		// This is `const {one} = goog.require('some.foo')` case.
		// For each imported symbol in destructuring we need to find corresponding declaration symbol
		// in `some.foo` namespace and connect them.
		jutil.CheckState(n.GetFirstChild().IsObjectPattern())
		// Iterate through all symbols in destructuring.
		for stringKey := range n.GetFirstChild().Children() {
			jutil.CheckState(stringKey.IsStringKey())
			varName := stringKey.GetString()
			var varDeclaration *SymbolTable_Symbol
			// AST of goog.module and goog.provide differs significantly so we need to lookup variables
			// differently.
			if moduleType == nil {
				panic(jutil.NewNullPointerException())
			}
			switch *moduleType {
			case ModuleMetadataMap_ModuleType_GOOG_MODULE, ModuleMetadataMap_ModuleType_LEGACY_GOOG_MODULE:
				if moduleScope == nil {
					panic(jutil.NewNullPointerException())
				}
				varDeclaration = moduleScope.GetOwnSlot(varName)
			case ModuleMetadataMap_ModuleType_GOOG_PROVIDE:
				qname := moduleName + "." + varName
				varDeclaration = t.GetSymbolForName(nil, &qname)
			default:
				// skip
			}
			if varDeclaration != nil {
				varDeclaration.DefineReferenceAt(stringKey)
			}
		}
	}
}

// FillGoogProvideModuleRequires connects goog.module/goog.provide calls with
// goog.require/goog.requireType. For each goog.module/goog.provide call
// creates a symbol and all goog.require/goog.requireType added as references.
//
// Java: SymbolTable.fillGoogProvideModuleRequires(Node, Node)
func (t *SymbolTable) FillGoogProvideModuleRequires(externs *rhino.Node, root *rhino.Node) {
	// small optimization to keep symbols created within this function for fast access instead of
	// going through SymbolTable.getSymbolForName() every time.
	declaredNamespaces := jutil.NewLinkedHashMap[string, *SymbolTable_Symbol]()
	// Map containind goog.module names to Scope object representing them. For goog.provide
	// namespaces there is no scope as elements of goog.provide are global.
	moduleScopes := jutil.NewLinkedHashMap[string, *SymbolTable_SymbolScope]()
	looksLikeGoogCall :=
		func(n *rhino.Node) bool {
			if !n.IsCall() {
				return false
			}
			left := n.GetFirstChild()
			arg := n.GetSecondChild()
			// We want only nodes of type `goog.xyz('somestring')`.
			if !(left.IsGetProp() &&
				left.GetFirstChild().IsName() &&
				left.GetFirstChild().GetString() == "goog") ||
				arg == nil ||
				!arg.IsStringLit() {
				return false
			}
			return true
		}
	collectModuleScopes :=
		&symbolTable_postOrderCallback{visit: func(nt *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
			if !looksLikeGoogCall(n) {
				return
			}
			arg := n.GetSecondChild()
			namespaceName := "ns$" + arg.GetString()
			switch n.GetFirstChild().GetString() {
			case "module", "provide":
				ns :=
					t.declareSymbol(
						namespaceName,
						/* type= */ nil,
						/* inferred= */ false,
						t.GetGlobalScope(),
						arg,
						/* info= */ nil)
				declaredNamespaces.Put(namespaceName, ns)
				if n.GetGrandparent().IsModuleBody() {
					moduleScopes.Put(arg.GetString(), t.scopes.Get(n.GetGrandparent()))
				}
			default:
				// do nothing. Some other goog.xyz call.
			}
		}}
	processRequireStatements :=
		&symbolTable_postOrderCallback{visit: func(nt *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
			if !looksLikeGoogCall(n) {
				return
			}
			arg := n.GetSecondChild()
			namespaceName := "ns$" + arg.GetString()
			switch n.GetFirstChild().GetString() {
			case "require", "requireType", "forwardDeclare":
				t.addRefsInGoogRequireStatement(n.GetParent(), moduleScopes)
				symbol := declaredNamespaces.Get(namespaceName)
				// We expect that namespace was already processed by that point, but in some broken
				// code it might be missing. In which case just skip it.
				if symbol != nil {
					symbol.DefineReferenceAt(arg)
				}
			default:
				// do nothing. Some other goog.xyz call.
			}
		}}
	NodeTraversal_TraverseRoots(t.compiler, collectModuleScopes, externs, root)
	NodeTraversal_TraverseRoots(t.compiler, processRequireStatements, externs, root)
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
	return node != nil && node.IsStringKey() && node.IsQuotedStringKey()
}

// isSymbolDuplicatedExternOnWindow is a heuristic method to check whether
// symbol was created by DeclaredGlobalExternsOnWindow.java pass.
//
// Java: SymbolTable.isSymbolDuplicatedExternOnWindow(Symbol)
func (t *SymbolTable) isSymbolDuplicatedExternOnWindow(symbol *SymbolTable_Symbol) bool {
	node := symbol.GetDeclarationNode()
	// Check that node is of type "window.foo";
	return !node.IsIndexable() &&
		node.IsGetProp() &&
		node.GetFirstChild().IsName() &&
		node.GetFirstChild().GetString() == "window"
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
	externSymbols := jutil.HashBasedTable_Create[string, int, *SymbolTable_Symbol]()
	for _, symbol := range t.symbols.ValueSlice() {
		if symbol.GetDeclarationNode() == nil ||
			jutil.IsNil(symbol.GetDeclarationNode().GetStaticSourceFile()) ||
			!symbol.GetDeclarationNode().GetStaticSourceFile().IsExtern() {
			continue
		}
		sourceFilePtr := symbol.GetSourceFileName()
		if sourceFilePtr == nil {
			// Java: HashBasedTable.put(null, ...) throws NullPointerException.
			panic(jutil.NewNullPointerException())
		}
		sourceFile := *sourceFilePtr
		position := symbol.GetDeclarationNode().GetSourcePosition()
		if !externSymbols.Contains(sourceFile, position) {
			externSymbols.Put(sourceFile, position, symbol)
			continue
		}
		existingSymbol, _ := externSymbols.Get(sourceFile, position)
		// Consider 2 possibilies: either symbol or existingSymbol might be the generated symbol we
		// are looking for.
		if t.isSymbolDuplicatedExternOnWindow(existingSymbol) {
			t.mergeSymbol(existingSymbol, symbol)
			externSymbols.Put(sourceFile, position, symbol)
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
		} else if symbol.GetDeclarationNode() != nil &&
			symbol.GetDeclarationNode().IsModuleBody() {
			// Symbols that represent whole module (goog.module) aren't needed for indexing. We already
			// recorded imports/references of that module in fillGoogProvideModuleRequires() function.
			// Individual exported functions/variables also have their own symbols.
			t.removeSymbol(symbol)
		}
	}
	t.mergeExternSymbolsDuplicatedOnWindow()
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

		if jutil.IsNil(otherScopeParent) {
			// The global scope must be created before any local scopes.
			jutil.CheckStateMsg(t.globalScope == nil, "Global scopes found at different roots")
		}

		var parent *SymbolTable_SymbolScope
		if !jutil.IsNil(otherScopeParent) {
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

// symbolTable_postOrderCallback is an anonymous AbstractPostOrderCallback of
// fillGoogProvideModuleRequires.
//
// Go-only: language-adapter — an anonymous class.
type symbolTable_postOrderCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	visit func(t *NodeTraversal, n *rhino.Node, parent *rhino.Node)
}

// Java: SymbolTable.fillGoogProvideModuleRequires(...).new AbstractPostOrderCallback() {...}.visit(NodeTraversal, Node, Node)
func (c *symbolTable_postOrderCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	c.visit(t, n, parent)
}

// SymbolTable_Symbol is a symbol-table entry.
//
// Java: public static final class SymbolTable.Symbol extends SimpleSlot
type SymbolTable_Symbol struct {
	rhino.SimpleSlotBase

	// Use a linked hash map, so that the results are deterministic
	// (and so the declaration always comes first).
	references *jutil.LinkedHashMap[*rhino.Node, *SymbolTable_Reference]

	scope *SymbolTable_SymbolScope

	propertyScope *SymbolTable_SymbolScope

	declaration *SymbolTable_Reference

	docInfo *rhino.JSDocInfo

	// Stored separately from docInfo, because the visibility stored in
	// JSDocInfo is not
	visibility *rhino.JSDocInfo_Visibility

	// A scope for symbols that are only documented in JSDoc.
	docScope *SymbolTable_SymbolScope
}

// Java: SymbolTable.Symbol(String, JSType, boolean, SymbolScope)
func NewSymbolTable_Symbol(name string, type_ rhino.JSType, inferred bool, scope *SymbolTable_SymbolScope) *SymbolTable_Symbol {
	s := &SymbolTable_Symbol{
		references: jutil.NewLinkedHashMap[*rhino.Node, *SymbolTable_Reference](),
		scope:      scope,
	}
	s.InitSimpleSlotBase(name, type_, inferred)
	return s
}

// Java: SymbolTable.Symbol.equals(Object)
func (s *SymbolTable_Symbol) Equals(o any) bool {
	other, ok := o.(*SymbolTable_Symbol)
	if !ok || other == nil {
		return false
	}

	return s.IsTypeInferred() == other.IsTypeInferred() &&
		s.GetName() == other.GetName() &&
		jutil.ObjectsEquals(s.GetType(), other.GetType()) &&
		s.scope == other.scope
}

// Java: SymbolTable.Symbol.hashCode()
func (s *SymbolTable_Symbol) HashCode() int32 {
	return jutil.ObjectsHash(s.IsTypeInferred(), s.GetName(), s.GetType(), jutil.IdentityHashCode(s.scope))
}

// GetDeclaration returns the declaration reference (a *SymbolTable_Reference,
// or nil).
//
// Java: SymbolTable.Symbol.getDeclaration()
func (s *SymbolTable_Symbol) GetDeclaration() rhino.StaticRef {
	if s.declaration == nil {
		return nil
	}
	return s.declaration
}

// Java: SymbolTable.Symbol.getFunctionType()
func (s *SymbolTable_Symbol) GetFunctionType() rhino.FunctionType {
	return rhino.JSType_ToMaybeFunctionType(s.GetType())
}

// Java: SymbolTable.Symbol.defineReferenceAt(Node)
func (s *SymbolTable_Symbol) DefineReferenceAt(n *rhino.Node) *SymbolTable_Reference {
	return s.references.ComputeIfAbsent(n, func(k *rhino.Node) *SymbolTable_Reference { return NewSymbolTable_Reference(s, k) })
}

// SetDeclaration sets the declaration node. May only be called once.
//
// Java: SymbolTable.Symbol.setDeclaration(Reference)
func (s *SymbolTable_Symbol) SetDeclaration(ref *SymbolTable_Reference) {
	jutil.CheckState(s.declaration == nil)
	s.declaration = ref
}

// Java: SymbolTable.Symbol.getDeclarationNode()
func (s *SymbolTable_Symbol) GetDeclarationNode() *rhino.Node {
	if s.declaration == nil {
		return nil
	}
	return s.declaration.GetNode()
}

// Java: SymbolTable.Symbol.getSourceFileName()
func (s *SymbolTable_Symbol) GetSourceFileName() *string {
	n := s.GetDeclarationNode()
	if n == nil {
		return nil
	}
	return n.GetSourceFileName()
}

// Java: SymbolTable.Symbol.getPropertyScope()
func (s *SymbolTable_Symbol) GetPropertyScope() *SymbolTable_SymbolScope {
	return s.propertyScope
}

// Java: SymbolTable.Symbol.setPropertyScope(SymbolScope)
func (s *SymbolTable_Symbol) SetPropertyScope(scope *SymbolTable_SymbolScope) {
	s.propertyScope = scope
	if scope != nil {
		s.propertyScope.SetSymbolForScope(s)
	}
}

// Java: SymbolTable.Symbol.getJSDocInfo()
func (s *SymbolTable_Symbol) GetJSDocInfo() *rhino.JSDocInfo {
	return s.docInfo
}

// Java: SymbolTable.Symbol.setJSDocInfo(JSDocInfo)
func (s *SymbolTable_Symbol) SetJSDocInfo(info *rhino.JSDocInfo) {
	s.docInfo = info
}

// Java: SymbolTable.Symbol.getVisibility()
func (s *SymbolTable_Symbol) GetVisibility() *rhino.JSDocInfo_Visibility {
	return s.visibility
}

// Java: SymbolTable.Symbol.setVisibility(Visibility)
func (s *SymbolTable_Symbol) SetVisibility(v rhino.JSDocInfo_Visibility) {
	s.visibility = &v
}

// IsProperty reports whether this is a property of another variable.
//
// Java: SymbolTable.Symbol.isProperty()
func (s *SymbolTable_Symbol) IsProperty() bool {
	return s.scope.IsPropertyScope()
}

// IsLexicalVariable reports whether this is a variable in a lexical scope.
//
// Java: SymbolTable.Symbol.isLexicalVariable()
func (s *SymbolTable_Symbol) IsLexicalVariable() bool {
	return s.scope.IsLexicalScope()
}

// IsDocOnlyParameter reports whether this is a variable that's only in JSDoc.
//
// Java: SymbolTable.Symbol.isDocOnlyParameter()
func (s *SymbolTable_Symbol) IsDocOnlyParameter() bool {
	return s.scope.IsDocScope()
}

// Java: SymbolTable.Symbol.toString()
func (s *SymbolTable_Symbol) String() string {
	n := s.GetDeclarationNode()
	lineNo := -1
	if n != nil {
		lineNo = n.GetLineno()
	}
	return s.GetName() + "@" + symbolTable_nullableString(s.GetSourceFileName()) + ":" + strconv.Itoa(lineNo)
}

// SymbolTable_Reference is a reference.
//
// Java: public static final class SymbolTable.Reference extends SimpleReference<Symbol>
type SymbolTable_Reference struct {
	rhino.SimpleReference[*SymbolTable_Symbol]
}

// Java: SymbolTable.Reference(Symbol, Node)
func NewSymbolTable_Reference(symbol *SymbolTable_Symbol, node *rhino.Node) *SymbolTable_Reference {
	r := &SymbolTable_Reference{}
	r.InitSimpleReference(symbol, node)
	return r
}

// SymbolTable_SymbolScope is the scope of a symbol.
//
// Java: public static final class SymbolTable.SymbolScope
type SymbolTable_SymbolScope struct {
	rootNode   *rhino.Node
	parent     *SymbolTable_SymbolScope
	typeOfThis rhino.JSType
	ownSymbols *jutil.LinkedHashMap[string, *SymbolTable_Symbol]
	scopeDepth int

	// Index of current scope in the parent scope. Used to generate unique names for variables
	// having the same name but defined in different scopes.
	indexInParent int

	numberOfChildScopes int

	// The symbol associated with a property scope or doc scope.
	mySymbol *SymbolTable_Symbol
}

// Java: SymbolTable.SymbolScope(Node, SymbolScope, JSType, Symbol)
func NewSymbolTable_SymbolScope(
	rootNode *rhino.Node,
	parent *SymbolTable_SymbolScope,
	typeOfThis rhino.JSType,
	mySymbol *SymbolTable_Symbol) *SymbolTable_SymbolScope {
	s := &SymbolTable_SymbolScope{
		rootNode:   rootNode,
		parent:     parent,
		typeOfThis: typeOfThis,
		ownSymbols: jutil.NewLinkedHashMap[string, *SymbolTable_Symbol](),
		mySymbol:   mySymbol,
	}
	if parent == nil {
		s.scopeDepth = 0
	} else {
		s.scopeDepth = parent.GetScopeDepth() + 1
	}
	if parent == nil {
		s.indexInParent = 0
	} else {
		s.indexInParent = parent.numberOfChildScopes
		parent.numberOfChildScopes++
	}
	return s
}

// Java: SymbolTable.SymbolScope.getSymbolForScope()
func (s *SymbolTable_SymbolScope) GetSymbolForScope() *SymbolTable_Symbol {
	return s.mySymbol
}

// Java: SymbolTable.SymbolScope.setSymbolForScope(Symbol)
func (s *SymbolTable_SymbolScope) SetSymbolForScope(sym *SymbolTable_Symbol) {
	s.mySymbol = sym
}

// GetIndexOfSymbol gets a unique index for the symbol in this scope.
//
// Java: SymbolTable.SymbolScope.getIndexOfSymbol(Symbol)
func (s *SymbolTable_SymbolScope) GetIndexOfSymbol(sym *SymbolTable_Symbol) int {
	i := 0
	for own := range s.ownSymbols.Values() {
		// Predicates.equalTo(sym): sym.equals(own), false for a null sym.
		if sym != nil && sym.Equals(own) {
			return i
		}
		i++
	}
	return -1
}

// Java: SymbolTable.SymbolScope.getRootNode()
func (s *SymbolTable_SymbolScope) GetRootNode() *rhino.Node {
	return s.rootNode
}

// Java: SymbolTable.SymbolScope.getParentScope()
func (s *SymbolTable_SymbolScope) GetParentScope() *SymbolTable_SymbolScope {
	return s.parent
}

// GetQualifiedSlot gets the slot for a fully-qualified name (e.g., "a.b.c")
// by trying to find property scopes at each part of the path.
//
// Java: SymbolTable.SymbolScope.getQualifiedSlot(String)
func (s *SymbolTable_SymbolScope) GetQualifiedSlot(name string) *SymbolTable_Symbol {
	fullyNamedSym := s.GetSlot(name)
	if fullyNamedSym != nil {
		return fullyNamedSym
	}

	// Go-only: byte index of the ASCII '.' (only used to cut the string).
	dot := strings.LastIndexByte(name, '.')
	if dot != -1 {
		owner := s.GetQualifiedSlot(name[:dot])
		if owner != nil && owner.GetPropertyScope() != nil {
			return owner.GetPropertyScope().GetSlot(name[dot+1:])
		}
	}

	return nil
}

// Java: SymbolTable.SymbolScope.getSlot(String)
func (s *SymbolTable_SymbolScope) GetSlot(name string) *SymbolTable_Symbol {
	own := s.GetOwnSlot(name)
	if own != nil {
		return own
	}

	var ancestor *SymbolTable_Symbol
	if s.parent != nil {
		ancestor = s.parent.GetSlot(name)
	}
	if ancestor != nil {
		return ancestor
	}
	return nil
}

// Java: SymbolTable.SymbolScope.getOwnSlot(String)
func (s *SymbolTable_SymbolScope) GetOwnSlot(name string) *SymbolTable_Symbol {
	return s.ownSymbols.Get(name)
}

// Java: SymbolTable.SymbolScope.getTypeOfThis()
func (s *SymbolTable_SymbolScope) GetTypeOfThis() rhino.JSType {
	return s.typeOfThis
}

// Java: SymbolTable.SymbolScope.isGlobalScope()
func (s *SymbolTable_SymbolScope) IsGlobalScope() bool {
	return s.GetParentScope() == nil && s.GetRootNode() != nil
}

// Java: SymbolTable.SymbolScope.isModuleScope()
func (s *SymbolTable_SymbolScope) IsModuleScope() bool {
	return s.GetRootNode() != nil && s.GetRootNode().IsModuleBody()
}

// IsDocScope returns whether this is a doc scope. A doc scope is a table for
// symbols that are documented solely within a JSDoc comment.
//
// Java: SymbolTable.SymbolScope.isDocScope()
func (s *SymbolTable_SymbolScope) IsDocScope() bool {
	return s.GetRootNode() == nil && s.mySymbol != nil && s.mySymbol.docScope == s
}

// Java: SymbolTable.SymbolScope.isPropertyScope()
func (s *SymbolTable_SymbolScope) IsPropertyScope() bool {
	return s.GetRootNode() == nil && !s.IsDocScope()
}

// Java: SymbolTable.SymbolScope.isLexicalScope()
func (s *SymbolTable_SymbolScope) IsLexicalScope() bool {
	return s.GetRootNode() != nil
}

// Java: SymbolTable.SymbolScope.isBlockScope()
func (s *SymbolTable_SymbolScope) IsBlockScope() bool {
	return s.GetRootNode() != nil && NodeUtil_CreatesBlockScope(s.GetRootNode())
}

// Java: SymbolTable.SymbolScope.getScopeDepth()
func (s *SymbolTable_SymbolScope) GetScopeDepth() int {
	return s.scopeDepth
}

// Java: SymbolTable.SymbolScope.getIndexInParent()
func (s *SymbolTable_SymbolScope) GetIndexInParent() int {
	return s.indexInParent
}

// Java: SymbolTable.SymbolScope.toString()
func (s *SymbolTable_SymbolScope) String() string {
	n := s.GetRootNode()
	if s.IsGlobalScope() {
		return "GlobalScope"
	} else if n != nil {
		type_ := "Scope"
		if s.IsModuleScope() {
			type_ = "ModuleScope"
		} else if s.IsBlockScope() {
			type_ = "BlockScope"
		} else if s.IsLexicalScope() {
			type_ = "LexicalScope"
		}
		return type_ + "@" + symbolTable_nullableString(n.GetSourceFileName()) + ":" + strconv.Itoa(n.GetLineno())
	} else {
		sym := "null"
		if s.GetSymbolForScope() != nil {
			sym = s.GetSymbolForScope().String()
		}
		return "PropertyScope@" + sym
	}
}

// SymbolTable_PropertyRefCollector collects property references.
//
// Java: private class SymbolTable.PropertyRefCollector extends AbstractPostOrderCallback implements CompilerPass
type SymbolTable_PropertyRefCollector struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *SymbolTable
}

// Java: SymbolTable.PropertyRefCollector()
func NewSymbolTable_PropertyRefCollector(outer *SymbolTable) *SymbolTable_PropertyRefCollector {
	return &SymbolTable_PropertyRefCollector{outer: outer}
}

// Java: SymbolTable.PropertyRefCollector.process(Node, Node)
func (c *SymbolTable_PropertyRefCollector) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_TraverseRoots(c.outer.compiler, c, externs, root)
}

// Java: SymbolTable.PropertyRefCollector.maybeDefineReference(Node, String, Symbol)
func (c *SymbolTable_PropertyRefCollector) maybeDefineReference(n *rhino.Node, propName string, ownerSymbol *SymbolTable_Symbol) bool {
	if ownerSymbol == nil {
		return false
	}
	// If current symbol is ObjectType - try using ObjectType.getPropertyNode to find a node.
	// That function searches for property in all extended classes and interfaces, which supports
	// multiple extended interfaces.
	owner := ownerSymbol.GetType()
	if !jutil.IsNil(owner) && owner.IsObjectType() {
		propNode := owner.ToObjectType().GetPropertyNode(propName)
		propSymbol, _ := c.outer.symbols.Get(propNode, propName)
		if propSymbol != nil {
			propSymbol.DefineReferenceAt(n)
			return true
		}
	}
	// Use PropertyScope to find property. PropertyScope doesn't support multiple parents so all
	// cases where multiple parents can occur (e.g. class implementing 2+ interfaces) must be
	// handled above.
	if ownerSymbol.GetPropertyScope() != nil {
		prop := ownerSymbol.GetPropertyScope().GetSlot(propName)
		if prop != nil {
			prop.DefineReferenceAt(n)
			return true
		}
	}
	return false
}

// tryDefineLexicalQualifiedNameRef tries to find the symbol by its fully
// qualified name.
//
// Java: SymbolTable.PropertyRefCollector.tryDefineLexicalQualifiedNameRef(String, Node)
func (c *SymbolTable_PropertyRefCollector) tryDefineLexicalQualifiedNameRef(name *string, n *rhino.Node) bool {
	if name != nil {
		lexicalSym := c.outer.GetEnclosingScope(n).GetQualifiedSlot(*name)
		if lexicalSym != nil {
			lexicalSym.DefineReferenceAt(n)
			return true
		}
	}
	return false
}

// tryRemoveLexicalQualifiedNameRef tries to remove a reference by its fully
// qualified name. If the symbol has no references left, remove it completely.
//
// Java: SymbolTable.PropertyRefCollector.tryRemoveLexicalQualifiedNameRef(String, Node)
func (c *SymbolTable_PropertyRefCollector) tryRemoveLexicalQualifiedNameRef(name *string, n *rhino.Node) {
	if name != nil {
		lexicalSym := c.outer.GetEnclosingScope(n).GetQualifiedSlot(*name)
		if lexicalSym != nil &&
			lexicalSym.IsLexicalVariable() &&
			lexicalSym.GetDeclaration().GetNode() == n {
			c.outer.removeSymbol(lexicalSym)
		}
	}
}

// Java: SymbolTable.PropertyRefCollector.maybeDefineTypedReference(Node, String, JSType)
func (c *SymbolTable_PropertyRefCollector) maybeDefineTypedReference(n *rhino.Node, propName string, owner rhino.JSType) bool {
	t := c.outer
	if owner.IsGlobalThisType() {
		sym := t.globalScope.GetSlot(propName)
		if sym != nil {
			sym.DefineReferenceAt(n)
			return true
		}
	} else if owner.IsNominalConstructorOrInterface() {
		return c.maybeDefineReference(n, propName, t.GetSymbolDeclaredBy(owner.ToMaybeFunctionType()))
	} else if owner.IsEnumType() {
		return c.maybeDefineReference(n, propName, t.GetSymbolDeclaredByEnum(owner.ToMaybeEnumType()))
	} else {
		defined := false
		for _, ctor := range t.GetAllSymbolsForType(owner) {
			if ctor.GetPropertyScope() != nil {
				ctorType := ctor.GetType()
				isConstructor :=
					!jutil.IsNil(ctorType) && (ctorType.IsConstructor() || ctorType.IsInterface())
				if !isConstructor && c.maybeDefineReference(n, propName, ctor) {
					defined = true
				}
			}
			if c.maybeDefineReference(n, propName, t.GetSymbolForInstancesOf(ctor)) {
				defined = true
			}
		}
		if defined {
			return true
		}
	}
	// Handle cases like accessing properties on a imported module.
	// const foo = goog.require('some.module.foo');
	// foo.doSomething();
	// here n is foo, propName is doSomething and owner is module type.
	if owner.IsObjectType() && !jutil.IsNil(owner.ToMaybeObjectType()) {
		propNodeDecl := owner.AssertObjectType().GetPropertyNode(propName)
		if propNodeDecl == nil {
			return false
		}
		// There is no handy way to find symbol object the given property node. So we do
		// property node => namespace node => namespace symbol => property symbol.
		if propNodeDecl != nil &&
			(propNodeDecl.IsGetProp() || propNodeDecl.IsOptChainGetProp()) {
			namespace := propNodeDecl.GetFirstChild()
			namespaceSym := t.GetSymbolForName(namespace, namespace.GetQualifiedName())
			if namespaceSym != nil && namespaceSym.GetPropertyScope() != nil {
				propSym := namespaceSym.GetPropertyScope().GetSlot(propName)
				if propSym != nil && propSym.GetDeclarationNode() != n {
					propSym.DefineReferenceAt(n)
					return true
				}
			}
		}

		// Here handle the case of goog.module exports that have shape of:
		//
		// goog.module('foo.bar');
		// function doOne() { ... }
		// exports = {doOne};
		//
		// propNodeDecl is `doOne` node on the last line with exports. From that node and propName
		// `doOne` we need to get symbol for the `function doOne`. To achieve that we get scope of
		// `doOne` node which should be whole module scope and look for symbol with corresponding
		// name. This relies on the assumption that function name in module scope exported name. But
		// if export is renamed: `exports = {doOneRenamed: doOne}`, then this approach will fail.
		propNodeDeclScope := t.GetEnclosingScope(propNodeDecl)
		if propNodeDeclScope == nil || !propNodeDeclScope.IsModuleScope() {
			return false
		}
		symbol := t.GetSymbolForName(propNodeDecl, &propName)
		if symbol != nil && symbol.GetDeclarationNode() != n {
			symbol.DefineReferenceAt(n)
			return true
		}
	}
	return false
}

// Java: SymbolTable.PropertyRefCollector.visit(NodeTraversal, Node, Node)
func (c *SymbolTable_PropertyRefCollector) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	// There are two ways to define a property reference:
	// 1) As a fully qualified lexical symbol (e.g., x.y)
	// 2) As a property of another object (e.g., x's y)
	// Property definitions should take precedence over lexical
	// definitions. e.g., for "a.b", it's more useful to record
	// this as "property b of the type of a", than as "symbol a.b".

	if n.IsGetProp() || n.IsOptChainGetProp() {
		owner := n.GetFirstChild().GetJSType()
		if !jutil.IsNil(owner) {
			defined := c.maybeDefineTypedReference(n, n.GetString(), owner)

			if defined {
				c.tryRemoveLexicalQualifiedNameRef(n.GetQualifiedName(), n)
				return
			}
		}

		c.tryDefineLexicalQualifiedNameRef(n.GetQualifiedName(), n)
	} else if n.IsStringKey() {
		owner := parent.GetJSType()
		if !jutil.IsNil(owner) {
			defined := c.maybeDefineTypedReference(n, n.GetString(), owner)

			if defined {
				c.tryRemoveLexicalQualifiedNameRef(NodeUtil_GetBestLValueName(n), n)
				return
			}
		}

		c.tryDefineLexicalQualifiedNameRef(NodeUtil_GetBestLValueName(n), n)
	}
}

// SymbolTable_ThisRefCollector collects references to "this".
//
// Java: private class SymbolTable.ThisRefCollector extends NodeTraversal.AbstractScopedCallback implements CompilerPass
type SymbolTable_ThisRefCollector struct {
	NodeTraversal_AbstractScopedCallbackBase
	outer *SymbolTable

	// The 'this' symbols in the current scope chain.
	//
	// If we don't know how to declare 'this' in a scope chain,
	// then null should be on the stack. But this should be a rare
	// occurrence. We should strive to always be able to come up
	// with some symbol for 'this'.
	//
	// This list only has entries for function scopes and the global scope, and doesn't store a
	// separate `this` value for other block scopes.
	thisStack []*SymbolTable_Symbol
}

// Java: SymbolTable.ThisRefCollector()
func NewSymbolTable_ThisRefCollector(outer *SymbolTable) *SymbolTable_ThisRefCollector {
	return &SymbolTable_ThisRefCollector{outer: outer}
}

// Java: SymbolTable.ThisRefCollector.process(Node, Node)
func (c *SymbolTable_ThisRefCollector) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_TraverseRoots(c.outer.compiler, c, externs, root)
}

// Java: SymbolTable.ThisRefCollector.enterScope(NodeTraversal)
func (c *SymbolTable_ThisRefCollector) EnterScope(t *NodeTraversal) {
	st := c.outer
	var symbol *SymbolTable_Symbol
	if t.InGlobalScope() {
		// Declare the global this at the first input root.
		// This is a bizarre place to put it, but we need some
		// location with a real file path (because all symbols
		// must have a path).
		// Note that root.lastChild.firstChild is the first non-extern input.
		firstInputRoot := t.GetScopeRoot().GetLastChild().GetFirstChild()
		if firstInputRoot != nil {
			symbol =
				st.addSymbol(
					SymbolTable_GLOBAL_THIS,
					st.registry.GetNativeType(rhino.JSTypeNative_GLOBAL_THIS),
					false, /* declared */
					st.globalScope,
					firstInputRoot)
			symbol.SetDeclaration(NewSymbolTable_Reference(symbol, firstInputRoot))
		}
		c.thisStack = append(c.thisStack, symbol)
	} else if t.GetScopeRoot().IsFunction() {
		// Otherwise, declare a "this" property when possible.
		scopeRoot := t.GetScopeRoot()
		scope := st.scopes.Get(scopeRoot)
		if NodeUtil_GetFunctionBody(scopeRoot).HasChildren() {
			scopeSymbol := st.GetSymbolForScope(scope)
			if scopeSymbol != nil {
				propScope := scopeSymbol.GetPropertyScope()
				if propScope != nil {
					// If a function is assigned multiple times, we only want
					// one addressable "this" symbol.
					symbol = propScope.GetOwnSlot("this")
					if symbol == nil {
						rootType := t.GetScopeRoot().GetJSType()
						var fnType rhino.FunctionType
						if !jutil.IsNil(rootType) {
							fnType = rootType.ToMaybeFunctionType()
						}
						var type_ rhino.JSType
						if !jutil.IsNil(fnType) {
							type_ = fnType.GetTypeOfThis()
						}
						symbol = st.addSymbol("this", type_ /* inferred= */, false, scope, scopeRoot)
					}
				}
			}
		} else {
			symbolTable_logger.Fine("Skipping empty function: " + scopeRoot.String())
		}
		c.thisStack = append(c.thisStack, symbol)
	}
	// Don't add to the `thisStack` for other block scopes.
}

// Java: SymbolTable.ThisRefCollector.exitScope(NodeTraversal)
func (c *SymbolTable_ThisRefCollector) ExitScope(t *NodeTraversal) {
	if t.InGlobalScope() || t.GetScopeRoot().IsFunction() {
		c.thisStack = c.thisStack[:len(c.thisStack)-1]
	}
}

// Java: SymbolTable.ThisRefCollector.visit(NodeTraversal, Node, Node)
func (c *SymbolTable_ThisRefCollector) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if !n.IsThis() {
		return
	}

	if len(c.thisStack) == 0 {
		// Java: Iterables.getLast of an empty list.
		panic(jutil.NewNoSuchElementException())
	}
	symbol := c.thisStack[len(c.thisStack)-1]
	if symbol != nil {
		ref := symbol.DefineReferenceAt(n)
		if symbol.GetDeclaration() == nil {
			symbol.SetDeclaration(ref)
		}
	}
}

// SymbolTable_JSDocInfoCollector collects references to types in JSDocInfo.
//
// Java: private class SymbolTable.JSDocInfoCollector extends AbstractPostOrderCallback
type SymbolTable_JSDocInfoCollector struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer        *SymbolTable
	typeRegistry *rhino.JSTypeRegistry
}

// Java: SymbolTable.JSDocInfoCollector(JSTypeRegistry)
func NewSymbolTable_JSDocInfoCollector(outer *SymbolTable, registry *rhino.JSTypeRegistry) *SymbolTable_JSDocInfoCollector {
	return &SymbolTable_JSDocInfoCollector{outer: outer, typeRegistry: registry}
}

// Java: SymbolTable.JSDocInfoCollector.visit(NodeTraversal, Node, Node)
func (c *SymbolTable_JSDocInfoCollector) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	info := n.GetJSDocInfo()
	if info == nil {
		return
	}
	c.outer.docInfos.Add(n)

	for it := info.TypeExprIter(); it.Next(); {
		scope := c.outer.scopes.Get(t.GetScopeRoot())
		if scope == nil {
			scope = c.outer.globalScope
		}
		c.VisitTypeNode(info.GetTemplateTypes(), scope, it.Node())
	}
}

// Java: SymbolTable.JSDocInfoCollector.isNativeSourcelessType(String)
func (c *SymbolTable_JSDocInfoCollector) isNativeSourcelessType(name string) bool {
	switch name {
	case "null", "undefined", "void":
		return true
	default:
		return false
	}
}

// Java: SymbolTable.JSDocInfoCollector.visitTypeNode(ImmutableMap<String, JSTypeExpression>, SymbolScope, Node)
func (c *SymbolTable_JSDocInfoCollector) VisitTypeNode(
	templateTypeNames *jutil.LinkedHashMap[string, *rhino.JSTypeExpression], scope *SymbolTable_SymbolScope, n *rhino.Node) {
	if n.IsStringLit() &&
		!c.isNativeSourcelessType(n.GetString()) &&
		!templateTypeNames.ContainsKey(n.GetString()) {
		symbol := c.lookupPossiblyDottedName(scope, n.GetString())
		if symbol != nil {
			ref := n
			typeString := n.GetString()
			if n.GetOriginalName() != nil {
				typeString = *n.GetOriginalName()
			}
			// Qualified names in JSDoc types are kept as a single string: "foo.bar.MyType". In order
			// to have good indexing we need to make the SymbolTable reference include only "MyType"
			// instead of "foo.bar.MyType". To do that we clone the type node and change the source
			// info of the clone to include only the last part of the type ("MyType").
			if strings.Contains(typeString, ".") {
				lastPart := typeString[strings.LastIndexByte(typeString, '.')+1:]
				copy_ := n.CloneNode()
				copy_.SetLinenoCharno(
					copy_.GetLineno(), copy_.GetCharno()+copy_.GetLength()-jutil.StringLength(lastPart))
				copy_.SetLength(jutil.StringLength(lastPart))
				ref = copy_
			}
			symbol.DefineReferenceAt(ref)
		}
	}

	for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
		c.VisitTypeNode(templateTypeNames, scope, child)
	}
}

// Java: SymbolTable.JSDocInfoCollector.lookupPossiblyDottedName(SymbolScope, String)
func (c *SymbolTable_JSDocInfoCollector) lookupPossiblyDottedName(scope *SymbolTable_SymbolScope, dottedName string) *SymbolTable_Symbol {
	// Try the dotted name to start.
	result := scope.GetQualifiedSlot(dottedName)
	if result != nil {
		return result
	}
	// If we can't find this type, it might be a reference to a
	// primitive type (like {string}). Autobox it to check.
	type_ := c.typeRegistry.GetGlobalType(dottedName)
	var autobox rhino.JSType
	if !jutil.IsNil(type_) {
		autobox = type_.AutoboxesTo()
	}
	if jutil.IsNil(autobox) {
		return c.outer.GetSymbolForTypeHelper(type_, true)
	}
	return c.outer.GetSymbolForTypeHelper(autobox, true)
}

// SymbolTable_VisibilityCollector collects the visibility information for
// each name/property.
//
// Java: private class SymbolTable.VisibilityCollector extends AbstractPostOrderCallback
type SymbolTable_VisibilityCollector struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer             *SymbolTable
	fileVisibilityMap *jutil.LinkedHashMap[rhino.StaticSourceFile, rhino.JSDocInfo_Visibility]
}

// Java: SymbolTable.VisibilityCollector(ImmutableMap<StaticSourceFile, Visibility>)
func NewSymbolTable_VisibilityCollector(
	outer *SymbolTable, fileVisibilityMap *jutil.LinkedHashMap[rhino.StaticSourceFile, rhino.JSDocInfo_Visibility]) *SymbolTable_VisibilityCollector {
	return &SymbolTable_VisibilityCollector{outer: outer, fileVisibilityMap: fileVisibilityMap}
}

// Java: SymbolTable.VisibilityCollector.visit(NodeTraversal, Node, Node)
func (c *SymbolTable_VisibilityCollector) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsName() {
		c.visitName(t, n)
	} else if n.IsGetProp() || n.IsOptChainGetProp() {
		c.visitProperty(n, parent)
	}
}

// Java: SymbolTable.VisibilityCollector.visitName(NodeTraversal, Node)
func (c *SymbolTable_VisibilityCollector) visitName(t *NodeTraversal, n *rhino.Node) {
	symbol, _ := c.outer.symbols.Get(n, n.GetString())
	if symbol == nil {
		return
	}
	// Visibility already set.
	if symbol.GetVisibility() != nil {
		return
	}
	var_ := t.GetScope().GetVar(n.GetString())
	if var_ == nil {
		return
	}
	// Go-only: AccessControlUtils.getEffectiveNameVisibility never returns
	// null here (see AccessControlUtils), so Java's null check is dropped.
	v := AccessControlUtils_GetEffectiveNameVisibility(n, var_, c.fileVisibilityMap)
	symbol.SetVisibility(v)
}

// Java: SymbolTable.VisibilityCollector.visitProperty(Node, Node)
func (c *SymbolTable_VisibilityCollector) visitProperty(getprop *rhino.Node, parent *rhino.Node) {
	propertyName := getprop.GetString()
	symbol, _ := c.outer.symbols.Get(getprop, propertyName)
	if symbol == nil {
		return
	}
	// Visibility already set.
	if symbol.GetVisibility() != nil {
		return
	}
	jsType := getprop.GetFirstChild().GetJSType()
	if jutil.IsNil(jsType) {
		return
	}
	isOverride :=
		parent.GetJSDocInfo() != nil && parent.IsAssign() && parent.GetFirstChild() == getprop
	if isOverride {
		// Don't bother with AccessControlUtils for overridden properties.
		// AccessControlUtils currently has complicated logic for detecting
		// visibility mismatches for overridden properties that is still
		// too tightly coupled to CheckAccessControls. TODO(brndn): simplify.
		symbol.SetVisibility(rhino.JSDocInfo_Visibility_INHERITED)
	} else {
		referenceType := rhino.ObjectType_Cast(jsType.Dereference())
		// Go-only: AccessControlUtils.getEffectivePropertyVisibility never
		// returns null (see AccessControlUtils), so Java's null check is dropped.
		v :=
			AccessControlUtils_GetEffectivePropertyVisibility(
				getprop, referenceType, c.fileVisibilityMap)
		symbol.SetVisibility(v)
	}
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

	t.nodeOrdering = func(a, b *rhino.Node) int {
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
		return jutil.StringCompareTo(a.GetName(), b.GetName())
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
	if sym.GetJSDocInfo() != nil && sym.GetJSDocInfo().HasTypedefType() {
		type_ := t.registry.GetGlobalType(sym.GetName())
		if !jutil.IsNil(type_) {
			return type_
		}
		return sym.GetJSDocInfo().GetTypedefType().Evaluate(symbolTable_topScope(t.compiler), t.registry)
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
