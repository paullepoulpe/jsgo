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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// PotentialDeclaration encapsulates something that could be a declaration.
//
// This includes: var/let/const declarations, function/class declarations,
// method declarations, assignments, goog.define calls, and even valueless
// property accesses (e.g. `/** @type {number} * / Foo.prototype.bar`)
type PotentialDeclaration interface {
	GetPotentialDeclarationBase() *PotentialDeclarationBase

	// Java: PotentialDeclaration.getFullyQualifiedName()
	GetFullyQualifiedName() string
	// Java: PotentialDeclaration.getLhs()
	GetLhs() *rhino.Node
	// Java: PotentialDeclaration.getRhs()
	GetRhs() *rhino.Node
	// Java: PotentialDeclaration.getJsDoc()
	GetJsDoc() *rhino.JSDocInfo
	// Java: PotentialDeclaration.isDetached()
	IsDetached() bool
	// Java: PotentialDeclaration.getRemovableNode()
	GetRemovableNode() *rhino.Node
	// Java: PotentialDeclaration.remove(AbstractCompiler)
	Remove(compiler AbstractCompiler)
	// Java: PotentialDeclaration.simplify(AbstractCompiler)
	Simplify(compiler AbstractCompiler)
	// Java: PotentialDeclaration.breakDownDestructure(AbstractCompiler)
	BreakDownDestructure(compiler AbstractCompiler) bool
	// Java: PotentialDeclaration.isDefiniteDeclaration(AbstractCompiler)
	IsDefiniteDeclaration(compiler AbstractCompiler) bool
	// Java: PotentialDeclaration.shouldPreserve()
	ShouldPreserve() bool
	// Java: PotentialDeclaration.isConstToBeInferred()
	IsConstToBeInferred() bool
}

// PotentialDeclarationBase holds the fields and base methods of the
// abstract class PotentialDeclaration.
type PotentialDeclarationBase struct {
	self PotentialDeclaration

	// The fully qualified name of the declaration.
	fullyQualifiedName string
	// The LHS node of the declaration.
	lhs *rhino.Node
	// The RHS node of the declaration, if it exists.
	rhs *rhino.Node
}

// The threshold for trunctating string enum items. Any string literal longer
// than this will get truncated in the generated IJS file. This value is
// decided by the length of the longest security-sensitive attribute name.
const potentialDeclaration_STRING_ENUM_RETAIN_CAP = 10

// Java: PotentialDeclaration(String, Node, Node)
func (b *PotentialDeclarationBase) InitPotentialDeclarationBase(self PotentialDeclaration, fullyQualifiedName string, lhs *rhino.Node, rhs *rhino.Node) {
	b.self = self
	b.fullyQualifiedName = fullyQualifiedName // checkNotNull: a Go string cannot be null
	b.lhs = jutil.CheckNotNull(lhs)
	b.rhs = rhs
}

// GetPotentialDeclarationBase gives access to the base through the interface.
func (b *PotentialDeclarationBase) GetPotentialDeclarationBase() *PotentialDeclarationBase { return b }

// Java: PotentialDeclaration.fromName(Node)
func PotentialDeclaration_FromName(nameNode *rhino.Node) PotentialDeclaration {
	jutil.CheckArgumentMsg1(nameNode.IsQualifiedName(), "%s", nameNode)
	rhs := NodeUtil_GetRValueOfLValue(nameNode)
	if ClassUtil_IsThisPropInsideClassWithName(nameNode) {
		name := ClassUtil_GetFullyQualifiedNameOfThisProp(nameNode)
		return NewPotentialDeclaration_ThisPropDeclaration(name, nameNode, rhs)
	}
	return NewPotentialDeclaration_NameDeclaration(*nameNode.GetQualifiedName(), nameNode, rhs)
}

// Java: PotentialDeclaration.fromMemberFieldDef(Node)
func PotentialDeclaration_FromMemberFieldDef(fieldNode *rhino.Node) PotentialDeclaration {
	jutil.CheckArgument(fieldNode.IsMemberFieldDef())
	rhs := NodeUtil_GetRValueOfLValue(fieldNode)
	name := ClassUtil_GetFullyQualifiedNameOfMemberFieldDef(fieldNode)
	// lhs is MEMBER_FIELD_DEF  rhs is the value assigned to the field if any.
	return NewPotentialDeclaration_MemberFieldDeclaration(name, fieldNode, rhs)
}

// Java: PotentialDeclaration.fromMethod(Node)
func PotentialDeclaration_FromMethod(functionNode *rhino.Node) PotentialDeclaration {
	jutil.CheckArgument(ClassUtil_IsClassMethod(functionNode))
	name := ClassUtil_GetFullyQualifiedNameOfMethod(functionNode)
	return NewPotentialDeclaration_MethodDeclaration(name, functionNode)
}

