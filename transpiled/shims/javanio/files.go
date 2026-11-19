package javanio

import (
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaregex"
)

// StandardOpenOption models java.nio.file.StandardOpenOption.
type StandardOpenOption int

const (
	StandardOpenOption_READ StandardOpenOption = iota
	StandardOpenOption_WRITE
	StandardOpenOption_APPEND
	StandardOpenOption_TRUNCATE_EXISTING
	StandardOpenOption_CREATE
	StandardOpenOption_CREATE_NEW
)

// FileTime models java.nio.file.attribute.FileTime.
type FileTime struct {
	t time.Time
}

func FileTime_FromMillis(millis int64) *FileTime {
	return &FileTime{t: time.UnixMilli(millis).UTC()}
}

func (ft *FileTime) ToMillis() int64 {
	if ft == nil {
		return 0
	}
	return ft.t.UnixMilli()
}

func (ft *FileTime) CompareTo(other *FileTime) int {
	if other == nil {
		panic(javalang.NewNullPointerExceptionMsg("other is null"))
	}
	if ft.t.Before(other.t) {
		return -1
	}
	if ft.t.After(other.t) {
		return 1
	}
	return 0
}

func (ft *FileTime) Equals(other any) bool {
	o, ok := other.(*FileTime)
	return ok && o != nil && ft != nil && ft.t.Equal(o.t)
}

func (ft *FileTime) HashCode() int32 {
	if ft == nil {
		return 0
	}
	ns := ft.t.UnixNano()
	return int32(uint64(ns) ^ (uint64(ns) >> 32))
}

func (ft *FileTime) String() string {
	if ft == nil {
		return ""
	}
	return ft.t.UTC().Format(time.RFC3339Nano)
}

// --- Files Static Helpers ---

func osPath(p *Path) string {
	if p == nil {
		panic(javalang.NewNullPointerExceptionMsg("path is null"))
	}
	if p.s == "" {
		return "."
	}
	return p.s
}

func Files_Exists(p *Path) bool {
	if p == nil {
		panic(javalang.NewNullPointerExceptionMsg("path is null"))
	}
	if p.s == "" {
		return true
	}
	_, err := os.Stat(p.s)
	return err == nil
}

func Files_IsDirectory(p *Path) bool {
	st, err := os.Stat(osPath(p))
	return err == nil && st.IsDir()
}

func Files_IsRegularFile(p *Path) bool {
	st, err := os.Stat(osPath(p))
	return err == nil && st.Mode().IsRegular()
}

func Files_GetLastModifiedTime(p *Path) (*FileTime, error) {
	st, err := os.Stat(osPath(p))
	if err != nil {
		return nil, javaio.NewIOException(err.Error())
	}
	return &FileTime{t: st.ModTime().UTC()}, nil
}

func Files_CreateDirectories(dir *Path) (*Path, error) {
	if dir == nil {
		panic(javalang.NewNullPointerExceptionMsg("dir is null"))
	}
	if dir.s == "" {
		return dir, nil
	}
	if err := os.MkdirAll(dir.s, 0o755); err != nil {
		return nil, javaio.NewIOException(err.Error())
	}
	return dir, nil
}

func Files_ReadAllBytes(p *Path) ([]byte, error) {
	b, err := os.ReadFile(osPath(p))
	if err != nil {
		return nil, javaio.NewIOException(err.Error())
	}
	return b, nil
}

const readStringScratchSize = 256 * 1024

var (
	readStringScratchChan = make(chan *[readStringScratchSize]byte, 16)
	readStringScratchPool = sync.Pool{
		New: func() any {
			return new([readStringScratchSize]byte)
		},
	}
)

func getReadStringScratch() *[readStringScratchSize]byte {
	select {
	case b := <-readStringScratchChan:
		return b
	default:
		return readStringScratchPool.Get().(*[readStringScratchSize]byte)
	}
}

func putReadStringScratch(b *[readStringScratchSize]byte) {
	select {
	case readStringScratchChan <- b:
	default:
		readStringScratchPool.Put(b)
	}
}

func writeISO88591Chunk(sb *strings.Builder, chunk []byte) {
	if isBytesASCII(chunk) {
		sb.Write(chunk)
		return
	}
	start := 0
	for i, x := range chunk {
		if x >= 0x80 {
			if i > start {
				sb.Write(chunk[start:i])
			}
			sb.WriteByte(0xC0 | (x >> 6))
			sb.WriteByte(0x80 | (x & 0x3F))
			start = i + 1
		}
	}
	if start < len(chunk) {
		sb.Write(chunk[start:])
	}
}

