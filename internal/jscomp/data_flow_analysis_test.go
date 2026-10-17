// Port of com/google/javascript/jscomp/DataFlowAnalysisTest.java @ closure-compiler v20260415
// (test/com/google/javascript/jscomp/DataFlowAnalysisTest.java)

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

// The computeEscaped tests (testEscaped*, computeEscapedLocals) are in
// data_flow_analysis_escaped_test.go: they need LiveVariablesAnalysis and
// GoogleCodingConvention, which are ported in a later wave.

import (
	"strconv"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// DataFlowAnalysisTest is a test suite with a very small programming
// language that has two types of instructions: BranchInstruction and
// ArithmeticInstruction. Test cases must construct a small program with
// these instructions and manually put each instruction in a
// ControlFlowGraph.
type DataFlowAnalysisTest struct{ t *testing.T }

func (c *DataFlowAnalysisTest) JUnitSetT(t *testing.T) { c.t = t }
func (c *DataFlowAnalysisTest) JUnitBefore()           {}
func (c *DataFlowAnalysisTest) JUnitAfter()            {}

func TestDataFlowAnalysisTest(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return &DataFlowAnalysisTest{} })
}

// dataFlowAnalysisTest_Operation: operations supported by
// ArithmeticInstruction.
type dataFlowAnalysisTest_Operation int32

const (
	dataFlowAnalysisTest_Operation_ADD dataFlowAnalysisTest_Operation = iota
	dataFlowAnalysisTest_Operation_SUB
	dataFlowAnalysisTest_Operation_DIV
	dataFlowAnalysisTest_Operation_MUL
)

// Java: DataFlowAnalysisTest.Operation(String) (stringRep)
var dataFlowAnalysisTest_Operation_stringRep = [...]string{"+", "-", "/", "*"}

// Java: DataFlowAnalysisTest.Operation.toString()
func (o dataFlowAnalysisTest_Operation) String() string {
	return dataFlowAnalysisTest_Operation_stringRep[o]
}

// dataFlowAnalysisTest_Value is a simple value.
type dataFlowAnalysisTest_Value interface {
	// Java: DataFlowAnalysisTest.Value.isNumber()
	IsNumber() bool
	// Java: DataFlowAnalysisTest.Value.isVariable()
	IsVariable() bool
	String() string
}

// dataFlowAnalysisTest_valueIsNumber is Value.isNumber(): `this instanceof
// NumberValue`.
func dataFlowAnalysisTest_valueIsNumber(v dataFlowAnalysisTest_Value) bool {
	_, ok := v.(*dataFlowAnalysisTest_NumberValue)
	return ok
}

// dataFlowAnalysisTest_valueIsVariable is Value.isVariable(): `this
// instanceof Variable`.
func dataFlowAnalysisTest_valueIsVariable(v dataFlowAnalysisTest_Value) bool {
	_, ok := v.(*dataFlowAnalysisTest_Variable)
	return ok
}

// dataFlowAnalysisTest_Variable is a variable.
type dataFlowAnalysisTest_Variable struct {
	name string
}

// newDataFlowAnalysisTest_Variable: n is the name of the variable.
//
// Java: DataFlowAnalysisTest.Variable(String)
func newDataFlowAnalysisTest_Variable(n string) *dataFlowAnalysisTest_Variable {
	return &dataFlowAnalysisTest_Variable{name: n}
}

func (v *dataFlowAnalysisTest_Variable) IsNumber() bool { return dataFlowAnalysisTest_valueIsNumber(v) }
func (v *dataFlowAnalysisTest_Variable) IsVariable() bool {
	return dataFlowAnalysisTest_valueIsVariable(v)
}

// Java: DataFlowAnalysisTest.Variable.getName()
func (v *dataFlowAnalysisTest_Variable) GetName() string {
	return v.name
}

// Java: DataFlowAnalysisTest.Variable.equals(Object)
func (v *dataFlowAnalysisTest_Variable) Equals(other any) bool {
	// Use the String's .equals()
	variable, ok := other.(*dataFlowAnalysisTest_Variable)
	if !ok || variable == nil {
		return false
	}
	return variable.name == v.name
}

// Java: DataFlowAnalysisTest.Variable.hashCode()
func (v *dataFlowAnalysisTest_Variable) HashCode() int32 {
	return jutil.StringHashCode(v.name)
}

// Java: DataFlowAnalysisTest.Variable.toString()
func (v *dataFlowAnalysisTest_Variable) String() string {
	return v.name
}

// dataFlowAnalysisTest_NumberValue is a number constant.
type dataFlowAnalysisTest_NumberValue struct {
	value int32
}

// newDataFlowAnalysisTest_NumberValue: v is the value.
//
// Java: DataFlowAnalysisTest.NumberValue(int)
func newDataFlowAnalysisTest_NumberValue(v int32) *dataFlowAnalysisTest_NumberValue {
	return &dataFlowAnalysisTest_NumberValue{value: v}
}

func (v *dataFlowAnalysisTest_NumberValue) IsNumber() bool {
	return dataFlowAnalysisTest_valueIsNumber(v)
}
func (v *dataFlowAnalysisTest_NumberValue) IsVariable() bool {
	return dataFlowAnalysisTest_valueIsVariable(v)
}

// Java: DataFlowAnalysisTest.NumberValue.getValue()
func (v *dataFlowAnalysisTest_NumberValue) GetValue() int32 {
	return v.value
}

// Java: DataFlowAnalysisTest.NumberValue.toString()
func (v *dataFlowAnalysisTest_NumberValue) String() string {
	return "" + strconv.FormatInt(int64(v.value), 10)
}

// Java: DataFlowAnalysisTest.NumberValue.equals(Object)
func (v *dataFlowAnalysisTest_NumberValue) Equals(other any) bool {
	numberValue, ok := other.(*dataFlowAnalysisTest_NumberValue)
	if !ok || numberValue == nil {
		return false
	}
	return numberValue.value == v.value
}

// Java: DataFlowAnalysisTest.NumberValue.hashCode()
func (v *dataFlowAnalysisTest_NumberValue) HashCode() int32 {
	return v.value
}

