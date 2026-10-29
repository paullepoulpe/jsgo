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
	parseErr            func(h *OptionHandler, params Parameters) (int, *CmdLineException)
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
		parseErr: func(h *OptionHandler, params Parameters) (int, *CmdLineException) {
			return parse(h, params), nil
		},
	}
}

func newBuiltinOptionHandler(
	parser *CmdLineParser,
	option *OptionDef,
	setter Setter,
	defaultMetaVar string,
	parseErr func(h *OptionHandler, params Parameters) (int, *CmdLineException),
) *OptionHandler {
	return &OptionHandler{
		Owner:          parser,
		Option:         option,
		Setter:         setter,
		defaultMetaVar: defaultMetaVar,
		parseErr:       parseErr,
	}
}

func (h *OptionHandler) getMetaVariable() string {
	if h.ignoreOptionMetaVar || len(h.Option.metaVar) == 0 {
		return h.defaultMetaVar
	}
	return h.Option.metaVar
}

func (h *OptionHandler) getNameAndMeta(properties *ParserProperties) string {
	str := ""
	if h.Option.isNamed {
		str = h.Option.String()
	}
	if meta := h.getMetaVariable(); len(meta) > 0 {
		if len(str) > 0 {
			str += properties.optionValueDelimiter
		}
		str += meta
	}
	return str
}

func (h *OptionHandler) printDefaultValue() (string, bool) {
	fs, ok := h.Setter.(*FieldSetter)
	if !ok {
		return "", false
	}
	defaultValues := fs.getValueList()
	if len(defaultValues) == 0 {
		return "", false
	}
	if len(defaultValues) == 1 {
		return jutil.ToJavaString(defaultValues[0]), true
	}
	var buf strings.Builder
	buf.WriteByte('[')
	for i, v := range defaultValues {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(jutil.ToJavaString(v))
	}
	buf.WriteByte(']')
	return buf.String(), true
}

var defaultOptionHandlers = map[*Class]OptionHandlerClass{
	classBoolean:     newBooleanOptionHandler,
	classBooleanPrim: newBooleanOptionHandler,
	classInteger:     NewIntOptionHandler,
	classIntPrim:     NewIntOptionHandler,
	Class_String:     NewStringOptionHandler,
	classMap:         newMapOptionHandler,
}

func createOptionHandler(parser *CmdLineParser, o *OptionDef, setter Setter) *OptionHandler {
	if o.handler != nil {
		return o.handler(parser, o, setter)
	}
	t := setter.GetType()
	if t.isEnum {
		return newEnumOptionHandler(parser, o, setter, t)
	}
	cons, ok := defaultOptionHandlers[t]
	if !ok {
		panic(newIllegalAnnotationError(msgUnknownHandler.format(t)))
	}
	return cons(parser, o, setter)
}

var (
	acceptableBooleanTrue  = []string{"true", "on", "yes", "1"}
	acceptableBooleanFalse = []string{"false", "off", "no", "0"}
)

func newBooleanOptionHandler(parser *CmdLineParser, option *OptionDef, setter Setter) *OptionHandler {
	return newBuiltinOptionHandler(parser, option, setter, "", func(h *OptionHandler, params Parameters) (int, *CmdLineException) {
		if h.Option.isNamed {
			h.Setter.AddValue(true)
			return 0, nil
		}
		raw, err := params.TryGetParameter(0)
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
	return newBuiltinOptionHandler(parser, option, setter, defaultMetaString, func(h *OptionHandler, params Parameters) (int, *CmdLineException) {
		v, err := params.TryGetParameter(0)
		if err != nil {
			return 0, err
		}
		h.Setter.AddStringValue(v)
		return 1, nil
	})
}

// NewIntOptionHandler parses a single base-10 32-bit integer operand.
func NewIntOptionHandler(parser *CmdLineParser, option *OptionDef, setter Setter) *OptionHandler {
	return newBuiltinOptionHandler(parser, option, setter, "N", func(h *OptionHandler, params Parameters) (int, *CmdLineException) {
		token, err := params.TryGetParameter(0)
		if err != nil {
			return 0, err
		}
		val, ok := jutil.TryIntegerParseInt(token, 10)
		if !ok {
			prev, prevErr := params.TryGetParameter(-1)
			if prevErr != nil {
				return 0, prevErr
			}
			return 0, newCmdLineException(h.Owner, spiMsgIllegalOperand, prev, token)
		}
		h.Setter.AddValue(val)
		return 1, nil
	})
}

func newEnumOptionHandler(parser *CmdLineParser, option *OptionDef, setter Setter, enumType *Class) *OptionHandler {
	var buf strings.Builder
	buf.WriteByte('[')
	for i, c := range enumType.constants {
		if i > 0 {
			buf.WriteString(" | ")
		}
		buf.WriteString(jutil.ToJavaString(c))
	}
	buf.WriteByte(']')
	meta := buf.String()

	h := newBuiltinOptionHandler(parser, option, setter, meta, func(h *OptionHandler, params Parameters) (int, *CmdLineException) {
		raw, err := params.TryGetParameter(0)
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
			if !h.Option.isNamed {
				optName = h.Option.String()
			} else {
				var prevErr *CmdLineException
				optName, prevErr = params.TryGetParameter(-1)
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

func newMapOptionHandler(parser *CmdLineParser, option *OptionDef, setter Setter) *OptionHandler {
	if setter.AsFieldSetter() == nil {
		panic(jutil.NewIllegalArgumentExceptionMsg("MapOptionHandler can only work with fields"))
	}
	return newBuiltinOptionHandler(parser, option, setter, "", func(h *OptionHandler, params Parameters) (int, *CmdLineException) {
		fs := h.Setter.AsFieldSetter()
		var m *jutil.LinkedHashMap[string, *string]
		if o := fs.getValue(); o != nil {
			m = o.(*jutil.LinkedHashMap[string, *string])
		}
		if m == nil {
			m = jutil.NewLinkedHashMap[string, *string]()
			fs.AddValue(m)
		}
		raw, err := params.TryGetParameter(0)
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
	return newBuiltinOptionHandler(parser, option, setter, defaultMetaStringArray, func(h *OptionHandler, params Parameters) (int, *CmdLineException) {
		counter := 0
		for ; counter < params.Size(); counter++ {
			param, err := params.TryGetParameter(counter)
			if err != nil {
				return 0, err
			}
			if strings.HasPrefix(param, "-") {
				break
			}
			for _, p := range jutil.StringSplitLiteral(param, " ") {
				h.Setter.AddStringValue(p)
			}
		}
		return counter, nil
	})
}
