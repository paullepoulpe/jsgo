// Port of com/google/javascript/jscomp/CompilerOptions.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/CompilerOptionsPreprocessor.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/FlagUsageException.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/Platform.java @ closure-compiler v20261021

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

/*
 * Copyright 2018 The Closure Compiler Authors.
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

/*
 * Copyright 2020 The Closure Compiler Authors.
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
	"os"
	"runtime"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type defineReplacementKind uint8

const (
	defineReplacementBool defineReplacementKind = iota
	defineReplacementInt
	defineReplacementDouble
	defineReplacementString
)

type defineReplacement struct {
	kind    defineReplacementKind
	boolVal bool
	intVal  int32
	numVal  float64
	strVal  string
}

func (d defineReplacement) toNode() rhino.Node {
	switch d.kind {
	case defineReplacementBool:
		return NodeUtil_BooleanNode(d.boolVal)
	case defineReplacementInt:
		return NodeUtil_NumberNode(float64(d.intVal), rhino.NilNode)
	case defineReplacementDouble:
		return NodeUtil_NumberNode(d.numVal, rhino.NilNode)
	case defineReplacementString:
		return rhino.IR_String(d.strVal)
	default:
		panic(jutil.NewIllegalStateException())
	}
}

var compilerOptionsEmptyDefineReplacementsMap = jutil.NewLinkedHashMap[string, rhino.Node]()

// Returns the map of define replacements.
//
// Java: CompilerOptions.getDefineReplacements()
func (o *CompilerOptions) GetDefineReplacements() *jutil.LinkedHashMap[string, rhino.Node] {
	if len(o.defineReplacements.entries) == 0 {
		return compilerOptionsEmptyDefineReplacementsMap
	}
	m := jutil.NewLinkedHashMap[string, rhino.Node]()
	for _, e := range o.defineReplacements.entries {
		m.Put(e.Key, e.Value.toNode())
	}
	return m
}

// Sets the value of the @define variable in JS to a boolean literal.
//
// Java: CompilerOptions.setDefineToBooleanLiteral(String, boolean)
func (o *CompilerOptions) SetDefineToBooleanLiteral(defineName string, value bool) {
	o.defineReplacements.Put(defineName, defineReplacement{kind: defineReplacementBool, boolVal: value})
}

// Sets the value of the @define variable in JS to a String literal.
//
// Java: CompilerOptions.setDefineToStringLiteral(String, String)
func (o *CompilerOptions) SetDefineToStringLiteral(defineName string, value string) {
	o.defineReplacements.Put(defineName, defineReplacement{kind: defineReplacementString, strVal: value})
}

// Sets the value of the @define variable in JS to a number literal.
//
// Java: CompilerOptions.setDefineToNumberLiteral(String, int)
func (o *CompilerOptions) SetDefineToNumberLiteral(defineName string, value int) {
	o.defineReplacements.Put(defineName, defineReplacement{kind: defineReplacementInt, intVal: int32(value)})
}

// Sets the value of the @define variable in JS to a number literal.
//
// Java: CompilerOptions.setDefineToDoubleLiteral(String, double)
func (o *CompilerOptions) SetDefineToDoubleLiteral(defineName string, value float64) {
	o.defineReplacements.Put(defineName, defineReplacement{kind: defineReplacementDouble, numVal: value})
}

// Java: CompilerOptions.setUnknownDefinesToIgnore(ImmutableList<String>)
func (o *CompilerOptions) SetUnknownDefinesToIgnore(unknownDefinesToIgnore []string) {
	o.unknownDefinesToIgnore = unknownDefinesToIgnore
}

// Java: CompilerOptions.addUnknownDefinesToIgnore(ImmutableList<String>)
func (o *CompilerOptions) AddUnknownDefinesToIgnore(unknownDefinesToIgnore []string) {
	combined := make([]string, 0, len(o.unknownDefinesToIgnore)+len(unknownDefinesToIgnore))
	combined = append(combined, o.unknownDefinesToIgnore...)
	combined = append(combined, unknownDefinesToIgnore...)
	o.unknownDefinesToIgnore = combined
}

// Java: CompilerOptions.getUnknownDefinesToIgnore()
func (o *CompilerOptions) GetUnknownDefinesToIgnore() []string {
	return o.unknownDefinesToIgnore
}

// CompilerOptionsPreprocessor checks for combinations of options that are
// incompatible, i.e. will produce incorrect code.
//
// This is run by Compiler#compileInternal, which is not run during unit
// tests. The catch is that it's run after Compiler#initOptions, so if for
// example you want to change the warningsGuard, you can't do it here.
//
// Also, turns off options if the provided options don't make sense
// together.
//
// Don't instantiate (Java: private CompilerOptionsPreprocessor()).
type CompilerOptionsPreprocessor struct{}

// Java: CompilerOptionsPreprocessor.preprocess(CompilerOptions)
func CompilerOptionsPreprocessor_Preprocess(options *CompilerOptions) {
	if options.GetInlineFunctionsLevel() == CompilerOptions_Reach_NONE &&
		options.GetMaxFunctionSizeAfterInlining() != CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING {
		panic(newCompilerOptionsPreprocessor_InvalidOptionsException(
			"max_function_size_after_inlining has no effect if inlining is disabled."))
	}
}

// CompilerOptionsPreprocessor_InvalidOptionsException is an exception to
// indicate incompatible options in the CompilerOptions.
type CompilerOptionsPreprocessor_InvalidOptionsException struct {
	jutil.RuntimeException
}

// Java: CompilerOptionsPreprocessor.InvalidOptionsException(String)
func newCompilerOptionsPreprocessor_InvalidOptionsException(message string) *CompilerOptionsPreprocessor_InvalidOptionsException {
	e := &CompilerOptionsPreprocessor_InvalidOptionsException{}
	jutil.InitThrowable(e, "com.google.javascript.jscomp.CompilerOptionsPreprocessor$InvalidOptionsException", &message, nil)
	return e
}

// FlagUsageException is an exception thrown when command-line flags are used
// incorrectly.
type FlagUsageException struct{ jutil.RuntimeException }

func (*FlagUsageException) flagUsageException() {}

// IFlagUsageException matches FlagUsageException.
type IFlagUsageException interface {
	jutil.IRuntimeException
	flagUsageException()
}

// Java: FlagUsageException.serialVersionUID
const FlagUsageException_serialVersionUID int64 = 1

// Java: FlagUsageException(String)
func NewFlagUsageException(message string) *FlagUsageException {
	e := &FlagUsageException{}
	jutil.InitThrowable(e, "com.google.javascript.jscomp.FlagUsageException", &message, nil)
	return e
}

// Platform is a utility to abstract platform specific logic so it could be
// super-sourced for Web.
type Platform struct{}

// Java: Platform.getFileSeperator()
//
// Go-only: language-adapter — File.separator is os.PathSeparator.
func Platform_GetFileSeperator() string {
	return string(os.PathSeparator)
}

// Java: Platform.isThreadInterrupted()
//
// Go-only: language-adapter — goroutines cannot be interrupted, so the
// interrupted flag is never set.
func Platform_IsThreadInterrupted() bool {
	// Note that this clears interruption. Thread.isInterruped should be
	// preferred instead.
	return false
}

// Java: Platform.freeMemory()
//
// Go-only: language-adapter — Runtime.freeMemory() is the heap memory that
// is reserved but unused (HeapIdle).
func Platform_FreeMemory() int64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return int64(m.HeapIdle)
}

// Java: Platform.totalMemory()
//
// Go-only: language-adapter — Runtime.totalMemory() is the heap memory
// reserved from the OS (HeapSys).
func Platform_TotalMemory() int64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return int64(m.HeapSys)
}

// Java: Platform.formatMessage(String, String...)
func Platform_FormatMessage(message string, arguments ...string) string {
	// Note that MessageFormat is removing single quotes and in many cases
	// intended ones. Consider moving to a simpler formatting version like the
	// Web one.
	return parseDiagPattern(message).format(arguments...)
}
