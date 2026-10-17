package jscomp

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Go-only: MapDiagnosticToSuppression has no Java test. These compare its
// output with what MapDiagnosticToSuppression.main prints from the pinned
// jar (stdout and stderr interleaved; the --output goldens in
// testdata/map_diagnostic_to_suppression/ are its stdout).

func runMapDiagnosticToSuppression(args ...string) string {
	var buf bytes.Buffer
	ps := jutil.NewPrintStream(&buf)
	out, err := jutil.System_Out(), jutil.System_Err()
	jutil.System_SetOut(ps)
	jutil.System_SetErr(ps)
	defer func() {
		jutil.System_SetOut(out)
		jutil.System_SetErr(err)
	}()
	MapDiagnosticToSuppression_Main(args)
	return buf.String()
}

func TestMapDiagnosticToSuppressionArguments(t *testing.T) {
	const help = "Usage: MapDiagnosticToSuppression --output={json,md}\n"
	const cle = "Could not parse arguments: com.google.javascript.jscomp.jarjar.org.kohsuke.args4j.CmdLineException: "
	for _, c := range []struct {
		args []string
		want string
	}{
		{[]string{"--bogus"}, cle + "\"--bogus\" is not a valid option\n" + help},
		{[]string{"--output"}, cle + "Option \"--output\" takes an operand\n" + help},
		{[]string{"--output", "xml"}, cle + "\"xml\" is not a valid value for \"--output\"\n" + help},
		{[]string{"--help"}, help},
		{[]string{"extra"}, cle + "No argument is allowed: extra\n" + help},
		{[]string{"--output=md", "--help"}, help},
	} {
		if got := runMapDiagnosticToSuppression(c.args...); got != c.want {
			t.Errorf("MapDiagnosticToSuppression %s:\n got %q\nwant %q", strings.Join(c.args, " "), got, c.want)
		}
	}
}

func TestMapDiagnosticToSuppressionOutput(t *testing.T) {
	for _, format := range []string{"json", "md"} {
		want, err := os.ReadFile("testdata/map_diagnostic_to_suppression/" + format + ".txt")
		if err != nil {
			t.Fatal(err)
		}
		if got := runMapDiagnosticToSuppression("--output", format); got != string(want) {
			t.Errorf("--output %s differs from the jar's output:\n got %d bytes\nwant %d bytes", format, len(got), len(want))
		}
	}
}
