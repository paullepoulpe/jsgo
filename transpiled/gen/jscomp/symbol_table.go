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
// Port of com/google/javascript/jscomp/SymbolTable.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_symbol_table_1 j2gort.Ext_java_util_function_Predicate
var lambdaSingleton_symbol_table_2 j2gort.Ext_com_google_common_base_Supplier

func init() {
	lambdaSingleton_symbol_table_1 = j2gort.NewNodeLambda_Ext_java_util_function_Predicate(func(rawNode_59 uint32) bool {
		n := pkg_rhino.Node(rawNode_59)
		_ = n
		if !(n.IsCall()) {
			return false
		}
		var left pkg_rhino.Node = n.GetFirstChild()
		_ = left
		var arg pkg_rhino.Node = n.GetSecondChild()
		_ = arg
		if (!((left.IsGetProp() && left.GetFirstChild().IsName()) && (j2gort.CheckNotNullString(left.GetFirstChild().GetString()) == "goog")) || (arg == 0)) || !(arg.IsStringLit()) {
			return false
		}
		return true
	})
	lambdaSingleton_symbol_table_2 = j2gort.NewLambda_Ext_com_google_common_base_Supplier(func() any { return j2gort.NewExt_java_util_LinkedHashMap() })
}

// Java class: com.google.javascript.jscomp.SymbolTable
type SymbolTable struct {
	symbols                *SimpleTable[pkg_rhino.Node, string, *SymbolTable_Symbol]              `gson:"symbols"`
	scopes                 j2gort.OrderedMap[pkg_rhino.Node, *SymbolTable_SymbolScope]            `gson:"scopes"`
	docInfos               []pkg_rhino.Node                                                       `gson:"docInfos"`
	globalScope            *SymbolTable_SymbolScope                                               `gson:"globalScope"`
	compiler               AbstractCompiler                                                       `gson:"compiler"`
	registry               *pkg_rhino.JSTypeRegistry                                              `gson:"registry"`
	symbolThatDeclaresType *pkg_base.LinkedIdentityHashMap[pkg_rhino.JSType, *SymbolTable_Symbol] `gson:"symbolThatDeclaresType"`
	sourceNameOrdering     j2gort.Ext_com_google_common_collect_Ordering                          `gson:"sourceNameOrdering"`
	nodeOrdering           j2gort.Ext_com_google_common_collect_Ordering                          `gson:"nodeOrdering"`
	lexicalScopeOrdering   j2gort.Ext_com_google_common_collect_Ordering                          `gson:"lexicalScopeOrdering"`
	symbolOrdering         j2gort.Ext_com_google_common_collect_Ordering                          `gson:"symbolOrdering"`
	idHash_                atomic.Int32
}

func (this_ *SymbolTable) BindSelf_(s any) {
	_ = s
}
func (this_ *SymbolTable) IsJavaNil_() bool { return this_ == nil }
func (*SymbolTable) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SymbolTable", "SymbolTable"
}
func (this_ *SymbolTable) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SymbolTable) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SymbolTable", "com.google.javascript.jscomp.SymbolTable", func(v_ any) bool {
		return j2gort.InstanceOf[*SymbolTable](v_)
	})
}

var static_SymbolTable_logger *j2gort.Ext_java_util_logging_Logger

func StaticGet_SymbolTable_Logger() *j2gort.Ext_java_util_logging_Logger {
	clinit_SymbolTable()
	return static_SymbolTable_logger
}
func StaticSet_SymbolTable_Logger(v *j2gort.Ext_java_util_logging_Logger) {
	clinit_SymbolTable()
	static_SymbolTable_logger = v
}
func StaticGet_SymbolTable_GLOBAL_THIS() string { return "*global*" }

var clinitDone_SymbolTable uint32
var clinitState_SymbolTable j2gort.ClinitState

func clinit_SymbolTable() {
	if atomic.LoadUint32(&clinitDone_SymbolTable) != 0 {
		return
	}
	clinitSlow_SymbolTable()
}

func clinitSlow_SymbolTable() {
	if !j2gort.EnterClinit(&clinitState_SymbolTable, &clinitDone_SymbolTable) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_SymbolTable, &clinitDone_SymbolTable)
	static_SymbolTable_logger = j2gort.Ext_java_util_logging_Logger_GetLogger_1(j2gort.CheckNotNullPtr(j2gort.ClassLiteral("com.google.javascript.jscomp.SymbolTable", "SymbolTable")).GetName())
}

func (this_ *SymbolTable) initFields_SymbolTable() {
	_ = this_
	this_.symbols = SimpleTable_Create[pkg_rhino.Node, string, *SymbolTable_Symbol]()
	this_.scopes = j2gort.NewOrderedMap[pkg_rhino.Node, *SymbolTable_SymbolScope]()
	this_.docInfos = []pkg_rhino.Node{}

	this_.symbolThatDeclaresType = pkg_base.NewLinkedIdentityHashMap[pkg_rhino.JSType, *SymbolTable_Symbol]()
	this_.sourceNameOrdering = j2gort.Ext_com_google_common_collect_Ordering_Natural().NullsFirst()
	this_.nodeOrdering = NewSymbolTable__Anon_94(this_)
	this_.lexicalScopeOrdering = NewSymbolTable__Anon_95(this_)
	this_.symbolOrdering = NewSymbolTable__Anon_96(this_)
}

// Java constructor: SymbolTable(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func NewSymbolTable(compiler AbstractCompiler, registry *pkg_rhino.JSTypeRegistry) *SymbolTable {
	clinit_SymbolTable()
	this_ := &SymbolTable{}
	InitSymbolTable(this_, compiler, registry)
	return this_
}

func InitSymbolTable(this_ *SymbolTable, compiler AbstractCompiler, registry *pkg_rhino.JSTypeRegistry) {
	_ = compiler
	_ = registry
	this_.initFields_SymbolTable()
	this_.compiler = compiler
	this_.registry = registry
}

// Java method: SymbolTable.getReferences(com.google.javascript.jscomp.SymbolTable.Symbol symbol)
func (this_ *SymbolTable) GetReferences(symbol *SymbolTable_Symbol) j2gort.Ext_java_lang_Iterable {
	_ = symbol
	return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_java_util_Collections_UnmodifiableCollection(j2gort.AsIface[j2gort.Ext_java_util_Collection](symbol.references.Values())))
}

// Java method: SymbolTable.getReferenceList(com.google.javascript.jscomp.SymbolTable.Symbol symbol)
func (this_ *SymbolTable) GetReferenceList(symbol *SymbolTable_Symbol) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = symbol
	return j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](symbol.references.Values()))
}

// Java method: SymbolTable.getAllSymbols()
func (this_ *SymbolTable) GetAllSymbols() j2gort.Ext_com_google_common_collect_ImmutableList {
	return j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(this_.symbols.Values())
}

// Java method: SymbolTable.getAllSymbolsSorted()
func (this_ *SymbolTable) GetAllSymbolsSorted() j2gort.Ext_java_util_List {
	return this_.GetNaturalSymbolOrdering().SortedCopy(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.symbols.Values()))
}

// Java method: SymbolTable.getNaturalSymbolOrdering()
func (this_ *SymbolTable) GetNaturalSymbolOrdering() j2gort.Ext_com_google_common_collect_Ordering {
	return this_.symbolOrdering
}

// Java method: SymbolTable.getScope(com.google.javascript.jscomp.SymbolTable.Symbol slot)
func (this_ *SymbolTable) GetScope(slot *SymbolTable_Symbol) *SymbolTable_SymbolScope {
	_ = slot
	return slot.scope
}

// Java method: SymbolTable.getAllJSDocInfoNodes()
func (this_ *SymbolTable) GetAllJSDocInfoNodes() j2gort.Ext_java_util_Collection {
	return j2gort.SliceUnmodifiableListLive(&this_.docInfos)
}

// Java method: SymbolTable.getEnclosingScope(com.google.javascript.rhino.Node n)
func (this_ *SymbolTable) GetEnclosingScope(n pkg_rhino.Node) *SymbolTable_SymbolScope {
	_ = n
	var current pkg_rhino.Node = n.GetParent()
	_ = current
	if n.IsName() && n.GetParent().IsFunction() {
		current = current.GetParent()
	}
	for ; current != 0; current = current.GetParent() {
		if this_.scopes.ContainsKey(current) {
			return this_.scopes.Get(current)
		}
	}
	return nil
}

// Java method: SymbolTable.getEnclosingFunctionScope(com.google.javascript.rhino.Node n)
func (this_ *SymbolTable) GetEnclosingFunctionScope(n pkg_rhino.Node) *SymbolTable_SymbolScope {
	_ = n
	var current pkg_rhino.Node = n.GetParent()
	_ = current
	if (n.IsName() && (current != 0)) && current.IsFunction() {
		current = current.GetParent()
	}
	for ; current != 0; current = current.GetParent() {
		var scope *SymbolTable_SymbolScope = this_.scopes.Get(current)
		_ = scope
		if (scope != nil) && !(scope.IsBlockScope()) {
			return scope
		}
	}
	return this_.globalScope
}

// Java method: SymbolTable.getParameterInFunction(com.google.javascript.jscomp.SymbolTable.Symbol sym, java.lang.String paramName)
func (this_ *SymbolTable) GetParameterInFunction(sym *SymbolTable_Symbol, paramName string) *SymbolTable_Symbol {
	_ = sym
	_ = paramName
	var scope *SymbolTable_SymbolScope = this_.priv_SymbolTable_getScopeInFunction_5205(sym)
	_ = scope
	if scope != nil {
		var param *SymbolTable_Symbol = scope.GetSlot(paramName)
		_ = param
		if (param != nil) && (param.scope == scope) {
			return param
		}
	}
	return nil
}

// Java method: SymbolTable.getScopeInFunction(com.google.javascript.jscomp.SymbolTable.Symbol sym)
func (this_ *SymbolTable) priv_SymbolTable_getScopeInFunction_5205(sym *SymbolTable_Symbol) *SymbolTable_SymbolScope {
	_ = sym
	var type_ *pkg_rhino.FunctionType = sym.GetFunctionType()
	_ = type_
	if type_ == nil {
		return nil
	}
	var functionNode pkg_rhino.Node = type_.GetSource()
	_ = functionNode
	if functionNode == 0 {
		return nil
	}
	return this_.scopes.Get(functionNode)
}

// Java method: SymbolTable.getSymbolForScope(com.google.javascript.jscomp.SymbolTable.SymbolScope scope)
func (this_ *SymbolTable) GetSymbolForScope(scope *SymbolTable_SymbolScope) *SymbolTable_Symbol {
	_ = scope
	if scope.GetSymbolForScope() == nil {
		scope.SetSymbolForScope(this_.priv_SymbolTable_findSymbolForScope_5207(scope))
	}
	return scope.GetSymbolForScope()
}

// Java method: SymbolTable.findSymbolForScope(com.google.javascript.jscomp.SymbolTable.SymbolScope scope)
func (this_ *SymbolTable) priv_SymbolTable_findSymbolForScope_5207(scope *SymbolTable_SymbolScope) *SymbolTable_Symbol {
	_ = scope
	var rootNode pkg_rhino.Node = scope.GetRootNode()
	_ = rootNode
	if rootNode.GetParent() == 0 {
		return this_.globalScope.GetSlot("*global*")
	}
	if !(rootNode.IsFunction()) {
		return nil
	}
	var name string = NodeUtil_GetBestLValueName(NodeUtil_GetBestLValue(rootNode))
	_ = name
	if name == j2gort.NullString {
		return nil
	} else {
		return scope.GetParentScope().GetQualifiedSlot(name)
	}
}

// Java method: SymbolTable.getAllSymbolsForTypeOf(com.google.javascript.jscomp.SymbolTable.Symbol sym)
func (this_ *SymbolTable) GetAllSymbolsForTypeOf(sym *SymbolTable_Symbol) j2gort.Ext_java_lang_Iterable {
	_ = sym
	return j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.GetAllSymbolsForType(this_.priv_SymbolTable_getType_5258(pkg_rhino.ToIface_StaticSlot(sym))))
}

// Java method: SymbolTable.getGlobalScope()
func (this_ *SymbolTable) GetGlobalScope() *SymbolTable_SymbolScope {
	return this_.globalScope
}

// Java method: SymbolTable.getSymbolDeclaredBy(com.google.javascript.rhino.jstype.FunctionType fn)
func (this_ *SymbolTable) GetSymbolDeclaredBy_1_FunctionType(fn *pkg_rhino.FunctionType) *SymbolTable_Symbol {
	_ = fn
	if !(fn.IsConstructor() || fn.IsInterface()) {
		j2gort.CheckState(false)
	}
	var instanceType pkg_rhino.ObjectType = fn.GetInstanceType()
	_ = instanceType
	return this_.priv_SymbolTable_getSymbolForName_5214(fn.GetSource(), instanceType.GetReferenceName())
}

// Java method: SymbolTable.getSymbolDeclaredBy(com.google.javascript.rhino.jstype.EnumType enumType)
func (this_ *SymbolTable) GetSymbolDeclaredBy_1_EnumType(enumType *pkg_rhino.EnumType) *SymbolTable_Symbol {
	_ = enumType
	return this_.priv_SymbolTable_getSymbolForName_5214(enumType.GetSource(), enumType.GetElementsType().GetReferenceName())
}

// Java method: SymbolTable.getSymbolForInstancesOf(com.google.javascript.jscomp.SymbolTable.Symbol sym)
func (this_ *SymbolTable) GetSymbolForInstancesOf_1_SymbolTable_Symbol(sym *SymbolTable_Symbol) *SymbolTable_Symbol {
	_ = sym
	var fn *pkg_rhino.FunctionType = sym.GetFunctionType()
	_ = fn
	if (fn != nil) && fn.IsNominalConstructorOrInterface() {
		return this_.GetSymbolForInstancesOf_1_FunctionType(fn)
	}
	return nil
}

// Java method: SymbolTable.getSymbolForInstancesOf(com.google.javascript.rhino.jstype.FunctionType fn)
func (this_ *SymbolTable) GetSymbolForInstancesOf_1_FunctionType(fn *pkg_rhino.FunctionType) *SymbolTable_Symbol {
	_ = fn
	if !(fn.IsConstructor() || fn.IsInterface()) {
		j2gort.CheckState(false)
	}
	var pType pkg_rhino.ObjectType = fn.GetPrototype()
	_ = pType
	return this_.priv_SymbolTable_getSymbolForName_5214(fn.GetSource(), pType.GetReferenceName())
}

// Java method: SymbolTable.getSymbolForName(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node source, java.lang.String name)
func (this_ *SymbolTable) priv_SymbolTable_getSymbolForName_5214(source pkg_rhino.Node, name string) *SymbolTable_Symbol {
	_ = source
	_ = name
	if (name == j2gort.NullString) || (this_.globalScope == nil) {
		return nil
	}
	var scope *SymbolTable_SymbolScope
	_ = scope
	if source == 0 {
		scope = this_.globalScope
	} else {
		scope = this_.GetEnclosingScope(source)
	}
	if scope == nil {
		return nil
	}
	if scope.IsModuleScope() {
		var moduleNameWithDot string = (j2gort.StrValString(j2gort.CastString(j2gort.Ext_com_google_common_collect_Iterables_GetFirst(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.AsIface[ModuleMetadataMap_ModuleMetadata](j2gort.MapGetStr(this_.compiler.GetModuleMetadataMap().GetModulesByPath(), scope.GetRootNode().GetSourceFileName())).GoogNamespaces()), nil))) + ".")
		_ = moduleNameWithDot
		if javalang.StringStartsWith(j2gort.CheckNotNullString(name), j2gort.CheckNotNullString(moduleNameWithDot)) {
			name = j2gort.JStrSubFrom(name, j2gort.JStrLen(moduleNameWithDot))
		}
	}
	return scope.GetQualifiedSlot(name)
}

// Java method: SymbolTable.getAllSymbolsForType(com.google.javascript.rhino.jstype.JSType type)
func (this_ *SymbolTable) GetAllSymbolsForType(type_ pkg_rhino.JSType) j2gort.Ext_java_util_List {
	_ = type_
	if type_ == nil || type_.IsJavaNil_() {
		return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	}
	if type_.GetJSTypeBase().IsUnionType() {
		var result []*SymbolTable_Symbol = []*SymbolTable_Symbol{}
		_ = result
		var union *pkg_rhino.UnionType = type_.ToMaybeUnionType()
		_ = union
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = union.GetAlternatesCount()
			_ = n
			for ; i < n; i = i + 1 {
				var altSym *SymbolTable_Symbol = this_.priv_SymbolTable_getSymbolForTypeHelper_5216(union.GetAlternateAt(i), true)
				_ = altSym
				if altSym != nil {
					result = append(result, j2gort.CheckNotNullElem(altSym))
				}
			}
		}
		return j2gort.ImmListFromTypedSlice(result)
	}
	var result_2 *SymbolTable_Symbol = this_.priv_SymbolTable_getSymbolForTypeHelper_5216(type_, true)
	_ = result_2
	if result_2 == nil {
		return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
	} else {
		return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfPtr(result_2))
	}
}

