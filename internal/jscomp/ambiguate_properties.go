// Port of com/google/javascript/jscomp/disambiguate/AmbiguateProperties.java @ closure-compiler v20261028

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
	"slices"
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/graph"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

type ambiguateNodeToRename struct {
	node    rhino.Node
	propIdx uint32
}

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

	nodesToRename []ambiguateNodeToRename
	// Can't use these to start property names.
	reservedFirstCharacters ASCIICharSet
	// Can't use these at all in property names.
	reservedNonFirstCharacters ASCIICharSet

	// Properties discovered in deterministic encounter order, indexed by name ID.
	properties    []AmbiguateProperties_Property
	propIndexByID map[strpool.ID]uint32

	// Property names that don't get renamed
	externedIDs smallBitSet

	// Names to which properties shouldn't be renamed, to avoid name conflicts
	quotedIDs smallBitSet

	colorRegistry *colors.ColorRegistry

	// Map from original property name to new name. Only used by tests.
	renamingMap *OrderedStringMap[string] // may be nil

	graphNodeFactory *ColorGraphNodeFactory // may be nil
}

var ambiguateProperties_logger = jutil.Logger_GetLogger("com.google.javascript.jscomp.disambiguate.AmbiguateProperties")

// ambiguateProperties_FREQUENCY_COMPARATOR sorts Property objects by their
// count, breaking ties alphabetically to ensure a deterministic total
// ordering.
//
// Java: AmbiguateProperties.FREQUENCY_COMPARATOR
func ambiguateProperties_FREQUENCY_COMPARATOR(p1 *AmbiguateProperties_Property, p2 *AmbiguateProperties_Property) int {
	if p1.numOccurrences != p2.numOccurrences {
		return p2.numOccurrences - p1.numOccurrences
	}
	return jutil.StringCompareTo(p1.oldName, p2.oldName)
}

// Java: AmbiguateProperties(AbstractCompiler, Set<Character>, Set<Character>, Set<String>)
func NewAmbiguateProperties(
	compiler AbstractCompiler,
	reservedFirstCharacters ASCIICharSet,
	reservedNonFirstCharacters ASCIICharSet,
	externProperties jutil.Set[string]) *AmbiguateProperties {
	jutil.CheckState(compiler.GetLifeCycleStage().IsNormalized())
	// ImmutableSet.<String>builder().add("prototype").addAll(externProperties).build()
	var externedIDs smallBitSet
	externedIDs.Set(int(idPrototype))
	for name := range externProperties.All() {
		externedIDs.Set(int(strpool.Intern(name)))
	}
	return &AmbiguateProperties{
		compiler:                   compiler,
		reservedFirstCharacters:    reservedFirstCharacters,
		reservedNonFirstCharacters: reservedNonFirstCharacters,
		properties:                 make([]AmbiguateProperties_Property, 0, 2048),
		propIndexByID:              make(map[strpool.ID]uint32, 2048),
		externedIDs:                externedIDs,
		colorRegistry:              compiler.GetColorRegistry(),
	}
}

// Java: AmbiguateProperties.makePassForTesting(AbstractCompiler, Set<Character>, Set<Character>, Set<String>)
func AmbiguateProperties_MakePassForTesting(
	compiler AbstractCompiler,
	reservedFirstCharacters ASCIICharSet,
	reservedNonFirstCharacters ASCIICharSet,
	externProperties jutil.Set[string]) *AmbiguateProperties {
	ap := NewAmbiguateProperties(
		compiler, reservedFirstCharacters, reservedNonFirstCharacters, externProperties)
	ap.renamingMap = NewOrderedStringMap[string]()
	return ap
}

// Java: AmbiguateProperties.getRenamingMap()
func (a *AmbiguateProperties) GetRenamingMap() *OrderedStringMap[string] {
	jutil.CheckNotNull(a.renamingMap)
	return a.renamingMap
}

