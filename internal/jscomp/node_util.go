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
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// NodeUtil contains generally useful AST utilities.
type NodeUtil struct{}

// Value of JavaScript's Number.MAX_SAFE_INTEGER
//
// Java: NodeUtil.MAX_POSITIVE_INTEGER_NUMBER
const NodeUtil_MAX_POSITIVE_INTEGER_NUMBER int64 = (1 << 53) - 1

// Java: NodeUtil.JSC_PROPERTY_NAME_FN
const NodeUtil_JSC_PROPERTY_NAME_FN = "JSCompiler_renameProperty"

// Java: NodeUtil.LARGEST_BASIC_LATIN
const NodeUtil_LARGEST_BASIC_LATIN uint16 = 0x7f

// Java: NodeUtil.GOOG_MODULE_DECLARE_LEGACY_NAMESPACE
var NodeUtil_GOOG_MODULE_DECLARE_LEGACY_NAMESPACE = rhino.QualifiedName_Of("goog.module.declareLegacyNamespace")

// Java: NodeUtil.GOOG_SET_TEST_ONLY
var NodeUtil_GOOG_SET_TEST_ONLY = rhino.QualifiedName_Of("goog.setTestOnly")

// Java: NodeUtil.GOOG_PROVIDE
var NodeUtil_GOOG_PROVIDE = rhino.QualifiedName_Of("goog.provide")

// Java: NodeUtil.GOOG_MODULE
var NodeUtil_GOOG_MODULE = rhino.QualifiedName_Of("goog.module")

// Java: NodeUtil.GOOG_MODULE_GET
var NodeUtil_GOOG_MODULE_GET = rhino.QualifiedName_Of("goog.module.get")

// Java: NodeUtil.GOOG_REQUIRE
var NodeUtil_GOOG_REQUIRE = rhino.QualifiedName_Of("goog.require")

// Java: NodeUtil.GOOG_REQUIRE_TYPE
var NodeUtil_GOOG_REQUIRE_TYPE = rhino.QualifiedName_Of("goog.requireType")

// Java: NodeUtil.GOOG_FORWARD_DECLARE
var NodeUtil_GOOG_FORWARD_DECLARE = rhino.QualifiedName_Of("goog.forwardDeclare")

// Java: NodeUtil.GOOG_REQUIRE_DYNAMIC
var NodeUtil_GOOG_REQUIRE_DYNAMIC = rhino.QualifiedName_Of("goog.requireDynamic")

// Java: NodeUtil.GOOG_WEAK_USAGE
var NodeUtil_GOOG_WEAK_USAGE = rhino.QualifiedName_Of("goog.weakUsage")

// Java: NodeUtil.GOOG_WEAK_USAGE_MANGLED
var NodeUtil_GOOG_WEAK_USAGE_MANGLED = rhino.QualifiedName_Of("goog$weakUsage")

// Utility class; do not instantiate.
//
// Java: NodeUtil()
func NewNodeUtil() *NodeUtil {
	return &NodeUtil{}
}

// NodeUtil_IsSymmetricOperation returns true if the operator on this node is
// symmetric.
//
// Java: NodeUtil.isSymmetricOperation(Node)
func NodeUtil_IsSymmetricOperation(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_EQ, // equal
		rhino.Token_NE,   // not equal
		rhino.Token_SHEQ, // exactly equal
		rhino.Token_SHNE, // exactly not equal
		rhino.Token_MUL:  // multiply, unlike add it only works on numbers
		// or results NaN if any of the operators is not a number
		return true
	default:
		return false
	}
}

// NodeUtil_IsRelationalOperation returns true if the operator on this node is
// relational. the returned set does not include the equalities.
//
// Java: NodeUtil.isRelationalOperation(Node)
func NodeUtil_IsRelationalOperation(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_GT, rhino.Token_GE, rhino.Token_LT, rhino.Token_LE:
		return true
	default:
		return false
	}
}

