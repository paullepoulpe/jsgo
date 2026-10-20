package jscomp

import (
	"reflect"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type peepholePublicFeatureProbe struct {
	AbstractPeepholeOptimizationBase
	name            string
	events          *[]string
	getters, clears int
}

func newPeepholePublicFeatureProbe(name string, events *[]string) *peepholePublicFeatureProbe {
	p := &peepholePublicFeatureProbe{name: name, events: events}
	p.InitAbstractPeepholeOptimizationBase(p)
	return p
}

func (p *peepholePublicFeatureProbe) BeginTraversal(compiler AbstractCompiler) {
	*p.events = append(*p.events, p.name+".begin")
	p.AbstractPeepholeOptimizationBase.BeginTraversal(compiler)
}

func (p *peepholePublicFeatureProbe) OptimizeSubtree(n *rhino.Node) *rhino.Node {
	if n.IsScript() {
		*p.events = append(*p.events, p.name+".optimize")
		p.AddFeatureToEnclosingScript(parser.FeatureSet_Feature_LET_DECLARATIONS)
	}
	return n
}

func (p *peepholePublicFeatureProbe) GetAbstractPeepholeOptimizationBase() *AbstractPeepholeOptimizationBase {
	p.getters++
	*p.events = append(*p.events, p.name+".get")
	return &p.AbstractPeepholeOptimizationBase
}

func (p *peepholePublicFeatureProbe) ClearNewFeatures() {
	p.clears++
	*p.events = append(*p.events, p.name+".clear")
	p.AbstractPeepholeOptimizationBase.ClearNewFeatures()
}

func (p *peepholePublicFeatureProbe) EndTraversal() {
	*p.events = append(*p.events, p.name+".end")
	p.AbstractPeepholeOptimizationBase.EndTraversal()
}

func TestPeepholePublicFeatureCallbacksKeepOrder(t *testing.T) {
	for _, count := range []int{2, 17} {
		name := "masked"
		if count > 16 {
			name = "fallback"
		}
		t.Run(name, func(t *testing.T) {
			fixture := newSyntacticScopeCreatorTest(t)
			var events []string
			var probes []*peepholePublicFeatureProbe
			var optimizations []AbstractPeepholeOptimization
			for i := 0; i < count; i++ {
				p := newPeepholePublicFeatureProbe(string(rune('a'+i)), &events)
				probes = append(probes, p)
				optimizations = append(optimizations, p)
			}
			pass := NewPeepholeOptimizationsPassOfList(fixture.compiler, "public-feature-order", optimizations)
			if len(events) != 0 {
				t.Errorf("constructor called custom feature getters: %v", events)
			}
			events = nil
			for _, p := range probes {
				p.getters, p.clears = 0, 0
			}
			pass.SetRetraverseOnChange(false)
			root := fixture.getRoot("consume();")
			fixture.compiler.SetAllowableFeatures(parser.FeatureSet_ES_NEXT)
			pass.Process(rhino.IR_Root(), root)
			var want []string
			for _, p := range probes {
				want = append(want, p.name+".begin")
			}
			for _, p := range probes {
				want = append(want, p.name+".optimize")
			}
			for _, p := range probes {
				want = append(want, p.name+".get", p.name+".clear")
			}
			for _, p := range probes {
				want = append(want, p.name+".end")
			}
			if !reflect.DeepEqual(events, want) {
				t.Errorf("feature callback order:\n got %v\nwant %v", events, want)
			}
			for _, p := range probes {
				if p.getters != 1 || p.clears != 1 {
					t.Errorf("%s getter/virtual clear counts = %d/%d, want 1/1", p.name, p.getters, p.clears)
				}
			}
			features := NodeUtil_GetFeatureSetOfScript(root)
			if features == nil || !features.ContainsFeature(parser.FeatureSet_Feature_LET_DECLARATIONS) {
				t.Fatal("custom optimizer's feature was not published to script")
			}
		})
	}
}

func TestPeepholeBuiltinPendingFeatureSurvivesDispatchSetup(t *testing.T) {
	fixture := newSyntacticScopeCreatorTest(t)
	root := fixture.getRoot("consume();")
	fixture.compiler.SetAllowableFeatures(parser.FeatureSet_ES_NEXT)
	fold := NewPeepholeFoldConstants(false, false)
	// All optimizations expose this existing API, even if their normal stock
	// implementation does not itself introduce a feature on this input.
	fold.AddFeatureToEnclosingScript(parser.FeatureSet_Feature_LET_DECLARATIONS)
	pass := NewPeepholeOptimizationsPass(fixture.compiler, "builtin-pending-feature", fold)
	pass.SetRetraverseOnChange(false)
	pass.Process(rhino.IR_Root(), root)
	features := NodeUtil_GetFeatureSetOfScript(root)
	if features == nil || !features.ContainsFeature(parser.FeatureSet_Feature_LET_DECLARATIONS) {
		t.Fatal("pending feature from exact builtin was not published")
	}
	if !fold.GetNewFeatures().IsEmpty() {
		t.Fatal("pending builtin feature was not cleared at scope exit")
	}
}
