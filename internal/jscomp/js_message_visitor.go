// Port of com/google/javascript/jscomp/JsMessageVisitor.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/JsMessageExtractor.java @ closure-compiler v20261021

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"errors"
	"slices"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const (
	jsMessageVisitor_MSG_FUNCTION_NAME          = "getMsg"
	jsMessageVisitor_MSG_FUNCTION_QNAME         = "goog." + jsMessageVisitor_MSG_FUNCTION_NAME
	jsMessageVisitor_ICU_MSG_FUNCTION_NAME      = "declareIcuTemplate"
	jsMessageVisitor_ICU_MSG_FUNCTION_QNAME     = "goog.i18n.messages." + jsMessageVisitor_ICU_MSG_FUNCTION_NAME
	jsMessageVisitor_MSG_FALLBACK_FUNCTION_NAME = "goog.getMsgWithFallback"

	// jsMessageVisitor_MSG_EXTERNAL_PREFIX identifies a message with a specific
	// ID which doesn't get extracted from the JS code containing its declaration.
	jsMessageVisitor_MSG_EXTERNAL_PREFIX = "MSG_EXTERNAL_"

	JsMessageVisitor_MSG_PREFIX = "MSG_"

	jsMessageVisitor_SCOPED_ALIASES_PREFIX = "$jscomp$scope$"
)

// JsMessageVisitor_FUNCTIONS_USED_IN_MESSAGE_EXTRACTION exposes a list of
// functions that are used in message extraction.
var JsMessageVisitor_FUNCTIONS_USED_IN_MESSAGE_EXTRACTION = []string{
	jsMessageVisitor_MSG_FUNCTION_NAME, jsMessageVisitor_ICU_MSG_FUNCTION_NAME,
}

var JsMessageVisitor_MESSAGE_HAS_NO_DESCRIPTION = DiagnosticType_Warning(
	"JSC_MSG_HAS_NO_DESCRIPTION", "Message {0} has no description. Add @desc JsDoc tag.")

var JsMessageVisitor_MESSAGE_HAS_NO_TEXT = DiagnosticType_Warning(
	"JSC_MSG_HAS_NO_TEXT",
	"Message value of {0} is just an empty string. Empty messages are forbidden.")

var JsMessageVisitor_MESSAGE_TREE_MALFORMED = DiagnosticType_Error(
	"JSC_MSG_TREE_MALFORMED", "Message parse tree malformed. {0}")

var JsMessageVisitor_MESSAGE_HAS_NO_VALUE = DiagnosticType_Error(
	"JSC_MSG_HAS_NO_VALUE", "message node {0} has no value")

var JsMessageVisitor_MESSAGE_DUPLICATE_KEY = DiagnosticType_Error(
	"JSC_MSG_KEY_DUPLICATED",
	"duplicate message variable name found for {0}, initial definition {1}:{2}")

var JsMessageVisitor_MESSAGE_NODE_IS_ORPHANED = DiagnosticType_Error(
	"JSC_MSG_ORPHANED_NODE",
	"{0}() function may be used only with MSG_* property or variable")

var JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY = DiagnosticType_Warning(
	"JSC_MSG_NOT_INITIALIZED_CORRECTLY",
	"Message must be initialized using a call to "+
		jsMessageVisitor_MSG_FUNCTION_QNAME+
		" or "+
		jsMessageVisitor_ICU_MSG_FUNCTION_QNAME)

var JsMessageVisitor_BAD_FALLBACK_SYNTAX = DiagnosticType_Error(
	"JSC_MSG_BAD_FALLBACK_SYNTAX",
	"Bad syntax. Expected syntax: "+jsMessageVisitor_MSG_FALLBACK_FUNCTION_NAME+"(MSG_1, MSG_2)")

var JsMessageVisitor_FALLBACK_ARG_ERROR = DiagnosticType_Error(
	"JSC_MSG_FALLBACK_ARG_ERROR", "Could not find message entry for fallback argument {0}")

type jsMessageLocation struct {
	message     *JsMessage
	messageNode rhino.Node
}

// malformedMessageError reports a malformed message AST node without panicking.
type malformedMessageError struct {
	msg  string
	node rhino.Node
}

func (e *malformedMessageError) Error() string {
	return "com.google.javascript.jscomp.JsMessageVisitor$MalformedException: " + e.msg
}

// JsMessageVisitor locates JS code that declares localizable messages and
// dispatches found definitions to its configured callbacks.
type JsMessageVisitor struct {
	compiler          AbstractCompiler
	idGenerator       JsMessage_IdGenerator
	onJsMessage       func(JsMessageDefinition)
	onIcuTemplate     func(IcuTemplateDefinition)
	onMessageFallback func(callNode rhino.Node, message1, message2 *JsMessage)

	// messageNames tracks named messages with their defining node for duplicate
	// detection and fallback resolution.
	messageNames map[string]jsMessageLocation

	// unnamedMessages tracks unnamed messages by Var, as they are not
	// guaranteed to be globally unique.
	unnamedMessages map[*Var]*JsMessage

	// googMsgOrder + googMsgActive hold visited goog.getMsg / declareIcuTemplate
	// call nodes that have not yet been matched to an enclosing MSG_*
	// declaration. Any nodes left at the end of traversal are orphaned.
	googMsgOrder  []rhino.Node
	googMsgActive map[rhino.Node]struct{}
}

// NewJsMessageVisitor constructs a JsMessageVisitor with the given callbacks.
func NewJsMessageVisitor(
	compiler AbstractCompiler,
	idGenerator JsMessage_IdGenerator,
	onJsMessage func(JsMessageDefinition),
	onIcuTemplate func(IcuTemplateDefinition),
	onMessageFallback func(callNode rhino.Node, message1, message2 *JsMessage),
) *JsMessageVisitor {
	v := &JsMessageVisitor{}
	v.InitJsMessageVisitor(compiler, idGenerator, onJsMessage, onIcuTemplate, onMessageFallback)
	return v
}

// InitJsMessageVisitor initializes an embedded JsMessageVisitor.
func (v *JsMessageVisitor) InitJsMessageVisitor(
	compiler AbstractCompiler,
	idGenerator JsMessage_IdGenerator,
	onJsMessage func(JsMessageDefinition),
	onIcuTemplate func(IcuTemplateDefinition),
	onMessageFallback func(callNode rhino.Node, message1, message2 *JsMessage),
) {
	v.compiler = compiler
	v.idGenerator = idGenerator
	v.onJsMessage = onJsMessage
	v.onIcuTemplate = onIcuTemplate
	v.onMessageFallback = onMessageFallback
}

// ShouldTraverse implements NodeTraversal_Callback (post-order traversal).
func (v *JsMessageVisitor) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return true
}

// Java: JsMessageVisitor.process(Node, Node)
func (v *JsMessageVisitor) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(v.compiler, root, v)

	if len(v.googMsgActive) > 0 {
		for _, msgNode := range v.googMsgOrder {
			if _, ok := v.googMsgActive[msgNode]; ok {
				v.compiler.Report(
					JSError_MakeWithNode(
						msgNode, JsMessageVisitor_MESSAGE_NODE_IS_ORPHANED, *msgNode.GetFirstChild().GetQualifiedName()))
			}
		}
	}
}

// Java: JsMessageVisitor.visit(NodeTraversal, Node, Node)
func (v *JsMessageVisitor) Visit(traversal *NodeTraversal, node rhino.Node, unused rhino.Node) {
	v.collectGetMsgCall(traversal, node)
	v.checkMessageInitialization(traversal, node)
}

func isJavaTrimEmpty(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > ' ' {
			return false
		}
	}
	return true
}

