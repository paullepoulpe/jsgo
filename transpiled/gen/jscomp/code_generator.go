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
// Port of com/google/javascript/jscomp/CodeGenerator.java (generated mechanically by j2go).

package jscomp

import (
	"math"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	pkg_sourcemap "github.com/paullepoulpe/jsgo/transpiled/gen/sourcemap"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_code_generator_1 j2gort.Ext_com_google_common_base_Predicate

func init() {
	lambdaSingleton_code_generator_1 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_ uint32) bool { return pkg_rhino.Node(rawNode_).IsCall() })
}

// Java polymorphic class interface: com.google.javascript.jscomp.CodeGenerator
type CodeGenerator interface {
	Implements_jscomp_CodeGenerator()
	IsJavaNil_() bool
	GetCodeGeneratorBase() *CodeGeneratorBase
	TagAsTypeSummary()
	TagAsStrict()
	PrintLeadingCommentsInOrder(pkg_rhino.Node)
	PrintTrailingComment(pkg_rhino.Node)
	Add_1_string(string)
	Add_1_Node(pkg_rhino.Node)
	Add_2(pkg_rhino.Node, CodeGenerator_Context)
	Add_3(pkg_rhino.Node, CodeGenerator_Context, bool)
	AddList_1(pkg_rhino.Node)
	AddList_2(pkg_rhino.Node, string)
	AddList_4(pkg_rhino.Node, bool, CodeGenerator_Context, string)
	AddStringKey(pkg_rhino.Node)
	AddObjectPattern(pkg_rhino.Node)
	AddArrayList(pkg_rhino.Node)
	AddCaseBody(pkg_rhino.Node)
	AddAllSiblings(pkg_rhino.Node)
	RegexpEscape(string) string
}

func ToIface_CodeGenerator[T interface {
	CodeGenerator
	comparable
}](v T) CodeGenerator {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.CodeGenerator
type CodeGeneratorBase struct {
	self_                   CodeGenerator
	cc                      CodeConsumer           `gson:"cc"`
	outputCharsetEncoder    *OutputCharsetEncoder  `gson:"outputCharsetEncoder"`
	outputFeatureSet        *pkg_parser.FeatureSet `gson:"outputFeatureSet"`
	jsDocInfoPrinter        *JSDocInfoPrinter      `gson:"jsDocInfoPrinter"`
	idHash_                 atomic.Int32
	preferSingleQuotes      bool `gson:"preferSingleQuotes"`
	preserveTypeAnnotations bool `gson:"preserveTypeAnnotations"`
	printNonJsDocComments   bool `gson:"printNonJsDocComments"`
	gentsMode               bool `gson:"gentsMode"`
	trustedStrings          bool `gson:"trustedStrings"`
	quoteKeywordProperties  bool `gson:"quoteKeywordProperties"`
	useOriginalName         bool `gson:"useOriginalName"`
}

func (this_ *CodeGeneratorBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(CodeGenerator)
}
func (this_ *CodeGeneratorBase) IsJavaNil_() bool { return this_ == nil }
func (*CodeGeneratorBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CodeGenerator", "CodeGenerator"
}
func (this_ *CodeGeneratorBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CodeGenerator", "com.google.javascript.jscomp.CodeGenerator", func(v_ any) bool {
		return j2gort.InstanceOf[CodeGenerator](v_)
	})
}
func (b *CodeGeneratorBase) GetCodeGeneratorBase() *CodeGeneratorBase { return b }
func (*CodeGeneratorBase) Implements_jscomp_CodeGenerator()           {}

func StaticGet_CodeGenerator_LT_ESCAPED() string { return "\\x3c" }
func StaticGet_CodeGenerator_GT_ESCAPED() string { return "\\x3e" }

var static_CodeGenerator_JSCOMP_SCOPE pkg_rhino.QualifiedName

func StaticGet_CodeGenerator_JSCOMP_SCOPE() pkg_rhino.QualifiedName {
	clinit_CodeGenerator()
	return static_CodeGenerator_JSCOMP_SCOPE
}
func StaticSet_CodeGenerator_JSCOMP_SCOPE(v pkg_rhino.QualifiedName) {
	clinit_CodeGenerator()
	static_CodeGenerator_JSCOMP_SCOPE = v
}

var clinitDone_CodeGenerator uint32
var clinitState_CodeGenerator j2gort.ClinitState

func clinit_CodeGenerator() {
	if atomic.LoadUint32(&clinitDone_CodeGenerator) != 0 {
		return
	}
	clinitSlow_CodeGenerator()
}

func clinitSlow_CodeGenerator() {
	if !j2gort.EnterClinit(&clinitState_CodeGenerator, &clinitDone_CodeGenerator) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CodeGenerator, &clinitDone_CodeGenerator)
	static_CodeGenerator_JSCOMP_SCOPE = pkg_rhino.QualifiedName_Of("$jscomp.scope")
}

func (this_ *CodeGeneratorBase) initFields_CodeGenerator() {
	_ = this_
}

// Java constructor: CodeGenerator(com.google.javascript.jscomp.CodeConsumer consumer)
func NewCodeGenerator_1(consumer CodeConsumer) CodeGenerator {
	this_ := &CodeGeneratorBase{}
	this_.BindSelf_(this_)
	InitCodeGenerator_1(this_, consumer)
	return this_
}

func InitCodeGenerator_1(this_ *CodeGeneratorBase, consumer CodeConsumer) {
	_ = consumer
	this_.initFields_CodeGenerator()
	this_.cc = consumer
	this_.outputCharsetEncoder = nil
	this_.preferSingleQuotes = false
	this_.trustedStrings = true
	this_.preserveTypeAnnotations = false
	this_.printNonJsDocComments = false
	this_.gentsMode = false
	this_.quoteKeywordProperties = false
	this_.useOriginalName = false
	this_.outputFeatureSet = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM()
	this_.jsDocInfoPrinter = NewJSDocInfoPrinter_1(false)
}

// Java constructor: CodeGenerator(com.google.javascript.jscomp.CodeConsumer consumer, com.google.javascript.jscomp.CompilerOptions options)
func NewCodeGenerator_2(consumer CodeConsumer, options *CompilerOptions) CodeGenerator {
	this_ := &CodeGeneratorBase{}
	this_.BindSelf_(this_)
	InitCodeGenerator_2(this_, consumer, options)
	return this_
}

func InitCodeGenerator_2(this_ *CodeGeneratorBase, consumer CodeConsumer, options *CompilerOptions) {
	_ = consumer
	_ = options
	this_.initFields_CodeGenerator()
	this_.cc = consumer
	this_.outputCharsetEncoder = NewOutputCharsetEncoder(options.GetOutputCharset())
	this_.preferSingleQuotes = options.ShouldPreferSingleQuotes()
	this_.trustedStrings = options.AssumeTrustedStrings()
	this_.preserveTypeAnnotations = options.ShouldPreserveTypeAnnotations()
	this_.printNonJsDocComments = options.GetPreserveNonJSDocComments()
	this_.gentsMode = options.GetGentsMode()
	this_.quoteKeywordProperties = options.ShouldQuoteKeywordProperties()
	this_.useOriginalName = options.GetUseOriginalNamesInOutput()
	this_.outputFeatureSet = options.GetOutputFeatureSet()
	this_.jsDocInfoPrinter = NewJSDocInfoPrinter_1(this_.useOriginalName)
}

// Java static method: CodeGenerator.forCostEstimation(com.google.javascript.jscomp.CodeConsumer consumer)
func CodeGenerator_ForCostEstimation(consumer CodeConsumer) CodeGenerator {
	_ = consumer
	return NewCodeGenerator_1(consumer)
}

// Java method: CodeGenerator.tagAsTypeSummary()
func (this_ *CodeGeneratorBase) TagAsTypeSummary() {
	this_.Add_1_string("/** @fileoverview @typeSummary */\n")
}

// Java method: CodeGenerator.tagAsStrict()
func (this_ *CodeGeneratorBase) TagAsStrict() {
	this_.Add_1_string("'use strict';")
	this_.cc.EndLine()
}

// Java method: CodeGenerator.printJsDocComment(com.google.javascript.rhino.Node node, com.google.javascript.rhino.JSDocInfo jsDocInfo)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_printJsDocComment_1306(node pkg_rhino.Node, jsDocInfo *pkg_rhino.JSDocInfo) {
	_ = node
	_ = jsDocInfo
	var jsdocAsString string = j2gort.NullString
	_ = jsdocAsString
	if this_.gentsMode {
		jsdocAsString = jsDocInfo.GetOriginalCommentString()
	} else {
		jsdocAsString = this_.jsDocInfoPrinter.Print(jsDocInfo)
	}
	if (jsdocAsString != j2gort.NullString) && !(j2gort.CheckNotNullString(jsdocAsString) == "/** */ ") {
		this_.Add_1_string(jsdocAsString)
		if !(node.IsCast()) {
			this_.cc.EndLine()
		}
	}
}

// Java method: CodeGenerator.printNonJsDocComment(com.google.javascript.rhino.Node node, com.google.javascript.rhino.NonJSDocComment nonJsDocComment)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_printNonJsDocComment_1307(node pkg_rhino.Node, nonJsDocComment *pkg_rhino.NonJSDocComment) {
	_ = node
	_ = nonJsDocComment
	var nonJsDocCommentString string = nonJsDocComment.GetCommentString()
	_ = nonJsDocCommentString
	if !(j2gort.CheckNotNullString(nonJsDocCommentString) == "") {
		this_.priv_CodeGenerator_addNonJsDocNonTrailing_1342(node, nonJsDocComment)
	}
}

// Java method: CodeGenerator.printLeadingCommentsInOrder(com.google.javascript.rhino.Node node)
func (this_ *CodeGeneratorBase) PrintLeadingCommentsInOrder(node pkg_rhino.Node) {
	_ = node
	var jsDocInfo *pkg_rhino.JSDocInfo = node.GetJSDocInfo()
	_ = jsDocInfo
	var nonJsDocComment *pkg_rhino.NonJSDocComment = node.GetNonJSDocComment()
	_ = nonJsDocComment
	var printJsDoc bool = (this_.preserveTypeAnnotations && (jsDocInfo != nil))
	_ = printJsDoc
	var printNonJsDoc bool = (this_.printNonJsDocComments && (nonJsDocComment != nil))
	_ = printNonJsDoc
	if printJsDoc && printNonJsDoc {
		if jsDocInfo.GetOriginalCommentPosition() < nonJsDocComment.GetStartPosition().GetOffset() {
			this_.priv_CodeGenerator_printJsDocComment_1306(node, jsDocInfo)
			this_.priv_CodeGenerator_printNonJsDocComment_1307(node, nonJsDocComment)
		} else {
			this_.priv_CodeGenerator_printNonJsDocComment_1307(node, nonJsDocComment)
			this_.priv_CodeGenerator_printJsDocComment_1306(node, jsDocInfo)
		}
		return
	}
	if printJsDoc {
		this_.priv_CodeGenerator_printJsDocComment_1306(node, jsDocInfo)
		return
	}
	if printNonJsDoc {
		this_.priv_CodeGenerator_printNonJsDocComment_1307(node, nonJsDocComment)
		return
	}
}

// Java method: CodeGenerator.hasTrailingCommentOnSameLine(com.google.javascript.rhino.Node node)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_hasTrailingCommentOnSameLine_1309(node pkg_rhino.Node) bool {
	_ = node
	if !(this_.printNonJsDocComments) {
		return false
	}
	return !(j2gort.CheckNotNullString(node.GetTrailingNonJSDocCommentString()) == "")
}

// Java method: CodeGenerator.printTrailingComment(com.google.javascript.rhino.Node node)
func (this_ *CodeGeneratorBase) PrintTrailingComment(node pkg_rhino.Node) {
	_ = node
	if !(this_.printNonJsDocComments) {
		return
	}
	var nonJsDocComment *pkg_rhino.NonJSDocComment = node.GetTrailingNonJSDocComment()
	_ = nonJsDocComment
	if nonJsDocComment != nil {
		var nonJsDocCommentString string = node.GetTrailingNonJSDocCommentString()
		_ = nonJsDocCommentString
		if !(j2gort.CheckNotNullString(nonJsDocCommentString) == "") {
			this_.priv_CodeGenerator_addNonJsDoctrailing_1343(nonJsDocComment, this_.priv_CodeGenerator_hasTrailingCommentOnSameLine_1309(node))
		}
	}
}

// Java method: CodeGenerator.add(java.lang.String str)
func (this_ *CodeGeneratorBase) Add_1_string(str string) {
	_ = str
	this_.cc.Add_1_string(str)
}

// Java method: CodeGenerator.add(com.google.javascript.rhino.Node n)
func (this_ *CodeGeneratorBase) Add_1_Node(n pkg_rhino.Node) {
	_ = n
	this_.self_.Add_2(n, StaticGet_CodeGenerator_Context_OTHER())
}

// Java method: CodeGenerator.add(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.CodeGenerator.Context context)
func (this_ *CodeGeneratorBase) Add_2(node pkg_rhino.Node, context CodeGenerator_Context) {
	_ = node
	_ = context
	this_.Add_3(node, context, true)
}

