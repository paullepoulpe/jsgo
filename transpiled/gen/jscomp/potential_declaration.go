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
// Port of com/google/javascript/jscomp/ijs/PotentialDeclaration.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java polymorphic class interface: com.google.javascript.jscomp.ijs.PotentialDeclaration
type PotentialDeclaration interface {
	Implements_jscomp_PotentialDeclaration()
	IsJavaNil_() bool
	GetPotentialDeclarationBase() *PotentialDeclarationBase
	GetFullyQualifiedName() string
	GetLhs() pkg_rhino.Node
	GetRhs() pkg_rhino.Node
	GetJsDoc() *pkg_rhino.JSDocInfo
	IsDetached() bool
	GetRemovableNode() pkg_rhino.Node
	Remove(AbstractCompiler)
	Simplify(AbstractCompiler)
	BreakDownDestructure(AbstractCompiler) bool
	IsDefiniteDeclaration(AbstractCompiler) bool
	ShouldPreserve() bool
	IsConstToBeInferred() bool
}

func ToIface_PotentialDeclaration[T interface {
	PotentialDeclaration
	comparable
}](v T) PotentialDeclaration {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_PotentialDeclaration[T PotentialDeclaration](v T) PotentialDeclaration {
	var vi PotentialDeclaration = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.ijs.PotentialDeclaration
type PotentialDeclarationBase struct {
	self_              PotentialDeclaration
	fullyQualifiedName string         `gson:"fullyQualifiedName"`
	lhs                pkg_rhino.Node `gson:"lhs"`
	rhs                pkg_rhino.Node `gson:"rhs"`
	idHash_            atomic.Int32
}

func (this_ *PotentialDeclarationBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(PotentialDeclaration)
}
func (this_ *PotentialDeclarationBase) IsJavaNil_() bool { return this_ == nil }
func (*PotentialDeclarationBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ijs.PotentialDeclaration", "PotentialDeclaration"
}
func (this_ *PotentialDeclarationBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ijs.PotentialDeclaration", "com.google.javascript.jscomp.ijs.PotentialDeclaration", func(v_ any) bool {
		return j2gort.InstanceOf[PotentialDeclaration](v_)
	})
}
func (b *PotentialDeclarationBase) GetPotentialDeclarationBase() *PotentialDeclarationBase { return b }
func (*PotentialDeclarationBase) Implements_jscomp_PotentialDeclaration()                  {}

func StaticGet_PotentialDeclaration_STRING_ENUM_RETAIN_CAP() int32 { return int32(10) }

var static_PotentialDeclaration_GOOG_ABSTRACTMETHOD pkg_rhino.QualifiedName

func StaticGet_PotentialDeclaration_GOOG_ABSTRACTMETHOD() pkg_rhino.QualifiedName {
	clinit_PotentialDeclaration()
	return static_PotentialDeclaration_GOOG_ABSTRACTMETHOD
}
func StaticSet_PotentialDeclaration_GOOG_ABSTRACTMETHOD(v pkg_rhino.QualifiedName) {
	clinit_PotentialDeclaration()
	static_PotentialDeclaration_GOOG_ABSTRACTMETHOD = v
}

var static_PotentialDeclaration_GOOG_REQUIRE pkg_rhino.QualifiedName

func StaticGet_PotentialDeclaration_GOOG_REQUIRE() pkg_rhino.QualifiedName {
	clinit_PotentialDeclaration()
	return static_PotentialDeclaration_GOOG_REQUIRE
}
func StaticSet_PotentialDeclaration_GOOG_REQUIRE(v pkg_rhino.QualifiedName) {
	clinit_PotentialDeclaration()
	static_PotentialDeclaration_GOOG_REQUIRE = v
}

var static_PotentialDeclaration_GOOG_REQUIRETYPE pkg_rhino.QualifiedName

func StaticGet_PotentialDeclaration_GOOG_REQUIRETYPE() pkg_rhino.QualifiedName {
	clinit_PotentialDeclaration()
	return static_PotentialDeclaration_GOOG_REQUIRETYPE
}
func StaticSet_PotentialDeclaration_GOOG_REQUIRETYPE(v pkg_rhino.QualifiedName) {
	clinit_PotentialDeclaration()
	static_PotentialDeclaration_GOOG_REQUIRETYPE = v
}

var static_PotentialDeclaration_GOOG_FORWARDDECLARE pkg_rhino.QualifiedName

func StaticGet_PotentialDeclaration_GOOG_FORWARDDECLARE() pkg_rhino.QualifiedName {
	clinit_PotentialDeclaration()
	return static_PotentialDeclaration_GOOG_FORWARDDECLARE
}
func StaticSet_PotentialDeclaration_GOOG_FORWARDDECLARE(v pkg_rhino.QualifiedName) {
	clinit_PotentialDeclaration()
	static_PotentialDeclaration_GOOG_FORWARDDECLARE = v
}

var static_PotentialDeclaration_MODULE_EXPORTS pkg_rhino.QualifiedName

func StaticGet_PotentialDeclaration_MODULE_EXPORTS() pkg_rhino.QualifiedName {
	clinit_PotentialDeclaration()
	return static_PotentialDeclaration_MODULE_EXPORTS
}
func StaticSet_PotentialDeclaration_MODULE_EXPORTS(v pkg_rhino.QualifiedName) {
	clinit_PotentialDeclaration()
	static_PotentialDeclaration_MODULE_EXPORTS = v
}

var static_PotentialDeclaration_GOOG_MODULE_GET pkg_rhino.QualifiedName

func StaticGet_PotentialDeclaration_GOOG_MODULE_GET() pkg_rhino.QualifiedName {
	clinit_PotentialDeclaration()
	return static_PotentialDeclaration_GOOG_MODULE_GET
}
func StaticSet_PotentialDeclaration_GOOG_MODULE_GET(v pkg_rhino.QualifiedName) {
	clinit_PotentialDeclaration()
	static_PotentialDeclaration_GOOG_MODULE_GET = v
}

var clinitDone_PotentialDeclaration uint32
var clinitState_PotentialDeclaration j2gort.ClinitState

func clinit_PotentialDeclaration() {
	if atomic.LoadUint32(&clinitDone_PotentialDeclaration) != 0 {
		return
	}
	clinitSlow_PotentialDeclaration()
}

func clinitSlow_PotentialDeclaration() {
	if !j2gort.EnterClinit(&clinitState_PotentialDeclaration, &clinitDone_PotentialDeclaration) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_PotentialDeclaration, &clinitDone_PotentialDeclaration)
	static_PotentialDeclaration_GOOG_ABSTRACTMETHOD = pkg_rhino.QualifiedName_Of("goog.abstractMethod")
	static_PotentialDeclaration_GOOG_REQUIRE = pkg_rhino.QualifiedName_Of("goog.require")
	static_PotentialDeclaration_GOOG_REQUIRETYPE = pkg_rhino.QualifiedName_Of("goog.requireType")
	static_PotentialDeclaration_GOOG_FORWARDDECLARE = pkg_rhino.QualifiedName_Of("goog.forwardDeclare")
	static_PotentialDeclaration_MODULE_EXPORTS = pkg_rhino.QualifiedName_Of("module.exports")
	static_PotentialDeclaration_GOOG_MODULE_GET = pkg_rhino.QualifiedName_Of("goog.module.get")
}

func (this_ *PotentialDeclarationBase) initFields_PotentialDeclaration() {
	_ = this_
	this_.fullyQualifiedName = j2gort.NullString
}

