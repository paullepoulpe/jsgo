// Port of com/google/javascript/jscomp/RenameVars.java @ closure-compiler v20261021

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
	"slices"
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type renameVarsNodeRef struct {
	node       rhino.Node
	assignment *RenameVars_Assignment
}

// RenameVars renames all the variables names into short names, to reduce
// code size and also to obfuscate the code.
type RenameVars struct {
	compiler   AbstractCompiler
	convention CodingConvention

	globalRefChunks [][]renameVarsNodeRef
	localRefChunks  [][]renameVarsNodeRef

	// Maps a name node to its pseudo name, null if we are not generating so
	// there will be no overhead unless we are debugging.
	pseudoNameMap map[rhino.Node]string

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
	localBleedingFunctions         map[*Var]struct{}
	localBleedingFunctionsPerScope map[*Scope][]*Var

	// Maps an old name to a new name assignment
	assignments  *jutil.LinkedHashMap[string, *RenameVars_Assignment]
	localByIndex []*RenameVars_Assignment
	seenVars     map[*Var]renameVarDecision

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
	newNameVal        string
	count             int // Number of times this is referenced
}

// Java: RenameVars.Assignment(String)
func NewRenameVars_Assignment(outer *RenameVars, name string) *RenameVars_Assignment {
	return outer.allocAssignment(name, jutil.StringStartsWith(name, RenameVars_LOCAL_VAR_PREFIX))
}

type renameVarDecision struct {
	assignment *RenameVars_Assignment
	invertName bool
}

func (r *RenameVars) allocAssignment(name string, isLocal bool) *RenameVars_Assignment {
	a := &RenameVars_Assignment{
		isLocal:           isLocal,
		oldName:           name,
		orderOfOccurrence: r.assignmentCount,
	}
	r.assignmentCount++
	return a
}

func (r *RenameVars) appendLocalRef(n rhino.Node, a *RenameVars_Assignment) {
	nChunks := len(r.localRefChunks)
	if nChunks == 0 || len(r.localRefChunks[nChunks-1]) == cap(r.localRefChunks[nChunks-1]) {
		newCap := 1024
		if nChunks > 0 {
			newCap = min(cap(r.localRefChunks[nChunks-1])*2, 16384)
		}
		r.localRefChunks = append(r.localRefChunks, make([]renameVarsNodeRef, 0, newCap))
		nChunks++
	}
	r.localRefChunks[nChunks-1] = append(r.localRefChunks[nChunks-1], renameVarsNodeRef{node: n, assignment: a})
}

func (r *RenameVars) appendGlobalRef(n rhino.Node, a *RenameVars_Assignment) {
	nChunks := len(r.globalRefChunks)
	if nChunks == 0 || len(r.globalRefChunks[nChunks-1]) == cap(r.globalRefChunks[nChunks-1]) {
		newCap := 512
		if nChunks > 0 {
			newCap = min(cap(r.globalRefChunks[nChunks-1])*2, 8192)
		}
		r.globalRefChunks = append(r.globalRefChunks, make([]renameVarsNodeRef, 0, newCap))
		nChunks++
	}
	r.globalRefChunks[nChunks-1] = append(r.globalRefChunks[nChunks-1], renameVarsNodeRef{node: n, assignment: a})
}

