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
// Port of com/google/javascript/jscomp/ijs/JsdocUtil.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.ijs.JsdocUtil
type JsdocUtil struct {
	idHash_ atomic.Int32
}

func (this_ *JsdocUtil) BindSelf_(s any) {
	_ = s
}
func (this_ *JsdocUtil) IsJavaNil_() bool { return this_ == nil }
func (*JsdocUtil) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ijs.JsdocUtil", "JsdocUtil"
}
func (this_ *JsdocUtil) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JsdocUtil) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ijs.JsdocUtil", "com.google.javascript.jscomp.ijs.JsdocUtil", func(v_ any) bool {
		return j2gort.InstanceOf[*JsdocUtil](v_)
	})
}

func StaticGet_JsdocUtil_SYNTHETIC_FILE_NAME() string { return "<synthetic>" }

var static_JsdocUtil_SYNTHETIC_FILE pkg_rhino.StaticSourceFile

func StaticGet_JsdocUtil_SYNTHETIC_FILE() pkg_rhino.StaticSourceFile {
	clinit_JsdocUtil()
	return static_JsdocUtil_SYNTHETIC_FILE
}
func StaticSet_JsdocUtil_SYNTHETIC_FILE(v pkg_rhino.StaticSourceFile) {
	clinit_JsdocUtil()
	static_JsdocUtil_SYNTHETIC_FILE = v
}

var static_JsdocUtil_SYNTETIC_SRCINFO_NODE pkg_rhino.Node

func StaticGet_JsdocUtil_SYNTETIC_SRCINFO_NODE() pkg_rhino.Node {
	clinit_JsdocUtil()
	return static_JsdocUtil_SYNTETIC_SRCINFO_NODE
}
func StaticSet_JsdocUtil_SYNTETIC_SRCINFO_NODE(v pkg_rhino.Node) {
	clinit_JsdocUtil()
	static_JsdocUtil_SYNTETIC_SRCINFO_NODE = v
}

var clinitDone_JsdocUtil uint32
var clinitState_JsdocUtil j2gort.ClinitState

func clinit_JsdocUtil() {
	if atomic.LoadUint32(&clinitDone_JsdocUtil) != 0 {
		return
	}
	clinitSlow_JsdocUtil()
}

func clinitSlow_JsdocUtil() {
	if !j2gort.EnterClinit(&clinitState_JsdocUtil, &clinitDone_JsdocUtil) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_JsdocUtil, &clinitDone_JsdocUtil)
	static_JsdocUtil_SYNTHETIC_FILE = pkg_rhino.NewSimpleSourceFile("<synthetic>", pkg_rhino.StaticGet_StaticSourceFile_SourceKind_EXTERN())
	static_JsdocUtil_SYNTETIC_SRCINFO_NODE = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_EMPTY()).SetStaticSourceFile(StaticGet_JsdocUtil_SYNTHETIC_FILE())
}

func (this_ *JsdocUtil) initFields_JsdocUtil() {
	_ = this_
}

// Java constructor: JsdocUtil()
func NewJsdocUtil() *JsdocUtil {
	clinit_JsdocUtil()
	this_ := &JsdocUtil{}
	InitJsdocUtil(this_)
	return this_
}

func InitJsdocUtil(this_ *JsdocUtil) {
	this_.initFields_JsdocUtil()
}

// Java static method: JsdocUtil.isPrivate(com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo jsdoc)
func JsdocUtil_IsPrivate(jsdoc *pkg_rhino.JSDocInfo) bool {
	clinit_JsdocUtil()
	_ = jsdoc
	return ((jsdoc != nil) && func() bool {
		re_1 := jsdoc.GetVisibility()
		ae_2 := pkg_rhino.StaticGet_JSDocInfo_Visibility_PRIVATE()
		if re_1 == 0 {
			j2gort.ThrowNPE()
		}
		return re_1 == ae_2
	}())
}

// Java static method: JsdocUtil.getUnusableTypeJSDoc(com.google.javascript.rhino.JSDocInfo oldJSDoc)
func JsdocUtil_GetUnusableTypeJSDoc(oldJSDoc *pkg_rhino.JSDocInfo) *pkg_rhino.JSDocInfo {
	clinit_JsdocUtil()
	_ = oldJSDoc
	return JsdocUtil_GetConstJSDoc_2_JSDocInfo_string(oldJSDoc, "UnusableType")
}

