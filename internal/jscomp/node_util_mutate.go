// Port of com/google/javascript/jscomp/NodeUtil.java @ closure-compiler v20261021

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
	"math"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// NodeUtil_IteratesImpureIterable returns true if node might execute an
// `Iterable` iteration that has side-effects, false if there are definitely
// no such side-effects.
//
// This function only considers purity of the iteration. Other expressions
// within the node subtree may still have side-effects.
//
// Throws IllegalStateException if node is of a kind that does not trigger
// iteration. An explicit goal of this function is to record all the kinds of
// nodes that do.
//
// Java: NodeUtil.iteratesImpureIterable(Node)
func NodeUtil_IteratesImpureIterable(node rhino.Node) bool {
	parent := node.GetParent()

	var iterable rhino.Node
	switch node.GetToken() {
	case rhino.Token_ITER_SPREAD:
		iterable = node.GetOnlyChild()
	case rhino.Token_YIELD:
		if !node.IsYieldAll() {
			return false // Regular `yield` does not iterate, only `yield*`.
		}
		iterable = node.GetOnlyChild()
	case rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		iterable = node.GetSecondChild()
	case rhino.Token_ITER_REST:
		switch parent.GetToken() {
		// Rest arguments are flat at the call-site.
		case rhino.Token_PARAM_LIST:
			return false
		// TODO(b/127862986): We assume the r-value to be an impure iterable.
		case rhino.Token_ARRAY_PATTERN:
			return true
		default:
			panic(jutil.NewIllegalStateExceptionMsg(
				"Unexpected parent of ITRE_REST: " + parent.ToStringTree()))
		}
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			"Expected a kind of node that may trigger iteration: " + node.ToStringTree()))
	}

	return !NodeUtil_IsPureIterable(iterable)
}

// NodeUtil_DeleteFunctionCall permanently deletes the given call from the
// AST while maintaining a valid node structure, as well as report the
// related AST changes to the given compiler. In some cases, this is done by
// deleting the parent from the AST and is come cases expression is replaced
// by undefined.
//
// Java: NodeUtil.deleteFunctionCall(Node, AbstractCompiler)
func NodeUtil_DeleteFunctionCall(n rhino.Node, compiler AbstractCompiler) {
	jutil.CheckState(n.IsCall())

	parent := n.GetParent()
	if parent.IsExprResult() {
		grandParent := parent.GetParent()
		parent.Detach()
		parent = grandParent
	} else {
		// Seems like part of more complex expression, fallback to replacing with no-op.
		n.ReplaceWith(NodeUtil_NewUndefinedNode(n))
	}

	NodeUtil_MarkFunctionsDeleted(n, compiler)
	compiler.ReportChangeToEnclosingScope(parent)
}

// NodeUtil_DeleteChildren permanently deletes all the children of the given
// node, including reporting changes.
//
// Java: NodeUtil.deleteChildren(Node, AbstractCompiler)
func NodeUtil_DeleteChildren(n rhino.Node, compiler AbstractCompiler) {
	for n.HasChildren() {
		NodeUtil_DeleteNode(n.GetFirstChild(), compiler)
	}
}

