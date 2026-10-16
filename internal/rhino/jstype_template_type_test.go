// Port of com/google/javascript/rhino/jstype/TemplateTypeTest.java @ closure-compiler v20260415

package rhino_test

import (
	"testing"

	. "github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

type templateTypeTest struct {
	*rhinotesting.BaseJSTypeTestCase
}

func newTemplateTypeTest(t *testing.T) *templateTypeTest {
	s := &templateTypeTest{BaseJSTypeTestCase: rhinotesting.NewBaseJSTypeTestCase(t)}
	return s
}

func TestTemplateType_TemplateTypes_boundedByUnknown_areNotEqual(t *testing.T) {
	s := newTemplateTypeTest(t)
	unknownType := s.Registry.GetNativeType(JSTypeNative_UNKNOWN_TYPE)
	t1 := NewTemplateTypeWithBound(s.Registry, "T", unknownType)
	t2 := NewTemplateTypeWithBound(s.Registry, "T", unknownType)

	t1.Resolve(s.Registry.GetErrorReporter())
	t2.Resolve(s.Registry.GetErrorReporter())

	s.AssertType(t1).IsNotEqualTo(t2)
	s.AssertType(t1).IsSubtypeOf(t2)
}

func TestTemplateType_TemplateTypes_boundedByNominalTypes_areNotEqual(t *testing.T) {
	s := newTemplateTypeTest(t)
	fooType :=
		s.Registry.CreateObjectType(ptr("Foo"), s.Registry.GetNativeObjectType(JSTypeNative_OBJECT_TYPE))
	t1 := NewTemplateTypeWithBound(s.Registry, "T", fooType)
	t2 := NewTemplateTypeWithBound(s.Registry, "T", fooType)

	t1.Resolve(s.Registry.GetErrorReporter())
	t2.Resolve(s.Registry.GetErrorReporter())

	s.AssertType(t1).IsNotEqualTo(t2)
	s.AssertType(t1).IsNotSubtypeOf(t2)
}

func TestTemplateType_ProxiesOf_templateTypes_boundedByNominalTypes_areNotEqual(t *testing.T) {
	s := newTemplateTypeTest(t)
	fooType :=
		s.Registry.CreateObjectType(ptr("Foo"), s.Registry.GetNativeObjectType(JSTypeNative_OBJECT_TYPE))

	t1 := NewTemplateTypeWithBound(s.Registry, "T", fooType)
	t2 := NewTemplateTypeWithBound(s.Registry, "T", fooType)

	t1Proxy :=
		NamedType_Builder_(s.Registry, "T").
			SetResolutionKind(NamedType_ResolutionKind_NONE).
			SetReferencedType(t1).
			Build()
	t2Proxy :=
		NamedType_Builder_(s.Registry, "T").
			SetResolutionKind(NamedType_ResolutionKind_NONE).
			SetReferencedType(t2).
			Build()

	s.AssertType(t1Proxy).IsNotEqualTo(t2Proxy)
	s.AssertType(t1Proxy).IsNotSubtypeOf(t2Proxy)
}

func TestTemplateType_TemplateTypes_areNotEqualToUnknown(t *testing.T) {
	s := newTemplateTypeTest(t)
	// Given
	t_ := NewTemplateType(s.Registry, "T")
	s.AssertThat(t_.IsUnknownType()).IsTrue()

	// Then
	s.AssertType(t_).IsNotEqualTo(s.UNKNOWN_TYPE)
	s.AssertType(t_).IsNotEqualTo(s.CHECKED_UNKNOWN_TYPE)
}

func TestTemplateType_TemplateTypes_withBound_areNotEqualToBound(t *testing.T) {
	s := newTemplateTypeTest(t)
	// Given
	t_ := NewTemplateTypeWithBound(s.Registry, "T", s.NUMBER_TYPE)

	// Then
	s.AssertType(t_).IsNotEqualTo(s.NUMBER_TYPE)
}
