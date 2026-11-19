/*
 * Copyright 2015 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/RewritePolyfills.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_resources "github.com/paullepoulpe/jsgo/transpiled/gen/resources"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

var lambdaSingleton_rewrite_polyfills_1 j2gort.Ext_java_util_function_Function
var lambdaSingleton_rewrite_polyfills_2 j2gort.Ext_java_util_function_Predicate

func init() {
	lambdaSingleton_rewrite_polyfills_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_3 any) any {
		p := j2gort.CastPtr[PolyfillUsageFinder_Polyfill](lp_3)
		_ = p
		return j2gort.BoxString(p.library)
	})
	lambdaSingleton_rewrite_polyfills_2 = j2gort.NewLambda_Ext_java_util_function_Predicate(func(lp_4 any) bool {
		p := j2gort.CastPtr[PolyfillUsageFinder_Polyfill](lp_4)
		_ = p
		return !(j2gort.CheckNotNullString(p.library) == "")
	})
}

// Java class: com.google.javascript.jscomp.RewritePolyfills
type RewritePolyfills struct {
	compiler                 AbstractCompiler               `gson:"compiler"`
	runtimeJsLibManager      *RuntimeJsLibManager           `gson:"runtimeJsLibManager"`
	polyfills                *PolyfillUsageFinder_Polyfills `gson:"polyfills"`
	libraries                j2gort.Ext_java_util_Set       `gson:"libraries"`
	idHash_                  atomic.Int32
	injectPolyfills          bool                         `gson:"injectPolyfills"`
	isolatePolyfills         bool                         `gson:"isolatePolyfills"`
	injectPolyfillsNewerThan CompilerOptions_LanguageMode `gson:"injectPolyfillsNewerThan"`
	disambiguateProperties   bool                         `gson:"disambiguateProperties"`
}

func (this_ *RewritePolyfills) BindSelf_(s any) {
	_ = s
}
func (this_ *RewritePolyfills) IsJavaNil_() bool { return this_ == nil }
func (*RewritePolyfills) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.RewritePolyfills", "RewritePolyfills"
}
func (this_ *RewritePolyfills) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*RewritePolyfills) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.RewritePolyfills", "com.google.javascript.jscomp.RewritePolyfills", func(v_ any) bool {
		return j2gort.InstanceOf[*RewritePolyfills](v_)
	})
}
func (*RewritePolyfills) Implements_jscomp_CompilerPass() {}

var static_RewritePolyfills_INSUFFICIENT_OUTPUT_VERSION_ERROR *DiagnosticType

func StaticGet_RewritePolyfills_INSUFFICIENT_OUTPUT_VERSION_ERROR() *DiagnosticType {
	clinit_RewritePolyfills()
	return static_RewritePolyfills_INSUFFICIENT_OUTPUT_VERSION_ERROR
}
func StaticSet_RewritePolyfills_INSUFFICIENT_OUTPUT_VERSION_ERROR(v *DiagnosticType) {
	clinit_RewritePolyfills()
	static_RewritePolyfills_INSUFFICIENT_OUTPUT_VERSION_ERROR = v
}

var static_RewritePolyfills_JSCOMP_POLYFILL pkg_rhino.QualifiedName

func StaticGet_RewritePolyfills_JSCOMP_POLYFILL() pkg_rhino.QualifiedName {
	clinit_RewritePolyfills()
	return static_RewritePolyfills_JSCOMP_POLYFILL
}
func StaticSet_RewritePolyfills_JSCOMP_POLYFILL(v pkg_rhino.QualifiedName) {
	clinit_RewritePolyfills()
	static_RewritePolyfills_JSCOMP_POLYFILL = v
}

var clinitDone_RewritePolyfills uint32
var clinitState_RewritePolyfills j2gort.ClinitState

func clinit_RewritePolyfills() {
	if atomic.LoadUint32(&clinitDone_RewritePolyfills) != 0 {
		return
	}
	clinitSlow_RewritePolyfills()
}

func clinitSlow_RewritePolyfills() {
	if !j2gort.EnterClinit(&clinitState_RewritePolyfills, &clinitDone_RewritePolyfills) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_RewritePolyfills, &clinitDone_RewritePolyfills)
	static_RewritePolyfills_INSUFFICIENT_OUTPUT_VERSION_ERROR = DiagnosticType_Disabled("JSC_INSUFFICIENT_OUTPUT_VERSION", "Built-in ''{0}'' not supported in output version {1}")
	static_RewritePolyfills_JSCOMP_POLYFILL = pkg_rhino.QualifiedName_Of("$jscomp.polyfill")
}

func (this_ *RewritePolyfills) initFields_RewritePolyfills() {
	_ = this_
}

// Java constructor: RewritePolyfills(com.google.javascript.jscomp.AbstractCompiler compiler, boolean injectPolyfills, boolean isolatePolyfills, com.google.javascript.jscomp.CompilerOptions.LanguageMode injectPolyfillsNewerThan, boolean disambiguateProperties)
func NewRewritePolyfills_5(compiler AbstractCompiler, injectPolyfills bool, isolatePolyfills bool, injectPolyfillsNewerThan CompilerOptions_LanguageMode, disambiguateProperties bool) *RewritePolyfills {
	clinit_RewritePolyfills()
	this_ := &RewritePolyfills{}
	InitRewritePolyfills_5(this_, compiler, injectPolyfills, isolatePolyfills, injectPolyfillsNewerThan, disambiguateProperties)
	return this_
}

func InitRewritePolyfills_5(this_ *RewritePolyfills, compiler AbstractCompiler, injectPolyfills bool, isolatePolyfills bool, injectPolyfillsNewerThan CompilerOptions_LanguageMode, disambiguateProperties bool) {
	_ = compiler
	_ = injectPolyfills
	_ = isolatePolyfills
	_ = injectPolyfillsNewerThan
	_ = disambiguateProperties
	InitRewritePolyfills_7(this_, compiler, compiler.GetRuntimeJsLibManager(), PolyfillUsageFinder_Polyfills_FromTable(pkg_resources.ResourceLoader_LoadTextResource(j2gort.ClassLiteral("com.google.javascript.jscomp.RewritePolyfills", "RewritePolyfills"), "js/polyfills.txt")), injectPolyfills, isolatePolyfills, injectPolyfillsNewerThan, disambiguateProperties)
}

// Java constructor: RewritePolyfills(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.js.RuntimeJsLibManager runtimeJsLibManager, com.google.javascript.jscomp.PolyfillUsageFinder.Polyfills polyfills, boolean injectPolyfills, boolean isolatePolyfills, com.google.javascript.jscomp.CompilerOptions.LanguageMode injectPolyfillsNewerThan, boolean disambiguateProperties)
func NewRewritePolyfills_7(compiler AbstractCompiler, runtimeJsLibManager *RuntimeJsLibManager, polyfills *PolyfillUsageFinder_Polyfills, injectPolyfills bool, isolatePolyfills bool, injectPolyfillsNewerThan CompilerOptions_LanguageMode, disambiguateProperties bool) *RewritePolyfills {
	clinit_RewritePolyfills()
	this_ := &RewritePolyfills{}
	InitRewritePolyfills_7(this_, compiler, runtimeJsLibManager, polyfills, injectPolyfills, isolatePolyfills, injectPolyfillsNewerThan, disambiguateProperties)
	return this_
}

func InitRewritePolyfills_7(this_ *RewritePolyfills, compiler AbstractCompiler, runtimeJsLibManager *RuntimeJsLibManager, polyfills *PolyfillUsageFinder_Polyfills, injectPolyfills bool, isolatePolyfills bool, injectPolyfillsNewerThan CompilerOptions_LanguageMode, disambiguateProperties bool) {
	_ = compiler
	_ = runtimeJsLibManager
	_ = polyfills
	_ = injectPolyfills
	_ = isolatePolyfills
	_ = injectPolyfillsNewerThan
	_ = disambiguateProperties
	this_.initFields_RewritePolyfills()
	this_.compiler = compiler
	this_.runtimeJsLibManager = runtimeJsLibManager
	this_.polyfills = polyfills
	this_.injectPolyfills = injectPolyfills
	this_.isolatePolyfills = isolatePolyfills
	this_.injectPolyfillsNewerThan = injectPolyfillsNewerThan
	this_.disambiguateProperties = disambiguateProperties
}

// Java method: RewritePolyfills.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *RewritePolyfills) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if this_.isolatePolyfills {
		var jscompLookupMethodDecl pkg_rhino.Node = pkg_rhino.IR_Var_1(pkg_rhino.IR_Name("$jscomp$lookupPolyfilledValue"))
		_ = jscompLookupMethodDecl
		var synthesizedExternsAstRoot pkg_rhino.Node = this_.compiler.GetSynthesizedExternsInput().GetAstRoot(this_.compiler)
		_ = synthesizedExternsAstRoot
		jscompLookupMethodDecl.SrcrefTree(synthesizedExternsAstRoot)
		synthesizedExternsAstRoot.AddChildToBack(jscompLookupMethodDecl)
		this_.compiler.ReportChangeToEnclosingScope(jscompLookupMethodDecl)
	}
	if !(this_.injectPolyfills) && (this_.injectPolyfillsNewerThan == 0) {
		return
	}
	if this_.injectPolyfills {
		this_.libraries = j2gort.NewExt_java_util_LinkedHashSet()
		{
			clinit_PolyfillUsageFinder()
			var val_1 PolyfillUsageFinder
			InitPolyfillUsageFinder(&val_1, this_.compiler, this_.polyfills, this_.disambiguateProperties)
			(&val_1).TraverseExcludingGuarded(root, func() j2gort.Ext_java_util_function_Consumer {
				bound_2 := this_
				return j2gort.NewLambda_Ext_java_util_function_Consumer(func(mrp0_ any) {
					bound_2.priv_RewritePolyfills_inject_4991(j2gort.CastPtr[PolyfillUsageFinder_PolyfillUsage](mrp0_))
				})
			}())
		}
	}
	var librariesToInject j2gort.Ext_com_google_common_collect_ImmutableList = nil
	_ = librariesToInject
	if this_.injectPolyfillsNewerThan != 0 {
		var polyfillsToInject j2gort.Ext_com_google_common_collect_ImmutableList = this_.polyfills.GetPolyfillsNewerThan(this_.injectPolyfillsNewerThan)
		_ = polyfillsToInject
		librariesToInject = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableList](polyfillsToInject.Stream().Filter(lambdaSingleton_rewrite_polyfills_2).Map(lambdaSingleton_rewrite_polyfills_1).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableList_ToImmutableList()))
	} else {
		librariesToInject = j2gort.Ext_com_google_common_collect_ImmutableList_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.libraries))
	}
	this_.priv_RewritePolyfills_injectAll_4988(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](librariesToInject), (this_.injectPolyfillsNewerThan != 0))
}

// Java method: RewritePolyfills.injectAll(java.lang.Iterable<java.lang.String> librariesToInject, boolean forceInjection)
func (this_ *RewritePolyfills) priv_RewritePolyfills_injectAll_4988(librariesToInject j2gort.Ext_java_lang_Iterable, forceInjection bool) {
	_ = librariesToInject
	_ = forceInjection
	var lastNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = lastNode
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(librariesToInject))
		for cur_5.HasNext() {
			library := cur_5.NextStr()
			_ = library
			j2gort.CheckNotNullString(library)
			if !(!(j2gort.CheckNotNullString(library) == "")) {
				j2gort.CheckStateStr(false, "unexpected empty library")
			}
			lastNode = this_.runtimeJsLibManager.EnsureLibraryInjected(library, forceInjection)
		}
	}
	if lastNode != 0 {
		var parent pkg_rhino.Node = lastNode.GetParent()
		_ = parent
		this_.priv_RewritePolyfills_removeUnneededPolyfills_4989(parent, lastNode.GetNext())
		this_.compiler.ReportChangeToEnclosingScope(parent)
	}
}

// Java method: RewritePolyfills.removeUnneededPolyfills(com.google.javascript.rhino.Node parent, com.google.javascript.rhino.Node runtimeEnd)
func (this_ *RewritePolyfills) priv_RewritePolyfills_removeUnneededPolyfills_4989(parent pkg_rhino.Node, runtimeEnd pkg_rhino.Node) {
	_ = parent
	_ = runtimeEnd
	var node pkg_rhino.Node = parent.GetFirstChild()
	_ = node
	var outputFeatureSet *pkg_parser.FeatureSet = this_.compiler.GetOptions().GetOutputFeatureSet()
	_ = outputFeatureSet
	for (node != 0) && (node != runtimeEnd) {
		var next pkg_rhino.Node = node.GetNext()
		_ = next
		var polyfillSupportedFeatureSet *pkg_parser.FeatureSet = this_.priv_RewritePolyfills_getPolyfillSupportedFeatureSet_4990(node)
		_ = polyfillSupportedFeatureSet
		if (polyfillSupportedFeatureSet != nil) && outputFeatureSet.Contains_1_FeatureSet(polyfillSupportedFeatureSet) {
			NodeUtil_RemoveChild(parent, node)
			NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(node, this_.compiler)
		}
		node = next
	}
}

// Java method: RewritePolyfills.getPolyfillSupportedFeatureSet(com.google.javascript.rhino.Node maybePolyfill)
func (this_ *RewritePolyfills) priv_RewritePolyfills_getPolyfillSupportedFeatureSet_4990(maybePolyfill pkg_rhino.Node) *pkg_parser.FeatureSet {
	_ = maybePolyfill
	var polyfillSupportFeatureSet *pkg_parser.FeatureSet = nil
	_ = polyfillSupportFeatureSet
	if NodeUtil_IsExprCall(maybePolyfill) {
		var call pkg_rhino.Node = maybePolyfill.GetFirstChild()
		_ = call
		var name pkg_rhino.Node = call.GetFirstChild()
		_ = name
		if StaticGet_RewritePolyfills_JSCOMP_POLYFILL().Matches(name) {
			var nativeVersionStr string = name.GetNext().GetNext().GetNext().GetString()
			_ = nativeVersionStr
			polyfillSupportFeatureSet = pkg_parser.FeatureSet_ValueOf(nativeVersionStr)
			polyfillSupportFeatureSet = PolyfillUsageFinder_GetPolyfillSupportedFeatureSet(nativeVersionStr)
		}
	}
	return polyfillSupportFeatureSet
}

// Java method: RewritePolyfills.inject(com.google.javascript.jscomp.PolyfillUsageFinder.PolyfillUsage polyfillUsage)
func (this_ *RewritePolyfills) priv_RewritePolyfills_inject_4991(polyfillUsage *PolyfillUsageFinder_PolyfillUsage) {
	_ = polyfillUsage
	var polyfill *PolyfillUsageFinder_Polyfill = polyfillUsage.Polyfill()
	_ = polyfill
	var outputFeatureSet *pkg_parser.FeatureSet = this_.compiler.GetOptions().GetOutputFeatureSet()
	_ = outputFeatureSet
	var featuresRequiredByPolyfill *pkg_parser.FeatureSet = pkg_parser.FeatureSet_ValueOf(polyfill.polyfillVersion)
	_ = featuresRequiredByPolyfill
	{
		var iifeCond_1 bool
		{

			re_6 := polyfill.kind
			ae_7 := StaticGet_PolyfillUsageFinder_Polyfill_Kind_STATIC()
			if re_6 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_6 == ae_7
		}
		if iifeCond_1 && !(outputFeatureSet.Contains_1_FeatureSet(featuresRequiredByPolyfill)) {
			this_.compiler.Report_1(JSError_Make_3(polyfillUsage.Node(), StaticGet_RewritePolyfills_INSUFFICIENT_OUTPUT_VERSION_ERROR(), []string{polyfillUsage.Name(), outputFeatureSet.Version()}))
		}
	}
	if !(outputFeatureSet.Contains_1_FeatureSet(pkg_parser.FeatureSet_ValueOf(polyfill.nativeVersion))) && !(j2gort.CheckNotNullString(polyfill.library) == "") {
		j2gort.SetAddStr(this_.libraries, polyfill.library)
	}
}
