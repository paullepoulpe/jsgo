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

// Go-only: DiagnosticTypeError ports DiagnosticType.error(String, String) using the
// package-qualified naming convention; Java's static method is named just `error`.
func DiagnosticTypeError(name, descriptionFormat string) *DiagnosticType {
	return DiagnosticTypeMake(name, CheckLevelError, descriptionFormat)
}

// Go-only: DiagnosticTypeWarning ports DiagnosticType.warning(String, String) using the
// package-qualified naming convention; Java's static method is named just `warning`.
func DiagnosticTypeWarning(name, descriptionFormat string) *DiagnosticType {
	return DiagnosticTypeMake(name, CheckLevelWarning, descriptionFormat)
}

// Go-only: DiagnosticTypeDisabled ports DiagnosticType.disabled(String, String) using the
// package-qualified naming convention; Java's static method is named just `disabled`.
func DiagnosticTypeDisabled(name, descriptionFormat string) *DiagnosticType {
	return DiagnosticTypeMake(name, CheckLevelOff, descriptionFormat)
}

// Go-only: DiagnosticTypeMake ports DiagnosticType.make(String, CheckLevel, String) using the
// package-qualified naming convention; Java's static method is named just `make`.
func DiagnosticTypeMake(name string, level CheckLevel, descriptionFormat string) *DiagnosticType {
	return &DiagnosticType{Key: name, Level: level, Format: descriptionFormat}
}

// Go-only: FormatMessage ports DiagnosticType.format(String...) under a different name
// to avoid shadowing the Format field; Java's method is named just `format`.
func (dt *DiagnosticType) FormatMessage(arguments ...string) string {
	result := dt.Format
	for i, arg := range arguments {
		placeholder := fmt.Sprintf("{%d}", i)
		// Simple string replacement — sufficient for the {N} style Closure uses.
		result = replaceAll(result, placeholder, arg)
	}
	return result
}

// Go-only: replaceAll is a string-replacement helper; Java uses java.text.MessageFormat.
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

// Go-only: indexOf is a string-search helper used by replaceAll; Java uses String.indexOf.
func indexOf(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

// Equals mirrors DiagnosticType.equals(Object) — two DiagnosticTypes are equal iff their keys match.
func (dt *DiagnosticType) Equals(other *DiagnosticType) bool {
	if other == nil {
		return false
	}
	return dt.Key == other.Key
}

// HashCode mirrors DiagnosticType.hashCode() — returns a hash of the key string.
func (dt *DiagnosticType) HashCode() int {
	h := 0
	for _, c := range dt.Key {
		h = 31*h + int(c)
	}
	return h
}

// CompareTo mirrors DiagnosticType.compareTo(DiagnosticType) — lexicographic key comparison.
func (dt *DiagnosticType) CompareTo(other *DiagnosticType) int {
	if dt.Key < other.Key {
		return -1
	}
	if dt.Key > other.Key {
		return 1
	}
	return 0
}

// String mirrors DiagnosticType.toString().
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

// Go-only: IsOn is a convenience predicate absent from the Java CheckLevel enum.
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
