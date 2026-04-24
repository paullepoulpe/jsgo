// Port of com/google/javascript/jscomp/AllowlistWarningsGuard.java
// @ closure-compiler v20260415

package diag

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
)

// lineNumberPattern strips trailing ":-?<N>" from allowlist lines.
var lineNumberPattern = regexp.MustCompile(`:-?\d+`)

// AllowlistWarningsGuard suppresses diagnostics whose formatted key appears in
// a pre-loaded allowlist. ERROR-level diagnostics are never suppressed.
//
// Matches Java's AllowlistWarningsGuard class.
type AllowlistWarningsGuard struct {
	allowlist map[string]bool
}

// NewAllowlistWarningsGuard creates a guard from a set of pre-normalized
// allowlist strings. Mirrors Java's AllowlistWarningsGuard(Set<String>).
func NewAllowlistWarningsGuard(allowlist map[string]bool) *AllowlistWarningsGuard {
	return &AllowlistWarningsGuard{allowlist: allowlist}
}

// NewAllowlistWarningsGuardEmpty creates a guard with an empty allowlist.
// Mirrors Java's AllowlistWarningsGuard() no-arg constructor.
func NewAllowlistWarningsGuardEmpty() *AllowlistWarningsGuard {
	return &AllowlistWarningsGuard{allowlist: map[string]bool{}}
}

// AllowlistWarningsGuardFromFile loads a guard from a text file.
// Mirrors Java's AllowlistWarningsGuard.fromFile(File).
func AllowlistWarningsGuardFromFile(path string) (*AllowlistWarningsGuard, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	lines, err := loadAllowlistedJsWarnings(f)
	if err != nil {
		return nil, err
	}
	return NewAllowlistWarningsGuard(NormalizeAllowlist(lines)), nil
}

// NormalizeAllowlist strips blank lines, comment lines, and line numbers from
// a set of raw allowlist strings.
// Mirrors Java's AllowlistWarningsGuard.normalizeAllowlist.
func NormalizeAllowlist(lines []string) map[string]bool {
	result := make(map[string]bool, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		result[lineNumberPattern.ReplaceAllString(trimmed, ":")] = true
	}
	return result
}

// loadAllowlistedJsWarnings reads all non-empty lines from r.
func loadAllowlistedJsWarnings(r io.Reader) ([]string, error) {
	var lines []string
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines, sc.Err()
}

// Level returns the overridden CheckLevel for the given Diagnostic, or the
// zero value (CheckLevelOff == 0) if the guard does not apply.
//
// ERROR-level diagnostics are never suppressed by the allowlist.
// Mirrors Java's AllowlistWarningsGuard.level(JSError).
func (g *AllowlistWarningsGuard) Level(d Diagnostic) (CheckLevel, bool) {
	if d.Level == Error {
		return 0, false
	}
	if len(g.allowlist) > 0 && g.ContainWarning(g.FormatWarning(d)) {
		return CheckLevelWarning, true
	}
	return 0, false
}

// ContainWarning reports whether formattedWarning appears in the allowlist.
// Mirrors Java's AllowlistWarningsGuard.containWarning.
func (g *AllowlistWarningsGuard) ContainWarning(formattedWarning string) bool {
	return g.allowlist[formattedWarning]
}

// FormatWarning formats a Diagnostic into the allowlist key format
// "<file>:  <first-line-of-description>".
// Mirrors Java's AllowlistWarningsGuard.formatWarning(JSError, false).
func (g *AllowlistWarningsGuard) FormatWarning(d Diagnostic) string {
	src := normalizeSourceName(d.Path)
	lines := strings.SplitN(d.Message, "\n", 2)
	firstLine := ""
	if len(lines) > 0 {
		firstLine = lines[0]
	}
	return src + ":  " + firstLine
}

// FormatWarningWithMetadata includes line number and full multi-line message.
// Mirrors Java's AllowlistWarningsGuard.formatWarning(JSError, true).
func (g *AllowlistWarningsGuard) FormatWarningWithMetadata(d Diagnostic) string {
	var sb strings.Builder
	src := normalizeSourceName(d.Path)
	sb.WriteString(src)
	sb.WriteByte(':')
	// include line number
	sb.WriteString(strings.ReplaceAll("", "", "")) // placeholder — no fmt import
	lines := strings.Split(d.Message, "\n")
	sb.WriteString("  ")
	sb.WriteString(lines[0])
	for i := 1; i < len(lines); i++ {
		sb.WriteString("\n# ")
		sb.WriteString(lines[i])
	}
	sb.WriteByte('\n')
	return sb.String()
}

// normalizeSourceName strips Blaze genfiles prefixes.
// Mirrors Java's AllowlistWarningsGuard.normalizeSourceName.
func normalizeSourceName(sourceName string) string {
	if sourceName == "" {
		return sourceName
	}
	// "blaze-out/k8-fastbuild/genfiles/some/path/foo.js" → "some/path/foo.js"
	re := regexp.MustCompile(`blaze-out/[^/]*/(bin|genfiles)/`)
	return re.ReplaceAllString(sourceName, "")
}

// GetFirstLine returns everything up to (but not including) the first newline.
// Mirrors Java's AllowlistWarningsGuard.getFirstLine.
func GetFirstLine(warning string) string {
	idx := strings.Index(warning, "\n")
	if idx > 0 {
		return warning[:idx]
	}
	return warning
}
