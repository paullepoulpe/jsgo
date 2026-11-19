// Clean-room org.kohsuke.args4j runtime implementation for j2gort.
package j2gort

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

type args4jParserCore struct {
	bean       any
	options    []args4jHandlerEntry
	arguments  []args4jHandlerEntry
	byName     map[string]Ext_org_kohsuke_args4j_spi_OptionHandler
	usageWidth int
}

type args4jHandlerEntry struct {
	def     Ext_org_kohsuke_args4j_OptionDef
	handler Ext_org_kohsuke_args4j_spi_OptionHandler
}

type args4jOptionBinder interface {
	RegisterArgs4jOptions(p *Ext_org_kohsuke_args4j_CmdLineParser)
}

func initArgs4jParser(p *Ext_org_kohsuke_args4j_CmdLineParser, bean any) {
	core := &args4jParserCore{
		bean:       bean,
		byName:     make(map[string]Ext_org_kohsuke_args4j_spi_OptionHandler),
		usageWidth: 80,
	}
	p.Core_ = core
	if binder, ok := bean.(args4jOptionBinder); ok && !javalang.IsNil(bean) {
		binder.RegisterArgs4jOptions(p)
	}
}

func NewArgs4jStringMap() Ext_java_util_Map {
	return wrapLinkedHashMap(NewJMapCore(false))
}

func NewArgs4jList() Ext_java_util_List {
	return wrapArrayList(nil)
}

func Args4jMapPut(m Ext_java_util_Map, k, v string) {
	if mc, ok := any(m).(interface{ mc() *JMapCore }); ok {
		mc.mc().Put(BoxString(k), BoxString(v))
	}
}

func Args4jListAdd(l Ext_java_util_List, elem any) {
	if lc, ok := any(l).(interface{ lc() *JListCore }); ok {
		lc.lc().Add(elem)
	}
}

func (p *Ext_org_kohsuke_args4j_CmdLineParser) args4jCore() *args4jParserCore {
	if p == nil {
		return nil
	}
	if c, ok := p.Core_.(*args4jParserCore); ok && c != nil {
		return c
	}
	c := &args4jParserCore{byName: make(map[string]Ext_org_kohsuke_args4j_spi_OptionHandler), usageWidth: 80}
	p.Core_ = c
	return c
}

func NewArgs4jCmdLineException(p *Ext_org_kohsuke_args4j_CmdLineParser, msg string) Ext_org_kohsuke_args4j_CmdLineException {
	_ = p
	ex := &Ext_org_kohsuke_args4j_CmdLineExceptionBase{}
	ex.Self_ = ex
	javalang.InitThrowable(ex, "org.kohsuke.args4j.CmdLineException", msg, nil)
	return ex
}

func NewArgs4jCmdLineExceptionCause(p *Ext_org_kohsuke_args4j_CmdLineParser, cause javalang.Throwable) Ext_org_kohsuke_args4j_CmdLineException {
	_ = p
	ex := &Ext_org_kohsuke_args4j_CmdLineExceptionBase{}
	ex.Self_ = ex
	msg := ""
	if cause != nil {
		msg = cause.GetMessage()
	}
	javalang.InitThrowable(ex, "org.kohsuke.args4j.CmdLineException", msg, cause)
	return ex
}

type args4jOptionDefData struct {
	isNamed     bool
	name        string
	aliases     []string
	usage       string
	metaVar     string
	required    bool
	help        bool
	hidden      bool
	multiValued bool
	index       int
}

func NewArgs4jNamedOptionDef(
	name string,
	aliases []string,
	usage string,
	metaVar string,
	required bool,
	help bool,
	hidden bool,
	multiValued bool,
) *Ext_org_kohsuke_args4j_NamedOptionDef {
	d := &args4jOptionDefData{
		isNamed:     true,
		name:        name,
		aliases:     aliases,
		usage:       usage,
		metaVar:     metaVar,
		required:    required,
		help:        help,
		hidden:      hidden,
		multiValued: multiValued,
	}
	obj := &Ext_org_kohsuke_args4j_NamedOptionDef{}
	obj.Core_ = d
	obj.Ext_org_kohsuke_args4j_OptionDefBase.Core_ = d
	obj.Ext_org_kohsuke_args4j_OptionDefBase.Self_ = obj
	return obj
}

