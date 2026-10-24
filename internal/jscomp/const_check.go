// Port of com/google/javascript/jscomp/ConstCheck.java @ closure-compiler v20261021

/*
 * Copyright 2004 The Closure Compiler Authors.
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
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ConstCheck verifies that constants are only assigned a value once.
//
// e.g.
//
//	XX = 3;    // error!
//	XX++;      // error!
//
// TODO(tbreisacher): Consider merging this with CheckAccessControls so that
// all const-related checks are in the same place.
type ConstCheck struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler             AbstractCompiler
	initializedConstants map[*rhino.Node]struct{}
}

// Java: ConstCheck.CONST_REASSIGNED_VALUE_ERROR
var ConstCheck_CONST_REASSIGNED_VALUE_ERROR = DiagnosticType_Warning(
	"JSC_CONSTANT_REASSIGNED_VALUE_ERROR",
	"constant {0} assigned a value more than once.\nOriginal definition at {1}")

// NewConstCheck creates an instance.
//
// Java: ConstCheck(AbstractCompiler)
func NewConstCheck(compiler AbstractCompiler) *ConstCheck {
	c := &ConstCheck{}
	c.compiler = compiler
	c.initializedConstants = make(map[*rhino.Node]struct{})
	return c
}

func (c *ConstCheck) addInitializedConstant(var_ *Var) bool {
	key := var_.nameNode
	if key == nil {
		key = var_.scope.rootNode
	}
	if _, exists := c.initializedConstants[key]; exists {
		return false
	}
	c.initializedConstants[key] = struct{}{}
	return true
}

// Java: ConstCheck.process(Node, Node)
func (c *ConstCheck) Process(externs *rhino.Node, root *rhino.Node) {
	sc := NewSyntacticScopeCreator(c.compiler)
	defer func() {
		clear(c.initializedConstants)
	}()
	b := NodeTraversal_Builder{
		compiler:     c.compiler,
		callback:     c,
		scopeCreator: sc,
	}
	b.TraverseRoots(externs, root)
}

func (*ConstCheck) onlyNameLeaves() {}

// Java: ConstCheck.visit(NodeTraversal, Node, Node)
func (c *ConstCheck) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_NAME:
		if NodeUtil_IsNameDeclaration(parent) {
			name := n.GetString()
			var_ := t.LookupVar(name)
			if c.IsConstant(var_, n) {
				// If a constant is declared in externs, add it to initializedConstants to indicate
				// that it is initialized externally.
				if n.IsFromExterns() {
					c.addInitializedConstant(var_)
				} else if n.HasChildren() {
					if !c.addInitializedConstant(var_) {
						c.ReportError(n, var_, name)
					}
				}
			}
		}
	case rhino.Token_ASSIGN,
		rhino.Token_ASSIGN_BITOR,
		rhino.Token_ASSIGN_BITXOR,
		rhino.Token_ASSIGN_BITAND,
		rhino.Token_ASSIGN_LSH,
		rhino.Token_ASSIGN_RSH,
		rhino.Token_ASSIGN_URSH,
		rhino.Token_ASSIGN_ADD,
		rhino.Token_ASSIGN_SUB,
		rhino.Token_ASSIGN_MUL,
		rhino.Token_ASSIGN_DIV,
		rhino.Token_ASSIGN_MOD,
		rhino.Token_ASSIGN_EXPONENT:
		lhs := n.GetFirstChild()
		if lhs.IsName() {
			name := lhs.GetString()
			var_ := t.LookupVar(name)
			if c.IsConstant(var_, lhs) && !c.addInitializedConstant(var_) {
				c.ReportError(n, var_, name)
			} else if var_ != nil && var_.IsGoogModuleExports() && !c.addInitializedConstant(var_) {
				c.compiler.Report(
					JSError_MakeWithNode(n, ConstCheck_CONST_REASSIGNED_VALUE_ERROR, "exports",
						jutil.ToJavaString(n.GetSourceFileName())))
			}
		}
	case rhino.Token_INC, rhino.Token_DEC:
		lhs := n.GetFirstChild()
		if lhs.IsName() {
			name := lhs.GetString()
			var_ := t.LookupVar(name)
			if c.IsConstant(var_, lhs) {
				c.ReportError(n, var_, name)
			}
		}
	default:
	}
}

// IsConstant gets whether a variable is a constant initialized to a literal
// value at the point where it is declared.
//
// Java: ConstCheck.isConstant(Var, Node)
func (c *ConstCheck) IsConstant(var_ *Var, nameNode *rhino.Node) bool {
	jutil.CheckNotNullMsg1(var_, "Found unexpected undeclared name %s", nameNode)
	return var_.IsConst() || var_.IsDeclaredOrInferredConst()
}

// ReportError reports a reassigned constant error.
//
// Java: ConstCheck.reportError(Node, Var, String)
func (c *ConstCheck) ReportError(n *rhino.Node, var_ *Var, name string) {
	info := NodeUtil_GetBestJSDocInfo(n)
	if info == nil || !info.HasSuppression("const") {
		declNode := var_.GetNode()
		declaredPosition := jutil.ToJavaString(declNode.GetSourceFileName()) + ":" +
			strconv.Itoa(declNode.GetLineno())
		c.compiler.Report(JSError_MakeWithNode(n, ConstCheck_CONST_REASSIGNED_VALUE_ERROR, name, declaredPosition))
	}
}
