// Port of com/google/javascript/jscomp/ChromePass.java @ closure-compiler v20261028

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var (
	chromePassVisitTokens       = makeVisitTokenSet(rhino.Token_CALL)
	chromePassRenameVisitTokens = makeVisitTokenSet(rhino.Token_FUNCTION, rhino.Token_CLASS, rhino.Token_NAME)

	_ nodeTraversal_visitTokens    = (*ChromePass)(nil)
	_ nodeTraversal_noVisitLeaves  = (*ChromePass)(nil)
	_ nodeTraversal_visitTokens    = (*chromePass_RenameInternalsToExternalsCallback)(nil)
	_ nodeTraversal_onlyNameLeaves = (*chromePass_RenameInternalsToExternalsCallback)(nil)
)

func (*ChromePass) noVisitLeaves()                                     {}
func (*ChromePass) visitTokens() *ast.TokenSet                         { return &chromePassVisitTokens }
func (*chromePass_RenameInternalsToExternalsCallback) onlyNameLeaves() {}
func (*chromePass_RenameInternalsToExternalsCallback) visitTokens() *ast.TokenSet {
	return &chromePassRenameVisitTokens
}

// ChromePass is the compiler pass for Chrome-specific needs. It handles the
// following Chrome JS features:
//
//   - namespace declaration using cr.define(),
//   - unquoted property declaration using {cr|Object}.defineProperty().
//
// For the details, see tests inside ChromePassTest.java.
//
// Java: public class ChromePass extends AbstractPostOrderCallback implements CompilerPass
type ChromePass struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler AbstractCompiler

	createdObjects map[string]struct{}
}

const chromePass_CR_DEFINE = "cr.define"
const chromePass_OBJECT_DEFINE_PROPERTY = "Object.defineProperty"
const chromePass_CR_DEFINE_PROPERTY = "cr.defineProperty"
const chromePass_VIRTUAL_FILE = "<ChromePass.java>"

var chromePassVirtualFile = SourceFile_FromCodeWithKind(chromePass_VIRTUAL_FILE, "", rhino.StaticSourceFile_SourceKind_EXTERN)
var chromePass_CR_PROPERTYKIND_JS = rhino.QualifiedName_Of("cr.PropertyKind.JS")
var chromePass_CR_PROPERTYKIND_ATTR = rhino.QualifiedName_Of("cr.PropertyKind.ATTR")
var chromePass_CR_PROPERTYKIND_BOOL_ATTR = rhino.QualifiedName_Of("cr.PropertyKind.BOOL_ATTR")

const chromePass_CR_DEFINE_COMMON_EXPLANATION = "It should be called like this:" +
	" cr.define('name.space', function() '{ ... return {Export: Internal}; }');"

var ChromePass_CR_DEFINE_WRONG_NUMBER_OF_ARGUMENTS = DiagnosticType_Error(
	"JSC_CR_DEFINE_WRONG_NUMBER_OF_ARGUMENTS",
	"cr.define() should have exactly 2 arguments. "+chromePass_CR_DEFINE_COMMON_EXPLANATION)

var ChromePass_CR_DEFINE_INVALID_FIRST_ARGUMENT = DiagnosticType_Error(
	"JSC_CR_DEFINE_INVALID_FIRST_ARGUMENT",
	"Invalid first argument for cr.define(). "+chromePass_CR_DEFINE_COMMON_EXPLANATION)

var ChromePass_CR_DEFINE_INVALID_SECOND_ARGUMENT = DiagnosticType_Error(
	"JSC_CR_DEFINE_INVALID_SECOND_ARGUMENT",
	"Invalid second argument for cr.define(). "+chromePass_CR_DEFINE_COMMON_EXPLANATION)

var ChromePass_CR_DEFINE_INVALID_RETURN_IN_FUNCTION = DiagnosticType_Error(
	"JSC_CR_DEFINE_INVALID_RETURN_IN_SECOND_ARGUMENT",
	"Function passed as second argument of cr.define() should return the"+
		" dictionary in its last statement. "+
		chromePass_CR_DEFINE_COMMON_EXPLANATION)

var ChromePass_CR_DEFINE_PROPERTY_TOO_FEW_ARGUMENTS = DiagnosticType_Error(
	"JSC_CR_DEFINE_PROPERTY_TOO_FEW_ARGUMENTS",
	"cr.defineProperty() requires at least 2 arguments.")

var ChromePass_CR_DEFINE_PROPERTY_INVALID_PROPERTY_KIND = DiagnosticType_Error(
	"JSC_CR_DEFINE_PROPERTY_INVALID_PROPERTY_KIND",
	"Invalid cr.PropertyKind passed to cr.defineProperty(): expected ATTR,"+
		" BOOL_ATTR or JS, found \"{0}\".")

