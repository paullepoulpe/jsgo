// Port of com/google/javascript/jscomp/Reference.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/BasicBlock.java @ closure-compiler v20261021

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: declaration and l-value queries read parent/child roles through
// package ast while keeping Reference's rhino.Node API intact.

// Reference represents a single declaration or reference to a variable.
// Note that references can only be used with untyped scopes and traversals.
type Reference struct {
	nameNode   rhino.Node
	basicBlock *BasicBlock
	scope      *Scope
	inputId    *rhino.InputId
}

// Java: Reference(Node, NodeTraversal, BasicBlock)
func NewReference(nameNode rhino.Node, t *NodeTraversal, basicBlock *BasicBlock) *Reference {
	return newReferenceWithScopeAndInputId(nameNode, basicBlock, t.GetScope(), t.GetInput().GetInputId())
}

// Java: Reference(Node, BasicBlock, Scope, InputId)
func newReferenceWithScopeAndInputId(
	nameNode rhino.Node, basicBlock *BasicBlock, scope *Scope, inputId *rhino.InputId) *Reference {
	return &Reference{
		nameNode:   nameNode,
		basicBlock: basicBlock,
		scope:      scope,
		inputId:    inputId,
	}
}

// Java: Reference.toString()
func (r *Reference) String() string {
	return r.nameNode.String()
}

// Reference_CreateRefForTest creates a variable reference in a given script
// file name, used in tests.
//
// Java: Reference.createRefForTest(CompilerInput)
func Reference_CreateRefForTest(input *CompilerInput) *Reference {
	return newReferenceWithScopeAndInputId(rhino.NewNode(rhino.Token_NAME), nil, nil, input.GetInputId())
}

// CloneWithNewScope makes a copy of the current reference using a new
// Scope instance.
//
// Java: Reference.cloneWithNewScope(Scope)
func (r *Reference) CloneWithNewScope(newScope *Scope) *Reference {
	return newReferenceWithScopeAndInputId(r.nameNode, r.basicBlock, newScope, r.inputId)
}

// GetSymbol returns Var in Java, narrowing rhino.StaticRef.getSymbol();
// it keeps rhino's signature and never returns a typed nil (assert
// `.(*Var)`; see the note at the top of abstract_scope.go).
//
// Java: Reference.getSymbol()
func (r *Reference) GetSymbol() rhino.StaticSlot {
	if v := r.getSymbolV(); v != nil {
		return v
	}
	return nil
}

// getSymbolV is getSymbol with its Java return type.
//
// Go-only: language-adapter — see GetSymbol.
func (r *Reference) getSymbolV() *Var {
	return r.scope.GetVar(r.nameNode.GetString())
}

// Java: Reference.getNode()
func (r *Reference) GetNode() rhino.Node {
	return r.nameNode
}

// Java: Reference.getInputId()
func (r *Reference) GetInputId() *rhino.InputId {
	return r.inputId
}

// Java: Reference.getSourceFile()
func (r *Reference) GetSourceFile() rhino.StaticSourceFile {
	return r.nameNode.GetStaticSourceFile()
}

// Java: Reference.isDeclaration()
func (r *Reference) IsDeclaration() bool {
	return reference_isDeclarationHelper(r.nameNode)
}

// Java: Reference.isDeclarationHelper(Node)
func reference_isDeclarationHelper(node rhino.Node) bool {
	parent := node.GetParent()
	parentToken := parent.GetToken()

	// Special case for class B extends A, A is not a declaration.
	if cls, ok := ast.AsClass(parent); ok && cls.NameNode() != node {
		return false
	}

	grandparent := parent.GetParent()
	// This condition can be true during InlineVariables.
	if grandparent.IsNil() {
		return false
	}

	if _, isDecl := ast.AsNameDeclaration(grandparent); isDecl &&
		ast.RoleIn(node, parent) == ast.RoleDestructuringLhsInit {
		// This is the RHS of a var/let/const and thus not a declaration.
		return false
	}

	switch parentToken {
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST,
		rhino.Token_PARAM_LIST, rhino.Token_CLASS, rhino.Token_CATCH,
		rhino.Token_IMPORT:
		return true
	case rhino.Token_FUNCTION:
		fn := ast.MustFunction(parent)
		if fn.IsArrow() {
			return fn.NameNode().Is(node)
		}
		return true
	case rhino.Token_DESTRUCTURING_LHS, rhino.Token_OBJECT_PATTERN, rhino.Token_ARRAY_PATTERN,
		rhino.Token_ITER_REST, rhino.Token_OBJECT_REST:
		return reference_isDeclarationHelper(parent)
	case rhino.Token_STRING_KEY:
		if _, ok := ast.AsObjectPattern(grandparent); ok {
			return reference_isDeclarationHelper(parent)
		}
		return false
	case rhino.Token_COMPUTED_PROP:
		if _, ok := ast.AsObjectPattern(grandparent); ok {
			if val, hasVal := ast.MustComputedProp(parent).Value(); hasVal && val.Is(node) {
				return reference_isDeclarationHelper(parent)
			}
		}
		return false
	case rhino.Token_DEFAULT_VALUE:
		if ast.MustDefaultValue(parent).Target().Is(node) {
			return reference_isDeclarationHelper(parent)
		}
		return false
	case rhino.Token_IMPORT_SPEC:
		return ast.MustImportSpec(parent).Local().Is(node)
	default:
		return false
	}
}