func NewArgs4jArgumentDef(
	usage string,
	metaVar string,
	required bool,
	hidden bool,
	multiValued bool,
	index int,
) Ext_org_kohsuke_args4j_OptionDef {
	d := &args4jOptionDefData{
		isNamed:     false,
		usage:       usage,
		metaVar:     metaVar,
		required:    required,
		hidden:      hidden,
		multiValued: multiValued,
		index:       index,
	}
	obj := &Ext_org_kohsuke_args4j_OptionDefBase{Core_: d}
	obj.Self_ = obj
	return obj
}

func getOptionDefData(core any) *args4jOptionDefData {
	if d, ok := core.(*args4jOptionDefData); ok && d != nil {
		return d
	}
	return &args4jOptionDefData{}
}

type args4jFieldSetterData struct {
	multiValued bool
	typ         *javalang.Class
	addFn       func(any)
	getFn       func() any
}

func NewArgs4jFieldSetter(multiValued bool, typ *javalang.Class, addFn func(any), getFn func() any) *Ext_org_kohsuke_args4j_spi_FieldSetter {
	d := &args4jFieldSetterData{
		multiValued: multiValued,
		typ:         typ,
		addFn:       addFn,
		getFn:       getFn,
	}
	return &Ext_org_kohsuke_args4j_spi_FieldSetter{Core_: d}
}

type args4jParametersImpl struct {
	p       *Ext_org_kohsuke_args4j_CmdLineParser
	optName string
	args    []string
	pos     int
}

func (*args4jParametersImpl) Implements_j2gort_Ext_org_kohsuke_args4j_spi_Parameters() {}

func (pr *args4jParametersImpl) GetParameter(idx int32) string {
	i := pr.pos + int(idx)
	if i < 0 || i >= len(pr.args) {
		if pr.optName != "" {
			panic(NewArgs4jCmdLineException(pr.p, fmt.Sprintf("Option %q takes an operand", pr.optName)))
		}
		panic(NewArgs4jCmdLineException(pr.p, "Option takes an operand"))
	}
	return pr.args[i]
}

func (pr *args4jParametersImpl) Size() int32 {
	return int32(len(pr.args) - pr.pos)
}

type args4jBuiltinHandler struct {
	Ext_org_kohsuke_args4j_spi_OptionHandlerBase
	kind    string
	enumFqn string
}

func NewArgs4jBuiltinHandler(
	p *Ext_org_kohsuke_args4j_CmdLineParser,
	opt Ext_org_kohsuke_args4j_OptionDef,
	setter Ext_org_kohsuke_args4j_spi_Setter,
	kind string,
	enumFqn string,
) Ext_org_kohsuke_args4j_spi_OptionHandler {
	h := &args4jBuiltinHandler{kind: kind, enumFqn: enumFqn}
	h.Self_ = h
	InitExt_org_kohsuke_args4j_spi_OptionHandler(&h.Ext_org_kohsuke_args4j_spi_OptionHandlerBase, p, opt, setter)
	return h
}

func (h *args4jBuiltinHandler) ParseArguments(params Ext_org_kohsuke_args4j_spi_Parameters) int32 {
	switch h.kind {
	case "boolean":
		h.F_Setter.AddValue(ValueOfBoolean(true))
		return 0
	case "explicit_boolean":
		var tok string
		hasTok := false
		func() {
			defer func() { _ = recover() }()
			tok = params.GetParameter(0)
			hasTok = true
		}()
		if !hasTok {
			h.F_Setter.AddValue(ValueOfBoolean(true))
			return 0
		}
		low := strings.ToLower(tok)
		switch low {
		case "true", "on", "yes", "1":
			h.F_Setter.AddValue(ValueOfBoolean(true))
			return 1
		case "false", "off", "no", "0":
			h.F_Setter.AddValue(ValueOfBoolean(false))
			return 1
		default:
			h.F_Setter.AddValue(ValueOfBoolean(true))
			return 0
		}
	default:
		s := params.GetParameter(0)
		h.F_Setter.AddValue(BoxString(s))
		return 1
	}
}

