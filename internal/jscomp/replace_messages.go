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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const (
	// ReplaceMessagesConstants_DEFINE_MSG_CALLEE: `goog.getMsg()` calls will
	// be converted into a call to this method which is defined in synthetic
	// externs.
	ReplaceMessagesConstants_DEFINE_MSG_CALLEE = "__jscomp_define_msg__"

	// ReplaceMessagesConstants_FALLBACK_MSG_CALLEE:
	// `goog.getMsgWithFallback(MSG_NEW, MSG_OLD)` will be converted into a
	// call to this method which is defined in synthetic externs.
	ReplaceMessagesConstants_FALLBACK_MSG_CALLEE = "__jscomp_msg_fallback__"
)

// Java: ReplaceMessagesConstants.isProtectedMessage(Node)
func ReplaceMessagesConstants_IsProtectedMessage(n rhino.Node) bool {
	call, ok := ast.AsCall(n)
	if !ok {
		return false
	}
	callee, ok := ast.AsName(call.Callee().Node())
	if !ok {
		return false
	}
	name := callee.Ident()
	return name == ReplaceMessagesConstants_DEFINE_MSG_CALLEE ||
		name == ReplaceMessagesConstants_FALLBACK_MSG_CALLEE
}

// ReplaceMessages replaces user-visible messages with alternatives.
type ReplaceMessages struct {
	compiler                 AbstractCompiler
	bundle                   MessageBundle
	astFactory               *AstFactory
	strictReplacement        bool
	collapsePropertiesHasRun bool
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
	return &ReplaceMessages{
		compiler:          compiler,
		astFactory:        compiler.CreateAstFactory(),
		bundle:            bundle,
		strictReplacement: strictReplacement,
		collapsePropertiesHasRun: compiler.GetOptions().GetPropertyCollapseLevel() == CompilerOptions_PropertyCollapseLevel_ALL &&
			compiler.GetOptions().DoLateLocalization(),
	}
}

// GetMsgProtectionPass returns a pass that replaces `goog.getMsg()` etc. calls
// with a protected form that survives optimizations.
//
// Java: ReplaceMessages.getMsgProtectionPass()
func (r *ReplaceMessages) GetMsgProtectionPass() CompilerPass {
	return NewReplaceMessages_MsgProtectionPass(r)
}

// ReplaceMessages_MsgProtectionData holds inputs for message protection.
type ReplaceMessages_MsgProtectionData struct {
	message               *JsMessage
	messageNode           rhino.Node
	templateTextNode      rhino.Node
	placeholderValuesNode rhino.Node
	messageOptions        ReplaceMessages_MsgOptions
}

// ReplaceMessages_MsgProtectionPass protects messages before optimizations.
type ReplaceMessages_MsgProtectionPass struct {
	JsMessageVisitor
	outer *ReplaceMessages
}

// Java: ReplaceMessages.MsgProtectionPass()
func NewReplaceMessages_MsgProtectionPass(outer *ReplaceMessages) *ReplaceMessages_MsgProtectionPass {
	p := &ReplaceMessages_MsgProtectionPass{outer: outer}
	p.InitJsMessageVisitor(
		outer.compiler,
		outer.bundle.IdGenerator(),
		p.ProcessJsMessageDefinition,
		p.ProcessIcuTemplateDefinition,
		p.ProcessMessageFallback,
	)
	return p
}

// Java: ReplaceMessages.MsgProtectionPass.process(Node, Node)
func (p *ReplaceMessages_MsgProtectionPass) Process(externs rhino.Node, root rhino.Node) {
	compiler := p.outer.compiler
	NodeUtil_CreateSynthesizedExternsSymbol(compiler, ReplaceMessagesConstants_DEFINE_MSG_CALLEE)
	NodeUtil_CreateSynthesizedExternsSymbol(compiler, ReplaceMessagesConstants_FALLBACK_MSG_CALLEE)
	p.JsMessageVisitor.Process(externs, root)
}

// Java: ReplaceMessages.MsgProtectionPass.processIcuTemplateDefinition(IcuTemplateDefinition)
func (p *ReplaceMessages_MsgProtectionPass) ProcessIcuTemplateDefinition(definition IcuTemplateDefinition) {
	p.PerformMessageProtection(ReplaceMessages_MsgProtectionData{
		message:               definition.GetMessage(),
		messageNode:           definition.GetMessageNode(),
		templateTextNode:      definition.GetTemplateTextNode(),
		placeholderValuesNode: rhino.NilNode,
		messageOptions:        replaceMessages_ICU_MSG_OPTIONS,
	})
}

// Java: ReplaceMessages.MsgProtectionPass.processJsMessageDefinition(JsMessageDefinition)
func (p *ReplaceMessages_MsgProtectionPass) ProcessJsMessageDefinition(definition JsMessageDefinition) {
	p.PerformMessageProtection(ReplaceMessages_MsgProtectionData{
		message:               definition.GetMessage(),
		messageNode:           definition.GetMessageNode(),
		templateTextNode:      definition.GetTemplateTextNode(),
		placeholderValuesNode: definition.GetPlaceholderValuesNode(),
		messageOptions: ReplaceMessages_MsgOptions{
			escapeLessThan:       definition.ShouldEscapeLessThan(),
			unescapeHtmlEntities: definition.ShouldUnescapeHtmlEntities(),
		},
	})
}

