package jscomp

// Differential test against tools/oracle/java/InlineDump.java: the inlining
// passes (InlineVariables in LOCALS_ONLY and ALL modes,
// FlowSensitiveInlineVariables, InlineFunctions, InlineAndCollapseProperties
// with PropertyCollapseLevel ALL and NONE, InlineSimpleMethods) run on
// normalized files, plus a GlobalNamespace dump, printed in the same format
// and compared file by file.
//
//   - TestInlineGolden replays the checked-in golden dump
//     testdata/inline/golden.txt.gz (with --full trees): the hand-written
//     inputs testdata/inline/js/*.js, followed by the corpus files of
//     testdata/inline/golden-files.txt (skipped when the corpus is not
//     present). Regenerate with
//       CP=.cache/$CLOSURE_JAR:<InlineDump build dir>
//       J="java -Xss1g -XX:-ShowCodeDetailsInExceptionMessages -cp $CP"
//       $J com.google.javascript.jscomp.InlineDump --full a.gz \
//         internal/jscomp/testdata/inline/js <the *.js names>
//       $J com.google.javascript.jscomp.InlineDump --full b.gz \
//         .cache/corpus --list=internal/jscomp/testdata/inline/golden-files.txt
//       cat a.gz b.gz > internal/jscomp/testdata/inline/golden.txt.gz
//   - TestInlineOracle compares any other dump (hashes, or --full):
//     INLINE_ORACLE=<dump.gz> [INLINE_ROOT=<root>] (default root: the corpus
//     under /home/user/jsgo/.cache/corpus). The held-out corpus
//     (.cache/corpus-holdout, corpus-holdout-files.txt) is for a final check
//     only.

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// inlineOracleExterns is InlineDump.EXTERNS.
const inlineOracleExterns = "var window; var Arguments; var arguments; /** @constructor */ function String() {}\n" +
	"/** @nosideeffects */ function nochg(){}\nfunction chg(){}\n"

// inlineOraclePasses is InlineDump.PASSES.
var inlineOraclePasses = []string{
	"inline-variables-locals",
	"inline-variables-all",
	"flow-sensitive",
	"inline-functions",
	"collapse-all",
	"collapse-none",
	"inline-simple-methods",
}

type inlineOracleDumper struct {
	full bool
	sb   strings.Builder
}

func (d *inlineOracleDumper) line(s string) {
	d.sb.WriteString(normalizeOracleEscape(s))
	d.sb.WriteByte('\n')
}

// newCompiler mirrors InlineDump.newCompiler.
func (d *inlineOracleDumper) newCompiler(rel string, src string) *Compiler {
	compiler := NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
	options := NewCompilerOptions()
	options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
	options.SetLanguageOut(CompilerOptions_LanguageMode_NO_TRANSPILE)
	options.SetWarningLevel(DiagnosticGroups_MODULE_LOAD, CheckLevel_OFF)
	options.SetCodingConvention(NewGoogleCodingConvention())
	compiler.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", inlineOracleExterns)},
		[]*SourceFile{SourceFile_FromCode(rel, src)},
		options)
	compiler.Parse()
	return compiler
}

// newPass mirrors InlineDump.newPass.
func inlineOracleNewPass(name string, c *Compiler) CompilerPass {
	switch name {
	case "inline-variables-locals":
		return NewInlineVariables(c, InlineVariables_Mode_LOCALS_ONLY)
	case "inline-variables-all":
		return NewInlineVariables(c, InlineVariables_Mode_ALL)
	case "flow-sensitive":
		return NewFlowSensitiveInlineVariables(c)
	case "inline-functions":
		return NewInlineFunctions(
			c,
			c.GetUniqueNameIdSupplier(),
			CompilerOptions_Reach_ALL,
			/* assumeStrictThis= */ false,
			CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING)
	case "collapse-all":
		return InlineAndCollapseProperties_Builder_(c).
			SetPropertyCollapseLevel(CompilerOptions_PropertyCollapseLevel_ALL).
			SetChunkOutputType(CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE).
			SetHaveModulesBeenRewritten(false).
			SetModuleResolutionMode(ModuleLoader_ResolutionMode_BROWSER).
			Build()
	case "collapse-none":
		return InlineAndCollapseProperties_Builder_(c).
			SetPropertyCollapseLevel(CompilerOptions_PropertyCollapseLevel_NONE).
			SetChunkOutputType(CompilerOptions_ChunkOutputType_GLOBAL_NAMESPACE).
			SetHaveModulesBeenRewritten(false).
			SetModuleResolutionMode(ModuleLoader_ResolutionMode_BROWSER).
			Build()
	case "inline-simple-methods":
		c.SetExternProperties(jutil.NewLinkedHashSet[string]())
		return NewInlineSimpleMethods(c)
	}
	panic(jutil.NewIllegalArgumentExceptionMsg(name))
}

