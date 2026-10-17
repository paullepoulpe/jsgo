// Port of org/kohsuke/args4j/OptionDef.java @ args4j 2.33

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

// OptionDef is the run-time copy of the Option or Argument annotation. By
// definition, unnamed options are arguments (and instances of this class).
// Named options are actually a subclass.
//
// Java: class OptionDef
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

// OptionDefBase holds the fields of OptionDef and its implementation.
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

// Java: OptionDef(Argument, boolean)
func NewOptionDef(a *Argument, forceMultiValued bool) *OptionDefBase {
	d := &OptionDefBase{}
	d.InitOptionDefBase(d, a.Usage, a.MetaVar, a.Required, false, a.Hidden, a.Handler, a.MultiValued || forceMultiValued)
	return d
}

// InitOptionDefBase is the protected constructor, called by subclasses.
//
// Java: OptionDef(String, String, boolean, boolean, boolean, Class, boolean)
func (d *OptionDefBase) InitOptionDefBase(self OptionDef, usage string, metaVar string, required bool,
	help bool, hidden bool,
	handler OptionHandlerClass,
	multiValued bool) {
	d.self = self
	d.usage = usage
	d.metaVar = metaVar
	d.required = required
	d.help = help
	d.hidden = hidden
	d.handler = handler
	d.multiValued = multiValued
}

// GetOptionDefBase returns the OptionDef fields.
func (d *OptionDefBase) GetOptionDefBase() *OptionDefBase { return d }

// Java: OptionDef.usage()
func (d *OptionDefBase) Usage() string {
	return d.usage
}

// Java: OptionDef.metaVar()
func (d *OptionDefBase) MetaVar() string {
	return d.metaVar
}

// Java: OptionDef.required()
func (d *OptionDefBase) Required() bool {
	return d.required
}

// Java: OptionDef.help()
func (d *OptionDefBase) Help() bool {
	return d.help
}

// Value from Option.Hidden or Argument.Hidden.
//
// Java: OptionDef.hidden()
func (d *OptionDefBase) Hidden() bool {
	return d.hidden
}

// Java: OptionDef.handler()
func (d *OptionDefBase) Handler() OptionHandlerClass {
	return d.handler
}

// Java: OptionDef.isMultiValued()
func (d *OptionDefBase) IsMultiValued() bool {
	return d.multiValued
}

// Java: OptionDef.isArgument()
func (d *OptionDefBase) IsArgument() bool {
	return true
}

// Java: OptionDef.toString()
func (d *OptionDefBase) String() string {
	// Go-only: language-adapter — a Go string is never null: the Java test
	// metaVar() != null ? metaVar() : "ARG" always takes metaVar() for the
	// annotation values the port reads (an annotation element is never null).
	return d.self.MetaVar()
}
