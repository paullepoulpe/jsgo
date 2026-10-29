// Port of org/kohsuke/args4j/CmdLineParser.java @ args4j 2.33
// Port of org/kohsuke/args4j/ClassParser.java @ args4j 2.33
// Port of org/kohsuke/args4j/Localizable.java @ args4j 2.33
// Port of org/kohsuke/args4j/Messages.java @ args4j 2.33
// Port of org/kohsuke/args4j/spi/Messages.java @ args4j 2.33
// Port of org/kohsuke/args4j/spi/Parameters.java @ args4j 2.33
// Port of org/kohsuke/args4j/Utilities.java @ args4j 2.33

/*
 * Copyright (c) 2013 Kohsuke Kawaguchi and other contributors
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to
 * deal in the Software without restriction, including without limitation the
 * rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
 * sell copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
 * IN THE SOFTWARE.
 */

// Package args4j is a port of the subset of args4j 2.33 used by Closure
// Compiler's command-line flags (CmdLineParser, @-file expansion,
// required/help/depends/forbids validation, usage formatting, and typed field
// setters and handlers).
//
// Instead of Java reflection and field annotations, an option struct implements
// Bean by returning its AnnotatedField descriptors in declaration order.
package args4j

import (
	"io"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Bean is implemented by structs whose fields carry Option or Argument metadata.
type Bean interface {
	Args4jFields() []AnnotatedField
}

// AnnotatedField binds a Bean Field to its Option and/or Argument metadata.
type AnnotatedField struct {
	Field    *Field
	Option   *Option
	Argument *Argument
}

// Parameters provides lookahead over remaining command-line tokens to an OptionHandler.
type Parameters interface {
	GetParameter(idx int) string
	TryGetParameter(idx int) (string, *CmdLineException)
	PeekParameter(idx int) (string, bool)
	Size() int
}

// CmdLineParser parses command-line flags and positional arguments into a Bean
// and formats usage/help output.
type CmdLineParser struct {
	options              []*OptionHandler
	optionByName         map[string]*OptionHandler
	arguments            []*OptionHandler
	parsingOptions       bool
	currentOptionHandler *OptionHandler
	parserProperties     *ParserProperties
}

// NewCmdLineParser creates a CmdLineParser for bean with default ParserProperties.
func NewCmdLineParser(bean any) *CmdLineParser {
	return NewCmdLineParserWithProperties(bean, ParserProperties_Defaults())
}

// NewCmdLineParserWithProperties creates a CmdLineParser for bean with custom properties.
func NewCmdLineParserWithProperties(bean any, parserProperties *ParserProperties) *CmdLineParser {
	p := &CmdLineParser{
		parsingOptions:   true,
		parserProperties: parserProperties,
	}
	if bean == nil {
		return p
	}
	if b, ok := bean.(Bean); ok {
		for _, af := range b.Args4jFields() {
			if af.Option != nil {
				p.AddOption(newSetterForField(af.Field), af.Option)
			}
			if af.Argument != nil {
				p.AddArgument(newSetterForField(af.Field), af.Argument)
			}
		}
	}
	if sorter := parserProperties.GetOptionSorter(); sorter != nil {
		slices.SortStableFunc(p.options, sorter)
		p.optionByName = nil
	}
	return p
}

func (p *CmdLineParser) GetProperties() *ParserProperties { return p.parserProperties }
func (p *CmdLineParser) GetArguments() []*OptionHandler   { return p.arguments }
func (p *CmdLineParser) GetOptions() []*OptionHandler     { return p.options }

func (p *CmdLineParser) AddArgument(setter Setter, a *Argument) {
	if setter == nil {
		panic(jutil.NewNullPointerExceptionMsg("Setter is null"))
	}
	if a == nil {
		panic(jutil.NewNullPointerExceptionMsg("Argument is null"))
	}

	h := p.CreateOptionHandler(NewOptionDef(a, setter.IsMultiValued()), setter)
	index := a.Index
	for index >= len(p.arguments) {
		p.arguments = append(p.arguments, nil)
	}
	if p.arguments[index] != nil {
		panic(NewIllegalAnnotationError(msgMultipleUseOfArgument.format(index)))
	}
	p.arguments[index] = h
}

func (p *CmdLineParser) AddOption(setter Setter, o *Option) {
	if setter == nil {
		panic(jutil.NewNullPointerExceptionMsg("Setter is null"))
	}
	if o == nil {
		panic(jutil.NewNullPointerExceptionMsg("Option is null"))
	}

	p.checkOptionNotInMap(o.Name)
	for _, alias := range o.Aliases {
		p.checkOptionNotInMap(alias)
	}
	h := p.CreateOptionHandler(NewNamedOptionDef(o), setter)
	p.options = append(p.options, h)
	if p.optionByName != nil {
		if _, exists := p.optionByName[o.Name]; !exists {
			p.optionByName[o.Name] = h
		}
		for _, alias := range o.Aliases {
			if _, exists := p.optionByName[alias]; !exists {
				p.optionByName[alias] = h
			}
		}
	}
}

func (p *CmdLineParser) checkOptionNotInMap(name string) {
	if p.findOptionByName(name) != nil {
		panic(NewIllegalAnnotationError(msgMultipleUseOfOption.format(name)))
	}
}

func (p *CmdLineParser) CreateOptionHandler(o *OptionDef, setter Setter) *OptionHandler {
	return createOptionHandler(p, o, setter)
}

func (p *CmdLineParser) PrintExample(filter OptionHandlerFilterFunc) string {
	return p.PrintExampleWithResourceBundle(filter, nil)
}

func (p *CmdLineParser) PrintExampleWithResourceBundle(mode OptionHandlerFilterFunc, rb ResourceBundle) string {
	if mode == nil {
		panic(jutil.NewNullPointerExceptionMsg("mode is null"))
	}
	var buf strings.Builder
	for _, h := range p.options {
		if len(h.Option.Usage()) == 0 || !mode(h) {
			continue
		}
		buf.WriteByte(' ')
		buf.WriteString(h.GetNameAndMetaWithProperties(rb, p.parserProperties))
	}
	return buf.String()
}

func (p *CmdLineParser) PrintUsage(out io.Writer) {
	p.PrintUsageToWriter(jutil.NewOutputStreamWriter(outputStream(out)), nil)
}

func (p *CmdLineParser) PrintUsageToWriter(out jutil.Writer, rb ResourceBundle) {
	p.PrintUsageToWriterWithFilter(out, rb, OptionHandlerFilter_PUBLIC)
}

func (p *CmdLineParser) PrintUsageToWriterWithFilter(out jutil.Writer, rb ResourceBundle, filter OptionHandlerFilterFunc) {
	if out == nil {
		panic(jutil.NewNullPointerException())
	}
	prefixLen := 0
	for _, h := range p.arguments {
		prefixLen = max(prefixLen, p.getPrefixLen(h, rb))
	}
	for _, h := range p.options {
		prefixLen = max(prefixLen, p.getPrefixLen(h, rb))
	}
	for _, h := range p.arguments {
		p.PrintOption(out, h, prefixLen, rb, filter)
	}
	for _, h := range p.options {
		p.PrintOption(out, h, prefixLen, rb, filter)
	}
	_ = out.TryFlush()
}

func (p *CmdLineParser) PrintOption(out jutil.Writer, handler *OptionHandler, maxMetaLen int, rb ResourceBundle, filter OptionHandlerFilterFunc) {
	option := handler.Option
	if len(option.Usage()) == 0 {
		return
	}
	if filter == nil {
		panic(jutil.NewNullPointerException())
	}
	if !filter(handler) {
		return
	}

	totalUsageWidth := p.parserProperties.GetUsageWidth()
	widthMetadata := min(maxMetaLen, (totalUsageWidth-4)/2)
	widthUsage := totalUsageWidth - 4 - widthMetadata

	defaultValuePart := p.createDefaultValuePart(handler)
	namesAndMetas := p.wrapLines(handler.GetNameAndMetaWithProperties(rb, p.parserProperties), widthMetadata)
	usages := p.wrapLines(localize(option.Usage(), rb)+defaultValuePart, widthUsage)

	var lineBuf strings.Builder
	for i := 0; i < max(len(namesAndMetas), len(usages)); i++ {
		nameAndMeta := ""
		if i < len(namesAndMetas) {
			nameAndMeta = namesAndMetas[i]
		}
		usage := ""
		if i < len(usages) {
			usage = usages[i]
		}
		lineBuf.Reset()
		lineBuf.WriteByte(' ')
		lineBuf.WriteString(nameAndMeta)
		for pad := widthMetadata - utf8.RuneCountInString(nameAndMeta); pad > 0; pad-- {
			lineBuf.WriteByte(' ')
		}
		if len(nameAndMeta) > 0 && i == 0 {
			lineBuf.WriteString(" : ")
		} else {
			lineBuf.WriteString("   ")
		}
		if len(usage) == 0 {
			lineBuf.WriteByte(' ')
		} else {
			lineBuf.WriteString(usage)
		}
		lineBuf.WriteByte('\n')
		_ = out.TryWrite(lineBuf.String())
	}
}

func (p *CmdLineParser) createDefaultValuePart(handler *OptionHandler) string {
	if _, isGetter := handler.Setter.(Getter); p.parserProperties.GetShowDefaults() && !handler.Option.Required() && isGetter {
		if v, ok := handler.PrintDefaultValue(); ok {
			return " " + msgDefaultValue.format(v)
		}
	}
	return ""
}

func localize(s string, rb ResourceBundle) string {
	if rb != nil {
		return rb.GetString(s)
	}
	return s
}

func (p *CmdLineParser) wrapLines(line string, maxLength int) []string {
	var rv []string
	for _, restOfLine := range jutil.StringSplitLiteral(line, "\n") {
		for utf8.RuneCountInString(restOfLine) > maxLength {
			byteAtMax := len(restOfLine)
			lastSpaceRune := -1
			lastSpaceByte := -1
			runeCount := 0
			for byteIdx, r := range restOfLine {
				if runeCount == maxLength {
					byteAtMax = byteIdx
					break
				}
				if r == ' ' {
					lastSpaceRune = runeCount
					lastSpaceByte = byteIdx
				}
				runeCount++
			}
			cutByte := byteAtMax
			if lastSpaceRune > maxLength*3/5 {
				cutByte = lastSpaceByte
			}
			rv = append(rv, restOfLine[:cutByte])
			restOfLine = jutil.StringTrim(restOfLine[cutByte:])
		}
		rv = append(rv, restOfLine)
	}
	return rv
}

func (p *CmdLineParser) getPrefixLen(h *OptionHandler, rb ResourceBundle) int {
	if len(h.Option.Usage()) == 0 {
		return 0
	}
	return utf8.RuneCountInString(h.GetNameAndMetaWithProperties(rb, p.parserProperties))
}

type cmdLineParams struct {
	owner *CmdLineParser
	args  []string
	pos   int
}

func (c *cmdLineParams) hasMore() bool        { return c.pos < len(c.args) }
func (c *cmdLineParams) currentToken() string { return c.args[c.pos] }
func (c *cmdLineParams) proceed(n int)        { c.pos += n }
func (c *cmdLineParams) Size() int            { return len(c.args) - c.pos }

func (c *cmdLineParams) GetParameter(idx int) string {
	val, err := c.TryGetParameter(idx)
	if err != nil {
		panic(err)
	}
	return val
}

func (c *cmdLineParams) PeekParameter(idx int) (string, bool) {
	if c.pos+idx >= len(c.args) || c.pos+idx < 0 {
		return "", false
	}
	return c.args[c.pos+idx], true
}

func (c *cmdLineParams) TryGetParameter(idx int) (string, *CmdLineException) {
	if c.pos+idx >= len(c.args) || c.pos+idx < 0 {
		optName := c.owner.currentOptionHandler.Option.String()
		return "", newCmdLineException(c.owner, msgMissingOperand, optName)
	}
	return c.args[c.pos+idx], nil
}

// splitToken replaces "-option=value" at c.pos with "value".
func (c *cmdLineParams) splitToken() {
	if c.pos >= 0 && c.pos < len(c.args) {
		if idx := strings.IndexByte(c.args[c.pos], '='); idx > 0 {
			c.args[c.pos] = c.args[c.pos][idx+1:]
		}
	}
}

func (p *CmdLineParser) ParseArgumentCollection(args []string) {
	p.ParseArgument(slices.Clone(args)...)
}

func (p *CmdLineParser) ParseArgument(args ...string) {
	if err := p.TryParseArgument(args...); err != nil {
		panic(err)
	}
}

// TryParseArgument parses command-line flags and positional arguments, returning
// a *CmdLineException instead of panicking on user argument errors.
func (p *CmdLineParser) TryParseArgument(args ...string) *CmdLineException {
	expandedArgs := args
	if p.parserProperties.GetAtSyntax() {
		var err *CmdLineException
		expandedArgs, err = p.tryExpandAtFiles(args)
		if err != nil {
			return err
		}
	}
	cmdLine := &cmdLineParams{owner: p, args: expandedArgs}

	presentSet := make(map[*OptionHandler]bool)
	var presentOrder []*OptionHandler
	argIndex := 0

	for cmdLine.hasMore() {
		arg := cmdLine.currentToken()
		if p.IsOption(arg) {
			delim := p.parserProperties.GetOptionValueDelimiter()
			isKeyValuePair := strings.Contains(arg, delim) || strings.IndexByte(arg, '=') != -1
			if isKeyValuePair {
				p.currentOptionHandler = p.findOptionHandler(arg)
			} else {
				p.currentOptionHandler = p.findOptionByName(arg)
			}
			if p.currentOptionHandler == nil {
				return newCmdLineException(p, msgUndefinedOption, arg)
			}
			if isKeyValuePair {
				cmdLine.splitToken()
			} else {
				cmdLine.proceed(1)
			}
		} else {
			if argIndex >= len(p.arguments) {
				msg := msgTooManyArguments
				if len(p.arguments) == 0 {
					msg = msgNoArgumentAllowed
				}
				return newCmdLineException(p, msg, arg)
			}
			p.currentOptionHandler = p.arguments[argIndex]
			if p.currentOptionHandler == nil {
				panic(jutil.NewIllegalStateExceptionMsg("@Argument with index=" + strconv.Itoa(argIndex) + " is undefined"))
			}
			if !p.currentOptionHandler.Option.IsMultiValued() {
				argIndex++
			}
		}
		diff, err := p.currentOptionHandler.tryParseArguments(cmdLine)
		if err != nil {
			return err
		}
		cmdLine.proceed(diff)
		if !presentSet[p.currentOptionHandler] {
			presentSet[p.currentOptionHandler] = true
			presentOrder = append(presentOrder, p.currentOptionHandler)
		}
	}

	for _, handler := range p.options {
		if handler.Option.Help() && presentSet[handler] {
			return nil
		}
	}
	return p.tryCheckRequiredOptionsAndArguments(presentSet, presentOrder)
}

func (p *CmdLineParser) tryExpandAtFiles(args []string) ([]string, *CmdLineException) {
	var result []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "@") {
			file := jutil.NewFile(arg[1:])
			if !file.Exists() {
				return nil, newCmdLineException(p, msgNoSuchFile, file.GetPath())
			}
			lines, ex := jutil.Files_AsCharSource(file, jutil.StandardCharsets_UTF_8).TryReadLines()
			if ex != nil {
				return nil, NewCmdLineExceptionWithParserMessageAndCause(p, "Failed to parse "+file.String(), ex)
			}
			result = append(result, lines...)
		} else {
			result = append(result, arg)
		}
	}
	return result, nil
}

