// Port of com/google/javascript/jscomp/Scope.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Scope contains information about a variable scope in JavaScript. Scopes
// can be nested, a scope points back to its parent scope. A Scope contains
// information about variables defined in that scope.
//
// Most methods come from AbstractScopeBase[*Scope, *Var] and are typed:
// GetVar returns *Var, GetParent and GetClosestHoistScope return *Scope,
// GetVarIterable is an iter.Seq[*Var]. GetSlot, GetOwnSlot and
// GetParentScope keep rhino.StaticScope's signatures; see the note at the
// top of abstract_scope.go.
//
// See NodeTraversal.
type Scope struct {
	AbstractScopeBase[*Scope, *Var]

	parent *Scope
	arena  *scopeArena
}

// scopeArena slab-allocates Scope, Var, and LinkedHashMap slot slices for a
// SyntacticScopeCreator and its descendant scopes.
type scopeArena struct {
	scopeChunk    []Scope
	scopeUsed     int
	varChunk      []Var
	varUsed       int
	slotChunk     []jutil.LhmSlot[string, *Var]
	slotUsed      int
	implicitChunk [][4]*Var
	implicitUsed  int
}

func (a *scopeArena) allocScope() *Scope {
	if a.scopeUsed == len(a.scopeChunk) {
		n := 4 * len(a.scopeChunk)
		if n < 8 {
			n = 8
		} else if n > 256 {
			n = 256
		}
		a.scopeChunk = make([]Scope, n)
		a.scopeUsed = 0
	}
	s := &a.scopeChunk[a.scopeUsed]
	a.scopeUsed++
	return s
}

func (a *scopeArena) allocVar() *Var {
	if a.varUsed == len(a.varChunk) {
		n := 4 * len(a.varChunk)
		if n < 16 {
			n = 16
		} else if n > 512 {
			n = 512
		}
		a.varChunk = make([]Var, n)
		a.varUsed = 0
	}
	v := &a.varChunk[a.varUsed]
	a.varUsed++
	return v
}

func (a *scopeArena) allocVarSlots(n int) []jutil.LhmSlot[string, *Var] {
	if len(a.slotChunk)-a.slotUsed < n {
		size := 4 * len(a.slotChunk)
		if size < 32 {
			size = 32
		} else if size > 512 {
			size = 512
		}
		if size < n {
			size = n
		}
		a.slotChunk = make([]jutil.LhmSlot[string, *Var], size)
		a.slotUsed = 0
	}
	u := a.slotUsed
	a.slotUsed = u + n
	return a.slotChunk[u : u+n : u+n]
}

func (a *scopeArena) allocImplicitVars() *[4]*Var {
	if a.implicitUsed == len(a.implicitChunk) {
		n := 4 * len(a.implicitChunk)
		if n < 8 {
			n = 8
		} else if n > 128 {
			n = 128
		}
		a.implicitChunk = make([][4]*Var, n)
		a.implicitUsed = 0
	}
	p := &a.implicitChunk[a.implicitUsed]
	a.implicitUsed++
	return p
}

// Java: Scope.createGlobalScope(Node)
func Scope_CreateGlobalScope(rootNode *rhino.Node) *Scope {
	return newScope(rootNode)
}

// Java: Scope.createChildScope(Scope, Node)
func Scope_CreateChildScope(parent *Scope, rootNode *rhino.Node) *Scope {
	return newScopeWithParent(parent, rootNode)
}

// Java: Scope(Scope, Node)
func newScopeWithParent(parent *Scope, rootNode *rhino.Node) *Scope {
	var s *Scope
	if parent != nil && parent.arena != nil {
		s = parent.arena.allocScope()
		s.arena = parent.arena
	} else {
		s = &Scope{}
	}
	initScopeWithParent(s, parent, rootNode)
	return s
}

