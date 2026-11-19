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
// Port of com/google/javascript/jscomp/parsing/parser/trees/ArgumentListTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.ArgumentListTree
type ArgumentListTree struct {
	ParseTreeBase
	arg0             ParseTree        `gson:"arg0"`
	arg1             ParseTree        `gson:"arg1"`
	arg2             ParseTree        `gson:"arg2"`
	restArguments    []ParseTree      `gson:"restArguments"`
	comma0           SourcePosition   `gson:"comma0"`
	comma1           SourcePosition   `gson:"comma1"`
	comma2           SourcePosition   `gson:"comma2"`
	restCommas       []SourcePosition `gson:"restCommas"`
	argumentCount    int32            `gson:"argumentCount"`
	commaCount       int32            `gson:"commaCount"`
	hasTrailingComma bool             `gson:"hasTrailingComma"`
}

func (this_ *ArgumentListTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *ArgumentListTree) IsJavaNil_() bool { return this_ == nil }
func (*ArgumentListTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.ArgumentListTree", "ArgumentListTree"
}
func (*ArgumentListTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.ArgumentListTree", "com.google.javascript.jscomp.parsing.parser.trees.ArgumentListTree", func(v_ any) bool {
		return j2gort.InstanceOf[*ArgumentListTree](v_)
	})
}
func (*ArgumentListTree) Implements_parser_ParseTree()          {}
func (this_ *ArgumentListTree) FieldGet_argumentCount() int32   { return this_.argumentCount }
func (this_ *ArgumentListTree) FieldGet_arg0() ParseTree        { return this_.arg0 }
func (this_ *ArgumentListTree) FieldGet_arg1() ParseTree        { return this_.arg1 }
func (this_ *ArgumentListTree) FieldGet_hasTrailingComma() bool { return this_.hasTrailingComma }

func (this_ *ArgumentListTree) initFields_ArgumentListTree() {
	_ = this_
}

// Java constructor: ArgumentListTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location)
func NewArgumentListTree_1(location SourceRange) *ArgumentListTree {
	this_ := &ArgumentListTree{}
	this_.BindSelf_(this_)
	InitArgumentListTree_1(this_, location)
	return this_
}
func allocPooled_ArgumentListTree_1(slab_ *j2gort.PooledSlab[ArgumentListTree], location SourceRange) *ArgumentListTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitArgumentListTree_1(this_, location)
	return this_
}

func InitArgumentListTree_1(this_ *ArgumentListTree, location SourceRange) {
	_ = location
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_ARGUMENT_LIST(), location)
	this_.initFields_ArgumentListTree()
	this_.argumentCount = int32(0)
	this_.arg0 = nil
	this_.arg1 = nil
	this_.arg2 = nil
	this_.restArguments = nil
	this_.hasTrailingComma = false
	this_.commaCount = int32(0)
	this_.comma0 = (SourcePosition{})
	this_.comma1 = (SourcePosition{})
	this_.comma2 = (SourcePosition{})
	this_.restCommas = nil
}

// Java constructor: ArgumentListTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree arg0, boolean hasTrailingComma, int commaCount, com.google.javascript.jscomp.parsing.parser.util.@org.jspecify.annotations.Nullable SourcePosition comma0)
func NewArgumentListTree_5(location SourceRange, arg0 ParseTree, hasTrailingComma bool, commaCount int32, comma0 SourcePosition) *ArgumentListTree {
	this_ := &ArgumentListTree{}
	this_.BindSelf_(this_)
	InitArgumentListTree_5(this_, location, arg0, hasTrailingComma, commaCount, comma0)
	return this_
}
func allocPooled_ArgumentListTree_5(slab_ *j2gort.PooledSlab[ArgumentListTree], location SourceRange, arg0 ParseTree, hasTrailingComma bool, commaCount int32, comma0 SourcePosition) *ArgumentListTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitArgumentListTree_5(this_, location, arg0, hasTrailingComma, commaCount, comma0)
	return this_
}

func InitArgumentListTree_5(this_ *ArgumentListTree, location SourceRange, arg0 ParseTree, hasTrailingComma bool, commaCount int32, comma0 SourcePosition) {
	_ = location
	_ = arg0
	_ = hasTrailingComma
	_ = commaCount
	_ = comma0
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_ARGUMENT_LIST(), location)
	this_.initFields_ArgumentListTree()
	this_.argumentCount = int32(1)
	this_.arg0 = arg0
	this_.arg1 = nil
	this_.arg2 = nil
	this_.restArguments = nil
	this_.hasTrailingComma = hasTrailingComma
	this_.commaCount = commaCount
	this_.comma0 = comma0
	this_.comma1 = (SourcePosition{})
	this_.comma2 = (SourcePosition{})
	this_.restCommas = nil
}

