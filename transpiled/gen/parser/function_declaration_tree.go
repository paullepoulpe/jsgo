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
// Port of com/google/javascript/jscomp/parsing/parser/trees/FunctionDeclarationTree.java (generated mechanically by j2go).

package parser

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.parsing.parser.trees.FunctionDeclarationTree
type FunctionDeclarationTree struct {
	ParseTreeBase
	name                Token                        `gson:"name"`
	formalParameterList *FormalParameterListTree     `gson:"formalParameterList"`
	functionBody        ParseTree                    `gson:"functionBody"`
	isClassMember       bool                         `gson:"isClassMember"`
	isStatic            bool                         `gson:"isStatic"`
	isGenerator         bool                         `gson:"isGenerator"`
	isOptional          bool                         `gson:"isOptional"`
	isAsync             bool                         `gson:"isAsync"`
	kind                FunctionDeclarationTree_Kind `gson:"kind"`
}

func (this_ *FunctionDeclarationTree) BindSelf_(s any) {
	_ = s
	this_.ParseTreeBase.BindSelf_(s)
}
func (this_ *FunctionDeclarationTree) IsJavaNil_() bool { return this_ == nil }
func (*FunctionDeclarationTree) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.FunctionDeclarationTree", "FunctionDeclarationTree"
}
func (*FunctionDeclarationTree) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.FunctionDeclarationTree", "com.google.javascript.jscomp.parsing.parser.trees.FunctionDeclarationTree", func(v_ any) bool {
		return j2gort.InstanceOf[*FunctionDeclarationTree](v_)
	})
}
func (*FunctionDeclarationTree) Implements_parser_ParseTree() {}
func (this_ *FunctionDeclarationTree) FieldGet_name() Token   { return this_.name }
func (this_ *FunctionDeclarationTree) FieldGet_formalParameterList() *FormalParameterListTree {
	return this_.formalParameterList
}
func (this_ *FunctionDeclarationTree) FieldGet_functionBody() ParseTree            { return this_.functionBody }
func (this_ *FunctionDeclarationTree) FieldGet_isClassMember() bool                { return this_.isClassMember }
func (this_ *FunctionDeclarationTree) FieldGet_isStatic() bool                     { return this_.isStatic }
func (this_ *FunctionDeclarationTree) FieldGet_isGenerator() bool                  { return this_.isGenerator }
func (this_ *FunctionDeclarationTree) FieldGet_isOptional() bool                   { return this_.isOptional }
func (this_ *FunctionDeclarationTree) FieldGet_isAsync() bool                      { return this_.isAsync }
func (this_ *FunctionDeclarationTree) FieldGet_kind() FunctionDeclarationTree_Kind { return this_.kind }

func (this_ *FunctionDeclarationTree) initFields_FunctionDeclarationTree() {
	_ = this_
}

// Java constructor: FunctionDeclarationTree(com.google.javascript.jscomp.parsing.parser.trees.FunctionDeclarationTree.Builder builder)
func NewFunctionDeclarationTree(builder *FunctionDeclarationTree_Builder) *FunctionDeclarationTree {
	this_ := &FunctionDeclarationTree{}
	this_.BindSelf_(this_)
	InitFunctionDeclarationTree(this_, builder)
	return this_
}
func allocPooled_FunctionDeclarationTree(slab_ *j2gort.PooledSlab[FunctionDeclarationTree], builder *FunctionDeclarationTree_Builder) *FunctionDeclarationTree {
	this_ := slab_.Alloc()
	this_.BindSelf_(this_)
	InitFunctionDeclarationTree(this_, builder)
	return this_
}

func InitFunctionDeclarationTree(this_ *FunctionDeclarationTree, builder *FunctionDeclarationTree_Builder) {
	_ = builder
	InitParseTree(&this_.ParseTreeBase, StaticGet_ParseTreeType_FUNCTION_DECLARATION(), builder.location)
	this_.initFields_FunctionDeclarationTree()
	this_.name = builder.name
	this_.isClassMember = builder.isClassMember
	this_.isStatic = builder.isStatic
	this_.isGenerator = builder.isGenerator
	this_.isOptional = builder.isOptional
	this_.kind = j2gort.CheckNotNullEnum(builder.kind)
	this_.formalParameterList = j2gort.CheckNotNullPtr(builder.formalParameterList)
	this_.functionBody = j2gort.CheckNotNullIface(builder.functionBody)
	this_.isAsync = builder.isAsync
}

// Java static method: FunctionDeclarationTree.builder(com.google.javascript.jscomp.parsing.parser.trees.FunctionDeclarationTree.Kind kind)
func FunctionDeclarationTree_S_Builder(kind FunctionDeclarationTree_Kind) *FunctionDeclarationTree_Builder {
	_ = kind
	return NewFunctionDeclarationTree_Builder(kind)
}
func allocPooledVia_FunctionDeclarationTree_S_Builder(slab_ *j2gort.PooledSlab[FunctionDeclarationTree_Builder], kind FunctionDeclarationTree_Kind) *FunctionDeclarationTree_Builder {
	_ = kind
	return allocPooled_FunctionDeclarationTree_Builder(slab_, kind)
}