// Java: AmbiguateProperties.process(Node, Node)
func (a *AmbiguateProperties) Process(externs rhino.Node, root rhino.Node) {
	a.graphNodeFactory = ColorGraphNodeFactory_CreateFactory(a.colorRegistry)

	// Find all property references and record the types on which they occur.
	// Populate nodesToRename, propertyMap, quotedNames.
	NodeTraversal_TraverseVisitor(a.compiler, root, a.visitor())

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
			 * elements. This prevents us from needing to check cardinality or equality, which
			 * would otherwise dominate the cost of computing the fixed point.
			 *
			 * <p>We're guaranteed to converge because the sizes will be equal after the OR
			 * operation.
			 */
			superSet := supertype.GetSubtypeIndices()
			if !subtype.hasSubtypeIndices {
				idx := subtype.GetIndex()
				if superSet.Get(idx) {
					return false
				}
				superSet.Set(idx)
				return true
			}
			subSet := &subtype.subtypeIndices
			if subSet.WordCapacityBits() > superSet.WordCapacityBits() {
				superSet.Or(subSet)
				return true
			}

			startSize := superSet.Cardinality()
			superSet.Or(subSet)
			return superSet.Cardinality() > startSize
		})).
		ComputeFixedPoint(colorGraph)

	// Fill in all transitive edges in subtyping graph per property.
	for i := range a.properties {
		prop := &a.properties[i]
		if prop.skipAmbiguating || prop.firstColorSeed == nil {
			prop.firstColorSeed = nil
			prop.extraSeeds = nil
			continue
		}
		first := prop.firstColorSeed
		prop.firstColorSeed = nil
		if len(prop.extraSeeds) == 0 {
			if !first.hasSubtypeIndices {
				prop.singleColorIndex = int32(first.index)
			} else {
				prop.relatedColors.Or(&first.subtypeIndices)
			}
			continue
		}
		if first.hasSubtypeIndices {
			prop.relatedColors.Or(&first.subtypeIndices)
		} else {
			prop.relatedColors.Set(first.index)
		}
		for _, color := range prop.extraSeeds {
			if color.hasSubtypeIndices {
				prop.relatedColors.Or(&color.subtypeIndices)
			} else {
				prop.relatedColors.Set(color.index)
			}
		}
		clear(prop.extraSeeds)
		prop.extraSeeds = nil
	}
	colorGraph = nil
	a.graphNodeFactory = nil
	a.propIndexByID = nil

	// ImmutableSet.<String>builder().addAll(externedNames).addAll(quotedNames)
	numProperties := len(a.properties)
	reservedNames := make(stringMapSet, a.externedIDs.Cardinality()+a.quotedIDs.Cardinality()+256)
	for id := a.externedIDs.NextSetBit(0); id >= 0; id = a.externedIDs.NextSetBit(id + 1) {
		reservedNames[strpool.String(strpool.ID(id))] = struct{}{}
	}
	for id := a.quotedIDs.NextSetBit(0); id >= 0; id = a.quotedIDs.NextSetBit(id + 1) {
		reservedNames[strpool.String(strpool.ID(id))] = struct{}{}
	}
	a.externedIDs = smallBitSet{}
	a.quotedIDs = smallBitSet{}
	numSkippedPropertyNames := 0
	props := make([]*AmbiguateProperties_Property, 0, numProperties)
	for i := range a.properties {
		prop := &a.properties[i]
		if prop.skipAmbiguating {
			numSkippedPropertyNames++
			reservedNames[prop.oldName] = struct{}{}
		} else {
			props = append(props, prop)
		}
	}
	finalNumRenamedPropertyNames := len(props)
	finalNumSkippedPropertyNames := numSkippedPropertyNames

	// Greedily color unrelated properties and assign generated names directly.
	nameGen := NewDefaultNameGeneratorWithFirstAndNonFirstChars(
		reservedNames,
		"",
		a.reservedFirstCharacters,
		a.reservedNonFirstCharacters)
	reservedNames = nil

	var numNewPropertyNames int
	if len(props) == 0 {
		// Preserves GreedyGraphColoring's do-while count on an empty graph.
		numNewPropertyNames = 1
	} else {
		var origProps []*AmbiguateProperties_Property
		if a.renamingMap != nil {
			origProps = slices.Clone(props)
		}
		slices.SortStableFunc(props, ambiguateProperties_FREQUENCY_COMPARATOR)

		var relatedTypes smallBitSet
		worklist := props
		for len(worklist) > 0 {
			relatedTypes.ClearAll()
			colorName := nameGen.GenerateNextName()
			colorID := strpool.Intern(colorName)
			numNewPropertyNames++
			kept := worklist[:0]
			for _, prop := range worklist {
				var independent bool
				if prop.singleColorIndex >= 0 {
					independent = !relatedTypes.Get(int(prop.singleColorIndex))
				} else {
					independent = !relatedTypes.Intersects(&prop.relatedColors)
				}
				if independent {
					if prop.singleColorIndex >= 0 {
						relatedTypes.Set(int(prop.singleColorIndex))
					} else {
						relatedTypes.Or(&prop.relatedColors)
					}
					prop.newName = colorName
					prop.newNameID = colorID
				} else {
					kept = append(kept, prop)
				}
			}
			clear(worklist[len(kept):])
			worklist = kept
		}
		if a.renamingMap != nil {
			for _, prop := range origProps {
				a.renamingMap.Put(prop.oldName, prop.newName)
			}
		}
	}
	nameGen = nil
	props = nil

	// Assign the new names to the relevant STRING nodes in the AST.
	for _, item := range a.nodesToRename {
		p := &a.properties[item.propIdx]
		if p.newNameID != 0 && p.newNameID != p.oldNameID && nodeStringID(item.node) != p.newNameID {
			setNodeStringID(item.node, p.newNameID)
			a.compiler.ReportChangeToEnclosingScope(item.node)
		}
	}
	a.nodesToRename = nil
	a.properties = nil

	// We may have renamed getter / setter properties.
	// TODO(b/161947315): this shouldn't be the responsibility of AmbiguateProperties
	GatherGetterAndSetterProperties_Update(a.compiler, externs, root)

	// Suppliers.memoize(...)
	var summary string
	var computed bool
	summarySupplier := func() string {
		if !computed {
			summary = "Collapsed " +
				strconv.Itoa(finalNumRenamedPropertyNames) +
				" properties into " +
				strconv.Itoa(numNewPropertyNames) +
				" and skipped renaming " +
				strconv.Itoa(finalNumSkippedPropertyNames) +
				" properties."
			computed = true
		}
		return summary
	}
	if ambiguateProperties_logger.IsLoggable(jutil.Level_FINE) {
		ambiguateProperties_logger.Fine(summarySupplier())
	}
	a.compiler.ReportAmbiguatePropertiesSummary(summarySupplier)
}

