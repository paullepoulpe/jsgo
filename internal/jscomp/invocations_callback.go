// Port of com/google/javascript/jscomp/InvocationsCallback.java @ closure-compiler v20261021

/*
 * Copyright 2007 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// InvocationsCallback is a traversal callback that finds method
// invocations of the form
//
//	call
//	  getprop
//	    ...
//	    string
//	  ...
//
// and invokes a method defined by subclasses for processing these
// invocations.
//
// Go-only: language-adapter — Java's two visit overloads are Visit (the
// Callback method) and VisitWithCallName (the abstract template method).
type InvocationsCallback interface {
	NodeTraversal_Callback

	// VisitWithCallName is called for each callnode that is a method
	// invocation.
	//
	// callNode is a node of type call; parent the parent of callNode;
	// callName the name of method invoked by first child of call.
	//
	// Java: InvocationsCallback.visit(NodeTraversal, Node, Node, String)
	VisitWithCallName(t *NodeTraversal, callNode *rhino.Node, parent *rhino.Node, callName string)

	GetInvocationsCallbackBase() *InvocationsCallbackBase
}

// InvocationsCallbackBase holds InvocationsCallback's implementation.
type InvocationsCallbackBase struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	self InvocationsCallback
}

// InitInvocationsCallbackBase is the InvocationsCallback() constructor;
// self is the callback being constructed.
//
// Java: InvocationsCallback()
func (b *InvocationsCallbackBase) InitInvocationsCallbackBase(self InvocationsCallback) {
	b.self = self
}

// GetInvocationsCallbackBase gives access to the base through the
// interface.
func (b *InvocationsCallbackBase) GetInvocationsCallbackBase() *InvocationsCallbackBase { return b }

// Java: InvocationsCallback.visit(NodeTraversal, Node, Node)
func (b *InvocationsCallbackBase) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	if !n.IsCall() {
		return
	}

	callee := n.GetFirstChild()
	if !callee.IsGetProp() {
		return
	}

	b.self.VisitWithCallName(t, n, parent, callee.GetString())
}
