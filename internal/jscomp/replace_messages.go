// Port of com/google/javascript/jscomp/ReplaceMessages.java @ closure-compiler v20261021

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ReplaceMessages replaces user-visible messages with alternatives. It uses
// Google specific JsMessageVisitor implementation.
//
// Java: public final class ReplaceMessages
type ReplaceMessages struct {
	compiler                 AbstractCompiler
	bundle                   MessageBundle
	strictReplacement        bool
	collapsePropertiesHasRun bool
	astFactory               *AstFactory
}

// Java: ReplaceMessages.BUNDLE_DOES_NOT_HAVE_THE_MESSAGE
var ReplaceMessages_BUNDLE_DOES_NOT_HAVE_THE_MESSAGE = DiagnosticType_Error(
	"JSC_BUNDLE_DOES_NOT_HAVE_THE_MESSAGE",
	"Message with id = {0} could not be found in replacement bundle")

// Java: ReplaceMessages.INVALID_ALTERNATE_MESSAGE_PLACEHOLDERS
var ReplaceMessages_INVALID_ALTERNATE_MESSAGE_PLACEHOLDERS = DiagnosticType_Error(
	"JSC_INVALID_ALTERNATE_MESSAGE_PLACEHOLDERS",
	"Alternate message ID={0} placeholders ({1}) differs from {2} placeholders ({3}).")

// Java: ReplaceMessages(AbstractCompiler, MessageBundle, boolean)
func NewReplaceMessages(compiler AbstractCompiler, bundle MessageBundle, strictReplacement bool) *ReplaceMessages {
	r := &ReplaceMessages{}
	r.compiler = compiler
	r.astFactory = compiler.CreateAstFactory()
	r.bundle = bundle
	r.strictReplacement = strictReplacement
	r.collapsePropertiesHasRun =
		compiler.GetOptions().GetPropertyCollapseLevel() == CompilerOptions_PropertyCollapseLevel_ALL &&
			compiler.GetOptions().DoLateLocalization()
	return r
}

// GetMsgProtectionPass: when the returned pass is executed, the original
// `goog.getMsg()` etc. calls will be replaced with a form that will survive
// unchanged through optimizations unless eliminated as unused.
//
// After all optimizations are complete, the pass returned by
// `getReplacementCompletionPass()`.
//
// Java: ReplaceMessages.getMsgProtectionPass()
func (r *ReplaceMessages) GetMsgProtectionPass() CompilerPass {
	return NewReplaceMessages_MsgProtectionPass(r)
}

// Java: record ReplaceMessages.MsgProtectionData
type ReplaceMessages_MsgProtectionData struct {
	message               *JsMessage
	messageNode           *rhino.Node
	templateTextNode      *rhino.Node
	placeholderValuesNode *rhino.Node // may be nil
	messageOptions        *ReplaceMessages_MsgOptions
}

// Java: ReplaceMessages.MsgProtectionData(JsMessage, Node, Node, Node, ReplaceMessages.MsgOptions)
func NewReplaceMessages_MsgProtectionData(
	message *JsMessage,
	messageNode *rhino.Node,
	templateTextNode *rhino.Node,
	placeholderValuesNode *rhino.Node,
	messageOptions *ReplaceMessages_MsgOptions) *ReplaceMessages_MsgProtectionData {
	return &ReplaceMessages_MsgProtectionData{
		message:               message,
		messageNode:           messageNode,
		templateTextNode:      templateTextNode,
		placeholderValuesNode: placeholderValuesNode,
		messageOptions:        messageOptions,
	}
}

// Java: ReplaceMessages.MsgProtectionData.message()
func (d *ReplaceMessages_MsgProtectionData) Message() *JsMessage { return d.message }

// Java: ReplaceMessages.MsgProtectionData.messageNode()
func (d *ReplaceMessages_MsgProtectionData) MessageNode() *rhino.Node { return d.messageNode }

// Java: ReplaceMessages.MsgProtectionData.templateTextNode()
func (d *ReplaceMessages_MsgProtectionData) TemplateTextNode() *rhino.Node {
	return d.templateTextNode
}

// Java: ReplaceMessages.MsgProtectionData.placeholderValuesNode()
func (d *ReplaceMessages_MsgProtectionData) PlaceholderValuesNode() *rhino.Node {
	return d.placeholderValuesNode
}

// Java: ReplaceMessages.MsgProtectionData.messageOptions()
func (d *ReplaceMessages_MsgProtectionData) MessageOptions() *ReplaceMessages_MsgOptions {
	return d.messageOptions
}

// Java: class ReplaceMessages.MsgProtectionPass extends JsMessageVisitor
type ReplaceMessages_MsgProtectionPass struct {
	JsMessageVisitorBase
	outer *ReplaceMessages
}

// Java: ReplaceMessages.MsgProtectionPass()
func NewReplaceMessages_MsgProtectionPass(outer *ReplaceMessages) *ReplaceMessages_MsgProtectionPass {
	p := &ReplaceMessages_MsgProtectionPass{outer: outer}
	p.InitJsMessageVisitorBase(p, outer.compiler, outer.bundle.IdGenerator())
	return p
}

// Java: ReplaceMessages.MsgProtectionPass.process(Node, Node)
func (p *ReplaceMessages_MsgProtectionPass) Process(externs *rhino.Node, root *rhino.Node) {
	compiler := p.outer.compiler
	// Add externs declarations for the function names we use in our replacements.
	NodeUtil_CreateSynthesizedExternsSymbol(compiler, ReplaceMessagesConstants_DEFINE_MSG_CALLEE)
	NodeUtil_CreateSynthesizedExternsSymbol(
		compiler, ReplaceMessagesConstants_FALLBACK_MSG_CALLEE)

	// JsMessageVisitor.process() does the traversal that calls the processX() methods below.
	p.JsMessageVisitorBase.Process(externs, root)
}

// Java: ReplaceMessages.MsgProtectionPass.processIcuTemplateDefinition(IcuTemplateDefinition)
func (p *ReplaceMessages_MsgProtectionPass) ProcessIcuTemplateDefinition(definition IcuTemplateDefinition) {
	p.PerformMessageProtection(
		NewReplaceMessages_MsgProtectionData(
			/* message= */ definition.GetMessage(),
			/* messageNode= */ definition.GetMessageNode(),
			/* templateTextNode= */ definition.GetTemplateTextNode(),
			// There are no compile-time placeholder replacements for ICU templates
			/* placeholderValuesNode= */
			nil,
			/* messageOptions= */ replaceMessages_ICU_MSG_OPTIONS))
}

// Java: ReplaceMessages.MsgProtectionPass.processJsMessageDefinition(JsMessageDefinition)
func (p *ReplaceMessages_MsgProtectionPass) ProcessJsMessageDefinition(definition JsMessageDefinition) {
	// This is the currently preferred form.
	// `MSG_A = goog.getMsg('hello, {$name}', {name: getName()}, {html: true})`
	p.PerformMessageProtection(
		NewReplaceMessages_MsgProtectionData(
			/* message= */ definition.GetMessage(),
			/* messageNode= */ definition.GetMessageNode(),
			/* templateTextNode= */ definition.GetTemplateTextNode(),
			/* placeholderValuesNode= */ definition.GetPlaceholderValuesNode(),
			/* messageOptions= */ p.outer.GetMsgOptionsFromDefinition(definition)))
}