// dataFlowAnalysisTest_Instruction is an instruction of the dummy program.
type dataFlowAnalysisTest_Instruction interface {
	// IsArithmetic checks whether this is an arithmetic instruction.
	//
	// Java: DataFlowAnalysisTest.Instruction.isArithmetic()
	IsArithmetic() bool
	// IsBranch checks whether this is a branch instruction.
	//
	// Java: DataFlowAnalysisTest.Instruction.isBranch()
	IsBranch() bool
	GetInstructionBase() *dataFlowAnalysisTest_InstructionBase
}

// dataFlowAnalysisTest_InstructionBase holds Instruction's fields.
type dataFlowAnalysisTest_InstructionBase struct {
	order int32
}

func (b *dataFlowAnalysisTest_InstructionBase) GetInstructionBase() *dataFlowAnalysisTest_InstructionBase {
	return b
}

// dataFlowAnalysisTest_ArithmeticInstruction is a basic arithmetic
// instruction that only takes the form of:
//
//	Result = Operand1 operator Operand2
type dataFlowAnalysisTest_ArithmeticInstruction struct {
	dataFlowAnalysisTest_InstructionBase
	operation dataFlowAnalysisTest_Operation
	operand1  dataFlowAnalysisTest_Value
	operand2  dataFlowAnalysisTest_Value
	result    *dataFlowAnalysisTest_Variable
}

// Java: DataFlowAnalysisTest.ArithmeticInstruction(Variable, int, Operation, int)
func newDataFlowAnalysisTest_ArithmeticInstruction(res *dataFlowAnalysisTest_Variable, op1 int32, o dataFlowAnalysisTest_Operation, op2 int32) *dataFlowAnalysisTest_ArithmeticInstruction {
	return newDataFlowAnalysisTest_ArithmeticInstructionOfValues(res, newDataFlowAnalysisTest_NumberValue(op1), o, newDataFlowAnalysisTest_NumberValue(op2))
}

// Java: DataFlowAnalysisTest.ArithmeticInstruction(Variable, Value, Operation, int)
func newDataFlowAnalysisTest_ArithmeticInstructionOfValueAndInt(res *dataFlowAnalysisTest_Variable, op1 dataFlowAnalysisTest_Value, o dataFlowAnalysisTest_Operation, op2 int32) *dataFlowAnalysisTest_ArithmeticInstruction {
	return newDataFlowAnalysisTest_ArithmeticInstructionOfValues(res, op1, o, newDataFlowAnalysisTest_NumberValue(op2))
}

// Java: DataFlowAnalysisTest.ArithmeticInstruction(Variable, int, Operation, Value)
func newDataFlowAnalysisTest_ArithmeticInstructionOfIntAndValue(res *dataFlowAnalysisTest_Variable, op1 int32, o dataFlowAnalysisTest_Operation, op2 dataFlowAnalysisTest_Value) *dataFlowAnalysisTest_ArithmeticInstruction {
	return newDataFlowAnalysisTest_ArithmeticInstructionOfValues(res, newDataFlowAnalysisTest_NumberValue(op1), o, op2)
}

// Java: DataFlowAnalysisTest.ArithmeticInstruction(Variable, Value, Operation, Value)
func newDataFlowAnalysisTest_ArithmeticInstructionOfValues(res *dataFlowAnalysisTest_Variable, op1 dataFlowAnalysisTest_Value, o dataFlowAnalysisTest_Operation, op2 dataFlowAnalysisTest_Value) *dataFlowAnalysisTest_ArithmeticInstruction {
	return &dataFlowAnalysisTest_ArithmeticInstruction{
		result:    res,
		operand1:  op1,
		operand2:  op2,
		operation: o,
	}
}

func (i *dataFlowAnalysisTest_ArithmeticInstruction) IsArithmetic() bool { return true }
func (i *dataFlowAnalysisTest_ArithmeticInstruction) IsBranch() bool     { return false }

// Java: DataFlowAnalysisTest.ArithmeticInstruction.getOperator()
func (i *dataFlowAnalysisTest_ArithmeticInstruction) GetOperator() dataFlowAnalysisTest_Operation {
	return i.operation
}

// Java: DataFlowAnalysisTest.ArithmeticInstruction.setOperator(Operation)
func (i *dataFlowAnalysisTest_ArithmeticInstruction) SetOperator(op dataFlowAnalysisTest_Operation) {
	i.operation = op
}

// Java: DataFlowAnalysisTest.ArithmeticInstruction.getOperand1()
func (i *dataFlowAnalysisTest_ArithmeticInstruction) GetOperand1() dataFlowAnalysisTest_Value {
	return i.operand1
}

// Java: DataFlowAnalysisTest.ArithmeticInstruction.setOperand1(Value)
func (i *dataFlowAnalysisTest_ArithmeticInstruction) SetOperand1(operand1 dataFlowAnalysisTest_Value) {
	i.operand1 = operand1
}

// Java: DataFlowAnalysisTest.ArithmeticInstruction.getOperand2()
func (i *dataFlowAnalysisTest_ArithmeticInstruction) GetOperand2() dataFlowAnalysisTest_Value {
	return i.operand2
}

// Java: DataFlowAnalysisTest.ArithmeticInstruction.setOperand2(Value)
func (i *dataFlowAnalysisTest_ArithmeticInstruction) SetOperand2(operand2 dataFlowAnalysisTest_Value) {
	i.operand2 = operand2
}

// Java: DataFlowAnalysisTest.ArithmeticInstruction.getResult()
func (i *dataFlowAnalysisTest_ArithmeticInstruction) GetResult() *dataFlowAnalysisTest_Variable {
	return i.result
}

// Java: DataFlowAnalysisTest.ArithmeticInstruction.setResult(Variable)
func (i *dataFlowAnalysisTest_ArithmeticInstruction) SetResult(result *dataFlowAnalysisTest_Variable) {
	i.result = result
}

// Java: DataFlowAnalysisTest.ArithmeticInstruction.toString()
func (i *dataFlowAnalysisTest_ArithmeticInstruction) String() string {
	out := jutil.NewStringBuilder()
	out.Append(i.result.String())
	out.Append(" = ")
	out.Append(i.operand1.String())
	out.Append(i.operation.String())
	out.Append(i.operand2.String())
	return out.String()
}

