/*
 * Copyright 2008 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/disambiguate/AmbiguateProperties.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	pkg_diagnostic "github.com/paullepoulpe/jsgo/transpiled/gen/diagnostic"
	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_ambiguate_properties_1 j2gort.Ext_java_util_Comparator
var lambdaSingleton_ambiguate_properties_2 pkg_graph.LowestCommonAncestorFinder_Factory[*ColorGraphNode, any]
var lambdaSingleton_ambiguate_properties_3 pkg_graph.FixedPointGraphTraversal_EdgeCallback[*ColorGraphNode, any]

func init() {
	lambdaSingleton_ambiguate_properties_1 = j2gort.NewLambda_Ext_java_util_Comparator(func(lp_1 any, lp_2 any) int32 {
		p1 := j2gort.CastPtr[AmbiguateProperties_Property](lp_1)
		_ = p1
		p2 := j2gort.CastPtr[AmbiguateProperties_Property](lp_2)
		_ = p2
		if p1.numOccurrences != p2.numOccurrences {
			return (p2.numOccurrences - p1.numOccurrences)
		}
		return j2gort.JStrCompareTo(p1.oldName, p2.oldName)
	})
	lambdaSingleton_ambiguate_properties_2 = pkg_graph.NewLambda_LowestCommonAncestorFinder_Factory[*ColorGraphNode, any](func(mrp0_ pkg_graph.DiGraph[*ColorGraphNode, any]) *pkg_graph.LowestCommonAncestorFinder[*ColorGraphNode, any] {
		return pkg_graph.NewLowestCommonAncestorFinder[*ColorGraphNode, any](mrp0_)
	})
	lambdaSingleton_ambiguate_properties_3 = pkg_graph.NewLambda_FixedPointGraphTraversal_EdgeCallback[*ColorGraphNode, any](func(subtype *ColorGraphNode, e any, supertype *ColorGraphNode) bool {
		_ = subtype
		_ = e
		_ = supertype
		if subtype.GetSubtypeIndices().Size() > supertype.GetSubtypeIndices().Size() {
			supertype.GetSubtypeIndices().Or(subtype.GetSubtypeIndices())
			return true
		}
		var startSize int32 = supertype.GetSubtypeIndices().Cardinality()
		_ = startSize
		supertype.GetSubtypeIndices().Or(subtype.GetSubtypeIndices())
		return (supertype.GetSubtypeIndices().Cardinality() > startSize)
	})
}

// Java class: com.google.javascript.jscomp.disambiguate.AmbiguateProperties
type AmbiguateProperties struct {
	compiler                   AbstractCompiler                                         `gson:"compiler"`
	stringNodesToRename        []pkg_rhino.Node                                         `gson:"stringNodesToRename"`
	reservedFirstCharacters    j2gort.Ext_java_util_Set                                 `gson:"reservedFirstCharacters"`
	reservedNonFirstCharacters j2gort.Ext_java_util_Set                                 `gson:"reservedNonFirstCharacters"`
	propertyMap                j2gort.OrderedMap[string, *AmbiguateProperties_Property] `gson:"propertyMap"`
	externedNames              j2gort.Ext_com_google_common_collect_ImmutableSet        `gson:"externedNames"`
	quotedNames                j2gort.Ext_java_util_Set                                 `gson:"quotedNames"`
	colorRegistry              *pkg_colors.ColorRegistry                                `gson:"colorRegistry"`
	renamingMap                j2gort.Ext_java_util_Map                                 `gson:"renamingMap"`
	graphNodeFactory           *ColorGraphNodeFactory                                   `gson:"graphNodeFactory"`
	idHash_                    atomic.Int32
}

func (this_ *AmbiguateProperties) BindSelf_(s any) {
	_ = s
}
func (this_ *AmbiguateProperties) IsJavaNil_() bool { return this_ == nil }
func (*AmbiguateProperties) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.disambiguate.AmbiguateProperties", "AmbiguateProperties"
}
func (this_ *AmbiguateProperties) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AmbiguateProperties) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.AmbiguateProperties", "com.google.javascript.jscomp.disambiguate.AmbiguateProperties", func(v_ any) bool {
		return j2gort.InstanceOf[*AmbiguateProperties](v_)
	})
}
func (*AmbiguateProperties) Implements_jscomp_CompilerPass() {}

var static_AmbiguateProperties_logger *j2gort.Ext_java_util_logging_Logger

func StaticGet_AmbiguateProperties_Logger() *j2gort.Ext_java_util_logging_Logger {
	clinit_AmbiguateProperties()
	return static_AmbiguateProperties_logger
}
func StaticSet_AmbiguateProperties_Logger(v *j2gort.Ext_java_util_logging_Logger) {
	clinit_AmbiguateProperties()
	static_AmbiguateProperties_logger = v
}

var static_AmbiguateProperties_FREQUENCY_COMPARATOR j2gort.Ext_java_util_Comparator

func StaticGet_AmbiguateProperties_FREQUENCY_COMPARATOR() j2gort.Ext_java_util_Comparator {
	clinit_AmbiguateProperties()
	return static_AmbiguateProperties_FREQUENCY_COMPARATOR
}
func StaticSet_AmbiguateProperties_FREQUENCY_COMPARATOR(v j2gort.Ext_java_util_Comparator) {
	clinit_AmbiguateProperties()
	static_AmbiguateProperties_FREQUENCY_COMPARATOR = v
}

var clinitDone_AmbiguateProperties uint32
var clinitState_AmbiguateProperties j2gort.ClinitState

func clinit_AmbiguateProperties() {
	if atomic.LoadUint32(&clinitDone_AmbiguateProperties) != 0 {
		return
	}
	clinitSlow_AmbiguateProperties()
}

func clinitSlow_AmbiguateProperties() {
	if !j2gort.EnterClinit(&clinitState_AmbiguateProperties, &clinitDone_AmbiguateProperties) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_AmbiguateProperties, &clinitDone_AmbiguateProperties)
	static_AmbiguateProperties_logger = j2gort.Ext_java_util_logging_Logger_GetLogger_1(j2gort.CheckNotNullPtr(j2gort.ClassLiteral("com.google.javascript.jscomp.disambiguate.AmbiguateProperties", "AmbiguateProperties")).GetName())
	static_AmbiguateProperties_FREQUENCY_COMPARATOR = lambdaSingleton_ambiguate_properties_1
}

func (this_ *AmbiguateProperties) initFields_AmbiguateProperties() {
	_ = this_
	this_.stringNodesToRename = []pkg_rhino.Node{}
	this_.propertyMap = j2gort.NewOrderedMap[string, *AmbiguateProperties_Property]()
	this_.quotedNames = j2gort.NewExt_java_util_LinkedHashSet()

}

// Java constructor: AmbiguateProperties(com.google.javascript.jscomp.AbstractCompiler compiler, java.util.Set<java.lang.Character> reservedFirstCharacters, java.util.Set<java.lang.Character> reservedNonFirstCharacters, java.util.Set<java.lang.String> externProperties)
func NewAmbiguateProperties(compiler AbstractCompiler, reservedFirstCharacters j2gort.Ext_java_util_Set, reservedNonFirstCharacters j2gort.Ext_java_util_Set, externProperties j2gort.Ext_java_util_Set) *AmbiguateProperties {
	clinit_AmbiguateProperties()
	this_ := &AmbiguateProperties{}
	InitAmbiguateProperties(this_, compiler, reservedFirstCharacters, reservedNonFirstCharacters, externProperties)
	return this_
}

func InitAmbiguateProperties(this_ *AmbiguateProperties, compiler AbstractCompiler, reservedFirstCharacters j2gort.Ext_java_util_Set, reservedNonFirstCharacters j2gort.Ext_java_util_Set, externProperties j2gort.Ext_java_util_Set) {
	_ = compiler
	_ = reservedFirstCharacters
	_ = reservedNonFirstCharacters
	_ = externProperties
	this_.initFields_AmbiguateProperties()
	if !(compiler.GetLifeCycleStage().IsNormalized()) {
		j2gort.CheckState(false)
	}
	this_.compiler = compiler
	this_.reservedFirstCharacters = reservedFirstCharacters
	this_.reservedNonFirstCharacters = reservedNonFirstCharacters
	this_.externedNames = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.ImmSetBuilderAddStr(j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder(), "prototype").AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](externProperties))).Build())
	this_.colorRegistry = compiler.GetColorRegistry()
}

// Java static method: AmbiguateProperties.makePassForTesting(com.google.javascript.jscomp.AbstractCompiler compiler, java.util.Set<java.lang.Character> reservedFirstCharacters, java.util.Set<java.lang.Character> reservedNonFirstCharacters, java.util.Set<java.lang.String> externProperties)
func AmbiguateProperties_MakePassForTesting(compiler AbstractCompiler, reservedFirstCharacters j2gort.Ext_java_util_Set, reservedNonFirstCharacters j2gort.Ext_java_util_Set, externProperties j2gort.Ext_java_util_Set) *AmbiguateProperties {
	clinit_AmbiguateProperties()
	_ = compiler
	_ = reservedFirstCharacters
	_ = reservedNonFirstCharacters
	_ = externProperties
	clinit_AmbiguateProperties()
	var val_3 AmbiguateProperties
	InitAmbiguateProperties(&val_3, compiler, reservedFirstCharacters, reservedNonFirstCharacters, externProperties)
	var ap *AmbiguateProperties = &val_3
	_ = ap
	{
		r_4 := ap
		r_4.renamingMap = j2gort.NewExt_java_util_LinkedHashMap()
	}
	return ap
}

// Java method: AmbiguateProperties.getRenamingMap()
func (this_ *AmbiguateProperties) GetRenamingMap() j2gort.Ext_java_util_Map {
	j2gort.CheckNotNull(this_.renamingMap)
	return this_.renamingMap
}

// Java method: AmbiguateProperties.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *AmbiguateProperties) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	this_.graphNodeFactory = ColorGraphNodeFactory_CreateFactory(this_.colorRegistry)
	NodeTraversal_Traverse(this_.compiler, root, NewAmbiguateProperties_ProcessPropertiesAndConstructors(this_))
	var val_5 ColorGraphBuilder
	InitColorGraphBuilder(&val_5, this_.graphNodeFactory, lambdaSingleton_ambiguate_properties_2, this_.colorRegistry)
	var graphBuilder *ColorGraphBuilder = &val_5
	_ = graphBuilder
	graphBuilder.AddAll(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.graphNodeFactory.GetAllKnownTypes()))
	var colorGraph pkg_graph.LinkedDirectedGraph[*ColorGraphNode, any] = graphBuilder.Build()
	_ = colorGraph
	func() {
		var chunkGraphLog pkg_diagnostic.LogFile = this_.compiler.CreateOrReopenLog(j2gort.GetObjectClass(this_), "color_graph.dot", []string{})
		_ = chunkGraphLog
		defer j2gort.CloseResource(chunkGraphLog)
		if chunkGraphLog.IsLogging() {
			chunkGraphLog.Log_1_string(DotFormatter_ToDot_1_GraphvizGraph(pkg_graph.Upcast_GraphvizGraph(colorGraph)))
		}
	}()
	{
		cur_6 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.graphNodeFactory.GetAllKnownTypes()))
		for cur_6.HasNext() {
			rawElem_7 := cur_6.Next()
			node := j2gort.CastPtr[ColorGraphNode](rawElem_7)
			_ = node
			node.GetSubtypeIndices().Set_1(node.GetIndex())
		}
	}
	pkg_graph.FixedPointGraphTraversal_NewReverseTraversal[*ColorGraphNode, any](lambdaSingleton_ambiguate_properties_3).ComputeFixedPoint_1(j2gort.AsIface[pkg_graph.DiGraph[*ColorGraphNode, any]](colorGraph))
	{
		cur_8 := this_.propertyMap.IterValues()
		for cur_8.HasNext() {

			prop := cur_8.
				Next()
			_ = prop
			if prop.relatedColorsSeeds == nil {
				continue
			}
			{
				cur_10 := prop.relatedColorsSeeds.IterKeys()
				for cur_10.HasNext() {

					color := cur_10.
						Next()
					_ = color
					prop.relatedColors.Or(color.GetSubtypeIndices())
				}
			}
			{
				r_12 := prop
				r_12.relatedColorsSeeds = nil
			}
		}
	}
	var reservedNames j2gort.Ext_com_google_common_collect_ImmutableSet_Builder = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder().AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.externedNames))).AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.quotedNames)))
	_ = reservedNames
	var numRenamedPropertyNames int32 = int32(0)
	_ = numRenamedPropertyNames
	var numSkippedPropertyNames int32 = int32(0)
	_ = numSkippedPropertyNames
	var nodes j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList_1_int32(this_.propertyMap.Size())
	_ = nodes
	{
		cur_13 := this_.propertyMap.IterValues()
		for cur_13.HasNext() {

			prop_2 := cur_13.
				Next()
			_ = prop_2
			if prop_2.skipAmbiguating {
				numSkippedPropertyNames = numSkippedPropertyNames + 1
				j2gort.ImmSetBuilderAddStrVoid(reservedNames, prop_2.oldName)
			} else {
				numRenamedPropertyNames = numRenamedPropertyNames + 1
				nodes.Add_1(NewAmbiguateProperties_PropertyGraphNode(prop_2))
			}
		}
	}
	var finalNumRenamedPropertyNames int32 = numRenamedPropertyNames
	_ = finalNumRenamedPropertyNames
	var finalNumSkippedPropertyNames int32 = numSkippedPropertyNames
	_ = finalNumSkippedPropertyNames
	var val_15 AmbiguateProperties_PropertyGraph
	InitAmbiguateProperties_PropertyGraph(&val_15, nodes)
	var propertyGraph *AmbiguateProperties_PropertyGraph = &val_15
	_ = propertyGraph
	var coloring pkg_graph.GraphColoring[*AmbiguateProperties_Property, any] = pkg_graph.ToIface_GraphColoring[*AmbiguateProperties_Property, any](pkg_graph.NewGraphColoring_GreedyGraphColoring_2[*AmbiguateProperties_Property, any](pkg_graph.ToIface_AdjacencyGraph[*AmbiguateProperties_Property, any](propertyGraph), StaticGet_AmbiguateProperties_FREQUENCY_COMPARATOR()))
	_ = coloring
	var numNewPropertyNames int32 = coloring.Color()
	_ = numNewPropertyNames
	var nameGen NameGenerator = NewDefaultNameGenerator_4_Ext_java_util_Set_string_Ext_java_util_Set_Ext_java_util_Set(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](reservedNames.Build()), "", this_.reservedFirstCharacters, this_.reservedNonFirstCharacters)
	_ = nameGen
	var colorMap []string = j2gort.NewStringArray(numNewPropertyNames)
	_ = colorMap
	{
		var i int32 = int32(0)
		_ = i
		for ; i < numNewPropertyNames; i = i + 1 {
			{
				a_16 := colorMap
				i_17 := i
				a_16[i_17] = nameGen.GenerateNextName()
			}
		}
	}
	{
		cur_18 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_List](propertyGraph.GetNodes())))
		for cur_18.HasNext() {
			rawElem_19 := cur_18.Next()
			node_2 := j2gort.CastPtr[AmbiguateProperties_PropertyGraphNode](rawElem_19)
			_ = node_2
			{
				r_20 := node_2.GetValue()
				r_20.newName = colorMap[javalang.ObjectHashCode(j2gort.CheckNotNull(node_2.GetAnnotation()))]
			}
			if this_.renamingMap != nil {
				j2gort.MapPutStr(this_.renamingMap, node_2.GetValue().oldName, j2gort.BoxString(node_2.GetValue().newName))
			}
		}
	}
	for _, n := range this_.stringNodesToRename {
		_ = n
		var oldName string = n.GetString()
		_ = oldName
		var p *AmbiguateProperties_Property = this_.propertyMap.Get(oldName)
		_ = p
		if (p != nil) && (p.newName != j2gort.NullString) {
			{
				var iifeCond_1 bool
				{

					sr_22 := oldName
					se_23 := p.oldName
					iifeCond_1 = j2gort.CheckNotNullString(sr_22) == se_23
				}
				if !(iifeCond_1) {
					j2gort.CheckState(false)
				}
			}
			if !(j2gort.CheckNotNullString(p.newName) == oldName) {
				n.SetString(p.newName)
				this_.compiler.ReportChangeToEnclosingScope(n)
			}
		}
	}

	GatherGetterAndSetterProperties_Update(this_.compiler, externs, root)
	var summarySupplier j2gort.Ext_java_util_function_Supplier = j2gort.AsIface[j2gort.Ext_java_util_function_Supplier](j2gort.Ext_com_google_common_base_Suppliers_Memoize(j2gort.NewLambda_Ext_com_google_common_base_Supplier(func() any {
		return j2gort.BoxString(("Collapsed " + j2gort.StrValInt32(finalNumRenamedPropertyNames) + " properties into " + j2gort.StrValInt32(numNewPropertyNames) + " and skipped renaming " + j2gort.StrValInt32(finalNumSkippedPropertyNames) + " properties."))
	})))
	_ = summarySupplier
	if StaticGet_AmbiguateProperties_Logger().IsLoggable(j2gort.Ext_java_util_logging_Level_F_FINE) {
		StaticGet_AmbiguateProperties_Logger().Fine_1_string(j2gort.CastString(summarySupplier.Get()))
	}
	this_.compiler.ReportAmbiguatePropertiesSummary(summarySupplier)
}

// Java method: AmbiguateProperties.getProperty(java.lang.String name)
func (this_ *AmbiguateProperties) priv_AmbiguateProperties_getProperty_6203(name string) *AmbiguateProperties_Property {
	_ = name
	var prop *AmbiguateProperties_Property = this_.propertyMap.Get(name)
	_ = prop
	if prop == nil {
		prop = NewAmbiguateProperties_Property(this_, name)
		this_.propertyMap.Put(name, prop)
	}
	return prop
}

// Java method: AmbiguateProperties.getColor(com.google.javascript.rhino.Node n)
func (this_ *AmbiguateProperties) priv_AmbiguateProperties_getColor_6204(n pkg_rhino.Node) *pkg_colors.Color {
	_ = n
	var type_ *pkg_colors.Color = n.GetColor()
	_ = type_
	if type_ == nil {
		return pkg_colors.StaticGet_StandardColors_UNKNOWN()
	} else {
		return type_
	}
}

// Java class: com.google.javascript.jscomp.disambiguate.AmbiguateProperties.PropertyGraph
type AmbiguateProperties_PropertyGraph struct {
	nodes   j2gort.Ext_java_util_ArrayList `gson:"nodes"`
	idHash_ atomic.Int32
}

func (this_ *AmbiguateProperties_PropertyGraph) BindSelf_(s any) {
	_ = s
}
func (this_ *AmbiguateProperties_PropertyGraph) IsJavaNil_() bool { return this_ == nil }
func (*AmbiguateProperties_PropertyGraph) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.disambiguate.AmbiguateProperties$PropertyGraph", "PropertyGraph"
}
func (this_ *AmbiguateProperties_PropertyGraph) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AmbiguateProperties_PropertyGraph) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.AmbiguateProperties.PropertyGraph", "com.google.javascript.jscomp.disambiguate.AmbiguateProperties$PropertyGraph", func(v_ any) bool {
		return j2gort.InstanceOf[*AmbiguateProperties_PropertyGraph](v_)
	})
}
func (*AmbiguateProperties_PropertyGraph) Implements_graph_AdjacencyGraph() {}

func (this_ *AmbiguateProperties_PropertyGraph) initFields_AmbiguateProperties_PropertyGraph() {
	_ = this_
}

// Java constructor: PropertyGraph(java.util.ArrayList<com.google.javascript.jscomp.disambiguate.AmbiguateProperties.PropertyGraphNode> nodes)
func NewAmbiguateProperties_PropertyGraph(nodes j2gort.Ext_java_util_ArrayList) *AmbiguateProperties_PropertyGraph {
	this_ := &AmbiguateProperties_PropertyGraph{}
	InitAmbiguateProperties_PropertyGraph(this_, nodes)
	return this_
}

func InitAmbiguateProperties_PropertyGraph(this_ *AmbiguateProperties_PropertyGraph, nodes j2gort.Ext_java_util_ArrayList) {
	_ = nodes
	this_.initFields_AmbiguateProperties_PropertyGraph()
	this_.nodes = nodes
}

// Java method: PropertyGraph.getNodes()
func (this_ *AmbiguateProperties_PropertyGraph) GetNodes() j2gort.Ext_java_util_Collection {
	return j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.nodes)
}

// Java method: PropertyGraph.getNodeCount()
func (this_ *AmbiguateProperties_PropertyGraph) GetNodeCount() int32 {
	return this_.nodes.Size()
}

// Java method: PropertyGraph.getNode(com.google.javascript.jscomp.disambiguate.AmbiguateProperties.Property property)
func (this_ *AmbiguateProperties_PropertyGraph) GetNode(property *AmbiguateProperties_Property) any {
	_ = property
	panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", "PropertyGraph#getNode is never called.")))
}

// Java method: PropertyGraph.newSubGraph()
func (this_ *AmbiguateProperties_PropertyGraph) NewSubGraph() pkg_graph.SubGraph[*AmbiguateProperties_Property, any] {
	return pkg_graph.ToIface_SubGraph[*AmbiguateProperties_Property, any](NewAmbiguateProperties_PropertySubGraph())
}

// Java method: PropertyGraph.clearNodeAnnotations()
func (this_ *AmbiguateProperties_PropertyGraph) ClearNodeAnnotations() {
	{
		cur_24 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.nodes))
		for cur_24.HasNext() {
			rawElem_25 := cur_24.Next()
			node := j2gort.CastPtr[AmbiguateProperties_PropertyGraphNode](rawElem_25)
			_ = node
			node.SetAnnotation(nil)
		}
	}
}

// Java method: PropertyGraph.getWeight(com.google.javascript.jscomp.disambiguate.AmbiguateProperties.Property value)
func (this_ *AmbiguateProperties_PropertyGraph) GetWeight(value *AmbiguateProperties_Property) int32 {
	_ = value
	return value.numOccurrences
}

// Java class: com.google.javascript.jscomp.disambiguate.AmbiguateProperties.PropertySubGraph
type AmbiguateProperties_PropertySubGraph struct {
	relatedTypes *j2gort.Ext_java_util_BitSet `gson:"relatedTypes"`
	idHash_      atomic.Int32
}

func (this_ *AmbiguateProperties_PropertySubGraph) BindSelf_(s any) {
	_ = s
}
func (this_ *AmbiguateProperties_PropertySubGraph) IsJavaNil_() bool { return this_ == nil }
func (*AmbiguateProperties_PropertySubGraph) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.disambiguate.AmbiguateProperties$PropertySubGraph", "PropertySubGraph"
}
func (this_ *AmbiguateProperties_PropertySubGraph) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AmbiguateProperties_PropertySubGraph) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.AmbiguateProperties.PropertySubGraph", "com.google.javascript.jscomp.disambiguate.AmbiguateProperties$PropertySubGraph", func(v_ any) bool {
		return j2gort.InstanceOf[*AmbiguateProperties_PropertySubGraph](v_)
	})
}
func (*AmbiguateProperties_PropertySubGraph) Implements_graph_SubGraph() {}

func (this_ *AmbiguateProperties_PropertySubGraph) initFields_AmbiguateProperties_PropertySubGraph() {
	_ = this_
	this_.relatedTypes = j2gort.NewExt_java_util_BitSet()
}

// Java constructor: PropertySubGraph()
func NewAmbiguateProperties_PropertySubGraph() *AmbiguateProperties_PropertySubGraph {
	this_ := &AmbiguateProperties_PropertySubGraph{}
	InitAmbiguateProperties_PropertySubGraph(this_)
	return this_
}

func InitAmbiguateProperties_PropertySubGraph(this_ *AmbiguateProperties_PropertySubGraph) {
	this_.initFields_AmbiguateProperties_PropertySubGraph()
}

// Java method: PropertySubGraph.isIndependentOf(com.google.javascript.jscomp.disambiguate.AmbiguateProperties.Property prop)
func (this_ *AmbiguateProperties_PropertySubGraph) IsIndependentOf(prop *AmbiguateProperties_Property) bool {
	_ = prop
	return !(this_.relatedTypes.Intersects(prop.relatedColors))
}

// Java method: PropertySubGraph.addNode(com.google.javascript.jscomp.disambiguate.AmbiguateProperties.Property prop)
func (this_ *AmbiguateProperties_PropertySubGraph) AddNode(prop *AmbiguateProperties_Property) {
	_ = prop
	this_.relatedTypes.Or(prop.relatedColors)
}

// Java class: com.google.javascript.jscomp.disambiguate.AmbiguateProperties.PropertyGraphNode
type AmbiguateProperties_PropertyGraphNode struct {
	property   *AmbiguateProperties_Property `gson:"property"`
	annotation pkg_graph.Annotation          `gson:"annotation"`
	idHash_    atomic.Int32
}

func (this_ *AmbiguateProperties_PropertyGraphNode) BindSelf_(s any) {
	_ = s
}
func (this_ *AmbiguateProperties_PropertyGraphNode) IsJavaNil_() bool { return this_ == nil }
func (*AmbiguateProperties_PropertyGraphNode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.disambiguate.AmbiguateProperties$PropertyGraphNode", "PropertyGraphNode"
}
func (this_ *AmbiguateProperties_PropertyGraphNode) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AmbiguateProperties_PropertyGraphNode) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.AmbiguateProperties.PropertyGraphNode", "com.google.javascript.jscomp.disambiguate.AmbiguateProperties$PropertyGraphNode", func(v_ any) bool {
		return j2gort.InstanceOf[*AmbiguateProperties_PropertyGraphNode](v_)
	})
}
func (*AmbiguateProperties_PropertyGraphNode) Implements_graph_GraphNode()   {}
func (*AmbiguateProperties_PropertyGraphNode) Implements_graph_Annotatable() {}

func (this_ *AmbiguateProperties_PropertyGraphNode) initFields_AmbiguateProperties_PropertyGraphNode() {
	_ = this_
}

// Java constructor: PropertyGraphNode(com.google.javascript.jscomp.disambiguate.AmbiguateProperties.Property property)
func NewAmbiguateProperties_PropertyGraphNode(property *AmbiguateProperties_Property) *AmbiguateProperties_PropertyGraphNode {
	this_ := &AmbiguateProperties_PropertyGraphNode{}
	InitAmbiguateProperties_PropertyGraphNode(this_, property)
	return this_
}

func InitAmbiguateProperties_PropertyGraphNode(this_ *AmbiguateProperties_PropertyGraphNode, property *AmbiguateProperties_Property) {
	_ = property
	this_.initFields_AmbiguateProperties_PropertyGraphNode()
	this_.property = property
}

// Java method: PropertyGraphNode.getValue()
func (this_ *AmbiguateProperties_PropertyGraphNode) GetValue() *AmbiguateProperties_Property {
	return this_.property
}

// Java method: PropertyGraphNode.getAnnotation()
func (this_ *AmbiguateProperties_PropertyGraphNode) GetAnnotation() pkg_graph.Annotation {
	return this_.annotation
}

// Java method: PropertyGraphNode.setAnnotation(com.google.javascript.jscomp.graph.@org.jspecify.annotations.Nullable Annotation data)
func (this_ *AmbiguateProperties_PropertyGraphNode) SetAnnotation(data pkg_graph.Annotation) {
	_ = data
	this_.annotation = data
}

// Java class: com.google.javascript.jscomp.disambiguate.AmbiguateProperties.ProcessPropertiesAndConstructors
type AmbiguateProperties_ProcessPropertiesAndConstructors struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *AmbiguateProperties
}

func (this_ *AmbiguateProperties_ProcessPropertiesAndConstructors) BindSelf_(s any) {
	_ = s
}
func (this_ *AmbiguateProperties_ProcessPropertiesAndConstructors) IsJavaNil_() bool {
	return this_ == nil
}
func (*AmbiguateProperties_ProcessPropertiesAndConstructors) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.disambiguate.AmbiguateProperties$ProcessPropertiesAndConstructors", "ProcessPropertiesAndConstructors"
}
func (this_ *AmbiguateProperties_ProcessPropertiesAndConstructors) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*AmbiguateProperties_ProcessPropertiesAndConstructors) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.AmbiguateProperties.ProcessPropertiesAndConstructors", "com.google.javascript.jscomp.disambiguate.AmbiguateProperties$ProcessPropertiesAndConstructors", func(v_ any) bool {
		return j2gort.InstanceOf[*AmbiguateProperties_ProcessPropertiesAndConstructors](v_)
	})
}
func (this_ *AmbiguateProperties_ProcessPropertiesAndConstructors) SetOuter_(o *AmbiguateProperties) {
	this_.outer_ = o
}
func (this_ *AmbiguateProperties_ProcessPropertiesAndConstructors) GetOuter_() *AmbiguateProperties {
	return this_.outer_
}
func (*AmbiguateProperties_ProcessPropertiesAndConstructors) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*AmbiguateProperties_ProcessPropertiesAndConstructors) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *AmbiguateProperties_ProcessPropertiesAndConstructors) initFields_AmbiguateProperties_ProcessPropertiesAndConstructors() {
	_ = this_
}

// Java constructor: ProcessPropertiesAndConstructors()
func NewAmbiguateProperties_ProcessPropertiesAndConstructors(outer_ *AmbiguateProperties) *AmbiguateProperties_ProcessPropertiesAndConstructors {
	this_ := &AmbiguateProperties_ProcessPropertiesAndConstructors{}
	this_.outer_ = outer_
	InitAmbiguateProperties_ProcessPropertiesAndConstructors(this_)
	return this_
}

func InitAmbiguateProperties_ProcessPropertiesAndConstructors(this_ *AmbiguateProperties_ProcessPropertiesAndConstructors) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_AmbiguateProperties_ProcessPropertiesAndConstructors()
}

// Java method: ProcessPropertiesAndConstructors.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *AmbiguateProperties_ProcessPropertiesAndConstructors) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(26) /* Token.GETPROP */, int32(29) /* Token.OPTCHAIN_GETPROP */ :
		this_.priv_ProcessPropertiesAndConstructors_processGetProp_6216(n)
		return
	case int32(28) /* Token.CALL */ :
		this_.priv_ProcessPropertiesAndConstructors_processCall_6217(n)
		return
	case int32(32) /* Token.NAME */ :
		if NodeUtil_IsNameDeclaration(parent) || parent.IsFunction() {
			this_.outer_.graphNodeFactory.CreateNode(this_.outer_.priv_AmbiguateProperties_getColor_6204(n))
		}
		return
	case int32(47) /* Token.OBJECTLIT */, int32(103) /* Token.OBJECT_PATTERN */ :
		this_.priv_ProcessPropertiesAndConstructors_processObjectLitOrPattern_6219(n)
		return
	case int32(27) /* Token.GETELEM */ :
		this_.priv_ProcessPropertiesAndConstructors_processGetElem_6220(n)
		return
	case int32(105) /* Token.CLASS */ :
		this_.priv_ProcessPropertiesAndConstructors_processClass_6221(n)
		return
	default:
	}
}

