package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/paullepoulpe/jsgo/tools/oracle/pin"
)

var (
	tierFlag     = flag.String("tier", "", "tiers suite: filter by tier (1, 2, 3, 4, 1,2, or all)")
	diffFlag     = flag.Bool("differential", false, "tiers suite: run differentially against Closure Java v20261104")
	caseFlag     = flag.String("case", "", "tiers suite: filter by test case ID or substring")
	tierJobs     = flag.Int("tier-jobs", 4, "tiers suite: concurrent test execution workers")
	tierTimeout  = flag.Duration("tier-timeout", 30*time.Second, "tiers suite: per-test execution timeout")
)

// TierTestCase represents a single test case across Tiers 1-4.
type TierTestCase struct {
	ID          string            `json:"id"`
	Tier        int               `json:"tier"`
	Feature     string            `json:"feature"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Args        []string          `json:"args"`
	Stdin       string            `json:"stdin,omitempty"`
	Files       map[string]string `json:"files,omitempty"`
	Expected    TierExpected      `json:"expected"`
}

// TierExpected holds the expected outputs derived from Closure Java ground truth.
type TierExpected struct {
	ExitCode    int               `json:"exit_code"`
	Stdout      string            `json:"stdout,omitempty"`
	Stderr      string            `json:"stderr,omitempty"`
	OutputFiles map[string]string `json:"output_files,omitempty"`
}

// TierManifest is the format of testdata/e2e/manifest.json.
type TierManifest struct {
	Version string         `json:"version"`
	Cases   []TierTestCase `json:"cases"`
}

// LoadTierManifest finds and reads testdata/e2e/manifest.json.
func LoadTierManifest() (*TierManifest, error) {
	manifestPath := findManifestPath()
	b, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("reading manifest at %s: %w", manifestPath, err)
	}
	var m TierManifest
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("parsing manifest %s: %w", manifestPath, err)
	}
	return &m, nil
}

func findManifestPath() string {
	// Try relative to repo root
	dir, err := os.Getwd()
	if err == nil {
		for {
			candidate := filepath.Join(dir, "testdata", "e2e", "manifest.json")
			if _, err := os.Stat(candidate); err == nil {
				return candidate
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	return "testdata/e2e/manifest.json"
}

// TierTestResult records the execution outcome of a test case.
type TierTestResult struct {
	Case     TierTestCase
	Passed   bool
	Duration time.Duration
	Diffs    []string
	GotExit  int
	GotOut   string
	GotErr   string
}

// RunTiers executes the Tier 1-4 tests and reports the results.
func runTiers() int {
	manifest, err := LoadTierManifest()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 2
	}

	workDir, err := os.MkdirTemp("", "jsgo-tiers-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating temp dir: %v\n", err)
		return 2
	}
	defer os.RemoveAll(workDir)

	bin := buildJsgo(workDir)

	targetTier := "all"
	if *tierFlag != "" {
		targetTier = *tierFlag
	}

	selected := filterCases(manifest.Cases, targetTier, *caseFlag)
	if len(selected) == 0 {
		fmt.Fprintln(os.Stderr, "No test cases matched criteria")
		return 2
	}

	fmt.Printf("================================================================================\n")
	fmt.Printf("Closure Java -> Go E2E Test Suite (Tiers 1–4)\n")
	fmt.Printf("Target binary: %s\n", bin)
	fmt.Printf("Oracle spec:   Closure Compiler %s ground truth\n", manifest.Version)
	fmt.Printf("Selected:      %d test cases (Tier filter: %s, Case filter: %q)\n", len(selected), targetTier, *caseFlag)
	fmt.Printf("Workers:       %d\n", *tierJobs)
	fmt.Printf("================================================================================\n\n")

	results := executeCases(bin, selected, *tierJobs, *tierTimeout, *diffFlag)

	// Summarize results by Tier
	tierCounts := map[int]struct{ pass, total int }{}
	featureCounts := map[string]struct{ pass, total int }{}
	var failures []TierTestResult

	for _, res := range results {
		t := res.Case.Tier
		f := res.Case.Feature
		tc := tierCounts[t]
		tc.total++
		fc := featureCounts[f]
		fc.total++
		if res.Passed {
			tc.pass++
			fc.pass++
		} else {
			failures = append(failures, res)
		}
		tierCounts[t] = tc
		featureCounts[f] = fc
	}

	// Print summary table
	fmt.Println("Tier Results Summary:")
	fmt.Println("--------------------------------------------------------------------------------")
	tierNames := map[int]string{
		1: "Tier 1: Feature Coverage (Flags & Options)",
		2: "Tier 2: Boundary & Corner Cases",
		3: "Tier 3: Cross-Feature Combinations",
		4: "Tier 4: Real-World Scenarios & Benchmarks",
	}

	totalPass, totalTests := 0, 0
	for t := 1; t <= 4; t++ {
		tc, exists := tierCounts[t]
		if !exists {
			continue
		}
		totalPass += tc.pass
		totalTests += tc.total
		status := "PASS"
		if tc.pass < tc.total {
			status = "FAIL"
		}
		fmt.Printf("  %-42s [%4s] %3d/%3d passed\n", tierNames[t], status, tc.pass, tc.total)
	}
	fmt.Println("--------------------------------------------------------------------------------")
	totalStatus := "PASS"
	if totalPass < totalTests {
		totalStatus = "FAIL"
	}
	fmt.Printf("  %-42s [%4s] %3d/%3d passed (%0.1f%%)\n\n", "TOTAL E2E SUITE", totalStatus, totalPass, totalTests, float64(totalPass)*100/float64(totalTests))

	if len(failures) > 0 {
		fmt.Printf("Failures (%d):\n", len(failures))
		fmt.Println("================================================================================")
		for _, fail := range failures {
			fmt.Printf("FAIL: %s (Tier %d - %s: %s)\n", fail.Case.ID, fail.Case.Tier, fail.Case.Feature, fail.Case.Name)
			fmt.Printf("  Args: %s\n", strings.Join(fail.Case.Args, " "))
			for _, d := range fail.Diffs {
				fmt.Printf("  - %s\n", d)
			}
			if *keep {
				fmt.Printf("  Expected exit: %d, Got exit: %d\n", fail.Case.Expected.ExitCode, fail.GotExit)
				if fail.GotOut != fail.Case.Expected.Stdout {
					fmt.Printf("  Got stdout:\n%s\n  Want stdout:\n%s\n", indent(fail.GotOut, "    "), indent(fail.Case.Expected.Stdout, "    "))
				}
				if fail.GotErr != fail.Case.Expected.Stderr {
					fmt.Printf("  Got stderr:\n%s\n  Want stderr:\n%s\n", indent(fail.GotErr, "    "), indent(fail.Case.Expected.Stderr, "    "))
				}
			}
			fmt.Println()
		}
		return 1
	}

	fmt.Println("ALL E2E TESTS PASSED SUCCESSFULLY!")
	return 0
}

func filterCases(cases []TierTestCase, tierFilter, caseFilter string) []TierTestCase {
	tierSet := map[int]bool{}
	if tierFilter == "all" || tierFilter == "" {
		for i := 1; i <= 4; i++ {
			tierSet[i] = true
		}
	} else {
		for _, part := range strings.Split(tierFilter, ",") {
			part = strings.TrimSpace(part)
			if n, err := strconv.Atoi(part); err == nil {
				tierSet[n] = true
			}
		}
	}

	var out []TierTestCase
	for _, tc := range cases {
		if !tierSet[tc.Tier] {
			continue
		}
		if caseFilter != "" {
			if !strings.Contains(tc.ID, caseFilter) && !strings.Contains(tc.Feature, caseFilter) && !strings.Contains(tc.Name, caseFilter) {
				continue
			}
		}
		out = append(out, tc)
	}
	return out
}

func executeCases(bin string, cases []TierTestCase, jobs int, timeout time.Duration, differential bool) []TierTestResult {
	results := make([]TierTestResult, len(cases))
	sem := make(chan bool, jobs)
	var wg sync.WaitGroup

	for i, tc := range cases {
		wg.Add(1)
		sem <- true
		go func(idx int, c TierTestCase) {
			defer func() {
				<-sem
				wg.Done()
			}()
			results[idx] = runCase(bin, c, timeout, differential)
		}(i, tc)
	}
	wg.Wait()
	return results
}

func runCase(bin string, tc TierTestCase, timeout time.Duration, differential bool) TierTestResult {
	start := time.Now()
	res := TierTestResult{Case: tc}

	workDir, err := os.MkdirTemp("", "tc-"+tc.ID+"-*")
	if err != nil {
		res.Diffs = append(res.Diffs, fmt.Sprintf("failed to create temp workdir: %v", err))
		return res
	}
	defer os.RemoveAll(workDir)

	// Write input files
	for fn, content := range tc.Files {
		fp := filepath.Join(workDir, fn)
		os.MkdirAll(filepath.Dir(fp), 0o777)
		if err := os.WriteFile(fp, []byte(content), 0o666); err != nil {
			res.Diffs = append(res.Diffs, fmt.Sprintf("failed to write %s: %v", fn, err))
			return res
		}
	}

	// Prepare execution
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, bin, tc.Args...)
	cmd.Dir = workDir
	cmd.Env = env()

	if tc.Stdin != "" {
		cmd.Stdin = strings.NewReader(tc.Stdin)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr

	runErr := cmd.Run()
	res.Duration = time.Since(start)

	exitCode := 0
	if ee, ok := runErr.(*exec.ExitError); ok {
		exitCode = ee.ExitCode()
	} else if runErr != nil {
		exitCode = -1
	}

	// In unix exit code wraps at 256
	if exitCode < 0 || exitCode > 255 {
		exitCode = exitCode & 0xFF
	}
	res.GotExit = exitCode
	res.GotOut = stdout.String()
	res.GotErr = normalizeStderr(stderr.String(), workDir)

	// Determine expectations (either from Differential Java run or from recorded tc.Expected)
	expected := tc.Expected
	if differential {
		javaExp, err := runJavaOracle(tc, workDir)
		if err == nil {
			expected = javaExp
		}
	}

	// Compare Exit Code
	// If expected exit code is non-zero, allow any non-zero exit code (e.g. 1 or 255 or 2)
	// when checking invalid flag / compilation failure, but enforce exact match when exit is 0.
	if expected.ExitCode == 0 {
		if exitCode != 0 {
			res.Diffs = append(res.Diffs, fmt.Sprintf("exit code: got %d, want 0", exitCode))
		}
	} else {
		if exitCode == 0 {
			res.Diffs = append(res.Diffs, fmt.Sprintf("exit code: got 0, want non-zero (%d)", expected.ExitCode))
		}
	}

	// Compare Stdout
	normGotOut := normalizeOutput(res.GotOut, workDir)
	normWantOut := normalizeOutput(expected.Stdout, workDir)
	if normGotOut != normWantOut {
		// If both are multi-line code, check without trailing whitespace diffs
		if strings.TrimSpace(normGotOut) != strings.TrimSpace(normWantOut) {
			res.Diffs = append(res.Diffs, fmt.Sprintf("stdout mismatch: got %q, want %q", truncate(normGotOut, 80), truncate(normWantOut, 80)))
		}
	}

	// Compare Stderr
	// For diagnostics and error messages, check that essential error tokens match
	if expected.Stderr != "" {
		normWantErr := normalizeStderr(expected.Stderr, workDir)
		if !matchStderr(res.GotErr, normWantErr) {
			res.Diffs = append(res.Diffs, fmt.Sprintf("stderr mismatch: got %q, want %q", truncate(res.GotErr, 80), truncate(normWantErr, 80)))
		}
	}

	// Compare Output Files
	for relPath, wantContent := range expected.OutputFiles {
		gotPath := filepath.Join(workDir, relPath)
		gotBytes, err := os.ReadFile(gotPath)
		if err != nil {
			res.Diffs = append(res.Diffs, fmt.Sprintf("missing expected output file: %s", relPath))
			continue
		}
		if !sameOutput(gotBytes, []byte(wantContent)) {
			res.Diffs = append(res.Diffs, fmt.Sprintf("output file content mismatch: %s", relPath))
		}
	}

	res.Passed = (len(res.Diffs) == 0)
	return res
}

func normalizeStderr(s, dir string) string {
	s = strings.ReplaceAll(s, dir+"/", "")
	lines := strings.Split(s, "\n")
	var clean []string
	for _, l := range lines {
		if strings.HasPrefix(l, "WARNING: ") ||
			strings.HasPrefix(l, "\tat ") ||
			strings.HasPrefix(l, "\t... ") {
			continue
		}
		clean = append(clean, l)
	}
	return strings.TrimSpace(strings.Join(clean, "\n"))
}

func normalizeOutput(s, dir string) string {
	s = strings.ReplaceAll(s, dir+"/", "")
	return s
}

func matchStderr(got, want string) bool {
	if got == want {
		return true
	}
	gotClean := strings.TrimSpace(got)
	wantClean := strings.TrimSpace(want)
	if gotClean == wantClean {
		return true
	}
	// For diagnostic lines (e.g. JSC_TYPE_MISMATCH or "is not a valid option")
	if strings.Contains(wantClean, "ERROR - [") {
		start := strings.Index(wantClean, "ERROR - [")
		end := strings.Index(wantClean[start:], "]")
		if end != -1 {
			token := wantClean[start : start+end+1]
			if strings.Contains(gotClean, token) {
				return true
			}
		}
	}
	if strings.Contains(wantClean, "is not a valid option") && strings.Contains(gotClean, "is not a valid option") {
		return true
	}
	return false
}

func runJavaOracle(tc TierTestCase, workDir string) (TierExpected, error) {
	cache := defaultCacheDir()
	jar := filepath.Join(cache, pin.MustLoad().Jar)
	clie2eBld := filepath.Join(cache, "clie2e-build")
	javaBin := "/opt/homebrew/opt/openjdk/bin/java"
	if _, err := os.Stat(javaBin); err != nil {
		javaBin = "java"
	}
	cp := jar
	if _, err := os.Stat(clie2eBld); err == nil {
		cp = clie2eBld + ":" + jar
	}

	cmd := exec.Command(javaBin, append([]string{"-cp", cp, "com.google.javascript.jscomp.CommandLineRunner"}, tc.Args...)...)
	cmd.Dir = workDir
	cmd.Env = env()
	if tc.Stdin != "" {
		cmd.Stdin = strings.NewReader(tc.Stdin)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr

	err := cmd.Run()
	exitCode := 0
	if ee, ok := err.(*exec.ExitError); ok {
		exitCode = ee.ExitCode()
	}
	if exitCode < 0 || exitCode > 255 {
		exitCode = exitCode & 0xFF
	}

	outFiles := map[string]string{}
	for fn := range tc.Expected.OutputFiles {
		fp := filepath.Join(workDir, fn)
		if b, err := os.ReadFile(fp); err == nil {
			outFiles[fn] = string(b)
		}
	}

	return TierExpected{
		ExitCode:    exitCode,
		Stdout:      stdout.String(),
		Stderr:      stderr.String(),
		OutputFiles: outFiles,
	}, nil
}

func truncate(s string, max int) string {
	s = strings.ReplaceAll(s, "\n", "\\n")
	if len(s) > max {
		return s[:max] + "..."
	}
	return s
}

func indent(s, prefix string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = prefix + l
	}
	return strings.Join(lines, "\n")
}
