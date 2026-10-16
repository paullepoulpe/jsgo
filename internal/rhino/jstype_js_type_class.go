// Port of com/google/javascript/rhino/jstype/JSTypeClass.java @ closure-compiler v20260415

package rhino

// JSTypeClass is the set of Java classes that are subtypes of JSType.
//
// Each class reports the corresponding element as its instance
// JSType#getTypeClass method.
//
// These values are not intended to be compatible with Liskov substitution;
// JSTypes have exactly one associated JSTypeClass and it must not affect the
// outcome of operations such as Object#equals. Use of this enum is
// tantamount to reflection (e.g. Object#getClass or instanceof) and should
// only be used where such a comparison is appropriate and desirable.
//
// Use of this enum is preferred to true reflection, when possible, for
// simplicity. This is safe because this package controls all JSType
// subclasses; they can all be known a priori.
//
// Abstract classes are not listed here because there would be no way to use
// those elements.
type JSTypeClass int32

const (
	JSTypeClass_ALL JSTypeClass = iota
	JSTypeClass_ARROW
	JSTypeClass_BOOLEAN
	JSTypeClass_BIGINT
	JSTypeClass_ENUM
	JSTypeClass_ENUM_ELEMENT
	JSTypeClass_FUNCTION
	JSTypeClass_INSTANCE_OBJECT
	JSTypeClass_NAMED
	JSTypeClass_NO
	JSTypeClass_NO_OBJECT
	JSTypeClass_NO_RESOLVED
	JSTypeClass_NULL
	JSTypeClass_NUMBER
	JSTypeClass_PROTOTYPE_OBJECT
	JSTypeClass_PROXY_OBJECT
	JSTypeClass_RECORD
	JSTypeClass_STRING
	JSTypeClass_SYMBOL
	JSTypeClass_WELL_KNOWN_SYMBOL
	JSTypeClass_TEMPLATE
	JSTypeClass_TEMPLATIZED
	JSTypeClass_UNION
	JSTypeClass_UNKNOWN
	JSTypeClass_VOID
)

var jsTypeClass_names = [...]string{
	"ALL", "ARROW", "BOOLEAN", "BIGINT", "ENUM", "ENUM_ELEMENT", "FUNCTION", "INSTANCE_OBJECT",
	"NAMED", "NO", "NO_OBJECT", "NO_RESOLVED", "NULL", "NUMBER", "PROTOTYPE_OBJECT",
	"PROXY_OBJECT", "RECORD", "STRING", "SYMBOL", "WELL_KNOWN_SYMBOL", "TEMPLATE", "TEMPLATIZED",
	"UNION", "UNKNOWN", "VOID",
}

// JSTypeClass_values is JSTypeClass.values().
func JSTypeClass_values() []JSTypeClass {
	out := make([]JSTypeClass, len(jsTypeClass_names))
	for i := range out {
		out[i] = JSTypeClass(i)
	}
	return out
}

// String is Enum.name().
func (c JSTypeClass) String() string { return jsTypeClass_names[c] }

// Java: JSTypeClass.isTypeOf(JSType)
func (c JSTypeClass) IsTypeOf(t JSType) bool {
	return c == t.GetTypeClass()
}
