package javaio

import (
	"os"
	"sort"
	"strings"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

const (
	File_separator     = "/"
	File_separatorChar = '/'
	File_pathSeparator = ":"
)

// File models java.io.File on a Unix filesystem.
type File struct {
	path string
}

func NewFile(pathname string) *File {
	return &File{path: normalizeFilePathname(pathname)}
}

func NewFileWithParent(parent *File, child string) *File {
	if parent == nil {
		return NewFile(child)
	}
	if parent.path == "" {
		return NewFile("/" + normalizeFilePathname(child))
	}
	return NewFile(parent.path + "/" + child)
}

func NewFileParentString(parent, child string) *File {
	if parent == "" {
		return NewFile(child)
	}
	return NewFileWithParent(NewFile(parent), child)
}

func normalizeFilePathname(pathname string) string {
	if !strings.Contains(pathname, "//") && (len(pathname) <= 1 || !strings.HasSuffix(pathname, "/")) {
		return pathname
	}
	var names []string
	for _, name := range strings.Split(pathname, "/") {
		if name != "" {
			names = append(names, name)
		}
	}
	s := strings.Join(names, "/")
	if strings.HasPrefix(pathname, "/") {
		s = "/" + s
	}
	return s
}

func (f *File) GetPath() string { return f.path }
func (f *File) String() string  { return f.path }

func (f *File) GetName() string {
	idx := strings.LastIndexByte(f.path, '/')
	return f.path[idx+1:]
}

func (f *File) Parent() (string, bool) {
	idx := strings.LastIndexByte(f.path, '/')
	if idx < 0 {
		return "", false
	}
	if idx == 0 {
		if len(f.path) > 1 {
			return "/", true
		}
		return "", false
	}
	return f.path[:idx], true
}

func (f *File) GetParent() string {
	p, _ := f.Parent()
	return p
}

func (f *File) GetParentFile() *File {
	p, ok := f.Parent()
	if !ok {
		return nil
	}
	return NewFile(p)
}

func (f *File) IsAbsolute() bool {
	return strings.HasPrefix(f.path, "/")
}

func (f *File) GetAbsolutePath() string {
	if strings.HasPrefix(f.path, "/") {
		return f.path
	}
	wd, err := os.Getwd()
	if err != nil {
		return f.path
	}
	if f.path == "" {
		return wd
	}
	return NewFile(wd + "/" + f.path).path
}

func (f *File) GetAbsoluteFile() *File {
	return NewFile(f.GetAbsolutePath())
}

func (f *File) Exists() bool {
	if f.path == "" {
		return false
	}
	_, err := os.Stat(f.path)
	return err == nil
}

func (f *File) IsDirectory() bool {
	if f.path == "" {
		return false
	}
	st, err := os.Stat(f.path)
	return err == nil && st.IsDir()
}

func (f *File) IsFile() bool {
	if f.path == "" {
		return false
	}
	st, err := os.Stat(f.path)
	return err == nil && st.Mode().IsRegular()
}

func (f *File) Length() int64 {
	if f.path == "" {
		return 0
	}
	st, err := os.Stat(f.path)
	if err != nil {
		return 0
	}
	return st.Size()
}

func (f *File) Delete() bool {
	if f.path == "" {
		return false
	}
	return os.Remove(f.path) == nil
}

func (f *File) Mkdir() bool {
	if f.path == "" {
		return false
	}
	return os.Mkdir(f.path, 0o755) == nil
}

func (f *File) Mkdirs() bool {
	if f.path == "" || f.Exists() {
		return false
	}
	return os.MkdirAll(f.path, 0o755) == nil
}

func (f *File) RenameTo(dest *File) bool {
	if dest == nil || f.path == "" || dest.path == "" {
		return false
	}
	return os.Rename(f.path, dest.path) == nil
}

// ListFiles returns directory entries sorted deterministically by name, or nil if not a directory.
func (f *File) ListFiles() []*File {
	if f.path == "" {
		return nil
	}
	entries, err := os.ReadDir(f.path)
	if err != nil {
		return nil
	}
	names := make([]string, len(entries))
	for i, e := range entries {
		names[i] = e.Name()
	}
	sort.Strings(names)
	out := make([]*File, len(names))
	for i, name := range names {
		out[i] = NewFileWithParent(f, name)
	}
	return out
}

func (f *File) CreateNewFile() (bool, error) {
	if f.path == "" {
		return false, &IOException{}
	}
	fd, err := os.OpenFile(f.path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o666)
	if err != nil {
		if os.IsExist(err) {
			return false, nil
		}
		return false, err
	}
	_ = fd.Close()
	return true, nil
}

func (f *File) CompareTo(other *File) int {
	return javalang.StringCompareTo(f.path, other.path)
}

func (f *File) Equals(other any) bool {
	o, ok := other.(*File)
	return ok && o != nil && f.path == o.path
}

func (f *File) HashCode() int32 {
	// UnixFileSystem.hashCode(File): path.hashCode() ^ 1234321
	return javalang.StringHashCode(f.path) ^ 1234321
}