// Java: DataFlowAnalysisTest.ArithmeticInstruction.equals(Object)
func (i *dataFlowAnalysisTest_ArithmeticInstruction) Equals(other any) bool {
	that, ok := other.(*dataFlowAnalysisTest_ArithmeticInstruction)
	if !ok || that == nil {
		return false
	}
	return that.order == i.order &&
		that.operation == i.operation &&
		jutil.ObjectsEquals(that.operand1, i.operand1) &&
		jutil.ObjectsEquals(that.operand2, i.operand2) &&
		that.result.Equals(i.result)
}

// Java: DataFlowAnalysisTest.ArithmeticInstruction.hashCode()
func (i *dataFlowAnalysisTest_ArithmeticInstruction) HashCode() int32 {
	return jutil.StringHashCode(i.String())
}

// Java: DataFlowAnalysisTest.newAssignNumberToVariableInstruction(Variable, int)
func DataFlowAnalysisTest_NewAssignNumberToVariableInstruction(res *dataFlowAnalysisTest_Variable, num int32) *dataFlowAnalysisTest_ArithmeticInstruction {
	return newDataFlowAnalysisTest_ArithmeticInstruction(res, num, dataFlowAnalysisTest_Operation_ADD, 0)
}

// Java: DataFlowAnalysisTest.newAssignVariableToVariableInstruction(Variable, Variable)
func DataFlowAnalysisTest_NewAssignVariableToVariableInstruction(
	lhs *dataFlowAnalysisTest_Variable, rhs *dataFlowAnalysisTest_Variable) *dataFlowAnalysisTest_ArithmeticInstruction {
	return newDataFlowAnalysisTest_ArithmeticInstructionOfValueAndInt(lhs, rhs, dataFlowAnalysisTest_Operation_ADD, 0)
}

// dataFlowAnalysisTest_BranchInstruction is a branch instruction based on a
// Value as a condition.
type dataFlowAnalysisTest_BranchInstruction struct {
	dataFlowAnalysisTest_InstructionBase
	condition dataFlowAnalysisTest_Value
}

// Java: DataFlowAnalysisTest.BranchInstruction(Value)
func newDataFlowAnalysisTest_BranchInstruction(cond dataFlowAnalysisTest_Value) *dataFlowAnalysisTest_BranchInstruction {
	return &dataFlowAnalysisTest_BranchInstruction{condition: cond}
}

func (i *dataFlowAnalysisTest_BranchInstruction) IsArithmetic() bool { return false }
func (i *dataFlowAnalysisTest_BranchInstruction) IsBranch() bool     { return true }

// Java: DataFlowAnalysisTest.BranchInstruction.getCondition()
func (i *dataFlowAnalysisTest_BranchInstruction) GetCondition() dataFlowAnalysisTest_Value {
	return i.condition
}

// Java: DataFlowAnalysisTest.BranchInstruction.setCondition(Value)
func (i *dataFlowAnalysisTest_BranchInstruction) SetCondition(condition dataFlowAnalysisTest_Value) {
	i.condition = condition
}

// dataFlowAnalysisTest_ConstPropLatticeElement is a lattice to represent
// constant states. Each variable of the program will have a lattice
// defined as:
//
//	      TOP
//	 / / |         \
//	0  1 2 3 ..... MAX_VALUE
//	\  \ |         /
//	     BOTTOM
//
// Where BOTTOM represents the variable is not a constant.
//
// This class will represent a product lattice of each variable's lattice.
// The whole lattice is store in a HashMap. If variable x is defined to be
// constant 10. The map will contain the value 10 with the variable x as key.
// Otherwise, x is not a constant.
//
// Go-only: language-adapter — the HashMap<Variable, Integer> keys use
// Variable.equals, so it is an ObjectLinkedHashMap (its iteration order only
// reaches toString).
type dataFlowAnalysisTest_ConstPropLatticeElement struct {
	constMap *jutil.ObjectLinkedHashMap[*dataFlowAnalysisTest_Variable, int32]
	isTop    bool
}

// newDataFlowAnalysisTest_ConstPropLatticeElementWithTop: isTop is to
// define if the lattice is top.
//
// Java: DataFlowAnalysisTest.ConstPropLatticeElement(boolean)
func newDataFlowAnalysisTest_ConstPropLatticeElementWithTop(isTop bool) *dataFlowAnalysisTest_ConstPropLatticeElement {
	return &dataFlowAnalysisTest_ConstPropLatticeElement{
		isTop:    isTop,
		constMap: jutil.NewObjectLinkedHashMap[*dataFlowAnalysisTest_Variable, int32](),
	}
}

// newDataFlowAnalysisTest_ConstPropLatticeElement creates a lattice where
// every variable is defined to be not constant.
//
// Java: DataFlowAnalysisTest.ConstPropLatticeElement()
func newDataFlowAnalysisTest_ConstPropLatticeElement() *dataFlowAnalysisTest_ConstPropLatticeElement {
	return newDataFlowAnalysisTest_ConstPropLatticeElementWithTop(false)
}

// Java: DataFlowAnalysisTest.ConstPropLatticeElement(ConstPropLatticeElement)
func newDataFlowAnalysisTest_ConstPropLatticeElementCopy(other *dataFlowAnalysisTest_ConstPropLatticeElement) *dataFlowAnalysisTest_ConstPropLatticeElement {
	e := &dataFlowAnalysisTest_ConstPropLatticeElement{
		isTop:    other.isTop,
		constMap: jutil.NewObjectLinkedHashMap[*dataFlowAnalysisTest_Variable, int32](),
	}
	for k, v := range other.constMap.All() {
		e.constMap.Put(k, v)
	}
	return e
}

// Java: DataFlowAnalysisTest.ConstPropLatticeElement.toString()
func (e *dataFlowAnalysisTest_ConstPropLatticeElement) String() string {
	if e.isTop {
		return "TOP"
	}
	out := jutil.NewStringBuilder()

	out.Append("{")
	for var_ := range e.constMap.Keys() {
		out.Append(var_.String())
		out.Append("=")
		out.Append(strconv.FormatInt(int64(e.constMap.Get(var_)), 10))
		out.Append(" ")
	}
	out.Append("}")
	return out.String()
}