func (h *args4jBuiltinHandler) GetDefaultMetaVariable() string {
	switch h.kind {
	case "boolean":
		return NullString
	case "explicit_boolean":
		return "[true | false]"
	case "int", "double":
		return "N"
	case "map":
		return "key=value"
	case "enum":
		consts := LookupEnumConstantsByName(h.enumFqn)
		if len(consts) > 0 {
			names := make([]string, 0, len(consts))
			for _, c := range consts {
				if n, ok := c.(interface{ Name() string }); ok {
					names = append(names, n.Name())
				}
			}
			if len(names) > 0 {
				return "[" + strings.Join(names, " | ") + "]"
			}
		}
		return "VAL"
	default:
		return "VAL"
	}
}

func Args4jRegisterHandler(
	p *Ext_org_kohsuke_args4j_CmdLineParser,
	def Ext_org_kohsuke_args4j_OptionDef,
	handler Ext_org_kohsuke_args4j_spi_OptionHandler,
	isArgument bool,
) {
	c := p.args4jCore()
	if isArgument {
		c.arguments = append(c.arguments, args4jHandlerEntry{def: def, handler: handler})
		return
	}
	c.options = append(c.options, args4jHandlerEntry{def: def, handler: handler})
	if nd, ok := def.(*Ext_org_kohsuke_args4j_NamedOptionDef); ok && nd != nil {
		c.byName[nd.Name()] = handler
		d := getOptionDefData(nd.Core_)
		for _, al := range d.aliases {
			c.byName[al] = handler
		}
	}
}

func Args4jAsString(v any) string {
	if v == nil {
		return NullString
	}
	if s, ok := v.(string); ok {
		return s
	}
	return StrVal(v)
}

func Args4jAsBool(p *Ext_org_kohsuke_args4j_CmdLineParser, v any) bool {
	switch x := v.(type) {
	case bool:
		return x
	case *BoxBoolean:
		return x != nil && x.Val
	}
	s := strings.ToLower(strings.TrimSpace(Args4jAsString(v)))
	switch s {
	case "true", "on", "yes", "1":
		return true
	case "false", "off", "no", "0":
		return false
	default:
		panic(NewArgs4jCmdLineException(p, fmt.Sprintf("%q is not a valid boolean value", s)))
	}
}

func Args4jAsInt(p *Ext_org_kohsuke_args4j_CmdLineParser, v any) int32 {
	switch x := v.(type) {
	case int32:
		return x
	case *BoxInt:
		if x != nil {
			return x.Val
		}
	}
	s := Args4jAsString(v)
	n, err := strconv.ParseInt(strings.TrimSpace(s), 10, 32)
	if err != nil {
		panic(NewArgs4jCmdLineException(p, fmt.Sprintf("%q is not a valid number", s)))
	}
	return int32(n)
}

func Args4jAsLong(p *Ext_org_kohsuke_args4j_CmdLineParser, v any) int64 {
	switch x := v.(type) {
	case int64:
		return x
	case *BoxLong:
		if x != nil {
			return x.Val
		}
	}
	s := Args4jAsString(v)
	n, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil {
		panic(NewArgs4jCmdLineException(p, fmt.Sprintf("%q is not a valid number", s)))
	}
	return n
}

func Args4jAsDouble(p *Ext_org_kohsuke_args4j_CmdLineParser, v any) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case *BoxDouble:
		if x != nil {
			return x.Val
		}
	}
	s := Args4jAsString(v)
	n, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		panic(NewArgs4jCmdLineException(p, fmt.Sprintf("%q is not a valid number", s)))
	}
	return n
}

func Args4jParseEnum(p *Ext_org_kohsuke_args4j_CmdLineParser, fqn string, v any) any {
	if v != nil {
		if _, isStr := v.(string); !isStr {
			return v
		}
	}
	s := Args4jAsString(v)
	norm := strings.ReplaceAll(s, "-", "_")
	consts := LookupEnumConstantsByName(fqn)
	for _, c := range consts {
		if n, ok := c.(interface{ Name() string }); ok && n.Name() == s {
			return c
		}
	}
	for _, c := range consts {
		if n, ok := c.(interface{ Name() string }); ok && (strings.EqualFold(n.Name(), s) || strings.EqualFold(n.Name(), norm)) {
			return c
		}
	}
	panic(NewArgs4jCmdLineException(p, fmt.Sprintf("%q is not a valid value", s)))
}

