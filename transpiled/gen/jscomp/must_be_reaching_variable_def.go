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
// Port of com/google/javascript/jscomp/MustBeReachingVariableDef.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.MustBeReachingVariableDef
type MustBeReachingVariableDef struct {
	DataFlowAnalysisBase[pkg_rhino.Node, *MustBeReachingVariableDef_MustDef]
	compiler         AbstractCompiler                                  `gson:"compiler"`
	varNames         []string                                          `gson:"varNames"`
	vars             []*Var                                            `gson:"vars"`
	isEscaped        []bool                                            `gson:"isEscaped"`
	paramIndices     []int32                                           `gson:"paramIndices"`
	nameTable        []string                                          `gson:"nameTable"`
	indexTable       []int32                                           `gson:"indexTable"`
	emptyLattice     *MustBeReachingVariableDef_MustDef                `gson:"emptyLattice"`
	flowOutput       *MustBeReachingVariableDef_MustDef                `gson:"flowOutput"`
	currentDepTarget *MustBeReachingVariableDef_Definition             `gson:"currentDepTarget"`
	depCollector     ControlFlowGraph_AbstractCfgNodeTraversalCallback `gson:"depCollector"`
	numVars          int32                                             `gson:"numVars"`
}

func (this_ *MustBeReachingVariableDef) BindSelf_(s any) {
	_ = s
	this_.DataFlowAnalysisBase.BindSelf_(s)
}
func (this_ *MustBeReachingVariableDef) IsJavaNil_() bool { return this_ == nil }
func (*MustBeReachingVariableDef) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.MustBeReachingVariableDef", "MustBeReachingVariableDef"
}
func (this_ *MustBeReachingVariableDef) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.DataFlowAnalysisBase.IdentityHash_()
}
func (*MustBeReachingVariableDef) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MustBeReachingVariableDef", "com.google.javascript.jscomp.MustBeReachingVariableDef", func(v_ any) bool {
		return j2gort.InstanceOf[*MustBeReachingVariableDef](v_)
	})
}
func (*MustBeReachingVariableDef) Implements_jscomp_DataFlowAnalysis() {}

var static_MustBeReachingVariableDef_EMPTY_INT_ARRAY []int32

func StaticGet_MustBeReachingVariableDef_EMPTY_INT_ARRAY() []int32 {
	clinit_MustBeReachingVariableDef()
	return static_MustBeReachingVariableDef_EMPTY_INT_ARRAY
}
func StaticSet_MustBeReachingVariableDef_EMPTY_INT_ARRAY(v []int32) {
	clinit_MustBeReachingVariableDef()
	static_MustBeReachingVariableDef_EMPTY_INT_ARRAY = v
}

var clinitDone_MustBeReachingVariableDef uint32
var clinitState_MustBeReachingVariableDef j2gort.ClinitState

func clinit_MustBeReachingVariableDef() {
	if atomic.LoadUint32(&clinitDone_MustBeReachingVariableDef) != 0 {
		return
	}
	clinitSlow_MustBeReachingVariableDef()
}

func clinitSlow_MustBeReachingVariableDef() {
	if !j2gort.EnterClinit(&clinitState_MustBeReachingVariableDef, &clinitDone_MustBeReachingVariableDef) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_MustBeReachingVariableDef, &clinitDone_MustBeReachingVariableDef)
	static_MustBeReachingVariableDef_EMPTY_INT_ARRAY = make([]int32, int32(0))
}

func (this_ *MustBeReachingVariableDef) initFields_MustBeReachingVariableDef() {
	_ = this_
	this_.depCollector = NewMustBeReachingVariableDef__Anon_63(this_)
}

// Java constructor: MustBeReachingVariableDef(com.google.javascript.jscomp.ControlFlowGraph<com.google.javascript.rhino.Node> cfg, com.google.javascript.jscomp.AbstractCompiler compiler, java.util.Set<com.google.javascript.jscomp.Var> escaped, java.util.Map<java.lang.String,com.google.javascript.jscomp.Var> allVarsInFn)
func NewMustBeReachingVariableDef(cfg ControlFlowGraph[pkg_rhino.Node], compiler AbstractCompiler, escaped j2gort.Ext_java_util_Set, allVarsInFn j2gort.Ext_java_util_Map) *MustBeReachingVariableDef {
	clinit_MustBeReachingVariableDef()
	this_ := &MustBeReachingVariableDef{}
	this_.BindSelf_(this_)
	InitMustBeReachingVariableDef(this_, cfg, compiler, escaped, allVarsInFn)
	return this_
}

func InitMustBeReachingVariableDef(this_ *MustBeReachingVariableDef, cfg ControlFlowGraph[pkg_rhino.Node], compiler AbstractCompiler, escaped j2gort.Ext_java_util_Set, allVarsInFn j2gort.Ext_java_util_Map) {
	_ = cfg
	_ = compiler
	_ = escaped
	_ = allVarsInFn
	InitDataFlowAnalysis[pkg_rhino.Node, *MustBeReachingVariableDef_MustDef](&this_.DataFlowAnalysisBase, cfg)
	this_.initFields_MustBeReachingVariableDef()
	this_.compiler = compiler
	var n int32 = allVarsInFn.Size()
	_ = n
	this_.numVars = n
	this_.varNames = j2gort.NewStringArray(n)
	this_.vars = make([]*Var, n)
	this_.isEscaped = make([]bool, n)
	var paramCount int32 = int32(0)
	_ = paramCount
	var i int32 = int32(0)
	_ = i
	{
		cur_1 := j2gort.NewMapEntryCursor(allVarsInFn)
		for cur_1.NextEntry() {
			var name string = cur_1.KeyStr()
			_ = name
			var v *Var = j2gort.CastPtr[Var](cur_1.GetValue())
			_ = v
			{
				a_2 := this_.varNames
				i_3 := i
				a_2[i_3] = name
			}
			{
				a_4 := this_.vars
				i_5 := i
				a_4[i_5] = v
			}
			{
				a_6 := this_.isEscaped
				i_7 := i
				a_6[i_7] = escaped.Contains(j2gort.PtrToAny(v))
			}
			if v.IsParam() {
				paramCount = paramCount + 1
			}
			i = i + 1
		}
	}
	if paramCount == int32(0) {
		this_.paramIndices = StaticGet_MustBeReachingVariableDef_EMPTY_INT_ARRAY()
	} else {
		var pIndices []int32 = make([]int32, paramCount)
		_ = pIndices
		var p int32 = int32(0)
		_ = p
		{
			var j int32 = int32(0)
			_ = j
			for ; j < n; j = j + 1 {
				if this_.vars[j].IsParam() {
					{
						a_8 := pIndices
						var i_9 int32
						{

							old_10 := p
							nw_11 := old_10 + 1
							p = nw_11
							i_9 = old_10
						}
						a_8[i_9] = j
					}
				}
			}
		}
		this_.paramIndices = pIndices
	}
	if n > int32(8) {
		var cap_ int32 = javalang.Shl32(j2gort.IntHighestOneBit(((n * int32(2)) - int32(1))), int(int32(1)))
		_ = cap_
		var mask int32 = (cap_ - int32(1))
		_ = mask
		var keys []string = j2gort.NewStringArray(cap_)
		_ = keys
		var indices []int32 = make([]int32, cap_)
		_ = indices
		{
			var j_2 int32 = int32(0)
			_ = j_2
			for ; j_2 < n; j_2 = j_2 + 1 {
				var key string = this_.varNames[j_2]
				_ = key
				var slot int32 = (javalang.StringHashCode(j2gort.CheckNotNullString(key)) & mask)
				_ = slot
				for keys[slot] != j2gort.NullString {
					slot = ((slot + int32(1)) & mask)
				}
				{
					a_12 := keys
					i_13 := slot
					a_12[i_13] = key
				}
				{
					a_14 := indices
					i_15 := slot
					a_14[i_15] = j_2
				}
			}
		}
		this_.nameTable = keys
		this_.indexTable = indices
	} else {
		this_.nameTable = nil
		this_.indexTable = nil
	}
	this_.emptyLattice = NewMustBeReachingVariableDef_MustDef_1_int32(n)
}

