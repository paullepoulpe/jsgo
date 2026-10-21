// Port of com/google/javascript/jscomp/ClosureCheckModule.java @ closure-compiler v20261015

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
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ClosureCheckModule checks that goog.module() is used correctly.
//
// Note that this file only does checks that can be done per-file. Whole
// program checks happen during goog.module rewriting, in
// ClosureRewriteModule.
type ClosureCheckModule struct {
	NodeTraversal_AbstractModuleCallbackBase

	currentModuleInfo *ClosureCheckModule_ModuleInfo
}

var ClosureCheckModule_AT_EXPORT_IN_GOOG_MODULE = DiagnosticType_Error(
	"JSC_AT_EXPORT_IN_GOOG_MODULE",
	"@export has no effect on top-level names in a goog.module."+
		" Consider using goog.exportSymbol instead.")

var ClosureCheckModule_AT_EXPORT_IN_NON_LEGACY_GOOG_MODULE = DiagnosticType_Error(
	"JSC_AT_EXPORT_IN_NON_LEGACY_GOOG_MODULE",
	"@export is not allowed here in a non-legacy goog.module."+
		" Consider using goog.exportSymbol instead.")

var ClosureCheckModule_GOOG_MODULE_IN_NON_MODULE = DiagnosticType_Error(
	"JSC_GOOG_MODULE_IN_NON_MODULE",
	"goog.module() call must be the first statement in a module.")

var ClosureCheckModule_GOOG_MODULE_MISPLACED = DiagnosticType_Error(
	"JSC_GOOG_MODULE_MISPLACED", "goog.module() call must be the first statement in a file.")

var ClosureCheckModule_DECLARE_LEGACY_NAMESPACE_IN_NON_MODULE = DiagnosticType_Error(
	"JSC_DECLARE_LEGACY_NAMESPACE_IN_NON_MODULE",
	"goog.module.declareLegacyNamespace may only be called in a goog.module.")

var ClosureCheckModule_GOOG_MODULE_REFERENCES_THIS = DiagnosticType_Error(
	"JSC_GOOG_MODULE_REFERENCES_THIS", "The body of a goog.module cannot reference 'this'.")

var ClosureCheckModule_GOOG_MODULE_USES_THROW = DiagnosticType_Error(
	"JSC_GOOG_MODULE_USES_THROW", "The body of a goog.module cannot use 'throw'.")

var ClosureCheckModule_DUPLICATE_NAME_SHORT_REQUIRE = DiagnosticType_Error(
	"JSC_DUPLICATE_NAME_SHORT_REQUIRE",
	"Found multiple goog.require statements importing identifier ''{0}''.")

var ClosureCheckModule_INVALID_DESTRUCTURING_REQUIRE = DiagnosticType_Error(
	"JSC_INVALID_DESTRUCTURING_REQUIRE",
	"Destructuring goog.require must be a simple object pattern.")

var ClosureCheckModule_LET_GOOG_REQUIRE = DiagnosticType_Disabled(
	"JSC_LET_GOOG_REQUIRE",
	"Module imports must be constant. Please use ''const'' instead of ''let''.")

var ClosureCheckModule_MULTIPLE_MODULES_IN_FILE = DiagnosticType_Error(
	"JSC_MULTIPLE_MODULES_IN_FILE",
	"There should only be a single goog.module() statement per file.")

var ClosureCheckModule_ONE_REQUIRE_PER_DECLARATION = DiagnosticType_Error(
	"JSC_ONE_REQUIRE_PER_DECLARATION",
	"There may only be one goog.require() per var/let/const declaration.")

var ClosureCheckModule_INCORRECT_SHORTNAME_CAPITALIZATION = DiagnosticType_Disabled(
	"JSC_INCORRECT_SHORTNAME_CAPITALIZATION",
	"The capitalization of short name {0} is incorrect; it should be {1}.")

var ClosureCheckModule_EXPORT_NOT_AT_MODULE_SCOPE = DiagnosticType_Error(
	"JSC_EXPORT_NOT_AT_MODULE_SCOPE", "Exports must be at the top-level of a module")

var ClosureCheckModule_EXPORT_NOT_A_STATEMENT = DiagnosticType_Error("JSC_EXPORT_NOT_A_STATEMENT", "Exports should be a statement.")

var ClosureCheckModule_EXPORT_REPEATED_ERROR = DiagnosticType_Error(
	"JSC_EXPORT_REPEATED_ERROR",
	"Name cannot be exported multiple times. Previous export on line {0}.")

