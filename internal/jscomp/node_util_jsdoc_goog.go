// Port of com/google/javascript/jscomp/NodeUtil.java @ closure-compiler v20261021

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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// NodeUtil_IsFromTypeSummary determines if the given SCRIPT is a
// @typeSummary file, like an i.js file.
//
// Java: NodeUtil.isFromTypeSummary(Node)
func NodeUtil_IsFromTypeSummary(n rhino.Node) bool {
	jutil.CheckArgumentMsg1(n.IsScript(), "%s", n)
	info := n.GetJSDocInfo()
	return info != nil && info.IsTypeSummary()
}

// NodeUtil_NewExpr creates an EXPR_RESULT.
//
// child is the expression itself; returns a newly created EXPR node with the
// child as subexpression.
//
// Java: NodeUtil.newExpr(Node)
func NodeUtil_NewExpr(child rhino.Node) rhino.Node {
	return rhino.IR_ExprResult(child).Srcref(child)
}

// Java: NodeUtil.isInSyntheticScript(Node)
//
// TODO(tbreisacher): Add a method for detecting nodes injected as runtime
// libraries.
func NodeUtil_IsInSyntheticScript(n rhino.Node) bool {
	sourceFileName := n.GetSourceFileName()
	return sourceFileName != nil &&
		(strings.HasPrefix(*sourceFileName, " [synthetic:") ||
			strings.HasPrefix(*sourceFileName, AbstractCompiler_RUNTIME_LIB_DIR))
}

// NodeUtil_DeleteNode permanently deletes the given node from the AST, as
// well as report the related AST changes/deletions to the given compiler.
//
// Java: NodeUtil.deleteNode(Node, AbstractCompiler)
func NodeUtil_DeleteNode(n rhino.Node, compiler AbstractCompiler) {
	parent := n.GetParent()
	NodeUtil_MarkFunctionsDeleted(n, compiler)
	n.Detach()
	compiler.ReportChangeToEnclosingScope(parent)
}

// NodeUtil_GetSourceName returns the source name property on the node or its
// ancestors.
//
// Java: NodeUtil.getSourceName(Node)
func NodeUtil_GetSourceName(n rhino.Node) *string {
	var sourceName *string
	for sourceName == nil && !n.IsNil() {
		sourceName = n.GetSourceFileName()
		n = n.GetParent()
	}
	return sourceName
}

// NodeUtil_GetSourceFile returns the source name property on the node or its
// ancestors.
//
// Java: NodeUtil.getSourceFile(Node)
func NodeUtil_GetSourceFile(n rhino.Node) rhino.StaticSourceFile {
	var sourceName rhino.StaticSourceFile
	for sourceName == nil && !n.IsNil() {
		sourceName = n.GetStaticSourceFile()
		n = n.GetParent()
	}
	return sourceName
}

// NodeUtil_GetInputId returns the InputId property on the node or its
// ancestors.
//
// Java: NodeUtil.getInputId(Node)
func NodeUtil_GetInputId(n rhino.Node) *rhino.InputId {
	for !n.IsNil() && !n.IsScript() {
		n = n.GetParent()
	}

	if !n.IsNil() && n.IsScript() {
		return n.GetInputId()
	}
	return nil
}

// NodeUtil_NewCallNode returns a new CALL node with the "FREE_CALL" set
// based on call target.
//
// Java: NodeUtil.newCallNode(Node, Node...)
func NodeUtil_NewCallNode(callTarget rhino.Node, parameters ...rhino.Node) rhino.Node {
	isFreeCall := !NodeUtil_IsNormalGet(callTarget)
	call := rhino.IR_Call(callTarget)
	call.PutBooleanProp(rhino.Node_FREE_CALL, isFreeCall)
	for _, parameter := range parameters {
		call.AddChildToBack(parameter)
	}
	return call
}

