// Port of test/com/google/javascript/jscomp/VerboseMessageFormatterTest.java @ closure-compiler v20261021

/*
 * Copyright 2026 The Closure Compiler Authors.
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
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

// VerboseMessageFormatterTest is the JUnit 4 class VerboseMessageFormatterTest.
type VerboseMessageFormatterTest struct {
	t *testing.T
}

// Java: VerboseMessageFormatterTest.FOO_TYPE
var VerboseMessageFormatterTest_FOO_TYPE = DiagnosticType_Error("TEST_FOO", "error description here")

func (c *VerboseMessageFormatterTest) JUnitSetT(t *testing.T) { c.t = t }
func (c *VerboseMessageFormatterTest) JUnitBefore()           {}
func (c *VerboseMessageFormatterTest) JUnitAfter()            {}

func TestVerboseMessageFormatterTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return &VerboseMessageFormatterTest{} })
}

// Java: VerboseMessageFormatterTest.testFormatError_validLineNumber()
func (c *VerboseMessageFormatterTest) TestFormatError_validLineNumber() {
	formatter := NewVerboseMessageFormatter(verboseMessageFormatterTest_emptySourceProvider())
	error := JSError_MakeWithSourceLocation(jutil.StringPtr("test.js"), 10, 5, VerboseMessageFormatterTest_FOO_TYPE)

	formatted := formatter.FormatError(error)

	AssertThat(formatted).Contains("line 10")
	AssertThat(formatted).DoesNotContain("(unknown line)")
	AssertThat(formatted).IsEqualTo("ERROR: error description here at test.js line 10 .")
}

// Java: VerboseMessageFormatterTest.testFormatError_unknownLineNumber()
func (c *VerboseMessageFormatterTest) TestFormatError_unknownLineNumber() {
	formatter := NewVerboseMessageFormatter(verboseMessageFormatterTest_emptySourceProvider())
	error := JSError_MakeWithSourceLocation(jutil.StringPtr("test.js"), -1, -1, VerboseMessageFormatterTest_FOO_TYPE)

	formatted := formatter.FormatError(error)

	AssertThat(formatted).Contains("line (unknown line)")
	AssertThat(formatted).DoesNotContain("-1")
	AssertThat(formatted).
		IsEqualTo("ERROR: error description here at test.js line (unknown line) .")
}

// Java: VerboseMessageFormatterTest.testFormatError_unknownSourceName()
func (c *VerboseMessageFormatterTest) TestFormatError_unknownSourceName() {
	formatter := NewVerboseMessageFormatter(verboseMessageFormatterTest_emptySourceProvider())
	error := JSError_MakeWithSourceLocation(nil, 10, 5, VerboseMessageFormatterTest_FOO_TYPE)

	formatted := formatter.FormatError(error)

	AssertThat(formatted).Contains("(unknown source)")
	AssertThat(formatted).Contains("line 10")
	AssertThat(formatted).IsEqualTo("ERROR: error description here at (unknown source) line 10 .")
}

// Java: VerboseMessageFormatterTest.testFormatError_unknownSourceAndLineNumber()
func (c *VerboseMessageFormatterTest) TestFormatError_unknownSourceAndLineNumber() {
	formatter := NewVerboseMessageFormatter(verboseMessageFormatterTest_emptySourceProvider())
	error := JSError_MakeWithSourceLocation(nil, -1, -1, VerboseMessageFormatterTest_FOO_TYPE)

	formatted := formatter.FormatError(error)

	AssertThat(formatted).Contains("(unknown source)")
	AssertThat(formatted).Contains("line (unknown line)")
	AssertThat(formatted).
		IsEqualTo("ERROR: error description here at (unknown source) line (unknown line) .")
}

// verboseMessageFormatterTest_regionProvider is the anonymous
// SourceExcerptProvider of testFormatError_withSourceRegion.
type verboseMessageFormatterTest_regionProvider struct {
	verboseMessageFormatterTest_emptyProvider
}

// Java: VerboseMessageFormatterTest.testFormatError_withSourceRegion().new SourceExcerptProvider() {...}.getSourceRegion(String, int)
func (p *verboseMessageFormatterTest_regionProvider) GetSourceRegion(sourceName *string, lineNumber int) Region {
	// Objects.equals(sourceName, "test.js")
	if sourceName != nil && *sourceName == "test.js" && lineNumber == 10 {
		return NewSimpleRegion(10, 10, "const x = 1;")
	}
	return nil
}

// Java: VerboseMessageFormatterTest.testFormatError_withSourceRegion()
func (c *VerboseMessageFormatterTest) TestFormatError_withSourceRegion() {
	provider := &verboseMessageFormatterTest_regionProvider{}

	formatter := NewVerboseMessageFormatter(provider)
	error := JSError_MakeWithSourceLocation(jutil.StringPtr("test.js"), 10, 0, VerboseMessageFormatterTest_FOO_TYPE)

	formatted := formatter.FormatError(error)

	AssertThat(formatted).
		IsEqualTo("ERROR: error description here at test.js line 10 :\n\nconst x = 1;")
}

// Java: VerboseMessageFormatterTest.testFormatWarning_validLineNumber()
func (c *VerboseMessageFormatterTest) TestFormatWarning_validLineNumber() {
	formatter := NewVerboseMessageFormatter(verboseMessageFormatterTest_emptySourceProvider())
	warningType := DiagnosticType_Warning("TEST_WARN", "warning description here")
	warning := JSError_MakeWithSourceLocation(jutil.StringPtr("test.js"), 42, 0, warningType)

	formatted := formatter.FormatWarning(warning)

	AssertThat(formatted).StartsWith("WARNING: ")
	AssertThat(formatted).Contains("line 42")
	AssertThat(formatted).IsEqualTo("WARNING: warning description here at test.js line 42 .")
}

// verboseMessageFormatterTest_emptyProvider is the anonymous
// SourceExcerptProvider of emptySourceProvider().
type verboseMessageFormatterTest_emptyProvider struct{}

func (p *verboseMessageFormatterTest_emptyProvider) GetSourceLine(sourceName *string, lineNumber int) *string {
	return nil
}

func (p *verboseMessageFormatterTest_emptyProvider) GetSourceLines(sourceName *string, lineNumber int, length int) Region {
	return nil
}

func (p *verboseMessageFormatterTest_emptyProvider) GetSourceRegion(sourceName *string, lineNumber int) Region {
	return nil
}

func (p *verboseMessageFormatterTest_emptyProvider) GetSourceMapping(
	sourceName *string, lineNumber int, columnNumber int) *sourcemap.Mapping_OriginalMapping {
	return nil
}

// Java: VerboseMessageFormatterTest.emptySourceProvider()
func verboseMessageFormatterTest_emptySourceProvider() SourceExcerptProvider {
	return &verboseMessageFormatterTest_emptyProvider{}
}
