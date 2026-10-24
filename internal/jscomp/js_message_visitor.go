// Port of com/google/javascript/jscomp/JsMessageVisitor.java @ closure-compiler v20261015

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

package jscomp

import (
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const jsMessageVisitor_MSG_FUNCTION_NAME = "getMsg"
const jsMessageVisitor_MSG_FUNCTION_QNAME = "goog." + jsMessageVisitor_MSG_FUNCTION_NAME
const jsMessageVisitor_ICU_MSG_FUNCTION_NAME = "declareIcuTemplate"
const jsMessageVisitor_ICU_MSG_FUNCTION_QNAME = "goog.i18n.messages." + jsMessageVisitor_ICU_MSG_FUNCTION_NAME
const jsMessageVisitor_MSG_FALLBACK_FUNCTION_NAME = "goog.getMsgWithFallback"

// JsMessageVisitor_FUNCTIONS_USED_IN_MESSAGE_EXTRACTION exposes a list of
// functions that are used in message extraction. Standalone tools that run
// message extraction can use this list to filter out files that do not
// contain any of these functions.
var JsMessageVisitor_FUNCTIONS_USED_IN_MESSAGE_EXTRACTION = []string{
	jsMessageVisitor_MSG_FUNCTION_NAME, jsMessageVisitor_ICU_MSG_FUNCTION_NAME}

// jsMessageVisitor_MSG_EXTERNAL_PREFIX identifies a message with a specific
// ID which doesn't get extracted from the JS code containing its
// declaration.
//
//	// A message with ID 1234 whose original definition comes from some source other than this
//	// JS code.
//	const MSG_EXTERNAL_1234 = goog.getMsg("Some message.");
const jsMessageVisitor_MSG_EXTERNAL_PREFIX = "MSG_EXTERNAL_"

var JsMessageVisitor_MESSAGE_HAS_NO_DESCRIPTION = DiagnosticType_Warning(
	"JSC_MSG_HAS_NO_DESCRIPTION", "Message {0} has no description. Add @desc JsDoc tag.")

var JsMessageVisitor_MESSAGE_HAS_NO_TEXT = DiagnosticType_Warning(
	"JSC_MSG_HAS_NO_TEXT",
	"Message value of {0} is just an empty string. "+"Empty messages are forbidden.")

var JsMessageVisitor_MESSAGE_TREE_MALFORMED = DiagnosticType_Error(
	"JSC_MSG_TREE_MALFORMED", "Message parse tree malformed. {0}")

var JsMessageVisitor_MESSAGE_HAS_NO_VALUE = DiagnosticType_Error(
	"JSC_MSG_HAS_NO_VALUE", "message node {0} has no value")

var JsMessageVisitor_MESSAGE_DUPLICATE_KEY = DiagnosticType_Error(
	"JSC_MSG_KEY_DUPLICATED",
	"duplicate message variable name found for {0}, "+"initial definition {1}:{2}")

var JsMessageVisitor_MESSAGE_NODE_IS_ORPHANED = DiagnosticType_Error(
	"JSC_MSG_ORPHANED_NODE",
	"{0}() function may be used only with MSG_* property or variable")

var JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY = DiagnosticType_Warning(
	"JSC_MSG_NOT_INITIALIZED_CORRECTLY",
	"Message must be initialized using a call to "+
		jsMessageVisitor_MSG_FUNCTION_QNAME+
		" or "+
		jsMessageVisitor_ICU_MSG_FUNCTION_QNAME)

// JsMessageVisitor_BAD_FALLBACK_SYNTAX: the Java message is built with
// String.format("Bad syntax. Expected syntax: %s(MSG_1, MSG_2)",
// MSG_FALLBACK_FUNCTION_NAME).
var JsMessageVisitor_BAD_FALLBACK_SYNTAX = DiagnosticType_Error(
	"JSC_MSG_BAD_FALLBACK_SYNTAX",
	"Bad syntax. "+"Expected syntax: "+jsMessageVisitor_MSG_FALLBACK_FUNCTION_NAME+"(MSG_1, MSG_2)")

var JsMessageVisitor_FALLBACK_ARG_ERROR = DiagnosticType_Error(
	"JSC_MSG_FALLBACK_ARG_ERROR", "Could not find message entry for fallback argument {0}")

const JsMessageVisitor_MSG_PREFIX = "MSG_"

// jsMessageVisitor_SCOPED_ALIASES_PREFIX_PATTERN: the ScopedAliases pass
// transforms the goog.scope declarations to have a unique id as prefix.
//
// After '$jscomp$scope$' expect some sequence of digits and dollar signs
// ending in '$MSG_
var jsMessageVisitor_SCOPED_ALIASES_PREFIX_PATTERN = jutil.Pattern_Compile(`\$jscomp\$scope\$\S+\$MSG_`)

// jsMessageVisitor_MSG_UNNAMED_PATTERN is the pattern for unnamed messages.
//
// All JS messages in JS code should have unique name but messages in
// generated code (i.e. from soy template) could have duplicated message
// names. Later we replace the message names with ids constructed as a hash
// of the message content.
//
// Soy (https://github.com/google/closure-templates) generates messages with
// names MSG_UNNAMED.* . This pattern recognizes such messages.
var jsMessageVisitor_MSG_UNNAMED_PATTERN = jutil.Pattern_Compile("MSG_UNNAMED.*")

// JsMessageVisitor locates JS code that is intended to declare localizable
// messages.
//
// It passes each found message to either
// JsMessageVisitor.ProcessJsMessageDefinition(JsMessageDefinition) for
// goog.getMsg() calls or
// JsMessageVisitor.ProcessIcuTemplateDefinition(IcuTemplateDefinition) for
// goog.i18n.messages.declareIcuTemplate() calls.
//
// Java: abstract class JsMessageVisitor extends AbstractPostOrderCallback
// implements CompilerPass
type JsMessageVisitor interface {
	NodeTraversal_Callback
	CompilerPass

	// ProcessJsMessageDefinition processes a found JS message that was
	// defined with goog.getMsg().
	//
	// Several examples of "standard" processing routines are:
	//
	//  1. extract all JS messages
	//  2. replace JS messages with localized versions for some specific
	//     language
	//  3. check that messages have correct syntax and present in
	//     localization bundle
	//
	// Java: JsMessageVisitor.processJsMessageDefinition(JsMessageDefinition)
	ProcessJsMessageDefinition(definition JsMessageDefinition)

	// ProcessIcuTemplateDefinition processes a found call to
	// goog.i18n.messages.declareIcuTemplate()
	//
	// Java: JsMessageVisitor.processIcuTemplateDefinition(IcuTemplateDefinition)
	ProcessIcuTemplateDefinition(definition IcuTemplateDefinition)

	// ProcessMessageFallback processes the goog.getMsgWithFallback
	// primitive. goog.getMsgWithFallback(MSG_1, MSG_2);
	//
	// By default, does nothing.
	//
	// Java: JsMessageVisitor.processMessageFallback(Node, JsMessage, JsMessage)
	ProcessMessageFallback(callNode rhino.Node, message1 *JsMessage, message2 *JsMessage)

	// IsMessageName returns whether the given JS identifier is a valid JS
	// message name.
	//
	// Java: JsMessageVisitor.isMessageName(String)
	IsMessageName(identifier string) bool

	// CheckNode checks a node's type.
	//
	// Panics with a JsMessageVisitor_MalformedException if the node is null
	// or the wrong type.
	//
	// Java: JsMessageVisitor.checkNode(Node, Token)
	CheckNode(node rhino.Node, type_ rhino.Token)

	GetJsMessageVisitorBase() *JsMessageVisitorBase
}

// JsMessageVisitorBase holds the fields and base implementations of
// JsMessageVisitor; every subclass embeds it.
type JsMessageVisitorBase struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	self JsMessageVisitor

	idGenerator JsMessage_IdGenerator
	compiler    AbstractCompiler

	// messageNames are the names encountered associated with their
	// defining node and source. We use it for tracking duplicated message
	// ids in the source code.
	messageNames *jutil.LinkedHashMap[string, *jsMessageVisitor_MessageLocation]

	// unnamedMessages tracks unnamed messages by Var, not string, as they
	// are not guaranteed to be globally unique
	unnamedMessages *jutil.LinkedHashMap[*Var, *JsMessage]

	// googMsgNodes is the set of found goog.getMsg call nodes.
	//
	// When we visit goog.getMsg() node we add it, and later when we visit
	// its parent we remove it. All nodes that are left at the end of
	// traversing are orphaned nodes. It means have no corresponding var or
	// property node.
	googMsgNodes *jutil.LinkedHashSet[rhino.Node]
}

// InitJsMessageVisitorBase creates JS message visitor.
//
// compiler is the compiler instance; idGenerator is the generator that used
// for creating unique ID for the message.
//
// Java: JsMessageVisitor(AbstractCompiler, JsMessage.IdGenerator)
func (b *JsMessageVisitorBase) InitJsMessageVisitorBase(
	self JsMessageVisitor, compiler AbstractCompiler, idGenerator JsMessage_IdGenerator) {
	b.self = self
	b.messageNames = jutil.NewLinkedHashMap[string, *jsMessageVisitor_MessageLocation]()
	b.unnamedMessages = jutil.NewLinkedHashMap[*Var, *JsMessage]()
	b.googMsgNodes = jutil.NewLinkedHashSet[rhino.Node]()

	b.compiler = compiler
	b.idGenerator = idGenerator

	// TODO(anatol): add flag that decides whether to process UNNAMED messages.
	// Some projects would not want such functionality (unnamed) as they don't
	// use SOY templates.
}

// GetJsMessageVisitorBase gives access to the base struct.
func (b *JsMessageVisitorBase) GetJsMessageVisitorBase() *JsMessageVisitorBase { return b }

// Java: JsMessageVisitor.process(Node, Node)
func (b *JsMessageVisitorBase) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(b.compiler, root, b.self)

	for msgNode := range b.googMsgNodes.All() {
		b.compiler.Report(
			JSError_MakeWithNode(
				msgNode, JsMessageVisitor_MESSAGE_NODE_IS_ORPHANED, *msgNode.GetFirstChild().GetQualifiedName()))
	}
}

// Java: JsMessageVisitor.visit(NodeTraversal, Node, Node)
func (b *JsMessageVisitorBase) Visit(traversal *NodeTraversal, node rhino.Node, unused rhino.Node) {
	b.collectGetMsgCall(traversal, node)
	b.checkMessageInitialization(traversal, node)
}

