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
// Port of com/google/javascript/rhino/StaticScope.java (generated mechanically by j2go).

package rhino

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java interface: com.google.javascript.rhino.StaticScope
type StaticScope interface {
	Implements_rhino_StaticScope()
	IsJavaNil_() bool
	GetRootNode_0() Node
	GetParentScope() any
	GetSlot(string) any
	GetOwnSlot(string) any
	GetTopmostScopeOfEventualDeclaration(string) StaticScope
}
func ToIface_StaticScope[T interface {
	StaticScope
	comparable
}](v T) StaticScope {
	var z T
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.StaticScope", "com.google.javascript.rhino.StaticScope", func(v_ any) bool {
		return j2gort.InstanceOf[StaticScope](v_)
	})
}
// Java default interface method: StaticScope.getTopmostScopeOfEventualDeclaration(java.lang.String name)
func IfaceDefault_StaticScope_GetTopmostScopeOfEventualDeclaration(this_ StaticScope, name string) StaticScope {
	_ = this_
	_ = name
	var slot StaticSlot = j2gort.AsIface[StaticSlot](this_.GetOwnSlot(name))
	_ = slot
	if ((slot != nil && !slot.IsJavaNil_())) {
		return j2gort.AsIface[StaticScope](slot.GetScope())
	}
	if !(j2gort.IfaceEquals[StaticScope](j2gort.AsIface[StaticScope](this_.GetParentScope()), nil)) {
		return j2gort.AsIface[StaticScope](this_.GetParentScope()).GetTopmostScopeOfEventualDeclaration(name)
	} else {
		return nil
	}
}

