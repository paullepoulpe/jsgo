// Port of org/kohsuke/args4j/spi/OptionHandler.java @ args4j 2.33
// Port of org/kohsuke/args4j/OptionHandlerRegistry.java @ args4j 2.33
// Port of org/kohsuke/args4j/spi/OneArgumentOptionHandler.java @ args4j 2.33
// Port of org/kohsuke/args4j/spi/BooleanOptionHandler.java @ args4j 2.33
// Port of org/kohsuke/args4j/spi/StringOptionHandler.java @ args4j 2.33
// Port of org/kohsuke/args4j/spi/IntOptionHandler.java @ args4j 2.33
// Port of org/kohsuke/args4j/spi/EnumOptionHandler.java @ args4j 2.33
// Port of org/kohsuke/args4j/spi/MapOptionHandler.java @ args4j 2.33
// Port of org/kohsuke/args4j/spi/StringArrayOptionHandler.java @ args4j 2.33

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

// OptionHandler parses command-line operands and formats metadata for a single
// Option or Argument.
type OptionHandler struct {
	Option              *OptionDef
	Setter              Setter
	Owner               *CmdLineParser
	parse               func(h *OptionHandler, params Parameters) int
	parseErr            func(h *OptionHandler, params Parameters) (int, ICmdLineException)
	defaultMetaVar      string
	ignoreOptionMetaVar bool
}

// NewOptionHandler constructs a concrete OptionHandler with custom parsing logic.
func NewOptionHandler(
	parser *CmdLineParser,
	option *OptionDef,
	setter Setter,
	defaultMetaVar string,
	parse func(h *OptionHandler, params Parameters) int,
) *OptionHandler {
	return &OptionHandler{
		Owner:          parser,
		Option:         option,
		Setter:         setter,
		defaultMetaVar: defaultMetaVar,
		parse:          parse,
	}
}

func newBuiltinOptionHandler(
	parser *CmdLineParser,
	option *OptionDef,
	setter Setter,
	defaultMetaVar string,
	parseErr func(h *OptionHandler, params Parameters) (int, ICmdLineException),
) *OptionHandler {
	return &OptionHandler{
		Owner:          parser,
		Option:         option,
		Setter:         setter,
		defaultMetaVar: defaultMetaVar,
		parseErr:       parseErr,
	}
}

func (h *OptionHandler) ParseArguments(params Parameters) int {
	if h.parseErr != nil {
		n, err := h.parseErr(h, params)
		if err != nil {
			panic(err)
		}
		return n
	}
	return h.parse(h, params)
}

func (h *OptionHandler) tryParseArguments(params Parameters) (int, ICmdLineException) {
	if h.parseErr != nil {
		return h.parseErr(h, params)
	}
	return h.parse(h, params), nil
}

func tryGetParameter(params Parameters, idx int) (string, ICmdLineException) {
	if cp, ok := params.(*cmdLineParams); ok {
		return cp.tryGetParameter(idx)
	}
	return params.GetParameter(idx), nil
}

func (h *OptionHandler) GetDefaultMetaVariable() string {
	return h.defaultMetaVar
}

func (h *OptionHandler) GetMetaVariable(rb ResourceBundle) string {
	if h.ignoreOptionMetaVar {
		return h.defaultMetaVar
	}
	token := h.Option.MetaVar()
	if len(token) == 0 {
		token = h.defaultMetaVar
	}
	if len(token) == 0 {
		return ""
	}
	if rb != nil {
		return rb.GetString(token)
	}
	return token
}

func (h *OptionHandler) GetNameAndMeta(rb ResourceBundle) string {
	return h.GetNameAndMetaWithProperties(rb, ParserProperties_Defaults())
}

func (h *OptionHandler) GetNameAndMetaWithProperties(rb ResourceBundle, properties *ParserProperties) string {
	str := ""
	if !h.Option.IsArgument() {
		str = h.Option.String()
	}
	if meta := h.GetMetaVariable(rb); len(meta) > 0 {
		if len(str) > 0 {
			str += properties.GetOptionValueDelimiter()
		}
		str += meta
	}
	return str
}

