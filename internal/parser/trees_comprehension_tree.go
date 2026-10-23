// Port of com/google/javascript/jscomp/parsing/parser/trees/ComprehensionTree.java @ closure-compiler v20261015

/*
 * Copyright 2014 The Closure Compiler Authors.
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

// ComprehensionTree is the parse tree for ParseTreeType COMPREHENSION.
type ComprehensionTree struct {
	ParseTreeBase

	// Hides ParseTreeBase.Type, exactly like the Java field hides ParseTree.type.
	Type ComprehensionTree_ComprehensionType
	// The "for of" and "if" trees.
	Children []ParseTree
	// The AssignmentExpression at the end of the ComprehensionTail.
	TailExpression ParseTree
}

// ComprehensionTree_ComprehensionType is ComprehensionTree.ComprehensionType.
type ComprehensionTree_ComprehensionType int32

const (
	ComprehensionTree_ComprehensionType_ARRAY ComprehensionTree_ComprehensionType = iota
	ComprehensionTree_ComprehensionType_GENERATOR
)

var comprehensionTree_ComprehensionType_names = [...]string{"ARRAY", "GENERATOR"}

// ComprehensionTree_ComprehensionType_values is ComprehensionType.values().
//
// Java: ComprehensionTree.ComprehensionType.values()
func ComprehensionTree_ComprehensionType_values() []ComprehensionTree_ComprehensionType {
	return []ComprehensionTree_ComprehensionType{ComprehensionTree_ComprehensionType_ARRAY, ComprehensionTree_ComprehensionType_GENERATOR}
}

// Name is Enum.name().
//
// Java: ComprehensionTree.ComprehensionType.name()
func (t ComprehensionTree_ComprehensionType) Name() string {
	return comprehensionTree_ComprehensionType_names[t]
}

// String is Enum.toString().
//
// Java: ComprehensionTree.ComprehensionType.toString()
func (t ComprehensionTree_ComprehensionType) String() string { return t.Name() }

// NewComprehensionTree constructs a ComprehensionTree.
//
// Java: ComprehensionTree(SourceRange, ComprehensionType, ImmutableList, ParseTree)
func NewComprehensionTree(location SourceRange, type_ ComprehensionTree_ComprehensionType, children []ParseTree, tailExpression ParseTree) *ComprehensionTree {
	t := &ComprehensionTree{}
	t.InitParseTreeBase(ParseTreeType_COMPREHENSION, location)
	t.Type = type_
	t.Children = children
	t.TailExpression = tailExpression
	return t
}