// Java: ReplaceMessages.MsgProtectionPass.performMessageProtection(ReplaceMessages.MsgProtectionData)
func (p *ReplaceMessages_MsgProtectionPass) PerformMessageProtection(msgProtectionData *ReplaceMessages_MsgProtectionData) {
	r := p.outer
	message := msgProtectionData.Message()
	callNode := msgProtectionData.MessageNode()
	originalMessageString := msgProtectionData.TemplateTextNode()
	placeholdersNode := msgProtectionData.PlaceholderValuesNode()
	msgOptions := msgProtectionData.MessageOptions()

	jutil.CheckStateMsg1(callNode.IsCall(), "%s", callNode)
	// `goog.getMsg('message string', {<substitutions>}, {<options>})`
	googGetMsg := callNode.GetFirstChild()

	// Construct
	// `__jscomp_define_msg__({<msg properties>}, {<substitutions>})`
	protectionFunctionName := ReplaceMessagesConstants_DEFINE_MSG_CALLEE
	newCallee := p.CreateProtectionFunctionCallee(protectionFunctionName).Srcref(googGetMsg)
	msgPropertiesNode := r.CreateMsgPropertiesNode(message, msgOptions).SrcrefTree(originalMessageString)
	newCallNode := r.astFactory.CreateCall(newCallee, AstFactory_Type_(callNode), msgPropertiesNode).Srcref(callNode)
	// If the result of this call (the message) is unused, there is no reason for optimizations
	// to preserve it.
	newCallNode.SetSideEffectFlags(rhino.Node_SideEffectFlags_NO_SIDE_EFFECTS)
	if placeholdersNode != nil {
		jutil.CheckStateMsg1(placeholdersNode.IsObjectLit(), "%s", placeholdersNode)
		// put quotes around the keys so they won't get renamed.
		for strKey := placeholdersNode.GetFirstChild(); strKey != nil; strKey = strKey.GetNext() {
			jutil.CheckStateMsg1(strKey.IsStringKey(), "%s", strKey)
			strKey.SetQuotedStringKey()
		}
		newCallNode.AddChildToBack(placeholdersNode.Detach())
	}
	callNode.ReplaceWith(newCallNode)
	r.compiler.ReportChangeToEnclosingScope(newCallNode)
}

// Java: ReplaceMessages.MsgProtectionPass.createProtectionFunctionCallee(String)
func (p *ReplaceMessages_MsgProtectionPass) CreateProtectionFunctionCallee(protectionFunctionName string) *rhino.Node {
	callee := p.outer.astFactory.CreateNameWithUnknownType(protectionFunctionName)
	// The name is declared constant in the externs definition we created, so all references
	// to it must also be marked as constant for consistency's sake.
	callee.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	return callee
}

// Java: ReplaceMessages.MsgProtectionPass.processMessageFallback(Node, JsMessage, JsMessage)
func (p *ReplaceMessages_MsgProtectionPass) ProcessMessageFallback(callNode *rhino.Node, message1 *JsMessage, message2 *JsMessage) {
	r := p.outer
	originalCallee := jutil.CheckNotNullMsg1(callNode.GetFirstChild(), "%s", callNode)
	fallbackCallee := p.CreateProtectionFunctionCallee(ReplaceMessagesConstants_FALLBACK_MSG_CALLEE).
		Srcref(originalCallee)

	originalFirstArg := jutil.CheckNotNullMsg1(originalCallee.GetNext(), "%s", callNode)
	firstMsgKey := r.astFactory.CreateString(message1.GetKey()).Srcref(originalFirstArg)

	originalSecondArg := jutil.CheckNotNullMsg1(originalFirstArg.GetNext(), "%s", callNode)
	secondMsgKey := r.astFactory.CreateString(message2.GetKey()).Srcref(originalSecondArg)

	// `__jscomp_msg_fallback__('MSG_ONE', MSG_ONE, 'MSG_TWO', MSG_TWO)`
	newCallNode := r.astFactory.
		CreateCall(
			fallbackCallee,
			AstFactory_Type_(callNode),
			firstMsgKey,
			originalFirstArg.Detach(),
			secondMsgKey,
			originalSecondArg.Detach()).
		Srcref(callNode)
	// If the result of this call (the message) is unused, there is no reason for optimizations
	// to preserve it.
	newCallNode.SetSideEffectFlags(rhino.Node_SideEffectFlags_NO_SIDE_EFFECTS)
	callNode.ReplaceWith(newCallNode)
	r.compiler.ReportChangeToEnclosingScope(newCallNode)
}

// Java: ReplaceMessages.getMsgOptionsFromDefinition(JsMessageDefinition)
func (r *ReplaceMessages) GetMsgOptionsFromDefinition(definition JsMessageDefinition) *ReplaceMessages_MsgOptions {
	return NewReplaceMessages_MsgOptions(
		/* isIcuTemplate= */ false,
		/* escapeLessThan= */ definition.ShouldEscapeLessThan(),
		/* unescapeHtmlEntities= */ definition.ShouldUnescapeHtmlEntities())
}

// Java: ReplaceMessages.createMsgPropertiesNode(JsMessage, ReplaceMessages.MsgOptions)
func (r *ReplaceMessages) CreateMsgPropertiesNode(message *JsMessage, msgOptions *ReplaceMessages_MsgOptions) *rhino.Node {
	msgPropsBuilder := NewReplaceMessages_QuotedKeyObjectLitBuilder(r)
	msgPropsBuilder.AddString("key", message.GetKey())
	if msgOptions.IsIcuTemplate() && !message.CanonicalPlaceholderNames().IsEmpty() {
		// ICU messages created using `declareIcuTemplate` can get stored into the XMB file as
		// multiple parts if necessary to record example or original code text.
		// `icu_placeholder_names` stores these parts of the ICU message, which allows us to
		// correctly calculate the message ID in the protected message.
		namesArrayLit := r.astFactory.CreateArraylit()
		for name := range message.CanonicalPlaceholderNames().All() {
			namesArrayLit.AddChildToBack(r.astFactory.CreateString(name))
		}
		// Example:
		// declareIcuTemplate('blah blah {PH1} blah {PH2}', ... );
		// icu_placeholder_names: ['PH1', 'PH2']
		msgPropsBuilder.AddNode("icu_placeholder_names", namesArrayLit)
	}
	altId := message.GetAlternateId()
	if altId != nil {
		msgPropsBuilder.AddString("alt_id", *altId)
	}
	meaning := message.GetMeaning()
	if meaning != nil {
		msgPropsBuilder.AddString("meaning", *meaning)
	}
	if msgOptions.IsIcuTemplate() {
		msgPropsBuilder.AddString("msg_text", message.AsIcuMessageString())
		// Just being present is what records this option as true
		msgPropsBuilder.AddString("isIcuTemplate", "")
	} else {
		msgPropsBuilder.AddString("msg_text", message.AsJsMessageString())
	}
	if msgOptions.EscapeLessThan() {
		// Just being present is what records this option as true
		msgPropsBuilder.AddString("escapeLessThan", "")
	}
	if msgOptions.UnescapeHtmlEntities() {
		// Just being present is what records this option as true
		msgPropsBuilder.AddString("unescapeHtmlEntities", "")
	}
	return msgPropsBuilder.Build()
}

// Java: private final class ReplaceMessages.QuotedKeyObjectLitBuilder
type ReplaceMessages_QuotedKeyObjectLitBuilder struct {
	outer *ReplaceMessages
	// LinkedHashMap to keep the keys in the order we set them so our output is deterministic.
	keyToValueNodeMap *jutil.LinkedHashMap[string, *rhino.Node]
}

// Java: ReplaceMessages.QuotedKeyObjectLitBuilder()
func NewReplaceMessages_QuotedKeyObjectLitBuilder(outer *ReplaceMessages) *ReplaceMessages_QuotedKeyObjectLitBuilder {
	return &ReplaceMessages_QuotedKeyObjectLitBuilder{
		outer:             outer,
		keyToValueNodeMap: jutil.NewLinkedHashMap[string, *rhino.Node](),
	}
}

// Java: ReplaceMessages.QuotedKeyObjectLitBuilder.addString(String, String)
func (b *ReplaceMessages_QuotedKeyObjectLitBuilder) AddString(key string, value string) *ReplaceMessages_QuotedKeyObjectLitBuilder {
	return b.AddNode(key, b.outer.astFactory.CreateString(value))
}

// Java: ReplaceMessages.QuotedKeyObjectLitBuilder.addNode(String, Node)
func (b *ReplaceMessages_QuotedKeyObjectLitBuilder) AddNode(key string, node *rhino.Node) *ReplaceMessages_QuotedKeyObjectLitBuilder {
	jutil.CheckStateMsg1(!b.keyToValueNodeMap.ContainsKey(key), "repeated key: %s", key)
	b.keyToValueNodeMap.Put(key, node)
	return b
}

