// Port of com/google/javascript/jscomp/parsing/parser/trees/FunctionDeclarationTree.java @ closure-compiler v20261015

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

// FunctionDeclarationTree is the parse tree for ParseTreeType
// FUNCTION_DECLARATION (declarations, expressions, methods and arrows).
type FunctionDeclarationTree struct {
	ParseTreeBase

	Name                *IdentifierToken
	FormalParameterList *FormalParameterListTree
	FunctionBody        ParseTree
	IsClassMember       bool
	IsStatic            bool
	IsGenerator         bool
	IsOptional          bool
	IsAsync             bool
	Kind                FunctionDeclarationTree_Kind
}

// FunctionDeclarationTree_Kind is FunctionDeclarationTree.Kind.
type FunctionDeclarationTree_Kind int32

const (
	FunctionDeclarationTree_Kind_DECLARATION FunctionDeclarationTree_Kind = iota
	FunctionDeclarationTree_Kind_EXPRESSION
	FunctionDeclarationTree_Kind_MEMBER
	FunctionDeclarationTree_Kind_ARROW
)

var functionDeclarationTree_Kind_names = [...]string{"DECLARATION", "EXPRESSION", "MEMBER", "ARROW"}

// FunctionDeclarationTree_Kind_values is Kind.values().
//
// Java: FunctionDeclarationTree.Kind.values()
func FunctionDeclarationTree_Kind_values() []FunctionDeclarationTree_Kind {
	return []FunctionDeclarationTree_Kind{
		FunctionDeclarationTree_Kind_DECLARATION,
		FunctionDeclarationTree_Kind_EXPRESSION,
		FunctionDeclarationTree_Kind_MEMBER,
		FunctionDeclarationTree_Kind_ARROW,
	}
}

// Name is Enum.name().
//
// Java: FunctionDeclarationTree.Kind.name()
func (k FunctionDeclarationTree_Kind) Name() string { return functionDeclarationTree_Kind_names[k] }

// String is Enum.toString().
//
// Java: FunctionDeclarationTree.Kind.toString()
func (k FunctionDeclarationTree_Kind) String() string { return k.Name() }

// FunctionDeclarationTree_NewBuilder returns a builder for a function of the
// given kind.
//
// Go-only: language-adapter — the mechanical name FunctionDeclarationTree_Builder
// is taken by the nested Builder type.
//
// Java: FunctionDeclarationTree.builder(Kind)
func FunctionDeclarationTree_NewBuilder(kind FunctionDeclarationTree_Kind) *FunctionDeclarationTree_Builder {
	return NewFunctionDeclarationTree_Builder(kind)
}

// NewFunctionDeclarationTree constructs a FunctionDeclarationTree from a
// builder.
//
// Java: FunctionDeclarationTree(Builder)
func NewFunctionDeclarationTree(builder *FunctionDeclarationTree_Builder) *FunctionDeclarationTree {
	t := &FunctionDeclarationTree{}
	t.InitParseTreeBase(t, ParseTreeType_FUNCTION_DECLARATION, builder.location)

	t.Name = builder.name
	t.IsClassMember = builder.isClassMember
	t.IsStatic = builder.isStatic
	t.IsGenerator = builder.isGenerator
	t.IsOptional = builder.isOptional
	t.Kind = builder.kind // Java: checkNotNull(builder.kind); an enum value is never null in Go.
	t.FormalParameterList = jutil.CheckNotNull(builder.formalParameterList)
	t.FunctionBody = jutil.CheckNotNull(builder.functionBody)
	t.IsAsync = builder.isAsync
	return t
}

// FunctionDeclarationTree_Builder builds a FunctionDeclarationTree.
type FunctionDeclarationTree_Builder struct {
	kind FunctionDeclarationTree_Kind

	name                *IdentifierToken
	formalParameterList *FormalParameterListTree
	functionBody        ParseTree
	isClassMember       bool
	isStatic            bool
	isGenerator         bool
	isOptional          bool
	isAsync             bool
	location            SourceRange
	arena               *ParseArena
}

