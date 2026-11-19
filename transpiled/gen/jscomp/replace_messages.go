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
// Port of com/google/javascript/jscomp/ReplaceMessages.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.ReplaceMessages
type ReplaceMessages struct {
	compiler                 AbstractCompiler `gson:"compiler"`
	bundle                   MessageBundle    `gson:"bundle"`
	astFactory               *AstFactory      `gson:"astFactory"`
	idHash_                  atomic.Int32
	strictReplacement        bool `gson:"strictReplacement"`
	collapsePropertiesHasRun bool `gson:"collapsePropertiesHasRun"`
}

func (this_ *ReplaceMessages) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceMessages) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceMessages) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceMessages", "ReplaceMessages"
}
func (this_ *ReplaceMessages) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReplaceMessages) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceMessages", "com.google.javascript.jscomp.ReplaceMessages", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceMessages](v_)
	})
}

var static_ReplaceMessages_BUNDLE_DOES_NOT_HAVE_THE_MESSAGE *DiagnosticType

func StaticGet_ReplaceMessages_BUNDLE_DOES_NOT_HAVE_THE_MESSAGE() *DiagnosticType {
	clinit_ReplaceMessages()
	return static_ReplaceMessages_BUNDLE_DOES_NOT_HAVE_THE_MESSAGE
}
func StaticSet_ReplaceMessages_BUNDLE_DOES_NOT_HAVE_THE_MESSAGE(v *DiagnosticType) {
	clinit_ReplaceMessages()
	static_ReplaceMessages_BUNDLE_DOES_NOT_HAVE_THE_MESSAGE = v
}

var static_ReplaceMessages_INVALID_ALTERNATE_MESSAGE_PLACEHOLDERS *DiagnosticType

func StaticGet_ReplaceMessages_INVALID_ALTERNATE_MESSAGE_PLACEHOLDERS() *DiagnosticType {
	clinit_ReplaceMessages()
	return static_ReplaceMessages_INVALID_ALTERNATE_MESSAGE_PLACEHOLDERS
}
func StaticSet_ReplaceMessages_INVALID_ALTERNATE_MESSAGE_PLACEHOLDERS(v *DiagnosticType) {
	clinit_ReplaceMessages()
	static_ReplaceMessages_INVALID_ALTERNATE_MESSAGE_PLACEHOLDERS = v
}

var static_ReplaceMessages_ICU_MSG_OPTIONS *ReplaceMessages_MsgOptions

func StaticGet_ReplaceMessages_ICU_MSG_OPTIONS() *ReplaceMessages_MsgOptions {
	clinit_ReplaceMessages()
	return static_ReplaceMessages_ICU_MSG_OPTIONS
}
func StaticSet_ReplaceMessages_ICU_MSG_OPTIONS(v *ReplaceMessages_MsgOptions) {
	clinit_ReplaceMessages()
	static_ReplaceMessages_ICU_MSG_OPTIONS = v
}

var clinitDone_ReplaceMessages uint32
var clinitState_ReplaceMessages j2gort.ClinitState

func clinit_ReplaceMessages() {
	if atomic.LoadUint32(&clinitDone_ReplaceMessages) != 0 {
		return
	}
	clinitSlow_ReplaceMessages()
}

func clinitSlow_ReplaceMessages() {
	if !j2gort.EnterClinit(&clinitState_ReplaceMessages, &clinitDone_ReplaceMessages) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ReplaceMessages, &clinitDone_ReplaceMessages)
	static_ReplaceMessages_BUNDLE_DOES_NOT_HAVE_THE_MESSAGE = DiagnosticType_Error("JSC_BUNDLE_DOES_NOT_HAVE_THE_MESSAGE", "Message with id = {0} could not be found in replacement bundle")
	static_ReplaceMessages_INVALID_ALTERNATE_MESSAGE_PLACEHOLDERS = DiagnosticType_Error("JSC_INVALID_ALTERNATE_MESSAGE_PLACEHOLDERS", "Alternate message ID={0} placeholders ({1}) differs from {2} placeholders ({3}).")
	static_ReplaceMessages_ICU_MSG_OPTIONS = NewReplaceMessages_MsgOptions(true, false, false)
}

func (this_ *ReplaceMessages) initFields_ReplaceMessages() {
	_ = this_
}

// Java constructor: ReplaceMessages(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.MessageBundle bundle, boolean strictReplacement)
func NewReplaceMessages(compiler AbstractCompiler, bundle MessageBundle, strictReplacement bool) *ReplaceMessages {
	clinit_ReplaceMessages()
	this_ := &ReplaceMessages{}
	InitReplaceMessages(this_, compiler, bundle, strictReplacement)
	return this_
}

func InitReplaceMessages(this_ *ReplaceMessages, compiler AbstractCompiler, bundle MessageBundle, strictReplacement bool) {
	_ = compiler
	_ = bundle
	_ = strictReplacement
	this_.initFields_ReplaceMessages()
	this_.compiler = compiler
	this_.astFactory = compiler.CreateAstFactory()
	this_.bundle = bundle
	this_.strictReplacement = strictReplacement
	this_.collapsePropertiesHasRun = ((compiler.GetOptions().GetPropertyCollapseLevel() == StaticGet_CompilerOptions_PropertyCollapseLevel_ALL()) && compiler.GetOptions().DoLateLocalization())
}

// Java method: ReplaceMessages.getMsgProtectionPass()
func (this_ *ReplaceMessages) GetMsgProtectionPass() CompilerPass {
	return NewReplaceMessages_MsgProtectionPass(this_)
}

// Java method: ReplaceMessages.getMsgOptionsFromDefinition(com.google.javascript.jscomp.JsMessageDefinition definition)
func (this_ *ReplaceMessages) priv_ReplaceMessages_getMsgOptionsFromDefinition_4857(definition JsMessageDefinition) *ReplaceMessages_MsgOptions {
	_ = definition
	return NewReplaceMessages_MsgOptions(false, definition.ShouldEscapeLessThan(), definition.ShouldUnescapeHtmlEntities())
}

// Java method: ReplaceMessages.createMsgPropertiesNode(com.google.javascript.jscomp.JsMessage message, com.google.javascript.jscomp.ReplaceMessages.MsgOptions msgOptions)
func (this_ *ReplaceMessages) priv_ReplaceMessages_createMsgPropertiesNode_4858(message *JsMessage, msgOptions *ReplaceMessages_MsgOptions) pkg_rhino.Node {
	_ = message
	_ = msgOptions
	var val_1 ReplaceMessages_QuotedKeyObjectLitBuilder
	val_1.outer_ = this_
	InitReplaceMessages_QuotedKeyObjectLitBuilder(&val_1)
	var msgPropsBuilder *ReplaceMessages_QuotedKeyObjectLitBuilder = &val_1
	_ = msgPropsBuilder
	msgPropsBuilder.priv_QuotedKeyObjectLitBuilder_addString_4876("key", message.GetKey())
	if msgOptions.IsIcuTemplate() && !(message.CanonicalPlaceholderNames().IsEmpty()) {
		var namesArrayLit pkg_rhino.Node = this_.astFactory.CreateArraylit_1_Node([]pkg_rhino.Node{})
		_ = namesArrayLit
		{
			cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(message.CanonicalPlaceholderNames()))
			for cur_2.HasNext() {
				name := cur_2.NextStr()
				_ = name
				namesArrayLit.AddChildToBack(this_.astFactory.CreateString(name))
			}
		}
		msgPropsBuilder.priv_QuotedKeyObjectLitBuilder_addNode_4877("icu_placeholder_names", namesArrayLit)
	}
	var altId string = message.GetAlternateId()
	_ = altId
	if altId != j2gort.NullString {
		msgPropsBuilder.priv_QuotedKeyObjectLitBuilder_addString_4876("alt_id", altId)
	}
	var meaning string = message.GetMeaning()
	_ = meaning
	if meaning != j2gort.NullString {
		msgPropsBuilder.priv_QuotedKeyObjectLitBuilder_addString_4876("meaning", meaning)
	}
	if msgOptions.IsIcuTemplate() {
		msgPropsBuilder.priv_QuotedKeyObjectLitBuilder_addString_4876("msg_text", message.AsIcuMessageString_0())
		msgPropsBuilder.priv_QuotedKeyObjectLitBuilder_addString_4876("isIcuTemplate", "")
	} else {
		msgPropsBuilder.priv_QuotedKeyObjectLitBuilder_addString_4876("msg_text", message.AsJsMessageString_0())
	}
	if msgOptions.EscapeLessThan() {
		msgPropsBuilder.priv_QuotedKeyObjectLitBuilder_addString_4876("escapeLessThan", "")
	}
	if msgOptions.UnescapeHtmlEntities() {
		msgPropsBuilder.priv_QuotedKeyObjectLitBuilder_addString_4876("unescapeHtmlEntities", "")
	}
	return msgPropsBuilder.priv_QuotedKeyObjectLitBuilder_build_4878()
}

// Java method: ReplaceMessages.getReplacementCompletionPass()
func (this_ *ReplaceMessages) GetReplacementCompletionPass() CompilerPass {
	return NewReplaceMessages_ReplacementCompletionPass(this_)
}

// Java static method: ReplaceMessages.extractPlaceholderValuesMapOrThrow(com.google.javascript.rhino.Node valuesObjLit)
func ReplaceMessages_ExtractPlaceholderValuesMapOrThrow(valuesObjLit pkg_rhino.Node) j2gort.Ext_com_google_common_collect_ImmutableMap {
	clinit_ReplaceMessages()
	_ = valuesObjLit
	var didRet_3 bool
	var retVal_4 j2gort.Ext_com_google_common_collect_ImmutableMap
	func() {
		defer func() {
			if rec_5 := recover(); rec_5 != nil {
				jex_6 := javalang.RecoverException(rec_5)
				_ = jex_6
				if j2gort.InstanceOf[*JsMessageVisitor_MalformedException](jex_6) {
					e := j2gort.AsIface[*JsMessageVisitor_MalformedException](jex_6)
					_ = e
					panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.PtrToAny(e))))
				} else {
					panic(rec_5)
				}
			}
		}()
		retVal_4 = JsMessageVisitor_ExtractObjectLiteralMap(valuesObjLit).ExtractAsValueMap()
		didRet_3 = true
		return
	}()
	if didRet_3 {
		return retVal_4
	}
	panic("unreachable")
}

