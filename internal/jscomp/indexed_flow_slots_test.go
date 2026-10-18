// Copyright 2026 The jsgo Authors.
// Licensed under the Apache License, Version 2.0.

package jscomp

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Run the original lattice suite unchanged against the alternative storage,
// including bottom joins, qualified slots, block scope trimming and long chains.
func TestIndexedFlowSlotsOriginalLatticeSuite(t *testing.T) {
	saved := indexedFlowSlotsEnabled
	indexedFlowSlotsEnabled = true
	defer func() { indexedFlowSlotsEnabled = saved }()
	RunJUnit4(t, func() JUnitTestCase { return NewLinkedFlowScopeTest() })
}

// Type-join/equality callbacks can resolve lazy aliases. Compare complete typed
// ASTs and ordered diagnostics, rather than checking generated JS alone.
func TestIndexedFlowSlotsCompilerDifferential(t *testing.T) {
	saved := indexedFlowSlotsEnabled
	defer func() { indexedFlowSlotsEnabled = saved }()
	fixtures := []string{
		`/** @typedef {(Later|string)} */ var Early;
   /** @param {Early} x @param {boolean} c */
   function f(x,c) { if(c) {x=1;} else if(typeof x==='string') {x=x.length;} return x; }
   /** @typedef {number} */ var Later; f('x',true);`,
		`/** @constructor */ function Link() { /** @type {?Link} */ this.next=null; }
   /** @param {?Link} node @param {boolean} c */
   function walk(node,c) { while(node && c) {if(node.next) {node=node.next;} else {break;}} return node; }`,
		`/** @param {?} x @param {boolean} c */
   function unknown(x,c) { if(c) {x=3;} else {x='x';} if(typeof x==='string') {return x.length;} return x; }`,
		`/** @param {(number|string)} Aa @param {(number|string)} BB @param {boolean} c */
   function collisions(Aa,BB,c) {if(c) {Aa=1;BB='x';} else {Aa='x';BB=2;} return Aa+BB;}`,
		`/** @param {{value:(number|string)}} x @param {boolean} c */
   function shadow(x,c) { if(c) {let x={value:1};x.value=2;} else {x.value='x';}
    try {if(c) {throw x;}} catch(e) {e.value=3;} return x.value; }`,
	}
	var large strings.Builder
	large.WriteString("/** @param {boolean} c */ function large(c) {\n")
	for i := 0; i < 280; i++ {
		fmt.Fprintf(&large, "var x%d = %d;\n", i, i)
	}
	large.WriteString("if(c) {\n")
	for i := 0; i < 280; i++ {
		fmt.Fprintf(&large, "x%d = 'x';\n", i)
	}
	large.WriteString("} return x0+x279; }")
	fixtures = append(fixtures, large.String())
	type snapshot struct {
		success          bool
		source, ast      string
		errors, warnings []string
	}
	compile := func(code string, native bool) snapshot {
		indexedFlowSlotsEnabled = native
		options := NewCompilerOptions()
		options.SetChecksOnly(true)
		options.SetCheckTypes(true)
		options.SetContinueAfterErrors(true)
		compiler := NewCompilerWithErrorManager(NewBlackHoleErrorManager())
		compiler.DisableThreads()
		result := compiler.Compile([]*SourceFile{SourceFile_FromCode("externs.js", CompilerTypeTestCase_DEFAULT_EXTERNS)}, []*SourceFile{SourceFile_FromCode("input.js", code)}, options)
		out := snapshot{success: result.Success, source: compiler.ToSource()}
		var types strings.Builder
		var visit func(*rhino.Node)
		visit = func(node *rhino.Node) {
			fmt.Fprintf(&types, "%s@%d:%d", node.GetToken(), node.GetLineno(), node.GetCharno())
			if typ := node.GetJSType(); typ != nil {
				types.WriteString(":" + typ.String())
			}
			types.WriteByte('\n')
			for child := node.GetFirstChild(); child != nil; child = child.GetNext() {
				visit(child)
			}
		}
		visit(compiler.GetJsRoot())
		out.ast = types.String()
		for _, err := range result.Errors {
			out.errors = append(out.errors, err.String())
		}
		for _, warning := range result.Warnings {
			out.warnings = append(out.warnings, warning.String())
		}
		return out
	}
	for i, code := range fixtures {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			want, got := compile(code, false), compile(code, true)
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("indexed slots changed typed compilation\ngot=%+v\nwant=%+v", got, want)
			}
		})
	}
}