// Java method: ProcessPropertiesAndConstructors.processGetProp(com.google.javascript.rhino.Node getProp)
func (this_ *AmbiguateProperties_ProcessPropertiesAndConstructors) priv_ProcessPropertiesAndConstructors_processGetProp_6216(getProp pkg_rhino.Node) {
	_ = getProp
	var type_ *pkg_colors.Color = this_.outer_.priv_AmbiguateProperties_getColor_6204(getProp.GetFirstChild())
	_ = type_
	this_.priv_ProcessPropertiesAndConstructors_maybeMarkCandidate_6222(getProp, type_)
	if NodeUtil_IsLhsOfAssign(getProp) || NodeUtil_IsStatement(getProp.GetParent()) {
		this_.outer_.graphNodeFactory.CreateNode(type_)
	}
}

// Java method: ProcessPropertiesAndConstructors.processCall(com.google.javascript.rhino.Node call)
func (this_ *AmbiguateProperties_ProcessPropertiesAndConstructors) priv_ProcessPropertiesAndConstructors_processCall_6217(call pkg_rhino.Node) {
	_ = call
	var target pkg_rhino.Node = call.GetFirstChild()
	_ = target
	if !(target.IsQualifiedName()) {
		return
	}
	if this_.outer_.compiler.GetCodingConvention().IsPropertyRenameFunction(target) {
		var propName pkg_rhino.Node = call.GetSecondChild()
		_ = propName
		if (propName == 0) || !(propName.IsStringLit()) {
			return
		}
		var p *AmbiguateProperties_Property = this_.outer_.priv_AmbiguateProperties_getProperty_6203(propName.GetString())
		_ = p
		{
			r_26 := p
			r_26.skipAmbiguating = true
		}
	} else if NodeUtil_IsObjectDefinePropertiesDefinition(call) {
		var typeObj pkg_rhino.Node = call.GetSecondChild()
		_ = typeObj
		var type_ *pkg_colors.Color = this_.outer_.priv_AmbiguateProperties_getColor_6204(typeObj)
		_ = type_
		var objectLiteral pkg_rhino.Node = typeObj.GetNext()
		_ = objectLiteral
		if !(objectLiteral.IsObjectLit()) {
			return
		}
		{
			var key pkg_rhino.Node = objectLiteral.GetFirstChild()
			_ = key
			for ; key != 0; key = key.GetNext() {
				this_.priv_ProcessPropertiesAndConstructors_processObjectProperty_6218(objectLiteral, key, type_)
			}
		}
	}
}

