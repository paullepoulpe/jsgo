/*
 * Copyright 2014 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ReportUntranspilableFeatures.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_parsing "github.com/paullepoulpe/jsgo/transpiled/gen/parsing"
	pkg_regex "github.com/paullepoulpe/jsgo/transpiled/gen/regex"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)


var lambdaSingleton_report_untranspilable_features_1 j2gort.Ext_java_util_function_Predicate
var lambdaSingleton_report_untranspilable_features_2 j2gort.Ext_java_util_function_Predicate
var lambdaSingleton_report_untranspilable_features_3 j2gort.Ext_java_util_function_Predicate
func init() {
	lambdaSingleton_report_untranspilable_features_1 = j2gort.NewLambda_Ext_java_util_function_Predicate(func(lp_4 any) bool {
	t := j2gort.AsIface[pkg_regex.RegExpTree](lp_4)
	_ = t
	return j2gort.InstanceOf[*pkg_regex.RegExpTree_LookbehindAssertion](t)
})
	lambdaSingleton_report_untranspilable_features_2 = j2gort.NewLambda_Ext_java_util_function_Predicate(func(lp_5 any) bool {
	t := j2gort.AsIface[pkg_regex.RegExpTree](lp_5)
	_ = t
	return j2gort.InstanceOf[*pkg_regex.RegExpTree_NamedCaptureGroup](t)
})
	lambdaSingleton_report_untranspilable_features_3 = j2gort.NewLambda_Ext_java_util_function_Predicate(func(lp_6 any) bool {
	t := j2gort.AsIface[pkg_regex.RegExpTree](lp_6)
	_ = t
	return j2gort.InstanceOf[*pkg_regex.RegExpTree_UnicodePropertyEscape](t)
})
}
// Java class: com.google.javascript.jscomp.ReportUntranspilableFeatures
type ReportUntranspilableFeatures struct {
	AbstractPeepholeTranspilationBase
	compiler AbstractCompiler `gson:"compiler"`
	browserFeatureSetYear *CompilerOptions_BrowserFeaturesetYear `gson:"browserFeatureSetYear"`
	untranspilableFeaturesToRemove *pkg_parser.FeatureSet `gson:"untranspilableFeaturesToRemove"`
}

func (this_ *ReportUntranspilableFeatures) BindSelf_(s any) {
	_ = s
}
func (this_ *ReportUntranspilableFeatures) IsJavaNil_() bool { return this_ == nil }
func (*ReportUntranspilableFeatures) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.ReportUntranspilableFeatures", "ReportUntranspilableFeatures" }
func (this_ *ReportUntranspilableFeatures) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.AbstractPeepholeTranspilationBase.IdentityHash_()
}
func (*ReportUntranspilableFeatures) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReportUntranspilableFeatures", "com.google.javascript.jscomp.ReportUntranspilableFeatures", func(v_ any) bool {
		return j2gort.InstanceOf[*ReportUntranspilableFeatures](v_)
	})
}
func (*ReportUntranspilableFeatures) Implements_jscomp_AbstractPeepholeTranspilation() {}

var static_ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT *DiagnosticType
func StaticGet_ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT() *DiagnosticType {
	clinit_ReportUntranspilableFeatures()
	return static_ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT
}
func StaticSet_ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT(v *DiagnosticType) {
	clinit_ReportUntranspilableFeatures()
	static_ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT = v
}
var static_ReportUntranspilableFeatures_UNSUPPORTED_FEATURE_PRESENT *DiagnosticType
func StaticGet_ReportUntranspilableFeatures_UNSUPPORTED_FEATURE_PRESENT() *DiagnosticType {
	clinit_ReportUntranspilableFeatures()
	return static_ReportUntranspilableFeatures_UNSUPPORTED_FEATURE_PRESENT
}
func StaticSet_ReportUntranspilableFeatures_UNSUPPORTED_FEATURE_PRESENT(v *DiagnosticType) {
	clinit_ReportUntranspilableFeatures()
	static_ReportUntranspilableFeatures_UNSUPPORTED_FEATURE_PRESENT = v
}
var static_ReportUntranspilableFeatures_UNTRANSPILABLE_ES5_FEATURES *pkg_parser.FeatureSet
func StaticGet_ReportUntranspilableFeatures_UNTRANSPILABLE_ES5_FEATURES() *pkg_parser.FeatureSet {
	clinit_ReportUntranspilableFeatures()
	return static_ReportUntranspilableFeatures_UNTRANSPILABLE_ES5_FEATURES
}
func StaticSet_ReportUntranspilableFeatures_UNTRANSPILABLE_ES5_FEATURES(v *pkg_parser.FeatureSet) {
	clinit_ReportUntranspilableFeatures()
	static_ReportUntranspilableFeatures_UNTRANSPILABLE_ES5_FEATURES = v
}
var static_ReportUntranspilableFeatures_UNTRANSPILABLE_2018_FEATURES *pkg_parser.FeatureSet
func StaticGet_ReportUntranspilableFeatures_UNTRANSPILABLE_2018_FEATURES() *pkg_parser.FeatureSet {
	clinit_ReportUntranspilableFeatures()
	return static_ReportUntranspilableFeatures_UNTRANSPILABLE_2018_FEATURES
}
func StaticSet_ReportUntranspilableFeatures_UNTRANSPILABLE_2018_FEATURES(v *pkg_parser.FeatureSet) {
	clinit_ReportUntranspilableFeatures()
	static_ReportUntranspilableFeatures_UNTRANSPILABLE_2018_FEATURES = v
}
var static_ReportUntranspilableFeatures_UNTRANSPILABLE_2019_FEATURES *pkg_parser.FeatureSet
func StaticGet_ReportUntranspilableFeatures_UNTRANSPILABLE_2019_FEATURES() *pkg_parser.FeatureSet {
	clinit_ReportUntranspilableFeatures()
	return static_ReportUntranspilableFeatures_UNTRANSPILABLE_2019_FEATURES
}
func StaticSet_ReportUntranspilableFeatures_UNTRANSPILABLE_2019_FEATURES(v *pkg_parser.FeatureSet) {
	clinit_ReportUntranspilableFeatures()
	static_ReportUntranspilableFeatures_UNTRANSPILABLE_2019_FEATURES = v
}
var static_ReportUntranspilableFeatures_UNTRANSPILABLE_2020_FEATURES *pkg_parser.FeatureSet
func StaticGet_ReportUntranspilableFeatures_UNTRANSPILABLE_2020_FEATURES() *pkg_parser.FeatureSet {
	clinit_ReportUntranspilableFeatures()
	return static_ReportUntranspilableFeatures_UNTRANSPILABLE_2020_FEATURES
}
func StaticSet_ReportUntranspilableFeatures_UNTRANSPILABLE_2020_FEATURES(v *pkg_parser.FeatureSet) {
	clinit_ReportUntranspilableFeatures()
	static_ReportUntranspilableFeatures_UNTRANSPILABLE_2020_FEATURES = v
}
var static_ReportUntranspilableFeatures_UNTRANSPILABLE_2022_FEATURES *pkg_parser.FeatureSet
func StaticGet_ReportUntranspilableFeatures_UNTRANSPILABLE_2022_FEATURES() *pkg_parser.FeatureSet {
	clinit_ReportUntranspilableFeatures()
	return static_ReportUntranspilableFeatures_UNTRANSPILABLE_2022_FEATURES
}
func StaticSet_ReportUntranspilableFeatures_UNTRANSPILABLE_2022_FEATURES(v *pkg_parser.FeatureSet) {
	clinit_ReportUntranspilableFeatures()
	static_ReportUntranspilableFeatures_UNTRANSPILABLE_2022_FEATURES = v
}
var static_ReportUntranspilableFeatures_ALL_UNTRANSPILABLE_FEATURES *pkg_parser.FeatureSet
func StaticGet_ReportUntranspilableFeatures_ALL_UNTRANSPILABLE_FEATURES() *pkg_parser.FeatureSet {
	clinit_ReportUntranspilableFeatures()
	return static_ReportUntranspilableFeatures_ALL_UNTRANSPILABLE_FEATURES
}
func StaticSet_ReportUntranspilableFeatures_ALL_UNTRANSPILABLE_FEATURES(v *pkg_parser.FeatureSet) {
	clinit_ReportUntranspilableFeatures()
	static_ReportUntranspilableFeatures_ALL_UNTRANSPILABLE_FEATURES = v
}
var clinitDone_ReportUntranspilableFeatures uint32
var clinitState_ReportUntranspilableFeatures j2gort.ClinitState
func clinit_ReportUntranspilableFeatures() {
	if atomic.LoadUint32(&clinitDone_ReportUntranspilableFeatures) != 0 {
		return
	}
	clinitSlow_ReportUntranspilableFeatures()
}

func clinitSlow_ReportUntranspilableFeatures() {
	if !j2gort.EnterClinit(&clinitState_ReportUntranspilableFeatures, &clinitDone_ReportUntranspilableFeatures) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ReportUntranspilableFeatures, &clinitDone_ReportUntranspilableFeatures)
	static_ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT = DiagnosticType_Error("JSC_UNTRANSPILABLE", "Cannot convert feature \"{0}\" to targeted output language. Feature requires at minimum {1}.{2}")
	static_ReportUntranspilableFeatures_UNSUPPORTED_FEATURE_PRESENT = DiagnosticType_Error("JSC_UNSUPPORTED", "The feature \"{0}\" is currently unsupported for transpilation.")
	static_ReportUntranspilableFeatures_UNTRANSPILABLE_ES5_FEATURES = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature_7186([]*pkg_parser.FeatureSet_Feature{pkg_parser.StaticGet_FeatureSet_Feature_GETTER(), pkg_parser.StaticGet_FeatureSet_Feature_SETTER()})
	static_ReportUntranspilableFeatures_UNTRANSPILABLE_2018_FEATURES = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature_7186([]*pkg_parser.FeatureSet_Feature{pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_FLAG_S(), pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_LOOKBEHIND(), pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_NAMED_GROUPS(), pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_UNICODE_PROPERTY_ESCAPE()})
	static_ReportUntranspilableFeatures_UNTRANSPILABLE_2019_FEATURES = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_UNESCAPED_UNICODE_LINE_OR_PARAGRAPH_SEP())
	static_ReportUntranspilableFeatures_UNTRANSPILABLE_2020_FEATURES = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_BIGINT())
	static_ReportUntranspilableFeatures_UNTRANSPILABLE_2022_FEATURES = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature_7186([]*pkg_parser.FeatureSet_Feature{pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_FLAG_D(), pkg_parser.StaticGet_FeatureSet_Feature_TOP_LEVEL_AWAIT()})
	static_ReportUntranspilableFeatures_ALL_UNTRANSPILABLE_FEATURES = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().Union(StaticGet_ReportUntranspilableFeatures_UNTRANSPILABLE_ES5_FEATURES()).Union(StaticGet_ReportUntranspilableFeatures_UNTRANSPILABLE_2018_FEATURES()).Union(StaticGet_ReportUntranspilableFeatures_UNTRANSPILABLE_2019_FEATURES()).Union(StaticGet_ReportUntranspilableFeatures_UNTRANSPILABLE_2022_FEATURES()).Union(StaticGet_ReportUntranspilableFeatures_UNTRANSPILABLE_2020_FEATURES())
}

func (this_ *ReportUntranspilableFeatures) initFields_ReportUntranspilableFeatures() {
	_ = this_
}

// Java constructor: ReportUntranspilableFeatures(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.CompilerOptions.@org.jspecify.annotations.Nullable BrowserFeaturesetYear browserFeatureSetYear, com.google.javascript.jscomp.parsing.parser.FeatureSet outputFeatures)
func NewReportUntranspilableFeatures(compiler AbstractCompiler, browserFeatureSetYear *CompilerOptions_BrowserFeaturesetYear, outputFeatures *pkg_parser.FeatureSet) *ReportUntranspilableFeatures {
	clinit_ReportUntranspilableFeatures()
	this_ := &ReportUntranspilableFeatures{}
	InitReportUntranspilableFeatures(this_, compiler, browserFeatureSetYear, outputFeatures)
	return this_
}

func InitReportUntranspilableFeatures(this_ *ReportUntranspilableFeatures, compiler AbstractCompiler, browserFeatureSetYear *CompilerOptions_BrowserFeaturesetYear, outputFeatures *pkg_parser.FeatureSet) {
	_ = compiler
	_ = browserFeatureSetYear
	_ = outputFeatures
	InitAbstractPeepholeTranspilation(&this_.AbstractPeepholeTranspilationBase)
	this_.initFields_ReportUntranspilableFeatures()
	j2gort.CheckNotNullIface(compiler)
	j2gort.CheckNotNullPtr(outputFeatures)
	this_.compiler = compiler
	this_.browserFeatureSetYear = browserFeatureSetYear
	this_.untranspilableFeaturesToRemove = StaticGet_ReportUntranspilableFeatures_ALL_UNTRANSPILABLE_FEATURES().Without_1(outputFeatures)
}

// Java method: ReportUntranspilableFeatures.getTranspiledAwayFeatures()
func (this_ *ReportUntranspilableFeatures) GetTranspiledAwayFeatures() *pkg_parser.FeatureSet {
	return this_.untranspilableFeaturesToRemove
}

// Java method: ReportUntranspilableFeatures.getAdditionalFeaturesToRunOn()
func (this_ *ReportUntranspilableFeatures) GetAdditionalFeaturesToRunOn() *pkg_parser.FeatureSet {
	return pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_SYNTAX())
}

// Java method: ReportUntranspilableFeatures.checkForUntranspilable(com.google.javascript.rhino.Node root)
func (this_ *ReportUntranspilableFeatures) priv_ReportUntranspilableFeatures_checkForUntranspilable_4895(root pkg_rhino.Node) {
	_ = root
	switch root.TokenSwitchOrd_() {
	case int32(42) /* Token.REGEXP */:
		var pattern string = root.GetFirstChild().GetString()
		_ = pattern
		var flags string
		_ = flags
		if root.HasTwoChildren() {
			flags = root.GetLastChild().GetString()
		} else {
			flags = ""
		}
		var reg pkg_regex.RegExpTree = nil
		_ = reg
		var tryFlow_1 int
		func() {
			defer func() {
				if rec_2 := recover(); rec_2 != nil {
					jex_3 := javalang.RecoverException(rec_2)
					_ = jex_3
					if j2gort.InstanceOf[javalang.IIllegalArgumentException](jex_3) || j2gort.InstanceOf[javalang.IIndexOutOfBoundsException](jex_3) {
						ex := jex_3
						_ = ex
						this_.compiler.Report_1(JSError_Make_3(root, StaticGet_CheckRegExp_MALFORMED_REGEXP(), []string{j2gort.ThrowableGetMessage(ex)}))
						tryFlow_1 = 1
						return
					} else {
						panic(rec_2)
					}
				}
			}()
			reg = pkg_regex.RegExpTree_ParseRegExp(pattern, flags)
		}()
		if tryFlow_1 == 1 { break }
		if (this_.untranspilableFeaturesToRemove.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_FLAG_S())) {
			this_.priv_ReportUntranspilableFeatures_checkForRegExpSFlag_4897(root)
		}
		if (this_.untranspilableFeaturesToRemove.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_LOOKBEHIND())) {
			this_.priv_ReportUntranspilableFeatures_checkForLookbehind_4898(root, reg)
		}
		if (this_.untranspilableFeaturesToRemove.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_NAMED_GROUPS())) {
			this_.priv_ReportUntranspilableFeatures_checkForNamedGroups_4899(root, reg)
		}
		if (this_.untranspilableFeaturesToRemove.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_UNICODE_PROPERTY_ESCAPE())) {
			this_.priv_ReportUntranspilableFeatures_checkForUnicodePropertyEscape_4900(root, reg)
		}
		if (this_.untranspilableFeaturesToRemove.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_FLAG_D())) {
			this_.priv_ReportUntranspilableFeatures_checkForRegExpDFlag_4901(root)
		}
	case int32(34) /* Token.BIGINT */:
		if (this_.untranspilableFeaturesToRemove.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_BIGINT())) {
			this_.priv_ReportUntranspilableFeatures_reportUntranspilable_4896(pkg_parser.StaticGet_FeatureSet_Feature_BIGINT(), root)
		}
	case int32(95) /* Token.GETTER_DEF */:
		if (this_.untranspilableFeaturesToRemove.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_GETTER())) {
			this_.priv_ReportUntranspilableFeatures_reportUntranspilable_4896(pkg_parser.StaticGet_FeatureSet_Feature_GETTER(), root)
		}
	case int32(96) /* Token.SETTER_DEF */:
		if (this_.untranspilableFeaturesToRemove.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_SETTER())) {
			this_.priv_ReportUntranspilableFeatures_reportUntranspilable_4896(pkg_parser.StaticGet_FeatureSet_Feature_SETTER(), root)
		}
	case int32(115) /* Token.AWAIT */, int32(113) /* Token.FOR_AWAIT_OF */:
		if (this_.untranspilableFeaturesToRemove.Contains_1_FeatureSet_Feature(pkg_parser.StaticGet_FeatureSet_Feature_TOP_LEVEL_AWAIT())) {
			this_.priv_ReportUntranspilableFeatures_checkForTopLevelAwait_4902(root)
		}
	default:
	}
}

