/*
 * Copyright 2013 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/parsing/IRFactory.java (generated mechanically by j2go).

package parsing

import (
	"fmt"
	"sync/atomic"

	pkg_dtoa "github.com/paullepoulpe/jsgo/transpiled/gen/dtoa"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
)

var lambdaSingleton_ir_factory_1 j2gort.Ext_java_util_function_ToIntFunction
var lambdaSingleton_ir_factory_2 j2gort.Ext_java_util_function_ToIntFunction
var lambdaSingleton_ir_factory_3 j2gort.Ext_java_util_function_ToIntFunction
var lambdaSingleton_ir_factory_4 j2gort.Ext_java_util_function_ToIntFunction

func init() {
	lambdaSingleton_ir_factory_1 = j2gort.NewLambda_Ext_java_util_function_ToIntFunction(func(lp_1 any) int32 {
		x := j2gort.CastPtr[pkg_parser.Comment](lp_1)
		_ = x
		return x.FieldGet_location().FieldGet_start().FieldGet_column()
	})
	lambdaSingleton_ir_factory_2 = j2gort.NewLambda_Ext_java_util_function_ToIntFunction(func(lp_2 any) int32 {
		x := j2gort.CastPtr[pkg_parser.Comment](lp_2)
		_ = x
		return x.FieldGet_location().FieldGet_start().FieldGet_line()
	})
	lambdaSingleton_ir_factory_3 = j2gort.NewLambda_Ext_java_util_function_ToIntFunction(func(mrp0_ any) int32 { return j2gort.CastPtr[IRFactory_LineAndColumn](mrp0_).Column() })
	lambdaSingleton_ir_factory_4 = j2gort.NewLambda_Ext_java_util_function_ToIntFunction(func(mrp0_ any) int32 { return j2gort.CastPtr[IRFactory_LineAndColumn](mrp0_).LineNo() })
}

// Java class: com.google.javascript.jscomp.parsing.IRFactory
type IRFactory struct {
	sourceFile                 pkg_rhino.StaticSourceFile                        `gson:"sourceFile"`
	sourceName                 string                                            `gson:"sourceName"`
	fileWithContent            *pkg_parser.SourceFile                            `gson:"fileWithContent"`
	config                     Config                                            `gson:"config"`
	errorReporter              pkg_rhino.ErrorReporter                           `gson:"errorReporter"`
	transformDispatcher        *IRFactory_TransformDispatcher                    `gson:"transformDispatcher"`
	reusableJsDocStream        *JsDocTokenStream                                 `gson:"reusableJsDocStream"`
	reusableJsDocParser        *JsDocInfoParser                                  `gson:"reusableJsDocParser"`
	reservedKeywords           j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"reservedKeywords"`
	licenseBuilder             *j2gort.Ext_java_util_LinkedHashSet               `gson:"licenseBuilder"`
	firstFileoverview          *pkg_rhino.JSDocInfo                              `gson:"firstFileoverview"`
	jsdocTracker               *IRFactory_CommentTracker                         `gson:"jsdocTracker"`
	nonJsdocTracker            *IRFactory_CommentTracker                         `gson:"nonJsdocTracker"`
	features                   *pkg_parser.FeatureSet                            `gson:"features"`
	closureUnawareCodeRanges   j2gort.Ext_com_google_common_collect_RangeSet     `gson:"closureUnawareCodeRanges"`
	scratchParamNames          []string                                          `gson:"scratchParamNames"`
	scratchParamTable          []int32                                           `gson:"scratchParamTable"`
	scratchParamSlots          []int32                                           `gson:"scratchParamSlots"`
	cachedParamValidator       j2gort.Ext_java_util_function_Consumer            `gson:"cachedParamValidator"`
	templateNode               pkg_rhino.Node                                    `gson:"templateNode"`
	defaultTemplateNode        pkg_rhino.Node                                    `gson:"defaultTemplateNode"`
	closureUnawareTemplateNode pkg_rhino.Node                                    `gson:"closureUnawareTemplateNode"`
	resultNode                 pkg_rhino.Node                                    `gson:"resultNode"`
	scratchParamCount          int32                                             `gson:"scratchParamCount"`
	idHash_                    atomic.Int32
	jsDocSourceKind            JsDocInfoParser_JsDocSourceKind `gson:"jsDocSourceKind"`
	currentFileIsExterns       bool                            `gson:"currentFileIsExterns"`
	nodeAlloc_                 nodestore.NodeAllocator
}

func (this_ *IRFactory) BindSelf_(s any) {
	_ = s
}
func (this_ *IRFactory) IsJavaNil_() bool { return this_ == nil }
func (*IRFactory) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.IRFactory", "IRFactory"
}
func (this_ *IRFactory) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*IRFactory) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.IRFactory", "com.google.javascript.jscomp.parsing.IRFactory", func(v_ any) bool {
		return j2gort.InstanceOf[*IRFactory](v_)
	})
}

func StaticGet_IRFactory_GETTER_ERROR_MESSAGE() string {
	return "getters are not supported in older versions of JavaScript. If you are targeting newer versions of JavaScript, set the appropriate language_in option."
}
func StaticGet_IRFactory_SETTER_ERROR_MESSAGE() string {
	return "setters are not supported in older versions of JavaScript. If you are targeting newer versions of JavaScript, set the appropriate language_in option."
}
func StaticGet_IRFactory_INVALID_ES3_PROP_NAME() string {
	return "Keywords and reserved words are not allowed as unquoted property names in older versions of JavaScript. If you are targeting newer versions of JavaScript, set the appropriate language_in option."
}
func StaticGet_IRFactory_INVALID_ES5_STRICT_OCTAL() string {
	return "Octal integer literals are not supported in strict mode."
}
func StaticGet_IRFactory_INVALID_OCTAL_DIGIT() string { return "Invalid octal digit in octal literal." }
func StaticGet_IRFactory_STRING_CONTINUATION_WARNING() string {
	return "String continuations are not recommended. See https://google.github.io/styleguide/jsguide.html#features-strings-no-line-continuations"
}
func StaticGet_IRFactory_OCTAL_STRING_LITERAL_WARNING() string {
	return "Octal literals in strings are not supported in this language mode."
}
func StaticGet_IRFactory_DUPLICATE_PARAMETER() string { return "Duplicate parameter name \"%s\"" }
func StaticGet_IRFactory_DUPLICATE_LABEL() string     { return "Duplicate label \"%s\"" }
func StaticGet_IRFactory_UNLABELED_BREAK() string {
	return "unlabelled break must be inside loop or switch"
}
func StaticGet_IRFactory_UNEXPECTED_CONTINUE() string { return "continue must be inside loop" }
func StaticGet_IRFactory_UNEXPECTED_LABELLED_CONTINUE() string {
	return "continue can only use labels of iteration statements"
}
func StaticGet_IRFactory_UNEXPECTED_RETURN() string { return "return must be inside function" }
func StaticGet_IRFactory_UNEXPECTED_YIELD() string  { return "yield must be inside generator function" }
func StaticGet_IRFactory_UNEXPECTED_AWAIT() string {
	return "await must be inside asynchronous function"
}
func StaticGet_IRFactory_UNEXPECTED_FOR_AWAIT_OF() string {
	return "'for-await-of' used in a non-async function context"
}
func StaticGet_IRFactory_UNEXPECTED_NEW_DOT_TARGET() string {
	return "new.target must be inside a function"
}
func StaticGet_IRFactory_NEW_DOT_TARGET_IN_CLASS_FIELD_OR_STATIC_BLOCK() string {
	return "new.target is not supported in class field initializers or static blocks"
}
func StaticGet_IRFactory_UNDEFINED_LABEL() string { return "undefined label \"%s\"" }

var static_IRFactory_ES5_RESERVED_KEYWORDS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_IRFactory_ES5_RESERVED_KEYWORDS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_IRFactory()
	return static_IRFactory_ES5_RESERVED_KEYWORDS
}
func StaticSet_IRFactory_ES5_RESERVED_KEYWORDS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_IRFactory()
	static_IRFactory_ES5_RESERVED_KEYWORDS = v
}

var static_IRFactory_ES5_STRICT_RESERVED_KEYWORDS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_IRFactory_ES5_STRICT_RESERVED_KEYWORDS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_IRFactory()
	return static_IRFactory_ES5_STRICT_RESERVED_KEYWORDS
}
func StaticSet_IRFactory_ES5_STRICT_RESERVED_KEYWORDS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_IRFactory()
	static_IRFactory_ES5_STRICT_RESERVED_KEYWORDS = v
}

var static_IRFactory_COMMENT_START_POSITION_COMPARATOR j2gort.Ext_java_util_Comparator

func StaticGet_IRFactory_COMMENT_START_POSITION_COMPARATOR() j2gort.Ext_java_util_Comparator {
	clinit_IRFactory()
	return static_IRFactory_COMMENT_START_POSITION_COMPARATOR
}
func StaticSet_IRFactory_COMMENT_START_POSITION_COMPARATOR(v j2gort.Ext_java_util_Comparator) {
	clinit_IRFactory()
	static_IRFactory_COMMENT_START_POSITION_COMPARATOR = v
}
func StaticGet_IRFactory_JSCOMP_CLOSURE_UNAWARE_CODE_SHADOW_HOST_NAME() string {
	return "$jscomp_wrap_closure_unaware_code"
}
func StaticGet_IRFactory_JSCOMP_CLOSURE_UNAWARE_SINK_VALUE() string {
	return "$jscomp_sink_closure_unaware_impl"
}

var static_IRFactory_GOOG_MODULE pkg_rhino.QualifiedName

func StaticGet_IRFactory_GOOG_MODULE() pkg_rhino.QualifiedName {
	clinit_IRFactory()
	return static_IRFactory_GOOG_MODULE
}
func StaticSet_IRFactory_GOOG_MODULE(v pkg_rhino.QualifiedName) {
	clinit_IRFactory()
	static_IRFactory_GOOG_MODULE = v
}

var clinitDone_IRFactory uint32
var clinitState_IRFactory j2gort.ClinitState

func clinit_IRFactory() {
	if atomic.LoadUint32(&clinitDone_IRFactory) != 0 {
		return
	}
	clinitSlow_IRFactory()
}

func clinitSlow_IRFactory() {
	if !j2gort.EnterClinit(&clinitState_IRFactory, &clinitDone_IRFactory) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_IRFactory, &clinitDone_IRFactory)
	static_IRFactory_ES5_RESERVED_KEYWORDS = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_7("class", "const", "enum", "export", "extends", "import", []any{"super"})
	static_IRFactory_ES5_STRICT_RESERVED_KEYWORDS = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_7("class", "const", "enum", "export", "extends", "import", []any{"super", "implements", "interface", "let", "package", "private", "protected", "public", "static", "yield"})
	static_IRFactory_COMMENT_START_POSITION_COMPARATOR = j2gort.Ext_java_util_Comparator_ComparingInt(lambdaSingleton_ir_factory_2).ThenComparingInt(lambdaSingleton_ir_factory_1)
	static_IRFactory_GOOG_MODULE = pkg_rhino.QualifiedName_Of("goog.module")
}

func (this_ *IRFactory) initFields_IRFactory() {
	_ = this_
	this_.sourceName = j2gort.NullString

	this_.features = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM()

}

// Java constructor: IRFactory(com.google.javascript.rhino.StaticSourceFile sourceFile, com.google.javascript.jscomp.parsing.Config config, com.google.javascript.rhino.ErrorReporter errorReporter, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.parsing.parser.trees.Comment> comments, com.google.javascript.jscomp.parsing.parser.SourceFile fileWithContent, com.google.javascript.jscomp.parsing.JsDocInfoParser.JsDocSourceKind jsDocSourceKind)
func NewIRFactory(sourceFile pkg_rhino.StaticSourceFile, config Config, errorReporter pkg_rhino.ErrorReporter, comments j2gort.Ext_com_google_common_collect_ImmutableList, fileWithContent *pkg_parser.SourceFile, jsDocSourceKind JsDocInfoParser_JsDocSourceKind) *IRFactory {
	clinit_IRFactory()
	this_ := &IRFactory{}
	InitIRFactory(this_, sourceFile, config, errorReporter, comments, fileWithContent, jsDocSourceKind)
	return this_
}

func InitIRFactory(this_ *IRFactory, sourceFile pkg_rhino.StaticSourceFile, config Config, errorReporter pkg_rhino.ErrorReporter, comments j2gort.Ext_com_google_common_collect_ImmutableList, fileWithContent *pkg_parser.SourceFile, jsDocSourceKind JsDocInfoParser_JsDocSourceKind) {
	_ = sourceFile
	_ = config
	_ = errorReporter
	_ = comments
	_ = fileWithContent
	_ = jsDocSourceKind
	this_.initFields_IRFactory()
	this_.nodeAlloc_.InitLocal()
	if comments != nil {
		{
			var i int32 = int32(0)
			_ = i
			var len_ int32 = comments.Size()
			_ = len_
			for ; i < len_; i = i + 1 {
				{
					r_3 := j2gort.CastPtr[pkg_parser.Comment](comments.Get(i))
					r_3.FieldSet_parsed(false)
				}
			}
		}
	}
	this_.jsdocTracker = NewIRFactory_CommentTracker(comments, true)
	this_.nonJsdocTracker = NewIRFactory_CommentTracker(comments, false)
	this_.sourceFile = sourceFile
	this_.fileWithContent = fileWithContent
	this_.defaultTemplateNode = this_.priv_IRFactory_createTemplateNode_6809()
	this_.templateNode = this_.defaultTemplateNode
	this_.nodeAlloc_.SetTemplate(nodestore.Node(this_.templateNode))
	this_.closureUnawareTemplateNode = this_.priv_IRFactory_createClosureUnawareTemplateNode_6810()
	if sourceFile == nil || sourceFile.IsJavaNil_() {
		this_.sourceName = j2gort.NullString
	} else {
		this_.sourceName = sourceFile.GetName()
	}
	this_.config = config
	this_.errorReporter = errorReporter
	this_.transformDispatcher = NewIRFactory_TransformDispatcher(this_)
	if config.StrictMode().IsStrict() {
		this_.reservedKeywords = StaticGet_IRFactory_ES5_STRICT_RESERVED_KEYWORDS()
	} else if config.LanguageMode() == StaticGet_Config_LanguageMode_ECMASCRIPT3() {
		this_.reservedKeywords = nil
	} else {
		this_.reservedKeywords = StaticGet_IRFactory_ES5_RESERVED_KEYWORDS()
	}
	this_.jsDocSourceKind = jsDocSourceKind
}

// Java method: IRFactory.createTemplateNode()
func (this_ *IRFactory) priv_IRFactory_createTemplateNode_6809() pkg_rhino.Node {
	var templateNode pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_SCRIPT())
	_ = templateNode
	templateNode.SetStaticSourceFile(this_.sourceFile)
	return templateNode
}

// Java method: IRFactory.createClosureUnawareTemplateNode()
func (this_ *IRFactory) priv_IRFactory_createClosureUnawareTemplateNode_6810() pkg_rhino.Node {
	var closureUnawareTemplateNode pkg_rhino.Node = this_.defaultTemplateNode.CloneNode_0()
	_ = closureUnawareTemplateNode
	closureUnawareTemplateNode.SetIsInClosureUnawareSubtree(true)
	return closureUnawareTemplateNode
}

// Java static method: IRFactory.transformTree(com.google.javascript.jscomp.parsing.parser.trees.ProgramTree tree, com.google.javascript.rhino.StaticSourceFile sourceFile, com.google.javascript.jscomp.parsing.Config config, com.google.javascript.rhino.ErrorReporter errorReporter, com.google.javascript.jscomp.parsing.parser.SourceFile file)
func IRFactory_TransformTree(tree *pkg_parser.ProgramTree, sourceFile pkg_rhino.StaticSourceFile, config Config, errorReporter pkg_rhino.ErrorReporter, file *pkg_parser.SourceFile) *IRFactory {
	clinit_IRFactory()
	_ = tree
	_ = sourceFile
	_ = config
	_ = errorReporter
	_ = file
	var jsDocSourceKind JsDocInfoParser_JsDocSourceKind
	_ = jsDocSourceKind
	if sourceFile.IsTypeScriptSource() {
		jsDocSourceKind = StaticGet_JsDocInfoParser_JsDocSourceKind_TSICKLE()
	} else {
		jsDocSourceKind = StaticGet_JsDocInfoParser_JsDocSourceKind_NORMAL()
	}
	clinit_IRFactory()
	var val_4 IRFactory
	InitIRFactory(&val_4, sourceFile, config, errorReporter, tree.FieldGet_sourceComments(), file, jsDocSourceKind)
	var irFactory *IRFactory = &val_4
	_ = irFactory
	defer irFactory.nodeAlloc_.Flush()
	var n pkg_rhino.Node = irFactory.transformDispatcher.Process(pkg_parser.ToIface_ParseTree(tree))
	_ = n
	irFactory.SetSourceInfo_2_Node_ParseTree(n, pkg_parser.ToIface_ParseTree(tree))
	if tree.FieldGet_sourceComments() != nil {
		{
			var i int32 = int32(0)
			_ = i
			var len_ int32 = tree.FieldGet_sourceComments().Size()
			_ = len_
			for ; i < len_; i = i + 1 {
				var comment *pkg_parser.Comment = j2gort.CastPtr[pkg_parser.Comment](tree.FieldGet_sourceComments().Get(i))
				_ = comment
				if ((comment.FieldGet_type_() == pkg_parser.StaticGet_Comment_Type_JSDOC()) || (comment.FieldGet_type_() == pkg_parser.StaticGet_Comment_Type_IMPORTANT())) && !(comment.FieldGet_parsed()) {
					var useLicensesOnlyConfig bool = irFactory.priv_IRFactory_withinClosureUnawareCodeRange_6842(comment.FieldGet_location().FieldGet_start().FieldGet_line(), comment.FieldGet_location().FieldGet_start().FieldGet_column())
					_ = useLicensesOnlyConfig
					irFactory.priv_IRFactory_handlePossibleFileOverviewJsDoc_6831(comment, useLicensesOnlyConfig)
				}
			}
		}
		if n.IsScript() {
			var endOfFilePos pkg_parser.SourcePosition = tree.FieldGet_location().FieldGet_end()
			_ = endOfFilePos
			var nonJSDocComment *pkg_rhino.NonJSDocComment = irFactory.priv_IRFactory_parseNonJSDocCommentAt_6838(endOfFilePos, false)
			_ = nonJSDocComment
			if nonJSDocComment != nil {
				n.SetTrailingNonJSDocComment(nonJSDocComment)
			}
		}
	}
	irFactory.priv_IRFactory_setFileOverviewJsDoc_6828(n)
	irFactory.priv_IRFactory_validateAll_6813(n)
	{
		r_5 := irFactory
		r_5.resultNode = n
	}
	return irFactory
}

// Java method: IRFactory.getResultNode()
func (this_ *IRFactory) GetResultNode() pkg_rhino.Node {
	return this_.resultNode
}

// Java method: IRFactory.getFeatures()
func (this_ *IRFactory) GetFeatures() *pkg_parser.FeatureSet {
	return this_.features
}

// Java method: IRFactory.validateAll(com.google.javascript.rhino.Node n)
func (this_ *IRFactory) priv_IRFactory_validateAll_6813(n pkg_rhino.Node) {
	_ = n
	var work j2gort.Deque[pkg_rhino.Node]
	_ = work
	for n != 0 {
		this_.priv_IRFactory_validate_6814(n)
		var nextSibling pkg_rhino.Node = n.GetNext()
		_ = nextSibling
		var firstChild pkg_rhino.Node = n.GetFirstChild()
		_ = firstChild
		if firstChild != 0 {
			if nextSibling != 0 {
				work.Push(j2gort.CheckNotNullNode(nextSibling))
			}
			n = firstChild
		} else if nextSibling != 0 {
			n = nextSibling
		} else {
			n = work.PollFirst()
		}
	}
	if !(work.IsEmpty()) {
		j2gort.CheckState(false)
	}
}

// Java method: IRFactory.validate(com.google.javascript.rhino.Node n)
func (this_ *IRFactory) priv_IRFactory_validate_6814(n pkg_rhino.Node) {
	_ = n
	this_.priv_IRFactory_validateParameters_6826(n)
	this_.priv_IRFactory_validateBreakContinue_6820(n)
	this_.priv_IRFactory_validateReturn_6818(n)
	this_.priv_IRFactory_validateYield_6817(n)
	this_.priv_IRFactory_validateAwait_6816(n)
	this_.priv_IRFactory_validateNewDotTarget_6819(n)
	this_.priv_IRFactory_validateLabel_6821(n)
	this_.priv_IRFactory_validateBlockScopedFunctions_6827(n)
	this_.priv_IRFactory_validateShadowHost_6815(n)
}

// Java method: IRFactory.validateShadowHost(com.google.javascript.rhino.Node n)
func (this_ *IRFactory) priv_IRFactory_validateShadowHost_6815(n pkg_rhino.Node) {
	_ = n
	if n.GetClosureUnawareShadow() == 0 {
		return
	}
	if (this_.firstFileoverview == nil) || !(this_.firstFileoverview.IsClosureUnawareCode()) {
		this_.priv_IRFactory_error_6876(n, "invalid @closureUnaware node: must be in a file annotated with @closureUnaware")
		return
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if parent.IsCall() {
		return
	}
	var grandparent pkg_rhino.Node = parent.GetParent()
	_ = grandparent
	if (parent.IsGetProp() && (grandparent != 0)) && grandparent.IsCall() {
		var calledFn string = parent.GetString()
		_ = calledFn
		if (j2gort.CheckNotNullString(calledFn) == "call") || (j2gort.CheckNotNullString(calledFn) == "apply") {
			return
		}
	}
	this_.priv_IRFactory_error_6876(n, "invalid @closureUnaware node: must be directly called")
}

// Java method: IRFactory.validateAwait(com.google.javascript.rhino.Node n)
func (this_ *IRFactory) priv_IRFactory_validateAwait_6816(n pkg_rhino.Node) {
	_ = n
	var errorMsg string = j2gort.NullString
	_ = errorMsg
	if n.IsAwait() {
		errorMsg = "await must be inside asynchronous function"
	} else if n.IsForAwaitOf() {
		errorMsg = "'for-await-of' used in a non-async function context"
	} else {
		return
	}
	var parent pkg_rhino.Node = n
	_ = parent
	for (func() pkg_rhino.Node { var asgn_6 pkg_rhino.Node = parent.GetParent(); parent = asgn_6; return asgn_6 }()) != 0 {
		if (parent.IsBlock() || parent.IsMemberFieldDef()) && parent.GetParent().IsClassMembers() {
			this_.errorReporter.M_Error(errorMsg, this_.sourceName, n.GetLineno(), n.GetCharno())
			return
		}
		if parent.IsAsyncFunction() {
			return
		}
		if parent.IsFunction() {
			this_.errorReporter.M_Error(errorMsg, this_.sourceName, n.GetLineno(), n.GetCharno())
			return
		}
		if parent.IsScript() && parent.GetBooleanProp(pkg_rhino.StaticGet_Node_ES6_MODULE()) {
			this_.MaybeWarnForFeature_2_Node_FeatureSet_Feature(n, pkg_parser.StaticGet_FeatureSet_Feature_TOP_LEVEL_AWAIT())
			return
		}
	}
	this_.errorReporter.M_Error(errorMsg, this_.sourceName, n.GetLineno(), n.GetCharno())
}

// Java method: IRFactory.validateYield(com.google.javascript.rhino.Node n)
func (this_ *IRFactory) priv_IRFactory_validateYield_6817(n pkg_rhino.Node) {
	_ = n
	if !(n.IsYield()) {
		return
	}
	var parent pkg_rhino.Node = n
	_ = parent
	for (func() pkg_rhino.Node { var asgn_7 pkg_rhino.Node = parent.GetParent(); parent = asgn_7; return asgn_7 }()) != 0 {
		if parent.IsClassMembers() {
			this_.errorReporter.M_Error("yield must be inside generator function", this_.sourceName, n.GetLineno(), n.GetCharno())
			return
		}
		if parent.IsGeneratorFunction() {
			return
		}
	}
}

// Java method: IRFactory.validateReturn(com.google.javascript.rhino.Node n)
func (this_ *IRFactory) priv_IRFactory_validateReturn_6818(n pkg_rhino.Node) {
	_ = n
	if !(n.IsReturn()) {
		return
	}
	var parent pkg_rhino.Node = n
	_ = parent
	for (func() pkg_rhino.Node { var asgn_8 pkg_rhino.Node = parent.GetParent(); parent = asgn_8; return asgn_8 }()) != 0 {
		if parent.IsClassMembers() {
			this_.errorReporter.M_Error("return must be inside function", this_.sourceName, n.GetLineno(), n.GetCharno())
			return
		}
		if parent.IsFunction() {
			return
		}
	}
	this_.errorReporter.M_Error("return must be inside function", this_.sourceName, n.GetLineno(), n.GetCharno())
}

// Java method: IRFactory.validateNewDotTarget(com.google.javascript.rhino.Node n)
func (this_ *IRFactory) priv_IRFactory_validateNewDotTarget_6819(n pkg_rhino.Node) {
	_ = n
	if n.RawToken_() != uint8(135) {
		return
	}
	var child pkg_rhino.Node = n
	_ = child
	{
		var parent pkg_rhino.Node = n.GetParent()
		_ = parent
		for ; parent != 0; func() { child = parent; parent = parent.GetParent() }() {
			if parent.IsFunction() && !(parent.IsArrowFunction()) {
				return
			}
			if (parent.IsMemberFieldDef() || (parent.IsComputedFieldDef() && child.IsSecondChildOf(parent))) || (parent.IsClassMembers() && child.IsBlock()) {
				this_.errorReporter.M_Error("new.target is not supported in class field initializers or static blocks", this_.sourceName, n.GetLineno(), n.GetCharno())
				return
			}
		}
	}
	this_.errorReporter.M_Error("new.target must be inside a function", this_.sourceName, n.GetLineno(), n.GetCharno())
}

// Java method: IRFactory.validateBreakContinue(com.google.javascript.rhino.Node n)
func (this_ *IRFactory) priv_IRFactory_validateBreakContinue_6820(n pkg_rhino.Node) {
	_ = n
	if !(n.IsBreak()) && !(n.IsContinue()) {
		return
	}
	var labelName pkg_rhino.Node = n.GetFirstChild()
	_ = labelName
	if labelName != 0 {
		var parent pkg_rhino.Node = n.GetParent()
		_ = parent
		for !(parent.IsLabel()) || !(IRFactory_LabelsMatch(parent, labelName)) {
			if (parent.IsFunction() || parent.IsScript()) || parent.IsClassMembers() {
				this_.errorReporter.M_Error(j2gort.Ext_java_lang_String_Format_2("undefined label \"%s\"", []any{j2gort.BoxString(labelName.GetString())}), this_.sourceName, n.GetLineno(), n.GetCharno())
				break
			}
			parent = parent.GetParent()
		}
		if parent.IsLabel() && IRFactory_LabelsMatch(parent, labelName) {
			if n.IsContinue() && !(IRFactory_IsContinueTarget(parent.GetLastChild())) {
				this_.errorReporter.M_Error("continue can only use labels of iteration statements", this_.sourceName, n.GetLineno(), n.GetCharno())
			}
		}
	} else if n.IsContinue() {
		var parent_2 pkg_rhino.Node = n.GetParent()
		_ = parent_2
		for !(IRFactory_IsContinueTarget(parent_2)) {
			if (parent_2.IsFunction() || parent_2.IsScript()) || parent_2.IsClassMembers() {
				this_.errorReporter.M_Error("continue must be inside loop", this_.sourceName, n.GetLineno(), n.GetCharno())
				break
			}
			parent_2 = parent_2.GetParent()
		}
	} else {
		var parent_3 pkg_rhino.Node = n.GetParent()
		_ = parent_3
		for !(IRFactory_IsBreakTarget(parent_3)) {
			if (parent_3.IsFunction() || parent_3.IsScript()) || parent_3.IsClassMembers() {
				this_.errorReporter.M_Error("unlabelled break must be inside loop or switch", this_.sourceName, n.GetLineno(), n.GetCharno())
				break
			}
			parent_3 = parent_3.GetParent()
		}
	}
}

// Java static method: IRFactory.isBreakTarget(com.google.javascript.rhino.Node n)
func IRFactory_IsBreakTarget(n pkg_rhino.Node) bool {
	clinit_IRFactory()
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(80) /* Token.FOR */, int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */, int32(78) /* Token.WHILE */, int32(79) /* Token.DO */, int32(75) /* Token.SWITCH */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: IRFactory.isContinueTarget(com.google.javascript.rhino.Node n)
func IRFactory_IsContinueTarget(n pkg_rhino.Node) bool {
	clinit_IRFactory()
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(80) /* Token.FOR */, int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */, int32(78) /* Token.WHILE */, int32(79) /* Token.DO */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java static method: IRFactory.labelsMatch(com.google.javascript.rhino.Node label, com.google.javascript.rhino.Node labelName)
func IRFactory_LabelsMatch(label pkg_rhino.Node, labelName pkg_rhino.Node) bool {
	clinit_IRFactory()
	_ = label
	_ = labelName
	{

		sr_9 := label.GetFirstChild().GetString()
		se_10 := labelName.GetString()
		return j2gort.CheckNotNullString(sr_9) == se_10
	}

}

// Java method: IRFactory.validateLabel(com.google.javascript.rhino.Node n)
func (this_ *IRFactory) priv_IRFactory_validateLabel_6821(n pkg_rhino.Node) {
	_ = n
	if !(n.IsLabel()) {
		return
	}
	var labelName pkg_rhino.Node = n.GetFirstChild()
	_ = labelName
	{
		var parent pkg_rhino.Node = n.GetParent()
		_ = parent
		for ; (parent != 0) && !(parent.IsFunction()); parent = parent.GetParent() {
			if parent.IsLabel() && IRFactory_LabelsMatch(parent, labelName) {
				this_.errorReporter.M_Error(j2gort.Ext_java_lang_String_Format_2("Duplicate label \"%s\"", []any{j2gort.BoxString(labelName.GetString())}), this_.sourceName, n.GetLineno(), n.GetCharno())
				break
			}
		}
	}
}

// Java static method: IRFactory.unwrapSimpleParamNameNode(com.google.javascript.rhino.Node c)
func IRFactory_UnwrapSimpleParamNameNode(c pkg_rhino.Node) pkg_rhino.Node {
	clinit_IRFactory()
	_ = c
	if c.IsName() {
		return c
	}
	if c.IsRest() || c.IsDefaultValue() {
		var first pkg_rhino.Node = c.GetFirstChild()
		_ = first
		if (first != 0) && first.IsName() {
			return first
		}
	}
	return pkg_rhino.Node(0)
}

// Java method: IRFactory.reportDuplicateParameter(com.google.javascript.rhino.Node param, java.lang.String paramName)
func (this_ *IRFactory) priv_IRFactory_reportDuplicateParameter_6822(param pkg_rhino.Node, paramName string) {
	_ = param
	_ = paramName
	this_.errorReporter.Warning(("Duplicate parameter name \"" + j2gort.StrValString(paramName) + "\""), this_.sourceName, param.GetLineno(), param.GetCharno())
}

// Java method: IRFactory.addScratchParamName(java.lang.String name)
func (this_ *IRFactory) priv_IRFactory_addScratchParamName_6823(name string) bool {
	_ = name
	var table []int32 = this_.scratchParamTable
	_ = table
	if table == nil {
		table = make([]int32, int32(32))
		this_.scratchParamTable = table
		this_.scratchParamNames = j2gort.NewStringArray(int32(24))
		this_.scratchParamSlots = make([]int32, int32(24))
	}
	var mask int32 = (int32(len(table)) - int32(1))
	_ = mask
	var h int32 = javalang.StringHashCode(j2gort.CheckNotNullString(name))
	_ = h
	var slot int32 = ((h ^ (javalang.Ushr32(h, int(int32(16))))) & mask)
	_ = slot
	for {
		var entry int32 = table[slot]
		_ = entry
		if entry == int32(0) {
			break
		}
		if j2gort.CheckNotNullString(this_.scratchParamNames[(entry-int32(1))]) == name {
			return false
		}
		slot = ((slot + int32(1)) & mask)
	}
	var idx int32 = this_.scratchParamCount
	_ = idx
	if idx == int32(len(this_.scratchParamNames)) {
		var newLen int32 = (idx * int32(2))
		_ = newLen
		this_.scratchParamNames = func() []string {
			arrSrc_11 := this_.scratchParamNames
			arrLen_12 := newLen
			j2gort.CheckNotNullSlice(arrSrc_11)
			if arrLen_12 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_13 := j2gort.NewStringArray(arrLen_12)
			copy(arrOut_13, arrSrc_11)
			return arrOut_13
		}()
		this_.scratchParamSlots = func() []int32 {
			arrSrc_14 := this_.scratchParamSlots
			arrLen_15 := newLen
			j2gort.CheckNotNullSlice(arrSrc_14)
			if arrLen_15 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_16 := make([]int32, arrLen_15)
			copy(arrOut_16, arrSrc_14)
			return arrOut_16
		}()
	}
	{
		a_17 := this_.scratchParamNames
		i_18 := idx
		a_17[i_18] = name
	}
	{
		a_19 := this_.scratchParamSlots
		i_20 := idx
		a_19[i_20] = slot
	}
	{
		a_21 := table
		i_22 := slot
		a_21[i_22] = (idx + int32(1))
	}
	this_.scratchParamCount = (idx + int32(1))
	if (this_.scratchParamCount * int32(4)) >= (int32(len(table)) * int32(3)) {
		var newCap int32 = (int32(len(table)) * int32(2))
		_ = newCap
		var newTable []int32 = make([]int32, newCap)
		_ = newTable
		var newMask int32 = (newCap - int32(1))
		_ = newMask
		{
			var i int32 = int32(0)
			_ = i
			for ; i < this_.scratchParamCount; i = i + 1 {
				var nh int32 = javalang.StringHashCode(j2gort.CheckNotNullString(this_.scratchParamNames[i]))
				_ = nh
				var nslot int32 = ((nh ^ (javalang.Ushr32(nh, int(int32(16))))) & newMask)
				_ = nslot
				for newTable[nslot] != int32(0) {
					nslot = ((nslot + int32(1)) & newMask)
				}
				{
					a_23 := newTable
					i_24 := nslot
					a_23[i_24] = (i + int32(1))
				}
				{
					a_25 := this_.scratchParamSlots
					i_26 := i
					a_25[i_26] = nslot
				}
			}
		}
		this_.scratchParamTable = newTable
	}
	return true
}

// Java method: IRFactory.clearScratchParamNames()
func (this_ *IRFactory) priv_IRFactory_clearScratchParamNames_6824() {
	var count int32 = this_.scratchParamCount
	_ = count
	if count > int32(0) {
		var table []int32 = this_.scratchParamTable
		_ = table
		var slots []int32 = this_.scratchParamSlots
		_ = slots
		var names []string = this_.scratchParamNames
		_ = names
		{
			var i int32 = int32(0)
			_ = i
			for ; i < count; i = i + 1 {
				{
					a_27 := table
					i_28 := slots[i]
					a_27[i_28] = int32(0)
				}
				{
					a_29 := names
					i_30 := i
					a_29[i_30] = j2gort.NullString
				}
			}
		}
		this_.scratchParamCount = int32(0)
	}
}

// Java method: IRFactory.validateParamNodeCallback(com.google.javascript.rhino.Node param)
func (this_ *IRFactory) priv_IRFactory_validateParamNodeCallback_6825(param pkg_rhino.Node) {
	_ = param
	var paramName string = param.GetString()
	_ = paramName
	if !(this_.priv_IRFactory_addScratchParamName_6823(paramName)) {
		this_.priv_IRFactory_reportDuplicateParameter_6822(param, paramName)
	}
}

// Java method: IRFactory.validateParameters(com.google.javascript.rhino.Node n)
func (this_ *IRFactory) priv_IRFactory_validateParameters_6826(n pkg_rhino.Node) {
	_ = n
	if !(n.IsParamList()) {
		return
	}
	var c0 pkg_rhino.Node = n.GetFirstChild()
	_ = c0
	if c0 == 0 {
		return
	}
	var c1 pkg_rhino.Node = c0.GetNext()
	_ = c1
	var p0 pkg_rhino.Node = IRFactory_UnwrapSimpleParamNameNode(c0)
	_ = p0
	if (c1 == 0) && (p0 != 0) {
		return
	}
	if (p0 != 0) && (c1 != 0) {
		var p1 pkg_rhino.Node = IRFactory_UnwrapSimpleParamNameNode(c1)
		_ = p1
		var c2 pkg_rhino.Node = c1.GetNext()
		_ = c2
		if (p1 != 0) && (c2 == 0) {
			var s0 string = p0.GetString()
			_ = s0
			var s1 string = p1.GetString()
			_ = s1
			if j2gort.CheckNotNullString(s0) == s1 {
				this_.priv_IRFactory_reportDuplicateParameter_6822(p1, s1)
			}
			return
		}
		if (p1 != 0) && (c2 != 0) {
			var p2 pkg_rhino.Node = IRFactory_UnwrapSimpleParamNameNode(c2)
			_ = p2
			if (p2 != 0) && (c2.GetNext() == 0) {
				var s0_2 string = p0.GetString()
				_ = s0_2
				var s1_2 string = p1.GetString()
				_ = s1_2
				var s2 string = p2.GetString()
				_ = s2
				if (!(j2gort.CheckNotNullString(s0_2) == s1_2) && !(j2gort.CheckNotNullString(s0_2) == s2)) && !(j2gort.CheckNotNullString(s1_2) == s2) {
					return
				}
			}
		}
	}
	var validator j2gort.Ext_java_util_function_Consumer = this_.cachedParamValidator
	_ = validator
	if validator == nil {
		{

			bound_31 := this_
			validator = j2gort.NewLambda_Ext_java_util_function_Consumer(func(mrp0_ any) {
				bound_31.priv_IRFactory_validateParamNodeCallback_6825(j2gort.CastNode[pkg_rhino.Node](mrp0_))
			})
		}

		this_.cachedParamValidator = validator
	}
	defer func() {
		this_.priv_IRFactory_clearScratchParamNames_6824()
	}()
	{
		{
			var c pkg_rhino.Node = c0
			_ = c
			for ; c != 0; c = c.GetNext() {
				ParsingUtil_GetParamOrPatternNames(c, validator)
			}
		}
	}
}

// Java method: IRFactory.validateBlockScopedFunctions(com.google.javascript.rhino.Node n)
func (this_ *IRFactory) priv_IRFactory_validateBlockScopedFunctions_6827(n pkg_rhino.Node) {
	_ = n
	if (n.IsFunction() && n.GetParent().IsBlock()) && !(n.GetGrandparent().IsFunction()) {
		this_.MaybeWarnForFeature_2_Node_FeatureSet_Feature(n, pkg_parser.StaticGet_FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION())
	}
}

// Java method: IRFactory.setFileOverviewJsDoc(com.google.javascript.rhino.Node irNode)
func (this_ *IRFactory) priv_IRFactory_setFileOverviewJsDoc_6828(irNode pkg_rhino.Node) {
	_ = irNode
	var fileoverview *pkg_rhino.JSDocInfo_Builder
	_ = fileoverview
	if this_.firstFileoverview == nil {
		fileoverview = pkg_rhino.JSDocInfo_S_Builder()
	} else {
		fileoverview = this_.firstFileoverview.ToBuilder_0()
	}
	if (this_.licenseBuilder != nil) && !(this_.licenseBuilder.IsEmpty()) {
		fileoverview.RecordLicense(j2gort.Ext_java_lang_String_Join_2_any_Ext_java_lang_Iterable("", j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.licenseBuilder)))
	}
	irNode.SetJSDocInfo(fileoverview.Build_1(false))
}