// Java method: ProcessPropertiesAndConstructors.processObjectProperty(com.google.javascript.rhino.Node objectLit, com.google.javascript.rhino.Node key, com.google.javascript.jscomp.colors.Color type)
func (this_ *AmbiguateProperties_ProcessPropertiesAndConstructors) priv_ProcessPropertiesAndConstructors_processObjectProperty_6218(objectLit pkg_rhino.Node, key pkg_rhino.Node, type_ *pkg_colors.Color) {
	_ = objectLit
	_ = key
	_ = type_
	if !(objectLit.IsObjectLit() || objectLit.IsObjectPattern()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(objectLit))
	}
	switch key.TokenSwitchOrd_() {
	case int32(129) /* Token.COMPUTED_PROP */ :
		if key.GetFirstChild().IsStringLit() {
			j2gort.SetAddStr(this_.outer_.quotedNames, key.GetFirstChild().GetString())
		}
	case int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */, int32(100) /* Token.STRING_KEY */ :
		if key.IsQuotedStringKey() {
			j2gort.SetAddStr(this_.outer_.quotedNames, key.GetString())
		} else {
			this_.priv_ProcessPropertiesAndConstructors_maybeMarkCandidate_6222(key, type_)
		}
	case int32(126) /* Token.OBJECT_REST */, int32(128) /* Token.OBJECT_SPREAD */ :
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(("Unexpected child of " + j2gort.StrVal(objectLit.GetToken()) + ": " + j2gort.StrValString(key.ToStringTree()))))))
	}
}