// Java method: ReportUntranspilableFeatures.reportUntranspilable(com.google.javascript.jscomp.parsing.parser.FeatureSet.Feature feature, com.google.javascript.rhino.Node node)
func (this_ *ReportUntranspilableFeatures) priv_ReportUntranspilableFeatures_reportUntranspilable_4896(feature *pkg_parser.FeatureSet_Feature, node pkg_rhino.Node) {
	_ = feature
	_ = node
	var minimumLanguageMode *pkg_parsing.Config_LanguageMode = pkg_parsing.Config_LanguageMode_MinimumRequiredFor(feature)
	_ = minimumLanguageMode
	if ((minimumLanguageMode == pkg_parsing.StaticGet_Config_LanguageMode_UNSUPPORTED())) {
		this_.compiler.Report_1(JSError_Make_3(node, StaticGet_ReportUntranspilableFeatures_UNSUPPORTED_FEATURE_PRESENT(), []string{feature.String()}))
		return
	}
	var minimum string = j2gort.NullString
	_ = minimum
	var suggestion string = " Consider targeting a more modern output."
	_ = suggestion
	if ((this_.browserFeatureSetYear != nil)) {
		var minimumYear *CompilerOptions_BrowserFeaturesetYear = CompilerOptions_BrowserFeaturesetYear_MinimumRequiredFor(feature)
		_ = minimumYear
		if ((minimumYear == nil)) {
			minimum = (j2gort.StrVal(minimumLanguageMode) + ", which is not yet supported by any browser featureset year")
			suggestion = ""
		} else {
			minimum = ("browser featureset year " + j2gort.StrValInt32(minimumYear.GetYear()))
			suggestion = (j2gort.StrValString(suggestion) + "\nCurrent browser featureset year: " + j2gort.StrValInt32(this_.browserFeatureSetYear.GetYear()))
		}
	} else {
		minimum = minimumLanguageMode.String()
	}
	this_.compiler.Report_1(JSError_Make_3(node, StaticGet_ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT(), []string{feature.String(), minimum, suggestion}))
}

