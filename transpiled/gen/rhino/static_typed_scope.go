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
// Port of com/google/javascript/rhino/jstype/StaticTypedScope.java (generated mechanically by j2go).

package rhino

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java interface: com.google.javascript.rhino.jstype.StaticTypedScope
type StaticTypedScope interface {
	Implements_rhino_StaticTypedScope()
	IsJavaNil_() bool
	StaticScope
	GetTypeOfThis() JSType
	LookupQualifiedName_1_QualifiedName(QualifiedName) JSType
	LookupQualifiedName_1_string(string) JSType
	LookupQualifiedNameNode(Node) JSType
}

func ToIface_StaticTypedScope[T interface {
	StaticTypedScope
	comparable
}](v T) StaticTypedScope {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_StaticTypedScope[T StaticTypedScope](v T) StaticTypedScope {
	var vi StaticTypedScope = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.jstype.StaticTypedScope", "com.google.javascript.rhino.jstype.StaticTypedScope", func(v_ any) bool {
		return j2gort.InstanceOf[StaticTypedScope](v_)
	})
}

// Java default interface method: StaticTypedScope.lookupQualifiedName(com.google.javascript.rhino.QualifiedName qname)
func IfaceDefault_StaticTypedScope_LookupQualifiedName_1_QualifiedName(this_ StaticTypedScope, qname QualifiedName) JSType {
	_ = this_
	_ = qname
	var slot StaticTypedSlot = j2gort.AsIface[StaticTypedSlot](this_.GetSlot(qname.Join_0()))
	_ = slot
	if (slot != nil && !slot.IsJavaNil_()) && !(slot.IsTypeInferred()) {
		return slot.GetType()
	} else if !(qname.IsSimple()) {
		var type_ JSType = this_.LookupQualifiedName_1_QualifiedName(qname.GetOwner())
		_ = type_
		if (type_ != nil && !type_.IsJavaNil_()) && !(type_.GetJSTypeBase().IsUnknownType()) {
			return type_.GetJSTypeBase().FindPropertyType_1_string(qname.GetComponent())
		}
	}
	return nil
}

// Java default interface method: StaticTypedScope.lookupQualifiedName(java.lang.String qname)
func IfaceDefault_StaticTypedScope_LookupQualifiedName_1_string(this_ StaticTypedScope, qname string) JSType {
	_ = this_
	_ = qname
	var slot StaticTypedSlot = j2gort.AsIface[StaticTypedSlot](this_.GetSlot(qname))
	_ = slot
	if (slot != nil && !slot.IsJavaNil_()) && !(slot.IsTypeInferred()) {
		return slot.GetType()
	}
	var lastDot int32 = j2gort.JStrLastIndexOfCp(qname, int32(uint16(46)))
	_ = lastDot
	if lastDot >= int32(0) {
		var type_ JSType = this_.LookupQualifiedName_1_string(j2gort.JStrSub(qname, int32(0), lastDot))
		_ = type_
		if (type_ != nil && !type_.IsJavaNil_()) && !(type_.GetJSTypeBase().IsUnknownType()) {
			return type_.GetJSTypeBase().FindPropertyType_1_string(j2gort.JStrSubFrom(qname, (lastDot + int32(1))))
		}
	}
	return nil
}

// Java default interface method: StaticTypedScope.lookupQualifiedNameNode(com.google.javascript.rhino.Node qnameNode)
func IfaceDefault_StaticTypedScope_LookupQualifiedNameNode(this_ StaticTypedScope, qnameNode Node) JSType {
	_ = this_
	_ = qnameNode
	if qnameNode.IsName() {
		var name string = qnameNode.GetString()
		_ = name
		if j2gort.CheckNotNullString(name) == "" {
			return nil
		}
		var slot StaticTypedSlot = j2gort.AsIface[StaticTypedSlot](this_.GetSlot(name))
		_ = slot
		if (slot != nil && !slot.IsJavaNil_()) && !(slot.IsTypeInferred()) {
			return slot.GetType()
		} else {
			return nil
		}
	}
	var qname string = qnameNode.GetQualifiedName()
	_ = qname
	if qname != j2gort.NullString {
		return this_.LookupQualifiedName_1_string(qname)
	} else {
		return nil
	}
}