var ClosureCheckModule_REFERENCE_TO_MODULE_GLOBAL_NAME = DiagnosticType_Error(
	"JSC_REFERENCE_TO_MODULE_GLOBAL_NAME",
	"References to the global name of a module are not allowed. Perhaps you meant exports?")

var ClosureCheckModule_REFERENCE_TO_FULLY_QUALIFIED_IMPORT_NAME = DiagnosticType_Disabled(
	"JSC_REFERENCE_TO_FULLY_QUALIFIED_IMPORT_NAME",
	"Reference to fully qualified import name ''{0}''."+
		" Imports in goog.module should use the return value of"+
		" goog.require / goog.forwardDeclare instead.")

var ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME = DiagnosticType_Disabled(
	"JSC_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME",
	"Reference to fully qualified import name ''{0}''."+
		" Please use the short name ''{1}'' instead.")

var ClosureCheckModule_USE_OF_GOOG_PROVIDE = DiagnosticType_Disabled(
	"JSC_USE_OF_GOOG_PROVIDE",
	"goog.provide is deprecated in favor of goog.module.")

var ClosureCheckModule_REQUIRE_NOT_AT_TOP_LEVEL = DiagnosticType_Error(
	"JSC_REQUIRE_NOT_AT_TOP_LEVEL", "goog.require() must be called at file scope.")

var ClosureCheckModule_LEGACY_NAMESPACE_NOT_AT_TOP_LEVEL = DiagnosticType_Error(
	"JSC_LEGACY_NAMESPACE_NOT_AT_TOP_LEVEL",
	"goog.module.declareLegacyNamespace() does not return a value")

var ClosureCheckModule_LEGACY_NAMESPACE_NOT_AFTER_GOOG_MODULE = DiagnosticType_Error(
	"JSC_LEGACY_NAMESPACE_NOT_AT_TOP_LEVEL",
	"goog.module.declareLegacyNamespace() must be immediately after the"+
		" goog.module('...'); call")

var ClosureCheckModule_LEGACY_NAMESPACE_ARGUMENT = DiagnosticType_Error(
	"JSC_LEGACY_NAMESPACE_ARGUMENT",
	"goog.module.declareLegacyNamespace() takes no arguments")

var closureCheckModule_GOOG_MODULE = rhino.QualifiedName_Of("goog.module")
var closureCheckModule_GOOG_PROVIDE = rhino.QualifiedName_Of("goog.provide")
var closureCheckModule_GOOG_REQUIRE = rhino.QualifiedName_Of("goog.require")
var closureCheckModule_GOOG_REQUIRE_TYPE = rhino.QualifiedName_Of("goog.requireType")
var closureCheckModule_GOOG_MODULE_GET = rhino.QualifiedName_Of("goog.module.get")
var closureCheckModule_GOOG_FORWARD_DECLARE = rhino.QualifiedName_Of("goog.forwardDeclare")
var closureCheckModule_GOOG_MODULE_DECLARE_LEGACY_NAMESPACE = rhino.QualifiedName_Of("goog.module.declareLegacyNamespace")

// ClosureCheckModule_ModuleInfo is ModuleInfo.
type ClosureCheckModule_ModuleInfo struct {
	// Name of the module in question (i.e. the argument to goog.module)
	name string
	// Mapping from fully qualified goog.required names to the import LHS
	// node. For standalone goog.require()s the value is the EXPR_RESULT
	// node.
	//
	// Go-only: language-adapter — the Java map may receive a null key (the
	// first argument of a malformed goog.require call); it is only ever read
	// with non-null keys, so null keys are not stored.
	importsByLongRequiredName *jutil.LinkedHashMap[string, *rhino.Node]
	// Module-local short names for goog.required symbols.
	shortImportNames *jutil.LinkedHashSet[string]
	// A map from the export names "exports.name" to the nodes of those named
	// exports. The default export is keyed with just "exports"
	exportNodesByName *jutil.LinkedHashMap[string, *rhino.Node]
}

// Java: ClosureCheckModule.ModuleInfo(String)
func newClosureCheckModule_ModuleInfo(moduleName string) *ClosureCheckModule_ModuleInfo {
	i := &ClosureCheckModule_ModuleInfo{}
	i.importsByLongRequiredName = jutil.NewLinkedHashMap[string, *rhino.Node]()
	i.shortImportNames = jutil.NewLinkedHashSet[string]()
	i.exportNodesByName = jutil.NewLinkedHashMap[string, *rhino.Node]()
	i.name = moduleName
	return i
}

