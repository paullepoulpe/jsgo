package jscomp

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type inlineFnCacheStatsSnapshot struct {
	hits       int64
	misses     int64
	leafSkips  int64
	seHits     int64
	fnInfoSeed int64
	costHits   int64
}

func snapshotInlineFnCacheStats() inlineFnCacheStatsSnapshot {
	return inlineFnCacheStatsSnapshot{
		hits:       inlineFnCacheCandidateHits.Load(),
		misses:     inlineFnCacheCandidateMisses.Load(),
		leafSkips:  inlineFnCacheLeafWalkSkips.Load(),
		seHits:     inlineFnCacheSideEffectHits.Load(),
		fnInfoSeed: inlineFnCacheFnInfoSeeds.Load(),
		costHits:   inlineFnCacheCostHits.Load(),
	}
}

func (after inlineFnCacheStatsSnapshot) sub(before inlineFnCacheStatsSnapshot) inlineFnCacheStatsSnapshot {
	return inlineFnCacheStatsSnapshot{
		hits:       after.hits - before.hits,
		misses:     after.misses - before.misses,
		leafSkips:  after.leafSkips - before.leafSkips,
		seHits:     after.seHits - before.seHits,
		fnInfoSeed: after.fnInfoSeed - before.fnInfoSeed,
		costHits:   after.costHits - before.costHits,
	}
}

func withInlineFnCacheConfig(t testing.TB, env string, force, verify bool) {
	t.Helper()
	savedEnv, savedForce, savedVerify := inlineFnCacheEnv, inlineFnCacheForce, inlineFnCacheVerify
	inlineFnCacheEnv, inlineFnCacheForce, inlineFnCacheVerify = env, force, verify
	t.Cleanup(func() {
		inlineFnCacheEnv, inlineFnCacheForce, inlineFnCacheVerify = savedEnv, savedForce, savedVerify
	})
}

func newNormalizedTestCompiler(t testing.TB, externCode string, scripts ...string) *Compiler {
	t.Helper()
	options := NewCompilerOptions()
	CompilationLevel_SIMPLE_OPTIMIZATIONS.SetOptionsForCompilationLevel(options)
	c := NewCompiler()
	var inputs []*SourceFile
	for i, code := range scripts {
		inputs = append(inputs, SourceFile_FromCode(fmt.Sprintf("in_%d.js", i), code))
	}
	c.Init(
		[]*SourceFile{SourceFile_FromCode("externs.js", externCode)},
		inputs,
		options,
	)
	c.ParseInputs()
	if c.HasErrors() {
		t.Fatalf("parse errors: %v", c.GetErrors())
	}
	NewNormalize(c, false).Process(c.externsRoot, c.jsRoot)
	c.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
	return c
}

func buildIIFEBundleSource(numSiblings int) string {
	var b strings.Builder
	b.WriteString("(function(global) {\n")
	for i := 0; i < numSiblings; i++ {
		switch i % 4 {
		case 0:
			// Multi-call block-bodied function with side effects + parameter reads (kept uninlined due to cost).
			fmt.Fprintf(&b,
				"  function sib_%d(a, b, c) { sink(a + %d); var t = a * b + c; sink(t + b); return t - c + %d; }\n",
				i, i, i+1)
		case 1:
			// Direct-return function referencing y once after side-effecting call (exercises calleeSideEffectSummary).
			fmt.Fprintf(&b,
				"  var sib_%d = function(x, y) { return (sink(x + %d) ? sink(%d) + %d : sink(%d) + %d) + y; };\n",
				i, i, i+1, i+2, i+3, i+4)
		case 2:
			// Function with inner closure + parameter (shallowLocalClosureDisallowed).
			fmt.Fprintf(&b,
				"  function sib_%d(k) { return function(v) { return k + v + %d; }; }\n",
				i, i)
		case 3:
			// Non-inlinable default parameter function.
			fmt.Fprintf(&b,
				"  function sib_%d(a = %d, b = %d) { return a + b; }\n",
				i, i, i+2)
		}
	}
	// Mutating target function that remains uninlined across iterations.
	b.WriteString("  function mutTarget(x, y) { var r = x + 100; sink(r); sink(y); return r + y; }\n")
	b.WriteString("  var acc = 0;\n")
	for i := 0; i < numSiblings; i++ {
		switch i % 4 {
		case 0:
			fmt.Fprintf(&b, "  acc += sib_%d(acc, %d, %d) + sib_%d(%d, acc, 1) + sib_%d(1, 2, %d);\n",
				i, i, i+1, i, i+2, i, i+3)
		case 2:
			fmt.Fprintf(&b, "  acc += sib_%d(%d)(1) + sib_%d(%d)(2);\n", i, i, i, i+1)
		case 1:
			if i%8 == 1 {
				fmt.Fprintf(&b, "  acc += sib_%d(acc, global.out) + sib_%d(%d, global.out) + sib_%d(1, global.out);\n",
					i, i, i+2, i)
			} else {
				fmt.Fprintf(&b, "  acc += sib_%d(1, 2) + sib_%d(3, 4) + sib_%d(5, 6);\n", i, i, i)
			}
		default:
			fmt.Fprintf(&b, "  acc += sib_%d(acc, %d) + sib_%d(%d, acc) + sib_%d(1, %d);\n",
				i, i+1, i, i+2, i, i+3)
		}
	}
	b.WriteString("  acc += mutTarget(acc, 1) + mutTarget(2, acc) + mutTarget(3, 4);\n")
	b.WriteString("  global.out = acc;\n")
	b.WriteString("})(this);\n")
	return b.String()
}