// Java method: CodeGenerator.add(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.CodeGenerator.Context context, boolean printComments)
func (this_ *CodeGeneratorBase) Add_3(node pkg_rhino.Node, context CodeGenerator_Context, printComments bool) {
	_ = node
	_ = context
	_ = printComments
	if !(this_.cc.ContinueProcessing()) {
		return
	}
	var shadow pkg_rhino.Node = node.GetClosureUnawareShadow()
	_ = shadow
	if shadow != 0 {
		var script pkg_rhino.Node = shadow.GetOnlyChild()
		_ = script
		var sinkCall pkg_rhino.Node = script.GetOnlyChild().GetOnlyChild()
		_ = sinkCall
		if !(sinkCall.IsCall()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(sinkCall))
		}
		this_.Add_3(sinkCall.GetLastChild(), context, printComments)
		return
	}
	if printComments {
		this_.PrintLeadingCommentsInOrder(node)
	}
	this_.cc.TrackLicenses(node)
	var type_ pkg_rhino.Token = node.GetToken()
	_ = type_
	var opstr string = NodeUtil_OpToStr(type_)
	_ = opstr
	var childCount int32 = node.GetChildCount()
	_ = childCount
	var first pkg_rhino.Node = node.GetFirstChild()
	_ = first
	var last pkg_rhino.Node = node.GetLastChild()
	_ = last
	if (opstr != j2gort.NullString) && (first != last) {
		if !(childCount == int32(2)) {
			j2gort.CheckStateFmt(false, "Bad binary operator \"%s\": expected 2 arguments but got %s", j2gort.BoxString(opstr), j2gort.ValueOfInt(childCount))
		}
		var p int32 = this_.priv_CodeGenerator_precedence_1317(node)
		_ = p
		var rhsContext CodeGenerator_Context = CodeGenerator_GetContextForNoInOperator(context)
		_ = rhsContext
		var needsParens bool = (((context == StaticGet_CodeGenerator_Context_START_OF_EXPR()) || context.AtArrowFunctionBody()) && first.IsObjectPattern())
		_ = needsParens
		if node.IsAssign() && needsParens {
			this_.Add_1_string("(")
		}
		if NodeUtil_IsAssignmentOp(node) || (type_ == pkg_rhino.StaticGet_Token_EXPONENT()) {
			this_.priv_CodeGenerator_addExpr_1328(first, (p + int32(1)), context)
			this_.cc.AddOp(opstr, true)
			this_.priv_CodeGenerator_addExpr_1328(last, p, rhsContext)
		} else {
			this_.priv_CodeGenerator_unrollBinaryOperator_1326(node, type_, opstr, context, rhsContext, p, (p + int32(1)))
		}
		if node.IsAssign() && needsParens {
			this_.Add_1_string(")")
		}
		return
	}
	this_.cc.StartSourceMapping(node)
	switch (type_).Ordinal() {
	case int32(48) /* Token.TRY */ :
		if !(first.GetNext().IsBlock() && !(first.GetNext().HasMoreThanOneChild())) {
			j2gort.CheckState(false)
		}
		if !((childCount >= int32(2)) && (childCount <= int32(3))) {
			j2gort.CheckState(false)
		}
		this_.Add_1_string("try")
		this_.Add_1_Node(first)
		var catchblock pkg_rhino.Node = first.GetNext().GetFirstChild()
		_ = catchblock
		if catchblock != 0 {
			this_.Add_1_Node(catchblock)
		}
		if childCount == int32(3) {
			this_.cc.MaybeInsertSpace()
			this_.Add_1_string("finally")
			this_.Add_1_Node(last)
		}
	case int32(86) /* Token.CATCH */ :
		if !(childCount == int32(2)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.cc.MaybeInsertSpace()
		this_.Add_1_string("catch")
		this_.cc.MaybeInsertSpace()
		if !(first.IsEmpty()) {
			this_.Add_1_string("(")
			this_.Add_1_Node(first)
			this_.Add_1_string(")")
		}
		this_.Add_1_Node(last)
	case int32(43) /* Token.THROW */ :
		if !(childCount == int32(1)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.Add_1_string("throw")
		this_.cc.MaybeInsertSpace()
		this_.Add_1_Node(first)
		this_.cc.EndStatement_2(true, this_.priv_CodeGenerator_hasTrailingCommentOnSameLine_1309(node))
	case int32(0) /* Token.RETURN */ :
		this_.Add_1_string("return")
		if childCount == int32(1) {
			this_.cc.MaybeInsertSpace()
			if this_.preserveTypeAnnotations && (first.GetJSDocInfo() != nil) {
				this_.Add_1_string("(")
				this_.Add_1_Node(first)
				this_.Add_1_string(")")
			} else {
				this_.Add_1_Node(first)
			}
		} else {
			if !(childCount == int32(0)) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
			}
		}
		this_.cc.EndStatement_1(this_.priv_CodeGenerator_hasTrailingCommentOnSameLine_1309(node))
	case int32(84) /* Token.VAR */ :
		this_.Add_1_string("var ")
		this_.AddList_4(first, false, CodeGenerator_GetContextForNoInOperator(context), ",")
		if (node.GetParent() == 0) || NodeUtil_IsStatement(node) {
			this_.cc.EndStatement_1(this_.priv_CodeGenerator_hasTrailingCommentOnSameLine_1309(node))
		}
	case int32(97) /* Token.CONST */ :
		this_.Add_1_string("const ")
		this_.AddList_4(first, false, CodeGenerator_GetContextForNoInOperator(context), ",")
		if (node.GetParent() == 0) || NodeUtil_IsStatement(node) {
			this_.cc.EndStatement_1(this_.priv_CodeGenerator_hasTrailingCommentOnSameLine_1309(node))
		}
	case int32(111) /* Token.LET */ :
		this_.Add_1_string("let ")
		this_.AddList_4(first, false, CodeGenerator_GetContextForNoInOperator(context), ",")
		if (node.GetParent() == 0) || NodeUtil_IsStatement(node) {
			this_.cc.EndStatement_1(this_.priv_CodeGenerator_hasTrailingCommentOnSameLine_1309(node))
		}
	case int32(99) /* Token.LABEL_NAME */ :
		if !(!(j2gort.CheckNotNullString(node.GetString()) == "")) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.priv_CodeGenerator_addIdentifier_1315(node.GetString())
	case int32(104) /* Token.DESTRUCTURING_LHS */ :
		this_.Add_1_Node(first)
		if first != last {
			if !(childCount == int32(2)) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
			}
			this_.cc.AddOp("=", true)
			this_.priv_CodeGenerator_addExpr_1328(last, NodeUtil_Precedence(pkg_rhino.StaticGet_Token_ASSIGN()), CodeGenerator_GetContextForNoInOperator(context))
		}
	case int32(32) /* Token.NAME */ :
		if this_.useOriginalName && (node.GetOriginalName() != j2gort.NullString) {
			this_.priv_CodeGenerator_addIdentifier_1315(node.GetOriginalName())
		} else {
			this_.priv_CodeGenerator_addIdentifier_1315(node.GetString())
		}
		this_.priv_CodeGenerator_maybeAddOptional_1325(node)
		this_.priv_CodeGenerator_maybeAddTypeDecl_1323(node)
		if (first != 0) && !(first.IsEmpty()) {
			if !(childCount == int32(1)) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
			}
			this_.cc.AddOp("=", true)
			this_.priv_CodeGenerator_addExpr_1328(first, NodeUtil_Precedence(pkg_rhino.StaticGet_Token_ASSIGN()), CodeGenerator_GetContextForNoInOperator(context))
		}
	case int32(46) /* Token.ARRAYLIT */ :
		this_.Add_1_string("[")
		this_.AddArrayList(first)
		this_.Add_1_string("]")
	case int32(102) /* Token.ARRAY_PATTERN */ :
		this_.Add_1_string("[")
		this_.AddArrayList(first)
		this_.Add_1_string("]")
		this_.priv_CodeGenerator_maybeAddTypeDecl_1323(node)
	case int32(49) /* Token.PARAM_LIST */ :
		if ((node.GetParent().IsArrowFunction() && node.HasOneChild()) && first.IsName()) && !(this_.gentsMode) {
			this_.Add_1_Node(first)
		} else {
			this_.Add_1_string("(")
			this_.AddList_1(first)
			this_.Add_1_string(")")
		}
	case int32(134) /* Token.DEFAULT_VALUE */ :
		this_.Add_1_Node(first)
		this_.priv_CodeGenerator_maybeAddTypeDecl_1323(node)
		this_.cc.AddOp("=", true)
		this_.priv_CodeGenerator_addExpr_1328(first.GetNext(), int32(1), StaticGet_CodeGenerator_Context_OTHER())
	case int32(50) /* Token.COMMA */ :
		if !(childCount == int32(2)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.priv_CodeGenerator_unrollBinaryOperator_1326(node, pkg_rhino.StaticGet_Token_COMMA(), ",", context, CodeGenerator_GetContextForNoInOperator(context), int32(0), int32(0))
	case int32(33) /* Token.NUMBER */ :
		if !(childCount == int32(0)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.cc.AddNumber(node.GetDouble(), node)
	case int32(34) /* Token.BIGINT */ :
		if !(childCount == int32(0)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.cc.AddBigInt(node.GetBigInt())
	case int32(25) /* Token.TYPEOF */, int32(87) /* Token.VOID */, int32(19) /* Token.NOT */, int32(20) /* Token.BITNOT */, int32(21) /* Token.POS */, int32(22) /* Token.NEG */ :
		if !(childCount == int32(1)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.cc.AddOp(NodeUtil_OpToStrNoFail(type_), false)
		this_.priv_CodeGenerator_addExpr_1328(first, NodeUtil_Precedence(type_), StaticGet_CodeGenerator_Context_OTHER())
	case int32(67) /* Token.HOOK */ :
		if !(childCount == int32(3)) {
			j2gort.CheckStateFmt(false, "%s wrong number of children: %s", j2gort.BoxNode(node), j2gort.ValueOfInt(childCount))
		}
		var p_2 int32 = NodeUtil_Precedence(type_)
		_ = p_2
		var rhsContext_2 CodeGenerator_Context = CodeGenerator_GetContextForNoInOperator(context)
		_ = rhsContext_2
		this_.priv_CodeGenerator_addExpr_1328(first, (p_2 + int32(1)), context)
		this_.cc.AddOp("?", true)
		this_.priv_CodeGenerator_addExpr_1328(first.GetNext(), int32(1), rhsContext_2)
		this_.cc.AddOp(":", true)
		this_.priv_CodeGenerator_addExpr_1328(last, int32(1), rhsContext_2)
	case int32(42) /* Token.REGEXP */ :
		if !(first.IsStringLit()) || !(last.IsStringLit()) {
			panic(j2gort.AsIface[javalang.IError](j2gort.ConstructShimThrowable("java.lang.Error", "Expected children to be strings")))
		}
		var regexp string = this_.RegexpEscape(first.GetString())
		_ = regexp
		if childCount == int32(2) {
			this_.Add_1_string((j2gort.StrValString(regexp) + j2gort.StrValString(last.GetString())))
		} else {
			if !(childCount == int32(1)) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
			}
			this_.Add_1_string(regexp)
		}
	case int32(73) /* Token.FUNCTION */ :
		if j2gort.GetObjectClass(j2gort.BoxNode(node)) != j2gort.ClassLiteral("com.google.javascript.rhino.Node", "Node") {
			panic(j2gort.AsIface[javalang.IError](j2gort.ConstructShimThrowable("java.lang.Error", "Unexpected Node subclass.")))
		}
		if !(childCount == int32(3)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		if node.IsArrowFunction() {
			this_.priv_CodeGenerator_addArrowFunction_1320(node, first, last, context)
		} else {
			this_.priv_CodeGenerator_addFunction_1321(node, first, last, context)
		}
	case int32(125) /* Token.ITER_REST */, int32(126) /* Token.OBJECT_REST */ :
		this_.Add_1_string("...")
		this_.Add_1_Node(first)
		this_.priv_CodeGenerator_maybeAddTypeDecl_1323(node)
	case int32(127) /* Token.ITER_SPREAD */, int32(128) /* Token.OBJECT_SPREAD */ :
		this_.Add_1_string("...")
		this_.priv_CodeGenerator_addExpr_1328(first, NodeUtil_Precedence(type_), StaticGet_CodeGenerator_Context_OTHER())
	case int32(120) /* Token.EXPORT */ :
		this_.Add_1_string("export")
		if node.GetBooleanProp(pkg_rhino.StaticGet_Node_EXPORT_DEFAULT()) {
			this_.Add_1_string("default")
		}
		if node.GetBooleanProp(pkg_rhino.StaticGet_Node_EXPORT_ALL_FROM()) {
			this_.Add_1_string("*")
			if !((first != 0) && first.IsEmpty()) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
			}
		} else {
			this_.Add_1_Node(first)
		}
		if childCount == int32(2) {
			this_.Add_1_string("from")
			this_.Add_1_Node(last)
		}
		this_.priv_CodeGenerator_processEnd_1353(first, context)
	case int32(116) /* Token.IMPORT */ :
		this_.Add_1_string("import")
		var second pkg_rhino.Node = first.GetNext()
		_ = second
		if !(first.IsEmpty()) {
			this_.Add_1_Node(first)
			if !(second.IsEmpty()) {
				this_.cc.ListSeparator()
			}
		}
		if !(second.IsEmpty()) {
			this_.Add_1_Node(second)
		}
		if !(first.IsEmpty()) || !(second.IsEmpty()) {
			this_.Add_1_string("from")
		}
		this_.Add_1_Node(last)
		this_.cc.EndStatement_1(this_.priv_CodeGenerator_hasTrailingCommentOnSameLine_1309(node))
	case int32(121) /* Token.EXPORT_SPECS */, int32(117) /* Token.IMPORT_SPECS */ :
		this_.Add_1_string("{")
		{
			var c pkg_rhino.Node = first
			_ = c
			for ; c != 0; c = c.GetNext() {
				if c != first {
					this_.cc.ListSeparator()
				}
				this_.Add_1_Node(c)
			}
		}
		this_.Add_1_string("}")
	case int32(122) /* Token.EXPORT_SPEC */, int32(118) /* Token.IMPORT_SPEC */ :
		this_.Add_1_Node(first)
		if node.IsShorthandProperty() && func() bool {
			sr_1 := first.GetString()
			se_2 := last.GetString()
			return j2gort.CheckNotNullString(sr_1) == se_2
		}() {
			break
		}
		this_.Add_1_string("as")
		this_.Add_1_Node(last)
	case int32(119) /* Token.IMPORT_STAR */ :
		this_.Add_1_string("*")
		this_.Add_1_string("as")
		this_.Add_1_string(node.GetString())
	case int32(124) /* Token.DYNAMIC_IMPORT */ :
		this_.Add_1_string("import(")
		this_.priv_CodeGenerator_addExpr_1328(first, NodeUtil_Precedence(type_), context)
		this_.Add_1_string(")")
	case int32(136) /* Token.IMPORT_META */ :
		this_.Add_1_string("import.meta")
	case int32(105) /* Token.CLASS */ :
		{
			if !(childCount == int32(3)) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
			}
			var classNeedsParens bool = (context == StaticGet_CodeGenerator_Context_START_OF_EXPR())
			_ = classNeedsParens
			if classNeedsParens {
				this_.Add_1_string("(")
			}
			var name pkg_rhino.Node = first
			_ = name
			var superClass pkg_rhino.Node = first.GetNext()
			_ = superClass
			var members pkg_rhino.Node = last
			_ = members
			this_.Add_1_string("class")
			if !(name.IsEmpty()) {
				this_.Add_1_Node(name)
			}
			this_.priv_CodeGenerator_maybeAddGenericTypes_1324(first)
			if !(superClass.IsEmpty()) {
				this_.Add_1_string("extends")
				this_.priv_CodeGenerator_addExpr_1328(superClass, int32(1), StaticGet_CodeGenerator_Context_OTHER())
			}
			var interfaces pkg_rhino.Node = node.GetImplements()
			_ = interfaces
			if interfaces != 0 {
				this_.Add_1_string("implements")
				var child pkg_rhino.Node = interfaces.GetFirstChild()
				_ = child
				this_.Add_1_Node(child)
				for (func() pkg_rhino.Node { var asgn_3 pkg_rhino.Node = child.GetNext(); child = asgn_3; return asgn_3 }()) != 0 {
					this_.Add_1_string(",")
					this_.cc.MaybeInsertSpace()
					this_.Add_1_Node(child)
				}
			}
			this_.Add_1_Node(members)
			this_.cc.EndClass((context == StaticGet_CodeGenerator_Context_STATEMENT()))
			if classNeedsParens {
				this_.Add_1_string(")")
			}
		}
	case int32(106) /* Token.CLASS_MEMBERS */, int32(166) /* Token.INTERFACE_MEMBERS */, int32(176) /* Token.NAMESPACE_ELEMENTS */ :
		this_.cc.BeginBlock()
		{
			var c_2 pkg_rhino.Node = first
			_ = c_2
			for ; c_2 != 0; c_2 = c_2.GetNext() {
				this_.Add_1_Node(c_2)
				this_.priv_CodeGenerator_processEnd_1353(c_2, context)
				this_.cc.EndLine()
			}
		}
		this_.cc.EndBlock_1(false)
	case int32(168) /* Token.ENUM_MEMBERS */ :
		this_.cc.BeginBlock()
		{
			var c_3 pkg_rhino.Node = first
			_ = c_3
			for ; c_3 != 0; c_3 = c_3.GetNext() {
				this_.Add_1_Node(c_3)
				if c_3.GetNext() != 0 {
					this_.Add_1_string(",")
				}
				this_.cc.EndLine()
			}
		}
		this_.cc.EndBlock_1(false)
	case int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */, int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(172) /* Token.MEMBER_VARIABLE_DEF */ :
		if !((((node.GetParent().IsObjectLit() || node.GetParent().IsClassMembers()) || node.GetParent().IsInterfaceMembers()) || node.GetParent().IsRecordType()) || node.GetParent().IsIndexSignature()) {
			j2gort.CheckState(false)
		}
		this_.priv_CodeGenerator_maybeAddAccessibilityModifier_1322(node)
		if node.IsStaticMember() {
			this_.Add_1_string("static ")
		}
		if node.IsMemberFunctionDef() && node.GetFirstChild().IsAsyncFunction() {
			this_.Add_1_string("async ")
		}
		if !(node.IsMemberVariableDef()) && node.GetFirstChild().IsGeneratorFunction() {
			if !(type_ == pkg_rhino.StaticGet_Token_MEMBER_FUNCTION_DEF()) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
			}
			this_.Add_1_string("*")
		}
		switch (type_).Ordinal() {
		case int32(95) /* Token.GETTER_DEF */ :
			if !(!(first.GetSecondChild().HasChildren())) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
			}
			this_.Add_1_string("get ")
		case int32(96) /* Token.SETTER_DEF */ :
			if !(first.GetSecondChild().HasOneChild()) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
			}
			this_.Add_1_string("set ")
		case int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(172) /* Token.MEMBER_VARIABLE_DEF */ :
		default:
		}
		var name_2 string = node.GetString()
		_ = name_2
		if node.IsMemberVariableDef() {
			this_.Add_1_string(node.GetString())
			this_.priv_CodeGenerator_maybeAddOptional_1325(node)
			this_.priv_CodeGenerator_maybeAddTypeDecl_1323(node)
		} else {
			if !(childCount == int32(1)) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
			}
			if !(first.IsFunction()) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(first))
			}
			if !(j2gort.CheckNotNullString(first.GetFirstChild().GetString()) == "") {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(first))
			}
			var fn pkg_rhino.Node = first
			_ = fn
			var parameters pkg_rhino.Node = fn.GetSecondChild()
			_ = parameters
			var body pkg_rhino.Node = fn.GetLastChild()
			_ = body
			if node.IsPrivateIdentifier() {
				this_.priv_CodeGenerator_addIdentifier_1315(name_2)
				this_.priv_CodeGenerator_maybeAddGenericTypes_1324(fn.GetFirstChild())
			} else if (!(node.IsQuotedStringKey()) && pkg_rhino.TokenStream_IsJSIdentifier(name_2)) && NodeUtil_IsLatin(name_2) {
				this_.Add_1_string(name_2)
				this_.priv_CodeGenerator_maybeAddGenericTypes_1324(fn.GetFirstChild())
			} else {
				var d float64 = CodeGenerator_GetSimpleNumber(name_2)
				_ = d
				if !(math.IsNaN(d)) {
					this_.cc.AddNumber(d, node)
				} else {
					this_.priv_CodeGenerator_addJsString_1344(node)
				}
			}
			this_.priv_CodeGenerator_maybeAddOptional_1325(fn)
			this_.Add_1_Node(parameters)
			this_.priv_CodeGenerator_maybeAddTypeDecl_1323(fn)
			this_.Add_1_Node(body)
		}
	case int32(108) /* Token.MEMBER_FIELD_DEF */, int32(109) /* Token.COMPUTED_FIELD_DEF */ :
		if !(node.GetParent().IsClassMembers()) {
			j2gort.CheckState(false)
		}
		if node.GetBooleanProp(pkg_rhino.StaticGet_Node_STATIC_MEMBER()) {
			this_.Add_1_string("static ")
		}
		var init pkg_rhino.Node = pkg_rhino.Node(0)
		_ = init
		switch (type_).Ordinal() {
		case int32(108) /* Token.MEMBER_FIELD_DEF */ :
			this_.priv_CodeGenerator_addIdentifier_1315(node.GetString())
			init = first
		case int32(109) /* Token.COMPUTED_FIELD_DEF */ :
			this_.Add_1_string("[")
			this_.priv_CodeGenerator_addExpr_1328(first, int32(1), StaticGet_CodeGenerator_Context_OTHER())
			this_.Add_1_string("]")
			init = node.GetSecondChild()
		default:
		}
		if init != 0 {
			this_.cc.AddOp("=", true)
			this_.priv_CodeGenerator_addExpr_1328(init, int32(1), StaticGet_CodeGenerator_Context_OTHER())
		}
		this_.Add_1_string(";")
	case int32(94) /* Token.SCRIPT */, int32(123) /* Token.MODULE_BODY */, int32(90) /* Token.BLOCK */, int32(89) /* Token.ROOT */ :
		if j2gort.GetObjectClass(j2gort.BoxNode(node)) != j2gort.ClassLiteral("com.google.javascript.rhino.Node", "Node") {
			panic(j2gort.AsIface[javalang.IError](j2gort.ConstructShimThrowable("java.lang.Error", "Unexpected Node subclass.")))
		}
		if node.HasParent() {
			var staticBlock bool = (node.IsBlock() && node.GetParent().IsClassMembers())
			_ = staticBlock
			if staticBlock {
				this_.Add_1_string("static")
			}
		}
		var preserveBlock bool = (node.IsBlock() && !(node.IsSyntheticBlock()))
		_ = preserveBlock
		if preserveBlock {
			this_.cc.BeginBlock()
		}
		var preferLineBreaks bool = ((type_ == pkg_rhino.StaticGet_Token_SCRIPT()) || (((type_ == pkg_rhino.StaticGet_Token_BLOCK()) && !(preserveBlock)) && node.GetParent().IsScript()))
		_ = preferLineBreaks
		{
			var c_4 pkg_rhino.Node = first
			_ = c_4
			for ; c_4 != 0; c_4 = c_4.GetNext() {
				this_.self_.Add_2(c_4, StaticGet_CodeGenerator_Context_STATEMENT())
				if c_4.IsFunction() || c_4.IsClass() {
					this_.cc.MaybeLineBreak()
				}
				if preferLineBreaks {
					this_.cc.NotePreferredLineBreak()
				}
			}
		}
		if preserveBlock {
			this_.cc.EndBlock_1(this_.cc.BreakAfterBlockFor(node, (context == StaticGet_CodeGenerator_Context_STATEMENT())))
		}
	case int32(80) /* Token.FOR */ :
		if !(childCount == int32(4)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.Add_1_string("for")
		this_.cc.MaybeInsertSpace()
		this_.Add_1_string("(")
		if NodeUtil_IsNameDeclaration(first) {
			this_.self_.Add_2(first, StaticGet_CodeGenerator_Context_IN_FOR_INIT_CLAUSE())
		} else {
			this_.priv_CodeGenerator_addExpr_1328(first, int32(0), StaticGet_CodeGenerator_Context_IN_FOR_INIT_CLAUSE())
		}
		this_.Add_1_string(";")
		if !(first.GetNext().IsEmpty()) {
			this_.cc.MaybeInsertSpace()
		}
		this_.Add_1_Node(first.GetNext())
		this_.Add_1_string(";")
		if !(first.GetNext().GetNext().IsEmpty()) {
			this_.cc.MaybeInsertSpace()
		}
		this_.Add_1_Node(first.GetNext().GetNext())
		this_.Add_1_string(")")
		this_.priv_CodeGenerator_addNonEmptyStatement_1327(last, CodeGenerator_GetContextForNonEmptyExpression(context), false)
	case int32(81) /* Token.FOR_IN */ :
		if !(childCount == int32(3)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.Add_1_string("for")
		this_.cc.MaybeInsertSpace()
		this_.Add_1_string("(")
		this_.Add_1_Node(first)
		this_.Add_1_string("in")
		this_.Add_1_Node(first.GetNext())
		this_.Add_1_string(")")
		this_.priv_CodeGenerator_addNonEmptyStatement_1327(last, CodeGenerator_GetContextForNonEmptyExpression(context), false)
	case int32(112) /* Token.FOR_OF */ :
		if !(childCount == int32(3)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.Add_1_string("for")
		this_.cc.MaybeInsertSpace()
		this_.Add_1_string("(")
		this_.Add_1_Node(first)
		this_.cc.MaybeInsertSpace()
		this_.Add_1_string("of")
		this_.cc.MaybeInsertSpace()
		this_.priv_CodeGenerator_addExpr_1328(first.GetNext(), NodeUtil_Precedence(pkg_rhino.StaticGet_Token_ASSIGN()), StaticGet_CodeGenerator_Context_OTHER())
		this_.Add_1_string(")")
		this_.priv_CodeGenerator_addNonEmptyStatement_1327(last, CodeGenerator_GetContextForNonEmptyExpression(context), false)
	case int32(113) /* Token.FOR_AWAIT_OF */ :
		if !(childCount == int32(3)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.Add_1_string("for await")
		this_.cc.MaybeInsertSpace()
		this_.Add_1_string("(")
		this_.Add_1_Node(first)
		this_.cc.MaybeInsertSpace()
		this_.Add_1_string("of")
		this_.cc.MaybeInsertSpace()
		this_.priv_CodeGenerator_addExpr_1328(first.GetNext(), NodeUtil_Precedence(pkg_rhino.StaticGet_Token_ASSIGN()), StaticGet_CodeGenerator_Context_OTHER())
		this_.Add_1_string(")")
		this_.priv_CodeGenerator_addNonEmptyStatement_1327(last, CodeGenerator_GetContextForNonEmptyExpression(context), false)
	case int32(79) /* Token.DO */ :
		if !(childCount == int32(2)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.Add_1_string("do")
		this_.priv_CodeGenerator_addNonEmptyStatement_1327(first, StaticGet_CodeGenerator_Context_OTHER(), false)
		this_.cc.MaybeInsertSpace()
		this_.Add_1_string("while")
		this_.cc.MaybeInsertSpace()
		this_.Add_1_string("(")
		this_.Add_1_Node(last)
		this_.Add_1_string(")")
		this_.cc.EndStatement_1(this_.priv_CodeGenerator_hasTrailingCommentOnSameLine_1309(node))
	case int32(78) /* Token.WHILE */ :
		if !(childCount == int32(2)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.Add_1_string("while")
		this_.cc.MaybeInsertSpace()
		this_.Add_1_string("(")
		this_.Add_1_Node(first)
		this_.Add_1_string(")")
		this_.priv_CodeGenerator_addNonEmptyStatement_1327(last, CodeGenerator_GetContextForNonEmptyExpression(context), false)
	case int32(88) /* Token.EMPTY */ :
		if !(childCount == int32(0)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
	case int32(29) /* Token.OPTCHAIN_GETPROP */ :
		this_.priv_CodeGenerator_addExpr_1328(first, NodeUtil_Precedence(type_), context)
		this_.Add_1_string(func() string {
			if node.IsOptionalChainStart() {
				return "?."
			}
			return "."
		}())
		this_.priv_CodeGenerator_addGetpropIdentifier_1316(node)
	case int32(26) /* Token.GETPROP */ :
		if this_.useOriginalName && (node.GetOriginalName() != j2gort.NullString) {
			if StaticGet_CodeGenerator_JSCOMP_SCOPE().Matches(node.GetFirstChild()) && node.GetParent().IsAssign() {
				this_.Add_1_string("var ")
			}
			this_.priv_CodeGenerator_addGetpropIdentifier_1316(node)
			break
		}
		var breakOutOfOptionalChain bool = NodeUtil_IsOptChainNode(first)
		_ = breakOutOfOptionalChain
		var needsParens_2 bool = (first.IsNumber() || breakOutOfOptionalChain)
		_ = needsParens_2
		if needsParens_2 {
			this_.Add_1_string("(")
		}
		this_.priv_CodeGenerator_addExpr_1328(first, NodeUtil_Precedence(type_), context)
		if needsParens_2 {
			this_.Add_1_string(")")
		}
		if this_.quoteKeywordProperties && pkg_rhino.TokenStream_IsKeyword(node.GetString()) {
			this_.Add_1_string(("[\"" + j2gort.StrValString(node.GetString()) + "\"]"))
		} else {
			this_.Add_1_string(".")
			this_.priv_CodeGenerator_addGetpropIdentifier_1316(node)
		}
	case int32(30) /* Token.OPTCHAIN_GETELEM */ :
		if !(childCount == int32(2)) {
			j2gort.CheckStateFmt(false, "Bad GETELEM node: Expected 2 children but got %s. For node: %s", j2gort.ValueOfInt(childCount), j2gort.BoxNode(node))
		}
		this_.priv_CodeGenerator_addExpr_1328(first, NodeUtil_Precedence(type_), context)
		if node.IsOptionalChainStart() {
			this_.Add_1_string("?.")
		}
		this_.Add_1_string("[")
		this_.Add_1_Node(first.GetNext())
		this_.Add_1_string("]")
	case int32(27) /* Token.GETELEM */ :
		if !(childCount == int32(2)) {
			j2gort.CheckStateFmt(false, "Bad GETELEM node: Expected 2 children but got %s. For node: %s", j2gort.ValueOfInt(childCount), j2gort.BoxNode(node))
		}
		var needsParens_3 bool = NodeUtil_IsOptChainNode(first)
		_ = needsParens_3
		if needsParens_3 {
			this_.Add_1_string("(")
		}
		this_.priv_CodeGenerator_addExpr_1328(first, NodeUtil_Precedence(type_), context)
		if needsParens_3 {
			this_.Add_1_string(")")
		}
		this_.Add_1_string("[")
		this_.Add_1_Node(first.GetNext())
		this_.Add_1_string("]")
	case int32(85) /* Token.WITH */ :
		if !(childCount == int32(2)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.Add_1_string("with(")
		this_.Add_1_Node(first)
		this_.Add_1_string(")")
		this_.priv_CodeGenerator_addNonEmptyStatement_1327(last, CodeGenerator_GetContextForNonEmptyExpression(context), false)
	case int32(71) /* Token.INC */, int32(72) /* Token.DEC */ :
		if !(childCount == int32(1)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		var o string
		_ = o
		if type_ == pkg_rhino.StaticGet_Token_INC() {
			o = "++"
		} else {
			o = "--"
		}
		var postProp bool = node.GetBooleanProp(pkg_rhino.StaticGet_Node_INCRDECR_PROP())
		_ = postProp
		if postProp {
			this_.priv_CodeGenerator_addExpr_1328(first, NodeUtil_Precedence(type_), context)
			this_.cc.AddOp(o, false)
		} else {
			this_.cc.AddOp(o, false)
			this_.Add_1_Node(first)
		}
	case int32(31) /* Token.OPTCHAIN_CALL */ :
		if CodeGenerator_IsIndirectEval(first) || (node.GetBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL()) && NodeUtil_IsNormalOrOptChainGet(first)) {
			this_.Add_1_string("(0,")
			this_.priv_CodeGenerator_addExpr_1328(first, NodeUtil_Precedence(pkg_rhino.StaticGet_Token_COMMA()), StaticGet_CodeGenerator_Context_OTHER())
			this_.Add_1_string(")")
		} else {
			this_.priv_CodeGenerator_addExpr_1328(first, NodeUtil_Precedence(type_), context)
		}
		var args pkg_rhino.Node = first.GetNext()
		_ = args
		if node.IsOptionalChainStart() {
			this_.Add_1_string("?.")
		}
		this_.Add_1_string("(")
		this_.AddList_1(args)
		this_.Add_1_string(")")
	case int32(28) /* Token.CALL */ :
		this_.priv_CodeGenerator_addInvocationTarget_1318(node, context)
		this_.Add_1_string("(")
		this_.AddList_1(first.GetNext())
		this_.Add_1_string(")")
	case int32(74) /* Token.IF */ :
		if !((childCount == int32(2)) || (childCount == int32(3))) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		var hasElse bool = (childCount == int32(3))
		_ = hasElse
		var ambiguousElseClause bool = ((context == StaticGet_CodeGenerator_Context_BEFORE_DANGLING_ELSE()) && !(hasElse))
		_ = ambiguousElseClause
		if ambiguousElseClause {
			this_.cc.BeginBlock()
		}
		this_.Add_1_string("if")
		this_.cc.MaybeInsertSpace()
		this_.Add_1_string("(")
		this_.Add_1_Node(first)
		this_.Add_1_string(")")
		if hasElse {
			this_.priv_CodeGenerator_addNonEmptyStatement_1327(first.GetNext(), StaticGet_CodeGenerator_Context_BEFORE_DANGLING_ELSE(), false)
			this_.cc.MaybeInsertSpace()
			this_.Add_1_string("else")
			this_.priv_CodeGenerator_addNonEmptyStatement_1327(last, CodeGenerator_GetContextForNonEmptyExpression(context), false)
		} else {
			this_.priv_CodeGenerator_addNonEmptyStatement_1327(first.GetNext(), StaticGet_CodeGenerator_Context_OTHER(), false)
		}
		if ambiguousElseClause {
			this_.cc.EndBlock_0()
		}
	case int32(36) /* Token.NULL */ :
		if !(childCount == int32(0)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.cc.AddConstant("null")
	case int32(37) /* Token.THIS */ :
		if !(childCount == int32(0)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.Add_1_string("this")
	case int32(110) /* Token.SUPER */ :
		if !(childCount == int32(0)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.Add_1_string("super")
	case int32(135) /* Token.NEW_TARGET */ :
		if !(childCount == int32(0)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.Add_1_string("new.target")
	case int32(114) /* Token.YIELD */ :
		this_.Add_1_string("yield")
		if node.IsYieldAll() {
			j2gort.CheckNotNullNode(first)
			this_.Add_1_string("*")
		}
		if first != 0 {
			this_.cc.MaybeInsertSpace()
			this_.priv_CodeGenerator_addExpr_1328(first, NodeUtil_Precedence(type_), StaticGet_CodeGenerator_Context_OTHER())
		}
	case int32(115) /* Token.AWAIT */ :
		this_.Add_1_string("await ")
		this_.priv_CodeGenerator_addExpr_1328(first, NodeUtil_Precedence(type_), StaticGet_CodeGenerator_Context_OTHER())
	case int32(38) /* Token.FALSE */ :
		if !(childCount == int32(0)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.cc.AddConstant("false")
	case int32(39) /* Token.TRUE */ :
		if !(childCount == int32(0)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.cc.AddConstant("true")
	case int32(83) /* Token.CONTINUE */ :
		if !(childCount <= int32(1)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.Add_1_string("continue")
		if childCount == int32(1) {
			if !(first.IsLabelName()) {
				panic(j2gort.AsIface[javalang.IError](j2gort.ConstructShimThrowable("java.lang.Error", "Unexpected token type. Should be LABEL_NAME.")))
			}
			this_.Add_1_string(" ")
			this_.Add_1_Node(first)
		}
		this_.cc.EndStatement_1(this_.priv_CodeGenerator_hasTrailingCommentOnSameLine_1309(node))
	case int32(98) /* Token.DEBUGGER */ :
		if !(childCount == int32(0)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.Add_1_string("debugger")
		this_.cc.EndStatement_1(this_.priv_CodeGenerator_hasTrailingCommentOnSameLine_1309(node))
	case int32(82) /* Token.BREAK */ :
		if !(childCount <= int32(1)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.Add_1_string("break")
		if childCount == int32(1) {
			if !(first.IsLabelName()) {
				panic(j2gort.AsIface[javalang.IError](j2gort.ConstructShimThrowable("java.lang.Error", "Unexpected token type. Should be LABEL_NAME.")))
			}
			this_.Add_1_string(" ")
			this_.Add_1_Node(first)
		}
		this_.cc.EndStatement_1(this_.priv_CodeGenerator_hasTrailingCommentOnSameLine_1309(node))
	case int32(93) /* Token.EXPR_RESULT */ :
		if !(childCount == int32(1)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.self_.Add_2(first, StaticGet_CodeGenerator_Context_START_OF_EXPR())
		this_.cc.EndStatement_1(this_.priv_CodeGenerator_hasTrailingCommentOnSameLine_1309(node))
	case int32(23) /* Token.NEW */ :
		this_.Add_1_string("new ")
		var precedence int32 = NodeUtil_Precedence(type_)
		_ = precedence
		var precedenceOfFirst int32 = NodeUtil_Precedence(first.GetToken())
		_ = precedenceOfFirst
		if precedenceOfFirst == precedence {
			precedence = (precedence + int32(1))
		}
		if NodeUtil_Has(first, lambdaSingleton_code_generator_1, StaticGet_NodeUtil_MATCH_NOT_FUNCTION()) || NodeUtil_IsOptChainNode(first) {
			precedence = (NodeUtil_Precedence(first.GetToken()) + int32(1))
		}
		this_.priv_CodeGenerator_addExpr_1328(first, precedence, StaticGet_CodeGenerator_Context_OTHER())
		var next pkg_rhino.Node = first.GetNext()
		_ = next
		if next != 0 {
			this_.Add_1_string("(")
			this_.AddList_1(next)
			this_.Add_1_string(")")
		} else {
			if this_.cc.ShouldPreserveExtras(node) {
				this_.Add_1_string("(")
				this_.Add_1_string(")")
			}
		}
	case int32(100) /* Token.STRING_KEY */ :
		this_.AddStringKey(node)
	case int32(35) /* Token.STRINGLIT */ :
		if !(childCount == int32(0)) {
			j2gort.CheckStateFmt(false, "String node %s may not have children", j2gort.BoxNode(node))
		}
		this_.priv_CodeGenerator_addJsString_1344(node)
	case int32(24) /* Token.DELPROP */ :
		if !(childCount == int32(1)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.Add_1_string("delete ")
		this_.Add_1_Node(first)
	case int32(47) /* Token.OBJECTLIT */ :
		var needsParens_4 bool = ((context == StaticGet_CodeGenerator_Context_START_OF_EXPR()) || context.AtArrowFunctionBody())
		_ = needsParens_4
		if needsParens_4 {
			this_.Add_1_string("(")
		}
		this_.Add_1_string("{")
		{
			var c_5 pkg_rhino.Node = first
			_ = c_5
			for ; c_5 != 0; c_5 = c_5.GetNext() {
				if c_5 != first {
					this_.cc.ListSeparator()
				}
				if !(NodeUtil_IsObjLitProperty(c_5) || c_5.IsSpread()) {
					j2gort.CheckStateMsg(false, j2gort.BoxNode(c_5))
				}
				this_.Add_1_Node(c_5)
			}
		}
		if (first != 0) && node.HasTrailingComma() {
			this_.cc.OptionalListSeparator()
		}
		this_.Add_1_string("}")
		if needsParens_4 {
			this_.Add_1_string(")")
		}
	case int32(129) /* Token.COMPUTED_PROP */ :
		this_.priv_CodeGenerator_maybeAddAccessibilityModifier_1322(node)
		if node.GetBooleanProp(pkg_rhino.StaticGet_Node_STATIC_MEMBER()) {
			this_.Add_1_string("static ")
		}
		if node.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_GETTER()) {
			this_.Add_1_string("get ")
		} else if node.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_SETTER()) {
			this_.Add_1_string("set ")
		} else if node.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_METHOD()) {
			if last.IsAsyncFunction() {
				this_.Add_1_string("async")
			}
			if last.GetBooleanProp(pkg_rhino.StaticGet_Node_GENERATOR_FN()) {
				this_.Add_1_string("*")
			}
		}
		this_.Add_1_string("[")
		this_.priv_CodeGenerator_addExpr_1328(first, int32(1), StaticGet_CodeGenerator_Context_OTHER())
		this_.Add_1_string("]")
		this_.priv_CodeGenerator_maybeAddTypeDecl_1323(node)
		if (node.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_METHOD()) || node.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_GETTER())) || node.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_SETTER()) {
			var function pkg_rhino.Node = first.GetNext()
			_ = function
			var params pkg_rhino.Node = function.GetSecondChild()
			_ = params
			var body_2 pkg_rhino.Node = function.GetLastChild()
			_ = body_2
			this_.Add_1_Node(params)
			this_.Add_1_Node(body_2)
		} else {
			var isInClass bool = node.GetParent().IsClassMembers()
			_ = isInClass
			var initializer pkg_rhino.Node = first.GetNext()
			_ = initializer
			if initializer != 0 {
				if !(!(isInClass)) {
					j2gort.CheckStateStr(false, "initializers should only exist in object literals, not classes")
				}
				this_.cc.Add_1_string(":")
				this_.priv_CodeGenerator_addExpr_1328(initializer, int32(1), StaticGet_CodeGenerator_Context_OTHER())
			} else {
				if !(node.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_VARIABLE())) {
					j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
				}
			}
		}
	case int32(103) /* Token.OBJECT_PATTERN */ :
		this_.AddObjectPattern(node)
		this_.priv_CodeGenerator_maybeAddTypeDecl_1323(node)
	case int32(75) /* Token.SWITCH */ :
		this_.Add_1_string("switch(")
		this_.Add_1_Node(first)
		this_.Add_1_string(")")
		this_.self_.Add_2(last, context)
	case int32(91) /* Token.SWITCH_BODY */ :
		this_.cc.BeginBlock()
		this_.AddAllSiblings(first)
		this_.cc.EndBlock_1((context == StaticGet_CodeGenerator_Context_STATEMENT()))
	case int32(76) /* Token.CASE */ :
		if !(childCount == int32(2)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.Add_1_string("case ")
		this_.Add_1_Node(first)
		this_.AddCaseBody(last)
	case int32(77) /* Token.DEFAULT_CASE */ :
		if !(childCount == int32(1)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		this_.Add_1_string("default")
		this_.AddCaseBody(first)
	case int32(92) /* Token.LABEL */ :
		if !(childCount == int32(2)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		if !(first.IsLabelName()) {
			panic(j2gort.AsIface[javalang.IError](j2gort.ConstructShimThrowable("java.lang.Error", "Unexpected token type. Should be LABEL_NAME.")))
		}
		this_.Add_1_Node(first)
		this_.Add_1_string(":")
		if !(last.IsBlock()) {
			this_.cc.MaybeInsertSpace()
		}
		this_.priv_CodeGenerator_addNonEmptyStatement_1327(last, CodeGenerator_GetContextForNonEmptyExpression(context), true)
	case int32(101) /* Token.CAST */ :
		if this_.preserveTypeAnnotations {
			this_.Add_1_string("(")
			this_.Add_1_Node(first)
			this_.Add_1_string(")")
		} else {
			this_.self_.Add_2(first, context)
		}
	case int32(130) /* Token.TAGGED_TEMPLATELIT */ :
		this_.priv_CodeGenerator_addInvocationTarget_1318(node, context)
		this_.Add_1_Node(first.GetNext())
	case int32(131) /* Token.TEMPLATELIT */ :
		this_.cc.BeginTemplateLit()
		{
			var c_6 pkg_rhino.Node = first
			_ = c_6
			for ; c_6 != 0; c_6 = c_6.GetNext() {
				if c_6.IsTemplateLitString() {
					this_.Add_1_string(this_.priv_CodeGenerator_escapeUnrecognizedCharacters_1352(c_6.GetRawString()))
				} else {
					this_.cc.BeginTemplateLitSub()
					this_.self_.Add_2(c_6.GetFirstChild(), StaticGet_CodeGenerator_Context_START_OF_EXPR())
					this_.cc.EndTemplateLitSub()
				}
			}
		}
		this_.cc.EndTemplateLit()
	case int32(137) /* Token.STRING_TYPE */ :
		this_.Add_1_string("string")
	case int32(138) /* Token.BOOLEAN_TYPE */ :
		this_.Add_1_string("boolean")
	case int32(139) /* Token.NUMBER_TYPE */ :
		this_.Add_1_string("number")
	case int32(143) /* Token.ANY_TYPE */ :
		this_.Add_1_string("any")
	case int32(145) /* Token.VOID_TYPE */ :
		this_.Add_1_string("void")
	case int32(147) /* Token.NAMED_TYPE */ :
		this_.Add_1_Node(first)
	case int32(151) /* Token.ARRAY_TYPE */ :
		this_.priv_CodeGenerator_addExpr_1328(first, NodeUtil_Precedence(pkg_rhino.StaticGet_Token_ARRAY_TYPE()), context)
		this_.Add_1_string("[]")
	case int32(140) /* Token.FUNCTION_TYPE */ :
		var returnType pkg_rhino.Node = first
		_ = returnType
		this_.Add_1_string("(")
		this_.AddList_1(first.GetNext())
		this_.Add_1_string(")")
		this_.cc.AddOp("=>", true)
		this_.Add_1_Node(returnType)
	case int32(142) /* Token.UNION_TYPE */ :
		this_.AddList_2(first, "|")
	case int32(149) /* Token.RECORD_TYPE */ :
		this_.Add_1_string("{")
		this_.AddList_4(first, false, StaticGet_CodeGenerator_Context_OTHER(), ",")
		this_.Add_1_string("}")
	case int32(141) /* Token.PARAMETERIZED_TYPE */ :
		this_.Add_1_Node(first)
		this_.Add_1_string("<")
		this_.AddList_1(first.GetNext())
		this_.Add_1_string(">")
	case int32(153) /* Token.GENERIC_TYPE_LIST */ :
		this_.Add_1_string("<")
		this_.AddList_4(first, false, StaticGet_CodeGenerator_Context_STATEMENT(), ",")
		this_.Add_1_string(">")
	case int32(152) /* Token.GENERIC_TYPE */ :
		this_.priv_CodeGenerator_addIdentifier_1315(node.GetString())
		if node.HasChildren() {
			this_.Add_1_string("extends")
			this_.cc.MaybeInsertSpace()
			this_.Add_1_Node(node.GetFirstChild())
		}
	case int32(164) /* Token.INTERFACE */ :
		{
			if !(childCount == int32(3)) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
			}
			var name_3 pkg_rhino.Node = first
			_ = name_3
			var superTypes pkg_rhino.Node = first.GetNext()
			_ = superTypes
			var members_2 pkg_rhino.Node = last
			_ = members_2
			this_.Add_1_string("interface")
			this_.Add_1_Node(name_3)
			this_.priv_CodeGenerator_maybeAddGenericTypes_1324(name_3)
			if !(superTypes.IsEmpty()) {
				this_.Add_1_string("extends")
				var superType pkg_rhino.Node = superTypes.GetFirstChild()
				_ = superType
				this_.Add_1_Node(superType)
				for (func() pkg_rhino.Node {
					var asgn_4 pkg_rhino.Node = superType.GetNext()
					superType = asgn_4
					return asgn_4
				}()) != 0 {
					this_.Add_1_string(",")
					this_.cc.MaybeInsertSpace()
					this_.Add_1_Node(superType)
				}
			}
			this_.Add_1_Node(members_2)
		}
	case int32(167) /* Token.ENUM */ :
		if !(childCount == int32(2)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		var name_4 pkg_rhino.Node = first
		_ = name_4
		var members_3 pkg_rhino.Node = last
		_ = members_3
		this_.Add_1_string("enum")
		this_.Add_1_Node(name_4)
		this_.Add_1_Node(members_3)
	case int32(175) /* Token.NAMESPACE */ :
		if !(childCount == int32(2)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(node))
		}
		var name_5 pkg_rhino.Node = first
		_ = name_5
		var elements pkg_rhino.Node = last
		_ = elements
		this_.Add_1_string("namespace")
		this_.Add_1_Node(name_5)
		this_.Add_1_Node(elements)
	case int32(170) /* Token.TYPE_ALIAS */ :
		this_.Add_1_string("type")
		this_.Add_1_string(node.GetString())
		this_.cc.AddOp("=", true)
		this_.Add_1_Node(last)
		this_.cc.EndStatement_2(true, this_.priv_CodeGenerator_hasTrailingCommentOnSameLine_1309(node))
	case int32(171) /* Token.DECLARE */ :
		this_.Add_1_string("declare")
		this_.Add_1_Node(first)
		this_.priv_CodeGenerator_processEnd_1353(node, context)
	case int32(173) /* Token.INDEX_SIGNATURE */ :
		this_.Add_1_string("[")
		this_.Add_1_Node(first)
		this_.Add_1_string("]")
		this_.priv_CodeGenerator_maybeAddTypeDecl_1323(node)
		this_.cc.EndStatement_2(true, this_.priv_CodeGenerator_hasTrailingCommentOnSameLine_1309(node))
	case int32(174) /* Token.CALL_SIGNATURE */ :
		if node.GetBooleanProp(pkg_rhino.StaticGet_Node_CONSTRUCT_SIGNATURE()) {
			this_.Add_1_string("new ")
		}
		this_.priv_CodeGenerator_maybeAddGenericTypes_1324(node)
		this_.Add_1_Node(first)
		this_.priv_CodeGenerator_maybeAddTypeDecl_1323(node)
		this_.cc.EndStatement_2(true, this_.priv_CodeGenerator_hasTrailingCommentOnSameLine_1309(node))
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(("Unknown token " + j2gort.StrVal(type_) + "\n" + j2gort.StrValString(node.ToStringTree()))))))
	}
	this_.PrintTrailingComment(node)
	this_.cc.EndSourceMapping(node)
}

// Java method: CodeGenerator.addIdentifier(java.lang.String identifier)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_addIdentifier_1315(identifier string) {
	_ = identifier
	this_.cc.AddIdentifier(CodeGenerator_IdentifierEscape(identifier))
}

// Java method: CodeGenerator.addGetpropIdentifier(com.google.javascript.rhino.Node getprop)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_addGetpropIdentifier_1316(getprop pkg_rhino.Node) {
	_ = getprop
	this_.cc.StartSourceMapping(getprop)
	this_.priv_CodeGenerator_addIdentifier_1315(getprop.GetString())
	this_.cc.EndSourceMapping(getprop)
}

// Java method: CodeGenerator.precedence(com.google.javascript.rhino.Node n)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_precedence_1317(n pkg_rhino.Node) int32 {
	_ = n
	if n.IsCast() {
		return this_.priv_CodeGenerator_precedence_1317(n.GetFirstChild())
	}
	return NodeUtil_Precedence(n.GetToken())
}

// Java method: CodeGenerator.addInvocationTarget(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.CodeGenerator.Context context)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_addInvocationTarget_1318(node pkg_rhino.Node, context CodeGenerator_Context) {
	_ = node
	_ = context
	var first pkg_rhino.Node = node.GetFirstChild()
	_ = first
	var needsParens bool = NodeUtil_IsOptChainNode(first)
	_ = needsParens
	if CodeGenerator_IsIndirectEval(first) || (node.GetBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL()) && NodeUtil_IsNormalOrOptChainGet(first)) {
		this_.Add_1_string("(0,")
		this_.priv_CodeGenerator_addExpr_1328(first, NodeUtil_Precedence(pkg_rhino.StaticGet_Token_COMMA()), StaticGet_CodeGenerator_Context_OTHER())
		this_.Add_1_string(")")
	} else {
		if needsParens {
			this_.Add_1_string("(")
		}
		this_.priv_CodeGenerator_addExpr_1328(first, NodeUtil_Precedence(node.GetToken()), context)
		if needsParens {
			this_.Add_1_string(")")
		}
	}
}

// Java method: CodeGenerator.arrowFunctionNeedsParens(com.google.javascript.rhino.Node n)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_arrowFunctionNeedsParens_1319(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	var expressionOrEnclosingCast pkg_rhino.Node = n
	_ = expressionOrEnclosingCast
	for (parent != 0) && parent.IsCast() {
		if this_.preserveTypeAnnotations {
			return false
		}
		expressionOrEnclosingCast = parent
		parent = parent.GetParent()
	}
	if parent == 0 {
		return false
	} else if ((((((NodeUtil_IsBinaryOperator(parent) || NodeUtil_IsUnaryOperator(parent)) || NodeUtil_IsUpdateOperator(parent)) || parent.IsTaggedTemplateLit()) || parent.IsGetProp()) || parent.IsOptChainGetProp()) || parent.IsAwait()) || parent.IsYield() {
		return true
	} else if ((((parent.IsGetElem() || parent.IsCall()) || parent.IsHook()) || parent.IsOptChainGetElem()) || parent.IsOptChainCall()) || parent.IsNew() {
		return expressionOrEnclosingCast.IsFirstChildOf(parent)
	} else {
		return false
	}
}

// Java method: CodeGenerator.addArrowFunction(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node first, com.google.javascript.rhino.Node last, com.google.javascript.jscomp.CodeGenerator.Context context)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_addArrowFunction_1320(n pkg_rhino.Node, first pkg_rhino.Node, last pkg_rhino.Node, context CodeGenerator_Context) {
	_ = n
	_ = first
	_ = last
	_ = context
	if !(j2gort.CheckNotNullString(first.GetString()) == "") {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(first))
	}
	var funcNeedsParens bool = this_.priv_CodeGenerator_arrowFunctionNeedsParens_1319(n)
	_ = funcNeedsParens
	if funcNeedsParens {
		this_.Add_1_string("(")
	}
	this_.priv_CodeGenerator_maybeAddGenericTypes_1324(first)
	if n.IsAsyncFunction() {
		this_.Add_1_string("async")
	}
	this_.Add_1_Node(first.GetNext())
	this_.priv_CodeGenerator_maybeAddTypeDecl_1323(n)
	this_.cc.AddOp("=>", true)
	if last.IsBlock() {
		this_.Add_1_Node(last)
	} else {
		this_.priv_CodeGenerator_addExpr_1328(last, (NodeUtil_Precedence(pkg_rhino.StaticGet_Token_COMMA()) + int32(1)), CodeGenerator_GetContextForArrowFunctionBody(context))
	}
	this_.cc.EndFunction((context == StaticGet_CodeGenerator_Context_STATEMENT()))
	if funcNeedsParens {
		this_.Add_1_string(")")
	}
}

// Java method: CodeGenerator.addFunction(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node first, com.google.javascript.rhino.Node last, com.google.javascript.jscomp.CodeGenerator.Context context)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_addFunction_1321(n pkg_rhino.Node, first pkg_rhino.Node, last pkg_rhino.Node, context CodeGenerator_Context) {
	_ = n
	_ = first
	_ = last
	_ = context
	var funcNeedsParens bool = (context == StaticGet_CodeGenerator_Context_START_OF_EXPR())
	_ = funcNeedsParens
	if funcNeedsParens {
		this_.Add_1_string("(")
	}
	this_.Add_1_string(func() string {
		if n.IsAsyncFunction() {
			return "async function"
		}
		return "function"
	}())
	if n.IsGeneratorFunction() {
		this_.Add_1_string("*")
		if !(j2gort.CheckNotNullString(first.GetString()) == "") {
			this_.cc.MaybeInsertSpace()
		}
	}
	this_.Add_1_Node(first)
	this_.priv_CodeGenerator_maybeAddGenericTypes_1324(first)
	this_.Add_1_Node(first.GetNext())
	this_.priv_CodeGenerator_maybeAddTypeDecl_1323(n)
	this_.Add_1_Node(last)
	this_.cc.EndFunction((context == StaticGet_CodeGenerator_Context_STATEMENT()))
	if funcNeedsParens {
		this_.Add_1_string(")")
	}
}

