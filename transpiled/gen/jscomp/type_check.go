/*
 * Copyright 2006 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/TypeCheck.java (generated mechanically by j2go).

package jscomp

import (
	"strconv"
	"sync/atomic"

	pkg_base "github.com/paullepoulpe/jsgo/transpiled/gen/base"
	pkg_diagnostic "github.com/paullepoulpe/jsgo/transpiled/gen/diagnostic"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_jscomptype "github.com/paullepoulpe/jsgo/transpiled/gen/jscomptype"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_type_check_1 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_type_check_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return j2gort.BoxString(j2gort.AsIface[pkg_rhino.ObjectType](mrp0_).GetDisplayName())
	})
}

// Java class: com.google.javascript.jscomp.TypeCheck
type TypeCheck struct {
	compiler                   AbstractCompiler                          `gson:"compiler"`
	validator                  *TypeValidator                            `gson:"validator"`
	reverseInterpreter         pkg_jscomptype.ReverseAbstractInterpreter `gson:"reverseInterpreter"`
	typeRegistry               *pkg_rhino.JSTypeRegistry                 `gson:"typeRegistry"`
	topScope                   *TypedScope                               `gson:"topScope"`
	scopeCreator               *TypedScopeCreator                        `gson:"scopeCreator"`
	inferJSDocInfo             *InferJSDocInfo                           `gson:"inferJSDocInfo"`
	debugTypeLogger            *TypeCheck_DebugTypeLogger                `gson:"debugTypeLogger"`
	nonStringifiableKeyVisited j2gort.CustomOrderedSet[pkg_rhino.JSType] `gson:"nonStringifiableKeyVisited"`
	typedCount                 int32                                     `gson:"typedCount"`
	nullCount                  int32                                     `gson:"nullCount"`
	unknownCount               int32                                     `gson:"unknownCount"`
	idHash_                    atomic.Int32
	reportUnknownTypes         bool                           `gson:"reportUnknownTypes"`
	subtypingMode              pkg_rhino.JSType_SubtypingMode `gson:"subtypingMode"`
	reportMissingProperties    bool                           `gson:"reportMissingProperties"`
	inExterns                  bool                           `gson:"inExterns"`
}

func (this_ *TypeCheck) BindSelf_(s any) {
	_ = s
}
func (this_ *TypeCheck) IsJavaNil_() bool { return this_ == nil }
func (*TypeCheck) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypeCheck", "TypeCheck"
}
func (this_ *TypeCheck) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypeCheck) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypeCheck", "com.google.javascript.jscomp.TypeCheck", func(v_ any) bool {
		return j2gort.InstanceOf[*TypeCheck](v_)
	})
}
func (*TypeCheck) Implements_jscomp_NodeTraversal_Callback() {}
func (*TypeCheck) Implements_jscomp_CompilerPass()           {}

var static_TypeCheck_UNEXPECTED_TOKEN *DiagnosticType

func StaticGet_TypeCheck_UNEXPECTED_TOKEN() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_UNEXPECTED_TOKEN
}
func StaticSet_TypeCheck_UNEXPECTED_TOKEN(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_UNEXPECTED_TOKEN = v
}

var static_TypeCheck_DETERMINISTIC_TEST *DiagnosticType

func StaticGet_TypeCheck_DETERMINISTIC_TEST() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_DETERMINISTIC_TEST
}
func StaticSet_TypeCheck_DETERMINISTIC_TEST(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_DETERMINISTIC_TEST = v
}

var static_TypeCheck_INEXISTENT_ENUM_ELEMENT *DiagnosticType

func StaticGet_TypeCheck_INEXISTENT_ENUM_ELEMENT() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_INEXISTENT_ENUM_ELEMENT
}
func StaticSet_TypeCheck_INEXISTENT_ENUM_ELEMENT(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_INEXISTENT_ENUM_ELEMENT = v
}

var static_TypeCheck_INEXISTENT_PROPERTY *DiagnosticType

func StaticGet_TypeCheck_INEXISTENT_PROPERTY() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_INEXISTENT_PROPERTY
}
func StaticSet_TypeCheck_INEXISTENT_PROPERTY(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_INEXISTENT_PROPERTY = v
}
func StaticGet_TypeCheck_POSSIBLE_INEXISTENT_PROPERTY_EXPLANATION() string {
	return "\n\nThis property is accessed on a \"loose\" type, but is not defined anywhere in the program, so it must not exist."
}

var static_TypeCheck_POSSIBLE_INEXISTENT_PROPERTY *DiagnosticType

func StaticGet_TypeCheck_POSSIBLE_INEXISTENT_PROPERTY() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_POSSIBLE_INEXISTENT_PROPERTY
}
func StaticSet_TypeCheck_POSSIBLE_INEXISTENT_PROPERTY(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_POSSIBLE_INEXISTENT_PROPERTY = v
}

var static_TypeCheck_INEXISTENT_PROPERTY_WITH_SUGGESTION *DiagnosticType

func StaticGet_TypeCheck_INEXISTENT_PROPERTY_WITH_SUGGESTION() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_INEXISTENT_PROPERTY_WITH_SUGGESTION
}
func StaticSet_TypeCheck_INEXISTENT_PROPERTY_WITH_SUGGESTION(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_INEXISTENT_PROPERTY_WITH_SUGGESTION = v
}

var static_TypeCheck_STRICT_INEXISTENT_PROPERTY *DiagnosticType

func StaticGet_TypeCheck_STRICT_INEXISTENT_PROPERTY() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_STRICT_INEXISTENT_PROPERTY
}
func StaticSet_TypeCheck_STRICT_INEXISTENT_PROPERTY(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_STRICT_INEXISTENT_PROPERTY = v
}

var static_TypeCheck_STRICT_INEXISTENT_UNION_PROPERTY *DiagnosticType

func StaticGet_TypeCheck_STRICT_INEXISTENT_UNION_PROPERTY() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_STRICT_INEXISTENT_UNION_PROPERTY
}
func StaticSet_TypeCheck_STRICT_INEXISTENT_UNION_PROPERTY(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_STRICT_INEXISTENT_UNION_PROPERTY = v
}

var static_TypeCheck_STRICT_INEXISTENT_PROPERTY_WITH_SUGGESTION *DiagnosticType

func StaticGet_TypeCheck_STRICT_INEXISTENT_PROPERTY_WITH_SUGGESTION() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_STRICT_INEXISTENT_PROPERTY_WITH_SUGGESTION
}
func StaticSet_TypeCheck_STRICT_INEXISTENT_PROPERTY_WITH_SUGGESTION(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_STRICT_INEXISTENT_PROPERTY_WITH_SUGGESTION = v
}

var static_TypeCheck_NOT_A_CONSTRUCTOR *DiagnosticType

func StaticGet_TypeCheck_NOT_A_CONSTRUCTOR() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_NOT_A_CONSTRUCTOR
}
func StaticSet_TypeCheck_NOT_A_CONSTRUCTOR(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_NOT_A_CONSTRUCTOR = v
}

var static_TypeCheck_INSTANTIATE_ABSTRACT_CLASS *DiagnosticType

func StaticGet_TypeCheck_INSTANTIATE_ABSTRACT_CLASS() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_INSTANTIATE_ABSTRACT_CLASS
}
func StaticSet_TypeCheck_INSTANTIATE_ABSTRACT_CLASS(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_INSTANTIATE_ABSTRACT_CLASS = v
}

var static_TypeCheck_BIT_OPERATION *DiagnosticType

func StaticGet_TypeCheck_BIT_OPERATION() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_BIT_OPERATION
}
func StaticSet_TypeCheck_BIT_OPERATION(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_BIT_OPERATION = v
}

var static_TypeCheck_UNARY_OPERATION *DiagnosticType

func StaticGet_TypeCheck_UNARY_OPERATION() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_UNARY_OPERATION
}
func StaticSet_TypeCheck_UNARY_OPERATION(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_UNARY_OPERATION = v
}

var static_TypeCheck_BINARY_OPERATION *DiagnosticType

func StaticGet_TypeCheck_BINARY_OPERATION() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_BINARY_OPERATION
}
func StaticSet_TypeCheck_BINARY_OPERATION(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_BINARY_OPERATION = v
}

var static_TypeCheck_NOT_CALLABLE *DiagnosticType

func StaticGet_TypeCheck_NOT_CALLABLE() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_NOT_CALLABLE
}
func StaticSet_TypeCheck_NOT_CALLABLE(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_NOT_CALLABLE = v
}

var static_TypeCheck_CONSTRUCTOR_NOT_CALLABLE *DiagnosticType

func StaticGet_TypeCheck_CONSTRUCTOR_NOT_CALLABLE() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_CONSTRUCTOR_NOT_CALLABLE
}
func StaticSet_TypeCheck_CONSTRUCTOR_NOT_CALLABLE(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_CONSTRUCTOR_NOT_CALLABLE = v
}

var static_TypeCheck_ABSTRACT_SUPER_METHOD_NOT_USABLE *DiagnosticType

func StaticGet_TypeCheck_ABSTRACT_SUPER_METHOD_NOT_USABLE() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_ABSTRACT_SUPER_METHOD_NOT_USABLE
}
func StaticSet_TypeCheck_ABSTRACT_SUPER_METHOD_NOT_USABLE(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_ABSTRACT_SUPER_METHOD_NOT_USABLE = v
}

var static_TypeCheck_FUNCTION_MASKS_VARIABLE *DiagnosticType

func StaticGet_TypeCheck_FUNCTION_MASKS_VARIABLE() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_FUNCTION_MASKS_VARIABLE
}
func StaticSet_TypeCheck_FUNCTION_MASKS_VARIABLE(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_FUNCTION_MASKS_VARIABLE = v
}

var static_TypeCheck_MULTIPLE_VAR_DEF *DiagnosticType

func StaticGet_TypeCheck_MULTIPLE_VAR_DEF() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_MULTIPLE_VAR_DEF
}
func StaticSet_TypeCheck_MULTIPLE_VAR_DEF(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_MULTIPLE_VAR_DEF = v
}

var static_TypeCheck_INVALID_INTERFACE_MEMBER_DECLARATION *DiagnosticType

func StaticGet_TypeCheck_INVALID_INTERFACE_MEMBER_DECLARATION() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_INVALID_INTERFACE_MEMBER_DECLARATION
}
func StaticSet_TypeCheck_INVALID_INTERFACE_MEMBER_DECLARATION(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_INVALID_INTERFACE_MEMBER_DECLARATION = v
}

var static_TypeCheck_INTERFACE_METHOD_NOT_EMPTY *DiagnosticType

func StaticGet_TypeCheck_INTERFACE_METHOD_NOT_EMPTY() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_INTERFACE_METHOD_NOT_EMPTY
}
func StaticSet_TypeCheck_INTERFACE_METHOD_NOT_EMPTY(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_INTERFACE_METHOD_NOT_EMPTY = v
}

var static_TypeCheck_CONFLICTING_EXTENDED_TYPE *DiagnosticType

func StaticGet_TypeCheck_CONFLICTING_EXTENDED_TYPE() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_CONFLICTING_EXTENDED_TYPE
}
func StaticSet_TypeCheck_CONFLICTING_EXTENDED_TYPE(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_CONFLICTING_EXTENDED_TYPE = v
}

var static_TypeCheck_ES5_CLASS_EXTENDING_ES6_CLASS *DiagnosticType

func StaticGet_TypeCheck_ES5_CLASS_EXTENDING_ES6_CLASS() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_ES5_CLASS_EXTENDING_ES6_CLASS
}
func StaticSet_TypeCheck_ES5_CLASS_EXTENDING_ES6_CLASS(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_ES5_CLASS_EXTENDING_ES6_CLASS = v
}

var static_TypeCheck_DICT_EXTEND_STRUCT_TYPE *DiagnosticType

func StaticGet_TypeCheck_DICT_EXTEND_STRUCT_TYPE() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_DICT_EXTEND_STRUCT_TYPE
}
func StaticSet_TypeCheck_DICT_EXTEND_STRUCT_TYPE(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_DICT_EXTEND_STRUCT_TYPE = v
}

var static_TypeCheck_STRUCT_EXTEND_DICT_TYPE *DiagnosticType

func StaticGet_TypeCheck_STRUCT_EXTEND_DICT_TYPE() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_STRUCT_EXTEND_DICT_TYPE
}
func StaticSet_TypeCheck_STRUCT_EXTEND_DICT_TYPE(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_STRUCT_EXTEND_DICT_TYPE = v
}

var static_TypeCheck_ES6_CLASS_EXTENDING_CLASS_WITH_GOOG_INHERITS *DiagnosticType

func StaticGet_TypeCheck_ES6_CLASS_EXTENDING_CLASS_WITH_GOOG_INHERITS() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_ES6_CLASS_EXTENDING_CLASS_WITH_GOOG_INHERITS
}
func StaticSet_TypeCheck_ES6_CLASS_EXTENDING_CLASS_WITH_GOOG_INHERITS(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_ES6_CLASS_EXTENDING_CLASS_WITH_GOOG_INHERITS = v
}

var static_TypeCheck_INTERFACE_EXTENDS_LOOP *DiagnosticType

func StaticGet_TypeCheck_INTERFACE_EXTENDS_LOOP() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_INTERFACE_EXTENDS_LOOP
}
func StaticSet_TypeCheck_INTERFACE_EXTENDS_LOOP(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_INTERFACE_EXTENDS_LOOP = v
}

var static_TypeCheck_CONFLICTING_IMPLEMENTED_TYPE *DiagnosticType

func StaticGet_TypeCheck_CONFLICTING_IMPLEMENTED_TYPE() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_CONFLICTING_IMPLEMENTED_TYPE
}
func StaticSet_TypeCheck_CONFLICTING_IMPLEMENTED_TYPE(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_CONFLICTING_IMPLEMENTED_TYPE = v
}

var static_TypeCheck_BAD_IMPLEMENTED_TYPE *DiagnosticType

func StaticGet_TypeCheck_BAD_IMPLEMENTED_TYPE() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_BAD_IMPLEMENTED_TYPE
}
func StaticSet_TypeCheck_BAD_IMPLEMENTED_TYPE(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_BAD_IMPLEMENTED_TYPE = v
}

var static_TypeCheck_HIDDEN_SUPERCLASS_PROPERTY *DiagnosticType

func StaticGet_TypeCheck_HIDDEN_SUPERCLASS_PROPERTY() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_HIDDEN_SUPERCLASS_PROPERTY
}
func StaticSet_TypeCheck_HIDDEN_SUPERCLASS_PROPERTY(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_HIDDEN_SUPERCLASS_PROPERTY = v
}

var static_TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY *DiagnosticType

func StaticGet_TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY
}
func StaticSet_TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY = v
}

var static_TypeCheck_HIDDEN_INTERFACE_PROPERTY *DiagnosticType

func StaticGet_TypeCheck_HIDDEN_INTERFACE_PROPERTY() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_HIDDEN_INTERFACE_PROPERTY
}
func StaticSet_TypeCheck_HIDDEN_INTERFACE_PROPERTY(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_HIDDEN_INTERFACE_PROPERTY = v
}

var static_TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY_MISMATCH *DiagnosticType

func StaticGet_TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY_MISMATCH() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY_MISMATCH
}
func StaticSet_TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY_MISMATCH(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY_MISMATCH = v
}

var static_TypeCheck_UNKNOWN_OVERRIDE *DiagnosticType

func StaticGet_TypeCheck_UNKNOWN_OVERRIDE() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_UNKNOWN_OVERRIDE
}
func StaticSet_TypeCheck_UNKNOWN_OVERRIDE(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_UNKNOWN_OVERRIDE = v
}

var static_TypeCheck_UNKNOWN_PROTOTYPAL_OVERRIDE *DiagnosticType

func StaticGet_TypeCheck_UNKNOWN_PROTOTYPAL_OVERRIDE() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_UNKNOWN_PROTOTYPAL_OVERRIDE
}
func StaticSet_TypeCheck_UNKNOWN_PROTOTYPAL_OVERRIDE(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_UNKNOWN_PROTOTYPAL_OVERRIDE = v
}

var static_TypeCheck_UNKNOWN_EXPR_TYPE *DiagnosticType

func StaticGet_TypeCheck_UNKNOWN_EXPR_TYPE() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_UNKNOWN_EXPR_TYPE
}
func StaticSet_TypeCheck_UNKNOWN_EXPR_TYPE(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_UNKNOWN_EXPR_TYPE = v
}

var static_TypeCheck_WRONG_ARGUMENT_COUNT *DiagnosticType

func StaticGet_TypeCheck_WRONG_ARGUMENT_COUNT() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_WRONG_ARGUMENT_COUNT
}
func StaticSet_TypeCheck_WRONG_ARGUMENT_COUNT(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_WRONG_ARGUMENT_COUNT = v
}

var static_TypeCheck_ILLEGAL_IMPLICIT_CAST *DiagnosticType

func StaticGet_TypeCheck_ILLEGAL_IMPLICIT_CAST() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_ILLEGAL_IMPLICIT_CAST
}
func StaticSet_TypeCheck_ILLEGAL_IMPLICIT_CAST(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_ILLEGAL_IMPLICIT_CAST = v
}

var static_TypeCheck_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE *DiagnosticType

func StaticGet_TypeCheck_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE
}
func StaticSet_TypeCheck_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE = v
}

var static_TypeCheck_EXPECTED_THIS_TYPE *DiagnosticType

func StaticGet_TypeCheck_EXPECTED_THIS_TYPE() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_EXPECTED_THIS_TYPE
}
func StaticSet_TypeCheck_EXPECTED_THIS_TYPE(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_EXPECTED_THIS_TYPE = v
}

var static_TypeCheck_IN_USED_WITH_STRUCT *DiagnosticType

func StaticGet_TypeCheck_IN_USED_WITH_STRUCT() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_IN_USED_WITH_STRUCT
}
func StaticSet_TypeCheck_IN_USED_WITH_STRUCT(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_IN_USED_WITH_STRUCT = v
}

var static_TypeCheck_ILLEGAL_PROPERTY_CREATION *DiagnosticType

func StaticGet_TypeCheck_ILLEGAL_PROPERTY_CREATION() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_ILLEGAL_PROPERTY_CREATION
}
func StaticSet_TypeCheck_ILLEGAL_PROPERTY_CREATION(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_ILLEGAL_PROPERTY_CREATION = v
}

var static_TypeCheck_ILLEGAL_PROPERTY_CREATION_ON_UNION_TYPE *DiagnosticType

func StaticGet_TypeCheck_ILLEGAL_PROPERTY_CREATION_ON_UNION_TYPE() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_ILLEGAL_PROPERTY_CREATION_ON_UNION_TYPE
}
func StaticSet_TypeCheck_ILLEGAL_PROPERTY_CREATION_ON_UNION_TYPE(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_ILLEGAL_PROPERTY_CREATION_ON_UNION_TYPE = v
}

var static_TypeCheck_ILLEGAL_OBJLIT_KEY *DiagnosticType

func StaticGet_TypeCheck_ILLEGAL_OBJLIT_KEY() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_ILLEGAL_OBJLIT_KEY
}
func StaticSet_TypeCheck_ILLEGAL_OBJLIT_KEY(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_ILLEGAL_OBJLIT_KEY = v
}

var static_TypeCheck_ILLEGAL_CLASS_KEY *DiagnosticType

func StaticGet_TypeCheck_ILLEGAL_CLASS_KEY() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_ILLEGAL_CLASS_KEY
}
func StaticSet_TypeCheck_ILLEGAL_CLASS_KEY(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_ILLEGAL_CLASS_KEY = v
}

var static_TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY *DiagnosticType

func StaticGet_TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY
}
func StaticSet_TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY = v
}

var static_TypeCheck_ABSTRACT_METHOD_IN_CONCRETE_CLASS *DiagnosticType

func StaticGet_TypeCheck_ABSTRACT_METHOD_IN_CONCRETE_CLASS() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_ABSTRACT_METHOD_IN_CONCRETE_CLASS
}
func StaticSet_TypeCheck_ABSTRACT_METHOD_IN_CONCRETE_CLASS(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_ABSTRACT_METHOD_IN_CONCRETE_CLASS = v
}

var static_TypeCheck_CONFLICTING_GETTER_SETTER_TYPE *DiagnosticType

func StaticGet_TypeCheck_CONFLICTING_GETTER_SETTER_TYPE() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_CONFLICTING_GETTER_SETTER_TYPE
}
func StaticSet_TypeCheck_CONFLICTING_GETTER_SETTER_TYPE(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_CONFLICTING_GETTER_SETTER_TYPE = v
}

var static_TypeCheck_SAME_INTERFACE_MULTIPLE_IMPLEMENTS *DiagnosticType

func StaticGet_TypeCheck_SAME_INTERFACE_MULTIPLE_IMPLEMENTS() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_SAME_INTERFACE_MULTIPLE_IMPLEMENTS
}
func StaticSet_TypeCheck_SAME_INTERFACE_MULTIPLE_IMPLEMENTS(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_SAME_INTERFACE_MULTIPLE_IMPLEMENTS = v
}

var static_TypeCheck_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE *DiagnosticType

func StaticGet_TypeCheck_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE() *DiagnosticType {
	clinit_TypeCheck()
	return static_TypeCheck_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE
}
func StaticSet_TypeCheck_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE(v *DiagnosticType) {
	clinit_TypeCheck()
	static_TypeCheck_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE = v
}

var static_TypeCheck_ALL_DIAGNOSTICS *DiagnosticGroup

func StaticGet_TypeCheck_ALL_DIAGNOSTICS() *DiagnosticGroup {
	clinit_TypeCheck()
	return static_TypeCheck_ALL_DIAGNOSTICS
}
func StaticSet_TypeCheck_ALL_DIAGNOSTICS(v *DiagnosticGroup) {
	clinit_TypeCheck()
	static_TypeCheck_ALL_DIAGNOSTICS = v
}

var static_TypeCheck_ES5_INHERITANCE_DIAGNOSTIC_GROUP *DiagnosticGroup

func StaticGet_TypeCheck_ES5_INHERITANCE_DIAGNOSTIC_GROUP() *DiagnosticGroup {
	clinit_TypeCheck()
	return static_TypeCheck_ES5_INHERITANCE_DIAGNOSTIC_GROUP
}
func StaticSet_TypeCheck_ES5_INHERITANCE_DIAGNOSTIC_GROUP(v *DiagnosticGroup) {
	clinit_TypeCheck()
	static_TypeCheck_ES5_INHERITANCE_DIAGNOSTIC_GROUP = v
}

var static_TypeCheck_GOOG_INHERITS pkg_rhino.QualifiedName

func StaticGet_TypeCheck_GOOG_INHERITS() pkg_rhino.QualifiedName {
	clinit_TypeCheck()
	return static_TypeCheck_GOOG_INHERITS
}
func StaticSet_TypeCheck_GOOG_INHERITS(v pkg_rhino.QualifiedName) {
	clinit_TypeCheck()
	static_TypeCheck_GOOG_INHERITS = v
}

var clinitDone_TypeCheck uint32
var clinitState_TypeCheck j2gort.ClinitState

func clinit_TypeCheck() {
	if atomic.LoadUint32(&clinitDone_TypeCheck) != 0 {
		return
	}
	clinitSlow_TypeCheck()
}

func clinitSlow_TypeCheck() {
	if !j2gort.EnterClinit(&clinitState_TypeCheck, &clinitDone_TypeCheck) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_TypeCheck, &clinitDone_TypeCheck)
	static_TypeCheck_UNEXPECTED_TOKEN = DiagnosticType_Error("JSC_INTERNAL_ERROR_UNEXPECTED_TOKEN", "Internal Error: TypeCheck doesn''t know how to handle {0}")
	static_TypeCheck_DETERMINISTIC_TEST = DiagnosticType_Warning("JSC_DETERMINISTIC_TEST", "condition always evaluates to {2}\nleft : {0}\nright: {1}")
	static_TypeCheck_INEXISTENT_ENUM_ELEMENT = DiagnosticType_Warning("JSC_INEXISTENT_ENUM_ELEMENT", "element {0} does not exist on this enum")
	static_TypeCheck_INEXISTENT_PROPERTY = DiagnosticType_Warning("JSC_INEXISTENT_PROPERTY", "Property {0} never defined on {1}")
	static_TypeCheck_POSSIBLE_INEXISTENT_PROPERTY = DiagnosticType_Disabled("JSC_POSSIBLE_INEXISTENT_PROPERTY", ("Property {0} never defined on {1}" + j2gort.StrValString(javalang.StringReplace("\n\nThis property is accessed on a \"loose\" type, but is not defined anywhere in the program, so it must not exist.", "'", "''"))))
	static_TypeCheck_INEXISTENT_PROPERTY_WITH_SUGGESTION = DiagnosticType_Warning("JSC_INEXISTENT_PROPERTY_WITH_SUGGESTION", "Property {0} never defined on {1}. Did you mean {2}?")
	static_TypeCheck_STRICT_INEXISTENT_PROPERTY = DiagnosticType_Disabled("JSC_STRICT_INEXISTENT_PROPERTY", "Property {0} never defined on {1}")
	static_TypeCheck_STRICT_INEXISTENT_UNION_PROPERTY = DiagnosticType_Disabled("JSC_STRICT_INEXISTENT_UNION_PROPERTY", "Property {0} not defined on all member types of {1}")
	static_TypeCheck_STRICT_INEXISTENT_PROPERTY_WITH_SUGGESTION = DiagnosticType_Disabled("JSC_STRICT_INEXISTENT_PROPERTY_WITH_SUGGESTION", "Property {0} never defined on {1}. Did you mean {2}?")
	static_TypeCheck_NOT_A_CONSTRUCTOR = DiagnosticType_Warning("JSC_NOT_A_CONSTRUCTOR", "cannot instantiate non-constructor, found type: {0}")
	static_TypeCheck_INSTANTIATE_ABSTRACT_CLASS = DiagnosticType_Warning("JSC_INSTANTIATE_ABSTRACT_CLASS", "cannot instantiate abstract class")
	static_TypeCheck_BIT_OPERATION = DiagnosticType_Warning("JSC_BAD_TYPE_FOR_BIT_OPERATION", "operator {0} cannot be applied to {1}")
	static_TypeCheck_UNARY_OPERATION = DiagnosticType_Warning("JSC_BAD_TYPE_FOR_UNARY_OPERATION", "unary operator {0} cannot be applied to {1}")
	static_TypeCheck_BINARY_OPERATION = DiagnosticType_Warning("JSC_BAD_TYPES_FOR_BINARY_OPERATION", "operator {0} cannot be applied to {1} and {2}")
	static_TypeCheck_NOT_CALLABLE = DiagnosticType_Warning("JSC_NOT_FUNCTION_TYPE", "{0} expressions are not callable")
	static_TypeCheck_CONSTRUCTOR_NOT_CALLABLE = DiagnosticType_Warning("JSC_CONSTRUCTOR_NOT_CALLABLE", "Constructor {0} should be called with the \"new\" keyword")
	static_TypeCheck_ABSTRACT_SUPER_METHOD_NOT_USABLE = DiagnosticType_Warning("JSC_ABSTRACT_SUPER_METHOD_NOT_USABLE", "Abstract super method {0} cannot be dereferenced")
	static_TypeCheck_FUNCTION_MASKS_VARIABLE = DiagnosticType_Warning("JSC_FUNCTION_MASKS_VARIABLE", "function {0} masks variable (IE bug)")
	static_TypeCheck_MULTIPLE_VAR_DEF = DiagnosticType_Warning("JSC_MULTIPLE_VAR_DEF", "declaration of multiple variables with shared type information")
	static_TypeCheck_INVALID_INTERFACE_MEMBER_DECLARATION = DiagnosticType_Warning("JSC_INVALID_INTERFACE_MEMBER_DECLARATION", "interface members can only be empty property declarations, empty functions{0}")
	static_TypeCheck_INTERFACE_METHOD_NOT_EMPTY = DiagnosticType_Warning("JSC_INTERFACE_METHOD_NOT_EMPTY", "interface member functions must have an empty body")
	static_TypeCheck_CONFLICTING_EXTENDED_TYPE = DiagnosticType_Warning("JSC_CONFLICTING_EXTENDED_TYPE", "{1} cannot extend this type; {0}s can only extend {0}s")
	static_TypeCheck_ES5_CLASS_EXTENDING_ES6_CLASS = DiagnosticType_Warning("JSC_ES5_CLASS_EXTENDING_ES6_CLASS", "ES5 class {0} cannot extend ES6 class {1}")
	static_TypeCheck_DICT_EXTEND_STRUCT_TYPE = DiagnosticType_Warning("JSC_DICT_EXTEND_STRUCT_TYPE", "@dict class {0} cannot extend @struct class {1}")
	static_TypeCheck_STRUCT_EXTEND_DICT_TYPE = DiagnosticType_Warning("JSC_DICT_EXTEND_STRUCT_TYPE", "@struct class {0} cannot extend @dict class {1}")
	static_TypeCheck_ES6_CLASS_EXTENDING_CLASS_WITH_GOOG_INHERITS = DiagnosticType_Warning("JSC_ES6_CLASS_EXTENDING_CLASS_WITH_GOOG_INHERITS", "Do not use goog.inherits with ES6 classes. Use the ES6 `extends` keyword to inherit instead.")
	static_TypeCheck_INTERFACE_EXTENDS_LOOP = DiagnosticType_Warning("JSC_INTERFACE_EXTENDS_LOOP", "extends loop involving {0}, loop: {1}")
	static_TypeCheck_CONFLICTING_IMPLEMENTED_TYPE = DiagnosticType_Warning("JSC_CONFLICTING_IMPLEMENTED_TYPE", "{0} cannot implement this type; an interface can only extend, but not implement interfaces")
	static_TypeCheck_BAD_IMPLEMENTED_TYPE = DiagnosticType_Warning("JSC_IMPLEMENTS_NON_INTERFACE", "can only implement interfaces")
	static_TypeCheck_HIDDEN_SUPERCLASS_PROPERTY = DiagnosticType_Disabled("JSC_HIDDEN_SUPERCLASS_PROPERTY", "property {0} already defined on superclass {1}; use @override to override it")
	static_TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY = DiagnosticType_Disabled("JSC_PROTOTYPAL_HIDDEN_SUPERCLASS_PROPERTY", "property {0} already defined on supertype {1}; use @override to override it")
	static_TypeCheck_HIDDEN_INTERFACE_PROPERTY = DiagnosticType_Disabled("JSC_HIDDEN_INTERFACE_PROPERTY", "property {0} already defined on interface {1}; use @override to override it")
	static_TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY_MISMATCH = DiagnosticType_Warning("JSC_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY_MISMATCH", "mismatch of the {0} property type and the type of the property it overrides from supertype {1}\noriginal: {2}\noverride: {3}")
	static_TypeCheck_UNKNOWN_OVERRIDE = DiagnosticType_Warning("JSC_UNKNOWN_OVERRIDE", "property {0} not defined on any superclass of {1}")
	static_TypeCheck_UNKNOWN_PROTOTYPAL_OVERRIDE = DiagnosticType_Warning("JSC_UNKNOWN_PROTOTYPAL_OVERRIDE", "property {0} not defined on any supertype of {1}")
	static_TypeCheck_UNKNOWN_EXPR_TYPE = DiagnosticType_Warning("JSC_UNKNOWN_EXPR_TYPE", "could not determine the type of this expression")
	static_TypeCheck_WRONG_ARGUMENT_COUNT = DiagnosticType_Warning("JSC_WRONG_ARGUMENT_COUNT", "Function {0}: called with {1} argument(s). Function requires at least {2} argument(s){3}.")
	static_TypeCheck_ILLEGAL_IMPLICIT_CAST = DiagnosticType_Warning("JSC_ILLEGAL_IMPLICIT_CAST", "Illegal annotation on {0}. @implicitCast may only be used in externs.")
	static_TypeCheck_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE = DiagnosticType_Warning("JSC_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE", "Interface {0} has a property {1} with incompatible types in its super interfaces {2} and {3}")
	static_TypeCheck_EXPECTED_THIS_TYPE = DiagnosticType_Warning("JSC_EXPECTED_THIS_TYPE", "\"{0}\" must be called with a \"this\" type")
	static_TypeCheck_IN_USED_WITH_STRUCT = DiagnosticType_Warning("JSC_IN_USED_WITH_STRUCT", "Cannot use the IN operator with structs")
	static_TypeCheck_ILLEGAL_PROPERTY_CREATION = DiagnosticType_Warning("JSC_ILLEGAL_PROPERTY_CREATION", "Cannot add a property to a struct instance after it is constructed. (If you already declared the property, make sure to give it a type.)")
	static_TypeCheck_ILLEGAL_PROPERTY_CREATION_ON_UNION_TYPE = DiagnosticType_Disabled("JSC_ILLEGAL_PROPERTY_CREATION_ON_UNION_TYPE", "Cannot add a property to an instance of union type.")
	static_TypeCheck_ILLEGAL_OBJLIT_KEY = DiagnosticType_Warning("JSC_ILLEGAL_OBJLIT_KEY", "Illegal key, the object literal is a {0}")
	static_TypeCheck_ILLEGAL_CLASS_KEY = DiagnosticType_Warning("JSC_ILLEGAL_CLASS_KEY", "Illegal key, the class is a {0}")
	static_TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY = DiagnosticType_Warning("JSC_NON_STRINGIFIABLE_OBJECT_KEY", "Object type \"{0}\" contains non-stringifiable key and it may lead to an error. Please use ES6 Map instead or implement your own Map structure.")
	static_TypeCheck_ABSTRACT_METHOD_IN_CONCRETE_CLASS = DiagnosticType_Warning("JSC_ABSTRACT_METHOD_IN_CONCRETE_CLASS", "Abstract methods can only appear in abstract classes. Please declare the class as @abstract")
	static_TypeCheck_CONFLICTING_GETTER_SETTER_TYPE = DiagnosticType_Warning("JSC_CONFLICTING_GETTER_SETTER_TYPE", "The types of the getter and setter for property ''{0}'' do not match.\ngetter type is: {1}\nsetter type is: {2}")
	static_TypeCheck_SAME_INTERFACE_MULTIPLE_IMPLEMENTS = DiagnosticType_Warning("JSC_SAME_INTERFACE_MULTIPLE_IMPLEMENTS", "Cannot @implement the same interface more than once\nRepeated interface: {0}")
	static_TypeCheck_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE = DiagnosticType_Error("JSC_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE", "Should not assign to a property of readonly type ''{0}''")
	static_TypeCheck_ALL_DIAGNOSTICS = NewDiagnosticGroup_1_DiagnosticType_2([]*DiagnosticType{StaticGet_TypeCheck_ABSTRACT_METHOD_IN_CONCRETE_CLASS(), StaticGet_TypeCheck_ABSTRACT_SUPER_METHOD_NOT_USABLE(), StaticGet_TypeCheck_BAD_IMPLEMENTED_TYPE(), StaticGet_TypeCheck_BINARY_OPERATION(), StaticGet_TypeCheck_BIT_OPERATION(), StaticGet_TypeCheck_CONFLICTING_EXTENDED_TYPE(), StaticGet_TypeCheck_CONFLICTING_GETTER_SETTER_TYPE(), StaticGet_TypeCheck_CONFLICTING_IMPLEMENTED_TYPE(), StaticGet_TypeCheck_CONSTRUCTOR_NOT_CALLABLE(), StaticGet_TypeCheck_DETERMINISTIC_TEST(), StaticGet_TypeCheck_ES5_CLASS_EXTENDING_ES6_CLASS(), StaticGet_TypeCheck_EXPECTED_THIS_TYPE(), StaticGet_TypeCheck_FUNCTION_MASKS_VARIABLE(), StaticGet_TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY_MISMATCH(), StaticGet_TypeCheck_ILLEGAL_CLASS_KEY(), StaticGet_TypeCheck_ILLEGAL_IMPLICIT_CAST(), StaticGet_TypeCheck_ILLEGAL_OBJLIT_KEY(), StaticGet_TypeCheck_ILLEGAL_PROPERTY_CREATION(), StaticGet_TypeCheck_ILLEGAL_PROPERTY_CREATION_ON_UNION_TYPE(), StaticGet_TypeCheck_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE(), StaticGet_TypeCheck_INEXISTENT_ENUM_ELEMENT(), StaticGet_TypeCheck_INEXISTENT_PROPERTY(), StaticGet_TypeCheck_INEXISTENT_PROPERTY_WITH_SUGGESTION(), StaticGet_TypeCheck_INSTANTIATE_ABSTRACT_CLASS(), StaticGet_TypeCheck_INTERFACE_METHOD_NOT_EMPTY(), StaticGet_TypeCheck_INVALID_INTERFACE_MEMBER_DECLARATION(), StaticGet_TypeCheck_IN_USED_WITH_STRUCT(), StaticGet_TypeCheck_MULTIPLE_VAR_DEF(), StaticGet_TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY(), StaticGet_TypeCheck_NOT_A_CONSTRUCTOR(), StaticGet_TypeCheck_NOT_CALLABLE(), StaticGet_TypeCheck_STRUCT_EXTEND_DICT_TYPE(), StaticGet_TypeCheck_DICT_EXTEND_STRUCT_TYPE(), StaticGet_TypeCheck_POSSIBLE_INEXISTENT_PROPERTY(), StaticGet_TypeCheck_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE(), StaticGet_RhinoErrorReporter_CYCLIC_INHERITANCE_ERROR(), StaticGet_RhinoErrorReporter_TOO_MANY_TEMPLATE_PARAMS(), StaticGet_RhinoErrorReporter_TYPE_PARSE_ERROR(), StaticGet_RhinoErrorReporter_UNRECOGNIZED_TYPE_ERROR(), StaticGet_TypeCheck_SAME_INTERFACE_MULTIPLE_IMPLEMENTS(), StaticGet_TypeValidator_HIDDEN_SUPERCLASS_PROPERTY_MISMATCH(), StaticGet_TypedScopeCreator_CTOR_INITIALIZER(), StaticGet_TypedScopeCreator_IFACE_INITIALIZER(), StaticGet_TypedScopeCreator_LENDS_ON_NON_OBJECT(), StaticGet_TypedScopeCreator_UNKNOWN_LENDS(), StaticGet_TypeCheck_UNARY_OPERATION(), StaticGet_TypeCheck_UNKNOWN_OVERRIDE(), StaticGet_TypeCheck_UNKNOWN_PROTOTYPAL_OVERRIDE(), StaticGet_TypeCheck_WRONG_ARGUMENT_COUNT()})
	static_TypeCheck_ES5_INHERITANCE_DIAGNOSTIC_GROUP = NewDiagnosticGroup_1_DiagnosticType_2([]*DiagnosticType{StaticGet_TypeCheck_ES5_CLASS_EXTENDING_ES6_CLASS()})
	static_TypeCheck_GOOG_INHERITS = pkg_rhino.QualifiedName_Of("goog.inherits")
}

func (this_ *TypeCheck) initFields_TypeCheck() {
	_ = this_

	this_.subtypingMode = pkg_rhino.StaticGet_JSType_SubtypingMode_NORMAL()
	this_.reportMissingProperties = true

	this_.nonStringifiableKeyVisited = j2gort.NewCustomOrderedSet[pkg_rhino.JSType]()
}

// Java constructor: TypeCheck(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.type.ReverseAbstractInterpreter reverseInterpreter, com.google.javascript.rhino.jstype.JSTypeRegistry typeRegistry, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable TypedScope topScope, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable TypedScopeCreator scopeCreator)
func NewTypeCheck_5(compiler AbstractCompiler, reverseInterpreter pkg_jscomptype.ReverseAbstractInterpreter, typeRegistry *pkg_rhino.JSTypeRegistry, topScope *TypedScope, scopeCreator *TypedScopeCreator) *TypeCheck {
	clinit_TypeCheck()
	this_ := &TypeCheck{}
	InitTypeCheck_5(this_, compiler, reverseInterpreter, typeRegistry, topScope, scopeCreator)
	return this_
}

func InitTypeCheck_5(this_ *TypeCheck, compiler AbstractCompiler, reverseInterpreter pkg_jscomptype.ReverseAbstractInterpreter, typeRegistry *pkg_rhino.JSTypeRegistry, topScope *TypedScope, scopeCreator *TypedScopeCreator) {
	_ = compiler
	_ = reverseInterpreter
	_ = typeRegistry
	_ = topScope
	_ = scopeCreator
	this_.initFields_TypeCheck()
	this_.compiler = compiler
	this_.validator = compiler.GetTypeValidator()
	this_.reverseInterpreter = reverseInterpreter
	this_.typeRegistry = typeRegistry
	this_.topScope = topScope
	this_.scopeCreator = scopeCreator
	this_.inferJSDocInfo = NewInferJSDocInfo(compiler)
}

// Java constructor: TypeCheck(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.type.ReverseAbstractInterpreter reverseInterpreter, com.google.javascript.rhino.jstype.JSTypeRegistry typeRegistry)
func NewTypeCheck_3(compiler AbstractCompiler, reverseInterpreter pkg_jscomptype.ReverseAbstractInterpreter, typeRegistry *pkg_rhino.JSTypeRegistry) *TypeCheck {
	clinit_TypeCheck()
	this_ := &TypeCheck{}
	InitTypeCheck_3(this_, compiler, reverseInterpreter, typeRegistry)
	return this_
}

func InitTypeCheck_3(this_ *TypeCheck, compiler AbstractCompiler, reverseInterpreter pkg_jscomptype.ReverseAbstractInterpreter, typeRegistry *pkg_rhino.JSTypeRegistry) {
	_ = compiler
	_ = reverseInterpreter
	_ = typeRegistry
	InitTypeCheck_5(this_, compiler, reverseInterpreter, typeRegistry, nil, nil)
}

// Java method: TypeCheck.reportMissingProperties(boolean report)
func (this_ *TypeCheck) ReportMissingProperties(report bool) *TypeCheck {
	_ = report
	this_.reportMissingProperties = report
	return this_
}

// Java method: TypeCheck.reportUnknownTypes(boolean report)
func (this_ *TypeCheck) ReportUnknownTypes(report bool) *TypeCheck {
	_ = report
	this_.reportUnknownTypes = report
	return this_
}

// Java method: TypeCheck.process(com.google.javascript.rhino.Node externsRoot, com.google.javascript.rhino.Node jsRoot)
func (this_ *TypeCheck) Process_2(externsRoot pkg_rhino.Node, jsRoot pkg_rhino.Node) {
	_ = externsRoot
	_ = jsRoot
	j2gort.CheckNotNullPtr(this_.scopeCreator)
	j2gort.CheckNotNullPtr(this_.topScope)
	var externsAndJs pkg_rhino.Node = jsRoot.GetParent()
	_ = externsAndJs
	if !(externsAndJs != 0) {
		j2gort.CheckState(false)
	}
	if !((externsRoot == 0) || externsAndJs.HasChild(externsRoot)) {
		j2gort.CheckState(false)
	}
	func() {
		var debugTypeLogger *TypeCheck_DebugTypeLogger = NewTypeCheck_DebugTypeLogger(this_)
		_ = debugTypeLogger
		defer j2gort.CloseResource(debugTypeLogger)
		this_.debugTypeLogger = debugTypeLogger
		if externsRoot != 0 {
			this_.Check(externsRoot, true)
		}
		this_.Check(jsRoot, false)
	}()
}

// Java method: TypeCheck.processForTesting(com.google.javascript.rhino.Node externsRoot, com.google.javascript.rhino.Node jsRoot)
func (this_ *TypeCheck) ProcessForTesting(externsRoot pkg_rhino.Node, jsRoot pkg_rhino.Node) *TypedScope {
	_ = externsRoot
	_ = jsRoot
	if !(this_.scopeCreator == nil) {
		j2gort.CheckState(false)
	}
	if !(this_.topScope == nil) {
		j2gort.CheckState(false)
	}
	if !((externsRoot == 0) || externsRoot.IsRoot()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(externsRoot))
	}
	if !(jsRoot.IsRoot()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(jsRoot))
	}
	if !(jsRoot.HasParent() && jsRoot.GetParent().IsRoot()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(jsRoot.GetParent()))
	}
	if !((externsRoot == 0) || (externsRoot.GetNext() == jsRoot)) {
		j2gort.CheckStateStr(false, "externs root must be the preceding sibling of the js root")
	}
	this_.scopeCreator = NewTypedScopeCreator_1(this_.compiler)
	this_.topScope = NewTypeInferencePass(this_.compiler, this_.reverseInterpreter, this_.scopeCreator).InferAllScopes(jsRoot.GetParent())
	this_.Process_2(externsRoot, jsRoot)
	this_.compiler.SetTypeCheckingHasRun(true)
	this_.compiler.SetTopScope(this_.topScope)
	return this_.topScope
}

// Java method: TypeCheck.check(com.google.javascript.rhino.Node node, boolean externs)
func (this_ *TypeCheck) Check(node pkg_rhino.Node, externs bool) {
	_ = node
	_ = externs
	j2gort.CheckNotNullNode(node)
	this_.inExterns = externs
	{
		clinit_NodeTraversal()
		var ntb_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_)
		(&ntb_).SetCompiler(this_.compiler)
		(&ntb_).SetCallback_1_NodeTraversal_Callback(this_)
		(&ntb_).SetScopeCreator(ToIface_ScopeCreator(this_.scopeCreator))
		(&ntb_).TraverseWithScope(node, ToIface_AbstractScope(this_.topScope))
	}
	if externs {
		this_.inferJSDocInfo.Process_2(node, pkg_rhino.Node(0))
	} else {
		this_.inferJSDocInfo.Process_2(pkg_rhino.Node(0), node)
	}
}

// Java method: TypeCheck.report(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.DiagnosticType diagnosticType, java.lang.String[] arguments)
func (this_ *TypeCheck) priv_TypeCheck_report_5386(n pkg_rhino.Node, diagnosticType *DiagnosticType, arguments []string) {
	_ = n
	_ = diagnosticType
	_ = arguments
	this_.compiler.Report_1(JSError_Make_3(n, diagnosticType, arguments))
}

// Java method: TypeCheck.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *TypeCheck) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	this_.debugTypeLogger.MaybeStartLoggingAt(n)
	if n.IsScript() {
		if NodeUtil_IsFromTypeSummary(n) {
			return false
		}
		if J2clSourceUtils_IsJ2clSource_1_Node(n) {
			this_.subtypingMode = pkg_rhino.StaticGet_JSType_SubtypingMode_IGNORE_NULL_UNDEFINED()
		} else {
			this_.subtypingMode = pkg_rhino.StaticGet_JSType_SubtypingMode_NORMAL()
		}
		this_.validator.SetSubtypingMode(this_.subtypingMode)
	} else if n.IsFunction() {
		var outerScope *TypedScope = t.GetTypedScope()
		_ = outerScope
		var var_ *TypedVar = outerScope.LookupTypedVar(n.GetFirstChild().GetString())
		_ = var_
		if (((var_ != nil) && j2gort.CastPtr[TypedScope](var_.GetScope()).HasSameContainerScope(ToIface_AbstractScope(outerScope))) && !(j2gort.InstanceOf[*pkg_rhino.FunctionType](var_.GetType()))) && !(TypeValidator_HasDuplicateDeclarationSuppression(this_.compiler, var_.GetNameNode())) {
			this_.priv_TypeCheck_report_5386(n, StaticGet_TypeCheck_FUNCTION_MASKS_VARIABLE(), []string{var_.GetName()})
		}
	}
	return true
}

// Java method: TypeCheck.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *TypeCheck) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var childType pkg_rhino.JSType = nil
	_ = childType
	var leftType pkg_rhino.JSType = nil
	_ = leftType
	var rightType pkg_rhino.JSType = nil
	_ = rightType
	var left pkg_rhino.Node = pkg_rhino.Node(0)
	_ = left
	var right pkg_rhino.Node = pkg_rhino.Node(0)
	_ = right
	var typeable bool = true
	_ = typeable
	this_.validator.ExpectWellFormedTemplatizedType(n)
	switch n.TokenSwitchOrd_() {
	case int32(101) /* Token.CAST */ :
		var expr pkg_rhino.Node = n.GetFirstChild()
		_ = expr
		var exprType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(expr)
		_ = exprType
		var castType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(n)
		_ = castType
		if !(expr.IsObjectLit()) {
			this_.validator.ExpectCanCast(n, castType, exprType)
		}
		this_.priv_TypeCheck_ensureTyped_5461(n, castType)
		expr.SetJSTypeBeforeCast(exprType)
		if castType.RestrictByNotNullOrUndefined().IsSubtypeOf_1(exprType) || expr.IsObjectLit() {
			expr.SetJSType(castType)
		}
	case int32(32) /* Token.NAME */ :
		typeable = this_.VisitName(t, n, parent)
	case int32(50) /* Token.COMMA */ :
		this_.priv_TypeCheck_ensureTyped_5461(n, this_.priv_TypeCheck_getJSType_5456(n.GetLastChild()))
	case int32(37) /* Token.THIS */ :
		this_.priv_TypeCheck_ensureTyped_5461(n, t.GetTypedScope().GetTypeOfThis())
	case int32(36) /* Token.NULL */ :
		this_.priv_TypeCheck_ensureTyped_5460(n, pkg_rhino.StaticGet_JSTypeNative_NULL_TYPE())
	case int32(33) /* Token.NUMBER */ :
		this_.priv_TypeCheck_ensureTyped_5460(n, pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE())
	case int32(34) /* Token.BIGINT */ :
		this_.priv_TypeCheck_ensureTyped_5460(n, pkg_rhino.StaticGet_JSTypeNative_BIGINT_TYPE())
	case int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */ :
	case int32(46) /* Token.ARRAYLIT */ :
		this_.priv_TypeCheck_ensureTyped_5460(n, pkg_rhino.StaticGet_JSTypeNative_ARRAY_TYPE())
	case int32(42) /* Token.REGEXP */ :
		this_.priv_TypeCheck_ensureTyped_5460(n, pkg_rhino.StaticGet_JSTypeNative_REGEXP_TYPE())
	case int32(26) /* Token.GETPROP */ :
		this_.priv_TypeCheck_visitGetProp_5419(t, n)
		typeable = !(parent.IsAssign() && (parent.GetFirstChild() == n))
	case int32(29) /* Token.OPTCHAIN_GETPROP */ :
		this_.priv_TypeCheck_visitOptChainGetProp_5420(n)
	case int32(30) /* Token.OPTCHAIN_GETELEM */ :
		this_.priv_TypeCheck_visitOptChainGetElem_5432(n)
	case int32(27) /* Token.GETELEM */ :
		this_.priv_TypeCheck_visitGetElem_5431(n)
		typeable = false
	case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */ :
		this_.priv_TypeCheck_visitVar_5433(t, n)
		typeable = false
	case int32(23) /* Token.NEW */ :
		this_.priv_TypeCheck_visitNew_5434(n)
	case int32(31) /* Token.OPTCHAIN_CALL */ :
		this_.priv_TypeCheck_visitCall_5445(t, n)
	case int32(28) /* Token.CALL */ :
		this_.priv_TypeCheck_visitCall_5445(t, n)
		typeable = !(parent.IsExprResult())
	case int32(0) /* Token.RETURN */ :
		this_.priv_TypeCheck_visitReturn_5449(t, n)
		typeable = false
	case int32(114) /* Token.YIELD */ :
		this_.priv_TypeCheck_visitYield_5450(t, n)
	case int32(72) /* Token.DEC */, int32(71) /* Token.INC */ :
		left = n.GetFirstChild()
		this_.priv_TypeCheck_checkPropCreation_5399(left)
		if this_.priv_TypeCheck_getJSType_5456(n).IsNumber() {
			this_.validator.ExpectNumber(left, this_.priv_TypeCheck_getJSType_5456(left), "increment/decrement")
			this_.priv_TypeCheck_ensureTyped_5460(n, pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE())
		} else {
			this_.validator.ExpectBigIntOrNumber(left, this_.priv_TypeCheck_getJSType_5456(left), "increment/decrement")
		}
	case int32(87) /* Token.VOID */ :
		this_.priv_TypeCheck_ensureTyped_5460(n, pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE())
	case int32(35) /* Token.STRINGLIT */, int32(25) /* Token.TYPEOF */, int32(131) /* Token.TEMPLATELIT */, int32(133) /* Token.TEMPLATELIT_STRING */ :
		this_.priv_TypeCheck_ensureTyped_5460(n, pkg_rhino.StaticGet_JSTypeNative_STRING_TYPE())
	case int32(130) /* Token.TAGGED_TEMPLATELIT */ :
		this_.priv_TypeCheck_visitTaggedTemplateLit_5451(n)
		this_.priv_TypeCheck_ensureTyped_5459(n)
	case int32(20) /* Token.BITNOT */ :
		this_.priv_TypeCheck_visitBitwiseNOT_5389(n)
	case int32(21) /* Token.POS */ :
		this_.priv_TypeCheck_visitUnaryPlus_5387(n)
	case int32(22) /* Token.NEG */ :
		this_.priv_TypeCheck_visitUnaryMinus_5388(n)
	case int32(4) /* Token.EQ */, int32(5) /* Token.NE */, int32(40) /* Token.SHEQ */, int32(41) /* Token.SHNE */ :
		left = n.GetFirstChild()
		right = n.GetLastChild()
		if left.IsTypeOf() {
			if right.IsStringLit() {
				this_.priv_TypeCheck_checkTypeofString_5391(right, right.GetString())
			}
		} else if right.IsTypeOf() && left.IsStringLit() {
			this_.priv_TypeCheck_checkTypeofString_5391(left, left.GetString())
		}
		leftType = this_.priv_TypeCheck_getJSType_5456(left)
		rightType = this_.priv_TypeCheck_getJSType_5456(right)
		var leftTypeRestricted pkg_rhino.JSType = leftType.RestrictByNotNullOrUndefined()
		_ = leftTypeRestricted
		var rightTypeRestricted pkg_rhino.JSType = rightType.RestrictByNotNullOrUndefined()
		_ = rightTypeRestricted
		var result pkg_base.Tri = pkg_base.StaticGet_Tri_UNKNOWN()
		_ = result
		if n.IsEQ() || n.IsNE() {
			result = leftTypeRestricted.TestForEquality(rightTypeRestricted)
			if n.IsNE() {
				result = result.Not()
			}
		} else {
			if !(leftTypeRestricted.CanTestForShallowEqualityWith(rightTypeRestricted)) {
				if n.IsSHEQ() {
					result = pkg_base.StaticGet_Tri_FALSE()
				} else {
					result = pkg_base.StaticGet_Tri_TRUE()
				}
			}
		}
		if result != pkg_base.StaticGet_Tri_UNKNOWN() {
			this_.priv_TypeCheck_report_5386(n, StaticGet_TypeCheck_DETERMINISTIC_TEST(), []string{leftType.String(), rightType.String(), result.String()})
		}
		this_.priv_TypeCheck_ensureTyped_5460(n, pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_TYPE())
	case int32(6) /* Token.LT */, int32(7) /* Token.LE */, int32(8) /* Token.GT */, int32(9) /* Token.GE */ :
		var leftSide pkg_rhino.Node = n.GetFirstChild()
		_ = leftSide
		var rightSide pkg_rhino.Node = n.GetLastChild()
		_ = rightSide
		leftType = this_.priv_TypeCheck_getJSType_5456(leftSide)
		rightType = this_.priv_TypeCheck_getJSType_5456(rightSide)
		if rightType.IsUnknownType() {
			this_.validator.ExpectUnknownOrComparable(leftSide, leftType, "left side of comparison")
		} else if leftType.IsUnknownType() {
			this_.validator.ExpectUnknownOrComparable(rightSide, rightType, "right side of comparison")
		} else if rightType.IsBigIntOrNumber() {
			this_.validator.ExpectBigIntOrNumber(leftSide, leftType, "left side of numeric comparison")
		} else if leftType.IsBigIntOrNumber() {
			this_.validator.ExpectBigIntOrNumber(rightSide, rightType, "right side of numeric comparison")
		} else {
			var errorMsg string = "expected matching types in comparison"
			_ = errorMsg
			this_.validator.ExpectMatchingTypesStrict(n, leftType, rightType, errorMsg)
			if !(leftType.MatchesNumberContext()) || !(rightType.MatchesNumberContext()) {
				var message string = "left side of comparison"
				_ = message
				this_.validator.ExpectString(leftSide, leftType, message)
				this_.validator.ExpectNotNullOrUndefined(t, leftSide, leftType, message, this_.priv_TypeCheck_getNativeType_5463(pkg_rhino.StaticGet_JSTypeNative_STRING_TYPE()))
				message = "right side of comparison"
				this_.validator.ExpectString(rightSide, rightType, message)
				this_.validator.ExpectNotNullOrUndefined(t, rightSide, rightType, message, this_.priv_TypeCheck_getNativeType_5463(pkg_rhino.StaticGet_JSTypeNative_STRING_TYPE()))
			}
		}
		this_.priv_TypeCheck_ensureTyped_5460(n, pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_TYPE())
	case int32(44) /* Token.IN */ :
		left = n.GetFirstChild()
		right = n.GetLastChild()
		rightType = this_.priv_TypeCheck_getJSType_5456(right)
		this_.validator.ExpectStringOrSymbol(left, this_.priv_TypeCheck_getJSType_5456(left), "left side of 'in'")
		this_.validator.ExpectObject(n, rightType, "'in' requires an object")
		if rightType.IsStruct() {
			this_.priv_TypeCheck_report_5386(right, StaticGet_TypeCheck_IN_USED_WITH_STRUCT(), []string{})
		}
		this_.priv_TypeCheck_ensureTyped_5460(n, pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_TYPE())
	case int32(45) /* Token.INSTANCEOF */ :
		left = n.GetFirstChild()
		right = n.GetLastChild()
		rightType = this_.priv_TypeCheck_getJSType_5456(right).RestrictByNotNullOrUndefined()
		this_.validator.ExpectAnyObject(left, this_.priv_TypeCheck_getJSType_5456(left), "deterministic instanceof yields false")
		this_.validator.ExpectActualObject(right, rightType, "instanceof requires an object")
		this_.priv_TypeCheck_ensureTyped_5460(n, pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_TYPE())
	case int32(51) /* Token.ASSIGN */, int32(64) /* Token.ASSIGN_OR */, int32(65) /* Token.ASSIGN_AND */, int32(66) /* Token.ASSIGN_COALESCE */ :
		this_.priv_TypeCheck_visitAssign_5395(t, n)
		typeable = false
	case int32(55) /* Token.ASSIGN_LSH */, int32(56) /* Token.ASSIGN_RSH */, int32(57) /* Token.ASSIGN_URSH */, int32(61) /* Token.ASSIGN_DIV */, int32(62) /* Token.ASSIGN_MOD */, int32(52) /* Token.ASSIGN_BITOR */, int32(53) /* Token.ASSIGN_BITXOR */, int32(54) /* Token.ASSIGN_BITAND */, int32(59) /* Token.ASSIGN_SUB */, int32(58) /* Token.ASSIGN_ADD */, int32(60) /* Token.ASSIGN_MUL */, int32(63) /* Token.ASSIGN_EXPONENT */ :
		this_.priv_TypeCheck_checkPropCreation_5399(n.GetFirstChild())
		this_.priv_TypeCheck_visitBinaryOperator_5452(n.GetToken(), n)
	case int32(10) /* Token.LSH */, int32(11) /* Token.RSH */, int32(12) /* Token.URSH */, int32(16) /* Token.DIV */, int32(17) /* Token.MOD */, int32(1) /* Token.BITOR */, int32(2) /* Token.BITXOR */, int32(3) /* Token.BITAND */, int32(14) /* Token.SUB */, int32(13) /* Token.ADD */, int32(15) /* Token.MUL */, int32(18) /* Token.EXPONENT */ :
		this_.priv_TypeCheck_visitBinaryOperator_5452(n.GetToken(), n)
	case int32(39) /* Token.TRUE */, int32(38) /* Token.FALSE */, int32(19) /* Token.NOT */, int32(24) /* Token.DELPROP */ :
		this_.priv_TypeCheck_ensureTyped_5460(n, pkg_rhino.StaticGet_JSTypeNative_BOOLEAN_TYPE())
	case int32(76) /* Token.CASE */ :
		var switchConditionType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(parent.GetPrevious())
		_ = switchConditionType
		var caseType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(n.GetFirstChild())
		_ = caseType
		this_.validator.ExpectSwitchMatchesCase(n, switchConditionType, caseType)
		typeable = false
	case int32(85) /* Token.WITH */ :
		var child pkg_rhino.Node = n.GetFirstChild()
		_ = child
		childType = this_.priv_TypeCheck_getJSType_5456(child)
		this_.validator.ExpectObject(child, childType, "with requires an object")
		typeable = false
	case int32(73) /* Token.FUNCTION */ :
		this_.priv_TypeCheck_visitFunction_5437(n)
	case int32(105) /* Token.CLASS */ :
		this_.priv_TypeCheck_visitClass_5438(n)
	case int32(123) /* Token.MODULE_BODY */ :
		this_.priv_TypeCheck_visitModuleBody_5453(t, n)
	case int32(49) /* Token.PARAM_LIST */, int32(100) /* Token.STRING_KEY */, int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(129) /* Token.COMPUTED_PROP */, int32(108) /* Token.MEMBER_FIELD_DEF */, int32(109) /* Token.COMPUTED_FIELD_DEF */, int32(92) /* Token.LABEL */, int32(99) /* Token.LABEL_NAME */, int32(75) /* Token.SWITCH */, int32(91) /* Token.SWITCH_BODY */, int32(82) /* Token.BREAK */, int32(86) /* Token.CATCH */, int32(48) /* Token.TRY */, int32(94) /* Token.SCRIPT */, int32(120) /* Token.EXPORT */, int32(122) /* Token.EXPORT_SPEC */, int32(121) /* Token.EXPORT_SPECS */, int32(116) /* Token.IMPORT */, int32(118) /* Token.IMPORT_SPEC */, int32(117) /* Token.IMPORT_SPECS */, int32(119) /* Token.IMPORT_STAR */, int32(93) /* Token.EXPR_RESULT */, int32(90) /* Token.BLOCK */, int32(89) /* Token.ROOT */, int32(88) /* Token.EMPTY */, int32(77) /* Token.DEFAULT_CASE */, int32(83) /* Token.CONTINUE */, int32(98) /* Token.DEBUGGER */, int32(43) /* Token.THROW */, int32(79) /* Token.DO */, int32(74) /* Token.IF */, int32(78) /* Token.WHILE */, int32(80) /* Token.FOR */, int32(132) /* Token.TEMPLATELIT_SUB */, int32(125) /* Token.ITER_REST */, int32(126) /* Token.OBJECT_REST */, int32(104) /* Token.DESTRUCTURING_LHS */ :
		typeable = false
	case int32(124) /* Token.DYNAMIC_IMPORT */ :
		this_.priv_TypeCheck_visitDynamicImport_5454(t, n)
	case int32(102) /* Token.ARRAY_PATTERN */ :
		this_.priv_TypeCheck_ensureTyped_5459(n)
		this_.validator.ExpectAutoboxesToIterable(n, this_.priv_TypeCheck_getJSType_5456(n), "array pattern destructuring requires an Iterable")
	case int32(103) /* Token.OBJECT_PATTERN */ :
		this_.priv_TypeCheck_visitObjectPattern_5408(n)
	case int32(134) /* Token.DEFAULT_VALUE */ :
		this_.priv_TypeCheck_checkCanAssignToWithScope_5396(t, n, n.GetFirstChild(), this_.priv_TypeCheck_getJSType_5456(n.GetSecondChild()), nil, "default value has wrong type")
		var lhs pkg_rhino.Node = n.GetFirstChild()
		_ = lhs
		var rhs pkg_rhino.Node = n.GetSecondChild()
		_ = rhs
		if lhs.IsArrayPattern() {
			this_.validator.ExpectAutoboxesToIterable(rhs, this_.priv_TypeCheck_getJSType_5456(rhs), "array pattern destructuring requires an Iterable")
		} else if lhs.IsObjectPattern() {
			this_.validator.ExpectObject(rhs, this_.priv_TypeCheck_getJSType_5456(rhs), "cannot destructure a 'null' or 'undefined' default value")
		}
		typeable = false
	case int32(106) /* Token.CLASS_MEMBERS */ :
		var typ pkg_rhino.JSType = j2gort.AsIface[pkg_rhino.JSType](parent.GetJSType().ToMaybeFunctionType().GetInstanceType())
		_ = typ
		{
			var child_2 pkg_rhino.Node = n.GetFirstChild()
			_ = child_2
			for ; child_2 != 0; child_2 = child_2.GetNext() {
				this_.priv_TypeCheck_visitObjectOrClassLiteralKey_5409(child_2, n.GetParent(), typ)
				if child_2.IsSetterDef() || child_2.IsGetterDef() {
					this_.priv_TypeCheck_checkGetterOrSetterType_5393(child_2, parent.GetJSType().ToMaybeFunctionType())
				}
			}
		}
		typeable = false
	case int32(81) /* Token.FOR_IN */ :
		var obj pkg_rhino.Node = n.GetSecondChild()
		_ = obj
		if this_.priv_TypeCheck_getJSType_5456(obj).IsStruct() {
			this_.priv_TypeCheck_report_5386(obj, StaticGet_TypeCheck_IN_USED_WITH_STRUCT(), []string{})
		}
		typeable = false
	case int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */ :
		this_.priv_TypeCheck_ensureTyped_5459(n.GetSecondChild())
		typeable = false
	case int32(110) /* Token.SUPER */, int32(135) /* Token.NEW_TARGET */, int32(136) /* Token.IMPORT_META */, int32(115) /* Token.AWAIT */, int32(69) /* Token.AND */, int32(67) /* Token.HOOK */, int32(68) /* Token.OR */, int32(70) /* Token.COALESCE */ :
		this_.priv_TypeCheck_ensureTyped_5459(n)
	case int32(47) /* Token.OBJECTLIT */ :
		if j2gort.InstanceOf[*pkg_rhino.EnumType](parent.GetJSType()) {
			this_.priv_TypeCheck_ensureTyped_5461(n, parent.GetJSType())
		} else {
			this_.priv_TypeCheck_ensureTyped_5459(n)
		}
		var typ_2 pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(n)
		_ = typ_2
		{
			var key pkg_rhino.Node = n.GetFirstChild()
			_ = key
			for ; key != 0; key = key.GetNext() {
				this_.priv_TypeCheck_visitObjectOrClassLiteralKey_5409(key, n, typ_2)
			}
		}
	case int32(127) /* Token.ITER_SPREAD */, int32(128) /* Token.OBJECT_SPREAD */ :
		this_.priv_TypeCheck_checkSpread_5390(n)
		typeable = false
	default:
		this_.priv_TypeCheck_report_5386(n, StaticGet_TypeCheck_UNEXPECTED_TOKEN(), []string{j2gort.CheckNotNullEnum(n.GetToken()).String()})
		this_.priv_TypeCheck_ensureTyped_5459(n)
	}
	if NodeUtil_IsBlocklessArrowFunctionResult(n) {
		this_.priv_TypeCheck_visitImplicitReturnExpression_5448(t, n)
	}
	if (n.GetParent().IsForOf() || n.GetParent().IsForAwaitOf()) && (n.GetParent().GetFirstChild() == n) {
		this_.priv_TypeCheck_checkForOfTypes_5418(t, n.GetParent())
	}
	typeable = (typeable && !(this_.inExterns))
	if typeable {
		this_.debugTypeLogger.MaybeLogTypeOfNode(n)
		this_.priv_TypeCheck_doPercentTypedAccounting_5392(n)
	}
	this_.priv_TypeCheck_checkJsdocInfoContainsObjectWithBadKey_5464(n)
	this_.debugTypeLogger.StopLoggingIfThisIsWhereWeStarted(n)
}

