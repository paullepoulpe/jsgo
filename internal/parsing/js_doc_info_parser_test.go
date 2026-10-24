// Port of com/google/javascript/jscomp/parsing/JsDocInfoParserTest.java @ closure-compiler v20261015

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

package parsing

// The Java test class extends rhino.testing.BaseJSTypeTestCase: it owns a
// JSTypeRegistry and compares parsed type expressions with JSTypes built in
// that registry (assertTypeEquals, testParseType, resolve, createNamedType,
// ...). Here jsDocInfoParserTest embeds rhinotesting.BaseJSTypeTestCase.

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

const MISSING_TYPE_DECL_WARNING_TEXT = "Missing type declaration."

var jsDocInfoParserTest_EMPTY_SCOPE = rhinotesting.MapBasedScope_EmptyScope()

// jsDocInfoParserTest holds the per-test state of JsDocInfoParserTest (its
// fields, set up by setUp()) on top of BaseJSTypeTestCase.
type jsDocInfoParserTest struct {
	*rhinotesting.BaseJSTypeTestCase
	t *testing.T

	extraAnnotations        *jutil.LinkedHashSet[string]
	extraSuppressions       *jutil.LinkedHashSet[string]
	extraPrimitives         *jutil.LinkedHashSet[string]
	allowClosureUnawareCode bool
	prevLicense             *string
	jsDocSourceKind         JsDocInfoParser_JsDocSourceKind
}

// newJsDocInfoParserTest is setUp().
func newJsDocInfoParserTest(t *testing.T) *jsDocInfoParserTest {
	jt := &jsDocInfoParserTest{BaseJSTypeTestCase: rhinotesting.NewBaseJSTypeTestCase(t), t: t}
	// Go-only: language-adapter — new HashSet<>(...): the sets are only
	// queried and passed to Config.Builder, so insertion order stands in for
	// HashSet order.
	jt.extraAnnotations = jutil.NewLinkedHashSet(
		ParserRunner_CreateConfig(Config_LanguageMode_ECMASCRIPT3, nil, Config_StrictMode_SLOPPY).
			AnnotationNames().ToSlice()...)
	jt.extraSuppressions = jutil.NewLinkedHashSet(
		ParserRunner_CreateConfig(Config_LanguageMode_ECMASCRIPT3, nil, Config_StrictMode_SLOPPY).
			SuppressionNames().ToSlice()...)

	jt.extraSuppressions.Add("x")
	jt.extraSuppressions.Add("y")
	jt.extraSuppressions.Add("z")

	jt.extraPrimitives = jutil.NewLinkedHashSet[string]()
	jt.extraPrimitives.Add("id")
	jt.extraPrimitives.Add("idA")
	jt.extraPrimitives.Add("idB")

	jt.jsDocSourceKind = JsDocInfoParser_JsDocSourceKind_NORMAL
	return jt
}

// Java: JsDocInfoParserTest.assertTemplatizedTypeEquals(TemplateType, JSType, JSTypeExpression)
func (jt *jsDocInfoParserTest) assertTemplatizedTypeEquals(key *rhino.TemplateType, expected rhino.JSType, te *rhino.JSTypeExpression) {
	jt.t.Helper()
	assertThat(jt.t, jt.Resolve(te).GetTemplateTypeMap().GetResolvedTemplateType(key)).IsEqualTo(expected)
}

// Java: JsDocInfoParserTest.createNamedType(String)
func (jt *jsDocInfoParserTest) createNamedType(name string) *rhino.NamedType {
	jt.ErrorReporter.ExpectAllWarnings("Bad type annotation. Unknown type " + name)
	sourceName := ""
	return jt.Registry.CreateNamedType(jsDocInfoParserTest_EMPTY_SCOPE, name, &sourceName, -1, -1)
}

// Java: JsDocInfoParserTest.testParseType(String) / testParseType(String, String)
func (jt *jsDocInfoParserTest) testParseType(type_ string, typeExpected ...string) rhino.JSType {
	jt.t.Helper()
	if len(typeExpected) == 0 {
		return jt.testParseType(type_, type_)
	}
	info := jt.parse("@type {" + type_ + "}*/")

	assertThat(jt.t, info).IsNotNull()
	assertThat(jt.t, info.HasType()).IsTrue()

	actual := jt.Resolve(info.GetType())
	assertThat(jt.t, actual.String()).IsEqualTo(typeExpected[0])
	return actual
}

// Java: JsDocInfoParserTest.checkTokenPosition(Node, Token, int, int)
func checkTokenPosition(t *testing.T, n rhino.Node, tok rhino.Token, lineno int, charno int) {
	t.Helper()
	assertNode(t, n).HasType(tok)
	assertNode(t, n).HasLineno(lineno)
	assertNode(t, n).HasCharno(charno)
}

func (jt *jsDocInfoParserTest) checkTokenPosition(n rhino.Node, tok rhino.Token, lineno int, charno int) {
	jt.t.Helper()
	checkTokenPosition(jt.t, n, tok, lineno, charno)
}

// assertDocumentationInMarker asserts that a documentation field exists on
// the given marker.
//
// Java: JsDocInfoParserTest.assertDocumentationInMarker(Marker, String, int, int, int)
func (jt *jsDocInfoParserTest) assertDocumentationInMarker(
	marker *rhino.JSDocInfo_Marker, description string, startCharno int, endLineno int, endCharno int,
) *rhino.JSDocInfo_Marker {
	t := jt.t
	t.Helper()
	assertThat(t, marker.GetDescription()).IsNotNull()
	assertThat(t, marker.GetDescription().GetItem()).IsEqualTo(description)

	// Match positional information.
	assertThat(t, marker.GetDescription().GetStartLine()).IsEqualTo(marker.GetAnnotation().GetStartLine())
	assertThat(t, marker.GetDescription().GetPositionOnStartLine()).IsEqualTo(startCharno)
	assertThat(t, marker.GetDescription().GetEndLine()).IsEqualTo(endLineno)
	assertThat(t, marker.GetDescription().GetPositionOnEndLine()).IsEqualTo(endCharno)

	return marker
}

// assertTypeInMarker asserts that a type field exists on the given marker.
//
// Java: JsDocInfoParserTest.assertTypeInMarker(Marker, String, int, int, int, int, boolean)
func (jt *jsDocInfoParserTest) assertTypeInMarker(
	marker *rhino.JSDocInfo_Marker,
	typeName string,
	startLineno int,
	startCharno int,
	endLineno int,
	endCharno int,
	hasBrackets bool,
) *rhino.JSDocInfo_Marker {
	t := jt.t
	t.Helper()
	assertThat(t, marker.GetType()).IsNotNull()
	assertThat(t, marker.GetType().GetItem().IsStringLit()).IsTrue()

	// Match the name and brackets information.
	foundName := marker.GetType().GetItem().GetString()

	assertThat(t, foundName).IsEqualTo(typeName)
	assertThat(t, marker.GetType().HasBrackets()).IsEqualTo(hasBrackets)

	// Match position information.
	assertThat(t, marker.GetType().GetPositionOnStartLine()).IsEqualTo(startCharno)
	assertThat(t, marker.GetType().GetPositionOnEndLine()).IsEqualTo(endCharno)
	assertThat(t, marker.GetType().GetStartLine()).IsEqualTo(startLineno)
	assertThat(t, marker.GetType().GetEndLine()).IsEqualTo(endLineno)

	return marker
}

// assertNameInMarker asserts that a name field exists on the given marker.
//
// Java: JsDocInfoParserTest.assertNameInMarker(Marker, String, int, int)
func (jt *jsDocInfoParserTest) assertNameInMarker(
	marker *rhino.JSDocInfo_Marker, name string, startLine int, startCharno int,
) *rhino.JSDocInfo_Marker {
	t := jt.t
	t.Helper()
	assertThat(t, marker.GetNameNode()).IsNotNull()
	assertThat(t, marker.GetNameNode().GetItem().GetString()).IsEqualTo(name)

	assertThat(t, marker.GetNameNode().GetPositionOnStartLine()).IsEqualTo(startCharno)
	assertThat(t, marker.GetNameNode().GetPositionOnEndLine()).IsEqualTo(startCharno + jutil.StringLength(name))

	assertThat(t, marker.GetNameNode().GetStartLine()).IsEqualTo(startLine)
	assertThat(t, marker.GetNameNode().GetEndLine()).IsEqualTo(startLine)

	return marker
}

// assertAnnotationMarker asserts that the index-th annotation marker of a
// given annotation name is found in the given JSDocInfo (index defaults to
// 0).
//
// Java: JsDocInfoParserTest.assertAnnotationMarker(JSDocInfo, String, int, int[, int])
func (jt *jsDocInfoParserTest) assertAnnotationMarker(
	jsdoc *rhino.JSDocInfo, annotationName string, startLineno int, startCharno int, index ...int,
) *rhino.JSDocInfo_Marker {
	t := jt.t
	t.Helper()
	idx := 0
	if len(index) > 0 {
		idx = index[0]
	}

	markers := jsdoc.GetMarkers()

	assertThat(t, markers).IsNotEmpty()

	counter := 0

	for _, marker := range markers {
		if marker.GetAnnotation() != nil {
			if annotationName == marker.GetAnnotation().GetItem() {

				if counter == idx {
					assertThat(t, marker.GetAnnotation().GetStartLine()).IsEqualTo(startLineno)
					assertThat(t, marker.GetAnnotation().GetPositionOnStartLine()).IsEqualTo(startCharno)
					assertThat(t, marker.GetAnnotation().GetEndLine()).IsEqualTo(startLineno)
					assertThat(t, marker.GetAnnotation().GetPositionOnEndLine()).
						IsEqualTo(startCharno + jutil.StringLength(annotationName))

					return marker
				}

				counter++
			}
		}
	}

	t.Fatalf("No marker found")
	return nil
}

// Java: JsDocInfoParserTest.parseFull(String, String...)
func (jt *jsDocInfoParserTest) parseFull(code string, warnings ...string) rhino.Node {
	jt.t.Helper()
	testErrorReporter := newTestErrorReporter().expectAllWarnings(warnings...)
	config := Config_Builder_().
		SetExtraAnnotationNames(jt.extraAnnotations.All()).
		SetJsDocParsingMode(Config_JsDocParsing_INCLUDE_DESCRIPTIONS_NO_WHITESPACE).
		SetRunMode(Config_RunMode_KEEP_GOING).
		SetSuppressionNames(jt.extraSuppressions.All()).
		SetLanguageMode(Config_LanguageMode_ECMASCRIPT3).
		SetParseInlineSourceMaps(true).
		SetStrictMode(Config_StrictMode_SLOPPY).
		Build()

	result := ParserRunner_Parse(
		rhino.NewSimpleSourceFile("source", rhino.StaticSourceFile_SourceKind_STRONG), code, config, testErrorReporter)

	testErrorReporter.verifyHasEncounteredAllWarningsAndErrors(jt.t)
	return result.Ast
}

// Java: JsDocInfoParserTest.parseFileOverviewWithoutDoc(String, String...)
func (jt *jsDocInfoParserTest) parseFileOverviewWithoutDoc(comment string, warnings ...string) *rhino.JSDocInfo {
	jt.t.Helper()
	return jt.parseWith(comment, Config_JsDocParsing_INCLUDE_DESCRIPTIONS_NO_WHITESPACE, false, warnings...)
}

// Java: JsDocInfoParserTest.parseFileOverview(String, String...)
func (jt *jsDocInfoParserTest) parseFileOverview(comment string, warnings ...string) *rhino.JSDocInfo {
	jt.t.Helper()
	return jt.parseWith(comment, Config_JsDocParsing_INCLUDE_DESCRIPTIONS_NO_WHITESPACE, true, warnings...)
}

// Java: JsDocInfoParserTest.preserveWhitespaceParse(String, String...)
func (jt *jsDocInfoParserTest) preserveWhitespaceParse(comment string, warnings ...string) *rhino.JSDocInfo {
	jt.t.Helper()
	return jt.parseWith(comment, Config_JsDocParsing_INCLUDE_DESCRIPTIONS_WITH_WHITESPACE, false, warnings...)
}

// parse dispatches like the Java overloads parse(String, String...),
// parse(String, boolean, String...) and parse(String, JsDocParsing,
// String...).
func (jt *jsDocInfoParserTest) parse(comment string, rest ...any) *rhino.JSDocInfo {
	jt.t.Helper()
	mode := Config_JsDocParsing_TYPES_ONLY
	if len(rest) > 0 {
		switch r := rest[0].(type) {
		case bool:
			if r {
				mode = Config_JsDocParsing_INCLUDE_DESCRIPTIONS_NO_WHITESPACE
			}
			rest = rest[1:]
		case Config_JsDocParsing:
			mode = r
			rest = rest[1:]
		}
	}
	warnings := make([]string, len(rest))
	for i, w := range rest {
		warnings[i] = w.(string)
	}
	return jt.parseWith(comment, mode, false, warnings...)
}

// Java: JsDocInfoParserTest.parse(String, JsDocParsing, boolean, String...)
func (jt *jsDocInfoParserTest) parseWith(
	comment string,
	parseDocumentation Config_JsDocParsing,
	parseFileOverview bool,
	warnings ...string,
) *rhino.JSDocInfo {
	t := jt.t
	t.Helper()
	errorReporter := newTestErrorReporter().expectAllWarnings(warnings...)
	if jt.allowClosureUnawareCode {
		errorReporter.expectAllWarnings(
			"@closureUnaware annotation is not allowed in this compilation")
	}

	config := Config_Builder_().
		SetExtraAnnotationNames(jt.extraAnnotations.All()).
		SetJsDocParsingMode(parseDocumentation).
		SetSuppressionNames(jt.extraSuppressions.All()).
		SetClosurePrimitiveNames(jt.extraPrimitives.All()).
		SetLanguageMode(Config_LanguageMode_ECMASCRIPT3).
		SetParseInlineSourceMaps(true).
		SetStrictMode(Config_StrictMode_SLOPPY).
		Build()

	file := rhino.NewSimpleSourceFile("testcode", rhino.StaticSourceFile_SourceKind_STRONG)
	templateNode := rhino.IR_Script()
	templateNode.SetStaticSourceFile(file)

	jsdocParser := NewJsDocInfoParser(
		jt.stream(comment), comment, 0, templateNode, config, jt.jsDocSourceKind, errorReporter)

	jsdocParser.Parse()
	jt.prevLicense = jsdocParser.GetLicenseText()

	errorReporter.verifyHasEncounteredAllWarningsAndErrors(t)

	var result *rhino.JSDocInfo
	if parseFileOverview {
		result = jsdocParser.GetFileOverviewJSDocInfo()
	} else {
		result = jsdocParser.RetrieveAndResetParsedJSDocInfo()
	}

	if result != nil {
		// assertThat(result.getTypeNodes()).comparingElementsUsing(
		//     transforming(NodeUtil::getSourceName, "has source name that")).doesNotContain(null)
		for _, n := range result.GetTypeNodes() {
			if testGetSourceName(n) == nil {
				t.Fatalf("type node %v has no source name", n)
			}
		}
	}

	return result
}

// testGetSourceName is NodeUtil.getSourceName(Node) (jscomp, not ported
// yet).
func testGetSourceName(n rhino.Node) *string {
	var sourceName *string
	for sourceName == nil && !n.IsNil() {
		sourceName = n.GetSourceFileName()
		n = n.GetParent()
	}
	return sourceName
}

// Java: JsDocInfoParserTest.parseType(String)
func (jt *jsDocInfoParserTest) parseType(typeComment string) rhino.Node {
	return JsDocInfoParser_ParseTypeString(typeComment)
}

// Java: JsDocInfoParserTest.stream(String)
func (jt *jsDocInfoParserTest) stream(source string) *JsDocTokenStream {
	return NewJsDocTokenStreamWithLineno(source, 0)
}

// ---------------------------------------------------------------------------
// TestErrorReporter (port of rhino/testing/TestErrorReporter.java): verifies
// that messages reported to the reporter are expected.

type testErrorReporter struct {
	expectedErrors   []string
	expectedWarnings []string

	// Seen diagnostics are stored in a LinkedHashMap to deduplicate identical
	// diagnostics. The default com.google.javascript.jscomp.SortingErrorManager
	// behaves in this way.
	seenErrors   *jutil.LinkedHashMap[string, string]
	seenWarnings *jutil.LinkedHashMap[string, string]
}

func newTestErrorReporter() *testErrorReporter {
	return &testErrorReporter{
		seenErrors:   jutil.NewLinkedHashMap[string, string](),
		seenWarnings: jutil.NewLinkedHashMap[string, string](),
	}
}

func (r *testErrorReporter) Error(message string, sourceName string, line int, lineOffset int) {
	r.seenErrors.Put(r.fmtDiagnosticKey(message, sourceName, line, lineOffset), message)
}

func (r *testErrorReporter) Warning(message string, sourceName string, line int, lineOffset int) {
	r.seenWarnings.Put(r.fmtDiagnosticKey(message, sourceName, line, lineOffset), message)
}

func (r *testErrorReporter) expectAllErrors(errors ...string) *testErrorReporter {
	r.expectedErrors = append(r.expectedErrors, errors...)
	return r
}

func (r *testErrorReporter) expectAllWarnings(warnings ...string) *testErrorReporter {
	r.expectedWarnings = append(r.expectedWarnings, warnings...)
	return r
}

func (r *testErrorReporter) verifyHasEncounteredAllWarningsAndErrors(t *testing.T) {
	t.Helper()
	assertThat(t, slicesOf(r.seenWarnings.Values())).ContainsExactlyElementsIn(r.expectedWarnings).InOrder()
	assertThat(t, slicesOf(r.seenErrors.Values())).ContainsExactlyElementsIn(r.expectedErrors).InOrder()
}

func (r *testErrorReporter) fmtDiagnosticKey(message string, sourceName string, line int, lineOffset int) string {
	return message + ":" + sourceName + ":" + fmt.Sprint(line) + ":" + fmt.Sprint(lineOffset)
}

func slicesOf(seq func(func(string) bool)) []string {
	out := []string{}
	for s := range seq {
		out = append(out, s)
	}
	return out
}

// setOf is ImmutableSet.of(...).
func setOf(elems ...string) *jutil.LinkedHashSet[string] {
	return jutil.NewLinkedHashSet(elems...)
}

// ---------------------------------------------------------------------------
// A small subset of Truth (and of rhino/testing/NodeSubject) with Java
// equality: values are compared with jutil.ObjectsEquals after normalizing
// a nullable Java String (*string) to its value and Go integer types to
// int64.

type subject struct {
	t      *testing.T
	actual any
	msg    string
}

func assertThat(t *testing.T, actual any) *subject {
	return &subject{t: t, actual: actual}
}

type messageSubjectFactory struct {
	t   *testing.T
	msg string
}

func assertWithMessage(t *testing.T, msg string) *messageSubjectFactory {
	return &messageSubjectFactory{t: t, msg: msg}
}

func (f *messageSubjectFactory) That(actual any) *subject {
	return &subject{t: f.t, actual: actual, msg: f.msg}
}

func (f *messageSubjectFactory) Fail() {
	f.t.Helper()
	f.t.Fatal(f.msg)
}

func truthNorm(v any) any {
	if jutil.IsNil(v) {
		return nil
	}
	if rv := reflect.ValueOf(v); rv.Kind() == reflect.Pointer && rv.Elem().Kind() == reflect.Int32 {
		// A nullable Java enum (*E).
		return rv.Elem().Interface()
	}
	switch x := v.(type) {
	case *string:
		return *x
	case int:
		return int64(x)
	case int32:
		return int64(x)
	case int16:
		return int64(x)
	case int8:
		return int64(x)
	case uint16:
		return int64(x)
	}
	return v
}

func truthEquals(a, b any) bool {
	return jutil.ObjectsEquals(truthNorm(a), truthNorm(b))
}

func (s *subject) fail(format string, args ...any) {
	s.t.Helper()
	m := fmt.Sprintf(format, args...)
	if s.msg != "" {
		m = s.msg + ": " + m
	}
	s.t.Fatal(m)
}

func truthString(v any) string {
	v = truthNorm(v)
	if s, ok := v.(string); ok {
		return fmt.Sprintf("%q", s)
	}
	return jutil.ToJavaString(v)
}

func (s *subject) IsTrue() {
	s.t.Helper()
	if b, ok := s.actual.(bool); !ok || !b {
		s.fail("expected to be true, was %v", s.actual)
	}
}

func (s *subject) IsFalse() {
	s.t.Helper()
	if b, ok := s.actual.(bool); !ok || b {
		s.fail("expected to be false, was %v", s.actual)
	}
}

func (s *subject) IsNull() {
	s.t.Helper()
	if !jutil.IsNil(s.actual) {
		s.fail("expected null, was %s", truthString(s.actual))
	}
}

func (s *subject) IsNotNull() {
	s.t.Helper()
	if jutil.IsNil(s.actual) {
		s.fail("expected not to be null")
	}
}

func (s *subject) IsEqualTo(expected any) {
	s.t.Helper()
	if !truthEquals(s.actual, expected) {
		s.fail("expected %s, was %s", truthString(expected), truthString(s.actual))
	}
}

func (s *subject) IsNotEqualTo(expected any) {
	s.t.Helper()
	if truthEquals(s.actual, expected) {
		s.fail("expected not to be %s", truthString(expected))
	}
}

// truthElements returns the elements of a Java collection (a Go slice, a
// jutil list or set, or anything with ToSlice()).
func truthElements(t *testing.T, v any) []any {
	t.Helper()
	rv := reflect.ValueOf(v)
	if jutil.IsNil(v) && rv.Kind() != reflect.Slice {
		t.Fatalf("expected a collection, was null")
	}
	if str, ok := truthNorm(v).(string); ok {
		// StringSubject: its UTF-16 chars.
		out := []any{}
		for _, c := range jutil.StringToUTF16(str) {
			out = append(out, c)
		}
		return out
	}
	if m := rv.MethodByName("ToSlice"); m.IsValid() {
		rv = m.Call(nil)[0]
	} else if m := rv.MethodByName("All"); m.IsValid() && m.Type().NumOut() == 1 && m.Type().Out(0).NumIn() == 1 && m.Type().Out(0).In(0).NumIn() == 2 {
		// MapSubject: its entries.
		keys, _ := truthEntries(t, v)
		return keys
	}
	if rv.Kind() != reflect.Slice {
		t.Fatalf("not a collection: %T", v)
	}
	out := make([]any, rv.Len())
	for i := range out {
		out[i] = rv.Index(i).Interface()
	}
	return out
}

// truthEntries returns the entries of a Java map (anything with All()
// returning an iter.Seq2).
func truthEntries(t *testing.T, v any) (keys, values []any) {
	t.Helper()
	m := reflect.ValueOf(v).MethodByName("All")
	if !m.IsValid() {
		t.Fatalf("not a map: %T", v)
	}
	seq := m.Call(nil)[0]
	yieldType := seq.Type().In(0)
	yield := reflect.MakeFunc(yieldType, func(args []reflect.Value) []reflect.Value {
		keys = append(keys, args[0].Interface())
		values = append(values, args[1].Interface())
		return []reflect.Value{reflect.ValueOf(true)}
	})
	seq.Call([]reflect.Value{yield})
	return
}

func (s *subject) IsEmpty() {
	s.t.Helper()
	if n := len(truthElements(s.t, s.actual)); n != 0 {
		s.fail("expected to be empty, has %d elements", n)
	}
}

func (s *subject) IsNotEmpty() {
	s.t.Helper()
	if len(truthElements(s.t, s.actual)) == 0 {
		s.fail("expected not to be empty")
	}
}