// checkMessageInitialization is called for every Node in the sources AST.
//
// Java: JsMessageVisitor.checkMessageInitialization(NodeTraversal, Node)
func (v *JsMessageVisitor) checkMessageInitialization(traversal *NodeTraversal, node rhino.Node) {
	parent := node.GetParent()

	var possiblyObfuscatedMessageKey string
	var messageKeyFromLhs string
	var msgNode rhino.Node
	var jsDocInfo *rhino.JSDocInfo

	switch node.GetToken() {
	case rhino.Token_NAME:
		// Case: `var MSG_HELLO = 'Message';`
		if parent.IsNil() || !NodeUtil_IsNameDeclaration(parent) {
			return
		}
		nm := ast.MustName(node)
		possiblyObfuscatedMessageKey = nm.Ident()
		messageKeyFromLhs = possiblyObfuscatedMessageKey
		if orig := node.GetOriginalName(); orig != nil {
			messageKeyFromLhs = *orig
		}
		if init, ok := nm.Init(); ok {
			msgNode = init.Node()
		}
		jsDocInfo = parent.GetJSDocInfo()
	case rhino.Token_ASSIGN:
		// Case: `somenamespace.someclass.MSG_HELLO = 'Message';`
		assign := ast.MustAssign(node)
		getProp, ok := ast.AsGetProp(assign.Target().Node())
		if !ok {
			return
		}
		possiblyObfuscatedMessageKey = getProp.Property()
		messageKeyFromLhs = possiblyObfuscatedMessageKey
		if orig := getProp.Node().GetOriginalName(); orig != nil {
			messageKeyFromLhs = *orig
		}
		msgNode = assign.Value().Node()
		jsDocInfo = node.GetJSDocInfo()
	case rhino.Token_STRING_KEY:
		// Case: `var t = {MSG_HELLO: 'Message'}`;
		if node.IsQuotedStringKey() || !node.HasChildren() || parent.IsObjectPattern() {
			return
		}
		jutil.CheckStateMsg1(parent.IsObjectLit(), "%s", parent)
		sk := ast.MustStringKey(node)
		possiblyObfuscatedMessageKey = sk.Key()
		messageKeyFromLhs = possiblyObfuscatedMessageKey
		if orig := node.GetOriginalName(); orig != nil {
			messageKeyFromLhs = *orig
		}
		msgNode = sk.Value().Node()
		jsDocInfo = node.GetJSDocInfo()
	case rhino.Token_MEMBER_FIELD_DEF:
		// Case: `class Foo { MSG_HELLO = 'Message'; }`
		mfd := ast.MustMemberFieldDef(node)
		possiblyObfuscatedMessageKey = mfd.Key()
		messageKeyFromLhs = possiblyObfuscatedMessageKey
		if orig := node.GetOriginalName(); orig != nil {
			messageKeyFromLhs = *orig
		}
		if val, ok := mfd.Value(); ok {
			msgNode = val.Node()
		}
		jsDocInfo = node.GetJSDocInfo()
	default:
		return
	}

	if !v.IsMessageName(messageKeyFromLhs) {
		return
	}

	if msgNode.IsNil() {
		v.compiler.Report(JSError_MakeWithNode(node, JsMessageVisitor_MESSAGE_HAS_NO_VALUE, messageKeyFromLhs))
		return
	}

	if jsMessageVisitor_isLegalMessageVarAlias(msgNode) {
		return
	}

	msgCall, msgNodeIsACall := ast.AsCall(msgNode)
	if !msgNodeIsACall {
		v.compiler.Report(JSError_MakeWithNode(node, JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY))
		return
	}
	delete(v.googMsgActive, msgNode)

	srcPtr := traversal.GetSourceName()
	mapping := v.compiler.GetSourceMapping(srcPtr, node.GetLineno(), node.GetCharno())
	var sourceName string
	if mapping != nil {
		sourceName = mapping.GetOriginalFile() + ":" + strconv.Itoa(mapping.GetLineNumber())
	} else {
		src := "null"
		if srcPtr != nil {
			src = *srcPtr
		}
		sourceName = src + ":" + strconv.Itoa(node.GetLineno())
	}

	fnNameNode := msgCall.Callee().Node()
	if jsMessageVisitor_isDeclareIcuTemplateCallee(fnNameNode) {
		icuTemplateDefinition, err := v.extractIcuTemplateDefinition(msgCall, jsDocInfo, messageKeyFromLhs, sourceName)
		if err != nil {
			v.compiler.Report(JSError_MakeWithNode(err.node, JsMessageVisitor_MESSAGE_TREE_MALFORMED, err.msg))
			return
		}
		extractedMessage := icuTemplateDefinition.GetMessage()
		v.trackMessage(traversal, possiblyObfuscatedMessageKey, msgNode, extractedMessage)
		v.reportErrorIfEmptyMessage(node, extractedMessage)
		if v.onIcuTemplate != nil {
			v.onIcuTemplate(icuTemplateDefinition)
		}
	} else if fnNameNode.MatchesQualifiedName(jsMessageVisitor_MSG_FUNCTION_QNAME) {
		jsMessageDefinition, err := v.extractJsMessageDefinition(msgCall, jsDocInfo, messageKeyFromLhs, sourceName)
		if err != nil {
			v.compiler.Report(JSError_MakeWithNode(err.node, JsMessageVisitor_MESSAGE_TREE_MALFORMED, err.msg))
			return
		}
		extractedMessage := jsMessageDefinition.GetMessage()
		v.trackMessage(traversal, possiblyObfuscatedMessageKey, msgNode, extractedMessage)
		v.reportErrorIfEmptyMessage(node, extractedMessage)

		// goog.getMsg() calls are required to have `@desc` unless they are external.
		desc, hasDesc := extractedMessage.DescOpt()
		if (!hasDesc || isJavaTrimEmpty(desc)) && !extractedMessage.IsExternal() {
			v.compiler.Report(
				JSError_MakeWithNode(node, JsMessageVisitor_MESSAGE_HAS_NO_DESCRIPTION, extractedMessage.GetKey()))
		}

		if v.onJsMessage != nil {
			v.onJsMessage(jsMessageDefinition)
		}
	} else {
		v.compiler.Report(
			JSError_MakeWithNode(
				msgNode,
				JsMessageVisitor_MESSAGE_TREE_MALFORMED,
				"Message must be initialized using a call to "+
					jsMessageVisitor_MSG_FUNCTION_QNAME+
					" or "+
					jsMessageVisitor_ICU_MSG_FUNCTION_NAME+
					" (from goog.i18n.messages)."))
	}
}

// Java: JsMessageVisitor.isDeclareIcuTemplateCallee(Node)
func jsMessageVisitor_isDeclareIcuTemplateCallee(calleeNode rhino.Node) bool {
	if nm, ok := ast.AsName(calleeNode); ok {
		return strings.HasSuffix(nm.Ident(), jsMessageVisitor_ICU_MSG_FUNCTION_NAME)
	}
	if gp, ok := ast.AsGetProp(calleeNode); ok {
		return strings.HasSuffix(gp.Property(), jsMessageVisitor_ICU_MSG_FUNCTION_NAME) && calleeNode.IsQualifiedName()
	}
	return false
}

type jsMessageBuildInput struct {
	messageKeyFromLhs            string
	sourceName                   string
	messageText                  string
	description                  string
	meaning                      string
	alternateMessageId           string
	messageParts                 []JsMessage_Part
	placeholderExampleMap        map[string]string
	placeholderExampleOrder      []string
	placeholderOriginalCodeMap   map[string]string
	placeholderOriginalCodeOrder []string
	hasDescription               bool
	hasMeaning                   bool
	hasAlternateMessageId        bool
}

