// Port of com/google/javascript/jscomp/LateEs6ToEs3Converter.java @ closure-compiler v20261021

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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass uses the typed AST view (internal/ast) to read the
// shapes it rewrites: the target and value of the assignment that defines
// $jscomp.createTemplateTagFirstArgWithRaw, the key and function of a
// method, and the members of an object literal with computed properties
// (the key and value of each computed property, the value of each string
// key). Visit still switches on the token (it runs on every node) and hands
// typed nodes to the rewrites. The nodes it builds come from the AstFactory
// (which sets their colors), so they stay rhino.Nodes, and a string key still
// becomes the string literal of its property access in place (SetToken).

// LateEs6ToEs3Converter converts ES6 code to valid ES5 code. This class does
// most of the transpilation, and
// https://github.com/google/closure-compiler/wiki/ECMAScript6 lists which ES6
// features are supported. Other classes that start with "Es6" do other parts
// of the transpilation.
//
// In most cases, the output is valid as ES3 (hence the class name) but in
// some cases, if the output language is set to ES5, we rely on ES5 features
// such as getters, setters, and Object.defineProperties.
//
// TODO(tbreisacher): This class does too many things. Break it into smaller
// passes.
type LateEs6ToEs3Converter struct {
	compiler                 AbstractCompiler
	astFactory               *AstFactory
	templateLiteralConverter *Es6TemplateLiterals

	// We want to insert the call to `var tagFnFirstArg = $jscomp.createTemplateTagFirstArg...` just
	// before this node. For the first script, this node is right after the runtime injected function
	// definition as injecting to the top of the script causes runtime errors
	// https://github.com/google/closure-compiler/issues/3589. For the subsequent script(s), the call
	// is injected to the top of that script.
	templateLitInsertionPoint rhino.Node
}

// Java: LateEs6ToEs3Converter.transpiledFeatures
var LateEs6ToEs3Converter_transpiledFeatures = parser.FeatureSet_BARE_MINIMUM.WithFeatures(
	parser.FeatureSet_Feature_COMPUTED_PROPERTIES,
	parser.FeatureSet_Feature_MEMBER_DECLARATIONS,
	parser.FeatureSet_Feature_TEMPLATE_LITERALS)

// Java: LateEs6ToEs3Converter.FRESH_COMP_PROP_VAR
const LateEs6ToEs3Converter_FRESH_COMP_PROP_VAR = "$jscomp$compprop"

var lateEs6CreateTemplateTagFirstArgWithRaw = rhino.QualifiedName_Of("$jscomp.createTemplateTagFirstArgWithRaw")

var lateEs6ToEs3VisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_ASSIGN)
	s.Add(rhino.Token_OBJECTLIT)
	s.Add(rhino.Token_MEMBER_FUNCTION_DEF)
	s.Add(rhino.Token_TAGGED_TEMPLATELIT)
	s.Add(rhino.Token_TEMPLATELIT)
	return s
}()

var (
	_ nodeTraversal_visitTokens   = (*LateEs6ToEs3Converter)(nil)
	_ nodeTraversal_noVisitLeaves = (*LateEs6ToEs3Converter)(nil)
)

func (*LateEs6ToEs3Converter) noVisitLeaves()             {}
func (*LateEs6ToEs3Converter) visitTokens() *ast.TokenSet { return &lateEs6ToEs3VisitTokens }

// Java: LateEs6ToEs3Converter(AbstractCompiler)
func NewLateEs6ToEs3Converter(compiler AbstractCompiler) *LateEs6ToEs3Converter {
	l := &LateEs6ToEs3Converter{}
	l.templateLitInsertionPoint = rhino.NilNode
	l.compiler = compiler
	l.astFactory = compiler.CreateAstFactory()
	l.templateLiteralConverter = NewEs6TemplateLiterals(compiler)
	return l
}

// Java: LateEs6ToEs3Converter.process(Node, Node)
func (l *LateEs6ToEs3Converter) Process(externs rhino.Node, root rhino.Node) {
	TranspilationPasses_ProcessTranspile(l.compiler, root, LateEs6ToEs3Converter_transpiledFeatures, l)
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(l.compiler, root, LateEs6ToEs3Converter_transpiledFeatures)
}

// Java: LateEs6ToEs3Converter.shouldTraverse(NodeTraversal, Node, Node)
func (l *LateEs6ToEs3Converter) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsScript() {
		l.templateLitInsertionPoint = LateEs6ToEs3Converter_FindTemplateLitInsertionPoint(n)
	}
	return true
}