// NodeUtil_GetInverseOperator returns the inverse of an operator if it is
// invertible. ex. '>' ==> '<'
//
// Java: NodeUtil.getInverseOperator(Token)
func NodeUtil_GetInverseOperator(type_ rhino.Token) rhino.Token {
	switch type_ {
	case rhino.Token_GT:
		return rhino.Token_LT
	case rhino.Token_LT:
		return rhino.Token_GT
	case rhino.Token_GE:
		return rhino.Token_LE
	case rhino.Token_LE:
		return rhino.Token_GE
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg("Unexpected token: " + type_.String()))
	}
}

// NodeUtil_IsLiteralValue returns true if this is a literal value. We define
// a literal value as any node that evaluates to the same thing regardless of
// when or where it is evaluated. So /xyz/ and [3, 5] are literals, but the
// name a is not.
//
// Function literals do not meet this definition, because they lexically
// capture variables. For example, if you have
// `function() { return a; }`. If it is evaluated in a different scope, then
// it captures a different variable. Even if the function did not read any
// captured variables directly, it would still fail this definition, because
// it affects the lifecycle of variables in the enclosing scope.
//
// However, a function literal with respect to a particular scope is a
// literal.
//
// includeFunctions: if true, all function expressions will be treated as
// literals.
//
// Java: NodeUtil.isLiteralValue(Node, boolean)
func NodeUtil_IsLiteralValue(n rhino.Node, includeFunctions bool) bool {
	switch n.GetToken() {
	case rhino.Token_CAST:
		return NodeUtil_IsLiteralValue(n.GetFirstChild(), includeFunctions)
	case rhino.Token_ARRAYLIT:
		for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			if (!child.IsEmpty()) && !NodeUtil_IsLiteralValue(child, includeFunctions) {
				return false
			}
		}
		return true
	case rhino.Token_REGEXP:
		// Return true only if all descendants are const.
		for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			if !NodeUtil_IsLiteralValue(child, includeFunctions) {
				return false
			}
		}
		return true
	case rhino.Token_OBJECTLIT:
		for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			switch child.GetToken() {
			case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
				// { methodName() {...} }
				// { get propertyName() {...} }
				// { set propertyName(value) {...} }
				if !includeFunctions {
					return false
				}
				// keep going

			case rhino.Token_COMPUTED_PROP:
				// { [key_expression]: value, ... }
				// { [key_expression](args) {...}, ... }
				// { get [key_expression]() {...}, ... }
				// { set [key_expression](args) {...}, ... }
				if !NodeUtil_IsLiteralValue(child.GetFirstChild(), includeFunctions) ||
					!NodeUtil_IsLiteralValue(child.GetLastChild(), includeFunctions) {
					return false
				}
				// keep going

			case rhino.Token_OBJECT_SPREAD:
				if !NodeUtil_IsLiteralValue(child.GetOnlyChild(), includeFunctions) {
					return false
				}
				// keep going

			case rhino.Token_STRING_KEY:
				// { key: value, ... }
				// { "quoted_key": value, ... }
				if !NodeUtil_IsLiteralValue(child.GetOnlyChild(), includeFunctions) {
					return false
				}
				// keep going

			default:
				panic(jutil.NewIllegalArgumentExceptionMsg(
					"Unexpected child of OBJECTLIT: " + child.ToStringTree()))
			}
		}
		return true
	case rhino.Token_FUNCTION:
		return includeFunctions && !NodeUtil_IsFunctionDeclaration(n)
	case rhino.Token_TEMPLATELIT:
		for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			if child.IsTemplateLitSub() {
				if !NodeUtil_IsLiteralValue(child.GetFirstChild(), includeFunctions) {
					return false
				}
			}
		}
		return true
	default:
		return NodeUtil_IsImmutableValue(n)
	}
}

// Java: NodeUtil.isBinaryOperator(Node)
func NodeUtil_IsBinaryOperator(n rhino.Node) bool {
	return NodeUtil_IsBinaryOperatorType(n.GetToken())
}

