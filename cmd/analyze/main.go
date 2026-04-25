package main

import (
	"fmt"
	"os"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
)

func main() {
	path := "testdata/fixtures/extracted/testmodulereferencenotallowed.js"
	src, _ := os.ReadFile(path)
	diags, _ := jscomp.CheckSourceWithOptions(path, string(src), jscomp.NewCompilerOptions())
	for _, d := range diags {
		fmt.Printf("%d:%d [%s] %s\n", d.Line, d.Col, d.Code, d.Message)
	}
}