// Java constructor: ArgumentListTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree arg0, com.google.javascript.jscomp.parsing.parser.trees.ParseTree arg1, boolean hasTrailingComma, int commaCount, com.google.javascript.jscomp.parsing.parser.util.@org.jspecify.annotations.Nullable SourcePosition comma0, com.google.javascript.jscomp.parsing.parser.util.@org.jspecify.annotations.Nullable SourcePosition comma1)
func NewArgumentListTree_7(location SourceRange, arg0 ParseTree, arg1 ParseTree, hasTrailingComma bool, commaCount int32, comma0 SourcePosition, comma1 SourcePosition) *ArgumentListTree {
	this_ := &ArgumentListTree{}
	this_.BindSelf_(this_)
	InitArgumentListTree_7(this_, location, arg0, arg1, hasTrailingComma, commaCount, comma0, comma1)
	return this_
}
func allocPooled_ArgumentListTree_7(slab_ *j2gort.PooledSlab[ArgumentListTree], location SourceRange, arg0 ParseTree, arg1 ParseTree, hasTrailingComma bool, commaCount int32, comma0 SourcePosition, comma1 SourcePosition) *ArgumentListTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitArgumentListTree_7(this_, location, arg0, arg1, hasTrailingComma, commaCount, comma0, comma1)
	return this_
}

func InitArgumentListTree_7(this_ *ArgumentListTree, location SourceRange, arg0 ParseTree, arg1 ParseTree, hasTrailingComma bool, commaCount int32, comma0 SourcePosition, comma1 SourcePosition) {
	_ = location
	_ = arg0
	_ = arg1
	_ = hasTrailingComma
	_ = commaCount
	_ = comma0
	_ = comma1
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_ARGUMENT_LIST(), location)
	this_.initFields_ArgumentListTree()
	this_.argumentCount = int32(2)
	this_.arg0 = arg0
	this_.arg1 = arg1
	this_.arg2 = nil
	this_.restArguments = nil
	this_.hasTrailingComma = hasTrailingComma
	this_.commaCount = commaCount
	this_.comma0 = comma0
	this_.comma1 = comma1
	this_.comma2 = (SourcePosition{})
	this_.restCommas = nil
}

// Java constructor: ArgumentListTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree arg0, com.google.javascript.jscomp.parsing.parser.trees.ParseTree arg1, com.google.javascript.jscomp.parsing.parser.trees.ParseTree arg2, boolean hasTrailingComma, int commaCount, com.google.javascript.jscomp.parsing.parser.util.@org.jspecify.annotations.Nullable SourcePosition comma0, com.google.javascript.jscomp.parsing.parser.util.@org.jspecify.annotations.Nullable SourcePosition comma1, com.google.javascript.jscomp.parsing.parser.util.@org.jspecify.annotations.Nullable SourcePosition comma2)
func NewArgumentListTree_9(location SourceRange, arg0 ParseTree, arg1 ParseTree, arg2 ParseTree, hasTrailingComma bool, commaCount int32, comma0 SourcePosition, comma1 SourcePosition, comma2 SourcePosition) *ArgumentListTree {
	this_ := &ArgumentListTree{}
	this_.BindSelf_(this_)
	InitArgumentListTree_9(this_, location, arg0, arg1, arg2, hasTrailingComma, commaCount, comma0, comma1, comma2)
	return this_
}
func allocPooled_ArgumentListTree_9(slab_ *j2gort.PooledSlab[ArgumentListTree], location SourceRange, arg0 ParseTree, arg1 ParseTree, arg2 ParseTree, hasTrailingComma bool, commaCount int32, comma0 SourcePosition, comma1 SourcePosition, comma2 SourcePosition) *ArgumentListTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitArgumentListTree_9(this_, location, arg0, arg1, arg2, hasTrailingComma, commaCount, comma0, comma1, comma2)
	return this_
}

func InitArgumentListTree_9(this_ *ArgumentListTree, location SourceRange, arg0 ParseTree, arg1 ParseTree, arg2 ParseTree, hasTrailingComma bool, commaCount int32, comma0 SourcePosition, comma1 SourcePosition, comma2 SourcePosition) {
	_ = location
	_ = arg0
	_ = arg1
	_ = arg2
	_ = hasTrailingComma
	_ = commaCount
	_ = comma0
	_ = comma1
	_ = comma2
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_ARGUMENT_LIST(), location)
	this_.initFields_ArgumentListTree()
	this_.argumentCount = int32(3)
	this_.arg0 = arg0
	this_.arg1 = arg1
	this_.arg2 = arg2
	this_.restArguments = nil
	this_.hasTrailingComma = hasTrailingComma
	this_.commaCount = commaCount
	this_.comma0 = comma0
	this_.comma1 = comma1
	this_.comma2 = comma2
	this_.restCommas = nil
}

