// Port of com/google/javascript/rhino/jstype/NamedTypeTest.java @ closure-compiler v20260415

package rhino_test

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	. "github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

type namedTypeTest struct {
	*rhinotesting.BaseJSTypeTestCase
	fooCtorType FunctionType
	fooType     ObjectType
}

func newNamedTypeTest(t *testing.T) *namedTypeTest {
	s := &namedTypeTest{BaseJSTypeTestCase: rhinotesting.NewBaseJSTypeTestCase(t)}
	s.SetUp()
	return s
}

func (s *namedTypeTest) SetUp() {
	closer := s.Registry.GetResolver().OpenForDefinition()
	defer closer.Close()
	{
		s.fooCtorType = FunctionType_Builder_(s.Registry).ForConstructor().WithName(ptr("Foo")).Build()
		s.fooType = s.fooCtorType.GetInstanceType()
	}
}

func TestNamedType_ResolutionPropagatesNamedTypePropertiesToResolvedType(t *testing.T) {
	s := newNamedTypeTest(t)
	// Given
	fooToFooScope := rhinotesting.NewMapBasedScope(mapOf("Foo", s.fooCtorType))
	var namedFooType *NamedType
	func() {
		closer := s.Registry.GetResolver().OpenForDefinition()
		defer closer.Close()
		namedFooType = s.NamedTypeBuilder("Foo").SetScope(fooToFooScope).Build()

		namedFooType.DefineDeclaredProperty("myProperty", s.NUMBER_TYPE, nil)

		// The property should not be typed yet.
		s.AssertTypeNotEquals(s.NUMBER_TYPE, s.fooType.GetPropertyType("myProperty"))
	}()

	// Then
	s.AssertTypeEquals(s.NUMBER_TYPE, s.fooType.GetPropertyType("myProperty"))
}

func TestNamedType_StateOfForwardDeclaredType_Unresolved(t *testing.T) {
	s := newNamedTypeTest(t)
	s.ResetRegistryWithForwardDeclaredName("forwardDeclared")

	closer := s.Registry.GetResolver().OpenForDefinition()
	defer closer.Close()
	{
		// Given
		type_ := s.NamedTypeBuilder("forwardDeclared").Build()

		// Then
		s.AssertThat(type_.IsResolved()).IsFalse()
		s.AssertThat(type_.IsEmptyType()).IsFalse()
		s.AssertTypeNotEquals(s.UNKNOWN_TYPE, type_)
		s.AssertTypeEquals(s.UNKNOWN_TYPE, type_.GetReferencedType())
	}
}

func TestNamedType_StateOfForwardDeclaredType_UnsuccesfullyResolved(t *testing.T) {
	s := newNamedTypeTest(t)
	s.ResetRegistryWithForwardDeclaredName("forwardDeclared")

	// Given
	type_ := s.NamedTypeBuilder("forwardDeclared").Build()

	// Then
	s.AssertThat(type_.IsUnsuccessfullyResolved()).IsTrue()
	s.AssertThat(type_.IsUnknownType()).IsFalse()
}

func TestNamedType_StateOfForwardDeclaredType_SuccessfullyResolved(t *testing.T) {
	s := newNamedTypeTest(t)
	// Given
	fooToFooScope := rhinotesting.NewMapBasedScope(mapOf("Foo", s.fooCtorType))
	namedFooType := s.NamedTypeBuilder("Foo").SetScope(fooToFooScope).Build()

	// When
	namedFooType.Resolve(ErrorReporter_NULL_INSTANCE)

	// Then
	s.AssertThat(namedFooType.IsSuccessfullyResolved()).IsTrue()
	s.AssertTypeEquals(namedFooType, s.fooType)
}

