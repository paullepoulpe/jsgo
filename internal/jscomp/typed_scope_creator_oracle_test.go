package jscomp

// Go-only: test-helper — the differential oracle for TypedScopeCreator (and
// FunctionTypeBuilder): the Go mirror of tools/oracle/java/TypedScopeDump.java.
//
// For each input (a file, or a directory compiled as one unit), it runs the
// prerequisites of TypedScopeCreatorTest.getProcessor() (GatherModuleMetadata,
// ModuleMapCreator, InferConsts) and the scope-building half of
// TypeInferencePass.inferAllScopes (createScope(ROOT, null), a scoped
// traversal creating every scope inside JSTypeResolver.openForDefinition(),
// resolveWeakImportsPreResolution(), finishAndFreeze()), then prints the
// diagnostics, every scope and TypedVar, JSTypeRegistry's scoped name table and
// the types attached to the AST exactly as the Java dump does, and compares.
//
// By default it checks the checked-in golden set
// (testdata/typed_scope/golden.txt.gz over testdata/typed_scope/inputs).
// Corpus run:
//
//	java -Xss1g -cp closure-compiler.jar:oracle-build com.google.javascript.jscomp.TypedScopeDump \
//	    out.txt.gz internal/jscomp/testdata/typed_scope/default_externs.js .cache/corpus @.cache/corpus-files.txt
//	TYPED_SCOPE_ORACLE=out.txt.gz TYPED_SCOPE_ROOT=.cache/corpus go test -run TestTypedScopeCreatorOracle ./internal/jscomp
//
// TYPED_SCOPE_DIFF_DIR=<dir> writes the Go dump of every mismatching input
// there (as <n>.go.txt / <n>.java.txt). Known, accepted difference: the JVM's
// "helpful" NullPointerException messages cannot be reproduced, so X lines of
// a java.lang.NullPointerException are compared by class only.

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const typedScopeDump_BASELINE = "<baseline>"

// typedScopeDump_options are the options of CompilerTestCase.getOptions() as
// TypedScopeCreatorTest sets them up.
func typedScopeDump_options() *CompilerOptions {
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_UNSUPPORTED)
	options.SetEmitUseStrict(false)
	options.SetLanguageOut(CompilerOptions_LanguageMode_NO_TRANSPILE)
	options.SetModuleResolutionMode(ModuleLoader_ResolutionMode_BROWSER)
	options.SetParseJsDocDocumentation(parsing.Config_JsDocParsing_TYPES_ONLY)
	options.SetAssumeStaticInheritanceIsNotUsed(true)
	options.SetPreserveTypeAnnotations(true)
	options.SetAssumeGettersArePure(false)
	options.SetCheckSymbols(true)
	options.SetWarningLevel(DiagnosticGroups_INVALID_CASTS, CheckLevel_WARNING)
	options.SetWarningLevel(DiagnosticGroups_MISPLACED_MSG_ANNOTATION, CheckLevel_WARNING)
	options.SetWarningLevel(DiagnosticGroups_MISSING_PROPERTIES, CheckLevel_WARNING)
	options.SetCodingConvention(NewGoogleCodingConvention())
	polymerVersion := int32(1)
	options.SetPolymerVersion(&polymerVersion)
	options.SetRuntimeLibraryMode(RuntimeJsLibManager_RuntimeLibraryMode_RECORD_AND_VALIDATE_FIELDS)
	// TypedScopeCreatorTest.getOptions()
	options.SetClosurePass(false)
	return options
}

// typedScopeDump_errorManager records every report in order.
type typedScopeDump_errorManager struct {
	lines        []string
	errors       []*JSError
	warnings     []*JSError
	typedPercent float64
}

