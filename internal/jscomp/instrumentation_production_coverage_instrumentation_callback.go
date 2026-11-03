// Port of com/google/javascript/jscomp/instrumentation/ProductionCoverageInstrumentationCallback.java @ closure-compiler v20261028

/*
 * Copyright 2020 The Closure Compiler Authors.
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
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/sourcemap"
)

// ProductionCoverageInstrumentationCallback instruments production coverage
// for javascript. This type of coverage will instrument different levels of
// source code such as function and branch instrumentation. This
// instrumentation differs from the previous implementations
// (CoverageInstrumentationCallback and BranchCoverageInstrumentationCallback)
// in that it is properly optimized and obfuscated so that it can be run on
// client browsers with the goal of better detecting dead code. The callback
// will instrument by pushing a string onto an array which identifies what
// piece of code was executed.
//
// Java: final class ProductionCoverageInstrumentationCallback implements NodeTraversal.Callback
type ProductionCoverageInstrumentationCallback struct {
	// instrumentationArrayName is the name of the global array to which at
	// every instrumentation point a new encoded param will be added. This is
	// dynamically set by the command line flag --prod_instr_array_name.
	instrumentationArrayName string

	compiler         AbstractCompiler
	parameterMapping *productionCoverageInstrumentationCallback_ParameterMapping

	// functionNameStack stores a stack of function names that encapsulates
	// the children nodes being instrumented. The function name is popped off
	// the stack when the function node, and the entire subtree rooted at the
	// function node have been visited.
	//
	// Go-only: language-adapter — the ArrayDeque used as a stack is a slice
	// whose last element is the top.
	functionNameStack []string
}

const productionCoverageInstrumentationCallback_ANONYMOUS_FUNCTION_NAME = "<Anonymous>"

// Java: private enum ProductionCoverageInstrumentationCallback.Type
type productionCoverageInstrumentationCallback_Type int32

const (
	productionCoverageInstrumentationCallback_Type_FUNCTION productionCoverageInstrumentationCallback_Type = iota
	productionCoverageInstrumentationCallback_Type_BRANCH
	productionCoverageInstrumentationCallback_Type_BRANCH_DEFAULT
)

var productionCoverageInstrumentationCallback_Type_names = [...]string{"FUNCTION", "BRANCH", "BRANCH_DEFAULT"}

// Name is Enum.name().
func (t productionCoverageInstrumentationCallback_Type) Name() string {
	return productionCoverageInstrumentationCallback_Type_names[t]
}

// Java: ProductionCoverageInstrumentationCallback(AbstractCompiler, String)
func NewProductionCoverageInstrumentationCallback(
	compiler AbstractCompiler, instrumentationArrayName string) *ProductionCoverageInstrumentationCallback {
	return &ProductionCoverageInstrumentationCallback{
		compiler:                 compiler,
		parameterMapping:         &productionCoverageInstrumentationCallback_ParameterMapping{},
		instrumentationArrayName: instrumentationArrayName,
	}
}

// functionNamePeek returns the top of the function-name stack, or an absent nullableString.
func (c *ProductionCoverageInstrumentationCallback) functionNamePeek() nullableString {
	if len(c.functionNameStack) == 0 {
		return nullableString{}
	}
	return nullableString{val: c.functionNameStack[len(c.functionNameStack)-1], ok: true}
}

// Java: ProductionCoverageInstrumentationCallback.shouldTraverse(NodeTraversal, Node, Node)
func (c *ProductionCoverageInstrumentationCallback) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	// If origin of node is not from sourceFile, do not instrument. This typically occurs when
	// polyfill code is injected into the sourceFile AST and this check avoids instrumenting it. We
	// avoid instrumentation as this callback does not distinguish between sourceFile code and
	// injected code and can result in an error.
	if !n.IsRoot() {
		src, hasSrc := t.SourceNameOk()
		if (nullableString{val: src, ok: hasSrc}) != toNullableString(n.GetSourceFileName()) {
			return false
		}
	}

	if n.IsFunction() {
		fnName := NodeUtil_GetBestLValueName(NodeUtil_GetBestLValue(n))
		if fnName == nil {
			c.functionNameStack = append(c.functionNameStack, productionCoverageInstrumentationCallback_ANONYMOUS_FUNCTION_NAME)
		} else {
			c.functionNameStack = append(c.functionNameStack, *fnName)
		}
	}

	if len(c.functionNameStack) == 0 {
		c.functionNameStack = append(c.functionNameStack, productionCoverageInstrumentationCallback_ANONYMOUS_FUNCTION_NAME)
	}

	return true
}

// Java: ProductionCoverageInstrumentationCallback.visit(NodeTraversal, Node, Node)
func (c *ProductionCoverageInstrumentationCallback) Visit(traversal *NodeTraversal, node rhino.Node, parent rhino.Node) {
	switch node.GetToken() {
	case rhino.Token_FUNCTION:
		functionName := c.functionNamePeek()
		// If the function node has been visited by visit() then we can be assured that all its
		// children nodes have been visited and properly instrumented.
		c.functionNameStack = c.functionNameStack[:len(c.functionNameStack)-1]
		src, hasSrc := traversal.SourceNameOk()
		fileName := nullableString{val: src, ok: hasSrc}
		c.instrumentBlockNode(ast.MustFunction(node).Body().Node(), fileName, functionName, productionCoverageInstrumentationCallback_Type_FUNCTION)
	case rhino.Token_IF:
		sourceFileName := toNullableString(node.GetSourceFileName())
		functionName := c.functionNamePeek()
		ifNode := ast.MustIf(node)
		c.instrumentBlockNode(ifNode.Then().Node(), sourceFileName, functionName, productionCoverageInstrumentationCallback_Type_BRANCH)
		ifFalseBlock, hasElse := ifNode.Else()
		if !hasElse {
			ifFalseBlock = c.addElseBlock(ifNode)
		}
		// The compiler converts all sets of If-Else if-Else blocks into a combination of nested
		// If-Else blocks. This case checks if the else blocks first child is an If statement, and
		// if it is it will not instrument. This avoids adding multiple instrumentation calls.
		// Since we also make sure an Else case is added to every If statement, we are still
		// assured that the else statement is being reached through a later instrumentation call.
		firstStmt, hasStmts := ifFalseBlock.FirstStatement()
		if NodeUtil_IsEmptyBlock(ifFalseBlock.Node()) || (hasStmts && !firstStmt.Node().IsIf()) {
			c.instrumentBlockNode(ifFalseBlock.Node(), sourceFileName, functionName, productionCoverageInstrumentationCallback_Type_BRANCH_DEFAULT)
		}
	case rhino.Token_SWITCH:
		sourceFileName := toNullableString(node.GetSourceFileName())
		functionName := c.functionNamePeek()
		hasDefaultCase := false
		switchBody := ast.MustSwitch(node).Body()
		for sc := range switchBody.Cases() {
			if _, isDefault := sc.AsDefaultCase(); isDefault {
				c.instrumentBlockNode(
					sc.Body().Node(), sourceFileName, functionName, productionCoverageInstrumentationCallback_Type_BRANCH_DEFAULT)
				hasDefaultCase = true
			} else {
				c.instrumentBlockNode(sc.Body().Node(), sourceFileName, functionName, productionCoverageInstrumentationCallback_Type_BRANCH)
			}
		}
		if !hasDefaultCase {
			defaultBlock := rhino.IR_Block().SrcrefTreeIfMissing(node)
			defaultCase := rhino.IR_DefaultCase(defaultBlock).SrcrefTreeIfMissing(node)
			switchBody.AppendCase(ast.MustSwitchCase(defaultCase))
			c.instrumentBlockNode(defaultBlock, sourceFileName, functionName, productionCoverageInstrumentationCallback_Type_BRANCH_DEFAULT)
		}
	case rhino.Token_HOOK:
		sourceFileName := toNullableString(node.GetSourceFileName())
		functionName := c.functionNamePeek()
		hook := ast.MustHook(node)

		c.addInstrumentationNodeWithComma(
			hook.Then().Node(), sourceFileName, functionName, productionCoverageInstrumentationCallback_Type_BRANCH)
		c.addInstrumentationNodeWithComma(
			hook.Else().Node(), sourceFileName, functionName, productionCoverageInstrumentationCallback_Type_BRANCH)

		c.compiler.ReportChangeToEnclosingScope(node)
	case rhino.Token_OR, rhino.Token_AND, rhino.Token_COALESCE:
		sourceFileName := toNullableString(node.GetSourceFileName())
		functionName := c.functionNamePeek()
		// Only instrument the second child of the binary operation because the first child will
		// always execute, or the first child is part of a chain of binary operations and would have
		// already been instrumented.
		secondExpression := ast.MustBinary(node).Right().Node()
		c.addInstrumentationNodeWithComma(
			secondExpression, sourceFileName, functionName, productionCoverageInstrumentationCallback_Type_BRANCH)

		c.compiler.ReportChangeToEnclosingScope(node)
	default:
		if NodeUtil_IsLoopStructure(node) {
			sourceFileName := toNullableString(node.GetSourceFileName())
			functionName := c.functionNamePeek()
			blockNode := NodeUtil_GetLoopCodeBlock(node)
			rhino.CheckNotNull(blockNode)
			c.instrumentBlockNode(blockNode, sourceFileName, functionName, productionCoverageInstrumentationCallback_Type_BRANCH)
		}
	}
}

// addInstrumentationNodeWithComma: given a node, this function will create a
// new instrumentationNode and combine it with the original node using a
// COMMA node.
//
// Java: ProductionCoverageInstrumentationCallback.addInstrumentationNodeWithComma(Node, String, String, Type)
func (c *ProductionCoverageInstrumentationCallback) addInstrumentationNodeWithComma(
	originalNode rhino.Node,
	fileName nullableString,
	functionName nullableString,
	type_ productionCoverageInstrumentationCallback_Type) {
	cloneOfOriginal := originalNode.CloneTree()
	newInstrumentationNode :=
		c.newInstrumentationNode(cloneOfOriginal, fileName, functionName, type_)

	childOfInstrumentationNode := newInstrumentationNode.RemoveFirstChild()
	infusedExp := AstManipulations_FuseExpressions(childOfInstrumentationNode, cloneOfOriginal)
	originalNode.ReplaceWith(infusedExp)
}

// instrumentBlockNode consumes a block node and adds a new child to the
// front of the block node which is the instrumentation Node
//
// block is the block node to be instrumented; fileName the file name of the
// node being instrumented; fnName the function name of the node being
// instrumented; type_ the type of the node being instrumented.
//
// Java: ProductionCoverageInstrumentationCallback.instrumentBlockNode(Node, String, String, Type)
func (c *ProductionCoverageInstrumentationCallback) instrumentBlockNode(
	block rhino.Node,
	fileName nullableString,
	fnName nullableString,
	type_ productionCoverageInstrumentationCallback_Type) {
	block.AddChildToFront(c.newInstrumentationNode(block, fileName, fnName, type_))
	c.compiler.ReportChangeToEnclosingScope(block)
}

// newInstrumentationNode creates a function call to the Instrument Code
// function with properly encoded parameters. The instrumented function call
// will be of the following form: instrumentationArrayName.push(param).
// Where instrumentationArrayName is the name of the global array and param
// is the encoded param which will be pushed onto the array.
//
// node is the node to be instrumented; fileName the file name of the node
// being instrumented; fnName the function name of the node being
// instrumented; type_ the type of the node being instrumented. Returns the
// newly constructed function call node.
//
// Java: ProductionCoverageInstrumentationCallback.newInstrumentationNode(Node, String, String, Type)
func (c *ProductionCoverageInstrumentationCallback) newInstrumentationNode(
	node rhino.Node,
	fileName nullableString,
	fnName nullableString,
	type_ productionCoverageInstrumentationCallback_Type) rhino.Node {

	lineNo := node.GetLineno()
	columnNo := node.GetCharno()

	if node.IsBlock() {
		lineNo = node.GetParent().GetLineno()
		columnNo = node.GetParent().GetCharno()
	}

	encodedParam :=
		c.parameterMapping.getEncodedParam(fileName, fnName, type_, lineNo, columnNo)

	prop := rhino.IR_Getprop(rhino.IR_Name(c.instrumentationArrayName), "push")
	functionCall := rhino.IR_Call(prop, rhino.IR_String(encodedParam))
	exprNode := rhino.IR_ExprResult(functionCall)

	return exprNode.SrcrefTreeIfMissing(node)
}

// addElseBlock adds an else block for If statements if one is not already
// present.
//
// Java: ProductionCoverageInstrumentationCallback.addElseBlock(Node)
func (c *ProductionCoverageInstrumentationCallback) addElseBlock(ifNode ast.If) ast.Block {
	defaultBlock := ast.MustBlock(rhino.IR_Block().SrcrefTreeIfMissing(ifNode.Node()))
	ifNode.SetElse(defaultBlock)
	return defaultBlock
}

// Java: ProductionCoverageInstrumentationCallback.getInstrumentationMapping()
func (c *ProductionCoverageInstrumentationCallback) GetInstrumentationMapping() *VariableMap {
	return c.parameterMapping.getParamMappingAsVariableMap()
}

// nullableString represents an optional Java String used as an unboxed map/set key.
type nullableString struct {
	val string
	ok  bool
}

func toNullableString(s *string) nullableString {
	if s == nil {
		return nullableString{}
	}
	return nullableString{val: *s, ok: true}
}

// productionCoverageInstrumentationCallback_ParameterMapping maintains a
// mapping of unique identifiers to parameter values. It also generates unique
// identifiers by creating a counter starting from 0 and increments this value
// when assigning a new unique identifier. It converts the mapping to a
// VariableMap as well so that it can later be saved as a file.
//
// Java: private static final class ProductionCoverageInstrumentationCallback.ParameterMapping
type productionCoverageInstrumentationCallback_ParameterMapping struct {
	// Values are stored as a mapping of the String of the encoded array indices to a String of the
	// encoded uniqueIdentifier. This is so we can check if an encoded param has already been
	// defined so that we do not create a duplicate. (Ex. Key: ACA (Base64 VLQ encoding of [0,1,0]),
	// Value: C (Base64 VLQ encoding of the uniqueIdentifier). This map will later be inversed so
	// that it is printed in the following form: C:ACA.
	paramEncodings OrderedStringMap[string]

	// Insertion-ordered intern tables for file names, function names, and types.
	fileNames     OrderedSet[nullableString]
	functionNames OrderedSet[nullableString]
	typeNames     OrderedSet[string]

	nextUniqueIdentifier int64
	sb                   strings.Builder
}

// Java: ProductionCoverageInstrumentationCallback.ParameterMapping.getEncodedParam(String, String, Type, int, int)
func (m *productionCoverageInstrumentationCallback_ParameterMapping) getEncodedParam(
	fileKey nullableString,
	functionKey nullableString,
	type_ productionCoverageInstrumentationCallback_Type,
	lineNo int,
	colNo int) string {

	fileIdx := int32(m.fileNames.Intern(fileKey))
	functionIdx := int32(m.functionNames.Intern(functionKey))
	typeIdx := int32(m.typeNames.Intern(type_.Name()))

	m.sb.Reset()

	// The IOException that Java turns into an AssertionError cannot happen
	// when appending to a StringBuilder.
	sourcemap.Base64VLQ_Encode(&m.sb, fileIdx)
	sourcemap.Base64VLQ_Encode(&m.sb, functionIdx)
	sourcemap.Base64VLQ_Encode(&m.sb, typeIdx)
	sourcemap.Base64VLQ_Encode(&m.sb, int32(lineNo))
	sourcemap.Base64VLQ_Encode(&m.sb, int32(colNo))

	encodedParam := m.sb.String()

	if existing, ok := m.paramEncodings.Lookup(encodedParam); ok {
		return existing
	}

	uniqueIdentifier := m.generateUniqueIdentifier()
	if uniqueIdentifier > math.MaxInt32 {
		panic(jutil.NewArithmeticException(
			"Unique Identifier exceeds value of Integer.MAX_VALUE, could not encode with Base 64" +
				" VLQ"))
	}

	m.sb.Reset()
	sourcemap.Base64VLQ_Encode(&m.sb, int32(uniqueIdentifier))
	encodedId := m.sb.String()

	m.paramEncodings.Put(encodedParam, encodedId)

	return encodedId
}

// Java: ProductionCoverageInstrumentationCallback.ParameterMapping.generateUniqueIdentifier()
func (m *productionCoverageInstrumentationCallback_ParameterMapping) generateUniqueIdentifier() int64 {
	m.nextUniqueIdentifier++
	return m.nextUniqueIdentifier
}

func (m *productionCoverageInstrumentationCallback_ParameterMapping) formatNullableStringArrayJson(keys *OrderedSet[nullableString]) string {
	m.sb.Reset()
	m.sb.WriteByte('[')
	first := true
	for k := range keys.All() {
		if !first {
			m.sb.WriteByte(',')
		}
		first = false
		if k.ok {
			_ = jutil.WriteGsonString(&m.sb, k.val, false)
		} else {
			m.sb.WriteString("null")
		}
	}
	m.sb.WriteByte(']')
	return m.sb.String()
}

func (m *productionCoverageInstrumentationCallback_ParameterMapping) formatStringArrayJson(keys *OrderedSet[string]) string {
	m.sb.Reset()
	m.sb.WriteByte('[')
	first := true
	for k := range keys.All() {
		if !first {
			m.sb.WriteByte(',')
		}
		first = false
		_ = jutil.WriteGsonString(&m.sb, k, false)
	}
	m.sb.WriteByte(']')
	return m.sb.String()
}

// Java: ProductionCoverageInstrumentationCallback.ParameterMapping.getParamMappingAsVariableMap()
func (m *productionCoverageInstrumentationCallback_ParameterMapping) getParamMappingAsVariableMap() *VariableMap {
	pre := NewOrderedStringMapWithCapacity[string](m.paramEncodings.Len() + 3)
	for k, v := range m.paramEncodings.All() {
		pre.Put(k, v)
	}

	// Array names are given a " " (space) prefix since when writing to file, VariableMap.java
	// sorts the map by key values. This space will place the arrays at the top of the file.
	// The key and value entry are put in this order because the map will be inversed.
	pre.Put(m.formatNullableStringArrayJson(&m.fileNames), " FileNames")
	pre.Put(m.formatNullableStringArrayJson(&m.functionNames), " FunctionNames")
	pre.Put(m.formatStringArrayJson(&m.typeNames), " Types")

	origToNew := make(map[string]string, pre.Len())
	origKeys := make([]string, 0, pre.Len())
	for k, v := range pre.All() {
		if prevKey, dup := origToNew[v]; dup {
			panic(jutil.NewIllegalArgumentExceptionMsg(
				"Multiple entries with same value: " + k + "=" + v + " and " + prevKey + "=" + v))
		}
		origToNew[v] = k
		origKeys = append(origKeys, v)
	}
	return &VariableMap{
		origToNew: origToNew,
		newToOrig: pre.ToMap(),
		origKeys:  origKeys,
	}
}
