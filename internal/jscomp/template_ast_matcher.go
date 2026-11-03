// Port of com/google/javascript/jscomp/TemplateAstMatcher.java @ closure-compiler v20261028

/*
 * Copyright 2014 The Closure Compiler Authors.
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
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TemplateAstMatcher is a matcher that can take an arbitrary AST and use
// it as a template to find matches in another. As this matcher potentially
// matches against every node in the AST it is tuned to avoid generating GC
// garbage. It first checks the AST shape without types and then successful
// checks the associated types.
type TemplateAstMatcher struct {
	topScope     *TypedScope
	typeRegistry *rhino.JSTypeRegistry

	// The head of the Node list that should be used to start the matching
	// process.
	templateStart rhino.Node

	// The params declared in the template (in order)
	templateParams []string

	// Record the first Node to match a template parameter, only valid for
	// the last match if it was successful.
	paramNodeMatches []rhino.Node

	// The locals declared in the template (in order)
	templateLocals []string

	// Record the first name to match a template local variable, only valid
	// for the last match if it was successful.
	localVarMatches []base.OptString

	// Record whether the last successful was a loosely matched type, only
	// valid for the last match if it was successful.
	isLooseMatch bool

	// The strategy to use when matching the JSType of nodes.
	typeMatchingStrategy TypeMatchingStrategy
}

// Custom Token types for to use as placeholders in the template AST.
const (
	templateTypeParam     = rhino.Token_PLACEHOLDER1
	templateLocalName     = rhino.Token_PLACEHOLDER2
	templateStringLiteral = rhino.Token_PLACEHOLDER3
)

// NewTemplateAstMatcher constructs this matcher with a Function node that
// serves as the template to match all other nodes against. The body of the
// function will be used to match against.
//
// Java: TemplateAstMatcher(AbstractCompiler, Node, TypeMatchingStrategy)
func NewTemplateAstMatcher(
	compiler AbstractCompiler,
	templateFunctionNode rhino.Node,
	typeMatchingStrategy TypeMatchingStrategy,
) *TemplateAstMatcher {
	_, isFn := ast.AsFunction(templateFunctionNode)
	jutil.CheckStateMsg1(
		isFn,
		"Template node must be a function node. Received: %s",
		templateFunctionNode)

	m := &TemplateAstMatcher{
		// TopScope may be null if the template is used before type checking,
		// which is useful if just check code structure.
		topScope:     compiler.GetTopScope(),
		typeRegistry: jutil.CheckNotNull(compiler.GetTypeRegistry()),
		// Go-only: language-adapter — checkNotNull(typeMatchingStrategy): an
		// enum value, which cannot be null in Go.
		typeMatchingStrategy: typeMatchingStrategy,
	}
	m.templateStart = m.initTemplate(templateFunctionNode)
	return m
}

// Matches reports whether the node n matches the template.
//
// Java: TemplateAstMatcher.matches(Node)
func (m *TemplateAstMatcher) Matches(n rhino.Node) bool {
	if !m.matchesTemplateShape(m.templateStart, n) {
		return false
	}
	if len(m.paramNodeMatches) == 0 && len(m.localVarMatches) == 0 {
		// If there are no parameters or locals to match against, this has
		// been a successful match and there is no reason to traverse the
		// AST again.
		return true
	}
	m.reset()
	return m.matchesTemplate(m.templateStart, n)
}

// IsLooseMatch returns whether the last match succeeded due to loose type
// information.
//
// Java: TemplateAstMatcher.isLooseMatch()
func (m *TemplateAstMatcher) IsLooseMatch() bool {
	return m.isLooseMatch
}

// GetTemplateNodeToMatchMap returns a map from named template Nodes (such
// as parameters or local variables) to Nodes that were matches from the
// last matched template.
//
// Java: TemplateAstMatcher.getTemplateNodeToMatchMap()
func (m *TemplateAstMatcher) GetTemplateNodeToMatchMap() *OrderedStringMap[rhino.Node] {
	matches := NewOrderedStringMapWithCapacity[rhino.Node](len(m.templateParams) + len(m.templateLocals))
	for i, name := range m.templateParams {
		matches.Put(name, m.paramNodeMatches[i])
	}
	for i, name := range m.templateLocals {
		matchedName, ok := m.localVarMatches[i].Get()
		jutil.CheckState(ok)
		matches.Put(name, ast.NewName(matchedName).Node())
	}
	return matches
}

// initTemplate prepares a template AST to use when performing matches.
//
// templateFunctionNode is the template declaration function to extract the
// template AST from. Returns the first node of the template AST sequence to
// use when matching.
//
// Java: TemplateAstMatcher.initTemplate(Node)
func (m *TemplateAstMatcher) initTemplate(templateFunctionNode rhino.Node) rhino.Node {
	prepped := ast.MustFunction(templateFunctionNode.CloneTree())
	m.prepTemplatePlaceholders(prepped)

	var startNode rhino.Node
	if body, ok := prepped.Body().AsBlock(); ok {
		if onlyStmt, hasOnly := body.OnlyStatement(); hasOnly {
			if exprStmt, isExpr := onlyStmt.AsExprStmt(); isExpr {
				// When matching an expression, don't require it to be a complete
				// statement.
				startNode = exprStmt.Expr().Node()
			} else {
				startNode = onlyStmt.Node()
			}
		} else if firstStmt, hasFirst := body.FirstStatement(); hasFirst {
			startNode = firstStmt.Node()
		}
	} else {
		startNode = prepped.Body().Node().GetFirstChild()
	}

	m.localVarMatches = make([]base.OptString, len(m.templateLocals))
	m.paramNodeMatches = make([]rhino.Node, len(m.templateParams))
	return startNode
}

// prepTemplatePlaceholders builds parameter and local information for the
// template and replaces the references in the template 'fn' with
// placeholder nodes used to facilitate matching.
//
// Java: TemplateAstMatcher.prepTemplatePlaceholders(Node)
func (m *TemplateAstMatcher) prepTemplatePlaceholders(fn ast.Function) {
	// Drop the function name so it isn't included in the name maps.
	fnName := fn.NameNode().Ident()
	fn.RemoveName()

	// Build a list of parameter names and types.
	params := fn.Params()
	info := NodeUtil_GetBestJSDocInfo(fn.Node())
	var paramTypes []rhino.JSType
	if params.HasParams() {
		jutil.CheckNotNullMsg1(info,
			"Missing JSDoc declaration for template function %s", fnName)
		nParams := params.NumParams()
		m.templateParams = make([]string, 0, nParams)
		paramTypes = make([]rhino.JSType, 0, nParams)
	}
	// Go-only: language-adapter — a null topScope must reach the
	// registry as a null StaticTypedScope, not as a typed nil.
	var scope rhino.StaticTypedScope
	if m.topScope != nil {
		scope = m.topScope
	}
	for param := range params.Params() {
		var name string
		if nameNode, ok := param.AsName(); ok {
			name = nameNode.Ident()
		} else {
			name = param.Node().GetString()
		}
		expression := info.GetParameterType(name)
		jutil.CheckNotNullMsg2(expression,
			"Missing JSDoc for parameter %s of template function %s",
			name, fnName)
		paramType := m.typeRegistry.EvaluateTypeExpression(expression, scope)
		jutil.CheckNotNull(paramType)
		m.templateParams = append(m.templateParams, name)
		paramTypes = append(paramTypes, paramType)
	}

	// Find references to string literals, local variables and parameters
	// and replace them.
	traverseTemplatePostOrder(fn.Node(), func(n rhino.Node) {
		nameNode, ok := ast.AsName(n)
		if !ok {
			return
		}
		name := nameNode.Ident()
		if name != "" && n.GetParent().IsVar() && !slices.Contains(m.templateLocals, name) {
			m.templateLocals = append(m.templateLocals, name)
		}
		if paramIdx := slices.Index(m.templateParams, name); paramIdx >= 0 {
			paramType := paramTypes[paramIdx]
			isStringLiteral := paramType.IsStringValueType() && strings.HasPrefix(name, "string_literal")
			replaceTemplateNodeInPlace(n, createTemplateParameterNode(paramIdx, paramType, isStringLiteral))
		} else if localIdx := slices.Index(m.templateLocals, name); localIdx >= 0 {
			replaceTemplateNodeInPlace(n, createTemplateLocalNameNode(localIdx))
		}
	})
}

// Java: TemplateAstMatcher.replaceNodeInPlace(Node, Node)
func replaceTemplateNodeInPlace(n, replacement rhino.Node) {
	if n.HasChildren() {
		replacement.AddChildrenToFront(n.RemoveChildren())
	}
	n.ReplaceWith(replacement)
}

// Java: TemplateAstMatcher.traverse(Node, Visitor)
func traverseTemplatePostOrder(n rhino.Node, callback func(rhino.Node)) {
	for c := n.GetFirstChild(); !c.IsNil(); {
		next := c.GetNext() // in case the child is removed, grab the next node now
		traverseTemplatePostOrder(c, callback)
		c = next
	}
	callback(n)
}

// Java: TemplateAstMatcher.reset()
func (m *TemplateAstMatcher) reset() {
	m.isLooseMatch = false
	clear(m.localVarMatches)
	clear(m.paramNodeMatches)
}

// createTemplateParameterNode creates a template parameter or string
// literal template node.
//
// Java: TemplateAstMatcher.createTemplateParameterNode(int, JSType, boolean)
func createTemplateParameterNode(index int, paramType rhino.JSType, isStringLiteral bool) rhino.Node {
	jutil.CheckState(index >= 0)
	jutil.CheckNotNull(paramType)
	n := rhino.Node_NewNumber(float64(index))
	if isStringLiteral {
		n.SetToken(templateStringLiteral)
	} else {
		n.SetToken(templateTypeParam)
	}
	n.SetJSType(paramType)
	return n
}

// Java: TemplateAstMatcher.createTemplateLocalNameNode(int)
func createTemplateLocalNameNode(index int) rhino.Node {
	jutil.CheckState(index >= 0)
	n := rhino.Node_NewNumber(float64(index))
	n.SetToken(templateLocalName)
	return n
}

// matchesTemplateShape returns whether the template matches an AST
// structure node starting with n, taking into account the template
// parameters that were provided to this matcher. Here only the template
// shape is checked, template local declarations and parameters are checked
// later.
//
// Java: TemplateAstMatcher.matchesTemplateShape(Node, Node)
func (m *TemplateAstMatcher) matchesTemplateShape(template, n rhino.Node) bool {
	for !template.IsNil() {
		if n.IsNil() || !m.matchesNodeShape(template, n) {
			return false
		}
		template = template.GetNext()
		n = n.GetNext()
	}
	return true
}

// Java: TemplateAstMatcher.matchesNodeShape(Node, Node)
func (m *TemplateAstMatcher) matchesNodeShape(template, n rhino.Node) bool {
	switch template.GetToken() {
	case templateTypeParam:
		// Match the entire expression but only if it is an expression.
		return !n.IsNil() && !NodeUtil_IsStatement(n)
	case templateLocalName:
		// Match any name with matching initializer presence.
		if n.IsNil() || !n.IsName() || n.HasChildren() != template.HasChildren() {
			return false
		}
	case templateStringLiteral:
		// Matches parameters (in the refasterJS template) whose names start
		// with 'string_literal_'.
		return !n.IsNil() && NodeUtil_IsSomeCompileTimeConstStringValue(n)
	case rhino.Token_CALL:
		// Loosely match CALL nodes. IsEquivalentToShallow checks free calls
		// against non-free calls, but the template should ignore that
		// distinction.
		if n.IsNil() || !n.IsCall() || !n.HasXChildren(template.GetChildCount()) {
			return false
		}
	default:
		if !template.IsEquivalentToShallow(n) {
			return false
		}
	}

	// IsEquivalentToShallow (and the CALL / local-name checks above) guarantee
	// the child counts match.
	for tc, nc := template.GetFirstChild(), n.GetFirstChild(); !tc.IsNil(); tc, nc = tc.GetNext(), nc.GetNext() {
		if !m.matchesNodeShape(tc, nc) {
			return false
		}
	}
	return true
}

// Java: TemplateAstMatcher.matchesTemplate(Node, Node)
func (m *TemplateAstMatcher) matchesTemplate(template, n rhino.Node) bool {
	for !template.IsNil() {
		if n.IsNil() || !m.matchesNode(template, n) {
			return false
		}
		template = template.GetNext()
		n = n.GetNext()
	}
	return true
}

// matchesNode returns whether two nodes are equivalent, taking into account
// the template parameters that were provided to this matcher. If the
// template comparison node is a parameter node, then only the types of the
// node must match. If the template node is a string literal, only match
// string literals. Otherwise, the node must be equal and the child nodes
// must be equivalent according to the same function. This differs from the
// built in Node equivalence function with the special comparison.
//
// Java: TemplateAstMatcher.matchesNode(Node, Node)
func (m *TemplateAstMatcher) matchesNode(template, n rhino.Node) bool {
	switch template.GetToken() {
	case templateTypeParam:
		paramIndex := int(template.GetDouble())
		if previousMatch := m.paramNodeMatches[paramIndex]; !previousMatch.IsNil() {
			// If this named node has already been matched against, make sure
			// all subsequent usages of the same named node are equivalent.
			return n.IsEquivalentTo(previousMatch)
		}

		// Only the types need to match for the template parameters, which
		// allows the template function to express arbitrary expressions.
		templateType := template.GetJSType()
		jutil.CheckNotNullMsg(templateType, "null template parameter type.")

		// TODO(johnlenz): We shouldn't spend time checking template whose
		// types whose definitions aren't included (NoResolvedType).
		// Alternately we should treat them as "unknown" and perform loose
		// matches.
		if isUnresolvedTemplateType(templateType) {
			return false
		}

		matchResult := m.typeMatchingStrategy.Match(templateType, n.GetJSType())
		m.isLooseMatch = matchResult.IsLooseMatch()
		if !matchResult.IsMatch() {
			return false
		}
		m.paramNodeMatches[paramIndex] = n
		return true

	case templateLocalName:
		// If this template name node was already matched against, then make
		// sure all subsequent usages of the same template name node are
		// equivalent in the matched code.
		// For example, this code will handle the case:
		// function template() {
		//   var a = 'str';
		//   fn(a);
		// }
		//
		// will only match test code:
		//   var b = 'str';
		//   fn(b);
		//
		// but it will not match:
		//   var b = 'str';
		//   fn('str');
		localIndex := int(template.GetDouble())
		nameNode := ast.MustName(n)
		if prev, ok := m.localVarMatches[localIndex].Get(); ok {
			// If this named node has already been matched against, make sure
			// all subsequent usages of the same named node are equivalent.
			return nameNode.Ident() == prev
		}
		name, hasOrig := n.OriginalName()
		if !hasOrig {
			name = nameNode.Ident()
		}
		m.localVarMatches[localIndex] = base.SomeString(name)

	case templateStringLiteral:
		paramIndex := int(template.GetDouble())
		if previousMatch := m.paramNodeMatches[paramIndex]; !previousMatch.IsNil() {
			return n.IsEquivalentTo(previousMatch)
		}
		if NodeUtil_IsSomeCompileTimeConstStringValue(n) {
			m.paramNodeMatches[paramIndex] = n
			return true
		}
		return false
	}

	// Template and AST shape has already been checked, but continue look for
	// other template variables (parameters and locals) that must be checked.
	for tc, nc := template.GetFirstChild(), n.GetFirstChild(); !tc.IsNil(); tc, nc = tc.GetNext(), nc.GetNext() {
		if !m.matchesNode(tc, nc) {
			return false
		}
	}
	return true
}

// Java: TemplateAstMatcher.isUnresolvedType(JSType)
func isUnresolvedTemplateType(t rhino.JSType) bool {
	// TODO(b/146173738): When types are used in templates that do not appear
	// in the compilation unit being processed, the template type will be a
	// named type that resolves to unknown instead of being a no resolved
	// type. This should be fixed in the compiler such that it resolves to a
	// no resolved type, and then this code can be simplified to use that.
	if t.IsNoResolvedType() || (t.IsNamedType() && t.IsUnknownType()) {
		return true
	}
	if t.IsUnionType() {
		for _, alternate := range t.GetUnionMembers() {
			if isUnresolvedTemplateType(alternate) {
				return true
			}
		}
	}
	return false
}