// Java method: ReportUntranspilableFeatures.checkForRegExpSFlag(com.google.javascript.rhino.Node regexpNode)
func (this_ *ReportUntranspilableFeatures) priv_ReportUntranspilableFeatures_checkForRegExpSFlag_4897(regexpNode pkg_rhino.Node) {
	_ = regexpNode
	if !(regexpNode.IsRegExp()) { j2gort.CheckArgument(false) }
	var flags string
	_ = flags
	if regexpNode.HasTwoChildren() {
		flags = regexpNode.GetLastChild().GetString()
	} else {
		flags = ""
	}
	if ((j2gort.JStrIndexOf(flags, "s") >= 0)) {
		this_.priv_ReportUntranspilableFeatures_reportUntranspilable_4896(pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_FLAG_S(), regexpNode)
	}
}

// Java method: ReportUntranspilableFeatures.checkForLookbehind(com.google.javascript.rhino.Node regexpNode, com.google.javascript.jscomp.regex.RegExpTree tree)
func (this_ *ReportUntranspilableFeatures) priv_ReportUntranspilableFeatures_checkForLookbehind_4898(regexpNode pkg_rhino.Node, tree pkg_regex.RegExpTree) {
	_ = regexpNode
	_ = tree
	if !((regexpNode != 0)) { j2gort.CheckArgument(false) }
	if (ReportUntranspilableFeatures_AnySubtreeMeetsPredicate(tree, lambdaSingleton_report_untranspilable_features_1)) {
		this_.priv_ReportUntranspilableFeatures_reportUntranspilable_4896(pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_LOOKBEHIND(), regexpNode)
	}
}

