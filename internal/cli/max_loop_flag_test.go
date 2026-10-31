package cli

import "testing"

// Go-only: the hidden --max_optimization_loop_iterations flag sets
// CompilerOptions.setMaxOptimizationLoopIterations, and leaves it unset when absent.
func TestMaxOptimizationLoopIterationsFlag(t *testing.T) {
	for _, tc := range []struct {
		args []string
		want int
	}{
		{nil, 0},
		{[]string{"--max_optimization_loop_iterations=3"}, 3},
		{[]string{"--max_optimization_loop_iterations", "2", "--compilation_level=ADVANCED"}, 2},
	} {
		runner := NewCommandLineRunner(append([]string{}, tc.args...))
		if !runner.ShouldRunCompiler() {
			t.Fatalf("%v: flags rejected", tc.args)
		}
		if got := runner.CreateOptions().GetMaxOptimizationLoopIterations(); got != tc.want {
			t.Errorf("%v: GetMaxOptimizationLoopIterations() = %d, want %d", tc.args, got, tc.want)
		}
	}
}
