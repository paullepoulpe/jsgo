/*
 * Copyright 2008 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/deps/DependencyFile.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.deps.DependencyFile
type DependencyFile struct {
	dependencies j2gort.PlainMap[string, DependencyInfo] `gson:"dependencies"`
	delegate     Deps_SourceFile                         `gson:"delegate"`
	idHash_      atomic.Int32
}

func (this_ *DependencyFile) BindSelf_(s any) {
	_ = s
}
func (this_ *DependencyFile) IsJavaNil_() bool { return this_ == nil }
func (*DependencyFile) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.DependencyFile", "DependencyFile"
}
func (this_ *DependencyFile) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DependencyFile) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.DependencyFile", "com.google.javascript.jscomp.deps.DependencyFile", func(v_ any) bool {
		return j2gort.InstanceOf[*DependencyFile](v_)
	})
}
func (*DependencyFile) Implements_jscomp_Deps_SourceFile() {}

var static_DependencyFile_logger *j2gort.Ext_java_util_logging_Logger

func StaticGet_DependencyFile_Logger() *j2gort.Ext_java_util_logging_Logger {
	clinit_DependencyFile()
	return static_DependencyFile_logger
}
func StaticSet_DependencyFile_Logger(v *j2gort.Ext_java_util_logging_Logger) {
	clinit_DependencyFile()
	static_DependencyFile_logger = v
}

var clinitDone_DependencyFile uint32
var clinitState_DependencyFile j2gort.ClinitState

func clinit_DependencyFile() {
	if atomic.LoadUint32(&clinitDone_DependencyFile) != 0 {
		return
	}
	clinitSlow_DependencyFile()
}

func clinitSlow_DependencyFile() {
	if !j2gort.EnterClinit(&clinitState_DependencyFile, &clinitDone_DependencyFile) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_DependencyFile, &clinitDone_DependencyFile)
	static_DependencyFile_logger = j2gort.Ext_java_util_logging_Logger_GetLogger_1(j2gort.CheckNotNullPtr(j2gort.ClassLiteral("com.google.javascript.jscomp.deps.DependencyFile", "DependencyFile")).GetName())
}

func (this_ *DependencyFile) initFields_DependencyFile() {
	_ = this_
	this_.dependencies = j2gort.NewPlainMapVal[string, DependencyInfo]()
}

// Java constructor: DependencyFile(com.google.javascript.jscomp.deps.SourceFile delegate)
func NewDependencyFile(delegate Deps_SourceFile) *DependencyFile {
	clinit_DependencyFile()
	this_ := &DependencyFile{}
	InitDependencyFile(this_, delegate)
	return this_
}

func InitDependencyFile(this_ *DependencyFile, delegate Deps_SourceFile) {
	_ = delegate
	this_.initFields_DependencyFile()
	this_.delegate = delegate
}

// Java method: DependencyFile.getName()
func (this_ *DependencyFile) GetName() string {
	return this_.delegate.GetName()
}

// Java method: DependencyFile.getContent()
func (this_ *DependencyFile) GetContent() string {
	return this_.delegate.GetContent()
}

// Java method: DependencyFile.wasModified()
func (this_ *DependencyFile) WasModified() bool {
	return this_.delegate.WasModified()
}

// Java method: DependencyFile.ensureUpToDate()
func (this_ *DependencyFile) EnsureUpToDate() {
	if this_.dependencies.IsEmpty() || this_.WasModified() {
		this_.priv_DependencyFile_loadGraph_6112()
	}
}

// Java method: DependencyFile.getDependencyInfo(java.lang.String symbol)
func (this_ *DependencyFile) GetDependencyInfo(symbol string) DependencyInfo {
	_ = symbol
	return this_.dependencies.Get(symbol)
}

// Java method: DependencyFile.loadGraph()
func (this_ *DependencyFile) priv_DependencyFile_loadGraph_6112() {
	this_.dependencies.Clear()
	StaticGet_DependencyFile_Logger().Info_1_string("Loading dependency graph")
	var errorManager ErrorManager = NewLoggerErrorManager_1(StaticGet_DependencyFile_Logger())
	_ = errorManager
	var parser *DepsFileRegexParser = NewDepsFileRegexParser_1(errorManager)
	_ = parser
	var depInfos j2gort.Ext_java_util_List = parser.ParseFile_2(this_.GetName(), this_.GetContent())
	_ = depInfos
	if !(parser.DidParseSucceed()) {
		panic(NewServiceException(("Problem parsing " + j2gort.StrValString(this_.GetName()) + ". See logs for details.")))
	}
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(depInfos))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			depInfo := j2gort.AsIface[DependencyInfo](rawElem_2)
			_ = depInfo
			{
				cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(depInfo.GetProvides()))
				for cur_3.HasNext() {
					provide := cur_3.NextStr()
					_ = provide
					var existing DependencyInfo = this_.dependencies.Get(provide)
					_ = existing
					if (existing != nil && !existing.IsJavaNil_()) && !(javalang.ObjectsEquals(j2gort.CheckNotNull(existing), depInfo)) {
						panic(NewServiceException(("Duplicate provide of " + j2gort.StrValString(provide) + ". Was provided by " + j2gort.StrValString(existing.GetPathRelativeToClosureBase()) + " and " + j2gort.StrValString(depInfo.GetPathRelativeToClosureBase()))))
					}
					this_.dependencies.PutVoid(provide, depInfo)
				}
			}
		}
	}
	var provides j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = provides
	provides.Add_1("goog")
	this_.dependencies.PutVoid("goog", SimpleDependencyInfo_S_Builder("base.js", "base.js").SetProvides_1_Ext_java_util_Collection(provides).Build())
	errorManager.GenerateReport()
	StaticGet_DependencyFile_Logger().Info_1_string("Dependencies loaded")
}
