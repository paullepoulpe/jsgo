// Port of com/google/javascript/rhino/StaticScope.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/StaticSlot.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/StaticRef.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/StaticSymbolTable.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/StaticTypedScope.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/StaticTypedSlot.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/StaticTypedRef.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/SimpleSlot.java @ closure-compiler v20261021
// Port of com/google/javascript/rhino/jstype/SimpleReference.java @ closure-compiler v20261021

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Nick Santos
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Nick Santos
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"iter"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// StaticScope must be implemented by any object that defines variables for
// the purposes of static analysis. It is distinguished from the Scriptable
// class that Rhino normally uses to represent a run-time scope.
type StaticScope interface {
	// GetRootNode returns the root node associated with this scope. May be
	// null.
	GetRootNode() Node

	// GetParentScope returns the scope enclosing this one or null if none.
	GetParentScope() StaticScope

	// GetSlot returns any defined slot within this scope for this name. This
	// call continues searching through parent scopes if a slot with this name
	// is not found in the current scope. Returns null if no definition
	// exists.
	GetSlot(name string) StaticSlot

	// GetOwnSlot is like getSlot but does not recurse into parent scopes.
	GetOwnSlot(name string) StaticSlot

	// GetTopmostScopeOfEventualDeclaration returns the topmost slot containing
	// this name, or null if no slots do. May return true in cases where
	// getSlot returns null, though. Do not rely on this method if you need an
	// actual slot.
	//
	// This method is intended for use while scopes are still being built,
	// hence the name 'eventual' declaration. Once scope building is complete,
	// the scope returned from this method must be equivalent to
	// "getSlot(name).getScope()" or null.
	//
	// A Java default method: implementations that do not override it call
	// StaticScope_GetTopmostScopeOfEventualDeclaration.
	GetTopmostScopeOfEventualDeclaration(name string) StaticScope
}

// StaticScope_GetTopmostScopeOfEventualDeclaration is the default body of
// StaticScope.getTopmostScopeOfEventualDeclaration.
//
// Java: StaticScope.getTopmostScopeOfEventualDeclaration(String) (default method)
func StaticScope_GetTopmostScopeOfEventualDeclaration(s StaticScope, name string) StaticScope {
	slot := s.GetOwnSlot(name)
	if slot != nil {
		return slot.GetScope()
	}
	if s.GetParentScope() != nil {
		return s.GetParentScope().GetTopmostScopeOfEventualDeclaration(name)
	}
	return nil
}

// StaticSlot must be implemented by variables that can appear as members of
// a StaticScope.
type StaticSlot interface {
	// GetName gets the name of the slot.
	GetName() string

	// GetDeclaration gets the declaration of this symbol. May not exist.
	GetDeclaration() StaticRef

	// GetJSDocInfo gets the JSDoc for this slot.
	GetJSDocInfo() *JSDocInfo

	GetScope() StaticScope
}

// StaticRef tells us all the ways that a StaticSlot is used in a program.
type StaticRef interface {
	// GetSymbol returns the variable that this reference points to.
	GetSymbol() StaticSlot

	// GetNode returns the node where the reference lives.
	GetNode() Node

	// GetSourceFile returns the source file where the reference lives.
	GetSourceFile() StaticSourceFile
}

// StaticSymbolTable looks up references by the symbols that they refer to.
type StaticSymbolTable[S StaticSlot, R StaticRef] interface {
	// GetReferences returns the references that point to the given symbol.
	GetReferences(symbol S) iter.Seq[R]

	// GetScope returns the scope for a given symbol.
	GetScope(symbol S) StaticScope

	// GetAllSymbols returns all variables in this symbol table.
	GetAllSymbols() iter.Seq[S]
}

