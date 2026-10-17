// Port of com/google/javascript/jscomp/CheckJSDoc.java @ closure-compiler v20260415

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckJSDoc checks for misplaced, misused or deprecated JSDoc annotations.
type CheckJSDoc struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler  AbstractCompiler
	inExterns bool

	checkJsDocTypesVisitor *CheckJSDoc_CheckJsdocTypes
}

// Java: CheckJSDoc.MISPLACED_MSG_ANNOTATION
var CheckJSDoc_MISPLACED_MSG_ANNOTATION = DiagnosticType_Disabled(
	"JSC_MISPLACED_MSG_ANNOTATION",
	"Misplaced message annotation. @desc, @meaning, and @alternateMessageId"+
		" annotations should be only on message nodes."+
		"\nMessage constants must be prefixed with 'MSG_'.")

// Java: CheckJSDoc.MISPLACED_ANNOTATION
var CheckJSDoc_MISPLACED_ANNOTATION = DiagnosticType_Warning("JSC_MISPLACED_ANNOTATION", "Misplaced {0} annotation. {1}")

// Java: CheckJSDoc.ANNOTATION_DEPRECATED
var CheckJSDoc_ANNOTATION_DEPRECATED = DiagnosticType_Warning("JSC_ANNOTATION_DEPRECATED", "The {0} annotation is deprecated. {1}")

// Java: CheckJSDoc.DISALLOWED_MEMBER_JSDOC
var CheckJSDoc_DISALLOWED_MEMBER_JSDOC = DiagnosticType_Warning(
	"JSC_DISALLOWED_MEMBER_JSDOC",
	"Class level JSDocs (@interface, @extends, etc.) are not allowed on class members")

// Java: CheckJSDoc.ARROW_FUNCTION_AS_CONSTRUCTOR
var CheckJSDoc_ARROW_FUNCTION_AS_CONSTRUCTOR = DiagnosticType_Error(
	"JSC_ARROW_FUNCTION_AS_CONSTRUCTOR", "Arrow functions cannot be used as constructors")

// Java: CheckJSDoc.BAD_REST_PARAMETER_ANNOTATION
var CheckJSDoc_BAD_REST_PARAMETER_ANNOTATION = DiagnosticType_Warning(
	"JSC_BAD_REST_PARAMETER_ANNOTATION",
	"Missing \"...\" in type annotation for rest parameter.")

// Java: CheckJSDoc.DEFAULT_PARAM_MUST_BE_MARKED_OPTIONAL
var CheckJSDoc_DEFAULT_PARAM_MUST_BE_MARKED_OPTIONAL = DiagnosticType_Error(
	"JSC_DEFAULT_PARAM_MUST_BE_MARKED_OPTIONAL",
	"Inline JSDoc on default parameters must be marked as optional")

// Java: CheckJSDoc.INVALID_NO_SIDE_EFFECT_ANNOTATION
var CheckJSDoc_INVALID_NO_SIDE_EFFECT_ANNOTATION = DiagnosticType_Error(
	"JSC_INVALID_NO_SIDE_EFFECT_ANNOTATION",
	"@nosideeffects may only appear in externs files.")

// Java: CheckJSDoc.INVALID_MODIFIES_ANNOTATION
var CheckJSDoc_INVALID_MODIFIES_ANNOTATION = DiagnosticType_Error(
	"JSC_INVALID_MODIFIES_ANNOTATION", "@modifies may only appear in externs files.")

// Java: CheckJSDoc.INVALID_DEFINE_ON_LET
var CheckJSDoc_INVALID_DEFINE_ON_LET = DiagnosticType_Error(
	"JSC_INVALID_DEFINE_ON_LET",
	"variables annotated with @define may only be declared with VARs, ASSIGNs, or CONSTs")

// Java: CheckJSDoc.MISPLACED_SUPPRESS
var CheckJSDoc_MISPLACED_SUPPRESS = DiagnosticType_Warning(
	"JSC_MISPLACED_SUPPRESS",
	"@suppress annotation not allowed here. See"+
		" https://github.com/google/closure-compiler/wiki/@suppress-annotations")

// Java: CheckJSDoc.JSDOC_IN_BLOCK_COMMENT
var CheckJSDoc_JSDOC_IN_BLOCK_COMMENT = DiagnosticType_Warning(
	"JSC_JSDOC_IN_BLOCK_COMMENT",
	"Non-JSDoc comment has annotations. Did you mean to start it with '/**'?")

// Java: CheckJSDoc.JSDOC_ON_RETURN
var CheckJSDoc_JSDOC_ON_RETURN = DiagnosticType_Warning(
	"JSC_JSDOC_ON_RETURN", "JSDoc annotations are not supported on return.")