// Java method: ReportUntranspilableFeatures.checkForNamedGroups(com.google.javascript.rhino.Node regexpNode, com.google.javascript.jscomp.regex.RegExpTree tree)
func (this_ *ReportUntranspilableFeatures) priv_ReportUntranspilableFeatures_checkForNamedGroups_4899(regexpNode pkg_rhino.Node, tree pkg_regex.RegExpTree) {
	_ = regexpNode
	_ = tree
	if !((regexpNode != 0)) { j2gort.CheckArgument(false) }
	if (ReportUntranspilableFeatures_AnySubtreeMeetsPredicate(tree, lambdaSingleton_report_untranspilable_features_2)) {
		this_.priv_ReportUntranspilableFeatures_reportUntranspilable_4896(pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_NAMED_GROUPS(), regexpNode)
	}
}

// Java method: ReportUntranspilableFeatures.checkForUnicodePropertyEscape(com.google.javascript.rhino.Node regexpNode, com.google.javascript.jscomp.regex.RegExpTree tree)
func (this_ *ReportUntranspilableFeatures) priv_ReportUntranspilableFeatures_checkForUnicodePropertyEscape_4900(regexpNode pkg_rhino.Node, tree pkg_regex.RegExpTree) {
	_ = regexpNode
	_ = tree
	if !((regexpNode != 0)) { j2gort.CheckArgument(false) }
	if (ReportUntranspilableFeatures_AnySubtreeMeetsPredicate(tree, lambdaSingleton_report_untranspilable_features_3)) {
		this_.priv_ReportUntranspilableFeatures_reportUntranspilable_4896(pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_UNICODE_PROPERTY_ESCAPE(), regexpNode)
	}
}