// Java: DataFlowAnalysisTest.ConstPropLatticeElement.equals(Object)
func (e *dataFlowAnalysisTest_ConstPropLatticeElement) Equals(other any) bool {
	if otherLattice, ok := other.(*dataFlowAnalysisTest_ConstPropLatticeElement); ok && otherLattice != nil {
		return (e.isTop == otherLattice.isTop) && e.constMap.Equals(otherLattice.constMap)
	}
	return false
}

// Java: DataFlowAnalysisTest.ConstPropLatticeElement.hashCode()
func (e *dataFlowAnalysisTest_ConstPropLatticeElement) HashCode() int32 {
	return e.constMap.HashCode()
}

// dataFlowAnalysisTest_ConstPropJoinOp is ConstPropJoinOp.
type dataFlowAnalysisTest_ConstPropJoinOp struct {
	result *dataFlowAnalysisTest_ConstPropLatticeElement
}

// Java: DataFlowAnalysisTest.ConstPropJoinOp.joinFlow(ConstPropLatticeElement)
func (j *dataFlowAnalysisTest_ConstPropJoinOp) JoinFlow(input *dataFlowAnalysisTest_ConstPropLatticeElement) {
	if j.result == nil {
		j.result = input
	} else {
		j.result = dataFlowAnalysisTest_ConstPropJoinOp_apply(j.result, input)
	}
}

// Java: DataFlowAnalysisTest.ConstPropJoinOp.finish()
func (j *dataFlowAnalysisTest_ConstPropJoinOp) Finish() *dataFlowAnalysisTest_ConstPropLatticeElement {
	return j.result
}

// Java: DataFlowAnalysisTest.ConstPropJoinOp.apply(ConstPropLatticeElement, ConstPropLatticeElement)
func dataFlowAnalysisTest_ConstPropJoinOp_apply(
	a *dataFlowAnalysisTest_ConstPropLatticeElement, b *dataFlowAnalysisTest_ConstPropLatticeElement) *dataFlowAnalysisTest_ConstPropLatticeElement {
	result := newDataFlowAnalysisTest_ConstPropLatticeElement()
	// By the definition of TOP of the lattice.
	if a.isTop {
		return newDataFlowAnalysisTest_ConstPropLatticeElementCopy(a)
	}
	if b.isTop {
		return newDataFlowAnalysisTest_ConstPropLatticeElementCopy(b)
	}
	// Do the join for each variable's lattice.
	for var_ := range a.constMap.Keys() {
		if b.constMap.ContainsKey(var_) {
			number := b.constMap.Get(var_)

			// The result will contain that variable as a known constant
			// if both lattice has that variable the same constant.
			if a.constMap.Get(var_) == number {
				result.constMap.Put(var_, number)
			}
		}
	}
	return result
}

// dataFlowAnalysisTest_DummyConstPropagation is a simple forward constant
// propagation.
type dataFlowAnalysisTest_DummyConstPropagation struct {
	DataFlowAnalysisBase[dataFlowAnalysisTest_Instruction, *dataFlowAnalysisTest_ConstPropLatticeElement]
}

// newDataFlowAnalysisTest_DummyConstPropagation: targetCfg is the Control
// Flow Graph.
//
// Java: DataFlowAnalysisTest.DummyConstPropagation(ControlFlowGraph<Instruction>)
func newDataFlowAnalysisTest_DummyConstPropagation(targetCfg ControlFlowGraph[dataFlowAnalysisTest_Instruction]) *dataFlowAnalysisTest_DummyConstPropagation {
	a := &dataFlowAnalysisTest_DummyConstPropagation{}
	a.InitDataFlowAnalysisBase(a, targetCfg)
	return a
}

// Java: DataFlowAnalysisTest.DummyConstPropagation.isForward()
func (a *dataFlowAnalysisTest_DummyConstPropagation) IsForward() bool {
	return true
}

// Java: DataFlowAnalysisTest.DummyConstPropagation.createFlowJoiner()
func (a *dataFlowAnalysisTest_DummyConstPropagation) CreateFlowJoiner() DataFlowAnalysis_FlowJoiner[*dataFlowAnalysisTest_ConstPropLatticeElement] {
	return &dataFlowAnalysisTest_ConstPropJoinOp{}
}

// Java: DataFlowAnalysisTest.DummyConstPropagation.flowThrough(Instruction, ConstPropLatticeElement)
func (a *dataFlowAnalysisTest_DummyConstPropagation) FlowThrough(node dataFlowAnalysisTest_Instruction, input *dataFlowAnalysisTest_ConstPropLatticeElement) *dataFlowAnalysisTest_ConstPropLatticeElement {
	if node.IsBranch() {
		return newDataFlowAnalysisTest_ConstPropLatticeElementCopy(input)
	} else {
		return dataFlowAnalysisTest_flowThroughArithmeticInstruction(node.(*dataFlowAnalysisTest_ArithmeticInstruction), input)
	}
}

// Java: DataFlowAnalysisTest.DummyConstPropagation.createEntryLattice()
func (a *dataFlowAnalysisTest_DummyConstPropagation) CreateEntryLattice() *dataFlowAnalysisTest_ConstPropLatticeElement {
	return newDataFlowAnalysisTest_ConstPropLatticeElement()
}

// Java: DataFlowAnalysisTest.DummyConstPropagation.createInitialEstimateLattice()
func (a *dataFlowAnalysisTest_DummyConstPropagation) CreateInitialEstimateLattice() *dataFlowAnalysisTest_ConstPropLatticeElement {
	return newDataFlowAnalysisTest_ConstPropLatticeElementWithTop(true)
}

