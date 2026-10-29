// Port of com/google/javascript/jscomp/StrictModeCheck.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass uses the typed AST view (internal/ast): Visit hands
// each check the kind it reads, and the checks read by role (the target of
// an assignment, the operand of a delete, the object and property of a
// getprop). Two parts stay on rhino: the duplicate member scan walks the
// members of an object literal or a class body backwards, two different
// kinds of lists, and IsDeclaration reads a name's parent of any kind.

const (
	strictModeGetterBit       uint8 = 1 << 0
	strictModeSetterBit       uint8 = 1 << 1
	strictModeStaticGetterBit uint8 = 1 << 2
	strictModeStaticSetterBit uint8 = 1 << 3
)

type strictModeMemberEntry struct {
	id   uint32
	mask uint8
}

var (
	strictModeArgumentsID = rhino.IR_Name("arguments").GetStringID()
	strictModeEvalID      = rhino.IR_Name("eval").GetStringID()
	strictModeCalleeID    = rhino.IR_Name("callee").GetStringID()
	strictModeCallerID    = rhino.IR_Name("caller").GetStringID()
)

// StrictModeCheck checks that the code obeys the static restrictions of
// strict mode:
//
//  1. No use of "with".
//  2. No deleting variables, functions, or arguments.
//  3. No re-declarations or assignments of "eval" or arguments.
//  4. No use of arguments.callee
//  5. No use of arguments.caller
//  6. Class: Always under strict mode
//  7. In addition, no duplicate class method names
type StrictModeCheck struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler     AbstractCompiler
	defaultLevel CheckLevel
	memberFlags  map[uint32]uint8
}

// Java: StrictModeCheck.USE_OF_WITH
var StrictModeCheck_USE_OF_WITH = DiagnosticType_Error(
	"JSC_USE_OF_WITH", "The 'with' statement cannot be used in strict mode.")

// Java: StrictModeCheck.EVAL_DECLARATION
var StrictModeCheck_EVAL_DECLARATION = DiagnosticType_Error("JSC_EVAL_DECLARATION", "\"eval\" cannot be redeclared in strict mode")

// Java: StrictModeCheck.EVAL_ASSIGNMENT
var StrictModeCheck_EVAL_ASSIGNMENT = DiagnosticType_Error(
	"JSC_EVAL_ASSIGNMENT", "the \"eval\" object cannot be reassigned in strict mode")

// Java: StrictModeCheck.ARGUMENTS_DECLARATION
var StrictModeCheck_ARGUMENTS_DECLARATION = DiagnosticType_Error(
	"JSC_ARGUMENTS_DECLARATION", "\"arguments\" cannot be redeclared in strict mode")

// Java: StrictModeCheck.ARGUMENTS_ASSIGNMENT
var StrictModeCheck_ARGUMENTS_ASSIGNMENT = DiagnosticType_Error(
	"JSC_ARGUMENTS_ASSIGNMENT",
	"the \"arguments\" object cannot be reassigned in strict mode")

// Java: StrictModeCheck.ARGUMENTS_CALLEE_FORBIDDEN
var StrictModeCheck_ARGUMENTS_CALLEE_FORBIDDEN = DiagnosticType_Error(
	"JSC_ARGUMENTS_CALLEE_FORBIDDEN", "\"arguments.callee\" cannot be used in strict mode")

// Java: StrictModeCheck.ARGUMENTS_CALLER_FORBIDDEN
var StrictModeCheck_ARGUMENTS_CALLER_FORBIDDEN = DiagnosticType_Error(
	"JSC_ARGUMENTS_CALLER_FORBIDDEN", "\"arguments.caller\" cannot be used in strict mode")

// Java: StrictModeCheck.FUNCTION_CALLER_FORBIDDEN
var StrictModeCheck_FUNCTION_CALLER_FORBIDDEN = DiagnosticType_Error(
	"JSC_FUNCTION_CALLER_FORBIDDEN",
	"A function''s \"caller\" property cannot be used in strict mode")

// Java: StrictModeCheck.FUNCTION_ARGUMENTS_PROP_FORBIDDEN
var StrictModeCheck_FUNCTION_ARGUMENTS_PROP_FORBIDDEN = DiagnosticType_Error(
	"JSC_FUNCTION_ARGUMENTS_PROP_FORBIDDEN",
	"A function''s \"arguments\" property cannot be used in strict mode")

