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

// OptionHandler parses command-line operands for a single Option or Argument.
type OptionHandler interface {
	ParseArguments(params Parameters) int
	GetDefaultMetaVariable() *string
	GetMetaVariable(rb ResourceBundle) *string
	GetNameAndMeta(rb ResourceBundle) string
	GetNameAndMetaWithProperties(rb ResourceBundle, properties *ParserProperties) string
	Print(v any) string
	PrintDefaultValue() *string
	GetOptionHandlerBase() *OptionHandlerBase
}

// OptionHandlerBase implements the shared state and formatting methods of
// OptionHandler. Concrete handlers embed OptionHandlerBase and call
// InitOptionHandlerBase with themselves as self.
type OptionHandlerBase struct {
	self   OptionHandler
	Option OptionDef
	Setter Setter
	Owner  *CmdLineParser
}

func (h *OptionHandlerBase) InitOptionHandlerBase(self OptionHandler, parser *CmdLineParser, option OptionDef, setter Setter) {
	h.self = self
	h.Owner = parser
	h.Option = option
	h.Setter = setter
}

func (h *OptionHandlerBase) GetOptionHandlerBase() *OptionHandlerBase { return h }

func (h *OptionHandlerBase) GetMetaVariable(rb ResourceBundle) *string {
	token := h.Option.MetaVar()
	tokenOrNull := &token
	if len(token) == 0 {
		tokenOrNull = h.self.GetDefaultMetaVariable()
	}
	if tokenOrNull == nil {
		return nil
	}
	if rb != nil {
		localized := rb.GetString(*tokenOrNull)
		tokenOrNull = &localized
	}
	return tokenOrNull
}

func (h *OptionHandlerBase) GetNameAndMeta(rb ResourceBundle) string {
	return h.self.GetNameAndMetaWithProperties(rb, ParserProperties_Defaults())
}

func (h *OptionHandlerBase) GetNameAndMetaWithProperties(rb ResourceBundle, properties *ParserProperties) string {
	str := ""
	if !h.Option.IsArgument() {
		str = h.Option.String()
	}
	if meta := h.self.GetMetaVariable(rb); meta != nil {
		if len(str) > 0 {
			str += properties.GetOptionValueDelimiter()
		}
		str += *meta
	}
	return str
}

func (h *OptionHandlerBase) Print(v any) string {
	return jutil.ToJavaString(v)
}

func (h *OptionHandlerBase) PrintDefaultValue() *string {
	getter, ok := h.Setter.(Getter)
	if !ok {
		return nil
	}
	defaultValues := getter.GetValueList()
	if len(defaultValues) == 0 {
		return nil
	}
	if len(defaultValues) == 1 {
		s := h.self.Print(defaultValues[0])
		return &s
	}
	var buf strings.Builder
	buf.WriteByte('[')
	for i, v := range defaultValues {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(h.self.Print(v))
	}
	buf.WriteByte(']')
	s := buf.String()
	return &s
}

var defaultOptionHandlers = map[*Class]OptionHandlerClass{
	Class_Boolean: NewBooleanOptionHandler,
	Class_boolean: NewBooleanOptionHandler,
	Class_Integer: NewIntOptionHandler,
	Class_int:     NewIntOptionHandler,
	Class_String:  NewStringOptionHandler,
	Class_Map:     NewMapOptionHandler,
}

