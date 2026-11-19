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
// Port of com/google/javascript/jscomp/PreprocessorSymbolTable.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

var lambdaSingleton_preprocessor_symbol_table_1 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_preprocessor_symbol_table_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_7 any) any {
		k := j2gort.CastString(lp_7)
		_ = k
		return pkg_rhino.NewSimpleSlot(k, nil, true)
	})
}

// Java class: com.google.javascript.jscomp.PreprocessorSymbolTable
type PreprocessorSymbolTable struct {
	symbols j2gort.OrderedMap[string, pkg_rhino.SimpleSlot] `gson:"symbols"`
	refs    j2gort.Ext_com_google_common_collect_Multimap   `gson:"refs"`
	root    pkg_rhino.Node                                  `gson:"root"`
	idHash_ atomic.Int32
}

func (this_ *PreprocessorSymbolTable) BindSelf_(s any) {
	_ = s
}
func (this_ *PreprocessorSymbolTable) IsJavaNil_() bool { return this_ == nil }
func (*PreprocessorSymbolTable) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PreprocessorSymbolTable", "PreprocessorSymbolTable"
}
func (this_ *PreprocessorSymbolTable) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PreprocessorSymbolTable) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PreprocessorSymbolTable", "com.google.javascript.jscomp.PreprocessorSymbolTable", func(v_ any) bool {
		return j2gort.InstanceOf[*PreprocessorSymbolTable](v_)
	})
}
func (*PreprocessorSymbolTable) Implements_rhino_StaticTypedScope()  {}
func (*PreprocessorSymbolTable) Implements_rhino_StaticSymbolTable() {}
func (*PreprocessorSymbolTable) Implements_rhino_StaticScope()       {}

func (this_ *PreprocessorSymbolTable) initFields_PreprocessorSymbolTable() {
	_ = this_
	this_.symbols = j2gort.NewOrderedMap[string, pkg_rhino.SimpleSlot]()
	this_.refs = j2gort.AsIface[j2gort.Ext_com_google_common_collect_Multimap](j2gort.Ext_com_google_common_collect_ArrayListMultimap_Create_0())
}

// Java constructor: PreprocessorSymbolTable(com.google.javascript.rhino.Node root)
func NewPreprocessorSymbolTable(root pkg_rhino.Node) *PreprocessorSymbolTable {
	this_ := &PreprocessorSymbolTable{}
	InitPreprocessorSymbolTable(this_, root)
	return this_
}

func InitPreprocessorSymbolTable(this_ *PreprocessorSymbolTable, root pkg_rhino.Node) {
	_ = root
	this_.initFields_PreprocessorSymbolTable()
	this_.root = root
}

// Java method: PreprocessorSymbolTable.getRootNode()
func (this_ *PreprocessorSymbolTable) GetRootNode_0() pkg_rhino.Node {
	return this_.root
}

// Java method: PreprocessorSymbolTable.getTypeOfThis()
func (this_ *PreprocessorSymbolTable) GetTypeOfThis() pkg_rhino.JSType {
	return nil
}

// Java method: PreprocessorSymbolTable.getParentScope()
func (this_ *PreprocessorSymbolTable) GetParentScope() any {
	return nil
}

// Java method: PreprocessorSymbolTable.getSlot(java.lang.String name)
func (this_ *PreprocessorSymbolTable) GetSlot(name string) any {
	_ = name
	return this_.symbols.Get(name)
}

// Java method: PreprocessorSymbolTable.getOwnSlot(java.lang.String name)
func (this_ *PreprocessorSymbolTable) GetOwnSlot(name string) any {
	_ = name
	return j2gort.AsIface[pkg_rhino.SimpleSlot](this_.GetSlot(name))
}

// Java method: PreprocessorSymbolTable.getReferences(com.google.javascript.rhino.jstype.SimpleSlot symbol)
func (this_ *PreprocessorSymbolTable) GetReferences(arg_1 any) any {
	symbol := j2gort.AsIface[pkg_rhino.SimpleSlot](arg_1)
	_ = symbol
	return j2gort.Ext_java_util_Collections_UnmodifiableCollection(j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.MultimapGetStr(this_.refs, symbol.GetSimpleSlotBase(

	// Java method: PreprocessorSymbolTable.getAllSymbols()
	).GetName())))
}