// Java: ReplaceMessages.MsgProtectionPass.performMessageProtection(ReplaceMessages.MsgProtectionData)
func (p *ReplaceMessages_MsgProtectionPass) PerformMessageProtection(msgProtectionData ReplaceMessages_MsgProtectionData) {
	r := p.outer
	message := msgProtectionData.message
	callNode := msgProtectionData.messageNode
	originalMessageString := msgProtectionData.templateTextNode
	placeholdersNode := msgProtectionData.placeholderValuesNode
	msgOptions := msgProtectionData.messageOptions

	call, ok := ast.AsCall(callNode)
	jutil.CheckStateMsg1(ok, "%s", callNode)
	googGetMsg := call.Callee().Node()

	newCallee := p.CreateProtectionFunctionCallee(ReplaceMessagesConstants_DEFINE_MSG_CALLEE).Srcref(googGetMsg)
	msgPropertiesNode := r.CreateMsgPropertiesNode(message, msgOptions).SrcrefTree(originalMessageString)
	newCallNode := r.astFactory.CreateCall(newCallee, AstFactory_Type_(callNode), msgPropertiesNode).Srcref(callNode)
	newCallNode.SetSideEffectFlags(rhino.Node_SideEffectFlags_NO_SIDE_EFFECTS)
	if !placeholdersNode.IsNil() {
		placeholdersLit, ok := ast.AsObjectLit(placeholdersNode)
		jutil.CheckStateMsg1(ok, "%s", placeholdersNode)
		for member := range placeholdersLit.Members() {
			strKey := member.Node()
			jutil.CheckStateMsg1(strKey.IsStringKey(), "%s", strKey)
			strKey.SetQuotedStringKey()
		}
		newCallNode.AddChildToBack(placeholdersNode.Detach())
	}
	callNode.ReplaceWith(newCallNode)
	r.compiler.ReportChangeToEnclosingScope(newCallNode)
}

// Java: ReplaceMessages.MsgProtectionPass.createProtectionFunctionCallee(String)
func (p *ReplaceMessages_MsgProtectionPass) CreateProtectionFunctionCallee(protectionFunctionName string) rhino.Node {
	callee := p.outer.astFactory.CreateNameWithUnknownType(protectionFunctionName)
	callee.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	return callee
}

// Java: ReplaceMessages.MsgProtectionPass.processMessageFallback(Node, JsMessage, JsMessage)
func (p *ReplaceMessages_MsgProtectionPass) ProcessMessageFallback(callNode rhino.Node, message1 *JsMessage, message2 *JsMessage) {
	r := p.outer
	originalCallee := rhino.CheckNotNullMsg1(callNode.GetFirstChild(), "%s", callNode)
	fallbackCallee := p.CreateProtectionFunctionCallee(ReplaceMessagesConstants_FALLBACK_MSG_CALLEE).
		Srcref(originalCallee)

	originalFirstArg := rhino.CheckNotNullMsg1(originalCallee.GetNext(), "%s", callNode)
	firstMsgKey := r.astFactory.CreateString(message1.GetKey()).Srcref(originalFirstArg)

	originalSecondArg := rhino.CheckNotNullMsg1(originalFirstArg.GetNext(), "%s", callNode)
	secondMsgKey := r.astFactory.CreateString(message2.GetKey()).Srcref(originalSecondArg)

	newCallNode := r.astFactory.
		CreateCall(
			fallbackCallee,
			AstFactory_Type_(callNode),
			firstMsgKey,
			originalFirstArg.Detach(),
			secondMsgKey,
			originalSecondArg.Detach()).
		Srcref(callNode)
	newCallNode.SetSideEffectFlags(rhino.Node_SideEffectFlags_NO_SIDE_EFFECTS)
	callNode.ReplaceWith(newCallNode)
	r.compiler.ReportChangeToEnclosingScope(newCallNode)
}

// Java: ReplaceMessages.createMsgPropertiesNode(JsMessage, ReplaceMessages.MsgOptions)
func (r *ReplaceMessages) CreateMsgPropertiesNode(message *JsMessage, msgOptions ReplaceMessages_MsgOptions) rhino.Node {
	result := r.astFactory.CreateObjectLit()
	addString := func(key, value string) {
		result.AddChildToBack(r.astFactory.CreateQuotedStringKey(key, r.astFactory.CreateString(value)))
	}
	addString("key", message.GetKey())
	if msgOptions.isIcuTemplate && !message.CanonicalPlaceholderNames().IsEmpty() {
		namesArrayLit := r.astFactory.CreateArraylit()
		for name := range message.CanonicalPlaceholderNames().All() {
			namesArrayLit.AddChildToBack(r.astFactory.CreateString(name))
		}
		result.AddChildToBack(r.astFactory.CreateQuotedStringKey("icu_placeholder_names", namesArrayLit))
	}
	if altId, ok := message.AlternateIdOpt(); ok {
		addString("alt_id", altId)
	}
	if meaning, ok := message.MeaningOpt(); ok {
		addString("meaning", meaning)
	}
	if msgOptions.isIcuTemplate {
		addString("msg_text", message.AsIcuMessageString())
		addString("isIcuTemplate", "")
	} else {
		addString("msg_text", message.AsJsMessageString())
	}
	if msgOptions.escapeLessThan {
		addString("escapeLessThan", "")
	}
	if msgOptions.unescapeHtmlEntities {
		addString("unescapeHtmlEntities", "")
	}
	return result
}