func (p *CmdLineParser) tryCheckRequiredOptionsAndArguments(presentSet map[*OptionHandler]bool, presentOrder []*OptionHandler) *CmdLineException {
	for _, handler := range p.options {
		if handler.Option.Required() && !presentSet[handler] {
			return newCmdLineException(p, msgRequiredOptionMissing, handler.Option.String())
		}
	}
	for _, handler := range p.arguments {
		if handler.Option.Required() && !presentSet[handler] {
			return newCmdLineException(p, msgRequiredArgumentMissing, handler.Option.String())
		}
	}
	for _, handler := range presentOrder {
		if !handler.Option.IsNamed() {
			continue
		}
		for _, depend := range handler.Option.Depends() {
			if !presentSet[p.findOptionHandler(depend)] {
				return newCmdLineException(p, msgRequiresOptionMissing, handler.Option.String(), "["+strings.Join(handler.Option.Depends(), ", ")+"]")
			}
		}
	}
	for _, handler := range presentOrder {
		if !handler.Option.IsNamed() {
			continue
		}
		for _, forbid := range handler.Option.Forbids() {
			if presentSet[p.findOptionHandler(forbid)] {
				return newCmdLineException(p, msgForbiddenOptionPresent, handler.Option.String(), "["+strings.Join(handler.Option.Forbids(), ", ")+"]")
			}
		}
	}
	return nil
}

