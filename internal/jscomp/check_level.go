// Port of com/google/javascript/jscomp/CheckLevel.java @ closure-compiler v20260415

package jscomp

// CheckLevel controls checking levels of certain options. For all checks
// going forward, this should be used instead of booleans, so teams and
// individuals can control which checks are off, which produce only warnings,
// and which produce errors, without everyone having to agree.
type CheckLevel int32

const (
	CheckLevel_ERROR CheckLevel = iota
	CheckLevel_WARNING
	CheckLevel_OFF
)

var checkLevel_names = [...]string{"ERROR", "WARNING", "OFF"}

// CheckLevel_values is CheckLevel.values().
func CheckLevel_values() []CheckLevel {
	return []CheckLevel{CheckLevel_ERROR, CheckLevel_WARNING, CheckLevel_OFF}
}

// String is Enum.name().
func (l CheckLevel) String() string { return checkLevel_names[l] }

// Ordinal is Enum.ordinal().
func (l CheckLevel) Ordinal() int { return int(l) }

// CompareTo is Enum.compareTo(E): the ordinal difference.
func (l CheckLevel) CompareTo(o CheckLevel) int { return int(l) - int(o) }

// Java: CheckLevel.isOn()
func (l CheckLevel) IsOn() bool {
	return l != CheckLevel_OFF
}

// Go-only: language-adapter — a Java CheckLevel reference may be null; such
// values are *CheckLevel here (nil is null). Ptr returns the canonical
// pointer to a constant, so that pointers to the same level are identical,
// like the Java enum constants.
var checkLevel_ptrs = [...]CheckLevel{CheckLevel_ERROR, CheckLevel_WARNING, CheckLevel_OFF}

// Ptr returns the canonical non-null *CheckLevel for l.
func (l CheckLevel) Ptr() *CheckLevel { return &checkLevel_ptrs[l] }
