// Port of com/google/javascript/rhino/TypeDeclarationsIR.java @ closure-compiler v20260415

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Michael Zhou
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// TypeDeclarationsIR is an AST construction helper class for Node.
type TypeDeclarationsIR struct{}

// TypeDeclarationsIR_StringType returns a new node representing the string
// built-in type.
//
// Java: TypeDeclarationsIR.stringType()
func TypeDeclarationsIR_StringType() *Node {
	return NewNode(Token_STRING_TYPE)
}

// TypeDeclarationsIR_NumberType returns a new node representing the number
// built-in type.
//
// Java: TypeDeclarationsIR.numberType()
func TypeDeclarationsIR_NumberType() *Node {
	return NewNode(Token_NUMBER_TYPE)
}

// TypeDeclarationsIR_BooleanType returns a new node representing the
// boolean built-in type.
//
// Java: TypeDeclarationsIR.booleanType()
func TypeDeclarationsIR_BooleanType() *Node {
	return NewNode(Token_BOOLEAN_TYPE)
}

// TypeDeclarationsIR_AnyType is equivalent to the UNKNOWN type in Closure,
// expressed with {?}. Returns a new node representing any type, without type
// checking.
//
// Java: TypeDeclarationsIR.anyType()
func TypeDeclarationsIR_AnyType() *Node {
	return NewNode(Token_ANY_TYPE)
}

// TypeDeclarationsIR_VoidType returns a new node representing the Void type
// as defined by TypeScript.
//
// Java: TypeDeclarationsIR.voidType()
func TypeDeclarationsIR_VoidType() *Node {
	return NewNode(Token_VOID_TYPE)
}

// TypeDeclarationsIR_UndefinedType returns a new node representing the
// Undefined type as defined by TypeScript.
//
// Java: TypeDeclarationsIR.undefinedType()
func TypeDeclarationsIR_UndefinedType() *Node {
	return NewNode(Token_UNDEFINED_TYPE)
}

// TypeDeclarationsIR_NamedType splits a '.' separated qualified name into a
// tree of type segments.
//
// typeName: a qualified name such as "goog.ui.Window". Returns a new node
// representing the type
//
// Java: TypeDeclarationsIR.namedType(String)
func TypeDeclarationsIR_NamedType(typeName string) *Node {
	return TypeDeclarationsIR_NamedTypeOfSegments(JSTypeRegistry_DOT_SPLITTER_Split(typeName))
}

// TypeDeclarationsIR_NamedTypeOfSegments produces a tree structure similar
// to the Rhino AST of a qualified name expression, under a top-level
// NAMED_TYPE node.
//
// Example:
//
//	NAMED_TYPE
//	  NAME goog
//	    STRING ui
//	      STRING Window
//
// Java: TypeDeclarationsIR.namedType(Iterable<String>)
func TypeDeclarationsIR_NamedTypeOfSegments(segments []string) *Node {
	if len(segments) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	node := IR_Name(segments[0])
	for _, segment := range segments[1:] {
		node = IR_Getprop(node, segment)
	}
	return NewNodeWithChild(Token_NAMED_TYPE, node)
}

// TypeDeclarationsIR_RecordType represents a structural type. Closure calls
// this a Record Type and accepts the syntax {myNum: number, myObject}
//
// Example:
//
//	RECORD_TYPE
//	  STRING_KEY myNum
//	    NUMBER_TYPE
//	  STRING_KEY myObject
//
// properties: a map from property name to property type. Returns a new node
// representing the record type
//
// Java: TypeDeclarationsIR.recordType(LinkedHashMap<String, Node>)
func TypeDeclarationsIR_RecordType(properties *jutil.LinkedHashMap[string, *Node]) *Node {
	node := NewNode(Token_RECORD_TYPE)
	for key, value := range properties.All() {
		stringKey := IR_StringKey(key)
		node.AddChildToBack(stringKey)
		if value != nil {
			stringKey.AddChildToFront(value)
		}
	}
	return node
}

// Java: TypeDeclarationsIR.maybeAddType(Node, Node)
func TypeDeclarationsIR_MaybeAddType(node *Node, type_ *Node) *Node {
	if type_ != nil {
		node.SetDeclaredTypeExpression(type_)
	}
	return node
}