// SetNewName assigns the new name.
//
// Java: RenameVars.Assignment.setNewName(String)
func (a *RenameVars_Assignment) SetNewName(newName string) {
	jutil.CheckState(a.newName == nil)
	a.newNameVal = newName
	a.newName = &a.newNameVal
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
		compiler:    compiler,
		convention:  compiler.GetCodingConvention(),
		renameMap:   jutil.NewLinkedHashMap[string, string](),
		assignments: jutil.NewLinkedHashMap[string, *RenameVars_Assignment](),
	}
	// Strings.nullToEmpty
	if prefix != nil {
		r.prefix = *prefix
	} else {
		r.prefix = ""
	}
	r.localRenamingOnly = localRenamingOnly
	if generatePseudoNames {
		r.pseudoNameMap = make(map[rhino.Node]string)
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
	scopeRoot := t.GetScopeRoot()
	if !scopeRoot.IsFunction() || scopeRoot.GetFirstChild().GetString() == "" || !NodeUtil_IsFunctionExpression(scopeRoot) {
		return
	}
	r := p.outer
	if t.InGlobalHoistScope() || !r.ShouldTemporarilyRenameLocalsInScope(t.GetScope()) {
		return
	}
	scope := t.GetScope()
	for _, current := range scope.varList {
		if current == nil {
			continue
		}
		if current.IsBleedingFunction() {
			if r.localBleedingFunctions == nil {
				r.localBleedingFunctions = make(map[*Var]struct{})
				r.localBleedingFunctionsPerScope = make(map[*Scope][]*Var)
			}
			r.localBleedingFunctions[current] = struct{}{}
			parentScope := scope.GetParent()
			r.localBleedingFunctionsPerScope[parentScope] = append(r.localBleedingFunctionsPerScope[parentScope], current)
		}
	}
}

// Java: RenameVars.ProcessVars.exitScope(NodeTraversal)
func (p *RenameVars_ProcessVars) ExitScope(t *NodeTraversal) {}