// Java: Reference.isVarDeclaration()
func (r *Reference) IsVarDeclaration() bool {
	return r.GetParent().IsVar()
}

// Java: Reference.isLetDeclaration()
func (r *Reference) IsLetDeclaration() bool {
	return r.GetParent().IsLet()
}

// Java: Reference.isDotPropertyAccess()
func (r *Reference) IsDotPropertyAccess() bool {
	return r.GetParent().IsGetProp()
}

// Java: Reference.isAssignedToObjectDestructuringPattern()
func (r *Reference) IsAssignedToObjectDestructuringPattern() bool {
	parent := r.GetParent()
	// `let { propName: varName } = ref;`
	// or
	// `({ propName: varName } = ref);`
	if lhs, ok := ast.AsDestructuringLhs(parent); ok {
		_, isObj := lhs.Pattern().AsObjectPattern()
		return isObj
	}
	if assign, ok := ast.AsAssign(parent); ok {
		_, isObj := assign.Target().AsObjectPattern()
		return isObj
	}
	return false
}

// Java: Reference.isConstDeclaration()
func (r *Reference) IsConstDeclaration() bool {
	return r.GetParent().IsConst()
}

// Java: Reference.isHoistedFunction()
func (r *Reference) IsHoistedFunction() bool {
	return NodeUtil_IsHoistedFunctionDeclaration(r.GetParent())
}

// IsInitializingDeclaration determines whether the variable is initialized
// at the declaration.
//
// Java: Reference.isInitializingDeclaration()
func (r *Reference) IsInitializingDeclaration() bool {
	// VAR and LET are the only types of variable declarations that may not initialize
	// their variables. Catch blocks, named functions, and parameters all do.
	if name, ok := ast.AsName(r.nameNode); ok {
		if _, hasInit := name.Init(); hasInit {
			return true
		}
	}
	parent := r.GetParent()
	return r.IsDeclaration() && !parent.IsVar() && !parent.IsLet()
}

// GetAssignedValue returns, for an assignment, variable declaration, or
// function declaration, the assigned value, otherwise null.
//
// Java: Reference.getAssignedValue()
func (r *Reference) GetAssignedValue() rhino.Node {
	return NodeUtil_GetRValueOfLValue(r.nameNode)
}

// Java: Reference.getBasicBlock()
func (r *Reference) GetBasicBlock() *BasicBlock {
	return r.basicBlock
}

// Java: Reference.getParent()
func (r *Reference) GetParent() rhino.Node {
	return r.GetNode().GetParent()
}

// Java: Reference.getGrandparent()
func (r *Reference) GetGrandparent() rhino.Node {
	return r.GetNode().GetGrandparent()
}

// Java: Reference.isLhsOfEnhancedForExpression(Node)
func reference_isLhsOfEnhancedForExpression(n rhino.Node) bool {
	ef, ok := ast.AsEnhancedFor(n.GetParent())
	return ok && ef.Target().Is(n)
}

// Java: Reference.isSimpleAssignmentToName()
func (r *Reference) IsSimpleAssignmentToName() bool {
	assign, ok := ast.AsAssign(r.GetParent())
	return ok && assign.Target().Is(r.nameNode)
}

// IsLvalue returns whether the name node for this reference is an lvalue.
// TODO(tbreisacher): This method disagrees with NodeUtil#isLValue for "var
// x;" and "let x;". Consider updating it to match.
//
// Java: Reference.isLvalue()
func (r *Reference) IsLvalue() bool {
	parent := r.GetParent()
	switch parent.GetToken() {
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		if name, ok := ast.AsName(r.nameNode); ok {
			if _, hasInit := name.Init(); hasInit {
				return true
			}
		}
		return reference_isLhsOfEnhancedForExpression(r.nameNode)
	case rhino.Token_DEFAULT_VALUE:
		return ast.MustDefaultValue(parent).Target().Is(r.nameNode)
	case rhino.Token_INC, rhino.Token_DEC, rhino.Token_CATCH, rhino.Token_ITER_REST, rhino.Token_OBJECT_REST,
		rhino.Token_PARAM_LIST:
		return true
	case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		ef, ok := ast.AsEnhancedFor(parent)
		return ok && ef.Target().Is(r.nameNode)
	case rhino.Token_ARRAY_PATTERN, rhino.Token_STRING_KEY, rhino.Token_COMPUTED_PROP:
		return NodeUtil_IsLhsByDestructuring(r.nameNode)
	default:
		assign, ok := ast.AsAssignment(parent)
		return ok && assign.Target().Is(r.nameNode)
	}
}