// Java constructor: ArgumentListTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, int argumentCount, com.google.javascript.jscomp.parsing.parser.trees.ParseTree arg0, com.google.javascript.jscomp.parsing.parser.trees.ParseTree arg1, com.google.javascript.jscomp.parsing.parser.trees.ParseTree arg2, com.google.javascript.jscomp.parsing.parser.trees.ParseTree @org.jspecify.annotations.Nullable [] restArguments, boolean hasTrailingComma, int commaCount, com.google.javascript.jscomp.parsing.parser.util.@org.jspecify.annotations.Nullable SourcePosition comma0, com.google.javascript.jscomp.parsing.parser.util.@org.jspecify.annotations.Nullable SourcePosition comma1, com.google.javascript.jscomp.parsing.parser.util.@org.jspecify.annotations.Nullable SourcePosition comma2, com.google.javascript.jscomp.parsing.parser.util.SourcePosition @org.jspecify.annotations.Nullable [] restCommas)
func NewArgumentListTree_12(location SourceRange, argumentCount int32, arg0 ParseTree, arg1 ParseTree, arg2 ParseTree, restArguments []ParseTree, hasTrailingComma bool, commaCount int32, comma0 SourcePosition, comma1 SourcePosition, comma2 SourcePosition, restCommas []SourcePosition) *ArgumentListTree {
	this_ := &ArgumentListTree{}
	this_.BindSelf_(this_)
	InitArgumentListTree_12(this_, location, argumentCount, arg0, arg1, arg2, restArguments, hasTrailingComma, commaCount, comma0, comma1, comma2, restCommas)
	return this_
}
func allocPooled_ArgumentListTree_12(slab_ *j2gort.PooledSlab[ArgumentListTree], location SourceRange, argumentCount int32, arg0 ParseTree, arg1 ParseTree, arg2 ParseTree, restArguments []ParseTree, hasTrailingComma bool, commaCount int32, comma0 SourcePosition, comma1 SourcePosition, comma2 SourcePosition, restCommas []SourcePosition) *ArgumentListTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitArgumentListTree_12(this_, location, argumentCount, arg0, arg1, arg2, restArguments, hasTrailingComma, commaCount, comma0, comma1, comma2, restCommas)
	return this_
}

func InitArgumentListTree_12(this_ *ArgumentListTree, location SourceRange, argumentCount int32, arg0 ParseTree, arg1 ParseTree, arg2 ParseTree, restArguments []ParseTree, hasTrailingComma bool, commaCount int32, comma0 SourcePosition, comma1 SourcePosition, comma2 SourcePosition, restCommas []SourcePosition) {
	_ = location
	_ = argumentCount
	_ = arg0
	_ = arg1
	_ = arg2
	_ = restArguments
	_ = hasTrailingComma
	_ = commaCount
	_ = comma0
	_ = comma1
	_ = comma2
	_ = restCommas
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_ARGUMENT_LIST(), location)
	this_.initFields_ArgumentListTree()
	this_.argumentCount = argumentCount
	this_.arg0 = arg0
	this_.arg1 = arg1
	this_.arg2 = arg2
	this_.restArguments = restArguments
	this_.hasTrailingComma = hasTrailingComma
	this_.commaCount = commaCount
	this_.comma0 = comma0
	this_.comma1 = comma1
	this_.comma2 = comma2
	this_.restCommas = restCommas
}

// Java method: ArgumentListTree.getArgument(int index)
func (this_ *ArgumentListTree) GetArgument(index int32) ParseTree {
	_ = index
	switch index {
	case int32(0):
		return this_.arg0
	case int32(1):
		return this_.arg1
	case int32(2):
		return this_.arg2
	default:
		return this_.restArguments[(index - int32(3))]
	}
	panic("unreachable switch expression")
}

// Java method: ArgumentListTree.getCommaPosition(int index)
func (this_ *ArgumentListTree) GetCommaPosition(index int32) SourcePosition {
	_ = index
	if (index < int32(0)) || (index >= this_.commaCount) {
		panic(j2gort.AsIface[javalang.IArrayIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.ArrayIndexOutOfBoundsException", j2gort.ValueOfInt(index))))
	}
	switch index {
	case int32(0):
		return this_.comma0
	case int32(1):
		return this_.comma1
	case int32(2):
		return this_.comma2
	default:
		return this_.restCommas[(index - int32(3))]
	}
	panic("unreachable switch expression")
}
