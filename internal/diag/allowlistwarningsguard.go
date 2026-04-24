// Port of com/google/javascript/jscomp/AllowlistWarningsGuard.java
// @ closure-compiler v20260415

package diag

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
)

// lineNumberPattern strips trailing ":-?<N>" from allowlist lines.
var lineNumberPattern = regexp.MustCompile(`:-?\d+`)

// blazeGenfilesPattern strips Blaze genfiles prefixes from source names.
var blazeGenfilesPattern = regexp.MustCompile(`blaze-out/[^/]*/(bin|genfiles)/`)

// AllowlistWarningsGuard suppresses diagnostics whose formatted key appears in
// a pre-loaded allowlist. ERROR-level diagnostics are never suppressed.
//
// Matches Java's AllowlistWarningsGuard class.
//
// The allowlist format is one warning per line:
//
//	<file-name>:  <first-line-of-description>
//
// Lines starting with '#' are comments; blank lines are ignored.
// Line numbers embedded as ":<N>" or ":-<N>" are stripped before comparison
// so that the allowlist remains valid across small edits.
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
// Mirrors Java's AllowlistWarningsGuard.loadAllowlistedJsWarnings(Reader).
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
		// If the message matches the guard we use WARNING, so that it
		// - Shows up on stderr, and
		// - Gets caught by the AllowlistBuilder downstream in the pipeline.
		return CheckLevelWarning, true
	}
	return 0, false
}

// ContainWarning reports whether formattedWarning appears in the allowlist.
// Mirrors Java's AllowlistWarningsGuard.containWarning(String).
func (g *AllowlistWarningsGuard) ContainWarning(formattedWarning string) bool {
	return g.allowlist[formattedWarning]
}

// FormatWarning formats a Diagnostic into the allowlist key format
// "<file>:  <first-line-of-description>".
// Mirrors Java's AllowlistWarningsGuard.formatWarning(JSError, false).
func (g *AllowlistWarningsGuard) FormatWarning(d Diagnostic) string {
	return formatWarningImpl(d, false)
}

// FormatWarningWithMetadata includes line number and full multi-line message.
// Mirrors Java's AllowlistWarningsGuard.formatWarning(JSError, true).
func (g *AllowlistWarningsGuard) FormatWarningWithMetadata(d Diagnostic) string {
	return formatWarningImpl(d, true)
}