// Java method: ReplaceMessages.lookupMessage(com.google.javascript.rhino.Node callNode, com.google.javascript.jscomp.MessageBundle bundle, com.google.javascript.jscomp.JsMessage message)
func (this_ *ReplaceMessages) priv_ReplaceMessages_lookupMessage_4860(callNode pkg_rhino.Node, bundle MessageBundle, message *JsMessage) *JsMessage {
	_ = callNode
	_ = bundle
	_ = message
	var translatedMessage *JsMessage = bundle.GetMessage(message.GetId())
	_ = translatedMessage
	if translatedMessage != nil {
		return translatedMessage
	}
	var alternateId string = message.GetAlternateId()
	_ = alternateId
	if alternateId == j2gort.NullString {
		return nil
	}
	var alternateMessage *JsMessage = bundle.GetMessage(alternateId)
	_ = alternateMessage
	if alternateMessage != nil {
		var jsCodePlaceholderNames j2gort.Ext_com_google_common_collect_ImmutableSet = message.JsPlaceholderNames()
		_ = jsCodePlaceholderNames
		var alternateMsgPlaceholderNames j2gort.Ext_com_google_common_collect_ImmutableSet = alternateMessage.JsPlaceholderNames()
		_ = alternateMsgPlaceholderNames
		if !(javalang.ObjectsEquals(jsCodePlaceholderNames, alternateMsgPlaceholderNames)) {
			if jsCodePlaceholderNames.IsEmpty() && ReplaceMessages_IsStartOfIcuMessage(message.AsIcuMessageString_0()) {
				return alternateMessage
			} else {
				this_.compiler.Report_1(JSError_Make_3(callNode, StaticGet_ReplaceMessages_INVALID_ALTERNATE_MESSAGE_PLACEHOLDERS(), []string{alternateId, j2gort.StrVal(alternateMsgPlaceholderNames), message.GetKey(), j2gort.StrVal(jsCodePlaceholderNames)}))
				return nil
			}
		}
	}
	return alternateMessage
}

// Java method: ReplaceMessages.getFullReplacementPass()
func (this_ *ReplaceMessages) GetFullReplacementPass() CompilerPass {
	return NewReplaceMessages_FullReplacementPass(this_)
}

// Java method: ReplaceMessages.createNodeForGenderedMsgString(com.google.javascript.jscomp.JsMessage message, java.util.Map<java.lang.String,com.google.javascript.rhino.Node> placeholderMap, com.google.javascript.jscomp.ReplaceMessages.MsgOptions options, com.google.javascript.rhino.Node nodeToReplace)
func (this_ *ReplaceMessages) priv_ReplaceMessages_createNodeForGenderedMsgString_4862(message *JsMessage, placeholderMap j2gort.Ext_java_util_Map, options *ReplaceMessages_MsgOptions, nodeToReplace pkg_rhino.Node) pkg_rhino.Node {
	_ = message
	_ = placeholderMap
	_ = options
	_ = nodeToReplace
	var paramList pkg_rhino.Node = this_.astFactory.CreateParamList([]string{})
	_ = paramList
	var type_ AstFactory_Type = AstFactory_S_Type_1_Node(nodeToReplace)
	_ = type_
	var uniqueId string = this_.compiler.GetUniqueIdSupplier().GetUniqueId(this_.compiler.GetInput(NodeUtil_GetInputId(nodeToReplace)))
	_ = uniqueId
	var placeholderMapIds j2gort.OrderedMap[string, string] = j2gort.NewOrderedMap[string, string]()
	_ = placeholderMapIds
	{
		cur_7 := j2gort.NewMapKeyStrCursor(placeholderMap)
		for cur_7.HasNext() {
			placeholderName := cur_7.NextStr()
			_ = placeholderName
			var placeholderId string = (j2gort.StrValString(placeholderName) + j2gort.StrValString(uniqueId))
			_ = placeholderId
			paramList.AddChildToBack(this_.astFactory.CreateName_2_string_AstFactory_Type(placeholderId, type_))
			placeholderMapIds.Put(placeholderName, placeholderId)
		}
	}
	var hookExpression pkg_rhino.Node = this_.priv_ReplaceMessages_createHookExpressionForGenderedMsg_4864(type_, message, options, nodeToReplace, &placeholderMapIds)
	_ = hookExpression
	if placeholderMap.IsEmpty() {
		return hookExpression
	}
	var function pkg_rhino.Node = this_.astFactory.CreateFunction("", paramList, pkg_rhino.IR_Block_1_Node(pkg_rhino.IR_ReturnNode_1(hookExpression)), type_)
	_ = function
	function.SetColor(nodeToReplace.GetColor())
	this_.compiler.ReportChangeToChangeScope(function)
	var callNode pkg_rhino.Node = this_.astFactory.CreateCall(function, type_, []pkg_rhino.Node{})
	_ = callNode
	{
		cur_8 := j2gort.NewMapValCursor(placeholderMap)
		for cur_8.HasNext() {
			rawElem_9 := cur_8.Next()
			node := j2gort.CastNode[pkg_rhino.Node](rawElem_9)
			_ = node
			callNode.AddChildToBack(node.CloneTree_0())
		}
	}
	return callNode
}

// Java method: ReplaceMessages.createConditionForGrammaticalGender(com.google.javascript.jscomp.AstFactory.Type type, com.google.javascript.jscomp.JsMessage.GrammaticalGenderCase grammaticalGender, com.google.javascript.rhino.Node nodeToReplace)
func (this_ *ReplaceMessages) priv_ReplaceMessages_createConditionForGrammaticalGender_4863(type_ AstFactory_Type, grammaticalGender JsMessage_GrammaticalGenderCase, nodeToReplace pkg_rhino.Node) pkg_rhino.Node {
	_ = type_
	_ = grammaticalGender
	_ = nodeToReplace
	if this_.collapsePropertiesHasRun {
		var referenceName string = j2gort.NullString
		_ = referenceName
		switch (grammaticalGender).Ordinal() {
		case int32(0) /* GrammaticalGenderCase.MASCULINE */ :
			referenceName = "goog$msgKind$MASCULINE"
		case int32(1) /* GrammaticalGenderCase.FEMININE */ :
			referenceName = "goog$msgKind$FEMININE"
		case int32(2) /* GrammaticalGenderCase.NEUTER */ :
			referenceName = "goog$msgKind$NEUTER"
		default:
			referenceName = "goog$msgKind$OTHER"
		}
		var result pkg_rhino.Node = this_.astFactory.CreateName_2_string_AstFactory_Type(referenceName, AstFactory_S_Type_1_Node(nodeToReplace))
		_ = result
		result.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
		return result
	} else {
		var googNode pkg_rhino.Node = this_.astFactory.CreateName_2_string_AstFactory_Type("goog", type_)
		_ = googNode
		googNode.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
		return this_.astFactory.CreateGetProp(this_.astFactory.CreateGetProp(googNode, "msgKind", AstFactory_S_Type_1_Node(nodeToReplace)), j2gort.CheckNotNullEnum(grammaticalGender).String(), AstFactory_S_Type_1_Node(nodeToReplace))
	}
}

// Java method: ReplaceMessages.createHookExpressionForGenderedMsg(com.google.javascript.jscomp.AstFactory.Type type, com.google.javascript.jscomp.JsMessage msg, com.google.javascript.jscomp.ReplaceMessages.MsgOptions options, com.google.javascript.rhino.Node nodeToReplace, java.util.Map<java.lang.String,java.lang.String> placeholderMapIds)
func (this_ *ReplaceMessages) priv_ReplaceMessages_createHookExpressionForGenderedMsg_4864(type_ AstFactory_Type, msg *JsMessage, options *ReplaceMessages_MsgOptions, nodeToReplace pkg_rhino.Node, placeholderMapIds *j2gort.OrderedMap[string, string]) pkg_rhino.Node {
	_ = type_
	_ = msg
	_ = options
	_ = nodeToReplace
	_ = placeholderMapIds
	var hookHead pkg_rhino.Node = pkg_rhino.IR_Hook(pkg_rhino.IR_Name(""), pkg_rhino.IR_String(""), pkg_rhino.IR_String(""))
	_ = hookHead
	hookHead.SetColor(nodeToReplace.GetColor())
	var placeholderHook pkg_rhino.Node = hookHead
	_ = placeholderHook
	var placeholderNodeProvider ReplaceMessages_PlaceholderNodeProvider = NewLambda_ReplaceMessages_PlaceholderNodeProvider(func(placeholderName string) pkg_rhino.Node {
		_ = placeholderName
		var placeholderId string = j2gort.OrderedMapPtrGetStr(placeholderMapIds, placeholderName)
		_ = placeholderId
		if placeholderId != j2gort.NullString {
			return this_.astFactory.CreateName_2_string_AstFactory_Type(placeholderId, AstFactory_S_Type_1_Node(nodeToReplace))
		} else {
			return pkg_rhino.Node(0)
		}
	})
	_ = placeholderNodeProvider
	{
		cur_10 := j2gort.NewIterCursor(j2gort.CheckNotNull(msg.GetGenderedMessageVariants()))
		for cur_10.HasNext() {
			rawElem_11 := cur_10.Next()
			grammaticalGender := j2gort.Cast[JsMessage_GrammaticalGenderCase](rawElem_11)
			_ = grammaticalGender
			{
				var iifeCond_1 bool
				{

					re_12 := grammaticalGender
					ae_13 := StaticGet_JsMessage_GrammaticalGenderCase_OTHER()
					if re_12 == 0 {
						j2gort.ThrowNPE()
					}
					iifeCond_1 = re_12 == ae_13
				}
				if iifeCond_1 {
					continue
				}
			}
			var condition pkg_rhino.Node = this_.priv_ReplaceMessages_createConditionForGrammaticalGender_4863(type_, grammaticalGender, nodeToReplace)
			_ = condition
			var currentHook pkg_rhino.Node = pkg_rhino.IR_Hook(condition, this_.priv_ReplaceMessages_buildMessageNode_4865(j2gort.AsIface[j2gort.Ext_java_util_List](msg.GetGenderedMessageParts(grammaticalGender)), options, nodeToReplace, placeholderNodeProvider, true), this_.astFactory.CreateString(""))
			_ = currentHook
			currentHook.SetColor(nodeToReplace.GetColor())
			if placeholderHook.GetParent() != 0 {
				var parent pkg_rhino.Node = placeholderHook.GetParent()
				_ = parent
				parent.GetLastChild().ReplaceWith(currentHook)
				placeholderHook = currentHook.GetLastChild()
			} else {
				hookHead = currentHook
				placeholderHook = currentHook.GetLastChild()
			}
		}
	}
	placeholderHook.GetParent().GetLastChild().ReplaceWith(this_.priv_ReplaceMessages_buildMessageNode_4865(j2gort.AsIface[j2gort.Ext_java_util_List](msg.GetGenderedMessageParts(StaticGet_JsMessage_GrammaticalGenderCase_OTHER())), options, nodeToReplace, placeholderNodeProvider, true))
	return hookHead
}

