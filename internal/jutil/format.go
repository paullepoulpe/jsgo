package jutil

import (
	"strconv"
	"strings"
)

// Go-only: language-adapter — java.util.Formatter (String.format) for the
// conversions Closure uses: %s %S %b %d %c %x %X %o %% %n, with the '-' and
// '0' flags, a width, a precision (for %s) and explicit argument indices
// ("%1$s"). Formatting is locale-independent (Closure formats with the root
// locale in practice: ASCII digits, no grouping).

// illegalFormatException is java.util.IllegalFormatException (and its
// subclasses MissingFormatArgumentException, UnknownFormatConversionException).
type illegalFormatException struct{ IllegalArgumentException }

func newIllegalFormatException(className, message string) *illegalFormatException {
	e := &illegalFormatException{}
	InitThrowable(e, className, msgPtr(message), nil)
	return e
}

// StringFormat is String.format(String format, Object... args).
func StringFormat(format string, args ...any) string {
	var sb strings.Builder
	ordinary := 0 // index of the next ordinary (non-indexed) argument
	lastIndex := -1
	i := 0
	for i < len(format) {
		c := format[i]
		if c != '%' {
			sb.WriteByte(c)
			i++
			continue
		}
		start := i
		i++
		if i >= len(format) {
			panic(newIllegalFormatException("java.util.UnknownFormatConversionException", "Conversion = '%'"))
		}
		// [argument_index$]
		argIndex := -1
		j := i
		for j < len(format) && format[j] >= '0' && format[j] <= '9' {
			j++
		}
		if j < len(format) && format[j] == '$' && j > i {
			n, _ := strconv.Atoi(format[i:j])
			argIndex = n - 1
			i = j + 1
		}
		// [flags]
		leftJustify, zeroPad, relative := false, false, false
		for i < len(format) {
			switch format[i] {
			case '-':
				leftJustify = true
			case '0':
				zeroPad = true
			case '<':
				relative = true
			default:
				goto flagsDone
			}
			i++
		}
	flagsDone:
		// [width]
		width := -1
		j = i
		for j < len(format) && format[j] >= '0' && format[j] <= '9' {
			j++
		}
		if j > i {
			width, _ = strconv.Atoi(format[i:j])
			i = j
		}
		// [.precision]
		precision := -1
		if i < len(format) && format[i] == '.' {
			j = i + 1
			for j < len(format) && format[j] >= '0' && format[j] <= '9' {
				j++
			}
			precision, _ = strconv.Atoi(format[i+1 : j])
			i = j
		}
		if i >= len(format) {
			panic(newIllegalFormatException("java.util.UnknownFormatConversionException",
				"Conversion = '"+format[start:]+"'"))
		}
		conv := format[i]
		i++
		spec := format[start:i]
		if conv == '%' {
			sb.WriteString(pad("%", width, leftJustify, false))
			continue
		}
		if conv == 'n' {
			sb.WriteByte('\n')
			continue
		}
		var arg any
		switch {
		case relative:
			if lastIndex < 0 || lastIndex >= len(args) {
				panic(newIllegalFormatException("java.util.MissingFormatArgumentException", "Format specifier '"+spec+"'"))
			}
			arg = args[lastIndex]
		case argIndex >= 0:
			if argIndex >= len(args) {
				panic(newIllegalFormatException("java.util.MissingFormatArgumentException", "Format specifier '"+spec+"'"))
			}
			arg = args[argIndex]
			lastIndex = argIndex
		default:
			if ordinary >= len(args) {
				panic(newIllegalFormatException("java.util.MissingFormatArgumentException", "Format specifier '"+spec+"'"))
			}
			arg = args[ordinary]
			lastIndex = ordinary
			ordinary++
		}
		var s string
		numeric := false
		switch conv {
		case 's', 'S':
			s = ToJavaString(arg)
			if precision >= 0 {
				u := StringToUTF16(s)
				if precision < len(u) {
					s = StringFromUTF16(u[:precision])
				}
			}
			if conv == 'S' {
				s = strings.ToUpper(s)
			}
		case 'b', 'B':
			switch v := arg.(type) {
			case bool:
				s = strconv.FormatBool(v)
			default:
				s = strconv.FormatBool(!IsNil(arg))
			}
			if conv == 'B' {
				s = strings.ToUpper(s)
			}
		case 'c', 'C':
			switch v := arg.(type) {
			case uint16:
				s = StringFromChar(v)
			case int:
				s = string(rune(v))
			case int32:
				s = string(rune(v))
			default:
				if IsNil(arg) {
					s = "null"
				} else {
					panic(newIllegalFormatException("java.util.IllegalFormatConversionException", "c != "+javaClassNameOf(arg)))
				}
			}
			if conv == 'C' {
				s = strings.ToUpper(s)
			}
		case 'd':
			numeric = true
			v, ok := integerValue(arg)
			if !ok {
				if IsNil(arg) {
					s, numeric = "null", false
					break
				}
				panic(newIllegalFormatException("java.util.IllegalFormatConversionException", "d != "+javaClassNameOf(arg)))
			}
			s = strconv.FormatInt(v, 10)
		case 'x', 'X', 'o':
			numeric = true
			v, bits, ok := integerBits(arg)
			if !ok {
				if IsNil(arg) {
					s, numeric = "null", false
					break
				}
				panic(newIllegalFormatException("java.util.IllegalFormatConversionException",
					string(conv)+" != "+javaClassNameOf(arg)))
			}
			base := 16
			if conv == 'o' {
				base = 8
			}
			s = strconv.FormatUint(v&bits, base)
			if conv == 'X' {
				s = strings.ToUpper(s)
			}
		default:
			panic(newIllegalFormatException("java.util.UnknownFormatConversionException",
				"Conversion = '"+string(conv)+"'"))
		}
		sb.WriteString(pad(s, width, leftJustify, zeroPad && numeric))
	}
	return sb.String()
}

func pad(s string, width int, leftJustify, zeroPad bool) string {
	n := StringLength(s)
	if width <= n {
		return s
	}
	fill := width - n
	if zeroPad {
		if strings.HasPrefix(s, "-") {
			return "-" + strings.Repeat("0", fill) + s[1:]
		}
		return strings.Repeat("0", fill) + s
	}
	if leftJustify {
		return s + strings.Repeat(" ", fill)
	}
	return strings.Repeat(" ", fill) + s
}

func integerValue(arg any) (int64, bool) {
	switch v := arg.(type) {
	case int:
		return int64(v), true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	}
	return 0, false
}

// integerBits returns the value and the two's-complement mask of the Java
// type it stands for (Go int is Java int).
func integerBits(arg any) (uint64, uint64, bool) {
	switch v := arg.(type) {
	case int:
		return uint64(int64(v)), 0xFFFFFFFF, true
	case int8:
		return uint64(int64(v)), 0xFF, true
	case int16:
		return uint64(int64(v)), 0xFFFF, true
	case int32:
		return uint64(int64(v)), 0xFFFFFFFF, true
	case int64:
		return uint64(v), 0xFFFFFFFFFFFFFFFF, true
	}
	return 0, 0, false
}

func javaClassNameOf(arg any) string {
	switch arg.(type) {
	case string:
		return "java.lang.String"
	case uint16:
		return "java.lang.Character"
	case bool:
		return "java.lang.Boolean"
	case float64:
		return "java.lang.Double"
	case int, int32:
		return "java.lang.Integer"
	case int64:
		return "java.lang.Long"
	}
	return "java.lang.Object"
}
