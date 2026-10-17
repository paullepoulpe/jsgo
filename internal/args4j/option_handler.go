// Port of org/kohsuke/args4j/spi/OptionHandler.java @ args4j 2.33

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

// OptionHandler is code that parses operands of an option into Java.
//
// This class can be extended by application to support additional Java
// datatypes in option operands.
//
// Implementation of this class needs to be registered to args4j by using
// OptionHandlerRegistry.registerHandler(Class,Class). For registration to
// work, subclasses will need to implement the constructor with the signature
// OptionHandler(CmdLineParser, OptionDef, Setter).
//
// T is the component type of the field that this OptionHandler works with.
// When I say "component", I mean a field that can hold multiple values (such
// as Collection or array). This should refer to its component time.
//
// Setter implementations abstract away multi-value-ness by allowing
// OptionHandler to invoke its Setter.addValue(Object) multiple times.
//
// Go-only: language-adapter — generics: values of type T are passed as any.
//
// Java: abstract class spi.OptionHandler<T>
type OptionHandler interface {
	// ParseArguments is called if the option that this owner recognizes is
	// found.
	//
	// params is the rest of the arguments. This method can use this object
	// to access the arguments of the option if necessary. The object is
	// valid only during the method call.
	//
	// It returns the number of arguments consumed. (For example, returns 0
	// if this option doesn't take any parameters.) It may throw a
	// CmdLineException.
	//
	// Java: OptionHandler.parseArguments(Parameters)
	ParseArguments(params Parameters) int

	// GetDefaultMetaVariable gets the default meta variable name used to
	// print the usage screen.
	//
	// The value returned by this method can be a reference in the
	// ResourceBundle, if one was passed to CmdLineParser.
	//
	// It returns nil to hide a meta variable.
	//
	// Java: OptionHandler.getDefaultMetaVariable()
	GetDefaultMetaVariable() *string

	// Java: OptionHandler.getMetaVariable(ResourceBundle)
	GetMetaVariable(rb ResourceBundle) *string

	// Java: OptionHandler.getNameAndMeta(ResourceBundle)
	GetNameAndMeta(rb ResourceBundle) string

	// Java: OptionHandler.getNameAndMeta(ResourceBundle, ParserProperties)
	GetNameAndMetaWithProperties(rb ResourceBundle, properties *ParserProperties) string

	// Java: OptionHandler.print(T)
	Print(v any) string

	// Java: OptionHandler.printDefaultValue()
	PrintDefaultValue() *string

	GetOptionHandlerBase() *OptionHandlerBase
}

// OptionHandlerBase holds the fields of OptionHandler and the
// implementations of its non-abstract methods.
type OptionHandlerBase struct {
	self OptionHandler

	// The annotation.
	Option OptionDef

	// Object to be used for setting value.
	Setter Setter

	// The owner to which this handler belongs to.
	Owner *CmdLineParser
}

// InitOptionHandlerBase is the protected constructor, called by
// subclasses; self is the subclass instance.
//
// Java: OptionHandler(CmdLineParser, OptionDef, Setter)
func (h *OptionHandlerBase) InitOptionHandlerBase(self OptionHandler, parser *CmdLineParser, option OptionDef, setter Setter) {
	h.self = self
	h.Owner = parser
	h.Option = option
	h.Setter = setter
}

// GetOptionHandlerBase returns the OptionHandler fields.
func (h *OptionHandlerBase) GetOptionHandlerBase() *OptionHandlerBase { return h }

// Java: OptionHandler.getMetaVariable(ResourceBundle)
func (h *OptionHandlerBase) GetMetaVariable(rb ResourceBundle) *string {
	token := h.Option.MetaVar()
	var tokenOrNull *string = &token
	if jutil.StringLength(token) == 0 {
		tokenOrNull = h.self.GetDefaultMetaVariable()
	}
	if tokenOrNull == nil {
		return nil
	}

	if rb != nil {
		localized := rb.GetString(*tokenOrNull)
		// Go-only: ResourceBundle.getString never returns null (it throws
		// MissingResourceException), so Java's null test always passes.
		tokenOrNull = &localized
	}

	return tokenOrNull
}

// GetNameAndMeta gets string representing usage for this option, of the
// form "name metaval", e.g. "-foo VALUE" or "--foo VALUE". rb is the
// ResourceBundle to get localized version of meta string.
//
// Java: OptionHandler.getNameAndMeta(ResourceBundle)
func (h *OptionHandlerBase) GetNameAndMeta(rb ResourceBundle) string {
	return h.self.GetNameAndMetaWithProperties(rb, ParserProperties_Defaults())
}

// GetNameAndMetaWithProperties gets string representing usage for this
// option, of the form "name metaval" or "name=metaval, e.g. "--foo VALUE"
// or "--foo=VALUE". rb is the ResourceBundle to get localized version of
// meta string; properties affects the formatting behaviours.
//
// Java: OptionHandler.getNameAndMeta(ResourceBundle, ParserProperties)
func (h *OptionHandlerBase) GetNameAndMetaWithProperties(rb ResourceBundle, properties *ParserProperties) string {
	str := ""
	if !h.Option.IsArgument() {
		str = h.Option.String()
	}
	meta := h.self.GetMetaVariable(rb)
	if meta != nil {
		if jutil.StringLength(str) > 0 {
			str += properties.GetOptionValueDelimiter()
		}
		str += *meta
	}
	return str
}

// Print is the opposite of the parse operation.
//
// This method is used to print the usage screen.
//
// Java: OptionHandler.print(T)
func (h *OptionHandlerBase) Print(v any) string {
	return jutil.ToJavaString(v)
}

// PrintDefaultValue prints the default value by introspecting the current
// setter as Getter.
//
// It returns nil if the current value of the setter isn't available.
//
// Java: OptionHandler.printDefaultValue()
func (h *OptionHandlerBase) PrintDefaultValue() *string {
	if getter, ok := h.Setter.(Getter); ok {
		defaultValues := getter.GetValueList()
		if defaultValues != nil && len(defaultValues) != 0 {
			buf := jutil.NewStringBuilder()
			if len(defaultValues) > 1 {
				for _, v := range defaultValues {
					if buf.Length() == 0 {
						buf.AppendChar('[')
					} else {
						buf.AppendChar(',')
					}
					buf.Append(h.self.Print(v))
				}
				buf.AppendChar(']')
			} else {
				buf.Append(h.self.Print(defaultValues[0]))
			}
			s := buf.String()
			return &s
		}
	}
	return nil
}