// Java method: IRFactory.transformBlock(com.google.javascript.jscomp.parsing.parser.trees.ParseTree node)
func (this_ *IRFactory) TransformBlock(node pkg_parser.ParseTree) pkg_rhino.Node {
	_ = node
	var irNode pkg_rhino.Node = this_.Transform(node)
	_ = irNode
	if irNode.IsBlock() {
		return irNode
	}
	var newBlock pkg_rhino.Node
	_ = newBlock
	if irNode.IsEmpty() {
		newBlock = this_.NewNode_1(pkg_rhino.StaticGet_Token_BLOCK())
	} else {
		newBlock = this_.NewNode_2(pkg_rhino.StaticGet_Token_BLOCK(), irNode)
	}
	this_.SetSourceInfo_2_Node_Node(newBlock, irNode)
	newBlock.SetIsAddedBlock(true)
	return newBlock
}

// Java method: IRFactory.handlePossibleFileOverviewJsDoc(com.google.javascript.jscomp.parsing.JsDocInfoParser jsDocParser)
func (this_ *IRFactory) priv_IRFactory_handlePossibleFileOverviewJsDoc_6830(jsDocParser *JsDocInfoParser) bool {
	_ = jsDocParser
	if jsDocParser.GetLicenseText() != j2gort.NullString {
		if this_.licenseBuilder == nil {
			this_.licenseBuilder = j2gort.NewExt_java_util_LinkedHashSet()
		}
		j2gort.SetAddStr(this_.licenseBuilder, jsDocParser.GetLicenseText())
	}
	var newFileoverview *pkg_rhino.JSDocInfo = jsDocParser.GetFileOverviewJSDocInfo()
	_ = newFileoverview
	if newFileoverview == this_.firstFileoverview {
		return false
	}
	if this_.firstFileoverview == nil {
		this_.firstFileoverview = newFileoverview
		this_.currentFileIsExterns = newFileoverview.IsExterns()
	} else {
		var merged *pkg_rhino.JSDocInfo_Builder = this_.firstFileoverview.ToBuilder_0()
		_ = merged
		if !(newFileoverview.GetSuppressions().IsEmpty()) {
			merged.RecordSuppressions_1(newFileoverview.GetSuppressions())
		}
		if newFileoverview.IsExterns() {
			this_.currentFileIsExterns = true
			merged.RecordExterns()
		}
		if newFileoverview.IsNoCoverage() {
			merged.RecordNoCoverage()
		}
		if newFileoverview.IsTypeSummary() {
			merged.RecordTypeSummary()
		}
		this_.firstFileoverview = merged.Build_0()
	}
	return true
}

// Java method: IRFactory.handlePossibleFileOverviewJsDoc(com.google.javascript.jscomp.parsing.parser.trees.Comment comment, boolean useLicensesOnlyConfig)
func (this_ *IRFactory) priv_IRFactory_handlePossibleFileOverviewJsDoc_6831(comment *pkg_parser.Comment, useLicensesOnlyConfig bool) {
	_ = comment
	_ = useLicensesOnlyConfig
	var jsDocParser *JsDocInfoParser = this_.priv_IRFactory_createJsDocInfoParser_6853(comment, useLicensesOnlyConfig)
	_ = jsDocParser
	defer func() {
		this_.priv_IRFactory_releaseJsDocInfoParser_6854(jsDocParser)
	}()
	{
		{
			r_32 := comment
			r_32.FieldSet_parsed(true)
		}
		this_.priv_IRFactory_handlePossibleFileOverviewJsDoc_6830(jsDocParser)
	}
}

// Java method: IRFactory.getJSDocCommentAt(com.google.javascript.jscomp.parsing.parser.util.SourcePosition pos)
func (this_ *IRFactory) priv_IRFactory_getJSDocCommentAt_6832(pos pkg_parser.SourcePosition) *pkg_parser.Comment {
	_ = pos
	var closestPreviousComment *pkg_parser.Comment = nil
	_ = closestPreviousComment
	for this_.jsdocTracker.HasPendingCommentBefore(pos) {
		closestPreviousComment = this_.jsdocTracker.Current()
		this_.jsdocTracker.Advance()
	}
	return closestPreviousComment
}

// Java method: IRFactory.parseJSDocInfoFrom(com.google.javascript.jscomp.parsing.parser.trees.Comment comment)
func (this_ *IRFactory) priv_IRFactory_parseJSDocInfoFrom_6833(comment *pkg_parser.Comment) *pkg_rhino.JSDocInfo {
	_ = comment
	if comment == nil {
		return nil
	}
	if this_.priv_IRFactory_withinClosureUnawareCodeRange_6842(comment.FieldGet_location().FieldGet_start().FieldGet_line(), comment.FieldGet_location().FieldGet_start().FieldGet_column()) {
		this_.jsdocTracker.Backtrack()
		return nil
	}
	var jsDocParser *JsDocInfoParser = this_.priv_IRFactory_createJsDocInfoParser_6852(comment)
	_ = jsDocParser
	defer func() {
		this_.priv_IRFactory_releaseJsDocInfoParser_6854(jsDocParser)
	}()
	{
		{
			r_33 := comment
			r_33.FieldSet_parsed(true)
		}
		if this_.priv_IRFactory_handlePossibleFileOverviewJsDoc_6830(jsDocParser) {
			return nil
		}
		return jsDocParser.RetrieveAndResetParsedJSDocInfo()
	}
	panic("unreachable")
}

// Java method: IRFactory.parseJSDocInfoOnTree(com.google.javascript.jscomp.parsing.parser.trees.ParseTree tree)
func (this_ *IRFactory) priv_IRFactory_parseJSDocInfoOnTree_6834(tree pkg_parser.ParseTree) *pkg_rhino.JSDocInfo {
	_ = tree
	base_tree_ := tree.GetParseTreeBase()
	switch (base_tree_.FieldGet_type_()).Ordinal() {
	case int32(6) /* ParseTreeType.EXPRESSION_STATEMENT */, int32(20) /* ParseTreeType.LABELLED_STATEMENT */, int32(64) /* ParseTreeType.EXPORT_DECLARATION */, int32(74) /* ParseTreeType.TEMPLATE_SUBSTITUTION */ :
		return nil
	case int32(47) /* ParseTreeType.CALL_EXPRESSION */, int32(40) /* ParseTreeType.CONDITIONAL_EXPRESSION */, int32(38) /* ParseTreeType.COMMA_EXPRESSION */, int32(39) /* ParseTreeType.BINARY_OPERATOR */, int32(43) /* ParseTreeType.MEMBER_EXPRESSION */, int32(50) /* ParseTreeType.MEMBER_LOOKUP_EXPRESSION */, int32(42) /* ParseTreeType.UPDATE_EXPRESSION */ :
		{
			var nearest pkg_parser.ParseTree = IRFactory_FindNearestNode(tree)
			_ = nearest
			if nearest.GetParseTreeBase().FieldGet_type_() == pkg_parser.StaticGet_ParseTreeType_PAREN_EXPRESSION() {
				return nil
			}
		}
	default:
	}
	return this_.priv_IRFactory_parseJSDocInfoFrom_6833(this_.priv_IRFactory_getJSDocCommentAt_6832(base_tree_.GetStart()))
}

// Java method: IRFactory.parseJSDocInfoOnToken(com.google.javascript.jscomp.parsing.parser.Token token)
func (this_ *IRFactory) ParseJSDocInfoOnToken(token pkg_parser.Token) *pkg_rhino.JSDocInfo {
	_ = token
	return this_.priv_IRFactory_parseJSDocInfoFrom_6833(this_.priv_IRFactory_getJSDocCommentAt_6832(token.GetStart()))
}

// Java method: IRFactory.parseInlineJSDocAt(com.google.javascript.jscomp.parsing.parser.util.SourcePosition pos)
func (this_ *IRFactory) ParseInlineJSDocAt(pos pkg_parser.SourcePosition) *pkg_rhino.JSDocInfo {
	_ = pos
	if this_.priv_IRFactory_withinClosureUnawareCodeRange_6842(pos.FieldGet_line(), pos.FieldGet_column()) {
		return nil
	}
	var comment *pkg_parser.Comment = this_.priv_IRFactory_getJSDocCommentAt_6832(pos)
	_ = comment
	if (comment != nil) && !(j2gort.JStrIndexOf(comment.FieldGet_value(), "@") >= 0) {
		return this_.priv_IRFactory_parseInlineTypeDoc_6855(comment)
	} else {
		return this_.priv_IRFactory_parseJSDocInfoFrom_6833(comment)
	}
}

// Java method: IRFactory.hasAnyPendingCommentBefore(com.google.javascript.jscomp.parsing.parser.util.SourcePosition pos, boolean withinClosureUnawareCodeRange)
func (this_ *IRFactory) priv_IRFactory_hasAnyPendingCommentBefore_6837(pos pkg_parser.SourcePosition, withinClosureUnawareCodeRange bool) bool {
	_ = pos
	_ = withinClosureUnawareCodeRange
	return (this_.nonJsdocTracker.HasPendingCommentBefore(pos) || (withinClosureUnawareCodeRange && this_.jsdocTracker.HasPendingCommentBefore(pos)))
}

// Java method: IRFactory.parseNonJSDocCommentAt(com.google.javascript.jscomp.parsing.parser.util.SourcePosition pos, boolean isInline)
func (this_ *IRFactory) priv_IRFactory_parseNonJSDocCommentAt_6838(pos pkg_parser.SourcePosition, isInline bool) *pkg_rhino.NonJSDocComment {
	_ = pos
	_ = isInline
	var withinClosureUnawareCodeRange bool = this_.priv_IRFactory_withinClosureUnawareCodeRange_6842(pos.FieldGet_line(), pos.FieldGet_column())
	_ = withinClosureUnawareCodeRange
	if !(this_.priv_IRFactory_hasAnyPendingCommentBefore_6837(pos, withinClosureUnawareCodeRange)) {
		return nil
	}
	var includeAllComments bool = (this_.config.JsDocParsingMode() == StaticGet_Config_JsDocParsing_INCLUDE_ALL_COMMENTS())
	_ = includeAllComments
	var result *j2gort.Ext_java_lang_StringBuilder = nil
	_ = result
	var firstComment *pkg_parser.Comment = nil
	_ = firstComment
	var lastComment *pkg_parser.Comment = nil
	_ = lastComment
	for this_.priv_IRFactory_hasAnyPendingCommentBefore_6837(pos, withinClosureUnawareCodeRange) {
		var currentComment *pkg_parser.Comment = this_.nonJsdocTracker.Current()
		_ = currentComment
		var commentFromJsdocTracker bool = false
		_ = commentFromJsdocTracker
		if withinClosureUnawareCodeRange {
			var currentJsDocComment *pkg_parser.Comment = this_.jsdocTracker.Current()
			_ = currentJsDocComment
			if (currentComment == nil) || ((currentJsDocComment != nil) && (StaticGet_IRFactory_COMMENT_START_POSITION_COMPARATOR().Compare(j2gort.PtrToAny(currentComment), j2gort.PtrToAny(currentJsDocComment)) > int32(0))) {
				if j2gort.JStrIndexOf(currentJsDocComment.FieldGet_value(), "@license") >= 0 {
					this_.jsdocTracker.Advance()
					continue
				}
				commentFromJsdocTracker = true
				currentComment = currentJsDocComment
			}
		}
		if firstComment == nil {
			firstComment = currentComment
		}
		if includeAllComments {
			if result == nil {
				result = j2gort.NewExt_java_lang_StringBuilder()
			}
			if lastComment != nil {
				{
					var blankCount int32 = (currentComment.FieldGet_location().FieldGet_start().FieldGet_line() - lastComment.FieldGet_location().FieldGet_end().FieldGet_line())
					_ = blankCount
					for ; blankCount > int32(0); blankCount = blankCount - 1 {
						result.Append_1_string("\n")
					}
				}
			}
			result.Append_1_string(currentComment.FieldGet_value())
		}
		lastComment = currentComment
		if commentFromJsdocTracker {
			{
				r_34 := currentComment
				r_34.FieldSet_parsed(true)
			}
			this_.jsdocTracker.Advance()
		} else {
			this_.nonJsdocTracker.Advance()
		}
	}
	if (firstComment == nil) || !(includeAllComments) {
		return nil
	}
	var nonJSDocComment *pkg_rhino.NonJSDocComment = pkg_rhino.NewNonJSDocComment(firstComment.FieldGet_location().FieldGet_start(), lastComment.FieldGet_location().FieldGet_end(), result.String())
	_ = nonJSDocComment
	nonJSDocComment.SetEndsAsLineComment((lastComment.FieldGet_type_() == pkg_parser.StaticGet_Comment_Type_LINE()))
	nonJSDocComment.SetIsInline(isInline)
	return nonJSDocComment
}

// Java method: IRFactory.parseTrailingNonJSDocCommentAt(com.google.javascript.jscomp.parsing.parser.util.SourcePosition tokenEnd, com.google.javascript.jscomp.parsing.parser.util.SourcePosition possibleNextTokenStart)
func (this_ *IRFactory) priv_IRFactory_parseTrailingNonJSDocCommentAt_6839(tokenEnd pkg_parser.SourcePosition, possibleNextTokenStart pkg_parser.SourcePosition) *pkg_rhino.NonJSDocComment {
	_ = tokenEnd
	_ = possibleNextTokenStart
	if this_.config.JsDocParsingMode() != StaticGet_Config_JsDocParsing_INCLUDE_ALL_COMMENTS() {
		return nil
	}
	if !(this_.nonJsdocTracker.HasPendingCommentBefore(possibleNextTokenStart)) {
		return nil
	}
	var comment *pkg_parser.Comment = this_.nonJsdocTracker.Current()
	_ = comment
	if tokenEnd.FieldGet_line() != comment.FieldGet_location().FieldGet_start().FieldGet_line() {
		return nil
	}
	if comment.FieldGet_location().FieldGet_start().FieldGet_offset() <= tokenEnd.FieldGet_offset() {
		return nil
	}
	var preCommentText string = j2gort.JStrSub(this_.fileWithContent.FieldGet_contents(), (tokenEnd.FieldGet_offset() + int32(1)), comment.FieldGet_location().FieldGet_start().FieldGet_offset())
	_ = preCommentText
	if !(j2gort.CheckNotNullString(javalang.StringTrim(j2gort.CheckNotNullString(preCommentText))) == "") {
		return nil
	}
	this_.nonJsdocTracker.Advance()
	var nonJSDocComment *pkg_rhino.NonJSDocComment = pkg_rhino.NewNonJSDocComment(comment.FieldGet_location().FieldGet_start(), comment.FieldGet_location().FieldGet_end(), comment.FieldGet_value())
	_ = nonJSDocComment
	nonJSDocComment.SetEndsAsLineComment((comment.FieldGet_type_() == pkg_parser.StaticGet_Comment_Type_LINE()))
	return nonJSDocComment
}

// Java static method: IRFactory.findNearestNode(com.google.javascript.jscomp.parsing.parser.trees.ParseTree tree)
func IRFactory_FindNearestNode(tree pkg_parser.ParseTree) pkg_parser.ParseTree {
	clinit_IRFactory()
	_ = tree
	for {
		switch (tree.GetParseTreeBase().FieldGet_type_()).Ordinal() {
		case int32(6) /* ParseTreeType.EXPRESSION_STATEMENT */ :
			tree = tree.GetParseTreeBase().AsExpressionStatement().FieldGet_expression()
			continue
		case int32(47) /* ParseTreeType.CALL_EXPRESSION */ :
			tree = tree.GetParseTreeBase().AsCallExpression().FieldGet_operand()
			continue
		case int32(39) /* ParseTreeType.BINARY_OPERATOR */ :
			tree = tree.GetParseTreeBase().AsBinaryOperator().FieldGet_left()
			continue
		case int32(40) /* ParseTreeType.CONDITIONAL_EXPRESSION */ :
			tree = tree.GetParseTreeBase().AsConditionalExpression().FieldGet_condition()
			continue
		case int32(43) /* ParseTreeType.MEMBER_EXPRESSION */ :
			tree = tree.GetParseTreeBase().AsMemberExpression().FieldGet_operand()
			continue
		case int32(50) /* ParseTreeType.MEMBER_LOOKUP_EXPRESSION */ :
			tree = tree.GetParseTreeBase().AsMemberLookupExpression().FieldGet_operand()
			continue
		case int32(42) /* ParseTreeType.UPDATE_EXPRESSION */ :
			tree = tree.GetParseTreeBase().AsUpdateExpression().FieldGet_operand()
			continue
		case int32(38) /* ParseTreeType.COMMA_EXPRESSION */ :
			tree = j2gort.AsIface[pkg_parser.ParseTree](tree.GetParseTreeBase().AsCommaExpression().FieldGet_expressions().Get(int32(0)))
			continue
		default:
			return tree
		}
	}
	panic("unreachable")
}

// Java method: IRFactory.transform(com.google.javascript.jscomp.parsing.parser.trees.ParseTree tree)
func (this_ *IRFactory) Transform(tree pkg_parser.ParseTree) pkg_rhino.Node {
	_ = tree
	var info *pkg_rhino.JSDocInfo = this_.priv_IRFactory_parseJSDocInfoOnTree_6834(tree)
	_ = info
	var comment *pkg_rhino.NonJSDocComment = this_.priv_IRFactory_parseNonJSDocCommentAt_6838(tree.GetParseTreeBase().GetStart(), false)
	_ = comment
	var isRootOfClosureUnawareSubTree bool = ((info != nil) && info.IsClosureUnawareCode())
	_ = isRootOfClosureUnawareSubTree
	if isRootOfClosureUnawareSubTree {
		this_.templateNode = this_.closureUnawareTemplateNode
		this_.nodeAlloc_.SetTemplate(nodestore.Node(this_.templateNode))
		var nextClosureUnawareRange pkg_parser.SourceRange = tree.GetParseTreeBase().FieldGet_location()
		_ = nextClosureUnawareRange
		var nextClosureUnawareRangeStart *IRFactory_LineAndColumn = IRFactory_LineAndColumn_FromSourcePosition(nextClosureUnawareRange.FieldGet_start())
		_ = nextClosureUnawareRangeStart
		var nextClosureUnawareRangeEnd *IRFactory_LineAndColumn = IRFactory_LineAndColumn_FromSourcePosition(nextClosureUnawareRange.FieldGet_end())
		_ = nextClosureUnawareRangeEnd
		if this_.closureUnawareCodeRanges == nil {
			this_.closureUnawareCodeRanges = j2gort.AsIface[j2gort.Ext_com_google_common_collect_RangeSet](j2gort.Ext_com_google_common_collect_TreeRangeSet_Create_0())
		}
		this_.closureUnawareCodeRanges.Add(j2gort.Ext_com_google_common_collect_Range_Closed(j2gort.PtrToAny(nextClosureUnawareRangeStart), j2gort.PtrToAny(nextClosureUnawareRangeEnd)))
	}
	var node pkg_rhino.Node = this_.transformDispatcher.Process(tree)
	_ = node
	if info != nil {
		node = this_.priv_IRFactory_maybeInjectCastNode_6843(tree, info, node)
		node.SetJSDocInfo(info)
	}
	if comment != nil {
		node.SetNonJSDocComment(comment)
	}
	if isRootOfClosureUnawareSubTree {
		this_.SetSourceInfo_2_Node_ParseTree(node, tree)
		this_.templateNode = this_.defaultTemplateNode
		this_.nodeAlloc_.SetTemplate(nodestore.Node(this_.templateNode))
		node = this_.priv_IRFactory_hideClosureUnawareCodeRoot_6841(node)
	}
	this_.SetSourceInfo_2_Node_ParseTree(node, tree)
	return node
}

// Java method: IRFactory.hideClosureUnawareCodeRoot(com.google.javascript.rhino.Node fn)
func (this_ *IRFactory) priv_IRFactory_hideClosureUnawareCodeRoot_6841(fn pkg_rhino.Node) pkg_rhino.Node {
	_ = fn
	if !(fn.IsFunction()) {
		this_.priv_IRFactory_error_6876(fn, "@closureUnaware root must be function")
		return fn
	}
	var info *pkg_rhino.JSDocInfo = fn.GetJSDocInfo()
	_ = info
	if info.GetPerFileClosureUnawareMode() != pkg_rhino.StaticGet_JSDocInfo_PerFileClosureUnawareMode_UNSPECIFIED() {
		this_.priv_IRFactory_error_6876(fn, "@closureUnaware mode can only be specified at the fileoverview level")
	}
	fn.SetJSDocInfo(nil)
	var shadowNameNode pkg_rhino.Node = pkg_rhino.IR_Name("$jscomp_wrap_closure_unaware_code")
	_ = shadowNameNode
	var sinkCall pkg_rhino.Node = pkg_rhino.IR_Call(pkg_rhino.IR_Name("$jscomp_sink_closure_unaware_impl").Srcref(this_.templateNode), []pkg_rhino.Node{fn}).Srcref(this_.templateNode)
	_ = sinkCall
	sinkCall.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
	var script pkg_rhino.Node = pkg_rhino.IR_Script_1_Node([]pkg_rhino.Node{pkg_rhino.IR_ExprResult(sinkCall).Srcref(this_.templateNode)})
	_ = script
	var shadowJsRoot pkg_rhino.Node = pkg_rhino.IR_Root([]pkg_rhino.Node{script})
	_ = shadowJsRoot
	script.SetStaticSourceFileFrom(this_.templateNode)
	script.SetIsInClosureUnawareSubtree(true)
	shadowNameNode.ClonePropsFrom(this_.templateNode)
	shadowNameNode.SetClosureUnawareShadow(shadowJsRoot)
	return shadowNameNode
}

// Java method: IRFactory.withinClosureUnawareCodeRange(int line, int lineColumnNo)
func (this_ *IRFactory) priv_IRFactory_withinClosureUnawareCodeRange_6842(line int32, lineColumnNo int32) bool {
	_ = line
	_ = lineColumnNo
	if (this_.closureUnawareCodeRanges == nil) || this_.closureUnawareCodeRanges.IsEmpty() {
		return false
	}
	return this_.closureUnawareCodeRanges.Contains(NewIRFactory_LineAndColumn(line, lineColumnNo))
}

// Java method: IRFactory.maybeInjectCastNode(com.google.javascript.jscomp.parsing.parser.trees.ParseTree node, com.google.javascript.rhino.JSDocInfo info, com.google.javascript.rhino.Node irNode)
func (this_ *IRFactory) priv_IRFactory_maybeInjectCastNode_6843(node pkg_parser.ParseTree, info *pkg_rhino.JSDocInfo, irNode pkg_rhino.Node) pkg_rhino.Node {
	_ = node
	_ = info
	_ = irNode
	if (node.GetParseTreeBase().FieldGet_type_() == pkg_parser.StaticGet_ParseTreeType_PAREN_EXPRESSION()) && info.HasType() {
		irNode = this_.NewNode_2(pkg_rhino.StaticGet_Token_CAST(), irNode)
	}
	return irNode
}

// Java method: IRFactory.transformNodeWithInlineComments(com.google.javascript.jscomp.parsing.parser.trees.ParseTree tree)
func (this_ *IRFactory) TransformNodeWithInlineComments(tree pkg_parser.ParseTree) pkg_rhino.Node {
	_ = tree
	var info *pkg_rhino.JSDocInfo = this_.ParseInlineJSDocAt(tree.GetParseTreeBase().GetStart())
	_ = info
	var comment *pkg_rhino.NonJSDocComment = this_.priv_IRFactory_parseNonJSDocCommentAt_6838(tree.GetParseTreeBase().GetStart(), true)
	_ = comment
	var node pkg_rhino.Node = this_.transformDispatcher.Process(tree)
	_ = node
	if info != nil {
		node.SetJSDocInfo(info)
	}
	if comment != nil {
		node.SetNonJSDocComment(comment)
	}
	this_.SetSourceInfo_2_Node_ParseTree(node, tree)
	return node
}

// Java static method: IRFactory.lineno(com.google.javascript.jscomp.parsing.parser.trees.ParseTree node)
func IRFactory_Lineno_1_ParseTree(node pkg_parser.ParseTree) int32 {
	clinit_IRFactory()
	_ = node
	return IRFactory_Lineno_1_SourcePosition(node.GetParseTreeBase().FieldGet_location().FieldGet_start())
}

// Java static method: IRFactory.lineno(com.google.javascript.jscomp.parsing.parser.Token token)
func IRFactory_Lineno_1_Token(token pkg_parser.Token) int32 {
	clinit_IRFactory()
	_ = token
	return IRFactory_Lineno_1_SourcePosition(token.FieldGet_location().FieldGet_start())
}

// Java static method: IRFactory.lineno(com.google.javascript.jscomp.parsing.parser.util.SourcePosition location)
func IRFactory_Lineno_1_SourcePosition(location pkg_parser.SourcePosition) int32 {
	clinit_IRFactory()
	_ = location
	return (location.FieldGet_line() + int32(1))
}

// Java static method: IRFactory.charno(com.google.javascript.jscomp.parsing.parser.trees.ParseTree node)
func IRFactory_Charno_1_ParseTree(node pkg_parser.ParseTree) int32 {
	clinit_IRFactory()
	_ = node
	return IRFactory_Charno_1_SourcePosition(node.GetParseTreeBase().FieldGet_location().FieldGet_start())
}

// Java static method: IRFactory.charno(com.google.javascript.jscomp.parsing.parser.Token token)
func IRFactory_Charno_1_Token(token pkg_parser.Token) int32 {
	clinit_IRFactory()
	_ = token
	return IRFactory_Charno_1_SourcePosition(token.FieldGet_location().FieldGet_start())
}

// Java static method: IRFactory.charno(com.google.javascript.jscomp.parsing.parser.util.SourcePosition location)
func IRFactory_Charno_1_SourcePosition(location pkg_parser.SourcePosition) int32 {
	clinit_IRFactory()
	_ = location
	return location.FieldGet_column()
}

// Java static method: IRFactory.languageFeatureWarningMessage(com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature feature)
func IRFactory_LanguageFeatureWarningMessage(feature *pkg_parser.FeatureSet_Feature) string {
	clinit_IRFactory()
	_ = feature
	var forFeature *Config_LanguageMode = Config_LanguageMode_MinimumRequiredFor(feature)
	_ = forFeature
	if forFeature == StaticGet_Config_LanguageMode_UNSUPPORTED() {
		return ("This language feature is not currently supported by the compiler: " + j2gort.StrVal(feature))
	} else {
		return ("This language feature is only supported for " + j2gort.StrVal(Config_LanguageMode_MinimumRequiredFor(feature)) + " mode or better: " + j2gort.StrVal(feature))
	}
}

// Java method: IRFactory.maybeWarnForFeature(com.google.javascript.jscomp.parsing.parser.trees.ParseTree node, com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature feature)
func (this_ *IRFactory) MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(node pkg_parser.ParseTree, feature *pkg_parser.FeatureSet_Feature) {
	_ = node
	_ = feature
	this_.features = this_.features.With_1_FeatureSet_Feature(feature)
	if this_.priv_IRFactory_withinClosureUnawareCodeRange_6842(node.GetParseTreeBase().FieldGet_location().FieldGet_start().FieldGet_line(), node.GetParseTreeBase().FieldGet_location().FieldGet_start().FieldGet_column()) {
		return
	}
	if this_.currentFileIsExterns {
		return
	}
	if !(this_.IsSupportedForInputLanguageMode(feature)) {
		this_.errorReporter.Warning(IRFactory_LanguageFeatureWarningMessage(feature), this_.sourceName, IRFactory_Lineno_1_ParseTree(node), IRFactory_Charno_1_ParseTree(node))
	}
}

// Java method: IRFactory.maybeWarnForFeature(com.google.javascript.jscomp.parsing.parser.Token token, com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature feature)
func (this_ *IRFactory) MaybeWarnForFeature_2_Token_FeatureSet_Feature(token pkg_parser.Token, feature *pkg_parser.FeatureSet_Feature) {
	_ = token
	_ = feature
	this_.features = this_.features.With_1_FeatureSet_Feature(feature)
	if this_.priv_IRFactory_withinClosureUnawareCodeRange_6842(token.FieldGet_location().FieldGet_start().FieldGet_line(), token.FieldGet_location().FieldGet_start().FieldGet_column()) {
		return
	}
	if this_.currentFileIsExterns {
		return
	}
	if !(this_.IsSupportedForInputLanguageMode(feature)) {
		this_.errorReporter.Warning(IRFactory_LanguageFeatureWarningMessage(feature), this_.sourceName, IRFactory_Lineno_1_Token(token), IRFactory_Charno_1_Token(token))
	}
}

// Java method: IRFactory.maybeWarnForFeature(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature feature)
func (this_ *IRFactory) MaybeWarnForFeature_2_Node_FeatureSet_Feature(node pkg_rhino.Node, feature *pkg_parser.FeatureSet_Feature) {
	_ = node
	_ = feature
	this_.features = this_.features.With_1_FeatureSet_Feature(feature)
	if this_.priv_IRFactory_withinClosureUnawareCodeRange_6842(node.GetLineno(), node.GetCharno()) {
		return
	}
	if this_.currentFileIsExterns {
		return
	}
	if !(this_.IsSupportedForInputLanguageMode(feature)) {
		this_.errorReporter.Warning(IRFactory_LanguageFeatureWarningMessage(feature), this_.sourceName, node.GetLineno(), node.GetCharno())
	}
}

// Java method: IRFactory.setSourceInfo(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node ref)
func (this_ *IRFactory) SetSourceInfo_2_Node_Node(node pkg_rhino.Node, ref pkg_rhino.Node) {
	_ = node
	_ = ref
	node.SetLinenoCharno(ref.GetLineno(), ref.GetCharno())
	this_.SetLengthFrom(node, ref)
}

// Java method: IRFactory.setSourceInfo(com.google.javascript.rhino.Node irNode, com.google.javascript.jscomp.parsing.parser.trees.ParseTree node)
func (this_ *IRFactory) SetSourceInfo_2_Node_ParseTree(irNode pkg_rhino.Node, node pkg_parser.ParseTree) {
	_ = irNode
	_ = node
	if irNode.GetLineno() == int32(-1) {
		this_.SetSourceInfo_3(irNode, node.GetParseTreeBase().FieldGet_location().FieldGet_start(), node.GetParseTreeBase().FieldGet_location().FieldGet_end())
	}
}

// Java method: IRFactory.setSourceInfo(com.google.javascript.rhino.Node irNode, com.google.javascript.jscomp.parsing.parser.Token token)
func (this_ *IRFactory) SetSourceInfo_2_Node_Token(irNode pkg_rhino.Node, token pkg_parser.Token) {
	_ = irNode
	_ = token
	this_.SetSourceInfo_3(irNode, token.FieldGet_location().FieldGet_start(), token.FieldGet_location().FieldGet_end())
}

// Java method: IRFactory.setSourceInfo(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.parsing.parser.util.SourcePosition start, com.google.javascript.jscomp.parsing.parser.util.SourcePosition end)
func (this_ *IRFactory) SetSourceInfo_3(node pkg_rhino.Node, start pkg_parser.SourcePosition, end pkg_parser.SourcePosition) {
	_ = node
	_ = start
	_ = end
	if node.GetLineno() == int32(-1) {
		node.SetLinenoCharno(IRFactory_Lineno_1_SourcePosition(start), IRFactory_Charno_1_SourcePosition(start))
		this_.SetLength(node, start, end)
	}
}

