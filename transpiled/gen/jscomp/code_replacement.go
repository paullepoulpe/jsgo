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
// Port of com/google/javascript/refactoring/CodeReplacement.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

var lambdaSingleton_code_replacement_1 j2gort.Ext_java_util_function_Function
var lambdaSingleton_code_replacement_2 j2gort.Ext_java_util_function_Function
var lambdaSingleton_code_replacement_3 j2gort.Ext_java_util_function_ToIntFunction
var lambdaSingleton_code_replacement_4 j2gort.Ext_java_util_function_ToIntFunction

func init() {
	lambdaSingleton_code_replacement_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return j2gort.BoxString(j2gort.AsIface[CodeReplacement](mrp0_).GetNewContent()) })
	lambdaSingleton_code_replacement_2 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any { return j2gort.BoxString(j2gort.AsIface[CodeReplacement](mrp0_).GetSortKey()) })
	lambdaSingleton_code_replacement_3 = j2gort.NewLambda_Ext_java_util_function_ToIntFunction(func(mrp0_ any) int32 { return j2gort.AsIface[CodeReplacement](mrp0_).GetLength() })
	lambdaSingleton_code_replacement_4 = j2gort.NewLambda_Ext_java_util_function_ToIntFunction(func(mrp0_ any) int32 { return j2gort.AsIface[CodeReplacement](mrp0_).GetStartPosition() })
}

// Java polymorphic class interface: com.google.javascript.refactoring.CodeReplacement
type CodeReplacement interface {
	Implements_jscomp_CodeReplacement()
	IsJavaNil_() bool
	GetCodeReplacementBase() *CodeReplacementBase
	GetStartPosition() int32
	GetLength() int32
	GetEndPosition() int32
	GetNewContent() string
	GetSortKey() string
	CompareTo(any) int32
}

func ToIface_CodeReplacement[T interface {
	CodeReplacement
	comparable
}](v T) CodeReplacement {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.refactoring.CodeReplacement
type CodeReplacementBase struct {
	self_   CodeReplacement
	idHash_ atomic.Int32
}

func (this_ *CodeReplacementBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(CodeReplacement)
}
func (this_ *CodeReplacementBase) IsJavaNil_() bool { return this_ == nil }
func (*CodeReplacementBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.refactoring.CodeReplacement", "CodeReplacement"
}
func (this_ *CodeReplacementBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.refactoring.CodeReplacement", "com.google.javascript.refactoring.CodeReplacement", func(v_ any) bool {
		return j2gort.InstanceOf[CodeReplacement](v_)
	})
}
func (b *CodeReplacementBase) GetCodeReplacementBase() *CodeReplacementBase { return b }
func (*CodeReplacementBase) Implements_jscomp_CodeReplacement()             {}

var static_CodeReplacement_APPLICATION_ORDER j2gort.Ext_java_util_Comparator

func StaticGet_CodeReplacement_APPLICATION_ORDER() j2gort.Ext_java_util_Comparator {
	clinit_CodeReplacement()
	return static_CodeReplacement_APPLICATION_ORDER
}
func StaticSet_CodeReplacement_APPLICATION_ORDER(v j2gort.Ext_java_util_Comparator) {
	clinit_CodeReplacement()
	static_CodeReplacement_APPLICATION_ORDER = v
}

var clinitDone_CodeReplacement uint32
var clinitState_CodeReplacement j2gort.ClinitState

func clinit_CodeReplacement() {
	if atomic.LoadUint32(&clinitDone_CodeReplacement) != 0 {
		return
	}
	clinitSlow_CodeReplacement()
}

func clinitSlow_CodeReplacement() {
	if !j2gort.EnterClinit(&clinitState_CodeReplacement, &clinitDone_CodeReplacement) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CodeReplacement, &clinitDone_CodeReplacement)
	static_CodeReplacement_APPLICATION_ORDER = j2gort.Ext_java_util_Comparator_ComparingInt(lambdaSingleton_code_replacement_4).ThenComparingInt(lambdaSingleton_code_replacement_3).ThenComparing_1_Ext_java_util_function_Function(lambdaSingleton_code_replacement_2).ThenComparing_1_Ext_java_util_function_Function(lambdaSingleton_code_replacement_1)
}

func (this_ *CodeReplacementBase) initFields_CodeReplacement() {
	_ = this_
}

func InitCodeReplacement(this_ *CodeReplacementBase) {
	this_.initFields_CodeReplacement()
}

// Java static method: CodeReplacement.create(int startPosition, int length, java.lang.String newContent)
func CodeReplacement_Create_3(startPosition int32, length int32, newContent string) CodeReplacement {
	clinit_CodeReplacement()
	_ = startPosition
	_ = length
	_ = newContent
	return CodeReplacement_Create_4(startPosition, length, newContent, "")
}

// Java static method: CodeReplacement.create(int startPosition, int length, java.lang.String newContent, java.lang.String sortKey)
func CodeReplacement_Create_4(startPosition int32, length int32, newContent string, sortKey string) CodeReplacement {
	clinit_CodeReplacement()
	_ = startPosition
	_ = length
	_ = newContent
	_ = sortKey
	return NewAutoValue_CodeReplacement(startPosition, length, newContent, sortKey)
}

// Java method: CodeReplacement.getStartPosition()
func (this_ *CodeReplacementBase) GetStartPosition() int32 {
	panic("abstract or native method invoked")
}

// Java method: CodeReplacement.getLength()
func (this_ *CodeReplacementBase) GetLength() int32 {
	panic("abstract or native method invoked")
}

// Java method: CodeReplacement.getEndPosition()
func (this_ *CodeReplacementBase) GetEndPosition() int32 {
	return (this_.self_.GetStartPosition() + this_.self_.GetLength())
}

// Java method: CodeReplacement.getNewContent()
func (this_ *CodeReplacementBase) GetNewContent() string {
	panic("abstract or native method invoked")
}

// Java method: CodeReplacement.getSortKey()
func (this_ *CodeReplacementBase) GetSortKey() string {
	panic("abstract or native method invoked")
}

// Java method: CodeReplacement.compareTo(com.google.javascript.refactoring.CodeReplacement x)
func (this_ *CodeReplacementBase) CompareTo(arg_1 any) int32 {
	x := j2gort.AsIface[CodeReplacement](arg_1)
	_ = x
	return StaticGet_CodeReplacement_APPLICATION_ORDER().Compare(this_.self_, x)
}