// Java: LateEs6ToEs3Converter.findTemplateLitInsertionPoint(Node)
func LateEs6ToEs3Converter_FindTemplateLitInsertionPoint(script rhino.Node) rhino.Node {
	if script.GetIsInClosureUnawareSubtree() {
		// For all closure-unaware scripts, initialize the injection point within the closure
		// aware function.
		closureUnawareBlock := NodeUtil_FindClosureUnawareScriptRoot(script)
		return NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(closureUnawareBlock)
	}
	// For all other scripts, initialize the injection point to be top of the script. The spec
	// requires a single unique array per template literal, so for a template literal in a function
	// e.g. it would be incorrect to initialize a new array per every function call.
	return script.GetFirstChild()
}

// Java: LateEs6ToEs3Converter.visit(NodeTraversal, Node, Node)
func (l *LateEs6ToEs3Converter) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_ASSIGN:
		assign := ast.MustAssign(n)
		// Find whether this script contains the `$jscomp.createTemplateTagFirstArgWithRaw =
		// function(..) {..}` node. If yes, update the templateLitInsertionPoint.
		lhs, isGetProp := assign.Target().AsGetProp()
		if isGetProp && assign.Value().Token() == rhino.Token_FUNCTION && lhs.Object().Token() == rhino.Token_NAME {
			if lateEs6CreateTemplateTagFirstArgWithRaw.Matches(lhs.Node()) {
				exprStmt, ok := ast.AsExprStmt(n.GetParent())
				jutil.CheckStateMsg1(ok, "%s", n)
				l.templateLitInsertionPoint = exprStmt.Node().GetNext()
			}
		}
	case rhino.Token_OBJECTLIT:
		obj := ast.MustObjectLit(n)
		l.VisitObject(t, obj)
	case rhino.Token_MEMBER_FUNCTION_DEF:
		if member, ok := ast.AsMemberFunctionDef(n); ok && parent.IsObjectLit() {
			l.VisitMemberFunctionDefInObjectLit(member)
		}
	case rhino.Token_TAGGED_TEMPLATELIT:
		l.templateLiteralConverter.VisitTaggedTemplateLiteral(t, n, l.templateLitInsertionPoint)
	case rhino.Token_TEMPLATELIT:
		if !parent.IsTaggedTemplateLit() {
			l.templateLiteralConverter.VisitTemplateLiteral(t, n)
		}
	default:
	}
}

// VisitMemberFunctionDefInObjectLit converts a member definition in an
// object literal to an ES3 key/value pair. Member definitions in classes are
// handled in Es6RewriteClass.
//
// Java: LateEs6ToEs3Converter.visitMemberFunctionDefInObjectLit(Node)
func (l *LateEs6ToEs3Converter) VisitMemberFunctionDefInObjectLit(n ast.MemberFunctionDef) {
	name := n.Key()
	fn := n.Function()
	// The NAME child of the function: empty for a method, but at the
	// position of the method's name.
	nameNode := fn.Node().GetFirstChild()
	stringKey := l.astFactory.CreateStringKey(name, fn.Detach().Node())
	stringKey.SetJSDocInfo(n.Node().GetJSDocInfo())
	n.Node().ReplaceWith(stringKey)
	stringKey.Srcref(nameNode)
	l.compiler.ReportChangeToEnclosingScope(stringKey)
}

// Java: LateEs6ToEs3Converter.visitObject(NodeTraversal, Node)
func (l *LateEs6ToEs3Converter) VisitObject(t *NodeTraversal, obj ast.ObjectLit) {
	for member := range obj.Members() {
		if _, ok := member.AsComputedProp(); ok {
			l.VisitObjectWithComputedProperty(t, obj)
			return
		}
	}
}

