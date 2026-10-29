package jutil

import (
	"os"
	"strings"
)

// Go-only: language-adapter — java.nio.file.Path for the default file
// system on Linux, as Path.of builds it (see the Javadoc of
// java.nio.file.Path and FileSystem.getPath). A path is a byte string (the
// UTF-8 form of the path string, which is what the JDK uses with
// sun.jnu.encoding=UTF-8): an optional root "/" followed by name elements
// separated by "/". The empty path has one name element, the empty name.
//
// Observed behaviour (black-box, JDK 21 on Linux, UTF-8 locale):
//   - parsing drops repeated and trailing slashes ("a//b/" is "a/b", "//" is
//     "/"); it throws InvalidPathException("Nul character not allowed") for
//     a NUL in the (joined) string, and then InvalidPathException("Malformed
//     input or input contains unmappable characters") for a lone surrogate
//     in the path without the redundant slashes; neither has an index;
//   - hashCode is 31*h + b over the unsigned bytes, and compareTo returns the
//     difference of the first unequal unsigned bytes, or of the lengths;
//   - getName and subpath with bad indices throw IllegalArgumentException
//     without a message;
//   - normalize drops "." elements and each name followed by "..", keeps a
//     leading ".." of a relative path and drops one directly after the
//     root ("/.." is "/"); the result may be the empty path;
//   - relativize: the empty path relativized against p gives p unchanged;
//     otherwise both paths are normalized first, the common leading
//     elements dropped, and the rest of this path replaced by ".."
//     elements. When the rest of this path still holds "..", the result
//     cannot be computed: IllegalArgumentException("Unable to compute
//     relative  path from <this> to <other>") (two spaces). A relative and
//     an absolute path give IllegalArgumentException("'other' is different
//     type of Path").

// Path is java.nio.file.Path (a Unix path of the default file system).
type Path struct {
	s    string
	hash int32
	// elemStarts holds the byte offset where each name element begins; nil
	// until first needed.
	elemStarts []int
}

// Path_Of is Path.of(first, more...): when more is not empty, the non-empty
// strings among first and more are joined with "/".
func Path_Of(first string, more ...string) *Path {
	input := first
	if len(more) > 0 {
		var b strings.Builder
		b.Grow(len(first) + len(more)*8)
		if first != "" {
			b.WriteString(first)
		}
		for _, s := range more {
			if s == "" {
				continue
			}
			if b.Len() > 0 {
				b.WriteByte('/')
			}
			b.WriteString(s)
		}
		input = b.String()
	}
	return parsePath(input)
}

// ParsePathString validates and normalizes input identically to Path_Of(input).String()
// without allocating a *Path.
func ParsePathString(input string) string {
	slow := len(input) > 1 && input[len(input)-1] == '/'
	for i := 0; i < len(input); i++ {
		c := input[i]
		if c == 0 {
			panic(newInvalidPathException(input, "Nul character not allowed"))
		}
		if c == '/' && i > 0 && input[i-1] == '/' {
			slow = true
		} else if c == 0xED && i+1 < len(input) && input[i+1] >= 0xA0 && input[i+1] <= 0xBF {
			slow = true
		}
	}
	if !slow {
		return input
	}
	b := make([]byte, 0, len(input))
	for i := 0; i < len(input); i++ {
		c := input[i]
		if c == '/' && (len(b) > 0 && b[len(b)-1] == '/') {
			continue
		}
		b = append(b, c)
	}
	if len(b) > 1 && b[len(b)-1] == '/' {
		b = b[:len(b)-1]
	}
	for i := 0; i+1 < len(b); i++ {
		if b[i] == 0xED && b[i+1] >= 0xA0 && b[i+1] <= 0xBF {
			panic(newInvalidPathException(string(b), "Malformed input or input contains unmappable characters"))
		}
	}
	return string(b)
}

// parsePath checks a path string and removes its redundant slashes.
func parsePath(input string) *Path {
	return newPathFromString(ParsePathString(input))
}

func newPathFromString(s string) *Path {
	switch s {
	case "":
		return emptyPath
	case "/":
		return rootPath
	default:
		return &Path{s: s}
	}
}

var (
	emptyPath = &Path{s: "", elemStarts: []int{0}}
	rootPath  = &Path{s: "/", hash: '/', elemStarts: []int{}}
)

func (p *Path) isEmpty() bool { return len(p.s) == 0 }

// elementOffsets returns the start offset of every name element.
func (p *Path) elementOffsets() []int {
	if p.elemStarts != nil {
		return p.elemStarts
	}
	if p.isEmpty() {
		starts := []int{0}
		p.elemStarts = starts
		return starts
	}
	s := p.s
	n := 0
	for i := 0; i < len(s); i++ {
		if s[i] != '/' && (i == 0 || s[i-1] == '/') {
			n++
		}
	}
	if n == 0 {
		starts := []int{}
		p.elemStarts = starts
		return starts
	}
	starts := make([]int, 0, n)
	for i := 0; i < len(s); i++ {
		if s[i] != '/' && (i == 0 || s[i-1] == '/') {
			starts = append(starts, i)
		}
	}
	p.elemStarts = starts
	return starts
}