func (p *CmdLineParser) findOptionHandler(name string) *OptionHandler {
	delim := p.parserProperties.GetOptionValueDelimiter()
	pos := strings.Index(name, delim)
	if pos < 0 {
		pos = strings.IndexByte(name, '=')
	}
	if pos > 0 {
		name = name[:pos]
	}
	return p.findOptionByName(name)
}

func (p *CmdLineParser) findOptionByName(name string) *OptionHandler {
	if len(p.options) == 0 {
		return nil
	}
	if p.optionByName == nil {
		m := make(map[string]*OptionHandler, len(p.options)*2)
		for _, h := range p.options {
			optName := h.Option.Name()
			if _, exists := m[optName]; !exists {
				m[optName] = h
			}
			for _, alias := range h.Option.Aliases() {
				if _, exists := m[alias]; !exists {
					m[alias] = h
				}
			}
		}
		p.optionByName = m
	}
	return p.optionByName[name]
}

func (p *CmdLineParser) IsOption(arg string) bool {
	return p.parsingOptions && strings.HasPrefix(arg, "-")
}

func (p *CmdLineParser) SetUsageWidth(usageWidth int) {
	p.parserProperties.WithUsageWidth(usageWidth)
}

func (p *CmdLineParser) StopOptionParsing() {
	p.parsingOptions = false
}

