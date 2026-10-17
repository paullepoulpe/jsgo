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
	levelClass  = EnumClassOf("Args4jVectors$Level", []tLevel{tLevel_QUIET, tLevel_DEFAULT, tLevel_VERBOSE})
	modeClass   = EnumClassOf("Args4jVectors$Mode", []tMode{tMode_BROWSER, tMode_NODE, tMode_WEBPACK, tMode_BROWSER_WITH_TRANSFORMED_PREFIXES})
	fmtClass    = EnumClassOf("Args4jVectors$Fmt", []tFmt{tFmt_PRETTY_PRINT, tFmt_SINGLE_QUOTES})
	outFmtClass = EnumClassOf("Args4jVectors$OutFmt", []tOutFmt{tOutFmt_MD, tOutFmt_JSON})
)

// closureBool is CommandLineRunner.Flags.BooleanOptionHandler.
type closureBool struct{ OptionHandlerBase }

func newClosureBool(parser *CmdLineParser, option OptionDef, setter Setter) OptionHandler {
	h := &closureBool{}
	h.InitOptionHandlerBase(h, parser, option, setter)
	return h
}

func (h *closureBool) ParseArguments(params Parameters) int {
	var param *string
	jutil.Try(func() {
		p := params.GetParameter(0)
		param = &p
	}, func(e ICmdLineException) {
		param = nil
	})
	if param == nil {
		h.Setter.AddValue(true)
		return 0
	}
	lowerParam := jutil.StringToLowerCaseASCII(*param)
	switch lowerParam {
	case "true", "on", "yes", "1":
		h.Setter.AddValue(true)
	case "false", "off", "no", "0":
		h.Setter.AddValue(false)
	default:
		h.Setter.AddValue(true)
		return 0
	}
	return 1
}

func (h *closureBool) GetDefaultMetaVariable() *string { return nil }

var recorded []string

// recordingStringHandler is CommandLineRunner.Flags.JsOptionHandler with its
// MultiFlagSetter.
type recordingStringHandler struct{ StringOptionHandler }

func newRecordingStringHandler(parser *CmdLineParser, option OptionDef, setter Setter) OptionHandler {
	h := &recordingStringHandler{}
	h.InitStringOptionHandler(h, parser, option, &recording{proxy: setter})
	return h
}

type recording struct{ proxy Setter }

func (r *recording) IsMultiValued() bool                  { return r.proxy.IsMultiValued() }
func (r *recording) GetType() *Class                      { return r.proxy.GetType() }
func (r *recording) AsFieldSetter() *FieldSetter          { return r.proxy.AsFieldSetter() }
func (r *recording) AsAnnotatedElement() AnnotatedElement { return r.proxy.AsAnnotatedElement() }
func (r *recording) AddValue(v any) {
	value := v.(string)
	normalizedValue := value
	n := jutil.StringLength(value)
	if n > 0 {
		first := jutil.StringSubstring(value, 0, 1)
		if (first == "'" || first == "\"") && jutil.StringSubstringFrom(value, n-1) == first {
			normalizedValue = jutil.StringSubstring(value, 1, n-1)
		}
	}
	r.proxy.AddValue(normalizedValue)
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
	map_      *jutil.JavaHashMap[string, *string]
	nullMap   *jutil.JavaHashMap[string, *string]
	words     []string
	noUsage   string
	dep       *string
	forb      *string
	nullInt   *int32
	unicode   string
	arguments []string
}

func newBeanA() *beanA {
	return &beanA{zeta: "z", level: tLevel_DEFAULT, threads: 1, map_: jutil.NewJavaHashMap[string, *string](), noUsage: "nu", unicode: "é"}
}

