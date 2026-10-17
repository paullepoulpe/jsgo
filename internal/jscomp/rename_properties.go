// Port of com/google/javascript/jscomp/RenameProperties.java @ closure-compiler v20261015

/*
 * Copyright 2004 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RenameProperties renames properties (including methods) of all JavaScript
// objects. This includes prototypes, functions, object literals, etc.
//
// If provided a VariableMap of previously used names, it tries to reuse those
// names.
//
// To prevent a property from getting renamed you may extern it (add it to
// your externs file) or put it in quotes.
//
// To avoid run-time JavaScript errors, use quotes when accessing properties
// that are defined using quotes.
//
//	var a = {'myprop': 0}, b = a['myprop'];  // correct
//	var x = {'myprop': 0}, y = x.myprop;     // incorrect
//
// This pass also recognizes and replaces special renaming functions. They
// supply a property name as the string literal for the first argument. This
// pass will replace them as though they were JS property references. Here are
// two examples: JSCompiler_renameProperty('propertyName') -> 'jYq'
// JSCompiler_renameProperty('myProp.nestedProp.innerProp') -> 'e4.sW.C$'
//
// This class is not thread-safe.
//
// Java: class RenameProperties implements CompilerPass
type RenameProperties struct {
	compiler            AbstractCompiler
	generatePseudoNames bool

	// Property renaming map from a previous compilation.
	prevUsedPropertyMap *VariableMap

	toRemove                   *jutil.ArrayList[*rhino.Node]
	stringNodesToRename        *jutil.ArrayList[*rhino.Node]
	callNodeToParentMap        *jutil.LinkedHashMap[*rhino.Node, *rhino.Node]
	reservedFirstCharacters    jutil.Set[uint16]
	reservedNonFirstCharacters jutil.Set[uint16]

	// Map from property name to Property object
	propertyMap *jutil.LinkedHashMap[string, *RenameProperties_Property]

	// Property names that don't get renamed
	externedNames *jutil.LinkedHashSet[string]

	// Names to which properties shouldn't be renamed, to avoid name conflicts
	quotedNames *jutil.LinkedHashSet[string]

	// Shared name generator
	nameGenerator NameGenerator

	// Filter function for property names. The function takes a node as input and returns true if the
	// node should be renamed.
	propertyRenameEligibilityFilter func(*rhino.Node) bool
	// Property names that should not be renamed, based on the propertyNameFilter. Unlike
	// externedNames, this set is dynamically constructed during the pass.
	filteredOutNames *jutil.LinkedHashSet[string]
}

// renameProperties_DOT_SPLITTER is Splitter.on('.').split: every piece,
// empty ones included.
func renameProperties_DOT_SPLITTER(s string) []string {
	return strings.Split(s, ".")
}

// Java: RenameProperties.FREQUENCY_COMPARATOR
var renameProperties_FREQUENCY_COMPARATOR = func(p1 *RenameProperties_Property, p2 *RenameProperties_Property) int {
	/* First a frequently used names would always be picked first. */
	if p1.numOccurrences != p2.numOccurrences {
		return p2.numOccurrences - p1.numOccurrences
	}

	/* Finally, for determinism, we compare them based on the old name. */
	return jutil.StringCompareTo(p1.oldName, p2.oldName)
}

// Java: RenameProperties.BAD_CALL
var RenameProperties_BAD_CALL = DiagnosticType_Error(
	"JSC_BAD_RENAME_PROPERTY_FUNCTION_NAME_CALL",
	"Bad {0} call - the first argument must be a string literal")

// Java: RenameProperties.BAD_ARG
var RenameProperties_BAD_ARG = DiagnosticType_Error(
	"JSC_BAD_RENAME_PROPERTY_FUNCTION_NAME_ARG",
	"Bad {0} argument - ''{1}'' is not a valid JavaScript identifier")

// NewRenameProperties creates an instance.
//
// compiler is the JSCompiler. generatePseudoNames generates pseudo names, e.g
// foo -> $foo$ instead of compact obfuscated names; this is used for
// debugging. prevUsedPropertyMap is the property renaming map used in a
// previous compilation. If specified, reservedFirstCharacters won't be used in
// generated names for the first character, and reservedNonFirstCharacters
// won't be used in generated names for characters after the first.
// nameGenerator is a shared NameGenerator that this instance can use; the
// instance may reset or reconfigure it, so the caller should not expect any
// state to be preserved.
//
// Java: RenameProperties(AbstractCompiler, boolean, VariableMap, Set<Character>, Set<Character>, NameGenerator)
func NewRenameProperties(
	compiler AbstractCompiler,
	generatePseudoNames bool,
	prevUsedPropertyMap *VariableMap,
	reservedFirstCharacters jutil.Set[uint16],
	reservedNonFirstCharacters jutil.Set[uint16],
	nameGenerator NameGenerator) *RenameProperties {
	return NewRenamePropertiesWithFilter(
		compiler,
		generatePseudoNames,
		prevUsedPropertyMap,
		reservedFirstCharacters,
		reservedNonFirstCharacters,
		nameGenerator,
		func(n *rhino.Node) bool { return true })
}

