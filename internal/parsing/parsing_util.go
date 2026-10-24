// Port of com/google/javascript/jscomp/parsing/ParsingUtil.java @ closure-compiler v20261021

/*
 * Copyright 2019 The Closure Compiler Authors.
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

package parsing

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ParsingUtil holds utility functions for parsing that don't depend on the
// compiler.
//
// Java: final class ParsingUtil
type ParsingUtil struct{}

// Utility class; do not instantiate.
//
// Java: ParsingUtil()
func NewParsingUtil() *ParsingUtil { return &ParsingUtil{} }

// ParsingUtil_GetParamOrPatternNames calls cb with all NAMEs declared in a
// PARAM_LIST or destructuring pattern.
//
// Java: ParsingUtil.getParamOrPatternNames(Node, Consumer<Node>)
func ParsingUtil_GetParamOrPatternNames(n *rhino.Node, cb func(*rhino.Node)) {
	jutil.CheckNotNull(n)
	switch n.GetToken() {
	case rhino.Token_EMPTY, rhino.Token_GETELEM, rhino.Token_GETPROP:
		// let [,] = ...
		// [someArray[1]] = ...
		// [someObj.someProp] = ...
		// An empty node or assigned to an existing value means nothings here
		// and we can't recurse further.
	case rhino.Token_NAME:
		// let [someName] = ...
		// A name node generally means we found an assignment happening.
		cb(n)
	case rhino.Token_ITER_REST, rhino.Token_OBJECT_REST, rhino.Token_STRING_KEY, rhino.Token_DEFAULT_VALUE:
		// let [...name] = ...
		// let {...name} = ...
		// let {name} = ...
		// let {name = 10} = ...
		// For these node types there is a name node as the first child.
		ParsingUtil_GetParamOrPatternNames(n.GetFirstChild(), cb)
	case rhino.Token_COMPUTED_PROP:
		// let {[someProp]: name} = ...
		// For this node types there is a name node as the second child.
		ParsingUtil_GetParamOrPatternNames(n.GetSecondChild(), cb)
	case rhino.Token_OBJECT_PATTERN, rhino.Token_ARRAY_PATTERN, rhino.Token_PARAM_LIST:
		// let {<pattern>} = ...
		// let [<pattern>] = ...
		// function fn(<params>) ...
		// Each item in the destructuring pattern is scanned for assignments.
		// Although PARAM_LIST isn't technically a pattern, we know how to
		// handle it anyway.
		for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
			ParsingUtil_GetParamOrPatternNames(c, cb)
		}
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected parameter structure"))
	}
}