// Java method: TypeCheck.visitUnaryPlus(com.google.javascript.rhino.Node n)
func (this_ *TypeCheck) priv_TypeCheck_visitUnaryPlus_5387(n pkg_rhino.Node) {
	_ = n
	if this_.priv_TypeCheck_getJSType_5456(n).IsNoType() {
		this_.priv_TypeCheck_report_5386(n, StaticGet_TypeCheck_UNARY_OPERATION(), []string{NodeUtil_OpToStr(n.GetToken()), this_.priv_TypeCheck_getJSType_5456(n.GetFirstChild()).String()})
	} else {
		this_.priv_TypeCheck_ensureTyped_5460(n, pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE())
	}
}

// Java method: TypeCheck.visitUnaryMinus(com.google.javascript.rhino.Node n)
func (this_ *TypeCheck) priv_TypeCheck_visitUnaryMinus_5388(n pkg_rhino.Node) {
	_ = n
	var operatorType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(n)
	_ = operatorType
	var childType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(n.GetFirstChild())
	_ = childType
	if operatorType.IsNumber() {
		this_.validator.ExpectNumber(n, childType, "sign operator")
		this_.priv_TypeCheck_ensureTyped_5460(n, pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE())
	} else {
		this_.validator.ExpectBigIntOrNumber(n, childType, "unary minus operator")
	}
}