func Args4jSplitMapEntry(p *Ext_org_kohsuke_args4j_CmdLineParser, v any) (string, string) {
	s := Args4jAsString(v)
	idx := strings.IndexByte(s, '=')
	if idx < 0 {
		panic(NewArgs4jCmdLineException(p, "An argument for setting a Map must contain a \"=\""))
	}
	return s[:idx], s[idx+1:]
}

func args4jParseTokens(p *Ext_org_kohsuke_args4j_CmdLineParser, rawTokens []string) {
	c := p.args4jCore()
	tokens := make([]string, 0, len(rawTokens))
	for _, t := range rawTokens {
		if strings.HasPrefix(t, "@") && len(t) > 1 {
			filePath := t[1:]
			data, err := os.ReadFile(filePath)
			if err != nil {
				panic(NewArgs4jCmdLineException(p, "No such file: "+filePath))
			}
			for _, line := range strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n") {
				trimmed := strings.TrimSpace(line)
				if trimmed == "" || strings.HasPrefix(trimmed, "#") {
					continue
				}
				tokens = append(tokens, trimmed)
			}
		} else {
			tokens = append(tokens, t)
		}
	}

	endOfOpts := false
	i := 0
	argHandlerIdx := 0
	for i < len(tokens) {
		tok := tokens[i]
		if !endOfOpts && tok == "--" {
			endOfOpts = true
			i++
			continue
		}
		if !endOfOpts && strings.HasPrefix(tok, "-") && len(tok) > 1 {
			// Look up exact option name first, then --key=value split
			if h, ok := c.byName[tok]; ok {
				pr := &args4jParametersImpl{p: p, optName: tok, args: tokens, pos: i + 1}
				consumed := int(h.ParseArguments(pr))
				i += 1 + consumed
				continue
			}
			if eq := strings.IndexByte(tok, '='); eq > 0 {
				optKey := tok[:eq]
				optVal := tok[eq+1:]
				if len(optVal) >= 2 && ((optVal[0] == '\'' && optVal[len(optVal)-1] == '\'') || (optVal[0] == '"' && optVal[len(optVal)-1] == '"')) {
					optVal = optVal[1 : len(optVal)-1]
				}
				if h, ok := c.byName[optKey]; ok {
					pr := &args4jParametersImpl{p: p, optName: optKey, args: []string{optVal}, pos: 0}
					_ = h.ParseArguments(pr)
					i++
					continue
				}
			}
			panic(NewArgs4jCmdLineException(p, fmt.Sprintf("%q is not a valid option", tok)))
		}

		// Positional argument
		if len(c.arguments) == 0 {
			panic(NewArgs4jCmdLineException(p, fmt.Sprintf("No argument is allowed: %s", tok)))
		}
		idx := argHandlerIdx
		if idx >= len(c.arguments) {
			idx = len(c.arguments) - 1
		}
		ent := c.arguments[idx]
		pr := &args4jParametersImpl{p: p, args: tokens, pos: i}
		consumed := int(ent.handler.ParseArguments(pr))
		if consumed <= 0 {
			consumed = 1
		}
		i += consumed
		if !getOptionDefData(ent.def.GetExt_org_kohsuke_args4j_OptionDefBase().Core_).multiValued {
			argHandlerIdx++
		}
	}
}

func args4jDefaultSuffix(h Ext_org_kohsuke_args4j_spi_OptionHandler) string {
	if h == nil {
		return ""
	}
	setter := h.GetExt_org_kohsuke_args4j_spi_OptionHandlerBase().F_Setter
	fs, ok := setter.(*Ext_org_kohsuke_args4j_spi_FieldSetter)
	if !ok || fs == nil {
		return ""
	}
	d, ok := fs.Core_.(*args4jFieldSetterData)
	if !ok || d == nil || d.multiValued || d.getFn == nil {
		return ""
	}
	val := d.getFn()
	if javalang.IsNil(val) {
		return ""
	}
	return " (default: " + StrVal(val) + ")"
}

func args4jOptionHead(ent args4jHandlerEntry) (string, *args4jOptionDefData, bool) {
	nd, ok := ent.def.(*Ext_org_kohsuke_args4j_NamedOptionDef)
	if !ok || nd == nil {
		return "", nil, false
	}
	d := getOptionDefData(nd.Core_)
	if d.usage == "" {
		return "", d, false
	}
	var hb strings.Builder
	hb.WriteString(" ")
	hb.WriteString(d.name)
	if len(d.aliases) > 0 {
		hb.WriteString(" (")
		hb.WriteString(strings.Join(d.aliases, ", "))
		hb.WriteString(")")
	}
	meta := d.metaVar
	if meta == "" {
		mv := ent.handler.GetDefaultMetaVariable()
		if mv != NullString && mv != "" {
			meta = mv
		}
	}
	if meta != "" {
		hb.WriteString(" ")
		hb.WriteString(meta)
	}
	return hb.String(), d, true
}

