// Port of com/google/javascript/jscomp/ReportUntranspilableFeatures.java @ closure-compiler v20261015

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

package jscomp

import (
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/regex"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ReportUntranspilableFeatures looks for presence of features that are not
// supported for transpilation (mostly new RegExp features and bigint
// literal).
//
// Reports errors for any features are present in the root and not present
// in the targeted output language.
//
// Note: the errors reported by this pass are suppressible. In most cases
// the right thing to do if suppressing them is instead to bump the output
// language level; for example, someone who only needs to support ES2020 but
// is setting --language_out=ECMASCRIPT_2015 will see errors on ES2018 regex
// literal syntax. The best solution is to set
// --language_out=ECMASCRIPT_2020. A second solution, which we support but
// do not recommend, is to suppress these errors & let the compiler silently
// output untranspiled ES2018 regex literals.
type ReportUntranspilableFeatures struct {
	AbstractPeepholeTranspilationBase
	compiler                       AbstractCompiler
	browserFeatureSetYear          *CompilerOptions_BrowserFeaturesetYear
	untranspilableFeaturesToRemove *parser.FeatureSet
}

// Java: ReportUntranspilableFeatures.UNTRANSPILABLE_FEATURE_PRESENT
var ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT = DiagnosticType_Error(
	"JSC_UNTRANSPILABLE",
	"Cannot convert feature \"{0}\" to targeted output language. Feature requires at minimum"+
		" {1}.{2}")

// Java: ReportUntranspilableFeatures.UNSUPPORTED_FEATURE_PRESENT
var ReportUntranspilableFeatures_UNSUPPORTED_FEATURE_PRESENT = DiagnosticType_Error(
	"JSC_UNSUPPORTED", "The feature \"{0}\" is currently unsupported for transpilation.")

// Java: ReportUntranspilableFeatures.UNTRANSPILABLE_ES5_FEATURES
var reportUntranspilableFeatures_UNTRANSPILABLE_ES5_FEATURES = parser.FeatureSet_BARE_MINIMUM.WithFeatures(parser.FeatureSet_Feature_GETTER, parser.FeatureSet_Feature_SETTER)

// Java: ReportUntranspilableFeatures.UNTRANSPILABLE_2018_FEATURES
var reportUntranspilableFeatures_UNTRANSPILABLE_2018_FEATURES = parser.FeatureSet_BARE_MINIMUM.WithFeatures(
	parser.FeatureSet_Feature_REGEXP_FLAG_S,
	parser.FeatureSet_Feature_REGEXP_LOOKBEHIND,
	parser.FeatureSet_Feature_REGEXP_NAMED_GROUPS,
	parser.FeatureSet_Feature_REGEXP_UNICODE_PROPERTY_ESCAPE)

// Java: ReportUntranspilableFeatures.UNTRANSPILABLE_2019_FEATURES
var reportUntranspilableFeatures_UNTRANSPILABLE_2019_FEATURES = parser.FeatureSet_BARE_MINIMUM.With(
	// We could transpile this, but there's no point. We always escape these in the output,
	// no need to have a separate pass to escape them. So we'll piggy back off this pass to
	// mark it as transpiled. Note that we never complain that this feature won't be
	// transpiled below.
	parser.FeatureSet_Feature_UNESCAPED_UNICODE_LINE_OR_PARAGRAPH_SEP)

// Java: ReportUntranspilableFeatures.UNTRANSPILABLE_2020_FEATURES
var reportUntranspilableFeatures_UNTRANSPILABLE_2020_FEATURES = parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_BIGINT)

// Java: ReportUntranspilableFeatures.UNTRANSPILABLE_2022_FEATURES
var reportUntranspilableFeatures_UNTRANSPILABLE_2022_FEATURES = parser.FeatureSet_BARE_MINIMUM.WithFeatures(parser.FeatureSet_Feature_REGEXP_FLAG_D, parser.FeatureSet_Feature_TOP_LEVEL_AWAIT)

