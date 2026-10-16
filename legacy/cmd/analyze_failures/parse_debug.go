//go:build ignore

package main

import (
	"fmt"
	"github.com/paullepoulpe/jsgo/legacy/internal/jscomp"
	"github.com/paullepoulpe/jsgo/legacy/internal/diag"
)

func main() {
	src := `/** @type {string} */ var a; a &&= 0;`
	ds, err := jscomp.CheckSourceWithOptions("t.js", src, jscomp.NewCompilerOptions())
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	for _, d := range ds {
		lvl := "WARNING"
		if d.Level == diag.Error { lvl = "ERROR" }
		fmt.Printf("%s:%d:%d: %s - [%s] %s\n", d.Path, d.Line, d.Col, lvl, d.Code, d.Message)
	}
	if len(ds) == 0 {
		fmt.Println("no diags")
	}
}