// Java method: CodeGenerator.maybeAddAccessibilityModifier(com.google.javascript.rhino.Node n)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_maybeAddAccessibilityModifier_1322(n pkg_rhino.Node) {
	_ = n
	var access pkg_rhino.JSDocInfo_Visibility = n.GetAccessModifier()
	_ = access
	if access != 0 {
		this_.Add_1_string((j2gort.StrValString(j2gort.Ext_com_google_common_base_Ascii_ToLowerCase_1_string(j2gort.CheckNotNullEnum(access).String())) + " "))
	}
}

// Java method: CodeGenerator.maybeAddTypeDecl(com.google.javascript.rhino.Node n)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_maybeAddTypeDecl_1323(n pkg_rhino.Node) {
	_ = n
	if n.GetDeclaredTypeExpression() != 0 {
		this_.Add_1_string(":")
		this_.cc.MaybeInsertSpace()
		this_.Add_1_Node(n.GetDeclaredTypeExpression())
	}
}

// Java method: CodeGenerator.maybeAddGenericTypes(com.google.javascript.rhino.Node n)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_maybeAddGenericTypes_1324(n pkg_rhino.Node) {
	_ = n
	var generics pkg_rhino.Node = n.GetGenericType()
	_ = generics
	if generics != 0 {
		this_.Add_1_Node(generics)
	}
}