func (m *typedScopeDump_errorManager) Report(level *CheckLevel, error_ *JSError) {
	m.lines = append(m.lines, fmt.Sprintf("D %s %s %d:%d %s",
		jutil.ToJavaString(level), error_.Type().key, error_.Lineno(), error_.Charno(), error_.Description()))
	if level != nil && *level == CheckLevel_ERROR {
		m.errors = append(m.errors, error_)
	} else if level != nil && *level == CheckLevel_WARNING {
		m.warnings = append(m.warnings, error_)
	}
}
func (m *typedScopeDump_errorManager) GenerateReport()         {}
func (m *typedScopeDump_errorManager) GetErrorCount() int      { return len(m.errors) }
func (m *typedScopeDump_errorManager) GetWarningCount() int    { return len(m.warnings) }
func (m *typedScopeDump_errorManager) GetErrors() []*JSError   { return slices.Clone(m.errors) }
func (m *typedScopeDump_errorManager) GetWarnings() []*JSError { return slices.Clone(m.warnings) }
func (m *typedScopeDump_errorManager) SetTypedPercent(typedPercent float64) {
	m.typedPercent = typedPercent
}
func (m *typedScopeDump_errorManager) GetTypedPercent() float64 { return m.typedPercent }
func (m *typedScopeDump_errorManager) HasHaltingErrors() bool {
	return ErrorManager_HasHaltingErrors(m)
}
func (m *typedScopeDump_errorManager) ShouldReportConformanceViolation(
	requirement *Requirement,
	allowlistEntry *RequirementScopeEntry,
	diagnostic *JSError,
	behavior ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior,
	isAllowlisted bool) bool {
	return ErrorManager_ShouldReportConformanceViolation(m, requirement, allowlistEntry, diagnostic, behavior, isAllowlisted)
}

// typedScopeDump_firstScopeBuildingCallback is
// TypeInferencePass.FirstScopeBuildingCallback.
type typedScopeDump_firstScopeBuildingCallback struct {
	NodeTraversal_AbstractScopedCallbackBase
}

func (c *typedScopeDump_firstScopeBuildingCallback) EnterScope(t *NodeTraversal) { t.GetTypedScope() }
func (c *typedScopeDump_firstScopeBuildingCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
}

// typedScopeDump_exceptionLine renders an escaped exception like the Java
// "X <class>: <first message line>" of the innermost cause.
func typedScopeDump_exceptionLine(r any) string {
	th := jutil.AsThrowable(r)
	for th.GetCause() != nil && th.GetCause() != th {
		th = th.GetCause()
	}
	s := th.Error()
	msg := th.GetMessage()
	className := s
	if strings.HasSuffix(s, ": "+msg) {
		className = strings.TrimSuffix(s, ": "+msg)
	} else {
		msg = "null"
	}
	if nl := strings.IndexByte(msg, '\n'); nl >= 0 {
		msg = msg[:nl]
	}
	return "X " + className + ": " + msg
}

func typedScopeDump(externs string, inputs []*SourceFile, baseline *jutil.LinkedHashSet[string]) (out []string) {
	compiler := NewCompiler()
	compiler.SetAllowableFeatures(CompilerOptions_LanguageMode_UNSUPPORTED.ToFeatureSet())
	compiler.SetPreferRegexParser(false)
	em := &typedScopeDump_errorManager{}
	compiler.Init(
		[]*SourceFile{SourceFile_FromCode("externs", externs)},
		inputs,
		typedScopeDump_options())
	compiler.SetErrorManager(em)
	root := compiler.ParseInputs()
	if root.IsNil() || compiler.GetErrorCount() > 0 {
		out = append(out, fmt.Sprintf("P %d", compiler.GetErrorCount()))
		out = append(out, em.lines...)
		return out
	}
	em.lines = nil
	externsRoot := root.GetFirstChild()
	jsRoot := root.GetLastChild()
	var scopeCreator *TypedScopeCreator
	failed := false
	jutil.Try(func() {
		NewGatherModuleMetadata(compiler, false, ModuleLoader_ResolutionMode_BROWSER).
			Process(externsRoot, jsRoot)
		NewModuleMapCreator(compiler, compiler.GetModuleMetadataMap()).Process(externsRoot, jsRoot)
		NewInferConsts(compiler).Process(externsRoot, jsRoot)

		registry := compiler.GetTypeRegistry()
		scopeCreator = NewTypedScopeCreator(compiler)
		func() {
			closer := registry.GetResolver().OpenForDefinition()
			defer closer.Close()
			topScope := scopeCreator.CreateScope(root, nil).(*TypedScope)
			NodeTraversal_Builder_().
				SetCompiler(compiler).
				SetCallback(&typedScopeDump_firstScopeBuildingCallback{}).
				SetScopeCreator(scopeCreator).
				TraverseWithScope(root, topScope)
			scopeCreator.ResolveWeakImportsPreResolution()
		}()
		scopeCreator.FinishAndFreeze()
	}, func(e jutil.Throwable) {
		out = append(out, em.lines...)
		out = append(out, typedScopeDump_exceptionLine(e))
		failed = true
	})
	if failed {
		return out
	}
	out = append(out, em.lines...)

	scopes := scopeCreator.GetAllMemoizedScopes()
	slices.Reverse(scopes)
	scopeIndex := map[rhino.Node]int{}
	for i, s := range scopes {
		scopeIndex[s.GetRootNode()] = i
	}
	for i, s := range scopes {
		r := s.GetRootNode()
		if i > 0 && r.IsFromExterns() {
			continue
		}
		parent := -1
		if s.GetParent() != nil {
			parent = scopeIndex[s.GetParent().GetRootNode()]
		}
		out = append(out, fmt.Sprintf("S %d %d %s %s:%d:%d",
			i, parent, r.GetToken(), jutil.ToJavaString(r.GetSourceFileName()), r.GetLineno(), r.GetCharno()))
		for v := range s.GetVarIterable() {
			line := typedScopeDump_varLine(v)
			if i == 0 && baseline != nil && baseline.Contains(line) {
				continue
			}
			out = append(out, line)
		}
	}

	var table []string
	compiler.GetTypeRegistry().ForEachScopedType(func(rowKey rhino.Node, isGlobal bool, name string, type_ rhino.JSType) {
		var rowLabel string
		if isGlobal {
			rowLabel = "GLOBAL"
		} else {
			rowLabel = jutil.ToJavaString(rowKey.GetSourceFileName()) + ":" + typedScopeDump_pos(rowKey)
			if _, ok := scopeIndex[rowKey]; !ok {
				rowLabel += "!"
			}
		}
		line := "T " + rowLabel + " " + name + " | " + typedScopeDump_typeString(type_)
		if baseline != nil && baseline.Contains(line) {
			return
		}
		table = append(table, line)
	})
	slices.SortFunc(table, jutil.StringCompareTo)
	out = append(out, table...)

	out = typedScopeDump_nodeTypes(out, jsRoot)
	return out
}