func (p *CmdLineParser) PrintSingleLineUsage(out io.Writer) {
	if out == nil {
		panic(jutil.NewNullPointerExceptionMsg("OutputStream is null"))
	}
	p.PrintSingleLineUsageToWriter(jutil.NewOutputStreamWriter(outputStream(out)), nil)
}

func (p *CmdLineParser) PrintSingleLineUsageToWriter(w jutil.Writer, rb ResourceBundle) {
	if w == nil {
		panic(jutil.NewNullPointerExceptionMsg("Writer is null"))
	}
	var buf strings.Builder
	for _, h := range p.arguments {
		p.printSingleLineOption(&buf, h, rb)
	}
	for _, h := range p.options {
		p.printSingleLineOption(&buf, h, rb)
	}
	if buf.Len() > 0 {
		_ = w.TryWrite(buf.String())
	}
	_ = w.TryFlush()
}

func (p *CmdLineParser) printSingleLineOption(buf *strings.Builder, h *OptionHandler, rb ResourceBundle) {
	option := h.Option
	buf.WriteByte(' ')
	if !option.Required() {
		buf.WriteByte('[')
	}
	buf.WriteString(h.GetNameAndMetaWithProperties(rb, p.parserProperties))
	if option.IsMultiValued() {
		buf.WriteString(" ...")
	}
	if !option.Required() {
		buf.WriteByte(']')
	}
}