// Java: JsMessageVisitor.buildJsMessage(InputForBuildJsMsg)
func (v *JsMessageVisitor) buildJsMessage(input jsMessageBuildInput) *JsMessage {
	messageParts := input.messageParts
	messageKeyFromLhs := input.messageKeyFromLhs

	externalId, isExternal := jsMessageVisitor_externalMessageId(messageKeyFromLhs)

	var isAnonymous bool
	var messageId string
	var messageKeyFinal string
	if isExternal {
		// MSG_EXTERNAL_12345 = ...
		messageId = externalId
		messageKeyFinal = messageKeyFromLhs
	} else {
		meaningForIdGeneration := input.meaning
		if jsMessageVisitor_isUnnamedMessageName(messageKeyFromLhs) {
			isAnonymous = true
			messageKeyFinal = jsMessageVisitor_generateKeyFromMessageText(input.messageText)
			if !input.hasMeaning {
				meaningForIdGeneration = messageKeyFinal
			}
		} else {
			messageKeyFinal = messageKeyFromLhs
			if !input.hasMeaning {
				meaningForIdGeneration = JsMessageVisitor_RemoveScopedAliasesPrefix(messageKeyFromLhs)
			}
		}
		if v.idGenerator == nil {
			messageId = meaningForIdGeneration
		} else {
			messageId = generateMessageId(v.idGenerator, meaningForIdGeneration, messageParts)
		}
	}

	return NewJsMessage_Builder().
		AppendParts(messageParts).
		setPlaceholderNameToExampleOrdered(input.placeholderExampleMap, input.placeholderExampleOrder).
		setPlaceholderNameToOriginalCodeOrdered(input.placeholderOriginalCodeMap, input.placeholderOriginalCodeOrder).
		SetIsAnonymous(isAnonymous).
		SetIsExternalMsg(isExternal).
		SetKey(messageKeyFinal).
		SetSourceNameString(input.sourceName).
		SetDescOpt(input.description, input.hasDescription).
		SetMeaningOpt(input.meaning, input.hasMeaning).
		SetAlternateIdOpt(input.alternateMessageId, input.hasAlternateMessageId).
		SetId(messageId).
		Build()
}

// Java: JsMessageVisitor.trackMessage(NodeTraversal, String, Node, JsMessage)
func (v *JsMessageVisitor) trackMessage(
	traversal *NodeTraversal,
	possiblyObfuscatedMessageKey string,
	msgNode rhino.Node,
	extractedMessage *JsMessage,
) {
	if !extractedMessage.IsAnonymous() && !extractedMessage.IsExternal() {
		v.checkIfMessageDuplicated(extractedMessage.GetKey(), msgNode)
	}
	if extractedMessage.IsAnonymous() {
		v.trackUnnamedMessage(traversal, extractedMessage, possiblyObfuscatedMessageKey)
	} else {
		v.trackNormalMessage(extractedMessage, extractedMessage.GetKey(), msgNode)
	}
}

// Java: JsMessageVisitor.reportErrorIfEmptyMessage(Node, JsMessage)
func (v *JsMessageVisitor) reportErrorIfEmptyMessage(node rhino.Node, extractedMessage *JsMessage) {
	if extractedMessage.IsEmpty() {
		v.compiler.Report(JSError_MakeWithNode(node, JsMessageVisitor_MESSAGE_HAS_NO_TEXT, extractedMessage.GetKey()))
	}
}

func jsMessageVisitor_externalMessageId(messageKey string) (string, bool) {
	if strings.HasPrefix(messageKey, jsMessageVisitor_MSG_EXTERNAL_PREFIX) {
		start := len(jsMessageVisitor_MSG_EXTERNAL_PREFIX)
		end := start
		for ; end < len(messageKey); end++ {
			c := messageKey[end]
			if c < '0' || c > '9' {
				break
			}
		}
		if end > start {
			return messageKey[start:end], true
		}
	}
	return "", false
}

// JsMessageVisitor_GetExternalMessageId extracts an external message ID from
// the message key, if it contains one, or nil otherwise.
//
// Java: JsMessageVisitor.getExternalMessageId(String)
func JsMessageVisitor_GetExternalMessageId(messageKey string) *string {
	if id, ok := jsMessageVisitor_externalMessageId(messageKey); ok {
		return &id
	}
	return nil
}

// Java: JsMessageVisitor.generateKeyFromMessageText(String)
func jsMessageVisitor_generateKeyFromMessageText(msgText string) string {
	nonnegativeHash := int64(0x7fffffffffffffff) & jsMessage_Hash_hash64String(msgText)
	return JsMessageVisitor_MSG_PREFIX + jutil.StringToUpperCaseASCII(strconv.FormatInt(nonnegativeHash, 36))
}

// Java: JsMessageVisitor.collectGetMsgCall(NodeTraversal, Node)
func (v *JsMessageVisitor) collectGetMsgCall(traversal *NodeTraversal, n rhino.Node) {
	call, ok := ast.AsCall(n)
	if !ok {
		return
	}

	callee := call.Callee().Node()
	if callee.MatchesQualifiedName(jsMessageVisitor_MSG_FUNCTION_QNAME) || jsMessageVisitor_isDeclareIcuTemplateCallee(callee) {
		if v.googMsgActive == nil {
			v.googMsgActive = make(map[rhino.Node]struct{})
		}
		if _, exists := v.googMsgActive[n]; !exists {
			v.googMsgActive[n] = struct{}{}
			v.googMsgOrder = append(v.googMsgOrder, n)
		}
	} else if callee.MatchesQualifiedName(jsMessageVisitor_MSG_FALLBACK_FUNCTION_NAME) {
		v.visitFallbackFunctionCall(traversal, call)
	}
}

// Java: JsMessageVisitor.trackNormalMessage(JsMessage, String, Node)
func (v *JsMessageVisitor) trackNormalMessage(message *JsMessage, msgName string, msgNode rhino.Node) {
	if v.messageNames == nil {
		v.messageNames = make(map[string]jsMessageLocation)
	}
	v.messageNames[msgName] = jsMessageLocation{message: message, messageNode: msgNode}
}

// Java: JsMessageVisitor.trackUnnamedMessage(NodeTraversal, JsMessage, String)
func (v *JsMessageVisitor) trackUnnamedMessage(t *NodeTraversal, message *JsMessage, msgNameInScope string) {
	if var_ := t.GetScope().GetVar(msgNameInScope); var_ != nil {
		if v.unnamedMessages == nil {
			v.unnamedMessages = make(map[*Var]*JsMessage)
		}
		v.unnamedMessages[var_] = message
	}
}

// Java: JsMessageVisitor.isLegalMessageVarAlias(Node)
func jsMessageVisitor_isLegalMessageVarAlias(msgNode rhino.Node) bool {
	if gp, ok := ast.AsGetProp(msgNode); ok &&
		msgNode.IsQualifiedName() &&
		strings.HasPrefix(gp.Property(), JsMessageVisitor_MSG_PREFIX) {
		return true
	}

	nm, ok := ast.AsName(msgNode)
	if !ok {
		return false
	}

	originalName := nm.Ident()
	if orig := msgNode.GetOriginalName(); orig != nil {
		originalName = *orig
	}

	return strings.HasPrefix(originalName, JsMessageVisitor_MSG_PREFIX)
}

// Java: JsMessageVisitor.getTrackedUnnamedMessage(NodeTraversal, String)
func (v *JsMessageVisitor) getTrackedUnnamedMessage(t *NodeTraversal, msgNameInScope string) *JsMessage {
	if var_ := t.GetScope().GetVar(msgNameInScope); var_ != nil {
		return v.unnamedMessages[var_]
	}
	return nil
}

