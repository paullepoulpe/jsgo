/*
 * Copyright 2009 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ReplaceCssNames.java (generated mechanically by j2go).

package jscomp

import (
	"strconv"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.ReplaceCssNames
type ReplaceCssNames struct {
	compiler                     AbstractCompiler                  `gson:"compiler"`
	astFactory                   *AstFactory                       `gson:"astFactory"`
	cssNameCollector             ReplaceCssNames_CssNameCollector  `gson:"cssNameCollector"`
	cssNamesBySymbol             j2gort.OrderedMap[string, string] `gson:"cssNamesBySymbol"`
	classesObjectsQualifiedNames j2gort.OrderedSet[string]         `gson:"classesObjectsQualifiedNames"`
	symbolMap                    CssRenamingMap                    `gson:"symbolMap"`
	skiplist                     j2gort.Ext_java_util_Set          `gson:"skiplist"`
	idHash_                      atomic.Int32
}

func (this_ *ReplaceCssNames) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceCssNames) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceCssNames) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceCssNames", "ReplaceCssNames"
}
func (this_ *ReplaceCssNames) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReplaceCssNames) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceCssNames", "com.google.javascript.jscomp.ReplaceCssNames", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceCssNames](v_)
	})
}
func (*ReplaceCssNames) Implements_jscomp_CompilerPass() {}

var static_ReplaceCssNames_GET_CSS_NAME_FUNCTION pkg_rhino.Node

func StaticGet_ReplaceCssNames_GET_CSS_NAME_FUNCTION() pkg_rhino.Node {
	clinit_ReplaceCssNames()
	return static_ReplaceCssNames_GET_CSS_NAME_FUNCTION
}
func StaticSet_ReplaceCssNames_GET_CSS_NAME_FUNCTION(v pkg_rhino.Node) {
	clinit_ReplaceCssNames()
	static_ReplaceCssNames_GET_CSS_NAME_FUNCTION = v
}

var static_ReplaceCssNames_INVALID_NUM_ARGUMENTS_ERROR *DiagnosticType

func StaticGet_ReplaceCssNames_INVALID_NUM_ARGUMENTS_ERROR() *DiagnosticType {
	clinit_ReplaceCssNames()
	return static_ReplaceCssNames_INVALID_NUM_ARGUMENTS_ERROR
}
func StaticSet_ReplaceCssNames_INVALID_NUM_ARGUMENTS_ERROR(v *DiagnosticType) {
	clinit_ReplaceCssNames()
	static_ReplaceCssNames_INVALID_NUM_ARGUMENTS_ERROR = v
}

var static_ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR *DiagnosticType

func StaticGet_ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR() *DiagnosticType {
	clinit_ReplaceCssNames()
	return static_ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR
}
func StaticSet_ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR(v *DiagnosticType) {
	clinit_ReplaceCssNames()
	static_ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR = v
}

var static_ReplaceCssNames_UNEXPECTED_STRING_LITERAL_ERROR *DiagnosticType

func StaticGet_ReplaceCssNames_UNEXPECTED_STRING_LITERAL_ERROR() *DiagnosticType {
	clinit_ReplaceCssNames()
	return static_ReplaceCssNames_UNEXPECTED_STRING_LITERAL_ERROR
}
func StaticSet_ReplaceCssNames_UNEXPECTED_STRING_LITERAL_ERROR(v *DiagnosticType) {
	clinit_ReplaceCssNames()
	static_ReplaceCssNames_UNEXPECTED_STRING_LITERAL_ERROR = v
}

var static_ReplaceCssNames_NESTED_CALL_ERROR *DiagnosticType

func StaticGet_ReplaceCssNames_NESTED_CALL_ERROR() *DiagnosticType {
	clinit_ReplaceCssNames()
	return static_ReplaceCssNames_NESTED_CALL_ERROR
}
func StaticSet_ReplaceCssNames_NESTED_CALL_ERROR(v *DiagnosticType) {
	clinit_ReplaceCssNames()
	static_ReplaceCssNames_NESTED_CALL_ERROR = v
}

var static_ReplaceCssNames_UNKNOWN_SYMBOL_WARNING *DiagnosticType

func StaticGet_ReplaceCssNames_UNKNOWN_SYMBOL_WARNING() *DiagnosticType {
	clinit_ReplaceCssNames()
	return static_ReplaceCssNames_UNKNOWN_SYMBOL_WARNING
}
func StaticSet_ReplaceCssNames_UNKNOWN_SYMBOL_WARNING(v *DiagnosticType) {
	clinit_ReplaceCssNames()
	static_ReplaceCssNames_UNKNOWN_SYMBOL_WARNING = v
}

var static_ReplaceCssNames_UNEXPECTED_SASS_GENERATED_CSS_TS_ERROR *DiagnosticType

func StaticGet_ReplaceCssNames_UNEXPECTED_SASS_GENERATED_CSS_TS_ERROR() *DiagnosticType {
	clinit_ReplaceCssNames()
	return static_ReplaceCssNames_UNEXPECTED_SASS_GENERATED_CSS_TS_ERROR
}
func StaticSet_ReplaceCssNames_UNEXPECTED_SASS_GENERATED_CSS_TS_ERROR(v *DiagnosticType) {
	clinit_ReplaceCssNames()
	static_ReplaceCssNames_UNEXPECTED_SASS_GENERATED_CSS_TS_ERROR = v
}

var static_ReplaceCssNames_UNKNOWN_SYMBOL_ERROR *DiagnosticType

func StaticGet_ReplaceCssNames_UNKNOWN_SYMBOL_ERROR() *DiagnosticType {
	clinit_ReplaceCssNames()
	return static_ReplaceCssNames_UNKNOWN_SYMBOL_ERROR
}
func StaticSet_ReplaceCssNames_UNKNOWN_SYMBOL_ERROR(v *DiagnosticType) {
	clinit_ReplaceCssNames()
	static_ReplaceCssNames_UNKNOWN_SYMBOL_ERROR = v
}

var static_ReplaceCssNames_INVALID_USE_OF_CLASSES_OBJECT_ERROR *DiagnosticType

func StaticGet_ReplaceCssNames_INVALID_USE_OF_CLASSES_OBJECT_ERROR() *DiagnosticType {
	clinit_ReplaceCssNames()
	return static_ReplaceCssNames_INVALID_USE_OF_CLASSES_OBJECT_ERROR
}
func StaticSet_ReplaceCssNames_INVALID_USE_OF_CLASSES_OBJECT_ERROR(v *DiagnosticType) {
	clinit_ReplaceCssNames()
	static_ReplaceCssNames_INVALID_USE_OF_CLASSES_OBJECT_ERROR = v
}

var static_ReplaceCssNames_GOOG_SET_CSS_NAME_MAPPING pkg_rhino.Node

func StaticGet_ReplaceCssNames_GOOG_SET_CSS_NAME_MAPPING() pkg_rhino.Node {
	clinit_ReplaceCssNames()
	return static_ReplaceCssNames_GOOG_SET_CSS_NAME_MAPPING
}
func StaticSet_ReplaceCssNames_GOOG_SET_CSS_NAME_MAPPING(v pkg_rhino.Node) {
	clinit_ReplaceCssNames()
	static_ReplaceCssNames_GOOG_SET_CSS_NAME_MAPPING = v
}

var clinitDone_ReplaceCssNames uint32
var clinitState_ReplaceCssNames j2gort.ClinitState

func clinit_ReplaceCssNames() {
	if atomic.LoadUint32(&clinitDone_ReplaceCssNames) != 0 {
		return
	}
	clinitSlow_ReplaceCssNames()
}

func clinitSlow_ReplaceCssNames() {
	if !j2gort.EnterClinit(&clinitState_ReplaceCssNames, &clinitDone_ReplaceCssNames) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ReplaceCssNames, &clinitDone_ReplaceCssNames)
	static_ReplaceCssNames_GET_CSS_NAME_FUNCTION = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "getCssName")
	static_ReplaceCssNames_INVALID_NUM_ARGUMENTS_ERROR = DiagnosticType_Error("JSC_GETCSSNAME_NUM_ARGS", "goog.getCssName called with \"{0}\" arguments, expected 1 or 2.")
	static_ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR = DiagnosticType_Error("JSC_GETCSSNAME_STRING_LITERAL_EXPECTED", "goog.getCssName called with invalid argument, string literal expected.  Was \"{0}\".")
	static_ReplaceCssNames_UNEXPECTED_STRING_LITERAL_ERROR = DiagnosticType_Error("JSC_GETCSSNAME_UNEXPECTED_STRING_LITERAL", "goog.getCssName called with invalid arguments, string literal passed as first of two arguments.  Did you mean goog.getCssName(\"{0}-{1}\")?")
	static_ReplaceCssNames_NESTED_CALL_ERROR = DiagnosticType_Error("JSC_GETCSSNAME_NESTED_CALL", "goog.getCssName: nested call is not allowed.")
	static_ReplaceCssNames_UNKNOWN_SYMBOL_WARNING = DiagnosticType_Warning("JSC_GETCSSNAME_UNKNOWN_CSS_SYMBOL", "goog.getCssName called with unrecognized symbol \"{0}\" in class \"{1}\".")
	static_ReplaceCssNames_UNEXPECTED_SASS_GENERATED_CSS_TS_ERROR = DiagnosticType_Error("JSC_UNEXPECTED_SASS_GENERATED_CSS_TS", "@sass_generated_css_ts JSDoc annotation is only allowed on .css.closure.js files.")
	static_ReplaceCssNames_UNKNOWN_SYMBOL_ERROR = DiagnosticType_Error("JSC_UNKNOWN_CSS_SYMBOL_IN_CLASSES_OBJECT", "Symbol was not defined \"{0}\" in classes object.")
	static_ReplaceCssNames_INVALID_USE_OF_CLASSES_OBJECT_ERROR = DiagnosticType_Error("JSC_INVALID_USE_OF_CLASSES_OBJECT", "invalid use of generated classes object. Only accessing its members is allowed.")
	static_ReplaceCssNames_GOOG_SET_CSS_NAME_MAPPING = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "setCssNameMapping")
}

func (this_ *ReplaceCssNames) initFields_ReplaceCssNames() {
	_ = this_
}

// Java constructor: ReplaceCssNames(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable CssRenamingMap symbolMap, com.google.javascript.jscomp.ReplaceCssNames.CssNameCollector cssNameCollector, java.util.@org.jspecify.annotations.Nullable Set<java.lang.String> skiplist)
func NewReplaceCssNames(compiler AbstractCompiler, symbolMap CssRenamingMap, cssNameCollector ReplaceCssNames_CssNameCollector, skiplist j2gort.Ext_java_util_Set) *ReplaceCssNames {
	clinit_ReplaceCssNames()
	this_ := &ReplaceCssNames{}
	InitReplaceCssNames(this_, compiler, symbolMap, cssNameCollector, skiplist)
	return this_
}

func InitReplaceCssNames(this_ *ReplaceCssNames, compiler AbstractCompiler, symbolMap CssRenamingMap, cssNameCollector ReplaceCssNames_CssNameCollector, skiplist j2gort.Ext_java_util_Set) {
	_ = compiler
	_ = symbolMap
	_ = cssNameCollector
	_ = skiplist
	this_.initFields_ReplaceCssNames()
	this_.compiler = compiler
	this_.astFactory = compiler.CreateAstFactory()
	this_.symbolMap = symbolMap
	this_.cssNameCollector = cssNameCollector
	this_.skiplist = skiplist
	this_.cssNamesBySymbol = j2gort.NewOrderedMap[string, string]()
	this_.classesObjectsQualifiedNames = j2gort.NewOrderedSet[string]()
}

// Java method: ReplaceCssNames.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ReplaceCssNames) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, NewReplaceCssNames_FindSetCssNameTraversal(this_))
	var getCssNameInstances j2gort.Ext_java_util_List = this_.priv_ReplaceCssNames_gatherGetCssNameInstances_4829(root)
	_ = getCssNameInstances
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(getCssNameInstances))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			getCssNameInstance := j2gort.CastPtr[ReplaceCssNames_GetCssNameInstance](rawElem_2)
			_ = getCssNameInstance
			var cssClassName string = getCssNameInstance.GetCssClassName()
			_ = cssClassName
			if getCssNameInstance.IsCssClosureFileClassesMember() {
				this_.cssNamesBySymbol.Put(getCssNameInstance.GetCssClosureClassesMemberName(), cssClassName)
				this_.classesObjectsQualifiedNames.Add(getCssNameInstance.GetCssClosureClassesQualifiedName())
			} else {
				this_.cssNameCollector.Add(cssClassName)
			}
			getCssNameInstance.ReplaceWithExpression()
		}
	}
	NodeTraversal_Traverse(this_.compiler, root, NewReplaceCssNames_CountCssNamesBySymbol(this_))
}

// Java method: ReplaceCssNames.gatherGetCssNameInstances(com.google.javascript.rhino.Node root)
func (this_ *ReplaceCssNames) priv_ReplaceCssNames_gatherGetCssNameInstances_4829(root pkg_rhino.Node) j2gort.Ext_java_util_List {
	_ = root
	var val_3 ReplaceCssNames_GatherCssNamesTraversal
	val_3.outer_ = this_
	InitReplaceCssNames_GatherCssNamesTraversal(&val_3)
	var gatherCssNamesTraversal *ReplaceCssNames_GatherCssNamesTraversal = &val_3
	_ = gatherCssNamesTraversal
	NodeTraversal_Traverse(this_.compiler, root, gatherCssNamesTraversal)
	return gatherCssNamesTraversal.listOfCssNameInstances
}

// Java method: ReplaceCssNames.isNotInCssClosureFile(com.google.javascript.rhino.Node node)
func (this_ *ReplaceCssNames) priv_ReplaceCssNames_isNotInCssClosureFile_4830(node pkg_rhino.Node) bool {
	_ = node
	var sourceFileName string = node.GetSourceFileName()
	_ = sourceFileName
	return ((sourceFileName == j2gort.NullString) || !(javalang.StringEndsWith(j2gort.CheckNotNullString(sourceFileName), ".css.closure.js")))
}

// Java method: ReplaceCssNames.isGetCssNameCall(com.google.javascript.rhino.Node node)
func (this_ *ReplaceCssNames) priv_ReplaceCssNames_isGetCssNameCall_4831(node pkg_rhino.Node) bool {
	_ = node
	return (node.IsCall() && node.GetFirstChild().MatchesQualifiedName_1_Node(StaticGet_ReplaceCssNames_GET_CSS_NAME_FUNCTION()))
}

// Java method: ReplaceCssNames.createGetCssNameInstance(com.google.javascript.rhino.Node n, java.lang.String cssClosureClassesQualifiedName)
func (this_ *ReplaceCssNames) priv_ReplaceCssNames_createGetCssNameInstance_4832(n pkg_rhino.Node, cssClosureClassesQualifiedName string) *ReplaceCssNames_GetCssNameInstance {
	_ = n
	_ = cssClosureClassesQualifiedName
	var validationError *JSError = this_.priv_ReplaceCssNames_validateGetCssNameCall_4833(n)
	_ = validationError
	return NewReplaceCssNames_GetCssNameInstance(this_, n, validationError, cssClosureClassesQualifiedName)
}

// Java method: ReplaceCssNames.validateGetCssNameCall(com.google.javascript.rhino.Node callNode)
func (this_ *ReplaceCssNames) priv_ReplaceCssNames_validateGetCssNameCall_4833(callNode pkg_rhino.Node) *JSError {
	_ = callNode
	var childCount int32 = callNode.GetChildCount()
	_ = childCount
	switch childCount {
	case int32(2):
		return this_.priv_ReplaceCssNames_validateSingleArgGetCssNameCall_4834(callNode)
	case int32(3):
		return this_.priv_ReplaceCssNames_validateTwoArgGetCssNameCall_4835(callNode)
	default:
		return JSError_Make_3(callNode, StaticGet_ReplaceCssNames_INVALID_NUM_ARGUMENTS_ERROR(), []string{strconv.FormatInt(int64(childCount), 10)})
	}
	panic("unreachable switch expression")
}

// Java method: ReplaceCssNames.validateSingleArgGetCssNameCall(com.google.javascript.rhino.Node callNode)
func (this_ *ReplaceCssNames) priv_ReplaceCssNames_validateSingleArgGetCssNameCall_4834(callNode pkg_rhino.Node) *JSError {
	_ = callNode
	var stringLiteralArg pkg_rhino.Node = j2gort.CheckNotNullNode(callNode.GetLastChild())
	_ = stringLiteralArg
	if !(stringLiteralArg.IsStringLit()) {
		return JSError_Make_3(callNode, StaticGet_ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR(), []string{j2gort.CheckNotNullEnum(stringLiteralArg.GetToken()).String()})
	}
	return nil
}

// Java method: ReplaceCssNames.validateTwoArgGetCssNameCall(com.google.javascript.rhino.Node callNode)
func (this_ *ReplaceCssNames) priv_ReplaceCssNames_validateTwoArgGetCssNameCall_4835(callNode pkg_rhino.Node) *JSError {
	_ = callNode
	var firstArg pkg_rhino.Node = j2gort.CheckNotNullNode(callNode.GetSecondChild())
	_ = firstArg
	var secondArg pkg_rhino.Node = j2gort.CheckNotNullNode(firstArg.GetNext())
	_ = secondArg
	if !(secondArg.IsStringLit()) {
		return JSError_Make_3(callNode, StaticGet_ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR(), []string{j2gort.CheckNotNullEnum(secondArg.GetToken()).String()})
	}
	if firstArg.IsStringLit() {
		return JSError_Make_3(callNode, StaticGet_ReplaceCssNames_UNEXPECTED_STRING_LITERAL_ERROR(), []string{firstArg.GetString(), secondArg.GetString()})
	}
	if this_.priv_ReplaceCssNames_isGetCssNameCall_4831(firstArg) {
		return JSError_Make_3(firstArg, StaticGet_ReplaceCssNames_NESTED_CALL_ERROR(), []string{})
	}
	return nil
}

// Java method: ReplaceCssNames.processStringNode(com.google.javascript.rhino.Node n)
func (this_ *ReplaceCssNames) priv_ReplaceCssNames_processStringNode_4836(n pkg_rhino.Node) {
	_ = n
	var name string = n.GetString()
	_ = name
	if (this_.skiplist != nil) && j2gort.SetContainsStr(this_.skiplist, name) {
		return
	}
	var parts []string = j2gort.JStrSplit(name, "-")
	_ = parts
	if this_.symbolMap != nil && !this_.symbolMap.IsJavaNil_() {
		var replacement string = j2gort.NullString
		_ = replacement
		if javalang.StringStartsWith(j2gort.CheckNotNullString(name), "--") {
			replacement = this_.symbolMap.Get(name)
			if replacement == j2gort.NullString {
				this_.compiler.Report_1(JSError_Make_3(n, StaticGet_ReplaceCssNames_UNKNOWN_SYMBOL_WARNING(), []string{name, name}))
				return
			}
		} else {
			switch (this_.symbolMap.GetStyle()).Ordinal() {
			case int32(0) /* Style.BY_WHOLE */ :
				replacement = this_.symbolMap.Get(name)
				if replacement == j2gort.NullString {
					this_.compiler.Report_1(JSError_Make_3(n, StaticGet_ReplaceCssNames_UNKNOWN_SYMBOL_WARNING(), []string{name, name}))
					return
				}
			case int32(1) /* Style.BY_PART */ :
				var replaced []string = j2gort.NewStringArray(int32(len(parts)))
				_ = replaced
				{
					var i int32 = int32(0)
					_ = i
					for ; i < int32(len(parts)); i = i + 1 {
						var part string = this_.symbolMap.Get(parts[i])
						_ = part
						if part == j2gort.NullString {
							this_.compiler.Report_1(JSError_Make_3(n, StaticGet_ReplaceCssNames_UNKNOWN_SYMBOL_WARNING(), []string{parts[i], name}))
							return
						}
						{
							a_4 := replaced
							i_5 := i
							a_4[i_5] = part
						}
					}
				}
				replacement = j2gort.Ext_com_google_common_base_Joiner_On_1_string("-").Join_1_slice_any(func() []any {
					arr_6 := replaced
					if arr_6 == nil {
						return nil
					}
					out_7 := make([]any, len(arr_6))
					for i_, el_8 := range arr_6 {
						out_7[i_] = j2gort.BoxString(el_8)
					}
					return out_7
				}())
			}
		}
		n.SetString(replacement)
	}
}

