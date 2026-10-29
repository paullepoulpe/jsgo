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

func prewarmCompilerForParallelReads(compiler *Compiler) {
	if compiler == nil {
		return
	}
	_ = compiler.ensureSync()
	_ = compiler.GetAstAnalyzer()
	_ = compiler.GetCodingConvention()
	compiler.InitCompilerOptionsIfTesting()
}

func enableParallelTypeRegistryReads(compiler *Compiler) func() {
	if compiler == nil || compiler.typeRegistry == nil {
		return func() {}
	}
	prev := compiler.typeRegistry.SetParallelReadOnly(true)
	return func() { compiler.typeRegistry.SetParallelReadOnly(prev) }
}

func freezeGlobalScopeForParallelReads(gs *Scope) {
	if gs == nil {
		return
	}
	for imp := AbstractScope_ImplicitVar_ARGUMENTS; imp <= AbstractScope_ImplicitVar_THIS; imp++ {
		_ = gs.getOwnImplicitVar(imp)
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
		return rhino.NilNode, nil, false
	}
	scopeRoot := externs.GetParent()
	if scopeRoot.IsNil() || root.GetParent() != scopeRoot {
		return rhino.NilNode, nil, false
	}
	count := 0
	for c := externs.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if !c.IsScript() {
			return rhino.NilNode, nil, false
		}
		count++
	}
	for c := root.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if !c.IsScript() {
			return rhino.NilNode, nil, false
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

type recordingCompiler = Compiler

type parallelRecorder struct {
	events              []recordedCompilerEvent
	continueAfterErrors bool
	baseHalted          bool
	halted              bool
	errHandler          compiler_errorHandler
}

func newRecordingCompiler(base *Compiler) *Compiler {
	var rc Compiler
	initRecordingCompiler(&rc, base)
	return &rc
}

func initRecordingCompiler(rc *Compiler, base *Compiler) {
	rec := rc.parallelRecorder
	if rec == nil {
		rec = &parallelRecorder{}
	}
	_ = base.ensureSync()
	*rc = *base
	rec.events = rec.events[:0]
	opts := base.options
	rec.continueAfterErrors = opts != nil && opts.CanContinueAfterErrors()
	rec.baseHalted = !rec.continueAfterErrors && base.HasHaltingErrors()
	rec.halted = false
	rec.errHandler.c = rc
	rc.parallelRecorder = rec
	rc.asErrorHandler = &rec.errHandler
	rc.runtimeJsLibManager = nil
}

func resetRecordingCompiler(rc *Compiler, base *Compiler) {
	rec := rc.parallelRecorder
	rec.events = rec.events[:0]
	rec.halted = false
	rc.allowableFeatures = base.allowableFeatures
	rc.runtimeJsLibManager = nil
}

func (c *Compiler) takeEvents() []recordedCompilerEvent {
	r := c.parallelRecorder
	if r == nil || len(r.events) == 0 {
		return nil
	}
	out := make([]recordedCompilerEvent, len(r.events))
	copy(out, r.events)
	r.events = r.events[:0]
	return out
}

func (r *parallelRecorder) report(c *Compiler, err *JSError) {
	r.events = append(r.events, recordedCompilerEvent{
		kind: recordedEventReport,
		err:  err,
	})
	if !r.continueAfterErrors && !r.halted {
		if lvl := c.GetErrorLevel(err); lvl != nil && *lvl == CheckLevel_ERROR {
			r.halted = true
		}
	}
}

func (r *parallelRecorder) reportChangeToEnclosingScope(n rhino.Node) {
	scope := parallelChangeScopeForNode(n)
	r.events = append(r.events, recordedCompilerEvent{
		kind: recordedEventChangeScope,
		node: scope,
	})
}

func (r *parallelRecorder) reportChangeToChangeScope(changeScopeRoot rhino.Node) {
	jutil.CheckState(changeScopeRoot.IsScript() || changeScopeRoot.IsFunction())
	r.events = append(r.events, recordedCompilerEvent{
		kind: recordedEventChangeScope,
		node: changeScopeRoot,
	})
}

func (r *parallelRecorder) reportFunctionDeleted(node rhino.Node) {
	jutil.CheckState(node.IsFunction())
	r.events = append(r.events, recordedCompilerEvent{
		kind: recordedEventFunctionDeleted,
		node: node,
	})
}

func (c *Compiler) setWorkerAllowableFeatures(featureSet *parser.FeatureSet) {
	if featureSet.ContainsFeature(parser.FeatureSet_Feature_MODULES) &&
		featureSet.WithoutFeatures(parser.FeatureSet_Feature_MODULES).Equals(
			c.allowableFeatures.WithoutFeatures(parser.FeatureSet_Feature_MODULES)) {
		c.allowableFeatures = c.allowableFeatures.With(parser.FeatureSet_Feature_MODULES)
		c.parallelRecorder.events = append(c.parallelRecorder.events, recordedCompilerEvent{
			kind: recordedEventAllowModules,
		})
		return
	}
	panic(jutil.NewUnsupportedOperationExceptionMsg("SetAllowableFeatures unsupported in parallel worker"))
}

func (r *parallelRecorder) replayOnto(base *Compiler) {
	if r == nil || len(r.events) == 0 {
		return
	}
	replayCompilerEvents(base, r.events)
	r.events = r.events[:0]
}

func replayCompilerEvents(base *Compiler, events []recordedCompilerEvent) {
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
			initRecordingCompiler(&rc, base)
			for {
				idx := int(nextIdx.Add(1) - 1)
				if idx >= workCount {
					return
				}
				resetRecordingCompiler(&rc, base)
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

type parallelScriptOptions[S any] struct {
	hasExterns              bool
	externs                 rhino.Node
	useSharedGlobalScope    bool
	enableTypeRegistryReads bool
	stopOnHaltingError      bool
	prepare                 func(numWorkers, itemCount int) bool
	initWorker              func(workerID int, rc *Compiler) S
	callback                func(state S, idx int, parent rhino.Node, script rhino.Node, rc *Compiler) NodeTraversal_Callback
	afterScript             func(state S, idx int, parent rhino.Node, script rhino.Node, rc *Compiler)
	finishWorker            func(workerID int, state S)
}

func traverseScriptsParallel[S any](
	compiler *Compiler,
	root rhino.Node,
	opts parallelScriptOptions[S],
) bool {
	var scopeRoot rhino.Node
	var scripts []rhino.Node
	var items []scriptWithParent
	var itemCount int
	if opts.hasExterns {
		var ok bool
		scopeRoot, items, ok = collectExternAndRootScripts(opts.externs, root)
		if !ok {
			return false
		}
		itemCount = len(items)
	} else {
		var ok bool
		scripts, ok = collectRootScripts(root)
		if !ok {
			return false
		}
		scopeRoot = root
		itemCount = len(scripts)
	}
	numWorkers := parallelThreadCount(compiler, itemCount)
	if numWorkers <= 1 {
		return false
	}
	if opts.prepare != nil && !opts.prepare(numWorkers, itemCount) {
		return false
	}

	prewarmCompilerForParallelReads(compiler)
	if opts.enableTypeRegistryReads {
		restore := enableParallelTypeRegistryReads(compiler)
		defer restore()
	}
	var lazyGlobal *lazyGlobalScope
	if opts.useSharedGlobalScope {
		lazyGlobal = newLazyGlobalScope(compiler, scopeRoot)
	}

	runParallelWork(compiler, numWorkers, itemCount, func(workerID int) (func(int, *recordingCompiler), func()) {
		var workerRC *recordingCompiler
		var state S
		var sc *SyntacticScopeCreator
		initialized := false

		var cleanup func()
		if opts.finishWorker != nil {
			cleanup = func() {
				if initialized {
					opts.finishWorker(workerID, state)
				}
			}
		}

		runItem := func(idx int, rc *recordingCompiler) {
			if !initialized || workerRC != rc {
				workerRC = rc
				state = opts.initWorker(workerID, rc)
				initialized = true
				if opts.useSharedGlobalScope {
					sc = NewSyntacticScopeCreator(rc)
				}
			}
			parent := root
			var script rhino.Node
			if opts.hasExterns {
				parent = items[idx].parent
				script = items[idx].script
			} else {
				script = scripts[idx]
			}
			var cb NodeTraversal_Callback
			if opts.callback != nil {
				cb = opts.callback(state, idx, parent, script, rc)
			} else {
				cb = any(state).(NodeTraversal_Callback)
			}
			t := acquireNodeTraversalDirect(rc, cb, sc, false)
			traverseScriptUnderRoot(t, scopeRoot, parent, script, lazyGlobal)
			releaseNodeTraversal(t)
			if opts.afterScript != nil {
				opts.afterScript(state, idx, parent, script, rc)
			}
		}
		return runItem, cleanup
	}, opts.stopOnHaltingError)
	return true
}

func runParallelItemIndexed[W any](w *W, idx int, panicsByItem []any, runItem func(w *W, itemIdx int)) (panicked bool) {
	defer func() {
		if r := recover(); r != nil {
			panicsByItem[idx] = r
			panicked = true
		}
	}()
	runItem(w, idx)
	return false
}

func runParallelWorkersIndexed[W any](
	itemCount, numWorkers int,
	initWorker func(workerIdx int) W,
	runItem func(w *W, itemIdx int),
	finishWorker func(workerIdx int, w *W),
) {
	if itemCount <= 0 {
		return
	}
	if numWorkers > itemCount {
		numWorkers = itemCount
	}
	if numWorkers <= 1 {
		w := initWorker(0)
		if finishWorker != nil {
			defer finishWorker(0, &w)
		}
		for i := 0; i < itemCount; i++ {
			runItem(&w, i)
		}
		return
	}

	panicsByItem := make([]any, itemCount)
	var workerPanic any
	var workerPanicOnce sync.Once
	var nextIdx atomic.Uint32
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for workerIdx := 0; workerIdx < numWorkers; workerIdx++ {
		go func(wIdx int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					workerPanicOnce.Do(func() {
						workerPanic = r
					})
				}
			}()
			w := initWorker(wIdx)
			if finishWorker != nil {
				defer finishWorker(wIdx, &w)
			}
			for {
				idx := int(nextIdx.Add(1) - 1)
				if idx >= itemCount {
					return
				}
				if runParallelItemIndexed(&w, idx, panicsByItem, runItem) {
					return
				}
			}
		}(workerIdx)
	}
	wg.Wait()

	for _, p := range panicsByItem {
		if p != nil {
			panic(p)
		}
	}
	if workerPanic != nil {
		panic(workerPanic)
	}
}
