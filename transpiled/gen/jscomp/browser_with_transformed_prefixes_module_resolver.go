/*
 * Copyright 2017 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/deps/BrowserWithTransformedPrefixesModuleResolver.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_browser_with_transformed_prefixes_module_resolver_1 j2gort.Ext_java_util_function_Function
var lambdaSingleton_browser_with_transformed_prefixes_module_resolver_2 j2gort.Ext_java_util_function_ToIntFunction
var lambdaSingleton_browser_with_transformed_prefixes_module_resolver_3 j2gort.Ext_java_util_function_Function
var lambdaSingleton_browser_with_transformed_prefixes_module_resolver_4 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_browser_with_transformed_prefixes_module_resolver_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return j2gort.BoxString(j2gort.CastPtr[BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement](mrp0_).Prefix())
	})
	lambdaSingleton_browser_with_transformed_prefixes_module_resolver_2 = j2gort.NewLambda_Ext_java_util_function_ToIntFunction(func(lp_1 any) int32 {
		r := j2gort.CastPtr[BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement](lp_1)
		_ = r
		return j2gort.JStrLen(r.Prefix())
	})
	lambdaSingleton_browser_with_transformed_prefixes_module_resolver_3 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_2 any) any {
		entry := j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](lp_2)
		_ = entry
		return j2gort.PtrToAny(BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement_Of(j2gort.CastString(entry.GetKey()), j2gort.CastString(entry.GetValue())))
	})
	lambdaSingleton_browser_with_transformed_prefixes_module_resolver_4 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return j2gort.BoxString(j2gort.CastPtr[BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement](mrp0_).Prefix())
	})
}

// Java class: com.google.javascript.jscomp.deps.BrowserWithTransformedPrefixesModuleResolver
type BrowserWithTransformedPrefixesModuleResolver struct {
	ModuleResolverBase
	prefixReplacements j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"prefixReplacements"`
	expectedPrefixes   string                                            `gson:"expectedPrefixes"`
}

func (this_ *BrowserWithTransformedPrefixesModuleResolver) BindSelf_(s any) {
	_ = s
}
func (this_ *BrowserWithTransformedPrefixesModuleResolver) IsJavaNil_() bool { return this_ == nil }
func (*BrowserWithTransformedPrefixesModuleResolver) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.BrowserWithTransformedPrefixesModuleResolver", "BrowserWithTransformedPrefixesModuleResolver"
}
func (this_ *BrowserWithTransformedPrefixesModuleResolver) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ModuleResolverBase.IdentityHash_()
}
func (*BrowserWithTransformedPrefixesModuleResolver) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.BrowserWithTransformedPrefixesModuleResolver", "com.google.javascript.jscomp.deps.BrowserWithTransformedPrefixesModuleResolver", func(v_ any) bool {
		return j2gort.InstanceOf[*BrowserWithTransformedPrefixesModuleResolver](v_)
	})
}
func (*BrowserWithTransformedPrefixesModuleResolver) Implements_jscomp_ModuleResolver() {}

var static_BrowserWithTransformedPrefixesModuleResolver_INVALID_AMBIGUOUS_PATH *DiagnosticType

func StaticGet_BrowserWithTransformedPrefixesModuleResolver_INVALID_AMBIGUOUS_PATH() *DiagnosticType {
	clinit_BrowserWithTransformedPrefixesModuleResolver()
	return static_BrowserWithTransformedPrefixesModuleResolver_INVALID_AMBIGUOUS_PATH
}
func StaticSet_BrowserWithTransformedPrefixesModuleResolver_INVALID_AMBIGUOUS_PATH(v *DiagnosticType) {
	clinit_BrowserWithTransformedPrefixesModuleResolver()
	static_BrowserWithTransformedPrefixesModuleResolver_INVALID_AMBIGUOUS_PATH = v
}

var clinitDone_BrowserWithTransformedPrefixesModuleResolver uint32
var clinitState_BrowserWithTransformedPrefixesModuleResolver j2gort.ClinitState

func clinit_BrowserWithTransformedPrefixesModuleResolver() {
	if atomic.LoadUint32(&clinitDone_BrowserWithTransformedPrefixesModuleResolver) != 0 {
		return
	}
	clinitSlow_BrowserWithTransformedPrefixesModuleResolver()
}

func clinitSlow_BrowserWithTransformedPrefixesModuleResolver() {
	if !j2gort.EnterClinit(&clinitState_BrowserWithTransformedPrefixesModuleResolver, &clinitDone_BrowserWithTransformedPrefixesModuleResolver) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_BrowserWithTransformedPrefixesModuleResolver, &clinitDone_BrowserWithTransformedPrefixesModuleResolver)
	static_BrowserWithTransformedPrefixesModuleResolver_INVALID_AMBIGUOUS_PATH = DiagnosticType_Error("JSC_INVALID_AMBIGUOUS_PATH", "The path \"{0}\" is invalid. Expected any of the following prefixes for non-relative paths: {1}.")
}

func (this_ *BrowserWithTransformedPrefixesModuleResolver) initFields_BrowserWithTransformedPrefixesModuleResolver() {
	_ = this_
	this_.expectedPrefixes = j2gort.NullString
}

// Java constructor: BrowserWithTransformedPrefixesModuleResolver(com.google.common.collect.ImmutableSet<java.lang.String> modulePaths, com.google.common.collect.ImmutableList<java.lang.String> moduleRootPaths, com.google.javascript.jscomp.ErrorHandler errorHandler, com.google.javascript.jscomp.deps.ModuleLoader.PathEscaper pathEscaper, com.google.common.collect.ImmutableMap<java.lang.String,java.lang.String> prefixReplacements)
func NewBrowserWithTransformedPrefixesModuleResolver(modulePaths j2gort.Ext_com_google_common_collect_ImmutableSet, moduleRootPaths j2gort.Ext_com_google_common_collect_ImmutableList, errorHandler ErrorHandler, pathEscaper *ModuleLoader_PathEscaper, prefixReplacements j2gort.Ext_com_google_common_collect_ImmutableMap) *BrowserWithTransformedPrefixesModuleResolver {
	clinit_BrowserWithTransformedPrefixesModuleResolver()
	this_ := &BrowserWithTransformedPrefixesModuleResolver{}
	InitBrowserWithTransformedPrefixesModuleResolver(this_, modulePaths, moduleRootPaths, errorHandler, pathEscaper, prefixReplacements)
	return this_
}

func InitBrowserWithTransformedPrefixesModuleResolver(this_ *BrowserWithTransformedPrefixesModuleResolver, modulePaths j2gort.Ext_com_google_common_collect_ImmutableSet, moduleRootPaths j2gort.Ext_com_google_common_collect_ImmutableList, errorHandler ErrorHandler, pathEscaper *ModuleLoader_PathEscaper, prefixReplacements j2gort.Ext_com_google_common_collect_ImmutableMap) {
	_ = modulePaths
	_ = moduleRootPaths
	_ = errorHandler
	_ = pathEscaper
	_ = prefixReplacements
	InitModuleResolver(&this_.ModuleResolverBase, modulePaths, moduleRootPaths, errorHandler, pathEscaper)
	this_.initFields_BrowserWithTransformedPrefixesModuleResolver()
	this_.prefixReplacements = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](prefixReplacements.EntrySet()).Stream().Map(lambdaSingleton_browser_with_transformed_prefixes_module_resolver_3).Sorted_1(j2gort.Ext_java_util_Comparator_ComparingInt(lambdaSingleton_browser_with_transformed_prefixes_module_resolver_2).Reversed().ThenComparing_1_Ext_java_util_function_Function(lambdaSingleton_browser_with_transformed_prefixes_module_resolver_1)).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableSet_ToImmutableSet()))
	this_.expectedPrefixes = j2gort.CastString(this_.prefixReplacements.Stream().Map(lambdaSingleton_browser_with_transformed_prefixes_module_resolver_4).Sorted_0().Collect_1(j2gort.Ext_java_util_stream_Collectors_Joining_1(", ")))
}

// Java method: BrowserWithTransformedPrefixesModuleResolver.replacePathPrefixes(java.lang.String moduleAddress)
func (this_ *BrowserWithTransformedPrefixesModuleResolver) priv_BrowserWithTransformedPrefixesModuleResolver_replacePathPrefixes_6012(moduleAddress string) string {
	_ = moduleAddress
	if moduleAddress == j2gort.NullString {
		return j2gort.NullString
	}
	var transformedAddress string = moduleAddress
	_ = transformedAddress
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.prefixReplacements))
		for cur_3.HasNext() {
			rawElem_4 := cur_3.Next()
			prefixReplacement := j2gort.CastPtr[BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement](rawElem_4)
			_ = prefixReplacement
			if javalang.StringStartsWith(j2gort.CheckNotNullString(moduleAddress), j2gort.CheckNotNullString(prefixReplacement.Prefix())) {
				transformedAddress = (j2gort.StrValString(prefixReplacement.Replacement()) + j2gort.StrValString(j2gort.JStrSubFrom(moduleAddress, j2gort.JStrLen(prefixReplacement.Prefix()))))
				break
			}
		}
	}
	return transformedAddress
}

// Java method: BrowserWithTransformedPrefixesModuleResolver.resolveJsModule(java.lang.String scriptAddress, java.lang.String moduleAddress, java.lang.String sourcename, int lineno, int colno)
func (this_ *BrowserWithTransformedPrefixesModuleResolver) ResolveJsModule(scriptAddress string, moduleAddress string, sourcename string, lineno int32, colno int32) string {
	_ = scriptAddress
	_ = moduleAddress
	_ = sourcename
	_ = lineno
	_ = colno
	var transformedAddress string = this_.priv_BrowserWithTransformedPrefixesModuleResolver_replacePathPrefixes_6012(moduleAddress)
	_ = transformedAddress
	var replaced bool = ((transformedAddress != j2gort.NullString) && !(j2gort.CheckNotNullString(transformedAddress) == moduleAddress))
	_ = replaced
	if !(replaced) && ModuleLoader_IsAmbiguousIdentifier(transformedAddress) {
		this_.GetModuleResolverBase().errorHandler.Report_2(StaticGet_CheckLevel_WARNING(), JSError_Make_5(sourcename, lineno, colno, StaticGet_BrowserWithTransformedPrefixesModuleResolver_INVALID_AMBIGUOUS_PATH(), []string{transformedAddress, this_.expectedPrefixes}))
		return j2gort.NullString
	}
	return this_.Locate(scriptAddress, transformedAddress)
}

// Java method: BrowserWithTransformedPrefixesModuleResolver.resolveJsModuleSilently(java.lang.String scriptAddress, java.lang.String moduleAddress)
func (this_ *BrowserWithTransformedPrefixesModuleResolver) ResolveJsModuleSilently(scriptAddress string, moduleAddress string) string {
	_ = scriptAddress
	_ = moduleAddress
	var transformedAddress string = this_.priv_BrowserWithTransformedPrefixesModuleResolver_replacePathPrefixes_6012(moduleAddress)
	_ = transformedAddress
	var replaced bool = ((transformedAddress != j2gort.NullString) && !(j2gort.CheckNotNullString(transformedAddress) == moduleAddress))
	_ = replaced
	if !(replaced) && ModuleLoader_IsAmbiguousIdentifier(transformedAddress) {
		return j2gort.NullString
	}
	return this_.Locate(scriptAddress, transformedAddress)
}

// Java method: BrowserWithTransformedPrefixesModuleResolver.resolveModuleAsPath(java.lang.String scriptAddress, java.lang.String moduleAddress)
func (this_ *BrowserWithTransformedPrefixesModuleResolver) ResolveModuleAsPath(scriptAddress string, moduleAddress string) string {
	_ = scriptAddress
	_ = moduleAddress
	if ModuleLoader_IsRelativeIdentifier(moduleAddress) {
		return this_.ModuleResolverBase.ResolveModuleAsPath(scriptAddress, moduleAddress)
	}
	var transformedAddress string = moduleAddress
	_ = transformedAddress
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.prefixReplacements))
		for cur_5.HasNext() {
			rawElem_6 := cur_5.Next()
			prefixReplacement := j2gort.CastPtr[BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement](rawElem_6)
			_ = prefixReplacement
			if javalang.StringStartsWith(j2gort.CheckNotNullString(moduleAddress), j2gort.CheckNotNullString(prefixReplacement.Prefix())) {
				transformedAddress = (j2gort.StrValString(prefixReplacement.Replacement()) + j2gort.StrValString(j2gort.JStrSubFrom(moduleAddress, j2gort.JStrLen(prefixReplacement.Prefix()))))
				break
			}
		}
	}
	return ModuleLoader_Normalize(transformedAddress, j2gort.ImmListFromTypedSlice(this_.GetModuleResolverBase().moduleRootPaths))
}

// Java class: com.google.javascript.jscomp.deps.BrowserWithTransformedPrefixesModuleResolver.Factory
type BrowserWithTransformedPrefixesModuleResolver_Factory struct {
	prefixReplacements j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"prefixReplacements"`
	idHash_            atomic.Int32
}

func (this_ *BrowserWithTransformedPrefixesModuleResolver_Factory) BindSelf_(s any) {
	_ = s
}
func (this_ *BrowserWithTransformedPrefixesModuleResolver_Factory) IsJavaNil_() bool {
	return this_ == nil
}
func (*BrowserWithTransformedPrefixesModuleResolver_Factory) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.BrowserWithTransformedPrefixesModuleResolver$Factory", "Factory"
}
func (this_ *BrowserWithTransformedPrefixesModuleResolver_Factory) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*BrowserWithTransformedPrefixesModuleResolver_Factory) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.BrowserWithTransformedPrefixesModuleResolver.Factory", "com.google.javascript.jscomp.deps.BrowserWithTransformedPrefixesModuleResolver$Factory", func(v_ any) bool {
		return j2gort.InstanceOf[*BrowserWithTransformedPrefixesModuleResolver_Factory](v_)
	})
}
func (*BrowserWithTransformedPrefixesModuleResolver_Factory) Implements_jscomp_ModuleLoader_ModuleResolverFactory() {
}

func (this_ *BrowserWithTransformedPrefixesModuleResolver_Factory) initFields_BrowserWithTransformedPrefixesModuleResolver_Factory() {
	_ = this_
}

// Java constructor: Factory(com.google.common.collect.ImmutableMap<java.lang.String,java.lang.String> prefixReplacements)
func NewBrowserWithTransformedPrefixesModuleResolver_Factory(prefixReplacements j2gort.Ext_com_google_common_collect_ImmutableMap) *BrowserWithTransformedPrefixesModuleResolver_Factory {
	this_ := &BrowserWithTransformedPrefixesModuleResolver_Factory{}
	InitBrowserWithTransformedPrefixesModuleResolver_Factory(this_, prefixReplacements)
	return this_
}

func InitBrowserWithTransformedPrefixesModuleResolver_Factory(this_ *BrowserWithTransformedPrefixesModuleResolver_Factory, prefixReplacements j2gort.Ext_com_google_common_collect_ImmutableMap) {
	_ = prefixReplacements
	this_.initFields_BrowserWithTransformedPrefixesModuleResolver_Factory()
	this_.prefixReplacements = prefixReplacements
}

// Java method: Factory.create(com.google.common.collect.ImmutableSet<java.lang.String> modulePaths, com.google.common.collect.ImmutableList<java.lang.String> moduleRootPaths, com.google.javascript.jscomp.ErrorHandler errorHandler, com.google.javascript.jscomp.deps.ModuleLoader.PathEscaper pathEscaper)
func (this_ *BrowserWithTransformedPrefixesModuleResolver_Factory) Create(modulePaths j2gort.Ext_com_google_common_collect_ImmutableSet, moduleRootPaths j2gort.Ext_com_google_common_collect_ImmutableList, errorHandler ErrorHandler, pathEscaper *ModuleLoader_PathEscaper) ModuleResolver {
	_ = modulePaths
	_ = moduleRootPaths
	_ = errorHandler
	_ = pathEscaper
	return NewBrowserWithTransformedPrefixesModuleResolver(modulePaths, moduleRootPaths, errorHandler, pathEscaper, this_.prefixReplacements)
}

// Java record: com.google.javascript.jscomp.deps.BrowserWithTransformedPrefixesModuleResolver.PrefixReplacement
type BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement struct {
	prefix      string `gson:"prefix"`
	replacement string `gson:"replacement"`
}

func (this_ *BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement) BindSelf_(s any) {
	_ = s
}
func (this_ *BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement) IsJavaNil_() bool {
	return this_ == nil
}
func (*BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.BrowserWithTransformedPrefixesModuleResolver$PrefixReplacement", "PrefixReplacement"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.BrowserWithTransformedPrefixesModuleResolver.PrefixReplacement", "com.google.javascript.jscomp.deps.BrowserWithTransformedPrefixesModuleResolver$PrefixReplacement", func(v_ any) bool {
		return j2gort.InstanceOf[*BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement](v_)
	})
}
func (*BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement) Implements_j2gort_Ext_java_lang_Record() {
}

func (this_ *BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement) initFields_BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement() {
	_ = this_
}

// Java constructor: PrefixReplacement(java.lang.String prefix, java.lang.String replacement)
func NewBrowserWithTransformedPrefixesModuleResolver_PrefixReplacement(prefix string, replacement string) *BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement {
	this_ := &BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement{}
	InitBrowserWithTransformedPrefixesModuleResolver_PrefixReplacement(this_, prefix, replacement)
	return this_
}

func InitBrowserWithTransformedPrefixesModuleResolver_PrefixReplacement(this_ *BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement, prefix string, replacement string) {
	_ = prefix
	_ = replacement
	this_.initFields_BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement()
	j2gort.CheckNotNullString(prefix)
	j2gort.CheckNotNullString(replacement)
	this_.prefix = prefix
	this_.replacement = replacement
}

// Java static method: PrefixReplacement.of(java.lang.String prefix, java.lang.String replacement)
func BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement_Of(prefix string, replacement string) *BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement {
	_ = prefix
	_ = replacement
	return NewBrowserWithTransformedPrefixesModuleResolver_PrefixReplacement(prefix, replacement)
}

func (this_ *BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement) Prefix() string {
	return this_.prefix
}
func (this_ *BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement) Replacement() string {
	return this_.replacement
}
func (this_ *BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement) Equals(other_ any) bool {
	o_, ok_ := other_.(*BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.prefix != o_.prefix {
		return false
	}
	if this_.replacement != o_.replacement {
		return false
	}
	return true
}
func (this_ *BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement) HashCode() int32 {
	var h_ int32 = 0
	{
		hs_7 := this_.prefix
		h_ = 31 * h_
		if hs_7 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_7)
		}
	}
	{
		hs_8 := this_.replacement
		h_ = 31 * h_
		if hs_8 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_8)
		}
	}
	return h_
}
func (this_ *BrowserWithTransformedPrefixesModuleResolver_PrefixReplacement) String() string {
	_ = fmt.Sprint
	return "PrefixReplacement[" + "prefix=" + j2gort.StrVal(j2gort.BoxString(this_.prefix)) + ", " + "replacement=" + j2gort.StrVal(j2gort.BoxString(this_.replacement)) + "]"
}