func (s *subject) HasSize(n int) {
	s.t.Helper()
	if got := len(truthElements(s.t, s.actual)); got != n {
		s.fail("expected size %d, was %d", n, got)
	}
}

func (s *subject) Contains(e any) {
	s.t.Helper()
	if str, ok := truthNorm(s.actual).(string); ok {
		// StringSubject.contains
		if !strings.Contains(str, truthNorm(e).(string)) {
			s.fail("expected %q to contain %q", str, e)
		}
		return
	}
	for _, x := range truthElements(s.t, s.actual) {
		if truthEquals(x, e) {
			return
		}
	}
	s.fail("expected to contain %s, was %v", truthString(e), truthElements(s.t, s.actual))
}

func (s *subject) DoesNotContain(e any) {
	s.t.Helper()
	if str, ok := truthNorm(s.actual).(string); ok {
		if strings.Contains(str, truthNorm(e).(string)) {
			s.fail("expected %q not to contain %q", str, e)
		}
		return
	}
	for _, x := range truthElements(s.t, s.actual) {
		if truthEquals(x, e) {
			s.fail("expected not to contain %s", truthString(e))
		}
	}
}

type orderedSubject struct {
	s        *subject
	actual   []any
	expected []any
}

// ContainsExactly is IterableSubject.containsExactly(Object...), or
// MapSubject.containsExactly(k1, v1, ...) for a map.
func (s *subject) ContainsExactly(expected ...any) *orderedSubject {
	s.t.Helper()
	if m := reflect.ValueOf(s.actual).MethodByName("All"); m.IsValid() && m.Type().NumOut() == 1 && m.Type().Out(0).NumIn() == 1 && m.Type().Out(0).In(0).NumIn() == 2 {
		keys, values := truthEntries(s.t, s.actual)
		var actual []any
		for i := range keys {
			actual = append(actual, keys[i], values[i])
		}
		return s.containsExactly(actual, expected)
	}
	return s.containsExactly(truthElements(s.t, s.actual), expected)
}

func (s *subject) ContainsExactlyElementsIn(expected any) *orderedSubject {
	s.t.Helper()
	return s.containsExactly(truthElements(s.t, s.actual), truthElements(s.t, expected))
}

func (s *subject) containsExactly(actual, expected []any) *orderedSubject {
	s.t.Helper()
	used := make([]bool, len(actual))
	ok := len(actual) == len(expected)
	for _, e := range expected {
		found := false
		for i, a := range actual {
			if !used[i] && truthEquals(a, e) {
				used[i], found = true, true
				break
			}
		}
		if !found {
			ok = false
		}
	}
	if !ok {
		s.fail("expected exactly %v, was %v", truthStrings(expected), truthStrings(actual))
	}
	return &orderedSubject{s: s, actual: actual, expected: expected}
}

func truthStrings(xs []any) []string {
	out := make([]string, len(xs))
	for i, x := range xs {
		out[i] = truthString(x)
	}
	return out
}

func (o *orderedSubject) InOrder() {
	o.s.t.Helper()
	for i := range o.expected {
		if !truthEquals(o.actual[i], o.expected[i]) {
			o.s.fail("expected in order %v, was %v", truthStrings(o.expected), truthStrings(o.actual))
		}
	}
}

// ContainsEntry is MapSubject.containsEntry(Object, Object).
func (s *subject) ContainsEntry(k, v any) {
	s.t.Helper()
	keys, values := truthEntries(s.t, s.actual)
	for i := range keys {
		if truthEquals(keys[i], k) && truthEquals(values[i], v) {
			return
		}
	}
	s.fail("expected to contain entry %s=%s", truthString(k), truthString(v))
}

// nodeSubject is a subset of rhino/testing/NodeSubject.
type nodeSubject struct {
	t      *testing.T
	actual rhino.Node
}

func assertNode(t *testing.T, n rhino.Node) *nodeSubject {
	return &nodeSubject{t: t, actual: n}
}

func (s *nodeSubject) HasToken(tok rhino.Token) *nodeSubject {
	s.t.Helper()
	if s.actual.IsNil() {
		s.t.Fatalf("expected a node with token %v, was null", tok)
	}
	if s.actual.GetToken() != tok {
		s.t.Fatalf("expected token %v, was %v", tok, s.actual.GetToken())
	}
	return s
}

func (s *nodeSubject) HasType(tok rhino.Token) *nodeSubject {
	s.t.Helper()
	return s.HasToken(tok)
}

func (s *nodeSubject) HasLineno(lineno int) *nodeSubject {
	s.t.Helper()
	if got := s.actual.GetLineno(); got != lineno {
		s.t.Fatalf("expected lineno %d, was %d", lineno, got)
	}
	return s
}

func (s *nodeSubject) HasCharno(charno int) *nodeSubject {
	s.t.Helper()
	if got := s.actual.GetCharno(); got != charno {
		s.t.Fatalf("expected charno %d, was %d", charno, got)
	}
	return s
}

func (s *nodeSubject) HasLength(length int) *nodeSubject {
	s.t.Helper()
	if got := s.actual.GetLength(); got != length {
		s.t.Fatalf("expected length %d, was %d", length, got)
	}
	return s
}

// IsString asserts a string node with the given value.
func (s *nodeSubject) IsString(expected string) *nodeSubject {
	s.t.Helper()
	if s.actual.GetToken() != rhino.Token_STRINGLIT {
		s.t.Fatalf("getToken(): expected STRINGLIT, was %v", s.actual.GetToken())
	}
	if got := s.actual.GetString(); got != expected {
		s.t.Fatalf("expected string %q, was %q", expected, got)
	}
	return s
}

func (s *nodeSubject) HasStringThat() *subject {
	s.t.Helper()
	return assertThat(s.t, s.actual.GetString())
}

// markerAnnotationItems is markers.stream().map(m -> m.getAnnotation().getItem()).
func markerAnnotationItems(markers []*rhino.JSDocInfo_Marker) []string {
	out := []string{}
	for _, m := range markers {
		out = append(out, m.GetAnnotation().GetItem())
	}
	return out
}

// Tests.

func TestParseTypeViaStatic1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	typeNode := jt.parseType("null")
	jt.AssertTypeEqualsNode(jt.NULL_TYPE, typeNode)
}

func TestParseTypeViaStatic2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	typeNode := jt.parseType("string")
	jt.AssertTypeEqualsNode(jt.STRING_TYPE, typeNode)
}

func TestParseTypeViaStatic3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	typeNode := jt.parseType("!Date")
	jt.AssertTypeEqualsNode(jt.DATE_TYPE, typeNode)
}

func TestParseTypeViaStatic4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	typeNode := jt.parseType("boolean|string")
	jt.AssertTypeEqualsNode(jt.CreateUnionType(jt.BOOLEAN_TYPE, jt.STRING_TYPE), typeNode)
}

func TestParseInvalidTypeViaStatic(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	typeNode := jt.parseType("sometype<anothertype")
	assertThat(t, typeNode).IsNull()
}

func TestParseInvalidTypeViaStatic2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	typeNode := jt.parseType("")
	assertThat(t, typeNode).IsNull()
}

func TestParseNamedType1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@type {null}", "Unexpected end of file")).IsNull()
}

func TestParseNamedType2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info :=
		jt.parse(
			"@type null*/",
			"Bad type annotation. Type annotations should have curly braces."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
	jt.AssertTypeEqualsExpr(jt.NULL_TYPE, info.GetType())
}

func TestParseNamedType3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {string}*/")
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, info.GetType())
}

func TestParseNamedType4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	// Multi-line @type.
	info := jt.parse("@type \n {string}*/")
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, info.GetType())
}

func TestParseNamedType5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {!goog.\nBar}*/")
	jt.AssertTypeEqualsExpr(jt.createNamedType("goog.Bar"), info.GetType())
}

func TestParseNamedType6(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {!goog.\n * Bar.\n * Baz}*/")
	jt.AssertTypeEqualsExpr(jt.createNamedType("goog.Bar.Baz"), info.GetType())
}

func TestParseNamedTypeError1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	// To avoid parsing ambiguities, type names must end in a '.' to
	// get the continuation behavior.
	jt.parse(
		"@type {!goog\n * .Bar} */",
		"Bad type annotation. expected closing }"+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseNamedTypeError2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@type {!goog.\n * Bar\n * .Baz} */",
		"Bad type annotation. expected closing }"+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestTypedefType1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@typedef {string} */")
	assertThat(t, info.HasTypedefType()).IsTrue()
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, info.GetTypedefType())
}

func TestTypedefType2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@typedef \n {string}*/")
	assertThat(t, info.HasTypedefType()).IsTrue()
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, info.GetTypedefType())
}

func TestTypedefType3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@typedef \n {(string|number)}*/")
	assertThat(t, info.HasTypedefType()).IsTrue()
	jt.AssertTypeEqualsExpr(jt.CreateUnionType(jt.NUMBER_TYPE, jt.STRING_TYPE), info.GetTypedefType())
}

// https://github.com/google/closure-compiler/issues/2543
func TestTypedefType4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info :=
		jt.parse(
			"@typedef {{\n *  boo: ?,\n *  goo: ?\n * }}\n */\n")
	assertThat(t, info.HasTypedefType()).IsTrue()

	recordType := jt.CreateRecordTypeBuilder().AddProperty("boo", jt.UNKNOWN_TYPE, rhino.NilNode).AddProperty("goo", jt.UNKNOWN_TYPE, rhino.NilNode).Build()
	jt.AssertTypeEqualsExpr(recordType, info.GetTypedefType())
}

func TestParseStringType1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, jt.parse("@type {string}*/").GetType())
}

func TestParseStringType2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.STRING_OBJECT_TYPE, jt.parse("@type {!String}*/").GetType())
}

func TestParseBooleanType1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.BOOLEAN_TYPE, jt.parse("@type {boolean}*/").GetType())
}

func TestParseBooleanType2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.BOOLEAN_OBJECT_TYPE, jt.parse("@type {!Boolean}*/").GetType())
}

func TestParseNumberType1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.NUMBER_TYPE, jt.parse("@type {number}*/").GetType())
}

func TestParseNumberType2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.NUMBER_OBJECT_TYPE, jt.parse("@type {!Number}*/").GetType())
}

func TestParseNullType1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.NULL_TYPE, jt.parse("@type {null}*/").GetType())
}

func TestParseNullType2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.NULL_TYPE, jt.parse("@type {Null}*/").GetType())
}

func TestParseAllType1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("*")
}

func TestParseAllType2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("*?", "*")
}

func TestParseObjectType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.OBJECT_TYPE, jt.parse("@type {!Object}*/").GetType())
}

func TestParseDateType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.DATE_TYPE, jt.parse("@type {!Date}*/").GetType())
}

func TestParseFunctionType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.CreateNullableType(jt.FUNCTION_TYPE), jt.parse("@type {Function}*/").GetType())
}

func TestParseRegExpType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.REGEXP_TYPE, jt.parse("@type {!RegExp}*/").GetType())
}

func TestParseUndefinedType1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.VOID_TYPE, jt.parse("@type {undefined}*/").GetType())
}

func TestParseUndefinedType2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.VOID_TYPE, jt.parse("@type {Undefined}*/").GetType())
}

func TestParseUndefinedType3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.VOID_TYPE, jt.parse("@type {void}*/").GetType())
}

func TestParseTemplatizedTypeAlternateSyntax(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {!Array<number>} */")
	jt.AssertTypeEqualsExpr(jt.CreateTemplatizedType(jt.ARRAY_TYPE, jt.NUMBER_TYPE), info.GetType())
}

func TestParseTemplatizedType1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info :=
		jt.parse(
			"@type !Array<number> */",
			"Bad type annotation. Type annotations should have curly braces."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
	jt.AssertTypeEqualsExpr(jt.CreateTemplatizedType(jt.ARRAY_TYPE, jt.NUMBER_TYPE), info.GetType())
}

func TestParseTemplatizedType2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {!Array<number>}*/")
	jt.AssertTypeEqualsExpr(jt.CreateTemplatizedType(jt.ARRAY_TYPE, jt.NUMBER_TYPE), info.GetType())
}

func TestParseTemplatizedType4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {!Array<(number|null)>}*/")
	jt.AssertTypeEqualsExpr(jt.CreateTemplatizedType(jt.ARRAY_TYPE, jt.CreateUnionType(jt.NUMBER_TYPE, jt.NULL_TYPE)), info.GetType())
}

func TestParseTemplatizedType5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {!Array<Array<(number|null)>>}*/")
	jt.AssertTypeEqualsExpr(jt.CreateTemplatizedType(jt.ARRAY_TYPE, jt.CreateUnionType(jt.NULL_TYPE, jt.CreateTemplatizedType(jt.ARRAY_TYPE, jt.CreateUnionType(jt.NUMBER_TYPE, jt.NULL_TYPE)))),
		info.GetType())
}

func TestParseTemplatizedType6(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {!Array<!Array<(number|null)>>}*/")
	jt.AssertTypeEqualsExpr(jt.CreateTemplatizedType(jt.ARRAY_TYPE, jt.CreateTemplatizedType(jt.ARRAY_TYPE, jt.CreateUnionType(jt.NUMBER_TYPE, jt.NULL_TYPE))),
		info.GetType())
}

func TestParseTemplatizedType7(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {!Array<function():Date>}*/")
	jt.AssertTypeEqualsExpr(jt.CreateTemplatizedType(jt.ARRAY_TYPE, jt.Registry.CreateFunctionType(jt.CreateUnionType(jt.DATE_TYPE, jt.NULL_TYPE))),
		info.GetType())
}

func TestParseTemplatizedType8(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {!Array<function():!Date>}*/")
	jt.AssertTypeEqualsExpr(jt.CreateTemplatizedType(jt.ARRAY_TYPE, jt.Registry.CreateFunctionType(jt.DATE_TYPE)), info.GetType())
}

func TestParseTemplatizedType9(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {!Array<Date|number>}*/")
	jt.AssertTypeEqualsExpr(jt.CreateTemplatizedType(jt.ARRAY_TYPE, jt.CreateUnionType(jt.DATE_TYPE, jt.NUMBER_TYPE, jt.NULL_TYPE)),
		info.GetType())
}

func TestParseTemplatizedType10(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {!Array<Date|number|boolean>}*/")
	jt.AssertTypeEqualsExpr(jt.CreateTemplatizedType(jt.ARRAY_TYPE, jt.CreateUnionType(jt.DATE_TYPE, jt.NUMBER_TYPE, jt.BOOLEAN_TYPE, jt.NULL_TYPE)),
		info.GetType())
}

func TestParseTemplatizedType11(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {!Object<number>}*/")
	jt.AssertTypeEqualsExpr(jt.CreateTemplatizedTypeOfList(jt.OBJECT_TYPE, []rhino.JSType{jt.UNKNOWN_TYPE, jt.NUMBER_TYPE}),
		info.GetType())
	jt.assertTemplatizedTypeEquals(jt.Registry.GetObjectElementKey(), jt.NUMBER_TYPE, info.GetType())
}

func TestParseTemplatizedType12(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {!Object<string,number>}*/")
	jt.AssertTypeEqualsExpr(jt.CreateTemplatizedTypeOfList(jt.OBJECT_TYPE, []rhino.JSType{jt.STRING_TYPE, jt.NUMBER_TYPE}),
		info.GetType())
	jt.assertTemplatizedTypeEquals(jt.Registry.GetObjectElementKey(), jt.NUMBER_TYPE, info.GetType())
	jt.assertTemplatizedTypeEquals(jt.Registry.GetObjectIndexKey(), jt.STRING_TYPE, info.GetType())
}

func TestParseTemplatizedType13(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {!Array<?>} */")
	jt.AssertTypeEqualsExpr(jt.CreateTemplatizedType(jt.ARRAY_TYPE, jt.UNKNOWN_TYPE), info.GetType())
}

func TestParseUnionType2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {boolean|null}*/")
	jt.AssertTypeEqualsExpr(jt.CreateUnionType(jt.BOOLEAN_TYPE, jt.NULL_TYPE), info.GetType())
}

func TestParseUnionType6(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {Array<boolean>|null}*/")
	jt.AssertTypeEqualsExpr(jt.CreateUnionType(jt.CreateTemplatizedType(jt.ARRAY_TYPE, jt.BOOLEAN_TYPE), jt.NULL_TYPE),
		info.GetType())
}

func TestParseUnionType7(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {null|Array<boolean>}*/")
	jt.AssertTypeEqualsExpr(jt.CreateUnionType(jt.CreateTemplatizedType(jt.ARRAY_TYPE, jt.BOOLEAN_TYPE), jt.NULL_TYPE),
		info.GetType())
}

func TestParseUnionType10(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@type {string|}*/",
		"Bad type annotation. type not recognized due to syntax error."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseUnionType12(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@type {()}*/",
		"Bad type annotation. type not recognized due to syntax error."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseUnionType13(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("(function(this:Date)|function(this:String):number)", "Function")
}

func TestParseUnionType14(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType(
		"(function(...(function(number):boolean)):number)|function(this:String, string):number",
		"Function")
}

func TestParseUnionType15(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("*|number", "*")
}

func TestParseUnionType16(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("number|*", "*")
}

func TestParseUnionType17(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("string|number|*", "*")
}

func TestParseUnionType19(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {(?)} */")
	jt.AssertTypeEqualsExpr(jt.UNKNOWN_TYPE, info.GetType())
}

func TestParseUnionType_functionReturn_threeAlts(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	// Verify that this is parsed as a union of three types, not a function type returning a union.
	type_ :=
		jt.testParseType("function(): string|number|boolean", "(boolean|function(): string|number)")
	jt.AssertType(type_).IsUnionType()
	alts := type_.ToMaybeUnionType().GetAlternates()
	assertThat(t, alts).HasSize(3)
	jt.AssertType(alts[0]).IsFunctionTypeThat().HasReturnTypeThat().IsString()
	jt.AssertType(alts[1]).IsNumber()
	jt.AssertType(alts[2]).IsBoolean()
}

func TestParseUnionTypeError2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@type {string||number} */",
		"Bad type annotation. type not recognized due to syntax error."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseUnknownType1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("?")
}

func TestParseUnknownType2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("(?|number)", "?")
}

func TestParseUnknownType3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("(number|?)", "?")
}

func TestParseUnknownType_undefined(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("(undefined|?)", "(?|undefined)")
}

func TestParseUnknownType_null(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("(null|?)", "?")
}

func TestParseTemplatizedUnknown1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@ type {?<?>} */",
		"illegal use of unknown JSDoc tag \"\"; ignoring it. Place another character before the @"+
			" to stop JSCompiler from parsing it as an annotation.")
}

func TestParseTemplatizedUnknown2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@ type {{b: ?<?>}} */",
		"illegal use of unknown JSDoc tag \"\"; ignoring it. Place another character before the @"+
			" to stop JSCompiler from parsing it as an annotation.")
}

func TestParseFunctionalType1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("function(): number")
}

func TestParseFunctionalType2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("function(number, string): boolean")
}

func TestParseFunctionalType3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("function(this:Array)", "function(this:Array): ?")
}

func TestParseFunctionalType4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("function(...number): boolean")
}

func TestParseFunctionalType5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("function(number, ...string): boolean")
}

func TestParseFunctionalType6(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("function(this:Date, number): (boolean|number|string)")
}

func TestParseFunctionalType7(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("function()", "function(): ?")
}

func TestParseFunctionalType9(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType(
		"function(this:Array,!Date,...(boolean?))",
		"function(this:Array, Date, ...(boolean|null)): ?")
}

func TestParseFunctionalType10(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("function(...(Object?)):boolean?", "function(...(Object|null)): (boolean|null)")
}

func TestParseFunctionalType12(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("function(...)", "function(...?): ?")
}

func TestParseFunctionalType13(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("function(...): void", "function(...?): undefined")
}

func TestParseFunctionalType14(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("function(*, string, number): boolean")
}

func TestParseFunctionalType15(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("function(?, string): boolean")
}

func TestParseFunctionalType16(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("function(string, ?): ?")
}

func TestParseFunctionalType17(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("(function(?): ?|number)")
}

func TestParseFunctionalType18(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("function (?): (?|number)", "function(?): ?")
}

func TestParseFunctionalType19(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.testParseType("function(...?): void", "function(...?): undefined")
}

func TestStructuralConstructor(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	type_ := jt.testParseType("function (new:Object)", "function(new:Object): ?")
	assertThat(t, type_.IsConstructor()).IsTrue()
	assertThat(t, type_.IsNominalConstructorOrInterface()).IsFalse()
}

func TestStructuralConstructor2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	type_ :=
		jt.testParseType(
			"function (new:?)",
			// toString skips unknowns, but isConstructor reveals the truth.
			"function(): ?")
	assertThat(t, type_.IsConstructor()).IsTrue()
	assertThat(t, type_.IsNominalConstructorOrInterface()).IsFalse()
}

func TestStructuralConstructor3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.Resolve(
		jt.parse("@type {function (new:*)} */").GetType(), "constructed type must be an object type")
}

func TestNominalConstructor(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	type_ := jt.testParseType("Array", "(Array|null)").Dereference()
	assertThat(t, type_.GetConstructor().IsNominalConstructorOrInterface()).IsTrue()
}

func TestBug1419535(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@type {function(Object, string, *)?} */")
	jt.parse("@type {function(Object, string, *)|null} */")
}