func TestNamedType_Equality(t *testing.T) {
	s := newNamedTypeTest(t)
	// Given
	barType := s.CreateNominalType("Bar" /* resolve= */)
	anonType := FunctionType_Builder_(s.Registry).Build()

	closer := s.Registry.GetResolver().OpenForDefinition()

	namedFooBuilder := s.NamedTypeBuilder("Foo")
	namedFooUnsuccessfullyResolved :=
		s.ForceResolutionWith(NewNoResolvedType(s.Registry, "Foo", nil), namedFooBuilder)
	namedFooResolvedToFoo := s.ForceResolutionWith(s.fooType, namedFooBuilder)
	namedFooResolvedToAnon := s.ForceResolutionWith(anonType, namedFooBuilder)
	namedFooResolvedToBar := s.ForceResolutionWith(barType, namedFooBuilder)

	namedBarBuilder := s.NamedTypeBuilder("Bar")
	namedBarUnsuccessfullyResolved :=
		s.ForceResolutionWith(NewNoResolvedType(s.Registry, "Bar", nil), namedBarBuilder)
	namedBarResolvedToFoo := s.ForceResolutionWith(s.fooType, namedBarBuilder)
	namedBarResolvedToAnon := s.ForceResolutionWith(anonType, namedBarBuilder)
	namedBarResolvedToBar := s.ForceResolutionWith(barType, namedBarBuilder)

	closer.Close()

	// Then
	rhinotesting.NewEqualsTester().
		AddEqualityGroup(s.fooType, namedFooResolvedToFoo, namedBarResolvedToFoo).
		AddEqualityGroup(barType, namedFooResolvedToBar, namedBarResolvedToBar).
		AddEqualityGroup(anonType, namedFooResolvedToAnon, namedBarResolvedToAnon).
		AddEqualityGroup(namedFooUnsuccessfullyResolved).
		AddEqualityGroup(namedBarUnsuccessfullyResolved).
		// TODO(b/112425334): Either re-enable this equality group or remove the NO_RESOLVED_TYPE
		// from the typesystem.
		// .addEqualityGroup(NO_RESOLVED_TYPE)
		TestEquals(s.T)
}

func TestNamedType_EqualityOfTypesWithSameReferenceName_postResolution(t *testing.T) {
	s := newNamedTypeTest(t)
	// Given
	barTypeA := s.CreateNominalType("Bar")
	barTypeB := s.CreateNominalType("Bar")
	anonType := FunctionType_Builder_(s.Registry).Build()

	closer := s.Registry.GetResolver().OpenForDefinition()

	namedFooBuilder := NamedType_Builder_(s.Registry, "Foo")
	namedFooUnsuccessfullyResolved :=
		s.ForceResolutionWith(NewNoResolvedType(s.Registry, "Foo", nil), namedFooBuilder)
	namedFooResolvedToFoo := s.ForceResolutionWith(s.fooType, namedFooBuilder)
	namedFooResolvedToAnon := s.ForceResolutionWith(anonType, namedFooBuilder)
	namedFooResolvedToBarA := s.ForceResolutionWith(barTypeA, namedFooBuilder)
	namedFooResolvedToBarB := s.ForceResolutionWith(barTypeB, namedFooBuilder)
	namedFooResolvedToNumber := s.ForceResolutionWith(s.NUMBER_TYPE, namedFooBuilder)
	namedFooResolvedToString := s.ForceResolutionWith(s.STRING_TYPE, namedFooBuilder)

	namedBarBuilder := NamedType_Builder_(s.Registry, "Bar")
	namedBarUnsuccessfullyResolved :=
		s.ForceResolutionWith(NewNoResolvedType(s.Registry, "Bar", nil), namedBarBuilder)
	namedBarResolvedToFoo := s.ForceResolutionWith(s.fooType, namedBarBuilder)
	namedBarResolvedToAnon := s.ForceResolutionWith(anonType, namedBarBuilder)
	namedBarResolvedToBarA := s.ForceResolutionWith(barTypeA, namedBarBuilder)
	namedBarResolvedToBarB := s.ForceResolutionWith(barTypeB, namedBarBuilder)
	namedBarResolvedToNumber := s.ForceResolutionWith(s.NUMBER_TYPE, namedBarBuilder)
	namedBarResolvedToString := s.ForceResolutionWith(s.STRING_TYPE, namedBarBuilder)

	closer.Close()

	// Then
	rhinotesting.NewEqualsTester().
		AddEqualityGroup(s.fooType, namedFooResolvedToFoo, namedBarResolvedToFoo).
		AddEqualityGroup(barTypeA, namedFooResolvedToBarA, namedBarResolvedToBarA).
		AddEqualityGroup(barTypeB, namedFooResolvedToBarB, namedBarResolvedToBarB).
		AddEqualityGroup(anonType, namedFooResolvedToAnon, namedBarResolvedToAnon).
		AddEqualityGroup(s.NUMBER_TYPE, namedFooResolvedToNumber, namedBarResolvedToNumber).
		AddEqualityGroup(s.STRING_TYPE, namedFooResolvedToString, namedBarResolvedToString).
		AddEqualityGroup(namedFooUnsuccessfullyResolved).
		AddEqualityGroup(namedBarUnsuccessfullyResolved).
		// TODO(b/112425334): Either re-enable this equality group or remove the NO_RESOLVED_TYPE
		// from the typesystem.
		// .addEqualityGroup(NO_RESOLVED_TYPE)
		TestEquals(s.T)
}

