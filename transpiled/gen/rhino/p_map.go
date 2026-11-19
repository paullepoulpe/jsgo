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
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */
// Port of com/google/javascript/rhino/PMap.java (generated mechanically by j2go).

package rhino

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
)

// Java interface: com.google.javascript.rhino.PMap
type PMap[K comparable, V comparable] interface {
	Implements_rhino_PMap()
	IsJavaNil_() bool
	IsEmpty() bool
	Values() j2gort.Ext_java_lang_Iterable
	Keys() j2gort.Ext_java_lang_Iterable
	Get_1(K) V
	Plus_2(K, V) PMap[K, V]
	Minus(K) PMap[K, V]
	Reconcile(PMap[K, V], PMap_Reconciler[K, V]) PMap[K, V]
	Equivalent(PMap[K, V], j2gort.Ext_java_util_function_BiPredicate) bool
}
func ToIface_PMap[K comparable, V comparable, T_iface_ interface {
	PMap[K, V]
	comparable
}](v T_iface_) PMap[K, V] {
	var z T_iface_
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.PMap", "com.google.javascript.rhino.PMap", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			Implements_rhino_PMap()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
// Java interface: com.google.javascript.rhino.PMap.Reconciler
type PMap_Reconciler[K comparable, V comparable] interface {
	Implements_rhino_PMap_Reconciler()
	IsJavaNil_() bool
	Merge(K, V, V) V
}
func ToIface_PMap_Reconciler[K comparable, V comparable, T_iface_ interface {
	PMap_Reconciler[K, V]
	comparable
}](v T_iface_) PMap_Reconciler[K, V] {
	var z T_iface_
	if v == z {
		return nil
	}
	return v
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.rhino.PMap.Reconciler", "com.google.javascript.rhino.PMap$Reconciler", func(v_ any) bool {
		m_, ok_ := v_.(interface {
			Implements_rhino_PMap_Reconciler()
			IsJavaNil_() bool
		})
		return ok_ && !m_.IsJavaNil_()
	})
}
