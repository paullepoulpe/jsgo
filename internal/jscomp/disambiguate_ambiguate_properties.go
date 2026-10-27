// Port of com/google/javascript/jscomp/disambiguate/AmbiguateProperties.java @ closure-compiler v20261021

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

package jscomp

import (
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/graph"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// AmbiguateProperties renames unrelated properties to the same name, using
// Colors provided by the typechecker. This allows better compression as more
// properties can be given short names.
//
// Properties are considered unrelated if they are never referenced from the
// same color or from a subtype of each others' colors, thus this pass is only
// effective if type checking is enabled.
//
// Example:
//
//	Foo.fooprop = 0;
//	Foo.fooprop2 = 0;
//	Bar.barprop = 0;
//
// becomes:
//
//	Foo.a = 0;
//	Foo.b = 0;
//	Bar.a = 0;
//
// Java: public class AmbiguateProperties implements CompilerPass
type AmbiguateProperties struct {
	compiler AbstractCompiler

	stringNodesToRename []rhino.Node
	// Can't use these to start property names.
	reservedFirstCharacters jutil.Set[uint16]
	// Can't use these at all in property names.
	reservedNonFirstCharacters jutil.Set[uint16]

	// Map from property name to Property object
	propertyMap *jutil.LinkedHashMap[string, *AmbiguateProperties_Property]

	// Property names that don't get renamed
	externedNames *jutil.LinkedHashSet[string]

	// Names to which properties shouldn't be renamed, to avoid name conflicts
	quotedNames *jutil.LinkedHashSet[string]

	colorRegistry *colors.ColorRegistry

	// Map from original property name to new name. Only used by tests.
	renamingMap *jutil.LinkedHashMap[string, string] // may be nil

	graphNodeFactory *ColorGraphNodeFactory // may be nil
}

var ambiguateProperties_logger = jutil.Logger_GetLogger("com.google.javascript.jscomp.disambiguate.AmbiguateProperties")

// ambiguateProperties_FREQUENCY_COMPARATOR sorts Property objects by their
// count, breaking ties alphabetically to ensure a deterministic total
// ordering.
//
// Java: AmbiguateProperties.FREQUENCY_COMPARATOR
var ambiguateProperties_FREQUENCY_COMPARATOR = func(p1 *AmbiguateProperties_Property, p2 *AmbiguateProperties_Property) int {
	if p1.numOccurrences != p2.numOccurrences {
		return p2.numOccurrences - p1.numOccurrences
	}
	return jutil.StringCompareTo(p1.oldName, p2.oldName)
}

// Java: AmbiguateProperties(AbstractCompiler, Set<Character>, Set<Character>, Set<String>)
func NewAmbiguateProperties(
	compiler AbstractCompiler,
	reservedFirstCharacters jutil.Set[uint16],
	reservedNonFirstCharacters jutil.Set[uint16],
	externProperties jutil.Set[string]) *AmbiguateProperties {
	jutil.CheckState(compiler.GetLifeCycleStage().IsNormalized())
	// ImmutableSet.<String>builder().add("prototype").addAll(externProperties).build()
	externedNames := jutil.NewLinkedHashSetWithCapacity[string](externProperties.Size() + 1)
	externedNames.Add("prototype")
	for name := range externProperties.All() {
		externedNames.Add(name)
	}
	a := &AmbiguateProperties{
		compiler:                   compiler,
		reservedFirstCharacters:    reservedFirstCharacters,
		reservedNonFirstCharacters: reservedNonFirstCharacters,
		propertyMap:                jutil.NewLinkedHashMapWithCapacity[string, *AmbiguateProperties_Property](2048),
		externedNames:              externedNames,
		quotedNames:                jutil.NewLinkedHashSetWithCapacity[string](256),
		colorRegistry:              compiler.GetColorRegistry(),
		renamingMap:                nil,
		graphNodeFactory:           nil,
	}
	return a
}

// Java: AmbiguateProperties.makePassForTesting(AbstractCompiler, Set<Character>, Set<Character>, Set<String>)
func AmbiguateProperties_MakePassForTesting(
	compiler AbstractCompiler,
	reservedFirstCharacters jutil.Set[uint16],
	reservedNonFirstCharacters jutil.Set[uint16],
	externProperties jutil.Set[string]) *AmbiguateProperties {
	ap := NewAmbiguateProperties(
		compiler, reservedFirstCharacters, reservedNonFirstCharacters, externProperties)
	ap.renamingMap = jutil.NewLinkedHashMap[string, string]()
	return ap
}

// Java: AmbiguateProperties.getRenamingMap()
func (a *AmbiguateProperties) GetRenamingMap() *jutil.LinkedHashMap[string, string] {
	jutil.CheckNotNull(a.renamingMap)
	return a.renamingMap
}

// Java: AmbiguateProperties.process(Node, Node)
func (a *AmbiguateProperties) Process(externs rhino.Node, root rhino.Node) {
	a.graphNodeFactory = ColorGraphNodeFactory_CreateFactory(a.colorRegistry)

	// Find all property references and record the types on which they occur.
	// Populate stringNodesToRename, propertyMap, quotedNames.
	NodeTraversal_Traverse(a.compiler, root, &AmbiguateProperties_ProcessPropertiesAndConstructors{outer: a})

	graphBuilder := NewColorGraphBuilder(
		a.graphNodeFactory,
		func(g graph.DiGraph[*ColorGraphNode, ColorGraphBuilder_EdgeReason]) graph.LowestCommonAncestorFinder[*ColorGraphNode, ColorGraphBuilder_EdgeReason] {
			return graph.NewLowestCommonAncestorFinder(g)
		},
		a.colorRegistry)
	graphBuilder.AddAll(a.graphNodeFactory.typeIndex.Values())
	colorGraph := graphBuilder.Build()
	graphBuilder = nil
	func() {
		chunkGraphLog := a.compiler.CreateOrReopenLog("AmbiguateProperties", "color_graph.dot")
		defer chunkGraphLog.Close()
		if chunkGraphLog.IsLogging() {
			chunkGraphLog.LogString(DotFormatter_ToDotOfGraphvizGraph(colorGraph))
		}
	}()
	for node := range a.graphNodeFactory.typeIndex.Values() {
		node.reflexiveSubtype = true // Init subtyping as reflexive.
	}

	graph.FixedPointGraphTraversal_NewReverseTraversal[*ColorGraphNode, ColorGraphBuilder_EdgeReason](
		graph.FixedPointGraphTraversal_EdgeCallbackFunc[*ColorGraphNode, ColorGraphBuilder_EdgeReason](func(subtype *ColorGraphNode, _ ColorGraphBuilder_EdgeReason, supertype *ColorGraphNode) bool {
			/*
			 * Cheap path for when we're sure there's going to be a change.
			 *
			 * <p>Since bits only ever turn on, using more bits means there are definitely more
			 * elements. This prevents of from needing to check cardinality or equality, which
			 * would otherwise dominate the cost of computing the fixed point.
			 *
			 * <p>We're guaranteed to converge because the sizes will be euqal after the OR
			 * operation.
			 */
			superSet := supertype.GetSubtypeIndices()
			if subtype.subtypeIndices == nil {
				idx := subtype.GetIndex()
				if superSet.Get(idx) {
					return false
				}
				superSet.Set(idx)
				return true
			}
			subSet := subtype.subtypeIndices
			if subSet.Size() > superSet.Size() {
				superSet.Or(subSet)
				return true
			}

			startSize := superSet.Cardinality()
			superSet.Or(subSet)
			return superSet.Cardinality() > startSize
		})).
		ComputeFixedPoint(colorGraph)

	// Fill in all transitive edges in subtyping graph per property
	for prop := range a.propertyMap.Values() {
		if prop.skipAmbiguating || prop.firstColorSeed == nil {
			prop.firstColorSeed = nil
			prop.relatedColorsSeeds = nil
			continue
		}
		if prop.relatedColorsSeeds == nil {
			c := prop.firstColorSeed
			prop.firstColorSeed = nil
			if c.subtypeIndices == nil {
				prop.singleColorIndex = int32(c.index)
			} else {
				prop.relatedColors.Or(c.subtypeIndices)
			}
			continue
		}
		for color := range prop.relatedColorsSeeds.Keys() {
			if color.subtypeIndices != nil {
				prop.relatedColors.Or(color.subtypeIndices)
			} else {
				prop.relatedColors.Set(color.index)
			}
		}
		prop.firstColorSeed = nil
		prop.relatedColorsSeeds = nil
	}
	colorGraph = nil
	a.graphNodeFactory = nil

	// ImmutableSet.<String>builder().addAll(externedNames).addAll(quotedNames)
	numProperties := a.propertyMap.Size()
	reservedNames := jutil.NewLinkedHashSetWithCapacity[string](a.externedNames.Size() + a.quotedNames.Size() + numProperties + 256)
	for name := range a.externedNames.All() {
		reservedNames.Add(name)
	}
	for name := range a.quotedNames.All() {
		reservedNames.Add(name)
	}
	a.externedNames = nil
	a.quotedNames = nil
	numRenamedPropertyNames := 0
	numSkippedPropertyNames := 0
	nodeSlab := make([]AmbiguateProperties_PropertyGraphNode, numProperties)
	nodes := make([]*AmbiguateProperties_PropertyGraphNode, 0, numProperties)
	for prop := range a.propertyMap.Values() {
		if prop.skipAmbiguating {
			numSkippedPropertyNames++
			reservedNames.Add(prop.oldName)
		} else {
			node := &nodeSlab[numRenamedPropertyNames]
			node.property = prop
			numRenamedPropertyNames++
			nodes = append(nodes, node)
		}
	}
	finalNumRenamedPropertyNames := numRenamedPropertyNames
	finalNumSkippedPropertyNames := numSkippedPropertyNames

	propertyGraph := NewAmbiguateProperties_PropertyGraph(nodes)
	var coloring graph.GraphColoring[*AmbiguateProperties_Property, struct{}] = graph.NewGraphColoring_GreedyGraphColoringWithTieBreaker[*AmbiguateProperties_Property, struct{}](
		propertyGraph, ambiguateProperties_FREQUENCY_COMPARATOR)
	numNewPropertyNames := coloring.Color()
	coloring = nil

	// Generate new names for the properties that will be renamed.
	var nameGen NameGenerator = NewDefaultNameGeneratorWithFirstAndNonFirstChars(
		reservedNames,
		"",
		a.reservedFirstCharacters,
		a.reservedNonFirstCharacters)
	colorMap := make([]string, numNewPropertyNames)
	for i := 0; i < numNewPropertyNames; i++ {
		colorMap[i] = nameGen.GenerateNextName()
	}
	nameGen = nil
	reservedNames = nil

	// Translate the color of each Property instance to a name.
	for _, node := range propertyGraph.nodes {
		prop := node.GetValue()
		prop.newName = colorMap[graph.GraphColoring_ColorOf(node.GetAnnotation()).Value]
		if a.renamingMap != nil {
			a.renamingMap.Put(prop.oldName, prop.newName)
		}
	}
	propertyGraph = nil
	nodes = nil
	nodeSlab = nil
	colorMap = nil

	// Actually assign the new names to the relevant STRING nodes in the AST.
	for _, n := range a.stringNodesToRename {
		oldName := n.GetString()
		p := a.propertyMap.Get(oldName)
		if p != nil && p.newName != "" {
			jutil.CheckState(oldName == p.oldName)
			if p.newName != oldName {
				n.SetString(p.newName)
				a.compiler.ReportChangeToEnclosingScope(n)
			}
		}
	}
	a.stringNodesToRename = nil
	a.propertyMap = nil

	// We may have renamed getter / setter properties.
	// TODO(b/161947315): this shouldn't be the responsibility of AmbiguateProperties
	GatherGetterAndSetterProperties_Update(a.compiler, externs, root)

	// Suppliers.memoize(...)
	var summary *string
	summarySupplier := func() string {
		if summary == nil {
			s := "Collapsed " +
				strconv.Itoa(finalNumRenamedPropertyNames) +
				" properties into " +
				strconv.Itoa(numNewPropertyNames) +
				" and skipped renaming " +
				strconv.Itoa(finalNumSkippedPropertyNames) +
				" properties."
			summary = &s
		}
		return *summary
	}
	if ambiguateProperties_logger.IsLoggable(jutil.Level_FINE) {
		ambiguateProperties_logger.Fine(summarySupplier())
	}
	a.compiler.ReportAmbiguatePropertiesSummary(summarySupplier)
}

// Java: static class AmbiguateProperties.PropertyGraph implements AdjacencyGraph<Property, Void>
type AmbiguateProperties_PropertyGraph struct {
	nodes []*AmbiguateProperties_PropertyGraphNode
}

// Java: AmbiguateProperties.PropertyGraph(ArrayList<AmbiguateProperties.PropertyGraphNode>)
func NewAmbiguateProperties_PropertyGraph(nodes []*AmbiguateProperties_PropertyGraphNode) *AmbiguateProperties_PropertyGraph {
	return &AmbiguateProperties_PropertyGraph{nodes: nodes}
}

// Java: AmbiguateProperties.PropertyGraph.getNodes()
func (g *AmbiguateProperties_PropertyGraph) GetNodes() []graph.GraphNode[*AmbiguateProperties_Property, struct{}] {
	// Go-only: language-adapter — List<PropertyGraphNode> as the
	// interface's List<GraphNode<Property, Void>>.
	out := make([]graph.GraphNode[*AmbiguateProperties_Property, struct{}], len(g.nodes))
	for i, n := range g.nodes {
		out[i] = n
	}
	return out
}

// Java: AmbiguateProperties.PropertyGraph.getNodeCount()
func (g *AmbiguateProperties_PropertyGraph) GetNodeCount() int {
	return len(g.nodes)
}

// Java: AmbiguateProperties.PropertyGraph.getNode(AmbiguateProperties.Property)
func (g *AmbiguateProperties_PropertyGraph) GetNode(property *AmbiguateProperties_Property) graph.GraphNode[*AmbiguateProperties_Property, struct{}] {
	panic(jutil.NewRuntimeException("PropertyGraph#getNode is never called."))
}

// Java: AmbiguateProperties.PropertyGraph.newSubGraph()
func (g *AmbiguateProperties_PropertyGraph) NewSubGraph() graph.SubGraph[*AmbiguateProperties_Property, struct{}] {
	return NewAmbiguateProperties_PropertySubGraph()
}

// Java: AmbiguateProperties.PropertyGraph.clearNodeAnnotations()
func (g *AmbiguateProperties_PropertyGraph) ClearNodeAnnotations() {
	for _, node := range g.nodes {
		node.SetAnnotation(nil)
	}
}

// Java: AmbiguateProperties.PropertyGraph.getWeight(AmbiguateProperties.Property)
func (g *AmbiguateProperties_PropertyGraph) GetWeight(value *AmbiguateProperties_Property) int {
	return value.numOccurrences
}

// AmbiguateProperties_PropertySubGraph is a SubGraph that represents
// properties. The related types of the properties are used to efficiently
// calculate adjacency information.
//
// Java: static class AmbiguateProperties.PropertySubGraph implements SubGraph<Property, Void>
type AmbiguateProperties_PropertySubGraph struct {
	// Types related to properties referenced in this subgraph.
	relatedTypes jutil.BitSet
}

// Java: AmbiguateProperties.PropertySubGraph()
func NewAmbiguateProperties_PropertySubGraph() *AmbiguateProperties_PropertySubGraph {
	return &AmbiguateProperties_PropertySubGraph{}
}

// IsIndependentOf returns true if prop is in an independent set from all
// properties in this sub graph. That is, if none of its related types
// intersects with the related types for this sub graph.
//
// Java: AmbiguateProperties.PropertySubGraph.isIndependentOf(AmbiguateProperties.Property)
func (g *AmbiguateProperties_PropertySubGraph) IsIndependentOf(prop *AmbiguateProperties_Property) bool {
	if prop.singleColorIndex >= 0 {
		return !g.relatedTypes.Get(int(prop.singleColorIndex))
	}
	return !g.relatedTypes.Intersects(&prop.relatedColors)
}

// AddNode adds the node to the sub graph, adding all its related types to the
// related types for the sub graph.
//
// Java: AmbiguateProperties.PropertySubGraph.addNode(AmbiguateProperties.Property)
func (g *AmbiguateProperties_PropertySubGraph) AddNode(prop *AmbiguateProperties_Property) {
	if prop.singleColorIndex >= 0 {
		g.relatedTypes.Set(int(prop.singleColorIndex))
		return
	}
	g.relatedTypes.Or(&prop.relatedColors)
}

// Java: static class AmbiguateProperties.PropertyGraphNode implements GraphNode<Property, Void>
type AmbiguateProperties_PropertyGraphNode struct {
	property   *AmbiguateProperties_Property
	annotation graph.Annotation
}

// Java: AmbiguateProperties.PropertyGraphNode(AmbiguateProperties.Property)
func NewAmbiguateProperties_PropertyGraphNode(property *AmbiguateProperties_Property) *AmbiguateProperties_PropertyGraphNode {
	return &AmbiguateProperties_PropertyGraphNode{property: property}
}

// Java: AmbiguateProperties.PropertyGraphNode.getValue()
func (n *AmbiguateProperties_PropertyGraphNode) GetValue() *AmbiguateProperties_Property {
	return n.property
}

// Java: AmbiguateProperties.PropertyGraphNode.getAnnotation()
func (n *AmbiguateProperties_PropertyGraphNode) GetAnnotation() graph.Annotation {
	return n.annotation
}

// Java: AmbiguateProperties.PropertyGraphNode.setAnnotation(Annotation)
func (n *AmbiguateProperties_PropertyGraphNode) SetAnnotation(data graph.Annotation) {
	n.annotation = data
}

// AmbiguateProperties_ProcessPropertiesAndConstructors finds all property
// references, recording the types on which they occur, and records all
// constructors and their instance types in the ColorGraphNodeFactory.
//
// Java: private class AmbiguateProperties.ProcessPropertiesAndConstructors extends AbstractPostOrderCallback
type AmbiguateProperties_ProcessPropertiesAndConstructors struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *AmbiguateProperties
}

