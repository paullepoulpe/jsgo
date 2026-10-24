package jscomp

import (
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func parallelThreadCount(compiler AbstractCompiler, workCount int) int {
	if workCount < 2 || compiler == nil {
		return 1
	}
	opts := compiler.GetOptions()
	if opts == nil {
		return 1
	}
	n := opts.GetNumParallelThreads()
	if n <= 1 {
		return 1
	}
	if n > workCount {
		n = workCount
	}
	return n
}

func prewarmCompilerForParallelReads(compiler AbstractCompiler) {
	if compiler == nil {
		return
	}
	_ = compiler.GetAstAnalyzer()
	_ = compiler.GetCodingConvention()
	if c, ok := compiler.(*Compiler); ok && c != nil {
		c.InitCompilerOptionsIfTesting()
	}
}

func enableParallelTypeRegistryReads(compiler AbstractCompiler) func() {
	if c, ok := compiler.(*Compiler); ok {
		if c == nil || c.typeRegistry == nil {
			return func() {}
		}
		prev := c.typeRegistry.SetParallelReadOnly(true)
		return func() { c.typeRegistry.SetParallelReadOnly(prev) }
	}
	if compiler == nil || compiler.HasOptimizationColors() {
		return func() {}
	}
	reg := compiler.GetTypeRegistry()
	if reg == nil {
		return func() {}
	}
	prev := reg.SetParallelReadOnly(true)
	return func() { reg.SetParallelReadOnly(prev) }
}

func freezeGlobalScopeForParallelReads(gs *Scope) {
	if gs == nil {
		return
	}
	for i := range abstractScope_ImplicitVar_ptrs {
		_ = gs.GetOwnImplicitSlot(&abstractScope_ImplicitVar_ptrs[i])
	}
}

type lazyGlobalScope struct {
	once     sync.Once
	compiler AbstractCompiler
	root     rhino.Node
	scope    *Scope
}

func newLazyGlobalScope(compiler AbstractCompiler, root rhino.Node) *lazyGlobalScope {
	return &lazyGlobalScope{
		compiler: compiler,
		root:     root,
	}
}

func (l *lazyGlobalScope) get() *Scope {
	l.once.Do(func() {
		sc := NewSyntacticScopeCreator(l.compiler)
		gs := sc.CreateScopeOfScope(l.root, nil)
		freezeGlobalScopeForParallelReads(gs)
		l.scope = gs
	})
	return l.scope
}

func collectRootScripts(root rhino.Node) ([]rhino.Node, bool) {
	if root.IsNil() || !root.IsRoot() {
		return nil, false
	}
	count := 0
	for c := root.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if !c.IsScript() {
			return nil, false
		}
		count++
	}
	if count == 0 {
		return nil, true
	}
	scripts := make([]rhino.Node, 0, count)
	for c := root.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		scripts = append(scripts, c)
	}
	return scripts, true
}

type scriptWithParent struct {
	parent rhino.Node
	script rhino.Node
}

func collectExternAndRootScripts(externs rhino.Node, root rhino.Node) (rhino.Node, []scriptWithParent, bool) {
	if externs.IsNil() || root.IsNil() || !externs.IsRoot() || !root.IsRoot() {
		return rhino.Node{}, nil, false
	}
	scopeRoot := externs.GetParent()
	if scopeRoot.IsNil() || root.GetParent() != scopeRoot {
		return rhino.Node{}, nil, false
	}
	count := 0
	for c := externs.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if !c.IsScript() {
			return rhino.Node{}, nil, false
		}
		count++
	}
	for c := root.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if !c.IsScript() {
			return rhino.Node{}, nil, false
		}
		count++
	}
	if count == 0 {
		return scopeRoot, nil, true
	}
	items := make([]scriptWithParent, 0, count)
	for c := externs.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		items = append(items, scriptWithParent{parent: externs, script: c})
	}
	for c := root.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		items = append(items, scriptWithParent{parent: root, script: c})
	}
	return scopeRoot, items, true
}