// NodeUtil_RemoveChild safely removes children while maintaining a valid
// node structure. In some cases, this is done by removing the parent from
// the AST as well.
//
// Java: NodeUtil.removeChild(Node, Node)
func NodeUtil_RemoveChild(parent rhino.Node, node rhino.Node) {
	if NodeUtil_IsTryFinallyNode(parent, node) {
		if NodeUtil_HasCatchHandler(NodeUtil_GetCatchBlock(parent)) {
			// A finally can only be removed if there is a catch.
			node.Detach()
		} else {
			// Otherwise, only its children can be removed.
			node.DetachChildren()
		}
	} else if node.IsCatch() {
		// The CATCH can can only be removed if there is a finally clause.
		tryNode := node.GetGrandparent()
		jutil.CheckState(NodeUtil_HasFinally(tryNode))
		node.Detach()
	} else if NodeUtil_IsTryCatchNodeContainer(node) {
		// The container node itself can't be removed, but the contained CATCH
		// can if there is a 'finally' clause
		tryNode := node.GetParent()
		jutil.CheckState(NodeUtil_HasFinally(tryNode))
		node.DetachChildren()
	} else if node.IsBlock() {
		// Simply empty the block.  This maintains source location and
		// "synthetic"-ness.
		node.DetachChildren()
	} else if NodeUtil_IsStatementBlock(parent) || NodeUtil_IsSwitchCase(node) || node.IsMemberFunctionDef() {
		// A statement in a block or a member function can simply be removed
		node.Detach()
	} else if NodeUtil_IsNameDeclaration(parent) || parent.IsExprResult() {
		if parent.HasMoreThanOneChild() {
			node.Detach()
		} else {
			// Remove the node from the parent, so it can be reused.
			node.Detach()
			// This would leave an empty VAR, remove the VAR itself.
			NodeUtil_RemoveChild(parent.GetParent(), parent)
		}
	} else if parent.IsLabel() && node == parent.GetLastChild() {
		// Remove the node from the parent, so it can be reused.
		node.Detach()
		// A LABEL without children can not be referred to, remove it.
		NodeUtil_RemoveChild(parent.GetParent(), parent)
	} else if parent.IsVanillaFor() {
		// Only Token.FOR can have an Token.EMPTY other control structure
		// need something for the condition. Others need to be replaced
		// or the structure removed.
		node.ReplaceWith(rhino.IR_Empty())
	} else if parent.IsObjectPattern() {
		// Remove the name from the object pattern
		node.Detach()
	} else if parent.IsArrayPattern() {
		if node == parent.GetLastChild() {
			node.Detach()
		} else {
			node.ReplaceWith(rhino.IR_Empty())
		}
	} else if parent.IsDestructuringLhs() {
		// Destructuring is empty so we should remove the node
		node.Detach()
		if parent.GetParent().HasChildren() {
			// removing the destructuring could leave an empty variable declaration node, so we would
			// want to remove it from the AST
			NodeUtil_RemoveChild(parent.GetParent(), parent)
		}
	} else if parent.IsRest() {
		// Rest params can only ever have one child node
		parent.Detach()
	} else if parent.IsParamList() {
		node.Detach()
	} else if parent.IsImport() {
		// An import node must always have three child nodes. Only the first can be safely removed.
		if node == parent.GetFirstChild() {
			node.ReplaceWith(rhino.IR_Empty())
		} else {
			panic(jutil.NewIllegalStateExceptionMsg(
				"Invalid attempt to remove: " + node.String() + " from " + parent.String()))
		}
	} else {
		panic(jutil.NewIllegalStateExceptionMsg(
			"Invalid attempt to remove node: " + node.String() + " of " + parent.String()))
	}
}

// NodeUtil_ReplaceDeclarationChild replaces the child of a var/let/const
// declaration (usually a name) with a new statement. Preserves the order of
// side effects for all the other declaration children.
//
// declChild is the name node to be replaced; newStatement is the statement to
// replace with.
//
// Java: NodeUtil.replaceDeclarationChild(Node, Node)
func NodeUtil_ReplaceDeclarationChild(declChild rhino.Node, newStatement rhino.Node) {
	jutil.CheckArgument(NodeUtil_IsNameDeclaration(declChild.GetParent()))
	jutil.CheckArgument(newStatement.GetParent().IsNil())

	decl := declChild.GetParent()
	if decl.HasOneChild() {
		decl.ReplaceWith(newStatement)
	} else if declChild.GetNext().IsNil() {
		declChild.Detach()
		newStatement.InsertAfter(decl)
	} else if declChild.GetPrevious().IsNil() {
		declChild.Detach()
		newStatement.InsertBefore(decl)
	} else {
		jutil.CheckState(decl.HasMoreThanOneChild())
		newDecl := rhino.NewNode(decl.GetToken()).Srcref(decl)
		var next rhino.Node
		for after := declChild.GetNext(); !after.IsNil(); after = next {
			next = after.GetNext()
			newDecl.AddChildToBack(after.Detach())
		}
		declChild.Detach()
		newStatement.InsertAfter(decl)
		newDecl.InsertAfter(newStatement)
	}
}

