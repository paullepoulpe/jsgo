/*
 * Copyright 2006 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/deps/DefaultDependencyResolver.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.deps.DefaultDependencyResolver
type DefaultDependencyResolver struct {
	depsFiles      j2gort.Ext_java_util_List `gson:"depsFiles"`
	strictRequires bool                      `gson:"strictRequires"`
	idHash_        atomic.Int32
}

func (this_ *DefaultDependencyResolver) BindSelf_(s any) {
	_ = s
}
func (this_ *DefaultDependencyResolver) IsJavaNil_() bool { return this_ == nil }
func (*DefaultDependencyResolver) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.DefaultDependencyResolver", "DefaultDependencyResolver"
}
func (this_ *DefaultDependencyResolver) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DefaultDependencyResolver) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.DefaultDependencyResolver", "com.google.javascript.jscomp.deps.DefaultDependencyResolver", func(v_ any) bool {
		return j2gort.InstanceOf[*DefaultDependencyResolver](v_)
	})
}
func (*DefaultDependencyResolver) Implements_jscomp_DependencyResolver() {}

func StaticGet_DefaultDependencyResolver_CLOSURE_BASE() string         { return "base.js" }
func StaticGet_DefaultDependencyResolver_CLOSURE_BASE_PROVIDE() string { return "goog" }

var static_DefaultDependencyResolver_logger *j2gort.Ext_java_util_logging_Logger

func StaticGet_DefaultDependencyResolver_Logger() *j2gort.Ext_java_util_logging_Logger {
	clinit_DefaultDependencyResolver()
	return static_DefaultDependencyResolver_logger
}
func StaticSet_DefaultDependencyResolver_Logger(v *j2gort.Ext_java_util_logging_Logger) {
	clinit_DefaultDependencyResolver()
	static_DefaultDependencyResolver_logger = v
}

var clinitDone_DefaultDependencyResolver uint32
var clinitState_DefaultDependencyResolver j2gort.ClinitState

func clinit_DefaultDependencyResolver() {
	if atomic.LoadUint32(&clinitDone_DefaultDependencyResolver) != 0 {
		return
	}
	clinitSlow_DefaultDependencyResolver()
}

func clinitSlow_DefaultDependencyResolver() {
	if !j2gort.EnterClinit(&clinitState_DefaultDependencyResolver, &clinitDone_DefaultDependencyResolver) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_DefaultDependencyResolver, &clinitDone_DefaultDependencyResolver)
	static_DefaultDependencyResolver_logger = j2gort.Ext_java_util_logging_Logger_GetLogger_1(j2gort.CheckNotNullPtr(j2gort.ClassLiteral("com.google.javascript.jscomp.deps.DefaultDependencyResolver", "DefaultDependencyResolver")).GetName())
}

func (this_ *DefaultDependencyResolver) initFields_DefaultDependencyResolver() {
	_ = this_
}

// Java constructor: DefaultDependencyResolver(java.util.List<com.google.javascript.jscomp.deps.DependencyFile> depsFiles, boolean strictRequires)
func NewDefaultDependencyResolver(depsFiles j2gort.Ext_java_util_List, strictRequires bool) *DefaultDependencyResolver {
	clinit_DefaultDependencyResolver()
	this_ := &DefaultDependencyResolver{}
	InitDefaultDependencyResolver(this_, depsFiles, strictRequires)
	return this_
}

func InitDefaultDependencyResolver(this_ *DefaultDependencyResolver, depsFiles j2gort.Ext_java_util_List, strictRequires bool) {
	_ = depsFiles
	_ = strictRequires
	this_.initFields_DefaultDependencyResolver()
	this_.depsFiles = depsFiles
	this_.strictRequires = strictRequires
}

// Java method: DefaultDependencyResolver.getDependencies(java.lang.String code)
func (this_ *DefaultDependencyResolver) GetDependencies_1_string(code string) any {
	_ = code
	return j2gort.AsIface[j2gort.Ext_java_util_List](this_.GetDependencies_1_Ext_java_util_Collection(DefaultDependencyResolver_ParseRequires(code, true)))
}

// Java method: DefaultDependencyResolver.getDependencies(java.util.Collection<java.lang.String> symbols)
func (this_ *DefaultDependencyResolver) GetDependencies_1_Ext_java_util_Collection(symbols j2gort.Ext_java_util_Collection) any {
	_ = symbols
	return j2gort.AsIface[j2gort.Ext_java_util_List](this_.GetDependencies_2(symbols, j2gort.NewExt_java_util_HashSet()))
}

// Java method: DefaultDependencyResolver.getDependencies(java.lang.String code, java.util.Set<java.lang.String> seen, boolean addClosureBaseFile)
func (this_ *DefaultDependencyResolver) GetDependencies_3(code string, seen j2gort.Ext_java_util_Set, addClosureBaseFile bool) any {
	_ = code
	_ = seen
	_ = addClosureBaseFile
	return j2gort.AsIface[j2gort.Ext_java_util_List](this_.GetDependencies_2(DefaultDependencyResolver_ParseRequires(code, addClosureBaseFile), seen))
}

// Java method: DefaultDependencyResolver.getDependencies(java.util.Collection<java.lang.String> symbols, java.util.Set<java.lang.String> seen)
func (this_ *DefaultDependencyResolver) GetDependencies_2(symbols j2gort.Ext_java_util_Collection, seen j2gort.Ext_java_util_Set) any {
	_ = symbols
	_ = seen
	var list j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = list
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.depsFiles))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			depsFile := j2gort.CastPtr[DependencyFile](rawElem_2)
			_ = depsFile
			depsFile.EnsureUpToDate()
		}
	}
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(symbols))
		for cur_3.HasNext() {
			symbol := cur_3.NextStr()
			_ = symbol
			this_.priv_DefaultDependencyResolver_addDependency_6048(symbol, seen, list)
		}
	}
	return list
}

// Java method: DefaultDependencyResolver.addDependency(java.lang.String symbol, java.util.Set<java.lang.String> seen, java.util.List<java.lang.String> list)
func (this_ *DefaultDependencyResolver) priv_DefaultDependencyResolver_addDependency_6048(symbol string, seen j2gort.Ext_java_util_Set, list j2gort.Ext_java_util_List) {
	_ = symbol
	_ = seen
	_ = list
	var dependency DependencyInfo = this_.priv_DefaultDependencyResolver_getDependencyInfo_6049(symbol)
	_ = dependency
	if dependency == nil || dependency.IsJavaNil_() {
		if this_.strictRequires {
			panic(NewServiceException(("Unknown require of " + j2gort.StrValString(symbol))))
		}
	} else if !(seen.ContainsAll(j2gort.AsIface[j2gort.Ext_java_util_Collection](dependency.GetProvides()))) {
		seen.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](dependency.GetProvides()))
		var reqs j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_java_util_List](dependency.GetRequires())
		_ = reqs
		{
			var i int32 = int32(0)
			_ = i
			var len_ int32 = reqs.Size()
			_ = len_
			for ; i < len_; i = i + 1 {
				this_.priv_DefaultDependencyResolver_addDependency_6048(j2gort.CastPtr[DependencyInfo_Require](reqs.Get(i)).GetSymbol(), seen, list)
			}
		}
		list.Add_1(j2gort.BoxString(dependency.GetPathRelativeToClosureBase()))
	}
}

// Java static method: DefaultDependencyResolver.parseRequires(java.lang.String code, boolean addClosureBase)
func DefaultDependencyResolver_ParseRequires(code string, addClosureBase bool) j2gort.Ext_java_util_Collection {
	clinit_DefaultDependencyResolver()
	_ = code
	_ = addClosureBase
	var errorManager ErrorManager = NewLoggerErrorManager_1(StaticGet_DefaultDependencyResolver_Logger())
	_ = errorManager
	var parser *JsFileRegexParser = NewJsFileRegexParser(errorManager)
	_ = parser
	var deps DependencyInfo = parser.ParseFile("<unknown path>", "<unknown path>", code)
	_ = deps
	var requires []string = []string{}
	_ = requires
	if addClosureBase {
		requires = append(requires, "goog")
	}
	requires = j2gort.AppendStrCol(requires, j2gort.AsIface[j2gort.Ext_java_util_Collection](deps.GetRequiredSymbols()))
	errorManager.GenerateReport()
	return j2gort.WrapTypedArrayList(requires)
}

// Java method: DefaultDependencyResolver.getDependencyInfo(java.lang.String symbol)
func (this_ *DefaultDependencyResolver) priv_DefaultDependencyResolver_getDependencyInfo_6049(symbol string) DependencyInfo {
	_ = symbol
	{
		cur_4 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.depsFiles))
		for cur_4.HasNext() {
			rawElem_5 := cur_4.Next()
			depsFile := j2gort.CastPtr[DependencyFile](rawElem_5)
			_ = depsFile
			var di DependencyInfo = depsFile.GetDependencyInfo(symbol)
			_ = di
			if di != nil && !di.IsJavaNil_() {
				return di
			}
		}
	}
	return nil
}