// Java: JsMessageVisitor.getTrackedNormalMessage(String)
func (v *JsMessageVisitor) getTrackedNormalMessage(msgName string) *JsMessage {
	if location, ok := v.messageNames[msgName]; ok {
		return location.message
	}
	return nil
}

// Java: JsMessageVisitor.checkIfMessageDuplicated(String, Node)
func (v *JsMessageVisitor) checkIfMessageDuplicated(msgName string, msgNode rhino.Node) {
	if location, ok := v.messageNames[msgName]; ok {
		srcFile := "null"
		if s := location.messageNode.GetSourceFileName(); s != nil {
			srcFile = *s
		}
		v.compiler.Report(
			JSError_MakeWithNode(
				msgNode,
				JsMessageVisitor_MESSAGE_DUPLICATE_KEY,
				msgName,
				srcFile,
				strconv.Itoa(location.messageNode.GetLineno())))
	}
}

// jsMessageVisitor_extractStringFromStringExprNode returns the string value
// associated with a JS string literal, no-substitution template literal, or
// string concatenation.
//
// Java: JsMessageVisitor.extractStringFromStringExprNode(Node)
func jsMessageVisitor_extractStringFromStringExprNode(node rhino.Node) (string, *malformedMessageError) {
	switch node.GetToken() {
	case rhino.Token_STRINGLIT:
		return ast.MustStringLit(node).Value(), nil
	case rhino.Token_TEMPLATELIT:
		if node.HasOneChild() {
			return *jutil.CheckNotNull(node.GetFirstChild().GetCookedString()), nil
		}
		return "", &malformedMessageError{
			msg:  "Template literals with substitutions are not allowed.",
			node: node,
		}
	case rhino.Token_ADD:
		var sb strings.Builder
		for child := node.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			part, err := jsMessageVisitor_extractStringFromStringExprNode(child)
			if err != nil {
				return "", err
			}
			sb.WriteString(part)
		}
		return sb.String(), nil
	default:
		return "", &malformedMessageError{msg: "literal string or concatenation expected", node: node}
	}
}

// extractJsMessageDefinition extracts data from a call to goog.getMsg(...).
//
// Java: JsMessageVisitor.extractJsMessageDefinition(Node, JSDocInfo, String, String)
func (v *JsMessageVisitor) extractJsMessageDefinition(
	msgCall ast.Call, jsDocInfo *rhino.JSDocInfo, messageKeyFromLhs string, sourceName string,
) (JsMessageDefinition, *malformedMessageError) {
	msgNode := msgCall.Node()
	var description, meaning, alternateMessageId string
	var hasDescription, hasMeaning, hasAlternateMessageId bool
	if jsDocInfo != nil {
		if s := jsDocInfo.GetDescription(); s != nil {
			description = *s
			hasDescription = true
		}
		if s := jsDocInfo.GetMeaning(); s != nil {
			meaning = *s
			hasMeaning = true
		}
		if s := jsDocInfo.GetAlternateMessageId(); s != nil {
			alternateMessageId = *s
			hasAlternateMessageId = true
		}
	}

	firstArg, hasFirstArg := msgCall.FirstArg()
	if !hasFirstArg {
		return JsMessageDefinition{}, &malformedMessageError{msg: "Message string literal expected", node: msgNode}
	}
	msgTextNode := firstArg.Node()
	valuesObjLit := msgTextNode.GetNext()
	var optionsBagArgument rhino.Node
	if !valuesObjLit.IsNil() {
		optionsBagArgument = valuesObjLit.GetNext()
	}
	if !optionsBagArgument.IsNil() {
		if unexpectedArgument := optionsBagArgument.GetNext(); !unexpectedArgument.IsNil() {
			return JsMessageDefinition{}, &malformedMessageError{msg: "too many arguments", node: unexpectedArgument}
		}
	}

	msgTextString, malErr := jsMessageVisitor_extractStringFromStringExprNode(msgTextNode)
	if malErr != nil {
		return JsMessageDefinition{}, malErr
	}
	googGetMsgParsedText, fmtErr := jsMessageVisitor_extractGoogGetMsgParsedText(msgTextString)
	if fmtErr != nil {
		return JsMessageDefinition{}, &malformedMessageError{msg: fmtErr.Error(), node: msgTextNode}
	}

	placeholderObjectLiteralMap, malErr := extractObjectLiteralMap(valuesObjLit)
	if malErr != nil {
		return JsMessageDefinition{}, malErr
	}

	if malErr := placeholderObjectLiteralMap.checkForRequiredKeys(
		googGetMsgParsedText.placeholderOrder, msgTextNode, "Unrecognized message placeholder referenced: ",
	); malErr != nil {
		return JsMessageDefinition{}, malErr
	}
	if malErr := placeholderObjectLiteralMap.checkForUnexpectedKeys(
		googGetMsgParsedText.placeholderSet, "Unused message placeholder: ",
	); malErr != nil {
		return JsMessageDefinition{}, malErr
	}

	placeholderValuesMap, placeholderValuesOrder := placeholderObjectLiteralMap.extractAsValueMapAndOrder()

	jsMessageOptions, malErr := extractJsMessageOptions(optionsBagArgument)
	if malErr != nil {
		return JsMessageDefinition{}, malErr
	}
	if malErr := jsMessageOptions.checkForUnknownPlaceholders(googGetMsgParsedText.placeholderSet); malErr != nil {
		return JsMessageDefinition{}, malErr
	}

	jsExtractedMessage := v.buildJsMessage(jsMessageBuildInput{
		messageKeyFromLhs:            messageKeyFromLhs,
		sourceName:                   sourceName,
		messageText:                  googGetMsgParsedText.originalMsgText,
		messageParts:                 googGetMsgParsedText.parts,
		placeholderExampleMap:        jsMessageOptions.placeholderExamplesMap,
		placeholderExampleOrder:      jsMessageOptions.placeholderExamplesOrder,
		placeholderOriginalCodeMap:   jsMessageOptions.placeholderOriginalCodeMap,
		placeholderOriginalCodeOrder: jsMessageOptions.placeholderOriginalCodeOrder,
		description:                  description,
		hasDescription:               hasDescription,
		meaning:                      meaning,
		hasMeaning:                   hasMeaning,
		alternateMessageId:           alternateMessageId,
		hasAlternateMessageId:        hasAlternateMessageId,
	})

	return JsMessageDefinition{
		message:               jsExtractedMessage,
		messageNode:           msgNode,
		templateTextNode:      msgTextNode,
		placeholderValuesNode: valuesObjLit,
		placeholderValueMap:   placeholderValuesMap,
		placeholderValueOrder: placeholderValuesOrder,
		escapeLessThan:        jsMessageOptions.isEscapeLessThan,
		unescapeHtmlEntities:  jsMessageOptions.isUnescapeHtmlEntities,
	}, nil
}

