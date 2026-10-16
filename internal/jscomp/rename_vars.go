// Port of com/google/javascript/jscomp/RenameVars.java @ closure-compiler v20260415

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
	"cmp"
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RenameVars renames all the variables names into short names, to reduce
// code size and also to obfuscate the code.
type RenameVars struct {
	compiler AbstractCompiler

	// List of global NAME nodes
	globalNameNodes *jutil.ArrayList[*rhino.Node]

	// List of local NAME nodes
	localNameNodes *jutil.ArrayList[*rhino.Node]

	// Mapping of original names for change detection
	originalNameByNode *jutil.LinkedHashMap[*rhino.Node, string]

	// Maps a name node to its pseudo name, null if we are not generating so
	// there will be no overhead unless we are debugging.
	pseudoNameMap *jutil.LinkedHashMap[*rhino.Node, string]

	// Set of extern variable names
	externNames *jutil.LinkedHashSet[string]

	// Set of reserved variable names
	reservedNames *jutil.LinkedHashSet[string]

	// The renaming map
	renameMap *jutil.LinkedHashMap[string, string]

	// The previously used rename map.
	prevUsedRenameMap *VariableMap

	// The global name prefix
	prefix string

	// Counter for each assignment
	assignmentCount int

	// Logic for bleeding functions, where the name leaks into the outer
	// scope on IE but not on other browsers.
	localBleedingFunctions *jutil.LinkedHashSet[*Var]
	// Go-only: language-adapter — the ArrayListMultimap<Scope, Var> is a map
	// from scope to its value list; only get (an empty list when absent) and
	// put are used, so the key order is never observed.
	localBleedingFunctionsPerScope *jutil.LinkedHashMap[*Scope, *jutil.ArrayList[*Var]]

	// Maps an old name to a new name assignment
	assignments *jutil.LinkedHashMap[string, *RenameVars_Assignment]

	// Whether renaming should apply to local variables only.
	localRenamingOnly bool

	preferStableNames bool

	// Characters that shouldn't be used in variable names.
	reservedCharacters jutil.Set[uint16]

	// Shared name generator
	nameGenerator NameGenerator
}

// Limit on number of locals in a scope for temporary local renaming when
// preferStableNames is true.
const RenameVars_MAX_LOCALS_IN_SCOPE_TO_TEMP_RENAME = 1000

// RenameVars_Assignment is the inner class Assignment.
type RenameVars_Assignment struct {
	isLocal           bool
	oldName           string
	orderOfOccurrence int
	newName           *string
	count             int // Number of times this is referenced
}

// Java: RenameVars.Assignment(String)
func NewRenameVars_Assignment(outer *RenameVars, name string) *RenameVars_Assignment {
	a := &RenameVars_Assignment{}
	a.isLocal = jutil.StringStartsWith(name, RenameVars_LOCAL_VAR_PREFIX)
	a.oldName = name
	a.newName = nil
	a.count = 0

	// Represents the order at which a symbol appears in the source.
	a.orderOfOccurrence = outer.assignmentCount
	outer.assignmentCount++
	return a
}

// SetNewName assigns the new name.
//
// Java: RenameVars.Assignment.setNewName(String)
func (a *RenameVars_Assignment) SetNewName(newName string) {
	jutil.CheckState(a.newName == nil)
	a.newName = &newName
}

// A prefix to distinguish temporary local names from global names
const RenameVars_LOCAL_VAR_PREFIX = "L "