// GetReplacementCompletionPass returns the post-optimization pass that replaces
// protected messages with their final translated expressions.
//
// Java: ReplaceMessages.getReplacementCompletionPass()
func (r *ReplaceMessages) GetReplacementCompletionPass() CompilerPass {
	return NewReplaceMessages_ReplacementCompletionPass(r)
}

// ReplaceMessages_ReplacementCompletionPass replaces protected messages with
// localized expressions.
type ReplaceMessages_ReplacementCompletionPass struct {
	outer             *ReplaceMessages
	translatedMsgKeys map[string]struct{}
}

// Java: ReplaceMessages.ReplacementCompletionPass()
func NewReplaceMessages_ReplacementCompletionPass(outer *ReplaceMessages) *ReplaceMessages_ReplacementCompletionPass {
	return &ReplaceMessages_ReplacementCompletionPass{outer: outer}
}

// ShouldTraverse implements NodeTraversal_Callback (post-order traversal).
func (p *ReplaceMessages_ReplacementCompletionPass) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return true
}

// Visit implements NodeTraversal_Callback.
func (p *ReplaceMessages_ReplacementCompletionPass) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if protectedJsMessage, ok := ReplaceMessages_ProtectedJsMessage_FromAstNode(n, p.outer.bundle.IdGenerator()); ok {
		p.VisitMsgDefinition(protectedJsMessage)
	} else if protectedMsgFallback, ok := ReplaceMessages_ProtectedMsgFallback_FromAstNode(n); ok {
		p.VisitMsgFallback(protectedMsgFallback)
	}
}

// Java: ReplaceMessages.ReplacementCompletionPass.process(Node, Node)
func (p *ReplaceMessages_ReplacementCompletionPass) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(p.outer.compiler, root, p)
}

// Java: ReplaceMessages.ReplacementCompletionPass.visitMsgDefinition(ReplaceMessages.ProtectedJsMessage)
func (p *ReplaceMessages_ReplacementCompletionPass) VisitMsgDefinition(protectedJsMessage ReplaceMessages_ProtectedJsMessage) {
	r := p.outer
	originalMsg := protectedJsMessage.jsMessage
	nodeToReplace := protectedJsMessage.definitionNode
	translatedMsg := r.LookupMessage(nodeToReplace, r.bundle, originalMsg)
	var msgToUse *JsMessage
	if translatedMsg != nil {
		msgToUse = translatedMsg
		if p.translatedMsgKeys == nil {
			p.translatedMsgKeys = make(map[string]struct{})
		}
		p.translatedMsgKeys[originalMsg.GetKey()] = struct{}{}
	} else {
		if r.strictReplacement {
			r.compiler.Report(
				JSError_MakeWithNode(nodeToReplace, ReplaceMessages_BUNDLE_DOES_NOT_HAVE_THE_MESSAGE, originalMsg.GetId()))
		}
		msgToUse = originalMsg
	}
	placeholderMap := ReplaceMessages_ExtractPlaceholderValuesMapOrThrow(protectedJsMessage.substitutionsNode)
	finalMsgConstructionExpression, err := r.ConstructStringExprNode(msgToUse, placeholderMap, protectedJsMessage.options, nodeToReplace)
	if err != nil {
		r.compiler.Report(JSError_MakeWithNode(err.node, JsMessageVisitor_MESSAGE_TREE_MALFORMED, err.msg))
		return
	}
	finalMsgConstructionExpression.SrcrefTreeIfMissing(nodeToReplace)
	nodeToReplace.ReplaceWith(finalMsgConstructionExpression)
	r.compiler.ReportChangeToEnclosingScope(finalMsgConstructionExpression)
}

// Java: ReplaceMessages.ReplacementCompletionPass.visitMsgFallback(ReplaceMessages.ProtectedMsgFallback)
func (p *ReplaceMessages_ReplacementCompletionPass) VisitMsgFallback(protectedMsgFallback ReplaceMessages_ProtectedMsgFallback) {
	var valueNodeToUse rhino.Node
	if _, ok := p.translatedMsgKeys[protectedMsgFallback.firstMsgKey]; ok {
		valueNodeToUse = protectedMsgFallback.firstMsgValue
	} else if _, ok := p.translatedMsgKeys[protectedMsgFallback.secondMsgKey]; ok {
		valueNodeToUse = protectedMsgFallback.secondMsgValue
	} else {
		valueNodeToUse = protectedMsgFallback.firstMsgValue
	}
	valueNodeToUse.Detach()
	protectedMsgFallback.callNode.ReplaceWith(valueNodeToUse)
	p.outer.compiler.ReportChangeToEnclosingScope(valueNodeToUse)
}