func outputStream(out io.Writer) io.WriteCloser {
	if wc, ok := out.(io.WriteCloser); ok {
		return wc
	}
	return nopCloser{out}
}

type nopCloser struct{ io.Writer }

func (nopCloser) Close() error { return nil }

// ResourceBundle localizes option usage strings and meta variables when passed
// to CmdLineParser usage printers.
type ResourceBundle interface {
	GetString(key string) string
}

// messageTemplate is an English args4j MessageFormat pattern from
// org/kohsuke/args4j/Messages.properties or spi/Messages.properties.
type messageTemplate string

func (m messageTemplate) format(args ...any) string {
	s := string(m)
	var sb strings.Builder
	sb.Grow(len(s) + 16)
	quoted := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\'' {
			if i+1 < len(s) && s[i+1] == '\'' {
				sb.WriteByte('\'')
				i++
			} else {
				quoted = !quoted
			}
			continue
		}
		if c == '{' && !quoted {
			if end := strings.IndexByte(s[i+1:], '}'); end >= 0 {
				if idx, err := strconv.Atoi(s[i+1 : i+1+end]); err == nil && idx >= 0 {
					if idx < len(args) {
						sb.WriteString(jutil.ToJavaString(args[idx]))
					} else {
						sb.WriteByte('{')
						sb.WriteString(strconv.Itoa(idx))
						sb.WriteByte('}')
					}
					i += 1 + end
					continue
				}
			}
		}
		sb.WriteByte(c)
	}
	return sb.String()
}