// NewRenameVars: nameGenerator is a shared NameGenerator that this instance
// can use; the instance may reset or reconfigure it, so the caller should
// not expect any state to be preserved.
//
// Java: RenameVars(AbstractCompiler, String, boolean, boolean, boolean, VariableMap, Set<Character>, Set<String>, NameGenerator)
func NewRenameVars(
	compiler AbstractCompiler,
	prefix *string,
	localRenamingOnly bool,
	generatePseudoNames bool,
	preferStableNames bool,
	prevUsedRenameMap *VariableMap,
	reservedCharacters jutil.Set[uint16],
	reservedNames jutil.Set[string],
	nameGenerator NameGenerator) *RenameVars {
	r := &RenameVars{
		globalNameNodes:                jutil.NewArrayList[*rhino.Node](),
		localNameNodes:                 jutil.NewArrayList[*rhino.Node](),
		originalNameByNode:             jutil.NewLinkedHashMap[*rhino.Node, string](),
		renameMap:                      jutil.NewLinkedHashMap[string, string](),
		localBleedingFunctions:         jutil.NewLinkedHashSet[*Var](),
		localBleedingFunctionsPerScope: jutil.NewLinkedHashMap[*Scope, *jutil.ArrayList[*Var]](),
		assignments:                    jutil.NewLinkedHashMap[string, *RenameVars_Assignment](),
	}
	r.compiler = compiler
	// Strings.nullToEmpty
	if prefix != nil {
		r.prefix = *prefix
	} else {
		r.prefix = ""
	}
	r.localRenamingOnly = localRenamingOnly
	if generatePseudoNames {
		r.pseudoNameMap = jutil.NewLinkedHashMap[*rhino.Node, string]()
	} else {
		r.pseudoNameMap = nil
	}
	r.prevUsedRenameMap = prevUsedRenameMap
	r.reservedCharacters = reservedCharacters
	r.preferStableNames = preferStableNames
	if reservedNames == nil {
		r.reservedNames = jutil.NewLinkedHashSet[string]()
	} else {
		r.reservedNames = jutil.NewLinkedHashSet[string]()
		for name := range reservedNames.All() {
			r.reservedNames.Add(name)
		}
	}
	r.nameGenerator = nameGenerator
	return r
}

// RenameVars_ProcessVars iterates through the nodes, collect all the NAME
// nodes that need to be renamed, and count how many times each variable name
// is referenced.
//
// Keep track of all name references in globalNameNodes, and localNameNodes.
//
// To get shorter local variable renaming, we rename local variables to a
// temporary name "LOCAL_VAR_PREFIX + index" where index is the index of the
// variable declared in the local scope stack. e.g.
//
//	Foo(fa, fb) {
//	  var c = function(d, e) { return fa; }
//	}
//
// The indexes are: fa:0, fb:1, c:2, d:3, e:4
//
// In that way, local variable names are reused in each global function. e.g.
// the final code might look like
//
//	function x(a,b) { ... }
//	function y(a,b,c) { ... }
type RenameVars_ProcessVars struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *RenameVars
}

// Java: RenameVars.ProcessVars.enterScope(NodeTraversal)
func (p *RenameVars_ProcessVars) EnterScope(t *NodeTraversal) {
	r := p.outer
	if t.InGlobalHoistScope() || !r.ShouldTemporarilyRenameLocalsInScope(t.GetScope()) {
		return
	}
	scope := t.GetScope()
	for current := range scope.GetVarIterable() {
		if current.IsBleedingFunction() {
			r.localBleedingFunctions.Add(current)
			// localBleedingFunctionsPerScope.put(scope.getParent(), current)
			parentScope := scope.GetParent()
			list := r.localBleedingFunctionsPerScope.Get(parentScope)
			if list == nil {
				list = jutil.NewArrayList[*Var]()
				r.localBleedingFunctionsPerScope.Put(parentScope, list)
			}
			list.Add(current)
		}
	}
}

// Java: RenameVars.ProcessVars.exitScope(NodeTraversal)
func (p *RenameVars_ProcessVars) ExitScope(t *NodeTraversal) {}

// Java: RenameVars.ProcessVars.visit(NodeTraversal, Node, Node)
func (p *RenameVars_ProcessVars) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	r := p.outer
	if !(n.IsName() || n.IsImportStar()) {
		return
	}

	name := n.GetString()

	// Ignore anonymous functions and classes.
	if name == "" {
		return
	}

	// "import {x as y} from 'm';"
	// Skip x because it's not a variable in this scope.
	if parent.IsImportSpec() && parent.HasTwoChildren() && parent.GetFirstChild() == n {
		return
	}

	// Is this local or Global?
	// Bleeding functions should be treated as part of their outer
	// scope, because IE has bugs in how it handles bleeding
	// functions.
	var_ := t.GetScope().GetVar(name)
	local := var_ != nil &&
		var_.IsLocal() &&
		(var_.GetScope().(*Scope).GetParent().IsLocal() || !var_.IsBleedingFunction())

	// Never rename references to the arguments array
	if var_ != nil && var_.IsArguments() {
		r.reservedNames.Add(name)
		return
	}

	// Are we renaming global variables?
	if !local && r.localRenamingOnly {
		r.reservedNames.Add(name)
		return
	}

	// Check if we can rename this.
	if !r.OkToRenameVar(name, local) {
		if local {
			// Blindly de-uniquify for the Prototype library for
			// http://blickly.github.io/closure-compiler-issues/#103
			newName := MakeDeclaredNamesUnique_ContextualRenameInverter_GetOriginalName(name)
			if newName != name {
				n.SetString(newName)
			}
		}
		return
	}

	if r.pseudoNameMap != nil {
		r.RecordPseudoName(n)
	}

	if local && r.ShouldTemporarilyRenameLocalsInScope(var_.GetScope().(*Scope)) {
		// Give local variables a temporary name based on the
		// variable's index in the scope to enable name reuse across
		// locals in independent scopes.
		tempName := RenameVars_LOCAL_VAR_PREFIX + strconv.Itoa(r.GetLocalVarIndex(var_))
		p.IncCount(tempName)
		r.localNameNodes.Add(n)
		// Remember the original string in a name before it's temporarily filled with an "L".
		r.originalNameByNode.Put(n, n.GetString())
		n.SetString(tempName)
	} else if var_ != nil { // Not an extern
		// If it's global, increment global count
		p.IncCount(name)
		r.globalNameNodes.Add(n)
	}
}