// element returns name element i.
func (p *Path) element(i int) string {
	starts := p.elementOffsets()
	end := len(p.s)
	if i+1 < len(starts) {
		end = starts[i+1] - 1
	}
	return p.s[starts[i]:end]
}

func (p *Path) elements() []string {
	n := len(p.elementOffsets())
	out := make([]string, n)
	for i := range out {
		out[i] = p.element(i)
	}
	return out
}

// GetRoot is Path.getRoot(): "/" for an absolute path, else null.
func (p *Path) GetRoot() *Path {
	if p.IsAbsolute() {
		return rootPath
	}
	return nil
}

// GetFileName is Path.getFileName(): the last name element, null for the
// root, and the empty path itself for the empty path.
func (p *Path) GetFileName() *Path {
	if p.isEmpty() {
		return p
	}
	if p.s == "/" {
		return nil
	}
	idx := strings.LastIndexByte(p.s, '/')
	if idx < 0 {
		return p
	}
	return newPathFromString(p.s[idx+1:])
}

// GetParent is Path.getParent(): the path without its last name element,
// or null when there is none (the root, a one-element relative path).
func (p *Path) GetParent() *Path {
	if len(p.s) == 0 || p.s == "/" {
		return nil
	}
	idx := strings.LastIndexByte(p.s, '/')
	if idx < 0 {
		return nil
	}
	if idx == 0 {
		return rootPath
	}
	return newPathFromString(p.s[:idx])
}

// GetNameCount is Path.getNameCount().
func (p *Path) GetNameCount() int { return len(p.elementOffsets()) }

// GetName is Path.getName(int).
func (p *Path) GetName(index int) *Path {
	if index < 0 || index >= p.GetNameCount() {
		panic(NewIllegalArgumentException())
	}
	return newPathFromString(p.element(index))
}

// subpath is Path.subpath(int, int): the elements beginIndex (inclusive)
// to endIndex (exclusive), as a relative path.
func (p *Path) subpath(beginIndex, endIndex int) *Path {
	n := p.GetNameCount()
	if beginIndex < 0 || beginIndex >= n || endIndex > n || beginIndex >= endIndex {
		panic(NewIllegalArgumentException())
	}
	starts := p.elementOffsets()
	end := len(p.s)
	if endIndex < n {
		end = starts[endIndex] - 1
	}
	return newPathFromString(p.s[starts[beginIndex]:end])
}

// IsAbsolute is Path.isAbsolute(): whether the path has the root.
func (p *Path) IsAbsolute() bool { return len(p.s) > 0 && p.s[0] == '/' }

// Resolve is Path.resolve(Path): other if it is absolute, this path if other
// is empty, else other appended to this path.
func (p *Path) Resolve(obj *Path) *Path {
	if obj == nil {
		panic(NewNullPointerException())
	}
	switch {
	case obj.IsAbsolute():
		return obj
	case obj.isEmpty():
		return p
	case p.isEmpty():
		return obj
	}
	if p.s[len(p.s)-1] == '/' {
		return newPathFromString(p.s + obj.s)
	}
	return newPathFromString(p.s + "/" + obj.s)
}

// ResolveString is Path.resolve(String).
func (p *Path) ResolveString(other string) *Path { return p.Resolve(parsePath(other)) }

// ResolveSiblingString is Path.resolveSibling(String).
func (p *Path) ResolveSiblingString(other string) *Path {
	parsed := parsePath(other)
	parent := p.GetParent()
	if parent == nil {
		return parsed
	}
	return parent.Resolve(parsed)
}

// Relativize is Path.relativize(Path) (see the observed behaviour above).
func (p *Path) Relativize(obj *Path) *Path {
	if obj.IsAbsolute() != p.IsAbsolute() {
		panic(NewIllegalArgumentExceptionMsg("'other' is different type of Path"))
	}
	if p.Equals(obj) {
		return emptyPath
	}
	if p.isEmpty() {
		return obj
	}
	base, target := p.Normalize(), obj.Normalize()
	baseNames, targetNames := base.names(), target.names()
	common := 0
	for common < len(baseNames) && common < len(targetNames) && baseNames[common] == targetNames[common] {
		common++
	}
	var b strings.Builder
	for _, name := range baseNames[common:] {
		if name == ".." {
			panic(NewIllegalArgumentExceptionMsg("Unable to compute relative  path from " + p.String() + " to " + obj.String()))
		}
		if b.Len() > 0 {
			b.WriteByte('/')
		}
		b.WriteString("..")
	}
	for _, name := range targetNames[common:] {
		if b.Len() > 0 {
			b.WriteByte('/')
		}
		b.WriteString(name)
	}
	if b.Len() == 0 {
		return emptyPath
	}
	return newPathFromString(b.String())
}

// names returns the name elements, none for the empty path.
func (p *Path) names() []string {
	if p.isEmpty() {
		return nil
	}
	return p.elements()
}

func pathHasDotSegment(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == '.' && (i == 0 || s[i-1] == '/') {
			end := i + 1
			if end < len(s) && s[end] == '.' {
				end++
			}
			if end == len(s) || s[end] == '/' {
				return true
			}
		}
	}
	return false
}

