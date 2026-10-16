// Port of com/google/javascript/jscomp/J2clSourceFileChecker.java @ closure-compiler v20260415

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

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// J2clSourceFileChecker checks whether there are J2CL generated source files
// with pattern "*.java.js".
//
// Java: final class J2clSourceFileChecker implements CompilerPass
type J2clSourceFileChecker struct {
	compiler AbstractCompiler
}

// Java: J2clSourceFileChecker(AbstractCompiler)
func NewJ2clSourceFileChecker(compiler AbstractCompiler) *J2clSourceFileChecker {
	return &J2clSourceFileChecker{compiler: compiler}
}

// Java: J2clSourceFileChecker.hasJ2cl(Node)
func j2clSourceFileChecker_hasJ2cl(root *rhino.Node) bool {
	for script := root.GetFirstChild(); script != nil; script = script.GetNext() {
		jutil.CheckState(script.IsScript())
		if J2clSourceUtils_IsJ2clSourceOfNode(script) {
			return true
		}
	}
	return false
}

// Java: J2clSourceFileChecker.process(Node, Node)
func (c *J2clSourceFileChecker) Process(externs *rhino.Node, root *rhino.Node) {
	if j2clSourceFileChecker_hasJ2cl(root) {
		J2clSourceFileChecker_MarkToRunJ2clPasses(c.compiler)
	}
}

// Java: J2clSourceFileChecker.markToRunJ2clPasses(AbstractCompiler)
func J2clSourceFileChecker_MarkToRunJ2clPasses(compiler AbstractCompiler) {
	compiler.SetRunJ2clPasses(true)
}

// J2clSourceFileChecker_ShouldRunJ2clPasses indicates whether it should run
// future J2CL passes with information from the compiler. For example, if the
// compiler's HAS_J2CL annotation is false, it should.
//
// Java: J2clSourceFileChecker.shouldRunJ2clPasses(AbstractCompiler)
func J2clSourceFileChecker_ShouldRunJ2clPasses(compiler AbstractCompiler) bool {
	return compiler.RunJ2clPasses()
}