func (*AmbiguateProperties_ProcessPropertiesAndConstructors) onlyNameLeaves() {}

// Java: AmbiguateProperties.ProcessPropertiesAndConstructors.visit(NodeTraversal, Node, Node)
func (p *AmbiguateProperties_ProcessPropertiesAndConstructors) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	a := p.outer
	switch n.GetToken() {
	case rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP:
		p.ProcessGetProp(n)
		return
	case rhino.Token_CALL:
		p.ProcessCall(n)
		return
	case rhino.Token_NAME:
		// handle ES5-style classes
		if NodeUtil_IsNameDeclaration(parent) || parent.IsFunction() {
			a.graphNodeFactory.CreateNode(a.GetColor(n))
		}
		return
	case rhino.Token_OBJECTLIT, rhino.Token_OBJECT_PATTERN:
		p.ProcessObjectLitOrPattern(n)
		return
	case rhino.Token_GETELEM:
		p.ProcessGetElem(n)
		return
	case rhino.Token_CLASS:
		p.ProcessClass(n)
		return
	default:
		// Nothing to do.
	}
}

// Java: AmbiguateProperties.ProcessPropertiesAndConstructors.processGetProp(Node)
func (p *AmbiguateProperties_ProcessPropertiesAndConstructors) ProcessGetProp(getProp rhino.Node) {
	a := p.outer
	type_ := a.GetColor(getProp.GetFirstChild())
	p.MaybeMarkCandidate(getProp, type_)
	if NodeUtil_IsLhsOfAssign(getProp) || NodeUtil_IsStatement(getProp.GetParent()) {
		a.graphNodeFactory.CreateNode(type_)
	}
}

