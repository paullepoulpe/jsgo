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

var jsMessageVisitorVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_CALL)
	s.Add(rhino.Token_NAME)
	s.Add(rhino.Token_ASSIGN)
	s.Add(rhino.Token_STRING_KEY)
	s.Add(rhino.Token_MEMBER_FIELD_DEF)
	return s
}()

func (*JsMessageVisitor) nodeTraversalNoShouldTraverse() {}

func (*JsMessageVisitor) visitTokens() *ast.TokenSet {
	return &jsMessageVisitorVisitTokens
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