// Java method: TypeCheck.visitBitwiseNOT(com.google.javascript.rhino.Node n)
func (this_ *TypeCheck) priv_TypeCheck_visitBitwiseNOT_5389(n pkg_rhino.Node) {
	_ = n
	var operatorType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(n)
	_ = operatorType
	var childType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(n.GetFirstChild())
	_ = childType
	if operatorType.IsNumber() {
		if !(childType.MatchesNumberContext()) {
			this_.priv_TypeCheck_report_5386(n, StaticGet_TypeCheck_BIT_OPERATION(), []string{NodeUtil_OpToStr(n.GetToken()), childType.String()})
		} else {
			this_.validator.ExpectNumberStrict(n, childType, "bitwise NOT")
		}
		this_.priv_TypeCheck_ensureTyped_5460(n, pkg_rhino.StaticGet_JSTypeNative_NUMBER_TYPE())
	} else {
		this_.validator.ExpectBigIntOrNumber(n, childType, "bitwise NOT")
	}
}

// Java method: TypeCheck.checkSpread(com.google.javascript.rhino.Node spreadNode)
func (this_ *TypeCheck) priv_TypeCheck_checkSpread_5390(spreadNode pkg_rhino.Node) {
	_ = spreadNode
	var target pkg_rhino.Node = spreadNode.GetOnlyChild()
	_ = target
	this_.priv_TypeCheck_ensureTyped_5459(target)
	var targetType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(target)
	_ = targetType
	switch spreadNode.GetParent().TokenSwitchOrd_() {
	case int32(47) /* Token.OBJECTLIT */ :
	case int32(46) /* Token.ARRAYLIT */, int32(28) /* Token.CALL */, int32(31) /* Token.OPTCHAIN_CALL */, int32(23) /* Token.NEW */ :
		this_.validator.ExpectAutoboxesToIterable(target, targetType, "Spread operator only applies to Iterable types")
	default:
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected parent of SPREAD: "+j2gort.StrValString(spreadNode.GetParent().ToStringTree()))))
	}
}

// Java method: TypeCheck.checkTypeofString(com.google.javascript.rhino.Node n, java.lang.String s)
func (this_ *TypeCheck) priv_TypeCheck_checkTypeofString_5391(n pkg_rhino.Node, s string) {
	_ = n
	_ = s
	if !(((((((((j2gort.CheckNotNullString(s) == "number") || (j2gort.CheckNotNullString(s) == "string")) || (j2gort.CheckNotNullString(s) == "boolean")) || (j2gort.CheckNotNullString(s) == "undefined")) || (j2gort.CheckNotNullString(s) == "function")) || (j2gort.CheckNotNullString(s) == "object")) || (j2gort.CheckNotNullString(s) == "symbol")) || (j2gort.CheckNotNullString(s) == "unknown")) || (j2gort.CheckNotNullString(s) == "bigint")) {
		this_.validator.ExpectValidTypeofName(n, s)
	}
}

// Java method: TypeCheck.doPercentTypedAccounting(com.google.javascript.rhino.Node n)
func (this_ *TypeCheck) priv_TypeCheck_doPercentTypedAccounting_5392(n pkg_rhino.Node) {
	_ = n
	var type_ pkg_rhino.JSType = n.GetJSType()
	_ = type_
	if type_ == nil || type_.IsJavaNil_() {
		this_.nullCount = this_.nullCount + 1
	} else if type_.IsUnknownType() {
		if this_.reportUnknownTypes && !(n.GetParent().IsExprResult()) {
			this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeCheck_UNKNOWN_EXPR_TYPE(), []string{}))
		}
		this_.unknownCount = this_.unknownCount + 1
	} else {
		this_.typedCount = this_.typedCount + 1
	}
}

// Java method: TypeCheck.checkGetterOrSetterType(com.google.javascript.rhino.Node child, com.google.javascript.rhino.jstype.FunctionType classType)
func (this_ *TypeCheck) priv_TypeCheck_checkGetterOrSetterType_5393(child pkg_rhino.Node, classType *pkg_rhino.FunctionType) {
	_ = child
	_ = classType
	var propertyName string = child.GetString()
	_ = propertyName
	var methodType *pkg_rhino.FunctionType = child.GetLastChild().GetJSType().ToMaybeFunctionType()
	_ = methodType
	var propertyType pkg_rhino.JSType
	_ = propertyType
	if child.IsGetterDef() {
		propertyType = this_.priv_TypeCheck_determineGetterType_5394(methodType)
	} else if methodType.GetParameterCount() == int32(1) {
		propertyType = methodType.GetParameter(int32(0)).GetJSType()
	} else {
		propertyType = j2gort.Cast[pkg_rhino.FunctionType_Parameter](j2gort.Ext_com_google_common_collect_Iterables_GetOnlyElement_1(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](methodType.GetParameters()))).GetJSType()
	}
	var officialPropertyType pkg_rhino.JSType
	_ = officialPropertyType
	if child.IsStaticMember() {
		officialPropertyType = classType.GetPropertyType_1_string(propertyName)
	} else {
		officialPropertyType = classType.GetPrototype().GetPropertyType_1_string(propertyName)
	}
	if !(propertyType.Equals(officialPropertyType)) {
		this_.priv_TypeCheck_report_5386(child, StaticGet_TypeCheck_CONFLICTING_GETTER_SETTER_TYPE(), []string{propertyName, func() string {
			if child.IsGetterDef() {
				return propertyType.String()
			}
			return officialPropertyType.String()
		}(), func() string {
			if child.IsGetterDef() {
				return officialPropertyType.String()
			}
			return propertyType.String()
		}()})
	}
}

// Java method: TypeCheck.determineGetterType(com.google.javascript.rhino.jstype.FunctionType methodType)
func (this_ *TypeCheck) priv_TypeCheck_determineGetterType_5394(methodType *pkg_rhino.FunctionType) pkg_rhino.JSType {
	_ = methodType
	if !(methodType.IsReturnTypeInferred()) {
		return methodType.GetReturnType()
	} else {
		return this_.typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
	}
}

// Java method: TypeCheck.visitAssign(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node assign)
func (this_ *TypeCheck) priv_TypeCheck_visitAssign_5395(t *NodeTraversal, assign pkg_rhino.Node) {
	_ = t
	_ = assign
	var info *pkg_rhino.JSDocInfo = assign.GetJSDocInfo()
	_ = info
	var lvalue pkg_rhino.Node = assign.GetFirstChild()
	_ = lvalue
	var rvalue pkg_rhino.Node = assign.GetLastChild()
	_ = rvalue
	var rightType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(rvalue)
	_ = rightType
	this_.priv_TypeCheck_checkCanAssignToWithScope_5396(t, assign, lvalue, rightType, info, "assignment")
	this_.priv_TypeCheck_ensureTyped_5461(assign, rightType)
}

