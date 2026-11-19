/*
 * Copyright 2012 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ReplaceMessagesForChrome.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.ReplaceMessagesForChrome
type ReplaceMessagesForChrome struct {
	JsMessageVisitorBase
	astFactory             *AstFactory `gson:"astFactory"`
	introducesRegexpSyntax bool        `gson:"introducesRegexpSyntax"`
}

func (this_ *ReplaceMessagesForChrome) BindSelf_(s any) {
	_ = s
	this_.JsMessageVisitorBase.BindSelf_(s)
}
func (this_ *ReplaceMessagesForChrome) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceMessagesForChrome) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceMessagesForChrome", "ReplaceMessagesForChrome"
}
func (this_ *ReplaceMessagesForChrome) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.JsMessageVisitorBase.IdentityHash_()
}
func (*ReplaceMessagesForChrome) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceMessagesForChrome", "com.google.javascript.jscomp.ReplaceMessagesForChrome", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceMessagesForChrome](v_)
	})
}
func (*ReplaceMessagesForChrome) Implements_jscomp_JsMessageVisitor()                        {}
func (*ReplaceMessagesForChrome) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*ReplaceMessagesForChrome) Implements_jscomp_CompilerPass()                            {}
func (*ReplaceMessagesForChrome) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_ReplaceMessagesForChrome_DECLARE_ICU_TEMPLATE_NOT_SUPPORTED *DiagnosticType

func StaticGet_ReplaceMessagesForChrome_DECLARE_ICU_TEMPLATE_NOT_SUPPORTED() *DiagnosticType {
	clinit_ReplaceMessagesForChrome()
	return static_ReplaceMessagesForChrome_DECLARE_ICU_TEMPLATE_NOT_SUPPORTED
}
func StaticSet_ReplaceMessagesForChrome_DECLARE_ICU_TEMPLATE_NOT_SUPPORTED(v *DiagnosticType) {
	clinit_ReplaceMessagesForChrome()
	static_ReplaceMessagesForChrome_DECLARE_ICU_TEMPLATE_NOT_SUPPORTED = v
}

var clinitDone_ReplaceMessagesForChrome uint32
var clinitState_ReplaceMessagesForChrome j2gort.ClinitState

func clinit_ReplaceMessagesForChrome() {
	if atomic.LoadUint32(&clinitDone_ReplaceMessagesForChrome) != 0 {
		return
	}
	clinitSlow_ReplaceMessagesForChrome()
}

func clinitSlow_ReplaceMessagesForChrome() {
	if !j2gort.EnterClinit(&clinitState_ReplaceMessagesForChrome, &clinitDone_ReplaceMessagesForChrome) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ReplaceMessagesForChrome, &clinitDone_ReplaceMessagesForChrome)
	static_ReplaceMessagesForChrome_DECLARE_ICU_TEMPLATE_NOT_SUPPORTED = DiagnosticType_Error("JSC_DECLARE_ICU_TEMPLATE", "goog.i18n.messages.declareIcuTemplate() is not supported for Chrome i18n.")
}

func (this_ *ReplaceMessagesForChrome) initFields_ReplaceMessagesForChrome() {
	_ = this_

}

// Java constructor: ReplaceMessagesForChrome(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.JsMessage.IdGenerator idGenerator)
func NewReplaceMessagesForChrome(compiler AbstractCompiler, idGenerator JsMessage_IdGenerator) *ReplaceMessagesForChrome {
	clinit_ReplaceMessagesForChrome()
	this_ := &ReplaceMessagesForChrome{}
	this_.BindSelf_(this_)
	InitReplaceMessagesForChrome(this_, compiler, idGenerator)
	return this_
}

func InitReplaceMessagesForChrome(this_ *ReplaceMessagesForChrome, compiler AbstractCompiler, idGenerator JsMessage_IdGenerator) {
	_ = compiler
	_ = idGenerator
	InitJsMessageVisitor(&this_.JsMessageVisitorBase, compiler, idGenerator)
	this_.initFields_ReplaceMessagesForChrome()
	this_.astFactory = compiler.CreateAstFactory()
}

// Java method: ReplaceMessagesForChrome.getChromeI18nGetMessageNode(java.lang.String messageId)
func (this_ *ReplaceMessagesForChrome) priv_ReplaceMessagesForChrome_getChromeI18nGetMessageNode_4876(messageId string) pkg_rhino.Node {
	_ = messageId
	var getMessage pkg_rhino.Node = this_.astFactory.CreateQNameWithUnknownType_1_string("chrome.i18n.getMessage")
	_ = getMessage
	return this_.astFactory.CreateCall(getMessage, AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_STRING()), []pkg_rhino.Node{this_.astFactory.CreateString(messageId)})
}

// Java method: ReplaceMessagesForChrome.processIcuTemplateDefinition(com.google.javascript.jscomp.IcuTemplateDefinition definition)
func (this_ *ReplaceMessagesForChrome) ProcessIcuTemplateDefinition(definition IcuTemplateDefinition) {
	_ = definition
	this_.GetJsMessageVisitorBase().compiler.Report_1(JSError_Make_3(definition.GetMessageNode(), StaticGet_ReplaceMessagesForChrome_DECLARE_ICU_TEMPLATE_NOT_SUPPORTED(), []string{}))
}

// Java method: ReplaceMessagesForChrome.processJsMessageDefinition(com.google.javascript.jscomp.JsMessageDefinition definition)
func (this_ *ReplaceMessagesForChrome) ProcessJsMessageDefinition(definition JsMessageDefinition) {
	_ = definition
	this_.introducesRegexpSyntax = false
	var newValue pkg_rhino.Node = this_.priv_ReplaceMessagesForChrome_getNewValueNode_4877(definition.GetMessage(), definition)
	_ = newValue
	definition.GetMessageNode().ReplaceWith(newValue)
	if this_.introducesRegexpSyntax {
		NodeUtil_AddFeatureToScript(NodeUtil_GetEnclosingScript(newValue), pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_SYNTAX(), this_.GetJsMessageVisitorBase().compiler)
		this_.introducesRegexpSyntax = false
	}
	this_.GetJsMessageVisitorBase().compiler.ReportChangeToEnclosingScope(newValue)
}

// Java method: ReplaceMessagesForChrome.getNewValueNode(com.google.javascript.jscomp.JsMessage message, com.google.javascript.jscomp.JsMessageDefinition definition)
func (this_ *ReplaceMessagesForChrome) priv_ReplaceMessagesForChrome_getNewValueNode_4877(message *JsMessage, definition JsMessageDefinition) pkg_rhino.Node {
	_ = message
	_ = definition
	var newValueNode pkg_rhino.Node = this_.priv_ReplaceMessagesForChrome_getChromeI18nGetMessageNode_4876(message.GetId())
	_ = newValueNode
	var isHtml bool = definition.ShouldEscapeLessThan()
	_ = isHtml
	if !(message.JsPlaceholderNames().IsEmpty()) || isHtml {
		var placeholderNames j2gort.Ext_java_util_List = j2gort.Ext_com_google_common_collect_Ordering_Natural().SortedCopy(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](message.JsPlaceholderNames()))
		_ = placeholderNames
		var placeholderValueArray pkg_rhino.Node = this_.astFactory.CreateArraylit_1_Node([]pkg_rhino.Node{})
		_ = placeholderValueArray
		var placeholderValueMap j2gort.Ext_com_google_common_collect_ImmutableMap = definition.GetPlaceholderValueMap()
		_ = placeholderValueMap
		{
			cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(placeholderNames))
			for cur_1.HasNext() {
				name := cur_1.NextStr()
				_ = name
				var value pkg_rhino.Node = j2gort.CheckNotNullNode(j2gort.CastNode[pkg_rhino.Node](j2gort.MapGetStr(placeholderValueMap, name))).CloneTree_0()
				_ = value
				placeholderValueArray.AddChildToBack(value)
			}
		}
		if isHtml {
			var args pkg_rhino.Node = this_.astFactory.CreateArraylit_1_Node([]pkg_rhino.Node{this_.astFactory.CreateString(message.GetId()), placeholderValueArray})
			_ = args
			var options pkg_rhino.Node = this_.astFactory.CreateArraylit_1_Node([]pkg_rhino.Node{this_.astFactory.CreateObjectLit_1([]pkg_rhino.Node{this_.astFactory.CreateStringKey("escapeLt", this_.astFactory.CreateBoolean(true))})})
			_ = options
			var regexp pkg_rhino.Node = pkg_rhino.IR_Regexp_1(this_.astFactory.CreateString("Chrome\\/(\\d+)")).SetColor(pkg_colors.StaticGet_StandardColors_UNKNOWN())
			_ = regexp
			var userAgent pkg_rhino.Node = this_.astFactory.CreateQNameWithUnknownType_1_string("navigator.userAgent")
			_ = userAgent
			var version pkg_rhino.Node = this_.astFactory.CreateGetElem(this_.astFactory.CreateOr(this_.astFactory.CreateCall(this_.astFactory.CreateGetPropWithUnknownType(regexp, "exec"), AstFactory_S_Type_1_ColorId(pkg_colors.StaticGet_StandardColors_ARRAY_ID()), []pkg_rhino.Node{userAgent}), this_.astFactory.CreateArraylit_1_Node([]pkg_rhino.Node{})), this_.astFactory.CreateNumber(float64(int32(1))))
			_ = version
			var condition pkg_rhino.Node = pkg_rhino.IR_Ge(version, this_.astFactory.CreateNumber(float64(int32(79)))).SetColor(pkg_colors.StaticGet_StandardColors_BOOLEAN())
			_ = condition
			args = this_.astFactory.CreateCall(this_.astFactory.CreateGetPropWithUnknownType(args, "concat"), AstFactory_S_Type_1_ColorId(pkg_colors.StaticGet_StandardColors_ARRAY_ID()), []pkg_rhino.Node{this_.astFactory.CreateHook(condition, options, this_.astFactory.CreateArraylit_1_Node([]pkg_rhino.Node{}))})
			newValueNode = this_.astFactory.CreateCallWithUnknownType(this_.astFactory.CreateQNameWithUnknownType_1_string("chrome.i18n.getMessage.apply"), []pkg_rhino.Node{this_.astFactory.CreateNull(), args})
			this_.introducesRegexpSyntax = true
		} else {
			newValueNode.AddChildToBack(placeholderValueArray)
		}
	}
	newValueNode.SrcrefTreeIfMissing(definition.GetMessageNode())
	return newValueNode
}
func (*ReplaceMessagesForChrome) AlwaysShouldTraverse_() {}
