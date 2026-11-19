/*
 * Copyright 2021 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/serialization/ScriptNodeDeserializer.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parsing "github.com/paullepoulpe/jsgo/transpiled/gen/parsing"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
)

// Java class: com.google.javascript.jscomp.serialization.ScriptNodeDeserializer
type ScriptNodeDeserializer struct {
	sourceFile       *SourceFile                                `gson:"sourceFile"`
	scriptBytes      j2gort.Ext_com_google_protobuf_ByteString  `gson:"scriptBytes"`
	sourceMappingURL string                                     `gson:"sourceMappingURL"`
	colorPoolShard   j2gort.Ext_com_google_common_base_Optional `gson:"colorPoolShard"`
	stringPool       *StringPool                                `gson:"stringPool"`
	filePool         []*SourceFile                              `gson:"filePool"`
	idHash_          atomic.Int32
}

func (this_ *ScriptNodeDeserializer) BindSelf_(s any) {
	_ = s
}
func (this_ *ScriptNodeDeserializer) IsJavaNil_() bool { return this_ == nil }
func (*ScriptNodeDeserializer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.ScriptNodeDeserializer", "ScriptNodeDeserializer"
}
func (this_ *ScriptNodeDeserializer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ScriptNodeDeserializer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.ScriptNodeDeserializer", "com.google.javascript.jscomp.serialization.ScriptNodeDeserializer", func(v_ any) bool {
		return j2gort.InstanceOf[*ScriptNodeDeserializer](v_)
	})
}

func (this_ *ScriptNodeDeserializer) initFields_ScriptNodeDeserializer() {
	_ = this_
	this_.sourceMappingURL = j2gort.NullString
}

// Java constructor: ScriptNodeDeserializer(com.google.javascript.jscomp.serialization.LazyAst ast, com.google.javascript.jscomp.serialization.StringPool stringPool, com.google.common.base.Optional<com.google.javascript.jscomp.serialization.ColorPool.ShardView> colorPoolShard, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.SourceFile> filePool)
func NewScriptNodeDeserializer(ast *j2gort.Ext_com_google_javascript_jscomp_serialization_LazyAst, stringPool *StringPool, colorPoolShard j2gort.Ext_com_google_common_base_Optional, filePool j2gort.Ext_com_google_common_collect_ImmutableList) *ScriptNodeDeserializer {
	this_ := &ScriptNodeDeserializer{}
	InitScriptNodeDeserializer(this_, ast, stringPool, colorPoolShard, filePool)
	return this_
}

func InitScriptNodeDeserializer(this_ *ScriptNodeDeserializer, ast *j2gort.Ext_com_google_javascript_jscomp_serialization_LazyAst, stringPool *StringPool, colorPoolShard j2gort.Ext_com_google_common_base_Optional, filePool j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = ast
	_ = stringPool
	_ = colorPoolShard
	_ = filePool
	this_.initFields_ScriptNodeDeserializer()
	this_.scriptBytes = ast.GetScript()
	this_.sourceFile = j2gort.CastPtr[SourceFile](filePool.Get((ast.GetSourceFile() - int32(1))))
	this_.sourceMappingURL = ast.GetSourceMappingUrl()
	this_.colorPoolShard = colorPoolShard
	this_.stringPool = stringPool
	this_.filePool = j2gort.ImmListToSliceNullable[

	// Java method: ScriptNodeDeserializer.deserializeNew()
	*SourceFile](filePool)
}

func (this_ *ScriptNodeDeserializer) DeserializeNew() pkg_rhino.Node {
	var runner_ ScriptNodeDeserializer_Runner
	runner_.outer_ = this_
	InitScriptNodeDeserializer_Runner(&runner_)
	return runner_.Run()

}

// Java method: ScriptNodeDeserializer.getSourceMappingURL()
func (this_ *ScriptNodeDeserializer) GetSourceMappingURL() string {
	return this_.sourceMappingURL
}

// Java method: ScriptNodeDeserializer.getSourceFile()
func (this_ *ScriptNodeDeserializer) GetSourceFile() *SourceFile {
	return this_.sourceFile
}

// Java method: ScriptNodeDeserializer.createSourceInfoTemplate(com.google.javascript.jscomp.SourceFile file)
func (this_ *ScriptNodeDeserializer) priv_ScriptNodeDeserializer_createSourceInfoTemplate_7747(file *SourceFile) pkg_rhino.Node {
	_ = file
	var sourceInfoTemplate pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_SCRIPT())
	_ = sourceInfoTemplate
	sourceInfoTemplate.SetStaticSourceFile(pkg_rhino.ToIface_StaticSourceFile(file))
	return sourceInfoTemplate
}

// Java method: ScriptNodeDeserializer.setOriginalNameIfPresent(com.google.javascript.jscomp.serialization.AstNode astNode, com.google.javascript.rhino.Node n)
func (this_ *ScriptNodeDeserializer) priv_ScriptNodeDeserializer_setOriginalNameIfPresent_7748(astNode *j2gort.Ext_com_google_javascript_jscomp_serialization_AstNode, n pkg_rhino.Node, alloc_ *nodestore.NodeAllocator) {
	if off_ := astNode.GetOriginalNamePointer(); off_ != 0 {
		if n == 0 {
			j2gort.ThrowNPE()
		}
		if id_, ok_ := this_.stringPool.internAt_(off_, alloc_); !ok_ {
			nodestore.Node(n).SetOrigNameID(0)
		} else {
			nodestore.Node(n).SetOrigNameID(id_)
		}
	}

}

// Java method: ScriptNodeDeserializer.stringNode(com.google.javascript.rhino.Token token, com.google.javascript.jscomp.serialization.AstNode n, com.google.javascript.rhino.Node templateNode)
func (this_ *ScriptNodeDeserializer) priv_ScriptNodeDeserializer_stringNode_7749(token pkg_rhino.Token, n *j2gort.Ext_com_google_javascript_jscomp_serialization_AstNode, templateNode pkg_rhino.Node, alloc_ *nodestore.NodeAllocator) pkg_rhino.Node {
	if token == 0 {
		j2gort.ThrowNPE()
	}
	var strID_ nodestore.StringID
	if id_, ok_ := this_.stringPool.internAt_(n.GetStringValuePointer(), alloc_); ok_ {
		strID_ = id_
	}
	return pkg_rhino.Node(alloc_.AllocStringNode(uint8(token)-1, strID_, nodestore.Node(templateNode)))

}

// Java method: ScriptNodeDeserializer.newNodeFromTemplate(com.google.javascript.rhino.Token token, com.google.javascript.rhino.Node templateNode)
func (this_ *ScriptNodeDeserializer) priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(token pkg_rhino.Token, templateNode pkg_rhino.Node, alloc_ *nodestore.NodeAllocator) pkg_rhino.Node {
	if token == 0 {
		j2gort.ThrowNPE()
	}
	return pkg_rhino.Node(alloc_.AllocNode1(uint8(token)-1, nodestore.Node(templateNode)))

}

// Java method: ScriptNodeDeserializer.deserializeSingleNode(com.google.javascript.jscomp.serialization.AstNode n, com.google.javascript.rhino.Node templateNode)
func (this_ *ScriptNodeDeserializer) priv_ScriptNodeDeserializer_deserializeSingleNode_7751(n *j2gort.Ext_com_google_javascript_jscomp_serialization_AstNode, templateNode pkg_rhino.Node, alloc_ *nodestore.NodeAllocator) pkg_rhino.Node {
	_ = n
	_ = templateNode
	{
		var numberNode pkg_rhino.Node = pkg_rhino.Node(0)
		_ = numberNode
		var bigintString string = j2gort.NullString
		_ = bigintString
		var bi pkg_rhino.Node = pkg_rhino.Node(0)
		_ = bi
		var postInc pkg_rhino.Node = pkg_rhino.Node(0)
		_ = postInc
		var postDec pkg_rhino.Node = pkg_rhino.Node(0)
		_ = postDec
		var quotedStringKey pkg_rhino.Node = pkg_rhino.Node(0)
		_ = quotedStringKey
		var getterDef pkg_rhino.Node = pkg_rhino.Node(0)
		_ = getterDef
		var setterDef pkg_rhino.Node = pkg_rhino.Node(0)
		_ = setterDef
		switch j2gort.CheckNotNullPtr(n.GetKind()).Ordinal() {
		case int32(107) /* NodeKind.SOURCE_FILE */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_SCRIPT(), templateNode, alloc_)
		case int32(1) /* NodeKind.NUMBER_LITERAL */ :
			d_ := n.GetDoubleValue()
			if d_ != d_ || d_ < 0 || (d_ == 0 && 1/d_ < 0) {

				numberNode = pkg_rhino.IR_Number(n.GetDoubleValue())
				numberNode.ClonePropsFrom(templateNode)
				return numberNode
			}
			return pkg_rhino.Node(alloc_.AllocNumberNode(uint8(pkg_rhino.StaticGet_Token_NUMBER())-1, d_, nodestore.Node(templateNode)))

		case int32(2) /* NodeKind.STRING_LITERAL */ :
			return this_.priv_ScriptNodeDeserializer_stringNode_7749(pkg_rhino.StaticGet_Token_STRINGLIT(), n, templateNode, alloc_)
		case int32(12) /* NodeKind.IDENTIFIER */ :
			return this_.priv_ScriptNodeDeserializer_stringNode_7749(pkg_rhino.StaticGet_Token_NAME(), n, templateNode, alloc_)
		case int32(5) /* NodeKind.FALSE */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_FALSE(), templateNode, alloc_)
		case int32(4) /* NodeKind.TRUE */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_TRUE(), templateNode, alloc_)
		case int32(6) /* NodeKind.NULL */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_NULL(), templateNode, alloc_)
		case int32(13) /* NodeKind.THIS */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_THIS(), templateNode, alloc_)
		case int32(68) /* NodeKind.VOID */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_VOID(), templateNode, alloc_)
		case int32(3) /* NodeKind.BIGINT_LITERAL */ :
			bigintString = this_.stringPool.Get(n.GetStringValuePointer())
			bi = pkg_rhino.IR_Bigint(j2gort.NewExt_java_math_BigInteger_1_string(bigintString))
			bi.ClonePropsFrom(templateNode)
			return bi
		case int32(7) /* NodeKind.REGEX_LITERAL */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_REGEXP(), templateNode, alloc_)
		case int32(8) /* NodeKind.ARRAY_LITERAL */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_ARRAYLIT(), templateNode, alloc_)
		case int32(9) /* NodeKind.OBJECT_LITERAL */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_OBJECTLIT(), templateNode, alloc_)
		case int32(14) /* NodeKind.ASSIGNMENT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_ASSIGN(), templateNode, alloc_)
		case int32(16) /* NodeKind.CALL */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_CALL(), templateNode, alloc_)
		case int32(17) /* NodeKind.NEW */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_NEW(), templateNode, alloc_)
		case int32(23) /* NodeKind.PROPERTY_ACCESS */ :
			return this_.priv_ScriptNodeDeserializer_stringNode_7749(pkg_rhino.StaticGet_Token_GETPROP(), n, templateNode, alloc_)
		case int32(24) /* NodeKind.ELEMENT_ACCESS */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_GETELEM(), templateNode, alloc_)
		case int32(15) /* NodeKind.COMMA */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_COMMA(), templateNode, alloc_)
		case int32(20) /* NodeKind.BOOLEAN_OR */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_OR(), templateNode, alloc_)
		case int32(21) /* NodeKind.BOOLEAN_AND */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_AND(), templateNode, alloc_)
		case int32(22) /* NodeKind.HOOK */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_HOOK(), templateNode, alloc_)
		case int32(29) /* NodeKind.EQUAL */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_EQ(), templateNode, alloc_)
		case int32(31) /* NodeKind.NOT_EQUAL */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_NE(), templateNode, alloc_)
		case int32(25) /* NodeKind.LESS_THAN */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_LT(), templateNode, alloc_)
		case int32(26) /* NodeKind.LESS_THAN_EQUAL */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_LE(), templateNode, alloc_)
		case int32(27) /* NodeKind.GREATER_THAN */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_GT(), templateNode, alloc_)
		case int32(28) /* NodeKind.GREATER_THAN_EQUAL */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_GE(), templateNode, alloc_)
		case int32(30) /* NodeKind.TRIPLE_EQUAL */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_SHEQ(), templateNode, alloc_)
		case int32(32) /* NodeKind.NOT_TRIPLE_EQUAL */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_SHNE(), templateNode, alloc_)
		case int32(33) /* NodeKind.NOT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_NOT(), templateNode, alloc_)
		case int32(62) /* NodeKind.POSITIVE */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_POS(), templateNode, alloc_)
		case int32(63) /* NodeKind.NEGATIVE */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_NEG(), templateNode, alloc_)
		case int32(34) /* NodeKind.TYPEOF */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_TYPEOF(), templateNode, alloc_)
		case int32(35) /* NodeKind.INSTANCEOF */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_INSTANCEOF(), templateNode, alloc_)
		case int32(36) /* NodeKind.IN */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_IN(), templateNode, alloc_)
		case int32(40) /* NodeKind.ADD */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_ADD(), templateNode, alloc_)
		case int32(41) /* NodeKind.SUBTRACT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_SUB(), templateNode, alloc_)
		case int32(42) /* NodeKind.MULTIPLY */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_MUL(), templateNode, alloc_)
		case int32(43) /* NodeKind.DIVIDE */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_DIV(), templateNode, alloc_)
		case int32(44) /* NodeKind.MODULO */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_MOD(), templateNode, alloc_)
		case int32(45) /* NodeKind.EXPONENT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_EXPONENT(), templateNode, alloc_)
		case int32(67) /* NodeKind.BITWISE_NOT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_BITNOT(), templateNode, alloc_)
		case int32(64) /* NodeKind.BITWISE_OR */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_BITOR(), templateNode, alloc_)
		case int32(65) /* NodeKind.BITWISE_AND */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_BITAND(), templateNode, alloc_)
		case int32(66) /* NodeKind.BITWISE_XOR */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_BITXOR(), templateNode, alloc_)
		case int32(37) /* NodeKind.LEFT_SHIFT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_LSH(), templateNode, alloc_)
		case int32(38) /* NodeKind.RIGHT_SHIFT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_RSH(), templateNode, alloc_)
		case int32(39) /* NodeKind.UNSIGNED_RIGHT_SHIFT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_URSH(), templateNode, alloc_)
		case int32(58) /* NodeKind.PRE_INCREMENT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_INC(), templateNode, alloc_)
		case int32(59) /* NodeKind.POST_INCREMENT */ :
			postInc = this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_INC(), templateNode, alloc_)
			postInc.PutBooleanProp(pkg_rhino.StaticGet_Node_INCRDECR_PROP(), true)
			return postInc
		case int32(60) /* NodeKind.PRE_DECREMENT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_DEC(), templateNode, alloc_)
		case int32(61) /* NodeKind.POST_DECREMENT */ :
			postDec = this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_DEC(), templateNode, alloc_)
			postDec.PutBooleanProp(pkg_rhino.StaticGet_Node_INCRDECR_PROP(), true)
			return postDec
		case int32(46) /* NodeKind.ASSIGN_ADD */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_ASSIGN_ADD(), templateNode, alloc_)
		case int32(47) /* NodeKind.ASSIGN_SUBTRACT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_ASSIGN_SUB(), templateNode, alloc_)
		case int32(48) /* NodeKind.ASSIGN_MULTIPLY */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_ASSIGN_MUL(), templateNode, alloc_)
		case int32(49) /* NodeKind.ASSIGN_DIVIDE */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_ASSIGN_DIV(), templateNode, alloc_)
		case int32(50) /* NodeKind.ASSIGN_MODULO */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_ASSIGN_MOD(), templateNode, alloc_)
		case int32(51) /* NodeKind.ASSIGN_EXPONENT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_ASSIGN_EXPONENT(), templateNode, alloc_)
		case int32(52) /* NodeKind.ASSIGN_BITWISE_OR */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_ASSIGN_BITOR(), templateNode, alloc_)
		case int32(53) /* NodeKind.ASSIGN_BITWISE_AND */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_ASSIGN_BITAND(), templateNode, alloc_)
		case int32(54) /* NodeKind.ASSIGN_BITWISE_XOR */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_ASSIGN_BITXOR(), templateNode, alloc_)
		case int32(55) /* NodeKind.ASSIGN_LEFT_SHIFT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_ASSIGN_LSH(), templateNode, alloc_)
		case int32(56) /* NodeKind.ASSIGN_RIGHT_SHIFT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_ASSIGN_RSH(), templateNode, alloc_)
		case int32(57) /* NodeKind.ASSIGN_UNSIGNED_RIGHT_SHIFT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_ASSIGN_URSH(), templateNode, alloc_)
		case int32(18) /* NodeKind.YIELD */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_YIELD(), templateNode, alloc_)
		case int32(19) /* NodeKind.AWAIT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_AWAIT(), templateNode, alloc_)
		case int32(69) /* NodeKind.DELETE */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_DELPROP(), templateNode, alloc_)
		case int32(11) /* NodeKind.TAGGED_TEMPLATELIT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_TAGGED_TEMPLATELIT(), templateNode, alloc_)
		case int32(10) /* NodeKind.TEMPLATELIT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_TEMPLATELIT(), templateNode, alloc_)
		case int32(120) /* NodeKind.TEMPLATELIT_SUB */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_TEMPLATELIT_SUB(), templateNode, alloc_)
		case int32(121) /* NodeKind.TEMPLATELIT_STRING */ :
			{
				var templateStringValue *j2gort.Ext_com_google_javascript_jscomp_serialization_TemplateStringValue = n.GetTemplateStringValue()
				_ = templateStringValue
				var templateLitString pkg_rhino.Node = pkg_rhino.Node_NewTemplateLitString_3(this_.stringPool.GetInternedStrings(), templateStringValue.GetCookedStringPointer(), templateStringValue.GetRawStringPointer())
				_ = templateLitString
				templateLitString.ClonePropsFrom(templateNode)
				return templateLitString
			}
		case int32(70) /* NodeKind.NEW_TARGET */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_NEW_TARGET(), templateNode, alloc_)
		case int32(71) /* NodeKind.COMPUTED_PROP */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_COMPUTED_PROP(), templateNode, alloc_)
		case int32(72) /* NodeKind.IMPORT_META */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_IMPORT_META(), templateNode, alloc_)
		case int32(73) /* NodeKind.OPTCHAIN_PROPERTY_ACCESS */ :
			return this_.priv_ScriptNodeDeserializer_stringNode_7749(pkg_rhino.StaticGet_Token_OPTCHAIN_GETPROP(), n, templateNode, alloc_)
		case int32(74) /* NodeKind.OPTCHAIN_CALL */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_OPTCHAIN_CALL(), templateNode, alloc_)
		case int32(75) /* NodeKind.OPTCHAIN_ELEMENT_ACCESS */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_OPTCHAIN_GETELEM(), templateNode, alloc_)
		case int32(76) /* NodeKind.COALESCE */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_COALESCE(), templateNode, alloc_)
		case int32(77) /* NodeKind.DYNAMIC_IMPORT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_DYNAMIC_IMPORT(), templateNode, alloc_)
		case int32(78) /* NodeKind.ASSIGN_OR */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_ASSIGN_OR(), templateNode, alloc_)
		case int32(79) /* NodeKind.ASSIGN_AND */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_ASSIGN_AND(), templateNode, alloc_)
		case int32(80) /* NodeKind.ASSIGN_COALESCE */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_ASSIGN_COALESCE(), templateNode, alloc_)
		case int32(85) /* NodeKind.EXPRESSION_STATEMENT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_EXPR_RESULT(), templateNode, alloc_)
		case int32(81) /* NodeKind.BREAK_STATEMENT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_BREAK(), templateNode, alloc_)
		case int32(82) /* NodeKind.CONTINUE_STATEMENT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_CONTINUE(), templateNode, alloc_)
		case int32(83) /* NodeKind.DEBUGGER_STATEMENT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_DEBUGGER(), templateNode, alloc_)
		case int32(84) /* NodeKind.DO_STATEMENT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_DO(), templateNode, alloc_)
		case int32(89) /* NodeKind.FOR_STATEMENT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_FOR(), templateNode, alloc_)
		case int32(87) /* NodeKind.FOR_IN_STATEMENT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_FOR_IN(), templateNode, alloc_)
		case int32(88) /* NodeKind.FOR_OF_STATEMENT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_FOR_OF(), templateNode, alloc_)
		case int32(86) /* NodeKind.FOR_AWAIT_OF_STATEMENT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_FOR_AWAIT_OF(), templateNode, alloc_)
		case int32(90) /* NodeKind.IF_STATEMENT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_IF(), templateNode, alloc_)
		case int32(92) /* NodeKind.RETURN_STATEMENT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_RETURN(), templateNode, alloc_)
		case int32(93) /* NodeKind.SWITCH_STATEMENT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_SWITCH(), templateNode, alloc_)
		case int32(140) /* NodeKind.SWITCH_BODY */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_SWITCH_BODY(), templateNode, alloc_)
		case int32(94) /* NodeKind.THROW_STATEMENT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_THROW(), templateNode, alloc_)
		case int32(95) /* NodeKind.TRY_STATEMENT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_TRY(), templateNode, alloc_)
		case int32(96) /* NodeKind.WHILE_STATEMENT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_WHILE(), templateNode, alloc_)
		case int32(98) /* NodeKind.EMPTY */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_EMPTY(), templateNode, alloc_)
		case int32(101) /* NodeKind.WITH */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_WITH(), templateNode, alloc_)
		case int32(99) /* NodeKind.IMPORT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_IMPORT(), templateNode, alloc_)
		case int32(100) /* NodeKind.EXPORT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_EXPORT(), templateNode, alloc_)
		case int32(102) /* NodeKind.VAR_DECLARATION */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_VAR(), templateNode, alloc_)
		case int32(104) /* NodeKind.CONST_DECLARATION */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_CONST(), templateNode, alloc_)
		case int32(103) /* NodeKind.LET_DECLARATION */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_LET(), templateNode, alloc_)
		case int32(105) /* NodeKind.FUNCTION_LITERAL */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_FUNCTION(), templateNode, alloc_)
		case int32(106) /* NodeKind.CLASS_LITERAL */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_CLASS(), templateNode, alloc_)
		case int32(97) /* NodeKind.BLOCK */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_BLOCK(), templateNode, alloc_)
		case int32(91) /* NodeKind.LABELED_STATEMENT */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_LABEL(), templateNode, alloc_)
		case int32(116) /* NodeKind.LABELED_NAME */ :
			return this_.priv_ScriptNodeDeserializer_stringNode_7749(pkg_rhino.StaticGet_Token_LABEL_NAME(), n, templateNode, alloc_)
		case int32(111) /* NodeKind.CLASS_MEMBERS */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_CLASS_MEMBERS(), templateNode, alloc_)
		case int32(112) /* NodeKind.METHOD_DECLARATION */ :
			return this_.priv_ScriptNodeDeserializer_stringNode_7749(pkg_rhino.StaticGet_Token_MEMBER_FUNCTION_DEF(), n, templateNode, alloc_)
		case int32(138) /* NodeKind.FIELD_DECLARATION */ :
			return this_.priv_ScriptNodeDeserializer_stringNode_7749(pkg_rhino.StaticGet_Token_MEMBER_FIELD_DEF(), n, templateNode, alloc_)
		case int32(139) /* NodeKind.COMPUTED_PROP_FIELD */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_COMPUTED_FIELD_DEF(), templateNode, alloc_)
		case int32(113) /* NodeKind.PARAMETER_LIST */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_PARAM_LIST(), templateNode, alloc_)
		case int32(114) /* NodeKind.RENAMABLE_STRING_KEY */ :
			return this_.priv_ScriptNodeDeserializer_stringNode_7749(pkg_rhino.StaticGet_Token_STRING_KEY(), n, templateNode, alloc_)
		case int32(115) /* NodeKind.QUOTED_STRING_KEY */ :
			quotedStringKey = this_.priv_ScriptNodeDeserializer_stringNode_7749(pkg_rhino.StaticGet_Token_STRING_KEY(), n, templateNode, alloc_)
			quotedStringKey.SetQuotedStringKey()
			return quotedStringKey
		case int32(108) /* NodeKind.CASE */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_CASE(), templateNode, alloc_)
		case int32(109) /* NodeKind.DEFAULT_CASE */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_DEFAULT_CASE(), templateNode, alloc_)
		case int32(110) /* NodeKind.CATCH */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_CATCH(), templateNode, alloc_)
		case int32(122) /* NodeKind.SUPER */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_SUPER(), templateNode, alloc_)
		case int32(117) /* NodeKind.ARRAY_PATTERN */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_ARRAY_PATTERN(), templateNode, alloc_)
		case int32(118) /* NodeKind.OBJECT_PATTERN */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_OBJECT_PATTERN(), templateNode, alloc_)
		case int32(119) /* NodeKind.DESTRUCTURING_LHS */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_DESTRUCTURING_LHS(), templateNode, alloc_)
		case int32(123) /* NodeKind.DEFAULT_VALUE */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_DEFAULT_VALUE(), templateNode, alloc_)
		case int32(133) /* NodeKind.RENAMABLE_GETTER_DEF */, int32(134) /* NodeKind.QUOTED_GETTER_DEF */ :
			getterDef = this_.priv_ScriptNodeDeserializer_stringNode_7749(pkg_rhino.StaticGet_Token_GETTER_DEF(), n, templateNode, alloc_)
			if j2gort.RefEquals(n.GetKind(), j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_QUOTED_GETTER_DEF) {
				getterDef.SetQuotedStringKey()
			}
			return getterDef
		case int32(135) /* NodeKind.RENAMABLE_SETTER_DEF */, int32(136) /* NodeKind.QUOTED_SETTER_DEF */ :
			setterDef = this_.priv_ScriptNodeDeserializer_stringNode_7749(pkg_rhino.StaticGet_Token_SETTER_DEF(), n, templateNode, alloc_)
			if j2gort.RefEquals(n.GetKind(), j2gort.Ext_com_google_javascript_jscomp_serialization_NodeKind_F_QUOTED_SETTER_DEF) {
				setterDef.SetQuotedStringKey()
			}
			return setterDef
		case int32(124) /* NodeKind.IMPORT_SPECS */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_IMPORT_SPECS(), templateNode, alloc_)
		case int32(125) /* NodeKind.IMPORT_SPEC */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_IMPORT_SPEC(), templateNode, alloc_)
		case int32(126) /* NodeKind.IMPORT_STAR */ :
			return this_.priv_ScriptNodeDeserializer_stringNode_7749(pkg_rhino.StaticGet_Token_IMPORT_STAR(), n, templateNode, alloc_)
		case int32(127) /* NodeKind.EXPORT_SPECS */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_EXPORT_SPECS(), templateNode, alloc_)
		case int32(128) /* NodeKind.EXPORT_SPEC */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_EXPORT_SPEC(), templateNode, alloc_)
		case int32(137) /* NodeKind.MODULE_BODY */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_MODULE_BODY(), templateNode, alloc_)
		case int32(129) /* NodeKind.ITER_REST */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_ITER_REST(), templateNode, alloc_)
		case int32(130) /* NodeKind.ITER_SPREAD */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_ITER_SPREAD(), templateNode, alloc_)
		case int32(131) /* NodeKind.OBJECT_REST */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_OBJECT_REST(), templateNode, alloc_)
		case int32(132) /* NodeKind.OBJECT_SPREAD */ :
			return this_.priv_ScriptNodeDeserializer_newNodeFromTemplate_7750(pkg_rhino.StaticGet_Token_OBJECT_SPREAD(), templateNode, alloc_)
		case int32(0) /* NodeKind.NODE_KIND_UNSPECIFIED */, int32(141) /* NodeKind.UNRECOGNIZED */ :
		}
	}
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected serialized kind for AstNode: "+j2gort.StrVal(n))))
}