func TestIssue477(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@type {function} */", "Bad type annotation. missing opening ("+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestMalformedThisAnnotation(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@this */",
		"Bad type annotation. Type annotations should have curly braces."+JsDocInfoParser_BAD_TYPE_WIKI_LINK,
		"Bad type annotation. type not recognized due to syntax error."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseImportTypeError(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@type {import('http').Stream} */",
		"Bad type annotation. Import in typedef is not supported."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseFunctionalTypeError1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@type {function number):string}*/",
		"Bad type annotation. missing opening ("+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseFunctionalTypeError2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@type {function( number}*/",
		"Bad type annotation. missing closing )"+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseFunctionalTypeError3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@type {function(...[number], string)}*/",
		"Bad type annotation. type not recognized due to syntax error."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseFunctionalTypeError4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@type {function(string, ...number, boolean):string}*/",
		"Bad type annotation. variable length argument must be last."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseFunctionalTypeError5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@type {function (thi:Array)}*/",
		"Bad type annotation. missing closing )"+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseFunctionalTypeError6(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.Resolve(jt.parse("@type {function (this:number)}*/").GetType())
}

func TestParseFunctionalTypeError7(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@type {function(...[number)}*/",
		"Bad type annotation. type not recognized due to syntax error."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseFunctionalTypeError8(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@type {function(...number])}*/",
		"Bad type annotation. missing closing )"+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseFunctionalTypeError9(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@type {function (new:Array, this:Object)} */",
		"Bad type annotation. missing closing )"+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseFunctionalTypeError10(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@type {function (this:Array, new:Object)} */",
		"Bad type annotation. missing closing )"+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseFunctionalTypeError11(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@type {function (Array, new:Object)} */",
		"Bad type annotation. missing closing )"+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseFunctionalTypeError12(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.Resolve(
		jt.parse("@type {function (new:number)}*/").GetType(),
		"constructed type must be an object type")
}

func TestParseFunctionalTypeError13(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@type {function (...[number]): boolean} */",
		"Bad type annotation. type not recognized due to syntax error."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseFunctionalTypeError14(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@type {function (number, ...[string]): boolean} */",
		"Bad type annotation. type not recognized due to syntax error."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseFunctionalType8(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@type {function(this:Array,...[boolean])} */",
		"Bad type annotation. type not recognized due to syntax error."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseArrayTypeError1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@type {[number}*/",
		"Bad type annotation. type not recognized due to syntax error."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseArrayTypeError2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@type {number]}*/", "Bad type annotation. expected closing }"+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseArrayTypeError3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@type {[(number,boolean,Object?])]}*/",
		"Bad type annotation. type not recognized due to syntax error."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseArrayTypeError5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@type {[Object]}*/",
		"Bad type annotation. type not recognized due to syntax error."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseTypeofType1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	node := jt.parse("@type {typeof Foo}*/").GetType().GetRoot()
	assertNode(t, node).HasToken(rhino.Token_TYPEOF)
	assertThat(t, node.GetChildCount()).IsEqualTo(1)
	fooNode := node.GetFirstChild()
	assertNode(t, fooNode).HasToken(rhino.Token_STRINGLIT)
	assertNode(t, fooNode).HasStringThat().IsEqualTo("Foo")
	assertNode(t, fooNode).HasCharno(14)
}

func TestParseTypeofType2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	node := jt.parse("@type {(typeof Foo)}*/").GetType().GetRoot()
	assertNode(t, node).HasToken(rhino.Token_TYPEOF)
	assertThat(t, node.GetChildCount()).IsEqualTo(1)
	fooNode := node.GetFirstChild()
	assertNode(t, fooNode).HasToken(rhino.Token_STRINGLIT)
	assertNode(t, fooNode).HasStringThat().IsEqualTo("Foo")
	assertNode(t, fooNode).HasCharno(15)
}

func TestParseTypeofType3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	node := jt.parse("@type {typeof Foo|Bar<typeof Baz>}*/").GetType().GetRoot()
	assertNode(t, node).HasToken(rhino.Token_PIPE)
	assertNode(t, node.GetFirstChild()).HasToken(rhino.Token_TYPEOF)
	fooNode := node.GetFirstFirstChild()
	assertNode(t, fooNode).HasToken(rhino.Token_STRINGLIT)
	assertNode(t, fooNode).HasStringThat().IsEqualTo("Foo")
	assertNode(t, fooNode).HasCharno(14)
	barNode := node.GetLastChild()
	assertNode(t, barNode).HasToken(rhino.Token_STRINGLIT)
	assertNode(t, barNode).HasStringThat().IsEqualTo("Bar")
	assertNode(t, barNode).HasCharno(18)
	assertNode(t, barNode.GetFirstChild()).HasToken(rhino.Token_BLOCK)
	assertNode(t, barNode.GetFirstFirstChild()).HasToken(rhino.Token_TYPEOF)
	bazNode := barNode.GetFirstFirstChild().GetFirstChild()
	assertNode(t, bazNode).HasToken(rhino.Token_STRINGLIT)
	assertNode(t, bazNode).HasStringThat().IsEqualTo("Baz")
	assertNode(t, bazNode).HasCharno(29)
}

func TestParseTypeofTypeNewLine(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	node := jt.parse("@type {typeof \n  Foo}*/").GetType().GetRoot()
	assertNode(t, node).HasToken(rhino.Token_TYPEOF)
	assertThat(t, node.GetChildCount()).IsEqualTo(1)
	fooNode := node.GetFirstChild()
	assertNode(t, fooNode).HasToken(rhino.Token_STRINGLIT)
	assertNode(t, fooNode).HasStringThat().IsEqualTo("Foo")
	assertNode(t, fooNode).HasCharno(2)
	assertNode(t, fooNode).HasLineno(1)
}

func TestTypenameSourceInfo_simpleName(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	node := jt.parse("@type {Foo} */").GetType().GetRoot()
	assertNode(t, node).HasToken(rhino.Token_STRINGLIT)
	assertThat(t, node.GetLineno()).IsEqualTo(0)
	assertThat(t, node.GetCharno()).IsEqualTo(7)
	assertThat(t, node.GetLength()).IsEqualTo(3)
}

func TestTypenameSourceInfo_qualifiedName(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	node := jt.parse("@type {bar.baz.Foo} */").GetType().GetRoot()
	assertNode(t, node).HasToken(rhino.Token_STRINGLIT)
	assertThat(t, node.GetLineno()).IsEqualTo(0)
	assertThat(t, node.GetCharno()).IsEqualTo(7)
	assertThat(t, node.GetLength()).IsEqualTo(11)
}

func TestTypenameSourceInfo_qualifiedName_multiline(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	node :=
		jt.parse(
			"@type {bar.\n       baz.Foo} */\n").GetType().GetRoot()
	assertNode(t, node).HasToken(rhino.Token_STRINGLIT)
	assertThat(t, node.GetLineno()).IsEqualTo(0)
	assertThat(t, node.GetCharno()).IsEqualTo(7)
	assertThat(t, node.GetLength()).IsEqualTo(19) // Including newline and spaces.
}

func TestTypenameSourceInfo_qualifiedName_leadingWhitespace(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	node :=
		jt.parse(
			"@type {\n    Foo} */\n").GetType().GetRoot()
	assertNode(t, node).HasToken(rhino.Token_STRINGLIT)
	assertThat(t, node.GetLineno()).IsEqualTo(1)
	assertThat(t, node.GetCharno()).IsEqualTo(4)
	assertThat(t, node.GetLength()).IsEqualTo(3)
}

func TestParseNullableModifiers1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {string?}*/")
	jt.AssertTypeEqualsExpr(jt.CreateNullableType(jt.STRING_TYPE), info.GetType())
}

func TestParseNullableModifiers2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {!Array<string?>}*/")
	jt.AssertTypeEqualsExpr(jt.CreateTemplatizedType(jt.ARRAY_TYPE, jt.CreateUnionType(jt.STRING_TYPE, jt.NULL_TYPE)), info.GetType())
}

func TestParseNullableModifiers3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {Array<boolean>?}*/")
	jt.AssertTypeEqualsExpr(jt.CreateNullableType(jt.CreateTemplatizedType(jt.ARRAY_TYPE, jt.BOOLEAN_TYPE)), info.GetType())
}

func TestParseNullableModifiers7(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {string?|boolean}*/")
	jt.AssertTypeEqualsExpr(jt.CreateUnionType(jt.CreateNullableType(jt.STRING_TYPE), jt.BOOLEAN_TYPE), info.GetType())
}

func TestParseNullableModifiers8(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {string|boolean?}*/")
	jt.AssertTypeEqualsExpr(jt.CreateUnionType(jt.STRING_TYPE, jt.CreateNullableType(jt.BOOLEAN_TYPE)), info.GetType())
}

func TestParseNullableModifiers9(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {foo.Hello.World?}*/")
	jt.AssertTypeEqualsExpr(jt.CreateNullableType(jt.createNamedType("foo.Hello.World")), info.GetType())
}

func TestParseOptionalModifier(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {function(number=)}*/")
	jt.AssertTypeEqualsExpr(jt.Registry.CreateFunctionTypeWithParameters(jt.UNKNOWN_TYPE, jt.Registry.CreateOptionalParameters(jt.NUMBER_TYPE)),
		info.GetType())
}

func TestParseNewline1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {string\n* }\n*/")
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, info.GetType())
}

func TestParseNewline2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {!Array<\n* number\n* >} */")
	jt.AssertTypeEqualsExpr(jt.CreateTemplatizedType(jt.ARRAY_TYPE, jt.NUMBER_TYPE), info.GetType())
}

func TestParseNewline4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {!Array<(number|\n* null)>}*/")
	jt.AssertTypeEqualsExpr(jt.CreateTemplatizedType(jt.ARRAY_TYPE, jt.CreateUnionType(jt.NUMBER_TYPE, jt.NULL_TYPE)), info.GetType())
}

func TestParseNewline5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {!Array<function(\n* )\n* :\n* Date>}*/")
	jt.AssertTypeEqualsExpr(jt.CreateTemplatizedType(jt.ARRAY_TYPE, jt.Registry.CreateFunctionType(jt.CreateUnionType(jt.DATE_TYPE, jt.NULL_TYPE))),
		info.GetType())
}

func TestParseReturnType1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@return {null|string|Array<boolean>}*/")
	jt.AssertTypeEqualsExpr(jt.CreateUnionType(jt.CreateTemplatizedType(jt.ARRAY_TYPE, jt.BOOLEAN_TYPE), jt.NULL_TYPE, jt.STRING_TYPE),
		info.GetReturnType())
}

func TestParseReturnType2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@returns {null|(string|Array<boolean>)}*/")
	jt.AssertTypeEqualsExpr(jt.CreateUnionType(jt.CreateTemplatizedType(jt.ARRAY_TYPE, jt.BOOLEAN_TYPE), jt.NULL_TYPE, jt.STRING_TYPE),
		info.GetReturnType())
}

func TestParseReturnType3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@return {((null|Array<boolean>|string)|boolean)}*/")
	jt.AssertTypeEqualsExpr(jt.CreateUnionType(jt.CreateTemplatizedType(jt.ARRAY_TYPE, jt.BOOLEAN_TYPE), jt.NULL_TYPE, jt.STRING_TYPE, jt.BOOLEAN_TYPE),
		info.GetReturnType())
}

func TestParseThisType1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@this {goog.foo.Bar}*/")
	jt.AssertTypeEqualsExpr(jt.createNamedType("goog.foo.Bar"), info.GetThisType())
}

func TestParseThisType2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info :=
		jt.parse(
			"@this goog.foo.Bar*/",
			"Bad type annotation. Type annotations should have curly braces."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
	jt.AssertTypeEqualsExpr(jt.createNamedType("goog.foo.Bar"), info.GetThisType())
}

func TestParseThisType3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@type {number}\n@this {goog.foo.Bar}*/",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseThisType4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.Resolve(jt.parse("@this {number}*/").GetThisType())
}

func TestParseThisType5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@this {Date|Error}*/")
}

func TestParseThisType6(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.Resolve(jt.parse("@this {Date|number}*/").GetThisType())
}

func TestParseParam1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@param {number} index*/")
	assertThat(t, info.GetParameterCount()).IsEqualTo(1)
	jt.AssertTypeEqualsExpr(jt.NUMBER_TYPE, info.GetParameterType("index"))
}

func TestParseParam2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@param index*/", MISSING_TYPE_DECL_WARNING_TEXT)
	assertThat(t, info.GetParameterCount()).IsEqualTo(1)
	assertThat(t, info.GetParameterType("index")).IsNull()
}

func TestParseParam3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@param {number} index useful comments*/")
	assertThat(t, info.GetParameterCount()).IsEqualTo(1)
	jt.AssertTypeEqualsExpr(jt.NUMBER_TYPE, info.GetParameterType("index"))
}

func TestParseParam4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@param index useful comments*/", MISSING_TYPE_DECL_WARNING_TEXT)
	assertThat(t, info.GetParameterCount()).IsEqualTo(1)
	assertThat(t, info.GetParameterType("index")).IsNull()
}

func TestParseParam5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	// Test for multi-line @param.
	info := jt.parse("@param {number} \n index */")
	assertThat(t, info.GetParameterCount()).IsEqualTo(1)
	jt.AssertTypeEqualsExpr(jt.NUMBER_TYPE, info.GetParameterType("index"))
}

func TestParseParam6(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	// Test for multi-line @param.
	info := jt.parse("@param {number} \n * index */")
	assertThat(t, info.GetParameterCount()).IsEqualTo(1)
	jt.AssertTypeEqualsExpr(jt.NUMBER_TYPE, info.GetParameterType("index"))
}

func TestParseParam7(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	// Optional @param
	info := jt.parse("@param {number=} index */")
	jt.AssertTypeEqualsExpr(jt.Registry.CreateOptionalType(jt.NUMBER_TYPE), info.GetParameterType("index"))
}

func TestParseParam8(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	// Var args @param
	info := jt.parse("@param {...number} index */")
	jt.AssertTypeEqualsExpr(jt.NUMBER_TYPE, info.GetParameterType("index"))
}

func TestParseParam9(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@param {...number=} index */",
		"Bad type annotation. expected closing }"+JsDocInfoParser_BAD_TYPE_WIKI_LINK,
		"Bad type annotation. expecting a variable name in a @param tag."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseParam10(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@param {...number index */",
		"Bad type annotation. expected closing }"+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseParam11(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@param {number= index */", "Bad type annotation. expected closing }"+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseParam12(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@param {...number|string} index */")
	jt.AssertTypeEqualsExpr(jt.Registry.CreateUnionType(jt.STRING_TYPE, jt.NUMBER_TYPE), info.GetParameterType("index"))
}

func TestParseParam13(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@param {...(number|string)} index */")
	jt.AssertTypeEqualsExpr(jt.Registry.CreateUnionType(jt.STRING_TYPE, jt.NUMBER_TYPE), info.GetParameterType("index"))
}

func TestParseParam14(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@param {string} [index] */")
	assertThat(t, info.GetParameterCount()).IsEqualTo(1)
	jt.AssertTypeEqualsExpr(jt.Registry.CreateOptionalType(jt.STRING_TYPE), info.GetParameterType("index"))
}

func TestParseParam15(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info :=
		jt.parse(
			"@param {string} [index */",
			"Bad type annotation. missing closing ]"+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
	assertThat(t, info.GetParameterCount()).IsEqualTo(1)
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, info.GetParameterType("index"))
}

func TestParseParam16(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@param {string} index] */")
	assertThat(t, info.GetParameterCount()).IsEqualTo(1)
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, info.GetParameterType("index"))
}

func TestParseParam17(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@param {string=} [index] */")
	assertThat(t, info.GetParameterCount()).IsEqualTo(1)
	jt.AssertTypeEqualsExpr(jt.Registry.CreateOptionalType(jt.STRING_TYPE), info.GetParameterType("index"))
}

func TestParseParam18(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@param {...string} [index] */")
	assertThat(t, info.GetParameterCount()).IsEqualTo(1)
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, info.GetParameterType("index"))
}

func TestParseParam19(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@param {...} [index] */")
	assertThat(t, info.GetParameterCount()).IsEqualTo(1)
	jt.AssertTypeEqualsExpr(jt.Registry.CreateOptionalType(jt.UNKNOWN_TYPE), info.GetParameterType("index"))
	assertThat(t, info.GetParameterType("index").IsVarArgs()).IsTrue()
}

func TestParseParam20(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@param {?=} index */")
	assertThat(t, info.GetParameterCount()).IsEqualTo(1)
	jt.AssertTypeEqualsExpr(jt.UNKNOWN_TYPE, info.GetParameterType("index"))
}

func TestParseParam21(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@param {...?} index */")
	assertThat(t, info.GetParameterCount()).IsEqualTo(1)
	jt.AssertTypeEqualsExpr(jt.UNKNOWN_TYPE, info.GetParameterType("index"))
	assertThat(t, info.GetParameterType("index").IsVarArgs()).IsTrue()
}

func TestParseParam22(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@param {string} .index */", "invalid param name \".index\"")
	assertThat(t, info).IsNull()
}

func TestParseParam23(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@param {string} index. */", "invalid param name \"index.\"")
	assertThat(t, info).IsNull()
}

func TestParseParam24(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@param {string} foo.bar */", "invalid param name \"foo.bar\"")
	assertThat(t, info).IsNull()
}

func TestParseParam25(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info :=
		jt.parse(
			"@param {string} foo.bar\n * @param {string} baz */", "invalid param name \"foo.bar\"")
	assertThat(t, info.GetParameterCount()).IsEqualTo(1)
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, info.GetParameterType("baz"))
}

func TestParseParam26(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info :=
		jt.parse(
			"@param {{a: number, b: number}} {a, b}\n*/",
			"Bad type annotation. expecting a variable name in a @param tag."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
	assertThat(t, info).IsNull()
}

func TestParseParam27(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info :=
		jt.parse("@param {{a: number, b: number}} '{a, b}'\n*/", "invalid param name \"'\"")
	assertThat(t, info).IsNull()
}

func TestParseParam28(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@param {...} index */")
	assertThat(t, info.GetParameterCount()).IsEqualTo(1)
	jt.AssertTypeEqualsExpr(jt.Registry.CreateOptionalType(jt.UNKNOWN_TYPE), info.GetParameterType("index"))
	assertThat(t, info.GetParameterType("index").IsVarArgs()).IsTrue()
}

func TestParseThrows1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@throws {number} Some number */", true)
	assertThat(t, info.GetThrowsAnnotations()).ContainsExactly("{number} Some number")
}

func TestParseThrows2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info :=
		jt.parse(
			"@throws {not a type}\n *     Bla\n*@throws {String} A string */\n",
			true)
	assertThat(t, info.GetThrowsAnnotations()).ContainsExactly("{not a type} Bla", "{String} A string").InOrder()
}

func TestParseRecordType1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull("/** @param {{x}} n\n*/")
}

func TestParseRecordType2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull("/** @param {{z, y}} n\n*/")
}

func TestParseRecordType3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull("/** @param {{z, y, x, q, hello, thisisatest}} n\n*/")
}

func TestParseRecordType4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull("/** @param {{a, 'a', 'hello', 2, this, do, while, for}} n\n*/")
}

func TestParseRecordType5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull("/** @param {{x : hello}} n\n*/")
}

func TestParseRecordType6(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull("/** @param {{'x' : hello}} n\n*/")
}

func TestParseRecordType7(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull("/** @param {{'x' : !hello}} n\n*/")
}

func TestParseRecordType8(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull("/** @param {{'x' : !hello, y : bar}} n\n*/")
}

func TestParseRecordType9(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull("/** @param {{'x' : !hello, y : {z : bar, 3 : meh}}} n\n*/")
}

func TestParseRecordType10(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull("/** @param {{__proto__ : moo}} n\n*/")
}

func TestParseRecordType11(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull(
		"/** @param {{a : b} n\n*/",
		"Bad type annotation. expected closing }"+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseRecordType12(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull(
		"/** @param {{!hello : hey}} n\n*/",
		"Bad type annotation. type not recognized due to syntax error."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseRecordType13(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull("/** @param {{x}|number} n\n*/")
}

func TestParseRecordType14(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull("/** @param {{x : y}|number} n\n*/")
}

func TestParseRecordType15(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull("/** @param {{'x' : y}|number} n\n*/")
}

func TestParseRecordType16(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull("/** @param {{x, y}|number} n\n*/")
}

func TestParseRecordType17(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull("/** @param {{x : hello, 'y'}|number} n\n*/")
}

func TestParseRecordType18(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull("/** @param {number|{x : hello, 'y'}} n\n*/")
}

func TestParseRecordType19(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull("/** @param {?{x : hello, 'y'}} n\n*/")
}

func TestParseRecordType20(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull("/** @param {!{x : hello, 'y'}} n\n*/")
}

func TestParseRecordType21(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull("/** @param {{x : hello, 'y'}|boolean} n\n*/")
}

func TestParseRecordType22(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull("/** @param {{x : hello, 'y'}|function()} n\n*/")
}

func TestParseRecordType23(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull("/** @param {{x : function(), 'y'}|function()} n\n*/")
}

func TestParseRecordType24(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull("/** @param {{a : b, c,}} n\n*/")
}

func TestParseParamError1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull(
		"/** @param\n*/",
		"Bad type annotation. expecting a variable name in a @param tag."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseParamError2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull(
		"/** @param {Number}*/",
		"Bad type annotation. expecting a variable name in a @param tag."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseParamError3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull(
		"/** @param {Number}\n*/",
		"Bad type annotation. expecting a variable name in a @param tag."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseParamError4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull(
		"/** @param {Number}\n* * num */",
		"Bad type annotation. expecting a variable name in a @param tag."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseParamError5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@param {number} x \n * @param {string} x */",
		"Bad type annotation. duplicate variable name \"x\""+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseExtends1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.STRING_OBJECT_TYPE, jt.parse("@extends String*/").GetBaseType())
}

func TestParseExtends2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@extends com.google.Foo.Bar.Hello.World*/")
	jt.AssertTypeEqualsExpr(jt.createNamedType("com.google.Foo.Bar.Hello.World"), info.GetBaseType())
}

func TestParseExtendsGenerics(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@extends com.google.Foo.Bar.Hello.World<Boolean,number>*/")
	jt.AssertTypeEqualsExpr(jt.createNamedType("com.google.Foo.Bar.Hello.World"), info.GetBaseType())
}

func TestParseImplementsGenerics(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	// For types that are not templatized, <> annotations are ignored.
	interfaces :=
		jt.parse("@implements {SomeInterface<*>} */").GetImplementedInterfaces()
	assertThat(t, interfaces).HasSize(1)
	jt.AssertTypeEqualsExpr(jt.createNamedType("SomeInterface"), interfaces[0])
}

func TestParseExtends4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.STRING_OBJECT_TYPE, jt.parse("@extends {String}*/").GetBaseType())
}

func TestParseExtends5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.STRING_OBJECT_TYPE,
		jt.parse("@extends {String*/", "Bad type annotation. expected closing }"+JsDocInfoParser_BAD_TYPE_WIKI_LINK).GetBaseType())
}

func TestParseExtends6(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	// Multi-line extends
	jt.AssertTypeEqualsExpr(jt.STRING_OBJECT_TYPE, jt.parse("@extends \n * {String}*/").GetBaseType())
}

func TestParseExtendsInvalidName(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	// This looks bad, but for the time being it should be OK, as
	// we will not find a type with this name in the JS parsed tree.
	// If this is fixed in the future, change this test to check for a
	// warning/error message.
	jt.AssertTypeEqualsExpr(jt.createNamedType("some_++#%$%_UglyString"),
		jt.parse("@extends {some_++#%$%_UglyString} */").GetBaseType())
}

func TestParseExtendsNullable1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@extends {Base?} */", "Bad type annotation. expected closing }"+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseExtendsNullable2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@extends Base? */",
		"Bad type annotation. expected end of line or comment."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseExtendsWithQualifier(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@extends {!Base}*/",
		"Bad type annotation. @implements/@extends requires a bare interface/record name without !"+
			" or ?."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseExtends_multipleExtends(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.jsDocSourceKind = JsDocInfoParser_JsDocSourceKind_NORMAL
	jt.parse(
		"@extends {Foo}\n@extends {Bar} */",
		"Bad type annotation. type annotation incompatible with other annotations. See"+
			" https://github.com/google/closure-compiler/wiki/Annotating-JavaScript-for-the-Closure-Compiler"+
			" for more information.")
}

func TestParseExtends_multipleExtendsTsickleMode(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.jsDocSourceKind = JsDocInfoParser_JsDocSourceKind_TSICKLE
	result := jt.parse("@extends {Foo}\n@extends {Bar} */")
	// Treated as "@extends {JsDocInfoParser_TsickleMode_MissingSupertypePlaceholder}"
	// which has an implicit "!" (rhino.Token_BANG) when parsed.
	assertNode(t, result.GetBaseType().GetRoot()).HasToken(rhino.Token_BANG)
	assertNode(t, result.GetBaseType().GetRoot().GetOnlyChild()).IsString("JsDocInfoParser_TsickleMode_MissingSupertypePlaceholder")
}

func TestParseEnum1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.NUMBER_TYPE, jt.parse("@enum*/").GetEnumParameterType())
}

func TestParseEnum2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, jt.parse("@enum {string}*/").GetEnumParameterType())
}

func TestParseEnum3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.NUMBER_TYPE, jt.parse("@enum {number}*/").GetEnumParameterType())
}

func TestParseEnum4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.BOOLEAN_TYPE, jt.parse("@enum {boolean}*/").GetEnumParameterType())
}