// Java: DataFlowAnalysisTest.flowThroughArithmeticInstruction(ArithmeticInstruction, ConstPropLatticeElement)
func dataFlowAnalysisTest_flowThroughArithmeticInstruction(
	aInst *dataFlowAnalysisTest_ArithmeticInstruction, input *dataFlowAnalysisTest_ConstPropLatticeElement) *dataFlowAnalysisTest_ConstPropLatticeElement {

	out := newDataFlowAnalysisTest_ConstPropLatticeElementCopy(input)
	// Try to see if left is a number. If it is a variable, it might already
	// be a constant coming in.
	var leftConst *int32
	if aInst.operand1.IsNumber() {
		v := aInst.operand1.(*dataFlowAnalysisTest_NumberValue).value
		leftConst = &v
	} else {
		if variable, ok := aInst.operand1.(*dataFlowAnalysisTest_Variable); ok && input.constMap.ContainsKey(variable) {
			v := input.constMap.Get(variable)
			leftConst = &v
		}
	}

	// Do the same thing to the right.
	var rightConst *int32
	if aInst.operand2.IsNumber() {
		v := aInst.operand2.(*dataFlowAnalysisTest_NumberValue).value
		rightConst = &v
	} else {
		if variable, ok := aInst.operand2.(*dataFlowAnalysisTest_Variable); ok && input.constMap.ContainsKey(variable) {
			v := input.constMap.Get(variable)
			rightConst = &v
		}
	}

	// If both are known constant we can perform the operation.
	if leftConst != nil && rightConst != nil {
		var constResult int32
		switch aInst.operation {
		case dataFlowAnalysisTest_Operation_ADD:
			constResult = *leftConst + *rightConst
		case dataFlowAnalysisTest_Operation_SUB:
			constResult = *leftConst - *rightConst
		case dataFlowAnalysisTest_Operation_MUL:
			constResult = *leftConst * *rightConst
		case dataFlowAnalysisTest_Operation_DIV:
			if *rightConst == 0 {
				panic(jutil.NewArithmeticException("/ by zero"))
			}
			constResult = *leftConst / *rightConst
		}
		// Put it in the map. (Possibly replacing the existing constant value)
		out.constMap.Put(aInst.result, constResult)
	} else {
		// If we cannot find a constant for it
		out.constMap.Remove(aInst.result)
	}
	return out
}

// Java: DataFlowAnalysisTest.testSimpleIf()
func (c *DataFlowAnalysisTest) TestSimpleIf() {
	// if (a) { b = 1; } else { b = 1; } c = b;
	a := newDataFlowAnalysisTest_Variable("a")
	b := newDataFlowAnalysisTest_Variable("b")
	cv := newDataFlowAnalysisTest_Variable("c")
	var inst1 dataFlowAnalysisTest_Instruction = newDataFlowAnalysisTest_BranchInstruction(a)
	var inst2 dataFlowAnalysisTest_Instruction = DataFlowAnalysisTest_NewAssignNumberToVariableInstruction(b, 1)
	var inst3 dataFlowAnalysisTest_Instruction = DataFlowAnalysisTest_NewAssignNumberToVariableInstruction(b, 1)
	var inst4 dataFlowAnalysisTest_Instruction = DataFlowAnalysisTest_NewAssignVariableToVariableInstruction(cv, b)
	cfg := NewControlFlowGraph(inst1, true, true)
	n1 := cfg.CreateNode(inst1)
	n2 := cfg.CreateNode(inst2)
	n3 := cfg.CreateNode(inst3)
	n4 := cfg.CreateNode(inst4)
	cfg.Connect(inst1, ControlFlowGraph_Branch_ON_FALSE, inst2)
	cfg.Connect(inst1, ControlFlowGraph_Branch_ON_TRUE, inst3)
	cfg.Connect(inst2, ControlFlowGraph_Branch_UNCOND, inst4)
	cfg.Connect(inst3, ControlFlowGraph_Branch_UNCOND, inst4)

	constProp := newDataFlowAnalysisTest_DummyConstPropagation(cfg)
	constProp.Analyze()

	// We cannot conclude anything from if (a).
	dataFlowAnalysisTest_verifyInHas(n1, a, nil)
	dataFlowAnalysisTest_verifyInHas(n1, b, nil)
	dataFlowAnalysisTest_verifyInHas(n1, cv, nil)
	dataFlowAnalysisTest_verifyOutHas(n1, a, nil)
	dataFlowAnalysisTest_verifyOutHas(n1, b, nil)
	dataFlowAnalysisTest_verifyOutHas(n1, cv, nil)

	// We can conclude b = 1 after the instruction.
	dataFlowAnalysisTest_verifyInHas(n2, a, nil)
	dataFlowAnalysisTest_verifyInHas(n2, b, nil)
	dataFlowAnalysisTest_verifyInHas(n2, cv, nil)
	dataFlowAnalysisTest_verifyOutHas(n2, a, nil)
	dataFlowAnalysisTest_verifyOutHas(n2, b, dataFlowAnalysisTest_int(1))
	dataFlowAnalysisTest_verifyOutHas(n2, cv, nil)

	// Same as above.
	dataFlowAnalysisTest_verifyInHas(n3, a, nil)
	dataFlowAnalysisTest_verifyInHas(n3, b, nil)
	dataFlowAnalysisTest_verifyInHas(n3, cv, nil)
	dataFlowAnalysisTest_verifyOutHas(n3, a, nil)
	dataFlowAnalysisTest_verifyOutHas(n3, b, dataFlowAnalysisTest_int(1))
	dataFlowAnalysisTest_verifyOutHas(n3, cv, nil)

	// After the merge we should still have b = 1.
	dataFlowAnalysisTest_verifyInHas(n4, a, nil)
	dataFlowAnalysisTest_verifyInHas(n4, b, dataFlowAnalysisTest_int(1))
	dataFlowAnalysisTest_verifyInHas(n4, cv, nil)
	dataFlowAnalysisTest_verifyOutHas(n4, a, nil)
	// After the instruction both b and c are 1.
	dataFlowAnalysisTest_verifyOutHas(n4, b, dataFlowAnalysisTest_int(1))
	dataFlowAnalysisTest_verifyOutHas(n4, cv, dataFlowAnalysisTest_int(1))
}