// NodeUtil_IsBinaryOperatorType: an operator with two operands that does not
// assign a value to either. Once you cut through the layers of rules, these
// all parse similarly, taking LeftHandSideExpression operands on either side.
// Comma is not included, because it takes AssignmentExpression operands,
// making its syntax different.
//
// Java: NodeUtil.isBinaryOperatorType(Token)
func NodeUtil_IsBinaryOperatorType(type_ rhino.Token) bool {
	switch type_ {
	case rhino.Token_OR,
		rhino.Token_AND,
		rhino.Token_COALESCE,
		rhino.Token_BITOR,
		rhino.Token_BITXOR,
		rhino.Token_BITAND,
		rhino.Token_EQ,
		rhino.Token_NE,
		rhino.Token_SHEQ,
		rhino.Token_SHNE,
		rhino.Token_LT,
		rhino.Token_GT,
		rhino.Token_LE,
		rhino.Token_GE,
		rhino.Token_INSTANCEOF,
		rhino.Token_IN,
		rhino.Token_LSH,
		rhino.Token_RSH,
		rhino.Token_URSH,
		rhino.Token_ADD,
		rhino.Token_SUB,
		rhino.Token_MUL,
		rhino.Token_DIV,
		rhino.Token_MOD,
		rhino.Token_EXPONENT:
		return true
	default:
		return false
	}
}

// Java: NodeUtil.isUnaryOperator(Node)
func NodeUtil_IsUnaryOperator(n rhino.Node) bool {
	return NodeUtil_IsUnaryOperatorType(n.GetToken())
}

// NodeUtil_IsUnaryOperatorType: an operator taking only one operand. These
// all parse very similarly, taking LeftHandSideExpression operands.
//
// Java: NodeUtil.isUnaryOperatorType(Token)
func NodeUtil_IsUnaryOperatorType(type_ rhino.Token) bool {
	switch type_ {
	case rhino.Token_DELPROP, rhino.Token_VOID, rhino.Token_TYPEOF, rhino.Token_POS, rhino.Token_NEG,
		rhino.Token_BITNOT, rhino.Token_NOT:
		return true
	default:
		return false
	}
}

// Java: NodeUtil.isUpdateOperator(Node)
func NodeUtil_IsUpdateOperator(n rhino.Node) bool {
	return NodeUtil_IsUpdateOperatorType(n.GetToken())
}

// Java: NodeUtil.isUpdateOperatorType(Token)
func NodeUtil_IsUpdateOperatorType(type_ rhino.Token) bool {
	switch type_ {
	case rhino.Token_INC, rhino.Token_DEC:
		return true
	default:
		return false
	}
}

// Java: NodeUtil.isSimpleOperator(Node)
func NodeUtil_IsSimpleOperator(n rhino.Node) bool {
	return NodeUtil_IsSimpleOperatorType(n.GetToken())
}

// NodeUtil_IsSimpleOperatorType: a "simple" operator is one whose children
// are expressions, has no direct side-effects (unlike '+='), and has no
// conditional aspects (unlike '||').
//
// Java: NodeUtil.isSimpleOperatorType(Token)
func NodeUtil_IsSimpleOperatorType(type_ rhino.Token) bool {
	switch type_ {
	case rhino.Token_ADD,
		rhino.Token_BITAND,
		rhino.Token_BITNOT,
		rhino.Token_BITOR,
		rhino.Token_BITXOR,
		rhino.Token_COMMA,
		rhino.Token_DIV,
		rhino.Token_EQ,
		rhino.Token_EXPONENT,
		rhino.Token_GE,
		rhino.Token_GT,
		rhino.Token_IN,
		rhino.Token_INSTANCEOF,
		rhino.Token_LE,
		rhino.Token_LSH,
		rhino.Token_LT,
		rhino.Token_MOD,
		rhino.Token_MUL,
		rhino.Token_NE,
		rhino.Token_NOT,
		rhino.Token_RSH,
		rhino.Token_SHEQ,
		rhino.Token_SHNE,
		rhino.Token_SUB,
		rhino.Token_TYPEOF,
		rhino.Token_VOID,
		rhino.Token_POS,
		rhino.Token_NEG,
		rhino.Token_URSH:
		return true
	default:
		return false
	}
}