// normalized mirrors InlineDump.normalized: nil (after printing why) if
// Normalize threw.
func (d *inlineOracleDumper) normalized(sb *inlineOracleDumper, rel string, src string) (c *Compiler) {
	c = d.newCompiler(rel, src)
	defer func() {
		if r := recover(); r != nil {
			sb.line("# normalize exception " + normalizeOracleException(r))
			c = nil
		}
	}()
	Normalize_CreateNormalizeForOptimizations(c).Process(c.GetExternsRoot(), c.GetJsRoot())
	return c
}

// dumpFile mirrors InlineDump.dumpFile.
func (d *inlineOracleDumper) dumpFile(rel string, src string) {
	d.line("## file " + rel)
	compiler := d.newCompiler(rel, src)
	parseDiags := normalizeOracleDiagnostics(compiler)
	if compiler.HasErrors() {
		d.line("# parse errors")
	} else {
		d.line("# parse ok")
	}
	for _, diag := range parseDiags {
		d.line(diag)
	}
	if compiler.HasErrors() {
		return
	}
	for _, pass := range inlineOraclePasses {
		head := &inlineOracleDumper{}
		c := d.normalized(head, rel, src)
		if c == nil {
			d.sb.WriteString(head.sb.String())
			return
		}
		d.line("## pass " + pass)
		baseDiags := normalizeOracleDiagnostics(c)
		before := c.GetRoot().ToStringTree()
		exception := ""
		func() {
			defer func() {
				if r := recover(); r != nil {
					exception = "exception " + normalizeOracleException(r)
				}
			}()
			c.BeforePass(pass)
			inlineOracleNewPass(pass, c).Process(c.GetExternsRoot(), c.GetJsRoot())
		}()
		diags := normalizeOracleDiagnostics(c)
		for _, diag := range baseDiags {
			diags = normalizeOracleRemoveFirst(diags, diag)
		}
		for _, diag := range diags {
			d.line(diag)
		}
		if exception != "" {
			d.line(exception)
			continue
		}
		after := c.GetRoot().ToStringTree()
		if after == before {
			d.line("tree unchanged")
		} else {
			var lines []string
			for _, l := range strings.Split(after, "\n") {
				lines = append(lines, normalizeOracleEscape(l))
			}
			if lines[len(lines)-1] == "" {
				lines = lines[:len(lines)-1]
			}
			d.line("tree changed " + normalizeOracleSha16(lines))
			if d.full {
				for _, l := range lines {
					d.sb.WriteString(l + "\n")
				}
			}
		}
		var changes []string
		normalizeOracleCollectChanges(c.GetRoot(), &changes)
		var escapedChanges []string
		for _, l := range changes {
			escapedChanges = append(escapedChanges, normalizeOracleEscape(l))
		}
		d.line("changes " + strconv.Itoa(len(changes)) + " " + normalizeOracleSha16(escapedChanges))
		if d.full {
			for _, l := range escapedChanges {
				d.sb.WriteString(l + "\n")
			}
		}
	}
	c := d.normalized(d, rel, src)
	if c == nil {
		return
	}
	d.line("## namespace")
	var lines []string
	count := 0
	exception := ""
	func() {
		defer func() {
			if r := recover(); r != nil {
				exception = "exception " + normalizeOracleException(r)
			}
		}()
		namespace := NewGlobalNamespaceWithRoots(c, c.GetExternsRoot(), c.GetJsRoot())
		for name := range namespace.GetNameIndex().Values() {
			count++
			inlineOracleDumpName(&lines, name)
		}
	}()
	if exception != "" {
		d.line(exception)
		return
	}
	d.line("namespace " + strconv.Itoa(count) + " " + normalizeOracleSha16(lines))
	if d.full {
		for _, l := range lines {
			d.sb.WriteString(l + "\n")
		}
	}
}

// inlineOracleAttempt is InlineDump.attempt: String.valueOf(f.call()), or
// "!" + the Java class name of what it threw.
func inlineOracleAttempt[T any](f func() T) (s string) {
	defer func() {
		if r := recover(); r != nil {
			t := jutil.AsThrowable(r)
			full := t.Error()
			s = "!" + strings.TrimSuffix(full, ": "+t.GetMessage())
		}
	}()
	return inlineOracleValueOf(f())
}

func inlineOracleValueOf(v any) string {
	switch x := v.(type) {
	case bool:
		return strconv.FormatBool(x)
	case GlobalNamespace_Inlinability:
		return x.String()
	}
	panic("unexpected value")
}

