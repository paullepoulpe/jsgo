// Port of com/google/javascript/jscomp/disambiguate/DisambiguateProperties.java @ closure-compiler v20261021

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

package jscomp

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/diagnostic"
	"github.com/paullepoulpe/jsgo/internal/graph"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// DisambiguateProperties assembles the various parts of the diambiguator to
// execute them as a compiler pass.
//
// Java: public final class DisambiguateProperties implements CompilerPass
type DisambiguateProperties struct {
	compiler                       AbstractCompiler
	propertiesThatMustDisambiguate *jutil.LinkedHashSet[string]
	registry                       *colors.ColorRegistry

	mayHavePropertySeenSet *jutil.ObjectLinkedHashSet[*colors.Color]
}

// Java: DisambiguateProperties.PROPERTY_INVALIDATION
var DisambiguateProperties_PROPERTY_INVALIDATION = DiagnosticType_Error(
	"JSC_DISAMBIGUATE2_PROPERTY_INVALIDATION",
	"Property ''{0}'' was required to be disambiguated but was invalidated."+
		"{1}")

var disambiguateProperties_GSON = jutil.NewGson()

// Java: DisambiguateProperties(AbstractCompiler, ImmutableSet<String>)
func NewDisambiguateProperties(
	compiler AbstractCompiler, propertiesThatMustDisambiguate *jutil.LinkedHashSet[string]) *DisambiguateProperties {
	d := &DisambiguateProperties{
		mayHavePropertySeenSet: jutil.NewObjectLinkedHashSet[*colors.Color](),
	}
	d.compiler = compiler
	d.propertiesThatMustDisambiguate = propertiesThatMustDisambiguate
	d.registry = d.compiler.GetColorRegistry()
	return d
}