// NewRenamePropertiesWithFilter creates an instance.
//
// The parameters are those of NewRenameProperties, plus
// propertyRenameEligibilityFilter, which limits the nodes that are renamed
// using a predicate function. When this returns true, the property name is
// renamed if and only if all other property references for the same name
// match the predicate. If this returns false for any instance of a property,
// then no property references anywhere in the AST to that property are
// renamed.
//
// Java: RenameProperties(AbstractCompiler, boolean, VariableMap, Set<Character>, Set<Character>, NameGenerator, Predicate<Node>)
func NewRenamePropertiesWithFilter(
	compiler AbstractCompiler,
	generatePseudoNames bool,
	prevUsedPropertyMap *VariableMap,
	reservedFirstCharacters jutil.Set[uint16],
	reservedNonFirstCharacters jutil.Set[uint16],
	nameGenerator NameGenerator,
	propertyRenameEligibilityFilter func(*rhino.Node) bool) *RenameProperties {
	r := &RenameProperties{
		toRemove:            jutil.NewArrayList[*rhino.Node](),
		stringNodesToRename: jutil.NewArrayList[*rhino.Node](),
		callNodeToParentMap: jutil.NewLinkedHashMap[*rhino.Node, *rhino.Node](),
		propertyMap:         jutil.NewLinkedHashMap[string, *RenameProperties_Property](),
		externedNames:       jutil.NewLinkedHashSet[string]("prototype"),
		quotedNames:         jutil.NewLinkedHashSet[string](),
		filteredOutNames:    jutil.NewLinkedHashSet[string](),
	}
	r.compiler = compiler
	r.generatePseudoNames = generatePseudoNames
	r.prevUsedPropertyMap = prevUsedPropertyMap
	r.reservedFirstCharacters = reservedFirstCharacters
	r.reservedNonFirstCharacters = reservedNonFirstCharacters
	r.nameGenerator = nameGenerator
	r.propertyRenameEligibilityFilter = propertyRenameEligibilityFilter
	externProperties := compiler.GetExternProperties()
	if externProperties == nil {
		panic(jutil.NewNullPointerException()) // externedNames.addAll(null)
	}
	for name := range externProperties.All() {
		r.externedNames.Add(name)
	}
	return r
}

// Java: RenameProperties.process(Node, Node)
func (r *RenameProperties) Process(externs *rhino.Node, root *rhino.Node) {
	jutil.CheckState(r.compiler.GetLifeCycleStage().IsNormalized())

	NodeTraversal_Traverse(r.compiler, root, &RenameProperties_ProcessProperties{outer: r})

	// Sets.newLinkedHashSetWithExpectedSize(externedNames.size() + quotedNames.size())
	reservedNames := jutil.NewLinkedHashSet[string]()
	for name := range r.externedNames.All() {
		reservedNames.Add(name)
	}
	for name := range r.quotedNames.All() {
		reservedNames.Add(name)
	}

	// Assign names, sorted by descending frequency to minimize code size.
	propsByFreq := jutil.NewTreeSet[*RenameProperties_Property](renameProperties_FREQUENCY_COMPARATOR)
	propsByFreq.AddAll(r.propertyMap.Values())

	// First, try and reuse as many property names from the previous compilation
	// as possible.
	if r.prevUsedPropertyMap != nil {
		r.ReusePropertyNames(reservedNames, propsByFreq)
	}

	r.GenerateNames(propsByFreq, reservedNames)

	// Update the string nodes.
	for n := range r.stringNodesToRename.All() {
		oldName := n.GetString()
		if r.filteredOutNames.Contains(oldName) {
			continue
		}
		p := r.propertyMap.Get(oldName)
		if p != nil && p.newName != nil {
			jutil.CheckState(oldName == p.oldName)
			n.SetString(*p.newName)
			if *p.newName != oldName {
				r.compiler.ReportChangeToEnclosingScope(n)
			}
		}
	}

	// Update the call nodes.
	for key, parent := range r.callNodeToParentMap.All() {
		firstArg := key.GetSecondChild()
		var sb strings.Builder
		for _, oldName := range renameProperties_DOT_SPLITTER(firstArg.GetString()) {
			p := r.propertyMap.Get(oldName)
			var replacement string
			if p != nil && p.newName != nil {
				jutil.CheckState(oldName == p.oldName)
				if r.filteredOutNames.Contains(oldName) {
					replacement = oldName
				} else {
					replacement = *p.newName
				}
			} else {
				replacement = oldName
			}
			if sb.Len() > 0 {
				sb.WriteByte('.')
			}
			sb.WriteString(replacement)
		}
		key.ReplaceWith(rhino.IR_String(sb.String()))
		r.compiler.ReportChangeToEnclosingScope(parent)
	}

	// Complete queued removals.
	for n := range r.toRemove.All() {
		parent := n.GetParent()
		r.compiler.ReportChangeToEnclosingScope(n)
		n.Detach()
		NodeUtil_MarkFunctionsDeleted(n, r.compiler)
		if !parent.HasChildren() && !parent.IsScript() {
			parent.Detach()
		}
	}

	r.compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED_OBFUSCATED)
	// This pass may rename getter or setter properties
	GatherGetterAndSetterProperties_Update(r.compiler, externs, root)
}

