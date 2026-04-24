// Port of com/google/javascript/jscomp/DiagnosticType.java
// @ closure-compiler v20260415

package diag

import "fmt"

// DiagnosticType is the type of a compile or analysis error. It carries a
// stable string key, a default CheckLevel, and a java.text.MessageFormat-style
// description template. Matches Java's DiagnosticType final class.
type DiagnosticType struct {
	// Key is the stable error identifier (e.g. "JSC_TYPE_MISMATCH").
	Key string
	// Format is a java.text.MessageFormat-style template where {0}, {1}, …
	// are replaced by arguments at emission time.
	Format string
	// Level is the default reporting level for this diagnostic.
	Level CheckLevel
}

// Error creates a DiagnosticType at CheckLevel ERROR.
func DiagnosticTypeError(name, descriptionFormat string) *DiagnosticType {
	return DiagnosticTypeMake(name, CheckLevelError, descriptionFormat)
}

// Warning creates a DiagnosticType at CheckLevel WARNING.
func DiagnosticTypeWarning(name, descriptionFormat string) *DiagnosticType {
	return DiagnosticTypeMake(name, CheckLevelWarning, descriptionFormat)
}

// Disabled creates a DiagnosticType at CheckLevel OFF.
func DiagnosticTypeDisabled(name, descriptionFormat string) *DiagnosticType {
	return DiagnosticTypeMake(name, CheckLevelOff, descriptionFormat)
}

// DiagnosticTypeMake is the general factory — mirrors Java's DiagnosticType.make.
func DiagnosticTypeMake(name string, level CheckLevel, descriptionFormat string) *DiagnosticType {
	return &DiagnosticType{Key: name, Level: level, Format: descriptionFormat}
}

// FormatMessage substitutes positional arguments into a java.text.MessageFormat
// template. Only the simple {N} positional form is supported (no choice/date
// sub-formats). Matches Platform.formatMessage used by Java's DiagnosticType.format.
func (dt *DiagnosticType) FormatMessage(arguments ...string) string {
	result := dt.Format
	for i, arg := range arguments {
		placeholder := fmt.Sprintf("{%d}", i)
		// Simple string replacement — sufficient for the {N} style Closure uses.
		result = replaceAll(result, placeholder, arg)
	}
	return result
}

// replaceAll replaces all occurrences of old with new_ in s.
func replaceAll(s, old, new_ string) string {
	if old == "" {
		return s
	}
	out := ""
	for {
		idx := indexOf(s, old)
		if idx < 0 {
			out += s
			break
		}
		out += s[:idx] + new_
		s = s[idx+len(old):]
	}
	return out
}

func indexOf(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func (dt *DiagnosticType) String() string {
	return dt.Key + ": " + dt.Format
}

// CheckLevel controls reporting level for a diagnostic.
// Mirrors com/google/javascript/jscomp/CheckLevel.java.
type CheckLevel int

const (
	CheckLevelOff     CheckLevel = iota // OFF: suppress entirely
	CheckLevelWarning                   // WARNING: emit as warning
	CheckLevelError                     // ERROR: emit as error
)

func (cl CheckLevel) IsOn() bool {
	return cl != CheckLevelOff
}

func (cl CheckLevel) String() string {
	switch cl {
	case CheckLevelError:
		return "ERROR"
	case CheckLevelWarning:
		return "WARNING"
	default:
		return "OFF"
	}
}