// Java method: TypeCheck.checkCanAssignToWithScope(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node nodeToWarn, com.google.javascript.rhino.Node lvalue, com.google.javascript.rhino.jstype.JSType rightType, com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo info, java.lang.String msg)
func (this_ *TypeCheck) priv_TypeCheck_checkCanAssignToWithScope_5396(t *NodeTraversal, nodeToWarn pkg_rhino.Node, lvalue pkg_rhino.Node, rightType pkg_rhino.JSType, info *pkg_rhino.JSDocInfo, msg string) {
	_ = t
	_ = nodeToWarn
	_ = lvalue
	_ = rightType
	_ = info
	_ = msg
	if lvalue.IsDestructuringPattern() {
		this_.priv_TypeCheck_checkDestructuringAssignment_5397(t, nodeToWarn, lvalue, rightType, msg)
	} else {
		this_.priv_TypeCheck_checkCanAssignToNameGetpropOrGetelem_5398(t, nodeToWarn, lvalue, rightType, info, msg)
	}
}

// Java method: TypeCheck.checkDestructuringAssignment(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node nodeToWarn, com.google.javascript.rhino.Node pattern, com.google.javascript.rhino.jstype.JSType rightType, java.lang.String msg)
func (this_ *TypeCheck) priv_TypeCheck_checkDestructuringAssignment_5397(t *NodeTraversal, nodeToWarn pkg_rhino.Node, pattern pkg_rhino.Node, rightType pkg_rhino.JSType, msg string) {
	_ = t
	_ = nodeToWarn
	_ = pattern
	_ = rightType
	_ = msg
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(DestructuredTarget_CreateAllNonEmptyTargetsInPattern_3_JSTypeRegistry_JSType_Node(this_.typeRegistry, rightType, pattern)))
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			target := j2gort.CastPtr[DestructuredTarget](rawElem_2)
			_ = target
			this_.priv_TypeCheck_checkCanAssignToWithScope_5396(t, nodeToWarn, target.GetNode(), target.InferType(), nil, msg)
		}
	}
}

// Java method: TypeCheck.checkCanAssignToNameGetpropOrGetelem(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node nodeToWarn, com.google.javascript.rhino.Node lvalue, com.google.javascript.rhino.jstype.JSType rightType, com.google.javascript.rhino.JSDocInfo info, java.lang.String msg)
func (this_ *TypeCheck) priv_TypeCheck_checkCanAssignToNameGetpropOrGetelem_5398(t *NodeTraversal, nodeToWarn pkg_rhino.Node, lvalue pkg_rhino.Node, rightType pkg_rhino.JSType, info *pkg_rhino.JSDocInfo, msg string) {
	_ = t
	_ = nodeToWarn
	_ = lvalue
	_ = rightType
	_ = info
	_ = msg
	if !(((lvalue.IsName() || lvalue.IsGetProp()) || lvalue.IsGetElem()) || lvalue.IsCast()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(lvalue))
	}
	this_.priv_TypeCheck_checkNotReadonlyPropertyAssignment_5470(lvalue)
	if lvalue.IsGetProp() {
		var object pkg_rhino.Node = lvalue.GetFirstChild()
		_ = object
		var objectJsType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(object)
		_ = objectJsType
		var pname string = lvalue.GetString()
		_ = pname
		if object.IsGetProp() {
			var jsType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(object.GetFirstChild())
			_ = jsType
			if jsType.IsInterface() && (j2gort.CheckNotNullString(object.GetString()) == "prototype") {
				this_.priv_TypeCheck_visitInterfacePropertyAssignment_5415(object, lvalue)
			}
		}
		this_.priv_TypeCheck_checkEnumAlias_5455(t, info, rightType, nodeToWarn)
		this_.priv_TypeCheck_checkPropCreation_5399(lvalue)
		if j2gort.CheckNotNullString(pname) == "prototype" {
			this_.validator.ExpectCanAssignToPrototype(objectJsType, nodeToWarn, rightType)
			return
		}
		var objectCastType pkg_rhino.ObjectType = pkg_rhino.ObjectType_Cast(objectJsType.RestrictByNotNullOrUndefined())
		_ = objectCastType
		var expectedPropertyType pkg_rhino.JSType = this_.priv_TypeCheck_getPropertyTypeIfDeclared_5457(objectCastType, pname)
		_ = expectedPropertyType
		this_.priv_TypeCheck_checkPropertyInheritanceOnAssignment_5400(nodeToWarn, object, pname, info, expectedPropertyType)
		if !(expectedPropertyType.IsUnknownType()) {
			if !(TypeCheck_PropertyIsImplicitCast(objectCastType, pname)) {
				this_.validator.ExpectCanAssignToPropertyOf_5(nodeToWarn, rightType, expectedPropertyType, object, pname)
			}
			return
		}
	} else if lvalue.IsGetElem() && lvalue.GetSecondChild().GetJSType().IsKnownSymbolValueType() {
		var object_2 pkg_rhino.Node = lvalue.GetFirstChild()
		_ = object_2
		var objectJsType_2 pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(object_2)
		_ = objectJsType_2
		var objectCastType_2 pkg_rhino.ObjectType = pkg_rhino.ObjectType_Cast(objectJsType_2.RestrictByNotNullOrUndefined())
		_ = objectCastType_2
		var property *pkg_rhino.KnownSymbolType = lvalue.GetLastChild().GetJSType().ToMaybeKnownSymbolType()
		_ = property
		var expectedPropertyType_2 pkg_rhino.JSType = this_.priv_TypeCheck_getPropertyTypeIfDeclared_5458(objectCastType_2, pkg_rhino.NewProperty_SymbolKey(property))
		_ = expectedPropertyType_2
		this_.priv_TypeCheck_checkPropertyInheritanceOnAssignment_5401(nodeToWarn, object_2, pkg_rhino.NewProperty_SymbolKey(property), info, expectedPropertyType_2)
	}
	var leftType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(lvalue)
	_ = leftType
	if lvalue.IsQualifiedName() {
		var var_ *TypedVar = t.GetTypedScope().LookupTypedVar(lvalue.GetQualifiedName())
		_ = var_
		if var_ != nil {
			if var_.IsTypeInferred() {
				return
			}
			if NodeUtil_GetRootOfQualifiedName_1_Node(lvalue).IsThis() && (t.GetTypedScope() != j2gort.CastPtr[TypedScope](var_.GetScope())) {
				return
			}
			if !(j2gort.IfaceEquals[pkg_rhino.JSType](var_.GetType(), nil)) {
				leftType = var_.GetType()
			}
		}
	}
	this_.validator.ExpectCanAssignTo(nodeToWarn, rightType, leftType, msg)
}

// Java method: TypeCheck.checkPropCreation(com.google.javascript.rhino.Node lvalue)
func (this_ *TypeCheck) priv_TypeCheck_checkPropCreation_5399(lvalue pkg_rhino.Node) {
	_ = lvalue
	if lvalue.IsGetProp() {
		var objType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(lvalue.GetFirstChild())
		_ = objType
		if !(objType.IsEmptyType()) && !(objType.IsUnknownType()) {
			var propName string = lvalue.GetString()
			_ = propName
			var kind pkg_rhino.JSTypeRegistry_PropDefinitionKind = this_.typeRegistry.CanPropertyBeDefined(objType, propName)
			_ = kind
			{
				var iifeCond_1 bool
				{

					re_3 := kind
					ae_4 := pkg_rhino.StaticGet_JSTypeRegistry_PropDefinitionKind_KNOWN()
					if re_3 == 0 {
						j2gort.ThrowNPE()
					}
					iifeCond_1 = re_3 == ae_4
				}
				if !(iifeCond_1) {
					if objType.IsStruct() {
						if objType.RestrictByNotNullOrUndefined().IsUnionType() {
							this_.priv_TypeCheck_report_5386(lvalue, StaticGet_TypeCheck_ILLEGAL_PROPERTY_CREATION_ON_UNION_TYPE(), []string{})
						} else {
							this_.priv_TypeCheck_report_5386(lvalue, StaticGet_TypeCheck_ILLEGAL_PROPERTY_CREATION(), []string{})
						}
					} else {
						if (!(objType.IsNoType()) && !(objType.IsUnknownType())) && objType.IsSubtypeOf_1(this_.priv_TypeCheck_getNativeType_5463(pkg_rhino.StaticGet_JSTypeNative_NULL_VOID())) {
							return
						}
						this_.priv_TypeCheck_reportMissingProperty_5427(lvalue.GetFirstChild(), objType, lvalue, kind, true)
					}
				}
			}
		}
	}
}

// Java method: TypeCheck.checkPropertyInheritanceOnAssignment(com.google.javascript.rhino.Node assign, com.google.javascript.rhino.Node object, java.lang.String property, com.google.javascript.rhino.JSDocInfo info, com.google.javascript.rhino.jstype.JSType propertyType)
func (this_ *TypeCheck) priv_TypeCheck_checkPropertyInheritanceOnAssignment_5400(assign pkg_rhino.Node, object pkg_rhino.Node, property string, info *pkg_rhino.JSDocInfo, propertyType pkg_rhino.JSType) {
	_ = assign
	_ = object
	_ = property
	_ = info
	_ = propertyType
	if object.IsGetProp() && (j2gort.CheckNotNullString(object.GetString()) == "prototype") {
		var preObject pkg_rhino.Node = object.GetFirstChild()
		_ = preObject
		var ctorType *pkg_rhino.FunctionType = this_.priv_TypeCheck_getJSType_5456(preObject).ToMaybeFunctionType()
		_ = ctorType
		if (ctorType == nil) || !(ctorType.HasInstanceType()) {
			return
		}
		this_.priv_TypeCheck_checkDeclaredPropertyAgainstNominalInheritance_5410(assign, ctorType, property, info, propertyType)
		this_.priv_TypeCheck_checkAbstractMethodInConcreteClass_5414(assign, ctorType, info)
	} else {
		var ctorType_2 *pkg_rhino.FunctionType = this_.priv_TypeCheck_getJSType_5456(object).ToMaybeFunctionType()
		_ = ctorType_2
		if (ctorType_2 == nil) || !(ctorType_2.HasInstanceType()) {
			return
		}
		this_.priv_TypeCheck_checkDeclaredPropertyAgainstPrototypalInheritance_5411(assign, pkg_rhino.ToIface_ObjectType(ctorType_2), property, info, propertyType)
	}
}

// Java method: TypeCheck.checkPropertyInheritanceOnAssignment(com.google.javascript.rhino.Node assign, com.google.javascript.rhino.Node object, com.google.javascript.rhino.jstype.Property.Key property, com.google.javascript.rhino.JSDocInfo info, com.google.javascript.rhino.jstype.JSType propertyType)
func (this_ *TypeCheck) priv_TypeCheck_checkPropertyInheritanceOnAssignment_5401(assign pkg_rhino.Node, object pkg_rhino.Node, property pkg_rhino.Property_Key, info *pkg_rhino.JSDocInfo, propertyType pkg_rhino.JSType) {
	_ = assign
	_ = object
	_ = property
	_ = info
	_ = propertyType
	if object.IsGetProp() && (j2gort.CheckNotNullString(object.GetString()) == "prototype") {
		var preObject pkg_rhino.Node = object.GetFirstChild()
		_ = preObject
		var ctorType *pkg_rhino.FunctionType = this_.priv_TypeCheck_getJSType_5456(preObject).ToMaybeFunctionType()
		_ = ctorType
		if (ctorType == nil) || !(ctorType.HasInstanceType()) {
			return
		}
		this_.priv_TypeCheck_checkDeclaredPropertyAgainstNominalInheritance_5412(assign, ctorType, property, info, propertyType)
		this_.priv_TypeCheck_checkAbstractMethodInConcreteClass_5414(assign, ctorType, info)
	} else {
		var ctorType_2 *pkg_rhino.FunctionType = this_.priv_TypeCheck_getJSType_5456(object).ToMaybeFunctionType()
		_ = ctorType_2
		if (ctorType_2 == nil) || !(ctorType_2.HasInstanceType()) {
			return
		}
		this_.priv_TypeCheck_checkDeclaredPropertyAgainstPrototypalInheritance_5413(assign, pkg_rhino.ToIface_ObjectType(ctorType_2), property, info, propertyType)
	}
}

// Java method: TypeCheck.checkPropertyInheritanceOnPrototypeLitKey(com.google.javascript.rhino.Node key, java.lang.String propertyName, com.google.javascript.rhino.jstype.ObjectType type)
func (this_ *TypeCheck) priv_TypeCheck_checkPropertyInheritanceOnPrototypeLitKey_5402(key pkg_rhino.Node, propertyName string, type_ pkg_rhino.ObjectType) {
	_ = key
	_ = propertyName
	_ = type_
	this_.priv_TypeCheck_checkPropertyInheritance_5404(key, propertyName, type_.GetOwnerFunction(), type_)
}

// Java method: TypeCheck.checkPropertyInheritanceOnClassMember(com.google.javascript.rhino.Node key, java.lang.String propertyName, com.google.javascript.rhino.jstype.FunctionType ctorType)
func (this_ *TypeCheck) priv_TypeCheck_checkPropertyInheritanceOnClassMember_5403(key pkg_rhino.Node, propertyName string, ctorType *pkg_rhino.FunctionType) {
	_ = key
	_ = propertyName
	_ = ctorType
	if key.IsStaticMember() {
		this_.priv_TypeCheck_checkDeclaredPropertyAgainstPrototypalInheritance_5411(key, pkg_rhino.ToIface_ObjectType(ctorType), propertyName, key.GetJSDocInfo(), ctorType.GetPropertyType_1_string(propertyName))
	} else {
		this_.priv_TypeCheck_checkPropertyInheritance_5404(key, propertyName, ctorType, ctorType.GetInstanceType())
	}
}

// Java method: TypeCheck.checkPropertyInheritance(com.google.javascript.rhino.Node key, java.lang.String propertyName, com.google.javascript.rhino.jstype.FunctionType ctorType, com.google.javascript.rhino.jstype.ObjectType type)
func (this_ *TypeCheck) priv_TypeCheck_checkPropertyInheritance_5404(key pkg_rhino.Node, propertyName string, ctorType *pkg_rhino.FunctionType, type_ pkg_rhino.ObjectType) {
	_ = key
	_ = propertyName
	_ = ctorType
	_ = type_
	if (ctorType == nil) || !(ctorType.HasInstanceType()) {
		return
	}
	this_.priv_TypeCheck_checkDeclaredPropertyAgainstNominalInheritance_5410(key.GetFirstChild(), ctorType, propertyName, key.GetJSDocInfo(), type_.GetPropertyType_1_string(propertyName))
	this_.priv_TypeCheck_checkAbstractMethodInConcreteClass_5414(key, ctorType, key.GetJSDocInfo())
}

// Java method: TypeCheck.checkPropertyInheritanceOnPrototypeLitKey(com.google.javascript.rhino.Node key, com.google.javascript.rhino.jstype.Property.Key propertyName, com.google.javascript.rhino.jstype.ObjectType type)
func (this_ *TypeCheck) priv_TypeCheck_checkPropertyInheritanceOnPrototypeLitKey_5405(key pkg_rhino.Node, propertyName pkg_rhino.Property_Key, type_ pkg_rhino.ObjectType) {
	_ = key
	_ = propertyName
	_ = type_
	this_.priv_TypeCheck_checkPropertyInheritance_5407(key, propertyName, type_.GetOwnerFunction(), type_)
}

// Java method: TypeCheck.checkPropertyInheritanceOnClassMember(com.google.javascript.rhino.Node key, com.google.javascript.rhino.jstype.Property.Key propertyName, com.google.javascript.rhino.jstype.FunctionType ctorType)
func (this_ *TypeCheck) priv_TypeCheck_checkPropertyInheritanceOnClassMember_5406(key pkg_rhino.Node, propertyName pkg_rhino.Property_Key, ctorType *pkg_rhino.FunctionType) {
	_ = key
	_ = propertyName
	_ = ctorType
	if key.IsStaticMember() {
		this_.priv_TypeCheck_checkDeclaredPropertyAgainstPrototypalInheritance_5413(key, pkg_rhino.ToIface_ObjectType(ctorType), propertyName, key.GetJSDocInfo(), ctorType.GetPropertyType_1_Property_Key(propertyName))
	} else {
		this_.priv_TypeCheck_checkPropertyInheritance_5407(key, propertyName, ctorType, ctorType.GetInstanceType())
	}
}

// Java method: TypeCheck.checkPropertyInheritance(com.google.javascript.rhino.Node key, com.google.javascript.rhino.jstype.Property.Key propertyName, com.google.javascript.rhino.jstype.FunctionType ctorType, com.google.javascript.rhino.jstype.ObjectType type)
func (this_ *TypeCheck) priv_TypeCheck_checkPropertyInheritance_5407(key pkg_rhino.Node, propertyName pkg_rhino.Property_Key, ctorType *pkg_rhino.FunctionType, type_ pkg_rhino.ObjectType) {
	_ = key
	_ = propertyName
	_ = ctorType
	_ = type_
	if (ctorType == nil) || !(ctorType.HasInstanceType()) {
		return
	}
	this_.priv_TypeCheck_checkDeclaredPropertyAgainstNominalInheritance_5412(key.GetFirstChild(), ctorType, propertyName, key.GetJSDocInfo(), type_.GetPropertyType_1_Property_Key(propertyName))
	this_.priv_TypeCheck_checkAbstractMethodInConcreteClass_5414(key, ctorType, key.GetJSDocInfo())
}

// Java method: TypeCheck.visitObjectPattern(com.google.javascript.rhino.Node pattern)
func (this_ *TypeCheck) priv_TypeCheck_visitObjectPattern_5408(pattern pkg_rhino.Node) {
	_ = pattern
	var patternType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(pattern)
	_ = patternType
	this_.validator.ExpectObject(pattern, patternType, "cannot destructure 'null' or 'undefined'")
	{
		var child pkg_rhino.Node = pattern.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			var target *DestructuredTarget = DestructuredTarget_CreateTarget_3_JSTypeRegistry_JSType_Node(this_.typeRegistry, patternType, child)
			_ = target
			if target.HasComputedProperty() {
				var computedProperty pkg_rhino.Node = target.GetComputedProperty()
				_ = computedProperty
				this_.validator.ExpectIndexMatch(computedProperty, patternType, this_.priv_TypeCheck_getJSType_5456(computedProperty.GetFirstChild()))
			} else if target.HasStringKey() {
				var stringKey pkg_rhino.Node = target.GetStringKey()
				_ = stringKey
				if !(stringKey.IsQuotedStringKey()) {
					if patternType.IsDict() {
						this_.priv_TypeCheck_report_5386(stringKey, StaticGet_TypeValidator_ILLEGAL_PROPERTY_ACCESS(), []string{"unquoted", "dict"})
					}
					this_.priv_TypeCheck_checkPropertyAccessForDestructuring_5422(pattern, patternType, stringKey, this_.priv_TypeCheck_getJSType_5456(target.GetNode()))
				} else if patternType.IsStruct() {
					this_.priv_TypeCheck_report_5386(stringKey, StaticGet_TypeValidator_ILLEGAL_PROPERTY_ACCESS(), []string{"quoted", "struct"})
				}
			}
		}
	}
	this_.priv_TypeCheck_ensureTyped_5459(pattern)
}

// Java method: TypeCheck.visitObjectOrClassLiteralKey(com.google.javascript.rhino.Node key, com.google.javascript.rhino.Node owner, com.google.javascript.rhino.jstype.JSType ownerType)
func (this_ *TypeCheck) priv_TypeCheck_visitObjectOrClassLiteralKey_5409(key pkg_rhino.Node, owner pkg_rhino.Node, ownerType pkg_rhino.JSType) {
	_ = key
	_ = owner
	_ = ownerType
	if owner.IsFromExterns() {
		this_.priv_TypeCheck_ensureTyped_5459(key)
		return
	}
	if key.IsComputedProp() || key.IsComputedFieldDef() {
		var keyType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(key.GetFirstChild())
		_ = keyType
		this_.validator.ExpectIndexMatch(key, ownerType, keyType)
		if keyType.IsKnownSymbolValueType() {
			if owner.IsClass() {
				var classConstructorType *pkg_rhino.FunctionType = owner.GetJSType().AssertFunctionType()
				_ = classConstructorType
				this_.priv_TypeCheck_checkPropertyInheritanceOnClassMember_5406(key, pkg_rhino.NewProperty_SymbolKey(keyType.ToMaybeKnownSymbolType()), classConstructorType)
			} else {
				if !(owner.IsObjectLit()) {
					j2gort.CheckStateFmt(false, "Unexpected owner %s", j2gort.BoxNode(owner))
				}
				var objectCastType pkg_rhino.ObjectType = pkg_rhino.ObjectType_Cast(ownerType.RestrictByNotNullOrUndefined())
				_ = objectCastType
				var expectedPropertyType pkg_rhino.JSType = this_.priv_TypeCheck_getPropertyTypeIfDeclared_5458(objectCastType, pkg_rhino.NewProperty_SymbolKey(keyType.ToMaybeKnownSymbolType()))
				_ = expectedPropertyType
				this_.validator.ExpectCanAssignToPropertyOf_5(key, expectedPropertyType, this_.priv_TypeCheck_getJSType_5456(key.GetSecondChild()), owner, keyType.ToMaybeKnownSymbolType().GetDisplayName())
			}
		}
		return
	}
	if key.IsBlock() {
		return
	}
	if key.IsQuotedStringKey() {
		if ownerType.IsStruct() {
			this_.priv_TypeCheck_report_5386(key, func() *DiagnosticType {
				if owner.IsClass() {
					return StaticGet_TypeCheck_ILLEGAL_CLASS_KEY()
				}
				return StaticGet_TypeCheck_ILLEGAL_OBJLIT_KEY()
			}(), []string{"struct"})
		}
	} else {
		if ownerType.IsDict() && !(NodeUtil_IsEs6ConstructorMemberFunctionDef(key)) {
			this_.priv_TypeCheck_report_5386(key, func() *DiagnosticType {
				if owner.IsClass() {
					return StaticGet_TypeCheck_ILLEGAL_CLASS_KEY()
				}
				return StaticGet_TypeCheck_ILLEGAL_OBJLIT_KEY()
			}(), []string{"dict"})
		}
	}
	if key.IsSpread() {
		return
	}
	var rvalue pkg_rhino.Node = key.GetFirstChild()
	_ = rvalue
	if rvalue == 0 {
		this_.priv_TypeCheck_ensureTyped_5459(key)
		return
	}
	var rightType pkg_rhino.JSType = TypeCheck_GetObjectLitKeyTypeFromValueType(key, this_.priv_TypeCheck_getJSType_5456(rvalue))
	_ = rightType
	if rightType == nil || rightType.IsJavaNil_() {
		rightType = this_.priv_TypeCheck_getNativeType_5463(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
	}
	var keyType_2 pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(key)
	_ = keyType_2
	var propertyName string = NodeUtil_GetObjectOrClassLitKeyName(key)
	_ = propertyName
	var allowedValueType pkg_rhino.JSType = keyType_2
	_ = allowedValueType
	if allowedValueType.IsEnumElementType() {
		allowedValueType = allowedValueType.ToMaybeEnumElementType().GetPrimitiveType()
	}
	var valid bool = this_.validator.ExpectCanAssignToPropertyOf_5(key, rightType, allowedValueType, owner, propertyName)
	_ = valid
	if valid {
		this_.priv_TypeCheck_ensureTyped_5461(key, rightType)
	} else {
		this_.priv_TypeCheck_ensureTyped_5459(key)
	}
	j2gort.CheckNotNullString(propertyName)
	if owner.IsClass() {
		var classConstructorType_2 *pkg_rhino.FunctionType = owner.GetJSType().AssertFunctionType()
		_ = classConstructorType_2
		this_.priv_TypeCheck_checkPropertyInheritanceOnClassMember_5403(key, propertyName, classConstructorType_2)
	} else if !(j2gort.IfaceEquals[pkg_rhino.ObjectType](ownerType.ToMaybeObjectType(), nil)) {
		this_.priv_TypeCheck_checkPropertyInheritanceOnPrototypeLitKey_5402(key, propertyName, ownerType.ToMaybeObjectType())
	}
}

// Java static method: TypeCheck.propertyIsImplicitCast(com.google.javascript.rhino.jstype.ObjectType type, java.lang.String prop)
func TypeCheck_PropertyIsImplicitCast(type_ pkg_rhino.ObjectType, prop string) bool {
	clinit_TypeCheck()
	_ = type_
	_ = prop
	for ; type_ != nil && !type_.IsJavaNil_(); type_ = type_.GetImplicitPrototype() {
		var docInfo *pkg_rhino.JSDocInfo = type_.GetOwnPropertyJSDocInfo_1_string(prop)
		_ = docInfo
		if (docInfo != nil) && docInfo.IsImplicitCast() {
			return true
		}
	}
	return false
}

// Java method: TypeCheck.checkDeclaredPropertyAgainstNominalInheritance(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.FunctionType ctorType, java.lang.String propertyName, com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo info, com.google.javascript.rhino.jstype.JSType propertyType)
func (this_ *TypeCheck) priv_TypeCheck_checkDeclaredPropertyAgainstNominalInheritance_5410(n pkg_rhino.Node, ctorType *pkg_rhino.FunctionType, propertyName string, info *pkg_rhino.JSDocInfo, propertyType pkg_rhino.JSType) {
	_ = n
	_ = ctorType
	_ = propertyName
	_ = info
	_ = propertyType
	if ("__proto__" == propertyName) || ("constructor" == propertyName) {
		return
	}
	if ctorType.IsInterface() {
		return
	}
	if TypeCheck_HasUnknownOrEmptySupertype(ctorType) {
		return
	}
	var foundProperty bool = false
	_ = foundProperty
	var superCtor *pkg_rhino.FunctionType = ctorType.GetSuperClassConstructor()
	_ = superCtor
	if superCtor != nil {
		var propSlot *pkg_rhino.Property_OwnedProperty = superCtor.GetInstanceType().FindClosestDefinition_1_string(propertyName)
		_ = propSlot
		var superClassHasProperty bool = ((propSlot != nil) && !(propSlot.IsOwnedByInterface()))
		_ = superClassHasProperty
		foundProperty = func() bool { ba_5 := foundProperty; bb_6 := superClassHasProperty; return ba_5 || bb_6 }()
		if superClassHasProperty {
			var superClass pkg_rhino.ObjectType = propSlot.GetOwnerInstanceType()
			_ = superClass
			var superClassHasDeclaredProperty bool = !(propSlot.GetValue().IsTypeInferred())
			_ = superClassHasDeclaredProperty
			if superClassHasDeclaredProperty {
				if TypeCheck_IsDeclaredLocally_2_FunctionType_string(ctorType, propertyName) && !(TypeCheck_DeclaresOverride(info)) {
					this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeCheck_HIDDEN_SUPERCLASS_PROPERTY(), []string{propertyName, superClass.GetReferenceName()}))
				}
				this_.validator.CheckPropertyType_5_Node_JSType_ObjectType_string_JSType(n, ctorType.GetTypeOfThis(), superClass, propertyName, propertyType)
			}
		}
	}
	{
		cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(ctorType.GetAllImplementedInterfaces()))
		for cur_7.HasNext() {
			rawElem_8 := cur_7.Next()
			implementedInterface := j2gort.AsIface[pkg_rhino.ObjectType](rawElem_8)
			_ = implementedInterface
			if implementedInterface.IsUnknownType() || implementedInterface.IsEmptyType() {
				continue
			}
			if !(implementedInterface.IsInstanceType()) {
				j2gort.CheckStateMsg(false, implementedInterface)
			}
			var propSlot_2 *pkg_rhino.Property_OwnedProperty = implementedInterface.FindClosestDefinition_1_string(propertyName)
			_ = propSlot_2
			var interfaceHasProperty bool = (propSlot_2 != nil)
			_ = interfaceHasProperty
			foundProperty = func() bool { ba_9 := foundProperty; bb_10 := interfaceHasProperty; return ba_9 || bb_10 }()
			if interfaceHasProperty && !(TypeCheck_DeclaresOverride(info)) {
				this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeCheck_HIDDEN_INTERFACE_PROPERTY(), []string{propertyName, propSlot_2.GetOwnerInstanceType().GetReferenceName()}))
			}
		}
	}
	if !(foundProperty) && TypeCheck_DeclaresOverride(info) {
		this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeCheck_UNKNOWN_OVERRIDE(), []string{propertyName, ctorType.GetInstanceType().GetReferenceName()}))
	}
}

