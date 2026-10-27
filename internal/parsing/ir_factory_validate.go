// Port of com/google/javascript/jscomp/parsing/IRFactory.java @ closure-compiler v20261021

/*
 * Copyright 2013 The Closure Compiler Authors.
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

package parsing

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Java: IRFactory.validateAll(Node)
func (f *IRFactory) ValidateAll(n rhino.Node) {
	// Go-only: language-adapter — the ArrayDeque used as a stack (push/poll at
	// the head) is a slice used as a stack.
	var workBuf [32]rhino.Node
	work := workBuf[:0]
	for !n.IsNil() {
		f.Validate(n)
		nextSibling := n.GetNext()
		firstChild := n.GetFirstChild()
		if !firstChild.IsNil() {
			if !nextSibling.IsNil() {
				// handle the siblings later
				work = append(work, nextSibling)
			}
			n = firstChild
		} else if !nextSibling.IsNil() {
			// no children, handle the next sibling
			n = nextSibling
		} else {
			// no siblings, continue with work we have saved on the work queue
			if len(work) == 0 {
				n = rhino.NilNode
			} else {
				n = work[len(work)-1]
				work = work[:len(work)-1]
			}
		}
	}
	jutil.CheckState(len(work) == 0)
}

// Java: IRFactory.validate(Node)
func (f *IRFactory) Validate(n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_PARAM_LIST:
		f.ValidateParameters(n)
	case rhino.Token_BREAK, rhino.Token_CONTINUE:
		f.ValidateBreakContinue(n)
	case rhino.Token_RETURN:
		f.ValidateReturn(n)
	case rhino.Token_YIELD:
		f.ValidateYield(n)
	case rhino.Token_AWAIT, rhino.Token_FOR_AWAIT_OF:
		f.ValidateAwait(n)
	case rhino.Token_NEW_TARGET:
		f.ValidateNewDotTarget(n)
	case rhino.Token_LABEL:
		f.ValidateLabel(n)
	case rhino.Token_FUNCTION:
		f.ValidateBlockScopedFunctions(n)
	case rhino.Token_NAME:
		f.ValidateShadowHost(n)
	}
}

// Java: IRFactory.validateShadowHost(Node)
func (f *IRFactory) ValidateShadowHost(n rhino.Node) {
	if n.GetClosureUnawareShadow().IsNil() {
		return
	}

	if f.firstFileoverview == nil || !f.firstFileoverview.IsClosureUnawareCode() {
		f.Error(n, "invalid @closureUnaware node: must be in a file annotated with @closureUnaware")
		return
	}

	parent := n.GetParent()
	// We only support:

	//  - FUNCTION nodes that are direct CALL callees
	if parent.IsCall() {
		return
	}
	//  - FUNCTION nodes that are children of a GETPROP node that is a direct
	//    CALL callee (e.g. fn.call or fn.apply can be the CALL callee)
	grandparent := parent.GetParent()
	if parent.IsGetProp() && !grandparent.IsNil() && grandparent.IsCall() {
		calledFn := parent.GetString()
		if calledFn == "call" || calledFn == "apply" {
			return
		}
	}
	f.Error(n, "invalid @closureUnaware node: must be directly called")
}

// Java: IRFactory.validateAwait(Node)
func (f *IRFactory) ValidateAwait(n rhino.Node) {
	var errorMsg string
	if n.IsAwait() {
		errorMsg = IRFactory_UNEXPECTED_AWAIT
	} else if n.IsForAwaitOf() {
		errorMsg = IRFactory_UNEXPECTED_FOR_AWAIT_OF
	} else {
		return
	}

	parent := n
	for {
		parent = parent.GetParent()
		if parent.IsNil() {
			break
		}
		// The await is in a class static block or a field initializer.
		// e.g. `class C { static { await y; } }` or `class C { x = await y; }`
		if (parent.IsBlock() || parent.IsMemberFieldDef()) && parent.GetParent().IsClassMembers() {
			f.errorReporter.Error(errorMsg, f.sourceName, n.GetLineno(), n.GetCharno())
			return
		}
		if parent.IsAsyncFunction() {
			return
		}
		if parent.IsFunction() {
			// The await is in a non-async function.
			// e.g. `function f() { return await 5; }`
			// nested function e.g. `async function f() { function f2() { return await 5; } }`
			f.errorReporter.Error(errorMsg, f.sourceName, n.GetLineno(), n.GetCharno())
			return
		}
		if parent.IsScript() && parent.GetBooleanProp(rhino.Node_ES6_MODULE) {
			// This is a valid top-level await in an ES module.
			f.MaybeWarnForFeatureNode(n, parser.FeatureSet_Feature_TOP_LEVEL_AWAIT)
			return
		}
	}
	f.errorReporter.Error(errorMsg, f.sourceName, n.GetLineno(), n.GetCharno())
}

// Java: IRFactory.validateYield(Node)
func (f *IRFactory) ValidateYield(n rhino.Node) {
	if !n.IsYield() {
		return
	}

	parent := n
	for {
		parent = parent.GetParent()
		if parent.IsNil() {
			break
		}
		// The yield is in a class static block.
		// e.g. `class C { static { yield; } }`
		if parent.IsClassMembers() {
			f.errorReporter.Error(IRFactory_UNEXPECTED_YIELD, f.sourceName, n.GetLineno(), n.GetCharno())
			return
		}
		if parent.IsGeneratorFunction() {
			return
		}
	}
}

// Java: IRFactory.validateReturn(Node)
func (f *IRFactory) ValidateReturn(n rhino.Node) {
	if !n.IsReturn() {
		return
	}

	parent := n
	for {
		parent = parent.GetParent()
		if parent.IsNil() {
			break
		}
		// The return is in a class static block.
		// e.g. `class C { static { return; } }`
		if parent.IsClassMembers() {
			f.errorReporter.Error(IRFactory_UNEXPECTED_RETURN, f.sourceName, n.GetLineno(), n.GetCharno())
			return
		}
		if parent.IsFunction() {
			return
		}
	}
	f.errorReporter.Error(IRFactory_UNEXPECTED_RETURN, f.sourceName, n.GetLineno(), n.GetCharno())
}

// Java: IRFactory.validateNewDotTarget(Node)
func (f *IRFactory) ValidateNewDotTarget(n rhino.Node) {
	if n.GetToken() != rhino.Token_NEW_TARGET {
		return
	}

	child := n
	for parent := n.GetParent(); !parent.IsNil(); child, parent = parent, parent.GetParent() {
		// Skip arrow functions, which use the `new.target` of the code around them.
		if parent.IsFunction() && !parent.IsArrowFunction() {
			return
		}
		if parent.IsMemberFieldDef() ||
			(parent.IsComputedFieldDef() && child.IsSecondChildOf(parent)) ||
			(parent.IsClassMembers() && child.IsBlock()) {
			// Class field initializers and static blocks run with no [[NewTarget]], so `new.target` is
			// always `undefined` there. That is rarely what the code means, so it is an error.
			f.errorReporter.Error(
				IRFactory_NEW_DOT_TARGET_IN_CLASS_FIELD_OR_STATIC_BLOCK,
				f.sourceName,
				n.GetLineno(),
				n.GetCharno())
			return
		}
	}
	f.errorReporter.Error(IRFactory_UNEXPECTED_NEW_DOT_TARGET, f.sourceName, n.GetLineno(), n.GetCharno())
}

// Java: IRFactory.validateBreakContinue(Node)
func (f *IRFactory) ValidateBreakContinue(n rhino.Node) {
	if !n.IsBreak() && !n.IsContinue() {
		return
	}

	labelName := n.GetFirstChild()
	if !labelName.IsNil() {
		parent := n.GetParent()
		for !parent.IsLabel() || !IRFactory_LabelsMatch(parent, labelName) {
			if parent.IsFunction() || parent.IsScript() || parent.IsClassMembers() {
				// report missing label
				f.errorReporter.Error(
					jutil.StringFormat(IRFactory_UNDEFINED_LABEL, labelName.GetString()),
					f.sourceName,
					n.GetLineno(),
					n.GetCharno())
				break
			}
			parent = parent.GetParent()
		}

		if parent.IsLabel() && IRFactory_LabelsMatch(parent, labelName) {
			if n.IsContinue() && !IRFactory_IsContinueTarget(parent.GetLastChild()) {
				// report invalid continue target
				f.errorReporter.Error(
					IRFactory_UNEXPECTED_LABELLED_CONTINUE, f.sourceName, n.GetLineno(), n.GetCharno())
			}
		}
	} else if n.IsContinue() {
		parent := n.GetParent()
		for !IRFactory_IsContinueTarget(parent) {
			if parent.IsFunction() || parent.IsScript() || parent.IsClassMembers() {
				// report invalid continue
				f.errorReporter.Error(IRFactory_UNEXPECTED_CONTINUE, f.sourceName, n.GetLineno(), n.GetCharno())
				break
			}
			parent = parent.GetParent()
		}
	} else {
		parent := n.GetParent()
		for !IRFactory_IsBreakTarget(parent) {
			if parent.IsFunction() || parent.IsScript() || parent.IsClassMembers() {
				// report invalid break
				f.errorReporter.Error(IRFactory_UNLABELED_BREAK, f.sourceName, n.GetLineno(), n.GetCharno())
				break
			}
			parent = parent.GetParent()
		}
	}
}

// Java: IRFactory.isBreakTarget(Node)
func IRFactory_IsBreakTarget(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF,
		rhino.Token_WHILE, rhino.Token_DO, rhino.Token_SWITCH:
		return true
	default:
		return false
	}
}

// Java: IRFactory.isContinueTarget(Node)
func IRFactory_IsContinueTarget(n rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF,
		rhino.Token_WHILE, rhino.Token_DO:
		return true
	default:
		return false
	}
}

// Java: IRFactory.labelsMatch(Node, Node)
func IRFactory_LabelsMatch(label rhino.Node, labelName rhino.Node) bool {
	return label.GetFirstChild().GetString() == labelName.GetString()
}

// Java: IRFactory.validateLabel(Node)
func (f *IRFactory) ValidateLabel(n rhino.Node) {
	if !n.IsLabel() {
		return
	}

	labelName := n.GetFirstChild()
	for parent := n.GetParent(); !parent.IsNil() && !parent.IsFunction(); parent = parent.GetParent() {
		if parent.IsLabel() && IRFactory_LabelsMatch(parent, labelName) {
			f.errorReporter.Error(
				jutil.StringFormat(IRFactory_DUPLICATE_LABEL, labelName.GetString()),
				f.sourceName,
				n.GetLineno(),
				n.GetCharno())
			break
		}
	}
}

// Java: IRFactory.validateParameters(Node)
func (f *IRFactory) ValidateParameters(n rhino.Node) {
	if !n.IsParamList() {
		return
	}
	first := n.GetFirstChild()
	if first.IsNil() {
		return
	}
	if first.IsName() {
		second := first.GetNext()
		if second.IsNil() {
			return
		}
		var seen [8]string
		seen[0] = first.GetString()
		count := 1
		allSimple := true
		for c := second; !c.IsNil(); c = c.GetNext() {
			if !c.IsName() || count >= len(seen) {
				allSimple = false
				break
			}
			seen[count] = c.GetString()
			count++
		}
		if allSimple {
			hasDup := false
			for i := 1; i < count; i++ {
				for j := 0; j < i; j++ {
					if seen[i] == seen[j] {
						hasDup = true
						break
					}
				}
				if hasDup {
					break
				}
			}
			if !hasDup {
				return
			}
		}
	}
	f.validateParametersSlow(n)
}

func (f *IRFactory) validateParametersSlow(n rhino.Node) {
	seenNames := jutil.NewLinkedHashSet[string]()
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		ParsingUtil_GetParamOrPatternNames(
			c,
			func(param rhino.Node) {
				paramName := param.GetString()
				if !seenNames.Add(paramName) {
					f.errorReporter.Warning(
						jutil.StringFormat(IRFactory_DUPLICATE_PARAMETER, paramName),
						f.sourceName,
						param.GetLineno(),
						param.GetCharno())
				}
			})
	}
}

// Java: IRFactory.validateBlockScopedFunctions(Node)
func (f *IRFactory) ValidateBlockScopedFunctions(n rhino.Node) {
	if n.IsFunction() && n.GetParent().IsBlock() && !n.GetGrandparent().IsFunction() {
		f.MaybeWarnForFeatureNode(n, parser.FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION)
	}
}

// Java: IRFactory.reportErrorIfYieldOrAwaitInDefaultValue(Node)
func (f *IRFactory) ReportErrorIfYieldOrAwaitInDefaultValue(defaultValueNode rhino.Node) {
	yieldNode := IRFactory_FindNodeTypeInExpression(defaultValueNode, rhino.Token_YIELD)
	if !yieldNode.IsNil() {
		f.errorReporter.Error(
			"`yield` is illegal in parameter default value.",
			irFactory_nullableString(yieldNode.GetSourceFileName()),
			yieldNode.GetLineno(),
			yieldNode.GetCharno())
	}
	awaitNode := IRFactory_FindNodeTypeInExpression(defaultValueNode, rhino.Token_AWAIT)
	if !awaitNode.IsNil() {
		f.errorReporter.Error(
			"`await` is illegal in parameter default value.",
			irFactory_nullableString(awaitNode.GetSourceFileName()),
			awaitNode.GetLineno(),
			awaitNode.GetCharno())
	}
}

// irFactory_nullableString passes a nullable Java String to a rhino
// ErrorReporter.
//
// Go-only: language-adapter — rhino.ErrorReporter takes a Go string, so
// Java's null is "".
func irFactory_nullableString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// IRFactory_FindNodeTypeInExpression tries to find a node with the given
// token in the given expression. Returns the first one found in pre-order
// traversal or `null` if none found. Will not traverse into function or
// class expressions.
//
// Java: IRFactory.findNodeTypeInExpression(Node, Token)
func IRFactory_FindNodeTypeInExpression(expressionNode rhino.Node, token rhino.Token) rhino.Node {
	// Go-only: language-adapter — the ArrayDeque used as a FIFO queue
	// (add at the tail, remove at the head) is a slice with a head index.
	worklist := []rhino.Node{expressionNode}
	for head := 0; head < len(worklist); {
		node := worklist[head]
		head++
		if node.GetToken() == token {
			return node
		} else if !node.IsFunction() && !node.IsClass() {
			for child := node.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
				worklist = append(worklist, child)
			}
		}
	}
	return rhino.NilNode
}