// Java interface: com.google.javascript.jscomp.ReplaceCssNames.CssNameCollector
type ReplaceCssNames_CssNameCollector interface {
	Implements_jscomp_ReplaceCssNames_CssNameCollector()
	IsJavaNil_() bool
	Add(string)
}

func ToIface_ReplaceCssNames_CssNameCollector[T interface {
	ReplaceCssNames_CssNameCollector
	comparable
}](v T) ReplaceCssNames_CssNameCollector {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_ReplaceCssNames_CssNameCollector[T ReplaceCssNames_CssNameCollector](v T) ReplaceCssNames_CssNameCollector {
	var vi ReplaceCssNames_CssNameCollector = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceCssNames.CssNameCollector", "com.google.javascript.jscomp.ReplaceCssNames$CssNameCollector", func(v_ any) bool {
		return j2gort.InstanceOf[ReplaceCssNames_CssNameCollector](v_)
	})
}

type lambda_ReplaceCssNames_CssNameCollector struct {
	fn func(string)
}

func (l *lambda_ReplaceCssNames_CssNameCollector) IsJavaNil_() bool { return l == nil }
func NewLambda_ReplaceCssNames_CssNameCollector(fn func(string)) ReplaceCssNames_CssNameCollector {
	if fn == nil {
		return nil
	}
	return &lambda_ReplaceCssNames_CssNameCollector{fn: fn}
}

