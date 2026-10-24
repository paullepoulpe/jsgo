// Port of com/google/javascript/rhino/jstype/SimpleReference.java @ closure-compiler v20261015

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

import "github.com/paullepoulpe/jsgo/internal/jutil"

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
