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
// Port of com/google/javascript/jscomp/JsMessageVisitor.java (generated mechanically by j2go).

package jscomp

import (
	"strconv"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_js_message_visitor_1 j2gort.Ext_java_util_function_Function
var lambdaSingleton_js_message_visitor_2 j2gort.Ext_java_util_function_Function
var lambdaSingleton_js_message_visitor_3 j2gort.Ext_java_util_function_Function
var lambdaSingleton_js_message_visitor_4 j2gort.Ext_java_util_function_Function
var lambdaSingleton_js_message_visitor_5 j2gort.Ext_java_util_function_Function
var lambdaSingleton_js_message_visitor_6 j2gort.Ext_java_util_function_Function
var lambdaSingleton_js_message_visitor_7 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_js_message_visitor_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_8 any) any {
		placeholderName := j2gort.CastString(lp_8)
		_ = placeholderName
		return "Unused message placeholder: " + j2gort.StrValString(placeholderName)
	})
	lambdaSingleton_js_message_visitor_2 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_9 any) any {
		optionName := j2gort.CastString(lp_9)
		_ = optionName
		return "Unknown option: " + j2gort.StrValString(optionName)
	})
	lambdaSingleton_js_message_visitor_3 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_10 any) any {
		optionName := j2gort.CastString(lp_10)
		_ = optionName
		return "Unknown option: " + j2gort.StrValString(optionName)
	})
	lambdaSingleton_js_message_visitor_4 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_16 any) any {
		unknownName := j2gort.CastString(lp_16)
		_ = unknownName
		return "Unknown placeholder: " + j2gort.StrValString(unknownName)
	})
	lambdaSingleton_js_message_visitor_5 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_17 any) any {
		unknownName := j2gort.CastString(lp_17)
		_ = unknownName
		return "Unknown placeholder: " + j2gort.StrValString(unknownName)
	})
	lambdaSingleton_js_message_visitor_6 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_18 any) any {
		unknownName := j2gort.CastString(lp_18)
		_ = unknownName
		return "Unknown placeholder: " + j2gort.StrValString(unknownName)
	})
	lambdaSingleton_js_message_visitor_7 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_19 any) any {
		unknownName := j2gort.CastString(lp_19)
		_ = unknownName
		return "Unknown placeholder: " + j2gort.StrValString(unknownName)
	})
}

// Java polymorphic class interface: com.google.javascript.jscomp.JsMessageVisitor
type JsMessageVisitor interface {
	Implements_jscomp_JsMessageVisitor()
	IsJavaNil_() bool
	GetJsMessageVisitorBase() *JsMessageVisitorBase
	NodeTraversal_AbstractPostOrderCallback
	CompilerPass
	ProcessJsMessageDefinition(JsMessageDefinition)
	ProcessIcuTemplateDefinition(IcuTemplateDefinition)
	ProcessMessageFallback(pkg_rhino.Node, *JsMessage, *JsMessage)
	IsMessageName(string) bool
	CheckNode(pkg_rhino.Node, pkg_rhino.Token)
}

func ToIface_JsMessageVisitor[T interface {
	JsMessageVisitor
	comparable
}](v T) JsMessageVisitor {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_JsMessageVisitor[T JsMessageVisitor](v T) JsMessageVisitor {
	var vi JsMessageVisitor = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

// Java class: com.google.javascript.jscomp.JsMessageVisitor
type JsMessageVisitorBase struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	self_           JsMessageVisitor
	idGenerator     JsMessage_IdGenerator                                        `gson:"idGenerator"`
	compiler        AbstractCompiler                                             `gson:"compiler"`
	messageNames    j2gort.OrderedMap[string, *JsMessageVisitor_MessageLocation] `gson:"messageNames"`
	unnamedMessages j2gort.CustomOrderedMap[*Var, *JsMessage]                    `gson:"unnamedMessages"`
	googMsgNodes    j2gort.OrderedSet[pkg_rhino.Node]                            `gson:"googMsgNodes"`
}

func (this_ *JsMessageVisitorBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(JsMessageVisitor)
}
func (this_ *JsMessageVisitorBase) IsJavaNil_() bool { return this_ == nil }
func (*JsMessageVisitorBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JsMessageVisitor", "JsMessageVisitor"
}
func (this_ *JsMessageVisitorBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JsMessageVisitor", "com.google.javascript.jscomp.JsMessageVisitor", func(v_ any) bool {
		return j2gort.InstanceOf[JsMessageVisitor](v_)
	})
}
func (b *JsMessageVisitorBase) GetJsMessageVisitorBase() *JsMessageVisitorBase           { return b }
func (*JsMessageVisitorBase) Implements_jscomp_JsMessageVisitor()                        {}
func (*JsMessageVisitorBase) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*JsMessageVisitorBase) Implements_jscomp_CompilerPass()                            {}
func (*JsMessageVisitorBase) Implements_jscomp_NodeTraversal_Callback()                  {}

func StaticGet_JsMessageVisitor_MSG_FUNCTION_NAME() string     { return "getMsg" }
func StaticGet_JsMessageVisitor_MSG_FUNCTION_QNAME() string    { return "goog.getMsg" }
func StaticGet_JsMessageVisitor_ICU_MSG_FUNCTION_NAME() string { return "declareIcuTemplate" }
func StaticGet_JsMessageVisitor_ICU_MSG_FUNCTION_QNAME() string {
	return "goog.i18n.messages.declareIcuTemplate"
}
func StaticGet_JsMessageVisitor_MSG_FALLBACK_FUNCTION_NAME() string { return "goog.getMsgWithFallback" }

var static_JsMessageVisitor_FUNCTIONS_USED_IN_MESSAGE_EXTRACTION j2gort.Ext_java_util_List

func StaticGet_JsMessageVisitor_FUNCTIONS_USED_IN_MESSAGE_EXTRACTION() j2gort.Ext_java_util_List {
	clinit_JsMessageVisitor()
	return static_JsMessageVisitor_FUNCTIONS_USED_IN_MESSAGE_EXTRACTION
}
func StaticSet_JsMessageVisitor_FUNCTIONS_USED_IN_MESSAGE_EXTRACTION(v j2gort.Ext_java_util_List) {
	clinit_JsMessageVisitor()
	static_JsMessageVisitor_FUNCTIONS_USED_IN_MESSAGE_EXTRACTION = v
}
func StaticGet_JsMessageVisitor_MSG_EXTERNAL_PREFIX() string { return "MSG_EXTERNAL_" }

var static_JsMessageVisitor_MESSAGE_HAS_NO_DESCRIPTION *DiagnosticType

func StaticGet_JsMessageVisitor_MESSAGE_HAS_NO_DESCRIPTION() *DiagnosticType {
	clinit_JsMessageVisitor()
	return static_JsMessageVisitor_MESSAGE_HAS_NO_DESCRIPTION
}
func StaticSet_JsMessageVisitor_MESSAGE_HAS_NO_DESCRIPTION(v *DiagnosticType) {
	clinit_JsMessageVisitor()
	static_JsMessageVisitor_MESSAGE_HAS_NO_DESCRIPTION = v
}

var static_JsMessageVisitor_MESSAGE_HAS_NO_TEXT *DiagnosticType

func StaticGet_JsMessageVisitor_MESSAGE_HAS_NO_TEXT() *DiagnosticType {
	clinit_JsMessageVisitor()
	return static_JsMessageVisitor_MESSAGE_HAS_NO_TEXT
}
func StaticSet_JsMessageVisitor_MESSAGE_HAS_NO_TEXT(v *DiagnosticType) {
	clinit_JsMessageVisitor()
	static_JsMessageVisitor_MESSAGE_HAS_NO_TEXT = v
}

var static_JsMessageVisitor_MESSAGE_TREE_MALFORMED *DiagnosticType

func StaticGet_JsMessageVisitor_MESSAGE_TREE_MALFORMED() *DiagnosticType {
	clinit_JsMessageVisitor()
	return static_JsMessageVisitor_MESSAGE_TREE_MALFORMED
}
func StaticSet_JsMessageVisitor_MESSAGE_TREE_MALFORMED(v *DiagnosticType) {
	clinit_JsMessageVisitor()
	static_JsMessageVisitor_MESSAGE_TREE_MALFORMED = v
}

var static_JsMessageVisitor_MESSAGE_HAS_NO_VALUE *DiagnosticType

func StaticGet_JsMessageVisitor_MESSAGE_HAS_NO_VALUE() *DiagnosticType {
	clinit_JsMessageVisitor()
	return static_JsMessageVisitor_MESSAGE_HAS_NO_VALUE
}
func StaticSet_JsMessageVisitor_MESSAGE_HAS_NO_VALUE(v *DiagnosticType) {
	clinit_JsMessageVisitor()
	static_JsMessageVisitor_MESSAGE_HAS_NO_VALUE = v
}

var static_JsMessageVisitor_MESSAGE_DUPLICATE_KEY *DiagnosticType

func StaticGet_JsMessageVisitor_MESSAGE_DUPLICATE_KEY() *DiagnosticType {
	clinit_JsMessageVisitor()
	return static_JsMessageVisitor_MESSAGE_DUPLICATE_KEY
}
func StaticSet_JsMessageVisitor_MESSAGE_DUPLICATE_KEY(v *DiagnosticType) {
	clinit_JsMessageVisitor()
	static_JsMessageVisitor_MESSAGE_DUPLICATE_KEY = v
}

var static_JsMessageVisitor_MESSAGE_NODE_IS_ORPHANED *DiagnosticType

func StaticGet_JsMessageVisitor_MESSAGE_NODE_IS_ORPHANED() *DiagnosticType {
	clinit_JsMessageVisitor()
	return static_JsMessageVisitor_MESSAGE_NODE_IS_ORPHANED
}
func StaticSet_JsMessageVisitor_MESSAGE_NODE_IS_ORPHANED(v *DiagnosticType) {
	clinit_JsMessageVisitor()
	static_JsMessageVisitor_MESSAGE_NODE_IS_ORPHANED = v
}

var static_JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY *DiagnosticType

func StaticGet_JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY() *DiagnosticType {
	clinit_JsMessageVisitor()
	return static_JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY
}
func StaticSet_JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY(v *DiagnosticType) {
	clinit_JsMessageVisitor()
	static_JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY = v
}

var static_JsMessageVisitor_BAD_FALLBACK_SYNTAX *DiagnosticType

func StaticGet_JsMessageVisitor_BAD_FALLBACK_SYNTAX() *DiagnosticType {
	clinit_JsMessageVisitor()
	return static_JsMessageVisitor_BAD_FALLBACK_SYNTAX
}
func StaticSet_JsMessageVisitor_BAD_FALLBACK_SYNTAX(v *DiagnosticType) {
	clinit_JsMessageVisitor()
	static_JsMessageVisitor_BAD_FALLBACK_SYNTAX = v
}

var static_JsMessageVisitor_FALLBACK_ARG_ERROR *DiagnosticType

func StaticGet_JsMessageVisitor_FALLBACK_ARG_ERROR() *DiagnosticType {
	clinit_JsMessageVisitor()
	return static_JsMessageVisitor_FALLBACK_ARG_ERROR
}
func StaticSet_JsMessageVisitor_FALLBACK_ARG_ERROR(v *DiagnosticType) {
	clinit_JsMessageVisitor()
	static_JsMessageVisitor_FALLBACK_ARG_ERROR = v
}
func StaticGet_JsMessageVisitor_MSG_PREFIX() string { return "MSG_" }

var static_JsMessageVisitor_SCOPED_ALIASES_PREFIX_PATTERN *j2gort.Ext_java_util_regex_Pattern

func StaticGet_JsMessageVisitor_SCOPED_ALIASES_PREFIX_PATTERN() *j2gort.Ext_java_util_regex_Pattern {
	clinit_JsMessageVisitor()
	return static_JsMessageVisitor_SCOPED_ALIASES_PREFIX_PATTERN
}
func StaticSet_JsMessageVisitor_SCOPED_ALIASES_PREFIX_PATTERN(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_JsMessageVisitor()
	static_JsMessageVisitor_SCOPED_ALIASES_PREFIX_PATTERN = v
}

var static_JsMessageVisitor_MSG_UNNAMED_PATTERN *j2gort.Ext_java_util_regex_Pattern

func StaticGet_JsMessageVisitor_MSG_UNNAMED_PATTERN() *j2gort.Ext_java_util_regex_Pattern {
	clinit_JsMessageVisitor()
	return static_JsMessageVisitor_MSG_UNNAMED_PATTERN
}
func StaticSet_JsMessageVisitor_MSG_UNNAMED_PATTERN(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_JsMessageVisitor()
	static_JsMessageVisitor_MSG_UNNAMED_PATTERN = v
}

var static_JsMessageVisitor_ICU_PLACEHOLDER_RE *j2gort.Ext_java_util_regex_Pattern

func StaticGet_JsMessageVisitor_ICU_PLACEHOLDER_RE() *j2gort.Ext_java_util_regex_Pattern {
	clinit_JsMessageVisitor()
	return static_JsMessageVisitor_ICU_PLACEHOLDER_RE
}
func StaticSet_JsMessageVisitor_ICU_PLACEHOLDER_RE(v *j2gort.Ext_java_util_regex_Pattern) {
	clinit_JsMessageVisitor()
	static_JsMessageVisitor_ICU_PLACEHOLDER_RE = v
}

var static_JsMessageVisitor_MESSAGE_OPTION_NAMES j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_JsMessageVisitor_MESSAGE_OPTION_NAMES() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_JsMessageVisitor()
	return static_JsMessageVisitor_MESSAGE_OPTION_NAMES
}
func StaticSet_JsMessageVisitor_MESSAGE_OPTION_NAMES(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_JsMessageVisitor()
	static_JsMessageVisitor_MESSAGE_OPTION_NAMES = v
}

var static_JsMessageVisitor_ICU_TEMPLATE_OPTION_NAMES j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_JsMessageVisitor_ICU_TEMPLATE_OPTION_NAMES() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_JsMessageVisitor()
	return static_JsMessageVisitor_ICU_TEMPLATE_OPTION_NAMES
}
func StaticSet_JsMessageVisitor_ICU_TEMPLATE_OPTION_NAMES(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_JsMessageVisitor()
	static_JsMessageVisitor_ICU_TEMPLATE_OPTION_NAMES = v
}
func StaticGet_JsMessageVisitor_SCOPED_ALIASES_LITERAL_PREFIX() string { return "$jscomp$scope$" }

var clinitDone_JsMessageVisitor uint32
var clinitState_JsMessageVisitor j2gort.ClinitState

func clinit_JsMessageVisitor() {
	if atomic.LoadUint32(&clinitDone_JsMessageVisitor) != 0 {
		return
	}
	clinitSlow_JsMessageVisitor()
}

