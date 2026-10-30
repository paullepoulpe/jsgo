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

// NewMakeDeclaredNamesUnique creates a MakeDeclaredNamesUnique with the
// default ContextualRenamer, markChanges=true, and assertOnChange=false.
func NewMakeDeclaredNamesUnique() *MakeDeclaredNamesUnique {
	return NewMakeDeclaredNamesUniqueWith(nil, true, false)
}

// Java: MakeDeclaredNamesUnique(MakeDeclaredNamesUnique.Renamer, boolean, boolean)
func NewMakeDeclaredNamesUniqueWith(
	renamer MakeDeclaredNamesUnique_Renamer, markChanges bool, assertOnChange bool,
) *MakeDeclaredNamesUnique {
	m := &MakeDeclaredNamesUnique{}
	m.reset(renamer, markChanges, assertOnChange)
	return m
}

func (m *MakeDeclaredNamesUnique) reset(
	renamer MakeDeclaredNamesUnique_Renamer, markChanges bool, assertOnChange bool,
) {
	if renamer == nil {
		renamer = NewMakeDeclaredNamesUnique_ContextualRenamer()
	}
	clear(m.renamerStack)
	if cap(m.renamerStack) == 0 {
		m.renamerStack = m.renamerStackBuf[:0]
	} else {
		m.renamerStack = m.renamerStack[:0]
	}
	m.rootRenamer = renamer
	m.markChanges = markChanges
	m.assertOnChange = assertOnChange
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

type inverterRefLink struct {
	node rhino.Node
	next int32
}

type inverterVarLink struct {
	varID uint32
	next  int32
}

type inverterScopeVarKey struct {
	scopeSeq uint32
	varID    uint32
}

type inverterVar struct {
	name          string
	preferredName string
	inScopeMask   uint64
	nameID        uint32
	prevShadow    uint32
	scopeDepth    int32
	refHead       int32
	refTail       int32
	conflictHead  int32
	conflictTail  int32
	isCandidate   bool
}

type inverterScope struct {
	depth          int32
	scopeSeq       uint32
	declStart      int
	declEnd        int
	renamableStart int
	renamableEnd   int
	refHead        int32
	refTail        int32
}

// MakeDeclaredNamesUnique_ContextualRenameInverter inverts the
// transformation by ContextualRenamer, when possible.
type MakeDeclaredNamesUnique_ContextualRenameInverter struct {
	compiler AbstractCompiler

	// 1-indexed flat arena of active local and referenced outer/global variables
	// (vars[0] is reserved as the zero sentinel).
	vars           []inverterVar
	scopes         []inverterScope
	declVars       []uint32
	renamableDecls []uint32
	scopeRefNodes  []inverterVarLink
	refLinks       []inverterRefLink
	conflictLinks  []inverterVarLink

	// Maps StringPool ID to the innermost visible 1-based varID with that name;
	// outer shadowed variables are chained via inverterVar.prevShadow.
	shadowByID map[uint32]uint32

	// Fallback deduplication table only used when local scope nesting exceeds 64.
	deepScopeSeen map[inverterScopeVarKey]struct{}
	nextScopeSeq  uint32

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
		shadowByID: make(map[uint32]uint32),
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

func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) allocDeclaredVar(
	currentName string, id uint32, scopeDepth int32,
) uint32 {
	if len(c.vars) == 0 {
		c.vars = append(c.vars, inverterVar{})
	}
	varID := uint32(len(c.vars))
	v := inverterVar{
		name:         currentName,
		nameID:       id,
		scopeDepth:   scopeDepth,
		refHead:      -1,
		refTail:      -1,
		conflictHead: -1,
		conflictTail: -1,
	}
	if len(currentName) > len(MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR) {
		if idx := strings.LastIndex(currentName, MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR); idx > 0 {
			v.preferredName = currentName[:idx]
			v.isCandidate = true
		}
	}
	c.vars = append(c.vars, v)
	return varID
}

func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) pushVar(varID uint32) {
	v := &c.vars[varID]
	v.prevShadow = c.shadowByID[v.nameID]
	c.shadowByID[v.nameID] = varID
}

