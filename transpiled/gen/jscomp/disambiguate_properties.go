/*
 * Copyright 2019 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/disambiguate/DisambiguateProperties.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	pkg_diagnostic "github.com/paullepoulpe/jsgo/transpiled/gen/diagnostic"
	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

var lambdaSingleton_disambiguate_properties_1 pkg_graph.LowestCommonAncestorFinder_Factory[*ColorGraphNode, any]
var lambdaSingleton_disambiguate_properties_2 j2gort.Ext_java_util_function_ToIntFunction
var lambdaSingleton_disambiguate_properties_3 j2gort.Ext_java_util_function_Function
var lambdaSingleton_disambiguate_properties_4 j2gort.Ext_java_util_function_Function
var lambdaSingleton_disambiguate_properties_5 j2gort.Ext_java_util_function_Function
var lambdaSingleton_disambiguate_properties_6 j2gort.Ext_java_util_function_Function
var lambdaSingleton_disambiguate_properties_7 j2gort.Ext_java_util_function_Function
var lambdaSingleton_disambiguate_properties_8 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_disambiguate_properties_1 = pkg_graph.NewLambda_LowestCommonAncestorFinder_Factory[*ColorGraphNode, any](func(mrp0_ pkg_graph.DiGraph[*ColorGraphNode, any]) *pkg_graph.LowestCommonAncestorFinder[*ColorGraphNode, any] {
		return pkg_graph.NewLowestCommonAncestorFinder[*ColorGraphNode, any](mrp0_)
	})
	lambdaSingleton_disambiguate_properties_2 = j2gort.NewLambda_Ext_java_util_function_ToIntFunction(func(lp_8 any) int32 {
		x := j2gort.CastPtr[DisambiguateProperties_TypeNodeJson](lp_8)
		_ = x
		return x.index
	})
	lambdaSingleton_disambiguate_properties_3 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewDisambiguateProperties_TypeNodeJson(j2gort.AsIface[pkg_graph.DiGraph_DiGraphNode[*ColorGraphNode, any]](mrp0_))
	})
	lambdaSingleton_disambiguate_properties_4 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return j2gort.BoxString(j2gort.CastPtr[PropertyClustering](mrp0_).GetName()) })
	lambdaSingleton_disambiguate_properties_5 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewDisambiguateProperties_TypeEdgeJson(j2gort.AsIface[pkg_graph.DiGraph_DiGraphEdge[*ColorGraphNode, any]](mrp0_))
	})
	lambdaSingleton_disambiguate_properties_6 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_25 any) any {
		e := j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](lp_25)
		_ = e
		return j2gort.BoxString(j2gort.CastPtr[PropertyClustering](e.GetKey()).GetName())
	})
	lambdaSingleton_disambiguate_properties_7 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](mrp0_).GetValue() })
	lambdaSingleton_disambiguate_properties_8 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return NewDisambiguateProperties_PropertyReferenceIndexJson(j2gort.CastPtr[PropertyClustering](mrp0_))
	})
}

// Java class: com.google.javascript.jscomp.disambiguate.DisambiguateProperties
type DisambiguateProperties struct {
	compiler                       AbstractCompiler                                                                        `gson:"compiler"`
	propertiesThatMustDisambiguate j2gort.Ext_com_google_common_collect_ImmutableSet                                       `gson:"propertiesThatMustDisambiguate"`
	registry                       *pkg_colors.ColorRegistry                                                               `gson:"registry"`
	mayHavePropertySeenSet         j2gort.CustomOrderedSet[*pkg_colors.Color]                                              `gson:"mayHavePropertySeenSet"`
	largeOwnPropertiesCache        *pkg_base.LinkedIdentityHashMap[*pkg_colors.Color, *j2gort.Ext_java_util_LinkedHashSet] `gson:"largeOwnPropertiesCache"`
	idHash_                        atomic.Int32
}

func (this_ *DisambiguateProperties) BindSelf_(s any) {
	_ = s
}
func (this_ *DisambiguateProperties) IsJavaNil_() bool { return this_ == nil }
func (*DisambiguateProperties) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.disambiguate.DisambiguateProperties", "DisambiguateProperties"
}
func (this_ *DisambiguateProperties) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DisambiguateProperties) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.DisambiguateProperties", "com.google.javascript.jscomp.disambiguate.DisambiguateProperties", func(v_ any) bool {
		return j2gort.InstanceOf[*DisambiguateProperties](v_)
	})
}
func (*DisambiguateProperties) Implements_jscomp_CompilerPass() {}

var static_DisambiguateProperties_PROPERTY_INVALIDATION *DiagnosticType

func StaticGet_DisambiguateProperties_PROPERTY_INVALIDATION() *DiagnosticType {
	clinit_DisambiguateProperties()
	return static_DisambiguateProperties_PROPERTY_INVALIDATION
}
func StaticSet_DisambiguateProperties_PROPERTY_INVALIDATION(v *DiagnosticType) {
	clinit_DisambiguateProperties()
	static_DisambiguateProperties_PROPERTY_INVALIDATION = v
}

var static_DisambiguateProperties_GSON *j2gort.Ext_com_google_gson_Gson

func StaticGet_DisambiguateProperties_GSON() *j2gort.Ext_com_google_gson_Gson {
	clinit_DisambiguateProperties()
	return static_DisambiguateProperties_GSON
}
func StaticSet_DisambiguateProperties_GSON(v *j2gort.Ext_com_google_gson_Gson) {
	clinit_DisambiguateProperties()
	static_DisambiguateProperties_GSON = v
}

var clinitDone_DisambiguateProperties uint32
var clinitState_DisambiguateProperties j2gort.ClinitState

func clinit_DisambiguateProperties() {
	if atomic.LoadUint32(&clinitDone_DisambiguateProperties) != 0 {
		return
	}
	clinitSlow_DisambiguateProperties()
}

func clinitSlow_DisambiguateProperties() {
	if !j2gort.EnterClinit(&clinitState_DisambiguateProperties, &clinitDone_DisambiguateProperties) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_DisambiguateProperties, &clinitDone_DisambiguateProperties)
	static_DisambiguateProperties_PROPERTY_INVALIDATION = DiagnosticType_Error("JSC_DISAMBIGUATE2_PROPERTY_INVALIDATION", "Property ''{0}'' was required to be disambiguated but was invalidated.{1}")
	static_DisambiguateProperties_GSON = j2gort.NewExt_com_google_gson_Gson()
}

func (this_ *DisambiguateProperties) initFields_DisambiguateProperties() {
	_ = this_
	this_.mayHavePropertySeenSet = j2gort.NewCustomOrderedSet[*pkg_colors.Color]()
	this_.largeOwnPropertiesCache = pkg_base.NewLinkedIdentityHashMap[*pkg_colors.Color, *j2gort.Ext_java_util_LinkedHashSet]()
}

// Java constructor: DisambiguateProperties(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.common.collect.ImmutableSet<java.lang.String> propertiesThatMustDisambiguate)
func NewDisambiguateProperties(compiler AbstractCompiler, propertiesThatMustDisambiguate j2gort.Ext_com_google_common_collect_ImmutableSet) *DisambiguateProperties {
	clinit_DisambiguateProperties()
	this_ := &DisambiguateProperties{}
	InitDisambiguateProperties(this_, compiler, propertiesThatMustDisambiguate)
	return this_
}

func InitDisambiguateProperties(this_ *DisambiguateProperties, compiler AbstractCompiler, propertiesThatMustDisambiguate j2gort.Ext_com_google_common_collect_ImmutableSet) {
	_ = compiler
	_ = propertiesThatMustDisambiguate
	this_.initFields_DisambiguateProperties()
	this_.compiler = compiler
	this_.propertiesThatMustDisambiguate = propertiesThatMustDisambiguate
	this_.registry = this_.compiler.GetColorRegistry()
}

// Java method: DisambiguateProperties.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *DisambiguateProperties) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if !(externs.GetParent() == root.GetParent()) {
		j2gort.CheckArgument(false)
	}
	var flattener *ColorGraphNodeFactory = ColorGraphNodeFactory_CreateFactory(this_.registry)
	_ = flattener
	var findRefs *ColorFindPropertyReferences = NewColorFindPropertyReferences(flattener, func() ColorFindPropertyReferences_IsPropertyReflector {
		bound_1 := this_.compiler.GetCodingConvention()
		return NewLambda_ColorFindPropertyReferences_IsPropertyReflector(func(mrp0_ pkg_rhino.Node) bool { return bound_1.IsPropertyRenameFunction(mrp0_) })
	}())
	_ = findRefs
	var val_2 ColorGraphBuilder
	InitColorGraphBuilder(&val_2, flattener, lambdaSingleton_disambiguate_properties_1, this_.registry)
	var graphBuilder *ColorGraphBuilder = &val_2
	_ = graphBuilder
	var val_3 ClusterPropagator
	InitClusterPropagator(&val_3)
	var propagator *ClusterPropagator = &val_3
	_ = propagator
	clinit_UseSiteRenamer()
	var val_4 UseSiteRenamer
	InitUseSiteRenamer(&val_4, func() UseSiteRenamer_MutationCallback {
		bound_5 := this_.compiler
		return NewLambda_UseSiteRenamer_MutationCallback(func(mrp0_ pkg_rhino.Node) { bound_5.ReportChangeToEnclosingScope(mrp0_) })
	}())
	var renamer *UseSiteRenamer = &val_4
	_ = renamer
	NodeTraversal_Traverse(this_.compiler, externs.GetParent(), findRefs)
	DisambiguateProperties_InvalidateWellKnownProperties(findRefs)
	var allProps j2gort.Ext_java_util_ArrayList = findRefs.CollectProperties()
	_ = allProps
	this_.priv_DisambiguateProperties_logForDiagnostics_6237("prop_refs", NewDisambiguateProperties__Anon_113(this_, allProps))
	graphBuilder.AddAll(j2gort.AsIface[j2gort.Ext_java_util_Collection](flattener.GetAllKnownTypes()))
	var graph pkg_graph.DiGraph[*ColorGraphNode, any] = j2gort.AsIface[pkg_graph.DiGraph[*ColorGraphNode, any]](graphBuilder.Build())
	_ = graph
	{
		cur_6 := j2gort.NewIterCursor(j2gort.CheckNotNull(flattener.GetAllKnownTypes()))
		for cur_6.HasNext() {
			rawElem_7 := cur_6.Next()
			colorGraphNode := j2gort.CastPtr[ColorGraphNode](rawElem_7)
			_ = colorGraphNode
			if graph.GetOutEdges(colorGraphNode).IsEmpty() {
				continue
			}
			if colorGraphNode.GetColor().IsUnion() {
				continue
			}
			this_.priv_DisambiguateProperties_registerOwnDeclaredProperties_6235(colorGraphNode, findRefs)
		}
	}
	findRefs.ClearProperties()
	this_.priv_DisambiguateProperties_logForDiagnostics_6236("graph", j2gort.NewLambda_Ext_com_google_common_base_Supplier(func() any {
		return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableList](graph.GetNodes().Stream().Map(lambdaSingleton_disambiguate_properties_3).Sorted_1(j2gort.Ext_java_util_Comparator_ComparingInt(lambdaSingleton_disambiguate_properties_2)).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableList_ToImmutableList()))
	}))
	this_.priv_DisambiguateProperties_invalidateBasedOnType_6231(flattener)
	pkg_graph.FixedPointGraphTraversal_NewTraversal[*ColorGraphNode, any](pkg_graph.ToIface_FixedPointGraphTraversal_EdgeCallback[*ColorGraphNode, any](propagator)).ComputeFixedPoint_1(graph)
	var val_9 DisambiguateProperties_TrackerSummaryGenerator
	InitDisambiguateProperties_TrackerSummaryGenerator(&val_9)
	var trackerSummaryGenerator *DisambiguateProperties_TrackerSummaryGenerator = &val_9
	_ = trackerSummaryGenerator
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = allProps.Size()
		_ = n
		for ; i < n; i = i + 1 {
			var prop *PropertyClustering = j2gort.CastPtr[PropertyClustering](allProps.Get(i))
			_ = prop
			var renameUsesResult UseSiteRenamer_RenameUsesResult = renamer.RenameUses(prop)
			_ = renameUsesResult
			trackerSummaryGenerator.AddRenameUsesResult(renameUsesResult)
			if prop.IsInvalidated() && j2gort.SetContainsStr(this_.propertiesThatMustDisambiguate, prop.GetName()) {
				this_.compiler.Report_1(this_.priv_DisambiguateProperties_createInvalidationError_6230(prop))
			}
		}
	}
	this_.compiler.ReportDisambiguatePropertiesSummary(func() j2gort.Ext_java_util_function_Supplier {
		bound_10 := trackerSummaryGenerator
		return j2gort.NewLambda_Ext_java_util_function_Supplier(func() any { return j2gort.BoxString(bound_10.priv_TrackerSummaryGenerator_buildSummaryString_6239()) })
	}())
	this_.priv_DisambiguateProperties_logForDiagnostics_6236("renaming_index", j2gort.NewLambda_Ext_com_google_common_base_Supplier(func() any {
		return DisambiguateProperties_BuildRenamingIndex(j2gort.AsIface[j2gort.Ext_java_util_List](allProps), renamer)
	}))
	this_.priv_DisambiguateProperties_logForDiagnostics_6236("mismatches", func() j2gort.Ext_com_google_common_base_Supplier {
		bound_11 := this_.registry.GetMismatchLocationsForDebugging()
		return j2gort.NewLambda_Ext_com_google_common_base_Supplier(func() any { return bound_11.Inverse() })
	}())
	GatherGetterAndSetterProperties_Update(this_.compiler, externs, root)
}

// Java method: DisambiguateProperties.createInvalidationError(com.google.javascript.jscomp.disambiguate.PropertyClustering prop)
func (this_ *DisambiguateProperties) priv_DisambiguateProperties_createInvalidationError_6230(prop *PropertyClustering) *JSError {
	_ = prop
	var additionalContext string = ""
	_ = additionalContext
	return JSError_Make_5(j2gort.NullString, int32(-1), int32(-1), StaticGet_DisambiguateProperties_PROPERTY_INVALIDATION(), []string{prop.GetName(), additionalContext})
}

// Java static method: DisambiguateProperties.buildRenamingIndex(java.util.List<com.google.javascript.jscomp.disambiguate.PropertyClustering> props, com.google.javascript.jscomp.disambiguate.UseSiteRenamer renamer)
func DisambiguateProperties_BuildRenamingIndex(props j2gort.Ext_java_util_List, renamer *UseSiteRenamer) j2gort.Ext_com_google_common_collect_ImmutableMap {
	clinit_DisambiguateProperties()
	_ = props
	_ = renamer
	var newNames *j2gort.Ext_com_google_common_collect_ImmutableSetMultimap = renamer.GetRenamingIndex()
	_ = newNames
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](j2gort.CastPtr[j2gort.Ext_com_google_common_collect_ImmutableSortedMap](props.Stream().Collect_1(j2gort.Ext_com_google_common_collect_ImmutableSortedMap_ToImmutableSortedMap_3(j2gort.Ext_java_util_Comparator_NaturalOrder(), lambdaSingleton_disambiguate_properties_4, j2gort.NewLambda_Ext_java_util_function_Function(func(lp_12 any) any {
		prop := j2gort.CastPtr[PropertyClustering](lp_12)
		_ = prop
		{

			if prop.IsInvalidated() {
				return prop.GetLastInvalidation()
			}
			return j2gort.Ext_com_google_common_collect_ImmutableSortedSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](j2gort.MultimapGetStr(newNames, prop.GetName())))
		}

	})))))
}

// Java static method: DisambiguateProperties.invalidateWellKnownProperties(com.google.javascript.jscomp.disambiguate.ColorFindPropertyReferences findRefs)
func DisambiguateProperties_InvalidateWellKnownProperties(findRefs *ColorFindPropertyReferences) {
	clinit_DisambiguateProperties()
	_ = findRefs
	findRefs.GetOrCreateProperty("prototype").Invalidate(Invalidation_WellKnownProperty())
	findRefs.GetOrCreateProperty("constructor").Invalidate(Invalidation_WellKnownProperty())
	findRefs.GetOrCreateProperty("then").Invalidate(Invalidation_WellKnownProperty())
}

// Java method: DisambiguateProperties.invalidateBasedOnType(com.google.javascript.jscomp.disambiguate.ColorGraphNodeFactory flattener)
func (this_ *DisambiguateProperties) priv_DisambiguateProperties_invalidateBasedOnType_6231(flattener *ColorGraphNodeFactory) {
	_ = flattener
	{
		cur_13 := j2gort.NewIterCursor(j2gort.CheckNotNull(flattener.GetAllKnownTypes()))
		for cur_13.HasNext() {
			rawElem_14 := cur_13.Next()
			type_ := j2gort.CastPtr[ColorGraphNode](rawElem_14)
			_ = type_
			if type_.GetColor().IsInvalidating() {
				{
					cur_15 := type_.GetAssociatedProps().IterKeys()
					for cur_15.HasNext() {

						prop := cur_15.
							Next()
						_ = prop
						prop.Invalidate(Invalidation_InvalidatingType(type_.GetIndex()))
					}
				}
			} else {
				this_.priv_DisambiguateProperties_invalidateNonDeclaredPropertyAccesses_6232(type_)
			}
		}
	}
}

// Java method: DisambiguateProperties.invalidateNonDeclaredPropertyAccesses(com.google.javascript.jscomp.disambiguate.ColorGraphNode colorGraphNode)
func (this_ *DisambiguateProperties) priv_DisambiguateProperties_invalidateNonDeclaredPropertyAccesses_6232(colorGraphNode *ColorGraphNode) {
	_ = colorGraphNode
	var color *pkg_colors.Color = colorGraphNode.GetColor()
	_ = color
	if !(!(color.IsInvalidating())) {
		j2gort.CheckArgumentStr(false, "Not applicable to invalidating types. All their properties are invalidated")
	}
	{
		cur_17 := colorGraphNode.GetAssociatedProps().IterKeys()
		for cur_17.HasNext() {

			prop := cur_17.
				Next()
			_ = prop
			if prop.IsInvalidated() {
				continue
			}
			if !(this_.priv_DisambiguateProperties_mayHaveProperty_6234(color, prop.GetName())) {
				prop.Invalidate(Invalidation_UndeclaredAccess(colorGraphNode.GetIndex()))
			}
		}
	}
}

// Java method: DisambiguateProperties.hasOwnPropertyFast(com.google.javascript.jscomp.colors.Color color, java.lang.String propertyName)
func (this_ *DisambiguateProperties) priv_DisambiguateProperties_hasOwnPropertyFast_6233(color *pkg_colors.Color, propertyName string) bool {
	_ = color
	_ = propertyName
	var ownProps []string = color.GetOwnPropertiesArray()
	_ = ownProps
	var len_ int32 = int32(len(ownProps))
	_ = len_
	if len_ == int32(0) {
		return false
	}
	if len_ <= int32(8) {
		{
			var i int32 = int32(0)
			_ = i
			for ; i < len_; i = i + 1 {
				if j2gort.CheckNotNullString(ownProps[i]) == propertyName {
					return true
				}
			}
		}
		return false
	}
	var cached *j2gort.Ext_java_util_LinkedHashSet = this_.largeOwnPropertiesCache.Get(color)
	_ = cached
	if cached == nil {
		cached = j2gort.NewExt_java_util_LinkedHashSet_1_int32((len_ * int32(2)))
		{
			var i_2 int32 = int32(0)
			_ = i_2
			for ; i_2 < len_; i_2 = i_2 + 1 {
				j2gort.SetAddStr(cached, ownProps[i_2])
			}
		}
		this_.largeOwnPropertiesCache.Put(color, cached)
	}
	return j2gort.SetContainsStr(cached, propertyName)
}

// Java method: DisambiguateProperties.mayHaveProperty(com.google.javascript.jscomp.colors.Color color, java.lang.String propertyName)
func (this_ *DisambiguateProperties) priv_DisambiguateProperties_mayHaveProperty_6234(color *pkg_colors.Color, propertyName string) bool {
	_ = color
	_ = propertyName
	defer func() {
		this_.mayHavePropertySeenSet.Remove(color)
	}()
	{
		if !(this_.mayHavePropertySeenSet.Add(color)) {
			return false
		}
		if color.IsUnion() {
			var unionElems []*pkg_colors.Color = color.GetUnionElementsArray()
			_ = unionElems
			{
				var i int32 = int32(0)
				_ = i
				var n int32 = int32(len(unionElems))
				_ = n
				for ; i < n; i = i + 1 {
					if this_.priv_DisambiguateProperties_mayHaveProperty_6234(unionElems[i], propertyName) {
						return true
					}
				}
			}
			return false
		}
		if this_.priv_DisambiguateProperties_hasOwnPropertyFast_6233(color, propertyName) {
			return true
		}
		var supertypes []*pkg_colors.Color = this_.registry.GetDisambiguationSupertypesArray(color)
		_ = supertypes
		{
			var i_2 int32 = int32(0)
			_ = i_2
			var n_2 int32 = int32(len(supertypes))
			_ = n_2
			for ; i_2 < n_2; i_2 = i_2 + 1 {
				if this_.priv_DisambiguateProperties_mayHaveProperty_6234(supertypes[i_2], propertyName) {
					return true
				}
			}
		}
		return false
	}
	panic("unreachable")
}

// Java method: DisambiguateProperties.registerOwnDeclaredProperties(com.google.javascript.jscomp.disambiguate.ColorGraphNode colorGraphNode, com.google.javascript.jscomp.disambiguate.ColorFindPropertyReferences findRefs)
func (this_ *DisambiguateProperties) priv_DisambiguateProperties_registerOwnDeclaredProperties_6235(colorGraphNode *ColorGraphNode, findRefs *ColorFindPropertyReferences) {
	_ = colorGraphNode
	_ = findRefs
	if !(!(colorGraphNode.GetColor().IsUnion())) {
		j2gort.CheckArgument(false)
	}
	var ownProps []string = colorGraphNode.GetColor().GetOwnPropertiesArray()
	_ = ownProps
	{
		var i int32 = int32(0)
		_ = i
		var n int32 = int32(len(ownProps))
		_ = n
		for ; i < n; i = i + 1 {
			var propName string = ownProps[i]
			_ = propName
			var prop *PropertyClustering = findRefs.GetProperty(propName)
			_ = prop
			if prop == nil {
				continue
			} else if prop.IsInvalidated() {
				continue
			}
			prop.GetClusters().Add(colorGraphNode)
			colorGraphNode.GetAssociatedProps().Put(prop, StaticGet_ColorGraphNode_PropAssociation_TYPE_SYSTEM())
		}
	}
}

// Java method: DisambiguateProperties.logForDiagnostics(java.lang.String name, com.google.common.base.Supplier<java.lang.Object> data)
func (this_ *DisambiguateProperties) priv_DisambiguateProperties_logForDiagnostics_6236(name string, data j2gort.Ext_com_google_common_base_Supplier) {
	_ = name
	_ = data
	func() {
		var log pkg_diagnostic.LogFile = this_.compiler.CreateOrReopenLog(j2gort.GetObjectClass(this_), (j2gort.StrValString(name) + ".log"), []string{})
		_ = log
		defer j2gort.CloseResource(log)
		if log.IsLogging() {
			log.LogJson_1_Ext_java_util_function_Supplier(j2gort.AsIface[j2gort.Ext_java_util_function_Supplier](data))
		}
	}()
}

// Java method: DisambiguateProperties.logForDiagnostics(java.lang.String name, com.google.javascript.jscomp.diagnostic.LogFile.StreamedJsonProducer data)
func (this_ *DisambiguateProperties) priv_DisambiguateProperties_logForDiagnostics_6237(name string, data pkg_diagnostic.LogFile_StreamedJsonProducer) {
	_ = name
	_ = data
	func() {
		var log pkg_diagnostic.LogFile = this_.compiler.CreateOrReopenLog(j2gort.GetObjectClass(this_), (j2gort.StrValString(name) + ".log"), []string{})
		_ = log
		defer j2gort.CloseResource(log)
		log.LogJson_1_LogFile_StreamedJsonProducer(data)
	}()
}

// Java class: com.google.javascript.jscomp.disambiguate.DisambiguateProperties.TrackerSummaryGenerator
type DisambiguateProperties_TrackerSummaryGenerator struct {
	total             int32 `gson:"total"`
	numInvalidated    int32 `gson:"numInvalidated"`
	numDisambiguated  int32 `gson:"numDisambiguated"`
	numOnlyOneCluster int32 `gson:"numOnlyOneCluster"`
	idHash_           atomic.Int32
}

func (this_ *DisambiguateProperties_TrackerSummaryGenerator) BindSelf_(s any) {
	_ = s
}
func (this_ *DisambiguateProperties_TrackerSummaryGenerator) IsJavaNil_() bool { return this_ == nil }
func (*DisambiguateProperties_TrackerSummaryGenerator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.disambiguate.DisambiguateProperties$TrackerSummaryGenerator", "TrackerSummaryGenerator"
}
func (this_ *DisambiguateProperties_TrackerSummaryGenerator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DisambiguateProperties_TrackerSummaryGenerator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.DisambiguateProperties.TrackerSummaryGenerator", "com.google.javascript.jscomp.disambiguate.DisambiguateProperties$TrackerSummaryGenerator", func(v_ any) bool {
		return j2gort.InstanceOf[*DisambiguateProperties_TrackerSummaryGenerator](v_)
	})
}

func (this_ *DisambiguateProperties_TrackerSummaryGenerator) initFields_DisambiguateProperties_TrackerSummaryGenerator() {
	_ = this_

}

// Java constructor: TrackerSummaryGenerator()
func NewDisambiguateProperties_TrackerSummaryGenerator() *DisambiguateProperties_TrackerSummaryGenerator {
	this_ := &DisambiguateProperties_TrackerSummaryGenerator{}
	InitDisambiguateProperties_TrackerSummaryGenerator(this_)
	return this_
}

func InitDisambiguateProperties_TrackerSummaryGenerator(this_ *DisambiguateProperties_TrackerSummaryGenerator) {
	this_.initFields_DisambiguateProperties_TrackerSummaryGenerator()
}

// Java method: TrackerSummaryGenerator.addRenameUsesResult(com.google.javascript.jscomp.disambiguate.UseSiteRenamer.RenameUsesResult renameUsesResult)
func (this_ *DisambiguateProperties_TrackerSummaryGenerator) AddRenameUsesResult(renameUsesResult UseSiteRenamer_RenameUsesResult) {
	_ = renameUsesResult
	this_.total = this_.total + 1
	switch (renameUsesResult).Ordinal() {
	case int32(0) /* RenameUsesResult.INVALIDATED */ :
		this_.numInvalidated = this_.numInvalidated + 1
	case int32(1) /* RenameUsesResult.ONLY_ONE_CLUSTER */ :
		this_.numOnlyOneCluster = this_.numOnlyOneCluster + 1
	case int32(2) /* RenameUsesResult.DISAMBIGUATED */ :
		this_.numDisambiguated = this_.numDisambiguated + 1
	}
}