// extractIcuTemplateDefinition extracts message data from a
// `declareIcuTemplate()` call.
//
// Java: JsMessageVisitor.extractIcuTemplateDefinition(Node, JSDocInfo, String, String)
func (v *JsMessageVisitor) extractIcuTemplateDefinition(
	msgCall ast.Call, jsDocInfo *rhino.JSDocInfo, messageKeyFromLhs string, sourceName string,
) (IcuTemplateDefinition, *malformedMessageError) {
	msgNode := msgCall.Node()
	if jsDocInfo != nil {
		if jsDocInfo.GetAlternateMessageId() != nil {
			return IcuTemplateDefinition{}, &malformedMessageError{
				msg:  "Use the 'alternate_message_id' option, not the '@alternateMessageId' annotation",
				node: msgNode,
			}
		}
		if jsDocInfo.GetDescription() != nil {
			return IcuTemplateDefinition{}, &malformedMessageError{
				msg:  "Use the 'description' option, not the '@desc' annotation",
				node: msgNode,
			}
		}
		if jsDocInfo.GetMeaning() != nil {
			return IcuTemplateDefinition{}, &malformedMessageError{
				msg:  "Use the 'meaning' option, not the '@meaning' annotation",
				node: msgNode,
			}
		}
	}

	firstArg, hasFirstArg := msgCall.FirstArg()
	if !hasFirstArg {
		return IcuTemplateDefinition{}, &malformedMessageError{msg: "message string argument expected", node: msgNode}
	}
	stringLiteralExpression := firstArg.Node()
	templateString, malErr := jsMessageVisitor_extractStringFromStringExprNode(stringLiteralExpression)
	if malErr != nil {
		return IcuTemplateDefinition{}, malErr
	}
	icuMessageTemplateString := JsMessageVisitor_IcuMessageTemplateString{template: templateString}

	optionsBagNode := stringLiteralExpression.GetNext()
	if optionsBagNode.IsNil() {
		return IcuTemplateDefinition{}, &malformedMessageError{msg: "options argument expected", node: msgNode}
	}
	icuOptions, malErr := extractIcuTemplateOptions(optionsBagNode)
	if malErr != nil {
		return IcuTemplateDefinition{}, malErr
	}

	extractedIcuTemplateParts := icuMessageTemplateString.ExtractParts(icuOptions.placeholderNames)
	if malErr := icuOptions.checkForUnknownPlaceholders(extractedIcuTemplateParts.extractedPlaceholderNames); malErr != nil {
		return IcuTemplateDefinition{}, malErr
	}

	if unexpectedArgument := optionsBagNode.GetNext(); !unexpectedArgument.IsNil() {
		return IcuTemplateDefinition{}, &malformedMessageError{msg: "too many arguments", node: unexpectedArgument}
	}

	extractedMessage := v.buildJsMessage(jsMessageBuildInput{
		messageKeyFromLhs:            messageKeyFromLhs,
		sourceName:                   sourceName,
		messageText:                  icuMessageTemplateString.template,
		messageParts:                 extractedIcuTemplateParts.extractedParts,
		placeholderExampleMap:        icuOptions.placeholderExamplesMap,
		placeholderExampleOrder:      icuOptions.placeholderExamplesOrder,
		placeholderOriginalCodeMap:   icuOptions.placeholderOriginalCodeMap,
		placeholderOriginalCodeOrder: icuOptions.placeholderOriginalCodeOrder,
		description:                  icuOptions.description,
		hasDescription:               true,
		meaning:                      icuOptions.meaning,
		hasMeaning:                   icuOptions.hasMeaning,
		alternateMessageId:           icuOptions.alternateMessageId,
		hasAlternateMessageId:        icuOptions.hasAlternateMessageId,
	})

	return IcuTemplateDefinition{
		message:          extractedMessage,
		messageNode:      msgNode,
		templateTextNode: stringLiteralExpression,
	}, nil
}

// JsMessageVisitor_IcuMessageTemplateString is the template string of a
// declareIcuTemplate() call.
type JsMessageVisitor_IcuMessageTemplateString struct {
	template string
}

// Java: JsMessageVisitor.IcuMessageTemplateString(String)
func NewJsMessageVisitor_IcuMessageTemplateString(template string) JsMessageVisitor_IcuMessageTemplateString {
	return JsMessageVisitor_IcuMessageTemplateString{template: template}
}

// ExtractParts splits the template string into plain string parts and
// placeholders named in placeholderNamesToExtract.
//
// Java: JsMessageVisitor.IcuMessageTemplateString.extractParts(Set<String>)
func (s JsMessageVisitor_IcuMessageTemplateString) ExtractParts(
	placeholderNamesToExtract map[string]struct{},
) JsMessageVisitor_ExtractedIcuTemplateParts {
	tpl := s.template
	if len(tpl) == 0 {
		return JsMessageVisitor_ExtractedIcuTemplateParts{}
	}
	if len(placeholderNamesToExtract) == 0 {
		return JsMessageVisitor_ExtractedIcuTemplateParts{
			extractedParts: []JsMessage_Part{JsMessage_StringPart_Create(tpl)},
		}
	}

	var result JsMessageVisitor_ExtractedIcuTemplateParts
	remainingStart := 0
	for i := 0; i < len(tpl); {
		rel := strings.IndexByte(tpl[i:], '{')
		if rel < 0 {
			break
		}
		openIdx := i + rel
		j := openIdx + 1
		for j < len(tpl) {
			c := tpl[j]
			if (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
				j++
			} else {
				break
			}
		}
		if j > openIdx+1 && j < len(tpl) && tpl[j] == '}' {
			placeholderName := tpl[openIdx+1 : j]
			if _, ok := placeholderNamesToExtract[placeholderName]; ok {
				if result.extractedPlaceholderNames == nil {
					result.extractedPlaceholderNames = make(map[string]struct{})
				}
				result.extractedPlaceholderNames[placeholderName] = struct{}{}
				result.extractedParts = append(result.extractedParts,
					JsMessage_StringPart_Create(tpl[remainingStart:openIdx]),
					JsMessage_PlaceholderReference_CreateForCanonicalName(placeholderName),
				)
				remainingStart = j + 1
			}
			i = j + 1
		} else {
			i = j
			if i == openIdx+1 && i < len(tpl) && tpl[i] != '{' {
				i++
			}
		}
	}
	if remainingStart < len(tpl) {
		result.extractedParts = append(result.extractedParts, JsMessage_StringPart_Create(tpl[remainingStart:]))
	}
	return result
}

// JsMessageVisitor_ExtractedIcuTemplateParts is the result of
// IcuMessageTemplateString.ExtractParts.
type JsMessageVisitor_ExtractedIcuTemplateParts struct {
	extractedPlaceholderNames map[string]struct{}
	extractedParts            []JsMessage_Part
}

func isJsMessageOptionName(key string) bool {
	switch key {
	case "html", "unescapeHtmlEntities", "example", "original_code":
		return true
	}
	return false
}

type jsMessageOptions struct {
	placeholderExamplesMap       map[string]string
	placeholderExamplesOrder     []string
	placeholderOriginalCodeMap   map[string]string
	placeholderOriginalCodeOrder []string
	exampleObjectLiteralMap      jsMessageObjectLiteralMap
	originalCodeObjectLiteralMap jsMessageObjectLiteralMap
	isEscapeLessThan             bool
	isUnescapeHtmlEntities       bool
}

func (o *jsMessageOptions) checkForUnknownPlaceholders(knownPlaceholders map[string]struct{}) *malformedMessageError {
	if err := o.exampleObjectLiteralMap.checkForUnexpectedKeys(knownPlaceholders, "Unknown placeholder: "); err != nil {
		return err
	}
	return o.originalCodeObjectLiteralMap.checkForUnexpectedKeys(knownPlaceholders, "Unknown placeholder: ")
}