// Java: ReplaceMessages.extractPlaceholderValuesMapOrThrow(Node)
func ReplaceMessages_ExtractPlaceholderValuesMapOrThrow(valuesObjLit rhino.Node) *jutil.LinkedHashMap[string, rhino.Node] {
	if valuesObjLit.IsNil() {
		return &jsMessageEmptyNodeMap
	}
	m, err := extractObjectLiteralMap(valuesObjLit)
	if err != nil {
		panic(jutil.NewIllegalStateExceptionMsg(err.Error()))
	}
	if res := m.extractAsValueMap(); res != nil {
		return res
	}
	return &jsMessageEmptyNodeMap
}

// ReplaceMessages_ProtectedMsgFallback represents a `__jscomp_msg_fallback__` call.
type ReplaceMessages_ProtectedMsgFallback struct {
	firstMsgKey    string
	secondMsgKey   string
	callNode       rhino.Node
	firstMsgValue  rhino.Node
	secondMsgValue rhino.Node
}

// Java: ReplaceMessages.ProtectedMsgFallback.fromAstNode(Node)
func ReplaceMessages_ProtectedMsgFallback_FromAstNode(n rhino.Node) (ReplaceMessages_ProtectedMsgFallback, bool) {
	call, ok := ast.AsCall(n)
	if !ok {
		return ReplaceMessages_ProtectedMsgFallback{}, false
	}
	callee := call.Callee().Node()
	if !callee.MatchesName(ReplaceMessagesConstants_FALLBACK_MSG_CALLEE) {
		return ReplaceMessages_ProtectedMsgFallback{}, false
	}
	jutil.CheckStateMsg1(n.HasXChildren(5), "bad message fallback call: %s", n)
	firstMsgKeyNode := callee.GetNext()
	firstMsgValue := firstMsgKeyNode.GetNext()
	secondMsgKeyNode := firstMsgValue.GetNext()
	secondMsgValue := secondMsgKeyNode.GetNext()
	return ReplaceMessages_ProtectedMsgFallback{
		callNode:       n,
		firstMsgKey:    firstMsgKeyNode.GetString(),
		firstMsgValue:  firstMsgValue,
		secondMsgKey:   secondMsgKeyNode.GetString(),
		secondMsgValue: secondMsgValue,
	}, true
}

// Java: ReplaceMessages.lookupMessage(Node, MessageBundle, JsMessage)
func (r *ReplaceMessages) LookupMessage(
	callNode rhino.Node, bundle MessageBundle, message *JsMessage,
) *JsMessage {
	translatedMessage := bundle.GetMessage(message.GetId())
	if translatedMessage != nil {
		return translatedMessage
	}

	alternateId, hasAlt := message.AlternateIdOpt()
	if !hasAlt {
		return nil
	}

	alternateMessage := bundle.GetMessage(alternateId)
	if alternateMessage != nil {
		jsCodePlaceholderNames := message.JsPlaceholderNames()
		alternateMsgPlaceholderNames := alternateMessage.JsPlaceholderNames()
		if !jsMessage_stringSetEqual(jsCodePlaceholderNames, alternateMsgPlaceholderNames) {
			if jsCodePlaceholderNames.IsEmpty() && ReplaceMessages_IsStartOfIcuMessage(message.AsIcuMessageString()) {
				return alternateMessage
			}
			r.compiler.Report(
				JSError_MakeWithNode(
					callNode,
					ReplaceMessages_INVALID_ALTERNATE_MESSAGE_PLACEHOLDERS,
					alternateId,
					alternateMsgPlaceholderNames.String(),
					message.GetKey(),
					jsCodePlaceholderNames.String()))
			return nil
		}
	}
	return alternateMessage
}

// GetFullReplacementPass returns the single-phase message replacement pass.
//
// Java: ReplaceMessages.getFullReplacementPass()
func (r *ReplaceMessages) GetFullReplacementPass() CompilerPass {
	return NewReplaceMessages_FullReplacementPass(r)
}

// ReplaceMessages_FullReplacementMsgData holds inputs for single-phase replacement.
type ReplaceMessages_FullReplacementMsgData struct {
	message             *JsMessage
	placeholderValueMap *jutil.LinkedHashMap[string, rhino.Node]
	messageNode         rhino.Node
	messageOptions      ReplaceMessages_MsgOptions
}

// ReplaceMessages_FullReplacementPass performs single-phase message replacement.
type ReplaceMessages_FullReplacementPass struct {
	JsMessageVisitor
	outer *ReplaceMessages
}

// Java: ReplaceMessages.FullReplacementPass()
func NewReplaceMessages_FullReplacementPass(outer *ReplaceMessages) *ReplaceMessages_FullReplacementPass {
	p := &ReplaceMessages_FullReplacementPass{outer: outer}
	p.InitJsMessageVisitor(
		outer.compiler,
		outer.bundle.IdGenerator(),
		p.ProcessJsMessageDefinition,
		p.ProcessIcuTemplateDefinition,
		p.ProcessMessageFallback,
	)
	return p
}