// Java method: ReplaceMessages.buildMessageNode(java.util.List<com.google.javascript.jscomp.JsMessage.Part> msgParts, com.google.javascript.jscomp.ReplaceMessages.MsgOptions options, com.google.javascript.rhino.Node nodeToReplace, com.google.javascript.jscomp.ReplaceMessages.PlaceholderNodeProvider placeholderNodeProvider, boolean isGenderedMsg)
func (this_ *ReplaceMessages) priv_ReplaceMessages_buildMessageNode_4865(msgParts j2gort.Ext_java_util_List, options *ReplaceMessages_MsgOptions, nodeToReplace pkg_rhino.Node, placeholderNodeProvider ReplaceMessages_PlaceholderNodeProvider, isGenderedMsg bool) pkg_rhino.Node {
	_ = msgParts
	_ = options
	_ = nodeToReplace
	_ = placeholderNodeProvider
	_ = isGenderedMsg
	var parts j2gort.Ext_java_util_List = ReplaceMessages_MergeStringParts(msgParts)
	_ = parts
	if parts.IsEmpty() {
		return this_.astFactory.CreateString("")
	}
	var message pkg_rhino.Node = pkg_rhino.Node(0)
	_ = message
	{
		cur_14 := j2gort.NewIterCursor(j2gort.CheckNotNull(parts))
		for cur_14.HasNext() {
			rawElem_15 := cur_14.Next()
			msgPart := j2gort.AsIface[JsMessage_Part](rawElem_15)
			_ = msgPart
			var partNode pkg_rhino.Node = pkg_rhino.Node(0)
			_ = partNode
			if msgPart.IsPlaceholder() {
				var jsPlaceholderName string = msgPart.GetJsPlaceholderName()
				_ = jsPlaceholderName
				var valueNode pkg_rhino.Node = placeholderNodeProvider.Get(jsPlaceholderName)
				_ = valueNode
				if valueNode == 0 {
					if isGenderedMsg {
						partNode = this_.astFactory.CreateString(("{" + j2gort.StrValString(msgPart.GetCanonicalPlaceholderName()) + "}"))
					} else {
						panic(NewJsMessageVisitor_MalformedException(("Unrecognized message placeholder referenced: " + j2gort.StrValString(jsPlaceholderName)), nodeToReplace))
					}
				} else {
					partNode = valueNode.CloneTree_0()
				}
			} else {
				partNode = this_.priv_ReplaceMessages_createNodeForMsgString_4868(options, msgPart.GetString())
			}
			if message == 0 {
				message = partNode
			} else if partNode.IsString() && message.IsString() {
				message = this_.astFactory.CreateString((j2gort.StrValString(message.GetString()) + j2gort.StrValString(partNode.GetString())))
			} else {
				message = this_.astFactory.CreateAdd(message, partNode)
			}
		}
	}
	return message
}

// Java method: ReplaceMessages.constructStringExprNode(com.google.javascript.jscomp.JsMessage message, java.util.Map<java.lang.String,com.google.javascript.rhino.Node> placeholderMap, com.google.javascript.jscomp.ReplaceMessages.MsgOptions options, com.google.javascript.rhino.Node nodeToReplace)
func (this_ *ReplaceMessages) priv_ReplaceMessages_constructStringExprNode_4866(message *JsMessage, placeholderMap j2gort.Ext_java_util_Map, options *ReplaceMessages_MsgOptions, nodeToReplace pkg_rhino.Node) pkg_rhino.Node {
	_ = message
	_ = placeholderMap
	_ = options
	_ = nodeToReplace
	if !(message.GetGenderedMessageVariants().IsEmpty()) {
		return this_.priv_ReplaceMessages_createNodeForGenderedMsgString_4862(message, placeholderMap, options, nodeToReplace)
	} else {
		return this_.priv_ReplaceMessages_constructNonGenderedMsgNode_4867(message, placeholderMap, options, nodeToReplace)
	}
}

// Java method: ReplaceMessages.constructNonGenderedMsgNode(com.google.javascript.jscomp.JsMessage message, java.util.Map<java.lang.String,com.google.javascript.rhino.Node> placeholderMap, com.google.javascript.jscomp.ReplaceMessages.MsgOptions options, com.google.javascript.rhino.Node nodeToReplace)
func (this_ *ReplaceMessages) priv_ReplaceMessages_constructNonGenderedMsgNode_4867(message *JsMessage, placeholderMap j2gort.Ext_java_util_Map, options *ReplaceMessages_MsgOptions, nodeToReplace pkg_rhino.Node) pkg_rhino.Node {
	_ = message
	_ = placeholderMap
	_ = options
	_ = nodeToReplace
	if placeholderMap.IsEmpty() {
		if options.IsIcuTemplate() {
			return this_.priv_ReplaceMessages_createNodeForMsgString_4868(options, message.AsIcuMessageString_0())
		} else if message.JsPlaceholderNames().IsEmpty() {
			return this_.priv_ReplaceMessages_createNodeForMsgString_4868(options, message.AsJsMessageString_0())
		} else {
			var icuMsgString string = message.AsIcuMessageString_0()
			_ = icuMsgString
			if ReplaceMessages_IsStartOfIcuMessage(icuMsgString) {
				return this_.priv_ReplaceMessages_createNodeForMsgString_4868(options, icuMsgString)
			} else {
				panic(NewJsMessageVisitor_MalformedException("The translated message has placeholders, but the definition in the JS code does not.", nodeToReplace))
			}
		}
	}
	return this_.priv_ReplaceMessages_buildMessageNode_4865(j2gort.AsIface[j2gort.Ext_java_util_List](message.GetParts()), options, nodeToReplace, func() ReplaceMessages_PlaceholderNodeProvider {
		bound_16 := placeholderMap
		return NewLambda_ReplaceMessages_PlaceholderNodeProvider(func(mrp0_ string) pkg_rhino.Node {
			return j2gort.CastNode[pkg_rhino.Node](bound_16.Get(j2gort.BoxString(mrp0_)))
		})
	}(), false)
}

// Java method: ReplaceMessages.createNodeForMsgString(com.google.javascript.jscomp.ReplaceMessages.MsgOptions options, java.lang.String s)
func (this_ *ReplaceMessages) priv_ReplaceMessages_createNodeForMsgString_4868(options *ReplaceMessages_MsgOptions, s string) pkg_rhino.Node {
	_ = options
	_ = s
	if options.EscapeLessThan() {
		s = javalang.StringReplace(s, "<", "&lt;")
	}
	if options.UnescapeHtmlEntities() {
		s = javalang.StringReplace(javalang.StringReplace(javalang.StringReplace(javalang.StringReplace(javalang.StringReplace(s, "&lt;", "<"), "&gt;", ">"), "&apos;", "'"), "&quot;", "\""), "&amp;", "&")
	}
	return this_.astFactory.CreateString(s)
}

// Java static method: ReplaceMessages.mergeStringParts(java.util.List<com.google.javascript.jscomp.JsMessage.Part> parts)
func ReplaceMessages_MergeStringParts(parts j2gort.Ext_java_util_List) j2gort.Ext_java_util_List {
	clinit_ReplaceMessages()
	_ = parts
	var result []JsMessage_Part = []JsMessage_Part{}
	_ = result
	{
		cur_17 := j2gort.NewIterCursor(j2gort.CheckNotNull(parts))
		for cur_17.HasNext() {
			rawElem_18 := cur_17.Next()
			part := j2gort.AsIface[JsMessage_Part](rawElem_18)
			_ = part
			if part.IsPlaceholder() {
				result = append(result, part)
			} else {
				var lastPart JsMessage_Part
				_ = lastPart
				if len(result) == 0 {
					lastPart = nil
				} else {
					lastPart = j2gort.SliceGetLast(result)
				}
				if (lastPart == nil || lastPart.IsJavaNil_()) || lastPart.IsPlaceholder() {
					result = append(result, part)
				} else {
					result[(int32(len(result)) - int32(1))] = Upcast_JsMessage_Part(JsMessage_StringPart_Create((j2gort.StrValString(lastPart.GetString()) + j2gort.StrValString(part.GetString()))))
				}
			}
		}
	}
	return j2gort.WrapTypedArrayList(result)
}

// Java static method: ReplaceMessages.isStartOfIcuMessage(java.lang.String part)
func ReplaceMessages_IsStartOfIcuMessage(part string) bool {
	clinit_ReplaceMessages()
	_ = part
	if !(javalang.StringStartsWith(j2gort.CheckNotNullString(part), "{")) {
		return false
	}
	var commaIndex int32 = j2gort.JStrIndexOfCpFrom(part, int32(uint16(44)), int32(1))
	_ = commaIndex
	if commaIndex <= int32(1) {
		return false
	}
	var nextBracketIndex int32 = j2gort.JStrIndexOfCpFrom(part, int32(uint16(123)), int32(1))
	_ = nextBracketIndex
	return (((nextBracketIndex == int32(-1)) || (nextBracketIndex > commaIndex)) && (j2gort.JStrStartsWithOffset(part, "plural,", (commaIndex+int32(1))) || j2gort.JStrStartsWithOffset(part, "select,", (commaIndex+int32(1)))))
}

// Java record: com.google.javascript.jscomp.ReplaceMessages.MsgProtectionData
type ReplaceMessages_MsgProtectionData struct {
	message               *JsMessage                  `gson:"message"`
	messageOptions        *ReplaceMessages_MsgOptions `gson:"messageOptions"`
	messageNode           pkg_rhino.Node              `gson:"messageNode"`
	templateTextNode      pkg_rhino.Node              `gson:"templateTextNode"`
	placeholderValuesNode pkg_rhino.Node              `gson:"placeholderValuesNode"`
}

func (this_ *ReplaceMessages_MsgProtectionData) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceMessages_MsgProtectionData) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceMessages_MsgProtectionData) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceMessages$MsgProtectionData", "MsgProtectionData"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceMessages.MsgProtectionData", "com.google.javascript.jscomp.ReplaceMessages$MsgProtectionData", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceMessages_MsgProtectionData](v_)
	})
}
func (*ReplaceMessages_MsgProtectionData) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *ReplaceMessages_MsgProtectionData) initFields_ReplaceMessages_MsgProtectionData() {
	_ = this_
}

// Java constructor: MsgProtectionData(com.google.javascript.jscomp.JsMessage message, com.google.javascript.rhino.Node messageNode, com.google.javascript.rhino.Node templateTextNode, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node placeholderValuesNode, com.google.javascript.jscomp.ReplaceMessages.MsgOptions messageOptions)
func NewReplaceMessages_MsgProtectionData(message *JsMessage, messageNode pkg_rhino.Node, templateTextNode pkg_rhino.Node, placeholderValuesNode pkg_rhino.Node, messageOptions *ReplaceMessages_MsgOptions) *ReplaceMessages_MsgProtectionData {
	this_ := &ReplaceMessages_MsgProtectionData{}
	InitReplaceMessages_MsgProtectionData(this_, message, messageNode, templateTextNode, placeholderValuesNode, messageOptions)
	return this_
}

func InitReplaceMessages_MsgProtectionData(this_ *ReplaceMessages_MsgProtectionData, message *JsMessage, messageNode pkg_rhino.Node, templateTextNode pkg_rhino.Node, placeholderValuesNode pkg_rhino.Node, messageOptions *ReplaceMessages_MsgOptions) {
	_ = message
	_ = messageNode
	_ = templateTextNode
	_ = placeholderValuesNode
	_ = messageOptions
	this_.initFields_ReplaceMessages_MsgProtectionData()
	this_.message = message
	this_.messageNode = messageNode
	this_.templateTextNode = templateTextNode
	this_.placeholderValuesNode = placeholderValuesNode
	this_.messageOptions = messageOptions
}