// Java: AmbiguateProperties.ProcessPropertiesAndConstructors.processCall(Node)
func (p *AmbiguateProperties_ProcessPropertiesAndConstructors) ProcessCall(call rhino.Node) {
	a := p.outer
	target := call.GetFirstChild()
	if !target.IsQualifiedName() {
		return
	}

	if a.compiler.GetCodingConvention().IsPropertyRenameFunction(target) {
		propName := call.GetSecondChild()
		if propName.IsNil() || !propName.IsStringLit() {
			return
		}

		// Skip ambiguation for properties in renaming calls
		// NOTE (lharker@) - I'm not sure if this behavior is necessary, or if we could safely
		// ambiguate the property as long as we also updated the property renaming call
		prop := a.GetProperty(propName.GetString())
		prop.skipAmbiguating = true
	} else if NodeUtil_IsObjectDefinePropertiesDefinition(call) {
		typeObj := call.GetSecondChild()
		type_ := a.GetColor(typeObj)
		objectLiteral := typeObj.GetNext()

		if !objectLiteral.IsObjectLit() {
			return
		}

		for key := objectLiteral.GetFirstChild(); !key.IsNil(); key = key.GetNext() {
			p.ProcessObjectProperty(objectLiteral, key, type_)
		}
	}
}

// Java: AmbiguateProperties.ProcessPropertiesAndConstructors.processObjectProperty(Node, Node, Color)
func (p *AmbiguateProperties_ProcessPropertiesAndConstructors) ProcessObjectProperty(objectLit rhino.Node, key rhino.Node, type_ *colors.Color) {
	a := p.outer
	jutil.CheckArgumentMsg1(objectLit.IsObjectLit() || objectLit.IsObjectPattern(), "%s", objectLit)
	switch key.GetToken() {
	case rhino.Token_COMPUTED_PROP:
		if key.GetFirstChild().IsStringLit() {
			// If this quoted prop name is statically determinable, ensure we don't rename some
			// other property in a way that could conflict with it.
			//
			// This is largely because we store quoted member functions as computed properties and
			// want to be consistent with how other quoted properties invalidate property names.
			a.quotedNames.Add(key.GetFirstChild().GetString())
		}
	case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_STRING_KEY:
		if key.IsQuotedStringKey() {
			// If this quoted prop name is statically determinable, ensure we don't rename some
			// other property in a way that could conflict with it
			a.quotedNames.Add(key.GetString())
		} else {
			p.MaybeMarkCandidate(key, type_)
		}
	case rhino.Token_OBJECT_REST, rhino.Token_OBJECT_SPREAD:
		// Nothing to do.
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected child of " + objectLit.GetToken().String() + ": " + key.ToStringTree()))
	}
}

