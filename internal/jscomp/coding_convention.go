// Port of com/google/javascript/jscomp/CodingConvention.java @ closure-compiler v20261015
// (with com/google/javascript/jscomp/AutoValue_CodingConvention_AssertionFunctionSpec.java
// folded in)

/*
 * Copyright 2007 The Closure Compiler Authors.
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
	"strconv"
	"strings"
	"sync"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CodingConvention defines a set of hooks to customize the behavior of the
// Compiler for a specific team/company.
//
// Overloads (PORTING.md §3): isExported(String) keeps the plain name and
// isExported(String, boolean) is IsExportedWithLocal. Java String results
// that may be null are *string.
type CodingConvention interface {
	// IsConstant checks whether a given variable name, such as a name in
	// all-caps should be treated as if it had the @const annotation.
	//
	// variableName is a potentially constant variable name; returns true if
	// the name should be treated as a constant.
	//
	// Java: CodingConvention.isConstant(String)
	IsConstant(variableName string) bool

	// IsConstantKey checks whether a given key of an object literal, such as
	// a name in all-caps should be treated as if it had the @const
	// annotation.
	//
	// Java: CodingConvention.isConstantKey(String)
	IsConstantKey(keyName string) bool

	// IsValidEnumKey checks that a given key may be used as a key for an
	// enum.
	//
	// key is the potential key to an enum; returns true if the key may be
	// used as an enum key, false otherwise.
	//
	// Java: CodingConvention.isValidEnumKey(String)
	IsValidEnumKey(key string) bool

	// IsOptionalParameter checks whether a given parameter name should be
	// treated as an optional parameter as far as type checking or function
	// call arg count checking is concerned. Note that an optional function
	// parameter may be declared as a simple type and is automatically
	// converted to a union of the declared type and Undefined.
	//
	// parameter is the parameter's node; returns true if the parameter
	// should be treated as an optional parameter.
	//
	// Java: CodingConvention.isOptionalParameter(Node)
	IsOptionalParameter(parameter rhino.Node) bool

	// IsVarArgsParameter checks whether a given parameter should be treated
	// as a marker for a variable argument list function. A VarArgs parameter
	// must be the last parameter in a function declaration.
	//
	// parameter is the parameter's node; returns true if the parameter
	// should be treated as a variable length parameter.
	//
	// Java: CodingConvention.isVarArgsParameter(Node)
	IsVarArgsParameter(parameter rhino.Node) bool

	// IsFunctionCallThatAlwaysThrows is used by CheckMissingReturn. When a
	// function call always throws an error, it can be the last stm of a
	// block and we don't warn about missing return.
	//
	// Java: CodingConvention.isFunctionCallThatAlwaysThrows(Node)
	IsFunctionCallThatAlwaysThrows(n rhino.Node) bool

	// IsExportedWithLocal checks whether a global variable or function name
	// should be treated as exported, or externally referenceable.
	//
	// name is a global variable or function name; local is true if the name
	// is a local variable. Returns true if the name should be considered
	// exported.
	//
	// Java: CodingConvention.isExported(String, boolean)
	IsExportedWithLocal(name string, local bool) bool

	// IsExported is equivalent to `isExported(name, true) ||
	// isExported(name, false);`
	//
	// Should only be used to check if a property is exported. Variables
	// should always use IsExportedWithLocal, as in most cases local
	// variables should not be treated as exported.
	//
	// Do not override! Unfortunately, that cannot be enforced without making
	// this an abstract class.
	//
	// A Java default method: implementations that do not override it call
	// CodingConvention_IsExported.
	//
	// Java: CodingConvention.isExported(String)
	IsExported(name string) bool

	// BlockRenamingForProperty is a no-op convention that is never used.
	//
	// Deprecated: since this has a default value and is never checked, you
	// can safely remove any overrides of this method from your coding
	// conventions.
	//
	// A Java default method: implementations that do not override it call
	// CodingConvention_BlockRenamingForProperty.
	//
	// Java: CodingConvention.blockRenamingForProperty(String)
	BlockRenamingForProperty(name string) bool

	// GetPackageName returns the package name for the given source file, or
	// null if no package name is known.
	//
	// Java: CodingConvention.getPackageName(StaticSourceFile)
	GetPackageName(source rhino.StaticSourceFile) *string

	// GetClassesDefinedByCall checks if the given method defines a subclass
	// relationship, and if it does, returns information on that
	// relationship. By default, always returns null. Meant to be overridden
	// by subclasses.
	//
	// callNode is a CALL node.
	//
	// Java: CodingConvention.getClassesDefinedByCall(Node)
	GetClassesDefinedByCall(callNode rhino.Node) *CodingConvention_SubclassRelationship

	// IsClassFactoryCall checks if the given method is a call to a class
	// factory, such a factory returns a unique class.
	//
	// callNode is a CALL node.
	//
	// Java: CodingConvention.isClassFactoryCall(Node)
	IsClassFactoryCall(callNode rhino.Node) bool

	// IsSuperClassReference returns true if passed a string referring to the
	// superclass. The string will usually be from the string node at the
	// right of a GETPROP, e.g. this.superClass_.
	//
	// Java: CodingConvention.isSuperClassReference(String)
	IsSuperClassReference(propertyName string) bool

	// ExtractIsModuleFile is a convenience method for determining if the
	// node indicates the file is a "module" file (a file whose top level
	// symbols are not in global scope).
	//
	// Java: CodingConvention.extractIsModuleFile(Node, Node)
	ExtractIsModuleFile(node rhino.Node, parent rhino.Node) bool

	// ExtractClassNameIfProvide is a convenience method for determining
	// provided dependencies amongst different JS scripts.
	//
	// Java: CodingConvention.extractClassNameIfProvide(Node, Node)
	ExtractClassNameIfProvide(node rhino.Node, parent rhino.Node) *string

	// ExtractClassNameIfRequire is a convenience method for determining
	// required dependencies amongst different JS scripts.
	//
	// Java: CodingConvention.extractClassNameIfRequire(Node, Node)
	ExtractClassNameIfRequire(node rhino.Node, parent rhino.Node) *string

	// GetExportPropertyFunction is the function name used when exporting
	// properties. Signature: fn(object, publicName, symbol).
	//
	// Java: CodingConvention.getExportPropertyFunction()
	GetExportPropertyFunction() *string

	// GetExportSymbolFunction is the function name used when exporting
	// symbols. Signature: fn(publicPath, object).
	//
	// Java: CodingConvention.getExportSymbolFunction()
	GetExportSymbolFunction() *string

	// IdentifyTypeDeclarationCall checks if the given CALL node is
	// forward-declaring any types, and returns the name of the types if it
	// is.
	//
	// Java: CodingConvention.identifyTypeDeclarationCall(Node)
	IdentifyTypeDeclarationCall(n rhino.Node) []string

	// ApplySubclassRelationship: in many JS libraries, the function that
	// produces inheritance also adds properties to the superclass and/or
	// subclass.
	//
	// Java: CodingConvention.applySubclassRelationship(NominalTypeBuilder, NominalTypeBuilder, SubclassType)
	ApplySubclassRelationship(
		parent *rhino.NominalTypeBuilder, child *rhino.NominalTypeBuilder, type_ CodingConvention_SubclassType)

	// GetAbstractMethodName is the function name for abstract methods. An
	// abstract method can be assigned to an interface method instead of an
	// function expression in order to avoid linter warnings produced by
	// assigning a function without a return value where a return value is
	// expected.
	//
	// Java: CodingConvention.getAbstractMethodName()
	GetAbstractMethodName() *string

	// GetSingletonGetterClassName checks if the given method defines a
	// singleton getter, and if it does, returns the name of the class with
	// the singleton getter. By default, always returns null. Meant to be
	// overridden by subclasses.
	//
	// addSingletonGetter needs a coding convention because in the general
	// case, it can't be inlined. The function inliner sees that it creates
	// an alias to the given class in an inner closure, and bails out.
	//
	// callNode is a CALL node.
	//
	// Java: CodingConvention.getSingletonGetterClassName(Node)
	GetSingletonGetterClassName(callNode rhino.Node) *string

	// ApplySingletonGetter: in many JS libraries, the function that adds a
	// singleton getter to a class adds properties to the class.
	//
	// Java: CodingConvention.applySingletonGetter(NominalTypeBuilder, FunctionType)
	ApplySingletonGetter(classType *rhino.NominalTypeBuilder, getterType rhino.FunctionType)

	// DescribeFunctionBind returns a Bind instance or null.
	//
	// When seeing an expression exp1.bind(recv, arg1, ...); we only know
	// that it's a function bind if exp1 has type function. W/out type info,
	// exp1 has certainly a function type only if it's a function literal.
	//
	// If type checking has already happened, exp1's type is attached to the
	// AST node. When checkTypes is true, describeFunctionBind looks for that
	// type.
	//
	// checkTypes: check that the bound node has a function type.
	//
	// Java: CodingConvention.describeFunctionBind(Node, boolean)
	DescribeFunctionBind(n rhino.Node, checkTypes bool) *CodingConvention_Bind

	// DescribeCachingCall builds a Cache instance from the given call node
	// and returns that instance, or null if the Node does not resemble a
	// cache utility call.
	//
	// This should match calls to a cache utility method. This type of node
	// is specially considered for side-effects since conventionally storing
	// something on a cache object would be seen as a side-effect.
	//
	// Java: CodingConvention.describeCachingCall(Node)
	DescribeCachingCall(node rhino.Node) *CodingConvention_Cache

	// IsPropertyTestFunction: whether this CALL function is testing for the
	// existence of a property.
	//
	// Java: CodingConvention.isPropertyTestFunction(Node)
	IsPropertyTestFunction(call rhino.Node) bool

	// IsPrototypeAlias: whether this GETPROP node is an alias for an object
	// prototype.
	//
	// Java: CodingConvention.isPrototypeAlias(Node)
	IsPrototypeAlias(getProp rhino.Node) bool

	// IsPropertyRenameFunction: whether this GETPROP or NAME node is the
	// function is returning the string name for a property, but allows
	// renaming.
	//
	// Java: CodingConvention.isPropertyRenameFunction(Node)
	IsPropertyRenameFunction(nameNode rhino.Node) bool

	// GetObjectLiteralCast checks if the given method performs a object
	// literal cast, and if it does, returns information on the cast. By
	// default, always returns null. Meant to be overridden by subclasses.
	//
	// callNode is a CALL node.
	//
	// Java: CodingConvention.getObjectLiteralCast(Node)
	GetObjectLiteralCast(callNode rhino.Node) *CodingConvention_ObjectLiteralCast

	// GetAssertionFunctions returns the set of AssertionFunction.
	//
	// Java: CodingConvention.getAssertionFunctions()
	GetAssertionFunctions() []*CodingConvention_AssertionFunctionSpec
}

// CodingConvention_IsExported is the default body of
// CodingConvention.isExported(String).
//
// Java: CodingConvention.isExported(String) (default method)
func CodingConvention_IsExported(c CodingConvention, name string) bool {
	return c.IsExportedWithLocal(name, true) || c.IsExportedWithLocal(name, false)
}

// CodingConvention_BlockRenamingForProperty is the default body of
// CodingConvention.blockRenamingForProperty(String).
//
// Java: CodingConvention.blockRenamingForProperty(String) (default method)
func CodingConvention_BlockRenamingForProperty(c CodingConvention, name string) bool {
	return false
}

// CodingConvention_Bind is the Bind class.
type CodingConvention_Bind struct {
	// The target of the bind action
	target rhino.Node
	// The node representing the "this" value, maybe null
	thisValue rhino.Node
	// The head of a Node list representing the parameters
	parameters rhino.Node
}

// Java: CodingConvention.Bind(Node, Node, Node)
func NewCodingConvention_Bind(target rhino.Node, thisValue rhino.Node, parameters rhino.Node) *CodingConvention_Bind {
	b := &CodingConvention_Bind{}
	b.target = target
	b.thisValue = thisValue
	b.parameters = parameters
	return b
}

// GetBoundParameterCount is the number of parameters bound (not including
// the 'this' value).
//
// Java: CodingConvention.Bind.getBoundParameterCount()
func (b *CodingConvention_Bind) GetBoundParameterCount() int {
	if b.parameters.IsNil() {
		return 0
	}
	paramParent := b.parameters.GetParent()
	return paramParent.GetChildCount() - paramParent.GetIndexOfChild(b.parameters)
}

// CodingConvention_Cache is the Cache class.
type CodingConvention_Cache struct {
	cacheObj rhino.Node
	key      rhino.Node
	valueFn  rhino.Node
	keyFn    rhino.Node
}

// Java: CodingConvention.Cache(Node, Node, Node, Node)
func NewCodingConvention_Cache(cacheObj rhino.Node, key rhino.Node, valueFn rhino.Node, keyFn rhino.Node) *CodingConvention_Cache {
	c := &CodingConvention_Cache{}
	c.cacheObj = cacheObj
	c.key = key
	c.valueFn = valueFn
	c.keyFn = keyFn
	return c
}

// CodingConvention_SubclassType specifies the kind of inheritance.
type CodingConvention_SubclassType int32

const (
	CodingConvention_SubclassType_INHERITS CodingConvention_SubclassType = iota
	CodingConvention_SubclassType_MIXIN
)

var codingConvention_SubclassType_names = [...]string{"INHERITS", "MIXIN"}

// CodingConvention_SubclassType_values is SubclassType.values().
func CodingConvention_SubclassType_values() []CodingConvention_SubclassType {
	return []CodingConvention_SubclassType{
		CodingConvention_SubclassType_INHERITS, CodingConvention_SubclassType_MIXIN,
	}
}

// Name is Enum.name().
func (t CodingConvention_SubclassType) Name() string { return codingConvention_SubclassType_names[t] }

// String is Enum.toString().
func (t CodingConvention_SubclassType) String() string { return t.Name() }

// Ordinal is Enum.ordinal().
func (t CodingConvention_SubclassType) Ordinal() int { return int(t) }

// CodingConvention_SubclassRelationship records subclass relations.
type CodingConvention_SubclassRelationship struct {
	type_          CodingConvention_SubclassType
	subclassName   string
	superclassName string
}

type pooledSubclassRelSlab struct {
	chunk []CodingConvention_SubclassRelationship
}

var subclassRelSlabPool = sync.Pool{
	New: func() any { return &pooledSubclassRelSlab{} },
}

// Java: CodingConvention.SubclassRelationship(SubclassType, Node, Node)
func NewCodingConvention_SubclassRelationship(
	type_ CodingConvention_SubclassType, subclassNode rhino.Node, superclassNode rhino.Node) *CodingConvention_SubclassRelationship {
	subName, subOk := subclassNode.GetQualifiedNameString()
	if !subOk {
		jutil.CheckArgumentMsg1(false, "Expected qualified name, found: %s", subclassNode)
	}
	supName, supOk := superclassNode.GetQualifiedNameString()
	if !supOk {
		jutil.CheckArgumentMsg1(false, "Expected qualified name, found: %s", superclassNode)
	}
	s := subclassRelSlabPool.Get().(*pooledSubclassRelSlab)
	if len(s.chunk) == 0 {
		s.chunk = make([]CodingConvention_SubclassRelationship, 512) // 20 KiB <= 32 KiB
	}
	r := &s.chunk[0]
	s.chunk = s.chunk[1:]
	subclassRelSlabPool.Put(s)
	r.type_ = type_
	r.subclassName = subName
	r.superclassName = supName
	return r
}

// CodingConvention_ObjectLiteralCast: an object literal cast provides a
// mechanism to cast object literals to other types without a warning.
type CodingConvention_ObjectLiteralCast struct {
	// Type to cast to.
	typeName *string

	// Object to cast.
	objectNode rhino.Node

	// Error message
	diagnosticType *DiagnosticType
}

// Java: CodingConvention.ObjectLiteralCast(String, Node, DiagnosticType)
func NewCodingConvention_ObjectLiteralCast(
	typeName *string, objectNode rhino.Node, diagnosticType *DiagnosticType) *CodingConvention_ObjectLiteralCast {
	c := &CodingConvention_ObjectLiteralCast{}
	c.typeName = typeName
	c.objectNode = objectNode
	c.diagnosticType = diagnosticType
	return c
}

// CodingConvention_AssertionFunctionSpec is a description of a JavaScript
// function that will throw an exception when either:
//
//   - One of its parameters does not match the return type of the function
//   - One of its parameters is falsy. This has some special handling for
//     expressions that the match-return-type handling does not have.
//
// Java: @AutoValue abstract class CodingConvention.AssertionFunctionSpec,
// with its generated AutoValue_CodingConvention_AssertionFunctionSpec
// folded in: the struct holds the AutoValue fields.
type CodingConvention_AssertionFunctionSpec struct {
	functionName     *string
	closurePrimitive *rhino.ClosurePrimitive
	assertionKind    CodingConvention_AssertionFunctionSpec_AssertionKind
	paramIndex       int
}

// Java: AutoValue_CodingConvention_AssertionFunctionSpec(String, ClosurePrimitive, AssertionKind, int)
func newCodingConvention_AssertionFunctionSpec(
	functionName *string,
	closurePrimitive *rhino.ClosurePrimitive,
	assertionKind CodingConvention_AssertionFunctionSpec_AssertionKind,
	paramIndex int) *CodingConvention_AssertionFunctionSpec {
	s := &CodingConvention_AssertionFunctionSpec{}
	s.functionName = functionName
	s.closurePrimitive = closurePrimitive
	s.assertionKind = assertionKind
	s.paramIndex = paramIndex
	return s
}

// TODO(b/126254920): remove this field and always use ClosurePrimitive
//
// Java: CodingConvention.AssertionFunctionSpec.getFunctionName()
func (s *CodingConvention_AssertionFunctionSpec) GetFunctionName() *string {
	return s.functionName
}

// Java: CodingConvention.AssertionFunctionSpec.getClosurePrimitive()
func (s *CodingConvention_AssertionFunctionSpec) GetClosurePrimitive() *rhino.ClosurePrimitive {
	return s.closurePrimitive
}

// Java: CodingConvention.AssertionFunctionSpec.getAssertionKind()
func (s *CodingConvention_AssertionFunctionSpec) GetAssertionKind() CodingConvention_AssertionFunctionSpec_AssertionKind {
	return s.assertionKind
}

// GetParamIndex is the index of the formal parameter that is actually
// asserted.
//
// Java: CodingConvention.AssertionFunctionSpec.getParamIndex()
func (s *CodingConvention_AssertionFunctionSpec) GetParamIndex() int {
	return s.paramIndex
}

// Java: AutoValue_CodingConvention_AssertionFunctionSpec.toString()
func (s *CodingConvention_AssertionFunctionSpec) String() string {
	functionName := "null"
	if s.functionName != nil {
		functionName = *s.functionName
	}
	closurePrimitive := "null"
	if s.closurePrimitive != nil {
		closurePrimitive = s.closurePrimitive.String()
	}
	return "AssertionFunctionSpec{" +
		"functionName=" + functionName + ", " +
		"closurePrimitive=" + closurePrimitive + ", " +
		"assertionKind=" + s.assertionKind.String() + ", " +
		"paramIndex=" + strconv.Itoa(s.paramIndex) +
		"}"
}

// Java: AutoValue_CodingConvention_AssertionFunctionSpec.equals(Object)
func (s *CodingConvention_AssertionFunctionSpec) Equals(o any) bool {
	that, ok := o.(*CodingConvention_AssertionFunctionSpec)
	if !ok || that == nil {
		return false
	}
	if that == s {
		return true
	}
	var functionNameEquals bool
	if s.functionName == nil {
		functionNameEquals = that.GetFunctionName() == nil
	} else {
		functionNameEquals = that.GetFunctionName() != nil && *s.functionName == *that.GetFunctionName()
	}
	var closurePrimitiveEquals bool
	if s.closurePrimitive == nil {
		closurePrimitiveEquals = that.GetClosurePrimitive() == nil
	} else {
		closurePrimitiveEquals = that.GetClosurePrimitive() != nil && *s.closurePrimitive == *that.GetClosurePrimitive()
	}
	return functionNameEquals &&
		closurePrimitiveEquals &&
		s.assertionKind == that.GetAssertionKind() &&
		s.paramIndex == that.GetParamIndex()
}

// HashCode is AutoValue's hashCode.
//
// Go-only: language-adapter — Java enums hash by identity (a value that
// differs from run to run); the enum fields hash by ordinal here.
//
// Java: AutoValue_CodingConvention_AssertionFunctionSpec.hashCode()
func (s *CodingConvention_AssertionFunctionSpec) HashCode() int32 {
	var h int32 = 1
	h *= 1000003
	if s.functionName != nil {
		h ^= jutil.StringHashCode(*s.functionName)
	}
	h *= 1000003
	if s.closurePrimitive != nil {
		h ^= int32(*s.closurePrimitive)
	}
	h *= 1000003
	h ^= int32(s.assertionKind)
	h *= 1000003
	h ^= int32(s.paramIndex)
	return h
}

// CodingConvention_AssertionFunctionSpec_AssertionKind is
// AssertionFunctionSpec.AssertionKind.
type CodingConvention_AssertionFunctionSpec_AssertionKind int32

const (
	// an assertion that the parameter is 'truthy' and returns the parameter
	CodingConvention_AssertionFunctionSpec_AssertionKind_TRUTHY CodingConvention_AssertionFunctionSpec_AssertionKind = iota
	// an assertion that the parameter matches the inferred return kind
	CodingConvention_AssertionFunctionSpec_AssertionKind_MATCHES_RETURN_TYPE
)

var codingConvention_AssertionFunctionSpec_AssertionKind_names = [...]string{"TRUTHY", "MATCHES_RETURN_TYPE"}

// CodingConvention_AssertionFunctionSpec_AssertionKind_values is
// AssertionKind.values().
func CodingConvention_AssertionFunctionSpec_AssertionKind_values() []CodingConvention_AssertionFunctionSpec_AssertionKind {
	return []CodingConvention_AssertionFunctionSpec_AssertionKind{
		CodingConvention_AssertionFunctionSpec_AssertionKind_TRUTHY,
		CodingConvention_AssertionFunctionSpec_AssertionKind_MATCHES_RETURN_TYPE,
	}
}

// Name is Enum.name().
func (k CodingConvention_AssertionFunctionSpec_AssertionKind) Name() string {
	return codingConvention_AssertionFunctionSpec_AssertionKind_names[k]
}

// String is Enum.toString().
func (k CodingConvention_AssertionFunctionSpec_AssertionKind) String() string { return k.Name() }

// Ordinal is Enum.ordinal().
func (k CodingConvention_AssertionFunctionSpec_AssertionKind) Ordinal() int { return int(k) }

// CodingConvention_AssertionFunctionSpec_Builder_ is
// AssertionFunctionSpec.builder().
//
// Java: CodingConvention.AssertionFunctionSpec.builder()
func CodingConvention_AssertionFunctionSpec_Builder_() *CodingConvention_AssertionFunctionSpec_Builder {
	return newCodingConvention_AssertionFunctionSpec_Builder().SetParamIndex(0)
}

// Java: CodingConvention.AssertionFunctionSpec.forTruthy()
func CodingConvention_AssertionFunctionSpec_ForTruthy() *CodingConvention_AssertionFunctionSpec_Builder {
	return CodingConvention_AssertionFunctionSpec_Builder_().
		SetAssertionKind(CodingConvention_AssertionFunctionSpec_AssertionKind_TRUTHY)
}

// Java: CodingConvention.AssertionFunctionSpec.forMatchesReturn()
func CodingConvention_AssertionFunctionSpec_ForMatchesReturn() *CodingConvention_AssertionFunctionSpec_Builder {
	return CodingConvention_AssertionFunctionSpec_Builder_().
		SetAssertionKind(CodingConvention_AssertionFunctionSpec_AssertionKind_MATCHES_RETURN_TYPE)
}

// CodingConvention_AssertionFunctionSpec_Builder is
// AssertionFunctionSpec.Builder.
//
// Java: @AutoValue.Builder abstract class AssertionFunctionSpec.Builder, with
// the generated AutoValue_CodingConvention_AssertionFunctionSpec.Builder
// folded in: the struct holds the generated builder's fields.
type CodingConvention_AssertionFunctionSpec_Builder struct {
	functionName     *string
	closurePrimitive *rhino.ClosurePrimitive
	assertionKind    *CodingConvention_AssertionFunctionSpec_AssertionKind
	paramIndex       int
	set0             int8
}

// Java: AutoValue_CodingConvention_AssertionFunctionSpec.Builder()
func newCodingConvention_AssertionFunctionSpec_Builder() *CodingConvention_AssertionFunctionSpec_Builder {
	return &CodingConvention_AssertionFunctionSpec_Builder{}
}

// Java: CodingConvention.AssertionFunctionSpec.Builder.setFunctionName(String)
func (b *CodingConvention_AssertionFunctionSpec_Builder) SetFunctionName(name string) *CodingConvention_AssertionFunctionSpec_Builder {
	b.functionName = &name
	return b
}

// Java: CodingConvention.AssertionFunctionSpec.Builder.setClosurePrimitive(ClosurePrimitive)
func (b *CodingConvention_AssertionFunctionSpec_Builder) SetClosurePrimitive(primitive rhino.ClosurePrimitive) *CodingConvention_AssertionFunctionSpec_Builder {
	b.closurePrimitive = &primitive
	return b
}

// Java: CodingConvention.AssertionFunctionSpec.Builder.setParamIndex(int)
func (b *CodingConvention_AssertionFunctionSpec_Builder) SetParamIndex(paramIndex int) *CodingConvention_AssertionFunctionSpec_Builder {
	b.paramIndex = paramIndex
	b.set0 |= 1
	return b
}

// SetAssertionKind is the generated setter, which rejects null (a Go
// AssertionKind value is never null).
//
// Java: CodingConvention.AssertionFunctionSpec.Builder.setAssertionKind(AssertionKind)
func (b *CodingConvention_AssertionFunctionSpec_Builder) SetAssertionKind(
	kind CodingConvention_AssertionFunctionSpec_AssertionKind) *CodingConvention_AssertionFunctionSpec_Builder {
	b.assertionKind = &kind
	return b
}

// Java: AutoValue_CodingConvention_AssertionFunctionSpec.Builder.autoBuild()
func (b *CodingConvention_AssertionFunctionSpec_Builder) AutoBuild() *CodingConvention_AssertionFunctionSpec {
	if b.set0 != 1 || b.assertionKind == nil {
		missing := ""
		if b.assertionKind == nil {
			missing += " assertionKind"
		}
		if (b.set0 & 1) == 0 {
			missing += " paramIndex"
		}
		panic(jutil.NewIllegalStateExceptionMsg("Missing required properties:" + missing))
	}
	return newCodingConvention_AssertionFunctionSpec(
		b.functionName,
		b.closurePrimitive,
		*b.assertionKind,
		b.paramIndex)
}

// Java: CodingConvention.AssertionFunctionSpec.Builder.build()
func (b *CodingConvention_AssertionFunctionSpec_Builder) Build() *CodingConvention_AssertionFunctionSpec {
	spec := b.AutoBuild()
	jutil.CheckStateMsg(
		spec.GetFunctionName() != nil || spec.GetClosurePrimitive() != nil,
		"Must provide a function name or ClosurePrimitive for each spec")
	return spec
}

// GetId returns the spec's key: its ClosurePrimitive, or else its function
// name.
//
// Go-only: language-adapter — the Object result is a rhino.ClosurePrimitive
// or a string value (nil for Java null).
//
// Java: CodingConvention.AssertionFunctionSpec.getId()
func (s *CodingConvention_AssertionFunctionSpec) GetId() any {
	if s.GetClosurePrimitive() != nil {
		return *s.GetClosurePrimitive()
	}
	if s.GetFunctionName() != nil {
		return *s.GetFunctionName()
	}
	return nil
}

// GetAssertedArg returns which argument is actually being asserted, or null
// if fewer args than expected
//
// Java: CodingConvention.AssertionFunctionSpec.getAssertedArg(Node)
func (s *CodingConvention_AssertionFunctionSpec) GetAssertedArg(firstArg rhino.Node) rhino.Node {
	for i := 0; i < s.GetParamIndex(); i++ {
		if firstArg.IsNil() {
			// If there are fewer arguments than expected, return null instead of crashing in this
			// function.
			return rhino.Node{}
		}
		firstArg = firstArg.GetNext()
	}
	return firstArg
}

// CodingConvention_AssertionFunctionLookup stores a relation from either
// name or Closure Primitive to assertion function
type CodingConvention_AssertionFunctionLookup struct {
	// the key type 'Object' is mutable, but at runtime is only ever a string
	// or ClosurePrimitive
	//
	// Go-only: language-adapter — keys are string or rhino.ClosurePrimitive
	// values, whose Go == is Java equals.
	internal    *jutil.LinkedHashMap[any, *CodingConvention_AssertionFunctionSpec]
	byPrimitive map[rhino.ClosurePrimitive]*CodingConvention_AssertionFunctionSpec
	byName      map[string]*CodingConvention_AssertionFunctionSpec
	byLastProp  map[string][]*CodingConvention_AssertionFunctionSpec
}

// Java: CodingConvention.AssertionFunctionLookup(ImmutableMap<Object, AssertionFunctionSpec>)
func newCodingConvention_AssertionFunctionLookup(
	internal *jutil.LinkedHashMap[any, *CodingConvention_AssertionFunctionSpec]) *CodingConvention_AssertionFunctionLookup {
	l := &CodingConvention_AssertionFunctionLookup{}
	l.internal = internal
	if internal != nil && !internal.IsEmpty() {
		for k, spec := range internal.All() {
			switch key := k.(type) {
			case rhino.ClosurePrimitive:
				if l.byPrimitive == nil {
					l.byPrimitive = make(map[rhino.ClosurePrimitive]*CodingConvention_AssertionFunctionSpec)
				}
				l.byPrimitive[key] = spec
			case string:
				if lastDot := strings.LastIndexByte(key, '.'); lastDot >= 0 {
					if l.byLastProp == nil {
						l.byLastProp = make(map[string][]*CodingConvention_AssertionFunctionSpec)
					}
					lastProp := key[lastDot+1:]
					l.byLastProp[lastProp] = append(l.byLastProp[lastProp], spec)
				} else {
					if l.byName == nil {
						l.byName = make(map[string]*CodingConvention_AssertionFunctionSpec)
					}
					l.byName[key] = spec
				}
			}
		}
	}
	return l
}

// CodingConvention_AssertionFunctionLookup_Of returns a new map containing
// all the given AssertionFunctionSpecs
//
// Assumes that in the input, there is a unique mapping from string name to
// spec and closure primitive to spec.
//
// Java: CodingConvention.AssertionFunctionLookup.of(Collection<AssertionFunctionSpec>)
func CodingConvention_AssertionFunctionLookup_Of(
	specs []*CodingConvention_AssertionFunctionSpec) *CodingConvention_AssertionFunctionLookup {
	builder := jutil.NewImmutableMapBuilder[any, *CodingConvention_AssertionFunctionSpec]()
	for _, spec := range specs {
		builder.Put(spec.GetId(), spec)
	}
	idToSpecMap := builder.BuildOrThrow()

	return newCodingConvention_AssertionFunctionLookup(idToSpecMap)
}

// LookupByCallee returns the AssertionFunctionSpec matching the given
// function reference.
//
// This first looks up specs by their ClosurePrimitive, then falls back to
// qualified name
//
// Java: CodingConvention.AssertionFunctionLookup.lookupByCallee(Node)
func (l *CodingConvention_AssertionFunctionLookup) LookupByCallee(callee rhino.Node) *CodingConvention_AssertionFunctionSpec {
	fnType := rhino.JSType_ToMaybeFunctionType(callee.GetJSType())
	if fnType != nil && fnType.GetClosurePrimitive() != nil {
		spec := l.byPrimitive[*fnType.GetClosurePrimitive()]
		if spec != nil {
			return spec
		}
	}

	// TODO(b/126254920): remove this
	if callee.IsName() {
		return l.byName[callee.GetString()]
	}
	if callee.IsGetProp() && len(l.byLastProp) > 0 {
		if specs := l.byLastProp[callee.GetString()]; len(specs) > 0 {
			for _, spec := range specs {
				if spec.functionName != nil && callee.MatchesQualifiedName(*spec.functionName) {
					return spec
				}
			}
		}
	}

	return nil
}