// Java static method: JsdocUtil.getQmarkTypeJSDoc()
func JsdocUtil_GetQmarkTypeJSDoc() *pkg_rhino.JSDocInfo {
	clinit_JsdocUtil()
	return JsdocUtil_MakeBuilderWithType(nil, pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_QMARK())).Build_0()
}

// Java static method: JsdocUtil.makeBuilderWithType(com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo oldJSDoc, com.google.javascript.rhino.Node typeAst)
func JsdocUtil_MakeBuilderWithType(oldJSDoc *pkg_rhino.JSDocInfo, typeAst pkg_rhino.Node) *pkg_rhino.JSDocInfo_Builder {
	clinit_JsdocUtil()
	_ = oldJSDoc
	_ = typeAst
	var builder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(oldJSDoc)
	_ = builder
	builder.RecordType(pkg_rhino.NewJSTypeExpression(typeAst.SrcrefTree(StaticGet_JsdocUtil_SYNTETIC_SRCINFO_NODE()), "<synthetic>"))
	return builder
}

// Java static method: JsdocUtil.replaceWithBottomType(com.google.javascript.rhino.JSTypeExpression original)
func JsdocUtil_ReplaceWithBottomType(original pkg_rhino.JSTypeExpression) pkg_rhino.JSTypeExpression {
	clinit_JsdocUtil()
	_ = original
	if JsdocUtil_IsUnknownType(original) {
		return original
	}
	return JsdocUtil_NewTypeExpr(pkg_rhino.Node_NewString_1("IjsNoneType"))
}

// Java static method: JsdocUtil.replaceWithUnusableType(com.google.javascript.rhino.JSTypeExpression original)
func JsdocUtil_ReplaceWithUnusableType(original pkg_rhino.JSTypeExpression) pkg_rhino.JSTypeExpression {
	clinit_JsdocUtil()
	_ = original
	if JsdocUtil_IsUnknownType(original) {
		return original
	}
	return JsdocUtil_NewTypeExpr(pkg_rhino.Node_NewString_1("UnusableType"))
}

// Java static method: JsdocUtil.newTypeExpr(com.google.javascript.rhino.Node typeAst)
func JsdocUtil_NewTypeExpr(typeAst pkg_rhino.Node) pkg_rhino.JSTypeExpression {
	clinit_JsdocUtil()
	_ = typeAst
	return pkg_rhino.NewJSTypeExpression(typeAst.SrcrefTree(StaticGet_JsdocUtil_SYNTETIC_SRCINFO_NODE()), "<synthetic>")
}

// Java static method: JsdocUtil.getConstJSDoc(com.google.javascript.rhino.JSDocInfo oldJSDoc, java.lang.String contents)
func JsdocUtil_GetConstJSDoc_2_JSDocInfo_string(oldJSDoc *pkg_rhino.JSDocInfo, contents string) *pkg_rhino.JSDocInfo {
	clinit_JsdocUtil()
	_ = oldJSDoc
	_ = contents
	return JsdocUtil_GetConstJSDoc_2_JSDocInfo_Node(oldJSDoc, pkg_rhino.Node_NewString_1(contents))
}

// Java static method: JsdocUtil.getConstJSDoc(com.google.javascript.rhino.JSDocInfo oldJSDoc, com.google.javascript.rhino.Node typeAst)
func JsdocUtil_GetConstJSDoc_2_JSDocInfo_Node(oldJSDoc *pkg_rhino.JSDocInfo, typeAst pkg_rhino.Node) *pkg_rhino.JSDocInfo {
	clinit_JsdocUtil()
	_ = oldJSDoc
	_ = typeAst
	var builder *pkg_rhino.JSDocInfo_Builder = JsdocUtil_MakeBuilderWithType(oldJSDoc, typeAst)
	_ = builder
	builder.RecordConstancy()
	return builder.Build_0()
}

