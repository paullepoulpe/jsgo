package args4j

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Go-only: the differential test of the args4j port against the args4j that
// the pinned closure-compiler.jar shades. testdata/args4j_vectors.txt.gz is
// written by tools/oracle/java/Args4jVectors.java; this file declares the
// same option beans (and the two handlers it copies from Closure's
// CommandLineRunner), replays every "## " header of the oracle output and
// compares the whole text.

type tLevel int32

const (
	tLevel_QUIET tLevel = iota
	tLevel_DEFAULT
	tLevel_VERBOSE
)

func (l tLevel) String() string { return [...]string{"QUIET", "DEFAULT", "VERBOSE"}[l] }

type tMode int32

const (
	tMode_BROWSER tMode = iota
	tMode_NODE
	tMode_WEBPACK
	tMode_BROWSER_WITH_TRANSFORMED_PREFIXES
)

func (m tMode) String() string {
	return [...]string{"BROWSER", "NODE", "WEBPACK", "BROWSER_WITH_TRANSFORMED_PREFIXES"}[m]
}

// tFmt overrides toString(): String() is toString(), Name() is name().
type tFmt int32

const (
	tFmt_PRETTY_PRINT tFmt = iota
	tFmt_SINGLE_QUOTES
)

func (f tFmt) Name() string   { return [...]string{"PRETTY_PRINT", "SINGLE_QUOTES"}[f] }
func (f tFmt) String() string { return strings.ToLower(f.Name()) }

type tOutFmt int32

const (
	tOutFmt_MD tOutFmt = iota
	tOutFmt_JSON
)

func (f tOutFmt) String() string { return [...]string{"MD", "JSON"}[f] }

var (
	levelClass  = EnumClassOf([]tLevel{tLevel_QUIET, tLevel_DEFAULT, tLevel_VERBOSE})
	modeClass   = EnumClassOf([]tMode{tMode_BROWSER, tMode_NODE, tMode_WEBPACK, tMode_BROWSER_WITH_TRANSFORMED_PREFIXES})
	fmtClass    = EnumClassOf([]tFmt{tFmt_PRETTY_PRINT, tFmt_SINGLE_QUOTES})
	outFmtClass = EnumClassOf([]tOutFmt{tOutFmt_MD, tOutFmt_JSON})
)

// newClosureBool is CommandLineRunner.Flags.BooleanOptionHandler.
func newClosureBool(parser *CmdLineParser, option *OptionDef, setter Setter) *OptionHandler {
	target := setter.(TypedSetter[bool])
	return NewOptionHandler(parser, option, setter, "", func(h *OptionHandler, params Parameters) int {
		param, hasParam := params.PeekParameter(0)
		if !hasParam {
			target.AddTypedValue(true)
			return 0
		}
		switch jutil.StringToLowerCaseASCII(param) {
		case "true", "on", "yes", "1":
			target.AddTypedValue(true)
		case "false", "off", "no", "0":
			target.AddTypedValue(false)
		default:
			target.AddTypedValue(true)
			return 0
		}
		return 1
	})
}

var recorded []string

// newRecordingStringHandler is CommandLineRunner.Flags.JsOptionHandler with its
// MultiFlagSetter.
func newRecordingStringHandler(parser *CmdLineParser, option *OptionDef, setter Setter) *OptionHandler {
	return NewStringOptionHandler(parser, option, &recording{proxy: setter.(TypedSetter[string])})
}

type recording struct{ proxy TypedSetter[string] }

func (r *recording) IsMultiValued() bool                { return r.proxy.IsMultiValued() }
func (r *recording) DefaultHandler() OptionHandlerClass { return r.proxy.DefaultHandler() }
func (r *recording) AddTypedValue(value string) {
	normalizedValue := value
	n := jutil.StringLength(value)
	if n > 0 {
		first := jutil.StringSubstring(value, 0, 1)
		if (first == "'" || first == "\"") && jutil.StringSubstringFrom(value, n-1) == first {
			normalizedValue = jutil.StringSubstring(value, 1, n-1)
		}
	}
	r.proxy.AddTypedValue(normalizedValue)
	recorded = append(recorded, "JS:"+normalizedValue)
}

type beanA struct {
	zeta      string
	help      bool
	year      int32
	level     tLevel
	mode      tMode
	nullMode  *tMode
	fmt       []tFmt
	js        []string
	externs   []string
	define    []string
	threads   int
	debug     bool
	plain     bool
	boxed     *bool
	out       string
	prefix    *string
	map_      *jutil.OrderedMap[string, string]
	nullMap   *jutil.OrderedMap[string, string]
	words     []string
	noUsage   string
	dep       *string
	forb      *string
	nullInt   *int32
	unicode   string
	arguments []string
}

