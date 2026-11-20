/*
 * Copyright 2010 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ReplaceStrings.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_replace_strings_1 j2gort.Ext_com_google_common_base_Predicate

func init() {
	lambdaSingleton_replace_strings_1 = j2gort.NewLambda_Ext_com_google_common_base_Predicate(func(lp_1 any) bool {
		result := j2gort.CastPtr[ReplaceStrings_Result](lp_1)
		_ = result
		return result.didReplacement
	})
}

// Java class: com.google.javascript.jscomp.ReplaceStrings
type ReplaceStrings struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	placeholderToken string                                            `gson:"placeholderToken"`
	compiler         AbstractCompiler                                  `gson:"compiler"`
	functions        j2gort.OrderedMap[string, *ReplaceStrings_Config] `gson:"functions"`
	nameGenerator    *DefaultNameGenerator                             `gson:"nameGenerator"`
	results          j2gort.OrderedMap[string, *ReplaceStrings_Result] `gson:"results"`
}

func (this_ *ReplaceStrings) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceStrings) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceStrings) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceStrings", "ReplaceStrings"
}
func (this_ *ReplaceStrings) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*ReplaceStrings) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceStrings", "com.google.javascript.jscomp.ReplaceStrings", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceStrings](v_)
	})
}
func (*ReplaceStrings) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*ReplaceStrings) Implements_jscomp_CompilerPass()                            {}
func (*ReplaceStrings) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_ReplaceStrings_BAD_REPLACEMENT_CONFIGURATION *DiagnosticType

func StaticGet_ReplaceStrings_BAD_REPLACEMENT_CONFIGURATION() *DiagnosticType {
	clinit_ReplaceStrings()
	return static_ReplaceStrings_BAD_REPLACEMENT_CONFIGURATION
}
func StaticSet_ReplaceStrings_BAD_REPLACEMENT_CONFIGURATION(v *DiagnosticType) {
	clinit_ReplaceStrings()
	static_ReplaceStrings_BAD_REPLACEMENT_CONFIGURATION = v
}

var static_ReplaceStrings_STRING_REPLACEMENT_TAGGED_TEMPLATE *DiagnosticType

func StaticGet_ReplaceStrings_STRING_REPLACEMENT_TAGGED_TEMPLATE() *DiagnosticType {
	clinit_ReplaceStrings()
	return static_ReplaceStrings_STRING_REPLACEMENT_TAGGED_TEMPLATE
}
func StaticSet_ReplaceStrings_STRING_REPLACEMENT_TAGGED_TEMPLATE(v *DiagnosticType) {
	clinit_ReplaceStrings()
	static_ReplaceStrings_STRING_REPLACEMENT_TAGGED_TEMPLATE = v
}
func StaticGet_ReplaceStrings_DEFAULT_PLACEHOLDER_TOKEN() string { return "`" }
func StaticGet_ReplaceStrings_EXCLUSION_PREFIX() string          { return ":!" }
func StaticGet_ReplaceStrings_REPLACE_ONE_MARKER() string        { return "?" }
func StaticGet_ReplaceStrings_REPLACE_ALL_MARKER() string        { return "*" }

var static_ReplaceStrings_USED_RESULTS j2gort.Ext_com_google_common_base_Predicate

func StaticGet_ReplaceStrings_USED_RESULTS() j2gort.Ext_com_google_common_base_Predicate {
	clinit_ReplaceStrings()
	return static_ReplaceStrings_USED_RESULTS
}
func StaticSet_ReplaceStrings_USED_RESULTS(v j2gort.Ext_com_google_common_base_Predicate) {
	clinit_ReplaceStrings()
	static_ReplaceStrings_USED_RESULTS = v
}

var clinitDone_ReplaceStrings uint32
var clinitState_ReplaceStrings j2gort.ClinitState

func clinit_ReplaceStrings() {
	if atomic.LoadUint32(&clinitDone_ReplaceStrings) != 0 {
		return
	}
	clinitSlow_ReplaceStrings()
}

func clinitSlow_ReplaceStrings() {
	if !j2gort.EnterClinit(&clinitState_ReplaceStrings, &clinitDone_ReplaceStrings) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ReplaceStrings, &clinitDone_ReplaceStrings)
	static_ReplaceStrings_BAD_REPLACEMENT_CONFIGURATION = DiagnosticType_Error("JSC_BAD_REPLACEMENT_CONFIGURATION", "Bad replacement configuration \"{0}\": {1}")
	static_ReplaceStrings_STRING_REPLACEMENT_TAGGED_TEMPLATE = DiagnosticType_Warning("JSC_STRING_REPLACEMENT_TAGGED_TEMPLATE", "Cannot string-replace arguments of a template literal tag function.")
	static_ReplaceStrings_USED_RESULTS = lambdaSingleton_replace_strings_1
}

func (this_ *ReplaceStrings) initFields_ReplaceStrings() {
	_ = this_
	this_.placeholderToken = j2gort.NullString
	this_.functions = j2gort.NewOrderedMap[string, *ReplaceStrings_Config]()
	this_.results = j2gort.NewOrderedMap[string, *ReplaceStrings_Result]()
}

// Java constructor: ReplaceStrings(com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String placeholderToken, java.util.List<java.lang.String> functionsToInspect)
func NewReplaceStrings(compiler AbstractCompiler, placeholderToken string, functionsToInspect j2gort.Ext_java_util_List) *ReplaceStrings {
	clinit_ReplaceStrings()
	this_ := &ReplaceStrings{}
	InitReplaceStrings(this_, compiler, placeholderToken, functionsToInspect)
	return this_
}

func InitReplaceStrings(this_ *ReplaceStrings, compiler AbstractCompiler, placeholderToken string, functionsToInspect j2gort.Ext_java_util_List) {
	_ = compiler
	_ = placeholderToken
	_ = functionsToInspect
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_ReplaceStrings()
	this_.compiler = compiler
	if j2gort.CheckNotNullString(placeholderToken) == "" {
		this_.placeholderToken = "`"
	} else {
		this_.placeholderToken = placeholderToken
	}
	this_.nameGenerator = ReplaceStrings_CreateNameGenerator()
	this_.priv_ReplaceStrings_parseConfiguration_4926(functionsToInspect)
}

// Java method: ReplaceStrings.getResult()
func (this_ *ReplaceStrings) GetResult() j2gort.Ext_java_util_List {
	return j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_lang_Iterable(j2gort.Ext_com_google_common_collect_Iterables_Filter_2_Ext_java_lang_Iterable_Ext_com_google_common_base_Predicate(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.results.LiveValues()), StaticGet_ReplaceStrings_USED_RESULTS())))
}

// Java method: ReplaceStrings.getStringMap()
func (this_ *ReplaceStrings) GetStringMap() *VariableMap {
	var map_ j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
	_ = map_
	{
		cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.Ext_com_google_common_collect_Iterables_Filter_2_Ext_java_lang_Iterable_Ext_com_google_common_base_Predicate(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.results.LiveValues()), StaticGet_ReplaceStrings_USED_RESULTS())))
		for cur_2.HasNext() {
			rawElem_3 := cur_2.Next()
			result := j2gort.CastPtr[ReplaceStrings_Result](rawElem_3)
			_ = result
			map_.Put_2(j2gort.BoxString(result.replacement), j2gort.BoxString(result.original))
		}
	}
	return NewVariableMap(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](map_.BuildOrThrow()))
}

// Java method: ReplaceStrings.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ReplaceStrings) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}

// Java method: ReplaceStrings.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ReplaceStrings) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(23) /* Token.NEW */, int32(28) /* Token.CALL */, int32(130) /* Token.TAGGED_TEMPLATELIT */ :
		var calledFn pkg_rhino.Node = n.GetFirstChild()
		_ = calledFn
		var name string = calledFn.GetOriginalQualifiedName()
		_ = name
		if name != j2gort.NullString {
			var config *ReplaceStrings_Config = this_.priv_ReplaceStrings_findMatching_4920(name, n.GetSourceFileName())
			_ = config
			if config != nil {
				this_.priv_ReplaceStrings_doSubstitutions_4921(t, config, n)
				return
			}
		}
	default:
	}
}