func readLargeFDStringStrict(fd int, bufp *[readStringScratchSize]byte, total int, cs *Charset) (string, error) {
	initGrow := 0
	if sz, ok := fstatRegularSize(fd); ok && sz >= int64(total) && sz <= 1<<30 {
		initGrow = int(sz)
	}
	switch cs {
	case StandardCharsets_UTF_8:
		var sb strings.Builder
		if initGrow > 0 {
			sb.Grow(initGrow)
		}
		validLen, decErr := validateUTF8Chunk(bufp[:total], false)
		if decErr != nil {
			closeFD(fd)
			putReadStringScratch(bufp)
			return "", decErr
		}
		if validLen > 0 {
			sb.Write(bufp[:validLen])
		}
		carry := total - validLen
		if carry > 0 {
			copy(bufp[:carry], bufp[validLen:total])
		}
		for {
			m, rerr := readFD(fd, bufp[carry:])
			if rerr != nil {
				closeFD(fd)
				putReadStringScratch(bufp)
				return "", javaio.NewIOException(rerr.Error())
			}
			atEOF := m == 0
			chunkTotal := carry + m
			if chunkTotal > 0 {
				vLen, err := validateUTF8Chunk(bufp[:chunkTotal], atEOF)
				if err != nil {
					closeFD(fd)
					putReadStringScratch(bufp)
					return "", err
				}
				if vLen > 0 {
					sb.Write(bufp[:vLen])
				}
				carry = chunkTotal - vLen
				if carry > 0 {
					copy(bufp[:carry], bufp[vLen:chunkTotal])
				}
			}
			if atEOF {
				closeFD(fd)
				putReadStringScratch(bufp)
				return sb.String(), nil
			}
		}

	case StandardCharsets_US_ASCII:
		if !isBytesASCII(bufp[:total]) {
			closeFD(fd)
			putReadStringScratch(bufp)
			return "", &MalformedInputException{InputLength: 1}
		}
		var sb strings.Builder
		if initGrow > 0 {
			sb.Grow(initGrow)
		}
		sb.Write(bufp[:total])
		for {
			m, rerr := readFD(fd, bufp[:])
			if rerr != nil {
				closeFD(fd)
				putReadStringScratch(bufp)
				return "", javaio.NewIOException(rerr.Error())
			}
			if m == 0 {
				closeFD(fd)
				putReadStringScratch(bufp)
				return sb.String(), nil
			}
			if !isBytesASCII(bufp[:m]) {
				closeFD(fd)
				putReadStringScratch(bufp)
				return "", &MalformedInputException{InputLength: 1}
			}
			sb.Write(bufp[:m])
		}

	case StandardCharsets_ISO_8859_1:
		var sb strings.Builder
		if initGrow > 0 {
			sb.Grow(initGrow)
		}
		writeISO88591Chunk(&sb, bufp[:total])
		for {
			m, rerr := readFD(fd, bufp[:])
			if rerr != nil {
				closeFD(fd)
				putReadStringScratch(bufp)
				return "", javaio.NewIOException(rerr.Error())
			}
			if m == 0 {
				closeFD(fd)
				putReadStringScratch(bufp)
				return sb.String(), nil
			}
			writeISO88591Chunk(&sb, bufp[:m])
		}
	}

	capHint := total * 2
	if initGrow > capHint {
		capHint = initGrow
	}
	full := make([]byte, total, capHint)
	copy(full, bufp[:total])
	for {
		m, rerr := readFD(fd, bufp[:])
		if rerr != nil {
			closeFD(fd)
			putReadStringScratch(bufp)
			return "", javaio.NewIOException(rerr.Error())
		}
		if m == 0 {
			closeFD(fd)
			putReadStringScratch(bufp)
			return cs.DecodeStrict(full)
		}
		full = append(full, bufp[:m]...)
	}
}

