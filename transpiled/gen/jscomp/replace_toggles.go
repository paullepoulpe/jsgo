/*
 * Copyright 2023 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ReplaceToggles.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.ReplaceToggles
type ReplaceToggles struct {
	compiler       AbstractCompiler                                  `gson:"compiler"`
	astFactory     *AstFactory                                       `gson:"astFactory"`
	ordinalMapping j2gort.Ext_com_google_common_collect_ImmutableMap `gson:"ordinalMapping"`
	idHash_        atomic.Int32
}

func (this_ *ReplaceToggles) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceToggles) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceToggles) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceToggles", "ReplaceToggles"
}
func (this_ *ReplaceToggles) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ReplaceToggles) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceToggles", "com.google.javascript.jscomp.ReplaceToggles", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceToggles](v_)
	})
}
func (*ReplaceToggles) Implements_jscomp_CompilerPass() {}

var static_ReplaceToggles_INVALID_TOGGLE_PARAMETER *DiagnosticType

func StaticGet_ReplaceToggles_INVALID_TOGGLE_PARAMETER() *DiagnosticType {
	clinit_ReplaceToggles()
	return static_ReplaceToggles_INVALID_TOGGLE_PARAMETER
}
func StaticSet_ReplaceToggles_INVALID_TOGGLE_PARAMETER(v *DiagnosticType) {
	clinit_ReplaceToggles()
	static_ReplaceToggles_INVALID_TOGGLE_PARAMETER = v
}

var static_ReplaceToggles_INVALID_ORDINAL_MAPPING *DiagnosticType

func StaticGet_ReplaceToggles_INVALID_ORDINAL_MAPPING() *DiagnosticType {
	clinit_ReplaceToggles()
	return static_ReplaceToggles_INVALID_ORDINAL_MAPPING
}
func StaticSet_ReplaceToggles_INVALID_ORDINAL_MAPPING(v *DiagnosticType) {
	clinit_ReplaceToggles()
	static_ReplaceToggles_INVALID_ORDINAL_MAPPING = v
}
func StaticGet_ReplaceToggles_TRUE_VALUE() int32        { return int32(-2) }
func StaticGet_ReplaceToggles_FALSE_VALUE() int32       { return int32(-1) }
func StaticGet_ReplaceToggles_ORDINAL_VAR_NAME() string { return "CLOSURE_TOGGLE_ORDINALS" }

var static_ReplaceToggles_readToggleFunctionName pkg_rhino.QualifiedName

func StaticGet_ReplaceToggles_ReadToggleFunctionName() pkg_rhino.QualifiedName {
	clinit_ReplaceToggles()
	return static_ReplaceToggles_readToggleFunctionName
}
func StaticSet_ReplaceToggles_ReadToggleFunctionName(v pkg_rhino.QualifiedName) {
	clinit_ReplaceToggles()
	static_ReplaceToggles_readToggleFunctionName = v
}

var clinitDone_ReplaceToggles uint32
var clinitState_ReplaceToggles j2gort.ClinitState

func clinit_ReplaceToggles() {
	if atomic.LoadUint32(&clinitDone_ReplaceToggles) != 0 {
		return
	}
	clinitSlow_ReplaceToggles()
}

func clinitSlow_ReplaceToggles() {
	if !j2gort.EnterClinit(&clinitState_ReplaceToggles, &clinitDone_ReplaceToggles) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ReplaceToggles, &clinitDone_ReplaceToggles)
	static_ReplaceToggles_INVALID_TOGGLE_PARAMETER = DiagnosticType_Error("JSC_INVALID_TOGGLE_PARAMETER", "goog.readToggleInternalDoNotCallDirectly must be called with a string literal.")
	static_ReplaceToggles_INVALID_ORDINAL_MAPPING = DiagnosticType_Error("JSC_INVALID_ORDINAL_MAPPING", "CLOSURE_TOGGLE_ORDINALS must be initialized with an object literal mapping strings to booleans or unique whole numbers: {0}")
	static_ReplaceToggles_readToggleFunctionName = pkg_rhino.QualifiedName_Of("goog.readToggleInternalDoNotCallDirectly")
}

func (this_ *ReplaceToggles) initFields_ReplaceToggles() {
	_ = this_

}

// Java constructor: ReplaceToggles(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewReplaceToggles(compiler AbstractCompiler) *ReplaceToggles {
	clinit_ReplaceToggles()
	this_ := &ReplaceToggles{}
	InitReplaceToggles(this_, compiler)
	return this_
}

func InitReplaceToggles(this_ *ReplaceToggles, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_ReplaceToggles()
	this_.compiler = compiler
	this_.astFactory = compiler.CreateAstFactory()
}

// Java method: ReplaceToggles.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ReplaceToggles) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, NewReplaceToggles_Traversal(this_))
}

// Java class: com.google.javascript.jscomp.ReplaceToggles.Traversal
type ReplaceToggles_Traversal struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *ReplaceToggles
}

func (this_ *ReplaceToggles_Traversal) BindSelf_(s any) {
	_ = s
}
func (this_ *ReplaceToggles_Traversal) IsJavaNil_() bool { return this_ == nil }
func (*ReplaceToggles_Traversal) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ReplaceToggles$Traversal", "Traversal"
}
func (this_ *ReplaceToggles_Traversal) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*ReplaceToggles_Traversal) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ReplaceToggles.Traversal", "com.google.javascript.jscomp.ReplaceToggles$Traversal", func(v_ any) bool {
		return j2gort.InstanceOf[*ReplaceToggles_Traversal](v_)
	})
}
func (this_ *ReplaceToggles_Traversal) SetOuter_(o *ReplaceToggles)                          { this_.outer_ = o }
func (this_ *ReplaceToggles_Traversal) GetOuter_() *ReplaceToggles                           { return this_.outer_ }
func (*ReplaceToggles_Traversal) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*ReplaceToggles_Traversal) Implements_jscomp_NodeTraversal_Callback()                  {}

func (this_ *ReplaceToggles_Traversal) initFields_ReplaceToggles_Traversal() {
	_ = this_
}

// Java constructor: Traversal()
func NewReplaceToggles_Traversal(outer_ *ReplaceToggles) *ReplaceToggles_Traversal {
	this_ := &ReplaceToggles_Traversal{}
	this_.outer_ = outer_
	InitReplaceToggles_Traversal(this_)
	return this_
}

func InitReplaceToggles_Traversal(this_ *ReplaceToggles_Traversal) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_ReplaceToggles_Traversal()
}

// Java method: Traversal.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ReplaceToggles_Traversal) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if NodeUtil_IsNameDeclaration(n) && n.GetFirstChild().MatchesName_1_string("CLOSURE_TOGGLE_ORDINALS") {
		var rhs pkg_rhino.Node = n.GetFirstFirstChild()
		_ = rhs
		if (rhs == 0) && n.IsVar() {
			return
		} else if !(rhs.IsObjectLit()) {
			this_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_ReplaceToggles_INVALID_ORDINAL_MAPPING(), []string{"not an object literal"}))
			return
		} else if this_.outer_.ordinalMapping != nil {
			this_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_ReplaceToggles_INVALID_ORDINAL_MAPPING(), []string{"multiple initialized copies"}))
			return
		}
		var mapping j2gort.OrderedMap[string, *j2gort.BoxInt] = j2gort.NewOrderedMap[string, *j2gort.BoxInt]()
		_ = mapping
		var ordinals j2gort.LocalOrderedSet[int32] = j2gort.NewLocalOrderedSet[int32]()
		_ = ordinals
		{
			var c pkg_rhino.Node = rhs.GetFirstChild()
			_ = c
			for ; c != 0; c = c.GetNext() {
				if !(c.IsStringKey()) && !(c.IsStringLit()) {
					this_.outer_.compiler.Report_1(JSError_Make_3(c, StaticGet_ReplaceToggles_INVALID_ORDINAL_MAPPING(), []string{"non-string key"}))
					return
				}
				var key string = c.GetString()
				_ = key
				if mapping.ContainsKey(key) {
					this_.outer_.compiler.Report_1(JSError_Make_3(c, StaticGet_ReplaceToggles_INVALID_ORDINAL_MAPPING(), []string{("duplicate key: " + j2gort.StrValString(key))}))
					return
				}
				var child pkg_rhino.Node = c.GetFirstChild()
				_ = child
				var doubleValue *j2gort.BoxDouble = NodeUtil_GetNumberValue(child)
				_ = doubleValue
				var intValue int32
				_ = intValue
				if doubleValue != nil {
					intValue = doubleValue.IntValue()
				} else {
					intValue = int32(-1)
				}
				if child.IsTrue() || child.IsFalse() {
					if child.IsTrue() {
						intValue = int32(-2)
					} else {
						intValue = int32(-1)
					}
				} else if (!(child.IsNumber()) || (intValue < int32(0))) || (float64(intValue) != j2gort.UnboxDouble(doubleValue)) {
					this_.outer_.compiler.Report_1(JSError_Make_3(c, StaticGet_ReplaceToggles_INVALID_ORDINAL_MAPPING(), []string{"value not a boolean or whole number literal"}))
					return
				} else if ordinals.Contains(intValue) {
					this_.outer_.compiler.Report_1(JSError_Make_3(c, StaticGet_ReplaceToggles_INVALID_ORDINAL_MAPPING(), []string{("duplicate ordinal: " + j2gort.StrValInt32(intValue))}))
					return
				}
				mapping.Put(key, j2gort.ValueOfInt(intValue))
				ordinals.Add(intValue)
			}
		}
		{
			r_1 := this_.outer_
			r_1.ordinalMapping = mapping.ToImmutableMap()
		}
		return
	}
	if !(n.IsCall()) {
		return
	}
	var qname pkg_rhino.Node = NodeUtil_GetCallTargetResolvingIndirectCalls(n)
	_ = qname
	if !(StaticGet_ReplaceToggles_ReadToggleFunctionName().Matches(qname)) {
		return
	}
	var arg pkg_rhino.Node = n.GetSecondChild()
	_ = arg
	if ((arg == 0) || !(arg.IsStringLit())) || !(n.HasTwoChildren()) {
		this_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_ReplaceToggles_INVALID_TOGGLE_PARAMETER(), []string{}))
		return
	}
	var ordinal *j2gort.BoxInt
	_ = ordinal
	if this_.outer_.ordinalMapping != nil {
		ordinal = j2gort.CastPtr[j2gort.BoxInt](j2gort.MapGetStr(this_.outer_.ordinalMapping, arg.GetString()))
	} else {
		ordinal = nil
	}
	if (ordinal == nil) || (j2gort.UnboxInt(ordinal) < int32(0)) {
		n.ReplaceWith(this_.outer_.astFactory.CreateBoolean(((ordinal != nil) && (j2gort.UnboxInt(ordinal) == int32(-2)))).SrcrefTreeIfMissing(n))
		t.ReportCodeChange_0()
		return
	}
	var index int32 = (j2gort.UnboxInt(ordinal) / int32(30))
	_ = index
	var bit int32 = (j2gort.UnboxInt(ordinal) % int32(30))
	_ = bit
	var getElem pkg_rhino.Node = this_.outer_.astFactory.CreateGetElem(this_.outer_.astFactory.CreateQNameWithUnknownType_1_string("goog.TOGGLES_"), this_.outer_.astFactory.CreateNumber(float64(index)))
	_ = getElem
	var bitAnd pkg_rhino.Node = pkg_rhino.Node(0)
	_ = bitAnd
	if bit < int32(14) {
		bitAnd = this_.outer_.astFactory.CreateBitwiseAnd(getElem, this_.outer_.astFactory.CreateNumber(float64(javalang.Shl32(int32(1), int(bit)))))
	} else {
		bitAnd = this_.outer_.astFactory.CreateBitwiseAnd(this_.outer_.astFactory.CreateRightShift(getElem, this_.outer_.astFactory.CreateNumber(float64(bit))), this_.outer_.astFactory.CreateNumber(float64(int32(1))))
	}
	n.ReplaceWith(this_.outer_.astFactory.CreateNot(this_.outer_.astFactory.CreateNot(bitAnd)).SrcrefTreeIfMissing(n))
	t.ReportCodeChange_0()
}
func (*ReplaceToggles_Traversal) AlwaysShouldTraverse_()      {}
func (*ReplaceToggles_Traversal) CanReleaseSyntacticScopes_() {}
