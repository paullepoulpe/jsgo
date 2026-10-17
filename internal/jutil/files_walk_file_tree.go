package jutil

import (
	"os"
	"syscall"
)

// Go-only: language-adapter — java.nio.file.Files.walkFileTree with a
// FileVisitor, written from its Javadoc: a depth-first walk in which each
// directory's entries are visited in the order the directory stream returns
// them (the order of readdir, which is unsorted), the start file's
// attributes are read first (a failure is reported to visitFileFailed), a
// directory that cannot be opened is reported to visitFileFailed, and with
// FOLLOW_LINKS symbolic links are followed: a link whose target cannot be
// read is visited as a file, and a link back to an ancestor directory is a
// FileSystemLoopException reported to visitFileFailed. Exceptions passed
// to visitFileFailed are IOExceptions whose message is the file's path.

// FileVisitResult is java.nio.file.FileVisitResult.
type FileVisitResult int32

const (
	FileVisitResult_CONTINUE FileVisitResult = iota
	FileVisitResult_TERMINATE
	FileVisitResult_SKIP_SUBTREE
	FileVisitResult_SKIP_SIBLINGS
)

// FileVisitor is java.nio.file.FileVisitor<Path>, without the file
// attributes (no Closure visitor reads them).
type FileVisitor interface {
	PreVisitDirectory(dir *Path) FileVisitResult
	VisitFile(file *Path) FileVisitResult
	VisitFileFailed(file *Path, exc IIOException) FileVisitResult
	PostVisitDirectory(dir *Path, exc IIOException) FileVisitResult
}

// SimpleFileVisitor is java.nio.file.SimpleFileVisitor: CONTINUE
// everywhere, and visitFileFailed / postVisitDirectory rethrow the
// exception. Visitors embed it and override what they need.
type SimpleFileVisitor struct{}

// PreVisitDirectory is SimpleFileVisitor.preVisitDirectory.
func (SimpleFileVisitor) PreVisitDirectory(dir *Path) FileVisitResult {
	return FileVisitResult_CONTINUE
}

// VisitFile is SimpleFileVisitor.visitFile.
func (SimpleFileVisitor) VisitFile(file *Path) FileVisitResult { return FileVisitResult_CONTINUE }

// VisitFileFailed is SimpleFileVisitor.visitFileFailed: it rethrows.
func (SimpleFileVisitor) VisitFileFailed(file *Path, exc IIOException) FileVisitResult {
	panic(exc)
}

// PostVisitDirectory is SimpleFileVisitor.postVisitDirectory: it rethrows a
// non-null exception.
func (SimpleFileVisitor) PostVisitDirectory(dir *Path, exc IIOException) FileVisitResult {
	if exc != nil {
		panic(exc)
	}
	return FileVisitResult_CONTINUE
}

type fileKey struct{ dev, ino uint64 }

// Files_WalkFileTree is Files.walkFileTree(start, options, maxDepth,
// visitor), followLinks standing for EnumSet.of(FOLLOW_LINKS).
func Files_WalkFileTree(start *Path, followLinks bool, maxDepth int, visitor FileVisitor) *Path {
	w := &treeWalk{followLinks: followLinks, maxDepth: maxDepth, visitor: visitor}
	w.visit(start, 0, nil)
	return start
}

type treeWalk struct {
	followLinks bool
	maxDepth    int
	visitor     FileVisitor
}

func (w *treeWalk) stat(p *Path) (os.FileInfo, error) {
	if !w.followLinks {
		return os.Lstat(p.String())
	}
	st, err := os.Stat(p.String())
	if err != nil {
		// The target cannot be read: use the link's own attributes.
		if lst, lerr := os.Lstat(p.String()); lerr == nil {
			return lst, nil
		}
	}
	return st, err
}

// visit visits file; it returns false when the walk terminates.
func (w *treeWalk) visit(file *Path, depth int, ancestors []fileKey) bool {
	st, err := w.stat(file)
	if err != nil {
		return w.visitFileFailed(file)
	}
	if !st.IsDir() || depth >= w.maxDepth {
		return w.visitor.VisitFile(file) != FileVisitResult_TERMINATE
	}
	var key fileKey
	if sys, ok := st.Sys().(*syscall.Stat_t); ok {
		key = fileKey{uint64(sys.Dev), uint64(sys.Ino)}
		if w.followLinks {
			for _, a := range ancestors {
				if a == key {
					return w.visitFileFailedWith(file, NewIOException(file.String()))
				}
			}
		}
	}
	d, err := os.Open(file.String())
	if err != nil {
		return w.visitFileFailed(file)
	}
	names, err := d.Readdirnames(-1)
	d.Close()
	if err != nil {
		return w.visitFileFailed(file)
	}
	switch w.visitor.PreVisitDirectory(file) {
	case FileVisitResult_TERMINATE:
		return false
	case FileVisitResult_SKIP_SUBTREE, FileVisitResult_SKIP_SIBLINGS:
		return true
	}
	ancestors = append(ancestors, key)
	for _, name := range names {
		if !w.visit(file.ResolveString(name), depth+1, ancestors) {
			return false
		}
	}
	return w.visitor.PostVisitDirectory(file, nil) != FileVisitResult_TERMINATE
}

func (w *treeWalk) visitFileFailed(file *Path) bool {
	return w.visitFileFailedWith(file, NewIOException(file.String()))
}

func (w *treeWalk) visitFileFailedWith(file *Path, exc IIOException) bool {
	return w.visitor.VisitFileFailed(file, exc) != FileVisitResult_TERMINATE
}