// Java: ClosureCheckModule(AbstractCompiler, ModuleMetadataMap)
func NewClosureCheckModule(compiler AbstractCompiler, moduleMetadataMap *ModuleMetadataMap) *ClosureCheckModule {
	c := &ClosureCheckModule{}
	c.InitNodeTraversal_AbstractModuleCallbackBase(c, compiler, moduleMetadataMap)
	c.currentModuleInfo = nil
	return c
}

// Java: ClosureCheckModule.process(Node, Node)
func (c *ClosureCheckModule) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(c.compiler, root, c)
}

// Java: ClosureCheckModule.enterModule(ModuleMetadata, Node)
func (c *ClosureCheckModule) EnterModule(currentModule *ModuleMetadataMap_ModuleMetadata, moduleScopeRoot *rhino.Node) {
	if !currentModule.IsGoogModule() {
		return
	}

	jutil.CheckState(c.currentModuleInfo == nil)
	jutil.CheckState(!currentModule.GoogNamespaces().IsEmpty())
	first := ""
	for ns := range currentModule.GoogNamespaces().All() {
		first = ns
		break
	}
	c.currentModuleInfo = newClosureCheckModule_ModuleInfo(first)
}

// Java: ClosureCheckModule.exitModule(ModuleMetadata, Node)
func (c *ClosureCheckModule) ExitModule(currentModule *ModuleMetadataMap_ModuleMetadata, moduleScopeRoot *rhino.Node) {
	c.currentModuleInfo = nil
}

// Java: ClosureCheckModule.visit(NodeTraversal, Node, ModuleMetadata, Node)
func (c *ClosureCheckModule) VisitWithModule(
	t *NodeTraversal,
	n *rhino.Node,
	currentModule *ModuleMetadataMap_ModuleMetadata,
	moduleScopeRoot *rhino.Node) {
	parent := n.GetParent()
	if c.currentModuleInfo == nil {
		if NodeUtil_IsCallToQualifiedName(n, closureCheckModule_GOOG_MODULE) {
			t.Report(n, ClosureCheckModule_GOOG_MODULE_IN_NON_MODULE)
		} else if NodeUtil_IsGoogModuleDeclareLegacyNamespaceCall(n) {
			t.Report(n, ClosureCheckModule_DECLARE_LEGACY_NAMESPACE_IN_NON_MODULE)
		} else if NodeUtil_IsCallToQualifiedName(n, closureCheckModule_GOOG_PROVIDE) {
			// This error is reported here, rather than in a provide-specific pass, because it
			// must be reported prior to ClosureRewriteModule converting legacy goog.modules into
			// goog.provides.
			t.Report(n, ClosureCheckModule_USE_OF_GOOG_PROVIDE)
		}
		return
	}
	jsDoc := n.GetJSDocInfo()
	if jsDoc != nil {
		c.CheckJSDoc(t, jsDoc)
	}
	switch n.GetToken() {
	case rhino.Token_CALL:
		callee := n.GetFirstChild()
		if closureCheckModule_GOOG_MODULE.Matches(callee) {
			firstArgumentName := ClosureCheckModule_ExtractFirstArgumentName(n)
			if firstArgumentName == nil || c.currentModuleInfo.name != *firstArgumentName {
				t.Report(n, ClosureCheckModule_MULTIPLE_MODULES_IN_FILE)
			} else if !ClosureCheckModule_IsFirstExpressionInGoogModule(parent) {
				t.Report(n, ClosureCheckModule_GOOG_MODULE_MISPLACED)
			}
		} else if closureCheckModule_GOOG_REQUIRE.Matches(callee) ||
			closureCheckModule_GOOG_REQUIRE_TYPE.Matches(callee) ||
			closureCheckModule_GOOG_FORWARD_DECLARE.Matches(callee) {
			c.CheckRequireCall(t, n, parent)
		} else if closureCheckModule_GOOG_MODULE_GET.Matches(callee) && t.InModuleHoistScope() {
			t.Report(n, ClosurePrimitiveErrors_MODULE_USES_GOOG_MODULE_GET)
		} else if closureCheckModule_GOOG_MODULE_DECLARE_LEGACY_NAMESPACE.Matches(callee) {
			ClosureCheckModule_CheckLegacyNamespaceCall(t, n, parent)
		}
	case rhino.Token_ASSIGN:
		{
			if ClosureCheckModule_IsExportLhs(n.GetFirstChild()) {
				c.CheckModuleExport(t, n, parent)
			}
		}
	case rhino.Token_CLASS, rhino.Token_FUNCTION:
		if !NodeUtil_IsStatement(n) {
			break
		}
		fallthrough
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		if t.InModuleHoistScope() &&
			(n.IsClass() || NodeUtil_GetEnclosingClass(n) == nil) &&
			NodeUtil_GetEnclosingType(n, rhino.Token_OBJECTLIT) == nil {
			jsdoc := NodeUtil_GetBestJSDocInfo(n)
			if jsdoc != nil && jsdoc.IsExport() {
				t.Report(n, ClosureCheckModule_AT_EXPORT_IN_GOOG_MODULE)
			}
		}
	case rhino.Token_THIS:
		if t.InModuleHoistScope() {
			t.Report(n, ClosureCheckModule_GOOG_MODULE_REFERENCES_THIS)
		}
	case rhino.Token_THROW:
		if t.InModuleHoistScope() {
			t.Report(n, ClosureCheckModule_GOOG_MODULE_USES_THROW)
		}
	case rhino.Token_NAME, rhino.Token_GETPROP:
		if n.MatchesQualifiedName(c.currentModuleInfo.name) {
			if n.IsGetProp() {
				// This warning makes sense on NAME nodes as well,
				// but it would require a cleanup to land.
				t.Report(n, ClosureCheckModule_REFERENCE_TO_MODULE_GLOBAL_NAME)
			}
		} else if !parent.IsGetProp() && !c.currentModuleInfo.importsByLongRequiredName.IsEmpty() {
			if qname, ok := n.GetQualifiedNameString(); ok {
				node := n
				var nextPart string
				var nextQnamePart *string

				for {
					if !c.CheckImproperReferenceToImport(t, node, &qname, nextQnamePart) {
						break
					}
					if node.IsName() || !node.HasChildren() {
						// Don't get the value of a declaration (i.e. `var name = 2`)
						return
					}
					nextPart = node.GetString()
					nextQnamePart = &nextPart
					qname = qname[:len(qname)-1-len(nextPart)]
					node = node.GetOnlyChild()
					if node == nil {
						return
					}
				}
			}
		}
	default:
	}
}

