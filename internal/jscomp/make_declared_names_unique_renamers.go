// Port of com/google/javascript/jscomp/MakeDeclaredNamesUnique.java @ closure-compiler v20261028

/*
 * Copyright 2009 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

type contextualDecl struct {
	origID strpool.ID
	newID  strpool.ID // 0 means declared in this scope keeping original name
}

type renamerDeclMap struct {
	inlineDecls [8]contextualDecl
	inlineLen   int
	overflow    map[strpool.ID]strpool.ID
}

func (m *renamerDeclMap) isEmpty() bool {
	return m.inlineLen == 0 && len(m.overflow) == 0
}

func (m *renamerDeclMap) get(origID strpool.ID) (strpool.ID, bool) {
	for i := 0; i < m.inlineLen; i++ {
		if m.inlineDecls[i].origID == origID {
			return m.inlineDecls[i].newID, true
		}
	}
	if m.overflow != nil {
		v, ok := m.overflow[origID]
		return v, ok
	}
	return 0, false
}

func (m *renamerDeclMap) put(origID, newID strpool.ID) {
	if m.inlineLen < len(m.inlineDecls) {
		m.inlineDecls[m.inlineLen] = contextualDecl{origID: origID, newID: newID}
		m.inlineLen++
		return
	}
	if m.overflow == nil {
		m.overflow = make(map[strpool.ID]strpool.ID, 16)
	}
	m.overflow[origID] = newID
}

func (m *renamerDeclMap) lookupReplacementValue(oldName string) (string, bool) {
	if m.isEmpty() {
		return "", false
	}
	id := strpool.Lookup(oldName)
	if id == 0 {
		return "", false
	}
	if newID, ok := m.get(id); ok && newID != 0 {
		return strpool.String(newID), true
	}
	return "", false
}

func (m *renamerDeclMap) String() string {
	if m.isEmpty() {
		return "{}"
	}
	var sb strings.Builder
	sb.WriteByte('{')
	first := true
	writeEntry := func(origID, newID strpool.ID) {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(strpool.String(origID))
		sb.WriteByte('=')
		if newID == 0 {
			sb.WriteString("null")
		} else {
			sb.WriteString(strpool.String(newID))
		}
	}
	for i := 0; i < m.inlineLen; i++ {
		writeEntry(m.inlineDecls[i].origID, m.inlineDecls[i].newID)
	}
	if len(m.overflow) > 0 {
		keys := make([]strpool.ID, 0, len(m.overflow))
		// order-independent: overflow keys are sorted before formatting debug String()
		for k := range m.overflow {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		for _, k := range keys {
			writeEntry(k, m.overflow[k])
		}
	}
	sb.WriteByte('}')
	return sb.String()
}

// MakeDeclaredNamesUnique_ContextualRenamer renames every local name to be
// unique. The first encountered declaration of a given name (specifically a
// global declaration) is left in its original form. Those that are renamed
// are made unique by giving them a unique suffix based on the number of
// declarations of the name.
//
// The root ContextualRenamer is assumed to be in GlobalScope.
//
// Used by the Normalize pass.
//
// See Normalize.
type MakeDeclaredNamesUnique_ContextualRenamer struct {
	scopeRoot    rhino.Node
	nameUsage    map[strpool.ID]int
	decls        renamerDeclMap
	global       bool
	hoistRenamer MakeDeclaredNamesUnique_Renamer
}

const MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR = "$jscomp$"

func formatNameUsage(counts map[strpool.ID]int) string {
	if len(counts) == 0 {
		return "[]"
	}
	type entry struct {
		name  string
		count int
	}
	entries := make([]entry, 0, len(counts))
	// order-independent: entries are sorted lexicographically by name before formatting debug String()
	for id, c := range counts {
		entries = append(entries, entry{name: strpool.String(id), count: c})
	}
	slices.SortFunc(entries, func(a, b entry) int {
		return strings.Compare(a.name, b.name)
	})
	var sb strings.Builder
	sb.WriteByte('[')
	for idx, e := range entries {
		if idx > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(e.name)
		if e.count != 1 {
			sb.WriteString(" x ")
			sb.WriteString(strconv.Itoa(e.count))
		}
	}
	sb.WriteByte(']')
	return sb.String()
}

func (c *MakeDeclaredNamesUnique_ContextualRenamer) formatToString(className string) string {
	return className + "{scopeRoot=" + jutil.ToJavaString(c.scopeRoot) +
		", nameUsage=" + formatNameUsage(c.nameUsage) +
		", declarations=" + c.decls.String() +
		", global=" + jutil.ToJavaString(c.global) + "}"
}

// Java: MakeDeclaredNamesUnique.ContextualRenamer.toString()
func (c *MakeDeclaredNamesUnique_ContextualRenamer) String() string {
	return c.formatToString("ContextualRenamer")
}

// Java: MakeDeclaredNamesUnique.ContextualRenamer()
func NewMakeDeclaredNamesUnique_ContextualRenamer() *MakeDeclaredNamesUnique_ContextualRenamer {
	c := &MakeDeclaredNamesUnique_ContextualRenamer{
		scopeRoot: rhino.NilNode,
		global:    true,
		nameUsage: make(map[strpool.ID]int, 64),
	}
	c.hoistRenamer = c
	return c
}

// NewMakeDeclaredNamesUnique_ContextualRenamerForChildScope is the
// constructor for child scopes.
//
// Java: MakeDeclaredNamesUnique.ContextualRenamer(Node, Multiset<String>, boolean, MakeDeclaredNamesUnique.Renamer)
func NewMakeDeclaredNamesUnique_ContextualRenamerForChildScope(
	scopeRoot rhino.Node,
	nameUsage map[strpool.ID]int,
	hoistingTargetScope bool,
	parent MakeDeclaredNamesUnique_Renamer) *MakeDeclaredNamesUnique_ContextualRenamer {
	jutil.CheckStateMsg1(NodeUtil_CreatesScope(scopeRoot), "%s", scopeRoot)

	if scopeRoot.IsFunction() {
		jutil.CheckStateMsg1(!hoistingTargetScope, "%s", scopeRoot)
	}

	var hoistRenamer MakeDeclaredNamesUnique_Renamer
	if hoistingTargetScope {
		jutil.CheckStateMsg1(!NodeUtil_CreatesBlockScope(scopeRoot), "%s", scopeRoot)
	} else {
		jutil.CheckStateMsg1(NodeUtil_CreatesBlockScope(scopeRoot) || scopeRoot.IsFunction(), "%s", scopeRoot)
		hoistRenamer = parent.GetHoistRenamer()
	}

	c := &MakeDeclaredNamesUnique_ContextualRenamer{
		scopeRoot: scopeRoot,
		global:    false,
		nameUsage: nameUsage,
	}
	if hoistingTargetScope {
		c.hoistRenamer = c
	} else {
		c.hoistRenamer = hoistRenamer
	}
	return c
}

// CreateForChildScope creates a ContextualRenamer
//
// Java: MakeDeclaredNamesUnique.ContextualRenamer.createForChildScope(Node, boolean)
func (c *MakeDeclaredNamesUnique_ContextualRenamer) CreateForChildScope(
	scopeRoot rhino.Node, hoistingTargetScope bool) MakeDeclaredNamesUnique_Renamer {
	return NewMakeDeclaredNamesUnique_ContextualRenamerForChildScope(scopeRoot, c.nameUsage, hoistingTargetScope, c)
}

// AddDeclaredName adds a name to the map of names declared in this scope.
//
// Java: MakeDeclaredNamesUnique.ContextualRenamer.addDeclaredName(String, boolean)
func (c *MakeDeclaredNamesUnique_ContextualRenamer) AddDeclaredName(name string, hoisted bool) {
	c.addDeclaredNameID(strpool.Intern(name), hoisted)
}

func (c *MakeDeclaredNamesUnique_ContextualRenamer) addDeclaredNameID(origID strpool.ID, hoisted bool) {
	if hoisted && c.hoistRenamer != MakeDeclaredNamesUnique_Renamer(c) && !c.global {
		c.hoistRenamer.addDeclaredNameID(origID, true)
		return
	}
	if origID == idArguments {
		return
	}
	if c.global {
		if _, exists := c.nameUsage[origID]; !exists {
			c.nameUsage[origID] = 1
		}
		return
	}
	if _, exists := c.decls.get(origID); exists {
		return
	}
	id := c.nameUsage[origID]
	c.nameUsage[origID] = id + 1
	var newID strpool.ID
	if id != 0 {
		newID = allocateContextualUniqueNameID(c.nameUsage, origID, id, true)
	}
	c.decls.put(origID, newID)
}

func allocateContextualUniqueNameID(nameUsage map[strpool.ID]int, origID strpool.ID, id int, checkCollision bool) strpool.ID {
	name := strpool.String(origID)
	newID := strpool.Intern(MakeDeclaredNamesUnique_ContextualRenamer_GetUniqueName(name, id))
	if checkCollision {
		for {
			uid, exists := nameUsage[newID]
			if !exists || uid == 0 {
				break
			}
			nameUsage[newID] = uid + 1
			newID = strpool.Intern(MakeDeclaredNamesUnique_ContextualRenamer_GetUniqueName(name, uid))
		}
		if _, exists := nameUsage[newID]; !exists {
			nameUsage[newID] = 1
		}
	}
	return newID
}

type contextualRenameUndo struct {
	name    strpool.ID
	oldRepl strpool.ID
	hadOld  bool
}

// contextualRenamerState implements the single-pass LIFO undo-stack contextual
// renamer embedded in Normalize_NormalizeStatements, sharing $jscomp$ unique-suffix
// collision handling with MakeDeclaredNamesUnique_ContextualRenamer.
type contextualRenamerState struct {
	enabled             bool
	globalEntered       bool
	hasJscompVar        bool
	nameUsage           map[strpool.ID]int
	activeReplacements  map[strpool.ID]strpool.ID
	undoStack           []contextualRenameUndo
	scopeUndoOffsets    []int
	undoStackBuf        [8]contextualRenameUndo
	scopeUndoOffsetsBuf [8]int
}

func (r *contextualRenamerState) enable() {
	r.enabled = true
	r.undoStack = r.undoStackBuf[:0]
	r.scopeUndoOffsets = r.scopeUndoOffsetsBuf[:0]
}

func (r *contextualRenamerState) reset() {
	*r = contextualRenamerState{}
}

func contextualHasJscompSeparator(s string) bool {
	const sep = MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR
	if len(s) < len(sep)+2 {
		return false
	}
	last := s[len(s)-1]
	if last < '0' || last > '9' {
		return false
	}
	i := len(s) - 1
	for i > 0 && s[i-1] >= '0' && s[i-1] <= '9' {
		i--
	}
	if s[i] == '0' || i < len(sep)+1 {
		return false
	}
	return s[i-len(sep):i] == sep
}

func (r *contextualRenamerState) enableJscompVarMode() {
	if r.hasJscompVar {
		return
	}
	r.hasJscompVar = true
	if len(r.nameUsage) == 0 {
		return
	}
	type pending struct {
		nameID strpool.ID
		count  int
	}
	var toBackfill []pending
	// order-independent: backfills one entry per (nameID, i), and only absent ones, so the resulting map is the same in any order
	for origID, count := range r.nameUsage {
		if count > 1 {
			toBackfill = append(toBackfill, pending{nameID: origID, count: count})
		}
	}
	for _, p := range toBackfill {
		origName := strpool.String(p.nameID)
		for i := 1; i < p.count; i++ {
			u := strpool.Intern(MakeDeclaredNamesUnique_ContextualRenamer_GetUniqueName(origName, i))
			if _, exists := r.nameUsage[u]; !exists {
				r.nameUsage[u] = 1
			}
		}
	}
}

func (r *contextualRenamerState) enterScope(t *NodeTraversal) {
	if !r.enabled {
		return
	}
	declarationRoot := t.GetScopeRoot()
	if !r.globalEntered {
		jutil.CheckState(!declarationRoot.IsFunction())
		r.globalEntered = true
		varList := t.GetScope().varList
		if r.nameUsage == nil {
			r.nameUsage = make(map[strpool.ID]int, max(64, len(varList)+len(varList)/2))
		}
		for _, v := range varList {
			if v == nil || strpool.ID(v.nameID) == idArguments {
				continue
			}
			nameID := strpool.ID(v.nameID)
			if !r.hasJscompVar && contextualHasJscompSeparator(strpool.String(nameID)) {
				r.enableJscompVarMode()
			}
			if _, exists := r.nameUsage[nameID]; !exists {
				r.nameUsage[nameID] = 1
			}
		}
		return
	}
	r.scopeUndoOffsets = append(r.scopeUndoOffsets, len(r.undoStack))
	if !normalizeScopeMayHaveDeclarations(declarationRoot) {
		return
	}
	for _, v := range t.GetScope().varList {
		if v == nil || strpool.ID(v.nameID) == idArguments {
			continue
		}
		nameID := strpool.ID(v.nameID)
		id := r.nameUsage[nameID]
		if id == 0 && !r.hasJscompVar && contextualHasJscompSeparator(strpool.String(nameID)) {
			r.enableJscompVarMode()
			id = r.nameUsage[nameID]
		}
		r.nameUsage[nameID] = id + 1
		if id == 0 {
			continue
		}
		unique := allocateContextualUniqueNameID(r.nameUsage, nameID, id, r.hasJscompVar)
		oldRepl, hadOld := r.activeReplacements[nameID]
		r.undoStack = append(r.undoStack, contextualRenameUndo{
			name:    nameID,
			oldRepl: oldRepl,
			hadOld:  hadOld,
		})
		if r.activeReplacements == nil {
			r.activeReplacements = make(map[strpool.ID]strpool.ID, 16)
		}
		r.activeReplacements[nameID] = unique
	}
}

func (r *contextualRenamerState) exitScope(t *NodeTraversal) {
	if !r.enabled || t.InGlobalScope() {
		return
	}
	offset := r.scopeUndoOffsets[len(r.scopeUndoOffsets)-1]
	r.scopeUndoOffsets = r.scopeUndoOffsets[:len(r.scopeUndoOffsets)-1]
	for i := len(r.undoStack) - 1; i >= offset; i-- {
		u := r.undoStack[i]
		if u.hadOld {
			r.activeReplacements[u.name] = u.oldRepl
		} else {
			delete(r.activeReplacements, u.name)
		}
	}
	r.undoStack = r.undoStack[:offset]
}

func (r *contextualRenamerState) renameNode(t *NodeTraversal, n rhino.Node, assertOnChange bool) {
	if !r.enabled || len(r.activeReplacements) == 0 {
		return
	}
	if n.IsName() && NodeUtil_IsNonlocalModuleExportName(n) {
		return
	}
	newID, ok := r.activeReplacements[strpool.ID(n.GetStringID())]
	if !ok {
		return
	}
	if assertOnChange {
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected renaming in MakeDeclaredNamesUnique. new name: " + strpool.String(newID) + " for " + n.String()))
	}
	setNodeStringID(n, newID)
	t.ReportCodeChange()
	if parent := n.GetParent(); parent.IsFunction() && NodeUtil_IsFunctionDeclaration(parent) {
		t.GetCompiler().ReportChangeToEnclosingScope(parent)
	}
}

func (c *MakeDeclaredNamesUnique_ContextualRenamer) getReplacementNameID(origID strpool.ID) (strpool.ID, bool) {
	newID, ok := c.decls.get(origID)
	return newID, ok && newID != 0
}

// Java: MakeDeclaredNamesUnique.ContextualRenamer.getReplacementName(String)
func (c *MakeDeclaredNamesUnique_ContextualRenamer) GetReplacementNameValue(oldName string) (string, bool) {
	return c.decls.lookupReplacementValue(oldName)
}

// MakeDeclaredNamesUnique_ContextualRenamer_GetUniqueName: given a name
// and the associated id, create a new unique name.
//
// Java: MakeDeclaredNamesUnique.ContextualRenamer.getUniqueName(String, int)
func MakeDeclaredNamesUnique_ContextualRenamer_GetUniqueName(name string, id int) string {
	return name + MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR + strconv.Itoa(id)
}

// Java: MakeDeclaredNamesUnique.ContextualRenamer.stripConstIfReplaced()
func (c *MakeDeclaredNamesUnique_ContextualRenamer) StripConstIfReplaced() bool {
	return false
}

// Java: MakeDeclaredNamesUnique.ContextualRenamer.getHoistRenamer()
func (c *MakeDeclaredNamesUnique_ContextualRenamer) GetHoistRenamer() MakeDeclaredNamesUnique_Renamer {
	return c.hoistRenamer
}

// MakeDeclaredNamesUnique_InlineRenamer renames every declared name to be
// unique. Typically, this would be used when injecting code to ensure that
// names do not conflict with existing names.
//
// # Used by the FunctionInjector
//
// See FunctionInjector.
type MakeDeclaredNamesUnique_InlineRenamer struct {
	decls            renamerDeclMap
	uniqueIdSupplier func() string
	idPrefix         string
	removeConstness  bool
	convention       CodingConvention
	hoistRenamer     MakeDeclaredNamesUnique_Renamer
}

// Java: MakeDeclaredNamesUnique.InlineRenamer(CodingConvention, Supplier<String>, String, boolean, boolean, MakeDeclaredNamesUnique.Renamer)
func NewMakeDeclaredNamesUnique_InlineRenamer(
	convention CodingConvention,
	uniqueIdSupplier func() string,
	idPrefix string,
	removeConstness bool,
	hoistingTargetScope bool,
	parent MakeDeclaredNamesUnique_Renamer) *MakeDeclaredNamesUnique_InlineRenamer {
	// To ensure that the id does not conflict with the id from the
	// ContextualRenamer some prefix is needed.
	jutil.CheckArgument(idPrefix != "")
	r := &MakeDeclaredNamesUnique_InlineRenamer{
		convention:       convention,
		uniqueIdSupplier: uniqueIdSupplier,
		idPrefix:         idPrefix,
		removeConstness:  removeConstness,
	}
	if hoistingTargetScope {
		r.hoistRenamer = r
	} else {
		r.hoistRenamer = parent.GetHoistRenamer()
	}
	return r
}

func (r *MakeDeclaredNamesUnique_InlineRenamer) getReplacementNameID(origID strpool.ID) (strpool.ID, bool) {
	return r.decls.get(origID)
}

// Java: MakeDeclaredNamesUnique.InlineRenamer.addDeclaredName(String, boolean)
func (r *MakeDeclaredNamesUnique_InlineRenamer) AddDeclaredName(name string, hoisted bool) {
	r.addDeclaredNameID(strpool.Intern(name), hoisted)
}

func (r *MakeDeclaredNamesUnique_InlineRenamer) addDeclaredNameID(origID strpool.ID, hoisted bool) {
	jutil.CheckState(origID != idArguments)
	if hoisted && r.hoistRenamer != MakeDeclaredNamesUnique_Renamer(r) {
		r.hoistRenamer.addDeclaredNameID(origID, hoisted)
		return
	}
	if _, exists := r.decls.get(origID); exists {
		return
	}
	unique := r.GetUniqueName(strpool.String(origID))
	r.decls.put(origID, strpool.Intern(unique))
}

// Java: MakeDeclaredNamesUnique.InlineRenamer.getUniqueName(String)
func (r *MakeDeclaredNamesUnique_InlineRenamer) GetUniqueName(name string) string {
	if name == "" {
		return name
	}

	if idx := strings.LastIndex(name, MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR); idx >= 0 {
		name = name[:idx]
	}

	const sep = MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR
	if r.convention.IsExported(name) {
		// The google internal coding convention includes a naming convention
		// to export names starting with "_".  Simply strip "_" those to avoid
		// exporting names.
		return "JSCompiler_" + name + sep + r.idPrefix + r.uniqueIdSupplier()
	}

	// By using the same separator the id will be stripped if it isn't
	// needed when variable renaming is turned off.
	return name + sep + r.idPrefix + r.uniqueIdSupplier()
}

// Java: MakeDeclaredNamesUnique.InlineRenamer.getReplacementName(String)
func (r *MakeDeclaredNamesUnique_InlineRenamer) GetReplacementNameValue(oldName string) (string, bool) {
	return r.decls.lookupReplacementValue(oldName)
}

// Java: MakeDeclaredNamesUnique.InlineRenamer.createForChildScope(Node, boolean)
func (r *MakeDeclaredNamesUnique_InlineRenamer) CreateForChildScope(
	scopeRoot rhino.Node, hoistingTargetScope bool) MakeDeclaredNamesUnique_Renamer {
	return NewMakeDeclaredNamesUnique_InlineRenamer(
		r.convention, r.uniqueIdSupplier, r.idPrefix, r.removeConstness, hoistingTargetScope, r)
}

// Java: MakeDeclaredNamesUnique.InlineRenamer.stripConstIfReplaced()
func (r *MakeDeclaredNamesUnique_InlineRenamer) StripConstIfReplaced() bool {
	return r.removeConstness
}

// Java: MakeDeclaredNamesUnique.InlineRenamer.getHoistRenamer()
func (r *MakeDeclaredNamesUnique_InlineRenamer) GetHoistRenamer() MakeDeclaredNamesUnique_Renamer {
	return r.hoistRenamer
}

// MakeDeclaredNamesUnique_BoilerplateRenamer is for injecting boilerplate
// libraries. Leaves global names alone and renames local names like
// InlineRenamer.
type MakeDeclaredNamesUnique_BoilerplateRenamer struct {
	MakeDeclaredNamesUnique_ContextualRenamer
	uniqueIdSupplier func() string
	idPrefix         string
	convention       CodingConvention
}

// Java: MakeDeclaredNamesUnique.BoilerplateRenamer(CodingConvention, Supplier<String>, String)
func NewMakeDeclaredNamesUnique_BoilerplateRenamer(
	convention CodingConvention, uniqueIdSupplier func() string, idPrefix string) *MakeDeclaredNamesUnique_BoilerplateRenamer {
	b := &MakeDeclaredNamesUnique_BoilerplateRenamer{
		MakeDeclaredNamesUnique_ContextualRenamer: *NewMakeDeclaredNamesUnique_ContextualRenamer(),
		convention:       convention,
		uniqueIdSupplier: uniqueIdSupplier,
		idPrefix:         idPrefix,
	}
	b.hoistRenamer = b
	return b
}

// Java: MakeDeclaredNamesUnique.ContextualRenamer.toString()
func (b *MakeDeclaredNamesUnique_BoilerplateRenamer) String() string {
	return b.formatToString("BoilerplateRenamer")
}

// Java: MakeDeclaredNamesUnique.BoilerplateRenamer.createForChildScope(Node, boolean)
func (b *MakeDeclaredNamesUnique_BoilerplateRenamer) CreateForChildScope(
	scopeRoot rhino.Node, hoisted bool) MakeDeclaredNamesUnique_Renamer {
	return NewMakeDeclaredNamesUnique_InlineRenamer(b.convention, b.uniqueIdSupplier, b.idPrefix, false, hoisted, b)
}

// MakeDeclaredNamesUnique_TargettedRenamer only renames things that match
// specific names. Wraps another renamer.
type MakeDeclaredNamesUnique_TargettedRenamer struct {
	delegate   MakeDeclaredNamesUnique_Renamer
	targets    *jutil.LinkedHashSet[string]
	targetIDs  map[strpool.ID]struct{}
	cachedSize int
}

// Java: MakeDeclaredNamesUnique.TargettedRenamer(MakeDeclaredNamesUnique.Renamer, Set<String>)
func NewMakeDeclaredNamesUnique_TargettedRenamer(
	delegate MakeDeclaredNamesUnique_Renamer, targets *jutil.LinkedHashSet[string]) *MakeDeclaredNamesUnique_TargettedRenamer {
	r := &MakeDeclaredNamesUnique_TargettedRenamer{delegate: delegate, targets: targets}
	r.syncTargetIDs()
	return r
}

func (r *MakeDeclaredNamesUnique_TargettedRenamer) syncTargetIDs() {
	if r.targets == nil {
		return
	}
	n := r.targets.Size()
	if r.targetIDs != nil && r.cachedSize == n {
		return
	}
	r.targetIDs = make(map[strpool.ID]struct{}, n)
	for s := range r.targets.All() {
		r.targetIDs[strpool.Intern(s)] = struct{}{}
	}
	r.cachedSize = n
}

func (r *MakeDeclaredNamesUnique_TargettedRenamer) hasTargetID(id strpool.ID) bool {
	r.syncTargetIDs()
	_, ok := r.targetIDs[id]
	return ok
}

// Java: MakeDeclaredNamesUnique.TargettedRenamer.addDeclaredName(String, boolean)
func (r *MakeDeclaredNamesUnique_TargettedRenamer) AddDeclaredName(name string, hoisted bool) {
	if r.targets.Contains(name) {
		r.delegate.AddDeclaredName(name, hoisted)
	}
}

func (r *MakeDeclaredNamesUnique_TargettedRenamer) addDeclaredNameID(origID strpool.ID, hoisted bool) {
	if r.hasTargetID(origID) {
		r.delegate.addDeclaredNameID(origID, hoisted)
	}
}

func (r *MakeDeclaredNamesUnique_TargettedRenamer) getReplacementNameID(oldID strpool.ID) (strpool.ID, bool) {
	if r.hasTargetID(oldID) {
		return r.delegate.getReplacementNameID(oldID)
	}
	return 0, false
}

// Java: MakeDeclaredNamesUnique.TargettedRenamer.getReplacementName(String)
func (r *MakeDeclaredNamesUnique_TargettedRenamer) GetReplacementNameValue(oldName string) (string, bool) {
	if r.targets.Contains(oldName) {
		return r.delegate.GetReplacementNameValue(oldName)
	}
	return "", false
}

// Java: MakeDeclaredNamesUnique.TargettedRenamer.stripConstIfReplaced()
func (r *MakeDeclaredNamesUnique_TargettedRenamer) StripConstIfReplaced() bool {
	return r.delegate.StripConstIfReplaced()
}

// Java: MakeDeclaredNamesUnique.TargettedRenamer.createForChildScope(Node, boolean)
func (r *MakeDeclaredNamesUnique_TargettedRenamer) CreateForChildScope(
	scopeRoot rhino.Node, hoistingTargetScope bool) MakeDeclaredNamesUnique_Renamer {
	r.syncTargetIDs()
	return &MakeDeclaredNamesUnique_TargettedRenamer{
		delegate:   r.delegate.CreateForChildScope(scopeRoot, hoistingTargetScope),
		targets:    r.targets,
		targetIDs:  r.targetIDs,
		cachedSize: r.cachedSize,
	}
}

// Java: MakeDeclaredNamesUnique.TargettedRenamer.getHoistRenamer()
func (r *MakeDeclaredNamesUnique_TargettedRenamer) GetHoistRenamer() MakeDeclaredNamesUnique_Renamer {
	return r.delegate.GetHoistRenamer()
}
