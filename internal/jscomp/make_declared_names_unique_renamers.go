// Port of com/google/javascript/jscomp/MakeDeclaredNamesUnique.java @ closure-compiler v20261021

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
	scopeRoot rhino.Node

	// This multiset is shared between this ContextualRenamer and its parent (and its parent's
	// parent, etc.) because it tracks counts of variables across the entire JS program.
	nameUsage *makeDeclaredNamesUniqueMultiset

	// Inline small-scope declaration array + overflow map keyed by strpool.ID.
	inlineDecls [8]contextualDecl
	inlineLen   int
	overflow    map[strpool.ID]strpool.ID
	global      bool

	hoistRenamer MakeDeclaredNamesUnique_Renamer
}

const MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR = "$jscomp$"

func (c *MakeDeclaredNamesUnique_ContextualRenamer) declarationsString() string {
	if c.inlineLen == 0 && len(c.overflow) == 0 {
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
	for i := 0; i < c.inlineLen; i++ {
		writeEntry(c.inlineDecls[i].origID, c.inlineDecls[i].newID)
	}
	if len(c.overflow) > 0 {
		keys := make([]strpool.ID, 0, len(c.overflow))
		// order-independent: overflow keys are sorted before formatting debug String()
		for k := range c.overflow {
			keys = append(keys, k)
		}
		slices.Sort(keys)
		for _, k := range keys {
			writeEntry(k, c.overflow[k])
		}
	}
	sb.WriteByte('}')
	return sb.String()
}

func (c *MakeDeclaredNamesUnique_ContextualRenamer) formatToString(className string) string {
	return className + "{scopeRoot=" + jutil.ToJavaString(c.scopeRoot) +
		", nameUsage=" + c.nameUsage.String() +
		", declarations=" + c.declarationsString() +
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
		nameUsage: newMakeDeclaredNamesUniqueMultiset(),
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
	nameUsage *makeDeclaredNamesUniqueMultiset,
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

func (c *MakeDeclaredNamesUnique_ContextualRenamer) containsDeclaration(origID strpool.ID) bool {
	for i := 0; i < c.inlineLen; i++ {
		if c.inlineDecls[i].origID == origID {
			return true
		}
	}
	if c.overflow != nil {
		_, ok := c.overflow[origID]
		return ok
	}
	return false
}

func (c *MakeDeclaredNamesUnique_ContextualRenamer) putDeclaration(origID, newID strpool.ID) {
	if c.inlineLen < len(c.inlineDecls) {
		c.inlineDecls[c.inlineLen] = contextualDecl{origID: origID, newID: newID}
		c.inlineLen++
		return
	}
	if c.overflow == nil {
		c.overflow = make(map[strpool.ID]strpool.ID, 16)
	}
	c.overflow[origID] = newID
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
	if origID == makeDeclaredNamesUniqueArgumentsID {
		return
	}
	if c.global {
		c.nameUsage.SetCountConditionally(origID, 0, 1)
		return
	}
	// It hasn't been declared locally yet, so increment the count.
	if c.containsDeclaration(origID) {
		return
	}
	id := c.nameUsage.AddOccurrences(origID, 1)
	var newID strpool.ID
	if id != 0 {
		// this name exists at another source location, create another
		name := strpool.String(origID)
		newID = strpool.Intern(MakeDeclaredNamesUnique_ContextualRenamer_GetUniqueName(name, id))
		for c.nameUsage.Contains(newID) {
			// if the newName also exists at another location, create another
			id = c.nameUsage.AddOccurrences(newID, 1)
			newID = strpool.Intern(MakeDeclaredNamesUnique_ContextualRenamer_GetUniqueName(name, id))
		}
		c.nameUsage.SetCountConditionally(newID, 0, 1) // reserve this new name so that it's never reused
	}
	c.putDeclaration(origID, newID)
}

func (c *MakeDeclaredNamesUnique_ContextualRenamer) getReplacementNameID(origID strpool.ID) (strpool.ID, bool) {
	for i := 0; i < c.inlineLen; i++ {
		if c.inlineDecls[i].origID == origID {
			newID := c.inlineDecls[i].newID
			return newID, newID != 0
		}
	}
	if c.overflow != nil {
		if newID, ok := c.overflow[origID]; ok {
			return newID, newID != 0
		}
	}
	return 0, false
}

// Java: MakeDeclaredNamesUnique.ContextualRenamer.getReplacementName(String)
func (c *MakeDeclaredNamesUnique_ContextualRenamer) GetReplacementName(oldName string) *string {
	if v, ok := c.GetReplacementNameValue(oldName); ok {
		return &v
	}
	return nil
}

func (c *MakeDeclaredNamesUnique_ContextualRenamer) GetReplacementNameValue(oldName string) (string, bool) {
	if c.inlineLen == 0 && len(c.overflow) == 0 {
		return "", false
	}
	id := strpool.Lookup(oldName)
	if id == 0 {
		return "", false
	}
	if newID, ok := c.getReplacementNameID(id); ok {
		return strpool.String(newID), true
	}
	return "", false
}

// MakeDeclaredNamesUnique_ContextualRenamer_GetUniqueName: given a name
// and the associated id, create a new unique name.
//
// Java: MakeDeclaredNamesUnique.ContextualRenamer.getUniqueName(String, int)
func MakeDeclaredNamesUnique_ContextualRenamer_GetUniqueName(name string, id int) string {
	return name + MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR + strconv.Itoa(id)
}

// Java: MakeDeclaredNamesUnique.ContextualRenamer.reserveName(String)
func (c *MakeDeclaredNamesUnique_ContextualRenamer) ReserveName(name string) {
	c.nameUsage.SetCountConditionally(strpool.Intern(name), 0, 1)
}

// Java: MakeDeclaredNamesUnique.ContextualRenamer.incrementNameCount(String)
func (c *MakeDeclaredNamesUnique_ContextualRenamer) IncrementNameCount(name string) int {
	return c.nameUsage.AddOccurrences(strpool.Intern(name), 1)
}

// Java: MakeDeclaredNamesUnique.ContextualRenamer.stripConstIfReplaced()
func (c *MakeDeclaredNamesUnique_ContextualRenamer) StripConstIfReplaced() bool {
	return false
}

// Java: MakeDeclaredNamesUnique.ContextualRenamer.getHoistRenamer()
func (c *MakeDeclaredNamesUnique_ContextualRenamer) GetHoistRenamer() MakeDeclaredNamesUnique_Renamer {
	return c.hoistRenamer
}

// makeDeclaredNamesUniqueMultiset is the HashMultiset<String> of
// ContextualRenamer.nameUsage: count(e) (via Contains), add(e, n) and
// setCount(e, oldCount, newCount), keyed by strpool.ID.
type makeDeclaredNamesUniqueMultiset struct {
	counts map[strpool.ID]int
}

func newMakeDeclaredNamesUniqueMultiset() *makeDeclaredNamesUniqueMultiset {
	return &makeDeclaredNamesUniqueMultiset{counts: make(map[strpool.ID]int, 64)}
}

// Contains is Multiset.contains(Object).
func (m *makeDeclaredNamesUniqueMultiset) Contains(e strpool.ID) bool {
	_, ok := m.counts[e]
	return ok
}

// AddOccurrences is Multiset.add(E, int): returns the count before the
// operation.
func (m *makeDeclaredNamesUniqueMultiset) AddOccurrences(e strpool.ID, occurrences int) int {
	jutil.CheckArgumentMsg1(occurrences >= 0, "occurrences cannot be negative: %s", occurrences)
	old := m.counts[e]
	if occurrences == 0 {
		return old
	}
	m.counts[e] = old + occurrences
	return old
}

// SetCountConditionally is Multiset.setCount(E, int, int): sets the count
// to newCount only if it is currently oldCount.
func (m *makeDeclaredNamesUniqueMultiset) SetCountConditionally(e strpool.ID, oldCount int, newCount int) bool {
	if m.counts[e] != oldCount {
		return false
	}
	if newCount == 0 {
		delete(m.counts, e)
	} else {
		m.counts[e] = newCount
	}
	return true
}

// String is AbstractMultiset.toString(): entrySet().toString(), each entry
// "e" or "e x n".
func (m *makeDeclaredNamesUniqueMultiset) String() string {
	if len(m.counts) == 0 {
		return "[]"
	}
	type entry struct {
		name  string
		count int
	}
	entries := make([]entry, 0, len(m.counts))
	// order-independent: entries are sorted lexicographically by name before formatting debug String()
	for id, c := range m.counts {
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

// MakeDeclaredNamesUnique_InlineRenamer renames every declared name to be
// unique. Typically, this would be used when injecting code to ensure that
// names do not conflict with existing names.
//
// # Used by the FunctionInjector
//
// See FunctionInjector.
type MakeDeclaredNamesUnique_InlineRenamer struct {
	inlineDecls      [8]contextualDecl
	inlineLen        int
	overflow         map[strpool.ID]strpool.ID
	uniqueIdSupplier func() string
	idPrefix         string
	removeConstness  bool
	convention       CodingConvention

	hoistRenamer MakeDeclaredNamesUnique_Renamer
}

// Java: MakeDeclaredNamesUnique.InlineRenamer(CodingConvention, Supplier<String>, String, boolean, boolean, MakeDeclaredNamesUnique.Renamer)
func NewMakeDeclaredNamesUnique_InlineRenamer(
	convention CodingConvention,
	uniqueIdSupplier func() string,
	idPrefix string,
	removeConstness bool,
	hoistingTargetScope bool,
	parent MakeDeclaredNamesUnique_Renamer) *MakeDeclaredNamesUnique_InlineRenamer {
	r := &MakeDeclaredNamesUnique_InlineRenamer{}
	r.convention = convention
	r.uniqueIdSupplier = uniqueIdSupplier
	// To ensure that the id does not conflict with the id from the
	// ContextualRenamer some prefix is needed.
	jutil.CheckArgument(idPrefix != "")
	r.idPrefix = idPrefix
	r.removeConstness = removeConstness

	if hoistingTargetScope {
		r.hoistRenamer = r
	} else {
		r.hoistRenamer = parent.GetHoistRenamer()
	}
	return r
}

func (r *MakeDeclaredNamesUnique_InlineRenamer) getReplacementNameID(origID strpool.ID) (strpool.ID, bool) {
	for i := 0; i < r.inlineLen; i++ {
		if r.inlineDecls[i].origID == origID {
			return r.inlineDecls[i].newID, true
		}
	}
	if r.overflow != nil {
		v, ok := r.overflow[origID]
		return v, ok
	}
	return 0, false
}

// Java: MakeDeclaredNamesUnique.InlineRenamer.addDeclaredName(String, boolean)
func (r *MakeDeclaredNamesUnique_InlineRenamer) AddDeclaredName(name string, hoisted bool) {
	r.addDeclaredNameID(strpool.Intern(name), hoisted)
}

func (r *MakeDeclaredNamesUnique_InlineRenamer) addDeclaredNameID(origID strpool.ID, hoisted bool) {
	jutil.CheckState(origID != makeDeclaredNamesUniqueArgumentsID)
	if hoisted && r.hoistRenamer != MakeDeclaredNamesUnique_Renamer(r) {
		r.hoistRenamer.addDeclaredNameID(origID, hoisted)
		return
	}
	if _, exists := r.getReplacementNameID(origID); exists {
		return
	}
	unique := r.GetUniqueName(strpool.String(origID))
	newID := strpool.Intern(unique)
	if r.inlineLen < len(r.inlineDecls) {
		r.inlineDecls[r.inlineLen] = contextualDecl{origID: origID, newID: newID}
		r.inlineLen++
		return
	}
	if r.overflow == nil {
		r.overflow = make(map[strpool.ID]strpool.ID, 16)
	}
	r.overflow[origID] = newID
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
func (r *MakeDeclaredNamesUnique_InlineRenamer) GetReplacementName(oldName string) *string {
	if v, ok := r.GetReplacementNameValue(oldName); ok {
		return &v
	}
	return nil
}

func (r *MakeDeclaredNamesUnique_InlineRenamer) GetReplacementNameValue(oldName string) (string, bool) {
	if r.inlineLen == 0 && len(r.overflow) == 0 {
		return "", false
	}
	id := strpool.Lookup(oldName)
	if id == 0 {
		return "", false
	}
	if newID, ok := r.getReplacementNameID(id); ok {
		return strpool.String(newID), true
	}
	return "", false
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
func (r *MakeDeclaredNamesUnique_TargettedRenamer) GetReplacementName(oldName string) *string {
	if r.targets.Contains(oldName) {
		return r.delegate.GetReplacementName(oldName)
	}
	return nil
}

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
