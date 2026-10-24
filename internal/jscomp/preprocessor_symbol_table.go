// Port of com/google/javascript/jscomp/PreprocessorSymbolTable.java @ closure-compiler v20261021

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// PreprocessorSymbolTable is a symbol table for references that are
// removed by preprocessor passes (like ProcessClosurePrimitives).
//
// Go-only: language-adapter — getSlot and getOwnSlot narrow
// StaticScope's return type to SimpleSlot in Java; Go interface methods
// cannot, so they return the rhino.StaticSlot holding the SimpleSlot.
type PreprocessorSymbolTable struct {
	// All preprocessor symbols are globals.
	symbols *jutil.LinkedHashMap[string, rhino.SimpleSlot]

	refs *jutil.ListMultimap[string, *PreprocessorSymbolTable_Reference]

	root rhino.Node
}

// Java: PreprocessorSymbolTable(Node)
func NewPreprocessorSymbolTable(root rhino.Node) *PreprocessorSymbolTable {
	return &PreprocessorSymbolTable{
		symbols: jutil.NewLinkedHashMap[string, rhino.SimpleSlot](),
		refs:    jutil.ArrayListMultimap_Create[string, *PreprocessorSymbolTable_Reference](),
		root:    root,
	}
}

// Java: PreprocessorSymbolTable.getRootNode()
func (t *PreprocessorSymbolTable) GetRootNode() rhino.Node {
	return t.root
}

// Java: PreprocessorSymbolTable.getTypeOfThis()
func (t *PreprocessorSymbolTable) GetTypeOfThis() rhino.JSType {
	return nil
}

// Java: PreprocessorSymbolTable.getParentScope()
func (t *PreprocessorSymbolTable) GetParentScope() rhino.StaticScope {
	return nil
}

// Java: PreprocessorSymbolTable.getSlot(String)
func (t *PreprocessorSymbolTable) GetSlot(name string) rhino.StaticSlot {
	// A missing name is a nil SimpleSlot, which converts to a nil StaticSlot.
	return t.symbols.Get(name)
}

// Java: PreprocessorSymbolTable.getOwnSlot(String)
func (t *PreprocessorSymbolTable) GetOwnSlot(name string) rhino.StaticSlot {
	return t.GetSlot(name)
}

// GetTopmostScopeOfEventualDeclaration is the StaticScope default method.
//
// Java: StaticScope.getTopmostScopeOfEventualDeclaration(String) (default method)
func (t *PreprocessorSymbolTable) GetTopmostScopeOfEventualDeclaration(name string) rhino.StaticScope {
	return rhino.StaticScope_GetTopmostScopeOfEventualDeclaration(t, name)
}

// LookupQualifiedName is the StaticTypedScope default method.
//
// Java: StaticTypedScope.lookupQualifiedName(QualifiedName) (default method)
func (t *PreprocessorSymbolTable) LookupQualifiedName(qname rhino.QualifiedName) rhino.JSType {
	return rhino.StaticTypedScope_LookupQualifiedName(t, qname)
}

// Java: PreprocessorSymbolTable.getReferences(SimpleSlot)
func (t *PreprocessorSymbolTable) GetReferences(symbol rhino.SimpleSlot) iter.Seq[*PreprocessorSymbolTable_Reference] {
	// Collections.unmodifiableCollection(refs.get(symbol.getName()))
	return slices.Values(t.refs.Get(symbol.GetName()))
}

// Java: PreprocessorSymbolTable.getAllSymbols()
func (t *PreprocessorSymbolTable) GetAllSymbols() iter.Seq[rhino.SimpleSlot] {
	// Collections.unmodifiableCollection(symbols.values())
	return t.symbols.Values()
}

// Java: PreprocessorSymbolTable.getScope(SimpleSlot)
func (t *PreprocessorSymbolTable) GetScope(slot rhino.SimpleSlot) rhino.StaticScope {
	return t
}

// Java: PreprocessorSymbolTable.addReference(Node)
func (t *PreprocessorSymbolTable) AddReference(node rhino.Node) {
	name := t.GetQualifiedName(node)
	// Go-only: language-adapter — AddReferenceWithName takes a Go string, so
	// its checkNotNull(name) of a null qualified name happens here.
	jutil.CheckNotNull(name)
	t.AddReferenceWithName(node, *name)
}

// Java: PreprocessorSymbolTable.addReference(Node, String)
func (t *PreprocessorSymbolTable) AddReferenceWithName(node rhino.Node, name string) {
	// checkNotNull(name): a Go string is never null.

	t.symbols.ComputeIfAbsent(name, func(k string) rhino.SimpleSlot { return rhino.NewSimpleSlot(k, nil, true) })

	t.refs.Put(name, NewPreprocessorSymbolTable_Reference(t.symbols.Get(name), node))
}

// GetQualifiedName: this variant of Node#getQualifiedName is adds special
// support STRING nodes which represent module names.
//
// Java: PreprocessorSymbolTable.getQualifiedName(Node)
func (t *PreprocessorSymbolTable) GetQualifiedName(n rhino.Node) *string {
	if n.IsStringLit() {
		s := n.GetString()
		return &s
	}
	return n.GetQualifiedName()
}

// PreprocessorSymbolTable_Reference is a reference to a preprocessor
// symbol.
type PreprocessorSymbolTable_Reference struct {
	rhino.SimpleReference[rhino.SimpleSlot]
}

// Java: PreprocessorSymbolTable.Reference(SimpleSlot, Node)
func NewPreprocessorSymbolTable_Reference(symbol rhino.SimpleSlot, node rhino.Node) *PreprocessorSymbolTable_Reference {
	r := &PreprocessorSymbolTable_Reference{}
	r.InitSimpleReference(symbol, node)
	return r
}

// PreprocessorSymbolTable_CachedInstanceFactory is an object that maybe
// contains instance of the table. This object is needed because
// PreprocessorSymbolTable is used by multiple passes in different parts of
// code which initialized at different times (some even before compiler
// object is created). Instead instance of factory is passed around. Each
// pass that uses PreprocessorSymbolTable has to call MaybeInitialize()
// before getting instance.
type PreprocessorSymbolTable_CachedInstanceFactory struct {
	instance *PreprocessorSymbolTable
}

// Java: PreprocessorSymbolTable.CachedInstanceFactory()
func NewPreprocessorSymbolTable_CachedInstanceFactory() *PreprocessorSymbolTable_CachedInstanceFactory {
	return &PreprocessorSymbolTable_CachedInstanceFactory{}
}

// Java: PreprocessorSymbolTable.CachedInstanceFactory.maybeInitialize(AbstractCompiler)
func (f *PreprocessorSymbolTable_CachedInstanceFactory) MaybeInitialize(compiler AbstractCompiler) {
	if compiler.GetOptions().PreservesDetailedSourceInfo() {
		root := compiler.GetRoot()
		if f.instance == nil || f.instance.GetRootNode() != root {
			f.instance = NewPreprocessorSymbolTable(root)
		}
	}
}

// Java: PreprocessorSymbolTable.CachedInstanceFactory.getInstanceOrNull()
func (f *PreprocessorSymbolTable_CachedInstanceFactory) GetInstanceOrNull() *PreprocessorSymbolTable {
	return f.instance
}

// Compile-time checks of the interfaces PreprocessorSymbolTable implements.
var (
	_ rhino.StaticTypedScope                                                        = (*PreprocessorSymbolTable)(nil)
	_ rhino.StaticSymbolTable[rhino.SimpleSlot, *PreprocessorSymbolTable_Reference] = (*PreprocessorSymbolTable)(nil)
)
