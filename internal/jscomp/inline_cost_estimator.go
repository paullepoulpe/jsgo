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
	var cg CodeGenerator
	cg.initForCostEstimation(costThreshold)
	cg.AddNode(root)
	return cg.cc.cost
}

// InlineCostEstimator_CompiledSizeEstimator estimates compiled size by
// assuming names are shortened and all whitespace is stripped.
type InlineCostEstimator_CompiledSizeEstimator struct {
	cg CodeGenerator
}

// Java: InlineCostEstimator.CompiledSizeEstimator(int)
func NewInlineCostEstimator_CompiledSizeEstimator(costThreshold int) *InlineCostEstimator_CompiledSizeEstimator {
	e := &InlineCostEstimator_CompiledSizeEstimator{}
	e.cg.initForCostEstimation(costThreshold)
	return e
}

// AddNode is the Java add(Node) overload.
//
// Java: InlineCostEstimator.CompiledSizeEstimator.add(Node)
func (e *InlineCostEstimator_CompiledSizeEstimator) AddNode(root rhino.Node) {
	e.cg.AddNode(root)
}

// Java: InlineCostEstimator.CompiledSizeEstimator.add(String)
func (e *InlineCostEstimator_CompiledSizeEstimator) Add(newCode string) {
	e.cg.cc.Add(newCode)
}

// Java: InlineCostEstimator.CompiledSizeEstimator.getCost()
func (e *InlineCostEstimator_CompiledSizeEstimator) GetCost() int {
	return e.cg.cc.cost
}

// Java: InlineCostEstimator.CompiledSizeEstimator.continueProcessing()
func (e *InlineCostEstimator_CompiledSizeEstimator) ContinueProcessing() bool {
	return e.cg.cc.ContinueProcessing()
}

// Java: InlineCostEstimator.CompiledSizeEstimator.getLastChar()
func (e *InlineCostEstimator_CompiledSizeEstimator) GetLastChar() uint16 {
	return e.cg.cc.GetLastChar()
}