// inlineOracleDumpName mirrors InlineDump.dumpName.
func inlineOracleDumpName(lines *[]string, name *GlobalNamespace_Name) {
	*lines = append(*lines, normalizeOracleEscape("name "+name.String()))
	*lines = append(*lines, normalizeOracleEscape(
		"  decisions "+
			inlineOracleAttempt(name.CanCollapseOrInline)+
			" "+
			inlineOracleAttempt(name.CanCollapseOrInlineChildNames)+
			" "+
			inlineOracleAttempt(name.CalculateInlinability)+
			" needsStub="+
			inlineOracleAttempt(name.NeedsToBeStubbed)+
			" nsObjLit="+
			inlineOracleAttempt(name.IsNamespaceObjectLit)+
			" simpleStub="+
			inlineOracleAttempt(name.IsSimpleStubDeclaration)))
	for _, ref := range name.GetRefs() {
		n := ref.GetNode()
		parent := "null"
		if n != nil && n.GetParent() != nil {
			parent = n.GetParent().GetToken().String()
		}
		scope := "null"
		if ref.scope != nil {
			scope = inlineOracleNode(ref.scope.GetRootNode())
		}
		s := "  ref " + ref.type_.String() + " " + inlineOracleNode(n) + " parent=" + parent + " scope=" + scope
		if ref == name.getDeclarationRef() {
			s += " decl"
		}
		if ref == name.GetInitialization() {
			s += " init"
		}
		*lines = append(*lines, normalizeOracleEscape(s))
	}
}

// inlineOracleNode is InlineDump.node.
func inlineOracleNode(n *rhino.Node) string {
	if n == nil {
		return "null"
	}
	s := n.GetToken().String() + "@" + strconv.Itoa(n.GetLineno()) + ":" + strconv.Itoa(n.GetCharno())
	if n.IsQualifiedName() {
		s += " " + *n.GetQualifiedName()
	} else if n.IsStringKey() || n.IsMemberFunctionDef() || n.IsMemberFieldDef() || n.IsGetterDef() || n.IsSetterDef() {
		s += " " + n.GetString()
	}
	return s
}

// inlineOracleCompare dumps every file of the dump at path (sources under
// root) and compares; it returns the number of files compared.
func inlineOracleCompare(t *testing.T, path string, root string, onlyNames func(string) bool) int {
	t.Helper()
	names, sections, full := normalizeOracleSplit(t, path)
	for _, s := range sections {
		if strings.Contains(s, "\nname ") {
			full = true
			break
		}
	}
	n, failures := 0, 0
	for _, name := range names {
		rel := normalizeOracleUnescapeName(name)
		if onlyNames != nil && !onlyNames(rel) {
			continue
		}
		src, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			t.Logf("skipping %s: %v", rel, err)
			continue
		}
		d := &inlineOracleDumper{full: full}
		func() {
			defer func() {
				if r := recover(); r != nil {
					d.line("go panic " + normalizeOracleException(r))
				}
			}()
			d.dumpFile(rel, jutil.NewStringFromUTF8Bytes(src))
		}()
		n++
		if got, want := d.sb.String(), sections[name]; got != want {
			failures++
			if failures <= 10 {
				t.Errorf("%s: dump differs:\n%s", rel, normalizeOracleFirstDiff(got, want))
			}
		}
	}
	if failures > 0 {
		t.Errorf("%d of %d files differ", failures, n)
	}
	return n
}

func TestInlineGolden(t *testing.T) {
	golden := "testdata/inline/golden.txt.gz"
	jsDir := "testdata/inline/js"
	hand := map[string]bool{}
	entries, err := os.ReadDir(jsDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		hand[e.Name()] = true
	}
	if n := inlineOracleCompare(t, golden, jsDir, func(rel string) bool { return hand[rel] }); n != len(hand) {
		t.Errorf("compared %d hand-written files, want %d", n, len(hand))
	}
	corpus := oracleTestCorpusRoot(t)
	if _, err := os.Stat(corpus); err != nil {
		t.Logf("corpus not present; golden corpus files skipped")
		return
	}
	inlineOracleCompare(t, golden, corpus, func(rel string) bool { return !hand[rel] })
}

func TestInlineOracle(t *testing.T) {
	dump := os.Getenv("INLINE_ORACLE")
	if dump == "" {
		t.Skip("set INLINE_ORACLE=<InlineDump output .gz> (and INLINE_ROOT) to run")
	}
	root := os.Getenv("INLINE_ROOT")
	if root == "" {
		root = oracleTestCorpusRoot(t)
	}
	n := inlineOracleCompare(t, dump, root, nil)
	t.Logf("compared %d files", n)
}