// Java method: IRFactory.createJsDocInfoParser(com.google.javascript.jscomp.parsing.parser.trees.Comment node)
func (this_ *IRFactory) priv_IRFactory_createJsDocInfoParser_6852(node *pkg_parser.Comment) *JsDocInfoParser {
	_ = node
	return this_.priv_IRFactory_createJsDocInfoParser_6853(node, false)
}

// Java method: IRFactory.createJsDocInfoParser(com.google.javascript.jscomp.parsing.parser.trees.Comment node, boolean useLicensesOnlyConfig)
func (this_ *IRFactory) priv_IRFactory_createJsDocInfoParser_6853(node *pkg_parser.Comment, useLicensesOnlyConfig bool) *JsDocInfoParser {
	_ = node
	_ = useLicensesOnlyConfig
	var comment string = node.FieldGet_value()
	_ = comment
	var lineno int32 = IRFactory_Lineno_1_SourcePosition(node.FieldGet_location().FieldGet_start())
	_ = lineno
	var charno int32 = IRFactory_Charno_1_SourcePosition(node.FieldGet_location().FieldGet_start())
	_ = charno
	var position int32 = node.FieldGet_location().FieldGet_start().FieldGet_offset()
	_ = position
	var config Config = this_.config
	_ = config
	if useLicensesOnlyConfig {
		config = config.ToBuilder().SetJsDocParsingMode(StaticGet_Config_JsDocParsing_LICENSE_COMMENTS_ONLY()).Build()
	}
	var numOpeningChars int32 = int32(3)
	_ = numOpeningChars
	var body string = j2gort.JStrSubFrom(comment, numOpeningChars)
	_ = body
	var stream *JsDocTokenStream = this_.reusableJsDocStream
	_ = stream
	if stream == nil {
		stream = NewJsDocTokenStream_3(body, lineno, (charno + numOpeningChars))
	} else {
		this_.reusableJsDocStream = nil
		stream.Reset(body, lineno, (charno + numOpeningChars))
	}
	var jsdocParser *JsDocInfoParser = this_.reusableJsDocParser
	_ = jsdocParser
	if jsdocParser == nil {
		jsdocParser = NewJsDocInfoParser(stream, comment, position, this_.templateNode, config, this_.jsDocSourceKind, this_.errorReporter)
	} else {
		this_.reusableJsDocParser = nil
		jsdocParser.Reset(stream, comment, position, this_.templateNode, config, this_.jsDocSourceKind, this_.errorReporter)
	}
	jsdocParser.SetFileOverviewJSDocInfo(this_.firstFileoverview)
	var completed bool = false
	_ = completed
	func() {
		defer func() {
			stream.ClearSource()
			this_.reusableJsDocStream = stream
			if !(completed) {
				this_.priv_IRFactory_releaseJsDocInfoParser_6854(jsdocParser)
			}
		}()
		if (node.FieldGet_type_() == pkg_parser.StaticGet_Comment_Type_IMPORTANT()) && (j2gort.JStrLen(node.FieldGet_value()) > int32(0)) {
			jsdocParser.ParseImportantComment()
		} else {
			jsdocParser.Parse()
		}
		completed = true
	}()
	return jsdocParser
}

// Java method: IRFactory.releaseJsDocInfoParser(com.google.javascript.jscomp.parsing.JsDocInfoParser parser)
func (this_ *IRFactory) priv_IRFactory_releaseJsDocInfoParser_6854(parser *JsDocInfoParser) {
	_ = parser
	parser.ClearForReuse()
	this_.reusableJsDocParser = parser
}

// Java method: IRFactory.parseInlineTypeDoc(com.google.javascript.jscomp.parsing.parser.trees.Comment node)
func (this_ *IRFactory) priv_IRFactory_parseInlineTypeDoc_6855(node *pkg_parser.Comment) *pkg_rhino.JSDocInfo {
	_ = node
	var comment string = node.FieldGet_value()
	_ = comment
	var lineno int32 = IRFactory_Lineno_1_SourcePosition(node.FieldGet_location().FieldGet_start())
	_ = lineno
	var charno int32 = IRFactory_Charno_1_SourcePosition(node.FieldGet_location().FieldGet_start())
	_ = charno
	var numOpeningChars int32 = int32(3)
	_ = numOpeningChars
	var body string = j2gort.JStrSubFrom(comment, numOpeningChars)
	_ = body
	var stream *JsDocTokenStream = this_.reusableJsDocStream
	_ = stream
	if stream == nil {
		stream = NewJsDocTokenStream_3(body, lineno, (charno + numOpeningChars))
	} else {
		this_.reusableJsDocStream = nil
		stream.Reset(body, lineno, (charno + numOpeningChars))
	}
	var parser *JsDocInfoParser = this_.reusableJsDocParser
	_ = parser
	if parser == nil {
		parser = NewJsDocInfoParser(stream, comment, node.FieldGet_location().FieldGet_start().FieldGet_offset(), this_.templateNode, this_.config, this_.jsDocSourceKind, this_.errorReporter)
	} else {
		this_.reusableJsDocParser = nil
		parser.Reset(stream, comment, node.FieldGet_location().FieldGet_start().FieldGet_offset(), this_.templateNode, this_.config, this_.jsDocSourceKind, this_.errorReporter)
	}
	defer func() {
		stream.ClearSource()
		this_.reusableJsDocStream = stream
		this_.priv_IRFactory_releaseJsDocInfoParser_6854(parser)
	}()
	{
		return parser.ParseInlineTypeDoc()
	}
	panic("unreachable")
}

// Java method: IRFactory.setLength(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.parsing.parser.util.SourcePosition start, com.google.javascript.jscomp.parsing.parser.util.SourcePosition end)
func (this_ *IRFactory) SetLength(node pkg_rhino.Node, start pkg_parser.SourcePosition, end pkg_parser.SourcePosition) {
	_ = node
	_ = start
	_ = end
	node.SetLength((end.FieldGet_offset() - start.FieldGet_offset()))
}

// Java method: IRFactory.setLengthFrom(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node ref)
func (this_ *IRFactory) SetLengthFrom(node pkg_rhino.Node, ref pkg_rhino.Node) {
	_ = node
	_ = ref
	node.SetLength(ref.GetLength())
}

// Java method: IRFactory.reportErrorIfYieldOrAwaitInDefaultValue(com.google.javascript.rhino.Node defaultValueNode)
func (this_ *IRFactory) priv_IRFactory_reportErrorIfYieldOrAwaitInDefaultValue_6858(defaultValueNode pkg_rhino.Node) {
	_ = defaultValueNode
	var yieldNode pkg_rhino.Node = IRFactory_FindNodeTypeInExpression(defaultValueNode, pkg_rhino.StaticGet_Token_YIELD())
	_ = yieldNode
	if yieldNode != 0 {
		this_.errorReporter.M_Error("`yield` is illegal in parameter default value.", yieldNode.GetSourceFileName(), yieldNode.GetLineno(), yieldNode.GetCharno())
	}
	var awaitNode pkg_rhino.Node = IRFactory_FindNodeTypeInExpression(defaultValueNode, pkg_rhino.StaticGet_Token_AWAIT())
	_ = awaitNode
	if awaitNode != 0 {
		this_.errorReporter.M_Error("`await` is illegal in parameter default value.", awaitNode.GetSourceFileName(), awaitNode.GetLineno(), awaitNode.GetCharno())
	}
}

// Java static method: IRFactory.findNodeTypeInExpression(com.google.javascript.rhino.Node expressionNode, com.google.javascript.rhino.Token token)
func IRFactory_FindNodeTypeInExpression(expressionNode pkg_rhino.Node, token pkg_rhino.Token) pkg_rhino.Node {
	clinit_IRFactory()
	_ = expressionNode
	_ = token
	var worklist j2gort.Deque[pkg_rhino.Node]
	_ = worklist
	worklist.AddLast(j2gort.CheckNotNullNode(expressionNode))
	for !(worklist.IsEmpty()) {
		var node pkg_rhino.Node = worklist.Pop()
		_ = node
		if node.GetToken() == token {
			return node
		} else if !(node.IsFunction()) && !(node.IsClass()) {
			{
				var child pkg_rhino.Node = node.GetFirstChild()
				_ = child
				for ; child != 0; child = child.GetNext() {
					worklist.AddLast(j2gort.CheckNotNullNode(child))
				}
			}
		}
	}
	return pkg_rhino.Node(0)
}

// Java method: IRFactory.normalizeRegex(com.google.javascript.jscomp.parsing.parser.Token token)
func (this_ *IRFactory) NormalizeRegex(token pkg_parser.Token) string {
	_ = token
	var value string = token.FieldGet_value()
	_ = value
	var lastSlash int32 = j2gort.JStrLastIndexOfCp(value, int32(uint16(47)))
	_ = lastSlash
	var cur int32 = j2gort.JStrIndexOfCp(value, int32(uint16(92)))
	_ = cur
	if cur == int32(-1) {
		return j2gort.JStrSub(value, int32(1), lastSlash)
	}
	var result j2gort.LocalStringBuilder

	var start int32 = int32(1)
	_ = start
	for cur != int32(-1) {
		result.AppendSub(j2gort.BoxString(value), start, cur)
		cur = cur + 1
		var c uint16 = j2gort.JStrCharAt(value, cur)
		_ = c
		switch c {
		case uint16(94), uint16(36), uint16(92), uint16(47), uint16(46), uint16(42), uint16(43), uint16(63), uint16(40), uint16(41), uint16(91), uint16(93), uint16(123), uint16(125), uint16(124), uint16(45), uint16(98), uint16(66), uint16(99), uint16(100), uint16(68), uint16(102), uint16(110), uint16(107), uint16(112), uint16(80), uint16(114), uint16(115), uint16(83), uint16(116), uint16(117), uint16(118), uint16(119), uint16(87), uint16(120), uint16(48), uint16(49), uint16(50), uint16(51), uint16(52), uint16(53), uint16(54), uint16(55), uint16(56), uint16(57):
			result.AppendChar(uint16(92))
			fallthrough
		default:
			result.AppendChar(c)
		}
		start = (cur + int32(1))
		cur = j2gort.JStrIndexOfCpFrom(value, int32(uint16(92)), start)
	}
	result.AppendSub(j2gort.BoxString(value), start, lastSlash)
	return result.String()
}

// Java method: IRFactory.normalizeString(com.google.javascript.jscomp.parsing.parser.Token token, boolean templateLiteral)
func (this_ *IRFactory) NormalizeString(token pkg_parser.Token, templateLiteral bool) string {
	_ = token
	_ = templateLiteral
	var value string = token.FieldGet_value()
	_ = value
	if j2gort.JStrIndexOfCp(value, int32(uint16(13))) >= int32(0) {
		value = j2gort.JStrReplaceAll(value, "\r\n?", "\n")
	}
	var start int32
	_ = start
	if templateLiteral {
		start = int32(0)
	} else {
		start = int32(1)
	}
	var cur int32 = j2gort.JStrIndexOfCp(value, int32(uint16(92)))
	_ = cur
	if cur == int32(-1) {
		if templateLiteral {
			return value
		} else {
			return j2gort.JStrSub(value, int32(1), (j2gort.JStrLen(value) - int32(1)))
		}
	}
	var result j2gort.LocalStringBuilder

	for cur != int32(-1) {
		result.AppendSub(j2gort.BoxString(value), start, cur)
		cur = (cur + int32(1))
		var c uint16 = j2gort.JStrCharAt(value, cur)
		_ = c
		{
			var numDigits int32 = int32(0)
			_ = numDigits
			var escapeEnd int32 = int32(0)
			_ = escapeEnd
			var hexDigits string = j2gort.NullString
			_ = hexDigits
			var codePointValue int32 = int32(0)
			_ = codePointValue
			switch c {
			case uint16(98):
				result.AppendChar(uint16(8))
			case uint16(102):
				result.AppendChar(uint16(12))
			case uint16(110):
				result.AppendChar(uint16(10))
			case uint16(114):
				result.AppendChar(uint16(13))
			case uint16(116):
				result.AppendChar(uint16(9))
			case uint16(118):
				result.AppendChar(uint16(11))
			case uint16(10), uint16(8232), uint16(8233):
				this_.MaybeWarnForFeature_2_Token_FeatureSet_Feature(token, pkg_parser.StaticGet_FeatureSet_Feature_STRING_CONTINUATION())
				this_.errorReporter.Warning("String continuations are not recommended. See https://google.github.io/styleguide/jsguide.html#features-strings-no-line-continuations", this_.sourceName, IRFactory_Lineno_1_SourcePosition(token.FieldGet_location().FieldGet_start()), IRFactory_Charno_1_SourcePosition(token.FieldGet_location().FieldGet_start()))
			case uint16(48), uint16(49), uint16(50), uint16(51), uint16(52), uint16(53), uint16(54), uint16(55):
				if ((cur + int32(1)) < j2gort.JStrLen(value)) && IRFactory_IsOctalDigit(j2gort.JStrCharAt(value, (cur+int32(1)))) {
					if ((int32(c) <= int32(uint16(51))) && ((cur + int32(2)) < j2gort.JStrLen(value))) && IRFactory_IsOctalDigit(j2gort.JStrCharAt(value, (cur+int32(2)))) {
						numDigits = int32(3)
					} else {
						numDigits = int32(2)
					}
				} else {
					numDigits = int32(1)
				}
				if this_.priv_IRFactory_inStrictContext_6862() || templateLiteral {
					if (int32(c) == int32(uint16(48))) && (numDigits == int32(1)) {
					} else {
						this_.errorReporter.Warning("Octal literals in strings are not supported in this language mode.", this_.sourceName, IRFactory_Lineno_1_SourcePosition(token.FieldGet_location().FieldGet_start()), IRFactory_Charno_1_SourcePosition(token.FieldGet_location().FieldGet_start()))
					}
				}
				result.AppendChar(uint16(j2gort.ParseInt(j2gort.JStrSub(value, cur, (cur+numDigits)), int32(8))))
				cur = (cur + (numDigits - int32(1)))
			case uint16(120):
				result.AppendChar(uint16(((IRFactory_Hexdigit(j2gort.JStrCharAt(value, (cur+int32(1)))) * int32(16)) + IRFactory_Hexdigit(j2gort.JStrCharAt(value, (cur+int32(2)))))))
				cur = (cur + int32(2))
			case uint16(117):
				if int32(j2gort.JStrCharAt(value, (cur+int32(1)))) != int32(uint16(123)) {
					escapeEnd = (cur + int32(5))
					hexDigits = j2gort.JStrSub(value, (cur + int32(1)), escapeEnd)
				} else {
					escapeEnd = (cur + int32(2))
					for j2gort.BoxChar_Digit_2_uint16_int32(j2gort.JStrCharAt(value, escapeEnd), int32(16)) >= int32(0) {
						escapeEnd = escapeEnd + 1
					}
					hexDigits = j2gort.JStrSub(value, (cur + int32(2)), escapeEnd)
					escapeEnd = escapeEnd + 1
				}
				codePointValue = j2gort.ParseInt(hexDigits, int32(16))
				if codePointValue > int32(1114111) {
					this_.errorReporter.M_Error("Undefined Unicode code-point", this_.sourceName, IRFactory_Lineno_1_SourcePosition(token.FieldGet_location().FieldGet_start()), IRFactory_Charno_1_SourcePosition(token.FieldGet_location().FieldGet_start()))
					result.AppendStr("\\u{")
					result.AppendStr(hexDigits)
					result.AppendStr("}")
				} else {
					result.AppendChars(j2gort.BoxChar_ToChars_1(codePointValue))
				}
				cur = (escapeEnd - int32(1))
			case uint16(39), uint16(34), uint16(92):
				fallthrough
			default:
				result.AppendChar(c)
			}
		}
		start = (cur + int32(1))
		cur = j2gort.JStrIndexOfCpFrom(value, int32(uint16(92)), start)
	}
	result.AppendSub(j2gort.BoxString(value), start, func() int32 {
		if templateLiteral {
			return j2gort.JStrLen(value)
		}
		return (j2gort.JStrLen(value) - int32(1))
	}())
	return result.String()
}

// Java method: IRFactory.isSupportedForInputLanguageMode(com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature feature)
func (this_ *IRFactory) IsSupportedForInputLanguageMode(feature *pkg_parser.FeatureSet_Feature) bool {
	_ = feature
	return this_.config.LanguageMode().featureSet.Has(feature)
}

// Java method: IRFactory.inStrictContext()
func (this_ *IRFactory) priv_IRFactory_inStrictContext_6862() bool {
	return this_.config.StrictMode().IsStrict()
}

// Java method: IRFactory.normalizeNumber(com.google.javascript.jscomp.parsing.parser.Token token)
func (this_ *IRFactory) NormalizeNumber(token pkg_parser.Token) float64 {
	_ = token
	var value string = token.FieldGet_value()
	_ = value
	if j2gort.JStrIndexOf(value, "_") >= 0 {
		value = this_.priv_IRFactory_removeNumericSeparators_6865(value, token)
	}
	var location pkg_parser.SourceRange = token.FieldGet_location()
	_ = location
	var length int32 = j2gort.JStrLen(value)
	_ = length
	if !(length > int32(0)) {
		j2gort.CheckState(false)
	}
	if !((int32(j2gort.JStrCharAt(value, int32(0))) != int32(uint16(45))) && (int32(j2gort.JStrCharAt(value, int32(0))) != int32(uint16(43)))) {
		j2gort.CheckState(false)
	}
	if int32(j2gort.JStrCharAt(value, int32(0))) == int32(uint16(46)) {
		return j2gort.ParseDouble(("0" + j2gort.StrValString(value)))
	} else if (int32(j2gort.JStrCharAt(value, int32(0))) == int32(uint16(48))) && (length > int32(1)) {
		switch j2gort.JStrCharAt(value, int32(1)) {
		case uint16(46), uint16(101), uint16(69):
			return j2gort.ParseDouble(value)
		case uint16(98), uint16(66):
			this_.MaybeWarnForFeature_2_Token_FeatureSet_Feature(token, pkg_parser.StaticGet_FeatureSet_Feature_BINARY_LITERALS())
			var v float64 = float64(int32(0))
			_ = v
			var c int32 = int32(1)
			_ = c
			for func() int32 { old_35 := c; nw_36 := old_35 + 1; c = nw_36; return nw_36 }() < length {
				v = ((v * float64(int32(2))) + float64(IRFactory_Binarydigit(j2gort.JStrCharAt(value, c))))
			}
			return IRFactory_RoundLargeRadixLiteral(v, j2gort.JStrSubFrom(value, int32(2)), int32(2))
		case uint16(111), uint16(79):
			this_.MaybeWarnForFeature_2_Token_FeatureSet_Feature(token, pkg_parser.StaticGet_FeatureSet_Feature_OCTAL_LITERALS())
			var v_2 float64 = float64(int32(0))
			_ = v_2
			var c_2 int32 = int32(1)
			_ = c_2
			for func() int32 { old_37 := c_2; nw_38 := old_37 + 1; c_2 = nw_38; return nw_38 }() < length {
				v_2 = ((v_2 * float64(int32(8))) + float64(IRFactory_Octaldigit(j2gort.JStrCharAt(value, c_2))))
			}
			return IRFactory_RoundLargeRadixLiteral(v_2, j2gort.JStrSubFrom(value, int32(2)), int32(8))
		case uint16(120), uint16(88):
			var v_3 float64 = float64(int32(0))
			_ = v_3
			var c_3 int32 = int32(1)
			_ = c_3
			for func() int32 { old_39 := c_3; nw_40 := old_39 + 1; c_3 = nw_40; return nw_40 }() < length {
				v_3 = ((v_3 * float64(int32(16))) + float64(IRFactory_Hexdigit(j2gort.JStrCharAt(value, c_3))))
			}
			return IRFactory_RoundLargeRadixLiteral(v_3, j2gort.JStrSubFrom(value, int32(2)), int32(16))
		case uint16(48), uint16(49), uint16(50), uint16(51), uint16(52), uint16(53), uint16(54), uint16(55):
			var v_4 float64 = float64(int32(0))
			_ = v_4
			var c_4 int32 = int32(0)
			_ = c_4
			for func() int32 { old_41 := c_4; nw_42 := old_41 + 1; c_4 = nw_42; return nw_42 }() < length {
				var digit uint16 = j2gort.JStrCharAt(value, c_4)
				_ = digit
				if IRFactory_IsOctalDigit(digit) {
					v_4 = ((v_4 * float64(int32(8))) + float64(IRFactory_Octaldigit(digit)))
				} else {
					this_.errorReporter.M_Error("Invalid octal digit in octal literal.", this_.sourceName, IRFactory_Lineno_1_SourcePosition(location.FieldGet_start()), IRFactory_Charno_1_SourcePosition(location.FieldGet_start()))
					return float64(int32(0))
				}
			}
			if this_.priv_IRFactory_inStrictContext_6862() {
				this_.errorReporter.M_Error("Octal integer literals are not supported in strict mode.", this_.sourceName, IRFactory_Lineno_1_SourcePosition(location.FieldGet_start()), IRFactory_Charno_1_SourcePosition(location.FieldGet_start()))
			} else {
				this_.errorReporter.Warning("Octal integer literals are not supported in strict mode.", this_.sourceName, IRFactory_Lineno_1_SourcePosition(location.FieldGet_start()), IRFactory_Charno_1_SourcePosition(location.FieldGet_start()))
			}
			return IRFactory_RoundLargeRadixLiteral(v_4, j2gort.JStrSubFrom(value, int32(1)), int32(8))
		case uint16(56), uint16(57):
			this_.errorReporter.M_Error("Invalid octal digit in octal literal.", this_.sourceName, IRFactory_Lineno_1_SourcePosition(location.FieldGet_start()), IRFactory_Charno_1_SourcePosition(location.FieldGet_start()))
			return float64(int32(0))
		default:
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected character in number literal: "+j2gort.StrValChar(j2gort.JStrCharAt(value, int32(1))))))
		}
		panic("unreachable switch expression")
	} else {
		return j2gort.ParseDouble(value)
	}
}

// Java method: IRFactory.normalizeBigInt(com.google.javascript.jscomp.parsing.parser.Token token)
func (this_ *IRFactory) NormalizeBigInt(token pkg_parser.Token) *j2gort.Ext_java_math_BigInteger {
	_ = token
	var value string = token.FieldGet_value()
	_ = value
	value = j2gort.JStrSub(value, int32(0), j2gort.JStrIndexOfCp(value, int32(uint16(110))))
	if j2gort.JStrIndexOf(value, "_") >= 0 {
		value = this_.priv_IRFactory_removeNumericSeparators_6865(value, token)
	}
	var length int32 = j2gort.JStrLen(value)
	_ = length
	if !(length > int32(0)) {
		j2gort.CheckState(false)
	}
	if !((int32(j2gort.JStrCharAt(value, int32(0))) != int32(uint16(45))) && (int32(j2gort.JStrCharAt(value, int32(0))) != int32(uint16(43)))) {
		j2gort.CheckState(false)
	}
	if (int32(j2gort.JStrCharAt(value, int32(0))) == int32(uint16(48))) && (length > int32(1)) {
		switch j2gort.JStrCharAt(value, int32(1)) {
		case uint16(98), uint16(66):
			this_.MaybeWarnForFeature_2_Token_FeatureSet_Feature(token, pkg_parser.StaticGet_FeatureSet_Feature_BINARY_LITERALS())
			return j2gort.NewExt_java_math_BigInteger_2_string_int32(j2gort.JStrSubFrom(value, int32(2)), int32(2))
		case uint16(111), uint16(79):
			this_.MaybeWarnForFeature_2_Token_FeatureSet_Feature(token, pkg_parser.StaticGet_FeatureSet_Feature_OCTAL_LITERALS())
			return j2gort.NewExt_java_math_BigInteger_2_string_int32(j2gort.JStrSubFrom(value, int32(2)), int32(8))
		case uint16(120), uint16(88):
			return j2gort.NewExt_java_math_BigInteger_2_string_int32(j2gort.JStrSubFrom(value, int32(2)), int32(16))
		case uint16(48), uint16(49), uint16(50), uint16(51), uint16(52), uint16(53), uint16(54), uint16(55), uint16(56), uint16(57):
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Nonzero BigInts can't have a leading zero")))
		default:
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected character in bigint literal: "+j2gort.StrValChar(j2gort.JStrCharAt(value, int32(1))))))
		}
		panic("unreachable switch expression")
	} else {
		return j2gort.NewExt_java_math_BigInteger_1_string(value)
	}
}

// Java method: IRFactory.removeNumericSeparators(java.lang.String value, com.google.javascript.jscomp.parsing.parser.Token token)
func (this_ *IRFactory) priv_IRFactory_removeNumericSeparators_6865(value string, token pkg_parser.Token) string {
	_ = value
	_ = token
	this_.MaybeWarnForFeature_2_Token_FeatureSet_Feature(token, pkg_parser.StaticGet_FeatureSet_Feature_NUMERIC_SEPARATOR())
	return javalang.StringReplace(value, "_", "")
}

// Java static method: IRFactory.roundLargeRadixLiteral(double v, java.lang.String digits, int radix)
func IRFactory_RoundLargeRadixLiteral(v float64, digits string, radix int32) float64 {
	clinit_IRFactory()
	_ = v
	_ = digits
	_ = radix
	if v < float64(9.007199254740992e15) {
		return v
	} else {
		return j2gort.NewExt_java_math_BigInteger_2_string_int32(digits, radix).DoubleValue()
	}
}

// Java static method: IRFactory.binarydigit(char c)
func IRFactory_Binarydigit(c uint16) int32 {
	clinit_IRFactory()
	_ = c
	if (int32(c) >= int32(uint16(48))) && (int32(c) <= int32(uint16(49))) {
		return (int32(c) - int32(uint16(48)))
	}
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "unexpected: "+j2gort.StrValChar(c))))
}

// Java static method: IRFactory.isOctalDigit(char c)
func IRFactory_IsOctalDigit(c uint16) bool {
	clinit_IRFactory()
	_ = c
	return ((int32(c) >= int32(uint16(48))) && (int32(c) <= int32(uint16(55))))
}

// Java static method: IRFactory.octaldigit(char c)
func IRFactory_Octaldigit(c uint16) int32 {
	clinit_IRFactory()
	_ = c
	if IRFactory_IsOctalDigit(c) {
		return (int32(c) - int32(uint16(48)))
	}
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "unexpected: "+j2gort.StrValChar(c))))
}

// Java static method: IRFactory.hexdigit(char c)
func IRFactory_Hexdigit(c uint16) int32 {
	clinit_IRFactory()
	_ = c
	switch c {
	case uint16(48):
		return int32(0)
	case uint16(49):
		return int32(1)
	case uint16(50):
		return int32(2)
	case uint16(51):
		return int32(3)
	case uint16(52):
		return int32(4)
	case uint16(53):
		return int32(5)
	case uint16(54):
		return int32(6)
	case uint16(55):
		return int32(7)
	case uint16(56):
		return int32(8)
	case uint16(57):
		return int32(9)
	case uint16(97), uint16(65):
		return int32(10)
	case uint16(98), uint16(66):
		return int32(11)
	case uint16(99), uint16(67):
		return int32(12)
	case uint16(100), uint16(68):
		return int32(13)
	case uint16(101), uint16(69):
		return int32(14)
	case uint16(102), uint16(70):
		return int32(15)
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "unexpected: "+j2gort.StrValChar(c))))
	}
	panic("unreachable switch expression")
}

// Java static method: IRFactory.transformBooleanTokenType(com.google.javascript.jscomp.parsing.parser.TokenType token)
func IRFactory_TransformBooleanTokenType(token pkg_parser.TokenType) pkg_rhino.Token {
	clinit_IRFactory()
	_ = token
	switch (token).Ordinal() {
	case int32(102) /* TokenType.TRUE */ :
		return pkg_rhino.StaticGet_Token_TRUE()
	case int32(103) /* TokenType.FALSE */ :
		return pkg_rhino.StaticGet_Token_FALSE()
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.StrVal(token))))
	}
	panic("unreachable switch expression")
}

// Java static method: IRFactory.transformUpdateTokenType(com.google.javascript.jscomp.parsing.parser.TokenType token)
func IRFactory_TransformUpdateTokenType(token pkg_parser.TokenType) pkg_rhino.Token {
	clinit_IRFactory()
	_ = token
	switch (token).Ordinal() {
	case int32(68) /* TokenType.PLUS_PLUS */ :
		return pkg_rhino.StaticGet_Token_INC()
	case int32(69) /* TokenType.MINUS_MINUS */ :
		return pkg_rhino.StaticGet_Token_DEC()
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.StrVal(token))))
	}
	panic("unreachable switch expression")
}

// Java static method: IRFactory.transformUnaryTokenType(com.google.javascript.jscomp.parsing.parser.TokenType token)
func IRFactory_TransformUnaryTokenType(token pkg_parser.TokenType) pkg_rhino.Token {
	clinit_IRFactory()
	_ = token
	switch (token).Ordinal() {
	case int32(76) /* TokenType.BANG */ :
		return pkg_rhino.StaticGet_Token_NOT()
	case int32(77) /* TokenType.TILDE */ :
		return pkg_rhino.StaticGet_Token_BITNOT()
	case int32(63) /* TokenType.PLUS */ :
		return pkg_rhino.StaticGet_Token_POS()
	case int32(64) /* TokenType.MINUS */ :
		return pkg_rhino.StaticGet_Token_NEG()
	case int32(9) /* TokenType.DELETE */ :
		return pkg_rhino.StaticGet_Token_DELPROP()
	case int32(24) /* TokenType.TYPEOF */ :
		return pkg_rhino.StaticGet_Token_TYPEOF()
	case int32(26) /* TokenType.VOID */ :
		return pkg_rhino.StaticGet_Token_VOID()
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.StrVal(token))))
	}
	panic("unreachable switch expression")
}

// Java static method: IRFactory.transformBinaryTokenType(com.google.javascript.jscomp.parsing.parser.TokenType token)
func IRFactory_TransformBinaryTokenType(token pkg_parser.TokenType) pkg_rhino.Token {
	clinit_IRFactory()
	_ = token
	switch (token).Ordinal() {
	case int32(74) /* TokenType.BAR */ :
		return pkg_rhino.StaticGet_Token_BITOR()
	case int32(75) /* TokenType.CARET */ :
		return pkg_rhino.StaticGet_Token_BITXOR()
	case int32(73) /* TokenType.AMPERSAND */ :
		return pkg_rhino.StaticGet_Token_BITAND()
	case int32(59) /* TokenType.EQUAL_EQUAL */ :
		return pkg_rhino.StaticGet_Token_EQ()
	case int32(60) /* TokenType.NOT_EQUAL */ :
		return pkg_rhino.StaticGet_Token_NE()
	case int32(54) /* TokenType.OPEN_ANGLE */ :
		return pkg_rhino.StaticGet_Token_LT()
	case int32(56) /* TokenType.LESS_EQUAL */ :
		return pkg_rhino.StaticGet_Token_LE()
	case int32(55) /* TokenType.CLOSE_ANGLE */ :
		return pkg_rhino.StaticGet_Token_GT()
	case int32(57) /* TokenType.GREATER_EQUAL */ :
		return pkg_rhino.StaticGet_Token_GE()
	case int32(70) /* TokenType.LEFT_SHIFT */ :
		return pkg_rhino.StaticGet_Token_LSH()
	case int32(71) /* TokenType.RIGHT_SHIFT */ :
		return pkg_rhino.StaticGet_Token_RSH()
	case int32(72) /* TokenType.UNSIGNED_RIGHT_SHIFT */ :
		return pkg_rhino.StaticGet_Token_URSH()
	case int32(63) /* TokenType.PLUS */ :
		return pkg_rhino.StaticGet_Token_ADD()
	case int32(64) /* TokenType.MINUS */ :
		return pkg_rhino.StaticGet_Token_SUB()
	case int32(65) /* TokenType.STAR */ :
		return pkg_rhino.StaticGet_Token_MUL()
	case int32(96) /* TokenType.SLASH */ :
		return pkg_rhino.StaticGet_Token_DIV()
	case int32(67) /* TokenType.PERCENT */ :
		return pkg_rhino.StaticGet_Token_MOD()
	case int32(66) /* TokenType.STAR_STAR */ :
		return pkg_rhino.StaticGet_Token_EXPONENT()
	case int32(61) /* TokenType.EQUAL_EQUAL_EQUAL */ :
		return pkg_rhino.StaticGet_Token_SHEQ()
	case int32(62) /* TokenType.NOT_EQUAL_EQUAL */ :
		return pkg_rhino.StaticGet_Token_SHNE()
	case int32(16) /* TokenType.IN */ :
		return pkg_rhino.StaticGet_Token_IN()
	case int32(17) /* TokenType.INSTANCEOF */ :
		return pkg_rhino.StaticGet_Token_INSTANCEOF()
	case int32(53) /* TokenType.COMMA */ :
		return pkg_rhino.StaticGet_Token_COMMA()
	case int32(84) /* TokenType.EQUAL */ :
		return pkg_rhino.StaticGet_Token_ASSIGN()
	case int32(94) /* TokenType.BAR_EQUAL */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_BITOR()
	case int32(95) /* TokenType.CARET_EQUAL */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_BITXOR()
	case int32(93) /* TokenType.AMPERSAND_EQUAL */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_BITAND()
	case int32(90) /* TokenType.LEFT_SHIFT_EQUAL */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_LSH()
	case int32(91) /* TokenType.RIGHT_SHIFT_EQUAL */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_RSH()
	case int32(92) /* TokenType.UNSIGNED_RIGHT_SHIFT_EQUAL */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_URSH()
	case int32(85) /* TokenType.PLUS_EQUAL */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_ADD()
	case int32(86) /* TokenType.MINUS_EQUAL */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_SUB()
	case int32(87) /* TokenType.STAR_EQUAL */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_MUL()
	case int32(88) /* TokenType.STAR_STAR_EQUAL */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_EXPONENT()
	case int32(97) /* TokenType.SLASH_EQUAL */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_DIV()
	case int32(89) /* TokenType.PERCENT_EQUAL */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_MOD()
	case int32(79) /* TokenType.OR */ :
		return pkg_rhino.StaticGet_Token_OR()
	case int32(78) /* TokenType.AND */ :
		return pkg_rhino.StaticGet_Token_AND()
	case int32(81) /* TokenType.QUESTION_QUESTION */ :
		return pkg_rhino.StaticGet_Token_COALESCE()
	case int32(99) /* TokenType.OR_EQUAL */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_OR()
	case int32(98) /* TokenType.AND_EQUAL */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_AND()
	case int32(100) /* TokenType.QUESTION_QUESTION_EQUAL */ :
		return pkg_rhino.StaticGet_Token_ASSIGN_COALESCE()
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.StrVal(token))))
	}
	panic("unreachable switch expression")
}

// Java method: IRFactory.newNode(com.google.javascript.rhino.Token type)
func (this_ *IRFactory) NewNode_1(type_ pkg_rhino.Token) pkg_rhino.Node {
	return pkg_rhino.Node(this_.nodeAlloc_.AllocNode1Cached(uint8(type_) - 1))

}

// Java method: IRFactory.newNode(com.google.javascript.rhino.Token type, com.google.javascript.rhino.Node child1)
func (this_ *IRFactory) NewNode_2(type_ pkg_rhino.Token, child1 pkg_rhino.Node) pkg_rhino.Node {
	return pkg_rhino.Node(this_.nodeAlloc_.AllocNode2(uint8(type_)-1, nodestore.Node(child1), nodestore.Node(this_.templateNode)))

}