func newBeanA() *beanA {
	return &beanA{zeta: "z", level: tLevel_DEFAULT, threads: 1, map_: jutil.NewOrderedMap[string, string](), noUsage: "nu", unicode: "é"}
}

func (b *beanA) Args4jFields() []AnnotatedField {
	return []AnnotatedField{
		{Option: &Option{Name: "--zeta", Usage: "Last alphabetically but declared first"}, Field: String(&b.zeta)},
		{Option: &Option{Name: "--help", Handler: newClosureBool, Help: true, Usage: "Displays this message on stdout and exit"}, Field: Bool(&b.help)},
		{Option: &Option{Name: "--year", Usage: "shortcut for defining goog.FEATURESET_YEAR=YYYY. The minimum valid value of the browser year is 2012"}, Field: Int32(&b.year)},
		{Option: &Option{Name: "--level", Aliases: []string{"-W"}, Usage: "Specifies the warning level to use. Options: QUIET, DEFAULT, VERBOSE"}, Field: ScalarField(levelClass, &b.level)},
		{Option: &Option{Name: "--mode", Hidden: true, Usage: "A hidden enum option"}, Field: ScalarField(modeClass, &b.mode)},
		{Option: &Option{Name: "--nullmode", Usage: "A nullable enum option"}, Field: NullableField(modeClass, &b.nullMode)},
		{Option: &Option{Name: "--fmt", Usage: "Formatting options; may be repeated"}, Field: SliceField(fmtClass, &b.fmt)},
		{Option: &Option{Name: "--js", Handler: newRecordingStringHandler,
			Usage: "The JavaScript filename. You may specify multiple. The flag name is optional, " +
				"because args are interpreted as files by default. You may also use minimatch-style " +
				"glob patterns. For example, use --js='**.js' --js='!**_test.js' to recursively " +
				"include all js files that do not end in _test.js"}, Field: StringSlice(&b.js)},
		{Option: &Option{Name: "--externs", Usage: "The file containing JavaScript externs. You may specify multiple"}, Field: StringSlice(&b.externs)},
		{Option: &Option{Name: "--define", Aliases: []string{"--D", "-D"},
			Usage: "Override the value of a variable annotated @define.\nThe format is <name>[=<val>], " +
				"where <name> is the name of a @define variable and <val> is a boolean, number, or a " +
				"single-quoted string that contains no single quotes. If [=<val>] is omitted, the " +
				"variable is marked true"}, Field: StringSlice(&b.define)},
		{Option: &Option{Name: "--threads", Handler: NewIntOptionHandler, Usage: "Number of threads"}, Field: Int(&b.threads)},
		{Option: &Option{Name: "--debug", Handler: newClosureBool, Usage: "Enable debugging options"}, Field: Bool(&b.debug)},
		{Option: &Option{Name: "--plain", Usage: "A boolean handled by args4j's own handler"}, Field: Bool(&b.plain)},
		{Option: &Option{Name: "--boxed", Usage: "A boxed Boolean"}, Field: NullableBool(&b.boxed)},
		{Option: &Option{Name: "--out", MetaVar: "FILE", Usage: "Primary output filename. If not specified, output is written to stdout"}, Field: String(&b.out)},
		{Option: &Option{Name: "--prefix", Usage: "A nullable string"}, Field: NullableString(&b.prefix)},
		{Option: &Option{Name: "--map", Usage: "Prefixes to replace in ES6 import paths before resolving."}, Field: StringMap(&b.map_)},
		{Option: &Option{Name: "--nullmap", Usage: "A map that starts null"}, Field: StringMap(&b.nullMap)},
		{Option: &Option{Name: "--words", Handler: NewStringArrayOptionHandler, Usage: "Greedy words"}, Field: StringSlice(&b.words)},
		{Option: &Option{Name: "--nousage"}, Field: String(&b.noUsage)},
		{Option: &Option{Name: "--dep", Depends: []string{"--out"}, Usage: "Needs --out"}, Field: NullableString(&b.dep)},
		{Option: &Option{Name: "--forb", Forbids: []string{"--debug", "--plain"}, Usage: "Cannot be used with --debug"}, Field: NullableString(&b.forb)},
		{Option: &Option{Name: "--nullint", Usage: "A nullable Integer"}, Field: NullableInt32(&b.nullInt)},
		{Option: &Option{Name: "--ünï",
			Usage: "Non-ASCII option name: ünïcödé usage text that is long " +
				"enough to wrap around the column limit at least once   here"}, Field: String(&b.unicode)},
		{Argument: &Argument{}, Field: StringSlice(&b.arguments)},
	}
}

