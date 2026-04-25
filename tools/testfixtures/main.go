package main

import (
	"fmt"
	"os"
	"github.com/paullepoulpe/jsgo/internal/jscomp"
)

func main() {
	files := []string{
		"testdata/fixtures/extracted/testsetprop5.js",
		"testdata/fixtures/extracted/testcall8.js",
		"testdata/fixtures/extracted/testextendedinterfacepropertiescompatibility9.js",
		"testdata/fixtures/extracted/testinheritancecheck17.js",
		"testdata/fixtures/extracted/testbadextends3.js",
		"testdata/fixtures/extracted/testwarndatapropertyoninterface4.js",
		"testdata/fixtures/extracted/testinterfacenonemptyfunction.js",
	}
	for _, f := range files {
		b, _ := os.ReadFile(f)
		diags, _ := jscomp.CheckSource(f, string(b))
		fmt.Printf("=== %s ===\n", f)
		for _, d := range diags {
			fmt.Printf("  %s\n", d.FormatJavaLike())
		}
		if len(diags) == 0 {
			fmt.Printf("  (no diagnostics)\n")
		}
	}
}
