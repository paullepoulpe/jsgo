// Port of com/google/javascript/jscomp/JsMessageVisitor.java @ closure-compiler v20261021

/*
 * Copyright 2008 The Closure Compiler Authors.
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
	"errors"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// jsMessageVisitor_extractStringFromStringExprNode returns the string value
// associated with a JS string literal, no-substitution template literal, or
// string concatenation.
//
// Java: JsMessageVisitor.extractStringFromStringExprNode(Node)
func jsMessageVisitor_extractStringFromStringExprNode(node rhino.Node) (string, *malformedMessageError) {
	switch node.GetToken() {
	case rhino.Token_STRINGLIT:
		return ast.MustStringLit(node).Value(), nil
	case rhino.Token_TEMPLATELIT:
		if node.HasOneChild() {
			return *jutil.CheckNotNull(node.GetFirstChild().GetCookedString()), nil
		}
		return "", &malformedMessageError{
			msg:  "Template literals with substitutions are not allowed.",
			node: node,
		}
	case rhino.Token_ADD:
		var sb strings.Builder
		for child := node.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			part, err := jsMessageVisitor_extractStringFromStringExprNode(child)
			if err != nil {
				return "", err
			}
			sb.WriteString(part)
		}
		return sb.String(), nil
	default:
		return "", &malformedMessageError{msg: "literal string or concatenation expected", node: node}
	}
}

// extractJsMessageDefinition extracts data from a call to goog.getMsg(...).
//
// Java: JsMessageVisitor.extractJsMessageDefinition(Node, JSDocInfo, String, String)
func (v *JsMessageVisitor) extractJsMessageDefinition(
	msgCall ast.Call, jsDocInfo *rhino.JSDocInfo, messageKeyFromLhs string, sourceName string,
) (JsMessageDefinition, *malformedMessageError) {
	msgNode := msgCall.Node()
	var description, meaning, alternateMessageId string
	var hasDescription, hasMeaning, hasAlternateMessageId bool
	if jsDocInfo != nil {
		if s := jsDocInfo.GetDescription(); s != nil {
			description = *s
			hasDescription = true
		}
		if s := jsDocInfo.GetMeaning(); s != nil {
			meaning = *s
			hasMeaning = true
		}
		if s := jsDocInfo.GetAlternateMessageId(); s != nil {
			alternateMessageId = *s
			hasAlternateMessageId = true
		}
	}

	firstArg, hasFirstArg := msgCall.FirstArg()
	if !hasFirstArg {
		return JsMessageDefinition{}, &malformedMessageError{msg: "Message string literal expected", node: msgNode}
	}
	msgTextNode := firstArg.Node()
	valuesObjLit := msgTextNode.GetNext()
	var optionsBagArgument rhino.Node
	if !valuesObjLit.IsNil() {
		optionsBagArgument = valuesObjLit.GetNext()
	}
	if !optionsBagArgument.IsNil() {
		if unexpectedArgument := optionsBagArgument.GetNext(); !unexpectedArgument.IsNil() {
			return JsMessageDefinition{}, &malformedMessageError{msg: "too many arguments", node: unexpectedArgument}
		}
	}

	msgTextString, malErr := jsMessageVisitor_extractStringFromStringExprNode(msgTextNode)
	if malErr != nil {
		return JsMessageDefinition{}, malErr
	}
	googGetMsgParsedText, fmtErr := jsMessageVisitor_extractGoogGetMsgParsedText(msgTextString)
	if fmtErr != nil {
		return JsMessageDefinition{}, &malformedMessageError{msg: fmtErr.Error(), node: msgTextNode}
	}

	placeholderObjectLiteralMap, malErr := extractObjectLiteralMap(valuesObjLit)
	if malErr != nil {
		return JsMessageDefinition{}, malErr
	}

	if malErr := placeholderObjectLiteralMap.checkForRequiredKeys(
		googGetMsgParsedText.placeholderOrder, msgTextNode, "Unrecognized message placeholder referenced: ",
	); malErr != nil {
		return JsMessageDefinition{}, malErr
	}
	if malErr := placeholderObjectLiteralMap.checkForUnexpectedKeys(
		googGetMsgParsedText.placeholderSet, "Unused message placeholder: ",
	); malErr != nil {
		return JsMessageDefinition{}, malErr
	}

	placeholderValuesMap, placeholderValuesOrder := placeholderObjectLiteralMap.extractAsValueMapAndOrder()

	jsMessageOptions, malErr := extractJsMessageOptions(optionsBagArgument)
	if malErr != nil {
		return JsMessageDefinition{}, malErr
	}
	if malErr := jsMessageOptions.checkForUnknownPlaceholders(googGetMsgParsedText.placeholderSet); malErr != nil {
		return JsMessageDefinition{}, malErr
	}

	jsExtractedMessage := v.buildJsMessage(jsMessageBuildInput{
		messageKeyFromLhs:            messageKeyFromLhs,
		sourceName:                   sourceName,
		messageText:                  googGetMsgParsedText.originalMsgText,
		messageParts:                 googGetMsgParsedText.parts,
		placeholderExampleMap:        jsMessageOptions.placeholderExamplesMap,
		placeholderExampleOrder:      jsMessageOptions.placeholderExamplesOrder,
		placeholderOriginalCodeMap:   jsMessageOptions.placeholderOriginalCodeMap,
		placeholderOriginalCodeOrder: jsMessageOptions.placeholderOriginalCodeOrder,
		description:                  description,
		hasDescription:               hasDescription,
		meaning:                      meaning,
		hasMeaning:                   hasMeaning,
		alternateMessageId:           alternateMessageId,
		hasAlternateMessageId:        hasAlternateMessageId,
	})

	return JsMessageDefinition{
		message:               jsExtractedMessage,
		messageNode:           msgNode,
		templateTextNode:      msgTextNode,
		placeholderValuesNode: valuesObjLit,
		placeholderValueMap:   placeholderValuesMap,
		placeholderValueOrder: placeholderValuesOrder,
		escapeLessThan:        jsMessageOptions.isEscapeLessThan,
		unescapeHtmlEntities:  jsMessageOptions.isUnescapeHtmlEntities,
	}, nil
}

// extractIcuTemplateDefinition extracts message data from a
// `declareIcuTemplate()` call.
//
// Java: JsMessageVisitor.extractIcuTemplateDefinition(Node, JSDocInfo, String, String)
func (v *JsMessageVisitor) extractIcuTemplateDefinition(
	msgCall ast.Call, jsDocInfo *rhino.JSDocInfo, messageKeyFromLhs string, sourceName string,
) (IcuTemplateDefinition, *malformedMessageError) {
	msgNode := msgCall.Node()
	if jsDocInfo != nil {
		if jsDocInfo.GetAlternateMessageId() != nil {
			return IcuTemplateDefinition{}, &malformedMessageError{
				msg:  "Use the 'alternate_message_id' option, not the '@alternateMessageId' annotation",
				node: msgNode,
			}
		}
		if jsDocInfo.GetDescription() != nil {
			return IcuTemplateDefinition{}, &malformedMessageError{
				msg:  "Use the 'description' option, not the '@desc' annotation",
				node: msgNode,
			}
		}
		if jsDocInfo.GetMeaning() != nil {
			return IcuTemplateDefinition{}, &malformedMessageError{
				msg:  "Use the 'meaning' option, not the '@meaning' annotation",
				node: msgNode,
			}
		}
	}

	firstArg, hasFirstArg := msgCall.FirstArg()
	if !hasFirstArg {
		return IcuTemplateDefinition{}, &malformedMessageError{msg: "message string argument expected", node: msgNode}
	}
	stringLiteralExpression := firstArg.Node()
	templateString, malErr := jsMessageVisitor_extractStringFromStringExprNode(stringLiteralExpression)
	if malErr != nil {
		return IcuTemplateDefinition{}, malErr
	}
	icuMessageTemplateString := JsMessageVisitor_IcuMessageTemplateString{template: templateString}

	optionsBagNode := stringLiteralExpression.GetNext()
	if optionsBagNode.IsNil() {
		return IcuTemplateDefinition{}, &malformedMessageError{msg: "options argument expected", node: msgNode}
	}
	icuOptions, malErr := extractIcuTemplateOptions(optionsBagNode)
	if malErr != nil {
		return IcuTemplateDefinition{}, malErr
	}

	extractedIcuTemplateParts := icuMessageTemplateString.ExtractParts(icuOptions.placeholderNames)
	if malErr := icuOptions.checkForUnknownPlaceholders(extractedIcuTemplateParts.extractedPlaceholderNames); malErr != nil {
		return IcuTemplateDefinition{}, malErr
	}

	if unexpectedArgument := optionsBagNode.GetNext(); !unexpectedArgument.IsNil() {
		return IcuTemplateDefinition{}, &malformedMessageError{msg: "too many arguments", node: unexpectedArgument}
	}

	extractedMessage := v.buildJsMessage(jsMessageBuildInput{
		messageKeyFromLhs:            messageKeyFromLhs,
		sourceName:                   sourceName,
		messageText:                  icuMessageTemplateString.template,
		messageParts:                 extractedIcuTemplateParts.extractedParts,
		placeholderExampleMap:        icuOptions.placeholderExamplesMap,
		placeholderExampleOrder:      icuOptions.placeholderExamplesOrder,
		placeholderOriginalCodeMap:   icuOptions.placeholderOriginalCodeMap,
		placeholderOriginalCodeOrder: icuOptions.placeholderOriginalCodeOrder,
		description:                  icuOptions.description,
		hasDescription:               true,
		meaning:                      icuOptions.meaning,
		hasMeaning:                   icuOptions.hasMeaning,
		alternateMessageId:           icuOptions.alternateMessageId,
		hasAlternateMessageId:        icuOptions.hasAlternateMessageId,
	})

	return IcuTemplateDefinition{
		message:          extractedMessage,
		messageNode:      msgNode,
		templateTextNode: stringLiteralExpression,
	}, nil
}

// JsMessageVisitor_IcuMessageTemplateString is the template string of a
// declareIcuTemplate() call.
type JsMessageVisitor_IcuMessageTemplateString struct {
	template string
}

// Java: JsMessageVisitor.IcuMessageTemplateString(String)
func NewJsMessageVisitor_IcuMessageTemplateString(template string) JsMessageVisitor_IcuMessageTemplateString {
	return JsMessageVisitor_IcuMessageTemplateString{template: template}
}

// ExtractParts splits the template string into plain string parts and
// placeholders named in placeholderNamesToExtract.
//
// Java: JsMessageVisitor.IcuMessageTemplateString.extractParts(Set<String>)
func (s JsMessageVisitor_IcuMessageTemplateString) ExtractParts(
	placeholderNamesToExtract map[string]struct{},
) JsMessageVisitor_ExtractedIcuTemplateParts {
	tpl := s.template
	if len(tpl) == 0 {
		return JsMessageVisitor_ExtractedIcuTemplateParts{}
	}
	if len(placeholderNamesToExtract) == 0 {
		return JsMessageVisitor_ExtractedIcuTemplateParts{
			extractedParts: []JsMessage_Part{JsMessage_StringPart_Create(tpl)},
		}
	}

	var result JsMessageVisitor_ExtractedIcuTemplateParts
	remainingStart := 0
	for i := 0; i < len(tpl); {
		rel := strings.IndexByte(tpl[i:], '{')
		if rel < 0 {
			break
		}
		openIdx := i + rel
		j := openIdx + 1
		for j < len(tpl) {
			c := tpl[j]
			if (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
				j++
			} else {
				break
			}
		}
		if j > openIdx+1 && j < len(tpl) && tpl[j] == '}' {
			placeholderName := tpl[openIdx+1 : j]
			if _, ok := placeholderNamesToExtract[placeholderName]; ok {
				if result.extractedPlaceholderNames == nil {
					result.extractedPlaceholderNames = make(map[string]struct{})
				}
				result.extractedPlaceholderNames[placeholderName] = struct{}{}
				result.extractedParts = append(result.extractedParts,
					JsMessage_StringPart_Create(tpl[remainingStart:openIdx]),
					JsMessage_PlaceholderReference_CreateForCanonicalName(placeholderName),
				)
				remainingStart = j + 1
			}
			i = j + 1
		} else {
			i = j
			if i == openIdx+1 && i < len(tpl) && tpl[i] != '{' {
				i++
			}
		}
	}
	if remainingStart < len(tpl) {
		result.extractedParts = append(result.extractedParts, JsMessage_StringPart_Create(tpl[remainingStart:]))
	}
	return result
}

// JsMessageVisitor_ExtractedIcuTemplateParts is the result of
// IcuMessageTemplateString.ExtractParts.
type JsMessageVisitor_ExtractedIcuTemplateParts struct {
	extractedPlaceholderNames map[string]struct{}
	extractedParts            []JsMessage_Part
}

func isJsMessageOptionName(key string) bool {
	switch key {
	case "html", "unescapeHtmlEntities", "example", "original_code":
		return true
	}
	return false
}

type jsMessageOptions struct {
	placeholderExamplesMap       map[string]string
	placeholderExamplesOrder     []string
	placeholderOriginalCodeMap   map[string]string
	placeholderOriginalCodeOrder []string
	exampleObjectLiteralMap      jsMessageObjectLiteralMap
	originalCodeObjectLiteralMap jsMessageObjectLiteralMap
	isEscapeLessThan             bool
	isUnescapeHtmlEntities       bool
}

func (o *jsMessageOptions) checkForUnknownPlaceholders(knownPlaceholders map[string]struct{}) *malformedMessageError {
	if err := o.exampleObjectLiteralMap.checkForUnexpectedKeys(knownPlaceholders, "Unknown placeholder: "); err != nil {
		return err
	}
	return o.originalCodeObjectLiteralMap.checkForUnexpectedKeys(knownPlaceholders, "Unknown placeholder: ")
}

// Java: JsMessageVisitor.extractJsMessageOptions(Node)
func extractJsMessageOptions(optionsBag rhino.Node) (jsMessageOptions, *malformedMessageError) {
	if optionsBag.IsNil() {
		return jsMessageOptions{}, nil
	}
	objectLiteralMap, err := extractObjectLiteralMap(optionsBag)
	if err != nil {
		return jsMessageOptions{}, err
	}
	if err := objectLiteralMap.checkForUnexpectedKeysFn(isJsMessageOptionName, "Unknown option: "); err != nil {
		return jsMessageOptions{}, err
	}

	isEscapeLessThan, err := objectLiteralMap.getBooleanValueOrFalse("html")
	if err != nil {
		return jsMessageOptions{}, err
	}
	isUnescapeHtmlEntities, err := objectLiteralMap.getBooleanValueOrFalse("unescapeHtmlEntities")
	if err != nil {
		return jsMessageOptions{}, err
	}

	exampleObjectLiteralMap, err := extractObjectLiteralMap(objectLiteralMap.getValueNode("example"))
	if err != nil {
		return jsMessageOptions{}, err
	}
	placeholderExamplesMap, placeholderExamplesOrder, err := exampleObjectLiteralMap.extractAsStringToStringMap()
	if err != nil {
		return jsMessageOptions{}, err
	}

	originalCodeObjectLiteralMap, err := extractObjectLiteralMap(objectLiteralMap.getValueNode("original_code"))
	if err != nil {
		return jsMessageOptions{}, err
	}
	placeholderOriginalCodeMap, placeholderOriginalCodeOrder, err := originalCodeObjectLiteralMap.extractAsStringToStringMap()
	if err != nil {
		return jsMessageOptions{}, err
	}

	return jsMessageOptions{
		isEscapeLessThan:             isEscapeLessThan,
		isUnescapeHtmlEntities:       isUnescapeHtmlEntities,
		placeholderExamplesMap:       placeholderExamplesMap,
		placeholderExamplesOrder:     placeholderExamplesOrder,
		placeholderOriginalCodeMap:   placeholderOriginalCodeMap,
		placeholderOriginalCodeOrder: placeholderOriginalCodeOrder,
		exampleObjectLiteralMap:      exampleObjectLiteralMap,
		originalCodeObjectLiteralMap: originalCodeObjectLiteralMap,
	}, nil
}

func isIcuTemplateOptionName(key string) bool {
	switch key {
	case "description", "meaning", "alternate_message_id", "example", "original_code":
		return true
	}
	return false
}

type icuTemplateOptions struct {
	description                  string
	meaning                      string
	alternateMessageId           string
	placeholderExamplesMap       map[string]string
	placeholderExamplesOrder     []string
	placeholderOriginalCodeMap   map[string]string
	placeholderOriginalCodeOrder []string
	placeholderNames             map[string]struct{}
	exampleObjectLiteralMap      jsMessageObjectLiteralMap
	originalCodeObjectLiteralMap jsMessageObjectLiteralMap
	hasMeaning                   bool
	hasAlternateMessageId        bool
}

func (o *icuTemplateOptions) checkForUnknownPlaceholders(knownPlaceholders map[string]struct{}) *malformedMessageError {
	if err := o.exampleObjectLiteralMap.checkForUnexpectedKeys(knownPlaceholders, "Unknown placeholder: "); err != nil {
		return err
	}
	return o.originalCodeObjectLiteralMap.checkForUnexpectedKeys(knownPlaceholders, "Unknown placeholder: ")
}

// Java: JsMessageVisitor.extractIcuTemplateOptions(Node)
func extractIcuTemplateOptions(optionsBag rhino.Node) (icuTemplateOptions, *malformedMessageError) {
	objectLiteralMap, err := extractObjectLiteralMap(optionsBag)
	if err != nil {
		return icuTemplateOptions{}, err
	}
	if err := objectLiteralMap.checkForUnexpectedKeysFn(isIcuTemplateOptionName, "Unknown option: "); err != nil {
		return icuTemplateOptions{}, err
	}

	descriptionNode := objectLiteralMap.getValueNode("description")
	if descriptionNode.IsNil() {
		return icuTemplateOptions{}, &malformedMessageError{msg: "'description' option field is missing", node: optionsBag}
	}
	description, err := jsMessageVisitor_extractStringFromStringExprNode(descriptionNode)
	if err != nil {
		return icuTemplateOptions{}, err
	}

	var meaning string
	var hasMeaning bool
	if meaningNode := objectLiteralMap.getValueNode("meaning"); !meaningNode.IsNil() {
		m, err := jsMessageVisitor_extractStringFromStringExprNode(meaningNode)
		if err != nil {
			return icuTemplateOptions{}, err
		}
		meaning = m
		hasMeaning = true
	}

	var alternateMessageId string
	var hasAlternateMessageId bool
	if alternateMessageIdNode := objectLiteralMap.getValueNode("alternate_message_id"); !alternateMessageIdNode.IsNil() {
		alt, err := jsMessageVisitor_extractStringFromStringExprNode(alternateMessageIdNode)
		if err != nil {
			return icuTemplateOptions{}, err
		}
		alternateMessageId = alt
		hasAlternateMessageId = true
	}

	exampleObjectLiteralMap, err := extractObjectLiteralMap(objectLiteralMap.getValueNode("example"))
	if err != nil {
		return icuTemplateOptions{}, err
	}
	placeholderExamplesMap, placeholderExamplesOrder, err := exampleObjectLiteralMap.extractAsStringToStringMap()
	if err != nil {
		return icuTemplateOptions{}, err
	}

	originalCodeObjectLiteralMap, err := extractObjectLiteralMap(objectLiteralMap.getValueNode("original_code"))
	if err != nil {
		return icuTemplateOptions{}, err
	}
	placeholderOriginalCodeMap, placeholderOriginalCodeOrder, err := originalCodeObjectLiteralMap.extractAsStringToStringMap()
	if err != nil {
		return icuTemplateOptions{}, err
	}

	var placeholderNames map[string]struct{}
	if len(placeholderExamplesOrder) > 0 || len(placeholderOriginalCodeOrder) > 0 {
		placeholderNames = make(map[string]struct{}, len(placeholderExamplesOrder)+len(placeholderOriginalCodeOrder))
		for _, k := range placeholderExamplesOrder {
			if !JsMessage_IsCanonicalPlaceholderNameFormat(k) {
				return icuTemplateOptions{}, &malformedMessageError{
					msg:  "Placeholder not in UPPER_SNAKE_CASE: " + k,
					node: optionsBag,
				}
			}
			placeholderNames[k] = struct{}{}
		}
		for _, k := range placeholderOriginalCodeOrder {
			if _, seen := placeholderNames[k]; !seen && !JsMessage_IsCanonicalPlaceholderNameFormat(k) {
				return icuTemplateOptions{}, &malformedMessageError{
					msg:  "Placeholder not in UPPER_SNAKE_CASE: " + k,
					node: optionsBag,
				}
			}
			placeholderNames[k] = struct{}{}
		}
	}

	return icuTemplateOptions{
		description:                  description,
		meaning:                      meaning,
		hasMeaning:                   hasMeaning,
		alternateMessageId:           alternateMessageId,
		hasAlternateMessageId:        hasAlternateMessageId,
		placeholderExamplesMap:       placeholderExamplesMap,
		placeholderExamplesOrder:     placeholderExamplesOrder,
		placeholderOriginalCodeMap:   placeholderOriginalCodeMap,
		placeholderOriginalCodeOrder: placeholderOriginalCodeOrder,
		placeholderNames:             placeholderNames,
		exampleObjectLiteralMap:      exampleObjectLiteralMap,
		originalCodeObjectLiteralMap: originalCodeObjectLiteralMap,
	}, nil
}

type jsMessageObjEntry struct {
	key  string
	node rhino.Node
}

// jsMessageObjectLiteralMap represents an AST OBJECTLIT whose members are all
// unique STRING_KEY nodes. Small object literals (<=8 keys) avoid allocating a
// hash map altogether; larger ones lazily populate byKey.
type jsMessageObjectLiteralMap struct {
	entries []jsMessageObjEntry
	byKey   map[string]rhino.Node
}

func (m *jsMessageObjectLiteralMap) findNode(key string) (rhino.Node, bool) {
	if m.byKey != nil {
		n, ok := m.byKey[key]
		return n, ok
	}
	for _, e := range m.entries {
		if e.key == key {
			return e.node, true
		}
	}
	return rhino.NilNode, false
}

func extractObjectLiteralMap(objLit rhino.Node) (jsMessageObjectLiteralMap, *malformedMessageError) {
	if objLit.IsNil() {
		return jsMessageObjectLiteralMap{}, nil
	}
	lit, ok := ast.AsObjectLit(objLit)
	if !ok {
		return jsMessageObjectLiteralMap{}, &malformedMessageError{msg: "object literal expected", node: objLit}
	}
	var m jsMessageObjectLiteralMap
	for member := range lit.Members() {
		sk, ok := ast.AsStringKey(member.Node())
		if !ok {
			return jsMessageObjectLiteralMap{}, &malformedMessageError{msg: "string key expected", node: member.Node()}
		}
		key := sk.Key()
		if _, dup := m.findNode(key); dup {
			return jsMessageObjectLiteralMap{}, &malformedMessageError{msg: "duplicate string key: " + key, node: member.Node()}
		}
		if len(m.entries) == 8 && m.byKey == nil {
			m.byKey = make(map[string]rhino.Node, 16)
			for _, prev := range m.entries {
				m.byKey[prev.key] = prev.node
			}
		}
		m.entries = append(m.entries, jsMessageObjEntry{key: key, node: member.Node()})
		if m.byKey != nil {
			m.byKey[key] = member.Node()
		}
	}
	return m, nil
}

func (m *jsMessageObjectLiteralMap) getValueNode(key string) rhino.Node {
	stringKeyNode, ok := m.findNode(key)
	if !ok {
		return rhino.NilNode
	}
	return stringKeyNode.GetOnlyChild()
}

func (m *jsMessageObjectLiteralMap) getBooleanValueOrFalse(key string) (bool, *malformedMessageError) {
	stringKeyNode, ok := m.findNode(key)
	if !ok {
		return false, nil
	}
	valueNode := stringKeyNode.GetOnlyChild()
	if valueNode.IsTrue() {
		return true, nil
	}
	if valueNode.IsFalse() {
		return false, nil
	}
	return false, &malformedMessageError{
		msg:  stringKeyNode.GetString() + ": Literal true or false expected",
		node: valueNode,
	}
}

func (m *jsMessageObjectLiteralMap) extractAsValueMap() *jutil.LinkedHashMap[string, rhino.Node] {
	if len(m.entries) == 0 {
		return nil
	}
	out := jutil.NewLinkedHashMapWithCapacity[string, rhino.Node](len(m.entries))
	for _, e := range m.entries {
		out.PutNew(e.key, e.node.GetOnlyChild())
	}
	return out
}

func (m *jsMessageObjectLiteralMap) extractAsValueMapAndOrder() (map[string]rhino.Node, []string) {
	if len(m.entries) == 0 {
		return nil, nil
	}
	out := make(map[string]rhino.Node, len(m.entries))
	order := make([]string, len(m.entries))
	for i, e := range m.entries {
		out[e.key] = e.node.GetOnlyChild()
		order[i] = e.key
	}
	return out, order
}

func (m *jsMessageObjectLiteralMap) extractAsStringToStringMap() (map[string]string, []string, *malformedMessageError) {
	if len(m.entries) == 0 {
		return nil, nil, nil
	}
	out := make(map[string]string, len(m.entries))
	order := make([]string, len(m.entries))
	for i, e := range m.entries {
		val, err := jsMessageVisitor_extractStringFromStringExprNode(e.node.GetOnlyChild())
		if err != nil {
			return nil, nil, err
		}
		out[e.key] = val
		order[i] = e.key
	}
	return out, order, nil
}

func (m *jsMessageObjectLiteralMap) checkForUnexpectedKeys(
	expectedKeys map[string]struct{}, errorPrefix string,
) *malformedMessageError {
	for _, e := range m.entries {
		if _, ok := expectedKeys[e.key]; !ok {
			return &malformedMessageError{msg: errorPrefix + e.key, node: e.node}
		}
	}
	return nil
}

func (m *jsMessageObjectLiteralMap) checkForUnexpectedKeysFn(
	isExpected func(string) bool, errorPrefix string,
) *malformedMessageError {
	for _, e := range m.entries {
		if !isExpected(e.key) {
			return &malformedMessageError{msg: errorPrefix + e.key, node: e.node}
		}
	}
	return nil
}

func (m *jsMessageObjectLiteralMap) checkForRequiredKeys(
	requiredKeys []string, errNode rhino.Node, errorPrefix string,
) *malformedMessageError {
	for _, requiredKey := range requiredKeys {
		if _, ok := m.findNode(requiredKey); !ok {
			return &malformedMessageError{msg: errorPrefix + requiredKey, node: errNode}
		}
	}
	return nil
}

type googGetMsgParsedText struct {
	originalMsgText  string
	placeholderOrder []string
	placeholderSet   map[string]struct{}
	parts            []JsMessage_Part
}

func jsMessageVisitor_parseJsMessageTextIntoParts(originalMsgText string) ([]JsMessage_Part, error) {
	parsed, err := jsMessageVisitor_extractGoogGetMsgParsedText(originalMsgText)
	if err != nil {
		return nil, err
	}
	return parsed.parts, nil
}

var errPlaceholderMalformed = errors.New("Placeholder incorrectly formatted")

// Java: JsMessageVisitor.extractGoogGetMsgParsedText(String)
func jsMessageVisitor_extractGoogGetMsgParsedText(originalMsgText string) (googGetMsgParsedText, error) {
	msgText := originalMsgText
	var partsBuilder []JsMessage_Part
	var placeholderOrder []string
	var placeholderSet map[string]struct{}
	for {
		phBegin := strings.Index(msgText, JsMessage_PH_JS_PREFIX)
		if phBegin < 0 {
			partsBuilder = append(partsBuilder, JsMessage_StringPart_Create(msgText))
			break
		}
		if phBegin > 0 {
			partsBuilder = append(partsBuilder, JsMessage_StringPart_Create(msgText[:phBegin]))
		}

		afterPrefix := msgText[phBegin+len(JsMessage_PH_JS_PREFIX):]
		phEndRel := strings.Index(afterPrefix, JsMessage_PH_JS_SUFFIX)
		if phEndRel < 0 {
			return googGetMsgParsedText{}, errPlaceholderMalformed
		}

		phName := afterPrefix[:phEndRel]
		if !JsMessage_IsLowerCamelCaseWithNumericSuffixes(phName) {
			return googGetMsgParsedText{}, errors.New("Placeholder name not in lowerCamelCase: " + phName)
		}
		if placeholderSet == nil {
			placeholderSet = make(map[string]struct{})
		}
		if _, seen := placeholderSet[phName]; !seen {
			placeholderSet[phName] = struct{}{}
			placeholderOrder = append(placeholderOrder, phName)
		}
		partsBuilder = append(partsBuilder, JsMessage_PlaceholderReference_CreateForJsName(phName))
		rest := afterPrefix[phEndRel+len(JsMessage_PH_JS_SUFFIX):]
		if len(rest) == 0 {
			break
		}
		msgText = rest
	}
	return googGetMsgParsedText{
		originalMsgText:  originalMsgText,
		placeholderOrder: placeholderOrder,
		placeholderSet:   placeholderSet,
		parts:            partsBuilder,
	}, nil
}