func (this_ *ReplaceMessages_MsgProtectionData) Message() *JsMessage { return this_.message }
func (this_ *ReplaceMessages_MsgProtectionData) MessageNode() pkg_rhino.Node {
	return this_.messageNode
}
func (this_ *ReplaceMessages_MsgProtectionData) TemplateTextNode() pkg_rhino.Node {
	return this_.templateTextNode
}
func (this_ *ReplaceMessages_MsgProtectionData) PlaceholderValuesNode() pkg_rhino.Node {
	return this_.placeholderValuesNode
}
func (this_ *ReplaceMessages_MsgProtectionData) MessageOptions() *ReplaceMessages_MsgOptions {
	return this_.messageOptions
}
func (this_ *ReplaceMessages_MsgProtectionData) Equals(other_ any) bool {
	o_, ok_ := other_.(*ReplaceMessages_MsgProtectionData)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if !javalang.ObjectsEquals(j2gort.PtrToAny(this_.message), j2gort.PtrToAny(o_.message)) {
		return false
	}
	if this_.messageNode != o_.messageNode {
		return false
	}
	if this_.templateTextNode != o_.templateTextNode {
		return false
	}
	if this_.placeholderValuesNode != o_.placeholderValuesNode {
		return false
	}
	if !javalang.ObjectsEquals(j2gort.PtrToAny(this_.messageOptions), j2gort.PtrToAny(o_.messageOptions)) {
		return false
	}
	return true
}
func (this_ *ReplaceMessages_MsgProtectionData) HashCode() int32 {
	var h_ int32 = 0
	{
		hp_19 := this_.message
		h_ = 31 * h_
		if hp_19 != nil {
			h_ += hp_19.HashCode()
		}
	}
	h_ = 31*h_ + int32(this_.messageNode)
	h_ = 31*h_ + int32(this_.templateTextNode)
	h_ = 31*h_ + int32(this_.placeholderValuesNode)
	{
		hp_20 := this_.messageOptions
		h_ = 31 * h_
		if hp_20 != nil {
			h_ += hp_20.HashCode()
		}
	}
	return h_
}
func (this_ *ReplaceMessages_MsgProtectionData) String() string {
	_ = fmt.Sprint
	return "MsgProtectionData[" + "message=" + j2gort.StrVal(j2gort.PtrToAny(this_.message)) + ", " + "messageNode=" + j2gort.StrVal(j2gort.BoxNode(this_.messageNode)) + ", " + "templateTextNode=" + j2gort.StrVal(j2gort.BoxNode(this_.templateTextNode)) + ", " + "placeholderValuesNode=" + j2gort.StrVal(j2gort.BoxNode(this_.placeholderValuesNode)) + ", " + "messageOptions=" + j2gort.StrVal(j2gort.PtrToAny(this_.messageOptions)) + "]"
}

// Java class: com.google.javascript.jscomp.ReplaceMessages.MsgProtectionPass
type ReplaceMessages_MsgProtectionPass struct {
	JsMessageVisitorBase
	outer_ *ReplaceMessages
}

func (this_ *ReplaceMessages_MsgProtectionPass) BindSelf_(s any) {
	_ = s
	this_.JsMessageVisitorBase.BindSelf_(s)
}
func (this_ *ReplaceMessages_MsgProtectionPass) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceMessages_MsgProtectionPass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceMessages$MsgProtectionPass", "MsgProtectionPass"
}
func (this_ *ReplaceMessages_MsgProtectionPass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.JsMessageVisitorBase.IdentityHash_()
}
func (*ReplaceMessages_MsgProtectionPass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceMessages.MsgProtectionPass", "com.google.javascript.jscomp.ReplaceMessages$MsgProtectionPass", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceMessages_MsgProtectionPass](v_)
	})
}
func (this_ *ReplaceMessages_MsgProtectionPass) SetOuter_(o *ReplaceMessages)  { this_.outer_ = o }
func (this_ *ReplaceMessages_MsgProtectionPass) GetOuter_() *ReplaceMessages   { return this_.outer_ }
func (*ReplaceMessages_MsgProtectionPass) Implements_jscomp_JsMessageVisitor() {}
func (*ReplaceMessages_MsgProtectionPass) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*ReplaceMessages_MsgProtectionPass) Implements_jscomp_CompilerPass()           {}
func (*ReplaceMessages_MsgProtectionPass) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *ReplaceMessages_MsgProtectionPass) initFields_ReplaceMessages_MsgProtectionPass() {
	_ = this_
}

// Java constructor: MsgProtectionPass()
func NewReplaceMessages_MsgProtectionPass(outer_ *ReplaceMessages) *ReplaceMessages_MsgProtectionPass {
	this_ := &ReplaceMessages_MsgProtectionPass{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitReplaceMessages_MsgProtectionPass(this_)
	return this_
}

func InitReplaceMessages_MsgProtectionPass(this_ *ReplaceMessages_MsgProtectionPass) {
	InitJsMessageVisitor(&this_.JsMessageVisitorBase, this_.outer_.compiler, this_.outer_.bundle.IdGenerator())
	this_.initFields_ReplaceMessages_MsgProtectionPass()
}

// Java method: MsgProtectionPass.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ReplaceMessages_MsgProtectionPass) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeUtil_CreateSynthesizedExternsSymbol(this_.GetJsMessageVisitorBase().compiler, "__jscomp_define_msg__")
	NodeUtil_CreateSynthesizedExternsSymbol(this_.GetJsMessageVisitorBase().compiler, "__jscomp_msg_fallback__")
	this_.JsMessageVisitorBase.Process_2(externs, root)
}

// Java method: MsgProtectionPass.processIcuTemplateDefinition(com.google.javascript.jscomp.IcuTemplateDefinition definition)
func (this_ *ReplaceMessages_MsgProtectionPass) ProcessIcuTemplateDefinition(definition IcuTemplateDefinition) {
	_ = definition
	this_.priv_MsgProtectionPass_performMessageProtection_4874(NewReplaceMessages_MsgProtectionData(definition.GetMessage(), definition.GetMessageNode(), definition.GetTemplateTextNode(), pkg_rhino.Node(0), StaticGet_ReplaceMessages_ICU_MSG_OPTIONS()))
}

// Java method: MsgProtectionPass.processJsMessageDefinition(com.google.javascript.jscomp.JsMessageDefinition definition)
func (this_ *ReplaceMessages_MsgProtectionPass) ProcessJsMessageDefinition(definition JsMessageDefinition) {
	_ = definition
	this_.priv_MsgProtectionPass_performMessageProtection_4874(NewReplaceMessages_MsgProtectionData(definition.GetMessage(), definition.GetMessageNode(), definition.GetTemplateTextNode(), definition.GetPlaceholderValuesNode(), this_.outer_.priv_ReplaceMessages_getMsgOptionsFromDefinition_4857(definition)))
}

// Java method: MsgProtectionPass.performMessageProtection(com.google.javascript.jscomp.ReplaceMessages.MsgProtectionData msgProtectionData)
func (this_ *ReplaceMessages_MsgProtectionPass) priv_MsgProtectionPass_performMessageProtection_4874(msgProtectionData *ReplaceMessages_MsgProtectionData) {
	_ = msgProtectionData
	var message *JsMessage = msgProtectionData.Message()
	_ = message
	var callNode pkg_rhino.Node = msgProtectionData.MessageNode()
	_ = callNode
	var originalMessageString pkg_rhino.Node = msgProtectionData.TemplateTextNode()
	_ = originalMessageString
	var placeholdersNode pkg_rhino.Node = msgProtectionData.PlaceholderValuesNode()
	_ = placeholdersNode
	var msgOptions *ReplaceMessages_MsgOptions = msgProtectionData.MessageOptions()
	_ = msgOptions
	if !(callNode.IsCall()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(callNode))
	}
	var googGetMsg pkg_rhino.Node = callNode.GetFirstChild()
	_ = googGetMsg
	var protectionFunctionName string = "__jscomp_define_msg__"
	_ = protectionFunctionName
	var newCallee pkg_rhino.Node = this_.priv_MsgProtectionPass_createProtectionFunctionCallee_4875(protectionFunctionName).Srcref(googGetMsg)
	_ = newCallee
	var msgPropertiesNode pkg_rhino.Node = this_.outer_.priv_ReplaceMessages_createMsgPropertiesNode_4858(message, msgOptions).SrcrefTree(originalMessageString)
	_ = msgPropertiesNode
	var newCallNode pkg_rhino.Node = this_.outer_.astFactory.CreateCall(newCallee, AstFactory_S_Type_1_Node(callNode), []pkg_rhino.Node{msgPropertiesNode}).Srcref(callNode)
	_ = newCallNode
	newCallNode.SetSideEffectFlags_1_int32(int32(0))
	if placeholdersNode != 0 {
		if !(placeholdersNode.IsObjectLit()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(placeholdersNode))
		}
		{
			var strKey pkg_rhino.Node = placeholdersNode.GetFirstChild()
			_ = strKey
			for ; strKey != 0; strKey = strKey.GetNext() {
				if !(strKey.IsStringKey()) {
					j2gort.CheckStateMsg(false, j2gort.BoxNode(strKey))
				}
				strKey.SetQuotedStringKey()
			}
		}
		newCallNode.AddChildToBack(placeholdersNode.Detach())
	}
	callNode.ReplaceWith(newCallNode)
	this_.GetJsMessageVisitorBase().compiler.ReportChangeToEnclosingScope(newCallNode)
}

// Java method: MsgProtectionPass.createProtectionFunctionCallee(java.lang.String protectionFunctionName)
func (this_ *ReplaceMessages_MsgProtectionPass) priv_MsgProtectionPass_createProtectionFunctionCallee_4875(protectionFunctionName string) pkg_rhino.Node {
	_ = protectionFunctionName
	var callee pkg_rhino.Node = this_.outer_.astFactory.CreateNameWithUnknownType(protectionFunctionName)
	_ = callee
	callee.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
	return callee
}

// Java method: MsgProtectionPass.processMessageFallback(com.google.javascript.rhino.Node callNode, com.google.javascript.jscomp.JsMessage message1, com.google.javascript.jscomp.JsMessage message2)
func (this_ *ReplaceMessages_MsgProtectionPass) ProcessMessageFallback(callNode pkg_rhino.Node, message1 *JsMessage, message2 *JsMessage) {
	_ = callNode
	_ = message1
	_ = message2
	var originalCallee pkg_rhino.Node = j2gort.CheckNotNullNode(callNode.GetFirstChild())
	_ = originalCallee
	var fallbackCallee pkg_rhino.Node = this_.priv_MsgProtectionPass_createProtectionFunctionCallee_4875("__jscomp_msg_fallback__").Srcref(originalCallee)
	_ = fallbackCallee
	var originalFirstArg pkg_rhino.Node = j2gort.CheckNotNullNode(originalCallee.GetNext())
	_ = originalFirstArg
	var firstMsgKey pkg_rhino.Node = this_.outer_.astFactory.CreateString(message1.GetKey()).Srcref(originalFirstArg)
	_ = firstMsgKey
	var originalSecondArg pkg_rhino.Node = j2gort.CheckNotNullNode(originalFirstArg.GetNext())
	_ = originalSecondArg
	var secondMsgKey pkg_rhino.Node = this_.outer_.astFactory.CreateString(message2.GetKey()).Srcref(originalSecondArg)
	_ = secondMsgKey
	var newCallNode pkg_rhino.Node = this_.outer_.astFactory.CreateCall(fallbackCallee, AstFactory_S_Type_1_Node(callNode), []pkg_rhino.Node{firstMsgKey, originalFirstArg.Detach(), secondMsgKey, originalSecondArg.Detach()}).Srcref(callNode)
	_ = newCallNode
	newCallNode.SetSideEffectFlags_1_int32(int32(0))
	callNode.ReplaceWith(newCallNode)
	this_.GetJsMessageVisitorBase().compiler.ReportChangeToEnclosingScope(newCallNode)
}

