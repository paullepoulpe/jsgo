// Port of com/google/javascript/jscomp/ImplicitNullabilityCheck.java @ closure-compiler v20261021

/*
 * Copyright 2015 The Closure Compiler Authors.
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

// ImplicitNullabilityCheck warns about types in JSDoc that are implicitly
// nullable.
type ImplicitNullabilityCheck struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler AbstractCompiler
}

var ImplicitNullabilityCheck_IMPLICITLY_NULLABLE_JSDOC = DiagnosticType_Disabled(
	"JSC_IMPLICITLY_NULLABLE_JSDOC",
	"Name {0} in JSDoc is implicitly nullable, and is discouraged by the style guide.\n"+
		"Please add a '!' to make it non-nullable, or a '?' to make it explicitly nullable.")

var ImplicitNullabilityCheck_IMPLICITLY_NONNULL_JSDOC = DiagnosticType_Disabled(
	"JSC_IMPLICITLY_NONNULL_JSDOC",
	"Name {0} in JSDoc is implicitly non-null, and is discouraged by the style guide.\n"+
		"Please add a '!' to make it explicit.")

func implicitNullabilityCheck_isNullabilityOmittedType(typeName string) bool {
	switch typeName {
	case "*", "?", "bigint", "boolean", "null", "number", "string", "symbol", "undefined", "void":
		return true
	default:
		return false
	}
}

// Java: ImplicitNullabilityCheck(AbstractCompiler)
func NewImplicitNullabilityCheck(compiler AbstractCompiler) *ImplicitNullabilityCheck {
	c := &ImplicitNullabilityCheck{}
	c.compiler = compiler
	return c
}

// Java: ImplicitNullabilityCheck.process(Node, Node)
func (c *ImplicitNullabilityCheck) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_TraverseRoots(c.compiler, c, externs, root)
}

// ImplicitNullabilityCheck_Nullability represents the types of implicit
// nullability errors caught by this pass: a) implicitly nonnull (missing a
// "!"), and b) implicitly nullable (missing a "?").
type ImplicitNullabilityCheck_Nullability int32

const (
	ImplicitNullabilityCheck_Nullability_NONNULL ImplicitNullabilityCheck_Nullability = iota
	ImplicitNullabilityCheck_Nullability_NULLABLE
)

var implicitNullabilityCheck_Nullability_names = [...]string{"NONNULL", "NULLABLE"}

// Java: ImplicitNullabilityCheck.Nullability.values()
func ImplicitNullabilityCheck_Nullability_values() []ImplicitNullabilityCheck_Nullability {
	return []ImplicitNullabilityCheck_Nullability{
		ImplicitNullabilityCheck_Nullability_NONNULL, ImplicitNullabilityCheck_Nullability_NULLABLE}
}

// String is Enum.name().
func (n ImplicitNullabilityCheck_Nullability) String() string {
	return implicitNullabilityCheck_Nullability_names[n]
}

// Ordinal is Enum.ordinal().
func (n ImplicitNullabilityCheck_Nullability) Ordinal() int { return int(n) }

// Java: ImplicitNullabilityCheck.Nullability.isNullable()
func (n ImplicitNullabilityCheck_Nullability) IsNullable() bool {
	return n == ImplicitNullabilityCheck_Nullability_NULLABLE
}

// ImplicitNullabilityCheck_Result is information to represent a single
// "implicit nullability result", including the JSDoc string node that needs
// a "!" or "?" to be explicit, as well as an enum indicating which of the
// two nullability cases were found ("!" or "?").
type ImplicitNullabilityCheck_Result struct {
	node        rhino.Node
	nullability ImplicitNullabilityCheck_Nullability
}

// Java: ImplicitNullabilityCheck.Result(Node, Nullability)
func newImplicitNullabilityCheck_Result(node rhino.Node, nullability ImplicitNullabilityCheck_Nullability) *ImplicitNullabilityCheck_Result {
	r := &ImplicitNullabilityCheck_Result{}
	jutil.CheckArgument(node.IsStringLit())
	r.node = node
	r.nullability = nullability
	return r
}

// Java: ImplicitNullabilityCheck.Result.create(Node, Nullability)
func ImplicitNullabilityCheck_Result_Create(node rhino.Node, nullability ImplicitNullabilityCheck_Nullability) *ImplicitNullabilityCheck_Result {
	return newImplicitNullabilityCheck_Result(node, nullability)
}

// Java: ImplicitNullabilityCheck.Result.getNullability()
func (r *ImplicitNullabilityCheck_Result) GetNullability() ImplicitNullabilityCheck_Nullability {
	return r.nullability
}

// Java: ImplicitNullabilityCheck.Result.getNode()
func (r *ImplicitNullabilityCheck_Result) GetNode() rhino.Node {
	return r.node
}

// ImplicitNullabilityCheck_FindImplicitNullabilityResults finds and returns
// all the JSDoc nodes inside the given JSDoc object whose nullability is not
// explict, using the NodeTraversal the necessary state (current scope,
// etc.)
//
// Java: ImplicitNullabilityCheck.findImplicitNullabilityResults(JSDocInfo, NodeTraversal)
func ImplicitNullabilityCheck_FindImplicitNullabilityResults(
	info *rhino.JSDocInfo, t *NodeTraversal) []*ImplicitNullabilityCheck_Result {
	if info == nil {
		return []*ImplicitNullabilityCheck_Result{}
	}

	builder := []*ImplicitNullabilityCheck_Result{}
	for it := info.TypeExprIter(); it.Next(); {
		typeRoot := it.Node()
		NodeUtil_VisitPreOrderWithPred(
			typeRoot,
			NodeUtil_VisitorFunc(func(node rhino.Node) {
				if !node.IsStringLit() {
					return
				}
				parent := node.GetParent()
				if !parent.IsNil() {
					switch parent.GetToken() {
					case rhino.Token_BANG,
						rhino.Token_QMARK,
						rhino.Token_THIS,   // The names inside function(this:Foo) and
						rhino.Token_NEW,    // function(new:Bar) are already non-null.
						rhino.Token_TYPEOF: // Names after 'typeof' don't have nullability.
						return
					case rhino.Token_PIPE:
						// Inside a union
						gp := parent.GetParent()
						if !gp.IsNil() && gp.GetToken() == rhino.Token_QMARK {
							return // Inside an explicitly nullable union
						}
						for child := parent.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
							if (child.IsStringLit() && child.GetString() == "null") ||
								child.GetToken() == rhino.Token_QMARK {
								return // Inside a union that contains null or nullable type
							}
						}
					default:
					}
				}
				typeName := node.GetString()
				if implicitNullabilityCheck_isNullabilityOmittedType(typeName) {
					return
				}
				registry := t.GetCompiler().GetTypeRegistry()
				type_ := registry.GetType(t.GetScope(), typeName)
				if type_ == nil {
					return
				}
				isNonNullableName :=
					registry.IsNonNullableName(t.GetScope(), typeName) && !type_.IsNullable()
				nullability := ImplicitNullabilityCheck_Nullability_NULLABLE
				if isNonNullableName {
					nullability = ImplicitNullabilityCheck_Nullability_NONNULL
				}
				builder = append(builder, ImplicitNullabilityCheck_Result_Create(node, nullability))
			}),
			func(rhino.Node) bool { return true })
	}
	return builder
}

// Visit crawls the JSDoc of the given node to find any names in JSDoc that
// are implicitly null.
//
// Java: ImplicitNullabilityCheck.visit(NodeTraversal, Node, Node)
func (c *ImplicitNullabilityCheck) Visit(t *NodeTraversal, n rhino.Node, p rhino.Node) {
	info := n.GetJSDocInfo()
	for _, r := range ImplicitNullabilityCheck_FindImplicitNullabilityResults(info, t) {
		stringNode := r.GetNode()
		dt := ImplicitNullabilityCheck_IMPLICITLY_NONNULL_JSDOC
		if r.GetNullability().IsNullable() {
			dt = ImplicitNullabilityCheck_IMPLICITLY_NULLABLE_JSDOC
		}
		c.compiler.Report(JSError_MakeWithNode(stringNode, dt, stringNode.GetString()))
	}
}

var (
	_ NodeTraversal_Callback = (*ImplicitNullabilityCheck)(nil)
	_ CompilerPass           = (*ImplicitNullabilityCheck)(nil)
)
