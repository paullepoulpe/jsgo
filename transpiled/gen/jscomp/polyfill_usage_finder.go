/*
 * Copyright 2020 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/PolyfillUsageFinder.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_polyfill_usage_finder_1 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_polyfill_usage_finder_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_4 any) any {
		arg := j2gort.CastString(lp_4)
		_ = arg
		return j2gort.BoxString(j2gort.JStrSubFrom(arg, (j2gort.JStrLastIndexOfCp(arg, int32(uint16(46))) + int32(1))))
	})
}

// Java class: com.google.javascript.jscomp.PolyfillUsageFinder
type PolyfillUsageFinder struct {
	compiler                   AbstractCompiler                           `gson:"compiler"`
	polyfills                  *PolyfillUsageFinder_Polyfills             `gson:"polyfills"`
	registry                   *pkg_colors.ColorRegistry                  `gson:"registry"`
	hasInSupertypesListSeenSet j2gort.CustomOrderedSet[*pkg_colors.Color] `gson:"hasInSupertypesListSeenSet"`
	idHash_                    atomic.Int32
	disambiguateProperties     bool `gson:"disambiguateProperties"`
}

func (this_ *PolyfillUsageFinder) BindSelf_(s any) {
	_ = s
}
func (this_ *PolyfillUsageFinder) IsJavaNil_() bool { return this_ == nil }
func (*PolyfillUsageFinder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PolyfillUsageFinder", "PolyfillUsageFinder"
}
func (this_ *PolyfillUsageFinder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PolyfillUsageFinder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PolyfillUsageFinder", "com.google.javascript.jscomp.PolyfillUsageFinder", func(v_ any) bool {
		return j2gort.InstanceOf[*PolyfillUsageFinder](v_)
	})
}

var static_PolyfillUsageFinder_GLOBAL_NAMES j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_PolyfillUsageFinder_GLOBAL_NAMES() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_PolyfillUsageFinder()
	return static_PolyfillUsageFinder_GLOBAL_NAMES
}
func StaticSet_PolyfillUsageFinder_GLOBAL_NAMES(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_PolyfillUsageFinder()
	static_PolyfillUsageFinder_GLOBAL_NAMES = v
}

var static_PolyfillUsageFinder_STANDARD_COLLECTION_IDS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_PolyfillUsageFinder_STANDARD_COLLECTION_IDS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_PolyfillUsageFinder()
	return static_PolyfillUsageFinder_STANDARD_COLLECTION_IDS
}
func StaticSet_PolyfillUsageFinder_STANDARD_COLLECTION_IDS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_PolyfillUsageFinder()
	static_PolyfillUsageFinder_STANDARD_COLLECTION_IDS = v
}

var clinitDone_PolyfillUsageFinder uint32
var clinitState_PolyfillUsageFinder j2gort.ClinitState

func clinit_PolyfillUsageFinder() {
	if atomic.LoadUint32(&clinitDone_PolyfillUsageFinder) != 0 {
		return
	}
	clinitSlow_PolyfillUsageFinder()
}

func clinitSlow_PolyfillUsageFinder() {
	if !j2gort.EnterClinit(&clinitState_PolyfillUsageFinder, &clinitDone_PolyfillUsageFinder) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_PolyfillUsageFinder, &clinitDone_PolyfillUsageFinder)
	static_PolyfillUsageFinder_GLOBAL_NAMES = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_4("goog.global.", "window.", "goog$global.", "globalThis.")
	static_PolyfillUsageFinder_STANDARD_COLLECTION_IDS = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_7(j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_ARRAY_ID()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_ARRAY_PROTOTYPE_ID()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_READONLY_ARRAY_ID()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_READONLY_ARRAY_PROTOTYPE_ID()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_I_TEMPLATE_ARRAY_ID()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_I_TEMPLATE_ARRAY_PROTOTYPE_ID()), []any{j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_MAP_ID()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_MAP_PROTOTYPE_ID()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_READONLY_MAP_ID()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_READONLY_MAP_PROTOTYPE_ID()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_SET_ID()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_SET_PROTOTYPE_ID()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_READONLY_SET_ID()), j2gort.BoxNullableVal(pkg_colors.StaticGet_StandardColors_READONLY_SET_PROTOTYPE_ID())})
}

func (this_ *PolyfillUsageFinder) initFields_PolyfillUsageFinder() {
	_ = this_
	this_.hasInSupertypesListSeenSet = j2gort.NewCustomOrderedSet[*pkg_colors.Color]()
}

// Java constructor: PolyfillUsageFinder(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.PolyfillUsageFinder.Polyfills polyfills, boolean disambiguateProperties)
func NewPolyfillUsageFinder(compiler AbstractCompiler, polyfills *PolyfillUsageFinder_Polyfills, disambiguateProperties bool) *PolyfillUsageFinder {
	clinit_PolyfillUsageFinder()
	this_ := &PolyfillUsageFinder{}
	InitPolyfillUsageFinder(this_, compiler, polyfills, disambiguateProperties)
	return this_
}

func InitPolyfillUsageFinder(this_ *PolyfillUsageFinder, compiler AbstractCompiler, polyfills *PolyfillUsageFinder_Polyfills, disambiguateProperties bool) {
	_ = compiler
	_ = polyfills
	_ = disambiguateProperties
	this_.initFields_PolyfillUsageFinder()
	this_.polyfills = polyfills
	this_.compiler = compiler
	this_.disambiguateProperties = disambiguateProperties
	if compiler.HasOptimizationColors() {
		this_.registry = compiler.GetColorRegistry()
	} else {
		this_.registry = nil
	}
}

// Java static method: PolyfillUsageFinder.getPolyfillSupportedFeatureSet(java.lang.String nativeVersionStr)
func PolyfillUsageFinder_GetPolyfillSupportedFeatureSet(nativeVersionStr string) *pkg_parser.FeatureSet {
	clinit_PolyfillUsageFinder()
	_ = nativeVersionStr
	return pkg_parser.FeatureSet_BrowserFeatureSetValueOf(nativeVersionStr).Without_2(pkg_parser.StaticGet_FeatureSet_Feature_MODULES(), []*pkg_parser.FeatureSet_Feature{})
}

// Java method: PolyfillUsageFinder.traverseExcludingGuarded(com.google.javascript.rhino.Node root, java.util.function.Consumer<com.google.javascript.jscomp.PolyfillUsageFinder.PolyfillUsage> polyfillConsumer)
func (this_ *PolyfillUsageFinder) TraverseExcludingGuarded(root pkg_rhino.Node, polyfillConsumer j2gort.Ext_java_util_function_Consumer) {
	_ = root
	_ = polyfillConsumer
	{
		tr_ := NewPolyfillUsageFinder_Traverser(this_, this_.compiler, polyfillConsumer, StaticGet_PolyfillUsageFinder_Guard_ONLY_UNGUARDED())
		defer tr_.ReleaseSlabs_()
		NodeTraversal_Traverse(this_.compiler, root, tr_)
	}

}

// Java method: PolyfillUsageFinder.traverseIncludingGuarded(com.google.javascript.rhino.Node root, java.util.function.Consumer<com.google.javascript.jscomp.PolyfillUsageFinder.PolyfillUsage> polyfillConsumer)
func (this_ *PolyfillUsageFinder) TraverseIncludingGuarded(root pkg_rhino.Node, polyfillConsumer j2gort.Ext_java_util_function_Consumer) {
	_ = root
	_ = polyfillConsumer
	{
		tr_ := NewPolyfillUsageFinder_Traverser(this_, this_.compiler, polyfillConsumer, StaticGet_PolyfillUsageFinder_Guard_ALL())
		defer tr_.ReleaseSlabs_()
		NodeTraversal_Traverse(this_.compiler, root, tr_)
	}

}

// Java method: PolyfillUsageFinder.traverseOnlyGuarded(com.google.javascript.rhino.Node root, java.util.function.Consumer<com.google.javascript.jscomp.PolyfillUsageFinder.PolyfillUsage> polyfillConsumer)
func (this_ *PolyfillUsageFinder) TraverseOnlyGuarded(root pkg_rhino.Node, polyfillConsumer j2gort.Ext_java_util_function_Consumer) {
	_ = root
	_ = polyfillConsumer
	{
		tr_ := NewPolyfillUsageFinder_Traverser(this_, this_.compiler, polyfillConsumer, StaticGet_PolyfillUsageFinder_Guard_ONLY_GUARDED())
		defer tr_.ReleaseSlabs_()
		NodeTraversal_Traverse(this_.compiler, root, tr_)
	}

}

// Java method: PolyfillUsageFinder.maybeCreateStaticPolyfillUsageForGetPropChain(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node getPropNode)
func (this_ *PolyfillUsageFinder) priv_PolyfillUsageFinder_maybeCreateStaticPolyfillUsageForGetPropChain_4389(traversal *NodeTraversal, getPropNode pkg_rhino.Node) *PolyfillUsageFinder_PolyfillUsage {
	_ = traversal
	_ = getPropNode
	if !(getPropNode.IsGetProp() || getPropNode.IsOptChainGetProp()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(getPropNode))
	}
	var lastComponent string = getPropNode.GetString()
	_ = lastComponent
	if !(j2gort.SetContainsStr(this_.polyfills.suffixes, lastComponent)) {
		return nil
	}
	var components j2gort.Deque[string]
	_ = components
	components.AddFirst(lastComponent)
	var ownerNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = ownerNode
	{
		ownerNode = getPropNode.GetFirstChild()
		for ; ownerNode.IsGetProp() || ownerNode.IsOptChainGetProp(); ownerNode = ownerNode.GetFirstChild() {
			components.AddFirst(ownerNode.GetString())
		}
	}
	if !(ownerNode.IsName()) {
		return nil
	}
	var rootName string = ownerNode.GetString()
	_ = rootName
	components.AddFirst(rootName)
	var fullName string = j2gort.DequeJoinStrings(&components, ".")
	_ = fullName
	var globalPrefix string = PolyfillUsageFinder_FindGlobalPrefix(fullName)
	_ = globalPrefix
	if globalPrefix != j2gort.NullString {
		var polyfill *PolyfillUsageFinder_Polyfill = j2gort.CastPtr[PolyfillUsageFinder_Polyfill](j2gort.MapGetStr(this_.polyfills.statics, j2gort.JStrSubFrom(fullName, j2gort.JStrLen(globalPrefix))))
		_ = polyfill
		if polyfill != nil {
			return PolyfillUsageFinder_PolyfillUsage_CreateExplicit(polyfill, getPropNode, polyfill.nativeSymbol)
		}
	} else {
		var polyfill_2 *PolyfillUsageFinder_Polyfill = j2gort.CastPtr[PolyfillUsageFinder_Polyfill](j2gort.MapGetStr(this_.polyfills.statics, fullName))
		_ = polyfill_2
		if polyfill_2 != nil {
			if this_.priv_PolyfillUsageFinder_isPolyfillVisibleInScope_4390(traversal, rootName) {
				return PolyfillUsageFinder_PolyfillUsage_CreateNonExplicit(polyfill_2, getPropNode, polyfill_2.nativeSymbol)
			}
		}
	}
	return nil
}

// Java method: PolyfillUsageFinder.isPolyfillVisibleInScope(com.google.javascript.jscomp.NodeTraversal t, java.lang.String globalName)
func (this_ *PolyfillUsageFinder) priv_PolyfillUsageFinder_isPolyfillVisibleInScope_4390(t *NodeTraversal, globalName string) bool {
	_ = t
	_ = globalName
	return ((this_.compiler.GetLifeCycleStage() == StaticGet_AbstractCompiler_LifeCycleStage_NORMALIZED()) || (j2gort.CastPtr[Var](t.GetScope().GetVar(globalName)) == nil))
}

// Java static method: PolyfillUsageFinder.findGlobalPrefix(java.lang.String qualifiedName)
func PolyfillUsageFinder_FindGlobalPrefix(qualifiedName string) string {
	clinit_PolyfillUsageFinder()
	_ = qualifiedName
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(StaticGet_PolyfillUsageFinder_GLOBAL_NAMES()))
		for cur_1.HasNext() {
			global := cur_1.NextStr()
			_ = global
			if javalang.StringStartsWith(j2gort.CheckNotNullString(qualifiedName), j2gort.CheckNotNullString(global)) {
				return global
			}
		}
	}
	return j2gort.NullString
}

// Java static method: PolyfillUsageFinder.getMethodOwner(java.lang.String nativeSymbol)
func PolyfillUsageFinder_GetMethodOwner(nativeSymbol string) string {
	clinit_PolyfillUsageFinder()
	_ = nativeSymbol
	var protoIdx int32 = j2gort.JStrIndexOf(nativeSymbol, ".prototype.")
	_ = protoIdx
	if protoIdx != int32(-1) {
		return j2gort.JStrSub(nativeSymbol, int32(0), protoIdx)
	}
	return j2gort.NullString
}

// Java method: PolyfillUsageFinder.isDisjointFromIterator(com.google.javascript.rhino.Node getPropNode, java.lang.@org.jspecify.annotations.Nullable String ownerName)
func (this_ *PolyfillUsageFinder) IsDisjointFromIterator_2(getPropNode pkg_rhino.Node, ownerName string) bool {
	_ = getPropNode
	_ = ownerName
	if !("Iterator" == ownerName) {
		return false
	}
	if this_.registry == nil {
		return false
	}
	var receiverColor *pkg_colors.Color = getPropNode.GetFirstChild().GetColor()
	_ = receiverColor
	if receiverColor == nil {
		return false
	}
	return this_.IsDisjointFromIterator_1(receiverColor)
}

// Java method: PolyfillUsageFinder.isDisjointFromIterator(com.google.javascript.jscomp.colors.Color receiver)
func (this_ *PolyfillUsageFinder) IsDisjointFromIterator_1(receiver *pkg_colors.Color) bool {
	_ = receiver
	var iteratorColor *pkg_colors.Color = this_.registry.Get(pkg_colors.StaticGet_StandardColors_ITERATOR_ID())
	_ = iteratorColor
	if receiver.Equals(j2gort.PtrToAny(iteratorColor)) {
		return false
	}
	{
		var iifeCond_1 bool
		{

			rv_2 := receiver.GetId()
			av_3 := pkg_colors.StaticGet_StandardColors_GENERATOR_ID()
			_ = j2gort.CheckNotNullVal(rv_2)
			iifeCond_1 = rv_2 == av_3
		}
		if iifeCond_1 {
			return false
		}
	}
	if receiver.IsUnion() {
		var alts []*pkg_colors.Color = receiver.GetUnionElementsArray()
		_ = alts
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = int32(len(alts))
			_ = n
			for ; i < n; i = i + 1 {
				if !(this_.IsDisjointFromIterator_1(alts[i])) {
					return false
				}
			}
		}
		return true
	}
	if receiver.IsPrimitive() {
		return true
	}
	if receiver.IsInvalidating() && !(j2gort.SetContainsEqVal(StaticGet_PolyfillUsageFinder_STANDARD_COLLECTION_IDS(), receiver.GetId())) {
		return false
	}
	return (!(this_.priv_PolyfillUsageFinder_hasInSupertypesList_4393(receiver, iteratorColor)) && !(this_.priv_PolyfillUsageFinder_hasInSupertypesList_4393(iteratorColor, receiver)))
}

// Java method: PolyfillUsageFinder.hasInSupertypesList(com.google.javascript.jscomp.colors.Color subCtor, com.google.javascript.jscomp.colors.Color superCtor)
func (this_ *PolyfillUsageFinder) priv_PolyfillUsageFinder_hasInSupertypesList_4393(subCtor *pkg_colors.Color, superCtor *pkg_colors.Color) bool {
	_ = subCtor
	_ = superCtor
	if !(this_.hasInSupertypesListSeenSet.Add(subCtor)) {
		return false
	}
	defer func() {
		this_.hasInSupertypesListSeenSet.Remove(subCtor)
	}()
	{
		if subCtor.Equals(j2gort.PtrToAny(superCtor)) {
			return true
		}
		var supertypes []*pkg_colors.Color = this_.registry.GetDisambiguationSupertypesArray(subCtor)
		_ = supertypes
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = int32(len(supertypes))
			_ = n
			for ; i < n; i = i + 1 {
				var immediateSupertype *pkg_colors.Color = supertypes[i]
				_ = immediateSupertype
				if immediateSupertype.IsUnion() {
					var alts []*pkg_colors.Color = immediateSupertype.GetUnionElementsArray()
					_ = alts
					{
						var j int32 = int32(0)
						_ = j
						var m int32 = int32(len(alts))
						_ = m
						for ; j < m; j = j + 1 {
							if this_.priv_PolyfillUsageFinder_hasInSupertypesList_4393(alts[j], superCtor) {
								return true
							}
						}
					}
				} else if this_.priv_PolyfillUsageFinder_hasInSupertypesList_4393(immediateSupertype, superCtor) {
					return true
				}
			}
		}
		return false
	}
	panic("unreachable")
}

// Java class: com.google.javascript.jscomp.PolyfillUsageFinder.Polyfill
type PolyfillUsageFinder_Polyfill struct {
	nativeSymbol    string `gson:"nativeSymbol"`
	nativeVersion   string `gson:"nativeVersion"`
	polyfillVersion string `gson:"polyfillVersion"`
	library         string `gson:"library"`
	idHash_         atomic.Int32
	kind            PolyfillUsageFinder_Polyfill_Kind `gson:"kind"`
}

func (this_ *PolyfillUsageFinder_Polyfill) BindSelf_(s any) {
	_ = s
}
func (this_ *PolyfillUsageFinder_Polyfill) IsJavaNil_() bool { return this_ == nil }
func (*PolyfillUsageFinder_Polyfill) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PolyfillUsageFinder$Polyfill", "Polyfill"
}
func (this_ *PolyfillUsageFinder_Polyfill) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PolyfillUsageFinder_Polyfill) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PolyfillUsageFinder.Polyfill", "com.google.javascript.jscomp.PolyfillUsageFinder$Polyfill", func(v_ any) bool {
		return j2gort.InstanceOf[*PolyfillUsageFinder_Polyfill](v_)
	})
}

func (this_ *PolyfillUsageFinder_Polyfill) initFields_PolyfillUsageFinder_Polyfill() {
	_ = this_
	this_.nativeSymbol = j2gort.NullString
	this_.nativeVersion = j2gort.NullString
	this_.polyfillVersion = j2gort.NullString
	this_.library = j2gort.NullString
}

// Java constructor: Polyfill(java.lang.String nativeSymbol, java.lang.String nativeVersion, java.lang.String polyfillVersion, java.lang.String library, com.google.javascript.jscomp.PolyfillUsageFinder.Polyfill.Kind kind)
func NewPolyfillUsageFinder_Polyfill(nativeSymbol string, nativeVersion string, polyfillVersion string, library string, kind PolyfillUsageFinder_Polyfill_Kind) *PolyfillUsageFinder_Polyfill {
	this_ := &PolyfillUsageFinder_Polyfill{}
	InitPolyfillUsageFinder_Polyfill(this_, nativeSymbol, nativeVersion, polyfillVersion, library, kind)
	return this_
}

func InitPolyfillUsageFinder_Polyfill(this_ *PolyfillUsageFinder_Polyfill, nativeSymbol string, nativeVersion string, polyfillVersion string, library string, kind PolyfillUsageFinder_Polyfill_Kind) {
	_ = nativeSymbol
	_ = nativeVersion
	_ = polyfillVersion
	_ = library
	_ = kind
	this_.initFields_PolyfillUsageFinder_Polyfill()
	this_.nativeSymbol = nativeSymbol
	this_.nativeVersion = nativeVersion
	this_.polyfillVersion = polyfillVersion
	this_.library = library
	this_.kind = kind
}

// Java enum: com.google.javascript.jscomp.PolyfillUsageFinder.Polyfill.Kind
type PolyfillUsageFinder_Polyfill_Kind uint8

const (
	enum_PolyfillUsageFinder_Polyfill_Kind_STATIC PolyfillUsageFinder_Polyfill_Kind = 1
	enum_PolyfillUsageFinder_Polyfill_Kind_METHOD PolyfillUsageFinder_Polyfill_Kind = 2
)

var enumBaseTable_PolyfillUsageFinder_Polyfill_Kind = [...]j2gort.EnumBase{
	{Name_: "STATIC", Ordinal_: 0},
	{Name_: "METHOD", Ordinal_: 1},
}
var enumValues_PolyfillUsageFinder_Polyfill_Kind = []PolyfillUsageFinder_Polyfill_Kind{1, 2}

func StaticGet_PolyfillUsageFinder_Polyfill_Kind_STATIC() PolyfillUsageFinder_Polyfill_Kind {
	return enum_PolyfillUsageFinder_Polyfill_Kind_STATIC
}
func StaticGet_PolyfillUsageFinder_Polyfill_Kind_METHOD() PolyfillUsageFinder_Polyfill_Kind {
	return enum_PolyfillUsageFinder_Polyfill_Kind_METHOD
}

func PolyfillUsageFinder_Polyfill_Kind_Values() []PolyfillUsageFinder_Polyfill_Kind {
	return j2gort.CloneSlice(enumValues_PolyfillUsageFinder_Polyfill_Kind)
}

func PolyfillUsageFinder_Polyfill_Kind_ValueOf(name string) PolyfillUsageFinder_Polyfill_Kind {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_PolyfillUsageFinder_Polyfill_Kind {
		if enumBaseTable_PolyfillUsageFinder_Polyfill_Kind[i_].Name_ == name {
			return PolyfillUsageFinder_Polyfill_Kind(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.PolyfillUsageFinder.Polyfill.Kind." + name))
}

func (this_ PolyfillUsageFinder_Polyfill_Kind) IsJavaNil_() bool                              { return this_ == 0 }
func (PolyfillUsageFinder_Polyfill_Kind) HasIdentityEquals_()                                 {}
func (PolyfillUsageFinder_Polyfill_Kind) Implements_j2gort_Ext_java_lang_Enum()               {}
func (PolyfillUsageFinder_Polyfill_Kind) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ PolyfillUsageFinder_Polyfill_Kind) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (PolyfillUsageFinder_Polyfill_Kind) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PolyfillUsageFinder$Polyfill$Kind", "Kind"
}
func (this_ PolyfillUsageFinder_Polyfill_Kind) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_PolyfillUsageFinder_Polyfill_Kind[this_-1]
}
func (this_ PolyfillUsageFinder_Polyfill_Kind) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ PolyfillUsageFinder_Polyfill_Kind) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ PolyfillUsageFinder_Polyfill_Kind) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_PolyfillUsageFinder_Polyfill_Kind[this_-1].Name_
}
func (this_ PolyfillUsageFinder_Polyfill_Kind) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_PolyfillUsageFinder_Polyfill_Kind[this_-1].Name_
}
func (this_ PolyfillUsageFinder_Polyfill_Kind) Equals(other_ any) bool {
	o_, ok_ := other_.(PolyfillUsageFinder_Polyfill_Kind)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ PolyfillUsageFinder_Polyfill_Kind) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_PolyfillUsageFinder_Polyfill_Kind[this_-1].HashCode()
}
func (this_ PolyfillUsageFinder_Polyfill_Kind) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_PolyfillUsageFinder_Polyfill_Kind[this_-1].CompareTo(other_)
}
func (this_ PolyfillUsageFinder_Polyfill_Kind) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_PolyfillUsageFinder_Polyfill_Kind[this_-1].Name_)
}
func init() {
	vals_ := []any{PolyfillUsageFinder_Polyfill_Kind(1), PolyfillUsageFinder_Polyfill_Kind(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.PolyfillUsageFinder.Polyfill.Kind", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.PolyfillUsageFinder$Polyfill$Kind", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PolyfillUsageFinder.Polyfill.Kind", "com.google.javascript.jscomp.PolyfillUsageFinder$Polyfill$Kind", func(v_ any) bool {
		return j2gort.InstanceOf[PolyfillUsageFinder_Polyfill_Kind](v_)
	})
}

// Java class: com.google.javascript.jscomp.PolyfillUsageFinder.Polyfills
type PolyfillUsageFinder_Polyfills struct {
	methods         j2gort.Ext_com_google_common_collect_ImmutableMultimap `gson:"methods"`
	statics         j2gort.Ext_com_google_common_collect_ImmutableMap      `gson:"statics"`
	suffixes        j2gort.Ext_com_google_common_collect_ImmutableSet      `gson:"suffixes"`
	byNativeVersion j2gort.Ext_com_google_common_collect_ImmutableMultimap `gson:"byNativeVersion"`
	idHash_         atomic.Int32
}

func (this_ *PolyfillUsageFinder_Polyfills) BindSelf_(s any) {
	_ = s
}
func (this_ *PolyfillUsageFinder_Polyfills) IsJavaNil_() bool { return this_ == nil }
func (*PolyfillUsageFinder_Polyfills) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PolyfillUsageFinder$Polyfills", "Polyfills"
}
func (this_ *PolyfillUsageFinder_Polyfills) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PolyfillUsageFinder_Polyfills) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PolyfillUsageFinder.Polyfills", "com.google.javascript.jscomp.PolyfillUsageFinder$Polyfills", func(v_ any) bool {
		return j2gort.InstanceOf[*PolyfillUsageFinder_Polyfills](v_)
	})
}

func (this_ *PolyfillUsageFinder_Polyfills) initFields_PolyfillUsageFinder_Polyfills() {
	_ = this_
}

// Java constructor: Polyfills(com.google.common.collect.ImmutableMultimap<java.lang.String,com.google.javascript.jscomp.PolyfillUsageFinder.Polyfill> methods, com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.jscomp.PolyfillUsageFinder.Polyfill> statics, com.google.common.collect.ImmutableMultimap<java.lang.String,com.google.javascript.jscomp.PolyfillUsageFinder.Polyfill> byNativeVersion)
func NewPolyfillUsageFinder_Polyfills(methods j2gort.Ext_com_google_common_collect_ImmutableMultimap, statics j2gort.Ext_com_google_common_collect_ImmutableMap, byNativeVersion j2gort.Ext_com_google_common_collect_ImmutableMultimap) *PolyfillUsageFinder_Polyfills {
	this_ := &PolyfillUsageFinder_Polyfills{}
	InitPolyfillUsageFinder_Polyfills(this_, methods, statics, byNativeVersion)
	return this_
}

func InitPolyfillUsageFinder_Polyfills(this_ *PolyfillUsageFinder_Polyfills, methods j2gort.Ext_com_google_common_collect_ImmutableMultimap, statics j2gort.Ext_com_google_common_collect_ImmutableMap, byNativeVersion j2gort.Ext_com_google_common_collect_ImmutableMultimap) {
	_ = methods
	_ = statics
	_ = byNativeVersion
	this_.initFields_PolyfillUsageFinder_Polyfills()
	this_.methods = methods
	this_.statics = statics
	this_.byNativeVersion = byNativeVersion
	this_.suffixes = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](statics.KeySet_0()).Stream().Map(lambdaSingleton_polyfill_usage_finder_1).Collect_1(j2gort.Ext_java_util_stream_Collectors_ToList())))
}

// Java static method: Polyfills.fromTable(java.lang.String table)
func PolyfillUsageFinder_Polyfills_FromTable(table string) *PolyfillUsageFinder_Polyfills {
	_ = table
	var methods j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMultimap_S_Builder()
	_ = methods
	var statics j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
	_ = statics
	var byNativeVersion j2gort.Ext_com_google_common_collect_ImmutableMultimap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMultimap_S_Builder()
	_ = byNativeVersion
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.SplitterSplitStr(j2gort.Ext_com_google_common_base_Splitter_On_1_uint16(uint16(10)).OmitEmptyStrings(), table)))
		for cur_5.HasNext() {
			line := cur_5.NextStr()
			_ = line
			var tokens j2gort.Ext_java_util_List = j2gort.SplitterSplitToListStr(j2gort.Ext_com_google_common_base_Splitter_On_1_uint16(uint16(32)).OmitEmptyStrings(), javalang.StringTrim(j2gort.CheckNotNullString(line)))
			_ = tokens
			if (tokens.Size() == int32(1)) && (j2gort.CheckNotNullString(j2gort.ListGetStr(tokens, int32(0))) == "") {
				continue
			} else if tokens.Size() < int32(3) {
				panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Invalid table: too few tokens on line: "+j2gort.StrValString(line))))
			}
			var symbol string = j2gort.ListGetStr(tokens, int32(0))
			_ = symbol
			var isPrototypeMethod bool = (j2gort.JStrIndexOf(symbol, ".prototype.") >= 0)
			_ = isPrototypeMethod
			var nativeVersionStr string = j2gort.ListGetStr(tokens, int32(1))
			_ = nativeVersionStr
			var polyfillVersionStr string = j2gort.ListGetStr(tokens, int32(2))
			_ = polyfillVersionStr
			var iifeArg_1 string
		iifeSwitch_1:
			switch {
			default:

				if tokens.Size() > int32(3) {
					iifeArg_1 = j2gort.ListGetStr(tokens, int32(3))
					break iifeSwitch_1
				}
				iifeArg_1 = ""
				break iifeSwitch_1
			}
			var iifeArg_2 PolyfillUsageFinder_Polyfill_Kind
		iifeSwitch_2:
			switch {
			default:

				if isPrototypeMethod {
					iifeArg_2 = StaticGet_PolyfillUsageFinder_Polyfill_Kind_METHOD()
					break iifeSwitch_2
				}
				iifeArg_2 = StaticGet_PolyfillUsageFinder_Polyfill_Kind_STATIC()
				break iifeSwitch_2
			}
			var polyfill *PolyfillUsageFinder_Polyfill = NewPolyfillUsageFinder_Polyfill(symbol, nativeVersionStr, polyfillVersionStr, iifeArg_1,

				iifeArg_2,
			)
			_ = polyfill
			if isPrototypeMethod {
				j2gort.MultimapBuilderPutStr(methods, j2gort.JStrReplaceAll(symbol, ".*\\.prototype\\.", ""), polyfill)
			} else {
				statics.Put_2(j2gort.BoxString(symbol), polyfill)
			}
			j2gort.MultimapBuilderPutStr(byNativeVersion, nativeVersionStr, polyfill)
		}
	}
	return NewPolyfillUsageFinder_Polyfills(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap](methods.Build()), j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](statics.BuildOrThrow()), j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultimap](byNativeVersion.Build()))
}

// Java method: Polyfills.getPolyfillsNewerThan(com.google.javascript.jscomp.CompilerOptions.LanguageMode languageMode)
func (this_ *PolyfillUsageFinder_Polyfills) GetPolyfillsNewerThan(languageMode CompilerOptions_LanguageMode) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = languageMode
	var featureSet *pkg_parser.FeatureSet = languageMode.ToFeatureSet()
	_ = featureSet
	var result *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
	_ = result
	{
		cur_6 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](this_.byNativeVersion.KeySet())))
		for cur_6.HasNext() {
			nativeVersionStr := cur_6.NextStr()
			_ = nativeVersionStr
			var polyfillNativeFeatureSet *pkg_parser.FeatureSet = PolyfillUsageFinder_GetPolyfillSupportedFeatureSet(nativeVersionStr)
			_ = polyfillNativeFeatureSet
			if !(featureSet.Contains_1_FeatureSet(polyfillNativeFeatureSet)) {
				result.AddAllVoid_(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableCollection](j2gort.MultimapGetStr(this_.byNativeVersion, nativeVersionStr)))
			}
		}
	}
	return result.BuildImmList_()
}

// Java record: com.google.javascript.jscomp.PolyfillUsageFinder.PolyfillUsage
type PolyfillUsageFinder_PolyfillUsage struct {
	polyfill         *PolyfillUsageFinder_Polyfill `gson:"polyfill"`
	name             string                        `gson:"name"`
	node             pkg_rhino.Node                `gson:"node"`
	isExplicitGlobal bool                          `gson:"isExplicitGlobal"`
}

func (this_ *PolyfillUsageFinder_PolyfillUsage) BindSelf_(s any) {
	_ = s
}
func (this_ *PolyfillUsageFinder_PolyfillUsage) IsJavaNil_() bool { return this_ == nil }
func (*PolyfillUsageFinder_PolyfillUsage) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PolyfillUsageFinder$PolyfillUsage", "PolyfillUsage"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PolyfillUsageFinder.PolyfillUsage", "com.google.javascript.jscomp.PolyfillUsageFinder$PolyfillUsage", func(v_ any) bool {
		return j2gort.InstanceOf[*PolyfillUsageFinder_PolyfillUsage](v_)
	})
}
func (*PolyfillUsageFinder_PolyfillUsage) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *PolyfillUsageFinder_PolyfillUsage) initFields_PolyfillUsageFinder_PolyfillUsage() {
	_ = this_
}

// Java constructor: PolyfillUsage(com.google.javascript.jscomp.PolyfillUsageFinder.Polyfill polyfill, com.google.javascript.rhino.Node node, java.lang.String name, boolean isExplicitGlobal)
func NewPolyfillUsageFinder_PolyfillUsage(polyfill *PolyfillUsageFinder_Polyfill, node pkg_rhino.Node, name string, isExplicitGlobal bool) *PolyfillUsageFinder_PolyfillUsage {
	this_ := &PolyfillUsageFinder_PolyfillUsage{}
	InitPolyfillUsageFinder_PolyfillUsage(this_, polyfill, node, name, isExplicitGlobal)
	return this_
}

func InitPolyfillUsageFinder_PolyfillUsage(this_ *PolyfillUsageFinder_PolyfillUsage, polyfill *PolyfillUsageFinder_Polyfill, node pkg_rhino.Node, name string, isExplicitGlobal bool) {
	_ = polyfill
	_ = node
	_ = name
	_ = isExplicitGlobal
	this_.initFields_PolyfillUsageFinder_PolyfillUsage()
	j2gort.CheckNotNullPtr(polyfill)
	j2gort.CheckNotNullNode(node)
	j2gort.CheckNotNullString(name)
	this_.polyfill = polyfill
	this_.node = node
	this_.name = name
	this_.isExplicitGlobal = isExplicitGlobal
}

// Java static method: PolyfillUsage.createExplicit(com.google.javascript.jscomp.PolyfillUsageFinder.Polyfill polyfill, com.google.javascript.rhino.Node node, java.lang.String name)
func PolyfillUsageFinder_PolyfillUsage_CreateExplicit(polyfill *PolyfillUsageFinder_Polyfill, node pkg_rhino.Node, name string) *PolyfillUsageFinder_PolyfillUsage {
	_ = polyfill
	_ = node
	_ = name
	return NewPolyfillUsageFinder_PolyfillUsage(polyfill, node, name, true)
}

// Java static method: PolyfillUsage.createNonExplicit(com.google.javascript.jscomp.PolyfillUsageFinder.Polyfill polyfill, com.google.javascript.rhino.Node node, java.lang.String name)
func PolyfillUsageFinder_PolyfillUsage_CreateNonExplicit(polyfill *PolyfillUsageFinder_Polyfill, node pkg_rhino.Node, name string) *PolyfillUsageFinder_PolyfillUsage {
	_ = polyfill
	_ = node
	_ = name
	return NewPolyfillUsageFinder_PolyfillUsage(polyfill, node, name, false)
}

func (this_ *PolyfillUsageFinder_PolyfillUsage) Polyfill() *PolyfillUsageFinder_Polyfill {
	return this_.polyfill
}
func (this_ *PolyfillUsageFinder_PolyfillUsage) Node() pkg_rhino.Node { return this_.node }
func (this_ *PolyfillUsageFinder_PolyfillUsage) Name() string         { return this_.name }
func (this_ *PolyfillUsageFinder_PolyfillUsage) IsExplicitGlobal() bool {
	return this_.isExplicitGlobal
}
func (this_ *PolyfillUsageFinder_PolyfillUsage) Equals(other_ any) bool {
	o_, ok_ := other_.(*PolyfillUsageFinder_PolyfillUsage)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.polyfill != o_.polyfill {
		return false
	}
	if this_.node != o_.node {
		return false
	}
	if this_.name != o_.name {
		return false
	}
	if this_.isExplicitGlobal != o_.isExplicitGlobal {
		return false
	}
	return true
}
func (this_ *PolyfillUsageFinder_PolyfillUsage) HashCode() int32 {
	var h_ int32 = 0
	{
		hp_7 := this_.polyfill
		h_ = 31 * h_
		if hp_7 != nil {
			h_ += j2gort.IdentityHash(hp_7)
		}
	}
	h_ = 31*h_ + int32(this_.node)
	{
		hs_8 := this_.name
		h_ = 31 * h_
		if hs_8 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_8)
		}
	}
	if this_.isExplicitGlobal {
		h_ = 31*h_ + 1231
	} else {
		h_ = 31*h_ + 1237
	}
	return h_
}
func (this_ *PolyfillUsageFinder_PolyfillUsage) String() string {
	_ = fmt.Sprint
	return "PolyfillUsage[" + "polyfill=" + j2gort.StrVal(j2gort.PtrToAny(this_.polyfill)) + ", " + "node=" + j2gort.StrVal(j2gort.BoxNode(this_.node)) + ", " + "name=" + j2gort.StrVal(j2gort.BoxString(this_.name)) + ", " + "isExplicitGlobal=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.isExplicitGlobal)) + "]"
}

// Java enum: com.google.javascript.jscomp.PolyfillUsageFinder.Guard
type PolyfillUsageFinder_Guard uint8

const (
	enum_PolyfillUsageFinder_Guard_ONLY_GUARDED   PolyfillUsageFinder_Guard = 1
	enum_PolyfillUsageFinder_Guard_ONLY_UNGUARDED PolyfillUsageFinder_Guard = 2
	enum_PolyfillUsageFinder_Guard_ALL            PolyfillUsageFinder_Guard = 3
)

var enumBaseTable_PolyfillUsageFinder_Guard = [...]j2gort.EnumBase{
	{Name_: "ONLY_GUARDED", Ordinal_: 0},
	{Name_: "ONLY_UNGUARDED", Ordinal_: 1},
	{Name_: "ALL", Ordinal_: 2},
}
var enumValues_PolyfillUsageFinder_Guard = []PolyfillUsageFinder_Guard{1, 2, 3}

func StaticGet_PolyfillUsageFinder_Guard_ONLY_GUARDED() PolyfillUsageFinder_Guard {
	return enum_PolyfillUsageFinder_Guard_ONLY_GUARDED
}
func StaticGet_PolyfillUsageFinder_Guard_ONLY_UNGUARDED() PolyfillUsageFinder_Guard {
	return enum_PolyfillUsageFinder_Guard_ONLY_UNGUARDED
}
func StaticGet_PolyfillUsageFinder_Guard_ALL() PolyfillUsageFinder_Guard {
	return enum_PolyfillUsageFinder_Guard_ALL
}

func PolyfillUsageFinder_Guard_Values() []PolyfillUsageFinder_Guard {
	return j2gort.CloneSlice(enumValues_PolyfillUsageFinder_Guard)
}

func PolyfillUsageFinder_Guard_ValueOf(name string) PolyfillUsageFinder_Guard {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_PolyfillUsageFinder_Guard {
		if enumBaseTable_PolyfillUsageFinder_Guard[i_].Name_ == name {
			return PolyfillUsageFinder_Guard(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.PolyfillUsageFinder.Guard." + name))
}

// Java method: Guard.shouldInclude(boolean isGuarded)
func (this_ PolyfillUsageFinder_Guard) ShouldInclude(isGuarded bool) bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	_ = isGuarded
	switch (this_).Ordinal() {
	case int32(2) /* Guard.ALL */ :
		return true
	case int32(0) /* Guard.ONLY_GUARDED */ :
		return isGuarded
	case int32(1) /* Guard.ONLY_UNGUARDED */ :
		return !(isGuarded)
	}
	panic("unreachable switch expression")
}