// Java static method: JsdocUtil.markConstant(com.google.javascript.rhino.JSDocInfo oldJSDoc)
func JsdocUtil_MarkConstant(oldJSDoc *pkg_rhino.JSDocInfo) *pkg_rhino.JSDocInfo {
	clinit_JsdocUtil()
	_ = oldJSDoc
	var builder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(oldJSDoc)
	_ = builder
	builder.RecordConstancy()
	return builder.Build_0()
}

// Java static method: JsdocUtil.mergeJsdocs(com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo classicJsdoc, com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo inlineJsdoc)
func JsdocUtil_MergeJsdocs(classicJsdoc *pkg_rhino.JSDocInfo, inlineJsdoc *pkg_rhino.JSDocInfo) *pkg_rhino.JSDocInfo {
	clinit_JsdocUtil()
	_ = classicJsdoc
	_ = inlineJsdoc
	if (inlineJsdoc == nil) || !(inlineJsdoc.HasType()) {
		return classicJsdoc
	}
	var builder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(classicJsdoc)
	_ = builder
	builder.RecordType(inlineJsdoc.GetType())
	return builder.Build_0()
}

// Java static method: JsdocUtil.replaceAnnotatedTypesWithUnusableType(com.google.javascript.rhino.JSDocInfo jsdoc)
func JsdocUtil_ReplaceAnnotatedTypesWithUnusableType(jsdoc *pkg_rhino.JSDocInfo) *pkg_rhino.JSDocInfo {
	clinit_JsdocUtil()
	_ = jsdoc
	var builder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_S_Builder()
	_ = builder
	if jsdoc.HasType() {
		builder.RecordType(JsdocUtil_ReplaceWithUnusableType(jsdoc.GetType()))
	}
	if jsdoc.HasReturnType() {
		builder.RecordReturnType(JsdocUtil_ReplaceWithUnusableType(jsdoc.GetReturnType()))
	}
	{
		var i int32 = int32(0)
		_ = i
		for ; i < jsdoc.GetParameterCount(); i = i + 1 {
			var name string = jsdoc.GetParameterNameAt(i)
			_ = name
			var original pkg_rhino.JSTypeExpression = jsdoc.GetParameterType(name)
			_ = original
			if original == (pkg_rhino.JSTypeExpression{}) {
				continue
			}
			var newType pkg_rhino.JSTypeExpression = JsdocUtil_ReplaceWithBottomType(original)
			_ = newType
			if original.IsOptionalArg() {
				newType = pkg_rhino.JSTypeExpression_MakeOptionalArg(newType)
			} else if original.IsVarArgs() {
				newType = pkg_rhino.JSTypeExpression_MakeVarArgs(newType)
			}
			builder.RecordParameter(name, newType)
		}
	}
	if jsdoc.HasThisType() {
		builder.RecordThisType(JsdocUtil_ReplaceWithUnusableType(jsdoc.GetThisType()))
	}
	if jsdoc.IsConstant() {
		builder.RecordConstancy()
	}
	if jsdoc.GetVisibility() != 0 {
		builder.RecordVisibility(jsdoc.GetVisibility())
	}
	return builder.Build_0()
}

// Java static method: JsdocUtil.hasAnnotatedType(com.google.javascript.rhino.JSDocInfo jsdoc)
func JsdocUtil_HasAnnotatedType(jsdoc *pkg_rhino.JSDocInfo) bool {
	clinit_JsdocUtil()
	_ = jsdoc
	if jsdoc == nil {
		return false
	}
	return ((((((jsdoc.HasType() || jsdoc.HasReturnType()) || (jsdoc.GetParameterCount() > int32(0))) || jsdoc.IsConstructorOrInterface()) || jsdoc.HasTypedefType()) || jsdoc.HasThisType()) || jsdoc.HasEnumParameterType())
}

