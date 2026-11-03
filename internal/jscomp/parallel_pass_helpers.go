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
	if n := opts.GetNumParallelThreads(); n > 1 {
		return min(n, workCount)
	}
	return 1
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

func countRootScripts(root rhino.Node) (int, bool) {
	if root.IsNil() || !root.IsRoot() {
		return 0, false
	}
	count := 0
	for c := root.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if !c.IsScript() {
			return 0, false
		}
		count++
	}
	return count, true
}

func collectRootScripts(root rhino.Node) ([]rhino.Node, bool) {
	count, ok := countRootScripts(root)
	if !ok || count == 0 {
		return nil, ok
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
	extCount, ok1 := countRootScripts(externs)
	rootCount, ok2 := countRootScripts(root)
	if !ok1 || !ok2 {
		return rhino.NilNode, nil, false
	}
	scopeRoot := externs.GetParent()
	if scopeRoot.IsNil() || root.GetParent() != scopeRoot {
		return rhino.NilNode, nil, false
	}
	if extCount+rootCount == 0 {
		return scopeRoot, nil, true
	}
	items := make([]scriptWithParent, 0, extCount+rootCount)
	for _, parent := range [2]rhino.Node{externs, root} {
		for c := parent.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			items = append(items, scriptWithParent{parent: parent, script: c})
		}
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
	if !r.continueAfterErrors && !r.halted && reportHalts(c, err) {
		r.halted = true
	}
}

// reportHalts reports whether Compiler.Report(err) makes HasHaltingErrors
// true: Report gives err its default level unless the warnings guard sets
// one, and the error managers count as halting the errors reported at level
// ERROR whose type is an error (SortingErrorManager.originalErrorCount,
// ErrorManager_HasHaltingErrors), not the warnings a guard promotes.
//
// Go-only: workers used the guard's level alone, so an error at its default
// level did not stop a worker's script where the sequential traversal stops
// (a CombinedCompilerPass then reported later diagnostics of the same
// script), and a promoted warning stopped it where the sequential one goes
// on.
func reportHalts(c *Compiler, err *JSError) bool {
	if err.Type().level != CheckLevel_ERROR {
		return false
	}
	level := err.DefaultLevel()
	if c.warningsGuard != nil {
		if newLevel, ok := c.warningsGuard.Level(err); ok {
			level = newLevel
		}
	}
	return level == CheckLevel_ERROR
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

func (c *Compiler) setWorkerAllowableFeatures(featureSet parser.FeatureSet) {
	if featureSet.ContainsFeature(parser.FeatureSet_Feature_MODULES) &&
		featureSet.WithoutFeatures(parser.FeatureSet_Feature_MODULES) ==
			c.allowableFeatures.WithoutFeatures(parser.FeatureSet_Feature_MODULES) {
		c.allowableFeatures = c.allowableFeatures.With(parser.FeatureSet_Feature_MODULES)
		c.parallelRecorder.events = append(c.parallelRecorder.events, recordedCompilerEvent{
			kind: recordedEventAllowModules,
		})
		return
	}
	panic(jutil.NewUnsupportedOperationExceptionMsg("SetAllowableFeatures unsupported in parallel worker"))
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
			if n := ev.err.Node(); !n.IsNil() {
				n.PinForDiagnostic() // (a no-op unless the worker deferred it)
			}
			base.Report(ev.err)
		case recordedEventAllowModules:
			base.SetAllowableFeatures(base.GetAllowableFeatures().With(parser.FeatureSet_Feature_MODULES))
		}
	}
}

type parallelWorkState struct {
	rc      recordingCompiler
	runItem func(idx int, rc *recordingCompiler)
	cleanup func()
}

func rethrowParallelPanic(base AbstractCompiler, p any) {
	if ie, ok := p.(parallelInternalError); ok {
		base.ThrowInternalError(ie.msg, ie.cause)
	}
	panic(p)
}

