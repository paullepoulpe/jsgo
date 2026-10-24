// Port of com/google/javascript/jscomp/parsing/parser/trees/ProgramTree.java @ closure-compiler v20261021

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

// ProgramTree is the parse tree for ParseTreeType PROGRAM.
type ProgramTree struct {
	ParseTreeBase

	SourceElements []ParseTree
	SourceComments []*Comment
}

// NewProgramTree constructs a ProgramTree.
//
// Java: ProgramTree(SourceRange, ImmutableList, ImmutableList)
func NewProgramTree(location SourceRange, sourceElements []ParseTree, sourceComments []*Comment) *ProgramTree {
	t := &ProgramTree{}
	t.InitParseTreeBase(ParseTreeType_PROGRAM, location)
	t.SourceElements = sourceElements
	t.SourceComments = sourceComments
	return t
}
