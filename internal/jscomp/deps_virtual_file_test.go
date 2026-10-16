// Port of test/com/google/javascript/jscomp/deps/VirtualFileTest.java @ closure-compiler v20260415

package jscomp

import "testing"

func TestVirtualFile_Equals(t *testing.T) {
	vf := NewVirtualFile("name", "code")
	same := NewVirtualFile("name", "code")
	other := NewVirtualFile("otherName", "otherCode")
	jsaTrue(t, vf.Equals(same))
	jsaTrue(t, same.Equals(vf))
	jsaFalse(t, vf.Equals(other))
	jsaFalse(t, other.Equals(vf))
}

func TestVirtualFile_HashCode(t *testing.T) {
	vf := NewVirtualFile("name", "code")
	same := NewVirtualFile("name", "code")
	other := NewVirtualFile("otherName", "otherCode")
	jsaEq(t, same.HashCode(), vf.HashCode())
	jsaTrue(t, vf.HashCode() != other.HashCode())
}