// Java method: MustBeReachingVariableDef.lookupVarIndex(java.lang.String name)
func (this_ *MustBeReachingVariableDef) priv_MustBeReachingVariableDef_lookupVarIndex_3911(name string) int32 {
	_ = name
	var table []string = this_.nameTable
	_ = table
	if table != nil {
		var mask int32 = (int32(len(table)) - int32(1))
		_ = mask
		var slot int32 = (javalang.StringHashCode(j2gort.CheckNotNullString(name)) & mask)
		_ = slot
		for {
			var k string = table[slot]
			_ = k
			if k == j2gort.NullString {
				return int32(-1)
			}
			if j2gort.CheckNotNullString(k) == name {
				return this_.indexTable[slot]
			}
			slot = ((slot + int32(1)) & mask)
		}
	}
	var names []string = this_.varNames
	_ = names
	{
		var i int32 = int32(0)
		_ = i
		var len_ int32 = int32(len(names))
		_ = len_
		for ; i < len_; i = i + 1 {
			if j2gort.CheckNotNullString(names[i]) == name {
				return i
			}
		}
	}
	return int32(-1)
}

// Java method: MustBeReachingVariableDef.isForward()
func (this_ *MustBeReachingVariableDef) IsForward() bool {
	return true
}

// Java method: MustBeReachingVariableDef.needsNodeAnnotations()
func (this_ *MustBeReachingVariableDef) NeedsNodeAnnotations() bool {
	return false
}

// Java method: MustBeReachingVariableDef.isLatticeEqual(com.google.javascript.jscomp.MustBeReachingVariableDef.MustDef a, com.google.javascript.jscomp.MustBeReachingVariableDef.MustDef b)
func (this_ *MustBeReachingVariableDef) IsLatticeEqual(a *MustBeReachingVariableDef_MustDef, b *MustBeReachingVariableDef_MustDef) bool {
	_ = a
	_ = b
	return a.EqualsMustDef(b)
}

// Java method: MustBeReachingVariableDef.createEntryLattice()
func (this_ *MustBeReachingVariableDef) CreateEntryLattice() *MustBeReachingVariableDef_MustDef {
	return NewMustBeReachingVariableDef_MustDef_1_Var(this_.vars)
}

// Java method: MustBeReachingVariableDef.createInitialEstimateLattice()
func (this_ *MustBeReachingVariableDef) CreateInitialEstimateLattice() *MustBeReachingVariableDef_MustDef {
	return this_.emptyLattice
}

// Java method: MustBeReachingVariableDef.createFlowJoiner()
func (this_ *MustBeReachingVariableDef) CreateFlowJoiner() DataFlowAnalysis_FlowJoiner[*MustBeReachingVariableDef_MustDef] {
	return ToIface_DataFlowAnalysis_FlowJoiner[*MustBeReachingVariableDef_MustDef](NewMustBeReachingVariableDef_MustDefJoin(this_.numVars))
}

// Java method: MustBeReachingVariableDef.flowThrough(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.MustBeReachingVariableDef.MustDef input)
func (this_ *MustBeReachingVariableDef) FlowThrough(n pkg_rhino.Node, input *MustBeReachingVariableDef_MustDef) *MustBeReachingVariableDef_MustDef {
	_ = n
	_ = input
	this_.flowOutput = nil
	this_.priv_MustBeReachingVariableDef_computeMustDef_3913(n, n, input, false)
	var out *MustBeReachingVariableDef_MustDef = this_.flowOutput
	_ = out
	this_.flowOutput = nil
	if out != nil {
		return out
	} else {
		return input
	}
}

// Java method: MustBeReachingVariableDef.writableFlowOutput(com.google.javascript.jscomp.MustBeReachingVariableDef.MustDef input)
func (this_ *MustBeReachingVariableDef) priv_MustBeReachingVariableDef_writableFlowOutput_3912(input *MustBeReachingVariableDef_MustDef) *MustBeReachingVariableDef_MustDef {
	_ = input
	var out *MustBeReachingVariableDef_MustDef = this_.flowOutput
	_ = out
	if out == nil {
		out = NewMustBeReachingVariableDef_MustDef_1_MustBeReachingVariableDef_MustDef(input)
		this_.flowOutput = out
	}
	return out
}