// Java: ClosureCheckModule.checkJSDoc(NodeTraversal, JSDocInfo)
func (c *ClosureCheckModule) CheckJSDoc(t *NodeTraversal, jsDoc *rhino.JSDocInfo) {
	for _, typeNode := range jsDoc.GetTypeNodes() {
		c.CheckTypeExpression(t, typeNode)
	}
}

// closureCheckModule_checkTypeExpressionVisitor is the anonymous
// NodeUtil.Visitor of checkTypeExpression.
type closureCheckModule_checkTypeExpressionVisitor struct {
	outer *ClosureCheckModule
	t     *NodeTraversal
}

// Java: ClosureCheckModule.checkTypeExpression(NodeTraversal, Node).new Visitor() {...}.visit(Node)
func (v *closureCheckModule_checkTypeExpressionVisitor) Visit(node *rhino.Node) {
	if !node.IsStringLit() {
		return
	}

	qname := node.GetString()
	var nextQnamePart *string = nil

	for v.outer.CheckImproperReferenceToImport(v.t, node, &qname, nextQnamePart) {
		lastDot := jutil.StringLastIndexOf(qname, ".")
		if lastDot < 0 {
			return
		}
		next := jutil.StringSubstringFrom(qname, lastDot+1)
		nextQnamePart = &next
		qname = jutil.StringSubstring(qname, 0, lastDot)
	}
}

// Java: ClosureCheckModule.checkTypeExpression(NodeTraversal, Node)
func (c *ClosureCheckModule) CheckTypeExpression(t *NodeTraversal, typeNode *rhino.Node) {
	NodeUtil_VisitPreOrder(
		typeNode,
		&closureCheckModule_checkTypeExpressionVisitor{outer: c, t: t})
}

// Java: ClosureCheckModule.isLocalVar(Var)
func ClosureCheckModule_IsLocalVar(v *Var) bool {
	if v == nil {
		return false
	}
	return v.IsLocal()
}