// Java method: ProcessPropertiesAndConstructors.processObjectLitOrPattern(com.google.javascript.rhino.Node objectLit)
func (this_ *AmbiguateProperties_ProcessPropertiesAndConstructors) priv_ProcessPropertiesAndConstructors_processObjectLitOrPattern_6219(objectLit pkg_rhino.Node) {
	_ = objectLit
	if objectLit.GetParent().IsCall() && NodeUtil_IsObjectDefinePropertiesDefinition(objectLit.GetParent()) {
		return
	}
	var type_ *pkg_colors.Color = this_.outer_.priv_AmbiguateProperties_getColor_6204(objectLit)
	_ = type_
	{
		var key pkg_rhino.Node = objectLit.GetFirstChild()
		_ = key
		for ; key != 0; key = key.GetNext() {
			this_.priv_ProcessPropertiesAndConstructors_processObjectProperty_6218(objectLit, key, type_)
		}
	}
}

// Java method: ProcessPropertiesAndConstructors.processGetElem(com.google.javascript.rhino.Node n)
func (this_ *AmbiguateProperties_ProcessPropertiesAndConstructors) priv_ProcessPropertiesAndConstructors_processGetElem_6220(n pkg_rhino.Node) {
	_ = n
	var child pkg_rhino.Node = n.GetLastChild()
	_ = child
	if child.IsStringLit() {
		j2gort.SetAddStr(this_.outer_.quotedNames, child.GetString())
	}
}