// Java method: MustBeReachingVariableDef.computeMustDef(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node cfgNode, com.google.javascript.jscomp.MustBeReachingVariableDef.MustDef output, boolean conditional)
func (this_ *MustBeReachingVariableDef) priv_MustBeReachingVariableDef_computeMustDef_3913(n pkg_rhino.Node, cfgNode pkg_rhino.Node, output *MustBeReachingVariableDef_MustDef, conditional bool) {
	_ = n
	_ = cfgNode
	_ = output
	_ = conditional
	switch n.TokenSwitchOrd_() {
	case int32(90) /* Token.BLOCK */, int32(89) /* Token.ROOT */, int32(73) /* Token.FUNCTION */ :
		return
	case int32(78) /* Token.WHILE */, int32(79) /* Token.DO */, int32(74) /* Token.IF */, int32(80) /* Token.FOR */ :
		this_.priv_MustBeReachingVariableDef_computeMustDef_3913(NodeUtil_GetConditionExpression(n), cfgNode, output, conditional)
		return
	case int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */ :
		var lhs pkg_rhino.Node = n.GetFirstChild()
		_ = lhs
		var rhs pkg_rhino.Node = lhs.GetNext()
		_ = rhs
		if NodeUtil_IsNameDeclaration(lhs) {
			lhs = lhs.GetLastChild()
		}
		if lhs.IsName() {
			this_.priv_MustBeReachingVariableDef_addToDefIfLocal_3914(lhs.GetString(), pkg_rhino.Node(0), rhs, output)
		} else if lhs.IsDestructuringLhs() {
			lhs = lhs.GetFirstChild()
		}
		if lhs.IsDestructuringPattern() {
			this_.priv_MustBeReachingVariableDef_computeMustDef_3913(lhs, cfgNode, output, true)
		}
		return
	case int32(29) /* Token.OPTCHAIN_GETPROP */ :
		this_.priv_MustBeReachingVariableDef_computeMustDef_3913(n.GetFirstChild(), cfgNode, output, conditional)
		return
	case int32(69) /* Token.AND */, int32(68) /* Token.OR */, int32(70) /* Token.COALESCE */, int32(30) /* Token.OPTCHAIN_GETELEM */ :
		this_.priv_MustBeReachingVariableDef_computeMustDef_3913(n.GetFirstChild(), cfgNode, output, conditional)
		this_.priv_MustBeReachingVariableDef_computeMustDef_3913(n.GetLastChild(), cfgNode, output, true)
		return
	case int32(31) /* Token.OPTCHAIN_CALL */ :
		this_.priv_MustBeReachingVariableDef_computeMustDef_3913(n.GetFirstChild(), cfgNode, output, conditional)
		{
			var c pkg_rhino.Node = n.GetSecondChild()
			_ = c
			for ; c != 0; c = c.GetNext() {
				this_.priv_MustBeReachingVariableDef_computeMustDef_3913(c, cfgNode, output, true)
			}
		}
		return
	case int32(67) /* Token.HOOK */ :
		this_.priv_MustBeReachingVariableDef_computeMustDef_3913(n.GetFirstChild(), cfgNode, output, conditional)
		this_.priv_MustBeReachingVariableDef_computeMustDef_3913(n.GetSecondChild(), cfgNode, output, true)
		this_.priv_MustBeReachingVariableDef_computeMustDef_3913(n.GetLastChild(), cfgNode, output, true)
		return
	case int32(111) /* Token.LET */, int32(97) /* Token.CONST */, int32(84) /* Token.VAR */ :
		{
			var c_2 pkg_rhino.Node = n.GetFirstChild()
			_ = c_2
			for ; c_2 != 0; c_2 = c_2.GetNext() {
				if c_2.HasChildren() {
					if c_2.IsName() {
						this_.priv_MustBeReachingVariableDef_computeMustDef_3913(c_2.GetFirstChild(), cfgNode, output, conditional)
						this_.priv_MustBeReachingVariableDef_addToDefIfLocal_3914(c_2.GetString(), func() pkg_rhino.Node {
							if conditional {
								return pkg_rhino.Node(0)
							}
							return cfgNode
						}(), c_2.GetFirstChild(), output)
					} else {
						if !(c_2.IsDestructuringLhs()) {
							j2gort.CheckStateMsg(false, j2gort.BoxNode(c_2))
						}
						this_.priv_MustBeReachingVariableDef_computeMustDef_3913(c_2.GetSecondChild(), cfgNode, output, conditional)
						this_.priv_MustBeReachingVariableDef_computeMustDef_3913(c_2.GetFirstChild(), cfgNode, output, conditional)
					}
				}
			}
		}
		return
	case int32(134) /* Token.DEFAULT_VALUE */ :
		if n.GetFirstChild().IsDestructuringPattern() {
			this_.priv_MustBeReachingVariableDef_computeMustDef_3913(n.GetSecondChild(), cfgNode, output, true)
			this_.priv_MustBeReachingVariableDef_computeMustDef_3913(n.GetFirstChild(), cfgNode, output, conditional)
		} else if n.GetFirstChild().IsName() {
			this_.priv_MustBeReachingVariableDef_computeMustDef_3913(n.GetSecondChild(), cfgNode, output, true)
			this_.priv_MustBeReachingVariableDef_addToDefIfLocal_3914(n.GetFirstChild().GetString(), func() pkg_rhino.Node {
				if conditional {
					return pkg_rhino.Node(0)
				}
				return cfgNode
			}(), pkg_rhino.Node(0), output)
		} else {
			this_.priv_MustBeReachingVariableDef_computeMustDef_3913(n.GetFirstChild(), cfgNode, output, conditional)
			this_.priv_MustBeReachingVariableDef_computeMustDef_3913(n.GetSecondChild(), cfgNode, output, true)
		}
	case int32(32) /* Token.NAME */ :
		if NodeUtil_IsLhsByDestructuring(n) {
			this_.priv_MustBeReachingVariableDef_addToDefIfLocal_3914(n.GetString(), func() pkg_rhino.Node {
				if conditional {
					return pkg_rhino.Node(0)
				}
				return cfgNode
			}(), pkg_rhino.Node(0), output)
		} else {
			var iifeCond_1 bool
			{

				sr_16 := "arguments"
				se_17 := n.GetString()
				iifeCond_1 = j2gort.CheckNotNullString(sr_16) == se_17
			}
			if iifeCond_1 {
				this_.priv_MustBeReachingVariableDef_escapeParameters_3915(output)
			}
		}
		return
	default:
		if NodeUtil_IsAssignmentOp(n) {
			if n.GetFirstChild().IsName() {
				var name pkg_rhino.Node = n.GetFirstChild()
				_ = name
				this_.priv_MustBeReachingVariableDef_computeMustDef_3913(name.GetNext(), cfgNode, output, conditional)
				this_.priv_MustBeReachingVariableDef_addToDefIfLocal_3914(name.GetString(), func() pkg_rhino.Node {
					if conditional {
						return pkg_rhino.Node(0)
					}
					return cfgNode
				}(), n.GetLastChild(), output)
				return
			} else if NodeUtil_IsNormalGet(n.GetFirstChild()) {
				var obj pkg_rhino.Node = n.GetFirstFirstChild()
				_ = obj
				if obj.IsName() && func() bool {
					sr_18 := "arguments"
					se_19 := obj.GetString()
					return j2gort.CheckNotNullString(sr_18) == se_19
				}() {
					this_.priv_MustBeReachingVariableDef_escapeParameters_3915(output)
				}
			} else if n.GetFirstChild().IsDestructuringPattern() {
				this_.priv_MustBeReachingVariableDef_computeMustDef_3913(n.GetSecondChild(), cfgNode, output, conditional)
				this_.priv_MustBeReachingVariableDef_computeMustDef_3913(n.GetFirstChild(), cfgNode, output, conditional)
				return
			}
		}
		if n.IsDec() || n.IsInc() {
			var target pkg_rhino.Node = n.GetFirstChild()
			_ = target
			if target.IsName() {
				this_.priv_MustBeReachingVariableDef_addToDefIfLocal_3914(target.GetString(), func() pkg_rhino.Node {
					if conditional {
						return pkg_rhino.Node(0)
					}
					return cfgNode
				}(), pkg_rhino.Node(0), output)
				return
			}
		}
		{
			var c_3 pkg_rhino.Node = n.GetFirstChild()
			_ = c_3
			for ; c_3 != 0; c_3 = c_3.GetNext() {
				this_.priv_MustBeReachingVariableDef_computeMustDef_3913(c_3, cfgNode, output, conditional)
			}
		}
	}
}