func InitPotentialDeclaration(this_ *PotentialDeclarationBase, fullyQualifiedName string, lhs pkg_rhino.Node, rhs pkg_rhino.Node) {
	_ = fullyQualifiedName
	_ = lhs
	_ = rhs
	this_.initFields_PotentialDeclaration()
	this_.fullyQualifiedName = j2gort.CheckNotNullString(fullyQualifiedName)
	this_.lhs = j2gort.CheckNotNullNode(lhs)
	this_.rhs = rhs
}

// Java static method: PotentialDeclaration.fromName(com.google.javascript.rhino.Node nameNode)
func PotentialDeclaration_FromName(nameNode pkg_rhino.Node) PotentialDeclaration {
	_ = nameNode
	if !(nameNode.IsQualifiedName()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(nameNode))
	}
	var rhs pkg_rhino.Node = NodeUtil_GetRValueOfLValue(nameNode)
	_ = rhs
	if ClassUtil_IsThisPropInsideClassWithName(nameNode) {
		var name string = ClassUtil_GetFullyQualifiedNameOfThisProp(nameNode)
		_ = name
		return NewPotentialDeclaration_ThisPropDeclaration(name, nameNode, rhs)
	}
	return NewPotentialDeclaration_NameDeclaration(nameNode.GetQualifiedName(), nameNode, rhs)
}

// Java static method: PotentialDeclaration.fromMemberFieldDef(com.google.javascript.rhino.Node fieldNode)
func PotentialDeclaration_FromMemberFieldDef(fieldNode pkg_rhino.Node) PotentialDeclaration {
	_ = fieldNode
	if !(fieldNode.IsMemberFieldDef()) {
		j2gort.CheckArgument(false)
	}
	var rhs pkg_rhino.Node = NodeUtil_GetRValueOfLValue(fieldNode)
	_ = rhs
	var name string = ClassUtil_GetFullyQualifiedNameOfMemberFieldDef(fieldNode)
	_ = name
	return NewPotentialDeclaration_MemberFieldDeclaration(name, fieldNode, rhs)
}

// Java static method: PotentialDeclaration.fromMethod(com.google.javascript.rhino.Node functionNode)
func PotentialDeclaration_FromMethod(functionNode pkg_rhino.Node) PotentialDeclaration {
	_ = functionNode
	if !(ClassUtil_IsClassMethod(functionNode)) {
		j2gort.CheckArgument(false)
	}
	var name string = ClassUtil_GetFullyQualifiedNameOfMethod(functionNode)
	_ = name
	return NewPotentialDeclaration_MethodDeclaration(name, functionNode)
}

// Java static method: PotentialDeclaration.fromStringKey(com.google.javascript.rhino.Node stringKeyNode)
func PotentialDeclaration_FromStringKey(stringKeyNode pkg_rhino.Node) PotentialDeclaration {
	_ = stringKeyNode
	if !(stringKeyNode.IsStringKey()) {
		j2gort.CheckArgument(false)
	}
	if !(stringKeyNode.GetParent().IsObjectLit()) {
		j2gort.CheckArgument(false)
	}
	var name string = ("this." + j2gort.StrValString(stringKeyNode.GetString()))
	_ = name
	if j2gort.CheckNotNullString(stringKeyNode.GetString()) == "properties" {
		var objLitJsDoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(stringKeyNode.GetParent())
		_ = objLitJsDoc
		if (objLitJsDoc != nil) && objLitJsDoc.IsPolymerBehavior() {
			return NewPotentialDeclaration_PolymerBehaviorPropertiesDeclaration(name, stringKeyNode)
		}
	}
	return NewPotentialDeclaration_StringKeyDeclaration(name, stringKeyNode)
}

// Java static method: PotentialDeclaration.fromDefine(com.google.javascript.rhino.Node callNode)
func PotentialDeclaration_FromDefine(callNode pkg_rhino.Node) PotentialDeclaration {
	_ = callNode
	if !(NodeUtil_IsCallTo_2_Node_string(callNode, "goog.define")) {
		j2gort.CheckArgument(false)
	}
	return ToIface_PotentialDeclaration(PotentialDeclaration_DefineDeclaration_From(callNode))
}

// Java static method: PotentialDeclaration.fromAlias(com.google.javascript.rhino.Node nameNode)
func PotentialDeclaration_FromAlias(nameNode pkg_rhino.Node) PotentialDeclaration {
	_ = nameNode
	if !(nameNode.IsQualifiedName()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(nameNode))
	}
	return NewPotentialDeclaration_AliasDeclaration(nameNode.GetQualifiedName(), nameNode)
}

// Java method: PotentialDeclaration.getFullyQualifiedName()
func (this_ *PotentialDeclarationBase) GetFullyQualifiedName() string {
	return this_.fullyQualifiedName
}

// Java method: PotentialDeclaration.getLhs()
func (this_ *PotentialDeclarationBase) GetLhs() pkg_rhino.Node {
	return this_.lhs
}

// Java method: PotentialDeclaration.getRhs()
func (this_ *PotentialDeclarationBase) GetRhs() pkg_rhino.Node {
	return this_.rhs
}

// Java method: PotentialDeclaration.getJsDoc()
func (this_ *PotentialDeclarationBase) GetJsDoc() *pkg_rhino.JSDocInfo {
	return NodeUtil_GetBestJSDocInfo(this_.lhs)
}

// Java method: PotentialDeclaration.isDetached()
func (this_ *PotentialDeclarationBase) IsDetached() bool {
	{
		var current pkg_rhino.Node = this_.lhs
		_ = current
		for ; current != 0; current = current.GetParent() {
			if current.IsScript() {
				return false
			}
		}
	}
	return true
}

// Java method: PotentialDeclaration.getRemovableNode()
func (this_ *PotentialDeclarationBase) GetRemovableNode() pkg_rhino.Node {
	return NodeUtil_GetEnclosingStatement(this_.lhs)
}

// Java method: PotentialDeclaration.remove(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *PotentialDeclarationBase) Remove(compiler AbstractCompiler) {
	_ = compiler
	if this_.IsDetached() {
		return
	}
	var statement pkg_rhino.Node = this_.self_.GetRemovableNode()
	_ = statement
	NodeUtil_DeleteNode(statement, compiler)
	statement.RemoveChildren()
}

// Java method: PotentialDeclaration.simplify(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *PotentialDeclarationBase) Simplify(compiler AbstractCompiler) {
	_ = compiler
	panic("abstract or native method invoked")
}

// Java method: PotentialDeclaration.breakDownDestructure(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *PotentialDeclarationBase) BreakDownDestructure(compiler AbstractCompiler) bool {
	_ = compiler
	return false
}

// Java method: PotentialDeclaration.simplifyEnumValues(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *PotentialDeclarationBase) priv_PotentialDeclaration_simplifyEnumValues_6422(compiler AbstractCompiler) {
	_ = compiler
	if this_.GetRhs().IsObjectLit() && this_.GetRhs().HasChildren() {
		var changed bool = false
		_ = changed
		{
			var key pkg_rhino.Node = this_.GetRhs().GetFirstChild()
			_ = key
			for ; key != 0; key = key.GetNext() {
				var value pkg_rhino.Node = key.GetOnlyChild()
				_ = value
				if !(value.IsStringLit()) {
					PotentialDeclaration_RemoveStringKeyValue(key)
					changed = true
				} else {
					var content string = value.GetString()
					_ = content
					if j2gort.JStrLen(content) > int32(10) {
						PotentialDeclaration_TruncateStringKeyValue(key)
						changed = true
					}
				}
			}
		}
		if changed {
			compiler.ReportChangeToEnclosingScope(this_.GetRhs())
		}
	}
}