// Java method: ReplaceStrings.findMatching(java.lang.String name, java.lang.String callsiteSourceFileName)
func (this_ *ReplaceStrings) priv_ReplaceStrings_findMatching_4920(name string, callsiteSourceFileName string) *ReplaceStrings_Config {
	_ = name
	_ = callsiteSourceFileName
	var config *ReplaceStrings_Config = this_.functions.Get(name)
	_ = config
	if config == nil {
		name = j2gort.JStrReplaceChar(name, uint16(36), uint16(46))
		config = this_.functions.Get(name)
	}
	if config != nil {
		{
			cur_4 := j2gort.NewIterCursor(j2gort.CheckNotNull(config.excludedFilenameSuffixes))
			for cur_4.HasNext() {
				excludedSuffix := cur_4.NextStr()
				_ = excludedSuffix
				if javalang.StringEndsWith(j2gort.CheckNotNullString(callsiteSourceFileName), j2gort.CheckNotNullString(excludedSuffix)) {
					return nil
				}
			}
		}
	}
	return config
}

// Java method: ReplaceStrings.doSubstitutions(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.ReplaceStrings.Config config, com.google.javascript.rhino.Node n)
func (this_ *ReplaceStrings) priv_ReplaceStrings_doSubstitutions_4921(t *NodeTraversal, config *ReplaceStrings_Config, n pkg_rhino.Node) {
	_ = t
	_ = config
	_ = n
	if n.IsTaggedTemplateLit() {
		this_.compiler.Report_1(JSError_Make_3(n, StaticGet_ReplaceStrings_STRING_REPLACEMENT_TAGGED_TEMPLATE(), []string{}))
		return
	}
	if !(n.IsNew() || n.IsCall()) {
		j2gort.CheckState(false)
	}
	if !(config.IsReplaceAll()) {
		{
			cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(config.parameters))
			for cur_5.HasNext() {
				rawElem_6 := cur_5.Next()
				parameter := javalang.D2I(j2gort.UnboxToDouble(rawElem_6))
				_ = parameter
				var arg pkg_rhino.Node = n.GetChildAtIndex(parameter)
				_ = arg
				if arg != 0 {
					this_.priv_ReplaceStrings_replaceExpression_4922(t, arg)
				}
			}
		}
	} else {
		var firstParam pkg_rhino.Node = n.GetSecondChild()
		_ = firstParam
		{
			var arg_2 pkg_rhino.Node = firstParam
			_ = arg_2
			for ; arg_2 != 0; arg_2 = arg_2.GetNext() {
				arg_2 = this_.priv_ReplaceStrings_replaceExpression_4922(t, arg_2)
			}
		}
	}
}