// CheckImproperReferenceToImport returns true if the reference might be
// valid, false means an error was reported
//
// Java: ClosureCheckModule.checkImproperReferenceToImport(NodeTraversal, Node, String, String)
func (c *ClosureCheckModule) CheckImproperReferenceToImport(
	t *NodeTraversal, n *rhino.Node, qname *string, nextQnamePart *string) bool {
	if qname == nil {
		return true
	}

	importLhs := c.currentModuleInfo.importsByLongRequiredName.Get(*qname)
	if importLhs == nil {
		return true
	}

	if !(jutil.StringIndexOf(*qname, ".") >= 0) && ClosureCheckModule_IsLocalVar(t.GetScope().GetVar(*qname)) {
		return true
	}

	if importLhs.IsName() {
		c.compiler.Report(
			JSError_MakeWithNode(
				n,
				ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME,
				*qname,
				importLhs.GetString()))
		return false
	} else if importLhs.IsDestructuringLhs() && nextQnamePart != nil {
		objPattern := importLhs.GetFirstChild()
		if objPattern.IsObjectPattern() {
			for strKey := objPattern.GetFirstChild(); strKey != nil; strKey = strKey.GetNext() {
				// const {foo: barFoo} = goog.require('ns.bar');
				// Should use the short name "barFoo" instead of "ns.bar.foo".
				if strKey.IsStringKey() &&
					strKey.HasOneChild() &&
					strKey.GetFirstChild().IsName() &&
					strKey.GetString() == *nextQnamePart {
					parent := n.GetParent()
					reportNode := n
					if parent != nil && parent.IsGetProp() {
						reportNode = parent
					}
					c.compiler.Report(
						JSError_MakeWithNode(
							reportNode,
							ClosureCheckModule_REFERENCE_TO_SHORT_IMPORT_BY_LONG_NAME_INCLUDING_SHORT_NAME,
							*qname+"."+*nextQnamePart,
							strKey.GetFirstChild().GetString()))
					return false
				}
			}
		}
	}
	if importLhs.IsExprResult() && NodeUtil_IsGoogForwardDeclareCall(importLhs.GetOnlyChild()) {
		// goog.forwardDeclare does not need to be aliased.
		// This is because goog.forwardDeclares that are aliased have stricter semantics than those
		// that are not aliased: aliased goog.forwardDeclares always need to reference a namespace
		// that's actually defined, while non-aliased forwardDeclares may reference namespaces that
		// do not actually exist.
		return true
	}

	c.compiler.Report(JSError_MakeWithNode(n, ClosureCheckModule_REFERENCE_TO_FULLY_QUALIFIED_IMPORT_NAME, *qname))
	return false
}

// ClosureCheckModule_IsExportLhs: is this the LHS of a goog.module export?
// i.e. Either "exports" or "exports.name"
//
// Java: ClosureCheckModule.isExportLhs(Node)
func ClosureCheckModule_IsExportLhs(lhs *rhino.Node) bool {
	if !lhs.IsQualifiedName() {
		return false
	}
	return lhs.MatchesName("exports") ||
		(lhs.IsGetProp() && lhs.GetFirstChild().MatchesName("exports"))
}

// Java: ClosureCheckModule.checkModuleExport(NodeTraversal, Node, Node)
func (c *ClosureCheckModule) CheckModuleExport(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	jutil.CheckArgument(n.IsAssign())
	lhs := n.GetFirstChild()
	jutil.CheckState(ClosureCheckModule_IsExportLhs(lhs))
	// Check multiple exports of the same name
	previousDefinition := c.currentModuleInfo.exportNodesByName.Get(*lhs.GetQualifiedName())
	if previousDefinition != nil {
		previousLine := previousDefinition.GetLineno()
		t.Report(n, ClosureCheckModule_EXPORT_REPEATED_ERROR, strconv.Itoa(previousLine))
	}
	// Check exports in invalid program position
	defaultExportNode := c.currentModuleInfo.exportNodesByName.Get("exports")
	// If we have never seen an `exports =` default export assignment, or this is the
	// default export, then treat this assignment as an export and do the checks it is well formed.
	if defaultExportNode == nil || lhs.MatchesName("exports") {
		c.currentModuleInfo.exportNodesByName.Put(*lhs.GetQualifiedName(), lhs)
		if !t.InModuleScope() {
			t.Report(n, ClosureCheckModule_EXPORT_NOT_AT_MODULE_SCOPE)
		} else if !parent.IsExprResult() {
			t.Report(n, ClosureCheckModule_EXPORT_NOT_A_STATEMENT)
		}
	}
	// Check @export on a module local name
	if !NodeUtil_IsPrototypeProperty(lhs) &&
		!NodeUtil_IsLegacyGoogModuleFile(NodeUtil_GetEnclosingScript(n)) {
		jsDoc := n.GetJSDocInfo()
		if jsDoc != nil && jsDoc.IsExport() {
			t.Report(n, ClosureCheckModule_AT_EXPORT_IN_NON_LEGACY_GOOG_MODULE)
		}
	}
}