// Java: ReplaceMessages.QuotedKeyObjectLitBuilder.build()
func (b *ReplaceMessages_QuotedKeyObjectLitBuilder) Build() *rhino.Node {
	result := b.outer.astFactory.CreateObjectLit()
	for key, value := range b.keyToValueNodeMap.All() {
		result.AddChildToBack(b.outer.astFactory.CreateQuotedStringKey(key, value))
	}
	return result
}

// GetReplacementCompletionPass: when the returned pass is executed, the
// protected messages created by `getMsgProtectionPass()` will be replaced by
// the final message form read from the appropriate message bundle.
//
// Java: ReplaceMessages.getReplacementCompletionPass()
func (r *ReplaceMessages) GetReplacementCompletionPass() CompilerPass {
	return NewReplaceMessages_ReplacementCompletionPass(r)
}

// Java: class ReplaceMessages.ReplacementCompletionPass implements CompilerPass
type ReplaceMessages_ReplacementCompletionPass struct {
	outer *ReplaceMessages
	// Keep track of which messages actually got translated, so we know what do do when we
	// see a message fallback call.
	translatedMsgKeys *jutil.LinkedHashSet[string]
}

// Java: ReplaceMessages.ReplacementCompletionPass()
func NewReplaceMessages_ReplacementCompletionPass(outer *ReplaceMessages) *ReplaceMessages_ReplacementCompletionPass {
	return &ReplaceMessages_ReplacementCompletionPass{
		outer:             outer,
		translatedMsgKeys: jutil.NewLinkedHashSet[string](),
	}
}

// replaceMessages_replacementCompletionCallback is the anonymous
// AbstractPostOrderCallback of ReplacementCompletionPass.process.
type replaceMessages_replacementCompletionCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	pass *ReplaceMessages_ReplacementCompletionPass
}

// Java: ReplaceMessages.ReplacementCompletionPass.process(...).new AbstractPostOrderCallback() {...}.visit(NodeTraversal, Node, Node)
func (c *replaceMessages_replacementCompletionCallback) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	protectedJsMessage := ReplaceMessages_ProtectedJsMessage_FromAstNode(n, c.pass.outer.bundle.IdGenerator())
	if protectedJsMessage != nil {
		c.pass.VisitMsgDefinition(protectedJsMessage)
	} else {
		protectedMsgFallback := ReplaceMessages_ProtectedMsgFallback_FromAstNode(n)
		if protectedMsgFallback != nil {
			c.pass.VisitMsgFallback(protectedMsgFallback)
		}
	}
}

// Java: ReplaceMessages.ReplacementCompletionPass.process(Node, Node)
func (p *ReplaceMessages_ReplacementCompletionPass) Process(externs *rhino.Node, root *rhino.Node) {
	// for each `__jscomp_define_msg__` call in post-order traversal
	// replace it with the appropriate expression
	NodeTraversal_Traverse(
		p.outer.compiler,
		root,
		&replaceMessages_replacementCompletionCallback{pass: p})
}

// Java: ReplaceMessages.ReplacementCompletionPass.visitMsgDefinition(ReplaceMessages.ProtectedJsMessage)
func (p *ReplaceMessages_ReplacementCompletionPass) VisitMsgDefinition(protectedJsMessage *ReplaceMessages_ProtectedJsMessage) {
	r := p.outer
	jutil.Try(func() {
		originalMsg := protectedJsMessage.jsMessage
		nodeToReplace := protectedJsMessage.definitionNode
		translatedMsg := r.LookupMessage(protectedJsMessage.definitionNode, r.bundle, originalMsg)
		var msgToUse *JsMessage
		if translatedMsg != nil {
			msgToUse = translatedMsg
			// Remember that this one got translated in case it is used in a fallback.
			p.translatedMsgKeys.Add(originalMsg.GetKey())
		} else {
			if r.strictReplacement {
				r.compiler.Report(
					JSError_MakeWithNode(nodeToReplace, ReplaceMessages_BUNDLE_DOES_NOT_HAVE_THE_MESSAGE, originalMsg.GetId()))
			}
			msgToUse = originalMsg
		}
		msgOptions := protectedJsMessage.GetMsgOptions()
		placeholderMap := ReplaceMessages_ExtractPlaceholderValuesMapOrThrow(protectedJsMessage.substitutionsNode)
		finalMsgConstructionExpression := r.ConstructStringExprNode(msgToUse, placeholderMap, msgOptions, nodeToReplace)
		finalMsgConstructionExpression.SrcrefTreeIfMissing(nodeToReplace)
		nodeToReplace.ReplaceWith(finalMsgConstructionExpression)
		r.compiler.ReportChangeToEnclosingScope(finalMsgConstructionExpression)
	}, func(e IJsMessageVisitor_MalformedException) {
		r.compiler.Report(JSError_MakeWithNode(e.GetNode(), JsMessageVisitor_MESSAGE_TREE_MALFORMED, e.GetMessage()))
	})
}

// Java: ReplaceMessages.ReplacementCompletionPass.visitMsgFallback(ReplaceMessages.ProtectedMsgFallback)
func (p *ReplaceMessages_ReplacementCompletionPass) VisitMsgFallback(protectedMsgFallback *ReplaceMessages_ProtectedMsgFallback) {
	var valueNodeToUse *rhino.Node
	if p.translatedMsgKeys.Contains(protectedMsgFallback.firstMsgKey) {
		// Obviously use the first message, if it is translated.
		valueNodeToUse = protectedMsgFallback.firstMsgValue
	} else if p.translatedMsgKeys.Contains(protectedMsgFallback.secondMsgKey) {
		// Fallback to the second message if it has a translation.
		valueNodeToUse = protectedMsgFallback.secondMsgValue
	} else {
		// If neither is translated, then use the first message as it is defined in the source code.
		valueNodeToUse = protectedMsgFallback.firstMsgValue
	}
	valueNodeToUse.Detach()
	protectedMsgFallback.callNode.ReplaceWith(valueNodeToUse)
	p.outer.compiler.ReportChangeToEnclosingScope(valueNodeToUse)
}

// Java: ReplaceMessages.extractPlaceholderValuesMapOrThrow(Node)
func ReplaceMessages_ExtractPlaceholderValuesMapOrThrow(valuesObjLit *rhino.Node) *jutil.LinkedHashMap[string, *rhino.Node] {
	var result *jutil.LinkedHashMap[string, *rhino.Node]
	jutil.Try(func() {
		result = JsMessageVisitor_ExtractObjectLiteralMap(valuesObjLit).ExtractAsValueMap()
	}, func(e IJsMessageVisitor_MalformedException) {
		panic(jutil.NewIllegalStateExceptionWithMessageAndCause(e.Error(), e))
	})
	return result
}

// Java: private static class ReplaceMessages.ProtectedMsgFallback
type ReplaceMessages_ProtectedMsgFallback struct {
	callNode       *rhino.Node
	firstMsgKey    string
	firstMsgValue  *rhino.Node
	secondMsgKey   string
	secondMsgValue *rhino.Node
}

// Java: ReplaceMessages.ProtectedMsgFallback(Node, String, Node, String, Node)
func NewReplaceMessages_ProtectedMsgFallback(
	callNode *rhino.Node,
	firstMsgKey string,
	firstMsgValue *rhino.Node,
	secondMsgKey string,
	secondMsgValue *rhino.Node) *ReplaceMessages_ProtectedMsgFallback {
	return &ReplaceMessages_ProtectedMsgFallback{
		callNode:       callNode,
		firstMsgKey:    firstMsgKey,
		firstMsgValue:  firstMsgValue,
		secondMsgKey:   secondMsgKey,
		secondMsgValue: secondMsgValue,
	}
}