// NewFunctionDeclarationTree_Builder constructs a Builder.
//
// Java: FunctionDeclarationTree.Builder(Kind)
func NewFunctionDeclarationTree_Builder(kind FunctionDeclarationTree_Kind) *FunctionDeclarationTree_Builder {
	b := &FunctionDeclarationTree_Builder{}
	b.kind = kind
	return b
}

// SetName sets the optional function name.
//
// Default is nil.
//
// Java: FunctionDeclarationTree.Builder.setName(IdentifierToken)
func (b *FunctionDeclarationTree_Builder) SetName(name *IdentifierToken) *FunctionDeclarationTree_Builder {
	b.name = name
	return b
}

// SetFormalParameterList sets the required parameter list.
//
// Java: FunctionDeclarationTree.Builder.setFormalParameterList(FormalParameterListTree)
func (b *FunctionDeclarationTree_Builder) SetFormalParameterList(formalParameterList *FormalParameterListTree) *FunctionDeclarationTree_Builder {
	b.formalParameterList = formalParameterList
	return b
}

// SetFunctionBody sets the required function body.
//
// Java: FunctionDeclarationTree.Builder.setFunctionBody(ParseTree)
func (b *FunctionDeclarationTree_Builder) SetFunctionBody(functionBody ParseTree) *FunctionDeclarationTree_Builder {
	b.functionBody = functionBody
	return b
}

// SetIsClassMember: is the method an ES6 class member?
//
// Default is false. Only relevant for class method member declarations.
//
// Java: FunctionDeclarationTree.Builder.setIsClassMember(boolean)
func (b *FunctionDeclarationTree_Builder) SetIsClassMember(isClassMember bool) *FunctionDeclarationTree_Builder {
	b.isClassMember = isClassMember
	return b
}

// SetStatic: is the method static?
//
// Default is false. Only relevant for class method member declarations.
//
// Java: FunctionDeclarationTree.Builder.setStatic(boolean)
func (b *FunctionDeclarationTree_Builder) SetStatic(isStatic bool) *FunctionDeclarationTree_Builder {
	b.isStatic = isStatic
	return b
}

// SetGenerator: is this a generator function?
//
// Default is false.
//
// Java: FunctionDeclarationTree.Builder.setGenerator(boolean)
func (b *FunctionDeclarationTree_Builder) SetGenerator(isGenerator bool) *FunctionDeclarationTree_Builder {
	b.isGenerator = isGenerator
	return b
}

// SetOptional: is this the declaration of an optional function parameter?
// Default is false.
//
// Only relevant for function declaration as a parameter to another
// function.
//
// Java: FunctionDeclarationTree.Builder.setOptional(boolean)
func (b *FunctionDeclarationTree_Builder) SetOptional(isOptional bool) *FunctionDeclarationTree_Builder {
	b.isOptional = isOptional
	return b
}

// SetAsync: is this an asynchronous function?
//
// Default is false.
//
// Java: FunctionDeclarationTree.Builder.setAsync(boolean)
func (b *FunctionDeclarationTree_Builder) SetAsync(isAsync bool) *FunctionDeclarationTree_Builder {
	b.isAsync = isAsync
	return b
}

// Build returns a new FunctionDeclarationTree.
//
// The location is provided at this point because it cannot be correctly
// calculated until the whole function has been parsed.
//
// Java: FunctionDeclarationTree.Builder.build(SourceRange)
func (b *FunctionDeclarationTree_Builder) Build(location SourceRange) *FunctionDeclarationTree {
	b.location = location
	if b.arena != nil {
		t := b.arena.functionDeclarations.New()
		t.InitParseTreeBase(t, ParseTreeType_FUNCTION_DECLARATION, b.location)
		t.Name = b.name
		t.IsClassMember = b.isClassMember
		t.IsStatic = b.isStatic
		t.IsGenerator = b.isGenerator
		t.IsOptional = b.isOptional
		t.Kind = b.kind
		t.FormalParameterList = jutil.CheckNotNull(b.formalParameterList)
		t.FunctionBody = jutil.CheckNotNull(builderFunctionBody(b))
		t.IsAsync = b.isAsync
		return t
	}
	return NewFunctionDeclarationTree(b)
}

func builderFunctionBody(b *FunctionDeclarationTree_Builder) ParseTree {
	return b.functionBody
}