// Java: PotentialDeclaration.fromStringKey(Node)
func PotentialDeclaration_FromStringKey(stringKeyNode *rhino.Node) PotentialDeclaration {
	jutil.CheckArgument(stringKeyNode.IsStringKey())
	jutil.CheckArgument(stringKeyNode.GetParent().IsObjectLit())
	name := "this." + stringKeyNode.GetString()
	if stringKeyNode.GetString() == "properties" {
		objLitJsDoc := NodeUtil_GetBestJSDocInfo(stringKeyNode.GetParent())
		if objLitJsDoc != nil && objLitJsDoc.IsPolymerBehavior() {
			return NewPotentialDeclaration_PolymerBehaviorPropertiesDeclaration(name, stringKeyNode)
		}
	}
	return NewPotentialDeclaration_StringKeyDeclaration(name, stringKeyNode)
}

// Java: PotentialDeclaration.fromDefine(Node)
func PotentialDeclaration_FromDefine(callNode *rhino.Node) PotentialDeclaration {
	jutil.CheckArgument(NodeUtil_IsCallTo(callNode, "goog.define"))
	return PotentialDeclaration_DefineDeclaration_From(callNode)
}

// Java: PotentialDeclaration.fromAlias(Node)
func PotentialDeclaration_FromAlias(nameNode *rhino.Node) PotentialDeclaration {
	jutil.CheckArgumentMsg1(nameNode.IsQualifiedName(), "%s", nameNode)
	return NewPotentialDeclaration_AliasDeclaration(*nameNode.GetQualifiedName(), nameNode)
}

// Java: PotentialDeclaration.getFullyQualifiedName()
func (b *PotentialDeclarationBase) GetFullyQualifiedName() string {
	return b.fullyQualifiedName
}

// Java: PotentialDeclaration.getLhs()
func (b *PotentialDeclarationBase) GetLhs() *rhino.Node {
	return b.lhs
}

// Java: PotentialDeclaration.getRhs()
func (b *PotentialDeclarationBase) GetRhs() *rhino.Node {
	return b.rhs
}

// Java: PotentialDeclaration.getJsDoc()
func (b *PotentialDeclarationBase) GetJsDoc() *rhino.JSDocInfo {
	return NodeUtil_GetBestJSDocInfo(b.lhs)
}

// Java: PotentialDeclaration.isDetached()
func (b *PotentialDeclarationBase) IsDetached() bool {
	for current := b.lhs; current != nil; current = current.GetParent() {
		if current.IsScript() {
			return false
		}
	}
	return true
}

// Java: PotentialDeclaration.getRemovableNode()
func (b *PotentialDeclarationBase) GetRemovableNode() *rhino.Node {
	return NodeUtil_GetEnclosingStatement(b.lhs)
}

// Remove removes this "potential declaration" completely. Usually, this is
// because the same symbol has already been declared in this file. (final in
// Java.)
//
// Java: PotentialDeclaration.remove(AbstractCompiler)
func (b *PotentialDeclarationBase) Remove(compiler AbstractCompiler) {
	if b.self.IsDetached() {
		return
	}
	statement := b.self.GetRemovableNode()
	NodeUtil_DeleteNode(statement, compiler)
	statement.RemoveChildren()
}

// BreakDownDestructure breaks down this declaration if it's a destructuring
// LHS by replacing it with a VAR node and removing the RHS.
//
// TODO(lharker): can we merge this code into simplify(AbstractCompiler)?
//
// Returns true if the declaration is broken down. Otherwise, returns false.
//
// Java: PotentialDeclaration.breakDownDestructure(AbstractCompiler)
func (b *PotentialDeclarationBase) BreakDownDestructure(compiler AbstractCompiler) bool {
	return false
}

// PotentialDeclaration_NameDeclaration is a potential declaration that has a
// fully qualified name to describe it. This includes things like:
// var/let/const/function/class declarations, assignments to a fully
// qualified name, and goog.module exports This is the most common type of
// potential declaration.
type PotentialDeclaration_NameDeclaration struct {
	PotentialDeclarationBase
}

// Java: PotentialDeclaration.NameDeclaration(String, Node, Node)
func NewPotentialDeclaration_NameDeclaration(fullyQualifiedName string, lhs *rhino.Node, rhs *rhino.Node) *PotentialDeclaration_NameDeclaration {
	d := &PotentialDeclaration_NameDeclaration{}
	d.InitPotentialDeclarationBase(d, fullyQualifiedName, lhs, rhs)
	return d
}