// Java: ReplaceMessages.FullReplacementPass.processMessageFallback(Node, JsMessage, JsMessage)
func (p *ReplaceMessages_FullReplacementPass) ProcessMessageFallback(callNode rhino.Node, message1 *JsMessage, message2 *JsMessage) {
	r := p.outer
	isFirstMessageTranslated := r.LookupMessage(callNode, r.bundle, message1) != nil
	isSecondMessageTranslated := r.LookupMessage(callNode, r.bundle, message2) != nil
	var replacementNode rhino.Node
	if isSecondMessageTranslated && !isFirstMessageTranslated {
		replacementNode = callNode.GetChildAtIndex(2)
	} else {
		replacementNode = callNode.GetSecondChild()
	}
	callNode.ReplaceWith(replacementNode.Detach())
	changeScope := ChangeTracker_GetEnclosingChangeScopeRoot(replacementNode)
	if !changeScope.IsNil() {
		r.compiler.ReportChangeToChangeScope(changeScope)
	}
}

// Java: ReplaceMessages.FullReplacementPass.processIcuTemplateDefinition(IcuTemplateDefinition)
func (p *ReplaceMessages_FullReplacementPass) ProcessIcuTemplateDefinition(definition IcuTemplateDefinition) {
	p.ProcessFullReplacement(ReplaceMessages_FullReplacementMsgData{
		message:             definition.GetMessage(),
		messageNode:         definition.GetMessageNode(),
		messageOptions:      replaceMessages_ICU_MSG_OPTIONS,
		placeholderValueMap: &jsMessageEmptyNodeMap,
	})
}

// Java: ReplaceMessages.FullReplacementPass.processJsMessageDefinition(JsMessageDefinition)
func (p *ReplaceMessages_FullReplacementPass) ProcessJsMessageDefinition(definition JsMessageDefinition) {
	p.ProcessFullReplacement(ReplaceMessages_FullReplacementMsgData{
		message:     definition.GetMessage(),
		messageNode: definition.GetMessageNode(),
		messageOptions: ReplaceMessages_MsgOptions{
			escapeLessThan:       definition.ShouldEscapeLessThan(),
			unescapeHtmlEntities: definition.ShouldUnescapeHtmlEntities(),
		},
		placeholderValueMap: definition.GetPlaceholderValueMap(),
	})
}

// Java: ReplaceMessages.FullReplacementPass.processFullReplacement(ReplaceMessages.FullReplacementMsgData)
func (p *ReplaceMessages_FullReplacementPass) ProcessFullReplacement(fullReplacementMsgData ReplaceMessages_FullReplacementMsgData) {
	r := p.outer
	message := fullReplacementMsgData.message
	msgNode := fullReplacementMsgData.messageNode
	options := fullReplacementMsgData.messageOptions
	placeholderValueMap := fullReplacementMsgData.placeholderValueMap

	replacement := r.LookupMessage(msgNode, r.bundle, message)
	if replacement == nil {
		if r.strictReplacement {
			r.compiler.Report(JSError_MakeWithNode(msgNode, ReplaceMessages_BUNDLE_DOES_NOT_HAVE_THE_MESSAGE, message.GetId()))
			return
		}
		replacement = message
	}

	newValue, err := r.ConstructStringExprNode(replacement, placeholderValueMap, options, msgNode)
	if err != nil {
		r.compiler.Report(JSError_MakeWithNode(err.node, JsMessageVisitor_MESSAGE_TREE_MALFORMED, err.msg))
		return
	}

	if newValue != msgNode {
		newValue.SrcrefTreeIfMissing(msgNode)
		msgNode.ReplaceWith(newValue)
		r.compiler.ReportChangeToEnclosingScope(newValue)
	}
}

// Java: ReplaceMessages.createNodeForGenderedMsgString(JsMessage, Map<String, Node>, ReplaceMessages.MsgOptions, Node)
func (r *ReplaceMessages) CreateNodeForGenderedMsgString(
	message *JsMessage,
	placeholderMap *jutil.LinkedHashMap[string, rhino.Node],
	options ReplaceMessages_MsgOptions,
	nodeToReplace rhino.Node,
) (rhino.Node, *malformedMessageError) {
	paramList := r.astFactory.CreateParamList()
	type_ := AstFactory_Type_(nodeToReplace)
	uniqueId := r.compiler.
		GetUniqueIdSupplier().
		GetUniqueId(r.compiler.GetInput(NodeUtil_GetInputId(nodeToReplace)))
	var placeholderMapIds map[string]string
	if !placeholderMap.IsEmpty() {
		placeholderMapIds = make(map[string]string, placeholderMap.Size())
		for placeholderName := range placeholderMap.Keys() {
			placeholderId := placeholderName + uniqueId
			paramList.AddChildToBack(r.astFactory.CreateName(placeholderId, type_))
			placeholderMapIds[placeholderName] = placeholderId
		}
	}

	hookExpression, err := r.CreateHookExpressionForGenderedMsg(
		type_, message, options, nodeToReplace, placeholderMapIds)
	if err != nil {
		return rhino.NilNode, err
	}

	if placeholderMap.IsEmpty() {
		return hookExpression, nil
	}

	function := r.astFactory.CreateFunction("", paramList, rhino.IR_BlockOfNode(rhino.IR_ReturnNodeWithExpr(hookExpression)), type_)
	function.SetColor(nodeToReplace.GetColor())
	r.compiler.ReportChangeToChangeScope(function)

	callNode := r.astFactory.CreateCall(function, type_)
	for node := range placeholderMap.Values() {
		callNode.AddChildToBack(node.CloneTree())
	}
	return callNode, nil
}

