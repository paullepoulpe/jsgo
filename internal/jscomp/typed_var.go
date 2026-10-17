// Port of com/google/javascript/jscomp/TypedVar.java @ closure-compiler v20261015

/*
 * Copyright 2015 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TypedVar is the AbstractVar subclass for use with TypedScope.
//
// Note that this class inherits its equals and hashCode implementations
// from ScopedName, which does not include any type information. This is
// necessary because Var-keyed maps are used across multiple top scopes, but
// it comes with the caveat that if TypedVar instances are stored in a set,
// the type information is at risk of disappearing if an untyped (or
// differently typed) var is added for the same symbol.
//
// It is a rhino.StaticTypedSlot and rhino.StaticTypedRef; GetScope,
// GetDeclaration and GetSymbol keep rhino's signatures (assert
// `.(*TypedScope)` / `.(*TypedVar)`); see the note at the top of
// abstract_scope.go.
type TypedVar struct {
	AbstractVarBase[*TypedScope, *TypedVar]

	type_ rhino.JSType
	// The next two fields and the associated methods are only used by
	// TypeInference.java. Maybe there is a way to avoid having them in all typed variable instances.
	markedEscaped             bool
	markedAssignedExactlyOnce bool

	// Whether the variable's type has been inferred or is declared. An
	// inferred type may change over time (as more code is discovered),
	// whereas a declared type is a static contract that must be matched.
	typeInferred bool
}

// includes nodes that in plain JS semantics are not 'declarations', but that the type system
// & compiler treat as declarations.
var typedVar_NAME_NODE_TYPES = jutil.NewEnumSet(
	rhino.Token_NAME,
	rhino.Token_THIS, // `if ('prop' in this) {` creates a TypedVar `this.prop` requiring a node
	rhino.Token_IMPORT_STAR,
	rhino.Token_EXPR_RESULT, // implicit variables from goog.provide
	rhino.Token_EXPORT,      // tracks the *default* export
	rhino.Token_GETPROP,
	rhino.Token_FUNCTION,
	rhino.Token_STRING_KEY,
	rhino.Token_SETTER_DEF,
	rhino.Token_GETTER_DEF,
	rhino.Token_MEMBER_FUNCTION_DEF,
	rhino.Token_MODULE_BODY) // for the implicit exports object in a module

// Java: TypedVar(boolean, String, Node, JSType, TypedScope, int, CompilerInput)
func NewTypedVar(inferred bool, name string, nameNode *rhino.Node, type_ rhino.JSType,
	scope *TypedScope, index int, input *CompilerInput) *TypedVar {
	v := &TypedVar{}
	v.InitAbstractVarBase(v, name, nameNode, scope, index, input, nil /* implicitGoogNamespaceDefinition */)
	if nameNode != nil {
		jutil.CheckArgumentMsg1(
			typedVar_NAME_NODE_TYPES.Contains(nameNode.GetToken()),
			"Invalid name node token %s",
			nameNode.GetToken())
	}
	v.type_ = type_
	v.typeInferred = inferred
	return v
}

// GetType gets this variable's type. To know whether this type has been
// inferred, see IsTypeInferred.
//
// Java: TypedVar.getType()
func (v *TypedVar) GetType() rhino.JSType {
	return v.type_
}

// Java: TypedVar.setType(JSType)
func (v *TypedVar) SetType(type_ rhino.JSType) {
	v.type_ = type_
}

// Java: TypedVar.resolveType(ErrorReporter)
func (v *TypedVar) ResolveType(errorReporter rhino.ErrorReporter) {
	if v.type_ != nil {
		v.type_ = v.type_.Resolve(errorReporter)
	}
}

// IsTypeInferred returns whether this variable's type is inferred. To get
// the variable's type, see GetType.
//
// Java: TypedVar.isTypeInferred()
func (v *TypedVar) IsTypeInferred() bool {
	return v.typeInferred
}

// Java: TypedVar.getInputName()
func (v *TypedVar) GetInputName() string {
	if v.GetInput() == nil {
		return "<non-file>"
	}
	return v.GetInput().GetName()
}

// Java: TypedVar.toString()
func (v *TypedVar) String() string {
	return "Var " + v.GetName() + "{" + jutil.ToJavaString(v.type_) + "}"
}

// Java: TypedVar.markEscaped()
func (v *TypedVar) MarkEscaped() {
	v.markedEscaped = true
}

// Java: TypedVar.isMarkedEscaped()
func (v *TypedVar) IsMarkedEscaped() bool {
	return v.markedEscaped
}

// Java: TypedVar.markAssignedExactlyOnce()
func (v *TypedVar) MarkAssignedExactlyOnce() {
	v.markedAssignedExactlyOnce = true
}

// Java: TypedVar.isMarkedAssignedExactlyOnce()
func (v *TypedVar) IsMarkedAssignedExactlyOnce() bool {
	return v.markedAssignedExactlyOnce
}
