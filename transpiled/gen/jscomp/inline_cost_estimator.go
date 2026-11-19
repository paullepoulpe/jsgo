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
// Port of com/google/javascript/jscomp/InlineCostEstimator.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.InlineCostEstimator
type InlineCostEstimator struct {
	idHash_ atomic.Int32
}

func (this_ *InlineCostEstimator) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineCostEstimator) IsJavaNil_() bool { return this_ == nil }
func (*InlineCostEstimator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineCostEstimator", "InlineCostEstimator"
}
func (this_ *InlineCostEstimator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InlineCostEstimator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineCostEstimator", "com.google.javascript.jscomp.InlineCostEstimator", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineCostEstimator](v_)
	})
}

func StaticGet_InlineCostEstimator_ESTIMATED_IDENTIFIER() string { return "ab" }

var static_InlineCostEstimator_ESTIMATED_IDENTIFIER_COST int32

func StaticGet_InlineCostEstimator_ESTIMATED_IDENTIFIER_COST() int32 {
	clinit_InlineCostEstimator()
	return static_InlineCostEstimator_ESTIMATED_IDENTIFIER_COST
}
func StaticSet_InlineCostEstimator_ESTIMATED_IDENTIFIER_COST(v int32) {
	clinit_InlineCostEstimator()
	static_InlineCostEstimator_ESTIMATED_IDENTIFIER_COST = v
}

var clinitDone_InlineCostEstimator uint32
var clinitState_InlineCostEstimator j2gort.ClinitState

func clinit_InlineCostEstimator() {
	if atomic.LoadUint32(&clinitDone_InlineCostEstimator) != 0 {
		return
	}
	clinitSlow_InlineCostEstimator()
}

func clinitSlow_InlineCostEstimator() {
	if !j2gort.EnterClinit(&clinitState_InlineCostEstimator, &clinitDone_InlineCostEstimator) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_InlineCostEstimator, &clinitDone_InlineCostEstimator)
	static_InlineCostEstimator_ESTIMATED_IDENTIFIER_COST = j2gort.JStrLen("ab")
}

func (this_ *InlineCostEstimator) initFields_InlineCostEstimator() {
	_ = this_
}

// Java constructor: InlineCostEstimator()
func NewInlineCostEstimator() *InlineCostEstimator {
	clinit_InlineCostEstimator()
	this_ := &InlineCostEstimator{}
	InitInlineCostEstimator(this_)
	return this_
}

func InitInlineCostEstimator(this_ *InlineCostEstimator) {
	this_.initFields_InlineCostEstimator()
}

// Java static method: InlineCostEstimator.getCost(com.google.javascript.rhino.Node root)
func InlineCostEstimator_GetCost_1(root pkg_rhino.Node) int32 {
	clinit_InlineCostEstimator()
	_ = root
	return InlineCostEstimator_GetCost_2(root, int32(2147483647))
}

// Java static method: InlineCostEstimator.getCost(com.google.javascript.rhino.Node root, int costThreshold)
func InlineCostEstimator_GetCost_2(root pkg_rhino.Node, costThreshold int32) int32 {
	clinit_InlineCostEstimator()
	_ = root
	_ = costThreshold
	var estimator *InlineCostEstimator_CompiledSizeEstimator = NewInlineCostEstimator_CompiledSizeEstimator(costThreshold)
	_ = estimator
	estimator.Add_1_Node(root)
	return estimator.GetCost()
}

// Java class: com.google.javascript.jscomp.InlineCostEstimator.CompiledSizeEstimator
type InlineCostEstimator_CompiledSizeEstimator struct {
	CodeConsumerBase
	maxCost            int32  `gson:"maxCost"`
	cost               int32  `gson:"cost"`
	last               uint16 `gson:"last"`
	continueProcessing bool   `gson:"continueProcessing"`
}

func (this_ *InlineCostEstimator_CompiledSizeEstimator) BindSelf_(s any) {
	_ = s
	this_.CodeConsumerBase.BindSelf_(s)
}
func (this_ *InlineCostEstimator_CompiledSizeEstimator) IsJavaNil_() bool { return this_ == nil }
func (*InlineCostEstimator_CompiledSizeEstimator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineCostEstimator$CompiledSizeEstimator", "CompiledSizeEstimator"
}
func (this_ *InlineCostEstimator_CompiledSizeEstimator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.CodeConsumerBase.IdentityHash_()
}
func (*InlineCostEstimator_CompiledSizeEstimator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineCostEstimator.CompiledSizeEstimator", "com.google.javascript.jscomp.InlineCostEstimator$CompiledSizeEstimator", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineCostEstimator_CompiledSizeEstimator](v_)
	})
}
func (*InlineCostEstimator_CompiledSizeEstimator) Implements_jscomp_CodeConsumer() {}

func (this_ *InlineCostEstimator_CompiledSizeEstimator) initFields_InlineCostEstimator_CompiledSizeEstimator() {
	_ = this_

	this_.continueProcessing = true
}

// Java constructor: CompiledSizeEstimator(int costThreshold)
func NewInlineCostEstimator_CompiledSizeEstimator(costThreshold int32) *InlineCostEstimator_CompiledSizeEstimator {
	this_ := &InlineCostEstimator_CompiledSizeEstimator{}
	this_.BindSelf_(this_)
	InitInlineCostEstimator_CompiledSizeEstimator(this_, costThreshold)
	return this_
}

func InitInlineCostEstimator_CompiledSizeEstimator(this_ *InlineCostEstimator_CompiledSizeEstimator, costThreshold int32) {
	_ = costThreshold
	InitCodeConsumer(&this_.CodeConsumerBase)
	this_.initFields_InlineCostEstimator_CompiledSizeEstimator()
	this_.maxCost = costThreshold
}

// Java method: CompiledSizeEstimator.add(com.google.javascript.rhino.Node root)
func (this_ *InlineCostEstimator_CompiledSizeEstimator) Add_1_Node(root pkg_rhino.Node) {
	_ = root
	var cg CodeGenerator = CodeGenerator_ForCostEstimation(this_)
	_ = cg
	cg.GetCodeGeneratorBase(

	// Java method: CompiledSizeEstimator.getCost()
	).Add_1_Node(root)
}

func (this_ *InlineCostEstimator_CompiledSizeEstimator) GetCost() int32 {
	return this_.cost
}

// Java method: CompiledSizeEstimator.continueProcessing()
func (this_ *InlineCostEstimator_CompiledSizeEstimator) ContinueProcessing() bool {
	return this_.continueProcessing
}

// Java method: CompiledSizeEstimator.getLastChar()
func (this_ *InlineCostEstimator_CompiledSizeEstimator) GetLastChar() uint16 {
	return this_.last
}

// Java method: CompiledSizeEstimator.append(java.lang.String str)
func (this_ *InlineCostEstimator_CompiledSizeEstimator) Append(str string) {
	_ = str
	this_.last = j2gort.JStrCharAt(str, (j2gort.JStrLen(str) - int32(1)))
	this_.cost = (this_.cost + j2gort.JStrLen(str))
	if this_.maxCost <= this_.cost {
		this_.continueProcessing = false
	}
}

// Java method: CompiledSizeEstimator.addIdentifier(java.lang.String identifier)
func (this_ *InlineCostEstimator_CompiledSizeEstimator) AddIdentifier(identifier string) {
	_ = identifier
	this_.Add_1_string("ab")
}

// Java method: CompiledSizeEstimator.addConstant(java.lang.String newcode)
func (this_ *InlineCostEstimator_CompiledSizeEstimator) AddConstant(newcode string) {
	_ = newcode
	this_.Add_1_string("0")
}