type recordedCompilerEventKind uint8

const (
	recordedEventChangeScope recordedCompilerEventKind = iota + 1
	recordedEventFunctionDeleted
	recordedEventReport
	recordedEventAllowModules
)

type recordedCompilerEvent struct {
	kind recordedCompilerEventKind
	node rhino.Node
	err  *JSError
}

type parallelInternalError struct {
	msg   string
	cause jutil.Throwable
}

func parallelChangeScopeForNode(n rhino.Node) rhino.Node {
	if n.IsScript() {
		return n
	}
	enclosing := ChangeTracker_GetEnclosingChangeScopeRoot(n.GetParent())
	if enclosing.IsNil() {
		panic(jutil.NewIllegalStateExceptionMsg(
			"An enclosing scope is required for change reports but node " + n.String() + " doesn't have one."))
	}
	return enclosing
}

type recordingCompiler struct {
	AbstractCompiler

	events              []recordedCompilerEvent
	continueAfterErrors bool
	baseHalted          bool
	halted              bool
	allowModulesAdded   bool
}

func initRecordingCompiler(rc *recordingCompiler, base AbstractCompiler) {
	rc.AbstractCompiler = base
	rc.events = rc.events[:0]
	opts := base.GetOptions()
	rc.continueAfterErrors = opts != nil && opts.CanContinueAfterErrors()
	rc.baseHalted = !rc.continueAfterErrors && base.HasHaltingErrors()
	rc.halted = false
	rc.allowModulesAdded = false
}

func (rc *recordingCompiler) takeEvents() []recordedCompilerEvent {
	if len(rc.events) == 0 {
		return nil
	}
	out := make([]recordedCompilerEvent, len(rc.events))
	copy(out, rc.events)
	rc.events = rc.events[:0]
	return out
}

func (rc *recordingCompiler) Report(err *JSError) {
	rc.events = append(rc.events, recordedCompilerEvent{
		kind: recordedEventReport,
		err:  err,
	})
	if !rc.continueAfterErrors && !rc.halted {
		if lvl := rc.AbstractCompiler.GetErrorLevel(err); lvl != nil && *lvl == CheckLevel_ERROR {
			rc.halted = true
		}
	}
}

func (rc *recordingCompiler) HasHaltingErrors() bool {
	return rc.baseHalted || rc.halted
}

func (rc *recordingCompiler) ReportChangeToEnclosingScope(n rhino.Node) {
	scope := parallelChangeScopeForNode(n)
	rc.events = append(rc.events, recordedCompilerEvent{
		kind: recordedEventChangeScope,
		node: scope,
	})
}

func (rc *recordingCompiler) ReportChangeToChangeScope(changeScopeRoot rhino.Node) {
	jutil.CheckState(changeScopeRoot.IsScript() || changeScopeRoot.IsFunction())
	rc.events = append(rc.events, recordedCompilerEvent{
		kind: recordedEventChangeScope,
		node: changeScopeRoot,
	})
}

func (rc *recordingCompiler) ReportFunctionDeleted(node rhino.Node) {
	jutil.CheckState(node.IsFunction())
	rc.events = append(rc.events, recordedCompilerEvent{
		kind: recordedEventFunctionDeleted,
		node: node,
	})
}

func (rc *recordingCompiler) GetAllowableFeatures() *parser.FeatureSet {
	base := rc.AbstractCompiler.GetAllowableFeatures()
	if rc.allowModulesAdded {
		return base.With(parser.FeatureSet_Feature_MODULES)
	}
	return base
}

func (rc *recordingCompiler) SetAllowableFeatures(featureSet *parser.FeatureSet) {
	base := rc.AbstractCompiler.GetAllowableFeatures()
	if featureSet.ContainsFeature(parser.FeatureSet_Feature_MODULES) &&
		featureSet.WithoutFeatures(parser.FeatureSet_Feature_MODULES).Equals(
			base.WithoutFeatures(parser.FeatureSet_Feature_MODULES)) {
		rc.allowModulesAdded = true
		rc.events = append(rc.events, recordedCompilerEvent{
			kind: recordedEventAllowModules,
		})
		return
	}
	panic(jutil.NewUnsupportedOperationExceptionMsg("SetAllowableFeatures unsupported in parallel worker"))
}