// IncCount increments count of an assignment
//
// Java: RenameVars.ProcessVars.incCount(String)
func (p *RenameVars_ProcessVars) IncCount(name string) {
	// assignments.computeIfAbsent(name, Assignment::new)
	s := p.outer.assignments.Get(name)
	if s == nil {
		s = NewRenameVars_Assignment(p.outer, name)
		p.outer.assignments.Put(name, s)
	}
	s.count++
}

// RenameVars_FrequencyComparator sorts Assignment objects by their count,
// breaking ties by their order of occurrence in the source to ensure a
// deterministic total ordering.
//
// Java: RenameVars.frequencyComparator(Assignment, Assignment)
func RenameVars_FrequencyComparator(a1 *RenameVars_Assignment, a2 *RenameVars_Assignment) int {
	if a1.count != a2.count {
		return a2.count - a1.count
	}
	// Break a tie using the order in which the variable first appears in
	// the source.
	return RenameVars_ORDER_OF_OCCURRENCE_COMPARATOR(a1, a2)
}

// RenameVars_ORDER_OF_OCCURRENCE_COMPARATOR sorts Assignment objects by the
// order the variable name first appears in the source.
var RenameVars_ORDER_OF_OCCURRENCE_COMPARATOR = func(a1 *RenameVars_Assignment, a2 *RenameVars_Assignment) int {
	// Comparator.comparingInt: Integer.compare.
	return cmp.Compare(a1.orderOfOccurrence, a2.orderOfOccurrence)
}

// Java: RenameVars.process(Node, Node)
func (r *RenameVars) Process(externs *rhino.Node, root *rhino.Node) {
	r.externNames = NodeUtil_CollectExternVariableNames(r.compiler, externs)

	r.originalNameByNode.Clear()

	// Do variable reference counting.
	NodeTraversal_Traverse(r.compiler, root, &RenameVars_ProcessVars{outer: r})

	// Make sure that new names don't overlap with extern names.
	r.reservedNames.AddAll(r.externNames.ToSlice()...)

	// Rename vars, sorted by frequency of occurrence to minimize code size.
	varsByFrequency := jutil.NewTreeSet(RenameVars_FrequencyComparator)
	varsByFrequency.AddAll(r.assignments.Values())

	// First try to reuse names from an earlier compilation.
	if r.prevUsedRenameMap != nil {
		r.ReusePreviouslyUsedVariableMap(varsByFrequency)
	}

	// Assign names, sorted by descending frequency to minimize code size.
	r.AssignNames(varsByFrequency)

	// Rename the globals!
	for _, n := range r.globalNameNodes.Slice() {
		r.SetNameAndReport(n, r.GetNewGlobalName(n))
	}

	// Rename the locals!
	for _, n := range r.localNameNodes.Slice() {
		r.SetNameAndReport(n, r.GetNewLocalName(n))
	}
}

// Java: RenameVars.setNameAndReport(Node, String)
func (r *RenameVars) SetNameAndReport(n *rhino.Node, newName *string) {
	// A null newName, indicates it should not be renamed.
	if newName != nil && *newName != n.GetString() {
		n.SetString(*newName)

		// Only mark changes if the final name change is different than it was original before being
		// filled with the "L" temporary name.
		if originalName, ok := r.originalNameByNode.Lookup(n); !ok || *newName != originalName {
			r.compiler.ReportChangeToEnclosingScope(n)
			parent := n.GetParent()
			if parent.IsFunction() && NodeUtil_IsFunctionDeclaration(parent) {
				// If we are renaming a function declaration, make sure the containing scope
				// has the opportunity to act on the change.
				r.compiler.ReportChangeToEnclosingScope(parent)
			}
		}
	}
}