// checkMessageInitialization is called for every Node in the sources AST.
//
// Java: JsMessageVisitor.checkMessageInitialization(NodeTraversal, Node)
func (b *JsMessageVisitorBase) checkMessageInitialization(traversal *NodeTraversal, node rhino.Node) {
	parent := node.GetParent()

	var originalMessageKey *string
	var possiblyObfuscatedMessageKey string
	var msgNode rhino.Node
	var jsDocInfo *rhino.JSDocInfo

	switch node.GetToken() {
	case rhino.Token_NAME:
		// Case: `var MSG_HELLO = 'Message';`
		if parent.IsNil() || !NodeUtil_IsNameDeclaration(parent) {
			return
		}

		possiblyObfuscatedMessageKey = node.GetString()
		originalMessageKey = node.GetOriginalName()
		msgNode = node.GetFirstChild()
		jsDocInfo = parent.GetJSDocInfo()
	case rhino.Token_ASSIGN:
		// Case: `somenamespace.someclass.MSG_HELLO = 'Message';`
		getProp := node.GetFirstChild()
		if !getProp.IsGetProp() {
			return
		}

		possiblyObfuscatedMessageKey = getProp.GetString()
		originalMessageKey = getProp.GetOriginalName()
		msgNode = node.GetLastChild()
		jsDocInfo = node.GetJSDocInfo()
	case rhino.Token_STRING_KEY:
		// Case: `var t = {MSG_HELLO: 'Message'}`;
		if node.IsQuotedStringKey() || !node.HasChildren() || parent.IsObjectPattern() {
			// Don't require goog.getMsg() for quoted keys
			// Case: `var msgs = { 'MSG_QUOTED': anything };`
			//
			// Don't try to require goog.getMsg() for destructuring assignment targets.
			// goog.getMsg() needs to be used in a direct assignment to a variable or property
			// only.
			// Case: `var {MSG_HELLO} = anything;
			// Case: `var {something: MSG_HELLO} = anything;
			return
		}
		jutil.CheckStateMsg1(parent.IsObjectLit(), "%s", parent)

		possiblyObfuscatedMessageKey = node.GetString()
		originalMessageKey = node.GetOriginalName()
		msgNode = node.GetFirstChild()
		jsDocInfo = node.GetJSDocInfo()
	case rhino.Token_MEMBER_FIELD_DEF:
		// Case: `class Foo { MSG_HELLO = 'Message'; }`
		possiblyObfuscatedMessageKey = node.GetString()
		originalMessageKey = node.GetOriginalName()
		msgNode = node.GetFirstChild()
		jsDocInfo = node.GetJSDocInfo()
	default:
		return
	}

	messageKeyFromLhs := possiblyObfuscatedMessageKey
	if originalMessageKey != nil {
		messageKeyFromLhs = *originalMessageKey
	}

	// If we've reached this point, then messageKey is the name of a variable or a property that is
	// being assigned a value and msgNode is the Node representing the value being assigned.
	// However, we haven't actually determined yet that name looks like it should be a translatable
	// message or that the value is a call to goog.getMsg().

	// Is this a message name?
	msgNodeIsACall := !msgNode.IsNil() && msgNode.IsCall()

	if !b.self.IsMessageName(messageKeyFromLhs) {
		return
	}

	if msgNode.IsNil() {
		b.compiler.Report(JSError_MakeWithNode(node, JsMessageVisitor_MESSAGE_HAS_NO_VALUE, messageKeyFromLhs))
		return
	}

	if jsMessageVisitor_isLegalMessageVarAlias(msgNode) {
		return
	}

	// Report a warning if a qualified messageKey that looks like a message (e.g. "a.b.MSG_X")
	// doesn't use goog.getMsg().
	if msgNodeIsACall {
		b.googMsgNodes.Remove(msgNode)
	} else {
		b.compiler.Report(JSError_MakeWithNode(node, JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY))
		return
	}

	mapping := b.compiler.GetSourceMapping(traversal.GetSourceName(), node.GetLineno(), node.GetCharno())
	var sourceName string
	if mapping != nil {
		sourceName = mapping.GetOriginalFile() + ":" + strconv.Itoa(mapping.GetLineNumber())
	} else {
		sourceName = jutil.ToJavaString(traversal.GetSourceName()) + ":" + strconv.Itoa(node.GetLineno())
	}

	jutil.CheckState(msgNode.IsCall())
	fnNameNode := msgNode.GetFirstChild()
	jutil.Try(func() {
		if b.isDeclareIcuTemplateCallee(fnNameNode) {
			icuTemplateDefinition :=
				b.extractIcuTemplateDefinition(msgNode, jsDocInfo, messageKeyFromLhs, sourceName)
			extractedMessage := icuTemplateDefinition.GetMessage()
			b.trackMessage(traversal, possiblyObfuscatedMessageKey, msgNode, extractedMessage)
			b.reportErrorIfEmptyMessage(node, extractedMessage)
			b.self.ProcessIcuTemplateDefinition(icuTemplateDefinition)
		} else if fnNameNode.MatchesQualifiedName(jsMessageVisitor_MSG_FUNCTION_QNAME) {
			jsMessageDefinition :=
				b.extractJsMessageDefinition(msgNode, jsDocInfo, messageKeyFromLhs, sourceName)
			extractedMessage := jsMessageDefinition.GetMessage()
			b.trackMessage(traversal, possiblyObfuscatedMessageKey, msgNode, extractedMessage)
			b.reportErrorIfEmptyMessage(node, extractedMessage)

			// goog.getMsg() calls are required to have `@desc` unless they are external.
			desc := extractedMessage.GetDesc()
			if (desc == nil || jutil.StringTrim(*desc) == "") && !extractedMessage.IsExternal() {
				b.compiler.Report(
					JSError_MakeWithNode(node, JsMessageVisitor_MESSAGE_HAS_NO_DESCRIPTION, extractedMessage.GetKey()))
			}

			b.self.ProcessJsMessageDefinition(jsMessageDefinition)
		} else {
			b.compiler.Report(
				JSError_MakeWithNode(
					msgNode,
					JsMessageVisitor_MESSAGE_TREE_MALFORMED,
					"Message must be initialized using a call to "+
						jsMessageVisitor_MSG_FUNCTION_QNAME+
						" or "+
						jsMessageVisitor_ICU_MSG_FUNCTION_NAME+
						" (from goog.i18n.messages)."))
		}
	}, func(ex IJsMessageVisitor_MalformedException) {
		b.compiler.Report(JSError_MakeWithNode(ex.GetNode(), JsMessageVisitor_MESSAGE_TREE_MALFORMED, ex.GetMessage()))
	})
}

// Java: JsMessageVisitor.isDeclareIcuTemplateCallee(Node)
func (b *JsMessageVisitorBase) isDeclareIcuTemplateCallee(calleeNode rhino.Node) bool {
	// NOTE: We're requiring that the imported ICU template function name match its original name.
	// This is an intentional limitation, since it should be less confusing to users as well as
	// to this program.
	// TODO(bradfordcsmith): Make this check more robust, possibly building on ModuleMetadata or
	// ProcessClosurePrimitives.
	if calleeNode.IsName() {
		return jutil.StringEndsWith(calleeNode.GetString(), jsMessageVisitor_ICU_MSG_FUNCTION_NAME)
	}
	if calleeNode.IsGetProp() {
		return jutil.StringEndsWith(calleeNode.GetString(), jsMessageVisitor_ICU_MSG_FUNCTION_NAME) && calleeNode.IsQualifiedName()
	}
	return false
}

// JsMessageVisitor_InputForBuildJsMsg is the input of buildJsMessage.
//
// Java: interface JsMessageVisitor.InputForBuildJsMsg
type JsMessageVisitor_InputForBuildJsMsg interface {
	// Java: JsMessageVisitor.InputForBuildJsMsg.getMessageKeyFromLhs()
	GetMessageKeyFromLhs() string

	// Java: JsMessageVisitor.InputForBuildJsMsg.getSourceName()
	GetSourceName() string

	// Java: JsMessageVisitor.InputForBuildJsMsg.getMessageText()
	GetMessageText() string

	// Java: JsMessageVisitor.InputForBuildJsMsg.getMessageParts()
	GetMessageParts() []JsMessage_Part

	// Java: JsMessageVisitor.InputForBuildJsMsg.getPlaceholderExampleMap()
	GetPlaceholderExampleMap() *jutil.LinkedHashMap[string, string]

	// Java: JsMessageVisitor.InputForBuildJsMsg.getPlaceholderOriginalCodeMap()
	GetPlaceholderOriginalCodeMap() *jutil.LinkedHashMap[string, string]

	// Java: JsMessageVisitor.InputForBuildJsMsg.getDescription()
	GetDescription() *string

	// Java: JsMessageVisitor.InputForBuildJsMsg.getMeaning()
	GetMeaning() *string

	// Java: JsMessageVisitor.InputForBuildJsMsg.getAlternateMessageId()
	GetAlternateMessageId() *string
}

// Java: JsMessageVisitor.buildJsMessage(InputForBuildJsMsg)
func (b *JsMessageVisitorBase) buildJsMessage(input JsMessageVisitor_InputForBuildJsMsg) *JsMessage {
	messageParts := input.GetMessageParts()
	messageKeyFromLhs := input.GetMessageKeyFromLhs()

	// non-null for `MSG_EXTERNAL_12345`
	externalMessageId := JsMessageVisitor_GetExternalMessageId(messageKeyFromLhs)

	var isAnonymous bool
	var isExternal bool
	var messageId string
	var messageKeyFinal string
	if externalMessageId != nil {
		// MSG_EXTERNAL_12345 = ...
		isAnonymous = false
		isExternal = true
		messageId = *externalMessageId
		messageKeyFinal = messageKeyFromLhs
	} else {
		isExternal = false
		// We will need to generate the message ID from a combination of the message text and
		// its "meaning". If the code explicitly specifies a "meaning" string we'll use that,
		// otherwise we'll use the message key as the meaning
		meaningForIdGeneration := input.GetMeaning()
		if jsMessageVisitor_isUnnamedMessageName(messageKeyFromLhs) {
			// MSG_UNNAMED_XXXX = goog.getMsg(....);
			// JS code that is automatically generated uses this, since it is harder for it to create
			// message variable names that are guaranteed to be unique.
			isAnonymous = true
			messageKeyFinal = b.generateKeyFromMessageText(input.GetMessageText())
			if meaningForIdGeneration == nil {
				meaningForIdGeneration = jutil.StringPtr(messageKeyFinal)
			}
		} else {
			isAnonymous = false
			messageKeyFinal = messageKeyFromLhs
			if meaningForIdGeneration == nil {
				// Transpilation of goog.scope() may have added a prefix onto the variable name,
				// which we need to strip off when we treat it as a "meaning" for ID generation purposes.
				// Otherwise, the message ID would change if a goog.scope() call were added or removed.
				meaningForIdGeneration = jutil.StringPtr(JsMessageVisitor_RemoveScopedAliasesPrefix(messageKeyFromLhs))
			}
		}
		if b.idGenerator == nil {
			messageId = *meaningForIdGeneration
		} else {
			messageId = b.idGenerator.GenerateId(meaningForIdGeneration, messageParts)
		}
	}

	sourceName := input.GetSourceName()
	return NewJsMessage_Builder().
		AppendParts(messageParts).
		SetPlaceholderNameToExampleMap(input.GetPlaceholderExampleMap()).
		SetPlaceholderNameToOriginalCodeMap(input.GetPlaceholderOriginalCodeMap()).
		SetIsAnonymous(isAnonymous).
		SetIsExternalMsg(isExternal).
		SetKey(messageKeyFinal).
		SetSourceName(&sourceName).
		SetDesc(input.GetDescription()).
		// NOTE: JsMessageXmbWriter does NOT just take this value and write it to the XMB
		// file as the message meaning. It has its own code that defaults a null value to
		// the value of the key field, then prefixes that with a project ID, if any.
		// The intention of that seems to have been to make sure the meaning value in the XMB
		// file matches what was actually used by the ID generator above.
		SetMeaning(input.GetMeaning()).
		SetAlternateId(input.GetAlternateMessageId()).
		SetId(messageId).
		Build()
}

