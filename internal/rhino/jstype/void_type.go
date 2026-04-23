// Port of com/google/javascript/rhino/jstype/VoidType.java @ closure-compiler v20260415

package jstype

// voidType is the JS void / undefined type.  It is a subtype of itself, All,
// and Unknown.  TypeName returns "undefined" to match Java's getDisplayName().
type voidType struct{}

func (voidType) TypeName() string { return "undefined" }
func (voidType) IsSubtypeOf(o JSType) bool {
	if _, ok := o.(voidType); ok {
		return true
	}
	return isAll(o) || isUnknown(o)
}