func (*lambda_ReplaceCssNames_CssNameCollector) Implements_jscomp_ReplaceCssNames_CssNameCollector() {
}
func (l *lambda_ReplaceCssNames_CssNameCollector) Add(p0_ string) {
	l.fn(p0_)
}

// Java class: com.google.javascript.jscomp.ReplaceCssNames.FindSetCssNameTraversal
type ReplaceCssNames_FindSetCssNameTraversal struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *ReplaceCssNames
}

func (this_ *ReplaceCssNames_FindSetCssNameTraversal) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceCssNames_FindSetCssNameTraversal) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceCssNames_FindSetCssNameTraversal) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceCssNames$FindSetCssNameTraversal", "FindSetCssNameTraversal"
}
func (this_ *ReplaceCssNames_FindSetCssNameTraversal) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*ReplaceCssNames_FindSetCssNameTraversal) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceCssNames.FindSetCssNameTraversal", "com.google.javascript.jscomp.ReplaceCssNames$FindSetCssNameTraversal", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceCssNames_FindSetCssNameTraversal](v_)
	})
}
func (this_ *ReplaceCssNames_FindSetCssNameTraversal) SetOuter_(o *ReplaceCssNames) { this_.outer_ = o }
func (this_ *ReplaceCssNames_FindSetCssNameTraversal) GetOuter_() *ReplaceCssNames {
	return this_.outer_
}
func (*ReplaceCssNames_FindSetCssNameTraversal) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*ReplaceCssNames_FindSetCssNameTraversal) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *ReplaceCssNames_FindSetCssNameTraversal) initFields_ReplaceCssNames_FindSetCssNameTraversal() {
	_ = this_
}