func clinitSlow_JsMessageVisitor() {
	if !j2gort.EnterClinit(&clinitState_JsMessageVisitor, &clinitDone_JsMessageVisitor) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_JsMessageVisitor, &clinitDone_JsMessageVisitor)
	static_JsMessageVisitor_FUNCTIONS_USED_IN_MESSAGE_EXTRACTION = j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.ImmListOfStr2("getMsg", "declareIcuTemplate"))
	static_JsMessageVisitor_MESSAGE_HAS_NO_DESCRIPTION = DiagnosticType_Warning("JSC_MSG_HAS_NO_DESCRIPTION", "Message {0} has no description. Add @desc JsDoc tag.")
	static_JsMessageVisitor_MESSAGE_HAS_NO_TEXT = DiagnosticType_Warning("JSC_MSG_HAS_NO_TEXT", "Message value of {0} is just an empty string. Empty messages are forbidden.")
	static_JsMessageVisitor_MESSAGE_TREE_MALFORMED = DiagnosticType_Error("JSC_MSG_TREE_MALFORMED", "Message parse tree malformed. {0}")
	static_JsMessageVisitor_MESSAGE_HAS_NO_VALUE = DiagnosticType_Error("JSC_MSG_HAS_NO_VALUE", "message node {0} has no value")
	static_JsMessageVisitor_MESSAGE_DUPLICATE_KEY = DiagnosticType_Error("JSC_MSG_KEY_DUPLICATED", "duplicate message variable name found for {0}, initial definition {1}:{2}")
	static_JsMessageVisitor_MESSAGE_NODE_IS_ORPHANED = DiagnosticType_Error("JSC_MSG_ORPHANED_NODE", "{0}() function may be used only with MSG_* property or variable")
	static_JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY = DiagnosticType_Warning("JSC_MSG_NOT_INITIALIZED_CORRECTLY", ("Message must be initialized using a call to " + "goog.getMsg" + " or " + "goog.i18n.messages.declareIcuTemplate"))
	static_JsMessageVisitor_BAD_FALLBACK_SYNTAX = DiagnosticType_Error("JSC_MSG_BAD_FALLBACK_SYNTAX", j2gort.Ext_java_lang_String_Format_2("Bad syntax. Expected syntax: %s(MSG_1, MSG_2)", []any{"goog.getMsgWithFallback"}))
	static_JsMessageVisitor_FALLBACK_ARG_ERROR = DiagnosticType_Error("JSC_MSG_FALLBACK_ARG_ERROR", "Could not find message entry for fallback argument {0}")
	static_JsMessageVisitor_SCOPED_ALIASES_PREFIX_PATTERN = j2gort.Ext_java_util_regex_Pattern_Compile_1("\\$jscomp\\$scope\\$\\S+\\$MSG_")
	static_JsMessageVisitor_MSG_UNNAMED_PATTERN = j2gort.Ext_java_util_regex_Pattern_Compile_1("MSG_UNNAMED.*")
	static_JsMessageVisitor_ICU_PLACEHOLDER_RE = j2gort.Ext_java_util_regex_Pattern_Compile_1("\\{([A-Z_0-9]+)\\}")
	static_JsMessageVisitor_MESSAGE_OPTION_NAMES = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_4("html", "unescapeHtmlEntities", "example", "original_code")
	static_JsMessageVisitor_ICU_TEMPLATE_OPTION_NAMES = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_5("description", "meaning", "alternate_message_id", "example", "original_code")
}

func (this_ *JsMessageVisitorBase) initFields_JsMessageVisitor() {
	_ = this_
	this_.messageNames = j2gort.NewOrderedMap[string, *JsMessageVisitor_MessageLocation]()
	this_.unnamedMessages = j2gort.NewCustomOrderedMap[*Var, *JsMessage]()
	this_.googMsgNodes = j2gort.NewOrderedSet[pkg_rhino.Node]()
}

func InitJsMessageVisitor(this_ *JsMessageVisitorBase, compiler AbstractCompiler, idGenerator JsMessage_IdGenerator) {
	_ = compiler
	_ = idGenerator
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_JsMessageVisitor()
	this_.compiler = compiler
	this_.idGenerator = idGenerator
}

// Java method: JsMessageVisitor.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *JsMessageVisitorBase) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, Upcast_NodeTraversal_Callback(this_.self_))
	{
		cur_1 := this_.googMsgNodes.Iter()
		for cur_1.HasNext() {
			msgNode := cur_1.
				Next()
			_ = msgNode
			this_.compiler.Report_1(JSError_Make_3(msgNode, StaticGet_JsMessageVisitor_MESSAGE_NODE_IS_ORPHANED(), []string{msgNode.GetFirstChild().GetQualifiedName()}))
		}
	}
}

// Java method: JsMessageVisitor.visit(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node unused)
func (this_ *JsMessageVisitorBase) Visit_3(traversal *NodeTraversal, node pkg_rhino.Node, unused pkg_rhino.Node) {
	_ = traversal
	_ = node
	_ = unused
	this_.priv_JsMessageVisitor_collectGetMsgCall_3688(traversal, node)
	this_.priv_JsMessageVisitor_checkMessageInitialization_3682(traversal, node)
}

// Java method: JsMessageVisitor.checkMessageInitialization(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node node)
func (this_ *JsMessageVisitorBase) priv_JsMessageVisitor_checkMessageInitialization_3682(traversal *NodeTraversal, node pkg_rhino.Node) {
	_ = traversal
	_ = node
	var parent pkg_rhino.Node = node.GetParent()
	_ = parent
	var originalMessageKey string = j2gort.NullString
	_ = originalMessageKey
	var possiblyObfuscatedMessageKey string = j2gort.NullString
	_ = possiblyObfuscatedMessageKey
	var msgNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = msgNode
	var jsDocInfo *pkg_rhino.JSDocInfo = nil
	_ = jsDocInfo
	switch node.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		if (parent == 0) || !(NodeUtil_IsNameDeclaration(parent)) {
			return
		}
		possiblyObfuscatedMessageKey = node.GetString()
		originalMessageKey = node.GetOriginalName()
		msgNode = node.GetFirstChild()
		jsDocInfo = parent.GetJSDocInfo()
	case int32(51) /* Token.ASSIGN */ :
		var getProp pkg_rhino.Node = node.GetFirstChild()
		_ = getProp
		if !(getProp.IsGetProp()) {
			return
		}
		possiblyObfuscatedMessageKey = getProp.GetString()
		originalMessageKey = getProp.GetOriginalName()
		msgNode = node.GetLastChild()
		jsDocInfo = node.GetJSDocInfo()
	case int32(100) /* Token.STRING_KEY */ :
		if (node.IsQuotedStringKey() || !(node.HasChildren())) || parent.IsObjectPattern() {
			return
		}
		if !(parent.IsObjectLit()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(parent))
		}
		possiblyObfuscatedMessageKey = node.GetString()
		originalMessageKey = node.GetOriginalName()
		msgNode = node.GetFirstChild()
		jsDocInfo = node.GetJSDocInfo()
	case int32(108) /* Token.MEMBER_FIELD_DEF */ :
		possiblyObfuscatedMessageKey = node.GetString()
		originalMessageKey = node.GetOriginalName()
		msgNode = node.GetFirstChild()
		jsDocInfo = node.GetJSDocInfo()
	default:
		return
	}
	var messageKeyFromLhs string
	_ = messageKeyFromLhs
	if originalMessageKey != j2gort.NullString {
		messageKeyFromLhs = originalMessageKey
	} else {
		messageKeyFromLhs = possiblyObfuscatedMessageKey
	}
	var msgNodeIsACall bool = ((msgNode != 0) && msgNode.IsCall())
	_ = msgNodeIsACall
	if !(this_.IsMessageName(messageKeyFromLhs)) {
		return
	}
	if msgNode == 0 {
		this_.compiler.Report_1(JSError_Make_3(node, StaticGet_JsMessageVisitor_MESSAGE_HAS_NO_VALUE(), []string{messageKeyFromLhs}))
		return
	}
	if JsMessageVisitor_IsLegalMessageVarAlias(msgNode) {
		return
	}
	if msgNodeIsACall {
		this_.googMsgNodes.Remove(msgNode)
	} else {
		this_.compiler.Report_1(JSError_Make_3(node, StaticGet_JsMessageVisitor_MESSAGE_NOT_INITIALIZED_CORRECTLY(), []string{}))
		return
	}
	var mapping *j2gort.Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping = this_.compiler.GetSourceMapping(traversal.GetSourceName(), node.GetLineno(), node.GetCharno())
	_ = mapping
	var sourceName string = j2gort.NullString
	_ = sourceName
	if mapping != nil {
		sourceName = (j2gort.StrValString(mapping.GetOriginalFile()) + ":" + j2gort.StrValInt32(mapping.GetLineNumber()))
	} else {
		sourceName = (j2gort.StrValString(traversal.GetSourceName()) + ":" + j2gort.StrValInt32(node.GetLineno()))
	}
	if !(msgNode.IsCall()) {
		j2gort.CheckState(false)
	}
	var fnNameNode pkg_rhino.Node = msgNode.GetFirstChild()
	_ = fnNameNode
	func() {
		defer func() {
			if rec_2 := recover(); rec_2 != nil {
				jex_3 := javalang.RecoverException(rec_2)
				_ = jex_3
				if j2gort.InstanceOf[*JsMessageVisitor_MalformedException](jex_3) {
					ex := j2gort.AsIface[*JsMessageVisitor_MalformedException](jex_3)
					_ = ex
					this_.compiler.Report_1(JSError_Make_3(ex.GetNode(), StaticGet_JsMessageVisitor_MESSAGE_TREE_MALFORMED(), []string{j2gort.ThrowableGetMessage(ex)}))
				} else {
					panic(rec_2)
				}
			}
		}()
		if this_.priv_JsMessageVisitor_isDeclareIcuTemplateCallee_3683(fnNameNode) {
			var icuTemplateDefinition IcuTemplateDefinition = this_.priv_JsMessageVisitor_extractIcuTemplateDefinition_3695(msgNode, jsDocInfo, messageKeyFromLhs, sourceName)
			_ = icuTemplateDefinition
			var extractedMessage *JsMessage = icuTemplateDefinition.GetMessage()
			_ = extractedMessage
			this_.priv_JsMessageVisitor_trackMessage_3685(traversal, possiblyObfuscatedMessageKey, msgNode, extractedMessage)
			this_.priv_JsMessageVisitor_reportErrorIfEmptyMessage_3686(node, extractedMessage)
			this_.self_.ProcessIcuTemplateDefinition(icuTemplateDefinition)
		} else if fnNameNode.MatchesQualifiedName_1_string("goog.getMsg") {
			var jsMessageDefinition JsMessageDefinition = this_.priv_JsMessageVisitor_extractJsMessageDefinition_3694(msgNode, jsDocInfo, messageKeyFromLhs, sourceName)
			_ = jsMessageDefinition
			var extractedMessage_2 *JsMessage = jsMessageDefinition.GetMessage()
			_ = extractedMessage_2
			this_.priv_JsMessageVisitor_trackMessage_3685(traversal, possiblyObfuscatedMessageKey, msgNode, extractedMessage_2)
			this_.priv_JsMessageVisitor_reportErrorIfEmptyMessage_3686(node, extractedMessage_2)
			var desc string = extractedMessage_2.GetDesc()
			_ = desc
			if ((desc == j2gort.NullString) || (j2gort.CheckNotNullString(javalang.StringTrim(j2gort.CheckNotNullString(desc))) == "")) && !(extractedMessage_2.IsExternal()) {
				this_.compiler.Report_1(JSError_Make_3(node, StaticGet_JsMessageVisitor_MESSAGE_HAS_NO_DESCRIPTION(), []string{extractedMessage_2.GetKey()}))
			}
			this_.self_.ProcessJsMessageDefinition(jsMessageDefinition)
		} else {
			this_.compiler.Report_1(JSError_Make_3(msgNode, StaticGet_JsMessageVisitor_MESSAGE_TREE_MALFORMED(), []string{("Message must be initialized using a call to " + "goog.getMsg" + " or " + "declareIcuTemplate" + " (from goog.i18n.messages).")}))
		}
	}()
}

// Java method: JsMessageVisitor.isDeclareIcuTemplateCallee(com.google.javascript.rhino.Node calleeNode)
func (this_ *JsMessageVisitorBase) priv_JsMessageVisitor_isDeclareIcuTemplateCallee_3683(calleeNode pkg_rhino.Node) bool {
	_ = calleeNode
	if calleeNode.IsName() {
		return javalang.StringEndsWith(j2gort.CheckNotNullString(calleeNode.GetString()), "declareIcuTemplate")
	}
	if calleeNode.IsGetProp() {
		return (javalang.StringEndsWith(j2gort.CheckNotNullString(calleeNode.GetString()), "declareIcuTemplate") && calleeNode.IsQualifiedName())
	}
	return false
}

// Java method: JsMessageVisitor.buildJsMessage(com.google.javascript.jscomp.JsMessageVisitor.InputForBuildJsMsg input)
func (this_ *JsMessageVisitorBase) priv_JsMessageVisitor_buildJsMessage_3684(input JsMessageVisitor_InputForBuildJsMsg) *JsMessage {
	_ = input
	var messageParts j2gort.Ext_com_google_common_collect_ImmutableList = input.GetMessageParts()
	_ = messageParts
	var messageKeyFromLhs string = input.GetMessageKeyFromLhs()
	_ = messageKeyFromLhs
	var externalMessageId string = JsMessageVisitor_GetExternalMessageId(messageKeyFromLhs)
	_ = externalMessageId
	var isAnonymous bool = false
	_ = isAnonymous
	var isExternal bool = false
	_ = isExternal
	var messageId string = j2gort.NullString
	_ = messageId
	var messageKeyFinal string = j2gort.NullString
	_ = messageKeyFinal
	if externalMessageId != j2gort.NullString {
		isAnonymous = false
		isExternal = true
		messageId = externalMessageId
		messageKeyFinal = messageKeyFromLhs
	} else {
		isExternal = false
		var meaningForIdGeneration string = input.GetMeaning()
		_ = meaningForIdGeneration
		if JsMessageVisitor_IsUnnamedMessageName(messageKeyFromLhs) {
			isAnonymous = true
			messageKeyFinal = this_.priv_JsMessageVisitor_generateKeyFromMessageText_3687(input.GetMessageText())
			if meaningForIdGeneration == j2gort.NullString {
				meaningForIdGeneration = messageKeyFinal
			}
		} else {
			isAnonymous = false
			messageKeyFinal = messageKeyFromLhs
			if meaningForIdGeneration == j2gort.NullString {
				meaningForIdGeneration = JsMessageVisitor_RemoveScopedAliasesPrefix(messageKeyFromLhs)
			}
		}
		if this_.idGenerator == nil || this_.idGenerator.IsJavaNil_() {
			messageId = meaningForIdGeneration
		} else {
			messageId = this_.idGenerator.GenerateId(meaningForIdGeneration, j2gort.AsIface[j2gort.Ext_java_util_List](messageParts))
		}
	}
	return NewJsMessage_Builder().AppendParts(j2gort.AsIface[j2gort.Ext_java_util_List](messageParts)).SetPlaceholderNameToExampleMap(j2gort.AsIface[j2gort.Ext_java_util_Map](input.GetPlaceholderExampleMap())).SetPlaceholderNameToOriginalCodeMap(j2gort.AsIface[j2gort.Ext_java_util_Map](input.GetPlaceholderOriginalCodeMap())).SetIsAnonymous(isAnonymous).SetIsExternalMsg(isExternal).SetKey(messageKeyFinal).SetSourceName(input.GetSourceName()).SetDesc(input.GetDescription()).SetMeaning(input.GetMeaning()).SetAlternateId(input.GetAlternateMessageId()).SetId(messageId).Build()
}

