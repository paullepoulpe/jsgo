// Port of com/google/javascript/rhino/testing/TypeSubject.java @ closure-compiler v20260415

package rhinotesting

import (
	"fmt"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TypeSubject is a Truth Subject for JSType. Usage:
//
//	s.AssertType(type).IsString()
//	s.AssertType(type).IsObjectTypeWithProperty("propName").WithTypeOfProp("propName").IsNumber()
type TypeSubject struct {
	*Subject
	actual rhino.JSType
}

// AssertType is TypeSubject.assertType(JSType).
//
// Java: TypeSubject.assertType(JSType)
func (tr Truth) AssertType(type_ rhino.JSType) *TypeSubject {
	return newTypeSubject(&Subject{t: tr.T, actual: type_}, type_)
}

// TypeSubjectBuilder is Subject.Factory<TypeSubject, JSType>, the result of
// about(types()).
type TypeSubjectBuilder struct {
	t       testing.TB
	message string
}

// AboutTypes is StandardSubjectBuilder.about(TypeSubject.types()).
//
// Java: TypeSubject.types()
func (b *StandardSubjectBuilder) AboutTypes() *TypeSubjectBuilder {
	return &TypeSubjectBuilder{t: b.t, message: b.message}
}

// That is SimpleSubjectBuilder.that(JSType).
func (b *TypeSubjectBuilder) That(type_ rhino.JSType) *TypeSubject {
	return newTypeSubject(&Subject{t: b.t, actual: type_, message: b.message}, type_)
}

func newTypeSubject(s *Subject, type_ rhino.JSType) *TypeSubject {
	if type_ != nil && jutil.IsNil(type_) {
		// Go-only: a nil pointer inside a non-nil JSType interface is a
		// port bug (Go callers comparing with nil would see a type).
		s.t.Helper()
		s.fail("typed nil JSType (%T)", type_)
	}
	return &TypeSubject{Subject: s, actual: type_}
}

func (s *TypeSubject) check(format string, args ...any) *StandardSubjectBuilder {
	msg := "value of: type." + fmt.Sprintf(format, args...) + "\ntype was: " + debugStringOf(s.actual)
	if s.message != "" {
		msg = s.message + "\n" + msg
	}
	return &StandardSubjectBuilder{t: s.t, message: msg}
}

// IsEqualTo checks the actual type equals the provided one, symmetrically
// and with equal hash codes.
//
// Java: TypeSubject.isEqualTo(Object)
func (s *TypeSubject) IsEqualTo(provided rhino.JSType) {
	s.t.Helper()
	s.checkEqualityAgainst(provided, true)
}

// Java: TypeSubject.isNotEqualTo(JSType)
func (s *TypeSubject) IsNotEqualTo(provided rhino.JSType) {
	s.t.Helper()
	s.checkEqualityAgainst(provided, false)
}

// Java: TypeSubject.isNumber()
func (s *TypeSubject) IsNumber() {
	s.t.Helper()
	s.check("isNumberValueType()").That(s.actualNonNull().IsNumberValueType()).IsTrue()
}

// Java: TypeSubject.isOnlyBigInt()
func (s *TypeSubject) IsOnlyBigInt() {
	s.t.Helper()
	s.check("isOnlyBigInt()").That(s.actualNonNull().IsOnlyBigInt()).IsTrue()
}

// Java: TypeSubject.isNotOnlyBigInt()
func (s *TypeSubject) IsNotOnlyBigInt() {
	s.t.Helper()
	s.check("isOnlyBigInt()").That(s.actualNonNull().IsOnlyBigInt()).IsFalse()
}

// Java: TypeSubject.isString()
func (s *TypeSubject) IsString() {
	s.t.Helper()
	s.check("isStringValueType()").That(s.actualNonNull().IsStringValueType()).IsTrue()
}

// Java: TypeSubject.isBoolean()
func (s *TypeSubject) IsBoolean() {
	s.t.Helper()
	s.check("isBooleanValueType()").That(s.actualNonNull().IsBooleanValueType()).IsTrue()
}

// Java: TypeSubject.isVoid()
func (s *TypeSubject) IsVoid() {
	s.t.Helper()
	s.check("isVoidType()").That(s.actualNonNull().IsVoidType()).IsTrue()
}

// Java: TypeSubject.isNoType()
func (s *TypeSubject) IsNoType() {
	s.t.Helper()
	s.check("isNoType()").That(s.actualNonNull().IsNoType()).IsTrue()
}

// Java: TypeSubject.isUnknown()
func (s *TypeSubject) IsUnknown() {
	s.t.Helper()
	s.check("isUnknownType()").That(s.actualNonNull().IsUnknownType()).IsTrue()
}

// Java: TypeSubject.isNotUnknown()
func (s *TypeSubject) IsNotUnknown() {
	s.t.Helper()
	s.check("isUnknownType()").That(s.actualNonNull().IsUnknownType()).IsFalse()
}

// Java: TypeSubject.isNoResolvedType(String)
func (s *TypeSubject) IsNoResolvedType(referenceName string) {
	s.t.Helper()
	s.check("isNoResolvedType()").That(s.actual.IsNoResolvedType()).IsTrue()
	s.GetReferenceNameIsEqualTo(referenceName)
}

// Java: TypeSubject.isNotEmpty()
func (s *TypeSubject) IsNotEmpty() {
	s.t.Helper()
	s.check("isEmptyType()").That(s.actualNonNull().IsEmptyType()).IsFalse()
}

// Java: TypeSubject.isLiteralObject()
func (s *TypeSubject) IsLiteralObject() {
	s.t.Helper()
	s.check("isLiteralObject()").That(s.actualNonNull().IsLiteralObject()).IsTrue()
}

// Java: TypeSubject.isFunctionTypeThat()
func (s *TypeSubject) IsFunctionTypeThat() *TypeSubject_FunctionTypeSubject {
	s.t.Helper()
	s.check("isFunctionType()").That(s.actualNonNull().IsFunctionType()).IsTrue()
	return &TypeSubject_FunctionTypeSubject{outer: s}
}

// IsObjectTypeWithProperty asserts that the actual type is an object type
// with the given property, and returns the subject for chaining.
//
// Java: TypeSubject.isObjectTypeWithProperty(String)
func (s *TypeSubject) IsObjectTypeWithProperty(propName string) *TypeSubject {
	s.t.Helper()
	s.IsLiteralObject()
	s.WithTypeOfProp(propName).IsNotNull()
	return s
}

// Java: TypeSubject.isUnionType()
func (s *TypeSubject) IsUnionType() {
	s.t.Helper()
	s.check("isUnionType()").That(s.actualNonNull().IsUnionType()).IsTrue()
}

// Java: TypeSubject.isUnionOf(JSType...)
func (s *TypeSubject) IsUnionOf(alternates ...rhino.JSType) {
	s.t.Helper()
	s.IsUnionType()
	s.check("getAlternates().equals(%s)", Show(alternates)).
		That(s.actualNonNull().ToMaybeUnionType().GetAlternates()).
		ContainsExactlyElementsIn(alternates)
}

// Java: TypeSubject.isResolved()
func (s *TypeSubject) IsResolved() {
	s.t.Helper()
	s.check("isResolved()").That(s.actualNonNull().IsResolved()).IsTrue()
}

// Java: TypeSubject.isUnresolved()
func (s *TypeSubject) IsUnresolved() {
	s.t.Helper()
	s.check("isResolved()").That(s.actualNonNull().IsResolved()).IsFalse()
}

// WithTypeOfProp returns a TypeSubject for the type of the property with
// the given name. The actual type must be an object type.
//
// Java: TypeSubject.withTypeOfProp(String)
func (s *TypeSubject) WithTypeOfProp(propName string) *TypeSubject {
	s.t.Helper()
	s.check("isObjectType()").That(s.actualNonNull().IsObjectType()).IsTrue()
	return s.check("toMaybeObjectType().getPropertyType(%s)", propName).
		AboutTypes().
		That(s.actualNonNull().ToMaybeObjectType().GetPropertyType(propName))
}

// Java: TypeSubject.withTypeOfProp(KnownSymbolType)
func (s *TypeSubject) WithTypeOfPropSymbol(propName *rhino.KnownSymbolType) *TypeSubject {
	s.t.Helper()
	s.check("isObjectType()").That(s.actualNonNull().IsObjectType()).IsTrue()
	return s.check("toMaybeObjectType().getPropertyType(%s)", Show(propName)).
		AboutTypes().
		That(s.actualNonNull().ToMaybeObjectType().GetPropertyTypeByKey(rhino.NewProperty_SymbolKey(propName)))
}

// Java: TypeSubject.hasDeclaredProperty(String)
func (s *TypeSubject) HasDeclaredProperty(propName string) {
	s.t.Helper()
	s.check("isObjectType()").That(s.actualNonNull().IsObjectType()).IsTrue()
	s.check("toMaybeObjectType().isPropertyTypeDeclared(%s)", propName).
		That(s.actualNonNull().ToMaybeObjectType().IsPropertyTypeDeclared(propName)).
		IsTrue()
}

// Java: TypeSubject.hasProperty(KnownSymbolType)
func (s *TypeSubject) HasProperty(propName *rhino.KnownSymbolType) {
	s.t.Helper()
	s.check("isObjectType()").That(s.actualNonNull().IsObjectType()).IsTrue()
	s.check("toMaybeObjectType().isPropertyTypeDeclared(%s)", Show(propName)).
		That(s.actualNonNull().ToMaybeObjectType().HasPropertyByKey(rhino.NewProperty_SymbolKey(propName))).
		IsTrue()
}

// Java: TypeSubject.hasInferredProperty(String)
func (s *TypeSubject) HasInferredProperty(propName string) {
	s.t.Helper()
	s.check("isObjectType()").That(s.actualNonNull().IsObjectType()).IsTrue()
	s.check("toMaybeObjectType().isPropertyTypeInferred(%s)", propName).
		That(s.actualNonNull().ToMaybeObjectType().IsPropertyTypeInferred(propName)).
		IsTrue()
}

// Java: TypeSubject.isObjectTypeWithoutProperty(String)
func (s *TypeSubject) IsObjectTypeWithoutProperty(propName string) {
	s.t.Helper()
	s.IsLiteralObject()
	s.WithTypeOfProp(propName).IsNull()
}

// Java: TypeSubject.isSubtypeOf(JSType)
func (s *TypeSubject) IsSubtypeOf(superType rhino.JSType) {
	s.t.Helper()
	s.check("isSubtypeOf(%s)", Show(superType)).That(s.actualNonNull().IsSubtypeOf(superType)).IsTrue()
}

// Java: TypeSubject.isNotSubtypeOf(JSType)
func (s *TypeSubject) IsNotSubtypeOf(superType rhino.JSType) {
	s.t.Helper()
	s.check("isSubtypeOf(%s)", Show(superType)).That(s.actualNonNull().IsSubtypeOf(superType)).IsFalse()
}

// Java: TypeSubject.canTestForShallowEqualityWith(JSType)
func (s *TypeSubject) CanTestForShallowEqualityWith(other rhino.JSType) {
	s.t.Helper()
	s.check("canTestForShallowEqualityWith(%s)", Show(other)).
		That(s.actualNonNull().CanTestForShallowEqualityWith(other)).
		IsTrue()
}

// Java: TypeSubject.cannotTestForShallowEqualityWith(JSType)
func (s *TypeSubject) CannotTestForShallowEqualityWith(other rhino.JSType) {
	s.t.Helper()
	s.check("canTestForShallowEqualityWith(%s)", Show(other)).
		That(s.actualNonNull().CanTestForShallowEqualityWith(other)).
		IsFalse()
}

// Java: TypeSubject.toStringIsEqualTo(String)
func (s *TypeSubject) ToStringIsEqualTo(typeString string) {
	s.t.Helper()
	s.check("toString()").That(s.actualNonNull().String()).IsEqualTo(typeString)
}

// Java: TypeSubject.getReferenceNameIsEqualTo(String)
func (s *TypeSubject) GetReferenceNameIsEqualTo(referenceName string) {
	s.t.Helper()
	s.check("getReferenceName()").
		That(s.actualNonNull().ToMaybeObjectType().GetReferenceName()).
		IsEqualTo(referenceName)
	s.check("hasReferenceName()").
		That(s.actualNonNull().ToMaybeObjectType().HasReferenceName()).
		IsTrue()
}

// Java: TypeSubject.getReferenceNameIsNull()
func (s *TypeSubject) GetReferenceNameIsNull() {
	s.t.Helper()
	s.check("getReferenceName()").
		That(s.actualNonNull().ToMaybeObjectType().GetReferenceName()).
		IsNull()
	s.check("hasReferenceName()").
		That(s.actualNonNull().ToMaybeObjectType().HasReferenceName()).
		IsFalse()
}

// Java: TypeSubject.actualNonNull()
func (s *TypeSubject) actualNonNull() rhino.JSType {
	s.t.Helper()
	s.IsNotNull()
	return s.actual
}

// Java: TypeSubject.checkEqualityAgainst(JSType, boolean)
func (s *TypeSubject) checkEqualityAgainst(provided rhino.JSType, expectation bool) {
	s.t.Helper()
	if jutil.IsNil(provided) {
		provided = nil
	}
	providedString := debugStringOf(provided)
	actualString := debugStringOf(s.actual)

	actualEqualsProvided := jutil.ObjectsEquals(s.actual, provided)
	if actualEqualsProvided != expectation {
		s.fail("Types expected to be equal: %t\n%s: %t\nprovided: %s\nbut was: %s",
			expectation, equalsExpressionOf(actualString, providedString), actualEqualsProvided,
			providedString, actualString)
	}

	providedEqualsActual := jutil.ObjectsEquals(provided, s.actual)
	if actualEqualsProvided != providedEqualsActual {
		s.fail("Equality should be symmetric\n%s: %t\n%s: %t\nprovided: %s\nbut was: %s",
			equalsExpressionOf(actualString, providedString), actualEqualsProvided,
			equalsExpressionOf(providedString, actualString), providedEqualsActual,
			providedString, actualString)
	}

	if expectation {
		if jutil.ObjectHashCode(s.actual) != jutil.ObjectHashCode(provided) {
			s.fail("If two types are equal their hashcodes must also be equal\nhash of actual: %d\n"+
				"hash of provided: %d\nprovided: %s\nbut was: %s",
				jutil.ObjectHashCode(s.actual), jutil.ObjectHashCode(provided), providedString, actualString)
		}
	}
}

// Java: TypeSubject.debugStringOf(JSType)
func debugStringOf(type_ rhino.JSType) string {
	if jutil.IsNil(type_) {
		return "[Java null]"
	}
	return type_.String() + fmt.Sprintf(" [instanceof %T]", type_)
}

// Java: TypeSubject.equalsExpressionOf(Object, Object)
func equalsExpressionOf(receiver any, parameter any) string {
	return "(" + jutil.ToJavaString(receiver) + ").equals(" + jutil.ToJavaString(parameter) + ")"
}

// TypeSubject_FunctionTypeSubject holds the assertions that only apply to
// function types.
type TypeSubject_FunctionTypeSubject struct {
	outer *TypeSubject
}

func (f *TypeSubject_FunctionTypeSubject) actualFunctionType() rhino.FunctionType {
	f.outer.t.Helper()
	fn := f.outer.actualNonNull().ToMaybeFunctionType()
	if fn == nil {
		f.outer.fail("expected a function type, but was %s", debugStringOf(f.outer.actual))
	}
	return fn
}

// Java: TypeSubject.FunctionTypeSubject.hasTypeOfThisThat()
func (f *TypeSubject_FunctionTypeSubject) HasTypeOfThisThat() *TypeSubject {
	f.outer.t.Helper()
	return Truth{T: f.outer.t}.AssertType(f.actualFunctionType().GetTypeOfThis())
}

// Java: TypeSubject.FunctionTypeSubject.hasReturnTypeThat()
func (f *TypeSubject_FunctionTypeSubject) HasReturnTypeThat() *TypeSubject {
	f.outer.t.Helper()
	return Truth{T: f.outer.t}.AssertType(f.actualFunctionType().GetReturnType())
}

// Java: TypeSubject.FunctionTypeSubject.isConstructorFor(String)
func (f *TypeSubject_FunctionTypeSubject) IsConstructorFor(name string) {
	f.outer.t.Helper()
	f.outer.check("isConstructor()").That(f.outer.actual.IsConstructor()).IsTrue()
	f.outer.check("getInstanceType().getDisplayName()").
		That(f.actualFunctionType().GetInstanceType().GetDisplayName()).
		IsEqualTo(name)
}

// Java: TypeSubject.FunctionTypeSubject.isAbstract()
func (f *TypeSubject_FunctionTypeSubject) IsAbstract() {
	f.outer.t.Helper()
	f.outer.check("isAbstract()").That(f.actualFunctionType().IsAbstract()).IsTrue()
}

// Java: TypeSubject.FunctionTypeSubject.hasPrimitiveId(ClosurePrimitive)
func (f *TypeSubject_FunctionTypeSubject) HasPrimitiveId(id *rhino.ClosurePrimitive) {
	f.outer.t.Helper()
	f.outer.check("getClosurePrimitive()").
		That(f.outer.actualNonNull().ToMaybeFunctionType().GetClosurePrimitive()).
		IsEqualTo(id)
}
