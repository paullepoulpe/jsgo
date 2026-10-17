// Port of com/google/javascript/jscomp/MakeDeclaredNamesUnique.java @ closure-compiler v20261015

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

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
	renamerStack   []MakeDeclaredNamesUnique_Renamer
	rootRenamer    MakeDeclaredNamesUnique_Renamer
	markChanges    bool
	assertOnChange bool
}

// Arguments is special cased to handle cases where a local name shadows
// the arguments declaration.
const MakeDeclaredNamesUnique_ARGUMENTS = "arguments"

// Java: MakeDeclaredNamesUnique(MakeDeclaredNamesUnique.Renamer, boolean, boolean)
func NewMakeDeclaredNamesUnique(
	renamer MakeDeclaredNamesUnique_Renamer, markChanges bool, assertOnChange bool) *MakeDeclaredNamesUnique {
	m := &MakeDeclaredNamesUnique{}
	m.rootRenamer = renamer
	m.markChanges = markChanges
	m.assertOnChange = assertOnChange
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
		_, isContextualRenamer := m.rootRenamer.(MakeDeclaredNamesUnique_ContextualRenamer)
		jutil.CheckState(!declarationRoot.IsFunction() || !isContextualRenamer)
		renamer = m.rootRenamer
	} else {
		hoist := !declarationRoot.IsFunction() && !NodeUtil_CreatesBlockScope(declarationRoot)
		renamer = m.renamerStack[len(m.renamerStack)-1].CreateForChildScope(t.GetScopeRoot(), hoist)
	}

	m.renamerStack = append(m.renamerStack, renamer)

	m.FindDeclaredNames(t, declarationRoot)
}

// Java: MakeDeclaredNamesUnique.exitScope(NodeTraversal)
func (m *MakeDeclaredNamesUnique) ExitScope(t *NodeTraversal) {
	if !t.InGlobalScope() {
		if len(m.renamerStack) == 0 {
			panic(jutil.NewNoSuchElementException()) // ArrayDeque.pop()
		}
		m.renamerStack = m.renamerStack[:len(m.renamerStack)-1]
	}
}

// Java: MakeDeclaredNamesUnique.visit(NodeTraversal, Node, Node)
func (m *MakeDeclaredNamesUnique) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_NAME, rhino.Token_IMPORT_STAR:
		m.VisitNameOrImportStar(t, n, parent)
	default:
	}
}

// Java: MakeDeclaredNamesUnique.visitNameOrImportStar(NodeTraversal, Node, Node)
func (m *MakeDeclaredNamesUnique) VisitNameOrImportStar(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	// Don't rename the exported name foo in export {a as foo}; or import {foo as b};
	if n.IsName() && NodeUtil_IsNonlocalModuleExportName(n) {
		return
	}
	newName := m.GetReplacementName(n.GetString())
	if newName == nil {
		return
	}
	if m.assertOnChange {
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected renaming in MakeDeclaredNamesUnique. new name: " + *newName + " for " + n.String()))
	}

	renamer := m.renamerStack[len(m.renamerStack)-1]
	if renamer.StripConstIfReplaced() {
		n.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, false)
		jsDocInfoNode := NodeUtil_GetBestJSDocInfoNode(n)
		if jsDocInfoNode != nil && jsDocInfoNode.GetJSDocInfo() != nil {
			builder := rhino.JSDocInfo_Builder_CopyFrom(jsDocInfoNode.GetJSDocInfo())
			builder.RecordMutable()
			jsDocInfoNode.SetJSDocInfo(builder.Build())
		}
	}
	n.SetString(*newName)
	if m.markChanges {
		t.ReportCodeChange()
		// If we are renaming a function declaration, make sure the containing scope
		// has the opportunity to act on the change.
		if parent.IsFunction() && NodeUtil_IsFunctionDeclaration(parent) {
			t.GetCompiler().ReportChangeToEnclosingScope(parent)
		}
	}
}

// GetReplacementName walks the stack of name maps and finds the replacement
// name for the current scope.
//
// Java: MakeDeclaredNamesUnique.getReplacementName(String)
func (m *MakeDeclaredNamesUnique) GetReplacementName(oldName string) *string {
	for i := len(m.renamerStack) - 1; i >= 0; i-- {
		renamer := m.renamerStack[i]
		newName := renamer.GetReplacementName(oldName)
		if newName != nil {
			return newName
		}
	}
	return nil
}

