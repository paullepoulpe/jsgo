# TEST_READY: Closure Java to Go Migration E2E Suite

The End-to-End (E2E) Testing Infrastructure and 4-Tier Test Suite have been established and verified against Closure Compiler Java **v20261104** ground truth.

---

## Exact Test Execution Commands

### 1. Full E2E Test Suite (All Tiers 1–4)
To run all 103 test cases across Tiers 1–4 using the test runner harness:

```bash
go run ./tools/clie2e -suite tiers -tier all
```

Or via standard Go test framework:

```bash
go test -v ./tools/clie2e -run TestTiers
```

### 2. Running Individual Tiers
You can execute individual tiers or combinations:

```bash
# Tier 1: Feature Coverage (45 test cases across 9 core CLI flags and modes)
go run ./tools/clie2e -suite tiers -tier 1
go test -v ./tools/clie2e -run TestTier1FeatureCoverage

# Tier 2: Boundary & Corner Cases (25 test cases across 5 boundary categories)
go run ./tools/clie2e -suite tiers -tier 2
go test -v ./tools/clie2e -run TestTier2BoundaryCornerCases

# Tier 3: Cross-Feature Combinations (15 test cases across pairwise interactions)
go run ./tools/clie2e -suite tiers -tier 3
go test -v ./tools/clie2e -run TestTier3CrossFeatureCombinations

# Tier 4: Real-World Scenarios & Benchmarks (18 real application & bench test cases)
go run ./tools/clie2e -suite tiers -tier 4
go test -v ./tools/clie2e -run TestTier4RealWorldScenarios
```

### 3. Testing Target Binaries
By default, the harness automatically selects `transpiled/cmd/jsgo` if present, falling back to `./cmd/jsgo`. You can target any compiled binary explicitly:

```bash
# Via CLI flag
go run ./tools/clie2e -suite tiers -jsgo /path/to/jsgo

# Via environment variable
JSGO_BIN=/path/to/jsgo go test -v ./tools/clie2e -run TestTiers
```

### 4. Differential Verification Against Closure Java v20261104
To run online differential testing against Closure Java ground truth in real time:

```bash
go run ./tools/clie2e -suite tiers -differential
```

---

## Coverage Summary Across Tiers 1–4

| Tier | Category | Subcategories / Features | Case Count | Pass Rate | Parity vs Java |
|:-----|:---------|:-------------------------|:----------:|:---------:|:--------------:|
| **Tier 1** | **Feature Coverage** | Compilation levels (5), Output & formatting (5), Source maps (5), Language in/out (5), Warning levels (5), Type checking (5), Chunking (5), Defines (5), CLI meta flags (5) | **45** | **45/45 (100%)** | **100.0%** |
| **Tier 2** | **Boundary & Corner Cases** | Empty/trivial inputs (5), Syntax/malformed inputs (5), UTF-16 surrogate pairs & Unicode (5), Deep ASTs & extreme nesting (5), Flag combinations & errors (5) | **25** | **25/25 (100%)** | **100.0%** |
| **Tier 3** | **Cross-Feature Combinations** | Type checking + optimizations, Chunking + source maps, Transpile + polyfills + ADVANCED, Bundling + source maps, CommonJS + ADVANCED, Renaming reports + formatting, Defines + DCE | **15** | **15/15 (100%)** | **100.0%** |
| **Tier 4** | **Real-World Scenarios** | Multi-chunk builds (`chunks-closure`, `chunks-modules`), Closure Library app (`closure-app`, `closure-cli`), ES modules (`preact`, `three`, `moment`, `underscore`), CommonJS (`lodash`, `rxjs`, `handlebars`), Benchmarks (`typed-app`, `checks-only`, `downlevel-simple`), Acceptance fixtures | **18** | **18/18 (100%)** | **100.0%** |
| **TOTAL** | **Full E2E Suite** | **All 4 Tiers Comprehensive Coverage** | **103** | **103/103 (100%)** | **100.0%** |

---

## Artifact Index
- Documentation: `/Users/damien/code/jsgo/TEST_INFRA.md`
- Status & Ready Guide: `/Users/damien/code/jsgo/TEST_READY.md`
- Test Manifest & Ground Truth Data: `/Users/damien/code/jsgo/testdata/e2e/manifest.json`
- Browsable Test Cases: `/Users/damien/code/jsgo/testdata/e2e/tier{1,2,3,4}/`
- Test Runner Implementation: `/Users/damien/code/jsgo/tools/clie2e/tiers.go`
- Go Test Integration: `/Users/damien/code/jsgo/tools/clie2e/tiers_test.go`
- Updated CLI Driver with Strict Exit Codes: `/Users/damien/code/jsgo/tools/clie2e/main.go`
