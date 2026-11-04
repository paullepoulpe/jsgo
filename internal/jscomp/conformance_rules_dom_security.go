// Port of com/google/javascript/jscomp/ConformanceRules.java @ closure-compiler v20261028

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var (
	googDomQName        = rhino.QualifiedName_Of("goog.dom")
	googDomTagNameQName = rhino.QualifiedName_Of("goog.dom.TagName")
)

// conformanceBanCreateElementRule bans `document.createElement` and
// `goog.dom.createDom` calls with banned literal tag names (e.g. `'script'`).
type conformanceBanCreateElementRule struct {
	conformanceAbstractRule

	bannedTags    map[string]struct{}
	domHelperType rhino.JSType
	documentType  rhino.JSType
}

func newConformanceBanCreateElementRule(compiler AbstractCompiler, requirement *Requirement) (*conformanceBanCreateElementRule, error) {
	r := &conformanceBanCreateElementRule{}
	if err := r.initRule(compiler, requirement, false, r.checkNode); err != nil {
		return nil, err
	}
	values := requirement.GetValueList()
	if len(values) == 0 {
		return nil, errInvalidRequirementSpec("Specify one or more values.")
	}
	bannedTags := make(map[string]struct{}, len(values))
	for _, value := range values {
		bannedTags[jutil.StringToLowerCaseASCII(value)] = struct{}{}
	}
	r.bannedTags = bannedTags
	r.domHelperType = compiler.GetTypeRegistry().GetGlobalType("goog.dom.DomHelper")
	r.documentType = compiler.GetTypeRegistry().GetGlobalType("Document")
	return r, nil
}

