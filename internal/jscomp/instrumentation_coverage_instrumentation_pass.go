// Port of com/google/javascript/jscomp/instrumentation/CoverageInstrumentationPass.java @ closure-compiler v20261021

/*
 * Copyright 2009 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CoverageInstrumentationPass implements the instrumentation pass over the
// AST (returned by JSCompiler).
//
// Java: public class CoverageInstrumentationPass implements CompilerPass
type CoverageInstrumentationPass struct {
	compiler                  AbstractCompiler
	lineInstrumentationData   map[string]*FileInstrumentationData
	branchInstrumentationData *jutil.LinkedHashMap[string, *FileInstrumentationData]
	reach                     CoverageInstrumentationPass_CoverageReach
	instrumentOption          CompilerOptions_InstrumentOption

	productionInstrumentationArrayName string
}

// Java: CoverageInstrumentationPass.JS_INSTRUMENTATION_OBJECT_NAME
const CoverageInstrumentationPass_JS_INSTRUMENTATION_OBJECT_NAME = "__jscov"

// CoverageInstrumentationPass_CoverageReach configures which statements in
// the AST should be instrumented.
//
// Java: public enum CoverageInstrumentationPass.CoverageReach
type CoverageInstrumentationPass_CoverageReach int32

const (
	CoverageInstrumentationPass_CoverageReach_ALL         CoverageInstrumentationPass_CoverageReach = iota // Instrument all statements.
	CoverageInstrumentationPass_CoverageReach_CONDITIONAL                                                  // Do not instrument global statements.
)

var coverageInstrumentationPass_CoverageReach_names = [...]string{"ALL", "CONDITIONAL"}

// Java: CoverageInstrumentationPass.CoverageReach.name()
func (r CoverageInstrumentationPass_CoverageReach) String() string {
	return coverageInstrumentationPass_CoverageReach_names[r]
}

// Java: CoverageInstrumentationPass.CoverageReach.ordinal()
func (r CoverageInstrumentationPass_CoverageReach) Ordinal() int { return int(r) }

// NewCoverageInstrumentationPassWithInstrumentOption: compiler is the
// compiler which generates the AST.
//
// Java: CoverageInstrumentationPass(AbstractCompiler, CoverageInstrumentationPass.CoverageReach, CompilerOptions.InstrumentOption, String)
func NewCoverageInstrumentationPassWithInstrumentOption(
	compiler AbstractCompiler,
	reach CoverageInstrumentationPass_CoverageReach,
	instrumentOption CompilerOptions_InstrumentOption,
	productionInstrumentationArrayName string) *CoverageInstrumentationPass {
	return &CoverageInstrumentationPass{
		compiler:                           compiler,
		reach:                              reach,
		instrumentOption:                   instrumentOption,
		productionInstrumentationArrayName: productionInstrumentationArrayName,
	}
}

// Deprecated: use NewCoverageInstrumentationPassWithInstrumentOption(compiler,
// reach, CompilerOptions_InstrumentOption_LINE_ONLY, "").
//
// Java: CoverageInstrumentationPass(AbstractCompiler, CoverageInstrumentationPass.CoverageReach)
func NewCoverageInstrumentationPass(compiler AbstractCompiler, reach CoverageInstrumentationPass_CoverageReach) *CoverageInstrumentationPass {
	return NewCoverageInstrumentationPassWithInstrumentOption(compiler, reach, CompilerOptions_InstrumentOption_LINE_ONLY, "")
}

// AddHeaderCode creates the js code to be added to source. This code declares
// and initializes the variables required for collection of coverage data.
//
// Java: CoverageInstrumentationPass.addHeaderCode(Node)
func (p *CoverageInstrumentationPass) AddHeaderCode(script rhino.Node) {
	script.AddChildToFront(p.CreateConditionalObjectDecl(CoverageInstrumentationPass_JS_INSTRUMENTATION_OBJECT_NAME, script))

	// Make subsequent usages of "window" and "window.top" work in a Web Worker context.
	script.AddChildToFront(
		p.compiler.
			ParseSyntheticCode(
				"coverage_instrumentation_header",
				"(function(self) { if (!self.window) { self.window = self; self.window.top = self;"+
					" }})(typeof self !== 'undefined' ? self : globalThis)").
			RemoveFirstChild().
			SrcrefTreeIfMissing(script))
}

// Java: CoverageInstrumentationPass.addProductionHeaderCode(Node, String)
func (p *CoverageInstrumentationPass) AddProductionHeaderCode(script rhino.Node, arrayName string) {
	arrayLit := rhino.IR_Arraylit()
	name := rhino.IR_Name(arrayName)
	varNode := rhino.IR_VarWithValue(name, arrayLit)
	script.AddChildToFront(varNode.SrcrefTreeIfMissing(script))
}

// Java: CoverageInstrumentationPass.checkIfArrayNameExternDeclared(Node, String)
func (p *CoverageInstrumentationPass) CheckIfArrayNameExternDeclared(externsNode rhino.Node, arrayName string) {
	if !NodeUtil_CollectExternVariableNames(p.compiler, externsNode).Contains(arrayName) {
		panic(jutil.NewAssertionError(
			"The array name passed to --production_instrumentation_array_name was not " +
				"declared as an extern. This will result in undefined behaviour"))
	}
}

// Java: CoverageInstrumentationPass.process(Node, Node)
func (p *CoverageInstrumentationPass) Process(externsNode rhino.Node, rootNode rhino.Node) {
	if rootNode.HasChildren() {
		switch p.instrumentOption {
		case CompilerOptions_InstrumentOption_BRANCH_ONLY:
			if p.branchInstrumentationData == nil {
				p.branchInstrumentationData = jutil.NewLinkedHashMap[string, *FileInstrumentationData]()
			}
			NodeTraversal_Traverse(
				p.compiler,
				rootNode,
				NewBranchCoverageInstrumentationCallback(p.compiler, p.branchInstrumentationData))
		case CompilerOptions_InstrumentOption_PRODUCTION:
			productionCoverageInstrumentationCallback := NewProductionCoverageInstrumentationCallback(
				p.compiler, p.productionInstrumentationArrayName)

			NodeTraversal_Traverse(p.compiler, rootNode, productionCoverageInstrumentationCallback)

			p.CheckIfArrayNameExternDeclared(externsNode, p.productionInstrumentationArrayName)

			instrumentationMapping := productionCoverageInstrumentationCallback.GetInstrumentationMapping()
			p.compiler.SetInstrumentationMapping(instrumentationMapping)
		default:
			if p.lineInstrumentationData == nil {
				p.lineInstrumentationData = make(map[string]*FileInstrumentationData)
			}
			NodeTraversal_Traverse(
				p.compiler, rootNode, NewCoverageInstrumentationCallback(p.lineInstrumentationData, p.reach))
		}
		firstScript := rootNode.GetFirstChild()
		jutil.CheckState(firstScript.IsScript())
		// If any passes run after we need to preserve the MODULE_BODY structure of scripts - we can't
		// just add to a script if it is a module.
		if firstScript.HasChildren() && firstScript.GetFirstChild().IsModuleBody() {
			firstScript = firstScript.GetFirstChild()
		}

		if p.instrumentOption == CompilerOptions_InstrumentOption_PRODUCTION {
			p.AddProductionHeaderCode(firstScript, p.productionInstrumentationArrayName)
		} else {
			p.AddHeaderCode(firstScript)
		}
	}
}

// Java: CoverageInstrumentationPass.createConditionalObjectDecl(String, Node)
func (p *CoverageInstrumentationPass) CreateConditionalObjectDecl(name string, srcref rhino.Node) rhino.Node {
	// Make sure to quote properties so they are not renamed.
	var jscovData rhino.Node
	switch p.instrumentOption {
	case CompilerOptions_InstrumentOption_BRANCH_ONLY:
		jscovData = rhino.IR_Objectlit(
			rhino.IR_QuotedStringKey("fileNames", rhino.IR_Arraylit()),
			rhino.IR_QuotedStringKey("branchPresent", rhino.IR_Arraylit()),
			rhino.IR_QuotedStringKey("branchesInLine", rhino.IR_Arraylit()),
			rhino.IR_QuotedStringKey("branchesTaken", rhino.IR_Arraylit()))
	case CompilerOptions_InstrumentOption_LINE_ONLY:
		jscovData = rhino.IR_Objectlit(
			rhino.IR_QuotedStringKey("fileNames", rhino.IR_Arraylit()),
			rhino.IR_QuotedStringKey("instrumentedLines", rhino.IR_Arraylit()),
			rhino.IR_QuotedStringKey("executedLines", rhino.IR_Arraylit()))
	default:
		panic(jutil.NewAssertionError("Unexpected option: " + p.instrumentOption.String()))
	}

	// Add the __jscov var to the window as a quoted key so it can be found even if property
	// renaming is enabled.
	var_ := rhino.IR_VarWithValue(
		rhino.IR_Name(name),
		rhino.IR_Or(
			rhino.IR_Getelem(rhino.IR_Getprop(rhino.IR_Name("window"), "top"), rhino.IR_String("__jscov")),
			rhino.IR_Assign(
				rhino.IR_Getelem(rhino.IR_Getprop(rhino.IR_Name("window"), "top"), rhino.IR_String("__jscov")),
				jscovData)))
	return var_.SrcrefTreeIfMissing(srcref)
}