// Java: DisambiguateProperties.process(Node, Node)
func (d *DisambiguateProperties) Process(externs rhino.Node, root rhino.Node) {
	jutil.CheckArgument(externs.GetParent() == root.GetParent())

	flattener := ColorGraphNodeFactory_CreateFactory(d.registry)
	findRefs := NewColorFindPropertyReferences(
		flattener, d.compiler.GetCodingConvention().IsPropertyRenameFunction)
	graphBuilder := NewColorGraphBuilder(
		flattener,
		func(g graph.DiGraph[*ColorGraphNode, ColorGraphBuilder_EdgeReason]) graph.LowestCommonAncestorFinder[*ColorGraphNode, ColorGraphBuilder_EdgeReason] {
			return graph.NewLowestCommonAncestorFinder(g)
		},
		d.registry)
	propagator := NewClusterPropagator()
	renamingIndexLog := d.compiler.CreateOrReopenLog("DisambiguateProperties", "renaming_index.log")
	defer renamingIndexLog.Close()
	renamer := NewUseSiteRenamer( /* mutationCb= */ d.compiler.ReportChangeToEnclosingScope)
	renamer.trackRenamingIndex = renamingIndexLog.IsLogging()

	NodeTraversal_Traverse(d.compiler, externs.GetParent(), findRefs)
	propIndex := findRefs.GetPropertyIndex()
	findRefs = nil

	DisambiguateProperties_InvalidateWellKnownProperties(propIndex)
	d.LogForDiagnosticsProducer(
		"prop_refs",
		// use a StreamedJsonProducer instead of building up the entire json string at once to
		// prevent OOMs for very large projects.
		disambiguateProperties_propRefsProducer{propIndex: propIndex})

	graphBuilder.AddAll(flattener.typeIndex.Values())
	colorGraph := graphBuilder.Build()
	graphBuilder = nil

	// Model legacy behavior from the old (pre-January 2021) disambiguator.
	// TODO(b/177695515): delete this section.
	for colorGraphNode := range flattener.typeIndex.Values() {
		if len(colorGraph.GetOutEdges(colorGraphNode)) == 0 {
			// Skipping leaf types improves code size, especially as "namespace" types are all leaf
			// types and will often have their declared properties collapsed into variables.
			continue
		}
		if colorGraphNode.GetColor().IsUnion() {
			// Only need this step for SINGLE ColorGraphNodes because
			// and we will add properties of each alternate elsewhere in this loop.
			continue
		}
		d.RegisterOwnDeclaredProperties(colorGraphNode, propIndex)
	}

	d.LogForDiagnostics(
		"graph",
		func() any {
			var nodes []*disambiguateProperties_TypeNodeJson
			for _, n := range colorGraph.GetDiGraphNodes() {
				nodes = append(nodes, newDisambiguateProperties_TypeNodeJson(n))
			}
			slices.SortStableFunc(nodes, func(a, b *disambiguateProperties_TypeNodeJson) int { return cmp.Compare(a.Index, b.Index) })
			return nodes
		})

	d.InvalidateBasedOnType(flattener)
	d.mayHavePropertySeenSet = jutil.NewObjectLinkedHashSet[*colors.Color]()

	graph.FixedPointGraphTraversal_NewTraversal[*ColorGraphNode, ColorGraphBuilder_EdgeReason](propagator).ComputeFixedPoint(colorGraph)
	flattener = nil
	colorGraph = nil
	propagator = nil

	trackerSummaryGenerator := &DisambiguateProperties_TrackerSummaryGenerator{}
	for prop := range propIndex.Values() {
		renameUsesResult := renamer.RenameUses(prop)
		trackerSummaryGenerator.AddRenameUsesResult(renameUsesResult)
		if prop.IsInvalidated() && d.propertiesThatMustDisambiguate.Contains(prop.GetName()) {
			d.compiler.Report(d.CreateInvalidationError(prop))
		}
	}
	d.compiler.ReportDisambiguatePropertiesSummary(trackerSummaryGenerator.BuildSummaryString)

	renamingIndexLog.LogJsonSupplier(func() any { return DisambiguateProperties_BuildRenamingIndex(propIndex, renamer) })
	propIndex = nil
	renamer = nil
	d.LogForDiagnostics("mismatches", func() any { return d.registry.GetMismatchLocationsForDebugging().Inverse() })

	GatherGetterAndSetterProperties_Update(d.compiler, externs, root)
}

// disambiguateProperties_propRefsProducer is the anonymous
// LogFile.StreamedJsonProducer of process.
type disambiguateProperties_propRefsProducer struct {
	propIndex *jutil.LinkedHashMap[string, *PropertyClustering]
}

// Java: DisambiguateProperties.process(...).new StreamedJsonProducer() {...}.writeJson(JsonWriter)
func (p disambiguateProperties_propRefsProducer) WriteJson(writer *jutil.JsonWriter) {
	var propRefsJson []*disambiguateProperties_PropertyReferenceIndexJson
	for prop := range p.propIndex.Values() {
		propRefsJson = append(propRefsJson, newDisambiguateProperties_PropertyReferenceIndexJson(prop))
	}
	// collect(toImmutableSortedSet(naturalOrder()))
	slices.SortStableFunc(propRefsJson, (*disambiguateProperties_PropertyReferenceIndexJson).CompareTo)
	propRefsJson = slices.CompactFunc(propRefsJson, func(a, b *disambiguateProperties_PropertyReferenceIndexJson) bool { return a.CompareTo(b) == 0 })

	writer.BeginObject()
	for _, propRef := range propRefsJson {
		propRef.WriteJson(writer)
	}
	writer.EndObject()
}

// DisambiguateProperties_TrackerSummaryGenerator is a no-op class to use when
// tracer mode is not enabled.
//
// Java: private static class DisambiguateProperties.TrackerSummaryGenerator
type DisambiguateProperties_TrackerSummaryGenerator struct {
	total             int
	numInvalidated    int
	numDisambiguated  int
	numOnlyOneCluster int
}