// Java: DataFlowAnalysisTest.testSimpleLoop()
func (c *DataFlowAnalysisTest) TestSimpleLoop() {
	// a = 0; do { a = a + 1 } while (b); c = a;
	a := newDataFlowAnalysisTest_Variable("a")
	b := newDataFlowAnalysisTest_Variable("b")
	cv := newDataFlowAnalysisTest_Variable("c")
	var inst1 dataFlowAnalysisTest_Instruction = DataFlowAnalysisTest_NewAssignNumberToVariableInstruction(a, 0)
	var inst2 dataFlowAnalysisTest_Instruction = newDataFlowAnalysisTest_ArithmeticInstructionOfValueAndInt(a, a, dataFlowAnalysisTest_Operation_ADD, 1)
	var inst3 dataFlowAnalysisTest_Instruction = newDataFlowAnalysisTest_BranchInstruction(b)
	var inst4 dataFlowAnalysisTest_Instruction = DataFlowAnalysisTest_NewAssignVariableToVariableInstruction(cv, a)
	cfg := NewControlFlowGraph(inst1, true, true)
	n1 := cfg.CreateNode(inst1)
	n2 := cfg.CreateNode(inst2)
	n3 := cfg.CreateNode(inst3)
	n4 := cfg.CreateNode(inst4)
	cfg.Connect(inst1, ControlFlowGraph_Branch_UNCOND, inst2)
	cfg.Connect(inst2, ControlFlowGraph_Branch_UNCOND, inst3)
	cfg.Connect(inst3, ControlFlowGraph_Branch_ON_TRUE, inst2)
	cfg.Connect(inst3, ControlFlowGraph_Branch_ON_FALSE, inst4)

	constProp := newDataFlowAnalysisTest_DummyConstPropagation(cfg)
	// This will also show that the framework terminates properly.
	constProp.Analyze()

	// a = 0 is the only thing we know.
	dataFlowAnalysisTest_verifyInHas(n1, a, nil)
	dataFlowAnalysisTest_verifyInHas(n1, b, nil)
	dataFlowAnalysisTest_verifyInHas(n1, cv, nil)
	dataFlowAnalysisTest_verifyOutHas(n1, a, dataFlowAnalysisTest_int(0))
	dataFlowAnalysisTest_verifyOutHas(n1, b, nil)
	dataFlowAnalysisTest_verifyOutHas(n1, cv, nil)

	// Nothing is provable in this program, so confirm that we haven't
	// erroneously "proven" something.
	dataFlowAnalysisTest_verifyInHas(n2, a, nil)
	dataFlowAnalysisTest_verifyInHas(n2, b, nil)
	dataFlowAnalysisTest_verifyInHas(n2, cv, nil)
	dataFlowAnalysisTest_verifyOutHas(n2, a, nil)
	dataFlowAnalysisTest_verifyOutHas(n2, b, nil)
	dataFlowAnalysisTest_verifyOutHas(n2, cv, nil)

	dataFlowAnalysisTest_verifyInHas(n3, a, nil)
	dataFlowAnalysisTest_verifyInHas(n3, b, nil)
	dataFlowAnalysisTest_verifyInHas(n3, cv, nil)
	dataFlowAnalysisTest_verifyOutHas(n3, a, nil)
	dataFlowAnalysisTest_verifyOutHas(n3, b, nil)
	dataFlowAnalysisTest_verifyOutHas(n3, cv, nil)

	dataFlowAnalysisTest_verifyInHas(n4, a, nil)
	dataFlowAnalysisTest_verifyInHas(n4, b, nil)
	dataFlowAnalysisTest_verifyInHas(n4, cv, nil)
	dataFlowAnalysisTest_verifyOutHas(n4, a, nil)
	dataFlowAnalysisTest_verifyOutHas(n4, b, nil)
	dataFlowAnalysisTest_verifyOutHas(n4, cv, nil)
}

// dataFlowAnalysisTest_BranchedDummyConstPropagation is a simple forward
// constant propagation.
type dataFlowAnalysisTest_BranchedDummyConstPropagation struct {
	DataFlowAnalysisBase[dataFlowAnalysisTest_Instruction, *dataFlowAnalysisTest_ConstPropLatticeElement]
}

// Java: DataFlowAnalysisTest.BranchedDummyConstPropagation(ControlFlowGraph<Instruction>)
func newDataFlowAnalysisTest_BranchedDummyConstPropagation(targetCfg ControlFlowGraph[dataFlowAnalysisTest_Instruction]) *dataFlowAnalysisTest_BranchedDummyConstPropagation {
	a := &dataFlowAnalysisTest_BranchedDummyConstPropagation{}
	a.InitDataFlowAnalysisBase(a, targetCfg)
	return a
}

// Java: DataFlowAnalysisTest.BranchedDummyConstPropagation.isForward()
func (a *dataFlowAnalysisTest_BranchedDummyConstPropagation) IsForward() bool {
	return true
}

// Java: DataFlowAnalysisTest.BranchedDummyConstPropagation.isBranched()
func (a *dataFlowAnalysisTest_BranchedDummyConstPropagation) IsBranched() bool {
	return true
}

// Java: DataFlowAnalysisTest.BranchedDummyConstPropagation.flowThrough(Instruction, ConstPropLatticeElement)
func (a *dataFlowAnalysisTest_BranchedDummyConstPropagation) FlowThrough(node dataFlowAnalysisTest_Instruction, input *dataFlowAnalysisTest_ConstPropLatticeElement) *dataFlowAnalysisTest_ConstPropLatticeElement {
	if node.IsArithmetic() {
		return dataFlowAnalysisTest_flowThroughArithmeticInstruction(node.(*dataFlowAnalysisTest_ArithmeticInstruction), input)
	} else {
		return newDataFlowAnalysisTest_ConstPropLatticeElementCopy(input)
	}
}