// Java constructor: FindSetCssNameTraversal()
func NewReplaceCssNames_FindSetCssNameTraversal(outer_ *ReplaceCssNames) *ReplaceCssNames_FindSetCssNameTraversal {
	this_ := &ReplaceCssNames_FindSetCssNameTraversal{}
	this_.outer_ = outer_
	InitReplaceCssNames_FindSetCssNameTraversal(this_)
	return this_
}

func InitReplaceCssNames_FindSetCssNameTraversal(this_ *ReplaceCssNames_FindSetCssNameTraversal) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_ReplaceCssNames_FindSetCssNameTraversal()
}

// Java method: FindSetCssNameTraversal.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ReplaceCssNames_FindSetCssNameTraversal) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(n.IsCall()) || !(parent.IsExprResult()) {
		return
	}
	var callee pkg_rhino.Node = n.GetFirstChild()
	_ = callee
	if !(callee.MatchesQualifiedName_1_Node(StaticGet_ReplaceCssNames_GOOG_SET_CSS_NAME_MAPPING())) {
		return
	}
	var cssRenamingMap CssRenamingMap = ProcessClosurePrimitives_ProcessSetCssNameMapping(this_.outer_.compiler, n, parent)
	_ = cssRenamingMap
	this_.outer_.symbolMap = cssRenamingMap
	this_.outer_.compiler.ReportChangeToEnclosingScope(parent)
	parent.Detach()
}