// StaticTypedScope must be implemented by any object that defines variables
// for the purposes of static analysis. It is distinguished from the
// Scriptable class that Rhino normally uses to represent a run-time scope.
//
// Go-only: language-adapter — Java narrows the return types of
// getParentScope() (to StaticTypedScope), getSlot(String) and
// getOwnSlot(String) (to StaticTypedSlot). Go interfaces have no covariant
// overrides, so those methods keep StaticScope's signatures; every
// StaticTypedScope implementation returns StaticTypedScope / StaticTypedSlot
// values from them, and callers that rely on the narrower Java type
// type-assert (`slot.(StaticTypedSlot)`).
type StaticTypedScope interface {
	StaticScope

	// GetTypeOfThis returns the expected type of `this` in the current scope.
	GetTypeOfThis() JSType

	// LookupQualifiedName looks up a given qualified name in the scope. if
	// that fails, looks up the component properties off of any owner types
	// that are in scope.
	//
	// This is always more or equally expensive as calling getSlot(String
	// name), so should only be used when necessary.
	//
	// This only returns declared qualified names and known properties. It
	// returns null given an inferred name.
	//
	// A Java default method: implementations that do not override it call
	// StaticTypedScope_LookupQualifiedName.
	LookupQualifiedName(qname QualifiedName) JSType
}

// StaticTypedScope_LookupQualifiedName is the default body of
// StaticTypedScope.lookupQualifiedName.
//
// Java: StaticTypedScope.lookupQualifiedName(QualifiedName) (default method)
func StaticTypedScope_LookupQualifiedName(s StaticTypedScope, qname QualifiedName) JSType {
	switch q := qname.(type) {
	case *QualifiedName_NodeQname:
		if name, ok := q.node.getQualifiedNameString(); ok {
			return StaticTypedScope_LookupQualifiedNameString(s, name)
		}
		return StaticTypedScope_LookupQualifiedNameString(s, "null")
	case *QualifiedName_StringListQname:
		return StaticTypedScope_LookupQualifiedNameString(s, q.Join())
	}
	slot, _ := s.GetSlot(qname.Join()).(StaticTypedSlot)
	if slot != nil && !slot.IsTypeInferred() {
		return slot.GetType()
	} else if !qname.IsSimple() {
		type_ := s.LookupQualifiedName(qname.GetOwner())
		if type_ != nil && !type_.IsUnknownType() {
			return type_.FindPropertyType(qname.GetComponent())
		}
	}
	return nil
}

// StaticTypedScope_LookupQualifiedNameString is StaticTypedScope_LookupQualifiedName
// on a dot-separated string without allocating QualifiedName wrappers or prefix strings.
func StaticTypedScope_LookupQualifiedNameString(s StaticTypedScope, name string) JSType {
	slot, _ := s.GetSlot(name).(StaticTypedSlot)
	if slot != nil && !slot.IsTypeInferred() {
		return slot.GetType()
	}
	if dot := strings.LastIndexByte(name, '.'); dot >= 0 {
		type_ := StaticTypedScope_LookupQualifiedNameString(s, name[:dot])
		if type_ != nil && !type_.IsUnknownType() {
			return type_.FindPropertyType(name[dot+1:])
		}
	}
	return nil
}

// StaticTypedScope_LookupQualifiedNameNode looks up n's qualified name in s
// without allocating a QualifiedName_NodeQname wrapper.
func StaticTypedScope_LookupQualifiedNameNode(s StaticTypedScope, n Node) JSType {
	if name, ok := n.getQualifiedNameString(); ok {
		return StaticTypedScope_LookupQualifiedNameString(s, name)
	}
	return StaticTypedScope_LookupQualifiedNameString(s, "null")
}

// StaticTypedSlot must be implemented by variables that can appear as
// members of a StaticTypedScope.
//
// Go-only: language-adapter — getDeclaration() (StaticTypedRef in Java) and
// getScope() (StaticTypedScope in Java) keep StaticSlot's signatures; see
// StaticTypedScope.
type StaticTypedSlot interface {
	StaticSlot

	// GetType returns the type information, if any, for this slot, or null
	// if no type is declared for it.
	GetType() JSType

	// IsTypeInferred returns whether the type has been inferred (as opposed
	// to declared).
	IsTypeInferred() bool
}

// StaticTypedRef tells us all the ways that a StaticTypedSlot is used in a
// program.
//
// Go-only: language-adapter — getSymbol() (StaticTypedSlot in Java) keeps
// StaticRef's signature; see StaticTypedScope.
type StaticTypedRef = StaticRef