// Java: JsMessageVisitor.trackMessage(NodeTraversal, String, Node, JsMessage)
func (b *JsMessageVisitorBase) trackMessage(
	traversal *NodeTraversal,
	possiblyObfuscatedMessageKey string,
	msgNode rhino.Node,
	extractedMessage *JsMessage) {
	// If asked to check named internal messages.
	if !extractedMessage.IsAnonymous() && !extractedMessage.IsExternal() {
		b.checkIfMessageDuplicated(extractedMessage.GetKey(), msgNode)
	}
	if extractedMessage.IsAnonymous() {
		b.trackUnnamedMessage(traversal, extractedMessage, possiblyObfuscatedMessageKey)
	} else {
		b.trackNormalMessage(extractedMessage, extractedMessage.GetKey(), msgNode)
	}
}

// Java: JsMessageVisitor.reportErrorIfEmptyMessage(Node, JsMessage)
func (b *JsMessageVisitorBase) reportErrorIfEmptyMessage(node rhino.Node, extractedMessage *JsMessage) {
	if extractedMessage.IsEmpty() {
		// value of the message is an empty string. Translators do not like it.
		b.compiler.Report(JSError_MakeWithNode(node, JsMessageVisitor_MESSAGE_HAS_NO_TEXT, extractedMessage.GetKey()))
	}
}

// JsMessageVisitor_GetExternalMessageId extracts an external message ID
// from the message key, if it contains one.
//
// messageKey is the message key (usually the variable or property name).
// Returns the external ID if it is found, otherwise `null`
//
// Java: JsMessageVisitor.getExternalMessageId(String)
func JsMessageVisitor_GetExternalMessageId(messageKey string) *string {
	if jutil.StringStartsWith(messageKey, jsMessageVisitor_MSG_EXTERNAL_PREFIX) {
		start := jutil.StringLength(jsMessageVisitor_MSG_EXTERNAL_PREFIX)
		end := start
		for ; end < jutil.StringLength(messageKey); end++ {
			c := jutil.StringCharAt(messageKey, end)
			if c > '9' || c < '0' {
				break
			}
		}
		if end > start {
			return jutil.StringPtr(jutil.StringSubstring(messageKey, start, end))
		}
	}
	return nil
}

// Java: JsMessageVisitor.generateKeyFromMessageText(String)
func (b *JsMessageVisitorBase) generateKeyFromMessageText(msgText string) string {
	nonnegativeHash := int64(0x7fffffffffffffff) & JsMessage_Hash_Hash64(&msgText)
	return JsMessageVisitor_MSG_PREFIX + jutil.StringToUpperCaseASCII(strconv.FormatInt(nonnegativeHash, 36))
}

// Java: JsMessageVisitor.collectGetMsgCall(NodeTraversal, Node)
func (b *JsMessageVisitorBase) collectGetMsgCall(traversal *NodeTraversal, call rhino.Node) {
	if !call.IsCall() {
		return
	}

	// goog.getMsg()
	callee := call.GetFirstChild()
	if callee.MatchesQualifiedName(jsMessageVisitor_MSG_FUNCTION_QNAME) || b.isDeclareIcuTemplateCallee(callee) {
		b.googMsgNodes.Add(call)
	} else if callee.MatchesQualifiedName(jsMessageVisitor_MSG_FALLBACK_FUNCTION_NAME) {
		b.visitFallbackFunctionCall(traversal, call)
	}
}

// trackNormalMessage tracks a message for later retrieval.
//
// This is used for tracking duplicates, and for figuring out message
// fallback. Not all message types are trackable, because that would require
// a more sophisticated analysis. e.g., function f(s) { s.MSG_UNNAMED_X =
// 'Some untrackable message'; }
//
// Java: JsMessageVisitor.trackNormalMessage(JsMessage, String, Node)
func (b *JsMessageVisitorBase) trackNormalMessage(message *JsMessage, msgName string, msgNode rhino.Node) {
	location := newJsMessageVisitor_MessageLocation(message, msgNode)
	b.messageNames.Put(msgName, location)
}

// trackUnnamedMessage tracks an unnamed message for later retrieval.
//
// This is used for figuring out message fallbacks. Message duplicates are
// allowed for unnamed messages.
//
// Java: JsMessageVisitor.trackUnnamedMessage(NodeTraversal, JsMessage, String)
func (b *JsMessageVisitorBase) trackUnnamedMessage(t *NodeTraversal, message *JsMessage, msgNameInScope string) {
	var_ := t.GetScope().GetVar(msgNameInScope)
	if var_ != nil {
		b.unnamedMessages.Put(var_, message)
	}
}

// jsMessageVisitor_isLegalMessageVarAlias defines any special cases that
// are exceptions to what would otherwise be illegal message assignments.
//
// These exceptions are generally due to the pass being designed before new
// syntax was introduced.
//
// msgNode is the Node representing the value assigned to the message
// variable or property
//
// Java: JsMessageVisitor.isLegalMessageVarAlias(Node)
func jsMessageVisitor_isLegalMessageVarAlias(msgNode rhino.Node) bool {
	if msgNode.IsGetProp() &&
		msgNode.IsQualifiedName() &&
		jutil.StringStartsWith(msgNode.GetString(), JsMessageVisitor_MSG_PREFIX) {
		// Case: `foo.Thing.MSG_EXAMPLE_ALIAS = bar.OtherThing.MSG_EXAMPLE;`
		//
		// This kind of construct is created by TypeScript code generation and
		// ConcretizeStaticInheritanceForInlining. Just ignore it; the message will have already been
		// extracted
		// from the base class.
		return true
	}

	if !msgNode.IsName() {
		return false
	}

	var originalName string
	if msgNode.GetOriginalName() != nil {
		originalName = *msgNode.GetOriginalName()
	} else {
		originalName = msgNode.GetString()
	}

	if jutil.StringStartsWith(originalName, JsMessageVisitor_MSG_PREFIX) {
		// Creating an alias for a message is also allowed, and sometimes happens in generated code,
		// including some of the code generated by this compiler's transpilations.
		// e.g.
		// `var MSG_EXAMPLE_ALIAS = MSG_EXAMPLE;`
		// `var {MSG_HELLO_ALIAS} = MSG_HELLO;
		// `var {MSG_HELLO_ALIAS: goog$module$my$module_MSG_HELLO} = x;`).
		// `exports = {MSG_FOO}`
		// or `exports = {MSG_FOO: MSG_FOO}` when used with declareLegacyNamespace.
		return true
	}

	return false
}

// getTrackedUnnamedMessage gets a previously tracked unnamed message
//
// Java: JsMessageVisitor.getTrackedUnnamedMessage(NodeTraversal, String)
func (b *JsMessageVisitorBase) getTrackedUnnamedMessage(t *NodeTraversal, msgNameInScope string) *JsMessage {
	var_ := t.GetScope().GetVar(msgNameInScope)
	if var_ != nil {
		return b.unnamedMessages.Get(var_)
	}
	return nil
}

// getTrackedNormalMessage gets a previously tracked message.
//
// Java: JsMessageVisitor.getTrackedNormalMessage(String)
func (b *JsMessageVisitorBase) getTrackedNormalMessage(msgName string) *JsMessage {
	location := b.messageNames.Get(msgName)
	if location == nil {
		return nil
	}
	return location.message
}

// checkIfMessageDuplicated checks if message already processed. If so - it
// generates 'message duplicated' compiler error.
//
// msgName is the name of the message; msgNode is the node that represents
// JS message
//
// Java: JsMessageVisitor.checkIfMessageDuplicated(String, Node)
func (b *JsMessageVisitorBase) checkIfMessageDuplicated(msgName string, msgNode rhino.Node) {
	if b.messageNames.ContainsKey(msgName) {
		location := b.messageNames.Get(msgName)
		b.compiler.Report(
			JSError_MakeWithNode(
				msgNode,
				JsMessageVisitor_MESSAGE_DUPLICATE_KEY,
				msgName,
				jutil.ToJavaString(location.messageNode.GetSourceFileName()),
				strconv.Itoa(location.messageNode.GetLineno())))
	}
}

// jsMessageVisitor_extractStringFromStringExprNode returns the string value
// associated with a node representing a JS string or several JS strings
// added together (e.g. 'str' or 's' + 't' + 'r').
//
// node is the node from where we extract the string. Returns the String
// representation of the node. Panics with a
// JsMessageVisitor_MalformedException if the node is not a string literal
// or concatenation of them
//
// Java: JsMessageVisitor.extractStringFromStringExprNode(Node)
func jsMessageVisitor_extractStringFromStringExprNode(node rhino.Node) string {
	switch node.GetToken() {
	case rhino.Token_STRINGLIT:
		return node.GetString()
	case rhino.Token_TEMPLATELIT:
		if node.HasOneChild() {
			// Cooked string can be null only for tagged template literals.
			// A tagged template literal would hit the default case below.
			return *jutil.CheckNotNull(node.GetFirstChild().GetCookedString())
		} else {
			panic(NewJsMessageVisitor_MalformedException(
				"Template literals with substitutions are not allowed.", node))
		}
	case rhino.Token_ADD:
		sb := jutil.NewStringBuilder()
		for child := node.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			sb.Append(jsMessageVisitor_extractStringFromStringExprNode(child))
		}
		return sb.String()
	default:
		panic(NewJsMessageVisitor_MalformedException("literal string or concatenation expected", node))
	}
}