// Java: ReportUntranspilableFeatures.ALL_UNTRANSPILABLE_FEATURES
var reportUntranspilableFeatures_ALL_UNTRANSPILABLE_FEATURES = parser.FeatureSet_BARE_MINIMUM.
	Union(reportUntranspilableFeatures_UNTRANSPILABLE_ES5_FEATURES).
	Union(reportUntranspilableFeatures_UNTRANSPILABLE_2018_FEATURES).
	Union(reportUntranspilableFeatures_UNTRANSPILABLE_2019_FEATURES).
	Union(reportUntranspilableFeatures_UNTRANSPILABLE_2022_FEATURES).
	Union(reportUntranspilableFeatures_UNTRANSPILABLE_2020_FEATURES)

// NewReportUntranspilableFeatures takes the compiler instance, the given
// BrowserFeaturesetYear (or nil if it was not specified) and the set of
// features supported in the given output language.
//
// Java: ReportUntranspilableFeatures(AbstractCompiler, BrowserFeaturesetYear, FeatureSet)
func NewReportUntranspilableFeatures(
	compiler AbstractCompiler,
	browserFeatureSetYear *CompilerOptions_BrowserFeaturesetYear,
	outputFeatures *parser.FeatureSet) *ReportUntranspilableFeatures {
	jutil.CheckNotNull(compiler)
	jutil.CheckNotNull(outputFeatures)
	r := &ReportUntranspilableFeatures{}
	r.compiler = compiler
	r.browserFeatureSetYear = browserFeatureSetYear
	// Features that we can't transpile...
	// and do not exist in the output language features
	r.untranspilableFeaturesToRemove =
		reportUntranspilableFeatures_ALL_UNTRANSPILABLE_FEATURES.Without(outputFeatures)
	return r
}

// Java: ReportUntranspilableFeatures.getTranspiledAwayFeatures()
func (r *ReportUntranspilableFeatures) GetTranspiledAwayFeatures() *parser.FeatureSet {
	return r.untranspilableFeaturesToRemove
}

// Java: ReportUntranspilableFeatures.getAdditionalFeaturesToRunOn()
func (r *ReportUntranspilableFeatures) GetAdditionalFeaturesToRunOn() *parser.FeatureSet {
	/*
	 * This pass needs to run on all ES3 REGEXP_SYNTAX because it checks for the presence of
	 * non-flag RegExp features that are not parsed and not attached to nodes.
	 */
	return parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_REGEXP_SYNTAX)
}

// Java: ReportUntranspilableFeatures.checkForUntranspilable(Node)
func (r *ReportUntranspilableFeatures) CheckForUntranspilable(root *rhino.Node) {
	// Non-flag RegExp features are not attached to nodes, so we must force traversal.
	switch root.GetToken() {
	case rhino.Token_REGEXP:
		pattern := root.GetFirstChild().GetString()
		flags := ""
		if root.HasTwoChildren() {
			flags = root.GetLastChild().GetString()
		}
		var reg regex.RegExpTree
		malformed := false
		jutil.Try(func() {
			reg = regex.RegExpTree_ParseRegExp(pattern, flags)
		}, func(ex jutil.IRuntimeException) {
			// Go-only: language-adapter — Java's multi-catch
			// `catch (IllegalArgumentException | IndexOutOfBoundsException ex)`:
			// any other exception propagates unchanged.
			switch ex.(type) {
			case jutil.IIllegalArgumentException, jutil.IIndexOutOfBoundsException:
				// A null message formats as "null", as MessageFormat does.
				r.compiler.Report(JSError_MakeWithNode(root, CheckRegExp_MALFORMED_REGEXP, jutil.ToJavaString(jutil.GetMessageOrNull(ex))))
				malformed = true
			default:
				panic(ex)
			}
		})
		if malformed {
			break
		}
		if r.untranspilableFeaturesToRemove.ContainsFeature(parser.FeatureSet_Feature_REGEXP_FLAG_S) {
			r.CheckForRegExpSFlag(root)
		}
		if r.untranspilableFeaturesToRemove.ContainsFeature(parser.FeatureSet_Feature_REGEXP_LOOKBEHIND) {
			r.CheckForLookbehind(root, reg)
		}
		if r.untranspilableFeaturesToRemove.ContainsFeature(parser.FeatureSet_Feature_REGEXP_NAMED_GROUPS) {
			r.CheckForNamedGroups(root, reg)
		}
		if r.untranspilableFeaturesToRemove.ContainsFeature(parser.FeatureSet_Feature_REGEXP_UNICODE_PROPERTY_ESCAPE) {
			r.CheckForUnicodePropertyEscape(root, reg)
		}
		if r.untranspilableFeaturesToRemove.ContainsFeature(parser.FeatureSet_Feature_REGEXP_FLAG_D) {
			r.CheckForRegExpDFlag(root)
		}
	case rhino.Token_BIGINT:
		// Transpilation of BigInt is not supported
		if r.untranspilableFeaturesToRemove.ContainsFeature(parser.FeatureSet_Feature_BIGINT) {
			r.ReportUntranspilable(parser.FeatureSet_Feature_BIGINT, root)
		}
	case rhino.Token_GETTER_DEF:
		if r.untranspilableFeaturesToRemove.ContainsFeature(parser.FeatureSet_Feature_GETTER) {
			r.ReportUntranspilable(parser.FeatureSet_Feature_GETTER, root)
		}
	case rhino.Token_SETTER_DEF:
		if r.untranspilableFeaturesToRemove.ContainsFeature(parser.FeatureSet_Feature_SETTER) {
			r.ReportUntranspilable(parser.FeatureSet_Feature_SETTER, root)
		}
	case rhino.Token_AWAIT, rhino.Token_FOR_AWAIT_OF:
		if r.untranspilableFeaturesToRemove.ContainsFeature(parser.FeatureSet_Feature_TOP_LEVEL_AWAIT) {
			r.CheckForTopLevelAwait(root)
		}
	default:
	}
}