// NodeUtil_GetDeclaredTypeExpression returns declared JSDoc type for the
// given name declaration, or null if none present.
//
// Java: NodeUtil.getDeclaredTypeExpression(Node)
func NodeUtil_GetDeclaredTypeExpression(declaration rhino.Node) *rhino.JSTypeExpression {
	jutil.CheckArgument(declaration.IsName() || declaration.IsStringKey())
	nameJsdoc := NodeUtil_GetBestJSDocInfo(declaration)
	if nameJsdoc != nil {
		return nameJsdoc.GetType()
	}
	parent := declaration.GetParent()
	if parent.IsRest() || parent.IsDefaultValue() {
		parent = parent.GetParent()
	}
	if parent.IsParamList() {
		functionJsdoc := NodeUtil_GetBestJSDocInfo(parent.GetParent())
		if functionJsdoc != nil {
			return functionJsdoc.GetParameterType(declaration.GetString())
		}
	}
	return nil
}

// NodeUtil_GetBestJSDocInfo finds the best JSDoc for the given node.
//
// Java: NodeUtil.getBestJSDocInfo(Node)
func NodeUtil_GetBestJSDocInfo(n rhino.Node) *rhino.JSDocInfo {
	jsdocNode := NodeUtil_GetBestJSDocInfoNode(n)
	if jsdocNode.IsNil() {
		return nil
	}
	return jsdocNode.GetJSDocInfo()
}

// NodeUtil_GetBestJsDocInfoNodeStrict finds the best JSDocInfo node for the
// given node. This version is stricter than getBestJSDocInfoNode in that it
// only accepts a node that is either a lhs "let, const, var" name node, or a
// RHS class/function, or a class/function declaration, or an export.
//
// Throws IllegalStateException if any other input node token is provided.
//
// Java: NodeUtil.getBestJsDocInfoNodeStrict(Node)
func NodeUtil_GetBestJsDocInfoNodeStrict(n rhino.Node) rhino.Node {
	isDeclaredName := n.IsName() && NodeUtil_IsNameDeclaration(n.GetParent())
	if isDeclaredName {
		return NodeUtil_GetBestJsDocInfoNodeInternal(n)
	}
	isDeclaredClass := NodeUtil_IsClassDeclaration(n)
	if isDeclaredClass {
		return NodeUtil_GetBestJsDocInfoNodeInternal(n)
	}
	isDeclaredFunction := NodeUtil_IsFunctionDeclaration(n)
	if isDeclaredFunction {
		return NodeUtil_GetBestJsDocInfoNodeInternal(n)
	}
	isRhsClass := NodeUtil_IsRhsClass(n)
	if isRhsClass {
		return NodeUtil_GetBestJsDocInfoNodeInternal(n)
	}
	isRhsFunction := NodeUtil_IsRhsFunction(n)
	if isRhsFunction {
		return NodeUtil_GetBestJsDocInfoNodeInternal(n)
	}
	isRhsClassName := n.IsName() && NodeUtil_IsRhsClass(n.GetParent())
	if isRhsClassName {
		return NodeUtil_GetBestJsDocInfoNodeInternal(n)
	}
	isRhsFunctionName := n.IsName() && NodeUtil_IsRhsFunction(n.GetParent())
	if isRhsFunctionName {
		return NodeUtil_GetBestJsDocInfoNodeInternal(n)
	}
	isExport := n.IsExport()
	if isExport {
		return NodeUtil_GetBestJsDocInfoNodeInternal(n)
	}
	panic(jutil.NewIllegalStateExceptionMsg("Not allowed to get JSDocInfo node for node: " + n.String()))
}

// NodeUtil_GetBestJSDocInfoNode finds the best JSDocInfo node for the given
// node.
//
// Deprecated: because we want to control the input node tokens accepted by
// this function. Use getBestJSDocInfoNodeStrict instead.
//
// Java: NodeUtil.getBestJSDocInfoNode(Node)
func NodeUtil_GetBestJSDocInfoNode(n rhino.Node) rhino.Node {
	return NodeUtil_GetBestJsDocInfoNodeInternal(n)
}