// Java method: ReplaceStrings.replaceExpression(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node expr)
func (this_ *ReplaceStrings) priv_ReplaceStrings_replaceExpression_4922(t *NodeTraversal, expr pkg_rhino.Node) pkg_rhino.Node {
	_ = t
	_ = expr
	var replacement pkg_rhino.Node = pkg_rhino.Node(0)
	_ = replacement
	var key string = j2gort.NullString
	_ = key
	var replacementString string = j2gort.NullString
	_ = replacementString
	switch expr.TokenSwitchOrd_() {
	case int32(35) /* Token.STRINGLIT */ :
		key = expr.GetString()
		replacementString = this_.priv_ReplaceStrings_getReplacement_4923(key)
		replacement = pkg_rhino.IR_String(replacementString)
	case int32(131) /* Token.TEMPLATELIT */, int32(13) /* Token.ADD */, int32(32) /* Token.NAME */ :
		var keyBuilder *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
		_ = keyBuilder
		var keyNode pkg_rhino.Node = pkg_rhino.IR_String("")
		_ = keyNode
		replacement = this_.priv_ReplaceStrings_buildReplacement_4925(t, expr, keyNode, keyBuilder)
		key = keyBuilder.String()
		{
			var iifeCond_1 bool
			{

				sr_7 := key
				se_8 := this_.placeholderToken
				iifeCond_1 = j2gort.CheckNotNullString(sr_7) == se_8
			}
			if iifeCond_1 {
				return expr
			}
		}
		replacementString = this_.priv_ReplaceStrings_getReplacement_4923(key)
		keyNode.SetString(replacementString)
	default:
		return expr
	}
	j2gort.CheckNotNullString(key)
	j2gort.CheckNotNullString(replacementString)
	this_.priv_ReplaceStrings_recordReplacement_4924(key)
	replacement.SrcrefTreeIfMissing(expr)
	expr.ReplaceWith(replacement)
	t.ReportCodeChange_0()
	return replacement
}