// Java: ReplaceMessages.createConditionForGrammaticalGender(AstFactory.Type, JsMessage.GrammaticalGenderCase, Node)
func (r *ReplaceMessages) CreateConditionForGrammaticalGender(
	type_ AstFactory_Type, grammaticalGender JsMessage_GrammaticalGenderCase, nodeToReplace rhino.Node,
) rhino.Node {
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
	}
	googNode := r.astFactory.CreateName("goog", type_)
	googNode.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
	return r.astFactory.CreateGetProp(
		r.astFactory.CreateGetProp(googNode, "msgKind", AstFactory_Type_(nodeToReplace)),
		grammaticalGender.String(),
		AstFactory_Type_(nodeToReplace))
}

// Java: ReplaceMessages.createHookExpressionForGenderedMsg(AstFactory.Type, JsMessage, ReplaceMessages.MsgOptions, Node, Map<String, String>)
func (r *ReplaceMessages) CreateHookExpressionForGenderedMsg(
	type_ AstFactory_Type,
	msg *JsMessage,
	options ReplaceMessages_MsgOptions,
	nodeToReplace rhino.Node,
	placeholderMapIds map[string]string,
) (rhino.Node, *malformedMessageError) {
	hookHead := rhino.IR_Hook(rhino.IR_Name(""), rhino.IR_String(""), rhino.IR_String(""))
	hookHead.SetColor(nodeToReplace.GetColor())
	placeholderHook := hookHead

	placeholderNodeProvider := func(placeholderName string) rhino.Node {
		if placeholderId, ok := placeholderMapIds[placeholderName]; ok {
			return r.astFactory.CreateName(placeholderId, AstFactory_Type_(nodeToReplace))
		}
		return rhino.NilNode
	}

	for grammaticalGender, genderParts := range msg.GetGenderedMessagesMap().All() {
		if grammaticalGender == JsMessage_GrammaticalGenderCase_OTHER {
			continue
		}
		condition := r.CreateConditionForGrammaticalGender(type_, grammaticalGender, nodeToReplace)
		branchNode, err := r.BuildMessageNode(
			genderParts, options, nodeToReplace, placeholderNodeProvider, true)
		if err != nil {
			return rhino.NilNode, err
		}
		currentHook := rhino.IR_Hook(condition, branchNode, r.astFactory.CreateString(""))
		currentHook.SetColor(nodeToReplace.GetColor())

		if !placeholderHook.GetParent().IsNil() {
			parent := placeholderHook.GetParent()
			parent.GetLastChild().ReplaceWith(currentHook)
			placeholderHook = currentHook.GetLastChild()
		} else {
			hookHead = currentHook
			placeholderHook = currentHook.GetLastChild()
		}
	}

	otherNode, err := r.BuildMessageNode(
		msg.GetGenderedMessageParts(JsMessage_GrammaticalGenderCase_OTHER),
		options,
		nodeToReplace,
		placeholderNodeProvider,
		true,
	)
	if err != nil {
		return rhino.NilNode, err
	}
	placeholderHook.GetParent().GetLastChild().ReplaceWith(otherNode)
	return hookHead, nil
}

// Java: ReplaceMessages.buildMessageNode(List<Part>, ReplaceMessages.MsgOptions, Node, ReplaceMessages.PlaceholderNodeProvider, boolean)
func (r *ReplaceMessages) BuildMessageNode(
	msgParts []JsMessage_Part,
	options ReplaceMessages_MsgOptions,
	nodeToReplace rhino.Node,
	placeholderNodeProvider func(string) rhino.Node,
	isGenderedMsg bool,
) (rhino.Node, *malformedMessageError) {
	parts := ReplaceMessages_MergeStringParts(msgParts)
	if len(parts) == 0 {
		return r.astFactory.CreateString(""), nil
	}
	var message rhino.Node
	for _, msgPart := range parts {
		var partNode rhino.Node
		if msgPart.IsPlaceholder() {
			jsPlaceholderName := msgPart.GetJsPlaceholderName()
			valueNode := placeholderNodeProvider(jsPlaceholderName)
			if valueNode.IsNil() {
				if isGenderedMsg {
					partNode = r.astFactory.CreateString("{" + msgPart.GetCanonicalPlaceholderName() + "}")
				} else {
					return rhino.NilNode, &malformedMessageError{
						msg:  "Unrecognized message placeholder referenced: " + jsPlaceholderName,
						node: nodeToReplace,
					}
				}
			} else {
				partNode = valueNode.CloneTree()
			}
		} else {
			partNode = r.CreateNodeForMsgString(options, msgPart.GetString())
		}

		if message.IsNil() {
			message = partNode
		} else if partNode.IsString() && message.IsString() {
			message = r.astFactory.CreateString(message.GetString() + partNode.GetString())
		} else {
			message = r.astFactory.CreateAdd(message, partNode)
		}
	}
	return message, nil
}