func newScopeWithParentInArena(parent *Scope, rootNode *rhino.Node, arena *scopeArena) *Scope {
	s := arena.allocScope()
	s.arena = arena
	initScopeWithParent(s, parent, rootNode)
	return s
}

func initScopeWithParent(s *Scope, parent *Scope, rootNode *rhino.Node) {
	if rootNode == nil || parent == nil {
		panic(jutil.NewNullPointerException())
	}
	if !NodeUtil_CreatesScope(rootNode) {
		panic(jutil.NewIllegalArgumentExceptionMsg(jutil.LenientFormat("%s", rootNode)))
	}
	if rootNode == parent.rootNode {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			jutil.LenientFormat("rootNode should not be the parent's root node: %s", rootNode)))
	}
	s.self = s
	s.rootNode = rootNode
	s.parent = parent
	s.depth = parent.depth + 1
}

// Java: Scope(Node)
func newScope(rootNode *rhino.Node) *Scope {
	s := &Scope{}
	initScope(s, rootNode)
	return s
}

func newScopeInArena(rootNode *rhino.Node, arena *scopeArena) *Scope {
	s := arena.allocScope()
	s.arena = arena
	initScope(s, rootNode)
	return s
}

func initScope(s *Scope, rootNode *rhino.Node) {
	if rootNode == nil {
		panic(jutil.NewNullPointerException())
	}
	if !NodeUtil_CreatesScope(rootNode) && !rootNode.IsScript() && !rootNode.IsRoot() {
		panic(jutil.NewIllegalArgumentExceptionMsg(jutil.LenientFormat("%s", rootNode)))
	}
	s.self = s
	s.rootNode = rootNode
}

// Java: Scope.untyped()
func (s *Scope) Untyped() *Scope {
	return s
}

// Java: Scope.getParent()
func (s *Scope) GetParent() *Scope {
	return s.parent
}

// Java: AbstractScope.getParentScope()
func (s *Scope) GetParentScope() rhino.StaticScope {
	if s.parent != nil {
		return s.parent
	}
	return nil
}

// Java: AbstractScope.isGlobal()
func (s *Scope) IsGlobal() bool {
	return s.parent == nil
}

// Java: AbstractScope.isLocal()
func (s *Scope) IsLocal() bool {
	return s.parent != nil
}

// Java: AbstractScope.getGlobalScope()
func (s *Scope) GetGlobalScope() *Scope {
	for s.parent != nil {
		s = s.parent
	}
	return s
}

// Java: AbstractScope.isHoistScope()
func (s *Scope) IsHoistScope() bool {
	return s.parent == nil ||
		s.rootNode.IsFunction() ||
		NodeUtil_IsFunctionBlock(s.rootNode) ||
		s.rootNode.IsModuleBody() ||
		NodeUtil_IsClassStaticBlock(s.rootNode)
}

// Java: AbstractScope.getClosestHoistScope()
func (s *Scope) GetClosestHoistScope() *Scope {
	for ; s != nil; s = s.parent {
		if s.IsHoistScope() {
			return s
		}
	}
	return nil
}

// Java: AbstractScope.getClosestContainerScope()
func (s *Scope) GetClosestContainerScope() *Scope {
	scope := s.GetClosestHoistScope()
	if NodeUtil_IsFunctionBlock(scope.rootNode) {
		scope = scope.parent
		if scope.IsBlockScope() {
			panic(jutil.NewIllegalStateExceptionMsg(jutil.LenientFormat("%s", scope)))
		}
	}
	return scope
}

func (s *Scope) getOwnSlotV(name string) *Var {
	if s.vars.Size() > 0 {
		if var_ := s.vars.Get(name); var_ != nil {
			return var_
		}
	}
	if imp := AbstractScope_ImplicitVar_Of(name); imp != nil {
		return s.GetOwnImplicitSlot(imp)
	}
	return nil
}