// Java method: SymbolTable.getSymbolForTypeHelper(com.google.javascript.rhino.jstype.JSType type, boolean linkToCtor)
func (this_ *SymbolTable) priv_SymbolTable_getSymbolForTypeHelper_5216(type_ pkg_rhino.JSType, linkToCtor bool) *SymbolTable_Symbol {
	_ = type_
	_ = linkToCtor
	if type_ == nil || type_.IsJavaNil_() {
		return nil
	}
	base_type__ := type_.GetJSTypeBase()
	if base_type__.IsGlobalThisType() {
		return this_.globalScope.GetSlot("*global*")
	} else if base_type__.IsNominalConstructorOrInterface() {
		if linkToCtor {
			return this_.globalScope.GetSlot("Function")
		} else {
			return this_.GetSymbolDeclaredBy_1_FunctionType(type_.ToMaybeFunctionType())
		}
	} else if type_.IsFunctionPrototypeType() {
		var ownerFn *pkg_rhino.FunctionType = (j2gort.AsIface[pkg_rhino.ObjectType](type_)).GetOwnerFunction()
		_ = ownerFn
		if !(ownerFn.IsConstructor()) && !(ownerFn.IsInterface()) {
			return nil
		}
		if linkToCtor {
			return this_.GetSymbolDeclaredBy_1_FunctionType(ownerFn)
		} else {
			return this_.GetSymbolForInstancesOf_1_FunctionType(ownerFn)
		}
	} else if type_.IsInstanceType() {
		var ownerFn_2 *pkg_rhino.FunctionType = (j2gort.AsIface[pkg_rhino.ObjectType](type_)).GetConstructor()
		_ = ownerFn_2
		if linkToCtor {
			return this_.GetSymbolDeclaredBy_1_FunctionType(ownerFn_2)
		} else {
			return this_.GetSymbolForInstancesOf_1_FunctionType(ownerFn_2)
		}
	} else if base_type__.IsFunctionType() {
		if linkToCtor {
			return this_.globalScope.GetSlot("Function")
		} else {
			return this_.globalScope.GetQualifiedSlot("Function.prototype")
		}
	} else if !(j2gort.IfaceEquals[pkg_rhino.JSType](type_.AutoboxesTo(), nil)) {
		return this_.priv_SymbolTable_getSymbolForTypeHelper_5216(type_.AutoboxesTo(), linkToCtor)
	} else if base_type__.IsEnumType() {
		return this_.GetSymbolDeclaredBy_1_EnumType(j2gort.CastPtr[pkg_rhino.EnumType](type_))
	} else if base_type__.IsRecordType() {
		return this_.symbolThatDeclaresType.Get(type_)
	}
	return nil
}

// Java method: SymbolTable.toDebugStringTree()
func (this_ *SymbolTable) ToDebugStringTree_0() string {
	var childrenScopesBuilder *j2gort.Ext_com_google_common_collect_ImmutableListMultimap_Builder = j2gort.Ext_com_google_common_collect_ImmutableListMultimap_S_Builder()
	_ = childrenScopesBuilder
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetAllScopes()))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			scope := j2gort.CastPtr[SymbolTable_SymbolScope](rawElem_2)
			_ = scope
			if scope.GetParentScope() != nil {
				childrenScopesBuilder.Put_2(j2gort.PtrToAny(scope.GetParentScope()), j2gort.PtrToAny(scope))
			}
		}
	}
	var childrenScopes *j2gort.Ext_com_google_common_collect_ImmutableListMultimap = j2gort.CastPtr[j2gort.Ext_com_google_common_collect_ImmutableListMultimap](childrenScopesBuilder.Build())
	_ = childrenScopes
	var builder *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = builder
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetAllScopes()))
		for cur_3.HasNext() {
			rawElem_4 := cur_3.Next()
			scope_2 := j2gort.CastPtr[SymbolTable_SymbolScope](rawElem_4)
			_ = scope_2
			if scope_2.GetParentScope() == nil {
				this_.priv_SymbolTable_toDebugStringTree_5218(builder, "", scope_2, childrenScopes)
			}
		}
	}
	var visitedScopes j2gort.LocalOrderedSet[*SymbolTable_SymbolScope] = j2gort.LocalOrderedSetFromCol[*SymbolTable_SymbolScope](this_.GetAllScopes())
	_ = visitedScopes
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetAllSymbols()))
		for cur_5.HasNext() {
			rawElem_6 := cur_5.Next()
			symbol := j2gort.CastPtr[SymbolTable_Symbol](rawElem_6)
			_ = symbol
			if (symbol.propertyScope == nil) || visitedScopes.Contains(symbol.propertyScope) {
				continue
			}
			visitedScopes.Add(symbol.GetPropertyScope())
			this_.priv_SymbolTable_toDebugStringTree_5218(builder, "", symbol.propertyScope, childrenScopes)
		}
	}
	return builder.String()
}

// Java method: SymbolTable.toDebugStringTree(java.lang.StringBuilder builder, java.lang.String prefix, com.google.javascript.jscomp.SymbolTable.SymbolScope scope, com.google.common.collect.ImmutableListMultimap<com.google.javascript.jscomp.SymbolTable.SymbolScope,com.google.javascript.jscomp.SymbolTable.SymbolScope> childrenScopes)
func (this_ *SymbolTable) priv_SymbolTable_toDebugStringTree_5218(builder *j2gort.Ext_java_lang_StringBuilder, prefix string, scope *SymbolTable_SymbolScope, childrenScopes *j2gort.Ext_com_google_common_collect_ImmutableListMultimap) {
	_ = builder
	_ = prefix
	_ = scope
	_ = childrenScopes
	j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](builder.Append_1_string(prefix)).Append_1_any_9859(j2gort.PtrToAny(scope))
	var childrenPrefix string = (j2gort.StrValString(prefix) + "    ")
	_ = childrenPrefix
	var printedSomething bool = false
	_ = printedSomething
	if !(scope.ownSymbols.IsEmpty()) {
		printedSomething = true
		j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](builder.Append_1_uint16(uint16(10))).Append_1_string(prefix)).Append_1_string("  Symbols:\n")
		{
			cur_7 := scope.ownSymbols.IterValues()
			for cur_7.HasNext() {

				symbol := cur_7.
					Next()
				_ = symbol
				this_.priv_SymbolTable_toDebugString_5220(builder, childrenPrefix, symbol, false)
			}
		}
	}
	if childrenScopes.ContainsKey(j2gort.PtrToAny(scope)) {
		printedSomething = true
		j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](builder.Append_1_uint16(uint16(10))).Append_1_string(prefix)).Append_1_string("  Scopes:\n")
		{
			cur_9 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableList](childrenScopes.Get(j2gort.PtrToAny(scope)))))
			for cur_9.HasNext() {
				rawElem_10 := cur_9.Next()
				childScope := j2gort.CastPtr[SymbolTable_SymbolScope](rawElem_10)
				_ = childScope
				this_.priv_SymbolTable_toDebugStringTree_5218(builder, childrenPrefix, childScope, childrenScopes)
			}
		}
	}
	if !(printedSomething) {
		builder.Append_1_uint16(uint16(10))
	}
	builder.Append_1_uint16(uint16(10))
}

// Java method: SymbolTable.toDebugString()
func (this_ *SymbolTable) ToDebugString_0() string {
	var builder *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = builder
	{
		cur_11 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetAllSymbols()))
		for cur_11.HasNext() {
			rawElem_12 := cur_11.Next()
			symbol := j2gort.CastPtr[SymbolTable_Symbol](rawElem_12)
			_ = symbol
			this_.priv_SymbolTable_toDebugString_5220(builder, "", symbol, true)
		}
	}
	return builder.String()
}

// Java method: SymbolTable.toDebugString(java.lang.StringBuilder builder, java.lang.String prefix, com.google.javascript.jscomp.SymbolTable.Symbol symbol, boolean printScope)
func (this_ *SymbolTable) priv_SymbolTable_toDebugString_5220(builder *j2gort.Ext_java_lang_StringBuilder, prefix string, symbol *SymbolTable_Symbol, printScope bool) {
	_ = builder
	_ = prefix
	_ = symbol
	_ = printScope
	j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](builder.Append_1_string(prefix)).Append_1_string("'")).Append_1_string(symbol.GetName())).Append_1_string("'")
	if printScope {
		j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](builder.Append_1_string(" in ")).Append_1_any_9859(j2gort.PtrToAny(symbol.scope))
	}
	builder.Append_1_uint16(uint16(10))
	var refCount int32 = int32(0)
	_ = refCount
	{
		cur_13 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetReferences(symbol)))
		for cur_13.HasNext() {
			rawElem_14 := cur_13.Next()
			ref := j2gort.CastPtr[SymbolTable_Reference](rawElem_14)
			_ = ref
			var node pkg_rhino.Node = ref.GetNode()
			_ = node
			j2gort.CastPtr[j2gort.Ext_java_lang_StringBuilder](builder.Append_1_string(prefix)).Append_1_string(j2gort.Ext_java_lang_String_Format_2("  Ref %d: %s line: %d col: %d len: %d %s\n", []any{j2gort.ValueOfInt(refCount), j2gort.BoxString(node.GetSourceFileName()), j2gort.ValueOfInt(node.GetLineno()), j2gort.ValueOfInt(node.GetCharno()), j2gort.ValueOfInt(node.GetLength()), func() any {
				if node.IsIndexable() {
					return ""
				}
				return "non indexable"
			}()}))
			refCount = refCount + 1
		}
	}
}

// Java method: SymbolTable.addScopes(java.util.Collection<S> scopes)
func (this_ *SymbolTable) AddScopes(scopes j2gort.Ext_java_util_Collection) {
	_ = scopes
	{
		cur_15 := j2gort.NewIterCursor(j2gort.CheckNotNull(scopes))
		for cur_15.HasNext() {
			rawElem_16 := cur_15.Next()
			scope := j2gort.AsIface[pkg_rhino.StaticScope](rawElem_16)
			_ = scope
			this_.priv_SymbolTable_createScopeFrom_5256(scope)
		}
	}
}

// Java method: SymbolTable.findScopes(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *SymbolTable) FindScopes(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_TraverseRoots(this_.compiler, NewSymbolTable__Anon_90(this_), externs, root)
}

// Java method: SymbolTable.getAllScopes()
func (this_ *SymbolTable) GetAllScopes() j2gort.Ext_java_util_Collection {
	return j2gort.Ext_java_util_Collections_UnmodifiableCollection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.scopes.LiveValues()))
}

// Java method: SymbolTable.addAnonymousFunctions()
func (this_ *SymbolTable) AddAnonymousFunctions() {
	var scopes *j2gort.Ext_java_util_TreeSet = j2gort.NewExt_java_util_TreeSet_1_Ext_java_util_Comparator(j2gort.AsIface[j2gort.Ext_java_util_Comparator](this_.lexicalScopeOrdering))
	_ = scopes
	{
		cur_17 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetAllScopes()))
		for cur_17.HasNext() {
			rawElem_18 := cur_17.Next()
			scope := j2gort.CastPtr[SymbolTable_SymbolScope](rawElem_18)
			_ = scope
			if scope.IsLexicalScope() {
				scopes.Add_1(j2gort.PtrToAny(scope))
			}
		}
	}
	{
		cur_19 := j2gort.NewIterCursor(scopes)
		for cur_19.HasNext() {
			rawElem_20 := cur_19.Next()
			scope_2 := j2gort.CastPtr[SymbolTable_SymbolScope](rawElem_20)
			_ = scope_2
			this_.priv_SymbolTable_addAnonymousFunctionsInScope_5225(scope_2)
		}
	}
}

// Java method: SymbolTable.addAnonymousFunctionsInScope(com.google.javascript.jscomp.SymbolTable.SymbolScope scope)
func (this_ *SymbolTable) priv_SymbolTable_addAnonymousFunctionsInScope_5225(scope *SymbolTable_SymbolScope) {
	_ = scope
	var sym *SymbolTable_Symbol = this_.GetSymbolForScope(scope)
	_ = sym
	if sym == nil {
		var rootNode pkg_rhino.Node = scope.GetRootNode()
		_ = rootNode
		if ((((scope.IsLexicalScope() && !(scope.IsGlobalScope())) && (rootNode != 0)) && !(rootNode.IsFromExterns())) && (scope.GetParentScope() != nil)) && rootNode.IsFunction() {
			var parent *SymbolTable_SymbolScope = scope.GetParentScope()
			_ = parent
			var innerName string = ("function%" + j2gort.StrValInt32(scope.GetIndexInParent()))
			_ = innerName
			var type_ pkg_rhino.JSType = rootNode.GetJSType()
			_ = type_
			if rootNode.GetParent().IsMemberFunctionDef() {
				rootNode = rootNode.GetParent()
			}
			var anonymousFunctionSymbol *SymbolTable_Symbol = this_.priv_SymbolTable_declareSymbol_5233(innerName, type_, true, parent, rootNode, nil)
			_ = anonymousFunctionSymbol
			scope.SetSymbolForScope(anonymousFunctionSymbol)
		}
	}
}

// Java method: SymbolTable.addSymbolsFrom(com.google.javascript.rhino.StaticSymbolTable<S,R> otherSymbolTable)
func (this_ *SymbolTable) AddSymbolsFrom(otherSymbolTable pkg_rhino.StaticSymbolTable) {
	_ = otherSymbolTable
	{
		cur_21 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](otherSymbolTable.GetAllSymbols())))
		for cur_21.HasNext() {
			rawElem_22 := cur_21.Next()
			otherSymbol := j2gort.AsIface[pkg_rhino.StaticSlot](rawElem_22)
			_ = otherSymbol
			var name string = otherSymbol.GetName()
			_ = name
			var myScope *SymbolTable_SymbolScope = this_.priv_SymbolTable_createScopeFrom_5256(j2gort.AsIface[pkg_rhino.StaticScope](otherSymbolTable.GetScope(otherSymbol)))
			_ = myScope
			var decl pkg_rhino.StaticRef = this_.priv_SymbolTable_findBestDeclToAdd_5228(otherSymbolTable, otherSymbol)
			_ = decl
			var mySymbol *SymbolTable_Symbol = nil
			_ = mySymbol
			if decl != nil && !decl.IsJavaNil_() {
				var declNode pkg_rhino.Node = decl.GetNode()
				_ = declNode
				mySymbol = this_.priv_SymbolTable_isAnySymbolDeclared_5227(name, declNode, myScope)
				if mySymbol == nil {
					mySymbol = this_.priv_SymbolTable_copySymbolTo_5231(otherSymbol, declNode, myScope)
				}
			} else {
				mySymbol = myScope.GetOwnSlot(name)
			}
			if mySymbol != nil {
				{
					cur_23 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](otherSymbolTable.GetReferences(otherSymbol))))
					for cur_23.HasNext() {
						rawElem_24 := cur_23.Next()
						otherRef := j2gort.AsIface[pkg_rhino.StaticRef](rawElem_24)
						_ = otherRef
						if this_.priv_SymbolTable_isGoodRefToAdd_5229(otherRef) {
							mySymbol.DefineReferenceAt(otherRef.GetNode())
						}
					}
				}
			}
		}
	}
}

// Java method: SymbolTable.isAnySymbolDeclared(java.lang.String name, com.google.javascript.rhino.Node declNode, com.google.javascript.jscomp.SymbolTable.SymbolScope scope)
func (this_ *SymbolTable) priv_SymbolTable_isAnySymbolDeclared_5227(name string, declNode pkg_rhino.Node, scope *SymbolTable_SymbolScope) *SymbolTable_Symbol {
	_ = name
	_ = declNode
	_ = scope
	var sym *SymbolTable_Symbol = this_.symbols.Get(j2gort.BoxNode(declNode), j2gort.BoxString(name))
	_ = sym
	if sym == nil {
		return scope.ownSymbols.Get(name)
	}
	return sym
}

// Java method: SymbolTable.findBestDeclToAdd(com.google.javascript.rhino.StaticSymbolTable<S,R> otherSymbolTable, S slot)
func (this_ *SymbolTable) priv_SymbolTable_findBestDeclToAdd_5228(otherSymbolTable pkg_rhino.StaticSymbolTable, slot pkg_rhino.StaticSlot) pkg_rhino.StaticRef {
	_ = otherSymbolTable
	_ = slot
	var decl pkg_rhino.StaticRef = j2gort.AsIface[pkg_rhino.StaticRef](slot.GetDeclaration())
	_ = decl
	if this_.priv_SymbolTable_isGoodRefToAdd_5229(decl) {
		return decl
	}
	{
		cur_25 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](otherSymbolTable.GetReferences(slot))))
		for cur_25.HasNext() {
			rawElem_26 := cur_25.Next()
			ref := j2gort.AsIface[pkg_rhino.StaticRef](rawElem_26)
			_ = ref
			if this_.priv_SymbolTable_isGoodRefToAdd_5229(ref) {
				return ref
			}
		}
	}
	return nil
}

// Java method: SymbolTable.isGoodRefToAdd(com.google.javascript.rhino.@org.jspecify.annotations.Nullable StaticRef ref)
func (this_ *SymbolTable) priv_SymbolTable_isGoodRefToAdd_5229(ref pkg_rhino.StaticRef) bool {
	_ = ref
	return (((((ref != nil && !ref.IsJavaNil_()) && (ref.GetNode() != 0)) && !(j2gort.IfaceEquals[pkg_rhino.StaticSourceFile](ref.GetNode().GetStaticSourceFile(), nil))) && !(NodeUtil_IsInSyntheticScript(ref.GetNode()))) && !(NodeUtil_IsGoogProvideCall(ref.GetNode())))
}

// Java method: SymbolTable.copySymbolTo(com.google.javascript.rhino.StaticSlot sym, com.google.javascript.jscomp.SymbolTable.SymbolScope scope)
func (this_ *SymbolTable) priv_SymbolTable_copySymbolTo_5230(sym pkg_rhino.StaticSlot, scope *SymbolTable_SymbolScope) *SymbolTable_Symbol {
	_ = sym
	_ = scope
	return this_.priv_SymbolTable_copySymbolTo_5231(sym, j2gort.AsIface[pkg_rhino.StaticRef](sym.GetDeclaration()).GetNode(), scope)
}

// Java method: SymbolTable.copySymbolTo(com.google.javascript.rhino.StaticSlot sym, com.google.javascript.rhino.Node declNode, com.google.javascript.jscomp.SymbolTable.SymbolScope scope)
func (this_ *SymbolTable) priv_SymbolTable_copySymbolTo_5231(sym pkg_rhino.StaticSlot, declNode pkg_rhino.Node, scope *SymbolTable_SymbolScope) *SymbolTable_Symbol {
	_ = sym
	_ = declNode
	_ = scope
	j2gort.CheckNotNullNode(declNode)
	return this_.priv_SymbolTable_declareSymbol_5233(sym.GetName(), this_.priv_SymbolTable_getType_5258(sym), this_.priv_SymbolTable_isTypeInferred_5260(sym), scope, declNode, sym.GetJSDocInfo())
}