// Java method: PotentialDeclaration.isDefiniteDeclaration(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *PotentialDeclarationBase) IsDefiniteDeclaration(compiler AbstractCompiler) bool {
	_ = compiler
	var parent pkg_rhino.Node = this_.GetLhs().GetParent()
	_ = parent
	switch parent.TokenSwitchOrd_() {
	case int32(134) /* Token.DEFAULT_VALUE */ :
		if !(parent.GetParent().IsStringKey()) {
			return false
		}
		fallthrough
	case int32(129) /* Token.COMPUTED_PROP */, int32(100) /* Token.STRING_KEY */ :
		if NodeUtil_IsLhsByDestructuring(this_.GetLhs()) {
			var rootTarget pkg_rhino.Node = NodeUtil_GetRootTarget(this_.GetLhs())
			_ = rootTarget
			if !(rootTarget.GetParent().IsDestructuringLhs()) {
				j2gort.CheckState(false)
			}
			if NodeUtil_IsNameDeclaration(rootTarget.GetGrandparent()) {
				return true
			}
		}
		return false
	case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */, int32(105) /* Token.CLASS */, int32(73) /* Token.FUNCTION */ :
		return true
	default:
		return ((PotentialDeclaration_IsExportLhs(this_.GetLhs()) || ((this_.GetJsDoc() != nil) && this_.GetJsDoc().ContainsDeclaration())) || ((this_.GetRhs() != 0) && PotentialDeclaration_IsTypedRhs(this_.GetRhs())))
	}
	panic("unreachable")
}

// Java method: PotentialDeclaration.shouldPreserve()
func (this_ *PotentialDeclarationBase) ShouldPreserve() bool {
	return ((this_.GetRhs() != 0) && PotentialDeclaration_IsTypedRhs(this_.GetRhs()))
}

// Java method: PotentialDeclaration.isConstToBeInferred()
func (this_ *PotentialDeclarationBase) IsConstToBeInferred() bool {
	return PotentialDeclaration_IsConstToBeInferred(this_.GetLhs())
}

// Java static method: PotentialDeclaration.isConstToBeInferred(com.google.javascript.rhino.Node nameNode)
func PotentialDeclaration_IsConstToBeInferred(nameNode pkg_rhino.Node) bool {
	_ = nameNode
	var jsdoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(nameNode)
	_ = jsdoc
	var isConst bool = ((nameNode.GetParent().IsConst() || PotentialDeclaration_IsExportLhs(nameNode)) || ((jsdoc != nil) && jsdoc.IsConstant()))
	_ = isConst
	return ((isConst && !(JsdocUtil_HasAnnotatedType(jsdoc))) && !(NodeUtil_IsNamespaceDecl(nameNode)))
}

// Java static method: PotentialDeclaration.isTypedRhs(com.google.javascript.rhino.Node rhs)
func PotentialDeclaration_IsTypedRhs(rhs pkg_rhino.Node) bool {
	_ = rhs
	return ((rhs.IsFunction() || rhs.IsClass()) || (rhs.IsQualifiedName() && StaticGet_PotentialDeclaration_GOOG_ABSTRACTMETHOD().Matches(rhs)))
}

// Java static method: PotentialDeclaration.isExportLhs(com.google.javascript.rhino.Node lhs)
func PotentialDeclaration_IsExportLhs(lhs pkg_rhino.Node) bool {
	_ = lhs
	return (((lhs.IsName() && lhs.MatchesName_1_string("exports")) || (lhs.IsGetProp() && lhs.GetFirstChild().MatchesName_1_string("exports"))) || StaticGet_PotentialDeclaration_MODULE_EXPORTS().Matches(lhs))
}

// Java static method: PotentialDeclaration.isImportRhs(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node rhs)
func PotentialDeclaration_IsImportRhs(rhs pkg_rhino.Node) bool {
	_ = rhs
	if (rhs == 0) || !(rhs.IsCall()) {
		return false
	}
	var callee pkg_rhino.Node = rhs.GetFirstChild()
	_ = callee
	return (((StaticGet_PotentialDeclaration_GOOG_REQUIRE().Matches(callee) || StaticGet_PotentialDeclaration_GOOG_REQUIRETYPE().Matches(callee)) || StaticGet_PotentialDeclaration_GOOG_FORWARDDECLARE().Matches(callee)) || callee.MatchesName_1_string("require"))
}

// Java static method: PotentialDeclaration.isAliasDeclaration(com.google.javascript.rhino.Node lhs, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node rhs)
func PotentialDeclaration_IsAliasDeclaration(lhs pkg_rhino.Node, rhs pkg_rhino.Node) bool {
	_ = lhs
	_ = rhs
	return (((!(ClassUtil_IsThisPropInsideClassWithName(lhs)) && PotentialDeclaration_IsConstToBeInferred(lhs)) && (rhs != 0)) && PotentialDeclaration_IsQualifiedAliasExpression(rhs))
}

// Java static method: PotentialDeclaration.isQualifiedAliasExpression(com.google.javascript.rhino.Node n)
func PotentialDeclaration_IsQualifiedAliasExpression(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		return !(j2gort.CheckNotNullString(n.GetString()) == "")
	case int32(37) /* Token.THIS */, int32(110) /* Token.SUPER */ :
		return true
	case int32(26) /* Token.GETPROP */ :
		return PotentialDeclaration_IsQualifiedAliasExpression(n.GetFirstChild())
	case int32(28) /* Token.CALL */ :
		if StaticGet_PotentialDeclaration_GOOG_MODULE_GET().Matches(n.GetFirstChild()) {
			return true
		}
		return false
	case int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
		fallthrough
	default:
		return false
	}
	panic("unreachable")
}

// Java static method: PotentialDeclaration.removeStringKeyValue(com.google.javascript.rhino.Node stringKey)
func PotentialDeclaration_RemoveStringKeyValue(stringKey pkg_rhino.Node) {
	_ = stringKey
	var value pkg_rhino.Node = stringKey.GetOnlyChild()
	_ = value
	var replacementValue pkg_rhino.Node = pkg_rhino.IR_Number(float64(int32(0))).SrcrefTree(value)
	_ = replacementValue
	value.ReplaceWith(replacementValue)
}

// Java static method: PotentialDeclaration.truncateStringKeyValue(com.google.javascript.rhino.Node stringKey)
func PotentialDeclaration_TruncateStringKeyValue(stringKey pkg_rhino.Node) {
	_ = stringKey
	var value pkg_rhino.Node = stringKey.GetOnlyChild()
	_ = value
	var replacementValue pkg_rhino.Node = pkg_rhino.IR_String((j2gort.StrValString(j2gort.JStrSub(value.GetString(), int32(0), (int32(10)-int32(2)))) + "..")).SrcrefTree(value)
	_ = replacementValue
	value.ReplaceWith(replacementValue)
}

// Java class: com.google.javascript.jscomp.ijs.PotentialDeclaration.NameDeclaration
type PotentialDeclaration_NameDeclaration struct {
	PotentialDeclarationBase
}