// Java: DisambiguateProperties.TrackerSummaryGenerator.addRenameUsesResult(UseSiteRenamer.RenameUsesResult)
func (g *DisambiguateProperties_TrackerSummaryGenerator) AddRenameUsesResult(renameUsesResult UseSiteRenamer_RenameUsesResult) {
	g.total++
	switch renameUsesResult {
	case UseSiteRenamer_RenameUsesResult_INVALIDATED:
		g.numInvalidated++
	case UseSiteRenamer_RenameUsesResult_ONLY_ONE_CLUSTER:
		g.numOnlyOneCluster++
	case UseSiteRenamer_RenameUsesResult_DISAMBIGUATED:
		g.numDisambiguated++
	}
}

// Java: DisambiguateProperties.TrackerSummaryGenerator.buildSummaryString()
func (g *DisambiguateProperties_TrackerSummaryGenerator) BuildSummaryString() string {
	return fmt.Sprintf(
		"%d property names, %d disambiguated, %d invalidated, %d had a single cluster",
		g.total, g.numDisambiguated, g.numInvalidated, g.numOnlyOneCluster)
}

// Java: DisambiguateProperties.createInvalidationError(PropertyClustering)
func (d *DisambiguateProperties) CreateInvalidationError(prop *PropertyClustering) *JSError {
	additionalContext := ""

	return JSError_MakeWithSourceLocation(
		nil,
		-1,
		-1,
		DisambiguateProperties_PROPERTY_INVALIDATION,
		prop.GetName(),
		additionalContext)
}

// Java: DisambiguateProperties.buildRenamingIndex(LinkedHashMap<String, PropertyClustering>, UseSiteRenamer)
func DisambiguateProperties_BuildRenamingIndex(
	props *jutil.LinkedHashMap[string, *PropertyClustering], renamer *UseSiteRenamer) *jutil.LinkedHashMap[string, any] {
	newNames := renamer.GetRenamingIndex()
	// collect(toImmutableSortedMap(naturalOrder(), getName, ...)): duplicate keys throw.
	byName := jutil.NewLinkedHashMap[string, any]()
	for prop := range props.Values() {
		var value any
		if prop.IsInvalidated() {
			value = prop.GetLastInvalidation()
		} else {
			names := newNames.Get(prop.GetName()).ToSlice()
			slices.SortStableFunc(names, jutil.StringCompareTo)
			value = slices.Compact(names) // ImmutableSortedSet.copyOf
		}
		if byName.ContainsKey(prop.GetName()) {
			panic(jutil.NewIllegalArgumentExceptionMsg("Duplicate keys in toImmutableSortedMap: " + prop.GetName()))
		}
		byName.Put(prop.GetName(), value)
	}
	keys := byName.KeySlice()
	slices.SortStableFunc(keys, jutil.StringCompareTo)
	sorted := jutil.NewLinkedHashMap[string, any]()
	for _, key := range keys {
		sorted.Put(key, byName.Get(key))
	}
	return sorted
}

// Java: DisambiguateProperties.invalidateWellKnownProperties(LinkedHashMap<String, PropertyClustering>)
func DisambiguateProperties_InvalidateWellKnownProperties(
	propIndex *jutil.LinkedHashMap[string, *PropertyClustering]) {
	/*
	 * Expand this list as needed; it wasn't created exhaustively.
	 *
	 * <p>Good candidates are: props accessed by builtin functions, props accessed by syntax sugar,
	 * props used in strange ways by the language spec, etc.
	 *
	 * <p>TODO(b/169899789): consider instead omitting these properties entirely from the serialized
	 * colors format.
	 */
	names := []string{"prototype", "constructor", "then"}
	for _, name := range names {
		propIndex.
			ComputeIfAbsent(name, NewPropertyClustering).
			Invalidate(Invalidation_WellKnownProperty())
	}
}

