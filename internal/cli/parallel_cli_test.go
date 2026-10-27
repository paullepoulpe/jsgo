package cli

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestParallelPasses_MultiChunkSourceMapsAndWrappersParity(t *testing.T) {
	srcDir := t.TempDir()
	externPath := filepath.Join(srcDir, "externs.js")
	if err := os.WriteFile(externPath, []byte(`/** @externs */
/** @type {!Object} */ var window;
/** @param {...*} args */ function extSink(args) {}
`), 0o666); err != nil {
		t.Fatal(err)
	}

	const numChunks = 8
	var chunkFiles [numChunks][2]string
	for c := 0; c < numChunks; c++ {
		for f := 0; f < 2; f++ {
			name := filepath.Join(srcDir, fmt.Sprintf("c%d_f%d.js", c, f))
			chunkFiles[c][f] = name
			var body string
			if c == 0 && f == 0 {
				body = fmt.Sprintf(`/** @license Shared License Header */
function sharedUtil(v) {
  return v * 3 + 1;
}
window['root0'] = sharedUtil(1);
//# sourceMappingURL=c0_f0.js.map
`)
				// Provide an input source map for c0_f0.js so sharedUtil inlined across
				// multiple chunks exercises concurrent SourceMapConsumerV3 line decoding!
				inMap := `{"version":3,"file":"c0_f0.js","sources":["orig_shared.ts"],"sourcesContent":["export const sharedUtil = (v: number) => v * 3 + 1;\n"],"names":["sharedUtil","v"],"mappings":";AAAA,SAASA,UAAUC,CAAC;AAClB,SAAOA,CAAC,GAAG,CAAC,GAAG,CAAC;AAClB"}`
				if err := os.WriteFile(name+".map", []byte(inMap), 0o666); err != nil {
					t.Fatal(err)
				}
			} else {
				body = fmt.Sprintf(`/** @license Chunk %[1]d License */
function helper_%[1]d_%[2]d(x) {
  const y = sharedUtil(x + %[1]d);
  return y + %[2]d;
}
extSink(helper_%[1]d_%[2]d(%[1]d));
`, c, f)
			}
			if err := os.WriteFile(name, []byte(body), 0o666); err != nil {
				t.Fatal(err)
			}
		}
	}

	for _, includeContent := range []bool{true, false} {
		runCLI := func(threads int) (outDir string, stdout string, stderr string) {
			outDir = t.TempDir()
			prefix := filepath.Join(outDir, "chunk_")
			args := []string{
				"--env", "CUSTOM",
				"--compilation_level", "ADVANCED_OPTIMIZATIONS",
				"--externs", externPath,
				"--chunk_output_path_prefix", prefix,
				"--create_source_map", "%outname%.map",
				"--source_map_location_mapping", outDir + "|out",
				"--apply_input_source_maps",
				"--formatting", "PRINT_INPUT_DELIMITER",
				"--num_parallel_threads", fmt.Sprintf("%d", threads),
			}
			if includeContent {
				args = append(args, "--source_map_include_content")
			}
			for c := 0; c < numChunks; c++ {
				args = append(args, "--js", chunkFiles[c][0], "--js", chunkFiles[c][1])
				dep := ""
				if c > 0 {
					dep = ":c0"
				}
				args = append(args,
					"--chunk", fmt.Sprintf("c%d:2%s", c, dep),
					"--chunk_wrapper", fmt.Sprintf("c%d:/* wrapper %d %%basename%% */\n(function(){%%s})();\n//# sourceMappingURL=chunk_c%d.js.map", c, c, c),
				)
			}
			var outBuf, errBuf bytes.Buffer
			runner := NewCommandLineRunnerWithOutAndErr(
				args, jutil.NewPrintStream(&outBuf), jutil.NewPrintStream(&errBuf))
			exitCode := -1
			if runner.ShouldRunCompiler() {
				runner.SetExitCodeReceiver(func(code int) { exitCode = code })
				runner.Run()
			}
			if exitCode != 0 {
				t.Fatalf("includeContent=%v threads=%d failed with exit=%d stderr=%s", includeContent, threads, exitCode, errBuf.String())
			}
			return outDir, outBuf.String(), errBuf.String()
		}

		dir1, out1, err1 := runCLI(1)
		dir4, out4, err4 := runCLI(4)
		if out1 != out4 || err1 != err4 {
			t.Fatalf("includeContent=%v CLI stream mismatch:\nseq=(%q, %q)\npar=(%q, %q)", includeContent, out1, err1, out4, err4)
		}
		for c := 0; c < numChunks; c++ {
			for _, ext := range []string{".js", ".js.map"} {
				name := fmt.Sprintf("chunk_c%d%s", c, ext)
				b1, err := os.ReadFile(filepath.Join(dir1, name))
				if err != nil {
					t.Fatalf("missing seq %s: %v", name, err)
				}
				b4, err := os.ReadFile(filepath.Join(dir4, name))
				if err != nil {
					t.Fatalf("missing par %s: %v", name, err)
				}
				if !bytes.Equal(b1, b4) {
					t.Fatalf("includeContent=%v mismatch in %s:\n--- threads=1 ---\n%s\n--- threads=4 ---\n%s", includeContent, name, string(b1), string(b4))
				}
			}
		}
	}
}