func javaList[E any](l []E) string {
	parts := make([]string, len(l))
	for i, e := range l {
		parts[i] = jutil.ToJavaString(e)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

func javaStr(v any) string { return jutil.ToJavaString(v) }

func javaMap(m *jutil.OrderedMap[string, string]) string {
	if m == nil {
		return "null"
	}
	return formatStringMapDefault(m)
}

func (b *beanA) dump() string {
	return "zeta=" + b.zeta + " help=" + javaStr(b.help) + " year=" + javaStr(b.year) + " level=" + javaStr(b.level) +
		" mode=" + javaStr(b.mode) + " nullMode=" + javaStr(b.nullMode) + " fmt=" + javaList(b.fmt) +
		" js=" + javaList(b.js) + " externs=" + javaList(b.externs) + " define=" + javaList(b.define) +
		" threads=" + javaStr(b.threads) + " debug=" + javaStr(b.debug) + " plain=" + javaStr(b.plain) +
		" boxed=" + javaStr(b.boxed) + " out=" + b.out + " prefix=" + javaStr(b.prefix) + " map=" + javaMap(b.map_) +
		" nullMap=" + javaMap(b.nullMap) + " words=" + javaList(b.words) + " noUsage=" + b.noUsage +
		" dep=" + javaStr(b.dep) + " forb=" + javaStr(b.forb) + " nullInt=" + javaStr(b.nullInt) +
		" unicode=" + b.unicode + " arguments=" + javaList(b.arguments) + " recorded=" + javaList(recorded)
}

type beanB struct {
	mapping, reports, result string
	h                        bool
}

func (b *beanB) Args4jFields() []AnnotatedField {
	return []AnnotatedField{
		{Option: &Option{Name: "--mapping_file", Required: true, Usage: "The file name of the mapping generated by the production instrumentation pass."}, Field: String(&b.mapping)},
		{Option: &Option{Name: "--reports_directory", Required: true, Usage: "The folder/directory which contains all the reports created by the instrumented production code."}, Field: String(&b.reports)},
		{Option: &Option{Name: "--result_output", Required: true, Usage: "Use this flag to provide the name of the final report that will be generated bythis reporter."}, Field: String(&b.result)},
		{Option: &Option{Name: "-h", Help: true, Usage: "Help"}, Field: Bool(&b.h)},
	}
}

func (b *beanB) dump() string {
	return "mapping=" + b.mapping + " reports=" + b.reports + " result=" + b.result + " h=" + javaStr(b.h)
}

type beanC struct {
	fix            bool
	resolutionMode tMode
	files          []string
}

func (b *beanC) Args4jFields() []AnnotatedField {
	return []AnnotatedField{
		{Option: &Option{Name: "--fix", Usage: "Fix lint warnings automatically"}, Field: Bool(&b.fix)},
		{Option: &Option{Name: "--resolution_mode", Usage: "Expected path resolution mode for resolving module specifiers."}, Field: ScalarField(modeClass, &b.resolutionMode)},
		{Argument: &Argument{}, Field: StringSlice(&b.files)},
	}
}

func (b *beanC) dump() string {
	return "fix=" + javaStr(b.fix) + " resolutionMode=" + javaStr(b.resolutionMode) + " files=" + javaList(b.files)
}

type beanD struct {
	output        tOutFmt
	help          bool
	second, first *string
	flag          bool
}

func (b *beanD) Args4jFields() []AnnotatedField {
	return []AnnotatedField{
		{Option: &Option{Name: "--output", Usage: "Sets the desired output format.\n" + "Options: md, json" + "Default: json"}, Field: ScalarField(outFmtClass, &b.output)},
		{Option: &Option{Name: "--help", Usage: "Usage"}, Field: Bool(&b.help)},
		{Argument: &Argument{Index: 1, MetaVar: "SECOND", Usage: "The second argument"}, Field: NullableString(&b.second)},
		{Argument: &Argument{Index: 0, Required: true, Usage: "The first argument"}, Field: NullableString(&b.first)},
		{Argument: &Argument{Index: 2, Usage: "A boolean argument"}, Field: Bool(&b.flag)},
	}
}

func (b *beanD) dump() string {
	return "output=" + javaStr(b.output) + " help=" + javaStr(b.help) + " second=" + javaStr(b.second) +
		" first=" + javaStr(b.first) + " flag=" + javaStr(b.flag)
}

type beanDupName struct{ x, y *string }

func (b *beanDupName) Args4jFields() []AnnotatedField {
	return []AnnotatedField{
		{Option: &Option{Name: "-x", Usage: "x"}, Field: NullableString(&b.x)},
		{Option: &Option{Name: "-y", Aliases: []string{"-x"}, Usage: "y"}, Field: NullableString(&b.y)},
	}
}

type beanDupIndex struct{ a, b *string }

func (b *beanDupIndex) Args4jFields() []AnnotatedField {
	return []AnnotatedField{
		{Argument: &Argument{Index: 0}, Field: NullableString(&b.a)},
		{Argument: &Argument{Index: 0}, Field: NullableString(&b.b)},
	}
}

type beanFinal struct{ f string }

func finalField(declaringClass, name string) *Field {
	return &Field{
		finalErr: "Cannot set value to a final field '" + declaringClass + "." + name + "'.",
	}
}

func (b *beanFinal) Args4jFields() []AnnotatedField {
	return []AnnotatedField{
		{Option: &Option{Name: "-f", Usage: "final"}, Field: finalField("Args4jVectors$BeanFinal", "f")},
	}
}

type beanGap struct{ a *string }

func (b *beanGap) Args4jFields() []AnnotatedField {
	return []AnnotatedField{
		{Argument: &Argument{Index: 1}, Field: NullableString(&b.a)},
	}
}

func newVectorBean(name string) any {
	switch name {
	case "A":
		return newBeanA()
	case "B":
		return &beanB{}
	case "C":
		return &beanC{}
	case "D":
		return &beanD{output: tOutFmt_JSON}
	case "DupName":
		return &beanDupName{}
	case "DupIndex":
		return &beanDupIndex{}
	case "Final":
		return &beanFinal{}
	case "Gap":
		return &beanGap{}
	}
	panic(name)
}

func dumpVectorBean(bean any) string {
	if d, ok := bean.(interface{ dump() string }); ok {
		return d.dump()
	}
	return "-"
}

func vectorEsc(s string) string {
	var sb strings.Builder
	for _, c := range jutil.StringToUTF16(s) {
		switch {
		case c == '\\':
			sb.WriteString(`\\`)
		case c == '\t':
			sb.WriteString(`\t`)
		case c == '\n':
			sb.WriteString(`\n`)
		case c == '\r':
			sb.WriteString(`\r`)
		case c < 0x20 || c > 0x7e:
			fmt.Fprintf(&sb, `\u%04x`, c)
		default:
			sb.WriteByte(byte(c))
		}
	}
	return sb.String()
}

func vectorUnesc(s string) string {
	var units []uint16
	for i := 0; i < len(s); i++ {
		if s[i] != '\\' {
			units = append(units, uint16(s[i]))
			continue
		}
		i++
		switch s[i] {
		case 't':
			units = append(units, '\t')
		case 'n':
			units = append(units, '\n')
		case 'r':
			units = append(units, '\r')
		case 'u':
			v, _ := strconv.ParseUint(s[i+1:i+5], 16, 16)
			units = append(units, uint16(v))
			i += 4
		default:
			units = append(units, uint16(s[i]))
		}
	}
	return jutil.StringFromUTF16(units)
}

var vectorCategory = []string{"js", "level", "externs", "out", "define", "help"}

func vectorUsage(bean any, config string) string {
	props := ParserProperties_Defaults()
	switch config {
	case "nodefaults":
		props.WithShowDefaults(false)
	case "eqdelim":
		props.WithOptionValueDelimiter("=")
	case "unsorted":
		props.WithOptionSorter(nil)
	}
	parser := NewCmdLineParserWithProperties(bean, props)
	w := jutil.NewStringWriter()
	switch config {
	case "all":
		parser.PrintUsageToWriterWithFilter(w, nil, OptionHandlerFilter_ALL)
	case "w6", "w20", "w40", "w60", "w120", "w5000":
		n, _ := strconv.Atoi(config[1:])
		parser.SetUsageWidth(n)
		parser.PrintUsageToWriter(w, nil)
	case "category":
		parser.PrintUsageToWriterWithFilter(w, nil, OptionHandlerFilterFunc(func(optionHandler *OptionHandler) bool {
			option := optionHandler.Option
			if option.IsNamed() {
				return !option.Hidden() && slices.Contains(vectorCategory, strings.TrimPrefix(option.Name(), "--"))
			}
			return false
		}))
	case "example":
		w.Write("[" + parser.PrintExample(OptionHandlerFilter_ALL) + "]\n")
		w.Write("[" + parser.PrintExample(OptionHandlerFilter_REQUIRED) + "]\n")
		w.Write("[" + parser.PrintExample(OptionHandlerFilter_PUBLIC) + "]\n")
	case "single":
		parser.PrintSingleLineUsageToWriter(w, nil)
	default:
		parser.PrintUsageToWriter(w, nil)
	}
	return w.String()
}

func vectorBlock(out *[]string, text string) {
	for _, line := range strings.Split(text, "\n") {
		*out = append(*out, "| "+vectorEsc(line))
	}
}

func runVector(header string) []string {
	var out []string
	fields := strings.Split(header, "\t")
	words := strings.Fields(fields[0])
	switch words[1] {
	case "usage":
		bean, config := words[2], words[3]
		jutil.Try(func() {
			vectorBlock(&out, vectorUsage(newVectorBean(bean), config))
		}, func(t jutil.Throwable) {
			out = append(out, "throw "+vectorEsc(t.Error()))
		})
	case "parse":
		var args []string
		for _, f := range fields[1:] {
			args = append(args, vectorUnesc(f))
		}
		recorded = nil
		bean := newVectorBean(words[2])
		jutil.Try(func() {
			parser := NewCmdLineParser(bean)
			if e := parser.TryParseArgument(args...); e != nil {
				out = append(out, "error "+vectorEsc(e.Error()), "message "+vectorEsc(e.GetMessage()),
					"parser "+strconv.FormatBool(e.GetParser() == parser))
			} else {
				out = append(out, "ok")
			}
			out = append(out, "state "+vectorEsc(dumpVectorBean(bean)))
			if slices.Contains(args, "--help") {
				w := jutil.NewStringWriter()
				parser.PrintUsageToWriter(w, nil)
				vectorBlock(&out, w.String())
			}
		}, func(t jutil.Throwable) {
			out = append(out, "throw "+vectorEsc(t.Error()))
		})
	}
	return out
}

func sortStateMapField(line, startMarker, endMarker string) string {
	start := strings.Index(line, startMarker)
	if start < 0 {
		return line
	}
	contentStart := start + len(startMarker)
	end := strings.Index(line[contentStart:], endMarker)
	if end < 0 {
		return line
	}
	content := line[contentStart : contentStart+end]
	if !strings.Contains(content, ", ") {
		return line
	}
	parts := strings.Split(content, ", ")
	slices.Sort(parts)
	return line[:contentStart] + strings.Join(parts, ", ") + line[contentStart+end:]
}

func normalizeVectorLines(lines []string) []string {
	out := make([]string, len(lines))
	for i, line := range lines {
		if strings.HasPrefix(line, "state ") {
			line = sortStateMapField(line, " map={", "} nullMap=")
			line = sortStateMapField(line, " nullMap={", "} words=")
		}
		out[i] = line
	}
	return out
}

func TestArgs4jVectors(t *testing.T) {
	f, err := os.Open("testdata/args4j_vectors.txt.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	sc := bufio.NewScanner(gz)
	sc.Buffer(make([]byte, 1<<20), 1<<24)
	var lines []string
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	cases, diffs := 0, 0
	for i := 0; i < len(lines); {
		header := lines[i]
		if !strings.HasPrefix(header, "## ") {
			t.Fatalf("line %d: expected a header, got %q", i+1, header)
		}
		j := i + 1
		for j < len(lines) && !strings.HasPrefix(lines[j], "## ") {
			j++
		}
		want := normalizeVectorLines(lines[i+1 : j])
		got := normalizeVectorLines(runVector(header))
		cases++
		if !slices.Equal(got, want) {
			diffs++
			if diffs <= 10 {
				t.Errorf("%s\n got:\n  %s\nwant:\n  %s", header, strings.Join(got, "\n  "), strings.Join(want, "\n  "))
			}
		}
		i = j
	}
	if diffs > 0 {
		t.Errorf("%d of %d vectors differ", diffs, cases)
	}
	t.Logf("%d vectors", cases)
}