func createOptionHandler(parser *CmdLineParser, o OptionDef, setter Setter) OptionHandler {
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

// BooleanOptionHandler parses boolean flags and positional boolean arguments.
type BooleanOptionHandler struct {
	OptionHandlerBase
}

func NewBooleanOptionHandler(parser *CmdLineParser, option OptionDef, setter Setter) OptionHandler {
	h := &BooleanOptionHandler{}
	h.InitOptionHandlerBase(h, parser, option, setter)
	return h
}

func (h *BooleanOptionHandler) ParseArguments(params Parameters) int {
	if !h.Option.IsArgument() {
		h.Setter.AddValue(true)
		return 0
	}
	valueStr := jutil.StringToLowerCaseRoot(params.GetParameter(0))
	switch {
	case slices.Contains(acceptableBooleanTrue, valueStr):
		h.Setter.AddValue(true)
	case slices.Contains(acceptableBooleanFalse, valueStr):
		h.Setter.AddValue(false)
	default:
		panic(newCmdLineException(h.Owner, spiMsgIllegalBoolean, valueStr))
	}
	return 1
}

func (h *BooleanOptionHandler) GetDefaultMetaVariable() *string { return nil }

// StringOptionHandler parses a single string operand. Subclasses embed it and
// call InitStringOptionHandler.
type StringOptionHandler struct {
	OptionHandlerBase
}

func NewStringOptionHandler(parser *CmdLineParser, option OptionDef, setter Setter) OptionHandler {
	h := &StringOptionHandler{}
	h.InitStringOptionHandler(h, parser, option, setter)
	return h
}

func (h *StringOptionHandler) InitStringOptionHandler(self OptionHandler, parser *CmdLineParser, option OptionDef, setter Setter) {
	h.InitOptionHandlerBase(self, parser, option, setter)
}

func (h *StringOptionHandler) ParseArguments(params Parameters) int {
	v := params.GetParameter(0)
	if ss, ok := h.Setter.(StringSetter); ok {
		ss.AddStringValue(v)
	} else {
		h.Setter.AddValue(v)
	}
	return 1
}

func (h *StringOptionHandler) GetDefaultMetaVariable() *string {
	s := defaultMetaString
	return &s
}

// IntOptionHandler parses a single base-10 32-bit integer operand.
type IntOptionHandler struct {
	OptionHandlerBase
}

func NewIntOptionHandler(parser *CmdLineParser, option OptionDef, setter Setter) OptionHandler {
	h := &IntOptionHandler{}
	h.InitOptionHandlerBase(h, parser, option, setter)
	return h
}

func (h *IntOptionHandler) GetDefaultMetaVariable() *string {
	s := "N"
	return &s
}

func (h *IntOptionHandler) ParseArguments(params Parameters) int {
	token := params.GetParameter(0)
	jutil.Try(func() {
		h.Setter.AddValue(jutil.IntegerParseInt(token, 10))
	}, func(jutil.INumberFormatException) {
		panic(newCmdLineException(h.Owner, spiMsgIllegalOperand, params.GetParameter(-1), token))
	})
	return 1
}

// EnumOptionHandler parses a case-insensitive enum constant name (with '-' normalized to '_').
type EnumOptionHandler struct {
	OptionHandlerBase
	enumType *Class
}

func NewEnumOptionHandler(parser *CmdLineParser, option OptionDef, setter Setter, enumType *Class) *EnumOptionHandler {
	h := &EnumOptionHandler{enumType: enumType}
	h.InitOptionHandlerBase(h, parser, option, setter)
	return h
}

func (h *EnumOptionHandler) ParseArguments(params Parameters) int {
	s := strings.ReplaceAll(params.GetParameter(0), "-", "_")
	var value any
	for _, o := range h.enumType.GetEnumConstants() {
		if equalsIgnoreCaseUTF16(enumName(o), s) {
			value = o
			break
		}
	}
	if value == nil {
		optName := params.GetParameter(-1)
		if h.Option.IsArgument() {
			optName = h.Option.String()
		}
		panic(newCmdLineException(h.Owner, spiMsgIllegalOperand, optName, s))
	}
	h.Setter.AddValue(value)
	return 1
}

func equalsIgnoreCaseUTF16(a, b string) bool {
	ua, ub := jutil.StringToUTF16(a), jutil.StringToUTF16(b)
	return len(ua) == len(ub) && jutil.UTF16RegionMatches(true, ua, 0, ub, 0, len(ua))
}

func (h *EnumOptionHandler) GetDefaultMetaVariable() *string {
	constants := h.enumType.GetEnumConstants()
	var buf strings.Builder
	buf.WriteByte('[')
	for i, c := range constants {
		if i > 0 {
			buf.WriteString(" | ")
		}
		buf.WriteString(jutil.ToJavaString(c))
	}
	buf.WriteByte(']')
	s := buf.String()
	return &s
}

func (h *EnumOptionHandler) GetMetaVariable(_ ResourceBundle) *string {
	return h.GetDefaultMetaVariable()
}

// MapOptionHandler parses "key=value" operands into a Map field.
type MapOptionHandler struct {
	OptionHandlerBase
}

func NewMapOptionHandler(parser *CmdLineParser, option OptionDef, setter Setter) OptionHandler {
	h := &MapOptionHandler{}
	h.InitOptionHandlerBase(h, parser, option, setter)
	if setter.AsFieldSetter() == nil {
		panic(jutil.NewIllegalArgumentExceptionMsg("MapOptionHandler can only work with fields"))
	}
	return h
}

func (h *MapOptionHandler) GetDefaultMetaVariable() *string { return nil }

func (h *MapOptionHandler) ParseArguments(params Parameters) int {
	fs := h.Setter.AsFieldSetter()
	var m Map
	if o := fs.GetValue(); o != nil {
		m = o.(Map)
	}
	if m == nil {
		m = jutil.NewJavaHashMap[string, *string]()
		fs.AddValue(m)
	}
	h.addToMap(params.GetParameter(0), m)
	return 1
}

func (h *MapOptionHandler) addToMap(argument string, m Map) {
	mapKey, value, ok := strings.Cut(argument, "=")
	if !ok {
		panic(newCmdLineException(h.Owner, spiMsgFormatErrorForMap))
	}
	if len(mapKey) == 0 {
		panic(newCmdLineException(h.Owner, spiMsgMapHasNoKey))
	}
	var mapValue *string
	if len(value) > 0 {
		mapValue = &value
	}
	m.Put(mapKey, mapValue)
}

// StringArrayOptionHandler greedily consumes non-option tokens (splitting each
// on literal spaces) into a multi-valued String setter.
type StringArrayOptionHandler struct {
	OptionHandlerBase
}

func NewStringArrayOptionHandler(parser *CmdLineParser, option OptionDef, setter Setter) OptionHandler {
	h := &StringArrayOptionHandler{}
	h.InitOptionHandlerBase(h, parser, option, setter)
	return h
}

func (h *StringArrayOptionHandler) GetDefaultMetaVariable() *string {
	s := defaultMetaStringArray
	return &s
}

func (h *StringArrayOptionHandler) ParseArguments(params Parameters) int {
	counter := 0
	for ; counter < params.Size(); counter++ {
		param := params.GetParameter(counter)
		if strings.HasPrefix(param, "-") {
			break
		}
		for _, p := range jutil.StringSplitLiteral(param, " ") {
			h.Setter.AddValue(p)
		}
	}
	return counter
}
