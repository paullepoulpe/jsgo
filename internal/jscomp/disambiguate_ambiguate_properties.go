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
	propIndexByID map[uint32]uint32

	// Property names that don't get renamed
	externedIDs inlineBitSet

	// Names to which properties shouldn't be renamed, to avoid name conflicts
	quotedIDs inlineBitSet

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
	var externedIDs inlineBitSet
	externedIDs.Set(int(rhino.RhinoStringPool_InternID("prototype")))
	for name := range externProperties.All() {
		externedIDs.Set(int(rhino.RhinoStringPool_InternID(name)))
	}
	return &AmbiguateProperties{
		compiler:                   compiler,
		reservedFirstCharacters:    reservedFirstCharacters,
		reservedNonFirstCharacters: reservedNonFirstCharacters,
		properties:                 make([]AmbiguateProperties_Property, 0, 2048),
		propIndexByID:              make(map[uint32]uint32, 2048),
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
	// Populate nodesToRename, propertyMap, quotedNames.
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
		reservedNames[rhino.RhinoStringPool_StringFromID(uint32(id))] = struct{}{}
	}
	for id := a.quotedIDs.NextSetBit(0); id >= 0; id = a.quotedIDs.NextSetBit(id + 1) {
		reservedNames[rhino.RhinoStringPool_StringFromID(uint32(id))] = struct{}{}
	}
	a.externedIDs = inlineBitSet{}
	a.quotedIDs = inlineBitSet{}
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

		var relatedTypes inlineBitSet
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
		if p.newNameID != 0 && uint32(p.newNameID) != p.oldNameID && item.node.GetStringID() != uint32(p.newNameID) {
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

// AmbiguateProperties_ProcessPropertiesAndConstructors finds all property
// references, recording the types on which they occur, and records all
// constructors and their instance types in the ColorGraphNodeFactory.
//
// Java: private class AmbiguateProperties.ProcessPropertiesAndConstructors extends AbstractPostOrderCallback
type AmbiguateProperties_ProcessPropertiesAndConstructors struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *AmbiguateProperties
}

var ambiguatePropertiesVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_GETPROP)
	s.Add(rhino.Token_OPTCHAIN_GETPROP)
	s.Add(rhino.Token_CALL)
	s.Add(rhino.Token_NAME)
	s.Add(rhino.Token_OBJECTLIT)
	s.Add(rhino.Token_OBJECT_PATTERN)
	s.Add(rhino.Token_GETELEM)
	s.Add(rhino.Token_CLASS)
	return s
}()

var (
	_ nodeTraversal_onlyNameLeaves = (*AmbiguateProperties_ProcessPropertiesAndConstructors)(nil)
	_ nodeTraversal_visitTokens    = (*AmbiguateProperties_ProcessPropertiesAndConstructors)(nil)
)

func (*AmbiguateProperties_ProcessPropertiesAndConstructors) onlyNameLeaves() {}

func (*AmbiguateProperties_ProcessPropertiesAndConstructors) visitTokens() *ast.TokenSet {
	return &ambiguatePropertiesVisitTokens
}

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
		propIdx := a.getOrAddPropertyIndexByID(propName.GetStringID())
		a.properties[propIdx].skipAmbiguating = true
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
			a.quotedIDs.Set(int(key.GetFirstChild().GetStringID()))
		}
	case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_STRING_KEY:
		if key.IsQuotedStringKey() {
			// If this quoted prop name is statically determinable, ensure we don't rename some
			// other property in a way that could conflict with it
			a.quotedIDs.Set(int(key.GetStringID()))
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
		p.outer.quotedIDs.Set(int(child.GetStringID()))
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
	if len(possiblePrototypes) == 0 {
		classPrototype = colors.StandardColors_UNKNOWN
	} else {
		classPrototype = colors.Color_CreateUnionSlice(possiblePrototypes)
	}
	for member := NodeUtil_GetClassMembers(classNode).GetFirstChild(); !member.IsNil(); member = member.GetNext() {
		if member.IsQuotedStringKey() {
			// ignore get 'foo'() {} and prevent property name collisions
			// Note that only getters/setters are represented as quoted strings, not 'foo'() {}
			// see https://github.com/google/closure-compiler/issues/3071
			a.quotedIDs.Set(int(member.GetStringID()))
			continue
		} else if member.IsComputedProp() || member.IsComputedFieldDef() {
			// ignore ['foo']() {}
			// for simple cases, we also prevent renaming collisions
			if member.GetFirstChild().IsStringLit() {
				a.quotedIDs.Set(int(member.GetFirstChild().GetStringID()))
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
			if len(possibleInstances) == 0 {
				memberOwnerColor = colors.StandardColors_UNKNOWN
			} else {
				memberOwnerColor = colors.Color_CreateUnionSlice(possibleInstances)
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
	id := n.GetStringID()
	if !a.externedIDs.Get(int(id)) {
		propIdx := p.recordPropertyByID(id, type_)
		a.nodesToRename = append(a.nodesToRename, ambiguateNodeToRename{node: n, propIdx: propIdx})
	}
}

func (p *AmbiguateProperties_ProcessPropertiesAndConstructors) recordPropertyByID(id uint32, color *colors.Color) uint32 {
	a := p.outer
	propIdx := a.getOrAddPropertyIndexByID(id)
	a.properties[propIdx].addRelatedColor(a.graphNodeFactory, color)
	return propIdx
}

func (a *AmbiguateProperties) getOrAddPropertyIndexByID(id uint32) uint32 {
	if idx, ok := a.propIndexByID[id]; ok {
		return idx
	}
	idx := uint32(len(a.properties))
	a.properties = append(a.properties, AmbiguateProperties_Property{
		oldName:          rhino.RhinoStringPool_StringFromID(id),
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
	oldNameID        uint32
	newNameID        strpool.ID
	singleColorIndex int32
	skipAmbiguating  bool
	seedDedupBitSet  bool
	// All colors upon which this property was directly accessed. For "a.b" this includes "a"'s type.
	firstColorSeed *ColorGraphNode
	extraSeeds     []*ColorGraphNode
	// Includes seed colors + all subtypes of those seed colors. For example if this property
	// was accessed off of Iterable, then this bitset will include Array as well.
	relatedColors inlineBitSet
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
