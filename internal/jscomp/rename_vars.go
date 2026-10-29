// Port of com/google/javascript/jscomp/RenameVars.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/RenamingMap.java @ closure-compiler v20261021

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
 * Copyright 2013 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

var (
	renameVarsVisitTokens = func() ast.TokenSet {
		var s ast.TokenSet
		s.Add(rhino.Token_NAME)
		s.Add(rhino.Token_IMPORT_STAR)
		return s
	}()
	renameVarsLocalTempNames = func() [256]string {
		var t [256]string
		for i := range t {
			t[i] = RenameVars_LOCAL_VAR_PREFIX + strconv.Itoa(i)
		}
		return t
	}()
)

type renameVarsNodeRef struct {
	node       rhino.Node
	assignment *RenameVars_Assignment
}

type renameVarsMapEntry struct {
	oldName string
	newName string
}

// RenameVars renames all the variables names into short names, to reduce
// code size and also to obfuscate the code.
type RenameVars struct {
	compiler   AbstractCompiler
	convention CodingConvention

	globalRefChunks [][]renameVarsNodeRef
	localRefChunks  [][]renameVarsNodeRef

	// Maps a name node to its pseudo name, nil if we are not generating so
	// there will be no overhead unless we are debugging.
	pseudoNameMap map[rhino.Node]string

	// Set of extern variable names
	externNames map[string]struct{}

	// Set of reserved variable names
	reservedNames stringMapSet

	// Recorded finalized renames in deterministic assignment order
	finalized []renameVarsMapEntry

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

	// Maps an old name to a new name assignment, plus insertion-ordered list
	assignmentByName map[string]*RenameVars_Assignment
	assignmentList   []*RenameVars_Assignment
	assignmentSlab   []RenameVars_Assignment
	localByIndex     []*RenameVars_Assignment
	seenVars         map[*Var]renameVarDecision

	// Whether renaming should apply to local variables only.
	localRenamingOnly bool

	preferStableNames bool

	// Characters that shouldn't be used in variable names.
	reservedCharacters ASCIICharSet

	// Shared name generator
	nameGenerator NameGenerator
}

// Limit on number of locals in a scope for temporary local renaming when
// preferStableNames is true.
const RenameVars_MAX_LOCALS_IN_SCOPE_TO_TEMP_RENAME = 1000

// RenameVars_Assignment is the inner class Assignment.
type RenameVars_Assignment struct {
	isLocal           bool
	hasNewName        bool
	changedFromOld    bool
	oldName           string
	orderOfOccurrence int
	newName           string
	newNameID         strpool.ID
	count             int // Number of times this is referenced
}

// Java: RenameVars.Assignment(String)
func NewRenameVars_Assignment(outer *RenameVars, name string) *RenameVars_Assignment {
	return outer.allocAssignment(name, strings.HasPrefix(name, RenameVars_LOCAL_VAR_PREFIX))
}

type renameVarDecision struct {
	assignment *RenameVars_Assignment
	invertID   strpool.ID
}

func (r *RenameVars) allocAssignment(name string, isLocal bool) *RenameVars_Assignment {
	if len(r.assignmentSlab) == 0 {
		r.assignmentSlab = make([]RenameVars_Assignment, 128)
	}
	a := &r.assignmentSlab[0]
	r.assignmentSlab = r.assignmentSlab[1:]
	a.isLocal = isLocal
	a.oldName = name
	a.orderOfOccurrence = r.assignmentCount
	r.assignmentCount++
	return a
}