func (this_ PolyfillUsageFinder_Guard) IsJavaNil_() bool                              { return this_ == 0 }
func (PolyfillUsageFinder_Guard) HasIdentityEquals_()                                 {}
func (PolyfillUsageFinder_Guard) Implements_j2gort_Ext_java_lang_Enum()               {}
func (PolyfillUsageFinder_Guard) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ PolyfillUsageFinder_Guard) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (PolyfillUsageFinder_Guard) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PolyfillUsageFinder$Guard", "Guard"
}
func (this_ PolyfillUsageFinder_Guard) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_PolyfillUsageFinder_Guard[this_-1]
}
func (this_ PolyfillUsageFinder_Guard) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ PolyfillUsageFinder_Guard) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ PolyfillUsageFinder_Guard) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_PolyfillUsageFinder_Guard[this_-1].Name_
}
func (this_ PolyfillUsageFinder_Guard) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_PolyfillUsageFinder_Guard[this_-1].Name_
}
func (this_ PolyfillUsageFinder_Guard) Equals(other_ any) bool {
	o_, ok_ := other_.(PolyfillUsageFinder_Guard)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ PolyfillUsageFinder_Guard) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_PolyfillUsageFinder_Guard[this_-1].HashCode()
}
func (this_ PolyfillUsageFinder_Guard) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_PolyfillUsageFinder_Guard[this_-1].CompareTo(other_)
}
func (this_ PolyfillUsageFinder_Guard) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_PolyfillUsageFinder_Guard[this_-1].Name_)
}
func init() {
	vals_ := []any{PolyfillUsageFinder_Guard(1), PolyfillUsageFinder_Guard(2), PolyfillUsageFinder_Guard(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.PolyfillUsageFinder.Guard", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.PolyfillUsageFinder$Guard", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PolyfillUsageFinder.Guard", "com.google.javascript.jscomp.PolyfillUsageFinder$Guard", func(v_ any) bool {
		return j2gort.InstanceOf[PolyfillUsageFinder_Guard](v_)
	})
}

