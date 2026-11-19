package main

import (
	"os"
	"testing"
	"time"
)

// TestTiers runs the comprehensive 4-Tier E2E test suite via standard `go test`.
func TestTiers(t *testing.T) {
	manifest, err := LoadTierManifest()
	if err != nil {
		t.Fatalf("Failed to load tier manifest: %v", err)
	}

	tmpDir, err := os.MkdirTemp("", "jsgo-test-tiers-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	bin := buildJsgo(tmpDir)

	for _, tc := range manifest.Cases {
		t.Run(tc.ID, func(t *testing.T) {
			res := runCase(bin, tc, 30*time.Second, false)
			if !res.Passed {
				for _, d := range res.Diffs {
					t.Errorf("Diff: %s", d)
				}
				t.Logf("Got exit: %d, Want exit: %d", res.GotExit, tc.Expected.ExitCode)
				t.Logf("Got stdout:\n%s", res.GotOut)
				t.Logf("Want stdout:\n%s", tc.Expected.Stdout)
				t.Logf("Got stderr:\n%s", res.GotErr)
				t.Logf("Want stderr:\n%s", tc.Expected.Stderr)
			}
		})
	}
}

// TestTier1FeatureCoverage runs only Tier 1 feature tests.
func TestTier1FeatureCoverage(t *testing.T) {
	manifest, err := LoadTierManifest()
	if err != nil {
		t.Fatalf("Failed to load tier manifest: %v", err)
	}
	cases := filterCases(manifest.Cases, "1", "")
	if len(cases) == 0 {
		t.Fatal("No Tier 1 test cases found")
	}

	tmpDir, err := os.MkdirTemp("", "jsgo-test-tier1-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	bin := buildJsgo(tmpDir)
	for _, tc := range cases {
		t.Run(tc.ID, func(t *testing.T) {
			res := runCase(bin, tc, 30*time.Second, false)
			if !res.Passed {
				t.Errorf("Test %s failed: %v", tc.ID, res.Diffs)
			}
		})
	}
}

// TestTier2BoundaryCornerCases runs only Tier 2 boundary tests.
func TestTier2BoundaryCornerCases(t *testing.T) {
	manifest, err := LoadTierManifest()
	if err != nil {
		t.Fatalf("Failed to load tier manifest: %v", err)
	}
	cases := filterCases(manifest.Cases, "2", "")
	if len(cases) == 0 {
		t.Fatal("No Tier 2 test cases found")
	}

	tmpDir, err := os.MkdirTemp("", "jsgo-test-tier2-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	bin := buildJsgo(tmpDir)
	for _, tc := range cases {
		t.Run(tc.ID, func(t *testing.T) {
			res := runCase(bin, tc, 30*time.Second, false)
			if !res.Passed {
				t.Errorf("Test %s failed: %v", tc.ID, res.Diffs)
			}
		})
	}
}

// TestTier3CrossFeatureCombinations runs only Tier 3 combination tests.
func TestTier3CrossFeatureCombinations(t *testing.T) {
	manifest, err := LoadTierManifest()
	if err != nil {
		t.Fatalf("Failed to load tier manifest: %v", err)
	}
	cases := filterCases(manifest.Cases, "3", "")
	if len(cases) == 0 {
		t.Fatal("No Tier 3 test cases found")
	}

	tmpDir, err := os.MkdirTemp("", "jsgo-test-tier3-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	bin := buildJsgo(tmpDir)
	for _, tc := range cases {
		t.Run(tc.ID, func(t *testing.T) {
			res := runCase(bin, tc, 30*time.Second, false)
			if !res.Passed {
				t.Errorf("Test %s failed: %v", tc.ID, res.Diffs)
			}
		})
	}
}

// TestTier4RealWorldScenarios runs only Tier 4 real-world tests.
func TestTier4RealWorldScenarios(t *testing.T) {
	manifest, err := LoadTierManifest()
	if err != nil {
		t.Fatalf("Failed to load tier manifest: %v", err)
	}
	cases := filterCases(manifest.Cases, "4", "")
	if len(cases) == 0 {
		t.Fatal("No Tier 4 test cases found")
	}

	tmpDir, err := os.MkdirTemp("", "jsgo-test-tier4-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	bin := buildJsgo(tmpDir)
	for _, tc := range cases {
		t.Run(tc.ID, func(t *testing.T) {
			res := runCase(bin, tc, 30*time.Second, false)
			if !res.Passed {
				t.Errorf("Test %s failed: %v", tc.ID, res.Diffs)
			}
		})
	}
}
