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
// Port of com/google/javascript/jscomp/parsing/parser/trees/FormalParameterListTree.java (generated mechanically by j2go).

package parser

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.FormalParameterListTree
type FormalParameterListTree struct {
	ParseTreeBase
	param0           ParseTree        `gson:"param0"`
	param1           ParseTree        `gson:"param1"`
	param2           ParseTree        `gson:"param2"`
	restParameters   []ParseTree      `gson:"restParameters"`
	comma0           SourcePosition   `gson:"comma0"`
	comma1           SourcePosition   `gson:"comma1"`
	comma2           SourcePosition   `gson:"comma2"`
	restCommas       []SourcePosition `gson:"restCommas"`
	parameterCount   int32            `gson:"parameterCount"`
	commaCount       int32            `gson:"commaCount"`
	hasTrailingComma bool             `gson:"hasTrailingComma"`
}

func (this_ *FormalParameterListTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *FormalParameterListTree) IsJavaNil_() bool { return this_ == nil }
func (*FormalParameterListTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.FormalParameterListTree", "FormalParameterListTree"
}
func (*FormalParameterListTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.FormalParameterListTree", "com.google.javascript.jscomp.parsing.parser.trees.FormalParameterListTree", func(v_ any) bool {
		return j2gort.InstanceOf[*FormalParameterListTree](v_)
	})
}
func (*FormalParameterListTree) Implements_parser_ParseTree()          {}
func (this_ *FormalParameterListTree) FieldGet_parameterCount() int32  { return this_.parameterCount }
func (this_ *FormalParameterListTree) FieldGet_hasTrailingComma() bool { return this_.hasTrailingComma }

func (this_ *FormalParameterListTree) initFields_FormalParameterListTree() {
	_ = this_
}

// Java constructor: FormalParameterListTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location)
func NewFormalParameterListTree_1(location SourceRange) *FormalParameterListTree {
	this_ := &FormalParameterListTree{}
	this_.BindSelf_(this_)
	InitFormalParameterListTree_1(this_, location)
	return this_
}
func allocPooled_FormalParameterListTree_1(slab_ *j2gort.PooledSlab[FormalParameterListTree], location SourceRange) *FormalParameterListTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitFormalParameterListTree_1(this_, location)
	return this_
}

func InitFormalParameterListTree_1(this_ *FormalParameterListTree, location SourceRange) {
	_ = location
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_FORMAL_PARAMETER_LIST(), location)
	this_.initFields_FormalParameterListTree()
	this_.parameterCount = int32(0)
	this_.param0 = nil
	this_.param1 = nil
	this_.param2 = nil
	this_.restParameters = nil
	this_.hasTrailingComma = false
	this_.commaCount = int32(0)
	this_.comma0 = (SourcePosition{})
	this_.comma1 = (SourcePosition{})
	this_.comma2 = (SourcePosition{})
	this_.restCommas = nil
}

// Java constructor: FormalParameterListTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree param0, boolean hasTrailingComma, int commaCount, com.google.javascript.jscomp.parsing.parser.util.@org.jspecify.annotations.Nullable SourcePosition comma0)
func NewFormalParameterListTree_5(location SourceRange, param0 ParseTree, hasTrailingComma bool, commaCount int32, comma0 SourcePosition) *FormalParameterListTree {
	this_ := &FormalParameterListTree{}
	this_.BindSelf_(this_)
	InitFormalParameterListTree_5(this_, location, param0, hasTrailingComma, commaCount, comma0)
	return this_
}
func allocPooled_FormalParameterListTree_5(slab_ *j2gort.PooledSlab[FormalParameterListTree], location SourceRange, param0 ParseTree, hasTrailingComma bool, commaCount int32, comma0 SourcePosition) *FormalParameterListTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitFormalParameterListTree_5(this_, location, param0, hasTrailingComma, commaCount, comma0)
	return this_
}

func InitFormalParameterListTree_5(this_ *FormalParameterListTree, location SourceRange, param0 ParseTree, hasTrailingComma bool, commaCount int32, comma0 SourcePosition) {
	_ = location
	_ = param0
	_ = hasTrailingComma
	_ = commaCount
	_ = comma0
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_FORMAL_PARAMETER_LIST(), location)
	this_.initFields_FormalParameterListTree()
	if !((commaCount <= int32(1)) && (commaCount >= int32(0))) {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Unexpected # of comma and formal params.")))
	}
	this_.parameterCount = int32(1)
	this_.param0 = param0
	this_.param1 = nil
	this_.param2 = nil
	this_.restParameters = nil
	this_.hasTrailingComma = hasTrailingComma
	this_.commaCount = commaCount
	this_.comma0 = comma0
	this_.comma1 = (SourcePosition{})
	this_.comma2 = (SourcePosition{})
	this_.restCommas = nil
}