// Java static method: SymbolTable.sanitizeSpecialChars(java.lang.String s)
func SymbolTable_SanitizeSpecialChars(s string) string {
	clinit_SymbolTable()
	_ = s
	return javalang.StringReplace(javalang.StringReplace(javalang.StringReplace(s, "\\", "\\\\"), "\x00", "\\0"), "\n", "\\n")
}

// Java method: SymbolTable.addSymbol(java.lang.String name, com.google.javascript.rhino.jstype.JSType type, boolean inferred, com.google.javascript.jscomp.SymbolTable.SymbolScope scope, com.google.javascript.rhino.Node declNode)
func (this_ *SymbolTable) priv_SymbolTable_addSymbol_5232(name string, type_ pkg_rhino.JSType, inferred bool, scope *SymbolTable_SymbolScope, declNode pkg_rhino.Node) *SymbolTable_Symbol {
	_ = name
	_ = type_
	_ = inferred
	_ = scope
	_ = declNode
	name = SymbolTable_SanitizeSpecialChars(name)
	var symbol *SymbolTable_Symbol = NewSymbolTable_Symbol(name, type_, inferred, scope)
	_ = symbol
	var replacedSymbol *SymbolTable_Symbol = this_.symbols.Put(declNode, name, symbol)
	_ = replacedSymbol
	if !(replacedSymbol == nil) {
		j2gort.CheckStateFmt(false, "Found duplicate symbol %s in global index. Type %s", j2gort.BoxString(name), type_)
	}
	replacedSymbol = scope.ownSymbols.Put(name, symbol)
	if !(replacedSymbol == nil) {
		j2gort.CheckStateFmt(false, "Found duplicate symbol %s in its scope. Type %s", j2gort.BoxString(name), type_)
	}
	return symbol
}

// Java method: SymbolTable.declareSymbol(java.lang.String name, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType type, boolean inferred, com.google.javascript.jscomp.SymbolTable.SymbolScope scope, com.google.javascript.rhino.Node declNode, com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo info)
func (this_ *SymbolTable) priv_SymbolTable_declareSymbol_5233(name string, type_ pkg_rhino.JSType, inferred bool, scope *SymbolTable_SymbolScope, declNode pkg_rhino.Node, info *pkg_rhino.JSDocInfo) *SymbolTable_Symbol {
	_ = name
	_ = type_
	_ = inferred
	_ = scope
	_ = declNode
	_ = info
	var symbol *SymbolTable_Symbol = this_.priv_SymbolTable_addSymbol_5232(name, type_, inferred, scope, declNode)
	_ = symbol
	symbol.SetJSDocInfo(info)
	symbol.SetDeclaration(symbol.DefineReferenceAt(declNode))
	return symbol
}

// Java method: SymbolTable.mergeSymbol(com.google.javascript.jscomp.SymbolTable.Symbol from, com.google.javascript.jscomp.SymbolTable.Symbol to)
func (this_ *SymbolTable) priv_SymbolTable_mergeSymbol_5234(from *SymbolTable_Symbol, to *SymbolTable_Symbol) {
	_ = from
	_ = to
	{
		cur_27 := j2gort.NewMapKeyCursor(from.references)
		for cur_27.HasNext() {
			nodeToMove := pkg_rhino.Node(cur_27.NextNode())
			_ = nodeToMove
			{
				var iifeCond_1 bool
				{

					rn_28 := nodeToMove
					an_29 := from.GetDeclarationNode()
					iifeCond_1 = j2gort.CheckNotNullNode(rn_28) == an_29
				}
				if !(iifeCond_1) {
					to.DefineReferenceAt(nodeToMove)
				}
			}
		}
	}
	this_.priv_SymbolTable_removeSymbol_5235(from)
}

// Java method: SymbolTable.removeSymbol(com.google.javascript.jscomp.SymbolTable.Symbol s)
func (this_ *SymbolTable) priv_SymbolTable_removeSymbol_5235(s *SymbolTable_Symbol) {
	_ = s
	var scope *SymbolTable_SymbolScope = this_.GetScope(s)
	_ = scope
	if !(s.Equals(j2gort.PtrToAny(scope.ownSymbols.Remove(s.GetName())))) {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Symbol not found in scope "+j2gort.StrVal(s))))
	}
	if !(s.Equals(j2gort.PtrToAny(this_.symbols.Remove(j2gort.BoxNode(j2gort.CastPtr[SymbolTable_Reference](s.GetDeclaration()).GetNode()), j2gort.BoxString(s.GetName()))))) {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Symbol not found in table "+j2gort.StrVal(s))))
	}
	if (s.propertyScope != nil) && s.Equals(j2gort.PtrToAny(s.propertyScope.GetSymbolForScope())) {
		{
			cur_30 := j2gort.NewIterCursor(j2gort.CheckNotNull(s.propertyScope.ownSymbols.ValuesToImmutableList()))
			for cur_30.HasNext() {
				rawElem_31 := cur_30.Next()
				childSymbol := j2gort.CastPtr[SymbolTable_Symbol](rawElem_31)
				_ = childSymbol
				this_.priv_SymbolTable_removeSymbol_5235(childSymbol)
			}
		}
		this_.scopes.Remove(s.GetDeclarationNode())
	}
}

// Java method: SymbolTable.renameSymbol(com.google.javascript.jscomp.SymbolTable.SymbolScope scope, com.google.javascript.jscomp.SymbolTable.Symbol sym, java.lang.String newName)
func (this_ *SymbolTable) priv_SymbolTable_renameSymbol_5236(scope *SymbolTable_SymbolScope, sym *SymbolTable_Symbol, newName string) {
	_ = scope
	_ = sym
	_ = newName
	if !(sym.propertyScope == nil) {
		j2gort.CheckStateStr(false, "Renaming supported for simple symbols that don't have property scopes.")
	}
	var existingSym *SymbolTable_Symbol = scope.GetSlot(newName)
	_ = existingSym
	if existingSym == nil {
		existingSym = this_.priv_SymbolTable_declareSymbol_5233(newName, this_.priv_SymbolTable_getType_5258(pkg_rhino.ToIface_StaticSlot(sym)), this_.priv_SymbolTable_isTypeInferred_5260(pkg_rhino.ToIface_StaticSlot(sym)), scope, sym.GetDeclarationNode(), sym.GetJSDocInfo())
	}
	this_.priv_SymbolTable_mergeSymbol_5234(sym, existingSym)
}

// Java method: SymbolTable.fillNamespaceReferences()
func (this_ *SymbolTable) FillNamespaceReferences() {
	{
		cur_32 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetAllSymbols()))
		for cur_32.HasNext() {
			rawElem_33 := cur_32.Next()
			symbol := j2gort.CastPtr[SymbolTable_Symbol](rawElem_33)
			_ = symbol
			var qName string = symbol.GetName()
			_ = qName
			var rootIndex int32 = j2gort.JStrIndexOfCp(qName, int32(uint16(46)))
			_ = rootIndex
			if rootIndex == int32(-1) {
				continue
			}
			var root *SymbolTable_Symbol = symbol.scope.GetQualifiedSlot(j2gort.JStrSub(qName, int32(0), rootIndex))
			_ = root
			if root == nil {
				continue
			}
			{
				cur_34 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetReferences(symbol)))
				for cur_34.HasNext() {
					rawElem_35 := cur_34.Next()
					ref := j2gort.CastPtr[SymbolTable_Reference](rawElem_35)
					_ = ref
					var currentNode pkg_rhino.Node = ref.GetNode()
					_ = currentNode
					if !(currentNode.IsQualifiedName()) {
						continue
					}
					for currentNode.IsGetProp() {
						currentNode = currentNode.GetFirstChild()
						var name string = currentNode.GetQualifiedName()
						_ = name
						if name == j2gort.NullString {
							continue
						}
						var namespace *SymbolTable_Symbol = this_.priv_SymbolTable_isAnySymbolDeclared_5227(name, currentNode, root.scope)
						_ = namespace
						if namespace == nil {
							namespace = root.scope.GetQualifiedSlot(name)
						}
						if (namespace == nil) && root.scope.IsGlobalScope() {
							var nodeType pkg_rhino.JSType = currentNode.GetJSType()
							_ = nodeType
							var symbolType pkg_rhino.JSType
							_ = symbolType
							if (nodeType != nil && !nodeType.IsJavaNil_()) && !(nodeType.IsFunctionPrototypeType()) {
								symbolType = nodeType
							} else {
								symbolType = this_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
							}
							namespace = this_.priv_SymbolTable_declareSymbol_5233(name, symbolType, true, root.scope, currentNode, nil)
						}
						if namespace != nil {
							namespace.DefineReferenceAt(currentNode)
						}
					}
				}
			}
		}
	}
}

// Java method: SymbolTable.fillPropertyScopes()
func (this_ *SymbolTable) FillPropertyScopes() {
	var types j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = types
	var exports j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = exports
	{
		cur_36 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetAllSymbols()))
		for cur_36.HasNext() {
			rawElem_37 := cur_36.Next()
			sym := j2gort.CastPtr[SymbolTable_Symbol](rawElem_37)
			_ = sym
			if this_.priv_SymbolTable_needsPropertyScope_5239(sym) {
				var name string = sym.GetName()
				_ = name
				if j2gort.CheckNotNullString(name) == "exports" {
					exports.Add_1(j2gort.PtrToAny(sym))
				} else {
					types.Add_1(j2gort.PtrToAny(sym))
				}
			}
		}
	}
	j2gort.Ext_java_util_Collections_Sort_2(types, j2gort.AsIface[j2gort.Ext_java_util_Comparator](this_.GetNaturalSymbolOrdering().Reverse()))
	var allTypes j2gort.Ext_java_lang_Iterable = j2gort.Ext_com_google_common_collect_Iterables_Concat_2(types, exports)
	_ = allTypes
	{
		cur_38 := j2gort.NewIterCursor(j2gort.CheckNotNull(allTypes))
		for cur_38.HasNext() {
			rawElem_39 := cur_38.Next()
			s := j2gort.CastPtr[SymbolTable_Symbol](rawElem_39)
			_ = s
			var type_ pkg_rhino.JSType = this_.priv_SymbolTable_getType_5258(pkg_rhino.ToIface_StaticSlot(s))
			_ = type_
			var symbolForType *SymbolTable_Symbol = this_.priv_SymbolTable_getSymbolForTypeHelper_5216(type_, false)
			_ = symbolForType
			if ((type_.GetJSTypeBase().IsNominalConstructorOrInterface() || type_.IsFunctionPrototypeType()) || type_.GetJSTypeBase().IsEnumType()) && !(s.Equals(j2gort.PtrToAny(symbolForType))) {
				continue
			}
			this_.symbolThatDeclaresType.Put(type_, s)
		}
	}
	{
		cur_40 := j2gort.NewIterCursor(j2gort.CheckNotNull(allTypes))
		for cur_40.HasNext() {
			rawElem_41 := cur_40.Next()
			s_2 := j2gort.CastPtr[SymbolTable_Symbol](rawElem_41)
			_ = s_2
			var type__2 pkg_rhino.JSType = this_.priv_SymbolTable_getType_5258(pkg_rhino.ToIface_StaticSlot(s_2))
			_ = type__2
			if (type__2 == nil || type__2.IsJavaNil_()) || s_2.Equals(j2gort.PtrToAny(this_.symbolThatDeclaresType.Get(type__2))) {
				this_.priv_SymbolTable_createPropertyScopeFor_5244(s_2)
			}
		}
	}
	{
		cur_42 := j2gort.NewIterCursor(j2gort.CheckNotNull(allTypes))
		for cur_42.HasNext() {
			rawElem_43 := cur_42.Next()
			s_3 := j2gort.CastPtr[SymbolTable_Symbol](rawElem_43)
			_ = s_3
			var type__3 pkg_rhino.JSType = this_.priv_SymbolTable_getType_5258(pkg_rhino.ToIface_StaticSlot(s_3))
			_ = type__3
			if (type__3 != nil && !type__3.IsJavaNil_()) && (this_.symbolThatDeclaresType.Get(type__3) != nil) {
				{
					r_44 := s_3
					r_44.propertyScope = this_.symbolThatDeclaresType.Get(type__3).GetPropertyScope()
				}
			}
		}
	}
	this_.PruneOrphanedNames()
}

// Java method: SymbolTable.needsPropertyScope(com.google.javascript.jscomp.SymbolTable.Symbol sym)
func (this_ *SymbolTable) priv_SymbolTable_needsPropertyScope_5239(sym *SymbolTable_Symbol) bool {
	_ = sym
	var jsType pkg_rhino.JSType = this_.priv_SymbolTable_getType_5258(pkg_rhino.ToIface_StaticSlot(sym))
	_ = jsType
	if (jsType != nil && !jsType.IsJavaNil_()) && jsType.GetJSTypeBase().IsRecordType() {
		return true
	}
	var type_ pkg_rhino.ObjectType = pkg_rhino.ObjectType_Cast(this_.priv_SymbolTable_getType_5258(pkg_rhino.ToIface_StaticSlot(sym)))
	_ = type_
	if type_ == nil || type_.IsJavaNil_() {
		return false
	}
	base_type__ := type_.GetObjectTypeBase()
	if type_.GetReferenceName() == j2gort.NullString {
		return true
	}
	if base_type__.IsNominalConstructorOrInterface() || base_type__.IsFunctionPrototypeType() {
		return true
	}
	if base_type__.IsEnumType() {
		return true
	}
	return false
}

// Java method: SymbolTable.pruneOrphanedNames()
func (this_ *SymbolTable) PruneOrphanedNames() {
	{
		cur_45 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetAllSymbols()))
	nextSymbol:
		for cur_45.HasNext() {
			rawElem_46 := cur_45.Next()
			s := j2gort.CastPtr[SymbolTable_Symbol](rawElem_46)
			_ = s
			var currentName string = s.GetName()
			_ = currentName
			if s.IsProperty() || (this_.compiler.GetModuleMap().GetClosureModule(currentName) != nil) {
				continue
			}
			var dot int32 = int32(-1)
			_ = dot
			for int32(-1) != (func() int32 {
				var asgn_47 int32 = j2gort.JStrLastIndexOfCp(currentName, int32(uint16(46)))
				dot = asgn_47
				return asgn_47
			}()) {
				currentName = j2gort.JStrSub(currentName, int32(0), dot)
				var owner *SymbolTable_Symbol = s.scope.GetQualifiedSlot(currentName)
				_ = owner
				if ((owner != nil) && !(j2gort.IfaceEquals[pkg_rhino.JSType](this_.priv_SymbolTable_getType_5258(pkg_rhino.ToIface_StaticSlot(owner)), nil))) && ((this_.priv_SymbolTable_getType_5258(pkg_rhino.ToIface_StaticSlot(owner)).GetJSTypeBase().IsNominalConstructorOrInterface() || this_.priv_SymbolTable_getType_5258(pkg_rhino.ToIface_StaticSlot(owner)).IsFunctionPrototypeType()) || this_.priv_SymbolTable_getType_5258(pkg_rhino.ToIface_StaticSlot(owner)).GetJSTypeBase().IsEnumType()) {
					this_.priv_SymbolTable_removeSymbol_5235(s)
					continue nextSymbol
				}
			}
		}
	}
}

// Java method: SymbolTable.fillPropertySymbols(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *SymbolTable) FillPropertySymbols(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NewSymbolTable_PropertyRefCollector(this_).Process_2(externs, root)
}

// Java method: SymbolTable.fillJSDocInfo(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *SymbolTable) FillJSDocInfo(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_TraverseRoots(this_.compiler, NewSymbolTable_JSDocInfoCollector(this_, this_.compiler.GetTypeRegistry()), externs, root)
	{
		cur_48 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetAllSymbols()))
		for cur_48.HasNext() {
			rawElem_49 := cur_48.Next()
			sym := j2gort.CastPtr[SymbolTable_Symbol](rawElem_49)
			_ = sym
			var info *pkg_rhino.JSDocInfo = sym.GetJSDocInfo()
			_ = info
			if info == nil {
				continue
			}
			{
				cur_50 := j2gort.NewIterCursor(j2gort.CheckNotNull(info.GetMarkers()))
				for cur_50.HasNext() {
					rawElem_51 := cur_50.Next()
					marker := j2gort.CastPtr[pkg_rhino.JSDocInfo_Marker](rawElem_51)
					_ = marker
					var pos pkg_rhino.SourcePosition = pkg_rhino.ToIface_SourcePosition(marker.GetNameNode())
					_ = pos
					if pos == nil || pos.IsJavaNil_() {
						continue
					}
					var paramNode pkg_rhino.Node = j2gort.CastNode[pkg_rhino.Node](pos.GetItem())
					_ = paramNode
					var name string = paramNode.GetString()
					_ = name
					var param *SymbolTable_Symbol = this_.GetParameterInFunction(sym, name)
					_ = param
					if param == nil {
						var typePos pkg_rhino.SourcePosition = pkg_rhino.ToIface_SourcePosition(marker.GetType())
						_ = typePos
						var type_ pkg_rhino.JSType = nil
						_ = type_
						if typePos != nil && !typePos.IsJavaNil_() {
							type_ = j2gort.CastNode[pkg_rhino.Node](typePos.GetItem()).GetJSType()
						}
						if sym.docScope == nil {
							{
								r_52 := sym
								r_52.docScope = NewSymbolTable_SymbolScope(pkg_rhino.Node(0), nil, nil, sym)
							}
						}
						var existingSymbol *SymbolTable_Symbol = this_.priv_SymbolTable_isAnySymbolDeclared_5227(name, paramNode, sym.docScope)
						_ = existingSymbol
						if existingSymbol == nil {
							this_.priv_SymbolTable_declareSymbol_5233(name, type_, (type_ == nil || type_.IsJavaNil_()), sym.docScope, paramNode, nil)
						}
					} else {
						param.DefineReferenceAt(paramNode)
					}
				}
			}
		}
	}
}