// Java method: TrackerSummaryGenerator.buildSummaryString()
func (this_ *DisambiguateProperties_TrackerSummaryGenerator) priv_TrackerSummaryGenerator_buildSummaryString_6239() string {
	return j2gort.Ext_java_lang_String_Format_2("%d property names, %d disambiguated, %d invalidated, %d had a single cluster", []any{j2gort.ValueOfInt(this_.total), j2gort.ValueOfInt(this_.numDisambiguated), j2gort.ValueOfInt(this_.numInvalidated), j2gort.ValueOfInt(this_.numOnlyOneCluster)})
}

// Java class: com.google.javascript.jscomp.disambiguate.DisambiguateProperties.PropertyReferenceIndexJson
type DisambiguateProperties_PropertyReferenceIndexJson struct {
	name    string                                                  `gson:"name"`
	refs    j2gort.Ext_com_google_common_collect_ImmutableSortedSet `gson:"refs"`
	idHash_ atomic.Int32
}

func (this_ *DisambiguateProperties_PropertyReferenceIndexJson) BindSelf_(s any) {
	_ = s
}
func (this_ *DisambiguateProperties_PropertyReferenceIndexJson) IsJavaNil_() bool {
	return this_ == nil
}
func (*DisambiguateProperties_PropertyReferenceIndexJson) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.disambiguate.DisambiguateProperties$PropertyReferenceIndexJson", "PropertyReferenceIndexJson"
}
func (this_ *DisambiguateProperties_PropertyReferenceIndexJson) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DisambiguateProperties_PropertyReferenceIndexJson) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.DisambiguateProperties.PropertyReferenceIndexJson", "com.google.javascript.jscomp.disambiguate.DisambiguateProperties$PropertyReferenceIndexJson", func(v_ any) bool {
		return j2gort.InstanceOf[*DisambiguateProperties_PropertyReferenceIndexJson](v_)
	})
}