// Java: ClosureCheckModule.extractFirstArgumentName(Node)
func ClosureCheckModule_ExtractFirstArgumentName(callNode *rhino.Node) *string {
	firstArg := callNode.GetSecondChild()
	if firstArg != nil && firstArg.IsStringLit() {
		s := firstArg.GetString()
		return &s
	}
	return nil
}

// Java: ClosureCheckModule.checkRequireCall(NodeTraversal, Node, Node)
func (c *ClosureCheckModule) CheckRequireCall(t *NodeTraversal, callNode *rhino.Node, parent *rhino.Node) {
	jutil.CheckState(callNode.IsCall())
	jutil.CheckState(callNode.GetLastChild().IsStringLit())
	switch parent.GetToken() {
	case rhino.Token_EXPR_RESULT:
		key := ClosureCheckModule_ExtractFirstArgumentName(callNode)
		if key != nil {
			c.currentModuleInfo.importsByLongRequiredName.PutIfAbsent(*key, parent)
		}
	case rhino.Token_NAME, rhino.Token_DESTRUCTURING_LHS:
		c.CheckShortGoogRequireCall(t, callNode, parent.GetParent())
	case rhino.Token_AWAIT:
		grandParentToken := parent.GetParent().GetToken()
		if grandParentToken == rhino.Token_DESTRUCTURING_LHS ||
			grandParentToken == rhino.Token_NAME {
			c.CheckShortGoogRequireCall(t, callNode, parent.GetGrandparent())
		}
	default:
		t.Report(callNode, ClosureCheckModule_REQUIRE_NOT_AT_TOP_LEVEL)
	}
}

// Java: ClosureCheckModule.checkShortGoogRequireCall(NodeTraversal, Node, Node)
func (c *ClosureCheckModule) CheckShortGoogRequireCall(t *NodeTraversal, callNode *rhino.Node, declaration *rhino.Node) {
	if declaration.IsLet() && !closureCheckModule_GOOG_FORWARD_DECLARE.Matches(callNode.GetFirstChild()) {
		t.Report(declaration, ClosureCheckModule_LET_GOOG_REQUIRE)
	}
	if !declaration.HasOneChild() {
		t.Report(declaration, ClosureCheckModule_ONE_REQUIRE_PER_DECLARATION)
		return
	}
	lhs := declaration.GetFirstChild()
	if lhs.IsDestructuringLhs() {
		if !ClosureCheckModule_IsValidDestructuringImport(lhs) {
			t.Report(declaration, ClosureCheckModule_INVALID_DESTRUCTURING_REQUIRE)
		}
		if closureCheckModule_GOOG_FORWARD_DECLARE.Matches(callNode.GetFirstChild()) {
			t.Report(lhs, ClosurePrimitiveErrors_INVALID_DESTRUCTURING_FORWARD_DECLARE)
		}
	} else if lhs.IsName() {
		ClosureCheckModule_CheckShortName(t, lhs, callNode.GetLastChild().GetString())
	}
	if key := ClosureCheckModule_ExtractFirstArgumentName(callNode); key != nil {
		c.currentModuleInfo.importsByLongRequiredName.Put(*key, lhs)
	}
	NodeUtil_VisitLhsNodesInNode(
		declaration,
		func(nameNode *rhino.Node) {
			name := nameNode.GetString()
			if !c.currentModuleInfo.shortImportNames.Add(name) {
				t.Report(nameNode, ClosureCheckModule_DUPLICATE_NAME_SHORT_REQUIRE, name)
			}
		})
}

// closureCheckModule_asciiIsUpperCase is Guava Ascii.isUpperCase(char).
//
// Go-only: language-adapter — Guava Ascii.
func closureCheckModule_asciiIsUpperCase(c uint16) bool {
	return c >= 'A' && c <= 'Z'
}

