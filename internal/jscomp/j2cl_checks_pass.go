// Port of com/google/javascript/jscomp/J2clChecksPass.java @ closure-compiler v20260415

/*
 * Copyright 2016 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// J2clChecksPass performs correctness checks which are specific to
// J2CL-generated patterns.
type J2clChecksPass struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler AbstractCompiler
}

// Java: J2clChecksPass.J2CL_REFERENCE_EQUALITY
var J2clChecksPass_J2CL_REFERENCE_EQUALITY = DiagnosticType_Warning(
	"JSC_J2CL_REFERENCE_EQUALITY",
	"Reference equality may not be used with the specified type: {0}")

// J2clChecksPass_REFERENCE_EQUALITY_TYPE_PATTERNS: types for which using
// reference equality is an error. Mapped from name to filename.
//
// Java: J2clChecksPass.REFERENCE_EQUALITY_TYPE_PATTERNS
var J2clChecksPass_REFERENCE_EQUALITY_TYPE_PATTERNS = func() *jutil.LinkedHashMap[string, string] {
	m := jutil.NewLinkedHashMap[string, string]()
	m.Put("java.lang.Integer", "java/lang/Integer.impl.java.js")
	m.Put("java.lang.Float", "java/lang/Float.impl.java.js")
	m.Put("goog.math.Long", "javascript/closure/math/long.js")
	return m
}()

// Java: J2clChecksPass(AbstractCompiler)
func NewJ2clChecksPass(compiler AbstractCompiler) *J2clChecksPass {
	p := &J2clChecksPass{}
	p.compiler = compiler
	return p
}

// Java: J2clChecksPass.visit(NodeTraversal, Node, Node)
func (p *J2clChecksPass) Visit(unused *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	for typeName := range J2clChecksPass_REFERENCE_EQUALITY_TYPE_PATTERNS.Keys() {
		p.CheckReferenceEquality(n, typeName, J2clChecksPass_REFERENCE_EQUALITY_TYPE_PATTERNS.Get(typeName))
	}
}

// CheckReferenceEquality reports an error if the node is a reference
// equality check of the specified type.
//
// Java: J2clChecksPass.checkReferenceEquality(Node, String, String)
func (p *J2clChecksPass) CheckReferenceEquality(n *rhino.Node, typeName string, fileName string) {
	if n.IsSHEQ() || n.IsEQ() || n.IsSHNE() || n.IsNE() {
		firstJsType := n.GetFirstChild().GetJSType()
		lastJsType := n.GetLastChild().GetJSType()
		hasType := p.IsType(firstJsType, fileName) || p.IsType(lastJsType, fileName)
		hasNullType := p.IsNullType(firstJsType) || p.IsNullType(lastJsType)
		if hasType && !hasNullType {
			p.compiler.Report(JSError_MakeWithNode(n, J2clChecksPass_J2CL_REFERENCE_EQUALITY, typeName))
		}
	}
}

// Java: J2clChecksPass.isNullType(JSType)
func (p *J2clChecksPass) IsNullType(jsType rhino.JSType) bool {
	if jsType == nil {
		return false
	}
	return jsType.IsNullType() || jsType.IsVoidType()
}

// Java: J2clChecksPass.isType(JSType, String)
func (p *J2clChecksPass) IsType(jsType rhino.JSType, fileName string) bool {
	if jsType == nil {
		return false
	}
	jsType = jsType.RestrictByNotNullOrUndefined()
	if jsType.ToMaybeObjectType() == nil {
		return false
	}
	sourceName := p.GetSourceName(jsType)
	return sourceName != nil && strings.HasSuffix(*sourceName, fileName)
}

// GetSourceName returns a nullable Java String.
//
// Java: J2clChecksPass.getSourceName(JSType)
func (p *J2clChecksPass) GetSourceName(jsType rhino.JSType) *string {
	constructor := jsType.ToMaybeObjectType().GetConstructor()
	if constructor == nil {
		empty := ""
		return &empty
	}
	return NodeUtil_GetSourceName(constructor.GetSource())
}

// Java: J2clChecksPass.process(Node, Node)
func (p *J2clChecksPass) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(p.compiler, root, p)
}