func TestParallelPasses_SaveAndRestoreStateWithInlineSourceMapsParity(t *testing.T) {
	srcDir := t.TempDir()
	extPath := filepath.Join(srcDir, "externs.js")
	if err := os.WriteFile(extPath, []byte("/** @externs */\nvar window;\n/** @param {*} x */\nfunction extOut(x) {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	const numChunks = 4
	const filesPerChunk = 6
	var commonArgs []string
	commonArgs = append(commonArgs,
		"--compilation_level", "ADVANCED_OPTIMIZATIONS",
		"--warning_level", "VERBOSE",
		"--language_in", "ECMASCRIPT_2021",
		"--language_out", "ECMASCRIPT_2019",
		"--env", "CUSTOM",
		"--externs", extPath,
		"--parse_inline_source_maps",
		"--apply_input_source_maps",
		"--source_map_include_content",
	)

	for c := 0; c < numChunks; c++ {
		for f := 0; f < filesPerChunk; f++ {
			jsPath := filepath.Join(srcDir, fmt.Sprintf("m%d_f%d.js", c, f))
			origTS := filepath.Join(srcDir, fmt.Sprintf("orig_m%d_f%d.ts", c, f))
			inlineJSON := fmt.Sprintf(
				`{"version":3,"file":%q,"sources":[%q],"sourcesContent":[%q],"names":["val_%d_%d"],"mappings":"AAAA;AACA,SAASA,GAAG,CAAC,IAAI;AACjB,SAAO,IAAI,GAAG,CAAC;AACjB"}`,
				jsPath,
				origTS,
				fmt.Sprintf("// Café résumé λ %d:%d\nexport function fn_%d_%d(v: number) { return v + %d; }\n", c, f, c, f, c*10+f),
				c, f,
			)
			b64 := jutil.BaseEncoding_Base64().Encode([]byte(inlineJSON))
			var code string
			if c == 0 && f == 0 {
				code = fmt.Sprintf(
					"/** @license MultiStage InlineMap License */\n/** @constructor */\nfunction Box(v) { this.value = v; }\nBox.prototype.read = function() { return `box:${this.value}`; };\n/** @nosideeffects */\nfunction sharedCalc(n) { return new Box(n * 2 + 1).read(); }\nextOut(sharedCalc(1));\n//# sourceMappingURL=data:application/json;charset=utf-8;base64,%s\n",
					b64,
				)
			} else {
				code = fmt.Sprintf(
					"/** @desc Message %d %d */\nvar MSG_%d_%d = goog_msg_%d_%d;\nfunction goog_msg_%d_%d(x) { return sharedCalc(x + %d); }\nextOut(MSG_%d_%d(%d));\n//# sourceMappingURL=data:application/json;base64,%s\n",
					c, f, c, f, c, f, c, f, c+f, c, f, f, b64,
				)
			}
			if err := os.WriteFile(jsPath, []byte(code), 0o644); err != nil {
				t.Fatal(err)
			}
			commonArgs = append(commonArgs, "--js", jsPath)
		}
		if c == 0 {
			commonArgs = append(commonArgs, "--chunk", fmt.Sprintf("m0:%d", filesPerChunk))
		} else {
			commonArgs = append(commonArgs, "--chunk", fmt.Sprintf("m%d:%d:m0", c, filesPerChunk))
		}
	}

	gunzipFile := func(path string) []byte {
		f, err := os.Open(path)
		if err != nil {
			t.Fatalf("open %s: %v", path, err)
		}
		defer f.Close()
		zr, err := gzip.NewReader(f)
		if err != nil {
			t.Fatalf("gzip reader %s: %v", path, err)
		}
		defer zr.Close()
		raw, err := io.ReadAll(zr)
		if err != nil {
			t.Fatalf("gunzip %s: %v", path, err)
		}
		return raw
	}

	runThreeStage := func(threads int) (checksRaw, optRaw []byte, outDir string) {
		workDir := t.TempDir()
		checksState := filepath.Join(workDir, "checks.state")
		optState := filepath.Join(workDir, "opt.state")
		outDir = filepath.Join(workDir, "out")
		if err := os.MkdirAll(outDir, 0o755); err != nil {
			t.Fatal(err)
		}
		prefix := filepath.Join(outDir, "chunk_")

		execStage := func(stageName string, extra ...string) {
			args := append([]string{}, commonArgs...)
			args = append(args,
				"--num_parallel_threads", fmt.Sprintf("%d", threads),
				"--segment_of_compilation_to_run", stageName,
				"--chunk_output_path_prefix", prefix,
				"--create_source_map", "%outname%.map",
				"--source_map_location_mapping", outDir+"|out",
			)
			args = append(args, extra...)
			var outBuf, errBuf bytes.Buffer
			exitCode := -1
			runner := NewCommandLineRunnerWithOutAndErr(
				args, jutil.NewPrintStream(&outBuf), jutil.NewPrintStream(&errBuf))
			if runner.ShouldRunCompiler() {
				runner.SetExitCodeReceiver(func(code int) { exitCode = code })
				runner.Run()
			}
			if exitCode != 0 {
				t.Fatalf("threads=%d stage=%s failed exit=%d stderr=%s", threads, stageName, exitCode, errBuf.String())
			}
		}

		execStage("CHECKS", "--filename_to_save_to", checksState)
		execStage("OPTIMIZATIONS", "--filename_to_restore_from", checksState, "--filename_to_save_to", optState)
		execStage("FINALIZATIONS", "--filename_to_restore_from", optState)

		return gunzipFile(checksState), gunzipFile(optState), outDir
	}

	checks1, opt1, dir1 := runThreeStage(1)
	checks4, opt4, dir4 := runThreeStage(4)

	if !bytes.Equal(checks1, checks4) {
		t.Fatalf("CHECKS SaveState decompressed bytes differ: len(1)=%d len(4)=%d", len(checks1), len(checks4))
	}
	if !bytes.Equal(opt1, opt4) {
		t.Fatalf("OPTIMIZATIONS SaveState decompressed bytes differ: len(1)=%d len(4)=%d", len(opt1), len(opt4))
	}
	for c := 0; c < numChunks; c++ {
		for _, ext := range []string{".js", ".js.map"} {
			name := fmt.Sprintf("chunk_m%d%s", c, ext)
			b1, err := os.ReadFile(filepath.Join(dir1, name))
			if err != nil {
				t.Fatalf("missing seq %s: %v", name, err)
			}
			b4, err := os.ReadFile(filepath.Join(dir4, name))
			if err != nil {
				t.Fatalf("missing par %s: %v", name, err)
			}
			if !bytes.Equal(b1, b4) {
				t.Fatalf("mismatch in %s:\n--- threads=1 ---\n%s\n--- threads=4 ---\n%s", name, string(b1), string(b4))
			}
		}
	}
}
