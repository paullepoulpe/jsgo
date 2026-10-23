// Port of com/google/javascript/rhino/InputId.java @ closure-compiler v20261015

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
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"sync"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// InputId is an id used uniquely identify a CompilerInput.
type InputId struct {
	id   string
	hash int32
}

type pooledInputIdSlab struct {
	chunk []InputId
}

var inputIdSlabPool = sync.Pool{
	New: func() any { return &pooledInputIdSlab{} },
}

// Java: InputId.serialVersionUID
const InputId_serialVersionUID int64 = 1

// Java: InputId(String)
func NewInputId(id string) *InputId {
	h := jutil.StringHashCode(id)
	if !nodePoolingEnabled.Load() {
		return &InputId{id: id, hash: h}
	}
	s := inputIdSlabPool.Get().(*pooledInputIdSlab)
	if len(s.chunk) == 0 {
		s.chunk = make([]InputId, 1024) // 24 KiB <= 32 KiB
	}
	out := &s.chunk[0]
	s.chunk = s.chunk[1:]
	inputIdSlabPool.Put(s)
	out.id = id
	out.hash = h
	return out
}

// Java: InputId.getIdName()
func (i *InputId) GetIdName() string {
	return i.id
}

// Java: InputId.hashCode()
func (i *InputId) HashCode() int32 {
	return i.hash
}

// Java: InputId.equals(Object)
func (i *InputId) Equals(obj any) bool {
	other, ok := obj.(*InputId)
	if !ok || other == nil {
		return false
	}
	return i == other || (i.hash == other.hash && i.id == other.id)
}

// Java: InputId.toString()
func (i *InputId) String() string {
	return "InputId: " + i.GetIdName()
}