// Java method: CodeGenerator.maybeAddOptional(com.google.javascript.rhino.Node n)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_maybeAddOptional_1325(n pkg_rhino.Node) {
	_ = n
	if n.GetBooleanProp(pkg_rhino.StaticGet_Node_OPT_ES6_TYPED()) {
		this_.Add_1_string("?")
	}
}

// Java method: CodeGenerator.unrollBinaryOperator(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Token op, java.lang.String opStr, com.google.javascript.jscomp.CodeGenerator.Context context, com.google.javascript.jscomp.CodeGenerator.Context rhsContext, int leftPrecedence, int rightPrecedence)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_unrollBinaryOperator_1326(n pkg_rhino.Node, op pkg_rhino.Token, opStr string, context CodeGenerator_Context, rhsContext CodeGenerator_Context, leftPrecedence int32, rightPrecedence int32) {
	_ = n
	_ = op
	_ = opStr
	_ = context
	_ = rhsContext
	_ = leftPrecedence
	_ = rightPrecedence
	var firstNonOperator pkg_rhino.Node = n.GetFirstChild()
	_ = firstNonOperator
	for firstNonOperator.GetToken() == op {
		firstNonOperator = firstNonOperator.GetFirstChild()
	}
	this_.priv_CodeGenerator_addExpr_1328(firstNonOperator, leftPrecedence, context)
	var current pkg_rhino.Node = firstNonOperator
	_ = current
	for doFirst_5 := true; doFirst_5 || (current != n); doFirst_5 = false {
		current = current.GetParent()
		this_.cc.AddOp(opStr, true)
		this_.priv_CodeGenerator_addExpr_1328(current.GetSecondChild(), rightPrecedence, rhsContext)
	}
}