// Java: PotentialDeclaration.NameDeclaration.simplifyNamespace(AbstractCompiler)
func (d *PotentialDeclaration_NameDeclaration) SimplifyNamespace(compiler AbstractCompiler) {
	if d.GetRhs().IsOr() {
		objLit := d.GetRhs().GetLastChild().Detach()
		d.GetRhs().ReplaceWith(objLit)
		compiler.ReportChangeToEnclosingScope(d.GetLhs())
	}
}

// Java: PotentialDeclaration.NameDeclaration.simplifySymbol(AbstractCompiler)
func (d *PotentialDeclaration_NameDeclaration) SimplifySymbol(compiler AbstractCompiler) {
	jutil.CheckArgument(NodeUtil_IsCallTo(d.GetRhs(), "Symbol"))
	callNode := d.GetRhs()
	for callNode.HasMoreThanOneChild() {
		NodeUtil_DeleteNode(callNode.GetLastChild(), compiler)
	}
}

// Java: PotentialDeclaration.NameDeclaration.simplify(AbstractCompiler)
func (d *PotentialDeclaration_NameDeclaration) Simplify(compiler AbstractCompiler) {
	if d.GetRhs() == nil || d.ShouldPreserve() {
		return
	}
	nameNode := d.GetLhs()
	jsdoc := d.GetJsDoc()
	if jsdoc != nil && jsdoc.HasEnumParameterType() {
		d.PotentialDeclarationBase.SimplifyEnumValues(compiler)
		return
	}
	if NodeUtil_IsNamespaceDecl(nameNode) {
		d.SimplifyNamespace(compiler)
		return
	}
	if nameNode.MatchesName("exports") {
		// Replace the RHS of a default goog.module export with Unknown
		PotentialDeclaration_NameDeclaration_ReplaceRhsWithUnknown(d.GetRhs())
		compiler.ReportChangeToEnclosingScope(nameNode)
		return
	}
	if NodeUtil_IsCallTo(d.GetRhs(), "Symbol") {
		d.SimplifySymbol(compiler)
		return
	}
	if d.GetLhs().GetParent().IsConst() {
		jsdoc = JsdocUtil_MarkConstant(jsdoc)
	}
	// Just completely remove the RHS, and replace with a getprop.
	newStatement :=
		NodeUtil_NewQNameDeclaration(compiler, *nameNode.GetQualifiedName(), nil, jsdoc)
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

// Java: PotentialDeclaration.NameDeclaration.breakDownDestructure(AbstractCompiler)
func (d *PotentialDeclaration_NameDeclaration) BreakDownDestructure(compiler AbstractCompiler) bool {
	if !NodeUtil_IsLhsByDestructuring(d.GetLhs()) {
		return false
	}

	rootTarget := NodeUtil_GetRootTarget(d.GetLhs())

	jutil.CheckState(rootTarget.GetParent().IsDestructuringLhs())
	if !NodeUtil_IsNameDeclaration(rootTarget.GetGrandparent()) {
		return false
	}

	definitionNode := rootTarget.GetGrandparent()

	var prev *rhino.Node

	var lhsNodes []*rhino.Node
	NodeUtil_VisitLhsNodesInNode(rootTarget.GetParent(), func(n *rhino.Node) { lhsNodes = append(lhsNodes, n) })
	for _, n := range lhsNodes {
		n.Detach()
		temp := rhino.IR_Var(n)
		if prev == nil {
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

// Java: PotentialDeclaration.NameDeclaration.replaceRhsWithUnknown(Node)
func PotentialDeclaration_NameDeclaration_ReplaceRhsWithUnknown(rhs *rhino.Node) {
	rhs.ReplaceWith(rhino.IR_Cast(rhino.IR_Number(0), JsdocUtil_GetQmarkTypeJSDoc()).SrcrefTree(rhs))
}

// Java: PotentialDeclaration.NameDeclaration.shouldPreserve()
func (d *PotentialDeclaration_NameDeclaration) ShouldPreserve() bool {
	rhs := d.GetRhs()
	nameNode := d.GetLhs()
	jsdoc := d.GetJsDoc()
	isExport := PotentialDeclaration_IsExportLhs(nameNode)
	return d.PotentialDeclarationBase.ShouldPreserve() ||
		PotentialDeclaration_IsImportRhs(rhs) ||
		(isExport && rhs != nil && (rhs.IsQualifiedName() || rhs.IsObjectLit())) ||
		(jsdoc != nil && jsdoc.IsConstructor() && rhs != nil && rhs.IsQualifiedName()) ||
		(rhs != nil &&
			rhs.IsObjectLit() &&
			!rhs.HasChildren() &&
			(jsdoc == nil || !JsdocUtil_HasAnnotatedType(jsdoc))) ||
		(rhs != nil && NodeUtil_IsCallTo(rhs, "Polymer")) ||
		d.IsPolymerBehaviorAliasOrArray()
}

// IsPolymerBehaviorAliasOrArray: Polymer Behaviors can take 3 forms:
//
//  1. /** @polymerBehavior *\/ export const MyBehavior = { ... };
//  2. /** @polymerBehavior *\/ export const MyBehaviorAlias = MyBehavior;
//  3. /** @polymerBehavior *\/ export const MyBehaviorArray = [Behavior1, Behavior2];
//
// Form #1 will be simplified by PolymerBehaviorPropertiesDeclaration. Forms
// #2 and #3 need to be preserved here as-is so that the PolymerPass can
// follow the name references. Other forms annotated with @polymerBehavior
// are invalid and can be simplified or removed like any other variable.
//
// Java: PotentialDeclaration.NameDeclaration.isPolymerBehaviorAliasOrArray()
func (d *PotentialDeclaration_NameDeclaration) IsPolymerBehaviorAliasOrArray() bool {
	jsdoc := d.GetJsDoc()
	rhs := d.GetRhs()
	return jsdoc != nil &&
		jsdoc.IsPolymerBehavior() &&
		rhs != nil &&
		(rhs.IsName() || rhs.IsArrayLit())
}

// PotentialDeclaration_ThisPropDeclaration is a declaration of a property on
// `this` inside a constructor.
type PotentialDeclaration_ThisPropDeclaration struct {
	PotentialDeclarationBase

	insertionPoint *rhino.Node
}

// Java: PotentialDeclaration.ThisPropDeclaration(String, Node, Node)
func NewPotentialDeclaration_ThisPropDeclaration(fullyQualifiedName string, lhs *rhino.Node, rhs *rhino.Node) *PotentialDeclaration_ThisPropDeclaration {
	d := &PotentialDeclaration_ThisPropDeclaration{}
	d.InitPotentialDeclarationBase(d, fullyQualifiedName, lhs, rhs)
	thisPropDefinition := NodeUtil_GetEnclosingStatement(lhs)
	d.insertionPoint = NodeUtil_GetEnclosingStatement(thisPropDefinition.GetParent())
	return d
}

// Java: PotentialDeclaration.ThisPropDeclaration.simplify(AbstractCompiler)
func (d *PotentialDeclaration_ThisPropDeclaration) Simplify(compiler AbstractCompiler) {
	if d.ShouldPreserve() {
		return
	}
	// Just completely remove the RHS, if present, and replace with a getprop.
	newStatement :=
		NodeUtil_NewQNameDeclaration(compiler, d.GetFullyQualifiedName(), nil, d.GetJsDoc())
	newStatement.SrcrefTreeIfMissing(d.GetLhs())
	NodeUtil_DeleteNode(d.GetRemovableNode(), compiler)
	if d.insertionPoint.HasParent() {
		newStatement.InsertAfter(d.insertionPoint)
		compiler.ReportChangeToEnclosingScope(newStatement)
	}
}

// PotentialDeclaration_DefineDeclaration is a declaration declared by a call
// to `goog.define`. Note that a let, const, or var declaration annotated
// with @define in its JSDoc and no 'goog.define' would be a NameDeclaration
// instead.
type PotentialDeclaration_DefineDeclaration struct {
	PotentialDeclarationBase
}

// Java: PotentialDeclaration.DefineDeclaration(String, Node, Node)
func NewPotentialDeclaration_DefineDeclaration(qualifiedName string, lhs *rhino.Node, rhs *rhino.Node) *PotentialDeclaration_DefineDeclaration {
	d := &PotentialDeclaration_DefineDeclaration{}
	d.InitPotentialDeclarationBase(d, qualifiedName, lhs, rhs)
	return d
}

// Java: PotentialDeclaration.DefineDeclaration.simplify(AbstractCompiler)
func (d *PotentialDeclaration_DefineDeclaration) Simplify(compiler AbstractCompiler) {
	info := d.GetJsDoc()
	if info != nil && info.GetType() != nil {
		newRhs := PotentialDeclaration_DefineDeclaration_MakeEmptyValueNode(info.GetType())
		if newRhs != nil {
			d.GetRhs().ReplaceWith(newRhs)
			compiler.ReportChangeToEnclosingScope(newRhs)
			return
		}
	}
	NodeUtil_DeleteNode(d.GetRemovableNode(), compiler)
}

// Java: PotentialDeclaration.DefineDeclaration.from(Node)
func PotentialDeclaration_DefineDeclaration_From(callNode *rhino.Node) *PotentialDeclaration_DefineDeclaration {
	// Match a few different forms, depending on the call node's parent:
	//   1. EXPR_RESULT: goog.define('foo', 1);
	//   2. ASSIGN: a.b = goog.define('c', 2);
	//   3. NAME: var x = goog.define('d', 3);
	switch callNode.GetParent().GetToken() {
	case rhino.Token_EXPR_RESULT:
		return NewPotentialDeclaration_DefineDeclaration(
			callNode.GetSecondChild().GetString(), callNode, callNode.GetLastChild())
	case rhino.Token_ASSIGN:
		previous := callNode.GetPrevious()
		return NewPotentialDeclaration_DefineDeclaration(
			*previous.GetQualifiedName(), previous, callNode.GetLastChild())
	case rhino.Token_NAME:
		parent := callNode.GetParent()
		return NewPotentialDeclaration_DefineDeclaration(parent.GetString(), parent, callNode.GetLastChild())
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected parent: " + callNode.GetParent().GetToken().String()))
	}
}

// Java: PotentialDeclaration.DefineDeclaration.makeEmptyValueNode(JSTypeExpression)
func PotentialDeclaration_DefineDeclaration_MakeEmptyValueNode(type_ *rhino.JSTypeExpression) *rhino.Node {
	n := type_.GetRoot()
	for n != nil && !n.IsStringLit() && !n.IsName() {
		n = n.GetFirstChild()
	}
	s := ""
	if n != nil {
		s = n.GetString()
	}
	switch s {
	case "boolean":
		return rhino.NewNode(rhino.Token_FALSE)
	case "number":
		return rhino.Node_NewNumber(0)
	case "string":
		return rhino.Node_NewString("")
	default:
		return nil
	}
}

// PotentialDeclaration_MethodDeclaration is a declaration of a method
// defined using the ES6 method syntax. Note that a method defined as an
// assignment to a prototype property would be a NameDeclaration instead.
type PotentialDeclaration_MethodDeclaration struct {
	PotentialDeclarationBase
}

// Java: PotentialDeclaration.MethodDeclaration(String, Node)
func NewPotentialDeclaration_MethodDeclaration(name string, functionNode *rhino.Node) *PotentialDeclaration_MethodDeclaration {
	d := &PotentialDeclaration_MethodDeclaration{}
	d.InitPotentialDeclarationBase(d, name, functionNode.GetParent(), functionNode)
	return d
}

// Java: PotentialDeclaration.MethodDeclaration.simplify(AbstractCompiler)
func (d *PotentialDeclaration_MethodDeclaration) Simplify(compiler AbstractCompiler) {}

// Java: PotentialDeclaration.MethodDeclaration.getRemovableNode()
func (d *PotentialDeclaration_MethodDeclaration) GetRemovableNode() *rhino.Node {
	return d.GetLhs()
}

// PotentialDeclaration_MemberFieldDeclaration is the private static nested
// class PotentialDeclaration.MemberFieldDeclaration.
type PotentialDeclaration_MemberFieldDeclaration struct {
	PotentialDeclarationBase
}

// Java: PotentialDeclaration.MemberFieldDeclaration(String, Node, Node)
func NewPotentialDeclaration_MemberFieldDeclaration(name string, lhs *rhino.Node, rhs *rhino.Node) *PotentialDeclaration_MemberFieldDeclaration {
	d := &PotentialDeclaration_MemberFieldDeclaration{}
	d.InitPotentialDeclarationBase(d, name, lhs, rhs)
	return d
}

// Java: PotentialDeclaration.MemberFieldDeclaration.simplify(AbstractCompiler)
func (d *PotentialDeclaration_MemberFieldDeclaration) Simplify(compiler AbstractCompiler) {
	rhs := d.GetRhs()
	if rhs != nil {
		NodeUtil_DeleteNode(rhs, compiler)
	}
}

// Java: PotentialDeclaration.MemberFieldDeclaration.getRemovableNode()
func (d *PotentialDeclaration_MemberFieldDeclaration) GetRemovableNode() *rhino.Node {
	return d.GetLhs()
}

// PotentialDeclaration_StringKeyDeclaration is the private static nested
// class PotentialDeclaration.StringKeyDeclaration.
type PotentialDeclaration_StringKeyDeclaration struct {
	PotentialDeclarationBase
}

// Java: PotentialDeclaration.StringKeyDeclaration(String, Node)
func NewPotentialDeclaration_StringKeyDeclaration(name string, stringKeyNode *rhino.Node) *PotentialDeclaration_StringKeyDeclaration {
	d := &PotentialDeclaration_StringKeyDeclaration{}
	d.InitPotentialDeclarationBase(d, name, stringKeyNode, stringKeyNode.GetLastChild())
	return d
}

// Java: PotentialDeclaration.StringKeyDeclaration.simplify(AbstractCompiler)
func (d *PotentialDeclaration_StringKeyDeclaration) Simplify(compiler AbstractCompiler) {
	if d.ShouldPreserve() {
		return
	}
	jsdoc := d.GetJsDoc()
	if jsdoc != nil && jsdoc.HasEnumParameterType() {
		d.PotentialDeclarationBase.SimplifyEnumValues(compiler)
		return
	}
	key := d.GetLhs()
	PotentialDeclaration_RemoveStringKeyValue(key)
	compiler.ReportChangeToEnclosingScope(key)
	if jsdoc == nil || !jsdoc.ContainsDeclaration() || d.IsConstToBeInferred() {
		key.SetJSDocInfo(JsdocUtil_GetUnusableTypeJSDoc(jsdoc))
	}
}

// Java: PotentialDeclaration.StringKeyDeclaration.shouldPreserve()
func (d *PotentialDeclaration_StringKeyDeclaration) ShouldPreserve() bool {
	return d.PotentialDeclarationBase.IsDetached() || d.PotentialDeclarationBase.ShouldPreserve() || !d.IsInNamespace()
}

// Java: PotentialDeclaration.StringKeyDeclaration.isInNamespace()
func (d *PotentialDeclaration_StringKeyDeclaration) IsInNamespace() bool {
	stringKey := d.GetLhs()
	objLit := stringKey.GetParent()
	lvalue := NodeUtil_GetBestLValue(objLit)
	if lvalue == nil {
		return false
	}
	jsdoc := NodeUtil_GetBestJSDocInfo(lvalue)
	return !PotentialDeclaration_IsExportLhs(lvalue) &&
		!JsdocUtil_HasAnnotatedType(jsdoc) &&
		NodeUtil_IsNamespaceDecl(lvalue)
}

// Java: PotentialDeclaration.StringKeyDeclaration.getRemovableNode()
func (d *PotentialDeclaration_StringKeyDeclaration) GetRemovableNode() *rhino.Node {
	return d.GetLhs()
}

// PotentialDeclaration_PolymerBehaviorPropertiesDeclaration: Polymer
// Behaviors are mixin-like objects used in Polymer 1 for multiple
// inheritance. They are also supported by Polymer 2 and 3 for
// backwards-compatibility, though their use is discouraged in favor of
// regular JavaScript mixins.
//
// Example:
//
//	/** @polymerBehavior *\/
//	export const MyBehavior = {
//	  properties: {
//	    foo: String,
//	    bar: {
//	      type: Number,
//	      value: 123,
//	      readOnly: true
//	    }
//	  },
//	  baz: function() {}
//	};
//
// For incremental compilation, it is important that the "properties" object
// is preserved, because the PolymerPass injects the properties declared
// there onto the prototypes of the Polymer elements that apply that
// behavior. Note that method signatures are already preserved so don't need
// additional handling here.
type PotentialDeclaration_PolymerBehaviorPropertiesDeclaration struct {
	PotentialDeclarationBase
}

// Java: PotentialDeclaration.PolymerBehaviorPropertiesDeclaration(String, Node)
func NewPotentialDeclaration_PolymerBehaviorPropertiesDeclaration(name string, stringKeyNode *rhino.Node) *PotentialDeclaration_PolymerBehaviorPropertiesDeclaration {
	d := &PotentialDeclaration_PolymerBehaviorPropertiesDeclaration{}
	d.InitPotentialDeclarationBase(d, name, stringKeyNode, stringKeyNode.GetLastChild())
	return d
}

// Java: PotentialDeclaration.PolymerBehaviorPropertiesDeclaration.simplify(AbstractCompiler)
func (d *PotentialDeclaration_PolymerBehaviorPropertiesDeclaration) Simplify(compiler AbstractCompiler) {
	if d.IsDetached() {
		return
	}
	propertiesObject := d.GetRhs()
	if !propertiesObject.IsObjectLit() || !propertiesObject.HasChildren() {
		return
	}
	for propKey := propertiesObject.GetFirstChild(); propKey != nil; propKey = propKey.GetNext() {
		propDef := propKey.GetOnlyChild()
		// A property definition is either a function reference (e.g. String, Number), or another
		// object literal. If it's an object literal, only the "type" and "readOnly" sub-properties
		// matters for type checking, so we can delete everything else (which may include e.g. a
		// "value" sub-property with a function expression).
		if propDef.IsObjectLit() {
			for subProp := propDef.GetFirstChild(); subProp != nil; {
				next := subProp.GetNext()
				if subProp.GetString() != "type" && subProp.GetString() != "readOnly" {
					NodeUtil_DeleteNode(subProp, compiler)
				}
				subProp = next
			}
		}
	}
}

// Java: PotentialDeclaration.PolymerBehaviorPropertiesDeclaration.shouldPreserve()
func (d *PotentialDeclaration_PolymerBehaviorPropertiesDeclaration) ShouldPreserve() bool {
	return true
}

// Java: PotentialDeclaration.PolymerBehaviorPropertiesDeclaration.getRemovableNode()
func (d *PotentialDeclaration_PolymerBehaviorPropertiesDeclaration) GetRemovableNode() *rhino.Node {
	return d.GetLhs()
}

// PotentialDeclaration_AliasDeclaration is the private static nested class
// PotentialDeclaration.AliasDeclaration.
type PotentialDeclaration_AliasDeclaration struct {
	PotentialDeclarationBase
}

// NewPotentialDeclaration_AliasDeclaration: name is the alias name being
// declared; lhs the NAME node that represents the name of the individual
// alias.
//
// Java: PotentialDeclaration.AliasDeclaration(String, Node)
func NewPotentialDeclaration_AliasDeclaration(name string, lhs *rhino.Node) *PotentialDeclaration_AliasDeclaration {
	d := &PotentialDeclaration_AliasDeclaration{}
	d.InitPotentialDeclarationBase(d, name, lhs, nil)
	return d
}

// Java: PotentialDeclaration.AliasDeclaration.simplify(AbstractCompiler)
func (d *PotentialDeclaration_AliasDeclaration) Simplify(compiler AbstractCompiler) {
	// Does not simplify
}

// GetRemovableNode: if the declaration is a destructuring declaration: 1) If
// the lhs's destructuring pattern parent has only one child, e.g. const {Foo}
// = x; returns the enclosing statement to remove the entire statement. 2) If
// the parent has more than one children, e.g. const {Foo, Bar} = x; returns
// the lhs so that when Foo is removed, const {Foo, Bar} = x; becomes const
// {Bar} = x; Otherwise, returns the enclosing statement.
//
// Java: PotentialDeclaration.AliasDeclaration.getRemovableNode()
func (d *PotentialDeclaration_AliasDeclaration) GetRemovableNode() *rhino.Node {
	lhs := d.GetLhs()
	if lhs.GetParent().IsArrayPattern() && lhs.GetParent().HasMoreThanOneChild() {
		return lhs
	}
	if lhs.GetGrandparent().IsObjectPattern() && lhs.GetGrandparent().HasMoreThanOneChild() {
		return lhs.GetParent()
	}
	return NodeUtil_GetEnclosingStatement(lhs)
}

// Java: PotentialDeclaration.AliasDeclaration.isDefiniteDeclaration(AbstractCompiler)
func (d *PotentialDeclaration_AliasDeclaration) IsDefiniteDeclaration(compiler AbstractCompiler) bool {
	return true
}

// Java: PotentialDeclaration.AliasDeclaration.shouldPreserve()
func (d *PotentialDeclaration_AliasDeclaration) ShouldPreserve() bool {
	return true
}

// SimplifyEnumValues removes values from enums
//
// Java: PotentialDeclaration.simplifyEnumValues(AbstractCompiler)
func (b *PotentialDeclarationBase) SimplifyEnumValues(compiler AbstractCompiler) {
	if b.self.GetRhs().IsObjectLit() && b.self.GetRhs().HasChildren() {
		changed := false
		for key := b.self.GetRhs().GetFirstChild(); key != nil; key = key.GetNext() {
			value := key.GetOnlyChild()
			if !value.IsStringLit() {
				PotentialDeclaration_RemoveStringKeyValue(key)
				changed = true
			} else {
				content := value.GetString()
				if jutil.StringLength(content) > potentialDeclaration_STRING_ENUM_RETAIN_CAP {
					PotentialDeclaration_TruncateStringKeyValue(key)
					changed = true
				}
			}
		}
		if changed {
			compiler.ReportChangeToEnclosingScope(b.self.GetRhs())
		}
	}
}

// Java: PotentialDeclaration.isDefiniteDeclaration(AbstractCompiler)
func (b *PotentialDeclarationBase) IsDefiniteDeclaration(compiler AbstractCompiler) bool {
	parent := b.self.GetLhs().GetParent()
	switch parent.GetToken() {
	case rhino.Token_DEFAULT_VALUE:
		if !parent.GetParent().IsStringKey() {
			return false
		}
		// fall through
		fallthrough
	case rhino.Token_COMPUTED_PROP, rhino.Token_STRING_KEY:
		if NodeUtil_IsLhsByDestructuring(b.self.GetLhs()) {
			rootTarget := NodeUtil_GetRootTarget(b.self.GetLhs())
			jutil.CheckState(rootTarget.GetParent().IsDestructuringLhs())
			if NodeUtil_IsNameDeclaration(rootTarget.GetGrandparent()) {
				return true
			}
		}
		return false
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST, rhino.Token_CLASS, rhino.Token_FUNCTION:
		return true
	default:
		return PotentialDeclaration_IsExportLhs(b.self.GetLhs()) ||
			(b.self.GetJsDoc() != nil && b.self.GetJsDoc().ContainsDeclaration()) ||
			(b.self.GetRhs() != nil && PotentialDeclaration_IsTypedRhs(b.self.GetRhs()))
	}
}

// Java: PotentialDeclaration.shouldPreserve()
func (b *PotentialDeclarationBase) ShouldPreserve() bool {
	return b.self.GetRhs() != nil && PotentialDeclaration_IsTypedRhs(b.self.GetRhs())
}

// Java: PotentialDeclaration.isConstToBeInferred()
func (b *PotentialDeclarationBase) IsConstToBeInferred() bool {
	return PotentialDeclaration_IsConstToBeInferred(b.self.GetLhs())
}

// Java: PotentialDeclaration.isConstToBeInferred(Node)
func PotentialDeclaration_IsConstToBeInferred(nameNode *rhino.Node) bool {
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
func PotentialDeclaration_IsTypedRhs(rhs *rhino.Node) bool {
	return rhs.IsFunction() ||
		rhs.IsClass() ||
		(rhs.IsQualifiedName() && potentialDeclaration_GOOG_ABSTRACTMETHOD.Matches(rhs))
}

// Java: PotentialDeclaration.isExportLhs(Node)
func PotentialDeclaration_IsExportLhs(lhs *rhino.Node) bool {
	return (lhs.IsName() && lhs.MatchesName("exports")) ||
		(lhs.IsGetProp() && lhs.GetFirstChild().MatchesName("exports")) ||
		potentialDeclaration_MODULE_EXPORTS.Matches(lhs)
}

// Java: PotentialDeclaration.isImportRhs(Node)
func PotentialDeclaration_IsImportRhs(rhs *rhino.Node) bool {
	if rhs == nil || !rhs.IsCall() {
		return false
	}
	callee := rhs.GetFirstChild()
	return potentialDeclaration_GOOG_REQUIRE.Matches(callee) ||
		potentialDeclaration_GOOG_REQUIRETYPE.Matches(callee) ||
		potentialDeclaration_GOOG_FORWARDDECLARE.Matches(callee) ||
		callee.MatchesName("require")
}

// Java: PotentialDeclaration.isAliasDeclaration(Node, Node)
func PotentialDeclaration_IsAliasDeclaration(lhs *rhino.Node, rhs *rhino.Node) bool {
	return !ClassUtil_IsThisPropInsideClassWithName(lhs) &&
		PotentialDeclaration_IsConstToBeInferred(lhs) &&
		rhs != nil &&
		PotentialDeclaration_IsQualifiedAliasExpression(rhs)
}

var potentialDeclaration_GOOG_MODULE_GET = rhino.QualifiedName_Of("goog.module.get")

// PotentialDeclaration_IsQualifiedAliasExpression returns whether a node
// corresponds to a simple or a qualified name, such as x or a.b.c or this.a.
//
// Java: PotentialDeclaration.isQualifiedAliasExpression(Node)
func PotentialDeclaration_IsQualifiedAliasExpression(n *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_NAME:
		return n.GetString() != ""
	case rhino.Token_THIS, rhino.Token_SUPER:
		return true
	case rhino.Token_GETPROP:
		return PotentialDeclaration_IsQualifiedAliasExpression(n.GetFirstChild())
	case rhino.Token_CALL:
		if potentialDeclaration_GOOG_MODULE_GET.Matches(n.GetFirstChild()) {
			return true
		}
		return false

	case rhino.Token_MEMBER_FUNCTION_DEF:
		// These are explicitly *not* qualified name components.
		return false
	default:
		return false
	}
}

// Java: PotentialDeclaration.removeStringKeyValue(Node)
func PotentialDeclaration_RemoveStringKeyValue(stringKey *rhino.Node) {
	value := stringKey.GetOnlyChild()
	replacementValue := rhino.IR_Number(0).SrcrefTree(value)
	value.ReplaceWith(replacementValue)
}

// Java: PotentialDeclaration.truncateStringKeyValue(Node)
func PotentialDeclaration_TruncateStringKeyValue(stringKey *rhino.Node) {
	value := stringKey.GetOnlyChild()
	replacementValue :=
		rhino.IR_String(jutil.StringSubstring(value.GetString(), 0, potentialDeclaration_STRING_ENUM_RETAIN_CAP-2) + "..").
			SrcrefTree(value)
	value.ReplaceWith(replacementValue)
}