// Java: RenameVars.ProcessVars.visit(NodeTraversal, Node, Node)
func (p *RenameVars_ProcessVars) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	tok := n.GetToken()
	if tok != rhino.Token_NAME && tok != rhino.Token_IMPORT_STAR {
		return
	}

	id := n.GetStringID()

	// Ignore anonymous functions and classes.
	if id <= rhino.RhinoStringPool_EmptyID {
		return
	}

	// "import {x as y} from 'm';"
	// Skip x because it's not a variable in this scope.
	if parent.IsImportSpec() && parent.HasTwoChildren() && parent.GetFirstChild() == n {
		return
	}

	r := p.outer
	var_ := t.LookupVarByID(id)
	if var_ != nil {
		if dec, ok := r.seenVars[var_]; ok {
			if a := dec.assignment; a != nil {
				a.count++
				if r.pseudoNameMap != nil {
					r.RecordPseudoName(n)
				}
				if a.isLocal {
					r.appendLocalRef(n, a)
				} else {
					r.appendGlobalRef(n, a)
				}
			} else if dec.invertName {
				name := n.GetString()
				newName := MakeDeclaredNamesUnique_ContextualRenameInverter_GetOriginalName(name)
				if newName != name {
					n.SetString(newName)
				}
			}
			return
		}
		if r.seenVars == nil {
			r.seenVars = make(map[*Var]renameVarDecision, 256)
		}

		name := n.GetString()

		// First time seeing this var_.
		varScope := var_.scope
		local := var_.IsLocal() && (varScope.GetParent().IsLocal() || !var_.IsBleedingFunction())

		if var_.IsArguments() {
			r.reservedNames.Add(name)
			r.seenVars[var_] = renameVarDecision{}
			return
		}

		if !local && r.localRenamingOnly {
			r.reservedNames.Add(name)
			r.seenVars[var_] = renameVarDecision{}
			return
		}

		if !r.OkToRenameVar(name, local) {
			if local {
				newName := MakeDeclaredNamesUnique_ContextualRenameInverter_GetOriginalName(name)
				if newName != name {
					n.SetString(newName)
					r.seenVars[var_] = renameVarDecision{invertName: true}
					return
				}
			}
			r.seenVars[var_] = renameVarDecision{}
			return
		}

		if r.pseudoNameMap != nil {
			r.RecordPseudoName(n)
		}

		var a *RenameVars_Assignment
		if local && r.ShouldTemporarilyRenameLocalsInScope(varScope) {
			localIdx := r.GetLocalVarIndex(var_)
			if localIdx >= 0 && localIdx < len(r.localByIndex) {
				a = r.localByIndex[localIdx]
			}
			if a == nil {
				tempName := RenameVars_LOCAL_VAR_PREFIX + strconv.Itoa(localIdx)
				if localIdx < 0 {
					a = r.assignments.Get(tempName)
				}
				if a == nil {
					a = r.allocAssignment(tempName, true)
					r.assignments.Put(tempName, a)
				}
				if localIdx >= 0 {
					if localIdx >= cap(r.localByIndex) {
						newCap := max(localIdx+1, cap(r.localByIndex)*2, 32)
						grown := make([]*RenameVars_Assignment, localIdx+1, newCap)
						copy(grown, r.localByIndex)
						r.localByIndex = grown
					} else if localIdx >= len(r.localByIndex) {
						r.localByIndex = r.localByIndex[:localIdx+1]
					}
					r.localByIndex[localIdx] = a
				}
			}
			a.count++
			r.appendLocalRef(n, a)
		} else {
			a = r.assignments.Get(name)
			if a == nil {
				a = r.allocAssignment(name, false)
				r.assignments.Put(name, a)
			}
			a.count++
			r.appendGlobalRef(n, a)
		}
		r.seenVars[var_] = renameVarDecision{assignment: a}
		return
	}

	// var_ == nil (undeclared or extern)
	name := n.GetString()
	if r.localRenamingOnly {
		r.reservedNames.Add(name)
		return
	}
	if !r.OkToRenameVar(name, false) {
		return
	}
	if r.pseudoNameMap != nil {
		r.RecordPseudoName(n)
	}
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
func (r *RenameVars) Process(externs rhino.Node, root rhino.Node) {
	r.externNames = NodeUtil_CollectExternVariableNames(r.compiler, externs)

	r.globalRefChunks = nil
	r.localRefChunks = nil

	// Do variable reference counting.
	sc := NewSyntacticScopeCreator(r.compiler)
	b := NodeTraversal_Builder{
		compiler:     r.compiler,
		callback:     &RenameVars_ProcessVars{outer: r},
		scopeCreator: sc,
	}
	b.Traverse(root)
	r.localBleedingFunctions = nil
	r.localBleedingFunctionsPerScope = nil

	// Make sure that new names don't overlap with extern names.
	for name := range r.externNames.All() {
		r.reservedNames.Add(name)
	}

	// Rename vars, sorted by frequency of occurrence to minimize code size.
	varsByFrequency := make([]*RenameVars_Assignment, 0, r.assignments.Size())
	for a := range r.assignments.Values() {
		varsByFrequency = append(varsByFrequency, a)
	}
	slices.SortFunc(varsByFrequency, RenameVars_FrequencyComparator)

	// First try to reuse names from an earlier compilation.
	if r.prevUsedRenameMap != nil {
		r.reusePreviouslyUsedVariableMapSlice(varsByFrequency)
	}

	// Assign names, sorted by descending frequency to minimize code size.
	r.assignNamesSlice(varsByFrequency)

	// Rename the globals!
	for _, chunk := range r.globalRefChunks {
		for _, ref := range chunk {
			n := ref.node
			a := ref.assignment
			oldName := n.GetString()
			if a.newName != nil && *a.newName != oldName {
				newName := *a.newName
				if r.pseudoNameMap != nil {
					pName, ok := r.pseudoNameMap[n]
					if !ok {
						continue
					}
					newName = pName
				}
				if newName != oldName {
					n.SetString(newName)
					r.compiler.ReportChangeToEnclosingScope(n)
					parent := n.GetParent()
					if parent.IsFunction() && NodeUtil_IsFunctionDeclaration(parent) {
						r.compiler.ReportChangeToEnclosingScope(parent)
					}
				}
			}
		}
	}

	// Rename the locals!
	for _, chunk := range r.localRefChunks {
		for _, ref := range chunk {
			n := ref.node
			a := ref.assignment
			if *a.newName != a.oldName {
				newName := *a.newName
				if r.pseudoNameMap != nil {
					pName, ok := r.pseudoNameMap[n]
					if !ok {
						continue
					}
					newName = pName
				}
				originalName := n.GetString()
				if newName != originalName {
					n.SetString(newName)
					r.compiler.ReportChangeToEnclosingScope(n)
					parent := n.GetParent()
					if parent.IsFunction() && NodeUtil_IsFunctionDeclaration(parent) {
						r.compiler.ReportChangeToEnclosingScope(parent)
					}
				}
			}
		}
	}

	r.globalRefChunks = nil
	r.localRefChunks = nil
	r.seenVars = nil
}