// Java: AmbiguateProperties.ProcessPropertiesAndConstructors.processObjectLitOrPattern(Node)
func (p *AmbiguateProperties_ProcessPropertiesAndConstructors) ProcessObjectLitOrPattern(objectLit rhino.Node) {
	a := p.outer
	// Object.defineProperties literals are handled at the CALL node, as we determine the type
	// differently than for regular object literals.
	if objectLit.GetParent().IsCall() &&
		NodeUtil_IsObjectDefinePropertiesDefinition(objectLit.GetParent()) {
		return
	}

	// The children of an OBJECTLIT node are keys, where the values
	// are the children of the keys.
	type_ := a.GetColor(objectLit)
	for key := objectLit.GetFirstChild(); !key.IsNil(); key = key.GetNext() {
		p.ProcessObjectProperty(objectLit, key, type_)
	}
}

// Java: AmbiguateProperties.ProcessPropertiesAndConstructors.processGetElem(Node)
func (p *AmbiguateProperties_ProcessPropertiesAndConstructors) ProcessGetElem(n rhino.Node) {
	// If this is a quoted property access (e.g. x['myprop']), we need to
	// ensure that we never rename some other property in a way that
	// could conflict with this quoted name.
	child := n.GetLastChild()
	if child.IsStringLit() {
		p.outer.quotedNames.Add(child.GetString())
	}
}