// Java: DataFlowAnalysisTest.BranchedDummyConstPropagation.createFlowBrancher(Instruction, ConstPropLatticeElement)
func (a *dataFlowAnalysisTest_BranchedDummyConstPropagation) CreateFlowBrancher(
	node dataFlowAnalysisTest_Instruction, input *dataFlowAnalysisTest_ConstPropLatticeElement) DataFlowAnalysis_FlowBrancher[*dataFlowAnalysisTest_ConstPropLatticeElement] {
	var aResult *dataFlowAnalysisTest_ConstPropLatticeElement
	if node.IsArithmetic() {
		aResult = dataFlowAnalysisTest_flowThroughArithmeticInstruction(node.(*dataFlowAnalysisTest_ArithmeticInstruction), input)
	}

	return func(branch ControlFlowGraph_Branch) *dataFlowAnalysisTest_ConstPropLatticeElement {
		if aResult != nil {
			return aResult
		}

		branchInst := node.(*dataFlowAnalysisTest_BranchInstruction)
		edgeResult := newDataFlowAnalysisTest_ConstPropLatticeElementCopy(input)
		if branch == ControlFlowGraph_Branch_ON_FALSE && branchInst.GetCondition().IsVariable() {
			edgeResult.constMap.Put(branchInst.GetCondition().(*dataFlowAnalysisTest_Variable), 0)
		}
		return edgeResult
	}
}

// Java: DataFlowAnalysisTest.BranchedDummyConstPropagation.createEntryLattice()
func (a *dataFlowAnalysisTest_BranchedDummyConstPropagation) CreateEntryLattice() *dataFlowAnalysisTest_ConstPropLatticeElement {
	return newDataFlowAnalysisTest_ConstPropLatticeElement()
}

// Java: DataFlowAnalysisTest.BranchedDummyConstPropagation.createInitialEstimateLattice()
func (a *dataFlowAnalysisTest_BranchedDummyConstPropagation) CreateInitialEstimateLattice() *dataFlowAnalysisTest_ConstPropLatticeElement {
	return newDataFlowAnalysisTest_ConstPropLatticeElementWithTop(true)
}

// Java: DataFlowAnalysisTest.BranchedDummyConstPropagation.createFlowJoiner()
func (a *dataFlowAnalysisTest_BranchedDummyConstPropagation) CreateFlowJoiner() DataFlowAnalysis_FlowJoiner[*dataFlowAnalysisTest_ConstPropLatticeElement] {
	return &dataFlowAnalysisTest_ConstPropJoinOp{}
}

// Java: DataFlowAnalysisTest.testBranchedSimpleIf()
func (c *DataFlowAnalysisTest) TestBranchedSimpleIf() {
	// if (a) { a = 0; } else { b = 0; } c = b;
	a := newDataFlowAnalysisTest_Variable("a")
	b := newDataFlowAnalysisTest_Variable("b")
	cv := newDataFlowAnalysisTest_Variable("c")
	var inst1 dataFlowAnalysisTest_Instruction = newDataFlowAnalysisTest_BranchInstruction(a)
	var inst2 dataFlowAnalysisTest_Instruction = DataFlowAnalysisTest_NewAssignNumberToVariableInstruction(a, 0)
	var inst3 dataFlowAnalysisTest_Instruction = DataFlowAnalysisTest_NewAssignNumberToVariableInstruction(b, 0)
	var inst4 dataFlowAnalysisTest_Instruction = DataFlowAnalysisTest_NewAssignVariableToVariableInstruction(cv, b)
	cfg := NewControlFlowGraph(inst1, true, true)
	n1 := cfg.CreateNode(inst1)
	n2 := cfg.CreateNode(inst2)
	n3 := cfg.CreateNode(inst3)
	n4 := cfg.CreateNode(inst4)
	cfg.Connect(inst1, ControlFlowGraph_Branch_ON_TRUE, inst2)
	cfg.Connect(inst1, ControlFlowGraph_Branch_ON_FALSE, inst3)
	cfg.Connect(inst2, ControlFlowGraph_Branch_UNCOND, inst4)
	cfg.Connect(inst3, ControlFlowGraph_Branch_UNCOND, inst4)

	constProp := newDataFlowAnalysisTest_BranchedDummyConstPropagation(cfg)
	constProp.Analyze()

	// We cannot conclude anything from if (a).
	dataFlowAnalysisTest_verifyInHas(n1, a, nil)
	dataFlowAnalysisTest_verifyInHas(n1, b, nil)
	dataFlowAnalysisTest_verifyInHas(n1, cv, nil)

	// Nothing is known on the true branch.
	dataFlowAnalysisTest_verifyInHas(n2, a, nil)
	dataFlowAnalysisTest_verifyInHas(n2, b, nil)
	dataFlowAnalysisTest_verifyInHas(n2, cv, nil)

	// Verify that we have a = 0 on the false branch.
	dataFlowAnalysisTest_verifyInHas(n3, a, dataFlowAnalysisTest_int(0))
	dataFlowAnalysisTest_verifyInHas(n3, b, nil)
	dataFlowAnalysisTest_verifyInHas(n3, cv, nil)

	// After the merge we should still have a = 0.
	dataFlowAnalysisTest_verifyInHas(n4, a, dataFlowAnalysisTest_int(0))
}

// dataFlowAnalysisTest_DivergentAnalysis_Counter is DivergentAnalysis.Counter.
type dataFlowAnalysisTest_DivergentAnalysis_Counter struct {
	count int32
}

// dataFlowAnalysisTest_DivergentAnalysis_Step is DivergentAnalysis.Step.
type dataFlowAnalysisTest_DivergentAnalysis_Step struct {
	// Go-only: language-adapter — distinct identity.
	_ byte
}

// Java: DataFlowAnalysisTest.DivergentAnalysis.Step()
func newDataFlowAnalysisTest_DivergentAnalysis_Step() *dataFlowAnalysisTest_DivergentAnalysis_Step {
	return &dataFlowAnalysisTest_DivergentAnalysis_Step{}
}

// Java: DataFlowAnalysisTest.DivergentAnalysis.Step.hashCode()
func (s *dataFlowAnalysisTest_DivergentAnalysis_Step) HashCode() int32 {
	return 0
}

// Java: DataFlowAnalysisTest.DivergentAnalysis.Step.equals(Object)
func (s *dataFlowAnalysisTest_DivergentAnalysis_Step) Equals(o any) bool {
	return false
}