// Java class: com.google.javascript.jscomp.ReplaceCssNames.GatherCssNamesTraversal
type ReplaceCssNames_GatherCssNamesTraversal struct {
	outer_                 *ReplaceCssNames
	listOfCssNameInstances j2gort.Ext_java_util_List                               `gson:"listOfCssNameInstances"`
	traversalState         *ReplaceCssNames_GatherCssNamesTraversal_TraversalState `gson:"traversalState"`
	idHash_                atomic.Int32
}

func (this_ *ReplaceCssNames_GatherCssNamesTraversal) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceCssNames_GatherCssNamesTraversal) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceCssNames_GatherCssNamesTraversal) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceCssNames$GatherCssNamesTraversal", "GatherCssNamesTraversal"
}
func (this_ *ReplaceCssNames_GatherCssNamesTraversal) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReplaceCssNames_GatherCssNamesTraversal) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceCssNames.GatherCssNamesTraversal", "com.google.javascript.jscomp.ReplaceCssNames$GatherCssNamesTraversal", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceCssNames_GatherCssNamesTraversal](v_)
	})
}
func (this_ *ReplaceCssNames_GatherCssNamesTraversal) SetOuter_(o *ReplaceCssNames) { this_.outer_ = o }
func (this_ *ReplaceCssNames_GatherCssNamesTraversal) GetOuter_() *ReplaceCssNames {
	return this_.outer_
}
func (*ReplaceCssNames_GatherCssNamesTraversal) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *ReplaceCssNames_GatherCssNamesTraversal) initFields_ReplaceCssNames_GatherCssNamesTraversal() {
	_ = this_
	this_.listOfCssNameInstances = j2gort.NewExt_java_util_ArrayList()
	this_.traversalState = NewReplaceCssNames_GatherCssNamesTraversal_TraversalState()
}

// Java constructor: GatherCssNamesTraversal()
func NewReplaceCssNames_GatherCssNamesTraversal(outer_ *ReplaceCssNames) *ReplaceCssNames_GatherCssNamesTraversal {
	this_ := &ReplaceCssNames_GatherCssNamesTraversal{}
	this_.outer_ = outer_
	InitReplaceCssNames_GatherCssNamesTraversal(this_)
	return this_
}

func InitReplaceCssNames_GatherCssNamesTraversal(this_ *ReplaceCssNames_GatherCssNamesTraversal) {
	this_.initFields_ReplaceCssNames_GatherCssNamesTraversal()
}

// Java method: GatherCssNamesTraversal.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ReplaceCssNames_GatherCssNamesTraversal) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	var sassGeneratedCssTsExpert *ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert = this_.CreateSassGeneratedCssTsExpert(n)
	_ = sassGeneratedCssTsExpert
	if sassGeneratedCssTsExpert.sassGeneratedCssTsValidationError != nil {
		this_.outer_.compiler.Report_1(sassGeneratedCssTsExpert.sassGeneratedCssTsValidationError)
		return false
	}
	if n.IsScript() {
		{
			r_9 := this_.traversalState
			r_9.inSassGeneratedCssTsScript = sassGeneratedCssTsExpert.hasSassGeneratedCssTsJsDoc
		}
		{
			r_10 := this_.traversalState
			r_10.cssClosureClassesQualifiedName = j2gort.NullString
		}
	} else if this_.traversalState.inSassGeneratedCssTsScript && sassGeneratedCssTsExpert.isCssClosureClassesAssignment {
		{
			r_11 := this_.traversalState
			r_11.cssClosureClassesQualifiedName = sassGeneratedCssTsExpert.cssClosureClassesQualifiedName
		}
	}
	return true
}

// Java method: GatherCssNamesTraversal.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ReplaceCssNames_GatherCssNamesTraversal) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if this_.outer_.priv_ReplaceCssNames_isGetCssNameCall_4831(n) {
		var getCssNameInstance *ReplaceCssNames_GetCssNameInstance = this_.outer_.priv_ReplaceCssNames_createGetCssNameInstance_4832(n, this_.traversalState.cssClosureClassesQualifiedName)
		_ = getCssNameInstance
		if getCssNameInstance.IsValid() {
			this_.listOfCssNameInstances.Add_1(j2gort.PtrToAny(getCssNameInstance))
		} else {
			this_.outer_.compiler.Report_1(getCssNameInstance.GetValidationError())
		}
	}
}

// Java method: GatherCssNamesTraversal.createSassGeneratedCssTsExpert(com.google.javascript.rhino.Node n)
func (this_ *ReplaceCssNames_GatherCssNamesTraversal) CreateSassGeneratedCssTsExpert(n pkg_rhino.Node) *ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert {
	_ = n
	var jsdoc *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = jsdoc
	var hasSassGeneratedCssTsJsDoc bool = ((jsdoc != nil) && jsdoc.IsSassGeneratedCssTs())
	_ = hasSassGeneratedCssTsJsDoc
	if !(hasSassGeneratedCssTsJsDoc) && !(this_.traversalState.inSassGeneratedCssTsScript) {
		return StaticGet_ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert_DEFAULT()
	}
	var sassGeneratedCssTsValidationError *JSError = nil
	_ = sassGeneratedCssTsValidationError
	if hasSassGeneratedCssTsJsDoc && (!(n.IsScript()) || this_.outer_.priv_ReplaceCssNames_isNotInCssClosureFile_4830(n)) {
		sassGeneratedCssTsValidationError = JSError_Make_3(n, StaticGet_ReplaceCssNames_UNEXPECTED_SASS_GENERATED_CSS_TS_ERROR(), []string{})
	}
	var isCssClosureClassesAssignment bool = false
	_ = isCssClosureClassesAssignment
	var cssClosureClassesQualifiedName string = j2gort.NullString
	_ = cssClosureClassesQualifiedName
	if this_.traversalState.inSassGeneratedCssTsScript && n.IsAssign() {
		var assignmentTarget pkg_rhino.Node = n.GetFirstChild()
		_ = assignmentTarget
		if assignmentTarget.IsGetProp() && (j2gort.CheckNotNullString(assignmentTarget.GetString()) == "classes") {
			isCssClosureClassesAssignment = true
			cssClosureClassesQualifiedName = assignmentTarget.GetQualifiedName()
		}
	}
	if (!(hasSassGeneratedCssTsJsDoc) && (sassGeneratedCssTsValidationError == nil)) && !(isCssClosureClassesAssignment) {
		return StaticGet_ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert_DEFAULT()
	}
	return NewReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert(hasSassGeneratedCssTsJsDoc, sassGeneratedCssTsValidationError, isCssClosureClassesAssignment, cssClosureClassesQualifiedName)
}