func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) getOrCreateCurrentVar(
	id uint32, node rhino.Node,
) uint32 {
	if varID := c.shadowByID[id]; varID != 0 {
		return varID
	}
	if len(c.vars) == 0 {
		c.vars = append(c.vars, inverterVar{})
	}
	varID := uint32(len(c.vars))
	c.vars = append(c.vars, inverterVar{
		name:         node.GetString(),
		nameID:       id,
		refHead:      -1,
		refTail:      -1,
		conflictHead: -1,
		conflictTail: -1,
	})
	c.shadowByID[id] = varID
	return varID
}

func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) popVar(varID uint32) {
	v := &c.vars[varID]
	top := c.shadowByID[v.nameID]
	if top == 0 {
		panic(jutil.NewNullPointerExceptionMsg("Nonexistent variable info requested: '" + v.name + "'"))
	}
	jutil.CheckStateMsg1(
		top == varID,
		"Declared variable \"%s\" was not the top of the name stack",
		v.name)
	if v.prevShadow == 0 {
		delete(c.shadowByID, v.nameID)
	} else {
		c.shadowByID[v.nameID] = v.prevShadow
	}
	v.prevShadow = 0
}

func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) addReferenceNode(varID uint32, nameNode rhino.Node) {
	v := &c.vars[varID]
	if !v.isCandidate {
		return
	}
	linkIdx := int32(len(c.refLinks))
	c.refLinks = append(c.refLinks, inverterRefLink{node: nameNode, next: -1})
	if v.refHead < 0 {
		v.refHead = linkIdx
	} else {
		c.refLinks[v.refTail].next = linkIdx
	}
	v.refTail = linkIdx
}

func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) addConflictEdge(varID, otherID uint32) {
	if varID == otherID {
		return
	}
	v := &c.vars[varID]
	if !v.isCandidate {
		return
	}
	linkIdx := int32(len(c.conflictLinks))
	c.conflictLinks = append(c.conflictLinks, inverterVarLink{varID: otherID, next: -1})
	if v.conflictHead < 0 {
		v.conflictHead = linkIdx
	} else {
		c.conflictLinks[v.conflictTail].next = linkIdx
	}
	v.conflictTail = linkIdx
}

func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) isVarInScopeRefs(
	stackIdx int, scopeSeq uint32, varID uint32,
) bool {
	if stackIdx < 64 {
		return (c.vars[varID].inScopeMask & (uint64(1) << stackIdx)) != 0
	}
	_, ok := c.deepScopeSeen[inverterScopeVarKey{scopeSeq: scopeSeq, varID: varID}]
	return ok
}

func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) addScopeRef(stackIdx int, varID uint32) {
	s := &c.scopes[stackIdx]
	if stackIdx < 64 {
		bit := uint64(1) << stackIdx
		if c.vars[varID].inScopeMask&bit != 0 {
			return
		}
		c.vars[varID].inScopeMask |= bit
	} else {
		if c.deepScopeSeen == nil {
			c.deepScopeSeen = make(map[inverterScopeVarKey]struct{})
		}
		key := inverterScopeVarKey{scopeSeq: s.scopeSeq, varID: varID}
		if _, ok := c.deepScopeSeen[key]; ok {
			return
		}
		c.deepScopeSeen[key] = struct{}{}
	}
	linkIdx := int32(len(c.scopeRefNodes))
	c.scopeRefNodes = append(c.scopeRefNodes, inverterVarLink{varID: varID, next: -1})
	if s.refHead < 0 {
		s.refHead = linkIdx
	} else {
		c.scopeRefNodes[s.refTail].next = linkIdx
	}
	s.refTail = linkIdx
}

func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) tryToInvertName(varID uint32) {
	v := &c.vars[varID]
	if !v.isCandidate {
		return
	}
	newName := c.attemptRename(v)
	v.name = newName
	v.isCandidate = false
}

