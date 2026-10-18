package jscomp

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

type prebuildDependencyCallback struct {
	kind, source string
	level        CheckLevel
}

// The recorder is synchronized so a regression reports a useful assertion
// instead of introducing a race into the test itself. The goroutine check
// deterministically catches callbacks escaping ordered publication even if
// the worker completion order happens to match the input order.
type prebuildDependencyCallbackRecorder struct {
	owner    int64
	mu       sync.Mutex
	events   []prebuildDependencyCallback
	offOwner atomic.Bool
}

func (r *prebuildDependencyCallbackRecorder) record(kind string, level CheckLevel, error *JSError) {
	if jutil.GoroutineID() != r.owner {
		r.offOwner.Store(true)
	}
	source := ""
	if error.SourceName() != nil {
		source = *error.SourceName()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.events = append(r.events, prebuildDependencyCallback{kind, source, level})
}

func (r *prebuildDependencyCallbackRecorder) Report(level *CheckLevel, error *JSError) {
	r.record("handler", *level, error)
}

type prebuildDependencyWarningGuard struct {
	WarningsGuardBase
	recorder *prebuildDependencyCallbackRecorder
	level    CheckLevel
}

func (g *prebuildDependencyWarningGuard) Level(error *JSError) *CheckLevel {
	if error.Type() != ModuleLoader_LOAD_WARNING {
		return nil
	}
	g.recorder.record("guard", g.level, error)
	return g.level.Ptr()
}

func TestPrebuildDependencyDiagnosticsPreserveHandlerAndGuardRouting(t *testing.T) {
	for _, regex := range []bool{false, true} {
		for _, level := range []CheckLevel{CheckLevel_WARNING, CheckLevel_OFF} {
			t.Run(fmt.Sprintf("regex=%t/level=%s", regex, level.String()), func(t *testing.T) {
				recorder := &prebuildDependencyCallbackRecorder{owner: jutil.GoroutineID()}
				options := NewCompilerOptions()
				options.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_NEXT)
				options.SetErrorHandler(recorder)
				options.AddWarningsGuard(&prebuildDependencyWarningGuard{recorder: recorder, level: level})
				compiler := NewCompiler()
				compiler.InitOptions(options)
				compiler.SetErrorManager(NewSortingErrorManager(nil))
				compiler.SetPreferRegexParser(regex)
				inputs := []*CompilerInput{
					NewCompilerInput(SourceFile_FromCode("first.js", strings.Repeat("var padding = 0;\n", 200)+"import {first} from './missing-first.js';\n")),
					NewCompilerInput(SourceFile_FromCode("second.js", "import {second} from './missing-second.js';\n")),
				}
				known := []DependencyInfo{
					SimpleDependencyInfo_Builder_("first.js", "first.js").Build(),
					SimpleDependencyInfo_Builder_("second.js", "second.js").Build(),
				}
				compiler.moduleLoader = ModuleLoader_Builder_().
					SetModuleRoots([]string{"."}).SetInputs(known).
					SetFactory(BrowserModuleResolver_FACTORY).
					SetErrorHandler(compiler.AsErrorHandler()).Build()
				paths := make([]*ModuleLoader_ModulePath, len(inputs))
				for i, input := range inputs {
					input.SetCompiler(compiler)
					input.GetAstRoot(compiler) // exercise the already-parsed AST path too
					paths[i] = input.GetPath()
				}
				if len(recorder.events) != 0 {
					t.Fatalf("fixture unexpectedly reported diagnostics before dependency scanning: %v", recorder.events)
				}
				NewPrebuildDependencyInfo(2).Prebuild(slices.Values(inputs))
				if recorder.offOwner.Load() {
					t.Error("warning guard or custom error handler ran outside the publishing goroutine")
				}
				var want []prebuildDependencyCallback
				for _, source := range []string{"first.js", "second.js"} {
					want = append(want, prebuildDependencyCallback{"guard", source, level})
					if level != CheckLevel_OFF {
						want = append(want, prebuildDependencyCallback{"handler", source, level})
					}
				}
				if !reflect.DeepEqual(recorder.events, want) {
					t.Errorf("callback route/order got %#v, want %#v", recorder.events, want)
				}
				for i, input := range inputs {
					if input.GetPath() != paths[i] || input.GetPath().outer != compiler.moduleLoader {
						t.Errorf("input %d retained a private worker module path/loader", i)
					}
					if input.dependencyInfo == nil {
						t.Errorf("input %d did not publish its dependency info", i)
					}
				}
				wantWarnings := 0
				if level == CheckLevel_WARNING {
					wantWarnings = len(inputs)
				}
				if got := compiler.GetErrorManager().GetWarningCount(); got != wantWarnings {
					t.Errorf("manager warning count got %d, want %d", got, wantWarnings)
				}
				if got := compiler.GetErrorManager().GetErrorCount(); got != 0 {
					t.Errorf("manager error count got %d: warning guard was bypassed", got)
				}
			})
		}
	}
}
