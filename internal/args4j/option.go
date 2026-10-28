// Port of org/kohsuke/args4j/Option.java @ args4j 2.33
// Port of org/kohsuke/args4j/Argument.java @ args4j 2.33
// Port of org/kohsuke/args4j/OptionDef.java @ args4j 2.33
// Port of org/kohsuke/args4j/NamedOptionDef.java @ args4j 2.33
// Port of org/kohsuke/args4j/OptionHandlerFilter.java @ args4j 2.33
// Port of org/kohsuke/args4j/ParserProperties.java @ args4j 2.33
// Port of org/kohsuke/args4j/CmdLineException.java @ args4j 2.33
// Port of org/kohsuke/args4j/IllegalAnnotationError.java @ args4j 2.33

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
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Option describes a named command-line flag bound to a Bean field.
type Option struct {
	Name     string
	Aliases  []string
	Usage    string
	MetaVar  string
	Required bool
	Help     bool
	Hidden   bool
	Handler  OptionHandlerClass
	Depends  []string
	Forbids  []string
}

// OptionHandlerClass constructs an OptionHandler. A nil value instructs args4j
// to infer the handler from the field's Class.
type OptionHandlerClass func(parser *CmdLineParser, option OptionDef, setter Setter) OptionHandler

// Argument describes a positional command-line argument bound to a Bean field.
type Argument struct {
	Usage       string
	MetaVar     string
	Required    bool
	Hidden      bool
	Handler     OptionHandlerClass
	Index       int
	MultiValued bool
}

// OptionDef is the runtime metadata of an Option or Argument. Unnamed
// instances represent positional arguments; NamedOptionDef represents flags.
type OptionDef interface {
	Usage() string
	MetaVar() string
	Required() bool
	Help() bool
	Hidden() bool
	Handler() OptionHandlerClass
	IsMultiValued() bool
	IsArgument() bool
	String() string
	GetOptionDefBase() *OptionDefBase
}

// OptionDefBase holds the common OptionDef fields and default methods.
type OptionDefBase struct {
	self        OptionDef
	usage       string
	metaVar     string
	required    bool
	help        bool
	hidden      bool
	multiValued bool
	handler     OptionHandlerClass
}

// NewOptionDef constructs an OptionDef for a positional Argument.
func NewOptionDef(a *Argument, forceMultiValued bool) *OptionDefBase {
	d := &OptionDefBase{}
	d.InitOptionDefBase(d, a.Usage, a.MetaVar, a.Required, false, a.Hidden, a.Handler, a.MultiValued || forceMultiValued)
	return d
}

// InitOptionDefBase initializes OptionDefBase for a concrete OptionDef.
func (d *OptionDefBase) InitOptionDefBase(
	self OptionDef,
	usage, metaVar string,
	required, help, hidden bool,
	handler OptionHandlerClass,
	multiValued bool,
) {
	d.self = self
	d.usage = usage
	d.metaVar = metaVar
	d.required = required
	d.help = help
	d.hidden = hidden
	d.handler = handler
	d.multiValued = multiValued
}

func (d *OptionDefBase) GetOptionDefBase() *OptionDefBase { return d }
func (d *OptionDefBase) Usage() string                    { return d.usage }
func (d *OptionDefBase) MetaVar() string                  { return d.metaVar }
func (d *OptionDefBase) Required() bool                   { return d.required }
func (d *OptionDefBase) Help() bool                       { return d.help }
func (d *OptionDefBase) Hidden() bool                     { return d.hidden }
func (d *OptionDefBase) Handler() OptionHandlerClass      { return d.handler }
func (d *OptionDefBase) IsMultiValued() bool              { return d.multiValued }
func (d *OptionDefBase) IsArgument() bool                 { return true }
func (d *OptionDefBase) String() string                   { return d.self.MetaVar() }

// NamedOptionDef is the runtime copy of an Option annotation.
type NamedOptionDef struct {
	OptionDefBase
	name    string
	aliases []string
	depends []string
	forbids []string
}

// NewNamedOptionDef constructs a NamedOptionDef from an Option.
func NewNamedOptionDef(o *Option) *NamedOptionDef {
	d := &NamedOptionDef{
		name:    o.Name,
		aliases: nonNilStrings(o.Aliases),
		depends: nonNilStrings(o.Depends),
		forbids: nonNilStrings(o.Forbids),
	}
	d.InitOptionDefBase(d, o.Usage, o.MetaVar, o.Required, o.Help, o.Hidden, o.Handler, false)
	return d
}

func nonNilStrings(in []string) []string {
	if in == nil {
		return []string{}
	}
	return in
}

func (d *NamedOptionDef) Name() string      { return d.name }
func (d *NamedOptionDef) Aliases() []string { return slices.Clone(d.aliases) }
func (d *NamedOptionDef) Depends() []string { return slices.Clone(d.depends) }
func (d *NamedOptionDef) Forbids() []string { return slices.Clone(d.forbids) }
func (d *NamedOptionDef) IsArgument() bool  { return false }

func (d *NamedOptionDef) String() string {
	if len(d.aliases) > 0 {
		return d.name + " (" + strings.Join(d.aliases, ", ") + ")"
	}
	return d.name
}

// OptionHandlerFilter selects which OptionHandlers are included in usage output.
type OptionHandlerFilter interface {
	Select(o OptionHandler) bool
}

// OptionHandlerFilterFunc adapts a function to OptionHandlerFilter.
type OptionHandlerFilterFunc func(o OptionHandler) bool

