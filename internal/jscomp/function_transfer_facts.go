package jscomp

import (
	"os"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// functionTransferFacts is an opt-in, per-function immutable analysis
// generation. It canonicalizes the existing ordered name→Var map without
// inventing new Vars/Scopes, and lowers each CFG node's transfer into compact
// ordered actions. Both reaching-definition and reaching-use analyses share
// symbol IDs and escape facts. Their original recursive methods record the
// actions, preserving all stopping, conditional, and evaluation-order rules.
//
// Ownership: one function analysis creates this context after declaration/
// escape discovery; only its two serial fixed-point analyses may use it. The
// AST, declarations, escaped set, and CFG must remain unchanged throughout.
// invalidate ends that generation before any candidate rewrite and on panic.
// There is no global AST-pointer cache, cross-function reuse, or stale-node
// lookup. Once inactive, attached analyses fall back to their original code.
// Definition dependency sets are immutable after recording; consumers never
// mutate them. Lattice states and their ordered maps remain independently owned.
//
// Expected ceiling: removes repeated AST interpretation, symbol hashing, and
// identical definition/dependency allocations inside these two analyses. It
// does not accelerate parsing, type checking, global passes, or lattice joins.
type functionTransferFacts struct {
	active  bool
	vars    []*Var
	symbols map[string]uint32
	escaped []bool
	must    map[*rhino.Node]nativeMustProgram
	uses    map[*rhino.Node]nativeUseProgram
}

type nativeMustAction struct {
	definition       *MustBeReachingVariableDef_Definition
	symbol           uint32
	escapeParameters bool
}
type nativeMustProgram []nativeMustAction

type nativeUseAction struct {
	node   *rhino.Node
	symbol uint32
	kill   bool
}
type nativeUseProgram []nativeUseAction

func nativeFunctionFactsEnabled() bool {
	return os.Getenv("JSGO_EXPERIMENT_NATIVE_FACTS") == "1"
}

func newFunctionTransferFacts(vars *jutil.LinkedHashMap[string, *Var], escaped *jutil.ObjectLinkedHashSet[*Var]) *functionTransferFacts {
	f := &functionTransferFacts{
		active: true, symbols: make(map[string]uint32, vars.Size()),
		must: make(map[*rhino.Node]nativeMustProgram), uses: make(map[*rhino.Node]nativeUseProgram),
	}
	for name, v := range vars.All() {
		f.symbols[name] = uint32(len(f.vars))
		f.vars = append(f.vars, v)
		f.escaped = append(f.escaped, escaped.Contains(v))
	}
	return f
}

func (f *functionTransferFacts) invalidate() {
	f.active = false
	// Analyses can survive while inlining reads their final lattice states.
	// Release the action programs and their AST retention immediately.
	f.must, f.uses = nil, nil
	f.symbols, f.vars, f.escaped = nil, nil, nil
}

func (f *functionTransferFacts) recordDefinition(d *MustBeReachingVariableDef, name string, node, rValue *rhino.Node) {
	id, ok := f.symbols[name]
	if !ok {
		return
	}
	action := nativeMustAction{symbol: id}
	if node != nil && !f.escaped[id] {
		action.definition = NewMustBeReachingVariableDef_Definition(node)
		if rValue != nil {
			d.ComputeDependence(action.definition, rValue)
		}
	}
	*d.recording = append(*d.recording, action)
}

func (f *functionTransferFacts) flowMust(d *MustBeReachingVariableDef, node *rhino.Node, input *MustBeReachingVariableDef_MustDef) *MustBeReachingVariableDef_MustDef {
	program, found := f.must[node]
	if !found {
		func() {
			d.recording = &program
			defer func() { d.recording = nil }()
			d.ComputeMustDef(node, node, nil, false)
		}()
		f.must[node] = program
	}
	output := NewMustBeReachingVariableDef_MustDefCopy(input)
	for _, action := range program {
		if action.escapeParameters {
			d.EscapeParameters(output)
			continue
		}
		v := f.vars[action.symbol]
		for key, previous := range output.reachingDef.All() {
			if previous != nil && previous.depends.Contains(v) {
				output.reachingDef.Put(key, nil)
			}
		}
		if !f.escaped[action.symbol] {
			output.reachingDef.Put(v, action.definition)
		}
	}
	return output
}

func (f *functionTransferFacts) recordUse(u *MaybeReachingVariableUse, name string, node *rhino.Node, kill bool) {
	id, ok := f.symbols[name]
	if ok && !f.escaped[id] {
		*u.recording = append(*u.recording, nativeUseAction{symbol: id, node: node, kill: kill})
	}
}

func (f *functionTransferFacts) flowUses(u *MaybeReachingVariableUse, node *rhino.Node, input *MaybeReachingVariableUse_ReachingUses) *MaybeReachingVariableUse_ReachingUses {
	program, found := f.uses[node]
	if !found {
		func() {
			u.recording = &program
			defer func() { u.recording = nil }()
			u.ComputeMayUse(node, node, nil, u.HasExceptionHandler(node))
		}()
		f.uses[node] = program
	}
	output := NewMaybeReachingVariableUse_ReachingUsesCopy(input)
	for _, action := range program {
		v := f.vars[action.symbol]
		if action.kill {
			output.RemoveAll(v)
		} else {
			output.Put(v, action.node)
		}
	}
	return output
}
