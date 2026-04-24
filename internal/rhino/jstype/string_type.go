// Port of com/google/javascript/rhino/jstype/StringType.java @ closure-compiler v20260415

package jstype

// stringType is the primitive JS string type.  It is a subtype of itself,
// All, and Unknown.
type stringType struct{}

// Go-only: language-adapter
func (stringType) TypeName() string { return "string" }

// Go-only: language-adapter
func (stringType) IsSubtypeOf(o JSType) bool {
	if _, ok := o.(stringType); ok {
		return true
	}
	return isAll(o) || isUnknown(o)
}
