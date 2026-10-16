package rhino

import "testing"

func TestJSTypeSmoke(t *testing.T) {
	r := NewJSTypeRegistry(ErrorReporter_NULL_INSTANCE)
	for _, n := range JSTypeNative_values() {
		ty := r.GetNativeType(n)
		t.Logf("%s: %s / %s", n, ty.String(), ty.ToAnnotationString(JSType_Nullability_EXPLICIT))
	}
	num := r.GetNativeType(JSTypeNative_NUMBER_TYPE)
	str := r.GetNativeType(JSTypeNative_STRING_TYPE)
	u := r.CreateUnionType(num, str, r.GetNativeType(JSTypeNative_NULL_TYPE))
	t.Logf("union %s sub=%v", u, num.IsSubtypeOf(u))
	f := r.CreateFunctionType(num, str, u)
	t.Logf("fn %s", f)
}