// Java: CheckJSDoc.COMMENT_PATTERN
var checkJSDoc_COMMENT_PATTERN = jutil.Pattern_Compile("(/|(\n[ \t]*))\\*[ \t]*@[a-zA-Z]+[ \t\n{]")

// Java: CheckJSDoc(AbstractCompiler)
func NewCheckJSDoc(compiler AbstractCompiler) *CheckJSDoc {
	c := &CheckJSDoc{}
	c.checkJsDocTypesVisitor = &CheckJSDoc_CheckJsdocTypes{}
	c.compiler = compiler
	return c
}

// Java: CheckJSDoc.process(Node, Node)
func (c *CheckJSDoc) Process(externs *rhino.Node, root *rhino.Node) {
	c.inExterns = true
	NodeTraversal_Traverse(c.compiler, externs, c)
	c.inExterns = false
	NodeTraversal_Traverse(c.compiler, root, c)
}

// CheckJsDocInBlockComments checks for block comments (e.g. starting with /*)
// that look like they are JsDoc, and thus should start with /**.
//
// Java: CheckJSDoc.checkJsDocInBlockComments(String)
func (c *CheckJSDoc) CheckJsDocInBlockComments(fileName *string) {
	if !c.compiler.GetOptions().PreservesDetailedSourceInfo() {
		// Comments only available if preservesDetailedSourceInfo is true.
		return
	}

	for _, comment := range c.compiler.GetComments(*fileName) {
		if comment.Type == parser.Comment_Type_BLOCK {
			if checkJSDoc_COMMENT_PATTERN.Matcher(comment.Value).Find() {
				c.compiler.Report(
					JSError_MakeWithSourceLocation(
						fileName,
						comment.Location.Start.Line+1,
						comment.Location.Start.Column,
						CheckJSDoc_JSDOC_IN_BLOCK_COMMENT))
			}
		}
	}
}

// Java: CheckJSDoc.visit(NodeTraversal, Node, Node)
func (c *CheckJSDoc) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if n.IsScript() {
		c.CheckJsDocInBlockComments(n.GetSourceFileName())
	}
	info := n.GetJSDocInfo()
	c.ValidateTypeAnnotations(n, info)
	c.ValidateFunctionJsDoc(n, info)
	c.ValidateMsgJsDoc(n, info)
	c.ValidateNoCollapse(n, info)
	c.ValidateClassLevelJsDoc(n, info)
	c.ValidateArrowFunction(n)
	c.ValidateRestParameter(n)
	c.ValidateDefaultValue(n)
	c.ValidateTemplates(n, info)
	c.ValidateTypedefs(n, info)
	c.ValidateNoSideEffects(n, info)
	c.ValidateAbstractJsDoc(n, info)
	c.ValidateDefinesDeclaration(n, info)
	c.ValidateSuppress(n, info)
	c.ValidateImplicitCast(n, info)
	c.ValidateClosurePrimitive(n, info)
	c.ValidateReturnJsDoc(n, info)
	c.ValidateTsType(n, info)
	c.ValidateJsDocTypeNames(info)
	c.ValidateIsUsedViaDotConstructor(n, info)
}

// Java: CheckJSDoc.validateSuppress(Node, JSDocInfo)
func (c *CheckJSDoc) ValidateSuppress(n *rhino.Node, info *rhino.JSDocInfo) {
	if info == nil || info.GetSuppressions().IsEmpty() {
		return
	}
	switch n.GetToken() {
	case rhino.Token_FUNCTION,
		rhino.Token_CLASS,
		rhino.Token_VAR,
		rhino.Token_LET,
		rhino.Token_CONST,
		rhino.Token_SCRIPT,
		rhino.Token_MEMBER_FUNCTION_DEF,
		rhino.Token_GETTER_DEF,
		rhino.Token_SETTER_DEF,
		rhino.Token_MEMBER_FIELD_DEF,
		rhino.Token_COMPUTED_FIELD_DEF:
		// Suppressions are always valid here.
		return
	case rhino.Token_COMPUTED_PROP:
		if n.GetLastChild().IsFunction() {
			return // Suppressions are valid on computed properties that declare functions.
		}
	case rhino.Token_STRING_KEY:
		if n.GetParent().IsObjectLit() {
			return
		}
	case rhino.Token_WITH:
		if CheckJSDoc_ContainsOnlySuppressionFor(info, "with") {
			return
		}
	default:
	}
	if CheckJSDoc_ContainsOnlySuppressionFor(info, "missingRequire") {
		return
	}
	if n.GetParent().IsExprResult() {
		return
	}
	c.compiler.Report(JSError_MakeWithNode(n, CheckJSDoc_MISPLACED_SUPPRESS))
}