// Java method: IRFactory.newNode(com.google.javascript.rhino.Token type, com.google.javascript.rhino.Node child1, com.google.javascript.rhino.Node child2)
func (this_ *IRFactory) NewNode_3(type_ pkg_rhino.Token, child1 pkg_rhino.Node, child2 pkg_rhino.Node) pkg_rhino.Node {
	return pkg_rhino.Node(this_.nodeAlloc_.AllocNode3(uint8(type_)-1, nodestore.Node(child1), nodestore.Node(child2), nodestore.Node(this_.templateNode)))

}

// Java method: IRFactory.newNode(com.google.javascript.rhino.Token type, com.google.javascript.rhino.Node child1, com.google.javascript.rhino.Node child2, com.google.javascript.rhino.Node child3)
func (this_ *IRFactory) NewNode_4(type_ pkg_rhino.Token, child1 pkg_rhino.Node, child2 pkg_rhino.Node, child3 pkg_rhino.Node) pkg_rhino.Node {
	return pkg_rhino.Node(this_.nodeAlloc_.AllocNode4(uint8(type_)-1, nodestore.Node(child1), nodestore.Node(child2), nodestore.Node(child3), nodestore.Node(this_.templateNode)))

}

// Java method: IRFactory.newStringNode(java.lang.String value)
func (this_ *IRFactory) NewStringNode_1(value string) pkg_rhino.Node {
	if value == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	return pkg_rhino.Node(this_.nodeAlloc_.AllocStringNodeCached(uint8(pkg_rhino.StaticGet_Token_STRINGLIT())-1, this_.nodeAlloc_.Intern(value)))

}

// Java method: IRFactory.newStringNode(com.google.javascript.rhino.Token type, java.lang.String value)
func (this_ *IRFactory) NewStringNode_2(type_ pkg_rhino.Token, value string) pkg_rhino.Node {
	if value == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	return pkg_rhino.Node(this_.nodeAlloc_.AllocStringNodeCached(uint8(type_)-1, this_.nodeAlloc_.Intern(value)))

}

// Java method: IRFactory.maybeAddNonJsDocComment(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.parsing.parser.util.SourcePosition start)
func (this_ *IRFactory) MaybeAddNonJsDocComment(node pkg_rhino.Node, start pkg_parser.SourcePosition) pkg_rhino.Node {
	_ = node
	_ = start
	var comment *pkg_rhino.NonJSDocComment = this_.priv_IRFactory_parseNonJSDocCommentAt_6838(start, false)
	_ = comment
	if comment != nil {
		node.SetNonJSDocComment(comment)
	}
	return node
}

// Java method: IRFactory.newTemplateLitStringNode(java.lang.@org.jspecify.annotations.Nullable String cooked, java.lang.String raw)
func (this_ *IRFactory) NewTemplateLitStringNode(cooked string, raw string) pkg_rhino.Node {
	_ = cooked
	_ = raw
	return pkg_rhino.Node_NewTemplateLitString_2(cooked, raw).ClonePropsFrom(this_.templateNode)
}

// Java method: IRFactory.newNumberNode(java.lang.Double value)
func (this_ *IRFactory) NewNumberNode(value *j2gort.BoxDouble) pkg_rhino.Node {
	d_ := j2gort.UnboxDouble(value)
	if d_ != d_ || d_ < 0 || (d_ == 0 && 1/d_ < 0) {
		return pkg_rhino.IR_Number(d_).ClonePropsFrom(this_.templateNode)
	}
	return pkg_rhino.Node(this_.nodeAlloc_.AllocNumberNode(uint8(pkg_rhino.StaticGet_Token_NUMBER())-1, d_, nodestore.Node(this_.templateNode)))

}

// Java method: IRFactory.newBigIntNode(java.math.BigInteger value)
func (this_ *IRFactory) NewBigIntNode(value *j2gort.Ext_java_math_BigInteger) pkg_rhino.Node {
	_ = value
	return pkg_rhino.Node_NewBigInt(value).ClonePropsFrom(this_.templateNode)
}

// Java method: IRFactory.error(com.google.javascript.rhino.Node rhinoNode, java.lang.String message)
func (this_ *IRFactory) priv_IRFactory_error_6876(rhinoNode pkg_rhino.Node, message string) {
	_ = rhinoNode
	_ = message
	this_.errorReporter.M_Error(message, rhinoNode.GetSourceFileName(), rhinoNode.GetLineno(), rhinoNode.GetCharno())
}

// Java record: com.google.javascript.jscomp.parsing.IRFactory.LineAndColumn
type IRFactory_LineAndColumn struct {
	lineNo int32 `gson:"lineNo"`
	column int32 `gson:"column"`
}

func (this_ *IRFactory_LineAndColumn) BindSelf_(s any) {
	_ = s
}
func (this_ *IRFactory_LineAndColumn) IsJavaNil_() bool { return this_ == nil }
func (*IRFactory_LineAndColumn) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.IRFactory$LineAndColumn", "LineAndColumn"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.IRFactory.LineAndColumn", "com.google.javascript.jscomp.parsing.IRFactory$LineAndColumn", func(v_ any) bool {
		return j2gort.InstanceOf[*IRFactory_LineAndColumn](v_)
	})
}
func (*IRFactory_LineAndColumn) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *IRFactory_LineAndColumn) initFields_IRFactory_LineAndColumn() {
	_ = this_
}

// Java constructor: LineAndColumn(int lineNo, int column)
func NewIRFactory_LineAndColumn(lineNo int32, column int32) *IRFactory_LineAndColumn {
	this_ := &IRFactory_LineAndColumn{}
	InitIRFactory_LineAndColumn(this_, lineNo, column)
	return this_
}

func InitIRFactory_LineAndColumn(this_ *IRFactory_LineAndColumn, lineNo int32, column int32) {
	_ = lineNo
	_ = column
	this_.initFields_IRFactory_LineAndColumn()
	this_.lineNo = lineNo
	this_.column = column
}

// Java static method: LineAndColumn.fromSourcePosition(com.google.javascript.jscomp.parsing.parser.util.SourcePosition pos)
func IRFactory_LineAndColumn_FromSourcePosition(pos pkg_parser.SourcePosition) *IRFactory_LineAndColumn {
	_ = pos
	return NewIRFactory_LineAndColumn(pos.FieldGet_line(), pos.FieldGet_column())
}

// Java method: LineAndColumn.compareTo(com.google.javascript.jscomp.parsing.IRFactory.LineAndColumn o)
func (this_ *IRFactory_LineAndColumn) CompareTo(arg_43 any) int32 {
	o := j2gort.CastPtr[IRFactory_LineAndColumn](arg_43)
	_ = o
	return j2gort.Ext_java_util_Comparator_ComparingInt(lambdaSingleton_ir_factory_4).ThenComparingInt(lambdaSingleton_ir_factory_3).Compare(this_, j2gort.PtrToAny(o))
}

func (this_ *IRFactory_LineAndColumn) LineNo() int32 { return this_.lineNo }
func (this_ *IRFactory_LineAndColumn) Column() int32 { return this_.column }
func (this_ *IRFactory_LineAndColumn) Equals(other_ any) bool {
	o_, ok_ := other_.(*IRFactory_LineAndColumn)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.lineNo != o_.lineNo {
		return false
	}
	if this_.column != o_.column {
		return false
	}
	return true
}
func (this_ *IRFactory_LineAndColumn) HashCode() int32 {
	var h_ int32 = 0
	h_ = 31*h_ + (this_.lineNo)
	h_ = 31*h_ + (this_.column)
	return h_
}
func (this_ *IRFactory_LineAndColumn) String() string {
	_ = fmt.Sprint
	return "LineAndColumn[" + "lineNo=" + j2gort.StrVal(j2gort.ValueOfInt(this_.lineNo)) + ", " + "column=" + j2gort.StrVal(j2gort.ValueOfInt(this_.column)) + "]"
}

// Java class: com.google.javascript.jscomp.parsing.IRFactory.CommentTracker
type IRFactory_CommentTracker struct {
	source        []*pkg_parser.Comment `gson:"source"`
	index         int32                 `gson:"index"`
	previousIndex int32                 `gson:"previousIndex"`
	idHash_       atomic.Int32
	jsdocOnly     bool `gson:"jsdocOnly"`
}

func (this_ *IRFactory_CommentTracker) BindSelf_(s any) {
	_ = s
}
func (this_ *IRFactory_CommentTracker) IsJavaNil_() bool { return this_ == nil }
func (*IRFactory_CommentTracker) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.IRFactory$CommentTracker", "CommentTracker"
}
func (this_ *IRFactory_CommentTracker) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*IRFactory_CommentTracker) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.IRFactory.CommentTracker", "com.google.javascript.jscomp.parsing.IRFactory$CommentTracker", func(v_ any) bool {
		return j2gort.InstanceOf[*IRFactory_CommentTracker](v_)
	})
}

func (this_ *IRFactory_CommentTracker) initFields_IRFactory_CommentTracker() {
	_ = this_
	this_.index = int32(-1)
	this_.previousIndex = int32(-1)
}

// Java constructor: CommentTracker(com.google.common.collect.ImmutableList<com.google.javascript.jscomp.parsing.parser.trees.Comment> source, boolean jsdocOnly)
func NewIRFactory_CommentTracker(source j2gort.Ext_com_google_common_collect_ImmutableList, jsdocOnly bool) *IRFactory_CommentTracker {
	this_ := &IRFactory_CommentTracker{}
	InitIRFactory_CommentTracker(this_, source, jsdocOnly)
	return this_
}

func InitIRFactory_CommentTracker(this_ *IRFactory_CommentTracker, source j2gort.Ext_com_google_common_collect_ImmutableList, jsdocOnly bool) {
	_ = source
	_ = jsdocOnly
	this_.initFields_IRFactory_CommentTracker()
	this_.source = j2gort.ImmListToSliceNullable[*pkg_parser.Comment](source)
	this_.jsdocOnly = jsdocOnly
	this_.Advance()
}

// Java method: CommentTracker.current()
func (this_ *IRFactory_CommentTracker) Current() *pkg_parser.Comment {
	if this_.index >= int32(len(this_.source)) {
		return nil
	} else {
		return this_.source[this_.index]
	}
}

// Java method: CommentTracker.advance()
func (this_ *IRFactory_CommentTracker) Advance() {
	this_.previousIndex = this_.index
	for {
		this_.index = this_.index + 1
		var c *pkg_parser.Comment = this_.Current()
		_ = c
		if (c == nil) || ((c.FieldGet_type_() == pkg_parser.StaticGet_Comment_Type_JSDOC()) == this_.jsdocOnly) {
			break
		}
	}
}

// Java method: CommentTracker.backtrack()
func (this_ *IRFactory_CommentTracker) Backtrack() {
	this_.index = this_.previousIndex
}

// Java method: CommentTracker.hasPendingCommentBefore(com.google.javascript.jscomp.parsing.parser.util.SourcePosition pos)
func (this_ *IRFactory_CommentTracker) HasPendingCommentBefore(pos pkg_parser.SourcePosition) bool {
	_ = pos
	var c *pkg_parser.Comment = this_.Current()
	_ = c
	return (((c != nil) && (c.FieldGet_location().FieldGet_end().FieldGet_line() <= pos.FieldGet_line())) && (c.FieldGet_location().FieldGet_end().FieldGet_offset() <= pos.FieldGet_offset()))
}

// Java class: com.google.javascript.jscomp.parsing.IRFactory.TransformDispatcher
type IRFactory_TransformDispatcher struct {
	outer_                *IRFactory
	classScope            *IRFactory_TransformDispatcher_ScopeTracker `gson:"classScope"`
	privateIdLhsOfInScope *IRFactory_TransformDispatcher_ScopeTracker `gson:"privateIdLhsOfInScope"`
	idHash_               atomic.Int32
}

func (this_ *IRFactory_TransformDispatcher) BindSelf_(s any) {
	_ = s
}
func (this_ *IRFactory_TransformDispatcher) IsJavaNil_() bool { return this_ == nil }
func (*IRFactory_TransformDispatcher) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.IRFactory$TransformDispatcher", "TransformDispatcher"
}
func (this_ *IRFactory_TransformDispatcher) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*IRFactory_TransformDispatcher) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.IRFactory.TransformDispatcher", "com.google.javascript.jscomp.parsing.IRFactory$TransformDispatcher", func(v_ any) bool {
		return j2gort.InstanceOf[*IRFactory_TransformDispatcher](v_)
	})
}
func (this_ *IRFactory_TransformDispatcher) SetOuter_(o *IRFactory) { this_.outer_ = o }
func (this_ *IRFactory_TransformDispatcher) GetOuter_() *IRFactory  { return this_.outer_ }

func (this_ *IRFactory_TransformDispatcher) initFields_IRFactory_TransformDispatcher() {
	_ = this_
	this_.classScope = NewIRFactory_TransformDispatcher_ScopeTracker()
	this_.privateIdLhsOfInScope = NewIRFactory_TransformDispatcher_ScopeTracker()
}

// Java constructor: TransformDispatcher()
func NewIRFactory_TransformDispatcher(outer_ *IRFactory) *IRFactory_TransformDispatcher {
	this_ := &IRFactory_TransformDispatcher{}
	this_.outer_ = outer_
	InitIRFactory_TransformDispatcher(this_)
	return this_
}

func InitIRFactory_TransformDispatcher(this_ *IRFactory_TransformDispatcher) {
	this_.initFields_IRFactory_TransformDispatcher()
}

// Java method: TransformDispatcher.newStringNodeFromIdentifier(com.google.javascript.rhino.Token type, com.google.javascript.jscomp.parsing.IRFactory.TransformDispatcher.IdentifierType identifierType, com.google.javascript.jscomp.parsing.parser.Token identifierToken)
func (this_ *IRFactory_TransformDispatcher) NewStringNodeFromIdentifier(type_ pkg_rhino.Token, identifierType IRFactory_TransformDispatcher_IdentifierType, identifierToken pkg_parser.Token) pkg_rhino.Node {
	_ = type_
	_ = identifierType
	_ = identifierToken
	if !(this_.classScope.InScope()) {
		identifierType = StaticGet_IRFactory_TransformDispatcher_IdentifierType_STANDARD()
	}
	var value string = j2gort.NullString
	_ = value
	if identifierToken.IsPrivateIdentifier() {
		if identifierType == StaticGet_IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE() {
			this_.outer_.MaybeWarnForFeature_2_Token_FeatureSet_Feature(identifierToken, pkg_parser.StaticGet_FeatureSet_Feature_PRIVATE_ELEMENTS())
		} else {
			this_.outer_.errorReporter.M_Error("Private identifiers may not be used in this context", this_.outer_.sourceName, IRFactory_Lineno_1_Token(identifierToken), IRFactory_Charno_1_Token(identifierToken))
		}
		value = identifierToken.GetMaybePrivateValue()
	} else {
		value = identifierToken.GetValue()
	}
	var node pkg_rhino.Node = this_.outer_.NewStringNode_2(type_, value)
	_ = node
	if (identifierType == StaticGet_IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE()) && identifierToken.IsPrivateIdentifier() {
		node.SetPrivateIdentifier()
	}
	return node
}

// Java method: TransformDispatcher.processObjectLitKey(com.google.javascript.jscomp.parsing.parser.Token input, com.google.javascript.rhino.Token output, com.google.javascript.jscomp.parsing.IRFactory.TransformDispatcher.IdentifierType identifierType)
func (this_ *IRFactory_TransformDispatcher) priv_TransformDispatcher_processObjectLitKey_6884(input pkg_parser.Token, output pkg_rhino.Token, identifierType IRFactory_TransformDispatcher_IdentifierType) pkg_rhino.Node {
	_ = input
	_ = output
	_ = identifierType
	if (input).IsJavaNil_() {
		return this_.priv_TransformDispatcher_createMissingExpressionNode_6984()
	}
	if input.FieldGet_type_() == pkg_parser.StaticGet_TokenType_IDENTIFIER() {
		return this_.ProcessName_3(input, output, identifierType)
	}
	var literal pkg_parser.Token = input
	_ = literal
	var jsDocInfo *pkg_rhino.JSDocInfo = this_.outer_.ParseJSDocInfoOnToken(literal)
	_ = jsDocInfo
	var comment *pkg_rhino.NonJSDocComment = this_.outer_.priv_IRFactory_parseNonJSDocCommentAt_6838(literal.GetStart(), true)
	_ = comment
	var node pkg_rhino.Node = pkg_rhino.Node(0)
	_ = node
	switch (input.FieldGet_type_()).Ordinal() {
	case int32(104) /* TokenType.NUMBER */ :
		node = this_.outer_.NewStringNode_2(output, pkg_dtoa.DToA_NumberToString(this_.outer_.NormalizeNumber(literal)))
	case int32(106) /* TokenType.BIGINT */ :
		node = this_.outer_.NewStringNode_2(output, this_.outer_.NormalizeBigInt(literal).String())
	default:
		node = this_.outer_.NewStringNode_2(output, this_.outer_.NormalizeString(literal, false))
	}
	if jsDocInfo != nil {
		node.SetJSDocInfo(jsDocInfo)
	}
	if comment != nil {
		node.SetNonJSDocComment(comment)
	}
	this_.outer_.SetSourceInfo_2_Node_Token(node, literal)
	node.PutBooleanProp(pkg_rhino.StaticGet_Node_QUOTED_PROP(), true)
	return node
}

// Java method: TransformDispatcher.processArrayLiteral(com.google.javascript.jscomp.parsing.parser.trees.ArrayLiteralExpressionTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessArrayLiteral(tree *pkg_parser.ArrayLiteralExpressionTree) pkg_rhino.Node {
	_ = tree
	var node pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_ARRAYLIT())
	_ = node
	node.SetTrailingComma(tree.FieldGet_hasTrailingComma())
	{
		cur_44 := j2gort.NewIterCursor(j2gort.CheckNotNull(tree.FieldGet_elements()))
		for cur_44.HasNext() {
			rawElem_45 := cur_44.Next()
			child := j2gort.AsIface[pkg_parser.ParseTree](rawElem_45)
			_ = child
			var c pkg_rhino.Node = this_.outer_.Transform(child)
			_ = c
			node.AddChildToBack(c)
		}
	}
	return node
}

// Java method: TransformDispatcher.processArrayPattern(com.google.javascript.jscomp.parsing.parser.trees.ArrayPatternTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessArrayPattern(tree *pkg_parser.ArrayPatternTree) pkg_rhino.Node {
	_ = tree
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(tree), pkg_parser.StaticGet_FeatureSet_Feature_ARRAY_DESTRUCTURING())
	var node pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_ARRAY_PATTERN())
	_ = node
	{
		cur_46 := j2gort.NewIterCursor(j2gort.CheckNotNull(tree.FieldGet_elements()))
		for cur_46.HasNext() {
			rawElem_47 := cur_46.Next()
			child := j2gort.AsIface[pkg_parser.ParseTree](rawElem_47)
			_ = child
			var elementNode pkg_rhino.Node = pkg_rhino.Node(0)
			_ = elementNode
			switch (child.GetParseTreeBase().FieldGet_type_()).Ordinal() {
			case int32(71) /* ParseTreeType.DEFAULT_PARAMETER */ :
				elementNode = this_.ProcessDefaultParameter(child.GetParseTreeBase().AsDefaultParameter())
			case int32(61) /* ParseTreeType.ITER_REST */ :
				this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(child, pkg_parser.StaticGet_FeatureSet_Feature_ARRAY_PATTERN_REST())
				elementNode = this_.outer_.TransformNodeWithInlineComments(child)
			default:
				elementNode = this_.outer_.TransformNodeWithInlineComments(child)
			}
			node.AddChildToBack(elementNode)
		}
	}
	return node
}

// Java method: TransformDispatcher.processObjectPattern(com.google.javascript.jscomp.parsing.parser.trees.ObjectPatternTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessObjectPattern(tree *pkg_parser.ObjectPatternTree) pkg_rhino.Node {
	_ = tree
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(tree), pkg_parser.StaticGet_FeatureSet_Feature_OBJECT_DESTRUCTURING())
	var node pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_OBJECT_PATTERN())
	_ = node
	{
		cur_48 := j2gort.NewIterCursor(j2gort.CheckNotNull(tree.FieldGet_fields()))
		for cur_48.HasNext() {
			rawElem_49 := cur_48.Next()
			child := j2gort.AsIface[pkg_parser.ParseTree](rawElem_49)
			_ = child
			var childNode pkg_rhino.Node = this_.priv_TransformDispatcher_processObjectPatternElement_6888(child)
			_ = childNode
			node.AddChildToBack(childNode)
		}
	}
	return node
}

// Java method: TransformDispatcher.processObjectPatternElement(com.google.javascript.jscomp.parsing.parser.trees.ParseTree child)
func (this_ *IRFactory_TransformDispatcher) priv_TransformDispatcher_processObjectPatternElement_6888(child pkg_parser.ParseTree) pkg_rhino.Node {
	_ = child
	base_child_ := child.GetParseTreeBase()
	switch (base_child_.FieldGet_type_()).Ordinal() {
	case int32(71) /* ParseTreeType.DEFAULT_PARAMETER */ :
		return this_.priv_TransformDispatcher_processObjectPatternShorthandWithDefault_6889(base_child_.AsDefaultParameter())
	case int32(32) /* ParseTreeType.PROPERTY_NAME_ASSIGNMENT */ :
		return this_.priv_TransformDispatcher_processObjectPatternPropertyNameAssignment_6890(base_child_.AsPropertyNameAssignment())
	case int32(33) /* ParseTreeType.COMPUTED_PROPERTY_DEFINITION */ :
		var computedPropertyDefinition *pkg_parser.ComputedPropertyDefinitionTree = base_child_.AsComputedPropertyDefinition()
		_ = computedPropertyDefinition
		return this_.priv_TransformDispatcher_processObjectPatternComputedPropertyDefinition_6892(computedPropertyDefinition)
	case int32(62) /* ParseTreeType.OBJECT_REST */ :
		this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(child, pkg_parser.StaticGet_FeatureSet_Feature_OBJECT_PATTERN_REST())
		var target pkg_rhino.Node = this_.outer_.TransformNodeWithInlineComments(base_child_.AsObjectRest().FieldGet_assignmentTarget())
		_ = target
		var rest pkg_rhino.Node = this_.outer_.NewNode_2(pkg_rhino.StaticGet_Token_OBJECT_REST(), target)
		_ = rest
		this_.outer_.SetSourceInfo_2_Node_ParseTree(rest, child)
		return rest
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected object pattern element: "+j2gort.StrVal(child))))
	}
	panic("unreachable switch expression")
}

// Java method: TransformDispatcher.processObjectPatternShorthandWithDefault(com.google.javascript.jscomp.parsing.parser.trees.DefaultParameterTree defaultParameter)
func (this_ *IRFactory_TransformDispatcher) priv_TransformDispatcher_processObjectPatternShorthandWithDefault_6889(defaultParameter *pkg_parser.DefaultParameterTree) pkg_rhino.Node {
	_ = defaultParameter
	var defaultValueNode pkg_rhino.Node = this_.ProcessDefaultParameter(defaultParameter)
	_ = defaultValueNode
	var nameNode pkg_rhino.Node = defaultValueNode.GetFirstChild()
	_ = nameNode
	var stringKeyNode pkg_rhino.Node = this_.outer_.MaybeAddNonJsDocComment(this_.outer_.NewStringNode_2(pkg_rhino.StaticGet_Token_STRING_KEY(), nameNode.GetString()), defaultParameter.GetStart())
	_ = stringKeyNode
	this_.outer_.SetSourceInfo_2_Node_Node(stringKeyNode, nameNode)
	stringKeyNode.SetShorthandProperty(true)
	stringKeyNode.AddChildToBack(defaultValueNode)
	return stringKeyNode
}

// Java method: TransformDispatcher.processObjectPatternPropertyNameAssignment(com.google.javascript.jscomp.parsing.parser.trees.PropertyNameAssignmentTree propertyNameAssignment)
func (this_ *IRFactory_TransformDispatcher) priv_TransformDispatcher_processObjectPatternPropertyNameAssignment_6890(propertyNameAssignment *pkg_parser.PropertyNameAssignmentTree) pkg_rhino.Node {
	_ = propertyNameAssignment
	var key pkg_rhino.Node = this_.priv_TransformDispatcher_processObjectLitKey_6884(propertyNameAssignment.FieldGet_name(), pkg_rhino.StaticGet_Token_STRING_KEY(), StaticGet_IRFactory_TransformDispatcher_IdentifierType_STANDARD())
	_ = key
	var targetTree pkg_parser.ParseTree = propertyNameAssignment.FieldGet_value()
	_ = targetTree
	var valueNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = valueNode
	if targetTree == nil || targetTree.IsJavaNil_() {
		valueNode = this_.ProcessNameWithInlineComments_2_Token_IRFactory_TransformDispatcher_IdentifierType(propertyNameAssignment.FieldGet_name(), StaticGet_IRFactory_TransformDispatcher_IdentifierType_STANDARD())
		key.SetShorthandProperty(true)
	} else {
		valueNode = this_.priv_TransformDispatcher_processDestructuringElementTarget_6891(targetTree)
	}
	key.AddChildToFront(valueNode)
	return key
}

// Java method: TransformDispatcher.processDestructuringElementTarget(com.google.javascript.jscomp.parsing.parser.trees.ParseTree targetTree)
func (this_ *IRFactory_TransformDispatcher) priv_TransformDispatcher_processDestructuringElementTarget_6891(targetTree pkg_parser.ParseTree) pkg_rhino.Node {
	_ = targetTree
	var valueNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = valueNode
	base_targetTree_ := targetTree.GetParseTreeBase()
	if base_targetTree_.FieldGet_type_() == pkg_parser.StaticGet_ParseTreeType_DEFAULT_PARAMETER() {
		valueNode = this_.ProcessDefaultParameter(base_targetTree_.AsDefaultParameter())
	} else if base_targetTree_.FieldGet_type_() == pkg_parser.StaticGet_ParseTreeType_IDENTIFIER_EXPRESSION() {
		valueNode = this_.priv_TransformDispatcher_processNameWithInlineComments_6941(base_targetTree_.AsIdentifierExpression(), StaticGet_IRFactory_TransformDispatcher_IdentifierType_STANDARD())
	} else {
		valueNode = this_.outer_.TransformNodeWithInlineComments(targetTree)
	}
	return valueNode
}

// Java method: TransformDispatcher.processObjectPatternComputedPropertyDefinition(com.google.javascript.jscomp.parsing.parser.trees.ComputedPropertyDefinitionTree computedPropertyDefinition)
func (this_ *IRFactory_TransformDispatcher) priv_TransformDispatcher_processObjectPatternComputedPropertyDefinition_6892(computedPropertyDefinition *pkg_parser.ComputedPropertyDefinitionTree) pkg_rhino.Node {
	_ = computedPropertyDefinition
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(computedPropertyDefinition), pkg_parser.StaticGet_FeatureSet_Feature_COMPUTED_PROPERTIES())
	var expressionNode pkg_rhino.Node = this_.outer_.Transform(computedPropertyDefinition.FieldGet_property())
	_ = expressionNode
	var valueTree pkg_parser.ParseTree = computedPropertyDefinition.FieldGet_value()
	_ = valueTree
	var valueNode pkg_rhino.Node = this_.priv_TransformDispatcher_processDestructuringElementTarget_6891(valueTree)
	_ = valueNode
	var computedPropertyNode pkg_rhino.Node = this_.outer_.NewNode_3(pkg_rhino.StaticGet_Token_COMPUTED_PROP(), expressionNode, valueNode)
	_ = computedPropertyNode
	this_.outer_.SetSourceInfo_2_Node_ParseTree(computedPropertyNode, pkg_parser.ToIface_ParseTree(computedPropertyDefinition))
	return computedPropertyNode
}

// Java method: TransformDispatcher.processAstRoot(com.google.javascript.jscomp.parsing.parser.trees.ProgramTree rootNode)
func (this_ *IRFactory_TransformDispatcher) ProcessAstRoot(rootNode *pkg_parser.ProgramTree) pkg_rhino.Node {
	_ = rootNode
	var scriptNode pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_SCRIPT())
	_ = scriptNode
	{
		cur_50 := j2gort.NewIterCursor(j2gort.CheckNotNull(rootNode.FieldGet_sourceElements()))
		for cur_50.HasNext() {
			rawElem_51 := cur_50.Next()
			child := j2gort.AsIface[pkg_parser.ParseTree](rawElem_51)
			_ = child
			scriptNode.AddChildToBack(this_.outer_.Transform(child))
		}
	}
	this_.priv_TransformDispatcher_parseDirectives_6895(scriptNode)
	var isGoogModule bool = this_.priv_TransformDispatcher_isGoogModuleFile_6894(scriptNode)
	_ = isGoogModule
	if isGoogModule || this_.outer_.features.Has(pkg_parser.StaticGet_FeatureSet_Feature_MODULES()) {
		var moduleNode pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_MODULE_BODY())
		_ = moduleNode
		this_.outer_.SetSourceInfo_2_Node_ParseTree(moduleNode, pkg_parser.ToIface_ParseTree(rootNode))
		moduleNode.AddChildrenToBack(scriptNode.RemoveChildren())
		scriptNode.AddChildToBack(moduleNode)
		if isGoogModule {
			scriptNode.PutBooleanProp(pkg_rhino.StaticGet_Node_GOOG_MODULE(), true)
		} else {
			scriptNode.PutBooleanProp(pkg_rhino.StaticGet_Node_ES6_MODULE(), true)
		}
	}
	return scriptNode
}

// Java method: TransformDispatcher.isGoogModuleFile(com.google.javascript.rhino.Node scriptNode)
func (this_ *IRFactory_TransformDispatcher) priv_TransformDispatcher_isGoogModuleFile_6894(scriptNode pkg_rhino.Node) bool {
	_ = scriptNode
	if !(scriptNode.IsScript()) {
		j2gort.CheckArgument(false)
	}
	if !(scriptNode.HasChildren()) {
		return false
	}
	var exprResult pkg_rhino.Node = scriptNode.GetFirstChild()
	_ = exprResult
	if !(exprResult.IsExprResult()) {
		return false
	}
	var call pkg_rhino.Node = exprResult.GetFirstChild()
	_ = call
	if !(call.IsCall()) {
		return false
	}
	return ((StaticGet_IRFactory_GOOG_MODULE().Matches(call.GetFirstChild()) && call.HasTwoChildren()) && call.GetSecondChild().IsStringLit())
}

// Java method: TransformDispatcher.parseDirectives(com.google.javascript.rhino.Node node)
func (this_ *IRFactory_TransformDispatcher) priv_TransformDispatcher_parseDirectives_6895(node pkg_rhino.Node) {
	_ = node
	var useStrict bool = false
	_ = useStrict
	for {
		var statement pkg_rhino.Node = node.GetFirstChild()
		_ = statement
		if (statement == 0) || !(statement.IsExprResult()) {
			break
		}
		var directive pkg_rhino.Node = statement.GetFirstChild()
		_ = directive
		if !(directive.IsStringLit()) || !(j2gort.CheckNotNullString(directive.GetString()) == "use strict") {
			break
		}
		useStrict = true
		statement.Detach()
	}
	if useStrict {
		node.SetUseStrict(true)
	}
}

// Java method: TransformDispatcher.processBlock(com.google.javascript.jscomp.parsing.parser.trees.BlockTree blockNode)
func (this_ *IRFactory_TransformDispatcher) ProcessBlock(blockNode *pkg_parser.BlockTree) pkg_rhino.Node {
	_ = blockNode
	var node pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_BLOCK())
	_ = node
	{
		cur_52 := j2gort.NewIterCursor(j2gort.CheckNotNull(blockNode.FieldGet_statements()))
		for cur_52.HasNext() {
			rawElem_53 := cur_52.Next()
			child := j2gort.AsIface[pkg_parser.ParseTree](rawElem_53)
			_ = child
			var childNode pkg_rhino.Node = this_.outer_.Transform(child)
			_ = childNode
			node.AddChildToBack(childNode)
			this_.AttachPossibleTrailingComment(childNode, child.GetParseTreeBase().GetEnd())
		}
	}
	var lastComment *pkg_rhino.NonJSDocComment = this_.outer_.priv_IRFactory_parseNonJSDocCommentAt_6838(blockNode.GetEnd(), false)
	_ = lastComment
	this_.AddExtraTrailingComment(node, lastComment)
	return node
}

// Java method: TransformDispatcher.processBreakStatement(com.google.javascript.jscomp.parsing.parser.trees.BreakStatementTree statementNode)
func (this_ *IRFactory_TransformDispatcher) ProcessBreakStatement(statementNode *pkg_parser.BreakStatementTree) pkg_rhino.Node {
	_ = statementNode
	var node pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_BREAK())
	_ = node
	if statementNode.HasLabel() {
		var labelName pkg_rhino.Node = this_.TransformLabelName(statementNode.FieldGet_name())
		_ = labelName
		node.AddChildToBack(labelName)
	}
	return node
}

// Java method: TransformDispatcher.transformLabelName(com.google.javascript.jscomp.parsing.parser.Token token)
func (this_ *IRFactory_TransformDispatcher) TransformLabelName(token pkg_parser.Token) pkg_rhino.Node {
	_ = token
	var label pkg_rhino.Node = this_.outer_.MaybeAddNonJsDocComment(this_.NewStringNodeFromIdentifier(pkg_rhino.StaticGet_Token_LABEL_NAME(), StaticGet_IRFactory_TransformDispatcher_IdentifierType_STANDARD(), token), token.GetStart())
	_ = label
	this_.outer_.SetSourceInfo_2_Node_Token(label, token)
	return label
}

// Java method: TransformDispatcher.processConditionalExpression(com.google.javascript.jscomp.parsing.parser.trees.ConditionalExpressionTree exprNode)
func (this_ *IRFactory_TransformDispatcher) ProcessConditionalExpression(exprNode *pkg_parser.ConditionalExpressionTree) pkg_rhino.Node {
	_ = exprNode
	return this_.outer_.NewNode_4(pkg_rhino.StaticGet_Token_HOOK(), this_.outer_.Transform(exprNode.FieldGet_condition()), this_.outer_.Transform(exprNode.FieldGet_left()), this_.outer_.Transform(exprNode.FieldGet_right()))
}

// Java method: TransformDispatcher.processContinueStatement(com.google.javascript.jscomp.parsing.parser.trees.ContinueStatementTree statementNode)
func (this_ *IRFactory_TransformDispatcher) ProcessContinueStatement(statementNode *pkg_parser.ContinueStatementTree) pkg_rhino.Node {
	_ = statementNode
	var node pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_CONTINUE())
	_ = node
	if statementNode.HasLabel() {
		var labelName pkg_rhino.Node = this_.TransformLabelName(statementNode.FieldGet_name())
		_ = labelName
		node.AddChildToBack(labelName)
	}
	return node
}

// Java method: TransformDispatcher.processDoLoop(com.google.javascript.jscomp.parsing.parser.trees.DoWhileStatementTree loopNode)
func (this_ *IRFactory_TransformDispatcher) ProcessDoLoop(loopNode *pkg_parser.DoWhileStatementTree) pkg_rhino.Node {
	_ = loopNode
	return this_.outer_.NewNode_3(pkg_rhino.StaticGet_Token_DO(), this_.outer_.TransformBlock(loopNode.FieldGet_body()), this_.outer_.Transform(loopNode.FieldGet_condition()))
}