func typedScopeDump_varLine(v *TypedVar) string {
	type_ := v.GetType()
	ann := "-"
	if type_ != nil && !v.IsTypeInferred() {
		jutil.Try(func() {
			ann = type_.ToAnnotationString(rhino.JSType_Nullability_EXPLICIT)
		}, func(e jutil.IRuntimeException) {
			ann = "!" + typedScopeDump_simpleClassName(e)
		})
	}
	node := "null"
	if !v.GetNameNode().IsNil() {
		node = typedScopeDump_pos(v.GetNameNode())
	}
	return fmt.Sprintf("V %s | %s | %s | inferred=%t esc=%t once=%t input=%s node=%s",
		v.GetName(), typedScopeDump_typeString(type_), ann, v.IsTypeInferred(),
		v.IsMarkedEscaped(), v.IsMarkedAssignedExactlyOnce(), v.GetInputName(), node)
}

func typedScopeDump_simpleClassName(e jutil.Throwable) string {
	s := e.Error()
	if i := strings.Index(s, ": "); i >= 0 {
		s = s[:i]
	}
	return s[strings.LastIndexByte(s, '.')+1:]
}

func typedScopeDump_typeString(type_ rhino.JSType) (s string) {
	if type_ == nil {
		return "null"
	}
	jutil.Try(func() {
		s = type_.String()
	}, func(e jutil.IRuntimeException) {
		s = "!" + typedScopeDump_simpleClassName(e)
	})
	return s
}

func typedScopeDump_pos(n rhino.Node) string {
	return fmt.Sprintf("%s@%d:%d", n.GetToken(), n.GetLineno(), n.GetCharno())
}

func typedScopeDump_nodeTypes(out []string, n rhino.Node) []string {
	type_ := n.GetJSType()
	if type_ != nil && !typedScopeDump_isLiteral(n.GetToken()) {
		out = append(out, fmt.Sprintf("N %s %s:%d:%d | %s", n.GetToken(),
			jutil.ToJavaString(n.GetSourceFileName()), n.GetLineno(), n.GetCharno(), typedScopeDump_typeString(type_)))
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		out = typedScopeDump_nodeTypes(out, c)
	}
	return out
}

func typedScopeDump_isLiteral(t rhino.Token) bool {
	switch t {
	case rhino.Token_NULL, rhino.Token_VOID, rhino.Token_STRINGLIT, rhino.Token_TEMPLATELIT_STRING,
		rhino.Token_NUMBER, rhino.Token_BIGINT, rhino.Token_TRUE, rhino.Token_FALSE, rhino.Token_REGEXP,
		rhino.Token_ARRAYLIT:
		return true
	}
	return false
}

// typedScopeDump_escape is TypedScopeDump.escape: printable ASCII, "\\" for a
// backslash, \uXXXX for every other UTF-16 unit.
func typedScopeDump_escape(s string) string {
	var b strings.Builder
	for _, c := range jutil.StringToUTF16(s) {
		switch {
		case c == '\\':
			b.WriteString("\\\\")
		case c >= 0x20 && c < 0x7f:
			b.WriteByte(byte(c))
		default:
			fmt.Fprintf(&b, "\\u%04x", c)
		}
	}
	return b.String()
}

