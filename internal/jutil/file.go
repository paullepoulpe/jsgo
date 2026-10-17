package jutil

import (
	"errors"
	"io/fs"
	"os"
	"strings"
)

// Go-only: language-adapter — java.io.File, as an abstract pathname (the
// methods Closure uses). Paths use the Unix separator.

// File is java.io.File.
type File struct {
	path string
}

// NewFile is new File(String pathname). The pathname is normalized as
// observed through File.getPath() (file_test.go): runs of '/' collapse to
// one and a trailing '/' is dropped, except in the root path "/".
func NewFile(pathname string) *File {
	return &File{path: fileNormalize(pathname)}
}

// fileNormalize drops the empty names of a '/'-separated pathname, keeping
// a leading '/'.
func fileNormalize(pathname string) string {
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

// GetPath is File.getPath().
func (f *File) GetPath() string { return f.path }

// String is File.toString() (the path).
func (f *File) String() string { return f.path }

// ToPath is File.toPath().
func (f *File) ToPath() *Path { return Path_Of(f.path) }

// File_separator is File.separator ("/" on Unix).
const File_separator = "/"

// File_separatorChar is File.separatorChar ('/' on Unix).
const File_separatorChar = '/'

// NewFileWithParent is new File(File parent, String child): child alone
// when parent is null; resolved against the root directory when parent is
// the empty pathname; otherwise parent's path, a separator, then child.
func NewFileWithParent(parent *File, child string) *File {
	if parent == nil {
		return NewFile(child)
	}
	if parent.path == "" {
		return NewFile("/" + fileNormalize(child))
	}
	return NewFile(parent.path + "/" + child)
}

// lastSeparator is the index of the last name-separator character, or -1.
func (f *File) lastSeparator() int {
	return strings.LastIndexByte(f.path, '/')
}

// GetName is File.getName(): the last name in the pathname's name sequence
// ("" for an empty name sequence).
func (f *File) GetName() string {
	return f.path[f.lastSeparator()+1:]
}

// GetParent is File.getParent(): the pathname of the parent (the prefix and
// every name but the last), or nil (null) if the pathname names no parent.
func (f *File) GetParent() *string {
	index := f.lastSeparator()
	if index < 0 {
		return nil
	}
	// The prefix of an absolute pathname is "/".
	if index == 0 {
		if len(f.path) > 1 {
			root := "/"
			return &root
		}
		return nil
	}
	parent := f.path[:index]
	return &parent
}

// GetParentFile is File.getParentFile(): nil (null) when GetParent is.
func (f *File) GetParentFile() *File {
	p := f.GetParent()
	if p == nil {
		return nil
	}
	return NewFile(*p)
}

// IsDirectory is File.isDirectory(): whether the file exists and is a
// directory (following symbolic links).
func (f *File) IsDirectory() bool {
	st, err := os.Stat(f.path)
	return err == nil && st.IsDir()
}

// Mkdirs is File.mkdirs(): creates the directory and any missing parents,
// returning true if and only if the directory was created (false if it
// already existed or could not be created).
func (f *File) Mkdirs() bool {
	if f.Exists() {
		return false
	}
	return os.MkdirAll(f.path, 0o777) == nil
}

// CreateNewFile is File.createNewFile(): atomically creates the named file
// if it does not exist, returning true if it was created and false if it
// already existed; any other failure throws IOException with the C
// library's reason.
func (f *File) CreateNewFile() bool {
	file, err := os.OpenFile(f.path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o666)
	if err == nil {
		file.Close()
		return true
	}
	if errors.Is(err, fs.ErrExist) {
		return false
	}
	reason := FileNotFoundExceptionFromError(err, "").GetMessage()
	panic(NewIOException(strings.TrimSuffix(strings.TrimPrefix(reason, " ("), ")")))
}

// GetAbsolutePath is File.getAbsolutePath(): the pathname resolved against
// the current user directory if it is relative.
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

// Delete is File.delete(): deletes the file or (empty) directory, returning
// whether it was deleted.
func (f *File) Delete() bool {
	return os.Remove(f.path) == nil
}

// RenameTo is File.renameTo(File): renames the file, returning whether it
// succeeded.
func (f *File) RenameTo(dest *File) bool {
	CheckNotNull(dest)
	return os.Rename(f.path, dest.path) == nil
}

// ListFiles is File.listFiles(): the entries of the directory as Files
// (child of this pathname), in the order the directory stream returns them,
// or nil (null) if this pathname is not a readable directory.
func (f *File) ListFiles() []*File {
	d, err := os.Open(f.path)
	if err != nil {
		return nil
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return nil
	}
	files := make([]*File, len(names))
	for i, name := range names {
		files[i] = NewFileWithParent(f, name)
	}
	return files
}