// Java method: TransformDispatcher.processElementGet(com.google.javascript.jscomp.parsing.parser.trees.MemberLookupExpressionTree getNode)
func (this_ *IRFactory_TransformDispatcher) ProcessElementGet(getNode *pkg_parser.MemberLookupExpressionTree) pkg_rhino.Node {
	_ = getNode
	return this_.outer_.NewNode_3(pkg_rhino.StaticGet_Token_GETELEM(), this_.outer_.Transform(getNode.FieldGet_operand()), this_.outer_.Transform(getNode.FieldGet_memberExpression()))
}

// Java method: TransformDispatcher.processOptChainElementGet(com.google.javascript.jscomp.parsing.parser.trees.OptionalMemberLookupExpressionTree getNode)
func (this_ *IRFactory_TransformDispatcher) ProcessOptChainElementGet(getNode *pkg_parser.OptionalMemberLookupExpressionTree) pkg_rhino.Node {
	_ = getNode
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(getNode), pkg_parser.StaticGet_FeatureSet_Feature_OPTIONAL_CHAINING())
	var getElem pkg_rhino.Node = this_.outer_.NewNode_3(pkg_rhino.StaticGet_Token_OPTCHAIN_GETELEM(), this_.outer_.Transform(getNode.FieldGet_operand()), this_.outer_.Transform(getNode.FieldGet_memberExpression()))
	_ = getElem
	getElem.SetIsOptionalChainStart(getNode.FieldGet_isStartOfOptionalChain())
	return getElem
}

// Java method: TransformDispatcher.processEmptyStatement(com.google.javascript.jscomp.parsing.parser.trees.EmptyStatementTree unused)
func (this_ *IRFactory_TransformDispatcher) ProcessEmptyStatement(unused *pkg_parser.EmptyStatementTree) pkg_rhino.Node {
	_ = unused
	return this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_EMPTY())
}

// Java method: TransformDispatcher.processExpressionStatement(com.google.javascript.jscomp.parsing.parser.trees.ExpressionStatementTree statementNode)
func (this_ *IRFactory_TransformDispatcher) ProcessExpressionStatement(statementNode *pkg_parser.ExpressionStatementTree) pkg_rhino.Node {
	_ = statementNode
	var node pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_EXPR_RESULT())
	_ = node
	node.AddChildToBack(this_.outer_.Transform(statementNode.FieldGet_expression()))
	return node
}

// Java method: TransformDispatcher.processForInLoop(com.google.javascript.jscomp.parsing.parser.trees.ForInStatementTree loopNode)
func (this_ *IRFactory_TransformDispatcher) ProcessForInLoop(loopNode *pkg_parser.ForInStatementTree) pkg_rhino.Node {
	_ = loopNode
	var initializer pkg_rhino.Node = this_.outer_.Transform(loopNode.FieldGet_initializer())
	_ = initializer
	return this_.outer_.NewNode_4(pkg_rhino.StaticGet_Token_FOR_IN(), initializer, this_.outer_.Transform(loopNode.FieldGet_collection()), this_.outer_.TransformBlock(loopNode.FieldGet_body()))
}

// Java method: TransformDispatcher.processForOf(com.google.javascript.jscomp.parsing.parser.trees.ForOfStatementTree loopNode)
func (this_ *IRFactory_TransformDispatcher) ProcessForOf(loopNode *pkg_parser.ForOfStatementTree) pkg_rhino.Node {
	_ = loopNode
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(loopNode), pkg_parser.StaticGet_FeatureSet_Feature_FOR_OF())
	var initializer pkg_rhino.Node = this_.outer_.Transform(loopNode.FieldGet_initializer())
	_ = initializer
	return this_.outer_.NewNode_4(pkg_rhino.StaticGet_Token_FOR_OF(), initializer, this_.outer_.Transform(loopNode.FieldGet_collection()), this_.outer_.TransformBlock(loopNode.FieldGet_body()))
}

// Java method: TransformDispatcher.processForAwaitOf(com.google.javascript.jscomp.parsing.parser.trees.ForAwaitOfStatementTree loopNode)
func (this_ *IRFactory_TransformDispatcher) ProcessForAwaitOf(loopNode *pkg_parser.ForAwaitOfStatementTree) pkg_rhino.Node {
	_ = loopNode
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(loopNode), pkg_parser.StaticGet_FeatureSet_Feature_FOR_AWAIT_OF())
	var initializer pkg_rhino.Node = this_.outer_.Transform(loopNode.FieldGet_initializer())
	_ = initializer
	return this_.outer_.NewNode_4(pkg_rhino.StaticGet_Token_FOR_AWAIT_OF(), initializer, this_.outer_.Transform(loopNode.FieldGet_collection()), this_.outer_.TransformBlock(loopNode.FieldGet_body()))
}

// Java method: TransformDispatcher.processForLoop(com.google.javascript.jscomp.parsing.parser.trees.ForStatementTree loopNode)
func (this_ *IRFactory_TransformDispatcher) ProcessForLoop(loopNode *pkg_parser.ForStatementTree) pkg_rhino.Node {
	_ = loopNode
	var node pkg_rhino.Node = this_.outer_.NewNode_4(pkg_rhino.StaticGet_Token_FOR(), this_.TransformOrEmpty_2_ParseTree_ParseTree(loopNode.FieldGet_initializer(), pkg_parser.ToIface_ParseTree(loopNode)), this_.TransformOrEmpty_2_ParseTree_ParseTree(loopNode.FieldGet_condition(), pkg_parser.ToIface_ParseTree(loopNode)), this_.TransformOrEmpty_2_ParseTree_ParseTree(loopNode.FieldGet_increment(), pkg_parser.ToIface_ParseTree(loopNode)))
	_ = node
	node.AddChildToBack(this_.outer_.TransformBlock(loopNode.FieldGet_body()))
	return node
}

// Java method: TransformDispatcher.transformOrEmpty(com.google.javascript.jscomp.parsing.parser.trees.ParseTree tree, com.google.javascript.jscomp.parsing.parser.trees.ParseTree parent)
func (this_ *IRFactory_TransformDispatcher) TransformOrEmpty_2_ParseTree_ParseTree(tree pkg_parser.ParseTree, parent pkg_parser.ParseTree) pkg_rhino.Node {
	_ = tree
	_ = parent
	if tree == nil || tree.IsJavaNil_() {
		var n pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_EMPTY())
		_ = n
		this_.outer_.SetSourceInfo_2_Node_ParseTree(n, parent)
		return n
	}
	return this_.outer_.Transform(tree)
}

// Java method: TransformDispatcher.transformOrEmpty(com.google.javascript.jscomp.parsing.parser.Token token, com.google.javascript.jscomp.parsing.parser.trees.ParseTree parent)
func (this_ *IRFactory_TransformDispatcher) TransformOrEmpty_2_Token_ParseTree(token pkg_parser.Token, parent pkg_parser.ParseTree) pkg_rhino.Node {
	_ = token
	_ = parent
	if (token).IsJavaNil_() {
		var n pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_EMPTY())
		_ = n
		this_.outer_.SetSourceInfo_2_Node_ParseTree(n, parent)
		return n
	}
	return this_.ProcessName_3(token, pkg_rhino.StaticGet_Token_NAME(), StaticGet_IRFactory_TransformDispatcher_IdentifierType_STANDARD())
}

// Java method: TransformDispatcher.processFunctionCall(com.google.javascript.jscomp.parsing.parser.trees.CallExpressionTree callNode)
func (this_ *IRFactory_TransformDispatcher) ProcessFunctionCall(callNode *pkg_parser.CallExpressionTree) pkg_rhino.Node {
	_ = callNode
	var node pkg_rhino.Node = this_.outer_.NewNode_2(pkg_rhino.StaticGet_Token_CALL(), this_.outer_.Transform(callNode.FieldGet_operand()))
	_ = node
	var argumentsTree *pkg_parser.ArgumentListTree = callNode.FieldGet_arguments()
	_ = argumentsTree
	node.SetTrailingComma(argumentsTree.FieldGet_hasTrailingComma())
	var size int32 = argumentsTree.FieldGet_argumentCount()
	_ = size
	if size > int32(0) {
		var firstArgZone pkg_parser.SourcePosition = this_.GetEndOfFirstArgCommentZone(argumentsTree)
		_ = firstArgZone
		{
			var i int32 = int32(0)
			_ = i
			for ; i < size; i = i + 1 {
				var childNode pkg_rhino.Node = this_.outer_.Transform(argumentsTree.GetArgument(i))
				_ = childNode
				node.AddChildToBack(childNode)
				this_.AttachPossibleTrailingCommentsForArg(childNode, firstArgZone)
			}
		}
	}
	this_.priv_TransformDispatcher_annotateCalls_6913(node)
	return node
}

// Java method: TransformDispatcher.annotateCalls(com.google.javascript.rhino.Node n)
func (this_ *IRFactory_TransformDispatcher) priv_TransformDispatcher_annotateCalls_6913(n pkg_rhino.Node) {
	_ = n
	if !((n.IsCall() || n.IsOptChainCall()) || n.IsTaggedTemplateLit()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	var callee pkg_rhino.Node = n.GetFirstChild()
	_ = callee
	for callee.IsCast() {
		callee = callee.GetFirstChild()
	}
	if !(this_.IsNormalOrOptChainGet(callee)) {
		n.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
		if callee.IsName() && (j2gort.CheckNotNullString(callee.GetString()) == "eval") {
			callee.PutBooleanProp(pkg_rhino.StaticGet_Node_DIRECT_EVAL(), true)
		} else if callee.IsComma() && callee.GetFirstChild().IsNumber() {
			var realCallee pkg_rhino.Node = callee.GetSecondChild()
			_ = realCallee
			callee.ReplaceWith(realCallee.Detach())
		}
	}
}

// Java method: TransformDispatcher.isNormalOrOptChainGet(com.google.javascript.rhino.Node n)
func (this_ *IRFactory_TransformDispatcher) IsNormalOrOptChainGet(n pkg_rhino.Node) bool {
	_ = n
	return (((n.IsGetProp() || n.IsGetElem()) || n.IsOptChainGetProp()) || n.IsOptChainGetElem())
}

// Java method: TransformDispatcher.getEndOfFirstArgCommentZone(com.google.javascript.jscomp.parsing.parser.trees.ArgumentListTree argsTree)
func (this_ *IRFactory_TransformDispatcher) GetEndOfFirstArgCommentZone(argsTree *pkg_parser.ArgumentListTree) pkg_parser.SourcePosition {
	_ = argsTree
	if argsTree.FieldGet_argumentCount() > int32(1) {
		var arg pkg_parser.ParseTree = argsTree.FieldGet_arg0()
		_ = arg
		var nextParam pkg_parser.ParseTree = argsTree.FieldGet_arg1()
		_ = nextParam
		if nextParam.GetParseTreeBase().FieldGet_location().FieldGet_start().FieldGet_line() > arg.GetParseTreeBase().FieldGet_location().FieldGet_end().FieldGet_line() {
			return pkg_parser.NewSourcePosition(nil, int32(2147483647), arg.GetParseTreeBase().FieldGet_location().FieldGet_end().FieldGet_line(), int32(2147483647))
		}
		return argsTree.GetCommaPosition(int32(0))
	}
	return argsTree.FieldGet_location().FieldGet_end()
}

// Java method: TransformDispatcher.getEndOfParamCommentZone(com.google.javascript.jscomp.parsing.parser.trees.FormalParameterListTree tree, int index)
func (this_ *IRFactory_TransformDispatcher) GetEndOfParamCommentZone(tree *pkg_parser.FormalParameterListTree, index int32) pkg_parser.SourcePosition {
	_ = tree
	_ = index
	if tree.FieldGet_parameterCount() > (index + int32(1)) {
		var arg pkg_parser.ParseTree = tree.GetParameter(index)
		_ = arg
		var nextParam pkg_parser.ParseTree = tree.GetParameter((index + int32(1)))
		_ = nextParam
		if nextParam.GetParseTreeBase().FieldGet_location().FieldGet_start().FieldGet_line() > arg.GetParseTreeBase().FieldGet_location().FieldGet_end().FieldGet_line() {
			return pkg_parser.NewSourcePosition(nil, int32(2147483647), arg.GetParseTreeBase().FieldGet_location().FieldGet_end().FieldGet_line(), int32(2147483647))
		}
		return tree.GetCommaPosition(index)
	}
	return tree.FieldGet_location().FieldGet_end()
}

// Java method: TransformDispatcher.attachPossibleTrailingCommentsForArg(com.google.javascript.rhino.Node paramNode, com.google.javascript.jscomp.parsing.parser.util.SourcePosition endZone)
func (this_ *IRFactory_TransformDispatcher) AttachPossibleTrailingCommentsForArg(paramNode pkg_rhino.Node, endZone pkg_parser.SourcePosition) {
	_ = paramNode
	_ = endZone
	var trailingComment *pkg_rhino.NonJSDocComment = this_.outer_.priv_IRFactory_parseNonJSDocCommentAt_6838(endZone, true)
	_ = trailingComment
	if trailingComment == nil {
		return
	}
	paramNode.SetTrailingNonJSDocComment(trailingComment)
}

// Java method: TransformDispatcher.attachPossibleTrailingComment(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.parsing.parser.util.SourcePosition tokenEnd)
func (this_ *IRFactory_TransformDispatcher) AttachPossibleTrailingComment(node pkg_rhino.Node, tokenEnd pkg_parser.SourcePosition) {
	_ = node
	_ = tokenEnd
	var nextLine pkg_parser.SourcePosition = pkg_parser.NewSourcePosition(nil, int32(2147483647), (tokenEnd.FieldGet_line() + int32(1)), int32(0))
	_ = nextLine
	var trailingComment *pkg_rhino.NonJSDocComment = this_.outer_.priv_IRFactory_parseTrailingNonJSDocCommentAt_6839(tokenEnd, nextLine)
	_ = trailingComment
	if trailingComment == nil {
		return
	}
	node.SetTrailingNonJSDocComment(trailingComment)
}

// Java method: TransformDispatcher.addExtraTrailingComment(com.google.javascript.rhino.Node node, com.google.javascript.rhino.NonJSDocComment lastComment)
func (this_ *IRFactory_TransformDispatcher) AddExtraTrailingComment(node pkg_rhino.Node, lastComment *pkg_rhino.NonJSDocComment) {
	_ = node
	_ = lastComment
	if lastComment == nil {
		return
	}
	if !(node.HasChildren()) {
		node.SetTrailingNonJSDocComment(lastComment)
		return
	}
	var lastChild pkg_rhino.Node = node.GetLastChild()
	_ = lastChild
	var currentComment *pkg_rhino.NonJSDocComment = lastChild.GetTrailingNonJSDocComment()
	_ = currentComment
	if currentComment == nil {
		var newStart pkg_parser.SourcePosition = pkg_parser.NewSourcePosition(nil, (lastComment.GetEndPosition().FieldGet_offset() - int32(1)), (lastComment.GetEndPosition().FieldGet_line() - int32(1)), int32(0))
		_ = newStart
		var newlineComment *pkg_rhino.NonJSDocComment = pkg_rhino.NewNonJSDocComment(newStart, lastComment.GetEndPosition(), ("\n" + j2gort.StrValString(lastComment.GetCommentString())))
		_ = newlineComment
		node.GetLastChild().SetTrailingNonJSDocComment(newlineComment)
		return
	}
	var blankLines int32 = (lastComment.GetStartPosition().FieldGet_line() - currentComment.GetEndPosition().FieldGet_line())
	_ = blankLines
	var numWhiteSpace int32 = int32(0)
	_ = numWhiteSpace
	if blankLines == int32(0) {
		numWhiteSpace = ((lastComment.GetStartPosition().FieldGet_column() - currentComment.GetEndPosition().FieldGet_column()) - int32(1))
	}
	var comment j2gort.LocalStringBuilder
	comment.AppendStr(currentComment.GetCommentString())

	{
		var i int32 = int32(0)
		_ = i
		for ; i < numWhiteSpace; i = i + 1 {
			comment.AppendStr(" ")
		}
	}
	{
		var i_2 int32 = int32(0)
		_ = i_2
		for ; i_2 < blankLines; i_2 = i_2 + 1 {
			comment.AppendStr("\n")
		}
	}
	comment.AppendStr(lastComment.GetCommentString())
	var allComments *pkg_rhino.NonJSDocComment = pkg_rhino.NewNonJSDocComment(currentComment.GetStartPosition(), lastComment.GetEndPosition(), comment.String())
	_ = allComments
	allComments.SetEndsAsLineComment(lastComment.IsEndingAsLineComment())
	allComments.SetIsInline(true)
	node.GetLastChild().SetTrailingNonJSDocComment(allComments)
}

// Java method: TransformDispatcher.processOptChainFunctionCall(com.google.javascript.jscomp.parsing.parser.trees.OptChainCallExpressionTree callNode)
func (this_ *IRFactory_TransformDispatcher) ProcessOptChainFunctionCall(callNode *pkg_parser.OptChainCallExpressionTree) pkg_rhino.Node {
	_ = callNode
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(callNode), pkg_parser.StaticGet_FeatureSet_Feature_OPTIONAL_CHAINING())
	var node pkg_rhino.Node = this_.outer_.NewNode_2(pkg_rhino.StaticGet_Token_OPTCHAIN_CALL(), this_.outer_.Transform(callNode.FieldGet_operand()))
	_ = node
	node.SetTrailingComma(callNode.FieldGet_hasTrailingComma())
	var optChainArgs *pkg_parser.ArgumentListTree = callNode.FieldGet_arguments()
	_ = optChainArgs
	{
		var i int32 = int32(0)
		_ = i
		var len_ int32 = optChainArgs.FieldGet_argumentCount()
		_ = len_
		for ; i < len_; i = i + 1 {
			node.AddChildToBack(this_.outer_.Transform(optChainArgs.GetArgument(i)))
		}
	}
	node.SetIsOptionalChainStart(callNode.FieldGet_isStartOfOptionalChain())
	this_.priv_TransformDispatcher_annotateCalls_6913(node)
	return node
}

// Java method: TransformDispatcher.processFunction(com.google.javascript.jscomp.parsing.parser.trees.FunctionDeclarationTree functionTree)
func (this_ *IRFactory_TransformDispatcher) ProcessFunction(functionTree *pkg_parser.FunctionDeclarationTree) pkg_rhino.Node {
	_ = functionTree
	var isDeclaration bool = (functionTree.FieldGet_kind() == pkg_parser.StaticGet_FunctionDeclarationTree_Kind_DECLARATION())
	_ = isDeclaration
	var isMember bool = (functionTree.FieldGet_kind() == pkg_parser.StaticGet_FunctionDeclarationTree_Kind_MEMBER())
	_ = isMember
	var isArrow bool = (functionTree.FieldGet_kind() == pkg_parser.StaticGet_FunctionDeclarationTree_Kind_ARROW())
	_ = isArrow
	var isAsync bool = functionTree.FieldGet_isAsync()
	_ = isAsync
	var isGenerator bool = functionTree.FieldGet_isGenerator()
	_ = isGenerator
	var isSignature bool = (functionTree.FieldGet_functionBody().GetParseTreeBase().FieldGet_type_() == pkg_parser.StaticGet_ParseTreeType_EMPTY_STATEMENT())
	_ = isSignature
	var identifierType IRFactory_TransformDispatcher_IdentifierType
	_ = identifierType
	if functionTree.FieldGet_isClassMember() {
		identifierType = StaticGet_IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE()
	} else {
		identifierType = StaticGet_IRFactory_TransformDispatcher_IdentifierType_STANDARD()
	}
	if isGenerator {
		this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(functionTree), pkg_parser.StaticGet_FeatureSet_Feature_GENERATORS())
	}
	if isMember {
		this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(functionTree), pkg_parser.StaticGet_FeatureSet_Feature_MEMBER_DECLARATIONS())
	}
	if isArrow {
		this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(functionTree), pkg_parser.StaticGet_FeatureSet_Feature_ARROW_FUNCTIONS())
	}
	if isAsync {
		this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(functionTree), pkg_parser.StaticGet_FeatureSet_Feature_ASYNC_FUNCTIONS())
	}
	if isGenerator && isAsync {
		this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(functionTree), pkg_parser.StaticGet_FeatureSet_Feature_ASYNC_GENERATORS())
	}
	var name pkg_parser.Token = functionTree.FieldGet_name()
	_ = name
	var newName pkg_rhino.Node = pkg_rhino.Node(0)
	_ = newName
	if !(name).IsJavaNil_() {
		newName = this_.ProcessNameWithInlineComments_2_Token_IRFactory_TransformDispatcher_IdentifierType(name, identifierType)
	} else {
		if isDeclaration || isMember {
			this_.outer_.errorReporter.M_Error("unnamed function statement", this_.outer_.sourceName, IRFactory_Lineno_1_ParseTree(pkg_parser.ToIface_ParseTree(functionTree)), IRFactory_Charno_1_ParseTree(pkg_parser.ToIface_ParseTree(functionTree)))
			newName = this_.priv_TransformDispatcher_createMissingNameNode_6983()
		} else {
			newName = this_.outer_.NewStringNode_2(pkg_rhino.StaticGet_Token_NAME(), "")
		}
		this_.outer_.SetSourceInfo_2_Node_ParseTree(newName, pkg_parser.ToIface_ParseTree(functionTree))
	}
	var node pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_FUNCTION())
	_ = node
	if isMember {
		newName.SetString("")
	}
	node.AddChildToBack(newName)
	node.AddChildToBack(this_.outer_.Transform(pkg_parser.ToIface_ParseTree(functionTree.FieldGet_formalParameterList())))
	var bodyNode pkg_rhino.Node = this_.outer_.Transform(functionTree.FieldGet_functionBody())
	_ = bodyNode
	if (!(isArrow) && !(isSignature)) && !(bodyNode.IsBlock()) {
		if !(this_.outer_.config.RunMode() == StaticGet_Config_RunMode_KEEP_GOING()) {
			j2gort.CheckState(false)
		}
		bodyNode = pkg_rhino.IR_Block_0()
	}
	this_.priv_TransformDispatcher_parseDirectives_6895(bodyNode)
	node.AddChildToBack(bodyNode)
	node.SetIsGeneratorFunction(isGenerator)
	node.SetIsArrowFunction(isArrow)
	node.SetIsAsyncFunction(isAsync)
	node.PutBooleanProp(pkg_rhino.StaticGet_Node_OPT_ES6_TYPED(), functionTree.FieldGet_isOptional())
	this_.AttachPossibleTrailingComment(node, functionTree.GetEnd())
	var result pkg_rhino.Node = pkg_rhino.Node(0)
	_ = result
	if isMember {
		this_.outer_.SetSourceInfo_2_Node_ParseTree(node, pkg_parser.ToIface_ParseTree(functionTree))
		var member pkg_rhino.Node = this_.NewStringNodeFromIdentifier(pkg_rhino.StaticGet_Token_MEMBER_FUNCTION_DEF(), identifierType, name)
		_ = member
		member.AddChildToBack(node)
		member.SetStaticMember(functionTree.FieldGet_isStatic())
		this_.outer_.SetSourceInfo_2_Node_Token(member, name)
		result = member
	} else {
		result = node
	}
	return result
}

// Java method: TransformDispatcher.processField(com.google.javascript.jscomp.parsing.parser.trees.FieldDeclarationTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessField(tree *pkg_parser.FieldDeclarationTree) pkg_rhino.Node {
	_ = tree
	if tree.FieldGet_name().IsPrivateIdentifier() {
		this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(tree), pkg_parser.StaticGet_FeatureSet_Feature_PRIVATE_ELEMENTS())
	} else {
		this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(tree), pkg_parser.StaticGet_FeatureSet_Feature_PUBLIC_CLASS_FIELDS())
	}
	var node pkg_rhino.Node = this_.outer_.MaybeAddNonJsDocComment(this_.NewStringNodeFromIdentifier(pkg_rhino.StaticGet_Token_MEMBER_FIELD_DEF(), StaticGet_IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE(), tree.FieldGet_name()), tree.GetStart())
	_ = node
	if !(j2gort.IfaceEquals[pkg_parser.ParseTree](tree.FieldGet_initializer(), nil)) {
		var initializer pkg_rhino.Node = this_.outer_.Transform(tree.FieldGet_initializer())
		_ = initializer
		node.AddChildToBack(initializer)
		this_.outer_.SetLength(node, tree.FieldGet_location().FieldGet_start(), tree.FieldGet_location().FieldGet_end())
	}
	node.PutBooleanProp(pkg_rhino.StaticGet_Node_STATIC_MEMBER(), tree.FieldGet_isStatic())
	return node
}

// Java method: TransformDispatcher.processComputedPropertyField(com.google.javascript.jscomp.parsing.parser.trees.ComputedPropertyFieldTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessComputedPropertyField(tree *pkg_parser.ComputedPropertyFieldTree) pkg_rhino.Node {
	_ = tree
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(tree), pkg_parser.StaticGet_FeatureSet_Feature_PUBLIC_CLASS_FIELDS())
	var key pkg_rhino.Node = this_.outer_.Transform(tree.FieldGet_property())
	_ = key
	var node pkg_rhino.Node
	_ = node
	if !(j2gort.IfaceEquals[pkg_parser.ParseTree](tree.FieldGet_initializer(), nil)) {
		node = this_.outer_.NewNode_3(pkg_rhino.StaticGet_Token_COMPUTED_FIELD_DEF(), key, this_.outer_.Transform(tree.FieldGet_initializer()))
	} else {
		node = this_.outer_.NewNode_2(pkg_rhino.StaticGet_Token_COMPUTED_FIELD_DEF(), key)
	}
	node.PutBooleanProp(pkg_rhino.StaticGet_Node_STATIC_MEMBER(), tree.FieldGet_isStatic())
	return node
}

// Java method: TransformDispatcher.processFormalParameterList(com.google.javascript.jscomp.parsing.parser.trees.FormalParameterListTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessFormalParameterList(tree *pkg_parser.FormalParameterListTree) pkg_rhino.Node {
	_ = tree
	var params pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_PARAM_LIST())
	_ = params
	params.SetTrailingComma(tree.FieldGet_hasTrailingComma())
	if !(this_.priv_TransformDispatcher_checkParameters_7006(tree)) {
		return params
	}
	var size int32 = tree.FieldGet_parameterCount()
	_ = size
	{
		var argCount int32 = int32(0)
		_ = argCount
		for ; argCount < size; argCount = argCount + 1 {
			var param pkg_parser.ParseTree = tree.GetParameter(argCount)
			_ = param
			var paramNode pkg_rhino.Node = pkg_rhino.Node(0)
			_ = paramNode
			switch (param.GetParseTreeBase().FieldGet_type_()).Ordinal() {
			case int32(71) /* ParseTreeType.DEFAULT_PARAMETER */ :
				paramNode = this_.ProcessDefaultParameter(param.GetParseTreeBase().AsDefaultParameter())
			case int32(61) /* ParseTreeType.ITER_REST */ :
				this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(param, pkg_parser.StaticGet_FeatureSet_Feature_REST_PARAMETERS())
				paramNode = this_.outer_.TransformNodeWithInlineComments(param)
			default:
				paramNode = this_.outer_.TransformNodeWithInlineComments(param)
				this_.AttachPossibleTrailingCommentsForArg(paramNode, this_.GetEndOfParamCommentZone(tree, argCount))
			}
			if !((((paramNode.IsName() || paramNode.IsRest()) || paramNode.IsArrayPattern()) || paramNode.IsObjectPattern()) || paramNode.IsDefaultValue()) {
				j2gort.CheckState(false)
			}
			params.AddChildToBack(paramNode)
		}
	}
	return params
}

// Java method: TransformDispatcher.processDefaultParameter(com.google.javascript.jscomp.parsing.parser.trees.DefaultParameterTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessDefaultParameter(tree *pkg_parser.DefaultParameterTree) pkg_rhino.Node {
	_ = tree
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(tree), pkg_parser.StaticGet_FeatureSet_Feature_DEFAULT_PARAMETERS())
	var targetTree pkg_parser.ParseTree = tree.FieldGet_lhs()
	_ = targetTree
	var targetNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = targetNode
	base_targetTree_ := targetTree.GetParseTreeBase()
	if base_targetTree_.FieldGet_type_() == pkg_parser.StaticGet_ParseTreeType_IDENTIFIER_EXPRESSION() {
		targetNode = this_.priv_TransformDispatcher_processNameWithInlineComments_6941(base_targetTree_.AsIdentifierExpression(), StaticGet_IRFactory_TransformDispatcher_IdentifierType_STANDARD())
	} else {
		targetNode = this_.outer_.TransformNodeWithInlineComments(targetTree)
	}
	var defaultValueExpression pkg_rhino.Node = this_.outer_.Transform(tree.FieldGet_defaultValue())
	_ = defaultValueExpression
	var defaultValueNode pkg_rhino.Node = this_.outer_.NewNode_3(pkg_rhino.StaticGet_Token_DEFAULT_VALUE(), targetNode, defaultValueExpression)
	_ = defaultValueNode
	this_.outer_.priv_IRFactory_reportErrorIfYieldOrAwaitInDefaultValue_6858(defaultValueNode)
	this_.outer_.SetSourceInfo_2_Node_ParseTree(defaultValueNode, pkg_parser.ToIface_ParseTree(tree))
	return defaultValueNode
}

// Java method: TransformDispatcher.processIterRest(com.google.javascript.jscomp.parsing.parser.trees.IterRestTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessIterRest(tree *pkg_parser.IterRestTree) pkg_rhino.Node {
	_ = tree
	var target pkg_rhino.Node = this_.outer_.TransformNodeWithInlineComments(tree.FieldGet_assignmentTarget())
	_ = target
	return this_.outer_.NewNode_2(pkg_rhino.StaticGet_Token_ITER_REST(), target)
}

// Java method: TransformDispatcher.processIterSpread(com.google.javascript.jscomp.parsing.parser.trees.IterSpreadTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessIterSpread(tree *pkg_parser.IterSpreadTree) pkg_rhino.Node {
	_ = tree
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(tree), pkg_parser.StaticGet_FeatureSet_Feature_SPREAD_EXPRESSIONS())
	return this_.outer_.NewNode_2(pkg_rhino.StaticGet_Token_ITER_SPREAD(), this_.outer_.Transform(tree.FieldGet_expression()))
}

// Java method: TransformDispatcher.processObjectSpread(com.google.javascript.jscomp.parsing.parser.trees.ObjectSpreadTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessObjectSpread(tree *pkg_parser.ObjectSpreadTree) pkg_rhino.Node {
	_ = tree
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(tree), pkg_parser.StaticGet_FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD())
	return this_.outer_.NewNode_2(pkg_rhino.StaticGet_Token_OBJECT_SPREAD(), this_.outer_.Transform(tree.FieldGet_expression()))
}

// Java method: TransformDispatcher.processIfStatement(com.google.javascript.jscomp.parsing.parser.trees.IfStatementTree statementNode)
func (this_ *IRFactory_TransformDispatcher) ProcessIfStatement(statementNode *pkg_parser.IfStatementTree) pkg_rhino.Node {
	_ = statementNode
	var node pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_IF())
	_ = node
	node.AddChildToBack(this_.outer_.Transform(statementNode.FieldGet_condition()))
	node.AddChildToBack(this_.outer_.TransformBlock(statementNode.FieldGet_ifClause()))
	if !(j2gort.IfaceEquals[pkg_parser.ParseTree](statementNode.FieldGet_elseClause(), nil)) {
		node.AddChildToBack(this_.outer_.TransformBlock(statementNode.FieldGet_elseClause()))
	}
	return node
}

// Java method: TransformDispatcher.markBinaryExpressionFeatures(com.google.javascript.jscomp.parsing.parser.trees.BinaryOperatorTree exprNode)
func (this_ *IRFactory_TransformDispatcher) priv_TransformDispatcher_markBinaryExpressionFeatures_6930(exprNode *pkg_parser.BinaryOperatorTree) {
	_ = exprNode
	switch (exprNode.FieldGet_operator().FieldGet_type_()).Ordinal() {
	case int32(66) /* TokenType.STAR_STAR */, int32(88) /* TokenType.STAR_STAR_EQUAL */ :
		this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(exprNode), pkg_parser.StaticGet_FeatureSet_Feature_EXPONENT_OP())
	case int32(81) /* TokenType.QUESTION_QUESTION */ :
		this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(exprNode), pkg_parser.StaticGet_FeatureSet_Feature_NULL_COALESCE_OP())
	case int32(100) /* TokenType.QUESTION_QUESTION_EQUAL */ :
		this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(exprNode), pkg_parser.StaticGet_FeatureSet_Feature_NULL_COALESCE_OP())
		fallthrough
	case int32(99) /* TokenType.OR_EQUAL */, int32(98) /* TokenType.AND_EQUAL */ :
		this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(exprNode), pkg_parser.StaticGet_FeatureSet_Feature_LOGICAL_ASSIGNMENT())
	default:
	}
}

// Java method: TransformDispatcher.processBinaryExpression(com.google.javascript.jscomp.parsing.parser.trees.BinaryOperatorTree exprNode)
func (this_ *IRFactory_TransformDispatcher) ProcessBinaryExpression(exprNode *pkg_parser.BinaryOperatorTree) pkg_rhino.Node {
	_ = exprNode
	if this_.outer_.jsdocTracker.HasPendingCommentBefore(exprNode.FieldGet_right().GetParseTreeBase().FieldGet_location().FieldGet_start()) || this_.outer_.nonJsdocTracker.HasPendingCommentBefore(exprNode.FieldGet_right().GetParseTreeBase().FieldGet_location().FieldGet_start()) {
		this_.priv_TransformDispatcher_markBinaryExpressionFeatures_6930(exprNode)
		return this_.outer_.NewNode_3(IRFactory_TransformBinaryTokenType(exprNode.FieldGet_operator().FieldGet_type_()), this_.priv_TransformDispatcher_transformLhsOfBinaryTree_6933(exprNode), this_.outer_.Transform(exprNode.FieldGet_right()))
	} else {
		return this_.priv_TransformDispatcher_processBinaryExpressionHelper_6932(exprNode)
	}
}

// Java method: TransformDispatcher.processBinaryExpressionHelper(com.google.javascript.jscomp.parsing.parser.trees.BinaryOperatorTree exprTree)
func (this_ *IRFactory_TransformDispatcher) priv_TransformDispatcher_processBinaryExpressionHelper_6932(exprTree *pkg_parser.BinaryOperatorTree) pkg_rhino.Node {
	_ = exprTree
	var root pkg_rhino.Node = pkg_rhino.Node(0)
	_ = root
	var current pkg_rhino.Node = pkg_rhino.Node(0)
	_ = current
	var previous pkg_rhino.Node = pkg_rhino.Node(0)
	_ = previous
	for exprTree != nil {
		this_.priv_TransformDispatcher_markBinaryExpressionFeatures_6930(exprTree)
		previous = current
		current = this_.outer_.NewNode_2(IRFactory_TransformBinaryTokenType(exprTree.FieldGet_operator().FieldGet_type_()), this_.outer_.Transform(exprTree.FieldGet_right()))
		this_.outer_.SetSourceInfo_2_Node_ParseTree(current, pkg_parser.ToIface_ParseTree(exprTree))
		if previous != 0 {
			previous.AddChildToFront(current)
		}
		var binaryOperatorTree *pkg_parser.BinaryOperatorTree
		_ = binaryOperatorTree
		var ok_54 bool
		binaryOperatorTree, ok_54 = j2gort.TryCast[*pkg_parser.BinaryOperatorTree](exprTree.FieldGet_left())
		if ok_54 {
			exprTree = binaryOperatorTree
		} else {
			var leftNode pkg_rhino.Node = this_.priv_TransformDispatcher_transformLhsOfBinaryTree_6933(exprTree)
			_ = leftNode
			current.AddChildToFront(leftNode)
			exprTree = nil
		}
		if root == 0 {
			root = current
		}
	}
	return root
}

