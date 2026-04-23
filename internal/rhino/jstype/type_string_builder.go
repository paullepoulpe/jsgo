// Port of com/google/javascript/rhino/jstype/TypeStringBuilder.java
// @ closure-compiler v20260415
//
// Visits types to assemble an associated stringification.
//
// Instances are single use. They irreversibly accumulate state required
// during traversal.
//
// The Original Code is Rhino code, released May 6, 1999.
// The Initial Developer of the Original Code is Netscape Communications
// Corporation. Portions created by the Initial Developer are Copyright
// (C) 1997-1999 the Initial Developer. All Rights Reserved.
//
// Contributor(s): Bob Jervis, Google Inc.
//
// Licensed under the Mozilla Public License Version 1.1 (the "License")
// or, at your option, the GNU General Public License Version 2 or later.

package jstype

import "strings"

// typeStringBuilderAppender is the Go analogue of the Java abstract
// method JSType.appendTo(TypeStringBuilder). Types that want richer
// output (e.g. UnionType, FunctionType) implement this interface; all
// others fall back to TypeName().
type typeStringBuilderAppender interface {
	appendTo(sb *TypeStringBuilder)
}

// TypeStringBuilder visits types to assemble an associated
// stringification. Port of Java's final class TypeStringBuilder.
// Instances are single use; they irreversibly accumulate state
// required during traversal.
type TypeStringBuilder struct {
	builder          strings.Builder
	isForAnnotations bool
	indentation      string
}

// NewTypeStringBuilder constructs a TypeStringBuilder. When
// isForAnnotations is true the builder emits "!" prefixes and other
// annotation-specific syntax; mirrors the Java constructor.
func NewTypeStringBuilder(isForAnnotations bool) *TypeStringBuilder {
	return &TypeStringBuilder{isForAnnotations: isForAnnotations}
}

// IsForAnnotations reports whether this builder is rendering
// annotation-style output. Mirrors Java isForAnnotations().
func (b *TypeStringBuilder) IsForAnnotations() bool {
	return b.isForAnnotations
}

// Append appends the raw string x to the buffer. Mirrors Java
// TypeStringBuilder.append(String).
func (b *TypeStringBuilder) Append(x string) *TypeStringBuilder {
	b.builder.WriteString(x)
	return b
}

// AppendType appends the stringification of x. If x implements
// typeStringBuilderAppender, that method is called (mirrors
// Java's x.appendTo(this)); otherwise TypeName() is used as the
// fallback, which is correct for all primitive/leaf types.
// Mirrors Java TypeStringBuilder.append(JSType).
func (b *TypeStringBuilder) AppendType(x JSType) *TypeStringBuilder {
	if a, ok := x.(typeStringBuilderAppender); ok {
		a.appendTo(b)
	} else {
		b.builder.WriteString(x.TypeName())
	}
	return b
}

// AppendAll appends each element from elements, separated by separator.
// Elements may be JSType values or strings; any other concrete type
// panics (matches Java's unchecked cast). Mirrors Java
// TypeStringBuilder.appendAll(Iterable<?>, String).
func (b *TypeStringBuilder) AppendAll(elements []any, separator string) *TypeStringBuilder {
	separate := false
	for _, e := range elements {
		if separate {
			b.Append(separator)
		} else {
			separate = true
		}
		switch v := e.(type) {
		case JSType:
			b.AppendType(v)
		case string:
			b.Append(v)
		default:
			panic("TypeStringBuilder.AppendAll: unsupported element type")
		}
	}
	return b
}

// AppendNonNull appends type with an optional "!" prefix when rendering
// for annotations and type is a non-null object (i.e. not Unknown,
// TemplateType, record, FunctionType, UnionType, or anonymous object
// literal). Mirrors Java TypeStringBuilder.appendNonNull(JSType).
func (b *TypeStringBuilder) AppendNonNull(t JSType) *TypeStringBuilder {
	if b.isForAnnotations && isObjectType(t) && !isUnknown(t) && !isTemplateType(t) &&
		!isRecordType(t) && !isFunctionType(t) && !isUnionType(t) && !isLiteralObject(t) {
		b.Append("!")
	}
	return b.AppendType(t)
}

// BreakLineAndIndent appends a newline followed by the current
// indentation string. Mirrors Java TypeStringBuilder.breakLineAndIndent().
func (b *TypeStringBuilder) BreakLineAndIndent() *TypeStringBuilder {
	b.builder.WriteString("\n")
	b.builder.WriteString(b.indentation)
	return b
}

// Indent increases the indentation by two spaces for the duration of
// cb, then restores it. Mirrors Java TypeStringBuilder.indent(Runnable).
func (b *TypeStringBuilder) Indent(cb func()) *TypeStringBuilder {
	last := b.indentation
	b.indentation = last + "  "
	cb()
	b.indentation = last
	return b
}

// Build returns the assembled string. Mirrors Java TypeStringBuilder.build().
func (b *TypeStringBuilder) Build() string {
	return b.builder.String()
}

// CloneWithConfig returns a new TypeStringBuilder that shares the same
// isForAnnotations flag and current indentation but has an empty
// buffer. Used by UnionType.appendTo to stringify individual members
// in isolation before sorting. Mirrors Java TypeStringBuilder.cloneWithConfig().
func (b *TypeStringBuilder) CloneWithConfig() *TypeStringBuilder {
	return &TypeStringBuilder{
		isForAnnotations: b.isForAnnotations,
		indentation:      b.indentation,
	}
}

// ---------------------------------------------------------------------------
// Helpers that mirror Java's JSType.is*() predicates, implemented via
// Go type assertions so we avoid adding methods to the JSType interface.
// ---------------------------------------------------------------------------

// isObjectType mirrors JSType.isObject(): true for ObjectType,
// TemplatizedType (which wraps an ObjectType), FunctionType (functions
// are objects), EnumType, and EnumElementType.
func isObjectType(t JSType) bool {
	switch t.(type) {
	case *ObjectType, *TemplatizedType, *FunctionType, *EnumType, *EnumElementType:
		return true
	}
	return false
}

// isTemplateType mirrors JSType.isTemplateType().
func isTemplateType(t JSType) bool {
	_, ok := t.(*TemplateType)
	return ok
}

// isRecordType mirrors JSType.isRecordType(): an ObjectType whose
// IsStructural flag is set (corresponds to Java's RecordType / @record).
func isRecordType(t JSType) bool {
	if o, ok := t.(*ObjectType); ok {
		return o.IsStructural
	}
	return false
}

// isFunctionType mirrors JSType.isFunctionType().
func isFunctionType(t JSType) bool {
	_, ok := t.(*FunctionType)
	return ok
}

// isUnionType mirrors JSType.isUnionType().
func isUnionType(t JSType) bool {
	_, ok := t.(*UnionType)
	return ok
}

// isLiteralObject mirrors JSType.isLiteralObject(): an anonymous
// PrototypeObjectType (ObjectType with no Name) that is not structural.
func isLiteralObject(t JSType) bool {
	if o, ok := t.(*ObjectType); ok {
		return o.Name == "" && !o.IsStructural
	}
	return false
}