// NodeUtil_IsNamespaceDecl returns true iff this node defines a namespace,
// e.g.,
//
//	/** @const */ var goog = {};
//	/** @const */ var goog = goog || {};
//	/** @const */ goog.math = goog.math || {};
//
// Java: NodeUtil.isNamespaceDecl(Node)
func NodeUtil_IsNamespaceDecl(n rhino.Node) bool {
	jsdoc := NodeUtil_GetBestJSDocInfo(n)
	if jsdoc != nil && jsdoc.HasTypeNodes() && !jsdoc.HasTypedefType() {
		return false
	}
	// In externs, we allow namespace definitions without @const.
	// This is a worse design than always requiring @const, but it helps with
	// namespaces that are defined in many places, such as gapi.
	// Also, omitting @const in externs is not as confusing as in source code,
	// because assigning an object literal in externs only makes sense when
	// defining a namespace or enum.
	isMarkedConst := n.GetParent().IsConst() || (jsdoc != nil && jsdoc.IsConstant())
	if !n.IsFromExterns() && !isMarkedConst {
		return false
	}
	var qnameNode rhino.Node
	var initializer rhino.Node
	if NodeUtil_IsNameDeclaration(n.GetParent()) {
		qnameNode = n
		initializer = n.GetFirstChild()
	} else if n.IsExprResult() {
		expr := n.GetFirstChild()
		if !expr.IsAssign() || !expr.GetFirstChild().IsGetProp() {
			return false
		}
		qnameNode = expr.GetFirstChild()
		initializer = expr.GetLastChild()
	} else if n.IsGetProp() {
		parent := n.GetParent()
		if !parent.IsAssign() || !parent.GetParent().IsExprResult() {
			return false
		}
		qnameNode = n
		initializer = parent.GetLastChild()
	} else {
		return false
	}
	if initializer.IsNil() || qnameNode.IsNil() {
		return false
	}
	if initializer.IsObjectLit() {
		return true
	}
	return initializer.IsOr() &&
		qnameNode.MatchesQualifiedNameNode(initializer.GetFirstChild()) &&
		initializer.GetLastChild().IsObjectLit()
}

// NodeUtil_KNOWN_CONSTANTS: we will assume that side effects never change
// the values of these global variables.
//
// Java: NodeUtil.KNOWN_CONSTANTS
var NodeUtil_KNOWN_CONSTANTS = map[string]bool{
	"undefined": true,
	"Infinity":  true,
	"NaN":       true,
}

// Java: NodeUtil.isUndefined(Node)
func NodeUtil_IsUndefined(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_VOID:
		return true
	case rhino.Token_NAME:
		return n.GetString() == "undefined"
	default:
		return false
	}
}

// NodeUtil_IsCommutative returns true if the operator is commutative. e.g.
// (a * b) * c = c * (b * a) Note 1: "+" is not commutative because it is also
// the concatenation for strings. e.g. "a" + (1 + 2) is not "a" + 1 + 2 Note
// 2: only operations on literals and pure functions are commutative.
//
// Java: NodeUtil.isCommutative(Token)
func NodeUtil_IsCommutative(type_ rhino.Token) bool {
	switch type_ {
	case rhino.Token_MUL, rhino.Token_BITOR, rhino.Token_BITXOR, rhino.Token_BITAND:
		return true
	default:
		return false
	}
}

// NodeUtil_IsAssignmentOp returns true if the operator is an assignment type
// operator. Note: The logical assignments (i.e. ASSIGN_OR, ASSIGN_AND,
// ASSIGN_COALESCE) follow short-circuiting behavior, and the RHS may not
// always be evaluated. They are still considered AssignmentOps (may be
// optimized).
//
// Java: NodeUtil.isAssignmentOp(Node)
func NodeUtil_IsAssignmentOp(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_ASSIGN,
		rhino.Token_ASSIGN_BITOR,
		rhino.Token_ASSIGN_BITXOR,
		rhino.Token_ASSIGN_BITAND,
		rhino.Token_ASSIGN_LSH,
		rhino.Token_ASSIGN_RSH,
		rhino.Token_ASSIGN_URSH,
		rhino.Token_ASSIGN_ADD,
		rhino.Token_ASSIGN_SUB,
		rhino.Token_ASSIGN_MUL,
		rhino.Token_ASSIGN_EXPONENT,
		rhino.Token_ASSIGN_DIV,
		rhino.Token_ASSIGN_MOD,
		rhino.Token_ASSIGN_OR,
		rhino.Token_ASSIGN_AND,
		rhino.Token_ASSIGN_COALESCE:
		return true
	default:
		return false
	}
}