// NodeUtil_MaybeAddFinally adds a finally block if one does not exist.
//
// Java: NodeUtil.maybeAddFinally(Node)
func NodeUtil_MaybeAddFinally(tryNode rhino.Node) {
	jutil.CheckState(tryNode.IsTry())
	if !NodeUtil_HasFinally(tryNode) {
		tryNode.AddChildToBack(rhino.IR_Block().Srcref(tryNode))
	}
}

// NodeUtil_TryMergeBlock merges a block with its parent block.
//
// ignoreBlockScopedDeclarations: merge the block regardless of any inner
// block-scoped declarations that may cause name collisions. use if e.g. the
// AST is normalized. Returns whether the block was removed.
//
// Java: NodeUtil.tryMergeBlock(Node, boolean)
func NodeUtil_TryMergeBlock(block rhino.Node, ignoreBlockScopedDeclarations bool) bool {
	jutil.CheckState(block.IsBlock())
	parent := block.GetParent()
	canMerge := ignoreBlockScopedDeclarations || NodeUtil_CanMergeBlock(block)
	// Try to remove the block if its parent is a block/script or if its
	// parent is label and it has exactly one child.
	if NodeUtil_IsStatementBlock(parent) && canMerge {
		previous := block
		for block.HasChildren() {
			child := block.RemoveFirstChild()
			child.InsertAfter(previous)
			previous = child
		}
		block.Detach()
		return true
	} else {
		return false
	}
}

// NodeUtil_CanMergeBlock: a check inside a block to see if there are const,
// let, class, or function declarations to be safe and not hoist them into
// the upper block.
//
// Returns whether the block can be removed.
//
// Java: NodeUtil.canMergeBlock(Node)
func NodeUtil_CanMergeBlock(block rhino.Node) bool {
	for c := block.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		switch c.GetToken() {
		case rhino.Token_LABEL:
			if NodeUtil_CanMergeBlock(c) {
				continue
			} else {
				return false
			}
		case rhino.Token_CONST, rhino.Token_LET, rhino.Token_CLASS, rhino.Token_FUNCTION:
			return false
		default:
			continue
		}
	}
	return true
}

// NodeUtil_NewQName creates a node representing a qualified name.
//
// name is a qualified name (e.g. "foo" or "foo.bar.baz"); returns a NAME or
// GETPROP node.
//
// Java: NodeUtil.newQName(AbstractCompiler, String)
func NodeUtil_NewQName(compiler AbstractCompiler, name string) rhino.Node {
	dot := strings.IndexByte(name, '.')
	nodeName := name
	if dot != -1 {
		nodeName = name[:dot]
	}

	var qname rhino.Node
	if "this" == nodeName {
		qname = rhino.IR_ThisNode()
	} else if "super" == nodeName {
		qname = rhino.IR_SuperNode()
	} else {
		qname = NodeUtil_NewName(compiler, nodeName)
	}
	qname.SetLength(jutil.StringLength(nodeName))
	if dot == -1 {
		return qname
	}

	conv := compiler.GetCodingConvention()
	rest := name[dot+1:]
	for {
		dot = strings.IndexByte(rest, '.')
		part := rest
		if dot != -1 {
			part = rest[:dot]
			rest = rest[dot+1:]
		}

		qname = rhino.IR_Getprop(qname, part)
		if conv.IsConstantKey(part) {
			qname.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
		}
		qname.SetLength(jutil.StringLength(part))
		if dot == -1 {
			break
		}
	}

	return qname
}

