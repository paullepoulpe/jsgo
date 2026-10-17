// Port of com/google/javascript/rhino/jstype/TypeStringBuilder.java @ closure-compiler v20260415

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Bob Jervis
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// TypeStringBuilder visits types to assemble an associated stringification.
//
// Instances are single use. They irreversibly accumulate state required
// during traversal.
type TypeStringBuilder struct {
	builder          strings.Builder
	isForAnnotations bool
	indentation      string
}

// Java: TypeStringBuilder(boolean)
func NewTypeStringBuilder(isForAnnotations bool) *TypeStringBuilder {
	return &TypeStringBuilder{isForAnnotations: isForAnnotations}
}

// Java: TypeStringBuilder.isForAnnotations()
func (b *TypeStringBuilder) IsForAnnotations() bool {
	return b.isForAnnotations
}

// Java: TypeStringBuilder.append(String)
func (b *TypeStringBuilder) Append(x string) *TypeStringBuilder {
	b.builder.WriteString(x)
	return b
}

// Java: TypeStringBuilder.append(JSType)
func (b *TypeStringBuilder) AppendType(x JSType) *TypeStringBuilder {
	x.AppendTo(b)
	return b
}

// AppendAll appends the elements, which are JSTypes or Strings, separated by
// separator.
//
// Java: TypeStringBuilder.appendAll(Iterable<?>, String)
func (b *TypeStringBuilder) AppendAll(elements []any, separator string) *TypeStringBuilder {
	separate := false
	for _, e := range elements {
		if separate {
			b.Append(separator)
		} else {
			separate = true
		}

		if jSType, ok := e.(JSType); ok {
			b.AppendType(jSType)
		} else if s, ok := e.(string); ok {
			b.Append(s)
		} else {
			panic(jutil.NewClassCastException("java.lang.Object", "java.lang.String"))
		}
	}

	return b
}

// Java: TypeStringBuilder.appendNonNull(JSType)
func (b *TypeStringBuilder) AppendNonNull(type_ JSType) *TypeStringBuilder {
	if b.isForAnnotations &&
		type_.IsObject() &&
		!type_.IsUnknownType() &&
		!type_.IsTemplateType() &&
		!type_.IsRecordType() &&
		!type_.IsFunctionType() &&
		!type_.IsUnionType() &&
		!type_.IsLiteralObject() {
		b.Append("!")
	}
	return b.AppendType(type_)
}

// Java: TypeStringBuilder.breakLineAndIndent()
func (b *TypeStringBuilder) BreakLineAndIndent() *TypeStringBuilder {
	b.builder.WriteString("\n")
	b.builder.WriteString(b.indentation)
	return b
}

// Java: TypeStringBuilder.indent(Runnable)
func (b *TypeStringBuilder) Indent(cb func()) *TypeStringBuilder {
	lastIndent := b.indentation
	b.indentation = lastIndent + "  "
	cb()
	b.indentation = lastIndent
	return b
}

// Java: TypeStringBuilder.build()
func (b *TypeStringBuilder) Build() string {
	return b.builder.String()
}

// Java: TypeStringBuilder.cloneWithConfig()
func (b *TypeStringBuilder) CloneWithConfig() *TypeStringBuilder {
	clone := NewTypeStringBuilder(b.isForAnnotations)
	clone.indentation = b.indentation
	return clone
}
