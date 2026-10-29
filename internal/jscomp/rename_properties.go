// Port of com/google/javascript/jscomp/RenameProperties.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/RemovePropertyRenamingCalls.java @ closure-compiler v20261021

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

/*
 * Copyright 2023 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

var (
	renamePropertiesVisitTokens = func() ast.TokenSet {
		var s ast.TokenSet
		s.Add(rhino.Token_GETPROP)
		s.Add(rhino.Token_OPTCHAIN_GETPROP)
		s.Add(rhino.Token_MEMBER_FIELD_DEF)
		s.Add(rhino.Token_GETELEM)
		s.Add(rhino.Token_OPTCHAIN_GETELEM)
		s.Add(rhino.Token_CALL)
		s.Add(rhino.Token_MEMBER_FUNCTION_DEF)
		s.Add(rhino.Token_GETTER_DEF)
		s.Add(rhino.Token_SETTER_DEF)
		s.Add(rhino.Token_STRING_KEY)
		s.Add(rhino.Token_FUNCTION)
		return s
	}()
	renamePropertiesPrototypeID = strpool.Intern("prototype")
	renamePropertiesJscFnID     = strpool.Intern(NodeUtil_JSC_PROPERTY_NAME_FN)
)

type propRenameCandidate struct {
	node rhino.Node
	prop *RenameProperties_Property
}

type propCallSite struct {
	call   rhino.Node
	parent rhino.Node
}

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

	toRemove                   []rhino.Node
	candidates                 []propRenameCandidate
	callSites                  []propCallSite
	reservedFirstCharacters    ASCIICharSet
	reservedNonFirstCharacters ASCIICharSet

	// Map from interned property ID to Property object, plus insertion-ordered list.
	propertyByID map[strpool.ID]*RenameProperties_Property
	propertyList []*RenameProperties_Property
	propSlab     []RenameProperties_Property

	// Property names that don't get renamed
	externedNames map[strpool.ID]struct{}

	// Names to which properties shouldn't be renamed, to avoid name conflicts
	quotedNames map[strpool.ID]struct{}

	// Shared name generator
	nameGenerator NameGenerator

	// Filter function for property names. The function takes a node as input and returns true if the
	// node should be renamed.
	propertyRenameEligibilityFilter func(rhino.Node) bool
	allPropertiesEligible           bool
	// Property names that should not be renamed, based on the propertyNameFilter. Unlike
	// externedNames, this set is dynamically constructed during the pass.
	filteredOutNames map[strpool.ID]struct{}
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
	reservedFirstCharacters ASCIICharSet,
	reservedNonFirstCharacters ASCIICharSet,
	nameGenerator NameGenerator) *RenameProperties {
	r := newRenameProperties(
		compiler,
		generatePseudoNames,
		prevUsedPropertyMap,
		reservedFirstCharacters,
		reservedNonFirstCharacters,
		nameGenerator,
		nil)
	r.allPropertiesEligible = true
	return r
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
	reservedFirstCharacters ASCIICharSet,
	reservedNonFirstCharacters ASCIICharSet,
	nameGenerator NameGenerator,
	propertyRenameEligibilityFilter func(rhino.Node) bool) *RenameProperties {
	return newRenameProperties(
		compiler,
		generatePseudoNames,
		prevUsedPropertyMap,
		reservedFirstCharacters,
		reservedNonFirstCharacters,
		nameGenerator,
		propertyRenameEligibilityFilter)
}

func newRenameProperties(
	compiler AbstractCompiler,
	generatePseudoNames bool,
	prevUsedPropertyMap *VariableMap,
	reservedFirstCharacters ASCIICharSet,
	reservedNonFirstCharacters ASCIICharSet,
	nameGenerator NameGenerator,
	propertyRenameEligibilityFilter func(rhino.Node) bool) *RenameProperties {
	externProperties := compiler.GetExternProperties()
	if externProperties == nil {
		panic(jutil.NewNullPointerException()) // externedNames.addAll(null)
	}
	externedNames := make(map[strpool.ID]struct{}, externProperties.Size()+1)
	externedNames[renamePropertiesPrototypeID] = struct{}{}
	for name := range externProperties.All() {
		externedNames[strpool.Intern(name)] = struct{}{}
	}
	return &RenameProperties{
		compiler:                        compiler,
		generatePseudoNames:             generatePseudoNames,
		prevUsedPropertyMap:             prevUsedPropertyMap,
		reservedFirstCharacters:         reservedFirstCharacters,
		reservedNonFirstCharacters:      reservedNonFirstCharacters,
		propertyByID:                    make(map[strpool.ID]*RenameProperties_Property, 256),
		externedNames:                   externedNames,
		quotedNames:                     make(map[strpool.ID]struct{}, 64),
		nameGenerator:                   nameGenerator,
		propertyRenameEligibilityFilter: propertyRenameEligibilityFilter,
	}
}

// Java: RenameProperties.process(Node, Node)
func (r *RenameProperties) Process(externs rhino.Node, root rhino.Node) {
	jutil.CheckState(r.compiler.GetLifeCycleStage().IsNormalized())

	NodeTraversal_Traverse(r.compiler, root, &RenameProperties_ProcessProperties{outer: r})

	numProperties := len(r.propertyList)
	reservedNames := make(stringMapSet, len(r.externedNames)+len(r.quotedNames)+numProperties+64)
	// order-independent: populates unordered reservedNames membership set
	for id := range r.externedNames {
		reservedNames[strpool.String(id)] = struct{}{}
	}
	// order-independent: populates unordered reservedNames membership set
	for id := range r.quotedNames {
		reservedNames[strpool.String(id)] = struct{}{}
	}

	// Assign names, sorted by descending frequency to minimize code size.
	propsByFreq := make([]*RenameProperties_Property, numProperties)
	copy(propsByFreq, r.propertyList)
	slices.SortFunc(propsByFreq, renameProperties_FREQUENCY_COMPARATOR)

	// First, try and reuse as many property names from the previous compilation
	// as possible.
	if r.prevUsedPropertyMap != nil {
		r.reusePropertyNamesSlice(reservedNames, propsByFreq)
	}

	r.generateNamesSlice(propsByFreq, reservedNames)

	// Update the string nodes.
	for _, cand := range r.candidates {
		p := cand.prop
		if !p.externed && p.newNameID != 0 {
			setNodeStringID(cand.node, p.newNameID)
			if p.newNameID != p.oldNameID {
				r.compiler.ReportChangeToEnclosingScope(cand.node)
			}
		}
	}
	r.candidates = nil

	// Update the call nodes.
	for _, site := range r.callSites {
		firstArg := site.call.GetSecondChild()
		argStr := firstArg.GetString()
		var sb strings.Builder
		sb.Grow(len(argStr))
		for rest := argStr; ; {
			dot := strings.IndexByte(rest, '.')
			var oldName string
			if dot < 0 {
				oldName = rest
			} else {
				oldName = rest[:dot]
			}
			replacement := oldName
			if id := strpool.Lookup(oldName); id != 0 {
				if p := r.propertyByID[id]; p != nil && !p.externed && p.newNameID != 0 {
					replacement = p.newName
				}
			}
			if sb.Len() > 0 {
				sb.WriteByte('.')
			}
			sb.WriteString(replacement)
			if dot < 0 {
				break
			}
			rest = rest[dot+1:]
		}
		site.call.ReplaceWith(rhino.IR_String(sb.String()))
		r.compiler.ReportChangeToEnclosingScope(site.parent)
	}
	r.callSites = nil

	// Complete queued removals.
	for _, n := range r.toRemove {
		parent := n.GetParent()
		r.compiler.ReportChangeToEnclosingScope(n)
		n.Detach()
		NodeUtil_MarkFunctionsDeleted(n, r.compiler)
		if !parent.HasChildren() && !parent.IsScript() {
			parent.Detach()
		}
	}
	r.toRemove = nil

	r.compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED_OBFUSCATED)
	// This pass may rename getter or setter properties
	GatherGetterAndSetterProperties_Update(r.compiler, externs, root)
}

func (r *RenameProperties) reusePropertyNamesSlice(reservedNames stringMapSet, allProps []*RenameProperties_Property) {
	if r.generatePseudoNames {
		return
	}
	for _, prop := range allProps {
		prevName, ok := r.prevUsedPropertyMap.origToNew[prop.oldName]
		if ok {
			if _, reserved := reservedNames[prevName]; reserved {
				continue
			}

			prop.newName = prevName
			reservedNames[prevName] = struct{}{}
		}
	}
}

func (r *RenameProperties) generateNamesSlice(props []*RenameProperties_Property, reservedNames stringMapSet) {
	r.nameGenerator.ResetWithFirstAndNonFirstChars(
		reservedNames, "", r.reservedFirstCharacters, r.reservedNonFirstCharacters)
	for _, p := range props {
		if r.generatePseudoNames {
			p.newName = "$" + p.oldName + "$"
		} else if p.newName == "" {
			p.newName = r.nameGenerator.GenerateNextName()
		}
		reservedNames[p.newName] = struct{}{}
		if !p.externed {
			if _, filtered := r.filteredOutNames[p.oldNameID]; !filtered {
				p.newNameID = strpool.Intern(p.newName)
			}
		}
	}
}

// GetPropertyMap gets the property renaming map (the "answer key"): a mapping
// from original names to new names.
//
// Java: RenameProperties.getPropertyMap()
func (r *RenameProperties) GetPropertyMap() *VariableMap {
	entries := make([]variableMapEntry, 0, len(r.propertyList))
	for _, p := range r.propertyList {
		if p.newName != "" {
			entries = append(entries, variableMapEntry{orig: p.oldName, newName: p.newName})
		}
	}
	return newVariableMapFromEntries(entries)
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

func (*RenameProperties_ProcessProperties) visitTokens() *ast.TokenSet {
	return &renamePropertiesVisitTokens
}

// Java: RenameProperties.ProcessProperties.visit(NodeTraversal, Node, Node)
func (p *RenameProperties_ProcessProperties) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	r := p.outer
	switch n.GetToken() {
	case rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETPROP, rhino.Token_MEMBER_FIELD_DEF:
		p.MaybeMarkCandidate(n)
	case rhino.Token_GETELEM, rhino.Token_OPTCHAIN_GETELEM:
		// If this is a quoted property access (e.g. x['myprop']), we need to
		// ensure that we never rename some other property in a way that
		// could conflict with this quoted name.
		child := n.GetLastChild()
		if !child.IsNil() && child.IsStringLit() {
			r.quotedNames[strpool.ID(child.GetStringID())] = struct{}{}
		}
	case rhino.Token_CALL:
		// We replace property renaming function calls with a string
		// containing the renamed property.
		fnName := n.GetFirstChild()

		if r.compiler.GetCodingConvention().IsPropertyRenameFunction(fnName) {
			r.callSites = append(r.callSites, propCallSite{call: n, parent: parent})
			p.CountCallCandidates(t, n)
			if !r.allPropertiesEligible && !r.propertyRenameEligibilityFilter(n) {
				stringArgument := n.GetSecondChild()
				if r.filteredOutNames == nil {
					r.filteredOutNames = make(map[strpool.ID]struct{})
				}
				argStr := stringArgument.GetString()
				for rest := argStr; ; {
					dot := strings.IndexByte(rest, '.')
					var name string
					if dot < 0 {
						name = rest
					} else {
						name = rest[:dot]
					}
					r.filteredOutNames[strpool.Intern(name)] = struct{}{}
					if dot < 0 {
						break
					}
					rest = rest[dot+1:]
				}
			}
		}
	case rhino.Token_MEMBER_FUNCTION_DEF:
		jutil.CheckState(!n.IsQuotedStringKey())
		if NodeUtil_IsEs6ConstructorMemberFunctionDef(n) {
			r.externedNames[strpool.ID(n.GetStringID())] = struct{}{}
		} else {
			p.MaybeMarkCandidate(n)
		}
	case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_STRING_KEY:
		if n.IsQuotedStringKey() {
			// Ensure that we never rename some other property in a way
			// that could conflict with this quoted key.
			r.quotedNames[strpool.ID(n.GetStringID())] = struct{}{}
		} else {
			p.MaybeMarkCandidate(n)
		}
	case rhino.Token_FUNCTION:
		// We eliminate any stub implementations of JSCompiler_renameProperty
		// that we encounter.
		if NodeUtil_IsFunctionDeclaration(n) {
			fn := ast.MustFunction(n)
			if strpool.ID(fn.NameNode().Node().GetStringID()) == renamePropertiesJscFnID {
				r.toRemove = append(r.toRemove, n)
			}
		} else if parent.IsName() && strpool.ID(parent.GetStringID()) == renamePropertiesJscFnID {
			if decl, ok := ast.AsNameDeclaration(parent.GetParent()); ok && decl.Token() == rhino.Token_VAR {
				r.toRemove = append(r.toRemove, parent)
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
				r.toRemove = append(r.toRemove, exprResult)
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
func (p *RenameProperties_ProcessProperties) MaybeMarkCandidate(n rhino.Node) {
	r := p.outer
	id := strpool.ID(n.GetStringID())
	if _, externed := r.externedNames[id]; !externed {
		if !r.allPropertiesEligible && !r.propertyRenameEligibilityFilter(n) {
			if r.filteredOutNames == nil {
				r.filteredOutNames = make(map[strpool.ID]struct{})
			}
			r.filteredOutNames[id] = struct{}{}
		}
		prop := p.countPropertyOccurrenceID(id)
		r.candidates = append(r.candidates, propRenameCandidate{node: n, prop: prop})
	}
}

// CountCallCandidates counts references to property names that occur in a
// special function call.
//
// callNode is the CALL node for a property; t is the traversal.
//
// Java: RenameProperties.ProcessProperties.countCallCandidates(NodeTraversal, Node)
func (p *RenameProperties_ProcessProperties) CountCallCandidates(t *NodeTraversal, callNode rhino.Node) {
	r := p.outer
	fnNameNode := callNode.GetFirstChild()
	firstArg := callNode.GetSecondChild()
	if !firstArg.IsStringLit() {
		t.Report(callNode, RenameProperties_BAD_CALL, fnNameNode.GetString())
		return
	}

	argStr := firstArg.GetString()
	for rest := argStr; ; {
		dot := strings.IndexByte(rest, '.')
		var name string
		if dot < 0 {
			name = rest
		} else {
			name = rest[:dot]
		}
		if !rhino.TokenStream_IsJSIdentifier(name) {
			t.Report(callNode, RenameProperties_BAD_ARG, fnNameNode.GetString())
		} else {
			id := strpool.Intern(name)
			if _, externed := r.externedNames[id]; !externed {
				p.countPropertyOccurrenceID(id)
			}
		}
		if dot < 0 {
			break
		}
		rest = rest[dot+1:]
	}
}

// CountPropertyOccurrence increments the occurrence count for a property
// name.
//
// Java: RenameProperties.ProcessProperties.countPropertyOccurrence(String)
func (p *RenameProperties_ProcessProperties) CountPropertyOccurrence(name string) {
	p.countPropertyOccurrenceID(strpool.Intern(name))
}

func (p *RenameProperties_ProcessProperties) countPropertyOccurrenceID(id strpool.ID) *RenameProperties_Property {
	r := p.outer
	prop := r.propertyByID[id]
	if prop == nil {
		prop = r.allocProperty(id, strpool.String(id))
		r.propertyByID[id] = prop
		r.propertyList = append(r.propertyList, prop)
	}
	prop.numOccurrences++
	return prop
}

func (r *RenameProperties) allocProperty(id strpool.ID, name string) *RenameProperties_Property {
	if len(r.propSlab) == 0 {
		r.propSlab = make([]RenameProperties_Property, 128)
	}
	prop := &r.propSlab[0]
	r.propSlab = r.propSlab[1:]
	prop.oldNameID = id
	prop.oldName = name
	return prop
}

// -------------------------------------------------------------------------

// RenameProperties_Property encapsulates the information needed for renaming
// a property.
//
// Java: private static class RenameProperties.Property
type RenameProperties_Property struct {
	oldName        string
	newName        string
	oldNameID      strpool.ID
	newNameID      strpool.ID
	numOccurrences int
	externed       bool
}

// Java: RenameProperties.Property(String)
func NewRenameProperties_Property(name string) *RenameProperties_Property {
	return &RenameProperties_Property{oldName: name, oldNameID: strpool.Intern(name)}
}

// -------------------------------------------------------------------------

// RemovePropertyRenamingCalls replaces all calls to property-renaming
// functions like JSCompiler_renameProperty with the string literal property
// name.
//
// Intended for use when RenameProperties isn't running, since this
// functionality is also baked into RenameProperties.
type RemovePropertyRenamingCalls struct {
	compiler AbstractCompiler
}

// Java: RemovePropertyRenamingCalls(AbstractCompiler)
func NewRemovePropertyRenamingCalls(compiler AbstractCompiler) *RemovePropertyRenamingCalls {
	return &RemovePropertyRenamingCalls{compiler: compiler}
}

// Java: RemovePropertyRenamingCalls.process(Node, Node)
func (r *RemovePropertyRenamingCalls) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(r.compiler, root, &RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback{outer: r})
}

// RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback is the
// inner class RemoveCompilerPropertyRenamingCallback.
type RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *RemovePropertyRenamingCalls
}

var removePropRenamingCallsVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_CALL)
	return s
}()

var (
	_ nodeTraversal_noVisitLeaves = (*RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback)(nil)
	_ nodeTraversal_visitTokens   = (*RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback)(nil)
)

func (*RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback) noVisitLeaves() {}

func (*RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback) visitTokens() *ast.TokenSet {
	return &removePropRenamingCallsVisitTokens
}

// Java: RemovePropertyRenamingCalls.RemoveCompilerPropertyRenamingCallback.visit(NodeTraversal, Node, Node)
func (c *RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if call, ok := ast.AsCall(n); ok {
		if c.outer.compiler.GetCodingConvention().IsPropertyRenameFunction(call.Callee().Node()) {
			c.ReplacePropertyRenamingCall(n)
		}
	}
}

// ReplacePropertyRenamingCall replaces JSCompiler_renameProperty("bar", obj)
// with "bar".
//
// Java: RemovePropertyRenamingCalls.RemoveCompilerPropertyRenamingCallback.replacePropertyRenamingCall(Node)
func (c *RemovePropertyRenamingCalls_RemoveCompilerPropertyRenamingCallback) ReplacePropertyRenamingCall(callNode rhino.Node) {
	propName := NodeUtil_GetArgumentForCallOrNew(callNode, 0)
	if !propName.IsNil() {
		callNode.ReplaceWith(propName.Detach())
		c.outer.compiler.ReportChangeToEnclosingScope(propName)
	}
}