// Java method: ProcessPropertiesAndConstructors.processClass(com.google.javascript.rhino.Node classNode)
func (this_ *AmbiguateProperties_ProcessPropertiesAndConstructors) priv_ProcessPropertiesAndConstructors_processClass_6221(classNode pkg_rhino.Node) {
	_ = classNode
	var classConstructorType *pkg_colors.Color = this_.outer_.priv_AmbiguateProperties_getColor_6204(classNode)
	_ = classConstructorType
	this_.outer_.graphNodeFactory.CreateNode(classConstructorType)
	var possiblePrototypes []*pkg_colors.Color = classConstructorType.GetPrototypesArray()
	_ = possiblePrototypes
	var classPrototype *pkg_colors.Color
	_ = classPrototype
	if int32(len(possiblePrototypes)) == int32(0) {
		classPrototype = pkg_colors.StaticGet_StandardColors_UNKNOWN()
	} else {
		classPrototype = pkg_colors.Color_CreateUnion_1_Color(possiblePrototypes)
	}
	{
		var member pkg_rhino.Node = NodeUtil_GetClassMembers(classNode).GetFirstChild()
		_ = member
		for ; member != 0; member = member.GetNext() {
			if member.IsQuotedStringKey() {
				j2gort.SetAddStr(this_.outer_.quotedNames, member.GetString())
				continue
			} else if member.IsComputedProp() || member.IsComputedFieldDef() {
				if member.GetFirstChild().IsStringLit() {
					j2gort.SetAddStr(this_.outer_.quotedNames, member.GetFirstChild().GetString())
				}
				continue
			} else if NodeUtil_IsEs6ConstructorMemberFunctionDef(member) {
				continue
			} else if member.IsBlock() {
				continue
			}
			var memberOwnerColor *pkg_colors.Color = nil
			_ = memberOwnerColor
			if member.IsStaticMember() {
				memberOwnerColor = classConstructorType
			} else if member.IsMemberFieldDef() {
				var possibleInstances []*pkg_colors.Color = classConstructorType.GetInstanceColorsArray()
				_ = possibleInstances
				if int32(len(possibleInstances)) == int32(0) {
					memberOwnerColor = pkg_colors.StaticGet_StandardColors_UNKNOWN()
				} else {
					memberOwnerColor = pkg_colors.Color_CreateUnion_1_Color(possibleInstances)
				}
			} else {
				if !((member.IsMemberFunctionDef() || member.IsGetterDef()) || member.IsSetterDef()) {
					j2gort.CheckState(false)
				}
				memberOwnerColor = classPrototype
			}
			this_.priv_ProcessPropertiesAndConstructors_maybeMarkCandidate_6222(member, memberOwnerColor)
		}
	}
}