// Java: ReportUntranspilableFeatures.reportUntranspilable(Feature, Node)
func (r *ReportUntranspilableFeatures) ReportUntranspilable(feature parser.FeatureSet_Feature, node *rhino.Node) {
	minimumLanguageMode := parsing.Config_LanguageMode_MinimumRequiredFor(feature)
	if minimumLanguageMode == parsing.Config_LanguageMode_UNSUPPORTED {
		r.compiler.Report(JSError_MakeWithNode(node, ReportUntranspilableFeatures_UNSUPPORTED_FEATURE_PRESENT, feature.String()))
		return
	}

	// The compiler always has an output featureset configured. Sometimes this is configured
	// directly or via a LanguageMode, and in this case browserFeatureSetYear is null. Otherwise if
	// browserFeatureSetYear is not null, the user configured a browserFeatureSetYear specifically.
	// Report an error message based on the actual API the user invoked: sometimes the year in a
	// browser featureset year (e.g. 2020) does not correspond to e.g. ECMASCRIPT_2020.
	var minimum string
	suggestion := " Consider targeting a more modern output."
	if r.browserFeatureSetYear != nil {
		minimumYear := CompilerOptions_BrowserFeaturesetYear_MinimumRequiredFor(feature)
		if minimumYear == nil {
			minimum =
				minimumLanguageMode.String() + ", which is not yet supported by any browser featureset year"
			suggestion = ""
		} else {
			minimum = "browser featureset year " + strconv.Itoa(minimumYear.GetYear())
			suggestion =
				suggestion + "\nCurrent browser featureset year: " + strconv.Itoa(r.browserFeatureSetYear.GetYear())
		}
	} else {
		minimum = minimumLanguageMode.String()
	}
	r.compiler.Report(
		JSError_MakeWithNode(
			node, ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT, feature.String(), minimum, suggestion))
}

// Java: ReportUntranspilableFeatures.checkForRegExpSFlag(Node)
func (r *ReportUntranspilableFeatures) CheckForRegExpSFlag(regexpNode *rhino.Node) {
	jutil.CheckArgument(regexpNode.IsRegExp())
	flags := ""
	if regexpNode.HasTwoChildren() {
		flags = regexpNode.GetLastChild().GetString()
	}
	if strings.Contains(flags, "s") {
		r.ReportUntranspilable(parser.FeatureSet_Feature_REGEXP_FLAG_S, regexpNode)
	}
}