// Java enum: com.google.javascript.jscomp.parsing.parser.trees.FunctionDeclarationTree.Kind
type FunctionDeclarationTree_Kind uint8

const (
	enum_FunctionDeclarationTree_Kind_DECLARATION FunctionDeclarationTree_Kind = 1
	enum_FunctionDeclarationTree_Kind_EXPRESSION  FunctionDeclarationTree_Kind = 2
	enum_FunctionDeclarationTree_Kind_MEMBER      FunctionDeclarationTree_Kind = 3
	enum_FunctionDeclarationTree_Kind_ARROW       FunctionDeclarationTree_Kind = 4
)

var enumBaseTable_FunctionDeclarationTree_Kind = [...]j2gort.EnumBase{
	{Name_: "DECLARATION", Ordinal_: 0},
	{Name_: "EXPRESSION", Ordinal_: 1},
	{Name_: "MEMBER", Ordinal_: 2},
	{Name_: "ARROW", Ordinal_: 3},
}
var enumValues_FunctionDeclarationTree_Kind = []FunctionDeclarationTree_Kind{1, 2, 3, 4}

func StaticGet_FunctionDeclarationTree_Kind_DECLARATION() FunctionDeclarationTree_Kind {
	return enum_FunctionDeclarationTree_Kind_DECLARATION
}
func StaticGet_FunctionDeclarationTree_Kind_EXPRESSION() FunctionDeclarationTree_Kind {
	return enum_FunctionDeclarationTree_Kind_EXPRESSION
}
func StaticGet_FunctionDeclarationTree_Kind_MEMBER() FunctionDeclarationTree_Kind {
	return enum_FunctionDeclarationTree_Kind_MEMBER
}
func StaticGet_FunctionDeclarationTree_Kind_ARROW() FunctionDeclarationTree_Kind {
	return enum_FunctionDeclarationTree_Kind_ARROW
}

func FunctionDeclarationTree_Kind_Values() []FunctionDeclarationTree_Kind {
	return j2gort.CloneSlice(enumValues_FunctionDeclarationTree_Kind)
}

func FunctionDeclarationTree_Kind_ValueOf(name string) FunctionDeclarationTree_Kind {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_FunctionDeclarationTree_Kind {
		if enumBaseTable_FunctionDeclarationTree_Kind[i_].Name_ == name {
			return FunctionDeclarationTree_Kind(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.parsing.parser.trees.FunctionDeclarationTree.Kind." + name))
}

func (this_ FunctionDeclarationTree_Kind) IsJavaNil_() bool                              { return this_ == 0 }
func (FunctionDeclarationTree_Kind) HasIdentityEquals_()                                 {}
func (FunctionDeclarationTree_Kind) Implements_j2gort_Ext_java_lang_Enum()               {}
func (FunctionDeclarationTree_Kind) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ FunctionDeclarationTree_Kind) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (FunctionDeclarationTree_Kind) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.FunctionDeclarationTree$Kind", "Kind"
}
func (this_ FunctionDeclarationTree_Kind) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_FunctionDeclarationTree_Kind[this_-1]
}
func (this_ FunctionDeclarationTree_Kind) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ FunctionDeclarationTree_Kind) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ FunctionDeclarationTree_Kind) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_FunctionDeclarationTree_Kind[this_-1].Name_
}
func (this_ FunctionDeclarationTree_Kind) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_FunctionDeclarationTree_Kind[this_-1].Name_
}
func (this_ FunctionDeclarationTree_Kind) Equals(other_ any) bool {
	o_, ok_ := other_.(FunctionDeclarationTree_Kind)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ FunctionDeclarationTree_Kind) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_FunctionDeclarationTree_Kind[this_-1].HashCode()
}
func (this_ FunctionDeclarationTree_Kind) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_FunctionDeclarationTree_Kind[this_-1].CompareTo(other_)
}
func (this_ FunctionDeclarationTree_Kind) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_FunctionDeclarationTree_Kind[this_-1].Name_)
}
func init() {
	vals_ := []any{FunctionDeclarationTree_Kind(1), FunctionDeclarationTree_Kind(2), FunctionDeclarationTree_Kind(3), FunctionDeclarationTree_Kind(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.parser.trees.FunctionDeclarationTree.Kind", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.parsing.parser.trees.FunctionDeclarationTree$Kind", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.FunctionDeclarationTree.Kind", "com.google.javascript.jscomp.parsing.parser.trees.FunctionDeclarationTree$Kind", func(v_ any) bool {
		return j2gort.InstanceOf[FunctionDeclarationTree_Kind](v_)
	})
}

// Java class: com.google.javascript.jscomp.parsing.parser.trees.FunctionDeclarationTree.Builder
type FunctionDeclarationTree_Builder struct {
	name                Token                    `gson:"name"`
	formalParameterList *FormalParameterListTree `gson:"formalParameterList"`
	functionBody        ParseTree                `gson:"functionBody"`
	location            SourceRange              `gson:"location"`
	idHash_             atomic.Int32
	kind                FunctionDeclarationTree_Kind `gson:"kind"`
	isClassMember       bool                         `gson:"isClassMember"`
	isStatic            bool                         `gson:"isStatic"`
	isGenerator         bool                         `gson:"isGenerator"`
	isOptional          bool                         `gson:"isOptional"`
	isAsync             bool                         `gson:"isAsync"`
}

