// Port of com/google/javascript/jscomp/GatherExternProperties.java @ closure-compiler v20261021

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

// GatherExternProperties gathers property names defined in externs.
//
// The collection of these property names could easily happen during type
// checking. However, when exporting local property definitions, the externs
// may be modified after type checking, and we want to collect the new names
// as well.
//
// Mode.CHECK is set when checksOnly=true and only collects the externs in
// the type expressions in JSDoc. Mode.OPTIMIZE is set when checksOnly=false
// and only collects the externs in AST properties.
//
// To optimize build time, we collect only type expression externs in
// getChecks() and only AST property externs in getOptimizations() instead of
// collecting all in both.
type GatherExternProperties struct {
	compiler AbstractCompiler

	externProperties *jutil.LinkedHashSet[string]

	mode GatherExternProperties_Mode
}

// GatherExternProperties_Mode is the Java enum GatherExternProperties.Mode.
type GatherExternProperties_Mode int32

const (
	GatherExternProperties_Mode_CHECK GatherExternProperties_Mode = iota
	GatherExternProperties_Mode_OPTIMIZE
	GatherExternProperties_Mode_CHECK_AND_OPTIMIZE
)

var gatherExternProperties_Mode_names = [...]string{"CHECK", "OPTIMIZE", "CHECK_AND_OPTIMIZE"}

// Java: Mode.check
var gatherExternProperties_Mode_check = [...]bool{true, false, true}

// Java: Mode.optimize
var gatherExternProperties_Mode_optimize = [...]bool{false, true, true}

// Name is Enum.name().
func (m GatherExternProperties_Mode) Name() string { return gatherExternProperties_Mode_names[m] }

// String is Enum.toString().
func (m GatherExternProperties_Mode) String() string { return m.Name() }

// Ordinal is Enum.ordinal().
func (m GatherExternProperties_Mode) Ordinal() int { return int(m) }

// Java: GatherExternProperties(AbstractCompiler, Mode)
func NewGatherExternProperties(compiler AbstractCompiler, mode GatherExternProperties_Mode) *GatherExternProperties {
	g := &GatherExternProperties{externProperties: jutil.NewLinkedHashSet[string]()}
	g.compiler = compiler
	g.mode = mode

	if compiler.GetExternProperties() != nil {
		for p := range compiler.GetExternProperties().All() {
			g.externProperties.Add(p)
		}
	}
	return g
}

// Java: GatherExternProperties.process(Node, Node)
func (g *GatherExternProperties) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(g.compiler, externs, g)
	g.compiler.SetExternProperties(jutil.NewLinkedHashSet(g.externProperties.ToSlice()...))
}

// Java: GatherExternProperties.shouldTraverse(NodeTraversal, Node, Node)
func (g *GatherExternProperties) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	return !n.IsScript() || !NodeUtil_IsFromTypeSummary(n)
}

// Java: GatherExternProperties.visit(NodeTraversal, Node, Node)
func (g *GatherExternProperties) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if gatherExternProperties_Mode_optimize[g.mode] {
		switch n.GetToken() {
		case rhino.Token_GETPROP:
			// Gathers "name" from (someObject.name).
			g.externProperties.Add(n.GetString())
		case rhino.Token_STRING_KEY:
			if parent.IsObjectLit() {
				g.externProperties.Add(n.GetString())
			}
		case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_MEMBER_FIELD_DEF, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
			g.externProperties.Add(n.GetString())
		default:
		}
	}

	jsDocInfo := n.GetJSDocInfo()
	if jsDocInfo != nil && gatherExternProperties_Mode_check[g.mode] {
		g.GatherPropertiesFromJSDocInfo(jsDocInfo)
	}
}

// Java: GatherExternProperties.gatherPropertiesFromJSDocInfo(JSDocInfo)
func (g *GatherExternProperties) GatherPropertiesFromJSDocInfo(jsDocInfo *rhino.JSDocInfo) {
	for it := jsDocInfo.TypeExprIter(); it.Next(); {
		g.GatherPropertiesFromJsTypeExpressionNode(it.Node())
	}
}

// Java: GatherExternProperties.gatherPropertiesFromJsTypeExpressionNode(Node)
func (g *GatherExternProperties) GatherPropertiesFromJsTypeExpressionNode(jsTypeExpressionNode rhino.Node) {
	switch jsTypeExpressionNode.GetToken() {
	case rhino.Token_LB:
		g.GatherPropertiesFromJsDocRecordType(jsTypeExpressionNode)
	default:
		for child := jsTypeExpressionNode.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			g.GatherPropertiesFromJsTypeExpressionNode(child)
		}
	}
}

// Java: GatherExternProperties.gatherPropertiesFromJsDocRecordType(Node)
func (g *GatherExternProperties) GatherPropertiesFromJsDocRecordType(jsDocRecordNode rhino.Node) {
	jutil.CheckStateMsg1(jsDocRecordNode.GetToken() == rhino.Token_LB, "%s", jsDocRecordNode)
	for fieldNode := jsDocRecordNode.GetFirstChild(); !fieldNode.IsNil(); fieldNode = fieldNode.GetNext() {
		var fieldNameNode rhino.Node
		var fieldTypeNode rhino.Node
		if fieldNode.GetToken() == rhino.Token_COLON {
			fieldNameNode = fieldNode.GetFirstChild()
			fieldTypeNode = fieldNameNode.GetNext()
		} else {
			fieldNameNode = fieldNode
			fieldTypeNode = rhino.NilNode
		}
		jutil.CheckStateMsg1(fieldNameNode.IsStringKey(), "%s", fieldNameNode)
		fieldName := fieldNameNode.GetString()
		// TODO(bradfordcsmith): The JSDoc parser should do this.
		if jutil.StringStartsWith(fieldName, "'") || jutil.StringStartsWith(fieldName, "\"") {
			fieldName = jutil.StringSubstring(fieldName, 1, jutil.StringLength(fieldName)-1)
		}
		g.externProperties.Add(fieldName)
		if !fieldTypeNode.IsNil() {
			g.GatherPropertiesFromJsTypeExpressionNode(fieldTypeNode)
		}
	}
}