// Java static method: JsdocUtil.getJSDocForRhs(com.google.javascript.rhino.Node rhs, com.google.javascript.rhino.JSDocInfo oldJSDoc)
func JsdocUtil_GetJSDocForRhs(rhs pkg_rhino.Node, oldJSDoc *pkg_rhino.JSDocInfo) *pkg_rhino.JSDocInfo {
	clinit_JsdocUtil()
	_ = rhs
	_ = oldJSDoc
	switch (NodeUtil_GetKnownValueType(rhs)).Ordinal() {
	case int32(6) /* ValueType.BOOLEAN */ :
		return JsdocUtil_GetConstJSDoc_2_JSDocInfo_string(oldJSDoc, "boolean")
	case int32(3) /* ValueType.NUMBER */ :
		return JsdocUtil_GetConstJSDoc_2_JSDocInfo_string(oldJSDoc, "number")
	case int32(4) /* ValueType.BIGINT */ :
		return JsdocUtil_GetConstJSDoc_2_JSDocInfo_string(oldJSDoc, "bigint")
	case int32(5) /* ValueType.STRING */ :
		return JsdocUtil_GetConstJSDoc_2_JSDocInfo_string(oldJSDoc, "string")
	case int32(1) /* ValueType.NULL */ :
		return JsdocUtil_GetConstJSDoc_2_JSDocInfo_string(oldJSDoc, "null")
	case int32(2) /* ValueType.VOID */ :
		return JsdocUtil_GetConstJSDoc_2_JSDocInfo_string(oldJSDoc, "void")
	case int32(7) /* ValueType.OBJECT */ :
		if rhs.IsRegExp() {
			return JsdocUtil_GetConstJSDoc_2_JSDocInfo_Node(oldJSDoc, pkg_rhino.NewNode_2(pkg_rhino.StaticGet_Token_BANG(), pkg_rhino.IR_String("RegExp")))
		}
	case int32(0) /* ValueType.UNDETERMINED */ :
		if (oldJSDoc != nil) && (oldJSDoc.GetDescription() != j2gort.NullString) {
			return JsdocUtil_GetConstJSDoc_2_JSDocInfo_string(oldJSDoc, "string")
		}
	}
	if rhs.IsCast() {
		return JsdocUtil_GetConstJSDoc_2_JSDocInfo_Node(oldJSDoc, rhs.GetJSDocInfo().GetType().GetRoot())
	}
	return nil
}

// Java static method: JsdocUtil.getJSDocForName(com.google.javascript.jscomp.Var decl, com.google.javascript.rhino.JSDocInfo oldJSDoc)
func JsdocUtil_GetJSDocForName(decl *Var, oldJSDoc *pkg_rhino.JSDocInfo) *pkg_rhino.JSDocInfo {
	clinit_JsdocUtil()
	_ = decl
	_ = oldJSDoc
	if decl == nil {
		return nil
	}
	var expr pkg_rhino.JSTypeExpression = NodeUtil_GetDeclaredTypeExpression(decl.GetNameNode())
	_ = expr
	if expr == (pkg_rhino.JSTypeExpression{}) {
		return nil
	}
	var typeAst pkg_rhino.Node = expr.GetRoot()
	_ = typeAst
	switch typeAst.TokenSwitchOrd_() {
	case int32(160) /* Token.EQUALS */ :
		var typeRoot pkg_rhino.Node = typeAst.GetFirstChild().CloneTree_0()
		_ = typeRoot
		if !(decl.IsDefaultParam()) {
			typeRoot = pkg_rhino.NewNode_3(pkg_rhino.StaticGet_Token_PIPE(), typeRoot, pkg_rhino.IR_String("undefined"))
		}
		typeAst = typeRoot
	case int32(125) /* Token.ITER_REST */ :
		var newType pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_BANG())
		_ = newType
		var array pkg_rhino.Node = pkg_rhino.IR_String("Array")
		_ = array
		newType.AddChildToBack(array)
		var block pkg_rhino.Node = pkg_rhino.NewNode_2(pkg_rhino.StaticGet_Token_BLOCK(), typeAst.GetFirstChild().CloneTree_0())
		_ = block
		array.AddChildToBack(block)
		typeAst = newType
	default:
	}
	return JsdocUtil_GetConstJSDoc_2_JSDocInfo_Node(oldJSDoc, typeAst)
}

// Java static method: JsdocUtil.isUnknownType(com.google.javascript.rhino.JSTypeExpression type)
func JsdocUtil_IsUnknownType(type_ pkg_rhino.JSTypeExpression) bool {
	clinit_JsdocUtil()
	_ = type_
	var typeAst pkg_rhino.Node = type_.GetRoot()
	_ = typeAst
	return ((typeAst.RawToken_() == uint8(158)) && !(typeAst.HasChildren()))
}