func readStringUTF8Streaming(f *os.File, initGrow int) (string, error) {
	bufp := getReadStringScratch()
	var sb strings.Builder
	if initGrow > 0 {
		sb.Grow(initGrow)
	}
	carry := 0
	for {
		m, rerr := f.Read(bufp[carry:])
		total := carry + m
		atEOF := rerr == io.EOF
		if rerr != nil && !atEOF {
			_ = f.Close()
			putReadStringScratch(bufp)
			return "", javaio.NewIOException(rerr.Error())
		}
		if total > 0 {
			validLen, decErr := validateUTF8Chunk(bufp[:total], atEOF)
			if decErr != nil {
				_ = f.Close()
				putReadStringScratch(bufp)
				return "", decErr
			}
			if validLen > 0 {
				sb.Write(bufp[:validLen])
			}
			carry = total - validLen
			if carry > 0 {
				copy(bufp[:carry], bufp[validLen:total])
			}
		}
		if atEOF {
			_ = f.Close()
			putReadStringScratch(bufp)
			return sb.String(), nil
		}
	}
}

func readStringASCIIStrictStreaming(f *os.File, initGrow int) (string, error) {
	bufp := getReadStringScratch()
	var sb strings.Builder
	if initGrow > 0 {
		sb.Grow(initGrow)
	}
	for {
		m, rerr := f.Read(bufp[:])
		atEOF := rerr == io.EOF
		if rerr != nil && !atEOF {
			_ = f.Close()
			putReadStringScratch(bufp)
			return "", javaio.NewIOException(rerr.Error())
		}
		if m > 0 {
			if !isBytesASCII(bufp[:m]) {
				_ = f.Close()
				putReadStringScratch(bufp)
				return "", &MalformedInputException{InputLength: 1}
			}
			sb.Write(bufp[:m])
		}
		if atEOF {
			_ = f.Close()
			putReadStringScratch(bufp)
			return sb.String(), nil
		}
	}
}

func readStringISO88591Streaming(f *os.File, initGrow int) (string, error) {
	bufp := getReadStringScratch()
	var sb strings.Builder
	if initGrow > 0 {
		sb.Grow(initGrow)
	}
	for {
		m, rerr := f.Read(bufp[:])
		atEOF := rerr == io.EOF
		if rerr != nil && !atEOF {
			_ = f.Close()
			putReadStringScratch(bufp)
			return "", javaio.NewIOException(rerr.Error())
		}
		if m > 0 {
			writeISO88591Chunk(&sb, bufp[:m])
		}
		if atEOF {
			_ = f.Close()
			putReadStringScratch(bufp)
			return sb.String(), nil
		}
	}
}

// ReadFileDecodedReplace reads path using pooled scratch buffers without allocating
// a duplicate full-file []byte alongside the decoded replacement string.
func ReadFileDecodedReplace(path string, cs *Charset) (string, error) {
	if cs == nil {
		cs = StandardCharsets_UTF_8
	}
	if hasFastFDRead {
		fd, err := openReadOnlyFD(path)
		if err != nil {
			return "", javaio.NewIOException(err.Error())
		}
		bufp := getReadStringScratch()
		total := 0
		for total < readStringScratchSize {
			n, rerr := readFD(fd, bufp[total:])
			if n > 0 {
				total += n
			}
			if rerr != nil {
				closeFD(fd)
				putReadStringScratch(bufp)
				return "", javaio.NewIOException(rerr.Error())
			}
			if n == 0 {
				closeFD(fd)
				s := cs.Decode(bufp[:total])
				putReadStringScratch(bufp)
				return s, nil
			}
		}
		initGrow := total * 2
		if sz, ok := fstatRegularSize(fd); ok && sz > int64(total) && sz <= 1<<30 {
			initGrow = int(sz)
		}
		if cs == StandardCharsets_ISO_8859_1 {
			var sb strings.Builder
			sb.Grow(initGrow)
			writeISO88591Chunk(&sb, bufp[:total])
			for {
				m, rerr := readFD(fd, bufp[:])
				if rerr != nil {
					closeFD(fd)
					putReadStringScratch(bufp)
					return "", javaio.NewIOException(rerr.Error())
				}
				if m == 0 {
					closeFD(fd)
					putReadStringScratch(bufp)
					return sb.String(), nil
				}
				writeISO88591Chunk(&sb, bufp[:m])
			}
		}
		full := make([]byte, total, initGrow)
		copy(full, bufp[:total])
		for {
			m, rerr := readFD(fd, bufp[:])
			if rerr != nil {
				closeFD(fd)
				putReadStringScratch(bufp)
				return "", javaio.NewIOException(rerr.Error())
			}
			if m == 0 {
				closeFD(fd)
				putReadStringScratch(bufp)
				return cs.Decode(full), nil
			}
			full = append(full, bufp[:m]...)
		}
	}
	f, err := os.Open(path)
	if err != nil {
		return "", javaio.NewIOException(err.Error())
	}
	fi, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return "", javaio.NewIOException(err.Error())
	}
	if fi.Mode().IsRegular() {
		sz := fi.Size()
		if sz == 0 {
			_ = f.Close()
			return "", nil
		}
		if sz > 0 && sz <= readStringScratchSize {
			bufp := getReadStringScratch()
			n, rerr := io.ReadFull(f, bufp[:sz])
			_ = f.Close()
			if rerr != nil && rerr != io.EOF && rerr != io.ErrUnexpectedEOF {
				putReadStringScratch(bufp)
				return "", javaio.NewIOException(rerr.Error())
			}
			s := cs.Decode(bufp[:n])
			putReadStringScratch(bufp)
			return s, nil
		}
		if sz > readStringScratchSize && sz <= 1<<30 && cs == StandardCharsets_ISO_8859_1 {
			return readStringISO88591Streaming(f, int(sz))
		}
	}
	b, err := io.ReadAll(f)
	_ = f.Close()
	if err != nil {
		return "", javaio.NewIOException(err.Error())
	}
	return cs.Decode(b), nil
}