// Java static method: TypeCheck.isDeclaredLocally(com.google.javascript.rhino.jstype.FunctionType ctorType, java.lang.String propertyName)
func TypeCheck_IsDeclaredLocally_2_FunctionType_string(ctorType *pkg_rhino.FunctionType, propertyName string) bool {
	clinit_TypeCheck()
	_ = ctorType
	_ = propertyName
	if !(ctorType.IsConstructor()) {
		j2gort.CheckState(false)
	}
	return (ctorType.GetPrototype().HasOwnProperty_1_string(propertyName) || ctorType.GetInstanceType().HasOwnProperty_1_string(propertyName))
}

// Java method: TypeCheck.checkDeclaredPropertyAgainstPrototypalInheritance(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.ObjectType receiverType, java.lang.String propertyName, com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo info, com.google.javascript.rhino.jstype.JSType propertyType)
func (this_ *TypeCheck) priv_TypeCheck_checkDeclaredPropertyAgainstPrototypalInheritance_5411(n pkg_rhino.Node, receiverType pkg_rhino.ObjectType, propertyName string, info *pkg_rhino.JSDocInfo, propertyType pkg_rhino.JSType) {
	_ = n
	_ = receiverType
	_ = propertyName
	_ = info
	_ = propertyType
	var declaredOverride bool = TypeCheck_DeclaresOverride(info)
	_ = declaredOverride
	var supertypeWithProperty pkg_rhino.ObjectType = nil
	_ = supertypeWithProperty
	{
		var type_ pkg_rhino.ObjectType = receiverType.GetImplicitPrototype()
		_ = type_
		for ; type_ != nil && !type_.IsJavaNil_(); type_ = type_.GetImplicitPrototype() {
			if type_.HasOwnProperty_1_string(propertyName) {
				supertypeWithProperty = type_
				break
			}
		}
	}
	if supertypeWithProperty == nil || supertypeWithProperty.IsJavaNil_() {
		if declaredOverride && !(receiverType.LoosenTypecheckingDueToForwardReferencedSupertype()) {
			this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeCheck_UNKNOWN_PROTOTYPAL_OVERRIDE(), []string{propertyName, receiverType.String()}))
		}
	} else {
		if !(declaredOverride) {
			this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY(), []string{propertyName, supertypeWithProperty.String()}))
		}
		var overriddenPropertyType pkg_rhino.JSType = supertypeWithProperty.GetPropertyType_1_string(propertyName)
		_ = overriddenPropertyType
		if !(propertyType.IsSubtypeOf_1(overriddenPropertyType)) {
			this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY_MISMATCH(), []string{propertyName, supertypeWithProperty.String(), overriddenPropertyType.String(), propertyType.String()}))
		}
	}
}

// Java method: TypeCheck.checkDeclaredPropertyAgainstNominalInheritance(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.FunctionType ctorType, com.google.javascript.rhino.jstype.Property.Key propertyName, com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo info, com.google.javascript.rhino.jstype.JSType propertyType)
func (this_ *TypeCheck) priv_TypeCheck_checkDeclaredPropertyAgainstNominalInheritance_5412(n pkg_rhino.Node, ctorType *pkg_rhino.FunctionType, propertyName pkg_rhino.Property_Key, info *pkg_rhino.JSDocInfo, propertyType pkg_rhino.JSType) {
	_ = n
	_ = ctorType
	_ = propertyName
	_ = info
	_ = propertyType
	if propertyName.Matches("__proto__") || propertyName.Matches("constructor") {
		return
	}
	if ctorType.IsInterface() {
		return
	}
	if TypeCheck_HasUnknownOrEmptySupertype(ctorType) {
		return
	}
	var foundProperty bool = false
	_ = foundProperty
	var superCtor *pkg_rhino.FunctionType = ctorType.GetSuperClassConstructor()
	_ = superCtor
	if superCtor != nil {
		var propSlot *pkg_rhino.Property_OwnedProperty = superCtor.GetInstanceType().FindClosestDefinition_1_Property_Key(propertyName)
		_ = propSlot
		var superClassHasProperty bool = ((propSlot != nil) && !(propSlot.IsOwnedByInterface()))
		_ = superClassHasProperty
		foundProperty = func() bool { ba_11 := foundProperty; bb_12 := superClassHasProperty; return ba_11 || bb_12 }()
		if superClassHasProperty {
			var superClass pkg_rhino.ObjectType = propSlot.GetOwnerInstanceType()
			_ = superClass
			var superClassHasDeclaredProperty bool = !(propSlot.GetValue().IsTypeInferred())
			_ = superClassHasDeclaredProperty
			if superClassHasDeclaredProperty {
				if TypeCheck_IsDeclaredLocally_2_FunctionType_Property_Key(ctorType, propertyName) && !(TypeCheck_DeclaresOverride(info)) {
					this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeCheck_HIDDEN_SUPERCLASS_PROPERTY(), []string{propertyName.HumanReadableName(), superClass.GetReferenceName()}))
				}
				this_.validator.CheckPropertyType_5_Node_JSType_ObjectType_Property_Key_JSType(n, ctorType.GetTypeOfThis(), superClass, propertyName, propertyType)
			}
		}
	}
	{
		cur_13 := j2gort.NewIterCursor(j2gort.CheckNotNull(ctorType.GetAllImplementedInterfaces()))
		for cur_13.HasNext() {
			rawElem_14 := cur_13.Next()
			implementedInterface := j2gort.AsIface[pkg_rhino.ObjectType](rawElem_14)
			_ = implementedInterface
			if implementedInterface.IsUnknownType() || implementedInterface.IsEmptyType() {
				continue
			}
			if !(implementedInterface.IsInstanceType()) {
				j2gort.CheckStateMsg(false, implementedInterface)
			}
			var propSlot_2 *pkg_rhino.Property_OwnedProperty = implementedInterface.FindClosestDefinition_1_Property_Key(propertyName)
			_ = propSlot_2
			var interfaceHasProperty bool = (propSlot_2 != nil)
			_ = interfaceHasProperty
			foundProperty = func() bool { ba_15 := foundProperty; bb_16 := interfaceHasProperty; return ba_15 || bb_16 }()
			if interfaceHasProperty && !(TypeCheck_DeclaresOverride(info)) {
				this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeCheck_HIDDEN_INTERFACE_PROPERTY(), []string{propertyName.HumanReadableName(), propSlot_2.GetOwnerInstanceType().GetReferenceName()}))
			}
		}
	}
	if !(foundProperty) && TypeCheck_DeclaresOverride(info) {
		this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeCheck_UNKNOWN_OVERRIDE(), []string{propertyName.HumanReadableName(), ctorType.GetInstanceType().GetReferenceName()}))
	}
}

// Java static method: TypeCheck.isDeclaredLocally(com.google.javascript.rhino.jstype.FunctionType ctorType, com.google.javascript.rhino.jstype.Property.Key propertyName)
func TypeCheck_IsDeclaredLocally_2_FunctionType_Property_Key(ctorType *pkg_rhino.FunctionType, propertyName pkg_rhino.Property_Key) bool {
	clinit_TypeCheck()
	_ = ctorType
	_ = propertyName
	if !(ctorType.IsConstructor()) {
		j2gort.CheckState(false)
	}
	return (ctorType.GetPrototype().HasOwnProperty_1_Property_Key(propertyName) || ctorType.GetInstanceType().HasOwnProperty_1_Property_Key(propertyName))
}

// Java method: TypeCheck.checkDeclaredPropertyAgainstPrototypalInheritance(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.ObjectType receiverType, com.google.javascript.rhino.jstype.Property.Key propertyName, com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo info, com.google.javascript.rhino.jstype.JSType propertyType)
func (this_ *TypeCheck) priv_TypeCheck_checkDeclaredPropertyAgainstPrototypalInheritance_5413(n pkg_rhino.Node, receiverType pkg_rhino.ObjectType, propertyName pkg_rhino.Property_Key, info *pkg_rhino.JSDocInfo, propertyType pkg_rhino.JSType) {
	_ = n
	_ = receiverType
	_ = propertyName
	_ = info
	_ = propertyType
	var declaredOverride bool = TypeCheck_DeclaresOverride(info)
	_ = declaredOverride
	var supertypeWithProperty pkg_rhino.ObjectType = j2gort.AsIface[pkg_rhino.ObjectType](j2gort.Ext_com_google_common_collect_Streams_Stream_1_Ext_java_lang_Iterable(receiverType.GetImplicitPrototypeChain()).Filter(j2gort.NewLambda_Ext_java_util_function_Predicate(func(lp_17 any) bool {
		type_ := j2gort.AsIface[pkg_rhino.ObjectType](lp_17)
		_ = type_
		return type_.HasOwnProperty_1_Property_Key(propertyName)
	})).FindFirst().OrElse(nil))
	_ = supertypeWithProperty
	if supertypeWithProperty == nil || supertypeWithProperty.IsJavaNil_() {
		if declaredOverride && !(receiverType.LoosenTypecheckingDueToForwardReferencedSupertype()) {
			this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeCheck_UNKNOWN_PROTOTYPAL_OVERRIDE(), []string{propertyName.HumanReadableName(), receiverType.String()}))
		}
	} else {
		if !(declaredOverride) {
			this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY(), []string{propertyName.HumanReadableName(), supertypeWithProperty.String()}))
		}
		var overriddenPropertyType pkg_rhino.JSType = supertypeWithProperty.GetPropertyType_1_Property_Key(propertyName)
		_ = overriddenPropertyType
		if !(propertyType.IsSubtypeOf_1(overriddenPropertyType)) {
			this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeCheck_HIDDEN_PROTOTYPAL_SUPERTYPE_PROPERTY_MISMATCH(), []string{propertyName.HumanReadableName(), supertypeWithProperty.String(), overriddenPropertyType.String(), propertyType.String()}))
		}
	}
}

// Java method: TypeCheck.checkAbstractMethodInConcreteClass(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.FunctionType ctorType, com.google.javascript.rhino.JSDocInfo info)
func (this_ *TypeCheck) priv_TypeCheck_checkAbstractMethodInConcreteClass_5414(n pkg_rhino.Node, ctorType *pkg_rhino.FunctionType, info *pkg_rhino.JSDocInfo) {
	_ = n
	_ = ctorType
	_ = info
	if (info == nil) || !(info.IsAbstract()) {
		return
	}
	if ctorType.IsConstructor() && !(ctorType.IsAbstract()) {
		this_.priv_TypeCheck_report_5386(n, StaticGet_TypeCheck_ABSTRACT_METHOD_IN_CONCRETE_CLASS(), []string{})
	}
}

// Java static method: TypeCheck.hasUnknownOrEmptySupertype(com.google.javascript.rhino.jstype.FunctionType ctor)
func TypeCheck_HasUnknownOrEmptySupertype(ctor *pkg_rhino.FunctionType) bool {
	clinit_TypeCheck()
	_ = ctor
	if !(ctor.IsConstructor() || ctor.IsInterface()) {
		j2gort.CheckArgument(false)
	}
	if !(!(ctor.IsUnknownType())) {
		j2gort.CheckArgument(false)
	}
	for {
		var maybeSuperInstanceType pkg_rhino.ObjectType = ctor.GetPrototype().GetImplicitPrototype()
		_ = maybeSuperInstanceType
		if maybeSuperInstanceType == nil || maybeSuperInstanceType.IsJavaNil_() {
			return false
		}
		if maybeSuperInstanceType.IsUnknownType() || maybeSuperInstanceType.IsEmptyType() {
			return true
		}
		ctor = maybeSuperInstanceType.GetConstructor()
		if ctor == nil {
			return false
		}
		if !(ctor.IsConstructor() || ctor.IsInterface()) {
			j2gort.CheckState(false)
		}
	}
	panic("unreachable")
}

// Java static method: TypeCheck.getObjectLitKeyTypeFromValueType(com.google.javascript.rhino.Node key, com.google.javascript.rhino.jstype.JSType valueType)
func TypeCheck_GetObjectLitKeyTypeFromValueType(key pkg_rhino.Node, valueType pkg_rhino.JSType) pkg_rhino.JSType {
	clinit_TypeCheck()
	_ = key
	_ = valueType
	if valueType != nil && !valueType.IsJavaNil_() {
		switch key.TokenSwitchOrd_() {
		case int32(95) /* Token.GETTER_DEF */ :
			if valueType.IsFunctionType() {
				var fntype *pkg_rhino.FunctionType = valueType.ToMaybeFunctionType()
				_ = fntype
				valueType = fntype.GetReturnType()
			} else {
				return nil
			}
		case int32(96) /* Token.SETTER_DEF */ :
			if valueType.IsFunctionType() {
				var fntype_2 *pkg_rhino.FunctionType = valueType.ToMaybeFunctionType()
				_ = fntype_2
				var param pkg_rhino.FunctionType_Parameter = fntype_2.GetParameter(int32(0))
				_ = param
				valueType = param.GetJSType()
			} else {
				return nil
			}
		default:
		}
	}
	return valueType
}

// Java method: TypeCheck.visitInterfacePropertyAssignment(com.google.javascript.rhino.Node object, com.google.javascript.rhino.Node lvalue)
func (this_ *TypeCheck) priv_TypeCheck_visitInterfacePropertyAssignment_5415(object pkg_rhino.Node, lvalue pkg_rhino.Node) {
	_ = object
	_ = lvalue
	if !(lvalue.GetParent().IsAssign()) {
		this_.priv_TypeCheck_reportInvalidInterfaceMemberDeclaration_5416(object)
		return
	}
	var assign pkg_rhino.Node = lvalue.GetParent()
	_ = assign
	var rvalue pkg_rhino.Node = assign.GetSecondChild()
	_ = rvalue
	var rvalueType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(rvalue)
	_ = rvalueType
	if !(rvalueType.IsFunctionType()) {
		this_.priv_TypeCheck_reportInvalidInterfaceMemberDeclaration_5416(object)
	}
	if rvalue.IsFunction() && !(NodeUtil_IsEmptyBlock(NodeUtil_GetFunctionBody(rvalue))) {
		var abstractMethodName string = this_.compiler.GetCodingConvention().GetAbstractMethodName()
		_ = abstractMethodName
		this_.compiler.Report_1(JSError_Make_3(object, StaticGet_TypeCheck_INTERFACE_METHOD_NOT_EMPTY(), []string{abstractMethodName}))
	}
}

// Java method: TypeCheck.reportInvalidInterfaceMemberDeclaration(com.google.javascript.rhino.Node interfaceNode)
func (this_ *TypeCheck) priv_TypeCheck_reportInvalidInterfaceMemberDeclaration_5416(interfaceNode pkg_rhino.Node) {
	_ = interfaceNode
	var abstractMethodName string = this_.compiler.GetCodingConvention().GetAbstractMethodName()
	_ = abstractMethodName
	var abstractMethodMessage string
	_ = abstractMethodMessage
	if abstractMethodName != j2gort.NullString {
		abstractMethodMessage = (", or " + j2gort.StrValString(abstractMethodName))
	} else {
		abstractMethodMessage = ""
	}
	this_.compiler.Report_1(JSError_Make_3(interfaceNode, StaticGet_TypeCheck_INVALID_INTERFACE_MEMBER_DECLARATION(), []string{abstractMethodMessage}))
}

// Java method: TypeCheck.visitName(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *TypeCheck) VisitName(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	var isFunctionName bool = (n.GetParent().IsFunction() && n.IsFirstChildOf(parent))
	_ = isFunctionName
	if isFunctionName && (j2gort.CheckNotNullString(n.GetString()) == "") {
		return false
	}
	var reportIfMissingType bool = true
	_ = reportIfMissingType
	var parentNodeType pkg_rhino.Token = parent.GetToken()
	_ = parentNodeType
	if ((isFunctionName || (parentNodeType == pkg_rhino.StaticGet_Token_CATCH())) || (parentNodeType == pkg_rhino.StaticGet_Token_PARAM_LIST())) || NodeUtil_IsNameDeclaration(parent) {
		reportIfMissingType = false
	}
	if NodeUtil_IsEnhancedFor(parent) && (parent.GetFirstChild() == n) {
		reportIfMissingType = false
	}
	var type_ pkg_rhino.JSType = n.GetJSType()
	_ = type_
	if type_ == nil || type_.IsJavaNil_() {
		type_ = this_.priv_TypeCheck_getNativeType_5463(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
		var var_ *TypedVar = t.GetTypedScope().LookupTypedVar(n.GetString())
		_ = var_
		if var_ != nil {
			type_ = j2gort.CheckNotNullIface(var_.GetType())
		}
	}
	this_.priv_TypeCheck_ensureTyped_5461(n, type_)
	return reportIfMissingType
}

// Java method: TypeCheck.checkForOfTypes(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node forOf)
func (this_ *TypeCheck) priv_TypeCheck_checkForOfTypes_5418(t *NodeTraversal, forOf pkg_rhino.Node) {
	_ = t
	_ = forOf
	var lhs pkg_rhino.Node = forOf.GetFirstChild()
	_ = lhs
	var iterable pkg_rhino.Node = forOf.GetSecondChild()
	_ = iterable
	var iterableType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(iterable)
	_ = iterableType
	var actualType pkg_rhino.JSType = nil
	_ = actualType
	if forOf.IsForAwaitOf() {
		var maybeType *j2gort.Ext_java_util_Optional = this_.validator.ExpectAutoboxesToIterableOrAsyncIterable(iterable, iterableType, "Can only async iterate over a (non-null) Iterable or AsyncIterable type")
		_ = maybeType
		if maybeType.IsEmpty() {
			return
		}
		actualType = j2gort.AsIface[pkg_rhino.JSType](maybeType.Get())
	} else {
		this_.validator.ExpectAutoboxesToIterable(iterable, iterableType, "Can only iterate over a (non-null) Iterable type")
		actualType = iterableType.Autobox().GetTemplateTypeMap().GetResolvedTemplateType(this_.typeRegistry.GetIterableValueTemplate())
	}
	if NodeUtil_IsNameDeclaration(lhs) {
		lhs = lhs.GetFirstChild()
	}
	if lhs.IsDestructuringLhs() {
		lhs = lhs.GetFirstChild()
	}
	this_.priv_TypeCheck_checkCanAssignToWithScope_5396(t, forOf, lhs, actualType, lhs.GetJSDocInfo(), "declared type of for-of loop variable does not match inferred type")
}

// Java method: TypeCheck.visitGetProp(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *TypeCheck) priv_TypeCheck_visitGetProp_5419(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	var objNode pkg_rhino.Node = n.GetFirstChild()
	_ = objNode
	var childType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(objNode)
	_ = childType
	if childType.IsDict() {
		this_.priv_TypeCheck_report_5386(n, StaticGet_TypeValidator_ILLEGAL_PROPERTY_ACCESS(), []string{"'.'", "dict"})
	} else if this_.validator.ExpectNotNullOrUndefined(t, n, childType, "No properties on this expression", this_.priv_TypeCheck_getNativeType_5463(pkg_rhino.StaticGet_JSTypeNative_OBJECT_TYPE())) {
		this_.priv_TypeCheck_checkPropertyAccessForGetProp_5421(n)
	}
	this_.priv_TypeCheck_ensureTyped_5459(n)
}

// Java method: TypeCheck.visitOptChainGetProp(com.google.javascript.rhino.Node optChainGetProp)
func (this_ *TypeCheck) priv_TypeCheck_visitOptChainGetProp_5420(optChainGetProp pkg_rhino.Node) {
	_ = optChainGetProp
	var property pkg_rhino.Node = optChainGetProp.GetLastChild()
	_ = property
	var objNode pkg_rhino.Node = optChainGetProp.GetFirstChild()
	_ = objNode
	var childType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(objNode)
	_ = childType
	if childType.IsDict() {
		this_.priv_TypeCheck_report_5386(property, StaticGet_TypeValidator_ILLEGAL_PROPERTY_ACCESS(), []string{"'?.'", "dict"})
	} else if !(childType.IsUnknownType()) {
		this_.priv_TypeCheck_checkPropertyAccessForGetProp_5421(optChainGetProp)
	}
	this_.priv_TypeCheck_ensureTyped_5459(optChainGetProp)
}

// Java method: TypeCheck.checkPropertyAccessForGetProp(com.google.javascript.rhino.Node getProp)
func (this_ *TypeCheck) priv_TypeCheck_checkPropertyAccessForGetProp_5421(getProp pkg_rhino.Node) {
	_ = getProp
	if !(getProp.IsGetProp() || getProp.IsOptChainGetProp()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(getProp))
	}
	var objNode pkg_rhino.Node = getProp.GetFirstChild()
	_ = objNode
	var objType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(objNode)
	_ = objType
	var propType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(getProp)
	_ = propType
	this_.priv_TypeCheck_checkAbstractPropertyAccess_5423(getProp)
	this_.priv_TypeCheck_checkPropertyAccess_5424(objType, getProp, propType, objNode)
}

// Java method: TypeCheck.checkPropertyAccessForDestructuring(com.google.javascript.rhino.Node pattern, com.google.javascript.rhino.jstype.JSType objectType, com.google.javascript.rhino.Node stringKey, com.google.javascript.rhino.jstype.JSType inferredPropType)
func (this_ *TypeCheck) priv_TypeCheck_checkPropertyAccessForDestructuring_5422(pattern pkg_rhino.Node, objectType pkg_rhino.JSType, stringKey pkg_rhino.Node, inferredPropType pkg_rhino.JSType) {
	_ = pattern
	_ = objectType
	_ = stringKey
	_ = inferredPropType
	if !(pattern.IsDestructuringPattern()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(pattern))
	}
	if !(stringKey.IsStringKey()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(stringKey))
	}
	var objNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = objNode
	var patternParent pkg_rhino.Node = pattern.GetParent()
	_ = patternParent
	if (patternParent.IsAssign() || patternParent.IsDestructuringLhs()) && (pattern.GetNext() != 0) {
		objNode = pattern.GetNext()
	}
	this_.priv_TypeCheck_checkPropertyAccess_5424(objectType, stringKey, inferredPropType, objNode)
}

