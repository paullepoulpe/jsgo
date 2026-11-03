package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

type jsonStreamOutEntry struct {
	Src       string `json:"src"`
	Path      string `json:"path"`
	SourceMap string `json:"source_map"`
	WebpackID string `json:"webpack_id"`
}

func TestJsonStreamsOut_MultiChunkParallelParity(t *testing.T) {
	srcDir := t.TempDir()
	externPath := filepath.Join(srcDir, "externs.js")
	if err := os.WriteFile(externPath, []byte(`/** @externs */
/** @type {!Object} */ var window;
/** @param {...*} args */ function extSink(args) {}
`), 0o644); err != nil {
		t.Fatal(err)
	}

	const numFilledChunks = 6
	var chunkFiles [numFilledChunks][2]string
	for c := range numFilledChunks {
		for f := range 2 {
			path := filepath.Join(srcDir, fmt.Sprintf("chunk%d_file%d.js", c, f))
			chunkFiles[c][f] = path
			var body string
			switch {
			case c == 0 && f == 0:
				// Inlined helper whose script license must NOT bleed into chunk 1
				// under ScriptNodeLicensesOnlyTracker.
				body = `/** @license Base0 Exclusive License */
function sharedInlineHelper(v) {
  return v * 5 + 3;
}
window['root0'] = sharedInlineHelper(2);
//# sourceMappingURL=chunk0_file0.js.map
`
				inMap := `{"version":3,"file":"chunk0_file0.js","sources":["orig_base.ts"],"sourcesContent":["export const sharedInlineHelper = (v: number) => v * 5 + 3;\n"],"names":["sharedInlineHelper","v"],"mappings":";AAAA,SAASA,kBAAkBC,CAAC;AAC1B,SAAOA,CAAC,GAAG,CAAC,GAAG,CAAC;AAClB"}`
				if err := os.WriteFile(path+".map", []byte(inMap), 0o644); err != nil {
					t.Fatal(err)
				}
			case (c == 0 && f == 1) || (c == 3 && f == 0):
				// Identical script license across two different chunks (c0 and c3):
				// per-chunk ScriptNodeLicensesOnlyTracker must emit it in BOTH chunks.
				body = fmt.Sprintf(`/** @license Shared CrossChunk License */
function shared_lic_fn_%[1]d_%[2]d(x) {
  return sharedInlineHelper(x + %[1]d) + %[2]d;
}
extSink(shared_lic_fn_%[1]d_%[2]d(%[1]d));
`, c, f)
			case c == 2:
				// Identical script license on both files in chunk 2:
				// must be deduplicated to a single emission within chunk 2.
				body = fmt.Sprintf(`/** @license IntraChunk2 Dup License */
function c2_fn_%d(x) {
  return sharedInlineHelper(x + %d) + %d;
}
extSink(c2_fn_%d(%d));
`, f, c+f, f, f, c*10+f)
			default:
				body = fmt.Sprintf(`/** @license Chunk %[1]d File %[2]d License */
function fn_%[1]d_%[2]d(x) {
  return sharedInlineHelper(x + %[1]d) + %[2]d;
}
extSink(fn_%[1]d_%[2]d(%[1]d));
`, c, f)
			}
			if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}

	// Weak source file to ensure weak inputs + the synthetic $weak$ chunk are skipped identically.
	weakPath := filepath.Join(srcDir, "weak_unused.js")
	if err := os.WriteFile(weakPath, []byte(`/** @license Weak Unused License */
var unusedWeakSymbol = 999;
`), 0o644); err != nil {
		t.Fatal(err)
	}

	for _, includeContent := range []bool{false, true} {
		for _, trailingNL := range []bool{true, false} {
			t.Run(fmt.Sprintf("includeContent=%v/trailingNL=%v", includeContent, trailingNL), func(t *testing.T) {
				runOnce := func(threads int) (stdout, stderr, finalCompilerSM string) {
					args := []string{
						"--env", "CUSTOM",
						"--compilation_level", "ADVANCED_OPTIMIZATIONS",
						"--externs", externPath,
						"--json_streams", "OUT",
						"--chunk_output_path_prefix", "dist/bundles/",
						"--apply_input_source_maps",
						"--formatting", "PRINT_INPUT_DELIMITER",
						"--num_parallel_threads", fmt.Sprintf("%d", threads),
					}
					if includeContent {
						args = append(args, "--source_map_include_content")
					}
					for c := range numFilledChunks {
						args = append(args, "--js", chunkFiles[c][0], "--js", chunkFiles[c][1])
						dep := ""
						if c > 0 {
							dep = ":c0"
						}
						fileCount := 2
						if c == numFilledChunks-1 {
							// Include the trailing WEAKDEP source in the last chunk spec so
							// JSChunkGraph populates and retains the synthetic $weak$ chunk.
							fileCount = 3
						}
						args = append(args, "--chunk", fmt.Sprintf("c%d:%d%s", c, fileCount, dep))
						switch c {
						case 0:
							args = append(args, "--chunk_wrapper", "c0:/* hdr %basename% */%n%(function(){%s})();")
						case 1:
							args = append(args, "--chunk_wrapper", "c1:(function(){%output%})(); // %basename%")
						case 2:
							// Wrapper with empty prefix (%s at position 0).
							args = append(args, "--chunk_wrapper", "c2:%s%n%//# sourceMappingURL=%basename%.map")
						case 3:
							// No wrapper on c3.
						default:
							args = append(args, "--chunk_wrapper", fmt.Sprintf("c%d:/* c%d */(%s)", c, c, "%s"))
						}
					}
					// Include an empty active chunk depending on c0, plus a weakdep input.
					args = append(args,
						"--chunk", "c_empty:0:c0",
						"--chunk_wrapper", "c_empty:/* empty %basename% */%s",
					)

					var outBuf, errBuf bytes.Buffer
					runner := NewCommandLineRunnerWithOutAndErr(
						args, jutil.NewPrintStream(&outBuf), jutil.NewPrintStream(&errBuf))
					cfg := runner.GetCommandLineConfig()
					cfg.includeTrailingNewline = trailingNL
					cfg.chunkOutputFiles = []string{"c1:custom/c1.bundle.js"}
					// Also append a WEAKDEP entry to mixedJsSources to exercise weak-chunk handling.
					cfg.mixedJsSources = append(
						append([]*AbstractCommandLineRunner_FlagEntry[AbstractCommandLineRunner_JsSourceType]{}, cfg.mixedJsSources...),
						NewAbstractCommandLineRunner_FlagEntry(AbstractCommandLineRunner_JsSourceType_WEAKDEP, weakPath),
					)

					exitCode := -1
					if runner.ShouldRunCompiler() {
						runner.SetExitCodeReceiver(func(code int) { exitCode = code })
						runner.Run()
					}
					if exitCode != 0 {
						t.Fatalf("threads=%d failed exit=%d stderr=%s", threads, exitCode, errBuf.String())
					}
					if runner.GetCompiler().GetChunkGraph().GetChunkByName(jscomp.JSChunk_WEAK_CHUNK_NAME) == nil {
						t.Fatalf("expected synthetic %s chunk in chunk graph", jscomp.JSChunk_WEAK_CHUNK_NAME)
					}

					var smBuf strings.Builder
					lastName := "final_check.js"
					runner.GetCompiler().GetSourceMap().AppendTo(&smBuf, lastName)
					return outBuf.String(), errBuf.String(), smBuf.String()
				}

				out1, err1, sm1 := runOnce(1)
				out4, err4, sm4 := runOnce(4)

				if err1 != err4 {
					t.Fatalf("stderr mismatch:\nthreads=1: %q\nthreads=4: %q", err1, err4)
				}
				if out1 != out4 {
					t.Fatalf("JSON stream stdout mismatch:\n--- threads=1 ---\n%s\n--- threads=4 ---\n%s", out1, out4)
				}
				if sm1 != sm4 {
					t.Fatalf("final compiler.GetSourceMap() state mismatch:\n--- threads=1 ---\n%s\n--- threads=4 ---\n%s", sm1, sm4)
				}

				var entries []jsonStreamOutEntry
				if err := json.Unmarshal([]byte(out4), &entries); err != nil {
					t.Fatalf("unmarshal JSON stream output: %v", err)
				}
				// 6 filled chunks + 1 empty non-weak chunk = 7 (the $weak$ chunk must be omitted).
				if want := numFilledChunks + 1; len(entries) != want {
					t.Fatalf("len(entries) = %d, want %d", len(entries), want)
				}
				if entries[1].Path != "dist/bundles/custom/c1.bundle.js" {
					t.Fatalf("entries[1].Path = %q, want dist/bundles/custom/c1.bundle.js", entries[1].Path)
				}
				if !strings.Contains(entries[1].Src, "// c1.bundle.js") {
					t.Fatalf("entries[1].Src missing %%basename%% expansion: %q", entries[1].Src)
				}
				// Shared CrossChunk License must appear in BOTH c0 and c3 (per-chunk ScriptNodeLicensesOnlyTracker).
				if !strings.Contains(entries[0].Src, "Shared CrossChunk License") {
					t.Fatalf("c0 missing Shared CrossChunk License: %q", entries[0].Src)
				}
				if !strings.Contains(entries[3].Src, "Shared CrossChunk License") {
					t.Fatalf("c3 missing Shared CrossChunk License: %q", entries[3].Src)
				}
				// IntraChunk2 Dup License must appear exactly once in c2.
				if got := strings.Count(entries[2].Src, "IntraChunk2 Dup License"); got != 1 {
					t.Fatalf("c2 IntraChunk2 Dup License count = %d, want 1; src=%q", got, entries[2].Src)
				}
				// Base0 Exclusive License must NOT bleed into c1 despite sharedInlineHelper inlining.
				if strings.Contains(entries[1].Src, "Base0 Exclusive License") {
					t.Fatalf("c1 unexpectedly contains Base0 Exclusive License: %q", entries[1].Src)
				}
				// Every entry must have a valid JSON V3 source_map matching its Path.
				for i, e := range entries {
					var parsedSM map[string]any
					if err := json.Unmarshal([]byte(e.SourceMap), &parsedSM); err != nil {
						t.Fatalf("entry[%d] (%s) invalid source_map JSON: %v", i, e.Path, err)
					}
					if parsedSM["file"] != e.Path {
						t.Fatalf("entry[%d] source_map file = %v, want %q", i, parsedSM["file"], e.Path)
					}
				}
			})
		}
	}
}

func TestJsonStreamsBoth_MultiChunkParallelParity(t *testing.T) {
	type inSpec struct {
		Path      string `json:"path"`
		Src       string `json:"src"`
		SourceMap string `json:"source_map,omitempty"`
	}
	inSpecs := []inSpec{
		{
			Path: "pkg/base.js",
			Src: `/** @license Shared Stream License */
function baseCalc(a) { return a * 7 + 1; }
console.log(baseCalc(1));
`,
			SourceMap: `{"version":3,"file":"pkg/base.js","sources":["pkg/base.ts"],"sourcesContent":["export function baseCalc(a: number) { return a * 7 + 1; }\nconsole.log(baseCalc(1));\n"],"names":["baseCalc","a"],"mappings":";AAAA,SAASA,QAAQC,CAAC;AAAE,SAAOA,CAAC,GAAG,CAAC,GAAG,CAAC;AAAC;AACrC,QAAQ,CAAC,GAAG,CAACD,QAAQ,CAAC,CAAC,CAAC,CAAC"}`,
		},
		{
			Path: "pkg/mod_a1.js",
			Src:  "/** @license Shared Stream License */\nconsole.log(baseCalc(2));\n",
		},
		{
			Path: "pkg/mod_a2.js",
			Src:  "/** @license ModA License */\nconsole.log(baseCalc(3));\n",
		},
		{
			Path: "pkg/mod_b.js",
			Src:  "/** @license ModB License */\nconsole.log(baseCalc(4));\n",
		},
		{
			Path: "pkg/mod_c.js",
			Src:  "console.log(baseCalc(5));\n",
		},
	}
	stdinBytes, err := json.Marshal(inSpecs)
	if err != nil {
		t.Fatal(err)
	}

	runBoth := func(threads int) (string, string) {
		args := []string{
			"--json_streams=BOTH",
			"--compilation_level=SIMPLE_OPTIMIZATIONS",
			"--apply_input_source_maps",
			"--source_map_include_content",
			"--chunk=sub/dir/base:1",
			"--chunk=feature--alpha.v1:2:sub/dir/base",
			"--chunk=feature--beta.v1:1:sub/dir/base",
			"--chunk=feature--gamma.v1:1:feature--alpha.v1",
			"--chunk_wrapper=sub/dir/base:/* base:%basename% */%n%(function(){%s})();",
			"--chunk_wrapper=feature--alpha.v1:(function(){%output%})();",
			fmt.Sprintf("--num_parallel_threads=%d", threads),
		}
		var outBuf, errBuf bytes.Buffer
		runner := NewCommandLineRunnerWithStreams(
			args,
			bytes.NewReader(stdinBytes),
			jutil.NewPrintStream(&outBuf),
			jutil.NewPrintStream(&errBuf),
		)
		exitCode := -1
		if runner.ShouldRunCompiler() {
			runner.SetExitCodeReceiver(func(code int) { exitCode = code })
			runner.Run()
		}
		if exitCode != 0 {
			t.Fatalf("threads=%d failed exit=%d stderr=%s", threads, exitCode, errBuf.String())
		}
		return outBuf.String(), errBuf.String()
	}

	out1, err1 := runBoth(1)
	for iter := range 3 {
		out4, err4 := runBoth(4)
		if err1 != err4 {
			t.Fatalf("iter=%d stderr mismatch:\nthreads=1: %q\nthreads=4: %q", iter, err1, err4)
		}
		if out1 != out4 {
			t.Fatalf("iter=%d JSON stream BOTH output mismatch:\n--- threads=1 ---\n%s\n--- threads=4 ---\n%s", iter, out1, out4)
		}
	}
}

type customHookCommandLineRunner struct {
	*CommandLineRunner
	banner string
}

func (c *customHookCommandLineRunner) WriteChunkOutput(
	fileName string,
	out jutil.Appendable,
	lt jscomp.CodePrinter_LicenseTracker,
	m *jscomp.JSChunk,
) {
	out.AppendCharSequence(c.banner)
	c.CommandLineRunner.WriteChunkOutput(fileName, out, lt, m)
}

func TestJsonStreamsBoth_CustomRunnerWriteChunkOutputOverride(t *testing.T) {
	stdinJSON := `[
		{"path":"a.js","src":"var a = 1; console.log(a);"},
		{"path":"b.js","src":"var b = a + 2; console.log(b);"}
	]`
	args := []string{
		"--json_streams=BOTH",
		"--chunk=c0:1",
		"--chunk=c1:1:c0",
		"--num_parallel_threads=4",
	}
	var outBuf, errBuf bytes.Buffer
	inStream := io.Reader(strings.NewReader(stdinJSON))
	outPS := jutil.NewPrintStream(&outBuf)
	errPS := jutil.NewPrintStream(&errBuf)
	base := NewCommandLineRunnerWithStreams(args, inStream, outPS, errPS)
	custom := &customHookCommandLineRunner{
		CommandLineRunner: base,
		banner:            "/* CUSTOM_HOOK */\n",
	}
	custom.self = custom

	exitCode := custom.DoRun()
	if exitCode != 0 {
		t.Fatalf("custom runner DoRun exit=%d stderr=%s", exitCode, errBuf.String())
	}
	var entries []jsonStreamOutEntry
	if err := json.Unmarshal(outBuf.Bytes(), &entries); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("len(entries)=%d, want 2", len(entries))
	}
	for i, e := range entries {
		if !strings.HasPrefix(e.Src, "/* CUSTOM_HOOK */\n") {
			t.Fatalf("entry[%d] missing custom WriteChunkOutput banner: %q", i, e.Src)
		}
	}
}

