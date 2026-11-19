/*
 * Copyright 2016 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/deps/ModuleLoader.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_module_loader_1 j2gort.Ext_java_util_function_ToIntFunction
var lambdaSingleton_module_loader_2 ErrorHandler
var lambdaSingleton_module_loader_3 j2gort.Ext_java_util_function_Function
var lambdaSingleton_module_loader_4 j2gort.Ext_java_util_function_Function
var lambdaSingleton_module_loader_5 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_module_loader_1 = j2gort.NewLambda_Ext_java_util_function_ToIntFunction(func(mrp0_ any) int32 { return j2gort.JStrLen(j2gort.CastString(mrp0_)) })
	lambdaSingleton_module_loader_2 = NewLambda_ErrorHandler(func(level *CheckLevel, error_ *JSError) {
		_ = level
		_ = error_
	})
	lambdaSingleton_module_loader_3 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_1 any) any {
		n := j2gort.CastString(lp_1)
		_ = n
		return j2gort.BoxString(func() string {
			if ModuleLoader_IsAmbiguousIdentifier(n) {
				return ("/" + j2gort.StrValString(n))
			}
			return n
		}())
	})
	lambdaSingleton_module_loader_4 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_3 any) any {
		n := j2gort.CastString(lp_3)
		_ = n
		return j2gort.BoxString(func() string {
			if ModuleLoader_IsAmbiguousIdentifier(n) {
				return ("/" + j2gort.StrValString(n))
			}
			return n
		}())
	})
	lambdaSingleton_module_loader_5 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return j2gort.BoxString(j2gort.AsIface[DependencyInfo](mrp0_).GetName()) })
}

// Java class: com.google.javascript.jscomp.deps.ModuleLoader
type ModuleLoader struct {
	errorHandler    ErrorHandler                                       `gson:"errorHandler"`
	moduleRootPaths j2gort.Ext_com_google_common_collect_ImmutableList `gson:"moduleRootPaths"`
	pathResolver    *ModuleLoader_PathResolver                         `gson:"pathResolver"`
	pathEscaper     *ModuleLoader_PathEscaper                          `gson:"pathEscaper"`
	moduleResolver  ModuleResolver                                     `gson:"moduleResolver"`
	idHash_         atomic.Int32
}

func (this_ *ModuleLoader) BindSelf_(s any) {
	_ = s
}
func (this_ *ModuleLoader) IsJavaNil_() bool { return this_ == nil }
func (*ModuleLoader) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.ModuleLoader", "ModuleLoader"
}
func (this_ *ModuleLoader) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ModuleLoader) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.ModuleLoader", "com.google.javascript.jscomp.deps.ModuleLoader", func(v_ any) bool {
		return j2gort.InstanceOf[*ModuleLoader](v_)
	})
}

var static_ModuleLoader_MODULE_CONFLICT *DiagnosticType

func StaticGet_ModuleLoader_MODULE_CONFLICT() *DiagnosticType {
	clinit_ModuleLoader()
	return static_ModuleLoader_MODULE_CONFLICT
}
func StaticSet_ModuleLoader_MODULE_CONFLICT(v *DiagnosticType) {
	clinit_ModuleLoader()
	static_ModuleLoader_MODULE_CONFLICT = v
}
func StaticGet_ModuleLoader_MODULE_SLASH() string            { return "/" }
func StaticGet_ModuleLoader_DEFAULT_FILENAME_PREFIX() string { return "./" }
func StaticGet_ModuleLoader_JSC_BROWSER_SKIPLISTED_MARKER() string {
	return "$jscomp$browser$skiplisted"
}

var static_ModuleLoader_LOAD_WARNING *DiagnosticType

func StaticGet_ModuleLoader_LOAD_WARNING() *DiagnosticType {
	clinit_ModuleLoader()
	return static_ModuleLoader_LOAD_WARNING
}
func StaticSet_ModuleLoader_LOAD_WARNING(v *DiagnosticType) {
	clinit_ModuleLoader()
	static_ModuleLoader_LOAD_WARNING = v
}

var static_ModuleLoader_INVALID_MODULE_PATH *DiagnosticType

func StaticGet_ModuleLoader_INVALID_MODULE_PATH() *DiagnosticType {
	clinit_ModuleLoader()
	return static_ModuleLoader_INVALID_MODULE_PATH
}
func StaticSet_ModuleLoader_INVALID_MODULE_PATH(v *DiagnosticType) {
	clinit_ModuleLoader()
	static_ModuleLoader_INVALID_MODULE_PATH = v
}

var static_ModuleLoader_BEST_MATCH_PATH_ORDERING j2gort.Ext_java_util_Comparator

func StaticGet_ModuleLoader_BEST_MATCH_PATH_ORDERING() j2gort.Ext_java_util_Comparator {
	clinit_ModuleLoader()
	return static_ModuleLoader_BEST_MATCH_PATH_ORDERING
}
func StaticSet_ModuleLoader_BEST_MATCH_PATH_ORDERING(v j2gort.Ext_java_util_Comparator) {
	clinit_ModuleLoader()
	static_ModuleLoader_BEST_MATCH_PATH_ORDERING = v
}

var static_ModuleLoader_NOOP_ERROR_HANDER ErrorHandler

func StaticGet_ModuleLoader_NOOP_ERROR_HANDER() ErrorHandler {
	clinit_ModuleLoader()
	return static_ModuleLoader_NOOP_ERROR_HANDER
}
func StaticSet_ModuleLoader_NOOP_ERROR_HANDER(v ErrorHandler) {
	clinit_ModuleLoader()
	static_ModuleLoader_NOOP_ERROR_HANDER = v
}

var static_ModuleLoader_EMPTY *ModuleLoader

func StaticGet_ModuleLoader_EMPTY() *ModuleLoader {
	clinit_ModuleLoader()
	return static_ModuleLoader_EMPTY
}
func StaticSet_ModuleLoader_EMPTY(v *ModuleLoader) {
	clinit_ModuleLoader()
	static_ModuleLoader_EMPTY = v
}

var clinitDone_ModuleLoader uint32
var clinitState_ModuleLoader j2gort.ClinitState

func clinit_ModuleLoader() {
	if atomic.LoadUint32(&clinitDone_ModuleLoader) != 0 {
		return
	}
	clinitSlow_ModuleLoader()
}

func clinitSlow_ModuleLoader() {
	if !j2gort.EnterClinit(&clinitState_ModuleLoader, &clinitDone_ModuleLoader) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ModuleLoader, &clinitDone_ModuleLoader)
	static_ModuleLoader_MODULE_CONFLICT = DiagnosticType_Warning("JSC_MODULE_CONFLICT", "File cannot be a combination of goog.provide, goog.module, and/or ES6 module: {0}")
	static_ModuleLoader_LOAD_WARNING = DiagnosticType_Error("JSC_JS_MODULE_LOAD_WARNING", "Failed to load module \"{0}\"")
	static_ModuleLoader_INVALID_MODULE_PATH = DiagnosticType_Error("JSC_INVALID_MODULE_PATH", "Invalid module path \"{0}\" for resolution mode \"{1}\"")
	static_ModuleLoader_BEST_MATCH_PATH_ORDERING = j2gort.Ext_java_util_Comparator_ComparingInt(lambdaSingleton_module_loader_1).Reversed().ThenComparing_1_Ext_java_util_Comparator(j2gort.Ext_java_util_Comparator_NaturalOrder())
	static_ModuleLoader_NOOP_ERROR_HANDER = lambdaSingleton_module_loader_2
	static_ModuleLoader_EMPTY = ModuleLoader_S_Builder().SetModuleRoots(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())).SetInputs(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())).SetFactory(StaticGet_BrowserModuleResolver_FACTORY()).Build()
}

func (this_ *ModuleLoader) initFields_ModuleLoader() {
	_ = this_
}

// Java constructor: ModuleLoader(com.google.javascript.jscomp.deps.ModuleLoader.Builder builder)
func NewModuleLoader(builder *ModuleLoader_Builder) *ModuleLoader {
	clinit_ModuleLoader()
	this_ := &ModuleLoader{}
	InitModuleLoader(this_, builder)
	return this_
}

func InitModuleLoader(this_ *ModuleLoader, builder *ModuleLoader_Builder) {
	_ = builder
	this_.initFields_ModuleLoader()
	this_.errorHandler = j2gort.CheckNotNullIface(builder.errorHandler)
	this_.pathResolver = j2gort.CheckNotNullPtr(builder.pathResolver)
	this_.pathEscaper = j2gort.CheckNotNullPtr(builder.pathEscaper)
	this_.moduleRootPaths = this_.priv_ModuleLoader_createRootPaths_6105(j2gort.CheckNotNull(builder.moduleRoots))
	var modulePaths j2gort.Ext_com_google_common_collect_ImmutableSet = this_.priv_ModuleLoader_resolvePaths_6106(j2gort.CheckNotNull(builder.inputs))
	_ = modulePaths
	this_.moduleResolver = builder.factory.Create(modulePaths, this_.moduleRootPaths, this_.errorHandler, this_.pathEscaper)
}

// Java static method: ModuleLoader.builder()
func ModuleLoader_S_Builder() *ModuleLoader_Builder {
	clinit_ModuleLoader()
	return NewModuleLoader_Builder()
}

// Java method: ModuleLoader.getPackageJsonMainEntries()
func (this_ *ModuleLoader) GetPackageJsonMainEntries() j2gort.Ext_java_util_Map {
	return this_.moduleResolver.GetPackageJsonMainEntries()
}

// Java method: ModuleLoader.resolve(java.lang.String path)
func (this_ *ModuleLoader) Resolve(path string) *ModuleLoader_ModulePath {
	_ = path
	return NewModuleLoader_ModulePath(this_, ModuleLoader_Normalize(this_.pathEscaper.Escape(j2gort.CastString(this_.pathResolver.Apply(j2gort.BoxString(path)))), j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.moduleRootPaths)))
}

// Java static method: ModuleLoader.isRelativeIdentifier(java.lang.String name)
func ModuleLoader_IsRelativeIdentifier(name string) bool {
	clinit_ModuleLoader()
	_ = name
	return (((javalang.StringStartsWith(j2gort.CheckNotNullString(name), "."+"/") || javalang.StringStartsWith(j2gort.CheckNotNullString(name), ".."+"/")) || (j2gort.CheckNotNullString(name) == ".")) || (j2gort.CheckNotNullString(name) == ".."))
}

// Java static method: ModuleLoader.isAbsoluteIdentifier(java.lang.String name)
func ModuleLoader_IsAbsoluteIdentifier(name string) bool {
	clinit_ModuleLoader()
	_ = name
	return javalang.StringStartsWith(j2gort.CheckNotNullString(name), "/")
}

// Java static method: ModuleLoader.isAmbiguousIdentifier(java.lang.String name)
func ModuleLoader_IsAmbiguousIdentifier(name string) bool {
	clinit_ModuleLoader()
	_ = name
	return (!(ModuleLoader_IsAbsoluteIdentifier(name)) && !(ModuleLoader_IsRelativeIdentifier(name)))
}

// Java static method: ModuleLoader.isPathIdentifier(java.lang.String name)
func ModuleLoader_IsPathIdentifier(name string) bool {
	clinit_ModuleLoader()
	_ = name
	return (j2gort.JStrIndexOf(name, "/") >= 0)
}

// Java method: ModuleLoader.createRootPaths(java.lang.Iterable<java.lang.String> roots)
func (this_ *ModuleLoader) priv_ModuleLoader_createRootPaths_6105(roots j2gort.Ext_java_lang_Iterable) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = roots
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSortedSet](j2gort.Ext_com_google_common_collect_Streams_Stream_1_Ext_java_lang_Iterable(roots).Map(j2gort.ToIface_Ext_java_util_function_Function(this_.pathResolver)).Map(func() j2gort.Ext_java_util_function_Function {
		bound_2 := this_.pathEscaper
		return j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return j2gort.BoxString(bound_2.Escape(j2gort.CastString(mrp0_))) })
	}()).Map(lambdaSingleton_module_loader_3).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableSortedSet_ToImmutableSortedSet(StaticGet_ModuleLoader_BEST_MATCH_PATH_ORDERING()))).AsList()
}

// Java method: ModuleLoader.resolvePaths(java.lang.Iterable<? extends com.google.javascript.jscomp.deps.DependencyInfo> inputs)
func (this_ *ModuleLoader) priv_ModuleLoader_resolvePaths_6106(inputs j2gort.Ext_java_lang_Iterable) j2gort.Ext_com_google_common_collect_ImmutableSet {
	_ = inputs
	var dupeModulePaths j2gort.Ext_com_google_common_collect_ImmutableMultiset = j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMultiset](j2gort.Ext_com_google_common_collect_Streams_Stream_1_Ext_java_lang_Iterable(inputs).Map(lambdaSingleton_module_loader_5).Map(j2gort.ToIface_Ext_java_util_function_Function(this_.pathResolver)).Map(func() j2gort.Ext_java_util_function_Function {
		bound_5 := this_.pathEscaper
		return j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return j2gort.BoxString(bound_5.Escape(j2gort.CastString(mrp0_))) })
	}()).Map(j2gort.NewLambda_Ext_java_util_function_Function(func(lp_4 any) any {
		p := j2gort.CastString(lp_4)
		_ = p
		return j2gort.BoxString(ModuleLoader_Normalize(p, j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.moduleRootPaths)))
	})).Map(lambdaSingleton_module_loader_4).Sorted_1(StaticGet_ModuleLoader_BEST_MATCH_PATH_ORDERING()).Collect_1(j2gort.Ext_com_google_common_collect_ImmutableMultiset_ToImmutableMultiset_0()))
	_ = dupeModulePaths
	if !(dupeModulePaths.Size() == j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](dupeModulePaths.ElementSet()).Size()) {
		j2gort.CheckStateFmt(false, "Duplicate module paths after resolving: %s", dupeModulePaths)
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](dupeModulePaths.ElementSet())
}

// Java static method: ModuleLoader.normalize(java.lang.String path, java.lang.Iterable<java.lang.String> moduleRootPaths)
func ModuleLoader_Normalize(path string, moduleRootPaths j2gort.Ext_java_lang_Iterable) string {
	clinit_ModuleLoader()
	_ = path
	_ = moduleRootPaths
	var normalizedPath string = path
	_ = normalizedPath
	if ModuleLoader_IsAmbiguousIdentifier(normalizedPath) {
		normalizedPath = ("/" + j2gort.StrValString(normalizedPath))
	}
	{
		cur_6 := j2gort.NewIterCursor(j2gort.CheckNotNull(moduleRootPaths))
		for cur_6.HasNext() {
			moduleRoot := cur_6.NextStr()
			_ = moduleRoot
			if javalang.StringStartsWith(j2gort.CheckNotNullString(normalizedPath), j2gort.CheckNotNullString(moduleRoot)) {
				var trailing string = j2gort.JStrSubFrom(normalizedPath, j2gort.JStrLen(moduleRoot))
				_ = trailing
				if javalang.StringStartsWith(j2gort.CheckNotNullString(trailing), "/") {
					return j2gort.JStrSubFrom(trailing, j2gort.JStrLen("/"))
				}
			}
		}
	}
	return path
}

// Java static method: ModuleLoader.relativePathFrom(java.lang.String fromUriPath, java.lang.String toUriPath)
func ModuleLoader_RelativePathFrom(fromUriPath string, toUriPath string) string {
	clinit_ModuleLoader()
	_ = fromUriPath
	_ = toUriPath
	var fromPath j2gort.Ext_java_nio_file_Path = j2gort.Ext_java_nio_file_Path_Of_2(fromUriPath, []string{})
	_ = fromPath
	var toPath j2gort.Ext_java_nio_file_Path = j2gort.Ext_java_nio_file_Path_Of_2(toUriPath, []string{})
	_ = toPath
	var fromFolder j2gort.Ext_java_nio_file_Path = fromPath.GetParent()
	_ = fromFolder
	if (fromFolder == nil) && (toPath.GetParent() == nil) {
		return ("./" + j2gort.StrValString(toUriPath))
	} else if (fromFolder == nil) && ((javalang.StringStartsWith(j2gort.CheckNotNullString(toUriPath), ".") || javalang.StringStartsWith(j2gort.CheckNotNullString(toPath.String()), "/")) || javalang.StringStartsWith(j2gort.CheckNotNullString(toPath.String()), "\\")) {
		return toUriPath
	} else if fromFolder == nil {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Relative path between URIs cannot be calculated")))
	}
	var calculatedPath string = fromFolder.Relativize(toPath).String()
	_ = calculatedPath
	if javalang.StringStartsWith(j2gort.CheckNotNullString(calculatedPath), ".") || javalang.StringStartsWith(j2gort.CheckNotNullString(calculatedPath), "/") {
		return calculatedPath
	}
	return ("./" + j2gort.StrValString(calculatedPath))
}

// Java method: ModuleLoader.setErrorHandler(com.google.javascript.jscomp.ErrorHandler errorHandler)
func (this_ *ModuleLoader) SetErrorHandler(errorHandler ErrorHandler) {
	_ = errorHandler
	if errorHandler == nil || errorHandler.IsJavaNil_() {
		this_.errorHandler = StaticGet_ModuleLoader_NOOP_ERROR_HANDER()
	} else {
		this_.errorHandler = errorHandler
	}
	this_.moduleResolver.SetErrorHandler(this_.errorHandler)
}

// Java method: ModuleLoader.getErrorHandler()
func (this_ *ModuleLoader) GetErrorHandler() ErrorHandler {
	return this_.errorHandler
}

// Java class: com.google.javascript.jscomp.deps.ModuleLoader.Builder
type ModuleLoader_Builder struct {
	errorHandler ErrorHandler                       `gson:"errorHandler"`
	moduleRoots  j2gort.Ext_java_lang_Iterable      `gson:"moduleRoots"`
	inputs       j2gort.Ext_java_lang_Iterable      `gson:"inputs"`
	factory      ModuleLoader_ModuleResolverFactory `gson:"factory"`
	pathResolver *ModuleLoader_PathResolver         `gson:"pathResolver"`
	pathEscaper  *ModuleLoader_PathEscaper          `gson:"pathEscaper"`
	idHash_      atomic.Int32
}

func (this_ *ModuleLoader_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *ModuleLoader_Builder) IsJavaNil_() bool { return this_ == nil }
func (*ModuleLoader_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.ModuleLoader$Builder", "Builder"
}
func (this_ *ModuleLoader_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ModuleLoader_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.ModuleLoader.Builder", "com.google.javascript.jscomp.deps.ModuleLoader$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*ModuleLoader_Builder](v_)
	})
}

func (this_ *ModuleLoader_Builder) initFields_ModuleLoader_Builder() {
	_ = this_
	this_.errorHandler = StaticGet_ModuleLoader_NOOP_ERROR_HANDER()
	this_.pathResolver = StaticGet_ModuleLoader_PathResolver_RELATIVE()
	this_.pathEscaper = StaticGet_ModuleLoader_PathEscaper_ESCAPE()
}

// Java constructor: Builder()
func NewModuleLoader_Builder() *ModuleLoader_Builder {
	this_ := &ModuleLoader_Builder{}
	InitModuleLoader_Builder(this_)
	return this_
}

func InitModuleLoader_Builder(this_ *ModuleLoader_Builder) {
	this_.initFields_ModuleLoader_Builder()
}

// Java method: Builder.setErrorHandler(com.google.javascript.jscomp.ErrorHandler x)
func (this_ *ModuleLoader_Builder) SetErrorHandler(x ErrorHandler) *ModuleLoader_Builder {
	_ = x
	if x != nil && !x.IsJavaNil_() {
		this_.errorHandler = x
	}
	return this_
}

// Java method: Builder.setModuleRoots(java.lang.Iterable<java.lang.String> x)
func (this_ *ModuleLoader_Builder) SetModuleRoots(x j2gort.Ext_java_lang_Iterable) *ModuleLoader_Builder {
	_ = x
	this_.moduleRoots = x
	return this_
}

// Java method: Builder.setInputs(java.lang.Iterable<? extends com.google.javascript.jscomp.deps.DependencyInfo> x)
func (this_ *ModuleLoader_Builder) SetInputs(x j2gort.Ext_java_lang_Iterable) *ModuleLoader_Builder {
	_ = x
	this_.inputs = x
	return this_
}

// Java method: Builder.setFactory(com.google.javascript.jscomp.deps.ModuleLoader.ModuleResolverFactory x)
func (this_ *ModuleLoader_Builder) SetFactory(x ModuleLoader_ModuleResolverFactory) *ModuleLoader_Builder {
	_ = x
	this_.factory = x
	return this_
}

// Java method: Builder.setPathResolver(com.google.javascript.jscomp.deps.ModuleLoader.PathResolver x)
func (this_ *ModuleLoader_Builder) SetPathResolver(x *ModuleLoader_PathResolver) *ModuleLoader_Builder {
	_ = x
	this_.pathResolver = x
	return this_
}

// Java method: Builder.setPathEscaper(com.google.javascript.jscomp.deps.ModuleLoader.PathEscaper x)
func (this_ *ModuleLoader_Builder) SetPathEscaper(x *ModuleLoader_PathEscaper) *ModuleLoader_Builder {
	_ = x
	this_.pathEscaper = x
	return this_
}

// Java method: Builder.build()
func (this_ *ModuleLoader_Builder) Build() *ModuleLoader {
	return NewModuleLoader(this_)
}

// Java class: com.google.javascript.jscomp.deps.ModuleLoader.ModulePath
type ModuleLoader_ModulePath struct {
	outer_  *ModuleLoader
	path    string `gson:"path"`
	idHash_ atomic.Int32
}

func (this_ *ModuleLoader_ModulePath) BindSelf_(s any) {
	_ = s
}
func (this_ *ModuleLoader_ModulePath) IsJavaNil_() bool { return this_ == nil }
func (*ModuleLoader_ModulePath) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.ModuleLoader$ModulePath", "ModulePath"
}
func (this_ *ModuleLoader_ModulePath) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ModuleLoader_ModulePath) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.ModuleLoader.ModulePath", "com.google.javascript.jscomp.deps.ModuleLoader$ModulePath", func(v_ any) bool {
		return j2gort.InstanceOf[*ModuleLoader_ModulePath](v_)
	})
}
func (this_ *ModuleLoader_ModulePath) SetOuter_(o *ModuleLoader) { this_.outer_ = o }
func (this_ *ModuleLoader_ModulePath) GetOuter_() *ModuleLoader  { return this_.outer_ }

func (this_ *ModuleLoader_ModulePath) initFields_ModuleLoader_ModulePath() {
	_ = this_
	this_.path = j2gort.NullString
}

// Java constructor: ModulePath(java.lang.String path)
func NewModuleLoader_ModulePath(outer_ *ModuleLoader, path string) *ModuleLoader_ModulePath {
	this_ := &ModuleLoader_ModulePath{}
	this_.outer_ = outer_
	InitModuleLoader_ModulePath(this_, path)
	return this_
}

func InitModuleLoader_ModulePath(this_ *ModuleLoader_ModulePath, path string) {
	_ = path
	this_.initFields_ModuleLoader_ModulePath()
	this_.path = path
}

// Java method: ModulePath.toString()
func (this_ *ModuleLoader_ModulePath) String() string {
	return this_.path
}

// Java method: ModulePath.equalsIgnoreLeadingSlash(com.google.javascript.jscomp.deps.ModuleLoader.ModulePath other)
func (this_ *ModuleLoader_ModulePath) EqualsIgnoreLeadingSlash(other *ModuleLoader_ModulePath) bool {
	_ = other
	return ((other != nil) && func() bool {
		sr_7 := this_.ToModuleName()
		se_8 := other.ToModuleName()
		return j2gort.CheckNotNullString(sr_7) == se_8
	}())
}

// Java method: ModulePath.toModuleName()
func (this_ *ModuleLoader_ModulePath) ToModuleName() string {
	return ModuleNames_ToModuleName(this_.path)
}

// Java method: ModulePath.resolveJsModule(java.lang.String moduleAddress, java.lang.String sourcename, int lineno, int colno)
func (this_ *ModuleLoader_ModulePath) ResolveJsModule(moduleAddress string, sourcename string, lineno int32, colno int32) *ModuleLoader_ModulePath {
	_ = moduleAddress
	_ = sourcename
	_ = lineno
	_ = colno
	var loadAddress string = this_.outer_.moduleResolver.ResolveJsModule(this_.path, moduleAddress, sourcename, lineno, colno)
	_ = loadAddress
	if loadAddress != j2gort.NullString {
		return NewModuleLoader_ModulePath(this_.outer_, loadAddress)
	}
	return nil
}

// Java method: ModulePath.resolveJsModuleSilently(java.lang.String moduleAddress)
func (this_ *ModuleLoader_ModulePath) ResolveJsModuleSilently(moduleAddress string) *ModuleLoader_ModulePath {
	_ = moduleAddress
	var loadAddress string = this_.outer_.moduleResolver.ResolveJsModuleSilently(this_.path, moduleAddress)
	_ = loadAddress
	if loadAddress != j2gort.NullString {
		return NewModuleLoader_ModulePath(this_.outer_, loadAddress)
	}
	return nil
}

// Java method: ModulePath.resolveModuleAsPath(java.lang.String moduleAddress)
func (this_ *ModuleLoader_ModulePath) ResolveModuleAsPath(moduleAddress string) *ModuleLoader_ModulePath {
	_ = moduleAddress
	return NewModuleLoader_ModulePath(this_.outer_, this_.outer_.moduleResolver.ResolveModuleAsPath(this_.path, moduleAddress))
}

// Java enum: com.google.javascript.jscomp.deps.ModuleLoader.PathEscaper
type ModuleLoader_PathEscaper struct {
	j2gort.EnumBase
	fn_Escape func(*ModuleLoader_PathEscaper, string) string
}

func (this_ *ModuleLoader_PathEscaper) GetEnumBase() *j2gort.EnumBase {
	if this_ == nil {
		return nil
	}
	return &this_.EnumBase
}
func (this_ *ModuleLoader_PathEscaper) EnumOrdinal_() int32 {
	if this_ == nil {
		return -1
	}
	return this_.EnumBase.Ordinal_
}
func (*ModuleLoader_PathEscaper) Implements_j2gort_Ext_java_lang_Enum()               {}
func (*ModuleLoader_PathEscaper) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ *ModuleLoader_PathEscaper) initFields_ModuleLoader_PathEscaper() {
	_ = this_
}

func initEnumModuleLoader_PathEscaper(this_ *ModuleLoader_PathEscaper, enumName_ string, enumOrd_ int32) *ModuleLoader_PathEscaper {
	this_.EnumBase.Name_ = enumName_
	this_.EnumBase.Ordinal_ = enumOrd_
	this_.initFields_ModuleLoader_PathEscaper()
	return this_
}

var enum_ModuleLoader_PathEscaper_ESCAPE *ModuleLoader_PathEscaper

func StaticGet_ModuleLoader_PathEscaper_ESCAPE() *ModuleLoader_PathEscaper {
	clinit_ModuleLoader_PathEscaper()
	return enum_ModuleLoader_PathEscaper_ESCAPE
}

var enum_ModuleLoader_PathEscaper_CANONICALIZE_ONLY *ModuleLoader_PathEscaper

func StaticGet_ModuleLoader_PathEscaper_CANONICALIZE_ONLY() *ModuleLoader_PathEscaper {
	clinit_ModuleLoader_PathEscaper()
	return enum_ModuleLoader_PathEscaper_CANONICALIZE_ONLY
}

var clinitDone_ModuleLoader_PathEscaper uint32
var clinitState_ModuleLoader_PathEscaper j2gort.ClinitState

func clinit_ModuleLoader_PathEscaper() {
	if atomic.LoadUint32(&clinitDone_ModuleLoader_PathEscaper) != 0 {
		return
	}
	clinitSlow_ModuleLoader_PathEscaper()
}

func clinitSlow_ModuleLoader_PathEscaper() {
	if !j2gort.EnterClinit(&clinitState_ModuleLoader_PathEscaper, &clinitDone_ModuleLoader_PathEscaper) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ModuleLoader_PathEscaper, &clinitDone_ModuleLoader_PathEscaper)
	enum_ModuleLoader_PathEscaper_ESCAPE = &ModuleLoader_PathEscaper{}
	enum_ModuleLoader_PathEscaper_ESCAPE.fn_Escape = func(this_ *ModuleLoader_PathEscaper, path string) string {
		_ = path
		return ModuleNames_EscapePath(path)
	}
	initEnumModuleLoader_PathEscaper(enum_ModuleLoader_PathEscaper_ESCAPE, "ESCAPE", int32(0))
	enum_ModuleLoader_PathEscaper_CANONICALIZE_ONLY = &ModuleLoader_PathEscaper{}
	enum_ModuleLoader_PathEscaper_CANONICALIZE_ONLY.fn_Escape = func(this_ *ModuleLoader_PathEscaper, path string) string {
		_ = path
		return ModuleNames_CanonicalizePath(path)
	}
	initEnumModuleLoader_PathEscaper(enum_ModuleLoader_PathEscaper_CANONICALIZE_ONLY, "CANONICALIZE_ONLY", int32(1))
	enumValues_ModuleLoader_PathEscaper = []*ModuleLoader_PathEscaper{enum_ModuleLoader_PathEscaper_ESCAPE, enum_ModuleLoader_PathEscaper_CANONICALIZE_ONLY}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.deps.ModuleLoader.PathEscaper", []any{enum_ModuleLoader_PathEscaper_ESCAPE, enum_ModuleLoader_PathEscaper_CANONICALIZE_ONLY})
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.deps.ModuleLoader$PathEscaper", []any{enum_ModuleLoader_PathEscaper_ESCAPE, enum_ModuleLoader_PathEscaper_CANONICALIZE_ONLY})
}

func init() {
	j2gort.RegisterEnumConstantsLoader("com.google.javascript.jscomp.deps.ModuleLoader.PathEscaper", "com.google.javascript.jscomp.deps.ModuleLoader$PathEscaper", func() {
		clinit_ModuleLoader_PathEscaper()
	})
}

var enumValues_ModuleLoader_PathEscaper []*ModuleLoader_PathEscaper

func ModuleLoader_PathEscaper_Values() []*ModuleLoader_PathEscaper {
	clinit_ModuleLoader_PathEscaper()
	return j2gort.CloneSlice(enumValues_ModuleLoader_PathEscaper)
}

func ModuleLoader_PathEscaper_ValueOf(name string) *ModuleLoader_PathEscaper {
	clinit_ModuleLoader_PathEscaper()
	for _, c := range enumValues_ModuleLoader_PathEscaper {
		if c.Name() == name {
			return c
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.deps.ModuleLoader.PathEscaper." + name))
}

// Java method: PathEscaper.escape(java.lang.String path)
func (this_ *ModuleLoader_PathEscaper) Escape(path string) string {
	if this_.fn_Escape != nil {
		return this_.fn_Escape(this_, path)
	}
	return this_.baseEnum_Escape(path)
}

func (this_ *ModuleLoader_PathEscaper) baseEnum_Escape(path string) string {
	_ = path
	panic("abstract or native method invoked")
}

func (this_ *ModuleLoader_PathEscaper) IsJavaNil_() bool { return this_ == nil }
func (*ModuleLoader_PathEscaper) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.ModuleLoader$PathEscaper", "PathEscaper"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.ModuleLoader.PathEscaper", "com.google.javascript.jscomp.deps.ModuleLoader$PathEscaper", func(v_ any) bool {
		return j2gort.InstanceOf[*ModuleLoader_PathEscaper](v_)
	})
}

// Java enum: com.google.javascript.jscomp.deps.ModuleLoader.PathResolver
type ModuleLoader_PathResolver struct {
	j2gort.EnumBase
	fn_Apply func(*ModuleLoader_PathResolver, any) any
}

func (this_ *ModuleLoader_PathResolver) GetEnumBase() *j2gort.EnumBase {
	if this_ == nil {
		return nil
	}
	return &this_.EnumBase
}
func (this_ *ModuleLoader_PathResolver) EnumOrdinal_() int32 {
	if this_ == nil {
		return -1
	}
	return this_.EnumBase.Ordinal_
}
func (*ModuleLoader_PathResolver) Implements_j2gort_Ext_java_lang_Enum()                  {}
func (*ModuleLoader_PathResolver) Implements_j2gort_Ext_com_google_common_base_Function() {}
func (*ModuleLoader_PathResolver) Implements_j2gort_Ext_java_lang_constant_Constable()    {}
func (*ModuleLoader_PathResolver) Implements_j2gort_Ext_java_util_function_Function()     {}
func (this_ *ModuleLoader_PathResolver) initFields_ModuleLoader_PathResolver() {
	_ = this_
}

func initEnumModuleLoader_PathResolver(this_ *ModuleLoader_PathResolver, enumName_ string, enumOrd_ int32) *ModuleLoader_PathResolver {
	this_.EnumBase.Name_ = enumName_
	this_.EnumBase.Ordinal_ = enumOrd_
	this_.initFields_ModuleLoader_PathResolver()
	return this_
}

var enum_ModuleLoader_PathResolver_RELATIVE *ModuleLoader_PathResolver

func StaticGet_ModuleLoader_PathResolver_RELATIVE() *ModuleLoader_PathResolver {
	clinit_ModuleLoader_PathResolver()
	return enum_ModuleLoader_PathResolver_RELATIVE
}

var enum_ModuleLoader_PathResolver_ABSOLUTE *ModuleLoader_PathResolver

func StaticGet_ModuleLoader_PathResolver_ABSOLUTE() *ModuleLoader_PathResolver {
	clinit_ModuleLoader_PathResolver()
	return enum_ModuleLoader_PathResolver_ABSOLUTE
}

var clinitDone_ModuleLoader_PathResolver uint32
var clinitState_ModuleLoader_PathResolver j2gort.ClinitState

func clinit_ModuleLoader_PathResolver() {
	if atomic.LoadUint32(&clinitDone_ModuleLoader_PathResolver) != 0 {
		return
	}
	clinitSlow_ModuleLoader_PathResolver()
}

func clinitSlow_ModuleLoader_PathResolver() {
	if !j2gort.EnterClinit(&clinitState_ModuleLoader_PathResolver, &clinitDone_ModuleLoader_PathResolver) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ModuleLoader_PathResolver, &clinitDone_ModuleLoader_PathResolver)
	enum_ModuleLoader_PathResolver_RELATIVE = &ModuleLoader_PathResolver{}
	enum_ModuleLoader_PathResolver_RELATIVE.fn_Apply = func(this_ *ModuleLoader_PathResolver, arg_9 any) any {
		path := j2gort.CastString(arg_9)
		_ = path
		return j2gort.BoxString(path)
	}
	initEnumModuleLoader_PathResolver(enum_ModuleLoader_PathResolver_RELATIVE, "RELATIVE", int32(0))
	enum_ModuleLoader_PathResolver_ABSOLUTE = &ModuleLoader_PathResolver{}
	enum_ModuleLoader_PathResolver_ABSOLUTE.fn_Apply = func(this_ *ModuleLoader_PathResolver, arg_10 any) any {
		path := j2gort.CastString(arg_10)
		_ = path
		return j2gort.BoxString(j2gort.Ext_java_nio_file_Path_Of_2(path, []string{}).ToAbsolutePath().String())
	}
	initEnumModuleLoader_PathResolver(enum_ModuleLoader_PathResolver_ABSOLUTE, "ABSOLUTE", int32(1))
	enumValues_ModuleLoader_PathResolver = []*ModuleLoader_PathResolver{enum_ModuleLoader_PathResolver_RELATIVE, enum_ModuleLoader_PathResolver_ABSOLUTE}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.deps.ModuleLoader.PathResolver", []any{enum_ModuleLoader_PathResolver_RELATIVE, enum_ModuleLoader_PathResolver_ABSOLUTE})
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.deps.ModuleLoader$PathResolver", []any{enum_ModuleLoader_PathResolver_RELATIVE, enum_ModuleLoader_PathResolver_ABSOLUTE})
}

func init() {
	j2gort.RegisterEnumConstantsLoader("com.google.javascript.jscomp.deps.ModuleLoader.PathResolver", "com.google.javascript.jscomp.deps.ModuleLoader$PathResolver", func() {
		clinit_ModuleLoader_PathResolver()
	})
}

var enumValues_ModuleLoader_PathResolver []*ModuleLoader_PathResolver

func ModuleLoader_PathResolver_Values() []*ModuleLoader_PathResolver {
	clinit_ModuleLoader_PathResolver()
	return j2gort.CloneSlice(enumValues_ModuleLoader_PathResolver)
}

func ModuleLoader_PathResolver_ValueOf(name string) *ModuleLoader_PathResolver {
	clinit_ModuleLoader_PathResolver()
	for _, c := range enumValues_ModuleLoader_PathResolver {
		if c.Name() == name {
			return c
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.deps.ModuleLoader.PathResolver." + name))
}

func (this_ *ModuleLoader_PathResolver) Apply(p0_ any) any {
	if this_.fn_Apply != nil {
		return this_.fn_Apply(this_, p0_)
	}
	return this_.baseEnum_Apply(p0_)
}

func (this_ *ModuleLoader_PathResolver) baseEnum_Apply(p0_ any) any {
	_ = p0_
	panic("abstract or native method invoked")
}

func (this_ *ModuleLoader_PathResolver) IsJavaNil_() bool { return this_ == nil }
func (*ModuleLoader_PathResolver) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.ModuleLoader$PathResolver", "PathResolver"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.ModuleLoader.PathResolver", "com.google.javascript.jscomp.deps.ModuleLoader$PathResolver", func(v_ any) bool {
		return j2gort.InstanceOf[*ModuleLoader_PathResolver](v_)
	})
}

// Java interface: com.google.javascript.jscomp.deps.ModuleLoader.ModuleResolverFactory
type ModuleLoader_ModuleResolverFactory interface {
	Implements_jscomp_ModuleLoader_ModuleResolverFactory()
	IsJavaNil_() bool
	Create(j2gort.Ext_com_google_common_collect_ImmutableSet, j2gort.Ext_com_google_common_collect_ImmutableList, ErrorHandler, *ModuleLoader_PathEscaper) ModuleResolver
}

func ToIface_ModuleLoader_ModuleResolverFactory[T interface {
	ModuleLoader_ModuleResolverFactory
	comparable
}](v T) ModuleLoader_ModuleResolverFactory {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.ModuleLoader.ModuleResolverFactory", "com.google.javascript.jscomp.deps.ModuleLoader$ModuleResolverFactory", func(v_ any) bool {
		return j2gort.InstanceOf[ModuleLoader_ModuleResolverFactory](v_)
	})
}

type lambda_ModuleLoader_ModuleResolverFactory struct {
	fn func(j2gort.Ext_com_google_common_collect_ImmutableSet, j2gort.Ext_com_google_common_collect_ImmutableList, ErrorHandler, *ModuleLoader_PathEscaper) ModuleResolver
}

func (l *lambda_ModuleLoader_ModuleResolverFactory) IsJavaNil_() bool { return l == nil }
func NewLambda_ModuleLoader_ModuleResolverFactory(fn func(j2gort.Ext_com_google_common_collect_ImmutableSet, j2gort.Ext_com_google_common_collect_ImmutableList, ErrorHandler, *ModuleLoader_PathEscaper) ModuleResolver) ModuleLoader_ModuleResolverFactory {
	if fn == nil {
		return nil
	}
	return &lambda_ModuleLoader_ModuleResolverFactory{fn: fn}
}

func (*lambda_ModuleLoader_ModuleResolverFactory) Implements_jscomp_ModuleLoader_ModuleResolverFactory() {
}
func (l *lambda_ModuleLoader_ModuleResolverFactory) Create(p0_ j2gort.Ext_com_google_common_collect_ImmutableSet, p1_ j2gort.Ext_com_google_common_collect_ImmutableList, p2_ ErrorHandler, p3_ *ModuleLoader_PathEscaper) ModuleResolver {
	return l.fn(p0_, p1_, p2_, p3_)
}

// Java enum: com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode
type ModuleLoader_ResolutionMode uint8

const (
	enum_ModuleLoader_ResolutionMode_BROWSER                           ModuleLoader_ResolutionMode = 1
	enum_ModuleLoader_ResolutionMode_BROWSER_WITH_TRANSFORMED_PREFIXES ModuleLoader_ResolutionMode = 2
	enum_ModuleLoader_ResolutionMode_NODE                              ModuleLoader_ResolutionMode = 3
	enum_ModuleLoader_ResolutionMode_WEBPACK                           ModuleLoader_ResolutionMode = 4
)

var enumBaseTable_ModuleLoader_ResolutionMode = [...]j2gort.EnumBase{
	{Name_: "BROWSER", Ordinal_: 0},
	{Name_: "BROWSER_WITH_TRANSFORMED_PREFIXES", Ordinal_: 1},
	{Name_: "NODE", Ordinal_: 2},
	{Name_: "WEBPACK", Ordinal_: 3},
}
var enumValues_ModuleLoader_ResolutionMode = []ModuleLoader_ResolutionMode{1, 2, 3, 4}

func StaticGet_ModuleLoader_ResolutionMode_BROWSER() ModuleLoader_ResolutionMode {
	return enum_ModuleLoader_ResolutionMode_BROWSER
}
func StaticGet_ModuleLoader_ResolutionMode_BROWSER_WITH_TRANSFORMED_PREFIXES() ModuleLoader_ResolutionMode {
	return enum_ModuleLoader_ResolutionMode_BROWSER_WITH_TRANSFORMED_PREFIXES
}
func StaticGet_ModuleLoader_ResolutionMode_NODE() ModuleLoader_ResolutionMode {
	return enum_ModuleLoader_ResolutionMode_NODE
}
func StaticGet_ModuleLoader_ResolutionMode_WEBPACK() ModuleLoader_ResolutionMode {
	return enum_ModuleLoader_ResolutionMode_WEBPACK
}

func ModuleLoader_ResolutionMode_Values() []ModuleLoader_ResolutionMode {
	return j2gort.CloneSlice(enumValues_ModuleLoader_ResolutionMode)
}

func ModuleLoader_ResolutionMode_ValueOf(name string) ModuleLoader_ResolutionMode {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_ModuleLoader_ResolutionMode {
		if enumBaseTable_ModuleLoader_ResolutionMode[i_].Name_ == name {
			return ModuleLoader_ResolutionMode(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode." + name))
}

func (this_ ModuleLoader_ResolutionMode) IsJavaNil_() bool                              { return this_ == 0 }
func (ModuleLoader_ResolutionMode) HasIdentityEquals_()                                 {}
func (ModuleLoader_ResolutionMode) Implements_j2gort_Ext_java_lang_Enum()               {}
func (ModuleLoader_ResolutionMode) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ ModuleLoader_ResolutionMode) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (ModuleLoader_ResolutionMode) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.ModuleLoader$ResolutionMode", "ResolutionMode"
}
func (this_ ModuleLoader_ResolutionMode) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_ModuleLoader_ResolutionMode[this_-1]
}
func (this_ ModuleLoader_ResolutionMode) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ ModuleLoader_ResolutionMode) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ ModuleLoader_ResolutionMode) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ModuleLoader_ResolutionMode[this_-1].Name_
}
func (this_ ModuleLoader_ResolutionMode) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_ModuleLoader_ResolutionMode[this_-1].Name_
}
func (this_ ModuleLoader_ResolutionMode) Equals(other_ any) bool {
	o_, ok_ := other_.(ModuleLoader_ResolutionMode)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ ModuleLoader_ResolutionMode) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_ModuleLoader_ResolutionMode[this_-1].HashCode()
}
func (this_ ModuleLoader_ResolutionMode) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ModuleLoader_ResolutionMode[this_-1].CompareTo(other_)
}
func (this_ ModuleLoader_ResolutionMode) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_ModuleLoader_ResolutionMode[this_-1].Name_)
}
func init() {
	vals_ := []any{ModuleLoader_ResolutionMode(1), ModuleLoader_ResolutionMode(2), ModuleLoader_ResolutionMode(3), ModuleLoader_ResolutionMode(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.deps.ModuleLoader$ResolutionMode", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.ModuleLoader.ResolutionMode", "com.google.javascript.jscomp.deps.ModuleLoader$ResolutionMode", func(v_ any) bool {
		return j2gort.InstanceOf[ModuleLoader_ResolutionMode](v_)
	})
}