func findNamedFunctionNode(root rhino.Node, targetName string) rhino.Node {
	var found rhino.Node
	var walk func(n rhino.Node)
	walk = func(n rhino.Node) {
		if !found.IsNil() || n.IsNil() {
			return
		}
		if n.IsFunction() {
			if first := n.GetFirstChild(); first.IsName() && first.GetString() == targetName {
				found = n
				return
			}
			if p := n.GetParent(); p.IsName() && p.GetString() == targetName {
				found = n
				return
			}
		}
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			walk(c)
		}
	}
	walk(root)
	return found
}

// Test 1: Verifies cross-iteration cache hit reuse on an IIFE bundle where
// 1 sibling function mutates per PhaseOptimizer iteration while hundreds of
// sibling functions and the outer IIFE wrapper hit the cache, and verifies
// post-optimization cleanup (c.inlineFnCache == nil).
func TestInlineFunctionsCache_IIFEBundleReuseAndCleanup(t *testing.T) {
	withInlineFnCacheConfig(t, "", false, true)

	const externs = "function sink(x) {}\n"
	const numSiblings = 240
	src := buildIIFEBundleSource(numSiblings)

	c := newNormalizedTestCompiler(t, externs, src)
	mutFn := findNamedFunctionNode(c.jsRoot, "mutTarget")
	if mutFn.IsNil() {
		t.Fatal("mutTarget function not found")
	}

	var iterStats [3]inlineFnCacheStatsSnapshot
	runCount := 0
	mutatedAtRun1Before := false
	mutatedAtRun1After := false

	mutateMutTarget := func(comp AbstractCompiler) {
		var mutated bool
		var mutateNum func(n rhino.Node)
		mutateNum = func(n rhino.Node) {
			if mutated || n.IsNil() {
				return
			}
			if n.IsNumber() {
				n.SetDouble(n.GetDouble() + 1)
				comp.ReportChangeToEnclosingScope(n)
				mutated = true
				return
			}
			for ch := n.GetFirstChild(); !ch.IsNil(); ch = ch.GetNext() {
				mutateNum(ch)
			}
		}
		mutateNum(mutFn.GetLastChild())
		if !mutated {
			t.Fatal("expected to mutate number literal inside mutTarget")
		}
	}

	c.options.SetUseSizeHeuristicToStopOptimizationLoop(false)
	po := c.createPhaseOptimizer()
	c.phaseOptimizer = po
	po.SetValidityCheck(compiler_validityCheck)
	loop := po.AddFixedPointLoop()
	loop.AddLoopedPass(NewLoopablePassFactory("mutateBefore", func(comp AbstractCompiler) CompilerPass {
		return CompilerPassFunc(func(_, _ rhino.Node) {
			if runCount == 1 && !mutatedAtRun1Before {
				mutatedAtRun1Before = true
				mutateMutTarget(comp)
			}
		})
	}))
	loop.AddLoopedPass(NewLoopablePassFactory(PassNames_INLINE_FUNCTIONS, func(comp AbstractCompiler) CompilerPass {
		return CompilerPassFunc(func(ext, root rhino.Node) {
			before := snapshotInlineFnCacheStats()
			NewInlineFunctions(
				comp,
				comp.GetUniqueNameIdSupplier(),
				CompilerOptions_Reach_ALL,
				false,
				CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING,
			).Process(ext, root)
			after := snapshotInlineFnCacheStats()
			if runCount < len(iterStats) {
				iterStats[runCount] = after.sub(before)
			}
			runCount++
		})
	}))
	loop.AddLoopedPass(NewLoopablePassFactory("mutateAfter", func(comp AbstractCompiler) CompilerPass {
		return CompilerPassFunc(func(_, _ rhino.Node) {
			if runCount == 1 && !mutatedAtRun1After {
				mutatedAtRun1After = true
				mutateMutTarget(comp)
			} else if runCount == 2 && mutatedAtRun1After {
				mutatedAtRun1After = false
				mutateMutTarget(comp)
			}
		})
	}))

	po.Process(c.externsRoot, c.jsRoot)

	if runCount != 3 {
		t.Fatalf("expected 3 loop iterations, got %d", runCount)
	}
	if c.inlineFnCache == nil {
		t.Fatal("expected inlineFnCache to be populated during PhaseOptimizer loop")
	}
	// Iteration 0: cold cache — all candidates miss, 0 hits.
	if iterStats[0].hits != 0 || iterStats[0].misses < numSiblings+1 {
		t.Fatalf("iter 0 stats = %+v; expected 0 hits and >= %d misses", iterStats[0], numSiblings+1)
	}
	// Iterations 1 & 2: all 240 unchanged siblings + the outer IIFE wrapper MUST hit,
	// and ONLY mutTarget (1 function) should miss!
	for iter := 1; iter <= 2; iter++ {
		st := iterStats[iter]
		if st.misses != 1 {
			t.Fatalf("iter %d misses = %d, want exactly 1 (only mutTarget); full stats = %+v", iter, st.misses, st)
		}
		if st.hits < numSiblings+1 {
			t.Fatalf("iter %d hits = %d, want >= %d (all siblings + outer IIFE); full stats = %+v", iter, st.hits, numSiblings+1, st)
		}
		if st.seHits <= 0 {
			t.Fatalf("iter %d expected sideEffectHits > 0, got %+v", iter, st)
		}
		if st.fnInfoSeed <= 0 {
			t.Fatalf("iter %d expected fnInfoSeed > 0 (outer IIFE caller seeded), got %+v", iter, st)
		}
		if st.costHits <= 0 {
			t.Fatalf("iter %d expected costHits > 0, got %+v", iter, st)
		}
	}

	// Verify leafWalkSkips when verify mode is off.
	withInlineFnCacheConfig(t, "", true, false)
	beforeLeaf := snapshotInlineFnCacheStats()
	NewInlineFunctions(
		c,
		c.GetUniqueNameIdSupplier(),
		CompilerOptions_Reach_ALL,
		false,
		CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING,
	).Process(c.externsRoot, c.jsRoot)
	deltaLeaf := snapshotInlineFnCacheStats().sub(beforeLeaf)
	if deltaLeaf.leafSkips < int64(numSiblings*3/4) {
		t.Fatalf("expected leafSkips >= %d on warm pass, got %+v", numSiblings*3/4, deltaLeaf)
	}

	// Verify end-of-optimizations drops c.inlineFnCache.
	c.phaseOptimizer = nil
	c.maybeCompactAST("end of optimizations", nil, nil, nil)
	if c.inlineFnCache != nil {
		t.Fatal("expected maybeCompactAST('end of optimizations') to clear inlineFnCache")
	}

	// Verify full end-to-end Compile and SaveState also leave c.inlineFnCache == nil.
	withInlineFnCacheConfig(t, "", false, true)
	fullOpts := NewCompilerOptions()
	CompilationLevel_SIMPLE_OPTIMIZATIONS.SetOptionsForCompilationLevel(fullOpts)
	fullCompiler := NewCompiler()
	res := fullCompiler.Compile(
		[]*SourceFile{SourceFile_FromCode("externs.js", externs)},
		[]*SourceFile{SourceFile_FromCode("bundle.js", src)},
		fullOpts,
	)
	if !res.Success {
		t.Fatalf("full Compile failed: %v", res.Errors)
	}
	if fullCompiler.inlineFnCache != nil {
		t.Fatal("expected fullCompiler.inlineFnCache == nil after Compile completes")
	}

	// Verify NewTracer outside loop (as used by SaveState) drops inlineFnCache.
	fullCompiler.inlineFnCache = &inlineFunctionsCache{
		index: make(map[rhino.Node]uint32),
	}
	var stateBuf bytes.Buffer
	fullCompiler.SaveState(&stateBuf)
	if fullCompiler.inlineFnCache != nil {
		t.Fatal("expected SaveState to clear inlineFnCache")
	}
}

