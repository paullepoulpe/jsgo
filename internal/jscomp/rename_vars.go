// Port of com/google/javascript/jscomp/RenameVars.java @ closure-compiler v20261104
// Port of com/google/javascript/jscomp/RenamingMap.java @ closure-compiler v20261104

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
			t[i] = localVarTempPrefix + strconv.Itoa(i)
		}
		return t
	}()
)

type renameVarsNodeRef struct {
	node       rhino.Node
	assignment *renameVarAssignment
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
	assignmentByName map[string]*renameVarAssignment
	assignmentList   []*renameVarAssignment
	assignmentSlab   jutil.Slab[renameVarAssignment]
	localByIndex     []*renameVarAssignment
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
const maxLocalsInScopeToTempRename = 1000

// A prefix to distinguish temporary local names from global names.
const localVarTempPrefix = "L "

type renameVarAssignment struct {
	isLocal           bool
	hasNewName        bool
	changedFromOld    bool
	oldName           string
	orderOfOccurrence int
	newName           string
	newNameID         strpool.ID
	count             int // Number of times this is referenced
}

type renameVarDecision struct {
	assignment *renameVarAssignment
	invertID   strpool.ID
}

func (r *RenameVars) allocAssignment(name string, isLocal bool) *renameVarAssignment {
	a := r.assignmentSlab.New()
	a.isLocal = isLocal
	a.oldName = name
	a.orderOfOccurrence = r.assignmentCount
	r.assignmentCount++
	return a
}

func (r *RenameVars) putAssignment(name string, a *renameVarAssignment) {
	r.assignmentByName[name] = a
	r.assignmentList = append(r.assignmentList, a)
}

func appendChunkedNodeRef(chunks *[][]renameVarsNodeRef, n rhino.Node, a *renameVarAssignment) {
	s := *chunks
	if len(s) == 0 || len(s[len(s)-1]) == cap(s[len(s)-1]) {
		newCap := 1024
		if len(s) > 0 {
			newCap = min(cap(s[len(s)-1])*2, 16384)
		}
		s = append(s, make([]renameVarsNodeRef, 0, newCap))
		*chunks = s
	}
	s[len(s)-1] = append(s[len(s)-1], renameVarsNodeRef{node: n, assignment: a})
}

// Java: RenameVars.Assignment.setNewName(String)
func (a *renameVarAssignment) setNewName(newName string) {
	jutil.CheckState(!a.hasNewName)
	a.newName = newName
	a.newNameID = strpool.Intern(newName)
	a.changedFromOld = newName != a.oldName
	a.hasNewName = true
}

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
		assignmentByName:   make(map[string]*renameVarAssignment, 128),
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

// renameVarsCollector iterates through the nodes, collects all the NAME
// nodes that need to be renamed, and counts how many times each variable name
// is referenced.
//
// To get shorter local variable renaming, we rename local variables to a
// temporary name "localVarTempPrefix + index" where index is the index of the
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
type renameVarsCollector struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *RenameVars
}

func (*renameVarsCollector) visitTokens() *ast.TokenSet {
	return &renameVarsVisitTokens
}

