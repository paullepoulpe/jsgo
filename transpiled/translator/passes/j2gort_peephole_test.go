package passes_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/transpiled/translator/passes"
)

func TestJ2gortPeephole_HoistIfConditionClosureIIFEs(t *testing.T) {
	tmpDir := t.TempDir()
	mustWrite := func(rel, content string) {
		t.Helper()
		full := filepath.Join(tmpDir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	mustWrite("go.mod", "module example.com/w30bpeepholetest\n\ngo 1.24\n")
	mustWrite("dfa/dfa.go", `package dfa

type stepState struct {
	stepCount int32
}

// Generic receiver method matching DataFlowAnalysisBase[N, L].Analyze() in gen/jscomp/data_flow_analysis.go!
type DataFlowAnalysisBase[N comparable, L any] struct {
	maxSteps int32
	states   []*stepState
	zeroNode N
}

func (this_ *DataFlowAnalysisBase[N, L]) Analyze(limit int32) (tripped bool, finalSteps int32) {
	curState := this_.states[0]
	for curState.stepCount <= limit+2 {
		if func() int32 {
			r_12 := curState
			old_13 := r_12.stepCount
			nw_14 := old_13 + 1
			r_12.stepCount = nw_14
			return old_13
		}() > limit {
			tripped = true
			break
		}
	}
	return tripped, curState.stepCount
}

// Generic type parameter return type N in IIFE inside generic method!
func (this_ *DataFlowAnalysisBase[N, L]) CheckGenericReturn(candidate N) bool {
	steps := int32(0)
	if func() N {
		steps = steps + 1
		return candidate
	}() == this_.zeroNode {
		return steps == 1
	}
	return false
}

// Lexical hygiene / zero-shadowing + short-circuit + unary ! + else-if + nested if verification.
func EvaluateComplexConditions(start int32, allowRhs bool) (log []string, outerR12 int32) {
	r_12 := int32(999)
	counter := start

	rhsProbe := func() bool {
		log = append(log, "rhs")
		return allowRhs && r_12 == 999
	}

	// 1. Short-circuit && with LHS IIFE declaring its own local r_12 (MUST NOT shadow outer r_12!)
	if (func() bool {
		r_12 := &counter
		*r_12 = *r_12 + 1
		log = append(log, "iife1")
		return *r_12 > 10
	}()) && rhsProbe() {
		log = append(log, "branch1-then")
		if !(func() bool {
			r_12 := &counter
			*r_12 = *r_12 + 5
			return *r_12 > 100
		}()) {
			log = append(log, "nested-not-then")
		}
	} else if func() int32 {
		r_12 := &counter
		prev := *r_12
		*r_12 = prev + 2
		log = append(log, "elseif-iife")
		return prev
	}() == start+1 {
		if r_12 == 999 {
			log = append(log, "elseif-body-hygienic")
		}
	}

	return log, r_12
}

// Negative soundness guards: MUST NOT hoist!
func NegativeGuards(x int32) int32 {
	lhsCalls := int32(0)
	lhs := func() bool {
		lhsCalls++
		return false
	}
	// Guard A: RHS of && MUST NOT hoist before lhs()!
	if lhs() && (func() bool {
		lhsCalls += 100
		return true
	}()) {
		return -1
	}
	// Guard B: non-nil if Init statement MUST NOT hoist!
	if y := x + 1; (func() int32 {
		z := y + 1
		return z
	}()) > 10 {
		lhsCalls += 10
	}
	// Guard C: defer inside IIFE MUST NOT hoist to enclosing function!
	if (func() bool {
		defer func() {}()
		tmp := x > 0
		return tmp
	}()) {
		lhsCalls += 1
	}
	return lhsCalls
}
`)

	dfaDir := filepath.Join(tmpDir, "dfa")
	peephole := []passes.Pass{passes.NewJ2gortPeepholePass()}

	// Run 1
	stats1, err := passes.RunOnPackageDir(dfaDir, peephole, passes.Options{})
	if err != nil {
		t.Fatalf("Run 1 failed: %v", err)
	}
	if stats1.TotalSites != 5 {
		t.Fatalf("expected exact 5 hoisted if-condition IIFE sites on Run 1, got %+v", stats1)
	}

	// Run 2: 100% idempotency
	stats2, err := passes.RunOnPackageDir(dfaDir, peephole, passes.Options{})
	if err != nil {
		t.Fatalf("Run 2 failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("expected Run 2 idempotency (0 sites, 0 files modified), got %+v", stats2)
	}

	codeBytes, err := os.ReadFile(filepath.Join(dfaDir, "dfa.go"))
	if err != nil {
		t.Fatal(err)
	}
	code := string(codeBytes)

	for _, want := range []string{
		"var iifeCond_1 int32",
		"iifeCond_1 = old_13",
		"if iifeCond_1 > limit {",
		"var iifeCond_1 N",
		"if iifeCond_1 == this_.zeroNode {",
	} {
		if !strings.Contains(code, want) {
			t.Fatalf("expected lowered dfa.go to contain %q:\n%s", want, code)
		}
	}

	// Analyze() and CheckGenericReturn() must contain ZERO remaining closure literals (`func()`)!
	analyzeIdx := strings.Index(code, "func (this_ *DataFlowAnalysisBase[N, L]) Analyze")
	evalIdx := strings.Index(code, "func EvaluateComplexConditions")
	negIdx := strings.Index(code, "func NegativeGuards")
	if analyzeIdx < 0 || evalIdx <= analyzeIdx || negIdx <= evalIdx {
		t.Fatalf("unexpected function layout in dfa.go:\n%s", code)
	}
	if strings.Contains(code[analyzeIdx:evalIdx], "func() ") {
		t.Fatalf("expected zero remaining 0-arg FuncLit IIFEs in generic methods, got:\n%s", code[analyzeIdx:evalIdx])
	}
	// EvaluateComplexConditions has only `rhsProbe := func() bool`, zero IIFEs!
	if got := strings.Count(code[evalIdx:negIdx], "func() "); got != 1 {
		t.Fatalf("expected only rhsProbe func() in EvaluateComplexConditions, got %d:\n%s", got, code[evalIdx:negIdx])
	}

	// Negative guards must still retain their 3 un-hoisted FuncLits (+ lhs + defer func() == 5)!
	if got := strings.Count(code[negIdx:], "func() "); got != 5 {
		t.Fatalf("expected negative guards to retain all 3 IIFEs (+ lhs + defer func()), got %d `func() `:\n%s", got, code[negIdx:])
	}

	// Compile & execute runtime semantic tests on lowered code!
	mustWrite("dfa/dfa_test.go", `package dfa

import (
	"reflect"
	"testing"
)

func TestLoweredIfIIFERuntimeSemantics(t *testing.T) {
	d := &DataFlowAnalysisBase[string, int]{
		states:   []*stepState{{stepCount: 0}},
		zeroNode: "",
	}
	tripped, steps := d.Analyze(5)
	if !tripped || steps != 7 {
		t.Fatalf("Analyze(5) = (%v, %d), want (true, 7)", tripped, steps)
	}
	if !d.CheckGenericReturn("") || d.CheckGenericReturn("nonempty") {
		t.Fatalf("CheckGenericReturn mismatch")
	}

	// When start=0, first IIFE increments counter to 1 (<=10 -> false), so rhsProbe MUST NOT execute (short-circuit!),
	// then else-if IIFE sees prev==1 == start+1 and outer r_12 MUST remain 999!
	log1, r1 := EvaluateComplexConditions(0, true)
	wantLog1 := []string{"iife1", "elseif-iife", "elseif-body-hygienic"}
	if !reflect.DeepEqual(log1, wantLog1) || r1 != 999 {
		t.Fatalf("EvaluateComplexConditions(0, true) = (%v, %d), want (%v, 999)", log1, r1, wantLog1)
	}

	// When start=10 and allowRhs=true, first IIFE increments counter to 11 (>10 -> true), rhsProbe runs AFTER iife1,
	// and nested unary !IIFE runs!
	log2, r2 := EvaluateComplexConditions(10, true)
	wantLog2 := []string{"iife1", "rhs", "branch1-then", "nested-not-then"}
	if !reflect.DeepEqual(log2, wantLog2) || r2 != 999 {
		t.Fatalf("EvaluateComplexConditions(10, true) = (%v, %d), want (%v, 999)", log2, r2, wantLog2)
	}

	// Verify negative guards preserved short-circuit & init semantics
	if got := NegativeGuards(10); got != 12 {
		t.Fatalf("NegativeGuards(10) = %d, want 12 (1 from lhs + 10 from Guard B + 1 from Guard C)", got)
	}
}
`)

	cmd := exec.Command("go", "test", "-v", "./dfa")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test ./dfa failed: %v\n%s", err, string(out))
	}
}