// Java method: SymbolTable.fillSymbolVisibility(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *SymbolTable) FillSymbolVisibility(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	var val_53 CollectFileOverviewVisibility
	InitCollectFileOverviewVisibility(&val_53, this_.compiler)
	var collectPass *CollectFileOverviewVisibility = &val_53
	_ = collectPass
	collectPass.Process_2(externs, root)
	var visibilityMap j2gort.Ext_com_google_common_collect_ImmutableMap = collectPass.GetFileOverviewVisibilityMap()
	_ = visibilityMap
	NodeTraversal_TraverseRoots(this_.compiler, NewSymbolTable_VisibilityCollector(this_, visibilityMap), externs, root)
}

// Java method: SymbolTable.createPropertyScopeFor(com.google.javascript.jscomp.SymbolTable.Symbol s)
func (this_ *SymbolTable) priv_SymbolTable_createPropertyScopeFor_5244(s *SymbolTable_Symbol) {
	_ = s
	if s.propertyScope != nil {
		return
	}
	var jsType pkg_rhino.JSType = this_.priv_SymbolTable_getType_5258(pkg_rhino.ToIface_StaticSlot(s))
	_ = jsType
	var type_ pkg_rhino.ObjectType
	_ = type_
	if jsType == nil || jsType.IsJavaNil_() {
		type_ = nil
	} else {
		type_ = jsType.GetJSTypeBase().ToObjectType()
	}
	if type_ == nil || type_.IsJavaNil_() {
		return
	}
	if (j2gort.CheckNotNullString(s.GetName()) == "exports") && s.GetDeclarationNode().IsModuleBody() {
		{
			r_54 := s
			r_54.propertyScope = this_.scopes.Get(s.GetDeclarationNode())
		}
		return
	}
	var parentPropertyScope *SymbolTable_SymbolScope = this_.priv_SymbolTable_maybeGetParentPropertyScope_5245(type_)
	_ = parentPropertyScope
	s.SetPropertyScope(NewSymbolTable_SymbolScope(pkg_rhino.Node(0), parentPropertyScope, pkg_rhino.Upcast_JSType(type_), s))
	var instanceType pkg_rhino.ObjectType = type_
	_ = instanceType
	var propNames j2gort.Ext_java_lang_Iterable = j2gort.AsIface[j2gort.Ext_java_lang_Iterable](type_.GetOwnPropertyNames())
	_ = propNames
	if instanceType.GetObjectTypeBase().IsFunctionPrototypeType() {
		if instanceType.GetOwnerFunction().HasInstanceType() {
			instanceType = instanceType.GetOwnerFunction().GetInstanceType()
			propNames = j2gort.Ext_com_google_common_collect_Iterables_Concat_2(propNames, j2gort.AsIface[j2gort.Ext_java_lang_Iterable](instanceType.GetOwnPropertyNames()))
		}
	}
	{
		cur_55 := j2gort.NewIterCursor(j2gort.CheckNotNull(propNames))
		for cur_55.HasNext() {
			propName := cur_55.NextStr()
			_ = propName
			var newProp pkg_rhino.StaticSlot = pkg_rhino.ToIface_StaticSlot(instanceType.GetObjectTypeBase().GetSlot_1_string(propName))
			_ = newProp
			if j2gort.IfaceEquals[pkg_rhino.StaticRef](j2gort.AsIface[pkg_rhino.StaticRef](newProp.GetDeclaration()), nil) {
				continue
			}
			var oldProp *SymbolTable_Symbol = this_.symbols.Get(j2gort.BoxNode(j2gort.AsIface[pkg_rhino.StaticRef](newProp.GetDeclaration()).GetNode()), j2gort.BoxString((j2gort.StrValString(s.GetName()) + "." + j2gort.StrValString(propName))))
			_ = oldProp
			if (oldProp != nil) && (this_.compiler.GetModuleMap().GetClosureModule(oldProp.GetName()) != nil) {
				continue
			}
			if this_.symbols.Get(j2gort.BoxNode(j2gort.AsIface[pkg_rhino.StaticRef](newProp.GetDeclaration()).GetNode()), j2gort.BoxString(newProp.GetName())) != nil {
				if StaticGet_SymbolTable_Logger().IsLoggable(j2gort.Ext_java_util_logging_Level_F_FINE) {
					StaticGet_SymbolTable_Logger().Fine_1_string(("Found duplicate symbol " + j2gort.StrVal(newProp)))
				}
				continue
			}
			var newSym *SymbolTable_Symbol = this_.priv_SymbolTable_copySymbolTo_5230(newProp, s.propertyScope)
			_ = newSym
			if oldProp != nil {
				if newSym.GetJSDocInfo() == nil {
					newSym.SetJSDocInfo(oldProp.GetJSDocInfo())
				}
				newSym.SetPropertyScope(oldProp.propertyScope)
				{
					cur_56 := j2gort.NewMapValCursor(oldProp.references)
					for cur_56.HasNext() {
						rawElem_57 := cur_56.Next()
						ref := j2gort.CastPtr[SymbolTable_Reference](rawElem_57)
						_ = ref
						newSym.DefineReferenceAt(ref.GetNode())
					}
				}
				this_.priv_SymbolTable_removeSymbol_5235(oldProp)
			}
		}
	}
}

// Java method: SymbolTable.maybeGetParentPropertyScope(com.google.javascript.rhino.jstype.ObjectType symbolObjectType)
func (this_ *SymbolTable) priv_SymbolTable_maybeGetParentPropertyScope_5245(symbolObjectType pkg_rhino.ObjectType) *SymbolTable_SymbolScope {
	_ = symbolObjectType
	var proto pkg_rhino.ObjectType = symbolObjectType.GetImplicitPrototype()
	_ = proto
	if (proto == nil || proto.IsJavaNil_()) || (proto == symbolObjectType || ((proto == nil || proto.IsJavaNil_()) && (symbolObjectType == nil || symbolObjectType.IsJavaNil_()))) {
		return nil
	}
	var parentSymbol *SymbolTable_Symbol = nil
	_ = parentSymbol
	if this_.priv_SymbolTable_isEs6ClassConstructor_5246(pkg_rhino.Upcast_JSType(proto)) {
		parentSymbol = this_.GetSymbolDeclaredBy_1_FunctionType(proto.ToMaybeFunctionType())
	} else if proto.GetConstructor() != nil {
		parentSymbol = this_.GetSymbolForInstancesOf_1_FunctionType(proto.GetConstructor())
	} else {
		return nil
	}
	if parentSymbol == nil {
		return nil
	}
	this_.priv_SymbolTable_createPropertyScopeFor_5244(parentSymbol)
	return parentSymbol.GetPropertyScope()
}

// Java method: SymbolTable.isEs6ClassConstructor(com.google.javascript.rhino.jstype.JSType type)
func (this_ *SymbolTable) priv_SymbolTable_isEs6ClassConstructor_5246(type_ pkg_rhino.JSType) bool {
	_ = type_
	return ((type_.GetJSTypeBase().IsFunctionType() && (type_.ToMaybeFunctionType().GetSource() != 0)) && type_.ToMaybeFunctionType().GetSource().IsClass())
}

// Java method: SymbolTable.fillThisReferences(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *SymbolTable) FillThisReferences(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NewSymbolTable_ThisRefCollector(this_).Process_2(externs, root)
}

// Java method: SymbolTable.fillSuperReferences(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *SymbolTable) FillSuperReferences(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	var collectSuper NodeTraversal_Callback = NewSymbolTable__Anon_91(this_)
	_ = collectSuper
	NodeTraversal_TraverseRoots(this_.compiler, collectSuper, externs, root)
}

// Java method: SymbolTable.addRefsInGoogRequireStatement(com.google.javascript.rhino.Node n, java.util.Map<java.lang.String,com.google.javascript.jscomp.SymbolTable.SymbolScope> googModuleScopes)
func (this_ *SymbolTable) priv_SymbolTable_addRefsInGoogRequireStatement_5249(n pkg_rhino.Node, googModuleScopes j2gort.Ext_java_util_Map) {
	_ = n
	_ = googModuleScopes
	if this_.compiler.GetOptions().ShouldRewriteModules() {
		return
	}
	var require pkg_rhino.Node = n.GetLastChild()
	_ = require
	var moduleName string = require.GetSecondChild().GetString()
	_ = moduleName
	var module *Module = this_.compiler.GetModuleMap().GetClosureModule(moduleName)
	_ = module
	var moduleType ModuleMetadataMap_ModuleType
	_ = moduleType
	if module == nil {
		moduleType = ModuleMetadataMap_ModuleType(0)
	} else {
		moduleType = module.Metadata().ModuleType()
	}
	var moduleScope *SymbolTable_SymbolScope = j2gort.CastPtr[SymbolTable_SymbolScope](j2gort.MapGetStr(googModuleScopes, moduleName))
	_ = moduleScope
	if n.IsName() {
		var declaration *SymbolTable_Symbol = this_.priv_SymbolTable_getSymbolForTypeHelper_5216(n.GetJSType(), false)
		_ = declaration
		if declaration != nil {
			declaration.DefineReferenceAt(n)
		}
		switch (moduleType).Ordinal() {
		case int32(1) /* ModuleType.GOOG_MODULE */, int32(2) /* ModuleType.LEGACY_GOOG_MODULE */ :
			declaration = moduleScope.GetOwnSlot("exports")
		case int32(4) /* ModuleType.GOOG_PROVIDE */ :
			declaration = this_.priv_SymbolTable_getSymbolForName_5214(pkg_rhino.Node(0), moduleName)
		default:
		}
		if declaration != nil {
			declaration.DefineReferenceAt(n)
		}
	} else if n.IsDestructuringLhs() {
		if !(n.GetFirstChild().IsObjectPattern()) {
			j2gort.CheckState(false)
		}
		{
			cur_58 := j2gort.NewIterCursor(j2gort.CheckNotNull(n.GetFirstChild().Children()))
			for cur_58.HasNext() {
				stringKey := pkg_rhino.Node(cur_58.NextNode())
				_ = stringKey
				if !(stringKey.IsStringKey()) {
					j2gort.CheckState(false)
				}
				var varName string = stringKey.GetString()
				_ = varName
				var varDeclaration *SymbolTable_Symbol = nil
				_ = varDeclaration
				switch (moduleType).Ordinal() {
				case int32(1) /* ModuleType.GOOG_MODULE */, int32(2) /* ModuleType.LEGACY_GOOG_MODULE */ :
					varDeclaration = moduleScope.GetOwnSlot(varName)
				case int32(4) /* ModuleType.GOOG_PROVIDE */ :
					varDeclaration = this_.priv_SymbolTable_getSymbolForName_5214(pkg_rhino.Node(0), (j2gort.StrValString(moduleName) + "." + j2gort.StrValString(varName)))
				default:
				}
				if varDeclaration != nil {
					varDeclaration.DefineReferenceAt(stringKey)
				}
			}
		}
	}
}

// Java method: SymbolTable.fillGoogProvideModuleRequires(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *SymbolTable) FillGoogProvideModuleRequires(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	var declaredNamespaces j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = declaredNamespaces
	var moduleScopes j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = moduleScopes
	var looksLikeGoogCall j2gort.Ext_java_util_function_Predicate = lambdaSingleton_symbol_table_1
	_ = looksLikeGoogCall
	var collectModuleScopes NodeTraversal_Callback = NewSymbolTable__Anon_92(this_, looksLikeGoogCall, declaredNamespaces, moduleScopes)
	_ = collectModuleScopes
	var processRequireStatements NodeTraversal_Callback = NewSymbolTable__Anon_93(this_, looksLikeGoogCall, moduleScopes, declaredNamespaces)
	_ = processRequireStatements
	NodeTraversal_TraverseRoots(this_.compiler, collectModuleScopes, externs, root)
	NodeTraversal_TraverseRoots(this_.compiler, processRequireStatements, externs, root)
}

// Java method: SymbolTable.flattenGoogModuleExports()
func (this_ *SymbolTable) FlattenGoogModuleExports() {
	{
		cur_60 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.GetAllScopes()))
		for cur_60.HasNext() {
			rawElem_61 := cur_60.Next()
			moduleScope := j2gort.CastPtr[SymbolTable_SymbolScope](rawElem_61)
			_ = moduleScope
			if !(moduleScope.IsModuleScope()) {
				continue
			}
			{
				cur_62 := j2gort.NewIterCursor(j2gort.CheckNotNull(moduleScope.ownSymbols.ValuesToImmutableList()))
				for cur_62.HasNext() {
					rawElem_63 := cur_62.Next()
					sym := j2gort.CastPtr[SymbolTable_Symbol](rawElem_63)
					_ = sym
					if javalang.StringStartsWith(j2gort.CheckNotNullString(sym.GetName()), "exports.") {
						this_.priv_SymbolTable_renameSymbol_5236(moduleScope, sym, j2gort.JStrSubFrom(sym.GetName(), j2gort.JStrLen("exports.")))
					}
				}
			}
		}
	}
}

// Java method: SymbolTable.isSymbolAQuotedObjectKey(com.google.javascript.jscomp.SymbolTable.Symbol symbol)
func (this_ *SymbolTable) priv_SymbolTable_isSymbolAQuotedObjectKey_5252(symbol *SymbolTable_Symbol) bool {
	_ = symbol
	var node pkg_rhino.Node = symbol.GetDeclarationNode()
	_ = node
	return (((node != 0) && node.IsStringKey()) && node.IsQuotedStringKey())
}

// Java method: SymbolTable.isSymbolDuplicatedExternOnWindow(com.google.javascript.jscomp.SymbolTable.Symbol symbol)
func (this_ *SymbolTable) priv_SymbolTable_isSymbolDuplicatedExternOnWindow_5253(symbol *SymbolTable_Symbol) bool {
	_ = symbol
	var node pkg_rhino.Node = symbol.GetDeclarationNode()
	_ = node
	return (((!(node.IsIndexable()) && node.IsGetProp()) && node.GetFirstChild().IsName()) && (j2gort.CheckNotNullString(node.GetFirstChild().GetString()) == "window"))
}

// Java method: SymbolTable.mergeExternSymbolsDuplicatedOnWindow()
func (this_ *SymbolTable) priv_SymbolTable_mergeExternSymbolsDuplicatedOnWindow_5254() {
	var externSymbols j2gort.Ext_com_google_common_collect_Table = j2gort.Ext_com_google_common_collect_Tables_NewCustomTable(j2gort.NewExt_java_util_LinkedHashMap(), lambdaSingleton_symbol_table_2)
	_ = externSymbols
	{
		cur_64 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(this_.symbols.Values())))
		for cur_64.HasNext() {
			rawElem_65 := cur_64.Next()
			symbol := j2gort.CastPtr[SymbolTable_Symbol](rawElem_65)
			_ = symbol
			if ((symbol.GetDeclarationNode() == 0) || j2gort.IfaceEquals[pkg_rhino.StaticSourceFile](symbol.GetDeclarationNode().GetStaticSourceFile(), nil)) || !(symbol.GetDeclarationNode().GetStaticSourceFile().IsExtern()) {
				continue
			}
			var sourceFile string = symbol.GetSourceFileName()
			_ = sourceFile
			var position int32 = symbol.GetDeclarationNode().GetSourcePosition()
			_ = position
			if !(externSymbols.Contains(j2gort.BoxString(sourceFile), j2gort.ValueOfInt(position))) {
				externSymbols.Put(j2gort.BoxString(sourceFile), j2gort.ValueOfInt(position), j2gort.PtrToAny(symbol))
				continue
			}
			var existingSymbol *SymbolTable_Symbol = j2gort.CastPtr[SymbolTable_Symbol](externSymbols.Get(j2gort.BoxString(sourceFile), j2gort.ValueOfInt(position)))
			_ = existingSymbol
			if this_.priv_SymbolTable_isSymbolDuplicatedExternOnWindow_5253(existingSymbol) {
				this_.priv_SymbolTable_mergeSymbol_5234(existingSymbol, symbol)
				externSymbols.Put(j2gort.BoxString(sourceFile), j2gort.ValueOfInt(position), j2gort.PtrToAny(symbol))
			} else if this_.priv_SymbolTable_isSymbolDuplicatedExternOnWindow_5253(symbol) {
				this_.priv_SymbolTable_mergeSymbol_5234(symbol, existingSymbol)
			}
		}
	}
}

// Java method: SymbolTable.removeGeneratedSymbols()
func (this_ *SymbolTable) RemoveGeneratedSymbols() {
	{
		cur_66 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(this_.symbols.Values())))
		for cur_66.HasNext() {
			rawElem_67 := cur_66.Next()
			symbol := j2gort.CastPtr[SymbolTable_Symbol](rawElem_67)
			_ = symbol
			if this_.priv_SymbolTable_isSymbolAQuotedObjectKey_5252(symbol) {
				var symbolAlreadyRemoved bool = !(this_.GetScope(symbol).ownSymbols.ContainsKey(symbol.GetName()))
				_ = symbolAlreadyRemoved
				if !(symbolAlreadyRemoved) {
					this_.priv_SymbolTable_removeSymbol_5235(symbol)
				}
			} else if (symbol.GetDeclarationNode() != 0) && symbol.GetDeclarationNode().IsModuleBody() {
				this_.priv_SymbolTable_removeSymbol_5235(symbol)
			}
		}
	}
	this_.priv_SymbolTable_mergeExternSymbolsDuplicatedOnWindow_5254()
}