func (this_ *DisambiguateProperties_PropertyReferenceIndexJson) initFields_DisambiguateProperties_PropertyReferenceIndexJson() {
	_ = this_
	this_.name = j2gort.NullString
}

// Java constructor: PropertyReferenceIndexJson(com.google.javascript.jscomp.disambiguate.PropertyClustering prop)
func NewDisambiguateProperties_PropertyReferenceIndexJson(prop *PropertyClustering) *DisambiguateProperties_PropertyReferenceIndexJson {
	this_ := &DisambiguateProperties_PropertyReferenceIndexJson{}
	InitDisambiguateProperties_PropertyReferenceIndexJson(this_, prop)
	return this_
}

func InitDisambiguateProperties_PropertyReferenceIndexJson(this_ *DisambiguateProperties_PropertyReferenceIndexJson, prop *PropertyClustering) {
	_ = prop
	this_.initFields_DisambiguateProperties_PropertyReferenceIndexJson()
	this_.name = prop.GetName()
	if prop.IsInvalidated() {
		this_.refs = j2gort.Ext_com_google_common_collect_ImmutableSortedSet_Of_0()
	} else {
		var builder *j2gort.Ext_com_google_common_collect_ImmutableSortedSet_Builder = j2gort.Ext_com_google_common_collect_ImmutableSortedSet_NaturalOrder()
		_ = builder
		prop.ForEachUseSite(j2gort.NewLambda_Ext_java_util_function_BiConsumer(func(lp_19 any, lp_20 any) {
			site := j2gort.CastNode[pkg_rhino.Node](lp_19)
			_ = site
			owner := j2gort.CastPtr[ColorGraphNode](lp_20)
			_ = owner
			builder.Add_1_any(NewDisambiguateProperties_PropertyReferenceJson(site, owner))
		}))
		this_.refs = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSortedSet](builder.Build())
	}
}