// Java: DisambiguateProperties.invalidateBasedOnType(ColorGraphNodeFactory)
func (d *DisambiguateProperties) InvalidateBasedOnType(flattener *ColorGraphNodeFactory) {
	for type_ := range flattener.typeIndex.Values() {
		if type_.GetColor().IsInvalidating() {
			inv := Invalidation_InvalidatingType(type_.GetIndex())
			for prop := range type_.GetAssociatedProps().Keys() {
				prop.Invalidate(inv)
			}
		} else {
			d.InvalidateNonDeclaredPropertyAccesses(type_)
		}
	}
}

// InvalidateNonDeclaredPropertyAccesses invalidates all property accesses
// that cause "missing property" warnings.
//
// This behavior is not inherently necessary for correctness. It only exists
// because the older version of the disambiguator invalidated these
// properties and some projects began to rely on this.
//
// TODO(b/177695515): delete this method
//
// Java: DisambiguateProperties.invalidateNonDeclaredPropertyAccesses(ColorGraphNode)
func (d *DisambiguateProperties) InvalidateNonDeclaredPropertyAccesses(colorGraphNode *ColorGraphNode) {
	color := colorGraphNode.GetColor()
	jutil.CheckArgumentMsg(
		!color.IsInvalidating(),
		"Not applicable to invalidating types. All their properties are invalidated")

	for prop := range colorGraphNode.GetAssociatedProps().Keys() {
		if prop.IsInvalidated() {
			continue // Skip unnecessary `hasProperty` lookups which can be expensive.
		}

		if !d.MayHaveProperty(color, prop.GetName()) {
			prop.Invalidate(Invalidation_UndeclaredAccess(colorGraphNode.GetIndex()))
		}
	}
}

// MayHaveProperty returns true if the color or any of its ancestors has the
// given property.
//
// If this is a union, returns true if /any/ union alternate has the property.
//
// TODO(b/177695515): delete this method
//
// Java: DisambiguateProperties.mayHaveProperty(Color, String)
func (d *DisambiguateProperties) MayHaveProperty(color *colors.Color, propertyName string) bool {
	defer d.mayHavePropertySeenSet.Remove(color) // finally
	if !d.mayHavePropertySeenSet.Add(color) {
		return false
	}

	if color.IsUnion() {
		for element := range color.GetUnionElements().All() {
			if d.MayHaveProperty(element, propertyName) {
				return true
			}
		}
		return false
	}

	if color.GetOwnProperties().Contains(propertyName) {
		return true
	}
	for element := range d.registry.GetDisambiguationSupertypes(color).All() {
		if d.MayHaveProperty(element, propertyName) {
			return true
		}
	}
	return false
}

// Java: DisambiguateProperties.registerOwnDeclaredProperties(ColorGraphNode, Map<String, PropertyClustering>)
func (d *DisambiguateProperties) RegisterOwnDeclaredProperties(
	colorGraphNode *ColorGraphNode, propIndex *jutil.LinkedHashMap[string, *PropertyClustering]) {
	jutil.CheckArgument(!colorGraphNode.GetColor().IsUnion())

	// For each type, get the list of its "own properties" and add them to its clustering. This
	// is only to mimic the behavior of the old disambiguator and could be removed if we were
	// confident we could update all existing code to be compatible.
	for propName := range colorGraphNode.GetColor().GetOwnProperties().All() {
		prop := propIndex.Get(propName)
		if prop == nil {
			// Ignore declared properties without any visible references to rename.
			continue
		} else if prop.IsInvalidated() {
			continue
		}
		prop.GetClusters().Add(colorGraphNode)
		colorGraphNode.GetAssociatedProps().Put(prop, ColorGraphNode_PropAssociation_TYPE_SYSTEM)
	}
}

// Java: DisambiguateProperties.logForDiagnostics(String, Supplier<Object>)
func (d *DisambiguateProperties) LogForDiagnostics(name string, data func() any) {
	log := d.compiler.CreateOrReopenLog("DisambiguateProperties", name+".log")
	defer log.Close()
	log.LogJsonSupplier(data)
}