// Java method: SymbolTable.createScopeFrom(com.google.javascript.rhino.StaticScope otherScope)
func (this_ *SymbolTable) priv_SymbolTable_createScopeFrom_5256(otherScope pkg_rhino.StaticScope) *SymbolTable_SymbolScope {
	_ = otherScope
	var otherScopeRoot pkg_rhino.Node = otherScope.GetRootNode_0()
	_ = otherScopeRoot
	var myScope *SymbolTable_SymbolScope = this_.scopes.Get(otherScopeRoot)
	_ = myScope
	if myScope == nil {
		var otherScopeParent pkg_rhino.StaticScope = j2gort.AsIface[pkg_rhino.StaticScope](otherScope.GetParentScope())
		_ = otherScopeParent
		if otherScopeParent == nil || otherScopeParent.IsJavaNil_() {
			if !(this_.globalScope == nil) {
				j2gort.CheckStateStr(false, "Global scopes found at different roots")
			}
		}
		myScope = NewSymbolTable_SymbolScope(otherScopeRoot, func() *SymbolTable_SymbolScope {
			if otherScopeParent == nil || otherScopeParent.IsJavaNil_() {
				return nil
			}
			return this_.priv_SymbolTable_createScopeFrom_5256(otherScopeParent)
		}(), this_.priv_SymbolTable_getTypeOfThis_5259(otherScope), nil)
		this_.scopes.Put(otherScopeRoot, myScope)
		if myScope.IsGlobalScope() {
			this_.globalScope = myScope
		}
	}
	return myScope
}

// Java method: SymbolTable.getLexicalScopeDepth(com.google.javascript.jscomp.SymbolTable.SymbolScope scope)
func (this_ *SymbolTable) priv_SymbolTable_getLexicalScopeDepth_5257(scope *SymbolTable_SymbolScope) int32 {
	_ = scope
	if scope.IsLexicalScope() || scope.IsDocScope() {
		return scope.GetScopeDepth()
	} else {
		if !(scope.IsPropertyScope()) {
			j2gort.CheckState(false)
		}
		var sym *SymbolTable_Symbol = scope.GetSymbolForScope()
		_ = sym
		j2gort.CheckNotNullPtr(sym)
		return (this_.priv_SymbolTable_getLexicalScopeDepth_5257(this_.GetScope(sym)) + int32(1))
	}
}

// Java method: SymbolTable.getType(com.google.javascript.rhino.StaticSlot sym)
func (this_ *SymbolTable) priv_SymbolTable_getType_5258(sym pkg_rhino.StaticSlot) pkg_rhino.JSType {
	_ = sym
	if (sym.GetJSDocInfo() != nil) && sym.GetJSDocInfo().HasTypedefType() {
		var type_ pkg_rhino.JSType = this_.registry.GetGlobalType(sym.GetName())
		_ = type_
		if type_ != nil && !type_.IsJavaNil_() {
			return type_
		}
		return sym.GetJSDocInfo().GetTypedefType().Evaluate(pkg_rhino.ToIface_StaticTypedScope(this_.compiler.GetTopScope()), this_.registry)
	}
	var staticTypedSlot pkg_rhino.StaticTypedSlot
	_ = staticTypedSlot
	var ok_68 bool
	staticTypedSlot, ok_68 = j2gort.TryCast[pkg_rhino.StaticTypedSlot](sym)
	if ok_68 {
		return staticTypedSlot.GetType()
	}
	return nil
}

// Java method: SymbolTable.getTypeOfThis(com.google.javascript.rhino.StaticScope s)
func (this_ *SymbolTable) priv_SymbolTable_getTypeOfThis_5259(s pkg_rhino.StaticScope) pkg_rhino.JSType {
	_ = s
	var staticTypedScope pkg_rhino.StaticTypedScope
	_ = staticTypedScope
	var ok_69 bool
	staticTypedScope, ok_69 = j2gort.TryCast[pkg_rhino.StaticTypedScope](s)
	if ok_69 {
		return staticTypedScope.GetTypeOfThis()
	}
	return nil
}

// Java method: SymbolTable.isTypeInferred(com.google.javascript.rhino.StaticSlot sym)
func (this_ *SymbolTable) priv_SymbolTable_isTypeInferred_5260(sym pkg_rhino.StaticSlot) bool {
	_ = sym
	var staticTypedSlot pkg_rhino.StaticTypedSlot
	_ = staticTypedSlot
	var ok_70 bool
	staticTypedSlot, ok_70 = j2gort.TryCast[pkg_rhino.StaticTypedSlot](sym)
	if ok_70 {
		return staticTypedSlot.IsTypeInferred()
	}
	return true
}

// Java class: com.google.javascript.jscomp.SymbolTable.Symbol
type SymbolTable_Symbol struct {
	pkg_rhino.SimpleSlotBase
	references    j2gort.Ext_java_util_Map       `gson:"references"`
	scope         *SymbolTable_SymbolScope       `gson:"scope"`
	propertyScope *SymbolTable_SymbolScope       `gson:"propertyScope"`
	declaration   *SymbolTable_Reference         `gson:"declaration"`
	docInfo       *pkg_rhino.JSDocInfo           `gson:"docInfo"`
	docScope      *SymbolTable_SymbolScope       `gson:"docScope"`
	visibility    pkg_rhino.JSDocInfo_Visibility `gson:"visibility"`
}

func (this_ *SymbolTable_Symbol) BindSelf_(s any) {
	_ = s
}
func (this_ *SymbolTable_Symbol) IsJavaNil_() bool { return this_ == nil }
func (*SymbolTable_Symbol) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SymbolTable$Symbol", "Symbol"
}
func (this_ *SymbolTable_Symbol) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.SimpleSlotBase.IdentityHash_()
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SymbolTable.Symbol", "com.google.javascript.jscomp.SymbolTable$Symbol", func(v_ any) bool {
		return j2gort.InstanceOf[*SymbolTable_Symbol](v_)
	})
}
func (*SymbolTable_Symbol) Implements_rhino_SimpleSlot()      {}
func (*SymbolTable_Symbol) Implements_rhino_StaticTypedSlot() {}
func (*SymbolTable_Symbol) Implements_rhino_StaticSlot()      {}

func (this_ *SymbolTable_Symbol) initFields_SymbolTable_Symbol() {
	_ = this_
	this_.references = j2gort.NewExt_java_util_LinkedHashMap()

}

// Java constructor: Symbol(java.lang.String name, com.google.javascript.rhino.jstype.JSType type, boolean inferred, com.google.javascript.jscomp.SymbolTable.SymbolScope scope)
func NewSymbolTable_Symbol(name string, type_ pkg_rhino.JSType, inferred bool, scope *SymbolTable_SymbolScope) *SymbolTable_Symbol {
	this_ := &SymbolTable_Symbol{}
	InitSymbolTable_Symbol(this_, name, type_, inferred, scope)
	return this_
}

func InitSymbolTable_Symbol(this_ *SymbolTable_Symbol, name string, type_ pkg_rhino.JSType, inferred bool, scope *SymbolTable_SymbolScope) {
	_ = name
	_ = type_
	_ = inferred
	_ = scope
	pkg_rhino.InitSimpleSlot(&this_.SimpleSlotBase, name, type_, inferred)
	this_.initFields_SymbolTable_Symbol()
	this_.scope = scope
}

// Java method: Symbol.equals(java.lang.Object o)
func (this_ *SymbolTable_Symbol) Equals(arg_71 any) bool {
	o := arg_71
	_ = o
	var other *SymbolTable_Symbol
	_ = other
	var ok_72 bool
	other, ok_72 = j2gort.TryCast[*SymbolTable_Symbol](o)
	if !ok_72 {
		return false
	}
	return ((((this_.IsTypeInferred() == other.IsTypeInferred()) && (this_.GetName() == other.GetName())) && j2gort.ObjEqIface[pkg_rhino.JSType](this_.GetType(), other.GetType())) && (this_.scope == other.scope))
}

// Java method: Symbol.hashCode()
func (this_ *SymbolTable_Symbol) HashCode() int32 {
	return func() int32 {
		var oh_73 int32 = 1
		if this_.IsTypeInferred() {
			oh_73 = 31*oh_73 + 1231
		} else {
			oh_73 = 31*oh_73 + 1237
		}
		{
			hs_74 := this_.GetName()
			oh_73 = 31 * oh_73
			if hs_74 != j2gort.NullString {
				oh_73 += javalang.StringHashCode(hs_74)
			}
		}
		{
			hi_75 := this_.GetType()
			oh_73 = 31 * oh_73
			if hi_75 != nil && !hi_75.IsJavaNil_() {
				oh_73 += hi_75.GetJSTypeBase().HashCode()
			}
		}
		{
			hp_76 := this_.scope
			oh_73 = 31 * oh_73
			if hp_76 != nil {
				oh_73 += j2gort.IdentityHash(hp_76)
			}
		}
		return oh_73
	}()
}

// Java method: Symbol.getDeclaration()
func (this_ *SymbolTable_Symbol) GetDeclaration() any {
	return j2gort.PtrToAny(this_.declaration)
}

// Java method: Symbol.getFunctionType()
func (this_ *SymbolTable_Symbol) GetFunctionType() *pkg_rhino.FunctionType {
	return pkg_rhino.JSType_ToMaybeFunctionType(this_.GetType())
}

// Java method: Symbol.defineReferenceAt(com.google.javascript.rhino.Node n)
func (this_ *SymbolTable_Symbol) DefineReferenceAt(n pkg_rhino.Node) *SymbolTable_Reference {
	_ = n
	return func() *SymbolTable_Reference {
		m_77 := this_.references
		k_78 := n
		if v_79, ok_80 := j2gort.MapTryGetComputeNode(m_77, k_78); ok_80 {
			return j2gort.CastPtr[SymbolTable_Reference](v_79)
		}
		return j2gort.CastPtr[SymbolTable_Reference](j2gort.MapComputeIfAbsentNode(m_77, k_78, j2gort.NewLambda_Ext_java_util_function_Function(func(lp_81 any) any {
			k := j2gort.CastNode[pkg_rhino.Node](lp_81)
			_ = k
			return NewSymbolTable_Reference(this_, k)
		})))
	}()
}

// Java method: Symbol.setDeclaration(com.google.javascript.jscomp.SymbolTable.Reference ref)
func (this_ *SymbolTable_Symbol) SetDeclaration(ref *SymbolTable_Reference) {
	_ = ref
	if !(this_.declaration == nil) {
		j2gort.CheckState(false)
	}
	this_.declaration = ref
}

// Java method: Symbol.getDeclarationNode()
func (this_ *SymbolTable_Symbol) GetDeclarationNode() pkg_rhino.Node {
	if this_.declaration == nil {
		return pkg_rhino.Node(0)
	} else {
		return this_.declaration.GetNode()
	}
}

// Java method: Symbol.getSourceFileName()
func (this_ *SymbolTable_Symbol) GetSourceFileName() string {
	var n pkg_rhino.Node = this_.GetDeclarationNode()
	_ = n
	if n == 0 {
		return j2gort.NullString
	} else {
		return n.GetSourceFileName()
	}
}

// Java method: Symbol.getPropertyScope()
func (this_ *SymbolTable_Symbol) GetPropertyScope() *SymbolTable_SymbolScope {
	return this_.propertyScope
}

// Java method: Symbol.setPropertyScope(com.google.javascript.jscomp.SymbolTable.SymbolScope scope)
func (this_ *SymbolTable_Symbol) SetPropertyScope(scope *SymbolTable_SymbolScope) {
	_ = scope
	this_.propertyScope = scope
	if scope != nil {
		this_.propertyScope.SetSymbolForScope(this_)
	}
}

// Java method: Symbol.getJSDocInfo()
func (this_ *SymbolTable_Symbol) GetJSDocInfo() *pkg_rhino.JSDocInfo {
	return this_.docInfo
}

// Java method: Symbol.setJSDocInfo(com.google.javascript.rhino.JSDocInfo info)
func (this_ *SymbolTable_Symbol) SetJSDocInfo(info *pkg_rhino.JSDocInfo) {
	_ = info
	this_.docInfo = info
}

// Java method: Symbol.getVisibility()
func (this_ *SymbolTable_Symbol) GetVisibility() pkg_rhino.JSDocInfo_Visibility {
	return this_.visibility
}

// Java method: Symbol.setVisibility(com.google.javascript.rhino.JSDocInfo.Visibility v)
func (this_ *SymbolTable_Symbol) SetVisibility(v pkg_rhino.JSDocInfo_Visibility) {
	_ = v
	this_.visibility = v
}

// Java method: Symbol.isProperty()
func (this_ *SymbolTable_Symbol) IsProperty() bool {
	return this_.scope.IsPropertyScope()
}

// Java method: Symbol.isLexicalVariable()
func (this_ *SymbolTable_Symbol) IsLexicalVariable() bool {
	return this_.scope.IsLexicalScope()
}

// Java method: Symbol.isDocOnlyParameter()
func (this_ *SymbolTable_Symbol) IsDocOnlyParameter() bool {
	return this_.scope.IsDocScope()
}

// Java method: Symbol.toString()
func (this_ *SymbolTable_Symbol) String() string {
	var n pkg_rhino.Node = this_.GetDeclarationNode()
	_ = n
	var lineNo int32
	_ = lineNo
	if n == 0 {
		lineNo = int32(-1)
	} else {
		lineNo = n.GetLineno()
	}
	return (j2gort.StrValString(this_.GetName()) + "@" + j2gort.StrValString(this_.GetSourceFileName()) + ":" + j2gort.StrValInt32(lineNo))
}

// Java class: com.google.javascript.jscomp.SymbolTable.Reference
type SymbolTable_Reference struct {
	pkg_rhino.SimpleReferenceBase
}

func (this_ *SymbolTable_Reference) BindSelf_(s any) {
	_ = s
}
func (this_ *SymbolTable_Reference) IsJavaNil_() bool { return this_ == nil }
func (*SymbolTable_Reference) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SymbolTable$Reference", "Reference"
}
func (this_ *SymbolTable_Reference) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.SimpleReferenceBase.IdentityHash_()
}
func (*SymbolTable_Reference) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SymbolTable.Reference", "com.google.javascript.jscomp.SymbolTable$Reference", func(v_ any) bool {
		return j2gort.InstanceOf[*SymbolTable_Reference](v_)
	})
}
func (*SymbolTable_Reference) Implements_rhino_SimpleReference() {}
func (*SymbolTable_Reference) Implements_rhino_StaticTypedRef()  {}
func (*SymbolTable_Reference) Implements_rhino_StaticRef()       {}

func (this_ *SymbolTable_Reference) initFields_SymbolTable_Reference() {
	_ = this_
}

// Java constructor: Reference(com.google.javascript.jscomp.SymbolTable.Symbol symbol, com.google.javascript.rhino.Node node)
func NewSymbolTable_Reference(symbol *SymbolTable_Symbol, node pkg_rhino.Node) *SymbolTable_Reference {
	this_ := &SymbolTable_Reference{}
	InitSymbolTable_Reference(this_, symbol, node)
	return this_
}

func InitSymbolTable_Reference(this_ *SymbolTable_Reference, symbol *SymbolTable_Symbol, node pkg_rhino.Node) {
	_ = symbol
	_ = node
	pkg_rhino.InitSimpleReference(&this_.SimpleReferenceBase, pkg_rhino.ToIface_StaticTypedSlot(symbol), node)
	this_.initFields_SymbolTable_Reference()
}

// Java class: com.google.javascript.jscomp.SymbolTable.SymbolScope
type SymbolTable_SymbolScope struct {
	parent              *SymbolTable_SymbolScope                       `gson:"parent"`
	typeOfThis          pkg_rhino.JSType                               `gson:"typeOfThis"`
	ownSymbols          j2gort.OrderedMap[string, *SymbolTable_Symbol] `gson:"ownSymbols"`
	mySymbol            *SymbolTable_Symbol                            `gson:"mySymbol"`
	rootNode            pkg_rhino.Node                                 `gson:"rootNode"`
	scopeDepth          int32                                          `gson:"scopeDepth"`
	indexInParent       int32                                          `gson:"indexInParent"`
	numberOfChildScopes int32                                          `gson:"numberOfChildScopes"`
	idHash_             atomic.Int32
}

func (this_ *SymbolTable_SymbolScope) BindSelf_(s any) {
	_ = s
}
func (this_ *SymbolTable_SymbolScope) IsJavaNil_() bool { return this_ == nil }
func (*SymbolTable_SymbolScope) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SymbolTable$SymbolScope", "SymbolScope"
}
func (this_ *SymbolTable_SymbolScope) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SymbolTable_SymbolScope) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SymbolTable.SymbolScope", "com.google.javascript.jscomp.SymbolTable$SymbolScope", func(v_ any) bool {
		return j2gort.InstanceOf[*SymbolTable_SymbolScope](v_)
	})
}

func (this_ *SymbolTable_SymbolScope) initFields_SymbolTable_SymbolScope() {
	_ = this_
	this_.ownSymbols = j2gort.NewOrderedMap[string, *SymbolTable_Symbol]()

}

// Java constructor: SymbolScope(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node rootNode, com.google.javascript.jscomp.SymbolTable.@org.jspecify.annotations.Nullable SymbolScope parent, com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable JSType typeOfThis, com.google.javascript.jscomp.SymbolTable.@org.jspecify.annotations.Nullable Symbol mySymbol)
func NewSymbolTable_SymbolScope(rootNode pkg_rhino.Node, parent *SymbolTable_SymbolScope, typeOfThis pkg_rhino.JSType, mySymbol *SymbolTable_Symbol) *SymbolTable_SymbolScope {
	this_ := &SymbolTable_SymbolScope{}
	InitSymbolTable_SymbolScope(this_, rootNode, parent, typeOfThis, mySymbol)
	return this_
}