// Java: NodeUtil.getBestJsDocInfoNodeInternal(Node)
func NodeUtil_GetBestJsDocInfoNodeInternal(n rhino.Node) rhino.Node {
	if exprStmt, ok := ast.AsExprStmt(n); ok {
		return NodeUtil_GetBestJSDocInfoNode(exprStmt.Expr().Node())
	}
	info := n.GetJSDocInfo()
	if info == nil {
		parent := n.GetParent()
		if parent.IsNil() || n.IsExprResult() {
			return rhino.NilNode
		}

		role := ast.RoleIn(n, parent)
		if parent.IsName() ||
			parent.IsExport() ||
			NodeUtil_IsNameDeclaration(parent) ||
			parent.IsDeclare() ||
			parent.IsAssign() {
			return NodeUtil_GetBestJSDocInfoNode(parent)
		} else if NodeUtil_MayBeObjectLitKey(parent) || parent.IsComputedProp() {
			return parent
		} else if role == ast.RoleFunctionName || role == ast.RoleClassName {
			// n is the NAME node of the function/class.
			return NodeUtil_GetBestJSDocInfoNode(parent)
		} else if role == ast.RoleHookThen ||
			role == ast.RoleHookElse ||
			parent.IsOr() ||
			parent.IsAnd() ||
			(parent.IsComma() && role == ast.RoleBinaryRight) {
			return NodeUtil_GetBestJSDocInfoNode(parent)
		}
	}
	return n
}

// Java: NodeUtil.isGoogModuleCall(Node)
func NodeUtil_IsGoogModuleCall(n rhino.Node) bool {
	if exprStmt, ok := ast.AsExprStmt(n); ok {
		if call, ok := ast.AsCall(exprStmt.Expr().Node()); ok {
			return NodeUtil_GOOG_MODULE.Matches(call.Callee().Node())
		}
	}
	return false
}

// Java: NodeUtil.isGoogModuleGetCall(Node)
func NodeUtil_IsGoogModuleGetCall(callNode rhino.Node) bool {
	call, ok := ast.AsCall(callNode)
	if !ok || !NodeUtil_GOOG_MODULE_GET.Matches(call.Callee().Node()) {
		return false
	}
	arg, ok := call.OnlyArg()
	return ok && arg.Node().IsStringLit()
}

// Java: NodeUtil.isGoogRequireCall(Node)
func NodeUtil_IsGoogRequireCall(call rhino.Node) bool {
	if c, ok := ast.AsCall(call); ok {
		return NodeUtil_GOOG_REQUIRE.Matches(c.Callee().Node())
	}
	return false
}

// Java: NodeUtil.isGoogRequireTypeCall(Node)
func NodeUtil_IsGoogRequireTypeCall(call rhino.Node) bool {
	if c, ok := ast.AsCall(call); ok {
		return NodeUtil_GOOG_REQUIRE_TYPE.Matches(c.Callee().Node())
	}
	return false
}

// Java: NodeUtil.isGoogForwardDeclareCall(Node)
func NodeUtil_IsGoogForwardDeclareCall(call rhino.Node) bool {
	if c, ok := ast.AsCall(call); ok {
		return NodeUtil_GOOG_FORWARD_DECLARE.Matches(c.Callee().Node())
	}
	return false
}

// Java: NodeUtil.isGoogRequireDynamicCall(Node)
func NodeUtil_IsGoogRequireDynamicCall(call rhino.Node) bool {
	if c, ok := ast.AsCall(call); ok {
		return NodeUtil_GOOG_REQUIRE_DYNAMIC.Matches(c.Callee().Node())
	}
	return false
}

// Java: NodeUtil.isGoogWeakUsageCall(Node)
func NodeUtil_IsGoogWeakUsageCall(call rhino.Node) bool {
	if c, ok := ast.AsCall(call); ok {
		target := c.Callee().Node()
		return NodeUtil_GOOG_WEAK_USAGE.Matches(target) || NodeUtil_GOOG_WEAK_USAGE_MANGLED.Matches(target)
	}
	return false
}