// Java: CheckJSDoc.containsOnlySuppressionFor(JSDocInfo, String)
func CheckJSDoc_ContainsOnlySuppressionFor(jsdoc *rhino.JSDocInfo, allowedSuppression string) bool {
	suppressions := jsdoc.GetSuppressions()
	if suppressions.Size() != 1 {
		return false
	}
	// Iterables.getOnlyElement(suppressions): the set holds exactly one element.
	for only := range suppressions.All() {
		return only == allowedSuppression
	}
	panic(jutil.NewNoSuchElementException())
}

// Java: CheckJSDoc.validateTypedefs(Node, JSDocInfo)
func (c *CheckJSDoc) ValidateTypedefs(n *rhino.Node, info *rhino.JSDocInfo) {
	if info == nil || !info.HasTypedefType() {
		return
	}
	if c.IsClassDecl(n) {
		c.ReportMisplaced(n, "typedef", "@typedef is not allowed on a class declaration.")
		return
	}
	lvalue := n
	if NodeUtil_IsNameDeclaration(n) || n.IsAssign() {
		lvalue = n.GetFirstChild()
	}
	if !lvalue.IsQualifiedName() {
		c.ReportMisplaced(
			n,
			"typedef",
			"@typedef is only allowed on qualified name declarations. Did you mean @type?")
	} else if CheckJSDoc_IsPrototypeOrInstanceDecl(lvalue) {
		c.ReportMisplaced(
			n,
			"typedef",
			"@typedef is not allowed on instance or prototype properties. Did you mean @type?")
	}
}

// Java: CheckJSDoc.validateTemplates(Node, JSDocInfo)
func (c *CheckJSDoc) ValidateTemplates(n *rhino.Node, info *rhino.JSDocInfo) {
	if info != nil &&
		len(info.GetTemplateTypeNames()) != 0 &&
		!info.IsConstructorOrInterface() &&
		!c.IsClassDecl(n) &&
		!info.ContainsFunctionDeclaration() &&
		CheckJSDoc_GetFunctionDecl(n) == nil {
		c.ReportMisplaced(
			n,
			"template",
			"@template is only allowed in class, constructor, interface, function "+
				"or method declarations")
	}
}

// CheckJSDoc_GetFunctionDecl returns the function node associated with the
// function declaration associated with the specified node, no null if no such
// function exists.
//
// Java: CheckJSDoc.getFunctionDecl(Node)
func CheckJSDoc_GetFunctionDecl(n *rhino.Node) *rhino.Node {
	if n.IsFunction() {
		return n
	}
	if n.IsMemberFunctionDef() {
		return n.GetFirstChild()
	}
	if NodeUtil_IsNameDeclaration(n) &&
		n.GetFirstFirstChild() != nil &&
		n.GetFirstFirstChild().IsFunction() {
		return n.GetFirstFirstChild()
	}

	if n.IsAssign() && n.GetFirstChild().IsQualifiedName() && n.GetLastChild().IsFunction() {
		return n.GetLastChild()
	}

	if n.IsGetterDef() || n.IsSetterDef() {
		return n.GetFirstChild()
	}

	if n.IsComputedProp() && n.GetLastChild().IsFunction() {
		return n.GetLastChild()
	}

	return nil
}

// Java: CheckJSDoc.isClassDecl(Node)
func (c *CheckJSDoc) IsClassDecl(n *rhino.Node) bool {
	return c.IsClass(n) ||
		(n.IsAssign() && c.IsClass(n.GetLastChild())) ||
		(NodeUtil_IsNameDeclaration(n) && c.IsNameInitializeWithClass(n.GetFirstChild())) ||
		c.IsNameInitializeWithClass(n)
}

// Java: CheckJSDoc.isNameInitializeWithClass(Node)
func (c *CheckJSDoc) IsNameInitializeWithClass(n *rhino.Node) bool {
	return n != nil && n.IsName() && n.HasChildren() && c.IsClassDecl(n.GetFirstChild())
}

// Java: CheckJSDoc.isClass(Node)
func (c *CheckJSDoc) IsClass(n *rhino.Node) bool {
	return n.IsClass() || (n.IsCall() && c.compiler.GetCodingConvention().IsClassFactoryCall(n))
}