// Java method: ProcessPropertiesAndConstructors.maybeMarkCandidate(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.colors.Color type)
func (this_ *AmbiguateProperties_ProcessPropertiesAndConstructors) priv_ProcessPropertiesAndConstructors_maybeMarkCandidate_6222(n pkg_rhino.Node, type_ *pkg_colors.Color) {
	_ = n
	_ = type_
	var name string = n.GetString()
	_ = name
	if !(j2gort.SetContainsStr(this_.outer_.externedNames, name)) {
		this_.outer_.stringNodesToRename = append(this_.outer_.stringNodesToRename, n)
		this_.priv_ProcessPropertiesAndConstructors_recordProperty_6223(name, type_)
	}
}

// Java method: ProcessPropertiesAndConstructors.recordProperty(java.lang.String name, com.google.javascript.jscomp.colors.Color color)
func (this_ *AmbiguateProperties_ProcessPropertiesAndConstructors) priv_ProcessPropertiesAndConstructors_recordProperty_6223(name string, color *pkg_colors.Color) *AmbiguateProperties_Property {
	_ = name
	_ = color
	var prop *AmbiguateProperties_Property = this_.outer_.priv_AmbiguateProperties_getProperty_6203(name)
	_ = prop
	prop.AddRelatedColor(color)
	return prop
}