// Java: ChromePass(AbstractCompiler)
func NewChromePass(compiler AbstractCompiler) *ChromePass {
	return &ChromePass{
		compiler: compiler,
		// The global variable "cr" is declared in ui/webui/resources/js/cr.js.
		createdObjects: map[string]struct{}{"cr": {}},
	}
}

// Java: ChromePass.process(Node, Node)
func (p *ChromePass) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(p.compiler, root, p)
}

// Java: ChromePass.visit(NodeTraversal, Node, Node)
func (p *ChromePass) Visit(t *NodeTraversal, node rhino.Node, parent rhino.Node) {
	call := ast.MustCall(node)
	callee := call.Callee().Node()
	if callee.MatchesQualifiedName(chromePass_CR_DEFINE) {
		p.visitNamespaceDefinition(call, parent)
	} else if callee.MatchesQualifiedName(chromePass_OBJECT_DEFINE_PROPERTY) ||
		callee.MatchesQualifiedName(chromePass_CR_DEFINE_PROPERTY) {
		p.visitPropertyDefinition(call, parent)
	}
}

// Java: ChromePass.visitPropertyDefinition(Node, Node)
func (p *ChromePass) visitPropertyDefinition(call ast.Call, parent rhino.Node) {
	callNode := call.Node()
	isCrDefinePropertyCall := call.Callee().Node().MatchesQualifiedName(chromePass_CR_DEFINE_PROPERTY)

	targetArg, hasTarget := call.FirstArg()
	var propArg, kindArg ast.Arg
	var hasProp, hasKind bool
	if hasTarget {
		propArg, hasProp = targetArg.Next()
		if hasProp {
			kindArg, hasKind = propArg.Next()
		}
	}
	if isCrDefinePropertyCall {
		if !hasProp {
			p.compiler.Report(JSError_MakeWithNode(callNode, ChromePass_CR_DEFINE_PROPERTY_TOO_FEW_ARGUMENTS))
			return
		}
	} else if !hasKind {
		// Note: Object.defineProperty() with less than 3 args throws at runtime.
		// We could also report some type of error here, but it seems odd to do this from the
		// ChromePass as there's nothing particularly Chrome-specific about Object.defineProperty()
		// with too few args.
		return
	}

	// Java string concatenation renders a null qualified name as "null".
	target := jutil.ToJavaString(targetArg.Node().GetQualifiedName())
	if isCrDefinePropertyCall && !strings.HasSuffix(target, ".prototype") {
		target += ".prototype"
	}

	getPropNode := NodeUtil_NewQName(p.compiler, target+"."+propArg.Node().GetString()).SrcrefTree(callNode)

	if isCrDefinePropertyCall {
		// The 3rd argument (PropertyKind) is actually used at runtime, so it takes precedence.
		if propertyType := p.getTypeByCrPropertyKind(kindArg.Node()); !propertyType.IsNil() {
			chromePass_setJsDocWithType(getPropNode, propertyType)
		} else if sourceJsDocInfo := callNode.GetJSDocInfo(); sourceJsDocInfo != nil && sourceJsDocInfo.HasType() {
			// Otherwise, if there's a @type above the cr.defineProperty() call, move it to the getter.
			getPropNode.SetJSDocInfo(sourceJsDocInfo)
		} else {
			// Else, just give it a @type {?}.
			chromePass_setJsDocWithType(getPropNode, rhino.NewNode(rhino.Token_QMARK))
		}
		// Nuke the JsDoc info above cr.defineProperty() calls so that the CheckJSDocInfo pass won't
		// report any warnings.
		callNode.SetJSDocInfo(nil)
	} else {
		chromePass_setJsDocWithType(getPropNode, rhino.NewNode(rhino.Token_QMARK))
	}

	rhino.IR_ExprResult(getPropNode).Srcref(parent).InsertAfter(parent)
	if isCrDefinePropertyCall {
		p.compiler.ReportChangeToEnclosingScope(callNode)
	} else {
		p.compiler.ReportChangeToEnclosingScope(getPropNode)
	}
}