// Java method: JsMessageVisitor.trackMessage(com.google.javascript.jscomp.NodeTraversal traversal, java.lang.String possiblyObfuscatedMessageKey, com.google.javascript.rhino.Node msgNode, com.google.javascript.jscomp.JsMessage extractedMessage)
func (this_ *JsMessageVisitorBase) priv_JsMessageVisitor_trackMessage_3685(traversal *NodeTraversal, possiblyObfuscatedMessageKey string, msgNode pkg_rhino.Node, extractedMessage *JsMessage) {
	_ = traversal
	_ = possiblyObfuscatedMessageKey
	_ = msgNode
	_ = extractedMessage
	if !(extractedMessage.IsAnonymous()) && !(extractedMessage.IsExternal()) {
		this_.priv_JsMessageVisitor_checkIfMessageDuplicated_3693(extractedMessage.GetKey(), msgNode)
	}
	if extractedMessage.IsAnonymous() {
		this_.priv_JsMessageVisitor_trackUnnamedMessage_3690(traversal, extractedMessage, possiblyObfuscatedMessageKey)
	} else {
		this_.priv_JsMessageVisitor_trackNormalMessage_3689(extractedMessage, extractedMessage.GetKey(), msgNode)
	}
}

// Java method: JsMessageVisitor.reportErrorIfEmptyMessage(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.JsMessage extractedMessage)
func (this_ *JsMessageVisitorBase) priv_JsMessageVisitor_reportErrorIfEmptyMessage_3686(node pkg_rhino.Node, extractedMessage *JsMessage) {
	_ = node
	_ = extractedMessage
	if extractedMessage.IsEmpty() {
		this_.compiler.Report_1(JSError_Make_3(node, StaticGet_JsMessageVisitor_MESSAGE_HAS_NO_TEXT(), []string{extractedMessage.GetKey()}))
	}
}

// Java static method: JsMessageVisitor.getExternalMessageId(java.lang.String messageKey)
func JsMessageVisitor_GetExternalMessageId(messageKey string) string {
	clinit_JsMessageVisitor()
	_ = messageKey
	if javalang.StringStartsWith(j2gort.CheckNotNullString(messageKey), "MSG_EXTERNAL_") {
		var start int32 = j2gort.JStrLen("MSG_EXTERNAL_")
		_ = start
		var end int32 = start
		_ = end
		for ; end < j2gort.JStrLen(messageKey); end = end + 1 {
			var c uint16 = j2gort.JStrCharAt(messageKey, end)
			_ = c
			if (int32(c) > int32(uint16(57))) || (int32(c) < int32(uint16(48))) {
				break
			}
		}
		if end > start {
			return j2gort.JStrSub(messageKey, start, end)
		}
	}
	return j2gort.NullString
}

// Java method: JsMessageVisitor.generateKeyFromMessageText(java.lang.String msgText)
func (this_ *JsMessageVisitorBase) priv_JsMessageVisitor_generateKeyFromMessageText_3687(msgText string) string {
	_ = msgText
	var nonnegativeHash int64 = (int64(9223372036854775807) & JsMessage_Hash_Hash64_1(msgText))
	_ = nonnegativeHash
	return ("MSG_" + j2gort.StrValString(j2gort.Ext_com_google_common_base_Ascii_ToUpperCase_1_string(j2gort.BoxLong_ToString_2(nonnegativeHash, int32(36)))))
}

// Java method: JsMessageVisitor.collectGetMsgCall(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node call)
func (this_ *JsMessageVisitorBase) priv_JsMessageVisitor_collectGetMsgCall_3688(traversal *NodeTraversal, call pkg_rhino.Node) {
	_ = traversal
	_ = call
	if !(call.IsCall()) {
		return
	}
	var callee pkg_rhino.Node = call.GetFirstChild()
	_ = callee
	if callee.MatchesQualifiedName_1_string("goog.getMsg") || this_.priv_JsMessageVisitor_isDeclareIcuTemplateCallee_3683(callee) {
		this_.googMsgNodes.Add(call)
	} else if callee.MatchesQualifiedName_1_string("goog.getMsgWithFallback") {
		this_.priv_JsMessageVisitor_visitFallbackFunctionCall_3699(traversal, call)
	}
}

// Java method: JsMessageVisitor.trackNormalMessage(com.google.javascript.jscomp.JsMessage message, java.lang.String msgName, com.google.javascript.rhino.Node msgNode)
func (this_ *JsMessageVisitorBase) priv_JsMessageVisitor_trackNormalMessage_3689(message *JsMessage, msgName string, msgNode pkg_rhino.Node) {
	_ = message
	_ = msgName
	_ = msgNode
	var val_4 JsMessageVisitor_MessageLocation
	InitJsMessageVisitor_MessageLocation(&val_4, message, msgNode)
	var location *JsMessageVisitor_MessageLocation = &val_4
	_ = location
	this_.messageNames.Put(msgName, location)
}

// Java method: JsMessageVisitor.trackUnnamedMessage(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.JsMessage message, java.lang.String msgNameInScope)
func (this_ *JsMessageVisitorBase) priv_JsMessageVisitor_trackUnnamedMessage_3690(t *NodeTraversal, message *JsMessage, msgNameInScope string) {
	_ = t
	_ = message
	_ = msgNameInScope
	var var_ *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(msgNameInScope))
	_ = var_
	if var_ != nil {
		this_.unnamedMessages.Put(var_, message)
	}
}

// Java static method: JsMessageVisitor.isLegalMessageVarAlias(com.google.javascript.rhino.Node msgNode)
func JsMessageVisitor_IsLegalMessageVarAlias(msgNode pkg_rhino.Node) bool {
	clinit_JsMessageVisitor()
	_ = msgNode
	if (msgNode.IsGetProp() && msgNode.IsQualifiedName()) && javalang.StringStartsWith(j2gort.CheckNotNullString(msgNode.GetString()), "MSG_") {
		return true
	}
	if !(msgNode.IsName()) {
		return false
	}
	var originalName string
	_ = originalName
	if msgNode.GetOriginalName() != j2gort.NullString {
		originalName = msgNode.GetOriginalName()
	} else {
		originalName = msgNode.GetString()
	}
	if javalang.StringStartsWith(j2gort.CheckNotNullString(originalName), "MSG_") {
		return true
	}
	return false
}

// Java method: JsMessageVisitor.getTrackedUnnamedMessage(com.google.javascript.jscomp.NodeTraversal t, java.lang.String msgNameInScope)
func (this_ *JsMessageVisitorBase) priv_JsMessageVisitor_getTrackedUnnamedMessage_3691(t *NodeTraversal, msgNameInScope string) *JsMessage {
	_ = t
	_ = msgNameInScope
	var var_ *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(msgNameInScope))
	_ = var_
	if var_ != nil {
		return this_.unnamedMessages.Get(var_)
	}
	return nil
}

// Java method: JsMessageVisitor.getTrackedNormalMessage(java.lang.String msgName)
func (this_ *JsMessageVisitorBase) priv_JsMessageVisitor_getTrackedNormalMessage_3692(msgName string) *JsMessage {
	_ = msgName
	var location *JsMessageVisitor_MessageLocation = this_.messageNames.Get(msgName)
	_ = location
	if location == nil {
		return nil
	} else {
		return location.message
	}
}

// Java method: JsMessageVisitor.checkIfMessageDuplicated(java.lang.String msgName, com.google.javascript.rhino.Node msgNode)
func (this_ *JsMessageVisitorBase) priv_JsMessageVisitor_checkIfMessageDuplicated_3693(msgName string, msgNode pkg_rhino.Node) {
	_ = msgName
	_ = msgNode
	if this_.messageNames.ContainsKey(msgName) {
		var location *JsMessageVisitor_MessageLocation = this_.messageNames.Get(msgName)
		_ = location
		this_.compiler.Report_1(JSError_Make_3(msgNode, StaticGet_JsMessageVisitor_MESSAGE_DUPLICATE_KEY(), []string{msgName, location.messageNode.GetSourceFileName(), strconv.FormatInt(int64(location.messageNode.GetLineno()), 10)}))
	}
}

// Java static method: JsMessageVisitor.extractStringFromStringExprNode(com.google.javascript.rhino.Node node)
func JsMessageVisitor_ExtractStringFromStringExprNode(node pkg_rhino.Node) string {
	clinit_JsMessageVisitor()
	_ = node
	switch node.TokenSwitchOrd_() {
	case int32(35) /* Token.STRINGLIT */ :
		return node.GetString()
	case int32(131) /* Token.TEMPLATELIT */ :
		if node.HasOneChild() {
			return j2gort.CheckNotNullString(node.GetFirstChild().GetCookedString())
		} else {
			panic(NewJsMessageVisitor_MalformedException("Template literals with substitutions are not allowed.", node))
		}
	case int32(13) /* Token.ADD */ :
		var sb j2gort.LocalStringBuilder

		{
			var child pkg_rhino.Node = node.GetFirstChild()
			_ = child
			for ; child != 0; child = child.GetNext() {
				sb.AppendStr(JsMessageVisitor_ExtractStringFromStringExprNode(child))
			}
		}
		return sb.String()
	default:
		panic(NewJsMessageVisitor_MalformedException("literal string or concatenation expected", node))
	}
	panic("unreachable")
}

// Java method: JsMessageVisitor.extractJsMessageDefinition(com.google.javascript.rhino.Node msgNode, com.google.javascript.rhino.JSDocInfo jsDocInfo, java.lang.String messageKeyFromLhs, java.lang.String sourceName)
func (this_ *JsMessageVisitorBase) priv_JsMessageVisitor_extractJsMessageDefinition_3694(msgNode pkg_rhino.Node, jsDocInfo *pkg_rhino.JSDocInfo, messageKeyFromLhs string, sourceName string) JsMessageDefinition {
	_ = msgNode
	_ = jsDocInfo
	_ = messageKeyFromLhs
	_ = sourceName
	var description string = j2gort.NullString
	_ = description
	var meaning string = j2gort.NullString
	_ = meaning
	var alternateMessageId string = j2gort.NullString
	_ = alternateMessageId
	if jsDocInfo == nil {
		description = j2gort.NullString
		meaning = j2gort.NullString
		alternateMessageId = j2gort.NullString
	} else {
		description = jsDocInfo.GetDescription()
		meaning = jsDocInfo.GetMeaning()
		alternateMessageId = jsDocInfo.GetAlternateMessageId()
	}
	var msgTextNode pkg_rhino.Node = msgNode.GetSecondChild()
	_ = msgTextNode
	if msgTextNode == 0 {
		panic(NewJsMessageVisitor_MalformedException("Message string literal expected", msgNode))
	}
	var valuesObjLit pkg_rhino.Node = msgTextNode.GetNext()
	_ = valuesObjLit
	var optionsBagArgument pkg_rhino.Node
	_ = optionsBagArgument
	if valuesObjLit == 0 {
		optionsBagArgument = pkg_rhino.Node(0)
	} else {
		optionsBagArgument = valuesObjLit.GetNext()
	}
	var unexpectedArgument pkg_rhino.Node
	_ = unexpectedArgument
	if optionsBagArgument == 0 {
		unexpectedArgument = pkg_rhino.Node(0)
	} else {
		unexpectedArgument = optionsBagArgument.GetNext()
	}
	if unexpectedArgument != 0 {
		panic(NewJsMessageVisitor_MalformedException("too many arguments", unexpectedArgument))
	}
	var msgTextString string = JsMessageVisitor_ExtractStringFromStringExprNode(msgTextNode)
	_ = msgTextString
	var googGetMsgParsedText JsMessageVisitor_GoogGetMsgParsedText = nil
	_ = googGetMsgParsedText
	func() {
		defer func() {
			if rec_5 := recover(); rec_5 != nil {
				jex_6 := javalang.RecoverException(rec_5)
				_ = jex_6
				if j2gort.InstanceOf[*JsMessage_PlaceholderFormatException](jex_6) {
					e := j2gort.AsIface[*JsMessage_PlaceholderFormatException](jex_6)
					_ = e
					panic(NewJsMessageVisitor_MalformedException(j2gort.ThrowableGetMessage(e), msgTextNode))
				} else {
					panic(rec_5)
				}
			}
		}()
		googGetMsgParsedText = JsMessageVisitor_ExtractGoogGetMsgParsedText(msgTextString)
	}()
	var placeholderObjectLiteralMap JsMessageVisitor_ObjectLiteralMap = JsMessageVisitor_ExtractObjectLiteralMap(valuesObjLit)
	_ = placeholderObjectLiteralMap
	var placeholderNamesFromText j2gort.Ext_com_google_common_collect_ImmutableSet = googGetMsgParsedText.GetPlaceholderNames()
	_ = placeholderNamesFromText
	placeholderObjectLiteralMap.CheckForRequiredKeys(j2gort.AsIface[j2gort.Ext_java_util_Set](placeholderNamesFromText), j2gort.NewLambda_Ext_java_util_function_Function(func(lp_7 any) any {
		placeholderName := j2gort.CastString(lp_7)
		_ = placeholderName
		return NewJsMessageVisitor_MalformedException(("Unrecognized message placeholder referenced: " + j2gort.StrValString(placeholderName)), msgTextNode)
	}))
	placeholderObjectLiteralMap.CheckForUnexpectedKeys(j2gort.AsIface[j2gort.Ext_java_util_Set](placeholderNamesFromText), lambdaSingleton_js_message_visitor_1)
	var placeholderValuesMap j2gort.Ext_com_google_common_collect_ImmutableMap = placeholderObjectLiteralMap.ExtractAsValueMap()
	_ = placeholderValuesMap
	var jsMessageOptions JsMessageVisitor_JsMessageOptions = this_.priv_JsMessageVisitor_extractJsMessageOptions_3697(optionsBagArgument)
	_ = jsMessageOptions
	jsMessageOptions.CheckForUnknownPlaceholders(j2gort.AsIface[j2gort.Ext_java_util_Set](placeholderNamesFromText))
	var jsExtractedMessage *JsMessage = this_.priv_JsMessageVisitor_buildJsMessage_3684(NewJsMessageVisitor__Anon_56(this_.self_, messageKeyFromLhs, sourceName, googGetMsgParsedText, jsMessageOptions, description, meaning, alternateMessageId))
	_ = jsExtractedMessage
	return NewJsMessageVisitor__Anon_57(this_.self_, jsExtractedMessage, msgNode, msgTextNode, valuesObjLit, placeholderValuesMap, jsMessageOptions)
}

