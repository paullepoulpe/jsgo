/*
 * Copyright 2026 The Closure Compiler Authors.
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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestNodeUtil_MalformedGoogModuleAndRequireCalls(t *testing.T) {
	if NodeUtil_IsGoogModuleCall(rhino.NilNode) {
		t.Fatal("NodeUtil_IsGoogModuleCall(NilNode) = true, want false")
	}

	compiler := NewCompilerWithErrorManager(NewBlackHoleErrorManager())
	compiler.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", "")},
		[]*SourceFile{
			SourceFile_FromCode("m0.js", "goog.module(); const a = goog.require(); const {b} = goog.require();"),
			SourceFile_FromCode("m1.js", "goog.module(123); const c = goog.require(123); const {d} = goog.require(123);"),
		},
		NewCompilerOptions())
	compiler.ParseInputs()

	jsRoot := compiler.GetJsRoot()
	for script := jsRoot.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
		container := script
		if container.GetFirstChild().IsModuleBody() {
			container = container.GetFirstChild()
		}
		googModStmt := container.GetFirstChild()
		if !NodeUtil_IsGoogModuleCall(googModStmt) {
			t.Fatalf("NodeUtil_IsGoogModuleCall(%v) = false, want true", googModStmt)
		}
		constA := googModStmt.GetNext()
		nameA := constA.GetFirstChild()
		if got := NodeUtil_GetGoogRequireInfoNode(nameA); got != nil {
			t.Fatalf("NodeUtil_GetGoogRequireInfoNode(%v) = %v, want nil", nameA, got)
		}
		constDestruct := constA.GetNext()
		nameB := constDestruct.GetFirstFirstChild().GetFirstChild().GetFirstChild()
		if got := NodeUtil_GetGoogRequireInfoNode(nameB); got != nil {
			t.Fatalf("NodeUtil_GetGoogRequireInfoNode(%v) = %v, want nil", nameB, got)
		}
	}

	// Verify ClosureRewriteModule still emits INVALID_MODULE_ID_ARG on invalid goog.module() calls.
	errCompiler := NewCompiler()
	opts := NewCompilerOptions()
	opts.SetContinueAfterErrors(true)
	errCompiler.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", "")},
		[]*SourceFile{SourceFile_FromCode("bad.js", "goog.module();")},
		opts)
	errCompiler.ParseInputs()
	NewClosureRewriteModule(errCompiler, nil, nil).Process(errCompiler.GetExternsRoot(), errCompiler.GetJsRoot())
	if len(errCompiler.GetErrors()) == 0 || errCompiler.GetErrors()[0].GetType() != ClosureRewriteModule_INVALID_MODULE_ID_ARG {
		t.Fatalf("ClosureRewriteModule errors = %v, want %v", errCompiler.GetErrors(), ClosureRewriteModule_INVALID_MODULE_ID_ARG)
	}
}