// Java method: ReportUntranspilableFeatures.checkForRegExpDFlag(com.google.javascript.rhino.Node regexpNode)
func (this_ *ReportUntranspilableFeatures) priv_ReportUntranspilableFeatures_checkForRegExpDFlag_4901(regexpNode pkg_rhino.Node) {
	_ = regexpNode
	if !(regexpNode.IsRegExp()) { j2gort.CheckArgument(false) }
	var flags string
	_ = flags
	if regexpNode.HasTwoChildren() {
		flags = regexpNode.GetLastChild().GetString()
	} else {
		flags = ""
	}
	if ((j2gort.JStrIndexOf(flags, "d") >= 0)) {
		this_.priv_ReportUntranspilableFeatures_reportUntranspilable_4896(pkg_parser.StaticGet_FeatureSet_Feature_REGEXP_FLAG_D(), regexpNode)
	}
}

// Java method: ReportUntranspilableFeatures.checkForTopLevelAwait(com.google.javascript.rhino.Node awaitNode)
func (this_ *ReportUntranspilableFeatures) priv_ReportUntranspilableFeatures_checkForTopLevelAwait_4902(awaitNode pkg_rhino.Node) {
	_ = awaitNode
	if !((awaitNode.IsAwait() || awaitNode.IsForAwaitOf())) { j2gort.CheckArgument(false) }
	if ((NodeUtil_GetEnclosingFunction(awaitNode) == 0)) {
		if !(NodeUtil_GetEnclosingScript(awaitNode).GetBooleanProp(pkg_rhino.StaticGet_Node_ES6_MODULE())) { j2gort.CheckArgumentStr(false, "Top-level await is only allowed in ES module sources") }
		this_.priv_ReportUntranspilableFeatures_reportUntranspilable_4896(pkg_parser.StaticGet_FeatureSet_Feature_TOP_LEVEL_AWAIT(), awaitNode)
	}
}