// extractJsMessageDefinition extracts data from a call to goog.getMsg(...).
//
// Here's an example with all possible arguments and options present.
//
//	const MSG_HELLO =
//	    goog.getMsg(
//	        'Hello, {$name}', // message template string
//	        { name: getName() }, // placeholder replacement values
//	        { // options bag
//	          html: false,
//	          unescapeHtmlEntities: true,
//	          example: { name: 'George' },
//	          original_code: { name: 'getName()' },
//	        });
//
// Panics with a JsMessageVisitor_MalformedException if the code does not
// match the expected format.
//
// Java: JsMessageVisitor.extractJsMessageDefinition(Node, JSDocInfo, String, String)
func (b *JsMessageVisitorBase) extractJsMessageDefinition(
	msgNode rhino.Node, jsDocInfo *rhino.JSDocInfo, messageKeyFromLhs string, sourceName string) JsMessageDefinition {
	// Extract data from a call to `goog.getMsg(...)`
	// TODO(bradfordcsmith): Add these three fields to the options bag argument for goog.getMsg().
	// Specifying them there instead of in annotations will make them available to the uncompiled
	// `goog.getMsg()` call during runtime, which could enable runtime lookup of translated messages
	// in uncompiled code.
	var description *string
	var meaning *string
	var alternateMessageId *string
	if jsDocInfo == nil {
		description = nil
		meaning = nil
		alternateMessageId = nil
	} else {
		description = jsDocInfo.GetDescription()
		meaning = jsDocInfo.GetMeaning()
		alternateMessageId = jsDocInfo.GetAlternateMessageId()
	}

	// first child of the call is the `goog.getMsg` name
	// second is the message string value
	msgTextNode := msgNode.GetSecondChild()
	if msgTextNode.IsNil() {
		panic(NewJsMessageVisitor_MalformedException("Message string literal expected", msgNode))
	}
	// third is the optional object literal mapping placeholder names to value expressions
	valuesObjLit := msgTextNode.GetNext()
	// fourth is the optional options-bag argument
	var optionsBagArgument rhino.Node
	if !valuesObjLit.IsNil() {
		optionsBagArgument = valuesObjLit.GetNext()
	}

	var unexpectedArgument rhino.Node
	if !optionsBagArgument.IsNil() {
		unexpectedArgument = optionsBagArgument.GetNext()
	}
	if !unexpectedArgument.IsNil() {
		panic(NewJsMessageVisitor_MalformedException("too many arguments", unexpectedArgument))
	}

	// Extract and parse the message text
	// The message string can be a string literal or a concatenation of string literals.
	// We want the whole thing as a single string here.
	msgTextString := jsMessageVisitor_extractStringFromStringExprNode(msgTextNode)
	var googGetMsgParsedText jsMessageVisitor_GoogGetMsgParsedText
	jutil.Try(func() {
		googGetMsgParsedText = jsMessageVisitor_extractGoogGetMsgParsedText(msgTextString)
	}, func(e IJsMessage_PlaceholderFormatException) {
		panic(NewJsMessageVisitor_MalformedException(e.GetMessage(), msgTextNode))
	})

	// Extract the placeholder values object
	// NOTE: The extract method returns an effectively "empty" value for `null`
	placeholderObjectLiteralMap := JsMessageVisitor_ExtractObjectLiteralMap(valuesObjLit)

	// Confirm that we can find a value map entry for every placeholder name referenced in the text.
	placeholderNamesFromText := googGetMsgParsedText.GetPlaceholderNames()
	placeholderObjectLiteralMap.CheckForRequiredKeys(
		placeholderNamesFromText,
		func(placeholderName string) *JsMessageVisitor_MalformedException {
			return NewJsMessageVisitor_MalformedException(
				"Unrecognized message placeholder referenced: "+placeholderName, msgTextNode)
		})

	// Confirm that every placeholder name for which we have a value is actually referenced in the
	// text.
	placeholderObjectLiteralMap.CheckForUnexpectedKeys(
		placeholderNamesFromText,
		func(placeholderName string) string { return "Unused message placeholder: " + placeholderName })

	// Get a map from placeholder name to Node value
	placeholderValuesMap := placeholderObjectLiteralMap.ExtractAsValueMap()

	// Extract the Options bag data
	// NOTE: the extract method returns an "empty" value for `null`
	jsMessageOptions := b.extractJsMessageOptions(optionsBagArgument)

	// Confirm that any example or orignal_code placeholder information we have refers only to
	// placeholder names that appear in the message text.
	jsMessageOptions.CheckForUnknownPlaceholders(placeholderNamesFromText)

	jsExtractedMessage :=
		b.buildJsMessage(
			&jsMessageVisitor_googGetMsgInput{
				messageKeyFromLhs:    messageKeyFromLhs,
				sourceName:           sourceName,
				googGetMsgParsedText: googGetMsgParsedText,
				jsMessageOptions:     jsMessageOptions,
				description:          description,
				meaning:              meaning,
				alternateMessageId:   alternateMessageId,
			})

	return &jsMessageVisitor_jsMessageDefinition{
		jsExtractedMessage:   jsExtractedMessage,
		msgNode:              msgNode,
		msgTextNode:          msgTextNode,
		valuesObjLit:         valuesObjLit,
		placeholderValuesMap: placeholderValuesMap,
		jsMessageOptions:     jsMessageOptions,
	}
}

// jsMessageVisitor_googGetMsgInput is the anonymous InputForBuildJsMsg of
// extractJsMessageDefinition.
type jsMessageVisitor_googGetMsgInput struct {
	messageKeyFromLhs    string
	sourceName           string
	googGetMsgParsedText jsMessageVisitor_GoogGetMsgParsedText
	jsMessageOptions     jsMessageVisitor_JsMessageOptions
	description          *string
	meaning              *string
	alternateMessageId   *string
}

// Java: JsMessageVisitor.extractJsMessageDefinition(...).new InputForBuildJsMsg() {...}.getMessageKeyFromLhs()
func (i *jsMessageVisitor_googGetMsgInput) GetMessageKeyFromLhs() string { return i.messageKeyFromLhs }

// Java: JsMessageVisitor.extractJsMessageDefinition(...).new InputForBuildJsMsg() {...}.getSourceName()
func (i *jsMessageVisitor_googGetMsgInput) GetSourceName() string { return i.sourceName }

// Java: JsMessageVisitor.extractJsMessageDefinition(...).new InputForBuildJsMsg() {...}.getMessageText()
func (i *jsMessageVisitor_googGetMsgInput) GetMessageText() string {
	return i.googGetMsgParsedText.GetText()
}

// Java: JsMessageVisitor.extractJsMessageDefinition(...).new InputForBuildJsMsg() {...}.getMessageParts()
func (i *jsMessageVisitor_googGetMsgInput) GetMessageParts() []JsMessage_Part {
	return i.googGetMsgParsedText.GetParts()
}

// Java: JsMessageVisitor.extractJsMessageDefinition(...).new InputForBuildJsMsg() {...}.getPlaceholderExampleMap()
func (i *jsMessageVisitor_googGetMsgInput) GetPlaceholderExampleMap() *jutil.LinkedHashMap[string, string] {
	return i.jsMessageOptions.GetPlaceholderExampleMap()
}

// Java: JsMessageVisitor.extractJsMessageDefinition(...).new InputForBuildJsMsg() {...}.getPlaceholderOriginalCodeMap()
func (i *jsMessageVisitor_googGetMsgInput) GetPlaceholderOriginalCodeMap() *jutil.LinkedHashMap[string, string] {
	return i.jsMessageOptions.GetPlaceholderOriginalCodeMap()
}

// Java: JsMessageVisitor.extractJsMessageDefinition(...).new InputForBuildJsMsg() {...}.getDescription()
func (i *jsMessageVisitor_googGetMsgInput) GetDescription() *string { return i.description }

// Java: JsMessageVisitor.extractJsMessageDefinition(...).new InputForBuildJsMsg() {...}.getMeaning()
func (i *jsMessageVisitor_googGetMsgInput) GetMeaning() *string { return i.meaning }

// Java: JsMessageVisitor.extractJsMessageDefinition(...).new InputForBuildJsMsg() {...}.getAlternateMessageId()
func (i *jsMessageVisitor_googGetMsgInput) GetAlternateMessageId() *string {
	return i.alternateMessageId
}

// jsMessageVisitor_jsMessageDefinition is the anonymous JsMessageDefinition
// of extractJsMessageDefinition.
type jsMessageVisitor_jsMessageDefinition struct {
	jsExtractedMessage   *JsMessage
	msgNode              rhino.Node
	msgTextNode          rhino.Node
	valuesObjLit         rhino.Node
	placeholderValuesMap *jutil.LinkedHashMap[string, rhino.Node]
	jsMessageOptions     jsMessageVisitor_JsMessageOptions
}

// Java: JsMessageVisitor.extractJsMessageDefinition(...).new JsMessageDefinition() {...}.getMessage()
func (d *jsMessageVisitor_jsMessageDefinition) GetMessage() *JsMessage { return d.jsExtractedMessage }

// Java: JsMessageVisitor.extractJsMessageDefinition(...).new JsMessageDefinition() {...}.getMessageNode()
func (d *jsMessageVisitor_jsMessageDefinition) GetMessageNode() rhino.Node { return d.msgNode }

// Java: JsMessageVisitor.extractJsMessageDefinition(...).new JsMessageDefinition() {...}.getTemplateTextNode()
func (d *jsMessageVisitor_jsMessageDefinition) GetTemplateTextNode() rhino.Node {
	return d.msgTextNode
}

// Java: JsMessageVisitor.extractJsMessageDefinition(...).new JsMessageDefinition() {...}.getPlaceholderValuesNode()
func (d *jsMessageVisitor_jsMessageDefinition) GetPlaceholderValuesNode() rhino.Node {
	return d.valuesObjLit
}

// Java: JsMessageVisitor.extractJsMessageDefinition(...).new JsMessageDefinition() {...}.getPlaceholderValueMap()
func (d *jsMessageVisitor_jsMessageDefinition) GetPlaceholderValueMap() *jutil.LinkedHashMap[string, rhino.Node] {
	return d.placeholderValuesMap
}

// Java: JsMessageVisitor.extractJsMessageDefinition(...).new JsMessageDefinition() {...}.shouldEscapeLessThan()
func (d *jsMessageVisitor_jsMessageDefinition) ShouldEscapeLessThan() bool {
	return d.jsMessageOptions.IsEscapeLessThan()
}

// Java: JsMessageVisitor.extractJsMessageDefinition(...).new JsMessageDefinition() {...}.shouldUnescapeHtmlEntities()
func (d *jsMessageVisitor_jsMessageDefinition) ShouldUnescapeHtmlEntities() bool {
	return d.jsMessageOptions.IsUnescapeHtmlEntities()
}

// extractIcuTemplateDefinition extracts message data from a
// `declareIcuTemplate()` call.
//
// Example:
//
//	const MSG_ICU_EXAMPLE =
//	    declareIcuTemplate(
//	        `{NUM_PEOPLE, plural, offset:1
//	        =0 {I see {BEGIN_BOLD}no one at all{END_BOLD} in {INTERPOLATION_2}.}
//	        =1 {I see {INTERPOLATION_1} in {INTERPOLATION_2}.}
//	        =2 {I see {INTERPOLATION_1} and one other person in {INTERPOLATION_2}.}
//	        other {I see {INTERPOLATION_1} and some other people in {INTERPOLATION_2}.}
//	    }`,
//	    {
//	      description: 'required message description here',
//	      meaning: 'optional meaning here',
//	      alternate_message_id: 'optional alternate message ID here',
//	      // optional
//	      original_code: {
//	        'INTERPOLATION_1': '{{getPerson()}}',
//	        'INTERPOLATION_2': '{{getPlaceName()}}',
//	      },
//	      // optional
//	      example: {
//	        'INTERPOLATION_1': 'Jane Doe',
//	        'INTERPOLATION_2': 'Paris, France',
//	      }
//	    });
//
// Java: JsMessageVisitor.extractIcuTemplateDefinition(Node, JSDocInfo, String, String)
func (b *JsMessageVisitorBase) extractIcuTemplateDefinition(
	msgNode rhino.Node, jsDocInfo *rhino.JSDocInfo, messageKeyFromLhs string, sourceName string) IcuTemplateDefinition {
	if jsDocInfo != nil {
		// For declareIcuTemplateData() it is not valid to use the @desc, @meaning, and
		// @alternateMessageId annotations.
		// Instead, that information should go into the options bag parameter.
		if jsDocInfo.GetAlternateMessageId() != nil {
			panic(NewJsMessageVisitor_MalformedException(
				"Use the 'alternate_message_id' option, not the '@alternateMessageId' annotation",
				msgNode))
		}
		if jsDocInfo.GetDescription() != nil {
			panic(NewJsMessageVisitor_MalformedException(
				"Use the 'description' option, not the '@desc' annotation", msgNode))
		}
		if jsDocInfo.GetMeaning() != nil {
			panic(NewJsMessageVisitor_MalformedException(
				"Use the 'meaning' option, not the '@meaning' annotation", msgNode))
		}
	}

	// The first argument is the message string
	stringLiteralExpression := msgNode.GetSecondChild()
	if stringLiteralExpression.IsNil() {
		panic(NewJsMessageVisitor_MalformedException("message string argument expected", msgNode))
	}
	icuMessageTemplateString := b.extractIcuMessageTemplateString(stringLiteralExpression)

	// The second argument is the options bag
	optionsBagNode := stringLiteralExpression.GetNext()
	if optionsBagNode.IsNil() {
		panic(NewJsMessageVisitor_MalformedException("options argument expected", msgNode))
	}
	icuTemplateOptions := b.extractIcuTemplateOptions(optionsBagNode)

	// Parsing of the template string into parts depends on which placeholders are mentioned
	// in the options. Placeholders that do not have example or original_code text  will not be
	// extracted from the template as separate parts.
	extractedIcuTemplateParts := icuMessageTemplateString.ExtractParts(icuTemplateOptions.GetPlaceholderNames())
	icuTemplateOptions.CheckForUnknownPlaceholders(
		extractedIcuTemplateParts.extractedPlaceholderNames)

	// There shouldn't be another argument
	unexpectedArgument := optionsBagNode.GetNext()
	if !unexpectedArgument.IsNil() {
		panic(NewJsMessageVisitor_MalformedException("too many arguments", unexpectedArgument))
	}

	extractedMessage :=
		b.buildJsMessage(
			&jsMessageVisitor_icuTemplateInput{
				messageKeyFromLhs:         messageKeyFromLhs,
				sourceName:                sourceName,
				icuMessageTemplateString:  icuMessageTemplateString,
				extractedIcuTemplateParts: extractedIcuTemplateParts,
				icuTemplateOptions:        icuTemplateOptions,
			})

	return &jsMessageVisitor_icuTemplateDefinition{
		extractedMessage:        extractedMessage,
		msgNode:                 msgNode,
		stringLiteralExpression: stringLiteralExpression,
	}
}