// Java method: JsMessageVisitor.extractIcuTemplateDefinition(com.google.javascript.rhino.Node msgNode, com.google.javascript.rhino.JSDocInfo jsDocInfo, java.lang.String messageKeyFromLhs, java.lang.String sourceName)
func (this_ *JsMessageVisitorBase) priv_JsMessageVisitor_extractIcuTemplateDefinition_3695(msgNode pkg_rhino.Node, jsDocInfo *pkg_rhino.JSDocInfo, messageKeyFromLhs string, sourceName string) IcuTemplateDefinition {
	_ = msgNode
	_ = jsDocInfo
	_ = messageKeyFromLhs
	_ = sourceName
	if jsDocInfo != nil {
		if jsDocInfo.GetAlternateMessageId() != j2gort.NullString {
			panic(NewJsMessageVisitor_MalformedException("Use the 'alternate_message_id' option, not the '@alternateMessageId' annotation", msgNode))
		}
		if jsDocInfo.GetDescription() != j2gort.NullString {
			panic(NewJsMessageVisitor_MalformedException("Use the 'description' option, not the '@desc' annotation", msgNode))
		}
		if jsDocInfo.GetMeaning() != j2gort.NullString {
			panic(NewJsMessageVisitor_MalformedException("Use the 'meaning' option, not the '@meaning' annotation", msgNode))
		}
	}
	var stringLiteralExpression pkg_rhino.Node = msgNode.GetSecondChild()
	_ = stringLiteralExpression
	if stringLiteralExpression == 0 {
		panic(NewJsMessageVisitor_MalformedException("message string argument expected", msgNode))
	}
	var icuMessageTemplateString *JsMessageVisitor_IcuMessageTemplateString = this_.priv_JsMessageVisitor_extractIcuMessageTemplateString_3696(stringLiteralExpression)
	_ = icuMessageTemplateString
	var optionsBagNode pkg_rhino.Node = stringLiteralExpression.GetNext()
	_ = optionsBagNode
	if optionsBagNode == 0 {
		panic(NewJsMessageVisitor_MalformedException("options argument expected", msgNode))
	}
	var icuTemplateOptions JsMessageVisitor_IcuTemplateOptions = this_.priv_JsMessageVisitor_extractIcuTemplateOptions_3698(optionsBagNode)
	_ = icuTemplateOptions
	var extractedIcuTemplateParts *JsMessageVisitor_ExtractedIcuTemplateParts = icuMessageTemplateString.ExtractParts(j2gort.AsIface[j2gort.Ext_java_util_Set](icuTemplateOptions.GetPlaceholderNames()))
	_ = extractedIcuTemplateParts
	icuTemplateOptions.CheckForUnknownPlaceholders(j2gort.AsIface[j2gort.Ext_java_util_Set](extractedIcuTemplateParts.extractedPlaceholderNames))
	var unexpectedArgument pkg_rhino.Node = optionsBagNode.GetNext()
	_ = unexpectedArgument
	if unexpectedArgument != 0 {
		panic(NewJsMessageVisitor_MalformedException("too many arguments", unexpectedArgument))
	}
	var extractedMessage *JsMessage = this_.priv_JsMessageVisitor_buildJsMessage_3684(NewJsMessageVisitor__Anon_58(this_.self_, messageKeyFromLhs, sourceName, icuMessageTemplateString, extractedIcuTemplateParts, icuTemplateOptions))
	_ = extractedMessage
	return NewJsMessageVisitor__Anon_59(this_.self_, extractedMessage, msgNode, stringLiteralExpression)
}

// Java method: JsMessageVisitor.extractIcuMessageTemplateString(com.google.javascript.rhino.Node stringExpression)
func (this_ *JsMessageVisitorBase) priv_JsMessageVisitor_extractIcuMessageTemplateString_3696(stringExpression pkg_rhino.Node) *JsMessageVisitor_IcuMessageTemplateString {
	_ = stringExpression
	var templateString string = JsMessageVisitor_ExtractStringFromStringExprNode(stringExpression)
	_ = templateString
	return NewJsMessageVisitor_IcuMessageTemplateString(templateString)
}

// Java method: JsMessageVisitor.extractJsMessageOptions(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node optionsBag)
func (this_ *JsMessageVisitorBase) priv_JsMessageVisitor_extractJsMessageOptions_3697(optionsBag pkg_rhino.Node) JsMessageVisitor_JsMessageOptions {
	_ = optionsBag
	var objectLiteralMap JsMessageVisitor_ObjectLiteralMap = JsMessageVisitor_ExtractObjectLiteralMap(optionsBag)
	_ = objectLiteralMap
	objectLiteralMap.CheckForUnexpectedKeys(j2gort.AsIface[j2gort.Ext_java_util_Set](StaticGet_JsMessageVisitor_MESSAGE_OPTION_NAMES()), lambdaSingleton_js_message_visitor_2)
	var isEscapeLessThan bool = objectLiteralMap.GetBooleanValueOrFalse("html")
	_ = isEscapeLessThan
	var isUnescapeHtmlEntities bool = objectLiteralMap.GetBooleanValueOrFalse("unescapeHtmlEntities")
	_ = isUnescapeHtmlEntities
	var exampleValueNode pkg_rhino.Node = objectLiteralMap.GetValueNode("example")
	_ = exampleValueNode
	var exampleObjectLiteralMap JsMessageVisitor_ObjectLiteralMap = JsMessageVisitor_ExtractObjectLiteralMap(exampleValueNode)
	_ = exampleObjectLiteralMap
	var placeholderExamplesMap j2gort.Ext_com_google_common_collect_ImmutableMap = exampleObjectLiteralMap.ExtractAsStringToStringMap()
	_ = placeholderExamplesMap
	var originalCodeValueNode pkg_rhino.Node = objectLiteralMap.GetValueNode("original_code")
	_ = originalCodeValueNode
	var originalCodeObjectLiteralMap JsMessageVisitor_ObjectLiteralMap = JsMessageVisitor_ExtractObjectLiteralMap(originalCodeValueNode)
	_ = originalCodeObjectLiteralMap
	var placeholderOriginalCodeMap j2gort.Ext_com_google_common_collect_ImmutableMap = originalCodeObjectLiteralMap.ExtractAsStringToStringMap()
	_ = placeholderOriginalCodeMap
	return NewJsMessageVisitor__Anon_60(this_.self_, isEscapeLessThan, isUnescapeHtmlEntities, placeholderExamplesMap, placeholderOriginalCodeMap, exampleObjectLiteralMap, originalCodeObjectLiteralMap)
}

// Java method: JsMessageVisitor.extractIcuTemplateOptions(com.google.javascript.rhino.Node optionsBag)
func (this_ *JsMessageVisitorBase) priv_JsMessageVisitor_extractIcuTemplateOptions_3698(optionsBag pkg_rhino.Node) JsMessageVisitor_IcuTemplateOptions {
	_ = optionsBag
	var objectLiteralMap JsMessageVisitor_ObjectLiteralMap = JsMessageVisitor_ExtractObjectLiteralMap(optionsBag)
	_ = objectLiteralMap
	objectLiteralMap.CheckForUnexpectedKeys(j2gort.AsIface[j2gort.Ext_java_util_Set](StaticGet_JsMessageVisitor_ICU_TEMPLATE_OPTION_NAMES()), lambdaSingleton_js_message_visitor_3)
	var descriptionNode pkg_rhino.Node = objectLiteralMap.GetValueNode("description")
	_ = descriptionNode
	if descriptionNode == 0 {
		panic(NewJsMessageVisitor_MalformedException("'description' option field is missing", optionsBag))
	}
	var description string = JsMessageVisitor_ExtractStringFromStringExprNode(descriptionNode)
	_ = description
	var meaningNode pkg_rhino.Node = objectLiteralMap.GetValueNode("meaning")
	_ = meaningNode
	var meaning string
	_ = meaning
	if meaningNode == 0 {
		meaning = j2gort.NullString
	} else {
		meaning = JsMessageVisitor_ExtractStringFromStringExprNode(meaningNode)
	}
	var alternateMessageIdNode pkg_rhino.Node = objectLiteralMap.GetValueNode("alternate_message_id")
	_ = alternateMessageIdNode
	var alternateMessageId string
	_ = alternateMessageId
	if alternateMessageIdNode == 0 {
		alternateMessageId = j2gort.NullString
	} else {
		alternateMessageId = JsMessageVisitor_ExtractStringFromStringExprNode(alternateMessageIdNode)
	}
	var exampleValueNode pkg_rhino.Node = objectLiteralMap.GetValueNode("example")
	_ = exampleValueNode
	var exampleObjectLiteralMap JsMessageVisitor_ObjectLiteralMap = JsMessageVisitor_ExtractObjectLiteralMap(exampleValueNode)
	_ = exampleObjectLiteralMap
	var placeholderExamplesMap j2gort.Ext_com_google_common_collect_ImmutableMap = exampleObjectLiteralMap.ExtractAsStringToStringMap()
	_ = placeholderExamplesMap
	var originalCodeValueNode pkg_rhino.Node = objectLiteralMap.GetValueNode("original_code")
	_ = originalCodeValueNode
	var originalCodeObjectLiteralMap JsMessageVisitor_ObjectLiteralMap = JsMessageVisitor_ExtractObjectLiteralMap(originalCodeValueNode)
	_ = originalCodeObjectLiteralMap
	var placeholderOriginalCodeMap j2gort.Ext_com_google_common_collect_ImmutableMap = originalCodeObjectLiteralMap.ExtractAsStringToStringMap()
	_ = placeholderOriginalCodeMap
	var placeholderNames j2gort.Ext_com_google_common_collect_ImmutableSet = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet_Builder](j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder().AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](placeholderExamplesMap.KeySet_0()))).AddAll_1_Ext_java_lang_Iterable(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](placeholderOriginalCodeMap.KeySet_0()))).Build())
	_ = placeholderNames
	{
		cur_11 := j2gort.NewIterCursor(j2gort.CheckNotNull(placeholderNames))
		for cur_11.HasNext() {
			placeholderName := cur_11.NextStr()
			_ = placeholderName
			if !(JsMessage_IsCanonicalPlaceholderNameFormat(placeholderName)) {
				panic(NewJsMessageVisitor_MalformedException(j2gort.Ext_java_lang_String_Format_2("Placeholder not in UPPER_SNAKE_CASE: %s", []any{j2gort.BoxString(placeholderName)}), optionsBag))
			}
		}
	}
	return NewJsMessageVisitor__Anon_61(this_.self_, description, meaning, alternateMessageId, placeholderExamplesMap, placeholderOriginalCodeMap, placeholderNames, exampleObjectLiteralMap, originalCodeObjectLiteralMap)
}

// Java static method: JsMessageVisitor.extractBooleanStringKeyValue(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node stringKeyNode)
func JsMessageVisitor_ExtractBooleanStringKeyValue(stringKeyNode pkg_rhino.Node) bool {
	clinit_JsMessageVisitor()
	_ = stringKeyNode
	if stringKeyNode == 0 {
		return false
	} else {
		var valueNode pkg_rhino.Node = stringKeyNode.GetOnlyChild()
		_ = valueNode
		if valueNode.IsTrue() {
			return true
		} else if valueNode.IsFalse() {
			return false
		} else {
			panic(NewJsMessageVisitor_MalformedException((j2gort.StrValString(stringKeyNode.GetString()) + ": Literal true or false expected"), valueNode))
		}
	}
}

// Java static method: JsMessageVisitor.extractObjectLiteralMap(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node objLit)
func JsMessageVisitor_ExtractObjectLiteralMap(objLit pkg_rhino.Node) JsMessageVisitor_ObjectLiteralMap {
	clinit_JsMessageVisitor()
	_ = objLit
	if objLit == 0 {
		return NewJsMessageVisitor_ObjectLiteralMapImpl(j2gort.AsIface[j2gort.Ext_java_util_Map](j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()))
	}
	if !(objLit.IsObjectLit()) {
		panic(NewJsMessageVisitor_MalformedException("object literal expected", objLit))
	}
	var stringToStringKeyMap *j2gort.Ext_java_util_LinkedHashMap = j2gort.NewExt_java_util_LinkedHashMap()
	_ = stringToStringKeyMap
	{
		var stringKey pkg_rhino.Node = objLit.GetFirstChild()
		_ = stringKey
		for ; stringKey != 0; stringKey = stringKey.GetNext() {
			if !(stringKey.IsStringKey()) {
				panic(NewJsMessageVisitor_MalformedException("string key expected", stringKey))
			}
			var key string = stringKey.GetString()
			_ = key
			if j2gort.MapContainsKeyStr(stringToStringKeyMap, key) {
				panic(NewJsMessageVisitor_MalformedException(("duplicate string key: " + j2gort.StrValString(key)), stringKey))
			}
			j2gort.MapPutStr(stringToStringKeyMap, stringKey.GetString(), j2gort.BoxNode(stringKey))
		}
	}
	return NewJsMessageVisitor_ObjectLiteralMapImpl(stringToStringKeyMap)
}

// Java static method: JsMessageVisitor.parseJsMessageTextIntoParts(java.lang.String originalMsgText)
func JsMessageVisitor_ParseJsMessageTextIntoParts(originalMsgText string) j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_JsMessageVisitor()
	_ = originalMsgText
	var googGetMsgParsedText JsMessageVisitor_GoogGetMsgParsedText = JsMessageVisitor_ExtractGoogGetMsgParsedText(originalMsgText)
	_ = googGetMsgParsedText
	return googGetMsgParsedText.GetParts()
}

// Java static method: JsMessageVisitor.extractGoogGetMsgParsedText(java.lang.String originalMsgText)
func JsMessageVisitor_ExtractGoogGetMsgParsedText(originalMsgText string) JsMessageVisitor_GoogGetMsgParsedText {
	clinit_JsMessageVisitor()
	_ = originalMsgText
	var msgText string = originalMsgText
	_ = msgText
	var partsBuilder *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
	_ = partsBuilder
	var placeholderNamesBuilder j2gort.Ext_com_google_common_collect_ImmutableSet_Builder = j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder()
	_ = placeholderNamesBuilder
	for {
		var phBegin int32 = j2gort.JStrIndexOf(msgText, "{$")
		_ = phBegin
		if phBegin < int32(0) {
			partsBuilder.AddVoid_(j2gort.PtrToAny(JsMessage_StringPart_Create(msgText)))
			break
		} else {
			if phBegin > int32(0) {
				partsBuilder.AddVoid_(j2gort.PtrToAny(JsMessage_StringPart_Create(j2gort.JStrSub(msgText, int32(0), phBegin))))
			}
			var phEnd int32 = j2gort.JStrIndexOfFrom(msgText, "}", phBegin)
			_ = phEnd
			if phEnd < int32(0) {
				panic(NewJsMessage_PlaceholderFormatException("Placeholder incorrectly formatted"))
			}
			var phName string = j2gort.JStrSub(msgText, (phBegin + j2gort.JStrLen("{$")), phEnd)
			_ = phName
			if !(JsMessage_IsLowerCamelCaseWithNumericSuffixes(phName)) {
				panic(NewJsMessage_PlaceholderFormatException(("Placeholder name not in lowerCamelCase: " + j2gort.StrValString(phName))))
			}
			j2gort.ImmSetBuilderAddStrVoid(placeholderNamesBuilder, phName)
			partsBuilder.AddVoid_(j2gort.PtrToAny(JsMessage_PlaceholderReference_CreateForJsName(phName)))
			var nextPos int32 = (phEnd + j2gort.JStrLen("}"))
			_ = nextPos
			if nextPos < j2gort.JStrLen(msgText) {
				msgText = j2gort.JStrSubFrom(msgText, nextPos)
			} else {
				break
			}
		}
	}
	var placeholderNames j2gort.Ext_com_google_common_collect_ImmutableSet = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](placeholderNamesBuilder.Build())
	_ = placeholderNames
	var parts j2gort.Ext_com_google_common_collect_ImmutableList = partsBuilder.BuildImmList_()
	_ = parts
	return NewJsMessageVisitor__Anon_62(originalMsgText, placeholderNames, parts)
}

