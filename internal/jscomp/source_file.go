// Port of com/google/javascript/jscomp/SourceFile.java @ closure-compiler v20261015

/*
 * Copyright 2009 The Closure Compiler Authors.
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

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// SourceFile is an abstract representation of a source file that provides
// access to language-neutral features.
//
// The source file can be loaded from various locations, such as from disk
// or from a preloaded string. Loading is done as lazily as possible to
// minimize IO delays and memory cost of source text.
//
// Go-only: language-adapter — the code is a Go string, whose indices are
// bytes; the line offsets that the Java API exposes are in UTF-16 code
// units (lineOffsets), and their byte positions are kept alongside
// (lineByteOffsets) to cut the code at the same places.
type SourceFile struct {
	// The file name of the source file.
	//
	// It does not necessarily need to correspond to a real path. But it
	// should be unique. Will appear in warning messages emitted by the
	// compiler.
	fileName string

	kind rhino.StaticSourceFile_SourceKind

	loader SourceFile_CodeLoader

	isStubSourceFileForAlreadyProvidedInput bool

	// Source Line Information
	lineOffsets     []int // nil is null
	lineByteOffsets []int

	code atomic.Pointer[string]

	// Code statistics used for metrics.
	// For both of these, a negative value means we haven't counted yet.
	numLines int
	numBytes int

	isClosureUnawareCode bool

	// Go-only: language-adapter — synchronizes lazy code loading in
	// getCode and getCodeReader.
	mu sync.Mutex

	onDiskLoader SourceFile_CodeLoader_OnDisk

	sourceInfoTemplateMu    sync.Mutex
	sourceInfoTemplateReady atomic.Bool
	sourceInfoTemplate      rhino.Node
}

var _ rhino.StaticSourceFile = (*SourceFile)(nil)
var _ rhino.StaticSourceFileFast = (*SourceFile)(nil)

func (f *SourceFile) GetNamePtr() *string { return &f.fileName }

func (f *SourceFile) getOrCreateSourceInfoTemplate() rhino.Node {
	if f == nil {
		return rhino.NewNode(rhino.Token_SCRIPT)
	}
	if f.sourceInfoTemplateReady.Load() {
		return f.sourceInfoTemplate
	}
	f.sourceInfoTemplateMu.Lock()
	if !f.sourceInfoTemplateReady.Load() {
		// A template in the static store: a SourceFile can outlive the
		// compilation (tests share them) and must not keep its store alive.
		t := rhino.NewTemplateNode(rhino.Token_SCRIPT)
		t.SetStaticSourceFile(f)
		f.sourceInfoTemplate = t
		f.sourceInfoTemplateReady.Store(true)
	}
	f.sourceInfoTemplateMu.Unlock()
	return f.sourceInfoTemplate
}

const sourceFile_UTF8_BOM = "\uFEFF"

var sourceFile_BANG_SLASH = "!" + Platform_GetFileSeperator()

// Number of lines in the region returned by GetRegion(int). This length
// must be odd.
const sourceFile_SOURCE_EXCERPT_REGION_LENGTH = 5

func newSourceFile(loader SourceFile_CodeLoader, fileName string, kind rhino.StaticSourceFile_SourceKind) *SourceFile {
	if fileName == "" {
		panic(jutil.NewIllegalArgumentExceptionMsg("a source must have a name"))
	}

	if Platform_GetFileSeperator() != "/" {
		fileName = strings.ReplaceAll(fileName, Platform_GetFileSeperator(), "/")
	}

	return &SourceFile{
		loader:   loader,
		fileName: fileName,
		kind:     kind,
		numLines: -1,
		numBytes: -1,
	}
}

// Java: SourceFile.isStubSourceFileForAlreadyProvidedInput()
func (f *SourceFile) IsStubSourceFileForAlreadyProvidedInput() bool {
	return f.isStubSourceFileForAlreadyProvidedInput
}

func (f *SourceFile) SetIsStubSourceFileForAlreadyProvidedInput(isStubSourceFileForAlreadyProvidedInput bool) {
	f.isStubSourceFileForAlreadyProvidedInput = isStubSourceFileForAlreadyProvidedInput
}

// Java: SourceFile.getLineOffset(int)
func (f *SourceFile) GetLineOffset(lineno int) int {
	f.FindLineOffsets()
	if lineno < 1 || lineno > len(f.lineOffsets) {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"Expected line number between 1 and " + strconv.Itoa(len(f.lineOffsets)) + "\nActual: " + strconv.Itoa(lineno)))
	}
	return f.lineOffsets[lineno-1]
}

// GetNumLines returns the number of lines in the source file.
//
// NOTE: this may be stale if the SourceFile has been cleared and the file
// has changed on disk. If being fully accurate is required call GetCode and
// then call this method.
//
// Java: SourceFile.getNumLines()
func (f *SourceFile) GetNumLines() int {
	if f.numLines < 0 {
		// A negative value means we need to read in the code and calculate
		// this information. Otherwise, assume the file hasn't changed since we
		// read it.
		failed := false
		jutil.Try(func() {
			_ = f.GetCode()
		}, func(e jutil.IIOException) {
			// this is consistent with old behavior of this method.
			failed = true
		})
		if failed {
			return 0
		}
	}
	return f.numLines
}

// GetNumBytes returns the number of bytes in the source file.
//
// Oddly, bytes is defined as 'length of code as a Java String' this is
// accurate assuming files only use Latin1 characters but may be inaccurate
// outside of that.
//
// NOTE: this may be stale if the SourceFile has been cleared and the file
// has changed on disk. If being fully accurate is required call GetCode and
// then call this method.
//
// Java: SourceFile.getNumBytes()
func (f *SourceFile) GetNumBytes() int {
	if f.numBytes < 0 {
		// A negative value means we need to read in the code and calculate
		// this information. Otherwise, assume the file hasn't changed since we
		// read it.
		failed := false
		jutil.Try(func() {
			_ = f.GetCode()
		}, func(e jutil.IIOException) {
			// this is consistent with old behavior of this method.
			failed = true
		})
		if failed {
			return 0
		}
	}
	return f.numBytes
}

// Java: SourceFile.isLineTerminator(char)
func SourceFile_IsLineTerminator(ch uint16) bool {
	switch ch {
	case '\n', // Line Feed
		'\r',   // Carriage Return
		0x2028, // Line Separator
		0x2029: // Paragraph Separator
		return true
	default:
		return false
	}
}

// sourceFile_lineTerminatorAt is isLineTerminator(ch) for the char ch whose
// WTF-8 encoding starts at byte i of s: the length in bytes of that
// encoding (1 for \n and \r, 3 for \u2028 and \u2029), or 0 when ch is not
// a line terminator.
//
// Go-only: performance — the scans below read the WTF-8 bytes of the code
// instead of its UTF-16 chars.
func sourceFile_lineTerminatorAt(s string, i int) int {
	switch s[i] {
	case '\n', '\r':
		return 1
	case 0xE2:
		if i+2 < len(s) && s[i+1] == 0x80 && (s[i+2] == 0xA8 || s[i+2] == 0xA9) {
			return 3
		}
	}
	return 0
}

// sourceFile_lineTerminatorBefore is isLineTerminator(ch) for the char ch
// whose WTF-8 encoding ends at byte end of s: the length in bytes of that
// encoding, or 0 when ch is not a line terminator.
//
// Go-only: performance — see sourceFile_lineTerminatorAt.
func sourceFile_lineTerminatorBefore(s string, end int) int {
	switch s[end-1] {
	case '\n', '\r':
		return 1
	case 0xA8, 0xA9:
		if end >= 3 && s[end-3] == 0xE2 && s[end-2] == 0x80 {
			return 3
		}
	}
	return 0
}

// sourceFile_unitsOfByte is the number of UTF-16 units of the WTF-8
// sequence that starts with the byte b.
func sourceFile_unitsOfByte(b byte) int {
	if b >= 0xF0 && b < 0xF8 {
		return 2
	}
	if b&0xC0 == 0x80 {
		return 0 // continuation byte
	}
	return 1
}

func (f *SourceFile) FindLineOffsets() {
	if f.lineOffsets != nil {
		return
	}

	var localCode string
	if p := f.code.Load(); p != nil {
		localCode = *p
	} else {
		failed := false
		jutil.Try(func() {
			localCode = f.GetCode()
		}, func(e jutil.IIOException) {
			failed = true
		})
		if failed {
			f.lineOffsets = make([]int, 1)
			f.lineByteOffsets = make([]int, 1)
			return
		}
	}

	// getCode() updates numLines, so this is always in sync with the code.
	combined := make([]int, f.numLines*2)
	offsets := combined[:f.numLines:f.numLines]
	byteOffsets := combined[f.numLines:]
	index := 1 // start at 1 since the offset for line 0 is always at byte 0
	units := 0
	for i := 0; i < len(localCode); i++ {
		units += sourceFile_unitsOfByte(localCode[i])
		if n := sourceFile_lineTerminatorAt(localCode, i); n != 0 {
			ch := localCode[i]
			i += n - 1
			if ch == '\r' && i+1 < len(localCode) && localCode[i+1] == '\n' {
				i++
				units++
			}
			offsets[index] = units
			byteOffsets[index] = i + 1
			index++
		}
	}
	jutil.CheckState(index == len(offsets))
	f.lineOffsets = offsets
	f.lineByteOffsets = byteOffsets
}

// GetCode gets all the code in this source file. It panics with an
// IOException if the code cannot be loaded.
//
// Java: SourceFile.getCode()
func (f *SourceFile) GetCode() string {
	if p := f.code.Load(); p != nil {
		return *p
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if p := f.code.Load(); p != nil {
		return *p
	}
	f.SetCodeAndDoBookkeeping(jutil.StringPtr(f.loader.LoadUncachedCode()))
	return *f.code.Load()
}

// Java: SourceFile.setCodeDeprecated(String)
//
// Deprecated.
func (f *SourceFile) SetCodeDeprecated(code string) {
	f.SetCodeAndDoBookkeeping(&code)
}

// GetCodeReader gets a reader for the code in this source file. It panics
// with an IOException.
//
// Java: SourceFile.getCodeReader()
func (f *SourceFile) GetCodeReader() jutil.JavaReader {
	// Only synchronize if we need to
	if f.code.Load() == nil {
		var uncachedReader jutil.JavaReader
		func() {
			f.mu.Lock()
			defer f.mu.Unlock()
			// Make sure another thread hasn't loaded the code while we waited.
			if f.code.Load() == nil {
				uncachedReader = f.loader.OpenUncachedReader()
			}
		}()
		if uncachedReader != nil {
			return uncachedReader
		}
	}
	return jutil.NewStringReader(f.GetCode())
}

func (f *SourceFile) SetCodeAndDoBookkeeping(sourceCode *string) {
	f.code.Store(nil)
	// Force recalculation of all of these values when they are requested.
	f.lineOffsets = nil
	f.lineByteOffsets = nil

	if sourceCode != nil {
		s := *sourceCode
		if strings.HasPrefix(s, sourceFile_UTF8_BOM) {
			s = s[len(sourceFile_UTF8_BOM):]
		}

		// Update numLines and numBytes before publishing code via atomic Store.
		// NOTE: In some edit/refresh development workflows, we may end up
		// re-reading the file and getting different numbers than we got last
		// time we read it, so we should not have checkState() calls here to
		// assert they have not changed.
		// Misleading variable name. This really stores the 'number of utf16'
		// code points which is not the same as number of bytes.
		f.numBytes = jutil.StringLength(s)
		numLines := 1 // there is always at least one line
		for i := 0; i < len(s); i++ {
			if n := sourceFile_lineTerminatorAt(s, i); n != 0 {
				ch := s[i]
				i += n - 1
				if ch == '\r' && i+1 < len(s) && s[i+1] == '\n' {
					i++
				}
				numLines++
			}
		}
		f.numLines = numLines
		f.code.Store(&s)
	}
}

// GetOriginalPath is an alias of GetName().
//
// Deprecated: use GetName instead.
//
// Java: SourceFile.getOriginalPath()
func (f *SourceFile) GetOriginalPath() string {
	return f.GetName()
}

// ClearCachedSource flushes the cache, for SourceFile types which cache
// source code that can be regenerated easily.
//
// We maintain the cache mostly to speed up generating source when displaying
// error messages, so dumping the file contents after the compile is a fine
// thing to do.
//
// Java: SourceFile.clearCachedSource()
func (f *SourceFile) ClearCachedSource() {
	f.SetCodeAndDoBookkeeping(nil)
}

// Java: SourceFile.hasSourceInMemory()
func (f *SourceFile) HasSourceInMemory() bool {
	return f.code.Load() != nil
}

// GetName returns a unique name for the source file.
//
// This name is not required to be an actual file path on disk.
//
// Java: SourceFile.getName()
func (f *SourceFile) GetName() string {
	return f.fileName
}

// GetKind returns the source kind.
//
// Java: SourceFile.getKind()
func (f *SourceFile) GetKind() rhino.StaticSourceFile_SourceKind {
	return f.kind
}

// SetKind sets the source kind.
//
// Java: SourceFile.setKind(SourceKind)
func (f *SourceFile) SetKind(kind rhino.StaticSourceFile_SourceKind) {
	f.kind = kind
}

// Java: SourceFile.isClosureUnawareCode()
func (f *SourceFile) IsClosureUnawareCode() bool {
	return f.isClosureUnawareCode
}

// Java: SourceFile.markAsClosureUnawareCode()
func (f *SourceFile) MarkAsClosureUnawareCode() {
	f.isClosureUnawareCode = true
}

// Java: SourceFile.getLineOfOffset(int)
func (f *SourceFile) GetLineOfOffset(offset int) int {
	f.FindLineOffsets()
	search := jutil.ArraysBinarySearchInts(f.lineOffsets, offset)
	if search >= 0 {
		return search + 1 // lines are 1-based.
	}
	insertionPoint := -1 * (search + 1)
	return min(insertionPoint-1, len(f.lineOffsets)-1) + 1
}

// Java: SourceFile.getColumnOfOffset(int)
func (f *SourceFile) GetColumnOfOffset(offset int) int {
	line := f.GetLineOfOffset(offset)
	return offset - f.lineOffsets[line-1]
}

// GetLine gets the source line for the indicated line number (1 being the
// first line of the file). It does not include the newline at the end of
// the file. Returns nil if it does not exist, or if there was an IO
// exception.
//
// Java: SourceFile.getLine(int)
func (f *SourceFile) GetLine(lineNumber int) *string {
	var js string
	failed := false
	jutil.Try(func() {
		js = f.GetCode()
	}, func(e jutil.IIOException) {
		failed = true
	})
	if failed {
		return nil
	}
	f.FindLineOffsets()
	if lineNumber > len(f.lineOffsets) {
		return nil
	}

	if lineNumber < 1 {
		lineNumber = 1
	}

	pos := f.lineByteOffsets[lineNumber-1]

	if lineNumber < len(f.lineOffsets) {
		nextPos := f.lineByteOffsets[lineNumber]
		// The line terminator before nextPos.
		endPos := nextPos - sourceFile_lineTerminatorBefore(js, nextPos)
		if js[endPos] == '\n' && endPos > pos && js[endPos-1] == '\r' {
			endPos--
		}
		s := js[pos:endPos]
		return &s
	} else {
		// If next new line cannot be found, there are two cases
		// 1. pos already reaches the end of file, then null should be returned
		// 2. otherwise, return the contents between pos and the end of file.
		if pos >= len(js) {
			return nil
		} else {
			s := js[pos:]
			return &s
		}
	}
}

// GetLines gets the source lines starting at `lineNumber` (1 being the
// first line of the file) and continuing until `length` (the number of
// characters desired, starting at the 0th character of the specified line;
// if negative or 0, returns a single line). Omits any trailing newlines.
// Returns nil if it does not exist or if there was an IO exception.
//
// Java: SourceFile.getLines(int, int)
func (f *SourceFile) GetLines(lineNumber int, length int) Region {
	var js string
	failed := false
	jutil.Try(func() {
		js = f.GetCode()
	}, func(e jutil.IIOException) {
		failed = true
	})
	if failed {
		return nil
	}
	f.FindLineOffsets()
	if lineNumber > len(f.lineOffsets) {
		return nil
	}

	if lineNumber < 1 {
		lineNumber = 1
	}
	if length <= 0 {
		length = 1
	}

	jsLength := f.numBytes
	pos := f.lineOffsets[lineNumber-1]
	posByte := f.lineByteOffsets[lineNumber-1]
	if pos == jsLength {
		return NewSimpleRegion(lineNumber, lineNumber, "") // Happens when asking for the last empty line in a file.
	}
	endChar := pos
	endCharByte := posByte
	endLine := lineNumber
	// go through lines until we've reached the end of the file or met the
	// specified length.
	for ; endChar < pos+length && endLine <= len(f.lineOffsets); endLine++ {
		if endLine < len(f.lineOffsets) {
			endChar = f.lineOffsets[endLine]
			endCharByte = f.lineByteOffsets[endLine]
		} else {
			endChar = jsLength
			endCharByte = len(js)
		}
	}

	if endChar > pos {
		if js[endCharByte-1] == '\n' {
			if endChar-1 > pos && js[endCharByte-2] == '\r' {
				return NewSimpleRegion(lineNumber, endLine, js[posByte:endCharByte-2])
			}
			return NewSimpleRegion(lineNumber, endLine, js[posByte:endCharByte-1])
		} else if n := sourceFile_lineTerminatorBefore(js, endCharByte); n != 0 {
			return NewSimpleRegion(lineNumber, endLine, js[posByte:endCharByte-n])
		}
	}
	return NewSimpleRegion(lineNumber, endLine, js[posByte:endCharByte])
}

// GetRegion gets a region around the indicated line number (1 being the
// first line of the file). The exact definition of a region is
// implementation specific, but it must contain the line indicated by the
// line number. A region must not start or end by a carriage return. Returns
// nil if it does not exist, or if there was an IO exception.
//
// Java: SourceFile.getRegion(int)
func (f *SourceFile) GetRegion(lineNumber int) Region {
	js := ""
	failed := false
	jutil.Try(func() {
		js = f.GetCode()
	}, func(e jutil.IIOException) {
		failed = true
	})
	if failed {
		return nil
	}
	pos := 0
	startLine := max(1, lineNumber-(sourceFile_SOURCE_EXCERPT_REGION_LENGTH+1)/2+1)
	for n := 1; n < startLine; n++ {
		nextpos := strings.IndexByte(js[pos:], '\n')
		if nextpos == -1 {
			break
		}
		pos = pos + nextpos + 1
	}
	end := pos
	endLine := startLine
	for n := 0; n < sourceFile_SOURCE_EXCERPT_REGION_LENGTH; n, endLine = n+1, endLine+1 {
		i := strings.IndexByte(js[end:], '\n')
		if i == -1 {
			end = -1
			break
		}
		end = end + i + 1
	}
	if lineNumber >= endLine {
		return nil
	}
	if end == -1 {
		last := len(js) - 1
		if last < 0 {
			jutil.StringCharAt(js, -1) // throws StringIndexOutOfBoundsException
		}
		if js[last] == '\n' {
			return NewSimpleRegion(startLine, endLine, js[pos:last])
		}
		return NewSimpleRegion(startLine, endLine, js[pos:])
	}
	return NewSimpleRegion(startLine, endLine, js[pos:end])
}

// Java: SourceFile.toString()
func (f *SourceFile) String() string {
	return f.fileName
}

// SourceFile_FromZipFile panics with an IOException.
//
// Java: SourceFile.fromZipFile(String, Charset)
func SourceFile_FromZipFile(zipName string, inputCharset *jutil.Charset) []*SourceFile {
	data, err := os.ReadFile(zipName)
	if err != nil {
		panic(sourceFile_fileInputStreamException(err, zipName))
	}
	return SourceFile_FromZipInput(zipName, bytes.NewReader(data), inputCharset)
}

// sourceFile_fileInputStreamException is the exception new
// FileInputStream(name) throws: a FileNotFoundException whose message is
// "name (reason)".
//
// Go-only: language-adapter — maps an os error to the java.io exception.
func sourceFile_fileInputStreamException(err error, name string) jutil.Throwable {
	reason := err.Error()
	if pe, ok := err.(*os.PathError); ok {
		reason = pe.Err.Error()
	}
	if reason != "" {
		reason = strings.ToUpper(reason[:1]) + reason[1:]
	}
	msg := name + " (" + reason + ")"
	e := &jutil.IOException{}
	jutil.InitThrowable(e, "java.io.FileNotFoundException", &msg, nil)
	return e
}

// SourceFile_FromZipInput lists the .js entries of a zip stream. It panics
// with an IOException.
//
// Go-only: language-adapter — java.util.zip.ZipInputStream reads the local
// file headers in stream order; archive/zip reads the central directory,
// which lists the same entries in the same order in well-formed archives.
// Entry names without the UTF-8 flag are decoded with inputCharset, as
// ZipInputStream does.
//
// Java: SourceFile.fromZipInput(String, InputStream, Charset)
func SourceFile_FromZipInput(zipName string, input io.Reader, inputCharset *jutil.Charset) []*SourceFile {
	absoluteZipPath := jutil.Path_Of(zipName).ToAbsolutePath().String()
	var sourceFiles []*SourceFile

	data, err := io.ReadAll(input)
	if err != nil {
		panic(jutil.IOExceptionFromError(err, ""))
	}
	in, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		msg := err.Error()
		e := &jutil.IOException{}
		jutil.InitThrowable(e, "java.util.zip.ZipException", &msg, nil)
		panic(e)
	}
	for _, zipEntry := range in.File {
		entryName := zipEntry.Name
		if zipEntry.Flags&0x800 == 0 {
			entryName = inputCharset.Decode([]byte(zipEntry.Name))
		}
		if !strings.HasSuffix(entryName, ".js") { // Only accept js files
			continue
		}
		sourceFiles = append(sourceFiles,
			SourceFile_Builder_().
				WithCharset(inputCharset).
				WithOriginalPath(zipName+sourceFile_BANG_SLASH+entryName).
				WithZipEntryPath(absoluteZipPath, entryName).
				Build())
	}
	return sourceFiles
}

// Java: SourceFile.fromFile(String, Charset)
func SourceFile_FromFileWithCharset(fileName string, charset *jutil.Charset) *SourceFile {
	return SourceFile_Builder_().WithPath(fileName).WithCharset(charset).Build()
}

// Java: SourceFile.fromFile(String)
func SourceFile_FromFile(fileName string) *SourceFile {
	return SourceFile_Builder_().WithPath(fileName).Build()
}

// Java: SourceFile.fromPath(Path, Charset)
func SourceFile_FromPath(path *jutil.Path, charset *jutil.Charset) *SourceFile {
	return SourceFile_Builder_().WithPathFromPath(path).WithCharset(charset).Build()
}

// Java: SourceFile.fromCode(String, String, SourceKind)
func SourceFile_FromCodeWithKind(fileName string, code string, kind rhino.StaticSourceFile_SourceKind) *SourceFile {
	return SourceFile_Builder_().WithPath(fileName).WithKind(kind).WithContent(code).Build()
}

// Java: SourceFile.fromCode(String, String)
func SourceFile_FromCode(fileName string, code string) *SourceFile {
	return SourceFile_Builder_().WithPath(fileName).WithContent(code).Build()
}

// SourceFile_StubSourceFile returns a stub SourceFile with the given name
// and kind. It is used when the actual sources are provided via TypedASTs to
// stage2 and stage3.
//
// Java: SourceFile.stubSourceFile(String, SourceKind)
func SourceFile_StubSourceFile(fileName string, kind rhino.StaticSourceFile_SourceKind) *SourceFile {
	return SourceFile_Builder_().
		WithPath(fileName).
		WithKind(kind).
		SetIsStubSourceFileForAlreadyProvidedInput().
		Build()
}

// SourceFile_StripConfigSegment drops configuration-specific output root
// prefixes (such as "blaze-out/<config>/{bin,genfiles,testlogs}/",
// "bazel-out/<config>/{bin,genfiles,testlogs}/", or bare "bin/",
// "genfiles/", "testlogs/") from the given path.
//
// This allows paths to be matched consistently regardless of whether path
// mapping is enabled (see go/path-mapping-guide).
//
// Java: SourceFile.stripConfigSegment(String)
func SourceFile_StripConfigSegment(path string) string {
	// Case 1: Starts with bare config segment (e.g. mapped paths or relative roots)
	if jutil.StringStartsWith(path, "bin/") {
		return jutil.StringSubstringFrom(path, 4)
	}
	if jutil.StringStartsWith(path, "genfiles/") {
		return jutil.StringSubstringFrom(path, 9)
	}
	if jutil.StringStartsWith(path, "testlogs/") {
		return jutil.StringSubstringFrom(path, 9)
	}

	// Case 2: Contains blaze-out/ or bazel-out/
	blazeOutIdx := jutil.StringIndexOf(path, "blaze-out/")
	if blazeOutIdx == -1 {
		blazeOutIdx = jutil.StringIndexOf(path, "bazel-out/")
	}

	if blazeOutIdx != -1 {
		// Find the next slash after blaze-out/ or bazel-out/
		nextSlash := jutil.StringIndexOfCharFrom(path, '/', blazeOutIdx+10)
		if nextSlash != -1 {
			// String.startsWith(String, int)
			if jutil.StringRegionMatches(path, nextSlash+1, "bin/", 0, 4) {
				return jutil.StringSubstringFrom(path, nextSlash+5)
			}
			if jutil.StringRegionMatches(path, nextSlash+1, "genfiles/", 0, 9) {
				return jutil.StringSubstringFrom(path, nextSlash+10)
			}
			if jutil.StringRegionMatches(path, nextSlash+1, "testlogs/", 0, 9) {
				return jutil.StringSubstringFrom(path, nextSlash+10)
			}
		}
	}

	return path
}

// SourceFile_CanonicalizePathForMatching canonicalizes a file path for
// matching across compilation stages and TypedAST protos by stripping
// configuration prefix segments (see stripConfigSegment and
// go/path-mapping-guide).
//
// Java: SourceFile.canonicalizePathForMatching(String)
func SourceFile_CanonicalizePathForMatching(path string) string {
	return SourceFile_StripConfigSegment(path)
}

// RestoreCachedStateFrom reconciles serialized state in a SourceFileProto
// with the existing state in this file.
//
// This should be called whenever initializing a compilation based on
// TypedAST protos. For these compilations, the compiler initalization
// methods require creating SourceFiles before deserializing TypedAST
// protos, so we sometimes get two copies of the same SourceFile.)
//
// Java: SourceFile.restoreCachedStateFrom(SourceFileProto)
func (f *SourceFile) RestoreCachedStateFrom(protoSourceFile *SourceFileProto) {
	jutil.CheckStateMsg2(
		protoSourceFile.GetFilename() == f.GetName() ||
			SourceFile_CanonicalizePathForMatching(protoSourceFile.GetFilename()) ==
				SourceFile_CanonicalizePathForMatching(f.GetName()),
		"Cannot restore state for %s from %s",
		f.GetName(),
		protoSourceFile.GetFilename())
	// TypedAST proto information is expected to be more accurate for:
	//  1) whether a SourceFile contains an @extern annotation or not. In
	//    non-TypedAST builds, we allow passing @extern files under the --js
	//    flag. For TypedAST builds, we could support this, but it's an
	//    uncommon pattern and trickier to support than ban.
	//  2) tracking some extra state that is lazily computed in a SourceFile,
	//     like the number of lines and bytes in a file.
	//     SourceFile::restoreCachedStateFrom handles this case.
	// Note: the state in the proto might be incorrect in other cases, since
	// some state cannot be computed during library-level typechecking (e.g.
	// what files are in the weak chunk)
	if protoSourceFile.GetSourceKind() == SourceFileProto_SourceKind_EXTERN {
		jutil.CheckStateMsg1(
			f.GetKind() == rhino.StaticSourceFile_SourceKind_EXTERN,
			"TypedAST compilations must pass all extern files as externs, not js, but found %s",
			f.GetName())
	}

	// Restore the number of lines/bytes from the proto, unless we already
	// have cached numLines and numBytes. Offset by 1. the proto "unset"
	// value is 0, where as in SourceFile we use "-1"
	protoNumLines := int(protoSourceFile.GetNumLinesPlusOne()) - 1
	protoNumBytes := int(protoSourceFile.GetNumBytesPlusOne()) - 1
	// It's possible that this.numLines and protoNumLines are both not "-1"
	// but contain conflicting values. This would happen if a file changes on
	// disk after TypedAST proto serialization. We choose the proto value over
	// this.numLines because this data is intended for metrics recording,
	// which should care most about the size of a file when compilation began
	// and parsing ran.
	if protoNumLines != -1 {
		f.numLines = protoNumLines
	}
	if protoNumBytes != -1 {
		f.numBytes = protoNumBytes
	}
	f.isClosureUnawareCode = protoSourceFile.GetIsClosureUnawareCode()
}

// Java: SourceFile.fromProto(SourceFileProto)
func SourceFile_FromProto(protoSourceFile *SourceFileProto) *SourceFile {
	return sourceFileFromProtoWithSlab(protoSourceFile, nil)
}

func sourceFileFromProtoWithSlab(protoSourceFile *SourceFileProto, slab *jutil.Slab[SourceFile]) *SourceFile {
	sourceKind := sourceFile_getSourceKindFromProto(protoSourceFile)
	sourceFile := sourceFile_fromProtoWithKindAndSlab(protoSourceFile, sourceKind, slab)
	// Restore the number of lines/bytes, which are offset by 1 in the proto.
	sourceFile.numLines = int(protoSourceFile.GetNumLinesPlusOne()) - 1
	sourceFile.numBytes = int(protoSourceFile.GetNumBytesPlusOne()) - 1
	if protoSourceFile.GetIsClosureUnawareCode() {
		sourceFile.MarkAsClosureUnawareCode()
	}
	return sourceFile
}

// Java: SourceFile.fromProto(SourceFileProto, SourceKind)
func sourceFile_fromProtoWithKind(protoSourceFile *SourceFileProto, sourceKind rhino.StaticSourceFile_SourceKind) *SourceFile {
	return sourceFile_fromProtoWithKindAndSlab(protoSourceFile, sourceKind, nil)
}

func sourceFile_fromProtoWithKindAndSlab(
	protoSourceFile *SourceFileProto,
	sourceKind rhino.StaticSourceFile_SourceKind,
	slab *jutil.Slab[SourceFile],
) *SourceFile {
	switch protoSourceFile.GetLoaderCase() {
	case SourceFileProto_LoaderCase_PRELOADED_CONTENTS:
		return SourceFile_FromCodeWithKind(
			protoSourceFile.GetFilename(), protoSourceFile.GetPreloadedContents(), sourceKind)
	case SourceFileProto_LoaderCase_FILE_ON_DISK:
		fileOnDisk := protoSourceFile.GetFileOnDisk()
		pathOnDisk := fileOnDisk.GetActualPath()
		if pathOnDisk == "" {
			pathOnDisk = protoSourceFile.GetFilename()
		}
		charset := sourceFile_toCharset(fileOnDisk.GetCharset())
		if !sourceFile_isZipEntryPath(pathOnDisk) {
			return newOnDiskSourceFileFromSlab(slab, pathOnDisk, protoSourceFile.GetFilename(), charset, sourceKind)
		}
		var b SourceFile_Builder
		b.kind = sourceKind
		b.charset = charset
		return b.WithOriginalPath(protoSourceFile.GetFilename()).
			WithPath(pathOnDisk).
			Build()
	case SourceFileProto_LoaderCase_ZIP_ENTRY:
		zipEntry := protoSourceFile.GetZipEntry()
		var b SourceFile_Builder
		b.kind = sourceKind
		b.charset = sourceFile_toCharset(zipEntry.GetCharset())
		return b.WithOriginalPath(protoSourceFile.GetFilename()).
			WithZipEntryPath(zipEntry.GetZipPath(), zipEntry.GetEntryName()).
			Build()
	case SourceFileProto_LoaderCase_STUB_FILE:
		var b SourceFile_Builder
		b.kind = sourceKind
		b.charset = jutil.StandardCharsets_UTF_8
		return b.WithOriginalPath(protoSourceFile.GetFilename()).
			SetIsStubSourceFileForAlreadyProvidedInput().
			Build()
	case SourceFileProto_LoaderCase_LOADER_NOT_SET:
	}
	panic(jutil.NewAssertionErrorNoMessage())
}

func sourceFile_getSourceKindFromProto(protoSourceFile *SourceFileProto) rhino.StaticSourceFile_SourceKind {
	switch protoSourceFile.GetSourceKind() {
	case SourceFileProto_SourceKind_EXTERN:
		return rhino.StaticSourceFile_SourceKind_EXTERN
	case SourceFileProto_SourceKind_CODE:
		return rhino.StaticSourceFile_SourceKind_STRONG
	default: // NOT_SPECIFIED, UNRECOGNIZED
		panic(jutil.NewAssertionError(protoSourceFile.GetSourceKind().String()))
	}
}

func sourceFile_toCharset(protoCharset string) *jutil.Charset {
	if protoCharset == "" {
		return jutil.StandardCharsets_UTF_8
	}
	return jutil.Charset_ForName(protoCharset)
}

// SourceFile_Builder_ creates a new builder for source files
// (SourceFile.builder()).
//
// Go-only: language-adapter — the trailing underscore keeps the mechanical
// name SourceFile_Builder for the nested Builder type.
func SourceFile_Builder_() *SourceFile_Builder {
	return newSourceFile_Builder()
}

// SourceFile_Builder is a builder interface for source files.
//
// Allows users to customize the Charset, and the original path of the
// source file (if it differs from the path on disk).
type SourceFile_Builder struct {
	kind                                    rhino.StaticSourceFile_SourceKind
	charset                                 *jutil.Charset
	originalPath                            string
	hasOriginalPath                         bool
	isStubSourceFileForAlreadyProvidedInput bool

	path               string
	hasPath            bool
	pathWithFilesystem *jutil.Path
	zipEntryPath       string
	hasZipEntryPath    bool

	lazyContent func() string
}

func newSourceFile_Builder() *SourceFile_Builder {
	return &SourceFile_Builder{kind: rhino.StaticSourceFile_SourceKind_STRONG, charset: jutil.StandardCharsets_UTF_8}
}

// WithKind sets the source kind.
//
// Java: SourceFile.Builder.withKind(SourceKind)
func (b *SourceFile_Builder) WithKind(kind rhino.StaticSourceFile_SourceKind) *SourceFile_Builder {
	b.kind = kind
	return b
}

// WithCharset sets the charset to use when reading from an input stream or
// file.
//
// Java: SourceFile.Builder.withCharset(Charset)
func (b *SourceFile_Builder) WithCharset(charset *jutil.Charset) *SourceFile_Builder {
	b.charset = charset
	return b
}

// Java: SourceFile.Builder.withPath(String)
func (b *SourceFile_Builder) WithPath(path string) *SourceFile_Builder {
	return b.WithPathInternal(path, nil)
}

// Java: SourceFile.Builder.withPath(Path)
func (b *SourceFile_Builder) WithPathFromPath(path *jutil.Path) *SourceFile_Builder {
	return b.WithPathInternal(path.String(), path)
}

// Java: SourceFile.Builder.withContent(String)
func (b *SourceFile_Builder) WithContent(x string) *SourceFile_Builder {
	b.lazyContent = func() string { return x }
	return b
}

// Java: SourceFile.Builder.withContent(InputStream)
func (b *SourceFile_Builder) WithContentFromStream(x io.Reader) *SourceFile_Builder {
	b.lazyContent = func() string {
		jutil.CheckState(b.charset != nil)
		var result string
		jutil.Try(func() {
			result = jutil.CharStreams_ToString(jutil.NewInputStreamReader(x, b.charset))
		}, func(e jutil.IIOException) {
			panic(jutil.NewRuntimeExceptionWithCause(e))
		})
		return result
	}
	return b
}

// Java: SourceFile.Builder.withZipEntryPath(String, String)
func (b *SourceFile_Builder) WithZipEntryPath(zipPath string, entryPath string) *SourceFile_Builder {
	b.path = zipPath
	b.hasPath = true
	b.zipEntryPath = entryPath
	b.hasZipEntryPath = true
	return b
}

// WithOriginalPath sets a name for this source file that does not need to
// correspond to a path on disk.
//
// Allow passing a reasonable human-readable name in cases like for zip
// files and for generated files with unstable artifact prefixes.
//
// The name must still be unique.
//
// Java: SourceFile.Builder.withOriginalPath(String)
func (b *SourceFile_Builder) WithOriginalPath(originalPath string) *SourceFile_Builder {
	b.originalPath = originalPath
	b.hasOriginalPath = true
	return b
}

// SetIsStubSourceFileForAlreadyProvidedInput sets that this SourceFile is a
// stub file for an input that was already provided to the compiler in the
// form of TypedAST.
//
// Java: SourceFile.Builder.setIsStubSourceFileForAlreadyProvidedInput()
func (b *SourceFile_Builder) SetIsStubSourceFileForAlreadyProvidedInput() *SourceFile_Builder {
	b.isStubSourceFileForAlreadyProvidedInput = true
	return b
}

// Java: SourceFile.Builder.build()
func (b *SourceFile_Builder) Build() *SourceFile {
	var displayPath string
	if b.hasOriginalPath {
		displayPath = b.originalPath
	} else if !b.hasZipEntryPath {
		displayPath = sourceFile_nullableStringVal(b.path, b.hasPath)
	} else {
		displayPath = sourceFile_nullableStringVal(b.path, b.hasPath) + sourceFile_BANG_SLASH + b.zipEntryPath
	}

	if b.isStubSourceFileForAlreadyProvidedInput {
		file := newSourceFile(NewSourceFile_CodeLoader_StubSourceFileCodeLoader(), displayPath, b.kind)
		file.SetIsStubSourceFileForAlreadyProvidedInput(true)
		return file
	}

	if b.lazyContent != nil {
		return newSourceFile(NewSourceFile_CodeLoader_Preloaded(b.lazyContent()), displayPath, b.kind)
	}

	if b.hasZipEntryPath {
		return newSourceFile(
			NewSourceFile_CodeLoader_AtZip(b.path, b.zipEntryPath, b.charset),
			displayPath,
			b.kind)
	}

	if b.pathWithFilesystem != nil {
		return newOnDiskSourceFileWithNormalizedPath(nil, b.pathWithFilesystem.String(), displayPath, b.charset, b.kind)
	}
	if !b.hasPath {
		panic(jutil.NewNullPointerException())
	}
	return newOnDiskSourceFileFromSlab(nil, b.path, displayPath, b.charset, b.kind)
}

func sourceFile_nullableStringVal(s string, ok bool) string {
	if !ok {
		return "null"
	}
	return s
}

func sourceFile_isZipEntryPath(path string) bool {
	if strings.IndexByte(path, '!') < 0 {
		return false
	}
	bangSlashIndex := strings.Index(path, sourceFile_BANG_SLASH)
	if bangSlashIndex >= 0 {
		zipPath := path[:bangSlashIndex]
		entryPath := path[bangSlashIndex+len(sourceFile_BANG_SLASH):]
		return strings.HasSuffix(zipPath, ".zip") &&
			(strings.HasSuffix(entryPath, ".js") || strings.HasSuffix(entryPath, ".js.map"))
	}
	return false
}

func newOnDiskSourceFileFromSlab(
	slab *jutil.Slab[SourceFile],
	pathOnDisk string,
	displayPath string,
	charset *jutil.Charset,
	kind rhino.StaticSourceFile_SourceKind,
) *SourceFile {
	normalizedPath := jutil.ParsePathString(pathOnDisk)
	return newOnDiskSourceFileWithNormalizedPath(slab, normalizedPath, displayPath, charset, kind)
}

func newOnDiskSourceFileWithNormalizedPath(
	slab *jutil.Slab[SourceFile],
	normalizedPath string,
	displayPath string,
	charset *jutil.Charset,
	kind rhino.StaticSourceFile_SourceKind,
) *SourceFile {
	charsetName := charset.Name()
	if displayPath == "" {
		panic(jutil.NewIllegalArgumentExceptionMsg("a source must have a name"))
	}
	if Platform_GetFileSeperator() != "/" {
		displayPath = strings.ReplaceAll(displayPath, Platform_GetFileSeperator(), "/")
	}
	var f *SourceFile
	if slab != nil {
		f = slab.New()
	} else {
		f = &SourceFile{}
	}
	f.onDiskLoader.serializableCharset = charsetName
	f.onDiskLoader.relativePathStr = normalizedPath
	f.loader = &f.onDiskLoader
	f.fileName = displayPath
	f.kind = kind
	f.numLines = -1
	f.numBytes = -1
	return f
}

func (b *SourceFile_Builder) WithPathInternal(path string, pathWithFilesystem *jutil.Path) *SourceFile_Builder {
	// Check if this path should be inferred as a ZIP entry path.
	if strings.IndexByte(path, '!') >= 0 {
		bangSlashIndex := strings.Index(path, sourceFile_BANG_SLASH)
		if bangSlashIndex >= 0 {
			zipPath := path[:bangSlashIndex]
			entryPath := path[bangSlashIndex+len(sourceFile_BANG_SLASH):]

			if strings.HasSuffix(zipPath, ".zip") &&
				(strings.HasSuffix(entryPath, ".js") || strings.HasSuffix(entryPath, ".js.map")) {
				return b.WithZipEntryPath(zipPath, entryPath)
			}
		}
	}

	// Path instances have an implicit reference to a FileSystem. Make sure
	// to preserve it.
	b.path = path
	b.hasPath = true
	b.pathWithFilesystem = pathWithFilesystem

	return b
}

//////////////////////////////////////////////////////////////////////////////
// Implementations

// SourceFile_CodeLoader is SourceFile.CodeLoader.
type SourceFile_CodeLoader interface {
	// LoadUncachedCode returns the source text of this file from its
	// original storage. It panics with an IOException.
	//
	// The implementation may be a slow operation such as reading from a
	// file. SourceFile guarantees that this method is only called under
	// synchronization.
	LoadUncachedCode() string

	// OpenUncachedReader returns a Reader for the source text of this file
	// from its original storage, or nil. It panics with an IOException.
	OpenUncachedReader() jutil.JavaReader

	// ToProtoLocationBuilder returns a representation of this loader that can
	// be serialized/deserialized to reconstruct this SourceFile.
	ToProtoLocationBuilder(fileName string) *SourceFileProto_Builder
}

// SourceFile_CodeLoaderBase holds CodeLoader's base implementations.
type SourceFile_CodeLoaderBase struct{}

// Java: SourceFile.CodeLoader.loadUncachedCode()
func (SourceFile_CodeLoaderBase) LoadUncachedCode() string {
	panic(jutil.NewAssertionErrorNoMessage())
}

// Java: SourceFile.CodeLoader.openUncachedReader()
func (SourceFile_CodeLoaderBase) OpenUncachedReader() jutil.JavaReader {
	return nil
}

// SourceFile_CodeLoader_Preloaded is SourceFile.CodeLoader.Preloaded.
type SourceFile_CodeLoader_Preloaded struct {
	SourceFile_CodeLoaderBase
	preloadedCode string
}

// Java: SourceFile.CodeLoader.Preloaded(String)
func NewSourceFile_CodeLoader_Preloaded(preloadedCode string) *SourceFile_CodeLoader_Preloaded {
	return &SourceFile_CodeLoader_Preloaded{preloadedCode: preloadedCode}
}

// Java: SourceFile.CodeLoader.Preloaded.loadUncachedCode()
func (l *SourceFile_CodeLoader_Preloaded) LoadUncachedCode() string {
	return l.preloadedCode
}

// Java: SourceFile.CodeLoader.Preloaded.toProtoLocationBuilder(String)
func (l *SourceFile_CodeLoader_Preloaded) ToProtoLocationBuilder(fileName string) *SourceFileProto_Builder {
	return SourceFileProto_NewBuilder().SetPreloadedContents(l.preloadedCode)
}

// SourceFile_CodeLoader_StubSourceFileCodeLoader is a code loader which
// throws an error when loading code for a stub SourceFile.
type SourceFile_CodeLoader_StubSourceFileCodeLoader struct {
	SourceFile_CodeLoaderBase
}

// Java: SourceFile.CodeLoader.StubSourceFileCodeLoader()
func NewSourceFile_CodeLoader_StubSourceFileCodeLoader() *SourceFile_CodeLoader_StubSourceFileCodeLoader {
	return &SourceFile_CodeLoader_StubSourceFileCodeLoader{}
}

// Java: SourceFile.CodeLoader.StubSourceFileCodeLoader.loadUncachedCode()
func (l *SourceFile_CodeLoader_StubSourceFileCodeLoader) LoadUncachedCode() string {
	panic(jutil.NewUnsupportedOperationExceptionMsg("Attempting to load code from a stub SourceFile."))
}

// Java: SourceFile.CodeLoader.StubSourceFileCodeLoader.toProtoLocationBuilder(String)
func (l *SourceFile_CodeLoader_StubSourceFileCodeLoader) ToProtoLocationBuilder(fileName string) *SourceFileProto_Builder {
	return SourceFileProto_NewBuilder().SetFilename(fileName).SetStubFile(true)
}

// SourceFile_CodeLoader_OnDisk is SourceFile.CodeLoader.OnDisk.
type SourceFile_CodeLoader_OnDisk struct {
	SourceFile_CodeLoaderBase
	serializableCharset string
	relativePathStr     string
}

var sourceFileEmptyFileOnDisk = &SourceFileProto_FileOnDisk{frozen: true, memoizedSize: 0}

// Java: SourceFile.CodeLoader.OnDisk(Path, Charset)
func NewSourceFile_CodeLoader_OnDisk(relativePath *jutil.Path, c *jutil.Charset) *SourceFile_CodeLoader_OnDisk {
	return &SourceFile_CodeLoader_OnDisk{serializableCharset: c.Name(), relativePathStr: relativePath.String()}
}

// Java: SourceFile.CodeLoader.OnDisk.loadUncachedCode()
func (l *SourceFile_CodeLoader_OnDisk) LoadUncachedCode() string {
	var result string
	jutil.Try(func() {
		result = jutil.Files_ReadStringRaw(l.relativePathStr, l.GetCharset())
	}, func(e jutil.ICharacterCodingException) {
		msg := "Failed to read: " + l.relativePathStr + ", is this input UTF-8 encoded?"
		ex := &jutil.IOException{}
		jutil.InitThrowable(ex, "java.io.IOException", &msg, e)
		panic(ex)
	})
	return result
}

// Java: SourceFile.CodeLoader.OnDisk.openUncachedReader()
func (l *SourceFile_CodeLoader_OnDisk) OpenUncachedReader() jutil.JavaReader {
	return jutil.Files_NewBufferedReaderRaw(l.relativePathStr, l.GetCharset())
}

func (l *SourceFile_CodeLoader_OnDisk) GetCharset() *jutil.Charset {
	if l.serializableCharset == "UTF-8" {
		return jutil.StandardCharsets_UTF_8
	}
	return jutil.Charset_ForName(l.serializableCharset)
}

// Java: SourceFile.CodeLoader.OnDisk.toProtoLocationBuilder(String)
func (l *SourceFile_CodeLoader_OnDisk) ToProtoLocationBuilder(fileName string) *SourceFileProto_Builder {
	actualPath := l.relativePathStr
	// to save space, don't serialize the path if equal to the fileName.
	if fileName == actualPath {
		actualPath = ""
	}
	// save space by not serializing UTF_8 (the default charset)
	charset := l.serializableCharset
	if charset == "UTF-8" || l.GetCharset() == jutil.StandardCharsets_UTF_8 {
		charset = ""
	}
	return SourceFileProto_NewBuilder().
		SetFileOnDiskFromBuilder(
			SourceFileProto_FileOnDisk_NewBuilder().
				SetActualPath(actualPath).
				SetCharset(charset))
}

// SourceFile_CodeLoader_AtZip is SourceFile.CodeLoader.AtZip.
type SourceFile_CodeLoader_AtZip struct {
	SourceFile_CodeLoaderBase
	zipName             string
	entryName           string
	serializableCharset string
}

// Java: SourceFile.CodeLoader.AtZip(String, String, Charset)
func NewSourceFile_CodeLoader_AtZip(zipName string, entryName string, c *jutil.Charset) *SourceFile_CodeLoader_AtZip {
	return &SourceFile_CodeLoader_AtZip{zipName: zipName, entryName: entryName, serializableCharset: c.Name()}
}

// Java: SourceFile.CodeLoader.AtZip.loadUncachedCode()
func (l *SourceFile_CodeLoader_AtZip) LoadUncachedCode() string {
	return jutil.CharStreams_ToString(l.OpenUncachedReader())
}

// Java: SourceFile.CodeLoader.AtZip.openUncachedReader()
func (l *SourceFile_CodeLoader_AtZip) OpenUncachedReader() jutil.JavaReader {
	return jutil.NewInputStreamReader(
		JSCompZipFileCache_GetEntryStream(l.zipName, l.entryName), l.GetCharset())
}

func (l *SourceFile_CodeLoader_AtZip) GetCharset() *jutil.Charset {
	return jutil.Charset_ForName(l.serializableCharset)
}

// Java: SourceFile.CodeLoader.AtZip.toProtoLocationBuilder(String)
func (l *SourceFile_CodeLoader_AtZip) ToProtoLocationBuilder(fileName string) *SourceFileProto_Builder {
	// save space by not serializing UTF_8 (the default charset)
	charset := l.serializableCharset
	if l.GetCharset() == jutil.StandardCharsets_UTF_8 {
		charset = ""
	}
	return SourceFileProto_NewBuilder().
		SetFilename(fileName).
		SetZipEntry(
			SourceFileProto_ZipEntryOnDisk_NewBuilder().
				SetEntryName(l.entryName).
				SetZipPath(l.zipName).
				SetCharset(charset).
				Build())
}

// Java: SourceFile.getProto()
func (f *SourceFile) GetProto() *SourceFileProto {
	if l, ok := f.loader.(*SourceFile_CodeLoader_OnDisk); ok {
		actualPath := l.relativePathStr
		if f.fileName == actualPath {
			actualPath = ""
		}
		charset := l.serializableCharset
		if charset == "UTF-8" || l.GetCharset() == jutil.StandardCharsets_UTF_8 {
			charset = ""
		}
		fileOnDisk := sourceFileEmptyFileOnDisk
		if actualPath != "" || charset != "" {
			fileOnDisk = &SourceFileProto_FileOnDisk{
				actualPath_:  actualPath,
				charset_:     charset,
				frozen:       true,
				memoizedSize: -1,
			}
		}
		return &SourceFileProto{
			filename_:             f.fileName,
			fileOnDisk_:           fileOnDisk,
			sourceKind_:           sourceFile_sourceKindToProto(f.kind).GetNumber(),
			isClosureUnawareCode_: f.isClosureUnawareCode,
			numBytesPlusOne_:      int32(f.numBytes + 1),
			numLinesPlusOne_:      int32(f.numLines + 1),
			loaderCase_:           SourceFileProto_FILE_ON_DISK_FIELD_NUMBER,
			frozen:                true,
			memoizedSize:          -1,
		}
	}
	return f.loader.
		ToProtoLocationBuilder(f.GetName()).
		SetFilename(f.GetName()).
		SetSourceKind(sourceFile_sourceKindToProto(f.GetKind())).
		SetIsClosureUnawareCode(f.isClosureUnawareCode).
		SetNumLinesPlusOne(int32(f.numLines + 1)).
		SetNumBytesPlusOne(int32(f.numBytes + 1)).
		Build()
}

func sourceFile_sourceKindToProto(sourceKind rhino.StaticSourceFile_SourceKind) SourceFileProto_SourceKind {
	switch sourceKind {
	case rhino.StaticSourceFile_SourceKind_EXTERN:
		return SourceFileProto_SourceKind_EXTERN
	case rhino.StaticSourceFile_SourceKind_STRONG, rhino.StaticSourceFile_SourceKind_WEAK:
		return SourceFileProto_SourceKind_CODE
	default: // NON_CODE
		panic(jutil.NewAssertionError(sourceKind.String()))
	}
}

// Inherited StaticSourceFile default methods.

// Java: StaticSourceFile.isStrong()
func (f *SourceFile) IsStrong() bool { return rhino.StaticSourceFile_IsStrong(f) }

// Java: StaticSourceFile.isWeak()
func (f *SourceFile) IsWeak() bool { return rhino.StaticSourceFile_IsWeak(f) }

// Java: StaticSourceFile.isExtern()
func (f *SourceFile) IsExtern() bool { return rhino.StaticSourceFile_IsExtern(f) }

// Java: StaticSourceFile.isNonCode()
func (f *SourceFile) IsNonCode() bool { return rhino.StaticSourceFile_IsNonCode(f) }

// Java: StaticSourceFile.isTypeScriptSource()
func (f *SourceFile) IsTypeScriptSource() bool { return rhino.StaticSourceFile_IsTypeScriptSource(f) }