// jsMessageVisitor_icuTemplateInput is the anonymous InputForBuildJsMsg of
// extractIcuTemplateDefinition.
type jsMessageVisitor_icuTemplateInput struct {
	messageKeyFromLhs         string
	sourceName                string
	icuMessageTemplateString  *JsMessageVisitor_IcuMessageTemplateString
	extractedIcuTemplateParts *JsMessageVisitor_ExtractedIcuTemplateParts
	icuTemplateOptions        jsMessageVisitor_IcuTemplateOptions
}

// Java: JsMessageVisitor.extractIcuTemplateDefinition(...).new InputForBuildJsMsg() {...}.getMessageKeyFromLhs()
func (i *jsMessageVisitor_icuTemplateInput) GetMessageKeyFromLhs() string { return i.messageKeyFromLhs }

// Java: JsMessageVisitor.extractIcuTemplateDefinition(...).new InputForBuildJsMsg() {...}.getSourceName()
func (i *jsMessageVisitor_icuTemplateInput) GetSourceName() string { return i.sourceName }

// Java: JsMessageVisitor.extractIcuTemplateDefinition(...).new InputForBuildJsMsg() {...}.getMessageText()
func (i *jsMessageVisitor_icuTemplateInput) GetMessageText() string {
	return i.icuMessageTemplateString.template
}

// Java: JsMessageVisitor.extractIcuTemplateDefinition(...).new InputForBuildJsMsg() {...}.getMessageParts()
func (i *jsMessageVisitor_icuTemplateInput) GetMessageParts() []JsMessage_Part {
	return i.extractedIcuTemplateParts.extractedParts
}

// Java: JsMessageVisitor.extractIcuTemplateDefinition(...).new InputForBuildJsMsg() {...}.getPlaceholderExampleMap()
func (i *jsMessageVisitor_icuTemplateInput) GetPlaceholderExampleMap() *jutil.LinkedHashMap[string, string] {
	return i.icuTemplateOptions.GetPlaceholderExampleMap()
}

// Java: JsMessageVisitor.extractIcuTemplateDefinition(...).new InputForBuildJsMsg() {...}.getPlaceholderOriginalCodeMap()
func (i *jsMessageVisitor_icuTemplateInput) GetPlaceholderOriginalCodeMap() *jutil.LinkedHashMap[string, string] {
	return i.icuTemplateOptions.GetPlaceholderOriginalCodeMap()
}

// Java: JsMessageVisitor.extractIcuTemplateDefinition(...).new InputForBuildJsMsg() {...}.getDescription()
func (i *jsMessageVisitor_icuTemplateInput) GetDescription() *string {
	return jutil.StringPtr(i.icuTemplateOptions.GetDescription())
}

// Java: JsMessageVisitor.extractIcuTemplateDefinition(...).new InputForBuildJsMsg() {...}.getMeaning()
func (i *jsMessageVisitor_icuTemplateInput) GetMeaning() *string {
	return i.icuTemplateOptions.GetMeaning()
}

// Java: JsMessageVisitor.extractIcuTemplateDefinition(...).new InputForBuildJsMsg() {...}.getAlternateMessageId()
func (i *jsMessageVisitor_icuTemplateInput) GetAlternateMessageId() *string {
	return i.icuTemplateOptions.GetAlternateMessageId()
}

// jsMessageVisitor_icuTemplateDefinition is the anonymous
// IcuTemplateDefinition of extractIcuTemplateDefinition.
type jsMessageVisitor_icuTemplateDefinition struct {
	extractedMessage        *JsMessage
	msgNode                 rhino.Node
	stringLiteralExpression rhino.Node
}

// Java: JsMessageVisitor.extractIcuTemplateDefinition(...).new IcuTemplateDefinition() {...}.getMessage()
func (d *jsMessageVisitor_icuTemplateDefinition) GetMessage() *JsMessage { return d.extractedMessage }

// Java: JsMessageVisitor.extractIcuTemplateDefinition(...).new IcuTemplateDefinition() {...}.getMessageNode()
func (d *jsMessageVisitor_icuTemplateDefinition) GetMessageNode() rhino.Node { return d.msgNode }

// Java: JsMessageVisitor.extractIcuTemplateDefinition(...).new IcuTemplateDefinition() {...}.getTemplateTextNode()
func (d *jsMessageVisitor_icuTemplateDefinition) GetTemplateTextNode() rhino.Node {
	return d.stringLiteralExpression
}

var jsMessageVisitor_ICU_PLACEHOLDER_RE = jutil.Pattern_Compile(`\{([A-Z_0-9]+)\}`)

// JsMessageVisitor_IcuMessageTemplateString is the template string of a
// declareIcuTemplate() call.
//
// Java: static final class JsMessageVisitor.IcuMessageTemplateString
type JsMessageVisitor_IcuMessageTemplateString struct {
	template string
}

// Java: JsMessageVisitor.IcuMessageTemplateString(String)
func NewJsMessageVisitor_IcuMessageTemplateString(template string) *JsMessageVisitor_IcuMessageTemplateString {
	return &JsMessageVisitor_IcuMessageTemplateString{template: template}
}

// ExtractParts splits the template string into plain string parts and
// placeholders.
//
// Although this method will recognize all ICU placeholders in the string
// ("{NAME}"), it will only extract as separate parts those that are named
// in placholderNamesToExtract. Those are the ones for which we will need to
// create placeholder elements, so we can attach example text or original
// code snippets to them when we put the message into an XMB file.
//
// Java: JsMessageVisitor.IcuMessageTemplateString.extractParts(Set<String>)
func (s *JsMessageVisitor_IcuMessageTemplateString) ExtractParts(
	placholderNamesToExtract *jutil.LinkedHashSet[string]) *JsMessageVisitor_ExtractedIcuTemplateParts {
	extractedPlaceholderNames := jutil.NewLinkedHashSet[string]()
	var partsBuilder []JsMessage_Part
	matcher := jsMessageVisitor_ICU_PLACEHOLDER_RE.Matcher(s.template)
	remainingTemplateStartIndex := 0
	for matcher.Find() {
		placeholderName := *matcher.Group(1)
		if placholderNamesToExtract.Contains(placeholderName) {
			extractedPlaceholderNames.Add(placeholderName)
			nonPlaceholderPrefix :=
				jutil.StringSubstring(s.template, remainingTemplateStartIndex, matcher.Start())
			remainingTemplateStartIndex = matcher.End()
			partsBuilder = append(partsBuilder, JsMessage_StringPart_Create(nonPlaceholderPrefix))
			partsBuilder = append(partsBuilder, JsMessage_PlaceholderReference_CreateForCanonicalName(placeholderName))
		} // else we have no need to create a separate part for this placeholder.
	}
	if remainingTemplateStartIndex < jutil.StringLength(s.template) {
		remainingTemplate := jutil.StringSubstringFrom(s.template, remainingTemplateStartIndex)
		partsBuilder = append(partsBuilder, JsMessage_StringPart_Create(remainingTemplate))
	}
	return NewJsMessageVisitor_ExtractedIcuTemplateParts(extractedPlaceholderNames, partsBuilder)
}

// JsMessageVisitor_ExtractedIcuTemplateParts is the result of
// IcuMessageTemplateString.ExtractParts.
//
// Java: static class JsMessageVisitor.ExtractedIcuTemplateParts
type JsMessageVisitor_ExtractedIcuTemplateParts struct {
	extractedPlaceholderNames *jutil.LinkedHashSet[string]
	extractedParts            []JsMessage_Part
}

// Java: JsMessageVisitor.ExtractedIcuTemplateParts(ImmutableSet<String>, ImmutableList<Part>)
func NewJsMessageVisitor_ExtractedIcuTemplateParts(
	extractedPlaceholderNames *jutil.LinkedHashSet[string], extractedParts []JsMessage_Part) *JsMessageVisitor_ExtractedIcuTemplateParts {
	return &JsMessageVisitor_ExtractedIcuTemplateParts{
		extractedPlaceholderNames: extractedPlaceholderNames,
		extractedParts:            extractedParts,
	}
}

// Java: JsMessageVisitor.extractIcuMessageTemplateString(Node)
func (b *JsMessageVisitorBase) extractIcuMessageTemplateString(stringExpression rhino.Node) *JsMessageVisitor_IcuMessageTemplateString {
	templateString := jsMessageVisitor_extractStringFromStringExprNode(stringExpression)
	return NewJsMessageVisitor_IcuMessageTemplateString(templateString)
}

// Java: private interface JsMessageVisitor.JsMessageOptions
type jsMessageVisitor_JsMessageOptions interface {
	// IsEscapeLessThan: Replace `'<'` with `'&lt;'` in the message.
	//
	// Java: JsMessageVisitor.JsMessageOptions.isEscapeLessThan()
	IsEscapeLessThan() bool

	// IsUnescapeHtmlEntities: Replace these escaped entities with their
	// literal characters in the message (Overrides escapeLessThan)
	//
	//	'&lt;' -> '<'
	//	'&gt;' -> '>'
	//	'&apos;' -> "'"
	//	'&quot;' -> '"'
	//	'&amp;' -> '&'
	//
	// Java: JsMessageVisitor.JsMessageOptions.isUnescapeHtmlEntities()
	IsUnescapeHtmlEntities() bool

	// Java: JsMessageVisitor.JsMessageOptions.getPlaceholderExampleMap()
	GetPlaceholderExampleMap() *jutil.LinkedHashMap[string, string]

	// Java: JsMessageVisitor.JsMessageOptions.getPlaceholderOriginalCodeMap()
	GetPlaceholderOriginalCodeMap() *jutil.LinkedHashMap[string, string]

	// Java: JsMessageVisitor.JsMessageOptions.checkForUnknownPlaceholders(Set<String>)
	CheckForUnknownPlaceholders(knownPlaceholders *jutil.LinkedHashSet[string])
}

