/*
 * Copyright 2008 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
// Port of com/google/javascript/jscomp/deps/VirtualFile.java (generated mechanically by j2go).

package jscomp

import (
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.deps.VirtualFile
type VirtualFile struct {
	name string `gson:"name"`
	code string `gson:"code"`
}

func (this_ *VirtualFile) BindSelf_(s any) {
	_ = s
}
func (this_ *VirtualFile) IsJavaNil_() bool { return this_ == nil }
func (*VirtualFile) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.deps.VirtualFile", "VirtualFile"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.deps.VirtualFile", "com.google.javascript.jscomp.deps.VirtualFile", func(v_ any) bool {
		return j2gort.InstanceOf[*VirtualFile](v_)
	})
}
func (*VirtualFile) Implements_jscomp_Deps_SourceFile() {}

func (this_ *VirtualFile) initFields_VirtualFile() {
	_ = this_
	this_.name = j2gort.NullString
	this_.code = j2gort.NullString
}

// Java constructor: VirtualFile(java.lang.String name, java.lang.String code)
func NewVirtualFile(name string, code string) *VirtualFile {
	this_ := &VirtualFile{}
	InitVirtualFile(this_, name, code)
	return this_
}

func InitVirtualFile(this_ *VirtualFile, name string, code string) {
	_ = name
	_ = code
	this_.initFields_VirtualFile()
	this_.name = name
	this_.code = code
}

// Java method: VirtualFile.getName()
func (this_ *VirtualFile) GetName() string {
	return this_.name
}

// Java method: VirtualFile.getContent()
func (this_ *VirtualFile) GetContent() string {
	return this_.code
}

// Java method: VirtualFile.wasModified()
func (this_ *VirtualFile) WasModified() bool {
	return false
}

// Java method: VirtualFile.hashCode()
func (this_ *VirtualFile) HashCode() int32 {
	{
		var oh_1 int32 = 1
		{
			hs_2 := this_.name
			oh_1 = 31 * oh_1
			if hs_2 != j2gort.NullString {
				oh_1 += javalang.StringHashCode(hs_2)
			}
		}
		{
			hs_3 := this_.code
			oh_1 = 31 * oh_1
			if hs_3 != j2gort.NullString {
				oh_1 += javalang.StringHashCode(hs_3)
			}
		}
		return oh_1
	}
}

// Java method: VirtualFile.equals(java.lang.Object o)
func (this_ *VirtualFile) Equals(arg_4 any) bool {
	o := arg_4
	_ = o
	var vf *VirtualFile
	_ = vf
	var ok_5 bool
	vf, ok_5 = j2gort.TryCast[*VirtualFile](o)
	if ok_5 {
		return ((this_.name == vf.name) && (this_.code == vf.code))
	}
	return false
}