// Java: RenameVars.getNewGlobalName(Node)
func (r *RenameVars) GetNewGlobalName(n *rhino.Node) *string {
	oldName := n.GetString()
	a := r.assignments.Get(oldName)
	if a.newName != nil && *a.newName != oldName {
		if r.pseudoNameMap != nil {
			return r.pseudoNameMapGet(n)
		}
		return a.newName
	} else {
		return nil
	}
}

// Java: RenameVars.getNewLocalName(Node)
func (r *RenameVars) GetNewLocalName(n *rhino.Node) *string {
	oldTempName := n.GetString()
	a := r.assignments.Get(oldTempName)
	if *a.newName != oldTempName {
		if r.pseudoNameMap != nil {
			return r.pseudoNameMapGet(n)
		}
		return a.newName
	}
	return nil
}

// pseudoNameMapGet is pseudoNameMap.get(n): null when absent.
//
// Go-only: language-adapter — Map.get returning a nullable String.
func (r *RenameVars) pseudoNameMapGet(n *rhino.Node) *string {
	if v, ok := r.pseudoNameMap.Lookup(n); ok {
		return &v
	}
	return nil
}

// Java: RenameVars.recordPseudoName(Node)
func (r *RenameVars) RecordPseudoName(n *rhino.Node) {
	// Variable names should be in a different name space than
	// property pseudo names.
	r.pseudoNameMap.Put(n, "$"+n.GetString()+"$$")
}

// ReusePreviouslyUsedVariableMap runs through the assignments and reuses as
// many names as possible from the previously used variable map. Updates
// reservedNames with the set of names that were reused.
//
// Java: RenameVars.reusePreviouslyUsedVariableMap(SortedSet<Assignment>)
func (r *RenameVars) ReusePreviouslyUsedVariableMap(varsToRename *jutil.TreeSet[*RenameVars_Assignment]) {
	// If prevUsedRenameMap had duplicate values then this pass would be
	// non-deterministic.
	// In such a case, the following will throw an IllegalArgumentException.
	jutil.CheckNotNull(r.prevUsedRenameMap.GetNewNameToOriginalNameMap())
	for a := range varsToRename.All() {
		prevNewName := r.prevUsedRenameMap.LookupNewName(a.oldName)
		if prevNewName == nil || r.reservedNames.Contains(*prevNewName) {
			continue
		}

		if a.isLocal ||
			(!r.externNames.Contains(a.oldName) &&
				jutil.StringStartsWith(*prevNewName, r.prefix)) {
			r.reservedNames.Add(*prevNewName)
			r.FinalizeNameAssignment(a, *prevNewName)
		}
	}
}

// AssignNames determines which new names to substitute for the original
// names.
//
// Java: RenameVars.assignNames(SortedSet<Assignment>)
func (r *RenameVars) AssignNames(varsToRename *jutil.TreeSet[*RenameVars_Assignment]) {
	r.nameGenerator.Reset(r.reservedNames, r.prefix, r.reservedCharacters)

	globalNameGenerator := r.nameGenerator
	// Local variables never need a prefix.
	// Also, we need to avoid conflicts between global and local variable
	// names; we do this by having using the same generator (not two
	// instances). The case where global variables have a prefix (and
	// therefore we use two different generators) but a local variable name
	// might nevertheless conflict with a global one is not handled.
	var localNameGenerator NameGenerator
	if r.prefix == "" {
		localNameGenerator = globalNameGenerator
	} else {
		localNameGenerator = r.nameGenerator.Clone(r.reservedNames, "", r.reservedCharacters)
	}

	// Generated names and the assignments for non-local vars.
	var pendingAssignments []*RenameVars_Assignment
	var generatedNamesForAssignments []string

	for a := range varsToRename.All() {
		if a.newName != nil {
			continue
		}

		if r.externNames.Contains(a.oldName) {
			continue
		}

		var newName string
		if a.isLocal {
			// For local variable, we make the assignment right away.
			newName = localNameGenerator.GenerateNextName()
			r.FinalizeNameAssignment(a, newName)
		} else {
			// For non-local variable, delay finalizing the name assignment
			// until we know how many new names we'll have of length 2, 3, etc.
			newName = globalNameGenerator.GenerateNextName()
			pendingAssignments = append(pendingAssignments, a)
			generatedNamesForAssignments = append(generatedNamesForAssignments, newName)
		}
		r.reservedNames.Add(newName)
	}

	// Now that we have a list of generated names, and a list of variable
	// Assignment objects, we assign the generated names to the vars as
	// follows:
	// 1) The most frequent vars get the shorter names.
	// 2) If N number of vars are going to be assigned names of the same
	//    length, we assign the N names based on the order at which the vars
	//    first appear in the source. This makes the output somewhat less
	//    random, because symbols declared close together are assigned names
	//    that are quite similar. With this heuristic, the output is more
	//    compressible.
	//    For instance, the output may look like:
	//    var da = "..", ea = "..";
	//    function fa() { .. } function ga() { .. }

	numPendingAssignments := len(generatedNamesForAssignments)
	for i := 0; i < numPendingAssignments; {
		varsByOrderOfOccurrence := jutil.NewTreeSet(RenameVars_ORDER_OF_OCCURRENCE_COMPARATOR)

		// Add k number of Assignment to the set, where k is the number of
		// generated names of the same length.
		len_ := jutil.StringLength(generatedNamesForAssignments[i])
		for j := i; j < numPendingAssignments &&
			jutil.StringLength(generatedNamesForAssignments[j]) == len_; j++ {
			varsByOrderOfOccurrence.Add(pendingAssignments[j])
		}

		// Now, make the assignments
		for a := range varsByOrderOfOccurrence.All() {
			r.FinalizeNameAssignment(a, generatedNamesForAssignments[i])
			i++
		}
	}
}