func TestNamedType_ForwardDeclaredNamedType(t *testing.T) {
	s := newNamedTypeTest(t)
	s.ErrorReporter.ExpectAllWarnings("Bad type annotation. Unknown type Unresolvable")
	a := s.NamedTypeBuilder("Unresolvable").Build()

	s.AssertTypeEquals(s.UNKNOWN_TYPE, a.GetLeastSupertype(s.UNKNOWN_TYPE))
	s.AssertTypeEquals(s.CHECKED_UNKNOWN_TYPE, a.GetLeastSupertype(s.CHECKED_UNKNOWN_TYPE))
	s.AssertTypeEquals(s.UNKNOWN_TYPE, s.UNKNOWN_TYPE.GetLeastSupertype(a))
	s.AssertTypeEquals(s.CHECKED_UNKNOWN_TYPE, s.CHECKED_UNKNOWN_TYPE.GetLeastSupertype(a))
}

func TestNamedType_ActiveXObjectResolve(t *testing.T) {
	s := newNamedTypeTest(t)
	var activeXObject *NamedType
	func() {
		closer := s.Registry.GetResolver().OpenForDefinition()
		defer closer.Close()
		scope := rhinotesting.NewMapBasedScope(mapOf("ActiveXObject", s.NO_OBJECT_TYPE))
		activeXObject = s.NamedTypeBuilder("ActiveXObject").SetScope(scope).Build()

		s.AssertThat(activeXObject.String()).IsEqualTo("ActiveXObject")
	}()
	s.AssertThat(activeXObject.String()).IsEqualTo("NoObject")
	s.AssertTypeEquals(s.NO_OBJECT_TYPE, activeXObject.GetReferencedType())
}

func TestNamedType_ResolveToGoogModule(t *testing.T) {
	s := newNamedTypeTest(t)
	s.Registry.RegisterNonLegacyClosureNamespace("mod.Foo", nil, s.fooCtorType)
	modDotFoo := s.NamedTypeBuilder("mod.Foo").Build()

	modDotFoo.Resolve(ErrorReporter_NULL_INSTANCE)
	s.AssertTypeEquals(s.fooType, modDotFoo.GetReferencedType())
}

func TestNamedType_ResolveToPropertyOnGoogModule(t *testing.T) {
	s := newNamedTypeTest(t)
	exportsType := s.Registry.CreateRecordType(mapOf("Foo", s.fooCtorType))
	s.Registry.RegisterNonLegacyClosureNamespace("mod.Bar", nil, exportsType)
	modBarFoo := s.NamedTypeBuilder("mod.Bar.Foo").Build()

	modBarFoo.Resolve(ErrorReporter_NULL_INSTANCE)
	s.AssertTypeEquals(s.fooType, modBarFoo.GetReferencedType())
}

func TestNamedType_ResolveToGoogModule_failsIfLegacyNamespaceIsLongerPrefix(t *testing.T) {
	s := newNamedTypeTest(t)
	s.ErrorReporter.ExpectAllWarnings("Bad type annotation. Unknown type mod.Bar.Foo")
	var modDotFoo *NamedType
	func() {
		closer := s.Registry.GetResolver().OpenForDefinition()
		defer closer.Close()
		exportsType := s.Registry.CreateRecordType(mapOf("Foo", s.fooCtorType))
		s.Registry.RegisterNonLegacyClosureNamespace("mod.Bar", nil, exportsType)
		modDotFoo = s.NamedTypeBuilder("mod.Bar.Foo").Build()
		s.Registry.RegisterLegacyClosureNamespace("mod.Bar.Foo")
	}()

	s.AssertType(modDotFoo.GetReferencedType()).IsUnknown()
}

func TestNamedType_ResolveToGoogModule_usesLongestModulePrefix(t *testing.T) {
	s := newNamedTypeTest(t)
	barType := s.CreateNominalType("Bar" /* resolve= */)
	s.Registry.RegisterNonLegacyClosureNamespace(
		"mod.Foo", nil, s.Registry.CreateRecordType(mapOf("Bar", s.fooCtorType)))
	s.Registry.RegisterNonLegacyClosureNamespace("mod.Foo.Bar", nil, barType.GetConstructor())
	scope := rhinotesting.NewMapBasedScope(mapOf())
	modDotFoo := s.NamedTypeBuilder("mod.Foo.Bar").SetScope(scope).Build()

	modDotFoo.Resolve(ErrorReporter_NULL_INSTANCE)
	s.AssertTypeEquals(barType, modDotFoo.GetReferencedType())
}