// Java: ReplaceMessages.ProtectedMsgFallback.fromAstNode(Node)
func ReplaceMessages_ProtectedMsgFallback_FromAstNode(n *rhino.Node) *ReplaceMessages_ProtectedMsgFallback {
	if !n.IsCall() {
		return nil
	}
	callee := n.GetFirstChild()
	if !callee.MatchesName(ReplaceMessagesConstants_FALLBACK_MSG_CALLEE) {
		return nil
	}
	jutil.CheckStateMsg1(n.HasXChildren(5), "bad message fallback call: %s", n)
	firstMsgKeyNode := callee.GetNext()
	firstMsgKey := firstMsgKeyNode.GetString()
	firstMsgValue := firstMsgKeyNode.GetNext()
	secondMsgKeyNode := firstMsgValue.GetNext()
	secondMsgKey := secondMsgKeyNode.GetString()
	secondMsgValue := secondMsgKeyNode.GetNext()
	return NewReplaceMessages_ProtectedMsgFallback(n, firstMsgKey, firstMsgValue, secondMsgKey, secondMsgValue)
}

// Java: ReplaceMessages.lookupMessage(Node, MessageBundle, JsMessage)
func (r *ReplaceMessages) LookupMessage(
	callNode *rhino.Node, bundle MessageBundle, message *JsMessage) *JsMessage {
	translatedMessage := bundle.GetMessage(message.GetId())
	if translatedMessage != nil {
		return translatedMessage
	}

	alternateId := message.GetAlternateId()
	if alternateId == nil {
		return nil
	}

	alternateMessage := bundle.GetMessage(*alternateId)
	if alternateMessage != nil {
		// Validate that the alternate message is compatible with this message. Ideally we'd also
		// check meaning and description, but they're not populated by `MessageBundle.getMessage`.
		jsCodePlaceholderNames := message.JsPlaceholderNames()
		alternateMsgPlaceholderNames := alternateMessage.JsPlaceholderNames()
		if !jutil.ObjectsEquals(jsCodePlaceholderNames, alternateMsgPlaceholderNames) {
			// The JS code definition has no placeholders, but the message we got from the translation
			// bundle does have placeholders.
			//
			// This can happen for ICU - formatted messages. They can contain placeholders in the
			// message string in the form "{PH_NAME}", which this compiler ignores, because they will
			// be replaced at runtime by passing the whole string to a message formatter method.
			// However, sometimes the translation bundle will represent the "{PH_NAME}" substring as
			// an actual placeholder reference, because it was generated by some tool other than
			// closure-compiler. In that case, we need to join all the parts of the message back
			// together as a single string with the "{PH_NAME}" references back in place.
			//
			// Messages with this form always start with a particular formula, which we can test for
			// here to make sure this isn't actually a case of something being broken in the translation
			// pipeline.
			if jsCodePlaceholderNames.IsEmpty() && ReplaceMessages_IsStartOfIcuMessage(message.AsIcuMessageString()) {
				return alternateMessage
			} else {
				r.compiler.Report(
					JSError_MakeWithNode(
						callNode,
						ReplaceMessages_INVALID_ALTERNATE_MESSAGE_PLACEHOLDERS,
						*alternateId,
						jutil.ToJavaString(alternateMsgPlaceholderNames),
						message.GetKey(),
						jutil.ToJavaString(jsCodePlaceholderNames)))
				return nil
			}
		}
	}
	return alternateMessage
}

// GetFullReplacementPass: when the returned pass is executed, the original
// `goog.getMsg()` etc. calls will all be replaced with the final message
// form read from the message bundle.
//
// This is the original way of running this pass as a single operation.
//
// Java: ReplaceMessages.getFullReplacementPass()
func (r *ReplaceMessages) GetFullReplacementPass() CompilerPass {
	return NewReplaceMessages_FullReplacementPass(r)
}

// Java: record ReplaceMessages.FullReplacementMsgData
type ReplaceMessages_FullReplacementMsgData struct {
	message             *JsMessage
	messageNode         *rhino.Node
	messageOptions      *ReplaceMessages_MsgOptions
	placeholderValueMap *jutil.LinkedHashMap[string, *rhino.Node]
}

// Java: ReplaceMessages.FullReplacementMsgData(JsMessage, Node, ReplaceMessages.MsgOptions, ImmutableMap<String, Node>)
func NewReplaceMessages_FullReplacementMsgData(
	message *JsMessage,
	messageNode *rhino.Node,
	messageOptions *ReplaceMessages_MsgOptions,
	placeholderValueMap *jutil.LinkedHashMap[string, *rhino.Node]) *ReplaceMessages_FullReplacementMsgData {
	return &ReplaceMessages_FullReplacementMsgData{
		message:             message,
		messageNode:         messageNode,
		messageOptions:      messageOptions,
		placeholderValueMap: placeholderValueMap,
	}
}

// Java: ReplaceMessages.FullReplacementMsgData.message()
func (d *ReplaceMessages_FullReplacementMsgData) Message() *JsMessage { return d.message }

// Java: ReplaceMessages.FullReplacementMsgData.messageNode()
func (d *ReplaceMessages_FullReplacementMsgData) MessageNode() *rhino.Node { return d.messageNode }

// Java: ReplaceMessages.FullReplacementMsgData.messageOptions()
func (d *ReplaceMessages_FullReplacementMsgData) MessageOptions() *ReplaceMessages_MsgOptions {
	return d.messageOptions
}

// Java: ReplaceMessages.FullReplacementMsgData.placeholderValueMap()
func (d *ReplaceMessages_FullReplacementMsgData) PlaceholderValueMap() *jutil.LinkedHashMap[string, *rhino.Node] {
	return d.placeholderValueMap
}

// Java: class ReplaceMessages.FullReplacementPass extends JsMessageVisitor
type ReplaceMessages_FullReplacementPass struct {
	JsMessageVisitorBase
	outer *ReplaceMessages
}

// Java: ReplaceMessages.FullReplacementPass()
func NewReplaceMessages_FullReplacementPass(outer *ReplaceMessages) *ReplaceMessages_FullReplacementPass {
	p := &ReplaceMessages_FullReplacementPass{outer: outer}
	p.InitJsMessageVisitorBase(p, outer.compiler, outer.bundle.IdGenerator())
	return p
}

// Java: ReplaceMessages.FullReplacementPass.processMessageFallback(Node, JsMessage, JsMessage)
func (p *ReplaceMessages_FullReplacementPass) ProcessMessageFallback(callNode *rhino.Node, message1 *JsMessage, message2 *JsMessage) {
	r := p.outer
	isFirstMessageTranslated := r.LookupMessage(callNode, r.bundle, message1) != nil
	isSecondMessageTranslated := r.LookupMessage(callNode, r.bundle, message2) != nil
	var replacementNode *rhino.Node
	if isSecondMessageTranslated && !isFirstMessageTranslated {
		replacementNode = callNode.GetChildAtIndex(2)
	} else {
		replacementNode = callNode.GetSecondChild()
	}
	callNode.ReplaceWith(replacementNode.Detach())
	changeScope := ChangeTracker_GetEnclosingChangeScopeRoot(replacementNode)
	if changeScope != nil {
		r.compiler.ReportChangeToChangeScope(changeScope)
	}
}

// Java: ReplaceMessages.FullReplacementPass.processIcuTemplateDefinition(IcuTemplateDefinition)
func (p *ReplaceMessages_FullReplacementPass) ProcessIcuTemplateDefinition(definition IcuTemplateDefinition) {
	p.ProcessFullReplacement(
		NewReplaceMessages_FullReplacementMsgData(
			/* message= */ definition.GetMessage(),
			/* messageNode= */ definition.GetMessageNode(),
			/* messageOptions= */ replaceMessages_ICU_MSG_OPTIONS,
			// There are no compile-time placeholder replacements for an ICU template message.
			/* placeholderValueMap= */
			jutil.NewLinkedHashMap[string, *rhino.Node]()))
}

// Java: ReplaceMessages.FullReplacementPass.processJsMessageDefinition(JsMessageDefinition)
func (p *ReplaceMessages_FullReplacementPass) ProcessJsMessageDefinition(definition JsMessageDefinition) {
	p.ProcessFullReplacement(
		NewReplaceMessages_FullReplacementMsgData(
			/* message= */ definition.GetMessage(),
			/* messageNode= */ definition.GetMessageNode(),
			/* messageOptions= */ p.outer.GetMsgOptionsFromDefinition(definition),
			/* placeholderValueMap= */ definition.GetPlaceholderValueMap()))
}