// Java static method: ReportUntranspilableFeatures.anySubtreeMeetsPredicate(com.google.javascript.jscomp.regex.RegExpTree tree, java.util.function.Predicate<com.google.javascript.jscomp.regex.RegExpTree> p)
func ReportUntranspilableFeatures_AnySubtreeMeetsPredicate(tree pkg_regex.RegExpTree, p j2gort.Ext_java_util_function_Predicate) bool {
	clinit_ReportUntranspilableFeatures()
	_ = tree
	_ = p
	if (p.Test(tree)) {
		return true
	}
	{
		cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_List](tree.Children())))
		for cur_7.HasNext() {
			rawElem_8 := cur_7.Next()
			subTree := j2gort.AsIface[pkg_regex.RegExpTree](rawElem_8)
			_ = subTree
			if (ReportUntranspilableFeatures_AnySubtreeMeetsPredicate(subTree, p)) {
				return true
			}
		}
	}
	return false
}

// Java method: ReportUntranspilableFeatures.transpileSubtree(com.google.javascript.rhino.Node subtree)
func (this_ *ReportUntranspilableFeatures) TranspileSubtree(subtree pkg_rhino.Node) pkg_rhino.Node {
	_ = subtree
	this_.priv_ReportUntranspilableFeatures_checkForUntranspilable_4895(subtree)
	return subtree
}

