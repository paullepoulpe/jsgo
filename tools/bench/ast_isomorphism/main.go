// Command ast_isomorphism verifies lockstep AST + lexical-scope alpha-isomorphism
// between two compiled JS files (e.g. compat vs. fast_mode output) using jsgo's
// parser and SyntacticScopeCreator.
package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type output struct {
	Isomorphic bool                 `json:"isomorphic"`
	Error      string               `json:"error,omitempty"`
	Stats      *AstIsomorphismStats `json:"stats,omitempty"`
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: %s <compat.js> <fast.js>\n", os.Args[0])
		os.Exit(2)
	}
	bA, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v\n", os.Args[1], err)
		os.Exit(2)
	}
	bB, err := os.ReadFile(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v\n", os.Args[2], err)
		os.Exit(2)
	}

	stats, isoErr := VerifyAstAlphaIsomorphism(string(bA), string(bB))
	res := output{
		Isomorphic: isoErr == nil,
		Stats:      stats,
	}
	if isoErr != nil {
		res.Error = isoErr.Error()
	}
	enc := json.NewEncoder(os.Stdout)
	_ = enc.Encode(res)
	if isoErr != nil {
		os.Exit(1)
	}
}