// Java class: com.google.javascript.jscomp.disambiguate.AmbiguateProperties.Property
type AmbiguateProperties_Property struct {
	outer_             *AmbiguateProperties
	oldName            string                                     `gson:"oldName"`
	newName            string                                     `gson:"newName"`
	relatedColorsSeeds *j2gort.OrderedMap[*ColorGraphNode, int32] `gson:"relatedColorsSeeds"`
	relatedColors      *j2gort.Ext_java_util_BitSet               `gson:"relatedColors"`
	numOccurrences     int32                                      `gson:"numOccurrences"`
	idHash_            atomic.Int32
	skipAmbiguating    bool `gson:"skipAmbiguating"`
}

func (this_ *AmbiguateProperties_Property) BindSelf_(s any) {
	_ = s
}
func (this_ *AmbiguateProperties_Property) IsJavaNil_() bool { return this_ == nil }
func (*AmbiguateProperties_Property) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.disambiguate.AmbiguateProperties$Property", "Property"
}
func (this_ *AmbiguateProperties_Property) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AmbiguateProperties_Property) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.disambiguate.AmbiguateProperties.Property", "com.google.javascript.jscomp.disambiguate.AmbiguateProperties$Property", func(v_ any) bool {
		return j2gort.InstanceOf[*AmbiguateProperties_Property](v_)
	})
}
func (this_ *AmbiguateProperties_Property) SetOuter_(o *AmbiguateProperties) { this_.outer_ = o }
func (this_ *AmbiguateProperties_Property) GetOuter_() *AmbiguateProperties  { return this_.outer_ }

