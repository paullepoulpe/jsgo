package args4j

import "github.com/paullepoulpe/jsgo/internal/jutil"

// Go-only: language-adapter — reflection. Field is java.lang.reflect.Field
// as args4j's setters use it: the declared type of an option bean's field
// (and, for a parameterized List, its type argument), whether it is final,
// and accessors bound to the bean that read and write it with values boxed
// as args4j sees them (see the package documentation). The bean argument of
// Get and Set is the bean that FieldSetter and MultiValueFieldSetter carry;
// the accessors are already bound to it.
type Field struct {
	declaringClass string
	name           string
	type_          *Class
	typeArgument   *Class
	final          bool
	get            func() any
	set            func(value any)
	setString      func(value string)
	addString      func(value string)
}

// NewField is a Field of the class named declaringClass (its binary name),
// called name, of declared type type_; typeArgument is the class of the
// first type argument of a parameterized type (List<String>: Class_String),
// or nil. get returns the field's value (nil for Java null); set stores one.
func NewField(declaringClass, name string, type_, typeArgument *Class, final bool, get func() any, set func(value any)) *Field {
	return &Field{declaringClass: declaringClass, name: name, type_: type_, typeArgument: typeArgument, final: final, get: get, set: set}
}

// GetDeclaringClassName is Field.getDeclaringClass().getName().
func (f *Field) GetDeclaringClassName() string { return f.declaringClass }

// GetName is Field.getName().
func (f *Field) GetName() string { return f.name }

// GetType is Field.getType().
func (f *Field) GetType() *Class { return f.type_ }

// GetGenericTypeArgument is, for a field whose Field.getGenericType() is a
// ParameterizedType whose first actual type argument is a Class, that Class;
// otherwise nil.
func (f *Field) GetGenericTypeArgument() *Class { return f.typeArgument }

// IsFinal is Modifier.isFinal(Field.getModifiers()).
func (f *Field) IsFinal() bool { return f.final }

// Get is Field.get(Object bean).
func (f *Field) Get(bean any) any { return f.get() }

// Set is Field.set(Object bean, Object value).
func (f *Field) Set(bean any, value any) { f.set(value) }

// String is Field.toString() without the modifiers, which a Go Field does
// not record: "<type> <declaring class>.<name>".
func (f *Field) String() string {
	return f.type_.GetName() + " " + f.declaringClass + "." + f.name
}

func (*Field) annotatedElement() {}

// List is the raw java.util.List that MultiValueFieldSetter adds the values
// of a multi-valued option to, and reads its current values from.
// *jutil.ArrayList[any] is one.
type List interface {
	Add(e any) bool
	Size() int
	Get(index int) any
}

// Map is the raw java.util.Map that MapOptionHandler fills: it puts String
// keys and String values, nil being null. *jutil.JavaHashMap[string, *string]
// (a Java HashMap<String, String>) is one.
type Map interface {
	Put(key string, value *string) *string
}

// valueField is a Field over a Go variable of type T that holds the boxed
// value itself.
func valueField[T any](declaringClass, name string, type_ *Class, p *T) *Field {
	return NewField(declaringClass, name, type_, nil, false,
		func() any { return *p },
		func(v any) { *p = v.(T) })
}

// nullableField is a Field over a Go variable that holds a pointer to the
// boxed value, nil being Java null.
func nullableField[T any](declaringClass, name string, type_ *Class, p **T) *Field {
	return NewField(declaringClass, name, type_, nil, false,
		func() any {
			if *p == nil {
				return nil
			}
			return **p
		},
		func(v any) {
			if v == nil {
				*p = nil
				return
			}
			x := v.(T)
			*p = &x
		})
}

// BooleanField is a field of Java type boolean.
func BooleanField(declaringClass, name string, p *bool) *Field {
	return valueField(declaringClass, name, Class_boolean, p)
}

// NullableBooleanField is a field of Java type Boolean.
func NullableBooleanField(declaringClass, name string, p **bool) *Field {
	return nullableField(declaringClass, name, Class_Boolean, p)
}

// IntField is a field of Java type int held in a Go int; its values are
// boxed as int32 (Integer).
func IntField(declaringClass, name string, p *int) *Field {
	return NewField(declaringClass, name, Class_int, nil, false,
		func() any { return int32(*p) },
		func(v any) { *p = int(v.(int32)) })
}

// IntegerField is a field of Java type Integer that is never null.
func IntegerField(declaringClass, name string, p *int32) *Field {
	return valueField(declaringClass, name, Class_Integer, p)
}

// NullableIntegerField is a field of Java type Integer.
func NullableIntegerField(declaringClass, name string, p **int32) *Field {
	return nullableField(declaringClass, name, Class_Integer, p)
}

// StringField is a field of Java type String that is never null.
func StringField(declaringClass, name string, p *string) *Field {
	f := valueField(declaringClass, name, Class_String, p)
	f.setString = func(v string) { *p = v }
	return f
}

// NullableStringField is a field of Java type String.
func NullableStringField(declaringClass, name string, p **string) *Field {
	f := nullableField(declaringClass, name, Class_String, p)
	f.setString = func(v string) { *p = &v }
	return f
}

// EnumField is a field of the Java enum type enumClass (see NewEnumClass)
// that is never null; E is the Go type of its constants.
func EnumField[E any](declaringClass, name string, enumClass *Class, p *E) *Field {
	return valueField(declaringClass, name, enumClass, p)
}

// NullableEnumField is a field of the Java enum type enumClass.
func NullableEnumField[E any](declaringClass, name string, enumClass *Class, p **E) *Field {
	return nullableField(declaringClass, name, enumClass, p)
}

// ListField is a field of Java type List<X>, where elementType is the class
// of X and E the Go type of its boxed values, held in a Go slice. Its value
// is never null: a nil slice is an empty list, and the list's add appends to
// the slice. Setting the field replaces the slice with the elements of the
// List it is given.
func ListField[E any](declaringClass, name string, elementType *Class, p *[]E) *Field {
	sl := &sliceList[E]{p}
	f := NewField(declaringClass, name, Class_List, elementType, false,
		func() any { return sl },
		func(v any) {
			l := v.(List)
			s := make([]E, l.Size())
			for i := range s {
				s[i] = l.Get(i).(E)
			}
			*p = s
		})
	if sp, ok := any(p).(*[]string); ok {
		f.addString = func(v string) {
			*sp = append(*sp, v)
		}
	}
	return f
}

// sliceList is the List view of a Go slice variable.
type sliceList[E any] struct{ p *[]E }

func (l *sliceList[E]) Add(e any) bool {
	*l.p = append(*l.p, e.(E))
	return true
}

func (l *sliceList[E]) Size() int { return len(*l.p) }

func (l *sliceList[E]) Get(index int) any { return (*l.p)[index] }

// StringMapField is a field of Java type Map<String, String> holding a Java
// HashMap (or null).
func StringMapField(declaringClass, name string, p **jutil.JavaHashMap[string, *string]) *Field {
	return NewField(declaringClass, name, Class_Map, Class_String, false,
		func() any {
			if *p == nil {
				return nil
			}
			return *p
		},
		func(v any) {
			if v == nil {
				*p = nil
				return
			}
			*p = v.(*jutil.JavaHashMap[string, *string])
		})
}