func Files_ReadString(p *Path) (string, error) {
	return Files_ReadStringCharset(p, StandardCharsets_UTF_8)
}

func Files_ReadStringCharset(p *Path, cs *Charset) (string, error) {
	if cs == nil {
		panic(javalang.NewNullPointerExceptionMsg("charset is null"))
	}
	target := osPath(p)
	if hasFastFDRead {
		fd, err := openReadOnlyFD(target)
		if err != nil {
			return "", javaio.NewIOException(err.Error())
		}
		bufp := getReadStringScratch()
		total := 0
		for total < readStringScratchSize {
			n, rerr := readFD(fd, bufp[total:])
			if n > 0 {
				total += n
			}
			if rerr != nil {
				closeFD(fd)
				putReadStringScratch(bufp)
				return "", javaio.NewIOException(rerr.Error())
			}
			if n == 0 {
				closeFD(fd)
				s, decErr := cs.DecodeStrict(bufp[:total])
				putReadStringScratch(bufp)
				return s, decErr
			}
		}
		return readLargeFDStringStrict(fd, bufp, total, cs)
	}
	f, err := os.Open(target)
	if err != nil {
		return "", javaio.NewIOException(err.Error())
	}
	fi, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return "", javaio.NewIOException(err.Error())
	}
	if fi.Mode().IsRegular() {
		sz := fi.Size()
		if sz == 0 {
			_ = f.Close()
			return "", nil
		}
		if sz > 0 && sz <= readStringScratchSize {
			bufp := getReadStringScratch()
			n, rerr := io.ReadFull(f, bufp[:sz])
			_ = f.Close()
			if rerr != nil && rerr != io.EOF && rerr != io.ErrUnexpectedEOF {
				putReadStringScratch(bufp)
				return "", javaio.NewIOException(rerr.Error())
			}
			s, decErr := cs.DecodeStrict(bufp[:n])
			putReadStringScratch(bufp)
			return s, decErr
		}
		if sz > readStringScratchSize && sz <= 1<<30 {
			switch cs {
			case StandardCharsets_UTF_8:
				return readStringUTF8Streaming(f, int(sz))
			case StandardCharsets_US_ASCII:
				return readStringASCIIStrictStreaming(f, int(sz))
			case StandardCharsets_ISO_8859_1:
				return readStringISO88591Streaming(f, int(sz))
			}
			b := make([]byte, sz)
			n, rerr := io.ReadFull(f, b)
			_ = f.Close()
			if rerr != nil && rerr != io.EOF && rerr != io.ErrUnexpectedEOF {
				return "", javaio.NewIOException(rerr.Error())
			}
			return cs.DecodeStrict(b[:n])
		}
	}
	switch cs {
	case StandardCharsets_UTF_8:
		return readStringUTF8Streaming(f, 0)
	case StandardCharsets_US_ASCII:
		return readStringASCIIStrictStreaming(f, 0)
	case StandardCharsets_ISO_8859_1:
		return readStringISO88591Streaming(f, 0)
	}
	b, err := io.ReadAll(f)
	_ = f.Close()
	if err != nil {
		return "", javaio.NewIOException(err.Error())
	}
	return cs.DecodeStrict(b)
}