// Java class: com.google.javascript.jscomp.ReplaceCssNames.GatherCssNamesTraversal.SassGeneratedCssTsExpert
type ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert struct {
	sassGeneratedCssTsValidationError *JSError `gson:"sassGeneratedCssTsValidationError"`
	cssClosureClassesQualifiedName    string   `gson:"cssClosureClassesQualifiedName"`
	idHash_                           atomic.Int32
	hasSassGeneratedCssTsJsDoc        bool `gson:"hasSassGeneratedCssTsJsDoc"`
	isCssClosureClassesAssignment     bool `gson:"isCssClosureClassesAssignment"`
}

func (this_ *ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert) IsJavaNil_() bool {
	return this_ == nil
}
func (*ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceCssNames$GatherCssNamesTraversal$SassGeneratedCssTsExpert", "SassGeneratedCssTsExpert"
}
func (this_ *ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceCssNames.GatherCssNamesTraversal.SassGeneratedCssTsExpert", "com.google.javascript.jscomp.ReplaceCssNames$GatherCssNamesTraversal$SassGeneratedCssTsExpert", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert](v_)
	})
}

var static_ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert_DEFAULT *ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert

func StaticGet_ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert_DEFAULT() *ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert {
	clinit_ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert()
	return static_ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert_DEFAULT
}
func StaticSet_ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert_DEFAULT(v *ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert) {
	clinit_ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert()
	static_ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert_DEFAULT = v
}

var clinitDone_ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert uint32
var clinitState_ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert j2gort.ClinitState

func clinit_ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert() {
	if atomic.LoadUint32(&clinitDone_ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert) != 0 {
		return
	}
	clinitSlow_ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert()
}

func clinitSlow_ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert() {
	if !j2gort.EnterClinit(&clinitState_ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert, &clinitDone_ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert, &clinitDone_ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert)
	static_ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert_DEFAULT = NewReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert(false, nil, false, j2gort.NullString)
}

func (this_ *ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert) initFields_ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert() {
	_ = this_
	this_.cssClosureClassesQualifiedName = j2gort.NullString
}

// Java constructor: SassGeneratedCssTsExpert(boolean hasSassGeneratedCssTsJsDoc, com.google.javascript.jscomp.JSError sassGeneratedCssTsValidationError, boolean isCssClosureClassesAssignment, java.lang.String cssClosureClassesQualifiedName)
func NewReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert(hasSassGeneratedCssTsJsDoc bool, sassGeneratedCssTsValidationError *JSError, isCssClosureClassesAssignment bool, cssClosureClassesQualifiedName string) *ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert {
	clinit_ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert()
	this_ := &ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert{}
	InitReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert(this_, hasSassGeneratedCssTsJsDoc, sassGeneratedCssTsValidationError, isCssClosureClassesAssignment, cssClosureClassesQualifiedName)
	return this_
}

func InitReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert(this_ *ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert, hasSassGeneratedCssTsJsDoc bool, sassGeneratedCssTsValidationError *JSError, isCssClosureClassesAssignment bool, cssClosureClassesQualifiedName string) {
	_ = hasSassGeneratedCssTsJsDoc
	_ = sassGeneratedCssTsValidationError
	_ = isCssClosureClassesAssignment
	_ = cssClosureClassesQualifiedName
	this_.initFields_ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert()
	this_.hasSassGeneratedCssTsJsDoc = hasSassGeneratedCssTsJsDoc
	this_.sassGeneratedCssTsValidationError = sassGeneratedCssTsValidationError
	this_.isCssClosureClassesAssignment = isCssClosureClassesAssignment
	this_.cssClosureClassesQualifiedName = cssClosureClassesQualifiedName
}

// Java class: com.google.javascript.jscomp.ReplaceCssNames.GatherCssNamesTraversal.TraversalState
type ReplaceCssNames_GatherCssNamesTraversal_TraversalState struct {
	cssClosureClassesQualifiedName string `gson:"cssClosureClassesQualifiedName"`
	idHash_                        atomic.Int32
	inSassGeneratedCssTsScript     bool `gson:"inSassGeneratedCssTsScript"`
}

func (this_ *ReplaceCssNames_GatherCssNamesTraversal_TraversalState) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceCssNames_GatherCssNamesTraversal_TraversalState) IsJavaNil_() bool {
	return this_ == nil
}
func (*ReplaceCssNames_GatherCssNamesTraversal_TraversalState) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceCssNames$GatherCssNamesTraversal$TraversalState", "TraversalState"
}
func (this_ *ReplaceCssNames_GatherCssNamesTraversal_TraversalState) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReplaceCssNames_GatherCssNamesTraversal_TraversalState) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceCssNames.GatherCssNamesTraversal.TraversalState", "com.google.javascript.jscomp.ReplaceCssNames$GatherCssNamesTraversal$TraversalState", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceCssNames_GatherCssNamesTraversal_TraversalState](v_)
	})
}

func (this_ *ReplaceCssNames_GatherCssNamesTraversal_TraversalState) initFields_ReplaceCssNames_GatherCssNamesTraversal_TraversalState() {
	_ = this_
	this_.cssClosureClassesQualifiedName = j2gort.NullString
}

// Java constructor: TraversalState()
func NewReplaceCssNames_GatherCssNamesTraversal_TraversalState() *ReplaceCssNames_GatherCssNamesTraversal_TraversalState {
	this_ := &ReplaceCssNames_GatherCssNamesTraversal_TraversalState{}
	InitReplaceCssNames_GatherCssNamesTraversal_TraversalState(this_)
	return this_
}

func InitReplaceCssNames_GatherCssNamesTraversal_TraversalState(this_ *ReplaceCssNames_GatherCssNamesTraversal_TraversalState) {
	this_.initFields_ReplaceCssNames_GatherCssNamesTraversal_TraversalState()
}

// Java class: com.google.javascript.jscomp.ReplaceCssNames.GetCssNameInstance
type ReplaceCssNames_GetCssNameInstance struct {
	outer_                         *ReplaceCssNames
	validationError                *JSError       `gson:"validationError"`
	cssClosureClassesQualifiedName string         `gson:"cssClosureClassesQualifiedName"`
	callNode                       pkg_rhino.Node `gson:"callNode"`
	idHash_                        atomic.Int32
}

