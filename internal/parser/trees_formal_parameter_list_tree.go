// Port of com/google/javascript/jscomp/parsing/parser/trees/FormalParameterListTree.java @ closure-compiler v20260415

/*
 * Copyright 2011 The Closure Compiler Authors.
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

package parser

import "github.com/paullepoulpe/jsgo/internal/jutil"

// FormalParameterListTree is the parse tree for ParseTreeType FORMAL_PARAMETER_LIST.
type FormalParameterListTree struct {
	ParseTreeBase

	Parameters       []ParseTree
	HasTrailingComma bool
	CommaPositions   []*SourcePosition
}

// NewFormalParameterListTree constructs a FormalParameterListTree.
//
// Java: FormalParameterListTree(SourceRange, ImmutableList, boolean, ImmutableList)
func NewFormalParameterListTree(location *SourceRange, parameters []ParseTree, hasTrailingComma bool, commaPositions []*SourcePosition) *FormalParameterListTree {
	t := &FormalParameterListTree{}
	t.InitParseTreeBase(t, ParseTreeType_FORMAL_PARAMETER_LIST, location)
	numParams := len(parameters)
	numCommas := len(commaPositions)
	if !(numCommas <= numParams && numCommas >= numParams-1) {
		// Java: checkArgument(cond, "...", parameters, commaPositions); the lists
		// are rendered with AbstractCollection.toString only on failure.
		panic(jutil.NewIllegalArgumentExceptionMsg(jutil.LenientFormat(
			"Unexpected # of comma and formal params.\nparams: %s\ncomma positions: %s",
			jutil.CollectionToString(parameters),
			jutil.CollectionToString(commaPositions))))
	}
	t.Parameters = parameters
	t.HasTrailingComma = hasTrailingComma
	t.CommaPositions = commaPositions
	return t
}