// Java static method: CodeGenerator.isSimpleNumber(java.lang.String s)
func CodeGenerator_IsSimpleNumber(s string) bool {
	_ = s
	var len_ int32 = j2gort.JStrLen(s)
	_ = len_
	if len_ == int32(0) {
		return false
	}
	{
		var index int32 = int32(0)
		_ = index
		for ; index < len_; index = index + 1 {
			var c uint16 = j2gort.JStrCharAt(s, index)
			_ = c
			if (int32(c) < int32(uint16(48))) || (int32(c) > int32(uint16(57))) {
				return false
			}
		}
	}
	return ((len_ == int32(1)) || (int32(j2gort.JStrCharAt(s, int32(0))) != int32(uint16(48))))
}

// Java static method: CodeGenerator.getSimpleNumber(java.lang.String s)
func CodeGenerator_GetSimpleNumber(s string) float64 {
	_ = s
	if CodeGenerator_IsSimpleNumber(s) {
		var didRet_6 bool
		var retVal_7 float64
		func() {
			defer func() {
				if rec_8 := recover(); rec_8 != nil {
					jex_9 := javalang.RecoverException(rec_8)
					_ = jex_9
					if j2gort.InstanceOf[javalang.INumberFormatException](jex_9) {
						e := j2gort.AsIface[javalang.INumberFormatException](jex_9)
						_ = e
					} else {
						panic(rec_8)
					}
				}
			}()
			var l int64 = j2gort.ParseLong(s, 10)
			_ = l
			if l <= int64(9007199254740991) {
				retVal_7 = float64(l)
				didRet_6 = true
				return
			}
		}()
		if didRet_6 {
			return retVal_7
		}
	}
	return math.NaN()
}

