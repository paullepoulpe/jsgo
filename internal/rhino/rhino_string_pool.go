// Port of com/google/javascript/rhino/RhinoStringPool.java @ closure-compiler v20261028

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

package rhino

import (
	"iter"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// rhinoStringID is a 32-bit pointer-free handle to an interned string (see
// package strpool): 0 is null / unset, 1 is "".
type rhinoStringID = strpool.ID

func rhinoStringPoolIntern(s string) rhinoStringID { return strpool.Intern(s) }
func rhinoStringFromID(id rhinoStringID) string    { return strpool.String(id) }

func rhinoStringIDFromPtr(s *string) rhinoStringID {
	if s == nil {
		return 0
	}
	return strpool.Intern(*s)
}

// Java: RhinoStringPool.addOrGet(String)
func RhinoStringPool_AddOrGet(s string) string {
	return strpool.String(strpool.Intern(s))
}

// RhinoStringPool_EmptyID is the canonical 32-bit ID of the empty string "".
const RhinoStringPool_EmptyID uint32 = strpool.EmptyID

// RhinoStringPool_InternID interns s and returns its 32-bit pointer-free ID.
func RhinoStringPool_InternID(s string) uint32 {
	return uint32(strpool.Intern(s))
}

// RhinoStringPool_InternBytesID interns b and returns its 32-bit pointer-free ID
// without allocating a string when b is already interned.
func RhinoStringPool_InternBytesID(b []byte) uint32 {
	return uint32(strpool.InternBytes(b))
}

// RhinoStringPool_LookupID returns the 32-bit ID of s if already interned, or 0 if not.
func RhinoStringPool_LookupID(s string) uint32 {
	return uint32(strpool.Lookup(s))
}

// RhinoStringPool_StringFromID returns the interned string for a 32-bit ID.
func RhinoStringPool_StringFromID(id uint32) string {
	return strpool.String(strpool.ID(id))
}

// RhinoStringPool_LazyInternedStringList is a list of interned strings backed
// directly by pre-interned 32-bit rhinoStringIDs (noscan, zero pointers).
type RhinoStringPool_LazyInternedStringList struct {
	ids []uint32
}

// NewRhinoStringPool_LazyInternedStringListFromIDs creates a LazyInternedStringList
// backed directly by pre-interned 32-bit rhinoStringIDs (noscan, zero pointers).
func NewRhinoStringPool_LazyInternedStringListFromIDs(ids []uint32) *RhinoStringPool_LazyInternedStringList {
	return &RhinoStringPool_LazyInternedStringList{ids: ids}
}

//go:noinline
func (l *RhinoStringPool_LazyInternedStringList) outOfBounds(offset int) rhinoStringID {
	panic(jutil.NewArrayIndexOutOfBoundsExceptionForLength(offset, len(l.ids)))
}

// GetID returns the interned rhinoStringID at offset.
func (l *RhinoStringPool_LazyInternedStringList) GetID(offset int) rhinoStringID {
	if uint(offset) < uint(len(l.ids)) {
		return rhinoStringID(l.ids[offset])
	}
	return l.outOfBounds(offset)
}

// Get returns the interned string at offset.
//
// Java: RhinoStringPool.LazyInternedStringList.get(int)
func (l *RhinoStringPool_LazyInternedStringList) Get(offset int) string {
	return strpool.String(l.GetID(offset))
}

// Java: RhinoStringPool.LazyInternedStringList.stream()
func (l *RhinoStringPool_LazyInternedStringList) Stream() iter.Seq[string] {
	return func(yield func(string) bool) {
		for _, id := range l.ids {
			if !yield(strpool.String(strpool.ID(id))) {
				return
			}
		}
	}
}

// Slice returns the underlying string slice.
func (l *RhinoStringPool_LazyInternedStringList) Slice() []string {
	out := make([]string, len(l.ids))
	for i, id := range l.ids {
		out[i] = strpool.String(strpool.ID(id))
	}
	return out
}

// Len returns the number of strings in the list.
func (l *RhinoStringPool_LazyInternedStringList) Len() int {
	return len(l.ids)
}
