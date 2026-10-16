// Port of com/google/javascript/jscomp/RemoveWeakSources.java @ closure-compiler v20260415

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// RemoveWeakSources removes the contents of weak source files.
type RemoveWeakSources struct {
	compiler AbstractCompiler
}

// Java: RemoveWeakSources(AbstractCompiler)
func NewRemoveWeakSources(compiler AbstractCompiler) *RemoveWeakSources {
	return &RemoveWeakSources{compiler: compiler}
}

// Java: RemoveWeakSources.process(Node, Node)
func (r *RemoveWeakSources) Process(externs *rhino.Node, root *rhino.Node) {
	for script := root.GetFirstChild(); script != nil; script = script.GetNext() {
		jutil.CheckState(script.IsScript())
		if script.GetStaticSourceFile().IsWeak() {
			// Keep the file but remove the contents, since some users expect the number of input and
			// output files to be the same.
			NodeUtil_DeleteChildren(script, r.compiler)
		}
	}
}