// Java constructor: FormalParameterListTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree param0, com.google.javascript.jscomp.parsing.parser.trees.ParseTree param1, boolean hasTrailingComma, int commaCount, com.google.javascript.jscomp.parsing.parser.util.SourcePosition comma0, com.google.javascript.jscomp.parsing.parser.util.@org.jspecify.annotations.Nullable SourcePosition comma1)
func NewFormalParameterListTree_7(location SourceRange, param0 ParseTree, param1 ParseTree, hasTrailingComma bool, commaCount int32, comma0 SourcePosition, comma1 SourcePosition) *FormalParameterListTree {
	this_ := &FormalParameterListTree{}
	this_.BindSelf_(this_)
	InitFormalParameterListTree_7(this_, location, param0, param1, hasTrailingComma, commaCount, comma0, comma1)
	return this_
}
func allocPooled_FormalParameterListTree_7(slab_ *j2gort.PooledSlab[FormalParameterListTree], location SourceRange, param0 ParseTree, param1 ParseTree, hasTrailingComma bool, commaCount int32, comma0 SourcePosition, comma1 SourcePosition) *FormalParameterListTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitFormalParameterListTree_7(this_, location, param0, param1, hasTrailingComma, commaCount, comma0, comma1)
	return this_
}

func InitFormalParameterListTree_7(this_ *FormalParameterListTree, location SourceRange, param0 ParseTree, param1 ParseTree, hasTrailingComma bool, commaCount int32, comma0 SourcePosition, comma1 SourcePosition) {
	_ = location
	_ = param0
	_ = param1
	_ = hasTrailingComma
	_ = commaCount
	_ = comma0
	_ = comma1
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_FORMAL_PARAMETER_LIST(), location)
	this_.initFields_FormalParameterListTree()
	if !((commaCount <= int32(2)) && (commaCount >= int32(1))) {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Unexpected # of comma and formal params.")))
	}
	this_.parameterCount = int32(2)
	this_.param0 = param0
	this_.param1 = param1
	this_.param2 = nil
	this_.restParameters = nil
	this_.hasTrailingComma = hasTrailingComma
	this_.commaCount = commaCount
	this_.comma0 = comma0
	this_.comma1 = comma1
	this_.comma2 = (SourcePosition{})
	this_.restCommas = nil
}

// Java constructor: FormalParameterListTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, com.google.javascript.jscomp.parsing.parser.trees.ParseTree param0, com.google.javascript.jscomp.parsing.parser.trees.ParseTree param1, com.google.javascript.jscomp.parsing.parser.trees.ParseTree param2, boolean hasTrailingComma, int commaCount, com.google.javascript.jscomp.parsing.parser.util.SourcePosition comma0, com.google.javascript.jscomp.parsing.parser.util.SourcePosition comma1, com.google.javascript.jscomp.parsing.parser.util.@org.jspecify.annotations.Nullable SourcePosition comma2)
func NewFormalParameterListTree_9(location SourceRange, param0 ParseTree, param1 ParseTree, param2 ParseTree, hasTrailingComma bool, commaCount int32, comma0 SourcePosition, comma1 SourcePosition, comma2 SourcePosition) *FormalParameterListTree {
	this_ := &FormalParameterListTree{}
	this_.BindSelf_(this_)
	InitFormalParameterListTree_9(this_, location, param0, param1, param2, hasTrailingComma, commaCount, comma0, comma1, comma2)
	return this_
}
func allocPooled_FormalParameterListTree_9(slab_ *j2gort.PooledSlab[FormalParameterListTree], location SourceRange, param0 ParseTree, param1 ParseTree, param2 ParseTree, hasTrailingComma bool, commaCount int32, comma0 SourcePosition, comma1 SourcePosition, comma2 SourcePosition) *FormalParameterListTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitFormalParameterListTree_9(this_, location, param0, param1, param2, hasTrailingComma, commaCount, comma0, comma1, comma2)
	return this_
}