// Java class: com.google.javascript.jscomp.ReplaceMessages.QuotedKeyObjectLitBuilder
type ReplaceMessages_QuotedKeyObjectLitBuilder struct {
	outer_            *ReplaceMessages
	keyToValueNodeMap j2gort.OrderedMap[string, pkg_rhino.Node] `gson:"keyToValueNodeMap"`
	idHash_           atomic.Int32
}

func (this_ *ReplaceMessages_QuotedKeyObjectLitBuilder) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceMessages_QuotedKeyObjectLitBuilder) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceMessages_QuotedKeyObjectLitBuilder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceMessages$QuotedKeyObjectLitBuilder", "QuotedKeyObjectLitBuilder"
}
func (this_ *ReplaceMessages_QuotedKeyObjectLitBuilder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReplaceMessages_QuotedKeyObjectLitBuilder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceMessages.QuotedKeyObjectLitBuilder", "com.google.javascript.jscomp.ReplaceMessages$QuotedKeyObjectLitBuilder", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceMessages_QuotedKeyObjectLitBuilder](v_)
	})
}
func (this_ *ReplaceMessages_QuotedKeyObjectLitBuilder) SetOuter_(o *ReplaceMessages) {
	this_.outer_ = o
}
func (this_ *ReplaceMessages_QuotedKeyObjectLitBuilder) GetOuter_() *ReplaceMessages {
	return this_.outer_
}

func (this_ *ReplaceMessages_QuotedKeyObjectLitBuilder) initFields_ReplaceMessages_QuotedKeyObjectLitBuilder() {
	_ = this_
	this_.keyToValueNodeMap = j2gort.NewOrderedMap[string, pkg_rhino.Node]()
}

// Java constructor: QuotedKeyObjectLitBuilder()
func NewReplaceMessages_QuotedKeyObjectLitBuilder(outer_ *ReplaceMessages) *ReplaceMessages_QuotedKeyObjectLitBuilder {
	this_ := &ReplaceMessages_QuotedKeyObjectLitBuilder{}
	this_.outer_ = outer_
	InitReplaceMessages_QuotedKeyObjectLitBuilder(this_)
	return this_
}

func InitReplaceMessages_QuotedKeyObjectLitBuilder(this_ *ReplaceMessages_QuotedKeyObjectLitBuilder) {
	this_.initFields_ReplaceMessages_QuotedKeyObjectLitBuilder()
}

// Java method: QuotedKeyObjectLitBuilder.addString(java.lang.String key, java.lang.String value)
func (this_ *ReplaceMessages_QuotedKeyObjectLitBuilder) priv_QuotedKeyObjectLitBuilder_addString_4876(key string, value string) *ReplaceMessages_QuotedKeyObjectLitBuilder {
	_ = key
	_ = value
	return this_.priv_QuotedKeyObjectLitBuilder_addNode_4877(key, this_.outer_.astFactory.CreateString(value))
}

// Java method: QuotedKeyObjectLitBuilder.addNode(java.lang.String key, com.google.javascript.rhino.Node node)
func (this_ *ReplaceMessages_QuotedKeyObjectLitBuilder) priv_QuotedKeyObjectLitBuilder_addNode_4877(key string, node pkg_rhino.Node) *ReplaceMessages_QuotedKeyObjectLitBuilder {
	_ = key
	_ = node
	if !(!(this_.keyToValueNodeMap.ContainsKey(key))) {
		j2gort.CheckStateFmt(false, "repeated key: %s", j2gort.BoxString(key))
	}
	this_.keyToValueNodeMap.Put(key, node)
	return this_
}

// Java method: QuotedKeyObjectLitBuilder.build()
func (this_ *ReplaceMessages_QuotedKeyObjectLitBuilder) priv_QuotedKeyObjectLitBuilder_build_4878() pkg_rhino.Node {
	var result pkg_rhino.Node = this_.outer_.astFactory.CreateObjectLit_1([]pkg_rhino.Node{})
	_ = result
	{
		cur_19 := this_.keyToValueNodeMap.IterEntries()
		for cur_19.NextEntry() {
			result.AddChildToBack(this_.outer_.astFactory.CreateQuotedStringKey(cur_19.Key(), cur_19.Value()))
		}
	}
	return result
}

// Java class: com.google.javascript.jscomp.ReplaceMessages.ReplacementCompletionPass
type ReplaceMessages_ReplacementCompletionPass struct {
	outer_            *ReplaceMessages
	translatedMsgKeys j2gort.OrderedSet[string] `gson:"translatedMsgKeys"`
	idHash_           atomic.Int32
}

func (this_ *ReplaceMessages_ReplacementCompletionPass) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceMessages_ReplacementCompletionPass) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceMessages_ReplacementCompletionPass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceMessages$ReplacementCompletionPass", "ReplacementCompletionPass"
}
func (this_ *ReplaceMessages_ReplacementCompletionPass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReplaceMessages_ReplacementCompletionPass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceMessages.ReplacementCompletionPass", "com.google.javascript.jscomp.ReplaceMessages$ReplacementCompletionPass", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceMessages_ReplacementCompletionPass](v_)
	})
}
func (this_ *ReplaceMessages_ReplacementCompletionPass) SetOuter_(o *ReplaceMessages) {
	this_.outer_ = o
}
func (this_ *ReplaceMessages_ReplacementCompletionPass) GetOuter_() *ReplaceMessages {
	return this_.outer_
}
func (*ReplaceMessages_ReplacementCompletionPass) Implements_jscomp_CompilerPass() {}

func (this_ *ReplaceMessages_ReplacementCompletionPass) initFields_ReplaceMessages_ReplacementCompletionPass() {
	_ = this_
	this_.translatedMsgKeys = j2gort.NewOrderedSet[string]()
}

// Java constructor: ReplacementCompletionPass()
func NewReplaceMessages_ReplacementCompletionPass(outer_ *ReplaceMessages) *ReplaceMessages_ReplacementCompletionPass {
	this_ := &ReplaceMessages_ReplacementCompletionPass{}
	this_.outer_ = outer_
	InitReplaceMessages_ReplacementCompletionPass(this_)
	return this_
}

func InitReplaceMessages_ReplacementCompletionPass(this_ *ReplaceMessages_ReplacementCompletionPass) {
	this_.initFields_ReplaceMessages_ReplacementCompletionPass()
}

// Java method: ReplacementCompletionPass.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ReplaceMessages_ReplacementCompletionPass) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.outer_.compiler, root, NewReplaceMessages__Anon_77(this_))
}

// Java method: ReplacementCompletionPass.visitMsgDefinition(com.google.javascript.jscomp.ReplaceMessages.ProtectedJsMessage protectedJsMessage)
func (this_ *ReplaceMessages_ReplacementCompletionPass) VisitMsgDefinition(protectedJsMessage *ReplaceMessages_ProtectedJsMessage) {
	_ = protectedJsMessage
	func() {
		defer func() {
			if rec_20 := recover(); rec_20 != nil {
				jex_21 := javalang.RecoverException(rec_20)
				_ = jex_21
				if j2gort.InstanceOf[*JsMessageVisitor_MalformedException](jex_21) {
					e := j2gort.AsIface[*JsMessageVisitor_MalformedException](jex_21)
					_ = e
					this_.outer_.compiler.Report_1(JSError_Make_3(e.GetNode(), StaticGet_JsMessageVisitor_MESSAGE_TREE_MALFORMED(), []string{j2gort.ThrowableGetMessage(e)}))
				} else {
					panic(rec_20)
				}
			}
		}()
		var originalMsg *JsMessage = protectedJsMessage.jsMessage
		_ = originalMsg
		var nodeToReplace pkg_rhino.Node = protectedJsMessage.definitionNode
		_ = nodeToReplace
		var translatedMsg *JsMessage = this_.outer_.priv_ReplaceMessages_lookupMessage_4860(protectedJsMessage.definitionNode, this_.outer_.bundle, originalMsg)
		_ = translatedMsg
		var msgToUse *JsMessage = nil
		_ = msgToUse
		if translatedMsg != nil {
			msgToUse = translatedMsg
			this_.translatedMsgKeys.Add(originalMsg.GetKey())
		} else {
			if this_.outer_.strictReplacement {
				this_.outer_.compiler.Report_1(JSError_Make_3(nodeToReplace, StaticGet_ReplaceMessages_BUNDLE_DOES_NOT_HAVE_THE_MESSAGE(), []string{originalMsg.GetId()}))
			}
			msgToUse = originalMsg
		}
		var msgOptions *ReplaceMessages_MsgOptions = protectedJsMessage.GetMsgOptions()
		_ = msgOptions
		var placeholderMap j2gort.Ext_java_util_Map = j2gort.AsIface[j2gort.Ext_java_util_Map](ReplaceMessages_ExtractPlaceholderValuesMapOrThrow(protectedJsMessage.substitutionsNode))
		_ = placeholderMap
		var finalMsgConstructionExpression pkg_rhino.Node = this_.outer_.priv_ReplaceMessages_constructStringExprNode_4866(msgToUse, placeholderMap, msgOptions, nodeToReplace)
		_ = finalMsgConstructionExpression
		finalMsgConstructionExpression.SrcrefTreeIfMissing(nodeToReplace)
		nodeToReplace.ReplaceWith(finalMsgConstructionExpression)
		this_.outer_.compiler.ReportChangeToEnclosingScope(finalMsgConstructionExpression)
	}()
}

// Java method: ReplacementCompletionPass.visitMsgFallback(com.google.javascript.jscomp.ReplaceMessages.ProtectedMsgFallback protectedMsgFallback)
func (this_ *ReplaceMessages_ReplacementCompletionPass) priv_ReplacementCompletionPass_visitMsgFallback_4880(protectedMsgFallback *ReplaceMessages_ProtectedMsgFallback) {
	_ = protectedMsgFallback
	var valueNodeToUse pkg_rhino.Node = pkg_rhino.Node(0)
	_ = valueNodeToUse
	if this_.translatedMsgKeys.Contains(protectedMsgFallback.firstMsgKey) {
		valueNodeToUse = protectedMsgFallback.firstMsgValue
	} else if this_.translatedMsgKeys.Contains(protectedMsgFallback.secondMsgKey) {
		valueNodeToUse = protectedMsgFallback.secondMsgValue
	} else {
		valueNodeToUse = protectedMsgFallback.firstMsgValue
	}
	valueNodeToUse.Detach()
	protectedMsgFallback.callNode.ReplaceWith(valueNodeToUse)
	this_.outer_.compiler.ReportChangeToEnclosingScope(valueNodeToUse)
}

// Java class: com.google.javascript.jscomp.ReplaceMessages.ProtectedMsgFallback
type ReplaceMessages_ProtectedMsgFallback struct {
	firstMsgKey    string         `gson:"firstMsgKey"`
	secondMsgKey   string         `gson:"secondMsgKey"`
	callNode       pkg_rhino.Node `gson:"callNode"`
	firstMsgValue  pkg_rhino.Node `gson:"firstMsgValue"`
	secondMsgValue pkg_rhino.Node `gson:"secondMsgValue"`
	idHash_        atomic.Int32
}