// Java: AmbiguateProperties.ProcessPropertiesAndConstructors.processClass(Node)
func (p *AmbiguateProperties_ProcessPropertiesAndConstructors) ProcessClass(classNode rhino.Node) {
	a := p.outer
	classConstructorType := a.GetColor(classNode)
	a.graphNodeFactory.CreateNode(classConstructorType)
	// In theory all CLASS colors should be a function with a known prototype, but in
	// practice typecasts mean that this is not always the case.

	possiblePrototypes := classConstructorType.GetPrototypes()
	var classPrototype *colors.Color
	if possiblePrototypes.IsEmpty() {
		classPrototype = colors.StandardColors_UNKNOWN
	} else {
		classPrototype = colors.Color_CreateUnion(possiblePrototypes)
	}
	for member := NodeUtil_GetClassMembers(classNode).GetFirstChild(); !member.IsNil(); member = member.GetNext() {
		if member.IsQuotedStringKey() {
			// ignore get 'foo'() {} and prevent property name collisions
			// Note that only getters/setters are represented as quoted strings, not 'foo'() {}
			// see https://github.com/google/closure-compiler/issues/3071
			a.quotedNames.Add(member.GetString())
			continue
		} else if member.IsComputedProp() || member.IsComputedFieldDef() {
			// ignore ['foo']() {}
			// for simple cases, we also prevent renaming collisions
			if member.GetFirstChild().IsStringLit() {
				a.quotedNames.Add(member.GetFirstChild().GetString())
			}
			continue
		} else if NodeUtil_IsEs6ConstructorMemberFunctionDef(member) {
			// don't rename `class C { constructor() {} }` !
			// This only applies for ES6 classes, not generic properties called 'constructor', which
			// is why it's handled in this method specifically.
			continue
		} else if member.IsBlock() {
			// ES2022 static initialization blocks don't have names so can't be renamed.
			// Example: `class C { static { alert('foo'); } }`
			continue
		}

		var memberOwnerColor *colors.Color
		if member.IsStaticMember() {
			memberOwnerColor = classConstructorType
		} else if member.IsMemberFieldDef() {
			possibleInstances := classConstructorType.GetInstanceColors()
			if possibleInstances.IsEmpty() {
				memberOwnerColor = colors.StandardColors_UNKNOWN
			} else {
				memberOwnerColor = colors.Color_CreateUnion(possibleInstances)
			}
		} else {
			jutil.CheckState(member.IsMemberFunctionDef() || member.IsGetterDef() || member.IsSetterDef())
			memberOwnerColor = classPrototype
		}
		// member could be a MEMBER_FUNCTION_DEF, MEMBER_FIELD_DEF, GETTER_DEF, or SETTER_DEF
		p.MaybeMarkCandidate(member, memberOwnerColor)
	}
}