// Java method: JsMessageVisitor.visitFallbackFunctionCall(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node call)
func (this_ *JsMessageVisitorBase) priv_JsMessageVisitor_visitFallbackFunctionCall_3699(t *NodeTraversal, call pkg_rhino.Node) {
	_ = t
	_ = call
	if (!(call.HasXChildren(int32(3))) || !(JsMessageVisitor_IsMessageIdentifier(call.GetSecondChild()))) || !(JsMessageVisitor_IsMessageIdentifier(call.GetLastChild())) {
		this_.compiler.Report_1(JSError_Make_3(call, StaticGet_JsMessageVisitor_BAD_FALLBACK_SYNTAX(), []string{}))
		return
	}
	var firstArg pkg_rhino.Node = call.GetSecondChild()
	_ = firstArg
	var firstMessage *JsMessage = this_.priv_JsMessageVisitor_getJsMessageFromNode_3702(t, firstArg)
	_ = firstMessage
	if firstMessage == nil {
		this_.compiler.Report_1(JSError_Make_3(firstArg, StaticGet_JsMessageVisitor_FALLBACK_ARG_ERROR(), []string{firstArg.GetQualifiedName()}))
		return
	}
	var secondArg pkg_rhino.Node = firstArg.GetNext()
	_ = secondArg
	var secondMessage *JsMessage = this_.priv_JsMessageVisitor_getJsMessageFromNode_3702(t, secondArg)
	_ = secondMessage
	if secondMessage == nil {
		this_.compiler.Report_1(JSError_Make_3(secondArg, StaticGet_JsMessageVisitor_FALLBACK_ARG_ERROR(), []string{secondArg.GetQualifiedName()}))
		return
	}
	this_.self_.ProcessMessageFallback(call, firstMessage, secondMessage)
}

// Java method: JsMessageVisitor.processJsMessageDefinition(com.google.javascript.jscomp.JsMessageDefinition definition)
func (this_ *JsMessageVisitorBase) ProcessJsMessageDefinition(definition JsMessageDefinition) {
	_ = definition
	panic("abstract or native method invoked")
}

// Java method: JsMessageVisitor.processIcuTemplateDefinition(com.google.javascript.jscomp.IcuTemplateDefinition definition)
func (this_ *JsMessageVisitorBase) ProcessIcuTemplateDefinition(definition IcuTemplateDefinition) {
	_ = definition
	panic("abstract or native method invoked")
}

// Java method: JsMessageVisitor.processMessageFallback(com.google.javascript.rhino.Node callNode, com.google.javascript.jscomp.JsMessage message1, com.google.javascript.jscomp.JsMessage message2)
func (this_ *JsMessageVisitorBase) ProcessMessageFallback(callNode pkg_rhino.Node, message1 *JsMessage, message2 *JsMessage) {
	_ = callNode
	_ = message1
	_ = message2
}

// Java method: JsMessageVisitor.isMessageName(java.lang.String identifier)
func (this_ *JsMessageVisitorBase) IsMessageName(identifier string) bool {
	_ = identifier
	return (javalang.StringStartsWith(j2gort.CheckNotNullString(identifier), "MSG_") || JsMessageVisitor_IsScopedAliasesPrefix(identifier))
}

// Java static method: JsMessageVisitor.isMessageIdentifier(com.google.javascript.rhino.Node node)
func JsMessageVisitor_IsMessageIdentifier(node pkg_rhino.Node) bool {
	clinit_JsMessageVisitor()
	_ = node
	if !(node.IsQualifiedName()) {
		return false
	}
	{
		var cur pkg_rhino.Node = node
		_ = cur
		for ; cur != 0; cur = func() pkg_rhino.Node {
			if cur.IsGetProp() {
				return cur.GetFirstChild()
			}
			return pkg_rhino.Node(0)
		}() {
			if (cur.IsGetProp() || cur.IsName()) && (j2gort.JStrIndexOf(cur.GetString(), "MSG_") >= 0) {
				return true
			}
		}
	}
	return false
}

// Java method: JsMessageVisitor.getJsMessageFromNode(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node node)
func (this_ *JsMessageVisitorBase) priv_JsMessageVisitor_getJsMessageFromNode_3702(t *NodeTraversal, node pkg_rhino.Node) *JsMessage {
	_ = t
	_ = node
	var messageName string = node.GetQualifiedName()
	_ = messageName
	if (messageName == j2gort.NullString) || !(j2gort.JStrIndexOf(messageName, "MSG_") >= 0) {
		return nil
	}
	var messageKey string = j2gort.JStrSubFrom(messageName, j2gort.JStrIndexOf(messageName, "MSG_"))
	_ = messageKey
	if JsMessageVisitor_IsUnnamedMessageName(messageKey) {
		return this_.priv_JsMessageVisitor_getTrackedUnnamedMessage_3691(t, messageName)
	} else {
		return this_.priv_JsMessageVisitor_getTrackedNormalMessage_3692(messageKey)
	}
}

// Java static method: JsMessageVisitor.isUnnamedMessageName(java.lang.String identifier)
func JsMessageVisitor_IsUnnamedMessageName(identifier string) bool {
	clinit_JsMessageVisitor()
	_ = identifier
	return StaticGet_JsMessageVisitor_MSG_UNNAMED_PATTERN().Matcher(j2gort.BoxString(identifier)).Matches()
}

// Java method: JsMessageVisitor.checkNode(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node node, com.google.javascript.rhino.Token type)
func (this_ *JsMessageVisitorBase) CheckNode(node pkg_rhino.Node, type_ pkg_rhino.Token) {
	_ = node
	_ = type_
	if node == 0 {
		panic(NewJsMessageVisitor_MalformedException(("Expected node type " + j2gort.StrVal(type_) + "; found: null"), node))
	}
	if node.GetToken() != type_ {
		panic(NewJsMessageVisitor_MalformedException(("Expected node type " + j2gort.StrVal(type_) + "; found: " + j2gort.StrVal(node.GetToken())), node))
	}
}

// Java static method: JsMessageVisitor.isScopedAliasesPrefix(java.lang.String name)
func JsMessageVisitor_IsScopedAliasesPrefix(name string) bool {
	clinit_JsMessageVisitor()
	_ = name
	if (j2gort.JStrLen(name) < int32(19)) || !(javalang.StringStartsWith(j2gort.CheckNotNullString(name), "$jscomp$scope$")) {
		return false
	}
	return StaticGet_JsMessageVisitor_SCOPED_ALIASES_PREFIX_PATTERN().Matcher(j2gort.BoxString(name)).LookingAt()
}

// Java static method: JsMessageVisitor.removeScopedAliasesPrefix(java.lang.String name)
func JsMessageVisitor_RemoveScopedAliasesPrefix(name string) string {
	clinit_JsMessageVisitor()
	_ = name
	return StaticGet_JsMessageVisitor_SCOPED_ALIASES_PREFIX_PATTERN().Matcher(j2gort.BoxString(name)).ReplaceFirst_1_string("MSG_")
}

// Java interface: com.google.javascript.jscomp.JsMessageVisitor.InputForBuildJsMsg
type JsMessageVisitor_InputForBuildJsMsg interface {
	Implements_jscomp_JsMessageVisitor_InputForBuildJsMsg()
	IsJavaNil_() bool
	GetMessageKeyFromLhs() string
	GetSourceName() string
	GetMessageText() string
	GetMessageParts() j2gort.Ext_com_google_common_collect_ImmutableList
	GetPlaceholderExampleMap() any
	GetPlaceholderOriginalCodeMap() any
	GetDescription() string
	GetMeaning() string
	GetAlternateMessageId() string
}

func ToIface_JsMessageVisitor_InputForBuildJsMsg[T interface {
	JsMessageVisitor_InputForBuildJsMsg
	comparable
}](v T) JsMessageVisitor_InputForBuildJsMsg {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_JsMessageVisitor_InputForBuildJsMsg[T JsMessageVisitor_InputForBuildJsMsg](v T) JsMessageVisitor_InputForBuildJsMsg {
	var vi JsMessageVisitor_InputForBuildJsMsg = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JsMessageVisitor.InputForBuildJsMsg", "com.google.javascript.jscomp.JsMessageVisitor$InputForBuildJsMsg", func(v_ any) bool {
		return j2gort.InstanceOf[JsMessageVisitor_InputForBuildJsMsg](v_)
	})
}

// Java class: com.google.javascript.jscomp.JsMessageVisitor.IcuMessageTemplateString
type JsMessageVisitor_IcuMessageTemplateString struct {
	template string `gson:"template"`
	idHash_  atomic.Int32
}

func (this_ *JsMessageVisitor_IcuMessageTemplateString) BindSelf_(s any) {
	_ = s
}
func (this_ *JsMessageVisitor_IcuMessageTemplateString) IsJavaNil_() bool { return this_ == nil }
func (*JsMessageVisitor_IcuMessageTemplateString) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JsMessageVisitor$IcuMessageTemplateString", "IcuMessageTemplateString"
}
func (this_ *JsMessageVisitor_IcuMessageTemplateString) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JsMessageVisitor_IcuMessageTemplateString) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JsMessageVisitor.IcuMessageTemplateString", "com.google.javascript.jscomp.JsMessageVisitor$IcuMessageTemplateString", func(v_ any) bool {
		return j2gort.InstanceOf[*JsMessageVisitor_IcuMessageTemplateString](v_)
	})
}

func (this_ *JsMessageVisitor_IcuMessageTemplateString) initFields_JsMessageVisitor_IcuMessageTemplateString() {
	_ = this_
	this_.template = j2gort.NullString
}

// Java constructor: IcuMessageTemplateString(java.lang.String template)
func NewJsMessageVisitor_IcuMessageTemplateString(template string) *JsMessageVisitor_IcuMessageTemplateString {
	this_ := &JsMessageVisitor_IcuMessageTemplateString{}
	InitJsMessageVisitor_IcuMessageTemplateString(this_, template)
	return this_
}

func InitJsMessageVisitor_IcuMessageTemplateString(this_ *JsMessageVisitor_IcuMessageTemplateString, template string) {
	_ = template
	this_.initFields_JsMessageVisitor_IcuMessageTemplateString()
	this_.template = template
}

// Java method: IcuMessageTemplateString.extractParts(java.util.Set<java.lang.String> placholderNamesToExtract)
func (this_ *JsMessageVisitor_IcuMessageTemplateString) ExtractParts(placholderNamesToExtract j2gort.Ext_java_util_Set) *JsMessageVisitor_ExtractedIcuTemplateParts {
	_ = placholderNamesToExtract
	var extractedPlaceholderNames j2gort.Ext_com_google_common_collect_ImmutableSet_Builder = j2gort.Ext_com_google_common_collect_ImmutableSet_S_Builder()
	_ = extractedPlaceholderNames
	var partsBuilder *j2gort.Ext_com_google_common_collect_ImmutableList_Builder = (&j2gort.Ext_com_google_common_collect_ImmutableList_Builder{})
	_ = partsBuilder
	var matcher *j2gort.Ext_java_util_regex_Matcher = StaticGet_JsMessageVisitor_ICU_PLACEHOLDER_RE().Matcher(j2gort.BoxString(this_.template))
	_ = matcher
	var remainingTemplateStartIndex int32 = int32(0)
	_ = remainingTemplateStartIndex
	for matcher.Find_0() {
		var placeholderName string = matcher.Group_1_int32(int32(1))
		_ = placeholderName
		if j2gort.SetContainsStr(placholderNamesToExtract, placeholderName) {
			j2gort.ImmSetBuilderAddStrVoid(extractedPlaceholderNames, placeholderName)
			var nonPlaceholderPrefix string = j2gort.JStrSub(this_.template, remainingTemplateStartIndex, matcher.Start_0())
			_ = nonPlaceholderPrefix
			remainingTemplateStartIndex = matcher.End_0()
			partsBuilder.AddVoid_(j2gort.PtrToAny(JsMessage_StringPart_Create(nonPlaceholderPrefix)))
			partsBuilder.AddVoid_(j2gort.PtrToAny(JsMessage_PlaceholderReference_CreateForCanonicalName(placeholderName)))
		}
	}
	if remainingTemplateStartIndex < j2gort.JStrLen(this_.template) {
		var remainingTemplate string = j2gort.JStrSubFrom(this_.template, remainingTemplateStartIndex)
		_ = remainingTemplate
		partsBuilder.AddVoid_(j2gort.PtrToAny(JsMessage_StringPart_Create(remainingTemplate)))
	}
	return NewJsMessageVisitor_ExtractedIcuTemplateParts(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](extractedPlaceholderNames.Build()), partsBuilder.BuildImmList_())
}

// Java class: com.google.javascript.jscomp.JsMessageVisitor.ExtractedIcuTemplateParts
type JsMessageVisitor_ExtractedIcuTemplateParts struct {
	extractedPlaceholderNames j2gort.Ext_com_google_common_collect_ImmutableSet  `gson:"extractedPlaceholderNames"`
	extractedParts            j2gort.Ext_com_google_common_collect_ImmutableList `gson:"extractedParts"`
	idHash_                   atomic.Int32
}

func (this_ *JsMessageVisitor_ExtractedIcuTemplateParts) BindSelf_(s any) {
	_ = s
}
func (this_ *JsMessageVisitor_ExtractedIcuTemplateParts) IsJavaNil_() bool { return this_ == nil }
func (*JsMessageVisitor_ExtractedIcuTemplateParts) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JsMessageVisitor$ExtractedIcuTemplateParts", "ExtractedIcuTemplateParts"
}
func (this_ *JsMessageVisitor_ExtractedIcuTemplateParts) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JsMessageVisitor_ExtractedIcuTemplateParts) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JsMessageVisitor.ExtractedIcuTemplateParts", "com.google.javascript.jscomp.JsMessageVisitor$ExtractedIcuTemplateParts", func(v_ any) bool {
		return j2gort.InstanceOf[*JsMessageVisitor_ExtractedIcuTemplateParts](v_)
	})
}