// Java: RenameVars.recordPseudoName(Node)
func (r *RenameVars) RecordPseudoName(n rhino.Node) {
	// Variable names should be in a different name space than
	// property pseudo names.
	r.pseudoNameMap[n] = "$" + n.GetString() + "$$"
}

func (r *RenameVars) reusePreviouslyUsedVariableMapSlice(varsToRename []*RenameVars_Assignment) {
	jutil.CheckNotNullPtr(r.prevUsedRenameMap)
	for _, a := range varsToRename {
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

func (r *RenameVars) assignNamesSlice(varsToRename []*RenameVars_Assignment) {
	r.nameGenerator.Reset(r.reservedNames, r.prefix, r.reservedCharacters)

	globalNameGenerator := r.nameGenerator
	var localNameGenerator NameGenerator
	if r.prefix == "" {
		localNameGenerator = globalNameGenerator
	} else {
		localNameGenerator = r.nameGenerator.Clone(r.reservedNames, "", r.reservedCharacters)
	}

	var pendingAssignments []*RenameVars_Assignment
	var generatedNamesForAssignments []string

	for _, a := range varsToRename {
		if a.newName != nil {
			continue
		}

		if r.externNames.Contains(a.oldName) {
			continue
		}

		var newName string
		if a.isLocal {
			newName = localNameGenerator.GenerateNextName()
			r.FinalizeNameAssignment(a, newName)
		} else {
			newName = globalNameGenerator.GenerateNextName()
			pendingAssignments = append(pendingAssignments, a)
			generatedNamesForAssignments = append(generatedNamesForAssignments, newName)
		}
		r.reservedNames.Add(newName)
	}

	numPendingAssignments := len(generatedNamesForAssignments)
	for i := 0; i < numPendingAssignments; {
		len_ := jutil.StringLength(generatedNamesForAssignments[i])
		j := i + 1
		for j < numPendingAssignments &&
			jutil.StringLength(generatedNamesForAssignments[j]) == len_ {
			j++
		}
		batch := pendingAssignments[i:j]
		slices.SortFunc(batch, RenameVars_ORDER_OF_OCCURRENCE_COMPARATOR)
		for _, a := range batch {
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
	return !r.convention.IsExportedWithLocal(name, isLocal /* local= */)
}

// GetLocalVarIndex returns the index within the scope stack. e.g. function
// Foo(a) { var b; function c(d) { } } a = 0, b = 1, c = 2, d = 3
//
// Java: RenameVars.getLocalVarIndex(Var)
func (r *RenameVars) GetLocalVarIndex(v *Var) int {
	num := v.GetIndex()
	s := v.scope.GetParent()
	if s == nil {
		panic(jutil.NewIllegalArgumentExceptionMsg("Var is not local"))
	}

	if r.localBleedingFunctions == nil {
		for s.GetParent() != nil {
			if r.ShouldTemporarilyRenameLocalsInScope(s) {
				num += s.GetVarCount()
			}
			s = s.GetParent()
		}
		return num
	}

	_, hasBleeding := r.localBleedingFunctions[v]
	isBleedingIntoScope := s.GetParent() != nil && hasBleeding

	for s.GetParent() != nil {
		bleeding := r.localBleedingFunctionsPerScope[s]
		if isBleedingIntoScope {
			num += slices.Index(bleeding, v) + 1
			isBleedingIntoScope = false
		} else {
			num += len(bleeding)
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