// Test 2: Verifies exact compositional cache invalidation across:
//   - Deep inner closure mutation (non-candidate nested arrow inside candidate)
//   - Parameter mutation (adding default param) and body mutation (modifying param)
//   - Parent @noinline JSDoc added on VAR and NAME without mutating FUNCTION
//   - Declaring block let/const insertion (CAPTURED_VARIABLE_NOT_IN_SCOPE)
//   - Parent VAR rename turning function recursive (summaryFnNameID mismatch)
//   - HasRegExpGlobalReferences change invalidating cached sideEffects
//   - compactAST node relocation clearing c.inlineFnCache safely
func TestInlineFunctionsCache_CompositionalInvalidations(t *testing.T) {
	const externs = "function use(a, b) {}\n"

	runInlinePass := func(c *Compiler) inlineFnCacheStatsSnapshot {
		before := snapshotInlineFnCacheStats()
		NewInlineFunctions(
			c,
			c.GetUniqueNameIdSupplier(),
			CompilerOptions_Reach_ALL,
			false,
			CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING,
		).Process(c.externsRoot, c.jsRoot)
		return snapshotInlineFnCacheStats().sub(before)
	}

	t.Run("InnerClosureMutationInvalidatesOuterSubtree", func(t *testing.T) {
		withInlineFnCacheConfig(t, "", true, true)
		// outer has 0 params and 0 shallow declarations, returning nested arrows.
		// Initially called twice with large body so it stays uninlined by cost on pass 1,
		// then we remove one call site AND mutate the deepest arrow to reference `eval`.
		src := `
var outer = function() {
  return () => () => use(1, 2) + use(3, 4) + use(5, 6);
};
var a = outer();
var b = outer();
use(a, b);
`
		c := newNormalizedTestCompiler(t, externs, src)
		outerFn := findNamedFunctionNode(c.jsRoot, "outer")
		if outerFn.IsNil() {
			t.Fatal("outer not found")
		}
		outerTimeBefore := outerFn.GetChangeTime()
		st1 := runInlinePass(c)
		if st1.misses == 0 {
			t.Fatalf("expected cold miss, got %+v", st1)
		}

		// Warm pass without mutation: outer must hit cache!
		st2 := runInlinePass(c)
		if st2.hits != 1 || st2.misses != 0 {
			t.Fatalf("expected 1 hit / 0 misses on warm pass, got %+v", st2)
		}

		// Now mutate ONLY the deepest arrow to call eval("1"), advancing ONLY deepestArrow's changeTime!
		// After Normalize, arrow bodies are BLOCK -> RETURN -> expr.
		ret1 := outerFn.GetLastChild().GetFirstChild()
		midArrow := ret1.GetFirstChild()
		deepArrow := midArrow.GetLastChild().GetFirstChild().GetFirstChild()
		if !deepArrow.IsFunction() {
			t.Fatalf("expected deepArrow FUNCTION, got %v", deepArrow)
		}
		evalRet := rhino.NewNodeWithChild(rhino.Token_RETURN, rhino.IR_Call(rhino.IR_Name("eval"), rhino.IR_String("1")))
		deepArrow.GetLastChild().GetFirstChild().ReplaceWith(evalRet)
		c.changeTracker.IncrementChangeStamp()
		c.ReportChangeToEnclosingScope(evalRet)

		if outerFn.GetChangeTime() != outerTimeBefore {
			t.Fatalf("expected outerFn own ChangeTime to remain %d, got %d", outerTimeBefore, outerFn.GetChangeTime())
		}

		// Remove the second call `var b = outer();` so outer has a single call site
		// and WOULD be inlined if the stale cache entry were wrongly reused!
		script := c.jsRoot.GetFirstChild()
		for stmt := script.GetFirstChild(); !stmt.IsNil(); stmt = stmt.GetNext() {
			if stmt.IsVar() && stmt.GetFirstChild().GetString() == "b" {
				replacement := rhino.IR_VarWithValue(rhino.IR_Name("b"), rhino.IR_Number(0))
				stmt.ReplaceWith(replacement)
				c.ReportChangeToEnclosingScope(replacement)
				break
			}
		}

		st3 := runInlinePass(c)
		if st3.misses != 1 || st3.hits != 0 {
			t.Fatalf("expected inner arrow mutation to invalidate outer (1 miss, 0 hits), got %+v", st3)
		}
		if got := findNamedFunctionNode(c.jsRoot, "outer"); got.IsNil() {
			t.Fatal("outer was wrongly inlined despite deepest arrow referencing eval!")
		}
	})

	t.Run("ParamAndBodyMutationsInvalidate", func(t *testing.T) {
		withInlineFnCacheConfig(t, "", true, true)
		src := `
function fParam(x) { use(x, 100); use(x, 200); return use(x, 300) + use(x, 400); }
function fBody(y) { use(y, 100); use(y, 200); return use(y, 300) + use(y, 400); }
use(fParam(10), fParam(20));
use(fBody(30), fBody(40));
`
		c := newNormalizedTestCompiler(t, externs, src)
		runInlinePass(c)
		if st := runInlinePass(c); st.hits != 2 || st.misses != 0 {
			t.Fatalf("expected 2 warm hits, got %+v", st)
		}

		// 1. Mutate fParam's parameter x -> x = 1 (non-inlinable default value).
		fParam := findNamedFunctionNode(c.jsRoot, "fParam")
		paramList := fParam.GetSecondChild()
		xParam := paramList.GetFirstChild().Detach()
		defParam := rhino.NewNodeWith2Children(rhino.Token_DEFAULT_VALUE, xParam, rhino.IR_Number(1))
		paramList.AddChildToFront(defParam)
		c.changeTracker.IncrementChangeStamp()
		c.ReportChangeToEnclosingScope(defParam)

		if st := runInlinePass(c); st.hits != 1 || st.misses != 1 {
			t.Fatalf("expected 1 hit (fBody) and 1 miss (fParam), got %+v", st)
		}

		// 2. Mutate fBody's body to reassign parameter y (`y = 99; ...`).
		fBody := findNamedFunctionNode(c.jsRoot, "fBody")
		body := fBody.GetLastChild()
		assignStmt := rhino.IR_ExprResult(rhino.IR_Assign(rhino.IR_Name("y"), rhino.IR_Number(99)))
		body.AddChildToFront(assignStmt)
		c.changeTracker.IncrementChangeStamp()
		c.ReportChangeToEnclosingScope(assignStmt)

		if st := runInlinePass(c); st.hits != 1 || st.misses != 1 {
			t.Fatalf("expected 1 hit (fParam shallow hit) and 1 miss (fBody), got %+v", st)
		}
	})

	t.Run("NoInlineJSDocOnVarAndNameRespectedWithoutFnMutation", func(t *testing.T) {
		for _, onName := range []bool{false, true} {
			withInlineFnCacheConfig(t, "", true, true)
			src := `
var target = function(x) { use(x, 100); use(x, 200); return use(x, 300) + use(x, 400); };
var keep1 = target(1);
var keep2 = target(2);
use(keep1, keep2);
`
			c := newNormalizedTestCompiler(t, externs, src)
			runInlinePass(c) // warm cache

			targetFn := findNamedFunctionNode(c.jsRoot, "target")
			if targetFn.IsNil() {
				t.Fatal("target was unexpectedly inlined on warm pass")
			}
			nameNode := targetFn.GetParent()
			varNode := nameNode.GetParent()
			fnTimeBefore := targetFn.GetChangeTime()

			docBuilder := rhino.NewJSDocInfo_Builder()
			docBuilder.RecordNoInline()
			if onName {
				nameNode.SetJSDocInfo(docBuilder.Build())
			} else {
				varNode.SetJSDocInfo(docBuilder.Build())
			}

			// Drop second call so target would be inlined if @noinline were missed!
			script := c.jsRoot.GetFirstChild()
			for stmt := script.GetFirstChild(); !stmt.IsNil(); stmt = stmt.GetNext() {
				if stmt.IsVar() && stmt.GetFirstChild().GetString() == "keep2" {
					rep := rhino.IR_VarWithValue(rhino.IR_Name("keep2"), rhino.IR_Number(0))
					stmt.ReplaceWith(rep)
					c.changeTracker.IncrementChangeStamp()
					c.ReportChangeToEnclosingScope(rep)
					break
				}
			}
			if targetFn.GetChangeTime() != fnTimeBefore {
				t.Fatalf("expected targetFn ChangeTime unchanged")
			}

			runInlinePass(c)
			if findNamedFunctionNode(c.jsRoot, "target").IsNil() {
				t.Fatalf("onName=%v: target was inlined despite @noinline JSDoc!", onName)
			}
		}
	})

	t.Run("LetConstInDeclaringBlockInvalidatesVarCandidate", func(t *testing.T) {
		withInlineFnCacheConfig(t, "", true, true)
		src := `
{
  var blkFn = function(x) { use(x, 100); use(x, 200); return use(x, 300) + use(x, 400); };
  var r1 = blkFn(1);
  var r2 = blkFn(2);
  use(r1, r2);
}
`
		c := newNormalizedTestCompiler(t, externs, src)
		runInlinePass(c) // warm cache: blkFn is inlinable (block has no let/const yet)

		blkFn := findNamedFunctionNode(c.jsRoot, "blkFn")
		if blkFn.IsNil() {
			t.Fatal("blkFn was unexpectedly inlined on warm pass")
		}
		fnTimeBefore := blkFn.GetChangeTime()
		declBlock := blkFn.GetGrandparent().GetParent()

		// Insert `let blockScoped = 1;` into declBlock AND remove second call `blkFn(2)`
		// so blkFn WOULD be inlined if declaring-block let/const were cached on fnNode!
		letStmt := rhino.IR_LetWithValue(rhino.IR_Name("blockScoped"), rhino.IR_Number(1))
		declBlock.AddChildToFront(letStmt)
		for stmt := declBlock.GetFirstChild(); !stmt.IsNil(); stmt = stmt.GetNext() {
			if stmt.IsVar() && stmt.GetFirstChild().GetString() == "r2" {
				rep := rhino.IR_VarWithValue(rhino.IR_Name("r2"), rhino.IR_Number(0))
				stmt.ReplaceWith(rep)
				break
			}
		}
		c.changeTracker.IncrementChangeStamp()
		c.ReportChangeToEnclosingScope(letStmt)

		if blkFn.GetChangeTime() != fnTimeBefore {
			t.Fatalf("expected blkFn ChangeTime unchanged")
		}

		st := runInlinePass(c)
		if st.hits != 1 {
			t.Fatalf("expected cache hit with uncached parent block check, got %+v", st)
		}
		if findNamedFunctionNode(c.jsRoot, "blkFn").IsNil() {
			t.Fatal("blkFn was wrongly inlined despite let declaration in declaring block!")
		}
	})

	t.Run("VarRenameTurningRecursiveForcesMiss", func(t *testing.T) {
		withInlineFnCacheConfig(t, "", true, true)
		src := `
function rec() { use(10, 20); use(30, 40); return use(50, 60); }
var orig = function(x) { use(x, 100); return rec() + rec() + use(x, 200); };
use(orig(1), orig(2));
`
		c := newNormalizedTestCompiler(t, externs, src)
		runInlinePass(c) // warm cache under name "orig"

		origFn := findNamedFunctionNode(c.jsRoot, "orig")
		if origFn.IsNil() {
			t.Fatal("orig was unexpectedly inlined on warm pass")
		}
		fnTimeBefore := origFn.GetChangeTime()
		// Rename parent VAR name from "orig" to "rec" without mutating origFn's subtree!
		// Remove `function rec()` first and rename `var orig` -> `var rec` with single call `rec(1)`.
		script := c.jsRoot.GetFirstChild()
		recDecl := script.GetFirstChild()
		recDecl.Detach()
		NodeUtil_MarkFunctionsDeleted(recDecl, c)

		origNameNode := origFn.GetParent()
		newNameNode := rhino.IR_Name("rec")
		origFn.Detach()
		newNameNode.AddChildToFront(origFn)
		origNameNode.ReplaceWith(newNameNode)

		lastCallStmt := script.GetLastChild()
		repCall := rhino.IR_ExprResult(rhino.IR_Call(rhino.IR_Name("use"),
			rhino.IR_Call(rhino.IR_Name("rec"), rhino.IR_Number(1)),
			rhino.IR_Number(0)))
		lastCallStmt.ReplaceWith(repCall)
		c.changeTracker.IncrementChangeStamp()
		c.ReportChangeToEnclosingScope(repCall)

		if origFn.GetChangeTime() != fnTimeBefore {
			t.Fatalf("expected origFn own ChangeTime unchanged")
		}

		st := runInlinePass(c)
		if st.misses != 1 || st.hits != 0 {
			t.Fatalf("expected summaryFnNameID mismatch to trigger cache miss, got %+v", st)
		}
		if findNamedFunctionNode(c.jsRoot, "rec").IsNil() {
			t.Fatal("recursive renamed VAR function was wrongly inlined!")
		}
	})

	t.Run("HasRegExpGlobalReferencesInvalidatesSideEffectSummary", func(t *testing.T) {
		withInlineFnCacheConfig(t, "", true, true)
		src := `
function rxCallee(val) {
  return (/abc/.test("abc") ? use(1, 2) + use(3, 4) + use(5, 6) : use(7, 8) + use(9, 10)) + val;
}
var o = {p: 1};
use(rxCallee(o.p), rxCallee(o.p) + rxCallee(o.p));
`
		c := newNormalizedTestCompiler(t, externs, src)
		c.SetHasRegExpGlobalReferences(true)
		runInlinePass(c) // computes sideEffects with hasRegExpGlobalReferences=true

		stWarmSame := runInlinePass(c)
		if stWarmSame.seHits != 1 {
			t.Fatalf("expected 1 sideEffectHit when config unchanged, got %+v", stWarmSame)
		}

		// Flip HasRegExpGlobalReferences -> false; verify mode checks sideEffects parity!
		c.SetHasRegExpGlobalReferences(false)
		stFlipped := runInlinePass(c)
		if stFlipped.hits != 1 || stFlipped.seHits != 0 {
			t.Fatalf("expected candidate hit with 0 stale sideEffectHits after RegExp flag flip, got %+v", stFlipped)
		}
	})

	t.Run("ASTCompactionClearsCacheAndRebuildsSafely", func(t *testing.T) {
		withInlineFnCacheConfig(t, "", true, true)
		src := buildIIFEBundleSource(40)
		c := newNormalizedTestCompiler(t, "function sink(x) {}\n", src)
		c.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_COLORS_AND_SIMPLIFIED_JSDOC)
		// Normalized is required by InlineFunctions; note NORMALIZED >= COLORS_AND_SIMPLIFIED_JSDOC.
		c.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
		c.beginASTCompaction()

		st1 := runInlinePass(c)
		if c.inlineFnCache == nil || st1.misses == 0 {
			t.Fatalf("expected populated inlineFnCache before compaction, got %+v", st1)
		}

		// Force a full pre-order AST compaction that relocates live node IDs.
		c.astCompact.enabled = true
		c.astCompact.preorder = true
		c.compactAST("forced-test-compact", true, &c.externsRoot, &c.jsRoot, nil)
		if c.inlineFnCache != nil {
			t.Fatal("expected compactAST relocation to clear c.inlineFnCache")
		}

		// Subsequent pass must rebuild cache cleanly on relocated node IDs, and the pass after that must hit!
		st2 := runInlinePass(c)
		if st2.hits != 0 || st2.misses == 0 {
			t.Fatalf("expected cold rebuild after compaction, got %+v", st2)
		}
		st3 := runInlinePass(c)
		if st3.hits == 0 || st3.misses != 0 {
			t.Fatalf("expected 100%% cache hits after post-compaction rebuild, got %+v", st3)
		}
	})
}