// NodeUtil_NewQNameWithOriginalName creates a node representing a qualified
// name, copying over the source location information from the basis node and
// assigning the given original name to the node.
//
// name is a qualified name (e.g. "foo" or "foo.bar.baz"); basisNode is the
// node that represents the name as currently found in the AST; originalName
// is the original name of the item being represented by the NAME node. Used
// for debugging information. Returns a NAME or GETPROP node.
//
// Java: NodeUtil.newQName(AbstractCompiler, String, Node, String)
func NodeUtil_NewQNameWithOriginalName(
	compiler AbstractCompiler, name string, basisNode rhino.Node, originalName string) rhino.Node {
	node := NodeUtil_NewQName(compiler, name)
	node.SrcrefTreeIfMissing(basisNode)
	if nodeOriginalName := node.GetOriginalName(); nodeOriginalName == nil || originalName != *nodeOriginalName {
		// If basisNode already had the correct original name, then it will already be set correctly.
		// Setting it again will force the QName node to have a different property list from all of
		// its children, causing greater memory consumption.
		node.SetOriginalName(jutil.StringPtr(originalName))
	}
	return node
}

// NodeUtil_GetDeclarationFromName attaches nameNode to a new qualified name
// declaration and returns the new qualified declaration.
//
// Returns a new qualified name declaration.
//
// Java: NodeUtil.getDeclarationFromName(Node, Node, Token, JSDocInfo)
func NodeUtil_GetDeclarationFromName(nameNode rhino.Node, value rhino.Node, type_ rhino.Token, info *rhino.JSDocInfo) rhino.Node {
	var result rhino.Node
	if nameNode.IsName() {
		if value.IsNil() {
			result = rhino.IR_Declaration(nameNode, type_)
		} else {
			result = rhino.IR_DeclarationWithValue(nameNode, value, type_)
		}
		result.SetJSDocInfo(info)
	} else if !value.IsNil() {
		result = rhino.IR_ExprResult(rhino.IR_Assign(nameNode, value))
		result.GetFirstChild().SetJSDocInfo(info)
	} else {
		result = rhino.IR_ExprResult(nameNode)
		result.GetFirstChild().SetJSDocInfo(info)
	}
	return result
}

// NodeUtil_NewPropertyAccess creates a property access on the context tree.
//
// Java: NodeUtil.newPropertyAccess(AbstractCompiler, Node, String)
func NodeUtil_NewPropertyAccess(compiler AbstractCompiler, context rhino.Node, name string) rhino.Node {
	propNode := rhino.IR_Getprop(context, name)
	if compiler.GetCodingConvention().IsConstantKey(name) {
		propNode.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	}
	return propNode
}

// NodeUtil_NewQNameDeclaration creates a node representing a qualified name.
//
// name is a qualified name (e.g. "foo" or "foo.bar.baz"); returns a VAR
// node, or an EXPR_RESULT node containing an ASSIGN or NAME node.
//
// Java: NodeUtil.newQNameDeclaration(AbstractCompiler, String, Node, JSDocInfo)
func NodeUtil_NewQNameDeclaration(
	compiler AbstractCompiler, name string, value rhino.Node, info *rhino.JSDocInfo) rhino.Node {
	return NodeUtil_NewQNameDeclarationWithType(compiler, name, value, info, rhino.Token_VAR)
}

// NodeUtil_NewQNameDeclarationWithType creates a node representing a
// qualified name.
//
// name is a qualified name (e.g. "foo" or "foo.bar.baz"); type_ must be VAR,
// CONST, or LET. Ignored if name is dotted. Returns a VAR/CONST/LET node, or
// an EXPR_RESULT node containing an ASSIGN or NAME node.
//
// Java: NodeUtil.newQNameDeclaration(AbstractCompiler, String, Node, JSDocInfo, Token)
func NodeUtil_NewQNameDeclarationWithType(
	compiler AbstractCompiler, name string, value rhino.Node, info *rhino.JSDocInfo, type_ rhino.Token) rhino.Node {
	jutil.CheckStateMsg1(type_ == rhino.Token_VAR || type_ == rhino.Token_LET || type_ == rhino.Token_CONST, "%s", type_)
	nameNode := NodeUtil_NewQName(compiler, name)
	return NodeUtil_GetDeclarationFromName(nameNode, value, type_, info)
}

// NodeUtil_GetRootOfQualifiedName gets the root node of a qualified name.
// Must be either NAME, THIS or SUPER.
//
// Java: NodeUtil.getRootOfQualifiedName(Node)
func NodeUtil_GetRootOfQualifiedName(qName rhino.Node) rhino.Node {
	for current := qName; ; current = current.GetFirstChild() {
		if current.IsName() || current.IsThis() || current.IsSuper() {
			return current
		}
		jutil.CheckStateMsg1(current.IsGetProp(), "Not a getprop node:  (%s)", current)
	}
}