// Java: CheckJSDoc.isPrototypeOrInstanceDecl(Node)
func CheckJSDoc_IsPrototypeOrInstanceDecl(n *rhino.Node) bool {
	if n.IsStringKey() {
		return false
	}
	if NodeUtil_IsPrototypeProperty(n) {
		return true
	}
	receiver := NodeUtil_GetRootOfQualifiedName(n)
	return receiver.IsThis() || receiver.IsSuper()
}

// ValidateClassLevelJsDoc checks that class-level annotations like
// @interface/@extends are not used on member functions.
//
// Java: CheckJSDoc.validateClassLevelJsDoc(Node, JSDocInfo)
func (c *CheckJSDoc) ValidateClassLevelJsDoc(n *rhino.Node, info *rhino.JSDocInfo) {
	if info != nil && n.IsMemberFunctionDef() && CheckJSDoc_HasClassLevelJsDoc(info) {
		c.Report(n, CheckJSDoc_DISALLOWED_MEMBER_JSDOC)
	}
}

// Java: CheckJSDoc.validateAbstractJsDoc(Node, JSDocInfo)
func (c *CheckJSDoc) ValidateAbstractJsDoc(n *rhino.Node, info *rhino.JSDocInfo) {
	if info == nil || !info.IsAbstract() {
		return
	}
	if c.IsClassDecl(n) {
		return
	}

	// @abstract annotation on a function written as a class field.
	if n.IsGetProp() {
		if !(n.GetFirstChild().IsThis() && info.GetReturnType() != nil) {
			c.Report(
				n,
				CheckJSDoc_MISPLACED_ANNOTATION,
				"@abstract",
				"abstract undeclared methods can only be written as class fields")
		} else if !(n.GetParent().IsExprResult() &&
			n.GetGrandparent().IsBlock() &&
			n.GetParent().GetGrandparent().IsFunction()) {
			c.Report(
				n,
				CheckJSDoc_MISPLACED_ANNOTATION,
				"@abstract",
				"abstract methods without an initializer must be declared in a constructor function.")
			return
		}
		return
	}

	// @abstract annotation on a non-function
	functionNode := CheckJSDoc_GetFunctionDecl(n)
	if functionNode == nil {
		c.Report(
			n,
			CheckJSDoc_MISPLACED_ANNOTATION,
			"@abstract",
			"only functions or non-static methods can be abstract")
		return
	}

	if !info.IsConstructor() && NodeUtil_GetFunctionBody(functionNode).HasChildren() {
		// @abstract annotation on a function with a non-empty body
		c.Report(
			n,
			CheckJSDoc_MISPLACED_ANNOTATION,
			"@abstract",
			"function with a non-empty body cannot be abstract")
		return
	}

	if NodeUtil_IsEs6ConstructorMemberFunctionDef(n) {
		// @abstract annotation on an ES6 constructor
		c.Report(n, CheckJSDoc_MISPLACED_ANNOTATION, "@abstract", "constructors cannot be abstract")
		return
	}

	if !info.IsConstructor() &&
		!n.IsMemberFunctionDef() &&
		!n.IsStringKey() &&
		!n.IsComputedProp() &&
		!n.IsGetterDef() &&
		!n.IsSetterDef() &&
		!NodeUtil_IsPrototypeMethod(functionNode) {
		// @abstract annotation on a non-method (or static method) in ES5
		c.Report(
			n,
			CheckJSDoc_MISPLACED_ANNOTATION,
			"@abstract",
			"only functions or non-static methods can be abstract")
		return
	}

	if n.IsStaticMember() {
		// @abstract annotation on a static method in ES6
		c.Report(n, CheckJSDoc_MISPLACED_ANNOTATION, "@abstract", "static methods cannot be abstract")
		return
	}
}

// Java: CheckJSDoc.hasClassLevelJsDoc(JSDocInfo)
func CheckJSDoc_HasClassLevelJsDoc(info *rhino.JSDocInfo) bool {
	return info.IsConstructorOrInterface() ||
		info.HasBaseType() ||
		info.GetImplementedInterfaceCount() != 0 ||
		info.GetExtendedInterfacesCount() != 0
}

// ValidateNoCollapse warns when nocollapse annotations are present on nodes
// which are not eligible for property collapsing.
//
// Java: CheckJSDoc.validateNoCollapse(Node, JSDocInfo)
func (c *CheckJSDoc) ValidateNoCollapse(n *rhino.Node, info *rhino.JSDocInfo) {
	if info == nil || !info.IsNoCollapse() {
		return
	}
	if n.IsFromExterns() {
		// @nocollapse has no effect in externs
		c.ReportMisplaced(n, "nocollapse", "This JSDoc has no effect in externs.")
		return
	}
	if NodeUtil_IsPrototypePropertyDeclaration(n.GetParent()) ||
		(n.GetParent().IsClassMembers() && !n.IsStaticMember()) {
		c.ReportMisplaced(
			n,
			"nocollapse",
			"This JSDoc has no effect on prototype properties and non-static fields.")
	}
	if n.IsAssign() {
		assignee := n.GetFirstChild()
		if assignee.IsQualifiedName() {
			rootOfQname := NodeUtil_GetRootOfQualifiedName(assignee)
			if !rootOfQname.IsName() {
				c.ReportMisplaced(n, "nocollapse", "This JSDoc has no effect.")
			}
		}
	}
}