func (f OptionHandlerFilterFunc) Select(o OptionHandler) bool { return f(o) }

var (
	OptionHandlerFilter_ALL OptionHandlerFilter = OptionHandlerFilterFunc(func(OptionHandler) bool {
		return true
	})
	OptionHandlerFilter_PUBLIC OptionHandlerFilter = OptionHandlerFilterFunc(func(o OptionHandler) bool {
		return !o.GetOptionHandlerBase().Option.Hidden()
	})
	OptionHandlerFilter_REQUIRED OptionHandlerFilter = OptionHandlerFilterFunc(func(o OptionHandler) bool {
		return o.GetOptionHandlerBase().Option.Required()
	})
)

const defaultUsageWidth = 80

// ParserProperties configures CmdLineParser formatting and @-file behavior.
type ParserProperties struct {
	usageWidth           int
	optionSorter         func(o1, o2 OptionHandler) int
	optionValueDelimiter string
	atSyntax             bool
	showDefaults         bool
}

// ParserProperties_Defaults returns a new ParserProperties with standard defaults.
func ParserProperties_Defaults() *ParserProperties {
	return &ParserProperties{
		usageWidth:           defaultUsageWidth,
		optionSorter:         ParserProperties_DEFAULT_COMPARATOR,
		optionValueDelimiter: " ",
		atSyntax:             true,
		showDefaults:         true,
	}
}

func (p *ParserProperties) WithAtSyntax(atSyntax bool) *ParserProperties {
	p.atSyntax = atSyntax
	return p
}

func (p *ParserProperties) GetAtSyntax() bool { return p.atSyntax }

func (p *ParserProperties) WithShowDefaults(showDefaults bool) *ParserProperties {
	p.showDefaults = showDefaults
	return p
}

func (p *ParserProperties) GetShowDefaults() bool { return p.showDefaults }

func (p *ParserProperties) WithUsageWidth(usageWidth int) *ParserProperties {
	if usageWidth < 0 {
		panic(jutil.NewIllegalArgumentExceptionMsg("Usage width is negative"))
	}
	p.usageWidth = usageWidth
	return p
}

func (p *ParserProperties) GetUsageWidth() int { return p.usageWidth }

func (p *ParserProperties) WithOptionSorter(sorter func(o1, o2 OptionHandler) int) *ParserProperties {
	p.optionSorter = sorter
	return p
}

func (p *ParserProperties) GetOptionSorter() func(o1, o2 OptionHandler) int { return p.optionSorter }

func (p *ParserProperties) WithOptionValueDelimiter(v string) *ParserProperties {
	p.optionValueDelimiter = v
	return p
}

func (p *ParserProperties) GetOptionValueDelimiter() string { return p.optionValueDelimiter }

// ParserProperties_DEFAULT_COMPARATOR sorts options alphabetically by OptionDef.String().
var ParserProperties_DEFAULT_COMPARATOR = func(o1, o2 OptionHandler) int {
	return jutil.StringCompareTo(o1.GetOptionHandlerBase().Option.String(), o2.GetOptionHandlerBase().Option.String())
}

// cmdLineExceptionClassName matches the shaded class name in closure-compiler.jar.
const cmdLineExceptionClassName = "com.google.javascript.jscomp.jarjar.org.kohsuke.args4j.CmdLineException"

// CmdLineException signals an error in user command-line arguments.
type CmdLineException struct {
	jutil.JavaException
	parser *CmdLineParser
}

func (*CmdLineException) cmdLineException() {}

// ICmdLineException matches CmdLineException and its subclasses.
type ICmdLineException interface {
	jutil.IException
	cmdLineException()
	GetLocalizedMessage() string
	GetParser() *CmdLineParser
}

func newCmdLineException(parser *CmdLineParser, msg messageTemplate, args ...string) *CmdLineException {
	objs := make([]any, len(args))
	for i, a := range args {
		objs[i] = a
	}
	return NewCmdLineExceptionWithParser(parser, msg.format(objs...))
}

func NewCmdLineExceptionWithParser(parser *CmdLineParser, message string) *CmdLineException {
	e := &CmdLineException{parser: parser}
	jutil.InitThrowable(e, cmdLineExceptionClassName, &message, nil)
	return e
}

func NewCmdLineExceptionWithParserMessageAndCause(parser *CmdLineParser, message string, cause jutil.Throwable) *CmdLineException {
	e := &CmdLineException{parser: parser}
	jutil.InitThrowable(e, cmdLineExceptionClassName, &message, cause)
	return e
}

func (e *CmdLineException) GetLocalizedMessage() string { return e.GetMessage() }
func (e *CmdLineException) GetParser() *CmdLineParser   { return e.parser }

func (e *CmdLineException) Error() string {
	if jutil.GetMessageOrNull(e) == nil {
		return cmdLineExceptionClassName
	}
	return cmdLineExceptionClassName + ": " + e.GetMessage()
}

const illegalAnnotationErrorClassName = "com.google.javascript.jscomp.jarjar.org.kohsuke.args4j.IllegalAnnotationError"

// IllegalAnnotationError signals invalid option metadata on a Bean.
type IllegalAnnotationError struct {
	jutil.JavaLangError
}

func NewIllegalAnnotationError(message string) *IllegalAnnotationError {
	e := &IllegalAnnotationError{}
	jutil.InitThrowable(e, illegalAnnotationErrorClassName, &message, nil)
	return e
}