func (a *AmbiguateProperties) visitor() *ast.Visitor[*NodeTraversal] {
	return &ast.Visitor[*NodeTraversal]{
		GetProp: func(_ *NodeTraversal, gp ast.GetProp) {
			a.processGetProp(gp.Node(), gp.Object().Node())
		},
		OptChainGetProp: func(_ *NodeTraversal, gp ast.OptChainGetProp) {
			a.processGetProp(gp.Node(), gp.Object().Node())
		},
		Call: func(_ *NodeTraversal, call ast.Call) {
			a.processCall(call)
		},
		Name: func(_ *NodeTraversal, name ast.Name) {
			// handle ES5-style classes
			parent := name.Node().GetParent()
			if NodeUtil_IsNameDeclaration(parent) || parent.IsFunction() {
				a.graphNodeFactory.CreateNode(a.GetColor(name.Node()))
			}
		},
		ObjectLit: func(_ *NodeTraversal, lit ast.ObjectLit) {
			// Object.defineProperties literals are handled at the CALL node, as we determine the type
			// differently than for regular object literals.
			if !isInObjectDefinePropertiesCall(lit.Node()) {
				a.processObjectLit(lit, a.GetColor(lit.Node()))
			}
		},
		ObjectPattern: func(_ *NodeTraversal, pat ast.ObjectPattern) {
			type_ := a.GetColor(pat.Node())
			for prop := range pat.Props() {
				a.processObjectProperty(pat.Node(), prop.Node(), type_)
			}
		},
		GetElem: func(_ *NodeTraversal, ge ast.GetElem) {
			// If this is a quoted property access (e.g. x['myprop']), we need to
			// ensure that we never rename some other property in a way that
			// could conflict with this quoted name.
			if lit, ok := ge.Index().AsStringLit(); ok {
				a.quotedIDs.Set(int(lit.Node().GetStringID()))
			}
		},
		Class: func(_ *NodeTraversal, cls ast.Class) {
			a.processClass(cls)
		},
	}
}

// Java: AmbiguateProperties.ProcessPropertiesAndConstructors.processGetProp(Node)
func (a *AmbiguateProperties) processGetProp(getProp, obj rhino.Node) {
	type_ := a.GetColor(obj)
	a.maybeMarkCandidate(getProp, type_)
	if NodeUtil_IsLhsOfAssign(getProp) || NodeUtil_IsStatement(getProp.GetParent()) {
		a.graphNodeFactory.CreateNode(type_)
	}
}