// NodeUtil_NewUndefinedNode creates a node for an empty result expression:
// "void 0".
//
// Java: NodeUtil.newUndefinedNode(Node)
func NodeUtil_NewUndefinedNode(srcReferenceNode rhino.Node) rhino.Node {
	node := rhino.IR_VoidNode(rhino.IR_Number(0))
	if !srcReferenceNode.IsNil() {
		node.SrcrefTree(srcReferenceNode)
	}
	return node
}

// NodeUtil_NewVarNode creates a VAR node containing the given name and
// initial value expression.
//
// Java: NodeUtil.newVarNode(String, Node)
func NodeUtil_NewVarNode(name string, value rhino.Node) rhino.Node {
	lhs := rhino.IR_Name(name)
	if !value.IsNil() {
		lhs.Srcref(value)
	}
	return NodeUtil_NewVarNodeWithLhs(lhs, value)
}

// NodeUtil_NewVarNodeWithLhs creates a VAR node containing the given lhs
// (name or destructuring pattern) and initial value expression.
//
// Java: NodeUtil.newVarNode(Node, Node)
func NodeUtil_NewVarNodeWithLhs(lhs rhino.Node, value rhino.Node) rhino.Node {
	if lhs.IsDestructuringPattern() {
		// TODO(b/314005766): Ensure that destructuring declarations are generated as normalized (i.e
		// individual names get their own stub declarations)
		rhino.CheckNotNull(value)
		return rhino.IR_Var(rhino.NewNodeWith2Children(rhino.Token_DESTRUCTURING_LHS, lhs, value).Srcref(lhs)).Srcref(lhs)
	} else {
		jutil.CheckState(lhs.IsName() && !lhs.HasChildren())
		if !value.IsNil() {
			lhs.AddChildToBack(value)
		}
		return rhino.IR_Var(lhs).Srcref(lhs)
	}
}

// Java: NodeUtil.emptyFunction()
func NodeUtil_EmptyFunction() rhino.Node {
	return rhino.IR_Function(rhino.IR_Name(""), rhino.IR_ParamList(), rhino.IR_Block())
}

// NodeUtil_MatchNameNode is a predicate for matching name nodes with the
// specified node.
//
// Java: NodeUtil.MatchNameNode
type NodeUtil_MatchNameNode struct {
	name string
}

