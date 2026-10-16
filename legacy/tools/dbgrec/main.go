package main

import (
	"fmt"
	"os"

	"github.com/paullepoulpe/jsgo/legacy/internal/jscomp"
)

func main() {
	f := os.Args[1]
	b, _ := os.ReadFile(f)
	opts := jscomp.NewCompilerOptions()
	opts.SetWarningLevel(jscomp.GroupStrictCheckTypes, jscomp.LevelWarning)
	d, _ := jscomp.CheckSourceWithOptions(f, string(b), opts)
	fmt.Printf("got %d diags\n", len(d))
	for _, x := range d {
		fmt.Printf("%d:%d %s %s\n", x.Line, x.Col, x.Code, x.Message)
	}
}