// Java: ChromePass.getTypeByCrPropertyKind(Node)
func (p *ChromePass) getTypeByCrPropertyKind(propertyKind rhino.Node) rhino.Node {
	if propertyKind.IsNil() || chromePass_CR_PROPERTYKIND_JS.Matches(propertyKind) {
		// This is valid, it just doesn't tell us much about the type.
		return rhino.NilNode
	}
	if chromePass_CR_PROPERTYKIND_ATTR.Matches(propertyKind) {
		return rhino.IR_String("string")
	}
	if chromePass_CR_PROPERTYKIND_BOOL_ATTR.Matches(propertyKind) {
		return rhino.IR_String("boolean")
	}
	p.compiler.Report(
		JSError_MakeWithNode(
			propertyKind,
			ChromePass_CR_DEFINE_PROPERTY_INVALID_PROPERTY_KIND,
			jutil.ToJavaString(propertyKind.GetQualifiedName())))
	return rhino.NilNode
}

// Java: ChromePass.setJsDocWithType(Node, Node)
func chromePass_setJsDocWithType(target rhino.Node, type_ rhino.Node) {
	builder := rhino.JSDocInfo_Builder_()
	builder.RecordType(rhino.NewJSTypeExpression(type_.SetStaticSourceFile(chromePassVirtualFile), chromePass_VIRTUAL_FILE))
	target.SetJSDocInfo(builder.Build())
}

// Java: ChromePass.createAndInsertObjectsForQualifiedName(Node, String)
func (p *ChromePass) createAndInsertObjectsForQualifiedName(scriptChild rhino.Node, namespace string) {
	objectsForQualifiedName := p.createObjectsForQualifiedName(namespace)
	for _, n := range objectsForQualifiedName {
		n.SrcrefTree(scriptChild).InsertBefore(scriptChild)
	}
	if len(objectsForQualifiedName) > 0 {
		p.compiler.ReportChangeToEnclosingScope(scriptChild)
	}
}

// Java: ChromePass.visitNamespaceDefinition(Node, Node)
func (p *ChromePass) visitNamespaceDefinition(call ast.Call, parent rhino.Node) {
	if call.Node().GetChildCount() != 3 {
		p.compiler.Report(JSError_MakeWithNode(call.Node(), ChromePass_CR_DEFINE_WRONG_NUMBER_OF_ARGUMENTS))
	}

	nsArg, hasNs := call.FirstArg()
	var fnArg ast.Arg
	if hasNs {
		fnArg, _ = nsArg.Next()
	}
	namespaceArg := nsArg.Node()
	nsLit, ok := ast.AsStringLit(namespaceArg)
	if !ok {
		p.compiler.Report(JSError_MakeWithNode(namespaceArg, ChromePass_CR_DEFINE_INVALID_FIRST_ARGUMENT))
		return
	}

	// TODO(vitalyp): Check namespace name for validity here. It should be a valid chain of
	// identifiers.
	namespace := nsLit.Value()
	p.createAndInsertObjectsForQualifiedName(parent, namespace)

	fn, ok := ast.AsFunction(fnArg.Node())
	if !ok {
		p.compiler.Report(JSError_MakeWithNode(namespaceArg, ChromePass_CR_DEFINE_INVALID_SECOND_ARGUMENT))
		return
	}

	functionBlock := fn.Body().Node()
	ret, ok := ast.AsReturn(functionBlock.GetLastChild())
	if !ok {
		p.compiler.Report(JSError_MakeWithNode(namespaceArg, ChromePass_CR_DEFINE_INVALID_RETURN_IN_FUNCTION))
		return
	}
	retVal, hasVal := ret.Value()
	if !hasVal {
		p.compiler.Report(JSError_MakeWithNode(namespaceArg, ChromePass_CR_DEFINE_INVALID_RETURN_IN_FUNCTION))
		return
	}
	objectLit, ok := retVal.AsObjectLit()
	if !ok {
		p.compiler.Report(JSError_MakeWithNode(namespaceArg, ChromePass_CR_DEFINE_INVALID_RETURN_IN_FUNCTION))
		return
	}

	NodeTraversal_Traverse(
		p.compiler,
		functionBlock,
		&chromePass_RenameInternalsToExternalsCallback{
			outer:          p,
			namespaceName:  namespace,
			exports:        chromePass_objectLitToMap(objectLit),
			namespaceBlock: functionBlock,
		})
}

// Java: ChromePass.objectLitToMap(Node)
func chromePass_objectLitToMap(objectLit ast.ObjectLit) map[string]string {
	res := make(map[string]string)
	for member := range objectLit.Members() {
		keyNode := member.Node()
		if val, ok := ast.AsName(keyNode.GetFirstChild()); ok {
			res[val.Ident()] = keyNode.GetString()
		}
	}
	return res
}