// Java: ReplaceMessages.FullReplacementPass.processFullReplacement(ReplaceMessages.FullReplacementMsgData)
func (p *ReplaceMessages_FullReplacementPass) ProcessFullReplacement(fullReplacementMsgData *ReplaceMessages_FullReplacementMsgData) {
	r := p.outer
	message := fullReplacementMsgData.Message()
	msgNode := fullReplacementMsgData.MessageNode()
	options := fullReplacementMsgData.MessageOptions()
	placeholderValueMap := fullReplacementMsgData.PlaceholderValueMap()

	// Get the replacement.
	replacement := r.LookupMessage(msgNode, r.bundle, message)
	if replacement == nil {
		if r.strictReplacement {
			r.compiler.Report(JSError_MakeWithNode(msgNode, ReplaceMessages_BUNDLE_DOES_NOT_HAVE_THE_MESSAGE, message.GetId()))
			// Fallback to the default message
			return
		} else {
			// In case if it is not a strict replacement we could leave original
			// message.
			replacement = message
		}
	}

	// Replace the message.
	var newValue *rhino.Node
	jutil.Try(func() {
		// Build the replacement tree.
		newValue = r.ConstructStringExprNode(replacement, placeholderValueMap, options, msgNode)
	}, func(e IJsMessageVisitor_MalformedException) {
		r.compiler.Report(JSError_MakeWithNode(e.GetNode(), JsMessageVisitor_MESSAGE_TREE_MALFORMED, e.GetMessage()))
		newValue = msgNode
	})

	if newValue != msgNode {
		newValue.SrcrefTreeIfMissing(msgNode)
		msgNode.ReplaceWith(newValue)
		r.compiler.ReportChangeToEnclosingScope(newValue)
	}
}

// CreateNodeForGenderedMsgString outputs a message with a hook expression
// that returns the correct variant based on the value of
// `goog.viewerGrammaticalGender` from XTB files with gendered messages.
//
// With placeholders:
//
//	var WELCOME_MSG = function(name) {
//	   return goog.msgKind.MASCULINE ? "Bienvenido " + name :
//	          goog.msgKind.FEMININE ? "Bienvenida " + name :
//	          goog.msgKind.NEUTER ? "Les damos la bienvenida " + name :
//	          "Les damos la bienvenida " + name;
//	}(user.getName());
//
// Panics with a JsMessageVisitor_MalformedException (Java: throws
// MalformedException).
//
// Java: ReplaceMessages.createNodeForGenderedMsgString(JsMessage, Map<String, Node>, ReplaceMessages.MsgOptions, Node)
func (r *ReplaceMessages) CreateNodeForGenderedMsgString(
	message *JsMessage, placeholderMap *jutil.LinkedHashMap[string, *rhino.Node], options *ReplaceMessages_MsgOptions, nodeToReplace *rhino.Node) *rhino.Node {

	// Dynamically build parameter list with unique ids for each placeholder
	paramList := r.astFactory.CreateParamList()
	type_ := AstFactory_Type_(nodeToReplace)
	uniqueId := r.compiler.
		GetUniqueIdSupplier().
		GetUniqueId(r.compiler.GetInput(NodeUtil_GetInputId(nodeToReplace)))
	placeholderMapIds := jutil.NewLinkedHashMap[string, string]()
	for placeholderName := range placeholderMap.Keys() {
		placeholderId := placeholderName + uniqueId
		paramList.AddChildToBack(r.astFactory.CreateName(placeholderId, type_))
		placeholderMapIds.Put(placeholderName, placeholderId)
	}

	hookExpression := r.CreateHookExpressionForGenderedMsg(
		type_, message, options, nodeToReplace, placeholderMapIds)

	if placeholderMap.IsEmpty() {
		// The translated message read from the bundle is one of the following:
		// 1. has gendered variants and no placeholders
		// 2. has gendered variants and is an icu template because icu templates do not have
		// placeholders

		// Create the hook expression for the gendered message. This will be used to create the call
		// node if there are placeholders, or returned directly if there are no placeholders.
		return hookExpression
	}

	function := r.astFactory.CreateFunction("", paramList, rhino.IR_BlockOfNode(rhino.IR_ReturnNodeWithExpr(hookExpression)), type_)
	function.SetColor(nodeToReplace.GetColor())
	r.compiler.ReportChangeToChangeScope(function)

	callNode := r.astFactory.CreateCall(function, type_)

	// Add the placeholder values to the call
	for node := range placeholderMap.Values() {
		callNode.AddChildToBack(node.CloneTree())
	}
	return callNode
}

// CreateConditionForGrammaticalGender creates a node representing the
// grammatical gender condition.
//
// For example:
//
//	goog.msgKind.MASCULINE ? "Bienvenido + name" :
//	goog.msgKind.FEMININE ? "Bienvenida + name" :
//	goog.msgKind.NEUTER ? "Les damos la bienvenida + name" :
//	"Les damos la bienvenida + name";
//
// Java: ReplaceMessages.createConditionForGrammaticalGender(AstFactory.Type, JsMessage.GrammaticalGenderCase, Node)
func (r *ReplaceMessages) CreateConditionForGrammaticalGender(
	type_ AstFactory_Type, grammaticalGender JsMessage_GrammaticalGenderCase, nodeToReplace *rhino.Node) *rhino.Node {

	// NOTE: Collapse properties isn't guarantee to collapse any given property but if
	// we get a partial collapse of "goog.msgKind" then something has gone very wrong
	// so this seems reasonable rather than the alternative (traversing the AST to find the values).
	if r.collapsePropertiesHasRun {
		var referenceName string
		switch grammaticalGender {
		case JsMessage_GrammaticalGenderCase_MASCULINE:
			referenceName = "goog$msgKind$MASCULINE"
		case JsMessage_GrammaticalGenderCase_FEMININE:
			referenceName = "goog$msgKind$FEMININE"
		case JsMessage_GrammaticalGenderCase_NEUTER:
			referenceName = "goog$msgKind$NEUTER"
		default:
			referenceName = "goog$msgKind$OTHER"
		}
		result := r.astFactory.CreateName(referenceName, AstFactory_Type_(nodeToReplace))
		result.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
		return result
	} else {
		googNode := r.astFactory.CreateName("goog", type_)
		googNode.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
		return r.astFactory.CreateGetProp(
			r.astFactory.CreateGetProp(googNode, "msgKind", AstFactory_Type_(nodeToReplace)),
			grammaticalGender.String(),
			AstFactory_Type_(nodeToReplace))
	}
}