// Java method: PropertyReferenceIndexJson.compareTo(com.google.javascript.jscomp.disambiguate.DisambiguateProperties.PropertyReferenceIndexJson x)
func (this_ *DisambiguateProperties_PropertyReferenceIndexJson) CompareTo(arg_21 any) int32 {
	x := j2gort.CastPtr[DisambiguateProperties_PropertyReferenceIndexJson](arg_21)
	_ = x
	return j2gort.JStrCompareTo(this_.name, x.name)
}

// Java method: PropertyReferenceIndexJson.writeJson(com.google.gson.stream.JsonWriter writer)
func (this_ *DisambiguateProperties_PropertyReferenceIndexJson) priv_PropertyReferenceIndexJson_writeJson_6240(writer *j2gort.Ext_com_google_gson_stream_JsonWriter) {
	_ = writer
	writer.Name(this_.name)
	writer.BeginObject()
	writer.Name("name").Value_1_string(this_.name)
	writer.Name("refs").BeginArray()
	{
		cur_22 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.refs))
		for cur_22.HasNext() {
			rawElem_23 := cur_22.Next()
			ref := j2gort.CastPtr[DisambiguateProperties_PropertyReferenceJson](rawElem_23)
			_ = ref
			StaticGet_DisambiguateProperties_GSON().ToJson_3_any_any_Ext_com_google_gson_stream_JsonWriter(j2gort.PtrToAny(ref), j2gort.PtrToAny(j2gort.ClassLiteral("com.google.javascript.jscomp.disambiguate.DisambiguateProperties.PropertyReferenceJson", "PropertyReferenceJson")), writer)
		}
	}
	writer.EndArray()
	writer.EndObject()
}