// dataFlowAnalysisTest_DivergentAnalysis is DivergentAnalysis.
type dataFlowAnalysisTest_DivergentAnalysis struct {
	DataFlowAnalysisBase[*dataFlowAnalysisTest_DivergentAnalysis_Counter, *dataFlowAnalysisTest_DivergentAnalysis_Step]
}

// newDataFlowAnalysisTest_DivergentAnalysis: targetCfg is the Control Flow
// Graph.
//
// Java: DataFlowAnalysisTest.DivergentAnalysis(ControlFlowGraph<Counter>)
func newDataFlowAnalysisTest_DivergentAnalysis(targetCfg ControlFlowGraph[*dataFlowAnalysisTest_DivergentAnalysis_Counter]) *dataFlowAnalysisTest_DivergentAnalysis {
	a := &dataFlowAnalysisTest_DivergentAnalysis{}
	a.InitDataFlowAnalysisBase(a, targetCfg)
	return a
}

// Java: DataFlowAnalysisTest.DivergentAnalysis.isForward()
func (a *dataFlowAnalysisTest_DivergentAnalysis) IsForward() bool {
	return true
}

// Java: DataFlowAnalysisTest.DivergentAnalysis.flowThrough(Counter, Step)
func (a *dataFlowAnalysisTest_DivergentAnalysis) FlowThrough(node *dataFlowAnalysisTest_DivergentAnalysis_Counter, input *dataFlowAnalysisTest_DivergentAnalysis_Step) *dataFlowAnalysisTest_DivergentAnalysis_Step {
	node.count++
	return input
}

// Java: DataFlowAnalysisTest.DivergentAnalysis.createEntryLattice()
func (a *dataFlowAnalysisTest_DivergentAnalysis) CreateEntryLattice() *dataFlowAnalysisTest_DivergentAnalysis_Step {
	return newDataFlowAnalysisTest_DivergentAnalysis_Step()
}

// Java: DataFlowAnalysisTest.DivergentAnalysis.createInitialEstimateLattice()
func (a *dataFlowAnalysisTest_DivergentAnalysis) CreateInitialEstimateLattice() *dataFlowAnalysisTest_DivergentAnalysis_Step {
	return newDataFlowAnalysisTest_DivergentAnalysis_Step()
}

// dataFlowAnalysisTest_DivergentAnalysis_joiner is the anonymous FlowJoiner
// of DivergentAnalysis.createFlowJoiner().
type dataFlowAnalysisTest_DivergentAnalysis_joiner struct {
	// Go-only: language-adapter — distinct identity.
	_ byte
}

func (j *dataFlowAnalysisTest_DivergentAnalysis_joiner) JoinFlow(x *dataFlowAnalysisTest_DivergentAnalysis_Step) {
}

func (j *dataFlowAnalysisTest_DivergentAnalysis_joiner) Finish() *dataFlowAnalysisTest_DivergentAnalysis_Step {
	return newDataFlowAnalysisTest_DivergentAnalysis_Step()
}

// Java: DataFlowAnalysisTest.DivergentAnalysis.createFlowJoiner()
func (a *dataFlowAnalysisTest_DivergentAnalysis) CreateFlowJoiner() DataFlowAnalysis_FlowJoiner[*dataFlowAnalysisTest_DivergentAnalysis_Step] {
	return &dataFlowAnalysisTest_DivergentAnalysis_joiner{}
}

// Java: DataFlowAnalysisTest.testMaxIterationsExceededException()
func (c *DataFlowAnalysisTest) TestMaxIterationsExceededException() {
	entrypoint := &dataFlowAnalysisTest_DivergentAnalysis_Counter{}
	cfg := NewControlFlowGraph(entrypoint, true, true)

	cfg.Connect(entrypoint, ControlFlowGraph_Branch_UNCOND, entrypoint)
	constProp := newDataFlowAnalysisTest_DivergentAnalysis(cfg)

	// assertThrows(Exception.class, ...): every exception the analysis
	// throws is a RuntimeException.
	e := AssertThrows[jutil.IRuntimeException](constProp.Analyze)
	AssertThat(entrypoint.count).IsEqualTo(int32(DataFlowAnalysis_MAX_STEPS_PER_NODE + 1))
	AssertThat(e.GetMessage()).StartsWith("Dataflow analysis appears to diverge around: ")
}

// dataFlowAnalysisTest_int boxes an Integer constant.
func dataFlowAnalysisTest_int(v int32) *int32 { return &v }

// Java: DataFlowAnalysisTest.verifyInHas(GraphNode<Instruction, Branch>, Variable, Integer)
func dataFlowAnalysisTest_verifyInHas(
	node GraphNode[dataFlowAnalysisTest_Instruction, ControlFlowGraph_Branch], var_ *dataFlowAnalysisTest_Variable, constant *int32) {
	fState := node.GetAnnotation().(*DataFlowAnalysis_LinearFlowState[*dataFlowAnalysisTest_ConstPropLatticeElement])
	dataFlowAnalysisTest_veritfyLatticeElementHas(fState.GetIn(), var_, constant)
}

// Java: DataFlowAnalysisTest.verifyOutHas(GraphNode<Instruction, Branch>, Variable, Integer)
func dataFlowAnalysisTest_verifyOutHas(
	node GraphNode[dataFlowAnalysisTest_Instruction, ControlFlowGraph_Branch], var_ *dataFlowAnalysisTest_Variable, constant *int32) {
	fState := node.GetAnnotation().(*DataFlowAnalysis_LinearFlowState[*dataFlowAnalysisTest_ConstPropLatticeElement])
	dataFlowAnalysisTest_veritfyLatticeElementHas(fState.GetOut(), var_, constant)
}

// Java: DataFlowAnalysisTest.veritfyLatticeElementHas(ConstPropLatticeElement, Variable, Integer)
func dataFlowAnalysisTest_veritfyLatticeElementHas(el *dataFlowAnalysisTest_ConstPropLatticeElement, var_ *dataFlowAnalysisTest_Variable, constant *int32) {
	if constant == nil {
		AssertThat(el.constMap).DoesNotContainKey(var_)
	} else {
		AssertThat(el.constMap).ContainsEntry(var_, *constant)
	}
}