func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) attemptRename(v *inverterVar) string {
	// If we somehow ended up with "arguments$jscomp$..." it's not safe to rename
	// that to "arguments", because that name is special, but it would still be good
	// to simplify its name, if possible.
	conflict := v.preferredName == MakeDeclaredNamesUnique_ARGUMENTS
	if !conflict {
		for idx := v.conflictHead; idx >= 0; idx = c.conflictLinks[idx].next {
			if c.vars[c.conflictLinks[idx].varID].name == v.preferredName {
				conflict = true
				break
			}
		}
	}
	newName := v.preferredName
	if conflict {
		if c.disallowedNames == nil {
			c.disallowedNames = make(map[string]struct{}, 16)
		} else {
			clear(c.disallowedNames)
		}
		c.disallowedNames[MakeDeclaredNamesUnique_ARGUMENTS] = struct{}{}
		for idx := v.conflictHead; idx >= 0; idx = c.conflictLinks[idx].next {
			// Look up names immediately before making our rename decision, because
			// some of these variables could themselves have been renamed.
			c.disallowedNames[c.vars[c.conflictLinks[idx].varID].name] = struct{}{}
		}
		baseName := v.preferredName + MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR
		for i := 0; ; i++ {
			newName = baseName + strconv.Itoa(i)
			if _, exists := c.disallowedNames[newName]; !exists {
				break
			}
		}
	}

	if newName == v.name {
		return v.name
	}

	for idx := v.refHead; idx >= 0; idx = c.refLinks[idx].next {
		referenceNode := c.refLinks[idx].node
		referenceNode.SetString(newName)
		if c.compiler != nil {
			c.compiler.ReportChangeToEnclosingScope(referenceNode)
			parent := referenceNode.GetParent()
			if parent.IsFunction() && NodeUtil_IsFunctionDeclaration(parent) {
				c.compiler.ReportChangeToEnclosingScope(parent)
			}
		}
	}

	return newName
}

func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) pushScope(scopeDepth int32, varList []*Var) {
	declStart := len(c.declVars)
	renamableStart := len(c.renamableDecls)
	for _, var_ := range varList {
		if var_ == nil {
			continue
		}
		currentName := var_.GetName()
		id := var_.nameID
		if id == 0 {
			id = rhino.RhinoStringPool_InternID(currentName)
		}
		varID := c.allocDeclaredVar(currentName, id, scopeDepth)
		c.declVars = append(c.declVars, varID)
		if c.vars[varID].isCandidate {
			c.renamableDecls = append(c.renamableDecls, varID)
		}
		c.pushVar(varID)
	}
	c.nextScopeSeq++
	c.scopes = append(c.scopes, inverterScope{
		depth:          scopeDepth,
		scopeSeq:       c.nextScopeSeq,
		declStart:      declStart,
		declEnd:        len(c.declVars),
		renamableStart: renamableStart,
		renamableEnd:   len(c.renamableDecls),
		refHead:        -1,
		refTail:        -1,
	})
}