// Normalize is Path.normalize() (see the observed behaviour above).
func (p *Path) Normalize() *Path {
	if !pathHasDotSegment(p.s) {
		return p
	}
	var kept []string
	changed := false
	for _, name := range p.names() {
		switch name {
		case ".":
			changed = true
			continue
		case "..":
			if len(kept) > 0 && kept[len(kept)-1] != ".." {
				kept = kept[:len(kept)-1]
				changed = true
				continue
			}
			if p.IsAbsolute() {
				changed = true
				continue
			}
		}
		kept = append(kept, name)
	}
	if !changed {
		return p
	}
	if len(kept) == 0 {
		if p.IsAbsolute() {
			return rootPath
		}
		return emptyPath
	}
	var b strings.Builder
	b.Grow(len(p.s))
	if p.IsAbsolute() {
		b.WriteByte('/')
	}
	for i, name := range kept {
		if i > 0 {
			b.WriteByte('/')
		}
		b.WriteString(name)
	}
	return newPathFromString(b.String())
}

// StartsWith is Path.startsWith(Path): same root and the same leading name
// elements.
func (p *Path) StartsWith(other *Path) bool {
	if p.IsAbsolute() != other.IsAbsolute() {
		return false
	}
	if p.isEmpty() || other.isEmpty() {
		return p.isEmpty() && other.isEmpty()
	}
	if other.s == "/" {
		return true
	}
	if len(other.s) > len(p.s) {
		return false
	}
	if len(other.s) == len(p.s) {
		return p.s == other.s
	}
	return p.s[len(other.s)] == '/' && p.s[:len(other.s)] == other.s
}

// startsWithString is Path.startsWith(String).
func (p *Path) startsWithString(other string) bool { return p.StartsWith(parsePath(other)) }

// EndsWith is Path.endsWith(Path): the same trailing name elements; an
// absolute other must be equal to this path.
func (p *Path) EndsWith(other *Path) bool {
	if other.IsAbsolute() {
		return p.Equals(other)
	}
	if p.isEmpty() || other.isEmpty() {
		return p.isEmpty() && other.isEmpty()
	}
	if len(other.s) > len(p.s) {
		return false
	}
	if len(other.s) == len(p.s) {
		return p.s == other.s
	}
	off := len(p.s) - len(other.s)
	return p.s[off-1] == '/' && p.s[off:] == other.s
}

// endsWithString is Path.endsWith(String).
func (p *Path) endsWithString(other string) bool { return p.EndsWith(parsePath(other)) }

// CompareTo is Path.compareTo(Path): the unsigned bytes, lexicographically.
func (p *Path) CompareTo(other *Path) int {
	a, b := p.s, other.s
	if a == b {
		return 0
	}
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			return int(a[i]) - int(b[i])
		}
	}
	return len(a) - len(b)
}

// Equals is Path.equals(Object): another Path with the same bytes.
func (p *Path) Equals(ob any) bool {
	o, ok := ob.(*Path)
	return ok && o != nil && p.s == o.s
}

// HashCode is Path.hashCode().
func (p *Path) HashCode() int32 {
	if len(p.s) == 0 {
		return 0
	}
	if p.hash == 0 {
		var h int32
		for i := 0; i < len(p.s); i++ {
			h = 31*h + int32(p.s[i])
		}
		p.hash = h
	}
	return p.hash
}

// String is Path.toString().
func (p *Path) String() string { return p.s }

// Path_NormalizeString returns Path_Of(input).String() without allocating a
// Path struct on the common case where input has no NULs, lone surrogates,
// duplicate slashes, or trailing slashes.
func Path_NormalizeString(input string) string {
	return ParsePathString(input)
}

// ToAbsolutePath is Path.toAbsolutePath(): a relative path resolved against
// the working directory (Java's user.dir).
func (p *Path) ToAbsolutePath() *Path {
	if p.IsAbsolute() {
		return p
	}
	wd, err := os.Getwd()
	if err != nil {
		// Java reads user.dir once at startup and cannot fail here.
		panic(NewJavaLangError("cannot read the working directory: " + err.Error()))
	}
	return parsePath(wd).Resolve(p)
}

// Path_WorkingDir returns Path.of(System.getProperty("user.dir")) so callers
// resolving many relative paths in a single batch can call os.Getwd() once.
func Path_WorkingDir() *Path {
	wd, err := os.Getwd()
	if err != nil {
		panic(NewJavaLangError("cannot read the working directory: " + err.Error()))
	}
	return parsePath(wd)
}

// ToAbsolutePathWithWd resolves p against wd when p is relative.
func (p *Path) ToAbsolutePathWithWd(wd *Path) *Path {
	if p.IsAbsolute() {
		return p
	}
	return wd.Resolve(p)
}

// invalidPathException is java.nio.file.InvalidPathException.
type invalidPathException struct {
	IllegalArgumentException
}

func newInvalidPathException(input, reason string) *invalidPathException {
	e := &invalidPathException{}
	msg := reason + ": " + input
	InitThrowable(e, "java.nio.file.InvalidPathException", &msg, nil)
	return e
}