// Java method: MustBeReachingVariableDef.addToDefIfLocal(java.lang.String name, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node node, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node rValue, com.google.javascript.jscomp.MustBeReachingVariableDef.MustDef input)
func (this_ *MustBeReachingVariableDef) priv_MustBeReachingVariableDef_addToDefIfLocal_3914(name string, node pkg_rhino.Node, rValue pkg_rhino.Node, input *MustBeReachingVariableDef_MustDef) {
	_ = name
	_ = node
	_ = rValue
	_ = input
	var varIdx int32 = this_.priv_MustBeReachingVariableDef_lookupVarIndex_3911(name)
	_ = varIdx
	if varIdx < int32(0) {
		return
	}
	var var_ *Var = this_.vars[varIdx]
	_ = var_
	var def *MustBeReachingVariableDef_MustDef = this_.priv_MustBeReachingVariableDef_writableFlowOutput_3912(input)
	_ = def
	var defs []*MustBeReachingVariableDef_Definition = def.reachingDef
	_ = defs
	{
		var i int32 = int32(0)
		_ = i
		var len_ int32 = int32(len(defs))
		_ = len_
		for ; i < len_; i = i + 1 {
			var otherDef *MustBeReachingVariableDef_Definition = defs[i]
			_ = otherDef
			if (otherDef == nil) || (otherDef == StaticGet_MustBeReachingVariableDef_Definition_BOTTOM()) {
				continue
			}
			if otherDef.HasDependency(var_) {
				{
					a_20 := defs
					i_21 := i
					a_20[i_21] = StaticGet_MustBeReachingVariableDef_Definition_BOTTOM()
				}
			}
		}
	}
	if !(this_.isEscaped[varIdx]) {
		if node == 0 {
			{
				a_22 := defs
				i_23 := varIdx
				a_22[i_23] = StaticGet_MustBeReachingVariableDef_Definition_BOTTOM()
			}
		} else {
			clinit_MustBeReachingVariableDef_Definition()
			var val_24 MustBeReachingVariableDef_Definition
			InitMustBeReachingVariableDef_Definition(&val_24, node)
			var definition *MustBeReachingVariableDef_Definition = &val_24
			_ = definition
			if rValue != 0 {
				this_.priv_MustBeReachingVariableDef_computeDependence_3916(definition, rValue)
			}
			{
				a_25 := defs
				i_26 := varIdx
				a_25[i_26] = definition
			}
		}
	}
}

// Java method: MustBeReachingVariableDef.escapeParameters(com.google.javascript.jscomp.MustBeReachingVariableDef.MustDef input)
func (this_ *MustBeReachingVariableDef) priv_MustBeReachingVariableDef_escapeParameters_3915(input *MustBeReachingVariableDef_MustDef) {
	_ = input
	var output *MustBeReachingVariableDef_MustDef = this_.priv_MustBeReachingVariableDef_writableFlowOutput_3912(input)
	_ = output
	var defs []*MustBeReachingVariableDef_Definition = output.reachingDef
	_ = defs
	var pIndices []int32 = this_.paramIndices
	_ = pIndices
	{
		var i int32 = int32(0)
		_ = i
		var len_ int32 = int32(len(pIndices))
		_ = len_
		for ; i < len_; i = i + 1 {
			{
				a_27 := defs
				i_28 := pIndices[i]
				a_27[i_28] = StaticGet_MustBeReachingVariableDef_Definition_BOTTOM()
			}
		}
	}
	{
		var i_2 int32 = int32(0)
		_ = i_2
		var len__2 int32 = int32(len(defs))
		_ = len__2
		for ; i_2 < len__2; i_2 = i_2 + 1 {
			var value *MustBeReachingVariableDef_Definition = defs[i_2]
			_ = value
			if (value == nil) || (value == StaticGet_MustBeReachingVariableDef_Definition_BOTTOM()) {
				continue
			}
			if value.DependsOnParameter() {
				{
					a_29 := defs
					i_30 := i_2
					a_29[i_30] = StaticGet_MustBeReachingVariableDef_Definition_BOTTOM()
				}
			}
		}
	}
}

// Java method: MustBeReachingVariableDef.computeDependence(com.google.javascript.jscomp.MustBeReachingVariableDef.Definition def, com.google.javascript.rhino.Node rValue)
func (this_ *MustBeReachingVariableDef) priv_MustBeReachingVariableDef_computeDependence_3916(def *MustBeReachingVariableDef_Definition, rValue pkg_rhino.Node) {
	_ = def
	_ = rValue
	var prevTarget *MustBeReachingVariableDef_Definition = this_.currentDepTarget
	_ = prevTarget
	this_.currentDepTarget = def
	defer func() {
		this_.currentDepTarget = prevTarget
	}()
	{
		NodeTraversal_Traverse(this_.compiler, rValue, Upcast_NodeTraversal_Callback(this_.depCollector))
	}
}