// closureCheckModule_asciiToLowerCase is Guava Ascii.toLowerCase(char).
//
// Go-only: language-adapter — Guava Ascii.
func closureCheckModule_asciiToLowerCase(c uint16) uint16 {
	if closureCheckModule_asciiIsUpperCase(c) {
		return c ^ 0x20
	}
	return c
}

// closureCheckModule_asciiToUpperCase is Guava Ascii.toUpperCase(char).
//
// Go-only: language-adapter — Guava Ascii.
func closureCheckModule_asciiToUpperCase(c uint16) uint16 {
	if c >= 'a' && c <= 'z' {
		return c ^ 0x20
	}
	return c
}

// Java: ClosureCheckModule.checkShortName(NodeTraversal, Node, String)
func ClosureCheckModule_CheckShortName(t *NodeTraversal, shortNameNode *rhino.Node, namespace string) {
	nextQnamePart := shortNameNode.GetString()
	if jutil.StringStartsWith(namespace, "google3.") {
		// `google3` namespaces don't provide capitalization context for the import name
		return
	}

	lastSegment := jutil.StringSubstringFrom(namespace, jutil.StringLastIndexOf(namespace, ".")+1)
	if nextQnamePart == lastSegment || lastSegment == "" {
		return
	}

	if closureCheckModule_asciiIsUpperCase(jutil.StringCharAt(nextQnamePart, 0)) != closureCheckModule_asciiIsUpperCase(jutil.StringCharAt(lastSegment, 0)) {
		var newStartChar uint16
		if closureCheckModule_asciiIsUpperCase(jutil.StringCharAt(nextQnamePart, 0)) {
			newStartChar = closureCheckModule_asciiToLowerCase(jutil.StringCharAt(nextQnamePart, 0))
		} else {
			newStartChar = closureCheckModule_asciiToUpperCase(jutil.StringCharAt(nextQnamePart, 0))
		}
		correctedName := jutil.StringFromChar(newStartChar) + jutil.StringSubstringFrom(nextQnamePart, 1)
		t.Report(shortNameNode, ClosureCheckModule_INCORRECT_SHORTNAME_CAPITALIZATION, nextQnamePart, correctedName)
	}
}

// Java: ClosureCheckModule.isValidDestructuringImport(Node)
func ClosureCheckModule_IsValidDestructuringImport(destructuringLhs *rhino.Node) bool {
	jutil.CheckArgument(destructuringLhs.IsDestructuringLhs())
	objectPattern := destructuringLhs.GetFirstChild()
	if !objectPattern.IsObjectPattern() {
		return false
	}
	for stringKey := objectPattern.GetFirstChild(); stringKey != nil; stringKey = stringKey.GetNext() {
		if !stringKey.IsStringKey() {
			return false
		} else if !stringKey.HasChildren() || !stringKey.GetFirstChild().IsName() {
			return false
		}
	}
	return true
}

// ClosureCheckModule_CheckLegacyNamespaceCall validates the position of a
// goog.module.declareLegacyNamespace(); call
//
// Java: ClosureCheckModule.checkLegacyNamespaceCall(NodeTraversal, Node, Node)
func ClosureCheckModule_CheckLegacyNamespaceCall(t *NodeTraversal, callNode *rhino.Node, parent *rhino.Node) {
	jutil.CheckArgument(callNode.IsCall())
	if callNode.GetChildCount() > 1 {
		t.Report(callNode, ClosureCheckModule_LEGACY_NAMESPACE_ARGUMENT)
	}
	if !parent.IsExprResult() {
		t.Report(callNode, ClosureCheckModule_LEGACY_NAMESPACE_NOT_AT_TOP_LEVEL)
		return
	}
	prev := parent.GetPrevious()
	if prev == nil || !NodeUtil_IsGoogModuleCall(prev) {
		t.Report(callNode, ClosureCheckModule_LEGACY_NAMESPACE_NOT_AFTER_GOOG_MODULE)
	}
}

// ClosureCheckModule_IsFirstExpressionInGoogModule is whether this is the
// first statement in a module
//
// Java: ClosureCheckModule.isFirstExpressionInGoogModule(Node)
func ClosureCheckModule_IsFirstExpressionInGoogModule(node *rhino.Node) bool {
	if !node.IsExprResult() || node.GetPrevious() != nil {
		return false
	}
	statementBlock := node.GetParent()
	return statementBlock.IsModuleBody() || NodeUtil_IsBundledGoogModuleScopeRoot(statementBlock)
}