// NodeUtil_IsLogicalAssignmentOp returns true if the operator is a logical
// assignment type operator.
//
// Java: NodeUtil.isLogicalAssignmentOp(Node)
func NodeUtil_IsLogicalAssignmentOp(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_ASSIGN_OR, rhino.Token_ASSIGN_AND, rhino.Token_ASSIGN_COALESCE:
		return true
	default:
		return false
	}
}

// Java: NodeUtil.isCompoundAssignmentOp(Node)
func NodeUtil_IsCompoundAssignmentOp(n rhino.Node) bool {
	return NodeUtil_IsAssignmentOp(n) && !n.IsAssign()
}

// Java: NodeUtil.getOpFromAssignmentOp(Node)
func NodeUtil_GetOpFromAssignmentOp(n rhino.Node) rhino.Token {
	switch n.GetToken() {
	case rhino.Token_ASSIGN_BITOR:
		return rhino.Token_BITOR
	case rhino.Token_ASSIGN_BITXOR:
		return rhino.Token_BITXOR
	case rhino.Token_ASSIGN_BITAND:
		return rhino.Token_BITAND
	case rhino.Token_ASSIGN_LSH:
		return rhino.Token_LSH
	case rhino.Token_ASSIGN_RSH:
		return rhino.Token_RSH
	case rhino.Token_ASSIGN_URSH:
		return rhino.Token_URSH
	case rhino.Token_ASSIGN_ADD:
		return rhino.Token_ADD
	case rhino.Token_ASSIGN_SUB:
		return rhino.Token_SUB
	case rhino.Token_ASSIGN_MUL:
		return rhino.Token_MUL
	case rhino.Token_ASSIGN_EXPONENT:
		return rhino.Token_EXPONENT
	case rhino.Token_ASSIGN_DIV:
		return rhino.Token_DIV
	case rhino.Token_ASSIGN_MOD:
		return rhino.Token_MOD
	case rhino.Token_ASSIGN_OR:
		return rhino.Token_OR
	case rhino.Token_ASSIGN_AND:
		return rhino.Token_AND
	case rhino.Token_ASSIGN_COALESCE:
		return rhino.Token_COALESCE
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg("Not an assignment op:" + n.String()))
	}
}

// NodeUtil_GetEnclosingType gets the closest ancestor to the given node of
// the provided type.
//
// Java: NodeUtil.getEnclosingType(Node, Token)
func NodeUtil_GetEnclosingType(n rhino.Node, type_ rhino.Token) rhino.Node {
	return NodeUtil_GetEnclosingNode(n, func(n1 rhino.Node) bool { return n1.GetToken() == type_ })
}

var nodeUtil_opToStrTable = func() []*string {
	values := rhino.Token_values()
	table := make([]*string, len(values))
	for _, t := range values {
		if s, ok := nodeUtil_opToStrSwitch(t); ok {
			table[t] = jutil.StringPtr(s)
		}
	}
	return table
}()