// CreateHookExpressionForGenderedMsg creates a ternary expression with the
// gendered message variants.
//
// For example:
//
//	goog.msgKind.MASCULINE ? "Bienvenido + name" :
//	goog.msgKind.FEMININE ? "Bienvenida + name" :
//	goog.msgKind.NEUTER ? "Les damos la bienvenida + name" :
//	"Les damos la bienvenida + name";
//
// Panics with a JsMessageVisitor_MalformedException (Java: throws
// MalformedException).
//
// Java: ReplaceMessages.createHookExpressionForGenderedMsg(AstFactory.Type, JsMessage, ReplaceMessages.MsgOptions, Node, Map<String, String>)
func (r *ReplaceMessages) CreateHookExpressionForGenderedMsg(
	type_ AstFactory_Type,
	msg *JsMessage,
	options *ReplaceMessages_MsgOptions,
	nodeToReplace *rhino.Node,
	placeholderMapIds *jutil.LinkedHashMap[string, string]) *rhino.Node {

	hookHead := rhino.IR_Hook(rhino.IR_Name(""), rhino.IR_String(""), rhino.IR_String(""))
	hookHead.SetColor(nodeToReplace.GetColor())
	// The previous hook condition needing to be replaced
	placeholderHook := hookHead

	placeholderNodeProvider := ReplaceMessages_PlaceholderNodeProviderFunc(func(placeholderName string) *rhino.Node {
		if placeholderId, ok := placeholderMapIds.Lookup(placeholderName); ok {
			return r.astFactory.CreateName(placeholderId, AstFactory_Type_(nodeToReplace))
		}
		return nil
	})

	for _, grammaticalGender := range msg.GetGenderedMessageVariants() {
		// Skip the OTHER case as it should be the last case in the hook expression
		if grammaticalGender == JsMessage_GrammaticalGenderCase_OTHER {
			continue
		}
		// Hook condition ex: `goog.msgKind.MASCULINE?` or `goog.msgKind.FEMININE ?` or
		// `goog.msgKind.NEUTER ?`
		condition := r.CreateConditionForGrammaticalGender(type_, grammaticalGender, nodeToReplace)

		// The last child of the hook expression will be replaced with the next currentHook
		currentHook := rhino.IR_Hook(
			condition,
			r.BuildMessageNode(
				msg.GetGenderedMessageParts(grammaticalGender),
				options,
				nodeToReplace,
				placeholderNodeProvider,
				/* isGenderedMsg= */ true),
			r.astFactory.CreateString(""))
		currentHook.SetColor(nodeToReplace.GetColor())

		if placeholderHook.GetParent() != nil {
			// Replace the previous hook condition with the current hook condition
			parent := placeholderHook.GetParent()
			parent.GetLastChild().ReplaceWith(currentHook)
			placeholderHook = currentHook.GetLastChild()
		} else {
			hookHead = currentHook
			placeholderHook = currentHook.GetLastChild()
		}
	}
	// The OTHER is always the last case in the hook expression
	placeholderHook.
		GetParent().
		GetLastChild().
		ReplaceWith(
			r.BuildMessageNode(
				msg.GetGenderedMessageParts(JsMessage_GrammaticalGenderCase_OTHER),
				options,
				nodeToReplace,
				placeholderNodeProvider,
				/* isGenderedMsg= */ true))
	return hookHead
}

// Java: interface ReplaceMessages.PlaceholderNodeProvider
type ReplaceMessages_PlaceholderNodeProvider interface {
	// Get may return nil.
	//
	// Java: ReplaceMessages.PlaceholderNodeProvider.get(String)
	Get(placeholderName string) *rhino.Node
}

// ReplaceMessages_PlaceholderNodeProviderFunc lets a Go closure stand for a
// Java lambda or method reference implementing PlaceholderNodeProvider.
//
// Go-only: language-adapter
type ReplaceMessages_PlaceholderNodeProviderFunc func(placeholderName string) *rhino.Node

// Get calls f.
func (f ReplaceMessages_PlaceholderNodeProviderFunc) Get(placeholderName string) *rhino.Node {
	return f(placeholderName)
}

// BuildMessageNode creates a parse tree corresponding to a list of message
// parts. The result consists of one or more STRING nodes, placeholder
// replacement value nodes (which can be arbitrary expressions), and ADD
// nodes.
//
// msgParts are the parts of the message to build; options are the message
// escaping options; nodeToReplace is the original node being replaced (for
// context); placeholderNodeProvider is a function that provides the value
// Node for a given placeholder name; isGenderedMsg is whether this message is
// part of a gendered message structure. It returns the root of the
// constructed parse tree.
//
// Panics with a JsMessageVisitor_MalformedException (Java: throws
// MalformedException).
//
// Java: ReplaceMessages.buildMessageNode(List<Part>, ReplaceMessages.MsgOptions, Node, ReplaceMessages.PlaceholderNodeProvider, boolean)
func (r *ReplaceMessages) BuildMessageNode(
	msgParts []JsMessage_Part,
	options *ReplaceMessages_MsgOptions,
	nodeToReplace *rhino.Node,
	placeholderNodeProvider ReplaceMessages_PlaceholderNodeProvider,
	isGenderedMsg bool) *rhino.Node {

	// A message might have normal string parts that are consecutive.
	// Join them together before processing them further. This allows split HTML entities
	// to be escaped.
	parts := ReplaceMessages_MergeStringParts(msgParts)
	if len(parts) == 0 {
		return r.astFactory.CreateString("")
	}
	var message *rhino.Node
	for _, msgPart := range parts {
		var partNode *rhino.Node
		if msgPart.IsPlaceholder() {
			jsPlaceholderName := msgPart.GetJsPlaceholderName()
			valueNode := placeholderNodeProvider.Get(jsPlaceholderName)
			if valueNode == nil {
				if isGenderedMsg {
					// Add the ICU placeholder directly to the message ex: 'Hello {NAME}'
					partNode = r.astFactory.CreateString("{" + msgPart.GetCanonicalPlaceholderName() + "}")
				} else {
					panic(NewJsMessageVisitor_MalformedException(
						"Unrecognized message placeholder referenced: "+jsPlaceholderName, nodeToReplace))
				}
			} else {
				partNode = valueNode.CloneTree()
			}
		} else {
			// The part is just a string literal.
			partNode = r.CreateNodeForMsgString(options, msgPart.GetString())
		}

		if message == nil {
			message = partNode
		} else if partNode.IsString() && message.IsString() {
			message = r.astFactory.CreateString(message.GetString() + partNode.GetString())
		} else {
			message = r.astFactory.CreateAdd(message, partNode)
		}
	}
	return message
}

// ConstructStringExprNode creates a parse tree corresponding to the remaining
// message parts in an iteration. The result consists of one or more STRING
// nodes, placeholder replacement value nodes (which can be arbitrary
// expressions), and ADD nodes.
//
// placeholderMap maps placeholder names to value Nodes. It returns the root
// of the constructed parse tree.
//
// Panics with a JsMessageVisitor_MalformedException (Java: throws
// MalformedException).
//
// Java: ReplaceMessages.constructStringExprNode(JsMessage, Map<String, Node>, ReplaceMessages.MsgOptions, Node)
func (r *ReplaceMessages) ConstructStringExprNode(
	message *JsMessage, placeholderMap *jutil.LinkedHashMap[string, *rhino.Node], options *ReplaceMessages_MsgOptions, nodeToReplace *rhino.Node) *rhino.Node {

	// If the message has gendered variants, create a hook expression containing the gendered
	// variants.
	if len(message.GetGenderedMessageVariants()) != 0 {
		return r.CreateNodeForGenderedMsgString(message, placeholderMap, options, nodeToReplace)
	} else {
		return r.ConstructNonGenderedMsgNode(message, placeholderMap, options, nodeToReplace)
	}
}

// ConstructNonGenderedMsgNode panics with a
// JsMessageVisitor_MalformedException (Java: throws MalformedException).
//
// Java: ReplaceMessages.constructNonGenderedMsgNode(JsMessage, Map<String, Node>, ReplaceMessages.MsgOptions, Node)
func (r *ReplaceMessages) ConstructNonGenderedMsgNode(
	message *JsMessage, placeholderMap *jutil.LinkedHashMap[string, *rhino.Node], options *ReplaceMessages_MsgOptions, nodeToReplace *rhino.Node) *rhino.Node {
	if placeholderMap.IsEmpty() {
		// The compiler does not expect to do any placeholder substitution, because the message
		// definition in the JS code doesn't have any placeholders.
		if options.IsIcuTemplate() {
			// This message was declared as an ICU template. Its placeholders, if any, will not be
			// replaced during compilation, but rather at runtime by a special localization method.
			// We just need to put the string back together again, if it was broken up to include
			// placeholders.
			return r.CreateNodeForMsgString(options, message.AsIcuMessageString())
		} else if message.JsPlaceholderNames().IsEmpty() {
			// The translated message read from the bundle also has no placeholders.
			// It doesn't really matter which asXMessageString() call we make, since they return
			// the same thing when there are no placeholders.
			return r.CreateNodeForMsgString(options, message.AsJsMessageString())
		} else {
			// The JS code definition has no placeholders, but the message we got from the translation
			// bundle does have placeholders.
			//
			// This can happen for ICU - formatted messages. They can contain placeholders in the
			// message string in the form "{PH_NAME}", which this compiler ignores, because they will
			// be replaced at runtime by passing the whole string to a message formatter method.
			// However, sometimes the translation bundle will represent the "{PH_NAME}" substring as
			// an actual placeholder reference, because it was generated by some tool other than
			// closure-compiler. In that case, we need to join all the parts of the message back
			// together as a single string with the "{PH_NAME}" references back in place.
			//
			// Messages with this form always start with a particular formula, which we can test for
			// here to make sure this isn't actually a case of something being broken in the translation
			// pipeline.
			icuMsgString := message.AsIcuMessageString()
			if ReplaceMessages_IsStartOfIcuMessage(icuMsgString) {
				return r.CreateNodeForMsgString(options, icuMsgString)
			} else {
				panic(NewJsMessageVisitor_MalformedException(
					"The translated message has placeholders, but the definition in the JS code does"+
						" not.",
					nodeToReplace))
			}
		}
	}

	return r.BuildMessageNode(
		message.GetParts(),
		options,
		nodeToReplace,
		ReplaceMessages_PlaceholderNodeProviderFunc(placeholderMap.Get),
		/* isGenderedMsg= */ false)
}