func args4jWrapLines(s string, w int, indentCont bool) []string {
	if w <= 0 {
		return strings.Split(s, "\n")
	}
	var out []string
	for _, para := range strings.Split(s, "\n") {
		cur := para
		for len(cur) > w {
			idx := strings.LastIndexByte(cur[:w], ' ')
			if idx <= w*3/5 {
				out = append(out, cur[:w])
				rest := strings.TrimPrefix(cur[w:], " ")
				if indentCont {
					cur = " " + rest
				} else {
					cur = rest
				}
			} else {
				out = append(out, cur[:idx])
				if indentCont {
					cur = " " + cur[idx+1:]
				} else {
					cur = cur[idx+1:]
				}
			}
		}
		out = append(out, cur)
	}
	return out
}

func args4jPrintUsage(
	p *Ext_org_kohsuke_args4j_CmdLineParser,
	w javaio.Writer,
	filter Ext_org_kohsuke_args4j_OptionHandlerFilter,
) {
	if w == nil {
		return
	}
	c := p.args4jCore()
	opts := make([]args4jHandlerEntry, len(c.options))
	copy(opts, c.options)
	sort.SliceStable(opts, func(i, j int) bool {
		ni, nj := "", ""
		if ndi, ok := opts[i].def.(*Ext_org_kohsuke_args4j_NamedOptionDef); ok && ndi != nil {
			ni = getOptionDefData(ndi.Core_).name
		}
		if ndj, ok := opts[j].def.(*Ext_org_kohsuke_args4j_NamedOptionDef); ok && ndj != nil {
			nj = getOptionDefData(ndj.Core_).name
		}
		return ni < nj
	})
	maxHead := 0
	for _, ent := range opts {
		if head, _, ok := args4jOptionHead(ent); ok {
			if len(head) > maxHead {
				maxHead = len(head)
			}
		}
	}
	width := c.usageWidth
	if width <= 0 {
		width = 80
	}
	headCol := maxHead
	if headCol > width/2-1 {
		headCol = width/2 - 1
	}
	usageCol := width - headCol - 3
	for _, ent := range opts {
		if filter != nil && !filter.Select(ent.handler) {
			continue
		}
		head, d, ok := args4jOptionHead(ent)
		if !ok {
			continue
		}
		headLines := args4jWrapLines(head, headCol, true)
		usageLines := args4jWrapLines(d.usage+args4jDefaultSuffix(ent.handler), usageCol, false)
		nLines := len(headLines)
		if len(usageLines) > nLines {
			nLines = len(usageLines)
		}
		var lb strings.Builder
		for li := 0; li < nLines; li++ {
			hl := ""
			if li < len(headLines) {
				hl = headLines[li]
			}
			ul := ""
			if li < len(usageLines) {
				ul = usageLines[li]
			}
			lb.WriteString(hl)
			if li == 0 {
				if len(hl) < headCol {
					lb.WriteString(strings.Repeat(" ", headCol-len(hl)))
				}
				lb.WriteString(" : ")
				lb.WriteString(ul)
			} else if ul != "" {
				if len(hl) < headCol {
					lb.WriteString(strings.Repeat(" ", headCol-len(hl)))
				}
				lb.WriteString("   ")
				lb.WriteString(ul)
			}
			lb.WriteString("\n")
		}
		w.Write(lb.String())
	}
}

func (x *args4jParserCore) IsJavaNil_() bool      { return x == nil }
func (x *args4jHandlerEntry) IsJavaNil_() bool    { return x == nil }
func (x *args4jOptionDefData) IsJavaNil_() bool   { return x == nil }
func (x *args4jFieldSetterData) IsJavaNil_() bool { return x == nil }
func (x *args4jParametersImpl) IsJavaNil_() bool  { return x == nil }
func (x *args4jBuiltinHandler) IsJavaNil_() bool  { return x == nil }
