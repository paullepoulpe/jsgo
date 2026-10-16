package jutil

import "os"

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
	path []byte
	hash int32
	// offsets holds the start of each name element; nil until computed.
	offsets []int
}

// Path_Of is Path.of(first, more...): when more is not empty, the non-empty
// strings among first and more are joined with "/".
func Path_Of(first string, more ...string) *Path {
	input := first
	if len(more) > 0 {
		b := make([]byte, 0, len(first)+len(more)*8)
		for _, s := range append([]string{first}, more...) {
			if s == "" {
				continue
			}
			if len(b) > 0 {
				b = append(b, '/')
			}
			b = append(b, s...)
		}
		input = string(b)
	}
	return parsePath(input)
}

// parsePath checks a path string and removes its redundant slashes.
func parsePath(input string) *Path {
	for i := 0; i < len(input); i++ {
		if input[i] == 0 {
			panic(NewInvalidPathException(input, "Nul character not allowed"))
		}
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
	// A WTF-8 lone surrogate (ED A0..BF xx) has no UTF-8 encoding. Observed
	// behaviour (black-box): this error names the path with its redundant
	// slashes removed, the NUL error the string as given.
	for i := 0; i+1 < len(b); i++ {
		if b[i] == 0xED && b[i+1] >= 0xA0 && b[i+1] <= 0xBF {
			panic(NewInvalidPathException(string(b), "Malformed input or input contains unmappable characters"))
		}
	}
	return &Path{path: b}
}

func newPathFromBytes(b []byte) *Path { return &Path{path: b} }

var emptyPath = &Path{path: []byte{}}

func (p *Path) isEmpty() bool { return len(p.path) == 0 }

// elementOffsets returns the start offset of every name element.
func (p *Path) elementOffsets() []int {
	if p.offsets != nil {
		return p.offsets
	}
	offsets := []int{}
	if p.isEmpty() {
		offsets = append(offsets, 0)
	} else {
		for i := 0; i < len(p.path); i++ {
			if p.path[i] != '/' && (i == 0 || p.path[i-1] == '/') {
				offsets = append(offsets, i)
			}
		}
	}
	p.offsets = offsets
	return offsets
}

// element returns the bytes of name element i.
func (p *Path) element(i int) []byte {
	offsets := p.elementOffsets()
	end := len(p.path)
	if i+1 < len(offsets) {
		end = offsets[i+1] - 1
	}
	return p.path[offsets[i]:end]
}

func (p *Path) elements() [][]byte {
	n := len(p.elementOffsets())
	out := make([][]byte, n)
	for i := range out {
		out[i] = p.element(i)
	}
	return out
}

// GetRoot is Path.getRoot(): "/" for an absolute path, else null.
func (p *Path) GetRoot() *Path {
	if p.IsAbsolute() {
		return newPathFromBytes([]byte{'/'})
	}
	return nil
}

// GetFileName is Path.getFileName(): the last name element, null for the
// root, and the empty path itself for the empty path.
func (p *Path) GetFileName() *Path {
	if p.isEmpty() {
		return p
	}
	offsets := p.elementOffsets()
	if len(offsets) == 0 {
		return nil
	}
	if len(offsets) == 1 && p.path[0] != '/' {
		return p
	}
	return newPathFromBytes(p.path[offsets[len(offsets)-1]:])
}

// GetParent is Path.getParent(): the path without its last name element,
// or null when there is none (the root, a one-element relative path).
func (p *Path) GetParent() *Path {
	offsets := p.elementOffsets()
	switch {
	case len(offsets) == 0:
		return nil
	case len(offsets) == 1:
		return p.GetRoot()
	}
	return newPathFromBytes(p.path[:offsets[len(offsets)-1]-1])
}

// GetNameCount is Path.getNameCount().
func (p *Path) GetNameCount() int { return len(p.elementOffsets()) }

// GetName is Path.getName(int).
func (p *Path) GetName(index int) *Path {
	if index < 0 || index >= p.GetNameCount() {
		panic(NewIllegalArgumentException())
	}
	return newPathFromBytes(p.element(index))
}

// Subpath is Path.subpath(int, int): the elements beginIndex (inclusive)
// to endIndex (exclusive), as a relative path.
func (p *Path) Subpath(beginIndex, endIndex int) *Path {
	n := p.GetNameCount()
	if beginIndex < 0 || beginIndex >= n || endIndex > n || beginIndex >= endIndex {
		panic(NewIllegalArgumentException())
	}
	offsets := p.elementOffsets()
	end := len(p.path)
	if endIndex < n {
		end = offsets[endIndex] - 1
	}
	return newPathFromBytes(p.path[offsets[beginIndex]:end])
}

// IsAbsolute is Path.isAbsolute(): whether the path has the root.
func (p *Path) IsAbsolute() bool { return len(p.path) > 0 && p.path[0] == '/' }

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
	b := make([]byte, 0, len(p.path)+1+len(obj.path))
	b = append(b, p.path...)
	if b[len(b)-1] != '/' {
		b = append(b, '/')
	}
	return newPathFromBytes(append(b, obj.path...))
}

