// Port of com/google/javascript/jscomp/SourceFile.java @ closure-compiler v20261021

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
		if _, err := f.getCode(); err != nil {
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
		if _, err := f.getCode(); err != nil {
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

	localCode, err := f.getCode()
	if err != nil {
		f.lineOffsets = make([]int, 1)
		f.lineByteOffsets = make([]int, 1)
		return
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

func ioErrorMessage(err error) string {
	if t, ok := err.(jutil.Throwable); ok {
		return t.GetMessage()
	}
	return err.Error()
}

// GetCode gets all the code in this source file. It panics with an
// IOException if the code cannot be loaded.
//
// Java: SourceFile.getCode()
func (f *SourceFile) GetCode() string {
	code, err := f.getCode()
	if err != nil {
		panic(err)
	}
	return code
}

func (f *SourceFile) getCode() (string, jutil.IIOException) {
	if p := f.code.Load(); p != nil {
		return *p, nil
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if p := f.code.Load(); p != nil {
		return *p, nil
	}
	var raw string
	switch l := f.loader.(type) {
	case *SourceFile_CodeLoader_Preloaded:
		raw = l.preloadedCode
	case *SourceFile_CodeLoader_OnDisk:
		var err jutil.IIOException
		raw, err = l.loadCode()
		if err != nil {
			return "", err
		}
	case *SourceFile_CodeLoader_AtZip:
		var err jutil.IIOException
		raw, err = l.loadCode()
		if err != nil {
			return "", err
		}
	default:
		var err jutil.IIOException
		raw, err = sourceFile_loadFromCustomLoader(f.loader)
		if err != nil {
			return "", err
		}
	}
	f.SetCodeAndDoBookkeeping(&raw)
	return *f.code.Load(), nil
}

func sourceFile_loadFromCustomLoader(loader SourceFile_CodeLoader) (code string, ioErr jutil.IIOException) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(jutil.IIOException); ok {
				ioErr = e
				return
			}
			panic(r)
		}
	}()
	return loader.LoadUncachedCode(), nil
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
	jutil.CheckArgumentMsg1(offset >= 0, "offset must be non-negative: %s", offset)
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
	js, err := f.getCode()
	if err != nil {
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
	js, err := f.getCode()
	if err != nil {
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
	js, err := f.getCode()
	if err != nil {
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
	if strings.HasPrefix(path, "bin/") {
		return path[4:]
	}
	if strings.HasPrefix(path, "genfiles/") {
		return path[9:]
	}
	if strings.HasPrefix(path, "testlogs/") {
		return path[9:]
	}

	// Case 2: Contains blaze-out/ or bazel-out/
	blazeOutIdx := strings.Index(path, "blaze-out/")
	if blazeOutIdx == -1 {
		blazeOutIdx = strings.Index(path, "bazel-out/")
	}

	if blazeOutIdx != -1 {
		// Find the next slash after blaze-out/ or bazel-out/
		afterBlazeOut := path[blazeOutIdx+10:]
		if nextSlash := strings.IndexByte(afterBlazeOut, '/'); nextSlash != -1 {
			rest := afterBlazeOut[nextSlash+1:]
			if strings.HasPrefix(rest, "bin/") {
				return rest[4:]
			}
			if strings.HasPrefix(rest, "genfiles/") {
				return rest[9:]
			}
			if strings.HasPrefix(rest, "testlogs/") {
				return rest[9:]
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