// Java: AbstractScope.hasOwnSlot(String)
func (s *Scope) HasOwnSlot(name string) bool {
	if s.vars.Size() > 0 && s.vars.ContainsKey(name) {
		return true
	}
	if imp := AbstractScope_ImplicitVar_Of(name); imp != nil {
		return s.HasOwnImplicitSlot(imp)
	}
	return false
}

// Java: AbstractScope.hasSlot(String)
func (s *Scope) HasSlot(name string) bool {
	implicit := AbstractScope_ImplicitVar_Of(name)
	for ; s != nil; s = s.parent {
		if s.vars.Size() > 0 && s.vars.ContainsKey(name) {
			return true
		}
		if implicit != nil && implicit.IsMadeByScope(s) {
			return true
		}
	}
	return false
}

func (s *Scope) inheritedInput() (*CompilerInput, bool) {
	for anc := s.parent; anc != nil && anc.parent != nil; anc = anc.parent {
		if anc.vars.Size() > 0 {
			for v := range anc.vars.Values() {
				if !v.IsImplicitGoogNamespace() {
					return v.GetInput(), true
				}
			}
		}
	}
	return nil, false
}

// GetVar uses the concrete parent links of Java's final Scope class. The
// generic base must support TypedScope's virtual methods; untyped scopes can
// avoid those dictionary calls and classify implicit names once per lookup.
func (s *Scope) GetVar(name string) *Var {
	implicit := AbstractScope_ImplicitVar_Of(name)
	for {
		if s.vars.Size() > 0 {
			if v := s.vars.Get(name); v != nil {
				return v
			}
		}
		if implicit != nil {
			if v := s.GetOwnImplicitSlot(implicit); v != nil {
				return v
			}
		}
		s = s.parent
		if s == nil {
			return nil
		}
	}
}

// Java: AbstractScope.hasOwnImplicitSlot(ImplicitVar)
func (s *Scope) HasOwnImplicitSlot(name *AbstractScope_ImplicitVar) bool {
	return name != nil && name.IsMadeByScope(s)
}

// Java: AbstractScope.getOwnImplicitSlot(ImplicitVar)
func (s *Scope) GetOwnImplicitSlot(name *AbstractScope_ImplicitVar) *Var {
	if !s.HasOwnImplicitSlot(name) {
		return nil
	}
	if s.implicitVars == nil {
		if s.arena != nil {
			s.implicitVars = s.arena.allocImplicitVars()
		} else {
			s.implicitVars = new([4]*Var)
		}
	}
	if v := s.implicitVars[*name]; v != nil {
		return v
	}
	v := s.MakeImplicitVar(*name)
	if v != nil {
		s.implicitVars[*name] = v
	}
	return v
}

// Java: AbstractScope.getSlot(String) (final)
func (s *Scope) GetSlot(name string) rhino.StaticSlot {
	if v := s.GetVar(name); v != nil {
		return v
	}
	return nil
}

// Declare declares a variable.
//
// name is the name of the variable, nameNode the NAME node declaring the
// variable, input the input in which this variable is defined.
//
// Java: Scope.declare(String, Node, CompilerInput)
func (s *Scope) Declare(name string, nameNode *rhino.Node, input *CompilerInput) *Var {
	jutil.CheckArgument(name != "")
	// Make sure that it's declared only once
	jutil.CheckState(s.getOwnSlotV(name) == nil)
	var var_ *Var
	if s.arena != nil {
		var_ = s.arena.allocVar()
		initVar(var_, name, nameNode, s, s.GetVarCount(), input, nil)
	} else {
		var_ = NewVar(
			name,
			nameNode,
			s,
			s.GetVarCount(),
			input,
			/* implicitGoogNamespaceDefinition= */ nil)
	}
	s.DeclareInternal(name, var_)
	return var_
}