// Java class: com.google.javascript.jscomp.PolyfillUsageFinder.Traverser
type PolyfillUsageFinder_Traverser struct {
	GuardedCallbackBase
	outer_               *PolyfillUsageFinder
	polyfillConsumer     j2gort.Ext_java_util_function_Consumer `gson:"polyfillConsumer"`
	includeGuardedUsages PolyfillUsageFinder_Guard              `gson:"includeGuardedUsages"`
}

func (this_ *PolyfillUsageFinder_Traverser) BindSelf_(s any) {
	_ = s
	this_.GuardedCallbackBase.BindSelf_(s)
}
func (this_ *PolyfillUsageFinder_Traverser) IsJavaNil_() bool { return this_ == nil }
func (*PolyfillUsageFinder_Traverser) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PolyfillUsageFinder$Traverser", "Traverser"
}
func (this_ *PolyfillUsageFinder_Traverser) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.GuardedCallbackBase.IdentityHash_()
}
func (*PolyfillUsageFinder_Traverser) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PolyfillUsageFinder.Traverser", "com.google.javascript.jscomp.PolyfillUsageFinder$Traverser", func(v_ any) bool {
		return j2gort.InstanceOf[*PolyfillUsageFinder_Traverser](v_)
	})
}
func (this_ *PolyfillUsageFinder_Traverser) SetOuter_(o *PolyfillUsageFinder)    { this_.outer_ = o }
func (this_ *PolyfillUsageFinder_Traverser) GetOuter_() *PolyfillUsageFinder     { return this_.outer_ }
func (*PolyfillUsageFinder_Traverser) Implements_jscomp_GuardedCallback()        {}
func (*PolyfillUsageFinder_Traverser) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *PolyfillUsageFinder_Traverser) initFields_PolyfillUsageFinder_Traverser() {
	_ = this_
}