// Java: NodeUtil.isModuleScopeRoot(Node)
func NodeUtil_IsModuleScopeRoot(n rhino.Node) bool {
	return n.IsModuleBody() || NodeUtil_IsBundledGoogModuleScopeRoot(n)
}

// NodeUtil_IsGoogModuleExportsReference: whether the given node is
// referencing the `exports` object in some goog.module.
//
// Java: NodeUtil.isGoogModuleExportsReference(Scope, Node)
func NodeUtil_IsGoogModuleExportsReference(scope *Scope, possibleName rhino.Node) bool {
	if !possibleName.IsName() || possibleName.GetString() != "exports" {
		return false
	}
	var_ := scope.GetVarByID(possibleName.GetStringID())
	var scopeRoot rhino.Node
	if var_ != nil {
		scopeRoot = var_.GetScopeRoot()
	}
	return !scopeRoot.IsNil() &&
		(scopeRoot.IsModuleBody() ||
			(scopeRoot.IsFunction() &&
				NodeUtil_IsBundledGoogModuleScopeRoot(NodeUtil_GetFunctionBody(scopeRoot))))
}

// Java: NodeUtil.isBundledGoogModuleCall(Node)
func NodeUtil_IsBundledGoogModuleCall(n rhino.Node) bool {
	// TODO(lharker): take an EXPR_RESULT to align with NodeUtil.isGoogModuleCall and
	// NodeUtil.isGoogProvideCall.
	call, ok := ast.AsCall(n)
	if !ok || !NodeUtil_GOOG_LOADMODULE.Matches(call.Callee().Node()) {
		return false
	}
	if _, ok := call.OnlyArg(); !ok {
		return false
	}
	if ast.RoleOf(n) != ast.RoleExprStmtExpr {
		return false
	}
	_, ok = ast.AsScript(n.GetGrandparent())
	return ok
}

// Java: NodeUtil.isBundledGoogModuleScopeRoot(Node)
func NodeUtil_IsBundledGoogModuleScopeRoot(n rhino.Node) bool {
	block, ok := ast.AsBlock(n)
	if !ok {
		return false
	}
	firstStmt, ok := block.FirstStatement()
	if !ok || !NodeUtil_IsGoogModuleCall(firstStmt.Node()) {
		return false
	}
	fn, ok := ast.AsFunction(n.GetParent())
	if !ok {
		return false
	}
	param, ok := fn.Params().OnlyParam()
	if !ok || !param.Node().MatchesName("exports") {
		return false
	}
	return NodeUtil_IsBundledGoogModuleCall(fn.Node().GetParent())
}

// Java: NodeUtil.isGoogModuleDeclareLegacyNamespaceCall(Node)
func NodeUtil_IsGoogModuleDeclareLegacyNamespaceCall(n rhino.Node) bool {
	if exprStmt, ok := ast.AsExprStmt(n); ok {
		if call, ok := ast.AsCall(exprStmt.Expr().Node()); ok {
			return NodeUtil_GOOG_MODULE_DECLARE_LEGACY_NAMESPACE.Matches(call.Callee().Node())
		}
	}
	return false
}

// Java: NodeUtil.isGoogSetTestOnlyCall(Node)
func NodeUtil_IsGoogSetTestOnlyCall(n rhino.Node) bool {
	if exprStmt, ok := ast.AsExprStmt(n); ok {
		if call, ok := ast.AsCall(exprStmt.Expr().Node()); ok {
			return NodeUtil_GOOG_SET_TEST_ONLY.Matches(call.Callee().Node())
		}
	}
	return false
}

// Java: NodeUtil.isTopLevel(Node)
func NodeUtil_IsTopLevel(n rhino.Node) bool {
	return n.IsScript() || n.IsModuleBody()
}