// Java method: TransformDispatcher.transformLhsOfBinaryTree(com.google.javascript.jscomp.parsing.parser.trees.BinaryOperatorTree exprTree)
func (this_ *IRFactory_TransformDispatcher) priv_TransformDispatcher_transformLhsOfBinaryTree_6933(exprTree *pkg_parser.BinaryOperatorTree) pkg_rhino.Node {
	_ = exprTree
	var lhsOfInIsPrivateId bool = (((exprTree.FieldGet_operator().FieldGet_type_() == pkg_parser.StaticGet_TokenType_IN()) && j2gort.InstanceOf[*pkg_parser.IdentifierExpressionTree](exprTree.FieldGet_left())) && exprTree.FieldGet_left().GetParseTreeBase().AsIdentifierExpression().FieldGet_identifierToken().IsPrivateIdentifier())
	_ = lhsOfInIsPrivateId
	if !(lhsOfInIsPrivateId) {
		return this_.outer_.Transform(exprTree.FieldGet_left())
	}
	this_.privateIdLhsOfInScope.Enter()
	defer func() {
		this_.privateIdLhsOfInScope.Exit()
	}()
	{
		return this_.outer_.Transform(exprTree.FieldGet_left())
	}
	panic("unreachable")
}

// Java method: TransformDispatcher.processDebuggerStatement(com.google.javascript.jscomp.parsing.parser.trees.DebuggerStatementTree unused)
func (this_ *IRFactory_TransformDispatcher) ProcessDebuggerStatement(unused *pkg_parser.DebuggerStatementTree) pkg_rhino.Node {
	_ = unused
	return this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_DEBUGGER())
}

// Java method: TransformDispatcher.processThisExpression(com.google.javascript.jscomp.parsing.parser.trees.ThisExpressionTree unused)
func (this_ *IRFactory_TransformDispatcher) ProcessThisExpression(unused *pkg_parser.ThisExpressionTree) pkg_rhino.Node {
	_ = unused
	return this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_THIS())
}

// Java method: TransformDispatcher.processLabeledStatement(com.google.javascript.jscomp.parsing.parser.trees.LabelledStatementTree labelTree)
func (this_ *IRFactory_TransformDispatcher) ProcessLabeledStatement(labelTree *pkg_parser.LabelledStatementTree) pkg_rhino.Node {
	_ = labelTree
	var statement pkg_rhino.Node = this_.outer_.Transform(labelTree.FieldGet_statement())
	_ = statement
	if ((statement.IsFunction() || statement.IsClass()) || statement.IsLet()) || statement.IsConst() {
		this_.outer_.errorReporter.M_Error("Lexical declarations are only allowed at top level or inside a block.", this_.outer_.sourceName, IRFactory_Lineno_1_ParseTree(pkg_parser.ToIface_ParseTree(labelTree)), IRFactory_Charno_1_ParseTree(pkg_parser.ToIface_ParseTree(labelTree)))
		return statement
	}
	return this_.outer_.NewNode_3(pkg_rhino.StaticGet_Token_LABEL(), this_.TransformLabelName(labelTree.FieldGet_name()), statement)
}

// Java method: TransformDispatcher.processName(com.google.javascript.jscomp.parsing.parser.trees.IdentifierExpressionTree nameNode, com.google.javascript.jscomp.parsing.IRFactory.TransformDispatcher.IdentifierType identifierType)
func (this_ *IRFactory_TransformDispatcher) ProcessName_2(nameNode *pkg_parser.IdentifierExpressionTree, identifierType IRFactory_TransformDispatcher_IdentifierType) pkg_rhino.Node {
	_ = nameNode
	_ = identifierType
	return this_.ProcessName_3(nameNode.FieldGet_identifierToken(), pkg_rhino.StaticGet_Token_NAME(), identifierType)
}

// Java method: TransformDispatcher.processName(com.google.javascript.jscomp.parsing.parser.Token identifierToken, com.google.javascript.rhino.Token output, com.google.javascript.jscomp.parsing.IRFactory.TransformDispatcher.IdentifierType identifierType)
func (this_ *IRFactory_TransformDispatcher) ProcessName_3(identifierToken pkg_parser.Token, output pkg_rhino.Token, identifierType IRFactory_TransformDispatcher_IdentifierType) pkg_rhino.Node {
	_ = identifierToken
	_ = output
	_ = identifierType
	var comment *pkg_rhino.NonJSDocComment = this_.outer_.priv_IRFactory_parseNonJSDocCommentAt_6838(identifierToken.GetStart(), true)
	_ = comment
	var node pkg_rhino.Node = this_.NewStringNodeFromIdentifier(output, identifierType, identifierToken)
	_ = node
	if output == pkg_rhino.StaticGet_Token_NAME() {
		this_.priv_TransformDispatcher_maybeWarnReservedKeyword_6944(identifierToken)
		var info *pkg_rhino.JSDocInfo = this_.outer_.ParseJSDocInfoOnToken(identifierToken)
		_ = info
		if info != nil {
			node.SetJSDocInfo(info)
		}
	}
	if comment != nil {
		node.SetNonJSDocComment(comment)
	}
	this_.outer_.SetSourceInfo_2_Node_Token(node, identifierToken)
	return node
}

// Java method: TransformDispatcher.processString(com.google.javascript.jscomp.parsing.parser.Token token)
func (this_ *IRFactory_TransformDispatcher) ProcessString(token pkg_parser.Token) pkg_rhino.Node {
	_ = token
	if !(token.FieldGet_type_() == pkg_parser.StaticGet_TokenType_STRING()) {
		j2gort.CheckArgument(false)
	}
	var comment *pkg_rhino.NonJSDocComment = this_.outer_.priv_IRFactory_parseNonJSDocCommentAt_6838(token.GetStart(), true)
	_ = comment
	var node pkg_rhino.Node = this_.outer_.NewStringNode_2(pkg_rhino.StaticGet_Token_STRINGLIT(), this_.outer_.NormalizeString(token, false))
	_ = node
	if comment != nil {
		node.SetNonJSDocComment(comment)
	}
	this_.outer_.SetSourceInfo_2_Node_Token(node, token)
	return node
}

// Java method: TransformDispatcher.processTemplateLiteralToken(com.google.javascript.jscomp.parsing.parser.Token token)
func (this_ *IRFactory_TransformDispatcher) ProcessTemplateLiteralToken(token pkg_parser.Token) pkg_rhino.Node {
	_ = token
	if !((((token.FieldGet_type_() == pkg_parser.StaticGet_TokenType_NO_SUBSTITUTION_TEMPLATE()) || (token.FieldGet_type_() == pkg_parser.StaticGet_TokenType_TEMPLATE_HEAD())) || (token.FieldGet_type_() == pkg_parser.StaticGet_TokenType_TEMPLATE_MIDDLE())) || (token.FieldGet_type_() == pkg_parser.StaticGet_TokenType_TEMPLATE_TAIL())) {
		j2gort.CheckArgument(false)
	}
	var node pkg_rhino.Node = pkg_rhino.Node(0)
	_ = node
	if token.HasError() {
		node = this_.outer_.NewTemplateLitStringNode(j2gort.NullString, token.FieldGet_value())
	} else {
		node = this_.outer_.NewTemplateLitStringNode(this_.outer_.NormalizeString(token, true), token.FieldGet_value())
	}
	this_.outer_.SetSourceInfo_2_Node_Token(node, token)
	return node
}

// Java method: TransformDispatcher.processNameWithInlineComments(com.google.javascript.jscomp.parsing.parser.trees.IdentifierExpressionTree identifierExpression, com.google.javascript.jscomp.parsing.IRFactory.TransformDispatcher.IdentifierType identifierType)
func (this_ *IRFactory_TransformDispatcher) priv_TransformDispatcher_processNameWithInlineComments_6941(identifierExpression *pkg_parser.IdentifierExpressionTree, identifierType IRFactory_TransformDispatcher_IdentifierType) pkg_rhino.Node {
	_ = identifierExpression
	_ = identifierType
	return this_.ProcessNameWithInlineComments_2_Token_IRFactory_TransformDispatcher_IdentifierType(identifierExpression.FieldGet_identifierToken(), identifierType)
}

// Java method: TransformDispatcher.processNameWithInlineComments(com.google.javascript.jscomp.parsing.parser.Token identifierToken, com.google.javascript.jscomp.parsing.IRFactory.TransformDispatcher.IdentifierType identifierType)
func (this_ *IRFactory_TransformDispatcher) ProcessNameWithInlineComments_2_Token_IRFactory_TransformDispatcher_IdentifierType(identifierToken pkg_parser.Token, identifierType IRFactory_TransformDispatcher_IdentifierType) pkg_rhino.Node {
	_ = identifierToken
	_ = identifierType
	var info *pkg_rhino.JSDocInfo = this_.outer_.ParseInlineJSDocAt(identifierToken.GetStart())
	_ = info
	var comment *pkg_rhino.NonJSDocComment = this_.outer_.priv_IRFactory_parseNonJSDocCommentAt_6838(identifierToken.GetStart(), false)
	_ = comment
	this_.priv_TransformDispatcher_maybeWarnReservedKeyword_6944(identifierToken)
	var node pkg_rhino.Node = this_.NewStringNodeFromIdentifier(pkg_rhino.StaticGet_Token_NAME(), identifierType, identifierToken)
	_ = node
	if info != nil {
		node.SetJSDocInfo(info)
	}
	if comment != nil {
		node.SetNonJSDocComment(comment)
	}
	this_.outer_.SetSourceInfo_2_Node_Token(node, identifierToken)
	return node
}

// Java method: TransformDispatcher.maybeWarnKeywordProperty(com.google.javascript.rhino.Node node)
func (this_ *IRFactory_TransformDispatcher) priv_TransformDispatcher_maybeWarnKeywordProperty_6943(node pkg_rhino.Node) {
	_ = node
	if this_.outer_.currentFileIsExterns {
		return
	}
	if !(pkg_rhino.TokenStream_IsKeyword(node.GetString())) {
		return
	}
	this_.outer_.features = this_.outer_.features.With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_KEYWORDS_AS_PROPERTIES())
	if this_.outer_.config.LanguageMode() == StaticGet_Config_LanguageMode_ECMASCRIPT3() {
		this_.outer_.errorReporter.Warning("Keywords and reserved words are not allowed as unquoted property names in older versions of JavaScript. If you are targeting newer versions of JavaScript, set the appropriate language_in option.", this_.outer_.sourceName, node.GetLineno(), node.GetCharno())
	}
}

// Java method: TransformDispatcher.maybeWarnReservedKeyword(com.google.javascript.jscomp.parsing.parser.Token token)
func (this_ *IRFactory_TransformDispatcher) priv_TransformDispatcher_maybeWarnReservedKeyword_6944(token pkg_parser.Token) {
	_ = token
	var identifier string = token.GetMaybePrivateValue()
	_ = identifier
	var isIdentifier bool = false
	_ = isIdentifier
	if pkg_rhino.TokenStream_IsKeyword(identifier) {
		this_.outer_.features = this_.outer_.features.With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_ES3_KEYWORDS_AS_IDENTIFIERS())
		isIdentifier = (this_.outer_.config.LanguageMode() == StaticGet_Config_LanguageMode_ECMASCRIPT3())
	}
	if (this_.outer_.reservedKeywords != nil) && j2gort.SetContainsStr(this_.outer_.reservedKeywords, identifier) {
		this_.outer_.features = this_.outer_.features.With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_KEYWORDS_AS_PROPERTIES())
		isIdentifier = (this_.outer_.config.LanguageMode() == StaticGet_Config_LanguageMode_ECMASCRIPT3())
	}
	if isIdentifier {
		this_.outer_.errorReporter.M_Error("identifier is a reserved word", this_.outer_.sourceName, IRFactory_Lineno_1_SourcePosition(token.FieldGet_location().FieldGet_start()), IRFactory_Charno_1_SourcePosition(token.FieldGet_location().FieldGet_start()))
	}
}

// Java method: TransformDispatcher.processNewExpression(com.google.javascript.jscomp.parsing.parser.trees.NewExpressionTree exprNode)
func (this_ *IRFactory_TransformDispatcher) ProcessNewExpression(exprNode *pkg_parser.NewExpressionTree) pkg_rhino.Node {
	_ = exprNode
	var node pkg_rhino.Node = this_.outer_.NewNode_2(pkg_rhino.StaticGet_Token_NEW(), this_.outer_.Transform(exprNode.FieldGet_operand()))
	_ = node
	node.SetTrailingComma(exprNode.FieldGet_hasTrailingComma())
	if exprNode.FieldGet_arguments() != nil {
		var newArgs *pkg_parser.ArgumentListTree = exprNode.FieldGet_arguments()
		_ = newArgs
		{
			var i int32 = int32(0)
			_ = i
			var len_ int32 = newArgs.FieldGet_argumentCount()
			_ = len_
			for ; i < len_; i = i + 1 {
				node.AddChildToBack(this_.outer_.Transform(newArgs.GetArgument(i)))
			}
		}
	}
	return node
}

// Java method: TransformDispatcher.processNumberLiteral(com.google.javascript.jscomp.parsing.parser.trees.LiteralExpressionTree literalNode)
func (this_ *IRFactory_TransformDispatcher) ProcessNumberLiteral(literalNode *pkg_parser.LiteralExpressionTree) pkg_rhino.Node {
	_ = literalNode
	var value float64 = this_.outer_.NormalizeNumber(literalNode.FieldGet_literalToken())
	_ = value
	var number pkg_rhino.Node = this_.outer_.NewNumberNode(j2gort.ValueOfDouble(value))
	_ = number
	this_.outer_.SetSourceInfo_2_Node_ParseTree(number, pkg_parser.ToIface_ParseTree(literalNode))
	return number
}

// Java method: TransformDispatcher.processBigIntLiteral(com.google.javascript.jscomp.parsing.parser.trees.LiteralExpressionTree literalNode)
func (this_ *IRFactory_TransformDispatcher) ProcessBigIntLiteral(literalNode *pkg_parser.LiteralExpressionTree) pkg_rhino.Node {
	_ = literalNode
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(literalNode), pkg_parser.StaticGet_FeatureSet_Feature_BIGINT())
	var value *j2gort.Ext_java_math_BigInteger = this_.outer_.NormalizeBigInt(literalNode.FieldGet_literalToken())
	_ = value
	var bigint pkg_rhino.Node = this_.outer_.NewBigIntNode(value)
	_ = bigint
	this_.outer_.SetSourceInfo_2_Node_ParseTree(bigint, pkg_parser.ToIface_ParseTree(literalNode))
	return bigint
}

// Java method: TransformDispatcher.processObjectLiteral(com.google.javascript.jscomp.parsing.parser.trees.ObjectLiteralExpressionTree objTree)
func (this_ *IRFactory_TransformDispatcher) ProcessObjectLiteral(objTree *pkg_parser.ObjectLiteralExpressionTree) pkg_rhino.Node {
	_ = objTree
	var node pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_OBJECTLIT())
	_ = node
	node.SetTrailingComma(objTree.FieldGet_hasTrailingComma())
	var maybeWarn bool = false
	_ = maybeWarn
	{
		cur_55 := j2gort.NewIterCursor(j2gort.CheckNotNull(objTree.FieldGet_propertyNameAndValues()))
		for cur_55.HasNext() {
			rawElem_56 := cur_55.Next()
			el := j2gort.AsIface[pkg_parser.ParseTree](rawElem_56)
			_ = el
			if el.GetParseTreeBase().FieldGet_type_() == pkg_parser.StaticGet_ParseTreeType_DEFAULT_PARAMETER() {
				this_.outer_.errorReporter.M_Error("Default value cannot appear at top level of an object literal.", this_.outer_.sourceName, IRFactory_Lineno_1_ParseTree(el), int32(0))
				continue
			} else if (el.GetParseTreeBase().FieldGet_type_() == pkg_parser.StaticGet_ParseTreeType_GET_ACCESSOR()) && this_.MaybeReportGetter(el) {
				continue
			} else if (el.GetParseTreeBase().FieldGet_type_() == pkg_parser.StaticGet_ParseTreeType_SET_ACCESSOR()) && this_.MaybeReportSetter(el) {
				continue
			}
			var key pkg_rhino.Node = this_.outer_.Transform(el)
			_ = key
			if ((!(key.IsComputedProp()) && !(key.IsQuotedStringKey())) && !(key.IsSpread())) && !(this_.outer_.currentFileIsExterns) {
				this_.priv_TransformDispatcher_maybeWarnKeywordProperty_6943(key)
			}
			if key.IsShorthandProperty() {
				maybeWarn = true
			}
			node.AddChildToBack(key)
		}
	}
	if maybeWarn {
		this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(objTree), pkg_parser.StaticGet_FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES())
	}
	return node
}

// Java method: TransformDispatcher.processComputedPropertyDefinition(com.google.javascript.jscomp.parsing.parser.trees.ComputedPropertyDefinitionTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessComputedPropertyDefinition(tree *pkg_parser.ComputedPropertyDefinitionTree) pkg_rhino.Node {
	_ = tree
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(tree), pkg_parser.StaticGet_FeatureSet_Feature_COMPUTED_PROPERTIES())
	return this_.outer_.NewNode_3(pkg_rhino.StaticGet_Token_COMPUTED_PROP(), this_.outer_.Transform(tree.FieldGet_property()), this_.outer_.Transform(tree.FieldGet_value()))
}

// Java method: TransformDispatcher.processComputedPropertyMethod(com.google.javascript.jscomp.parsing.parser.trees.ComputedPropertyMethodTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessComputedPropertyMethod(tree *pkg_parser.ComputedPropertyMethodTree) pkg_rhino.Node {
	_ = tree
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(tree), pkg_parser.StaticGet_FeatureSet_Feature_COMPUTED_PROPERTIES())
	var n pkg_rhino.Node = this_.outer_.NewNode_3(pkg_rhino.StaticGet_Token_COMPUTED_PROP(), this_.outer_.Transform(tree.FieldGet_property()), this_.outer_.Transform(tree.FieldGet_method()))
	_ = n
	n.PutBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_METHOD(), true)
	if tree.FieldGet_method().GetParseTreeBase().AsFunctionDeclaration().FieldGet_isStatic() {
		n.SetStaticMember(true)
	}
	return n
}

// Java method: TransformDispatcher.processComputedPropertyGetter(com.google.javascript.jscomp.parsing.parser.trees.ComputedPropertyGetterTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessComputedPropertyGetter(tree *pkg_parser.ComputedPropertyGetterTree) pkg_rhino.Node {
	_ = tree
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(tree), pkg_parser.StaticGet_FeatureSet_Feature_COMPUTED_PROPERTIES())
	var key pkg_rhino.Node = this_.outer_.Transform(tree.FieldGet_property())
	_ = key
	var body pkg_rhino.Node = this_.outer_.Transform(pkg_parser.ToIface_ParseTree(tree.FieldGet_body()))
	_ = body
	var function pkg_rhino.Node = pkg_rhino.IR_Function(pkg_rhino.IR_Name(""), pkg_rhino.IR_ParamList([]pkg_rhino.Node{}), body)
	_ = function
	function.SrcrefTreeIfMissing(body)
	var n pkg_rhino.Node = this_.outer_.NewNode_3(pkg_rhino.StaticGet_Token_COMPUTED_PROP(), key, function)
	_ = n
	n.PutBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_GETTER(), true)
	n.PutBooleanProp(pkg_rhino.StaticGet_Node_STATIC_MEMBER(), tree.FieldGet_isStatic())
	return n
}

// Java method: TransformDispatcher.processComputedPropertySetter(com.google.javascript.jscomp.parsing.parser.trees.ComputedPropertySetterTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessComputedPropertySetter(tree *pkg_parser.ComputedPropertySetterTree) pkg_rhino.Node {
	_ = tree
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(tree), pkg_parser.StaticGet_FeatureSet_Feature_COMPUTED_PROPERTIES())
	var key pkg_rhino.Node = this_.outer_.Transform(tree.FieldGet_property())
	_ = key
	var paramList pkg_rhino.Node = this_.ProcessFormalParameterList(tree.FieldGet_parameter())
	_ = paramList
	this_.outer_.SetSourceInfo_2_Node_ParseTree(paramList, pkg_parser.ToIface_ParseTree(tree.FieldGet_parameter()))
	var body pkg_rhino.Node = this_.outer_.Transform(pkg_parser.ToIface_ParseTree(tree.FieldGet_body()))
	_ = body
	var function pkg_rhino.Node = pkg_rhino.IR_Function(pkg_rhino.IR_Name(""), paramList, body)
	_ = function
	function.SrcrefTreeIfMissing(body)
	var n pkg_rhino.Node = this_.outer_.NewNode_3(pkg_rhino.StaticGet_Token_COMPUTED_PROP(), key, function)
	_ = n
	n.PutBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_SETTER(), true)
	n.PutBooleanProp(pkg_rhino.StaticGet_Node_STATIC_MEMBER(), tree.FieldGet_isStatic())
	return n
}

// Java method: TransformDispatcher.processGetAccessor(com.google.javascript.jscomp.parsing.parser.trees.GetAccessorTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessGetAccessor(tree *pkg_parser.GetAccessorTree) pkg_rhino.Node {
	_ = tree
	var key pkg_rhino.Node = this_.priv_TransformDispatcher_processObjectLitKey_6884(tree.FieldGet_propertyName(), pkg_rhino.StaticGet_Token_GETTER_DEF(), func() IRFactory_TransformDispatcher_IdentifierType {
		if tree.FieldGet_isClassMember() {
			return StaticGet_IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE()
		}
		return StaticGet_IRFactory_TransformDispatcher_IdentifierType_STANDARD()
	}())
	_ = key
	var body pkg_rhino.Node = this_.outer_.Transform(pkg_parser.ToIface_ParseTree(tree.FieldGet_body()))
	_ = body
	var dummyName pkg_rhino.Node = this_.outer_.NewStringNode_2(pkg_rhino.StaticGet_Token_NAME(), "")
	_ = dummyName
	this_.outer_.SetSourceInfo_2_Node_ParseTree(dummyName, pkg_parser.ToIface_ParseTree(tree.FieldGet_body()))
	var paramList pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_PARAM_LIST())
	_ = paramList
	this_.outer_.SetSourceInfo_2_Node_ParseTree(paramList, pkg_parser.ToIface_ParseTree(tree.FieldGet_body()))
	var value pkg_rhino.Node = this_.outer_.NewNode_4(pkg_rhino.StaticGet_Token_FUNCTION(), dummyName, paramList, body)
	_ = value
	this_.outer_.SetSourceInfo_2_Node_ParseTree(value, pkg_parser.ToIface_ParseTree(tree.FieldGet_body()))
	key.AddChildToFront(value)
	key.SetStaticMember(tree.FieldGet_isStatic())
	return key
}

// Java method: TransformDispatcher.processSetAccessor(com.google.javascript.jscomp.parsing.parser.trees.SetAccessorTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessSetAccessor(tree *pkg_parser.SetAccessorTree) pkg_rhino.Node {
	_ = tree
	var key pkg_rhino.Node = this_.priv_TransformDispatcher_processObjectLitKey_6884(tree.FieldGet_propertyName(), pkg_rhino.StaticGet_Token_SETTER_DEF(), func() IRFactory_TransformDispatcher_IdentifierType {
		if tree.FieldGet_isClassMember() {
			return StaticGet_IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE()
		}
		return StaticGet_IRFactory_TransformDispatcher_IdentifierType_STANDARD()
	}())
	_ = key
	var paramList pkg_rhino.Node = this_.ProcessFormalParameterList(tree.FieldGet_parameter())
	_ = paramList
	this_.outer_.SetSourceInfo_2_Node_ParseTree(paramList, pkg_parser.ToIface_ParseTree(tree.FieldGet_parameter()))
	var body pkg_rhino.Node = this_.outer_.Transform(pkg_parser.ToIface_ParseTree(tree.FieldGet_body()))
	_ = body
	var dummyName pkg_rhino.Node = this_.outer_.NewStringNode_2(pkg_rhino.StaticGet_Token_NAME(), "")
	_ = dummyName
	this_.outer_.SetSourceInfo_2_Node_Token(dummyName, tree.FieldGet_propertyName())
	var value pkg_rhino.Node = this_.outer_.NewNode_4(pkg_rhino.StaticGet_Token_FUNCTION(), dummyName, paramList, body)
	_ = value
	this_.outer_.SetSourceInfo_2_Node_ParseTree(value, pkg_parser.ToIface_ParseTree(tree.FieldGet_body()))
	key.AddChildToFront(value)
	key.SetStaticMember(tree.FieldGet_isStatic())
	return key
}

// Java method: TransformDispatcher.processPropertyNameAssignment(com.google.javascript.jscomp.parsing.parser.trees.PropertyNameAssignmentTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessPropertyNameAssignment(tree *pkg_parser.PropertyNameAssignmentTree) pkg_rhino.Node {
	_ = tree
	var key pkg_rhino.Node = this_.priv_TransformDispatcher_processObjectLitKey_6884(tree.FieldGet_name(), pkg_rhino.StaticGet_Token_STRING_KEY(), StaticGet_IRFactory_TransformDispatcher_IdentifierType_STANDARD())
	_ = key
	if !(j2gort.IfaceEquals[pkg_parser.ParseTree](tree.FieldGet_value(), nil)) {
		key.AddChildToFront(this_.outer_.Transform(tree.FieldGet_value()))
	} else {
		var value pkg_rhino.Node = this_.outer_.MaybeAddNonJsDocComment(this_.outer_.NewStringNode_2(pkg_rhino.StaticGet_Token_NAME(), key.GetString()), tree.FieldGet_name().GetStart()).Srcref(key)
		_ = value
		key.SetShorthandProperty(true)
		key.AddChildToFront(value)
	}
	return key
}

// Java method: TransformDispatcher.checkParenthesizedExpression(com.google.javascript.jscomp.parsing.parser.trees.ParenExpressionTree exprNode)
func (this_ *IRFactory_TransformDispatcher) priv_TransformDispatcher_checkParenthesizedExpression_6956(exprNode *pkg_parser.ParenExpressionTree) {
	_ = exprNode
	if exprNode.FieldGet_expression().GetParseTreeBase().FieldGet_type_() == pkg_parser.StaticGet_ParseTreeType_COMMA_EXPRESSION() {
		var commaNodes []pkg_parser.ParseTree = j2gort.ImmListToSliceNullable[pkg_parser.ParseTree](exprNode.FieldGet_expression().GetParseTreeBase().AsCommaExpression().FieldGet_expressions())
		_ = commaNodes
		var lastChild pkg_parser.ParseTree = j2gort.SliceGetLast(commaNodes)
		_ = lastChild
		if lastChild.GetParseTreeBase().FieldGet_type_() == pkg_parser.StaticGet_ParseTreeType_ITER_REST() {
			this_.outer_.errorReporter.M_Error("A rest parameter must be in a parameter list.", this_.outer_.sourceName, IRFactory_Lineno_1_ParseTree(lastChild), IRFactory_Charno_1_ParseTree(lastChild))
		}
	}
}

// Java method: TransformDispatcher.processParenthesizedExpression(com.google.javascript.jscomp.parsing.parser.trees.ParenExpressionTree exprNode)
func (this_ *IRFactory_TransformDispatcher) ProcessParenthesizedExpression(exprNode *pkg_parser.ParenExpressionTree) pkg_rhino.Node {
	_ = exprNode
	this_.priv_TransformDispatcher_checkParenthesizedExpression_6956(exprNode)
	var expr pkg_rhino.Node = this_.outer_.Transform(exprNode.FieldGet_expression())
	_ = expr
	expr.SetIsParenthesized(true)
	return expr
}

// Java method: TransformDispatcher.processPropertyGet(com.google.javascript.jscomp.parsing.parser.trees.MemberExpressionTree getNode)
func (this_ *IRFactory_TransformDispatcher) ProcessPropertyGet(getNode *pkg_parser.MemberExpressionTree) pkg_rhino.Node {
	_ = getNode
	var leftChild pkg_rhino.Node = this_.outer_.Transform(getNode.FieldGet_operand())
	_ = leftChild
	var propName pkg_parser.Token = getNode.FieldGet_memberName()
	_ = propName
	if (propName).IsJavaNil_() {
		return leftChild
	}
	var identifierType IRFactory_TransformDispatcher_IdentifierType
	_ = identifierType
	if leftChild.IsSuper() {
		identifierType = StaticGet_IRFactory_TransformDispatcher_IdentifierType_STANDARD()
	} else {
		identifierType = StaticGet_IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE()
	}
	var getProp pkg_rhino.Node = this_.outer_.MaybeAddNonJsDocComment(this_.NewStringNodeFromIdentifier(pkg_rhino.StaticGet_Token_GETPROP(), identifierType, propName), getNode.FieldGet_memberName().GetStart())
	_ = getProp
	getProp.AddChildToBack(leftChild)
	this_.outer_.SetSourceInfo_2_Node_Token(getProp, propName)
	this_.priv_TransformDispatcher_maybeWarnKeywordProperty_6943(getProp)
	return getProp
}

// Java method: TransformDispatcher.processOptChainPropertyGet(com.google.javascript.jscomp.parsing.parser.trees.OptionalMemberExpressionTree getNode)
func (this_ *IRFactory_TransformDispatcher) ProcessOptChainPropertyGet(getNode *pkg_parser.OptionalMemberExpressionTree) pkg_rhino.Node {
	_ = getNode
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(getNode), pkg_parser.StaticGet_FeatureSet_Feature_OPTIONAL_CHAINING())
	var leftChild pkg_rhino.Node = this_.outer_.Transform(getNode.FieldGet_operand())
	_ = leftChild
	var propName pkg_parser.Token = getNode.FieldGet_memberName()
	_ = propName
	if (propName).IsJavaNil_() {
		return leftChild
	}
	var identifierType IRFactory_TransformDispatcher_IdentifierType
	_ = identifierType
	if leftChild.IsSuper() {
		identifierType = StaticGet_IRFactory_TransformDispatcher_IdentifierType_STANDARD()
	} else {
		identifierType = StaticGet_IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE()
	}
	var getProp pkg_rhino.Node = this_.outer_.MaybeAddNonJsDocComment(this_.NewStringNodeFromIdentifier(pkg_rhino.StaticGet_Token_OPTCHAIN_GETPROP(), identifierType, propName), getNode.FieldGet_memberName().GetStart())
	_ = getProp
	getProp.AddChildToBack(leftChild)
	getProp.SetIsOptionalChainStart(getNode.FieldGet_isStartOfOptionalChain())
	this_.outer_.SetSourceInfo_2_Node_Token(getProp, propName)
	this_.priv_TransformDispatcher_maybeWarnKeywordProperty_6943(getProp)
	return getProp
}

// Java method: TransformDispatcher.processRegExpLiteral(com.google.javascript.jscomp.parsing.parser.trees.LiteralExpressionTree literalTree)
func (this_ *IRFactory_TransformDispatcher) ProcessRegExpLiteral(literalTree *pkg_parser.LiteralExpressionTree) pkg_rhino.Node {
	_ = literalTree
	var token pkg_parser.Token = literalTree.FieldGet_literalToken()
	_ = token
	var literalStringNode pkg_rhino.Node = this_.outer_.NewStringNode_1(this_.outer_.NormalizeRegex(token))
	_ = literalStringNode
	this_.outer_.SetSourceInfo_2_Node_Token(literalStringNode, token)
	var node pkg_rhino.Node = this_.outer_.NewNode_2(pkg_rhino.StaticGet_Token_REGEXP(), literalStringNode)
	_ = node
	var rawRegex string = token.FieldGet_value()
	_ = rawRegex
	var lastSlash int32 = j2gort.JStrLastIndexOfCp(rawRegex, int32(uint16(47)))
	_ = lastSlash
	var flags string = ""
	_ = flags
	if lastSlash < j2gort.JStrLen(rawRegex) {
		flags = j2gort.JStrSubFrom(rawRegex, (lastSlash + int32(1)))
	}
	this_.priv_TransformDispatcher_validateRegExpFlags_6961(literalTree, flags)
	if !(j2gort.CheckNotNullString(flags) == "") {
		var flagsNode pkg_rhino.Node = this_.outer_.NewStringNode_1(flags)
		_ = flagsNode
		this_.outer_.SetSourceInfo_2_Node_Token(flagsNode, token)
		node.AddChildToBack(flagsNode)
	}
	return node
}

// Java method: TransformDispatcher.validateRegExpFlags(com.google.javascript.jscomp.parsing.parser.trees.LiteralExpressionTree tree, java.lang.String flags)
func (this_ *IRFactory_TransformDispatcher) priv_TransformDispatcher_validateRegExpFlags_6961(tree *pkg_parser.LiteralExpressionTree, flags string) {
	_ = tree
	_ = flags
	{
		cur_57 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.Ext_com_google_common_collect_Lists_CharactersOf_1_string(flags)))
		for cur_57.HasNext() {
			rawElem_58 := cur_57.Next()
			flag := j2gort.UnboxChar(j2gort.Cast[*j2gort.BoxChar](rawElem_58))
			_ = flag
			switch flag {
			case uint16(103), uint16(105), uint16(109):
			case uint16(117), uint16(121):
				var feature *pkg_parser.FeatureSet_Feature
				_ = feature
				if int32(flag) == int32(uint16(117)) {
					feature = pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_FLAG_U()
				} else {
					feature = pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_FLAG_Y()
				}
				this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(tree), feature)
			case uint16(115):
				this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(tree), pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_FLAG_S())
			case uint16(100):
				this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(tree), pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_FLAG_D())
			default:
				this_.outer_.errorReporter.M_Error(("Invalid RegExp flag '" + j2gort.StrValChar(flag) + "'"), this_.outer_.sourceName, IRFactory_Lineno_1_ParseTree(pkg_parser.ToIface_ParseTree(tree)), IRFactory_Charno_1_ParseTree(pkg_parser.ToIface_ParseTree(tree)))
			}
		}
	}
}

// Java method: TransformDispatcher.processReturnStatement(com.google.javascript.jscomp.parsing.parser.trees.ReturnStatementTree statementNode)
func (this_ *IRFactory_TransformDispatcher) ProcessReturnStatement(statementNode *pkg_parser.ReturnStatementTree) pkg_rhino.Node {
	_ = statementNode
	var node pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_RETURN())
	_ = node
	if !(j2gort.IfaceEquals[pkg_parser.ParseTree](statementNode.FieldGet_expression(), nil)) {
		node.AddChildToBack(this_.outer_.Transform(statementNode.FieldGet_expression()))
	}
	return node
}

// Java method: TransformDispatcher.processStringLiteral(com.google.javascript.jscomp.parsing.parser.trees.LiteralExpressionTree literalTree)
func (this_ *IRFactory_TransformDispatcher) ProcessStringLiteral(literalTree *pkg_parser.LiteralExpressionTree) pkg_rhino.Node {
	_ = literalTree
	return this_.ProcessString(literalTree.FieldGet_literalToken())
}