func functionScopeVarSlotCount(fn *rhino.Node) int {
	fnNameNode := fn.GetFirstChild()
	n := 0
	if fnNameNode.GetString() != "" && NodeUtil_IsFunctionExpression(fn) {
		n = 1
	}
	for c := fnNameNode.GetNext().GetFirstChild(); c != nil; c = c.GetNext() {
		switch c.GetToken() {
		case rhino.Token_NAME:
			n++
		case rhino.Token_DEFAULT_VALUE, rhino.Token_ITER_REST:
			if c.GetFirstChild().IsName() {
				n++
			} else {
				n += 4
			}
		default:
			n += 4
		}
	}
	if n < 1 {
		return 1
	}
	if n > 64 {
		return 64
	}
	return n
}

func functionBlockScopeVarSlotCount(block *rhino.Node) int {
	n := 0
	steps := 0
	for c := block.GetFirstChild(); c != nil; c = c.GetNext() {
		steps++
		if steps > 32 {
			if n < 16 {
				return 16
			}
			return n
		}
		switch c.GetToken() {
		case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
			for lhs := c.GetFirstChild(); lhs != nil; lhs = lhs.GetNext() {
				if lhs.IsName() {
					n++
				} else {
					n += 2
				}
			}
		case rhino.Token_FUNCTION, rhino.Token_CLASS:
			n++
		case rhino.Token_EXPR_RESULT, rhino.Token_RETURN, rhino.Token_THROW,
			rhino.Token_BREAK, rhino.Token_CONTINUE, rhino.Token_EMPTY, rhino.Token_DEBUGGER:
		default:
			n++
		}
		if n >= 32 {
			return 32
		}
	}
	if n < 1 {
		return 1
	}
	return n
}

const maxPreallocScopeVars = 512

func countLhsVars(n *rhino.Node) int {
	if n.HasOneChild() && n.GetFirstChild().IsName() {
		return 1
	}
	c := 0
	for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
		if child.IsName() {
			c++
		} else {
			c += 2
		}
	}
	if c == 0 {
		return 1
	}
	return c
}

func countScopeVars(n *rhino.Node, isHoist bool, inBlockScope bool, root *rhino.Node, count int, steps *int) int {
	if count >= maxPreallocScopeVars {
		return maxPreallocScopeVars
	}
	*steps++
	if *steps > 256 {
		if count < 16 {
			return 32
		}
		if count*2 > maxPreallocScopeVars {
			return maxPreallocScopeVars
		}
		return count * 2
	}
	switch n.GetToken() {
	case rhino.Token_VAR, rhino.Token_IMPORT:
		if isHoist {
			count += countLhsVars(n)
		}
		return count
	case rhino.Token_LET, rhino.Token_CONST:
		if inBlockScope {
			count += countLhsVars(n)
		}
		return count
	case rhino.Token_EXPORT:
		if first := n.GetFirstChild(); first != nil {
			return countScopeVars(first, isHoist, inBlockScope, root, count, steps)
		}
		return count
	case rhino.Token_FUNCTION, rhino.Token_CLASS:
		if inBlockScope && !NodeUtil_IsFunctionExpression(n) && !NodeUtil_IsClassExpression(n) {
			count++
		}
		return count
	case rhino.Token_CATCH:
		if inBlockScope {
			count++
		}
		return countScopeVars(n.GetSecondChild(), isHoist, inBlockScope, root, count, steps)
	case rhino.Token_MODULE_BODY:
		if n != root {
			return count + 1
		}
	case rhino.Token_BLOCK:
		if n != root && scanVarsBlockCreatesScope(n) {
			if !isHoist {
				return count
			}
			inBlockScope = false
		}
	case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF, rhino.Token_SWITCH_BODY:
		if n != root {
			if !isHoist {
				return count
			}
			inBlockScope = false
		}
	case rhino.Token_SCRIPT, rhino.Token_ROOT,
		rhino.Token_DO, rhino.Token_WHILE, rhino.Token_WITH, rhino.Token_IF,
		rhino.Token_LABEL, rhino.Token_TRY, rhino.Token_SWITCH, rhino.Token_CASE, rhino.Token_DEFAULT_CASE:
	default:
		return count
	}
	for child := n.GetFirstChild(); child != nil; child = child.GetNext() {
		if scanVarsTokenKind[uint8(child.GetToken())] == 1 {
			count = countScopeVars(child, isHoist, inBlockScope, root, count, steps)
			if count >= maxPreallocScopeVars || *steps > 256 {
				return count
			}
		}
	}
	return count
}