// NodeUtil_IsGoogModuleFile returns whether the node is a goog.module file's
// SCRIPT node.
//
// Java: NodeUtil.isGoogModuleFile(Node)
func NodeUtil_IsGoogModuleFile(n rhino.Node) bool {
	script, ok := ast.AsScript(n)
	if !ok {
		return false
	}
	firstItem, ok := script.FirstItem()
	if !ok {
		return false
	}
	modBody, ok := firstItem.AsModuleBody()
	if !ok {
		return false
	}
	firstStmt, ok := modBody.FirstStatement()
	return ok && NodeUtil_IsGoogModuleCall(firstStmt.Node())
}

// NodeUtil_IsLegacyGoogModuleFile returns whether the node is a SCRIPT node
// for a goog.module that has a declareLegacyNamespace call.
//
// Java: NodeUtil.isLegacyGoogModuleFile(Node)
func NodeUtil_IsLegacyGoogModuleFile(n rhino.Node) bool {
	if !NodeUtil_IsGoogModuleFile(n) {
		return false
	}
	modBody := ast.MustModuleBody(n.GetFirstChild())
	firstStmt, _ := modBody.FirstStatement()
	secondStmt, ok := firstStmt.Next()
	return ok && NodeUtil_IsGoogModuleDeclareLegacyNamespaceCall(secondStmt.Node())
}

// NodeUtil_GetFeatureSetOfScript returns a script node's FeatureSet, which
// is set at parse-time. This may not be up-to-date as passes can add/remove
// features from a script node's descendants.
//
// The feature set will be null if the script node was created artificially
// or if the parser didn't detect any interesting features.
//
// Java: NodeUtil.getFeatureSetOfScript(Node)
func NodeUtil_GetFeatureSetOfScript(scriptNode rhino.Node) *parser.FeatureSet {
	jutil.CheckStateMsg1(scriptNode.IsScript(), "%s", scriptNode)
	prop := scriptNode.GetProp(rhino.Node_FEATURE_SET)
	if prop == nil {
		return nil
	}
	return prop.(*parser.FeatureSet)
}

// NodeUtil_AddFeatureToScript adds the given features to a SCRIPT node's
// FeatureSet property.
//
// Also updates the compiler's FeatureSet.
//
// Java: NodeUtil.addFeatureToScript(Node, Feature, AbstractCompiler)
func NodeUtil_AddFeatureToScript(scriptNode rhino.Node, feature parser.FeatureSet_Feature, compiler AbstractCompiler) {
	jutil.CheckStateMsg1(scriptNode.IsScript(), "%s", scriptNode)
	currentFeatures := NodeUtil_GetFeatureSetOfScript(scriptNode)
	var newFeatures *parser.FeatureSet
	if currentFeatures != nil {
		newFeatures = currentFeatures.With(feature)
	} else {
		newFeatures = parser.FeatureSet_BARE_MINIMUM.With(feature)
	}
	scriptNode.PutProp(rhino.Node_FEATURE_SET, newFeatures)
	if feature != parser.FeatureSet_Feature_MODULES {
		// Except MODULES (which can get reintroduced later in ConvertChunksToEsModules pass),
		// prohibit passes from reintroducing any feature that already got transpiled.
		jutil.CheckStateMsg1(
			compiler.GetAllowableFeatures().ContainsFeature(feature),
			"Cannot add feature: %s. It is not supported in the output language, and either 1) its"+
				" corresponding transpilation pass has already run or 2) transpilation of this"+
				" feature is unsupported entirely",
			feature)
	} else {
		compiler.SetAllowableFeatures(compiler.GetAllowableFeatures().With(feature))
	}
}