func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) popScope() {
	if len(c.scopes) == 0 {
		panic(jutil.NewNoSuchElementException()) // ArrayDeque.pop()
	}
	stackIdx := len(c.scopes) - 1
	s := c.scopes[stackIdx]
	c.scopes = c.scopes[:stackIdx]

	decls := c.declVars[s.declStart:s.declEnd]
	renamable := c.renamableDecls[s.renamableStart:s.renamableEnd]

	// Tell all variables referenced in this scope (or propagated from child
	// scopes) and the variables declared in this scope about each other before
	// trying to rename the variables declared in this scope.
	for refIdx := s.refHead; refIdx >= 0; refIdx = c.scopeRefNodes[refIdx].next {
		refVarID := c.scopeRefNodes[refIdx].varID
		refVar := &c.vars[refVarID]
		if refVar.scopeDepth < s.depth {
			if refVar.isCandidate {
				for _, declID := range decls {
					c.addConflictEdge(refVarID, declID)
					c.addConflictEdge(declID, refVarID)
				}
			} else {
				for _, declID := range renamable {
					c.addConflictEdge(declID, refVarID)
				}
			}
		} else if refVar.scopeDepth == s.depth {
			if refVar.isCandidate {
				for _, declID := range decls {
					if declID == refVarID {
						continue
					}
					c.addConflictEdge(refVarID, declID)
					if !c.isVarInScopeRefs(stackIdx, s.scopeSeq, declID) {
						c.addConflictEdge(declID, refVarID)
					}
				}
			} else {
				for _, declID := range renamable {
					if !c.isVarInScopeRefs(stackIdx, s.scopeSeq, declID) {
						c.addConflictEdge(declID, refVarID)
					}
				}
			}
		}
	}

	var clearBit uint64
	if stackIdx < 64 {
		clearBit = uint64(1) << stackIdx
	}
	hasParent := stackIdx > 0
	var parentDepth int32
	if hasParent {
		parentDepth = c.scopes[stackIdx-1].depth
	}
	for refIdx := s.refHead; refIdx >= 0; refIdx = c.scopeRefNodes[refIdx].next {
		refVarID := c.scopeRefNodes[refIdx].varID
		if clearBit != 0 {
			c.vars[refVarID].inScopeMask &^= clearBit
		} else {
			delete(c.deepScopeSeen, inverterScopeVarKey{scopeSeq: s.scopeSeq, varID: refVarID})
		}
		if hasParent && c.vars[refVarID].scopeDepth <= parentDepth {
			c.addScopeRef(stackIdx-1, refVarID)
		}
	}

	for _, declID := range decls {
		c.popVar(declID)
		c.tryToInvertName(declID)
	}
	c.declVars = c.declVars[:s.declStart]
	c.renamableDecls = c.renamableDecls[:s.renamableStart]
}

func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) resetArena() {
	clear(c.shadowByID)
	clear(c.vars)
	c.vars = c.vars[:0]
	c.scopeRefNodes = c.scopeRefNodes[:0]
	c.refLinks = c.refLinks[:0]
	c.conflictLinks = c.conflictLinks[:0]
	c.nextScopeSeq = 0
	if len(c.deepScopeSeen) > 0 {
		clear(c.deepScopeSeen)
	}
}

// EnterScope prepares state for the new scope.
//
// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.enterScope(NodeTraversal)
func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) EnterScope(t *NodeTraversal) {
	if t.InGlobalScope() {
		// Don't track any information for the global scope, because it is likely to be
		// very large, and we will never rename any variables it contains. MakeDeclaredNamesUnique
		// favors keeping global variables unchanged, so they won't need to be restored.
		return
	}
	if len(c.scopes) == 0 && t.GetScopeRoot().IsModuleBody() &&
		!renameInverterSubtreeHasSeparator(t.GetScopeRoot()) {
		c.skippingModuleBody = true
		return
	}

	scope := t.GetScope()
	c.pushScope(int32(scope.GetDepth()), scope.varList)
}

// ExitScope renames vars for the current scope, and propagates referenced
// outer variables into the parent scope reference list.
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

	c.popScope()
	if len(c.scopes) == 0 {
		c.resetArena()
	}
}

// Java: MakeDeclaredNamesUnique.ContextualRenameInverter.shouldTraverse(NodeTraversal, Node, Node)
func (c *MakeDeclaredNamesUnique_ContextualRenameInverter) ShouldTraverse(
	t *NodeTraversal, n rhino.Node, parent rhino.Node,
) bool {
	if c.skippingModuleBody {
		return false
	}
	if len(c.scopes) == 0 && !parent.IsNil() {
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
	if len(c.scopes) == 0 {
		return
	}
	tok := node.GetToken()
	if tok == rhino.Token_NAME || tok == rhino.Token_IMPORT_STAR {
		id := node.GetStringID()
		if tok == rhino.Token_NAME && id <= rhino.RhinoStringPool_EmptyID {
			return
		}
		varID := c.getOrCreateCurrentVar(id, node)
		c.addReferenceNode(varID, node)
		c.addScopeRef(len(c.scopes)-1, varID)
	}
}