func TestParseEnum5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.SYMBOL_TYPE, jt.parse("@enum {symbol}*/").GetEnumParameterType())
}

func TestParseEnum6(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc := jt.parse(" @enum {Foo} */")
	enumTypeNode := jsdoc.GetEnumParameterType().GetRoot()
	assertNode(t, enumTypeNode).HasToken(rhino.Token_BANG)
}

func TestParseEnum7(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE,
		jt.parse(
			"@enum string*/",
			"Bad type annotation. Type annotations should have curly braces."+
				JsDocInfoParser_BAD_TYPE_WIKI_LINK).GetEnumParameterType())
}

func TestParseBadEnumNoCrash(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.NUMBER_TYPE,
		jt.parse(
			"@enum {@enum {string}}*/",
			"Bad type annotation. type not recognized due to syntax error."+
				JsDocInfoParser_BAD_TYPE_WIKI_LINK,
			"Bad type annotation. type annotation incompatible with other annotations."+
				JsDocInfoParser_BAD_TYPE_WIKI_LINK).GetEnumParameterType())
}

func TestParseJsDocAfterEnum(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@enum {string} @suppress {x} */")
	assertThat(t, info.GetSuppressions()).IsEqualTo(setOf("x"))
}

func TestParseDesc1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@desc hello world!*/").GetDescription()).IsEqualTo("hello world!")
}

func TestParseDesc2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@desc hello world!\n*/").GetDescription()).IsEqualTo("hello world!")
}

func TestParseDesc3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@desc*/").GetDescription()).IsEmpty()
}

func TestParseDesc4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@desc\n*/").GetDescription()).IsEmpty()
}

func TestParseDesc5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@desc hello\nworld!\n*/").GetDescription()).IsEqualTo("hello world!")
}

func TestParseDesc6(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@desc hello\n* world!\n*/").GetDescription()).IsEqualTo("hello world!")
}

func TestParseDesc7(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@desc a\n\nb\nc*/").GetDescription()).IsEqualTo("a b c")
}

func TestParseDesc8(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@desc a\n      *b\n\n  *c\n\nd*/").GetDescription()).IsEqualTo("a b c d")
}

func TestParseDesc9(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	comment := "@desc\n.\n,\n{\n)\n}\n|\n.<\n>\n<\n?\n~\n+\n-\n;\n:\n*/"

	assertThat(t, jt.parse(comment).GetDescription()).IsEqualTo(". , { ) } | < > < ? ~ + - ; :")
}

func TestParseDesc10(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	comment := "@desc\n?\n?\n?\n?*/"

	assertThat(t, jt.parse(comment).GetDescription()).IsEqualTo("? ? ? ?")
}

func TestParseDesc11(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	comment := "@desc :[]*/"

	assertThat(t, jt.parse(comment).GetDescription()).IsEqualTo(":[]")
}

func TestParseDesc12(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	comment := "@desc\n:\n[\n]\n...*/"

	assertThat(t, jt.parse(comment).GetDescription()).IsEqualTo(": [ ] ...")
}

func TestParseMeaning1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@meaning tigers   */").GetMeaning()).IsEqualTo("tigers")
}

func TestParseMeaning2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@meaning tigers\n * and lions\n * and bears */").GetMeaning()).IsEqualTo("tigers and lions and bears")
}

func TestParseMeaning3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@meaning  tigers\n * and lions\n * @desc  and bears */")
	assertThat(t, info.GetMeaning()).IsEqualTo("tigers and lions")
	assertThat(t, info.GetDescription()).IsEqualTo("and bears")
}

func TestParseMeaning4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@meaning  tigers\n * @meaning and lions  */", "extra @meaning tag")
}

func TestParseAlternateMessageId1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@alternateMessageId A   */").GetAlternateMessageId()).IsEqualTo("A")
}

func TestParseAlternateMessageId2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@alternateMessageId A\n * B\n * C */").GetAlternateMessageId()).IsEqualTo("A")
}

func TestParseAlternateMessageId3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@alternateMessageId  A\n * B\n * @desc  C */")
	assertThat(t, info.GetAlternateMessageId()).IsEqualTo("A")
	assertThat(t, info.GetDescription()).IsEqualTo("C")
}

func TestParseAlternateMessageId4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@alternateMessageId  A\n * @alternateMessageId B  */", "extra @alternateMessageId tag")
}

func TestParseClosurePrimitive(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@closurePrimitive {id} */")
	assertThat(t, info.GetClosurePrimitiveId()).IsEqualTo("id")
}

func TestParseClosurePrimitiveOnNewLine(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@closurePrimitive    \n\n {id} */")
	assertThat(t, info.GetClosurePrimitiveId()).IsEqualTo("id")
}

func TestParseClosurePrimitiveWithAnnotationAfterId(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@closurePrimitive {id} @const */")
	assertThat(t, info.GetClosurePrimitiveId()).IsEqualTo("id")
	assertThat(t, info.IsConstant()).IsTrue()
}

func TestParseClosurePrimitiveWithStringAfterId(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@closurePrimitive {id} some comment\n@const */")
	assertThat(t, info.GetClosurePrimitiveId()).IsEqualTo("id")
	assertThat(t, info.IsConstant()).IsTrue()
}

func TestParseTsType_typesOnlyMode(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@tsType ():this */", Config_JsDocParsing_TYPES_ONLY).GetTsTypes()).ContainsExactly("JsDocInfoParser_tsType_PlaceHolder")
}

func TestParseTsType_includeDescriptions(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t,
		jt.parse("@tsType ():this */", Config_JsDocParsing_INCLUDE_DESCRIPTIONS_NO_WHITESPACE).GetTsTypes()).ContainsExactly("():this")
}

func TestParseClosurePrimitiveIdMissingIdentifier(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@closurePrimitive */", "missing opening {")
}

func TestParseClosurePrimitiveIdInvalidId(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@closurePrimitive {unrecognizedId } */", "invalid id in @closurePrimitive tag.")
}

func TestParseClosurePrimitiveIdMissingLeftCurly(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@closurePrimitive id} */", "missing opening {")
}

func TestParseClosurePrimitiveIdMissingRightCurly(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@closurePrimitive {id */", "expected closing }")
}

func TestParseClosurePrimitiveIdMissingStringInCurly(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@closurePrimitive {} */", "missing id in @closurePrimitive tag.")
}

func TestParseClosurePrimitiveDuplicateTags(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@closurePrimitive {idA}\n@closurePrimitive {idB} */", "conflicting @closurePrimitive tag")
}

func TestParseLends1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@lends {name} */")
	assertThat(t, info.GetLendsName().GetRoot().GetString()).IsEqualTo("name")
}

func TestParseLends2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@lends   foo.bar  */")
	assertThat(t, info.GetLendsName().GetRoot().GetString()).IsEqualTo("foo.bar")
}

func TestParseLends3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@lends {name */", "Bad type annotation. expected closing }"+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseLends4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@lends {} */",
		"Bad type annotation. missing object name in @lends tag."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseLends5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@lends } */",
		"Bad type annotation. missing object name in @lends tag."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseLends6(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@lends {string} \n * @lends {string} */",
		"Bad type annotation. @lends tag incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseLends7(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@type {string} \n * @lends {string} */",
		"Bad type annotation. @lends tag incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseAbstract(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@abstract */")
	assertThat(t, info.IsAbstract()).IsTrue()
}

func TestParseAbstract_alreadyAbstract(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@abstract \n * @abstract */",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseAbstract_typedefBeforeAbstract(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@typedef {Object<T,T>} \n * @abstract */",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseAbstract_typeBeforeAbstract(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"* @type {Object<T,T>} \n * @abstract */",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseAbstract_interfaceBeforeAbstract(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"* @interface \n * @abstract */",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseAbstract_abstractBeforeTypedef(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@abstract \n * @typedef {Object<T,T>} */",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseAbstract_abstractBeforeType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"* @abstract \n * @type {Object<T,T>} */",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseAbstract_abstractBeforeInterface(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"* @abstract \n * @interface */",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseAbstract_abstractAndNotPrivate(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info1 := jt.parse("* @public \n * @abstract */")
	assertThat(t, info1.IsAbstract()).IsTrue()

	info2 := jt.parse("* @protected \n * @abstract */")
	assertThat(t, info2.IsAbstract()).IsTrue()

	info3 := jt.parse("* @package \n * @abstract */")
	assertThat(t, info3.IsAbstract()).IsTrue()
}

func TestParseAbstract_abstractAndPrivate(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"* @private \n * @abstract */",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseAbstract_abstractAndFinal(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"* @final \n * @abstract */",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestStackedAnnotation(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@const @type {string}*/")
	assertThat(t, info.IsConstant()).IsTrue()
	assertThat(t, info.HasType()).IsTrue()
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, info.GetType())
}

func TestStackedAnnotation2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@type {string} @const */")
	assertThat(t, info.IsConstant()).IsTrue()
	assertThat(t, info.HasType()).IsTrue()
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, info.GetType())
}

func TestStackedAnnotation3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@const @see {string}*/")
	assertThat(t, info.IsConstant()).IsTrue()
	assertThat(t, info.HasType()).IsFalse()
}

func TestStackedAnnotation4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@constructor @extends {Foo} @implements {Bar}*/")
	assertThat(t, info.IsConstructor()).IsTrue()
	assertThat(t, info.HasBaseType()).IsTrue()
	assertThat(t, info.GetImplementedInterfaceCount()).IsEqualTo(1)
}

func TestStackedAnnotation5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@param {number} x @constructor */")
	assertThat(t, info.HasParameterType("x")).IsTrue()
	assertThat(t, info.IsConstructor()).IsTrue()
}

func TestStackedAnnotation6(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@return {number} @constructor */", true)
	assertThat(t, info.HasReturnType()).IsTrue()
	assertThat(t, info.IsConstructor()).IsTrue()

	info = jt.parse("@return {number} @constructor */", false)
	assertThat(t, info.HasReturnType()).IsTrue()
	assertThat(t, info.IsConstructor()).IsTrue()
}

func TestStackedAnnotation7(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@return @constructor */", MISSING_TYPE_DECL_WARNING_TEXT)
	assertThat(t, info.HasReturnType()).IsTrue()
	assertThat(t, info.IsConstructor()).IsTrue()
}

func TestStackedAnnotation8(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@throws {number} @constructor */", true)
	assertThat(t, info.GetThrowsAnnotations()).IsNotEmpty()
	// The @constructor annotation gets swallowed as part of the throws annotation.
	assertThat(t, info.IsConstructor()).IsFalse()

	info = jt.parse("@return {number} @constructor */", false)
	assertThat(t, info.HasReturnType()).IsTrue()
	assertThat(t, info.IsConstructor()).IsTrue()
}

func TestStackedAnnotation9(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@const @private {string} */", true)
	assertThat(t, info.IsConstant()).IsTrue()
	assertThat(t, info.HasType()).IsTrue()
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, info.GetType())
}

func TestStackedAnnotation10(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@private @const {string} */", true)
	assertThat(t, info.GetVisibility()).IsEqualTo(rhino.JSDocInfo_Visibility_PRIVATE)
	assertThat(t, info.HasType()).IsTrue()
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, info.GetType())
	assertThat(t, info.IsConstant()).IsTrue()
}

func TestStackedAnnotation11(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@private @const */", true)
	assertThat(t, info.HasType()).IsFalse()
	assertThat(t, info.GetVisibility()).IsEqualTo(rhino.JSDocInfo_Visibility_PRIVATE)
	assertThat(t, info.IsConstant()).IsTrue()
}

func TestStackedAnnotation12(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@const @private */", true)
	assertThat(t, info.HasType()).IsFalse()
	assertThat(t, info.GetVisibility()).IsEqualTo(rhino.JSDocInfo_Visibility_PRIVATE)
	assertThat(t, info.IsConstant()).IsTrue()
}

func TestStackedAnnotation13(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@final @constructor */", true)
	assertThat(t, info.IsConstructor()).IsTrue()
	assertThat(t, info.IsFinal()).IsTrue()
}

func TestStackedAnnotation14(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@constructor @final */", true)
	assertThat(t, info.IsConstructor()).IsTrue()
	assertThat(t, info.IsFinal()).IsTrue()
}

func TestParsePreserve(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	comment := "@preserve Foo\nBar\n\nBaz*/"
	jt.parse(comment)
	assertThat(t, jt.prevLicense).IsEqualTo(" Foo\nBar\n\nBaz")
}

func TestParseLicense(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	comment := "@license Foo\nBar\n\nBaz*/"
	jt.parse(comment)
	assertThat(t, jt.prevLicense).IsEqualTo(" Foo\nBar\n\nBaz")
}

func TestParseLicenseAscii(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	comment := "@license Foo\n *   Bar\n\n  Baz*/"
	jt.parse(comment)
	assertThat(t, jt.prevLicense).IsEqualTo(" Foo\n   Bar\n\n  Baz")
}

func TestParseLicenseWithAnnotation(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	comment := "@license Foo \n * @author Charlie Brown */"
	jt.parse(comment)
	assertThat(t, jt.prevLicense).IsEqualTo(" Foo \n @author Charlie Brown ")
}

func TestParseLicenseMultiple(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	comment := "@license Foo @license Bar */"
	jt.parse(comment)
	assertThat(t, jt.prevLicense).IsEqualTo(" Foo @license Bar ")
}

func TestParseLicenseOnlyContentParsesLicenseWithInvalidJsDocTag(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	comment := "@date 11/12/2023\n@nosideeffects\n@deprecated\n@license Foo\n*/"
	r := jt.parse(comment, Config_JsDocParsing_LICENSE_COMMENTS_ONLY)
	assertThat(t, r).IsNotNull()
	assertThat(t, jt.prevLicense).IsEqualTo(" Foo\n")
	assertThat(t, r.IsNoSideEffects()).IsFalse()
	assertThat(t, r.IsDeprecated()).IsFalse()
	assertThat(t, markerAnnotationItems(r.GetMarkers())).ContainsExactlyElementsIn([]string{"license"})
}

func TestParseLicenseOnlyContentParsesPreservedCommentWithInvalidJsDocTag(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	comment := "@date 11/12/2023\n@nosideeffects\n@deprecated\n@preserve Foo\n*/"
	r := jt.parse(comment, Config_JsDocParsing_LICENSE_COMMENTS_ONLY)
	assertThat(t, r).IsNotNull()
	assertThat(t, jt.prevLicense).IsEqualTo(" Foo\n")
	assertThat(t, r.IsNoSideEffects()).IsFalse()
	assertThat(t, r.IsDeprecated()).IsFalse()
	assertThat(t, markerAnnotationItems(r.GetMarkers())).ContainsExactlyElementsIn([]string{"preserve"})
}

func TestCantParseLicenseWithInvalidJsDocTag(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	comment := "@date 11/12/2023\n@license Foo\n*/"
	jt.parse(
		comment,
		"illegal use of unknown JSDoc tag \"date\"; ignoring it. Place another character before the"+
			" @ to stop JSCompiler from parsing it as an annotation.")
	assertThat(t, jt.prevLicense).IsEqualTo(" Foo\n")
}

func TestParseDefine1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, jt.parse("@define {string}*/").GetType())
}

func TestParseDefine2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE,
		jt.parse("@define {string*/", "Bad type annotation. expected closing }"+JsDocInfoParser_BAD_TYPE_WIKI_LINK).GetType())
}

func TestParseDefine3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@define {boolean}*/")
	assertThat(t, info.IsConstant()).IsTrue()
	assertThat(t, info.IsDefine()).IsTrue()
	jt.AssertTypeEqualsExpr(jt.BOOLEAN_TYPE, info.GetType())
}

func TestParseDefine4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.NUMBER_TYPE, jt.parse("@define {number}*/").GetType())
}

func TestParseDefine5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.AssertTypeEqualsExpr(jt.CreateUnionType(jt.NUMBER_TYPE, jt.BOOLEAN_TYPE), jt.parse("@define {number|boolean}*/").GetType())
}

func TestParseDefineDescription(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	doc := jt.parse("@define {string} description of element \n next line*/", true)
	defineMarker := doc.GetMarkers()[0]
	assertThat(t, defineMarker.GetAnnotation().GetItem()).IsEqualTo("define")
	assertThat(t, defineMarker.GetDescription().GetItem()).Contains("description of element")
	assertThat(t, defineMarker.GetDescription().GetItem()).Contains("next line")
}

func TestParsePrivateDescription(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	doc := jt.parse("@private {string} description \n next line*/", true)
	defineMarker := doc.GetMarkers()[0]
	assertThat(t, defineMarker.GetAnnotation().GetItem()).IsEqualTo("private")
	assertThat(t, defineMarker.GetDescription().GetItem()).Contains("description ")
	assertThat(t, defineMarker.GetDescription().GetItem()).Contains("next line")
}

func TestParsePackagePrivateDescription(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	doc := jt.parse("@package {string} description \n next line */", true)
	defineMarker := doc.GetMarkers()[0]
	assertThat(t, defineMarker.GetAnnotation().GetItem()).IsEqualTo("package")
	assertThat(t, defineMarker.GetDescription().GetItem()).Contains("description ")
	assertThat(t, defineMarker.GetDescription().GetItem()).Contains("next line")
}

func TestParseProtectedDescription(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	doc := jt.parse("@protected {string} description \n next line*/", true)
	defineMarker := doc.GetMarkers()[0]
	assertThat(t, defineMarker.GetAnnotation().GetItem()).IsEqualTo("protected")
	assertThat(t, defineMarker.GetDescription().GetItem()).Contains("description ")
	assertThat(t, defineMarker.GetDescription().GetItem()).Contains("next line")
}

func TestParseDefineErrors1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@enum {string}\n @define {string} */", "conflicting @define tag")
}

func TestParseDefineErrors2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@define {string}\n@enum {string} */\n",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseDefineErrors3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@const\n @define {string} */", "conflicting @define tag")
}

func TestParseDefineErrors4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@type {string} \n @define {string} */", "conflicting @define tag")
}

func TestParseDefineErrors5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@return {string}\n @define {string} */", "conflicting @define tag")
}

func TestParseDefineErrors7(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@define {string}\n @const */").IsConstant()).IsTrue()
}

func TestParseDefineErrors8(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@define {string}\n @type {string} */",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseFinal1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@final*/").IsFinal()).IsTrue()
}

func TestParseFinal2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@final\n@final*/", "Bad type annotation. extra @final tag."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseFinal3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@final*/").IsConstant()).IsTrue()
}

func TestParseOverride1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@override*/").IsOverride()).IsTrue()
}

func TestParseOverride2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@override\n@override*/",
		"Bad type annotation. extra @override/@inheritDoc tag."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseInheritDoc1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@inheritDoc*/").IsOverride()).IsTrue()
}

func TestParseInheritDoc2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@override\n@inheritDoc*/",
		"Bad type annotation. extra @override/@inheritDoc tag."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseInheritDoc3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@inheritDoc\n@inheritDoc*/",
		"Bad type annotation. extra @override/@inheritDoc tag."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseDeprecated1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@deprecated*/").IsDeprecated()).IsTrue()
}

func TestParseDeprecated2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@deprecated\n@deprecated*/", "extra @deprecated tag")
}

func TestParseExport1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@export*/").IsExport()).IsTrue()
}

func TestParseExport2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@export\n@export*/", "extra @export tag")
}

func TestParseExterns1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parseFileOverview("@externs*/").IsExterns()).IsTrue()
}

func TestParseExterns2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFileOverview("@externs\n@externs*/", "extra @externs tag")
}

func TestParseExterns3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@externs*/")).IsNull()
}

func TestParseTypeSummary1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parseFileOverview("@typeSummary*/").IsTypeSummary()).IsTrue()
}

func TestParseTypeSummary2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFileOverview("@typeSummary\n@typeSummary*/", "extra @typeSummary tag")
}

func TestParseTypeSummary3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@typeSummary*/")).IsNull()
}

func TestParseNoCompile1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parseFileOverview("@nocompile*/").IsNoCompile()).IsTrue()
}

func TestParseNoCompile2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFileOverview("@nocompile\n@nocompile*/", "extra @nocompile tag")
}

func TestParseClosureUnawareCode1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t,
		jt.parseFileOverview(
			"@fileoverview\n@closureUnaware*/",
			"@closureUnaware annotation is not allowed in this compilation").IsClosureUnawareCode()).IsTrue()
}

func TestParseClosureUnawareCode2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFileOverview(
		"@closureUnaware\n@closureUnaware*/",
		"@closureUnaware annotation is not allowed in this compilation",
		"@closureUnaware annotation is not allowed in this compilation",
		"extra @closureUnaware tag")
}

func TestParseNoClosureUnawareCode(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parseFileOverview("@fileoverview*/").IsClosureUnawareCode()).IsFalse()
}

func TestParseClosureUnawareCode(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t,
		jt.parse(
			"@closureUnaware*/",
			"@closureUnaware annotation is not allowed in this compilation").IsClosureUnawareCode()).IsTrue()
}

func TestParseClosureUnawareCode_allowed(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.allowClosureUnawareCode = true
	assertThat(t, jt.parse("@closureUnaware*/").IsClosureUnawareCode()).IsTrue()
}

func TestParseClosureUnawareCode_fileoverview_allowed(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.allowClosureUnawareCode = true
	assertThat(t, jt.parseFileOverview("@fileoverview\n@closureUnaware*/").IsClosureUnawareCode()).IsTrue()
}

func TestParseClosureUnawareCode_doesNotCreateFileoverviewCommentStandalone(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.allowClosureUnawareCode = true
	assertThat(t, jt.parseFileOverview("@closureUnaware*/")).IsNull()
}

func TestParseNoCoverage_fileOverview(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parseFileOverview("@nocoverage*/").IsNoCoverage()).IsTrue()
}

func TestParseNoCoverage_extraTag(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFileOverview("@nocoverage\n@nocoverage*/", "extra @nocoverage tag")
}

func TestParseNoCoverage_nonFileOverview(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@nocoverage*/")).IsNull()
}

func TestParseNoDts1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	doc := jt.parse("@nodts*/", true)
	assertThat(t, doc.IsNoDts()).IsTrue()
}

func TestParseNoDts2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@nodts\n@nodts*/", Config_JsDocParsing_INCLUDE_DESCRIPTIONS_NO_WHITESPACE, "extra @nodts tag")
}

func TestBugAnnotation(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@bug */")
}

func TestDescriptionAnnotation(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@description */")
}

func TestRegression1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	comment :=
		"* @param {number} index the index of blah\n* @return {boolean} whatever\n* @private\n*/\n"

	info := jt.parse(comment)
	assertThat(t, info.GetParameterCount()).IsEqualTo(1)
	jt.AssertTypeEqualsExpr(jt.NUMBER_TYPE, info.GetParameterType("index"))
	jt.AssertTypeEqualsExpr(jt.BOOLEAN_TYPE, info.GetReturnType())
	assertThat(t, info.GetVisibility()).IsEqualTo(rhino.JSDocInfo_Visibility_PRIVATE)
}

func TestRegression2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	comment :=
		" * @return {boolean} whatever\n * but important\n *\n * @param {number} index the index of blah\n * some more comments here\n * @param name the name of the person\n *\n * @protected\n */\n"

	info := jt.parse(comment, MISSING_TYPE_DECL_WARNING_TEXT)

	assertThat(t, info.GetParameterCount()).IsEqualTo(2)
	jt.AssertTypeEqualsExpr(jt.NUMBER_TYPE, info.GetParameterType("index"))
	assertThat(t, info.GetParameterType("name")).IsNull()
	jt.AssertTypeEqualsExpr(jt.BOOLEAN_TYPE, info.GetReturnType())
	assertThat(t, info.GetVisibility()).IsEqualTo(rhino.JSDocInfo_Visibility_PROTECTED)
}