func Files_Write(p *Path, bytes []byte, opts ...StandardOpenOption) (*Path, error) {
	out, err := Files_NewOutputStream(p, opts...)
	if err != nil {
		return nil, err
	}
	if _, err := out.Write(bytes); err != nil {
		_ = out.Close()
		return nil, javaio.NewIOException(err.Error())
	}
	if err := out.Close(); err != nil {
		return nil, javaio.NewIOException(err.Error())
	}
	return p, nil
}

func Files_WriteString(p *Path, csq string, opts ...StandardOpenOption) (*Path, error) {
	return Files_WriteStringCharset(p, csq, StandardCharsets_UTF_8, opts...)
}

func Files_WriteStringCharset(p *Path, csq string, cs *Charset, opts ...StandardOpenOption) (*Path, error) {
	if cs == nil {
		panic(javalang.NewNullPointerExceptionMsg("charset is null"))
	}
	if (cs == StandardCharsets_UTF_8 && (javalang.IsASCII(csq) || isValidUTF8String(csq))) ||
		((cs == StandardCharsets_US_ASCII || cs == StandardCharsets_ISO_8859_1) && javalang.IsASCII(csq)) {
		out, err := Files_NewOutputStream(p, opts...)
		if err != nil {
			return nil, err
		}
		if _, err := out.WriteString(csq); err != nil {
			_ = out.Close()
			return nil, javaio.NewIOException(err.Error())
		}
		if err := out.Close(); err != nil {
			return nil, javaio.NewIOException(err.Error())
		}
		return p, nil
	}
	return Files_Write(p, cs.Encode(csq), opts...)
}

func Files_NewInputStream(p *Path) (*javaio.FileInputStream, error) {
	return javaio.NewFileInputStream(osPath(p))
}

func Files_NewOutputStream(p *Path, opts ...StandardOpenOption) (*javaio.FileOutputStream, error) {
	target := osPath(p)
	if len(opts) == 0 {
		return javaio.NewFileOutputStream(target)
	}
	flags := 0
	hasRead := false
	hasWrite := false
	for _, opt := range opts {
		switch opt {
		case StandardOpenOption_READ:
			hasRead = true
		case StandardOpenOption_WRITE:
			hasWrite = true
		case StandardOpenOption_APPEND:
			hasWrite = true
			flags |= os.O_APPEND
		case StandardOpenOption_TRUNCATE_EXISTING:
			flags |= os.O_TRUNC
		case StandardOpenOption_CREATE:
			flags |= os.O_CREATE
		case StandardOpenOption_CREATE_NEW:
			flags |= os.O_CREATE | os.O_EXCL
		}
	}
	if hasRead && hasWrite {
		flags |= os.O_RDWR
	} else if !hasRead {
		flags |= os.O_WRONLY
	}
	return javaio.NewFileOutputStreamFlags(target, flags)
}

func Files_NewBufferedReader(p *Path, cs *Charset) (*javaio.BufferedReader, error) {
	if cs == nil {
		cs = StandardCharsets_UTF_8
	}
	in, err := Files_NewInputStream(p)
	if err != nil {
		return nil, err
	}
	return javaio.NewBufferedReader(NewInputStreamReaderWithCharset(in, cs)), nil
}

func Files_NewBufferedWriter(p *Path, cs *Charset, opts ...StandardOpenOption) (*javaio.BufferedWriter, error) {
	if cs == nil {
		cs = StandardCharsets_UTF_8
	}
	out, err := Files_NewOutputStream(p, opts...)
	if err != nil {
		return nil, err
	}
	return javaio.NewBufferedWriter(NewOutputStreamWriterWithCharset(out, cs)), nil
}

func Files_Delete(p *Path) error {
	if err := os.Remove(osPath(p)); err != nil {
		return javaio.NewIOException(err.Error())
	}
	return nil
}

func Files_DeleteIfExists(p *Path) (bool, error) {
	err := os.Remove(osPath(p))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, javaio.NewIOException(err.Error())
}

// --- FileSystem, PathMatcher & WalkFileTree ---

type PathMatcher interface {
	Matches(path *Path) bool
}

type regexPathMatcher struct {
	re *regexp.Regexp
}