// Java method: ScriptNodeDeserializer.filterOutCastProp(long nodeProperties)
func (this_ *ScriptNodeDeserializer) priv_ScriptNodeDeserializer_filterOutCastProp_7752(nodeProperties int64) int64 {
	_ = nodeProperties
	if this_.colorPoolShard.IsPresent() {
		return nodeProperties
	}
	return (nodeProperties & ^(javalang.Shl64(int64(1), int(j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_COLOR_FROM_CAST.GetNumber()))))
}

// Java class: com.google.javascript.jscomp.serialization.ScriptNodeDeserializer.Runner
type ScriptNodeDeserializer_Runner struct {
	outer_           *ScriptNodeDeserializer
	featureCollector *pkg_parsing.FeatureCollector `gson:"featureCollector"`
	previousLine     int32                         `gson:"previousLine"`
	previousColumn   int32                         `gson:"previousColumn"`
	idHash_          atomic.Int32
	nodeAlloc_       nodestore.NodeAllocator
	cachedTmpl_      pkg_rhino.Node
}

func (this_ *ScriptNodeDeserializer_Runner) BindSelf_(s any) {
	_ = s
}
func (this_ *ScriptNodeDeserializer_Runner) IsJavaNil_() bool { return this_ == nil }
func (*ScriptNodeDeserializer_Runner) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.serialization.ScriptNodeDeserializer$Runner", "Runner"
}
func (this_ *ScriptNodeDeserializer_Runner) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ScriptNodeDeserializer_Runner) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.serialization.ScriptNodeDeserializer.Runner", "com.google.javascript.jscomp.serialization.ScriptNodeDeserializer$Runner", func(v_ any) bool {
		return j2gort.InstanceOf[*ScriptNodeDeserializer_Runner](v_)
	})
}
func (this_ *ScriptNodeDeserializer_Runner) SetOuter_(o *ScriptNodeDeserializer) { this_.outer_ = o }
func (this_ *ScriptNodeDeserializer_Runner) GetOuter_() *ScriptNodeDeserializer  { return this_.outer_ }