// Test 3: Verifies 100% byte-identical output and identical ValidityCheck /
// ChangeVerifier behavior between JSGO_INLINE_FN_CACHE=0 and =1 across
// numParallelThreads = 1 and 4 on both single-script IIFE bundles and
// multi-script programs (including forced mid-loop AST compaction).
func TestInlineFunctionsCache_ByteIdenticalSerialAndParallel(t *testing.T) {
	const externs = `
function sink(x) {}
/** @constructor */ function Box(v) { this.v = v; }
`
	singleBundle := []*SourceFile{
		SourceFile_FromCode("bundle.js", buildIIFEBundleSource(60)+`
(function() {
  function stepA(x) { return x + 1; }
  function stepB(x) { var v = stepA(x); return v * 2; }
  function stepC(x) { return stepB(x) + stepA(3); }
  sink(stepC(10));
})();
`),
	}

	var multiScripts []*SourceFile
	for s := 0; s < 4; s++ {
		var sb strings.Builder
		sb.WriteString("(function(global) {\n")
		for i := 0; i < 12; i++ {
			fmt.Fprintf(&sb,
				"  function m_%d_%d(a, b) { var r = a * b + %d; sink(r); return r - b + %d; }\n",
				s, i, i+1, s+1)
		}
		fmt.Fprintf(&sb,
			"  function stepA_%d(z) { return z + %d; }\n", s, s+1)
		fmt.Fprintf(&sb,
			"  function stepB_%d(z) { var q = stepA_%d(z); return q * 2; }\n", s, s)
		fmt.Fprintf(&sb,
			"  function caller_%d(v) { var t = stepB_%d(v) + stepA_%d(3); return m_%d_0(t, 1) + m_%d_0(2, t) + m_%d_1(t, 3) + m_%d_1(4, t); }\n",
			s, s, s, s, s, s, s)
		for i := 2; i < 12; i++ {
			fmt.Fprintf(&sb, "  sink(m_%d_%d(1, 2) + m_%d_%d(3, 4) + m_%d_%d(5, 6));\n", s, i, s, i, s, i)
		}
		fmt.Fprintf(&sb, "  global['out_%d'] = caller_%d(%d) + caller_%d(%d) + caller_%d(%d);\n", s, s, s, s, s+5, s, s+9)
		sb.WriteString("})(this);\n")
		multiScripts = append(multiScripts, SourceFile_FromCode(fmt.Sprintf("s_%d.js", s), sb.String()))
	}

	cases := []struct {
		name      string
		inputs    []*SourceFile
		advanced  bool
		maxFnSize int
		compactAt string
	}{
		{
			name:      "SingleScriptBundle_SimpleDefault",
			inputs:    singleBundle,
			advanced:  false,
			maxFnSize: CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING,
		},
		{
			name:      "SingleScriptBundle_AdvancedWithCompaction",
			inputs:    singleBundle,
			advanced:  true,
			maxFnSize: CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING,
			compactAt: "loop2.1",
		},
		{
			name:      "MultiScript_SimpleDefault",
			inputs:    multiScripts,
			advanced:  false,
			maxFnSize: CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING,
		},
		{
			name:      "MultiScript_AdvancedWithMaxSizeCap",
			inputs:    multiScripts,
			advanced:  true,
			maxFnSize: 75,
		},
	}

	compileOne := func(t *testing.T, tcInputs []*SourceFile, advanced bool, maxFnSize, threads int, compactAt string) string {
		t.Helper()
		if compactAt != "" {
			savedForce := astLayoutForce
			astLayoutForce = []string{compactAt}
			defer func() { astLayoutForce = savedForce }()
		}
		opts := NewCompilerOptions()
		if advanced {
			CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(opts)
		} else {
			CompilationLevel_SIMPLE_OPTIMIZATIONS.SetOptionsForCompilationLevel(opts)
		}
		opts.SetDevMode(CompilerOptions_DevMode_EVERY_PASS)
		opts.SetNumParallelThreads(threads)
		if maxFnSize != CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING {
			opts.SetMaxFunctionSizeAfterInlining(maxFnSize)
		}

		// Clone SourceFiles so each compilation gets fresh inputs.
		freshInputs := make([]*SourceFile, len(tcInputs))
		for i, sf := range tcInputs {
			freshInputs[i] = SourceFile_FromCode(sf.GetName(), sf.GetCode())
		}

		c := NewCompiler()
		res := c.Compile(
			[]*SourceFile{SourceFile_FromCode("externs.js", externs)},
			freshInputs,
			opts,
		)
		if !res.Success {
			t.Fatalf("compile failed: %v", res.Errors)
		}
		if c.inlineFnCache != nil {
			t.Fatal("expected inlineFnCache to be nil after Compile")
		}
		return c.ToSource()
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			withInlineFnCacheConfig(t, "0", false, false)
			baseline1 := compileOne(t, tc.inputs, tc.advanced, tc.maxFnSize, 1, tc.compactAt)
			baseline4 := compileOne(t, tc.inputs, tc.advanced, tc.maxFnSize, 4, tc.compactAt)
			if baseline1 != baseline4 {
				t.Fatalf("uncached baseline differs between 1 and 4 threads")
			}

			for _, threads := range []int{1, 4} {
				// Test both verified mode and fast leaf-skipping default mode.
				for _, verify := range []bool{true, false} {
					withInlineFnCacheConfig(t, "", false, verify)
					before := snapshotInlineFnCacheStats()
					got := compileOne(t, tc.inputs, tc.advanced, tc.maxFnSize, threads, tc.compactAt)
					delta := snapshotInlineFnCacheStats().sub(before)
					if got != baseline1 {
						t.Fatalf("threads=%d verify=%v output mismatch!\ngot:\n%s\nwant:\n%s",
							threads, verify, got, baseline1)
					}
					if tc.compactAt == "" && delta.hits <= 0 {
						t.Fatalf("threads=%d verify=%v expected cache hits > 0, got %+v",
							threads, verify, delta)
					}
				}
			}
		})
	}
}

