package args4j

import "fmt"

// Go-only: language-adapter — reflection. Class is java.lang.Class as args4j
// inspects it: the type of the values a Setter accepts (Setter.getType()),
// which OptionHandlerRegistry maps to a handler, and the type of a Field. It
// answers the questions args4j asks of a class: its name and toString(),
// whether it is an enum (and its constants), and isAssignableFrom.
//
// Classes are compared by identity, like Java's. The classes of the Java
// types args4j and Closure's option beans use are predeclared below; an enum
// class is made with NewEnumClass.
type Class struct {
	name        string
	isInterface bool
	isPrimitive bool
	// isEnum is set for the class of a Java enum (a subclass of
	// java.lang.Enum); constants are then its constants in ordinal order.
	isEnum    bool
	constants []any
	// supertypes are the direct superclass and interfaces, besides
	// java.lang.Object.
	supertypes []*Class
}

// The classes of the Java types that the ported handlers and Closure's option
// beans use.
var (
	Class_Object     = &Class{name: "java.lang.Object"}
	Class_Enum       = &Class{name: "java.lang.Enum"}
	Class_String     = &Class{name: "java.lang.String"}
	Class_Boolean    = &Class{name: "java.lang.Boolean"}
	Class_Integer    = &Class{name: "java.lang.Integer"}
	Class_Collection = &Class{name: "java.util.Collection", isInterface: true}
	Class_List       = &Class{name: "java.util.List", isInterface: true, supertypes: []*Class{Class_Collection}}
	Class_Map        = &Class{name: "java.util.Map", isInterface: true}
	Class_boolean    = &Class{name: "boolean", isPrimitive: true}
	Class_int        = &Class{name: "int", isPrimitive: true}
)

// NewEnumClass is the Class of the Java enum named name (its binary name,
// e.g. "com.google.javascript.jscomp.CompilerOptions$DevMode") whose
// constants, in ordinal order, are constants. A constant's Java name() is
// its Name() method if it has one and its String() method otherwise (the
// convention of the ported enums: String() is toString(), which is name()
// unless the enum overrides it); its toString() is String().
func NewEnumClass(name string, constants []any) *Class {
	return &Class{name: name, isEnum: true, constants: constants, supertypes: []*Class{Class_Enum}}
}

// EnumClassOf is NewEnumClass for the values of a Go enum type, as returned
// by its Foo_values() function.
func EnumClassOf[E fmt.Stringer](name string, values []E) *Class {
	constants := make([]any, len(values))
	for i, v := range values {
		constants[i] = v
	}
	return NewEnumClass(name, constants)
}

// GetName is Class.getName().
func (c *Class) GetName() string { return c.name }

// String is Class.toString(): "interface <name>", "class <name>", or the
// name of a primitive type.
func (c *Class) String() string {
	if c.isInterface {
		return "interface " + c.name
	}
	if c.isPrimitive {
		return c.name
	}
	return "class " + c.name
}

// IsPrimitive is Class.isPrimitive().
func (c *Class) IsPrimitive() bool { return c.isPrimitive }

// IsArray is Class.isArray(). No Go Field has an array type (array fields
// are not ported), so it is always false.
func (c *Class) IsArray() bool { return false }

// GetEnumConstants is Class.getEnumConstants(): the constants of an enum
// class in ordinal order, or nil.
func (c *Class) GetEnumConstants() []any {
	if !c.isEnum {
		return nil
	}
	return append([]any(nil), c.constants...)
}

// IsAssignableFrom is Class.isAssignableFrom(Class): whether cls is c, or a
// subclass or subinterface of c (every class and interface is assignable to
// java.lang.Object, a primitive type only to itself).
func (c *Class) IsAssignableFrom(cls *Class) bool {
	if c == cls {
		return true
	}
	if cls.isPrimitive {
		return false
	}
	if c == Class_Object {
		return true
	}
	for _, s := range cls.supertypes {
		if c.IsAssignableFrom(s) {
			return true
		}
	}
	return false
}

// enumName is Enum.name() of an enum constant (see NewEnumClass).
func enumName(constant any) string {
	if n, ok := constant.(interface{ Name() string }); ok {
		return n.Name()
	}
	return constant.(fmt.Stringer).String()
}
