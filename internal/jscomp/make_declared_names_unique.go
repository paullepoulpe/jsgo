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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

var makeDeclaredNamesUniqueArgumentsID = strpool.Intern(MakeDeclaredNamesUnique_ARGUMENTS)

// MakeDeclaredNamesUnique finds all Functions, VARs, and Exception names
// and make them unique. Specifically, it will not modify object properties.
type MakeDeclaredNamesUnique struct {
	NodeTraversal_AbstractScopedCallbackBase

	// There is one renamer on the stack for each scope. This was added before any support for ES6 was
	// in place, so it was necessary to maintain this separate stack, rather than just using the
	// NodeTraversal's stack of scopes, in order to handle catch blocks and function names correctly.
	//
	// Go-only: language-adapter — the ArrayDeque used as a stack is a slice
	// whose top is the last element.
	renamerStack    []MakeDeclaredNamesUnique_Renamer
	renamerStackBuf [16]MakeDeclaredNamesUnique_Renamer
	rootRenamer     MakeDeclaredNamesUnique_Renamer
	markChanges     bool
	assertOnChange  bool
}

func (*MakeDeclaredNamesUnique) visitTokens() *ast.TokenSet {
	return &renameInverterVisitTokens
}

// Arguments is special cased to handle cases where a local name shadows
// the arguments declaration.
const MakeDeclaredNamesUnique_ARGUMENTS = "arguments"

// Java: MakeDeclaredNamesUnique(MakeDeclaredNamesUnique.Renamer, boolean, boolean)
func NewMakeDeclaredNamesUnique(
	renamer MakeDeclaredNamesUnique_Renamer, markChanges bool, assertOnChange bool) *MakeDeclaredNamesUnique {
	m := &MakeDeclaredNamesUnique{
		rootRenamer:    renamer,
		markChanges:    markChanges,
		assertOnChange: assertOnChange,
	}
	m.renamerStack = m.renamerStackBuf[:0]
	return m
}

// MakeDeclaredNamesUnique_Builder is MakeDeclaredNamesUnique.Builder.
type MakeDeclaredNamesUnique_Builder struct {
	renamer        MakeDeclaredNamesUnique_Renamer
	markChanges    bool
	assertOnChange bool
}

// Java: MakeDeclaredNamesUnique.Builder()
func NewMakeDeclaredNamesUnique_Builder() *MakeDeclaredNamesUnique_Builder {
	return &MakeDeclaredNamesUnique_Builder{markChanges: true, assertOnChange: false}
}

// Java: MakeDeclaredNamesUnique.Builder.withRenamer(MakeDeclaredNamesUnique.Renamer)
func (b *MakeDeclaredNamesUnique_Builder) WithRenamer(renamer MakeDeclaredNamesUnique_Renamer) *MakeDeclaredNamesUnique_Builder {
	b.renamer = renamer
	return b
}

// WithMarkChanges enables the option to report any code changes to the
// compiler object.
//
// This option is true by default.
//
// Java: MakeDeclaredNamesUnique.Builder.withMarkChanges(boolean)
func (b *MakeDeclaredNamesUnique_Builder) WithMarkChanges(markChanges bool) *MakeDeclaredNamesUnique_Builder {
	b.markChanges = markChanges
	return b
}

// WithAssertOnChange: if this class finds any names that are not unique,
// throws an exception
//
// This is intended for use in validating that an AST is already normalized.
//
// This option is false by default.
//
// Java: MakeDeclaredNamesUnique.Builder.withAssertOnChange(boolean)
func (b *MakeDeclaredNamesUnique_Builder) WithAssertOnChange(assertOnChange bool) *MakeDeclaredNamesUnique_Builder {
	b.assertOnChange = assertOnChange
	return b
}