// Java method: ReplaceStrings.getReplacement(java.lang.String key)
func (this_ *ReplaceStrings) priv_ReplaceStrings_getReplacement_4923(key string) string {
	_ = key
	var result *ReplaceStrings_Result = this_.results.Get(key)
	_ = result
	if result != nil {
		return result.replacement
	}
	var replacement string = this_.nameGenerator.GenerateNextName()
	_ = replacement
	result = NewReplaceStrings_Result(key, replacement)
	this_.results.Put(key, result)
	return replacement
}

// Java method: ReplaceStrings.recordReplacement(java.lang.String key)
func (this_ *ReplaceStrings) priv_ReplaceStrings_recordReplacement_4924(key string) {
	_ = key
	var result *ReplaceStrings_Result = this_.results.Get(key)
	_ = result
	if !(result != nil) {
		j2gort.CheckState(false)
	}
	{
		r_9 := result
		r_9.didReplacement = true
	}
}

// Java method: ReplaceStrings.buildReplacement(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node expr, com.google.javascript.rhino.Node prefix, java.lang.StringBuilder keyBuilder)
func (this_ *ReplaceStrings) priv_ReplaceStrings_buildReplacement_4925(t *NodeTraversal, expr pkg_rhino.Node, prefix pkg_rhino.Node, keyBuilder *j2gort.Ext_java_lang_StringBuilder) pkg_rhino.Node {
	_ = t
	_ = expr
	_ = prefix
	_ = keyBuilder
	{
		var left pkg_rhino.Node = pkg_rhino.Node(0)
		_ = left
		var right pkg_rhino.Node = pkg_rhino.Node(0)
		_ = right
		var var_ *Var = nil
		_ = var_
		switch expr.TokenSwitchOrd_() {
		case int32(13) /* Token.ADD */ :
			left = expr.GetFirstChild()
			right = left.GetNext()
			prefix = this_.priv_ReplaceStrings_buildReplacement_4925(t, left, prefix, keyBuilder)
			return this_.priv_ReplaceStrings_buildReplacement_4925(t, right, prefix, keyBuilder)
		case int32(131) /* Token.TEMPLATELIT */ :
			{
				var child pkg_rhino.Node = expr.GetFirstChild()
				_ = child
				for ; child != 0; child = child.GetNext() {
					switch child.TokenSwitchOrd_() {
					case int32(133) /* Token.TEMPLATELIT_STRING */ :
						keyBuilder.Append_1_string(child.GetCookedString())
					case int32(132) /* Token.TEMPLATELIT_SUB */ :
						prefix = this_.priv_ReplaceStrings_buildReplacement_4925(t, child.GetFirstChild(), prefix, keyBuilder)
					default:
						panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected TEMPLATELIT child: "+j2gort.StrVal(child))))
					}
				}
			}
			return prefix
		case int32(35) /* Token.STRINGLIT */ :
			keyBuilder.Append_1_string(expr.GetString())
			return prefix
		case int32(32) /* Token.NAME */ :
			var_ = j2gort.CastPtr[Var](t.GetScope().GetVar(expr.GetString()))
			if (var_ != nil) && (var_.IsDeclaredOrInferredConst() || var_.IsConst()) {
				var initialValue pkg_rhino.Node = var_.GetInitialValue()
				_ = initialValue
				if initialValue != 0 {
					var newKeyNode pkg_rhino.Node = pkg_rhino.IR_String("")
					_ = newKeyNode
					var newKeyBuilder *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
					_ = newKeyBuilder
					var replacement pkg_rhino.Node = this_.priv_ReplaceStrings_buildReplacement_4925(t, initialValue, newKeyNode, newKeyBuilder)
					_ = replacement
					if replacement == newKeyNode {
						keyBuilder.Append_1_any(newKeyBuilder)
						return prefix
					}
				}
			}
			fallthrough
		default:
			keyBuilder.Append_1_string(this_.placeholderToken)
			prefix = pkg_rhino.IR_Add(prefix, pkg_rhino.IR_String(this_.placeholderToken))
			return pkg_rhino.IR_Add(prefix, expr.CloneTree_0())
		}
	}
	panic("unreachable")
}