// ValidateFunctionJsDoc checks that JSDoc intended for a function is actually
// attached to a function.
//
// Java: CheckJSDoc.validateFunctionJsDoc(Node, JSDocInfo)
func (c *CheckJSDoc) ValidateFunctionJsDoc(n *rhino.Node, info *rhino.JSDocInfo) {
	if info == nil {
		return
	}

	if info.ContainsFunctionDeclaration() && !info.HasType() && !c.IsJSDocOnFunctionNode(n, info) {
		// This JSDoc should be attached to a FUNCTION node, or an assignment
		// with a function as the RHS, etc.

		c.ReportMisplaced(
			n,
			"function",
			"This JSDoc is not attached to a function node. "+"Are you missing parentheses?")
	}
}

// IsJSDocOnFunctionNode reports whether this node's JSDoc may apply to a
// function.
//
// This has some false positive cases, to allow for patterns like
// goog.abstractMethod.
//
// Java: CheckJSDoc.isJSDocOnFunctionNode(Node, JSDocInfo)
func (c *CheckJSDoc) IsJSDocOnFunctionNode(n *rhino.Node, info *rhino.JSDocInfo) bool {
	switch n.GetToken() {
	case rhino.Token_FUNCTION,
		rhino.Token_GETTER_DEF,
		rhino.Token_SETTER_DEF,
		rhino.Token_MEMBER_FUNCTION_DEF,
		rhino.Token_STRING_KEY,
		rhino.Token_COMPUTED_PROP,
		rhino.Token_EXPORT,
		rhino.Token_MEMBER_FIELD_DEF,
		rhino.Token_COMPUTED_FIELD_DEF:
		return true
	case rhino.Token_GETELEM, rhino.Token_GETPROP:
		if n.GetFirstChild().IsQualifiedName() {
			// assume qualified names may be function declarations
			return true
		}
		return false
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST, rhino.Token_ASSIGN:
		lhs := n.GetFirstChild()
		rhs := NodeUtil_GetRValueOfLValue(lhs)
		if rhs != nil && c.IsClass(rhs) && !info.IsConstructor() {
			return false
		}

		// TODO(b/124081098): Check that the RHS of the assignment is a
		// function. Note that it can be a FUNCTION node, but it can also be
		// a call to goog.abstractMethod, goog.functions.constant, etc.
		return true
	default:
		return false
	}
}

// ValidateMsgJsDoc checks that annotations for messages (@desc, @meaning and
// @alternateMessageId) are in the proper place, namely on names starting
// with MSG_ which indicates they should be extracted for translation. A later
// pass checks that the right side is a call to goog.getMsg.
//
// Java: CheckJSDoc.validateMsgJsDoc(Node, JSDocInfo)
func (c *CheckJSDoc) ValidateMsgJsDoc(n *rhino.Node, info *rhino.JSDocInfo) {
	if info == nil {
		return
	}

	hasNonDescMsgTag := info.GetMeaning() != nil || info.GetAlternateMessageId() != nil

	if hasNonDescMsgTag ||
		// Don't error on TS gencode using @desc on a non-message. There's a lot of code that
		// uses @desc as a general purpose "@desc" tag
		(info.GetDescription() != nil && !CheckJSDoc_IsFromTs(n)) {
		descOkay := false
		switch n.GetToken() {
		case rhino.Token_ASSIGN, rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
			descOkay = CheckJSDoc_IsValidMsgName(n.GetFirstChild())
		case rhino.Token_STRING_KEY:
			descOkay = CheckJSDoc_IsValidMsgName(n)
		case rhino.Token_GETPROP:
			if n.IsFromExterns() && n.IsQualifiedName() {
				descOkay = CheckJSDoc_IsValidMsgName(n)
			}
		default:
		}
		if !descOkay {
			c.Report(n, CheckJSDoc_MISPLACED_MSG_ANNOTATION)
		}
	}
}

