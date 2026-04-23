// Port of com/google/javascript/rhino/jstype/NumberType.java @ closure-compiler v20260415

package jstype

// numberType is the primitive JS number type.  It is a subtype of itself,
// All, and Unknown.
type numberType struct{}

func (numberType) TypeName() string { return "number" }
func (numberType) IsSubtypeOf(o JSType) bool {
	if _, ok := o.(numberType); ok {
		return true
	}
	return isAll(o) || isUnknown(o)
}