// Java: ReplaceMessages.constructStringExprNode(JsMessage, Map<String, Node>, ReplaceMessages.MsgOptions, Node)
func (r *ReplaceMessages) ConstructStringExprNode(
	message *JsMessage,
	placeholderMap *jutil.LinkedHashMap[string, rhino.Node],
	options ReplaceMessages_MsgOptions,
	nodeToReplace rhino.Node,
) (rhino.Node, *malformedMessageError) {
	if message.HasGenderedVariants() {
		return r.CreateNodeForGenderedMsgString(message, placeholderMap, options, nodeToReplace)
	}
	return r.ConstructNonGenderedMsgNode(message, placeholderMap, options, nodeToReplace)
}

// Java: ReplaceMessages.constructNonGenderedMsgNode(JsMessage, Map<String, Node>, ReplaceMessages.MsgOptions, Node)
func (r *ReplaceMessages) ConstructNonGenderedMsgNode(
	message *JsMessage,
	placeholderMap *jutil.LinkedHashMap[string, rhino.Node],
	options ReplaceMessages_MsgOptions,
	nodeToReplace rhino.Node,
) (rhino.Node, *malformedMessageError) {
	if placeholderMap.IsEmpty() {
		if options.isIcuTemplate {
			return r.CreateNodeForMsgString(options, message.AsIcuMessageString()), nil
		}
		if message.JsPlaceholderNames().IsEmpty() {
			return r.CreateNodeForMsgString(options, message.AsJsMessageString()), nil
		}
		icuMsgString := message.AsIcuMessageString()
		if ReplaceMessages_IsStartOfIcuMessage(icuMsgString) {
			return r.CreateNodeForMsgString(options, icuMsgString), nil
		}
		return rhino.NilNode, &malformedMessageError{
			msg:  "The translated message has placeholders, but the definition in the JS code does not.",
			node: nodeToReplace,
		}
	}

	return r.BuildMessageNode(
		message.GetParts(),
		options,
		nodeToReplace,
		placeholderMap.Get,
		false,
	)
}