// NodeUtil_RemoveFeatureFromScript removes the given feature from the SCRIPT
// node's FeatureSet property.
//
// Removing a feature from a single script does not mean that it's removed
// from all scripts. Hence this method does not remove it from the compiler's
// featureSet. The caller must remove it from the compiler's featureSet if
// it's calling this method for each script.
//
// Java: NodeUtil.removeFeatureFromScript(Node, Feature)
func NodeUtil_RemoveFeatureFromScript(scriptNode rhino.Node, feature parser.FeatureSet_Feature) {
	currentFeatures := NodeUtil_GetFeatureSetOfScript(scriptNode)
	if currentFeatures == nil {
		scriptNode.PutProp(rhino.Node_FEATURE_SET, nil)
		return
	}
	updated := currentFeatures.WithoutFeatures(feature)
	if updated != currentFeatures {
		scriptNode.PutProp(rhino.Node_FEATURE_SET, updated)
	}
}

// NodeUtil_RemoveFeaturesFromScript removes the given featureSet from the
// SCRIPT node's FeatureSet property.
//
// Removing a feature from a single script does not mean that it's removed
// from all scripts. Hence this method does not remove it from the compiler's
// featureSet. The caller must remove it from the compiler's featureSet using
// AbstractCompiler.markFeatureSetNotAllowed or
// AbstractCompiler.markFeatureNotAllowed if it's calling this method for each
// script.
//
// Java: NodeUtil.removeFeaturesFromScript(Node, FeatureSet)
func NodeUtil_RemoveFeaturesFromScript(scriptNode rhino.Node, featureSet *parser.FeatureSet) {
	currentFeatures := NodeUtil_GetFeatureSetOfScript(scriptNode)
	if currentFeatures != nil {
		updated := currentFeatures.Without(featureSet)
		if updated != currentFeatures {
			scriptNode.PutProp(rhino.Node_FEATURE_SET, updated)
		}
	}
}

// NodeUtil_RemoveFeatureFromAllScripts removes the given feature from the
// FEATURE_SET prop of all SCRIPT nodes under root.
//
// Also updates the compiler's FeatureSet.
//
// Java: NodeUtil.removeFeatureFromAllScripts(Node, Feature, AbstractCompiler)
func NodeUtil_RemoveFeatureFromAllScripts(root rhino.Node, feature parser.FeatureSet_Feature, compiler AbstractCompiler) {
	jutil.CheckArgumentMsg1(root.IsRoot(), "%s", root)
	for childNode := root.GetFirstChild(); !childNode.IsNil(); childNode = childNode.GetNext() {
		jutil.CheckState(childNode.IsScript())
		NodeUtil_RemoveFeatureFromScript(childNode, feature)
	}
	compiler.MarkFeatureNotAllowed(feature)
}

// NodeUtil_RemoveFeaturesFromAllScripts removes the given feature from the
// FEATURE_SET prop of all SCRIPT nodes under root.
//
// Also updates the compiler's FeatureSet.
//
// Java: NodeUtil.removeFeaturesFromAllScripts(Node, FeatureSet, AbstractCompiler)
func NodeUtil_RemoveFeaturesFromAllScripts(
	root rhino.Node, featureSet *parser.FeatureSet, compiler AbstractCompiler) {
	jutil.CheckArgumentMsg1(root.IsRoot(), "%s", root)
	for childNode := root.GetFirstChild(); !childNode.IsNil(); childNode = childNode.GetNext() {
		jutil.CheckState(childNode.IsScript())
		NodeUtil_RemoveFeaturesFromScript(childNode, featureSet)
	}
	compiler.MarkFeatureSetNotAllowed(featureSet)
}

// NodeUtil_AddFeatureToAllScripts adds the given feature to the FEATURE_SET
// prop of all scripts under externs and root.
//
// Java: NodeUtil.addFeatureToAllScripts(Node, Feature, AbstractCompiler)
func NodeUtil_AddFeatureToAllScripts(root rhino.Node, feature parser.FeatureSet_Feature, compiler AbstractCompiler) {
	jutil.CheckArgumentMsg1(root.IsRoot(), "%s", root)
	for childNode := root.GetFirstChild(); !childNode.IsNil(); childNode = childNode.GetNext() {
		jutil.CheckState(childNode.IsScript())
		NodeUtil_AddFeatureToScript(childNode, feature, compiler)
	}
}