func TestNamedType_ReferenceGoogModuleByType_resolvesToModuleRatherThanRegistryType(t *testing.T) {
	s := newNamedTypeTest(t)
	// Note: this logic was put in place to mimic the semantics of type resolution when
	// modules are rewritten before typechecking.
	s.Registry.RegisterNonLegacyClosureNamespace("mod.Foo", nil, s.fooCtorType)
	scope := rhinotesting.NewMapBasedScope(mapOf())
	modDotFoo := s.NamedTypeBuilder("mod.Foo").SetScope(scope).Build()
	s.Registry.DeclareType(scope, "mod.Foo", s.CreateNominalType("other.mod.Foo"))

	s.AssertTypeEquals(s.fooType, modDotFoo.GetReferencedType())
}

func TestNamedType_GetBangType_onUnresolvedType(t *testing.T) {
	s := newNamedTypeTest(t)
	var barNonNull JSType
	func() {
		closer := s.Registry.GetResolver().OpenForDefinition()
		defer closer.Close()
		barType := s.NamedTypeBuilder("Bar").Build()

		barNonNull = barType.GetBangType()
		s.AssertThat(barNonNull.IsResolved()).IsFalse()

		s.Registry.DeclareType(nil, "Bar", s.Registry.CreateUnionType(s.fooType, s.NULL_TYPE))
	}()

	s.AssertType(barNonNull).ToStringIsEqualTo("Foo")
}

func TestNamedType_GetBangType_onResolvedType(t *testing.T) {
	s := newNamedTypeTest(t)
	var barType *NamedType
	func() {
		closer := s.Registry.GetResolver().OpenForDefinition()
		defer closer.Close()
		barType = s.NamedTypeBuilder("Bar").Build()
		s.Registry.DeclareType(nil, "Bar", s.Registry.CreateUnionType(s.fooType, s.NULL_TYPE))
	}()

	barNonNull := barType.GetBangType()
	s.AssertType(barNonNull).ToStringIsEqualTo("Foo")
}

func TestNamedType_BuilderForTypeof_emitsUnrecognizedTypeError(t *testing.T) {
	s := newNamedTypeTest(t)
	typeofFooBuilder := NamedType_Builder_(s.Registry, "typeof Foo").
		SetScope(rhinotesting.MapBasedScope_EmptyScope()).
		SetResolutionKind(NamedType_ResolutionKind_TYPEOF)

	s.AssertType(typeofFooBuilder.Build()).IsUnknown()
	s.ErrorReporter.ExpectAllWarnings(
		"Missing type for `typeof` value. The value must be declared and const.")
}

func TestNamedType_BuilderForTypeof_yieldsNoResolvedTypeWithForwardDeclaredName(t *testing.T) {
	s := newNamedTypeTest(t)
	s.ResetRegistryWithForwardDeclaredName("Foo")

	typeofFooBuilder := NamedType_Builder_(s.Registry, "typeof Foo").
		SetScope(rhinotesting.MapBasedScope_EmptyScope()).
		SetResolutionKind(NamedType_ResolutionKind_TYPEOF)

	s.AssertType(typeofFooBuilder.Build()).IsNoResolvedType("typeof Foo")
}

func TestNamedType_BuilderForTypeof_requiresReferenceToStartWithTypeof(t *testing.T) {
	s := newNamedTypeTest(t)
	typeofFooBuilder := NamedType_Builder_(s.Registry, "Foo").
		SetScope(rhinotesting.MapBasedScope_EmptyScope()).
		SetResolutionKind(NamedType_ResolutionKind_TYPEOF)

	rhinotesting.AssertThrows[jutil.Throwable](s.T, func() { typeofFooBuilder.Build() })
}

func (s *namedTypeTest) ForceResolutionWith(type_ JSType, proxy *NamedType_Builder) *NamedType {
	return proxy.SetReferencedType(type_).SetResolutionKind(NamedType_ResolutionKind_NONE).Build()
}

func (s *namedTypeTest) CreateNominalType(name string) ObjectType {
	closer := s.Registry.GetResolver().OpenForDefinition()
	defer closer.Close()
	{
		ctorType :=
			FunctionType_Builder_(s.Registry).ForConstructor().WithName(ptr(name)).Build()
		return ctorType.GetInstanceType()
	}
}

func (s *namedTypeTest) NamedTypeBuilder(name string) *NamedType_Builder {
	return NamedType_Builder_(s.Registry, name).
		SetResolutionKind(NamedType_ResolutionKind_TYPE_NAME).
		SetScope(rhinotesting.MapBasedScope_EmptyScope()).
		SetErrorReportingLocation(ptr("source"), 1, 0)
}