// Java: ReplaceMessages.createNodeForMsgString(ReplaceMessages.MsgOptions, String)
func (r *ReplaceMessages) CreateNodeForMsgString(options *ReplaceMessages_MsgOptions, s string) *rhino.Node {
	if options.EscapeLessThan() {
		// Note that "&" is not replaced because the translation can contain HTML entities.
		s = strings.ReplaceAll(s, "<", "&lt;")
	}
	if options.UnescapeHtmlEntities() {
		// Unescape entities that need to be escaped when embedding HTML or XML in data/attributes
		// of an HTML/XML document. See https://www.w3.org/TR/xml/#sec-predefined-ent.
		// Note that "&amp;" must be the last to avoid "creating" new entities.
		// To print an html entity in the resulting message, double-escape: `&amp;amp;`.
		s = strings.ReplaceAll(s, "&lt;", "<")
		s = strings.ReplaceAll(s, "&gt;", ">")
		s = strings.ReplaceAll(s, "&apos;", "'")
		s = strings.ReplaceAll(s, "&quot;", "\"")
		s = strings.ReplaceAll(s, "&amp;", "&")
	}
	return r.astFactory.CreateString(s)
}

// ReplaceMessages_MsgOptions holds the options for escaping characters in
// translated messages.
//
// Java: record ReplaceMessages.MsgOptions
type ReplaceMessages_MsgOptions struct {
	// True if the message is defined using the ICU template declaration method.
	isIcuTemplate bool
	// Replace `'<'` with `'&lt;'` in the message.
	escapeLessThan bool
	// Replace these escaped entities with their literal characters in the message (Overrides
	// escapeLessThan)
	//
	//	'&lt;' -> '<'
	//	'&gt;' -> '>'
	//	'&apos;' -> "'"
	//	'&quot;' -> '"'
	//	'&amp;' -> '&'
	unescapeHtmlEntities bool
}

// Java: ReplaceMessages.MsgOptions(boolean, boolean, boolean)
func NewReplaceMessages_MsgOptions(isIcuTemplate bool, escapeLessThan bool, unescapeHtmlEntities bool) *ReplaceMessages_MsgOptions {
	return &ReplaceMessages_MsgOptions{
		isIcuTemplate:        isIcuTemplate,
		escapeLessThan:       escapeLessThan,
		unescapeHtmlEntities: unescapeHtmlEntities,
	}
}

// Java: ReplaceMessages.MsgOptions.isIcuTemplate()
func (o *ReplaceMessages_MsgOptions) IsIcuTemplate() bool { return o.isIcuTemplate }

// Java: ReplaceMessages.MsgOptions.escapeLessThan()
func (o *ReplaceMessages_MsgOptions) EscapeLessThan() bool { return o.escapeLessThan }

// Java: ReplaceMessages.MsgOptions.unescapeHtmlEntities()
func (o *ReplaceMessages_MsgOptions) UnescapeHtmlEntities() bool { return o.unescapeHtmlEntities }

// Java: ReplaceMessages.ICU_MSG_OPTIONS
var replaceMessages_ICU_MSG_OPTIONS = NewReplaceMessages_MsgOptions(
	/* isIcuTemplate= */ true,
	/* escapeLessThan= */ false,
	/* unescapeHtmlEntities= */ false)

// ReplaceMessages_MergeStringParts merges consecutive string parts in the
// list of message parts.
//
// Java: ReplaceMessages.mergeStringParts(List<Part>)
func ReplaceMessages_MergeStringParts(parts []JsMessage_Part) []JsMessage_Part {
	var result []JsMessage_Part
	for _, part := range parts {
		if part.IsPlaceholder() {
			result = append(result, part)
		} else {
			var lastPart JsMessage_Part
			if len(result) != 0 {
				lastPart = result[len(result)-1]
			}
			if lastPart == nil || lastPart.IsPlaceholder() {
				result = append(result, part)
			} else {
				result[len(result)-1] = JsMessage_StringPart_Create(lastPart.GetString() + part.GetString())
			}
		}
	}
	return result
}

// ReplaceMessages_ProtectedJsMessage holds information about the protected
// form of a translatable message that appears in the AST.
//
// The original translatable messages are replaced with this protected form
// by the logic in `ReplaceMessages` to protect the message information
// through the optimization passes.
//
//	// original
//	var MSG_GREETING = goog.getMsg('Hello, {$name}!', {name: person.getName()}, {html: false});
//	// protected form
//	var MSG_GREETING = __jscomp_define_msg__(
//	    {
//	      'key': 'MSG_GREETING',
//	      'msg_text': 'Hello, {$name}!',
//	      'escapeLessThan': ''
//	    },
//	    {'name': person.getName()});
//
// Java: public static class ReplaceMessages.ProtectedJsMessage
type ReplaceMessages_ProtectedJsMessage struct {
	jsMessage *JsMessage
	// The expression Node that defines the message and should be replaced with the localized
	// message.
	definitionNode *rhino.Node
	// e.g. `{ name: x.getName(), age: x.getAgeString() }`
	substitutionsNode *rhino.Node // may be nil
	// The message was defined in JS code using the ICU template method.
	isIcuTemplate bool
	// Replace `'<'` with `'&lt;'` in the message.
	escapeLessThan bool
	// Replace these escaped entities with their literal characters in the message
	// (Overrides escapeLessThan)
	// '&lt;' -> '<'
	// '&gt;' -> '>'
	// '&apos;' -> "'"
	// '&quot;' -> '"'
	// '&amp;' -> '&'
	unescapeHtmlEntities bool
}

// Java: ReplaceMessages.ProtectedJsMessage(JsMessage, Node, Node, boolean, boolean, boolean)
func NewReplaceMessages_ProtectedJsMessage(
	jsMessage *JsMessage,
	definitionNode *rhino.Node,
	substitutionsNode *rhino.Node,
	isIcuTemplate bool,
	escapeLessThan bool,
	unescapeHtmlEntities bool) *ReplaceMessages_ProtectedJsMessage {
	return &ReplaceMessages_ProtectedJsMessage{
		jsMessage:            jsMessage,
		definitionNode:       definitionNode,
		substitutionsNode:    substitutionsNode,
		isIcuTemplate:        isIcuTemplate,
		escapeLessThan:       escapeLessThan,
		unescapeHtmlEntities: unescapeHtmlEntities,
	}
}