// Java: Reference.getScope()
func (r *Reference) GetScope() *Scope {
	return r.scope
}

var nextBasicBlockRootID atomic.Uint32

// BasicBlock represents a section of code that is uninterrupted by control
// structures (conditional or iterative logic).
type BasicBlock struct {
	parent *BasicBlock

	root   rhino.Node
	depth  int32
	rootID uint32

	// Whether this block denotes a function scope.
	isFunction bool

	// Whether this block denotes a loop.
	isLoop bool

	hasFnAncestor   bool
	hasLoopBeforeFn bool
}

// NewBasicBlock creates a new block.
//
// parent is the containing block; root is the root node of the block.
//
// Java: BasicBlock(BasicBlock, Node)
func NewBasicBlock(parent *BasicBlock, root rhino.Node) *BasicBlock {
	return initBasicBlock(&BasicBlock{}, parent, root)
}

func initBasicBlock(b *BasicBlock, parent *BasicBlock, root rhino.Node) *BasicBlock {
	b.parent = parent
	b.root = root
	if parent != nil {
		b.depth = parent.depth + 1
		b.rootID = parent.rootID
	} else {
		b.depth = 0
		b.rootID = nextBasicBlockRootID.Add(1)
	}

	b.isFunction = root.IsFunction()

	if p := root.GetParent(); !p.IsNil() {
		pType := p.GetToken()
		b.isLoop = pType == rhino.Token_DO ||
			pType == rhino.Token_WHILE ||
			pType == rhino.Token_FOR ||
			pType == rhino.Token_FOR_OF ||
			pType == rhino.Token_FOR_AWAIT_OF ||
			pType == rhino.Token_FOR_IN
	} else {
		b.isLoop = false
	}

	if b.isFunction {
		b.hasFnAncestor = true
		b.hasLoopBeforeFn = false
	} else if parent != nil {
		b.hasFnAncestor = parent.hasFnAncestor
		b.hasLoopBeforeFn = b.isLoop || parent.hasLoopBeforeFn
	} else {
		b.hasFnAncestor = false
		b.hasLoopBeforeFn = b.isLoop
	}
	return b
}

// Java: BasicBlock.getParent()
func (b *BasicBlock) GetParent() *BasicBlock {
	return b.parent
}

// IsGlobalScopeBlock determines whether this block is equivalent to the
// very first block that is created when reference collection traversal
// enters global scope. Note that when traversing a single script in a
// hot-swap fashion a new instance of BasicBlock is created.
//
// Java: BasicBlock.isGlobalScopeBlock()
func (b *BasicBlock) IsGlobalScopeBlock() bool {
	return b.parent == nil
}

// ProvablyExecutesBefore determines whether this block is guaranteed to
// begin executing before the given block does.
//
// Java: BasicBlock.provablyExecutesBefore(BasicBlock)
func (b *BasicBlock) ProvablyExecutesBefore(thatBlock *BasicBlock) bool {
	if thatBlock == b {
		return true
	}
	if b.parent == nil {
		return thatBlock.parent == nil || thatBlock.rootID == b.rootID
	}
	// If thatBlock is a descendant of this block, and there are no hoisted
	// blocks between them, then this block must start before thatBlock.
	bDepth := b.depth
	if thatBlock.depth <= bDepth || thatBlock.rootID != b.rootID {
		return false
	}
	currentBlock := thatBlock.parent
	for currentBlock.depth > bDepth {
		currentBlock = currentBlock.parent
	}
	return currentBlock == b
}

// Java: BasicBlock.isFunction()
func (b *BasicBlock) IsFunction() bool {
	return b.isFunction
}

// Java: BasicBlock.isLoop()
func (b *BasicBlock) IsLoop() bool {
	return b.isLoop
}

// Java: BasicBlock.getRoot()
func (b *BasicBlock) GetRoot() rhino.Node {
	return b.root
}

// Java: BasicBlock.toString()
func (b *BasicBlock) String() string {
	return "BasicBlock @ " + b.root.String()
}