// Java: ReportUntranspilableFeatures.checkForLookbehind(Node, RegExpTree)
func (r *ReportUntranspilableFeatures) CheckForLookbehind(regexpNode *rhino.Node, tree regex.RegExpTree) {
	jutil.CheckArgument(regexpNode != nil)
	if ReportUntranspilableFeatures_AnySubtreeMeetsPredicate(tree, func(t regex.RegExpTree) bool {
		_, ok := t.(*regex.RegExpTree_LookbehindAssertion)
		return ok
	}) {
		r.ReportUntranspilable(parser.FeatureSet_Feature_REGEXP_LOOKBEHIND, regexpNode)
	}
}

// Java: ReportUntranspilableFeatures.checkForNamedGroups(Node, RegExpTree)
func (r *ReportUntranspilableFeatures) CheckForNamedGroups(regexpNode *rhino.Node, tree regex.RegExpTree) {
	jutil.CheckArgument(regexpNode != nil)
	if ReportUntranspilableFeatures_AnySubtreeMeetsPredicate(tree, func(t regex.RegExpTree) bool {
		_, ok := t.(*regex.RegExpTree_NamedCaptureGroup)
		return ok
	}) {
		r.ReportUntranspilable(parser.FeatureSet_Feature_REGEXP_NAMED_GROUPS, regexpNode)
	}
}

// Java: ReportUntranspilableFeatures.checkForUnicodePropertyEscape(Node, RegExpTree)
func (r *ReportUntranspilableFeatures) CheckForUnicodePropertyEscape(regexpNode *rhino.Node, tree regex.RegExpTree) {
	jutil.CheckArgument(regexpNode != nil)
	if ReportUntranspilableFeatures_AnySubtreeMeetsPredicate(tree, func(t regex.RegExpTree) bool {
		_, ok := t.(*regex.RegExpTree_UnicodePropertyEscape)
		return ok
	}) {
		r.ReportUntranspilable(parser.FeatureSet_Feature_REGEXP_UNICODE_PROPERTY_ESCAPE, regexpNode)
	}
}

// Java: ReportUntranspilableFeatures.checkForRegExpDFlag(Node)
func (r *ReportUntranspilableFeatures) CheckForRegExpDFlag(regexpNode *rhino.Node) {
	jutil.CheckArgument(regexpNode.IsRegExp())
	flags := ""
	if regexpNode.HasTwoChildren() {
		flags = regexpNode.GetLastChild().GetString()
	}
	if strings.Contains(flags, "d") {
		r.ReportUntranspilable(parser.FeatureSet_Feature_REGEXP_FLAG_D, regexpNode)
	}
}

// Java: ReportUntranspilableFeatures.checkForTopLevelAwait(Node)
func (r *ReportUntranspilableFeatures) CheckForTopLevelAwait(awaitNode *rhino.Node) {
	jutil.CheckArgument(awaitNode.IsAwait() || awaitNode.IsForAwaitOf())

	if NodeUtil_GetEnclosingFunction(awaitNode) == nil {
		jutil.CheckArgumentMsg(
			NodeUtil_GetEnclosingScript(awaitNode).GetBooleanProp(rhino.Node_ES6_MODULE),
			"Top-level await is only allowed in ES module sources")
		r.ReportUntranspilable(parser.FeatureSet_Feature_TOP_LEVEL_AWAIT, awaitNode)
	}
}

// Java: ReportUntranspilableFeatures.anySubtreeMeetsPredicate(RegExpTree, Predicate<RegExpTree>)
func ReportUntranspilableFeatures_AnySubtreeMeetsPredicate(tree regex.RegExpTree, p func(regex.RegExpTree) bool) bool {
	if p(tree) {
		return true
	}
	for _, subTree := range tree.Children() {
		if ReportUntranspilableFeatures_AnySubtreeMeetsPredicate(subTree, p) {
			return true
		}
	}
	return false
}

// Java: ReportUntranspilableFeatures.transpileSubtree(Node)
func (r *ReportUntranspilableFeatures) TranspileSubtree(subtree *rhino.Node) *rhino.Node {
	r.CheckForUntranspilable(subtree)
	return subtree
}