// Java method: TypeCheck.checkAbstractPropertyAccess(com.google.javascript.rhino.Node method)
func (this_ *TypeCheck) priv_TypeCheck_checkAbstractPropertyAccess_5423(method pkg_rhino.Node) {
	_ = method
	if NodeUtil_IsLhsOfAssign(method) {
		return
	}
	var methodType *pkg_rhino.FunctionType = this_.priv_TypeCheck_getJSType_5456(method).ToMaybeFunctionType()
	_ = methodType
	if ((methodType == nil) || !(methodType.IsAbstract())) || methodType.IsConstructor() {
		return
	}
	var objectNode pkg_rhino.Node = method.GetFirstChild()
	_ = objectNode
	if objectNode.IsSuper() {
		this_.priv_TypeCheck_report_5386(method, StaticGet_TypeCheck_ABSTRACT_SUPER_METHOD_NOT_USABLE(), []string{methodType.GetDisplayName()})
	} else if objectNode.IsGetProp() {
		var objectProp string = objectNode.GetString()
		_ = objectProp
		if (j2gort.CheckNotNullString(objectProp) == "prototype") || this_.compiler.GetCodingConvention().IsSuperClassReference(objectProp) {
			this_.priv_TypeCheck_report_5386(method, StaticGet_TypeCheck_ABSTRACT_SUPER_METHOD_NOT_USABLE(), []string{methodType.GetDisplayName()})
		}
	}
}

// Java method: TypeCheck.checkPropertyAccess(com.google.javascript.rhino.jstype.JSType childType, com.google.javascript.rhino.Node propNode, com.google.javascript.rhino.jstype.JSType propType, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node objNode)
func (this_ *TypeCheck) priv_TypeCheck_checkPropertyAccess_5424(childType pkg_rhino.JSType, propNode pkg_rhino.Node, propType pkg_rhino.JSType, objNode pkg_rhino.Node) {
	_ = childType
	_ = propNode
	_ = propType
	_ = objNode
	var isGetprop bool = NodeUtil_IsNormalOrOptChainGetProp(propNode)
	_ = isGetprop
	var propName string = propNode.GetString()
	_ = propName
	if propType.Equals(this_.typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())) {
		childType = childType.Autobox()
		var objectType pkg_rhino.ObjectType = pkg_rhino.ObjectType_Cast(childType)
		_ = objectType
		if objectType != nil && !objectType.IsJavaNil_() {
			if !(objectType.HasProperty_1_string(propName)) || objectType.Equals(this_.typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())) {
				if j2gort.InstanceOf[*pkg_rhino.EnumType](objectType) {
					this_.priv_TypeCheck_report_5386(propNode, StaticGet_TypeCheck_INEXISTENT_ENUM_ELEMENT(), []string{propName})
				} else {
					this_.priv_TypeCheck_checkPropertyAccessHelper_5426(j2gort.AsIface[pkg_rhino.JSType](objectType), propNode, objNode, false)
				}
			}
		} else {
			this_.priv_TypeCheck_checkPropertyAccessHelper_5426(childType, propNode, objNode, false)
		}
	} else if (childType.IsUnionType() && isGetprop) && !(this_.IsLValueGetProp(propNode)) {
		this_.priv_TypeCheck_checkPropertyAccessHelper_5426(childType, propNode, objNode, true)
	}
}

// Java method: TypeCheck.isLValueGetProp(com.google.javascript.rhino.Node n)
func (this_ *TypeCheck) IsLValueGetProp(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	return ((NodeUtil_IsUpdateOperator(parent) || NodeUtil_IsAssignmentOp(parent)) && (parent.GetFirstChild() == n))
}

// Java method: TypeCheck.checkPropertyAccessHelper(com.google.javascript.rhino.jstype.JSType objectType, com.google.javascript.rhino.Node propNode, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node objNode, boolean strictCheck)
func (this_ *TypeCheck) priv_TypeCheck_checkPropertyAccessHelper_5426(objectType pkg_rhino.JSType, propNode pkg_rhino.Node, objNode pkg_rhino.Node, strictCheck bool) {
	_ = objectType
	_ = propNode
	_ = objNode
	_ = strictCheck
	var isGetprop bool = NodeUtil_IsNormalOrOptChainGetProp(propNode)
	_ = isGetprop
	var propName string = propNode.GetString()
	_ = propName
	if (!(this_.reportMissingProperties) || objectType.IsEmptyType()) || (isGetprop && this_.priv_TypeCheck_allowStrictPropertyAccessOnNode_5428(propNode)) {
		return
	}
	var kind pkg_rhino.JSTypeRegistry_PropDefinitionKind = this_.typeRegistry.CanPropertyBeDefined(objectType, propName)
	_ = kind
	{
		var iifeCond_1 bool
		{

			re_18 := kind
			ae_19 := pkg_rhino.StaticGet_JSTypeRegistry_PropDefinitionKind_KNOWN()
			if re_18 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_18 == ae_19
		}
		if iifeCond_1 {
			return
		}
	}
	var isLooselyAssociated bool = (func() bool {
		re_20 := kind
		ae_21 := pkg_rhino.StaticGet_JSTypeRegistry_PropDefinitionKind_LOOSE()
		if re_20 == 0 {
			j2gort.ThrowNPE()
		}
		return re_20 == ae_21
	}() || func() bool {
		re_22 := kind
		ae_23 := pkg_rhino.StaticGet_JSTypeRegistry_PropDefinitionKind_LOOSE_UNION()
		if re_22 == 0 {
			j2gort.ThrowNPE()
		}
		return re_22 == ae_23
	}())
	_ = isLooselyAssociated
	var isUnknownType bool = objectType.IsUnknownType()
	_ = isUnknownType
	if isLooselyAssociated && isUnknownType {
		return
	}
	var isStruct bool = objectType.IsStruct()
	_ = isStruct
	var loosePropertyDeclaration bool = ((!(isStruct) && isGetprop) && this_.priv_TypeCheck_isQNameAssignmentTarget_5430(propNode))
	_ = loosePropertyDeclaration
	var maybePropExistenceCheck bool = ((!(isStruct) && isGetprop) && this_.priv_TypeCheck_allowLoosePropertyAccessOnNode_5429(propNode))
	_ = maybePropExistenceCheck
	var strictReport bool = (((strictCheck || isLooselyAssociated) || loosePropertyDeclaration) || maybePropExistenceCheck)
	_ = strictReport
	this_.priv_TypeCheck_reportMissingProperty_5427(objNode, objectType, propNode, kind, strictReport)
}

// Java method: TypeCheck.reportMissingProperty(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node objNode, com.google.javascript.rhino.jstype.JSType objectType, com.google.javascript.rhino.Node propNode, com.google.javascript.rhino.jstype.JSTypeRegistry.PropDefinitionKind kind, boolean strictReport)
func (this_ *TypeCheck) priv_TypeCheck_reportMissingProperty_5427(objNode pkg_rhino.Node, objectType pkg_rhino.JSType, propNode pkg_rhino.Node, kind pkg_rhino.JSTypeRegistry_PropDefinitionKind, strictReport bool) {
	_ = objNode
	_ = objectType
	_ = propNode
	_ = kind
	_ = strictReport
	var propName string = propNode.GetString()
	_ = propName
	var isObjectType bool = objectType.Equals(this_.priv_TypeCheck_getNativeType_5463(pkg_rhino.StaticGet_JSTypeNative_OBJECT_TYPE()))
	_ = isObjectType
	var lowConfidence bool = ((objectType.IsUnknownType() || objectType.IsAllType()) || isObjectType)
	_ = lowConfidence
	var isKnownToUnionMember bool = func() bool {
		re_24 := kind
		ae_25 := pkg_rhino.StaticGet_JSTypeRegistry_PropDefinitionKind_LOOSE_UNION()
		if re_24 == 0 {
			j2gort.ThrowNPE()
		}
		return re_24 == ae_25
	}()
	_ = isKnownToUnionMember
	var pair *TypeCheck_SuggestionPair = nil
	_ = pair
	if !(lowConfidence) && !(isKnownToUnionMember) {
		pair = TypeCheck_GetClosestPropertySuggestion(objectType, propName, ((j2gort.JStrLen(propName) - int32(1)) / int32(4)))
	}
	if pair != nil {
		var reportType *DiagnosticType = nil
		_ = reportType
		if strictReport {
			reportType = StaticGet_TypeCheck_STRICT_INEXISTENT_PROPERTY_WITH_SUGGESTION()
		} else {
			reportType = StaticGet_TypeCheck_INEXISTENT_PROPERTY_WITH_SUGGESTION()
		}
		this_.priv_TypeCheck_report_5386(propNode, reportType, []string{propName, func() string {
			if objNode != 0 {
				return this_.typeRegistry.GetReadableTypeName(objNode)
			}
			return objectType.String()
		}(), pair.suggestion})
	} else {
		var reportType_2 *DiagnosticType = nil
		_ = reportType_2
		if strictReport {
			if isKnownToUnionMember {
				reportType_2 = StaticGet_TypeCheck_STRICT_INEXISTENT_UNION_PROPERTY()
			} else {
				reportType_2 = StaticGet_TypeCheck_STRICT_INEXISTENT_PROPERTY()
			}
		} else if lowConfidence {
			reportType_2 = StaticGet_TypeCheck_POSSIBLE_INEXISTENT_PROPERTY()
		} else {
			reportType_2 = StaticGet_TypeCheck_INEXISTENT_PROPERTY()
		}
		this_.priv_TypeCheck_report_5386(propNode, reportType_2, []string{propName, func() string {
			if objNode != 0 {
				return this_.typeRegistry.GetReadableTypeName(objNode)
			}
			return objectType.String()
		}()})
	}
}

// Java method: TypeCheck.allowStrictPropertyAccessOnNode(com.google.javascript.rhino.Node n)
func (this_ *TypeCheck) priv_TypeCheck_allowStrictPropertyAccessOnNode_5428(n pkg_rhino.Node) bool {
	_ = n
	return n.GetParent().IsTypeOf()
}

// Java method: TypeCheck.allowLoosePropertyAccessOnNode(com.google.javascript.rhino.Node n)
func (this_ *TypeCheck) priv_TypeCheck_allowLoosePropertyAccessOnNode_5429(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	return ((NodeUtil_IsPropertyTest(this_.compiler, n) || NodeUtil_IsPropertyAbsenceTest(n)) || (n.IsQualifiedName() && parent.IsExprResult()))
}

// Java method: TypeCheck.isQNameAssignmentTarget(com.google.javascript.rhino.Node n)
func (this_ *TypeCheck) priv_TypeCheck_isQNameAssignmentTarget_5430(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	return ((n.IsQualifiedName() && parent.IsAssign()) && (parent.GetFirstChild() == n))
}

// Java static method: TypeCheck.getClosestPropertySuggestion(com.google.javascript.rhino.jstype.JSType objectType, java.lang.String propName, int maxDistance)
func TypeCheck_GetClosestPropertySuggestion(objectType pkg_rhino.JSType, propName string, maxDistance int32) *TypeCheck_SuggestionPair {
	clinit_TypeCheck()
	_ = objectType
	_ = propName
	_ = maxDistance
	return nil
}

// Java method: TypeCheck.visitGetElem(com.google.javascript.rhino.Node n)
func (this_ *TypeCheck) priv_TypeCheck_visitGetElem_5431(n pkg_rhino.Node) {
	_ = n
	this_.validator.ExpectIndexMatch(n, this_.priv_TypeCheck_getJSType_5456(n.GetFirstChild()), this_.priv_TypeCheck_getJSType_5456(n.GetLastChild()))
	this_.priv_TypeCheck_ensureTyped_5459(n)
}

// Java method: TypeCheck.visitOptChainGetElem(com.google.javascript.rhino.Node optChainGetElem)
func (this_ *TypeCheck) priv_TypeCheck_visitOptChainGetElem_5432(optChainGetElem pkg_rhino.Node) {
	_ = optChainGetElem
	var obj pkg_rhino.Node = optChainGetElem.GetFirstChild()
	_ = obj
	var objType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(obj)
	_ = objType
	if objType.IsNullType() || objType.IsVoidType() {
		this_.priv_TypeCheck_ensureTyped_5460(optChainGetElem, pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE())
	} else {
		this_.priv_TypeCheck_visitGetElem_5431(optChainGetElem)
	}
}

// Java method: TypeCheck.visitVar(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *TypeCheck) priv_TypeCheck_visitVar_5433(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if (n.GetParent().IsForOf() || n.GetParent().IsForIn()) || n.GetParent().IsForAwaitOf() {
		return
	}
	var varInfo *pkg_rhino.JSDocInfo
	_ = varInfo
	if n.HasOneChild() {
		varInfo = n.GetJSDocInfo()
	} else {
		varInfo = nil
	}
	{
		var child pkg_rhino.Node = n.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			if child.IsName() {
				var value pkg_rhino.Node = child.GetFirstChild()
				_ = value
				if value != 0 {
					var valueType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(value)
					_ = valueType
					var info *pkg_rhino.JSDocInfo = child.GetJSDocInfo()
					_ = info
					if info == nil {
						info = varInfo
					}
					this_.priv_TypeCheck_checkEnumAlias_5455(t, info, valueType, value)
					this_.priv_TypeCheck_checkCanAssignToWithScope_5396(t, value, child, valueType, info, "initializing variable")
				}
			} else {
				if !(child.IsDestructuringLhs()) {
					j2gort.CheckStateMsg(false, j2gort.BoxNode(child))
				}
				var name pkg_rhino.Node = child.GetFirstChild()
				_ = name
				var value_2 pkg_rhino.Node = child.GetSecondChild()
				_ = value_2
				var valueType_2 pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(value_2)
				_ = valueType_2
				this_.priv_TypeCheck_checkCanAssignToWithScope_5396(t, child, name, valueType_2, nil, "initializing variable")
			}
		}
	}
}

// Java method: TypeCheck.visitNew(com.google.javascript.rhino.Node n)
func (this_ *TypeCheck) priv_TypeCheck_visitNew_5434(n pkg_rhino.Node) {
	_ = n
	var constructor pkg_rhino.Node = n.GetFirstChild()
	_ = constructor
	var type_ pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(constructor).RestrictByNotNullOrUndefined()
	_ = type_
	if (!(this_.priv_TypeCheck_couldBeAConstructor_5435(type_)) || type_.Equals(this_.typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_SYMBOL_OBJECT_FUNCTION_TYPE()))) || type_.Equals(this_.typeRegistry.GetNativeType(pkg_rhino.StaticGet_JSTypeNative_BIGINT_OBJECT_FUNCTION_TYPE())) {
		this_.priv_TypeCheck_report_5386(n, StaticGet_TypeCheck_NOT_A_CONSTRUCTOR(), []string{type_.String()})
		this_.priv_TypeCheck_ensureTyped_5459(n)
		return
	}
	var fnType *pkg_rhino.FunctionType = type_.ToMaybeFunctionType()
	_ = fnType
	if (fnType != nil) && fnType.HasInstanceType() {
		if fnType.IsAbstract() {
			this_.priv_TypeCheck_report_5386(n, StaticGet_TypeCheck_INSTANTIATE_ABSTRACT_CLASS(), []string{})
		}
		this_.priv_TypeCheck_visitArgumentList_5446(n, fnType)
		var objType pkg_rhino.ObjectType = fnType.GetInstanceType()
		_ = objType
		this_.priv_TypeCheck_ensureTyped_5461(n, func() pkg_rhino.JSType {
			if objType != nil && !objType.IsJavaNil_() {
				return j2gort.AsIface[pkg_rhino.JSType](objType)
			}
			return this_.priv_TypeCheck_getNativeType_5463(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
		}())
	} else {
		this_.priv_TypeCheck_ensureTyped_5459(n)
	}
}

// Java method: TypeCheck.couldBeAConstructor(com.google.javascript.rhino.jstype.JSType type)
func (this_ *TypeCheck) priv_TypeCheck_couldBeAConstructor_5435(type_ pkg_rhino.JSType) bool {
	_ = type_
	return ((type_.IsConstructor() || type_.IsEmptyType()) || type_.IsUnknownType())
}

// Java method: TypeCheck.checkInterfaceConflictProperties(com.google.javascript.rhino.Node n, java.lang.String functionName, java.util.LinkedHashMap<com.google.javascript.rhino.jstype.Property.Key,com.google.javascript.rhino.jstype.ObjectType> properties, java.util.LinkedHashMap<com.google.javascript.rhino.jstype.Property.Key,com.google.javascript.rhino.jstype.ObjectType> currentProperties, com.google.javascript.rhino.jstype.ObjectType interfaceType)
func (this_ *TypeCheck) priv_TypeCheck_checkInterfaceConflictProperties_5436(n pkg_rhino.Node, functionName string, properties *j2gort.Ext_java_util_LinkedHashMap, currentProperties *j2gort.Ext_java_util_LinkedHashMap, interfaceType pkg_rhino.ObjectType) {
	_ = n
	_ = functionName
	_ = properties
	_ = currentProperties
	_ = interfaceType
	var implicitProto pkg_rhino.ObjectType = interfaceType.GetImplicitPrototype()
	_ = implicitProto
	var currentPropertyNames j2gort.Ext_com_google_common_collect_ImmutableSet = nil
	_ = currentPropertyNames
	if implicitProto == nil || implicitProto.IsJavaNil_() {
		currentPropertyNames = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()
	} else {
		currentPropertyNames = implicitProto.GetOwnPropertyKeys()
	}
	{
		cur_26 := j2gort.NewIterCursor(j2gort.CheckNotNull(currentPropertyNames))
		for cur_26.HasNext() {
			rawElem_27 := cur_26.Next()
			name := j2gort.AsIface[pkg_rhino.Property_Key](rawElem_27)
			_ = name
			var oType pkg_rhino.ObjectType = j2gort.AsIface[pkg_rhino.ObjectType](properties.Get(name))
			_ = oType
			currentProperties.Put_2(name, interfaceType)
			if oType != nil && !oType.IsJavaNil_() {
				var thisPropType pkg_rhino.JSType = interfaceType.GetPropertyType_1_Property_Key(name)
				_ = thisPropType
				var oPropType pkg_rhino.JSType = oType.GetPropertyType_1_Property_Key(name)
				_ = oPropType
				if (thisPropType.IsSubtypeOf_2(oPropType, this_.subtypingMode) || oPropType.IsSubtypeOf_2(thisPropType, this_.subtypingMode)) || ((thisPropType.IsFunctionType() && oPropType.IsFunctionType()) && thisPropType.ToMaybeFunctionType().HasEqualCallType(oPropType.ToMaybeFunctionType())) {
					continue
				}
				this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeCheck_INCOMPATIBLE_EXTENDED_PROPERTY_TYPE(), []string{functionName, name.HumanReadableName(), oType.String(), interfaceType.String()}))
			}
		}
	}
	{
		cur_28 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](interfaceType.GetCtorExtendedInterfaces())))
		for cur_28.HasNext() {
			rawElem_29 := cur_28.Next()
			iType := j2gort.AsIface[pkg_rhino.ObjectType](rawElem_29)
			_ = iType
			this_.priv_TypeCheck_checkInterfaceConflictProperties_5436(n, functionName, properties, currentProperties, iType)
		}
	}
}

// Java method: TypeCheck.visitFunction(com.google.javascript.rhino.Node n)
func (this_ *TypeCheck) priv_TypeCheck_visitFunction_5437(n pkg_rhino.Node) {
	_ = n
	if NodeUtil_IsEs6Constructor(n) {
		return
	}
	var functionType *pkg_rhino.FunctionType = pkg_rhino.JSType_ToMaybeFunctionType(n.GetJSType())
	_ = functionType
	if functionType.IsConstructor() {
		this_.priv_TypeCheck_checkConstructor_5439(n, functionType)
	} else if functionType.IsInterface() {
		this_.priv_TypeCheck_checkInterface_5442(n, functionType)
	} else if n.IsAsyncGeneratorFunction() {
		var returnType pkg_rhino.JSType = functionType.GetReturnType()
		_ = returnType
		this_.validator.ExpectAsyncGeneratorSupertype(n, returnType, "An async generator function must return a (supertype of) AsyncGenerator")
	} else if n.IsGeneratorFunction() {
		var returnType_2 pkg_rhino.JSType = functionType.GetReturnType()
		_ = returnType_2
		this_.validator.ExpectGeneratorSupertype(n, returnType_2, "A generator function must return a (supertype of) Generator")
	} else if n.IsAsyncFunction() {
		var returnType_3 pkg_rhino.JSType = functionType.GetReturnType()
		_ = returnType_3
		this_.validator.ExpectValidAsyncReturnType(n, returnType_3)
	}
}

// Java method: TypeCheck.visitClass(com.google.javascript.rhino.Node n)
func (this_ *TypeCheck) priv_TypeCheck_visitClass_5438(n pkg_rhino.Node) {
	_ = n
	var functionType *pkg_rhino.FunctionType = pkg_rhino.JSType_ToMaybeFunctionType(n.GetJSType())
	_ = functionType
	var extendsClause pkg_rhino.Node = n.GetSecondChild()
	_ = extendsClause
	if !(extendsClause.IsEmpty()) {
		var superType pkg_rhino.JSType = extendsClause.GetJSType()
		_ = superType
		if superType.IsConstructor() || superType.IsInterface() {
			this_.validator.ExpectExtends(n, functionType, superType.ToMaybeFunctionType())
		} else if !(superType.IsUnknownType()) {
			this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeCheck_CONFLICTING_EXTENDED_TYPE(), []string{func() string {
				if functionType.IsConstructor() {
					return "constructor"
				}
				return "interface"
			}(), this_.priv_TypeCheck_getBestFunctionName_5443(n)}))
		}
	}
	if functionType.IsConstructor() {
		this_.priv_TypeCheck_checkConstructor_5439(n, functionType)
	} else if functionType.IsInterface() {
		this_.priv_TypeCheck_checkInterface_5442(n, functionType)
	} else {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "CLASS node's type must be either constructor or interface: "+j2gort.StrVal(functionType))))
	}
}

// Java method: TypeCheck.checkConstructor(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.FunctionType functionType)
func (this_ *TypeCheck) priv_TypeCheck_checkConstructor_5439(n pkg_rhino.Node, functionType *pkg_rhino.FunctionType) {
	_ = n
	_ = functionType
	var baseConstructor *pkg_rhino.FunctionType = functionType.GetSuperClassConstructor()
	_ = baseConstructor
	if (!(j2gort.ObjEqIface[pkg_rhino.JSType](baseConstructor, this_.priv_TypeCheck_getNativeType_5463(pkg_rhino.StaticGet_JSTypeNative_OBJECT_FUNCTION_TYPE()))) && (baseConstructor != nil)) && baseConstructor.IsInterface() {
		this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeCheck_CONFLICTING_EXTENDED_TYPE(), []string{"constructor", this_.priv_TypeCheck_getBestFunctionName_5443(n)}))
	} else {
		if (((n.IsFunction() && (baseConstructor != nil)) && (baseConstructor.GetSource() != 0)) && baseConstructor.GetSource().IsClass()) && !(functionType.GetSource().IsClass()) {
			this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeCheck_ES5_CLASS_EXTENDING_ES6_CLASS(), []string{functionType.GetDisplayName(), baseConstructor.GetDisplayName()}))
		}
		if baseConstructor != nil {
			this_.priv_TypeCheck_checkStructDictSubtyping_5441(n, functionType, baseConstructor)
		}
		var alreadySeenInterfaces j2gort.CustomOrderedSet[pkg_rhino.JSType] = j2gort.NewCustomOrderedSet[pkg_rhino.JSType]()
		_ = alreadySeenInterfaces
		{
			cur_30 := j2gort.NewIterCursor(j2gort.CheckNotNull(functionType.GetOwnImplementedInterfaces()))
			for cur_30.HasNext() {
				rawElem_31 := cur_30.Next()
				baseInterface := j2gort.AsIface[pkg_rhino.JSType](rawElem_31)
				_ = baseInterface
				var badImplementedType bool = false
				_ = badImplementedType
				var baseInterfaceObj pkg_rhino.ObjectType = pkg_rhino.ObjectType_Cast(baseInterface)
				_ = baseInterfaceObj
				if baseInterfaceObj != nil && !baseInterfaceObj.IsJavaNil_() {
					var interfaceConstructor *pkg_rhino.FunctionType = baseInterfaceObj.GetConstructor()
					_ = interfaceConstructor
					if (interfaceConstructor != nil) && !(interfaceConstructor.IsInterface()) {
						badImplementedType = true
					}
				} else {
					badImplementedType = true
				}
				if badImplementedType {
					this_.priv_TypeCheck_report_5386(n, StaticGet_TypeCheck_BAD_IMPLEMENTED_TYPE(), []string{this_.priv_TypeCheck_getBestFunctionName_5443(n)})
				}
				baseInterface = this_.priv_TypeCheck_normalizeTemplatizedType_5440(baseInterface)
				if !(alreadySeenInterfaces.Add(baseInterface)) {
					this_.priv_TypeCheck_report_5386(n, StaticGet_TypeCheck_SAME_INTERFACE_MULTIPLE_IMPLEMENTS(), []string{baseInterface.String()})
				}
			}
		}
		this_.validator.ExpectAllInterfaceProperties(n, functionType)
		if !(functionType.IsAbstract()) {
			this_.validator.ExpectAbstractMethodsImplemented(n, functionType)
		}
	}
}

// Java method: TypeCheck.normalizeTemplatizedType(com.google.javascript.rhino.jstype.JSType maybeTemplatizedType)
func (this_ *TypeCheck) priv_TypeCheck_normalizeTemplatizedType_5440(maybeTemplatizedType pkg_rhino.JSType) pkg_rhino.JSType {
	_ = maybeTemplatizedType
	if !(maybeTemplatizedType.IsTemplatizedType()) {
		return maybeTemplatizedType
	}
	return j2gort.AsIface[pkg_rhino.JSType](maybeTemplatizedType.ToMaybeTemplatizedType().GetRawType())
}

