// Port of test/com/google/javascript/jscomp/bundle/CachedTransformerTest.java @ closure-compiler v20260415

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

var (
	jsaCtFOO = Source_Builder_().SetCode("foo").Build()
	jsaCtBAR = Source_Builder_().SetCode("bar").Build()
	jsaCtBAZ = Source_Builder_().SetCode("baz").Build()
	jsaCtQUX = Source_Builder_().SetCode("qux").Build()
)

// jsaMockDelegate stands in for the Mockito mock: it answers from a list of
// stubbed (input, outputs...) and records calls.
type jsaMockDelegate struct {
	stubs map[*Source][]*Source
	calls map[*Source]int
}

func (m *jsaMockDelegate) apply(s *Source) *Source {
	m.calls[s]++
	answers := m.stubs[s]
	if len(answers) == 0 {
		return nil
	}
	a := answers[0]
	if len(answers) > 1 {
		m.stubs[s] = answers[1:]
	}
	return a
}

func TestCachedTransformer_Delegates(t *testing.T) {
	delegate := &jsaMockDelegate{stubs: map[*Source][]*Source{}, calls: map[*Source]int{}}
	cached := NewCachedTransformerWithBuilder(Source_Transformer_Of(delegate.apply), jutil.CacheBuilder_NewBuilder())
	delegate.stubs[jsaCtFOO] = []*Source{jsaCtBAR}
	delegate.stubs[jsaCtBAZ] = []*Source{jsaCtQUX}

	jsaTrue(t, cached.Transform(jsaCtFOO) == jsaCtBAR)
	jsaTrue(t, cached.Transform(jsaCtBAZ) == jsaCtQUX)
}

func TestCachedTransformer_Caches(t *testing.T) {
	delegate := &jsaMockDelegate{stubs: map[*Source][]*Source{}, calls: map[*Source]int{}}
	cached := NewCachedTransformerWithBuilder(Source_Transformer_Of(delegate.apply), jutil.CacheBuilder_NewBuilder())
	delegate.stubs[jsaCtFOO] = []*Source{jsaCtBAR, nil}

	jsaTrue(t, cached.Transform(jsaCtFOO) == jsaCtBAR)
	jsaTrue(t, cached.Transform(jsaCtFOO) == jsaCtBAR)
	jsaEq(t, delegate.calls[jsaCtFOO], 1)
}