func TestJsonStreamsBoth_SingleActiveChunkPlusWeakParity(t *testing.T) {
	// With --chunk=solo:2, chunkGraph.GetAllChunks() has len == 2 ([solo, $weak$]),
	// exercising len(chunks) > 1 with a single active chunk and parallel threads.
	stdinJSON := `[
		{"path":"s0.js","src":"/** @license Solo Lic */\nfunction solo(x) { return x + 10; }\n"},
		{"path":"s1.js","src":"/** @license Solo Lic */\nconsole.log(solo(5));\n"}
	]`
	run := func(threads int) string {
		args := []string{
			"--json_streams=BOTH",
			"--chunk=solo:2",
			"--chunk_wrapper=solo:/* %basename% */\n(%s)",
			fmt.Sprintf("--num_parallel_threads=%d", threads),
		}
		var outBuf, errBuf bytes.Buffer
		runner := NewCommandLineRunnerWithStreams(
			args,
			strings.NewReader(stdinJSON),
			jutil.NewPrintStream(&outBuf),
			jutil.NewPrintStream(&errBuf),
		)
		exitCode := -1
		if runner.ShouldRunCompiler() {
			runner.SetExitCodeReceiver(func(code int) { exitCode = code })
			runner.Run()
		}
		if exitCode != 0 {
			t.Fatalf("threads=%d exit=%d stderr=%s", threads, exitCode, errBuf.String())
		}
		return outBuf.String()
	}

	out1 := run(1)
	out4 := run(4)
	if out1 != out4 {
		t.Fatalf("single-active-chunk JSON stream mismatch:\nthreads=1: %s\nthreads=4: %s", out1, out4)
	}
}