// Java: AmbiguateProperties.ProcessPropertiesAndConstructors.processCall(Node)
func (a *AmbiguateProperties) processCall(call ast.Call) {
	target := call.Callee().Node()
	if !target.IsQualifiedName() {
		return
	}

	if a.compiler.GetCodingConvention().IsPropertyRenameFunction(target) {
		if firstArg, ok := call.FirstArg(); ok {
			if propName, ok := firstArg.AsStringLit(); ok {
				// Skip ambiguation for properties in renaming calls
				// NOTE (lharker@) - I'm not sure if this behavior is necessary, or if we could safely
				// ambiguate the property as long as we also updated the property renaming call
				propIdx := a.getOrAddPropertyIndexByID(nodeStringID(propName.Node()))
				a.properties[propIdx].skipAmbiguating = true
			}
		}
	} else if typeObj, objectLit, ok := unpackObjectDefineProperties(call); ok {
		a.processObjectLit(objectLit, a.GetColor(typeObj.Node()))
	}
}

func (a *AmbiguateProperties) processObjectLit(lit ast.ObjectLit, type_ *colors.Color) {
	for member := range lit.Members() {
		a.processObjectProperty(lit.Node(), member.Node(), type_)
	}
}

// Java: AmbiguateProperties.ProcessPropertiesAndConstructors.processObjectProperty(Node, Node, Color)
func (a *AmbiguateProperties) processObjectProperty(objectLit rhino.Node, key rhino.Node, type_ *colors.Color) {
	switch key.GetToken() {
	case rhino.Token_COMPUTED_PROP:
		if lit, ok := ast.MustComputedProp(key).Key().AsStringLit(); ok {
			// If this quoted prop name is statically determinable, ensure we don't rename some
			// other property in a way that could conflict with it.
			//
			// This is largely because we store quoted member functions as computed properties and
			// want to be consistent with how other quoted properties invalidate property names.
			a.quotedIDs.Set(int(lit.Node().GetStringID()))
		}
	case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_STRING_KEY:
		if key.IsQuotedStringKey() {
			// If this quoted prop name is statically determinable, ensure we don't rename some
			// other property in a way that could conflict with it
			a.quotedIDs.Set(int(key.GetStringID()))
		} else {
			a.maybeMarkCandidate(key, type_)
		}
	case rhino.Token_OBJECT_REST, rhino.Token_OBJECT_SPREAD:
		// Nothing to do.
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected child of " + objectLit.GetToken().String() + ": " + key.ToStringTree()))
	}
}

// Java: AmbiguateProperties.ProcessPropertiesAndConstructors.processClass(Node)
func (a *AmbiguateProperties) processClass(cls ast.Class) {
	classConstructorType := a.GetColor(cls.Node())
	a.graphNodeFactory.CreateNode(classConstructorType)
	// In theory all CLASS colors should be a function with a known prototype, but in
	// practice typecasts mean that this is not always the case.
	owners := classMemberOwners{ctor: classConstructorType}
	for member := range cls.Body().Members() {
		n := member.Node()
		if n.IsQuotedStringKey() {
			// ignore get 'foo'() {} and prevent property name collisions
			// Note that only getters/setters are represented as quoted strings, not 'foo'() {}
			// see https://github.com/google/closure-compiler/issues/3071
			a.quotedIDs.Set(int(n.GetStringID()))
			continue
		}
		if cp, ok := member.AsComputedProp(); ok {
			// ignore ['foo']() {}
			// for simple cases, we also prevent renaming collisions
			if lit, ok := cp.Key().AsStringLit(); ok {
				a.quotedIDs.Set(int(lit.Node().GetStringID()))
			}
			continue
		}
		if cf, ok := member.AsComputedFieldDef(); ok {
			if lit, ok := cf.Key().AsStringLit(); ok {
				a.quotedIDs.Set(int(lit.Node().GetStringID()))
			}
			continue
		}
		if mf, ok := member.AsMemberFunctionDef(); ok && mf.IsEs6Constructor() {
			// don't rename `class C { constructor() {} }` !
			// This only applies for ES6 classes, not generic properties called 'constructor', which
			// is why it's handled in this method specifically.
			continue
		}
		if _, ok := member.AsBlock(); ok {
			// ES2022 static initialization blocks don't have names so can't be renamed.
			// Example: `class C { static { alert('foo'); } }`
			continue
		}
		// member could be a MEMBER_FUNCTION_DEF, MEMBER_FIELD_DEF, GETTER_DEF, or SETTER_DEF
		a.maybeMarkCandidate(n, owners.ownerOf(member))
	}
}

