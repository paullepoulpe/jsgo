// Port of com/google/javascript/rhino/jstype/PrototypeObjectTypeTest.java @ closure-compiler v20260415

package rhino_test

import (
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

type prototypeObjectTypeTest struct {
	*rhinotesting.BaseJSTypeTestCase
}

func newPrototypeObjectTypeTest(t *testing.T) *prototypeObjectTypeTest {
	s := &prototypeObjectTypeTest{BaseJSTypeTestCase: rhinotesting.NewBaseJSTypeTestCase(t)}
	return s
}

func TestPrototypeObjectType_ToString(t *testing.T) {
	s := newPrototypeObjectTypeTest(t)
	type_ := s.Registry.CreateAnonymousObjectType(nil)
	s.AssertThat(type_.String()).IsEqualTo("{}")

	type_.DefineDeclaredProperty("foo", s.NUMBER_TYPE, nil)
	s.AssertThat(type_.String()).IsEqualTo("{foo: number}")

	type_.DefineDeclaredProperty("bar", type_, nil)
	s.AssertThat(type_.String()).
		IsEqualTo(
			strings.Join([]string{
				"{", //
				"  bar: {...},",
				"  foo: number",
				"}"}, "\n"))
}
