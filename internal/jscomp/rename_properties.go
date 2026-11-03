// Port of com/google/javascript/jscomp/RenameProperties.java @ closure-compiler v20261028
// Port of com/google/javascript/jscomp/RemovePropertyRenamingCalls.java @ closure-compiler v20261028

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

type propRenameCandidate struct {
	node rhino.Node
	prop *RenameProperties_Property
}

type propCallSite struct {
	call   ast.Call
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
	propSlab     jutil.Slab[RenameProperties_Property]

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
	externProperties := compiler.GetExternPropSet()
	if externProperties == nil {
		panic(jutil.NewNullPointerException()) // externedNames.addAll(null)
	}
	externedNames := make(map[strpool.ID]struct{}, externProperties.Len()+1)
	externedNames[idPrototype] = struct{}{}
	for _, id := range externProperties.IDs() {
		externedNames[id] = struct{}{}
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

	NodeTraversal_TraverseVisitor(r.compiler, root, r.visitor())

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
		if p.newNameID != 0 && p.newNameID != p.oldNameID && nodeStringID(cand.node) != p.newNameID {
			setNodeStringID(cand.node, p.newNameID)
			r.compiler.ReportChangeToEnclosingScope(cand.node)
		}
	}
	r.candidates = nil

	// Update the call nodes.
	for _, site := range r.callSites {
		firstArg, _ := site.call.FirstArg()
		argStr := firstArg.Node().GetString()
		var sb strings.Builder
		sb.Grow(len(argStr))
		for oldName := range strings.SplitSeq(argStr, ".") {
			replacement := oldName
			if id := strpool.Lookup(oldName); id != 0 {
				if p := r.propertyByID[id]; p != nil && p.newNameID != 0 {
					replacement = p.newName
				}
			}
			if sb.Len() > 0 {
				sb.WriteByte('.')
			}
			sb.WriteString(replacement)
		}
		site.call.Node().ReplaceWith(rhino.IR_String(sb.String()))
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
	if fastFault(r.compiler.GetOptions(), faultBenignRename) {
		reservedNames[r.nameGenerator.GenerateNextName()] = struct{}{}
	}
	for _, p := range props {
		if r.generatePseudoNames {
			p.newName = "$" + p.oldName + "$"
		} else if p.newName == "" {
			p.newName = r.nameGenerator.GenerateNextName()
		}
		reservedNames[p.newName] = struct{}{}
		if _, filtered := r.filteredOutNames[p.oldNameID]; !filtered {
			p.newNameID = strpool.Intern(p.newName)
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

func (r *RenameProperties) visitor() *ast.Visitor[*NodeTraversal] {
	return &ast.Visitor[*NodeTraversal]{
		GetProp:         func(_ *NodeTraversal, gp ast.GetProp) { r.maybeMarkCandidate(gp.Node()) },
		OptChainGetProp: func(_ *NodeTraversal, gp ast.OptChainGetProp) { r.maybeMarkCandidate(gp.Node()) },
		MemberFieldDef:  func(_ *NodeTraversal, mf ast.MemberFieldDef) { r.maybeMarkCandidate(mf.Node()) },
		GetElem:         func(_ *NodeTraversal, ge ast.GetElem) { r.recordQuotedIndex(ge.Index()) },
		OptChainGetElem: func(_ *NodeTraversal, ge ast.OptChainGetElem) { r.recordQuotedIndex(ge.Index()) },
		Call:            r.visitCall,
		MemberFunctionDef: func(_ *NodeTraversal, mf ast.MemberFunctionDef) {
			jutil.CheckState(!mf.Node().IsQuotedStringKey())
			if mf.IsEs6Constructor() {
				r.externedNames[nodeStringID(mf.Node())] = struct{}{}
			} else {
				r.maybeMarkCandidate(mf.Node())
			}
		},
		GetterDef: func(_ *NodeTraversal, gd ast.GetterDef) { r.visitKeyOrAccessor(gd.Node()) },
		SetterDef: func(_ *NodeTraversal, sd ast.SetterDef) { r.visitKeyOrAccessor(sd.Node()) },
		StringKey: func(_ *NodeTraversal, sk ast.StringKey) { r.visitKeyOrAccessor(sk.Node()) },
		Function:  func(_ *NodeTraversal, fn ast.Function) { r.visitFunction(fn) },
	}
}

func (r *RenameProperties) recordQuotedIndex(index ast.Expr) {
	// If this is a quoted property access (e.g. x['myprop']), we need to
	// ensure that we never rename some other property in a way that
	// could conflict with this quoted name.
	if lit, ok := index.AsStringLit(); ok {
		r.quotedNames[nodeStringID(lit.Node())] = struct{}{}
	}
}

func (r *RenameProperties) visitKeyOrAccessor(n rhino.Node) {
	if n.IsQuotedStringKey() {
		// Ensure that we never rename some other property in a way
		// that could conflict with this quoted key.
		r.quotedNames[nodeStringID(n)] = struct{}{}
	} else {
		r.maybeMarkCandidate(n)
	}
}

func (r *RenameProperties) visitCall(t *NodeTraversal, call ast.Call) {
	// We replace property renaming function calls with a string
	// containing the renamed property.
	callee := call.Callee().Node()
	if !r.compiler.GetCodingConvention().IsPropertyRenameFunction(callee) {
		return
	}
	n := call.Node()
	r.callSites = append(r.callSites, propCallSite{call: call, parent: n.GetParent()})
	r.countCallCandidates(t, call)
	if !r.allPropertiesEligible && !r.propertyRenameEligibilityFilter(n) {
		if r.filteredOutNames == nil {
			r.filteredOutNames = make(map[strpool.ID]struct{})
		}
		if firstArg, ok := call.FirstArg(); ok {
			for name := range strings.SplitSeq(firstArg.Node().GetString(), ".") {
				r.filteredOutNames[strpool.Intern(name)] = struct{}{}
			}
		}
	}
}

func (r *RenameProperties) visitFunction(fn ast.Function) {
	// We eliminate any stub implementations of JSCompiler_renameProperty
	// that we encounter.
	n := fn.Node()
	parent := n.GetParent()
	if NodeUtil_IsFunctionDeclaration(n) {
		if nodeStringID(fn.NameNode().Node()) == idJscPropertyNameFn {
			r.toRemove = append(r.toRemove, n)
		}
	} else if parent.IsName() && nodeStringID(parent) == idJscPropertyNameFn {
		if decl, ok := ast.AsNameDeclaration(parent.GetParent()); ok && decl.Token() == rhino.Token_VAR {
			r.toRemove = append(r.toRemove, parent)
		}
	} else if assign, ok := ast.AsAssign(parent); ok && NodeUtil_IsFunctionExpression(n) {
		// JSCompiler does not handle optional calls to property rename function
		if gp, ok := assign.Target().AsGetProp(); ok &&
			r.compiler.GetCodingConvention().IsPropertyRenameFunction(gp.Node()) {
			if exprStmt, ok := ast.AsExprStmt(parent.GetParent()); ok &&
				NodeUtil_IsStatementBlock(exprStmt.Node().GetParent()) {
				r.toRemove = append(r.toRemove, exprStmt.Node())
			}
		}
	}
}

// maybeMarkCandidate: if a property node is eligible for renaming, stashes a
// reference to it and increments the property name's access count.
//
// n is the STRING node for a property.
//
// Java: RenameProperties.ProcessProperties.maybeMarkCandidate(Node)
func (r *RenameProperties) maybeMarkCandidate(n rhino.Node) {
	id := nodeStringID(n)
	if _, externed := r.externedNames[id]; !externed {
		if !r.allPropertiesEligible && !r.propertyRenameEligibilityFilter(n) {
			if r.filteredOutNames == nil {
				r.filteredOutNames = make(map[strpool.ID]struct{})
			}
			r.filteredOutNames[id] = struct{}{}
		}
		prop := r.countPropertyOccurrenceID(id)
		r.candidates = append(r.candidates, propRenameCandidate{node: n, prop: prop})
	}
}

// countCallCandidates counts references to property names that occur in a
// special function call.
//
// Java: RenameProperties.ProcessProperties.countCallCandidates(NodeTraversal, Node)
func (r *RenameProperties) countCallCandidates(t *NodeTraversal, call ast.Call) {
	fnName := call.Callee().Node().GetString()
	firstArg, hasArg := call.FirstArg()
	lit, isStr := firstArg.AsStringLit()
	if !hasArg || !isStr {
		t.Report(call.Node(), RenameProperties_BAD_CALL, fnName)
		return
	}

	for name := range strings.SplitSeq(lit.Value(), ".") {
		if !rhino.TokenStream_IsJSIdentifier(name) {
			t.Report(call.Node(), RenameProperties_BAD_ARG, fnName)
		} else {
			id := strpool.Intern(name)
			if _, externed := r.externedNames[id]; !externed {
				r.countPropertyOccurrenceID(id)
			}
		}
	}
}

func (r *RenameProperties) countPropertyOccurrenceID(id strpool.ID) *RenameProperties_Property {
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
	prop := r.propSlab.New()
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
	NodeTraversal_TraverseVisitor(r.compiler, root, &ast.Visitor[*NodeTraversal]{
		Call: func(_ *NodeTraversal, call ast.Call) {
			if r.compiler.GetCodingConvention().IsPropertyRenameFunction(call.Callee().Node()) {
				if propName, ok := call.FirstArg(); ok {
					call.Node().ReplaceWith(propName.Detach().Node())
					r.compiler.ReportChangeToEnclosingScope(propName.Node())
				}
			}
		},
	})
}
