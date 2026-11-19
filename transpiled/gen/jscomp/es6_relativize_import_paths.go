/*
 * Copyright 2018 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/Es6RelativizeImportPaths.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.Es6RelativizeImportPaths
type Es6RelativizeImportPaths struct {
	compiler AbstractCompiler `gson:"compiler"`
	idHash_  atomic.Int32
}

func (this_ *Es6RelativizeImportPaths) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RelativizeImportPaths) IsJavaNil_() bool { return this_ == nil }
func (*Es6RelativizeImportPaths) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RelativizeImportPaths", "Es6RelativizeImportPaths"
}
func (this_ *Es6RelativizeImportPaths) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6RelativizeImportPaths) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RelativizeImportPaths", "com.google.javascript.jscomp.Es6RelativizeImportPaths", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RelativizeImportPaths](v_)
	})
}
func (*Es6RelativizeImportPaths) Implements_jscomp_CompilerPass() {}

func (this_ *Es6RelativizeImportPaths) initFields_Es6RelativizeImportPaths() {
	_ = this_
}

// Java constructor: Es6RelativizeImportPaths(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewEs6RelativizeImportPaths(compiler AbstractCompiler) *Es6RelativizeImportPaths {
	this_ := &Es6RelativizeImportPaths{}
	InitEs6RelativizeImportPaths(this_, compiler)
	return this_
}

func InitEs6RelativizeImportPaths(this_ *Es6RelativizeImportPaths, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_Es6RelativizeImportPaths()
	this_.compiler = compiler
}

// Java method: Es6RelativizeImportPaths.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *Es6RelativizeImportPaths) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	{
		var script pkg_rhino.Node = root.GetFirstChild()
		_ = script
		for ; script != 0; script = script.GetNext() {
			if Es6RewriteModules_IsEs6ModuleRoot(script) {
				NodeTraversal_Traverse(this_.compiler, script, NewEs6RelativizeImportPaths_Rewriter())
			}
		}
	}
}

// Java class: com.google.javascript.jscomp.Es6RelativizeImportPaths.Rewriter
type Es6RelativizeImportPaths_Rewriter struct {
	NodeTraversal_AbstractPreOrderCallbackBase
}

func (this_ *Es6RelativizeImportPaths_Rewriter) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RelativizeImportPaths_Rewriter) IsJavaNil_() bool { return this_ == nil }
func (*Es6RelativizeImportPaths_Rewriter) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RelativizeImportPaths$Rewriter", "Rewriter"
}
func (this_ *Es6RelativizeImportPaths_Rewriter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPreOrderCallbackBase.IdentityHash_()
}
func (*Es6RelativizeImportPaths_Rewriter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RelativizeImportPaths.Rewriter", "com.google.javascript.jscomp.Es6RelativizeImportPaths$Rewriter", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RelativizeImportPaths_Rewriter](v_)
	})
}
func (*Es6RelativizeImportPaths_Rewriter) Implements_jscomp_NodeTraversal_AbstractPreOrderCallback() {
}
func (*Es6RelativizeImportPaths_Rewriter) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *Es6RelativizeImportPaths_Rewriter) initFields_Es6RelativizeImportPaths_Rewriter() {
	_ = this_
}

// Java constructor: Rewriter()
func NewEs6RelativizeImportPaths_Rewriter() *Es6RelativizeImportPaths_Rewriter {
	this_ := &Es6RelativizeImportPaths_Rewriter{}
	InitEs6RelativizeImportPaths_Rewriter(this_)
	return this_
}

func InitEs6RelativizeImportPaths_Rewriter(this_ *Es6RelativizeImportPaths_Rewriter) {
	InitNodeTraversal_AbstractPreOrderCallback(&this_.NodeTraversal_AbstractPreOrderCallbackBase)
	this_.initFields_Es6RelativizeImportPaths_Rewriter()
}

// Java method: Rewriter.shouldTraverse(com.google.javascript.jscomp.NodeTraversal nodeTraversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RelativizeImportPaths_Rewriter) ShouldTraverse_3(nodeTraversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = nodeTraversal
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(89) /* Token.ROOT */, int32(123) /* Token.MODULE_BODY */, int32(94) /* Token.SCRIPT */ :
		return true
	case int32(116) /* Token.IMPORT */ :
		this_.priv_Rewriter_visitImport_2561(nodeTraversal, n)
		return false
	case int32(120) /* Token.EXPORT */ :
		this_.priv_Rewriter_visitExport_2562(nodeTraversal, n)
		return false
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: Rewriter.visitImport(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node importDecl)
func (this_ *Es6RelativizeImportPaths_Rewriter) priv_Rewriter_visitImport_2561(t *NodeTraversal, importDecl pkg_rhino.Node) {
	_ = t
	_ = importDecl
	var specifierNode pkg_rhino.Node = importDecl.GetLastChild()
	_ = specifierNode
	var specifier string = specifierNode.GetString()
	_ = specifier
	if ModuleLoader_IsRelativeIdentifier(specifier) {
		return
	} else {
		var didRet_1 bool
		func() {
			defer func() {
				if rec_2 := recover(); rec_2 != nil {
					jex_3 := javalang.RecoverException(rec_2)
					_ = jex_3
					if j2gort.InstanceOf[j2gort.Ext_java_net_URISyntaxException](jex_3) {
						e := j2gort.AsIface[j2gort.Ext_java_net_URISyntaxException](jex_3)
						_ = e
						didRet_1 = true
						return
					} else {
						panic(rec_2)
					}
				}
			}()
			var specifierURI *j2gort.Ext_java_net_URI = j2gort.NewExt_java_net_URI_1(specifier)
			_ = specifierURI
			if specifierURI.IsAbsolute() {
				didRet_1 = true
				return
			}
		}()
		if didRet_1 {
			return
		}
	}
	var scriptPath string = t.GetInput().GetPath().String()
	_ = scriptPath
	var didRet_4 bool
	func() {
		defer func() {
			if rec_5 := recover(); rec_5 != nil {
				jex_6 := javalang.RecoverException(rec_5)
				_ = jex_6
				if j2gort.InstanceOf[j2gort.Ext_java_net_URISyntaxException](jex_6) {
					e_2 := j2gort.AsIface[j2gort.Ext_java_net_URISyntaxException](jex_6)
					_ = e_2
					didRet_4 = true
					return
				} else {
					panic(rec_5)
				}
			}
		}()
		scriptPath = j2gort.NewExt_java_net_URI_1(scriptPath).GetPath()
	}()
	if didRet_4 {
		return
	}
	var newSpecifier string = t.GetInput().GetPath().ResolveModuleAsPath(specifier).String()
	_ = newSpecifier
	if !(javalang.StringStartsWith(j2gort.CheckNotNullString(newSpecifier), "/")) {
		newSpecifier = ("/" + j2gort.StrValString(newSpecifier))
	}
	newSpecifier = j2gort.Ext_java_nio_file_Path_Of_2(scriptPath, []string{}).GetParent().Relativize(j2gort.Ext_java_nio_file_Path_Of_2(newSpecifier, []string{})).String()
	if ModuleLoader_IsAmbiguousIdentifier(newSpecifier) {
		newSpecifier = ("./" + j2gort.StrValString(newSpecifier))
	}
	if !(j2gort.CheckNotNullString(newSpecifier) == specifier) {
		specifierNode.SetString(newSpecifier)
		t.ReportCodeChange_1(specifierNode)
	}
}

// Java method: Rewriter.visitExport(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node export)
func (this_ *Es6RelativizeImportPaths_Rewriter) priv_Rewriter_visitExport_2562(t *NodeTraversal, export pkg_rhino.Node) {
	_ = t
	_ = export
	if export.HasTwoChildren() {
		this_.priv_Rewriter_visitImport_2561(t, export)
	}
}
func (*Es6RelativizeImportPaths_Rewriter) NoOpVisit_()                 {}
func (*Es6RelativizeImportPaths_Rewriter) CanReleaseSyntacticScopes_() {}
