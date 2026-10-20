package rhino_test

import (
	"testing"

	. "github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestTemplateTypeMap_FastPathNormalizesConstructedPrimitives(t *testing.T) {
	cases := []struct {
		name   string
		create func(*JSTypeRegistry) JSType
		native JSTypeNative
	}{
		{"number", func(r *JSTypeRegistry) JSType { return NewNumberType(r) }, JSTypeNative_NUMBER_TYPE},
		{"string", func(r *JSTypeRegistry) JSType { return NewStringType(r) }, JSTypeNative_STRING_TYPE},
		{"boolean", func(r *JSTypeRegistry) JSType { return NewBooleanType(r) }, JSTypeNative_BOOLEAN_TYPE},
		{"null", func(r *JSTypeRegistry) JSType { return NewNullType(r) }, JSTypeNative_NULL_TYPE},
		{"void", func(r *JSTypeRegistry) JSType { return NewVoidType(r) }, JSTypeNative_VOID_TYPE},
		{"all", func(r *JSTypeRegistry) JSType { return NewAllType(r) }, JSTypeNative_ALL_TYPE},
		{"bigint", func(r *JSTypeRegistry) JSType { return NewBigIntType(r) }, JSTypeNative_BIGINT_TYPE},
		{"noObject", func(r *JSTypeRegistry) JSType { return NewNoObjectType(r) }, JSTypeNative_NO_OBJECT_TYPE},
		{"unknown", func(r *JSTypeRegistry) JSType { return NewUnknownType(r, false) }, JSTypeNative_UNKNOWN_TYPE},
		{"checkedUnknown", func(r *JSTypeRegistry) JSType { return NewUnknownType(r, true) }, JSTypeNative_UNKNOWN_TYPE},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := newTemplateTypeMapTest(t)
			input := tc.create(s.Registry)
			want := s.Registry.GetNativeType(tc.native)
			if input == want {
				t.Fatal("fixture must construct a distinct primitive")
			}
			key := NewTemplateType(s.Registry, "T")
			m := s.emptyMap.CopyWithExtensionOfKeys([]*TemplateType{key}, []JSType{input})
			if got := m.GetResolvedTemplateType(key); got != want {
				t.Fatalf("resolved value = %T %p, want canonical %T %p", got, got, want, want)
			}
		})
	}
}
