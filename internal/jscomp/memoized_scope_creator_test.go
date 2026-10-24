// Port of test/com/google/javascript/jscomp/MemoizedScopeCreatorTest.java @ closure-compiler v20261021

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

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestMemoizedScopeCreator_Memoization(t *testing.T) {
	root1 := rhino.NewNode(rhino.Token_ROOT)
	root2 := rhino.NewNode(rhino.Token_ROOT)
	compiler := newTestCompiler(NewCompilerOptions())
	var creator ScopeCreator = NewMemoizedScopeCreator(NewSyntacticScopeCreator(compiler))
	scopeA := creator.CreateScope(root1, nil).(*Scope)
	if got := creator.CreateScope(root1, nil); got != AbstractScope(scopeA) {
		t.Errorf("createScope(root1) = %v, want the same instance %v", got, scopeA)
	}
	if got := creator.CreateScope(root2, nil); got == AbstractScope(scopeA) {
		t.Errorf("createScope(root2) is the instance of root1")
	}
}

func TestMemoizedScopeCreator_PreconditionCheck(t *testing.T) {
	compiler := newTestCompiler(NewCompilerOptions())
	root := rhino.NewNode(rhino.Token_ROOT)
	var creator ScopeCreator = NewMemoizedScopeCreator(NewSyntacticScopeCreator(compiler))
	scopeA := creator.CreateScope(root, nil).(*Scope)

	handled := false
	jutil.Try(func() {
		creator.CreateScope(root, scopeA)
	}, func(e jutil.IIllegalStateException) {
		handled = true
	})
	if !handled {
		t.Error("expected an IllegalStateException")
	}
}