// Java static method: CodeGenerator.isIndirectEval(com.google.javascript.rhino.Node n)
func CodeGenerator_IsIndirectEval(n pkg_rhino.Node) bool {
	_ = n
	return ((n.IsName() && func() bool { sr_10 := "eval"; se_11 := n.GetString(); return j2gort.CheckNotNullString(sr_10) == se_11 }()) && !(n.GetBooleanProp(pkg_rhino.StaticGet_Node_DIRECT_EVAL())))
}

// Java method: CodeGenerator.addNonEmptyStatement(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.CodeGenerator.Context context, boolean allowNonBlockChild)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_addNonEmptyStatement_1327(n pkg_rhino.Node, context CodeGenerator_Context, allowNonBlockChild bool) {
	_ = n
	_ = context
	_ = allowNonBlockChild
	var nodeToProcess pkg_rhino.Node = n
	_ = nodeToProcess
	if !(allowNonBlockChild) && !(n.IsBlock()) {
		panic(j2gort.AsIface[javalang.IError](j2gort.ConstructShimThrowable("java.lang.Error", "Missing BLOCK child.")))
	}
	if n.IsBlock() {
		var count int32 = CodeGenerator_GetNonEmptyChildCount(n, int32(2))
		_ = count
		if count == int32(0) {
			if this_.cc.ShouldPreserveExtras(n) {
				this_.cc.BeginBlock()
				this_.PrintTrailingComment(n)
				this_.cc.EndBlock_1(this_.cc.BreakAfterBlockFor(n, (context == StaticGet_CodeGenerator_Context_STATEMENT())))
			} else {
				this_.PrintTrailingComment(n)
				this_.cc.EndStatement_2(true, false)
			}
			return
		}
		if count == int32(1) {
			var firstAndOnlyChild pkg_rhino.Node = CodeGenerator_GetFirstNonEmptyChild(n)
			_ = firstAndOnlyChild
			var alwaysWrapInBlock bool = this_.cc.ShouldPreserveExtras(n)
			_ = alwaysWrapInBlock
			if alwaysWrapInBlock || CodeGenerator_IsBlockDeclOrDo(firstAndOnlyChild) {
				this_.cc.BeginBlock()
				this_.self_.Add_2(firstAndOnlyChild, StaticGet_CodeGenerator_Context_STATEMENT())
				this_.PrintTrailingComment(n)
				this_.cc.MaybeLineBreak()
				this_.cc.EndBlock_1(this_.cc.BreakAfterBlockFor(n, (context == StaticGet_CodeGenerator_Context_STATEMENT())))
				return
			} else {
				nodeToProcess = firstAndOnlyChild
			}
		}
	}
	if nodeToProcess.IsEmpty() {
		this_.PrintTrailingComment(n)
		this_.cc.EndStatement_2(true, false)
	} else {
		this_.self_.Add_2(nodeToProcess, context)
		this_.PrintTrailingComment(n)
	}
}

// Java static method: CodeGenerator.isBlockDeclOrDo(com.google.javascript.rhino.Node n)
func CodeGenerator_IsBlockDeclOrDo(n pkg_rhino.Node) bool {
	_ = n
	if n.IsLabel() {
		var labeledStatement pkg_rhino.Node = n.GetLastChild()
		_ = labeledStatement
		if !(labeledStatement.IsBlock()) {
			return CodeGenerator_IsBlockDeclOrDo(labeledStatement)
		} else {
			if CodeGenerator_GetNonEmptyChildCount(n, int32(2)) == int32(1) {
				return CodeGenerator_IsBlockDeclOrDo(CodeGenerator_GetFirstNonEmptyChild(n))
			} else {
				return false
			}
		}
	} else {
		switch n.TokenSwitchOrd_() {
		case int32(111) /* Token.LET */, int32(97) /* Token.CONST */, int32(73) /* Token.FUNCTION */, int32(105) /* Token.CLASS */, int32(79) /* Token.DO */ :
			return true
		default:
			return false
		}
		panic("unreachable switch expression")
	}
}

// Java method: CodeGenerator.addExpr(com.google.javascript.rhino.Node n, int minPrecedence, com.google.javascript.jscomp.CodeGenerator.Context context)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_addExpr_1328(n pkg_rhino.Node, minPrecedence int32, context CodeGenerator_Context) {
	_ = n
	_ = minPrecedence
	_ = context
	if this_.priv_CodeGenerator_opRequiresParentheses_1329(n, minPrecedence, context) {
		this_.Add_1_string("(")
		this_.self_.Add_2(n, StaticGet_CodeGenerator_Context_OTHER())
		this_.Add_1_string(")")
	} else {
		this_.self_.Add_2(n, context)
	}
}

// Java method: CodeGenerator.opRequiresParentheses(com.google.javascript.rhino.Node n, int minPrecedence, com.google.javascript.jscomp.CodeGenerator.Context context)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_opRequiresParentheses_1329(n pkg_rhino.Node, minPrecedence int32, context CodeGenerator_Context) bool {
	_ = n
	_ = minPrecedence
	_ = context
	if context.InForInInitClause() && n.IsIn() {
		return true
	} else if NodeUtil_IsUnaryOperator(n) && this_.priv_CodeGenerator_isFirstOperandOfExponentiationExpression_1333(n) {
		return true
	} else if this_.priv_CodeGenerator_isLogicalAndOrLogicalOrChildOfNullishCoalesce_1330(n) || this_.priv_CodeGenerator_isNullishCoalesceChildOfLogicalAndOrLogicalOr_1331(n) {
		return true
	} else if n.IsAssign() && n.GetParent().IsClass() {
		return true
	} else {
		return (this_.priv_CodeGenerator_precedence_1317(n) < minPrecedence)
	}
}

// Java method: CodeGenerator.isLogicalAndOrLogicalOrChildOfNullishCoalesce(com.google.javascript.rhino.Node n)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_isLogicalAndOrLogicalOrChildOfNullishCoalesce_1330(n pkg_rhino.Node) bool {
	_ = n
	var expr pkg_rhino.Node = this_.priv_CodeGenerator_skipTransparentCasts_1332(n)
	_ = expr
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	var logicalAndOrLogicalOr bool = (expr.IsAnd() || expr.IsOr())
	_ = logicalAndOrLogicalOr
	var childOfNullishCoalesce bool = ((parent != 0) && parent.IsNullishCoalesce())
	_ = childOfNullishCoalesce
	return (logicalAndOrLogicalOr && childOfNullishCoalesce)
}

// Java method: CodeGenerator.isNullishCoalesceChildOfLogicalAndOrLogicalOr(com.google.javascript.rhino.Node n)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_isNullishCoalesceChildOfLogicalAndOrLogicalOr_1331(n pkg_rhino.Node) bool {
	_ = n
	var expr pkg_rhino.Node = this_.priv_CodeGenerator_skipTransparentCasts_1332(n)
	_ = expr
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	var childOfLogicalAndOrLogicalOr bool = ((parent != 0) && (parent.IsAnd() || parent.IsOr()))
	_ = childOfLogicalAndOrLogicalOr
	return (expr.IsNullishCoalesce() && childOfLogicalAndOrLogicalOr)
}

// Java method: CodeGenerator.skipTransparentCasts(com.google.javascript.rhino.Node n)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_skipTransparentCasts_1332(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	for !(this_.preserveTypeAnnotations) && n.IsCast() {
		n = n.GetFirstChild()
	}
	return n
}

// Java method: CodeGenerator.isFirstOperandOfExponentiationExpression(com.google.javascript.rhino.Node n)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_isFirstOperandOfExponentiationExpression_1333(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	return (((parent != 0) && parent.IsExponent()) && (parent.GetFirstChild() == n))
}

// Java method: CodeGenerator.addList(com.google.javascript.rhino.Node firstInList)
func (this_ *CodeGeneratorBase) AddList_1(firstInList pkg_rhino.Node) {
	_ = firstInList
	this_.AddList_4(firstInList, true, StaticGet_CodeGenerator_Context_OTHER(), ",")
}

// Java method: CodeGenerator.addList(com.google.javascript.rhino.Node firstInList, java.lang.String separator)
func (this_ *CodeGeneratorBase) AddList_2(firstInList pkg_rhino.Node, separator string) {
	_ = firstInList
	_ = separator
	this_.AddList_4(firstInList, true, StaticGet_CodeGenerator_Context_OTHER(), separator)
}