// Java method: TypeCheck.checkStructDictSubtyping(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.FunctionType subtype, com.google.javascript.rhino.jstype.FunctionType supertype)
func (this_ *TypeCheck) priv_TypeCheck_checkStructDictSubtyping_5441(n pkg_rhino.Node, subtype *pkg_rhino.FunctionType, supertype *pkg_rhino.FunctionType) {
	_ = n
	_ = subtype
	_ = supertype
	if subtype.MakesDicts() && supertype.MakesStructs() {
		this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeCheck_DICT_EXTEND_STRUCT_TYPE(), []string{subtype.GetDisplayName(), supertype.GetDisplayName()}))
	} else if subtype.MakesStructs() && supertype.MakesDicts() {
		this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeCheck_STRUCT_EXTEND_DICT_TYPE(), []string{subtype.GetDisplayName(), supertype.GetDisplayName()}))
	}
}

// Java method: TypeCheck.checkInterface(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.FunctionType functionType)
func (this_ *TypeCheck) priv_TypeCheck_checkInterface_5442(n pkg_rhino.Node, functionType *pkg_rhino.FunctionType) {
	_ = n
	_ = functionType
	{
		cur_32 := j2gort.NewIterCursor(j2gort.CheckNotNull(functionType.GetExtendedInterfaces()))
		for cur_32.HasNext() {
			rawElem_33 := cur_32.Next()
			extInterface := j2gort.AsIface[pkg_rhino.ObjectType](rawElem_33)
			_ = extInterface
			if (extInterface.GetConstructor() != nil) && !(extInterface.GetConstructor().IsInterface()) {
				this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeCheck_CONFLICTING_EXTENDED_TYPE(), []string{"interface", this_.priv_TypeCheck_getBestFunctionName_5443(n)}))
			}
		}
	}
	if functionType.GetExtendedInterfacesCount() > int32(1) {
		var properties *j2gort.Ext_java_util_LinkedHashMap = j2gort.NewExt_java_util_LinkedHashMap()
		_ = properties
		var currentProperties *j2gort.Ext_java_util_LinkedHashMap = j2gort.NewExt_java_util_LinkedHashMap()
		_ = currentProperties
		{
			cur_34 := j2gort.NewIterCursor(j2gort.CheckNotNull(functionType.GetExtendedInterfaces()))
			for cur_34.HasNext() {
				rawElem_35 := cur_34.Next()
				interfaceType := j2gort.AsIface[pkg_rhino.ObjectType](rawElem_35)
				_ = interfaceType
				currentProperties.Clear()
				this_.priv_TypeCheck_checkInterfaceConflictProperties_5436(n, this_.priv_TypeCheck_getBestFunctionName_5443(n), properties, currentProperties, interfaceType)
				properties.PutAll(currentProperties)
			}
		}
	}
	var loopPath j2gort.Ext_java_util_List = functionType.CheckExtendsLoop_0()
	_ = loopPath
	if loopPath != nil {
		var strPath string = j2gort.CastString(loopPath.Stream().Map(lambdaSingleton_type_check_1).Collect_1(j2gort.Ext_java_util_stream_Collectors_Joining_1(" -> ")))
		_ = strPath
		this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeCheck_INTERFACE_EXTENDS_LOOP(), []string{j2gort.CastPtr[pkg_rhino.FunctionType](loopPath.Get(int32(0))).GetDisplayName(), strPath}))
	}
	this_.validator.ExpectAllInterfaceProperties(n, functionType)
}

// Java method: TypeCheck.getBestFunctionName(com.google.javascript.rhino.Node n)
func (this_ *TypeCheck) priv_TypeCheck_getBestFunctionName_5443(n pkg_rhino.Node) string {
	_ = n
	if !(n.IsClass() || n.IsFunction()) {
		j2gort.CheckState(false)
	}
	var name string = NodeUtil_GetBestLValueName(NodeUtil_GetBestLValue(n))
	_ = name
	if name != j2gort.NullString {
		return name
	} else {
		return ("<anonymous@" + j2gort.StrValString(n.GetSourceFileName()) + ":" + j2gort.StrValInt32(n.GetLineno()) + ">")
	}
}

// Java method: TypeCheck.checkCallConventions(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *TypeCheck) priv_TypeCheck_checkCallConventions_5444(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	var relationship *CodingConvention_SubclassRelationship = this_.compiler.GetCodingConvention().GetClassesDefinedByCall(n)
	_ = relationship
	var scope *TypedScope = t.GetTypedScope()
	_ = scope
	if relationship != nil {
		var superClass pkg_rhino.JSType = scope.LookupQualifiedName_1_string(relationship.superclassName)
		_ = superClass
		var superClassInstanceType pkg_rhino.ObjectType = TypeValidator_GetInstanceOfCtor(superClass)
		_ = superClassInstanceType
		var subClass pkg_rhino.JSType = scope.LookupQualifiedName_1_string(relationship.subclassName)
		_ = subClass
		var subClassInstance pkg_rhino.ObjectType = TypeValidator_GetInstanceOfCtor(subClass)
		_ = subClassInstance
		if ((((relationship.type_ == StaticGet_CodingConvention_SubclassType_INHERITS()) && (superClassInstanceType != nil && !superClassInstanceType.IsJavaNil_())) && !(superClassInstanceType.IsEmptyType())) && (subClassInstance != nil && !subClassInstance.IsJavaNil_())) && !(subClassInstance.IsEmptyType()) {
			if (((n.GetFirstChild().IsQualifiedName() && StaticGet_TypeCheck_GOOG_INHERITS().Matches(n.GetFirstChild())) && (subClass.ToMaybeFunctionType() != nil)) && (subClass.ToMaybeFunctionType().GetSource() != 0)) && subClass.ToMaybeFunctionType().GetSource().IsClass() {
				this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeCheck_ES6_CLASS_EXTENDING_CLASS_WITH_GOOG_INHERITS(), []string{}))
			}
			this_.validator.ExpectSuperType(n, superClassInstanceType, subClassInstance)
		}
	}
}

// Java method: TypeCheck.visitCall(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *TypeCheck) priv_TypeCheck_visitCall_5445(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	this_.priv_TypeCheck_checkCallConventions_5444(t, n)
	var child pkg_rhino.Node = n.GetFirstChild()
	_ = child
	var childType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(child).RestrictByNotNullOrUndefined()
	_ = childType
	if !(childType.CanBeCalled()) {
		this_.priv_TypeCheck_report_5386(n, StaticGet_TypeCheck_NOT_CALLABLE(), []string{childType.String()})
		this_.priv_TypeCheck_ensureTyped_5459(n)
		return
	}
	if childType.IsFunctionType() {
		var functionType *pkg_rhino.FunctionType = childType.ToMaybeFunctionType()
		_ = functionType
		if ((functionType.IsConstructor() && !(functionType.IsNativeObjectType())) && (functionType.GetReturnType().IsUnknownType() || functionType.GetReturnType().IsVoidType())) && !(n.GetFirstChild().IsSuper()) {
			this_.priv_TypeCheck_report_5386(n, StaticGet_TypeCheck_CONSTRUCTOR_NOT_CALLABLE(), []string{childType.String()})
		}
		if functionType.IsOrdinaryFunction() && !(NodeUtil_IsNormalOrOptChainGet(child)) {
			var receiverType pkg_rhino.JSType = functionType.GetTypeOfThis()
			_ = receiverType
			if ((receiverType.IsUnknownType() || receiverType.IsAllType()) || receiverType.IsVoidType()) || (receiverType.IsObjectType() && receiverType.ToObjectType().IsNativeObjectType()) {
			} else {
				this_.priv_TypeCheck_report_5386(n, StaticGet_TypeCheck_EXPECTED_THIS_TYPE(), []string{functionType.String()})
			}
		}
		this_.priv_TypeCheck_visitArgumentList_5446(n, functionType)
		this_.priv_TypeCheck_ensureTyped_5461(n, functionType.GetReturnType())
	} else {
		this_.priv_TypeCheck_ensureTyped_5459(n)
	}
}

// Java method: TypeCheck.visitArgumentList(com.google.javascript.rhino.Node call, com.google.javascript.rhino.jstype.FunctionType functionType)
func (this_ *TypeCheck) priv_TypeCheck_visitArgumentList_5446(call pkg_rhino.Node, functionType *pkg_rhino.FunctionType) {
	_ = call
	_ = functionType
	this_.priv_TypeCheck_checkArgumentsMatchParameters_5447(call, functionType, int32(0))
}

// Java method: TypeCheck.checkArgumentsMatchParameters(com.google.javascript.rhino.Node call, com.google.javascript.rhino.jstype.FunctionType functionType, int firstParameterIndex)
func (this_ *TypeCheck) priv_TypeCheck_checkArgumentsMatchParameters_5447(call pkg_rhino.Node, functionType *pkg_rhino.FunctionType, firstParameterIndex int32) {
	_ = call
	_ = functionType
	_ = firstParameterIndex
	var parameters []pkg_rhino.FunctionType_Parameter = functionType.GetParameterArray()
	_ = parameters
	var paramLen int32 = int32(len(parameters))
	_ = paramLen
	var paramIdx int32 = firstParameterIndex
	_ = paramIdx
	var isTaggedTemplate bool = call.IsTaggedTemplateLit()
	_ = isTaggedTemplate
	var argCursor pkg_rhino.Node = NodeUtil_FirstInvocationArgCursor(call, isTaggedTemplate)
	_ = argCursor
	var spreadArgumentCount int32 = int32(0)
	_ = spreadArgumentCount
	var normalArgumentCount int32 = firstParameterIndex
	_ = normalArgumentCount
	var checkArgumentTypeAgainstParameter bool = true
	_ = checkArgumentTypeAgainstParameter
	var parameter pkg_rhino.FunctionType_Parameter = (pkg_rhino.FunctionType_Parameter{})
	_ = parameter
	for argCursor != 0 {
		var argument pkg_rhino.Node = NodeUtil_InvocationArgFromCursor(argCursor, isTaggedTemplate)
		_ = argument
		argCursor = NodeUtil_NextInvocationArgCursor(argCursor, isTaggedTemplate)
		if argument.IsSpread() {
			spreadArgumentCount = spreadArgumentCount + 1
			checkArgumentTypeAgainstParameter = false
		} else {
			normalArgumentCount = normalArgumentCount + 1
		}
		if checkArgumentTypeAgainstParameter {
			if paramIdx < paramLen {
				parameter = parameters[func() int32 { old_36 := paramIdx; nw_37 := old_36 + 1; paramIdx = nw_37; return old_36 }()]
			} else if (parameter != (pkg_rhino.FunctionType_Parameter{})) && parameter.IsVariadic() {
			} else {
				parameter = (pkg_rhino.FunctionType_Parameter{})
				checkArgumentTypeAgainstParameter = false
			}
		}
		if checkArgumentTypeAgainstParameter {
			this_.validator.ExpectArgumentMatchesParameter(argument, this_.priv_TypeCheck_getJSType_5456(argument), parameter.GetJSType(), call, normalArgumentCount)
		}
	}
	var minArity int32 = functionType.GetMinArity()
	_ = minArity
	var maxArity int32 = functionType.GetMaxArity()
	_ = maxArity
	if spreadArgumentCount > int32(0) {
		if normalArgumentCount > maxArity {
			this_.priv_TypeCheck_report_5386(call, StaticGet_TypeCheck_WRONG_ARGUMENT_COUNT(), []string{this_.typeRegistry.GetReadableTypeNameNoDeref(call.GetFirstChild()), ("at least " + j2gort.StrValInt32(normalArgumentCount)), strconv.FormatInt(int64(minArity), 10), func() string {
				if maxArity == int32(2147483647) {
					return ""
				}
				return (" and no more than " + j2gort.StrValInt32(maxArity) + " argument(s)")
			}()})
		}
	} else {
		if (minArity > normalArgumentCount) || (maxArity < normalArgumentCount) {
			this_.priv_TypeCheck_report_5386(call, StaticGet_TypeCheck_WRONG_ARGUMENT_COUNT(), []string{this_.typeRegistry.GetReadableTypeNameNoDeref(call.GetFirstChild()), strconv.FormatInt(int64(normalArgumentCount), 10), strconv.FormatInt(int64(minArity), 10), func() string {
				if maxArity == int32(2147483647) {
					return ""
				}
				return (" and no more than " + j2gort.StrValInt32(maxArity) + " argument(s)")
			}()})
		}
	}
}

// Java method: TypeCheck.visitImplicitReturnExpression(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node exprNode)
func (this_ *TypeCheck) priv_TypeCheck_visitImplicitReturnExpression_5448(t *NodeTraversal, exprNode pkg_rhino.Node) {
	_ = t
	_ = exprNode
	var enclosingFunction pkg_rhino.Node = t.GetEnclosingFunction()
	_ = enclosingFunction
	var jsType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(enclosingFunction)
	_ = jsType
	if jsType.IsFunctionType() {
		var functionType *pkg_rhino.FunctionType = jsType.ToMaybeFunctionType()
		_ = functionType
		var expectedReturnType pkg_rhino.JSType = functionType.GetReturnType()
		_ = expectedReturnType
		if expectedReturnType == nil || expectedReturnType.IsJavaNil_() {
			expectedReturnType = this_.priv_TypeCheck_getNativeType_5463(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE())
		} else if enclosingFunction.IsAsyncFunction() {
			expectedReturnType = Promises_CreateAsyncReturnableType(this_.typeRegistry, expectedReturnType)
		}
		var actualReturnType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(exprNode)
		_ = actualReturnType
		this_.validator.ExpectCanAssignTo(exprNode, actualReturnType, expectedReturnType, "inconsistent return type")
	}
}

// Java method: TypeCheck.visitReturn(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *TypeCheck) priv_TypeCheck_visitReturn_5449(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	var enclosingFunction pkg_rhino.Node = t.GetEnclosingFunction()
	_ = enclosingFunction
	var jsType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(enclosingFunction)
	_ = jsType
	if jsType.IsFunctionType() {
		var functionType *pkg_rhino.FunctionType = jsType.ToMaybeFunctionType()
		_ = functionType
		var returnType pkg_rhino.JSType = functionType.GetReturnType()
		_ = returnType
		if returnType == nil || returnType.IsJavaNil_() {
			returnType = this_.priv_TypeCheck_getNativeType_5463(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE())
		} else if enclosingFunction.IsGeneratorFunction() {
			returnType = JsIterables_GetReturnElementType(returnType, this_.typeRegistry)
			if enclosingFunction.IsAsyncGeneratorFunction() {
				returnType = Promises_CreateAsyncReturnableType(this_.typeRegistry, Promises_WrapInIThenable(this_.typeRegistry, returnType))
			}
		} else if enclosingFunction.IsAsyncFunction() {
			returnType = Promises_CreateAsyncReturnableType(this_.typeRegistry, returnType)
		} else if returnType.IsVoidType() && functionType.IsConstructor() {
			if !(n.HasChildren()) {
				return
			}
			returnType = j2gort.AsIface[pkg_rhino.JSType](functionType.GetInstanceType())
		}
		var valueNode pkg_rhino.Node = n.GetFirstChild()
		_ = valueNode
		var actualReturnType pkg_rhino.JSType = nil
		_ = actualReturnType
		if valueNode == 0 {
			actualReturnType = this_.priv_TypeCheck_getNativeType_5463(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE())
			valueNode = n
		} else {
			actualReturnType = this_.priv_TypeCheck_getJSType_5456(valueNode)
		}
		this_.validator.ExpectCanAssignTo(valueNode, actualReturnType, returnType, "inconsistent return type")
	}
}

// Java method: TypeCheck.visitYield(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *TypeCheck) priv_TypeCheck_visitYield_5450(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	var jsType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(t.GetEnclosingFunction())
	_ = jsType
	var declaredYieldType pkg_rhino.JSType = this_.priv_TypeCheck_getNativeType_5463(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
	_ = declaredYieldType
	if jsType.IsFunctionType() {
		var functionType *pkg_rhino.FunctionType = jsType.ToMaybeFunctionType()
		_ = functionType
		var returnType pkg_rhino.JSType = functionType.GetReturnType()
		_ = returnType
		declaredYieldType = JsIterables_GetElementType(returnType, this_.typeRegistry)
		if t.GetEnclosingFunction().IsAsyncGeneratorFunction() {
			declaredYieldType = Promises_CreateAsyncReturnableType(this_.typeRegistry, Promises_WrapInIThenable(this_.typeRegistry, declaredYieldType))
		}
	}
	var valueNode pkg_rhino.Node = n.GetFirstChild()
	_ = valueNode
	var actualYieldType pkg_rhino.JSType = nil
	_ = actualYieldType
	if valueNode == 0 {
		actualYieldType = this_.priv_TypeCheck_getNativeType_5463(pkg_rhino.StaticGet_JSTypeNative_VOID_TYPE())
		valueNode = n
	} else {
		actualYieldType = this_.priv_TypeCheck_getJSType_5456(valueNode)
	}
	if n.IsYieldAll() {
		if t.GetEnclosingFunction().IsAsyncGeneratorFunction() {
			var maybeActualYieldType *j2gort.Ext_java_util_Optional = this_.validator.ExpectAutoboxesToIterableOrAsyncIterable(n, actualYieldType, "Expression yield* expects an iterable or async iterable")
			_ = maybeActualYieldType
			if maybeActualYieldType.IsEmpty() {
				return
			}
			actualYieldType = j2gort.AsIface[pkg_rhino.JSType](maybeActualYieldType.Get())
		} else {
			if !(this_.validator.ExpectAutoboxesToIterable(n, actualYieldType, "Expression yield* expects an iterable")) {
				return
			}
			var templateType *pkg_rhino.TemplateType = this_.typeRegistry.GetIterableValueTemplate()
			_ = templateType
			actualYieldType = actualYieldType.Autobox().GetTemplateTypeMap().GetResolvedTemplateType(templateType)
		}
	}
	this_.validator.ExpectCanAssignTo(valueNode, actualYieldType, declaredYieldType, "Yielded type does not match declared return type.")
}

// Java method: TypeCheck.visitTaggedTemplateLit(com.google.javascript.rhino.Node n)
func (this_ *TypeCheck) priv_TypeCheck_visitTaggedTemplateLit_5451(n pkg_rhino.Node) {
	_ = n
	var tag pkg_rhino.Node = n.GetFirstChild()
	_ = tag
	var tagType pkg_rhino.JSType = tag.GetJSType().RestrictByNotNullOrUndefined()
	_ = tagType
	if !(tagType.CanBeCalled()) {
		this_.priv_TypeCheck_report_5386(n, StaticGet_TypeCheck_NOT_CALLABLE(), []string{tagType.String()})
		return
	} else if !(tagType.IsFunctionType()) {
		return
	}
	var tagFnType *pkg_rhino.FunctionType = tagType.ToMaybeFunctionType()
	_ = tagFnType
	var parameters []pkg_rhino.FunctionType_Parameter = tagFnType.GetParameterArray()
	_ = parameters
	if int32(len(parameters)) == int32(0) {
		this_.priv_TypeCheck_report_5386(n, StaticGet_TypeCheck_WRONG_ARGUMENT_COUNT(), []string{this_.typeRegistry.GetReadableTypeNameNoDeref(tag), strconv.FormatInt(int64(NodeUtil_GetInvocationArgsCount(n)), 10), "0", " and no more than 0 argument(s)"})
		return
	}
	var firstParameter pkg_rhino.FunctionType_Parameter = parameters[int32(0)]
	_ = firstParameter
	var parameterType pkg_rhino.JSType = firstParameter.GetJSType().RestrictByNotNullOrUndefined()
	_ = parameterType
	if parameterType != nil && !parameterType.IsJavaNil_() {
		this_.validator.ExpectITemplateArraySupertype(tag, parameterType, "Invalid type for the first parameter of tag function")
	}
	this_.priv_TypeCheck_checkArgumentsMatchParameters_5447(n, tagFnType, int32(1))
}

// Java method: TypeCheck.visitBinaryOperator(com.google.javascript.rhino.Token op, com.google.javascript.rhino.Node n)
func (this_ *TypeCheck) priv_TypeCheck_visitBinaryOperator_5452(op pkg_rhino.Token, n pkg_rhino.Node) {
	_ = op
	_ = n
	var operatorType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(n)
	_ = operatorType
	var left pkg_rhino.Node = n.GetFirstChild()
	_ = left
	var leftType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(left)
	_ = leftType
	var right pkg_rhino.Node = n.GetLastChild()
	_ = right
	var rightType pkg_rhino.JSType = this_.priv_TypeCheck_getJSType_5456(right)
	_ = rightType
	if operatorType.IsNoType() {
		this_.priv_TypeCheck_report_5386(n, StaticGet_TypeCheck_BINARY_OPERATION(), []string{NodeUtil_OpToStr(n.GetToken()), leftType.String(), rightType.String()})
		return
	}
	switch (op).Ordinal() {
	case int32(55) /* Token.ASSIGN_LSH */, int32(56) /* Token.ASSIGN_RSH */, int32(10) /* Token.LSH */, int32(11) /* Token.RSH */, int32(57) /* Token.ASSIGN_URSH */, int32(12) /* Token.URSH */ :
		if operatorType.IsNumber() {
			var opStr string = NodeUtil_OpToStr(n.GetToken())
			_ = opStr
			if !(leftType.MatchesNumberContext()) {
				this_.priv_TypeCheck_report_5386(left, StaticGet_TypeCheck_BIT_OPERATION(), []string{opStr, leftType.String()})
			} else {
				this_.validator.ExpectNumberStrict(n, leftType, ("operator " + j2gort.StrValString(opStr)))
			}
			if !(rightType.MatchesNumberContext()) {
				this_.priv_TypeCheck_report_5386(right, StaticGet_TypeCheck_BIT_OPERATION(), []string{opStr, rightType.String()})
			} else {
				this_.validator.ExpectNumberStrict(n, rightType, ("operator " + j2gort.StrValString(opStr)))
			}
		} else {
			this_.validator.ExpectBigIntOrNumber(left, leftType, "left operand")
			this_.validator.ExpectBigIntOrNumber(right, rightType, "right operand")
		}
	case int32(61) /* Token.ASSIGN_DIV */, int32(62) /* Token.ASSIGN_MOD */, int32(60) /* Token.ASSIGN_MUL */, int32(59) /* Token.ASSIGN_SUB */, int32(63) /* Token.ASSIGN_EXPONENT */, int32(16) /* Token.DIV */, int32(17) /* Token.MOD */, int32(15) /* Token.MUL */, int32(14) /* Token.SUB */, int32(18) /* Token.EXPONENT */ :
		if operatorType.IsNumber() {
			this_.validator.ExpectNumber(left, leftType, "left operand")
			this_.validator.ExpectNumber(right, rightType, "right operand")
		} else {
			this_.validator.ExpectBigIntOrNumber(left, leftType, "left operand")
			this_.validator.ExpectBigIntOrNumber(right, rightType, "right operand")
		}
	case int32(54) /* Token.ASSIGN_BITAND */, int32(53) /* Token.ASSIGN_BITXOR */, int32(52) /* Token.ASSIGN_BITOR */, int32(3) /* Token.BITAND */, int32(2) /* Token.BITXOR */, int32(1) /* Token.BITOR */ :
		if operatorType.IsNumber() {
			this_.validator.ExpectBitwiseable(left, leftType, "bad left operand to bitwise operator")
			this_.validator.ExpectBitwiseable(right, rightType, "bad right operand to bitwise operator")
		} else {
			this_.validator.ExpectBigIntOrNumber(left, leftType, "bad left operand to bitwise operator")
			this_.validator.ExpectBigIntOrNumber(right, rightType, "bad right operand to bitwise operator")
		}
	case int32(58) /* Token.ASSIGN_ADD */, int32(13) /* Token.ADD */ :
	default:
		this_.priv_TypeCheck_report_5386(n, StaticGet_TypeCheck_UNEXPECTED_TOKEN(), []string{j2gort.CheckNotNullEnum(op).String()})
	}
	this_.priv_TypeCheck_ensureTyped_5459(n)
}

// Java method: TypeCheck.visitModuleBody(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node moduleBody)
func (this_ *TypeCheck) priv_TypeCheck_visitModuleBody_5453(t *NodeTraversal, moduleBody pkg_rhino.Node) {
	_ = t
	_ = moduleBody
	var associatedModule *Module = ModuleImportResolver_GetModuleFromScopeRoot(this_.compiler.GetModuleMap(), NewLambda_CompilerInputProvider(func(b *pkg_rhino.InputId) *CompilerInput {
		_ = b
		return t.GetInput()
	}), moduleBody)
	_ = associatedModule
	if !(associatedModule.Metadata().IsLegacyGoogModule()) {
		return
	}
	var moduleName pkg_rhino.QualifiedName = pkg_rhino.QualifiedName_Of(associatedModule.ClosureNamespace())
	_ = moduleName
	var googModuleCall pkg_rhino.Node = moduleBody.GetFirstChild()
	_ = googModuleCall
	if moduleName.IsSimple() {
		var globalVar *TypedVar = this_.topScope.LookupTypedVar(moduleName.Join_0())
		_ = globalVar
		this_.validator.ExpectCanAssignTo(googModuleCall, moduleBody.GetJSType(), globalVar.GetType(), "legacy goog.module export")
	} else {
		var parentType pkg_rhino.JSType = this_.topScope.LookupQualifiedName_1_QualifiedName(moduleName.GetOwner())
		_ = parentType
		var parentObjectType pkg_rhino.ObjectType
		_ = parentObjectType
		if parentType != nil && !parentType.IsJavaNil_() {
			parentObjectType = parentType.ToMaybeObjectType()
		} else {
			parentObjectType = nil
		}
		if parentObjectType == nil || parentObjectType.IsJavaNil_() {
			return
		}
		this_.validator.ExpectCanAssignToPropertyOf_6_Node_JSType_JSType_JSType_QualifiedName_string(googModuleCall, moduleBody.GetJSType(), parentObjectType.GetPropertyType_1_string(moduleName.GetComponent()), j2gort.AsIface[pkg_rhino.JSType](parentObjectType), moduleName.GetOwner(), moduleName.GetComponent())
	}
}