// ReusePropertyNames runs through the list of properties and renames as many
// as possible with names from the previous compilation. Also, updates
// reservedNames with the set of reused names.
//
// Java: RenameProperties.reusePropertyNames(Set<String>, Collection<RenameProperties.Property>)
func (r *RenameProperties) ReusePropertyNames(reservedNames *jutil.LinkedHashSet[string], allProps *jutil.TreeSet[*RenameProperties_Property]) {
	for prop := range allProps.All() {
		// Check if this node can reuse a name from a previous compilation - if
		// it can set the newName for the property too.
		prevName := r.prevUsedPropertyMap.LookupNewName(prop.oldName)
		if !r.generatePseudoNames && prevName != nil {
			// We can reuse prevName if it's not reserved.
			if reservedNames.Contains(*prevName) {
				continue
			}

			prop.newName = prevName
			reservedNames.Add(*prevName)
		}
	}
}

// GenerateNames generates new names for properties.
//
// props are the properties to generate new names for; reservedNames is a set
// of names to which properties should not be renamed.
//
// Java: RenameProperties.generateNames(Set<RenameProperties.Property>, Set<String>)
func (r *RenameProperties) GenerateNames(props *jutil.TreeSet[*RenameProperties_Property], reservedNames *jutil.LinkedHashSet[string]) {
	r.nameGenerator.ResetWithFirstAndNonFirstChars(
		reservedNames, "", r.reservedFirstCharacters, r.reservedNonFirstCharacters)
	for p := range props.All() {
		if r.generatePseudoNames {
			newName := "$" + p.oldName + "$"
			p.newName = &newName
		} else {
			// If we haven't already given this property a reusable name.
			if p.newName == nil {
				newName := r.nameGenerator.GenerateNextName()
				p.newName = &newName
			}
		}
		reservedNames.Add(*p.newName)
	}
}

// GetPropertyMap gets the property renaming map (the "answer key"): a mapping
// from original names to new names.
//
// Java: RenameProperties.getPropertyMap()
func (r *RenameProperties) GetPropertyMap() *VariableMap {
	map_ := jutil.NewImmutableMapBuilder[string, string]()
	for p := range r.propertyMap.Values() {
		if p.newName != nil {
			map_.Put(p.oldName, *p.newName)
		}
	}
	return NewVariableMap(map_.BuildOrThrow())
}

// -------------------------------------------------------------------------

// RenameProperties_ProcessProperties is a traversal callback that collects
// property names and counts how frequently each property name occurs.
//
// Java: private class RenameProperties.ProcessProperties extends AbstractPostOrderCallback
type RenameProperties_ProcessProperties struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *RenameProperties
}