func (this_ *PotentialDeclaration_NameDeclaration) BindSelf_(s any) {
	_ = s
	this_.PotentialDeclarationBase.BindSelf_(s)
}
func (this_ *PotentialDeclaration_NameDeclaration) IsJavaNil_() bool { return this_ == nil }
func (*PotentialDeclaration_NameDeclaration) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ijs.PotentialDeclaration$NameDeclaration", "NameDeclaration"
}
func (this_ *PotentialDeclaration_NameDeclaration) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.PotentialDeclarationBase.IdentityHash_()
}
func (*PotentialDeclaration_NameDeclaration) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ijs.PotentialDeclaration.NameDeclaration", "com.google.javascript.jscomp.ijs.PotentialDeclaration$NameDeclaration", func(v_ any) bool {
		return j2gort.InstanceOf[*PotentialDeclaration_NameDeclaration](v_)
	})
}
func (*PotentialDeclaration_NameDeclaration) Implements_jscomp_PotentialDeclaration() {}

func (this_ *PotentialDeclaration_NameDeclaration) initFields_PotentialDeclaration_NameDeclaration() {
	_ = this_
}

// Java constructor: NameDeclaration(java.lang.String fullyQualifiedName, com.google.javascript.rhino.Node lhs, com.google.javascript.rhino.Node rhs)
func NewPotentialDeclaration_NameDeclaration(fullyQualifiedName string, lhs pkg_rhino.Node, rhs pkg_rhino.Node) *PotentialDeclaration_NameDeclaration {
	this_ := &PotentialDeclaration_NameDeclaration{}
	this_.BindSelf_(this_)
	InitPotentialDeclaration_NameDeclaration(this_, fullyQualifiedName, lhs, rhs)
	return this_
}

func InitPotentialDeclaration_NameDeclaration(this_ *PotentialDeclaration_NameDeclaration, fullyQualifiedName string, lhs pkg_rhino.Node, rhs pkg_rhino.Node) {
	_ = fullyQualifiedName
	_ = lhs
	_ = rhs
	InitPotentialDeclaration(&this_.PotentialDeclarationBase, fullyQualifiedName, lhs, rhs)
	this_.initFields_PotentialDeclaration_NameDeclaration()
}

// Java method: NameDeclaration.simplifyNamespace(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *PotentialDeclaration_NameDeclaration) priv_NameDeclaration_simplifyNamespace_6426(compiler AbstractCompiler) {
	_ = compiler
	if this_.GetRhs().IsOr() {
		var objLit pkg_rhino.Node = this_.GetRhs().GetLastChild().Detach()
		_ = objLit
		this_.GetRhs().ReplaceWith(objLit)
		compiler.ReportChangeToEnclosingScope(this_.GetLhs())
	}
}

// Java method: NameDeclaration.simplifySymbol(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *PotentialDeclaration_NameDeclaration) priv_NameDeclaration_simplifySymbol_6427(compiler AbstractCompiler) {
	_ = compiler
	if !(NodeUtil_IsCallTo_2_Node_string(this_.GetRhs(), "Symbol")) {
		j2gort.CheckArgument(false)
	}
	var callNode pkg_rhino.Node = this_.GetRhs()
	_ = callNode
	for callNode.HasMoreThanOneChild() {
		NodeUtil_DeleteNode(callNode.GetLastChild(), compiler)
	}
}

// Java method: NameDeclaration.simplify(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *PotentialDeclaration_NameDeclaration) Simplify(compiler AbstractCompiler) {
	_ = compiler
	if (this_.GetRhs() == 0) || this_.ShouldPreserve() {
		return
	}
	var nameNode pkg_rhino.Node = this_.GetLhs()
	_ = nameNode
	var jsdoc *pkg_rhino.JSDocInfo = this_.GetJsDoc()
	_ = jsdoc
	if (jsdoc != nil) && jsdoc.HasEnumParameterType() {
		this_.PotentialDeclarationBase.GetPotentialDeclarationBase().priv_PotentialDeclaration_simplifyEnumValues_6422(compiler)
		return
	}
	if NodeUtil_IsNamespaceDecl(nameNode) {
		this_.priv_NameDeclaration_simplifyNamespace_6426(compiler)
		return
	}
	if nameNode.MatchesName_1_string("exports") {
		PotentialDeclaration_NameDeclaration_ReplaceRhsWithUnknown(this_.GetRhs())
		compiler.ReportChangeToEnclosingScope(nameNode)
		return
	}
	if NodeUtil_IsCallTo_2_Node_string(this_.GetRhs(), "Symbol") {
		this_.priv_NameDeclaration_simplifySymbol_6427(compiler)
		return
	}
	if this_.GetLhs().GetParent().IsConst() {
		jsdoc = JsdocUtil_MarkConstant(jsdoc)
	}
	var newStatement pkg_rhino.Node = NodeUtil_NewQNameDeclaration_4_AbstractCompiler_string_Node_JSDocInfo(compiler, nameNode.GetQualifiedName(), pkg_rhino.Node(0), jsdoc)
	_ = newStatement
	newStatement.SrcrefTreeIfMissing(nameNode)
	var oldStatement pkg_rhino.Node = this_.GetRemovableNode()
	_ = oldStatement
	NodeUtil_DeleteChildren(oldStatement, compiler)
	if oldStatement.IsExport() {
		oldStatement.AddChildToBack(newStatement)
	} else {
		oldStatement.ReplaceWith(newStatement)
	}
	compiler.ReportChangeToEnclosingScope(newStatement)
}

// Java method: NameDeclaration.breakDownDestructure(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *PotentialDeclaration_NameDeclaration) BreakDownDestructure(compiler AbstractCompiler) bool {
	_ = compiler
	if !(NodeUtil_IsLhsByDestructuring(this_.GetLhs())) {
		return false
	}
	var rootTarget pkg_rhino.Node = NodeUtil_GetRootTarget(this_.GetLhs())
	_ = rootTarget
	if !(rootTarget.GetParent().IsDestructuringLhs()) {
		j2gort.CheckState(false)
	}
	if !(NodeUtil_IsNameDeclaration(rootTarget.GetGrandparent())) {
		return false
	}
	var definitionNode pkg_rhino.Node = rootTarget.GetGrandparent()
	_ = definitionNode
	var prev pkg_rhino.Node = pkg_rhino.Node(0)
	_ = prev
	var lhsNodes j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList()
	_ = lhsNodes
	NodeUtil_VisitLhsNodesInNode(rootTarget.GetParent(), func() j2gort.Ext_java_util_function_Consumer {
		bound_1 := lhsNodes
		return j2gort.NewLambda_Ext_java_util_function_Consumer(func(mrp0_ any) { bound_1.Add_1(mrp0_) })
	}())
	{
		cur_2 := j2gort.NewIterCursor(lhsNodes)
		for cur_2.HasNext() {
			n := pkg_rhino.Node(cur_2.NextNode())
			_ = n
			n.Detach()
			var temp pkg_rhino.Node = pkg_rhino.IR_Var_1(n)
			_ = temp
			if prev == 0 {
				definitionNode.ReplaceWith(temp)
				compiler.ReportChangeToEnclosingScope(temp)
			} else {
				temp.InsertAfter(prev)
				compiler.ReportChangeToEnclosingScope(temp)
				temp.SrcrefTree(prev)
			}
			prev = temp
		}
	}
	return true
}

// Java static method: NameDeclaration.replaceRhsWithUnknown(com.google.javascript.rhino.Node rhs)
func PotentialDeclaration_NameDeclaration_ReplaceRhsWithUnknown(rhs pkg_rhino.Node) {
	_ = rhs
	rhs.ReplaceWith(pkg_rhino.IR_Cast(pkg_rhino.IR_Number(float64(int32(0))), JsdocUtil_GetQmarkTypeJSDoc()).SrcrefTree(rhs))
}