// maybeMarkCandidate: if a property node is eligible for renaming, stashes a
// reference to it and increments the property name's access count.
//
// n is the STRING node for a property.
//
// Java: AmbiguateProperties.ProcessPropertiesAndConstructors.maybeMarkCandidate(Node, Color)
func (a *AmbiguateProperties) maybeMarkCandidate(n rhino.Node, type_ *colors.Color) {
	id := nodeStringID(n)
	if !a.externedIDs.Get(int(id)) {
		propIdx := a.recordPropertyByID(id, type_)
		a.nodesToRename = append(a.nodesToRename, ambiguateNodeToRename{node: n, propIdx: propIdx})
	}
}

func (a *AmbiguateProperties) recordPropertyByID(id strpool.ID, color *colors.Color) uint32 {
	propIdx := a.getOrAddPropertyIndexByID(id)
	a.properties[propIdx].addRelatedColor(a.graphNodeFactory, color)
	return propIdx
}

func (a *AmbiguateProperties) getOrAddPropertyIndexByID(id strpool.ID) uint32 {
	if idx, ok := a.propIndexByID[id]; ok {
		return idx
	}
	idx := uint32(len(a.properties))
	a.properties = append(a.properties, AmbiguateProperties_Property{
		oldName:          strpool.String(id),
		oldNameID:        id,
		singleColorIndex: -1,
	})
	a.propIndexByID[id] = idx
	return idx
}

// GetColor gets the Color from the Node argument or UNKNOWN if not present.
//
// Java: AmbiguateProperties.getColor(Node)
func (a *AmbiguateProperties) GetColor(n rhino.Node) *colors.Color {
	type_ := n.GetColor()
	if type_ == nil {
		// TODO(bradfordcsmith): This branch indicates a compiler bug. It should throw an exception.
		return colors.StandardColors_UNKNOWN
	}
	return type_
}

// AmbiguateProperties_Property encapsulates the information needed for
// renaming a property.
//
// Java: private class AmbiguateProperties.Property
type AmbiguateProperties_Property struct {
	oldName          string
	newName          string
	numOccurrences   int
	oldNameID        strpool.ID
	newNameID        strpool.ID
	singleColorIndex int32
	skipAmbiguating  bool
	seedDedupBitSet  bool
	// All colors upon which this property was directly accessed. For "a.b" this includes "a"'s type.
	firstColorSeed *ColorGraphNode
	extraSeeds     []*ColorGraphNode
	// Includes seed colors + all subtypes of those seed colors. For example if this property
	// was accessed off of Iterable, then this bitset will include Array as well.
	relatedColors smallBitSet
}

// addRelatedColor marks this color as related to this property.
//
// Java: AmbiguateProperties.Property.addRelatedColor(Color)
func (p *AmbiguateProperties_Property) addRelatedColor(factory *ColorGraphNodeFactory, color *colors.Color) {
	if p.skipAmbiguating {
		return
	}

	p.numOccurrences++

	if color.IsInvalidating() || color.GetPropertiesKeepOriginalName() {
		p.skipAmbiguating = true
		p.firstColorSeed = nil
		p.extraSeeds = nil
		p.seedDedupBitSet = false
		p.relatedColors.ClearAll()
		return
	}

	node := factory.CreateNode(color)
	if p.firstColorSeed == nil {
		p.firstColorSeed = node
		return
	}
	if p.firstColorSeed == node {
		return
	}
	if p.seedDedupBitSet {
		if !p.relatedColors.Get(node.index) {
			p.relatedColors.Set(node.index)
			p.extraSeeds = append(p.extraSeeds, node)
		}
		return
	}
	for _, s := range p.extraSeeds {
		if s == node {
			return
		}
	}
	if len(p.extraSeeds) == 8 {
		p.seedDedupBitSet = true
		p.relatedColors.Set(p.firstColorSeed.index)
		for _, s := range p.extraSeeds {
			p.relatedColors.Set(s.index)
		}
		p.relatedColors.Set(node.index)
	}
	p.extraSeeds = append(p.extraSeeds, node)
}