var jsMessageVisitor_MESSAGE_OPTION_NAMES = jutil.NewLinkedHashSet("html", "unescapeHtmlEntities", "example", "original_code")

// Java: JsMessageVisitor.extractJsMessageOptions(Node)
func (b *JsMessageVisitorBase) extractJsMessageOptions(optionsBag rhino.Node) jsMessageVisitor_JsMessageOptions {
	objectLiteralMap := JsMessageVisitor_ExtractObjectLiteralMap(optionsBag)
	objectLiteralMap.CheckForUnexpectedKeys(
		jsMessageVisitor_MESSAGE_OPTION_NAMES, func(optionName string) string { return "Unknown option: " + optionName })

	isEscapeLessThan := objectLiteralMap.GetBooleanValueOrFalse("html")
	isUnescapeHtmlEntities :=
		objectLiteralMap.GetBooleanValueOrFalse("unescapeHtmlEntities")

	exampleValueNode := objectLiteralMap.GetValueNode("example")
	exampleObjectLiteralMap := JsMessageVisitor_ExtractObjectLiteralMap(exampleValueNode)
	placeholderExamplesMap :=
		exampleObjectLiteralMap.ExtractAsStringToStringMap()

	originalCodeValueNode := objectLiteralMap.GetValueNode("original_code")
	originalCodeObjectLiteralMap :=
		JsMessageVisitor_ExtractObjectLiteralMap(originalCodeValueNode)
	placeholderOriginalCodeMap :=
		originalCodeObjectLiteralMap.ExtractAsStringToStringMap()

	// NOTE: The getX() methods below should all do little to no computation.
	// In particular, all checking for MalformedExceptions must be done before creating this object.
	return &jsMessageVisitor_jsMessageOptions{
		isEscapeLessThan:             isEscapeLessThan,
		isUnescapeHtmlEntities:       isUnescapeHtmlEntities,
		placeholderExamplesMap:       placeholderExamplesMap,
		placeholderOriginalCodeMap:   placeholderOriginalCodeMap,
		exampleObjectLiteralMap:      exampleObjectLiteralMap,
		originalCodeObjectLiteralMap: originalCodeObjectLiteralMap,
	}
}

// jsMessageVisitor_jsMessageOptions is the anonymous JsMessageOptions of
// extractJsMessageOptions.
type jsMessageVisitor_jsMessageOptions struct {
	isEscapeLessThan             bool
	isUnescapeHtmlEntities       bool
	placeholderExamplesMap       *jutil.LinkedHashMap[string, string]
	placeholderOriginalCodeMap   *jutil.LinkedHashMap[string, string]
	exampleObjectLiteralMap      JsMessageVisitor_ObjectLiteralMap
	originalCodeObjectLiteralMap JsMessageVisitor_ObjectLiteralMap
}

// Java: JsMessageVisitor.extractJsMessageOptions(...).new JsMessageOptions() {...}.isEscapeLessThan()
func (o *jsMessageVisitor_jsMessageOptions) IsEscapeLessThan() bool { return o.isEscapeLessThan }

// Java: JsMessageVisitor.extractJsMessageOptions(...).new JsMessageOptions() {...}.isUnescapeHtmlEntities()
func (o *jsMessageVisitor_jsMessageOptions) IsUnescapeHtmlEntities() bool {
	return o.isUnescapeHtmlEntities
}

// Java: JsMessageVisitor.extractJsMessageOptions(...).new JsMessageOptions() {...}.getPlaceholderExampleMap()
func (o *jsMessageVisitor_jsMessageOptions) GetPlaceholderExampleMap() *jutil.LinkedHashMap[string, string] {
	return o.placeholderExamplesMap
}

// Java: JsMessageVisitor.extractJsMessageOptions(...).new JsMessageOptions() {...}.getPlaceholderOriginalCodeMap()
func (o *jsMessageVisitor_jsMessageOptions) GetPlaceholderOriginalCodeMap() *jutil.LinkedHashMap[string, string] {
	return o.placeholderOriginalCodeMap
}

// Java: JsMessageVisitor.extractJsMessageOptions(...).new JsMessageOptions() {...}.checkForUnknownPlaceholders(Set<String>)
func (o *jsMessageVisitor_jsMessageOptions) CheckForUnknownPlaceholders(knownPlaceholders *jutil.LinkedHashSet[string]) {
	o.exampleObjectLiteralMap.CheckForUnexpectedKeys(
		knownPlaceholders, func(unknownName string) string { return "Unknown placeholder: " + unknownName })
	o.originalCodeObjectLiteralMap.CheckForUnexpectedKeys(
		knownPlaceholders, func(unknownName string) string { return "Unknown placeholder: " + unknownName })
}

// Java: private interface JsMessageVisitor.IcuTemplateOptions
type jsMessageVisitor_IcuTemplateOptions interface {
	// Java: JsMessageVisitor.IcuTemplateOptions.getDescription()
	GetDescription() string

	// Java: JsMessageVisitor.IcuTemplateOptions.getMeaning()
	GetMeaning() *string

	// Java: JsMessageVisitor.IcuTemplateOptions.getAlternateMessageId()
	GetAlternateMessageId() *string

	// Java: JsMessageVisitor.IcuTemplateOptions.getPlaceholderExampleMap()
	GetPlaceholderExampleMap() *jutil.LinkedHashMap[string, string]

	// Java: JsMessageVisitor.IcuTemplateOptions.getPlaceholderOriginalCodeMap()
	GetPlaceholderOriginalCodeMap() *jutil.LinkedHashMap[string, string]

	// Java: JsMessageVisitor.IcuTemplateOptions.checkForUnknownPlaceholders(Set<String>)
	CheckForUnknownPlaceholders(knownPlaceholders *jutil.LinkedHashSet[string])

	// GetPlaceholderNames returns all the placeholder names mentioned in the
	// options.
	//
	// Java: JsMessageVisitor.IcuTemplateOptions.getPlaceholderNames()
	GetPlaceholderNames() *jutil.LinkedHashSet[string]
}

// jsMessageVisitor_ICU_TEMPLATE_OPTION_NAMES are the property names that are
// expected to appear in the options bag argument of declareIcuTemplate().
var jsMessageVisitor_ICU_TEMPLATE_OPTION_NAMES = jutil.NewLinkedHashSet(
	"description", "meaning", "alternate_message_id", "example", "original_code")

// extractIcuTemplateOptions extracts the options from the second argument
// to a `declareIcuTemplate()` call
//
// Example:
//
//	{
//	  description: 'required message description here',
//	  meaning: 'optional meaning here',
//	  alternate_message_id: 'optional alternate message ID here',
//	  // optional
//	  original_code: {
//	    'INTERPOLATION_1': '{{getPerson()}}',
//	    'INTERPOLATION_2': '{{getPlaceName()}}',
//	  },
//	  // optional
//	  example: {
//	    'INTERPOLATION_1': 'Jane Doe',
//	    'INTERPOLATION_2': 'Paris, France',
//	  }
//	}
//
// Java: JsMessageVisitor.extractIcuTemplateOptions(Node)
func (b *JsMessageVisitorBase) extractIcuTemplateOptions(optionsBag rhino.Node) jsMessageVisitor_IcuTemplateOptions {
	objectLiteralMap := JsMessageVisitor_ExtractObjectLiteralMap(optionsBag)
	objectLiteralMap.CheckForUnexpectedKeys(
		jsMessageVisitor_ICU_TEMPLATE_OPTION_NAMES, func(optionName string) string { return "Unknown option: " + optionName })

	// required description string
	descriptionNode := objectLiteralMap.GetValueNode("description")
	if descriptionNode.IsNil() {
		panic(NewJsMessageVisitor_MalformedException("'description' option field is missing", optionsBag))
	}
	description := jsMessageVisitor_extractStringFromStringExprNode(descriptionNode)

	// optional meaning string
	meaningNode := objectLiteralMap.GetValueNode("meaning")
	var meaning *string
	if !meaningNode.IsNil() {
		meaning = jutil.StringPtr(jsMessageVisitor_extractStringFromStringExprNode(meaningNode))
	}

	// optional alternate message ID
	alternateMessageIdNode :=
		objectLiteralMap.GetValueNode("alternate_message_id")
	var alternateMessageId *string
	if !alternateMessageIdNode.IsNil() {
		alternateMessageId = jutil.StringPtr(jsMessageVisitor_extractStringFromStringExprNode(alternateMessageIdNode))
	}

	// optional map of placeholder names to example string values
	exampleValueNode := objectLiteralMap.GetValueNode("example")
	exampleObjectLiteralMap := JsMessageVisitor_ExtractObjectLiteralMap(exampleValueNode)
	placeholderExamplesMap :=
		exampleObjectLiteralMap.ExtractAsStringToStringMap()

	// optional map of placeholder names to original_code string values
	originalCodeValueNode := objectLiteralMap.GetValueNode("original_code")
	originalCodeObjectLiteralMap :=
		JsMessageVisitor_ExtractObjectLiteralMap(originalCodeValueNode)
	placeholderOriginalCodeMap :=
		originalCodeObjectLiteralMap.ExtractAsStringToStringMap()

	// All the placeholder names mentioned in the 2 optional maps.
	placeholderNames := jutil.NewLinkedHashSet[string]()
	for k := range placeholderExamplesMap.Keys() {
		placeholderNames.Add(k)
	}
	for k := range placeholderOriginalCodeMap.Keys() {
		placeholderNames.Add(k)
	}

	for placeholderName := range placeholderNames.All() {
		if !JsMessage_IsCanonicalPlaceholderNameFormat(placeholderName) {
			panic(NewJsMessageVisitor_MalformedException(
				"Placeholder not in UPPER_SNAKE_CASE: "+placeholderName, optionsBag))
		}
	}

	// NOTE: The getX() methods below should all do little to no computation.
	// In particular, all checking for MalformedExceptions must be done before creating this object.
	return &jsMessageVisitor_icuTemplateOptions{
		description:                  description,
		meaning:                      meaning,
		alternateMessageId:           alternateMessageId,
		placeholderExamplesMap:       placeholderExamplesMap,
		placeholderOriginalCodeMap:   placeholderOriginalCodeMap,
		placeholderNames:             placeholderNames,
		exampleObjectLiteralMap:      exampleObjectLiteralMap,
		originalCodeObjectLiteralMap: originalCodeObjectLiteralMap,
	}
}

// jsMessageVisitor_icuTemplateOptions is the anonymous IcuTemplateOptions
// of extractIcuTemplateOptions.
type jsMessageVisitor_icuTemplateOptions struct {
	description                  string
	meaning                      *string
	alternateMessageId           *string
	placeholderExamplesMap       *jutil.LinkedHashMap[string, string]
	placeholderOriginalCodeMap   *jutil.LinkedHashMap[string, string]
	placeholderNames             *jutil.LinkedHashSet[string]
	exampleObjectLiteralMap      JsMessageVisitor_ObjectLiteralMap
	originalCodeObjectLiteralMap JsMessageVisitor_ObjectLiteralMap
}