// Java class: com.google.javascript.jscomp.disambiguate.DisambiguateProperties.PropertyReferenceJson
type DisambiguateProperties_PropertyReferenceJson struct {
	location      string `gson:"location"`
	receiverIndex int32  `gson:"receiverIndex"`
	idHash_       atomic.Int32
}

func (this_ *DisambiguateProperties_PropertyReferenceJson) BindSelf_(s any) {
	_ = s
}
func (this_ *DisambiguateProperties_PropertyReferenceJson) IsJavaNil_() bool { return this_ == nil }
func (*DisambiguateProperties_PropertyReferenceJson) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.disambiguate.DisambiguateProperties$PropertyReferenceJson", "PropertyReferenceJson"
}
func (this_ *DisambiguateProperties_PropertyReferenceJson) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DisambiguateProperties_PropertyReferenceJson) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.DisambiguateProperties.PropertyReferenceJson", "com.google.javascript.jscomp.disambiguate.DisambiguateProperties$PropertyReferenceJson", func(v_ any) bool {
		return j2gort.InstanceOf[*DisambiguateProperties_PropertyReferenceJson](v_)
	})
}

func (this_ *DisambiguateProperties_PropertyReferenceJson) initFields_DisambiguateProperties_PropertyReferenceJson() {
	_ = this_
	this_.location = j2gort.NullString
}