func TestRegression3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	comment :=
		" * @param mediaTag this specified whether the @media tag is ....\n *\n\n@public\n *\n\n **********\n * @const\n */\n"

	info := jt.parse(comment, MISSING_TYPE_DECL_WARNING_TEXT)
	assertThat(t, info.GetParameterCount()).IsEqualTo(1)
	assertThat(t, info.GetParameterType("mediaTag")).IsNull()
	assertThat(t, info.GetVisibility()).IsEqualTo(rhino.JSDocInfo_Visibility_PUBLIC)
	assertThat(t, info.IsConstant()).IsTrue()
}

func TestRegression4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	comment :=
		"* @const\n* @private\n* @constructor\n*/\n"

	info := jt.parse(comment)
	assertThat(t, info.IsConstant()).IsTrue()
	assertThat(t, info.IsDefine()).IsFalse()
	assertThat(t, info.IsConstructor()).IsTrue()
	assertThat(t, info.GetVisibility()).IsEqualTo(rhino.JSDocInfo_Visibility_PRIVATE)
}

func TestRegression5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	comment := "@const\n@enum {string}\n@public*/"

	info := jt.parse(comment)
	assertThat(t, info.IsConstant()).IsTrue()
	assertThat(t, info.IsDefine()).IsFalse()
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, info.GetEnumParameterType())
	assertThat(t, info.GetVisibility()).IsEqualTo(rhino.JSDocInfo_Visibility_PUBLIC)
}

func TestRegression7(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	comment :=
		"* @desc description here\n* @param {boolean} flag and some more description\n*     nicely formatted\n*/\n"

	info := jt.parse(comment)
	assertThat(t, info.GetParameterCount()).IsEqualTo(1)
	jt.AssertTypeEqualsExpr(jt.BOOLEAN_TYPE, info.GetParameterType("flag"))
	assertThat(t, info.GetDescription()).IsEqualTo("description here")
}

func TestRegression8(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	comment :=
		" * @name random tag here\n * @desc description here\n *\n * @param {boolean} flag and some more description\n *     nicely formatted\n */\n"

	info := jt.parse(comment)
	assertThat(t, info.GetParameterCount()).IsEqualTo(1)
	jt.AssertTypeEqualsExpr(jt.BOOLEAN_TYPE, info.GetParameterType("flag"))
	assertThat(t, info.GetDescription()).IsEqualTo("description here")
}

func TestRegression9(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"* @param {string} p0 blah blah blah\n*/\n")

	assertThat(t, jsdoc.GetBaseType()).IsNull()
	assertThat(t, jsdoc.IsConstant()).IsFalse()
	assertThat(t, jsdoc.GetDescription()).IsNull()
	assertThat(t, jsdoc.GetEnumParameterType()).IsNull()
	assertThat(t, jsdoc.GetParameterCount()).IsEqualTo(1)
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, jsdoc.GetParameterType("p0"))
	assertThat(t, jsdoc.GetReturnType()).IsNull()
	assertThat(t, jsdoc.GetType()).IsNull()
	assertThat(t, jsdoc.GetVisibility()).IsEqualTo(rhino.JSDocInfo_Visibility_INHERITED)
}

func TestRegression10(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"* @param {!String} p0 blah blah blah\n* @param {boolean} p1 fobar\n* @return {!Date} jksjkash dshad\n*/\n")

	assertThat(t, jsdoc.GetBaseType()).IsNull()
	assertThat(t, jsdoc.IsConstant()).IsFalse()
	assertThat(t, jsdoc.GetDescription()).IsNull()
	assertThat(t, jsdoc.GetEnumParameterType()).IsNull()
	assertThat(t, jsdoc.GetParameterCount()).IsEqualTo(2)
	jt.AssertTypeEqualsExpr(jt.STRING_OBJECT_TYPE, jsdoc.GetParameterType("p0"))
	jt.AssertTypeEqualsExpr(jt.BOOLEAN_TYPE, jsdoc.GetParameterType("p1"))
	jt.AssertTypeEqualsExpr(jt.DATE_TYPE, jsdoc.GetReturnType())
	assertThat(t, jsdoc.GetType()).IsNull()
	assertThat(t, jsdoc.GetVisibility()).IsEqualTo(rhino.JSDocInfo_Visibility_INHERITED)
}

func TestRegression11(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"* @constructor\n*/\n")

	assertThat(t, jsdoc.GetBaseType()).IsNull()
	assertThat(t, jsdoc.IsConstant()).IsFalse()
	assertThat(t, jsdoc.GetDescription()).IsNull()
	assertThat(t, jsdoc.GetEnumParameterType()).IsNull()
	assertThat(t, jsdoc.GetParameterCount()).IsEqualTo(0)
	assertThat(t, jsdoc.GetReturnType()).IsNull()
	assertThat(t, jsdoc.GetType()).IsNull()
	assertThat(t, jsdoc.GetVisibility()).IsEqualTo(rhino.JSDocInfo_Visibility_INHERITED)
}

func TestRegression12(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"* @extends FooBar\n*/\n")

	jt.AssertTypeEqualsExpr(jt.createNamedType("FooBar"), jsdoc.GetBaseType())
	assertThat(t, jsdoc.IsConstant()).IsFalse()
	assertThat(t, jsdoc.GetDescription()).IsNull()
	assertThat(t, jsdoc.GetEnumParameterType()).IsNull()
	assertThat(t, jsdoc.GetParameterCount()).IsEqualTo(0)
	assertThat(t, jsdoc.GetReturnType()).IsNull()
	assertThat(t, jsdoc.GetType()).IsNull()
	assertThat(t, jsdoc.GetVisibility()).IsEqualTo(rhino.JSDocInfo_Visibility_INHERITED)
}

func TestRegression13(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"* @type {!RegExp}\n* @protected\n*/\n")

	assertThat(t, jsdoc.GetBaseType()).IsNull()
	assertThat(t, jsdoc.IsConstant()).IsFalse()
	assertThat(t, jsdoc.GetDescription()).IsNull()
	assertThat(t, jsdoc.GetEnumParameterType()).IsNull()
	assertThat(t, jsdoc.GetParameterCount()).IsEqualTo(0)
	assertThat(t, jsdoc.GetReturnType()).IsNull()
	jt.AssertTypeEqualsExpr(jt.REGEXP_TYPE, jsdoc.GetType())
	assertThat(t, jsdoc.GetVisibility()).IsEqualTo(rhino.JSDocInfo_Visibility_PROTECTED)
}

func TestRegression14(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"* @const\n* @private\n*/\n")

	assertThat(t, jsdoc.GetBaseType()).IsNull()
	assertThat(t, jsdoc.IsConstant()).IsTrue()
	assertThat(t, jsdoc.GetDescription()).IsNull()
	assertThat(t, jsdoc.GetEnumParameterType()).IsNull()
	assertThat(t, jsdoc.GetParameterCount()).IsEqualTo(0)
	assertThat(t, jsdoc.GetReturnType()).IsNull()
	assertThat(t, jsdoc.GetType()).IsNull()
	assertThat(t, jsdoc.GetVisibility()).IsEqualTo(rhino.JSDocInfo_Visibility_PRIVATE)
}

func TestRegression15(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"* @desc Hello,\n* World!\n*/\n")

	assertThat(t, jsdoc.GetBaseType()).IsNull()
	assertThat(t, jsdoc.IsConstant()).IsTrue()
	assertThat(t, jsdoc.GetDescription()).IsEqualTo("Hello, World!")
	assertThat(t, jsdoc.GetEnumParameterType()).IsNull()
	assertThat(t, jsdoc.GetParameterCount()).IsEqualTo(0)
	assertThat(t, jsdoc.GetReturnType()).IsNull()
	assertThat(t, jsdoc.GetType()).IsNull()
	assertThat(t, jsdoc.GetVisibility()).IsEqualTo(rhino.JSDocInfo_Visibility_INHERITED)
	assertThat(t, jsdoc.IsExport()).IsFalse()
}

func TestRegression16(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"Email is plp@foo.bar\n@type {string}\n*/\n")

	assertThat(t, jsdoc.GetBaseType()).IsNull()
	assertThat(t, jsdoc.IsConstant()).IsFalse()
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, jsdoc.GetType())
	assertThat(t, jsdoc.GetParameterCount()).IsEqualTo(0)
	assertThat(t, jsdoc.GetReturnType()).IsNull()
	assertThat(t, jsdoc.GetVisibility()).IsEqualTo(rhino.JSDocInfo_Visibility_INHERITED)
}

func TestRegression17(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	// verifying that if no @desc is present the description is empty
	assertThat(t, jt.parse("@private*/").GetDescription()).IsNull()
}

func TestFullRegression1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull(
		"/** @param (string,number) foo*/function bar(foo){}",
		"Bad type annotation. expecting a variable name in a @param tag."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestFullRegression2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull(
		"/** @param {string,number) foo*/function bar(foo){}",
		"Bad type annotation. expected closing }"+JsDocInfoParser_BAD_TYPE_WIKI_LINK,
		"Bad type annotation. expecting a variable name in a @param tag."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestFullRegression3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFull("/**..\n*/")
}

func TestBug907488(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@type {number,null} */", "Bad type annotation. expected closing }"+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestBug907494(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@return {Object,undefined} */",
		"Bad type annotation. expected closing }"+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestBug909468(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@extends {(x)}*/", "Bad type annotation. expecting a type name."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseInterface(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@interface*/").IsInterface()).IsTrue()
}

func TestParseImplicitCast1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@type {string} \n * @implicitCast*/").IsImplicitCast()).IsTrue()
}

func TestParseImplicitCast2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@type {string}*/").IsImplicitCast()).IsFalse()
}

func TestParseDuplicateImplicitCast(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@type {string} \n * @implicitCast \n * @implicitCast*/",
		"Bad type annotation. extra @implicitCast tag."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseInterfaceDoubled(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"* @interface\n* @interface\n*/\n",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestInvalidTemplatedTypedef1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T \n * @typedef {Object<T,T>} */",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestInvalidTemplatedTypedef2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@typedef {Array<T>} \n * @template T */",
		"Bad type annotation. Type name(s) for @template annotation declared twice."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseImplements(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	interfaces :=
		jt.parse("@implements {SomeInterface}*/").GetImplementedInterfaces()
	assertThat(t, interfaces).HasSize(1)
	jt.AssertTypeEqualsExpr(jt.createNamedType("SomeInterface"), interfaces[0])
}

func TestParseImplementsTwo(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	interfaces :=
		jt.parse(
			"* @implements {SomeInterface1}\n* @implements {SomeInterface2}\n*/\n").GetImplementedInterfaces()
	assertThat(t, interfaces).HasSize(2)
	jt.AssertTypeEqualsExpr(jt.createNamedType("SomeInterface1"), interfaces[0])
	jt.AssertTypeEqualsExpr(jt.createNamedType("SomeInterface2"), interfaces[1])
}

func TestParseImplementsSameTwice(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"* @implements {Smth}\n* @implements {Smth}\n*/\n",
		"Bad type annotation. duplicate @implements tag."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseImplementsNoName(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"* @implements {} */", "Bad type annotation. expecting a type name."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseImplementsMissingRC(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("* @implements {Smth */", "Bad type annotation. expected closing }"+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseImplementsNullable1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@implements {Base?} */", "Bad type annotation. expected closing }"+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseImplementsNullable2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@implements Base? */",
		"Bad type annotation. expected end of line or comment."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseImplementsWithQualifier(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@implements {!Base}*/",
		"Bad type annotation. @implements/@extends requires a bare interface/record name without !"+
			" or ?."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestInterfaceExtends(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"* @interface \n* @extends {Extended} */\n")
	assertThat(t, jsdoc.IsInterface()).IsTrue()
	assertThat(t, jsdoc.GetExtendedInterfacesCount()).IsEqualTo(1)
	types := jsdoc.GetExtendedInterfaces()
	jt.AssertTypeEqualsExpr(jt.createNamedType("Extended"), types[0])
}

func TestInterfaceMultiExtends1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"* @interface \n* @extends {Extended1} \n* @extends {Extended2} */\n")
	assertThat(t, jsdoc.IsInterface()).IsTrue()
	assertThat(t, jsdoc.GetBaseType()).IsNull()
	assertThat(t, jsdoc.GetExtendedInterfacesCount()).IsEqualTo(2)
	types := jsdoc.GetExtendedInterfaces()
	jt.AssertTypeEqualsExpr(jt.createNamedType("Extended1"), types[0])
	jt.AssertTypeEqualsExpr(jt.createNamedType("Extended2"), types[1])
}

func TestInterfaceMultiExtends2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"* @extends {Extended1} \n* @interface \n* @extends {Extended2} \n* @extends {Extended3} */\n")
	assertThat(t, jsdoc.IsInterface()).IsTrue()
	assertThat(t, jsdoc.GetBaseType()).IsNull()
	assertThat(t, jsdoc.GetExtendedInterfacesCount()).IsEqualTo(3)
	types := jsdoc.GetExtendedInterfaces()
	jt.AssertTypeEqualsExpr(jt.createNamedType("Extended1"), types[0])
	jt.AssertTypeEqualsExpr(jt.createNamedType("Extended2"), types[1])
	jt.AssertTypeEqualsExpr(jt.createNamedType("Extended3"), types[2])
}

func TestBadClassMultiExtends(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"* @extends {Extended1} \n* @constructor \n* @extends {Extended2} */\n",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestBadExtendsWithNullable(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"@constructor\n * @extends {Object?} */",
			"Bad type annotation. expected closing }"+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
	assertThat(t, jsdoc.IsConstructor()).IsTrue()
	jt.AssertTypeEqualsExpr(jt.OBJECT_TYPE, jsdoc.GetBaseType())
}

func TestBadImplementsWithNullable(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"@implements {Disposable?}\n * @constructor */",
			"Bad type annotation. expected closing }"+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
	assertThat(t, jsdoc.IsConstructor()).IsTrue()
	jt.AssertTypeEqualsExpr(jt.createNamedType("Disposable"), jsdoc.GetImplementedInterfaces()[0])
}

func TestBadTypeDefInterfaceAndConstructor1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"@interface\n@constructor*/",
			"Bad type annotation. cannot be both an interface and a constructor."+
				JsDocInfoParser_BAD_TYPE_WIKI_LINK)
	assertThat(t, jsdoc.IsInterface()).IsTrue()
}

func TestBadTypeDefInterfaceAndConstructor2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"@constructor\n@interface*/",
			"Bad type annotation. cannot be both an interface and a constructor."+
				JsDocInfoParser_BAD_TYPE_WIKI_LINK)
	assertThat(t, jsdoc.IsConstructor()).IsTrue()
}

// test structural interface matching */
func TestBadTypeDefInterfaceAndStructuralTyping1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"@constructor\n@record*/",
			"Bad type annotation. conflicting @record tag."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
	assertThat(t, jsdoc.IsInterface()).IsFalse()
}

// test structural interface matching */
func TestBadTypeDefInterfaceAndStructuralTyping2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"@type{number}\n@record*/",
			"Bad type annotation. conflicting @record tag."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
	assertThat(t, jsdoc.IsInterface()).IsFalse()
}

// test structural interface matching */
func TestBadTypeDefInterfaceAndStructuralTyping3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"@type{{x:number}}\n@record*/",
			"Bad type annotation. conflicting @record tag."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
	assertThat(t, jsdoc.IsInterface()).IsFalse()
}

// test structural interface matching */
func TestBadTypeDefInterfaceAndStructuralTyping4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@interface\n@record*/",
		"Bad type annotation. conflicting @record tag."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

// test structural interface matching */
func TestBadTypeDefInterfaceAndStructuralTyping5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc := jt.parse("@record*/")
	assertThat(t, jsdoc.IsInterface()).IsTrue()
}

func TestDocumentationParameter(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc := jt.parse("@param {Number} number42 This is a description.*/", true)

	assertThat(t, jsdoc.HasDescriptionForParameter("number42")).IsTrue()
	assertThat(t, jsdoc.GetDescriptionForParameter("number42")).IsEqualTo("This is a description.")
}

func TestMultilineDocumentationParameter(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"@param {Number} number42 This is a description"+"\n* on multiple \n* lines.*/",
			true)

	assertThat(t, jsdoc.HasDescriptionForParameter("number42")).IsTrue()
	assertThat(t, jsdoc.GetDescriptionForParameter("number42")).IsEqualTo("This is a description on multiple lines.")
}

func TestDocumentationMultipleParameter(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"@param {Number} number42 This is a description.\n* @param {Integer} number87 This is another description.*/\n",
			true)

	assertThat(t, jsdoc.HasDescriptionForParameter("number42")).IsTrue()
	assertThat(t, jsdoc.GetDescriptionForParameter("number42")).IsEqualTo("This is a description.")

	assertThat(t, jsdoc.HasDescriptionForParameter("number87")).IsTrue()
	assertThat(t, jsdoc.GetDescriptionForParameter("number87")).IsEqualTo("This is another description.")
}

func TestDocumentationMultipleParameter2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"@param {number} delta = 0 results in a redraw\n  != 0 ..... */\n",
			true)
	assertThat(t, jsdoc.HasDescriptionForParameter("delta")).IsTrue()
	assertThat(t, jsdoc.GetDescriptionForParameter("delta")).IsEqualTo("= 0 results in a redraw != 0 .....")
}

func TestAuthors(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"@param {Number} number42 This is a description.\n* @param {Integer} number87 This is another description.\n* @author a@google.com (A Person)\n* @author b@google.com (B Person)\n* @author c@google.com (C Person)*/\n",
			true)

	authors := jsdoc.GetAuthors()

	assertThat(t, authors).IsNotNull()
	assertThat(t, authors).HasSize(3)

	assertThat(t, authors).Contains("a@google.com (A Person)")
	assertThat(t, authors).Contains("b@google.com (B Person)")
	assertThat(t, authors).Contains("c@google.com (C Person)")
}

func TestSuppress1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@suppress {x} */")
	assertThat(t, info.GetSuppressions()).IsEqualTo(setOf("x"))
}

func TestSuppressFollowedByAnnotationOnSameLine_parseTypesOnly(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	sameLine :=
		jt.parse("@suppress {visibility} @return {boolean} */", Config_JsDocParsing_TYPES_ONLY)
	assertThat(t, sameLine.GetReturnType()).IsNotNull()
	assertThat(t, sameLine.GetSuppressions()).Contains("visibility")
}

func TestSuppressFollowedByAnnotationOnSameLine_parseWithDescription(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	// Regression test - the compiler used to include the @return {boolean} annotation in the
	// suppress description only if JsDocParsing.INCLUDE_DESCRIPTIONS_NO_WHITESPACE was used.
	sameLine :=
		jt.parse(
			"@suppress {visibility} @return {boolean} */",
			Config_JsDocParsing_INCLUDE_DESCRIPTIONS_NO_WHITESPACE)
	assertThat(t, sameLine.GetReturnType()).IsNotNull()
	assertThat(t, sameLine.GetSuppressions()).Contains("visibility")
}

func TestSuppressWithDescriptionFollowedByAnnotation_parseTypeOnly(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	sameLine :=
		jt.parse(
			"@suppress {visibility} this is a description @return {boolean} */",
			Config_JsDocParsing_TYPES_ONLY)
	// The @return {boolean} is considered part of the @suppress description.
	assertThat(t, sameLine.GetReturnType()).IsNull()
	assertThat(t, sameLine.GetSuppressions()).Contains("visibility")
}

func TestSuppressWithDescriptionFollowedByAnnotation_includeDescriptions(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	sameLine :=
		jt.parse(
			"@suppress {visibility} this is a description @return {boolean} */",
			Config_JsDocParsing_INCLUDE_DESCRIPTIONS_NO_WHITESPACE)
	// The @return {boolean} is considered part of the @suppress description.
	assertThat(t, sameLine.GetReturnType()).IsNull()
	assertThat(t, sameLine.GetSuppressions()).Contains("visibility")
}

func TestSuppressWithDescription(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsDocComment :=
		// No description for `@suppress {z}`.
		"@suppress {x,y} Some description.\n * @suppress {z}\n*/\n"
	info := jt.parse(jsDocComment /* parseDocumentation= */, true)
	assertThat(t, info.GetSuppressions()).IsEqualTo(setOf("x", "y", "z"))
	assertThat(t, info.GetSuppressionsAndTheirDescription()).ContainsEntry(setOf("x", "y"), "Some description.")
	assertThat(t, info.GetSuppressionsAndTheirDescription()).ContainsEntry(setOf("z"), "")
}

func TestSuppressWithDescription_multipleLines(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsDocComment :=
		// No description for `@suppress {z}`.
		"@suppress {x,y} Some description.\n * Spans across lines.\n * @suppress {z}\n*/\n"
	info := jt.parse(jsDocComment /* parseDocumentation= */, true)
	assertThat(t, info.GetSuppressions()).IsEqualTo(setOf("x", "y", "z"))
	assertThat(t, info.GetSuppressionsAndTheirDescription()).ContainsEntry(setOf("x", "y"), "Some description. Spans across lines.")
	assertThat(t, info.GetSuppressionsAndTheirDescription()).ContainsEntry(setOf("z"), "")
}

func TestSuppressWithDescription_multipleLines_withOtherAnnotations(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsDocComment :=
		// No description for `@suppress {z}`.
		"@author XYZ\n * @suppress {x,y} Some description.\n * Spans across lines.\n * @override\n * @suppress {z}\n*/\n"
	info := jt.parse(jsDocComment /* parseDocumentation= */, true)
	assertThat(t, info.GetAuthors()).Contains("XYZ")
	assertThat(t, info.IsOverride()).IsTrue()
	assertThat(t, info.GetSuppressions()).IsEqualTo(setOf("x", "y", "z"))
	assertThat(t, info.GetSuppressionsAndTheirDescription()).ContainsEntry(setOf("x", "y"), "Some description. Spans across lines.")
	assertThat(t, info.GetSuppressionsAndTheirDescription()).ContainsEntry(setOf("z"), "")
}

// When exact warning is repeated in another @suppress annotation, JSDocInfoParser does not save
// the second annotation.
func TestSuppressWithDescription_repeatedWarnings1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsDocComment :=
		"@suppress {x} Some description.\n * @suppress {x} Another description.\n*/\n"
	info := jt.parse(jsDocComment /* parseDocumentation= */, true)
	assertThat(t, info.GetSuppressions()).IsEqualTo(setOf("x"))
	assertThat(t, info.GetSuppressionsAndTheirDescription()).ContainsExactly(setOf("x"), "Some description.")
}

// A different @suppress annotation containing a non-repeated warning, we store both annotations.
func TestSuppressWithDescription_repeatedWarnings2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsDocComment :=
		"@suppress {x} Some description.\n * @suppress {x, y} Another description.\n*/\n"
	info := jt.parse(jsDocComment /* parseDocumentation= */, true)
	assertThat(t, info.GetSuppressions()).IsEqualTo(setOf("x", "y"))
	assertThat(t, info.GetSuppressionsAndTheirDescription()).HasSize(2)
	assertThat(t, info.GetSuppressionsAndTheirDescription()).ContainsEntry(setOf("x"), "Some description.")
	assertThat(t, info.GetSuppressionsAndTheirDescription()).ContainsEntry(setOf("x", "y"), "Another description.")
}

func TestSuppress2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@suppress {x|y|x|z} */")
	assertThat(t, info.GetSuppressions()).IsEqualTo(setOf("x", "y", "z"))
}

