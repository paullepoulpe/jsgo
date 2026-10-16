// Port of com/google/javascript/jscomp/Reference.java @ closure-compiler v20260415

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Reference represents a single declaration or reference to a variable.
// Note that references can only be used with untyped scopes and traversals.
type Reference struct {
	nameNode   *rhino.Node
	basicBlock *BasicBlock
	scope      *Scope
	inputId    *rhino.InputId
}

// Favor EnumSet over ImmutableSet for performance
var reference_DECLARATION_PARENTS = jutil.NewEnumSet(
	rhino.Token_VAR,
	rhino.Token_LET,
	rhino.Token_CONST,
	rhino.Token_PARAM_LIST,
	rhino.Token_FUNCTION,
	rhino.Token_CLASS,
	rhino.Token_CATCH)

// Java: Reference(Node, NodeTraversal, BasicBlock)
func NewReference(nameNode *rhino.Node, t *NodeTraversal, basicBlock *BasicBlock) *Reference {
	return newReferenceWithScopeAndInputId(nameNode, basicBlock, t.GetScope(), t.GetInput().GetInputId())
}

// Java: Reference(Node, BasicBlock, Scope, InputId)
func newReferenceWithScopeAndInputId(
	nameNode *rhino.Node, basicBlock *BasicBlock, scope *Scope, inputId *rhino.InputId) *Reference {
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
func (r *Reference) GetNode() *rhino.Node {
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
func reference_isDeclarationHelper(node *rhino.Node) bool {
	parent := node.GetParent()

	// Special case for class B extends A, A is not a declaration.
	if parent.IsClass() && !node.IsFirstChildOf(parent) {
		return false
	}

	// This condition can be true during InlineVariables.
	if parent.GetParent() == nil {
		return false
	}

	if NodeUtil_IsNameDeclaration(parent.GetParent()) && node == parent.GetSecondChild() {
		// This is the RHS of a var/let/const and thus not a declaration.
		return false
	}

	// Special cases for destructuring patterns.
	if parent.IsDestructuringLhs() ||
		parent.IsDestructuringPattern() ||
		parent.IsRest() ||
		(parent.IsStringKey() && parent.GetParent().IsObjectPattern()) ||
		(parent.IsComputedProp() &&
			parent.GetParent().IsObjectPattern() &&
			node == parent.GetLastChild()) ||
		(parent.IsDefaultValue() && node == parent.GetFirstChild()) {
		return reference_isDeclarationHelper(parent)
	}

	if parent.IsImport() {
		return true
	}

	if parent.IsImportSpec() && node == parent.GetLastChild() {
		return true
	}

	// Special case for arrow function
	if parent.IsArrowFunction() {
		return node == parent.GetFirstChild()
	}

	return reference_DECLARATION_PARENTS.Contains(parent.GetToken())
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
	return (parent.IsDestructuringLhs() || parent.IsAssign()) &&
		parent.GetFirstChild().IsObjectPattern()
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
	return (r.IsDeclaration() && !r.GetParent().IsVar() && !r.GetParent().IsLet()) ||
		r.nameNode.HasChildren()
}

// GetAssignedValue returns, for an assignment, variable declaration, or
// function declaration, the assigned value, otherwise null.
//
// Java: Reference.getAssignedValue()
func (r *Reference) GetAssignedValue() *rhino.Node {
	return NodeUtil_GetRValueOfLValue(r.nameNode)
}

// Java: Reference.getBasicBlock()
func (r *Reference) GetBasicBlock() *BasicBlock {
	return r.basicBlock
}

// Java: Reference.getParent()
func (r *Reference) GetParent() *rhino.Node {
	return r.GetNode().GetParent()
}

// Java: Reference.getGrandparent()
func (r *Reference) GetGrandparent() *rhino.Node {
	return r.GetNode().GetGrandparent()
}

// Java: Reference.isLhsOfEnhancedForExpression(Node)
func reference_isLhsOfEnhancedForExpression(n *rhino.Node) bool {
	parent := n.GetParent()
	return NodeUtil_IsEnhancedFor(parent) && parent.GetFirstChild() == n
}

// Java: Reference.isSimpleAssignmentToName()
func (r *Reference) IsSimpleAssignmentToName() bool {
	parent := r.GetParent()
	return parent.IsAssign() && parent.GetFirstChild() == r.nameNode
}

// IsLvalue returns whether the name node for this reference is an lvalue.
// TODO(tbreisacher): This method disagrees with NodeUtil#isLValue for "var
// x;" and "let x;". Consider updating it to match.
//
// Java: Reference.isLvalue()
func (r *Reference) IsLvalue() bool {
	parent := r.GetParent()
	parentType := parent.GetToken()
	switch parentType {
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		return r.nameNode.HasChildren() || reference_isLhsOfEnhancedForExpression(r.nameNode)
	case rhino.Token_DEFAULT_VALUE:
		return parent.GetFirstChild() == r.nameNode
	case rhino.Token_INC, rhino.Token_DEC, rhino.Token_CATCH, rhino.Token_ITER_REST, rhino.Token_OBJECT_REST,
		rhino.Token_PARAM_LIST:
		return true
	case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		return NodeUtil_IsEnhancedFor(parent) && parent.GetFirstChild() == r.nameNode
	case rhino.Token_ARRAY_PATTERN, rhino.Token_STRING_KEY, rhino.Token_COMPUTED_PROP:
		return NodeUtil_IsLhsByDestructuring(r.nameNode)
	default:
		return NodeUtil_IsAssignmentOp(parent) && parent.GetFirstChild() == r.nameNode
	}
}

// Java: Reference.getScope()
func (r *Reference) GetScope() *Scope {
	return r.scope
}