// NodeUtil_EvaluatesToLocalValue: whether the result of the expression node
// is known to be a primitive value or an object that has not yet escaped.
//
// This guarantee is different than that provided by isLiteralValue (where
// literal values are immune to side-effects if unescaped) or
// isImmutableValue (which can be safely aliased).
//
// The concept of "local values" allow for the containment of side-effect
// operations. For example, setting a property on a local value does not
// produce a global side-effect.
//
// Note that the concept of "local value" is not deep, it does not say
// anything about the properties of the "local value" (all class instances
// have "constructor" properties that are not local values for instance).
//
// Note that this method only provides the starting state of the expression
// result, it does not guarantee that the value is forever a local value. If
// the containing method has any non-local side-effect, "local values" may
// escape.
//
// Java: NodeUtil.evaluatesToLocalValue(Node)
func NodeUtil_EvaluatesToLocalValue(value rhino.Node) bool {
	switch value.GetToken() {
	case rhino.Token_ASSIGN:
		// A result that is aliased by a non-local name, is the effectively the
		// same as returning a non-local name, but this doesn't matter if the
		// value is immutable.
		return NodeUtil_IsImmutableValue(value.GetLastChild())
	case rhino.Token_COMMA:
		return NodeUtil_EvaluatesToLocalValue(value.GetLastChild())
	case rhino.Token_AND, rhino.Token_OR, rhino.Token_COALESCE:
		return NodeUtil_EvaluatesToLocalValue(value.GetFirstChild()) &&
			NodeUtil_EvaluatesToLocalValue(value.GetLastChild())
	case rhino.Token_HOOK:
		return NodeUtil_EvaluatesToLocalValue(value.GetSecondChild()) &&
			NodeUtil_EvaluatesToLocalValue(value.GetLastChild())
	case rhino.Token_DYNAMIC_IMPORT:
		// Dynamic import always returns a newly created Promise.
		return true
	case rhino.Token_THIS, rhino.Token_SUPER:
		return false
	case rhino.Token_NAME:
		return NodeUtil_IsImmutableValue(value)
	case rhino.Token_GETELEM, rhino.Token_GETPROP, rhino.Token_OPTCHAIN_GETELEM, rhino.Token_OPTCHAIN_GETPROP:
		// There is no information about the locality of object properties.
		return false
	case rhino.Token_CALL, rhino.Token_OPTCHAIN_CALL:
		return NodeUtil_IsToStringMethodCall(value)
	case rhino.Token_TAGGED_TEMPLATELIT:
		// No information about local values for tagged template literals
		return false
	case rhino.Token_NEW:
		return NodeUtil_NewHasLocalResult(value)
	case rhino.Token_DELPROP, rhino.Token_INC, rhino.Token_DEC, rhino.Token_CLASS, rhino.Token_FUNCTION,
		rhino.Token_REGEXP, rhino.Token_EMPTY, rhino.Token_ARRAYLIT, rhino.Token_OBJECTLIT, rhino.Token_TEMPLATELIT:
		return true
	case rhino.Token_CAST:
		return NodeUtil_EvaluatesToLocalValue(value.GetFirstChild())
	case rhino.Token_ITER_SPREAD, rhino.Token_OBJECT_SPREAD, rhino.Token_NEW_TARGET:
		// TODO(johnlenz): remove this case.
		// Returns an alias of a constructor (current or subclass).
		return false
	case rhino.Token_YIELD, rhino.Token_AWAIT:
		// TODO(johnlenz): we can do better for await if we use type information.  That is,
		// if we know the promise being awaited on is a immutable value type (string, etc)
		// we could return true here.
		return false
	default:
		// A logical assignment could evaluate to either the assignment target or the
		// right-hand side, and we generally have no information about the locality of the lhs.
		if NodeUtil_IsLogicalAssignmentOp(value) {
			return false
		}
		// Other non-logical assignment ops force a local value:
		//  '' + g (a local string)
		//  x -= g (x is now an local number)
		if NodeUtil_IsAssignmentOp(value) || NodeUtil_IsSimpleOperator(value) || NodeUtil_IsImmutableValue(value) {
			return true
		}

		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected expression node: " + value.String() + "\n parent:" + jutil.ToJavaString(value.GetParent())))
	}
}

// NodeUtil_BooleanNode returns an appropriate AST node for the boolean
// value.
//
// Java: NodeUtil.booleanNode(boolean)
func NodeUtil_BooleanNode(value bool) rhino.Node {
	if value {
		return rhino.IR_TrueNode()
	}
	return rhino.IR_FalseNode()
}

// NodeUtil_NumberNode returns an appropriate AST node for the double value.
//
// Java: NodeUtil.numberNode(double, Node)
func NodeUtil_NumberNode(value float64, srcref rhino.Node) rhino.Node {
	var result rhino.Node
	if math.IsNaN(value) {
		result = rhino.IR_Name("NaN")
		result.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	} else {
		if math.IsInf(value, 0) {
			result = rhino.IR_Name("Infinity")
			result.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
		} else {
			result = rhino.IR_Number(math.Abs(value))
		}
		if base.JSCompDoubles_IsNegative(value) {
			result = rhino.IR_Neg(result)
		}
	}
	if !srcref.IsNil() {
		result.SrcrefTree(srcref)
	}
	return result
}

// Java: NodeUtil.isGoogProvideCall(Node)
func NodeUtil_IsGoogProvideCall(n rhino.Node) bool {
	if NodeUtil_IsExprCall(n) {
		target := n.GetFirstFirstChild()
		return NodeUtil_GOOG_PROVIDE.Matches(target)
	}
	return false
}