// nodeUtil_opToStrSwitch is the switch of NodeUtil.opToStr(Token); ok is
// false where Java returns null.
func nodeUtil_opToStrSwitch(operator rhino.Token) (string, bool) {
	switch operator {
	case rhino.Token_COALESCE:
		return "??", true
	case rhino.Token_BITOR:
		return "|", true
	case rhino.Token_OR:
		return "||", true
	case rhino.Token_BITXOR:
		return "^", true
	case rhino.Token_AND:
		return "&&", true
	case rhino.Token_BITAND:
		return "&", true
	case rhino.Token_SHEQ:
		return "===", true
	case rhino.Token_EQ:
		return "==", true
	case rhino.Token_NOT:
		return "!", true
	case rhino.Token_NE:
		return "!=", true
	case rhino.Token_SHNE:
		return "!==", true
	case rhino.Token_LSH:
		return "<<", true
	case rhino.Token_IN:
		return "in", true
	case rhino.Token_LE:
		return "<=", true
	case rhino.Token_LT:
		return "<", true
	case rhino.Token_URSH:
		return ">>>", true
	case rhino.Token_RSH:
		return ">>", true
	case rhino.Token_GE:
		return ">=", true
	case rhino.Token_GT:
		return ">", true
	case rhino.Token_MUL:
		return "*", true
	case rhino.Token_DIV:
		return "/", true
	case rhino.Token_MOD:
		return "%", true
	case rhino.Token_EXPONENT:
		return "**", true
	case rhino.Token_BITNOT:
		return "~", true
	case rhino.Token_ADD, rhino.Token_POS:
		return "+", true
	case rhino.Token_SUB, rhino.Token_NEG:
		return "-", true
	case rhino.Token_ASSIGN:
		return "=", true
	case rhino.Token_ASSIGN_BITOR:
		return "|=", true
	case rhino.Token_ASSIGN_BITXOR:
		return "^=", true
	case rhino.Token_ASSIGN_BITAND:
		return "&=", true
	case rhino.Token_ASSIGN_LSH:
		return "<<=", true
	case rhino.Token_ASSIGN_RSH:
		return ">>=", true
	case rhino.Token_ASSIGN_URSH:
		return ">>>=", true
	case rhino.Token_ASSIGN_ADD:
		return "+=", true
	case rhino.Token_ASSIGN_SUB:
		return "-=", true
	case rhino.Token_ASSIGN_MUL:
		return "*=", true
	case rhino.Token_ASSIGN_EXPONENT:
		return "**=", true
	case rhino.Token_ASSIGN_DIV:
		return "/=", true
	case rhino.Token_ASSIGN_MOD:
		return "%=", true
	case rhino.Token_ASSIGN_OR:
		return "||=", true
	case rhino.Token_ASSIGN_AND:
		return "&&=", true
	case rhino.Token_ASSIGN_COALESCE:
		return "??=", true
	case rhino.Token_VOID:
		return "void", true
	case rhino.Token_TYPEOF:
		return "typeof", true
	case rhino.Token_INSTANCEOF:
		return "instanceof", true
	default:
		return "", false
	}
}

// NodeUtil_OpToStrNoFail converts an operator's token value (see Token) to a
// string representation or fails.
//
// operator is the operator's token value to convert; returns the string
// representation. Throws Error if the token value is not an operator.
//
// Java: NodeUtil.opToStrNoFail(Token)
func NodeUtil_OpToStrNoFail(operator rhino.Token) string {
	res := NodeUtil_OpToStr(operator)
	if res == nil {
		panic(jutil.NewJavaLangError("Unknown op " + operator.String()))
	}
	return *res
}

// NodeUtil_RedeclareVarsInsideBranch: given a node tree, finds all the VAR
// declarations in that tree that are not in an inner scope. Then adds a new
// VAR node at the top of the current scope that redeclares them, if
// necessary.
//
// Java: NodeUtil.redeclareVarsInsideBranch(Node)
func NodeUtil_RedeclareVarsInsideBranch(branch rhino.Node) {
	if !NodeUtil_CanContainHoistedVarsDecls(branch) {
		return
	}

	vars := NodeUtil_GetVarsDeclaredInBranch(branch)
	if len(vars) == 0 {
		return
	}

	parent := NodeUtil_GetAddingRoot(branch)
	for _, nameNode := range vars {
		var_ := rhino.IR_Var(rhino.IR_Name(nameNode.GetString()).Srcref(nameNode)).Srcref(nameNode)
		NodeUtil_CopyNameAnnotations(nameNode, var_.GetFirstChild())
		parent.AddChildToFront(var_)
	}
}