// Java method: MustBeReachingVariableDef.getDef(java.lang.String name, com.google.javascript.rhino.Node useNode)
func (this_ *MustBeReachingVariableDef) GetDef(name string, useNode pkg_rhino.Node) *MustBeReachingVariableDef_Definition {
	_ = name
	_ = useNode
	if !(this_.GetCfg().HasNode(useNode)) {
		j2gort.CheckArgument(false)
	}
	var varIdx int32 = this_.priv_MustBeReachingVariableDef_lookupVarIndex_3911(name)
	_ = varIdx
	if varIdx < int32(0) {
		return nil
	}
	var n *pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch] = j2gort.CastPtr[pkg_graph.LinkedDirectedGraph_LinkedDiGraphNode[pkg_rhino.Node, ControlFlowGraph_Branch]](this_.GetCfg().GetNode(useNode))
	_ = n
	var def *MustBeReachingVariableDef_Definition = this_.GetNodeInLattice(n).reachingDef[varIdx]
	_ = def
	if (def == nil) || (def == StaticGet_MustBeReachingVariableDef_Definition_BOTTOM()) {
		return nil
	} else {
		return def
	}
}

// Java method: MustBeReachingVariableDef.getDefNode(java.lang.String name, com.google.javascript.rhino.Node useNode)
func (this_ *MustBeReachingVariableDef) GetDefNode(name string, useNode pkg_rhino.Node) pkg_rhino.Node {
	_ = name
	_ = useNode
	var def *MustBeReachingVariableDef_Definition = this_.GetDef(name, useNode)
	_ = def
	if def == nil {
		return pkg_rhino.Node(0)
	} else {
		return def.node
	}
}

// Java method: MustBeReachingVariableDef.dependsOnOuterScopeVars(com.google.javascript.jscomp.MustBeReachingVariableDef.Definition def)
func (this_ *MustBeReachingVariableDef) DependsOnOuterScopeVars(def *MustBeReachingVariableDef_Definition) bool {
	_ = def
	if def.unknownDependencies {
		return true
	}
	return def.DependsOnCatchScopeVar()
}

// Java class: com.google.javascript.jscomp.MustBeReachingVariableDef.Definition
type MustBeReachingVariableDef_Definition struct {
	dep0                *Var                           `gson:"dep0"`
	dep1                *Var                           `gson:"dep1"`
	restDeps            *j2gort.CustomOrderedSet[*Var] `gson:"restDeps"`
	node                pkg_rhino.Node                 `gson:"node"`
	unknownDependencies bool                           `gson:"unknownDependencies"`
}

func (this_ *MustBeReachingVariableDef_Definition) BindSelf_(s any) {
	_ = s
}
func (this_ *MustBeReachingVariableDef_Definition) IsJavaNil_() bool { return this_ == nil }
func (*MustBeReachingVariableDef_Definition) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.MustBeReachingVariableDef$Definition", "Definition"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MustBeReachingVariableDef.Definition", "com.google.javascript.jscomp.MustBeReachingVariableDef$Definition", func(v_ any) bool {
		return j2gort.InstanceOf[*MustBeReachingVariableDef_Definition](v_)
	})
}

var static_MustBeReachingVariableDef_Definition_BOTTOM *MustBeReachingVariableDef_Definition

func StaticGet_MustBeReachingVariableDef_Definition_BOTTOM() *MustBeReachingVariableDef_Definition {
	clinit_MustBeReachingVariableDef_Definition()
	return static_MustBeReachingVariableDef_Definition_BOTTOM
}
func StaticSet_MustBeReachingVariableDef_Definition_BOTTOM(v *MustBeReachingVariableDef_Definition) {
	clinit_MustBeReachingVariableDef_Definition()
	static_MustBeReachingVariableDef_Definition_BOTTOM = v
}

var clinitDone_MustBeReachingVariableDef_Definition uint32
var clinitState_MustBeReachingVariableDef_Definition j2gort.ClinitState

func clinit_MustBeReachingVariableDef_Definition() {
	if atomic.LoadUint32(&clinitDone_MustBeReachingVariableDef_Definition) != 0 {
		return
	}
	clinitSlow_MustBeReachingVariableDef_Definition()
}

func clinitSlow_MustBeReachingVariableDef_Definition() {
	if !j2gort.EnterClinit(&clinitState_MustBeReachingVariableDef_Definition, &clinitDone_MustBeReachingVariableDef_Definition) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_MustBeReachingVariableDef_Definition, &clinitDone_MustBeReachingVariableDef_Definition)
	static_MustBeReachingVariableDef_Definition_BOTTOM = NewMustBeReachingVariableDef_Definition(pkg_rhino.Node(0))
}

func (this_ *MustBeReachingVariableDef_Definition) initFields_MustBeReachingVariableDef_Definition() {
	_ = this_

}

// Java constructor: Definition(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node node)
func NewMustBeReachingVariableDef_Definition(node pkg_rhino.Node) *MustBeReachingVariableDef_Definition {
	clinit_MustBeReachingVariableDef_Definition()
	this_ := &MustBeReachingVariableDef_Definition{}
	InitMustBeReachingVariableDef_Definition(this_, node)
	return this_
}

func InitMustBeReachingVariableDef_Definition(this_ *MustBeReachingVariableDef_Definition, node pkg_rhino.Node) {
	_ = node
	this_.initFields_MustBeReachingVariableDef_Definition()
	this_.node = node
}

// Java method: Definition.addDependency(com.google.javascript.jscomp.Var v)
func (this_ *MustBeReachingVariableDef_Definition) AddDependency(v *Var) {
	_ = v
	if (v == this_.dep0) || (v == this_.dep1) {
		return
	}
	if this_.dep0 == nil {
		this_.dep0 = v
		return
	}
	if this_.dep1 == nil {
		this_.dep1 = v
		return
	}
	if this_.restDeps == nil {
		this_.restDeps = j2gort.NewCustomOrderedSetPtrCap[*Var](int32(4))
	}
	this_.restDeps.Add(v)
}

// Java method: Definition.isEmptyDepends()
func (this_ *MustBeReachingVariableDef_Definition) IsEmptyDepends() bool {
	return (this_.dep0 == nil)
}