// Java: MakeDeclaredNamesUnique.Builder.build()
func (b *MakeDeclaredNamesUnique_Builder) Build() *MakeDeclaredNamesUnique {
	if b.renamer == nil {
		b.renamer = NewMakeDeclaredNamesUnique_ContextualRenamer()
	}
	return NewMakeDeclaredNamesUnique(b.renamer, b.markChanges, b.assertOnChange)
}

// Java: MakeDeclaredNamesUnique.builder()
func MakeDeclaredNamesUnique_Builder_() *MakeDeclaredNamesUnique_Builder {
	return NewMakeDeclaredNamesUnique_Builder()
}

// Java: MakeDeclaredNamesUnique.getContextualRenameInverter(AbstractCompiler)
func MakeDeclaredNamesUnique_GetContextualRenameInverter(compiler AbstractCompiler) CompilerPass {
	return NewMakeDeclaredNamesUnique_ContextualRenameInverter(compiler)
}

// Java: MakeDeclaredNamesUnique.enterScope(NodeTraversal)
func (m *MakeDeclaredNamesUnique) EnterScope(t *NodeTraversal) {
	declarationRoot := t.GetScopeRoot()

	var renamer MakeDeclaredNamesUnique_Renamer
	if len(m.renamerStack) == 0 {
		// If the contextual renamer is being used, the starting context can not
		// be a function.
		var isContextualRenamer bool
		switch m.rootRenamer.(type) {
		case *MakeDeclaredNamesUnique_ContextualRenamer, *MakeDeclaredNamesUnique_BoilerplateRenamer:
			isContextualRenamer = true
		}
		jutil.CheckState(!declarationRoot.IsFunction() || !isContextualRenamer)
		renamer = m.rootRenamer
	} else {
		hoist := !declarationRoot.IsFunction() && !NodeUtil_CreatesBlockScope(declarationRoot)
		renamer = m.renamerStack[len(m.renamerStack)-1].CreateForChildScope(declarationRoot, hoist)
	}

	m.renamerStack = append(m.renamerStack, renamer)

	m.FindDeclaredNames(t, declarationRoot)
}

// Java: MakeDeclaredNamesUnique.exitScope(NodeTraversal)
func (m *MakeDeclaredNamesUnique) ExitScope(t *NodeTraversal) {
	if !t.InGlobalScope() {
		top := len(m.renamerStack) - 1
		if top < 0 {
			panic(jutil.NewNoSuchElementException()) // ArrayDeque.pop()
		}
		m.renamerStack[top] = nil
		m.renamerStack = m.renamerStack[:top]
	}
}

// Java: MakeDeclaredNamesUnique.visit(NodeTraversal, Node, Node)
func (m *MakeDeclaredNamesUnique) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_NAME, rhino.Token_IMPORT_STAR:
		m.VisitNameOrImportStar(t, n, parent)
	default:
	}
}

// Java: MakeDeclaredNamesUnique.visitNameOrImportStar(NodeTraversal, Node, Node)
func (m *MakeDeclaredNamesUnique) VisitNameOrImportStar(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	// Don't rename the exported name foo in export {a as foo}; or import {foo as b};
	if n.IsName() && NodeUtil_IsNonlocalModuleExportName(n) {
		return
	}
	oldID := strpool.ID(n.GetStringID())
	newID, ok := m.getReplacementNameID(oldID)
	if !ok {
		return
	}
	if m.assertOnChange {
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected renaming in MakeDeclaredNamesUnique. new name: " + strpool.String(newID) + " for " + n.String()))
	}

	renamer := m.renamerStack[len(m.renamerStack)-1]
	if renamer.StripConstIfReplaced() {
		n.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, false)
		jsDocInfoNode := NodeUtil_GetBestJSDocInfoNode(n)
		if !jsDocInfoNode.IsNil() && jsDocInfoNode.GetJSDocInfo() != nil {
			builder := rhino.JSDocInfo_Builder_CopyFrom(jsDocInfoNode.GetJSDocInfo())
			builder.RecordMutable()
			jsDocInfoNode.SetJSDocInfo(builder.Build())
		}
	}
	setNodeStringID(n, newID)
	if m.markChanges {
		t.ReportCodeChange()
		// If we are renaming a function declaration, make sure the containing scope
		// has the opportunity to act on the change.
		if _, isFn := ast.AsFunction(parent); isFn && NodeUtil_IsFunctionDeclaration(parent) {
			t.GetCompiler().ReportChangeToEnclosingScope(parent)
		}
	}
}