func (m *regexPathMatcher) Matches(path *Path) bool {
	if path == nil {
		return false
	}
	return m.re.MatchString(path.String())
}

type javaRegexPathMatcher struct {
	pat *javaregex.Pattern
}

func (m *javaRegexPathMatcher) Matches(path *Path) bool {
	if path == nil {
		return false
	}
	return m.pat.Matcher(path.String()).Matches()
}

type FileSystem struct{}

var defaultFileSystem = &FileSystem{}

func FileSystems_GetDefault() *FileSystem {
	return defaultFileSystem
}

func (fs *FileSystem) GetPath(first string, more ...string) *Path {
	return Path_Of(first, more...)
}

func (fs *FileSystem) GetSeparator() string {
	return "/"
}

func (fs *FileSystem) GetPathMatcher(syntaxAndPattern string) PathMatcher {
	colon := strings.IndexByte(syntaxAndPattern, ':')
	if colon <= 0 {
		panic(javalang.NewIllegalArgumentException("Syntax must be 'glob:...' or 'regex:...': " + syntaxAndPattern))
	}
	syntax := strings.ToLower(syntaxAndPattern[:colon])
	expr := syntaxAndPattern[colon+1:]
	switch syntax {
	case "glob":
		reStr, err := globToRegex(expr)
		if err != nil {
			panic(javalang.NewIllegalArgumentException(err.Error()))
		}
		re, err := regexp.Compile(reStr)
		if err != nil {
			panic(javalang.NewIllegalArgumentException(err.Error()))
		}
		return &regexPathMatcher{re: re}
	case "regex":
		pat := javaregex.Pattern_Compile(expr)
		return &javaRegexPathMatcher{pat: pat}
	default:
		panic(javalang.NewUnsupportedOperationException("Unsupported syntax: " + syntax))
	}
}

func globToRegex(glob string) (string, error) {
	var sb strings.Builder
	sb.WriteString("^")
	inGroup := 0
	for i := 0; i < len(glob); {
		c := glob[i]
		switch c {
		case '\\':
			if i+1 >= len(glob) {
				return "", &InvalidPathException{Input: glob, Reason: "No character to escape"}
			}
			sb.WriteString(regexp.QuoteMeta(string(glob[i+1])))
			i += 2
		case '*':
			if i+1 < len(glob) && glob[i+1] == '*' {
				sb.WriteString(".*")
				i += 2
			} else {
				sb.WriteString("[^/]*")
				i++
			}
		case '?':
			sb.WriteString("[^/]")
			i++
		case '[':
			sb.WriteString("[")
			i++
			if i < len(glob) && glob[i] == '!' {
				sb.WriteString("^/")
				i++
			} else if i < len(glob) && glob[i] == '^' {
				sb.WriteString("\\^")
				i++
			}
			closed := false
			for i < len(glob) {
				ch := glob[i]
				i++
				if ch == ']' {
					sb.WriteByte(']')
					closed = true
					break
				}
				if ch == '\\' {
					sb.WriteString("\\\\")
				} else {
					sb.WriteByte(ch)
				}
			}
			if !closed {
				return "", &InvalidPathException{Input: glob, Reason: "Missing ']'"}
			}
		case '{':
			inGroup++
			sb.WriteString("(?:")
			i++
		case '}':
			if inGroup > 0 {
				inGroup--
				sb.WriteByte(')')
			} else {
				sb.WriteString("\\}")
			}
			i++
		case ',':
			if inGroup > 0 {
				sb.WriteByte('|')
			} else {
				sb.WriteByte(',')
			}
			i++
		default:
			sb.WriteString(regexp.QuoteMeta(string(c)))
			i++
		}
	}
	if inGroup != 0 {
		return "", &InvalidPathException{Input: glob, Reason: "Missing '}'"}
	}
	sb.WriteString("$")
	return sb.String(), nil
}

type FileVisitOption int

const (
	FileVisitOption_FOLLOW_LINKS FileVisitOption = iota
)

type FileVisitResult int

const (
	FileVisitResult_CONTINUE FileVisitResult = iota
	FileVisitResult_TERMINATE
	FileVisitResult_SKIP_SUBTREE
	FileVisitResult_SKIP_SIBLINGS
)

type BasicFileAttributes struct {
	info os.FileInfo
}