// NodeUtil_GetRootOfQualifiedNameString gets the root node of a string
// representing a qualified name.
//
// Java: NodeUtil.getRootOfQualifiedName(String)
func NodeUtil_GetRootOfQualifiedNameString(qName string) string {
	dot := jutil.StringIndexOfChar(qName, '.')
	if dot == -1 {
		return qName
	}
	return jutil.StringSubstring(qName, 0, dot)
}

// Java: NodeUtil.newName(AbstractCompiler, String)
func NodeUtil_NewName(compiler AbstractCompiler, name string) rhino.Node {
	nameNode := rhino.IR_Name(name)
	nameNode.SetLength(jutil.StringLength(name))
	if compiler.GetCodingConvention().IsConstant(name) {
		nameNode.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	}
	return nameNode
}

// NodeUtil_NewNameWithSrcref creates a new node representing an *existing*
// name, copying over the source location information from the basis node.
//
// name is the name for the new NAME node; srcref is the node that represents
// the name as currently found in the AST. Returns the node created.
//
// Java: NodeUtil.newName(AbstractCompiler, String, Node)
func NodeUtil_NewNameWithSrcref(compiler AbstractCompiler, name string, srcref rhino.Node) rhino.Node {
	return NodeUtil_NewName(compiler, name).Srcref(srcref)
}

// NodeUtil_NewNameWithOriginalName creates a new node representing an
// *existing* name, copying over the source location information from the
// basis node and assigning the given original name to the node.
//
// name is the name for the new NAME node; basisNode is the node that
// represents the name as currently found in the AST; originalName is the
// original name of the item being represented by the NAME node. Used for
// debugging information. Returns the node created.
//
// Java: NodeUtil.newName(AbstractCompiler, String, Node, String)
func NodeUtil_NewNameWithOriginalName(
	compiler AbstractCompiler, name string, basisNode rhino.Node, originalName string) rhino.Node {
	nameNode := NodeUtil_NewNameWithSrcref(compiler, name, basisNode)
	nameNode.SetOriginalNameString(originalName)
	return nameNode
}

// NodeUtil_IsLatin tests if all characters in the string are in the Basic
// Latin (aka ASCII) character set - that they have UTF-16 values equal to or
// below 0x7f. This check can find which identifiers with Unicode characters
// need to be escaped in order to allow resulting files to be processed by
// non-Unicode aware UNIX tools and editors. See
// http://en.wikipedia.org/wiki/Latin_characters_in_Unicode for more on Basic
// Latin.
//
// s is the string to be checked for ASCII-goodness. Returns true if all
// characters in the string are in Basic Latin set.
//
// Java: NodeUtil.isLatin(String)
func NodeUtil_IsLatin(s string) bool {
	// Go-only: language-adapter — a UTF-16 char above 0x7f exists iff the
	// WTF-8 encoding has a byte above 0x7f, so the bytes are scanned directly.
	len_ := len(s)
	for index := 0; index < len_; index++ {
		c := s[index]
		if uint16(c) > NodeUtil_LARGEST_BASIC_LATIN {
			return false
		}
	}
	return true
}

// NodeUtil_IsValidSimpleName determines whether the given name is a valid
// variable name.
//
// Java: NodeUtil.isValidSimpleName(String)
func NodeUtil_IsValidSimpleName(name string) bool {
	return rhino.TokenStream_IsJSIdentifier(name) &&
		!rhino.TokenStream_IsKeyword(name) &&
		// no Unicode escaped characters - some browsers are less tolerant
		// of Unicode characters that might be valid according to the
		// language spec.
		// Note that by this point, Unicode escapes have been converted
		// to UTF-16 characters, so we're only searching for character
		// values, not escapes.
		NodeUtil_IsLatin(name)
}