// Java constructor: PropertyReferenceJson(com.google.javascript.rhino.Node location, com.google.javascript.jscomp.disambiguate.ColorGraphNode receiver)
func NewDisambiguateProperties_PropertyReferenceJson(location pkg_rhino.Node, receiver *ColorGraphNode) *DisambiguateProperties_PropertyReferenceJson {
	this_ := &DisambiguateProperties_PropertyReferenceJson{}
	InitDisambiguateProperties_PropertyReferenceJson(this_, location, receiver)
	return this_
}

func InitDisambiguateProperties_PropertyReferenceJson(this_ *DisambiguateProperties_PropertyReferenceJson, location pkg_rhino.Node, receiver *ColorGraphNode) {
	_ = location
	_ = receiver
	this_.initFields_DisambiguateProperties_PropertyReferenceJson()
	this_.location = location.GetLocation()
	this_.receiverIndex = receiver.GetIndex()
}

// Java method: PropertyReferenceJson.compareTo(com.google.javascript.jscomp.disambiguate.DisambiguateProperties.PropertyReferenceJson x)
func (this_ *DisambiguateProperties_PropertyReferenceJson) CompareTo(arg_24 any) int32 {
	x := j2gort.CastPtr[DisambiguateProperties_PropertyReferenceJson](arg_24)
	_ = x
	return j2gort.Ext_com_google_common_collect_ComparisonChain_Start().Compare_2_int32_int32(this_.receiverIndex, x.receiverIndex).Compare_2_any_any(j2gort.BoxString(this_.location), j2gort.BoxString(x.location)).Result()
}

