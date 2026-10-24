// Port of com/google/javascript/jscomp/J2clSourceUtils.java @ closure-compiler v20261021

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// J2clSourceUtils holds static helpers for recognizing J2CL sources.
type J2clSourceUtils struct{}

// J2clSourceUtils_IsJ2clSource: sourceFileName is a nullable Java String.
//
// Java: J2clSourceUtils.isJ2clSource(String)
func J2clSourceUtils_IsJ2clSource(sourceFileName *string) bool {
	return sourceFileName != nil && strings.HasSuffix(*sourceFileName, ".java.js")
}

// Java: J2clSourceUtils.isJ2clSource(Node)
func J2clSourceUtils_IsJ2clSourceOfNode(scriptNode *rhino.Node) bool {
	jutil.CheckArgument(scriptNode.IsScript())
	return J2clSourceUtils_IsJ2clSource(scriptNode.GetSourceFileName())
}

// Java: J2clSourceUtils.isInsideJ2clSource(Node)
func J2clSourceUtils_IsInsideJ2clSource(node *rhino.Node) bool {
	var scriptNode *rhino.Node
	if node.IsScript() {
		scriptNode = node
	} else {
		scriptNode = NodeUtil_GetEnclosingScript(node)
	}
	if scriptNode == nil {
		return false
	}
	return J2clSourceUtils_IsJ2clSourceOfNode(scriptNode)
}

// Java: J2clSourceUtils()
func newJ2clSourceUtils() *J2clSourceUtils { return &J2clSourceUtils{} }