func (rc *recordingCompiler) ThrowInternalError(msg string, cause jutil.Throwable) {
	panic(parallelInternalError{msg: msg, cause: cause})
}

func (rc *recordingCompiler) GetChangeTracker() *ChangeTracker {
	panic(jutil.NewUnsupportedOperationExceptionMsg("GetChangeTracker unsupported in parallel worker"))
}

func (rc *recordingCompiler) EnsureLibraryInjected(resourceName string, force bool) rhino.Node {
	panic(jutil.NewUnsupportedOperationExceptionMsg("EnsureLibraryInjected unsupported in parallel worker"))
}

func (rc *recordingCompiler) GetUniqueNameIdSupplier() func() string {
	panic(jutil.NewUnsupportedOperationExceptionMsg("GetUniqueNameIdSupplier unsupported in parallel worker"))
}

func (rc *recordingCompiler) SynthesizeExternsInput(name string) *CompilerInput {
	panic(jutil.NewUnsupportedOperationExceptionMsg("SynthesizeExternsInput unsupported in parallel worker"))
}

func (rc *recordingCompiler) MarkFeatureNotAllowed(feature parser.FeatureSet_Feature) {
	panic(jutil.NewUnsupportedOperationExceptionMsg("MarkFeatureNotAllowed unsupported in parallel worker"))
}

func (rc *recordingCompiler) MarkFeatureSetNotAllowed(featureSet *parser.FeatureSet) {
	panic(jutil.NewUnsupportedOperationExceptionMsg("MarkFeatureSetNotAllowed unsupported in parallel worker"))
}

func replayCompilerEvents(base AbstractCompiler, events []recordedCompilerEvent) {
	for i := range events {
		ev := &events[i]
		switch ev.kind {
		case recordedEventChangeScope:
			base.ReportChangeToChangeScope(ev.node)
		case recordedEventFunctionDeleted:
			base.ReportFunctionDeleted(ev.node)
		case recordedEventReport:
			base.Report(ev.err)
		case recordedEventAllowModules:
			base.SetAllowableFeatures(base.GetAllowableFeatures().With(parser.FeatureSet_Feature_MODULES))
		}
	}
}

func runParallelWork(
	base AbstractCompiler,
	numWorkers int,
	workCount int,
	initWorker func(workerID int) (runItem func(idx int, rc *recordingCompiler), cleanup func()),
	stopOnHaltingError bool,
) {
	eventsByItem := make([][]recordedCompilerEvent, workCount)
	panicsByItem := make([]any, workCount)

	var nextIdx atomic.Uint32
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for w := 0; w < numWorkers; w++ {
		workerID := w
		go func() {
			defer wg.Done()
			runItem, cleanup := initWorker(workerID)
			if cleanup != nil {
				defer cleanup()
			}
			var rc recordingCompiler
			for {
				idx := int(nextIdx.Add(1) - 1)
				if idx >= workCount {
					return
				}
				initRecordingCompiler(&rc, base)
				func() {
					defer func() {
						if r := recover(); r != nil {
							panicsByItem[idx] = r
						}
					}()
					runItem(idx, &rc)
				}()
				eventsByItem[idx] = rc.takeEvents()
			}
		}()
	}
	wg.Wait()

	for i := 0; i < workCount; i++ {
		if stopOnHaltingError && base.HasHaltingErrors() {
			break
		}
		if len(eventsByItem[i]) != 0 {
			replayCompilerEvents(base, eventsByItem[i])
		}
		if p := panicsByItem[i]; p != nil {
			if ie, ok := p.(parallelInternalError); ok {
				base.ThrowInternalError(ie.msg, ie.cause)
			}
			panic(p)
		}
	}
}