// ResolveString is Path.resolve(String).
func (p *Path) ResolveString(other string) *Path { return p.Resolve(parsePath(other)) }

// ResolveSibling is Path.resolveSibling(Path): other resolved against the
// parent, or other itself when there is no parent.
func (p *Path) ResolveSibling(other *Path) *Path {
	if other == nil {
		panic(NewNullPointerException())
	}
	parent := p.GetParent()
	if parent == nil {
		return other
	}
	return parent.Resolve(other)
}

// ResolveSiblingString is Path.resolveSibling(String).
func (p *Path) ResolveSiblingString(other string) *Path {
	return p.ResolveSibling(parsePath(other))
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
	for common < len(baseNames) && common < len(targetNames) && string(baseNames[common]) == string(targetNames[common]) {
		common++
	}
	var b []byte
	for _, name := range baseNames[common:] {
		if string(name) == ".." {
			panic(NewIllegalArgumentExceptionMsg("Unable to compute relative  path from " + p.String() + " to " + obj.String()))
		}
		if len(b) > 0 {
			b = append(b, '/')
		}
		b = append(b, '.', '.')
	}
	for _, name := range targetNames[common:] {
		if len(b) > 0 {
			b = append(b, '/')
		}
		b = append(b, name...)
	}
	if b == nil {
		return emptyPath
	}
	return newPathFromBytes(b)
}

// names returns the name elements, none for the empty path.
func (p *Path) names() [][]byte {
	if p.isEmpty() {
		return nil
	}
	return p.elements()
}

// Normalize is Path.normalize() (see the observed behaviour above).
func (p *Path) Normalize() *Path {
	var kept [][]byte
	changed := false
	for _, name := range p.names() {
		switch string(name) {
		case ".":
			changed = true
			continue
		case "..":
			if len(kept) > 0 && string(kept[len(kept)-1]) != ".." {
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
	var b []byte
	if p.IsAbsolute() {
		b = append(b, '/')
	}
	for i, name := range kept {
		if i > 0 {
			b = append(b, '/')
		}
		b = append(b, name...)
	}
	if b == nil {
		return emptyPath
	}
	return newPathFromBytes(b)
}

// StartsWith is Path.startsWith(Path): same root and the same leading name
// elements.
func (p *Path) StartsWith(other *Path) bool {
	if p.IsAbsolute() != other.IsAbsolute() {
		return false
	}
	mine, theirs := p.elements(), other.elements()
	if p.isEmpty() || other.isEmpty() {
		return p.isEmpty() && other.isEmpty()
	}
	if len(theirs) > len(mine) {
		return false
	}
	for i := range theirs {
		if string(mine[i]) != string(theirs[i]) {
			return false
		}
	}
	return true
}

// StartsWithString is Path.startsWith(String).
func (p *Path) StartsWithString(other string) bool { return p.StartsWith(parsePath(other)) }

// EndsWith is Path.endsWith(Path): the same trailing name elements; an
// absolute other must be equal to this path.
func (p *Path) EndsWith(other *Path) bool {
	if other.IsAbsolute() {
		return p.Equals(other)
	}
	if p.isEmpty() || other.isEmpty() {
		return p.isEmpty() && other.isEmpty()
	}
	mine, theirs := p.elements(), other.elements()
	if len(theirs) > len(mine) {
		return false
	}
	off := len(mine) - len(theirs)
	for i := range theirs {
		if string(mine[off+i]) != string(theirs[i]) {
			return false
		}
	}
	return true
}

// EndsWithString is Path.endsWith(String).
func (p *Path) EndsWithString(other string) bool { return p.EndsWith(parsePath(other)) }

// CompareTo is Path.compareTo(Path): the unsigned bytes, lexicographically.
func (p *Path) CompareTo(other *Path) int {
	a, b := p.path, other.path
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
	return ok && o != nil && string(p.path) == string(o.path)
}

// HashCode is Path.hashCode().
func (p *Path) HashCode() int32 {
	if p.hash == 0 {
		var h int32
		for _, c := range p.path {
			h = 31*h + int32(c)
		}
		p.hash = h
	}
	return p.hash
}

// String is Path.toString().
func (p *Path) String() string { return string(p.path) }

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
