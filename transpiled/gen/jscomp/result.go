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
// Port of com/google/javascript/jscomp/Result.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java class: com.google.javascript.jscomp.Result
type Result struct {
	errors j2gort.Ext_com_google_common_collect_ImmutableList `gson:"errors"`
	warnings j2gort.Ext_com_google_common_collect_ImmutableList `gson:"warnings"`
	variableMap *VariableMap `gson:"variableMap"`
	propertyMap *VariableMap `gson:"propertyMap"`
	namedAnonFunctionMap *VariableMap `gson:"namedAnonFunctionMap"`
	stringMap *VariableMap `gson:"stringMap"`
	instrumentationMappings *VariableMap `gson:"instrumentationMappings"`
	sourceMap *SourceMap `gson:"sourceMap"`
	cssNames j2gort.Ext_java_util_Set `gson:"cssNames"`
	externExport string `gson:"externExport"`
	idGeneratorMap string `gson:"idGeneratorMap"`
	idHash_ atomic.Int32
	success bool `gson:"success"`
	transpiledFiles bool `gson:"transpiledFiles"`
}

func (this_ *Result) BindSelf_(s any) {
	_ = s
}
func (this_ *Result) IsJavaNil_() bool { return this_ == nil }
func (*Result) JavaClassInfo_() (string, string) { return "com.google.javascript.jscomp.Result", "Result" }
func (this_ *Result) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Result) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Result", "com.google.javascript.jscomp.Result", func(v_ any) bool {
		return j2gort.InstanceOf[*Result](v_)
	})
}

func (this_ *Result) initFields_Result() {
	_ = this_
	this_.externExport = j2gort.NullString
	this_.idGeneratorMap = j2gort.NullString
}

// Java constructor: Result(com.google.common.collect.ImmutableList<com.google.javascript.jscomp.JSError> errors, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.JSError> warnings, com.google.javascript.jscomp.VariableMap variableMap, com.google.javascript.jscomp.VariableMap propertyMap, com.google.javascript.jscomp.VariableMap namedAnonFunctionMap, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable VariableMap stringMap, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable VariableMap instrumentationMappings, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable SourceMap sourceMap, java.lang.String externExport, java.util.@org.jspecify.annotations.Nullable Set<java.lang.String> cssNames, java.lang.@org.jspecify.annotations.Nullable String idGeneratorMap, boolean transpiledFiles)
func NewResult_12(errors j2gort.Ext_com_google_common_collect_ImmutableList, warnings j2gort.Ext_com_google_common_collect_ImmutableList, variableMap *VariableMap, propertyMap *VariableMap, namedAnonFunctionMap *VariableMap, stringMap *VariableMap, instrumentationMappings *VariableMap, sourceMap *SourceMap, externExport string, cssNames j2gort.Ext_java_util_Set, idGeneratorMap string, transpiledFiles bool) *Result {
	this_ := &Result{}
	InitResult_12(this_, errors, warnings, variableMap, propertyMap, namedAnonFunctionMap, stringMap, instrumentationMappings, sourceMap, externExport, cssNames, idGeneratorMap, transpiledFiles)
	return this_
}

func InitResult_12(this_ *Result, errors j2gort.Ext_com_google_common_collect_ImmutableList, warnings j2gort.Ext_com_google_common_collect_ImmutableList, variableMap *VariableMap, propertyMap *VariableMap, namedAnonFunctionMap *VariableMap, stringMap *VariableMap, instrumentationMappings *VariableMap, sourceMap *SourceMap, externExport string, cssNames j2gort.Ext_java_util_Set, idGeneratorMap string, transpiledFiles bool) {
	_ = errors
	_ = warnings
	_ = variableMap
	_ = propertyMap
	_ = namedAnonFunctionMap
	_ = stringMap
	_ = instrumentationMappings
	_ = sourceMap
	_ = externExport
	_ = cssNames
	_ = idGeneratorMap
	_ = transpiledFiles
	this_.initFields_Result()
	this_.success = errors.IsEmpty()
	this_.errors = errors
	this_.warnings = warnings
	this_.variableMap = variableMap
	this_.propertyMap = propertyMap
	this_.namedAnonFunctionMap = namedAnonFunctionMap
	this_.stringMap = stringMap
	this_.instrumentationMappings = instrumentationMappings
	this_.sourceMap = sourceMap
	this_.externExport = externExport
	this_.cssNames = cssNames
	this_.idGeneratorMap = idGeneratorMap
	this_.transpiledFiles = transpiledFiles
}

// Java constructor: Result(com.google.common.collect.ImmutableList<com.google.javascript.jscomp.JSError> errors, com.google.common.collect.ImmutableList<com.google.javascript.jscomp.JSError> warnings, com.google.javascript.jscomp.VariableMap variableMap, com.google.javascript.jscomp.VariableMap propertyMap, com.google.javascript.jscomp.VariableMap namedAnonFunctionMap, com.google.javascript.jscomp.SourceMap sourceMap, java.lang.String externExport)
func NewResult_7(errors j2gort.Ext_com_google_common_collect_ImmutableList, warnings j2gort.Ext_com_google_common_collect_ImmutableList, variableMap *VariableMap, propertyMap *VariableMap, namedAnonFunctionMap *VariableMap, sourceMap *SourceMap, externExport string) *Result {
	this_ := &Result{}
	InitResult_7(this_, errors, warnings, variableMap, propertyMap, namedAnonFunctionMap, sourceMap, externExport)
	return this_
}

func InitResult_7(this_ *Result, errors j2gort.Ext_com_google_common_collect_ImmutableList, warnings j2gort.Ext_com_google_common_collect_ImmutableList, variableMap *VariableMap, propertyMap *VariableMap, namedAnonFunctionMap *VariableMap, sourceMap *SourceMap, externExport string) {
	_ = errors
	_ = warnings
	_ = variableMap
	_ = propertyMap
	_ = namedAnonFunctionMap
	_ = sourceMap
	_ = externExport
	InitResult_12(this_, errors, warnings, variableMap, propertyMap, namedAnonFunctionMap, nil, nil, sourceMap, externExport, nil, j2gort.NullString, false)
}

// Java static method: Result.pruneResultForPartialCompilation(com.google.javascript.jscomp.Result result)
func Result_PruneResultForPartialCompilation(result *Result) *Result {
	_ = result
	var val_1 VariableMap
	InitVariableMap(&val_1, j2gort.AsIface[j2gort.Ext_java_util_Map](j2gort.Ext_com_google_common_collect_ImmutableMap_Of_0()))
	var emptyVariableMap *VariableMap = &val_1
	_ = emptyVariableMap
	return NewResult_12(result.errors, result.warnings, emptyVariableMap, emptyVariableMap, emptyVariableMap, nil, emptyVariableMap, nil, "", nil, j2gort.NullString, false)
}

