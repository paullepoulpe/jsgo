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
	"bytes"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

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
	jutil.CheckStateMsg1(protoNumLines >= -1, "protoNumLines must be >= -1: %s", protoNumLines)
	jutil.CheckStateMsg1(protoNumBytes >= -1, "protoNumBytes must be >= -1: %s", protoNumBytes)
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
	sourceKind := sourceFile_getSourceKindFromProto(protoSourceFile)
	sourceFile := sourceFile_fromProtoWithKind(protoSourceFile, sourceKind)
	// Restore the number of lines/bytes, which are offset by 1 in the proto.
	protoNumLines := int(protoSourceFile.GetNumLinesPlusOne()) - 1
	protoNumBytes := int(protoSourceFile.GetNumBytesPlusOne()) - 1
	jutil.CheckStateMsg1(protoNumLines >= -1, "protoNumLines must be >= -1: %s", protoNumLines)
	jutil.CheckStateMsg1(protoNumBytes >= -1, "protoNumBytes must be >= -1: %s", protoNumBytes)
	sourceFile.numLines = protoNumLines
	sourceFile.numBytes = protoNumBytes
	if protoSourceFile.GetIsClosureUnawareCode() {
		sourceFile.MarkAsClosureUnawareCode()
	}
	return sourceFile
}

// Java: SourceFile.fromProto(SourceFileProto, SourceKind)
func sourceFile_fromProtoWithKind(
	protoSourceFile *SourceFileProto,
	sourceKind rhino.StaticSourceFile_SourceKind,
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
		var b SourceFile_Builder
		b.kind = sourceKind
		b.charset = sourceFile_toCharset(fileOnDisk.GetCharset())
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

	lazyContent       func() string
	reloadableContent func() string
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
		data, err := io.ReadAll(x)
		if err != nil {
			panic(jutil.NewRuntimeExceptionWithCause(jutil.IOExceptionFromError(err, "")))
		}
		return jutil.CharStreams_ToString(jutil.NewInputStreamReader(bytes.NewReader(data), b.charset))
	}
	return b
}

// WithReloadableContent sets the content to what load returns; load may be
// called again later.
//
// Go-only: performance — a file built from a stream keeps its text for the
// whole compilation (the Preloaded loader holds it). With a loader that can
// read it again, ClearCachedSource after parsing really frees the text, and
// it is reloaded if something needs it (an error message, a source map).
func (b *SourceFile_Builder) WithReloadableContent(load func() string) *SourceFile_Builder {
	b.reloadableContent = load
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

	if b.reloadableContent != nil {
		return newSourceFile(&sourceFile_CodeLoader_Reloadable{load: b.reloadableContent}, displayPath, b.kind)
	}

	if b.hasZipEntryPath {
		return newSourceFile(
			NewSourceFile_CodeLoader_AtZip(b.path, b.zipEntryPath, b.charset),
			displayPath,
			b.kind)
	}

	if b.pathWithFilesystem != nil {
		return newOnDiskSourceFileWithNormalizedPath(b.pathWithFilesystem.String(), displayPath, b.charset, b.kind)
	}
	if !b.hasPath {
		panic(jutil.NewNullPointerException())
	}
	return newOnDiskSourceFile(b.path, displayPath, b.charset, b.kind)
}

func sourceFile_nullableStringVal(s string, ok bool) string {
	if !ok {
		return "null"
	}
	return s
}

func newOnDiskSourceFile(
	pathOnDisk string,
	displayPath string,
	charset *jutil.Charset,
	kind rhino.StaticSourceFile_SourceKind,
) *SourceFile {
	normalizedPath := jutil.ParsePathString(pathOnDisk)
	return newOnDiskSourceFileWithNormalizedPath(normalizedPath, displayPath, charset, kind)
}