func (h *OptionHandler) Print(v any) string {
	return jutil.ToJavaString(v)
}

func (h *OptionHandler) PrintDefaultValue() (string, bool) {
	getter, ok := h.Setter.(Getter)
	if !ok {
		return "", false
	}
	defaultValues := getter.GetValueList()
	if len(defaultValues) == 0 {
		return "", false
	}
	if len(defaultValues) == 1 {
		return h.Print(defaultValues[0]), true
	}
	var buf strings.Builder
	buf.WriteByte('[')
	for i, v := range defaultValues {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(h.Print(v))
	}
	buf.WriteByte(']')
	return buf.String(), true
}

var defaultOptionHandlers = map[*Class]OptionHandlerClass{
	Class_Boolean: NewBooleanOptionHandler,
	Class_boolean: NewBooleanOptionHandler,
	Class_Integer: NewIntOptionHandler,
	Class_int:     NewIntOptionHandler,
	Class_String:  NewStringOptionHandler,
	Class_Map:     NewMapOptionHandler,
}

func createOptionHandler(parser *CmdLineParser, o *OptionDef, setter Setter) *OptionHandler {
	checkNonNull(o, "OptionDef")
	checkNonNull(setter, "Setter")

	if h := o.Handler(); h != nil {
		return h(parser, o, setter)
	}
	t := setter.GetType()
	if Class_Enum.IsAssignableFrom(t) {
		return NewEnumOptionHandler(parser, o, setter, t)
	}
	cons, ok := defaultOptionHandlers[t]
	if !ok {
		panic(NewIllegalAnnotationError(msgUnknownHandler.format(t)))
	}
	return cons(parser, o, setter)
}

var (
	acceptableBooleanTrue  = []string{"true", "on", "yes", "1"}
	acceptableBooleanFalse = []string{"false", "off", "no", "0"}
)

// NewBooleanOptionHandler parses boolean flags and positional boolean arguments.
func NewBooleanOptionHandler(parser *CmdLineParser, option *OptionDef, setter Setter) *OptionHandler {
	return newBuiltinOptionHandler(parser, option, setter, "", func(h *OptionHandler, params Parameters) (int, ICmdLineException) {
		if !h.Option.IsArgument() {
			h.Setter.AddValue(true)
			return 0, nil
		}
		raw, err := tryGetParameter(params, 0)
		if err != nil {
			return 0, err
		}
		valueStr := jutil.StringToLowerCaseRoot(raw)
		switch {
		case slices.Contains(acceptableBooleanTrue, valueStr):
			h.Setter.AddValue(true)
		case slices.Contains(acceptableBooleanFalse, valueStr):
			h.Setter.AddValue(false)
		default:
			return 0, newCmdLineException(h.Owner, spiMsgIllegalBoolean, valueStr)
		}
		return 1, nil
	})
}

// NewStringOptionHandler parses a single string operand.
func NewStringOptionHandler(parser *CmdLineParser, option *OptionDef, setter Setter) *OptionHandler {
	return newBuiltinOptionHandler(parser, option, setter, defaultMetaString, func(h *OptionHandler, params Parameters) (int, ICmdLineException) {
		v, err := tryGetParameter(params, 0)
		if err != nil {
			return 0, err
		}
		if ss, ok := h.Setter.(StringSetter); ok {
			ss.AddStringValue(v)
		} else {
			h.Setter.AddValue(v)
		}
		return 1, nil
	})
}

// NewIntOptionHandler parses a single base-10 32-bit integer operand.
func NewIntOptionHandler(parser *CmdLineParser, option *OptionDef, setter Setter) *OptionHandler {
	return newBuiltinOptionHandler(parser, option, setter, "N", func(h *OptionHandler, params Parameters) (int, ICmdLineException) {
		token, err := tryGetParameter(params, 0)
		if err != nil {
			return 0, err
		}
		val, ok := jutil.TryIntegerParseInt(token, 10)
		if !ok {
			prev, prevErr := tryGetParameter(params, -1)
			if prevErr != nil {
				return 0, prevErr
			}
			return 0, newCmdLineException(h.Owner, spiMsgIllegalOperand, prev, token)
		}
		h.Setter.AddValue(val)
		return 1, nil
	})
}

