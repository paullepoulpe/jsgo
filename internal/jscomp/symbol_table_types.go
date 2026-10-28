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
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// SymbolTable_Symbol is a symbol-table entry.
//
// Java: public static final class SymbolTable.Symbol extends SimpleSlot
type SymbolTable_Symbol struct {
	rhino.SimpleSlotBase

	// Use a linked hash map, so that the results are deterministic
	// (and so the declaration always comes first).
	references *jutil.LinkedHashMap[rhino.Node, *SymbolTable_Reference]

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
		references: jutil.NewLinkedHashMap[rhino.Node, *SymbolTable_Reference](),
		scope:      scope,
	}
	s.InitSimpleSlotBase(name, type_, inferred)
	return s
}

func jsTypeEquals(a, b rhino.JSType) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Equals(b)
}

// Java: SymbolTable.Symbol.equals(Object)
func (s *SymbolTable_Symbol) Equals(o any) bool {
	other, ok := o.(*SymbolTable_Symbol)
	if !ok || other == nil {
		return false
	}

	return s.IsTypeInferred() == other.IsTypeInferred() &&
		s.GetName() == other.GetName() &&
		jsTypeEquals(s.GetType(), other.GetType()) &&
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
func (s *SymbolTable_Symbol) DefineReferenceAt(n rhino.Node) *SymbolTable_Reference {
	return s.references.ComputeIfAbsent(n, func(k rhino.Node) *SymbolTable_Reference { return NewSymbolTable_Reference(s, k) })
}

// SetDeclaration sets the declaration node. May only be called once.
//
// Java: SymbolTable.Symbol.setDeclaration(Reference)
func (s *SymbolTable_Symbol) SetDeclaration(ref *SymbolTable_Reference) {
	jutil.CheckState(s.declaration == nil)
	s.declaration = ref
}

// Java: SymbolTable.Symbol.getDeclarationNode()
func (s *SymbolTable_Symbol) GetDeclarationNode() rhino.Node {
	if s.declaration == nil {
		return rhino.NilNode
	}
	return s.declaration.GetNode()
}

// Java: SymbolTable.Symbol.getSourceFileName()
func (s *SymbolTable_Symbol) GetSourceFileName() *string {
	n := s.GetDeclarationNode()
	if n.IsNil() {
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
	if !n.IsNil() {
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
func NewSymbolTable_Reference(symbol *SymbolTable_Symbol, node rhino.Node) *SymbolTable_Reference {
	r := &SymbolTable_Reference{}
	r.InitSimpleReference(symbol, node)
	return r
}

// SymbolTable_SymbolScope is the scope of a symbol.
//
// Java: public static final class SymbolTable.SymbolScope
type SymbolTable_SymbolScope struct {
	rootNode   rhino.Node
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
	rootNode rhino.Node,
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
func (s *SymbolTable_SymbolScope) GetRootNode() rhino.Node {
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
	return s.GetParentScope() == nil && !s.GetRootNode().IsNil()
}

// Java: SymbolTable.SymbolScope.isModuleScope()
func (s *SymbolTable_SymbolScope) IsModuleScope() bool {
	return !s.GetRootNode().IsNil() && s.GetRootNode().IsModuleBody()
}

// IsDocScope returns whether this is a doc scope. A doc scope is a table for
// symbols that are documented solely within a JSDoc comment.
//
// Java: SymbolTable.SymbolScope.isDocScope()
func (s *SymbolTable_SymbolScope) IsDocScope() bool {
	return s.GetRootNode().IsNil() && s.mySymbol != nil && s.mySymbol.docScope == s
}

// Java: SymbolTable.SymbolScope.isPropertyScope()
func (s *SymbolTable_SymbolScope) IsPropertyScope() bool {
	return s.GetRootNode().IsNil() && !s.IsDocScope()
}

// Java: SymbolTable.SymbolScope.isLexicalScope()
func (s *SymbolTable_SymbolScope) IsLexicalScope() bool {
	return !s.GetRootNode().IsNil()
}

// Java: SymbolTable.SymbolScope.isBlockScope()
func (s *SymbolTable_SymbolScope) IsBlockScope() bool {
	return !s.GetRootNode().IsNil() && NodeUtil_CreatesBlockScope(s.GetRootNode())
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
	} else if !n.IsNil() {
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
