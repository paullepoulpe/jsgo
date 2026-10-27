// Port of com/google/javascript/jscomp/InlineCostEstimator.java @ closure-compiler v20261021

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
	"math"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// InlineCostEstimator is for use with CodeGenerator to determine the cost of
// generated code.
//
// See CodeGenerator, CodePrinter.
type InlineCostEstimator struct{}

// For now simply assume identifiers are 2 characters.
const InlineCostEstimator_ESTIMATED_IDENTIFIER = "ab"

// Java: InlineCostEstimator.ESTIMATED_IDENTIFIER_COST
const InlineCostEstimator_ESTIMATED_IDENTIFIER_COST = len(InlineCostEstimator_ESTIMATED_IDENTIFIER)

// Java: InlineCostEstimator()
func NewInlineCostEstimator() *InlineCostEstimator {
	return &InlineCostEstimator{}
}

// InlineCostEstimator_GetCost determines the size of the JS code.
//
// Java: InlineCostEstimator.getCost(Node)
func InlineCostEstimator_GetCost(root rhino.Node) int {
	return InlineCostEstimator_GetCostWithThreshold(root, math.MaxInt32)
}

// InlineCostEstimator_GetCostWithThreshold determines the estimated size of
// the JS snippet represented by the node.
//
// Java: InlineCostEstimator.getCost(Node, int)
func InlineCostEstimator_GetCostWithThreshold(root rhino.Node, costThreshold int) int {
	var est InlineCostEstimator_CompiledSizeEstimator
	est.InitCodeConsumerBase(&est)
	est.maxCost = costThreshold
	est.continueProcessing = true
	var cg CodeGenerator
	cg.InitCodeGenerator(&est)
	cg.AddNode(root)
	return est.cost
}

// InlineCostEstimator_CompiledSizeEstimator is a code consumer that
// estimates compiled size by assuming names are shortened and all whitespace
// is stripped.
type InlineCostEstimator_CompiledSizeEstimator struct {
	CodeConsumerBase

	maxCost            int
	cost               int
	last               uint16
	continueProcessing bool
}

// Java: InlineCostEstimator.CompiledSizeEstimator(int)
func NewInlineCostEstimator_CompiledSizeEstimator(costThreshold int) *InlineCostEstimator_CompiledSizeEstimator {
	e := &InlineCostEstimator_CompiledSizeEstimator{
		cost:               0,
		last:               0,
		continueProcessing: true,
	}
	e.InitCodeConsumerBase(e)
	e.maxCost = costThreshold
	return e
}

// AddNode is the Java add(Node) overload.
//
// Java: InlineCostEstimator.CompiledSizeEstimator.add(Node)
func (e *InlineCostEstimator_CompiledSizeEstimator) AddNode(root rhino.Node) {
	cg := CodeGenerator_ForCostEstimation(e)
	cg.AddNode(root)
}

// Java: InlineCostEstimator.CompiledSizeEstimator.getCost()
func (e *InlineCostEstimator_CompiledSizeEstimator) GetCost() int {
	return e.cost
}

// Java: InlineCostEstimator.CompiledSizeEstimator.continueProcessing()
func (e *InlineCostEstimator_CompiledSizeEstimator) ContinueProcessing() bool {
	return e.continueProcessing
}

// Java: InlineCostEstimator.CompiledSizeEstimator.getLastChar()
func (e *InlineCostEstimator_CompiledSizeEstimator) GetLastChar() uint16 {
	return e.last
}

// Java: InlineCostEstimator.CompiledSizeEstimator.append(String)
func (e *InlineCostEstimator_CompiledSizeEstimator) Append(str string) {
	length := jutil.StringLength(str)
	e.last = jutil.StringCharAt(str, length-1)
	e.cost += length
	if e.maxCost <= e.cost {
		e.continueProcessing = false
	}
}

// Java: InlineCostEstimator.CompiledSizeEstimator.addIdentifier(String)
func (e *InlineCostEstimator_CompiledSizeEstimator) AddIdentifier(identifier string) {
	e.Add(InlineCostEstimator_ESTIMATED_IDENTIFIER)
}

// AddConstant: constants (true, false, null) are considered basically free,
// because it's likely that they will get folded when we're done.
//
// Java: InlineCostEstimator.CompiledSizeEstimator.addConstant(String)
func (e *InlineCostEstimator_CompiledSizeEstimator) AddConstant(newcode string) {
	e.Add("0")
}
