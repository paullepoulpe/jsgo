//go:build ignore

package main

import (
	"fmt"
	"github.com/paullepoulpe/jsgo/legacy/internal/jscomp"
	"github.com/paullepoulpe/jsgo/legacy/internal/diag"
)

func main() {
	ds, _ := jscomp.CheckSourceWithOptions("test.js", "(function(){}) == 'x'", jscomp.NewCompilerOptions())
	for _, d := range ds {
		lvl := "WARNING"
		if d.Level == diag.Error {
			lvl = "ERROR"
		}
		fmt.Printf("%s:%d:%d: %s - [%s] %s\n", d.Path, d.Line, d.Col, lvl, d.Code, d.Message)
	}
}
