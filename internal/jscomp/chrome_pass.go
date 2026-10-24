// Port of com/google/javascript/jscomp/ChromePass.java @ closure-compiler v20261015

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

	createdObjects *jutil.LinkedHashSet[string]
}

const chromePass_CR_DEFINE = "cr.define"
const chromePass_OBJECT_DEFINE_PROPERTY = "Object.defineProperty"
const chromePass_CR_DEFINE_PROPERTY = "cr.defineProperty"
const chromePass_VIRTUAL_FILE = "<ChromePass.java>"

var chromePass_VIRTUAL_NODE = rhino.IR_Empty().SetStaticSourceFile(SourceFile_FromCodeWithKind(chromePass_VIRTUAL_FILE, "", rhino.StaticSourceFile_SourceKind_EXTERN))
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
	p := &ChromePass{}
	p.compiler = compiler
	// The global variable "cr" is declared in ui/webui/resources/js/cr.js.
	p.createdObjects = jutil.NewLinkedHashSet("cr")
	return p
}

// Java: ChromePass.process(Node, Node)
func (p *ChromePass) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(p.compiler, root, p)
}

// Java: ChromePass.visit(NodeTraversal, Node, Node)
func (p *ChromePass) Visit(t *NodeTraversal, node rhino.Node, parent rhino.Node) {
	if node.IsCall() {
		callee := node.GetFirstChild()
		if callee.MatchesQualifiedName(chromePass_CR_DEFINE) {
			p.visitNamespaceDefinition(node, parent)
		} else if callee.MatchesQualifiedName(chromePass_OBJECT_DEFINE_PROPERTY) ||
			callee.MatchesQualifiedName(chromePass_CR_DEFINE_PROPERTY) {
			p.visitPropertyDefinition(node, parent)
		}
	}
}

