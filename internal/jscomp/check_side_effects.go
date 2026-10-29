// Port of com/google/javascript/jscomp/CheckSideEffects.java @ closure-compiler v20261021

/*
 * Copyright 2006 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// CheckSideEffects checks for non side effecting statements such as
//
//	var s = "this string is "
//	        "continued on the next line but you forgot the +";
//	x == foo();  // should that be '='?
//	foo();;  // probably just a stray-semicolon. Doesn't hurt to check though
//
// and generates warnings.
type CheckSideEffects struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	report bool

	problemNodes []rhino.Node

	noSideEffectExterns map[string]struct{}

	compiler AbstractCompiler

	protectSideEffectFreeCode bool

	// Whether the synthetic extern for JSCOMPILER_PRESERVE has been injected
	preserveFunctionInjected bool
}

// Java: CheckSideEffects.USELESS_CODE_ERROR
var CheckSideEffects_USELESS_CODE_ERROR = DiagnosticType_Warning("JSC_USELESS_CODE", "Suspicious code. {0}")

// Protect function is used to protect the underlying node from removal. At the
// very end, this protector function will be completely removed from production
// but the underlying node will be preserved.
//
// Java: CheckSideEffects.PROTECTOR_FN
const CheckSideEffects_PROTECTOR_FN = "JSCOMPILER_PRESERVE"

// J2CL protector function is used to protect the underlying node from removal
// until the very end so that referred prototype properties won't be removed by
// the compiler. At the very end, this protector function and its underyling
// node will be completely removed from production.
//
// Java: CheckSideEffects.J2CL_PROTECTOR_FN
const CheckSideEffects_J2CL_PROTECTOR_FN = "$J2CL_PRESERVE$"

var (
	checkSideEffectsProtectorID     = uint32(strpool.Intern(CheckSideEffects_PROTECTOR_FN))
	checkSideEffectsJ2clProtectorID = uint32(strpool.Intern(CheckSideEffects_J2CL_PROTECTOR_FN))
)

// Java: CheckSideEffects(AbstractCompiler, boolean, boolean)
func NewCheckSideEffects(compiler AbstractCompiler, report bool, protectSideEffectFreeCode bool) *CheckSideEffects {
	return &CheckSideEffects{
		compiler:                  compiler,
		report:                    report,
		protectSideEffectFreeCode: protectSideEffectFreeCode,
	}
}

// Java: CheckSideEffects.process(Node, Node)
func (c *CheckSideEffects) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(c.compiler, externs, &CheckSideEffects_GetNoSideEffectExterns{c: c})

	if !c.processRootParallel(root) {
		sc := NewSyntacticScopeCreator(c.compiler)
		b := NodeTraversal_Builder{
			compiler:     c.compiler,
			callback:     c,
			scopeCreator: sc,
		}
		b.Traverse(root)
	}

	// Code with hidden side-effect code is common, for example
	// accessing "el.offsetWidth" forces a reflow in browsers, to allow this
	// will still allowing local dead code removal in general,
	// protect the "side-effect free" code in the source.
	//
	// This also includes function calls such as with document.createElement
	if c.protectSideEffectFreeCode {
		c.ProtectSideEffects()
	}
}

func (c *CheckSideEffects) processRootParallel(root rhino.Node) bool {
	var problemsByScript [][]rhino.Node
	ok := traverseScriptsParallel(c.compiler, root, parallelScriptOptions[*CheckSideEffects]{
		useSharedGlobalScope: true,
		prepare: func(_, itemCount int) bool {
			problemsByScript = make([][]rhino.Node, itemCount)
			return true
		},
		initWorker: func(_ int, rc *Compiler) *CheckSideEffects {
			return &CheckSideEffects{
				report:                    c.report,
				noSideEffectExterns:       c.noSideEffectExterns,
				compiler:                  rc,
				protectSideEffectFreeCode: c.protectSideEffectFreeCode,
				preserveFunctionInjected:  c.preserveFunctionInjected,
			}
		},
		callback: func(w *CheckSideEffects, _ int, _, _ rhino.Node, _ *Compiler) NodeTraversal_Callback {
			w.problemNodes = w.problemNodes[:0]
			return w
		},
		afterScript: func(w *CheckSideEffects, idx int, _, _ rhino.Node, _ *Compiler) {
			if len(w.problemNodes) != 0 {
				copied := make([]rhino.Node, len(w.problemNodes))
				copy(copied, w.problemNodes)
				problemsByScript[idx] = copied
				w.problemNodes = w.problemNodes[:0]
			}
		},
	})
	if !ok {
		return false
	}
	for _, probs := range problemsByScript {
		if len(probs) != 0 {
			c.problemNodes = append(c.problemNodes, probs...)
		}
	}
	return true
}

// Java: CheckSideEffects.visit(NodeTraversal, Node, Node)
func (c *CheckSideEffects) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	// VOID nodes appear when there are extra semicolons at the BLOCK level.
	// I've been unable to think of any cases where this indicates a bug,
	// and apparently some people like keeping these semicolons around,
	// so we'll allow it.
	if n.IsEmpty() || n.IsComma() {
		return
	}

	if parent.IsNil() {
		return
	}

	// Do not try to remove a block or an expr result. We already handle
	// these cases when we visit the child, and the peephole passes will
	// fix up the tree in more clever ways when these are removed.
	if n.IsExprResult() || n.IsBlock() {
		return
	}

	// This no-op statement was there so that JSDoc information could
	// be attached to the name. This check should not complain about it.
	if n.IsQualifiedName() && n.GetJSDocInfo() != nil {
		return
	}

	// These are used by the compiler's built-in runtime libraries for dependency management
	// they have a special meaning used in Compiler.ensureLibraryInjected
	if n.IsStringLit() &&
		n.GetParent().IsExprResult() &&
		n.GetGrandparent().IsScript() &&
		strings.HasPrefix(*n.GetSourceFileName(), AbstractCompiler_RUNTIME_LIB_DIR) &&
		strings.HasPrefix(n.GetString(), "require ") {
		return
	}

	isResultUsed := NodeUtil_IsExpressionResultUsed(n)
	isSimpleOp := NodeUtil_IsSimpleOperator(n) || n.IsGetProp() || n.IsGetElem()
	if !isResultUsed && !n.IsVoid() {
		if isSimpleOp || !t.GetCompiler().GetAstAnalyzer().MayHaveSideEffects(n) {
			if c.report {
				msg := "This code lacks side-effects. Is there a bug?"
				if n.IsStringLit() || n.IsTemplateLit() {
					msg = "Is there a missing '+' on the previous line?"
				} else if isSimpleOp {
					msg = "The result of the '" +
						strings.ToLower(n.GetToken().String()) +
						"' operator is not being used."
				}

				t.Report(n, CheckSideEffects_USELESS_CODE_ERROR, msg)
			}
			// TODO(johnlenz): determine if it is necessary to
			// try to protect side-effect free statements as well.
			if !NodeUtil_IsStatement(n) {
				c.problemNodes = append(c.problemNodes, n)
			}
		} else if len(c.noSideEffectExterns) > 0 && n.IsCall() {
			first := n.GetFirstChild()
			if first.IsGetProp() || first.IsName() || first.IsStringLit() {
				var qnameStr string
				var ok bool
				if first.IsName() {
					qnameStr = first.GetString()
					ok = qnameStr != ""
				} else {
					qnameStr, ok = first.GetQualifiedNameString()
				}

				if ok {
					if _, isNoSideEffect := c.noSideEffectExterns[qnameStr]; isNoSideEffect {
						// The name should not be defined in src scopes - only externs
						isDefinedInSrc := false
						if first.IsGetProp() {
							rootNameNode := NodeUtil_GetRootOfQualifiedName(first)
							isDefinedInSrc = !rootNameNode.IsNil() &&
								rootNameNode.IsName() &&
								t.LookupVarByID(rootNameNode.GetStringID()) != nil
						} else if first.IsName() {
							isDefinedInSrc = t.LookupVarByID(first.GetStringID()) != nil
						} else {
							isDefinedInSrc = t.LookupVarByID(rhino.RhinoStringPool_LookupID(qnameStr)) != nil
						}

						if !isDefinedInSrc {
							c.problemNodes = append(c.problemNodes, n)
							if c.report {
								msg := "The result of the extern function call '" + qnameStr + "' is not being used."
								t.Report(n, CheckSideEffects_USELESS_CODE_ERROR, msg)
							}
						}
					}
				}
			}
		}
	}
}

// ProtectSideEffects protects side-effect free nodes by making them parameters
// to a extern function call. This call will be removed after all the
// optimizations passes have run.
//
// Java: CheckSideEffects.protectSideEffects()
func (c *CheckSideEffects) ProtectSideEffects() {
	if len(c.problemNodes) != 0 {
		if !c.preserveFunctionInjected {
			CheckSideEffects_AddExtern(c.compiler)
		}
		for _, n := range c.problemNodes {
			name := rhino.IR_Name(CheckSideEffects_PROTECTOR_FN).Srcref(n)
			name.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
			replacement := rhino.IR_Call(name).Srcref(n)
			replacement.PutBooleanProp(rhino.Node_FREE_CALL, true)
			n.ReplaceWith(replacement)
			replacement.AddChildToBack(n)
			c.compiler.ReportChangeToEnclosingScope(replacement)
		}
	}
}

// CheckSideEffects_AddExtern injects JSCOMPILER_PRESEVE into the synthetic
// externs.
//
// Java: CheckSideEffects.addExtern(AbstractCompiler)
func CheckSideEffects_AddExtern(compiler AbstractCompiler) {
	NodeUtil_CreateSynthesizedExternsSymbol(compiler, CheckSideEffects_PROTECTOR_FN)
}

// CheckSideEffects_StripProtection removes side-effect sync functions.
type CheckSideEffects_StripProtection struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler AbstractCompiler
}

// Java: CheckSideEffects.StripProtection(AbstractCompiler)
func NewCheckSideEffects_StripProtection(compiler AbstractCompiler) *CheckSideEffects_StripProtection {
	s := &CheckSideEffects_StripProtection{}
	s.compiler = compiler
	return s
}

// Java: CheckSideEffects.StripProtection.process(Node, Node)
func (s *CheckSideEffects_StripProtection) Process(externs rhino.Node, root rhino.Node) {
	if s.processParallel(root) {
		return
	}
	NodeTraversal_Traverse(s.compiler, root, s)
}

func (s *CheckSideEffects_StripProtection) processParallel(root rhino.Node) bool {
	return traverseScriptsParallel(s.compiler, root, parallelScriptOptions[*CheckSideEffects_StripProtection]{
		initWorker: func(_ int, rc *Compiler) *CheckSideEffects_StripProtection {
			return NewCheckSideEffects_StripProtection(rc)
		},
	})
}

var checkSideEffectsStripProtectionTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_CALL)
	return s
}()

var _ nodeTraversal_visitTokens = (*CheckSideEffects_StripProtection)(nil)

func (*CheckSideEffects_StripProtection) visitTokens() *ast.TokenSet {
	return &checkSideEffectsStripProtectionTokens
}

// Java: CheckSideEffects.StripProtection.visit(NodeTraversal, Node, Node)
func (s *CheckSideEffects_StripProtection) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsCall() {
		target := n.GetFirstChild()
		if !target.IsName() {
			return
		}

		targetID := target.GetStringID()
		if targetID == checkSideEffectsJ2clProtectorID {
			// Do not let non-J2CL code abuse it.
			jutil.CheckStateMsg(
				J2clSourceUtils_IsJ2clSource(n.GetSourceFileName()), "Only allowed for J2CL code")
			NodeUtil_DeleteFunctionCall(n, s.compiler)
			return
		}

		// TODO(johnlenz): add this to the coding convention
		// so we can remove goog.reflect.sinkValue as well.
		if targetID == checkSideEffectsProtectorID {
			expr := n.GetLastChild()
			n.DetachChildren()
			n.ReplaceWith(expr)
			t.ReportCodeChange()
		}
	}
}

// CheckSideEffects_GetNoSideEffectExterns gets fully qualified function names
// which are marked with @nosideeffects
//
// TODO(ChadKillingsworth) Add support for object literals
type CheckSideEffects_GetNoSideEffectExterns struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	c *CheckSideEffects
}

var checkSideEffectsNoSideEffectExternsTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_FUNCTION)
	s.Add(rhino.Token_NAME)
	return s
}()

var _ nodeTraversal_visitTokens = (*CheckSideEffects_GetNoSideEffectExterns)(nil)

func (*CheckSideEffects_GetNoSideEffectExterns) visitTokens() *ast.TokenSet {
	return &checkSideEffectsNoSideEffectExternsTokens
}

// Java: CheckSideEffects.GetNoSideEffectExterns.visit(NodeTraversal, Node, Node)
func (g *CheckSideEffects_GetNoSideEffectExterns) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsFunction() {
		jsDoc := NodeUtil_GetBestJSDocInfo(n)
		if jsDoc != nil && jsDoc.IsNoSideEffects() {
			if name, ok := NodeUtil_GetNameString(n); ok {
				if g.c.noSideEffectExterns == nil {
					g.c.noSideEffectExterns = make(map[string]struct{})
				}
				g.c.noSideEffectExterns[name] = struct{}{}
			}
		}
	}

	if n.IsName() && n.GetStringID() == checkSideEffectsProtectorID {
		g.c.preserveFunctionInjected = true
	}
}