func benchmarkIIFEBundleIterations(b *testing.B, numFns int, cacheEnv string) {
	withInlineFnCacheConfig(b, cacheEnv, false, false)
	const externs = "function sink(x) {}\n"
	src := buildIIFEBundleSource(numFns)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		c := newNormalizedTestCompiler(b, externs, src)
		c.options.SetUseSizeHeuristicToStopOptimizationLoop(false)
		mutFn := findNamedFunctionNode(c.jsRoot, "mutTarget")
		po := c.createPhaseOptimizer()
		c.phaseOptimizer = po
		loop := po.AddFixedPointLoop()
		inlineRuns := 0
		mutatedBefore := 0
		mutatedAfter := 0
		mutateNum := func(comp AbstractCompiler) {
			num := mutFn.GetLastChild().GetFirstChild().GetFirstChild().GetFirstChild().GetSecondChild()
			num.SetDouble(num.GetDouble() + 1)
			comp.ReportChangeToEnclosingScope(num)
		}
		loop.AddLoopedPass(NewLoopablePassFactory("mutateBefore", func(comp AbstractCompiler) CompilerPass {
			return CompilerPassFunc(func(_, _ rhino.Node) {
				if (inlineRuns == 1 && mutatedBefore == 0) || (inlineRuns == 3 && mutatedBefore == 1) {
					mutatedBefore++
					mutateNum(comp)
				}
			})
		}))
		loop.AddLoopedPass(NewLoopablePassFactory(PassNames_INLINE_FUNCTIONS, func(comp AbstractCompiler) CompilerPass {
			return CompilerPassFunc(func(ext, root rhino.Node) {
				NewInlineFunctions(
					comp,
					comp.GetUniqueNameIdSupplier(),
					CompilerOptions_Reach_ALL,
					false,
					CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING,
				).Process(ext, root)
				inlineRuns++
			})
		}))
		loop.AddLoopedPass(NewLoopablePassFactory("mutateAfter", func(comp AbstractCompiler) CompilerPass {
			return CompilerPassFunc(func(_, _ rhino.Node) {
				if (inlineRuns == 1 && mutatedAfter == 0) || (inlineRuns == 2 && mutatedAfter == 1) {
					mutatedAfter++
					mutateNum(comp)
				}
			})
		}))
		b.StartTimer()
		po.Process(c.externsRoot, c.jsRoot)
	}
}