// NewEnumOptionHandler parses a case-insensitive enum constant name (with '-' normalized to '_').
func NewEnumOptionHandler(parser *CmdLineParser, option *OptionDef, setter Setter, enumType *Class) *OptionHandler {
	constants := enumType.GetEnumConstants()
	var buf strings.Builder
	buf.WriteByte('[')
	for i, c := range constants {
		if i > 0 {
			buf.WriteString(" | ")
		}
		buf.WriteString(jutil.ToJavaString(c))
	}
	buf.WriteByte(']')
	meta := buf.String()

	h := newBuiltinOptionHandler(parser, option, setter, meta, func(h *OptionHandler, params Parameters) (int, ICmdLineException) {
		raw, err := tryGetParameter(params, 0)
		if err != nil {
			return 0, err
		}
		s := strings.ReplaceAll(raw, "-", "_")
		var value any
		for _, o := range enumType.constants {
			if strings.EqualFold(enumName(o), s) {
				value = o
				break
			}
		}
		if value == nil {
			var optName string
			if h.Option.IsArgument() {
				optName = h.Option.String()
			} else {
				var prevErr ICmdLineException
				optName, prevErr = tryGetParameter(params, -1)
				if prevErr != nil {
					return 0, prevErr
				}
			}
			return 0, newCmdLineException(h.Owner, spiMsgIllegalOperand, optName, s)
		}
		h.Setter.AddValue(value)
		return 1, nil
	})
	h.ignoreOptionMetaVar = true
	return h
}

// NewMapOptionHandler parses "key=value" operands into a Map field.
func NewMapOptionHandler(parser *CmdLineParser, option *OptionDef, setter Setter) *OptionHandler {
	if setter.AsFieldSetter() == nil {
		panic(jutil.NewIllegalArgumentExceptionMsg("MapOptionHandler can only work with fields"))
	}
	return newBuiltinOptionHandler(parser, option, setter, "", func(h *OptionHandler, params Parameters) (int, ICmdLineException) {
		fs := h.Setter.AsFieldSetter()
		var m Map
		if o := fs.GetValue(); o != nil {
			m = o.(Map)
		}
		if m == nil {
			m = jutil.NewLinkedHashMap[string, *string]()
			fs.AddValue(m)
		}
		raw, err := tryGetParameter(params, 0)
		if err != nil {
			return 0, err
		}
		mapKey, value, ok := strings.Cut(raw, "=")
		if !ok {
			return 0, newCmdLineException(h.Owner, spiMsgFormatErrorForMap)
		}
		if len(mapKey) == 0 {
			return 0, newCmdLineException(h.Owner, spiMsgMapHasNoKey)
		}
		var mapValue *string
		if len(value) > 0 {
			mapValue = &value
		}
		m.Put(mapKey, mapValue)
		return 1, nil
	})
}

// NewStringArrayOptionHandler greedily consumes non-option tokens (splitting each
// on literal spaces) into a multi-valued String setter.
func NewStringArrayOptionHandler(parser *CmdLineParser, option *OptionDef, setter Setter) *OptionHandler {
	return newBuiltinOptionHandler(parser, option, setter, defaultMetaStringArray, func(h *OptionHandler, params Parameters) (int, ICmdLineException) {
		counter := 0
		for ; counter < params.Size(); counter++ {
			param, err := tryGetParameter(params, counter)
			if err != nil {
				return 0, err
			}
			if strings.HasPrefix(param, "-") {
				break
			}
			for _, p := range jutil.StringSplitLiteral(param, " ") {
				h.Setter.AddValue(p)
			}
		}
		return counter, nil
	})
}