// Java: RenameProperties.ProcessProperties.visit(NodeTraversal, Node, Node)
func (p *RenameProperties_ProcessProperties) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	r := p.outer
	switch n.GetToken() {
	case rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP, rhino.Token_MEMBER_FIELD_DEF:
		p.MaybeMarkCandidate(n)
	case rhino.Token_GETELEM, rhino.Token_OPTCHAIN_GETELEM:
		// If this is a quoted property access (e.g. x['myprop']), we need to
		// ensure that we never rename some other property in a way that
		// could conflict with this quoted name.
		child := n.GetLastChild()
		if child != nil && child.IsStringLit() {
			r.quotedNames.Add(child.GetString())
		}
	case rhino.Token_CALL:
		// We replace property renaming function calls with a string
		// containing the renamed property.
		fnName := n.GetFirstChild()

		if r.compiler.GetCodingConvention().IsPropertyRenameFunction(fnName) {
			r.callNodeToParentMap.Put(n, parent)
			p.CountCallCandidates(t, n)
			stringArgument := n.GetSecondChild()
			if !r.propertyRenameEligibilityFilter(n) {
				for _, name := range renameProperties_DOT_SPLITTER(stringArgument.GetString()) {
					r.filteredOutNames.Add(name)
				}
			}
		}
	case rhino.Token_MEMBER_FUNCTION_DEF:
		jutil.CheckState(!n.IsQuotedStringKey())
		if NodeUtil_IsEs6ConstructorMemberFunctionDef(n) {
			r.externedNames.Add(n.GetString())
		} else {
			p.MaybeMarkCandidate(n)
		}
	case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_STRING_KEY:
		if n.IsQuotedStringKey() {
			// Ensure that we never rename some other property in a way
			// that could conflict with this quoted key.
			r.quotedNames.Add(n.GetString())
		} else {
			p.MaybeMarkCandidate(n)
		}
	case rhino.Token_FUNCTION:
		// We eliminate any stub implementations of JSCompiler_renameProperty
		// that we encounter.
		if NodeUtil_IsFunctionDeclaration(n) {
			name := n.GetFirstChild().GetString()
			if NodeUtil_JSC_PROPERTY_NAME_FN == name {
				r.toRemove.Add(n)
			}
		} else if parent.IsName() && NodeUtil_JSC_PROPERTY_NAME_FN == parent.GetString() {
			varNode := parent.GetParent()
			if varNode.IsVar() {
				r.toRemove.Add(parent)
			}
		} else if NodeUtil_IsFunctionExpression(n) &&
			parent.IsAssign() &&
			parent.
				GetFirstChild().
				IsGetProp() && // JSCompiler does not handle optional calls to property rename
			// function
			r.compiler.GetCodingConvention().IsPropertyRenameFunction(parent.GetFirstChild()) {
			exprResult := parent.GetParent()
			if exprResult.IsExprResult() &&
				NodeUtil_IsStatementBlock(exprResult.GetParent()) &&
				exprResult.GetFirstChild().IsAssign() {
				r.toRemove.Add(exprResult)
			}
		}
	default:
	}
}

// MaybeMarkCandidate: if a property node is eligible for renaming, stashes a
// reference to it and increments the property name's access count.
//
// n is the STRING node for a property.
//
// Java: RenameProperties.ProcessProperties.maybeMarkCandidate(Node)
func (p *RenameProperties_ProcessProperties) MaybeMarkCandidate(n *rhino.Node) {
	r := p.outer
	name := n.GetString()
	if !r.externedNames.Contains(name) {
		if !r.propertyRenameEligibilityFilter(n) {
			r.filteredOutNames.Add(n.GetString())
		}
		r.stringNodesToRename.Add(n)
		p.CountPropertyOccurrence(name)
	}
}

// CountCallCandidates counts references to property names that occur in a
// special function call.
//
// callNode is the CALL node for a property; t is the traversal.
//
// Java: RenameProperties.ProcessProperties.countCallCandidates(NodeTraversal, Node)
func (p *RenameProperties_ProcessProperties) CountCallCandidates(t *NodeTraversal, callNode *rhino.Node) {
	r := p.outer
	fnName := callNode.GetFirstChild().GetString()
	firstArg := callNode.GetSecondChild()
	if !firstArg.IsStringLit() {
		t.Report(callNode, RenameProperties_BAD_CALL, fnName)
		return
	}

	for _, name := range renameProperties_DOT_SPLITTER(firstArg.GetString()) {
		if !rhino.TokenStream_IsJSIdentifier(name) {
			t.Report(callNode, RenameProperties_BAD_ARG, fnName)
			continue
		}
		if !r.externedNames.Contains(name) {
			p.CountPropertyOccurrence(name)
		}
	}
}

// CountPropertyOccurrence increments the occurrence count for a property
// name.
//
// Java: RenameProperties.ProcessProperties.countPropertyOccurrence(String)
func (p *RenameProperties_ProcessProperties) CountPropertyOccurrence(name string) {
	prop := p.outer.propertyMap.ComputeIfAbsent(name, NewRenameProperties_Property)
	prop.numOccurrences++
}

// -------------------------------------------------------------------------

// RenameProperties_Property encapsulates the information needed for renaming
// a property.
//
// Java: private static class RenameProperties.Property
type RenameProperties_Property struct {
	oldName        string
	newName        *string // may be nil
	numOccurrences int
}

// Java: RenameProperties.Property(String)
func NewRenameProperties_Property(name string) *RenameProperties_Property {
	return &RenameProperties_Property{oldName: name}
}
