//go:build ignore

package main

import (
	"fmt"
	"github.com/paullepoulpe/jsgo/legacy/internal/jscomp"
	"github.com/paullepoulpe/jsgo/legacy/internal/rhino"
)

func main() {
	src := `function f(/** @type {string} */ x) { x ??= 5; }`
	root, err := jscomp.ParseForDebug("test.js", src)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	var walk func(*rhino.Node, string)
	walk = func(n *rhino.Node, indent string) {
		if n == nil { return }
		fmt.Printf("%s%v line=%d col=%d\n", indent, n.Token(), n.Lineno(), n.Charno())
		for _, c := range n.Children() {
			walk(c, indent+"  ")
		}
	}
	walk(root, "")
}
