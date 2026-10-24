// Port of com/google/javascript/jscomp/TemplateAstMatcher.java @ closure-compiler v20261021

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
	templateParams *jutil.ArrayList[string]

	// Record the first Node to match a template parameter, only valid for
	// the last match if it was successful.
	paramNodeMatches *jutil.ArrayList[rhino.Node]

	// The locals declared in the template (in order)
	templateLocals *jutil.ArrayList[string]

	// Record the first name to match a template local variable, only valid
	// for the last match if it was successful.
	localVarMatches *jutil.ArrayList[*string]

	// The names of all matched string literals, in order.
	//
	// This re-uses strings already present in the AST, which is faster and
	// simpler than keeping an additional layer of indirection.
	stringLiteralMatches *jutil.LinkedHashMap[string, rhino.Node]

	// Record whether the last successful was a loosely matched type, only
	// valid for the last match if it was successful.
	isLooseMatch bool

	// The strategy to use when matching the JSType of nodes.
	typeMatchingStrategy TypeMatchingStrategy
}

// Custom Token types for to use as placeholders in the template AST.
const (
	templateAstMatcher_TEMPLATE_TYPE_PARAM     = rhino.Token_PLACEHOLDER1
	templateAstMatcher_TEMPLATE_LOCAL_NAME     = rhino.Token_PLACEHOLDER2
	templateAstMatcher_TEMPLATE_STRING_LITERAL = rhino.Token_PLACEHOLDER3
)

// NewTemplateAstMatcher constructs this matcher with a Function node that
// serves as the template to match all other nodes against. The body of the
// function will be used to match against.
//
// Java: TemplateAstMatcher(AbstractCompiler, Node, TypeMatchingStrategy)
func NewTemplateAstMatcher(
	compiler AbstractCompiler,
	templateFunctionNode rhino.Node,
	typeMatchingStrategy TypeMatchingStrategy) *TemplateAstMatcher {
	jutil.CheckStateMsg1(
		templateFunctionNode.IsFunction(),
		"Template node must be a function node. Received: %s",
		templateFunctionNode)

	m := &TemplateAstMatcher{
		templateParams:       jutil.NewArrayList[string](),
		paramNodeMatches:     jutil.NewArrayList[rhino.Node](),
		templateLocals:       jutil.NewArrayList[string](),
		localVarMatches:      jutil.NewArrayList[*string](),
		stringLiteralMatches: jutil.NewLinkedHashMap[string, rhino.Node](),
		isLooseMatch:         false,
	}
	// TopScope may be null if the template is used before type checking,
	// which is useful if just check code structure.
	m.topScope = compiler.GetTopScope()
	m.typeRegistry = jutil.CheckNotNull(compiler.GetTypeRegistry())
	m.templateStart = m.InitTemplate(templateFunctionNode)
	// Go-only: language-adapter — checkNotNull(typeMatchingStrategy): an
	// enum value, which cannot be null in Go.
	m.typeMatchingStrategy = typeMatchingStrategy
	return m
}