// FinalizeNameAssignment makes a final name assignment.
//
// Java: RenameVars.finalizeNameAssignment(Assignment, String)
func (r *RenameVars) FinalizeNameAssignment(a *RenameVars_Assignment, newName string) {
	a.SetNewName(newName)

	// Keep track of the mapping
	r.renameMap.Put(a.oldName, newName)
}

// GetVariableMap gets the variable map.
//
// Java: RenameVars.getVariableMap()
func (r *RenameVars) GetVariableMap() *VariableMap {
	return NewVariableMap(jutil.NewLinkedHashMapFrom(r.renameMap))
}

// OkToRenameVar determines whether a variable name is okay to rename.
//
// Java: RenameVars.okToRenameVar(String, boolean)
func (r *RenameVars) OkToRenameVar(name string, isLocal bool) bool {
	return !r.compiler.GetCodingConvention().IsExportedWithLocal(name, isLocal /* local= */)
}

// GetLocalVarIndex returns the index within the scope stack. e.g. function
// Foo(a) { var b; function c(d) { } } a = 0, b = 1, c = 2, d = 3
//
// Java: RenameVars.getLocalVarIndex(Var)
func (r *RenameVars) GetLocalVarIndex(v *Var) int {
	num := v.GetIndex()
	s := v.GetScope().(*Scope).GetParent()
	if s == nil {
		panic(jutil.NewIllegalArgumentExceptionMsg("Var is not local"))
	}

	isBleedingIntoScope := s.GetParent() != nil && r.localBleedingFunctions.Contains(v)

	for s.GetParent() != nil {
		// localBleedingFunctionsPerScope.get(s): an empty list when absent.
		bleeding := r.localBleedingFunctionsPerScope.Get(s)
		if isBleedingIntoScope {
			index := -1
			if bleeding != nil {
				index = bleeding.IndexOf(v)
			}
			num += index + 1
			isBleedingIntoScope = false
		} else {
			if bleeding != nil {
				num += bleeding.Size()
			}
		}
		if r.ShouldTemporarilyRenameLocalsInScope(s) {
			num += s.GetVarCount()
		}
		s = s.GetParent()
	}
	return num
}

// ShouldTemporarilyRenameLocalsInScope returns true if the local variables
// in a scope should be given temporary names (eg, 'L 123') prior to renaming
// to allow reuse of names across scopes. With preferStableNames, temporary
// renaming is disabled if the number of locals in the scope is above a
// heuristic threshold to allow effective reuse of rename maps (see
// prevUsedRenameMap). In scopes with many variables the temporary name given
// to a variable is unlikely to be the same temporary name used when the
// rename map was created.
//
// Java: RenameVars.shouldTemporarilyRenameLocalsInScope(Scope)
func (r *RenameVars) ShouldTemporarilyRenameLocalsInScope(s *Scope) bool {
	return (!r.preferStableNames || s.GetVarCount() <= RenameVars_MAX_LOCALS_IN_SCOPE_TO_TEMP_RENAME)
}
