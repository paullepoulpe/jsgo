package pin

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	p, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(p.Version, "v") || !strings.HasSuffix(p.Jar, ".jar") {
		t.Errorf("Load() = %+v", p)
	}
	for _, s := range p.JarSHA256 {
		if len(s) != 64 {
			t.Errorf("CLOSURE_JAR_SHA256 %q is not a sha256 digest", s)
		}
	}
}

// The Go copies of the version (PORTING.md: the pinned release) must agree
// with the pin file.
func TestCopiesOfTheVersionAgree(t *testing.T) {
	p, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	for file, decl := range map[string]string{
		"internal/jscomp/command_line_runner_version.go": "const CommandLineRunnerVersion_COMPILER_VERSION = ",
		"tools/licenses/main.go":                         "const tag = ",
	} {
		b, err := os.ReadFile(filepath.Join("..", "..", "..", filepath.FromSlash(file)))
		if err != nil {
			t.Fatal(err)
		}
		if want := decl + strconv.Quote(p.Version); !strings.Contains(string(b), want+"\n") {
			t.Errorf("%s does not declare %s (%s)", file, want, Path)
		}
	}
}

func TestParse(t *testing.T) {
	p, err := Parse([]byte("# c\nCLOSURE_VERSION=v1\n\nCLOSURE_JAR=a.jar\nCLOSURE_JAR_SHA256=\"x y\"\n"))
	if err != nil || p.Version != "v1" || p.Jar != "a.jar" || len(p.JarSHA256) != 2 || p.JarSHA256[1] != "y" {
		t.Errorf("Parse = %+v, %v", p, err)
	}
	if _, err := Parse([]byte("CLOSURE_VERSION=v1\n")); err == nil {
		t.Error("Parse accepted a pin without a jar")
	}
}