// Java: JsMessageVisitor.extractJsMessageOptions(Node)
func extractJsMessageOptions(optionsBag rhino.Node) (jsMessageOptions, *malformedMessageError) {
	if optionsBag.IsNil() {
		return jsMessageOptions{}, nil
	}
	objectLiteralMap, err := extractObjectLiteralMap(optionsBag)
	if err != nil {
		return jsMessageOptions{}, err
	}
	if err := objectLiteralMap.checkForUnexpectedKeysFn(isJsMessageOptionName, "Unknown option: "); err != nil {
		return jsMessageOptions{}, err
	}

	isEscapeLessThan, err := objectLiteralMap.getBooleanValueOrFalse("html")
	if err != nil {
		return jsMessageOptions{}, err
	}
	isUnescapeHtmlEntities, err := objectLiteralMap.getBooleanValueOrFalse("unescapeHtmlEntities")
	if err != nil {
		return jsMessageOptions{}, err
	}

	exampleObjectLiteralMap, err := extractObjectLiteralMap(objectLiteralMap.getValueNode("example"))
	if err != nil {
		return jsMessageOptions{}, err
	}
	placeholderExamplesMap, placeholderExamplesOrder, err := exampleObjectLiteralMap.extractAsStringToStringMap()
	if err != nil {
		return jsMessageOptions{}, err
	}

	originalCodeObjectLiteralMap, err := extractObjectLiteralMap(objectLiteralMap.getValueNode("original_code"))
	if err != nil {
		return jsMessageOptions{}, err
	}
	placeholderOriginalCodeMap, placeholderOriginalCodeOrder, err := originalCodeObjectLiteralMap.extractAsStringToStringMap()
	if err != nil {
		return jsMessageOptions{}, err
	}

	return jsMessageOptions{
		isEscapeLessThan:             isEscapeLessThan,
		isUnescapeHtmlEntities:       isUnescapeHtmlEntities,
		placeholderExamplesMap:       placeholderExamplesMap,
		placeholderExamplesOrder:     placeholderExamplesOrder,
		placeholderOriginalCodeMap:   placeholderOriginalCodeMap,
		placeholderOriginalCodeOrder: placeholderOriginalCodeOrder,
		exampleObjectLiteralMap:      exampleObjectLiteralMap,
		originalCodeObjectLiteralMap: originalCodeObjectLiteralMap,
	}, nil
}

func isIcuTemplateOptionName(key string) bool {
	switch key {
	case "description", "meaning", "alternate_message_id", "example", "original_code":
		return true
	}
	return false
}

type icuTemplateOptions struct {
	description                  string
	meaning                      string
	alternateMessageId           string
	placeholderExamplesMap       map[string]string
	placeholderExamplesOrder     []string
	placeholderOriginalCodeMap   map[string]string
	placeholderOriginalCodeOrder []string
	placeholderNames             map[string]struct{}
	exampleObjectLiteralMap      jsMessageObjectLiteralMap
	originalCodeObjectLiteralMap jsMessageObjectLiteralMap
	hasMeaning                   bool
	hasAlternateMessageId        bool
}

func (o *icuTemplateOptions) checkForUnknownPlaceholders(knownPlaceholders map[string]struct{}) *malformedMessageError {
	if err := o.exampleObjectLiteralMap.checkForUnexpectedKeys(knownPlaceholders, "Unknown placeholder: "); err != nil {
		return err
	}
	return o.originalCodeObjectLiteralMap.checkForUnexpectedKeys(knownPlaceholders, "Unknown placeholder: ")
}

// Java: JsMessageVisitor.extractIcuTemplateOptions(Node)
func extractIcuTemplateOptions(optionsBag rhino.Node) (icuTemplateOptions, *malformedMessageError) {
	objectLiteralMap, err := extractObjectLiteralMap(optionsBag)
	if err != nil {
		return icuTemplateOptions{}, err
	}
	if err := objectLiteralMap.checkForUnexpectedKeysFn(isIcuTemplateOptionName, "Unknown option: "); err != nil {
		return icuTemplateOptions{}, err
	}

	descriptionNode := objectLiteralMap.getValueNode("description")
	if descriptionNode.IsNil() {
		return icuTemplateOptions{}, &malformedMessageError{msg: "'description' option field is missing", node: optionsBag}
	}
	description, err := jsMessageVisitor_extractStringFromStringExprNode(descriptionNode)
	if err != nil {
		return icuTemplateOptions{}, err
	}

	var meaning string
	var hasMeaning bool
	if meaningNode := objectLiteralMap.getValueNode("meaning"); !meaningNode.IsNil() {
		m, err := jsMessageVisitor_extractStringFromStringExprNode(meaningNode)
		if err != nil {
			return icuTemplateOptions{}, err
		}
		meaning = m
		hasMeaning = true
	}

	var alternateMessageId string
	var hasAlternateMessageId bool
	if alternateMessageIdNode := objectLiteralMap.getValueNode("alternate_message_id"); !alternateMessageIdNode.IsNil() {
		alt, err := jsMessageVisitor_extractStringFromStringExprNode(alternateMessageIdNode)
		if err != nil {
			return icuTemplateOptions{}, err
		}
		alternateMessageId = alt
		hasAlternateMessageId = true
	}

	exampleObjectLiteralMap, err := extractObjectLiteralMap(objectLiteralMap.getValueNode("example"))
	if err != nil {
		return icuTemplateOptions{}, err
	}
	placeholderExamplesMap, placeholderExamplesOrder, err := exampleObjectLiteralMap.extractAsStringToStringMap()
	if err != nil {
		return icuTemplateOptions{}, err
	}

	originalCodeObjectLiteralMap, err := extractObjectLiteralMap(objectLiteralMap.getValueNode("original_code"))
	if err != nil {
		return icuTemplateOptions{}, err
	}
	placeholderOriginalCodeMap, placeholderOriginalCodeOrder, err := originalCodeObjectLiteralMap.extractAsStringToStringMap()
	if err != nil {
		return icuTemplateOptions{}, err
	}

	var placeholderNames map[string]struct{}
	if len(placeholderExamplesOrder) > 0 || len(placeholderOriginalCodeOrder) > 0 {
		placeholderNames = make(map[string]struct{}, len(placeholderExamplesOrder)+len(placeholderOriginalCodeOrder))
		for _, k := range placeholderExamplesOrder {
			if !JsMessage_IsCanonicalPlaceholderNameFormat(k) {
				return icuTemplateOptions{}, &malformedMessageError{
					msg:  "Placeholder not in UPPER_SNAKE_CASE: " + k,
					node: optionsBag,
				}
			}
			placeholderNames[k] = struct{}{}
		}
		for _, k := range placeholderOriginalCodeOrder {
			if _, seen := placeholderNames[k]; !seen && !JsMessage_IsCanonicalPlaceholderNameFormat(k) {
				return icuTemplateOptions{}, &malformedMessageError{
					msg:  "Placeholder not in UPPER_SNAKE_CASE: " + k,
					node: optionsBag,
				}
			}
			placeholderNames[k] = struct{}{}
		}
	}

	return icuTemplateOptions{
		description:                  description,
		meaning:                      meaning,
		hasMeaning:                   hasMeaning,
		alternateMessageId:           alternateMessageId,
		hasAlternateMessageId:        hasAlternateMessageId,
		placeholderExamplesMap:       placeholderExamplesMap,
		placeholderExamplesOrder:     placeholderExamplesOrder,
		placeholderOriginalCodeMap:   placeholderOriginalCodeMap,
		placeholderOriginalCodeOrder: placeholderOriginalCodeOrder,
		placeholderNames:             placeholderNames,
		exampleObjectLiteralMap:      exampleObjectLiteralMap,
		originalCodeObjectLiteralMap: originalCodeObjectLiteralMap,
	}, nil
}

type jsMessageObjEntry struct {
	key  string
	node rhino.Node
}