func (this_ *ReplaceMessages_ProtectedMsgFallback) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceMessages_ProtectedMsgFallback) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceMessages_ProtectedMsgFallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceMessages$ProtectedMsgFallback", "ProtectedMsgFallback"
}
func (this_ *ReplaceMessages_ProtectedMsgFallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReplaceMessages_ProtectedMsgFallback) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceMessages.ProtectedMsgFallback", "com.google.javascript.jscomp.ReplaceMessages$ProtectedMsgFallback", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceMessages_ProtectedMsgFallback](v_)
	})
}

func (this_ *ReplaceMessages_ProtectedMsgFallback) initFields_ReplaceMessages_ProtectedMsgFallback() {
	_ = this_
	this_.firstMsgKey = j2gort.NullString
	this_.secondMsgKey = j2gort.NullString
}

// Java constructor: ProtectedMsgFallback(com.google.javascript.rhino.Node callNode, java.lang.String firstMsgKey, com.google.javascript.rhino.Node firstMsgValue, java.lang.String secondMsgKey, com.google.javascript.rhino.Node secondMsgValue)
func NewReplaceMessages_ProtectedMsgFallback(callNode pkg_rhino.Node, firstMsgKey string, firstMsgValue pkg_rhino.Node, secondMsgKey string, secondMsgValue pkg_rhino.Node) *ReplaceMessages_ProtectedMsgFallback {
	this_ := &ReplaceMessages_ProtectedMsgFallback{}
	InitReplaceMessages_ProtectedMsgFallback(this_, callNode, firstMsgKey, firstMsgValue, secondMsgKey, secondMsgValue)
	return this_
}

func InitReplaceMessages_ProtectedMsgFallback(this_ *ReplaceMessages_ProtectedMsgFallback, callNode pkg_rhino.Node, firstMsgKey string, firstMsgValue pkg_rhino.Node, secondMsgKey string, secondMsgValue pkg_rhino.Node) {
	_ = callNode
	_ = firstMsgKey
	_ = firstMsgValue
	_ = secondMsgKey
	_ = secondMsgValue
	this_.initFields_ReplaceMessages_ProtectedMsgFallback()
	this_.callNode = callNode
	this_.firstMsgKey = firstMsgKey
	this_.firstMsgValue = firstMsgValue
	this_.secondMsgKey = secondMsgKey
	this_.secondMsgValue = secondMsgValue
}

// Java static method: ProtectedMsgFallback.fromAstNode(com.google.javascript.rhino.Node n)
func ReplaceMessages_ProtectedMsgFallback_FromAstNode(n pkg_rhino.Node) *ReplaceMessages_ProtectedMsgFallback {
	_ = n
	if !(n.IsCall()) {
		return nil
	}
	var callee pkg_rhino.Node = n.GetFirstChild()
	_ = callee
	if !(callee.MatchesName_1_string("__jscomp_msg_fallback__")) {
		return nil
	}
	if !(n.HasXChildren(int32(5))) {
		j2gort.CheckStateFmt(false, "bad message fallback call: %s", j2gort.BoxNode(n))
	}
	var firstMsgKeyNode pkg_rhino.Node = callee.GetNext()
	_ = firstMsgKeyNode
	var firstMsgKey string = firstMsgKeyNode.GetString()
	_ = firstMsgKey
	var firstMsgValue pkg_rhino.Node = firstMsgKeyNode.GetNext()
	_ = firstMsgValue
	var secondMsgKeyNode pkg_rhino.Node = firstMsgValue.GetNext()
	_ = secondMsgKeyNode
	var secondMsgKey string = secondMsgKeyNode.GetString()
	_ = secondMsgKey
	var secondMsgValue pkg_rhino.Node = secondMsgKeyNode.GetNext()
	_ = secondMsgValue
	return NewReplaceMessages_ProtectedMsgFallback(n, firstMsgKey, firstMsgValue, secondMsgKey, secondMsgValue)
}