func (this_ *JsMessageVisitor_ExtractedIcuTemplateParts) initFields_JsMessageVisitor_ExtractedIcuTemplateParts() {
	_ = this_
}

// Java constructor: ExtractedIcuTemplateParts(com.google.common.collect.ImmutableSet<java.lang.String> extractedPlaceholderNames, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.JsMessage.Part> extractedParts)
func NewJsMessageVisitor_ExtractedIcuTemplateParts(extractedPlaceholderNames j2gort.Ext_com_google_common_collect_ImmutableSet, extractedParts j2gort.Ext_com_google_common_collect_ImmutableList) *JsMessageVisitor_ExtractedIcuTemplateParts {
	this_ := &JsMessageVisitor_ExtractedIcuTemplateParts{}
	InitJsMessageVisitor_ExtractedIcuTemplateParts(this_, extractedPlaceholderNames, extractedParts)
	return this_
}

func InitJsMessageVisitor_ExtractedIcuTemplateParts(this_ *JsMessageVisitor_ExtractedIcuTemplateParts, extractedPlaceholderNames j2gort.Ext_com_google_common_collect_ImmutableSet, extractedParts j2gort.Ext_com_google_common_collect_ImmutableList) {
	_ = extractedPlaceholderNames
	_ = extractedParts
	this_.initFields_JsMessageVisitor_ExtractedIcuTemplateParts()
	this_.extractedPlaceholderNames = extractedPlaceholderNames
	this_.extractedParts = extractedParts
}

// Java interface: com.google.javascript.jscomp.JsMessageVisitor.JsMessageOptions
type JsMessageVisitor_JsMessageOptions interface {
	Implements_jscomp_JsMessageVisitor_JsMessageOptions()
	IsJavaNil_() bool
	IsEscapeLessThan() bool
	IsUnescapeHtmlEntities() bool
	GetPlaceholderExampleMap() j2gort.Ext_com_google_common_collect_ImmutableMap
	GetPlaceholderOriginalCodeMap() j2gort.Ext_com_google_common_collect_ImmutableMap
	CheckForUnknownPlaceholders(j2gort.Ext_java_util_Set)
}

func ToIface_JsMessageVisitor_JsMessageOptions[T interface {
	JsMessageVisitor_JsMessageOptions
	comparable
}](v T) JsMessageVisitor_JsMessageOptions {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_JsMessageVisitor_JsMessageOptions[T JsMessageVisitor_JsMessageOptions](v T) JsMessageVisitor_JsMessageOptions {
	var vi JsMessageVisitor_JsMessageOptions = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JsMessageVisitor.JsMessageOptions", "com.google.javascript.jscomp.JsMessageVisitor$JsMessageOptions", func(v_ any) bool {
		return j2gort.InstanceOf[JsMessageVisitor_JsMessageOptions](v_)
	})
}

// Java interface: com.google.javascript.jscomp.JsMessageVisitor.IcuTemplateOptions
type JsMessageVisitor_IcuTemplateOptions interface {
	Implements_jscomp_JsMessageVisitor_IcuTemplateOptions()
	IsJavaNil_() bool
	GetDescription() string
	GetMeaning() string
	GetAlternateMessageId() string
	GetPlaceholderExampleMap() j2gort.Ext_com_google_common_collect_ImmutableMap
	GetPlaceholderOriginalCodeMap() j2gort.Ext_com_google_common_collect_ImmutableMap
	CheckForUnknownPlaceholders(j2gort.Ext_java_util_Set)
	GetPlaceholderNames() j2gort.Ext_com_google_common_collect_ImmutableSet
}

func ToIface_JsMessageVisitor_IcuTemplateOptions[T interface {
	JsMessageVisitor_IcuTemplateOptions
	comparable
}](v T) JsMessageVisitor_IcuTemplateOptions {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_JsMessageVisitor_IcuTemplateOptions[T JsMessageVisitor_IcuTemplateOptions](v T) JsMessageVisitor_IcuTemplateOptions {
	var vi JsMessageVisitor_IcuTemplateOptions = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JsMessageVisitor.IcuTemplateOptions", "com.google.javascript.jscomp.JsMessageVisitor$IcuTemplateOptions", func(v_ any) bool {
		return j2gort.InstanceOf[JsMessageVisitor_IcuTemplateOptions](v_)
	})
}

// Java interface: com.google.javascript.jscomp.JsMessageVisitor.ObjectLiteralMap
type JsMessageVisitor_ObjectLiteralMap interface {
	Implements_jscomp_JsMessageVisitor_ObjectLiteralMap()
	IsJavaNil_() bool
	GetBooleanValueOrFalse(string) bool
	ExtractAsValueMap() j2gort.Ext_com_google_common_collect_ImmutableMap
	ExtractAsStringToStringMap() j2gort.Ext_com_google_common_collect_ImmutableMap
	GetValueNode(string) pkg_rhino.Node
	CheckForUnexpectedKeys(j2gort.Ext_java_util_Set, j2gort.Ext_java_util_function_Function)
	CheckForRequiredKeys(j2gort.Ext_java_util_Set, j2gort.Ext_java_util_function_Function)
}

func ToIface_JsMessageVisitor_ObjectLiteralMap[T interface {
	JsMessageVisitor_ObjectLiteralMap
	comparable
}](v T) JsMessageVisitor_ObjectLiteralMap {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_JsMessageVisitor_ObjectLiteralMap[T JsMessageVisitor_ObjectLiteralMap](v T) JsMessageVisitor_ObjectLiteralMap {
	var vi JsMessageVisitor_ObjectLiteralMap = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JsMessageVisitor.ObjectLiteralMap", "com.google.javascript.jscomp.JsMessageVisitor$ObjectLiteralMap", func(v_ any) bool {
		return j2gort.InstanceOf[JsMessageVisitor_ObjectLiteralMap](v_)
	})
}

// Java class: com.google.javascript.jscomp.JsMessageVisitor.ObjectLiteralMapImpl
type JsMessageVisitor_ObjectLiteralMapImpl struct {
	stringToStringKeyMap j2gort.Ext_java_util_Map                          `gson:"stringToStringKeyMap"`
	valueMap             j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"valueMap"`
	stringMap            j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"stringMap"`
	idHash_              atomic.Int32
}

func (this_ *JsMessageVisitor_ObjectLiteralMapImpl) BindSelf_(s any) {
	_ = s
}
func (this_ *JsMessageVisitor_ObjectLiteralMapImpl) IsJavaNil_() bool { return this_ == nil }
func (*JsMessageVisitor_ObjectLiteralMapImpl) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JsMessageVisitor$ObjectLiteralMapImpl", "ObjectLiteralMapImpl"
}
func (this_ *JsMessageVisitor_ObjectLiteralMapImpl) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JsMessageVisitor_ObjectLiteralMapImpl) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JsMessageVisitor.ObjectLiteralMapImpl", "com.google.javascript.jscomp.JsMessageVisitor$ObjectLiteralMapImpl", func(v_ any) bool {
		return j2gort.InstanceOf[*JsMessageVisitor_ObjectLiteralMapImpl](v_)
	})
}
func (*JsMessageVisitor_ObjectLiteralMapImpl) Implements_jscomp_JsMessageVisitor_ObjectLiteralMap() {}

func (this_ *JsMessageVisitor_ObjectLiteralMapImpl) initFields_JsMessageVisitor_ObjectLiteralMapImpl() {
	_ = this_

}

// Java constructor: ObjectLiteralMapImpl(java.util.Map<java.lang.String,com.google.javascript.rhino.Node> stringToStringKeyMap)
func NewJsMessageVisitor_ObjectLiteralMapImpl(stringToStringKeyMap j2gort.Ext_java_util_Map) *JsMessageVisitor_ObjectLiteralMapImpl {
	this_ := &JsMessageVisitor_ObjectLiteralMapImpl{}
	InitJsMessageVisitor_ObjectLiteralMapImpl(this_, stringToStringKeyMap)
	return this_
}

func InitJsMessageVisitor_ObjectLiteralMapImpl(this_ *JsMessageVisitor_ObjectLiteralMapImpl, stringToStringKeyMap j2gort.Ext_java_util_Map) {
	_ = stringToStringKeyMap
	this_.initFields_JsMessageVisitor_ObjectLiteralMapImpl()
	this_.stringToStringKeyMap = stringToStringKeyMap
}

// Java method: ObjectLiteralMapImpl.getValueNode(java.lang.String key)
func (this_ *JsMessageVisitor_ObjectLiteralMapImpl) GetValueNode(key string) pkg_rhino.Node {
	_ = key
	var stringKeyNode pkg_rhino.Node = j2gort.CastNode[pkg_rhino.Node](j2gort.MapGetStr(this_.stringToStringKeyMap, key))
	_ = stringKeyNode
	if stringKeyNode == 0 {
		return pkg_rhino.Node(0)
	} else {
		return stringKeyNode.GetOnlyChild()
	}
}

// Java method: ObjectLiteralMapImpl.getBooleanValueOrFalse(java.lang.String key)
func (this_ *JsMessageVisitor_ObjectLiteralMapImpl) GetBooleanValueOrFalse(key string) bool {
	_ = key
	var stringKeyNode pkg_rhino.Node = j2gort.CastNode[pkg_rhino.Node](j2gort.MapGetStr(this_.stringToStringKeyMap, key))
	_ = stringKeyNode
	return JsMessageVisitor_ExtractBooleanStringKeyValue(stringKeyNode)
}

// Java method: ObjectLiteralMapImpl.extractAsValueMap()
func (this_ *JsMessageVisitor_ObjectLiteralMapImpl) ExtractAsValueMap() j2gort.Ext_com_google_common_collect_ImmutableMap {
	if this_.valueMap == nil {
		var builder j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
		_ = builder
		{
			cur_12 := j2gort.NewMapEntryCursor(this_.stringToStringKeyMap)
			for cur_12.NextEntry() {
				builder.Put_2(j2gort.BoxString(cur_12.KeyStr()), j2gort.BoxNode(j2gort.CastNode[pkg_rhino.Node](cur_12.GetValue()).GetOnlyChild()))
			}
		}
		this_.valueMap = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](builder.BuildOrThrow())
	}
	return this_.valueMap
}

// Java method: ObjectLiteralMapImpl.extractAsStringToStringMap()
func (this_ *JsMessageVisitor_ObjectLiteralMapImpl) ExtractAsStringToStringMap() j2gort.Ext_com_google_common_collect_ImmutableMap {
	if this_.stringMap == nil {
		var builder j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
		_ = builder
		{
			cur_13 := j2gort.NewMapEntryCursor(this_.stringToStringKeyMap)
			for cur_13.NextEntry() {
				builder.Put_2(j2gort.BoxString(cur_13.KeyStr()), j2gort.BoxString(JsMessageVisitor_ExtractStringFromStringExprNode(j2gort.CastNode[pkg_rhino.Node](cur_13.GetValue()).GetOnlyChild())))
			}
		}
		this_.stringMap = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](builder.BuildOrThrow())
	}
	return this_.stringMap
}

// Java method: ObjectLiteralMapImpl.checkForUnexpectedKeys(java.util.Set<java.lang.String> expectedKeys, java.util.function.Function<java.lang.String,java.lang.String> createErrorMessage)
func (this_ *JsMessageVisitor_ObjectLiteralMapImpl) CheckForUnexpectedKeys(expectedKeys j2gort.Ext_java_util_Set, createErrorMessage j2gort.Ext_java_util_function_Function) {
	_ = expectedKeys
	_ = createErrorMessage
	{
		cur_14 := j2gort.NewMapKeyStrCursor(this_.stringToStringKeyMap)
		for cur_14.HasNext() {
			key := cur_14.NextStr()
			_ = key
			if !(j2gort.SetContainsStr(expectedKeys, key)) {
				panic(NewJsMessageVisitor_MalformedException(j2gort.CastString(createErrorMessage.Apply(j2gort.BoxString(key))), j2gort.CastNode[pkg_rhino.Node](j2gort.MapGetStr(this_.stringToStringKeyMap, key))))
			}
		}
	}
}

// Java method: ObjectLiteralMapImpl.checkForRequiredKeys(java.util.Set<java.lang.String> requiredKeys, java.util.function.Function<java.lang.String,com.google.javascript.jscomp.JsMessageVisitor.MalformedException> createException)
func (this_ *JsMessageVisitor_ObjectLiteralMapImpl) CheckForRequiredKeys(requiredKeys j2gort.Ext_java_util_Set, createException j2gort.Ext_java_util_function_Function) {
	_ = requiredKeys
	_ = createException
	{
		cur_15 := j2gort.NewIterCursor(j2gort.CheckNotNull(requiredKeys))
		for cur_15.HasNext() {
			requiredKey := cur_15.NextStr()
			_ = requiredKey
			if !(j2gort.MapContainsKeyStr(this_.stringToStringKeyMap, requiredKey)) {
				panic(j2gort.CheckNotNullPtr(j2gort.CastPtr[JsMessageVisitor_MalformedException](createException.Apply(j2gort.BoxString(requiredKey)))))
			}
		}
	}
}

// Java interface: com.google.javascript.jscomp.JsMessageVisitor.GoogGetMsgParsedText
type JsMessageVisitor_GoogGetMsgParsedText interface {
	Implements_jscomp_JsMessageVisitor_GoogGetMsgParsedText()
	IsJavaNil_() bool
	GetText() string
	GetPlaceholderNames() j2gort.Ext_com_google_common_collect_ImmutableSet
	GetParts() j2gort.Ext_com_google_common_collect_ImmutableList
}

func ToIface_JsMessageVisitor_GoogGetMsgParsedText[T interface {
	JsMessageVisitor_GoogGetMsgParsedText
	comparable
}](v T) JsMessageVisitor_GoogGetMsgParsedText {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_JsMessageVisitor_GoogGetMsgParsedText[T JsMessageVisitor_GoogGetMsgParsedText](v T) JsMessageVisitor_GoogGetMsgParsedText {
	var vi JsMessageVisitor_GoogGetMsgParsedText = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JsMessageVisitor.GoogGetMsgParsedText", "com.google.javascript.jscomp.JsMessageVisitor$GoogGetMsgParsedText", func(v_ any) bool {
		return j2gort.InstanceOf[JsMessageVisitor_GoogGetMsgParsedText](v_)
	})
}

// Java class: com.google.javascript.jscomp.JsMessageVisitor.MalformedException
type JsMessageVisitor_MalformedException struct {
	javalang.JavaException
	node    pkg_rhino.Node `gson:"node"`
	idHash_ atomic.Int32
}

