// Port of com/google/javascript/jscomp/J2clChecksPass.java @ closure-compiler v20261028

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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// J2clChecksPass performs correctness checks which are specific to
// J2CL-generated patterns.
type J2clChecksPass struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler AbstractCompiler
}

var j2clChecksPassVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_SHEQ)
	s.Add(rhino.Token_EQ)
	s.Add(rhino.Token_SHNE)
	s.Add(rhino.Token_NE)
	return s
}()

func (*J2clChecksPass) noVisitLeaves() {}
func (*J2clChecksPass) visitTokens() *ast.TokenSet {
	return &j2clChecksPassVisitTokens
}

// Java: J2clChecksPass.J2CL_REFERENCE_EQUALITY
var J2clChecksPass_J2CL_REFERENCE_EQUALITY = DiagnosticType_Warning(
	"JSC_J2CL_REFERENCE_EQUALITY",
	"Reference equality may not be used with the specified type: {0}")

var j2clReferenceEqualityPatterns = [...]struct {
	typeName string
	fileName string
}{
	{"java.lang.Integer", "java/lang/Integer.impl.java.js"},
	{"java.lang.Float", "java/lang/Float.impl.java.js"},
	{"goog.math.Long", "javascript/closure/math/long.js"},
}

// Java: J2clChecksPass(AbstractCompiler)
func NewJ2clChecksPass(compiler AbstractCompiler) *J2clChecksPass {
	return &J2clChecksPass{compiler: compiler}
}

// Java: J2clChecksPass.visit(NodeTraversal, Node, Node)
func (p *J2clChecksPass) Visit(_ *NodeTraversal, n rhino.Node, _ rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_SHEQ, rhino.Token_EQ, rhino.Token_SHNE, rhino.Token_NE:
		bin := ast.MustBinary(n)
		firstJsType := bin.Left().Node().GetJSType()
		lastJsType := bin.Right().Node().GetJSType()
		if (firstJsType == nil && lastJsType == nil) || isNullOrVoidType(firstJsType) || isNullOrVoidType(lastJsType) {
			return
		}
		firstSrc := j2clCtorSourceName(firstJsType)
		lastSrc := j2clCtorSourceName(lastJsType)
		if firstSrc == "" && lastSrc == "" {
			return
		}
		for i := range j2clReferenceEqualityPatterns {
			pat := &j2clReferenceEqualityPatterns[i]
			if (firstSrc != "" && strings.HasSuffix(firstSrc, pat.fileName)) ||
				(lastSrc != "" && strings.HasSuffix(lastSrc, pat.fileName)) {
				p.compiler.Report(JSError_MakeWithNode(n, J2clChecksPass_J2CL_REFERENCE_EQUALITY, pat.typeName))
			}
		}
	}
}

func isNullOrVoidType(jsType rhino.JSType) bool {
	return jsType != nil && (jsType.IsNullType() || jsType.IsVoidType())
}

func j2clCtorSourceName(jsType rhino.JSType) string {
	if jsType == nil {
		return ""
	}
	objType := jsType.RestrictByNotNullOrUndefined().ToMaybeObjectType()
	if objType == nil {
		return ""
	}
	ctor := objType.GetConstructor()
	if ctor == nil {
		return ""
	}
	src := ctor.GetSource()
	if src.IsNil() {
		return ""
	}
	if sf := src.GetStaticSourceFile(); sf != nil {
		return sf.GetName()
	}
	sourceName, ok := NodeUtil_SourceName(src)
	if !ok {
		return ""
	}
	return sourceName
}

// Java: J2clChecksPass.process(Node, Node)
func (p *J2clChecksPass) Process(externs rhino.Node, root rhino.Node) {
	if root.IsNil() {
		return
	}
	if comp := p.compiler; comp != nil && !comp.IsDebugLoggingEnabled() {
		if traverseScriptsParallel(comp, root, parallelScriptOptions[*J2clChecksPass]{
			enableTypeRegistryReads: true,
			initWorker: func(_ int, rc *Compiler) *J2clChecksPass {
				return &J2clChecksPass{compiler: rc}
			},
		}) {
			return
		}
	}
	NodeTraversal_Traverse(p.compiler, root, p)
}