// GetReplacementName walks the stack of name maps and finds the replacement
// name for the current scope.
//
// Java: MakeDeclaredNamesUnique.getReplacementName(String)
func (m *MakeDeclaredNamesUnique) GetReplacementName(oldName string) *string {
	if v, ok := m.GetReplacementNameValue(oldName); ok {
		return &v
	}
	return nil
}

// GetReplacementNameValue walks the stack of name maps and returns the replacement
// name without heap-allocating a *string wrapper.
func (m *MakeDeclaredNamesUnique) GetReplacementNameValue(oldName string) (string, bool) {
	id := strpool.Lookup(oldName)
	if id == 0 {
		return "", false
	}
	if newID, ok := m.getReplacementNameID(id); ok {
		return strpool.String(newID), true
	}
	return "", false
}

func (m *MakeDeclaredNamesUnique) getReplacementNameID(oldID strpool.ID) (strpool.ID, bool) {
	for i := len(m.renamerStack) - 1; i >= 0; i-- {
		if newID, ok := m.renamerStack[i].getReplacementNameID(oldID); ok {
			return newID, true
		}
	}
	return 0, false
}

// FindDeclaredNames traverses the current scope and collects declared
// names by calling addDeclaredName on the Renamer that is at the top of the
// renamerStack.
//
// Java: MakeDeclaredNamesUnique.findDeclaredNames(NodeTraversal, Node)
func (m *MakeDeclaredNamesUnique) FindDeclaredNames(t *NodeTraversal, n rhino.Node) {
	jutil.CheckStateMsg1(NodeUtil_CreatesScope(n) || n.IsScript(), "%s", n)

	topRenamer := m.renamerStack[len(m.renamerStack)-1]
	for _, v := range t.GetScope().varList {
		if v != nil {
			origID := strpool.ID(v.nameID)
			if origID == 0 {
				origID = strpool.Intern(v.GetName())
			}
			topRenamer.addDeclaredNameID(origID, false)
		}
	}
}

// MakeDeclaredNamesUnique_Renamer is the declared names renaming policy
// interface.
type MakeDeclaredNamesUnique_Renamer interface {
	// AddDeclaredName is called when a declared name is found in the local
	// current scope.
	//
	// hoisted is whether this name should be declared in the nearest
	// enclosing "hoist scope" instead of the scope represented by this
	// Renamer.
	//
	// Java: MakeDeclaredNamesUnique.Renamer.addDeclaredName(String, boolean)
	AddDeclaredName(name string, hoisted bool)

	addDeclaredNameID(origID strpool.ID, hoisted bool)

	// GetReplacementName returns a replacement name, null if oldName is
	// unknown or should not be replaced.
	//
	// Java: MakeDeclaredNamesUnique.Renamer.getReplacementName(String)
	GetReplacementName(oldName string) *string

	GetReplacementNameValue(oldName string) (string, bool)

	getReplacementNameID(oldID strpool.ID) (strpool.ID, bool)

	// StripConstIfReplaced returns whether the constant-ness of a name
	// should be removed.
	//
	// Java: MakeDeclaredNamesUnique.Renamer.stripConstIfReplaced()
	StripConstIfReplaced() bool

	// CreateForChildScope: hoisted is true if this is a "hoist" scope: A
	// function, module, or global scope. Returns a Renamer for a scope
	// within the scope of the current Renamer.
	//
	// Java: MakeDeclaredNamesUnique.Renamer.createForChildScope(Node, boolean)
	CreateForChildScope(scopeRoot rhino.Node, hoisted bool) MakeDeclaredNamesUnique_Renamer

	// GetHoistRenamer returns the closest hoisting target for var and
	// function declarations.
	//
	// Java: MakeDeclaredNamesUnique.Renamer.getHoistRenamer()
	GetHoistRenamer() MakeDeclaredNamesUnique_Renamer
}