// MaybeMarkCandidate: if a property node is eligible for renaming, stashes a
// reference to it and increments the property name's access count.
//
// n is the STRING node for a property.
//
// Java: AmbiguateProperties.ProcessPropertiesAndConstructors.maybeMarkCandidate(Node, Color)
func (p *AmbiguateProperties_ProcessPropertiesAndConstructors) MaybeMarkCandidate(n rhino.Node, type_ *colors.Color) {
	a := p.outer
	name := n.GetString()
	if !a.externedNames.Contains(name) {
		a.stringNodesToRename = appendDoubling(a.stringNodesToRename, n)
		p.RecordProperty(name, type_)
	}
}

// Java: AmbiguateProperties.ProcessPropertiesAndConstructors.recordProperty(String, Color)
func (p *AmbiguateProperties_ProcessPropertiesAndConstructors) RecordProperty(name string, color *colors.Color) *AmbiguateProperties_Property {
	prop := p.outer.GetProperty(name)
	prop.AddRelatedColor(color)
	return prop
}

// Java: AmbiguateProperties.getProperty(String)
func (a *AmbiguateProperties) GetProperty(name string) *AmbiguateProperties_Property {
	if prop := a.propertyMap.Get(name); prop != nil {
		return prop
	}
	prop := a.allocProperty(name)
	a.propertyMap.PutNew(name, prop)
	return prop
}