func (b *beanA) Args4jFields() []AnnotatedField {
	const c = "Args4jVectors$BeanA"
	return []AnnotatedField{
		{Option: &Option{Name: "--zeta", Usage: "Last alphabetically but declared first"}, Field: StringField(c, "zeta", &b.zeta)},
		{Option: &Option{Name: "--help", Handler: newClosureBool, Help: true, Usage: "Displays this message on stdout and exit"}, Field: BooleanField(c, "help", &b.help)},
		{Option: &Option{Name: "--year", Usage: "shortcut for defining goog.FEATURESET_YEAR=YYYY. The minimum valid value of the browser year is 2012"}, Field: IntegerField(c, "year", &b.year)},
		{Option: &Option{Name: "--level", Aliases: []string{"-W"}, Usage: "Specifies the warning level to use. Options: QUIET, DEFAULT, VERBOSE"}, Field: EnumField(c, "level", levelClass, &b.level)},
		{Option: &Option{Name: "--mode", Hidden: true, Usage: "A hidden enum option"}, Field: EnumField(c, "mode", modeClass, &b.mode)},
		{Option: &Option{Name: "--nullmode", Usage: "A nullable enum option"}, Field: NullableEnumField(c, "nullMode", modeClass, &b.nullMode)},
		{Option: &Option{Name: "--fmt", Usage: "Formatting options; may be repeated"}, Field: ListField(c, "fmt", fmtClass, &b.fmt)},
		{Option: &Option{Name: "--js", Handler: newRecordingStringHandler,
			Usage: "The JavaScript filename. You may specify multiple. The flag name is optional, " +
				"because args are interpreted as files by default. You may also use minimatch-style " +
				"glob patterns. For example, use --js='**.js' --js='!**_test.js' to recursively " +
				"include all js files that do not end in _test.js"}, Field: ListField(c, "js", Class_String, &b.js)},
		{Option: &Option{Name: "--externs", Usage: "The file containing JavaScript externs. You may specify multiple"}, Field: ListField(c, "externs", Class_String, &b.externs)},
		{Option: &Option{Name: "--define", Aliases: []string{"--D", "-D"},
			Usage: "Override the value of a variable annotated @define.\nThe format is <name>[=<val>], " +
				"where <name> is the name of a @define variable and <val> is a boolean, number, or a " +
				"single-quoted string that contains no single quotes. If [=<val>] is omitted, the " +
				"variable is marked true"}, Field: ListField(c, "define", Class_String, &b.define)},
		{Option: &Option{Name: "--threads", Handler: NewIntOptionHandler, Usage: "Number of threads"}, Field: IntField(c, "threads", &b.threads)},
		{Option: &Option{Name: "--debug", Handler: newClosureBool, Usage: "Enable debugging options"}, Field: BooleanField(c, "debug", &b.debug)},
		{Option: &Option{Name: "--plain", Usage: "A boolean handled by args4j's own handler"}, Field: BooleanField(c, "plain", &b.plain)},
		{Option: &Option{Name: "--boxed", Usage: "A boxed Boolean"}, Field: NullableBooleanField(c, "boxed", &b.boxed)},
		{Option: &Option{Name: "--out", MetaVar: "FILE", Usage: "Primary output filename. If not specified, output is written to stdout"}, Field: StringField(c, "out", &b.out)},
		{Option: &Option{Name: "--prefix", Usage: "A nullable string"}, Field: NullableStringField(c, "prefix", &b.prefix)},
		{Option: &Option{Name: "--map", Usage: "Prefixes to replace in ES6 import paths before resolving."}, Field: StringMapField(c, "map", &b.map_)},
		{Option: &Option{Name: "--nullmap", Usage: "A map that starts null"}, Field: StringMapField(c, "nullMap", &b.nullMap)},
		{Option: &Option{Name: "--words", Handler: NewStringArrayOptionHandler, Usage: "Greedy words"}, Field: ListField(c, "words", Class_String, &b.words)},
		{Option: &Option{Name: "--nousage"}, Field: StringField(c, "noUsage", &b.noUsage)},
		{Option: &Option{Name: "--dep", Depends: []string{"--out"}, Usage: "Needs --out"}, Field: NullableStringField(c, "dep", &b.dep)},
		{Option: &Option{Name: "--forb", Forbids: []string{"--debug", "--plain"}, Usage: "Cannot be used with --debug"}, Field: NullableStringField(c, "forb", &b.forb)},
		{Option: &Option{Name: "--nullint", Usage: "A nullable Integer"}, Field: NullableIntegerField(c, "nullInt", &b.nullInt)},
		{Option: &Option{Name: "--ünï",
			Usage: "Non-ASCII option name: ünïcödé usage text that is long " +
				"enough to wrap around the column limit at least once   here"}, Field: StringField(c, "unicode", &b.unicode)},
		{Argument: &Argument{}, Field: ListField(c, "arguments", Class_String, &b.arguments)},
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

func (b *beanA) dump() string {
	return "zeta=" + b.zeta + " help=" + javaStr(b.help) + " year=" + javaStr(b.year) + " level=" + javaStr(b.level) +
		" mode=" + javaStr(b.mode) + " nullMode=" + javaStr(b.nullMode) + " fmt=" + javaList(b.fmt) +
		" js=" + javaList(b.js) + " externs=" + javaList(b.externs) + " define=" + javaList(b.define) +
		" threads=" + javaStr(b.threads) + " debug=" + javaStr(b.debug) + " plain=" + javaStr(b.plain) +
		" boxed=" + javaStr(b.boxed) + " out=" + b.out + " prefix=" + javaStr(b.prefix) + " map=" + javaStr(b.map_) +
		" nullMap=" + javaStr(b.nullMap) + " words=" + javaList(b.words) + " noUsage=" + b.noUsage +
		" dep=" + javaStr(b.dep) + " forb=" + javaStr(b.forb) + " nullInt=" + javaStr(b.nullInt) +
		" unicode=" + b.unicode + " arguments=" + javaList(b.arguments) + " recorded=" + javaList(recorded)
}

type beanB struct {
	mapping, reports, result string
	h                        bool
}

func (b *beanB) Args4jFields() []AnnotatedField {
	const c = "Args4jVectors$BeanB"
	return []AnnotatedField{
		{Option: &Option{Name: "--mapping_file", Required: true, Usage: "The file name of the mapping generated by the production instrumentation pass."}, Field: StringField(c, "mapping", &b.mapping)},
		{Option: &Option{Name: "--reports_directory", Required: true, Usage: "The folder/directory which contains all the reports created by the instrumented production code."}, Field: StringField(c, "reports", &b.reports)},
		{Option: &Option{Name: "--result_output", Required: true, Usage: "Use this flag to provide the name of the final report that will be generated bythis reporter."}, Field: StringField(c, "result", &b.result)},
		{Option: &Option{Name: "-h", Help: true, Usage: "Help"}, Field: BooleanField(c, "h", &b.h)},
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
	const c = "Args4jVectors$BeanC"
	return []AnnotatedField{
		{Option: &Option{Name: "--fix", Usage: "Fix lint warnings automatically"}, Field: BooleanField(c, "fix", &b.fix)},
		{Option: &Option{Name: "--resolution_mode", Usage: "Expected path resolution mode for resolving module specifiers."}, Field: EnumField(c, "resolutionMode", modeClass, &b.resolutionMode)},
		{Argument: &Argument{}, Field: ListField(c, "files", Class_String, &b.files)},
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
	const c = "Args4jVectors$BeanD"
	return []AnnotatedField{
		{Option: &Option{Name: "--output", Usage: "Sets the desired output format.\n" + "Options: md, json" + "Default: json"}, Field: EnumField(c, "output", outFmtClass, &b.output)},
		{Option: &Option{Name: "--help", Usage: "Usage"}, Field: BooleanField(c, "help", &b.help)},
		{Argument: &Argument{Index: 1, MetaVar: "SECOND", Usage: "The second argument"}, Field: NullableStringField(c, "second", &b.second)},
		{Argument: &Argument{Index: 0, Required: true, Usage: "The first argument"}, Field: NullableStringField(c, "first", &b.first)},
		{Argument: &Argument{Index: 2, Usage: "A boolean argument"}, Field: BooleanField(c, "flag", &b.flag)},
	}
}

func (b *beanD) dump() string {
	return "output=" + javaStr(b.output) + " help=" + javaStr(b.help) + " second=" + javaStr(b.second) +
		" first=" + javaStr(b.first) + " flag=" + javaStr(b.flag)
}

type beanDupName struct{ x, y *string }

func (b *beanDupName) Args4jFields() []AnnotatedField {
	const c = "Args4jVectors$BeanDupName"
	return []AnnotatedField{
		{Option: &Option{Name: "-x", Usage: "x"}, Field: NullableStringField(c, "x", &b.x)},
		{Option: &Option{Name: "-y", Aliases: []string{"-x"}, Usage: "y"}, Field: NullableStringField(c, "y", &b.y)},
	}
}

type beanDupIndex struct{ a, b *string }

func (b *beanDupIndex) Args4jFields() []AnnotatedField {
	const c = "Args4jVectors$BeanDupIndex"
	return []AnnotatedField{
		{Argument: &Argument{Index: 0}, Field: NullableStringField(c, "a", &b.a)},
		{Argument: &Argument{Index: 0}, Field: NullableStringField(c, "b", &b.b)},
	}
}

type beanFinal struct{ f string }

func (b *beanFinal) Args4jFields() []AnnotatedField {
	const c = "Args4jVectors$BeanFinal"
	return []AnnotatedField{
		{Option: &Option{Name: "-f", Usage: "final"}, Field: NewField(c, "f", Class_String, nil, true,
			func() any { return b.f }, func(v any) { b.f = v.(string) })},
	}
}

type beanGap struct{ a *string }

func (b *beanGap) Args4jFields() []AnnotatedField {
	return []AnnotatedField{
		{Argument: &Argument{Index: 1}, Field: NullableStringField("Args4jVectors$BeanGap", "a", &b.a)},
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
		parser.PrintUsageToWriterWithFilter(w, nil, OptionHandlerFilterFunc(func(optionHandler OptionHandler) bool {
			option := optionHandler.GetOptionHandlerBase().Option
			if namedOptionDef, ok := option.(*NamedOptionDef); ok {
				return !option.Hidden() && slices.Contains(vectorCategory, strings.TrimPrefix(namedOptionDef.Name(), "--"))
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
			jutil.Try(func() {
				parser.ParseArgument(args...)
				out = append(out, "ok")
			}, func(e ICmdLineException) {
				out = append(out, "error "+vectorEsc(e.Error()), "message "+vectorEsc(e.GetMessage()),
					"parser "+strconv.FormatBool(e.GetParser() == parser))
			})
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
		want := lines[i+1 : j]
		got := runVector(header)
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