// MakeDeclaredNamesUnique_ContextualRenameInverter inverts the
// transformation by ContextualRenamer, when possible.
//
// Go-only: performance — collapses Java's 5-class VariableInfo /
// DeclaredVariableInfo / UndeclaredVariableInfo / RenamingInfo /
// VariableInfoStackMap hierarchy onto a slab-allocated concrete inverterVar
// with an intrusive prevShadow pointer stack keyed by StringPool ID, skips
// non-NAME/IMPORT_STAR Visit callbacks via visitTokens(), and skips top-level
// local scope subtrees that contain no "$jscomp$" names.
type MakeDeclaredNamesUnique_ContextualRenameInverter struct {
	compiler AbstractCompiler

	// Stack of active local scope contexts (empty in the global scope).
	scopeContextStack []*renameInverterScopeContext

	// Maps StringPool ID to the innermost visible inverterVar with that name;
	// outer shadowed variables are chained via inverterVar.prevShadow.
	shadowByID map[uint32]*inverterVar

	varSlab            jutil.ReuseSlab[inverterVar]
	scopeSlab          jutil.ReuseSlab[renameInverterScopeContext]
	disallowedNames    map[string]struct{}
	skippingModuleBody bool
}

var renameInverterVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_NAME)
	s.Add(rhino.Token_IMPORT_STAR)
	return s
}()

func (*MakeDeclaredNamesUnique_ContextualRenameInverter) visitTokens() *ast.TokenSet {
	return &renameInverterVisitTokens
}

type renameInverterScopeCreator struct {
	sc *SyntacticScopeCreator
}