func (a *BasicFileAttributes) IsDirectory() bool   { return a.info != nil && a.info.IsDir() }
func (a *BasicFileAttributes) IsRegularFile() bool { return a.info != nil && a.info.Mode().IsRegular() }
func (a *BasicFileAttributes) IsSymbolicLink() bool {
	return a.info != nil && (a.info.Mode()&os.ModeSymlink) != 0
}
func (a *BasicFileAttributes) Size() int64 {
	if a.info == nil {
		return 0
	}
	return a.info.Size()
}
func (a *BasicFileAttributes) LastModifiedTime() *FileTime {
	if a.info == nil {
		return FileTime_FromMillis(0)
	}
	return &FileTime{t: a.info.ModTime().UTC()}
}

type FileVisitor interface {
	PreVisitDirectory(dir *Path, attrs *BasicFileAttributes) (FileVisitResult, error)
	VisitFile(file *Path, attrs *BasicFileAttributes) (FileVisitResult, error)
	VisitFileFailed(file *Path, exc error) (FileVisitResult, error)
	PostVisitDirectory(dir *Path, exc error) (FileVisitResult, error)
}

type SimpleFileVisitor struct{}

func (SimpleFileVisitor) PreVisitDirectory(*Path, *BasicFileAttributes) (FileVisitResult, error) {
	return FileVisitResult_CONTINUE, nil
}
func (SimpleFileVisitor) VisitFile(*Path, *BasicFileAttributes) (FileVisitResult, error) {
	return FileVisitResult_CONTINUE, nil
}
func (SimpleFileVisitor) VisitFileFailed(_ *Path, exc error) (FileVisitResult, error) {
	return FileVisitResult_TERMINATE, exc
}
func (SimpleFileVisitor) PostVisitDirectory(_ *Path, exc error) (FileVisitResult, error) {
	if exc != nil {
		return FileVisitResult_TERMINATE, exc
	}
	return FileVisitResult_CONTINUE, nil
}

func Files_WalkFileTree(start *Path, followLinks bool, maxDepth int, visitor FileVisitor) (*Path, error) {
	if start == nil || visitor == nil {
		panic(javalang.NewNullPointerExceptionMsg("start or visitor is null"))
	}
	if maxDepth < 0 {
		panic(javalang.NewIllegalArgumentException("maxDepth < 0"))
	}
	_, err := walkTreeNode(start, followLinks, 0, maxDepth, visitor)
	return start, err
}

func walkTreeNode(cur *Path, followLinks bool, depth, maxDepth int, visitor FileVisitor) (FileVisitResult, error) {
	target := osPath(cur)
	var info os.FileInfo
	var err error
	if followLinks {
		info, err = os.Stat(target)
	} else {
		info, err = os.Lstat(target)
	}
	if err != nil {
		return visitor.VisitFileFailed(cur, javaio.NewIOException(err.Error()))
	}
	attrs := &BasicFileAttributes{info: info}
	if !info.IsDir() || depth >= maxDepth {
		return visitor.VisitFile(cur, attrs)
	}
	res, err := visitor.PreVisitDirectory(cur, attrs)
	if err != nil || res != FileVisitResult_CONTINUE {
		if res == FileVisitResult_SKIP_SUBTREE {
			return FileVisitResult_CONTINUE, nil
		}
		return res, err
	}
	entries, readErr := os.ReadDir(target)
	var dirErr error
	if readErr != nil {
		dirErr = javaio.NewIOException(readErr.Error())
	} else {
		names := make([]string, len(entries))
		for i, e := range entries {
			names[i] = e.Name()
		}
		sort.Strings(names)
		for _, name := range names {
			var child *Path
			if cur.s == "." {
				child = Path_Of("./" + name)
				// Preserve "./name" prefix like Java Files.walkFileTree(Path.of(".")).
				child = &Path{s: "./" + filepath.Clean(name)}
			} else {
				child = cur.ResolveString(name)
			}
			childRes, childErr := walkTreeNode(child, followLinks, depth+1, maxDepth, visitor)
			if childErr != nil {
				return FileVisitResult_TERMINATE, childErr
			}
			if childRes == FileVisitResult_TERMINATE {
				return FileVisitResult_TERMINATE, nil
			}
			if childRes == FileVisitResult_SKIP_SIBLINGS {
				break
			}
		}
	}
	postRes, postErr := visitor.PostVisitDirectory(cur, dirErr)
	if postErr != nil {
		return FileVisitResult_TERMINATE, postErr
	}
	return postRes, nil
}