func InitSymbolTable_SymbolScope(this_ *SymbolTable_SymbolScope, rootNode pkg_rhino.Node, parent *SymbolTable_SymbolScope, typeOfThis pkg_rhino.JSType, mySymbol *SymbolTable_Symbol) {
	_ = rootNode
	_ = parent
	_ = typeOfThis
	_ = mySymbol
	this_.initFields_SymbolTable_SymbolScope()
	this_.rootNode = rootNode
	this_.parent = parent
	this_.typeOfThis = typeOfThis
	if parent == nil {
		this_.scopeDepth = int32(0)
	} else {
		this_.scopeDepth = (parent.GetScopeDepth() + int32(1))
	}
	this_.mySymbol = mySymbol
	if parent == nil {
		this_.indexInParent = int32(0)
	} else {
		this_.indexInParent = parent.numberOfChildScopes
		r_82 := parent
		r_82.numberOfChildScopes = r_82.numberOfChildScopes + 1
	}
}

// Java method: SymbolScope.getSymbolForScope()
func (this_ *SymbolTable_SymbolScope) GetSymbolForScope() *SymbolTable_Symbol {
	return this_.mySymbol
}

// Java method: SymbolScope.setSymbolForScope(com.google.javascript.jscomp.SymbolTable.Symbol sym)
func (this_ *SymbolTable_SymbolScope) SetSymbolForScope(sym *SymbolTable_Symbol) {
	_ = sym
	this_.mySymbol = sym
}

// Java method: SymbolScope.getIndexOfSymbol(com.google.javascript.jscomp.SymbolTable.Symbol sym)
func (this_ *SymbolTable_SymbolScope) GetIndexOfSymbol(sym *SymbolTable_Symbol) int32 {
	_ = sym
	return j2gort.Ext_com_google_common_collect_Iterables_IndexOf(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.ownSymbols.LiveValues()), j2gort.Ext_com_google_common_base_Predicates_EqualTo(j2gort.PtrToAny(sym)))
}

// Java method: SymbolScope.getRootNode()
func (this_ *SymbolTable_SymbolScope) GetRootNode() pkg_rhino.Node {
	return this_.rootNode
}

// Java method: SymbolScope.getParentScope()
func (this_ *SymbolTable_SymbolScope) GetParentScope() *SymbolTable_SymbolScope {
	return this_.parent
}

// Java method: SymbolScope.getQualifiedSlot(java.lang.String name)
func (this_ *SymbolTable_SymbolScope) GetQualifiedSlot(name string) *SymbolTable_Symbol {
	_ = name
	var fullyNamedSym *SymbolTable_Symbol = this_.GetSlot(name)
	_ = fullyNamedSym
	if fullyNamedSym != nil {
		return fullyNamedSym
	}
	var dot int32 = j2gort.JStrLastIndexOfCp(name, int32(uint16(46)))
	_ = dot
	if dot != int32(-1) {
		var owner *SymbolTable_Symbol = this_.GetQualifiedSlot(j2gort.JStrSub(name, int32(0), dot))
		_ = owner
		if (owner != nil) && (owner.GetPropertyScope() != nil) {
			return owner.GetPropertyScope().GetSlot(j2gort.JStrSubFrom(name, (dot + int32(1))))
		}
	}
	return nil
}

// Java method: SymbolScope.getSlot(java.lang.String name)
func (this_ *SymbolTable_SymbolScope) GetSlot(name string) *SymbolTable_Symbol {
	_ = name
	var own *SymbolTable_Symbol = this_.GetOwnSlot(name)
	_ = own
	if own != nil {
		return own
	}
	var ancestor *SymbolTable_Symbol
	_ = ancestor
	if this_.parent == nil {
		ancestor = nil
	} else {
		ancestor = this_.parent.GetSlot(name)
	}
	if ancestor != nil {
		return ancestor
	}
	return nil
}

// Java method: SymbolScope.getOwnSlot(java.lang.String name)
func (this_ *SymbolTable_SymbolScope) GetOwnSlot(name string) *SymbolTable_Symbol {
	_ = name
	return this_.ownSymbols.Get(name)
}

// Java method: SymbolScope.getTypeOfThis()
func (this_ *SymbolTable_SymbolScope) GetTypeOfThis() pkg_rhino.JSType {
	return this_.typeOfThis
}

// Java method: SymbolScope.isGlobalScope()
func (this_ *SymbolTable_SymbolScope) IsGlobalScope() bool {
	return ((this_.GetParentScope() == nil) && (this_.GetRootNode() != 0))
}

// Java method: SymbolScope.isModuleScope()
func (this_ *SymbolTable_SymbolScope) IsModuleScope() bool {
	return ((this_.GetRootNode() != 0) && this_.GetRootNode().IsModuleBody())
}

// Java method: SymbolScope.isDocScope()
func (this_ *SymbolTable_SymbolScope) IsDocScope() bool {
	return (((this_.GetRootNode() == 0) && (this_.mySymbol != nil)) && (this_.mySymbol.docScope == this_))
}

// Java method: SymbolScope.isPropertyScope()
func (this_ *SymbolTable_SymbolScope) IsPropertyScope() bool {
	return ((this_.GetRootNode() == 0) && !(this_.IsDocScope()))
}

// Java method: SymbolScope.isLexicalScope()
func (this_ *SymbolTable_SymbolScope) IsLexicalScope() bool {
	return (this_.GetRootNode() != 0)
}

// Java method: SymbolScope.isBlockScope()
func (this_ *SymbolTable_SymbolScope) IsBlockScope() bool {
	return ((this_.GetRootNode() != 0) && NodeUtil_CreatesBlockScope(this_.GetRootNode()))
}

// Java method: SymbolScope.getScopeDepth()
func (this_ *SymbolTable_SymbolScope) GetScopeDepth() int32 {
	return this_.scopeDepth
}

// Java method: SymbolScope.getIndexInParent()
func (this_ *SymbolTable_SymbolScope) GetIndexInParent() int32 {
	return this_.indexInParent
}

// Java method: SymbolScope.toString()
func (this_ *SymbolTable_SymbolScope) String() string {
	var n pkg_rhino.Node = this_.GetRootNode()
	_ = n
	if this_.IsGlobalScope() {
		return "GlobalScope"
	} else if n != 0 {
		var type_ string = "Scope"
		_ = type_
		if this_.IsModuleScope() {
			type_ = "ModuleScope"
		} else if this_.IsBlockScope() {
			type_ = "BlockScope"
		} else if this_.IsLexicalScope() {
			type_ = "LexicalScope"
		}
		return (j2gort.StrValString(type_) + "@" + j2gort.StrValString(n.GetSourceFileName()) + ":" + j2gort.StrValInt32(n.GetLineno()))
	} else {
		return ("PropertyScope@" + j2gort.StrVal(this_.GetSymbolForScope()))
	}
}

// Java class: com.google.javascript.jscomp.SymbolTable.PropertyRefCollector
type SymbolTable_PropertyRefCollector struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *SymbolTable
}

func (this_ *SymbolTable_PropertyRefCollector) BindSelf_(s any) {
	_ = s
}
func (this_ *SymbolTable_PropertyRefCollector) IsJavaNil_() bool { return this_ == nil }
func (*SymbolTable_PropertyRefCollector) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SymbolTable$PropertyRefCollector", "PropertyRefCollector"
}
func (this_ *SymbolTable_PropertyRefCollector) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*SymbolTable_PropertyRefCollector) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SymbolTable.PropertyRefCollector", "com.google.javascript.jscomp.SymbolTable$PropertyRefCollector", func(v_ any) bool {
		return j2gort.InstanceOf[*SymbolTable_PropertyRefCollector](v_)
	})
}
func (this_ *SymbolTable_PropertyRefCollector) SetOuter_(o *SymbolTable) { this_.outer_ = o }
func (this_ *SymbolTable_PropertyRefCollector) GetOuter_() *SymbolTable  { return this_.outer_ }
func (*SymbolTable_PropertyRefCollector) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*SymbolTable_PropertyRefCollector) Implements_jscomp_CompilerPass()           {}
func (*SymbolTable_PropertyRefCollector) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *SymbolTable_PropertyRefCollector) initFields_SymbolTable_PropertyRefCollector() {
	_ = this_
}

// Java constructor: PropertyRefCollector()
func NewSymbolTable_PropertyRefCollector(outer_ *SymbolTable) *SymbolTable_PropertyRefCollector {
	this_ := &SymbolTable_PropertyRefCollector{}
	this_.outer_ = outer_
	InitSymbolTable_PropertyRefCollector(this_)
	return this_
}

func InitSymbolTable_PropertyRefCollector(this_ *SymbolTable_PropertyRefCollector) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_SymbolTable_PropertyRefCollector()
}

// Java method: PropertyRefCollector.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *SymbolTable_PropertyRefCollector) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_TraverseRoots(this_.outer_.compiler, this_, externs, root)
}

// Java method: PropertyRefCollector.maybeDefineReference(com.google.javascript.rhino.Node n, java.lang.String propName, com.google.javascript.jscomp.SymbolTable.Symbol ownerSymbol)
func (this_ *SymbolTable_PropertyRefCollector) priv_PropertyRefCollector_maybeDefineReference_5295(n pkg_rhino.Node, propName string, ownerSymbol *SymbolTable_Symbol) bool {
	_ = n
	_ = propName
	_ = ownerSymbol
	if ownerSymbol == nil {
		return false
	}
	var owner pkg_rhino.JSType = ownerSymbol.GetType()
	_ = owner
	if (owner != nil && !owner.IsJavaNil_()) && owner.GetJSTypeBase().IsObjectType() {
		var propNode pkg_rhino.Node = owner.GetJSTypeBase().ToObjectType().GetObjectTypeBase().GetPropertyNode_1_string(propName)
		_ = propNode
		var propSymbol *SymbolTable_Symbol = this_.outer_.symbols.Get(j2gort.BoxNode(propNode), j2gort.BoxString(propName))
		_ = propSymbol
		if propSymbol != nil {
			propSymbol.DefineReferenceAt(n)
			return true
		}
	}
	if ownerSymbol.GetPropertyScope() != nil {
		var prop *SymbolTable_Symbol = ownerSymbol.GetPropertyScope().GetSlot(propName)
		_ = prop
		if prop != nil {
			prop.DefineReferenceAt(n)
			return true
		}
	}
	return false
}

// Java method: PropertyRefCollector.tryDefineLexicalQualifiedNameRef(java.lang.String name, com.google.javascript.rhino.Node n)
func (this_ *SymbolTable_PropertyRefCollector) priv_PropertyRefCollector_tryDefineLexicalQualifiedNameRef_5296(name string, n pkg_rhino.Node) bool {
	_ = name
	_ = n
	if name != j2gort.NullString {
		var lexicalSym *SymbolTable_Symbol = this_.outer_.GetEnclosingScope(n).GetQualifiedSlot(name)
		_ = lexicalSym
		if lexicalSym != nil {
			lexicalSym.DefineReferenceAt(n)
			return true
		}
	}
	return false
}

// Java method: PropertyRefCollector.tryRemoveLexicalQualifiedNameRef(java.lang.String name, com.google.javascript.rhino.Node n)
func (this_ *SymbolTable_PropertyRefCollector) priv_PropertyRefCollector_tryRemoveLexicalQualifiedNameRef_5297(name string, n pkg_rhino.Node) {
	_ = name
	_ = n
	if name != j2gort.NullString {
		var lexicalSym *SymbolTable_Symbol = this_.outer_.GetEnclosingScope(n).GetQualifiedSlot(name)
		_ = lexicalSym
		if ((lexicalSym != nil) && lexicalSym.IsLexicalVariable()) && (j2gort.CastPtr[SymbolTable_Reference](lexicalSym.GetDeclaration()).GetNode() == n) {
			this_.outer_.priv_SymbolTable_removeSymbol_5235(lexicalSym)
		}
	}
}

// Java method: PropertyRefCollector.maybeDefineTypedReference(com.google.javascript.rhino.Node n, java.lang.String propName, com.google.javascript.rhino.jstype.JSType owner)
func (this_ *SymbolTable_PropertyRefCollector) priv_PropertyRefCollector_maybeDefineTypedReference_5298(n pkg_rhino.Node, propName string, owner pkg_rhino.JSType) bool {
	_ = n
	_ = propName
	_ = owner
	base_owner_ := owner.GetJSTypeBase()
	if base_owner_.IsGlobalThisType() {
		var sym *SymbolTable_Symbol = this_.outer_.globalScope.GetSlot(propName)
		_ = sym
		if sym != nil {
			sym.DefineReferenceAt(n)
			return true
		}
	} else if base_owner_.IsNominalConstructorOrInterface() {
		return this_.priv_PropertyRefCollector_maybeDefineReference_5295(n, propName, this_.outer_.GetSymbolDeclaredBy_1_FunctionType(owner.ToMaybeFunctionType()))
	} else if base_owner_.IsEnumType() {
		return this_.priv_PropertyRefCollector_maybeDefineReference_5295(n, propName, this_.outer_.GetSymbolDeclaredBy_1_EnumType(owner.ToMaybeEnumType()))
	} else {
		var defined bool = false
		_ = defined
		{
			cur_83 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.outer_.GetAllSymbolsForType(owner)))
			for cur_83.HasNext() {
				rawElem_84 := cur_83.Next()
				ctor := j2gort.CastPtr[SymbolTable_Symbol](rawElem_84)
				_ = ctor
				if ctor.GetPropertyScope() != nil {
					var ctorType pkg_rhino.JSType = ctor.GetType()
					_ = ctorType
					var isConstructor bool = ((ctorType != nil && !ctorType.IsJavaNil_()) && (ctorType.GetJSTypeBase().IsConstructor() || ctorType.GetJSTypeBase().IsInterface()))
					_ = isConstructor
					if !(isConstructor) && this_.priv_PropertyRefCollector_maybeDefineReference_5295(n, propName, ctor) {
						defined = true
					}
				}
				if this_.priv_PropertyRefCollector_maybeDefineReference_5295(n, propName, this_.outer_.GetSymbolForInstancesOf_1_SymbolTable_Symbol(ctor)) {
					defined = true
				}
			}
		}
		if defined {
			return true
		}
	}
	if base_owner_.IsObjectType() && !(j2gort.IfaceEquals[pkg_rhino.ObjectType](base_owner_.ToMaybeObjectType(), nil)) {
		var propNodeDecl pkg_rhino.Node = base_owner_.AssertObjectType().GetObjectTypeBase().GetPropertyNode_1_string(propName)
		_ = propNodeDecl
		if propNodeDecl == 0 {
			return false
		}
		if (propNodeDecl != 0) && (propNodeDecl.IsGetProp() || propNodeDecl.IsOptChainGetProp()) {
			var namespace pkg_rhino.Node = propNodeDecl.GetFirstChild()
			_ = namespace
			var namespaceSym *SymbolTable_Symbol = this_.outer_.priv_SymbolTable_getSymbolForName_5214(namespace, namespace.GetQualifiedName())
			_ = namespaceSym
			if (namespaceSym != nil) && (namespaceSym.GetPropertyScope() != nil) {
				var propSym *SymbolTable_Symbol = namespaceSym.GetPropertyScope().GetSlot(propName)
				_ = propSym
				if (propSym != nil) && (propSym.GetDeclarationNode() != n) {
					propSym.DefineReferenceAt(n)
					return true
				}
			}
		}
		var propNodeDeclScope *SymbolTable_SymbolScope = this_.outer_.GetEnclosingScope(propNodeDecl)
		_ = propNodeDeclScope
		if (propNodeDeclScope == nil) || !(propNodeDeclScope.IsModuleScope()) {
			return false
		}
		var symbol *SymbolTable_Symbol = this_.outer_.priv_SymbolTable_getSymbolForName_5214(propNodeDecl, propName)
		_ = symbol
		if (symbol != nil) && (symbol.GetDeclarationNode() != n) {
			symbol.DefineReferenceAt(n)
			return true
		}
	}
	return false
}

// Java method: PropertyRefCollector.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *SymbolTable_PropertyRefCollector) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsGetProp() || n.IsOptChainGetProp() {
		var owner pkg_rhino.JSType = n.GetFirstChild().GetJSType()
		_ = owner
		if owner != nil && !owner.IsJavaNil_() {
			var defined bool = this_.priv_PropertyRefCollector_maybeDefineTypedReference_5298(n, n.GetString(), owner)
			_ = defined
			if defined {
				this_.priv_PropertyRefCollector_tryRemoveLexicalQualifiedNameRef_5297(n.GetQualifiedName(), n)
				return
			}
		}
		this_.priv_PropertyRefCollector_tryDefineLexicalQualifiedNameRef_5296(n.GetQualifiedName(), n)
	} else if n.IsStringKey() {
		var owner_2 pkg_rhino.JSType = parent.GetJSType()
		_ = owner_2
		if owner_2 != nil && !owner_2.IsJavaNil_() {
			var defined_2 bool = this_.priv_PropertyRefCollector_maybeDefineTypedReference_5298(n, n.GetString(), owner_2)
			_ = defined_2
			if defined_2 {
				this_.priv_PropertyRefCollector_tryRemoveLexicalQualifiedNameRef_5297(NodeUtil_GetBestLValueName(n), n)
				return
			}
		}
		this_.priv_PropertyRefCollector_tryDefineLexicalQualifiedNameRef_5296(NodeUtil_GetBestLValueName(n), n)
	}
}

// Java class: com.google.javascript.jscomp.SymbolTable.ThisRefCollector
type SymbolTable_ThisRefCollector struct {
	NodeTraversal_AbstractScopedCallbackBase
	outer_    *SymbolTable
	thisStack []*SymbolTable_Symbol `gson:"thisStack"`
}