func largeScopeVarSlotCount(sc *Scope) int {
	root := sc.rootNode
	if root.IsFunction() {
		n := functionScopeVarSlotCount(root)
		if n < 12 {
			return 12
		}
		return n
	}
	isHoist := sc.IsHoistScope()
	steps := 0
	n := countScopeVars(root, isHoist, true, root, 0, &steps)
	if n < 16 {
		return 16
	}
	return n
}

// declareNewVar declares a variable when the caller
// (SyntacticScopeCreator_ScopeScanner.DeclareVar) has already verified that
// s.getOwnSlotV(name) == nil.
func (s *Scope) declareNewVar(name string, nameNode *rhino.Node, input *CompilerInput) {
	jutil.CheckArgument(name != "")
	var var_ *Var
	if s.arena != nil {
		var_ = s.arena.allocVar()
		initVar(var_, name, nameNode, s, s.vars.Size(), input, nil)
	} else {
		var_ = NewVar(name, nameNode, s, s.vars.Size(), input, nil)
	}
	if s.IsFunctionBlockScope() && s.parent.HasOwnSlot(name) && !s.IsBleedingFunctionName(name) {
		panic(jutil.NewIllegalStateExceptionMsg(jutil.LenientFormat("Illegal shadow: %s", nameNode)))
	}
	if !s.vars.HasBacking() && s.arena != nil && s.parent != nil {
		n := 2
		if s.rootNode.IsFunction() {
			n = functionScopeVarSlotCount(s.rootNode)
		} else if s.parent.rootNode.IsFunction() && s.rootNode.IsBlock() {
			n = functionBlockScopeVarSlotCount(s.rootNode)
		}
		s.vars.InitWithBuffer(s.arena.allocVarSlots(n))
	}
	s.varsInitialized = true
	s.vars.PutNew(name, var_)
}

// DeclareImplicitGoogNamespaceIfAbsent declares an implicit goog.provide or
// goog.module namespace in this scope.
//
// name is the name of the variable, definition the STRINGLIT node holding
// the full namespace.
//
// Java: Scope.declareImplicitGoogNamespaceIfAbsent(String, Node)
func (s *Scope) DeclareImplicitGoogNamespaceIfAbsent(name string, definition *rhino.Node) *Var {
	jutil.CheckArgument(name != "")
	jutil.CheckStateMsg1(s.IsGlobal(), "Cannot declare implicit goog namespace in local scope %s", s)
	// Allow redeclarations of provides, since they are implicit and don't have a single
	// declaration site.
	var_ := s.getOwnSlotV(name)
	if var_ == nil {
		var_ = Var_CreateImplicitGoogNamespace(name, s, definition)
		s.DeclareInternal(name, var_)
	} else if var_.IsImplicitGoogNamespace() {
		var_.AddImplicitGoogNamespaceDefinition(definition)
	}
	return var_
}

// Java: Scope.makeImplicitVar(ImplicitVar)
func (s *Scope) MakeImplicitVar(var_ AbstractScope_ImplicitVar) *Var {
	if s.arena != nil {
		v := s.arena.allocVar()
		initVar(v, var_.name(), nil, s, -1, nil, nil)
		return v
	}
	return NewVar(
		var_.name(),
		/* nameNode= */ nil,
		s,
		/* index= */ -1,
		/* input= */ nil,
		/* implicitGoogNamespaceDefinition= */ nil)
}