// Java: JsMessageVisitor.extractIcuTemplateOptions(...).new IcuTemplateOptions() {...}.getDescription()
func (o *jsMessageVisitor_icuTemplateOptions) GetDescription() string { return o.description }

// Java: JsMessageVisitor.extractIcuTemplateOptions(...).new IcuTemplateOptions() {...}.getMeaning()
func (o *jsMessageVisitor_icuTemplateOptions) GetMeaning() *string { return o.meaning }

// Java: JsMessageVisitor.extractIcuTemplateOptions(...).new IcuTemplateOptions() {...}.getAlternateMessageId()
func (o *jsMessageVisitor_icuTemplateOptions) GetAlternateMessageId() *string {
	return o.alternateMessageId
}

// Java: JsMessageVisitor.extractIcuTemplateOptions(...).new IcuTemplateOptions() {...}.getPlaceholderExampleMap()
func (o *jsMessageVisitor_icuTemplateOptions) GetPlaceholderExampleMap() *jutil.LinkedHashMap[string, string] {
	return o.placeholderExamplesMap
}

// Java: JsMessageVisitor.extractIcuTemplateOptions(...).new IcuTemplateOptions() {...}.getPlaceholderOriginalCodeMap()
func (o *jsMessageVisitor_icuTemplateOptions) GetPlaceholderOriginalCodeMap() *jutil.LinkedHashMap[string, string] {
	return o.placeholderOriginalCodeMap
}

// Java: JsMessageVisitor.extractIcuTemplateOptions(...).new IcuTemplateOptions() {...}.getPlaceholderNames()
func (o *jsMessageVisitor_icuTemplateOptions) GetPlaceholderNames() *jutil.LinkedHashSet[string] {
	return o.placeholderNames
}

// Java: JsMessageVisitor.extractIcuTemplateOptions(...).new IcuTemplateOptions() {...}.checkForUnknownPlaceholders(Set<String>)
func (o *jsMessageVisitor_icuTemplateOptions) CheckForUnknownPlaceholders(knownPlaceholders *jutil.LinkedHashSet[string]) {
	o.exampleObjectLiteralMap.CheckForUnexpectedKeys(
		knownPlaceholders, func(unknownName string) string { return "Unknown placeholder: " + unknownName })
	o.originalCodeObjectLiteralMap.CheckForUnexpectedKeys(
		knownPlaceholders, func(unknownName string) string { return "Unknown placeholder: " + unknownName })
}

// Java: JsMessageVisitor.extractBooleanStringKeyValue(Node)
func jsMessageVisitor_extractBooleanStringKeyValue(stringKeyNode rhino.Node) bool {
	if stringKeyNode.IsNil() {
		return false
	} else {
		valueNode := stringKeyNode.GetOnlyChild()
		if valueNode.IsTrue() {
			return true
		} else if valueNode.IsFalse() {
			return false
		} else {
			panic(NewJsMessageVisitor_MalformedException(
				stringKeyNode.GetString()+": Literal true or false expected", valueNode))
		}
	}
}

// JsMessageVisitor_ObjectLiteralMap represents the contents of a object
// literal Node in the AST.
//
// The object literal may not have any computed keys or methods.
//
// Use JsMessageVisitor_ExtractObjectLiteralMap(Node) to get an instance.
//
// Java: public interface JsMessageVisitor.ObjectLiteralMap
type JsMessageVisitor_ObjectLiteralMap interface {
	// Java: JsMessageVisitor.ObjectLiteralMap.getBooleanValueOrFalse(String)
	GetBooleanValueOrFalse(key string) bool

	// ExtractAsValueMap returns a map from object property names the Node
	// values they have in the AST, building the map first, if necessary.
	//
	// Java: JsMessageVisitor.ObjectLiteralMap.extractAsValueMap()
	ExtractAsValueMap() *jutil.LinkedHashMap[string, rhino.Node]

	// ExtractAsStringToStringMap returns a map from object property names to
	// string values, building it first, if necessary.
	//
	// Panics with a JsMessageVisitor_MalformedException if any of the values
	// are not actually simple string literals or concatenations of string
	// literals.
	//
	// Java: JsMessageVisitor.ObjectLiteralMap.extractAsStringToStringMap()
	ExtractAsStringToStringMap() *jutil.LinkedHashMap[string, string]

	// GetValueNode gets the value node for a key.
	//
	// This method avoids the work of extracting the full value map.
	//
	// Java: JsMessageVisitor.ObjectLiteralMap.getValueNode(String)
	GetValueNode(key string) rhino.Node

	// CheckForUnexpectedKeys panics with a
	// JsMessageVisitor_MalformedException if the object literal has keys not
	// in the expected set.
	//
	// Java: JsMessageVisitor.ObjectLiteralMap.checkForUnexpectedKeys(Set<String>, Function<String, String>)
	CheckForUnexpectedKeys(expectedKeys *jutil.LinkedHashSet[string], createErrorMessage func(string) string)

	// CheckForRequiredKeys panics with a JsMessageVisitor_MalformedException
	// if the object literal does not have one all the keys in the required
	// set.
	//
	// Java: JsMessageVisitor.ObjectLiteralMap.checkForRequiredKeys(Set<String>, Function<String, MalformedException>)
	CheckForRequiredKeys(
		requiredKeys *jutil.LinkedHashSet[string], createException func(string) *JsMessageVisitor_MalformedException)
}

// Java: private static class JsMessageVisitor.ObjectLiteralMapImpl
type jsMessageVisitor_ObjectLiteralMapImpl struct {
	stringToStringKeyMap *jutil.LinkedHashMap[string, rhino.Node]
	// The result of extractValueMap(). It will be populated if requested.
	valueMap *jutil.LinkedHashMap[string, rhino.Node]
	// The result of extractStringMap(). It will be populated if requested.
	stringMap *jutil.LinkedHashMap[string, string]
}

// Java: JsMessageVisitor.ObjectLiteralMapImpl(Map<String, Node>)
func newJsMessageVisitor_ObjectLiteralMapImpl(stringToStringKeyMap *jutil.LinkedHashMap[string, rhino.Node]) *jsMessageVisitor_ObjectLiteralMapImpl {
	return &jsMessageVisitor_ObjectLiteralMapImpl{stringToStringKeyMap: stringToStringKeyMap}
}

// Java: JsMessageVisitor.ObjectLiteralMapImpl.getValueNode(String)
func (m *jsMessageVisitor_ObjectLiteralMapImpl) GetValueNode(key string) rhino.Node {
	stringKeyNode := m.stringToStringKeyMap.Get(key)
	if stringKeyNode.IsNil() {
		return rhino.NilNode
	}
	return stringKeyNode.GetOnlyChild()
}

// Java: JsMessageVisitor.ObjectLiteralMapImpl.getBooleanValueOrFalse(String)
func (m *jsMessageVisitor_ObjectLiteralMapImpl) GetBooleanValueOrFalse(key string) bool {
	stringKeyNode := m.stringToStringKeyMap.Get(key)
	return jsMessageVisitor_extractBooleanStringKeyValue(stringKeyNode)
}

// Java: JsMessageVisitor.ObjectLiteralMapImpl.extractAsValueMap()
func (m *jsMessageVisitor_ObjectLiteralMapImpl) ExtractAsValueMap() *jutil.LinkedHashMap[string, rhino.Node] {
	if m.valueMap == nil {
		builder := jutil.NewImmutableMapBuilder[string, rhino.Node]()
		for key, value := range m.stringToStringKeyMap.All() {
			builder.Put(key, value.GetOnlyChild())
		}
		m.valueMap = builder.BuildOrThrow()
	}
	return m.valueMap
}

// Java: JsMessageVisitor.ObjectLiteralMapImpl.extractAsStringToStringMap()
func (m *jsMessageVisitor_ObjectLiteralMapImpl) ExtractAsStringToStringMap() *jutil.LinkedHashMap[string, string] {
	if m.stringMap == nil {
		builder := jutil.NewImmutableMapBuilder[string, string]()
		for key, value := range m.stringToStringKeyMap.All() {
			builder.Put(
				key, jsMessageVisitor_extractStringFromStringExprNode(value.GetOnlyChild()))
		}
		m.stringMap = builder.BuildOrThrow()
	}
	return m.stringMap
}

// Java: JsMessageVisitor.ObjectLiteralMapImpl.checkForUnexpectedKeys(Set<String>, Function<String, String>)
func (m *jsMessageVisitor_ObjectLiteralMapImpl) CheckForUnexpectedKeys(
	expectedKeys *jutil.LinkedHashSet[string], createErrorMessage func(string) string) {
	for key := range m.stringToStringKeyMap.Keys() {
		if !expectedKeys.Contains(key) {
			panic(NewJsMessageVisitor_MalformedException(
				createErrorMessage(key), m.stringToStringKeyMap.Get(key)))
		}
	}
}

// Java: JsMessageVisitor.ObjectLiteralMapImpl.checkForRequiredKeys(Set<String>, Function<String, MalformedException>)
func (m *jsMessageVisitor_ObjectLiteralMapImpl) CheckForRequiredKeys(
	requiredKeys *jutil.LinkedHashSet[string], createException func(string) *JsMessageVisitor_MalformedException) {
	for requiredKey := range requiredKeys.All() {
		if !m.stringToStringKeyMap.ContainsKey(requiredKey) {
			panic(createException(requiredKey))
		}
	}
}

// JsMessageVisitor_ExtractObjectLiteralMap returns an object to represent an
// object literal Node from the AST.
//
// The object literal's members must all be `STRING_KEY` nodes.
//
// No duplicate string keys are allowed.
//
// A `null` argument is treated as if it were an empty object literal.
//
// Panics with a JsMessageVisitor_MalformedException if the Node does not
// meet the above requirements.
//
// Java: JsMessageVisitor.extractObjectLiteralMap(Node)
func JsMessageVisitor_ExtractObjectLiteralMap(objLit rhino.Node) JsMessageVisitor_ObjectLiteralMap {
	if objLit.IsNil() {
		return newJsMessageVisitor_ObjectLiteralMapImpl(jutil.NewLinkedHashMap[string, rhino.Node]())
	}
	if !objLit.IsObjectLit() {
		panic(NewJsMessageVisitor_MalformedException("object literal expected", objLit))
	}
	stringToStringKeyMap := jutil.NewLinkedHashMap[string, rhino.Node]()
	for stringKey := objLit.GetFirstChild(); !stringKey.IsNil(); stringKey = stringKey.GetNext() {
		if !stringKey.IsStringKey() {
			panic(NewJsMessageVisitor_MalformedException("string key expected", stringKey))
		}
		key := stringKey.GetString()
		if stringToStringKeyMap.ContainsKey(key) {
			panic(NewJsMessageVisitor_MalformedException("duplicate string key: "+key, stringKey))
		}
		stringToStringKeyMap.Put(stringKey.GetString(), stringKey)
	}
	return newJsMessageVisitor_ObjectLiteralMapImpl(stringToStringKeyMap)
}

// JsMessageVisitor_ParseJsMessageTextIntoParts panics with an
// IJsMessage_PlaceholderFormatException on a malformed placeholder.
//
// Java: JsMessageVisitor.parseJsMessageTextIntoParts(String)
func JsMessageVisitor_ParseJsMessageTextIntoParts(originalMsgText string) []JsMessage_Part {
	googGetMsgParsedText := jsMessageVisitor_extractGoogGetMsgParsedText(originalMsgText)
	return googGetMsgParsedText.GetParts()
}