// Java method: TransformDispatcher.processTemplateLiteral(com.google.javascript.jscomp.parsing.parser.trees.TemplateLiteralExpressionTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessTemplateLiteral(tree *pkg_parser.TemplateLiteralExpressionTree) pkg_rhino.Node {
	_ = tree
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(tree), pkg_parser.StaticGet_FeatureSet_Feature_TEMPLATE_LITERALS())
	var templateLitNode pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_TEMPLATELIT())
	_ = templateLitNode
	this_.outer_.SetSourceInfo_2_Node_ParseTree(templateLitNode, pkg_parser.ToIface_ParseTree(tree))
	var node pkg_rhino.Node
	_ = node
	if j2gort.IfaceEquals[pkg_parser.ParseTree](tree.FieldGet_operand(), nil) {
		node = templateLitNode
	} else {
		node = this_.outer_.NewNode_3(pkg_rhino.StaticGet_Token_TAGGED_TEMPLATELIT(), this_.outer_.Transform(tree.FieldGet_operand()), templateLitNode)
	}
	{
		cur_59 := j2gort.NewIterCursor(j2gort.CheckNotNull(tree.FieldGet_elements()))
		for cur_59.HasNext() {
			rawElem_60 := cur_59.Next()
			child := j2gort.AsIface[pkg_parser.ParseTree](rawElem_60)
			_ = child
			templateLitNode.AddChildToBack(this_.outer_.Transform(child))
		}
	}
	if node.IsTaggedTemplateLit() {
		this_.priv_TransformDispatcher_annotateCalls_6913(node)
	}
	return node
}

// Java method: TransformDispatcher.processTemplateLiteralPortion(com.google.javascript.jscomp.parsing.parser.trees.TemplateLiteralPortionTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessTemplateLiteralPortion(tree *pkg_parser.TemplateLiteralPortionTree) pkg_rhino.Node {
	_ = tree
	return this_.ProcessTemplateLiteralToken(tree.FieldGet_value())
}

// Java method: TransformDispatcher.processTemplateSubstitution(com.google.javascript.jscomp.parsing.parser.trees.TemplateSubstitutionTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessTemplateSubstitution(tree *pkg_parser.TemplateSubstitutionTree) pkg_rhino.Node {
	_ = tree
	return this_.outer_.NewNode_2(pkg_rhino.StaticGet_Token_TEMPLATELIT_SUB(), this_.outer_.Transform(tree.FieldGet_expression()))
}

// Java method: TransformDispatcher.processSwitchCase(com.google.javascript.jscomp.parsing.parser.trees.CaseClauseTree caseNode)
func (this_ *IRFactory_TransformDispatcher) ProcessSwitchCase(caseNode *pkg_parser.CaseClauseTree) pkg_rhino.Node {
	_ = caseNode
	var expr pkg_parser.ParseTree = caseNode.FieldGet_expression()
	_ = expr
	var node pkg_rhino.Node = this_.outer_.NewNode_2(pkg_rhino.StaticGet_Token_CASE(), this_.outer_.Transform(expr))
	_ = node
	var block pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_BLOCK())
	_ = block
	block.SetIsAddedBlock(true)
	this_.outer_.SetSourceInfo_2_Node_ParseTree(block, pkg_parser.ToIface_ParseTree(caseNode))
	if caseNode.FieldGet_statements() != nil {
		{
			cur_61 := j2gort.NewIterCursor(j2gort.CheckNotNull(caseNode.FieldGet_statements()))
			for cur_61.HasNext() {
				rawElem_62 := cur_61.Next()
				child := j2gort.AsIface[pkg_parser.ParseTree](rawElem_62)
				_ = child
				block.AddChildToBack(this_.outer_.Transform(child))
			}
		}
	}
	node.AddChildToBack(block)
	return node
}

// Java method: TransformDispatcher.processSwitchDefault(com.google.javascript.jscomp.parsing.parser.trees.DefaultClauseTree caseNode)
func (this_ *IRFactory_TransformDispatcher) ProcessSwitchDefault(caseNode *pkg_parser.DefaultClauseTree) pkg_rhino.Node {
	_ = caseNode
	var node pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_DEFAULT_CASE())
	_ = node
	var block pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_BLOCK())
	_ = block
	block.SetIsAddedBlock(true)
	this_.outer_.SetSourceInfo_2_Node_ParseTree(block, pkg_parser.ToIface_ParseTree(caseNode))
	if caseNode.FieldGet_statements() != nil {
		{
			cur_63 := j2gort.NewIterCursor(j2gort.CheckNotNull(caseNode.FieldGet_statements()))
			for cur_63.HasNext() {
				rawElem_64 := cur_63.Next()
				child := j2gort.AsIface[pkg_parser.ParseTree](rawElem_64)
				_ = child
				block.AddChildToBack(this_.outer_.Transform(child))
			}
		}
	}
	node.AddChildToBack(block)
	return node
}

// Java method: TransformDispatcher.processSwitchStatement(com.google.javascript.jscomp.parsing.parser.trees.SwitchStatementTree statementNode)
func (this_ *IRFactory_TransformDispatcher) ProcessSwitchStatement(statementNode *pkg_parser.SwitchStatementTree) pkg_rhino.Node {
	_ = statementNode
	var node pkg_rhino.Node = this_.outer_.NewNode_2(pkg_rhino.StaticGet_Token_SWITCH(), this_.outer_.Transform(statementNode.FieldGet_expression()))
	_ = node
	var switchBody pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_SWITCH_BODY())
	_ = switchBody
	this_.outer_.SetSourceInfo_2_Node_ParseTree(switchBody, pkg_parser.ToIface_ParseTree(statementNode))
	node.AddChildToBack(switchBody)
	{
		cur_65 := j2gort.NewIterCursor(j2gort.CheckNotNull(statementNode.FieldGet_caseClauses()))
		for cur_65.HasNext() {
			rawElem_66 := cur_65.Next()
			child := j2gort.AsIface[pkg_parser.ParseTree](rawElem_66)
			_ = child
			switchBody.AddChildToBack(this_.outer_.Transform(child))
		}
	}
	return node
}

// Java method: TransformDispatcher.processThrowStatement(com.google.javascript.jscomp.parsing.parser.trees.ThrowStatementTree statementNode)
func (this_ *IRFactory_TransformDispatcher) ProcessThrowStatement(statementNode *pkg_parser.ThrowStatementTree) pkg_rhino.Node {
	_ = statementNode
	return this_.outer_.NewNode_2(pkg_rhino.StaticGet_Token_THROW(), this_.outer_.Transform(statementNode.FieldGet_value()))
}

// Java method: TransformDispatcher.processTryStatement(com.google.javascript.jscomp.parsing.parser.trees.TryStatementTree statementNode)
func (this_ *IRFactory_TransformDispatcher) ProcessTryStatement(statementNode *pkg_parser.TryStatementTree) pkg_rhino.Node {
	_ = statementNode
	var node pkg_rhino.Node = this_.outer_.NewNode_2(pkg_rhino.StaticGet_Token_TRY(), this_.outer_.TransformBlock(statementNode.FieldGet_body()))
	_ = node
	var block pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_BLOCK())
	_ = block
	var lineSet bool = false
	_ = lineSet
	var cc pkg_parser.ParseTree = statementNode.FieldGet_catchBlock()
	_ = cc
	if cc != nil && !cc.IsJavaNil_() {
		this_.outer_.SetSourceInfo_2_Node_ParseTree(block, cc)
		lineSet = true
		block.AddChildToBack(this_.outer_.Transform(cc))
	}
	node.AddChildToBack(block)
	var finallyBlock pkg_parser.ParseTree = statementNode.FieldGet_finallyBlock()
	_ = finallyBlock
	if finallyBlock != nil && !finallyBlock.IsJavaNil_() {
		node.AddChildToBack(this_.outer_.TransformBlock(finallyBlock))
	}
	if !(lineSet) && (finallyBlock != nil && !finallyBlock.IsJavaNil_()) {
		this_.outer_.SetSourceInfo_2_Node_ParseTree(block, finallyBlock)
	}
	return node
}

// Java method: TransformDispatcher.processCatchClause(com.google.javascript.jscomp.parsing.parser.trees.CatchTree clauseNode)
func (this_ *IRFactory_TransformDispatcher) ProcessCatchClause(clauseNode *pkg_parser.CatchTree) pkg_rhino.Node {
	_ = clauseNode
	if clauseNode.FieldGet_exception().GetParseTreeBase().FieldGet_type_() == pkg_parser.StaticGet_ParseTreeType_EMPTY_STATEMENT() {
		this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(clauseNode), pkg_parser.StaticGet_FeatureSet_Feature_OPTIONAL_CATCH_BINDING())
	}
	return this_.outer_.NewNode_3(pkg_rhino.StaticGet_Token_CATCH(), this_.outer_.Transform(clauseNode.FieldGet_exception()), this_.outer_.TransformBlock(clauseNode.FieldGet_catchBody()))
}

// Java method: TransformDispatcher.processFinally(com.google.javascript.jscomp.parsing.parser.trees.FinallyTree finallyNode)
func (this_ *IRFactory_TransformDispatcher) ProcessFinally(finallyNode *pkg_parser.FinallyTree) pkg_rhino.Node {
	_ = finallyNode
	return this_.outer_.TransformBlock(finallyNode.FieldGet_block())
}

// Java method: TransformDispatcher.processUnaryExpression(com.google.javascript.jscomp.parsing.parser.trees.UnaryExpressionTree exprNode)
func (this_ *IRFactory_TransformDispatcher) ProcessUnaryExpression(exprNode *pkg_parser.UnaryExpressionTree) pkg_rhino.Node {
	_ = exprNode
	var type_ pkg_rhino.Token = IRFactory_TransformUnaryTokenType(exprNode.FieldGet_operator().FieldGet_type_())
	_ = type_
	var operand pkg_rhino.Node = this_.outer_.Transform(exprNode.FieldGet_operand())
	_ = operand
	switch (type_).Ordinal() {
	case int32(24) /* Token.DELPROP */ :
		if !((((operand.IsGetProp() || operand.IsGetElem()) || operand.IsName()) || operand.IsOptChainGetProp()) || operand.IsOptChainGetElem()) {
			var msg string = "Invalid delete operand. Only properties can be deleted."
			_ = msg
			this_.outer_.errorReporter.M_Error(msg, this_.outer_.sourceName, operand.GetLineno(), int32(0))
		}
	case int32(21) /* Token.POS */ :
		if operand.IsBigInt() {
			this_.outer_.errorReporter.M_Error("Cannot convert a BigInt value to a number", this_.outer_.sourceName, operand.GetLineno(), int32(0))
		}
	default:
	}
	return this_.outer_.NewNode_2(type_, operand)
}

// Java method: TransformDispatcher.processUpdateExpression(com.google.javascript.jscomp.parsing.parser.trees.UpdateExpressionTree updateExpr)
func (this_ *IRFactory_TransformDispatcher) ProcessUpdateExpression(updateExpr *pkg_parser.UpdateExpressionTree) pkg_rhino.Node {
	_ = updateExpr
	var type_ pkg_rhino.Token = IRFactory_TransformUpdateTokenType(updateExpr.FieldGet_operator().FieldGet_type_())
	_ = type_
	var operand pkg_rhino.Node = this_.outer_.Transform(updateExpr.FieldGet_operand())
	_ = operand
	return this_.priv_TransformDispatcher_createUpdateNode_6976(type_, (updateExpr.FieldGet_operatorPosition() == pkg_parser.StaticGet_UpdateExpressionTree_OperatorPosition_POSTFIX()), operand)
}

// Java method: TransformDispatcher.createUpdateNode(com.google.javascript.rhino.Token type, boolean postfix, com.google.javascript.rhino.Node operand)
func (this_ *IRFactory_TransformDispatcher) priv_TransformDispatcher_createUpdateNode_6976(type_ pkg_rhino.Token, postfix bool, operand pkg_rhino.Node) pkg_rhino.Node {
	_ = type_
	_ = postfix
	_ = operand
	var assignTarget pkg_rhino.Node
	_ = assignTarget
	if operand.IsCast() {
		assignTarget = operand.GetFirstChild()
	} else {
		assignTarget = operand
	}
	if !(assignTarget.IsValidAssignmentTarget()) {
		this_.outer_.errorReporter.M_Error(j2gort.Ext_java_lang_String_Format_2("Invalid %s %s operand.", []any{(func() any {
			if postfix {
				return "postfix"
			}
			return "prefix"
		}()), (func() any {
			if type_ == pkg_rhino.StaticGet_Token_INC() {
				return "increment"
			}
			return "decrement"
		}())}), this_.outer_.sourceName, operand.GetLineno(), operand.GetCharno())
	}
	var node pkg_rhino.Node = this_.outer_.NewNode_2(type_, operand)
	_ = node
	node.PutBooleanProp(pkg_rhino.StaticGet_Node_INCRDECR_PROP(), postfix)
	return node
}

// Java method: TransformDispatcher.processVariableStatement(com.google.javascript.jscomp.parsing.parser.trees.VariableStatementTree stmt)
func (this_ *IRFactory_TransformDispatcher) ProcessVariableStatement(stmt *pkg_parser.VariableStatementTree) pkg_rhino.Node {
	_ = stmt
	return this_.outer_.transformDispatcher.Process(pkg_parser.ToIface_ParseTree(stmt.FieldGet_declarations()))
}

// Java method: TransformDispatcher.processVariableDeclarationList(com.google.javascript.jscomp.parsing.parser.trees.VariableDeclarationListTree decl)
func (this_ *IRFactory_TransformDispatcher) ProcessVariableDeclarationList(decl *pkg_parser.VariableDeclarationListTree) pkg_rhino.Node {
	_ = decl
	var declType pkg_rhino.Token = pkg_rhino.Token(0)
	_ = declType
	switch (decl.FieldGet_declarationType()).Ordinal() {
	case int32(30) /* TokenType.CONST */ :
		this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(decl), pkg_parser.StaticGet_FeatureSet_Feature_CONST_DECLARATIONS())
		declType = pkg_rhino.StaticGet_Token_CONST()
	case int32(38) /* TokenType.LET */ :
		this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(decl), pkg_parser.StaticGet_FeatureSet_Feature_LET_DECLARATIONS())
		declType = pkg_rhino.StaticGet_Token_LET()
	case int32(25) /* TokenType.VAR */ :
		declType = pkg_rhino.StaticGet_Token_VAR()
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException")))
	}
	var node pkg_rhino.Node = this_.outer_.NewNode_1(declType)
	_ = node
	{

		for _, child := range j2gort.CheckNotNullSlice(j2gort.CheckNotNullPtr(decl).
			Slice_FieldGet_declarations()) {

			_ = child
			node.AddChildToBack(this_.outer_.TransformNodeWithInlineComments(pkg_parser.ToIface_ParseTree(child)))
		}
	}
	this_.AttachPossibleTrailingComment(node, decl.GetEnd())
	return node
}

// Java method: TransformDispatcher.processVariableDeclaration(com.google.javascript.jscomp.parsing.parser.trees.VariableDeclarationTree decl)
func (this_ *IRFactory_TransformDispatcher) ProcessVariableDeclaration(decl *pkg_parser.VariableDeclarationTree) pkg_rhino.Node {
	_ = decl
	var node pkg_rhino.Node = this_.outer_.TransformNodeWithInlineComments(decl.FieldGet_lvalue())
	_ = node
	var lhs pkg_rhino.Node
	_ = lhs
	if node.IsDestructuringPattern() {
		lhs = this_.outer_.NewNode_2(pkg_rhino.StaticGet_Token_DESTRUCTURING_LHS(), node)
	} else {
		lhs = node
	}
	if !(j2gort.IfaceEquals[pkg_parser.ParseTree](decl.FieldGet_initializer(), nil)) {
		var initializer pkg_rhino.Node = this_.outer_.Transform(decl.FieldGet_initializer())
		_ = initializer
		lhs.AddChildToBack(initializer)
		this_.outer_.SetLength(lhs, decl.FieldGet_location().FieldGet_start(), decl.FieldGet_location().FieldGet_end())
	}
	return lhs
}

// Java method: TransformDispatcher.processWhileLoop(com.google.javascript.jscomp.parsing.parser.trees.WhileStatementTree stmt)
func (this_ *IRFactory_TransformDispatcher) ProcessWhileLoop(stmt *pkg_parser.WhileStatementTree) pkg_rhino.Node {
	_ = stmt
	return this_.outer_.NewNode_3(pkg_rhino.StaticGet_Token_WHILE(), this_.outer_.Transform(stmt.FieldGet_condition()), this_.outer_.TransformBlock(stmt.FieldGet_body()))
}

// Java method: TransformDispatcher.processWithStatement(com.google.javascript.jscomp.parsing.parser.trees.WithStatementTree stmt)
func (this_ *IRFactory_TransformDispatcher) ProcessWithStatement(stmt *pkg_parser.WithStatementTree) pkg_rhino.Node {
	_ = stmt
	return this_.outer_.NewNode_3(pkg_rhino.StaticGet_Token_WITH(), this_.outer_.Transform(stmt.FieldGet_expression()), this_.outer_.TransformBlock(stmt.FieldGet_body()))
}

// Java method: TransformDispatcher.processMissingExpression(com.google.javascript.jscomp.parsing.parser.trees.MissingPrimaryExpressionTree unused)
func (this_ *IRFactory_TransformDispatcher) ProcessMissingExpression(unused *pkg_parser.MissingPrimaryExpressionTree) pkg_rhino.Node {
	_ = unused
	return this_.priv_TransformDispatcher_createMissingExpressionNode_6984()
}

// Java method: TransformDispatcher.createMissingNameNode()
func (this_ *IRFactory_TransformDispatcher) priv_TransformDispatcher_createMissingNameNode_6983() pkg_rhino.Node {
	return this_.outer_.NewStringNode_2(pkg_rhino.StaticGet_Token_NAME(), "__missing_name__")
}

// Java method: TransformDispatcher.createMissingExpressionNode()
func (this_ *IRFactory_TransformDispatcher) priv_TransformDispatcher_createMissingExpressionNode_6984() pkg_rhino.Node {
	return this_.outer_.NewStringNode_2(pkg_rhino.StaticGet_Token_NAME(), "__missing_expression__")
}

// Java method: TransformDispatcher.processIllegalToken(com.google.javascript.jscomp.parsing.parser.trees.ParseTree node)
func (this_ *IRFactory_TransformDispatcher) ProcessIllegalToken(node pkg_parser.ParseTree) pkg_rhino.Node {
	_ = node
	this_.outer_.errorReporter.M_Error(("Unsupported syntax: " + j2gort.StrVal(node.GetParseTreeBase().FieldGet_type_())), this_.outer_.sourceName, IRFactory_Lineno_1_ParseTree(node), int32(0))
	return this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_EMPTY())
}

// Java method: TransformDispatcher.maybeReportGetter(com.google.javascript.jscomp.parsing.parser.trees.ParseTree node)
func (this_ *IRFactory_TransformDispatcher) MaybeReportGetter(node pkg_parser.ParseTree) bool {
	_ = node
	this_.outer_.features = this_.outer_.features.With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_GETTER())
	if this_.outer_.config.LanguageMode() == StaticGet_Config_LanguageMode_ECMASCRIPT3() {
		this_.outer_.errorReporter.M_Error("getters are not supported in older versions of JavaScript. If you are targeting newer versions of JavaScript, set the appropriate language_in option.", this_.outer_.sourceName, IRFactory_Lineno_1_ParseTree(node), int32(0))
		return true
	}
	return false
}

// Java method: TransformDispatcher.maybeReportSetter(com.google.javascript.jscomp.parsing.parser.trees.ParseTree node)
func (this_ *IRFactory_TransformDispatcher) MaybeReportSetter(node pkg_parser.ParseTree) bool {
	_ = node
	this_.outer_.features = this_.outer_.features.With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_SETTER())
	if this_.outer_.config.LanguageMode() == StaticGet_Config_LanguageMode_ECMASCRIPT3() {
		this_.outer_.errorReporter.M_Error("setters are not supported in older versions of JavaScript. If you are targeting newer versions of JavaScript, set the appropriate language_in option.", this_.outer_.sourceName, IRFactory_Lineno_1_ParseTree(node), int32(0))
		return true
	}
	return false
}

// Java method: TransformDispatcher.processBooleanLiteral(com.google.javascript.jscomp.parsing.parser.trees.LiteralExpressionTree literal)
func (this_ *IRFactory_TransformDispatcher) ProcessBooleanLiteral(literal *pkg_parser.LiteralExpressionTree) pkg_rhino.Node {
	_ = literal
	return this_.outer_.NewNode_1(IRFactory_TransformBooleanTokenType(literal.FieldGet_literalToken().FieldGet_type_()))
}

// Java method: TransformDispatcher.processNullLiteral(com.google.javascript.jscomp.parsing.parser.trees.LiteralExpressionTree unused)
func (this_ *IRFactory_TransformDispatcher) ProcessNullLiteral(unused *pkg_parser.LiteralExpressionTree) pkg_rhino.Node {
	_ = unused
	return this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_NULL())
}

// Java method: TransformDispatcher.processNull(com.google.javascript.jscomp.parsing.parser.trees.NullTree unused)
func (this_ *IRFactory_TransformDispatcher) ProcessNull(unused *pkg_parser.NullTree) pkg_rhino.Node {
	_ = unused
	return this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_EMPTY())
}

// Java method: TransformDispatcher.processCommaExpression(com.google.javascript.jscomp.parsing.parser.trees.CommaExpressionTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessCommaExpression(tree *pkg_parser.CommaExpressionTree) pkg_rhino.Node {
	_ = tree
	var root pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_COMMA())
	_ = root
	var start pkg_parser.SourcePosition = j2gort.AsIface[pkg_parser.ParseTree](tree.FieldGet_expressions().Get(int32(0))).GetParseTreeBase().FieldGet_location().FieldGet_start()
	_ = start
	var end pkg_parser.SourcePosition = j2gort.AsIface[pkg_parser.ParseTree](tree.FieldGet_expressions().Get(int32(1))).GetParseTreeBase().FieldGet_location().FieldGet_end()
	_ = end
	this_.outer_.SetSourceInfo_3(root, start, end)
	{
		cur_69 := j2gort.NewIterCursor(j2gort.CheckNotNull(tree.FieldGet_expressions()))
		for cur_69.HasNext() {
			rawElem_70 := cur_69.Next()
			expr := j2gort.AsIface[pkg_parser.ParseTree](rawElem_70)
			_ = expr
			var count int32 = root.GetChildCount()
			_ = count
			if count < int32(2) {
				root.AddChildToBack(this_.outer_.Transform(expr))
			} else {
				end = expr.GetParseTreeBase().FieldGet_location().FieldGet_end()
				root = this_.outer_.NewNode_3(pkg_rhino.StaticGet_Token_COMMA(), root, this_.outer_.Transform(expr))
				this_.outer_.SetSourceInfo_3(root, start, end)
			}
		}
	}
	return root
}

// Java method: TransformDispatcher.processClassDeclaration(com.google.javascript.jscomp.parsing.parser.trees.ClassDeclarationTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessClassDeclaration(tree *pkg_parser.ClassDeclarationTree) pkg_rhino.Node {
	_ = tree
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(tree), pkg_parser.StaticGet_FeatureSet_Feature_CLASSES())
	var name pkg_rhino.Node = this_.TransformOrEmpty_2_Token_ParseTree(tree.FieldGet_name(), pkg_parser.ToIface_ParseTree(tree))
	_ = name
	var superClass pkg_rhino.Node = this_.TransformOrEmpty_2_ParseTree_ParseTree(tree.FieldGet_superClass(), pkg_parser.ToIface_ParseTree(tree))
	_ = superClass
	var classMembers pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_CLASS_MEMBERS())
	_ = classMembers
	this_.outer_.SetSourceInfo_2_Node_ParseTree(classMembers, pkg_parser.ToIface_ParseTree(tree))
	this_.classScope.Enter()
	func() {
		defer func() {
			this_.classScope.Exit()
		}()
		var hasConstructor bool = false
		_ = hasConstructor
		{
			cur_71 := j2gort.NewIterCursor(j2gort.CheckNotNull(tree.FieldGet_elements()))
			for cur_71.HasNext() {
				rawElem_72 := cur_71.Next()
				child := j2gort.AsIface[pkg_parser.ParseTree](rawElem_72)
				_ = child
				switch (child.GetParseTreeBase().FieldGet_type_()).Ordinal() {
				case int32(34) /* ParseTreeType.COMPUTED_PROPERTY_GETTER */, int32(36) /* ParseTreeType.COMPUTED_PROPERTY_SETTER */, int32(30) /* ParseTreeType.GET_ACCESSOR */, int32(31) /* ParseTreeType.SET_ACCESSOR */ :
					this_.outer_.features = this_.outer_.features.With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_CLASS_GETTER_SETTER())
				case int32(2) /* ParseTreeType.BLOCK */ :
					this_.outer_.features = this_.outer_.features.With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_CLASS_STATIC_BLOCK())
				default:
				}
				var childIsCtor bool = this_.priv_TransformDispatcher_validateClassConstructorMember_6995(child)
				_ = childIsCtor
				if childIsCtor {
					if hasConstructor {
						this_.outer_.errorReporter.M_Error("Class may have only one constructor.", this_.outer_.sourceName, IRFactory_Lineno_1_ParseTree(child), IRFactory_Charno_1_ParseTree(child))
					}
					hasConstructor = true
				}
				classMembers.AddChildToBack(this_.outer_.Transform(child))
			}
		}
	}()
	if !(this_.classScope.InScope()) {
		this_.priv_TransformDispatcher_validatePrivatePropertyUsage_6993(j2gort.Ext_java_util_Collections_UnmodifiableSet(this_.priv_TransformDispatcher_getPrivatePropsAndValidateUniqueness_6994(classMembers)), classMembers)
	}
	var classNode pkg_rhino.Node = this_.outer_.NewNode_4(pkg_rhino.StaticGet_Token_CLASS(), name, superClass, classMembers)
	_ = classNode
	this_.AttachPossibleTrailingComment(classNode, tree.GetEnd())
	return classNode
}

// Java method: TransformDispatcher.validatePrivatePropertyUsage(java.util.Set<java.lang.String> privatePropNames, com.google.javascript.rhino.Node classMembers)
func (this_ *IRFactory_TransformDispatcher) priv_TransformDispatcher_validatePrivatePropertyUsage_6993(privatePropNames j2gort.Ext_java_util_Set, classMembers pkg_rhino.Node) {
	_ = privatePropNames
	_ = classMembers
	if !(classMembers.IsClassMembers()) {
		j2gort.CheckState(false)
	}
	if !(classMembers.HasChildren()) {
		return
	}
	var visitStack []pkg_rhino.Node = make([]pkg_rhino.Node, int32(32))
	_ = visitStack
	var stackTop int32 = int32(0)
	_ = stackTop
	{
		a_73 := visitStack
		var i_74 int32
		{

			old_75 := stackTop
			nw_76 := old_75 + 1
			stackTop = nw_76
			i_74 = old_75
		}
		a_73[i_74] = classMembers.GetFirstChild()
	}
	for stackTop > int32(0) {
		var node pkg_rhino.Node = visitStack[func() int32 { old_77 := stackTop; nw_78 := old_77 - 1; stackTop = nw_78; return nw_78 }()]
		_ = node
		{
			a_79 := visitStack
			i_80 := stackTop
			a_79[i_80] = pkg_rhino.Node(0)
		}
		if (node.IsPrivateIdentifier() && ((node.IsGetProp() || node.IsOptChainGetProp()) || ((node.IsName() && node.HasParent()) && node.GetParent().IsIn()))) && !(j2gort.SetContainsStr(privatePropNames, node.GetString())) {
			var propType string
			_ = propType
			if node.GetParent().IsCall() {
				propType = "methods"
			} else {
				propType = "fields"
			}
			this_.outer_.errorReporter.M_Error(("Private " + j2gort.StrValString(propType) + " must be declared in an enclosing class"), this_.outer_.sourceName, node.GetLineno(), node.GetCharno())
		}
		if node.IsDelProp() {
			var firstChild pkg_rhino.Node = j2gort.CheckNotNullNode(node.GetFirstChild())
			_ = firstChild
			if firstChild.IsPrivateIdentifier() && (firstChild.IsGetProp() || firstChild.IsOptChainGetProp()) {
				this_.outer_.errorReporter.M_Error("Private fields cannot be deleted", this_.outer_.sourceName, node.GetLineno(), node.GetCharno())
			}
		}
		if node.IsClassMembers() {
			var innerClassPrivatePropNames j2gort.Ext_java_util_Set = this_.priv_TransformDispatcher_getPrivatePropsAndValidateUniqueness_6994(node)
			_ = innerClassPrivatePropNames
			innerClassPrivatePropNames.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](privatePropNames))
			this_.priv_TransformDispatcher_validatePrivatePropertyUsage_6993(j2gort.Ext_java_util_Collections_UnmodifiableSet(innerClassPrivatePropNames), node)
		} else {
			var next pkg_rhino.Node = node.GetNext()
			_ = next
			if next != 0 {
				if stackTop == int32(len(visitStack)) {
					{

						arrSrc_81 := visitStack
						arrLen_82 := (stackTop * int32(2))
						j2gort.CheckNotNullSlice(arrSrc_81)
						if arrLen_82 < 0 {
							panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
						}
						arrOut_83 := make([]pkg_rhino.Node, arrLen_82)
						copy(arrOut_83, arrSrc_81)
						visitStack = arrOut_83
					}

				}
				{
					a_84 := visitStack
					var i_85 int32
					{

						old_86 := stackTop
						nw_87 := old_86 + 1
						stackTop = nw_87
						i_85 = old_86
					}
					a_84[i_85] = next
				}
			}
			if node.HasChildren() {
				if stackTop == int32(len(visitStack)) {
					{

						arrSrc_88 := visitStack
						arrLen_89 := (stackTop * int32(2))
						j2gort.CheckNotNullSlice(arrSrc_88)
						if arrLen_89 < 0 {
							panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
						}
						arrOut_90 := make([]pkg_rhino.Node, arrLen_89)
						copy(arrOut_90, arrSrc_88)
						visitStack = arrOut_90
					}

				}
				{
					a_91 := visitStack
					var i_92 int32
					{

						old_93 := stackTop
						nw_94 := old_93 + 1
						stackTop = nw_94
						i_92 = old_93
					}
					a_91[i_92] = node.GetFirstChild()
				}
			}
		}
	}
}

// Java method: TransformDispatcher.getPrivatePropsAndValidateUniqueness(com.google.javascript.rhino.Node classMembers)
func (this_ *IRFactory_TransformDispatcher) priv_TransformDispatcher_getPrivatePropsAndValidateUniqueness_6994(classMembers pkg_rhino.Node) j2gort.Ext_java_util_Set {
	_ = classMembers
	var privateGetterNames j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_TreeSet()
	_ = privateGetterNames
	var privateSetterNames j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_TreeSet()
	_ = privateSetterNames
	var privateStaticGetterNames j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_TreeSet()
	_ = privateStaticGetterNames
	var privateStaticSetterNames j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_TreeSet()
	_ = privateStaticSetterNames
	var privateFieldAndMethodNames j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_TreeSet()
	_ = privateFieldAndMethodNames
	var privatePropNames j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_TreeSet()
	_ = privatePropNames
	{
		var curNode pkg_rhino.Node = classMembers.GetFirstChild()
		_ = curNode
		for ; curNode != 0; curNode = curNode.GetNext() {
			if !(curNode.IsPrivateIdentifier()) {
				continue
			}
			var propName string = curNode.GetString()
			_ = propName
			if !(((curNode.IsGetterDef() || curNode.IsSetterDef()) || curNode.IsMemberFieldDef()) || curNode.IsMemberFunctionDef()) {
				j2gort.CheckStateFmt(false, "Private property '%s' has an unsupported token type: %s", j2gort.BoxString(propName), j2gort.BoxNullableVal(curNode.GetToken()))
			}
			var isStatic bool = curNode.IsStaticMember()
			_ = isStatic
			var alreadyDeclared bool = j2gort.SetContainsStr(privatePropNames, propName)
			_ = alreadyDeclared
			switch curNode.TokenSwitchOrd_() {
			case int32(95) /* Token.GETTER_DEF */ :
				if alreadyDeclared {
					var alreadyDeclaredIsOtherThanItsSetter bool = (((j2gort.SetContainsStr(privateFieldAndMethodNames, propName) || j2gort.SetContainsStr(privateGetterNames, propName)) || j2gort.SetContainsStr(privateStaticGetterNames, propName)) || (func() bool {
						if isStatic {
							return j2gort.SetContainsStr(privateSetterNames, propName)
						}
						return j2gort.SetContainsStr(privateStaticSetterNames, propName)
					}()))
					_ = alreadyDeclaredIsOtherThanItsSetter
					if !(alreadyDeclaredIsOtherThanItsSetter || (func() bool {
						if isStatic {
							return j2gort.SetContainsStr(privateStaticSetterNames, propName)
						}
						return j2gort.SetContainsStr(privateSetterNames, propName)
					}())) {
						j2gort.CheckState(false)
					}
					alreadyDeclared = alreadyDeclaredIsOtherThanItsSetter
				}
				if isStatic {
					j2gort.SetAddStr(privateStaticGetterNames, propName)
				} else {
					j2gort.SetAddStr(privateGetterNames, propName)
				}
			case int32(96) /* Token.SETTER_DEF */ :
				if alreadyDeclared {
					var alreadyDeclaredIsOtherThanItsGetter bool = (((j2gort.SetContainsStr(privateFieldAndMethodNames, propName) || j2gort.SetContainsStr(privateSetterNames, propName)) || j2gort.SetContainsStr(privateStaticSetterNames, propName)) || (func() bool {
						if isStatic {
							return j2gort.SetContainsStr(privateGetterNames, propName)
						}
						return j2gort.SetContainsStr(privateStaticGetterNames, propName)
					}()))
					_ = alreadyDeclaredIsOtherThanItsGetter
					if !(alreadyDeclaredIsOtherThanItsGetter || (func() bool {
						if isStatic {
							return j2gort.SetContainsStr(privateStaticGetterNames, propName)
						}
						return j2gort.SetContainsStr(privateGetterNames, propName)
					}())) {
						j2gort.CheckState(false)
					}
					alreadyDeclared = alreadyDeclaredIsOtherThanItsGetter
				}
				if curNode.IsStaticMember() {
					j2gort.SetAddStr(privateStaticSetterNames, propName)
				} else {
					j2gort.SetAddStr(privateSetterNames, propName)
				}
			case int32(108) /* Token.MEMBER_FIELD_DEF */, int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
				j2gort.SetAddStr(privateFieldAndMethodNames, propName)
			default:
			}
			if alreadyDeclared {
				this_.outer_.errorReporter.M_Error(("Identifier '" + j2gort.StrValString(propName) + "' has already been declared"), this_.outer_.sourceName, curNode.GetLineno(), curNode.GetCharno())
			}
			j2gort.SetAddStr(privatePropNames, propName)
		}
	}
	return privatePropNames
}