// Java method: NameDeclaration.shouldPreserve()
func (this_ *PotentialDeclaration_NameDeclaration) ShouldPreserve() bool {
	var rhs pkg_rhino.Node = this_.GetRhs()
	_ = rhs
	var nameNode pkg_rhino.Node = this_.GetLhs()
	_ = nameNode
	var jsdoc *pkg_rhino.JSDocInfo = this_.GetJsDoc()
	_ = jsdoc
	var isExport bool = PotentialDeclaration_IsExportLhs(nameNode)
	_ = isExport
	return ((((((this_.PotentialDeclarationBase.ShouldPreserve() || PotentialDeclaration_IsImportRhs(rhs)) || ((isExport && (rhs != 0)) && (rhs.IsQualifiedName() || rhs.IsObjectLit()))) || ((((jsdoc != nil) && jsdoc.IsConstructor()) && (rhs != 0)) && rhs.IsQualifiedName())) || ((((rhs != 0) && rhs.IsObjectLit()) && !(rhs.HasChildren())) && ((jsdoc == nil) || !(JsdocUtil_HasAnnotatedType(jsdoc))))) || ((rhs != 0) && NodeUtil_IsCallTo_2_Node_string(rhs, "Polymer"))) || this_.IsPolymerBehaviorAliasOrArray())
}

// Java method: NameDeclaration.isPolymerBehaviorAliasOrArray()
func (this_ *PotentialDeclaration_NameDeclaration) IsPolymerBehaviorAliasOrArray() bool {
	var jsdoc *pkg_rhino.JSDocInfo = this_.GetJsDoc()
	_ = jsdoc
	var rhs pkg_rhino.Node = this_.GetRhs()
	_ = rhs
	return ((((jsdoc != nil) && jsdoc.IsPolymerBehavior()) && (rhs != 0)) && (rhs.IsName() || rhs.IsArrayLit()))
}

// Java class: com.google.javascript.jscomp.ijs.PotentialDeclaration.ThisPropDeclaration
type PotentialDeclaration_ThisPropDeclaration struct {
	PotentialDeclarationBase
	insertionPoint pkg_rhino.Node `gson:"insertionPoint"`
}

func (this_ *PotentialDeclaration_ThisPropDeclaration) BindSelf_(s any) {
	_ = s
	this_.PotentialDeclarationBase.BindSelf_(s)
}
func (this_ *PotentialDeclaration_ThisPropDeclaration) IsJavaNil_() bool { return this_ == nil }
func (*PotentialDeclaration_ThisPropDeclaration) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ijs.PotentialDeclaration$ThisPropDeclaration", "ThisPropDeclaration"
}
func (this_ *PotentialDeclaration_ThisPropDeclaration) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.PotentialDeclarationBase.IdentityHash_()
}
func (*PotentialDeclaration_ThisPropDeclaration) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ijs.PotentialDeclaration.ThisPropDeclaration", "com.google.javascript.jscomp.ijs.PotentialDeclaration$ThisPropDeclaration", func(v_ any) bool {
		return j2gort.InstanceOf[*PotentialDeclaration_ThisPropDeclaration](v_)
	})
}
func (*PotentialDeclaration_ThisPropDeclaration) Implements_jscomp_PotentialDeclaration() {}

func (this_ *PotentialDeclaration_ThisPropDeclaration) initFields_PotentialDeclaration_ThisPropDeclaration() {
	_ = this_
}

// Java constructor: ThisPropDeclaration(java.lang.String fullyQualifiedName, com.google.javascript.rhino.Node lhs, com.google.javascript.rhino.Node rhs)
func NewPotentialDeclaration_ThisPropDeclaration(fullyQualifiedName string, lhs pkg_rhino.Node, rhs pkg_rhino.Node) *PotentialDeclaration_ThisPropDeclaration {
	this_ := &PotentialDeclaration_ThisPropDeclaration{}
	this_.BindSelf_(this_)
	InitPotentialDeclaration_ThisPropDeclaration(this_, fullyQualifiedName, lhs, rhs)
	return this_
}

func InitPotentialDeclaration_ThisPropDeclaration(this_ *PotentialDeclaration_ThisPropDeclaration, fullyQualifiedName string, lhs pkg_rhino.Node, rhs pkg_rhino.Node) {
	_ = fullyQualifiedName
	_ = lhs
	_ = rhs
	InitPotentialDeclaration(&this_.PotentialDeclarationBase, fullyQualifiedName, lhs, rhs)
	this_.initFields_PotentialDeclaration_ThisPropDeclaration()
	var thisPropDefinition pkg_rhino.Node = NodeUtil_GetEnclosingStatement(lhs)
	_ = thisPropDefinition
	this_.insertionPoint = NodeUtil_GetEnclosingStatement(thisPropDefinition.GetParent())
}

// Java method: ThisPropDeclaration.simplify(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *PotentialDeclaration_ThisPropDeclaration) Simplify(compiler AbstractCompiler) {
	_ = compiler
	if this_.ShouldPreserve() {
		return
	}
	var newStatement pkg_rhino.Node = NodeUtil_NewQNameDeclaration_4_AbstractCompiler_string_Node_JSDocInfo(compiler, this_.GetFullyQualifiedName(), pkg_rhino.Node(0), this_.GetJsDoc())
	_ = newStatement
	newStatement.SrcrefTreeIfMissing(this_.GetLhs())
	NodeUtil_DeleteNode(this_.GetRemovableNode(), compiler)
	if this_.insertionPoint.HasParent() {
		newStatement.InsertAfter(this_.insertionPoint)
		compiler.ReportChangeToEnclosingScope(newStatement)
	}
}

// Java class: com.google.javascript.jscomp.ijs.PotentialDeclaration.DefineDeclaration
type PotentialDeclaration_DefineDeclaration struct {
	PotentialDeclarationBase
}

func (this_ *PotentialDeclaration_DefineDeclaration) BindSelf_(s any) {
	_ = s
	this_.PotentialDeclarationBase.BindSelf_(s)
}
func (this_ *PotentialDeclaration_DefineDeclaration) IsJavaNil_() bool { return this_ == nil }
func (*PotentialDeclaration_DefineDeclaration) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ijs.PotentialDeclaration$DefineDeclaration", "DefineDeclaration"
}
func (this_ *PotentialDeclaration_DefineDeclaration) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.PotentialDeclarationBase.IdentityHash_()
}
func (*PotentialDeclaration_DefineDeclaration) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ijs.PotentialDeclaration.DefineDeclaration", "com.google.javascript.jscomp.ijs.PotentialDeclaration$DefineDeclaration", func(v_ any) bool {
		return j2gort.InstanceOf[*PotentialDeclaration_DefineDeclaration](v_)
	})
}
func (*PotentialDeclaration_DefineDeclaration) Implements_jscomp_PotentialDeclaration() {}

func (this_ *PotentialDeclaration_DefineDeclaration) initFields_PotentialDeclaration_DefineDeclaration() {
	_ = this_
}

