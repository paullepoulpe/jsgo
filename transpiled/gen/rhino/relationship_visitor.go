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
 *   John Lenz
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */
// Port of com/google/javascript/rhino/jstype/RelationshipVisitor.java (generated mechanically by j2go).

package rhino

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java interface: com.google.javascript.rhino.jstype.RelationshipVisitor
type RelationshipVisitor interface {
	Implements_rhino_RelationshipVisitor()
	IsJavaNil_() bool
	CaseUnknownType(JSType, JSType) any
	CaseNoType(JSType) any
	CaseNoObjectType(JSType) any
	CaseAllType(JSType) any
	CaseValueType(ValueType, JSType) any
	CaseObjectType(ObjectType, JSType) any
	CaseFunctionType(*FunctionType, JSType) any
	CaseUnionType(*UnionType, JSType) any
	CaseTemplatizedType(*TemplatizedType, JSType) any
	CaseTemplateType(*TemplateType, JSType) any
	CaseEnumElementType(*EnumElementType, JSType) any
}
func ToIface_RelationshipVisitor[T interface {
	RelationshipVisitor
	comparable
}](v T) RelationshipVisitor {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.RelationshipVisitor", "com.google.javascript.rhino.jstype.RelationshipVisitor", func(v_ any) bool {
		return j2gort.InstanceOf[RelationshipVisitor](v_)
	})
}
