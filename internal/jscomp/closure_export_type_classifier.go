// Port of com/google/javascript/jscomp/ClosureExportTypeClassifier.java @ closure-compiler v20261104

/*
 * Copyright 2026 The Closure Compiler Authors.
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

const (
	closureExportTypeClassifier_EXPORTS         = "exports"
	closureExportTypeClassifier_MAX_ALIAS_DEPTH = 4
)

type closureExportTypeClassifier_Kind int32

const (
	closureExportTypeClassifier_Kind_NOT_NAMESPACE closureExportTypeClassifier_Kind = iota
	closureExportTypeClassifier_Kind_EMPTY_OBJECT
	closureExportTypeClassifier_Kind_NAMESPACE
)

type closureExportTypeClassifier_NameInfo struct {
	value                 rhino.Node
	jsdoc                 *rhino.JSDocInfo
	isNeverNamespace      bool
	propertyCount         int
	isGoogModuleExports   bool
	isNamedExportsLiteral bool
}

// ClosureExportTypeClassifier_TsExportType classifies whether a goog.module
// namespace is seen from TypeScript as a single default export or named exports.
type ClosureExportTypeClassifier_TsExportType int32

const (
	// ClosureExportTypeClassifier_TsExportType_DEFAULT indicates the module exports a single default value.
	ClosureExportTypeClassifier_TsExportType_DEFAULT ClosureExportTypeClassifier_TsExportType = iota
	// ClosureExportTypeClassifier_TsExportType_NAMED indicates the module exports multiple named values.
	ClosureExportTypeClassifier_TsExportType_NAMED
)

var closureExportTypeClassifier_TsExportType_names = [...]string{"DEFAULT", "NAMED"}

// Name is Enum.name().
func (e ClosureExportTypeClassifier_TsExportType) Name() string {
	return closureExportTypeClassifier_TsExportType_names[e]
}

// String is Enum.toString().
func (e ClosureExportTypeClassifier_TsExportType) String() string { return e.Name() }

// ClosureExportTypeClassifier syntactically classifies whether each goog.module
// namespace in a file is seen from TypeScript as a set of "named" exports or as
// a single "default" export.
type ClosureExportTypeClassifier struct {
	topLevelNames     map[string]*closureExportTypeClassifier_NameInfo
	forcedNamedExport map[string]struct{}
}

// ClosureExportTypeClassifier_Classify classifies a single goog.module namespace
// as named or default.
//
// Java: ClosureExportTypeClassifier.classify(Node, String)
func ClosureExportTypeClassifier_Classify(moduleRoot rhino.Node, googNamespace string) ClosureExportTypeClassifier_TsExportType {
	classifier := &ClosureExportTypeClassifier{
		topLevelNames:     make(map[string]*closureExportTypeClassifier_NameInfo),
		forcedNamedExport: make(map[string]struct{}),
	}
	classifier.addGoogModuleExports()
	closureExportTypeClassifier_forEachTopLevelStatement(moduleRoot, classifier.visitStatement)

	info := classifier.topLevelNames[closureExportTypeClassifier_EXPORTS]
	if _, forced := classifier.forcedNamedExport[googNamespace]; forced ||
		classifier.kindOfName(info, 0) == closureExportTypeClassifier_Kind_NAMESPACE {
		return ClosureExportTypeClassifier_TsExportType_NAMED
	}
	return ClosureExportTypeClassifier_TsExportType_DEFAULT
}

func (c *ClosureExportTypeClassifier) addGoogModuleExports() {
	c.topLevelNames[closureExportTypeClassifier_EXPORTS] = &closureExportTypeClassifier_NameInfo{
		isGoogModuleExports: true,
	}
}

func closureExportTypeClassifier_forEachTopLevelStatement(moduleRoot rhino.Node, callback func(rhino.Node)) {
	jutil.CheckArgumentMsg(
		moduleRoot.IsScript() && !moduleRoot.GetFirstChild().IsNil() && moduleRoot.GetFirstChild().IsModuleBody(),
		"%s", moduleRoot)
	body := moduleRoot.GetFirstChild()
	for statement := body.GetFirstChild(); !statement.IsNil(); statement = statement.GetNext() {
		callback(statement)
	}
}

func (c *ClosureExportTypeClassifier) visitStatement(statement rhino.Node) {
	switch statement.GetToken() {
	case rhino.Token_VAR, rhino.Token_LET, rhino.Token_CONST:
		c.visitNameDeclaration(statement)
	case rhino.Token_CLASS, rhino.Token_FUNCTION:
		name := statement.GetFirstChild()
		info := c.getOrCreate(name.GetString())
		info.isNeverNamespace = true
	case rhino.Token_EXPR_RESULT:
		c.maybeRecordLegacyNamedTsExport(statement)
		c.visitExpression(statement.GetFirstChild())
	default:
	}
}

func (c *ClosureExportTypeClassifier) maybeRecordLegacyNamedTsExport(statement rhino.Node) {
	if !NodeUtil_IsGoogModuleCall(statement) {
		return
	}
	call := statement.GetFirstChild()
	jsdoc := call.GetJSDocInfo()
	if jsdoc == nil || !jsdoc.IsLegacyNamedTsExport() {
		return
	}
	c.forcedNamedExport[call.GetSecondChild().GetString()] = struct{}{}
}

func (c *ClosureExportTypeClassifier) visitNameDeclaration(declaration rhino.Node) {
	for lhs := declaration.GetFirstChild(); !lhs.IsNil(); lhs = lhs.GetNext() {
		if lhs.IsName() {
			c.visitDeclaredName(lhs)
		}
	}
}

func (c *ClosureExportTypeClassifier) visitDeclaredName(lhs rhino.Node) {
	initializer := lhs.GetFirstChild()
	info := c.getOrCreate(lhs.GetString())
	if !initializer.IsNil() && initializer.IsAssign() {
		c.visitExpression(initializer)
	}
	jsdoc := NodeUtil_GetBestJSDocInfo(lhs)
	if !initializer.IsNil() {
		info.value = initializer
		info.jsdoc = jsdoc
	} else if info.value.IsNil() && jsdoc != nil {
		// A stub like `/** @type {T} */ let x;`. A bare redeclaration like `var x;` keeps any
		// earlier value.
		info.jsdoc = jsdoc
	}
}

func (c *ClosureExportTypeClassifier) visitExpression(expr rhino.Node) {
	switch expr.GetToken() {
	case rhino.Token_ASSIGN:
		c.visitAssign(expr)
	case rhino.Token_GETPROP:
		// Stub declarations like `/** @typedef {number} */ exports.T;`
		if expr.GetJSDocInfo() != nil && expr.GetFirstChild().IsName() {
			variableName := expr.GetFirstChild().GetString()
			c.getOrCreate(variableName).propertyCount++
		}
	default:
	}
}

func (c *ClosureExportTypeClassifier) visitAssign(expr rhino.Node) {
	lhs := expr.GetFirstChild()
	rhs := expr.GetSecondChild()
	jsdoc := expr.GetJSDocInfo()
	if lhs.IsName() {
		// `exports = ...;` or `localName = ...;`
		assigned := c.getOrCreate(lhs.GetString())
		assigned.value = rhs
		assigned.jsdoc = jsdoc
		assigned.isNeverNamespace = false
		assigned.isNamedExportsLiteral = assigned.isGoogModuleExports && NodeUtil_IsNamedExportsLiteral(rhs)
		return
	}
	if lhs.IsGetProp() && lhs.GetFirstChild().IsName() {
		// `exports.x = ...;` or `localName.x = ...;`
		c.getOrCreate(lhs.GetFirstChild().GetString()).propertyCount++
	}
}

func (c *ClosureExportTypeClassifier) getOrCreate(name string) *closureExportTypeClassifier_NameInfo {
	if info, ok := c.topLevelNames[name]; ok {
		return info
	}
	info := &closureExportTypeClassifier_NameInfo{}
	c.topLevelNames[name] = info
	return info
}

func (c *ClosureExportTypeClassifier) kindOfName(info *closureExportTypeClassifier_NameInfo, depth int) closureExportTypeClassifier_Kind {
	if info.isNeverNamespace {
		return closureExportTypeClassifier_Kind_NOT_NAMESPACE
	}
	if info.isNamedExportsLiteral {
		return closureExportTypeClassifier_Kind_NAMESPACE
	}
	var kind closureExportTypeClassifier_Kind
	if !info.value.IsNil() || info.jsdoc != nil {
		kind = c.kindOfValue(info.value, info.jsdoc, depth)
	} else if info.isGoogModuleExports {
		kind = closureExportTypeClassifier_Kind_EMPTY_OBJECT
	} else {
		kind = closureExportTypeClassifier_Kind_NOT_NAMESPACE
	}
	if kind == closureExportTypeClassifier_Kind_EMPTY_OBJECT && info.propertyCount > 0 {
		return closureExportTypeClassifier_Kind_NAMESPACE
	}
	return kind
}

func (c *ClosureExportTypeClassifier) kindOfValue(value rhino.Node, jsdoc *rhino.JSDocInfo, depth int) closureExportTypeClassifier_Kind {
	if jsdoc != nil {
		if jsdoc.HasEnumParameterType() || jsdoc.HasTypedefType() || jsdoc.IsConstructorOrInterface() {
			return closureExportTypeClassifier_Kind_NOT_NAMESPACE
		}
		if jsdoc.HasType() {
			return closureExportTypeClassifier_kindOfDeclaredType(jsdoc.GetType())
		}
	}
	if value.IsNil() {
		return closureExportTypeClassifier_Kind_NOT_NAMESPACE
	}
	switch value.GetToken() {
	case rhino.Token_OBJECTLIT:
		// Computed properties and spreads do not create named own properties.
		for member := value.GetFirstChild(); !member.IsNil(); member = member.GetNext() {
			if !member.IsComputedProp() && !member.IsSpread() {
				return closureExportTypeClassifier_Kind_NAMESPACE
			}
		}
		return closureExportTypeClassifier_Kind_EMPTY_OBJECT
	case rhino.Token_NAME:
		if depth >= closureExportTypeClassifier_MAX_ALIAS_DEPTH {
			return closureExportTypeClassifier_Kind_NOT_NAMESPACE
		}
		aliased := c.topLevelNames[value.GetString()]
		if aliased == nil {
			return closureExportTypeClassifier_Kind_NOT_NAMESPACE
		}
		return c.kindOfName(aliased, depth+1)
	case rhino.Token_CAST:
		return c.kindOfValue(value.GetFirstChild(), value.GetJSDocInfo(), depth)
	default:
		// Classes, functions, calls, `new`, literals, expressions, etc.
		return closureExportTypeClassifier_Kind_NOT_NAMESPACE
	}
}

func closureExportTypeClassifier_kindOfDeclaredType(type_ *rhino.JSTypeExpression) closureExportTypeClassifier_Kind {
	root := type_.GetRoot()
	if root.GetToken() == rhino.Token_LC {
		return closureExportTypeClassifier_Kind_NAMESPACE
	}
	return closureExportTypeClassifier_Kind_NOT_NAMESPACE
}
