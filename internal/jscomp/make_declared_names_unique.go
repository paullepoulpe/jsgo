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

// Java: MakeDeclaredNamesUnique.builder()
func MakeDeclaredNamesUnique_Builder_() MakeDeclaredNamesUnique_Builder {
	return MakeDeclaredNamesUnique_Builder{markChanges: true}
}

// Java: MakeDeclaredNamesUnique.Builder.withRenamer(MakeDeclaredNamesUnique.Renamer)
func (b MakeDeclaredNamesUnique_Builder) WithRenamer(renamer MakeDeclaredNamesUnique_Renamer) MakeDeclaredNamesUnique_Builder {
	b.renamer = renamer
	return b
}

// WithMarkChanges enables the option to report any code changes to the
// compiler object.
//
// This option is true by default.
//
// Java: MakeDeclaredNamesUnique.Builder.withMarkChanges(boolean)
func (b MakeDeclaredNamesUnique_Builder) WithMarkChanges(markChanges bool) MakeDeclaredNamesUnique_Builder {
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
func (b MakeDeclaredNamesUnique_Builder) WithAssertOnChange(assertOnChange bool) MakeDeclaredNamesUnique_Builder {
	b.assertOnChange = assertOnChange
	return b
}

// Java: MakeDeclaredNamesUnique.Builder.build()
func (b MakeDeclaredNamesUnique_Builder) Build() *MakeDeclaredNamesUnique {
	if b.renamer == nil {
		b.renamer = NewMakeDeclaredNamesUnique_ContextualRenamer()
	}
	return NewMakeDeclaredNamesUnique(b.renamer, b.markChanges, b.assertOnChange)
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

	// GetReplacementNameValue returns a replacement name and true, or ("", false)
	// if oldName is unknown or should not be replaced.
	//
	// Java: MakeDeclaredNamesUnique.Renamer.getReplacementName(String)
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
	sc := NewSyntacticScopeCreator(c.compiler)
	sc.skipGlobalScan = true
	t := acquireNodeTraversalDirect(c.compiler, c, sc, false)
	t.Traverse(js)
	releaseNodeTraversal(t)
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