func (a *AmbiguateProperties) allocProperty(name string) *AmbiguateProperties_Property {
	return &AmbiguateProperties_Property{outer: a, oldName: name, singleColorIndex: -1}
}

// GetColor gets the Color from the Node argument or UNKNOWN if not present.
//
// Java: AmbiguateProperties.getColor(Node)
func (a *AmbiguateProperties) GetColor(n rhino.Node) *colors.Color {
	type_ := n.GetColor()
	if type_ == nil {
		// TODO(bradfordcsmith): This branch indicates a compiler bug. It should throw an exception.
		return colors.StandardColors_UNKNOWN
	} else {
		return type_
	}
}

// AmbiguateProperties_Property encapsulates the information needed for
// renaming a property.
//
// Java: private class AmbiguateProperties.Property
type AmbiguateProperties_Property struct {
	outer *AmbiguateProperties

	oldName          string
	newName          string
	numOccurrences   int
	skipAmbiguating  bool
	singleColorIndex int32
	// All colors upon which this property was directly accessed. For "a.b" this includes "a"'s type
	firstColorSeed     *ColorGraphNode
	relatedColorsSeeds *jutil.LinkedHashMap[*ColorGraphNode, int] // may be nil
	// includes relatedTypesSeeds + all subtypes of those seed colors. For example if this property
	// was accessed off of Iterable, then this bitset will include Array as well.
	relatedColors jutil.BitSet
}