func newOnDiskSourceFileWithNormalizedPath(
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
	f := &SourceFile{}
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

// SourceFile_CodeLoader_Preloaded is SourceFile.CodeLoader.Preloaded.
type SourceFile_CodeLoader_Preloaded struct {
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

func (*SourceFile_CodeLoader_Preloaded) OpenUncachedReader() jutil.JavaReader {
	return nil
}

// Java: SourceFile.CodeLoader.Preloaded.toProtoLocationBuilder(String)
func (l *SourceFile_CodeLoader_Preloaded) ToProtoLocationBuilder(fileName string) *SourceFileProto_Builder {
	return SourceFileProto_NewBuilder().SetPreloadedContents(l.preloadedCode)
}

// sourceFile_CodeLoader_Reloadable loads the code with a function that can
// be called again (see SourceFile_Builder.WithReloadableContent).
//
// Go-only: performance.
type sourceFile_CodeLoader_Reloadable struct {
	load func() string
}

func (l *sourceFile_CodeLoader_Reloadable) LoadUncachedCode() string {
	return l.load()
}

func (*sourceFile_CodeLoader_Reloadable) OpenUncachedReader() jutil.JavaReader {
	return nil
}

// ToProtoLocationBuilder serializes the code itself, as Preloaded does.
func (l *sourceFile_CodeLoader_Reloadable) ToProtoLocationBuilder(fileName string) *SourceFileProto_Builder {
	return SourceFileProto_NewBuilder().SetPreloadedContents(l.load())
}

// sourceFile_CodeLoader_InlineBase64 defers base64 decoding of an inline
// //# sourceMappingURL=data:application/json;base64,... payload until its
// source map contents are actually requested (Stage 3 FINALIZATIONS).
//
// Go-only: performance.
type sourceFile_CodeLoader_InlineBase64 struct {
	rawBase64     string
	canonicalUTF8 bool
}

func newInlineBase64SourceFile(fileName string, rawBase64 string, canonicalUTF8 bool) *SourceFile {
	return newSourceFile(
		&sourceFile_CodeLoader_InlineBase64{
			rawBase64:     rawBase64,
			canonicalUTF8: canonicalUTF8,
		},
		fileName,
		rhino.StaticSourceFile_SourceKind_NON_CODE)
}

func (l *sourceFile_CodeLoader_InlineBase64) LoadUncachedCode() string {
	return jutil.NewStringFromUTF8Bytes(decodeInlineBase64(l.rawBase64))
}

func (*sourceFile_CodeLoader_InlineBase64) OpenUncachedReader() jutil.JavaReader {
	return nil
}

func (l *sourceFile_CodeLoader_InlineBase64) ToProtoLocationBuilder(fileName string) *SourceFileProto_Builder {
	return SourceFileProto_NewBuilder().SetPreloadedContents(l.LoadUncachedCode())
}

// SourceFile_CodeLoader_StubSourceFileCodeLoader is a code loader which
// throws an error when loading code for a stub SourceFile.
type SourceFile_CodeLoader_StubSourceFileCodeLoader struct{}

// Java: SourceFile.CodeLoader.StubSourceFileCodeLoader()
func NewSourceFile_CodeLoader_StubSourceFileCodeLoader() *SourceFile_CodeLoader_StubSourceFileCodeLoader {
	return &SourceFile_CodeLoader_StubSourceFileCodeLoader{}
}

// Java: SourceFile.CodeLoader.StubSourceFileCodeLoader.loadUncachedCode()
func (l *SourceFile_CodeLoader_StubSourceFileCodeLoader) LoadUncachedCode() string {
	panic(jutil.NewUnsupportedOperationExceptionMsg("Attempting to load code from a stub SourceFile."))
}

func (*SourceFile_CodeLoader_StubSourceFileCodeLoader) OpenUncachedReader() jutil.JavaReader {
	return nil
}

// Java: SourceFile.CodeLoader.StubSourceFileCodeLoader.toProtoLocationBuilder(String)
func (l *SourceFile_CodeLoader_StubSourceFileCodeLoader) ToProtoLocationBuilder(fileName string) *SourceFileProto_Builder {
	return SourceFileProto_NewBuilder().SetFilename(fileName).SetStubFile(true)
}

// SourceFile_CodeLoader_OnDisk is SourceFile.CodeLoader.OnDisk.
type SourceFile_CodeLoader_OnDisk struct {
	serializableCharset string
	relativePathStr     string
}

var sourceFileEmptyFileOnDisk = &SourceFileProto_FileOnDisk{frozen: true, memoizedSize: 0}

// Java: SourceFile.CodeLoader.OnDisk(Path, Charset)
func NewSourceFile_CodeLoader_OnDisk(relativePath *jutil.Path, c *jutil.Charset) *SourceFile_CodeLoader_OnDisk {
	return &SourceFile_CodeLoader_OnDisk{serializableCharset: c.Name(), relativePathStr: relativePath.String()}
}

func (l *SourceFile_CodeLoader_OnDisk) loadCode() (string, jutil.IIOException) {
	b, err := os.ReadFile(l.relativePathStr)
	if err != nil {
		return "", jutil.IOExceptionFromError(err, l.relativePathStr).(jutil.IIOException)
	}
	return sourceFile_decodeStrictBytes(l.GetCharset(), b, l.relativePathStr)
}

func sourceFile_decodeStrictBytes(cs *jutil.Charset, b []byte, path string) (s string, ioErr jutil.IIOException) {
	if (cs == nil || cs == jutil.StandardCharsets_UTF_8) && utf8.Valid(b) {
		return string(b), nil
	}
	defer func() {
		if r := recover(); r != nil {
			if ce, ok := r.(jutil.ICharacterCodingException); ok {
				msg := "Failed to read: " + path + ", is this input UTF-8 encoded?"
				ex := &jutil.IOException{}
				jutil.InitThrowable(ex, "java.io.IOException", &msg, ce)
				ioErr = ex
				return
			}
			panic(r)
		}
	}()
	return cs.DecodeStrict(b), nil
}

// Java: SourceFile.CodeLoader.OnDisk.loadUncachedCode()
func (l *SourceFile_CodeLoader_OnDisk) LoadUncachedCode() string {
	s, err := l.loadCode()
	if err != nil {
		panic(err)
	}
	return s
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
	zipName             string
	entryName           string
	serializableCharset string
}

// Java: SourceFile.CodeLoader.AtZip(String, String, Charset)
func NewSourceFile_CodeLoader_AtZip(zipName string, entryName string, c *jutil.Charset) *SourceFile_CodeLoader_AtZip {
	return &SourceFile_CodeLoader_AtZip{zipName: zipName, entryName: entryName, serializableCharset: c.Name()}
}

func (l *SourceFile_CodeLoader_AtZip) loadCode() (string, jutil.IIOException) {
	rc, ioErr := jSCompZipFileCache_getEntryStreamErr(l.zipName, l.entryName)
	if ioErr != nil {
		return "", ioErr
	}
	data, err := io.ReadAll(rc)
	_ = rc.Close()
	if err != nil {
		return "", jutil.IOExceptionFromError(err, "").(jutil.IIOException)
	}
	return jutil.CharStreams_ToString(jutil.NewInputStreamReader(bytes.NewReader(data), l.GetCharset())), nil
}

// Java: SourceFile.CodeLoader.AtZip.loadUncachedCode()
func (l *SourceFile_CodeLoader_AtZip) LoadUncachedCode() string {
	s, err := l.loadCode()
	if err != nil {
		panic(err)
	}
	return s
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
