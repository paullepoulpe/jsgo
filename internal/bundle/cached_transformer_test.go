// Port of test/com/google/javascript/jscomp/bundle/CachedTransformerTest.java @ closure-compiler v20261028

/*
 * Copyright 2017 The Closure Compiler Authors.
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

package bundle

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var (
	jsaCtFOO = NewSource(SourceConfig{Code: "foo"})
	jsaCtBAR = NewSource(SourceConfig{Code: "bar"})
	jsaCtBAZ = NewSource(SourceConfig{Code: "baz"})
	jsaCtQUX = NewSource(SourceConfig{Code: "qux"})
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
	cached := NewCachedTransformerWithMaxSize(Source_Transformer_Of(delegate.apply), -1)
	delegate.stubs[jsaCtFOO] = []*Source{jsaCtBAR}
	delegate.stubs[jsaCtBAZ] = []*Source{jsaCtQUX}

	jsaTrue(t, cached.Transform(jsaCtFOO) == jsaCtBAR)
	jsaTrue(t, cached.Transform(jsaCtBAZ) == jsaCtQUX)
}

func TestCachedTransformer_Caches(t *testing.T) {
	delegate := &jsaMockDelegate{stubs: map[*Source][]*Source{}, calls: map[*Source]int{}}
	cached := NewCachedTransformerWithMaxSize(Source_Transformer_Of(delegate.apply), -1)
	delegate.stubs[jsaCtFOO] = []*Source{jsaCtBAR, nil}

	jsaTrue(t, cached.Transform(jsaCtFOO) == jsaCtBAR)
	jsaTrue(t, cached.Transform(jsaCtFOO) == jsaCtBAR)
	jsaEq(t, delegate.calls[jsaCtFOO], 1)
}

func TestOptionsCompilerSupplierConcurrentAndNodeStoreReset(t *testing.T) {
	rhino.ResetNodeStore()
	defer rhino.ResetNodeStore()

	cached := NewCachedTransformerWithMaxSize(
		NewTranspiler(NewCompilerSupplier(), "es6_runtime"),
		64,
	)

	const workers = 8
	start := make(chan struct{})
	var wg sync.WaitGroup
	errCh := make(chan error, workers)

	for g := 0; g < workers; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			<-start
			src := NewSource(SourceConfig{
				Path: jutil.Path_Of(fmt.Sprintf("mod_%d.js", id)),
				Code: fmt.Sprintf("class Box%d { get() { return [1, 2].map(x => x + %d); } }", id, id),
			})
			out := cached.Transform(src)
			if !strings.Contains(out.Code(), fmt.Sprintf("Box%d", id)) {
				errCh <- fmt.Errorf("worker %d unexpected output: %q", id, out.Code())
			}
		}(g)
	}
	close(start)
	wg.Wait()
	close(errCh)
	for err := range errCh {
		t.Fatal(err)
	}

	if !rhino.CurrentNodeStore().IsStatic() {
		t.Fatalf("expected rhino.CurrentNodeStore().IsStatic() == true after standalone OptionsCompilerSupplier calls")
	}

	// Verify an already-active outer compilation NodeStore is NOT reset.
	outer := rhino.NewNodeStore()
	prev := rhino.SetCurrentNodeStore(outer)
	defer rhino.SetCurrentNodeStore(prev)

	res := NewCompilerSupplier().Compile("nested.js", "const y = () => 1;")
	if len(res.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
	if rhino.CurrentNodeStore() != outer {
		t.Fatalf("OptionsCompilerSupplier clobbered active outer NodeStore")
	}
}