// Java: JsMessageVisitor.extractGoogGetMsgParsedText(String)
func jsMessageVisitor_extractGoogGetMsgParsedText(originalMsgText string) jsMessageVisitor_GoogGetMsgParsedText {
	msgText := originalMsgText
	var partsBuilder []JsMessage_Part
	placeholderNamesBuilder := jutil.NewLinkedHashSet[string]()
	for {
		phBegin := jutil.StringIndexOf(msgText, JsMessage_PH_JS_PREFIX)
		if phBegin < 0 {
			// Just a string literal
			partsBuilder = append(partsBuilder, JsMessage_StringPart_Create(msgText))
			break
		} else {
			if phBegin > 0 {
				// A string literal followed by a placeholder
				partsBuilder = append(partsBuilder, JsMessage_StringPart_Create(jutil.StringSubstring(msgText, 0, phBegin)))
			}

			// A placeholder. Find where it ends
			phEnd := jutil.StringIndexOfFrom(msgText, JsMessage_PH_JS_SUFFIX, phBegin)
			if phEnd < 0 {
				panic(NewJsMessage_PlaceholderFormatException("Placeholder incorrectly formatted"))
			}

			phName := jutil.StringSubstring(msgText, phBegin+jutil.StringLength(JsMessage_PH_JS_PREFIX), phEnd)
			if !JsMessage_IsLowerCamelCaseWithNumericSuffixes(phName) {
				panic(NewJsMessage_PlaceholderFormatException("Placeholder name not in lowerCamelCase: " + phName))
			}
			placeholderNamesBuilder.Add(phName)
			partsBuilder = append(partsBuilder, JsMessage_PlaceholderReference_CreateForJsName(phName))
			nextPos := phEnd + jutil.StringLength(JsMessage_PH_JS_SUFFIX)
			if nextPos < jutil.StringLength(msgText) {
				// Iterate on the rest of the message value
				msgText = jutil.StringSubstringFrom(msgText, nextPos)
			} else {
				// The message is parsed
				break
			}
		}
	}
	placeholderNames := placeholderNamesBuilder
	parts := partsBuilder
	// NOTE: The methods defined below should do little to no computation, and definitely should
	// not throw exceptions.
	return &jsMessageVisitor_googGetMsgParsedText{
		originalMsgText:  originalMsgText,
		placeholderNames: placeholderNames,
		parts:            parts,
	}
}

// Java: private interface JsMessageVisitor.GoogGetMsgParsedText
type jsMessageVisitor_GoogGetMsgParsedText interface {
	// Java: JsMessageVisitor.GoogGetMsgParsedText.getText()
	GetText() string

	// Java: JsMessageVisitor.GoogGetMsgParsedText.getPlaceholderNames()
	GetPlaceholderNames() *jutil.LinkedHashSet[string]

	// Java: JsMessageVisitor.GoogGetMsgParsedText.getParts()
	GetParts() []JsMessage_Part
}

// jsMessageVisitor_googGetMsgParsedText is the anonymous
// GoogGetMsgParsedText of extractGoogGetMsgParsedText.
type jsMessageVisitor_googGetMsgParsedText struct {
	originalMsgText  string
	placeholderNames *jutil.LinkedHashSet[string]
	parts            []JsMessage_Part
}

// Java: JsMessageVisitor.extractGoogGetMsgParsedText(...).new GoogGetMsgParsedText() {...}.getText()
func (p *jsMessageVisitor_googGetMsgParsedText) GetText() string { return p.originalMsgText }

// Java: JsMessageVisitor.extractGoogGetMsgParsedText(...).new GoogGetMsgParsedText() {...}.getPlaceholderNames()
func (p *jsMessageVisitor_googGetMsgParsedText) GetPlaceholderNames() *jutil.LinkedHashSet[string] {
	return p.placeholderNames
}

// Java: JsMessageVisitor.extractGoogGetMsgParsedText(...).new GoogGetMsgParsedText() {...}.getParts()
func (p *jsMessageVisitor_googGetMsgParsedText) GetParts() []JsMessage_Part { return p.parts }

// visitFallbackFunctionCall visits a call to goog.getMsgWithFallback.
//
// Java: JsMessageVisitor.visitFallbackFunctionCall(NodeTraversal, Node)
func (b *JsMessageVisitorBase) visitFallbackFunctionCall(t *NodeTraversal, call rhino.Node) {
	// Check to make sure the function call looks like:
	// goog.getMsgWithFallback(MSG_1, MSG_2);
	// or:
	// goog.getMsgWithFallback(some.import.MSG_1, some.import.MSG_2);
	if !call.HasXChildren(3) ||
		!jsMessageVisitor_isMessageIdentifier(call.GetSecondChild()) ||
		!jsMessageVisitor_isMessageIdentifier(call.GetLastChild()) {
		b.compiler.Report(JSError_MakeWithNode(call, JsMessageVisitor_BAD_FALLBACK_SYNTAX))
		return
	}

	firstArg := call.GetSecondChild()
	firstMessage := b.getJsMessageFromNode(t, firstArg)
	if firstMessage == nil {
		b.compiler.Report(JSError_MakeWithNode(firstArg, JsMessageVisitor_FALLBACK_ARG_ERROR, *firstArg.GetQualifiedName()))
		return
	}

	secondArg := firstArg.GetNext()
	secondMessage := b.getJsMessageFromNode(t, secondArg)
	if secondMessage == nil {
		b.compiler.Report(JSError_MakeWithNode(secondArg, JsMessageVisitor_FALLBACK_ARG_ERROR, *secondArg.GetQualifiedName()))
		return
	}

	b.self.ProcessMessageFallback(call, firstMessage, secondMessage)
}

// ProcessMessageFallback processes the goog.getMsgWithFallback primitive.
// goog.getMsgWithFallback(MSG_1, MSG_2);
//
// By default, does nothing.
//
// Java: JsMessageVisitor.processMessageFallback(Node, JsMessage, JsMessage)
func (b *JsMessageVisitorBase) ProcessMessageFallback(callNode rhino.Node, message1 *JsMessage, message2 *JsMessage) {
}

// IsMessageName returns whether the given JS identifier is a valid JS
// message name.
//
// Java: JsMessageVisitor.isMessageName(String)
func (b *JsMessageVisitorBase) IsMessageName(identifier string) bool {
	return jutil.StringStartsWith(identifier, JsMessageVisitor_MSG_PREFIX) || JsMessageVisitor_IsScopedAliasesPrefix(identifier)
}

// Java: JsMessageVisitor.isMessageIdentifier(Node)
func jsMessageVisitor_isMessageIdentifier(node rhino.Node) bool {
	qname := node.GetQualifiedName()
	return qname != nil && jutil.StringIndexOf(*qname, JsMessageVisitor_MSG_PREFIX) >= 0
}

// getJsMessageFromNode extracts a message name (e.g. MSG_FOO) from either a
// NAME node or a GETPROP node. This should cover all of the following
// cases:
//
//  1. a NAME node (e.g. MSG_FOO)
//  2. a NAME node which is the product of renaming (e.g.
//     $module$contents$MSG_FOO)
//  3. a GETPROP node (e.g. some.import.MSG_FOO)
//
// Java: JsMessageVisitor.getJsMessageFromNode(NodeTraversal, Node)
func (b *JsMessageVisitorBase) getJsMessageFromNode(t *NodeTraversal, node rhino.Node) *JsMessage {
	messageNamePtr := node.GetQualifiedName()
	if messageNamePtr == nil || jutil.StringIndexOf(*messageNamePtr, JsMessageVisitor_MSG_PREFIX) < 0 {
		return nil
	}
	messageName := *messageNamePtr

	messageKey := jutil.StringSubstringFrom(messageName, jutil.StringIndexOf(messageName, JsMessageVisitor_MSG_PREFIX))
	if jsMessageVisitor_isUnnamedMessageName(messageKey) {
		return b.getTrackedUnnamedMessage(t, messageName)
	} else {
		return b.getTrackedNormalMessage(messageKey)
	}
}

// jsMessageVisitor_isUnnamedMessageName returns whether the given message
// name is in the unnamed namespace.
//
// Java: JsMessageVisitor.isUnnamedMessageName(String)
func jsMessageVisitor_isUnnamedMessageName(identifier string) bool {
	return jsMessageVisitor_MSG_UNNAMED_PATTERN.Matcher(identifier).Matches()
}

// CheckNode checks a node's type.
//
// Panics with a JsMessageVisitor_MalformedException if the node is null or
// the wrong type
//
// Java: JsMessageVisitor.checkNode(Node, Token)
func (b *JsMessageVisitorBase) CheckNode(node rhino.Node, type_ rhino.Token) {
	if node.IsNil() {
		panic(NewJsMessageVisitor_MalformedException("Expected node type "+type_.String()+"; found: null", node))
	}
	if node.GetToken() != type_ {
		panic(NewJsMessageVisitor_MalformedException(
			"Expected node type "+type_.String()+"; found: "+node.GetToken().String(), node))
	}
}

// JsMessageVisitor_MalformedException is thrown when a message definition
// does not have the expected structure.
//
// Java: static class JsMessageVisitor.MalformedException extends Exception
type JsMessageVisitor_MalformedException struct {
	jutil.JavaException

	node rhino.Node
}

func (*JsMessageVisitor_MalformedException) malformedException() {}

// IJsMessageVisitor_MalformedException matches MalformedException.
type IJsMessageVisitor_MalformedException interface {
	jutil.IException
	malformedException()
	GetNode() rhino.Node
}

// Java: JsMessageVisitor.MalformedException(String, Node)
func NewJsMessageVisitor_MalformedException(message string, node rhino.Node) *JsMessageVisitor_MalformedException {
	e := &JsMessageVisitor_MalformedException{}
	jutil.InitThrowable(e, "com.google.javascript.jscomp.JsMessageVisitor$MalformedException", &message, nil)
	e.node = node
	return e
}

// Java: JsMessageVisitor.MalformedException.getNode()
func (e *JsMessageVisitor_MalformedException) GetNode() rhino.Node {
	return e.node
}

// Java: private static class JsMessageVisitor.MessageLocation
type jsMessageVisitor_MessageLocation struct {
	message     *JsMessage
	messageNode rhino.Node
}

// Java: JsMessageVisitor.MessageLocation(JsMessage, Node)
func newJsMessageVisitor_MessageLocation(message *JsMessage, messageNode rhino.Node) *jsMessageVisitor_MessageLocation {
	return &jsMessageVisitor_MessageLocation{message: message, messageNode: messageNode}
}

// Java: JsMessageVisitor.isScopedAliasesPrefix(String)
func JsMessageVisitor_IsScopedAliasesPrefix(name string) bool {
	return jsMessageVisitor_SCOPED_ALIASES_PREFIX_PATTERN.Matcher(name).LookingAt()
}

// Java: JsMessageVisitor.removeScopedAliasesPrefix(String)
func JsMessageVisitor_RemoveScopedAliasesPrefix(name string) string {
	return jsMessageVisitor_SCOPED_ALIASES_PREFIX_PATTERN.Matcher(name).ReplaceFirst("MSG_")
}