// formatWarningImpl is the common implementation for FormatWarning and
// FormatWarningWithMetadata. Mirrors Java's AllowlistWarningsGuard.formatWarning(JSError, boolean).
func formatWarningImpl(d Diagnostic, withMetaData bool) string {
	var sb strings.Builder
	src := normalizeSourceName(d.Path)
	sb.WriteString(src)
	sb.WriteByte(':')
	if withMetaData && d.Line > 0 {
		sb.WriteString(intToString(d.Line))
	}
	lines := strings.SplitN(d.Message, "\n", -1)
	sb.WriteString("  ")
	if len(lines) > 0 {
		sb.WriteString(lines[0])
	}
	if withMetaData {
		for i := 1; i < len(lines); i++ {
			sb.WriteString("\n# ")
			sb.WriteString(lines[i])
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

// normalizeSourceName strips Blaze genfiles prefixes.
// Mirrors Java's AllowlistWarningsGuard.normalizeSourceName(String).
func normalizeSourceName(sourceName string) string {
	if sourceName == "" {
		return sourceName
	}
	// "blaze-out/k8-fastbuild/genfiles/some/path/foo.js" → "some/path/foo.js"
	return blazeGenfilesPattern.ReplaceAllString(sourceName, "")
}

// intToString converts an integer to its decimal string representation
// without importing fmt (to keep the hot path allocation-free).
func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + intToString(-n)
	}
	digits := make([]byte, 0, 10)
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}

// GetFirstLine returns everything up to (but not including) the first newline.
// Mirrors Java's AllowlistWarningsGuard.getFirstLine(String).
func GetFirstLine(warning string) string {
	idx := strings.Index(warning, "\n")
	if idx > 0 {
		return warning[:idx]
	}
	return warning
}

// GetPriority returns the priority of this guard. A higher number means the
// guard is applied later. Mirrors Java's AllowlistWarningsGuard.getPriority().
// Returns 1 to match Java's WarningsGuard.Priority.SUPPRESS_BY_ALLOWLIST.getValue().
func (g *AllowlistWarningsGuard) GetPriority() int {
	return 1
}

// ---------------------------------------------------------------------------
// AllowlistBuilder — mirrors Java's AllowlistWarningsGuard.AllowlistBuilder
// inner class.
// ---------------------------------------------------------------------------

// AllowlistBuilder collects warnings and can write them out in the allowlist
// format understood by AllowlistWarningsGuard. Mirrors Java's
// AllowlistWarningsGuard.AllowlistBuilder inner class.
type AllowlistBuilder struct {
	warnings        []warningEntry
	productName     string
	generatorTarget string
	headerNote      string
	// guard is the owning AllowlistWarningsGuard, used for FormatWarning.
	guard *AllowlistWarningsGuard
}

type warningEntry struct {
	diag     Diagnostic
	diagType *DiagnosticType
}

// NewAllowlistBuilder returns a new AllowlistBuilder associated with the
// given guard.
func NewAllowlistBuilder(g *AllowlistWarningsGuard) *AllowlistBuilder {
	return &AllowlistBuilder{guard: g}
}

// SetProductName records the product name for inclusion in the header comment.
// Mirrors Java's AllowlistBuilder.setProductName(String).
func (b *AllowlistBuilder) SetProductName(name string) *AllowlistBuilder {
	b.productName = name
	return b
}

// SetGeneratorTarget records how to regenerate the allowlist.
// Mirrors Java's AllowlistBuilder.setGeneratorTarget(String).
func (b *AllowlistBuilder) SetGeneratorTarget(name string) *AllowlistBuilder {
	b.generatorTarget = name
	return b
}

// SetNote sets the header note included in the allowlist file.
// Mirrors Java's AllowlistBuilder.setNote(String).
func (b *AllowlistBuilder) SetNote(note string) *AllowlistBuilder {
	b.headerNote = note
	return b
}

// Report adds a diagnostic to the builder. ERROR-level diagnostics are
// ignored, matching Java's AllowlistBuilder.report(CheckLevel, JSError).
func (b *AllowlistBuilder) Report(level CheckLevel, d Diagnostic, dt *DiagnosticType) {
	if d.Level == Error {
		// ERROR-level diagnostics are ignored by AllowlistWarningsGuard.
		return
	}
	b.warnings = append(b.warnings, warningEntry{diag: d, diagType: dt})
}

// AppendAllowlist writes the collected warnings to w in allowlist format.
// Mirrors Java's AllowlistBuilder.appendAllowlist(PrintStream).
func (b *AllowlistBuilder) AppendAllowlist(w io.Writer) error {
	write := func(s string) error {
		_, err := fmt.Fprint(w, s)
		return err
	}

	if err := write("# This is a list of legacy warnings that have yet to be fixed.\n"); err != nil {
		return err
	}

	if b.productName != "" && len(b.warnings) > 0 {
		msg := "# Please find some time and fix at least one of them " +
			"and it will be the happiest day for " + b.productName + ".\n"
		if err := write(msg); err != nil {
			return err
		}
	}

	if b.generatorTarget != "" {
		if err := write("# When you fix any of these warnings, run " + b.generatorTarget + " task.\n"); err != nil {
			return err
		}
	}

	if b.headerNote != "" {
		lines := strings.Split(b.headerNote, "\n")
		for _, line := range lines {
			if err := write("# " + line + "\n"); err != nil {
				return err
			}
		}
	}

	// Group warnings by DiagnosticType key.
	type typeGroup struct {
		key      string
		format   string
		warnings []string
	}
	byType := make(map[string]*typeGroup)
	for _, e := range b.warnings {
		key := e.diag.Code
		if e.diagType != nil {
			key = e.diagType.Key
		}
		tg, ok := byType[key]
		if !ok {
			fmtStr := ""
			if e.diagType != nil {
				fmtStr = e.diagType.Format
			}
			tg = &typeGroup{key: key, format: fmtStr}
			byType[key] = tg
		}
		formatted := b.guard.FormatWarningWithMetadata(e.diag)
		tg.warnings = append(tg.warnings, formatted)
	}

	// Sort by key for deterministic output.
	keys := make([]string, 0, len(byType))
	for k := range byType {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		tg := byType[key]
		firstLine := GetFirstLine(tg.format)
		if err := write("\n# Warning " + tg.key + ": " + firstLine + "\n"); err != nil {
			return err
		}
		sort.Strings(tg.warnings)
		for _, warning := range tg.warnings {
			if err := write(warning); err != nil {
				return err
			}
		}
	}

	return nil
}

// WriteAllowlist writes the collected warnings to a file at path.
// Mirrors Java's AllowlistBuilder.writeAllowlist(File).
func (b *AllowlistBuilder) WriteAllowlist(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return b.AppendAllowlist(f)
}