// Java: ReplaceMessages.createNodeForMsgString(ReplaceMessages.MsgOptions, String)
func (r *ReplaceMessages) CreateNodeForMsgString(options ReplaceMessages_MsgOptions, s string) rhino.Node {
	if options.escapeLessThan {
		s = strings.ReplaceAll(s, "<", "&lt;")
	}
	if options.unescapeHtmlEntities {
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
type ReplaceMessages_MsgOptions struct {
	isIcuTemplate        bool
	escapeLessThan       bool
	unescapeHtmlEntities bool
}

var replaceMessages_ICU_MSG_OPTIONS = ReplaceMessages_MsgOptions{isIcuTemplate: true}

// ReplaceMessages_MergeStringParts merges consecutive string parts in the
// list of message parts, returning parts unchanged when no merge is needed.
//
// Java: ReplaceMessages.mergeStringParts(List<Part>)
func ReplaceMessages_MergeStringParts(parts []JsMessage_Part) []JsMessage_Part {
	needsMerge := false
	for i := 1; i < len(parts); i++ {
		if !parts[i-1].IsPlaceholder() && !parts[i].IsPlaceholder() {
			needsMerge = true
			break
		}
	}
	if !needsMerge {
		return parts
	}
	result := make([]JsMessage_Part, 0, len(parts)-1)
	for _, part := range parts {
		if !part.IsPlaceholder() && len(result) > 0 && !result[len(result)-1].IsPlaceholder() {
			result[len(result)-1] = JsMessage_StringPart_Create(result[len(result)-1].GetString() + part.GetString())
		} else {
			result = append(result, part)
		}
	}
	return result
}

// ReplaceMessages_ProtectedJsMessage holds information about the protected
// form (`__jscomp_define_msg__`) of a translatable message in the AST.
type ReplaceMessages_ProtectedJsMessage struct {
	jsMessage         *JsMessage
	definitionNode    rhino.Node
	substitutionsNode rhino.Node
	options           ReplaceMessages_MsgOptions
}

// Java: ReplaceMessages.ProtectedJsMessage.fromAstNode(Node, JsMessage.IdGenerator)
func ReplaceMessages_ProtectedJsMessage_FromAstNode(
	node rhino.Node, idGenerator JsMessage_IdGenerator,
) (ReplaceMessages_ProtectedJsMessage, bool) {
	call, ok := ast.AsCall(node)
	if !ok {
		return ReplaceMessages_ProtectedJsMessage{}, false
	}
	calleeNode := call.Callee().Node()
	if !calleeNode.MatchesName(ReplaceMessagesConstants_DEFINE_MSG_CALLEE) {
		return ReplaceMessages_ProtectedJsMessage{}, false
	}
	propertiesNode := rhino.CheckNotNullMsg1(calleeNode.GetNext(), "%s", calleeNode)
	substitutionsNode := propertiesNode.GetNext()
	var options ReplaceMessages_MsgOptions
	jsMessageBuilder := NewJsMessage_Builder()
	propertiesLit, ok := ast.AsObjectLit(propertiesNode)
	jutil.CheckStateMsg1(ok, "%s", propertiesNode)

	var msgKey, meaning, messageText string
	var hasMsgKey, hasMeaning, hasMessageText bool
	var icuPlaceholderNames map[string]struct{}
	var messageTextNode rhino.Node

	for member := range propertiesLit.Members() {
		strKey, ok := ast.AsStringKey(member.Node())
		jutil.CheckStateMsg1(ok, "%s", member.Node())
		key := strKey.Key()
		valueNode := member.Node().GetOnlyChild()
		if key == "icu_placeholder_names" {
			arr, ok := ast.AsArrayLit(valueNode)
			jutil.CheckStateMsg(ok, "icu_placeholder_names must be an array")
			for elem := range arr.Elements() {
				if icuPlaceholderNames == nil {
					icuPlaceholderNames = make(map[string]struct{})
				}
				icuPlaceholderNames[elem.Node().GetString()] = struct{}{}
			}
			continue
		}
		strLit, ok := ast.AsStringLit(valueNode)
		jutil.CheckStateMsg1(ok, "%s", valueNode)
		value := strLit.Value()
		switch key {
		case "key":
			jsMessageBuilder.SetKey(value)
			msgKey = value
			hasMsgKey = true
		case "meaning":
			jsMessageBuilder.SetMeaningString(value)
			meaning = value
			hasMeaning = true
		case "alt_id":
			jsMessageBuilder.SetAlternateIdString(value)
		case "msg_text":
			messageText = value
			hasMessageText = true
			messageTextNode = valueNode
		case "isIcuTemplate":
			options.isIcuTemplate = true
		case "escapeLessThan":
			options.escapeLessThan = true
		case "unescapeHtmlEntities":
			options.unescapeHtmlEntities = true
		default:
			panic(jutil.NewIllegalStateExceptionMsg("unknown protected message key: " + member.Node().String()))
		}
	}

	if !hasMessageText {
		panic(jutil.NewNullPointerException())
	}
	if len(icuPlaceholderNames) > 0 {
		jutil.CheckStateMsg(
			options.isIcuTemplate,
			"Found icu_placeholder_names for a message that is not an ICU template.")
		icuMessageTemplateString := NewJsMessageVisitor_IcuMessageTemplateString(messageText)
		extractedIcuTemplateParts := icuMessageTemplateString.ExtractParts(icuPlaceholderNames)
		jsMessageBuilder.AppendParts(extractedIcuTemplateParts.extractedParts)
	} else {
		parts, err := jsMessageVisitor_parseJsMessageTextIntoParts(messageText)
		if err != nil {
			panic(jutil.NewIllegalStateExceptionMsg(
				messageTextNode.GetLocation() +
					": Placeholder incorrectly formatted: >" +
					messageText +
					"<"))
		}
		jsMessageBuilder.AppendParts(parts)
	}

	if !hasMsgKey {
		panic(jutil.NewNullPointerException())
	}
	if externalMessageId, isExternal := jsMessageVisitor_externalMessageId(msgKey); isExternal {
		jsMessageBuilder.SetIsExternalMsg(true).SetId(externalMessageId)
	} else {
		meaningForIdGeneration := meaning
		if !hasMeaning {
			meaningForIdGeneration = JsMessageVisitor_RemoveScopedAliasesPrefix(msgKey)
		}
		if idGenerator != nil {
			jsMessageBuilder.SetId(
				generateMessageId(idGenerator, meaningForIdGeneration, jsMessageBuilder.GetParts()))
		} else {
			jsMessageBuilder.SetId(meaningForIdGeneration)
		}
	}

	return ReplaceMessages_ProtectedJsMessage{
		jsMessage:         jsMessageBuilder.Build(),
		definitionNode:    node,
		substitutionsNode: substitutionsNode,
		options:           options,
	}, true
}

// ReplaceMessages_IsStartOfIcuMessage detects an ICU-formatted plural or
// select message.
//
// Java: ReplaceMessages.isStartOfIcuMessage(String)
func ReplaceMessages_IsStartOfIcuMessage(part string) bool {
	if !strings.HasPrefix(part, "{") {
		return false
	}
	commaIndex := strings.IndexByte(part[1:], ',')
	if commaIndex <= 0 {
		return false
	}
	commaIndex++ // offset in part
	nextBracketRel := strings.IndexByte(part[1:], '{')
	nextBracketIndex := -1
	if nextBracketRel >= 0 {
		nextBracketIndex = nextBracketRel + 1
	}
	afterComma := part[commaIndex+1:]
	return (nextBracketIndex == -1 || nextBracketIndex > commaIndex) &&
		(strings.HasPrefix(afterComma, "plural,") ||
			strings.HasPrefix(afterComma, "select,"))
}
