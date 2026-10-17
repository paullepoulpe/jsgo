package jscomp

// Differential test for the jscomp-a port: replays the Java reference
// vectors in testdata/jscomp_a.tsv (written by
// tools/oracle/java/JscompAVectors.java against closure-compiler v20261015)
// and requires byte-identical results. Each row records its inputs; the
// diagnostic types, groups and the source excerpt provider below mirror the
// oracle's TYPES, GROUPS and PROVIDER.

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

// jsaNull stands for a Java null column ("\N" in the file).
const jsaNull = "\x00null"

func jsaUnescape(s string) string {
	if s == `\N` {
		return jsaNull
	}
	if !strings.Contains(s, `\`) {
		return s
	}
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\\' && i+1 < len(s) {
			i++
			switch s[i] {
			case 'n':
				b.WriteByte('\n')
			case 't':
				b.WriteByte('\t')
			case 'r':
				b.WriteByte('\r')
			default:
				b.WriteByte(s[i])
			}
			continue
		}
		b.WriteByte(c)
	}
	return b.String()
}

func jsaReadVectors(t *testing.T) [][]string {
	f, err := os.Open("testdata/jscomp_a.tsv")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var rows [][]string
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<24)
	for sc.Scan() {
		cols := strings.Split(sc.Text(), "\t")
		for i := range cols {
			cols[i] = jsaUnescape(cols[i])
		}
		rows = append(rows, cols)
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	return rows
}

// jsaNullable maps a column to a Java String that may be null.
func jsaNullable(s string) *string {
	if s == jsaNull {
		return nil
	}
	return &s
}

// jsaJavaStr renders a nullable String as Java string concatenation does.
func jsaJavaStr(s *string) string {
	if s == nil {
		return "null"
	}
	return *s
}

func jsaLevelStr(l *CheckLevel) string {
	if l == nil {
		return "null"
	}
	return l.String()
}

// jsaExc renders a thrown exception as the oracle's exc(): "!" + simple
// class name + ": " + message.
func jsaExc(r any) string {
	s := jutil.AsThrowable(r).Error()
	class, msg, ok := strings.Cut(s, ": ")
	if !ok {
		msg = "null"
	}
	if i := strings.LastIndexAny(class, ".$"); i >= 0 {
		class = class[i+1:]
	}
	return "!" + class + ": " + msg
}

// jsaSafe runs f, rendering a thrown exception with jsaExc.
func jsaSafe(f func() string) (out string) {
	defer func() {
		if r := recover(); r != nil {
			out = jsaExc(r)
		}
	}()
	return f()
}

func jsaAtoi(t *testing.T, s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		t.Fatalf("bad int %q", s)
	}
	return n
}

// jsaSplit is String.split(sep) for the specs' lists, where "" is the empty
// list.
func jsaSplit(s, sep string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, sep)
}

func jsaMapSpec(m *jutil.LinkedHashMap[string, string]) string {
	var b strings.Builder
	for k, v := range m.All() {
		if b.Len() > 0 {
			b.WriteByte(',')
		}
		b.WriteString(k + "=" + v)
	}
	return b.String()
}

func jsaParseMapSpec(s string) *jutil.LinkedHashMap[string, string] {
	m := jutil.NewLinkedHashMap[string, string]()
	for _, e := range jsaSplit(s, ",") {
		k, v, _ := strings.Cut(e, "=")
		m.Put(k, v)
	}
	return m
}

var (
	jsaDiffTypes = []*DiagnosticType{
		DiagnosticType_Error("JSA_E0", "error zero {0}"),
		DiagnosticType_Warning("JSA_W1", "warning one <{0}>"),
		DiagnosticType_Disabled("JSA_D2", "disabled two"),
		ClosurePrimitiveErrors_DUPLICATE_NAMESPACE_AND_MODULE,
		RhinoErrorReporter_TYPE_PARSE_ERROR,
		StrictWarningsGuard_UNRAISABLE_WARNING,
		DiagnosticType_Warning("JSA_W6", "w6 '{0}' \"q\" & é"),
	}
	jsaDiffGroups = []*DiagnosticGroup{
		NewDiagnosticGroupWithName(jutil.StringPtr("g0"), jsaDiffTypes[0], jsaDiffTypes[1]),
		NewDiagnosticGroupWithName(jutil.StringPtr("g1"), jsaDiffTypes[1]),
		NewDiagnosticGroupWithName(jutil.StringPtr("g2"), jsaDiffTypes[2], jsaDiffTypes[6]),
		DiagnosticGroup_ForType(jsaDiffTypes[6]),
		NewDiagnosticGroupWithNameFromGroups(jutil.StringPtr("g4"),
			NewDiagnosticGroupWithName(jutil.StringPtr("g0"), jsaDiffTypes[0], jsaDiffTypes[1]),
			NewDiagnosticGroupWithName(jutil.StringPtr("g2"), jsaDiffTypes[2], jsaDiffTypes[6])),
		NewDiagnosticGroupWithName(jutil.StringPtr("g5"), jsaDiffTypes[4], jsaDiffTypes[3]),
	}
	jsaDiffLevels = []CheckLevel{CheckLevel_ERROR, CheckLevel_WARNING, CheckLevel_OFF}
)

func jsaCheckLevelValueOf(s string) CheckLevel {
	for _, l := range CheckLevel_values() {
		if l.String() == s {
			return l
		}
	}
	panic(jutil.NewIllegalArgumentExceptionMsg("No enum constant com.google.javascript.jscomp.CheckLevel." + s))
}

// jsaDiffProvider is the oracle's PROVIDER: excerpts from the "src" files;
// names starting with "orig/" are the reverse mapped originals (the same
// file); lines whose number is 3 mod 5 have such a mapping.
type jsaDiffProvider struct{ files map[string]*SourceFile }

func (p *jsaDiffProvider) file(name *string) *SourceFile {
	if name == nil {
		return nil
	}
	return p.files[strings.TrimPrefix(*name, "orig/")]
}

func (p *jsaDiffProvider) GetSourceLine(sourceName *string, lineNumber int) *string {
	f := p.file(sourceName)
	if f == nil {
		return nil
	}
	return f.GetLine(lineNumber)
}

func (p *jsaDiffProvider) GetSourceLines(sourceName *string, lineNumber int, length int) Region {
	f := p.file(sourceName)
	if f == nil {
		return nil
	}
	return f.GetLines(lineNumber, length)
}

func (p *jsaDiffProvider) GetSourceRegion(sourceName *string, lineNumber int) Region {
	f := p.file(sourceName)
	if f == nil {
		return nil
	}
	return f.GetRegion(lineNumber)
}

func (p *jsaDiffProvider) GetSourceMapping(sourceName *string, lineNumber int, columnNumber int) *sourcemap.Mapping_OriginalMapping {
	if sourceName == nil || strings.HasPrefix(*sourceName, "orig/") || lineNumber%5 != 3 {
		return nil
	}
	return sourcemap.Mapping_OriginalMapping_NewBuilder().
		SetOriginalFile("orig/" + *sourceName).
		SetLineNumber(lineNumber + 1).
		SetColumnPosition(columnNumber + 1).
		Build()
}

// jsaCollectHandler is an ErrorHandler that renders each report.
type jsaCollectHandler func(level *CheckLevel, error *JSError)

func (h jsaCollectHandler) Report(level *CheckLevel, error *JSError) { h(level, error) }

// jsaDiff accumulates mismatches per section.
type jsaDiff struct {
	t        *testing.T
	checked  map[string]int
	failed   map[string]int
	reported int
}

func (d *jsaDiff) check(section string, row []string, what string, got, want string) {
	d.checked[section]++
	if got == want {
		return
	}
	d.failed[section]++
	if d.reported < 40 {
		d.reported++
		d.t.Errorf("%s %q: %s\n got: %q\nwant: %q", section, row[1], what, got, want)
	}
}

func (d *jsaDiff) expectCols(row []string, n int) bool {
	if len(row) != n {
		d.t.Errorf("%s row has %d columns, want %d: %q", row[0], len(row), n, row)
		return false
	}
	return true
}

func TestJscompADifferential(t *testing.T) {
	rows := jsaReadVectors(t)
	d := &jsaDiff{t: t, checked: map[string]int{}, failed: map[string]int{}}
	provider := &jsaDiffProvider{files: map[string]*SourceFile{}}
	var manager *SortingErrorManagerBase
	loaders := map[string]*ModuleLoader{}
	var loaderErrors []string
	jsdocIndex := 0
	for _, row := range rows {
		switch row[0] {
		case "src":
			provider.files[row[1]] = SourceFile_FromCode(row[1], row[2])
			manager = NewSortingErrorManager(nil)
		case "fmt":
			if d.expectCols(row, 21) {
				jsaDiffFormat(t, d, row, provider, manager)
			}
		case "report":
			if d.expectCols(row, 8) {
				jsaDiffReport(t, d, row, provider, manager)
			}
		case "sorting":
			if d.expectCols(row, 4) {
				jsaDiffSorting(t, d, row)
			}
		case "guard":
			jsaDiffGuard(t, d, row)
		case "maperr":
			if d.expectCols(row, 4) {
				jsaDiffMapError(d, row)
			}
		case "modnames":
			if d.expectCols(row, 8) {
				fns := []struct {
					name string
					f    func(string) string
				}{
					{"escapePath", ModuleNames_EscapePath},
					{"canonicalizePath", ModuleNames_CanonicalizePath},
					{"toModuleName", ModuleNames_ToModuleName},
					{"toJSIdentifier", ModuleNames_ToJSIdentifier},
					{"fileToModuleName", ModuleNames_FileToModuleName},
					{"fileToJsIdentifier", ModuleNames_FileToJsIdentifier},
				}
				for i, fn := range fns {
					got := jsaSafe(func() string { return fn.f(row[1]) })
					d.check("modnames", row, fn.name, got, row[2+i])
				}
			}
		case "relpath":
			if d.expectCols(row, 4) {
				got := jsaSafe(func() string { return ModuleLoader_RelativePathFrom(row[1], row[2]) })
				d.check("relpath", row, "relativePathFrom", got, row[3])
			}
		case "modsetup":
			if d.expectCols(row, 7) {
				loader := jsaDiffModuleSetup(d, row, &loaderErrors)
				if loader != nil {
					loaders[row[1]] = loader
				}
			}
		case "modres":
			if d.expectCols(row, 9) {
				jsaDiffModuleResolve(d, row, loaders[row[1]], &loaderErrors)
			}
		case "adddep":
			if d.expectCols(row, 5) {
				info := jsaDepInfo(row[1])
				sb := jutil.NewStringBuilder()
				DependencyInfo_Util_WriteAddDependency(sb, info)
				d.check("adddep", row, "writeAddDependency", sb.String(), row[2])
				d.check("adddep", row, "toString", info.String(), row[3])
				d.check("adddep", row, "isEs6Module/isGoogModule",
					fmt.Sprintf("%t/%t", info.IsEs6Module(), info.IsGoogModule()), row[4])
			}
		case "deps":
			if d.expectCols(row, 7) {
				jsaDiffDeps(t, d, row)
			}
		case "smr":
			if d.expectCols(row, 5) {
				name := jsaSafe(func() string { return SourceMapResolver_GetRelativePath(row[1], row[2]).GetName() })
				d.check("smr", row, "getRelativePath", name, row[3])
				inline := jsaSafe(func() string {
					sf := SourceMapResolver_ExtractSourceMap(SourceFile_FromCode(row[1]+"x", ""), row[2], true)
					if sf == nil {
						return jsaNull
					}
					if strings.HasSuffix(sf.GetName(), ".inline.map") {
						return sf.GetName() + " " + sf.GetCode()
					}
					return sf.GetName()
				})
				d.check("smr", row, "extractSourceMap", inline, row[4])
			}
		case "jsdoc":
			if d.expectCols(row, 4) {
				jsaDiffJsdoc(t, d, row, jsdocIndex)
			}
			jsdocIndex++
		default:
			t.Fatalf("unknown section %q", row[0])
		}
	}
	total := 0
	for section, n := range d.checked {
		total += n
		if d.failed[section] > 0 {
			t.Errorf("%s: %d of %d checks differ", section, d.failed[section], n)
		}
	}
	if total < 5000 {
		t.Errorf("only %d checks ran", total)
	}
	t.Logf("%d rows, %d checks: %v", len(rows), total, d.checked)
}

func jsaDiffMakeError(t *testing.T, provider *jsaDiffProvider, row []string) *JSError {
	name := jsaNullable(row[1])
	line := jsaAtoi(t, row[2])
	charno := jsaAtoi(t, row[3])
	length := jsaAtoi(t, row[4])
	type_ := jsaDiffTypes[jsaAtoi(t, row[6])]
	arg := "arg" + strconv.Itoa(line)
	if row[5] == "0" {
		return JSError_MakeWithSourceLocation(name, line, charno, type_, arg)
	}
	n := rhino.Node_NewString("x").SetLinenoCharno(line, charno)
	n.SetLength(length)
	if name != nil {
		n.SetStaticSourceFile(provider.files[*name])
	}
	return JSError_MakeWithNode(n, type_, arg)
}

func jsaDiffFormat(t *testing.T, d *jsaDiff, row []string, provider *jsaDiffProvider, manager *SortingErrorManagerBase) {
	error := jsaDiffMakeError(t, provider, row)
	col := 8
	for _, f := range ErrorFormat_values() {
		for _, colorize := range []bool{false, true} {
			if colorize && f != ErrorFormat_SINGLELINE {
				continue
			}
			formatter := f.ToFormatter(provider, colorize)
			what := fmt.Sprintf("%s colorize=%t", f, colorize)
			d.check("fmt", row, what+" formatError", jsaSafe(func() string { return formatter.FormatError(error) }), row[col])
			d.check("fmt", row, what+" formatWarning", jsaSafe(func() string { return formatter.FormatWarning(error) }), row[col+1])
			col += 2
		}
	}
	verbose := NewVerboseMessageFormatter(provider)
	d.check("fmt", row, "verbose formatError", jsaSafe(func() string { return verbose.FormatError(error) }), row[col])
	noLevel := NewLightweightMessageFormatter(provider).SetIncludeLevel(false).SetIncludeLocation(false)
	d.check("fmt", row, "noLevel formatWarning", jsaSafe(func() string { return noLevel.FormatWarning(error) }), row[col+1])
	d.check("fmt", row, "getExcerptWithPosition",
		jsaSafe(func() string { return NewLightweightMessageFormatter(provider).GetExcerptWithPosition(error) }), row[col+2])
	manager.Report(jsaDiffLevels[jsaAtoi(t, row[7])].Ptr(), error)
}

func jsaDiffReport(t *testing.T, d *jsaDiff, row []string, provider *jsaDiffProvider, manager *SortingErrorManagerBase) {
	manager.SetTypedPercent([]float64{0.0, 12.25, 99.95}[jsaAtoi(t, row[2])])
	var json, ps1, ps3, ps2 bytes.Buffer
	NewJsonErrorReportGenerator(jutil.NewPrintStream(&json), provider).GenerateReport(manager)
	NewPrintStreamErrorReportGenerator(
		ErrorFormat_SINGLELINE.ToFormatter(provider, false), jutil.NewPrintStream(&ps1), 1).
		GenerateReport(manager)
	NewPrintStreamErrorReportGenerator(
		ErrorFormat_FULL.ToFormatter(provider, false), jutil.NewPrintStream(&ps3), 3).
		GenerateReport(manager)
	pem := NewPrintStreamErrorManager(jutil.NewPrintStream(&ps2))
	pem.SetSummaryDetailLevel(2)
	for _, e := range manager.GetErrors() {
		pem.Report(CheckLevel_ERROR.Ptr(), e)
	}
	for _, e := range manager.GetWarnings() {
		pem.Report(CheckLevel_WARNING.Ptr(), e)
	}
	pem.SetTypedPercent(manager.GetTypedPercent())
	pem.GenerateReport()
	d.check("report", row, "JsonErrorReportGenerator", json.String(), row[3])
	d.check("report", row, "PrintStreamErrorReportGenerator SINGLELINE/1", ps1.String(), row[4])
	d.check("report", row, "PrintStreamErrorReportGenerator FULL/3", ps3.String(), row[5])
	d.check("report", row, "PrintStreamErrorManager/2", ps2.String(), row[6])
	d.check("report", row, "counts",
		fmt.Sprintf("%d/%d/%t", manager.GetErrorCount(), manager.GetWarningCount(), manager.HasHaltingErrors()), row[7])
}

func jsaDiffSource(s string) *string {
	if s == "~" {
		return nil
	}
	return &s
}

func jsaDiffSorting(t *testing.T, d *jsaDiff, row []string) {
	m := NewSortingErrorManager(nil)
	for _, e := range strings.Split(row[1], ";") {
		f := strings.Split(e, "|")
		m.Report(jsaDiffLevels[jsaAtoi(t, f[4])].Ptr(), JSError_MakeWithSourceLocation(
			jsaDiffSource(f[0]), jsaAtoi(t, f[1]), jsaAtoi(t, f[2]), jsaDiffTypes[jsaAtoi(t, f[3])], "a"))
	}
	var sorted strings.Builder
	for _, e := range m.GetSortedDiagnostics() {
		fmt.Fprintf(&sorted, "%s %s:%d:%d %s\n", jsaLevelStr(e.Level), jsaJavaStr(e.Error.SourceName()),
			e.Error.GetLineNumber(), e.Error.Charno(), e.Error.Description())
	}
	d.check("sorting", row, "getSortedDiagnostics", sorted.String(), row[2])
	d.check("sorting", row, "counts", fmt.Sprintf("%d/%d/%t/%d/%d", m.GetErrorCount(), m.GetWarningCount(),
		m.HasHaltingErrors(), len(m.GetErrors()), len(m.GetWarnings())), row[3])
}

// jsaSplitTop splits a compose(...) body at top-level ';'.
func jsaSplitTop(s string) []string {
	var parts []string
	depth := 0
	start := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
		case ';':
			if depth == 0 {
				parts = append(parts, s[start:i])
				start = i + 1
			}
		}
	}
	if start < len(s) {
		parts = append(parts, s[start:])
	}
	return parts
}

// jsaParseGuard builds a guard from the oracle's spec grammar:
// bypath:F|X:LEVEL:paths, show:I|E:paths, strict, dg:group:LEVEL,
// dgps:group:part, summary:LEVEL, cu, compose(spec;...;=k) where =k repeats
// the k-th guard of the list.
func jsaParseGuard(spec string) WarningsGuard {
	f := strings.Split(spec, ":")
	group := func(s string) *DiagnosticGroup {
		n, _ := strconv.Atoi(s)
		return jsaDiffGroups[n]
	}
	switch f[0] {
	case "bypath":
		paths := strings.Split(f[3], ",")
		level := jsaCheckLevelValueOf(f[2])
		if f[1] == "F" {
			return ByPathWarningsGuard_ForPath(paths, level)
		}
		return ByPathWarningsGuard_ExceptPath(paths, level)
	case "show":
		showType := ShowByPathWarningsGuard_ShowType_EXCLUDE
		if f[1] == "I" {
			showType = ShowByPathWarningsGuard_ShowType_INCLUDE
		}
		return NewShowByPathWarningsGuardWithPathsAndShowType(strings.Split(f[2], ","), showType)
	case "strict":
		return NewStrictWarningsGuard()
	case "dg":
		return NewDiagnosticGroupWarningsGuard(group(f[1]), jsaCheckLevelValueOf(f[2]))
	case "dgps":
		return NewDiagnosticGroupPathSuppressingWarningsGuard(group(f[1]), f[2])
	case "summary":
		return NewCheckTypeSummaryWarningsGuard(jsaCheckLevelValueOf(f[1]))
	case "cu":
		return NewClosureUnawareCodeWarningsGuard()
	}
	inner := spec[len("compose(") : len(spec)-1]
	var guards []WarningsGuard
	for _, part := range jsaSplitTop(inner) {
		if strings.HasPrefix(part, "=") {
			k, _ := strconv.Atoi(part[1:])
			guards = append(guards, guards[k])
		} else {
			guards = append(guards, jsaParseGuard(part))
		}
	}
	return NewComposeWarningsGuardFromList(guards)
}

func jsaDiffGuard(t *testing.T, d *jsaDiff, row []string) {
	var guard WarningsGuard
	parsed := jsaSafe(func() string { guard = jsaParseGuard(row[1]); return "" })
	if len(row) == 4 {
		d.check("guard", row, "construction", parsed, row[3])
		return
	}
	if !d.expectCols(row, 7) {
		return
	}
	if parsed != "" {
		d.check("guard", row, "construction", parsed, "")
		return
	}
	var levels []string
	for _, e := range strings.Split(row[2], ";") {
		f := strings.Split(e, "|")
		err := JSError_MakeWithSourceLocation(jsaDiffSource(f[0]), jsaAtoi(t, f[2]), -1, jsaDiffTypes[jsaAtoi(t, f[1])])
		levels = append(levels, jsaSafe(func() string { return jsaLevelStr(guard.Level(err)) }))
	}
	d.check("guard", row, "level", strings.Join(levels, ";"), row[3])
	var must []string
	for _, g := range jsaDiffGroups {
		must = append(must, jsaSafe(func() string { return guard.MustRunChecks(g).String() }))
	}
	d.check("guard", row, "mustRunChecks", strings.Join(must, ";"), row[4])
	d.check("guard", row, "getPriority", strconv.Itoa(guard.GetPriority()), row[5])
	count := "-"
	if cg, ok := guard.(*ComposeWarningsGuard); ok {
		count = strconv.Itoa(len(cg.GetGuards()))
	}
	d.check("guard", row, "getGuards", count, row[6])
}

func jsaDiffMapError(d *jsaDiff, row []string) {
	message := row[1]
	key := jsaNull
	if ty := RhinoErrorReporter_MapError(message); ty != nil {
		key = ty.key
	}
	d.check("maperr", row, "mapError", key, row[2])
	var reported []string
	h := jsaCollectHandler(func(level *CheckLevel, error *JSError) {
		reported = append(reported, fmt.Sprintf("%s %s %s %s:%d:%d %s", jsaLevelStr(level), error.Type().key,
			error.DefaultLevel(), jsaJavaStr(error.SourceName()), error.GetLineNumber(), error.Charno(),
			error.Description()))
	})
	reporter := RhinoErrorReporter_ForOldRhino(h)
	reporter.Warning(message, "w.js", 3, 4)
	reporter.Error(message, "e.js", 5, 6)
	d.check("maperr", row, "forOldRhino", strings.Join(reported, "\n"), row[3])
}

func jsaDiffModuleFactory(spec string) ModuleLoader_ModuleResolverFactory {
	kind, m, _ := strings.Cut(spec, ":")
	switch kind {
	case "browser":
		return BrowserModuleResolver_FACTORY
	case "nodenull":
		return NewNodeModuleResolver_Factory()
	case "node":
		return NewNodeModuleResolver_FactoryWithEntries(jsaParseMapSpec(m))
	case "prefix":
		return NewBrowserWithTransformedPrefixesModuleResolver_Factory(jsaParseMapSpec(m))
	case "webpack":
		return NewWebpackModuleResolver_Factory(jsaParseMapSpec(m))
	}
	panic("bad factory spec " + spec)
}

func jsaDiffModuleSetup(d *jsaDiff, row []string, errors *[]string) *ModuleLoader {
	var inputs []DependencyInfo
	for _, in := range strings.Split(row[4], "\x01") {
		inputs = append(inputs, SimpleDependencyInfo_Builder_(in, in).Build())
	}
	escaper := ModuleLoader_PathEscaper_CANONICALIZE_ONLY
	if row[5] == "true" {
		escaper = ModuleLoader_PathEscaper_ESCAPE
	}
	var loader *ModuleLoader
	got := jsaSafe(func() string {
		loader = ModuleLoader_Builder_().
			SetModuleRoots(jsaSplit(row[3], ",")).
			SetInputs(inputs).
			SetFactory(jsaDiffModuleFactory(row[2])).
			SetErrorHandler(jsaCollectHandler(func(level *CheckLevel, error *JSError) {
				*errors = append(*errors, jsaLevelStr(level)+" "+error.Type().key+" "+error.Description())
			})).
			SetPathEscaper(escaper).
			Build()
		return jsaMapSpec(loader.GetPackageJsonMainEntries())
	})
	d.check("modsetup", row, "build/getPackageJsonMainEntries", got, row[6])
	if strings.HasPrefix(got, "!") {
		return nil
	}
	return loader
}

func jsaDiffModuleResolve(d *jsaDiff, row []string, loader *ModuleLoader, errors *[]string) {
	if loader == nil {
		d.check("modres", row, "loader", "missing", row[1])
		return
	}
	from, addr := row[2], row[3]
	*errors = nil
	var resolved, js, silent, asPath string
	func() {
		defer func() {
			if r := recover(); r != nil {
				resolved = jsaExc(r)
				js, silent, asPath = "", "", ""
			}
		}()
		p := loader.Resolve(from)
		resolved = p.String() + " " + p.ToModuleName()
		js = jsaNull
		if j := p.ResolveJsModule(addr, "s.js", 1, 2); j != nil {
			js = j.String()
		}
		silent = jsaNull
		if s := p.ResolveJsModuleSilently(addr); s != nil {
			silent = s.String()
		}
		asPath = p.ResolveModuleAsPath(addr).String()
	}()
	d.check("modres", row, "resolve", resolved, row[4])
	d.check("modres", row, "resolveJsModule", js, row[5])
	d.check("modres", row, "resolveJsModuleSilently", silent, row[6])
	d.check("modres", row, "resolveModuleAsPath", asPath, row[7])
	d.check("modres", row, "errors", strings.Join(*errors, "\n"), row[8])
}

// jsaDepInfo builds a SimpleDependencyInfo from the oracle's
// "name>provides|requires|typeRequires|loadFlags" spec.
func jsaDepInfo(spec string) *SimpleDependencyInfo {
	name, rest, _ := strings.Cut(spec, ">")
	f := strings.Split(rest, "|")
	var requires []*DependencyInfo_Require
	for _, q := range jsaSplit(f[1], ",") {
		requires = append(requires, DependencyInfo_Require_GoogRequireSymbol(q))
	}
	return SimpleDependencyInfo_Builder_("rel/"+name, name).
		SetProvides(jsaSplit(f[0], ",")).
		SetTypeRequires(jsaSplit(f[2], ",")).
		SetRequires(requires).
		SetLoadFlags(jsaParseMapSpec(f[3])).
		Build()
}

func jsaDepNames(l []*SimpleDependencyInfo) string {
	var out []string
	for _, d := range l {
		out = append(out, d.GetName())
	}
	return strings.Join(out, ",")
}

func jsaDiffDeps(t *testing.T, d *jsaDiff, row []string) {
	var infos []*SimpleDependencyInfo
	for _, e := range strings.Split(row[1], ";") {
		infos = append(infos, jsaDepInfo(e))
	}
	var roots []*SimpleDependencyInfo
	for _, i := range strings.Split(row[2], ",") {
		roots = append(roots, infos[jsaAtoi(t, i)])
	}
	var sorted, strong, weak, providing string
	func() {
		defer func() {
			if r := recover(); r != nil {
				sorted = jsaExc(r)
				strong, weak, providing = "", "", ""
			}
		}()
		sd := NewSortedDependencies(infos)
		sorted = jsaDepNames(sd.GetSortedList())
		strong = jsaDepNames(sd.GetSortedStrongDependenciesOf(roots))
		weak = jsaDepNames(sd.GetSortedWeakDependenciesOf(roots))
		var pb strings.Builder
		for k := 0; k < 8; k++ {
			if p := sd.MaybeGetInputProviding("n" + strconv.Itoa(k)); p == nil {
				pb.WriteString("-,")
			} else {
				pb.WriteString(p.GetName() + ",")
			}
		}
		pb.WriteString(jsaDepNames(sd.GetInputsWithoutProvides()))
		providing = pb.String()
	}()
	d.check("deps", row, "getSortedList", sorted, row[3])
	d.check("deps", row, "getSortedStrongDependenciesOf", strong, row[4])
	d.check("deps", row, "getSortedWeakDependenciesOf", weak, row[5])
	d.check("deps", row, "maybeGetInputProviding/getInputsWithoutProvides", providing, row[6])
}

func jsaJavaSet(s *jutil.LinkedHashSet[string]) string {
	if s == nil {
		return "[]"
	}
	return "[" + strings.Join(s.ToSlice(), ", ") + "]"
}

func jsaJsdocSummary(d *rhino.JSDocInfo) string {
	mode := "null"
	if m := d.GetPerFileClosureUnawareMode(); m != nil {
		mode = m.String()
	}
	return fmt.Sprintf("%t %t %s %s %s %s %s %t%t%t %t%t%t%t%t %s %t %s %s",
		d.IsConstant(), d.IsNoSideEffects(), jsaJavaSet(d.GetModifies()), jsaJavaStr(d.GetDescription()),
		jsaJavaStr(d.GetMeaning()), jsaJavaStr(d.GetAlternateMessageId()), jsaJavaSet(d.GetSuppressions()),
		d.IsStableIdGenerator(), d.IsConsistentIdGenerator(), d.IsIdGenerator(),
		d.IsConstructor(), d.IsInterface(), d.IsAbstract(), d.IsNoInline(), d.IsTypeSummary(),
		jsaJavaStr(d.GetLicense()), d.HasFileOverview(),
		"["+strings.Join(d.GetThrowsAnnotations(), ", ")+"]", mode)
}

// jsaDiffJsdoc replays the oracle's JSDoc builder ops; index is the case
// number, which the description includes.
func jsaDiffJsdoc(t *testing.T, d *jsaDiff, row []string, index int) {
	b := rhino.JSDocInfo_Builder_()
	c := strconv.Itoa(index)
	for _, op := range jsaSplit(strings.TrimSuffix(row[1], ","), ",") {
		switch jsaAtoi(t, op) {
		case 0:
			b.RecordConstancy()
		case 1:
			b.RecordNoSideEffects()
		case 2:
			b.RecordModifies(jutil.NewLinkedHashSet("this", "arguments"))
		case 3:
			b.RecordDescription("desc " + c)
		case 4:
			b.RecordMeaning("meaning")
		case 5:
			b.RecordAlternateMessageId("alt")
		case 6:
			b.RecordSuppressions(jutil.NewLinkedHashSet("partialAlias", "messageConventions"))
		case 7:
			b.RecordSuppressions(jutil.NewLinkedHashSet("untranspilableFeatures", "other"))
		case 8:
			b.RecordStableIdGenerator()
		case 9:
			b.RecordConsistentIdGenerator()
		case 10:
			b.RecordIdGenerator()
		case 11:
			b.RecordConstructor()
		case 12:
			b.RecordInterface()
		case 13:
			b.RecordAbstract()
		case 14:
			b.RecordNoInline()
		case 15:
			b.RecordTypeSummary()
		case 16:
			b.AddLicense("lic")
		case 17:
			b.RecordFileOverview("fo")
		case 18:
			b.RecordThrowsAnnotation("{Error}")
		default:
			b.RecordClosureUnawareCodeWithMode(rhino.JSDocInfo_PerFileClosureUnawareMode_SIMPLE)
		}
	}
	info := b.Build()
	pool := StringPool_Builder_()
	var ser, de string
	func() {
		defer func() {
			if r := recover(); r != nil {
				ser, de = jsaExc(r), ""
			}
		}()
		ser = jsaNull
		if s := JSDocSerializer_SerializeJsdoc(info, pool); s != nil {
			ser = s.String()
		}
		de = jsaNull
		if dd := JSDocSerializer_ConvertJSDocInfoForOptimizations(info); dd != nil {
			de = jsaJsdocSummary(dd)
		}
	}()
	d.check("jsdoc", row, "serializeJsdoc", ser, row[2])
	d.check("jsdoc", row, "convertJSDocInfoForOptimizations", de, row[3])
}