// Java method: Definition.hasDependency(com.google.javascript.jscomp.@org.jspecify.annotations.Nullable Var v)
func (this_ *MustBeReachingVariableDef_Definition) HasDependency(v *Var) bool {
	_ = v
	if (v == nil) || (this_.dep0 == nil) {
		return false
	}
	if (v == this_.dep0) || (v == this_.dep1) {
		return true
	}
	return ((this_.restDeps != nil) && this_.restDeps.Contains(v))
}

// Java method: Definition.containsAllDependencies(com.google.javascript.jscomp.MustBeReachingVariableDef.Definition other)
func (this_ *MustBeReachingVariableDef_Definition) ContainsAllDependencies(other *MustBeReachingVariableDef_Definition) bool {
	_ = other
	if (other.dep0 != nil) && !(this_.HasDependency(other.dep0)) {
		return false
	}
	if (other.dep1 != nil) && !(this_.HasDependency(other.dep1)) {
		return false
	}
	if other.restDeps != nil {
		{
			cur_31 := other.restDeps.Iter()
			for cur_31.HasNext() {

				v := cur_31.
					Next()
				_ = v
				if !(this_.HasDependency(v)) {
					return false
				}
			}
		}
	}
	return true
}

// Java method: Definition.addDependencyNamesTo(java.util.Set<java.lang.String> out)
func (this_ *MustBeReachingVariableDef_Definition) AddDependencyNamesTo(out j2gort.Ext_java_util_Set) {
	_ = out
	if this_.dep0 != nil {
		j2gort.SetAddStr(out, this_.dep0.GetName())
		if this_.dep1 != nil {
			j2gort.SetAddStr(out, this_.dep1.GetName())
			if this_.restDeps != nil {
				{
					cur_33 := this_.restDeps.Iter()
					for cur_33.HasNext() {

						v := cur_33.
							Next()
						_ = v
						j2gort.SetAddStr(out, v.GetName())
					}
				}
			}
		}
	}
}