func (r *renameInverterScopeCreator) CreateScope(n rhino.Node, parent AbstractScope) AbstractScope {
	if parent == nil {
		// ContextualRenameInverter never reads global variable declarations;
		// avoid scanning the global AST when instantiating a local scope.
		return newScope(n)
	}
	return r.sc.CreateScopeOfScope(n, parent.(*Scope))
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter(AbstractCompiler)
func NewMakeDeclaredNamesUnique_ContextualRenameInverter(
	compiler AbstractCompiler,
) *MakeDeclaredNamesUnique_ContextualRenameInverter {
	return &MakeDeclaredNamesUnique_ContextualRenameInverter{
		compiler:   compiler,
		shadowByID: make(map[uint32]*inverterVar),
	}
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.process(Node, Node)
func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) Process(externs rhino.Node, js rhino.Node) {
	sc := &renameInverterScopeCreator{sc: NewSyntacticScopeCreator(c.compiler)}
	b := NodeTraversal_Builder{
		compiler:     c.compiler,
		callback:     c,
		scopeCreator: sc,
	}
	b.Traverse(js)
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.getOriginalName(String)
func MakeDeclaredNamesUnique_ContextualRenameInverter_GetOriginalName(name string) string {
	index := MakeDeclaredNamesUnique_ContextualRenameInverter_IndexOfSeparator(name)
	if index <= 0 {
		return name
	}
	return jutil.StringSubstring(name, 0, index)
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.indexOfSeparator(String)
func MakeDeclaredNamesUnique_ContextualRenameInverter_IndexOfSeparator(name string) int {
	return jutil.StringLastIndexOf(name, MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR)
}

func renameInverterSubtreeHasSeparator(n rhino.Node) bool {
	tok := n.GetToken()
	if tok == rhino.Token_NAME || tok == rhino.Token_IMPORT_STAR {
		if n.GetStringID() > rhino.RhinoStringPool_EmptyID {
			s := n.GetString()
			if len(s) > len(MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR) &&
				strings.LastIndex(s, MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR) > 0 {
				return true
			}
		}
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if renameInverterSubtreeHasSeparator(c) {
			return true
		}
	}
	return false
}

// inverterVar unifies Java's VariableInfo, DeclaredVariableInfo,
// UndeclaredVariableInfo, and RenamingInfo into a single concrete struct.
type inverterVar struct {
	inverter *MakeDeclaredNamesUnique_ContextualRenameInverter

	name          string
	nameID        uint32
	preferredName string
	scopeDepth    int
	prevShadow    *inverterVar

	// Non-nil (pointing to this inverterVar) while renaming is pending;
	// nil when not renamable or after TryToInvertName().
	renamingInfo *inverterVar

	referenceNodes           []rhino.Node
	potentialShadowVariables jutil.LinkedHashSet[*inverterVar]
}

type MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo = inverterVar
type MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo = *inverterVar

func (v *inverterVar) recycle() {
	v.inverter = nil
	v.name = ""
	v.nameID = 0
	v.preferredName = ""
	v.scopeDepth = 0
	v.prevShadow = nil
	v.renamingInfo = nil
	if cap(v.referenceNodes) <= 64 {
		clear(v.referenceNodes)
		v.referenceNodes = v.referenceNodes[:0]
	} else {
		v.referenceNodes = nil
	}
	if v.potentialShadowVariables.Size() > 64 {
		v.potentialShadowVariables = jutil.LinkedHashSet[*inverterVar]{}
	} else {
		v.potentialShadowVariables.Clear()
	}
}

func NewMakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo(
	this_ *MakeDeclaredNamesUnique_ContextualRenameInverter,
	currentName string,
	preferredName string,
) *inverterVar {
	v := &inverterVar{
		inverter:      this_,
		name:          currentName,
		nameID:        rhino.RhinoStringPool_InternID(currentName),
		preferredName: preferredName,
	}
	v.renamingInfo = v
	return v
}

func NewMakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo(
	var_ *Var,
	renamingInfo *inverterVar,
) *inverterVar {
	currentName := var_.GetName()
	id := var_.nameID
	if id == 0 {
		id = rhino.RhinoStringPool_InternID(currentName)
	}
	depth := var_.GetScope().(*Scope).GetDepth()
	if renamingInfo != nil {
		renamingInfo.name = currentName
		renamingInfo.nameID = id
		renamingInfo.scopeDepth = depth
		renamingInfo.renamingInfo = renamingInfo
		return renamingInfo
	}
	return &inverterVar{
		name:       currentName,
		nameID:     id,
		scopeDepth: depth,
	}
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfo.getName()
func (v *inverterVar) GetName() string {
	return v.name
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfo.getScopeDepth()
func (v *inverterVar) GetScopeDepth() int {
	return v.scopeDepth
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfo.addPotentialConflict(MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfo)
func (v *inverterVar) AddPotentialConflict(other *inverterVar) {
	if v.renamingInfo != nil && other != v {
		v.potentialShadowVariables.Add(other)
	}
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfo.addReferenceNode(Node)
func (v *inverterVar) AddReferenceNode(nameNode rhino.Node) {
	if v.renamingInfo != nil {
		v.referenceNodes = append(v.referenceNodes, nameNode)
	}
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.DeclaredVariableInfo.tryToInvertName()
func (v *inverterVar) TryToInvertName() {
	if v.renamingInfo != nil {
		v.name = v.attemptRename()
		v.renamingInfo = nil
	}
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.RenamingInfo.attemptRename()
func (v *inverterVar) attemptRename() string {
	// If we somehow ended up with "arguments$jscomp$..." it's not safe to rename
	// that to "arguments", because that name is special, but it would still be good
	// to simplify its name, if possible.
	conflict := v.preferredName == MakeDeclaredNamesUnique_ARGUMENTS
	if !conflict {
		for potentialShadowVariable := range v.potentialShadowVariables.All() {
			if potentialShadowVariable.name == v.preferredName {
				conflict = true
				break
			}
		}
	}
	newName := v.preferredName
	if conflict {
		var disallowedNames map[string]struct{}
		if v.inverter != nil {
			if v.inverter.disallowedNames == nil {
				v.inverter.disallowedNames = make(map[string]struct{}, v.potentialShadowVariables.Size()+1)
			} else {
				clear(v.inverter.disallowedNames)
			}
			disallowedNames = v.inverter.disallowedNames
		} else {
			disallowedNames = make(map[string]struct{}, v.potentialShadowVariables.Size()+1)
		}
		disallowedNames[MakeDeclaredNamesUnique_ARGUMENTS] = struct{}{}
		for potentialShadowVariable := range v.potentialShadowVariables.All() {
			// Look up names immediately before making our rename decision, because
			// some of these variables could themselves have been renamed.
			disallowedNames[potentialShadowVariable.name] = struct{}{}
		}
		baseName := v.preferredName + MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR
		i := 0
		for {
			newName = baseName + strconv.Itoa(i)
			i++
			if _, exists := disallowedNames[newName]; !exists {
				break
			}
		}
	}

	if newName == v.name {
		return v.name
	}

	for _, referenceNode := range v.referenceNodes {
		referenceNode.SetString(newName)
		if v.inverter != nil && v.inverter.compiler != nil {
			v.inverter.compiler.ReportChangeToEnclosingScope(referenceNode)
			parent := referenceNode.GetParent()
			if parent.IsFunction() && NodeUtil_IsFunctionDeclaration(parent) {
				v.inverter.compiler.ReportChangeToEnclosingScope(parent)
			}
		}
	}

	return newName
}

func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) allocDeclaredVar(
	var_ *Var, scopeDepth int,
) *inverterVar {
	v := c.varSlab.Alloc()
	currentName := var_.GetName()
	id := var_.nameID
	if id == 0 {
		id = rhino.RhinoStringPool_InternID(currentName)
	}
	v.inverter = c
	v.name = currentName
	v.nameID = id
	v.scopeDepth = scopeDepth
	if len(currentName) > len(MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR) {
		if indexOfSeparator := strings.LastIndex(currentName, MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR); indexOfSeparator > 0 {
			// Note: indexOfSeparator is a byte index into currentName (and the
			// separator "$jscomp$" is ASCII), so currentName[:indexOfSeparator]
			// slices the exact prefix before "$jscomp$".
			v.preferredName = currentName[:indexOfSeparator]
			v.renamingInfo = v
		}
	}
	return v
}

func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) pushVar(v *inverterVar) {
	v.prevShadow = c.shadowByID[v.nameID]
	c.shadowByID[v.nameID] = v
}

func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) getOrCreateCurrentVar(
	id uint32, node rhino.Node,
) *inverterVar {
	if v := c.shadowByID[id]; v != nil {
		return v
	}
	v := c.varSlab.Alloc()
	v.inverter = c
	v.name = node.GetString()
	v.nameID = id
	c.shadowByID[id] = v
	return v
}

func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) popVar(expected *inverterVar) {
	top, ok := c.shadowByID[expected.nameID]
	if !ok {
		panic(jutil.NewNullPointerExceptionMsg("Nonexistent variable info requested: '" + expected.name + "'"))
	}
	jutil.CheckStateMsg1(
		top == expected,
		"Declared variable \"%s\" was not the top of the name stack",
		expected.name)
	if top.prevShadow == nil {
		delete(c.shadowByID, expected.nameID)
	} else {
		c.shadowByID[expected.nameID] = top.prevShadow
	}
	top.prevShadow = nil
}

// EnterScope prepares a set for the new scope.
//
// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.enterScope(NodeTraversal)
func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) EnterScope(t *NodeTraversal) {
	if t.InGlobalScope() {
		// Don't track any information for the global scope, because it is likely to be
		// very large, and we will never rename any variables it contains. MakeDeclaredNamesUnique
		// favors keeping global variables unchanged, so they won't need to be restored.
		return
	}
	if len(c.scopeContextStack) == 0 && t.GetScopeRoot().IsModuleBody() &&
		!renameInverterSubtreeHasSeparator(t.GetScopeRoot()) {
		c.skippingModuleBody = true
		return
	}

	scope := t.GetScope()
	var parent *renameInverterScopeContext
	if len(c.scopeContextStack) != 0 {
		parent = c.scopeContextStack[len(c.scopeContextStack)-1]
	}
	ctx := c.scopeSlab.Alloc()
	ctx.scope = scope
	ctx.parent = parent

	scopeDepth := scope.GetDepth()
	for _, var_ := range scope.varList {
		if var_ == nil {
			continue
		}
		variableInfo := c.allocDeclaredVar(var_, scopeDepth)
		ctx.declaredVariableInfos = append(ctx.declaredVariableInfos, variableInfo)
		if variableInfo.renamingInfo != nil {
			ctx.renamableDeclared = append(ctx.renamableDeclared, variableInfo)
		}
		c.pushVar(variableInfo)
	}
	c.scopeContextStack = append(c.scopeContextStack, ctx)
}

// ExitScope renames vars for the current scope, and merge any referenced
// names into the parent scope reference set.
//
// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.exitScope(NodeTraversal)
func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) ExitScope(t *NodeTraversal) {
	if t.InGlobalScope() {
		return
	}
	if c.skippingModuleBody {
		c.skippingModuleBody = false
		return
	}

	if len(c.scopeContextStack) == 0 {
		panic(jutil.NewNoSuchElementException()) // ArrayDeque.pop()
	}
	topIdx := len(c.scopeContextStack) - 1
	scopeContext := jutil.CheckNotNull(c.scopeContextStack[topIdx])
	c.scopeContextStack[topIdx] = nil
	c.scopeContextStack = c.scopeContextStack[:topIdx]

	// Tell all the variables we referenced and the variables that might conflict with them
	// about each other before trying to rename the variables in this scope.
	for referencedVariable := range scopeContext.referencedVariables.Keys() {
		scopeContext.RecordPotentialConflicts(referencedVariable)
	}

	for _, declaredVariableInfo := range scopeContext.declaredVariableInfos {
		c.popVar(declaredVariableInfo)
		declaredVariableInfo.TryToInvertName()
	}

	if len(c.scopeContextStack) == 0 {
		clear(c.shadowByID)
		c.varSlab.Reset((*inverterVar).recycle)
		c.scopeSlab.Reset((*renameInverterScopeContext).recycle)
	}
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.shouldTraverse(NodeTraversal, Node, Node)
func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) ShouldTraverse(
	t *NodeTraversal, n rhino.Node, parent rhino.Node,
) bool {
	if c.skippingModuleBody {
		return false
	}
	if len(c.scopeContextStack) == 0 && !parent.IsNil() {
		tok := n.GetToken()
		if tok == rhino.Token_FUNCTION || tok == rhino.Token_CLASS ||
			(!n.IsSyntheticBlock() && t.CreatesBlockScope(n)) {
			return renameInverterSubtreeHasSeparator(n)
		}
	}
	return true
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.visit(NodeTraversal, Node, Node)
func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) Visit(
	t *NodeTraversal, node rhino.Node, parent rhino.Node,
) {
	if len(c.scopeContextStack) == 0 {
		return
	}
	tok := node.GetToken()
	if tok == rhino.Token_NAME {
		id := node.GetStringID()
		if id <= rhino.RhinoStringPool_EmptyID {
			return
		}
		scopeContext := c.scopeContextStack[len(c.scopeContextStack)-1]
		variableInfo := c.getOrCreateCurrentVar(id, node)
		variableInfo.AddReferenceNode(node)
		scopeContext.AddReferencedVariableInfo(variableInfo)
	} else if tok == rhino.Token_IMPORT_STAR {
		id := node.GetStringID()
		scopeContext := c.scopeContextStack[len(c.scopeContextStack)-1]
		variableInfo := c.getOrCreateCurrentVar(id, node)
		variableInfo.AddReferenceNode(node)
		scopeContext.AddReferencedVariableInfo(variableInfo)
	}
}

// renameInverterScopeContext keeps track of information about variables
// declared and referenced in a local scope.
type renameInverterScopeContext struct {
	parent                *renameInverterScopeContext
	scope                 *Scope
	declaredVariableInfos []*inverterVar
	renamableDeclared     []*inverterVar

	// Will be filled with objects for all variables we reference in this scope.
	referencedVariables jutil.LinkedHashMap[*inverterVar, bool]
}

type MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext = renameInverterScopeContext

func (s *renameInverterScopeContext) recycle() {
	s.parent = nil
	s.scope = nil
	if cap(s.declaredVariableInfos) <= 64 {
		clear(s.declaredVariableInfos)
		s.declaredVariableInfos = s.declaredVariableInfos[:0]
	} else {
		s.declaredVariableInfos = nil
	}
	if cap(s.renamableDeclared) <= 64 {
		clear(s.renamableDeclared)
		s.renamableDeclared = s.renamableDeclared[:0]
	} else {
		s.renamableDeclared = nil
	}
	if s.referencedVariables.Size() > 64 {
		s.referencedVariables.Reset()
	} else {
		s.referencedVariables.Clear()
	}
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.ScopeContext(Scope, ImmutableList<MakeDeclaredNamesUnique.ContextualRenameInverter.DeclaredVariableInfo>, MakeDeclaredNamesUnique.ContextualRenameInverter.ScopeContext)
func NewMakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext(
	scope *Scope,
	declaredVariableInfos []*inverterVar,
	parent *renameInverterScopeContext,
) *renameInverterScopeContext {
	var renamable []*inverterVar
	for _, d := range declaredVariableInfos {
		if d.renamingInfo != nil {
			renamable = append(renamable, d)
		}
	}
	return &renameInverterScopeContext{
		scope:                 scope,
		declaredVariableInfos: declaredVariableInfos,
		renamableDeclared:     renamable,
		parent:                parent,
	}
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.ScopeContext.addReferencedVariableInfo(MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfo)
func (s *renameInverterScopeContext) AddReferencedVariableInfo(variableInfo *inverterVar) {
	s.referencedVariables.LoadOrStore(variableInfo, false)
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.ScopeContext.recordPotentialConflicts(MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfo)
func (s *renameInverterScopeContext) RecordPotentialConflicts(referencedVariable *inverterVar) {
	thisScopeDepth := s.scope.GetDepth()
	referencedVariableDeclarationDepth := referencedVariable.scopeDepth
	if thisScopeDepth >= referencedVariableDeclarationDepth {
		// A scope's declaration list is fixed. Recording this same pair again
		// can only re-add the same object identities to conflict sets. Remember
		// completed pairs, including recursive visits from distinct child scopes.
		// False values still represent this scope's queued own references; adding
		// an already processed child reference preserves the original own order.
		if len(s.declaredVariableInfos) >= 8 {
			if s.referencedVariables.Put(referencedVariable, true) {
				return
			}
		}
		if referencedVariable.renamingInfo != nil {
			for _, declaredVariableInfo := range s.declaredVariableInfos {
				referencedVariable.AddPotentialConflict(declaredVariableInfo)
				declaredVariableInfo.AddPotentialConflict(referencedVariable)
			}
		} else {
			for _, declaredVariableInfo := range s.renamableDeclared {
				declaredVariableInfo.AddPotentialConflict(referencedVariable)
			}
		}
		if s.parent != nil {
			s.parent.RecordPotentialConflicts(referencedVariable)
		}
	}
}

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