// CheckJSDoc_IsValidMsgName returns whether of not the given name is valid
// target for the result of goog.getMsg.
//
// Java: CheckJSDoc.isValidMsgName(Node)
func CheckJSDoc_IsValidMsgName(nameNode *rhino.Node) bool {
	if nameNode.IsName() || nameNode.IsStringKey() {
		return strings.HasPrefix(nameNode.GetString(), "MSG_")
	} else if nameNode.IsQualifiedName() {
		return strings.HasPrefix(nameNode.GetString(), "MSG_")
	} else {
		return false
	}
}

// ValidateTypeAnnotations checks that JSDoc with a @type annotation is in a
// valid place.
//
// Java: CheckJSDoc.validateTypeAnnotations(Node, JSDocInfo)
func (c *CheckJSDoc) ValidateTypeAnnotations(n *rhino.Node, info *rhino.JSDocInfo) {
	if info != nil && info.HasType() {
		valid := false
		switch n.GetToken() {
		case rhino.Token_FUNCTION:
			// Function declarations are valid
			valid = NodeUtil_IsFunctionDeclaration(n)
		// Object literal properties, catch declarations and variable
		// initializers are valid.
		case rhino.Token_NAME:
			valid = CheckJSDoc_IsTypeAnnotationAllowedForName(n)
		case rhino.Token_ARRAY_PATTERN, rhino.Token_OBJECT_PATTERN:
			// allow JSDoc like
			//   function f(/** !Object */ {x}) {}
			//   function f(/** !Array */ [x]) {}
			valid = n.GetParent().IsParamList()
		// Casts, exports, and Object literal properties are valid.
		case rhino.Token_CAST,
			rhino.Token_EXPORT,
			rhino.Token_STRING_KEY,
			rhino.Token_GETTER_DEF,
			rhino.Token_SETTER_DEF,
			rhino.Token_MEMBER_FIELD_DEF,
			rhino.Token_COMPUTED_FIELD_DEF:
			valid = true
		// Declarations are valid iff they only contain simple names
		//   /** @type {number} */ var x = 3; // ok
		//   /** @type {number} */ var {x} = obj; // forbidden
		case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
			valid = !NodeUtil_IsDestructuringDeclaration(n)
		// Property assignments are valid, if at the root of an expression.
		case rhino.Token_ASSIGN:
			lvalue := n.GetFirstChild()
			valid = n.GetParent().IsExprResult() &&
				(lvalue.IsGetProp() || lvalue.IsGetElem() || lvalue.MatchesName("exports"))
		case rhino.Token_GETPROP:
			valid = n.GetParent().IsExprResult() && n.IsQualifiedName()
		case rhino.Token_CALL:
			valid = info.IsDefine()
		default:
		}

		if !valid {
			c.ReportMisplaced(
				n, "type", "Type annotations are not allowed here. "+"Are you missing parentheses?")
		}
	}
}

// CheckJSDoc_IsTypeAnnotationAllowedForName: is it valid to have a type
// annotation on the given NAME node?
//
// Java: CheckJSDoc.isTypeAnnotationAllowedForName(Node)
func CheckJSDoc_IsTypeAnnotationAllowedForName(n *rhino.Node) bool {
	jutil.CheckStateMsg1(n.IsName(), "%s", n)
	// Only allow type annotations on nodes used as an lvalue.
	if !NodeUtil_IsLValue(n) {
		return false
	}
	// Don't allow JSDoc on a name in an assignment. Simple names should only have JSDoc on them
	// when originally declared.
	rootTarget := NodeUtil_GetRootTarget(n)
	return !NodeUtil_IsLhsOfAssign(rootTarget)
}

// Java: CheckJSDoc.reportMisplaced(Node, String, String)
func (c *CheckJSDoc) ReportMisplaced(n *rhino.Node, annotationName string, note string) {
	c.compiler.Report(JSError_MakeWithNode(n, CheckJSDoc_MISPLACED_ANNOTATION, annotationName, note))
}

// Java: CheckJSDoc.report(Node, DiagnosticType, String...)
func (c *CheckJSDoc) Report(n *rhino.Node, type_ *DiagnosticType, arguments ...string) {
	c.compiler.Report(JSError_MakeWithNode(n, type_, arguments...))
}

// ValidateArrowFunction checks that an arrow function is not annotated with
// {@constructor}.
//
// Java: CheckJSDoc.validateArrowFunction(Node)
func (c *CheckJSDoc) ValidateArrowFunction(n *rhino.Node) {
	if n.IsArrowFunction() {
		info := NodeUtil_GetBestJSDocInfo(n)
		if info != nil && info.IsConstructorOrInterface() {
			c.Report(n, CheckJSDoc_ARROW_FUNCTION_AS_CONSTRUCTOR)
		}
	}
}