// TypeDeclarationsIR_FunctionType represents a function type. Closure has
// syntax like {function(string, boolean):number} Closure doesn't include
// parameter names. If the parameter types are unnamed, arbitrary names can
// be substituted, eg. p1, p2, etc.
//
// Example:
//
//	FUNCTION_TYPE
//	  NUMBER_TYPE
//	  STRING_KEY p1 [declared_type_expr: STRING_TYPE]
//	  STRING_KEY p2 [declared_type_expr: BOOLEAN_TYPE]
//
// returnType: the type returned by the function, possibly ANY_TYPE.
// requiredParams: the names and types of the required parameters.
// optionalParams: the names and types of the optional parameters. restName:
// the name of the rest parameter, if any. restType: the type of the rest
// parameter, if any.
//
// Java: TypeDeclarationsIR.functionType(Node, LinkedHashMap<String, Node>, LinkedHashMap<String, Node>, String, Node)
func TypeDeclarationsIR_FunctionType(
	returnType *Node,
	requiredParams *jutil.LinkedHashMap[string, *Node],
	optionalParams *jutil.LinkedHashMap[string, *Node],
	restName *string,
	restType *Node) *Node {
	node := NewNodeWithChild(Token_FUNCTION_TYPE, returnType)
	jutil.CheckNotNull(requiredParams)
	jutil.CheckNotNull(optionalParams)

	for key, value := range requiredParams.All() {
		name := IR_Name(key)
		node.AddChildToBack(TypeDeclarationsIR_MaybeAddType(name, value))
	}

	for key, value := range optionalParams.All() {
		name := IR_Name(key)
		name.PutBooleanProp(Node_OPT_ES6_TYPED, true)
		node.AddChildToBack(TypeDeclarationsIR_MaybeAddType(name, value))
	}

	if restName != nil {
		rest := NewNodeWithChild(Token_ITER_REST, IR_Name(*restName))
		node.AddChildToBack(TypeDeclarationsIR_MaybeAddType(rest, restType))
	}
	return node
}

// TypeDeclarationsIR_ParameterizedType represents a parameterized, or
// generic, type. Closure calls this a Type Application and accepts syntax
// like {Object.<string, number>}
//
// Example:
//
//	PARAMETERIZED_TYPE
//	  NAMED_TYPE
//	    NAME Object
//	  STRING_TYPE
//	  NUMBER_TYPE
//
// Java: TypeDeclarationsIR.parameterizedType(Node, Iterable<Node>)
func TypeDeclarationsIR_ParameterizedType(baseType *Node, typeParameters []*Node) *Node {
	if len(typeParameters) == 0 {
		return baseType
	}
	node := NewNodeWithChild(Token_PARAMETERIZED_TYPE, baseType)
	for _, typeParameter := range typeParameters {
		node.AddChildToBack(typeParameter)
	}
	return node
}

// TypeDeclarationsIR_ArrayType represents an array type. In Closure, this is
// represented by a parameterized type of Array with elementType as the sole
// type parameter.
//
// Example
//
//	ARRAY_TYPE
//	  elementType
//
// Java: TypeDeclarationsIR.arrayType(Node)
func TypeDeclarationsIR_ArrayType(elementType *Node) *Node {
	return NewNodeWithChild(Token_ARRAY_TYPE, elementType)
}

// TypeDeclarationsIR_UnionType represents a union type, which can be one of
// the given types. Closure accepts syntax like {(number|boolean)}
//
// Example:
//
//	UNION_TYPE
//	  NUMBER_TYPE
//	  BOOLEAN_TYPE
//
// options: the types which are accepted. Returns a new node representing the
// union type
//
// Java: TypeDeclarationsIR.unionType(Iterable<Node>)
func TypeDeclarationsIR_UnionType(options []*Node) *Node {
	jutil.CheckArgumentMsg(len(options) != 0, "union must have at least one option")
	node := NewNode(Token_UNION_TYPE)
	for _, option := range options {
		node.AddChildToBack(option)
	}
	return node
}

// Java: TypeDeclarationsIR.unionType(Node...)
func TypeDeclarationsIR_UnionTypeOfNodes(options ...*Node) *Node {
	return TypeDeclarationsIR_UnionType(options)
}

// TypeDeclarationsIR_OptionalParameter represents a function parameter that
// is optional. In closure syntax, this is function(?string=, number=) In
// TypeScript syntax, it is (firstName: string, lastName?: string)=>string
//
// parameterType: the type of the parameter. Returns a new node representing
// the function parameter type
//
// Java: TypeDeclarationsIR.optionalParameter(Node)
func TypeDeclarationsIR_OptionalParameter(parameterType *Node) *Node {
	return NewNodeWithChild(Token_OPTIONAL_PARAMETER, parameterType)
}

// Java: TypeDeclarationsIR()
func NewTypeDeclarationsIR() *TypeDeclarationsIR {
	return &TypeDeclarationsIR{}
}