// Java: ReplaceMessages.ProtectedJsMessage.fromAstNode(Node, JsMessage.IdGenerator)
func ReplaceMessages_ProtectedJsMessage_FromAstNode(
	node *rhino.Node, idGenerator JsMessage_IdGenerator) *ReplaceMessages_ProtectedJsMessage {
	if !node.IsCall() {
		return nil
	}
	calleeNode := jutil.CheckNotNullMsg1(node.GetFirstChild(), "%s", node)
	if !calleeNode.MatchesName(ReplaceMessagesConstants_DEFINE_MSG_CALLEE) {
		return nil
	}
	propertiesNode := jutil.CheckNotNullMsg1(calleeNode.GetNext(), "%s", calleeNode)
	substitutionsNode := propertiesNode.GetNext()
	escapeLessThanOption := false
	unescapeHtmlEntitiesOption := false
	isIcuTemplate := false
	jsMessageBuilder := NewJsMessage_Builder()
	jutil.CheckStateMsg1(propertiesNode.IsObjectLit(), "%s", propertiesNode)
	var msgKey *string
	var meaning *string
	icuPlaceholderNames := jutil.NewLinkedHashSet[string]()
	var messageText *string
	var messageTextNode *rhino.Node
	for strKey := propertiesNode.GetFirstChild(); strKey != nil; strKey = strKey.GetNext() {
		jutil.CheckStateMsg1(strKey.IsStringKey(), "%s", strKey)
		key := strKey.GetString()
		valueNode := strKey.GetOnlyChild()
		if key == "icu_placeholder_names" {
			jutil.CheckStateMsg(valueNode.IsArrayLit(), "icu_placeholder_names must be an array")
			// If the message is an ICU template and `icu_placeholder_names` is present, then there
			// are placeholders in the message. These placeholders will be replaced at runtime, but
			// it is important that we keep track of these placeholders because it means that the
			// ICU template CANNOT be treated as a single string part, because having placeholders
			// means that the message has multiple parts.
			// When a message is a `declareIcuTemplate` with multiple parts, we generate a `msg id`
			// in the XMB, which is sent to the Translation Console so that the translators can
			// translate this. We generate the deterministic `msg id` using an algorithm that takes
			// into account how many parts the message has.
			// Now during JSCompiler compilation process, we protect the message by wrapping it in a
			// `__jscomp_define_msg__` (for safety because we don't want any of our optimization
			// passes to change the message).
			// Later in this method, we generate an ID (using `idGenerator.generateId()`) and use
			// this to lookup a message in the translated XTB file. As I mentioned earlier, the
			// algorithm for generating an ID needs to know the correct parts of the message, so we
			// fail to generate the same ID we did when we added the message to the XMB.
			// This `icu_placeholder_names` field is necessary to help us figure out the correct
			// parts of the message, in order to generate the correct message ID that matches the ID
			// we generated when we added the message to the XMB (which is the same ID in the XTB).
			for valueNodeChild := range valueNode.Children() {
				icuPlaceholderNames.Add(valueNodeChild.GetString())
			}
			continue
		}
		jutil.CheckStateMsg1(valueNode.IsStringLit(), "%s", valueNode)
		value := valueNode.GetString()
		switch key {
		case "key":
			jsMessageBuilder.SetKey(value)
			msgKey = &value
		case "meaning":
			jsMessageBuilder.SetMeaning(&value)
			meaning = &value
		case "alt_id":
			jsMessageBuilder.SetAlternateId(&value)
		case "msg_text":
			// This may be an ICU template that also has the `icu_placeholder_names` property, which
			// means we need to append multiple parts of the message to `jsMessageBuilder`. For now,
			// we will save the message text and current node, and we'll parse it once we know if
			// this is an ICU template with multiple parts (after this loop to run through all the
			// properties is finished).
			messageText = &value
			messageTextNode = valueNode
		case "isIcuTemplate":
			isIcuTemplate = true
		case "escapeLessThan":
			// Just being present enables this option
			escapeLessThanOption = true
		case "unescapeHtmlEntities":
			// just being present enables this option
			unescapeHtmlEntitiesOption = true
		default:
			panic(jutil.NewIllegalStateExceptionMsg("unknown protected message key: " + strKey.String()))
		}
	}

	jutil.Try(func() {
		if !icuPlaceholderNames.IsEmpty() {
			jutil.CheckStateMsg(
				isIcuTemplate,
				"Found icu_placeholder_names for a message that is not an ICU template.")
			// This is an ICU template with placeholders ("{$placeholderName}"). We cannot treat this
			// as a single string part, because it has multiple parts. Otherwise, we will generate the
			// wrong message id and we will not be able to find the correct translated message in the
			// XTB file (because when the XMB message was created during message extraction, we
			// treated this ICU template as having multiple parts).
			// new IcuMessageTemplateString(messageText): a null text throws.
			icuMessageTemplateString := NewJsMessageVisitor_IcuMessageTemplateString(*messageText)
			extractedIcuTemplateParts := icuMessageTemplateString.ExtractParts(icuPlaceholderNames)

			// Append the parts of the ICU template to the jsMessageBuilder.
			jsMessageBuilder.AppendParts(extractedIcuTemplateParts.extractedParts)
		} else {
			// This message is a single string part. It may be an ICU template without placeholders,
			// or it may be a normal `goog.getMsg()` message.
			jsMessageBuilder.AppendParts(JsMessageVisitor_ParseJsMessageTextIntoParts(*messageText))
		}
	}, func(unused IJsMessage_PlaceholderFormatException) {
		// Somehow we stored the protected message text incorrectly, which should never
		// happen 🙏
		panic(jutil.NewIllegalStateExceptionMsg(
			messageTextNode.GetLocation() +
				": Placeholder incorrectly formatted: >" +
				jutil.ToJavaString(messageText) +
				"<"))
	})

	// JsMessageVisitor.getExternalMessageId(msgKey): a null key throws.
	externalMessageId := JsMessageVisitor_GetExternalMessageId(*msgKey)
	if externalMessageId != nil {
		// MSG_EXTERNAL_12345 = ...
		jsMessageBuilder.SetIsExternalMsg(true).SetId(*externalMessageId)
	} else {
		// NOTE: If the message was anonymous (assigned to a variable or property named
		// MSG_UNNAMED_XXX), the key we have here will be the one generated from the message
		// text, and we won't end up setting the isAnonymous flag. Nothing seems to use that
		// flag...
		// TODO(bradfordcsmith): Maybe remove the isAnonymous flag for jsMessage objects?
		var meaningForIdGeneration string
		if meaning != nil {
			meaningForIdGeneration = *meaning
		} else {
			meaningForIdGeneration = JsMessageVisitor_RemoveScopedAliasesPrefix(*msgKey)
		}
		if idGenerator != nil {
			jsMessageBuilder.SetId(
				idGenerator.GenerateId(&meaningForIdGeneration, jsMessageBuilder.GetParts().Slice()))
		} else {
			jsMessageBuilder.SetId(meaningForIdGeneration)
		}
	}

	return NewReplaceMessages_ProtectedJsMessage(
		/* jsMessage= */ jsMessageBuilder.Build(),
		/* definitionNode= */ node,
		/* substitutionsNode= */ substitutionsNode,
		/* isIcuTemplate= */ isIcuTemplate,
		/* escapeLessThan= */ escapeLessThanOption,
		/* unescapeHtmlEntities= */ unescapeHtmlEntitiesOption)
}

// Java: ReplaceMessages.ProtectedJsMessage.getMsgOptions()
func (m *ReplaceMessages_ProtectedJsMessage) GetMsgOptions() *ReplaceMessages_MsgOptions {
	return NewReplaceMessages_MsgOptions(
		/* isIcuTemplate= */ m.isIcuTemplate,
		/* escapeLessThan= */ m.escapeLessThan,
		/* unescapeHtmlEntities= */ m.unescapeHtmlEntities)
}

// ReplaceMessages_IsStartOfIcuMessage detects an ICU-formatted plural or
// select message. Any placeholders occurring inside these messages must be
// rewritten in ICU format.
//
// Java: ReplaceMessages.isStartOfIcuMessage(String)
func ReplaceMessages_IsStartOfIcuMessage(part string) bool {
	// ICU messages start with a '{' followed by an identifier, followed by a ',' and then 'plural'
	// or 'select' followed by another comma.
	// the 'startsWith' check is redundant but should allow us to skip using the matcher
	if !strings.HasPrefix(part, "{") {
		return false
	}
	commaIndex := jutil.StringIndexOfCharFrom(part, ',', 1)
	// if commaIndex == 1 that means the identifier is empty, which isn't allowed.
	if commaIndex <= 1 {
		return false
	}
	nextBracketIndex := jutil.StringIndexOfCharFrom(part, '{', 1)
	// Java part.startsWith(prefix, commaIndex + 1); commaIndex < length here.
	afterComma := jutil.StringSubstringFrom(part, commaIndex+1)
	return (nextBracketIndex == -1 || nextBracketIndex > commaIndex) &&
		(strings.HasPrefix(afterComma, "plural,") ||
			strings.HasPrefix(afterComma, "select,"))
}
