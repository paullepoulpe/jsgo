// Port of com/google/javascript/jscomp/deps/VirtualFile.java @ closure-compiler v20260415

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
