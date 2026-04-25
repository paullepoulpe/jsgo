// Port of com/google/javascript/jscomp/type/ReverseAbstractInterpreter.java
// @ closure-compiler v20260415.
//
// Java's ReverseAbstractInterpreter is a one-method interface. In Go
// it becomes an interface with the same contract: given a condition
// Node, a FlowScope without any knowledge of the outcome, and the
// known outcome, return a more precise FlowScope.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ReverseAbstractInterpreter defines what reversed abstract interpreters
// provide. Abstract interpretation is the process of interpreting a
// program at an abstracted level (such as at the type level) instead of
// the concrete level (the flow of values). This reversed abstract
// interpreter reverses the abstract interpretation process by knowing
// the outcome of some computation and calculating a preciser view of
// the world than the view without knowing the outcome of the computation.
//
// Port of com/google/javascript/jscomp/type/ReverseAbstractInterpreter.java.
type ReverseAbstractInterpreter interface {
	// GetPreciserScopeKnowingConditionOutcome calculates a precise version
	// of the scope knowing the outcome of the condition.
	//
	//  condition  – the condition's expression
	//  blindScope – the scope without knowledge about the outcome of the condition
	//  outcome    – the outcome of the condition (true → truthy branch, false → falsy branch)
	GetPreciserScopeKnowingConditionOutcome(condition *rhino.Node, blindScope *FlowScope, outcome bool) *FlowScope
}
