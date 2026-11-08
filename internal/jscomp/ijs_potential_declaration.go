// Port of com/google/javascript/jscomp/ijs/PotentialDeclaration.java @ closure-compiler v20261104

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
	"github.com/paullepoulpe/jsgo/internal/strpool"
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

// potentialDeclaration encapsulates a candidate declaration in `.i.js` summary
// generation: var/let/const declarations, function/class declarations, method
// and class-field declarations, assignments, goog.define calls, and annotated
// valueless property accesses (e.g. `/** @type {number} */ Foo.prototype.bar`).
type potentialDeclaration struct {
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
const stringEnumRetainCap = 10

func mustPotentialDeclQualifiedName(n rhino.Node) string {
	s, ok := n.QualifiedName()
	if !ok {
		panic(jutil.NewNullPointerException())
	}
	return s
}

func newPotentialDeclaration(kind potentialDeclKind, fullyQualifiedName string, lhs, rhs rhino.Node) *potentialDeclaration {
	if lhs.IsNil() {
		panic(jutil.NewNullPointerException())
	}
	return &potentialDeclaration{
		fullyQualifiedName: fullyQualifiedName,
		lhs:                lhs,
		rhs:                rhs,
		kind:               kind,
	}
}

func potentialDeclFromName(nameNode rhino.Node) *potentialDeclaration {
	if nameNode.IsNil() || !nameNode.IsQualifiedName() {
		panic(jutil.NewIllegalArgumentExceptionMsg("expected qualified name"))
	}
	rhs := NodeUtil_GetRValueOfLValue(nameNode)
	if className, ok := classNameOfThisProp(nameNode); ok {
		name := className + ".prototype." + ast.MustGetProp(nameNode).Property()
		d := newPotentialDeclaration(potentialDeclThisProp, name, nameNode, rhs)
		if stmt := NodeUtil_GetEnclosingStatement(nameNode); !stmt.IsNil() && !stmt.GetParent().IsNil() {
			d.insertionPoint = NodeUtil_GetEnclosingStatement(stmt.GetParent())
		}
		return d
	}
	return newPotentialDeclaration(potentialDeclName, mustPotentialDeclQualifiedName(nameNode), nameNode, rhs)
}

func potentialDeclFromMemberField(field ast.MemberFieldDef) *potentialDeclaration {
	val, _ := field.Value()
	// lhs is MEMBER_FIELD_DEF; rhs is the value assigned to the field if any.
	return newPotentialDeclaration(
		potentialDeclMemberField,
		memberFieldQualifiedName(field),
		field.Node(),
		val.Node())
}

func potentialDeclFromMethod(method ast.MemberFunctionDef) *potentialDeclaration {
	return newPotentialDeclaration(
		potentialDeclMethod,
		classMethodQualifiedName(method),
		method.Node(),
		method.Function().Node())
}

func potentialDeclFromStringKey(sk ast.StringKey) *potentialDeclaration {
	stringKeyNode := sk.Node()
	parent := stringKeyNode.GetParent()
	if parent.IsNil() || !parent.IsObjectLit() {
		panic(jutil.NewIllegalArgumentExceptionMsg("expected STRING_KEY inside OBJECTLIT"))
	}
	keyName := sk.Key()
	name := "this." + keyName
	kind := potentialDeclStringKey
	if keyName == "properties" {
		if objLitJsDoc := NodeUtil_GetBestJSDocInfo(parent); objLitJsDoc != nil && objLitJsDoc.IsPolymerBehavior() {
			kind = potentialDeclPolymerBehaviorProperties
		}
	}
	return newPotentialDeclaration(kind, name, stringKeyNode, sk.Value().Node())
}

func potentialDeclFromDefine(call ast.Call) *potentialDeclaration {
	callNode := call.Node()
	if !NodeUtil_IsCallTo(callNode, "goog.define") {
		panic(jutil.NewIllegalArgumentExceptionMsg("expected goog.define call"))
	}
	var rhs rhino.Node
	if call.NumArgs() >= 2 {
		if lastArg, ok := call.LastArg(); ok {
			rhs = lastArg.Node()
		}
	}
	// Match a few different forms, depending on the call node's parent:
	//   1. EXPR_RESULT: goog.define('foo', 1);
	//   2. ASSIGN: a.b = goog.define('c', 2);
	//   3. NAME: var x = goog.define('d', 3);
	parent := callNode.GetParent()
	if parent.IsNil() {
		return nil
	}
	switch parent.GetToken() {
	case rhino.Token_EXPR_RESULT:
		firstArg, ok := call.FirstArg()
		str, isStr := firstArg.AsStringLit()
		if !ok || !isStr {
			return nil
		}
		return newPotentialDeclaration(potentialDeclDefine, str.Value(), callNode, rhs)
	case rhino.Token_ASSIGN:
		previous := ast.MustAssign(parent).Target().Node()
		qname, ok := previous.QualifiedName()
		if !ok {
			return nil
		}
		return newPotentialDeclaration(potentialDeclDefine, qname, previous, rhs)
	case rhino.Token_NAME:
		return newPotentialDeclaration(potentialDeclDefine, ast.MustName(parent).Ident(), parent, rhs)
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected parent: " + parent.GetToken().String()))
	}
}

func potentialDeclFromAlias(nameNode rhino.Node) *potentialDeclaration {
	if nameNode.IsNil() || !nameNode.IsQualifiedName() {
		panic(jutil.NewIllegalArgumentExceptionMsg("expected qualified alias name"))
	}
	return newPotentialDeclaration(potentialDeclAlias, mustPotentialDeclQualifiedName(nameNode), nameNode, rhino.NilNode)
}

func (d *potentialDeclaration) jsdoc() *rhino.JSDocInfo {
	return NodeUtil_GetBestJSDocInfo(d.lhs)
}

func (d *potentialDeclaration) isDetached() bool {
	for current := d.lhs; !current.IsNil(); current = current.GetParent() {
		if current.IsScript() {
			return false
		}
	}
	return true
}

func countNonEmptyArrayPatternElems(ap ast.ArrayPattern) int {
	n := 0
	for elem := range ap.Elements() {
		if !elem.Node().IsEmpty() {
			n++
		}
	}
	return n
}

func (d *potentialDeclaration) removableNode() rhino.Node {
	switch d.kind {
	case potentialDeclMethod,
		potentialDeclMemberField,
		potentialDeclStringKey,
		potentialDeclPolymerBehaviorProperties:
		return d.lhs
	case potentialDeclAlias:
		// If the declaration is a destructuring declaration:
		// 1) If the lhs's destructuring pattern parent has only one bound entry,
		//    e.g. `const {Foo} = x;`, return the enclosing statement to remove it all.
		// 2) If the pattern has more than one bound entry, e.g. `const {Foo, Bar} = x;`
		//    or `const [Foo, Bar] = x;`, return only the individual pattern entry so
		//    remaining sibling bindings stay intact.
		elem := d.lhs
		if parent := elem.GetParent(); !parent.IsNil() && (parent.IsDefaultValue() || parent.IsRest()) {
			elem = parent
		}
		parent := elem.GetParent()
		if ap, ok := ast.AsArrayPattern(parent); ok && countNonEmptyArrayPatternElems(ap) > 1 {
			return elem
		}
		if !parent.IsNil() && (parent.IsStringKey() || parent.IsComputedProp()) {
			if op, ok := ast.AsObjectPattern(parent.GetParent()); ok && op.NumProps() > 1 {
				return parent
			}
		}
		if op, ok := ast.AsObjectPattern(parent); ok && elem.IsObjectRest() && op.NumProps() > 1 {
			return elem
		}
		return NodeUtil_GetEnclosingStatement(d.lhs)
	default:
		return NodeUtil_GetEnclosingStatement(d.lhs)
	}
}

// remove deletes this potential declaration completely (typically when a
// higher-precedence declaration for the same symbol already exists in the file).
func (d *potentialDeclaration) remove(compiler AbstractCompiler) {
	if d.isDetached() {
		return
	}
	statement := d.removableNode()
	if statement.IsNil() {
		return
	}
	if parent := statement.GetParent(); !parent.IsNil() {
		if ap, ok := ast.AsArrayPattern(parent); ok {
			lastElem, _ := ap.LastElement()
			if !lastElem.Is(statement) {
				NodeUtil_MarkFunctionsDeleted(statement, compiler)
				statement.ReplaceWith(rhino.IR_Empty().Srcref(statement))
				compiler.ReportChangeToEnclosingScope(parent)
				statement.RemoveChildren()
				return
			}
		}
	}
	NodeUtil_DeleteNode(statement, compiler)
	statement.RemoveChildren()
}

// breakDownDestructure breaks down this declaration if it is a destructuring
// LHS by replacing it with separate VAR declarations and removing the RHS.
func (d *potentialDeclaration) breakDownDestructure(compiler AbstractCompiler) bool {
	if d.kind != potentialDeclName || !NodeUtil_IsLhsByDestructuring(d.lhs) {
		return false
	}

	rootTarget := NodeUtil_GetRootTarget(d.lhs)
	if rootTarget.IsNil() {
		return false
	}
	dlhs, ok := ast.AsDestructuringLhs(rootTarget.GetParent())
	if !ok {
		return false
	}
	definitionNode := dlhs.Node().GetParent()
	if !NodeUtil_IsNameDeclaration(definitionNode) {
		return false
	}

	isExport := ast.RoleOf(definitionNode) == ast.RoleExportBody
	statement := definitionNode
	if isExport {
		statement = definitionNode.GetParent()
	}

	var prev rhino.Node
	var lhsNodes []rhino.Node
	NodeUtil_VisitLhsNodesInNode(dlhs.Node(), func(n rhino.Node) { lhsNodes = append(lhsNodes, n) })
	if len(lhsNodes) == 0 {
		NodeUtil_DeleteNode(statement, compiler)
		return false
	}
	for _, n := range lhsNodes {
		n.Detach()
		stmt := rhino.IR_Var(n).Srcref(definitionNode)
		if isExport {
			stmt = ast.NewExport(ast.MustExportBody(stmt)).Node().Srcref(statement)
		}
		if prev.IsNil() {
			statement.ReplaceWith(stmt)
			compiler.ReportChangeToEnclosingScope(stmt)
		} else {
			stmt.InsertAfter(prev)
			compiler.ReportChangeToEnclosingScope(stmt)
			stmt.SrcrefTreeIfMissing(prev)
		}
		prev = stmt
	}
	return true
}

func (d *potentialDeclaration) simplify(compiler AbstractCompiler) {
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

func (d *potentialDeclaration) simplifyName(compiler AbstractCompiler) {
	if d.rhs.IsNil() || d.shouldPreserve() {
		return
	}
	nameNode := d.lhs
	jsdoc := d.jsdoc()
	if jsdoc != nil && jsdoc.HasEnumParameterType() {
		d.simplifyEnumValues(compiler)
		return
	}
	if NodeUtil_IsNamespaceDecl(nameNode) {
		if bin, ok := ast.AsBinary(d.rhs); ok && bin.Token() == rhino.Token_OR {
			objLit := bin.Right().Detach().Node()
			d.rhs.ReplaceWith(objLit)
			compiler.ReportChangeToEnclosingScope(d.lhs)
		}
		return
	}
	if nameNode.MatchesName("exports") {
		// Replace the RHS of a default goog.module export with Unknown.
		d.rhs.ReplaceWith(rhino.IR_Cast(rhino.IR_Number(0), qmarkTypeJSDoc()).SrcrefTree(d.rhs))
		compiler.ReportChangeToEnclosingScope(nameNode)
		return
	}
	if NodeUtil_IsCallTo(d.rhs, "Symbol") {
		call := ast.MustCall(d.rhs)
		for {
			lastArg, ok := call.LastArg()
			if !ok {
				break
			}
			NodeUtil_DeleteNode(lastArg.Node(), compiler)
		}
		return
	}
	if parent := d.lhs.GetParent(); !parent.IsNil() && parent.IsConst() {
		jsdoc = markJSDocConstant(jsdoc)
	}
	// Completely remove the RHS, and replace with a getprop.
	newStatement := NodeUtil_NewQNameDeclaration(compiler, mustPotentialDeclQualifiedName(nameNode), rhino.NilNode, jsdoc)
	newStatement.SrcrefTreeIfMissing(nameNode)
	oldStatement := d.removableNode()
	NodeUtil_DeleteChildren(oldStatement, compiler)
	if oldStatement.IsExport() {
		oldStatement.AddChildToBack(newStatement)
	} else {
		oldStatement.ReplaceWith(newStatement)
	}
	compiler.ReportChangeToEnclosingScope(newStatement)
}

func (d *potentialDeclaration) simplifyThisProp(compiler AbstractCompiler) {
	if d.shouldPreserve() {
		return
	}
	// Completely remove the RHS, if present, and replace with a getprop.
	newStatement := NodeUtil_NewQNameDeclaration(compiler, d.fullyQualifiedName, rhino.NilNode, d.jsdoc())
	newStatement.SrcrefTreeIfMissing(d.lhs)
	if rem := d.removableNode(); !rem.IsNil() {
		NodeUtil_DeleteNode(rem, compiler)
	}
	if !d.insertionPoint.IsNil() && d.insertionPoint.HasParent() {
		newStatement.InsertAfter(d.insertionPoint)
		compiler.ReportChangeToEnclosingScope(newStatement)
	}
}

func (d *potentialDeclaration) simplifyDefine(compiler AbstractCompiler) {
	info := d.jsdoc()
	if info != nil && info.GetType() != nil && !d.rhs.IsNil() {
		if newRhs := makeEmptyDefineValueNode(info.GetType()); !newRhs.IsNil() {
			d.rhs.ReplaceWith(newRhs)
			compiler.ReportChangeToEnclosingScope(newRhs)
			return
		}
	}
	if rem := d.removableNode(); !rem.IsNil() {
		NodeUtil_DeleteNode(rem, compiler)
	}
}

var (
	potentialDeclBooleanID  = strpool.Intern("boolean")
	potentialDeclNumberID   = strpool.Intern("number")
	potentialDeclStringID   = strpool.Intern("string")
	potentialDeclTypeID     = strpool.Intern("type")
	potentialDeclReadOnlyID = strpool.Intern("readOnly")
)

func makeEmptyDefineValueNode(typeExpr *rhino.JSTypeExpression) rhino.Node {
	if typeExpr == nil {
		return rhino.NilNode
	}
	// Walk down JSDoc type-expression wrappers (e.g. BANG/QMARK/PIPE), which are
	// JSDoc AST nodes rather than JS syntax nodes.
	n := typeExpr.GetRoot()
	for !n.IsNil() && !n.IsStringLit() && !n.IsName() {
		next := rhino.NilNode
		for child := range n.Children() {
			next = child
			break
		}
		n = next
	}
	if n.IsNil() {
		return rhino.NilNode
	}
	switch nodeStringID(n) {
	case potentialDeclBooleanID:
		return rhino.NewNode(rhino.Token_FALSE)
	case potentialDeclNumberID:
		return rhino.Node_NewNumber(0)
	case potentialDeclStringID:
		return rhino.Node_NewString("")
	default:
		return rhino.NilNode
	}
}

func (d *potentialDeclaration) simplifyStringKey(compiler AbstractCompiler) {
	if d.shouldPreserve() {
		return
	}
	jsdoc := d.jsdoc()
	if jsdoc != nil && jsdoc.HasEnumParameterType() {
		d.simplifyEnumValues(compiler)
		return
	}
	sk, ok := ast.AsStringKey(d.lhs)
	if !ok {
		return
	}
	replaceStringKeyValueWithZero(sk)
	compiler.ReportChangeToEnclosingScope(d.lhs)
	if jsdoc == nil || !jsdoc.ContainsDeclaration() || d.isConstToBeInferred() {
		d.lhs.SetJSDocInfo(unusableTypeJSDoc(jsdoc))
	}
}

// simplifyPolymerBehaviorProperties preserves the "properties" object on a
// @polymerBehavior literal while stripping non-"type"/non-"readOnly" keys.
func (d *potentialDeclaration) simplifyPolymerBehaviorProperties(compiler AbstractCompiler) {
	if d.isDetached() {
		return
	}
	propertiesObject, ok := ast.AsObjectLit(d.rhs)
	if !ok || !propertiesObject.HasMembers() {
		return
	}
	for prop := range propertiesObject.Members() {
		sk, isKey := prop.AsStringKey()
		if !isKey {
			continue
		}
		// A property definition is either a function reference (e.g. String, Number), or another
		// object literal. If it's an object literal, only the "type" and "readOnly" sub-properties
		// matter for type checking, so we can delete everything else.
		if defObj, isObj := sk.Value().AsObjectLit(); isObj {
			for subProp := range ast.ReadAhead(defObj.FirstMember()) {
				if subKey, isSubKey := subProp.AsStringKey(); isSubKey {
					subID := nodeStringID(subKey.Node())
					if subID == potentialDeclTypeID || subID == potentialDeclReadOnlyID {
						continue
					}
				}
				NodeUtil_DeleteNode(subProp.Node(), compiler)
			}
		}
	}
}

// simplifyEnumValues removes non-string or oversized string values from @enum literals.
func (d *potentialDeclaration) simplifyEnumValues(compiler AbstractCompiler) {
	objLit, ok := ast.AsObjectLit(d.rhs)
	if !ok || !objLit.HasMembers() {
		return
	}
	changed := false
	for member := range ast.ReadAhead(objLit.FirstMember()) {
		if sk, isKey := member.AsStringKey(); isKey {
			if str, isStr := sk.Value().AsStringLit(); !isStr {
				replaceStringKeyValueWithZero(sk)
				changed = true
			} else if jutil.StringLength(str.Value()) > stringEnumRetainCap {
				truncateStringKeyValue(sk, str.Value())
				changed = true
			}
		} else if _, isSpread := member.AsObjectSpread(); isSpread {
			NodeUtil_DeleteNode(member.Node(), compiler)
			changed = true
		}
	}
	if changed {
		compiler.ReportChangeToEnclosingScope(d.rhs)
	}
}

func (d *potentialDeclaration) isDefiniteDeclaration() bool {
	if d.kind == potentialDeclAlias {
		return true
	}
	if NodeUtil_IsLhsByDestructuring(d.lhs) {
		rootTarget := NodeUtil_GetRootTarget(d.lhs)
		return !rootTarget.IsNil() &&
			rootTarget.GetParent().IsDestructuringLhs() &&
			NodeUtil_IsNameDeclaration(rootTarget.GetGrandparent())
	}
	parent := d.lhs.GetParent()
	if parent.IsNil() {
		return false
	}
	switch parent.GetToken() {
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST, rhino.Token_CLASS, rhino.Token_FUNCTION:
		return true
	default:
		jsdoc := d.jsdoc()
		return isExportLhs(d.lhs) ||
			(jsdoc != nil && jsdoc.ContainsDeclaration()) ||
			isTypedIjsRhs(d.rhs)
	}
}

func (d *potentialDeclaration) shouldPreserveBase() bool {
	return isTypedIjsRhs(d.rhs)
}

func (d *potentialDeclaration) shouldPreserve() bool {
	switch d.kind {
	case potentialDeclPolymerBehaviorProperties, potentialDeclAlias:
		return true
	case potentialDeclStringKey:
		return d.isDetached() || d.shouldPreserveBase() || !d.isStringKeyInNamespace()
	case potentialDeclName:
		rhs := d.rhs
		nameNode := d.lhs
		jsdoc := d.jsdoc()
		isExport := isExportLhs(nameNode)
		return d.shouldPreserveBase() ||
			isImportRhs(rhs) ||
			(isExport && !rhs.IsNil() && (rhs.IsQualifiedName() || rhs.IsObjectLit())) ||
			(jsdoc != nil && jsdoc.IsConstructor() && !rhs.IsNil() && rhs.IsQualifiedName()) ||
			(!rhs.IsNil() &&
				rhs.IsObjectLit() &&
				!rhs.HasChildren() &&
				(jsdoc == nil || !hasAnnotatedType(jsdoc))) ||
			(!rhs.IsNil() && NodeUtil_IsCallTo(rhs, "Polymer")) ||
			d.isPolymerBehaviorAliasOrArray(jsdoc)
	default:
		return d.shouldPreserveBase()
	}
}

func (d *potentialDeclaration) isPolymerBehaviorAliasOrArray(jsdoc *rhino.JSDocInfo) bool {
	return jsdoc != nil &&
		jsdoc.IsPolymerBehavior() &&
		!d.rhs.IsNil() &&
		(d.rhs.IsName() || d.rhs.IsArrayLit())
}

func (d *potentialDeclaration) isStringKeyInNamespace() bool {
	objLit := d.lhs.GetParent()
	if objLit.IsNil() {
		return false
	}
	lvalue := NodeUtil_GetBestLValue(objLit)
	if lvalue.IsNil() {
		return false
	}
	jsdoc := NodeUtil_GetBestJSDocInfo(lvalue)
	return !isExportLhs(lvalue) &&
		!hasAnnotatedType(jsdoc) &&
		NodeUtil_IsNamespaceDecl(lvalue)
}

func (d *potentialDeclaration) isConstToBeInferred() bool {
	return isConstNodeToBeInferred(d.lhs)
}

func isConstNodeToBeInferred(nameNode rhino.Node) bool {
	if nameNode.IsNil() {
		return false
	}
	jsdoc := NodeUtil_GetBestJSDocInfo(nameNode)
	parent := nameNode.GetParent()
	isConst :=
		(!parent.IsNil() && parent.IsConst()) ||
			isExportLhs(nameNode) ||
			(jsdoc != nil && jsdoc.IsConstant())
	return isConst && !hasAnnotatedType(jsdoc) && !NodeUtil_IsNamespaceDecl(nameNode)
}

var (
	ijsGoogAbstractMethodQName = rhino.QualifiedName_Of("goog.abstractMethod")
	ijsGoogRequireQName        = rhino.QualifiedName_Of("goog.require")
	ijsGoogRequireTypeQName    = rhino.QualifiedName_Of("goog.requireType")
	ijsGoogForwardDeclareQName = rhino.QualifiedName_Of("goog.forwardDeclare")
	ijsModuleExportsQName      = rhino.QualifiedName_Of("module.exports")
	ijsGoogModuleGetQName      = rhino.QualifiedName_Of("goog.module.get")
)

func isTypedIjsRhs(rhs rhino.Node) bool {
	if rhs.IsNil() {
		return false
	}
	return rhs.IsFunction() ||
		rhs.IsClass() ||
		(rhs.IsQualifiedName() && ijsGoogAbstractMethodQName.Matches(rhs))
}

func isExportLhs(lhs rhino.Node) bool {
	if lhs.IsNil() {
		return false
	}
	if name, ok := ast.AsName(lhs); ok && name.Ident() == "exports" {
		return true
	}
	if gp, ok := ast.AsGetProp(lhs); ok && gp.Object().Node().MatchesName("exports") {
		return true
	}
	return ijsModuleExportsQName.Matches(lhs)
}

func isImportRhs(rhs rhino.Node) bool {
	call, ok := ast.AsCall(rhs)
	if !ok {
		return false
	}
	callee := call.Callee().Node()
	return ijsGoogRequireQName.Matches(callee) ||
		ijsGoogRequireTypeQName.Matches(callee) ||
		ijsGoogForwardDeclareQName.Matches(callee) ||
		callee.MatchesName("require")
}

func isAliasDeclaration(lhs, rhs rhino.Node) bool {
	return !isThisPropInsideNamedClass(lhs) &&
		isConstNodeToBeInferred(lhs) &&
		!rhs.IsNil() &&
		isQualifiedAliasExpression(rhs)
}

// isQualifiedAliasExpression reports whether n is a simple or qualified alias
// expression such as `x`, `a.b.c`, `this.a`, or `goog.module.get(...)`.
func isQualifiedAliasExpression(n rhino.Node) bool {
	if n.IsNil() {
		return false
	}
	switch n.GetToken() {
	case rhino.Token_NAME:
		return ast.MustName(n).Ident() != ""
	case rhino.Token_THIS, rhino.Token_SUPER:
		return true
	case rhino.Token_GETPROP:
		return isQualifiedAliasExpression(ast.MustGetProp(n).Object().Node())
	case rhino.Token_CALL:
		return ijsGoogModuleGetQName.Matches(ast.MustCall(n).Callee().Node())
	default:
		return false
	}
}

func replaceStringKeyValueWithZero(sk ast.StringKey) {
	value := sk.Value().Node()
	value.ReplaceWith(rhino.IR_Number(0).SrcrefTree(value))
}

func truncateStringKeyValue(sk ast.StringKey, strVal string) {
	value := sk.Value().Node()
	truncated := rhino.IR_String(jutil.StringSubstring(strVal, 0, stringEnumRetainCap-2) + "..").SrcrefTree(value)
	value.ReplaceWith(truncated)
}
