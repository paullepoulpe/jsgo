package javanio

import (
	"os"
	"strings"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

type InvalidPathException struct {
	Input  string
	Reason string
}

func (e *InvalidPathException) Error() string {
	return "java.nio.file.InvalidPathException: " + e.Reason + ": " + e.Input
}

// Path models java.nio.file.Path on the default Unix filesystem.
type Path struct {
	s          string
	hash       atomic.Int32
	elemStarts atomic.Pointer[[]int]
}

var (
	emptyPath       = &Path{s: ""}
	rootPath        = &Path{s: "/"}
	emptyElemStarts = []int{0}
	rootElemStarts  = []int{}
)

func init() {
	emptyPath.elemStarts.Store(&emptyElemStarts)
	rootPath.elemStarts.Store(&rootElemStarts)
	rootPath.hash.Store('/')
}

// Path_Of implements Path.of(first, more...) and Paths.get(first, more...).
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

func Paths_Get(first string, more ...string) *Path {
	return Path_Of(first, more...)
}

func FileToPath(f *javaio.File) *Path {
	if f == nil {
		panic(javalang.NewNullPointerExceptionMsg("file is null"))
	}
	return Path_Of(f.GetPath())
}

func ParsePathString(input string) string {
	slow := len(input) > 1 && input[len(input)-1] == '/'
	for i := 0; i < len(input); i++ {
		c := input[i]
		if c == 0 {
			panic(&InvalidPathException{Input: input, Reason: "Nul character not allowed"})
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
		if c == '/' && len(b) > 0 && b[len(b)-1] == '/' {
			continue
		}
		b = append(b, c)
	}
	if len(b) > 1 && b[len(b)-1] == '/' {
		b = b[:len(b)-1]
	}
	for i := 0; i+1 < len(b); i++ {
		if b[i] == 0xED && b[i+1] >= 0xA0 && b[i+1] <= 0xBF {
			panic(&InvalidPathException{
				Input:  string(b),
				Reason: "Malformed input or input contains unmappable characters",
			})
		}
	}
	return string(b)
}

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

func (p *Path) isEmpty() bool { return len(p.s) == 0 }

func (p *Path) elementOffsets() []int {
	if starts := p.elemStarts.Load(); starts != nil {
		return *starts
	}
	if p.isEmpty() {
		p.elemStarts.Store(&emptyElemStarts)
		return emptyElemStarts
	}
	s := p.s
	n := 0
	for i := 0; i < len(s); i++ {
		if s[i] != '/' && (i == 0 || s[i-1] == '/') {
			n++
		}
	}
	if n == 0 {
		p.elemStarts.Store(&rootElemStarts)
		return rootElemStarts
	}
	starts := make([]int, 0, n)
	for i := 0; i < len(s); i++ {
		if s[i] != '/' && (i == 0 || s[i-1] == '/') {
			starts = append(starts, i)
		}
	}
	p.elemStarts.Store(&starts)
	return starts
}

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

func (p *Path) IsAbsolute() bool {
	return len(p.s) > 0 && p.s[0] == '/'
}

func (p *Path) GetRoot() *Path {
	if p.IsAbsolute() {
		return rootPath
	}
	return nil
}

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

func (p *Path) GetNameCount() int {
	return len(p.elementOffsets())
}

func (p *Path) GetName(index int) *Path {
	if index < 0 || index >= p.GetNameCount() {
		panic(javalang.NewIllegalArgumentException(""))
	}
	return newPathFromString(p.element(index))
}

func (p *Path) Subpath(beginIndex, endIndex int) *Path {
	n := p.GetNameCount()
	if beginIndex < 0 || beginIndex >= n || endIndex > n || beginIndex >= endIndex {
		panic(javalang.NewIllegalArgumentException(""))
	}
	starts := p.elementOffsets()
	end := len(p.s)
	if endIndex < n {
		end = starts[endIndex] - 1
	}
	return newPathFromString(p.s[starts[beginIndex]:end])
}

func ResolveParsedPath(base, other string) string {
	switch {
	case len(other) > 0 && other[0] == '/':
		return other
	case len(other) == 0:
		return base
	case len(base) == 0:
		return other
	}
	if base[len(base)-1] == '/' {
		return base + other
	}
	return base + "/" + other
}

func (p *Path) Resolve(other *Path) *Path {
	if other == nil {
		panic(javalang.NewNullPointerExceptionMsg("other is null"))
	}
	switch {
	case other.IsAbsolute():
		return other
	case other.isEmpty():
		return p
	case p.isEmpty():
		return other
	}
	return newPathFromString(ResolveParsedPath(p.s, other.s))
}

func (p *Path) ResolveString(other string) *Path {
	return p.Resolve(parsePath(other))
}

func (p *Path) ResolveSibling(other *Path) *Path {
	if other == nil {
		panic(javalang.NewNullPointerExceptionMsg("other is null"))
	}
	parent := p.GetParent()
	if parent == nil {
		return other
	}
	return parent.Resolve(other)
}

func (p *Path) ResolveSiblingString(other string) *Path {
	return p.ResolveSibling(parsePath(other))
}

func (p *Path) Relativize(other *Path) *Path {
	if other == nil {
		panic(javalang.NewNullPointerExceptionMsg("other is null"))
	}
	if other.IsAbsolute() != p.IsAbsolute() {
		panic(javalang.NewIllegalArgumentException("'other' is different type of Path"))
	}
	if p.Equals(other) {
		return emptyPath
	}
	if p.isEmpty() {
		return other
	}
	base, target := p.Normalize(), other.Normalize()
	var baseNames, targetNames []string
	if !base.isEmpty() {
		baseNames = base.elements()
	}
	if !target.isEmpty() {
		targetNames = target.elements()
	}
	common := 0
	for common < len(baseNames) && common < len(targetNames) && baseNames[common] == targetNames[common] {
		common++
	}
	var b strings.Builder
	for _, name := range baseNames[common:] {
		if name == ".." {
			panic(javalang.NewIllegalArgumentException("Unable to compute relative  path from " + p.String() + " to " + other.String()))
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

func NormalizeParsedPath(s string) string {
	if !pathHasDotSegment(s) {
		return s
	}
	isAbs := len(s) > 0 && s[0] == '/'
	var kept []string
	changed := false
	for rest := s; len(rest) > 0; {
		if rest[0] == '/' {
			rest = rest[1:]
			continue
		}
		idx := strings.IndexByte(rest, '/')
		var name string
		if idx < 0 {
			name = rest
			rest = ""
		} else {
			name = rest[:idx]
			rest = rest[idx+1:]
		}
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
			if isAbs {
				changed = true
				continue
			}
		}
		kept = append(kept, name)
	}
	if !changed {
		return s
	}
	if len(kept) == 0 {
		if isAbs {
			return "/"
		}
		return ""
	}
	var b strings.Builder
	b.Grow(len(s))
	if isAbs {
		b.WriteByte('/')
	}
	for i, name := range kept {
		if i > 0 {
			b.WriteByte('/')
		}
		b.WriteString(name)
	}
	return b.String()
}

func (p *Path) Normalize() *Path {
	norm := NormalizeParsedPath(p.s)
	if norm == p.s {
		return p
	}
	return newPathFromString(norm)
}

func (p *Path) ToAbsolutePath() *Path {
	if p.IsAbsolute() {
		return p
	}
	wd, err := os.Getwd()
	if err != nil {
		return p
	}
	return newPathFromString(ResolveParsedPath(ParsePathString(wd), p.s))
}

func (p *Path) ToFile() *javaio.File {
	return javaio.NewFile(p.s)
}

func (p *Path) StartsWith(other *Path) bool {
	if other == nil {
		panic(javalang.NewNullPointerExceptionMsg("other is null"))
	}
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

func (p *Path) StartsWithString(other string) bool {
	return p.StartsWith(parsePath(other))
}

func (p *Path) EndsWith(other *Path) bool {
	if other == nil {
		panic(javalang.NewNullPointerExceptionMsg("other is null"))
	}
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

func (p *Path) EndsWithString(other string) bool {
	return p.EndsWith(parsePath(other))
}

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

func (p *Path) Equals(other any) bool {
	o, ok := other.(*Path)
	return ok && o != nil && p.s == o.s
}

func (p *Path) HashCode() int32 {
	if len(p.s) == 0 {
		return 0
	}
	h := p.hash.Load()
	if h == 0 {
		for i := 0; i < len(p.s); i++ {
			h = 31*h + int32(p.s[i])
		}
		p.hash.Store(h)
	}
	return h
}

func (p *Path) String() string { return p.s }
