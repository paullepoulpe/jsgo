// Port of com/google/javascript/jscomp/AbstractCompiler.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/CodeChangeHandler.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/RemoveWeakSources.java @ closure-compiler v20261021

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

package jscomp

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// AbstractCompiler aliases *Compiler so compiler passes dispatch directly on
// the concrete compiler pointer without a 109-method interface table. Parallel
// worker copies carry a non-nil recorder field that intercepts change/error
// reporting on the worker path.
type AbstractCompiler = *Compiler

// Java: AbstractCompiler.READ_ERROR
var AbstractCompiler_READ_ERROR = DiagnosticType_Error("JSC_READ_ERROR", "Cannot read file {0}: {1}")

// Java: AbstractCompiler.FILL_FILE_SUFFIX
const abstractCompiler_FILL_FILE_SUFFIX = "$fillFile"

// AbstractCompiler_CreateFillFileName: empty chunks get an empty "fill"
// file, so that we can move code into an empty chunk.
//
// Java: AbstractCompiler.createFillFileName(String)
func AbstractCompiler_CreateFillFileName(chunkName string) string {
	return chunkName + abstractCompiler_FILL_FILE_SUFFIX
}

// AbstractCompiler_IsFillFileName returns whether a file name was created by
// AbstractCompiler_CreateFillFileName.
//
// Java: AbstractCompiler.isFillFileName(String)
func AbstractCompiler_IsFillFileName(fileName string) bool {
	return strings.HasSuffix(fileName, abstractCompiler_FILL_FILE_SUFFIX)
}

// AbstractCompiler_ConfigContext represents the different contexts for
// which the compiler could have distinct configurations.
type AbstractCompiler_ConfigContext int32

const (
	// Normal JavaScript.
	AbstractCompiler_ConfigContext_DEFAULT AbstractCompiler_ConfigContext = iota

	// Externs files.
	AbstractCompiler_ConfigContext_EXTERNS
)

var abstractCompiler_ConfigContext_names = [...]string{"DEFAULT", "EXTERNS"}

// String is Enum.name().
func (e AbstractCompiler_ConfigContext) String() string {
	return abstractCompiler_ConfigContext_names[e]
}

// Ordinal is Enum.ordinal().
func (e AbstractCompiler_ConfigContext) Ordinal() int { return int(e) }

// AbstractCompiler_LifeCycleStage is what point in optimizations we're in.
// For use by compiler passes
type AbstractCompiler_LifeCycleStage int32

const (
	AbstractCompiler_LifeCycleStage_RAW AbstractCompiler_LifeCycleStage = iota

	// See constraints put on the AST by either running the
	// ConvertTypesToColors.java pass /or/ by having created the AST via
	// serialization/TypedAstDeserializer.java
	// NORMALIZED implies this constraint
	AbstractCompiler_LifeCycleStage_COLORS_AND_SIMPLIFIED_JSDOC

	// See constraints put on the tree by Normalize.java
	AbstractCompiler_LifeCycleStage_NORMALIZED

	// The normalize pass has put constraints on the tree,
	// but variables and properties have been renamed so
	// coding conventions no longer apply.
	AbstractCompiler_LifeCycleStage_NORMALIZED_OBFUSCATED
)

var abstractCompiler_LifeCycleStage_names = [...]string{
	"RAW", "COLORS_AND_SIMPLIFIED_JSDOC", "NORMALIZED", "NORMALIZED_OBFUSCATED",
}

// String is Enum.name().
func (s AbstractCompiler_LifeCycleStage) String() string {
	return abstractCompiler_LifeCycleStage_names[s]
}

// Ordinal is Enum.ordinal().
func (s AbstractCompiler_LifeCycleStage) Ordinal() int { return int(s) }

// Java: AbstractCompiler.LifeCycleStage.isNormalized()
func (s AbstractCompiler_LifeCycleStage) IsNormalized() bool {
	return s == AbstractCompiler_LifeCycleStage_NORMALIZED ||
		s == AbstractCompiler_LifeCycleStage_NORMALIZED_OBFUSCATED
}

// Java: AbstractCompiler.LifeCycleStage.isNormalizedUnobfuscated()
func (s AbstractCompiler_LifeCycleStage) IsNormalizedUnobfuscated() bool {
	return s == AbstractCompiler_LifeCycleStage_NORMALIZED
}

// Java: AbstractCompiler.LifeCycleStage.isNormalizedObfuscated()
func (s AbstractCompiler_LifeCycleStage) IsNormalizedObfuscated() bool {
	return s == AbstractCompiler_LifeCycleStage_NORMALIZED_OBFUSCATED
}

// Java: AbstractCompiler.LifeCycleStage.hasColorAndSimplifiedJSDoc()
func (s AbstractCompiler_LifeCycleStage) HasColorAndSimplifiedJSDoc() bool {
	return s != AbstractCompiler_LifeCycleStage_RAW
}

// Java: AbstractCompiler.RUNTIME_LIB_DIR
const AbstractCompiler_RUNTIME_LIB_DIR = RuntimeJsLibManager_RUNTIME_LIB_DIR

// CodeChangeHandler is a simple listener for code change events.
type CodeChangeHandler interface {
	// ReportChange reports a change to some unspecified node of the AST.
	//
	// Java: CodeChangeHandler.reportChange()
	ReportChange()
}

// ForbiddenChange is a change handler that throws an exception if any
// changes are made.
type ForbiddenChange struct{}

// NewForbiddenChange is new ForbiddenChange().
func NewForbiddenChange() *ForbiddenChange { return &ForbiddenChange{} }

// Java: ForbiddenChange.reportChange()
func (f *ForbiddenChange) ReportChange() {
	panic(jutil.NewIllegalStateExceptionMsg("Code changes forbidden"))
}

// RecentChange records whether the code has changed since the last reset.
type RecentChange struct {
	hasChanged bool
}

// NewRecentChange is new RecentChange().
func NewRecentChange() *RecentChange { return &RecentChange{} }

// Java: RecentChange.reportChange()
func (r *RecentChange) ReportChange() {
	r.hasChanged = true
}

// Java: RecentChange.hasCodeChanged()
func (r *RecentChange) HasCodeChanged() bool {
	return r.hasChanged
}

// Java: RecentChange.reset()
func (r *RecentChange) Reset() {
	r.hasChanged = false
}

// RemoveWeakSources removes the contents of weak source files.
type RemoveWeakSources struct {
	compiler AbstractCompiler
}

// Java: RemoveWeakSources(AbstractCompiler)
func NewRemoveWeakSources(compiler AbstractCompiler) *RemoveWeakSources {
	return &RemoveWeakSources{compiler: compiler}
}

// Java: RemoveWeakSources.process(Node, Node)
func (r *RemoveWeakSources) Process(externs rhino.Node, root rhino.Node) {
	for script := root.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
		jutil.CheckState(script.IsScript())
		if script.GetStaticSourceFile().IsWeak() {
			// Keep the file but remove the contents, since some users expect the number of input and
			// output files to be the same.
			NodeUtil_DeleteChildren(script, r.compiler)
		}
	}
}
