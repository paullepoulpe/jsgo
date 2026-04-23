// Port of com/google/javascript/rhino/jstype/BooleanType.java @ closure-compiler v20260415

package jstype

// booleanType is the primitive JS boolean type.  It is a subtype of itself,
// All, and Unknown.
type booleanType struct{}

func (booleanType) TypeName() string { return "boolean" }
func (booleanType) IsSubtypeOf(o JSType) bool {
	if _, ok := o.(booleanType); ok {
		return true
	}
	return isAll(o) || isUnknown(o)
}