// ValidateRestParameter checks that a rest parameter has JSDoc marked as
// variadic.
//
// Java: CheckJSDoc.validateRestParameter(Node)
func (c *CheckJSDoc) ValidateRestParameter(restParam *rhino.Node) {
	if !restParam.IsRest() || !restParam.GetParent().IsParamList() {
		return
	}
	paramList := restParam.GetParent()
	inlineInfo := restParam.GetFirstChild().GetJSDocInfo()
	functionInfo := NodeUtil_GetBestJSDocInfo(paramList.GetParent())
	var paramTypeAnnotation *rhino.JSTypeExpression
	if inlineInfo != nil {
		paramTypeAnnotation = inlineInfo.GetType()
	} else if functionInfo != nil {
		if restParam.GetFirstChild().IsName() {
			paramName := restParam.GetFirstChild().GetString()
			paramTypeAnnotation = functionInfo.GetParameterType(paramName)
		} else {
			// destructuring rest param. use the nth JSDoc parameter if present. the name will not match
			indexOfRest := paramList.GetIndexOfChild(restParam)
			if functionInfo.GetParameterCount() >= indexOfRest {
				// getParameterType(null) (no such parameter) is null.
				if paramName := functionInfo.GetParameterNameAt(indexOfRest); paramName != nil {
					paramTypeAnnotation = functionInfo.GetParameterType(*paramName)
				}
			} else {
				paramTypeAnnotation = nil
			}
		}
	} else {
		paramTypeAnnotation = nil
	}

	if paramTypeAnnotation != nil &&
		paramTypeAnnotation.GetRoot().GetToken() != rhino.Token_ITER_REST {
		c.compiler.Report(JSError_MakeWithNode(restParam, CheckJSDoc_BAD_REST_PARAMETER_ANNOTATION))
	}
}

// ValidateDefaultValue checks that a parameter with a default value is marked
// as optional. TODO(bradfordcsmith): This is redundant. We shouldn't require
// it.
//
// Java: CheckJSDoc.validateDefaultValue(Node)
func (c *CheckJSDoc) ValidateDefaultValue(n *rhino.Node) {
	if n.IsDefaultValue() && n.GetParent().IsParamList() {
		targetNode := n.GetFirstChild()
		info := targetNode.GetJSDocInfo()
		if info == nil {
			return
		}

		typeExpr := info.GetType()
		if typeExpr == nil {
			return
		}

		typeNode := typeExpr.GetRoot()
		if typeNode.GetToken() != rhino.Token_EQUALS {
			c.Report(typeNode, CheckJSDoc_DEFAULT_PARAM_MUST_BE_MARKED_OPTIONAL)
		}
	}
}

// ValidateNoSideEffects checks that @modifies annotations are only present in
// externs.
//
// Java: CheckJSDoc.validateNoSideEffects(Node, JSDocInfo)
func (c *CheckJSDoc) ValidateNoSideEffects(n *rhino.Node, info *rhino.JSDocInfo) {
	// Cannot have @modifies in regular (non externs) js. Report errors.
	if info == nil {
		return
	}

	if n.IsFromExterns() {
		return
	}

	if info.HasSideEffectsArgumentsAnnotation() || info.ModifiesThis() {
		c.Report(n, CheckJSDoc_INVALID_MODIFIES_ANNOTATION)
	}
}

// ValidateDefinesDeclaration checks that a let declaration is not used with
// {@defines}.
//
// Java: CheckJSDoc.validateDefinesDeclaration(Node, JSDocInfo)
func (c *CheckJSDoc) ValidateDefinesDeclaration(n *rhino.Node, info *rhino.JSDocInfo) {
	if info != nil && info.IsDefine() && n.IsLet() {
		c.Report(n, CheckJSDoc_INVALID_DEFINE_ON_LET)
	}
}

// ValidateImplicitCast checks that an @implicitCast annotation is in the
// externs.
//
// Java: CheckJSDoc.validateImplicitCast(Node, JSDocInfo)
func (c *CheckJSDoc) ValidateImplicitCast(n *rhino.Node, info *rhino.JSDocInfo) {
	if !c.inExterns && info != nil && info.IsImplicitCast() {
		c.Report(n, TypeCheck_ILLEGAL_IMPLICIT_CAST)
	}
}

// ValidateClosurePrimitive checks that a @closurePrimitive {id} is on a
// function.
//
// Java: CheckJSDoc.validateClosurePrimitive(Node, JSDocInfo)
func (c *CheckJSDoc) ValidateClosurePrimitive(n *rhino.Node, info *rhino.JSDocInfo) {
	if info == nil || !info.HasClosurePrimitiveId() {
		return
	}

	if !c.IsJSDocOnFunctionNode(n, info) {
		c.Report(n, CheckJSDoc_MISPLACED_ANNOTATION, "closurePrimitive", "must be on a function node")
	}
}