// Matches reports whether the node n matches the template.
//
// Java: TemplateAstMatcher.matches(Node)
func (m *TemplateAstMatcher) Matches(n rhino.Node) bool {
	if m.MatchesTemplateShape(m.templateStart, n) {
		if m.paramNodeMatches.IsEmpty() && m.localVarMatches.IsEmpty() {
			// If there are no parameters or locals to match against, this has
			// been a successful match and there is no reason to traverse the
			// AST again.
			return true
		}
		m.Reset()
		return m.MatchesTemplate(m.templateStart, n)
	}
	return false
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
func (m *TemplateAstMatcher) GetTemplateNodeToMatchMap() *jutil.LinkedHashMap[string, rhino.Node] {
	map_ := jutil.NewLinkedHashMapFrom(m.stringLiteralMatches)

	for i := 0; i < m.templateParams.Size(); i++ {
		name := m.templateParams.Get(i)
		map_.Put(name, m.paramNodeMatches.Get(i))
	}

	for i := 0; i < m.templateLocals.Size(); i++ {
		name := m.templateLocals.Get(i)
		map_.Put(name, rhino.IR_Name(*m.localVarMatches.Get(i)))
	}

	return map_
}

// InitTemplate prepares an template AST to use when performing matches.
//
// templateFunctionNode is the template declaration function to extract the
// template AST from. Returns the first node of the template AST sequence to
// use when matching.
//
// Java: TemplateAstMatcher.initTemplate(Node)
func (m *TemplateAstMatcher) InitTemplate(templateFunctionNode rhino.Node) rhino.Node {
	prepped := templateFunctionNode.CloneTree()
	m.PrepTemplatePlaceholders(prepped)

	body := prepped.GetLastChild()
	var startNode rhino.Node
	if body.HasOneChild() && body.GetFirstChild().IsExprResult() {
		// When matching an expression, don't require it to be a complete
		// statement.
		startNode = body.GetFirstFirstChild()
	} else {
		startNode = body.GetFirstChild()
	}

	for i := 0; i < m.templateLocals.Size(); i++ {
		// reserve space in the locals array.
		m.localVarMatches.Add(nil)
	}
	for i := 0; i < m.templateParams.Size(); i++ {
		// reserve space in the params array.
		m.paramNodeMatches.Add(rhino.NilNode)
	}

	return startNode
}

// PrepTemplatePlaceholders builds parameter and local information for the
// template and replaces the references in the template 'fn' with
// placeholder nodes use to facility matching.
//
// Java: TemplateAstMatcher.prepTemplatePlaceholders(Node)
func (m *TemplateAstMatcher) PrepTemplatePlaceholders(fn rhino.Node) {
	locals := m.templateLocals
	params := m.templateParams
	paramTypes := jutil.NewLinkedHashMap[string, rhino.JSType]()

	// drop the function name so it isn't include in the name maps
	fnName := fn.GetFirstChild().GetString()
	fn.GetFirstChild().SetString("")

	// Build a list of parameter names and types.
	templateParametersNode := fn.GetSecondChild()
	info := NodeUtil_GetBestJSDocInfo(fn)
	if templateParametersNode.HasChildren() {
		jutil.CheckNotNullMsg1(info,
			"Missing JSDoc declaration for template function %s", fnName)
	}
	for paramNode := templateParametersNode.GetFirstChild(); !paramNode.IsNil(); paramNode = paramNode.GetNext() {
		name := paramNode.GetString()
		expression := info.GetParameterType(name)
		jutil.CheckNotNullMsg2(expression,
			"Missing JSDoc for parameter %s of template function %s",
			name, fnName)
		// Go-only: language-adapter — a null topScope must reach the
		// registry as a null StaticTypedScope, not as a typed nil.
		var scope rhino.StaticTypedScope
		if m.topScope != nil {
			scope = m.topScope
		}
		type_ := m.typeRegistry.EvaluateTypeExpression(expression, scope)
		jutil.CheckNotNull(type_)
		params.Add(name)
		paramTypes.Put(name, type_)
	}

	// Find references to string literals, local variables and parameters
	// and replace them.
	m.Traverse(
		fn,
		func(n rhino.Node) {
			if n.IsName() {
				parent := n.GetParent()
				name := n.GetString()
				if name != "" && parent.IsVar() && !locals.Contains(name) {
					locals.Add(n.GetString())
				}

				if params.Contains(name) {
					type_ := paramTypes.Get(name)
					isStringLiteral := type_.IsStringValueType() && jutil.StringStartsWith(name, "string_literal")
					m.ReplaceNodeInPlace(
						n, m.CreateTemplateParameterNode(params.IndexOf(name), type_, isStringLiteral))
				} else if locals.Contains(name) {
					m.ReplaceNodeInPlace(n, m.CreateTemplateLocalNameNode(locals.IndexOf(name)))
				}
			}
		})
}

// Java: TemplateAstMatcher.replaceNodeInPlace(Node, Node)
func (m *TemplateAstMatcher) ReplaceNodeInPlace(n rhino.Node, replacement rhino.Node) {
	if n.HasChildren() {
		children := n.RemoveChildren()
		replacement.AddChildrenToFront(children)
	}
	n.ReplaceWith(replacement)
}

// templateAstMatcher_Visitor is the private Visitor interface.
//
// Go-only: language-adapter — a single-method interface implemented by an
// anonymous class, so a Go func type.
type templateAstMatcher_Visitor func(n rhino.Node)

// Java: TemplateAstMatcher.traverse(Node, Visitor)
func (m *TemplateAstMatcher) Traverse(n rhino.Node, callback templateAstMatcher_Visitor) {
	var next rhino.Node
	for c := n.GetFirstChild(); !c.IsNil(); c = next {
		next = c.GetNext() // in case the child is remove, grab the next node now
		m.Traverse(c, callback)
	}
	callback(n)
}

// Java: TemplateAstMatcher.reset()
func (m *TemplateAstMatcher) Reset() {
	m.isLooseMatch = false
	// Collections.fill(localVarMatches, null)
	for i := 0; i < m.localVarMatches.Size(); i++ {
		m.localVarMatches.Set(i, nil)
	}
	for i := 0; i < m.paramNodeMatches.Size(); i++ {
		m.paramNodeMatches.Set(i, rhino.NilNode)
	}
}

// Java: TemplateAstMatcher.isTemplateParameterNode(Node)
func (m *TemplateAstMatcher) IsTemplateParameterNode(n rhino.Node) bool {
	return n.GetToken() == templateAstMatcher_TEMPLATE_TYPE_PARAM
}

// IsTemplateParameterStringLiteralNode matches parameters (in the
// refasterJS template) whose names start with 'string_literal_'.
//
// Java: TemplateAstMatcher.isTemplateParameterStringLiteralNode(Node)
func (m *TemplateAstMatcher) IsTemplateParameterStringLiteralNode(n rhino.Node) bool {
	return n.GetToken() == templateAstMatcher_TEMPLATE_STRING_LITERAL
}

// CreateTemplateParameterNode creates a template parameter or string
// literal template node.
//
// Java: TemplateAstMatcher.createTemplateParameterNode(int, JSType, boolean)
func (m *TemplateAstMatcher) CreateTemplateParameterNode(index int, type_ rhino.JSType, isStringLiteral bool) rhino.Node {
	jutil.CheckState(index >= 0)
	jutil.CheckNotNull(type_)
	n := rhino.Node_NewNumber(float64(index))
	if isStringLiteral {
		n.SetToken(templateAstMatcher_TEMPLATE_STRING_LITERAL)
	} else {
		n.SetToken(templateAstMatcher_TEMPLATE_TYPE_PARAM)
	}
	n.SetJSType(type_)
	return n
}

// Java: TemplateAstMatcher.isTemplateLocalNameNode(Node)
func (m *TemplateAstMatcher) IsTemplateLocalNameNode(n rhino.Node) bool {
	return n.GetToken() == templateAstMatcher_TEMPLATE_LOCAL_NAME
}

// Java: TemplateAstMatcher.createTemplateLocalNameNode(int)
func (m *TemplateAstMatcher) CreateTemplateLocalNameNode(index int) rhino.Node {
	jutil.CheckState(index >= 0)
	n := rhino.Node_NewNumber(float64(index))
	n.SetToken(templateAstMatcher_TEMPLATE_LOCAL_NAME)
	return n
}

// MatchesTemplateShape returns whether the template matches an AST
// structure node starting with node, taking into account the template
// parameters that were provided to this matcher. Here only the template
// shape is checked, template local declarations and parameters are checked
// later.
//
// Java: TemplateAstMatcher.matchesTemplateShape(Node, Node)
func (m *TemplateAstMatcher) MatchesTemplateShape(template rhino.Node, ast rhino.Node) bool {
	for !template.IsNil() {
		if ast.IsNil() || !m.MatchesNodeShape(template, ast) {
			return false
		}
		template = template.GetNext()
		ast = ast.GetNext()
	}
	return true
}

// Java: TemplateAstMatcher.matchesNodeShape(Node, Node)
func (m *TemplateAstMatcher) MatchesNodeShape(template rhino.Node, ast rhino.Node) bool {
	if m.IsTemplateParameterNode(template) {
		// Match the entire expression but only if it is an expression.
		return !NodeUtil_IsStatement(ast)
	} else if m.IsTemplateLocalNameNode(template) {
		// Match any name. Maybe match locals here.
		if !ast.IsName() {
			return false
		}
	} else if m.IsTemplateParameterStringLiteralNode(template) {
		// Matches parameters (in the refasterJS template) whose names start
		// with 'string_literal_'.
		return NodeUtil_IsSomeCompileTimeConstStringValue(ast)
	} else if template.IsCall() {
		// Loosely match CALL nodes. isEquivalentToShallow checks free calls
		// against non-free calls, but the template should ignore that
		// distinction.
		if ast.IsNil() || !ast.IsCall() || !ast.HasXChildren(template.GetChildCount()) {
			return false
		}
		// But check any children.
	} else if !template.IsEquivalentToShallow(ast) {
		return false
	}

	// isEquivalentToShallow guarantees the child counts match
	templateChild := template.GetFirstChild()
	astChild := ast.GetFirstChild()
	for !templateChild.IsNil() {
		if !m.MatchesNodeShape(templateChild, astChild) {
			return false
		}
		templateChild = templateChild.GetNext()
		astChild = astChild.GetNext()
	}
	return true
}

// Java: TemplateAstMatcher.matchesTemplate(Node, Node)
func (m *TemplateAstMatcher) MatchesTemplate(template rhino.Node, ast rhino.Node) bool {
	for !template.IsNil() {
		if ast.IsNil() || !m.MatchesNode(template, ast) {
			return false
		}
		template = template.GetNext()
		ast = ast.GetNext()
	}
	return true
}

// MatchesNode returns whether two nodes are equivalent, taking into account
// the template parameters that were provided to this matcher. If the
// template comparison node is a parameter node, then only the types of the
// node must match. If the template node is a string literal, only match
// string literals. Otherwise, the node must be equal and the child nodes
// must be equivalent according to the same function. This differs from the
// built in Node equivalence function with the special comparison.
//
// Java: TemplateAstMatcher.matchesNode(Node, Node)
func (m *TemplateAstMatcher) MatchesNode(template rhino.Node, ast rhino.Node) bool {
	if m.IsTemplateParameterNode(template) {
		paramIndex := int(jutil.D2I(template.GetDouble()))
		previousMatch := m.paramNodeMatches.Get(paramIndex)
		if !previousMatch.IsNil() {
			// If this named node has already been matched against, make sure
			// all subsequent usages of the same named node are equivalent.
			return ast.IsEquivalentTo(previousMatch)
		}

		// Only the types need to match for the template parameters, which
		// allows the template function to express arbitrary expressions.
		templateType := template.GetJSType()

		jutil.CheckNotNullMsg(templateType, "null template parameter type.")

		// TODO(johnlenz): We shouldn't spend time checking template whose
		// types whose definitions aren't included (NoResolvedType).
		// Alternately we should treat them as "unknown" and perform loose
		// matches.
		if m.IsUnresolvedType(templateType) {
			return false
		}

		matchResult := m.typeMatchingStrategy.Match(templateType, ast.GetJSType())
		m.isLooseMatch = matchResult.IsLooseMatch()
		isMatch := matchResult.IsMatch()
		if isMatch && previousMatch.IsNil() {
			m.paramNodeMatches.Set(paramIndex, ast)
		}
		return isMatch
	} else if m.IsTemplateLocalNameNode(template) {
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
		paramIndex := int(jutil.D2I(template.GetDouble()))
		previouslyMatched := m.localVarMatches.Get(paramIndex) != nil
		if previouslyMatched {
			// If this named node has already been matched against, make sure
			// all subsequent usages of the same named node are equivalent.
			return ast.GetString() == *m.localVarMatches.Get(paramIndex)
		} else {
			originalName := ast.GetOriginalName()
			var name string
			if originalName != nil {
				name = *originalName
			} else {
				name = ast.GetString()
			}
			m.localVarMatches.Set(paramIndex, &name)
		}
	} else if m.IsTemplateParameterStringLiteralNode(template) {
		paramIndex := int(jutil.D2I(template.GetDouble()))
		previousMatch := m.paramNodeMatches.Get(paramIndex)
		if !previousMatch.IsNil() {
			return ast.IsEquivalentTo(previousMatch)
		}

		if NodeUtil_IsSomeCompileTimeConstStringValue(ast) {
			m.paramNodeMatches.Set(paramIndex, ast)
			return true
		}
		return false
	}

	// Template and AST shape has already been checked, but continue look for
	// other template variables (parameters and locals) that must be checked.
	templateChild := template.GetFirstChild()
	astChild := ast.GetFirstChild()
	for !templateChild.IsNil() {
		if !m.MatchesNode(templateChild, astChild) {
			return false
		}
		templateChild = templateChild.GetNext()
		astChild = astChild.GetNext()
	}

	return true
}

// Java: TemplateAstMatcher.isUnresolvedType(JSType)
func (m *TemplateAstMatcher) IsUnresolvedType(type_ rhino.JSType) bool {
	// TODO(b/146173738): When types are used in templates that do not appear
	// in the compilation unit being processed, the template type will be a
	// named type that resolves to unknown instead of being a no resolved
	// type. This should be fixed in the compiler such that it resolves to a
	// no resolved type, and then this code can be simplified to use that.
	if type_.IsNoResolvedType() || (type_.IsNamedType() && type_.IsUnknownType()) {
		return true
	}
	if type_.IsUnionType() {
		for _, alternate := range type_.GetUnionMembers() {
			if m.IsUnresolvedType(alternate) {
				return true
			}
		}
	}
	return false
}