// Java method: TransformDispatcher.validateClassConstructorMember(com.google.javascript.jscomp.parsing.parser.trees.ParseTree member)
func (this_ *IRFactory_TransformDispatcher) priv_TransformDispatcher_validateClassConstructorMember_6995(member pkg_parser.ParseTree) bool {
	_ = member
	var memberName pkg_parser.Token = (pkg_parser.Token{})
	_ = memberName
	var isStatic bool = false
	_ = isStatic
	var hasIllegalModifier bool = false
	_ = hasIllegalModifier
	base_member_ := member.GetParseTreeBase()
	switch (base_member_.FieldGet_type_()).Ordinal() {
	case int32(30) /* ParseTreeType.GET_ACCESSOR */ :
		var getter *pkg_parser.GetAccessorTree = base_member_.AsGetAccessor()
		_ = getter
		memberName = getter.FieldGet_propertyName()
		isStatic = getter.FieldGet_isStatic()
		hasIllegalModifier = true
	case int32(31) /* ParseTreeType.SET_ACCESSOR */ :
		var setter *pkg_parser.SetAccessorTree = base_member_.AsSetAccessor()
		_ = setter
		memberName = setter.FieldGet_propertyName()
		isStatic = setter.FieldGet_isStatic()
		hasIllegalModifier = true
	case int32(1) /* ParseTreeType.FUNCTION_DECLARATION */ :
		var method *pkg_parser.FunctionDeclarationTree = base_member_.AsFunctionDeclaration()
		_ = method
		memberName = method.FieldGet_name()
		isStatic = method.FieldGet_isStatic()
		hasIllegalModifier = (method.FieldGet_isGenerator() || method.FieldGet_isAsync())
	default:
		return false
	}
	if isStatic {
		return false
	}
	{
		var iifeCond_1 bool
		{

			re_95 := memberName.FieldGet_type_()
			ae_96 := pkg_parser.StaticGet_TokenType_IDENTIFIER()
			if re_95 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_95 == ae_96
		}
		if !(iifeCond_1) || !(memberName.ValueEquals("constructor")) {
			return false
		}
	}
	if hasIllegalModifier {
		this_.outer_.errorReporter.M_Error("Class constructor may not be getter, setter, async, or generator.", this_.outer_.sourceName, IRFactory_Lineno_1_ParseTree(member), IRFactory_Charno_1_ParseTree(member))
		return false
	}
	return true
}

// Java method: TransformDispatcher.processSuper(com.google.javascript.jscomp.parsing.parser.trees.SuperExpressionTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessSuper(tree *pkg_parser.SuperExpressionTree) pkg_rhino.Node {
	_ = tree
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(tree), pkg_parser.StaticGet_FeatureSet_Feature_SUPER())
	return this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_SUPER())
}

// Java method: TransformDispatcher.processNewTarget(com.google.javascript.jscomp.parsing.parser.trees.NewTargetExpressionTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessNewTarget(tree *pkg_parser.NewTargetExpressionTree) pkg_rhino.Node {
	_ = tree
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(tree), pkg_parser.StaticGet_FeatureSet_Feature_NEW_TARGET())
	return this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_NEW_TARGET())
}

// Java method: TransformDispatcher.processYield(com.google.javascript.jscomp.parsing.parser.trees.YieldExpressionTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessYield(tree *pkg_parser.YieldExpressionTree) pkg_rhino.Node {
	_ = tree
	var yield pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_YIELD())
	_ = yield
	if !(j2gort.IfaceEquals[pkg_parser.ParseTree](tree.FieldGet_expression(), nil)) {
		yield.AddChildToBack(this_.outer_.Transform(tree.FieldGet_expression()))
	}
	yield.SetYieldAll(tree.FieldGet_isYieldAll())
	return yield
}

// Java method: TransformDispatcher.processAwait(com.google.javascript.jscomp.parsing.parser.trees.AwaitExpressionTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessAwait(tree *pkg_parser.AwaitExpressionTree) pkg_rhino.Node {
	_ = tree
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(tree), pkg_parser.StaticGet_FeatureSet_Feature_ASYNC_FUNCTIONS())
	var await pkg_rhino.Node = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_AWAIT())
	_ = await
	await.AddChildToBack(this_.outer_.Transform(tree.FieldGet_expression()))
	return await
}

// Java method: TransformDispatcher.processExportDecl(com.google.javascript.jscomp.parsing.parser.trees.ExportDeclarationTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessExportDecl(tree *pkg_parser.ExportDeclarationTree) pkg_rhino.Node {
	_ = tree
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(tree), pkg_parser.StaticGet_FeatureSet_Feature_MODULES())
	var decls pkg_rhino.Node = pkg_rhino.Node(0)
	_ = decls
	if tree.FieldGet_isExportAll() {
		if !(j2gort.IfaceEquals[pkg_parser.ParseTree](tree.FieldGet_declaration(), nil) && (tree.FieldGet_exportSpecifierList() == nil)) {
			j2gort.CheckState(false)
		}
	} else if !(j2gort.IfaceEquals[pkg_parser.ParseTree](tree.FieldGet_declaration(), nil)) {
		if !(tree.FieldGet_exportSpecifierList() == nil) {
			j2gort.CheckState(false)
		}
		decls = this_.outer_.Transform(tree.FieldGet_declaration())
	} else {
		decls = this_.priv_TransformDispatcher_transformList_7007(pkg_rhino.StaticGet_Token_EXPORT_SPECS(), tree.FieldGet_exportSpecifierList())
	}
	if decls == 0 {
		decls = this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_EMPTY())
	}
	this_.outer_.SetSourceInfo_2_Node_ParseTree(decls, pkg_parser.ToIface_ParseTree(tree))
	var export pkg_rhino.Node = this_.outer_.NewNode_2(pkg_rhino.StaticGet_Token_EXPORT(), decls)
	_ = export
	if !(tree.FieldGet_from()).IsJavaNil_() {
		var from pkg_rhino.Node = this_.ProcessString(tree.FieldGet_from())
		_ = from
		export.AddChildToBack(from)
	}
	export.PutBooleanProp(pkg_rhino.StaticGet_Node_EXPORT_ALL_FROM(), tree.FieldGet_isExportAll())
	export.PutBooleanProp(pkg_rhino.StaticGet_Node_EXPORT_DEFAULT(), tree.FieldGet_isDefault())
	return export
}

// Java method: TransformDispatcher.processExportSpec(com.google.javascript.jscomp.parsing.parser.trees.ExportSpecifierTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessExportSpec(tree *pkg_parser.ExportSpecifierTree) pkg_rhino.Node {
	_ = tree
	var importedName pkg_rhino.Node = this_.ProcessName_3(tree.FieldGet_importedName(), pkg_rhino.StaticGet_Token_NAME(), StaticGet_IRFactory_TransformDispatcher_IdentifierType_STANDARD())
	_ = importedName
	var exportSpec pkg_rhino.Node = this_.outer_.NewNode_2(pkg_rhino.StaticGet_Token_EXPORT_SPEC(), importedName)
	_ = exportSpec
	if (tree.FieldGet_destinationName()).IsJavaNil_() {
		exportSpec.SetShorthandProperty(true)
		exportSpec.AddChildToBack(importedName.CloneTree_0())
	} else {
		var destinationName pkg_rhino.Node = this_.ProcessName_3(tree.FieldGet_destinationName(), pkg_rhino.StaticGet_Token_NAME(), StaticGet_IRFactory_TransformDispatcher_IdentifierType_STANDARD())
		_ = destinationName
		exportSpec.AddChildToBack(destinationName)
	}
	return exportSpec
}

// Java method: TransformDispatcher.processImportDecl(com.google.javascript.jscomp.parsing.parser.trees.ImportDeclarationTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessImportDecl(tree *pkg_parser.ImportDeclarationTree) pkg_rhino.Node {
	_ = tree
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(tree), pkg_parser.StaticGet_FeatureSet_Feature_MODULES())
	var firstChild pkg_rhino.Node = this_.TransformOrEmpty_2_Token_ParseTree(tree.FieldGet_defaultBindingIdentifier(), pkg_parser.ToIface_ParseTree(tree))
	_ = firstChild
	var secondChild pkg_rhino.Node = pkg_rhino.Node(0)
	_ = secondChild
	if (tree.FieldGet_nameSpaceImportIdentifier()).IsJavaNil_() {
		secondChild = this_.priv_TransformDispatcher_transformListOrEmpty_7008(pkg_rhino.StaticGet_Token_IMPORT_SPECS(), tree.FieldGet_importSpecifierList())
		this_.outer_.SetSourceInfo_2_Node_ParseTree(secondChild, pkg_parser.ToIface_ParseTree(tree))
	} else {
		secondChild = this_.NewStringNodeFromIdentifier(pkg_rhino.StaticGet_Token_IMPORT_STAR(), StaticGet_IRFactory_TransformDispatcher_IdentifierType_STANDARD(), tree.FieldGet_nameSpaceImportIdentifier())
		this_.outer_.SetSourceInfo_2_Node_Token(secondChild, tree.FieldGet_nameSpaceImportIdentifier())
	}
	var thirdChild pkg_rhino.Node = this_.ProcessString(tree.FieldGet_moduleSpecifier())
	_ = thirdChild
	return this_.outer_.NewNode_4(pkg_rhino.StaticGet_Token_IMPORT(), firstChild, secondChild, thirdChild)
}

// Java method: TransformDispatcher.processImportSpec(com.google.javascript.jscomp.parsing.parser.trees.ImportSpecifierTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessImportSpec(tree *pkg_parser.ImportSpecifierTree) pkg_rhino.Node {
	_ = tree
	var importedName pkg_rhino.Node = this_.ProcessName_3(tree.FieldGet_importedName(), pkg_rhino.StaticGet_Token_NAME(), StaticGet_IRFactory_TransformDispatcher_IdentifierType_STANDARD())
	_ = importedName
	var importSpec pkg_rhino.Node = this_.outer_.NewNode_2(pkg_rhino.StaticGet_Token_IMPORT_SPEC(), importedName)
	_ = importSpec
	if (tree.FieldGet_destinationName()).IsJavaNil_() {
		importSpec.SetShorthandProperty(true)
		importSpec.AddChildToBack(importedName.CloneTree_0())
	} else {
		importSpec.AddChildToBack(this_.ProcessName_3(tree.FieldGet_destinationName(), pkg_rhino.StaticGet_Token_NAME(), StaticGet_IRFactory_TransformDispatcher_IdentifierType_STANDARD()))
	}
	return importSpec
}

// Java method: TransformDispatcher.processDynamicImport(com.google.javascript.jscomp.parsing.parser.trees.DynamicImportTree dynamicImportNode)
func (this_ *IRFactory_TransformDispatcher) ProcessDynamicImport(dynamicImportNode *pkg_parser.DynamicImportTree) pkg_rhino.Node {
	_ = dynamicImportNode
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(dynamicImportNode), pkg_parser.StaticGet_FeatureSet_Feature_DYNAMIC_IMPORT())
	var argument pkg_rhino.Node = this_.outer_.Transform(dynamicImportNode.FieldGet_argument())
	_ = argument
	return this_.outer_.NewNode_2(pkg_rhino.StaticGet_Token_DYNAMIC_IMPORT(), argument)
}

// Java method: TransformDispatcher.processImportMeta(com.google.javascript.jscomp.parsing.parser.trees.ImportMetaExpressionTree tree)
func (this_ *IRFactory_TransformDispatcher) ProcessImportMeta(tree *pkg_parser.ImportMetaExpressionTree) pkg_rhino.Node {
	_ = tree
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(tree), pkg_parser.StaticGet_FeatureSet_Feature_MODULES())
	this_.outer_.MaybeWarnForFeature_2_ParseTree_FeatureSet_Feature(pkg_parser.ToIface_ParseTree(tree), pkg_parser.StaticGet_FeatureSet_Feature_IMPORT_META())
	return this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_IMPORT_META())
}

// Java method: TransformDispatcher.checkParameters(com.google.javascript.jscomp.parsing.parser.trees.FormalParameterListTree tree)
func (this_ *IRFactory_TransformDispatcher) priv_TransformDispatcher_checkParameters_7006(tree *pkg_parser.FormalParameterListTree) bool {
	_ = tree
	var len_ int32 = tree.FieldGet_parameterCount()
	_ = len_
	var good bool = true
	_ = good
	{
		var i int32 = int32(0)
		_ = i
		for ; i < len_; i = i + 1 {
			var param pkg_parser.ParseTree = tree.GetParameter(i)
			_ = param
			if param.GetParseTreeBase().FieldGet_type_() == pkg_parser.StaticGet_ParseTreeType_ITER_REST() {
				if i != (len_ - int32(1)) {
					this_.outer_.errorReporter.M_Error("A rest parameter must be last in a parameter list.", this_.outer_.sourceName, IRFactory_Lineno_1_ParseTree(param), IRFactory_Charno_1_ParseTree(param))
					good = false
				} else if tree.FieldGet_hasTrailingComma() {
					this_.outer_.errorReporter.M_Error("A trailing comma must not follow a rest parameter.", this_.outer_.sourceName, IRFactory_Lineno_1_ParseTree(param), IRFactory_Charno_1_ParseTree(param))
					good = false
				}
			}
		}
	}
	return good
}

// Java method: TransformDispatcher.transformList(com.google.javascript.rhino.Token type, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.parsing.parser.trees.ParseTree> list)
func (this_ *IRFactory_TransformDispatcher) priv_TransformDispatcher_transformList_7007(type_ pkg_rhino.Token, list j2gort.Ext_com_google_common_collect_ImmutableList) pkg_rhino.Node {
	_ = type_
	_ = list
	var n pkg_rhino.Node = this_.outer_.NewNode_1(type_)
	_ = n
	{
		cur_97 := j2gort.NewIterCursor(j2gort.CheckNotNull(list))
		for cur_97.HasNext() {
			rawElem_98 := cur_97.Next()
			tree := j2gort.AsIface[pkg_parser.ParseTree](rawElem_98)
			_ = tree
			n.AddChildToBack(this_.outer_.Transform(tree))
		}
	}
	return n
}

// Java method: TransformDispatcher.transformListOrEmpty(com.google.javascript.rhino.Token type, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.parsing.parser.trees.ParseTree> list)
func (this_ *IRFactory_TransformDispatcher) priv_TransformDispatcher_transformListOrEmpty_7008(type_ pkg_rhino.Token, list j2gort.Ext_com_google_common_collect_ImmutableList) pkg_rhino.Node {
	_ = type_
	_ = list
	if (list == nil) || list.IsEmpty() {
		return this_.outer_.NewNode_1(pkg_rhino.StaticGet_Token_EMPTY())
	} else {
		return this_.priv_TransformDispatcher_transformList_7007(type_, list)
	}
}

// Java method: TransformDispatcher.processLiteralExpression(com.google.javascript.jscomp.parsing.parser.trees.LiteralExpressionTree expr)
func (this_ *IRFactory_TransformDispatcher) ProcessLiteralExpression(expr *pkg_parser.LiteralExpressionTree) pkg_rhino.Node {
	_ = expr
	switch (expr.FieldGet_literalToken().FieldGet_type_()).Ordinal() {
	case int32(104) /* TokenType.NUMBER */ :
		return this_.ProcessNumberLiteral(expr)
	case int32(105) /* TokenType.STRING */ :
		return this_.ProcessStringLiteral(expr)
	case int32(106) /* TokenType.BIGINT */ :
		return this_.ProcessBigIntLiteral(expr)
	case int32(103) /* TokenType.FALSE */, int32(102) /* TokenType.TRUE */ :
		return this_.ProcessBooleanLiteral(expr)
	case int32(101) /* TokenType.NULL */ :
		return this_.ProcessNullLiteral(expr)
	case int32(107) /* TokenType.REGULAR_EXPRESSION */ :
		return this_.ProcessRegExpLiteral(expr)
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(("Unexpected literal type: " + j2gort.StrVal(j2gort.GetObjectClass(j2gort.BoxNullableVal(expr.FieldGet_literalToken()))) + " type: " + j2gort.StrVal(expr.FieldGet_literalToken().FieldGet_type_()))))))
	}
	panic("unreachable switch expression")
}

// Java method: TransformDispatcher.process(com.google.javascript.jscomp.parsing.parser.trees.ParseTree node)
func (this_ *IRFactory_TransformDispatcher) Process(node pkg_parser.ParseTree) pkg_rhino.Node {
	_ = node
	base_node_ := node.GetParseTreeBase()
	switch (base_node_.FieldGet_type_()).Ordinal() {
	case int32(39) /* ParseTreeType.BINARY_OPERATOR */ :
		return this_.ProcessBinaryExpression(base_node_.AsBinaryOperator())
	case int32(28) /* ParseTreeType.ARRAY_LITERAL_EXPRESSION */ :
		return this_.ProcessArrayLiteral(base_node_.AsArrayLiteralExpression())
	case int32(72) /* ParseTreeType.TEMPLATE_LITERAL_EXPRESSION */ :
		return this_.ProcessTemplateLiteral(base_node_.AsTemplateLiteralExpression())
	case int32(73) /* ParseTreeType.TEMPLATE_LITERAL_PORTION */ :
		return this_.ProcessTemplateLiteralPortion(base_node_.AsTemplateLiteralPortion())
	case int32(74) /* ParseTreeType.TEMPLATE_SUBSTITUTION */ :
		return this_.ProcessTemplateSubstitution(base_node_.AsTemplateSubstitution())
	case int32(41) /* ParseTreeType.UNARY_EXPRESSION */ :
		return this_.ProcessUnaryExpression(base_node_.AsUnaryExpression())
	case int32(2) /* ParseTreeType.BLOCK */ :
		return this_.ProcessBlock(base_node_.AsBlock())
	case int32(14) /* ParseTreeType.BREAK_STATEMENT */ :
		return this_.ProcessBreakStatement(base_node_.AsBreakStatement())
	case int32(47) /* ParseTreeType.CALL_EXPRESSION */ :
		return this_.ProcessFunctionCall(base_node_.AsCallExpression())
	case int32(48) /* ParseTreeType.OPT_CHAIN_CALL_EXPRESSION */ :
		return this_.ProcessOptChainFunctionCall(base_node_.AsOptChainCallExpression())
	case int32(17) /* ParseTreeType.CASE_CLAUSE */ :
		return this_.ProcessSwitchCase(base_node_.AsCaseClause())
	case int32(18) /* ParseTreeType.DEFAULT_CLAUSE */ :
		return this_.ProcessSwitchDefault(base_node_.AsDefaultClause())
	case int32(22) /* ParseTreeType.CATCH */ :
		return this_.ProcessCatchClause(base_node_.AsCatch())
	case int32(13) /* ParseTreeType.CONTINUE_STATEMENT */ :
		return this_.ProcessContinueStatement(base_node_.AsContinueStatement())
	case int32(8) /* ParseTreeType.DO_WHILE_STATEMENT */ :
		return this_.ProcessDoLoop(base_node_.AsDoWhileStatement())
	case int32(5) /* ParseTreeType.EMPTY_STATEMENT */ :
		return this_.ProcessEmptyStatement(base_node_.AsEmptyStatement())
	case int32(6) /* ParseTreeType.EXPRESSION_STATEMENT */ :
		return this_.ProcessExpressionStatement(base_node_.AsExpressionStatement())
	case int32(24) /* ParseTreeType.DEBUGGER_STATEMENT */ :
		return this_.ProcessDebuggerStatement(base_node_.AsDebuggerStatement())
	case int32(25) /* ParseTreeType.THIS_EXPRESSION */ :
		return this_.ProcessThisExpression(base_node_.AsThisExpression())
	case int32(11) /* ParseTreeType.FOR_STATEMENT */ :
		return this_.ProcessForLoop(base_node_.AsForStatement())
	case int32(10) /* ParseTreeType.FOR_IN_STATEMENT */ :
		return this_.ProcessForInLoop(base_node_.AsForInStatement())
	case int32(1) /* ParseTreeType.FUNCTION_DECLARATION */ :
		return this_.ProcessFunction(base_node_.AsFunctionDeclaration())
	case int32(50) /* ParseTreeType.MEMBER_LOOKUP_EXPRESSION */ :
		return this_.ProcessElementGet(base_node_.AsMemberLookupExpression())
	case int32(51) /* ParseTreeType.OPT_CHAIN_MEMBER_LOOKUP_EXPRESSION */ :
		return this_.ProcessOptChainElementGet(base_node_.AsOptionalMemberLookupExpression())
	case int32(43) /* ParseTreeType.MEMBER_EXPRESSION */ :
		return this_.ProcessPropertyGet(base_node_.AsMemberExpression())
	case int32(44) /* ParseTreeType.OPT_CHAIN_MEMBER_EXPRESSION */ :
		return this_.ProcessOptChainPropertyGet(base_node_.AsOptionalMemberExpression())
	case int32(40) /* ParseTreeType.CONDITIONAL_EXPRESSION */ :
		return this_.ProcessConditionalExpression(base_node_.AsConditionalExpression())
	case int32(7) /* ParseTreeType.IF_STATEMENT */ :
		return this_.ProcessIfStatement(base_node_.AsIfStatement())
	case int32(20) /* ParseTreeType.LABELLED_STATEMENT */ :
		return this_.ProcessLabeledStatement(base_node_.AsLabelledStatement())
	case int32(52) /* ParseTreeType.PAREN_EXPRESSION */ :
		return this_.ProcessParenthesizedExpression(base_node_.AsParenExpression())
	case int32(26) /* ParseTreeType.IDENTIFIER_EXPRESSION */ :
		return this_.ProcessName_2(base_node_.AsIdentifierExpression(), func() IRFactory_TransformDispatcher_IdentifierType {
			if this_.privateIdLhsOfInScope.InScope() {
				return StaticGet_IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE()
			}
			return StaticGet_IRFactory_TransformDispatcher_IdentifierType_STANDARD()
		}())
	case int32(45) /* ParseTreeType.NEW_EXPRESSION */ :
		return this_.ProcessNewExpression(base_node_.AsNewExpression())
	case int32(29) /* ParseTreeType.OBJECT_LITERAL_EXPRESSION */ :
		return this_.ProcessObjectLiteral(base_node_.AsObjectLiteralExpression())
	case int32(33) /* ParseTreeType.COMPUTED_PROPERTY_DEFINITION */ :
		return this_.ProcessComputedPropertyDefinition(base_node_.AsComputedPropertyDefinition())
	case int32(34) /* ParseTreeType.COMPUTED_PROPERTY_GETTER */ :
		return this_.ProcessComputedPropertyGetter(base_node_.AsComputedPropertyGetter())
	case int32(35) /* ParseTreeType.COMPUTED_PROPERTY_METHOD */ :
		return this_.ProcessComputedPropertyMethod(base_node_.AsComputedPropertyMethod())
	case int32(36) /* ParseTreeType.COMPUTED_PROPERTY_SETTER */ :
		return this_.ProcessComputedPropertySetter(base_node_.AsComputedPropertySetter())
	case int32(15) /* ParseTreeType.RETURN_STATEMENT */ :
		return this_.ProcessReturnStatement(base_node_.AsReturnStatement())
	case int32(42) /* ParseTreeType.UPDATE_EXPRESSION */ :
		return this_.ProcessUpdateExpression(base_node_.AsUpdateExpression())
	case int32(0) /* ParseTreeType.PROGRAM */ :
		return this_.ProcessAstRoot(base_node_.AsProgram())
	case int32(27) /* ParseTreeType.LITERAL_EXPRESSION */ :
		return this_.ProcessLiteralExpression(base_node_.AsLiteralExpression())
	case int32(19) /* ParseTreeType.SWITCH_STATEMENT */ :
		return this_.ProcessSwitchStatement(base_node_.AsSwitchStatement())
	case int32(21) /* ParseTreeType.THROW_STATEMENT */ :
		return this_.ProcessThrowStatement(base_node_.AsThrowStatement())
	case int32(23) /* ParseTreeType.TRY_STATEMENT */ :
		return this_.ProcessTryStatement(base_node_.AsTryStatement())
	case int32(3) /* ParseTreeType.VARIABLE_STATEMENT */ :
		return this_.ProcessVariableStatement(base_node_.AsVariableStatement())
	case int32(12) /* ParseTreeType.VARIABLE_DECLARATION_LIST */ :
		return this_.ProcessVariableDeclarationList(base_node_.AsVariableDeclarationList())
	case int32(4) /* ParseTreeType.VARIABLE_DECLARATION */ :
		return this_.ProcessVariableDeclaration(base_node_.AsVariableDeclaration())
	case int32(9) /* ParseTreeType.WHILE_STATEMENT */ :
		return this_.ProcessWhileLoop(base_node_.AsWhileStatement())
	case int32(16) /* ParseTreeType.WITH_STATEMENT */ :
		return this_.ProcessWithStatement(base_node_.AsWithStatement())
	case int32(38) /* ParseTreeType.COMMA_EXPRESSION */ :
		return this_.ProcessCommaExpression(base_node_.AsCommaExpression())
	case int32(60) /* ParseTreeType.NULL */ :
		return this_.ProcessNull(base_node_.AsNull())
	case int32(53) /* ParseTreeType.FINALLY */ :
		return this_.ProcessFinally(base_node_.AsFinally())
	case int32(37) /* ParseTreeType.MISSING_PRIMARY_EXPRESSION */ :
		return this_.ProcessMissingExpression(base_node_.AsMissingPrimaryExpression())
	case int32(32) /* ParseTreeType.PROPERTY_NAME_ASSIGNMENT */ :
		return this_.ProcessPropertyNameAssignment(base_node_.AsPropertyNameAssignment())
	case int32(30) /* ParseTreeType.GET_ACCESSOR */ :
		return this_.ProcessGetAccessor(base_node_.AsGetAccessor())
	case int32(31) /* ParseTreeType.SET_ACCESSOR */ :
		return this_.ProcessSetAccessor(base_node_.AsSetAccessor())
	case int32(57) /* ParseTreeType.FORMAL_PARAMETER_LIST */ :
		return this_.ProcessFormalParameterList(base_node_.AsFormalParameterList())
	case int32(49) /* ParseTreeType.CLASS_DECLARATION */ :
		return this_.ProcessClassDeclaration(base_node_.AsClassDeclaration())
	case int32(54) /* ParseTreeType.SUPER_EXPRESSION */ :
		return this_.ProcessSuper(base_node_.AsSuperExpression())
	case int32(76) /* ParseTreeType.NEW_TARGET_EXPRESSION */ :
		return this_.ProcessNewTarget(base_node_.AsNewTargetExpression())
	case int32(70) /* ParseTreeType.YIELD_EXPRESSION */ :
		return this_.ProcessYield(base_node_.AsYieldStatement())
	case int32(77) /* ParseTreeType.AWAIT_EXPRESSION */ :
		return this_.ProcessAwait(base_node_.AsAwaitExpression())
	case int32(68) /* ParseTreeType.FOR_OF_STATEMENT */ :
		return this_.ProcessForOf(base_node_.AsForOfStatement())
	case int32(69) /* ParseTreeType.FOR_AWAIT_OF_STATEMENT */ :
		return this_.ProcessForAwaitOf(base_node_.AsForAwaitOfStatement())
	case int32(64) /* ParseTreeType.EXPORT_DECLARATION */ :
		return this_.ProcessExportDecl(base_node_.AsExportDeclaration())
	case int32(65) /* ParseTreeType.EXPORT_SPECIFIER */ :
		return this_.ProcessExportSpec(base_node_.AsExportSpecifier())
	case int32(66) /* ParseTreeType.IMPORT_DECLARATION */ :
		return this_.ProcessImportDecl(base_node_.AsImportDeclaration())
	case int32(67) /* ParseTreeType.IMPORT_SPECIFIER */ :
		return this_.ProcessImportSpec(base_node_.AsImportSpecifier())
	case int32(78) /* ParseTreeType.DYNAMIC_IMPORT_EXPRESSION */ :
		return this_.ProcessDynamicImport(base_node_.AsDynamicImportExpression())
	case int32(79) /* ParseTreeType.IMPORT_META_EXPRESSION */ :
		return this_.ProcessImportMeta(base_node_.AsImportMetaExpression())
	case int32(55) /* ParseTreeType.ARRAY_PATTERN */ :
		return this_.ProcessArrayPattern(base_node_.AsArrayPattern())
	case int32(56) /* ParseTreeType.OBJECT_PATTERN */ :
		return this_.ProcessObjectPattern(base_node_.AsObjectPattern())
	case int32(71) /* ParseTreeType.DEFAULT_PARAMETER */ :
		return this_.ProcessDefaultParameter(base_node_.AsDefaultParameter())
	case int32(61) /* ParseTreeType.ITER_REST */ :
		return this_.ProcessIterRest(base_node_.AsIterRest())
	case int32(58) /* ParseTreeType.ITER_SPREAD */ :
		return this_.ProcessIterSpread(base_node_.AsIterSpread())
	case int32(62) /* ParseTreeType.OBJECT_REST */ :
		return this_.priv_TransformDispatcher_processObjectPatternElement_6888(pkg_parser.ToIface_ParseTree(base_node_.AsObjectRest()))
	case int32(59) /* ParseTreeType.OBJECT_SPREAD */ :
		return this_.ProcessObjectSpread(base_node_.AsObjectSpread())
	case int32(81) /* ParseTreeType.FIELD_DECLARATION */ :
		return this_.ProcessField(base_node_.AsFieldDeclaration())
	case int32(80) /* ParseTreeType.COMPUTED_PROPERTY_FIELD */ :
		return this_.ProcessComputedPropertyField(base_node_.AsComputedPropertyField())
	case int32(46) /* ParseTreeType.ARGUMENT_LIST */ :
	default:
	}
	return this_.ProcessIllegalToken(node)
}

// Java class: com.google.javascript.jscomp.parsing.IRFactory.TransformDispatcher.ScopeTracker
type IRFactory_TransformDispatcher_ScopeTracker struct {
	usageCount int32 `gson:"usageCount"`
	idHash_    atomic.Int32
}

func (this_ *IRFactory_TransformDispatcher_ScopeTracker) BindSelf_(s any) {
	_ = s
}
func (this_ *IRFactory_TransformDispatcher_ScopeTracker) IsJavaNil_() bool { return this_ == nil }
func (*IRFactory_TransformDispatcher_ScopeTracker) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.IRFactory$TransformDispatcher$ScopeTracker", "ScopeTracker"
}
func (this_ *IRFactory_TransformDispatcher_ScopeTracker) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*IRFactory_TransformDispatcher_ScopeTracker) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.IRFactory.TransformDispatcher.ScopeTracker", "com.google.javascript.jscomp.parsing.IRFactory$TransformDispatcher$ScopeTracker", func(v_ any) bool {
		return j2gort.InstanceOf[*IRFactory_TransformDispatcher_ScopeTracker](v_)
	})
}

func (this_ *IRFactory_TransformDispatcher_ScopeTracker) initFields_IRFactory_TransformDispatcher_ScopeTracker() {
	_ = this_

}

// Java constructor: ScopeTracker()
func NewIRFactory_TransformDispatcher_ScopeTracker() *IRFactory_TransformDispatcher_ScopeTracker {
	this_ := &IRFactory_TransformDispatcher_ScopeTracker{}
	InitIRFactory_TransformDispatcher_ScopeTracker(this_)
	return this_
}

func InitIRFactory_TransformDispatcher_ScopeTracker(this_ *IRFactory_TransformDispatcher_ScopeTracker) {
	this_.initFields_IRFactory_TransformDispatcher_ScopeTracker()
}

// Java method: ScopeTracker.inScope()
func (this_ *IRFactory_TransformDispatcher_ScopeTracker) InScope() bool {
	return (this_.usageCount > int32(0))
}

// Java method: ScopeTracker.enter()
func (this_ *IRFactory_TransformDispatcher_ScopeTracker) Enter() {
	if !(this_.usageCount >= int32(0)) {
		j2gort.CheckState(false)
	}
	this_.usageCount = this_.usageCount + 1
}

// Java method: ScopeTracker.exit()
func (this_ *IRFactory_TransformDispatcher_ScopeTracker) Exit() {
	this_.usageCount = this_.usageCount - 1
	if !(this_.usageCount >= int32(0)) {
		j2gort.CheckState(false)
	}
}

// Java enum: com.google.javascript.jscomp.parsing.IRFactory.TransformDispatcher.IdentifierType
type IRFactory_TransformDispatcher_IdentifierType uint8

const (
	enum_IRFactory_TransformDispatcher_IdentifierType_STANDARD       IRFactory_TransformDispatcher_IdentifierType = 1
	enum_IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE IRFactory_TransformDispatcher_IdentifierType = 2
)

var enumBaseTable_IRFactory_TransformDispatcher_IdentifierType = [...]j2gort.EnumBase{
	{Name_: "STANDARD", Ordinal_: 0},
	{Name_: "CAN_BE_PRIVATE", Ordinal_: 1},
}
var enumValues_IRFactory_TransformDispatcher_IdentifierType = []IRFactory_TransformDispatcher_IdentifierType{1, 2}

func StaticGet_IRFactory_TransformDispatcher_IdentifierType_STANDARD() IRFactory_TransformDispatcher_IdentifierType {
	return enum_IRFactory_TransformDispatcher_IdentifierType_STANDARD
}
func StaticGet_IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE() IRFactory_TransformDispatcher_IdentifierType {
	return enum_IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE
}

func IRFactory_TransformDispatcher_IdentifierType_Values() []IRFactory_TransformDispatcher_IdentifierType {
	return j2gort.CloneSlice(enumValues_IRFactory_TransformDispatcher_IdentifierType)
}

func IRFactory_TransformDispatcher_IdentifierType_ValueOf(name string) IRFactory_TransformDispatcher_IdentifierType {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_IRFactory_TransformDispatcher_IdentifierType {
		if enumBaseTable_IRFactory_TransformDispatcher_IdentifierType[i_].Name_ == name {
			return IRFactory_TransformDispatcher_IdentifierType(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.parsing.IRFactory.TransformDispatcher.IdentifierType." + name))
}

func (this_ IRFactory_TransformDispatcher_IdentifierType) IsJavaNil_() bool                { return this_ == 0 }
func (IRFactory_TransformDispatcher_IdentifierType) HasIdentityEquals_()                   {}
func (IRFactory_TransformDispatcher_IdentifierType) Implements_j2gort_Ext_java_lang_Enum() {}
func (IRFactory_TransformDispatcher_IdentifierType) Implements_j2gort_Ext_java_lang_constant_Constable() {
}
func (this_ IRFactory_TransformDispatcher_IdentifierType) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (IRFactory_TransformDispatcher_IdentifierType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.IRFactory$TransformDispatcher$IdentifierType", "IdentifierType"
}
func (this_ IRFactory_TransformDispatcher_IdentifierType) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_IRFactory_TransformDispatcher_IdentifierType[this_-1]
}
func (this_ IRFactory_TransformDispatcher_IdentifierType) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ IRFactory_TransformDispatcher_IdentifierType) EnumOrdinal_() int32 {
	return int32(this_) - 1
}
func (this_ IRFactory_TransformDispatcher_IdentifierType) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_IRFactory_TransformDispatcher_IdentifierType[this_-1].Name_
}
func (this_ IRFactory_TransformDispatcher_IdentifierType) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_IRFactory_TransformDispatcher_IdentifierType[this_-1].Name_
}
func (this_ IRFactory_TransformDispatcher_IdentifierType) Equals(other_ any) bool {
	o_, ok_ := other_.(IRFactory_TransformDispatcher_IdentifierType)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ IRFactory_TransformDispatcher_IdentifierType) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_IRFactory_TransformDispatcher_IdentifierType[this_-1].HashCode()
}
func (this_ IRFactory_TransformDispatcher_IdentifierType) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_IRFactory_TransformDispatcher_IdentifierType[this_-1].CompareTo(other_)
}
func (this_ IRFactory_TransformDispatcher_IdentifierType) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_IRFactory_TransformDispatcher_IdentifierType[this_-1].Name_)
}
func init() {
	vals_ := []any{IRFactory_TransformDispatcher_IdentifierType(1), IRFactory_TransformDispatcher_IdentifierType(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.IRFactory.TransformDispatcher.IdentifierType", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.IRFactory$TransformDispatcher$IdentifierType", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.IRFactory.TransformDispatcher.IdentifierType", "com.google.javascript.jscomp.parsing.IRFactory$TransformDispatcher$IdentifierType", func(v_ any) bool {
		return j2gort.InstanceOf[IRFactory_TransformDispatcher_IdentifierType](v_)
	})
}
