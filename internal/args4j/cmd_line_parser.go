// Port of org/kohsuke/args4j/CmdLineParser.java @ args4j 2.33

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

package args4j

import (
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// CmdLineParser is the command line argument owner.
//
// For typical usage, see this example:
// https://args4j.dev.java.net/source/browse/args4j/args4j/examples/SampleMain.java?view=markup
//
// Java: class CmdLineParser
type CmdLineParser struct {
	// Discovered OptionHandlers for options.
	options []OptionHandler

	// Discovered OptionHandlers for arguments.
	arguments []OptionHandler

	parsingOptions       bool
	currentOptionHandler OptionHandler

	// settings for the parser
	parserProperties *ParserProperties
}

// NewCmdLineParser creates a new command line owner that parses
// arguments/options and set them into the given object.
//
// bean is an instance of a class annotated by Option and Argument (a Bean).
// This object will receive values. If this is nil, the processing will be
// skipped, which is useful if you'd like to feed metadata from other
// sources.
//
// It throws IllegalAnnotationError if the option bean class is using args4j
// annotations incorrectly.
//
// Java: CmdLineParser(Object)
func NewCmdLineParser(bean any) *CmdLineParser {
	// for display purposes, we like the arguments in argument order, but the
	// options in alphabetical order
	return NewCmdLineParserWithProperties(bean, ParserProperties_Defaults())
}

// NewCmdLineParserWithProperties creates a new command line owner that
// parses arguments/options and set them into the given object, with
// parserProperties as its settings.
//
// Java: CmdLineParser(Object, ParserProperties)
func NewCmdLineParserWithProperties(bean any, parserProperties *ParserProperties) *CmdLineParser {
	p := &CmdLineParser{parsingOptions: true}
	p.parserProperties = parserProperties
	// A 'return' in the constructor just skips the rest of the implementation
	// and returns the new object directly.
	if jutil.IsNil(bean) {
		return p
	}

	// Parse the metadata and create the setters
	NewClassParser().Parse(bean, p)

	if parserProperties.GetOptionSorter() != nil {
		// Collections.sort is a stable sort.
		slices.SortStableFunc(p.options, parserProperties.GetOptionSorter())
	}
	return p
}

// Java: CmdLineParser.getProperties()
func (p *CmdLineParser) GetProperties() *ParserProperties {
	return p.parserProperties
}

// AddArgument programmatically defines an argument (instead of reading it
// from annotations as normal). It throws NullPointerException if setter or a
// is nil.
//
// Java: CmdLineParser.addArgument(Setter, Argument)
func (p *CmdLineParser) AddArgument(setter Setter, a *Argument) {
	Utilities_CheckNonNull(setter, "Setter")
	Utilities_CheckNonNull(a, "Argument")

	h := p.CreateOptionHandler(NewOptionDef(a, setter.IsMultiValued()), setter)
	index := a.Index
	// make sure the argument will fit in the list
	for index >= len(p.arguments) {
		p.arguments = append(p.arguments, nil)
	}
	if p.arguments[index] != nil {
		panic(NewIllegalAnnotationError(Messages_MULTIPLE_USE_OF_ARGUMENT.Format(index)))
	}
	p.arguments[index] = h
}

// AddOption programmatically defines an option (instead of reading it from
// annotations as normal). It throws NullPointerException if setter or o is
// nil, and IllegalAnnotationError if the option name or one of the aliases
// is already taken.
//
// Java: CmdLineParser.addOption(Setter, Option)
func (p *CmdLineParser) AddOption(setter Setter, o *Option) {
	Utilities_CheckNonNull(setter, "Setter")
	Utilities_CheckNonNull(o, "Option")

	p.checkOptionNotInMap(o.Name)
	for _, alias := range o.Aliases {
		p.checkOptionNotInMap(alias)
	}
	p.options = append(p.options, p.CreateOptionHandler(NewNamedOptionDef(o), setter))
}

// GetArguments lists up all the defined arguments in the order.
//
// Java: CmdLineParser.getArguments()
func (p *CmdLineParser) GetArguments() []OptionHandler {
	return p.arguments
}

// GetOptions lists up all the defined options.
//
// Java: CmdLineParser.getOptions()
func (p *CmdLineParser) GetOptions() []OptionHandler {
	return p.options
}

// Java: CmdLineParser.checkOptionNotInMap(String)
func (p *CmdLineParser) checkOptionNotInMap(name string) {
	Utilities_CheckNonNull(name, "name")

	if p.findOptionByName(name) != nil {
		panic(NewIllegalAnnotationError(Messages_MULTIPLE_USE_OF_OPTION.Format(name)))
	}
}

// CreateOptionHandler creates an OptionHandler that handles the given Option
// annotation and the Setter instance.
//
// Java: CmdLineParser.createOptionHandler(OptionDef, Setter)
func (p *CmdLineParser) CreateOptionHandler(o OptionDef, setter Setter) OptionHandler {
	Utilities_CheckNonNull(o, "OptionDef")
	Utilities_CheckNonNull(setter, "Setter")
	return OptionHandlerRegistry_GetRegistry().CreateOptionHandler(p, o, setter)
}

// PrintExample formats a command line example into a string. See
// PrintExampleWithResourceBundle for more details. filter must not be nil;
// the result is never nil.
//
// Java: CmdLineParser.printExample(OptionHandlerFilter)
func (p *CmdLineParser) PrintExample(filter OptionHandlerFilter) string {
	return p.PrintExampleWithResourceBundle(filter, nil)
}

// Deprecated: Use PrintExample.
//
// Java: CmdLineParser.printExample(ExampleMode)
func (p *CmdLineParser) PrintExampleWithExampleMode(mode ExampleMode) string {
	return p.PrintExampleWithExampleModeAndResourceBundle(mode, nil)
}

// PrintExampleWithResourceBundle formats a command line example into a
// string.
//
// This method produces a string like " -d <dir> -v -b". This is useful for
// printing a command line example (perhaps as a part of the usage screen).
//
// mode determines which options will be a part of the returned string. Must
// not be nil. If rb is non-nil, meta variables (<dir> in the above example)
// is treated as a key to this resource bundle, and the associated value is
// printed. See Option.MetaVar. This is to support localization. Passing nil
// would print Option.MetaVar directly.
//
// The result is never nil. If there's no option, this method returns just
// the empty string "". Otherwise, this method returns a string that contains
// a space at the beginning (but not at the end). This allows you to do
// something like:
// System.err.println("java -jar my.jar"+parser.printExample(REQUIRED)+" arg1 arg2");
// It throws NullPointerException if mode is nil.
//
// Java: CmdLineParser.printExample(OptionHandlerFilter, ResourceBundle)
func (p *CmdLineParser) PrintExampleWithResourceBundle(mode OptionHandlerFilter, rb ResourceBundle) string {
	buf := jutil.NewStringBuilder()

	Utilities_CheckNonNull(mode, "mode")

	for _, h := range p.options {
		option := h.GetOptionHandlerBase().Option
		if jutil.StringLength(option.Usage()) == 0 {
			continue // ignore
		}
		if !mode.Select(h) {
			continue
		}

		buf.AppendChar(' ')
		buf.Append(h.GetNameAndMetaWithProperties(rb, p.parserProperties))
	}

	return buf.String()
}

// Deprecated: Use PrintExampleWithResourceBundle.
//
// Java: CmdLineParser.printExample(ExampleMode, ResourceBundle)
func (p *CmdLineParser) PrintExampleWithExampleModeAndResourceBundle(mode ExampleMode, rb ResourceBundle) string {
	return p.PrintExampleWithResourceBundle(mode, rb)
}

// PrintUsage prints the list of options and their usages to the screen.
//
// This is a convenience method for calling
// printUsage(new OutputStreamWriter(out),null) so that you can do
// printUsage(System.err).
//
// Java: CmdLineParser.printUsage(OutputStream)
func (p *CmdLineParser) PrintUsage(out io.Writer) {
	p.PrintUsageToWriter(jutil.NewOutputStreamWriter(cmdLineParser_outputStream(out)), nil)
}

// PrintUsageToWriter prints the list of all the non-hidden options and their
// usages to the screen.
//
// Short for printUsage(out,rb,OptionHandlerFilter.PUBLIC).
//
// Java: CmdLineParser.printUsage(Writer, ResourceBundle)
func (p *CmdLineParser) PrintUsageToWriter(out jutil.Writer, rb ResourceBundle) {
	p.PrintUsageToWriterWithFilter(out, rb, OptionHandlerFilter_PUBLIC)
}

// PrintUsageToWriterWithFilter prints the list of all the non-hidden options
// and their usages to the screen.
//
// If rb is non-nil, Option.Usage is treated as a key to obtain the actual
// message from this resource bundle. filter controls which options to be
// printed.
//
// Java: CmdLineParser.printUsage(Writer, ResourceBundle, OptionHandlerFilter)
func (p *CmdLineParser) PrintUsageToWriterWithFilter(out jutil.Writer, rb ResourceBundle, filter OptionHandlerFilter) {
	w := jutil.NewPrintWriter(out)
	// determine the length of the option + metavar first
	len_ := 0
	for _, h := range p.arguments {
		curLen := p.getPrefixLen(h, rb)
		len_ = max(len_, curLen)
	}
	for _, h := range p.options {
		curLen := p.getPrefixLen(h, rb)
		len_ = max(len_, curLen)
	}

	// then print
	for _, h := range p.arguments {
		p.PrintOption(w, h, len_, rb, filter)
	}
	for _, h := range p.options {
		p.PrintOption(w, h, len_, rb, filter)
	}

	w.Flush()
}

// PrintOption prints usage information for a given option.
//
// Subtypes may override this method and determine which options get printed
// (or other things), based on OptionHandler (perhaps by using
// handler.setter.asAnnotatedElement()).
//
// out is the Writer to write into, handler the handler where to receive the
// information, len_ the maximum length of metadata column, rb the
// ResourceBundle for I18N.
//
// Java: CmdLineParser.printOption(PrintWriter, OptionHandler, int, ResourceBundle, OptionHandlerFilter)
func (p *CmdLineParser) PrintOption(out *jutil.PrintWriter, handler OptionHandler, len_ int, rb ResourceBundle, filter OptionHandlerFilter) {
	option := handler.GetOptionHandlerBase().Option
	// Hiding options without usage information
	//
	// Go-only: a Go string is never null, so Java's usage() == null test is
	// always false.
	if jutil.StringLength(option.Usage()) == 0 ||
		!filter.Select(handler) {
		return
	}

	// What is the width of the two data columns
	totalUsageWidth := p.parserProperties.GetUsageWidth()
	widthMetadata := min(len_, (totalUsageWidth-4)/2)
	widthUsage := totalUsageWidth - 4 - widthMetadata

	defaultValuePart := p.createDefaultValuePart(handler)

	// Line wrapping
	// the 'left' side
	namesAndMetas := p.wrapLines(handler.GetNameAndMetaWithProperties(rb, p.parserProperties), widthMetadata)
	// the 'right' side
	usages := p.wrapLines(p.localize(option.Usage(), rb)+defaultValuePart, widthUsage)

	// Output
	for i := 0; i < max(len(namesAndMetas), len(usages)); i++ {
		nameAndMeta := ""
		if i < len(namesAndMetas) {
			nameAndMeta = namesAndMetas[i]
		}
		usage := ""
		if i < len(usages) {
			usage = usages[i]
		}
		var format string
		if jutil.StringLength(nameAndMeta) > 0 && i == 0 {
			format = " %1$-" + strconv.Itoa(widthMetadata) + "s : %2$-1s"
		} else {
			format = " %1$-" + strconv.Itoa(widthMetadata) + "s   %2$-1s"
		}
		output := jutil.StringFormat(format, nameAndMeta, usage)

		out.Println(output)
	}
}

// Java: CmdLineParser.createDefaultValuePart(OptionHandler)
func (p *CmdLineParser) createDefaultValuePart(handler OptionHandler) string {
	b := handler.GetOptionHandlerBase()
	if _, isGetter := b.Setter.(Getter); p.parserProperties.GetShowDefaults() && !b.Option.Required() && isGetter {
		v := handler.PrintDefaultValue()
		if v != nil {
			return " " + Messages_DEFAULT_VALUE.Format(*v)
		}
	}
	return ""
}

// Java: CmdLineParser.localize(String, ResourceBundle)
func (p *CmdLineParser) localize(s string, rb ResourceBundle) string {
	if rb != nil {
		return rb.GetString(s)
	}
	return s
}

// wrapLines wraps a line so that the resulting parts are not longer than a
// given maximum length: line is the line to wrap, maxLength the maximum
// length for the resulting parts. It returns the list of all wrapped parts.
//
// Java: CmdLineParser.wrapLines(String, int)
func (p *CmdLineParser) wrapLines(line string, maxLength int) []string {
	var rv []string
	// String.split("\\n"): the regex matches only the literal newline.
	for _, restOfLine := range jutil.StringSplitLiteral(line, "\n") {
		for jutil.StringLength(restOfLine) > maxLength {
			// try to wrap at space, but don't try too hard as some languages
			// don't even have whitespaces.
			var lineLength int
			candidate := jutil.StringSubstring(restOfLine, 0, maxLength)
			sp := jutil.StringLastIndexOf(candidate, " ")
			if sp > maxLength*3/5 {
				lineLength = sp
			} else {
				lineLength = maxLength
			}
			rv = append(rv, jutil.StringSubstring(restOfLine, 0, lineLength))
			restOfLine = jutil.StringTrim(jutil.StringSubstringFrom(restOfLine, lineLength))
		}
		rv = append(rv, restOfLine)
	}
	return rv
}

// Java: CmdLineParser.getPrefixLen(OptionHandler, ResourceBundle)
func (p *CmdLineParser) getPrefixLen(h OptionHandler, rb ResourceBundle) int {
	if jutil.StringLength(h.GetOptionHandlerBase().Option.Usage()) == 0 {
		return 0
	}

	return jutil.StringLength(h.GetNameAndMetaWithProperties(rb, p.parserProperties))
}

// CmdLineParser_CmdLineImpl is essentially a pointer over a String array. It
// can move forward; can look ahead.
//
// Java: class CmdLineParser.CmdLineImpl implements Parameters
type CmdLineParser_CmdLineImpl struct {
	this *CmdLineParser // CmdLineParser.this
	args []string
	pos  int
}

// Java: CmdLineImpl(String[])
func newCmdLineParser_CmdLineImpl(this *CmdLineParser, args []string) *CmdLineParser_CmdLineImpl {
	c := &CmdLineParser_CmdLineImpl{this: this}
	c.args = args
	c.pos = 0
	return c
}

// Java: CmdLineImpl.hasMore()
func (c *CmdLineParser_CmdLineImpl) HasMore() bool {
	return c.pos < len(c.args)
}

// Java: CmdLineImpl.getCurrentToken()
func (c *CmdLineParser_CmdLineImpl) GetCurrentToken() string {
	return c.args[c.pos]
}

// Java: CmdLineImpl.proceed(int)
func (c *CmdLineParser_CmdLineImpl) proceed(n int) {
	c.pos += n
}

// Java: CmdLineImpl.getParameter(int)
func (c *CmdLineParser_CmdLineImpl) GetParameter(idx int) string {
	if c.pos+idx >= len(c.args) || c.pos+idx < 0 {
		panic(NewCmdLineExceptionWithLocalizable(c.this, Messages_MISSING_OPERAND, c.this.getOptionName()))
	}
	return c.args[c.pos+idx]
}

// Java: CmdLineImpl.size()
func (c *CmdLineParser_CmdLineImpl) Size() int {
	return len(c.args) - c.pos
}

// splitToken is used when the current token is of the form "-option=value",
// to replace the current token by "value", as if this was given as two
// tokens "-option value".
//
// Java: CmdLineImpl.splitToken()
func (c *CmdLineParser_CmdLineImpl) splitToken() {
	if c.pos < len(c.args) && c.pos >= 0 {
		idx := jutil.StringIndexOf(c.args[c.pos], "=")
		if idx > 0 {
			c.args[c.pos] = jutil.StringSubstringFrom(c.args[c.pos], idx+1)
		}
	}
}

// Java: CmdLineParser.getOptionName()
func (p *CmdLineParser) getOptionName() string {
	return p.currentOptionHandler.GetOptionHandlerBase().Option.String()
}

// ParseArgumentCollection is the same as ParseArgument.
//
// Java: CmdLineParser.parseArgument(Collection)
func (p *CmdLineParser) ParseArgumentCollection(args []string) {
	// Collection.toArray(new String[size]) copies the elements.
	p.ParseArgument(append(make([]string, 0, len(args)), args...)...)
}

// ParseArgument parses the command line arguments and set them to the option
// bean given in the constructor.
//
// It throws CmdLineException if there's any error parsing arguments, or if a
// required option was not given.
//
// Java: CmdLineParser.parseArgument(String...)
func (p *CmdLineParser) ParseArgument(args ...string) {
	// Go-only: Java first calls checkNonNull(args, "args"); a Go slice is
	// never null (a nil slice is an empty array).

	expandedArgs := args
	if p.parserProperties.GetAtSyntax() {
		expandedArgs = p.expandAtFiles(args)
	}
	cmdLine := newCmdLineParser_CmdLineImpl(p, expandedArgs)

	// Go-only: language-adapter — Java HashSet order: the handlers hash by
	// identity, so Java's iteration order below is not deterministic; the
	// port iterates in insertion order.
	present := jutil.NewLinkedHashSet[OptionHandler]()
	argIndex := 0

	for cmdLine.HasMore() {
		arg := cmdLine.GetCurrentToken()
		if p.IsOption(arg) {
			// '=' is for historical compatibility fallback
			isKeyValuePair := strings.Contains(arg, p.parserProperties.GetOptionValueDelimiter()) || jutil.StringIndexOf(arg, "=") != -1

			// parse this as an option.
			if isKeyValuePair {
				p.currentOptionHandler = p.findOptionHandler(arg)
			} else {
				p.currentOptionHandler = p.findOptionByName(arg)
			}

			if p.currentOptionHandler == nil {
				// TODO: insert dynamic handler processing
				panic(NewCmdLineExceptionWithLocalizable(p, Messages_UNDEFINED_OPTION, arg))
			}

			// known option; skip its name
			if isKeyValuePair {
				cmdLine.splitToken()
			} else {
				cmdLine.proceed(1)
			}
		} else {
			if argIndex >= len(p.arguments) {
				msg := Messages_TOO_MANY_ARGUMENTS
				if len(p.arguments) == 0 {
					msg = Messages_NO_ARGUMENT_ALLOWED
				}
				panic(NewCmdLineExceptionWithLocalizable(p, msg, arg))
			}

			// known argument
			p.currentOptionHandler = p.arguments[argIndex]
			if p.currentOptionHandler == nil { // this is a programmer error. arg index should be continuous
				panic(jutil.NewIllegalStateExceptionMsg("@Argument with index=" + strconv.Itoa(argIndex) + " is undefined"))
			}

			if !p.currentOptionHandler.GetOptionHandlerBase().Option.IsMultiValued() {
				argIndex++
			}
		}
		diff := p.currentOptionHandler.ParseArguments(cmdLine)
		cmdLine.proceed(diff)
		present.Add(p.currentOptionHandler)
	}

	// check whether a help option is set
	helpSet := false
	for _, handler := range p.options {
		if handler.GetOptionHandlerBase().Option.Help() && present.Contains(handler) {
			helpSet = true
		}
	}

	if !helpSet {
		p.checkRequiredOptionsAndArguments(present)
	}
}

// expandAtFiles expands every entry prefixed with the AT sign by reading the
// file. The AT sign is used to reference another file that contains command
// line options separated by line breaks. args are the command line
// arguments to be preprocessed. It returns args with the @ sequences
// replaced by the text files referenced by the @ sequences, split around
// the line breaks. It throws CmdLineException.
//
// Java: CmdLineParser.expandAtFiles(String[])
func (p *CmdLineParser) expandAtFiles(args []string) []string {
	result := []string{}
	for _, arg := range args {
		if jutil.StringStartsWith(arg, "@") {
			file := jutil.NewFile(jutil.StringSubstringFrom(arg, 1))
			if !file.Exists() {
				panic(NewCmdLineExceptionWithLocalizable(p, Messages_NO_SUCH_FILE, file.GetPath()))
			}
			jutil.Try(func() {
				result = append(result, cmdLineParser_readAllLines(file)...)
			}, func(ex jutil.IIOException) {
				panic(NewCmdLineExceptionWithParserMessageAndCause(p, "Failed to parse "+file.String(), ex))
			})
		} else {
			result = append(result, arg)
		}
	}
	return result
}

// cmdLineParser_readAllLines reads all lines of a file with the platform
// encoding.
//
// Go-only: language-adapter — new BufferedReader(new FileReader(f)): the
// file is opened as FileInputStream opens it (FileNotFoundException), decoded
// with the default charset, UTF-8, and split into lines at "\n", "\r" and
// "\r\n", as BufferedReader.readLine splits them.
//
// Java: CmdLineParser.readAllLines(File)
func cmdLineParser_readAllLines(f *jutil.File) []string {
	r := jutil.Files_AsCharSource(f, jutil.StandardCharsets_UTF_8).OpenStream()
	defer r.Close()
	result := []string{}
	result = append(result, jutil.CharStreams_ReadLines(r)...)
	return result
}

// Java: CmdLineParser.checkRequiredOptionsAndArguments(Set)
func (p *CmdLineParser) checkRequiredOptionsAndArguments(present *jutil.LinkedHashSet[OptionHandler]) {
	// make sure that all mandatory options are present
	for _, handler := range p.options {
		option := handler.GetOptionHandlerBase().Option
		if option.Required() && !present.Contains(handler) {
			panic(NewCmdLineExceptionWithLocalizable(p, Messages_REQUIRED_OPTION_MISSING, option.String()))
		}
	}

	// make sure that all mandatory arguments are present
	for _, handler := range p.arguments {
		option := handler.GetOptionHandlerBase().Option
		if option.Required() && !present.Contains(handler) {
			panic(NewCmdLineExceptionWithLocalizable(p, Messages_REQUIRED_ARGUMENT_MISSING, option.String()))
		}
	}

	//make sure that all requires arguments are present
	for handler := range present.All() {
		option := handler.GetOptionHandlerBase().Option
		if named, ok := option.(*NamedOptionDef); ok && !p.isHandlerHasHisOptions(named, present) {
			panic(NewCmdLineExceptionWithLocalizable(p, Messages_REQUIRES_OPTION_MISSING,
				option.String(), cmdLineParser_arraysToString(named.Depends())))
		}
	}

	//make sure that all forbids arguments are not present
	for handler := range present.All() {
		option := handler.GetOptionHandlerBase().Option
		if named, ok := option.(*NamedOptionDef); ok && !p.isHandlerAllowOtherOptions(named, present) {
			panic(NewCmdLineExceptionWithLocalizable(p, Messages_FORBIDDEN_OPTION_PRESENT,
				option.String(), cmdLineParser_arraysToString(named.Forbids())))
		}
	}
}

// cmdLineParser_arraysToString is Arrays.toString(Object[]) of a String
// array.
func cmdLineParser_arraysToString(a []string) string {
	return "[" + strings.Join(a, ", ") + "]"
}

// isHandlerHasHisOptions returns true if all options required by option are
// present, false otherwise.
//
// Java: CmdLineParser.isHandlerHasHisOptions(NamedOptionDef, Set)
func (p *CmdLineParser) isHandlerHasHisOptions(option *NamedOptionDef, present *jutil.LinkedHashSet[OptionHandler]) bool {
	for _, depend := range option.Depends() {
		if !present.Contains(p.findOptionHandler(depend)) {
			return false
		}
	}
	return true
}

// isHandlerAllowOtherOptions returns true if all options forbid by option
// are not present, false otherwise.
//
// Java: CmdLineParser.isHandlerAllowOtherOptions(NamedOptionDef, Set)
func (p *CmdLineParser) isHandlerAllowOtherOptions(option *NamedOptionDef, present *jutil.LinkedHashSet[OptionHandler]) bool {
	for _, forbid := range option.Forbids() {
		if present.Contains(p.findOptionHandler(forbid)) {
			return false
		}
	}
	return true
}

// Java: CmdLineParser.findOptionHandler(String)
func (p *CmdLineParser) findOptionHandler(name string) OptionHandler {
	// Look for key/value pair first.
	pos := jutil.StringIndexOf(name, p.parserProperties.GetOptionValueDelimiter())
	if pos < 0 {
		pos = jutil.StringIndexOf(name, "=") // historical compatibility fallback
	}
	if pos > 0 {
		name = jutil.StringSubstring(name, 0, pos)
	}
	return p.findOptionByName(name)
}

// findOptionByName finds a registered OptionHandler by its name or its
// alias. It returns the OptionHandler or nil.
//
// Java: CmdLineParser.findOptionByName(String)
func (p *CmdLineParser) findOptionByName(name string) OptionHandler {
	for _, h := range p.options {
		option := h.GetOptionHandlerBase().Option.(*NamedOptionDef)
		if name == option.Name() {
			return h
		}
		for _, alias := range option.Aliases() {
			if name == alias {
				return h
			}
		}
	}
	return nil
}

// IsOption returns true if the given token is an option (as opposed to an
// argument). It throws NullPointerException if arg is null.
//
// Java: CmdLineParser.isOption(String)
func (p *CmdLineParser) IsOption(arg string) bool {
	Utilities_CheckNonNull(arg, "arg")

	return p.parsingOptions && jutil.StringStartsWith(arg, "-")
}

// CmdLineParser_RegisterHandler registers a user-defined OptionHandler class
// with args4j.
//
// This method allows users to extend the behavior of args4j by writing their
// own OptionHandler implementation.
//
// valueType: the specified handler is used when the field/method annotated
// by Option is of this type. handlerClass: this class must have the
// constructor that has the same signature as
// OptionHandler(CmdLineParser, OptionDef, Setter).
//
// Deprecated: You should use OptionHandlerRegistry.registerHandler instead.
//
// Java: CmdLineParser.registerHandler(Class, Class)
func CmdLineParser_RegisterHandler(valueType *Class, handlerClass OptionHandlerClass) {
	Utilities_CheckNonNull(valueType, "valueType")
	Utilities_CheckNonNull(handlerClass, "handlerClass")

	OptionHandlerRegistry_GetRegistry().RegisterHandler(valueType, handlerClass)
}

// SetUsageWidth sets the width of the usage output in columns. It throws
// IllegalArgumentException if usageWidth is negative.
//
// Deprecated: Use ParserProperties.withUsageWidth(int) instead.
//
// Java: CmdLineParser.setUsageWidth(int)
func (p *CmdLineParser) SetUsageWidth(usageWidth int) {
	p.parserProperties.WithUsageWidth(usageWidth)
}

// StopOptionParsing signals the parser that parsing the options has
// finished.
//
// Everything seen after this call is treated as an argument as opposed to an
// option.
//
// Java: CmdLineParser.stopOptionParsing()
func (p *CmdLineParser) StopOptionParsing() {
	p.parsingOptions = false
}

// PrintSingleLineUsage prints a single-line usage to the screen.
//
// This is a convenience method for calling
// printUsage(new OutputStreamWriter(out),null) so that you can do
// printUsage(System.err). It throws NullPointerException if out is nil.
//
// Java: CmdLineParser.printSingleLineUsage(OutputStream)
func (p *CmdLineParser) PrintSingleLineUsage(out io.Writer) {
	Utilities_CheckNonNull(out, "OutputStream")

	p.PrintSingleLineUsageToWriter(jutil.NewOutputStreamWriter(cmdLineParser_outputStream(out)), nil)
}

// PrintSingleLineUsageToWriter prints a single-line usage to the screen.
//
// If rb is non-nil, Option.Usage is treated as a key to obtain the actual
// message from this resource bundle. It throws NullPointerException if w is
// nil.
//
// Java: CmdLineParser.printSingleLineUsage(Writer, ResourceBundle)
func (p *CmdLineParser) PrintSingleLineUsageToWriter(w jutil.Writer, rb ResourceBundle) {
	Utilities_CheckNonNull(w, "Writer")

	pw := jutil.NewPrintWriter(w)
	for _, h := range p.arguments {
		p.printSingleLineOption(pw, h, rb)
	}
	for _, h := range p.options {
		p.printSingleLineOption(pw, h, rb)
	}
	pw.Flush()
}

// Java: CmdLineParser.printSingleLineOption(PrintWriter, OptionHandler, ResourceBundle)
func (p *CmdLineParser) printSingleLineOption(pw *jutil.PrintWriter, h OptionHandler, rb ResourceBundle) {
	option := h.GetOptionHandlerBase().Option
	pw.PrintChar(' ')
	if !option.Required() {
		pw.PrintChar('[')
	}
	pw.Print(h.GetNameAndMetaWithProperties(rb, p.parserProperties))
	if option.IsMultiValued() {
		pw.Print(" ...")
	}
	if !option.Required() {
		pw.PrintChar(']')
	}
}

// cmdLineParser_outputStream is an OutputStream as the io.WriteCloser that
// jutil.OutputStreamWriter wraps; the usage printers never close it.
//
// Go-only: language-adapter
func cmdLineParser_outputStream(out io.Writer) io.WriteCloser {
	if wc, ok := out.(io.WriteCloser); ok {
		return wc
	}
	return cmdLineParser_nopCloser{out}
}

type cmdLineParser_nopCloser struct{ io.Writer }

func (cmdLineParser_nopCloser) Close() error { return nil }