func (this_ *JsMessageVisitor_MalformedException) BindSelf_(s any) {
	_ = s
}
func (this_ *JsMessageVisitor_MalformedException) IsJavaNil_() bool { return this_ == nil }
func (*JsMessageVisitor_MalformedException) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JsMessageVisitor$MalformedException", "MalformedException"
}
func (this_ *JsMessageVisitor_MalformedException) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JsMessageVisitor.MalformedException", "com.google.javascript.jscomp.JsMessageVisitor$MalformedException", func(v_ any) bool {
		return j2gort.InstanceOf[*JsMessageVisitor_MalformedException](v_)
	})
}

func StaticGet_JsMessageVisitor_MalformedException_SerialVersionUID() int64 { return int64(1) }
func (this_ *JsMessageVisitor_MalformedException) initFields_JsMessageVisitor_MalformedException() {
	_ = this_
}

// Java constructor: MalformedException(java.lang.String message, com.google.javascript.rhino.Node node)
func NewJsMessageVisitor_MalformedException(message string, node pkg_rhino.Node) *JsMessageVisitor_MalformedException {
	this_ := &JsMessageVisitor_MalformedException{}
	InitJsMessageVisitor_MalformedException(this_, message, node)
	return this_
}

func InitJsMessageVisitor_MalformedException(this_ *JsMessageVisitor_MalformedException, message string, node pkg_rhino.Node) {
	_ = message
	_ = node
	j2gort.InitUserThrowable(this_, "com.google.javascript.jscomp.JsMessageVisitor.MalformedException", j2gort.BoxString(message))
	this_.initFields_JsMessageVisitor_MalformedException()
	this_.node = node
}

// Java method: MalformedException.getNode()
func (this_ *JsMessageVisitor_MalformedException) GetNode() pkg_rhino.Node {
	return this_.node
}

// Java class: com.google.javascript.jscomp.JsMessageVisitor.MessageLocation
type JsMessageVisitor_MessageLocation struct {
	message     *JsMessage     `gson:"message"`
	messageNode pkg_rhino.Node `gson:"messageNode"`
	idHash_     atomic.Int32
}

func (this_ *JsMessageVisitor_MessageLocation) BindSelf_(s any) {
	_ = s
}
func (this_ *JsMessageVisitor_MessageLocation) IsJavaNil_() bool { return this_ == nil }
func (*JsMessageVisitor_MessageLocation) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JsMessageVisitor$MessageLocation", "MessageLocation"
}
func (this_ *JsMessageVisitor_MessageLocation) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JsMessageVisitor_MessageLocation) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.JsMessageVisitor.MessageLocation", "com.google.javascript.jscomp.JsMessageVisitor$MessageLocation", func(v_ any) bool {
		return j2gort.InstanceOf[*JsMessageVisitor_MessageLocation](v_)
	})
}

func (this_ *JsMessageVisitor_MessageLocation) initFields_JsMessageVisitor_MessageLocation() {
	_ = this_
}

// Java constructor: MessageLocation(com.google.javascript.jscomp.JsMessage message, com.google.javascript.rhino.Node messageNode)
func NewJsMessageVisitor_MessageLocation(message *JsMessage, messageNode pkg_rhino.Node) *JsMessageVisitor_MessageLocation {
	this_ := &JsMessageVisitor_MessageLocation{}
	InitJsMessageVisitor_MessageLocation(this_, message, messageNode)
	return this_
}

func InitJsMessageVisitor_MessageLocation(this_ *JsMessageVisitor_MessageLocation, message *JsMessage, messageNode pkg_rhino.Node) {
	_ = message
	_ = messageNode
	this_.initFields_JsMessageVisitor_MessageLocation()
	this_.message = message
	this_.messageNode = messageNode
}

// Java class:
type JsMessageVisitor__Anon_56 struct {
	outer_                     JsMessageVisitor
	cap_0_messageKeyFromLhs    string
	cap_1_sourceName           string
	cap_2_googGetMsgParsedText JsMessageVisitor_GoogGetMsgParsedText
	cap_3_jsMessageOptions     JsMessageVisitor_JsMessageOptions
	cap_4_description          string
	cap_5_meaning              string
	cap_6_alternateMessageId   string
	idHash_                    atomic.Int32
}

func (this_ *JsMessageVisitor__Anon_56) BindSelf_(s any) {
	_ = s
}
func (this_ *JsMessageVisitor__Anon_56) IsJavaNil_() bool { return this_ == nil }
func (*JsMessageVisitor__Anon_56) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JsMessageVisitor$1", ""
}
func (this_ *JsMessageVisitor__Anon_56) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JsMessageVisitor__Anon_56) HasIdentityEquals_()                                    {}
func (this_ *JsMessageVisitor__Anon_56) SetOuter_(o JsMessageVisitor)                     { this_.outer_ = o }
func (this_ *JsMessageVisitor__Anon_56) GetOuter_() JsMessageVisitor                      { return this_.outer_ }
func (*JsMessageVisitor__Anon_56) Implements_jscomp_JsMessageVisitor_InputForBuildJsMsg() {}

func (this_ *JsMessageVisitor__Anon_56) initFields_JsMessageVisitor__Anon_56() {
	_ = this_
}

func NewJsMessageVisitor__Anon_56(outer_ JsMessageVisitor, cap_p_messageKeyFromLhs string, cap_p_sourceName string, cap_p_googGetMsgParsedText JsMessageVisitor_GoogGetMsgParsedText, cap_p_jsMessageOptions JsMessageVisitor_JsMessageOptions, cap_p_description string, cap_p_meaning string, cap_p_alternateMessageId string) *JsMessageVisitor__Anon_56 {
	this_ := &JsMessageVisitor__Anon_56{}
	this_.outer_ = outer_
	this_.cap_0_messageKeyFromLhs = cap_p_messageKeyFromLhs
	this_.cap_1_sourceName = cap_p_sourceName
	this_.cap_2_googGetMsgParsedText = cap_p_googGetMsgParsedText
	this_.cap_3_jsMessageOptions = cap_p_jsMessageOptions
	this_.cap_4_description = cap_p_description
	this_.cap_5_meaning = cap_p_meaning
	this_.cap_6_alternateMessageId = cap_p_alternateMessageId
	this_.initFields_JsMessageVisitor__Anon_56()
	return this_
}

// Java method: .getMessageKeyFromLhs()
func (this_ *JsMessageVisitor__Anon_56) GetMessageKeyFromLhs() string {
	return this_.cap_0_messageKeyFromLhs
}

// Java method: .getSourceName()
func (this_ *JsMessageVisitor__Anon_56) GetSourceName() string {
	return this_.cap_1_sourceName
}

// Java method: .getMessageText()
func (this_ *JsMessageVisitor__Anon_56) GetMessageText() string {
	return this_.cap_2_googGetMsgParsedText.GetText()
}

// Java method: .getMessageParts()
func (this_ *JsMessageVisitor__Anon_56) GetMessageParts() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.cap_2_googGetMsgParsedText.GetParts()
}

// Java method: .getPlaceholderExampleMap()
func (this_ *JsMessageVisitor__Anon_56) GetPlaceholderExampleMap() any {
	return this_.cap_3_jsMessageOptions.GetPlaceholderExampleMap()
}

// Java method: .getPlaceholderOriginalCodeMap()
func (this_ *JsMessageVisitor__Anon_56) GetPlaceholderOriginalCodeMap() any {
	return this_.cap_3_jsMessageOptions.GetPlaceholderOriginalCodeMap()
}

// Java method: .getDescription()
func (this_ *JsMessageVisitor__Anon_56) GetDescription() string {
	return this_.cap_4_description
}

// Java method: .getMeaning()
func (this_ *JsMessageVisitor__Anon_56) GetMeaning() string {
	return this_.cap_5_meaning
}

// Java method: .getAlternateMessageId()
func (this_ *JsMessageVisitor__Anon_56) GetAlternateMessageId() string {
	return this_.cap_6_alternateMessageId
}

// Java class:
type JsMessageVisitor__Anon_57 struct {
	outer_                     JsMessageVisitor
	cap_0_jsExtractedMessage   *JsMessage
	cap_4_placeholderValuesMap j2gort.Ext_com_google_common_collect_ImmutableMap
	cap_5_jsMessageOptions     JsMessageVisitor_JsMessageOptions
	cap_1_msgNode              pkg_rhino.Node
	cap_2_msgTextNode          pkg_rhino.Node
	cap_3_valuesObjLit         pkg_rhino.Node
	idHash_                    atomic.Int32
}

func (this_ *JsMessageVisitor__Anon_57) BindSelf_(s any) {
	_ = s
}
func (this_ *JsMessageVisitor__Anon_57) IsJavaNil_() bool { return this_ == nil }
func (*JsMessageVisitor__Anon_57) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JsMessageVisitor$2", ""
}
func (this_ *JsMessageVisitor__Anon_57) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JsMessageVisitor__Anon_57) HasIdentityEquals_()                    {}
func (this_ *JsMessageVisitor__Anon_57) SetOuter_(o JsMessageVisitor)     { this_.outer_ = o }
func (this_ *JsMessageVisitor__Anon_57) GetOuter_() JsMessageVisitor      { return this_.outer_ }
func (*JsMessageVisitor__Anon_57) Implements_jscomp_JsMessageDefinition() {}

func (this_ *JsMessageVisitor__Anon_57) initFields_JsMessageVisitor__Anon_57() {
	_ = this_
}

func NewJsMessageVisitor__Anon_57(outer_ JsMessageVisitor, cap_p_jsExtractedMessage *JsMessage, cap_p_msgNode pkg_rhino.Node, cap_p_msgTextNode pkg_rhino.Node, cap_p_valuesObjLit pkg_rhino.Node, cap_p_placeholderValuesMap j2gort.Ext_com_google_common_collect_ImmutableMap, cap_p_jsMessageOptions JsMessageVisitor_JsMessageOptions) *JsMessageVisitor__Anon_57 {
	this_ := &JsMessageVisitor__Anon_57{}
	this_.outer_ = outer_
	this_.cap_0_jsExtractedMessage = cap_p_jsExtractedMessage
	this_.cap_1_msgNode = cap_p_msgNode
	this_.cap_2_msgTextNode = cap_p_msgTextNode
	this_.cap_3_valuesObjLit = cap_p_valuesObjLit
	this_.cap_4_placeholderValuesMap = cap_p_placeholderValuesMap
	this_.cap_5_jsMessageOptions = cap_p_jsMessageOptions
	this_.initFields_JsMessageVisitor__Anon_57()
	return this_
}

// Java method: .getMessage()
func (this_ *JsMessageVisitor__Anon_57) GetMessage() *JsMessage {
	return this_.cap_0_jsExtractedMessage
}

// Java method: .getMessageNode()
func (this_ *JsMessageVisitor__Anon_57) GetMessageNode() pkg_rhino.Node {
	return this_.cap_1_msgNode
}

// Java method: .getTemplateTextNode()
func (this_ *JsMessageVisitor__Anon_57) GetTemplateTextNode() pkg_rhino.Node {
	return this_.cap_2_msgTextNode
}

// Java method: .getPlaceholderValuesNode()
func (this_ *JsMessageVisitor__Anon_57) GetPlaceholderValuesNode() pkg_rhino.Node {
	return this_.cap_3_valuesObjLit
}

// Java method: .getPlaceholderValueMap()
func (this_ *JsMessageVisitor__Anon_57) GetPlaceholderValueMap() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.cap_4_placeholderValuesMap
}

// Java method: .shouldEscapeLessThan()
func (this_ *JsMessageVisitor__Anon_57) ShouldEscapeLessThan() bool {
	return this_.cap_5_jsMessageOptions.IsEscapeLessThan()
}

// Java method: .shouldUnescapeHtmlEntities()
func (this_ *JsMessageVisitor__Anon_57) ShouldUnescapeHtmlEntities() bool {
	return this_.cap_5_jsMessageOptions.IsUnescapeHtmlEntities()
}

// Java class:
type JsMessageVisitor__Anon_58 struct {
	outer_                          JsMessageVisitor
	cap_0_messageKeyFromLhs         string
	cap_1_sourceName                string
	cap_2_icuMessageTemplateString  *JsMessageVisitor_IcuMessageTemplateString
	cap_3_extractedIcuTemplateParts *JsMessageVisitor_ExtractedIcuTemplateParts
	cap_4_icuTemplateOptions        JsMessageVisitor_IcuTemplateOptions
	idHash_                         atomic.Int32
}

func (this_ *JsMessageVisitor__Anon_58) BindSelf_(s any) {
	_ = s
}
func (this_ *JsMessageVisitor__Anon_58) IsJavaNil_() bool { return this_ == nil }
func (*JsMessageVisitor__Anon_58) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JsMessageVisitor$3", ""
}
func (this_ *JsMessageVisitor__Anon_58) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JsMessageVisitor__Anon_58) HasIdentityEquals_()                                    {}
func (this_ *JsMessageVisitor__Anon_58) SetOuter_(o JsMessageVisitor)                     { this_.outer_ = o }
func (this_ *JsMessageVisitor__Anon_58) GetOuter_() JsMessageVisitor                      { return this_.outer_ }
func (*JsMessageVisitor__Anon_58) Implements_jscomp_JsMessageVisitor_InputForBuildJsMsg() {}

func (this_ *JsMessageVisitor__Anon_58) initFields_JsMessageVisitor__Anon_58() {
	_ = this_
}

func NewJsMessageVisitor__Anon_58(outer_ JsMessageVisitor, cap_p_messageKeyFromLhs string, cap_p_sourceName string, cap_p_icuMessageTemplateString *JsMessageVisitor_IcuMessageTemplateString, cap_p_extractedIcuTemplateParts *JsMessageVisitor_ExtractedIcuTemplateParts, cap_p_icuTemplateOptions JsMessageVisitor_IcuTemplateOptions) *JsMessageVisitor__Anon_58 {
	this_ := &JsMessageVisitor__Anon_58{}
	this_.outer_ = outer_
	this_.cap_0_messageKeyFromLhs = cap_p_messageKeyFromLhs
	this_.cap_1_sourceName = cap_p_sourceName
	this_.cap_2_icuMessageTemplateString = cap_p_icuMessageTemplateString
	this_.cap_3_extractedIcuTemplateParts = cap_p_extractedIcuTemplateParts
	this_.cap_4_icuTemplateOptions = cap_p_icuTemplateOptions
	this_.initFields_JsMessageVisitor__Anon_58()
	return this_
}

// Java method: .getMessageKeyFromLhs()
func (this_ *JsMessageVisitor__Anon_58) GetMessageKeyFromLhs() string {
	return this_.cap_0_messageKeyFromLhs
}

// Java method: .getSourceName()
func (this_ *JsMessageVisitor__Anon_58) GetSourceName() string {
	return this_.cap_1_sourceName
}

// Java method: .getMessageText()
func (this_ *JsMessageVisitor__Anon_58) GetMessageText() string {
	return this_.cap_2_icuMessageTemplateString.template
}

// Java method: .getMessageParts()
func (this_ *JsMessageVisitor__Anon_58) GetMessageParts() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.cap_3_extractedIcuTemplateParts.extractedParts
}