func (this_ *PreprocessorSymbolTable) GetAllSymbols() any {
	return j2gort.Ext_java_util_Collections_UnmodifiableCollection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.symbols.LiveValues()))
}

// Java method: PreprocessorSymbolTable.getScope(com.google.javascript.rhino.jstype.SimpleSlot slot)
func (this_ *PreprocessorSymbolTable) GetScope(arg_2 any) any {
	slot := j2gort.AsIface[pkg_rhino.SimpleSlot](arg_2)
	_ = slot
	return this_
}

// Java method: PreprocessorSymbolTable.addReference(com.google.javascript.rhino.Node node)
func (this_ *PreprocessorSymbolTable) AddReference_1(node pkg_rhino.Node) {
	_ = node
	this_.AddReference_2(node, this_.GetQualifiedName(node))
}

// Java method: PreprocessorSymbolTable.addReference(com.google.javascript.rhino.Node node, java.lang.String name)
func (this_ *PreprocessorSymbolTable) AddReference_2(node pkg_rhino.Node, name string) {
	_ = node
	_ = name
	j2gort.CheckNotNullString(name)
	j2gort.OrderedMapComputeIfAbsentFn(&this_.symbols, name, nil, lambdaSingleton_preprocessor_symbol_table_1)

	j2gort.MultimapPutStr(this_.refs, name, NewPreprocessorSymbolTable_Reference(this_.symbols.Get(name), node))
}

// Java method: PreprocessorSymbolTable.getQualifiedName(com.google.javascript.rhino.Node n)
func (this_ *PreprocessorSymbolTable) GetQualifiedName(n pkg_rhino.Node) string {
	_ = n
	if n.IsStringLit() {
		return n.GetString()
	} else {
		return n.GetQualifiedName()
	}
}

func (this_ *PreprocessorSymbolTable) LookupQualifiedName_1_QualifiedName(p0_ pkg_rhino.QualifiedName) pkg_rhino.JSType {
	return pkg_rhino.IfaceDefault_StaticTypedScope_LookupQualifiedName_1_QualifiedName(this_, p0_)
}

func (this_ *PreprocessorSymbolTable) LookupQualifiedName_1_string(p0_ string) pkg_rhino.JSType {
	return pkg_rhino.IfaceDefault_StaticTypedScope_LookupQualifiedName_1_string(this_, p0_)
}

func (this_ *PreprocessorSymbolTable) LookupQualifiedNameNode(p0_ pkg_rhino.Node) pkg_rhino.JSType {
	return pkg_rhino.IfaceDefault_StaticTypedScope_LookupQualifiedNameNode(this_, p0_)
}

func (this_ *PreprocessorSymbolTable) GetTopmostScopeOfEventualDeclaration(p0_ string) pkg_rhino.StaticScope {
	return pkg_rhino.IfaceDefault_StaticScope_GetTopmostScopeOfEventualDeclaration(this_, p0_)
}

// Java class: com.google.javascript.jscomp.PreprocessorSymbolTable.Reference
type PreprocessorSymbolTable_Reference struct {
	pkg_rhino.SimpleReferenceBase
}

func (this_ *PreprocessorSymbolTable_Reference) BindSelf_(s any) {
	_ = s
}
func (this_ *PreprocessorSymbolTable_Reference) IsJavaNil_() bool { return this_ == nil }
func (*PreprocessorSymbolTable_Reference) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PreprocessorSymbolTable$Reference", "Reference"
}
func (this_ *PreprocessorSymbolTable_Reference) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.SimpleReferenceBase.IdentityHash_()
}
func (*PreprocessorSymbolTable_Reference) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PreprocessorSymbolTable.Reference", "com.google.javascript.jscomp.PreprocessorSymbolTable$Reference", func(v_ any) bool {
		return j2gort.InstanceOf[*PreprocessorSymbolTable_Reference](v_)
	})
}
func (*PreprocessorSymbolTable_Reference) Implements_rhino_SimpleReference() {}
func (*PreprocessorSymbolTable_Reference) Implements_rhino_StaticTypedRef()  {}
func (*PreprocessorSymbolTable_Reference) Implements_rhino_StaticRef()       {}