func (this_ *ScriptNodeDeserializer_Runner) initFields_ScriptNodeDeserializer_Runner() {
	_ = this_
	this_.featureCollector = pkg_parsing.NewFeatureCollector()

}

// Java constructor: Runner()
func NewScriptNodeDeserializer_Runner(outer_ *ScriptNodeDeserializer) *ScriptNodeDeserializer_Runner {
	this_ := &ScriptNodeDeserializer_Runner{}
	this_.outer_ = outer_
	InitScriptNodeDeserializer_Runner(this_)
	return this_
}

func InitScriptNodeDeserializer_Runner(this_ *ScriptNodeDeserializer_Runner) {
	this_.initFields_ScriptNodeDeserializer_Runner()
}

// Java method: Runner.run()
func (this_ *ScriptNodeDeserializer_Runner) Run() pkg_rhino.Node {
	this_.nodeAlloc_.InitLocal()
	defer this_.nodeAlloc_.Flush()

	var didRet_1 bool
	var retVal_2 pkg_rhino.Node
	func() {
		defer func() {
			if rec_3 := recover(); rec_3 != nil {
				jex_4 := javalang.RecoverException(rec_3)
				_ = jex_4
				if j2gort.InstanceOf[javaio.IIOException](jex_4) {
					ex := j2gort.AsIface[javaio.IIOException](jex_4)
					_ = ex
					panic(NewMalformedTypedAstException_2(j2gort.PtrToAny(this_.priv_Runner_owner_7754().sourceFile), j2gort.AsIface[javalang.Throwable](ex)))
				} else {
					panic(rec_3)
				}
			}
		}()
		var astStream j2gort.Ext_com_google_protobuf_CodedInputStream = this_.priv_Runner_owner_7754().scriptBytes.NewCodedInput()
		_ = astStream
		astStream.SetRecursionLimit(int32(2147483647))
		var contextStack []pkg_parsing.FeatureCollector_FeatureContext = []pkg_parsing.FeatureCollector_FeatureContext{}
		_ = contextStack
		contextStack = append(contextStack, pkg_parsing.StaticGet_FeatureCollector_FeatureContext_NONE())
		var unmodifiableContextStack j2gort.Ext_java_util_List = j2gort.SliceUnmodifiableListLive(&contextStack)
		_ = unmodifiableContextStack
		astRoot_ := j2gort.Ext_com_google_javascript_jscomp_serialization_AstNode_ParseFrom_2_Ext_com_google_protobuf_CodedInputStream_Ext_com_google_protobuf_ExtensionRegistryLite(astStream, j2gort.AsIface[j2gort.Ext_com_google_protobuf_ExtensionRegistryLite](j2gort.Ext_com_google_protobuf_ExtensionRegistry_GetEmptyRegistry()))
		defer astRoot_.Release_()
		var scriptNode pkg_rhino.Node = this_.priv_Runner_visit_7755(astRoot_, &contextStack, unmodifiableContextStack, this_.priv_Runner_owner_7754().priv_ScriptNodeDeserializer_createSourceInfoTemplate_7747(this_.priv_Runner_owner_7754().sourceFile))
		_ = scriptNode
		scriptNode.SetFeatureSet(j2gort.PtrToAny(this_.featureCollector.AllFeatures()))
		retVal_2 = scriptNode
		didRet_1 = true
		return
	}()
	if didRet_1 {
		return retVal_2
	}
	panic("unreachable")
}