func (this_ *FunctionDeclarationTree_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *FunctionDeclarationTree_Builder) IsJavaNil_() bool { return this_ == nil }
func (*FunctionDeclarationTree_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.parsing.parser.trees.FunctionDeclarationTree$Builder", "Builder"
}
func (this_ *FunctionDeclarationTree_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*FunctionDeclarationTree_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.parsing.parser.trees.FunctionDeclarationTree.Builder", "com.google.javascript.jscomp.parsing.parser.trees.FunctionDeclarationTree$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*FunctionDeclarationTree_Builder](v_)
	})
}

func (this_ *FunctionDeclarationTree_Builder) initFields_FunctionDeclarationTree_Builder() {
	_ = this_
	this_.name = (Token{})

}

// Java constructor: Builder(com.google.javascript.jscomp.parsing.parser.trees.FunctionDeclarationTree.Kind kind)
func NewFunctionDeclarationTree_Builder(kind FunctionDeclarationTree_Kind) *FunctionDeclarationTree_Builder {
	this_ := &FunctionDeclarationTree_Builder{}
	InitFunctionDeclarationTree_Builder(this_, kind)
	return this_
}
func allocPooled_FunctionDeclarationTree_Builder(slab_ *j2gort.PooledSlab[FunctionDeclarationTree_Builder], kind FunctionDeclarationTree_Kind) *FunctionDeclarationTree_Builder {
	this_ := slab_.Alloc()
	InitFunctionDeclarationTree_Builder(this_, kind)
	return this_
}

func InitFunctionDeclarationTree_Builder(this_ *FunctionDeclarationTree_Builder, kind FunctionDeclarationTree_Kind) {
	_ = kind
	this_.initFields_FunctionDeclarationTree_Builder()
	this_.kind = kind
}

// Java method: Builder.setName(com.google.javascript.jscomp.parsing.parser.Token name)
func (this_ *FunctionDeclarationTree_Builder) SetName(name Token) *FunctionDeclarationTree_Builder {
	_ = name
	this_.name = name
	return this_
}

// Java method: Builder.setFormalParameterList(com.google.javascript.jscomp.parsing.parser.trees.FormalParameterListTree formalParameterList)
func (this_ *FunctionDeclarationTree_Builder) SetFormalParameterList(formalParameterList *FormalParameterListTree) *FunctionDeclarationTree_Builder {
	_ = formalParameterList
	this_.formalParameterList = formalParameterList
	return this_
}

// Java method: Builder.setFunctionBody(com.google.javascript.jscomp.parsing.parser.trees.ParseTree functionBody)
func (this_ *FunctionDeclarationTree_Builder) SetFunctionBody(functionBody ParseTree) *FunctionDeclarationTree_Builder {
	_ = functionBody
	this_.functionBody = functionBody
	return this_
}

// Java method: Builder.setIsClassMember(boolean isClassMember)
func (this_ *FunctionDeclarationTree_Builder) SetIsClassMember(isClassMember bool) *FunctionDeclarationTree_Builder {
	_ = isClassMember
	this_.isClassMember = isClassMember
	return this_
}

// Java method: Builder.setStatic(boolean isStatic)
func (this_ *FunctionDeclarationTree_Builder) SetStatic(isStatic bool) *FunctionDeclarationTree_Builder {
	_ = isStatic
	this_.isStatic = isStatic
	return this_
}

// Java method: Builder.setGenerator(boolean isGenerator)
func (this_ *FunctionDeclarationTree_Builder) SetGenerator(isGenerator bool) *FunctionDeclarationTree_Builder {
	_ = isGenerator
	this_.isGenerator = isGenerator
	return this_
}

// Java method: Builder.setOptional(boolean isOptional)
func (this_ *FunctionDeclarationTree_Builder) SetOptional(isOptional bool) *FunctionDeclarationTree_Builder {
	_ = isOptional
	this_.isOptional = isOptional
	return this_
}

// Java method: Builder.setAsync(boolean isAsync)
func (this_ *FunctionDeclarationTree_Builder) SetAsync(isAsync bool) *FunctionDeclarationTree_Builder {
	_ = isAsync
	this_.isAsync = isAsync
	return this_
}

// Java method: Builder.build(com.google.javascript.jscomp.parsing.parser.util.SourceRange location)
func (this_ *FunctionDeclarationTree_Builder) Build(location SourceRange) *FunctionDeclarationTree {
	_ = location
	this_.location = location
	return NewFunctionDeclarationTree(this_)
}
func (this_ *FunctionDeclarationTree_Builder) buildWithPooledSlab_FunctionDeclarationTree(slab_ *j2gort.PooledSlab[FunctionDeclarationTree], location SourceRange) *FunctionDeclarationTree {
	_ = location
	this_.location = location
	return allocPooled_FunctionDeclarationTree(slab_, this_)
}