// Java: ChromePass.visitPropertyDefinition(Node, Node)
func (p *ChromePass) visitPropertyDefinition(call rhino.Node, parent rhino.Node) {
	callee := call.GetFirstChild()
	isCrDefinePropertyCall := callee.MatchesQualifiedName(chromePass_CR_DEFINE_PROPERTY)

	if isCrDefinePropertyCall {
		if call.GetChildCount() < 3 {
			p.compiler.Report(JSError_MakeWithNode(call, ChromePass_CR_DEFINE_PROPERTY_TOO_FEW_ARGUMENTS))
			return
		}
	} else if call.GetChildCount() < 4 {
		// Note: Object.defineProperty() with less than 3 args throws at runtime.
		// We could also report some type of error here, but it seems odd to do this from the
		// ChromePass as there's nothing particularly Chrome-specific about Object.defineProperty()
		// with too few args.
		return
	}

	// Java string concatenation renders a null qualified name as "null".
	target := jutil.ToJavaString(call.GetSecondChild().GetQualifiedName())
	if isCrDefinePropertyCall && !jutil.StringEndsWith(target, ".prototype") {
		target += ".prototype"
	}

	property := call.GetChildAtIndex(2)
	getPropNode :=
		NodeUtil_NewQName(p.compiler, target+"."+property.GetString()).SrcrefTree(call)

	if isCrDefinePropertyCall {
		// The 3rd argument (PropertyKind) is actually used at runtime, so it takes precedence.
		propertyType := p.getTypeByCrPropertyKind(call.GetChildAtIndex(3))
		if !propertyType.IsNil() {
			chromePass_setJsDocWithType(getPropNode, propertyType)
		} else {
			// Otherwise, if there's a @type above the cr.defineProperty() call, move it to the getter.
			sourceJsDocInfo := call.GetJSDocInfo()
			if sourceJsDocInfo != nil && sourceJsDocInfo.HasType() {
				getPropNode.SetJSDocInfo(sourceJsDocInfo)
			} else {
				// Else, just give it a @type {?}.
				chromePass_setJsDocWithType(getPropNode, rhino.NewNode(rhino.Token_QMARK))
			}
		}
		// Nuke the JsDoc info above cr.defineProperty() calls so that the CheckJSDocInfo pass won't
		// report any warnings.
		call.SetJSDocInfo(nil)
	} else {
		chromePass_setJsDocWithType(getPropNode, rhino.NewNode(rhino.Token_QMARK))
	}

	definitionNode := rhino.IR_ExprResult(getPropNode).Srcref(parent)

	definitionNode.InsertAfter(parent)
	if isCrDefinePropertyCall {
		p.compiler.ReportChangeToEnclosingScope(call)
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
	builder.RecordType(rhino.NewJSTypeExpression(type_.SrcrefTree(chromePass_VIRTUAL_NODE), chromePass_VIRTUAL_FILE))
	target.SetJSDocInfo(builder.Build())
}

// Java: ChromePass.createAndInsertObjectsForQualifiedName(Node, String)
func (p *ChromePass) createAndInsertObjectsForQualifiedName(scriptChild rhino.Node, namespace string) {
	objectsForQualifiedName := p.createObjectsForQualifiedName(namespace)
	for n := range objectsForQualifiedName.All() {
		n.SrcrefTree(scriptChild).InsertBefore(scriptChild)
	}
	if !objectsForQualifiedName.IsEmpty() {
		p.compiler.ReportChangeToEnclosingScope(scriptChild)
	}
}

// Java: ChromePass.visitNamespaceDefinition(Node, Node)
func (p *ChromePass) visitNamespaceDefinition(crDefineCallNode rhino.Node, parent rhino.Node) {
	if crDefineCallNode.GetChildCount() != 3 {
		p.compiler.Report(JSError_MakeWithNode(crDefineCallNode, ChromePass_CR_DEFINE_WRONG_NUMBER_OF_ARGUMENTS))
	}

	namespaceArg := crDefineCallNode.GetSecondChild()
	function := crDefineCallNode.GetChildAtIndex(2)

	if !namespaceArg.IsStringLit() {
		p.compiler.Report(JSError_MakeWithNode(namespaceArg, ChromePass_CR_DEFINE_INVALID_FIRST_ARGUMENT))
		return
	}

	// TODO(vitalyp): Check namespace name for validity here. It should be a valid chain of
	// identifiers.
	namespace := namespaceArg.GetString()

	p.createAndInsertObjectsForQualifiedName(parent, namespace)

	if !function.IsFunction() {
		p.compiler.Report(JSError_MakeWithNode(namespaceArg, ChromePass_CR_DEFINE_INVALID_SECOND_ARGUMENT))
		return
	}

	var returnNode rhino.Node
	var objectLit rhino.Node
	functionBlock := function.GetLastChild()
	returnNode = functionBlock.GetLastChild()
	if returnNode.IsNil() || !returnNode.IsReturn() {
		p.compiler.Report(JSError_MakeWithNode(namespaceArg, ChromePass_CR_DEFINE_INVALID_RETURN_IN_FUNCTION))
		return
	}
	objectLit = returnNode.GetFirstChild()
	if objectLit.IsNil() || !objectLit.IsObjectLit() {
		p.compiler.Report(JSError_MakeWithNode(namespaceArg, ChromePass_CR_DEFINE_INVALID_RETURN_IN_FUNCTION))
		return
	}

	exports := chromePass_objectLitToMap(objectLit)

	NodeTraversal_Traverse(
		p.compiler,
		functionBlock,
		newChromePass_RenameInternalsToExternalsCallback(p, namespace, exports, functionBlock))
}

// Java: ChromePass.objectLitToMap(Node)
func chromePass_objectLitToMap(objectLit rhino.Node) *jutil.LinkedHashMap[string, string] {
	res := jutil.NewLinkedHashMap[string, string]()

	for keyNode := objectLit.GetFirstChild(); !keyNode.IsNil(); keyNode = keyNode.GetNext() {
		key := keyNode.GetString()

		valueNode := keyNode.GetFirstChild()
		if valueNode.IsName() {
			value := keyNode.GetFirstChild().GetString()
			res.Put(value, key)
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
func (p *ChromePass) createObjectsForQualifiedName(namespace string) *jutil.ArrayList[rhino.Node] {
	objects := jutil.NewArrayList[rhino.Node]()
	// namespace.split("\\.")
	parts := jutil.StringSplitLiteral(namespace, ".")

	p.createObjectIfNew(objects, parts[0], true)

	if len(parts) >= 2 {
		currPrefix := jutil.NewStringBuilder().Append(parts[0])
		for i := 1; i < len(parts); i++ {
			currPrefix.Append(".").Append(parts[i])
			p.createObjectIfNew(objects, currPrefix.String(), false)
		}
	}

	return objects
}

// Java: ChromePass.createObjectIfNew(List<Node>, String, boolean)
func (p *ChromePass) createObjectIfNew(objects *jutil.ArrayList[rhino.Node], name string, needVar bool) {
	if p.createdObjects.Contains(name) {
		return
	}
	p.createdObjects.Add(name)

	v := NodeUtil_NewQName(p.compiler, name)
	lhs := rhino.IR_Or(v.CloneTree(), rhino.IR_Objectlit())
	if needVar {
		objects.Add(rhino.IR_VarWithValue(v, lhs))
	} else {
		objects.Add(rhino.IR_ExprResult(rhino.IR_Assign(v, lhs)))
	}
}

// Java: private class ChromePass.RenameInternalsToExternalsCallback extends AbstractPostOrderCallback
type chromePass_RenameInternalsToExternalsCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	outer *ChromePass

	namespaceName  string
	exports        *jutil.LinkedHashMap[string, string]
	namespaceBlock rhino.Node
}

// Java: ChromePass.RenameInternalsToExternalsCallback(String, Map<String, String>, Node)
func newChromePass_RenameInternalsToExternalsCallback(
	outer *ChromePass,
	namespaceName string, exports *jutil.LinkedHashMap[string, string], namespaceBlock rhino.Node) *chromePass_RenameInternalsToExternalsCallback {
	return &chromePass_RenameInternalsToExternalsCallback{
		outer:          outer,
		namespaceName:  namespaceName,
		exports:        exports,
		namespaceBlock: namespaceBlock,
	}
}

// Java: ChromePass.RenameInternalsToExternalsCallback.visit(NodeTraversal, Node, Node)
func (c *chromePass_RenameInternalsToExternalsCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	compiler := c.outer.compiler
	if (n.IsFunction() || n.IsClass()) &&
		parent == c.namespaceBlock &&
		c.exports.ContainsKey(n.GetFirstChild().GetString()) {
		// It's a top-level function/constructor/class definition.
		//
		// Change
		//
		//   /** Some doc */
		//   function internalName() {}
		//
		// to
		//
		//   /** Some doc */
		//   my.namespace.name.externalName = function internalName() {};
		//
		// and change
		//
		//   /** Some doc */
		//   class InternalName {}
		//
		// to
		//
		//   /** Some doc */
		//   my.namespace.name.ExternalName = class InternalName {};
		//
		// by looking up in this.exports for internalName to find the correspondent
		// externalName.
		clone := n.CloneTree()
		if clone.IsClass() {
			className := clone.GetFirstChild()
			className.ReplaceWith(rhino.IR_Empty().Srcref(className))
		}
		NodeUtil_MarkNewScopesChanged(clone, compiler)
		exprResult :=
			rhino.IR_ExprResult(rhino.IR_Assign(c.buildQualifiedName(n.GetFirstChild()), clone).Srcref(n)).
				Srcref(n)

		if n.GetJSDocInfo() != nil {
			exprResult.GetFirstChild().SetJSDocInfo(n.GetJSDocInfo())
			clone.SetJSDocInfo(nil)
		}
		n.ReplaceWith(exprResult)
		NodeUtil_MarkFunctionsDeleted(n, compiler)
		compiler.ReportChangeToEnclosingScope(exprResult)
	} else if n.IsName() &&
		c.exports.ContainsKey(n.GetString()) &&
		!parent.IsFunction() &&
		!parent.IsClass() {
		if NodeUtil_IsNameDeclaration(parent) {
			if parent.GetParent() == c.namespaceBlock {
				// It's a top-level exported variable definition (maybe without an
				// assignment).
				// Change
				//
				//   var enum = { 'one': 1, 'two': 2 };
				//
				// to
				//
				//   my.namespace.name.enum = { 'one': 1, 'two': 2 };
				varContent := n.RemoveFirstChild()
				var exprResult rhino.Node
				if varContent.IsNil() {
					exprResult = rhino.IR_ExprResult(c.buildQualifiedName(n)).Srcref(parent)
				} else {
					exprResult =
						rhino.IR_ExprResult(rhino.IR_Assign(c.buildQualifiedName(n), varContent).Srcref(parent)).
							Srcref(parent)
				}
				if parent.GetJSDocInfo() != nil {
					exprResult.GetFirstChild().SetJSDocInfo(parent.GetJSDocInfo().Clone())
				}
				parent.ReplaceWith(exprResult)
				compiler.ReportChangeToEnclosingScope(exprResult)
			}
		} else {
			// It's a local name referencing exported entity. Change to its global name.
			newNode := c.buildQualifiedName(n)
			if n.GetJSDocInfo() != nil {
				newNode.SetJSDocInfo(n.GetJSDocInfo().Clone())
			}

			// If we alter the name of a called function, then it gets an explicit "this"
			// value.
			if parent.IsCall() {
				parent.PutBooleanProp(rhino.Node_FREE_CALL, false)
			}

			n.ReplaceWith(newNode)
			compiler.ReportChangeToEnclosingScope(newNode)
		}
	}
}

// Java: ChromePass.RenameInternalsToExternalsCallback.buildQualifiedName(Node)
func (c *chromePass_RenameInternalsToExternalsCallback) buildQualifiedName(internalName rhino.Node) rhino.Node {
	// Java string concatenation renders a missing export (null) as "null".
	externalName, ok := c.exports.Lookup(internalName.GetString())
	if !ok {
		externalName = "null"
	}
	return NodeUtil_NewQName(c.outer.compiler, c.namespaceName+"."+externalName).
		SrcrefTree(internalName)
}