// Java: StrictModeCheck.DELETE_VARIABLE
var StrictModeCheck_DELETE_VARIABLE = DiagnosticType_Error(
	"JSC_DELETE_VARIABLE",
	"variables, functions, and arguments cannot be deleted in strict mode")

// Java: StrictModeCheck.DUPLICATE_MEMBER
var StrictModeCheck_DUPLICATE_MEMBER = DiagnosticType_Warning(
	"JSC_DUPLICATE_MEMBER",
	"Class or object literal contains duplicate member \"{0}\". In non-strict code, the last"+
		" duplicate will overwrite the others.")

// Java: StrictModeCheck(AbstractCompiler, CheckLevel)
func NewStrictModeCheck(compiler AbstractCompiler, defaultLevel CheckLevel) *StrictModeCheck {
	c := &StrictModeCheck{}
	c.compiler = compiler
	c.defaultLevel = defaultLevel
	return c
}

// Java: StrictModeCheck.process(Node, Node)
func (c *StrictModeCheck) Process(externs rhino.Node, root rhino.Node) {
	if c.processParallel(externs, root) {
		return
	}
	NodeTraversal_TraverseRoots(c.compiler, c, externs, root)
	NodeTraversal_Traverse(c.compiler, root, NewStrictModeCheck_NonExternChecks(c))
}

func (c *StrictModeCheck) processParallel(externs rhino.Node, root rhino.Node) bool {
	scopeRoot, items, ok := collectExternAndRootScripts(externs, root)
	if !ok {
		return false
	}
	numWorkers := parallelThreadCount(c.compiler, len(items))
	if numWorkers <= 1 {
		return false
	}

	prewarmCompilerForParallelReads(c.compiler)
	restoreReg := enableParallelTypeRegistryReads(c.compiler)
	defer restoreReg()
	lazyGlobal := newLazyGlobalScope(c.compiler, scopeRoot)

	pass2EventsByItem := make([][]recordedCompilerEvent, len(items))
	runParallelWork(c.compiler, numWorkers, len(items), func(_ int) (func(int, *recordingCompiler), func()) {
		var workerRC *recordingCompiler
		var workerPass1 *StrictModeCheck
		var workerPass2 *StrictModeCheck_NonExternChecks
		var pass2RC recordingCompiler
		var sc *SyntacticScopeCreator

		runItem := func(idx int, rc *recordingCompiler) {
			if workerPass1 == nil || workerRC != rc {
				workerRC = rc
				workerPass1 = NewStrictModeCheck(rc, c.defaultLevel)
				workerPass2 = NewStrictModeCheck_NonExternChecks(NewStrictModeCheck(&pass2RC, c.defaultLevel))
				sc = NewSyntacticScopeCreator(rc)
			}
			item := items[idx]
			b1 := NodeTraversal_Builder{
				compiler:     rc,
				callback:     workerPass1,
				scopeCreator: sc,
			}
			b1.traverseScriptUnderRoot(scopeRoot, item.parent, item.script, lazyGlobal)

			if item.parent == root {
				initRecordingCompiler(&pass2RC, c.compiler)
				b2 := NodeTraversal_Builder{
					compiler: &pass2RC,
					callback: workerPass2,
				}
				b2.traverseScriptUnderRoot(root, root, item.script, nil)
				pass2EventsByItem[idx] = pass2RC.takeEvents()
			}
		}
		return runItem, nil
	}, false)

	for i := range pass2EventsByItem {
		if len(pass2EventsByItem[i]) != 0 {
			replayCompilerEvents(c.compiler, pass2EventsByItem[i])
		}
	}
	return true
}

// Java: StrictModeCheck.visit(NodeTraversal, Node, Node)
func (c *StrictModeCheck) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_ASSIGN:
		c.CheckAssignment(ast.MustAssign(n))
	case rhino.Token_DELPROP:
		c.CheckDelete(t, ast.MustUnary(n))
	case rhino.Token_OBJECTLIT:
		c.CheckObjectLiteralOrClass(n)
	case rhino.Token_CLASS:
		c.CheckObjectLiteralOrClass(ast.MustClass(n).Body().Node())
	case rhino.Token_WITH:
		c.CheckWith(ast.MustWith(n))
	}
}

