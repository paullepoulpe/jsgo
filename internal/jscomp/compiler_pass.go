// Port of com/google/javascript/jscomp/CompilerPass.java @ closure-compiler v20261021

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CompilerPass is the interface for classes that can compile JS.
//
// Class has single function "process", which is passed the root node of the
// parsed JS tree, as well as the root node of the external JS tree (used to
// provide a public API and prevent renaming of system functions).
//
// Use this class to support testing with BaseCompilerTest
type CompilerPass interface {
	// Process the JS with root node root. Can modify the contents of each
	// Node tree.
	//
	// externs is the top of external JS tree; root is the top of JS tree.
	//
	// Java: CompilerPass.process(Node, Node)
	Process(externs *rhino.Node, root *rhino.Node)
}

// CompilerPassFunc lets a Go closure stand for a Java lambda or method
// reference implementing CompilerPass, e.g. `(externs, root) -> ...`.
//
// Go-only: language-adapter — functional interface implemented by a lambda.
type CompilerPassFunc func(externs *rhino.Node, root *rhino.Node)

// Process calls f(externs, root).
func (f CompilerPassFunc) Process(externs *rhino.Node, root *rhino.Node) { f(externs, root) }
