// Go-only: test of the TemplateTypeMap fast path (resolveSimpleTemplateValue):
// a function-typed template value must go through the TemplateTypeReplacer,
// as Java always does, so that its signature is resolved.

package rhino_test

import (
	"testing"

	. "github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestTemplateTypeMap_FastPathResolvesFunctionSignatures(t *testing.T) {
	for _, tc := range []struct {
		name             string
		templateReturn   bool
		templateArgument bool
	}{
		{"return", true, false},
		{"argument", false, true},
		{"returnAndArgument", true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := newTemplateTypeMapTest(t)
			key := NewTemplateType(s.Registry, "T")
			functionKey := NewTemplateType(s.Registry, "F")
			number := s.Registry.GetNativeType(JSTypeNative_NUMBER_TYPE)
			text := s.Registry.GetNativeType(JSTypeNative_STRING_TYPE)
			returnType, argumentType := text, text
			wantReturn, wantArgument := text, text
			if tc.templateReturn {
				returnType, wantReturn = key, number
			}
			if tc.templateArgument {
				argumentType, wantArgument = key, number
			}
			function := s.Registry.CreateFunctionType(returnType, argumentType)
			// Separate extensions make T visible to F at the next-submap
			// boundary. Both signatures are ordinary built-in function types.
			outer := s.emptyMap.CopyWithExtensionOfKeys([]*TemplateType{functionKey}, []JSType{function})
			resolved := outer.CopyWithExtensionOfKeys([]*TemplateType{key}, []JSType{number})
			got, ok := resolved.GetResolvedTemplateType(functionKey).(FunctionType)
			if !ok || got.GetReturnType() != wantReturn || len(got.GetParameters()) != 1 || got.GetParameters()[0].GetJSType() != wantArgument {
				t.Fatal("template-map fast path skipped substitution in a function signature")
			}
			if got == function || function.GetReturnType() != returnType || function.GetParameters()[0].GetJSType() != argumentType {
				t.Fatal("function-signature resolution failed to preserve the original function type")
			}
		})
	}
}