// CheckWith reports a warning for with statements.
//
// Java: StrictModeCheck.checkWith(Node)
func (c *StrictModeCheck) CheckWith(with ast.With) {
	info := with.Node().GetJSDocInfo()
	allowWith := info != nil && info.HasSuppression("with")
	if !allowWith {
		c.Report(with.Node(), StrictModeCheck_USE_OF_WITH)
	}
}

// StrictModeCheck_IsDeclaration determines if the given name is a
// declaration, which can be a declaration of a variable, function, or
// argument.
//
// Java: StrictModeCheck.isDeclaration(Node)
func StrictModeCheck_IsDeclaration(n rhino.Node) bool {
	switch n.GetParent().GetToken() {
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_VAR, rhino.Token_CATCH:
		return true
	case rhino.Token_FUNCTION:
		return n == n.GetParent().GetFirstChild()
	case rhino.Token_PARAM_LIST:
		return n.GetGrandparent().IsFunction()
	default:
		return false
	}
}

// CheckAssignment checks that an assignment is not to the "arguments" object.
//
// Java: StrictModeCheck.checkAssignment(Node)
func (c *StrictModeCheck) CheckAssignment(assign ast.Assign) {
	if name, ok := assign.Target().AsName(); ok {
		switch name.Node().GetStringID() {
		case strictModeArgumentsID:
			c.Report(assign.Node(), StrictModeCheck_ARGUMENTS_ASSIGNMENT)
		case strictModeEvalID:
			// Note that assignment to eval is already illegal because any use of
			// that name is illegal.
			c.Report(assign.Node(), StrictModeCheck_EVAL_ASSIGNMENT)
		}
	}
}

// CheckDelete checks that variables, functions, and arguments are not
// deleted.
//
// Java: StrictModeCheck.checkDelete(NodeTraversal, Node)
func (c *StrictModeCheck) CheckDelete(t *NodeTraversal, del ast.Unary) {
	if name, ok := del.Operand().AsName(); ok {
		v := t.GetScope().GetVar(name.Ident())
		if v != nil {
			c.Report(del.Node(), StrictModeCheck_DELETE_VARIABLE)
		}
	}
}

// CheckObjectLiteralOrClass checks that object literal keys or class method
// names are valid.
//
// Java: StrictModeCheck.checkObjectLiteralOrClass(Node)
func (c *StrictModeCheck) CheckObjectLiteralOrClass(n rhino.Node) {
	if n.GetFirstChild() == n.GetLastChild() {
		return
	}

	var seenBuf [16]strictModeMemberEntry
	seenLen := 0
	usedMap := false

	/*
	 * Iterate backwards because the last duplicate is the one that will be used in sloppy or ES6
	 * code. The earlier duplicates are the ones that should be removed.
	 */
	for key := n.GetLastChild(); !key.IsNil(); key = key.GetPrevious() {
		if key.IsEmpty() ||
			key.IsComputedProp() ||
			key.IsSpread() ||
			key.IsComputedFieldDef() ||
			// Computed properties cannot be computed at compile time
			NodeUtil_IsClassStaticBlock(key) {
			// Will not check since whether duplicates are declared/assigned cannot be determined at
			// compile time since we do not know which code will be run
			continue
		}

		keyID := key.GetStringID()
		isSetter := key.IsSetterDef()
		isGetter := key.IsGetterDef()
		getterBit := strictModeGetterBit
		setterBit := strictModeSetterBit
		if key.IsStaticMember() {
			getterBit = strictModeStaticGetterBit
			setterBit = strictModeStaticSetterBit
		}
		var addMask uint8
		if !isSetter {
			addMask |= getterBit
		}
		if !isGetter {
			addMask |= setterBit
		}

		var oldMask uint8
		if !usedMap {
			found := false
			for i := 0; i < seenLen; i++ {
				if seenBuf[i].id == keyID {
					oldMask = seenBuf[i].mask
					seenBuf[i].mask = oldMask | addMask
					found = true
					break
				}
			}
			if !found {
				if seenLen < len(seenBuf) {
					seenBuf[seenLen] = strictModeMemberEntry{id: keyID, mask: addMask}
					seenLen++
				} else {
					usedMap = true
					if c.memberFlags == nil {
						c.memberFlags = make(map[uint32]uint8, 64)
					} else {
						clear(c.memberFlags)
					}
					for i := range seenBuf {
						c.memberFlags[seenBuf[i].id] = seenBuf[i].mask
					}
					c.memberFlags[keyID] = addMask
				}
			}
		} else {
			oldMask = c.memberFlags[keyID]
			c.memberFlags[keyID] = oldMask | addMask
		}

		if !isSetter && (oldMask&getterBit) != 0 {
			// normal property and getter cases
			c.Report(key, StrictModeCheck_DUPLICATE_MEMBER, key.GetString())
		}
		if !isGetter && (oldMask&setterBit) != 0 {
			// normal property and setter cases
			c.Report(key, StrictModeCheck_DUPLICATE_MEMBER, key.GetString())
		}
	}
}