func runParallelWork(
	base AbstractCompiler,
	numWorkers int,
	workCount int,
	initWorker func(workerID int) (runItem func(idx int, rc *recordingCompiler), cleanup func()),
	stopOnHaltingError bool,
) {
	if workCount <= 0 {
		return
	}
	eventsByItem := make([][]recordedCompilerEvent, workCount)
	panicsByItem, workerPanic := collectParallelWorkersDeferredPins(
		workCount,
		max(min(numWorkers, workCount), 1),
		func(workerID int, w *parallelWorkState) {
			w.runItem, w.cleanup = initWorker(workerID)
			initRecordingCompiler(&w.rc, base)
		},
		func(w *parallelWorkState, idx int) {
			resetRecordingCompiler(&w.rc, base)
			w.runItem(idx, &w.rc)
		},
		func(w *parallelWorkState, idx int) {
			eventsByItem[idx] = w.rc.takeEvents()
		},
		func(_ int, w *parallelWorkState) {
			if w.cleanup != nil {
				w.cleanup()
			}
		},
		false,
	)

	for i := 0; i < workCount; i++ {
		if stopOnHaltingError && base.HasHaltingErrors() {
			break
		}
		if len(eventsByItem[i]) != 0 {
			replayCompilerEvents(base, eventsByItem[i])
		}
		if p := panicsByItem[i]; p != nil {
			rethrowParallelPanic(base, p)
		}
	}
	if workerPanic != nil {
		rethrowParallelPanic(base, workerPanic)
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
	// runScript, when set, replaces the traversal of each script (callback
	// and afterScript are not used): for passes that are a function of each
	// script but not a NodeTraversal callback (AstValidator).
	runScript func(state S, idx int, parent rhino.Node, script rhino.Node, rc *Compiler)
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
			if !initialized {
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
			if opts.runScript != nil {
				opts.runScript(state, idx, parent, script, rc)
				return
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

func spawnIndexedWorkers(
	itemCount, numWorkers int,
	body func(workerIdx int, claim func() (int, bool)),
) any {
	n := min(max(numWorkers, 0), itemCount)
	if n <= 0 {
		return nil
	}
	var workerPanic any
	var workerPanicOnce sync.Once
	var nextIdx atomic.Uint32
	claim := func() (int, bool) {
		idx := int(nextIdx.Add(1) - 1)
		if idx >= itemCount {
			return 0, false
		}
		return idx, true
	}
	var wg sync.WaitGroup
	wg.Add(n)
	for workerIdx := range n {
		go func(wIdx int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					workerPanicOnce.Do(func() { workerPanic = r })
				}
			}()
			body(wIdx, claim)
		}(workerIdx)
	}
	wg.Wait()
	return workerPanic
}

func collectParallelWorkersDeferredPins[W any](
	itemCount, numWorkers int,
	initWorker func(workerIdx int, w *W),
	runItem func(w *W, itemIdx int),
	afterItem func(w *W, itemIdx int),
	finishWorker func(workerIdx int, w *W),
	stopWorkerOnPanic bool,
) ([]any, any) {
	jsErrorPinsDeferred.Add(1)
	defer jsErrorPinsDeferred.Add(-1)
	return collectParallelWorkers(
		itemCount,
		numWorkers,
		initWorker,
		runItem,
		afterItem,
		finishWorker,
		stopWorkerOnPanic,
	)
}

func collectParallelWorkers[W any](
	itemCount, numWorkers int,
	initWorker func(workerIdx int, w *W),
	runItem func(w *W, itemIdx int),
	afterItem func(w *W, itemIdx int),
	finishWorker func(workerIdx int, w *W),
	stopWorkerOnPanic bool,
) ([]any, any) {
	panicsByItem := make([]any, itemCount)
	workerPanic := spawnIndexedWorkers(itemCount, numWorkers, func(wIdx int, claim func() (int, bool)) {
		var w W
		initWorker(wIdx, &w)
		if finishWorker != nil {
			defer finishWorker(wIdx, &w)
		}
		for {
			idx, ok := claim()
			if !ok {
				return
			}
			panicked := runParallelItemIndexed(&w, idx, panicsByItem, runItem)
			if afterItem != nil {
				afterItem(&w, idx)
			}
			if panicked && stopWorkerOnPanic {
				return
			}
		}
	})
	return panicsByItem, workerPanic
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
	numWorkers = min(numWorkers, itemCount)
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

	panicsByItem, workerPanic := collectParallelWorkers(
		itemCount,
		numWorkers,
		func(wIdx int, dst *W) { *dst = initWorker(wIdx) },
		runItem,
		nil,
		finishWorker,
		true,
	)
	for _, p := range panicsByItem {
		if p != nil {
			panic(p)
		}
	}
	if workerPanic != nil {
		panic(workerPanic)
	}
}

type parallelNodeTask[S any] struct {
	allowDetached           bool
	enableTypeRegistryReads bool
	prepare                 func()
	nodeAt                  func(idx int) rhino.Node
	initWorker              func(workerID int, rc *Compiler) S
	runItem                 func(state S, idx int, rc *Compiler)
	onItemPanic             func(state S)
	finishWorker            func(workerID int, state S)
}

func runParallelNodesByScript[S any](
	compiler *Compiler,
	itemCount int,
	task parallelNodeTask[S],
) bool {
	if itemCount < 2 || parallelThreadCount(compiler, itemCount) <= 1 {
		return false
	}
	bucketByScript := make(map[rhino.Node]int, min(itemCount, 64))
	var buckets [][]int
	var detached []int
	for idx := 0; idx < itemCount; idx++ {
		script := NodeUtil_GetEnclosingScript(task.nodeAt(idx))
		if script.IsNil() {
			if !task.allowDetached {
				return false
			}
			detached = append(detached, idx)
			continue
		}
		bIdx, exists := bucketByScript[script]
		if !exists {
			bIdx = len(buckets)
			bucketByScript[script] = bIdx
			buckets = append(buckets, []int{idx})
		} else {
			buckets[bIdx] = append(buckets[bIdx], idx)
		}
	}
	numWorkers := parallelThreadCount(compiler, len(buckets))
	if numWorkers <= 1 {
		return false
	}

	if task.enableTypeRegistryReads {
		prewarmCompilerForParallelReads(compiler)
		restoreReg := enableParallelTypeRegistryReads(compiler)
		defer restoreReg()
	}
	if task.prepare != nil {
		task.prepare()
	}

	eventsByItem := make([][]recordedCompilerEvent, itemCount)
	panicsByItem := make([]any, itemCount)

	runParallelWork(compiler, numWorkers, len(buckets), func(workerID int) (func(int, *recordingCompiler), func()) {
		var state S
		initialized := false
		var cleanup func()
		if task.finishWorker != nil {
			cleanup = func() {
				if initialized {
					task.finishWorker(workerID, state)
				}
			}
		}
		runBucket := func(bIdx int, rc *recordingCompiler) {
			if !initialized {
				if task.initWorker != nil {
					state = task.initWorker(workerID, rc)
				}
				initialized = true
			}
			for _, origIdx := range buckets[bIdx] {
				resetRecordingCompiler(rc, compiler)
				panicked := runParallelItemIndexed(rc, origIdx, panicsByItem, func(rc *recordingCompiler, itemIdx int) {
					task.runItem(state, itemIdx, rc)
				})
				eventsByItem[origIdx] = rc.takeEvents()
				if panicked {
					if task.onItemPanic != nil {
						task.onItemPanic(state)
					}
					return
				}
			}
		}
		return runBucket, cleanup
	}, false)

	if len(detached) > 0 {
		var mainRC recordingCompiler
		initRecordingCompiler(&mainRC, compiler)
		var mainState S
		if task.initWorker != nil {
			mainState = task.initWorker(0, &mainRC)
			if task.finishWorker != nil {
				defer task.finishWorker(0, mainState)
			}
		}
		for _, origIdx := range detached {
			resetRecordingCompiler(&mainRC, compiler)
			panicked := runParallelItemIndexed(&mainRC, origIdx, panicsByItem, func(rc *recordingCompiler, itemIdx int) {
				task.runItem(mainState, itemIdx, rc)
			})
			eventsByItem[origIdx] = mainRC.takeEvents()
			if panicked && task.onItemPanic != nil {
				task.onItemPanic(mainState)
			}
		}
	}

	for i := 0; i < itemCount; i++ {
		if len(eventsByItem[i]) != 0 {
			replayCompilerEvents(compiler, eventsByItem[i])
		}
		if p := panicsByItem[i]; p != nil {
			rethrowParallelPanic(compiler, p)
		}
	}
	return true
}

// runThrowableTask invokes fn and catches any jutil.Throwable panic, returning
// it as a value while re-panicking non-Throwable panics.
func runThrowableTask(fn func()) (failure jutil.Throwable) {
	defer func() {
		if r := recover(); r != nil {
			if t, ok := r.(jutil.Throwable); ok {
				failure = t
				return
			}
			panic(r)
		}
	}()
	fn()
	return nil
}

// parallelPrebuild runs every submitted task, retaining Throwable failures in
// index order. Each slot is written by one worker and consumed after the join.
func parallelPrebuild(count, workers int, run func(int)) []jutil.Throwable {
	return runParallelIndexedTasks(count, workers, func(_, index int) { run(index) })
}

func runParallelThrowableWorkers(
	count, workers int,
	workerFn func(workerIdx int, claim func() (int, bool), recordFailure func(int, jutil.Throwable)),
) []jutil.Throwable {
	failures := make([]jutil.Throwable, count)
	recordFailure := func(index int, t jutil.Throwable) {
		failures[index] = t
	}
	if workerPanic := spawnIndexedWorkers(count, workers, func(wIdx int, claim func() (int, bool)) {
		workerFn(wIdx, claim, recordFailure)
	}); workerPanic != nil {
		panic(workerPanic)
	}
	return failures
}

func runParallelIndexedTasks(count, workers int, run func(worker, index int)) []jutil.Throwable {
	return runParallelThrowableWorkers(count, workers, func(wIdx int, claim func() (int, bool), recordFailure func(int, jutil.Throwable)) {
		for {
			index, ok := claim()
			if !ok {
				return
			}
			if t := runThrowableTask(func() { run(wIdx, index) }); t != nil {
				recordFailure(index, t)
			}
		}
	})
}

func parallelPrebuildWorkers(
	count, workers int,
	workerFn func(claim func() (int, bool), recordFailure func(int, jutil.Throwable)),
) []jutil.Throwable {
	return runParallelThrowableWorkers(count, workers, func(_ int, claim func() (int, bool), recordFailure func(int, jutil.Throwable)) {
		workerFn(claim, recordFailure)
	})
}
