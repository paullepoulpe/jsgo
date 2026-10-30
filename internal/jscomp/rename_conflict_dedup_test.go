package jscomp

import (
	"fmt"
	"math/rand"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type oracleInverterVar struct {
	id            uint32
	name          string
	preferredName string
	scopeDepth    int32
	isCandidate   bool
	conflicts     jutil.LinkedHashSet[*oracleInverterVar]
}

func (v *oracleInverterVar) addConflict(other *oracleInverterVar) {
	if v.isCandidate && other != v {
		v.conflicts.Add(other)
	}
}

func (v *oracleInverterVar) tryToInvertName() {
	if !v.isCandidate {
		return
	}
	conflict := v.preferredName == MakeDeclaredNamesUnique_ARGUMENTS
	if !conflict {
		for other := range v.conflicts.All() {
			if other.name == v.preferredName {
				conflict = true
				break
			}
		}
	}
	newName := v.preferredName
	if conflict {
		disallowed := map[string]struct{}{MakeDeclaredNamesUnique_ARGUMENTS: {}}
		for other := range v.conflicts.All() {
			disallowed[other.name] = struct{}{}
		}
		base := v.preferredName + MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR
		for i := 0; ; i++ {
			newName = base + strconv.Itoa(i)
			if _, exists := disallowed[newName]; !exists {
				break
			}
		}
	}
	v.name = newName
	v.isCandidate = false
}

type oracleInverterScope struct {
	parent   *oracleInverterScope
	depth    int32
	declared []*oracleInverterVar
	refs     jutil.LinkedHashSet[*oracleInverterVar]
}

func (s *oracleInverterScope) recordPotentialConflicts(ref *oracleInverterVar) {
	if s.depth >= ref.scopeDepth {
		for _, decl := range s.declared {
			ref.addConflict(decl)
			decl.addConflict(ref)
		}
		if s.parent != nil {
			s.parent.recordPotentialConflicts(ref)
		}
	}
}

func TestRenameConflictDedupPreservesOrderAndNames(t *testing.T) {
	for seed := int64(0); seed < 30; seed++ {
		rng := rand.New(rand.NewSource(seed))
		inverter := NewMakeDeclaredNamesUnique_ContextualRenameInverter(nil)

		var oracleStack []*oracleInverterScope
		var oracleVars []*oracleInverterVar // 1-indexed like inverter.vars
		oracleVars = append(oracleVars, nil)

		// Include global/undeclared variables at depth 0.
		for g := 0; g < 6; g++ {
			name := fmt.Sprintf("v%d", g)
			id := rhino.RhinoStringPool_InternID(name)
			if len(inverter.vars) == 0 {
				inverter.vars = append(inverter.vars, inverterVar{})
			}
			varID := uint32(len(inverter.vars))
			inverter.vars = append(inverter.vars, inverterVar{
				name: name, nameID: id, scopeDepth: 0,
				refHead: -1, refTail: -1, conflictHead: -1, conflictTail: -1,
			})
			inverter.shadowByID[id] = varID
			oracleVars = append(oracleVars, &oracleInverterVar{id: varID, name: name, scopeDepth: 0})
		}

		maxDepth := 12
		if seed == 29 {
			// Exercise >64 deep scope nesting fallback as well.
			maxDepth = 70
		}

		var visibleByDepth [][]uint32
		visibleByDepth = append(visibleByDepth, []uint32{1, 2, 3, 4, 5, 6})

		pushRandomScope := func() {
			depth := int32(len(oracleStack) + 1)
			nDecl := 1 + rng.Intn(6)
			scope := &Scope{}
			scope.depth = depth
			var decls []*Var
			var oDecls []*oracleInverterVar
			vis := slices.Clone(visibleByDepth[len(visibleByDepth)-1])

			for i := 0; i < nDecl; i++ {
				baseIdx := rng.Intn(8)
				var current string
				if baseIdx == 7 {
					current = fmt.Sprintf("arguments$jscomp$%d_%d", depth, i)
				} else if rng.Intn(4) == 0 {
					current = fmt.Sprintf("plain%d_%d", depth, i)
				} else {
					current = fmt.Sprintf("v%d$jscomp$%d_%d", baseIdx, depth, i)
				}
				v := NewVar(current, rhino.NilNode, scope, i, nil, rhino.NilNode)
				decls = append(decls, v)

				nextID := uint32(len(oracleVars))
				ov := &oracleInverterVar{
					id:         nextID,
					name:       current,
					scopeDepth: depth,
				}
				if sep := strings.LastIndex(current, MakeDeclaredNamesUnique_ContextualRenamer_UNIQUE_ID_SEPARATOR); sep > 0 {
					ov.preferredName = current[:sep]
					ov.isCandidate = true
				}
				oracleVars = append(oracleVars, ov)
				oDecls = append(oDecls, ov)
				vis = append(vis, nextID)
			}

			inverter.pushScope(depth, decls)
			var parent *oracleInverterScope
			if len(oracleStack) > 0 {
				parent = oracleStack[len(oracleStack)-1]
			}
			oracleStack = append(oracleStack, &oracleInverterScope{
				parent:   parent,
				depth:    depth,
				declared: oDecls,
			})
			visibleByDepth = append(visibleByDepth, vis)
		}

		popTopScope := func() {
			top := len(oracleStack) - 1
			os := oracleStack[top]
			oracleStack = oracleStack[:top]
			visibleByDepth = visibleByDepth[:len(visibleByDepth)-1]

			for ref := range os.refs.All() {
				os.recordPotentialConflicts(ref)
			}
			for _, d := range os.declared {
				d.tryToInvertName()
			}
			inverter.popScope()
		}

		addRandomRefs := func() {
			if len(oracleStack) == 0 {
				return
			}
			vis := visibleByDepth[len(visibleByDepth)-1]
			nRefs := 1 + rng.Intn(10)
			topIdx := len(oracleStack) - 1
			for r := 0; r < nRefs; r++ {
				varID := vis[rng.Intn(len(vis))]
				inverter.addScopeRef(topIdx, varID)
				oracleStack[topIdx].refs.Add(oracleVars[varID])
			}
		}

		for step := 0; step < 200; step++ {
			switch {
			case len(oracleStack) == 0 || (len(oracleStack) < maxDepth && rng.Intn(3) != 0):
				pushRandomScope()
				addRandomRefs()
			case rng.Intn(2) == 0:
				addRandomRefs()
			default:
				popTopScope()
			}
		}
		for len(oracleStack) > 0 {
			addRandomRefs()
			popTopScope()
		}

		for id := uint32(1); int(id) < len(oracleVars); id++ {
			gotVar := &inverter.vars[id]
			wantVar := oracleVars[id]
			if gotVar.name != wantVar.name {
				t.Fatalf("seed %d var %d: got name %q, want %q", seed, id, gotVar.name, wantVar.name)
			}

			// Verify zero duplicate conflict links and exact conflict set match.
			var gotConflicts []uint32
			seen := make(map[uint32]bool)
			for idx := gotVar.conflictHead; idx >= 0; idx = inverter.conflictLinks[idx].next {
				target := inverter.conflictLinks[idx].varID
				if seen[target] {
					t.Fatalf("seed %d var %d: duplicate conflict edge to %d", seed, id, target)
				}
				seen[target] = true
				gotConflicts = append(gotConflicts, target)
			}
			var wantConflicts []uint32
			for other := range wantVar.conflicts.All() {
				wantConflicts = append(wantConflicts, other.id)
			}
			slices.Sort(gotConflicts)
			slices.Sort(wantConflicts)
			if !reflect.DeepEqual(gotConflicts, wantConflicts) {
				t.Fatalf("seed %d var %d: got conflicts %v, want %v", seed, id, gotConflicts, wantConflicts)
			}
		}

		inverter.resetArena()
	}
}