func (r *conformanceBanCreateElementRule) checkNode(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if call, ok := ast.AsNormalOrOptChainCall(n); ok {
		if firstArg, ok := call.FirstArg(); ok {
			if tag, ok := ast.AsStringLit(firstArg.Node()); ok {
				if _, banned := r.bannedTags[jutil.StringToLowerCaseASCII(tag.Value())]; banned {
					return r.checkCreateElement(call)
				}
			}
		}
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

func (r *conformanceBanCreateElementRule) checkCreateElement(call ast.NormalOrOptChainCall) *ConformanceRules_ConformanceResult {
	target, ok := ast.AsNormalOrOptChainGetProp(call.Callee().Node())
	if !ok {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	functionName := target.Property()
	if functionName != "createElement" && functionName != "createDom" {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	srcObj := target.Object().Node()
	if googDomQName.Matches(srcObj) {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	type_ := srcObj.GetJSType()
	if isLooseConformanceType(type_) {
		if r.reportLooseTypeViolations {
			return conformancePossibleViolationLooseTypes
		}
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	if (r.domHelperType != nil && r.domHelperType.IsSubtypeOf(type_)) ||
		(r.documentType != nil && r.documentType.IsSubtypeOf(type_)) {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

// conformanceBanCreateDomRule bans `goog.dom.createDom` and
// `goog.dom.DomHelper#createDom` with banned `tagname.attribute` pairs.
type conformanceBanCreateDomRule struct {
	conformanceAbstractRule

	bannedTagAttrs [][2]string
	domHelperType  rhino.JSType
	classNameTypes rhino.JSType
}

func newConformanceBanCreateDomRule(compiler AbstractCompiler, requirement *Requirement) (*conformanceBanCreateDomRule, error) {
	r := &conformanceBanCreateDomRule{}
	if err := r.initRule(compiler, requirement, false, r.checkNode); err != nil {
		return nil, err
	}
	values := requirement.GetValueList()
	bannedTagAttrs := make([][2]string, 0, len(values))
	for _, value := range values {
		tagPart, attrPart, hasDot := strings.Cut(value, ".")
		if !hasDot || tagPart == "" || attrPart == "" || strings.Contains(attrPart, ".") {
			return nil, errInvalidRequirementSpec("Values must be in the format tagname.attribute.")
		}
		lowercasedTag := jutil.StringToLowerCaseRoot(tagPart)
		bannedTagAttrs = append(bannedTagAttrs, [2]string{lowercasedTag, attrPart})
	}
	if len(bannedTagAttrs) == 0 {
		return nil, errInvalidRequirementSpec("Specify one or more values.")
	}
	r.bannedTagAttrs = bannedTagAttrs
	registry := compiler.GetTypeRegistry()
	r.domHelperType = registry.GetGlobalType("goog.dom.DomHelper")
	r.classNameTypes = registry.CreateUnionType(
		registry.GetNativeType(rhino.JSTypeNative_STRING_TYPE),
		registry.GetNativeType(rhino.JSTypeNative_ARRAY_TYPE),
		registry.GetNativeType(rhino.JSTypeNative_NULL_TYPE),
		registry.GetNativeType(rhino.JSTypeNative_VOID_TYPE))
	return r, nil
}

func (r *conformanceBanCreateDomRule) checkNode(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	call, ok := r.asCreateDomCall(n)
	if !ok {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	cn := call.AsCallOrNew()
	firstArg, hasFirst := cn.ArgAt(0)
	attrsArg, hasAttrs := cn.ArgAt(1)
	if !hasFirst || !hasAttrs {
		// goog.dom.createDom('iframe') is fine.
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	attrsNode := attrsArg.Node()

	tagNames, tagNamesIsNull := r.tagNames(firstArg.Node())
	attrsType := attrsNode.GetJSType()
	isSpreadAttrs := attrsNode.IsSpread()
	if attrsType == nil && !isSpreadAttrs {
		// Type information is not available; don't run this check.
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	// String or array attribute sets the class.
	isClassName := attrsType != nil && !attrsType.IsUnknownType() && attrsType.IsSubtypeOf(r.classNameTypes)

	if attrsNode.IsNull() || (attrsType != nil && attrsType.IsVoidType()) {
		// goog.dom.createDom('iframe', null) is fine.
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	hasExtraChildren := call.NumArgs() > 2

	for _, tagAttr := range r.bannedTagAttrs {
		if !tagNamesIsNull && !slices.Contains(tagNames, tagAttr[0]) && tagAttr[0] != "*" {
			continue
		}
		var violation *ConformanceRules_ConformanceResult
		if !tagNamesIsNull || tagAttr[0] == "*" {
			violation = ConformanceRules_ConformanceResult_VIOLATION
		} else if r.reportLooseTypeViolations {
			violation = ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION
		} else {
			violation = ConformanceRules_ConformanceResult_CONFORMANCE
		}
		if isClassName {
			if tagAttr[1] != "class" {
				continue
			}
			return violation
		}
		if tagAttr[1] == "textContent" &&
			hasExtraChildren &&
			violation != ConformanceRules_ConformanceResult_CONFORMANCE {
			return violation
		}
		objLit, isObjLit := ast.AsObjectLit(attrsNode)
		if !isObjLit {
			// Attrs is not an object literal and tagName matches or is unknown.
			if r.reportLooseTypeViolations {
				return ConformanceRules_ConformanceResult_POSSIBLE_VIOLATION
			}
			return ConformanceRules_ConformanceResult_CONFORMANCE
		}
		prop := NodeUtil_GetFirstPropMatchingKey(attrsNode, tagAttr[1])
		if !prop.IsNil() {
			if NodeUtil_IsSomeCompileTimeConstStringValue(prop) {
				// Ignore string literal values.
				continue
			}
			return violation
		}

		for member := range objLit.Members() {
			if _, isComputed := member.AsComputedProp(); isComputed {
				// We don't know if the computed property matches 'src' or not.
				if r.reportLooseTypeViolations {
					return conformancePossibleViolationLooseTypes
				}
				return ConformanceRules_ConformanceResult_CONFORMANCE
			}
		}
	}

	return ConformanceRules_ConformanceResult_CONFORMANCE
}

func (r *conformanceBanCreateDomRule) tagNames(tag rhino.Node) (tagNames []string, isNull bool) {
	if s, ok := ast.AsStringLit(tag); ok {
		return []string{jutil.StringToLowerCaseRoot(s.Value())}, false
	} else if gp, ok := ast.AsNormalOrOptChainGetProp(tag); ok && googDomTagNameQName.Matches(gp.Object().Node()) {
		return []string{jutil.StringToLowerCaseRoot(gp.Property())}, false
	}
	type_ := tag.GetJSType()
	if type_ == nil || !type_.IsTemplatizedType() {
		return nil, true
	}
	typeAsObj := type_.ToMaybeObjectType()
	if typeAsObj != nil && typeAsObj.GetRawType() != nil {
		if rawDisp, ok := typeAsObj.GetRawType().DisplayName(); ok && rawDisp == "goog.dom.TagName" {
			tagType := onlyJSTypeElement(typeAsObj.GetTemplateTypes())
			var names []string
			if displayName, ok := tagType.DisplayName(); ok {
				names = htmlElementTagNames[displayName]
			}
			return names, false
		}
	}
	return nil, true
}

func onlyJSTypeElement(types []rhino.JSType) rhino.JSType {
	if len(types) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	if len(types) == 1 {
		return types[0]
	}
	var sb strings.Builder
	sb.WriteString("expected one element but was: <")
	sb.WriteString(types[0].String())
	i := 1
	for ; i < 5 && i < len(types); i++ {
		sb.WriteString(", ")
		sb.WriteString(types[i].String())
	}
	if i < len(types) {
		sb.WriteString(", ...")
	}
	sb.WriteByte('>')
	panic(jutil.NewIllegalArgumentExceptionMsg(sb.String()))
}

var htmlElementTagNames = map[string][]string{
	"HTMLAnchorElement":       {"a"},
	"HTMLAppletElement":       {"applet"},
	"HTMLAreaElement":         {"area"},
	"HTMLAudioElement":        {"audio"},
	"HTMLBRElement":           {"br"},
	"HTMLBaseElement":         {"base"},
	"HTMLBaseFontElement":     {"basefont"},
	"HTMLBodyElement":         {"body"},
	"HTMLButtonElement":       {"button"},
	"HTMLCanvasElement":       {"canvas"},
	"HTMLDListElement":        {"dl"},
	"HTMLDataListElement":     {"datalist"},
	"HTMLDetailsElement":      {"details"},
	"HTMLDialogElement":       {"dialog"},
	"HTMLDirectoryElement":    {"dir"},
	"HTMLDivElement":          {"div"},
	"HTMLEmbedElement":        {"embed"},
	"HTMLFieldSetElement":     {"fieldset"},
	"HTMLFontElement":         {"font"},
	"HTMLFormElement":         {"form"},
	"HTMLFrameElement":        {"frame"},
	"HTMLFrameSetElement":     {"frameset"},
	"HTMLHRElement":           {"hr"},
	"HTMLHeadElement":         {"head"},
	"HTMLHeadingElement":      {"h1", "h2", "h3", "h4", "h5", "h6"},
	"HTMLHtmlElement":         {"html"},
	"HTMLIFrameElement":       {"iframe"},
	"HTMLImageElement":        {"img"},
	"HTMLInputElement":        {"input"},
	"HTMLIsIndexElement":      {"isindex"},
	"HTMLLIElement":           {"li"},
	"HTMLLabelElement":        {"label"},
	"HTMLLegendElement":       {"legend"},
	"HTMLLinkElement":         {"link"},
	"HTMLMapElement":          {"map"},
	"HTMLMenuElement":         {"menu"},
	"HTMLMetaElement":         {"meta"},
	"HTMLMeterElement":        {"meter"},
	"HTMLModElement":          {"del", "ins"},
	"HTMLOListElement":        {"ol"},
	"HTMLObjectElement":       {"object"},
	"HTMLOptGroupElement":     {"optgroup"},
	"HTMLOptionElement":       {"option"},
	"HTMLOutputElement":       {"output"},
	"HTMLParagraphElement":    {"p"},
	"HTMLParamElement":        {"param"},
	"HTMLPreElement":          {"pre"},
	"HTMLProgressElement":     {"progress"},
	"HTMLQuoteElement":        {"blockquote", "q"},
	"HTMLScriptElement":       {"script"},
	"HTMLSelectElement":       {"select"},
	"HTMLSourceElement":       {"source"},
	"HTMLSpanElement":         {"span"},
	"HTMLStyleElement":        {"style"},
	"HTMLTableCaptionElement": {"caption"},
	"HTMLTableCellElement":    {"td", "th"},
	"HTMLTableColElement":     {"col", "colgroup"},
	"HTMLTableElement":        {"table"},
	"HTMLTableRowElement":     {"tr"},
	"HTMLTableSectionElement": {"tbody", "tfoot", "thead"},
	"HTMLTemplateElement":     {"template"},
	"HTMLTextAreaElement":     {"textarea"},
	"HTMLTitleElement":        {"title"},
	"HTMLTrackElement":        {"track"},
	"HTMLUListElement":        {"ul"},
	"HTMLVideoElement":        {"video"},
}

func (r *conformanceBanCreateDomRule) asCreateDomCall(n rhino.Node) (ast.NormalOrOptChainCall, bool) {
	call, ok := ast.AsNormalOrOptChainCall(n)
	if !ok {
		return ast.NormalOrOptChainCall{}, false
	}
	target, ok := ast.AsNormalOrOptChainGetProp(call.Callee().Node())
	if !ok || target.Property() != "createDom" {
		return ast.NormalOrOptChainCall{}, false
	}

	srcObj := target.Object().Node()
	if googDomQName.Matches(srcObj) {
		return call, true
	}
	type_ := srcObj.GetJSType()
	if type_ != nil && r.domHelperType != nil {
		if call.Node().IsOptChainCall() || target.Node().IsOptChainGetProp() {
			type_ = type_.RestrictByNotNullOrUndefined()
		}
		if type_.Equals(r.domHelperType) {
			return call, true
		}
	}
	return ast.NormalOrOptChainCall{}, false
}

// conformanceSecuritySensitiveAttributes checks attribute-name expressions
// against a blocklist of security-sensitive HTML attributes and event handlers.
type conformanceSecuritySensitiveAttributes struct {
	bannedAttrs             map[string]struct{}
	globalNamespaceSupplier func() *GlobalNamespace
}

var allBannedSecurityAttrs = [...]string{
	"href",
	"rel",
	"src",
	"srcdoc",
	"action",
	"formaction",
	"sandbox",
	"icon",
	"codebase",
	"data",
}

func newDefaultSecuritySensitiveAttributes() *conformanceSecuritySensitiveAttributes {
	return newSecuritySensitiveAttributes(allBannedSecurityAttrs[:], func() *GlobalNamespace { return nil })
}

func newSecuritySensitiveAttributes(
	bannedAttrs []string, globalNamespaceSupplier func() *GlobalNamespace) *conformanceSecuritySensitiveAttributes {
	m := make(map[string]struct{}, len(bannedAttrs))
	for _, a := range bannedAttrs {
		m[jutil.StringToLowerCaseRoot(a)] = struct{}{}
	}
	return &conformanceSecuritySensitiveAttributes{
		bannedAttrs:             m,
		globalNamespaceSupplier: globalNamespaceSupplier,
	}
}

func (s *conformanceSecuritySensitiveAttributes) contains(attributeName string) bool {
	_, ok := s.bannedAttrs[attributeName]
	return ok
}

func (s *conformanceSecuritySensitiveAttributes) checkAttributeName(
	traversal *NodeTraversal, attrName rhino.Node) *ConformanceRules_ConformanceResult {
	var scope *Scope
	if traversal != nil {
		scope = traversal.GetScope()
	}
	literalName, hasLiteral :=
		inferConformanceStringValue(scope, attrName, s.globalNamespaceSupplier)
	if !hasLiteral {
		// xid() obfuscates attribute names, thus never clashing with security-sensitive attributes.
		if !attrName.IsNil() && isXidType(attrName.GetJSType()) {
			return ConformanceRules_ConformanceResult_CONFORMANCE
		}
		return ConformanceRules_ConformanceResult_VIOLATION
	}

	lowerName := jutil.StringToLowerCaseRoot(literalName)
	if s.contains(lowerName) || isEventHandlerAttrName(lowerName) {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

func (s *conformanceSecuritySensitiveAttributes) checkAttributeNameHighConfidence(
	traversal *NodeTraversal, attrName rhino.Node) *ConformanceRules_ConformanceResult {
	var scope *Scope
	if traversal != nil {
		scope = traversal.GetScope()
	}
	literalName, hasLiteral :=
		inferConformanceStringValue(scope, attrName, s.globalNamespaceSupplier)
	if !hasLiteral {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	if s.contains(jutil.StringToLowerCaseRoot(literalName)) {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

type conformanceBanElementSetAttributeRule struct {
	conformanceAbstractRule

	performGlobalNamespaceAnalysis bool
	globalNamespace                *GlobalNamespace

	elementType                      rhino.JSType
	securitySensitiveAttributes      *conformanceSecuritySensitiveAttributes
	defaultDecisionForUncertainCases *ConformanceRules_ConformanceResult
}

const (
	elementTypeName                     = "Element"
	propSetAttribute                    = "setAttribute"
	propSetAttributeNS                  = "setAttributeNS"
	propSetAttributeNode                = "setAttributeNode"
	propSetAttributeNodeNS              = "setAttributeNodeNS"
	setAttrGlobalNamespaceAnalysisLimit = 10000
)

func isBannedSetAttributeProp(name string) bool {
	switch name {
	case propSetAttribute,
		propSetAttributeNS,
		propSetAttributeNode,
		propSetAttributeNodeNS:
		return true
	default:
		return false
	}
}

func newConformanceBanElementSetAttributeRule(compiler AbstractCompiler, requirement *Requirement) (*conformanceBanElementSetAttributeRule, error) {
	r := &conformanceBanElementSetAttributeRule{
		performGlobalNamespaceAnalysis: true,
	}
	if err := r.initRule(compiler, requirement, false, r.checkNode); err != nil {
		return nil, err
	}
	r.securitySensitiveAttributes =
		newSecuritySensitiveAttributes(requirement.GetValueList(), r.getGlobalNamespace)
	if requirement.GetReportLooseTypeViolations() {
		r.defaultDecisionForUncertainCases = ConformanceRules_ConformanceResult_VIOLATION
	} else {
		r.defaultDecisionForUncertainCases = ConformanceRules_ConformanceResult_CONFORMANCE
	}
	r.elementType = compiler.GetTypeRegistry().GetGlobalType(elementTypeName)
	return r, nil
}

func (r *conformanceBanElementSetAttributeRule) getGlobalNamespace() *GlobalNamespace {
	if r.performGlobalNamespaceAnalysis && r.globalNamespace == nil {
		if NodeUtil_CountAstSizeUpToLimit(r.compiler.GetJsRoot(), setAttrGlobalNamespaceAnalysisLimit) >=
			setAttrGlobalNamespaceAnalysisLimit {
			r.performGlobalNamespaceAnalysis = false
		} else {
			r.globalNamespace = NewGlobalNamespace(r.compiler, r.compiler.GetJsRoot())
		}
	}
	return r.globalNamespace
}

func (r *conformanceBanElementSetAttributeRule) checkNode(traversal *NodeTraversal, node rhino.Node) *ConformanceRules_ConformanceResult {
	if call, ok := ast.AsNormalOrOptChainCall(node); ok {
		return r.checkPropertyCall(traversal, call)
	} else if ge, ok := ast.AsNormalOrOptChainGetElem(node); ok && NodeUtil_IsLValue(node) {
		return r.checkGetElement(traversal, ge)
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

func (r *conformanceBanElementSetAttributeRule) checkPropertyCall(
	traversal *NodeTraversal, call ast.NormalOrOptChainCall) *ConformanceRules_ConformanceResult {
	calledProperty, ok := r.bannedPropertyName(call)
	if !ok {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	cn := call.AsCallOrNew()
	if calledProperty == propSetAttribute {
		firstArg, hasFirst := cn.ArgAt(0)
		if !hasFirst || (call.NumArgs() < 2 && !firstArg.Node().IsSpread()) {
			return ConformanceRules_ConformanceResult_CONFORMANCE
		}
		if r.requirement.GetReportLooseTypeViolations() {
			return r.securitySensitiveAttributes.checkAttributeName(traversal, firstArg.Node())
		}
		return r.securitySensitiveAttributes.checkAttributeNameHighConfidence(traversal, firstArg.Node())
	}

	if calledProperty == propSetAttributeNS {
		if call.NumArgs() < 3 {
			return ConformanceRules_ConformanceResult_CONFORMANCE
		}
		firstArg, _ := cn.ArgAt(0)
		secondArg, _ := cn.ArgAt(1)
		if firstArg.Node().IsNull() {
			// A null namespace makes this equivalent to a regular setAttribute call.
			if r.requirement.GetReportLooseTypeViolations() {
				return r.securitySensitiveAttributes.checkAttributeName(traversal, secondArg.Node())
			}
			return r.securitySensitiveAttributes.checkAttributeNameHighConfidence(traversal, secondArg.Node())
		}
	}

	// We haven't defined a security contract on setAttributeNode and setAttributeNodeNS yet, so
	// flag them as risky if report_loose_type_violations is set.
	return r.defaultDecisionForUncertainCases
}

func (r *conformanceBanElementSetAttributeRule) bannedPropertyName(call ast.NormalOrOptChainCall) (string, bool) {
	target, ok := ast.AsNormalOrOptChainGetProp(call.Callee().Node())
	if !ok {
		return "", false
	}
	propertyName := target.Property()
	if !isBannedSetAttributeProp(propertyName) {
		return "", false
	}
	type_ := target.Object().Node().GetJSType()
	if type_ == nil {
		return "", false
	}
	if r.elementType != nil && type_.RestrictByNotNullOrUndefined().IsSubtypeOf(r.elementType) {
		return propertyName, true
	}
	return "", false
}

func (r *conformanceBanElementSetAttributeRule) checkGetElement(
	traversal *NodeTraversal, ge ast.NormalOrOptChainGetElem) *ConformanceRules_ConformanceResult {
	key := ge.Index().Node()
	var scope *Scope
	if traversal != nil {
		scope = traversal.GetScope()
	}
	if keyName, ok := inferConformanceStringValue(scope, key, r.getGlobalNamespace); ok {
		if r.securitySensitiveAttributes.contains(jutil.StringToLowerCaseRoot(keyName)) ||
			(r.requirement.GetReportLooseTypeViolations() &&
				(keyName == "innerHTML" || keyName == "outerHTML")) {
			if r.hasElementType(ge) {
				return ConformanceRules_ConformanceResult_VIOLATION
			}
		}
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	if !r.requirement.GetReportLooseTypeViolations() {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	keyType := key.GetJSType()
	if keyType == nil || isXidType(keyType) {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	if r.hasElementType(ge) {
		if keyType.IsString() {
			return ConformanceRules_ConformanceResult_VIOLATION
		} else if keyType.IsUnionType() {
			for _, alternate := range keyType.ToMaybeUnionType().GetAlternates() {
				if alternate.IsString() {
					return ConformanceRules_ConformanceResult_VIOLATION
				}
			}
		}
	}

	return ConformanceRules_ConformanceResult_CONFORMANCE
}

func (r *conformanceBanElementSetAttributeRule) hasElementType(ge ast.NormalOrOptChainGetElem) bool {
	objType := ge.Object().Node().GetJSType()
	if objType == nil || r.elementType == nil {
		return false
	}
	objType = objType.RestrictByNotNullOrUndefined()
	if objType.IsUnknownType() || objType.IsEmptyType() || !objType.IsSubtypeOf(r.elementType) {
		return false
	}
	return true
}

type conformanceBanSettingAttributesRule struct {
	conformanceAbstractRule

	restrictions []attributeSettingRestriction
}

type attributeSettingRestriction struct {
	type_    rhino.JSType
	property string
}

var defaultBannedSettingAttrs = newDefaultSecuritySensitiveAttributes()

func newConformanceBanSettingAttributesRule(compiler AbstractCompiler, requirement *Requirement) (*conformanceBanSettingAttributesRule, error) {
	r := &conformanceBanSettingAttributesRule{}
	if err := r.initRule(compiler, requirement, false, r.checkNode); err != nil {
		return nil, err
	}

	values := requirement.GetValueList()
	if len(values) == 0 {
		return nil, errInvalidRequirementSpec("missing value")
	}

	registry := compiler.GetTypeRegistry()
	restrictions := make([]attributeSettingRestriction, 0, len(values))
	for _, value := range values {
		typeName, property, ok := splitConformanceClassAndProperty(value)
		if !ok {
			return nil, errInvalidRequirementSpec("bad prop value")
		}
		restrictions = append(restrictions, attributeSettingRestriction{
			type_:    registry.GetGlobalType(typeName),
			property: property,
		})
	}

	r.restrictions = restrictions
	return r, nil
}

func (r *conformanceBanSettingAttributesRule) checkNode(traversal *NodeTraversal, node rhino.Node) *ConformanceRules_ConformanceResult {
	if call, ok := ast.AsNormalOrOptChainCall(node); ok {
		return r.checkPropertyCall(traversal, call)
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

func (r *conformanceBanSettingAttributesRule) checkPropertyCall(traversal *NodeTraversal, call ast.NormalOrOptChainCall) *ConformanceRules_ConformanceResult {
	if !r.matchesBannedProperty(call) {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}
	firstArg, hasFirst := call.FirstArg()
	if !hasFirst || (call.NumArgs() < 2 && !firstArg.Node().IsSpread()) {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	return defaultBannedSettingAttrs.checkAttributeName(traversal, firstArg.Node())
}

func (r *conformanceBanSettingAttributesRule) matchesBannedProperty(call ast.NormalOrOptChainCall) bool {
	target, ok := ast.AsNormalOrOptChainGetProp(call.Callee().Node())
	if !ok {
		return false
	}
	type_ := target.Object().Node().GetJSType()
	if type_ == nil {
		return false
	}
	if call.Node().IsOptChainCall() || target.Node().IsOptChainGetProp() {
		type_ = type_.RestrictByNotNullOrUndefined()
	}
	propertyName := target.Property()
	for i := range r.restrictions {
		restricted := &r.restrictions[i]
		if propertyName == restricted.property &&
			restricted.type_ != nil &&
			type_.IsSubtypeOf(restricted.type_) {
			return true
		}
	}
	return false
}

type conformanceBanExecCommandRule struct {
	conformanceAbstractRule

	bannedAttrs []string
}

func newConformanceBanExecCommandRule(compiler AbstractCompiler, requirement *Requirement) (*conformanceBanExecCommandRule, error) {
	r := &conformanceBanExecCommandRule{}
	if err := r.initRule(compiler, requirement, false, r.checkNode); err != nil {
		return nil, err
	}
	values := requirement.GetValueList()
	bannedAttrs := make([]string, 0, len(values))
	for _, v := range values {
		bannedAttrs = append(bannedAttrs, jutil.StringToLowerCaseRoot(v))
	}
	r.bannedAttrs = bannedAttrs
	return r, nil
}

func (r *conformanceBanExecCommandRule) checkNode(traversal *NodeTraversal, node rhino.Node) *ConformanceRules_ConformanceResult {
	call, ok := r.asBannedExecCommandCall(node)
	if !ok {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	firstArg, hasFirst := call.FirstArg()
	if !hasFirst {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	attr, ok := ast.AsStringLit(firstArg.Node())
	if !ok {
		return ConformanceRules_ConformanceResult_VIOLATION
	}

	if slices.Contains(r.bannedAttrs, jutil.StringToLowerCaseRoot(attr.Value())) {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

func (r *conformanceBanExecCommandRule) asBannedExecCommandCall(node rhino.Node) (ast.NormalOrOptChainCall, bool) {
	call, ok := ast.AsNormalOrOptChainCall(node)
	if !ok {
		return ast.NormalOrOptChainCall{}, false
	}
	target, ok := ast.AsNormalOrOptChainGetProp(call.Callee().Node())
	if !ok || target.Property() != "execCommand" {
		return ast.NormalOrOptChainCall{}, false
	}
	type_ := target.Object().Node().GetJSType()
	if type_ == nil {
		return ast.NormalOrOptChainCall{}, false
	}
	documentType := r.compiler.GetTypeRegistry().GetGlobalType("Document")
	if documentType == nil || !type_.RestrictByNotNullOrUndefined().IsSubtypeOf(documentType) {
		return ast.NormalOrOptChainCall{}, false
	}
	return call, true
}

// conformanceBanStaticThisRule checks that `this` is not referenced directly
// inside a static class method.
type conformanceBanStaticThisRule struct {
	conformanceAbstractRule
}

func newConformanceBanStaticThisRule(compiler AbstractCompiler, requirement *Requirement) (*conformanceBanStaticThisRule, error) {
	r := &conformanceBanStaticThisRule{}
	if err := r.initTypeRule(compiler, requirement, false, r.checkNode); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *conformanceBanStaticThisRule) checkNode(t *NodeTraversal, n rhino.Node) *ConformanceRules_ConformanceResult {
	if n.IsNil() || !n.IsThis() {
		return ConformanceRules_ConformanceResult_CONFORMANCE
	}

	enclosingFunction := NodeUtil_GetEnclosingNonArrowFunction(n)
	if isStaticClassMethod(enclosingFunction) {
		return ConformanceRules_ConformanceResult_VIOLATION
	}
	return ConformanceRules_ConformanceResult_CONFORMANCE
}

func isStaticClassMethod(n rhino.Node) bool {
	if n.IsNil() || !n.IsFunction() {
		return false
	}
	parent := n.GetParent()
	if mfd, ok := ast.AsMemberFunctionDef(parent); ok {
		return mfd.IsStatic()
	}
	if cp, ok := ast.AsComputedProp(parent); ok && ast.RoleOf(parent) == ast.RoleClassBodyMembers {
		return cp.IsStatic()
	}
	return false
}