// Java method: Runner.owner()
func (this_ *ScriptNodeDeserializer_Runner) priv_Runner_owner_7754() *ScriptNodeDeserializer {
	return this_.outer_
}

// Java method: Runner.visit(com.google.javascript.jscomp.serialization.AstNode astNode, java.util.List<com.google.javascript.jscomp.parsing.FeatureCollector.FeatureContext> contextStack, java.util.List<com.google.javascript.jscomp.parsing.FeatureCollector.FeatureContext> unmodifiableContextStack, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node sourceFileTemplate)
func (this_ *ScriptNodeDeserializer_Runner) priv_Runner_visit_7755(astNode *j2gort.Ext_com_google_javascript_jscomp_serialization_AstNode, contextStack *[]pkg_parsing.FeatureCollector_FeatureContext, unmodifiableContextStack j2gort.Ext_java_util_List, sourceFileTemplate pkg_rhino.Node) pkg_rhino.Node {
	_ = astNode
	_ = contextStack
	_ = unmodifiableContextStack
	_ = sourceFileTemplate
	if (sourceFileTemplate == 0) || (astNode.GetSourceFile() != int32(0)) {
		sourceFileTemplate = this_.priv_Runner_owner_7754().priv_ScriptNodeDeserializer_createSourceInfoTemplate_7747(this_.priv_Runner_owner_7754().filePool[(astNode.GetSourceFile() - int32(1))])
	}
	var currentLine int32 = (this_.previousLine + astNode.GetRelativeLine())
	_ = currentLine
	var currentColumn int32 = (this_.previousColumn + astNode.GetRelativeColumn())
	_ = currentColumn
	if this_.cachedTmpl_ != sourceFileTemplate {
		this_.cachedTmpl_ = sourceFileTemplate
		this_.nodeAlloc_.SetTemplate(nodestore.Node(sourceFileTemplate))
	}

	var n pkg_rhino.Node = this_.priv_Runner_owner_7754().priv_ScriptNodeDeserializer_deserializeSingleNode_7751(astNode, sourceFileTemplate, &this_.nodeAlloc_)
	_ = n
	if astNode.HasType() && this_.priv_Runner_owner_7754().colorPoolShard.IsPresent() {
		n.SetColor(j2gort.CastPtr[ColorPool_ShardView](this_.priv_Runner_owner_7754().colorPoolShard.Get()).GetColor(astNode.GetType()))
	}
	var properties int64 = astNode.GetBooleanProperties()
	_ = properties
	if properties > int64(int32(0)) {
		n.DeserializeProperties(this_.outer_.priv_ScriptNodeDeserializer_filterOutCastProp_7752(properties), sourceFileTemplate.GetIsInClosureUnawareSubtree())
	}
	if astNode.HasJsdoc() {
		n.SetJSDocInfo(JSDocSerializer_DeserializeJsdoc(astNode.GetJsdoc(), this_.outer_.stringPool))
	}
	n.SetLinenoCharno(currentLine, currentColumn)
	this_.previousLine = currentLine
	this_.previousColumn = currentColumn
	if contextStack != nil {
		*contextStack = append(*contextStack, this_.featureCollector.VisitSingleNode(unmodifiableContextStack, n))
	}
	if pkg_rhino.Node_HasBitSet(properties, j2gort.Ext_com_google_javascript_jscomp_serialization_NodeProperty_F_CLOSURE_UNAWARE_SHADOW.GetNumber()) {
		var serializedShadowChild *j2gort.Ext_com_google_javascript_jscomp_serialization_AstNode = astNode.GetChild(int32(0))
		_ = serializedShadowChild
		var closureUnawareSubtreeTemplateNode pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_SCRIPT()).ClonePropsFrom(sourceFileTemplate)
		_ = closureUnawareSubtreeTemplateNode
		closureUnawareSubtreeTemplateNode.SetIsInClosureUnawareSubtree(true)
		var shadowedCode pkg_rhino.Node = this_.priv_Runner_visit_7755(serializedShadowChild, contextStack, unmodifiableContextStack, closureUnawareSubtreeTemplateNode)
		_ = shadowedCode
		this_.priv_Runner_owner_7754().priv_ScriptNodeDeserializer_setOriginalNameIfPresent_7748(serializedShadowChild, shadowedCode, &this_.nodeAlloc_)
		var shadowRoot pkg_rhino.Node = pkg_rhino.IR_Root([]pkg_rhino.Node{pkg_rhino.IR_Script_1_Node([]pkg_rhino.Node{pkg_rhino.IR_ExprResult(shadowedCode)})})
		_ = shadowRoot
		shadowRoot.ClonePropsFrom(closureUnawareSubtreeTemplateNode)
		shadowRoot.GetFirstChild().ClonePropsFrom(closureUnawareSubtreeTemplateNode)
		shadowRoot.GetFirstFirstChild().ClonePropsFrom(closureUnawareSubtreeTemplateNode)
		n.SetClosureUnawareShadow(shadowRoot)
		return n
	}
	var children int32 = astNode.GetChildCount()
	_ = children
	{
		var i int32 = int32(0)
		_ = i
		for ; i < children; i = i + 1 {
			var child *j2gort.Ext_com_google_javascript_jscomp_serialization_AstNode = astNode.GetChild(i)
			_ = child
			var deserializedChild pkg_rhino.Node = this_.priv_Runner_visit_7755(child, contextStack, unmodifiableContextStack, sourceFileTemplate)
			_ = deserializedChild
			n.AddChildToBack(deserializedChild)
			this_.priv_Runner_owner_7754().priv_ScriptNodeDeserializer_setOriginalNameIfPresent_7748(child, deserializedChild, &this_.nodeAlloc_)
		}
	}
	_ = j2gort.SliceGetLast(*contextStack)
	*contextStack = (*contextStack)[:len(*contextStack)-1]
	return n
}
