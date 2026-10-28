// Port of com/google/javascript/jscomp/ijs/PotentialDeclaration.java @ closure-compiler v20261021

/*
 * Copyright 2017 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type potentialDeclKind uint8

const (
	potentialDeclName potentialDeclKind = iota
	potentialDeclThisProp
	potentialDeclDefine
	potentialDeclMethod
	potentialDeclMemberField
	potentialDeclStringKey
	potentialDeclPolymerBehaviorProperties
	potentialDeclAlias
)

// PotentialDeclaration encapsulates something that could be a declaration.
//
// This includes: var/let/const declarations, function/class declarations,
// method declarations, assignments, goog.define calls, and even valueless
// property accesses (e.g. `/** @type {number} */ Foo.prototype.bar`).
type PotentialDeclaration struct {
	// The fully qualified name of the declaration.
	fullyQualifiedName string
	// The LHS node of the declaration.
	lhs rhino.Node
	// The RHS node of the declaration, if it exists.
	rhs rhino.Node
	// Statement insertion anchor for `this.prop` declarations inside constructors.
	insertionPoint rhino.Node
	kind           potentialDeclKind
}

// The threshold for truncating string enum items. Any string literal longer
// than this will get truncated in the generated IJS file. This value is
// decided by the length of the longest security-sensitive attribute name.
const potentialDeclaration_STRING_ENUM_RETAIN_CAP = 10

func newPotentialDeclaration(kind potentialDeclKind, fullyQualifiedName string, lhs rhino.Node, rhs rhino.Node) *PotentialDeclaration {
	return &PotentialDeclaration{
		fullyQualifiedName: fullyQualifiedName,
		lhs:                rhino.CheckNotNull(lhs),
		rhs:                rhs,
		kind:               kind,
	}
}

// Java: PotentialDeclaration.fromName(Node)
func PotentialDeclaration_FromName(nameNode rhino.Node) *PotentialDeclaration {
	jutil.CheckArgumentMsg1(nameNode.IsQualifiedName(), "%s", nameNode)
	rhs := NodeUtil_GetRValueOfLValue(nameNode)
	if ClassUtil_IsThisPropInsideClassWithName(nameNode) {
		name := ClassUtil_GetFullyQualifiedNameOfThisProp(nameNode)
		d := newPotentialDeclaration(potentialDeclThisProp, name, nameNode, rhs)
		thisPropDefinition := NodeUtil_GetEnclosingStatement(nameNode)
		d.insertionPoint = NodeUtil_GetEnclosingStatement(thisPropDefinition.GetParent())
		return d
	}
	return newPotentialDeclaration(potentialDeclName, *nameNode.GetQualifiedName(), nameNode, rhs)
}

// Java: PotentialDeclaration.fromMemberFieldDef(Node)
func PotentialDeclaration_FromMemberFieldDef(fieldNode rhino.Node) *PotentialDeclaration {
	jutil.CheckArgument(fieldNode.IsMemberFieldDef())
	rhs := NodeUtil_GetRValueOfLValue(fieldNode)
	name := ClassUtil_GetFullyQualifiedNameOfMemberFieldDef(fieldNode)
	// lhs is MEMBER_FIELD_DEF; rhs is the value assigned to the field if any.
	return newPotentialDeclaration(potentialDeclMemberField, name, fieldNode, rhs)
}

// Java: PotentialDeclaration.fromMethod(Node)
func PotentialDeclaration_FromMethod(functionNode rhino.Node) *PotentialDeclaration {
	jutil.CheckArgument(ClassUtil_IsClassMethod(functionNode))
	name := ClassUtil_GetFullyQualifiedNameOfMethod(functionNode)
	return newPotentialDeclaration(potentialDeclMethod, name, functionNode.GetParent(), functionNode)
}

// Java: PotentialDeclaration.fromStringKey(Node)
func PotentialDeclaration_FromStringKey(stringKeyNode rhino.Node) *PotentialDeclaration {
	sk, ok := ast.AsStringKey(stringKeyNode)
	jutil.CheckArgument(ok)
	jutil.CheckArgument(stringKeyNode.GetParent().IsObjectLit())
	keyName := sk.Key()
	name := "this." + keyName
	kind := potentialDeclStringKey
	if keyName == "properties" {
		objLitJsDoc := NodeUtil_GetBestJSDocInfo(stringKeyNode.GetParent())
		if objLitJsDoc != nil && objLitJsDoc.IsPolymerBehavior() {
			kind = potentialDeclPolymerBehaviorProperties
		}
	}
	return newPotentialDeclaration(kind, name, stringKeyNode, stringKeyNode.GetLastChild())
}

// Java: PotentialDeclaration.fromDefine(Node)
func PotentialDeclaration_FromDefine(callNode rhino.Node) *PotentialDeclaration {
	jutil.CheckArgument(NodeUtil_IsCallTo(callNode, "goog.define"))
	call := ast.MustCall(callNode)
	lastArg, _ := call.LastArg()
	rhs := lastArg.Node()
	// Match a few different forms, depending on the call node's parent:
	//   1. EXPR_RESULT: goog.define('foo', 1);
	//   2. ASSIGN: a.b = goog.define('c', 2);
	//   3. NAME: var x = goog.define('d', 3);
	parent := callNode.GetParent()
	switch parent.GetToken() {
	case rhino.Token_EXPR_RESULT:
		return newPotentialDeclaration(potentialDeclDefine, call.MustFirstArg().Node().GetString(), callNode, rhs)
	case rhino.Token_ASSIGN:
		previous := ast.MustAssign(parent).Target().Node()
		return newPotentialDeclaration(potentialDeclDefine, *previous.GetQualifiedName(), previous, rhs)
	case rhino.Token_NAME:
		return newPotentialDeclaration(potentialDeclDefine, parent.GetString(), parent, rhs)
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected parent: " + parent.GetToken().String()))
	}
}

// Java: PotentialDeclaration.fromAlias(Node)
func PotentialDeclaration_FromAlias(nameNode rhino.Node) *PotentialDeclaration {
	jutil.CheckArgumentMsg1(nameNode.IsQualifiedName(), "%s", nameNode)
	return newPotentialDeclaration(potentialDeclAlias, *nameNode.GetQualifiedName(), nameNode, rhino.NilNode)
}

// Java: PotentialDeclaration.getFullyQualifiedName()
func (d *PotentialDeclaration) GetFullyQualifiedName() string {
	return d.fullyQualifiedName
}

// Java: PotentialDeclaration.getLhs()
func (d *PotentialDeclaration) GetLhs() rhino.Node {
	return d.lhs
}

// Java: PotentialDeclaration.getRhs()
func (d *PotentialDeclaration) GetRhs() rhino.Node {
	return d.rhs
}

// Java: PotentialDeclaration.getJsDoc()
func (d *PotentialDeclaration) GetJsDoc() *rhino.JSDocInfo {
	return NodeUtil_GetBestJSDocInfo(d.lhs)
}

// Java: PotentialDeclaration.isDetached()
func (d *PotentialDeclaration) IsDetached() bool {
	for current := d.lhs; !current.IsNil(); current = current.GetParent() {
		if current.IsScript() {
			return false
		}
	}
	return true
}

// Java: PotentialDeclaration.getRemovableNode()
func (d *PotentialDeclaration) GetRemovableNode() rhino.Node {
	switch d.kind {
	case potentialDeclMethod,
		potentialDeclMemberField,
		potentialDeclStringKey,
		potentialDeclPolymerBehaviorProperties:
		return d.lhs
	case potentialDeclAlias:
		// If the declaration is a destructuring declaration:
		// 1) If the lhs's destructuring pattern parent has only one child, e.g.
		//    const {Foo} = x;, return the enclosing statement to remove it all.
		// 2) If the parent has more than one child, e.g. const {Foo, Bar} = x;,
		//    return the lhs (or its STRING_KEY parent) so only Foo is removed.
		if d.lhs.GetParent().IsArrayPattern() && d.lhs.GetParent().HasMoreThanOneChild() {
			return d.lhs
		}
		if d.lhs.GetGrandparent().IsObjectPattern() && d.lhs.GetGrandparent().HasMoreThanOneChild() {
			return d.lhs.GetParent()
		}
		return NodeUtil_GetEnclosingStatement(d.lhs)
	default:
		return NodeUtil_GetEnclosingStatement(d.lhs)
	}
}

// Remove removes this "potential declaration" completely. Usually, this is
// because the same symbol has already been declared in this file.
//
// Java: PotentialDeclaration.remove(AbstractCompiler)
func (d *PotentialDeclaration) Remove(compiler AbstractCompiler) {
	if d.IsDetached() {
		return
	}
	statement := d.GetRemovableNode()
	NodeUtil_DeleteNode(statement, compiler)
	statement.RemoveChildren()
}

// BreakDownDestructure breaks down this declaration if it's a destructuring
// LHS by replacing it with a VAR node and removing the RHS.
//
// Java: PotentialDeclaration.breakDownDestructure(AbstractCompiler)
func (d *PotentialDeclaration) BreakDownDestructure(compiler AbstractCompiler) bool {
	if d.kind != potentialDeclName || !NodeUtil_IsLhsByDestructuring(d.lhs) {
		return false
	}

	rootTarget := NodeUtil_GetRootTarget(d.lhs)
	jutil.CheckState(rootTarget.GetParent().IsDestructuringLhs())
	definitionNode := rootTarget.GetGrandparent()
	if !NodeUtil_IsNameDeclaration(definitionNode) {
		return false
	}

	var prev rhino.Node
	var lhsNodes []rhino.Node
	NodeUtil_VisitLhsNodesInNode(rootTarget.GetParent(), func(n rhino.Node) { lhsNodes = append(lhsNodes, n) })
	for _, n := range lhsNodes {
		n.Detach()
		temp := rhino.IR_Var(n)
		if prev.IsNil() {
			definitionNode.ReplaceWith(temp)
			compiler.ReportChangeToEnclosingScope(temp)
		} else {
			temp.InsertAfter(prev)
			compiler.ReportChangeToEnclosingScope(temp)
			temp.SrcrefTree(prev)
		}
		prev = temp
	}
	return true
}

// Java: PotentialDeclaration.simplify(AbstractCompiler)
func (d *PotentialDeclaration) Simplify(compiler AbstractCompiler) {
	switch d.kind {
	case potentialDeclName:
		d.simplifyName(compiler)
	case potentialDeclThisProp:
		d.simplifyThisProp(compiler)
	case potentialDeclDefine:
		d.simplifyDefine(compiler)
	case potentialDeclMemberField:
		if !d.rhs.IsNil() {
			NodeUtil_DeleteNode(d.rhs, compiler)
		}
	case potentialDeclStringKey:
		d.simplifyStringKey(compiler)
	case potentialDeclPolymerBehaviorProperties:
		d.simplifyPolymerBehaviorProperties(compiler)
	case potentialDeclMethod, potentialDeclAlias:
		// Does not simplify.
	}
}

func (d *PotentialDeclaration) simplifyName(compiler AbstractCompiler) {
	if d.rhs.IsNil() || d.ShouldPreserve() {
		return
	}
	nameNode := d.lhs
	jsdoc := d.GetJsDoc()
	if jsdoc != nil && jsdoc.HasEnumParameterType() {
		d.simplifyEnumValues(compiler)
		return
	}
	if NodeUtil_IsNamespaceDecl(nameNode) {
		if d.rhs.IsOr() {
			objLit := d.rhs.GetLastChild().Detach()
			d.rhs.ReplaceWith(objLit)
			compiler.ReportChangeToEnclosingScope(d.lhs)
		}
		return
	}
	if nameNode.MatchesName("exports") {
		// Replace the RHS of a default goog.module export with Unknown.
		d.rhs.ReplaceWith(rhino.IR_Cast(rhino.IR_Number(0), JsdocUtil_GetQmarkTypeJSDoc()).SrcrefTree(d.rhs))
		compiler.ReportChangeToEnclosingScope(nameNode)
		return
	}
	if NodeUtil_IsCallTo(d.rhs, "Symbol") {
		for d.rhs.HasMoreThanOneChild() {
			NodeUtil_DeleteNode(d.rhs.GetLastChild(), compiler)
		}
		return
	}
	if d.lhs.GetParent().IsConst() {
		jsdoc = JsdocUtil_MarkConstant(jsdoc)
	}
	// Completely remove the RHS, and replace with a getprop.
	newStatement := NodeUtil_NewQNameDeclaration(compiler, *nameNode.GetQualifiedName(), rhino.NilNode, jsdoc)
	newStatement.SrcrefTreeIfMissing(nameNode)
	oldStatement := d.GetRemovableNode()
	NodeUtil_DeleteChildren(oldStatement, compiler)
	if oldStatement.IsExport() {
		oldStatement.AddChildToBack(newStatement)
	} else {
		oldStatement.ReplaceWith(newStatement)
	}
	compiler.ReportChangeToEnclosingScope(newStatement)
}

func (d *PotentialDeclaration) simplifyThisProp(compiler AbstractCompiler) {
	if d.ShouldPreserve() {
		return
	}
	// Completely remove the RHS, if present, and replace with a getprop.
	newStatement := NodeUtil_NewQNameDeclaration(compiler, d.fullyQualifiedName, rhino.NilNode, d.GetJsDoc())
	newStatement.SrcrefTreeIfMissing(d.lhs)
	NodeUtil_DeleteNode(d.GetRemovableNode(), compiler)
	if d.insertionPoint.HasParent() {
		newStatement.InsertAfter(d.insertionPoint)
		compiler.ReportChangeToEnclosingScope(newStatement)
	}
}

func (d *PotentialDeclaration) simplifyDefine(compiler AbstractCompiler) {
	info := d.GetJsDoc()
	if info != nil && info.GetType() != nil {
		newRhs := potentialDeclaration_makeEmptyValueNode(info.GetType())
		if !newRhs.IsNil() {
			d.rhs.ReplaceWith(newRhs)
			compiler.ReportChangeToEnclosingScope(newRhs)
			return
		}
	}
	NodeUtil_DeleteNode(d.GetRemovableNode(), compiler)
}

func potentialDeclaration_makeEmptyValueNode(type_ *rhino.JSTypeExpression) rhino.Node {
	n := type_.GetRoot()
	for !n.IsNil() && !n.IsStringLit() && !n.IsName() {
		n = n.GetFirstChild()
	}
	if n.IsNil() {
		return rhino.NilNode
	}
	switch n.GetString() {
	case "boolean":
		return rhino.NewNode(rhino.Token_FALSE)
	case "number":
		return rhino.Node_NewNumber(0)
	case "string":
		return rhino.Node_NewString("")
	default:
		return rhino.NilNode
	}
}

func (d *PotentialDeclaration) simplifyStringKey(compiler AbstractCompiler) {
	if d.ShouldPreserve() {
		return
	}
	jsdoc := d.GetJsDoc()
	if jsdoc != nil && jsdoc.HasEnumParameterType() {
		d.simplifyEnumValues(compiler)
		return
	}
	PotentialDeclaration_RemoveStringKeyValue(d.lhs)
	compiler.ReportChangeToEnclosingScope(d.lhs)
	if jsdoc == nil || !jsdoc.ContainsDeclaration() || d.IsConstToBeInferred() {
		d.lhs.SetJSDocInfo(JsdocUtil_GetUnusableTypeJSDoc(jsdoc))
	}
}

// simplifyPolymerBehaviorProperties preserves the "properties" object on a
// @polymerBehavior literal while stripping non-"type"/non-"readOnly" keys.
func (d *PotentialDeclaration) simplifyPolymerBehaviorProperties(compiler AbstractCompiler) {
	if d.IsDetached() {
		return
	}
	propertiesObject, ok := ast.AsObjectLit(d.rhs)
	if !ok || !propertiesObject.HasMembers() {
		return
	}
	for prop := range propertiesObject.Members() {
		propDef := prop.Node().GetOnlyChild()
		// A property definition is either a function reference (e.g. String, Number), or another
		// object literal. If it's an object literal, only the "type" and "readOnly" sub-properties
		// matter for type checking, so we can delete everything else.
		if defObj, isObj := ast.AsObjectLit(propDef); isObj {
			for subProp := range ast.ReadAhead(defObj.FirstMember()) {
				subNode := subProp.Node()
				if subNode.GetString() != "type" && subNode.GetString() != "readOnly" {
					NodeUtil_DeleteNode(subNode, compiler)
				}
			}
		}
	}
}

// simplifyEnumValues removes values from enums.
func (d *PotentialDeclaration) simplifyEnumValues(compiler AbstractCompiler) {
	objLit, ok := ast.AsObjectLit(d.rhs)
	if !ok || !objLit.HasMembers() {
		return
	}
	changed := false
	for member := range objLit.Members() {
		key := member.Node()
		value := key.GetOnlyChild()
		if !value.IsStringLit() {
			PotentialDeclaration_RemoveStringKeyValue(key)
			changed = true
		} else if jutil.StringLength(value.GetString()) > potentialDeclaration_STRING_ENUM_RETAIN_CAP {
			PotentialDeclaration_TruncateStringKeyValue(key)
			changed = true
		}
	}
	if changed {
		compiler.ReportChangeToEnclosingScope(d.rhs)
	}
}

// Java: PotentialDeclaration.isDefiniteDeclaration(AbstractCompiler)
func (d *PotentialDeclaration) IsDefiniteDeclaration(compiler AbstractCompiler) bool {
	if d.kind == potentialDeclAlias {
		return true
	}
	parent := d.lhs.GetParent()
	switch parent.GetToken() {
	case rhino.Token_DEFAULT_VALUE:
		if !parent.GetParent().IsStringKey() {
			return false
		}
		fallthrough
	case rhino.Token_COMPUTED_PROP, rhino.Token_STRING_KEY:
		if NodeUtil_IsLhsByDestructuring(d.lhs) {
			rootTarget := NodeUtil_GetRootTarget(d.lhs)
			jutil.CheckState(rootTarget.GetParent().IsDestructuringLhs())
			if NodeUtil_IsNameDeclaration(rootTarget.GetGrandparent()) {
				return true
			}
		}
		return false
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST, rhino.Token_CLASS, rhino.Token_FUNCTION:
		return true
	default:
		jsdoc := d.GetJsDoc()
		return PotentialDeclaration_IsExportLhs(d.lhs) ||
			(jsdoc != nil && jsdoc.ContainsDeclaration()) ||
			(!d.rhs.IsNil() && PotentialDeclaration_IsTypedRhs(d.rhs))
	}
}

func (d *PotentialDeclaration) shouldPreserveBase() bool {
	return !d.rhs.IsNil() && PotentialDeclaration_IsTypedRhs(d.rhs)
}

// Java: PotentialDeclaration.shouldPreserve()
func (d *PotentialDeclaration) ShouldPreserve() bool {
	switch d.kind {
	case potentialDeclPolymerBehaviorProperties, potentialDeclAlias:
		return true
	case potentialDeclStringKey:
		return d.IsDetached() || d.shouldPreserveBase() || !d.isStringKeyInNamespace()
	case potentialDeclName:
		rhs := d.rhs
		nameNode := d.lhs
		jsdoc := d.GetJsDoc()
		isExport := PotentialDeclaration_IsExportLhs(nameNode)
		return d.shouldPreserveBase() ||
			PotentialDeclaration_IsImportRhs(rhs) ||
			(isExport && !rhs.IsNil() && (rhs.IsQualifiedName() || rhs.IsObjectLit())) ||
			(jsdoc != nil && jsdoc.IsConstructor() && !rhs.IsNil() && rhs.IsQualifiedName()) ||
			(!rhs.IsNil() &&
				rhs.IsObjectLit() &&
				!rhs.HasChildren() &&
				(jsdoc == nil || !JsdocUtil_HasAnnotatedType(jsdoc))) ||
			(!rhs.IsNil() && NodeUtil_IsCallTo(rhs, "Polymer")) ||
			d.isPolymerBehaviorAliasOrArray(jsdoc)
	default:
		return d.shouldPreserveBase()
	}
}

func (d *PotentialDeclaration) isPolymerBehaviorAliasOrArray(jsdoc *rhino.JSDocInfo) bool {
	return jsdoc != nil &&
		jsdoc.IsPolymerBehavior() &&
		!d.rhs.IsNil() &&
		(d.rhs.IsName() || d.rhs.IsArrayLit())
}

func (d *PotentialDeclaration) isStringKeyInNamespace() bool {
	objLit := d.lhs.GetParent()
	lvalue := NodeUtil_GetBestLValue(objLit)
	if lvalue.IsNil() {
		return false
	}
	jsdoc := NodeUtil_GetBestJSDocInfo(lvalue)
	return !PotentialDeclaration_IsExportLhs(lvalue) &&
		!JsdocUtil_HasAnnotatedType(jsdoc) &&
		NodeUtil_IsNamespaceDecl(lvalue)
}

// Java: PotentialDeclaration.isConstToBeInferred()
func (d *PotentialDeclaration) IsConstToBeInferred() bool {
	return PotentialDeclaration_IsConstToBeInferred(d.lhs)
}

// Java: PotentialDeclaration.isConstToBeInferred(Node)
func PotentialDeclaration_IsConstToBeInferred(nameNode rhino.Node) bool {
	jsdoc := NodeUtil_GetBestJSDocInfo(nameNode)
	isConst :=
		nameNode.GetParent().IsConst() ||
			PotentialDeclaration_IsExportLhs(nameNode) ||
			(jsdoc != nil && jsdoc.IsConstant())
	return isConst && !JsdocUtil_HasAnnotatedType(jsdoc) && !NodeUtil_IsNamespaceDecl(nameNode)
}

var potentialDeclaration_GOOG_ABSTRACTMETHOD = rhino.QualifiedName_Of("goog.abstractMethod")
var potentialDeclaration_GOOG_REQUIRE = rhino.QualifiedName_Of("goog.require")
var potentialDeclaration_GOOG_REQUIRETYPE = rhino.QualifiedName_Of("goog.requireType")
var potentialDeclaration_GOOG_FORWARDDECLARE = rhino.QualifiedName_Of("goog.forwardDeclare")
var potentialDeclaration_MODULE_EXPORTS = rhino.QualifiedName_Of("module.exports")

// Java: PotentialDeclaration.isTypedRhs(Node)
func PotentialDeclaration_IsTypedRhs(rhs rhino.Node) bool {
	return rhs.IsFunction() ||
		rhs.IsClass() ||
		(rhs.IsQualifiedName() && potentialDeclaration_GOOG_ABSTRACTMETHOD.Matches(rhs))
}

// Java: PotentialDeclaration.isExportLhs(Node)
func PotentialDeclaration_IsExportLhs(lhs rhino.Node) bool {
	if lhs.IsName() && lhs.MatchesName("exports") {
		return true
	}
	if gp, ok := ast.AsGetProp(lhs); ok && gp.Object().Node().MatchesName("exports") {
		return true
	}
	return potentialDeclaration_MODULE_EXPORTS.Matches(lhs)
}

// Java: PotentialDeclaration.isImportRhs(Node)
func PotentialDeclaration_IsImportRhs(rhs rhino.Node) bool {
	call, ok := ast.AsCall(rhs)
	if !ok {
		return false
	}
	callee := call.Callee().Node()
	return potentialDeclaration_GOOG_REQUIRE.Matches(callee) ||
		potentialDeclaration_GOOG_REQUIRETYPE.Matches(callee) ||
		potentialDeclaration_GOOG_FORWARDDECLARE.Matches(callee) ||
		callee.MatchesName("require")
}

// Java: PotentialDeclaration.isAliasDeclaration(Node, Node)
func PotentialDeclaration_IsAliasDeclaration(lhs rhino.Node, rhs rhino.Node) bool {
	return !ClassUtil_IsThisPropInsideClassWithName(lhs) &&
		PotentialDeclaration_IsConstToBeInferred(lhs) &&
		!rhs.IsNil() &&
		PotentialDeclaration_IsQualifiedAliasExpression(rhs)
}

var potentialDeclaration_GOOG_MODULE_GET = rhino.QualifiedName_Of("goog.module.get")

// PotentialDeclaration_IsQualifiedAliasExpression returns whether a node
// corresponds to a simple or a qualified name, such as x or a.b.c or this.a.
//
// Java: PotentialDeclaration.isQualifiedAliasExpression(Node)
func PotentialDeclaration_IsQualifiedAliasExpression(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_NAME:
		return n.GetString() != ""
	case rhino.Token_THIS, rhino.Token_SUPER:
		return true
	case rhino.Token_GETPROP:
		return PotentialDeclaration_IsQualifiedAliasExpression(ast.MustGetProp(n).Object().Node())
	case rhino.Token_CALL:
		return potentialDeclaration_GOOG_MODULE_GET.Matches(ast.MustCall(n).Callee().Node())
	case rhino.Token_MEMBER_FUNCTION_DEF:
		// These are explicitly *not* qualified name components.
		return false
	default:
		return false
	}
}

// Java: PotentialDeclaration.removeStringKeyValue(Node)
func PotentialDeclaration_RemoveStringKeyValue(stringKey rhino.Node) {
	value := stringKey.GetOnlyChild()
	replacementValue := rhino.IR_Number(0).SrcrefTree(value)
	value.ReplaceWith(replacementValue)
}

// Java: PotentialDeclaration.truncateStringKeyValue(Node)
func PotentialDeclaration_TruncateStringKeyValue(stringKey rhino.Node) {
	value := stringKey.GetOnlyChild()
	replacementValue :=
		rhino.IR_String(jutil.StringSubstring(value.GetString(), 0, potentialDeclaration_STRING_ENUM_RETAIN_CAP-2) + "..").
			SrcrefTree(value)
	value.ReplaceWith(replacementValue)
}