// Java constructor: Traverser(com.google.javascript.jscomp.AbstractCompiler compiler, java.util.function.Consumer<com.google.javascript.jscomp.PolyfillUsageFinder.PolyfillUsage> polyfillConsumer, com.google.javascript.jscomp.PolyfillUsageFinder.Guard includeGuardedUsages)
func NewPolyfillUsageFinder_Traverser(outer_ *PolyfillUsageFinder, compiler AbstractCompiler, polyfillConsumer j2gort.Ext_java_util_function_Consumer, includeGuardedUsages PolyfillUsageFinder_Guard) *PolyfillUsageFinder_Traverser {
	this_ := &PolyfillUsageFinder_Traverser{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitPolyfillUsageFinder_Traverser(this_, compiler, polyfillConsumer, includeGuardedUsages)
	return this_
}

func InitPolyfillUsageFinder_Traverser(this_ *PolyfillUsageFinder_Traverser, compiler AbstractCompiler, polyfillConsumer j2gort.Ext_java_util_function_Consumer, includeGuardedUsages PolyfillUsageFinder_Guard) {
	_ = compiler
	_ = polyfillConsumer
	_ = includeGuardedUsages
	InitGuardedCallback(&this_.GuardedCallbackBase, compiler)
	this_.initFields_PolyfillUsageFinder_Traverser()
	this_.polyfillConsumer = polyfillConsumer
	this_.includeGuardedUsages = includeGuardedUsages
}

// Java method: Traverser.visitGuarded(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *PolyfillUsageFinder_Traverser) VisitGuarded(traversal *NodeTraversal, node pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = traversal
	_ = node
	_ = parent
	switch node.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		this_.priv_Traverser_visitName_4400(traversal, node)
	case int32(26) /* Token.GETPROP */, int32(29) /* Token.OPTCHAIN_GETPROP */ :
		this_.priv_Traverser_visitGetPropChain_4401(traversal, node)
	default:
	}
}

