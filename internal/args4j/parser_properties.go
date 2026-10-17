// Port of org/kohsuke/args4j/ParserProperties.java @ args4j 2.33

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

import "github.com/paullepoulpe/jsgo/internal/jutil"

// ParserProperties is the set of properties that controls CmdLineParser
// behaviours.
//
// See CmdLineParser(Object, ParserProperties).
//
// Java: class ParserProperties
type ParserProperties struct {
	usageWidth           int
	optionSorter         func(o1, o2 OptionHandler) int
	optionValueDelimiter string
	atSyntax             bool
	showDefaults         bool
}

const parserProperties_DEFAULT_USAGE_WIDTH = 80

// Java: ParserProperties()
func newParserProperties() *ParserProperties {
	p := &ParserProperties{}
	p.usageWidth = parserProperties_DEFAULT_USAGE_WIDTH
	p.optionSorter = ParserProperties_DEFAULT_COMPARATOR
	p.optionValueDelimiter = " "
	p.atSyntax = true
	p.showDefaults = true
	return p
}

// ParserProperties_Defaults returns an instance of the default parser
// properties. This instance can be manipulated with the withXXX() methods in
// this class.
//
// Java: ParserProperties.defaults()
func ParserProperties_Defaults() *ParserProperties {
	return newParserProperties()
}

// WithAtSyntax toggles the parsing of @-prefixes in values. If a command line
// value starts with @, it is interpreted as being a file, loaded, and
// interpreted as if the file content would have been passed to the command
// line. atSyntax is true if at sign is being parsed, false if it is to be
// ignored. Defaults to true.
//
// Java: ParserProperties.withAtSyntax(boolean)
func (p *ParserProperties) WithAtSyntax(atSyntax bool) *ParserProperties {
	p.atSyntax = atSyntax
	return p
}

// GetAtSyntax gets whether @-prefix-parsing is enabled.
//
// Java: ParserProperties.getAtSyntax()
func (p *ParserProperties) GetAtSyntax() bool {
	return p.atSyntax
}

// WithShowDefaults toggles the showing of default values in the command line
// help. showDefaults is true if to show defaults, false otherweise. Defaults
// to true.
//
// Java: ParserProperties.withShowDefaults(boolean)
func (p *ParserProperties) WithShowDefaults(showDefaults bool) *ParserProperties {
	p.showDefaults = showDefaults
	return p
}

// GetShowDefaults gets whether show defaults is enabled.
//
// Java: ParserProperties.getShowDefaults()
func (p *ParserProperties) GetShowDefaults() bool {
	return p.showDefaults
}

// WithUsageWidth sets the width of a usage line. If the usage message is
// longer than this value, the parser wraps the line.
//
// Defaults to 80. It throws IllegalArgumentException if usageWidth is
// negative.
//
// Java: ParserProperties.withUsageWidth(int)
func (p *ParserProperties) WithUsageWidth(usageWidth int) *ParserProperties {
	if usageWidth < 0 {
		panic(jutil.NewIllegalArgumentExceptionMsg("Usage width is negative"))
	}
	p.usageWidth = usageWidth
	return p
}

// GetUsageWidth returns the width of a usage line.
//
// Java: ParserProperties.getUsageWidth()
func (p *ParserProperties) GetUsageWidth() int {
	return p.usageWidth
}

// WithOptionSorter controls how options are sorted in the usage screen.
// If sorter is non-nil, options are sorted in the order induced by this
// comparator.
//
// Java: ParserProperties.withOptionSorter(Comparator)
func (p *ParserProperties) WithOptionSorter(sorter func(o1, o2 OptionHandler) int) *ParserProperties {
	p.optionSorter = sorter
	return p
}

// GetOptionSorter returns nil if options are left unsorted and should be
// listed by their discovery order. Otherwise the returned comparator is used
// to sort options. The default value is a comparator that sorts options
// alphabetically.
//
// Java: ParserProperties.getOptionSorter()
func (p *ParserProperties) GetOptionSorter() func(o1, o2 OptionHandler) int {
	return p.optionSorter
}

// WithOptionValueDelimiter sets the string used to separate option name and
// its value (such as --foo=bar vs --foo bar).
//
// Default to whitespace. Note that the tokens separated in the argument array
// (such as '-foo','bar') is always recognized as a valid name/value
// separator.
//
// Java: ParserProperties.withOptionValueDelimiter(String)
func (p *ParserProperties) WithOptionValueDelimiter(v string) *ParserProperties {
	p.optionValueDelimiter = v
	return p
}

// Java: ParserProperties.getOptionValueDelimiter()
func (p *ParserProperties) GetOptionValueDelimiter() string {
	return p.optionValueDelimiter
}

// ParserProperties_DEFAULT_COMPARATOR sorts options by the string form of
// their OptionDef (String.compareTo).
//
// Java: ParserProperties.DEFAULT_COMPARATOR
var ParserProperties_DEFAULT_COMPARATOR = func(o1, o2 OptionHandler) int {
	return jutil.StringCompareTo(o1.GetOptionHandlerBase().Option.String(), o2.GetOptionHandlerBase().Option.String())
}