// Java method: CodeGenerator.addList(com.google.javascript.rhino.Node firstInList, boolean isArrayOrFunctionArgument, com.google.javascript.jscomp.CodeGenerator.Context lhsContext, java.lang.String separator)
func (this_ *CodeGeneratorBase) AddList_4(firstInList pkg_rhino.Node, isArrayOrFunctionArgument bool, lhsContext CodeGenerator_Context, separator string) {
	_ = firstInList
	_ = isArrayOrFunctionArgument
	_ = lhsContext
	_ = separator
	if firstInList == 0 {
		return
	}
	{
		var n pkg_rhino.Node = firstInList
		_ = n
		for ; n != 0; n = n.GetNext() {
			var isFirst bool = (n == firstInList)
			_ = isFirst
			var minPrecedence int32
			_ = minPrecedence
			if isArrayOrFunctionArgument {
				minPrecedence = int32(1)
			} else {
				minPrecedence = int32(0)
			}
			if isFirst {
				this_.priv_CodeGenerator_addExpr_1328(n, minPrecedence, lhsContext)
			} else {
				this_.cc.AddOp(separator, true)
				this_.priv_CodeGenerator_addExpr_1328(n, minPrecedence, CodeGenerator_GetContextForNoInOperator(lhsContext))
			}
		}
	}
	if isArrayOrFunctionArgument && j2gort.CheckNotNullNode(firstInList.GetParent()).HasTrailingComma() {
		this_.cc.OptionalListSeparator()
	}
}

// Java method: CodeGenerator.addStringKey(com.google.javascript.rhino.Node n)
func (this_ *CodeGeneratorBase) AddStringKey(n pkg_rhino.Node) {
	_ = n
	var key string = n.GetString()
	_ = key
	var mustBeQuoted bool = (((n.IsQuotedStringKey() || (this_.quoteKeywordProperties && pkg_rhino.TokenStream_IsKeyword(key))) || !(pkg_rhino.TokenStream_IsJSIdentifier(key))) || !(NodeUtil_IsLatin(key)))
	_ = mustBeQuoted
	if !(mustBeQuoted) {
		if n.IsShorthandProperty() {
			var child pkg_rhino.Node = n.GetFirstChild()
			_ = child
			if child.MatchesQualifiedName_1_string(key) || (child.IsDefaultValue() && child.GetFirstChild().MatchesQualifiedName_1_string(key)) {
				this_.Add_1_Node(child)
				return
			}
		}
		this_.Add_1_string(key)
	} else {
		var d float64 = CodeGenerator_GetSimpleNumber(key)
		_ = d
		if !(math.IsNaN(d)) {
			this_.cc.AddNumber(d, n)
		} else {
			this_.priv_CodeGenerator_addJsString_1344(n)
		}
	}
	if n.HasChildren() {
		this_.Add_1_string(":")
		this_.priv_CodeGenerator_addExpr_1328(n.GetFirstChild(), int32(1), StaticGet_CodeGenerator_Context_OTHER())
	}
}

// Java method: CodeGenerator.addObjectPattern(com.google.javascript.rhino.Node n)
func (this_ *CodeGeneratorBase) AddObjectPattern(n pkg_rhino.Node) {
	_ = n
	this_.Add_1_string("{")
	{
		var child pkg_rhino.Node = n.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			if child != n.GetFirstChild() {
				this_.cc.ListSeparator()
			}
			this_.Add_1_Node(child)
		}
	}
	this_.Add_1_string("}")
}

// Java method: CodeGenerator.addArrayList(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node firstInList)
func (this_ *CodeGeneratorBase) AddArrayList(firstInList pkg_rhino.Node) {
	_ = firstInList
	if firstInList == 0 {
		return
	}
	var lastWasEmpty bool = false
	_ = lastWasEmpty
	{
		var n pkg_rhino.Node = firstInList
		_ = n
		for ; n != 0; n = n.GetNext() {
			if n != firstInList {
				this_.cc.ListSeparator()
			}
			this_.priv_CodeGenerator_addExpr_1328(n, int32(1), StaticGet_CodeGenerator_Context_OTHER())
			lastWasEmpty = n.IsEmpty()
		}
	}
	if lastWasEmpty {
		this_.cc.ListSeparator()
	} else if firstInList.GetParent().HasTrailingComma() {
		this_.cc.OptionalListSeparator()
	}
}

// Java method: CodeGenerator.addCaseBody(com.google.javascript.rhino.Node caseBody)
func (this_ *CodeGeneratorBase) AddCaseBody(caseBody pkg_rhino.Node) {
	_ = caseBody
	if !(caseBody.IsBlock()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(caseBody))
	}
	this_.cc.BeginCaseBody()
	this_.AddAllSiblings(caseBody.GetFirstChild())
	this_.cc.EndCaseBody()
}

// Java method: CodeGenerator.addAllSiblings(com.google.javascript.rhino.Node n)
func (this_ *CodeGeneratorBase) AddAllSiblings(n pkg_rhino.Node) {
	_ = n
	{
		var c pkg_rhino.Node = n
		_ = c
		for ; c != 0; c = c.GetNext() {
			this_.Add_1_Node(c)
		}
	}
}

// Java method: CodeGenerator.addNonJsDocNonTrailing(com.google.javascript.rhino.Node node, com.google.javascript.rhino.NonJSDocComment nonJsDocComment)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_addNonJsDocNonTrailing_1342(node pkg_rhino.Node, nonJsDocComment *pkg_rhino.NonJSDocComment) {
	_ = node
	_ = nonJsDocComment
	var content string = nonJsDocComment.GetCommentString()
	_ = content
	var commentEndPosition pkg_parser.SourcePosition = nonJsDocComment.GetEndPosition()
	_ = commentEndPosition
	var nodeLineNumber int32 = (node.GetLineno() - int32(1))
	_ = nodeLineNumber
	if nonJsDocComment.IsEndingAsLineComment() {
		if !(this_.gentsMode || (commentEndPosition.FieldGet_line() < nodeLineNumber)) {
			j2gort.CheckStateStr(false, "Non trailing line comments can not be on the same line as the node.")
		}
		this_.Add_1_string((j2gort.StrValString(content) + "\n"))
	} else {
		if nodeLineNumber == commentEndPosition.FieldGet_line() {
			this_.Add_1_string((j2gort.StrValString(content) + " "))
		} else {
			this_.Add_1_string((j2gort.StrValString(content) + "\n"))
		}
	}
}

// Java method: CodeGenerator.addNonJsDoctrailing(com.google.javascript.rhino.NonJSDocComment nonJsDocComment, boolean sameLine)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_addNonJsDoctrailing_1343(nonJsDocComment *pkg_rhino.NonJSDocComment, sameLine bool) {
	_ = nonJsDocComment
	_ = sameLine
	var content string = nonJsDocComment.GetCommentString()
	_ = content
	if nonJsDocComment.IsEndingAsLineComment() {
		this_.Add_1_string((" " + j2gort.StrValString(content)))
		if sameLine {
			this_.cc.StartNewLine()
		}
	} else {
		if nonJsDocComment.IsInline() {
			this_.Add_1_string((" " + j2gort.StrValString(content)))
		} else {
			this_.Add_1_string((" " + j2gort.StrValString(content)))
			if sameLine {
				this_.cc.StartNewLine()
			}
		}
	}
}

// Java method: CodeGenerator.addJsString(com.google.javascript.rhino.Node n)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_addJsString_1344(n pkg_rhino.Node) {
	_ = n
	this_.Add_1_string(this_.priv_CodeGenerator_jsString_1345(n.GetString()))
}

// Java method: CodeGenerator.jsString(java.lang.String s)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_jsString_1345(s string) string {
	_ = s
	var singleq int32 = int32(0)
	_ = singleq
	var doubleq int32 = int32(0)
	_ = doubleq
	{
		var i int32 = int32(0)
		_ = i
		for ; i < j2gort.JStrLen(s); i = i + 1 {
			switch j2gort.JStrCharAt(s, i) {
			case uint16(34):
				doubleq = doubleq + 1
			case uint16(39):
				singleq = singleq + 1
			default:
			}
		}
	}
	var doublequote string = j2gort.NullString
	_ = doublequote
	var singlequote string = j2gort.NullString
	_ = singlequote
	var quote uint16 = uint16(0)
	_ = quote
	if func() bool {
		if this_.preferSingleQuotes {
			return (singleq <= doubleq)
		}
		return (singleq < doubleq)
	}() {
		quote = uint16(39)
		doublequote = "\""
		singlequote = "\\'"
	} else {
		quote = uint16(34)
		doublequote = "\\\""
		singlequote = "'"
	}
	return (j2gort.StrValChar(quote) + j2gort.StrValString(this_.priv_CodeGenerator_strEscape_1347(s, doublequote, singlequote, "`", "\\\\", "$", false)) + j2gort.StrValChar(quote))
}

// Java method: CodeGenerator.regexpEscape(java.lang.String s)
func (this_ *CodeGeneratorBase) RegexpEscape(s string) string {
	_ = s
	return ("/" + j2gort.StrValString(this_.priv_CodeGenerator_strEscape_1347(s, "\"", "'", "`", "\\", "$", true)) + "/")
}

// Java method: CodeGenerator.strEscape(java.lang.String s, java.lang.String doublequoteEscape, java.lang.String singlequoteEscape, java.lang.String backtickEscape, java.lang.String backslashEscape, java.lang.String dollarEscape, boolean isRegexp)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_strEscape_1347(s string, doublequoteEscape string, singlequoteEscape string, backtickEscape string, backslashEscape string, dollarEscape string, isRegexp bool) string {
	_ = s
	_ = doublequoteEscape
	_ = singlequoteEscape
	_ = backtickEscape
	_ = backslashEscape
	_ = dollarEscape
	_ = isRegexp
	var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder_1_int32((j2gort.JStrLen(s) + int32(2)))
	_ = sb
	{
		var i int32 = int32(0)
		_ = i
		for ; i < j2gort.JStrLen(s); i = i + 1 {
			var c uint16 = j2gort.JStrCharAt(s, i)
			_ = c
			switch c {
			case uint16(0):
				sb.Append_1_string("\\x00")
			case uint16(11):
				if !(isRegexp) {
					sb.Append_1_string("\\v")
				} else {
					sb.Append_1_string("\\x0B")
				}
			case uint16(8):
				sb.Append_1_string("\\b")
			case uint16(12):
				sb.Append_1_string("\\f")
			case uint16(10):
				sb.Append_1_string("\\n")
			case uint16(13):
				sb.Append_1_string("\\r")
			case uint16(9):
				sb.Append_1_string("\\t")
			case uint16(92):
				sb.Append_1_string(backslashEscape)
			case uint16(34):
				sb.Append_1_string(doublequoteEscape)
			case uint16(39):
				sb.Append_1_string(singlequoteEscape)
			case uint16(36):
				sb.Append_1_string(dollarEscape)
			case uint16(96):
				sb.Append_1_string(backtickEscape)
			case uint16(61):
				this_.priv_CodeGenerator_appendEscapedEq_1348(sb, isRegexp)
			case uint16(38):
				this_.priv_CodeGenerator_appendEscapedAmp_1349(sb, isRegexp)
			case uint16(62):
				this_.priv_CodeGenerator_appendEscapedGt_1350(sb, s, i, isRegexp)
			case uint16(60):
				this_.priv_CodeGenerator_appendEscapedLt_1351(sb, s, i, isRegexp)
			default:
				if isRegexp || !(this_.outputFeatureSet.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_UNESCAPED_UNICODE_LINE_OR_PARAGRAPH_SEP())) {
					if int32(c) == int32(uint16(8232)) {
						sb.Append_1_string("\\u2028")
						break
					}
					if int32(c) == int32(uint16(8233)) {
						sb.Append_1_string("\\u2029")
						break
					}
				}
				if ((this_.outputCharsetEncoder != nil) && this_.outputCharsetEncoder.CanEncode(c)) || ((int32(c) > int32(31)) && (int32(c) < int32(127))) {
					sb.Append_1_uint16(c)
				} else {
					pkg_sourcemap.Util_AppendHexJavaScriptRepresentation_2_Ext_java_lang_StringBuilder_uint16(sb, c)
				}
			}
		}
	}
	return sb.String()
}

// Java method: CodeGenerator.appendEscapedEq(java.lang.StringBuilder sb, boolean isRegexp)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_appendEscapedEq_1348(sb *j2gort.Ext_java_lang_StringBuilder, isRegexp bool) {
	_ = sb
	_ = isRegexp
	if this_.trustedStrings || isRegexp {
		sb.Append_1_uint16(uint16(61))
	} else {
		sb.Append_1_string("\\x3d")
	}
}

// Java method: CodeGenerator.appendEscapedAmp(java.lang.StringBuilder sb, boolean isRegexp)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_appendEscapedAmp_1349(sb *j2gort.Ext_java_lang_StringBuilder, isRegexp bool) {
	_ = sb
	_ = isRegexp
	if this_.trustedStrings || isRegexp {
		sb.Append_1_uint16(uint16(38))
	} else {
		sb.Append_1_string("\\x26")
	}
}

// Java method: CodeGenerator.appendEscapedGt(java.lang.StringBuilder sb, java.lang.String s, int index, boolean isRegexp)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_appendEscapedGt_1350(sb *j2gort.Ext_java_lang_StringBuilder, s string, index int32, isRegexp bool) {
	_ = sb
	_ = s
	_ = index
	_ = isRegexp
	if !(this_.trustedStrings) && !(isRegexp) {
		sb.Append_1_string("\\x3e")
		return
	}
	if (index >= int32(2)) && (((int32(j2gort.JStrCharAt(s, (index-int32(1)))) == int32(uint16(45))) && (int32(j2gort.JStrCharAt(s, (index-int32(2)))) == int32(uint16(45)))) || ((int32(j2gort.JStrCharAt(s, (index-int32(1)))) == int32(uint16(93))) && (int32(j2gort.JStrCharAt(s, (index-int32(2)))) == int32(uint16(93))))) {
		sb.Append_1_string("\\x3e")
	} else {
		sb.Append_1_uint16(uint16(62))
	}
}