// Java method: ReplaceStrings.parseConfiguration(java.util.List<java.lang.String> functionsToInspect)
func (this_ *ReplaceStrings) priv_ReplaceStrings_parseConfiguration_4926(functionsToInspect j2gort.Ext_java_util_List) {
	_ = functionsToInspect
	{
		cur_10 := j2gort.NewIterCursor(j2gort.CheckNotNull(functionsToInspect))
		for cur_10.HasNext() {
			function := cur_10.NextStr()
			_ = function
			var config *ReplaceStrings_Config = this_.priv_ReplaceStrings_parseConfiguration_4927(function)
			_ = config
			if config != nil {
				this_.functions.Put(config.name, config)
			}
		}
	}
}

// Java method: ReplaceStrings.parseConfiguration(java.lang.String function)
func (this_ *ReplaceStrings) priv_ReplaceStrings_parseConfiguration_4927(function string) *ReplaceStrings_Config {
	_ = function
	var first int32 = j2gort.JStrIndexOfCp(function, int32(uint16(40)))
	_ = first
	var last int32 = j2gort.JStrIndexOfCp(function, int32(uint16(41)))
	_ = last
	var colon int32 = j2gort.JStrIndexOf(function, ":!")
	_ = colon
	if !((first != int32(-1)) && (last != int32(-1))) {
		j2gort.CheckState(false)
	}
	if j2gort.JStrIndexOf(function, ".prototype.") >= 0 {
		this_.compiler.Report_1(JSError_Make_2(StaticGet_ReplaceStrings_BAD_REPLACEMENT_CONFIGURATION(), []string{function, "Cannot replace strings passed to prototype methods."}))
		return nil
	}
	var name string = j2gort.JStrSub(function, int32(0), first)
	_ = name
	var params string = j2gort.JStrSub(function, (first + int32(1)), last)
	_ = params
	var paramCount int32 = int32(0)
	_ = paramCount
	var replacementParameters j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = replacementParameters
	var parts j2gort.Ext_java_util_List = j2gort.SplitterSplitToListStr(j2gort.Ext_com_google_common_base_Splitter_On_1_uint16(uint16(44)), params)
	_ = parts
	{
		cur_11 := j2gort.NewIterCursor(j2gort.CheckNotNull(parts))
		for cur_11.HasNext() {
			param := cur_11.NextStr()
			_ = param
			paramCount = paramCount + 1
			{
				var iifeCond_2 bool
				{

					sr_12 := param
					se_13 := "*"
					iifeCond_2 = j2gort.CheckNotNullString(sr_12) == se_13
				}
				if iifeCond_2 {
					if !((paramCount == int32(1)) && (parts.Size() == int32(1))) {
						j2gort.CheckState(false)
					}
					replacementParameters.Add_1(j2gort.ValueOfInt(int32(0)))
				} else {
					var iifeCond_1 bool
					{

						sr_14 := param
						se_15 := "?"
						iifeCond_1 = j2gort.CheckNotNullString(sr_14) == se_15
					}
					if iifeCond_1 {
						if !(!(replacementParameters.Contains(j2gort.ValueOfInt(int32(0))))) {
							j2gort.CheckState(false)
						}
						replacementParameters.Add_1(j2gort.ValueOfInt(paramCount))
					} else {
						if !(j2gort.CheckNotNullString(param) == "") {
							j2gort.CheckStateFmt(false, "Unknown marker (%s)", j2gort.BoxString(param))
						}
					}
				}
			}
		}
	}
	if !(!(replacementParameters.IsEmpty())) {
		j2gort.CheckState(false)
	}
	var iifeArg_2 j2gort.
		Ext_com_google_common_collect_ImmutableSet
iifeSwitch_2:
	switch {
	default:

		if colon == int32(-1) {
			iifeArg_2 = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
			break iifeSwitch_2
		}
		var iifeArg_1 []any
	iifeSwitch_1:
		switch {
		default:

			arr_16 := j2gort.JStrSplit(j2gort.JStrSubFrom(function, (colon+j2gort.JStrLen(":!"))), ",")
			if arr_16 == nil {
				iifeArg_1 = nil
				break iifeSwitch_1
			}
			out_17 := make([]any, len(arr_16))
			for i_, el_18 := range arr_16 {
				out_17[i_] = j2gort.BoxString(el_18)
			}
			iifeArg_1 = out_17
			break iifeSwitch_1
		}
		iifeArg_2 = j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_slice_any(iifeArg_1)
		break iifeSwitch_2
	}
	return NewReplaceStrings_Config(name, replacementParameters, iifeArg_2)
}