// Deprecated: use NodeUtil_IsValidQualifiedNameFeatureSet(mode.ToFeatureSet(), name).
//
// Java: NodeUtil.isValidQualifiedName(LanguageMode, String)
func NodeUtil_IsValidQualifiedName(mode CompilerOptions_LanguageMode, name string) bool {
	return NodeUtil_IsValidQualifiedNameFeatureSet(mode.ToFeatureSet(), name)
}

// NodeUtil_IsValidQualifiedNameFeatureSet determines whether the given name
// is a valid qualified name.
//
// Java: NodeUtil.isValidQualifiedName(FeatureSet, String)
func NodeUtil_IsValidQualifiedNameFeatureSet(mode *parser.FeatureSet, name string) bool {
	if strings.HasSuffix(name, ".") || strings.HasPrefix(name, ".") {
		return false
	}

	// Go-only: language-adapter — Splitter.on('.').splitToList(name) keeps
	// empty parts, exactly like strings.Split with a one-character separator.
	parts := strings.Split(name, ".")
	for _, part := range parts {
		if !NodeUtil_IsValidPropertyName(mode, part) {
			return false
		}
	}
	return NodeUtil_IsValidSimpleName(parts[0])
}

// NodeUtil_IsValidPropertyName determines whether the given name can appear
// on the right side of the dot operator. Many properties (like reserved
// words) cannot, in ES3.
//
// Java: NodeUtil.isValidPropertyName(FeatureSet, String)
func NodeUtil_IsValidPropertyName(mode *parser.FeatureSet, name string) bool {
	if NodeUtil_IsValidSimpleName(name) {
		return true
	} else {
		return mode.Has(parser.FeatureSet_Feature_KEYWORDS_AS_PROPERTIES) && rhino.TokenStream_IsKeyword(name)
	}
}

// NodeUtil_VarCollector collects the names declared by VARs.
//
// Java: NodeUtil.VarCollector
type NodeUtil_VarCollector struct {
	vars *jutil.LinkedHashMap[string, rhino.Node]
}

// Java: NodeUtil.GLOBAL_OBJECT
var NodeUtil_GLOBAL_OBJECT = rhino.QualifiedName_Of("$jscomp.global.Object")

// Java: NodeUtil.GLOBAL_OBJECT_MANGLED
var NodeUtil_GLOBAL_OBJECT_MANGLED = rhino.QualifiedName_Of("$jscomp$global.Object")

// Java: NodeUtil.NUMBER_NAN
var NodeUtil_NUMBER_NAN = rhino.IR_Getprop(rhino.IR_Name("Number"), "NaN")

// Java: NodeUtil.GOOG_LOADMODULE
var NodeUtil_GOOG_LOADMODULE = rhino.QualifiedName_Of("goog.loadModule")

// NodeUtil_GetOriginalName returns the original name of a renamed variable
// by stripping unique suffixes and module prefixes.
//
// Java: NodeUtil.getOriginalName(String)
func NodeUtil_GetOriginalName(name string) string {
	// Strip unique name suffix: "foo$$1" -> "foo"
	name = jutil.StringReplaceAll(name, "\\$\\$\\d+$", "")

	// Strip contextual rename suffix: "foo$jscomp$1" -> "foo"
	name = MakeDeclaredNamesUnique_ContextualRenameInverter_GetOriginalName(name)

	// Strip module exports prefix: "module$exports$foo$Bar" -> "Bar"
	if strings.HasPrefix(name, "module$exports$") {
		lastDollar := jutil.StringLastIndexOf(name, "$")
		if lastDollar != -1 {
			name = jutil.StringSubstringFrom(name, lastDollar+1)
		}
	} else if strings.HasPrefix(name, "module$contents$") {
		// Strip module contents prefix: "module$contents$foo$Bar_baz" -> "baz"
		lastUnderscore := jutil.StringLastIndexOf(name, "_")
		if lastUnderscore != -1 {
			name = jutil.StringSubstringFrom(name, lastUnderscore+1)
		}
	}
	return name
}

// nodeUtilFormatInt formats v in base 10 (the uniqueNameIdSupplier hot path).
func nodeUtilFormatInt(v int) string {
	return strconv.Itoa(v)
}