// Java method: CodeGenerator.appendEscapedLt(java.lang.StringBuilder sb, java.lang.String s, int index, boolean isRegexp)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_appendEscapedLt_1351(sb *j2gort.Ext_java_lang_StringBuilder, s string, index int32, isRegexp bool) {
	_ = sb
	_ = s
	_ = index
	_ = isRegexp
	if !(this_.trustedStrings) && !(isRegexp) {
		sb.Append_1_string("\\x3c")
		return
	}
	var endScript string = "/script"
	_ = endScript
	var endStyle string = "/style"
	_ = endStyle
	var startComment string = "!--"
	_ = startComment
	if (j2gort.JStrRegionMatches(s, true, (index+int32(1)), endScript, int32(0), j2gort.JStrLen(endScript)) || j2gort.JStrRegionMatches(s, true, (index+int32(1)), endStyle, int32(0), j2gort.JStrLen(endStyle))) || j2gort.JStrRegionMatches(s, false, (index+int32(1)), startComment, int32(0), j2gort.JStrLen(startComment)) {
		sb.Append_1_string("\\x3c")
	} else {
		sb.Append_1_uint16(uint16(60))
	}
}

// Java method: CodeGenerator.escapeUnrecognizedCharacters(java.lang.String s)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_escapeUnrecognizedCharacters_1352(s string) string {
	_ = s
	var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = sb
	var escapeStart int32 = int32(-1)
	_ = escapeStart
	{
		var i int32 = int32(0)
		_ = i
		for ; i < j2gort.JStrLen(s); i = i + 1 {
			var c uint16 = j2gort.JStrCharAt(s, i)
			_ = c
			var start int32 = sb.Length()
			_ = start
			switch c {
			case uint16(8), uint16(12), uint16(10), uint16(13), uint16(9), uint16(92), uint16(34), uint16(39), uint16(36), uint16(96), uint16(8232), uint16(8233):
				sb.Append_1_uint16(c)
			case uint16(61):
				this_.priv_CodeGenerator_appendEscapedEq_1348(sb, false)
			case uint16(38):
				this_.priv_CodeGenerator_appendEscapedAmp_1349(sb, false)
			case uint16(62):
				this_.priv_CodeGenerator_appendEscapedGt_1350(sb, s, i, false)
			case uint16(60):
				this_.priv_CodeGenerator_appendEscapedLt_1351(sb, s, i, false)
			default:
				if ((this_.outputCharsetEncoder != nil) && this_.outputCharsetEncoder.CanEncode(c)) || ((int32(c) > int32(31)) && (int32(c) < int32(127))) {
					sb.Append_1_uint16(c)
				} else {
					pkg_sourcemap.Util_AppendHexJavaScriptRepresentation_2_Ext_java_lang_StringBuilder_uint16(sb, c)
				}
			}
			if escapeStart >= int32(0) {
				if (sb.Length() - start) > int32(1) {
					sb.DeleteCharAt(escapeStart)
				}
				escapeStart = int32(-1)
			} else if int32(c) == int32(uint16(92)) {
				escapeStart = start
			}
		}
	}
	return sb.String()
}

// Java static method: CodeGenerator.identifierEscape(java.lang.String s)
func CodeGenerator_IdentifierEscape(s string) string {
	_ = s
	if NodeUtil_IsLatin(s) {
		return s
	}
	var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = sb
	{
		var i int32 = int32(0)
		_ = i
		for ; i < j2gort.JStrLen(s); i = i + 1 {
			var c uint16 = j2gort.JStrCharAt(s, i)
			_ = c
			if (int32(c) > int32(31)) && (int32(c) < int32(127)) {
				sb.Append_1_uint16(c)
			} else {
				pkg_sourcemap.Util_AppendHexJavaScriptRepresentation_2_Ext_java_lang_StringBuilder_uint16(sb, c)
			}
		}
	}
	return sb.String()
}

// Java static method: CodeGenerator.getNonEmptyChildCount(com.google.javascript.rhino.Node n, int maxCount)
func CodeGenerator_GetNonEmptyChildCount(n pkg_rhino.Node, maxCount int32) int32 {
	_ = n
	_ = maxCount
	var i int32 = int32(0)
	_ = i
	var c pkg_rhino.Node = n.GetFirstChild()
	_ = c
	for ; (c != 0) && (i < maxCount); c = c.GetNext() {
		if c.IsBlock() {
			i = (i + CodeGenerator_GetNonEmptyChildCount(c, (maxCount-i)))
		} else if !(c.IsEmpty()) {
			i = i + 1
		}
	}
	return i
}

// Java static method: CodeGenerator.getFirstNonEmptyChild(com.google.javascript.rhino.Node n)
func CodeGenerator_GetFirstNonEmptyChild(n pkg_rhino.Node) pkg_rhino.Node {
	_ = n
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = c.GetNext() {
			if c.IsBlock() {
				var result pkg_rhino.Node = CodeGenerator_GetFirstNonEmptyChild(c)
				_ = result
				if result != 0 {
					return result
				}
			} else if !(c.IsEmpty()) {
				return c
			}
		}
	}
	return pkg_rhino.Node(0)
}

// Java static method: CodeGenerator.getContextForNonEmptyExpression(com.google.javascript.jscomp.CodeGenerator.Context currentContext)
func CodeGenerator_GetContextForNonEmptyExpression(currentContext CodeGenerator_Context) CodeGenerator_Context {
	_ = currentContext
	if currentContext == StaticGet_CodeGenerator_Context_BEFORE_DANGLING_ELSE() {
		return StaticGet_CodeGenerator_Context_BEFORE_DANGLING_ELSE()
	} else {
		return StaticGet_CodeGenerator_Context_OTHER()
	}
}

// Java static method: CodeGenerator.getContextForNoInOperator(com.google.javascript.jscomp.CodeGenerator.Context context)
func CodeGenerator_GetContextForNoInOperator(context CodeGenerator_Context) CodeGenerator_Context {
	_ = context
	if context.InForInInitClause() {
		return context
	} else {
		return StaticGet_CodeGenerator_Context_OTHER()
	}
}

// Java static method: CodeGenerator.getContextForArrowFunctionBody(com.google.javascript.jscomp.CodeGenerator.Context context)
func CodeGenerator_GetContextForArrowFunctionBody(context CodeGenerator_Context) CodeGenerator_Context {
	_ = context
	if context.InForInInitClause() {
		return StaticGet_CodeGenerator_Context_START_OF_ARROW_FN_IN_FOR_INIT()
	} else {
		return StaticGet_CodeGenerator_Context_START_OF_ARROW_FN_BODY()
	}
}

// Java method: CodeGenerator.processEnd(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.CodeGenerator.Context context)
func (this_ *CodeGeneratorBase) priv_CodeGenerator_processEnd_1353(n pkg_rhino.Node, context CodeGenerator_Context) {
	_ = n
	_ = context
	switch n.TokenSwitchOrd_() {
	case int32(105) /* Token.CLASS */, int32(164) /* Token.INTERFACE */, int32(167) /* Token.ENUM */, int32(175) /* Token.NAMESPACE */ :
		this_.cc.EndClass((context == StaticGet_CodeGenerator_Context_STATEMENT()))
	case int32(73) /* Token.FUNCTION */ :
		if n.GetLastChild().IsEmpty() {
			this_.cc.EndStatement_2(true, false)
		} else {
			this_.cc.EndFunction((context == StaticGet_CodeGenerator_Context_STATEMENT()))
		}
	case int32(171) /* Token.DECLARE */ :
		if n.GetParent().RawToken_() != uint8(176) {
			this_.priv_CodeGenerator_processEnd_1353(n.GetFirstChild(), context)
		}
	case int32(120) /* Token.EXPORT */ :
		if (n.GetParent().RawToken_() != uint8(176)) && !(n.GetFirstChild().IsDeclare()) {
			this_.priv_CodeGenerator_processEnd_1353(n.GetFirstChild(), context)
		}
	case int32(129) /* Token.COMPUTED_PROP */ :
		if n.HasOneChild() {
			this_.cc.EndStatement_2(true, false)
		}
	case int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */ :
		if n.GetFirstChild().GetLastChild().IsEmpty() {
			this_.cc.EndStatement_2(true, false)
		}
	case int32(172) /* Token.MEMBER_VARIABLE_DEF */ :
		this_.cc.EndStatement_2(true, false)
	default:
		if context == StaticGet_CodeGenerator_Context_STATEMENT() {
			this_.cc.EndStatement_1(false)
		}
	}
}

// Java enum: com.google.javascript.jscomp.CodeGenerator.Context
type CodeGenerator_Context uint8

const (
	enum_CodeGenerator_Context_STATEMENT                     CodeGenerator_Context = 1
	enum_CodeGenerator_Context_BEFORE_DANGLING_ELSE          CodeGenerator_Context = 2
	enum_CodeGenerator_Context_START_OF_EXPR                 CodeGenerator_Context = 3
	enum_CodeGenerator_Context_IN_FOR_INIT_CLAUSE            CodeGenerator_Context = 4
	enum_CodeGenerator_Context_START_OF_ARROW_FN_BODY        CodeGenerator_Context = 5
	enum_CodeGenerator_Context_START_OF_ARROW_FN_IN_FOR_INIT CodeGenerator_Context = 6
	enum_CodeGenerator_Context_OTHER                         CodeGenerator_Context = 7
)

var enumBaseTable_CodeGenerator_Context = [...]j2gort.EnumBase{
	{Name_: "STATEMENT", Ordinal_: 0},
	{Name_: "BEFORE_DANGLING_ELSE", Ordinal_: 1},
	{Name_: "START_OF_EXPR", Ordinal_: 2},
	{Name_: "IN_FOR_INIT_CLAUSE", Ordinal_: 3},
	{Name_: "START_OF_ARROW_FN_BODY", Ordinal_: 4},
	{Name_: "START_OF_ARROW_FN_IN_FOR_INIT", Ordinal_: 5},
	{Name_: "OTHER", Ordinal_: 6},
}
var enumValues_CodeGenerator_Context = []CodeGenerator_Context{1, 2, 3, 4, 5, 6, 7}

var enumField_CodeGenerator_Context_inForInitClause = [...]bool{
	false,
	false,
	false,
	true,
	false,
	true,
	false,
}

func (this_ CodeGenerator_Context) FieldGet_inForInitClause() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumField_CodeGenerator_Context_inForInitClause[this_-1]
}

var enumField_CodeGenerator_Context_atArrowFnBody = [...]bool{
	false,
	false,
	false,
	false,
	true,
	true,
	false,
}

func (this_ CodeGenerator_Context) FieldGet_atArrowFnBody() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumField_CodeGenerator_Context_atArrowFnBody[this_-1]
}

func StaticGet_CodeGenerator_Context_STATEMENT() CodeGenerator_Context {
	return enum_CodeGenerator_Context_STATEMENT
}
func StaticGet_CodeGenerator_Context_BEFORE_DANGLING_ELSE() CodeGenerator_Context {
	return enum_CodeGenerator_Context_BEFORE_DANGLING_ELSE
}
func StaticGet_CodeGenerator_Context_START_OF_EXPR() CodeGenerator_Context {
	return enum_CodeGenerator_Context_START_OF_EXPR
}
func StaticGet_CodeGenerator_Context_IN_FOR_INIT_CLAUSE() CodeGenerator_Context {
	return enum_CodeGenerator_Context_IN_FOR_INIT_CLAUSE
}
func StaticGet_CodeGenerator_Context_START_OF_ARROW_FN_BODY() CodeGenerator_Context {
	return enum_CodeGenerator_Context_START_OF_ARROW_FN_BODY
}
func StaticGet_CodeGenerator_Context_START_OF_ARROW_FN_IN_FOR_INIT() CodeGenerator_Context {
	return enum_CodeGenerator_Context_START_OF_ARROW_FN_IN_FOR_INIT
}
func StaticGet_CodeGenerator_Context_OTHER() CodeGenerator_Context {
	return enum_CodeGenerator_Context_OTHER
}

func CodeGenerator_Context_Values() []CodeGenerator_Context {
	return j2gort.CloneSlice(enumValues_CodeGenerator_Context)
}

func CodeGenerator_Context_ValueOf(name string) CodeGenerator_Context {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_CodeGenerator_Context {
		if enumBaseTable_CodeGenerator_Context[i_].Name_ == name {
			return CodeGenerator_Context(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.CodeGenerator.Context." + name))
}

// Java method: Context.inForInInitClause()
func (this_ CodeGenerator_Context) InForInInitClause() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return this_.FieldGet_inForInitClause()
}

// Java method: Context.atArrowFunctionBody()
func (this_ CodeGenerator_Context) AtArrowFunctionBody() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return this_.FieldGet_atArrowFnBody()
}

func (this_ CodeGenerator_Context) IsJavaNil_() bool                              { return this_ == 0 }
func (CodeGenerator_Context) HasIdentityEquals_()                                 {}
func (CodeGenerator_Context) Implements_j2gort_Ext_java_lang_Enum()               {}
func (CodeGenerator_Context) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ CodeGenerator_Context) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (CodeGenerator_Context) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CodeGenerator$Context", "Context"
}
func (this_ CodeGenerator_Context) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_CodeGenerator_Context[this_-1]
}
func (this_ CodeGenerator_Context) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ CodeGenerator_Context) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ CodeGenerator_Context) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CodeGenerator_Context[this_-1].Name_
}
func (this_ CodeGenerator_Context) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_CodeGenerator_Context[this_-1].Name_
}
func (this_ CodeGenerator_Context) Equals(other_ any) bool {
	o_, ok_ := other_.(CodeGenerator_Context)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ CodeGenerator_Context) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_CodeGenerator_Context[this_-1].HashCode()
}
func (this_ CodeGenerator_Context) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_CodeGenerator_Context[this_-1].CompareTo(other_)
}
func (this_ CodeGenerator_Context) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_CodeGenerator_Context[this_-1].Name_)
}
func init() {
	vals_ := []any{CodeGenerator_Context(1), CodeGenerator_Context(2), CodeGenerator_Context(3), CodeGenerator_Context(4), CodeGenerator_Context(5), CodeGenerator_Context(6), CodeGenerator_Context(7)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CodeGenerator.Context", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.CodeGenerator$Context", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CodeGenerator.Context", "com.google.javascript.jscomp.CodeGenerator$Context", func(v_ any) bool {
		return j2gort.InstanceOf[CodeGenerator_Context](v_)
	})
}