// createObjectsForQualifiedName: for a string "a.b.c" produce the following
// JS IR:
//
//	var a = a || {};
//	a.b = a.b || {};
//	a.b.c = a.b.c || {};
//
// Java: ChromePass.createObjectsForQualifiedName(String)
func (p *ChromePass) createObjectsForQualifiedName(namespace string) []rhino.Node {
	var objects []rhino.Node
	// namespace.split("\\.")
	parts := jutil.StringSplitLiteral(namespace, ".")

	objects = p.createObjectIfNew(objects, parts[0], true)

	if len(parts) >= 2 {
		var currPrefix strings.Builder
		currPrefix.WriteString(parts[0])
		for i := 1; i < len(parts); i++ {
			currPrefix.WriteByte('.')
			currPrefix.WriteString(parts[i])
			objects = p.createObjectIfNew(objects, currPrefix.String(), false)
		}
	}

	return objects
}

// Java: ChromePass.createObjectIfNew(List<Node>, String, boolean)
func (p *ChromePass) createObjectIfNew(objects []rhino.Node, name string, needVar bool) []rhino.Node {
	if _, exists := p.createdObjects[name]; exists {
		return objects
	}
	p.createdObjects[name] = struct{}{}

	v := NodeUtil_NewQName(p.compiler, name)
	lhs := rhino.IR_Or(v.CloneTree(), rhino.IR_Objectlit())
	if needVar {
		return append(objects, rhino.IR_VarWithValue(v, lhs))
	}
	return append(objects, rhino.IR_ExprResult(rhino.IR_Assign(v, lhs)))
}

// Java: private class ChromePass.RenameInternalsToExternalsCallback extends AbstractPostOrderCallback
type chromePass_RenameInternalsToExternalsCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	outer          *ChromePass
	namespaceName  string
	exports        map[string]string
	namespaceBlock rhino.Node
}

// Java: ChromePass.RenameInternalsToExternalsCallback.visit(NodeTraversal, Node, Node)
func (c *chromePass_RenameInternalsToExternalsCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	compiler := c.outer.compiler
	if (n.IsFunction() || n.IsClass()) && parent == c.namespaceBlock {
		nameNode := n.GetFirstChild()
		if _, exported := c.exports[nameNode.GetString()]; exported {
			clone := n.CloneTree()
			if cls, ok := ast.AsClass(clone); ok {
				cls.RemoveName()
			}
			NodeUtil_MarkNewScopesChanged(clone, compiler)
			assign := rhino.IR_Assign(c.buildQualifiedName(nameNode), clone).Srcref(n)
			if jsdoc := n.GetJSDocInfo(); jsdoc != nil {
				assign.SetJSDocInfo(jsdoc)
				clone.SetJSDocInfo(nil)
			}
			exprResult := rhino.IR_ExprResult(assign).Srcref(n)
			n.ReplaceWith(exprResult)
			NodeUtil_MarkFunctionsDeleted(n, compiler)
			compiler.ReportChangeToEnclosingScope(exprResult)
			return
		}
	}
	if name, ok := ast.AsName(n); ok && !parent.IsFunction() && !parent.IsClass() {
		if _, exported := c.exports[name.Ident()]; exported {
			if NodeUtil_IsNameDeclaration(parent) {
				if parent.GetParent() == c.namespaceBlock {
					varContent := n.RemoveFirstChild()
					inner := c.buildQualifiedName(n)
					if !varContent.IsNil() {
						inner = rhino.IR_Assign(inner, varContent).Srcref(parent)
					}
					if jsdoc := parent.GetJSDocInfo(); jsdoc != nil {
						inner.SetJSDocInfo(jsdoc.Clone())
					}
					exprResult := rhino.IR_ExprResult(inner).Srcref(parent)
					parent.ReplaceWith(exprResult)
					compiler.ReportChangeToEnclosingScope(exprResult)
				}
			} else {
				// It's a local name referencing exported entity. Change to its global name.
				newNode := c.buildQualifiedName(n)
				if jsdoc := n.GetJSDocInfo(); jsdoc != nil {
					newNode.SetJSDocInfo(jsdoc.Clone())
				}
				// If we alter the name of a called function, then it gets an explicit "this" value.
				if call, ok := ast.AsCall(parent); ok {
					call.SetFreeCall(false)
				}
				n.ReplaceWith(newNode)
				compiler.ReportChangeToEnclosingScope(newNode)
			}
		}
	}
}

// Java: ChromePass.RenameInternalsToExternalsCallback.buildQualifiedName(Node)
func (c *chromePass_RenameInternalsToExternalsCallback) buildQualifiedName(internalName rhino.Node) rhino.Node {
	// Java string concatenation renders a missing export (null) as "null".
	externalName, ok := c.exports[internalName.GetString()]
	if !ok {
		externalName = "null"
	}
	return NodeUtil_NewQName(c.outer.compiler, c.namespaceName+"."+externalName).
		SrcrefTree(internalName)
}