func (this_ *AmbiguateProperties_Property) initFields_AmbiguateProperties_Property() {
	_ = this_
	this_.oldName = j2gort.NullString
	this_.newName = j2gort.NullString

	this_.relatedColors = j2gort.NewExt_java_util_BitSet()
}

// Java constructor: Property(java.lang.String name)
func NewAmbiguateProperties_Property(outer_ *AmbiguateProperties, name string) *AmbiguateProperties_Property {
	this_ := &AmbiguateProperties_Property{}
	this_.outer_ = outer_
	InitAmbiguateProperties_Property(this_, name)
	return this_
}

func InitAmbiguateProperties_Property(this_ *AmbiguateProperties_Property, name string) {
	_ = name
	this_.initFields_AmbiguateProperties_Property()
	this_.oldName = name
}

// Java method: Property.addRelatedColor(com.google.javascript.jscomp.colors.Color color)
func (this_ *AmbiguateProperties_Property) AddRelatedColor(color *pkg_colors.Color) {
	_ = color
	if this_.skipAmbiguating {
		return
	}
	this_.numOccurrences = this_.numOccurrences + 1
	if color.IsInvalidating() || color.GetPropertiesKeepOriginalName() {
		this_.skipAmbiguating = true
		return
	}
	if this_.relatedColorsSeeds == nil {
		this_.relatedColorsSeeds = &j2gort.OrderedMap[*ColorGraphNode, int32]{}
	}
	var newColorGraphNode *ColorGraphNode = this_.outer_.graphNodeFactory.CreateNode(color)
	_ = newColorGraphNode
	this_.relatedColorsSeeds.Put(newColorGraphNode, int32(0))
}
func (*AmbiguateProperties_ProcessPropertiesAndConstructors) AlwaysShouldTraverse_()      {}
func (*AmbiguateProperties_ProcessPropertiesAndConstructors) CanReleaseSyntacticScopes_() {}