// Java: AmbiguateProperties.Property(String)
func NewAmbiguateProperties_Property(outer *AmbiguateProperties, name string) *AmbiguateProperties_Property {
	return &AmbiguateProperties_Property{
		outer:            outer,
		oldName:          name,
		singleColorIndex: -1,
	}
}

// AddRelatedColor marks this color as related to this property.
//
// Java: AmbiguateProperties.Property.addRelatedColor(Color)
func (p *AmbiguateProperties_Property) AddRelatedColor(color *colors.Color) {
	if p.skipAmbiguating {
		return
	}

	p.numOccurrences++

	if color.IsInvalidating() || color.GetPropertiesKeepOriginalName() {
		p.skipAmbiguating = true
		p.firstColorSeed = nil
		p.relatedColorsSeeds = nil
		return
	}

	newColorGraphNode := p.outer.graphNodeFactory.CreateNode(color)
	if p.firstColorSeed == nil {
		p.firstColorSeed = newColorGraphNode
		return
	}
	if p.firstColorSeed == newColorGraphNode {
		return
	}
	if p.relatedColorsSeeds == nil {
		p.relatedColorsSeeds = jutil.NewLinkedHashMap[*ColorGraphNode, int]()
		p.relatedColorsSeeds.Put(p.firstColorSeed, 0)
	}
	p.relatedColorsSeeds.Put(newColorGraphNode, 0)
}