// Java method: TypeCheck.visitDynamicImport(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node dynamicImport)
func (this_ *TypeCheck) priv_TypeCheck_visitDynamicImport_5454(t *NodeTraversal, dynamicImport pkg_rhino.Node) {
	_ = t
	_ = dynamicImport
	this_.priv_TypeCheck_ensureTyped_5461(dynamicImport, this_.priv_TypeCheck_getNativeType_5463(pkg_rhino.StaticGet_JSTypeNative_PROMISE_TYPE()))
	var importSpecifier pkg_rhino.Node = dynamicImport.GetFirstChild()
	_ = importSpecifier
	var importSpecifierType pkg_rhino.JSType = importSpecifier.GetJSType()
	_ = importSpecifierType
	if importSpecifierType == nil || importSpecifierType.IsJavaNil_() {
		this_.priv_TypeCheck_ensureTyped_5460(importSpecifier, pkg_rhino.StaticGet_JSTypeNative_STRING_TYPE())
	} else {
		this_.validator.ExpectNotNullOrUndefined(t, importSpecifier, importSpecifierType, "dynamic import specifier", this_.priv_TypeCheck_getNativeType_5463(pkg_rhino.StaticGet_JSTypeNative_STRING_TYPE()))
	}
}

// Java method: TypeCheck.checkEnumAlias(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.JSDocInfo declInfo, com.google.javascript.rhino.jstype.JSType valueType, com.google.javascript.rhino.Node nodeToWarn)
func (this_ *TypeCheck) priv_TypeCheck_checkEnumAlias_5455(t *NodeTraversal, declInfo *pkg_rhino.JSDocInfo, valueType pkg_rhino.JSType, nodeToWarn pkg_rhino.Node) {
	_ = t
	_ = declInfo
	_ = valueType
	_ = nodeToWarn
	if (declInfo == nil) || !(declInfo.HasEnumParameterType()) {
		return
	}
	if !(valueType.IsEnumType()) {
		return
	}
	var valueEnumType *pkg_rhino.EnumType = valueType.ToMaybeEnumType()
	_ = valueEnumType
	var valueEnumPrimitiveType pkg_rhino.JSType = valueEnumType.GetElementsType().GetPrimitiveType()
	_ = valueEnumPrimitiveType
	this_.validator.ExpectCanAssignTo(nodeToWarn, valueEnumPrimitiveType, declInfo.GetEnumParameterType().Evaluate(pkg_rhino.ToIface_StaticTypedScope(t.GetTypedScope()), this_.typeRegistry), "incompatible enum element types")
}

// Java method: TypeCheck.getJSType(com.google.javascript.rhino.Node n)
func (this_ *TypeCheck) priv_TypeCheck_getJSType_5456(n pkg_rhino.Node) pkg_rhino.JSType {
	_ = n
	var jsType pkg_rhino.JSType = n.GetJSType()
	_ = jsType
	if jsType == nil || jsType.IsJavaNil_() {
		return this_.priv_TypeCheck_getNativeType_5463(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
	} else {
		return jsType
	}
}

// Java method: TypeCheck.getPropertyTypeIfDeclared(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable ObjectType objectType, java.lang.String propertyName)
func (this_ *TypeCheck) priv_TypeCheck_getPropertyTypeIfDeclared_5457(objectType pkg_rhino.ObjectType, propertyName string) pkg_rhino.JSType {
	_ = objectType
	_ = propertyName
	j2gort.CheckNotNullString(propertyName)
	if ((objectType != nil && !objectType.IsJavaNil_()) && objectType.HasProperty_1_string(propertyName)) && !(objectType.IsPropertyTypeInferred_1_string(propertyName)) {
		return objectType.GetPropertyType_1_string(propertyName)
	}
	return this_.priv_TypeCheck_getNativeType_5463(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
}

// Java method: TypeCheck.getPropertyTypeIfDeclared(com.google.javascript.rhino.jstype.@org.jspecify.annotations.Nullable ObjectType objectType, com.google.javascript.rhino.jstype.Property.Key propertyName)
func (this_ *TypeCheck) priv_TypeCheck_getPropertyTypeIfDeclared_5458(objectType pkg_rhino.ObjectType, propertyName pkg_rhino.Property_Key) pkg_rhino.JSType {
	_ = objectType
	_ = propertyName
	if ((objectType != nil && !objectType.IsJavaNil_()) && objectType.HasProperty_1_Property_Key(propertyName)) && !(objectType.IsPropertyTypeInferred_1_Property_Key(propertyName)) {
		return objectType.GetPropertyType_1_Property_Key(propertyName)
	}
	return this_.priv_TypeCheck_getNativeType_5463(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE())
}

// Java method: TypeCheck.ensureTyped(com.google.javascript.rhino.Node n)
func (this_ *TypeCheck) priv_TypeCheck_ensureTyped_5459(n pkg_rhino.Node) {
	_ = n
	this_.priv_TypeCheck_ensureTyped_5461(n, this_.priv_TypeCheck_getNativeType_5463(pkg_rhino.StaticGet_JSTypeNative_UNKNOWN_TYPE()))
}

// Java method: TypeCheck.ensureTyped(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSTypeNative type)
func (this_ *TypeCheck) priv_TypeCheck_ensureTyped_5460(n pkg_rhino.Node, type_ pkg_rhino.JSTypeNative) {
	_ = n
	_ = type_
	this_.priv_TypeCheck_ensureTyped_5461(n, this_.priv_TypeCheck_getNativeType_5463(type_))
}

// Java method: TypeCheck.ensureTyped(com.google.javascript.rhino.Node n, com.google.javascript.rhino.jstype.JSType type)
func (this_ *TypeCheck) priv_TypeCheck_ensureTyped_5461(n pkg_rhino.Node, type_ pkg_rhino.JSType) {
	_ = n
	_ = type_
	if !((!(n.IsFunction()) || type_.IsFunctionType()) || type_.IsUnknownType()) {
		j2gort.CheckState(false)
	}
	if j2gort.IfaceEquals[pkg_rhino.JSType](n.GetJSType(), nil) {
		n.SetJSType(type_)
	}
}

// Java method: TypeCheck.getTypedPercent()
func (this_ *TypeCheck) GetTypedPercent() float64 {
	var total int32 = ((this_.nullCount + this_.unknownCount) + this_.typedCount)
	_ = total
	if total == int32(0) {
		return float64(0.0)
	} else {
		return ((float64(100.0) * float64(this_.typedCount)) / float64(total))
	}
}

// Java method: TypeCheck.getNativeType(com.google.javascript.rhino.jstype.JSTypeNative typeId)
func (this_ *TypeCheck) priv_TypeCheck_getNativeType_5463(typeId pkg_rhino.JSTypeNative) pkg_rhino.JSType {
	_ = typeId
	return this_.typeRegistry.GetNativeType(typeId)
}

// Java method: TypeCheck.checkJsdocInfoContainsObjectWithBadKey(com.google.javascript.rhino.Node n)
func (this_ *TypeCheck) priv_TypeCheck_checkJsdocInfoContainsObjectWithBadKey_5464(n pkg_rhino.Node) {
	_ = n
	if n.GetJSDocInfo() != nil {
		var info *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
		_ = info
		this_.priv_TypeCheck_checkTypeContainsObjectWithBadKey_5465(n, info.GetType())
		this_.priv_TypeCheck_checkTypeContainsObjectWithBadKey_5465(n, info.GetReturnType())
		this_.priv_TypeCheck_checkTypeContainsObjectWithBadKey_5465(n, info.GetTypedefType())
		{
			var i int32 = int32(0)
			_ = i
			var paramCount int32 = info.GetParameterCount()
			_ = paramCount
			for ; i < paramCount; i = i + 1 {
				this_.priv_TypeCheck_checkTypeContainsObjectWithBadKey_5465(n, info.GetParameterTypeAt(i))
			}
		}
	}
}

// Java method: TypeCheck.checkTypeContainsObjectWithBadKey(com.google.javascript.rhino.Node n, com.google.javascript.rhino.JSTypeExpression type)
func (this_ *TypeCheck) priv_TypeCheck_checkTypeContainsObjectWithBadKey_5465(n pkg_rhino.Node, type_ pkg_rhino.JSTypeExpression) {
	_ = n
	_ = type_
	if (type_ != (pkg_rhino.JSTypeExpression{})) && !(j2gort.IfaceEquals[pkg_rhino.JSType](type_.GetRoot().GetJSType(), nil)) {
		var realType pkg_rhino.JSType = type_.GetRoot().GetJSType()
		_ = realType
		if !(realType.MayHaveNonStringifiableKey()) {
			return
		}
		var objectWithBadKey pkg_rhino.JSType = nil
		_ = objectWithBadKey
		func() {
			defer func() {
				this_.nonStringifiableKeyVisited.Clear()
			}()
			objectWithBadKey = this_.priv_TypeCheck_findObjectWithNonStringifiableKey_5468(realType)
		}()
		if objectWithBadKey != nil && !objectWithBadKey.IsJavaNil_() {
			this_.compiler.Report_1(JSError_Make_3(n, StaticGet_TypeCheck_NON_STRINGIFIABLE_OBJECT_KEY(), []string{objectWithBadKey.String()}))
		}
	}
}

// Java method: TypeCheck.isReasonableObjectPropertyKey(com.google.javascript.rhino.jstype.JSType type)
func (this_ *TypeCheck) priv_TypeCheck_isReasonableObjectPropertyKey_5466(type_ pkg_rhino.JSType) bool {
	_ = type_
	if ((((((((((type_.IsUnknownType() || type_.IsNumber()) || type_.IsString()) || type_.IsSymbol()) || type_.IsBooleanObjectType()) || type_.IsBooleanValueType()) || type_.IsDateType()) || type_.IsRegexpType()) || type_.IsInterface()) || type_.IsRecordType()) || type_.IsNullType()) || type_.IsVoidType() {
		return true
	}
	if type_.ToMaybeEnumElementType() != nil {
		return this_.priv_TypeCheck_isReasonableObjectPropertyKey_5466(type_.ToMaybeEnumElementType().GetPrimitiveType())
	}
	if type_.IsArrayType() {
		return true
	}
	if type_.IsTemplatizedType() {
		var templatizedType *pkg_rhino.TemplatizedType = type_.ToMaybeTemplatizedType()
		_ = templatizedType
		if templatizedType.GetReferencedType().IsArrayType() {
			return this_.priv_TypeCheck_isReasonableObjectPropertyKey_5466(j2gort.CheckNotNullSlice(j2gort.CheckNotNullPtr(templatizedType).Slice_GetTemplateTypes())[int32(0)])
		}
	}
	var namedType *pkg_rhino.NamedType
	_ = namedType
	var ok_38 bool
	namedType, ok_38 = j2gort.TryCast[*pkg_rhino.NamedType](type_)
	if ok_38 {
		return this_.priv_TypeCheck_isReasonableObjectPropertyKey_5466(namedType.GetReferencedType())
	}
	if type_.IsUnionType() {
		var union *pkg_rhino.UnionType = type_.ToMaybeUnionType()
		_ = union
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = union.GetAlternatesCount()
			_ = n
			for ; i < n; i = i + 1 {
				if !(this_.priv_TypeCheck_isReasonableObjectPropertyKey_5466(union.GetAlternateAt(i))) {
					return false
				}
			}
		}
		return true
	}
	if type_.IsObject() {
		var objectType pkg_rhino.ObjectType = type_.ToMaybeObjectType()
		_ = objectType
		var constructor pkg_rhino.JSType = pkg_rhino.ToIface_JSType(objectType.GetConstructor())
		_ = constructor
		if (constructor != nil && !constructor.IsJavaNil_()) && constructor.IsInterface() {
			return true
		}
		return this_.priv_TypeCheck_classHasToString_5469(objectType)
	}
	return false
}

// Java method: TypeCheck.isObjectTypeWithNonStringifiableKey(com.google.javascript.rhino.jstype.JSType type)
func (this_ *TypeCheck) priv_TypeCheck_isObjectTypeWithNonStringifiableKey_5467(type_ pkg_rhino.JSType) bool {
	_ = type_
	if !(type_.IsTemplatizedType()) {
		return false
	}
	var templateTypeMap *pkg_rhino.TemplateTypeMap = type_.GetTemplateTypeMap()
	_ = templateTypeMap
	var objectIndexKey *pkg_rhino.TemplateType = this_.typeRegistry.GetObjectIndexKey()
	_ = objectIndexKey
	if templateTypeMap.HasTemplateKey_1(objectIndexKey) {
		return !(this_.priv_TypeCheck_isReasonableObjectPropertyKey_5466(templateTypeMap.GetResolvedTemplateType(objectIndexKey)))
	} else {
		return false
	}
}

// Java method: TypeCheck.findObjectWithNonStringifiableKey(com.google.javascript.rhino.jstype.JSType type)
func (this_ *TypeCheck) priv_TypeCheck_findObjectWithNonStringifiableKey_5468(type_ pkg_rhino.JSType) pkg_rhino.JSType {
	_ = type_
	if !(type_.MayHaveNonStringifiableKey()) {
		return nil
	}
	if !(this_.nonStringifiableKeyVisited.Add(type_)) {
		return nil
	}
	if this_.priv_TypeCheck_isObjectTypeWithNonStringifiableKey_5467(type_) {
		return type_
	}
	if type_.IsUnionType() {
		var union *pkg_rhino.UnionType = type_.ToMaybeUnionType()
		_ = union
		{
			var i int32 = int32(0)
			_ = i
			var n int32 = union.GetAlternatesCount()
			_ = n
			for ; i < n; i = i + 1 {
				var result pkg_rhino.JSType = this_.priv_TypeCheck_findObjectWithNonStringifiableKey_5468(union.GetAlternateAt(i))
				_ = result
				if result != nil && !result.IsJavaNil_() {
					return result
				}
			}
		}
	}
	if type_.IsTemplatizedType() {
		var templateTypes []pkg_rhino.JSType = j2gort.CheckNotNullPtr(type_.ToMaybeTemplatizedType()).Slice_GetTemplateTypes()
		_ = templateTypes
		{
			var i_2 int32 = int32(0)
			_ = i_2
			var n_2 int32 = int32(len(templateTypes))
			_ = n_2
			for ; i_2 < n_2; i_2 = i_2 + 1 {
				var result_2 pkg_rhino.JSType = this_.priv_TypeCheck_findObjectWithNonStringifiableKey_5468(templateTypes[i_2])
				_ = result_2
				if result_2 != nil && !result_2.IsJavaNil_() {
					return result_2
				}
			}
		}
	}
	if type_.IsOrdinaryFunction() {
		var function *pkg_rhino.FunctionType = type_.ToMaybeFunctionType()
		_ = function
		var params []pkg_rhino.FunctionType_Parameter = function.GetParameterArray()
		_ = params
		{
			var i_3 int32 = int32(0)
			_ = i_3
			var len_ int32 = int32(len(params))
			_ = len_
			for ; i_3 < len_; i_3 = i_3 + 1 {
				var result_3 pkg_rhino.JSType = this_.priv_TypeCheck_findObjectWithNonStringifiableKey_5468(params[i_3].GetJSType())
				_ = result_3
				if result_3 != nil && !result_3.IsJavaNil_() {
					return result_3
				}
			}
		}
		return this_.priv_TypeCheck_findObjectWithNonStringifiableKey_5468(function.GetReturnType())
	}
	return nil
}

// Java method: TypeCheck.classHasToString(com.google.javascript.rhino.jstype.ObjectType type)
func (this_ *TypeCheck) priv_TypeCheck_classHasToString_5469(type_ pkg_rhino.ObjectType) bool {
	_ = type_
	var toStringProperty *pkg_rhino.Property = type_.GetOwnSlot_1_string("toString")
	_ = toStringProperty
	if toStringProperty != nil {
		return toStringProperty.GetType().IsFunctionType()
	}
	var parent pkg_rhino.ObjectType = type_.GetImplicitPrototype()
	_ = parent
	if (parent != nil && !parent.IsJavaNil_()) && !(parent.IsNativeObjectType()) {
		return this_.priv_TypeCheck_classHasToString_5469(parent)
	}
	return false
}

// Java method: TypeCheck.checkNotReadonlyPropertyAssignment(com.google.javascript.rhino.Node lhs)
func (this_ *TypeCheck) priv_TypeCheck_checkNotReadonlyPropertyAssignment_5470(lhs pkg_rhino.Node) {
	_ = lhs
	if !(lhs.IsGetProp()) && !(lhs.IsGetElem()) {
		return
	}
	var lhsType pkg_rhino.JSType = lhs.GetFirstChild().GetJSType()
	_ = lhsType
	if lhsType == nil || lhsType.IsJavaNil_() {
		return
	}
	var roArray pkg_rhino.JSType = this_.priv_TypeCheck_getNativeType_5463(pkg_rhino.StaticGet_JSTypeNative_READONLY_ARRAY_TYPE())
	_ = roArray
	var union *pkg_rhino.UnionType = lhsType.ToMaybeUnionType()
	_ = union
	var altCount int32
	_ = altCount
	if union != nil {
		altCount = union.GetAlternatesCount()
	} else {
		altCount = int32(1)
	}
	{
		var i int32 = int32(0)
		_ = i
		for ; i < altCount; i = i + 1 {
			var type_ pkg_rhino.JSType
			_ = type_
			if union != nil {
				type_ = union.GetAlternateAt(i)
			} else {
				type_ = lhsType
			}
			if roArray.Equals(type_) || roArray.Equals(TypeCheck_MaybeReferencedType(type_)) {
				this_.compiler.Report_1(JSError_Make_3(lhs, StaticGet_TypeCheck_PROPERTY_ASSIGNMENT_TO_READONLY_VALUE(), []string{type_.String()}))
				break
			}
		}
	}
}

// Java static method: TypeCheck.maybeReferencedType(com.google.javascript.rhino.jstype.JSType type)
func TypeCheck_MaybeReferencedType(type_ pkg_rhino.JSType) pkg_rhino.JSType {
	clinit_TypeCheck()
	_ = type_
	var maybeTemplatized *pkg_rhino.TemplatizedType = type_.ToMaybeTemplatizedType()
	_ = maybeTemplatized
	if maybeTemplatized != nil {
		return j2gort.AsIface[pkg_rhino.JSType](maybeTemplatized.GetReferencedType())
	} else {
		return nil
	}
}

// Java static method: TypeCheck.declaresOverride(com.google.javascript.rhino.@org.jspecify.annotations.Nullable JSDocInfo jsdoc)
func TypeCheck_DeclaresOverride(jsdoc *pkg_rhino.JSDocInfo) bool {
	clinit_TypeCheck()
	_ = jsdoc
	return ((jsdoc != nil) && jsdoc.IsOverride())
}

// Java class: com.google.javascript.jscomp.TypeCheck.SuggestionPair
type TypeCheck_SuggestionPair struct {
	suggestion string `gson:"suggestion"`
	distance   int32  `gson:"distance"`
	idHash_    atomic.Int32
}

func (this_ *TypeCheck_SuggestionPair) BindSelf_(s any) {
	_ = s
}
func (this_ *TypeCheck_SuggestionPair) IsJavaNil_() bool { return this_ == nil }
func (*TypeCheck_SuggestionPair) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypeCheck$SuggestionPair", "SuggestionPair"
}
func (this_ *TypeCheck_SuggestionPair) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypeCheck_SuggestionPair) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypeCheck.SuggestionPair", "com.google.javascript.jscomp.TypeCheck$SuggestionPair", func(v_ any) bool {
		return j2gort.InstanceOf[*TypeCheck_SuggestionPair](v_)
	})
}

func (this_ *TypeCheck_SuggestionPair) initFields_TypeCheck_SuggestionPair() {
	_ = this_
	this_.suggestion = j2gort.NullString
}

// Java constructor: SuggestionPair(java.lang.String suggestion, int distance)
func NewTypeCheck_SuggestionPair(suggestion string, distance int32) *TypeCheck_SuggestionPair {
	this_ := &TypeCheck_SuggestionPair{}
	InitTypeCheck_SuggestionPair(this_, suggestion, distance)
	return this_
}

func InitTypeCheck_SuggestionPair(this_ *TypeCheck_SuggestionPair, suggestion string, distance int32) {
	_ = suggestion
	_ = distance
	this_.initFields_TypeCheck_SuggestionPair()
	this_.suggestion = suggestion
	this_.distance = distance
}

// Java class: com.google.javascript.jscomp.TypeCheck.DebugTypeLogger
type TypeCheck_DebugTypeLogger struct {
	outer_      *TypeCheck
	typeLogFile pkg_diagnostic.LogFile `gson:"typeLogFile"`
	parentNode  pkg_rhino.Node         `gson:"parentNode"`
	idHash_     atomic.Int32
}

func (this_ *TypeCheck_DebugTypeLogger) BindSelf_(s any) {
	_ = s
}
func (this_ *TypeCheck_DebugTypeLogger) IsJavaNil_() bool { return this_ == nil }
func (*TypeCheck_DebugTypeLogger) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TypeCheck$DebugTypeLogger", "DebugTypeLogger"
}
func (this_ *TypeCheck_DebugTypeLogger) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TypeCheck_DebugTypeLogger) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TypeCheck.DebugTypeLogger", "com.google.javascript.jscomp.TypeCheck$DebugTypeLogger", func(v_ any) bool {
		return j2gort.InstanceOf[*TypeCheck_DebugTypeLogger](v_)
	})
}
func (this_ *TypeCheck_DebugTypeLogger) SetOuter_(o *TypeCheck)                   { this_.outer_ = o }
func (this_ *TypeCheck_DebugTypeLogger) GetOuter_() *TypeCheck                    { return this_.outer_ }
func (*TypeCheck_DebugTypeLogger) Implements_j2gort_Ext_java_lang_AutoCloseable() {}

func (this_ *TypeCheck_DebugTypeLogger) initFields_TypeCheck_DebugTypeLogger() {
	_ = this_
}

// Java constructor: DebugTypeLogger()
func NewTypeCheck_DebugTypeLogger(outer_ *TypeCheck) *TypeCheck_DebugTypeLogger {
	this_ := &TypeCheck_DebugTypeLogger{}
	this_.outer_ = outer_
	InitTypeCheck_DebugTypeLogger(this_)
	return this_
}

func InitTypeCheck_DebugTypeLogger(this_ *TypeCheck_DebugTypeLogger) {
	this_.initFields_TypeCheck_DebugTypeLogger()
	if this_.outer_.compiler.IsDebugLoggingEnabled() {
		this_.typeLogFile = this_.outer_.compiler.CreateOrReopenLog(j2gort.GetObjectClass(j2gort.PtrToAny(this_)), "types_logged_in_compiler.log", []string{})
	} else {
		this_.typeLogFile = nil
	}
}

// Java method: DebugTypeLogger.maybeStartLoggingAt(com.google.javascript.rhino.Node n)
func (this_ *TypeCheck_DebugTypeLogger) MaybeStartLoggingAt(n pkg_rhino.Node) {
	_ = n
	if !(this_.outer_.compiler.IsDebugLoggingEnabled()) {
		return
	}
	if this_.parentNode != 0 {
		return
	}
	var info *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = info
	if (info == nil) || !(info.GetLogTypeInCompiler()) {
		return
	}
	this_.parentNode = n
}

// Java method: DebugTypeLogger.stopLoggingIfThisIsWhereWeStarted(com.google.javascript.rhino.Node n)
func (this_ *TypeCheck_DebugTypeLogger) StopLoggingIfThisIsWhereWeStarted(n pkg_rhino.Node) {
	_ = n
	if n == this_.parentNode {
		this_.parentNode = pkg_rhino.Node(0)
	}
}

// Java method: DebugTypeLogger.close()
func (this_ *TypeCheck_DebugTypeLogger) Close_0() {
	if this_.typeLogFile != nil && !this_.typeLogFile.IsJavaNil_() {
		this_.typeLogFile.Close_0()
	}
}

// Java method: DebugTypeLogger.maybeLogTypeOfNode(com.google.javascript.rhino.Node n)
func (this_ *TypeCheck_DebugTypeLogger) MaybeLogTypeOfNode(n pkg_rhino.Node) {
	_ = n
	if this_.parentNode == 0 {
		return
	}
	if j2gort.IfaceEquals[pkg_rhino.JSType](n.GetJSType(), nil) {
		return
	}
	if this_.typeLogFile.IsLogging() {
		this_.typeLogFile.Log_1_string(n.ToString_3(true, true, true))
	}
}