func (this_ *SymbolTable_ThisRefCollector) BindSelf_(s any) {
	_ = s
}
func (this_ *SymbolTable_ThisRefCollector) IsJavaNil_() bool { return this_ == nil }
func (*SymbolTable_ThisRefCollector) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SymbolTable$ThisRefCollector", "ThisRefCollector"
}
func (this_ *SymbolTable_ThisRefCollector) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractScopedCallbackBase.IdentityHash_()
}
func (*SymbolTable_ThisRefCollector) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SymbolTable.ThisRefCollector", "com.google.javascript.jscomp.SymbolTable$ThisRefCollector", func(v_ any) bool {
		return j2gort.InstanceOf[*SymbolTable_ThisRefCollector](v_)
	})
}
func (this_ *SymbolTable_ThisRefCollector) SetOuter_(o *SymbolTable)                          { this_.outer_ = o }
func (this_ *SymbolTable_ThisRefCollector) GetOuter_() *SymbolTable                           { return this_.outer_ }
func (*SymbolTable_ThisRefCollector) Implements_jscomp_NodeTraversal_AbstractScopedCallback() {}
func (*SymbolTable_ThisRefCollector) Implements_jscomp_CompilerPass()                         {}
func (*SymbolTable_ThisRefCollector) Implements_jscomp_NodeTraversal_ScopedCallback()         {}
func (*SymbolTable_ThisRefCollector) Implements_jscomp_NodeTraversal_Callback()               {}

func (this_ *SymbolTable_ThisRefCollector) initFields_SymbolTable_ThisRefCollector() {
	_ = this_
	this_.thisStack = []*SymbolTable_Symbol{}
}

// Java constructor: ThisRefCollector()
func NewSymbolTable_ThisRefCollector(outer_ *SymbolTable) *SymbolTable_ThisRefCollector {
	this_ := &SymbolTable_ThisRefCollector{}
	this_.outer_ = outer_
	InitSymbolTable_ThisRefCollector(this_)
	return this_
}

func InitSymbolTable_ThisRefCollector(this_ *SymbolTable_ThisRefCollector) {
	InitNodeTraversal_AbstractScopedCallback(&this_.NodeTraversal_AbstractScopedCallbackBase)
	this_.initFields_SymbolTable_ThisRefCollector()
}

// Java method: ThisRefCollector.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *SymbolTable_ThisRefCollector) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_TraverseRoots(this_.outer_.compiler, this_, externs, root)
}

// Java method: ThisRefCollector.enterScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *SymbolTable_ThisRefCollector) EnterScope(t *NodeTraversal) {
	_ = t
	var symbol *SymbolTable_Symbol = nil
	_ = symbol
	if t.InGlobalScope() {
		var firstInputRoot pkg_rhino.Node = t.GetScopeRoot().GetLastChild().GetFirstChild()
		_ = firstInputRoot
		if firstInputRoot != 0 {
			symbol = this_.outer_.priv_SymbolTable_addSymbol_5232("*global*", this_.outer_.registry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_GLOBAL_THIS()), false, this_.outer_.globalScope, firstInputRoot)
			symbol.SetDeclaration(NewSymbolTable_Reference(symbol, firstInputRoot))
		}
		this_.thisStack = append(this_.thisStack, symbol)
	} else if t.GetScopeRoot().IsFunction() {
		var scopeRoot pkg_rhino.Node = t.GetScopeRoot()
		_ = scopeRoot
		var scope *SymbolTable_SymbolScope = this_.outer_.scopes.Get(scopeRoot)
		_ = scope
		if NodeUtil_GetFunctionBody(scopeRoot).HasChildren() {
			var scopeSymbol *SymbolTable_Symbol = this_.outer_.GetSymbolForScope(scope)
			_ = scopeSymbol
			if scopeSymbol != nil {
				var propScope *SymbolTable_SymbolScope = scopeSymbol.GetPropertyScope()
				_ = propScope
				if propScope != nil {
					symbol = propScope.GetOwnSlot("this")
					if symbol == nil {
						var rootType pkg_rhino.JSType = t.GetScopeRoot().GetJSType()
						_ = rootType
						var fnType *pkg_rhino.FunctionType
						_ = fnType
						if rootType == nil || rootType.IsJavaNil_() {
							fnType = nil
						} else {
							fnType = rootType.ToMaybeFunctionType()
						}
						var type_ pkg_rhino.JSType
						_ = type_
						if fnType == nil {
							type_ = nil
						} else {
							type_ = fnType.GetTypeOfThis()
						}
						symbol = this_.outer_.priv_SymbolTable_addSymbol_5232("this", type_, false, scope, scopeRoot)
					}
				}
			}
		} else {
			StaticGet_SymbolTable_Logger().Fine_1_string(("Skipping empty function: " + j2gort.StrVal(scopeRoot)))
		}
		this_.thisStack = append(this_.thisStack, symbol)
	}
}

// Java method: ThisRefCollector.exitScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *SymbolTable_ThisRefCollector) ExitScope(t *NodeTraversal) {
	_ = t
	if t.InGlobalScope() || t.GetScopeRoot().IsFunction() {
		this_.thisStack = j2gort.SliceRemoveAt(this_.thisStack, (int32(len(this_.thisStack)) - int32(1)))
	}
}

// Java method: ThisRefCollector.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *SymbolTable_ThisRefCollector) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(n.IsThis()) {
		return
	}
	var symbol *SymbolTable_Symbol = j2gort.SliceGetLast(this_.thisStack)
	_ = symbol
	if symbol != nil {
		var ref *SymbolTable_Reference = symbol.DefineReferenceAt(n)
		_ = ref
		if j2gort.CastPtr[SymbolTable_Reference](symbol.GetDeclaration()) == nil {
			symbol.SetDeclaration(ref)
		}
	}
}

// Java class: com.google.javascript.jscomp.SymbolTable.JSDocInfoCollector
type SymbolTable_JSDocInfoCollector struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_       *SymbolTable
	typeRegistry *pkg_rhino.JSTypeRegistry `gson:"typeRegistry"`
}

func (this_ *SymbolTable_JSDocInfoCollector) BindSelf_(s any) {
	_ = s
}
func (this_ *SymbolTable_JSDocInfoCollector) IsJavaNil_() bool { return this_ == nil }
func (*SymbolTable_JSDocInfoCollector) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SymbolTable$JSDocInfoCollector", "JSDocInfoCollector"
}
func (this_ *SymbolTable_JSDocInfoCollector) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*SymbolTable_JSDocInfoCollector) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SymbolTable.JSDocInfoCollector", "com.google.javascript.jscomp.SymbolTable$JSDocInfoCollector", func(v_ any) bool {
		return j2gort.InstanceOf[*SymbolTable_JSDocInfoCollector](v_)
	})
}
func (this_ *SymbolTable_JSDocInfoCollector) SetOuter_(o *SymbolTable)                             { this_.outer_ = o }
func (this_ *SymbolTable_JSDocInfoCollector) GetOuter_() *SymbolTable                              { return this_.outer_ }
func (*SymbolTable_JSDocInfoCollector) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*SymbolTable_JSDocInfoCollector) Implements_jscomp_NodeTraversal_Callback()                  {}

func (this_ *SymbolTable_JSDocInfoCollector) initFields_SymbolTable_JSDocInfoCollector() {
	_ = this_
}

// Java constructor: JSDocInfoCollector(com.google.javascript.rhino.jstype.JSTypeRegistry registry)
func NewSymbolTable_JSDocInfoCollector(outer_ *SymbolTable, registry *pkg_rhino.JSTypeRegistry) *SymbolTable_JSDocInfoCollector {
	this_ := &SymbolTable_JSDocInfoCollector{}
	this_.outer_ = outer_
	InitSymbolTable_JSDocInfoCollector(this_, registry)
	return this_
}

func InitSymbolTable_JSDocInfoCollector(this_ *SymbolTable_JSDocInfoCollector, registry *pkg_rhino.JSTypeRegistry) {
	_ = registry
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_SymbolTable_JSDocInfoCollector()
	this_.typeRegistry = registry
}

// Java method: JSDocInfoCollector.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *SymbolTable_JSDocInfoCollector) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var info *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = info
	if info == nil {
		return
	}
	this_.outer_.docInfos = append(this_.outer_.docInfos, n)
	{
		var i int32 = int32(0)
		_ = i
		var count int32 = info.GetTypeNodeCount()
		_ = count
		for ; i < count; i = i + 1 {
			var scope *SymbolTable_SymbolScope = this_.outer_.scopes.Get(t.GetScopeRoot())
			_ = scope
			this_.VisitTypeNode(info.GetTemplateTypes(), func() *SymbolTable_SymbolScope {
				if scope == nil {
					return this_.outer_.globalScope
				}
				return scope
			}(), info.GetTypeNodeAt(i))
		}
	}
}