func benchmarkIIFEBundleWarmPass(b *testing.B, numFns int, cacheEnv string) {
	withInlineFnCacheConfig(b, cacheEnv, false, false)
	const externs = "function sink(x) {}\n"
	src := buildIIFEBundleSource(numFns)

	c := newNormalizedTestCompiler(b, externs, src)
	c.phaseOptimizer = &PhaseOptimizer{inLoop: true}
	mutFn := findNamedFunctionNode(c.jsRoot, "mutTarget")

	// Warm iteration 1 once before timing subsequent iterations where 1 function mutates.
	NewInlineFunctions(
		c,
		c.GetUniqueNameIdSupplier(),
		CompilerOptions_Reach_ALL,
		false,
		CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING,
	).Process(c.externsRoot, c.jsRoot)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		num := mutFn.GetLastChild().GetFirstChild().GetFirstChild().GetFirstChild().GetSecondChild()
		num.SetDouble(num.GetDouble() + 1)
		c.changeTracker.IncrementChangeStamp()
		c.ReportChangeToEnclosingScope(num)

		NewInlineFunctions(
			c,
			c.GetUniqueNameIdSupplier(),
			CompilerOptions_Reach_ALL,
			false,
			CompilerOptions_UNLIMITED_FUN_SIZE_AFTER_INLINING,
		).Process(c.externsRoot, c.jsRoot)
	}
}

func BenchmarkInlineFunctionsCache_IIFEBundle2000_Uncached(b *testing.B) {
	benchmarkIIFEBundleIterations(b, 2000, "0")
}

func BenchmarkInlineFunctionsCache_IIFEBundle2000_Cached(b *testing.B) {
	benchmarkIIFEBundleIterations(b, 2000, "1")
}

func BenchmarkInlineFunctionsCache_IIFEBundle5000_Uncached(b *testing.B) {
	benchmarkIIFEBundleIterations(b, 5000, "0")
}

func BenchmarkInlineFunctionsCache_IIFEBundle5000_Cached(b *testing.B) {
	benchmarkIIFEBundleIterations(b, 5000, "1")
}

func BenchmarkInlineFunctionsCache_IIFEBundle5000_WarmPass_Uncached(b *testing.B) {
	benchmarkIIFEBundleWarmPass(b, 5000, "0")
}

func BenchmarkInlineFunctionsCache_IIFEBundle5000_WarmPass_Cached(b *testing.B) {
	benchmarkIIFEBundleWarmPass(b, 5000, "1")
}