func (r *RenameVars) putAssignment(name string, a *RenameVars_Assignment) {
	r.assignmentByName[name] = a
	r.assignmentList = append(r.assignmentList, a)
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
	jutil.CheckState(!a.hasNewName)
	a.newName = newName
	a.newNameID = strpool.Intern(newName)
	a.changedFromOld = newName != a.oldName
	a.hasNewName = true
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
	prefix string,
	localRenamingOnly bool,
	generatePseudoNames bool,
	preferStableNames bool,
	prevUsedRenameMap *VariableMap,
	reservedCharacters ASCIICharSet,
	reservedNames jutil.Set[string],
	nameGenerator NameGenerator) *RenameVars {
	reservedCap := 32
	if reservedNames != nil {
		reservedCap = max(reservedCap, reservedNames.Size()+32)
	}
	resMap := make(stringMapSet, reservedCap)
	if reservedNames != nil {
		for name := range reservedNames.All() {
			resMap[name] = struct{}{}
		}
	}
	r := &RenameVars{
		compiler:           compiler,
		convention:         compiler.GetCodingConvention(),
		reservedNames:      resMap,
		prevUsedRenameMap:  prevUsedRenameMap,
		prefix:             prefix,
		assignmentByName:   make(map[string]*RenameVars_Assignment, 128),
		localRenamingOnly:  localRenamingOnly,
		preferStableNames:  preferStableNames,
		reservedCharacters: reservedCharacters,
		nameGenerator:      nameGenerator,
	}
	if generatePseudoNames {
		r.pseudoNameMap = make(map[rhino.Node]string)
	}
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

func (*RenameVars_ProcessVars) visitTokens() *ast.TokenSet {
	return &renameVarsVisitTokens
}

// Java: RenameVars.ProcessVars.enterScope(NodeTraversal)
func (p *RenameVars_ProcessVars) EnterScope(t *NodeTraversal) {
	scopeRoot := t.GetScopeRoot()
	fn, ok := ast.AsFunction(scopeRoot)
	if !ok || fn.NameNode().Node().GetStringID() <= rhino.RhinoStringPool_EmptyID || !NodeUtil_IsFunctionExpression(scopeRoot) {
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
			} else if dec.invertID != 0 {
				setNodeStringID(n, dec.invertID)
			}
			return
		}
		if r.seenVars == nil {
			r.seenVars = make(map[*Var]renameVarDecision, 256)
		}

		// First time seeing this var_.
		varScope := var_.scope
		local := var_.IsLocal() && (varScope.GetParent().IsLocal() || !var_.IsBleedingFunction())

		if var_.IsArguments() {
			r.reservedNames["arguments"] = struct{}{}
			r.seenVars[var_] = renameVarDecision{}
			return
		}

		name := n.GetString()

		if !local && r.localRenamingOnly {
			r.reservedNames[name] = struct{}{}
			r.seenVars[var_] = renameVarDecision{}
			return
		}

		if !r.OkToRenameVar(name, local) {
			if local {
				newName := MakeDeclaredNamesUnique_ContextualRenameInverter_GetOriginalName(name)
				if newName != name {
					invID := strpool.Intern(newName)
					setNodeStringID(n, invID)
					r.seenVars[var_] = renameVarDecision{invertID: invID}
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
				var tempName string
				if localIdx >= 0 && localIdx < len(renameVarsLocalTempNames) {
					tempName = renameVarsLocalTempNames[localIdx]
				} else {
					tempName = RenameVars_LOCAL_VAR_PREFIX + strconv.Itoa(localIdx)
				}
				if localIdx < 0 {
					a = r.assignmentByName[tempName]
				}
				if a == nil {
					a = r.allocAssignment(tempName, true)
					r.putAssignment(tempName, a)
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
			a = r.assignmentByName[name]
			if a == nil {
				a = r.allocAssignment(name, false)
				r.putAssignment(name, a)
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
		r.reservedNames[name] = struct{}{}
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
	externSet := NodeUtil_CollectExternVariableNames(r.compiler, externs)
	r.externNames = make(map[string]struct{}, externSet.Size())
	for name := range externSet.All() {
		r.externNames[name] = struct{}{}
	}

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
	// order-independent: populates unordered reservedNames membership set
	for name := range r.externNames {
		r.reservedNames[name] = struct{}{}
	}

	// Rename vars, sorted by frequency of occurrence to minimize code size.
	varsByFrequency := make([]*RenameVars_Assignment, len(r.assignmentList))
	copy(varsByFrequency, r.assignmentList)
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
			oldID := strpool.ID(n.GetStringID())
			if a.hasNewName && a.newNameID != oldID {
				newID := a.newNameID
				if r.pseudoNameMap != nil {
					pName, ok := r.pseudoNameMap[n]
					if !ok {
						continue
					}
					newID = strpool.Intern(pName)
				}
				if newID != oldID {
					setNodeStringID(n, newID)
					r.compiler.ReportChangeToEnclosingScope(n)
					parent := n.GetParent()
					if _, isFn := ast.AsFunction(parent); isFn && NodeUtil_IsFunctionDeclaration(parent) {
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
			if a.changedFromOld {
				oldID := strpool.ID(n.GetStringID())
				newID := a.newNameID
				if r.pseudoNameMap != nil {
					pName, ok := r.pseudoNameMap[n]
					if !ok {
						continue
					}
					newID = strpool.Intern(pName)
				}
				if newID != oldID {
					setNodeStringID(n, newID)
					r.compiler.ReportChangeToEnclosingScope(n)
					parent := n.GetParent()
					if _, isFn := ast.AsFunction(parent); isFn && NodeUtil_IsFunctionDeclaration(parent) {
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
		prevNewName, ok := r.prevUsedRenameMap.origToNew[a.oldName]
		if !ok {
			continue
		}
		if _, reserved := r.reservedNames[prevNewName]; reserved {
			continue
		}

		_, isExtern := r.externNames[a.oldName]
		if a.isLocal ||
			(!isExtern && strings.HasPrefix(prevNewName, r.prefix)) {
			r.reservedNames[prevNewName] = struct{}{}
			r.FinalizeNameAssignment(a, prevNewName)
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
		if a.hasNewName {
			continue
		}

		if _, isExtern := r.externNames[a.oldName]; isExtern {
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
		r.reservedNames[newName] = struct{}{}
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
	r.finalized = append(r.finalized, renameVarsMapEntry{oldName: a.oldName, newName: newName})
}

// GetVariableMap gets the variable map.
//
// Java: RenameVars.getVariableMap()
func (r *RenameVars) GetVariableMap() *VariableMap {
	n := len(r.finalized)
	origToNew := make(map[string]string, n)
	newToOrig := make(map[string]string, n)
	origKeys := make([]string, 0, n)
	for _, entry := range r.finalized {
		if prevNew, exists := origToNew[entry.oldName]; exists {
			delete(newToOrig, prevNew)
		} else {
			origKeys = append(origKeys, entry.oldName)
		}
		if prevOrig, dup := newToOrig[entry.newName]; dup && prevOrig != entry.oldName {
			panic(jutil.NewIllegalArgumentExceptionMsg(
				"Multiple entries with same value: " + entry.oldName + "=" + entry.newName + " and " + prevOrig + "=" + entry.newName))
		}
		origToNew[entry.oldName] = entry.newName
		newToOrig[entry.newName] = entry.oldName
	}
	return &VariableMap{
		origToNew: origToNew,
		newToOrig: newToOrig,
		origKeys:  origKeys,
	}
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

// -------------------------------------------------------------------------

// RenamingMap is the interface used by to substitute names.
type RenamingMap interface {
	// Get returns the replacement for value, or nil (Java null).
	//
	// Java: RenamingMap.get(String)
	Get(value string) *string
}

// RenamingToken is a marker object for the command line api to pass
// special values to the ReplaceIdGenerators pass.
type RenamingToken int32

const (
	RenamingToken_INCONSISTENT RenamingToken = iota
	RenamingToken_STABLE
	RenamingToken_DISABLE
)

var renamingToken_names = [...]string{"INCONSISTENT", "STABLE", "DISABLE"}

// Name is Enum.name().
func (r RenamingToken) Name() string { return renamingToken_names[r] }

// String is Enum.toString().
func (r RenamingToken) String() string { return r.Name() }

// Ordinal is Enum.ordinal().
func (r RenamingToken) Ordinal() int { return int(r) }

// Java: RenamingToken.get(String)
func (r RenamingToken) Get(value string) *string {
	return nil
}

// VariableRenamingPolicy holds the policies to determine which variables
// should be renamed.
type VariableRenamingPolicy int32

const (
	// VariableRenamingPolicy_OFF renames no variables.
	VariableRenamingPolicy_OFF VariableRenamingPolicy = iota

	// VariableRenamingPolicy_LOCAL renames local variables only.
	VariableRenamingPolicy_LOCAL

	// VariableRenamingPolicy_ALL renames all variables and functions unless
	// they are exported or externed.
	VariableRenamingPolicy_ALL
)

var variableRenamingPolicy_names = [...]string{"OFF", "LOCAL", "ALL"}

// VariableRenamingPolicy_values is VariableRenamingPolicy.values().
func VariableRenamingPolicy_values() []VariableRenamingPolicy {
	return []VariableRenamingPolicy{VariableRenamingPolicy_OFF, VariableRenamingPolicy_LOCAL, VariableRenamingPolicy_ALL}
}

// Name is Enum.name().
func (p VariableRenamingPolicy) Name() string { return variableRenamingPolicy_names[p] }

// String is Enum.toString().
func (p VariableRenamingPolicy) String() string { return p.Name() }

// Ordinal is Enum.ordinal().
func (p VariableRenamingPolicy) Ordinal() int { return int(p) }

// PropertyRenamingPolicy holds the policies to determine how properties
// should be renamed.
type PropertyRenamingPolicy int32

const (
	// PropertyRenamingPolicy_OFF renames no properties.
	PropertyRenamingPolicy_OFF PropertyRenamingPolicy = iota

	// PropertyRenamingPolicy_ALL_UNQUOTED renames all properties that
	// aren't explicitly quoted and aren't externally defined (i.e. declared
	// in an externs file). This policy achieves better compaction than the
	// others.
	//
	// See RenameProperties.
	PropertyRenamingPolicy_ALL_UNQUOTED
)

var propertyRenamingPolicy_names = [...]string{"OFF", "ALL_UNQUOTED"}

// PropertyRenamingPolicy_values is PropertyRenamingPolicy.values().
func PropertyRenamingPolicy_values() []PropertyRenamingPolicy {
	return []PropertyRenamingPolicy{PropertyRenamingPolicy_OFF, PropertyRenamingPolicy_ALL_UNQUOTED}
}

// Name is Enum.name().
func (p PropertyRenamingPolicy) Name() string { return propertyRenamingPolicy_names[p] }

// String is Enum.toString().
func (p PropertyRenamingPolicy) String() string { return p.Name() }

// Ordinal is Enum.ordinal().
func (p PropertyRenamingPolicy) Ordinal() int { return int(p) }
