// Port of com/google/javascript/jscomp/deps/VirtualFile.java @ closure-compiler v20260415

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

package jscomp

import "github.com/paullepoulpe/jsgo/internal/jutil"

// VirtualFile is used to represent a "virtual" file.
type VirtualFile struct {
	name string
	code string
}

var _ Deps_SourceFile = (*VirtualFile)(nil)

// Java: VirtualFile(String, String)
func NewVirtualFile(name string, code string) *VirtualFile {
	return &VirtualFile{name: name, code: code}
}

// Java: VirtualFile.getName()
func (f *VirtualFile) GetName() string { return f.name }

// Java: VirtualFile.getContent()
func (f *VirtualFile) GetContent() string { return f.code }

// Java: VirtualFile.wasModified()
func (f *VirtualFile) WasModified() bool { return false }

// Java: VirtualFile.hashCode() (Objects.hash(name, code))
func (f *VirtualFile) HashCode() int32 {
	h := int32(1)
	h = 31*h + jutil.StringHashCode(f.name)
	h = 31*h + jutil.StringHashCode(f.code)
	return h
}

// Java: VirtualFile.equals(Object)
func (f *VirtualFile) Equals(o any) bool {
	if vf, ok := o.(*VirtualFile); ok && vf != nil {
		return f.name == vf.name && f.code == vf.code
	}
	return false
}
