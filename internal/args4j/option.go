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
type OptionHandlerClass func(parser *CmdLineParser, option *OptionDef, setter Setter) *OptionHandler

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

// OptionDef is the runtime metadata of an Option (isNamed == true) or a
// positional Argument (isNamed == false).
type OptionDef struct {
	isNamed     bool
	name        string
	aliases     []string
	depends     []string
	forbids     []string
	usage       string
	metaVar     string
	required    bool
	help        bool
	hidden      bool
	multiValued bool
	handler     OptionHandlerClass
}

func newOptionDef(a *Argument, forceMultiValued bool) *OptionDef {
	return &OptionDef{
		usage:       a.Usage,
		metaVar:     a.MetaVar,
		required:    a.Required,
		hidden:      a.Hidden,
		handler:     a.Handler,
		multiValued: a.MultiValued || forceMultiValued,
	}
}

func newNamedOptionDef(o *Option) *OptionDef {
	return &OptionDef{
		isNamed:  true,
		name:     o.Name,
		aliases:  o.Aliases,
		depends:  o.Depends,
		forbids:  o.Forbids,
		usage:    o.Usage,
		metaVar:  o.MetaVar,
		required: o.Required,
		help:     o.Help,
		hidden:   o.Hidden,
		handler:  o.Handler,
	}
}

func (d *OptionDef) IsNamed() bool { return d.isNamed }
func (d *OptionDef) Name() string  { return d.name }
func (d *OptionDef) Hidden() bool  { return d.hidden }

func (d *OptionDef) String() string {
	if !d.isNamed {
		return d.metaVar
	}
	if len(d.aliases) > 0 {
		return d.name + " (" + strings.Join(d.aliases, ", ") + ")"
	}
	return d.name
}

// OptionHandlerFilterFunc selects which OptionHandlers are included in usage output.
type OptionHandlerFilterFunc func(o *OptionHandler) bool

var (
	OptionHandlerFilter_ALL OptionHandlerFilterFunc = func(*OptionHandler) bool {
		return true
	}
	OptionHandlerFilter_PUBLIC OptionHandlerFilterFunc = func(o *OptionHandler) bool {
		return !o.Option.hidden
	}
	OptionHandlerFilter_REQUIRED OptionHandlerFilterFunc = func(o *OptionHandler) bool {
		return o.Option.required
	}
)

const defaultUsageWidth = 80

// ParserProperties configures CmdLineParser formatting and @-file behavior.
type ParserProperties struct {
	usageWidth           int
	optionSorter         func(o1, o2 *OptionHandler) int
	optionValueDelimiter string
	atSyntax             bool
	showDefaults         bool
}

// ParserProperties_Defaults returns a new ParserProperties with standard defaults.
func ParserProperties_Defaults() *ParserProperties {
	return &ParserProperties{
		usageWidth:           defaultUsageWidth,
		optionSorter:         defaultOptionComparator,
		optionValueDelimiter: " ",
		atSyntax:             true,
		showDefaults:         true,
	}
}

func (p *ParserProperties) WithShowDefaults(showDefaults bool) *ParserProperties {
	p.showDefaults = showDefaults
	return p
}

func (p *ParserProperties) WithUsageWidth(usageWidth int) *ParserProperties {
	if usageWidth < 0 {
		panic(jutil.NewIllegalArgumentExceptionMsg("Usage width is negative"))
	}
	p.usageWidth = usageWidth
	return p
}

func (p *ParserProperties) WithOptionSorter(sorter func(o1, o2 *OptionHandler) int) *ParserProperties {
	p.optionSorter = sorter
	return p
}

func (p *ParserProperties) WithOptionValueDelimiter(v string) *ParserProperties {
	p.optionValueDelimiter = v
	return p
}

func defaultOptionComparator(o1, o2 *OptionHandler) int {
	return jutil.StringCompareTo(o1.Option.String(), o2.Option.String())
}

// cmdLineExceptionClassName matches the shaded class name in closure-compiler.jar.
const cmdLineExceptionClassName = "com.google.javascript.jscomp.jarjar.org.kohsuke.args4j.CmdLineException"

// CmdLineException signals an error in user command-line arguments.
type CmdLineException struct {
	jutil.JavaException
	parser *CmdLineParser
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

func newCmdLineExceptionWithCause(parser *CmdLineParser, message string, cause jutil.Throwable) *CmdLineException {
	e := &CmdLineException{parser: parser}
	jutil.InitThrowable(e, cmdLineExceptionClassName, &message, cause)
	return e
}

func (e *CmdLineException) GetParser() *CmdLineParser { return e.parser }

func (e *CmdLineException) Error() string {
	if jutil.GetMessageOrNull(e) == nil {
		return cmdLineExceptionClassName
	}
	return cmdLineExceptionClassName + ": " + e.GetMessage()
}

const illegalAnnotationErrorClassName = "com.google.javascript.jscomp.jarjar.org.kohsuke.args4j.IllegalAnnotationError"

type illegalAnnotationError struct {
	jutil.JavaLangError
}

func newIllegalAnnotationError(message string) *illegalAnnotationError {
	e := &illegalAnnotationError{}
	jutil.InitThrowable(e, illegalAnnotationErrorClassName, &message, nil)
	return e
}