// Java constructor: DefineDeclaration(java.lang.String qualifiedName, com.google.javascript.rhino.Node lhs, com.google.javascript.rhino.Node rhs)
func NewPotentialDeclaration_DefineDeclaration(qualifiedName string, lhs pkg_rhino.Node, rhs pkg_rhino.Node) *PotentialDeclaration_DefineDeclaration {
	this_ := &PotentialDeclaration_DefineDeclaration{}
	this_.BindSelf_(this_)
	InitPotentialDeclaration_DefineDeclaration(this_, qualifiedName, lhs, rhs)
	return this_
}

func InitPotentialDeclaration_DefineDeclaration(this_ *PotentialDeclaration_DefineDeclaration, qualifiedName string, lhs pkg_rhino.Node, rhs pkg_rhino.Node) {
	_ = qualifiedName
	_ = lhs
	_ = rhs
	InitPotentialDeclaration(&this_.PotentialDeclarationBase, qualifiedName, lhs, rhs)
	this_.initFields_PotentialDeclaration_DefineDeclaration()
}

// Java method: DefineDeclaration.simplify(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *PotentialDeclaration_DefineDeclaration) Simplify(compiler AbstractCompiler) {
	_ = compiler
	var info *pkg_rhino.JSDocInfo = this_.GetJsDoc()
	_ = info
	if (info != nil) && (info.GetType() != (pkg_rhino.JSTypeExpression{})) {
		var newRhs pkg_rhino.Node = PotentialDeclaration_DefineDeclaration_MakeEmptyValueNode(info.GetType())
		_ = newRhs
		if newRhs != 0 {
			this_.GetRhs().ReplaceWith(newRhs)
			compiler.ReportChangeToEnclosingScope(newRhs)
			return
		}
	}
	NodeUtil_DeleteNode(this_.GetRemovableNode(), compiler)
}

// Java static method: DefineDeclaration.from(com.google.javascript.rhino.Node callNode)
func PotentialDeclaration_DefineDeclaration_From(callNode pkg_rhino.Node) *PotentialDeclaration_DefineDeclaration {
	_ = callNode
	switch callNode.GetParent().TokenSwitchOrd_() {
	case int32(93) /* Token.EXPR_RESULT */ :
		return NewPotentialDeclaration_DefineDeclaration(callNode.GetSecondChild().GetString(), callNode, callNode.GetLastChild())
	case int32(51) /* Token.ASSIGN */ :
		var previous pkg_rhino.Node = callNode.GetPrevious()
		_ = previous
		return NewPotentialDeclaration_DefineDeclaration(previous.GetQualifiedName(), previous, callNode.GetLastChild())
	case int32(32) /* Token.NAME */ :
		var parent pkg_rhino.Node = callNode.GetParent()
		_ = parent
		return NewPotentialDeclaration_DefineDeclaration(parent.GetString(), parent, callNode.GetLastChild())
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected parent: "+j2gort.StrVal(callNode.GetParent().GetToken()))))
	}
	panic("unreachable switch expression")
}

// Java static method: DefineDeclaration.makeEmptyValueNode(com.google.javascript.rhino.JSTypeExpression type)
func PotentialDeclaration_DefineDeclaration_MakeEmptyValueNode(type_ pkg_rhino.JSTypeExpression) pkg_rhino.Node {
	_ = type_
	var n pkg_rhino.Node = type_.GetRoot()
	_ = n
	for ((n != 0) && !(n.IsStringLit())) && !(n.IsName()) {
		n = n.GetFirstChild()
	}
	switch j2gort.CheckNotNullString(func() string {
		if n != 0 {
			return n.GetString()
		}
		return ""
	}()) {
	case "boolean":
		return pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_FALSE())
	case "number":
		return pkg_rhino.Node_NewNumber(float64(int32(0)))
	case "string":
		return pkg_rhino.Node_NewString_1("")
	default:
		return pkg_rhino.Node(0)
	}
	panic("unreachable switch expression")
}

// Java class: com.google.javascript.jscomp.ijs.PotentialDeclaration.MethodDeclaration
type PotentialDeclaration_MethodDeclaration struct {
	PotentialDeclarationBase
}

func (this_ *PotentialDeclaration_MethodDeclaration) BindSelf_(s any) {
	_ = s
	this_.PotentialDeclarationBase.BindSelf_(s)
}
func (this_ *PotentialDeclaration_MethodDeclaration) IsJavaNil_() bool { return this_ == nil }
func (*PotentialDeclaration_MethodDeclaration) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ijs.PotentialDeclaration$MethodDeclaration", "MethodDeclaration"
}
func (this_ *PotentialDeclaration_MethodDeclaration) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.PotentialDeclarationBase.IdentityHash_()
}
func (*PotentialDeclaration_MethodDeclaration) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ijs.PotentialDeclaration.MethodDeclaration", "com.google.javascript.jscomp.ijs.PotentialDeclaration$MethodDeclaration", func(v_ any) bool {
		return j2gort.InstanceOf[*PotentialDeclaration_MethodDeclaration](v_)
	})
}
func (*PotentialDeclaration_MethodDeclaration) Implements_jscomp_PotentialDeclaration() {}

func (this_ *PotentialDeclaration_MethodDeclaration) initFields_PotentialDeclaration_MethodDeclaration() {
	_ = this_
}

// Java constructor: MethodDeclaration(java.lang.String name, com.google.javascript.rhino.Node functionNode)
func NewPotentialDeclaration_MethodDeclaration(name string, functionNode pkg_rhino.Node) *PotentialDeclaration_MethodDeclaration {
	this_ := &PotentialDeclaration_MethodDeclaration{}
	this_.BindSelf_(this_)
	InitPotentialDeclaration_MethodDeclaration(this_, name, functionNode)
	return this_
}

func InitPotentialDeclaration_MethodDeclaration(this_ *PotentialDeclaration_MethodDeclaration, name string, functionNode pkg_rhino.Node) {
	_ = name
	_ = functionNode
	InitPotentialDeclaration(&this_.PotentialDeclarationBase, name, functionNode.GetParent(), functionNode)
	this_.initFields_PotentialDeclaration_MethodDeclaration()
}

// Java method: MethodDeclaration.simplify(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *PotentialDeclaration_MethodDeclaration) Simplify(compiler AbstractCompiler) {
	_ = compiler
}

// Java method: MethodDeclaration.getRemovableNode()
func (this_ *PotentialDeclaration_MethodDeclaration) GetRemovableNode() pkg_rhino.Node {
	return this_.GetLhs()
}

// Java class: com.google.javascript.jscomp.ijs.PotentialDeclaration.MemberFieldDeclaration
type PotentialDeclaration_MemberFieldDeclaration struct {
	PotentialDeclarationBase
}

func (this_ *PotentialDeclaration_MemberFieldDeclaration) BindSelf_(s any) {
	_ = s
	this_.PotentialDeclarationBase.BindSelf_(s)
}
func (this_ *PotentialDeclaration_MemberFieldDeclaration) IsJavaNil_() bool { return this_ == nil }
func (*PotentialDeclaration_MemberFieldDeclaration) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ijs.PotentialDeclaration$MemberFieldDeclaration", "MemberFieldDeclaration"
}
func (this_ *PotentialDeclaration_MemberFieldDeclaration) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.PotentialDeclarationBase.IdentityHash_()
}
func (*PotentialDeclaration_MemberFieldDeclaration) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ijs.PotentialDeclaration.MemberFieldDeclaration", "com.google.javascript.jscomp.ijs.PotentialDeclaration$MemberFieldDeclaration", func(v_ any) bool {
		return j2gort.InstanceOf[*PotentialDeclaration_MemberFieldDeclaration](v_)
	})
}
func (*PotentialDeclaration_MemberFieldDeclaration) Implements_jscomp_PotentialDeclaration() {}

