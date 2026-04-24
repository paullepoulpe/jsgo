package main

import (
	"fmt"
	"testing"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func printTree(n *rhino.Node, depth int) {
	if n == nil { return }
	indent := ""
	for i := 0; i < depth; i++ { indent += "  " }
	fmt.Printf("%s[%d] %v %q\n", indent, depth, n.Token(), n.String_())
	for _, c := range n.Children() { printTree(c, depth+1) }
}

func TestDebugParse(t *testing.T) {
	src := "/** @type {string|undefined} */ var a; a &&= 0;"
	ast, err := parser.Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	printTree(ast, 0)
}