func (this_ *ReplaceCssNames_GetCssNameInstance) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceCssNames_GetCssNameInstance) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceCssNames_GetCssNameInstance) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceCssNames$GetCssNameInstance", "GetCssNameInstance"
}
func (this_ *ReplaceCssNames_GetCssNameInstance) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReplaceCssNames_GetCssNameInstance) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceCssNames.GetCssNameInstance", "com.google.javascript.jscomp.ReplaceCssNames$GetCssNameInstance", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceCssNames_GetCssNameInstance](v_)
	})
}
func (this_ *ReplaceCssNames_GetCssNameInstance) SetOuter_(o *ReplaceCssNames) { this_.outer_ = o }
func (this_ *ReplaceCssNames_GetCssNameInstance) GetOuter_() *ReplaceCssNames  { return this_.outer_ }

func (this_ *ReplaceCssNames_GetCssNameInstance) initFields_ReplaceCssNames_GetCssNameInstance() {
	_ = this_
	this_.cssClosureClassesQualifiedName = j2gort.NullString
}

// Java constructor: GetCssNameInstance(com.google.javascript.rhino.Node callNode, com.google.javascript.jscomp.JSError validationError, java.lang.String cssClosureClassesQualifiedName)
func NewReplaceCssNames_GetCssNameInstance(outer_ *ReplaceCssNames, callNode pkg_rhino.Node, validationError *JSError, cssClosureClassesQualifiedName string) *ReplaceCssNames_GetCssNameInstance {
	this_ := &ReplaceCssNames_GetCssNameInstance{}
	this_.outer_ = outer_
	InitReplaceCssNames_GetCssNameInstance(this_, callNode, validationError, cssClosureClassesQualifiedName)
	return this_
}

func InitReplaceCssNames_GetCssNameInstance(this_ *ReplaceCssNames_GetCssNameInstance, callNode pkg_rhino.Node, validationError *JSError, cssClosureClassesQualifiedName string) {
	_ = callNode
	_ = validationError
	_ = cssClosureClassesQualifiedName
	this_.initFields_ReplaceCssNames_GetCssNameInstance()
	this_.callNode = callNode
	this_.validationError = validationError
	this_.cssClosureClassesQualifiedName = cssClosureClassesQualifiedName
}

// Java method: GetCssNameInstance.isValid()
func (this_ *ReplaceCssNames_GetCssNameInstance) IsValid() bool {
	return (this_.validationError == nil)
}

// Java method: GetCssNameInstance.getValidationError()
func (this_ *ReplaceCssNames_GetCssNameInstance) GetValidationError() *JSError {
	j2gort.CheckNotNullPtr(this_.validationError)
	return this_.validationError
}

// Java method: GetCssNameInstance.replaceWithExpression()
func (this_ *ReplaceCssNames_GetCssNameInstance) ReplaceWithExpression() {
	if !(this_.IsValid()) {
		j2gort.CheckStateFmt(false, "not a valid goog.getCssName() call: %s", j2gort.BoxNode(this_.callNode))
	}
	j2gort.CheckNotNullNode(this_.callNode.GetParent())
	var childCount int32 = this_.callNode.GetChildCount()
	_ = childCount
	switch childCount {
	case int32(2):
		this_.priv_GetCssNameInstance_replaceWithConvertedSingleArg_4842()
	case int32(3):
		this_.priv_GetCssNameInstance_replaceWithConcatenatedArgs_4843()
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("invalid number of children: %s for: %s", []any{j2gort.ValueOfInt(childCount), j2gort.BoxNode(this_.callNode)})))))
	}
}

// Java method: GetCssNameInstance.replaceWithConvertedSingleArg()
func (this_ *ReplaceCssNames_GetCssNameInstance) priv_GetCssNameInstance_replaceWithConvertedSingleArg_4842() {
	var stringLitArg pkg_rhino.Node = j2gort.CheckNotNullNode(this_.callNode.GetSecondChild())
	_ = stringLitArg
	if !(stringLitArg.IsStringLit()) {
		j2gort.CheckStateFmt(false, "not a string literal: %s", j2gort.BoxNode(stringLitArg))
	}
	stringLitArg.Detach()
	this_.outer_.priv_ReplaceCssNames_processStringNode_4836(stringLitArg)
	this_.callNode.ReplaceWith(stringLitArg)
	this_.outer_.compiler.ReportChangeToEnclosingScope(stringLitArg)
}

// Java method: GetCssNameInstance.replaceWithConcatenatedArgs()
func (this_ *ReplaceCssNames_GetCssNameInstance) priv_GetCssNameInstance_replaceWithConcatenatedArgs_4843() {
	var firstArg pkg_rhino.Node = j2gort.CheckNotNullNode(this_.callNode.GetSecondChild())
	_ = firstArg
	var secondArg pkg_rhino.Node = j2gort.CheckNotNullNode(firstArg.GetNext())
	_ = secondArg
	firstArg.Detach()
	secondArg.Detach()
	this_.outer_.priv_ReplaceCssNames_processStringNode_4836(secondArg)
	secondArg.SetString(("-" + j2gort.StrValString(secondArg.GetString())))
	var replacement pkg_rhino.Node = this_.outer_.astFactory.CreateAdd(firstArg, secondArg).Srcref(this_.callNode)
	_ = replacement
	this_.callNode.ReplaceWith(replacement)
	this_.outer_.compiler.ReportChangeToEnclosingScope(replacement)
}

// Java method: GetCssNameInstance.isCssClosureFileClassesMember()
func (this_ *ReplaceCssNames_GetCssNameInstance) IsCssClosureFileClassesMember() bool {
	return (this_.cssClosureClassesQualifiedName != j2gort.NullString)
}

// Java method: GetCssNameInstance.getCssClosureClassesQualifiedName()
func (this_ *ReplaceCssNames_GetCssNameInstance) GetCssClosureClassesQualifiedName() string {
	j2gort.CheckNotNullString(this_.cssClosureClassesQualifiedName)
	return this_.cssClosureClassesQualifiedName
}