func (this_ *PotentialDeclaration_MemberFieldDeclaration) initFields_PotentialDeclaration_MemberFieldDeclaration() {
	_ = this_
}

// Java constructor: MemberFieldDeclaration(java.lang.@org.jspecify.annotations.Nullable String name, com.google.javascript.rhino.Node lhs, com.google.javascript.rhino.Node rhs)
func NewPotentialDeclaration_MemberFieldDeclaration(name string, lhs pkg_rhino.Node, rhs pkg_rhino.Node) *PotentialDeclaration_MemberFieldDeclaration {
	this_ := &PotentialDeclaration_MemberFieldDeclaration{}
	this_.BindSelf_(this_)
	InitPotentialDeclaration_MemberFieldDeclaration(this_, name, lhs, rhs)
	return this_
}

func InitPotentialDeclaration_MemberFieldDeclaration(this_ *PotentialDeclaration_MemberFieldDeclaration, name string, lhs pkg_rhino.Node, rhs pkg_rhino.Node) {
	_ = name
	_ = lhs
	_ = rhs
	InitPotentialDeclaration(&this_.PotentialDeclarationBase, name, lhs, rhs)
	this_.initFields_PotentialDeclaration_MemberFieldDeclaration()
}

// Java method: MemberFieldDeclaration.simplify(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *PotentialDeclaration_MemberFieldDeclaration) Simplify(compiler AbstractCompiler) {
	_ = compiler
	var rhs pkg_rhino.Node = this_.GetRhs()
	_ = rhs
	if rhs != 0 {
		NodeUtil_DeleteNode(rhs, compiler)
	}
}

// Java method: MemberFieldDeclaration.getRemovableNode()
func (this_ *PotentialDeclaration_MemberFieldDeclaration) GetRemovableNode() pkg_rhino.Node {
	return this_.GetLhs()
}

// Java class: com.google.javascript.jscomp.ijs.PotentialDeclaration.StringKeyDeclaration
type PotentialDeclaration_StringKeyDeclaration struct {
	PotentialDeclarationBase
}

func (this_ *PotentialDeclaration_StringKeyDeclaration) BindSelf_(s any) {
	_ = s
	this_.PotentialDeclarationBase.BindSelf_(s)
}
func (this_ *PotentialDeclaration_StringKeyDeclaration) IsJavaNil_() bool { return this_ == nil }
func (*PotentialDeclaration_StringKeyDeclaration) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ijs.PotentialDeclaration$StringKeyDeclaration", "StringKeyDeclaration"
}
func (this_ *PotentialDeclaration_StringKeyDeclaration) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.PotentialDeclarationBase.IdentityHash_()
}
func (*PotentialDeclaration_StringKeyDeclaration) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ijs.PotentialDeclaration.StringKeyDeclaration", "com.google.javascript.jscomp.ijs.PotentialDeclaration$StringKeyDeclaration", func(v_ any) bool {
		return j2gort.InstanceOf[*PotentialDeclaration_StringKeyDeclaration](v_)
	})
}
func (*PotentialDeclaration_StringKeyDeclaration) Implements_jscomp_PotentialDeclaration() {}

func (this_ *PotentialDeclaration_StringKeyDeclaration) initFields_PotentialDeclaration_StringKeyDeclaration() {
	_ = this_
}

// Java constructor: StringKeyDeclaration(java.lang.String name, com.google.javascript.rhino.Node stringKeyNode)
func NewPotentialDeclaration_StringKeyDeclaration(name string, stringKeyNode pkg_rhino.Node) *PotentialDeclaration_StringKeyDeclaration {
	this_ := &PotentialDeclaration_StringKeyDeclaration{}
	this_.BindSelf_(this_)
	InitPotentialDeclaration_StringKeyDeclaration(this_, name, stringKeyNode)
	return this_
}

func InitPotentialDeclaration_StringKeyDeclaration(this_ *PotentialDeclaration_StringKeyDeclaration, name string, stringKeyNode pkg_rhino.Node) {
	_ = name
	_ = stringKeyNode
	InitPotentialDeclaration(&this_.PotentialDeclarationBase, name, stringKeyNode, stringKeyNode.GetLastChild())
	this_.initFields_PotentialDeclaration_StringKeyDeclaration()
}

// Java method: StringKeyDeclaration.simplify(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *PotentialDeclaration_StringKeyDeclaration) Simplify(compiler AbstractCompiler) {
	_ = compiler
	if this_.ShouldPreserve() {
		return
	}
	var jsdoc *pkg_rhino.JSDocInfo = this_.GetJsDoc()
	_ = jsdoc
	if (jsdoc != nil) && jsdoc.HasEnumParameterType() {
		this_.PotentialDeclarationBase.GetPotentialDeclarationBase().priv_PotentialDeclaration_simplifyEnumValues_6422(compiler)
		return
	}
	var key pkg_rhino.Node = this_.GetLhs()
	_ = key
	PotentialDeclaration_RemoveStringKeyValue(key)
	compiler.ReportChangeToEnclosingScope(key)
	if ((jsdoc == nil) || !(jsdoc.ContainsDeclaration())) || this_.IsConstToBeInferred() {
		key.SetJSDocInfo(JsdocUtil_GetUnusableTypeJSDoc(jsdoc))
	}
}

// Java method: StringKeyDeclaration.shouldPreserve()
func (this_ *PotentialDeclaration_StringKeyDeclaration) ShouldPreserve() bool {
	return ((this_.PotentialDeclarationBase.IsDetached() || this_.PotentialDeclarationBase.ShouldPreserve()) || !(this_.priv_StringKeyDeclaration_isInNamespace_6429()))
}

// Java method: StringKeyDeclaration.isInNamespace()
func (this_ *PotentialDeclaration_StringKeyDeclaration) priv_StringKeyDeclaration_isInNamespace_6429() bool {
	var stringKey pkg_rhino.Node = this_.GetLhs()
	_ = stringKey
	var objLit pkg_rhino.Node = stringKey.GetParent()
	_ = objLit
	var lvalue pkg_rhino.Node = NodeUtil_GetBestLValue(objLit)
	_ = lvalue
	if lvalue == 0 {
		return false
	}
	var jsdoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(lvalue)
	_ = jsdoc
	return ((!(PotentialDeclaration_IsExportLhs(lvalue)) && !(JsdocUtil_HasAnnotatedType(jsdoc))) && NodeUtil_IsNamespaceDecl(lvalue))
}

// Java method: StringKeyDeclaration.getRemovableNode()
func (this_ *PotentialDeclaration_StringKeyDeclaration) GetRemovableNode() pkg_rhino.Node {
	return this_.GetLhs()
}

// Java class: com.google.javascript.jscomp.ijs.PotentialDeclaration.PolymerBehaviorPropertiesDeclaration
type PotentialDeclaration_PolymerBehaviorPropertiesDeclaration struct {
	PotentialDeclarationBase
}