// Java method: Definition.dependsOnParameter()
func (this_ *MustBeReachingVariableDef_Definition) DependsOnParameter() bool {
	if this_.dep0 != nil {
		if this_.dep0.IsParam() {
			return true
		}
		if this_.dep1 != nil {
			if this_.dep1.IsParam() {
				return true
			}
			if this_.restDeps != nil {
				{
					cur_35 := this_.restDeps.Iter()
					for cur_35.HasNext() {

						v := cur_35.
							Next()
						_ = v
						if v.IsParam() {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

// Java method: Definition.dependsOnCatchScopeVar()
func (this_ *MustBeReachingVariableDef_Definition) DependsOnCatchScopeVar() bool {
	if this_.dep0 != nil {
		if j2gort.CastPtr[Scope](this_.dep0.GetScope()).IsCatchScope() {
			return true
		}
		if this_.dep1 != nil {
			if j2gort.CastPtr[Scope](this_.dep1.GetScope()).IsCatchScope() {
				return true
			}
			if this_.restDeps != nil {
				{
					cur_37 := this_.restDeps.Iter()
					for cur_37.HasNext() {

						v := cur_37.
							Next()
						_ = v
						if j2gort.CastPtr[Scope](v.GetScope()).IsCatchScope() {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

// Java method: Definition.equals(java.lang.Object other)
func (this_ *MustBeReachingVariableDef_Definition) Equals(arg_39 any) bool {
	other := arg_39
	_ = other
	if j2gort.RefEquals(this_, other) {
		return true
	}
	var otherDef *MustBeReachingVariableDef_Definition
	_ = otherDef
	var ok_41 bool
	otherDef, ok_41 = j2gort.TryCast[*MustBeReachingVariableDef_Definition](other)
	ok_40 := !ok_41
	if ok_40 || (this_.node == 0) || (otherDef.node == 0) {
		return false
	}
	return (otherDef.node == this_.node)
}

// Java method: Definition.toString()
func (this_ *MustBeReachingVariableDef_Definition) String() string {
	return ("Definition@" + j2gort.StrVal(this_.node))
}

// Java method: Definition.hashCode()
func (this_ *MustBeReachingVariableDef_Definition) HashCode() int32 {
	if this_.node != 0 {
		return int32(j2gort.CheckNotNullNode(this_.node))
	} else {
		return int32(0)
	}
}

// Java class: com.google.javascript.jscomp.MustBeReachingVariableDef.MustDef
type MustBeReachingVariableDef_MustDef struct {
	reachingDef []*MustBeReachingVariableDef_Definition `gson:"reachingDef"`
}

func (this_ *MustBeReachingVariableDef_MustDef) BindSelf_(s any) {
	_ = s
}
func (this_ *MustBeReachingVariableDef_MustDef) IsJavaNil_() bool { return this_ == nil }
func (*MustBeReachingVariableDef_MustDef) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.MustBeReachingVariableDef$MustDef", "MustDef"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MustBeReachingVariableDef.MustDef", "com.google.javascript.jscomp.MustBeReachingVariableDef$MustDef", func(v_ any) bool {
		return j2gort.InstanceOf[*MustBeReachingVariableDef_MustDef](v_)
	})
}
func (*MustBeReachingVariableDef_MustDef) Implements_graph_LatticeElement() {}
func (*MustBeReachingVariableDef_MustDef) Implements_graph_Annotation()     {}

var static_MustBeReachingVariableDef_MustDef_EMPTY_DEFS []*MustBeReachingVariableDef_Definition

func StaticGet_MustBeReachingVariableDef_MustDef_EMPTY_DEFS() []*MustBeReachingVariableDef_Definition {
	clinit_MustBeReachingVariableDef_MustDef()
	return static_MustBeReachingVariableDef_MustDef_EMPTY_DEFS
}
func StaticSet_MustBeReachingVariableDef_MustDef_EMPTY_DEFS(v []*MustBeReachingVariableDef_Definition) {
	clinit_MustBeReachingVariableDef_MustDef()
	static_MustBeReachingVariableDef_MustDef_EMPTY_DEFS = v
}

var clinitDone_MustBeReachingVariableDef_MustDef uint32
var clinitState_MustBeReachingVariableDef_MustDef j2gort.ClinitState

func clinit_MustBeReachingVariableDef_MustDef() {
	if atomic.LoadUint32(&clinitDone_MustBeReachingVariableDef_MustDef) != 0 {
		return
	}
	clinitSlow_MustBeReachingVariableDef_MustDef()
}

func clinitSlow_MustBeReachingVariableDef_MustDef() {
	if !j2gort.EnterClinit(&clinitState_MustBeReachingVariableDef_MustDef, &clinitDone_MustBeReachingVariableDef_MustDef) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_MustBeReachingVariableDef_MustDef, &clinitDone_MustBeReachingVariableDef_MustDef)
	static_MustBeReachingVariableDef_MustDef_EMPTY_DEFS = make([]*MustBeReachingVariableDef_Definition, int32(0))
}

func (this_ *MustBeReachingVariableDef_MustDef) initFields_MustBeReachingVariableDef_MustDef() {
	_ = this_
}

// Java constructor: MustDef(int numVars)
func NewMustBeReachingVariableDef_MustDef_1_int32(numVars int32) *MustBeReachingVariableDef_MustDef {
	clinit_MustBeReachingVariableDef_MustDef()
	this_ := &MustBeReachingVariableDef_MustDef{}
	InitMustBeReachingVariableDef_MustDef_1_int32(this_, numVars)
	return this_
}

func InitMustBeReachingVariableDef_MustDef_1_int32(this_ *MustBeReachingVariableDef_MustDef, numVars int32) {
	_ = numVars
	this_.initFields_MustBeReachingVariableDef_MustDef()
	if numVars == int32(0) {
		this_.reachingDef = StaticGet_MustBeReachingVariableDef_MustDef_EMPTY_DEFS()
	} else {
		this_.reachingDef = make([]*MustBeReachingVariableDef_Definition, numVars)
	}
}

// Java constructor: MustDef(com.google.javascript.jscomp.Var[] vars)
func NewMustBeReachingVariableDef_MustDef_1_Var(vars []*Var) *MustBeReachingVariableDef_MustDef {
	clinit_MustBeReachingVariableDef_MustDef()
	this_ := &MustBeReachingVariableDef_MustDef{}
	InitMustBeReachingVariableDef_MustDef_1_Var(this_, vars)
	return this_
}

func InitMustBeReachingVariableDef_MustDef_1_Var(this_ *MustBeReachingVariableDef_MustDef, vars []*Var) {
	_ = vars
	this_.initFields_MustBeReachingVariableDef_MustDef()
	var len_ int32 = int32(len(vars))
	_ = len_
	if len_ == int32(0) {
		this_.reachingDef = StaticGet_MustBeReachingVariableDef_MustDef_EMPTY_DEFS()
	} else {
		var defs []*MustBeReachingVariableDef_Definition = make([]*MustBeReachingVariableDef_Definition, len_)
		_ = defs
		{
			var i int32 = int32(0)
			_ = i
			for ; i < len_; i = i + 1 {
				{
					a_42 := defs
					i_43 := i
					a_42[i_43] = NewMustBeReachingVariableDef_Definition(j2gort.CastPtr[Scope](vars[i].GetScope()).GetRootNode_0())
				}
			}
		}
		this_.reachingDef = defs
	}
}

// Java constructor: MustDef(com.google.javascript.jscomp.MustBeReachingVariableDef.MustDef other)
func NewMustBeReachingVariableDef_MustDef_1_MustBeReachingVariableDef_MustDef(other *MustBeReachingVariableDef_MustDef) *MustBeReachingVariableDef_MustDef {
	clinit_MustBeReachingVariableDef_MustDef()
	this_ := &MustBeReachingVariableDef_MustDef{}
	InitMustBeReachingVariableDef_MustDef_1_MustBeReachingVariableDef_MustDef(this_, other)
	return this_
}

func InitMustBeReachingVariableDef_MustDef_1_MustBeReachingVariableDef_MustDef(this_ *MustBeReachingVariableDef_MustDef, other *MustBeReachingVariableDef_MustDef) {
	_ = other
	this_.initFields_MustBeReachingVariableDef_MustDef()
	if int32(len(other.reachingDef)) == int32(0) {
		this_.reachingDef = StaticGet_MustBeReachingVariableDef_MustDef_EMPTY_DEFS()
	} else {
		this_.reachingDef = j2gort.CloneSlice(other.reachingDef)
	}
}

// Java method: MustDef.equals(java.lang.Object other)
func (this_ *MustBeReachingVariableDef_MustDef) Equals(arg_44 any) bool {
	other := arg_44
	_ = other
	var otherMustDef *MustBeReachingVariableDef_MustDef
	_ = otherMustDef
	var ok_45 bool
	otherMustDef, ok_45 = j2gort.TryCast[*MustBeReachingVariableDef_MustDef](other)
	if !ok_45 {
		return false
	}
	return this_.EqualsMustDef(otherMustDef)
}

// Java method: MustDef.equalsMustDef(com.google.javascript.jscomp.MustBeReachingVariableDef.MustDef otherMustDef)
func (this_ *MustBeReachingVariableDef_MustDef) EqualsMustDef(otherMustDef *MustBeReachingVariableDef_MustDef) bool {
	_ = otherMustDef
	if this_ == otherMustDef {
		return true
	}
	var a []*MustBeReachingVariableDef_Definition = this_.reachingDef
	_ = a
	var b []*MustBeReachingVariableDef_Definition = otherMustDef.reachingDef
	_ = b
	var len_ int32 = int32(len(a))
	_ = len_
	if len_ != int32(len(b)) {
		return false
	}
	{
		var i int32 = int32(0)
		_ = i
		for ; i < len_; i = i + 1 {
			var da *MustBeReachingVariableDef_Definition = a[i]
			_ = da
			var db *MustBeReachingVariableDef_Definition = b[i]
			_ = db
			if da == db {
				continue
			}
			if (((da == nil) || (db == nil)) || (da.node == 0)) || (da.node != db.node) {
				return false
			}
		}
	}
	return true
}

// Java method: MustDef.hashCode()
func (this_ *MustBeReachingVariableDef_MustDef) HashCode() int32 {
	var iifeArg_1 []any
iifeSwitch_1:
	switch {
	default:

		arr_46 := this_.reachingDef
		if arr_46 == nil {
			iifeArg_1 = nil
			break iifeSwitch_1
		}
		out_47 := make([]any, len(arr_46))
		for i_, el_48 := range arr_46 {
			out_47[i_] = j2gort.PtrToAny(el_48)
		}
		iifeArg_1 = out_47
		break iifeSwitch_1
	}
	return j2gort.Ext_java_util_Arrays_HashCode_1_slice_any(iifeArg_1)
}

// Java class: com.google.javascript.jscomp.MustBeReachingVariableDef.MustDefJoin
type MustBeReachingVariableDef_MustDefJoin struct {
	result  *MustBeReachingVariableDef_MustDef `gson:"result"`
	idHash_ atomic.Int32
}

func (this_ *MustBeReachingVariableDef_MustDefJoin) BindSelf_(s any) {
	_ = s
}
func (this_ *MustBeReachingVariableDef_MustDefJoin) IsJavaNil_() bool { return this_ == nil }
func (*MustBeReachingVariableDef_MustDefJoin) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.MustBeReachingVariableDef$MustDefJoin", "MustDefJoin"
}
func (this_ *MustBeReachingVariableDef_MustDefJoin) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*MustBeReachingVariableDef_MustDefJoin) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.MustBeReachingVariableDef.MustDefJoin", "com.google.javascript.jscomp.MustBeReachingVariableDef$MustDefJoin", func(v_ any) bool {
		return j2gort.InstanceOf[*MustBeReachingVariableDef_MustDefJoin](v_)
	})
}
func (*MustBeReachingVariableDef_MustDefJoin) Implements_jscomp_DataFlowAnalysis_FlowJoiner() {}

func (this_ *MustBeReachingVariableDef_MustDefJoin) initFields_MustBeReachingVariableDef_MustDefJoin() {
	_ = this_
}

// Java constructor: MustDefJoin(int numVars)
func NewMustBeReachingVariableDef_MustDefJoin(numVars int32) *MustBeReachingVariableDef_MustDefJoin {
	this_ := &MustBeReachingVariableDef_MustDefJoin{}
	InitMustBeReachingVariableDef_MustDefJoin(this_, numVars)
	return this_
}

func InitMustBeReachingVariableDef_MustDefJoin(this_ *MustBeReachingVariableDef_MustDefJoin, numVars int32) {
	_ = numVars
	this_.initFields_MustBeReachingVariableDef_MustDefJoin()
	this_.result = NewMustBeReachingVariableDef_MustDef_1_int32(numVars)
}

// Java method: MustDefJoin.joinFlow(com.google.javascript.jscomp.MustBeReachingVariableDef.MustDef input)
func (this_ *MustBeReachingVariableDef_MustDefJoin) JoinFlow(input *MustBeReachingVariableDef_MustDef) {
	_ = input
	var res []*MustBeReachingVariableDef_Definition = this_.result.reachingDef
	_ = res
	var inDefs []*MustBeReachingVariableDef_Definition = input.reachingDef
	_ = inDefs
	{
		var i int32 = int32(0)
		_ = i
		var len_ int32 = int32(len(res))
		_ = len_
		for ; i < len_; i = i + 1 {
			var inDef *MustBeReachingVariableDef_Definition = inDefs[i]
			_ = inDef
			if inDef == nil {
				continue
			}
			if inDef == StaticGet_MustBeReachingVariableDef_Definition_BOTTOM() {
				{
					a_49 := res
					i_50 := i
					a_49[i_50] = StaticGet_MustBeReachingVariableDef_Definition_BOTTOM()
				}
				continue
			}
			var cur *MustBeReachingVariableDef_Definition = res[i]
			_ = cur
			if cur == nil {
				{
					a_51 := res
					i_52 := i
					a_51[i_52] = inDef
				}
			} else if ((cur != StaticGet_MustBeReachingVariableDef_Definition_BOTTOM()) && (cur != inDef)) && (cur.node != inDef.node) {
				{
					a_53 := res
					i_54 := i
					a_53[i_54] = StaticGet_MustBeReachingVariableDef_Definition_BOTTOM()
				}
			}
		}
	}
}

// Java method: MustDefJoin.finish()
func (this_ *MustBeReachingVariableDef_MustDefJoin) Finish() *MustBeReachingVariableDef_MustDef {
	return this_.result
}

// Java class:
type MustBeReachingVariableDef__Anon_63 struct {
	ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase
	outer_ *MustBeReachingVariableDef
}

func (this_ *MustBeReachingVariableDef__Anon_63) BindSelf_(s any) {
	_ = s
}
func (this_ *MustBeReachingVariableDef__Anon_63) IsJavaNil_() bool { return this_ == nil }
func (*MustBeReachingVariableDef__Anon_63) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.MustBeReachingVariableDef$1", ""
}
func (this_ *MustBeReachingVariableDef__Anon_63) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase.IdentityHash_()
}
func (*MustBeReachingVariableDef__Anon_63) HasIdentityEquals_() {}
func (this_ *MustBeReachingVariableDef__Anon_63) SetOuter_(o *MustBeReachingVariableDef) {
	this_.outer_ = o
}
func (this_ *MustBeReachingVariableDef__Anon_63) GetOuter_() *MustBeReachingVariableDef {
	return this_.outer_
}
func (*MustBeReachingVariableDef__Anon_63) Implements_jscomp_ControlFlowGraph_AbstractCfgNodeTraversalCallback() {
}
func (*MustBeReachingVariableDef__Anon_63) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *MustBeReachingVariableDef__Anon_63) initFields_MustBeReachingVariableDef__Anon_63() {
	_ = this_
}

func NewMustBeReachingVariableDef__Anon_63(outer_ *MustBeReachingVariableDef) *MustBeReachingVariableDef__Anon_63 {
	this_ := &MustBeReachingVariableDef__Anon_63{}
	this_.outer_ = outer_
	InitControlFlowGraph_AbstractCfgNodeTraversalCallback(&this_.ControlFlowGraph_AbstractCfgNodeTraversalCallbackBase)
	this_.initFields_MustBeReachingVariableDef__Anon_63()
	return this_
}

// Java method: .visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *MustBeReachingVariableDef__Anon_63) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsName() {
		var idx int32 = this_.outer_.priv_MustBeReachingVariableDef_lookupVarIndex_3911(n.GetString())
		_ = idx
		var target *MustBeReachingVariableDef_Definition = this_.outer_.currentDepTarget
		_ = target
		if idx < int32(0) {
			{
				r_55 := target
				r_55.unknownDependencies = true
			}
		} else {
			target.AddDependency(this_.outer_.vars[idx])
		}
	}
}
func (*MustBeReachingVariableDef__Anon_63) CanReleaseSyntacticScopes_() {}