// Java method: GetCssNameInstance.getCssClosureClassesMemberName()
func (this_ *ReplaceCssNames_GetCssNameInstance) GetCssClosureClassesMemberName() string {
	j2gort.CheckNotNullString(this_.cssClosureClassesQualifiedName)

	{
		curAnc_ := this_.callNode

		if curAnc_ == 0 {
			j2gort.ThrowNPE()
		}
		curAnc_ = curAnc_.GetParent()
		for curAnc_ != 0 {
			ancestor := curAnc_
			curAnc_ = curAnc_.GetParent()

			_ = ancestor
			if ancestor.IsStringKey() {
				return (j2gort.StrValString(this_.cssClosureClassesQualifiedName) + "." + j2gort.StrValString(ancestor.GetString()))
			}
		}
	}
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("No css closure classes member name found: %s", []any{j2gort.BoxNode(this_.callNode)})))))
}

// Java method: GetCssNameInstance.getCssClassName()
func (this_ *ReplaceCssNames_GetCssNameInstance) GetCssClassName() string {
	if !(this_.IsValid()) {
		j2gort.CheckStateFmt(false, "not a valid goog.getCssName() call: %s", j2gort.BoxNode(this_.callNode))
	}
	j2gort.CheckNotNullNode(this_.callNode.GetParent())
	var childCount int32 = this_.callNode.GetChildCount()
	_ = childCount
	var firstArg pkg_rhino.Node = j2gort.CheckNotNullNode(this_.callNode.GetSecondChild())
	_ = firstArg
	switch childCount {
	case int32(2):
		if !(firstArg.IsStringLit()) {
			j2gort.CheckStateFmt(false, "not a string literal: %s", j2gort.BoxNode(firstArg))
		}
		return firstArg.GetString()
	case int32(3):
		var secondArg pkg_rhino.Node = j2gort.CheckNotNullNode(firstArg.GetNext())
		_ = secondArg
		if !(secondArg.IsStringLit()) {
			j2gort.CheckStateFmt(false, "not a string literal: %s", j2gort.BoxNode(secondArg))
		}
		return secondArg.GetString()
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("invalid number of children: %s for: %s", []any{j2gort.ValueOfInt(childCount), j2gort.BoxNode(this_.callNode)})))))
	}
	panic("unreachable switch expression")
}

// Java class: com.google.javascript.jscomp.ReplaceCssNames.CountCssNamesBySymbol
type ReplaceCssNames_CountCssNamesBySymbol struct {
	outer_  *ReplaceCssNames
	idHash_ atomic.Int32
}

func (this_ *ReplaceCssNames_CountCssNamesBySymbol) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceCssNames_CountCssNamesBySymbol) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceCssNames_CountCssNamesBySymbol) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceCssNames$CountCssNamesBySymbol", "CountCssNamesBySymbol"
}
func (this_ *ReplaceCssNames_CountCssNamesBySymbol) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReplaceCssNames_CountCssNamesBySymbol) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceCssNames.CountCssNamesBySymbol", "com.google.javascript.jscomp.ReplaceCssNames$CountCssNamesBySymbol", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceCssNames_CountCssNamesBySymbol](v_)
	})
}
func (this_ *ReplaceCssNames_CountCssNamesBySymbol) SetOuter_(o *ReplaceCssNames)        { this_.outer_ = o }
func (this_ *ReplaceCssNames_CountCssNamesBySymbol) GetOuter_() *ReplaceCssNames         { return this_.outer_ }
func (*ReplaceCssNames_CountCssNamesBySymbol) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *ReplaceCssNames_CountCssNamesBySymbol) initFields_ReplaceCssNames_CountCssNamesBySymbol() {
	_ = this_
}

// Java constructor: CountCssNamesBySymbol()
func NewReplaceCssNames_CountCssNamesBySymbol(outer_ *ReplaceCssNames) *ReplaceCssNames_CountCssNamesBySymbol {
	this_ := &ReplaceCssNames_CountCssNamesBySymbol{}
	this_.outer_ = outer_
	InitReplaceCssNames_CountCssNamesBySymbol(this_)
	return this_
}

func InitReplaceCssNames_CountCssNamesBySymbol(this_ *ReplaceCssNames_CountCssNamesBySymbol) {
	this_.initFields_ReplaceCssNames_CountCssNamesBySymbol()
}

// Java method: CountCssNamesBySymbol.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ReplaceCssNames_CountCssNamesBySymbol) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if n.IsScript() {
		return this_.outer_.priv_ReplaceCssNames_isNotInCssClosureFile_4830(n)
	} else {
		return true
	}
}

// Java method: CountCssNamesBySymbol.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ReplaceCssNames_CountCssNamesBySymbol) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if this_.outer_.classesObjectsQualifiedNames.IsEmpty() {
		return
	}
	if n.IsGetProp() && !(func() bool {
		sr_13 := "classes"
		se_14 := n.GetString()
		return j2gort.CheckNotNullString(sr_13) == se_14
	}()) {
		return
	}
	var classesObjectQualifiedName string = n.GetQualifiedName()
	_ = classesObjectQualifiedName
	if (classesObjectQualifiedName == j2gort.NullString) || !(this_.outer_.classesObjectsQualifiedNames.Contains(classesObjectQualifiedName)) {
		return
	}
	var classNameNode pkg_rhino.Node = n.GetParent()
	_ = classNameNode
	if !(n.IsGetProp()) || !(classNameNode.IsGetProp()) {
		this_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_ReplaceCssNames_INVALID_USE_OF_CLASSES_OBJECT_ERROR(), []string{}))
		return
	}
	var classQualifiedName string = classNameNode.GetQualifiedName()
	_ = classQualifiedName
	if (classQualifiedName == j2gort.NullString) || !(this_.outer_.cssNamesBySymbol.ContainsKey(classQualifiedName)) {
		this_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_ReplaceCssNames_UNKNOWN_SYMBOL_ERROR(), []string{classNameNode.GetString()}))
		return
	}
	var className string = j2gort.OrderedMapGetStr(this_.outer_.cssNamesBySymbol, classQualifiedName)
	_ = className
	this_.outer_.cssNameCollector.Add(className)
}
func (*ReplaceCssNames_FindSetCssNameTraversal) AlwaysShouldTraverse_()      {}
func (*ReplaceCssNames_FindSetCssNameTraversal) CanReleaseSyntacticScopes_() {}
func (*ReplaceCssNames_GatherCssNamesTraversal) CanReleaseSyntacticScopes_() {}
func (*ReplaceCssNames_CountCssNamesBySymbol) CanReleaseSyntacticScopes_()   {}