// ValidateReturnJsDoc checks that there are no annotations on return.
//
// Java: CheckJSDoc.validateReturnJsDoc(Node, JSDocInfo)
func (c *CheckJSDoc) ValidateReturnJsDoc(n *rhino.Node, info *rhino.JSDocInfo) {
	if !n.IsReturn() || info == nil {
		return
	}
	// @type and @typedef are handled separately
	if info.ContainsDeclaration() && !info.HasType() && !info.HasTypedefType() {
		c.Report(n, CheckJSDoc_JSDOC_ON_RETURN)
	}
}

// ValidateTsType checks that a @tsType is on a function in a supported file.
//
// Java: CheckJSDoc.validateTsType(Node, JSDocInfo)
func (c *CheckJSDoc) ValidateTsType(n *rhino.Node, info *rhino.JSDocInfo) {
	if info == nil || len(info.GetTsTypes()) == 0 {
		return
	}

	if !c.IsJSDocOnFunctionNode(n, info) {
		c.Report(n, CheckJSDoc_MISPLACED_ANNOTATION, "tsType", "must be on a function node")
	}
}

// Java: CheckJSDoc.isFromTs(Node)
func CheckJSDoc_IsFromTs(n *rhino.Node) bool {
	return n.GetStaticSourceFile().IsTypeScriptSource()
}

// Java: CheckJSDoc.validateJsDocTypeNames(JSDocInfo)
func (c *CheckJSDoc) ValidateJsDocTypeNames(info *rhino.JSDocInfo) {
	if info == nil {
		return
	}
	for _, typeNode := range info.GetTypeNodes() {
		NodeUtil_VisitPreOrder(typeNode, c.checkJsDocTypesVisitor)
	}
}

// CheckJSDoc_CheckJsdocTypes bans any references to compiler internal
// implementation details.
type CheckJSDoc_CheckJsdocTypes struct{}

// Java: CheckJSDoc.CheckJsdocTypes.visit(Node)
func (*CheckJSDoc_CheckJsdocTypes) Visit(typeRefNode *rhino.Node) {
	if !typeRefNode.IsStringLit() {
		return
	}
	// A type name that might be simple like "Foo" or qualified like "foo.Bar".
	typeName := typeRefNode.GetString()
	dot := jutil.StringIndexOfChar(typeName, '.')
	rootOfType := typeName
	if dot != -1 {
		rootOfType = jutil.StringSubstring(typeName, 0, dot)
	}

	// Prevent handwritten JS from referencing a module export or module content name that's
	// synthesized by ClosureRewriteModule. Prefix the JSDoc references with
	// "UnrecognizedType_" and leave it to the typechecker to report a JSC_UNRECOGNIZED_TYPE_ERROR
	//  * why not report an error here? even if we did report an error, we
	//    should still add the prefix to ensure the typechecker doesn't resolve this type. Some
	//    builds and/or files suppress type errors.
	//    TODO(lharker): consider reporting an unsuppressible error instead of doing this
	//    rewriting, if we can clean up all existing violations of this error.
	//  * why do this here instead of in the ClosureRewriteModule pass? the Es6RewriteModule runs
	//    before ClosureRewriteModule and may add references to these module export names.
	//  * note: for references in code, not JSDoc, undefined variable checks will handle this.
	// TODO(b/144593112): remove this when ClosureRewriteModule always runs after typechecking
	if ClosureRewriteModule_IsModuleExport(rootOfType) ||
		ClosureRewriteModule_IsModuleContent(rootOfType) {
		typeRefNode.SetString("UnrecognizedType_" + typeName)
	}
}

// ValidateIsUsedViaDotConstructor checks that @usedViaDotConstructor is only
// used on constructors.
//
// Java: CheckJSDoc.validateIsUsedViaDotConstructor(Node, JSDocInfo)
func (c *CheckJSDoc) ValidateIsUsedViaDotConstructor(n *rhino.Node, info *rhino.JSDocInfo) {
	if info == nil || !info.IsUsedViaDotConstructor() {
		return
	}
	if !(n.IsFunction() && info.IsConstructor()) &&
		!NodeUtil_IsEs6ConstructorMemberFunctionDef(n) {
		c.Report(n, CheckJSDoc_MISPLACED_ANNOTATION, "usedViaDotConstructor", "must be on a constructor")
	}
}