const (
	msgMissingOperand          messageTemplate = `Option "{0}" takes an operand`
	msgUndefinedOption         messageTemplate = `"{0}" is not a valid option`
	msgNoArgumentAllowed       messageTemplate = `No argument is allowed: {0}`
	msgRequiredOptionMissing   messageTemplate = `Option "{0}" is required`
	msgRequiredArgumentMissing messageTemplate = `Argument "{0}" is required`
	msgTooManyArguments        messageTemplate = `Too many arguments: {0}`
	msgMultipleUseOfArgument   messageTemplate = `Argument with index {0} is used more than once`
	msgMultipleUseOfOption     messageTemplate = `Option name {0} is used more than once`
	msgUnknownHandler          messageTemplate = `No OptionHandler is registered to handle {0}`
	msgRequiresOptionMissing   messageTemplate = `option "{0}" requires the option(s) {1}`
	msgForbiddenOptionPresent  messageTemplate = `option "{0}" cannot be used with the option(s) {1}`
	msgNoSuchFile              messageTemplate = `No such file: {0}`
	msgDefaultValue            messageTemplate = `(default: {0})`

	spiMsgIllegalOperand        messageTemplate = `"{1}" is not a valid value for "{0}"`
	spiMsgIllegalBoolean        messageTemplate = `"{0}" is not a legal boolean value`
	spiMsgIllegalFieldSignature messageTemplate = `Field of type {0} isn't supported by args4j`
	spiMsgIllegalList           messageTemplate = `type of {0} is not a List`
	spiMsgFormatErrorForMap     messageTemplate = `An argument for setting a Map must contain a "="`
	spiMsgMapHasNoKey           messageTemplate = `A key must be set.`

	defaultMetaStringArray = "STRING[]"
	defaultMetaString      = "VAL"
)