// StrictModeCheck_NonExternChecks holds checks that are performed on
// non-extern code only.
//
// Java: StrictModeCheck.NonExternChecks
type StrictModeCheck_NonExternChecks struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	// Go-only: language-adapter — the enclosing StrictModeCheck instance of
	// this inner class.
	outer *StrictModeCheck
}

// Java: StrictModeCheck.NonExternChecks()
func NewStrictModeCheck_NonExternChecks(outer *StrictModeCheck) *StrictModeCheck_NonExternChecks {
	return &StrictModeCheck_NonExternChecks{outer: outer}
}

// Java: StrictModeCheck.NonExternChecks.visit(NodeTraversal, Node, Node)
func (c *StrictModeCheck_NonExternChecks) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_NAME:
		if StrictModeCheck_IsDeclaration(n) {
			c.CheckDeclaration(ast.MustName(n))
		}
	case rhino.Token_GETPROP:
		c.CheckGetProp(ast.MustGetProp(n))
	}
}

// CheckDeclaration checks for illegal declarations.
//
// Java: StrictModeCheck.NonExternChecks.checkDeclaration(Node)
func (c *StrictModeCheck_NonExternChecks) CheckDeclaration(name ast.Name) {
	switch name.Node().GetStringID() {
	case strictModeEvalID:
		c.outer.Report(name.Node(), StrictModeCheck_EVAL_DECLARATION)
	case strictModeArgumentsID:
		c.outer.Report(name.Node(), StrictModeCheck_ARGUMENTS_DECLARATION)
	}
}

// CheckGetProp checks that the arguments.callee is not used.
//
// Java: StrictModeCheck.NonExternChecks.checkGetProp(Node)
func (c *StrictModeCheck_NonExternChecks) CheckGetProp(getProp ast.GetProp) {
	n := getProp.Node()
	target := getProp.Object()
	switch n.GetStringID() {
	case strictModeCalleeID:
		if strictModeCheck_isArguments(target) {
			c.outer.Report(n, StrictModeCheck_ARGUMENTS_CALLEE_FORBIDDEN)
		}
	case strictModeCallerID:
		if strictModeCheck_isArguments(target) {
			c.outer.Report(n, StrictModeCheck_ARGUMENTS_CALLER_FORBIDDEN)
		} else if StrictModeCheck_IsFunctionType(target.Node()) {
			c.outer.Report(n, StrictModeCheck_FUNCTION_CALLER_FORBIDDEN)
		}
	case strictModeArgumentsID:
		if StrictModeCheck_IsFunctionType(target.Node()) {
			c.outer.Report(n, StrictModeCheck_FUNCTION_ARGUMENTS_PROP_FORBIDDEN)
		}
	}
}

// strictModeCheck_isArguments reports whether e is the name arguments.
//
// Go-only: CheckGetProp's `target.isName() && target.getString().equals("arguments")`.
func strictModeCheck_isArguments(e ast.Expr) bool {
	name, ok := e.AsName()
	return ok && name.Node().GetStringID() == strictModeArgumentsID
}

// Java: StrictModeCheck.isFunctionType(Node)
func StrictModeCheck_IsFunctionType(n rhino.Node) bool {
	type_ := n.GetJSType()
	return type_ != nil && type_.IsFunctionType()
}

// Java: StrictModeCheck.report(Node, DiagnosticType, String...)
func (c *StrictModeCheck) Report(n rhino.Node, diagnostic *DiagnosticType, args ...string) {
	c.compiler.Report(
		JSError_Builder_(diagnostic, args...).SetLevel(c.defaultLevel).SetNode(n).Build())
}