// jsMessageObjectLiteralMap represents an AST OBJECTLIT whose members are all
// unique STRING_KEY nodes. Small object literals (<=8 keys) avoid allocating a
// hash map altogether; larger ones lazily populate byKey.
type jsMessageObjectLiteralMap struct {
	entries []jsMessageObjEntry
	byKey   map[string]rhino.Node
}

func (m *jsMessageObjectLiteralMap) findNode(key string) (rhino.Node, bool) {
	if m.byKey != nil {
		n, ok := m.byKey[key]
		return n, ok
	}
	for _, e := range m.entries {
		if e.key == key {
			return e.node, true
		}
	}
	return rhino.NilNode, false
}

func extractObjectLiteralMap(objLit rhino.Node) (jsMessageObjectLiteralMap, *malformedMessageError) {
	if objLit.IsNil() {
		return jsMessageObjectLiteralMap{}, nil
	}
	lit, ok := ast.AsObjectLit(objLit)
	if !ok {
		return jsMessageObjectLiteralMap{}, &malformedMessageError{msg: "object literal expected", node: objLit}
	}
	var m jsMessageObjectLiteralMap
	for member := range lit.Members() {
		sk, ok := ast.AsStringKey(member.Node())
		if !ok {
			return jsMessageObjectLiteralMap{}, &malformedMessageError{msg: "string key expected", node: member.Node()}
		}
		key := sk.Key()
		if _, dup := m.findNode(key); dup {
			return jsMessageObjectLiteralMap{}, &malformedMessageError{msg: "duplicate string key: " + key, node: member.Node()}
		}
		if len(m.entries) == 8 && m.byKey == nil {
			m.byKey = make(map[string]rhino.Node, 16)
			for _, prev := range m.entries {
				m.byKey[prev.key] = prev.node
			}
		}
		m.entries = append(m.entries, jsMessageObjEntry{key: key, node: member.Node()})
		if m.byKey != nil {
			m.byKey[key] = member.Node()
		}
	}
	return m, nil
}

func (m *jsMessageObjectLiteralMap) getValueNode(key string) rhino.Node {
	stringKeyNode, ok := m.findNode(key)
	if !ok {
		return rhino.NilNode
	}
	return stringKeyNode.GetOnlyChild()
}

func (m *jsMessageObjectLiteralMap) getBooleanValueOrFalse(key string) (bool, *malformedMessageError) {
	stringKeyNode, ok := m.findNode(key)
	if !ok {
		return false, nil
	}
	valueNode := stringKeyNode.GetOnlyChild()
	if valueNode.IsTrue() {
		return true, nil
	}
	if valueNode.IsFalse() {
		return false, nil
	}
	return false, &malformedMessageError{
		msg:  stringKeyNode.GetString() + ": Literal true or false expected",
		node: valueNode,
	}
}

func (m *jsMessageObjectLiteralMap) extractAsValueMap() *jutil.LinkedHashMap[string, rhino.Node] {
	if len(m.entries) == 0 {
		return nil
	}
	out := jutil.NewLinkedHashMapWithCapacity[string, rhino.Node](len(m.entries))
	for _, e := range m.entries {
		out.PutNew(e.key, e.node.GetOnlyChild())
	}
	return out
}

func (m *jsMessageObjectLiteralMap) extractAsValueMapAndOrder() (map[string]rhino.Node, []string) {
	if len(m.entries) == 0 {
		return nil, nil
	}
	out := make(map[string]rhino.Node, len(m.entries))
	order := make([]string, len(m.entries))
	for i, e := range m.entries {
		out[e.key] = e.node.GetOnlyChild()
		order[i] = e.key
	}
	return out, order
}

func (m *jsMessageObjectLiteralMap) extractAsStringToStringMap() (map[string]string, []string, *malformedMessageError) {
	if len(m.entries) == 0 {
		return nil, nil, nil
	}
	out := make(map[string]string, len(m.entries))
	order := make([]string, len(m.entries))
	for i, e := range m.entries {
		val, err := jsMessageVisitor_extractStringFromStringExprNode(e.node.GetOnlyChild())
		if err != nil {
			return nil, nil, err
		}
		out[e.key] = val
		order[i] = e.key
	}
	return out, order, nil
}

func (m *jsMessageObjectLiteralMap) checkForUnexpectedKeys(
	expectedKeys map[string]struct{}, errorPrefix string,
) *malformedMessageError {
	for _, e := range m.entries {
		if _, ok := expectedKeys[e.key]; !ok {
			return &malformedMessageError{msg: errorPrefix + e.key, node: e.node}
		}
	}
	return nil
}

func (m *jsMessageObjectLiteralMap) checkForUnexpectedKeysFn(
	isExpected func(string) bool, errorPrefix string,
) *malformedMessageError {
	for _, e := range m.entries {
		if !isExpected(e.key) {
			return &malformedMessageError{msg: errorPrefix + e.key, node: e.node}
		}
	}
	return nil
}

func (m *jsMessageObjectLiteralMap) checkForRequiredKeys(
	requiredKeys []string, errNode rhino.Node, errorPrefix string,
) *malformedMessageError {
	for _, requiredKey := range requiredKeys {
		if _, ok := m.findNode(requiredKey); !ok {
			return &malformedMessageError{msg: errorPrefix + requiredKey, node: errNode}
		}
	}
	return nil
}

type googGetMsgParsedText struct {
	originalMsgText  string
	placeholderOrder []string
	placeholderSet   map[string]struct{}
	parts            []JsMessage_Part
}

func jsMessageVisitor_parseJsMessageTextIntoParts(originalMsgText string) ([]JsMessage_Part, error) {
	parsed, err := jsMessageVisitor_extractGoogGetMsgParsedText(originalMsgText)
	if err != nil {
		return nil, err
	}
	return parsed.parts, nil
}

var errPlaceholderMalformed = errors.New("Placeholder incorrectly formatted")

// Java: JsMessageVisitor.extractGoogGetMsgParsedText(String)
func jsMessageVisitor_extractGoogGetMsgParsedText(originalMsgText string) (googGetMsgParsedText, error) {
	msgText := originalMsgText
	var partsBuilder []JsMessage_Part
	var placeholderOrder []string
	var placeholderSet map[string]struct{}
	for {
		phBegin := strings.Index(msgText, JsMessage_PH_JS_PREFIX)
		if phBegin < 0 {
			partsBuilder = append(partsBuilder, JsMessage_StringPart_Create(msgText))
			break
		}
		if phBegin > 0 {
			partsBuilder = append(partsBuilder, JsMessage_StringPart_Create(msgText[:phBegin]))
		}

		afterPrefix := msgText[phBegin+len(JsMessage_PH_JS_PREFIX):]
		phEndRel := strings.Index(afterPrefix, JsMessage_PH_JS_SUFFIX)
		if phEndRel < 0 {
			return googGetMsgParsedText{}, errPlaceholderMalformed
		}

		phName := afterPrefix[:phEndRel]
		if !JsMessage_IsLowerCamelCaseWithNumericSuffixes(phName) {
			return googGetMsgParsedText{}, errors.New("Placeholder name not in lowerCamelCase: " + phName)
		}
		if placeholderSet == nil {
			placeholderSet = make(map[string]struct{})
		}
		if _, seen := placeholderSet[phName]; !seen {
			placeholderSet[phName] = struct{}{}
			placeholderOrder = append(placeholderOrder, phName)
		}
		partsBuilder = append(partsBuilder, JsMessage_PlaceholderReference_CreateForJsName(phName))
		rest := afterPrefix[phEndRel+len(JsMessage_PH_JS_SUFFIX):]
		if len(rest) == 0 {
			break
		}
		msgText = rest
	}
	return googGetMsgParsedText{
		originalMsgText:  originalMsgText,
		placeholderOrder: placeholderOrder,
		placeholderSet:   placeholderSet,
		parts:            partsBuilder,
	}, nil
}