func TestSuppress3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@suppress {x,y} */")
	assertThat(t, info.GetSuppressions()).IsEqualTo(setOf("x", "y"))
}

func TestSuppressMultiline(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsDocComment :=
		`@suppress {
 *   x, y,
 *   z
 * } Auto-added to unblock migration
 */
`
	info := jt.parse(jsDocComment /* parseDocumentation= */, true, "malformed @suppress tag")
	assertThat(t, info.GetSuppressions()).IsEmpty()
	assertThat(t, info.GetSuppressionsAndTheirDescription()).IsEmpty()
}

func TestSuppressMultilineWithPipesAndLeadingSeparator(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsDocComment :=
		`@suppress {x
 *   | y
 *   , z}
 */
`
	info := jt.parse(jsDocComment, "malformed @suppress tag")
	assertThat(t, info).IsNull()
}

func TestJsDocAfterSuppress(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@suppress {x,y} @const {string} */")
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, info.GetType())
}

func TestMultipleSuppressTags(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@suppress {x} \n * @suppress {y} */")
	assertThat(t, info.GetSuppressions()).IsEqualTo(setOf("x", "y"))
}

func TestBadSuppress1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@suppress {} */", "malformed @suppress tag")
}

func TestBadSuppress2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@suppress {x|} */", "malformed @suppress tag")
}

func TestBadSuppress3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@suppress {|x} */", "malformed @suppress tag")
}

func TestBadSuppress4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@suppress {x|y */", "malformed @suppress tag")
}

func TestBadSuppress7(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@suppress {impossible} */", "unknown @suppress parameter: impossible")
}

func TestBadSuppress8(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@suppress */", "malformed @suppress tag")
}

func TestBadSuppressMultilineTrailingComma(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@suppress {x,\n * } */", "malformed @suppress tag")
}

func TestMods1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parseFileOverview("@mods {ns.Foo} */")
	assertThat(t, info.GetMods()).IsEqualTo("ns.Foo")
}

func TestMods2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parseFileOverview("@mods {google3.path.to.file} */")
	assertThat(t, info.GetMods()).IsEqualTo("google3.path.to.file")
}

func TestJsDocAfterMods(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parseFileOverview("@mods {ns.Foo} @modName {new_feature_flag} */")
	assertThat(t, info.GetMods()).IsEqualTo("ns.Foo")
}

func TestMultipleModsTags(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFileOverview("@mods {ns.Foo} \n * @mods {ns.Bar} */", "extra @mods tag")
}

func TestBadMods1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFileOverview("@mods {} */", "malformed @mods tag")
}

func TestBadMods2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFileOverview("@mods { */", "malformed @mods tag")
}

func TestBadMods3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFileOverview("@mods } */", "malformed @mods tag")
}

func TestBadMods4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFileOverview("@mods ns.Foo */", "malformed @mods tag")
}

func TestBadMods5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFileOverview("@mods {ns.Foo */", "malformed @mods tag")
}

func TestBadMods6(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFileOverview("@mods ns.Foo} */", "malformed @mods tag")
}

func TestModifies1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@modifies {this} */")
	assertThat(t, info.GetModifies()).IsEqualTo(setOf("this"))
}

func TestModifies2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@modifies {arguments} */")
	assertThat(t, info.GetModifies()).IsEqualTo(setOf("arguments"))
}

func TestModifies3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@modifies {this|arguments} */")
	assertThat(t, info.GetModifies()).IsEqualTo(setOf("this", "arguments"))
}

func TestModifies4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@param {*} x\n * @modifies {x} */")
	assertThat(t, info.GetModifies()).IsEqualTo(setOf("x"))
}

func TestModifies5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info :=
		jt.parse(
			"@param {*} x\n * @param {*} y\n * @modifies {x} */\n")
	assertThat(t, info.GetModifies()).IsEqualTo(setOf("x"))
}

func TestModifies6(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info :=
		jt.parse(
			"@param {*} x\n * @param {*} y\n * @modifies {x|y} */\n")
	assertThat(t, info.GetModifies()).IsEqualTo(setOf("x", "y"))
}

func TestBadModifies1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@modifies {} */", "malformed @modifies tag")
}

func TestBadModifies2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@modifies {this|} */", "malformed @modifies tag")
}

func TestBadModifies3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@modifies {|this} */", "malformed @modifies tag")
}

func TestBadModifies4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@modifies {this|arguments */", "malformed @modifies tag")
}

func TestBadModifies5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@modifies {this,arguments} */", "malformed @modifies tag")
}

func TestBadModifies6(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@modifies {this} \n * @modifies {this} */", "conflicting @modifies tag")
}

func TestBadModifies7(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@modifies {impossible} */", "unknown @modifies parameter: impossible")
}

func TestBadModifies8(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@modifies {this}\n@nosideeffects */\n",
		"@nosideeffects functions cannot have @modifies (modifying arguments/this is a side"+
			" effect)")
}

func TestBadModifies9(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@nosideeffects @modifies {this} */",
		"@nosideeffects functions cannot have @modifies (modifying arguments/this is a side"+
			" effect)")
}

func TestConflictingThrowsAndNoSideEffects1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@throws {Error}\n@nosideeffects */\n",
		"@nosideeffects functions cannot have @throws (throwing is a side effect)")
}

func TestConflictingThrowsAndNoSideEffects2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@nosideeffects\n@throws {Error} */\n",
		"@nosideeffects functions cannot have @throws (throwing is a side effect)")
}

func TestDuplicateNoSideEffects(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@nosideeffects\n@nosideeffects */\n",
		"conflicting @nosideeffects tag")
}

func TestThrowsAndModifiesCompatible(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info :=
		jt.parse(
			"@modifies {this}\n@throws {Error} Some error */\n",
			true)
	assertThat(t, info.ModifiesThis()).IsTrue()
	assertThat(t, info.GetThrowsAnnotations()).ContainsExactly("{Error} Some error")

	info =
		jt.parse(
			"@throws {Error} Some error\n@modifies {this} */\n",
			true)
	assertThat(t, info.ModifiesThis()).IsTrue()
	assertThat(t, info.GetThrowsAnnotations()).ContainsExactly("{Error} Some error")
}

func TestNoParseFileOverview(t *testing.T) {
	t.Skip("@Ignore in JsDocInfoParserTest.java")
	jt := newJsDocInfoParserTest(t)
	jsdoc := jt.parseFileOverviewWithoutDoc("@fileoverview Hi mom! */")
	assertThat(t, jsdoc.GetFileOverview()).IsNull()
	assertThat(t, jsdoc.HasFileOverview()).IsTrue()
}

func TestFileOverviewSingleLine(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc := jt.parseFileOverview("@fileoverview Hi mom! */")
	assertThat(t, jsdoc.GetFileOverview()).IsEqualTo("Hi mom!")
}

func TestFileOverviewMultiLine(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc := jt.parseFileOverview("@fileoverview Pie is \n * good! */")
	assertThat(t, jsdoc.GetFileOverview()).IsEqualTo("Pie is\ngood!")
}

func TestFileOverviewDuplicate(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFileOverview("@fileoverview Pie \n * @fileoverview Cake */", "extra @fileoverview tag")
}

func TestPublicVisibilityAllowedInFileOverview(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFileOverview("@fileoverview \n * @public */")
}

func TestPackageVisibilityAllowedInFileOverview(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFileOverview("@fileoverview \n * @package */")
}

func TestImplicitVisibilityAllowedInFileOverview(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFileOverview("@fileoverview */")
}

func TestProtectedVisibilityNotAllowedInFileOverview(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFileOverview(
		"@fileoverview \n * @protected */",
		"protected visibility not allowed in @fileoverview block")
	jt.parseFileOverview(
		"@fileoverview @protected */", "protected visibility not allowed in @fileoverview block")
}

func TestPrivateVisibilityNotAllowedInFileOverview(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parseFileOverview(
		"@fileoverview \n @private */", "private visibility not allowed in @fileoverview block")
	jt.parseFileOverview(
		"@fileoverview @private */", "private visibility not allowed in @fileoverview block")
}

func TestReferences(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"@see A cool place!\n* @see The world.\n* @see SomeClass#SomeMember\n* @see A boring test case*/\n",
			true)

	references := jsdoc.GetReferences()

	assertThat(t, references).IsNotNull()
	assertThat(t, references).HasSize(4)

	assertThat(t, references).Contains("A cool place!")
	assertThat(t, references).Contains("The world.")
	assertThat(t, references).Contains("SomeClass#SomeMember")
	assertThat(t, references).Contains("A boring test case")
}

func TestSingleTags(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"@deprecated In favor of the new one!\n* @return {SomeType} The most important object :-)*/\n",
			true)

	assertThat(t, jsdoc.IsDeprecated()).IsTrue()
	assertThat(t, jsdoc.GetDeprecationReason()).IsEqualTo("In favor of the new one!")
	assertThat(t, jsdoc.GetReturnDescription()).IsEqualTo("The most important object :-)")
}

func TestSingleTags2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc := jt.parse("@param {SomeType} a The most important object :-)*/", true)

	assertThat(t, jsdoc.GetDescriptionForParameter("a")).IsEqualTo("The most important object :-)")
}

func TestSingleTagsReordered(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"@deprecated In favor of the new one!\n * @return {SomeType} The most important object :-)*/\n",
			true)

	assertThat(t, jsdoc.IsDeprecated()).IsTrue()
	assertThat(t, jsdoc.GetDeprecationReason()).IsEqualTo("In favor of the new one!")
	assertThat(t, jsdoc.GetReturnDescription()).IsEqualTo("The most important object :-)")
}

func TestAuthorMissing(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("* @author */", true, "@author tag missing author")
}

func TestSeeMissing(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("* @see */", true, "@see tag missing description")
}

func TestParseBlockComment(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"this is a nice comment\n * that is multiline\n * @author abc@google.com */\n",
			true)

	assertThat(t, jsdoc.GetBlockDescription()).IsEqualTo("this is a nice comment\nthat is multiline")

	jt.assertDocumentationInMarker(
		jt.assertAnnotationMarker(jsdoc, "author", 2, 3), "abc@google.com", 10, 2, 24)
}

func TestParseBlockComment2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"this is a nice comment\n * that is *** multiline\n * @author abc@google.com */\n",
			true)

	assertThat(t, jsdoc.GetBlockDescription()).IsEqualTo("this is a nice comment\nthat is *** multiline")

	jt.assertDocumentationInMarker(
		jt.assertAnnotationMarker(jsdoc, "author", 2, 3), "abc@google.com", 10, 2, 24)
}

func TestParseBlockComment3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			" \n* hello world \n* @author abc@google.com */\n",
			true)

	assertThat(t, jsdoc.GetBlockDescription()).IsEqualTo("hello world")

	jt.assertDocumentationInMarker(
		jt.assertAnnotationMarker(jsdoc, "author", 2, 2), "abc@google.com", 9, 2, 23)
}

func TestParseWithMarkers1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc := jt.parse("@author abc@google.com */", true)

	jt.assertDocumentationInMarker(
		jt.assertAnnotationMarker(jsdoc, "author", 0, 0), "abc@google.com", 7, 0, 21)
}

func TestParseWithMarkers2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc := jt.parse("@param {Foo} somename abc@google.com */", true)

	jt.assertDocumentationInMarker(
		jt.assertAnnotationMarker(jsdoc, "param", 0, 0), "abc@google.com", 21, 0, 37)
}

func TestParseWithMarkers3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"@return {Foo} some long\n * multiline\n * description */\n",
			true)

	returnDoc := jt.assertAnnotationMarker(jsdoc, "return", 0, 0)
	jt.assertDocumentationInMarker(returnDoc, "some long multiline description", 14, 2, 15)
	assertThat(t, returnDoc.GetType().GetPositionOnStartLine()).IsEqualTo(8)
	assertThat(t, returnDoc.GetType().GetPositionOnEndLine()).IsEqualTo(12)
}

func TestParseWithMarkers4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"@author foobar\n * @param {Foo} somename abc@google.com */\n",
			true)

	jt.assertAnnotationMarker(jsdoc, "author", 0, 0)
	jt.assertAnnotationMarker(jsdoc, "param", 1, 3)
}

func TestParseWithMarkers5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"@return some long\n * multiline\n * description */\n",
			true,
			MISSING_TYPE_DECL_WARNING_TEXT)

	jt.assertDocumentationInMarker(
		jt.assertAnnotationMarker(jsdoc, "return", 0, 0), "some long multiline description", 8, 2, 15)
}

func TestParseWithMarkers6(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"@param x some long\n * multiline\n * description */\n",
			true,
			MISSING_TYPE_DECL_WARNING_TEXT)

	jt.assertDocumentationInMarker(
		jt.assertAnnotationMarker(jsdoc, "param", 0, 0), "some long multiline description", 8, 2, 15)
}

func TestParseWithMarkerNames1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc := jt.parse("@param {SomeType} name somedescription */", true)

	jt.assertNameInMarker(jt.assertAnnotationMarker(jsdoc, "param", 0, 0), "name", 0, 18)
}

func TestParseWithMarkerNames2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"@param {SomeType} name somedescription\n* @param {AnotherType} anothername des */\n",
			true)

	jt.assertTypeInMarker(
		jt.assertNameInMarker(jt.assertAnnotationMarker(jsdoc, "param", 0, 0, 0), "name", 0, 18),
		"SomeType",
		0,
		7,
		0,
		16,
		true)

	jt.assertTypeInMarker(
		jt.assertNameInMarker(jt.assertAnnotationMarker(jsdoc, "param", 1, 2, 1), "anothername", 1, 23),
		"AnotherType",
		1,
		9,
		1,
		21,
		true)
}

func TestParseWithMarkerNames3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc := jt.parse("@param {Some.Long.Type.\n *  Name} name somedescription */", true)

	jt.assertTypeInMarker(
		jt.assertNameInMarker(jt.assertAnnotationMarker(jsdoc, "param", 0, 0, 0), "name", 1, 10),
		"Some.Long.Type.Name",
		0,
		7,
		1,
		8,
		true)
}

func TestParseWithoutMarkerName(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc := jt.parse("@author helloworld*/", true)
	assertThat(t, jt.assertAnnotationMarker(jsdoc, "author", 0, 0).GetNameNode()).IsNull()
}

func TestParseWithMarkerType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc := jt.parse("@extends {FooBar}*/", true)

	jt.assertTypeInMarker(jt.assertAnnotationMarker(jsdoc, "extends", 0, 0), "FooBar", 0, 9, 0, 16, true)
}

func TestParseWithMarkerType2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc := jt.parse("@extends FooBar*/", true)

	jt.assertTypeInMarker(
		jt.assertAnnotationMarker(jsdoc, "extends", 0, 0), "FooBar", 0, 9, 0, 15, false)
}

func TestTypeTagConflict1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@constructor \n * @constructor */",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestTypeTagConflict2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@interface \n * @interface */",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestTypeTagConflict3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@constructor \n * @interface */",
		"Bad type annotation. cannot be both an interface and a constructor."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestTypeTagConflict4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@interface \n * @constructor */",
		"Bad type annotation. cannot be both an interface and a constructor."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestTypeTagConflict5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@interface \n * @type {string} */",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestTypeTagConflict6(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@typedef {string} \n * @type {string} */",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestTypeTagConflict7(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@typedef {string} \n * @constructor */",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestTypeTagConflict8(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@typedef {string} \n * @return {boolean} */",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestTypeTagConflict9(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@enum {string} \n * @return {boolean} */",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestTypeTagConflict10(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@this {Object} \n * @enum {boolean} */",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestTypeTagConflict11(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@param {Object} x \n * @type {boolean} */",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestTypeTagConflict12(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@typedef {boolean} \n * @param {Object} x */",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestTypeTagConflict13(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@typedef {boolean} \n * @extends {Object} */",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestTypeTagConflict14(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@return {?} x \n * @return {?} y */",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestTypeTagConflict15(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"/**\n * @struct\n * @struct\n */\nfunction StrStr() {}\n",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestTypeTagConflict16(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"/**\n * @struct\n * @interface\n */\nfunction StrIntf() {}\n")
}

func TestTypeTagConflict17(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"/**\n * @interface\n * @struct\n */\nfunction StrIntf() {}\n")
}

func TestTypeTagConflict18(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"/**\n * @dict\n * @dict\n */\nfunction DictDict() {}\n",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestTypeTagConflict19(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"/**\n * @dict\n * @interface\n */\nfunction DictDict() {}\n")
}

func TestTypeTagConflict20(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"/**\n * @interface\n * @dict\n */\nfunction DictDict() {}\n")
}

func TestTypeTagConflict21(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"/**\n * @private {string}\n * @type {number}\n */\nfunction DictDict() {}\n",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestTypeTagConflict22(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"/**\n * @protected {string}\n * @param {string} x\n */\nfunction DictDict(x) {}\n",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestTypeTagConflict23(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"/**\n * @public {string}\n * @return {string} x\n */\nfunction DictDict() {}\n",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestTypeTagConflict24(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"/**\n * @const {string}\n * @return {string} x\n */\nfunction DictDict() {}\n",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestTypeTagConflict25(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"/**\n * @package {string}\n * @return {string} x\n */\nfunction DictDict() {}\n",
		"Bad type annotation. type annotation incompatible with other annotations."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestPackageType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc := jt.parse("@package {string} */")
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, jsdoc.GetType())
}

func TestPrivateType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc := jt.parse("@private {string} */")
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, jsdoc.GetType())
}

func TestProtectedType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc := jt.parse("@protected {string} */")
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, jsdoc.GetType())
}

func TestPublicType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc := jt.parse("@public {string} */")
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, jsdoc.GetType())
}

func TestConstType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc := jt.parse("@const {string} */")
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, jsdoc.GetType())
}

func TestFinalType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc := jt.parse("@final {string} */")
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, jsdoc.GetType())
}

func TestExportType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc := jt.parse("@export {string} descr\n next line */", true)
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, jsdoc.GetType())

	assertThat(t, jsdoc.IsExport()).IsTrue()

	defineMarker := jsdoc.GetMarkers()[0]
	assertThat(t, defineMarker.GetAnnotation().GetItem()).IsEqualTo("export")
	assertThat(t, defineMarker.GetDescription().GetItem()).Contains("descr")
	assertThat(t, defineMarker.GetDescription().GetItem()).Contains("next line")
}

func TestMixedVisibility(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@public @private */", "extra visibility tag")
	jt.parse("@public @protected */", "extra visibility tag")
	jt.parse("@export @protected */")
	jt.parse("@export {string}\n * @private */")
	jt.parse("@export {string}\n * @public */")
}

func TestStableIdGenerator(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info :=
		jt.parse(
			"/**\n * @idGenerator {stable}\n */\nfunction getId() {}\n")
	assertThat(t, info.IsStableIdGenerator()).IsTrue()
}

func TestStableIdGeneratorConflict(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"/**\n * @idGenerator {stable}\n * @idGenerator {stable}\n */\nfunction getId() {}\n",
		"extra @idGenerator tag")
}

func TestXidGenerator(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info :=
		jt.parse(
			"/**\n * @idGenerator {xid}\n */\nfunction getId() {}\n")
	assertThat(t, info.IsXidGenerator()).IsTrue()
}

func TestXidGeneratorConflict(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"/**\n * @idGenerator {xid}\n * @idGenerator {xid}\n */\nfunction getId() {}\n",
		"extra @idGenerator tag")
}

func TestIdGenerator(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info :=
		jt.parse(
			"/**\n * @idGenerator\n */\nfunction getId() {}\n")
	assertThat(t, info.IsIdGenerator()).IsTrue()
}

func TestIdGeneratorConflict(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"/**\n * @idGenerator\n * @idGenerator\n */\nfunction getId() {}\n",
		"extra @idGenerator tag")
}

func TestIdGenerator1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@idGenerator {unique} */")
	assertThat(t, info.IsIdGenerator()).IsTrue()
}

func TestIdGenerator2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@idGenerator {consistent} */")
	assertThat(t, info.IsConsistentIdGenerator()).IsTrue()
}

func TestIdGenerator3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@idGenerator {stable} */")
	assertThat(t, info.IsStableIdGenerator()).IsTrue()
}

func TestIdGenerator4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@idGenerator {xid} */")
	assertThat(t, info.IsXidGenerator()).IsTrue()
}

func TestIdGenerator5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@idGenerator {mapped} */")
	assertThat(t, info.IsMappedIdGenerator()).IsTrue()
}

func TestBadIdGenerator1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@idGenerator {} */", "malformed @idGenerator tag")
}

func TestBadIdGenerator2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@idGenerator {impossible} */", "unknown @idGenerator parameter: impossible")
}

func TestBadIdGenerator3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@idGenerator {unique */", "malformed @idGenerator tag")
}

func TestParserWithTemplateTypeNameMissing(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template */",
		"name not recognized due to syntax error.",
		"Bad type annotation. @template tag missing type name."+JsDocInfoParser_BAD_TYPE_WIKI_LINK,
		"Unexpected end of file")
}

func TestParserWithTwoTemplates(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T,V */")
}

func TestParserWithInvalidTemplateType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template {T} */",
		"Bounded generic semantics are currently still in development",
		"name not recognized due to syntax error.",
		"Bad type annotation. @template tag missing type name."+JsDocInfoParser_BAD_TYPE_WIKI_LINK,
		"Unexpected end of file")
}

func TestBoundedGenericWithoutName(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template {number} */",
		"Bounded generic semantics are currently still in development",
		"name not recognized due to syntax error.",
		"Bad type annotation. @template tag missing type name."+JsDocInfoParser_BAD_TYPE_WIKI_LINK,
		"Unexpected end of file")
}

func TestParserWithValidTemplateType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template {string} T */", //
		"Bounded generic semantics are currently still in development")
}

func TestParserWithValidAndInvalidTemplateType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template S, {T} */", //
		"name not recognized due to syntax error.")
}

func TestParserWithTemplateDuplicated(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T\n@template V */")
}

func TestParserWithTemplateDuplicated2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T,R\n@template V,U */")
}

func TestParserWithTemplateDuplicated3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T\n@param {string} x\n@template V */")
}

func TestParserWithTemplateTypeNameDeclaredTwice(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T\n@template T */",
		"Bad type annotation. Type name(s) for @template annotation declared twice."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParserWithTemplateTypeNameDeclaredTwice2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := S =: \n @template T := R =:*/",
		"Bad type annotation. Type name(s) for @template annotation declared twice."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParserWithTemplateTypeNameDeclaredTwice3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T \n @template T := R =:*/",
		"Bad type annotation. Type name(s) for @template annotation declared twice."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParserWithTemplateTypeNameDeclaredTwice4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := R =: \n @template T*/",
		"Bad type annotation. Type name(s) for @template annotation declared twice."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParserWithDoubleTemplateDeclaration2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T,T */",
		"Bad type annotation. Type name(s) for @template annotation declared twice."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParserWithTemplateDuplicatedTypeNameMissing(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T,R\n@template */",
		"name not recognized due to syntax error.",
		"Bad type annotation. @template tag missing type name."+JsDocInfoParser_BAD_TYPE_WIKI_LINK,
		"Unexpected end of file")
}

// This is unusual, but allowed. */
func TestParserWithLowercaseTemplateTypeName(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template t */")
	jt.parse("@template s,t */")
	jt.parse("@template key,value */")
}

func TestParserWithTemplateTypeNameWithDollar(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@template A,B$,C */")
	assertThat(t, info.GetTemplateTypeNames()).ContainsExactly("A", "B$", "C")
}

func TestParserWithTypeTransformationNewline(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template R := \n 'string' =:*/")
}