// Java record: com.google.javascript.jscomp.ReplaceMessages.FullReplacementMsgData
type ReplaceMessages_FullReplacementMsgData struct {
	message             *JsMessage                                        `gson:"message"`
	messageOptions      *ReplaceMessages_MsgOptions                       `gson:"messageOptions"`
	placeholderValueMap j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"placeholderValueMap"`
	messageNode         pkg_rhino.Node                                    `gson:"messageNode"`
}

func (this_ *ReplaceMessages_FullReplacementMsgData) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceMessages_FullReplacementMsgData) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceMessages_FullReplacementMsgData) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceMessages$FullReplacementMsgData", "FullReplacementMsgData"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceMessages.FullReplacementMsgData", "com.google.javascript.jscomp.ReplaceMessages$FullReplacementMsgData", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceMessages_FullReplacementMsgData](v_)
	})
}
func (*ReplaceMessages_FullReplacementMsgData) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *ReplaceMessages_FullReplacementMsgData) initFields_ReplaceMessages_FullReplacementMsgData() {
	_ = this_
}

// Java constructor: FullReplacementMsgData(com.google.javascript.jscomp.JsMessage message, com.google.javascript.rhino.Node messageNode, com.google.javascript.jscomp.ReplaceMessages.MsgOptions messageOptions, com.google.common.collect.ImmutableMap<java.lang.String,com.google.javascript.rhino.Node> placeholderValueMap)
func NewReplaceMessages_FullReplacementMsgData(message *JsMessage, messageNode pkg_rhino.Node, messageOptions *ReplaceMessages_MsgOptions, placeholderValueMap j2gort.Ext_com_google_common_collect_ImmutableMap) *ReplaceMessages_FullReplacementMsgData {
	this_ := &ReplaceMessages_FullReplacementMsgData{}
	InitReplaceMessages_FullReplacementMsgData(this_, message, messageNode, messageOptions, placeholderValueMap)
	return this_
}

func InitReplaceMessages_FullReplacementMsgData(this_ *ReplaceMessages_FullReplacementMsgData, message *JsMessage, messageNode pkg_rhino.Node, messageOptions *ReplaceMessages_MsgOptions, placeholderValueMap j2gort.Ext_com_google_common_collect_ImmutableMap) {
	_ = message
	_ = messageNode
	_ = messageOptions
	_ = placeholderValueMap
	this_.initFields_ReplaceMessages_FullReplacementMsgData()
	this_.message = message
	this_.messageNode = messageNode
	this_.messageOptions = messageOptions
	this_.placeholderValueMap = placeholderValueMap
}

func (this_ *ReplaceMessages_FullReplacementMsgData) Message() *JsMessage { return this_.message }
func (this_ *ReplaceMessages_FullReplacementMsgData) MessageNode() pkg_rhino.Node {
	return this_.messageNode
}
func (this_ *ReplaceMessages_FullReplacementMsgData) MessageOptions() *ReplaceMessages_MsgOptions {
	return this_.messageOptions
}
func (this_ *ReplaceMessages_FullReplacementMsgData) PlaceholderValueMap() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.placeholderValueMap
}
func (this_ *ReplaceMessages_FullReplacementMsgData) Equals(other_ any) bool {
	o_, ok_ := other_.(*ReplaceMessages_FullReplacementMsgData)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if !javalang.ObjectsEquals(j2gort.PtrToAny(this_.message), j2gort.PtrToAny(o_.message)) {
		return false
	}
	if this_.messageNode != o_.messageNode {
		return false
	}
	if !javalang.ObjectsEquals(j2gort.PtrToAny(this_.messageOptions), j2gort.PtrToAny(o_.messageOptions)) {
		return false
	}
	if !javalang.ObjectsEquals(this_.placeholderValueMap, o_.placeholderValueMap) {
		return false
	}
	return true
}
func (this_ *ReplaceMessages_FullReplacementMsgData) HashCode() int32 {
	var h_ int32 = 0
	{
		hp_22 := this_.message
		h_ = 31 * h_
		if hp_22 != nil {
			h_ += hp_22.HashCode()
		}
	}
	h_ = 31*h_ + int32(this_.messageNode)
	{
		hp_23 := this_.messageOptions
		h_ = 31 * h_
		if hp_23 != nil {
			h_ += hp_23.HashCode()
		}
	}
	h_ = 31*h_ + j2gort.ObjectHash(this_.placeholderValueMap)
	return h_
}
func (this_ *ReplaceMessages_FullReplacementMsgData) String() string {
	_ = fmt.Sprint
	return "FullReplacementMsgData[" + "message=" + j2gort.StrVal(j2gort.PtrToAny(this_.message)) + ", " + "messageNode=" + j2gort.StrVal(j2gort.BoxNode(this_.messageNode)) + ", " + "messageOptions=" + j2gort.StrVal(j2gort.PtrToAny(this_.messageOptions)) + ", " + "placeholderValueMap=" + j2gort.StrVal(this_.placeholderValueMap) + "]"
}

// Java class: com.google.javascript.jscomp.ReplaceMessages.FullReplacementPass
type ReplaceMessages_FullReplacementPass struct {
	JsMessageVisitorBase
	outer_ *ReplaceMessages
}

func (this_ *ReplaceMessages_FullReplacementPass) BindSelf_(s any) {
	_ = s
	this_.JsMessageVisitorBase.BindSelf_(s)
}
func (this_ *ReplaceMessages_FullReplacementPass) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceMessages_FullReplacementPass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceMessages$FullReplacementPass", "FullReplacementPass"
}
func (this_ *ReplaceMessages_FullReplacementPass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.JsMessageVisitorBase.IdentityHash_()
}
func (*ReplaceMessages_FullReplacementPass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceMessages.FullReplacementPass", "com.google.javascript.jscomp.ReplaceMessages$FullReplacementPass", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceMessages_FullReplacementPass](v_)
	})
}
func (this_ *ReplaceMessages_FullReplacementPass) SetOuter_(o *ReplaceMessages)  { this_.outer_ = o }
func (this_ *ReplaceMessages_FullReplacementPass) GetOuter_() *ReplaceMessages   { return this_.outer_ }
func (*ReplaceMessages_FullReplacementPass) Implements_jscomp_JsMessageVisitor() {}
func (*ReplaceMessages_FullReplacementPass) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*ReplaceMessages_FullReplacementPass) Implements_jscomp_CompilerPass()           {}
func (*ReplaceMessages_FullReplacementPass) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *ReplaceMessages_FullReplacementPass) initFields_ReplaceMessages_FullReplacementPass() {
	_ = this_
}

// Java constructor: FullReplacementPass()
func NewReplaceMessages_FullReplacementPass(outer_ *ReplaceMessages) *ReplaceMessages_FullReplacementPass {
	this_ := &ReplaceMessages_FullReplacementPass{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitReplaceMessages_FullReplacementPass(this_)
	return this_
}

func InitReplaceMessages_FullReplacementPass(this_ *ReplaceMessages_FullReplacementPass) {
	InitJsMessageVisitor(&this_.JsMessageVisitorBase, this_.outer_.compiler, this_.outer_.bundle.IdGenerator())
	this_.initFields_ReplaceMessages_FullReplacementPass()
}

// Java method: FullReplacementPass.processMessageFallback(com.google.javascript.rhino.Node callNode, com.google.javascript.jscomp.JsMessage message1, com.google.javascript.jscomp.JsMessage message2)
func (this_ *ReplaceMessages_FullReplacementPass) ProcessMessageFallback(callNode pkg_rhino.Node, message1 *JsMessage, message2 *JsMessage) {
	_ = callNode
	_ = message1
	_ = message2
	var isFirstMessageTranslated bool = (this_.outer_.priv_ReplaceMessages_lookupMessage_4860(callNode, this_.outer_.bundle, message1) != nil)
	_ = isFirstMessageTranslated
	var isSecondMessageTranslated bool = (this_.outer_.priv_ReplaceMessages_lookupMessage_4860(callNode, this_.outer_.bundle, message2) != nil)
	_ = isSecondMessageTranslated
	var replacementNode pkg_rhino.Node
	_ = replacementNode
	if isSecondMessageTranslated && !(isFirstMessageTranslated) {
		replacementNode = callNode.GetChildAtIndex(int32(2))
	} else {
		replacementNode = callNode.GetSecondChild()
	}
	callNode.ReplaceWith(replacementNode.Detach())
	var changeScope pkg_rhino.Node = ChangeTracker_GetEnclosingChangeScopeRoot(replacementNode)
	_ = changeScope
	if changeScope != 0 {
		this_.GetJsMessageVisitorBase().compiler.ReportChangeToChangeScope(changeScope)
	}
}

// Java method: FullReplacementPass.processIcuTemplateDefinition(com.google.javascript.jscomp.IcuTemplateDefinition definition)
func (this_ *ReplaceMessages_FullReplacementPass) ProcessIcuTemplateDefinition(definition IcuTemplateDefinition) {
	_ = definition
	this_.priv_FullReplacementPass_processFullReplacement_4885(NewReplaceMessages_FullReplacementMsgData(definition.GetMessage(), definition.GetMessageNode(), StaticGet_ReplaceMessages_ICU_MSG_OPTIONS(), j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()))
}

// Java method: FullReplacementPass.processJsMessageDefinition(com.google.javascript.jscomp.JsMessageDefinition definition)
func (this_ *ReplaceMessages_FullReplacementPass) ProcessJsMessageDefinition(definition JsMessageDefinition) {
	_ = definition
	this_.priv_FullReplacementPass_processFullReplacement_4885(NewReplaceMessages_FullReplacementMsgData(definition.GetMessage(), definition.GetMessageNode(), this_.outer_.priv_ReplaceMessages_getMsgOptionsFromDefinition_4857(definition), definition.GetPlaceholderValueMap()))
}

// Java method: FullReplacementPass.processFullReplacement(com.google.javascript.jscomp.ReplaceMessages.FullReplacementMsgData fullReplacementMsgData)
func (this_ *ReplaceMessages_FullReplacementPass) priv_FullReplacementPass_processFullReplacement_4885(fullReplacementMsgData *ReplaceMessages_FullReplacementMsgData) {
	_ = fullReplacementMsgData
	var message *JsMessage = fullReplacementMsgData.Message()
	_ = message
	var msgNode pkg_rhino.Node = fullReplacementMsgData.MessageNode()
	_ = msgNode
	var options *ReplaceMessages_MsgOptions = fullReplacementMsgData.MessageOptions()
	_ = options
	var placeholderValueMap j2gort.Ext_com_google_common_collect_ImmutableMap = fullReplacementMsgData.PlaceholderValueMap()
	_ = placeholderValueMap
	var replacement *JsMessage = this_.outer_.priv_ReplaceMessages_lookupMessage_4860(msgNode, this_.outer_.bundle, message)
	_ = replacement
	if replacement == nil {
		if this_.outer_.strictReplacement {
			this_.GetJsMessageVisitorBase().compiler.Report_1(JSError_Make_3(msgNode, StaticGet_ReplaceMessages_BUNDLE_DOES_NOT_HAVE_THE_MESSAGE(), []string{message.GetId()}))
			return
		} else {
			replacement = message
		}
	}
	var newValue pkg_rhino.Node = pkg_rhino.Node(0)
	_ = newValue
	func() {
		defer func() {
			if rec_22 := recover(); rec_22 != nil {
				jex_23 := javalang.RecoverException(rec_22)
				_ = jex_23
				if j2gort.InstanceOf[*JsMessageVisitor_MalformedException](jex_23) {
					e := j2gort.AsIface[*JsMessageVisitor_MalformedException](jex_23)
					_ = e
					this_.GetJsMessageVisitorBase().compiler.Report_1(JSError_Make_3(e.GetNode(), StaticGet_JsMessageVisitor_MESSAGE_TREE_MALFORMED(), []string{j2gort.ThrowableGetMessage(e)}))
					newValue = msgNode
				} else {
					panic(rec_22)
				}
			}
		}()
		newValue = this_.outer_.priv_ReplaceMessages_constructStringExprNode_4866(replacement, j2gort.AsIface[j2gort.Ext_java_util_Map](placeholderValueMap), options, msgNode)
	}()
	if newValue != msgNode {
		newValue.SrcrefTreeIfMissing(msgNode)
		msgNode.ReplaceWith(newValue)
		this_.GetJsMessageVisitorBase().compiler.ReportChangeToEnclosingScope(newValue)
	}
}

// Java interface: com.google.javascript.jscomp.ReplaceMessages.PlaceholderNodeProvider
type ReplaceMessages_PlaceholderNodeProvider interface {
	Implements_jscomp_ReplaceMessages_PlaceholderNodeProvider()
	IsJavaNil_() bool
	Get(string) pkg_rhino.Node
}

func ToIface_ReplaceMessages_PlaceholderNodeProvider[T interface {
	ReplaceMessages_PlaceholderNodeProvider
	comparable
}](v T) ReplaceMessages_PlaceholderNodeProvider {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_ReplaceMessages_PlaceholderNodeProvider[T ReplaceMessages_PlaceholderNodeProvider](v T) ReplaceMessages_PlaceholderNodeProvider {
	var vi ReplaceMessages_PlaceholderNodeProvider = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceMessages.PlaceholderNodeProvider", "com.google.javascript.jscomp.ReplaceMessages$PlaceholderNodeProvider", func(v_ any) bool {
		return j2gort.InstanceOf[ReplaceMessages_PlaceholderNodeProvider](v_)
	})
}

type lambda_ReplaceMessages_PlaceholderNodeProvider struct {
	fn func(string) pkg_rhino.Node
}

func (l *lambda_ReplaceMessages_PlaceholderNodeProvider) IsJavaNil_() bool { return l == nil }
func NewLambda_ReplaceMessages_PlaceholderNodeProvider(fn func(string) pkg_rhino.Node) ReplaceMessages_PlaceholderNodeProvider {
	if fn == nil {
		return nil
	}
	return &lambda_ReplaceMessages_PlaceholderNodeProvider{fn: fn}
}

func (*lambda_ReplaceMessages_PlaceholderNodeProvider) Implements_jscomp_ReplaceMessages_PlaceholderNodeProvider() {
}
func (l *lambda_ReplaceMessages_PlaceholderNodeProvider) Get(p0_ string) pkg_rhino.Node {
	return l.fn(p0_)
}

// Java record: com.google.javascript.jscomp.ReplaceMessages.MsgOptions
type ReplaceMessages_MsgOptions struct {
	isIcuTemplate        bool `gson:"isIcuTemplate"`
	escapeLessThan       bool `gson:"escapeLessThan"`
	unescapeHtmlEntities bool `gson:"unescapeHtmlEntities"`
}

func (this_ *ReplaceMessages_MsgOptions) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceMessages_MsgOptions) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceMessages_MsgOptions) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceMessages$MsgOptions", "MsgOptions"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceMessages.MsgOptions", "com.google.javascript.jscomp.ReplaceMessages$MsgOptions", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceMessages_MsgOptions](v_)
	})
}
func (*ReplaceMessages_MsgOptions) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *ReplaceMessages_MsgOptions) initFields_ReplaceMessages_MsgOptions() {
	_ = this_
}

// Java constructor: MsgOptions(boolean isIcuTemplate, boolean escapeLessThan, boolean unescapeHtmlEntities)
func NewReplaceMessages_MsgOptions(isIcuTemplate bool, escapeLessThan bool, unescapeHtmlEntities bool) *ReplaceMessages_MsgOptions {
	this_ := &ReplaceMessages_MsgOptions{}
	InitReplaceMessages_MsgOptions(this_, isIcuTemplate, escapeLessThan, unescapeHtmlEntities)
	return this_
}

func InitReplaceMessages_MsgOptions(this_ *ReplaceMessages_MsgOptions, isIcuTemplate bool, escapeLessThan bool, unescapeHtmlEntities bool) {
	_ = isIcuTemplate
	_ = escapeLessThan
	_ = unescapeHtmlEntities
	this_.initFields_ReplaceMessages_MsgOptions()
	this_.isIcuTemplate = isIcuTemplate
	this_.escapeLessThan = escapeLessThan
	this_.unescapeHtmlEntities = unescapeHtmlEntities
}

func (this_ *ReplaceMessages_MsgOptions) IsIcuTemplate() bool  { return this_.isIcuTemplate }
func (this_ *ReplaceMessages_MsgOptions) EscapeLessThan() bool { return this_.escapeLessThan }
func (this_ *ReplaceMessages_MsgOptions) UnescapeHtmlEntities() bool {
	return this_.unescapeHtmlEntities
}
func (this_ *ReplaceMessages_MsgOptions) Equals(other_ any) bool {
	o_, ok_ := other_.(*ReplaceMessages_MsgOptions)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.isIcuTemplate != o_.isIcuTemplate {
		return false
	}
	if this_.escapeLessThan != o_.escapeLessThan {
		return false
	}
	if this_.unescapeHtmlEntities != o_.unescapeHtmlEntities {
		return false
	}
	return true
}
func (this_ *ReplaceMessages_MsgOptions) HashCode() int32 {
	var h_ int32 = 0
	if this_.isIcuTemplate {
		h_ = 31*h_ + 1231
	} else {
		h_ = 31*h_ + 1237
	}
	if this_.escapeLessThan {
		h_ = 31*h_ + 1231
	} else {
		h_ = 31*h_ + 1237
	}
	if this_.unescapeHtmlEntities {
		h_ = 31*h_ + 1231
	} else {
		h_ = 31*h_ + 1237
	}
	return h_
}
func (this_ *ReplaceMessages_MsgOptions) String() string {
	_ = fmt.Sprint
	return "MsgOptions[" + "isIcuTemplate=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.isIcuTemplate)) + ", " + "escapeLessThan=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.escapeLessThan)) + ", " + "unescapeHtmlEntities=" + j2gort.StrVal(j2gort.ValueOfBoolean(this_.unescapeHtmlEntities)) + "]"
}

// Java class: com.google.javascript.jscomp.ReplaceMessages.ProtectedJsMessage
type ReplaceMessages_ProtectedJsMessage struct {
	jsMessage            *JsMessage     `gson:"jsMessage"`
	definitionNode       pkg_rhino.Node `gson:"definitionNode"`
	substitutionsNode    pkg_rhino.Node `gson:"substitutionsNode"`
	idHash_              atomic.Int32
	isIcuTemplate        bool `gson:"isIcuTemplate"`
	escapeLessThan       bool `gson:"escapeLessThan"`
	unescapeHtmlEntities bool `gson:"unescapeHtmlEntities"`
}

func (this_ *ReplaceMessages_ProtectedJsMessage) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceMessages_ProtectedJsMessage) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceMessages_ProtectedJsMessage) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceMessages$ProtectedJsMessage", "ProtectedJsMessage"
}
func (this_ *ReplaceMessages_ProtectedJsMessage) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReplaceMessages_ProtectedJsMessage) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceMessages.ProtectedJsMessage", "com.google.javascript.jscomp.ReplaceMessages$ProtectedJsMessage", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceMessages_ProtectedJsMessage](v_)
	})
}

func (this_ *ReplaceMessages_ProtectedJsMessage) initFields_ReplaceMessages_ProtectedJsMessage() {
	_ = this_
}

// Java constructor: ProtectedJsMessage(com.google.javascript.jscomp.JsMessage jsMessage, com.google.javascript.rhino.Node definitionNode, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node substitutionsNode, boolean isIcuTemplate, boolean escapeLessThan, boolean unescapeHtmlEntities)
func NewReplaceMessages_ProtectedJsMessage(jsMessage *JsMessage, definitionNode pkg_rhino.Node, substitutionsNode pkg_rhino.Node, isIcuTemplate bool, escapeLessThan bool, unescapeHtmlEntities bool) *ReplaceMessages_ProtectedJsMessage {
	this_ := &ReplaceMessages_ProtectedJsMessage{}
	InitReplaceMessages_ProtectedJsMessage(this_, jsMessage, definitionNode, substitutionsNode, isIcuTemplate, escapeLessThan, unescapeHtmlEntities)
	return this_
}

func InitReplaceMessages_ProtectedJsMessage(this_ *ReplaceMessages_ProtectedJsMessage, jsMessage *JsMessage, definitionNode pkg_rhino.Node, substitutionsNode pkg_rhino.Node, isIcuTemplate bool, escapeLessThan bool, unescapeHtmlEntities bool) {
	_ = jsMessage
	_ = definitionNode
	_ = substitutionsNode
	_ = isIcuTemplate
	_ = escapeLessThan
	_ = unescapeHtmlEntities
	this_.initFields_ReplaceMessages_ProtectedJsMessage()
	this_.jsMessage = jsMessage
	this_.definitionNode = definitionNode
	this_.substitutionsNode = substitutionsNode
	this_.isIcuTemplate = isIcuTemplate
	this_.escapeLessThan = escapeLessThan
	this_.unescapeHtmlEntities = unescapeHtmlEntities
}

// Java static method: ProtectedJsMessage.fromAstNode(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.JsMessage.IdGenerator idGenerator)
func ReplaceMessages_ProtectedJsMessage_FromAstNode(node pkg_rhino.Node, idGenerator JsMessage_IdGenerator) *ReplaceMessages_ProtectedJsMessage {
	_ = node
	_ = idGenerator
	if !(node.IsCall()) {
		return nil
	}
	var calleeNode pkg_rhino.Node = j2gort.CheckNotNullNode(node.GetFirstChild())
	_ = calleeNode
	if !(calleeNode.MatchesName_1_string("__jscomp_define_msg__")) {
		return nil
	}
	var propertiesNode pkg_rhino.Node = j2gort.CheckNotNullNode(calleeNode.GetNext())
	_ = propertiesNode
	var substitutionsNode pkg_rhino.Node = propertiesNode.GetNext()
	_ = substitutionsNode
	var escapeLessThanOption bool = false
	_ = escapeLessThanOption
	var unescapeHtmlEntitiesOption bool = false
	_ = unescapeHtmlEntitiesOption
	var isIcuTemplate bool = false
	_ = isIcuTemplate
	var val_24 JsMessage_Builder
	InitJsMessage_Builder(&val_24)
	var jsMessageBuilder *JsMessage_Builder = &val_24
	_ = jsMessageBuilder
	if !(propertiesNode.IsObjectLit()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(propertiesNode))
	}
	var msgKey string = j2gort.NullString
	_ = msgKey
	var meaning string = j2gort.NullString
	_ = meaning
	var icuPlaceholderNames j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_LinkedHashSet()
	_ = icuPlaceholderNames
	var messageText string = j2gort.NullString
	_ = messageText
	var messageTextNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = messageTextNode
	{
		var strKey pkg_rhino.Node = propertiesNode.GetFirstChild()
		_ = strKey
		for ; strKey != 0; strKey = strKey.GetNext() {
			if !(strKey.IsStringKey()) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(strKey))
			}
			var key string = strKey.GetString()
			_ = key
			var valueNode pkg_rhino.Node = strKey.GetOnlyChild()
			_ = valueNode
			if j2gort.CheckNotNullString(key) == "icu_placeholder_names" {
				if !(valueNode.IsArrayLit()) {
					j2gort.CheckStateStr(false, "icu_placeholder_names must be an array")
				}
				{
					cur_25 := j2gort.NewIterCursor(j2gort.CheckNotNull(valueNode.Children()))
					for cur_25.HasNext() {
						valueNodeChild := pkg_rhino.Node(cur_25.NextNode())
						_ = valueNodeChild
						j2gort.SetAddStr(icuPlaceholderNames, valueNodeChild.GetString())
					}
				}
				continue
			}
			if !(valueNode.IsStringLit()) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(valueNode))
			}
			var value string = valueNode.GetString()
			_ = value
			switch j2gort.CheckNotNullString(key) {
			case "key":
				jsMessageBuilder.SetKey(value)
				msgKey = value
			case "meaning":
				jsMessageBuilder.SetMeaning(value)
				meaning = value
			case "alt_id":
				jsMessageBuilder.SetAlternateId(value)
			case "msg_text":
				messageText = value
				messageTextNode = valueNode
			case "isIcuTemplate":
				isIcuTemplate = true
			case "escapeLessThan":
				escapeLessThanOption = true
			case "unescapeHtmlEntities":
				unescapeHtmlEntitiesOption = true
			default:
				panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "unknown protected message key: "+j2gort.StrVal(strKey))))
			}
		}
	}
	func() {
		defer func() {
			if rec_26 := recover(); rec_26 != nil {
				jex_27 := javalang.RecoverException(rec_26)
				_ = jex_27
				if j2gort.InstanceOf[*JsMessage_PlaceholderFormatException](jex_27) {
					unused := j2gort.AsIface[*JsMessage_PlaceholderFormatException](jex_27)
					_ = unused
					panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString((j2gort.StrValString(messageTextNode.GetLocation()) + ": Placeholder incorrectly formatted: >" + j2gort.StrValString(messageText) + "<")))))
				} else {
					panic(rec_26)
				}
			}
		}()
		if !(icuPlaceholderNames.IsEmpty()) {
			if !(isIcuTemplate) {
				j2gort.CheckStateStr(false, "Found icu_placeholder_names for a message that is not an ICU template.")
			}
			var val_28 JsMessageVisitor_IcuMessageTemplateString
			InitJsMessageVisitor_IcuMessageTemplateString(&val_28, messageText)
			var icuMessageTemplateString *JsMessageVisitor_IcuMessageTemplateString = &val_28
			_ = icuMessageTemplateString
			var extractedIcuTemplateParts *JsMessageVisitor_ExtractedIcuTemplateParts = icuMessageTemplateString.ExtractParts(icuPlaceholderNames)
			_ = extractedIcuTemplateParts
			jsMessageBuilder.AppendParts(j2gort.AsIface[j2gort.Ext_java_util_List](extractedIcuTemplateParts.extractedParts))
		} else {
			jsMessageBuilder.AppendParts(j2gort.AsIface[j2gort.Ext_java_util_List](JsMessageVisitor_ParseJsMessageTextIntoParts(messageText)))
		}
	}()
	var externalMessageId string = JsMessageVisitor_GetExternalMessageId(msgKey)
	_ = externalMessageId
	if externalMessageId != j2gort.NullString {
		jsMessageBuilder.SetIsExternalMsg(true).SetId(externalMessageId)
	} else {
		var meaningForIdGeneration string
		_ = meaningForIdGeneration
		if meaning != j2gort.NullString {
			meaningForIdGeneration = meaning
		} else {
			meaningForIdGeneration = JsMessageVisitor_RemoveScopedAliasesPrefix(msgKey)
		}
		if idGenerator != nil && !idGenerator.IsJavaNil_() {
			jsMessageBuilder.SetId(idGenerator.GenerateId(meaningForIdGeneration, jsMessageBuilder.GetParts()))
		} else {
			jsMessageBuilder.SetId(meaningForIdGeneration)
		}
	}
	return NewReplaceMessages_ProtectedJsMessage(jsMessageBuilder.Build(), node, substitutionsNode, isIcuTemplate, escapeLessThanOption, unescapeHtmlEntitiesOption)
}

// Java method: ProtectedJsMessage.getMsgOptions()
func (this_ *ReplaceMessages_ProtectedJsMessage) GetMsgOptions() *ReplaceMessages_MsgOptions {
	return NewReplaceMessages_MsgOptions(this_.isIcuTemplate, this_.escapeLessThan, this_.unescapeHtmlEntities)
}

// Java class:
type ReplaceMessages__Anon_77 struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *ReplaceMessages_ReplacementCompletionPass
}

func (this_ *ReplaceMessages__Anon_77) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceMessages__Anon_77) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceMessages__Anon_77) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceMessages$ReplacementCompletionPass$1", ""
}
func (this_ *ReplaceMessages__Anon_77) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*ReplaceMessages__Anon_77) HasIdentityEquals_() {}
func (this_ *ReplaceMessages__Anon_77) SetOuter_(o *ReplaceMessages_ReplacementCompletionPass) {
	this_.outer_ = o
}
func (this_ *ReplaceMessages__Anon_77) GetOuter_() *ReplaceMessages_ReplacementCompletionPass {
	return this_.outer_
}
func (*ReplaceMessages__Anon_77) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*ReplaceMessages__Anon_77) Implements_jscomp_NodeTraversal_Callback()                  {}

func (this_ *ReplaceMessages__Anon_77) initFields_ReplaceMessages__Anon_77() {
	_ = this_
}

func NewReplaceMessages__Anon_77(outer_ *ReplaceMessages_ReplacementCompletionPass) *ReplaceMessages__Anon_77 {
	this_ := &ReplaceMessages__Anon_77{}
	this_.outer_ = outer_
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_ReplaceMessages__Anon_77()
	return this_
}

// Java method: .visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ReplaceMessages__Anon_77) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var protectedJsMessage *ReplaceMessages_ProtectedJsMessage = ReplaceMessages_ProtectedJsMessage_FromAstNode(n, this_.outer_.outer_.bundle.IdGenerator())
	_ = protectedJsMessage
	if protectedJsMessage != nil {
		this_.outer_.VisitMsgDefinition(protectedJsMessage)
	} else {
		var protectedMsgFallback *ReplaceMessages_ProtectedMsgFallback = ReplaceMessages_ProtectedMsgFallback_FromAstNode(n)
		_ = protectedMsgFallback
		if protectedMsgFallback != nil {
			this_.outer_.priv_ReplacementCompletionPass_visitMsgFallback_4880(protectedMsgFallback)
		}
	}
}
func (*ReplaceMessages_MsgProtectionPass) AlwaysShouldTraverse_()   {}
func (*ReplaceMessages_FullReplacementPass) AlwaysShouldTraverse_() {}
func (*ReplaceMessages__Anon_77) AlwaysShouldTraverse_()            {}
func (*ReplaceMessages__Anon_77) CanReleaseSyntacticScopes_()       {}