// visitFallbackFunctionCall visits a call to goog.getMsgWithFallback.
//
// Java: JsMessageVisitor.visitFallbackFunctionCall(NodeTraversal, Node)
func (v *JsMessageVisitor) visitFallbackFunctionCall(t *NodeTraversal, call ast.Call) {
	callNode := call.Node()
	if !callNode.HasXChildren(3) ||
		!jsMessageVisitor_isMessageIdentifier(callNode.GetSecondChild()) ||
		!jsMessageVisitor_isMessageIdentifier(callNode.GetLastChild()) {
		v.compiler.Report(JSError_MakeWithNode(callNode, JsMessageVisitor_BAD_FALLBACK_SYNTAX))
		return
	}

	firstArg := callNode.GetSecondChild()
	firstMessage := v.getJsMessageFromNode(t, firstArg)
	if firstMessage == nil {
		v.compiler.Report(JSError_MakeWithNode(firstArg, JsMessageVisitor_FALLBACK_ARG_ERROR, *firstArg.GetQualifiedName()))
		return
	}

	secondArg := firstArg.GetNext()
	secondMessage := v.getJsMessageFromNode(t, secondArg)
	if secondMessage == nil {
		v.compiler.Report(JSError_MakeWithNode(secondArg, JsMessageVisitor_FALLBACK_ARG_ERROR, *secondArg.GetQualifiedName()))
		return
	}

	if v.onMessageFallback != nil {
		v.onMessageFallback(callNode, firstMessage, secondMessage)
	}
}

// IsMessageName returns whether the given JS identifier is a valid JS
// message name.
//
// Java: JsMessageVisitor.isMessageName(String)
func (v *JsMessageVisitor) IsMessageName(identifier string) bool {
	return strings.HasPrefix(identifier, JsMessageVisitor_MSG_PREFIX) || JsMessageVisitor_IsScopedAliasesPrefix(identifier)
}

// Java: JsMessageVisitor.isMessageIdentifier(Node)
func jsMessageVisitor_isMessageIdentifier(node rhino.Node) bool {
	qname := node.GetQualifiedName()
	return qname != nil && strings.Contains(*qname, JsMessageVisitor_MSG_PREFIX)
}

// Java: JsMessageVisitor.getJsMessageFromNode(NodeTraversal, Node)
func (v *JsMessageVisitor) getJsMessageFromNode(t *NodeTraversal, node rhino.Node) *JsMessage {
	messageNamePtr := node.GetQualifiedName()
	if messageNamePtr == nil {
		return nil
	}
	messageName := *messageNamePtr
	msgIdx := strings.Index(messageName, JsMessageVisitor_MSG_PREFIX)
	if msgIdx < 0 {
		return nil
	}

	messageKey := messageName[msgIdx:]
	if jsMessageVisitor_isUnnamedMessageName(messageKey) {
		return v.getTrackedUnnamedMessage(t, messageName)
	}
	return v.getTrackedNormalMessage(messageKey)
}

// Java: JsMessageVisitor.isUnnamedMessageName(String)
func jsMessageVisitor_isUnnamedMessageName(identifier string) bool {
	return strings.HasPrefix(identifier, "MSG_UNNAMED") && !strings.ContainsAny(identifier, "\r\n")
}

// scopedAliasesPrefixEnd matches `^\$jscomp\$scope\$\S+\$MSG_` and returns the
// byte offset immediately after `$MSG_`, or -1 if it does not match.
func jsMessageVisitor_scopedAliasesPrefixEnd(name string) int {
	if !strings.HasPrefix(name, jsMessageVisitor_SCOPED_ALIASES_PREFIX) {
		return -1
	}
	rest := name[len(jsMessageVisitor_SCOPED_ALIASES_PREFIX):]
	wsIdx := len(rest)
	for i := 0; i < len(rest); i++ {
		switch rest[i] {
		case ' ', '\t', '\n', '\v', '\f', '\r':
			wsIdx = i
		}
		if wsIdx != len(rest) {
			break
		}
	}
	idx := strings.LastIndex(rest[:wsIdx], "$MSG_")
	if idx < 1 {
		return -1
	}
	return len(jsMessageVisitor_SCOPED_ALIASES_PREFIX) + idx + len("$MSG_")
}

// Java: JsMessageVisitor.isScopedAliasesPrefix(String)
func JsMessageVisitor_IsScopedAliasesPrefix(name string) bool {
	return jsMessageVisitor_scopedAliasesPrefixEnd(name) >= 0
}

// Java: JsMessageVisitor.removeScopedAliasesPrefix(String)
func JsMessageVisitor_RemoveScopedAliasesPrefix(name string) string {
	end := jsMessageVisitor_scopedAliasesPrefixEnd(name)
	if end < 0 {
		return name
	}
	return "MSG_" + name[end:]
}

// JsMessageExtractor extracts messages and message comments from JS code.
type JsMessageExtractor struct {
	idGenerator             JsMessage_IdGenerator
	options                 *CompilerOptions
	extractExternalMessages bool
}

// Java: JsMessageExtractor(JsMessage.IdGenerator)
func NewJsMessageExtractor(idGenerator JsMessage_IdGenerator) *JsMessageExtractor {
	return NewJsMessageExtractorWithOptions(idGenerator, NewCompilerOptions(), false)
}

// Java: JsMessageExtractor(JsMessage.IdGenerator, CompilerOptions, boolean)
func NewJsMessageExtractorWithOptions(
	idGenerator JsMessage_IdGenerator, options *CompilerOptions, extractExternalMessages bool,
) *JsMessageExtractor {
	return &JsMessageExtractor{
		idGenerator:             idGenerator,
		options:                 options,
		extractExternalMessages: extractExternalMessages,
	}
}

// ExtractMessages extracts JS messages from JavaScript code.
//
// Java: JsMessageExtractor.extractMessages(SourceFile...)
func (e *JsMessageExtractor) ExtractMessages(inputs ...*SourceFile) []*JsMessage {
	return e.ExtractMessagesFromIterable(inputs)
}

// ExtractMessagesFromIterable extracts JS messages from JavaScript code.
//
// Java: JsMessageExtractor.extractMessages(Iterable<SourceFile>)
func (e *JsMessageExtractor) ExtractMessagesFromIterable(inputs []*SourceFile) []*JsMessage {
	compiler := NewCompiler()
	compiler.Init([]*SourceFile{}, slices.Clone(inputs), e.options)
	Compiler_RunInCompilerThread(compiler,
		func() any {
			compiler.ParseInputs()
			return nil
		})

	var messages []*JsMessage
	recordMessage := func(jsMessage *JsMessage) {
		if e.extractExternalMessages || !jsMessage.IsExternal() {
			messages = append(messages, jsMessage)
		}
	}
	if len(compiler.GetErrors()) == 0 {
		visitor := NewJsMessageVisitor(
			compiler,
			e.idGenerator,
			func(def JsMessageDefinition) { recordMessage(def.GetMessage()) },
			func(def IcuTemplateDefinition) { recordMessage(def.GetMessage()) },
			nil,
		)
		visitor.Process(rhino.NilNode, compiler.GetRoot())
	}

	errs := compiler.GetErrors()
	if len(errs) != 0 {
		var msg strings.Builder
		msg.WriteString("JSCompiler errors\n")
		formatter := NewLightweightMessageFormatter(compiler)
		for _, err := range errs {
			msg.WriteString(formatter.FormatError(err))
		}
		panic(jutil.NewRuntimeException(msg.String()))
	}

	return messages
}
