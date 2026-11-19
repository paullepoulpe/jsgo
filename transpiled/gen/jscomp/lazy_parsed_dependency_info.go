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
// Port of com/google/javascript/jscomp/LazyParsedDependencyInfo.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
)

// Java class: com.google.javascript.jscomp.LazyParsedDependencyInfo
type LazyParsedDependencyInfo struct {
	delegate  DependencyInfo                                    `gson:"delegate"`
	ast       *CompilerInput                                    `gson:"ast"`
	compiler  AbstractCompiler                                  `gson:"-"`
	loadFlags j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"loadFlags"`
	idHash_   atomic.Int32
}

func (this_ *LazyParsedDependencyInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *LazyParsedDependencyInfo) IsJavaNil_() bool { return this_ == nil }
func (*LazyParsedDependencyInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.LazyParsedDependencyInfo", "LazyParsedDependencyInfo"
}
func (this_ *LazyParsedDependencyInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*LazyParsedDependencyInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.LazyParsedDependencyInfo", "com.google.javascript.jscomp.LazyParsedDependencyInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*LazyParsedDependencyInfo](v_)
	})
}
func (*LazyParsedDependencyInfo) Implements_jscomp_DependencyInfo() {}

func (this_ *LazyParsedDependencyInfo) initFields_LazyParsedDependencyInfo() {
	_ = this_
}

// Java constructor: LazyParsedDependencyInfo(com.google.javascript.jscomp.deps.DependencyInfo delegate, com.google.javascript.jscomp.CompilerInput ast, com.google.javascript.jscomp.AbstractCompiler compiler)
func NewLazyParsedDependencyInfo(delegate DependencyInfo, ast *CompilerInput, compiler AbstractCompiler) *LazyParsedDependencyInfo {
	this_ := &LazyParsedDependencyInfo{}
	InitLazyParsedDependencyInfo(this_, delegate, ast, compiler)
	return this_
}

func InitLazyParsedDependencyInfo(this_ *LazyParsedDependencyInfo, delegate DependencyInfo, ast *CompilerInput, compiler AbstractCompiler) {
	_ = delegate
	_ = ast
	_ = compiler
	this_.initFields_LazyParsedDependencyInfo()
	this_.delegate = j2gort.CheckNotNullIface(delegate)
	this_.ast = j2gort.CheckNotNullPtr(ast)
	this_.compiler = j2gort.CheckNotNullIface(compiler)
}

// Java method: LazyParsedDependencyInfo.isEs6Module()
func (this_ *LazyParsedDependencyInfo) IsEs6Module() bool {
	return this_.delegate.IsEs6Module()
}

// Java method: LazyParsedDependencyInfo.isGoogModule()
func (this_ *LazyParsedDependencyInfo) IsGoogModule() bool {
	return this_.delegate.IsGoogModule()
}

// Java method: LazyParsedDependencyInfo.getLoadFlags()
func (this_ *LazyParsedDependencyInfo) GetLoadFlags() j2gort.Ext_com_google_common_collect_ImmutableMap {
	if this_.loadFlags == nil {
		var loadFlagsBuilder j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_TreeMap()
		_ = loadFlagsBuilder
		var hadSourceInMemory bool = this_.ast.GetSourceFile().HasSourceInMemory()
		_ = hadSourceInMemory
		loadFlagsBuilder.PutAll(j2gort.AsIface[j2gort.Ext_java_util_Map](this_.delegate.GetLoadFlags()))
		var features *pkg_parser.FeatureSet = this_.ast.GetFeatures(this_.compiler)
		_ = features
		if features.Has(pkg_parser.StaticGet_FeatureSet_Feature_MODULES()) {
			var previousModule string = j2gort.CastString(j2gort.MapGetStr(loadFlagsBuilder, "module"))
			_ = previousModule
			if (previousModule != j2gort.NullString) && !(j2gort.CheckNotNullString(previousModule) == "es6") {
				this_.compiler.Report_1(JSError_Make_5(this_.GetName(), int32(-1), int32(-1), StaticGet_ModuleLoader_MODULE_CONFLICT(), []string{this_.GetName()}))
			}
			j2gort.MapPutStr(loadFlagsBuilder, "module", "es6")
		}
		var version string = features.Version()
		_ = version
		if !(j2gort.CheckNotNullString(version) == "es3") {
			j2gort.MapPutStr(loadFlagsBuilder, "lang", j2gort.BoxString(version))
		}
		this_.loadFlags = j2gort.Ext_com_google_common_collect_ImmutableMap_CopyOf_1_Ext_java_util_Map(loadFlagsBuilder)
		if !(hadSourceInMemory) {
			this_.ast.GetSourceFile().ClearCachedSource()
		}
		this_.ast = nil
	}
	return this_.loadFlags
}

// Java method: LazyParsedDependencyInfo.getName()
func (this_ *LazyParsedDependencyInfo) GetName() string {
	return this_.delegate.GetName()
}

// Java method: LazyParsedDependencyInfo.getPathRelativeToClosureBase()
func (this_ *LazyParsedDependencyInfo) GetPathRelativeToClosureBase() string {
	return this_.delegate.GetPathRelativeToClosureBase()
}

// Java method: LazyParsedDependencyInfo.getRequires()
func (this_ *LazyParsedDependencyInfo) GetRequires() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.delegate.GetRequires()
}

// Java method: LazyParsedDependencyInfo.getTypeRequires()
func (this_ *LazyParsedDependencyInfo) GetTypeRequires() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.delegate.GetTypeRequires()
}

// Java method: LazyParsedDependencyInfo.getProvides()
func (this_ *LazyParsedDependencyInfo) GetProvides() j2gort.Ext_com_google_common_collect_ImmutableList {
	return this_.delegate.GetProvides()
}

// Java method: LazyParsedDependencyInfo.getHasExternsAnnotation()
func (this_ *LazyParsedDependencyInfo) GetHasExternsAnnotation() bool {
	return this_.delegate.GetHasExternsAnnotation()
}

// Java method: LazyParsedDependencyInfo.getHasNoCompileAnnotation()
func (this_ *LazyParsedDependencyInfo) GetHasNoCompileAnnotation() bool {
	return this_.delegate.GetHasNoCompileAnnotation()
}

func (this_ *LazyParsedDependencyInfo) GetRequiredSymbols() j2gort.Ext_com_google_common_collect_ImmutableList {
	return IfaceDefault_DependencyInfo_GetRequiredSymbols(this_)
}