// Java: DisambiguateProperties.logForDiagnostics(String, LogFile.StreamedJsonProducer)
func (d *DisambiguateProperties) LogForDiagnosticsProducer(name string, data diagnostic.LogFile_StreamedJsonProducer) {
	log := d.compiler.CreateOrReopenLog("DisambiguateProperties", name+".log")
	defer log.Close()
	log.LogJsonProducer(data)
}

// Java: private static final class DisambiguateProperties.PropertyReferenceIndexJson implements Comparable<PropertyReferenceIndexJson>
type disambiguateProperties_PropertyReferenceIndexJson struct {
	name string
	refs []*disambiguateProperties_PropertyReferenceJson
}

// Java: DisambiguateProperties.PropertyReferenceIndexJson(PropertyClustering)
func newDisambiguateProperties_PropertyReferenceIndexJson(prop *PropertyClustering) *disambiguateProperties_PropertyReferenceIndexJson {
	j := &disambiguateProperties_PropertyReferenceIndexJson{}
	j.name = prop.GetName()
	if prop.IsInvalidated() {
		j.refs = []*disambiguateProperties_PropertyReferenceJson{}
	} else {
		for key, value := range prop.GetUseSites().All() {
			j.refs = append(j.refs, newDisambiguateProperties_PropertyReferenceJson(key, value))
		}
		// collect(toImmutableSortedSet(naturalOrder()))
		slices.SortStableFunc(j.refs, (*disambiguateProperties_PropertyReferenceJson).CompareTo)
		j.refs = slices.CompactFunc(j.refs, func(a, b *disambiguateProperties_PropertyReferenceJson) bool { return a.CompareTo(b) == 0 })
	}
	return j
}

// Java: DisambiguateProperties.PropertyReferenceIndexJson.compareTo(DisambiguateProperties.PropertyReferenceIndexJson)
func (j *disambiguateProperties_PropertyReferenceIndexJson) CompareTo(x *disambiguateProperties_PropertyReferenceIndexJson) int {
	return jutil.StringCompareTo(j.name, x.name)
}

// Java: DisambiguateProperties.PropertyReferenceIndexJson.writeJson(JsonWriter)
func (j *disambiguateProperties_PropertyReferenceIndexJson) WriteJson(writer *jutil.JsonWriter) {
	// creates an entry such as:
	//   foo: {name: 'foo', refs: [{location: 'bar.js:3:4', receiverIndex: 2}]}
	writer.Name(j.name)

	writer.BeginObject()
	writer.Name("name").ValueString(&j.name)
	writer.Name("refs").BeginArray()
	for _, ref := range j.refs {
		disambiguateProperties_GSON.ToJsonWithWriter(ref, writer)
	}
	writer.EndArray()
	writer.EndObject()
}

// Java: private static final class DisambiguateProperties.PropertyReferenceJson implements Comparable<PropertyReferenceJson>
type disambiguateProperties_PropertyReferenceJson struct {
	Location      string `gson:"location"`
	ReceiverIndex int    `gson:"receiverIndex"`
}

// Java: DisambiguateProperties.PropertyReferenceJson(Node, ColorGraphNode)
func newDisambiguateProperties_PropertyReferenceJson(location rhino.Node, receiver *ColorGraphNode) *disambiguateProperties_PropertyReferenceJson {
	return &disambiguateProperties_PropertyReferenceJson{
		Location:      location.GetLocation(),
		ReceiverIndex: receiver.GetIndex(),
	}
}

// Java: DisambiguateProperties.PropertyReferenceJson.compareTo(DisambiguateProperties.PropertyReferenceJson)
func (j *disambiguateProperties_PropertyReferenceJson) CompareTo(x *disambiguateProperties_PropertyReferenceJson) int {
	// ComparisonChain.start().compare(receiverIndex).compare(location).result()
	if c := cmp.Compare(j.ReceiverIndex, x.ReceiverIndex); c != 0 {
		return c
	}
	return jutil.StringCompareTo(j.Location, x.Location)
}