func (this_ *PreprocessorSymbolTable_Reference) initFields_PreprocessorSymbolTable_Reference() {
	_ = this_
}

// Java constructor: Reference(com.google.javascript.rhino.jstype.SimpleSlot symbol, com.google.javascript.rhino.Node node)
func NewPreprocessorSymbolTable_Reference(symbol pkg_rhino.SimpleSlot, node pkg_rhino.Node) *PreprocessorSymbolTable_Reference {
	this_ := &PreprocessorSymbolTable_Reference{}
	InitPreprocessorSymbolTable_Reference(this_, symbol, node)
	return this_
}

func InitPreprocessorSymbolTable_Reference(this_ *PreprocessorSymbolTable_Reference, symbol pkg_rhino.SimpleSlot, node pkg_rhino.Node) {
	_ = symbol
	_ = node
	pkg_rhino.InitSimpleReference(&this_.SimpleReferenceBase, pkg_rhino.Upcast_StaticTypedSlot(symbol), node)
	this_.initFields_PreprocessorSymbolTable_Reference()
}

// Java class: com.google.javascript.jscomp.PreprocessorSymbolTable.CachedInstanceFactory
type PreprocessorSymbolTable_CachedInstanceFactory struct {
	instance *PreprocessorSymbolTable `gson:"instance"`
	idHash_  atomic.Int32
}

func (this_ *PreprocessorSymbolTable_CachedInstanceFactory) BindSelf_(s any) {
	_ = s
}
func (this_ *PreprocessorSymbolTable_CachedInstanceFactory) IsJavaNil_() bool { return this_ == nil }
func (*PreprocessorSymbolTable_CachedInstanceFactory) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PreprocessorSymbolTable$CachedInstanceFactory", "CachedInstanceFactory"
}
func (this_ *PreprocessorSymbolTable_CachedInstanceFactory) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PreprocessorSymbolTable_CachedInstanceFactory) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PreprocessorSymbolTable.CachedInstanceFactory", "com.google.javascript.jscomp.PreprocessorSymbolTable$CachedInstanceFactory", func(v_ any) bool {
		return j2gort.InstanceOf[*PreprocessorSymbolTable_CachedInstanceFactory](v_)
	})
}

func (this_ *PreprocessorSymbolTable_CachedInstanceFactory) initFields_PreprocessorSymbolTable_CachedInstanceFactory() {
	_ = this_
}

// Java constructor: CachedInstanceFactory()
func NewPreprocessorSymbolTable_CachedInstanceFactory() *PreprocessorSymbolTable_CachedInstanceFactory {
	this_ := &PreprocessorSymbolTable_CachedInstanceFactory{}
	InitPreprocessorSymbolTable_CachedInstanceFactory(this_)
	return this_
}

func InitPreprocessorSymbolTable_CachedInstanceFactory(this_ *PreprocessorSymbolTable_CachedInstanceFactory) {
	this_.initFields_PreprocessorSymbolTable_CachedInstanceFactory()
}

// Java method: CachedInstanceFactory.maybeInitialize(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *PreprocessorSymbolTable_CachedInstanceFactory) MaybeInitialize(compiler AbstractCompiler) {
	_ = compiler
	if compiler.GetOptions().PreservesDetailedSourceInfo() {
		var root pkg_rhino.Node = compiler.GetRoot()
		_ = root
		if (this_.instance == nil) || (this_.instance.GetRootNode_0() != root) {
			this_.instance = NewPreprocessorSymbolTable(root)
		}
	}
}

// Java method: CachedInstanceFactory.getInstanceOrNull()
func (this_ *PreprocessorSymbolTable_CachedInstanceFactory) GetInstanceOrNull() *PreprocessorSymbolTable {
	return this_.instance
}