// Java static method: ReplaceStrings.createNameGenerator()
func ReplaceStrings_CreateNameGenerator() *DefaultNameGenerator {
	clinit_ReplaceStrings()
	var namePrefix string = ""
	_ = namePrefix
	var reservedChars j2gort.Ext_com_google_common_collect_ImmutableSet = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
	_ = reservedChars
	return NewDefaultNameGenerator_3(j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()), namePrefix, j2gort.AsIface[j2gort.Ext_java_util_Set](reservedChars))
}

// Java class: com.google.javascript.jscomp.ReplaceStrings.Config
type ReplaceStrings_Config struct {
	name                     string                                            `gson:"name"`
	parameters               j2gort.Ext_java_util_List                         `gson:"parameters"`
	excludedFilenameSuffixes j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"excludedFilenameSuffixes"`
	idHash_                  atomic.Int32
}

func (this_ *ReplaceStrings_Config) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceStrings_Config) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceStrings_Config) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceStrings$Config", "Config"
}
func (this_ *ReplaceStrings_Config) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReplaceStrings_Config) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceStrings.Config", "com.google.javascript.jscomp.ReplaceStrings$Config", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceStrings_Config](v_)
	})
}

func StaticGet_ReplaceStrings_Config_REPLACE_ALL_VALUE() int32 { return int32(0) }
func (this_ *ReplaceStrings_Config) initFields_ReplaceStrings_Config() {
	_ = this_
	this_.name = j2gort.NullString
}

// Java constructor: Config(java.lang.String name, java.util.List<java.lang.Integer> replacementParameters, com.google.common.collect.ImmutableSet<java.lang.String> excludedFilenameSuffixes)
func NewReplaceStrings_Config(name string, replacementParameters j2gort.Ext_java_util_List, excludedFilenameSuffixes j2gort.Ext_com_google_common_collect_ImmutableSet) *ReplaceStrings_Config {
	this_ := &ReplaceStrings_Config{}
	InitReplaceStrings_Config(this_, name, replacementParameters, excludedFilenameSuffixes)
	return this_
}

func InitReplaceStrings_Config(this_ *ReplaceStrings_Config, name string, replacementParameters j2gort.Ext_java_util_List, excludedFilenameSuffixes j2gort.Ext_com_google_common_collect_ImmutableSet) {
	_ = name
	_ = replacementParameters
	_ = excludedFilenameSuffixes
	this_.initFields_ReplaceStrings_Config()
	this_.name = name
	this_.parameters = replacementParameters
	this_.excludedFilenameSuffixes = excludedFilenameSuffixes
}

// Java method: Config.isReplaceAll()
func (this_ *ReplaceStrings_Config) IsReplaceAll() bool {
	return ((this_.parameters.Size() == int32(1)) && this_.parameters.Contains(j2gort.ValueOfInt(int32(0))))
}

// Java class: com.google.javascript.jscomp.ReplaceStrings.Result
type ReplaceStrings_Result struct {
	original       string `gson:"original"`
	replacement    string `gson:"replacement"`
	idHash_        atomic.Int32
	didReplacement bool `gson:"didReplacement"`
}

func (this_ *ReplaceStrings_Result) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceStrings_Result) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceStrings_Result) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceStrings$Result", "Result"
}
func (this_ *ReplaceStrings_Result) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReplaceStrings_Result) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceStrings.Result", "com.google.javascript.jscomp.ReplaceStrings$Result", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceStrings_Result](v_)
	})
}

func (this_ *ReplaceStrings_Result) initFields_ReplaceStrings_Result() {
	_ = this_
	this_.original = j2gort.NullString
	this_.replacement = j2gort.NullString

}

// Java constructor: Result(java.lang.String original, java.lang.String replacement)
func NewReplaceStrings_Result(original string, replacement string) *ReplaceStrings_Result {
	this_ := &ReplaceStrings_Result{}
	InitReplaceStrings_Result(this_, original, replacement)
	return this_
}

func InitReplaceStrings_Result(this_ *ReplaceStrings_Result, original string, replacement string) {
	_ = original
	_ = replacement
	this_.initFields_ReplaceStrings_Result()
	this_.original = original
	this_.replacement = replacement
}
func (*ReplaceStrings) AlwaysShouldTraverse_() {}