// Java: private static final class DisambiguateProperties.TypeNodeJson
type disambiguateProperties_TypeNodeJson struct {
	Index        int                                                          `gson:"index"`
	Invalidating bool                                                         `gson:"invalidating"`
	ColorId      string                                                       `gson:"colorId"`
	Edges        []*disambiguateProperties_TypeEdgeJson                       `gson:"edges"`
	Props        *jutil.LinkedHashMap[string, ColorGraphNode_PropAssociation] `gson:"props"`
}

// Java: DisambiguateProperties.TypeNodeJson(DiGraph.DiGraphNode<ColorGraphNode, Object>)
func newDisambiguateProperties_TypeNodeJson(n graph.DiGraph_DiGraphNode[*ColorGraphNode, ColorGraphBuilder_EdgeReason]) *disambiguateProperties_TypeNodeJson {
	t := n.GetValue()

	j := &disambiguateProperties_TypeNodeJson{}
	j.Index = t.GetIndex()
	j.ColorId = t.GetColor().GetId().String()
	j.Invalidating = t.GetColor().IsInvalidating()
	for _, e := range n.GetOutEdges() {
		j.Edges = append(j.Edges, newDisambiguateProperties_TypeEdgeJson(e))
	}
	// collect(toImmutableSortedSet(naturalOrder()))
	slices.SortStableFunc(j.Edges, (*disambiguateProperties_TypeEdgeJson).CompareTo)
	j.Edges = slices.CompactFunc(j.Edges, func(a, b *disambiguateProperties_TypeEdgeJson) bool { return a.CompareTo(b) == 0 })
	// collect(toImmutableSortedMap(naturalOrder(), e -> e.getKey().getName(), Map.Entry::getValue))
	var names []string
	byName := jutil.NewLinkedHashMap[string, ColorGraphNode_PropAssociation]()
	for prop, association := range t.GetAssociatedProps().All() {
		if byName.ContainsKey(prop.GetName()) {
			panic(jutil.NewIllegalArgumentExceptionMsg("Duplicate keys in toImmutableSortedMap: " + prop.GetName()))
		}
		byName.Put(prop.GetName(), association)
		names = append(names, prop.GetName())
	}
	slices.SortStableFunc(names, jutil.StringCompareTo)
	j.Props = jutil.NewLinkedHashMap[string, ColorGraphNode_PropAssociation]()
	for _, name := range names {
		j.Props.Put(name, byName.Get(name))
	}
	return j
}

// Java: private static final class DisambiguateProperties.TypeEdgeJson implements Comparable<TypeEdgeJson>
type disambiguateProperties_TypeEdgeJson struct {
	Dest  int                          `gson:"dest"`
	Value ColorGraphBuilder_EdgeReason `gson:"value"`
}

// Java: DisambiguateProperties.TypeEdgeJson(DiGraph.DiGraphEdge<ColorGraphNode, Object>)
func newDisambiguateProperties_TypeEdgeJson(e graph.DiGraph_DiGraphEdge[*ColorGraphNode, ColorGraphBuilder_EdgeReason]) *disambiguateProperties_TypeEdgeJson {
	return &disambiguateProperties_TypeEdgeJson{
		Dest:  e.GetDestination().GetValue().GetIndex(),
		Value: e.GetValue(),
	}
}

// Java: DisambiguateProperties.TypeEdgeJson.compareTo(DisambiguateProperties.TypeEdgeJson)
func (j *disambiguateProperties_TypeEdgeJson) CompareTo(x *disambiguateProperties_TypeEdgeJson) int {
	jutil.CheckArgument(j.Dest != x.Dest)
	return j.Dest - x.Dest
}