// Java: NodeUtil.createSynthesizedExternsSymbol(AbstractCompiler, String)
func NodeUtil_CreateSynthesizedExternsSymbol(compiler AbstractCompiler, nameToAdd string) {
	name := rhino.IR_Name(nameToAdd)
	name.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	var_ := rhino.IR_Var(name)
	input := compiler.GetSynthesizedExternsInput()
	root := input.GetAstRoot(compiler)
	name.SetStaticSourceFileFrom(root)
	var_.SetStaticSourceFileFrom(root)
	root.AddChildToBack(var_)
	compiler.ReportChangeToEnclosingScope(var_)
}

// NodeUtil_MarkNewScopesChanged recurses through a tree, marking all
// function nodes as changed.
//
// Java: NodeUtil.markNewScopesChanged(Node, AbstractCompiler)
func NodeUtil_MarkNewScopesChanged(node rhino.Node, compiler AbstractCompiler) {
	if node.IsFunction() {
		compiler.ReportChangeToChangeScope(node)
	}
	for child := node.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		NodeUtil_MarkNewScopesChanged(child, compiler)
	}
}

// NodeUtil_MarkFunctionsDeleted recurses through a tree, marking all
// function nodes deleted.
//
// Java: NodeUtil.markFunctionsDeleted(Node, AbstractCompiler)
func NodeUtil_MarkFunctionsDeleted(node rhino.Node, compiler AbstractCompiler) {
	if node.IsFunction() {
		compiler.ReportFunctionDeleted(node)
	}

	for child := node.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		NodeUtil_MarkFunctionsDeleted(child, compiler)
	}

	shadowed := node.GetClosureUnawareShadow()
	if !shadowed.IsNil() {
		NodeUtil_MarkFunctionsDeleted(shadowed, compiler)
	}
}

// NodeUtil_GetParentChangeScopeNodes returns the list of scope nodes which
// are parents of the provided list of scope nodes.
//
// Java: NodeUtil.getParentChangeScopeNodes(List<Node>)
func NodeUtil_GetParentChangeScopeNodes(scopeNodes []rhino.Node) []rhino.Node {
	if len(scopeNodes) == 0 {
		return []rhino.Node{}
	}
	seen := make(map[rhino.Node]struct{}, len(scopeNodes)*2)
	out := make([]rhino.Node, 0, len(scopeNodes)*2)
	for _, scopeNode := range scopeNodes {
		if _, ok := seen[scopeNode]; !ok {
			seen[scopeNode] = struct{}{}
			out = append(out, scopeNode)
		}
	}
	for _, scopeNode := range scopeNodes {
		parent := ChangeTracker_GetEnclosingChangeScopeRoot(scopeNode)
		if _, ok := seen[parent]; !ok {
			seen[parent] = struct{}{}
			out = append(out, parent)
		}
	}
	return out
}

// NodeUtil_RemoveNestedChangeScopeNodes removes any scope nodes from the
// provided list that are nested within some other scope node also in the
// list. Returns the modified list.
//
// Java: NodeUtil.removeNestedChangeScopeNodes(List<Node>)
func NodeUtil_RemoveNestedChangeScopeNodes(scopeNodes []rhino.Node) []rhino.Node {
	if len(scopeNodes) == 0 {
		return []rhino.Node{}
	}
	seen := make(map[rhino.Node]struct{}, len(scopeNodes))
	out := make([]rhino.Node, 0, len(scopeNodes))
	for _, scopeNode := range scopeNodes {
		if _, ok := seen[scopeNode]; !ok {
			seen[scopeNode] = struct{}{}
			out = append(out, scopeNode)
		}
	}
	n := 0
	for _, scopeNode := range out {
		nested := false
		for ancestor := scopeNode.GetParent(); !ancestor.IsNil(); ancestor = ancestor.GetParent() {
			if ChangeTracker_IsChangeScopeRoot(ancestor) {
				if _, ok := seen[ancestor]; ok {
					nested = true
					break
				}
			}
		}
		if !nested {
			out[n] = scopeNode
			n++
		}
	}
	clear(out[n:])
	return out[:n]
}