// Java method: JSDocInfoCollector.isNativeSourcelessType(java.lang.String name)
func (this_ *SymbolTable_JSDocInfoCollector) priv_JSDocInfoCollector_isNativeSourcelessType_5299(name string) bool {
	_ = name
	switch j2gort.CheckNotNullString(name) {
	case "null", "undefined", "void":
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: JSDocInfoCollector.visitTypeNode(com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.rhino.JSTypeExpression> templateTypeNames, com.google.javascript.jscomp.SymbolTable.SymbolScope scope, com.google.javascript.rhino.Node n)
func (this_ *SymbolTable_JSDocInfoCollector) VisitTypeNode(templateTypeNames j2gort.Ext_com_google_common_collect_ImmutableMap, scope *SymbolTable_SymbolScope, n pkg_rhino.Node) {
	_ = templateTypeNames
	_ = scope
	_ = n
	if (n.IsStringLit() && !(this_.priv_JSDocInfoCollector_isNativeSourcelessType_5299(n.GetString()))) && !(j2gort.MapContainsKeyStr(templateTypeNames, n.GetString())) {
		var symbol *SymbolTable_Symbol = this_.priv_JSDocInfoCollector_lookupPossiblyDottedName_5301(scope, n.GetString())
		_ = symbol
		if symbol != nil {
			var ref pkg_rhino.Node = n
			_ = ref
			var typeString string
			_ = typeString
			if n.GetOriginalName() != j2gort.NullString {
				typeString = n.GetOriginalName()
			} else {
				typeString = n.GetString()
			}
			if j2gort.JStrIndexOf(typeString, ".") >= 0 {
				var lastPart string = j2gort.JStrSubFrom(typeString, (j2gort.JStrLastIndexOfCp(typeString, int32(uint16(46))) + int32(1)))
				_ = lastPart
				var copy_ pkg_rhino.Node = n.CloneNode_0()
				_ = copy_
				copy_.SetLinenoCharno(copy_.GetLineno(), ((copy_.GetCharno() + copy_.GetLength()) - j2gort.JStrLen(lastPart)))
				copy_.SetLength(j2gort.JStrLen(lastPart))
				ref = copy_
			}
			symbol.DefineReferenceAt(ref)
		}
	}
	{
		var child pkg_rhino.Node = n.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			this_.VisitTypeNode(templateTypeNames, scope, child)
		}
	}
}

// Java method: JSDocInfoCollector.lookupPossiblyDottedName(com.google.javascript.jscomp.SymbolTable.SymbolScope scope, java.lang.String dottedName)
func (this_ *SymbolTable_JSDocInfoCollector) priv_JSDocInfoCollector_lookupPossiblyDottedName_5301(scope *SymbolTable_SymbolScope, dottedName string) *SymbolTable_Symbol {
	_ = scope
	_ = dottedName
	var result *SymbolTable_Symbol = scope.GetQualifiedSlot(dottedName)
	_ = result
	if result != nil {
		return result
	}
	var type_ pkg_rhino.JSType = this_.typeRegistry.GetGlobalType(dottedName)
	_ = type_
	var autobox pkg_rhino.JSType
	_ = autobox
	if type_ == nil || type_.IsJavaNil_() {
		autobox = nil
	} else {
		autobox = type_.AutoboxesTo()
	}
	if autobox == nil || autobox.IsJavaNil_() {
		return this_.outer_.priv_SymbolTable_getSymbolForTypeHelper_5216(type_, true)
	} else {
		return this_.outer_.priv_SymbolTable_getSymbolForTypeHelper_5216(autobox, true)
	}
}

// Java class: com.google.javascript.jscomp.SymbolTable.VisibilityCollector
type SymbolTable_VisibilityCollector struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_            *SymbolTable
	fileVisibilityMap j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"fileVisibilityMap"`
}

func (this_ *SymbolTable_VisibilityCollector) BindSelf_(s any) {
	_ = s
}
func (this_ *SymbolTable_VisibilityCollector) IsJavaNil_() bool { return this_ == nil }
func (*SymbolTable_VisibilityCollector) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SymbolTable$VisibilityCollector", "VisibilityCollector"
}
func (this_ *SymbolTable_VisibilityCollector) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*SymbolTable_VisibilityCollector) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SymbolTable.VisibilityCollector", "com.google.javascript.jscomp.SymbolTable$VisibilityCollector", func(v_ any) bool {
		return j2gort.InstanceOf[*SymbolTable_VisibilityCollector](v_)
	})
}
func (this_ *SymbolTable_VisibilityCollector) SetOuter_(o *SymbolTable)                             { this_.outer_ = o }
func (this_ *SymbolTable_VisibilityCollector) GetOuter_() *SymbolTable                              { return this_.outer_ }
func (*SymbolTable_VisibilityCollector) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*SymbolTable_VisibilityCollector) Implements_jscomp_NodeTraversal_Callback()                  {}

func (this_ *SymbolTable_VisibilityCollector) initFields_SymbolTable_VisibilityCollector() {
	_ = this_
}

// Java constructor: VisibilityCollector(com.google.common.collect.ImmutableMap<com.google.javascript.rhino.StaticSourceFile,com.google.javascript.rhino.JSDocInfo.Visibility> fileVisibilityMap)
func NewSymbolTable_VisibilityCollector(outer_ *SymbolTable, fileVisibilityMap j2gort.Ext_com_google_common_collect_ImmutableMap) *SymbolTable_VisibilityCollector {
	this_ := &SymbolTable_VisibilityCollector{}
	this_.outer_ = outer_
	InitSymbolTable_VisibilityCollector(this_, fileVisibilityMap)
	return this_
}

func InitSymbolTable_VisibilityCollector(this_ *SymbolTable_VisibilityCollector, fileVisibilityMap j2gort.Ext_com_google_common_collect_ImmutableMap) {
	_ = fileVisibilityMap
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_SymbolTable_VisibilityCollector()
	this_.fileVisibilityMap = fileVisibilityMap
}

// Java method: VisibilityCollector.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *SymbolTable_VisibilityCollector) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsName() {
		this_.priv_VisibilityCollector_visitName_5302(t, n)
	} else if n.IsGetProp() || n.IsOptChainGetProp() {
		this_.priv_VisibilityCollector_visitProperty_5303(n, parent)
	}
}

// Java method: VisibilityCollector.visitName(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *SymbolTable_VisibilityCollector) priv_VisibilityCollector_visitName_5302(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	var symbol *SymbolTable_Symbol = this_.outer_.symbols.Get(j2gort.BoxNode(n), j2gort.BoxString(n.GetString()))
	_ = symbol
	if symbol == nil {
		return
	}
	if symbol.GetVisibility() != 0 {
		return
	}
	var var_ *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(n.GetString()))
	_ = var_
	if var_ == nil {
		return
	}
	var v pkg_rhino.JSDocInfo_Visibility = AccessControlUtils_GetEffectiveNameVisibility(n, var_, this_.fileVisibilityMap)
	_ = v
	if v == 0 {
		return
	}
	symbol.SetVisibility(v)
}

// Java method: VisibilityCollector.visitProperty(com.google.javascript.rhino.Node getprop, com.google.javascript.rhino.Node parent)
func (this_ *SymbolTable_VisibilityCollector) priv_VisibilityCollector_visitProperty_5303(getprop pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = getprop
	_ = parent
	var propertyName string = getprop.GetString()
	_ = propertyName
	var symbol *SymbolTable_Symbol = this_.outer_.symbols.Get(j2gort.BoxNode(getprop), j2gort.BoxString(propertyName))
	_ = symbol
	if symbol == nil {
		return
	}
	if symbol.GetVisibility() != 0 {
		return
	}
	var jsType pkg_rhino.JSType = getprop.GetFirstChild().GetJSType()
	_ = jsType
	if jsType == nil || jsType.IsJavaNil_() {
		return
	}
	var isOverride bool = (((parent.GetJSDocInfo() != nil) && parent.IsAssign()) && (parent.GetFirstChild() == getprop))
	_ = isOverride
	if isOverride {
		symbol.SetVisibility(pkg_rhino.StaticGet_JSDocInfo_Visibility_INHERITED())
	} else {
		var referenceType pkg_rhino.ObjectType = pkg_rhino.ObjectType_Cast(pkg_rhino.Upcast_JSType(jsType.GetJSTypeBase().Dereference()))
		_ = referenceType
		var v pkg_rhino.JSDocInfo_Visibility = AccessControlUtils_GetEffectivePropertyVisibility(getprop, referenceType, this_.fileVisibilityMap)
		_ = v
		if v == 0 {
			return
		}
		symbol.SetVisibility(v)
	}
}

// Java class:
type SymbolTable__Anon_94 struct {
	j2gort.Ext_com_google_common_collect_OrderingBase
	outer_  *SymbolTable
	idHash_ atomic.Int32
}

func (this_ *SymbolTable__Anon_94) BindSelf_(s any) {
	_ = s
	this_.Ext_com_google_common_collect_OrderingBase.BindExtSelf_(s)
}
func (this_ *SymbolTable__Anon_94) IsJavaNil_() bool { return this_ == nil }
func (*SymbolTable__Anon_94) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SymbolTable$5", ""
}
func (this_ *SymbolTable__Anon_94) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (this_ *SymbolTable__Anon_94) SetOuter_(o *SymbolTable)                            { this_.outer_ = o }
func (this_ *SymbolTable__Anon_94) GetOuter_() *SymbolTable                             { return this_.outer_ }
func (*SymbolTable__Anon_94) Implements_j2gort_Ext_com_google_common_collect_Ordering() {}
func (*SymbolTable__Anon_94) Implements_j2gort_Ext_java_util_Comparator()               {}

func (this_ *SymbolTable__Anon_94) initFields_SymbolTable__Anon_94() {
	_ = this_
}

func NewSymbolTable__Anon_94(outer_ *SymbolTable) *SymbolTable__Anon_94 {
	this_ := &SymbolTable__Anon_94{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	j2gort.InitExt_com_google_common_collect_Ordering(&this_.Ext_com_google_common_collect_OrderingBase)
	this_.initFields_SymbolTable__Anon_94()
	return this_
}

// Java method: .compare(com.google.javascript.rhino.Node a, com.google.javascript.rhino.Node b)
func (this_ *SymbolTable__Anon_94) Compare(arg_85 any, arg_86 any) int32 {
	a := j2gort.CastNode[pkg_rhino.Node](arg_85)
	_ = a
	b := j2gort.CastNode[pkg_rhino.Node](arg_86)
	_ = b
	var result int32 = this_.outer_.sourceNameOrdering.Compare(j2gort.BoxString(a.GetSourceFileName()), j2gort.BoxString(b.GetSourceFileName()))
	_ = result
	if result != int32(0) {
		return result
	}
	return (a.GetSourcePosition() - b.GetSourcePosition())
}

// Java class:
type SymbolTable__Anon_95 struct {
	j2gort.Ext_com_google_common_collect_OrderingBase
	outer_  *SymbolTable
	idHash_ atomic.Int32
}

func (this_ *SymbolTable__Anon_95) BindSelf_(s any) {
	_ = s
	this_.Ext_com_google_common_collect_OrderingBase.BindExtSelf_(s)
}
func (this_ *SymbolTable__Anon_95) IsJavaNil_() bool { return this_ == nil }
func (*SymbolTable__Anon_95) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SymbolTable$6", ""
}
func (this_ *SymbolTable__Anon_95) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (this_ *SymbolTable__Anon_95) SetOuter_(o *SymbolTable)                            { this_.outer_ = o }
func (this_ *SymbolTable__Anon_95) GetOuter_() *SymbolTable                             { return this_.outer_ }
func (*SymbolTable__Anon_95) Implements_j2gort_Ext_com_google_common_collect_Ordering() {}
func (*SymbolTable__Anon_95) Implements_j2gort_Ext_java_util_Comparator()               {}

func (this_ *SymbolTable__Anon_95) initFields_SymbolTable__Anon_95() {
	_ = this_
}

func NewSymbolTable__Anon_95(outer_ *SymbolTable) *SymbolTable__Anon_95 {
	this_ := &SymbolTable__Anon_95{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	j2gort.InitExt_com_google_common_collect_Ordering(&this_.Ext_com_google_common_collect_OrderingBase)
	this_.initFields_SymbolTable__Anon_95()
	return this_
}

// Java method: .compare(com.google.javascript.jscomp.SymbolTable.SymbolScope a, com.google.javascript.jscomp.SymbolTable.SymbolScope b)
func (this_ *SymbolTable__Anon_95) Compare(arg_87 any, arg_88 any) int32 {
	a := j2gort.CastPtr[SymbolTable_SymbolScope](arg_87)
	_ = a
	b := j2gort.CastPtr[SymbolTable_SymbolScope](arg_88)
	_ = b
	if !(a.IsLexicalScope() && b.IsLexicalScope()) {
		j2gort.CheckStateStr(false, "We can only sort lexical scopes")
	}
	var result int32 = this_.outer_.nodeOrdering.Compare(j2gort.BoxNode(a.GetRootNode()), j2gort.BoxNode(b.GetRootNode()))
	_ = result
	if result != int32(0) {
		return result
	}
	return (a.GetScopeDepth() - b.GetScopeDepth())
}

// Java class:
type SymbolTable__Anon_96 struct {
	j2gort.Ext_com_google_common_collect_OrderingBase
	outer_  *SymbolTable
	idHash_ atomic.Int32
}

func (this_ *SymbolTable__Anon_96) BindSelf_(s any) {
	_ = s
	this_.Ext_com_google_common_collect_OrderingBase.BindExtSelf_(s)
}
func (this_ *SymbolTable__Anon_96) IsJavaNil_() bool { return this_ == nil }
func (*SymbolTable__Anon_96) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SymbolTable$7", ""
}
func (this_ *SymbolTable__Anon_96) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (this_ *SymbolTable__Anon_96) SetOuter_(o *SymbolTable)                            { this_.outer_ = o }
func (this_ *SymbolTable__Anon_96) GetOuter_() *SymbolTable                             { return this_.outer_ }
func (*SymbolTable__Anon_96) Implements_j2gort_Ext_com_google_common_collect_Ordering() {}
func (*SymbolTable__Anon_96) Implements_j2gort_Ext_java_util_Comparator()               {}

func (this_ *SymbolTable__Anon_96) initFields_SymbolTable__Anon_96() {
	_ = this_
}

func NewSymbolTable__Anon_96(outer_ *SymbolTable) *SymbolTable__Anon_96 {
	this_ := &SymbolTable__Anon_96{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	j2gort.InitExt_com_google_common_collect_Ordering(&this_.Ext_com_google_common_collect_OrderingBase)
	this_.initFields_SymbolTable__Anon_96()
	return this_
}

// Java method: .compare(com.google.javascript.jscomp.SymbolTable.Symbol a, com.google.javascript.jscomp.SymbolTable.Symbol b)
func (this_ *SymbolTable__Anon_96) Compare(arg_89 any, arg_90 any) int32 {
	a := j2gort.CastPtr[SymbolTable_Symbol](arg_89)
	_ = a
	b := j2gort.CastPtr[SymbolTable_Symbol](arg_90)
	_ = b
	var scopeA *SymbolTable_SymbolScope = this_.outer_.GetScope(a)
	_ = scopeA
	var scopeB *SymbolTable_SymbolScope = this_.outer_.GetScope(b)
	_ = scopeB
	var result int32 = (this_.outer_.priv_SymbolTable_getLexicalScopeDepth_5257(scopeA) - this_.outer_.priv_SymbolTable_getLexicalScopeDepth_5257(scopeB))
	_ = result
	if result != int32(0) {
		return result
	}
	result = j2gort.JStrCompareTo(a.GetName(), b.GetName())
	if result != int32(0) {
		return result
	}
	if scopeA.IsLexicalScope() && scopeB.IsLexicalScope() {
		return this_.outer_.lexicalScopeOrdering.Compare(j2gort.PtrToAny(scopeA), j2gort.PtrToAny(scopeB))
	}
	return this_.outer_.nodeOrdering.NullsFirst().Compare(j2gort.BoxNode(a.GetDeclarationNode()), j2gort.BoxNode(b.GetDeclarationNode()))
}

// Java class:
type SymbolTable__Anon_90 struct {
	NodeTraversal_AbstractScopedCallbackBase
	outer_ *SymbolTable
}

func (this_ *SymbolTable__Anon_90) BindSelf_(s any) {
	_ = s
}
func (this_ *SymbolTable__Anon_90) IsJavaNil_() bool { return this_ == nil }
func (*SymbolTable__Anon_90) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SymbolTable$1", ""
}
func (this_ *SymbolTable__Anon_90) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractScopedCallbackBase.IdentityHash_()
}
func (*SymbolTable__Anon_90) HasIdentityEquals_()                                     {}
func (this_ *SymbolTable__Anon_90) SetOuter_(o *SymbolTable)                          { this_.outer_ = o }
func (this_ *SymbolTable__Anon_90) GetOuter_() *SymbolTable                           { return this_.outer_ }
func (*SymbolTable__Anon_90) Implements_jscomp_NodeTraversal_AbstractScopedCallback() {}
func (*SymbolTable__Anon_90) Implements_jscomp_NodeTraversal_ScopedCallback()         {}
func (*SymbolTable__Anon_90) Implements_jscomp_NodeTraversal_Callback()               {}

func (this_ *SymbolTable__Anon_90) initFields_SymbolTable__Anon_90() {
	_ = this_
}

func NewSymbolTable__Anon_90(outer_ *SymbolTable) *SymbolTable__Anon_90 {
	this_ := &SymbolTable__Anon_90{}
	this_.outer_ = outer_
	InitNodeTraversal_AbstractScopedCallback(&this_.NodeTraversal_AbstractScopedCallbackBase)
	this_.initFields_SymbolTable__Anon_90()
	return this_
}

// Java method: .enterScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *SymbolTable__Anon_90) EnterScope(t *NodeTraversal) {
	_ = t
	this_.outer_.priv_SymbolTable_createScopeFrom_5256(pkg_rhino.ToIface_StaticScope(t.GetScope()))
}

// Java method: .visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node p)
func (this_ *SymbolTable__Anon_90) Visit_3(t *NodeTraversal, n pkg_rhino.Node, p pkg_rhino.Node) {
	_ = t
	_ = n
	_ = p
}

// Java class:
type SymbolTable__Anon_91 struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *SymbolTable
}

func (this_ *SymbolTable__Anon_91) BindSelf_(s any) {
	_ = s
}
func (this_ *SymbolTable__Anon_91) IsJavaNil_() bool { return this_ == nil }
func (*SymbolTable__Anon_91) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SymbolTable$2", ""
}
func (this_ *SymbolTable__Anon_91) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*SymbolTable__Anon_91) HasIdentityEquals_()                                        {}
func (this_ *SymbolTable__Anon_91) SetOuter_(o *SymbolTable)                             { this_.outer_ = o }
func (this_ *SymbolTable__Anon_91) GetOuter_() *SymbolTable                              { return this_.outer_ }
func (*SymbolTable__Anon_91) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*SymbolTable__Anon_91) Implements_jscomp_NodeTraversal_Callback()                  {}

func (this_ *SymbolTable__Anon_91) initFields_SymbolTable__Anon_91() {
	_ = this_
}

func NewSymbolTable__Anon_91(outer_ *SymbolTable) *SymbolTable__Anon_91 {
	this_ := &SymbolTable__Anon_91{}
	this_.outer_ = outer_
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_SymbolTable__Anon_91()
	return this_
}

// Java method: .visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *SymbolTable__Anon_91) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(n.IsSuper()) || j2gort.IfaceEquals[pkg_rhino.JSType](n.GetJSType(), nil) {
		return
	}
	var classSymbol *SymbolTable_Symbol = this_.outer_.priv_SymbolTable_getSymbolForTypeHelper_5216(n.GetJSType(), false)
	_ = classSymbol
	if classSymbol != nil {
		classSymbol.DefineReferenceAt(n)
	}
}

// Java class:
type SymbolTable__Anon_92 struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_                   *SymbolTable
	cap_0_looksLikeGoogCall  j2gort.Ext_java_util_function_Predicate
	cap_1_declaredNamespaces j2gort.Ext_java_util_Map
	cap_2_moduleScopes       j2gort.Ext_java_util_Map
}

func (this_ *SymbolTable__Anon_92) BindSelf_(s any) {
	_ = s
}
func (this_ *SymbolTable__Anon_92) IsJavaNil_() bool { return this_ == nil }
func (*SymbolTable__Anon_92) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SymbolTable$3", ""
}
func (this_ *SymbolTable__Anon_92) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*SymbolTable__Anon_92) HasIdentityEquals_()                                        {}
func (this_ *SymbolTable__Anon_92) SetOuter_(o *SymbolTable)                             { this_.outer_ = o }
func (this_ *SymbolTable__Anon_92) GetOuter_() *SymbolTable                              { return this_.outer_ }
func (*SymbolTable__Anon_92) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*SymbolTable__Anon_92) Implements_jscomp_NodeTraversal_Callback()                  {}

func (this_ *SymbolTable__Anon_92) initFields_SymbolTable__Anon_92() {
	_ = this_
}

func NewSymbolTable__Anon_92(outer_ *SymbolTable, cap_p_looksLikeGoogCall j2gort.Ext_java_util_function_Predicate, cap_p_declaredNamespaces j2gort.Ext_java_util_Map, cap_p_moduleScopes j2gort.Ext_java_util_Map) *SymbolTable__Anon_92 {
	this_ := &SymbolTable__Anon_92{}
	this_.outer_ = outer_
	this_.cap_0_looksLikeGoogCall = cap_p_looksLikeGoogCall
	this_.cap_1_declaredNamespaces = cap_p_declaredNamespaces
	this_.cap_2_moduleScopes = cap_p_moduleScopes
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_SymbolTable__Anon_92()
	return this_
}

// Java method: .visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *SymbolTable__Anon_92) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(j2gort.PredTestNode(this_.cap_0_looksLikeGoogCall, uint32(n))) {
		return
	}
	var arg pkg_rhino.Node = n.GetSecondChild()
	_ = arg
	var namespaceName string = ("ns$" + j2gort.StrValString(arg.GetString()))
	_ = namespaceName
	switch j2gort.CheckNotNullString(n.GetFirstChild().GetString()) {
	case "module", "provide":
		var ns *SymbolTable_Symbol = this_.outer_.priv_SymbolTable_declareSymbol_5233(namespaceName, nil, false, this_.outer_.GetGlobalScope(), arg, nil)
		_ = ns
		j2gort.MapPutStr(this_.cap_1_declaredNamespaces, namespaceName, j2gort.PtrToAny(ns))
		if n.GetGrandparent().IsModuleBody() {
			j2gort.MapPutStr(this_.cap_2_moduleScopes, arg.GetString(), j2gort.PtrToAny(this_.outer_.scopes.Get(n.GetGrandparent())))
		}
	default:
	}
}

// Java class:
type SymbolTable__Anon_93 struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_                   *SymbolTable
	cap_0_looksLikeGoogCall  j2gort.Ext_java_util_function_Predicate
	cap_1_moduleScopes       j2gort.Ext_java_util_Map
	cap_2_declaredNamespaces j2gort.Ext_java_util_Map
}

func (this_ *SymbolTable__Anon_93) BindSelf_(s any) {
	_ = s
}
func (this_ *SymbolTable__Anon_93) IsJavaNil_() bool { return this_ == nil }
func (*SymbolTable__Anon_93) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SymbolTable$4", ""
}
func (this_ *SymbolTable__Anon_93) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*SymbolTable__Anon_93) HasIdentityEquals_()                                        {}
func (this_ *SymbolTable__Anon_93) SetOuter_(o *SymbolTable)                             { this_.outer_ = o }
func (this_ *SymbolTable__Anon_93) GetOuter_() *SymbolTable                              { return this_.outer_ }
func (*SymbolTable__Anon_93) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*SymbolTable__Anon_93) Implements_jscomp_NodeTraversal_Callback()                  {}

func (this_ *SymbolTable__Anon_93) initFields_SymbolTable__Anon_93() {
	_ = this_
}

func NewSymbolTable__Anon_93(outer_ *SymbolTable, cap_p_looksLikeGoogCall j2gort.Ext_java_util_function_Predicate, cap_p_moduleScopes j2gort.Ext_java_util_Map, cap_p_declaredNamespaces j2gort.Ext_java_util_Map) *SymbolTable__Anon_93 {
	this_ := &SymbolTable__Anon_93{}
	this_.outer_ = outer_
	this_.cap_0_looksLikeGoogCall = cap_p_looksLikeGoogCall
	this_.cap_1_moduleScopes = cap_p_moduleScopes
	this_.cap_2_declaredNamespaces = cap_p_declaredNamespaces
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_SymbolTable__Anon_93()
	return this_
}

// Java method: .visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *SymbolTable__Anon_93) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(j2gort.PredTestNode(this_.cap_0_looksLikeGoogCall, uint32(n))) {
		return
	}
	var arg pkg_rhino.Node = n.GetSecondChild()
	_ = arg
	var namespaceName string = ("ns$" + j2gort.StrValString(arg.GetString()))
	_ = namespaceName
	switch j2gort.CheckNotNullString(n.GetFirstChild().GetString()) {
	case "require", "requireType", "forwardDeclare":
		this_.outer_.priv_SymbolTable_addRefsInGoogRequireStatement_5249(n.GetParent(), this_.cap_1_moduleScopes)
		var symbol *SymbolTable_Symbol = j2gort.CastPtr[SymbolTable_Symbol](j2gort.MapGetStr(this_.cap_2_declaredNamespaces, namespaceName))
		_ = symbol
		if symbol != nil {
			symbol.DefineReferenceAt(arg)
		}
	default:
	}
}
func (*SymbolTable_PropertyRefCollector) AlwaysShouldTraverse_()      {}
func (*SymbolTable_PropertyRefCollector) CanReleaseSyntacticScopes_() {}
func (*SymbolTable_ThisRefCollector) AlwaysShouldTraverse_()          {}
func (*SymbolTable_ThisRefCollector) CanReleaseSyntacticScopes_()     {}
func (*SymbolTable_JSDocInfoCollector) AlwaysShouldTraverse_()        {}
func (*SymbolTable_JSDocInfoCollector) CanReleaseSyntacticScopes_()   {}
func (*SymbolTable_VisibilityCollector) AlwaysShouldTraverse_()       {}
func (*SymbolTable_VisibilityCollector) CanReleaseSyntacticScopes_()  {}
func (*SymbolTable__Anon_90) AlwaysShouldTraverse_()                  {}
func (*SymbolTable__Anon_90) NoOpVisit_()                             {}
func (*SymbolTable__Anon_91) AlwaysShouldTraverse_()                  {}
func (*SymbolTable__Anon_91) CanReleaseSyntacticScopes_()             {}
func (*SymbolTable__Anon_92) AlwaysShouldTraverse_()                  {}
func (*SymbolTable__Anon_92) CanReleaseSyntacticScopes_()             {}
func (*SymbolTable__Anon_93) AlwaysShouldTraverse_()                  {}
func (*SymbolTable__Anon_93) CanReleaseSyntacticScopes_()             {}