// Java method: .getPlaceholderExampleMap()
func (this_ *JsMessageVisitor__Anon_58) GetPlaceholderExampleMap() any {
	return this_.cap_4_icuTemplateOptions.GetPlaceholderExampleMap()
}

// Java method: .getPlaceholderOriginalCodeMap()
func (this_ *JsMessageVisitor__Anon_58) GetPlaceholderOriginalCodeMap() any {
	return this_.cap_4_icuTemplateOptions.GetPlaceholderOriginalCodeMap()
}

// Java method: .getDescription()
func (this_ *JsMessageVisitor__Anon_58) GetDescription() string {
	return this_.cap_4_icuTemplateOptions.GetDescription()
}

// Java method: .getMeaning()
func (this_ *JsMessageVisitor__Anon_58) GetMeaning() string {
	return this_.cap_4_icuTemplateOptions.GetMeaning()
}

// Java method: .getAlternateMessageId()
func (this_ *JsMessageVisitor__Anon_58) GetAlternateMessageId() string {
	return this_.cap_4_icuTemplateOptions.GetAlternateMessageId()
}

// Java class:
type JsMessageVisitor__Anon_59 struct {
	outer_                        JsMessageVisitor
	cap_0_extractedMessage        *JsMessage
	cap_1_msgNode                 pkg_rhino.Node
	cap_2_stringLiteralExpression pkg_rhino.Node
	idHash_                       atomic.Int32
}

func (this_ *JsMessageVisitor__Anon_59) BindSelf_(s any) {
	_ = s
}
func (this_ *JsMessageVisitor__Anon_59) IsJavaNil_() bool { return this_ == nil }
func (*JsMessageVisitor__Anon_59) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JsMessageVisitor$4", ""
}
func (this_ *JsMessageVisitor__Anon_59) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JsMessageVisitor__Anon_59) HasIdentityEquals_()                      {}
func (this_ *JsMessageVisitor__Anon_59) SetOuter_(o JsMessageVisitor)       { this_.outer_ = o }
func (this_ *JsMessageVisitor__Anon_59) GetOuter_() JsMessageVisitor        { return this_.outer_ }
func (*JsMessageVisitor__Anon_59) Implements_jscomp_IcuTemplateDefinition() {}

func (this_ *JsMessageVisitor__Anon_59) initFields_JsMessageVisitor__Anon_59() {
	_ = this_
}

func NewJsMessageVisitor__Anon_59(outer_ JsMessageVisitor, cap_p_extractedMessage *JsMessage, cap_p_msgNode pkg_rhino.Node, cap_p_stringLiteralExpression pkg_rhino.Node) *JsMessageVisitor__Anon_59 {
	this_ := &JsMessageVisitor__Anon_59{}
	this_.outer_ = outer_
	this_.cap_0_extractedMessage = cap_p_extractedMessage
	this_.cap_1_msgNode = cap_p_msgNode
	this_.cap_2_stringLiteralExpression = cap_p_stringLiteralExpression
	this_.initFields_JsMessageVisitor__Anon_59()
	return this_
}

// Java method: .getMessage()
func (this_ *JsMessageVisitor__Anon_59) GetMessage() *JsMessage {
	return this_.cap_0_extractedMessage
}

// Java method: .getMessageNode()
func (this_ *JsMessageVisitor__Anon_59) GetMessageNode() pkg_rhino.Node {
	return this_.cap_1_msgNode
}

// Java method: .getTemplateTextNode()
func (this_ *JsMessageVisitor__Anon_59) GetTemplateTextNode() pkg_rhino.Node {
	return this_.cap_2_stringLiteralExpression
}

// Java class:
type JsMessageVisitor__Anon_60 struct {
	outer_                             JsMessageVisitor
	cap_2_placeholderExamplesMap       j2gort.Ext_com_google_common_collect_ImmutableMap
	cap_3_placeholderOriginalCodeMap   j2gort.Ext_com_google_common_collect_ImmutableMap
	cap_4_exampleObjectLiteralMap      JsMessageVisitor_ObjectLiteralMap
	cap_5_originalCodeObjectLiteralMap JsMessageVisitor_ObjectLiteralMap
	idHash_                            atomic.Int32
	cap_0_isEscapeLessThan             bool
	cap_1_isUnescapeHtmlEntities       bool
}

func (this_ *JsMessageVisitor__Anon_60) BindSelf_(s any) {
	_ = s
}
func (this_ *JsMessageVisitor__Anon_60) IsJavaNil_() bool { return this_ == nil }
func (*JsMessageVisitor__Anon_60) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JsMessageVisitor$5", ""
}
func (this_ *JsMessageVisitor__Anon_60) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JsMessageVisitor__Anon_60) HasIdentityEquals_()                                  {}
func (this_ *JsMessageVisitor__Anon_60) SetOuter_(o JsMessageVisitor)                   { this_.outer_ = o }
func (this_ *JsMessageVisitor__Anon_60) GetOuter_() JsMessageVisitor                    { return this_.outer_ }
func (*JsMessageVisitor__Anon_60) Implements_jscomp_JsMessageVisitor_JsMessageOptions() {}

func (this_ *JsMessageVisitor__Anon_60) initFields_JsMessageVisitor__Anon_60() {
	_ = this_
}

func NewJsMessageVisitor__Anon_60(outer_ JsMessageVisitor, cap_p_isEscapeLessThan bool, cap_p_isUnescapeHtmlEntities bool, cap_p_placeholderExamplesMap j2gort.Ext_com_google_common_collect_ImmutableMap, cap_p_placeholderOriginalCodeMap j2gort.Ext_com_google_common_collect_ImmutableMap, cap_p_exampleObjectLiteralMap JsMessageVisitor_ObjectLiteralMap, cap_p_originalCodeObjectLiteralMap JsMessageVisitor_ObjectLiteralMap) *JsMessageVisitor__Anon_60 {
	this_ := &JsMessageVisitor__Anon_60{}
	this_.outer_ = outer_
	this_.cap_0_isEscapeLessThan = cap_p_isEscapeLessThan
	this_.cap_1_isUnescapeHtmlEntities = cap_p_isUnescapeHtmlEntities
	this_.cap_2_placeholderExamplesMap = cap_p_placeholderExamplesMap
	this_.cap_3_placeholderOriginalCodeMap = cap_p_placeholderOriginalCodeMap
	this_.cap_4_exampleObjectLiteralMap = cap_p_exampleObjectLiteralMap
	this_.cap_5_originalCodeObjectLiteralMap = cap_p_originalCodeObjectLiteralMap
	this_.initFields_JsMessageVisitor__Anon_60()
	return this_
}

// Java method: .isEscapeLessThan()
func (this_ *JsMessageVisitor__Anon_60) IsEscapeLessThan() bool {
	return this_.cap_0_isEscapeLessThan
}

// Java method: .isUnescapeHtmlEntities()
func (this_ *JsMessageVisitor__Anon_60) IsUnescapeHtmlEntities() bool {
	return this_.cap_1_isUnescapeHtmlEntities
}

// Java method: .getPlaceholderExampleMap()
func (this_ *JsMessageVisitor__Anon_60) GetPlaceholderExampleMap() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.cap_2_placeholderExamplesMap
}

// Java method: .getPlaceholderOriginalCodeMap()
func (this_ *JsMessageVisitor__Anon_60) GetPlaceholderOriginalCodeMap() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.cap_3_placeholderOriginalCodeMap
}

// Java method: .checkForUnknownPlaceholders(java.util.Set<java.lang.String> knownPlaceholders)
func (this_ *JsMessageVisitor__Anon_60) CheckForUnknownPlaceholders(knownPlaceholders j2gort.Ext_java_util_Set) {
	_ = knownPlaceholders
	this_.cap_4_exampleObjectLiteralMap.CheckForUnexpectedKeys(knownPlaceholders, lambdaSingleton_js_message_visitor_4)
	this_.cap_5_originalCodeObjectLiteralMap.CheckForUnexpectedKeys(knownPlaceholders, lambdaSingleton_js_message_visitor_5)
}

// Java class:
type JsMessageVisitor__Anon_61 struct {
	outer_                             JsMessageVisitor
	cap_0_description                  string
	cap_1_meaning                      string
	cap_2_alternateMessageId           string
	cap_3_placeholderExamplesMap       j2gort.Ext_com_google_common_collect_ImmutableMap
	cap_4_placeholderOriginalCodeMap   j2gort.Ext_com_google_common_collect_ImmutableMap
	cap_5_placeholderNames             j2gort.Ext_com_google_common_collect_ImmutableSet
	cap_6_exampleObjectLiteralMap      JsMessageVisitor_ObjectLiteralMap
	cap_7_originalCodeObjectLiteralMap JsMessageVisitor_ObjectLiteralMap
	idHash_                            atomic.Int32
}

func (this_ *JsMessageVisitor__Anon_61) BindSelf_(s any) {
	_ = s
}
func (this_ *JsMessageVisitor__Anon_61) IsJavaNil_() bool { return this_ == nil }
func (*JsMessageVisitor__Anon_61) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JsMessageVisitor$6", ""
}
func (this_ *JsMessageVisitor__Anon_61) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JsMessageVisitor__Anon_61) HasIdentityEquals_()                                    {}
func (this_ *JsMessageVisitor__Anon_61) SetOuter_(o JsMessageVisitor)                     { this_.outer_ = o }
func (this_ *JsMessageVisitor__Anon_61) GetOuter_() JsMessageVisitor                      { return this_.outer_ }
func (*JsMessageVisitor__Anon_61) Implements_jscomp_JsMessageVisitor_IcuTemplateOptions() {}

func (this_ *JsMessageVisitor__Anon_61) initFields_JsMessageVisitor__Anon_61() {
	_ = this_
}

func NewJsMessageVisitor__Anon_61(outer_ JsMessageVisitor, cap_p_description string, cap_p_meaning string, cap_p_alternateMessageId string, cap_p_placeholderExamplesMap j2gort.Ext_com_google_common_collect_ImmutableMap, cap_p_placeholderOriginalCodeMap j2gort.Ext_com_google_common_collect_ImmutableMap, cap_p_placeholderNames j2gort.Ext_com_google_common_collect_ImmutableSet, cap_p_exampleObjectLiteralMap JsMessageVisitor_ObjectLiteralMap, cap_p_originalCodeObjectLiteralMap JsMessageVisitor_ObjectLiteralMap) *JsMessageVisitor__Anon_61 {
	this_ := &JsMessageVisitor__Anon_61{}
	this_.outer_ = outer_
	this_.cap_0_description = cap_p_description
	this_.cap_1_meaning = cap_p_meaning
	this_.cap_2_alternateMessageId = cap_p_alternateMessageId
	this_.cap_3_placeholderExamplesMap = cap_p_placeholderExamplesMap
	this_.cap_4_placeholderOriginalCodeMap = cap_p_placeholderOriginalCodeMap
	this_.cap_5_placeholderNames = cap_p_placeholderNames
	this_.cap_6_exampleObjectLiteralMap = cap_p_exampleObjectLiteralMap
	this_.cap_7_originalCodeObjectLiteralMap = cap_p_originalCodeObjectLiteralMap
	this_.initFields_JsMessageVisitor__Anon_61()
	return this_
}

// Java method: .getDescription()
func (this_ *JsMessageVisitor__Anon_61) GetDescription() string {
	return this_.cap_0_description
}

// Java method: .getMeaning()
func (this_ *JsMessageVisitor__Anon_61) GetMeaning() string {
	return this_.cap_1_meaning
}

// Java method: .getAlternateMessageId()
func (this_ *JsMessageVisitor__Anon_61) GetAlternateMessageId() string {
	return this_.cap_2_alternateMessageId
}

// Java method: .getPlaceholderExampleMap()
func (this_ *JsMessageVisitor__Anon_61) GetPlaceholderExampleMap() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.cap_3_placeholderExamplesMap
}

// Java method: .getPlaceholderOriginalCodeMap()
func (this_ *JsMessageVisitor__Anon_61) GetPlaceholderOriginalCodeMap() j2gort.Ext_com_google_common_collect_ImmutableMap {
	return this_.cap_4_placeholderOriginalCodeMap
}

// Java method: .getPlaceholderNames()
func (this_ *JsMessageVisitor__Anon_61) GetPlaceholderNames() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.cap_5_placeholderNames
}

// Java method: .checkForUnknownPlaceholders(java.util.Set<java.lang.String> knownPlaceholders)
func (this_ *JsMessageVisitor__Anon_61) CheckForUnknownPlaceholders(knownPlaceholders j2gort.Ext_java_util_Set) {
	_ = knownPlaceholders
	this_.cap_6_exampleObjectLiteralMap.CheckForUnexpectedKeys(knownPlaceholders, lambdaSingleton_js_message_visitor_6)
	this_.cap_7_originalCodeObjectLiteralMap.CheckForUnexpectedKeys(knownPlaceholders, lambdaSingleton_js_message_visitor_7)
}

// Java class:
type JsMessageVisitor__Anon_62 struct {
	cap_0_originalMsgText  string
	cap_1_placeholderNames j2gort.Ext_com_google_common_collect_ImmutableSet
	cap_2_parts            j2gort.Ext_com_google_common_collect_ImmutableList
	idHash_                atomic.Int32
}

func (this_ *JsMessageVisitor__Anon_62) BindSelf_(s any) {
	_ = s
}
func (this_ *JsMessageVisitor__Anon_62) IsJavaNil_() bool { return this_ == nil }
func (*JsMessageVisitor__Anon_62) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.JsMessageVisitor$7", ""
}
func (this_ *JsMessageVisitor__Anon_62) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*JsMessageVisitor__Anon_62) HasIdentityEquals_()                                      {}
func (*JsMessageVisitor__Anon_62) Implements_jscomp_JsMessageVisitor_GoogGetMsgParsedText() {}

func (this_ *JsMessageVisitor__Anon_62) initFields_JsMessageVisitor__Anon_62() {
	_ = this_
}

func NewJsMessageVisitor__Anon_62(cap_p_originalMsgText string, cap_p_placeholderNames j2gort.Ext_com_google_common_collect_ImmutableSet, cap_p_parts j2gort.Ext_com_google_common_collect_ImmutableList) *JsMessageVisitor__Anon_62 {
	this_ := &JsMessageVisitor__Anon_62{}
	this_.cap_0_originalMsgText = cap_p_originalMsgText
	this_.cap_1_placeholderNames = cap_p_placeholderNames
	this_.cap_2_parts = cap_p_parts
	this_.initFields_JsMessageVisitor__Anon_62()
	return this_
}

// Java method: .getText()
func (this_ *JsMessageVisitor__Anon_62) GetText() string {
	return this_.cap_0_originalMsgText
}

// Java method: .getPlaceholderNames()
func (this_ *JsMessageVisitor__Anon_62) GetPlaceholderNames() j2gort.Ext_com_google_common_collect_ImmutableSet {
	return this_.cap_1_placeholderNames
}

// Java method: .getParts()
func (this_ *JsMessageVisitor__Anon_62) GetParts() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.cap_2_parts
}