// NodeUtil_AddFeaturesToScript adds the given features to a SCRIPT node's
// FeatureSet property.
//
// Also updates the compiler's FeatureSet.
//
// Java: NodeUtil.addFeaturesToScript(Node, FeatureSet, AbstractCompiler)
func NodeUtil_AddFeaturesToScript(scriptNode rhino.Node, features *parser.FeatureSet, compiler AbstractCompiler) {
	jutil.CheckStateMsg1(scriptNode.IsScript(), "%s", scriptNode)
	for feature := range features.All() {
		NodeUtil_AddFeatureToScript(scriptNode, feature, compiler)
	}
}

// NodeUtil_GetParamOrPatternNames calls cb with all NAMEs declared in a
// PARAM_LIST or destructuring pattern.
//
// Java: NodeUtil.getParamOrPatternNames(Node, Consumer<Node>)
func NodeUtil_GetParamOrPatternNames(n rhino.Node, cb func(rhino.Node)) {
	parsing.ParsingUtil_GetParamOrPatternNames(n, cb)
}

// NodeUtil_GoogRequire represents a goog.require'd namespace and property
// inside a module.
//
// namespace is the Closure namespace inside the require call; property is
// the property off the required namespace. Needed only for destructuring
// requires. isStrongRequire is true if the require is a goog.require, false
// if it is a goog.requireType.
//
// Java: NodeUtil.GoogRequire
type NodeUtil_GoogRequire struct {
	namespace       string
	property        *string
	isStrongRequire bool
}

// Java: NodeUtil.GoogRequire(String, String, boolean)
func NewNodeUtil_GoogRequire(namespace string, property *string, isStrongRequire bool) *NodeUtil_GoogRequire {
	// Go-only: language-adapter — requireNonNull(namespace, "namespace"): a Go
	// string cannot be null.
	return &NodeUtil_GoogRequire{namespace: namespace, property: property, isStrongRequire: isStrongRequire}
}

// Java: NodeUtil.GoogRequire.namespace()
func (r *NodeUtil_GoogRequire) Namespace() string { return r.namespace }

// Java: NodeUtil.GoogRequire.property()
func (r *NodeUtil_GoogRequire) Property() *string { return r.property }

// Java: NodeUtil.GoogRequire.isStrongRequire()
func (r *NodeUtil_GoogRequire) IsStrongRequire() bool { return r.isStrongRequire }

// Equals is the record's component-wise equals.
//
// Java: NodeUtil.GoogRequire.equals(Object)
func (r *NodeUtil_GoogRequire) Equals(o any) bool {
	other, ok := o.(*NodeUtil_GoogRequire)
	if !ok || other == nil {
		return false
	}
	if r == other {
		return true
	}
	if (r.property == nil) != (other.property == nil) ||
		(r.property != nil && *r.property != *other.property) {
		return false
	}
	return r.namespace == other.namespace && r.isStrongRequire == other.isStrongRequire
}

// HashCode is the record's hashCode: 31 * h + hash(component), in
// declaration order.
//
// Java: NodeUtil.GoogRequire.hashCode()
func (r *NodeUtil_GoogRequire) HashCode() int32 {
	var h int32
	h = 31*h + jutil.StringHashCode(r.namespace)
	var propertyHash int32
	if r.property != nil {
		propertyHash = jutil.StringHashCode(*r.property)
	}
	h = 31*h + propertyHash
	var strongHash int32 = 1237
	if r.isStrongRequire {
		strongHash = 1231
	}
	h = 31*h + strongHash
	return h
}

// String is the record's toString.
//
// Java: NodeUtil.GoogRequire.toString()
func (r *NodeUtil_GoogRequire) String() string {
	return "GoogRequire[namespace=" + r.namespace + ", property=" + jutil.ToJavaString(r.property) +
		", isStrongRequire=" + jutil.ToJavaString(r.isStrongRequire) + "]"
}

