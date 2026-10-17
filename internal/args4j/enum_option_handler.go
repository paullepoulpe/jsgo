// Port of org/kohsuke/args4j/spi/EnumOptionHandler.java @ args4j 2.33

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

// EnumOptionHandler is the Enum OptionHandler.
//
// Java: class spi.EnumOptionHandler<T extends Enum<T>> extends OptionHandler<T>
type EnumOptionHandler struct {
	OptionHandlerBase

	enumType *Class
}

// Java: EnumOptionHandler(CmdLineParser, OptionDef, Setter, Class)
func NewEnumOptionHandler(parser *CmdLineParser, option OptionDef, setter Setter, enumType *Class) *EnumOptionHandler {
	h := &EnumOptionHandler{}
	h.InitOptionHandlerBase(h, parser, option, setter)
	h.enumType = enumType
	return h
}

// Java: EnumOptionHandler.parseArguments(Parameters)
func (h *EnumOptionHandler) ParseArguments(params Parameters) int {
	s := jutil.StringReplaceAll(params.GetParameter(0), "-", "_")
	var value any = nil
	for _, o := range h.enumType.GetEnumConstants() {
		if enumOptionHandler_equalsIgnoreCase(enumName(o), s) {
			value = o
			break
		}
	}

	if value == nil {
		if h.Option.IsArgument() {
			panic(NewCmdLineExceptionWithLocalizable(h.Owner, SpiMessages_ILLEGAL_OPERAND, h.Option.String(), s))
		} else {
			panic(NewCmdLineExceptionWithLocalizable(h.Owner, SpiMessages_ILLEGAL_OPERAND, params.GetParameter(-1), s))
		}
	}
	h.Setter.AddValue(value)
	return 1
}

// enumOptionHandler_equalsIgnoreCase is String.equalsIgnoreCase(String).
func enumOptionHandler_equalsIgnoreCase(a, b string) bool {
	ua, ub := jutil.StringToUTF16(a), jutil.StringToUTF16(b)
	return len(ua) == len(ub) && jutil.UTF16RegionMatches(true, ua, 0, ub, 0, len(ua))
}

// GetDefaultMetaVariable returns all values of an enum type split by pipe:
// [ one | two | three].
//
// Java: EnumOptionHandler.getDefaultMetaVariable()
func (h *EnumOptionHandler) GetDefaultMetaVariable() *string {
	rv := jutil.NewStringBuilder()
	rv.Append("[")
	for _, t := range h.enumType.GetEnumConstants() {
		rv.Append(jutil.ToJavaString(t)).Append(" | ")
	}
	// Go-only: StringBuffer.delete(length-3, length); for an enum without
	// constants Java throws StringIndexOutOfBoundsException there.
	rv.SetLength(rv.Length() - 3)
	rv.Append("]")
	s := rv.String()
	return &s
}

// Java: EnumOptionHandler.getMetaVariable(ResourceBundle)
func (h *EnumOptionHandler) GetMetaVariable(rb ResourceBundle) *string {
	return h.GetDefaultMetaVariable()
}