func TestParserWithTypeTransformation(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := 'string' =:*/")
}

func TestParserWithTypeTransformation2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := 'string' =:\nRandom text*/\n")
}

func TestParserWithTypeTransformationMultipleNames(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T, R := 'string' =:*/",
		"Bad type annotation. Multiple template names cannot be declared with bounds or TTL."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParserWithMissingTypeTransformationExpression(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := */",
		"Bad type annotation. Expected end delimiter for a type transformation."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK,
		"Bad type annotation. Missing type transformation expression."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParserWithMissingTypeTransformationExpression2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := =:*/",
		"Bad type annotation. Missing type transformation expression."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestBug16129690(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@param {T} x\n@template T\nrandom documentation text*/\n")
}

func TestTTLLineNoCharNo(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info :=
		jt.parse(
			"Some text on line 0\nMore text! This is line 1\n@template T := foo =:*/\n")
	assertThat(t, info.GetTypeTransformations()).HasSize(1)
	n := info.GetTypeTransformations().Get("T")
	assertNode(t, n).HasLineno(2) // lineno of "@" in "@template"
	assertNode(t, n).HasCharno(0) // charno of "@" in "@template"
}

func TestMultilineTTLLineNoCharNo(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info :=
		jt.parse(
			"Some text on line 0\nMore text! This is line 1\n@template T :=\n  foo\n=:*/\n")
	assertThat(t, info.GetTypeTransformations()).HasSize(1)
	n := info.GetTypeTransformations().Get("T")
	assertNode(t, n).HasLineno(2) // lineno of "@" in "@template"
	assertNode(t, n).HasCharno(0) // charno of "@" in "@template"
}

func TestParserWithTTLInvalidOperation(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := foo() =:*/", "Bad type annotation. Invalid type transformation expression")
}

func TestParserWithTTLInvalidTypeTransformation(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := var a; =:*/",
		"Bad type annotation. Invalid type transformation expression")
}

func TestParserWithTTLValidTypename(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := foo =:*/")
}

func TestParserWithTTLValidTypename2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := R =:*/")
}

func TestParserWithTTLValidTypename3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := _Foo =:*/")
}

func TestParserWithTTLValidTypename4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := $foo =:*/")
}

func TestParserWithTTLBasicType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := 'string' =:*/")
}

func TestParserWithTTLValidUnionType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := union('string', 'number') =:*/")
}

func TestParserWithTTLValidUnionType2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := union(R, S) =:*/")
}

func TestParserWithTTLValidUnionType3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := union(R, 'string', S) =:*/")
}

func TestParserWithTTLEmptyUnionType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := union() =:*/", "Bad type annotation. Missing parameter in union")
}

func TestParserWithTTLSingletonUnionType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := union('string') =:*/", "Bad type annotation. Missing parameter in union")
}

func TestParserWithTTLInvalidUnionType2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := union(function(a){}, T) =:*/",
		"Bad type annotation. Invalid type transformation expression",
		"Bad type annotation. Invalid expression inside union type")
}

func TestParserWithNestedUnionFirstParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := union(union(N, 'null'), S) =:*/")
}

func TestParserWithNestedUnionSecondParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := union(N, union('null', S)) =:*/")
}

func TestParserWithNestedBooleanFirstParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond( eq(cond(eq(N, N), 'string', 'number'), 'string'),\n'string',\n'number') =: */\n")
}

func TestParserWithNestedBooleanSecondParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond( eq('string', cond(eq(N, N), 'string', 'number')),\n'string',\n'number') =:*/\n")
}

func TestParserWithTTLConditional(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := cond(eq(T, R), R, S) =: */")
}

func TestParserWithTTLConditional2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := cond(sub(T, R), R, S) =: */")
}

func TestParserWithTTLConditionalStringEquivalence(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := cond(streq(R, S), R, S) =: */")
}

func TestParserWithTTLConditionalStringEquivalence2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := cond(streq(R, 'foo'), R, S) =: */")
}

func TestParserWithTTLConditionalStringEquivalence3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := cond(streq('foo', 'bar'), R, S) =: */")
}

func TestParserWithTTLConditionalIsConstructor(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := cond(isCtor(R), R, S) =: */")
}

func TestParserWithTTLConditionalIsConstructor2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := cond(isCtor('foo'), R, S) =: */")
}

func TestParserWithTTLConditionalIsTemplatized(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := cond(isTemplatized(R), R, S) =: */")
}

func TestParserWithTTLConditionalIsTemplatized2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := cond(isTemplatized('foo'), R, S) =: */")
}

func TestParserWithTTLConditionalIsRecord(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := cond(isRecord(R), R, S) =: */")
}

func TestParserWithTTLConditionalAndOperation(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := cond(isCtor(R) && isCtor(S), R, S) =: */")
}

func TestParserWithTTLConditionalOrOperation(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := cond(isCtor(R) || isCtor(S), R, S) =: */")
}

func TestParserWithTTLConditionalNotOperation(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := cond(!isCtor(R), R, S) =: */")
}

func TestParserWithTTLConditionalNestedBoolOperation(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := cond(((isCtor(R) || isCtor(S)) && !isCtor(R)), R, S) =: */")
}

func TestParserWithTTLConditionalIsRecord2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := cond(isRecord('foo'), R, S) =: */")
}

func TestParserWithTTLConditionalIsDefined(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := cond(isDefined(R), R, S) =: */")
}

func TestParserWithTTLConditionalIsUnknown(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := cond(isUnknown(R), R, S) =: */")
}

func TestParserWithTTLConditionalStringEquivalenceInvalidParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(streq('foo', foo()), R, S) =: */",
		"Bad type annotation. Invalid string",
		"Bad type annotation. Invalid expression inside boolean",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLConditionalStringEquivalenceInvalidParamEmptyStr(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(streq('', S), R, S) =: */",
		"Bad type annotation. Invalid string parameter",
		"Bad type annotation. Invalid expression inside boolean",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLConditionalIsConstructorInvalidParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(isCtor(foo()), R, S) =: */",
		"Bad type annotation. Invalid type transformation expression",
		"Bad type annotation. Invalid expression inside boolean",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLConditionalIsTemplatizedInvalidParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(isTemplatized(foo()), R, S) =: */",
		"Bad type annotation. Invalid type transformation expression",
		"Bad type annotation. Invalid expression inside boolean",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLConditionalIsRecordInvalidParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(isRecord(foo()), R, S) =: */",
		"Bad type annotation. Invalid type transformation expression",
		"Bad type annotation. Invalid expression inside boolean",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLConditionalIsDefinedInvalidParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(isDefined('foo'), R, S) =: */",
		"Bad type annotation. Invalid name",
		"Bad type annotation. Invalid expression inside boolean",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLConditionalIsUnknownInvalidParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(isUnknown(foo()), R, S) =: */",
		"Bad type annotation. Invalid type transformation expression",
		"Bad type annotation. Invalid expression inside boolean",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLConditionalAndInvalidParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond('foo' && isCtor(R), R, S) =: */",
		"Bad type annotation. Invalid boolean expression",
		"Bad type annotation. Invalid expression inside boolean",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLConditionalAndInvalidParam2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(isCtor(R) && 'foo', R, S) =: */",
		"Bad type annotation. Invalid boolean expression",
		"Bad type annotation. Invalid expression inside boolean",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLConditionalOrInvalidParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond('foo' || isCtor(R), R, S) =: */",
		"Bad type annotation. Invalid boolean expression",
		"Bad type annotation. Invalid expression inside boolean",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLConditionalOrInvalidParam2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(isCtor(R) || 'foo', R, S) =: */",
		"Bad type annotation. Invalid boolean expression",
		"Bad type annotation. Invalid expression inside boolean",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLConditionalNotInvalidParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(!'foo', R, S) =: */",
		"Bad type annotation. Invalid boolean expression",
		"Bad type annotation. Invalid expression inside boolean",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLExtraParamBoolean(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(eq(T, R, S), R, S) =: */",
		"Bad type annotation. Found extra parameter in eq",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLExtraParamStringEq(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(streq(T, R, S), R, S) =: */",
		"Bad type annotation. Found extra parameter in streq",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLExtraParamIsConstructor(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(isCtor(T, R, S), R, S) =: */",
		"Bad type annotation. Found extra parameter in isCtor",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLExtraParamIsTemplatized(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(isTemplatized(T, R, S), R, S) =: */",
		"Bad type annotation. Found extra parameter in isTemplatized",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLExtraParamIsRecord(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(isRecord(T, R), R, S) =: */",
		"Bad type annotation. Found extra parameter in isRecord",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLExtraParamIsDefined(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(isDefined(T, R), R, S) =: */",
		"Bad type annotation. Found extra parameter in isDefined",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLExtraParamIsUnknown(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(isUnknown(T, R), R, S) =: */",
		"Bad type annotation. Found extra parameter in isUnknown",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLMissingParamBoolean(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(eq(T), R, S) =: */",
		"Bad type annotation. Missing parameter in eq",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLMissingParamStringEquivalence(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(streq(T), R, S) =: */",
		"Bad type annotation. Missing parameter in streq",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLMissingParamIsConstructor(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(isCtor(), R, S) =: */",
		"Bad type annotation. Missing parameter in isCtor",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLMissingParamIsTemplatized(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(isTemplatized(), R, S) =: */",
		"Bad type annotation. Missing parameter in isTemplatized",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLMissingParamIsRecord(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(isRecord(), R, S) =: */",
		"Bad type annotation. Missing parameter in isRecord",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLMissingParamIsDefined(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(isDefined(), R, S) =: */",
		"Bad type annotation. Missing parameter in isDefined",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLMissingParamIsUnknown(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(isUnknown(), R, S) =: */",
		"Bad type annotation. Missing parameter in isUnknown",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLInvalidBooleanConditional(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(aaa, R, S) =: */",
		"Bad type annotation. Invalid boolean expression",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLInvalidBooleanConditional2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(foo(T, R), S, R) =: */",
		"Bad type annotation. Invalid boolean predicate",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLInvalidBooleanConditional3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(eq(T, foo()), R, S) =: */",
		"Bad type annotation. Invalid type transformation expression",
		"Bad type annotation. Invalid expression inside boolean",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLInvalidConditionalMissingParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(sub(T, R), T) =: */",
		"Bad type annotation. Missing parameter in cond")
}

func TestParserWithTTLInvalidConditionalExtraParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(sub(T, R), T, R, R) =: */",
		"Bad type annotation. Found extra parameter in cond")
}

func TestParserWithTTLInvalidConditional(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(eq(T, R), foo(), S) =: */",
		"Bad type annotation. Invalid type transformation expression",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLInvalidConditional2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := cond(eq(T, R), S, foo()) =: */",
		"Bad type annotation. Invalid type transformation expression",
		"Bad type annotation. Invalid expression inside conditional")
}

func TestParserWithTTLValidMapunion(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := mapunion(T, (S) => S) =: */")
}

func TestParserWithTTLValidMapunion2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := mapunion(union('string', 'number'), (S) => S) =: */")
}

func TestParserWithTTLInvalidMapunionType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := mapunion(foo(), (S) => S) =: */",
		"Bad type annotation. Invalid type transformation expression",
		"Bad type annotation. Invalid expression inside mapunion")
}

func TestParserWithTTLInvalidMapunionFn(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := mapunion(R, S) =: */",
		"Bad type annotation. Invalid map function",
		"Bad type annotation. Invalid expression inside mapunion")
}

func TestParserWithTTLInvalidMapunionMissingParams(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := mapunion(T) =: */", "Bad type annotation. Missing parameter in mapunion")
}

func TestParserWithTTLInvalidMapunionExtraParams(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := mapunion(T, (S) => S, R) =: */",
		"Bad type annotation. Found extra parameter in mapunion")
}

func TestParserWithTTLInvalidMapunionMissingFnParams(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := mapunion(T, () => S) =: */",
		"Bad type annotation. Missing parameter in map function",
		"Bad type annotation. Invalid expression inside mapunion")
}

func TestParserWithTTLInvalidMapunionExtraFnParams(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := mapunion(T, (S, R) => S) =: */",
		"Bad type annotation. Found extra parameter in map function",
		"Bad type annotation. Invalid expression inside mapunion")
}

func TestParserWithTTLInvalidMapunionFunctionBody(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := mapunion(T, (S) => foo()) =: */",
		"Bad type annotation. Invalid type transformation expression",
		"Bad type annotation. Invalid expression inside map function body")
}

func TestParserWithTTLUseCaseObject(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := mapunion(T, (x) => cond(eq(x, 'string'), 'String',\ncond(eq(x, 'number'), 'Number',\ncond(eq(x, 'boolean'), 'Boolean',\ncond(eq(x, 'null'), 'Object',\ncond(eq(x, 'undefined'), 'Object',\nx)))))) =: */\n")
}

func TestParserWithTTLNoneType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := none() =: */")
}

func TestParserWithTTLNoneType2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := cond(eq(S, none()), S, T) =: */")
}

func TestParserWithTTLInvalidNoneType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := none(foo) =: */", "Bad type annotation. Found extra parameter in none")
}

func TestParserWithTTLInvalidNoneType2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := none(a, b, c) =: */", "Bad type annotation. Found extra parameter in none")
}

func TestParserWithTTLAllType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := all() =: */")
}

func TestParserWithTTLAllType2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := cond(eq(S, all()), S, T) =: */")
}

func TestParserWithTTLInvalidAllType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := all(foo) =: */", "Bad type annotation. Found extra parameter in all")
}

func TestParserWithTTLInvalidAllType2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := all(a, b, c) =: */", "Bad type annotation. Found extra parameter in all")
}

func TestParserWithTTLUnknownType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := unknown() =: */")
}

func TestParserWithTTLUnknownType2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := cond(eq(S, unknown()), S, T) =: */")
}

func TestParserWithTTLInvalidUnknownType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := unknown(foo) =: */",
		"Bad type annotation. Found extra parameter in unknown")
}

func TestParserWithTTLInvalidUnknownType2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := unknown(a, b, c) =: */",
		"Bad type annotation. Found extra parameter in unknown")
}

func TestParserWithTTLTemplateTypeOperation(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := type('Map', 'string', 'number') =: */")
}

func TestParserWithTTLTemplateTypeOperationGeneric(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := type('Array', T) =: */")
}

func TestParserWithTTLTemplateTypeOperationGeneric2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := type(T, R) =: */")
}

func TestParserWithTTLTemplateTypeOperationNestedGeneric(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := type(T, type(R, S)) =: */")
}

func TestParserWithTTLTemplateTypeOperationGenericWithUnion(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := type(T, union(R, S)) =: */")
}

func TestParserWithTTLInvalidTemplateTypeOperationGenericUnion(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := type(union(R, S), T) =: */",
		"Bad type annotation. Invalid type name or type variable",
		"Bad type annotation. Invalid expression inside template type operation")
}

func TestParserWithTTLInvalidTypeOperationNestedGeneric(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := type(T, foo()) =: */",
		"Bad type annotation. Invalid type transformation expression",
		"Bad type annotation. Invalid expression inside template type operation")
}

func TestParserWithTTLValidRawTypeOperation(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := rawTypeOf(T) =: */")
}

func TestParserWithTTLValidRawTypeOperation2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := rawTypeOf(type(T, R)) =: */")
}

func TestParserWithTTLInvalidRawTypeOperation(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := rawTypeOf(foo()) =: */",
		"Bad type annotation. Invalid type transformation expression",
		"Bad type annotation. Invalid expression inside rawTypeOf")
}

func TestParserWithTTLInvalidRawTypeOperationExtraParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := rawTypeOf(R, S) =: */",
		"Bad type annotation. Found extra parameter in rawTypeOf")
}

func TestParserWithTTLInvalidRawTypeOperationMissingParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := rawTypeOf() =: */", "Bad type annotation. Missing parameter in rawTypeOf")
}

func TestParserWithTTLNestedRawTypeOperation(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := rawTypeOf(type(T, rawTypeOf(type(R, S)))) =: */")
}

func TestParserWithTTLValidTemplateTypeOfOperation(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := templateTypeOf(T, 1) =: */")
}

func TestParserWithTTLValidTemplateTypeOfOperation2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := templateTypeOf(type(T, R), 1) =: */")
}

func TestParserWithTTLInvalidFirstParamTemplateTypeOf(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := templateTypeOf(foo(), 1) =: */",
		"Bad type annotation. Invalid type transformation expression",
		"Bad type annotation. Invalid expression inside templateTypeOf")
}

func TestParserWithTTLInvalidSecondParamTemplateTypeOf(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := templateTypeOf(R, foo()) =: */",
		"Bad type annotation. Invalid index",
		"Bad type annotation. Invalid expression inside templateTypeOf")
}

func TestParserWithTTLInvalidSecondParamTemplateTypeOf2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := templateTypeOf(R, 1.5) =: */",
		"Bad type annotation. Invalid index",
		"Bad type annotation. Invalid expression inside templateTypeOf")
}

func TestParserWithTTLInvalidSecondParamTemplateTypeOf3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := templateTypeOf(R, -1) =: */",
		"Bad type annotation. Invalid index",
		"Bad type annotation. Invalid expression inside templateTypeOf")
}

func TestParserWithTTLInvalidTemplateTypeOfExtraParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := templateTypeOf(R, 1, S) =: */",
		"Bad type annotation. Found extra parameter in templateTypeOf")
}

func TestParserWithTTLInvalidTemplateTypeOfMissingParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := templateTypeOf() =: */",
		"Bad type annotation. Missing parameter in templateTypeOf")
}

func TestParserWithTTLInvalidTemplateTypeOfMissingParam2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := templateTypeOf(T) =: */",
		"Bad type annotation. Missing parameter in templateTypeOf")
}

func TestParserWithTTLNestedTemplateTypeOfOperation(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := templateTypeOf(templateTypeOf(type(T, type(R, S)), 0),0) =: */")
}

func TestParserWithTTLValidPrintType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := printType('msg', R) =: */")
}

func TestParserWithTTLInvalidFirstParamPrintType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := printType(foo(), R) =: */",
		"Bad type annotation. Invalid message",
		"Bad type annotation. Invalid expression inside printType")
}

func TestParserWithTTLInvalidSecondParamPrintType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := printType('msg', foo()) =: */",
		"Bad type annotation. Invalid type transformation expression",
		"Bad type annotation. Invalid expression inside printType")
}

func TestParserWithTTLInvalidPrintTypeExtraParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := printType(R, S, U) =: */",
		"Bad type annotation. Found extra parameter in printType")
}

func TestParserWithTTLInvalidPrintTypeOfMissingParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := printType() =: */", "Bad type annotation. Missing parameter in printType")
}

func TestParserWithTTLInvalidPrintTypeOfMissingParam2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := printType(R) =: */", "Bad type annotation. Missing parameter in printType")
}

func TestParserWithTTLValidPropType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := propType('p', R) =: */")
}

func TestParserWithTTLInvalidFirstParamPropType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := propType(foo(), R) =: */",
		"Bad type annotation. Invalid property name",
		"Bad type annotation. Invalid expression inside propType")
}

func TestParserWithTTLInvalidSecondParamPropType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := propType('msg', foo()) =: */",
		"Bad type annotation. Invalid type transformation expression",
		"Bad type annotation. Invalid expression inside propType")
}

func TestParserWithTTLInvalidPropTypeExtraParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := propType(R, S, U) =: */",
		"Bad type annotation. Found extra parameter in propType")
}

func TestParserWithTTLInvalidPropTypeOfMissingParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := propType() =: */", "Bad type annotation. Missing parameter in propType")
}

func TestParserWithTTLInvalidPropTypeOfMissingParam2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := propType(R) =: */", "Bad type annotation. Missing parameter in propType")
}

func TestParserWithTTLRecordType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := record({prop:T}) =: */")
}

func TestParserWithTTLNestedRecordType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := record({prop: record({p1:'number', p2:'boolean'}), x:'string' })\n=: */\n")
}

func TestParserWithTTLInvalidRecordTypeMissingParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := record() =: */", "Bad type annotation. Missing parameter in record")
}

func TestParserWithTTLMergeRecords(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := record({x:'number'}, {y:'number'}) =: */")
}

func TestParserWithTTLInvalidMergeRecords(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := record({x:'number'}, foo()) =: */",
		"Bad type annotation. Invalid type transformation expression",
		"Bad type annotation. Invalid expression inside record")
}

func TestParserWithTTLInvalidRecordTypeWithInvalidTypeInProperty(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := record({prop:foo()}) =: */",
		"Bad type annotation. Invalid type transformation expression",
		"Bad type annotation. Invalid expression inside record")
}

func TestParserWithTTLInvalidRecordTypeMissingTypeInProperty(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := record({prop}) =: */",
		"Bad type annotation. Invalid property, missing type",
		"Bad type annotation. Invalid expression inside record")
}

func TestParserWithTTLInvalidRecordTypeInvalidRecordExpression(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := record(foo()) =: */",
		"Bad type annotation. Invalid type transformation expression",
		"Bad type annotation. Invalid expression inside record")
}

func TestParserWithTTLRecordTypeTypeVars(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := record(T, R) =: */")
}

func TestParserWithTTLEmptyRecordType(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := record({}) =: */")
}

func TestParserWithTTLTypeTransformationInFirstParamMapunion(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := mapunion(templateTypeOf(type(R, union(S, U)), 0), (x) => x) =: */\n")
}

func TestParserWithTTLValidMaprecord(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := maprecord(R, (K, V) => V) =: */")
}

func TestParserWithTTLValidMaprecord2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := maprecord(record({x:'string', y:'number'}), (K, V) => V) =: */")
}

func TestParserWithTTLInvalidMaprecordFirstParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := maprecord(foo(), (K, V) => V) =: */",
		"Bad type annotation. Invalid type transformation expression",
		"Bad type annotation. Invalid expression inside maprecord")
}

func TestParserWithTTLInvalidMaprecordNotAFunction(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := maprecord(R, S) =: */",
		"Bad type annotation. Invalid map function",
		"Bad type annotation. Invalid expression inside maprecord")
}

func TestParserWithTTLInvalidMaprecordMissingParams(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := maprecord(R) =: */", "Bad type annotation. Missing parameter in maprecord")
}

func TestParserWithTTLInvalidMaprecordExtraParams(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := maprecord(R, (K, V) => V, R) =: */",
		"Bad type annotation. Found extra parameter in maprecord")
}

func TestParserWithTTLInvalidMaprecordMissingParamsInMapFunction(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := maprecord(R, () => S) =: */",
		"Bad type annotation. Missing parameter in map function",
		"Bad type annotation. Invalid expression inside maprecord")
}

func TestParserWithTTLInvalidMaprecordMissingParamsInMapFunction2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := maprecord(R, (K) => S) =: */",
		"Bad type annotation. Missing parameter in map function",
		"Bad type annotation. Invalid expression inside maprecord")
}

func TestParserWithTTLInvalidMaprecordExtraParamsInMapFunction(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := maprecord(R, (K, V, S) => S) =: */",
		"Bad type annotation. Found extra parameter in map function",
		"Bad type annotation. Invalid expression inside maprecord")
}

func TestParserWithTTLInvalidMaprecordInvalidFunctionBody(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := maprecord(R, (K, V) => foo()) =: */",
		"Bad type annotation. Invalid type transformation expression",
		"Bad type annotation. Invalid expression inside map function body")
}

