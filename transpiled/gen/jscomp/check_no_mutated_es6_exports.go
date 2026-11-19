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
// Port of com/google/javascript/jscomp/lint/CheckNoMutatedEs6Exports.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.lint.CheckNoMutatedEs6Exports
type CheckNoMutatedEs6Exports struct {
	compiler           AbstractCompiler                              `gson:"compiler"`
	mutatedNames       j2gort.Ext_com_google_common_collect_Multimap `gson:"mutatedNames"`
	exportedLocalNames j2gort.Ext_java_util_Set                      `gson:"exportedLocalNames"`
	idHash_            atomic.Int32
}

func (this_ *CheckNoMutatedEs6Exports) BindSelf_(s any) {
	_ = s
}
func (this_ *CheckNoMutatedEs6Exports) IsJavaNil_() bool { return this_ == nil }
func (*CheckNoMutatedEs6Exports) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.lint.CheckNoMutatedEs6Exports", "CheckNoMutatedEs6Exports"
}
func (this_ *CheckNoMutatedEs6Exports) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CheckNoMutatedEs6Exports) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.lint.CheckNoMutatedEs6Exports", "com.google.javascript.jscomp.lint.CheckNoMutatedEs6Exports", func(v_ any) bool {
		return j2gort.InstanceOf[*CheckNoMutatedEs6Exports](v_)
	})
}
func (*CheckNoMutatedEs6Exports) Implements_jscomp_NodeTraversal_Callback() {}
func (*CheckNoMutatedEs6Exports) Implements_jscomp_CompilerPass()           {}

var static_CheckNoMutatedEs6Exports_MUTATED_EXPORT *DiagnosticType

func StaticGet_CheckNoMutatedEs6Exports_MUTATED_EXPORT() *DiagnosticType {
	clinit_CheckNoMutatedEs6Exports()
	return static_CheckNoMutatedEs6Exports_MUTATED_EXPORT
}
func StaticSet_CheckNoMutatedEs6Exports_MUTATED_EXPORT(v *DiagnosticType) {
	clinit_CheckNoMutatedEs6Exports()
	static_CheckNoMutatedEs6Exports_MUTATED_EXPORT = v
}

var clinitDone_CheckNoMutatedEs6Exports uint32
var clinitState_CheckNoMutatedEs6Exports j2gort.ClinitState

func clinit_CheckNoMutatedEs6Exports() {
	if atomic.LoadUint32(&clinitDone_CheckNoMutatedEs6Exports) != 0 {
		return
	}
	clinitSlow_CheckNoMutatedEs6Exports()
}

func clinitSlow_CheckNoMutatedEs6Exports() {
	if !j2gort.EnterClinit(&clinitState_CheckNoMutatedEs6Exports, &clinitDone_CheckNoMutatedEs6Exports) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_CheckNoMutatedEs6Exports, &clinitDone_CheckNoMutatedEs6Exports)
	static_CheckNoMutatedEs6Exports_MUTATED_EXPORT = DiagnosticType_Warning("JSC_MUTATED_EXPORT", "The name \"{0}\" is exported and should not be mutated outside of module initialization. Mutable exports are generally difficult to reason about. You can work around this by exporting getter/setter functions, or an object with mutable properties instead.")
}

func (this_ *CheckNoMutatedEs6Exports) initFields_CheckNoMutatedEs6Exports() {
	_ = this_
	this_.mutatedNames = j2gort.AsIface[j2gort.Ext_com_google_common_collect_SetMultimap](j2gort.Ext_com_google_common_collect_MultimapBuilder_HashKeys_0().HashSetValues_0().Build_0())
	this_.exportedLocalNames = j2gort.NewExt_java_util_LinkedHashSet()
}

// Java constructor: CheckNoMutatedEs6Exports(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewCheckNoMutatedEs6Exports(compiler AbstractCompiler) *CheckNoMutatedEs6Exports {
	clinit_CheckNoMutatedEs6Exports()
	this_ := &CheckNoMutatedEs6Exports{}
	InitCheckNoMutatedEs6Exports(this_, compiler)
	return this_
}

func InitCheckNoMutatedEs6Exports(this_ *CheckNoMutatedEs6Exports, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_CheckNoMutatedEs6Exports()
	this_.compiler = compiler
}