func (this_ *PotentialDeclaration_PolymerBehaviorPropertiesDeclaration) BindSelf_(s any) {
	_ = s
	this_.PotentialDeclarationBase.BindSelf_(s)
}
func (this_ *PotentialDeclaration_PolymerBehaviorPropertiesDeclaration) IsJavaNil_() bool {
	return this_ == nil
}
func (*PotentialDeclaration_PolymerBehaviorPropertiesDeclaration) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ijs.PotentialDeclaration$PolymerBehaviorPropertiesDeclaration", "PolymerBehaviorPropertiesDeclaration"
}
func (this_ *PotentialDeclaration_PolymerBehaviorPropertiesDeclaration) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.PotentialDeclarationBase.IdentityHash_()
}
func (*PotentialDeclaration_PolymerBehaviorPropertiesDeclaration) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ijs.PotentialDeclaration.PolymerBehaviorPropertiesDeclaration", "com.google.javascript.jscomp.ijs.PotentialDeclaration$PolymerBehaviorPropertiesDeclaration", func(v_ any) bool {
		return j2gort.InstanceOf[*PotentialDeclaration_PolymerBehaviorPropertiesDeclaration](v_)
	})
}
func (*PotentialDeclaration_PolymerBehaviorPropertiesDeclaration) Implements_jscomp_PotentialDeclaration() {
}

func (this_ *PotentialDeclaration_PolymerBehaviorPropertiesDeclaration) initFields_PotentialDeclaration_PolymerBehaviorPropertiesDeclaration() {
	_ = this_
}

// Java constructor: PolymerBehaviorPropertiesDeclaration(java.lang.String name, com.google.javascript.rhino.Node stringKeyNode)
func NewPotentialDeclaration_PolymerBehaviorPropertiesDeclaration(name string, stringKeyNode pkg_rhino.Node) *PotentialDeclaration_PolymerBehaviorPropertiesDeclaration {
	this_ := &PotentialDeclaration_PolymerBehaviorPropertiesDeclaration{}
	this_.BindSelf_(this_)
	InitPotentialDeclaration_PolymerBehaviorPropertiesDeclaration(this_, name, stringKeyNode)
	return this_
}

func InitPotentialDeclaration_PolymerBehaviorPropertiesDeclaration(this_ *PotentialDeclaration_PolymerBehaviorPropertiesDeclaration, name string, stringKeyNode pkg_rhino.Node) {
	_ = name
	_ = stringKeyNode
	InitPotentialDeclaration(&this_.PotentialDeclarationBase, name, stringKeyNode, stringKeyNode.GetLastChild())
	this_.initFields_PotentialDeclaration_PolymerBehaviorPropertiesDeclaration()
}

// Java method: PolymerBehaviorPropertiesDeclaration.simplify(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *PotentialDeclaration_PolymerBehaviorPropertiesDeclaration) Simplify(compiler AbstractCompiler) {
	_ = compiler
	if this_.IsDetached() {
		return
	}
	var propertiesObject pkg_rhino.Node = this_.GetRhs()
	_ = propertiesObject
	if !(propertiesObject.IsObjectLit()) || !(propertiesObject.HasChildren()) {
		return
	}
	{
		var propKey pkg_rhino.Node = propertiesObject.GetFirstChild()
		_ = propKey
		for ; propKey != 0; propKey = propKey.GetNext() {
			var propDef pkg_rhino.Node = propKey.GetOnlyChild()
			_ = propDef
			if propDef.IsObjectLit() {
				{
					var subProp pkg_rhino.Node = propDef.GetFirstChild()
					_ = subProp
					for subProp != 0 {
						var next pkg_rhino.Node = subProp.GetNext()
						_ = next
						if !(j2gort.CheckNotNullString(subProp.GetString()) == "type") && !(j2gort.CheckNotNullString(subProp.GetString()) == "readOnly") {
							NodeUtil_DeleteNode(subProp, compiler)
						}
						subProp = next
					}
				}
			}
		}
	}
}

// Java method: PolymerBehaviorPropertiesDeclaration.shouldPreserve()
func (this_ *PotentialDeclaration_PolymerBehaviorPropertiesDeclaration) ShouldPreserve() bool {
	return true
}

// Java method: PolymerBehaviorPropertiesDeclaration.getRemovableNode()
func (this_ *PotentialDeclaration_PolymerBehaviorPropertiesDeclaration) GetRemovableNode() pkg_rhino.Node {
	return this_.GetLhs()
}

// Java class: com.google.javascript.jscomp.ijs.PotentialDeclaration.AliasDeclaration
type PotentialDeclaration_AliasDeclaration struct {
	PotentialDeclarationBase
}

func (this_ *PotentialDeclaration_AliasDeclaration) BindSelf_(s any) {
	_ = s
	this_.PotentialDeclarationBase.BindSelf_(s)
}
func (this_ *PotentialDeclaration_AliasDeclaration) IsJavaNil_() bool { return this_ == nil }
func (*PotentialDeclaration_AliasDeclaration) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ijs.PotentialDeclaration$AliasDeclaration", "AliasDeclaration"
}
func (this_ *PotentialDeclaration_AliasDeclaration) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.PotentialDeclarationBase.IdentityHash_()
}
func (*PotentialDeclaration_AliasDeclaration) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ijs.PotentialDeclaration.AliasDeclaration", "com.google.javascript.jscomp.ijs.PotentialDeclaration$AliasDeclaration", func(v_ any) bool {
		return j2gort.InstanceOf[*PotentialDeclaration_AliasDeclaration](v_)
	})
}
func (*PotentialDeclaration_AliasDeclaration) Implements_jscomp_PotentialDeclaration() {}

func (this_ *PotentialDeclaration_AliasDeclaration) initFields_PotentialDeclaration_AliasDeclaration() {
	_ = this_
}

// Java constructor: AliasDeclaration(java.lang.String name, com.google.javascript.rhino.Node lhs)
func NewPotentialDeclaration_AliasDeclaration(name string, lhs pkg_rhino.Node) *PotentialDeclaration_AliasDeclaration {
	this_ := &PotentialDeclaration_AliasDeclaration{}
	this_.BindSelf_(this_)
	InitPotentialDeclaration_AliasDeclaration(this_, name, lhs)
	return this_
}

func InitPotentialDeclaration_AliasDeclaration(this_ *PotentialDeclaration_AliasDeclaration, name string, lhs pkg_rhino.Node) {
	_ = name
	_ = lhs
	InitPotentialDeclaration(&this_.PotentialDeclarationBase, name, lhs, pkg_rhino.Node(0))
	this_.initFields_PotentialDeclaration_AliasDeclaration()
}

// Java method: AliasDeclaration.simplify(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *PotentialDeclaration_AliasDeclaration) Simplify(compiler AbstractCompiler) {
	_ = compiler
}

// Java method: AliasDeclaration.getRemovableNode()
func (this_ *PotentialDeclaration_AliasDeclaration) GetRemovableNode() pkg_rhino.Node {
	var lhs pkg_rhino.Node = this_.GetLhs()
	_ = lhs
	if lhs.GetParent().IsArrayPattern() && lhs.GetParent().HasMoreThanOneChild() {
		return lhs
	}
	if lhs.GetGrandparent().IsObjectPattern() && lhs.GetGrandparent().HasMoreThanOneChild() {
		return lhs.GetParent()
	}
	return NodeUtil_GetEnclosingStatement(lhs)
}

// Java method: AliasDeclaration.isDefiniteDeclaration(com.google.javascript.jscomp.AbstractCompiler compiler)
func (this_ *PotentialDeclaration_AliasDeclaration) IsDefiniteDeclaration(compiler AbstractCompiler) bool {
	_ = compiler
	return true
}

// Java method: AliasDeclaration.shouldPreserve()
func (this_ *PotentialDeclaration_AliasDeclaration) ShouldPreserve() bool {
	return true
}