// VisitObjectWithComputedProperty transpiles an object node with computed
// property, and add type information to the new nodes if this pass ran after
// type checking. For example,
//
//	var obj = {a: 1, [i++]: 2}
//	is transpiled to
//	var $jscomp$compprop0 = {};
//	var obj = ($jscomp$compprop0.a = 1, ($jscomp$compprop0[i++] = 2, $jscomp$compprop0));
//
// Note that when adding type information to the nodes, the NAME node
// $jscomp$compprop0 would always be assigned the type of the entire object
// (in the above example {a: number}). This is because we do not have
// sufficient type information during transpilation to know, for example,
// $jscomp$compprop0 has type Object{} in the expression
// $jscomp$compprop0.a = 1
//
// Java: LateEs6ToEs3Converter.visitObjectWithComputedProperty(NodeTraversal, Node)
func (l *LateEs6ToEs3Converter) VisitObjectWithComputedProperty(t *NodeTraversal, obj ast.ObjectLit) {
	props := []ast.ObjectMember{}
	objectType := AstFactory_Type_(obj.Node())

	// ReadAhead: the current member may be detached.
	for currElement := range ast.ReadAhead(obj.FirstMember()) {
		if currElement.Node().GetBooleanProp(rhino.Node_COMPUTED_PROP_GETTER) ||
			currElement.Node().GetBooleanProp(rhino.Node_COMPUTED_PROP_SETTER) {
			l.compiler.Report(
				JSError_MakeWithNode(
					currElement.Node(),
					ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT,
					"computed getter/setter in an object literal",
					"ES2015",
					""))
			return
		} else if tok := currElement.Token(); tok == rhino.Token_GETTER_DEF || tok == rhino.Token_SETTER_DEF {
			// Getters and setters stay in the object literal.
		} else {
			props = append(props, currElement.Detach())
		}
	}

	objName := l.compiler.GetUniqueIdSupplier().UniqueVarName(LateEs6ToEs3Converter_FRESH_COMP_PROP_VAR, t.GetInput())

	slices.Reverse(props) // Lists.reverse(props)
	result := l.astFactory.CreateName(objName, objectType)
	for _, member := range props {
		if computedProp, ok := member.AsComputedProp(); ok {
			value, ok := computedProp.Value()
			jutil.CheckStateMsg1(ok, "computed property without a value: %s", computedProp.Node())
			propertyExpression := computedProp.Key().Detach()
			value.Detach()
			result =
				l.astFactory.CreateComma(
					l.astFactory.CreateAssign(
						l.astFactory.CreateGetElem(
							l.astFactory.CreateName(objName, objectType), propertyExpression.Node()),
						value.Node()),
					result)
		} else {
			// Go-only: the other members are string keys (a method was
			// rewritten to a string key when it was visited, and getters and
			// setters stayed in the object literal).
			stringKey, ok := member.AsStringKey()
			jutil.CheckStateMsg1(ok, "unexpected object literal member: %s", member.Node())
			val := stringKey.Value().Detach()
			isQuotedAccess := stringKey.Node().IsQuotedStringKey()

			// The string key node itself becomes the string literal of the
			// property access.
			propdef := stringKey.Node()
			propdef.SetToken(rhino.Token_STRINGLIT)
			propdef.SetColor(colors.StandardColors_STRING)
			propdef.PutBooleanProp(rhino.Node_QUOTED_PROP, false)

			objNameNode := l.astFactory.CreateName(objName, objectType)
			var access rhino.Node
			if isQuotedAccess {
				access = l.astFactory.CreateGetElem(objNameNode, propdef)
			} else {
				access = l.astFactory.CreateGetProp(objNameNode, propdef.GetString(), AstFactory_Type_(propdef))
			}
			result = l.astFactory.CreateComma(l.astFactory.CreateAssign(access, val.Node()), result)
		}
	}

	statement := obj.Node()
	for !NodeUtil_IsStatement(statement) {
		statement = statement.GetParent()
	}

	isInsideLoopHeader := NodeUtil_IsLoopStructure(statement)

	if isInsideLoopHeader {
		placeholder := rhino.IR_Empty()
		obj.Node().ReplaceWith(placeholder)

		assignEmptyObj :=
			l.astFactory.CreateAssign(l.astFactory.CreateName(objName, objectType), obj.Node())
		result = l.astFactory.CreateComma(assignEmptyObj, result)
		result.SrcrefTreeIfMissing(obj.Node())

		placeholder.ReplaceWith(result)

		var_ := rhino.IR_Var(l.astFactory.CreateName(objName, objectType))
		var_.SrcrefTreeIfMissing(statement)
		var_.InsertBefore(statement)
		l.compiler.ReportChangeToEnclosingScope(var_)
	} else {
		result.SrcrefTreeIfMissing(obj.Node())
		obj.Node().ReplaceWith(result)

		var_ := rhino.IR_VarWithValue(l.astFactory.CreateName(objName, objectType), obj.Node())
		var_.SrcrefTreeIfMissing(statement)
		var_.InsertBefore(statement)
		l.compiler.ReportChangeToEnclosingScope(var_)
	}
}