// SimpleSlot is the minimum implementation of StaticTypedSlot.
//
// Java: class SimpleSlot (subclassed by jscomp's SymbolTable.Symbol).
type SimpleSlot = StaticTypedSlot

// SimpleSlotBase holds SimpleSlot's fields and methods; plain slots are
// *SimpleSlotBase values. The base methods make no virtual calls, so no self
// field is needed.
type SimpleSlotBase struct {
	name     string
	type_    JSType
	inferred bool
}

// NewSimpleSlot constructs a plain SimpleSlot.
//
// Java: SimpleSlot(String, JSType, boolean)
func NewSimpleSlot(name string, type_ JSType, inferred bool) *SimpleSlotBase {
	s := &SimpleSlotBase{}
	s.InitSimpleSlotBase(name, type_, inferred)
	return s
}

// InitSimpleSlotBase is the Go form of the constructor for subclasses.
//
// Java: SimpleSlot(String, JSType, boolean)
func (s *SimpleSlotBase) InitSimpleSlotBase(name string, type_ JSType, inferred bool) {
	s.name = name
	s.type_ = type_
	s.inferred = inferred
}

// Java: SimpleSlot.getName()
func (s *SimpleSlotBase) GetName() string {
	return s.name
}

// Java: SimpleSlot.getType()
func (s *SimpleSlotBase) GetType() JSType {
	return s.type_
}

// Java: SimpleSlot.isTypeInferred()
func (s *SimpleSlotBase) IsTypeInferred() bool {
	return s.inferred
}

// Java: SimpleSlot.getDeclaration()
func (s *SimpleSlotBase) GetDeclaration() StaticRef {
	return nil
}

// Java: SimpleSlot.getJSDocInfo()
func (s *SimpleSlotBase) GetJSDocInfo() *JSDocInfo {
	return nil
}

// Java: SimpleSlot.getScope()
func (s *SimpleSlotBase) GetScope() StaticScope {
	return nil
}

// SimpleReference is a simple immutable reference.
//
// Go-only: language-adapter — SimpleReference<T extends StaticTypedSlot> is
// a generic class whose methods make no virtual calls; subclasses (jscomp's
// SymbolTable.Reference, PreprocessorSymbolTable.Reference) embed the
// struct. getSymbol() keeps StaticRef's return type (see StaticTypedScope);
// GetTypedSymbol returns the T.
type SimpleReference[T StaticTypedSlot] struct {
	symbol T
	node   Node
}

// Java: SimpleReference(T, Node)
func NewSimpleReference[T StaticTypedSlot](symbol T, node Node) *SimpleReference[T] {
	r := &SimpleReference[T]{}
	r.InitSimpleReference(symbol, node)
	return r
}

// InitSimpleReference is the Go form of the constructor for subclasses.
//
// Java: SimpleReference(T, Node)
func (r *SimpleReference[T]) InitSimpleReference(symbol T, node Node) {
	r.symbol = symbol
	r.node = node
}

// Java: SimpleReference.getSymbol()
func (r *SimpleReference[T]) GetSymbol() StaticSlot {
	if jutil.IsNil(r.symbol) {
		return nil
	}
	return r.symbol
}

// GetTypedSymbol is getSymbol() with its Java return type T.
//
// Go-only: language-adapter — covariant return of getSymbol().
func (r *SimpleReference[T]) GetTypedSymbol() T {
	return r.symbol
}

// Java: SimpleReference.getNode()
func (r *SimpleReference[T]) GetNode() Node {
	return r.node
}

// Java: SimpleReference.getSourceFile()
func (r *SimpleReference[T]) GetSourceFile() StaticSourceFile {
	return r.node.GetStaticSourceFile()
}

// Java: SimpleReference.toString()
func (r *SimpleReference[T]) String() string {
	var sourceName *string
	if !r.node.IsNil() {
		sourceName = r.node.GetSourceFileName()
	}
	lineNo := -1
	if !r.node.IsNil() {
		lineNo = r.node.GetLineno()
	}
	return jutil.ToJavaString(r.node.GetQualifiedName()) + "@" + jutil.ToJavaString(sourceName) + ":" + jutil.ToJavaString(lineNo)
}