// Java method: CheckNoMutatedEs6Exports.checkNoMutations()
func (this_ *CheckNoMutatedEs6Exports) priv_CheckNoMutatedEs6Exports_checkNoMutations_6588() {
	var mutatedExports j2gort.Ext_java_util_Set = j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_Sets_Intersection(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.mutatedNames.KeySet()), this_.exportedLocalNames))
	_ = mutatedExports
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(mutatedExports))
		for cur_1.HasNext() {
			mutatedExport := cur_1.NextStr()
			_ = mutatedExport
			{
				cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.MultimapGetStr(this_.mutatedNames, mutatedExport))))
				for cur_2.HasNext() {
					mutation := pkg_rhino.Node(cur_2.NextNode())
					_ = mutation
					this_.compiler.Report_1(JSError_Make_3(mutation, StaticGet_CheckNoMutatedEs6Exports_MUTATED_EXPORT(), []string{mutatedExport}))
				}
			}
		}
	}
	this_.mutatedNames.Clear()
	this_.exportedLocalNames.Clear()
}

// Java method: CheckNoMutatedEs6Exports.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckNoMutatedEs6Exports) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsModuleBody() {
		this_.priv_CheckNoMutatedEs6Exports_checkNoMutations_6588()
	}
}

// Java method: CheckNoMutatedEs6Exports.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *CheckNoMutatedEs6Exports) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(94) /* Token.SCRIPT */ :
		return n.GetBooleanProp(pkg_rhino.StaticGet_Node_ES6_MODULE())
	case int32(120) /* Token.EXPORT */ :
		this_.priv_CheckNoMutatedEs6Exports_visitExport_6589(n)
		return true
	case int32(32) /* Token.NAME */ :
		this_.priv_CheckNoMutatedEs6Exports_visitName_6590(t, n)
		return true
	default:
		return true
	}
	panic("unreachable switch expression")
}

// Java method: CheckNoMutatedEs6Exports.visitExport(com.google.javascript.rhino.Node export)
func (this_ *CheckNoMutatedEs6Exports) priv_CheckNoMutatedEs6Exports_visitExport_6589(export pkg_rhino.Node) {
	_ = export
	if export.HasOneChild() && export.GetFirstChild().IsExportSpecs() {
		{
			var exportSpec pkg_rhino.Node = export.GetFirstFirstChild()
			_ = exportSpec
			for ; exportSpec != 0; exportSpec = exportSpec.GetNext() {
				if !(exportSpec.HasTwoChildren()) {
					j2gort.CheckState(false)
				}
				j2gort.SetAddStr(this_.exportedLocalNames, exportSpec.GetFirstChild().GetString())
			}
		}
	} else if export.HasOneChild() && !(export.GetBooleanProp(pkg_rhino.StaticGet_Node_EXPORT_ALL_FROM())) {
		var declaration pkg_rhino.Node = export.GetFirstChild()
		_ = declaration
		if NodeUtil_IsNameDeclaration(declaration) {
			NodeUtil_VisitLhsNodesInNode(declaration, j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_3 any) {
				lhs := j2gort.CastNode[pkg_rhino.Node](lp_3)
				_ = lhs
				if !(lhs.IsName()) {
					j2gort.CheckState(false)
				}
				j2gort.SetAddStr(this_.exportedLocalNames, lhs.GetString())
			}))
		} else if export.GetBooleanProp(pkg_rhino.StaticGet_Node_EXPORT_DEFAULT()) {
			if !(declaration.IsClass()) && !(declaration.IsFunction()) {
				return
			}
			var nameNode pkg_rhino.Node = declaration.GetFirstChild()
			_ = nameNode
			if !(nameNode.IsEmpty()) && !(j2gort.CheckNotNullString(nameNode.GetString()) == "") {
				j2gort.SetAddStr(this_.exportedLocalNames, nameNode.GetString())
			}
		} else {
			if !(declaration.IsClass() || declaration.IsFunction()) {
				j2gort.CheckState(false)
			}
			var nameNode_2 pkg_rhino.Node = declaration.GetFirstChild()
			_ = nameNode_2
			j2gort.SetAddStr(this_.exportedLocalNames, nameNode_2.GetString())
		}
	}
}

// Java method: CheckNoMutatedEs6Exports.visitName(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node name)
func (this_ *CheckNoMutatedEs6Exports) priv_CheckNoMutatedEs6Exports_visitName_6590(t *NodeTraversal, name pkg_rhino.Node) {
	_ = t
	_ = name
	var scope *Scope = t.GetScope()
	_ = scope
	if NodeUtil_IsLValue(name) && !(j2gort.CastPtr[Scope](scope.GetClosestHoistScope()).IsModuleScope()) {
		var var_ *Var = j2gort.CastPtr[Var](scope.GetVar(name.GetString()))
		_ = var_
		if (var_ != nil) && j2gort.CastPtr[Scope](var_.GetScope()).IsModuleScope() {
			j2gort.MultimapPutStr(this_.mutatedNames, name.GetString(), j2gort.BoxNode(name))
		}
	}
}

// Java method: CheckNoMutatedEs6Exports.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CheckNoMutatedEs6Exports) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
}
func (*CheckNoMutatedEs6Exports) CanReleaseSyntacticScopes_() {}