// Java method: Traverser.visitName(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node nameNode)
func (this_ *PolyfillUsageFinder_Traverser) priv_Traverser_visitName_4400(traversal *NodeTraversal, nameNode pkg_rhino.Node) {
	_ = traversal
	_ = nameNode
	var name string = nameNode.GetString()
	_ = name
	var polyfill *PolyfillUsageFinder_Polyfill = j2gort.CastPtr[PolyfillUsageFinder_Polyfill](j2gort.MapGetStr(this_.outer_.polyfills.statics, name))
	_ = polyfill
	if polyfill == nil {
		return
	}
	if this_.outer_.priv_PolyfillUsageFinder_isPolyfillVisibleInScope_4390(traversal, name) && this_.includeGuardedUsages.ShouldInclude(this_.IsGuarded(j2gort.BoxString(name))) {
		this_.polyfillConsumer.Accept(j2gort.PtrToAny(PolyfillUsageFinder_PolyfillUsage_CreateNonExplicit(polyfill, nameNode, name)))
	}
}

// Java method: Traverser.visitGetPropChain(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node getPropNode)
func (this_ *PolyfillUsageFinder_Traverser) priv_Traverser_visitGetPropChain_4401(traversal *NodeTraversal, getPropNode pkg_rhino.Node) {
	_ = traversal
	_ = getPropNode
	var staticPolyfillUsage *PolyfillUsageFinder_PolyfillUsage = this_.outer_.priv_PolyfillUsageFinder_maybeCreateStaticPolyfillUsageForGetPropChain_4389(traversal, getPropNode)
	_ = staticPolyfillUsage
	if staticPolyfillUsage != nil {
		if this_.includeGuardedUsages.ShouldInclude(this_.IsGuarded(j2gort.BoxString(staticPolyfillUsage.Name()))) {
			this_.polyfillConsumer.Accept(j2gort.PtrToAny(staticPolyfillUsage))
		}
	} else {
		var propertyName string = getPropNode.GetString()
		_ = propertyName
		var methodPolyfills j2gort.Ext_com_google_common_collect_ImmutableCollection = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableCollection](j2gort.MultimapGetStr(this_.outer_.polyfills.methods, propertyName))
		_ = methodPolyfills
		if !(methodPolyfills.IsEmpty()) && this_.includeGuardedUsages.ShouldInclude(this_.IsGuarded("."+j2gort.StrValString(propertyName))) {
			{
				cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(methodPolyfills))
				for cur_7.HasNext() {
					rawElem_8 := cur_7.Next()
					polyfill := j2gort.CastPtr[PolyfillUsageFinder_Polyfill](rawElem_8)
					_ = polyfill
					var ownerName string = PolyfillUsageFinder_GetMethodOwner(polyfill.nativeSymbol)
					_ = ownerName
					if this_.outer_.disambiguateProperties && this_.outer_.IsDisjointFromIterator_2(getPropNode, ownerName) {
						continue
					}
					this_.polyfillConsumer.Accept(j2gort.PtrToAny(PolyfillUsageFinder_PolyfillUsage_CreateNonExplicit(polyfill, getPropNode, propertyName)))
				}
			}
		}
	}
}
func (*PolyfillUsageFinder_Traverser) CanReleaseSyntacticScopes_() {}