func InitFormalParameterListTree_9(this_ *FormalParameterListTree, location SourceRange, param0 ParseTree, param1 ParseTree, param2 ParseTree, hasTrailingComma bool, commaCount int32, comma0 SourcePosition, comma1 SourcePosition, comma2 SourcePosition) {
	_ = location
	_ = param0
	_ = param1
	_ = param2
	_ = hasTrailingComma
	_ = commaCount
	_ = comma0
	_ = comma1
	_ = comma2
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_FORMAL_PARAMETER_LIST(), location)
	this_.initFields_FormalParameterListTree()
	if !((commaCount <= int32(3)) && (commaCount >= int32(2))) {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Unexpected # of comma and formal params.")))
	}
	this_.parameterCount = int32(3)
	this_.param0 = param0
	this_.param1 = param1
	this_.param2 = param2
	this_.restParameters = nil
	this_.hasTrailingComma = hasTrailingComma
	this_.commaCount = commaCount
	this_.comma0 = comma0
	this_.comma1 = comma1
	this_.comma2 = comma2
	this_.restCommas = nil
}

// Java constructor: FormalParameterListTree(com.google.javascript.jscomp.parsing.parser.util.SourceRange location, int parameterCount, com.google.javascript.jscomp.parsing.parser.trees.ParseTree param0, com.google.javascript.jscomp.parsing.parser.trees.ParseTree param1, com.google.javascript.jscomp.parsing.parser.trees.ParseTree param2, com.google.javascript.jscomp.parsing.parser.trees.ParseTree @org.jspecify.annotations.Nullable [] restParameters, boolean hasTrailingComma, int commaCount, com.google.javascript.jscomp.parsing.parser.util.SourcePosition comma0, com.google.javascript.jscomp.parsing.parser.util.SourcePosition comma1, com.google.javascript.jscomp.parsing.parser.util.SourcePosition comma2, com.google.javascript.jscomp.parsing.parser.util.SourcePosition @org.jspecify.annotations.Nullable [] restCommas)
func NewFormalParameterListTree_12(location SourceRange, parameterCount int32, param0 ParseTree, param1 ParseTree, param2 ParseTree, restParameters []ParseTree, hasTrailingComma bool, commaCount int32, comma0 SourcePosition, comma1 SourcePosition, comma2 SourcePosition, restCommas []SourcePosition) *FormalParameterListTree {
	this_ := &FormalParameterListTree{}
	this_.BindSelf_(this_)
	InitFormalParameterListTree_12(this_, location, parameterCount, param0, param1, param2, restParameters, hasTrailingComma, commaCount, comma0, comma1, comma2, restCommas)
	return this_
}
func allocPooled_FormalParameterListTree_12(slab_ *j2gort.PooledSlab[FormalParameterListTree], location SourceRange, parameterCount int32, param0 ParseTree, param1 ParseTree, param2 ParseTree, restParameters []ParseTree, hasTrailingComma bool, commaCount int32, comma0 SourcePosition, comma1 SourcePosition, comma2 SourcePosition, restCommas []SourcePosition) *FormalParameterListTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitFormalParameterListTree_12(this_, location, parameterCount, param0, param1, param2, restParameters, hasTrailingComma, commaCount, comma0, comma1, comma2, restCommas)
	return this_
}

func InitFormalParameterListTree_12(this_ *FormalParameterListTree, location SourceRange, parameterCount int32, param0 ParseTree, param1 ParseTree, param2 ParseTree, restParameters []ParseTree, hasTrailingComma bool, commaCount int32, comma0 SourcePosition, comma1 SourcePosition, comma2 SourcePosition, restCommas []SourcePosition) {
	_ = location
	_ = parameterCount
	_ = param0
	_ = param1
	_ = param2
	_ = restParameters
	_ = hasTrailingComma
	_ = commaCount
	_ = comma0
	_ = comma1
	_ = comma2
	_ = restCommas
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_FORMAL_PARAMETER_LIST(), location)
	this_.initFields_FormalParameterListTree()
	if !((commaCount <= parameterCount) && (commaCount >= (parameterCount - int32(1)))) {
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "Unexpected # of comma and formal params.")))
	}
	this_.parameterCount = parameterCount
	this_.param0 = param0
	this_.param1 = param1
	this_.param2 = param2
	this_.restParameters = restParameters
	this_.hasTrailingComma = hasTrailingComma
	this_.commaCount = commaCount
	this_.comma0 = comma0
	this_.comma1 = comma1
	this_.comma2 = comma2
	this_.restCommas = restCommas
}

// Java method: FormalParameterListTree.getParameter(int index)
func (this_ *FormalParameterListTree) GetParameter(index int32) ParseTree {
	_ = index
	switch index {
	case int32(0):
		return this_.param0
	case int32(1):
		return this_.param1
	case int32(2):
		return this_.param2
	default:
		return this_.restParameters[(index - int32(3))]
	}
	panic("unreachable switch expression")
}

// Java method: FormalParameterListTree.getCommaPosition(int index)
func (this_ *FormalParameterListTree) GetCommaPosition(index int32) SourcePosition {
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