// Java class: com.google.javascript.jscomp.disambiguate.DisambiguateProperties.TypeNodeJson
type DisambiguateProperties_TypeNodeJson struct {
	index        int32                                                    `gson:"index"`
	invalidating bool                                                     `gson:"invalidating"`
	colorId      string                                                   `gson:"colorId"`
	edges        j2gort.Ext_com_google_common_collect_ImmutableSortedSet  `gson:"edges"`
	props        *j2gort.Ext_com_google_common_collect_ImmutableSortedMap `gson:"props"`
	idHash_      atomic.Int32
}

func (this_ *DisambiguateProperties_TypeNodeJson) BindSelf_(s any) {
	_ = s
}
func (this_ *DisambiguateProperties_TypeNodeJson) IsJavaNil_() bool { return this_ == nil }
func (*DisambiguateProperties_TypeNodeJson) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.disambiguate.DisambiguateProperties$TypeNodeJson", "TypeNodeJson"
}
func (this_ *DisambiguateProperties_TypeNodeJson) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DisambiguateProperties_TypeNodeJson) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.DisambiguateProperties.TypeNodeJson", "com.google.javascript.jscomp.disambiguate.DisambiguateProperties$TypeNodeJson", func(v_ any) bool {
		return j2gort.InstanceOf[*DisambiguateProperties_TypeNodeJson](v_)
	})
}

func (this_ *DisambiguateProperties_TypeNodeJson) initFields_DisambiguateProperties_TypeNodeJson() {
	_ = this_
	this_.colorId = j2gort.NullString
}

// Java constructor: TypeNodeJson(com.google.javascript.jscomp.graph.DiGraph.DiGraphNode<com.google.javascript.jscomp.disambiguate.ColorGraphNode,java.lang.Object> n)
func NewDisambiguateProperties_TypeNodeJson(n pkg_graph.DiGraph_DiGraphNode[*ColorGraphNode, any]) *DisambiguateProperties_TypeNodeJson {
	this_ := &DisambiguateProperties_TypeNodeJson{}
	InitDisambiguateProperties_TypeNodeJson(this_, n)
	return this_
}

func InitDisambiguateProperties_TypeNodeJson(this_ *DisambiguateProperties_TypeNodeJson, n pkg_graph.DiGraph_DiGraphNode[*ColorGraphNode, any]) {
	_ = n
	this_.initFields_DisambiguateProperties_TypeNodeJson()
	var t *ColorGraphNode = n.GetValue()
	_ = t
	this_.index = t.GetIndex()
	this_.colorId = t.GetColor().GetId().String()
	this_.invalidating = t.GetColor().IsInvalidating()
	this_.edges = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSortedSet](n.GetOutEdges().Stream().Map(lambdaSingleton_disambiguate_properties_5).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableSortedSet_ToImmutableSortedSet(j2gort.Ext_java_util_Comparator_NaturalOrder())))
	this_.props = j2gort.CastPtr[j2gort.Ext_com_google_common_collect_ImmutableSortedMap](j2gort.AsIface[j2gort.Ext_java_util_Set](t.GetAssociatedProps().ToJavaLinkedHashMap().EntrySet()).Stream().Collect_1(j2gort.Ext_com_google_common_collect_ImmutableSortedMap_ToImmutableSortedMap_3(j2gort.Ext_java_util_Comparator_NaturalOrder(), lambdaSingleton_disambiguate_properties_6, lambdaSingleton_disambiguate_properties_7)))
}