// Java: RenameVars.ProcessVars.enterScope(NodeTraversal)
func (p *renameVarsCollector) EnterScope(t *NodeTraversal) {
	scopeRoot := t.GetScopeRoot()
	fn, ok := ast.AsFunction(scopeRoot)
	if !ok || fn.NameNode().Node().GetStringID() <= strpool.EmptyID || !NodeUtil_IsFunctionExpression(scopeRoot) {
		return
	}
	r := p.outer
	if t.InGlobalHoistScope() || !r.shouldTemporarilyRenameLocalsInScope(t.GetScope()) {
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
func (p *renameVarsCollector) ExitScope(t *NodeTraversal) {}

// Java: RenameVars.ProcessVars.visit(NodeTraversal, Node, Node)
func (p *renameVarsCollector) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	tok := n.GetToken()
	if tok != rhino.Token_NAME && tok != rhino.Token_IMPORT_STAR {
		return
	}

	id := n.GetStringID()

	// Ignore anonymous functions and classes.
	if id <= strpool.EmptyID {
		return
	}

	// "import {x as y} from 'm';"
	// Skip x because it's not a variable in this scope.
	if spec, ok := ast.AsImportSpec(parent); ok && spec.Imported().Is(n) {
		return
	}

	r := p.outer
	var_ := t.LookupVarByID(id)
	if var_ != nil {
		if dec, ok := r.seenVars[var_]; ok {
			if a := dec.assignment; a != nil {
				a.count++
				if r.pseudoNameMap != nil {
					r.recordPseudoName(n)
				}
				if a.isLocal {
					appendChunkedNodeRef(&r.localRefChunks, n, a)
				} else {
					appendChunkedNodeRef(&r.globalRefChunks, n, a)
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

		if !r.okToRenameVar(name, local) {
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
			r.recordPseudoName(n)
		}

		var a *renameVarAssignment
		if local && r.shouldTemporarilyRenameLocalsInScope(varScope) {
			localIdx := r.getLocalVarIndex(var_)
			if localIdx >= 0 && localIdx < len(r.localByIndex) {
				a = r.localByIndex[localIdx]
			}
			if a == nil {
				var tempName string
				if localIdx >= 0 && localIdx < len(renameVarsLocalTempNames) {
					tempName = renameVarsLocalTempNames[localIdx]
				} else {
					tempName = localVarTempPrefix + strconv.Itoa(localIdx)
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
						grown := make([]*renameVarAssignment, localIdx+1, newCap)
						copy(grown, r.localByIndex)
						r.localByIndex = grown
					} else if localIdx >= len(r.localByIndex) {
						r.localByIndex = r.localByIndex[:localIdx+1]
					}
					r.localByIndex[localIdx] = a
				}
			}
			a.count++
			appendChunkedNodeRef(&r.localRefChunks, n, a)
		} else {
			a = r.assignmentByName[name]
			if a == nil {
				a = r.allocAssignment(name, false)
				r.putAssignment(name, a)
			}
			a.count++
			appendChunkedNodeRef(&r.globalRefChunks, n, a)
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
	if !r.okToRenameVar(name, false) {
		return
	}
	if r.pseudoNameMap != nil {
		r.recordPseudoName(n)
	}
}

// compareAssignmentFrequency sorts renameVarAssignment objects by their count,
// breaking ties by their order of occurrence in the source to ensure a
// deterministic total ordering.
//
// Java: RenameVars.frequencyComparator(Assignment, Assignment)
func compareAssignmentFrequency(a1, a2 *renameVarAssignment) int {
	if a1.count != a2.count {
		return a2.count - a1.count
	}
	return compareAssignmentOccurrence(a1, a2)
}

// compareAssignmentOccurrence sorts renameVarAssignment objects by the order
// the variable name first appears in the source.
func compareAssignmentOccurrence(a1, a2 *renameVarAssignment) int {
	return cmp.Compare(a1.orderOfOccurrence, a2.orderOfOccurrence)
}

func (r *RenameVars) applyNodeRename(n rhino.Node, newID, oldID strpool.ID) {
	if r.pseudoNameMap != nil {
		pName, ok := r.pseudoNameMap[n]
		if !ok {
			return
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

// Java: RenameVars.process(Node, Node)
func (r *RenameVars) Process(externs rhino.Node, root rhino.Node) {
	r.externNames = make(map[string]struct{}, 64)
	collectExternVariableNamesInto(r.compiler, externs, r.externNames, r.reservedNames)

	r.globalRefChunks = nil
	r.localRefChunks = nil

	// Do variable reference counting.
	sc := NewSyntacticScopeCreator(r.compiler)
	b := NodeTraversal_Builder{
		compiler:     r.compiler,
		callback:     &renameVarsCollector{outer: r},
		scopeCreator: sc,
	}
	b.Traverse(root)
	r.localBleedingFunctions = nil
	r.localBleedingFunctionsPerScope = nil

	// Rename vars, sorted by frequency of occurrence to minimize code size.
	varsByFrequency := make([]*renameVarAssignment, len(r.assignmentList))
	copy(varsByFrequency, r.assignmentList)
	slices.SortFunc(varsByFrequency, compareAssignmentFrequency)

	// First try to reuse names from an earlier compilation.
	if r.prevUsedRenameMap != nil {
		r.reusePreviouslyUsedVariableMapSlice(varsByFrequency)
	}

	// Assign names, sorted by descending frequency to minimize code size.
	r.assignNamesSlice(varsByFrequency)

	// Rename the globals!
	for _, chunk := range r.globalRefChunks {
		for _, ref := range chunk {
			a := ref.assignment
			oldID := strpool.ID(ref.node.GetStringID())
			if a.hasNewName && a.newNameID != oldID {
				r.applyNodeRename(ref.node, a.newNameID, oldID)
			}
		}
	}

	// Rename the locals!
	for _, chunk := range r.localRefChunks {
		for _, ref := range chunk {
			a := ref.assignment
			if a.changedFromOld {
				r.applyNodeRename(ref.node, a.newNameID, strpool.ID(ref.node.GetStringID()))
			}
		}
	}

	r.globalRefChunks = nil
	r.localRefChunks = nil
	r.seenVars = nil
}

// Java: RenameVars.recordPseudoName(Node)
func (r *RenameVars) recordPseudoName(n rhino.Node) {
	// Variable names should be in a different name space than
	// property pseudo names.
	r.pseudoNameMap[n] = "$" + n.GetString() + "$$"
}

func (r *RenameVars) reusePreviouslyUsedVariableMapSlice(varsToRename []*renameVarAssignment) {
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
			r.finalizeNameAssignment(a, prevNewName)
		}
	}
}

func (r *RenameVars) assignNamesSlice(varsToRename []*renameVarAssignment) {
	r.nameGenerator.Reset(r.reservedNames, r.prefix, r.reservedCharacters)

	globalNameGenerator := r.nameGenerator
	var localNameGenerator NameGenerator
	if r.prefix == "" {
		localNameGenerator = globalNameGenerator
	} else {
		localNameGenerator = r.nameGenerator.Clone(r.reservedNames, "", r.reservedCharacters)
	}

	var pendingAssignments []*renameVarAssignment
	var generatedNamesForAssignments []string

	if fastFault(r.compiler.GetOptions(), faultBenignRename) {
		r.reservedNames[globalNameGenerator.GenerateNextName()] = struct{}{}
	}

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
			r.finalizeNameAssignment(a, newName)
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
		slices.SortFunc(batch, compareAssignmentOccurrence)
		for _, a := range batch {
			r.finalizeNameAssignment(a, generatedNamesForAssignments[i])
			i++
		}
	}
}

// finalizeNameAssignment makes a final name assignment.
//
// Java: RenameVars.finalizeNameAssignment(Assignment, String)
func (r *RenameVars) finalizeNameAssignment(a *renameVarAssignment, newName string) {
	a.setNewName(newName)

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

// okToRenameVar determines whether a variable name is okay to rename.
//
// Java: RenameVars.okToRenameVar(String, boolean)
func (r *RenameVars) okToRenameVar(name string, isLocal bool) bool {
	return !r.convention.IsExportedWithLocal(name, isLocal)
}

// getLocalVarIndex returns the index within the scope stack. e.g. function
// Foo(a) { var b; function c(d) { } } a = 0, b = 1, c = 2, d = 3
//
// Java: RenameVars.getLocalVarIndex(Var)
func (r *RenameVars) getLocalVarIndex(v *Var) int {
	num := v.GetIndex()
	s := v.scope.GetParent()
	if s == nil {
		panic(jutil.NewIllegalArgumentExceptionMsg("Var is not local"))
	}

	if r.localBleedingFunctions == nil {
		for s.GetParent() != nil {
			if r.shouldTemporarilyRenameLocalsInScope(s) {
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
		if r.shouldTemporarilyRenameLocalsInScope(s) {
			num += s.GetVarCount()
		}
		s = s.GetParent()
	}
	return num
}

// shouldTemporarilyRenameLocalsInScope returns true if the local variables
// in a scope should be given temporary names (eg, 'L 123') prior to renaming
// to allow reuse of names across scopes. With preferStableNames, temporary
// renaming is disabled if the number of locals in the scope is above a
// heuristic threshold to allow effective reuse of rename maps (see
// prevUsedRenameMap). In scopes with many variables the temporary name given
// to a variable is unlikely to be the same temporary name used when the
// rename map was created.
//
// Java: RenameVars.shouldTemporarilyRenameLocalsInScope(Scope)
func (r *RenameVars) shouldTemporarilyRenameLocalsInScope(s *Scope) bool {
	return !r.preferStableNames || s.GetVarCount() <= maxLocalsInScopeToTempRename
}

// -------------------------------------------------------------------------

// RenamingMap is the interface used by to substitute names.
type RenamingMap interface {
	// Get returns the replacement for value and true, or ("", false) when unmapped.
	//
	// Java: RenamingMap.get(String)
	Get(value string) (string, bool)
}

// RenamingToken is a marker object for the command line api to pass
// special values to the ReplaceIdGenerators pass.
type RenamingToken int32

const (
	RenamingToken_INCONSISTENT RenamingToken = iota
	RenamingToken_STABLE
	RenamingToken_DISABLE
)

var renamingTokenNames = [...]string{"INCONSISTENT", "STABLE", "DISABLE"}

// Name is Enum.name().
func (r RenamingToken) Name() string { return renamingTokenNames[r] }

// String is Enum.toString().
func (r RenamingToken) String() string { return r.Name() }

// Java: RenamingToken.get(String)
func (r RenamingToken) Get(value string) (string, bool) {
	return "", false
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

var variableRenamingPolicyNames = [...]string{"OFF", "LOCAL", "ALL"}

// VariableRenamingPolicy_values is VariableRenamingPolicy.values().
func VariableRenamingPolicy_values() []VariableRenamingPolicy {
	return []VariableRenamingPolicy{VariableRenamingPolicy_OFF, VariableRenamingPolicy_LOCAL, VariableRenamingPolicy_ALL}
}

// Name is Enum.name().
func (p VariableRenamingPolicy) Name() string { return variableRenamingPolicyNames[p] }

// String is Enum.toString().
func (p VariableRenamingPolicy) String() string { return p.Name() }

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

var propertyRenamingPolicyNames = [...]string{"OFF", "ALL_UNQUOTED"}

// PropertyRenamingPolicy_values is PropertyRenamingPolicy.values().
func PropertyRenamingPolicy_values() []PropertyRenamingPolicy {
	return []PropertyRenamingPolicy{PropertyRenamingPolicy_OFF, PropertyRenamingPolicy_ALL_UNQUOTED}
}

// Name is Enum.name().
func (p PropertyRenamingPolicy) Name() string { return propertyRenamingPolicyNames[p] }

// String is Enum.toString().
func (p PropertyRenamingPolicy) String() string { return p.Name() }