// typedScopeDump_inputs reads one oracle input: a file, or a directory whose
// files (in name order) form one compilation.
func typedScopeDump_inputs(t *testing.T, root, rel string) []*SourceFile {
	p := filepath.Join(root, rel)
	st, err := os.Stat(p)
	if err != nil {
		t.Fatal(err)
	}
	if !st.IsDir() {
		src, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		return []*SourceFile{SourceFile_FromCode(rel, string(src))}
	}
	entries, err := os.ReadDir(p)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	slices.SortFunc(names, jutil.StringCompareTo)
	var inputs []*SourceFile
	for _, name := range names {
		src, err := os.ReadFile(filepath.Join(p, name))
		if err != nil {
			t.Fatal(err)
		}
		inputs = append(inputs, SourceFile_FromCode(rel+"/"+name, string(src)))
	}
	return inputs
}

// typedScopeDump_sameLine compares a Go and a Java line, masking the
// message of a java.lang.NullPointerException (see the file comment).
func typedScopeDump_sameLine(goLine, javaLine string) bool {
	const npe = "X java.lang.NullPointerException"
	if strings.HasPrefix(goLine, npe) && strings.HasPrefix(javaLine, npe) {
		return true
	}
	return goLine == javaLine
}

func TestTypedScopeCreatorOracle(t *testing.T) {
	oraclePath := os.Getenv("TYPED_SCOPE_ORACLE")
	root := os.Getenv("TYPED_SCOPE_ROOT")
	if oraclePath == "" {
		oraclePath = "testdata/typed_scope/golden.txt.gz"
		root = "testdata/typed_scope/inputs"
	}
	diffDir := os.Getenv("TYPED_SCOPE_DIFF_DIR")
	externsBytes, err := os.ReadFile("testdata/typed_scope/default_externs.js")
	if err != nil {
		t.Fatal(err)
	}
	externs := string(externsBytes)

	f, err := os.Open(oraclePath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	sc := bufio.NewScanner(gz)
	sc.Buffer(make([]byte, 1<<20), 1<<30)

	var baseline *jutil.LinkedHashSet[string]
	files, mismatches := 0, 0
	check := func(name string, javaLines []string) {
		var goLines []string
		if name == typedScopeDump_BASELINE {
			raw := typedScopeDump(externs, []*SourceFile{SourceFile_FromCode(typedScopeDump_BASELINE, "")}, nil)
			baseline = jutil.NewLinkedHashSet(raw...)
			goLines = raw
		} else {
			goLines = typedScopeDump(externs, typedScopeDump_inputs(t, root, name), baseline)
		}
		for i := range goLines {
			goLines[i] = typedScopeDump_escape(goLines[i])
		}
		files++
		same := len(goLines) == len(javaLines)
		first := -1
		for i := 0; i < min(len(goLines), len(javaLines)); i++ {
			if !typedScopeDump_sameLine(goLines[i], javaLines[i]) {
				same = false
				first = i
				break
			}
		}
		if same {
			return
		}
		mismatches++
		if first < 0 {
			first = min(len(goLines), len(javaLines))
		}
		get := func(lines []string, i int) string {
			if i < len(lines) {
				return lines[i]
			}
			return "<end>"
		}
		if mismatches <= 20 {
			t.Errorf("%s: first difference at line %d\n  go:   %s\n  java: %s",
				name, first+1, get(goLines, first), get(javaLines, first))
		}
		if diffDir != "" {
			base := filepath.Join(diffDir, fmt.Sprintf("%05d", mismatches))
			_ = os.WriteFile(base+".go.txt", []byte(name+"\n"+strings.Join(goLines, "\n")+"\n"), 0o644)
			_ = os.WriteFile(base+".java.txt", []byte(name+"\n"+strings.Join(javaLines, "\n")+"\n"), 0o644)
		}
	}

	var current string
	var lines []string
	started := false
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "=== ") {
			if started {
				check(current, lines)
			}
			current, lines, started = strings.TrimPrefix(line, "=== "), nil, true
			continue
		}
		lines = append(lines, line)
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	if started {
		check(current, lines)
	}
	t.Logf("%d inputs, %d mismatches", files, mismatches)
	if files < 2 {
		t.Fatalf("only %d inputs in %s", files, oraclePath)
	}
}
