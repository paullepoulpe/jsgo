// Port of test/com/google/javascript/jscomp/SourceFileTest.java @ closure-compiler v20260415

/*
 * Copyright 2007 The Closure Compiler Authors.
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
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestSourceFile_SourceKind(t *testing.T) {
	sf1 := SourceFile_FromCode("test1.js", "1")
	jsaTrue(t, sf1.IsStrong())
	jsaFalse(t, sf1.IsWeak())
	jsaFalse(t, sf1.IsExtern())

	sf1.SetKind(rhino.StaticSourceFile_SourceKind_WEAK)
	jsaFalse(t, sf1.IsStrong())
	jsaTrue(t, sf1.IsWeak())
	jsaFalse(t, sf1.IsExtern())

	sf2 := SourceFile_FromCodeWithKind("test2.js", "2", rhino.StaticSourceFile_SourceKind_WEAK)
	jsaFalse(t, sf2.IsStrong())
	jsaTrue(t, sf2.IsWeak())
	jsaFalse(t, sf2.IsExtern())

	sf2.SetKind(rhino.StaticSourceFile_SourceKind_EXTERN)
	jsaFalse(t, sf2.IsStrong())
	jsaFalse(t, sf2.IsWeak())
	jsaTrue(t, sf2.IsExtern())
}

func TestSourceFile_LineOffset(t *testing.T) {
	factory := func(code string) *SourceFile { return SourceFile_FromCode("test.js", code) }
	f0 := factory("")
	jsaEq(t, f0.GetLineOfOffset(0), 1)
	jsaEq(t, f0.GetColumnOfOffset(0), 0)
	jsaEq(t, f0.GetLineOfOffset(10), 1)
	jsaEq(t, f0.GetColumnOfOffset(10), 10)
	jsaEq(t, f0.GetNumBytes(), 0)
	jsaEq(t, f0.GetNumLines(), 1)

	f1 := factory("'1';\n'2';\n'3'\n")
	jsaEq(t, f1.GetLineOffset(1), 0)
	jsaEq(t, f1.GetLineOffset(2), 5)
	jsaEq(t, f1.GetLineOffset(3), 10)
	jsaEq(t, f1.GetNumBytes(), 14)
	jsaEq(t, f1.GetNumLines(), 4)

	f2 := factory("'100';\n'200;'\n'300'\n")
	jsaEq(t, f2.GetLineOffset(1), 0)
	jsaEq(t, f2.GetLineOffset(2), 7)
	jsaEq(t, f2.GetLineOffset(3), 14)
	jsaEq(t, f2.GetNumBytes(), 20)
	jsaEq(t, f2.GetNumLines(), 4)

	longLine := strings.Repeat("a", 300)
	f3 := factory(longLine + "\n" + longLine + "\n" + longLine + "\n")
	jsaEq(t, f3.GetLineOffset(1), 0)
	jsaEq(t, f3.GetLineOffset(2), 301)
	jsaEq(t, f3.GetLineOffset(3), 602)

	jsaEq(t, f3.GetLineOfOffset(0), 1)
	jsaEq(t, f3.GetLineOfOffset(300), 1)
	jsaEq(t, f3.GetLineOfOffset(301), 2)
	jsaEq(t, f3.GetLineOfOffset(601), 2)
	jsaEq(t, f3.GetLineOfOffset(602), 3)
	jsaEq(t, f3.GetLineOfOffset(902), 3)
	jsaEq(t, f3.GetLineOfOffset(903), 4)

	jsaEq(t, f3.GetNumBytes(), 903)
	jsaEq(t, f3.GetNumLines(), 4)

	// TODO(nickreid): This seems like a bug.
	jsaEq(t, f3.GetLineOfOffset(-1), 0)
	jsaThrows[jutil.Throwable](t, func() { f3.GetColumnOfOffset(-1) })

	startsWithNewline := factory("\n'a'\n'b'")
	jsaEq(t, startsWithNewline.GetLineOffset(1), 0)
	jsaEq(t, startsWithNewline.GetLineOffset(2), 1)
	jsaEq(t, startsWithNewline.GetLineOffset(3), 5)
	jsaEq(t, startsWithNewline.GetNumBytes(), 8)
	jsaEq(t, startsWithNewline.GetNumLines(), 3)
}

// Go-only: extra coverage — offsets and "byte" counts are UTF-16 units, as
// in Java.
func TestSourceFile_NonAsciiOffsets(t *testing.T) {
	f := SourceFile_FromCode("test.js", "'\u00e9\U0001F600';\nx\n")
	jsaEq(t, f.GetLineOffset(2), 7)
	jsaEq(t, f.GetNumBytes(), 9)
	jsaEq(t, f.GetLineOfOffset(6), 1)
	jsaEq(t, f.GetLineOfOffset(7), 2)
	jsaEq(t, *f.GetLine(1), "'\u00e9\U0001F600';")
	jsaEq(t, *f.GetLine(2), "x")
}

func jsaWriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestSourceFile_CachingFile(t *testing.T) {
	expectedContent := "// content content content"
	newExpectedContent := "// new content new content new content"

	jsPath := filepath.Join(t.TempDir(), "test.js")
	jsaWriteFile(t, jsPath, expectedContent)
	sourceFile := SourceFile_FromPath(jutil.Path_Of(jsPath), jutil.StandardCharsets_UTF_8)

	jsaEq(t, sourceFile.GetCode(), expectedContent)

	jsaWriteFile(t, jsPath, newExpectedContent)
	sourceFile.ClearCachedSource()

	jsaEq(t, sourceFile.GetCode(), newExpectedContent)
}

func jsaCreateZipWithContent(t *testing.T, zipFile string, content string) {
	t.Helper()
	lastModified := time.Now()
	if st, err := os.Stat(zipFile); err == nil {
		// Ensure that file modified date is updated, otherwise could cause flakiness (b/123962282).
		lastModified = st.ModTime().Add(time.Second)
		os.Remove(zipFile)
	}
	f, err := os.Create(zipFile)
	if err != nil {
		t.Fatal(err)
	}
	zos := zip.NewWriter(f)
	w, err := zos.Create("foo.js")
	if err != nil {
		t.Fatal(err)
	}
	w.Write([]byte(content))
	if err := zos.Close(); err != nil {
		t.Fatal(err)
	}
	f.Close()
	if err := os.Chtimes(zipFile, lastModified, lastModified); err != nil {
		t.Fatal(err)
	}
}

func TestSourceFile_CachingZipFile(t *testing.T) {
	expectedContent := "// content content content"
	newExpectedContent := "// new content new content new content"
	jsZipFile := filepath.Join(t.TempDir(), "test.js.zip")
	jsaCreateZipWithContent(t, jsZipFile, expectedContent)
	zipSourceFile := SourceFile_Builder_().
		WithZipEntryPath(jutil.Path_Of(jsZipFile).ToAbsolutePath().String(), "foo.js").
		Build()

	jsaEq(t, zipSourceFile.GetCode(), expectedContent)

	jsaCreateZipWithContent(t, jsZipFile, newExpectedContent)
	// Verify cache is consistent unless cleared.
	jsaEq(t, zipSourceFile.GetCode(), expectedContent)
	zipSourceFile.ClearCachedSource()

	jsaEq(t, zipSourceFile.GetCode(), newExpectedContent)
}

func TestSourceFile_SourceFileResolvesZipEntries(t *testing.T) {
	expectedContent := "// <program goes here>"
	jsZipPath := filepath.Join(t.TempDir(), "test.js.zip")
	jsaCreateZipWithContent(t, jsZipPath, expectedContent)
	abs := jutil.Path_Of(jsZipPath).ToAbsolutePath().String()

	sourceFileFromZipEntry := SourceFile_Builder_().
		WithKind(rhino.StaticSourceFile_SourceKind_WEAK).
		WithZipEntryPath(abs, "foo.js").
		Build()
	jsaEq(t, sourceFileFromZipEntry.GetCode(), expectedContent)
	jsaEq(t, sourceFileFromZipEntry.GetKind(), rhino.StaticSourceFile_SourceKind_WEAK)

	sourceFileFromZipEntryDefaultKind := SourceFile_Builder_().
		WithZipEntryPath(abs, "foo.js").
		Build()
	jsaEq(t, sourceFileFromZipEntryDefaultKind.GetCode(), expectedContent)
	jsaEq(t, sourceFileFromZipEntryDefaultKind.GetKind(), rhino.StaticSourceFile_SourceKind_STRONG)

	sourceFileFromFileString := SourceFile_FromFileWithCharset(jsZipPath+"!/foo.js", jutil.StandardCharsets_UTF_8)
	jsaEq(t, sourceFileFromFileString.GetCode(), expectedContent)

	sourceFileFromFileStringCharset := SourceFile_FromFileWithCharset(jsZipPath+"!/foo.js", jutil.StandardCharsets_UTF_8)
	jsaEq(t, sourceFileFromFileStringCharset.GetCode(), expectedContent)

	zipEntryPath := jutil.Path_Of(jsZipPath + "!/foo.js")
	sourceFileFromPathCharset := SourceFile_FromPath(zipEntryPath, jutil.StandardCharsets_UTF_8)
	jsaEq(t, sourceFileFromPathCharset.GetCode(), expectedContent)
}

func TestSourceFile_SourceFileFromZipFile(t *testing.T) {
	expectedContent := "// <program goes here>"
	jsZipPath := filepath.Join(t.TempDir(), "test.js.zip")
	jsaCreateZipWithContent(t, jsZipPath, "// <program goes here>")

	sourceFiles := SourceFile_FromZipFile(jsZipPath, jutil.StandardCharsets_UTF_8)
	jsaEq(t, len(sourceFiles), 1)

	sourceFile := sourceFiles[0]
	jsaEq(t, sourceFile.GetName(), jsZipPath+"!/foo.js")
	jsaEq(t, sourceFile.GetCode(), expectedContent)
}

func TestSourceFile_SourceFileFromZipInput(t *testing.T) {
	expectedContent := "// <program goes here>"
	jsZipPath := filepath.Join(t.TempDir(), "test.js.zip")
	jsaCreateZipWithContent(t, jsZipPath, "// <program goes here>")

	in, err := os.Open(jsZipPath)
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()
	sourceFiles := SourceFile_FromZipInput(jsZipPath, in, jutil.StandardCharsets_UTF_8)
	jsaEq(t, len(sourceFiles), 1)

	sourceFile := sourceFiles[0]
	jsaEq(t, sourceFile.GetName(), jsZipPath+"!/foo.js")
	jsaEq(t, sourceFile.GetCode(), expectedContent)
}

// Go-only: extra coverage — a missing zip throws FileNotFoundException.
func TestSourceFile_FromZipFileMissing(t *testing.T) {
	e := jsaThrows[jutil.IIOException](t, func() { SourceFile_FromZipFile("/nonexistent/x.zip", jutil.StandardCharsets_UTF_8) })
	if e != nil {
		jsaEq(t, e.Error(), "java.io.FileNotFoundException: /nonexistent/x.zip (No such file or directory)")
	}
}

func TestSourceFile_StubSourceFile_created_crashesOnGetCode(t *testing.T) {
	fileName := "file.js"
	stubSourceFile := SourceFile_StubSourceFile(fileName, rhino.StaticSourceFile_SourceKind_WEAK)
	jsaEq(t, stubSourceFile.GetName(), fileName)
	jsaEq(t, stubSourceFile.GetKind(), rhino.StaticSourceFile_SourceKind_WEAK)
	jsaThrows[jutil.IUnsupportedOperationException](t, func() { stubSourceFile.GetCode() })
}

func TestSourceFile_DiskFile(t *testing.T) {
	expectedContent := "var c;"
	tempFile := filepath.Join(t.TempDir(), "test.js")
	jsaWriteFile(t, tempFile, expectedContent)

	newFile := SourceFile_FromFile(tempFile)
	jsaEq(t, *newFile.GetLine(1), expectedContent)

	newFile.ClearCachedSource()

	jsaFalse(t, newFile.HasSourceInMemory())
	jsaEq(t, *newFile.GetLine(1), expectedContent)
}

func TestSourceFile_DiskFileWithOriginalPath(t *testing.T) {
	expectedContent := "var c;"
	tempFile := filepath.Join(t.TempDir(), "test.js")
	jsaWriteFile(t, tempFile, expectedContent)

	newFile := SourceFile_Builder_().
		WithOriginalPath("original_test.js").
		WithPath(tempFile).
		Build()
	jsaEq(t, *newFile.GetLine(1), expectedContent)

	newFile.ClearCachedSource()

	jsaFalse(t, newFile.HasSourceInMemory())
	jsaEq(t, *newFile.GetLine(1), expectedContent)

	jsaEq(t, newFile.GetName(), "original_test.js")
}

func TestSourceFile_GetLines(t *testing.T) {
	sourceFile := SourceFile_FromCode("file.js", "const a = 0;\nconst b = 1;\nconst c = 2;")

	jsaEq(t, sourceFile.GetLines(1, 1).GetSourceExcerpt(), "const a = 0;")
	jsaEq(t, sourceFile.GetLines(2, 1).GetSourceExcerpt(), "const b = 1;")
	jsaEq(t, sourceFile.GetLines(3, 1).GetSourceExcerpt(), "const c = 2;")

	jsaEq(t, sourceFile.GetLines(1, len("const a = 0;\n")).GetSourceExcerpt(), "const a = 0;")
	jsaEq(t, sourceFile.GetLines(1, len("const a = 0;\nconst b")).GetSourceExcerpt(), "const a = 0;\nconst b = 1;")
	jsaEq(t, sourceFile.GetLines(2, len("const b = 1;\nconst c")).GetSourceExcerpt(), "const b = 1;\nconst c = 2;")

	jsaEq(t, sourceFile.GetLines(3, len("const c = 2;")).GetSourceExcerpt(), "const c = 2;")
}

func TestSourceFile_GetLines_invalidLengths(t *testing.T) {
	sourceFile := SourceFile_FromCode("file.js", "const a = 0;\nconst b = 1;\nconst c = 2;")

	jsaEq(t, sourceFile.GetLines(0, -3).GetSourceExcerpt(), "const a = 0;")
	jsaEq(t, sourceFile.GetLines(0, 0).GetSourceExcerpt(), "const a = 0;")
	jsaEq(t, sourceFile.GetLines(1, 100000).GetSourceExcerpt(), "const a = 0;\nconst b = 1;\nconst c = 2;")
	jsaEq(t, sourceFile.GetLines(2, 10000).GetSourceExcerpt(), "const b = 1;\nconst c = 2;")
	jsaEq(t, sourceFile.GetLines(3, 10000).GetSourceExcerpt(), "const c = 2;")
}

func TestSourceFile_GetLines_whenFileEndsWithNewline(t *testing.T) {
	sourceFile := SourceFile_FromCode("file.js", "const a = 0;\nconst b = 1;\nconst c = 2;\n")

	jsaEq(t, sourceFile.GetLines(1, 100000).GetSourceExcerpt(), "const a = 0;\nconst b = 1;\nconst c = 2;")
	jsaEq(t, sourceFile.GetLines(2, 10000).GetSourceExcerpt(), "const b = 1;\nconst c = 2;")
	jsaEq(t, sourceFile.GetLines(3, 1).GetSourceExcerpt(), "const c = 2;")
	jsaEq(t, sourceFile.GetLines(4, 1).GetSourceExcerpt(), "")
}

func TestSourceFile_GetLines_invalidLineNumbers(t *testing.T) {
	sourceFile := SourceFile_FromCode("file.js", "const a = 0;\nconst b = 1;\nconst c = 2;")

	jsaEq(t, sourceFile.GetLines(-20, 1).GetSourceExcerpt(), "const a = 0;")
	jsaEq(t, sourceFile.GetLines(0, 1).GetSourceExcerpt(), "const a = 0;")
	jsaTrue(t, sourceFile.GetLines(4, 1) == nil)
}

func TestSourceFile_FromProtoPreloadedContents(t *testing.T) {
	sourceFile := SourceFile_FromProto(
		SourceFileProto_NewBuilder().
			SetFilename("file.js").
			SetSourceKind(SourceFileProto_SourceKind_CODE).
			SetPreloadedContents("42").
			Build())

	jsaEq(t, sourceFile.GetName(), "file.js")
	jsaEq(t, sourceFile.GetCode(), "42")
	jsaEq(t, sourceFile.GetNumLines(), 1)
	jsaEq(t, sourceFile.GetNumBytes(), 2)
}

func TestSourceFile_FromProto_getNumLinesAvoidsFileRead(t *testing.T) {
	// Create a SourceFile pointing to a path that doesn't actually exist.
	sourceFile := SourceFile_FromProto(
		SourceFileProto_NewBuilder().
			SetFilename("file.js").
			SetSourceKind(SourceFileProto_SourceKind_CODE).
			SetFileOnDisk(SourceFileProto_FileOnDisk_GetDefaultInstance()).
			SetNumLinesPlusOne(2).
			SetNumBytesPlusOne(3).
			Build())

	jsaEq(t, sourceFile.GetNumLines(), 1)
	jsaEq(t, sourceFile.GetNumBytes(), 2)

	// Verify reading the file from disk does fail.
	jsaThrows[jutil.IIOException](t, func() { sourceFile.GetCode() })
}

func TestSourceFile_FromProtoPreloadedContents_withNumLinesAndBytes(t *testing.T) {
	sourceFile := SourceFile_FromProto(
		SourceFileProto_NewBuilder().
			SetFilename("file.js").
			SetSourceKind(SourceFileProto_SourceKind_CODE).
			SetPreloadedContents("42").
			SetNumLinesPlusOne(2).
			SetNumBytesPlusOne(3).
			Build())

	jsaEq(t, sourceFile.GetNumLines(), 1)
	jsaEq(t, sourceFile.GetNumBytes(), 2)
}

func TestSourceFile_GetProto_withoutGetCodeCalls(t *testing.T) {
	sourceFile := SourceFile_FromCode("file.js", "42")

	jsaJavaEq(t, sourceFile.GetProto(),
		SourceFileProto_NewBuilder().
			SetFilename("file.js").
			SetPreloadedContents("42").
			SetSourceKind(SourceFileProto_SourceKind_CODE).
			Build())
}

func TestSourceFile_GetProto_withGetCodeCalls(t *testing.T) {
	sourceFile := SourceFile_FromCode("file.js", "42")
	_ = sourceFile.GetCode()

	jsaJavaEq(t, sourceFile.GetProto(),
		SourceFileProto_NewBuilder().
			SetFilename("file.js").
			SetPreloadedContents("42").
			// note that numLines and numBytes are only present after a 'getCode' call.
			SetNumLinesPlusOne(2).
			SetNumBytesPlusOne(3).
			SetSourceKind(SourceFileProto_SourceKind_CODE).
			Build())
}

func TestSourceFile_GetProto_recordsClosureUnawareCodeMark(t *testing.T) {
	sourceFile := SourceFile_FromCode("file.js", "42")
	sourceFile.MarkAsClosureUnawareCode()

	jsaJavaEq(t, sourceFile.GetProto(),
		SourceFileProto_NewBuilder().
			SetFilename("file.js").
			SetPreloadedContents("42").
			SetIsClosureUnawareCode(true).
			SetSourceKind(SourceFileProto_SourceKind_CODE).
			Build())
}

func TestSourceFile_GetProto_withGetCodeCalls_andClearCachedSource(t *testing.T) {
	sourceFile := SourceFile_FromCode("file.js", "42")
	_ = sourceFile.GetCode()
	sourceFile.ClearCachedSource()

	jsaJavaEq(t, sourceFile.GetProto(),
		SourceFileProto_NewBuilder().
			SetFilename("file.js").
			SetPreloadedContents("42").
			// note that numLines and numBytes are still present, i.e. were not cleared by
			// .clearCachedSource()
			SetNumLinesPlusOne(2).
			SetNumBytesPlusOne(3).
			SetSourceKind(SourceFileProto_SourceKind_CODE).
			Build())
}

func TestSourceFile_RestoreCachedState_correctNumLines(t *testing.T) {
	sourceFile := SourceFile_FromCode("file.js", "42")

	// In the proto, numLines and numBytes are increased by 1 from the actual value.
	sourceFile.RestoreCachedStateFrom(
		SourceFileProto_NewBuilder().
			SetFilename("file.js").
			SetNumLinesPlusOne(2).
			SetNumBytesPlusOne(3).
			Build())

	jsaEq(t, sourceFile.GetNumLines(), 1)
	jsaEq(t, sourceFile.GetNumBytes(), 2)
}

func TestSourceFile_RestoreCachedState_setsClosureUnawareCodeTrue(t *testing.T) {
	sourceFile := SourceFile_FromCode("file.js", "42")
	jsaFalse(t, sourceFile.IsClosureUnawareCode())

	sourceFile.RestoreCachedStateFrom(
		SourceFileProto_NewBuilder().SetFilename("file.js").SetIsClosureUnawareCode(true).Build())

	jsaTrue(t, sourceFile.IsClosureUnawareCode())
}

func TestSourceFile_RestoreCachedState_setsClosureUnawareCodeFalse_protoIsSourceOfTruth(t *testing.T) {
	sourceFile := SourceFile_FromCode("file.js", "42")
	sourceFile.MarkAsClosureUnawareCode()
	jsaTrue(t, sourceFile.IsClosureUnawareCode())

	sourceFile.RestoreCachedStateFrom(
		SourceFileProto_NewBuilder().SetFilename("file.js").SetIsClosureUnawareCode(false).Build())

	jsaFalse(t, sourceFile.IsClosureUnawareCode())
}

func TestSourceFile_RestoreCachedState_forcesRecalculationIfValuesUnsetInProto(t *testing.T) {
	sourceFile := SourceFile_FromCode("file.js", "42")

	sourceFile.RestoreCachedStateFrom(
		SourceFileProto_NewBuilder().
			SetFilename("file.js").
			// avoid setting numLines and numBytes, so that they default to '0'
			Build())

	jsaEq(t, sourceFile.GetNumLines(), 1)
	jsaEq(t, sourceFile.GetNumBytes(), 2)
}

func TestSourceFile_RestoreCachedState_treatsProtoNumLinesAsCanonical(t *testing.T) {
	sourceFile := SourceFile_FromCode("file.js", "42")

	jsaEq(t, sourceFile.GetNumLines(), 1)

	sourceFile.RestoreCachedStateFrom(
		SourceFileProto_NewBuilder().
			SetFilename("file.js").
			SetNumLinesPlusOne(1000).
			SetNumBytesPlusOne(-1).
			Build())

	jsaEq(t, sourceFile.GetNumLines(), 999)
}