// Java: NodeUtil.GoogRequire.fromNamespace(String, boolean)
func NodeUtil_GoogRequire_FromNamespace(namespace string, isStrongRequire bool) *NodeUtil_GoogRequire {
	return NewNodeUtil_GoogRequire(namespace /* property= */, nil, isStrongRequire)
}

// Java: NodeUtil.GoogRequire.fromNamespaceAndProperty(String, String, boolean)
func NodeUtil_GoogRequire_FromNamespaceAndProperty(
	namespace string, property string, isStrongRequire bool) *NodeUtil_GoogRequire {
	return NewNodeUtil_GoogRequire(namespace, jutil.StringPtr(property), isStrongRequire)
}

// Java: NodeUtil.getGoogRequireInfo(String, Scope)
func NodeUtil_GetGoogRequireInfoWithScope(name string, scope *Scope) *NodeUtil_GoogRequire {
	var_ := scope.GetVar(name)
	if var_ == nil {
		return nil
	}
	return NodeUtil_GetGoogRequireInfo(var_)
}

// Java: NodeUtil.getGoogRequireInfo(Var)
func NodeUtil_GetGoogRequireInfo(var_ *Var) *NodeUtil_GoogRequire {
	jutil.CheckNotNull(var_)
	if !var_.GetScopeRoot().IsModuleBody() || var_.GetNameNode().IsNil() {
		return nil
	}

	return NodeUtil_GetGoogRequireInfoNode(var_.GetNameNode())
}

// Java: NodeUtil.getGoogRequireInfo(Node)
func NodeUtil_GetGoogRequireInfoNode(nameNode rhino.Node) *NodeUtil_GoogRequire {
	// IMPORT_STAR is not part of a GoogRequire call.
	if nameNode.IsImportStar() {
		return nil
	}
	jutil.CheckStateMsg1(nameNode.IsName(), "unexpected node type: %s", nameNode)
	name := ast.MustName(nameNode)

	if NodeUtil_IsNameDeclaration(nameNode.GetParent()) {
		requireCall, ok := name.Init()
		if !ok {
			return nil
		}
		isStrongRequire := NodeUtil_IsGoogRequireCall(requireCall.Node())
		if !(isStrongRequire || NodeUtil_IsGoogRequireTypeCall(requireCall.Node())) {
			return nil
		}
		namespace := requireCall.Node().GetSecondChild().GetString()
		return NodeUtil_GoogRequire_FromNamespace(namespace, isStrongRequire)
	} else if strKey, ok := ast.AsStringKey(nameNode.GetParent()); ok && nameNode.GetGrandparent().IsObjectPattern() {
		requireCall := nameNode.GetGrandparent().GetNext()
		if requireCall.IsNil() {
			return nil
		}
		isStrongRequire := NodeUtil_IsGoogRequireCall(requireCall)
		if !(isStrongRequire || NodeUtil_IsGoogRequireTypeCall(requireCall)) {
			return nil
		}
		property := strKey.Key()
		namespace := requireCall.GetSecondChild().GetString()
		return NodeUtil_GoogRequire_FromNamespaceAndProperty(namespace, property, isStrongRequire)
	}
	return nil
}

// NodeUtil_EstimateNumLines estimates the number of lines in the file of
// this script node.
//
// This method returns the line number of the last node in the script +1.
// This is not strictly the number of lines in the file (consider trailing
// comments or whitespace), but should be strongly correlated with it. If
// perfect accuracy is desired the original source file will have to be read.
//
// Java: NodeUtil.estimateNumLines(Node)
func NodeUtil_EstimateNumLines(scriptNode rhino.Node) int {
	jutil.CheckArgument(scriptNode.IsScript())
	current := scriptNode
	for current.HasChildren() {
		current = current.GetLastChild()
	}

	// +1 since this is the start line of the last AST node and we can expect at least one
	// trailing newline
	return current.GetLineno() + 1
}