// FindDeclaredNames traverses the current scope and collects declared
// names by calling addDeclaredName on the Renamer that is at the top of the
// renamerStack.
//
// Java: MakeDeclaredNamesUnique.findDeclaredNames(NodeTraversal, Node)
func (m *MakeDeclaredNamesUnique) FindDeclaredNames(t *NodeTraversal, n *rhino.Node) {
	jutil.CheckStateMsg1(NodeUtil_CreatesScope(n) || n.IsScript(), "%s", n)

	for v := range t.GetScope().GetVarIterable() {
		m.renamerStack[len(m.renamerStack)-1].AddDeclaredName(v.GetName(), false)
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

	// GetReplacementName returns a replacement name, null if oldName is
	// unknown or should not be replaced.
	//
	// Java: MakeDeclaredNamesUnique.Renamer.getReplacementName(String)
	GetReplacementName(oldName string) *string

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
	CreateForChildScope(scopeRoot *rhino.Node, hoisted bool) MakeDeclaredNamesUnique_Renamer

	// GetHoistRenamer returns the closest hoisting target for var and
	// function declarations.
	//
	// Java: MakeDeclaredNamesUnique.Renamer.getHoistRenamer()
	GetHoistRenamer() MakeDeclaredNamesUnique_Renamer
}

// MakeDeclaredNamesUnique_ContextualRenameInverter inverts the
// transformation by ContextualRenamer, when possible.
type MakeDeclaredNamesUnique_ContextualRenameInverter struct {
	compiler AbstractCompiler

	// The top of this stack is always an object storing information about the current variable
	// scope or null for the global scope.
	//
	// Go-only: language-adapter — the ArrayDeque used as a stack is a slice
	// whose top is the last element.
	scopeContextStack []*MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext

	// Keeps track of all VariableInfo objects for all variables visible in the current
	// variable scope.
	//
	// Basically this is a map whose keys are names and whose values are stacks of
	// VariableInfo objects. The top of each stack represents the variable with that name that is
	// visible in the scope we are currently traversing.
	variablesStackedByName *MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter(AbstractCompiler)
func NewMakeDeclaredNamesUnique_ContextualRenameInverter(
	compiler AbstractCompiler) *MakeDeclaredNamesUnique_ContextualRenameInverter {
	c := &MakeDeclaredNamesUnique_ContextualRenameInverter{}
	c.variablesStackedByName = NewMakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap(c)
	c.compiler = compiler
	return c
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.process(Node, Node)
func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) Process(externs *rhino.Node, js *rhino.Node) {
	NodeTraversal_Traverse(c.compiler, js, c)
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

	scope := t.GetScope()

	// Create DeclaredVariableInfo objects for all variables declared in this scope.
	// Push them onto the variablesByName data structure to make them visible,
	// and store them into a ScopeContext object for the scope we're entering.
	var declaredVariablesBuilder []*MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo
	for var_ := range scope.GetVarIterable() {
		variableInfo := c.CreateDeclaredVariableInfo(var_)
		declaredVariablesBuilder = append(declaredVariablesBuilder, variableInfo)
		c.variablesStackedByName.PushVariableInfo(variableInfo)
	}
	var parent *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext
	if len(c.scopeContextStack) != 0 {
		parent = c.scopeContextStack[len(c.scopeContextStack)-1]
	}
	c.scopeContextStack = append(c.scopeContextStack,
		NewMakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext(scope, declaredVariablesBuilder, parent))
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.createDeclaredVariableInfo(Var)
func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) CreateDeclaredVariableInfo(
	var_ *Var) *MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo {
	currentName := var_.GetName()
	indexOfSeparator := MakeDeclaredNamesUnique_ContextualRenameInverter_IndexOfSeparator(currentName)
	var renamingInfo *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo
	if indexOfSeparator > 0 {
		// We check for an index > 0, because a variable created directly by the compiler generally
		// begin with `$jscomp$` (the separator). We can't and don't want to try to rename any
		// variable to an empty string.
		invertedName := jutil.StringSubstring(currentName, 0, indexOfSeparator)
		renamingInfo = NewMakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo(c, currentName, invertedName)
	} else {
		renamingInfo = nil
	}
	return NewMakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo(var_, renamingInfo)
}

// MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo keeps track
// of information needed to rename a DeclaredVariableInfo.
type MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo struct {
	this_ *MakeDeclaredNamesUnique_ContextualRenameInverter // the enclosing instance

	// Will be filled with all nodes referring to the variable.
	referenceNodes []*rhino.Node

	// Will be filled with objects representing variables whose names must not conflict with this
	// one.
	potentialShadowVariables *jutil.LinkedHashSet[MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo]

	// The name the variable has currently
	currentName string

	// The name we want to use when renaming the variable.
	preferredName string
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.RenamingInfo(String, String)
func NewMakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo(
	this_ *MakeDeclaredNamesUnique_ContextualRenameInverter,
	currentName string,
	preferredName string) *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo {
	r := &MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo{this_: this_}
	r.potentialShadowVariables = jutil.NewLinkedHashSet[MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo]()
	r.currentName = currentName
	r.preferredName = preferredName
	return r
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.RenamingInfo.addPotentialShadowVariable(MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfo)
func (r *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo) AddPotentialShadowVariable(
	variableInfo MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo) {
	r.potentialShadowVariables.Add(variableInfo)
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.RenamingInfo.attemptRename()
func (r *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo) AttemptRename() string {
	disallowedNames := jutil.NewLinkedHashSet[string]()
	// If we somehow ended up with "arguments$jscomp$..." it's not safe to rename
	// that to "arguments", because that name is special, but it would still be good
	// to simplify its name, if possible. See the note below regarding why we
	// rename variables to something other than their original name.
	disallowedNames.Add(MakeDeclaredNamesUnique_ARGUMENTS)
	for potentialShadowVariable := range r.potentialShadowVariables.All() {
		// NOTE: We need to look up these names immediately before making our rename
		// decision, because some of these variables could themselves have been renamed.
		potentialShadowName := potentialShadowVariable.GetName()
		disallowedNames.Add(potentialShadowName)
	}
	newName := r.preferredName
	if disallowedNames.Contains(r.preferredName) {
		// Why are we bothering to find a name other than the original one?
		// The reason is that we would like the final output name to depend more on the shape
		// of the output code than on the process the compiler followed to reach that shape.
		// This keeps our unit tests more stable. Without it, a tweak to the details of some
		// optimization pass would be more likely to make a trivial change to the output that
		// breaks unit tests.
		baseName := r.preferredName + MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR
		i := 0
		for {
			newName = baseName + strconv.Itoa(i)
			i++
			if !disallowedNames.Contains(newName) {
				break
			}
		}
	}

	// It's possible we ended up generating the same name we already had.
	if newName == r.currentName {
		return r.currentName
	}

	for _, referenceNode := range r.referenceNodes {
		referenceNode.SetString(newName)
		// NOTE: This could probably be made more efficient by only reporting after doing all the
		// nodes that occur in the same scope.
		r.this_.compiler.ReportChangeToEnclosingScope(referenceNode)
		parent := referenceNode.GetParent()
		// If we are renaming a function declaration, make sure the containing scope
		// has the opportunity to act on the change.
		if parent.IsFunction() && NodeUtil_IsFunctionDeclaration(parent) {
			r.this_.compiler.ReportChangeToEnclosingScope(parent)
		}
	}

	return newName
}

// MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo
// records information about a variable for which we have seen a non-global
// declaration.
type MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo struct {
	var_ *Var
	name string

	// Keep track of information needed to rename this variable, if that is possible.
	//
	// This will be null if renaming has already been done or cannot be done.
	renamingInfo *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.DeclaredVariableInfo(Var, MakeDeclaredNamesUnique.ContextualRenameInverter.RenamingInfo)
func NewMakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo(
	var_ *Var,
	renamingInfo *MakeDeclaredNamesUnique_ContextualRenameInverter_RenamingInfo,
) *MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo {
	d := &MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo{}
	d.var_ = var_
	d.name = var_.GetName()
	d.renamingInfo = renamingInfo
	return d
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.DeclaredVariableInfo.getName()
func (d *MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo) GetName() string {
	return d.name
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.DeclaredVariableInfo.getScopeDepth()
func (d *MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo) GetScopeDepth() int {
	// Go-only: language-adapter — AbstractVar.getScope() returns the StaticScope
	// interface (see abstract_scope.go); a Var's scope is a *Scope.
	return d.var_.GetScope().(*Scope).GetDepth()
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.DeclaredVariableInfo.addPotentialConflict(MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfo)
func (d *MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo) AddPotentialConflict(
	variableInfo MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo) {
	if d.renamingInfo != nil && variableInfo != MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo(d) {
		d.renamingInfo.AddPotentialShadowVariable(variableInfo)
	} // else we won't rename, so we don't care.
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.DeclaredVariableInfo.addReferenceNode(Node)
func (d *MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo) AddReferenceNode(nameNode *rhino.Node) {
	if d.renamingInfo != nil {
		d.renamingInfo.referenceNodes = append(d.renamingInfo.referenceNodes, nameNode)
	} // else we won't rename, so we don't care
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.DeclaredVariableInfo.tryToInvertName()
func (d *MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo) TryToInvertName() {
	if d.renamingInfo != nil {
		// update the name for the sake of other variables that may check this one for
		// conflicts.
		d.name = d.renamingInfo.AttemptRename()
		d.renamingInfo = nil // allow garbage collection
	}
}

// ExitScope renames vars for the current scope, and merge any referenced
// names into the parent scope reference set.
//
// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.exitScope(NodeTraversal)
func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) ExitScope(t *NodeTraversal) {
	if t.InGlobalScope() {
		return
	}

	if len(c.scopeContextStack) == 0 {
		panic(jutil.NewNoSuchElementException()) // ArrayDeque.pop()
	}
	scopeContext := jutil.CheckNotNull(c.scopeContextStack[len(c.scopeContextStack)-1])
	c.scopeContextStack = c.scopeContextStack[:len(c.scopeContextStack)-1]

	// Tell all the variables we referenced and the variables that might conflict with them
	// about each other.
	// This needs to happen before we start trying to rename the variables in this scope next.
	for referencedVariable := range scopeContext.referencedVariables.All() {
		scopeContext.RecordPotentialConflicts(referencedVariable)
	}

	// Pop each variable declared in the scope we're exiting off of the variable stack
	// data structure.
	// Try to rename each one.
	for _, declaredVariableInfo := range scopeContext.declaredVariableInfos {
		variableName := declaredVariableInfo.GetName()
		// Out of an abundance of caution:
		// 1. Pop the variable from its name stack and make sure it matches
		// 2. Do this before renaming causes the name to be different.
		oldStackTop := c.variablesStackedByName.PopVariableInfo(variableName)
		jutil.CheckStateMsg1(
			oldStackTop == MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo(declaredVariableInfo),
			"Declared variable \"%s\" was not the top of the name stack",
			variableName)
		declaredVariableInfo.TryToInvertName()
	}

	if len(c.scopeContextStack) == 0 {
		// clear the records for any global or undeclared variables before we start traversing
		// a new local scope.
		c.variablesStackedByName.Clear()
	}
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.shouldTraverse(NodeTraversal, Node, Node)
func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) ShouldTraverse(
	t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	return true
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.visit(NodeTraversal, Node, Node)
func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) Visit(
	t *NodeTraversal, node *rhino.Node, parent *rhino.Node) {
	var scopeContext *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext
	if len(c.scopeContextStack) != 0 {
		scopeContext = c.scopeContextStack[len(c.scopeContextStack)-1]
	}
	if scopeContext != nil {
		// We're not in the global scope
		if NodeUtil_IsReferenceName(node) || node.IsImportStar() {
			name := node.GetString()
			// Get the corresponding VariableInfo, creating one if we haven't found a declaration
			// for it.
			variableInfo := c.variablesStackedByName.GetOrCreateCurrentVariableInfo(name)
			// add this reference to the variable
			variableInfo.AddReferenceNode(node)
			// tell the current scope that it references the variable
			scopeContext.AddReferencedVariableInfo(variableInfo)
		}
	} // else nothing to do for the global scope
}

// MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext keeps track
// of information about variables declared and referenced in a local scope.
//
// We do not create one of these objects for the global scope, because we
// know those won't be renamed. It isn't worthwhile to track all of them.
type MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext struct {
	parent                *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext
	scope                 *Scope
	declaredVariableInfos []*MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo

	// Will be filled with objects for all variables we reference in this scope.
	referencedVariables *jutil.LinkedHashSet[MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo]
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.ScopeContext(Scope, ImmutableList<MakeDeclaredNamesUnique.ContextualRenameInverter.DeclaredVariableInfo>, MakeDeclaredNamesUnique.ContextualRenameInverter.ScopeContext)
func NewMakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext(
	scope *Scope,
	declaredVariableInfos []*MakeDeclaredNamesUnique_ContextualRenameInverter_DeclaredVariableInfo,
	parent *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext,
) *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext {
	s := &MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext{}
	s.referencedVariables = jutil.NewLinkedHashSet[MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo]()
	s.scope = scope
	s.declaredVariableInfos = declaredVariableInfos
	s.parent = parent
	return s
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.ScopeContext.addReferencedVariableInfo(MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfo)
func (s *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext) AddReferencedVariableInfo(
	variableInfo MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo) {
	s.referencedVariables.Add(variableInfo)
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.ScopeContext.recordPotentialConflicts(MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfo)
func (s *MakeDeclaredNamesUnique_ContextualRenameInverter_ScopeContext) RecordPotentialConflicts(
	referencedVariable MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo) {
	thisScopeDepth := s.scope.GetDepth()
	referencedVariableDeclarationDepth := referencedVariable.GetScopeDepth()
	if thisScopeDepth >= referencedVariableDeclarationDepth {
		for _, declaredVariableInfo := range s.declaredVariableInfos {
			// NOTE: No need to check whether the 2 variables are actually the same here.
			// The logic in addPotentialConflict() can do that more efficiently.

			// The referenced variable must be sure not to "hide behind" another variable
			// declared between its declaration and the location where it is referenced.
			referencedVariable.AddPotentialConflict(declaredVariableInfo)

			// The declared variable must be sure not to "shadow" a variable declared in its own
			// or a higher scope, thus preventing the reference from seeing it.
			declaredVariableInfo.AddPotentialConflict(referencedVariable)
		}
		if s.parent != nil {
			s.parent.RecordPotentialConflicts(referencedVariable)
		}
	}
}

// MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap
// handles storing and looking up VariableInfo objects for the variables
// that are visible in the current ScopeContext.
type MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap struct {
	this_ *MakeDeclaredNamesUnique_ContextualRenameInverter // the enclosing instance

	// Go-only: language-adapter — each ArrayDeque used as a stack is a
	// *[]VariableInfo whose top is the last element.
	mapOfStacks *jutil.LinkedHashMap[string, *[]MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo]
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfoStackMap()
func NewMakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap(
	this_ *MakeDeclaredNamesUnique_ContextualRenameInverter,
) *MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap {
	return &MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap{
		this_:       this_,
		mapOfStacks: jutil.NewLinkedHashMap[string, *[]MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo](),
	}
}

// PushVariableInfo pushes a VariableInfo onto the stack corresponding to
// its name.
//
// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfoStackMap.pushVariableInfo(MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfo)
func (v *MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap) PushVariableInfo(
	info MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo) {
	name := info.GetName()
	variableInfoStack := v.GetVariableInfoStack(name)
	*variableInfoStack = append(*variableInfoStack, info)
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfoStackMap.getVariableInfoStack(String)
func (v *MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap) GetVariableInfoStack(
	name string) *[]MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo {
	return v.mapOfStacks.ComputeIfAbsent(name,
		func(ignoredKey string) *[]MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo {
			return &[]MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo{}
		})
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfoStackMap.getOrCreateCurrentVariableInfo(String)
func (v *MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap) GetOrCreateCurrentVariableInfo(
	name string) MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo {
	variableInfoStack := v.GetVariableInfoStack(name)
	var variableInfo MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo
	if len(*variableInfoStack) == 0 {
		variableInfo = v.this_.CreateUndeclaredVariableInfo(name)
		*variableInfoStack = append(*variableInfoStack, variableInfo)
	} else {
		variableInfo = (*variableInfoStack)[len(*variableInfoStack)-1]
	}
	return variableInfo
}

// PopVariableInfo pops the topmost VariableInfo off of the stack for the
// given variable name.
//
// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfoStackMap.popVariableInfo(String)
func (v *MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap) PopVariableInfo(
	name string) MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo {
	variableInfos := v.mapOfStacks.Get(name)
	jutil.CheckNotNullMsg1(variableInfos, "Nonexistent variable info requested: '%s'", name)
	if len(*variableInfos) == 0 {
		panic(jutil.NewNoSuchElementException()) // ArrayDeque.pop()
	}
	info := (*variableInfos)[len(*variableInfos)-1]
	*variableInfos = (*variableInfos)[:len(*variableInfos)-1]
	if len(*variableInfos) == 0 {
		// Remove empty stacks, so they can be garbage collected.
		v.mapOfStacks.Remove(name)
	}
	return info
}

// Clear empties all the variable name stacks.
//
// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfoStackMap.clear()
func (v *MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfoStackMap) Clear() {
	v.mapOfStacks.Clear()
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.createUndeclaredVariableInfo(String)
func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) CreateUndeclaredVariableInfo(
	name string) MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo {
	return NewMakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo(name)
}

// MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo
// represents a variable we've seen referenced but not declared.
//
// One of these will be created when a global variable or simply undeclared
// variable is referenced. It doesn't do much more than keep track of the
// name and the fact that we don't have a declaration for it (as a 0 scope
// depth).
type MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo struct {
	name string
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.UndeclaredVariableInfo(String)
func NewMakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo(
	name string) *MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo {
	return &MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo{name: name}
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.UndeclaredVariableInfo.addReferenceNode(Node)
func (u *MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo) AddReferenceNode(nameNode *rhino.Node) {
	// no need to track references
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.UndeclaredVariableInfo.getName()
func (u *MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo) GetName() string {
	return u.name
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.UndeclaredVariableInfo.getScopeDepth()
func (u *MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo) GetScopeDepth() int {
	// NOTE: Pretend that any variable for which we didn't see a declaration was declared in
	// the global scope, which is depth 0. This will most often be actually true, since we skip
	// creating variables for the global scope to avoid wasting time and space on variables that
	// will never be renamed.
	return 0
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.UndeclaredVariableInfo.addPotentialConflict(MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfo)
func (u *MakeDeclaredNamesUnique_ContextualRenameInverter_UndeclaredVariableInfo) AddPotentialConflict(
	variableInfo MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo) {
	// No need to record the potential conflict, because this variable will never be renamed.
}

// MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo: one of
// these will be created for every variable we encounter.
type MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo interface {
	// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfo.addReferenceNode(Node)
	AddReferenceNode(nameNode *rhino.Node)

	// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfo.getName()
	GetName() string

	// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfo.getScopeDepth()
	GetScopeDepth() int

	// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfo.addPotentialConflict(MakeDeclaredNamesUnique.ContextualRenameInverter.VariableInfo)
	AddPotentialConflict(variableInfo MakeDeclaredNamesUnique_ContextualRenameInverter_VariableInfo)
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
//
// ContextualRenamer has a subclass (BoilerplateRenamer), so it is this
// interface plus MakeDeclaredNamesUnique_ContextualRenamerBase.
type MakeDeclaredNamesUnique_ContextualRenamer interface {
	MakeDeclaredNamesUnique_Renamer
	GetMakeDeclaredNamesUnique_ContextualRenamerBase() *MakeDeclaredNamesUnique_ContextualRenamerBase
	String() string
}

// MakeDeclaredNamesUnique_ContextualRenamerBase holds the fields and base
// implementations of ContextualRenamer.
type MakeDeclaredNamesUnique_ContextualRenamerBase struct {
	self MakeDeclaredNamesUnique_ContextualRenamer

	scopeRoot *rhino.Node

	// This multiset is shared between this ContextualRenamer and its parent (and its parent's
	// parent, etc.) because it tracks counts of variables across the entire JS program.
	nameUsage *makeDeclaredNamesUniqueMultiset

	// By contrast, this is a different map for each ContextualRenamer because it's just keeping
	// track of the names used by this renamer.
	//
	// Go-only: language-adapter — values are nullable Java Strings.
	declarations *jutil.LinkedHashMap[string, *string]
	global       bool

	hoistRenamer MakeDeclaredNamesUnique_Renamer
}

const MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR = "$jscomp$"

// Java: MakeDeclaredNamesUnique.ContextualRenamer.toString()
func (c *MakeDeclaredNamesUnique_ContextualRenamerBase) String() string {
	// Go-only: language-adapter — MoreObjects.toStringHelper(this) prints
	// this.getClass().getSimpleName().
	className := "ContextualRenamer"
	if _, ok := c.self.(*MakeDeclaredNamesUnique_BoilerplateRenamer); ok {
		className = "BoilerplateRenamer"
	}
	return className + "{scopeRoot=" + jutil.ToJavaString(c.scopeRoot) +
		", nameUsage=" + c.nameUsage.String() +
		", declarations=" + c.declarations.String() +
		", global=" + jutil.ToJavaString(c.global) + "}"
}

// Java: MakeDeclaredNamesUnique.ContextualRenamer.getMakeDeclaredNamesUnique_ContextualRenamerBase()
//
// Go-only: language-adapter — field access through the interface.
func (c *MakeDeclaredNamesUnique_ContextualRenamerBase) GetMakeDeclaredNamesUnique_ContextualRenamerBase() *MakeDeclaredNamesUnique_ContextualRenamerBase {
	return c
}

// Java: MakeDeclaredNamesUnique.ContextualRenamer()
func NewMakeDeclaredNamesUnique_ContextualRenamer() *MakeDeclaredNamesUnique_ContextualRenamerBase {
	c := &MakeDeclaredNamesUnique_ContextualRenamerBase{}
	c.InitMakeDeclaredNamesUnique_ContextualRenamerBase(c)
	return c
}

// InitMakeDeclaredNamesUnique_ContextualRenamerBase is the body of the
// root constructor ContextualRenamer().
//
// Java: MakeDeclaredNamesUnique.ContextualRenamer()
func (c *MakeDeclaredNamesUnique_ContextualRenamerBase) InitMakeDeclaredNamesUnique_ContextualRenamerBase(
	self MakeDeclaredNamesUnique_ContextualRenamer) {
	c.self = self
	c.declarations = jutil.NewLinkedHashMap[string, *string]()
	c.scopeRoot = nil
	c.global = true
	c.nameUsage = newMakeDeclaredNamesUniqueMultiset()

	c.hoistRenamer = self
}

// NewMakeDeclaredNamesUnique_ContextualRenamerForChildScope is the
// constructor for child scopes.
//
// Java: MakeDeclaredNamesUnique.ContextualRenamer(Node, Multiset<String>, boolean, MakeDeclaredNamesUnique.Renamer)
func NewMakeDeclaredNamesUnique_ContextualRenamerForChildScope(
	scopeRoot *rhino.Node,
	nameUsage *makeDeclaredNamesUniqueMultiset,
	hoistingTargetScope bool,
	parent MakeDeclaredNamesUnique_Renamer) *MakeDeclaredNamesUnique_ContextualRenamerBase {
	c := &MakeDeclaredNamesUnique_ContextualRenamerBase{}
	c.self = c
	c.declarations = jutil.NewLinkedHashMap[string, *string]()
	jutil.CheckStateMsg1(NodeUtil_CreatesScope(scopeRoot), "%s", scopeRoot)

	if scopeRoot.IsFunction() {
		jutil.CheckStateMsg1(!hoistingTargetScope, "%s", scopeRoot)
	}

	c.scopeRoot = scopeRoot
	c.global = false
	c.nameUsage = nameUsage

	if hoistingTargetScope {
		jutil.CheckStateMsg1(!NodeUtil_CreatesBlockScope(scopeRoot), "%s", scopeRoot)
		c.hoistRenamer = c.self
	} else {
		jutil.CheckStateMsg1(NodeUtil_CreatesBlockScope(scopeRoot) || scopeRoot.IsFunction(), "%s", scopeRoot)
		c.hoistRenamer = parent.GetHoistRenamer()
	}
	return c
}

// CreateForChildScope creates a ContextualRenamer
//
// Java: MakeDeclaredNamesUnique.ContextualRenamer.createForChildScope(Node, boolean)
func (c *MakeDeclaredNamesUnique_ContextualRenamerBase) CreateForChildScope(
	scopeRoot *rhino.Node, hoistingTargetScope bool) MakeDeclaredNamesUnique_Renamer {
	return NewMakeDeclaredNamesUnique_ContextualRenamerForChildScope(scopeRoot, c.nameUsage, hoistingTargetScope, c.self)
}

// AddDeclaredName adds a name to the map of names declared in this scope.
//
// Java: MakeDeclaredNamesUnique.ContextualRenamer.addDeclaredName(String, boolean)
func (c *MakeDeclaredNamesUnique_ContextualRenamerBase) AddDeclaredName(name string, hoisted bool) {
	if hoisted && c.hoistRenamer != MakeDeclaredNamesUnique_Renamer(c.self) {
		c.hoistRenamer.AddDeclaredName(name, true)
	} else {
		if name != MakeDeclaredNamesUnique_ARGUMENTS {
			if c.global {
				c.ReserveName(name)
			} else {
				// It hasn't been declared locally yet, so increment the count.
				if !c.declarations.ContainsKey(name) {
					id := c.IncrementNameCount(name)
					var newName *string = nil
					if id != 0 {
						// this name is exists at another source location, create another
						// The stack of renamers traverse all scopes and track all declared
						// names and their counts. This means, they could generate conflicting names if it
						// does not see the entire program (e.g. when used via library level checks).
						// Currently, the renamer stack is only used by Normalize and ScopedAliases pass.
						// 1. When used in normalize, renamer stack see the whole program and tracks the
						// name counts. Hence, when generating a unique name, it's guaranteed to generate
						// names unique across the entire program.
						// 2. When used in ScopedAliases, the renamer stack sees only the current library
						// sources. But, in this use case, is only renames locals within a function body.
						// Hence it does not generate conflicting declarations.
						// We could stop tracking the counts and switch the renamer to use the file hashcode
						// based UniqueIdSupplier, but that requires updating several unit tests that run
						// normalize.
						unique := MakeDeclaredNamesUnique_ContextualRenamer_GetUniqueName(name, id)
						for c.nameUsage.Contains(unique) {
							// if the newName also exists at another location, create another
							id = c.IncrementNameCount(unique)
							unique = MakeDeclaredNamesUnique_ContextualRenamer_GetUniqueName(name, id)
						}
						c.ReserveName(unique) // reserve this new name so that it's never reused
						newName = &unique
					}
					c.declarations.Put(name, newName)
				}
			}
		}
	}
}

// Java: MakeDeclaredNamesUnique.ContextualRenamer.getReplacementName(String)
func (c *MakeDeclaredNamesUnique_ContextualRenamerBase) GetReplacementName(oldName string) *string {
	return c.declarations.Get(oldName)
}

// MakeDeclaredNamesUnique_ContextualRenamer_GetUniqueName: given a name
// and the associated id, create a new unique name.
//
// Java: MakeDeclaredNamesUnique.ContextualRenamer.getUniqueName(String, int)
func MakeDeclaredNamesUnique_ContextualRenamer_GetUniqueName(name string, id int) string {
	return name + MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR + strconv.Itoa(id)
}

// Java: MakeDeclaredNamesUnique.ContextualRenamer.reserveName(String)
func (c *MakeDeclaredNamesUnique_ContextualRenamerBase) ReserveName(name string) {
	c.nameUsage.SetCountConditionally(name, 0, 1)
}

// Java: MakeDeclaredNamesUnique.ContextualRenamer.incrementNameCount(String)
func (c *MakeDeclaredNamesUnique_ContextualRenamerBase) IncrementNameCount(name string) int {
	return c.nameUsage.AddOccurrences(name, 1)
}

// Java: MakeDeclaredNamesUnique.ContextualRenamer.stripConstIfReplaced()
func (c *MakeDeclaredNamesUnique_ContextualRenamerBase) StripConstIfReplaced() bool {
	return false
}

// Java: MakeDeclaredNamesUnique.ContextualRenamer.getHoistRenamer()
func (c *MakeDeclaredNamesUnique_ContextualRenamerBase) GetHoistRenamer() MakeDeclaredNamesUnique_Renamer {
	return c.hoistRenamer
}

// makeDeclaredNamesUniqueMultiset is the HashMultiset<String> of
// ContextualRenamer.nameUsage: count(e) (via Contains), add(e, n) and
// setCount(e, oldCount, newCount).
//
// Go-only: language-adapter — Guava HashMultiset. Its iteration order
// (Java HashMap order) is only observable through ContextualRenamer's
// debugging toString; this one iterates in insertion order.
type makeDeclaredNamesUniqueMultiset struct {
	counts *jutil.LinkedHashMap[string, int]
}

func newMakeDeclaredNamesUniqueMultiset() *makeDeclaredNamesUniqueMultiset {
	return &makeDeclaredNamesUniqueMultiset{counts: jutil.NewLinkedHashMap[string, int]()}
}

// Contains is Multiset.contains(Object).
func (m *makeDeclaredNamesUniqueMultiset) Contains(e string) bool {
	return m.counts.ContainsKey(e)
}

// AddOccurrences is Multiset.add(E, int): returns the count before the
// operation.
func (m *makeDeclaredNamesUniqueMultiset) AddOccurrences(e string, occurrences int) int {
	jutil.CheckArgumentMsg1(occurrences >= 0, "occurrences cannot be negative: %s", occurrences)
	old := m.counts.GetOrDefault(e, 0)
	if occurrences == 0 {
		return old
	}
	m.counts.Put(e, old+occurrences)
	return old
}

// SetCountConditionally is Multiset.setCount(E, int, int): sets the count
// to newCount only if it is currently oldCount.
func (m *makeDeclaredNamesUniqueMultiset) SetCountConditionally(e string, oldCount int, newCount int) bool {
	if m.counts.GetOrDefault(e, 0) != oldCount {
		return false
	}
	if newCount == 0 {
		m.counts.Remove(e)
	} else {
		m.counts.Put(e, newCount)
	}
	return true
}

// String is AbstractMultiset.toString(): entrySet().toString(), each entry
// "e" or "e x n".
func (m *makeDeclaredNamesUniqueMultiset) String() string {
	var sb strings.Builder
	sb.WriteByte('[')
	first := true
	for e, n := range m.counts.All() {
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(e)
		if n != 1 {
			sb.WriteString(" x ")
			sb.WriteString(strconv.Itoa(n))
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
	declarations     *jutil.LinkedHashMap[string, string]
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
	r.declarations = jutil.NewLinkedHashMap[string, string]()
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

// Java: MakeDeclaredNamesUnique.InlineRenamer.addDeclaredName(String, boolean)
func (r *MakeDeclaredNamesUnique_InlineRenamer) AddDeclaredName(name string, hoisted bool) {
	jutil.CheckState(name != MakeDeclaredNamesUnique_ARGUMENTS)
	if hoisted && r.hoistRenamer != MakeDeclaredNamesUnique_Renamer(r) {
		r.hoistRenamer.AddDeclaredName(name, hoisted)
	} else {
		r.declarations.ComputeIfAbsent(name, r.GetUniqueName)
	}
}

// Java: MakeDeclaredNamesUnique.InlineRenamer.getUniqueName(String)
func (r *MakeDeclaredNamesUnique_InlineRenamer) GetUniqueName(name string) string {
	if name == "" {
		return name
	}

	if strings.Contains(name, MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR) {
		name = jutil.StringSubstring(
			name, 0, jutil.StringLastIndexOf(name, MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR))
	}

	if r.convention.IsExported(name) {
		// The google internal coding convention includes a naming convention
		// to export names starting with "_".  Simply strip "_" those to avoid
		// exporting names.
		name = "JSCompiler_" + name
	}

	// By using the same separator the id will be stripped if it isn't
	// needed when variable renaming is turned off.
	return name + MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR + r.idPrefix + r.uniqueIdSupplier()
}

// Java: MakeDeclaredNamesUnique.InlineRenamer.getReplacementName(String)
func (r *MakeDeclaredNamesUnique_InlineRenamer) GetReplacementName(oldName string) *string {
	if v, ok := r.declarations.Lookup(oldName); ok {
		return &v
	}
	return nil
}

// Java: MakeDeclaredNamesUnique.InlineRenamer.createForChildScope(Node, boolean)
func (r *MakeDeclaredNamesUnique_InlineRenamer) CreateForChildScope(
	scopeRoot *rhino.Node, hoistingTargetScope bool) MakeDeclaredNamesUnique_Renamer {
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
	MakeDeclaredNamesUnique_ContextualRenamerBase
	uniqueIdSupplier func() string
	idPrefix         string
	convention       CodingConvention
}

// Java: MakeDeclaredNamesUnique.BoilerplateRenamer(CodingConvention, Supplier<String>, String)
func NewMakeDeclaredNamesUnique_BoilerplateRenamer(
	convention CodingConvention, uniqueIdSupplier func() string, idPrefix string) *MakeDeclaredNamesUnique_BoilerplateRenamer {
	b := &MakeDeclaredNamesUnique_BoilerplateRenamer{}
	b.InitMakeDeclaredNamesUnique_ContextualRenamerBase(b)
	b.convention = convention
	b.uniqueIdSupplier = uniqueIdSupplier
	b.idPrefix = idPrefix
	return b
}

// Java: MakeDeclaredNamesUnique.BoilerplateRenamer.createForChildScope(Node, boolean)
func (b *MakeDeclaredNamesUnique_BoilerplateRenamer) CreateForChildScope(
	scopeRoot *rhino.Node, hoisted bool) MakeDeclaredNamesUnique_Renamer {
	return NewMakeDeclaredNamesUnique_InlineRenamer(b.convention, b.uniqueIdSupplier, b.idPrefix, false, hoisted, b)
}

// MakeDeclaredNamesUnique_TargettedRenamer only renames things that match
// specific names. Wraps another renamer.
type MakeDeclaredNamesUnique_TargettedRenamer struct {
	delegate MakeDeclaredNamesUnique_Renamer
	targets  *jutil.LinkedHashSet[string]
}

// Java: MakeDeclaredNamesUnique.TargettedRenamer(MakeDeclaredNamesUnique.Renamer, Set<String>)
func NewMakeDeclaredNamesUnique_TargettedRenamer(
	delegate MakeDeclaredNamesUnique_Renamer, targets *jutil.LinkedHashSet[string]) *MakeDeclaredNamesUnique_TargettedRenamer {
	return &MakeDeclaredNamesUnique_TargettedRenamer{delegate: delegate, targets: targets}
}

// Java: MakeDeclaredNamesUnique.TargettedRenamer.addDeclaredName(String, boolean)
func (r *MakeDeclaredNamesUnique_TargettedRenamer) AddDeclaredName(name string, hoisted bool) {
	if r.targets.Contains(name) {
		r.delegate.AddDeclaredName(name, hoisted)
	}
}

// Java: MakeDeclaredNamesUnique.TargettedRenamer.getReplacementName(String)
func (r *MakeDeclaredNamesUnique_TargettedRenamer) GetReplacementName(oldName string) *string {
	if r.targets.Contains(oldName) {
		return r.delegate.GetReplacementName(oldName)
	}
	return nil
}

// Java: MakeDeclaredNamesUnique.TargettedRenamer.stripConstIfReplaced()
func (r *MakeDeclaredNamesUnique_TargettedRenamer) StripConstIfReplaced() bool {
	return r.delegate.StripConstIfReplaced()
}

// Java: MakeDeclaredNamesUnique.TargettedRenamer.createForChildScope(Node, boolean)
func (r *MakeDeclaredNamesUnique_TargettedRenamer) CreateForChildScope(
	scopeRoot *rhino.Node, hoistingTargetScope bool) MakeDeclaredNamesUnique_Renamer {
	return NewMakeDeclaredNamesUnique_TargettedRenamer(
		r.delegate.CreateForChildScope(scopeRoot, hoistingTargetScope), r.targets)
}

// Java: MakeDeclaredNamesUnique.TargettedRenamer.getHoistRenamer()
func (r *MakeDeclaredNamesUnique_TargettedRenamer) GetHoistRenamer() MakeDeclaredNamesUnique_Renamer {
	return r.delegate.GetHoistRenamer()
}