func TestParserWithTTLTypeOfVar(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := typeOfVar('x') =: */")
}

func TestParserWithTTLTypeOfVar2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := typeOfVar('x.y') =: */")
}

func TestParserWithTTLTypeOfVarInvalidName(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := typeOfVar(foo()) =: */",
		"Bad type annotation. Invalid name",
		"Bad type annotation. Invalid expression inside typeOfVar")
}

func TestParserWithTTLTypeOfVarMissingParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := typeOfVar() =: */", "Bad type annotation. Missing parameter in typeOfVar")
}

func TestParserWithTTLTypeOfVarExtraParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := typeOfVar(a, b) =: */",
		"Bad type annotation. Found extra parameter in typeOfVar")
}

func TestParserWithTTLInstanceOf(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := instanceOf('x') =: */")
}

func TestParserWithTTLInstanceOf2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := instanceOf(R) =: */")
}

func TestParserWithTTLInstanceOfInvalidName(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := instanceOf(foo()) =: */",
		"Bad type annotation. Invalid type transformation expression",
		"Bad type annotation. Invalid expression inside instanceOf")
}

func TestParserWithTTLInstanceOfMissingParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := instanceOf() =: */",
		"Bad type annotation. Missing parameter in instanceOf")
}

func TestParserWithTTLInstanceOfExtraParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := instanceOf(a, b) =: */",
		"Bad type annotation. Found extra parameter in instanceOf")
}

func TestParserWithTTLNativeTypeExprBasic(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := typeExpr('string') =: */")
}

func TestParserWithTTLNativeTypeExprBasic2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := typeExpr('goog.ui.Menu') =: */")
}

func TestParserWithTTLNativeTypeExprUnion(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := typeExpr('number|boolean') =: */")
}

func TestParserWithTTLNativeTypeExprRecord(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := typeExpr('{myNum: number, myObject} ') =: */")
}

func TestParserWithTTLNativeTypeExprNullable(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := typeExpr('?number') =: */")
}

func TestParserWithTTLNativeTypeExprNonNullable(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := typeExpr('!Object') =: */")
}

func TestParserWithTTLNativeTypeExprFunction(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := typeExpr('function(string, boolean)') =: */")
}

func TestParserWithTTLNativeTypeExprFunctionReturn(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := typeExpr('function(): number') =: */")
}

func TestParserWithTTLNativeTypeExprFunctionThis(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := typeExpr('function(this:goog.ui.Menu, string)') =: */")
}

func TestParserWithTTLNativeTypeExprFunctionNew(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := typeExpr('function(new:goog.ui.Menu, string)') =: */")
}

func TestParserWithTTLNativeTypeExprFunctionVarargs(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := typeExpr('function(string, ...number): number') =: */")
}

func TestParserWithTTLNativeTypeExprFunctionOptional(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := typeExpr('function(?string=, number=)') =: */")
}

func TestParserWithTTLNativeTypeExprMissingParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@template T := typeExpr() =: */", "Bad type annotation. Missing parameter in typeExpr")
}

func TestParserWithTTLNativeTypeExprExtraParam(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := typeExpr('a', 'b') =: */",
		"Bad type annotation. Found extra parameter in typeExpr")
}

func TestParserWithTTLNativeInvalidTypeExpr(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template T := typeExpr(foo) =: */",
		"Bad type annotation. Invalid native type expression",
		"Bad type annotation. Invalid expression inside typeExpr")
}

func TestParserWithTTLAsynchUseCase(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template R := cond(eq(T, 'Object'),\nmaprecord(T,\n(K, V) => cond(eq(rawTypeOf(V), 'Promise'),\ntemplateTypeOf(V, 0),\n'undefined')\n),\nT)\n=: */\n")
}

func TestAllowlistedNewAnnotations(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@foobar */",
		"illegal use of unknown JSDoc tag \"foobar\"; ignoring it. Place another character before"+
			" the @ to stop JSCompiler from parsing it as an annotation.")
	jt.extraAnnotations.Add("foobar")
	jt.parse("@foobar */")
}

func TestEscapingIllegalAnnotations(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	// no error when illegal @annotation is preceded by a character.
	jt.parse("_@foobar */")
	jt.parse("\\@foobar */")
}

func TestAllowlistedConflictingAnnotation(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.extraAnnotations.Add("param")
	info := jt.parse("@param {number} index */")
	jt.AssertTypeEqualsExpr(jt.NUMBER_TYPE, info.GetParameterType("index"))
}

func TestJsDocAfterStatic(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@static @type {number} */")
	jt.AssertTypeEqualsExpr(jt.NUMBER_TYPE, info.GetType())
}

func TestNonIdentifierAnnotation(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	// Try to allowlist an annotation that is not a valid JS identifier.
	// It should not work.
	jt.extraAnnotations.Add("123")
	jt.parse(
		"@123 */",
		"illegal use of unknown JSDoc tag \"\"; ignoring it. Place another character before the @"+
			" to stop JSCompiler from parsing it as an annotation.")
}

func TestUnsupportedJsDocSyntax1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info := jt.parse("@param {string} [accessLevel=\"author\"] The user level */", true)
	assertThat(t, info.GetParameterCount()).IsEqualTo(1)
	jt.AssertTypeEqualsExpr(jt.Registry.CreateOptionalType(jt.STRING_TYPE), info.GetParameterType("accessLevel"))
	assertThat(t, info.GetDescriptionForParameter("accessLevel")).IsEqualTo("The user level")
}

func TestUnsupportedJsDocSyntax2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info :=
		jt.parse(
			"@param {?} userInfo The user info.\n * @param {?} userInfo.name The name of the user */\n",
			true,
			"invalid param name \"userInfo.name\"")
	assertThat(t, info.GetParameterCount()).IsEqualTo(1)
	assertThat(t, info.GetDescriptionForParameter("userInfo")).IsEqualTo("The user info.")
}

func TestAllowlistedAnnotations(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"* @addon \n* @augments \n* @base \n* @borrows \n* @bug \n* @channel \n* @class \n* @config \n* @constructs \n* @default \n* @description \n* @enhance \n* @enhanceable \n* @event \n* @example \n* @exception \n* @exec \n* @externs \n* @field \n* @function \n* @delcall \n* @deltemplate \n* @hassoydelcall \n* @hassoydeltemplate \n* @id \n* @ignore \n* @jsx \n* @jsxFrag \n* @inner \n* @lends {string} \n* @link \n* @mayhaveextraedge \n* @member \n* @memberOf \n* @modName \n* @name \n* @namespace \n* @ngInject \n* @nocompile \n* @pintomodule \n* @property \n* @requirecss \n* @requires \n* @since \n* @soyTemplate \n* @soyModule \n* @static \n* @supported\n* @wizaction \n* @wizcallback \n*/\n")
}

// TestJsDocInfoPosition (JsDocInfoParserTest.testJsDocInfoPosition) is in
// js_doc_info_parser_position_test.go: it reads positions through a jscomp
// SourceFile, so it is in the external test package parsing_test.

// JsDocInfoParserTest_ParseFull is parseFull(code) on a fresh test
// instance, for the external test package.
//
// Go-only: test-helper — exported to parsing_test (see above).
func JsDocInfoParserTest_ParseFull(t *testing.T, code string) rhino.Node {
	t.Helper()
	return newJsDocInfoParserTest(t).parseFull(code)
}

func TestGetOriginalCommentString(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	comment := "* @desc This is a comment */"
	info := jt.parse(comment)
	assertThat(t, info.GetOriginalCommentString()).IsNull()
	info = jt.parse(comment /* parseDocumentation= */, true)
	assertThat(t, info.GetOriginalCommentString()).IsEqualTo(comment)
}

func TestGetMultilineDesc(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	comment :=
		" * @desc description here\n         continued here\n */\n"
	info := jt.parse(comment)
	assertThat(t, info.GetDescription()).IsEqualTo("description here continued here")
}

func TestParseNgInject1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@ngInject*/").IsNgInject()).IsTrue()
}

func TestParseNgInject2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@ngInject \n@ngInject*/", "extra @ngInject tag")
}

func TestParsePolymerBehavior(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@polymerBehavior*/").IsPolymerBehavior()).IsTrue()
}

func TestParsePolymerBehaviorExtra(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@polymerBehavior \n@polymerBehavior*/", "extra @polymerBehavior tag")
}

func TestParsePolymer(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@polymer*/").IsPolymer()).IsTrue()
}

func TestParsePolymerExtra(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@polymer \n@polymer*/", "extra @polymer tag")
}

func TestParseCustomElement(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@customElement*/").IsCustomElement()).IsTrue()
}

func TestParseCustomElementExtra(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@customElement \n@customElement*/", "extra @customElement tag")
}

func TestParseMixinClass(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@mixinClass*/").IsMixinClass()).IsTrue()
}

func TestParseMixinClassExtra(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@mixinClass \n@mixinClass*/", "extra @mixinClass tag")
}

func TestParseMixinFunction(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@mixinFunction*/").IsMixinFunction()).IsTrue()
}

func TestParseMixinFunctionExtra(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@mixinFunction \n@mixinFunction*/", "extra @mixinFunction tag")
}

func TestParseSassGeneratedCssTs(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@sassGeneratedCssTs*/").IsSassGeneratedCssTs()).IsTrue()
}

func TestParseSassGeneratedCssTsExtra(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@sassGeneratedCssTs \n@sassGeneratedCssTs*/", "extra @sassGeneratedCssTs tag")
}

func TestParseWizaction1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@wizaction*/").IsWizaction()).IsTrue()
}

func TestParseWizaction2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@wizaction \n@wizaction*/", "extra @wizaction tag")
}

func TestParseProvideAlreadyProvided(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@provideAlreadyProvided*/").IsProvideAlreadyProvided()).IsTrue()
}

func TestParseProvideAlreadyProvidedExtra(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@provideAlreadyProvided \n@provideAlreadyProvided*/", "extra @provideAlreadyProvided tag")
}

func TestParseWizcallback(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	assertThat(t, jt.parse("@wizcallback*/").IsWizcallback()).IsTrue()
}

func TestParseWizcallbackExpectDuplicateError(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@wizcallback \n@wizcallback*/", "extra @wizcallback tag")
}

func TestTextExtents(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@return {@code foo} bar \n *    baz. */",
		true,
		"Bad type annotation. type not recognized due to syntax error."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestParseDuplicateNoCollapse(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse("@nocollapse \n * @nocollapse \n * @type {Object}*/", "extra @nocollapse tag")
}

func TestParsePureOrBreakMyCode(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsDocInfo := jt.parse("@pureOrBreakMyCode */")
	assertThat(t, jsDocInfo.IsPureOrBreakMyCode()).IsTrue()
}

func TestParseCollapsibleOrBreakMyCode(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsDocInfo := jt.parse("@collapsibleOrBreakMyCode */")
	assertThat(t, jsDocInfo.IsCollapsibleOrBreakMyCode()).IsTrue()
}

func TestParseDuplicatePureOrBreakMyCode(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@pureOrBreakMyCode \n * @pureOrBreakMyCode \n * @type {Object}*/",
		"extra @pureOrBreakMyCode tag")
}

func TestPreserveWhitespace1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.preserveWhitespaceParse(
			"this is a nice comment\n * that spans multiple lines\n *     with custom\n *     formatting\n * @author abc@google.com */\n")

	assertThat(t, jsdoc.GetBlockDescription()).IsEqualTo(
		"this is a nice comment\n that spans multiple lines\n     with custom\n     formatting\n" +
			" ")

	jt.assertDocumentationInMarker(
		jt.assertAnnotationMarker(jsdoc, "author", 4, 3), "abc@google.com", 10, 4, 24)
}

func TestPreserveWhitespace2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.preserveWhitespaceParse(
			"@param {string} x this is a nice comment\n * that spans multiple lines\n *     with custom\n *     formatting\n * @param {string} y */\n")

	assertThat(t, jsdoc.GetDescriptionForParameter("x")).IsEqualTo(
		" this is a nice comment\n that spans multiple lines\n     with custom\n     formatting\n" +
			" ")
}

func TestPreserveWhitespace3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.preserveWhitespaceParse(
			"@return {string} this is a nice comment\n * that spans multiple lines\n *     with custom\n *     formatting\n * @param {string} y */\n")

	assertThat(t, jsdoc.GetReturnDescription()).IsEqualTo(
		"this is a nice comment\n that spans multiple lines\n     with custom\n     formatting\n" +
			" ")
}

func TestPreserveWhitespace4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.preserveWhitespaceParse(
			"@throws {string} this is a nice comment\n * that spans multiple lines\n *     with custom\n *     formatting\n * @param {string} y */\n")

	assertThat(t, jsdoc.GetMarkers()[0].GetDescription().GetItem()).IsEqualTo(
		"{string} this is a nice comment\n that spans multiple lines\n     with custom\n     formatting\n" +
			" ")
}

func TestPreserveWhitespace5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.preserveWhitespaceParse(
			"@desc this is a nice comment\n * that spans multiple lines\n *     with custom\n *     formatting\n */\n")

	assertThat(t, jsdoc.GetDescription()).IsEqualTo(
		" this is a nice comment\n that spans multiple lines\n     with custom\n     formatting\n")
}

func TestPreserveWhitespace6(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info :=
		jt.preserveWhitespaceParse(
			"This is the typedef description\n *     with multiple lines\n * @typedef \n {string}*/\n")
	assertThat(t, info.GetBlockDescription()).IsEqualTo(
		"This is the typedef description\n     with multiple lines\n" +
			" ")
	assertThat(t, info.HasTypedefType()).IsTrue()
	jt.AssertTypeEqualsExpr(jt.STRING_TYPE, info.GetTypedefType())
}

func TestParseCommentWithStarsAfterLeadingSpace(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.preserveWhitespaceParse(
			"@desc this comment has extra\n * * stars on new lines\n */\n")

	assertThat(t, jsdoc.GetDescription()).IsEqualTo(
		" this comment has extra\n * stars on new lines\n")
}

func TestParseCommentWithThickLeadingStarBlockPreserveWhitespace(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.preserveWhitespaceParse(
			"@desc line 2 has extra stars\n **** that should pad content\n */\n")

	assertThat(t, jsdoc.GetDescription()).IsEqualTo(
		" line 2 has extra stars\n*** that should pad content\n")
}

func TestParseCommentWithThickLeadingStarBlockSingeLine(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jsdoc :=
		jt.parse(
			"@desc line 2 has extra stars\n **** that should pad content\n */\n")

	assertThat(t, jsdoc.GetDescription()).IsEqualTo("line 2 has extra stars *** that should pad content")
}

func TestParseCommentWithStarsOnOpenCommentLine(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	script :=
		jt.parseFull(
			"/******\n * This is a typedef comment with ASCII art.\n * @desc this is a description\n *****/\nfunction x() {}\n")
	jutil.CheckState(script.IsScript())
	fn := script.GetFirstChild()
	jutil.CheckState(fn.IsFunction())

	info := fn.GetJSDocInfo()
	assertThat(t, info.GetBlockDescription()).IsEqualTo("****\nThis is a typedef comment with ASCII art.")
	assertThat(t, info.GetDescription()).IsEqualTo("this is a description ***")
}

func TestParseTokenPositions1(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	// BANG
	//   STRING
	info := jt.parse("@type\n {!Object} */")
	assertThat(t, info).IsNotNull()
	assertThat(t, info.GetType()).IsNotNull()
	root := info.GetType().GetRoot()
	jt.checkTokenPosition(root, rhino.Token_BANG, 1, 2)
	jt.checkTokenPosition(root.GetFirstChild(), rhino.Token_STRINGLIT, 1, 3)
}

func TestParseTokenPositions2(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	// QMARK
	//   STRING
	info := jt.parse("@type {?Object} */")
	root := info.GetType().GetRoot()
	jt.checkTokenPosition(root, rhino.Token_QMARK, 0, 7)
	jt.checkTokenPosition(root.GetFirstChild(), rhino.Token_STRINGLIT, 0, 8)
}

func TestParseTokenPositions3(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	// BANG
	//   STRING
	info := jt.parse("@type {Object!} */")
	root := info.GetType().GetRoot()
	jt.checkTokenPosition(root, rhino.Token_BANG, 0, 13)
	jt.checkTokenPosition(root.GetFirstChild(), rhino.Token_STRINGLIT, 0, 7)
}

func TestParseTokenPositions4(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	// QMARK
	//   STRING
	info := jt.parse("@type {Object?} */")
	root := info.GetType().GetRoot()
	jt.checkTokenPosition(root, rhino.Token_QMARK, 0, 13)
	jt.checkTokenPosition(root.GetFirstChild(), rhino.Token_STRINGLIT, 0, 7)
}

func TestParseTokenPositions5(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	// PIPE
	//   STRING string
	//   BANG
	//     STRING Object
	info := jt.parse("@type {(string|!Object)} */")
	root := info.GetType().GetRoot()
	jt.checkTokenPosition(root.GetLastChild(), rhino.Token_BANG, 0, 15)
	jt.checkTokenPosition(root.GetLastChild().GetFirstChild(), rhino.Token_STRINGLIT, 0, 16)
}

func TestParseTokenPositions6(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	// STRING Set
	//   BANG
	//     STRING string
	info := jt.parse("@type {Set<!string>} */")
	root := info.GetType().GetRoot()
	jt.checkTokenPosition(root.GetFirstFirstChild(), rhino.Token_BANG, 0, 11)
	jt.checkTokenPosition(root.GetFirstChild().GetFirstFirstChild(), rhino.Token_STRINGLIT, 0, 12)
}

func TestParseTokenPositions7(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	// QMARK
	//   LC
	//     LB
	//       COLON
	//         STRING_KEY foo
	//         BANG
	//           STRING Object
	info := jt.parse("@type {?{foo: !Object}} */")
	root := info.GetType().GetRoot()
	jt.checkTokenPosition(root, rhino.Token_QMARK, 0, 7)
	jt.checkTokenPosition(root.GetFirstFirstChild().GetFirstFirstChild(), rhino.Token_STRING_KEY, 0, 9)
	jt.checkTokenPosition(root.GetFirstChild().GetFirstFirstChild().GetLastChild(), rhino.Token_BANG, 0, 14)
	jt.checkTokenPosition(
		root.GetFirstChild().GetFirstFirstChild().GetLastChild().GetFirstChild(),
		rhino.Token_STRINGLIT,
		0,
		15)
}

func TestParseTokenPositions8(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	// FUNCTION
	//   PARAM_LIST
	//     BANG
	//       STRING Object
	//   QMARK
	//     STRING string
	info := jt.parse("@type {function(!Object): ?string} */")
	root := info.GetType().GetRoot()
	jt.checkTokenPosition(root.GetFirstFirstChild(), rhino.Token_BANG, 0, 16)
	jt.checkTokenPosition(root.GetFirstChild().GetFirstFirstChild(), rhino.Token_STRINGLIT, 0, 17)
	jt.checkTokenPosition(root.GetLastChild(), rhino.Token_QMARK, 0, 26)
	jt.checkTokenPosition(root.GetLastChild().GetFirstChild(), rhino.Token_STRINGLIT, 0, 27)
}

func TestBoundedGeneric(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template {number} T */", "Bounded generic semantics are currently still in development")
}

func TestUnbalancedBracesBoundedGeneric(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template {number T */",
		"Bounded generic semantics are currently still in development",
		"Bad type annotation. expected closing } See"+
			" https://github.com/google/closure-compiler/wiki/Annotating-JavaScript-for-the-Closure-Compiler"+
			" for more information.")
}

func TestMultipleBoundedGeneric(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template {string} T,U */",
		"Bounded generic semantics are currently still in development",
		"Bad type annotation. Multiple template names cannot be declared with bounds or TTL."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestMultipleBoundsBoundedGeneric(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template {string} T, {number} U */",
		"Bounded generic semantics are currently still in development",
		"name not recognized due to syntax error.")
}

func TestBadTypeExpressionBoundedGeneric(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template {string || number} T */",
		"Bounded generic semantics are currently still in development",
		"Bad type annotation. type not recognized due to syntax error."+JsDocInfoParser_BAD_TYPE_WIKI_LINK,
		"name not recognized due to syntax error.",
		"Bad type annotation. @template tag missing type name."+JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestReuseTemplateTypeName(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		"@template {string} T \n* @template {number} T */",
		"Bounded generic semantics are currently still in development",
		"Bounded generic semantics are currently still in development",
		"Bad type annotation. Type name(s) for @template annotation declared twice."+
			JsDocInfoParser_BAD_TYPE_WIKI_LINK)
}

func TestMultipleTemplateBoundDeclarations(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	jt.parse(
		" @template {string} T \n* @template {number} U */",
		"Bounded generic semantics are currently still in development",
		"Bounded generic semantics are currently still in development")
}

func TestTemplateJSTypeExpression(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info :=
		jt.parse(
			"@template {string|number} T */",
			"Bounded generic semantics are currently still in development")

	assertThat(t, (info.GetTemplateTypes().Size() == 1)).IsTrue()
	assertThat(t, info.GetTemplateTypeNames()[0] == "T").IsTrue()

	root := info.GetTemplateTypes().Get("T").GetRoot()
	jt.checkTokenPosition(root, rhino.Token_PIPE, 0, 18)

	jt.checkTokenPosition(root.GetFirstChild(), rhino.Token_STRINGLIT, 0, 11)
	assertThat(t, root.GetFirstChild().GetString() == "string").IsTrue()

	jt.checkTokenPosition(root.GetSecondChild(), rhino.Token_STRINGLIT, 0, 18)
	assertThat(t, root.GetSecondChild().GetString() == "number").IsTrue()
}

func TestClosureUnawareMode_fileoverview_whitespace(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info :=
		jt.parseFileOverview(
			"@fileoverview\n@closureUnaware {WHITESPACE} */",
			"@closureUnaware annotation is not allowed in this compilation")
	assertThat(t, info.IsClosureUnawareCode()).IsTrue()
	assertThat(t, info.GetPerFileClosureUnawareMode()).IsEqualTo(rhino.JSDocInfo_PerFileClosureUnawareMode_WHITESPACE)
}

func TestClosureUnawareMode_fileoverview_simple(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info :=
		jt.parseFileOverview(
			"@fileoverview\n@closureUnaware {SIMPLE} */",
			"@closureUnaware annotation is not allowed in this compilation")
	assertThat(t, info.IsClosureUnawareCode()).IsTrue()
	assertThat(t, info.GetPerFileClosureUnawareMode()).IsEqualTo(rhino.JSDocInfo_PerFileClosureUnawareMode_SIMPLE)
}

func TestSupportsClosureUnaware_invalidValue_warnsAndDefaultsToSimple(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info :=
		jt.parseFileOverview(
			"@fileoverview @closureUnaware {FOOBAR} */",
			"@closureUnaware annotation is not allowed in this compilation",
			"invalid value for @closureUnaware: FOOBAR")
	assertThat(t, info.IsClosureUnawareCode()).IsTrue()
	assertThat(t, info.GetPerFileClosureUnawareMode()).IsEqualTo(rhino.JSDocInfo_PerFileClosureUnawareMode_SIMPLE)
}

func TestClosureUnawareMode_outsideFileoverview_errorIfSpecified(t *testing.T) {
	jt := newJsDocInfoParserTest(t)
	info :=
		jt.parse(
			"@closureUnaware {SIMPLE} */",
			"@closureUnaware annotation is not allowed in this compilation")
	assertThat(t, info.IsClosureUnawareCode()).IsTrue()
}