// Java class: com.google.javascript.jscomp.disambiguate.DisambiguateProperties.TypeEdgeJson
type DisambiguateProperties_TypeEdgeJson struct {
	dest    int32 `gson:"dest"`
	value   any   `gson:"value"`
	idHash_ atomic.Int32
}

func (this_ *DisambiguateProperties_TypeEdgeJson) BindSelf_(s any) {
	_ = s
}
func (this_ *DisambiguateProperties_TypeEdgeJson) IsJavaNil_() bool { return this_ == nil }
func (*DisambiguateProperties_TypeEdgeJson) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.disambiguate.DisambiguateProperties$TypeEdgeJson", "TypeEdgeJson"
}
func (this_ *DisambiguateProperties_TypeEdgeJson) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DisambiguateProperties_TypeEdgeJson) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.DisambiguateProperties.TypeEdgeJson", "com.google.javascript.jscomp.disambiguate.DisambiguateProperties$TypeEdgeJson", func(v_ any) bool {
		return j2gort.InstanceOf[*DisambiguateProperties_TypeEdgeJson](v_)
	})
}

func (this_ *DisambiguateProperties_TypeEdgeJson) initFields_DisambiguateProperties_TypeEdgeJson() {
	_ = this_
}

// Java constructor: TypeEdgeJson(com.google.javascript.jscomp.graph.DiGraph.DiGraphEdge<com.google.javascript.jscomp.disambiguate.ColorGraphNode,java.lang.Object> e)
func NewDisambiguateProperties_TypeEdgeJson(e pkg_graph.DiGraph_DiGraphEdge[*ColorGraphNode, any]) *DisambiguateProperties_TypeEdgeJson {
	this_ := &DisambiguateProperties_TypeEdgeJson{}
	InitDisambiguateProperties_TypeEdgeJson(this_, e)
	return this_
}

func InitDisambiguateProperties_TypeEdgeJson(this_ *DisambiguateProperties_TypeEdgeJson, e pkg_graph.DiGraph_DiGraphEdge[*ColorGraphNode, any]) {
	_ = e
	this_.initFields_DisambiguateProperties_TypeEdgeJson()
	this_.dest = e.GetDestination().GetValue().GetIndex()
	this_.value = e.GetValue()
}

// Java method: TypeEdgeJson.compareTo(com.google.javascript.jscomp.disambiguate.DisambiguateProperties.TypeEdgeJson x)
func (this_ *DisambiguateProperties_TypeEdgeJson) CompareTo(arg_26 any) int32 {
	x := j2gort.CastPtr[DisambiguateProperties_TypeEdgeJson](arg_26)
	_ = x
	if !(this_.dest != x.dest) {
		j2gort.CheckArgument(false)
	}
	return (this_.dest - x.dest)
}

// Java class:
type DisambiguateProperties__Anon_113 struct {
	outer_         *DisambiguateProperties
	cap_0_allProps j2gort.Ext_java_util_ArrayList
	idHash_        atomic.Int32
}

func (this_ *DisambiguateProperties__Anon_113) BindSelf_(s any) {
	_ = s
}
func (this_ *DisambiguateProperties__Anon_113) IsJavaNil_() bool { return this_ == nil }
func (*DisambiguateProperties__Anon_113) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.disambiguate.DisambiguateProperties$1", ""
}
func (this_ *DisambiguateProperties__Anon_113) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DisambiguateProperties__Anon_113) HasIdentityEquals_()                       {}
func (this_ *DisambiguateProperties__Anon_113) SetOuter_(o *DisambiguateProperties) { this_.outer_ = o }
func (this_ *DisambiguateProperties__Anon_113) GetOuter_() *DisambiguateProperties {
	return this_.outer_
}
func (*DisambiguateProperties__Anon_113) Implements_diagnostic_LogFile_StreamedJsonProducer() {}

func (this_ *DisambiguateProperties__Anon_113) initFields_DisambiguateProperties__Anon_113() {
	_ = this_
}

func NewDisambiguateProperties__Anon_113(outer_ *DisambiguateProperties, cap_p_allProps j2gort.Ext_java_util_ArrayList) *DisambiguateProperties__Anon_113 {
	this_ := &DisambiguateProperties__Anon_113{}
	this_.outer_ = outer_
	this_.cap_0_allProps = cap_p_allProps
	this_.initFields_DisambiguateProperties__Anon_113()
	return this_
}

// Java method: .writeJson(com.google.gson.stream.JsonWriter writer)
func (this_ *DisambiguateProperties__Anon_113) WriteJson(writer *j2gort.Ext_com_google_gson_stream_JsonWriter) {
	_ = writer
	var propRefsJson j2gort.Ext_com_google_common_collect_ImmutableSortedSet = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSortedSet](this_.cap_0_allProps.Stream().Map(lambdaSingleton_disambiguate_properties_8).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableSortedSet_ToImmutableSortedSet(j2gort.Ext_java_util_Comparator_NaturalOrder())))
	_ = propRefsJson
	writer.BeginObject()
	{
		cur_27 := j2gort.NewIterCursor(j2gort.CheckNotNull(propRefsJson))
		for cur_27.HasNext() {
			rawElem_28 := cur_27.Next()
			propRef := j2gort.CastPtr[DisambiguateProperties_PropertyReferenceIndexJson](rawElem_28)
			_ = propRef
			propRef.priv_PropertyReferenceIndexJson_writeJson_6240(writer)
		}
	}
	writer.EndObject()
}
