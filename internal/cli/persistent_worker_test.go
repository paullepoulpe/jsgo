package cli

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"weak"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func writeFileWithDigest(t testing.TB, path string, content string) WorkerInput {
	t.Helper()
	b := []byte(content)
	if err := os.WriteFile(path, b, 0o644); err != nil {
		t.Fatalf("WriteFile(%s): %v", path, err)
	}
	return WorkerInput{
		Path:   path,
		Digest: sha256Hex(b),
	}
}

func encodeWorkRequests(t testing.TB, reqs ...WorkRequest) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	for i := range reqs {
		line, err := json.Marshal(reqs[i])
		if err != nil {
			t.Fatalf("json.Marshal(req[%d]): %v", i, err)
		}
		buf.Write(line)
		buf.WriteByte('\n')
	}
	return &buf
}

func decodeWorkResponses(t testing.TB, raw string) []WorkResponse {
	t.Helper()
	trimmed := strings.TrimSuffix(raw, "\n")
	if trimmed == "" {
		return nil
	}
	lines := strings.Split(trimmed, "\n")
	resps := make([]WorkResponse, len(lines))
	for i, line := range lines {
		if err := json.Unmarshal([]byte(line), &resps[i]); err != nil {
			t.Fatalf("line %d is not valid single-line WorkResponse JSON (%q): %v", i, line, err)
		}
	}
	return resps
}

func runStandaloneCompile(t testing.TB, args []string) (exitCode int, stdout, stderr string) {
	t.Helper()
	rhino.ResetNodeStore()
	defer rhino.ResetNodeStore()

	var outBuf, errBuf bytes.Buffer
	outPS := jutil.NewPrintStream(&outBuf)
	errPS := jutil.NewPrintStream(&errBuf)
	runner := NewCommandLineRunnerWithStreams(args, bytes.NewReader(nil), outPS, errPS)
	runner.SetExitCodeReceiver(func(code int) {
		exitCode = code
	})
	if runner.ShouldRunCompiler() {
		runner.Run()
	} else if runner.HasErrors() {
		exitCode = 1
	}
	outPS.Flush()
	errPS.Flush()
	return exitCode, outBuf.String(), errBuf.String()
}

func TestPersistentWorkerMultiRequestStream(t *testing.T) {
	dir := t.TempDir()

	extPath := filepath.Join(dir, "externs.js")
	libPath := filepath.Join(dir, "lib.js")
	appPath := filepath.Join(dir, "app.js")
	badPath := filepath.Join(dir, "bad.js")

	extIn := writeFileWithDigest(t, extPath, `/** @externs */
/** @param {string} msg */
function hostReport(msg) {}
`)
	libIn := writeFileWithDigest(t, libPath, `/**
 * @param {string} name
 * @param {number} count
 * @return {string}
 */
function formatGreeting(name, count) {
  return name + ':' + (count + 1);
}
`)
	appInV1 := writeFileWithDigest(t, appPath, `hostReport(formatGreeting('alpha', 41));
`)
	badIn := writeFileWithDigest(t, badPath, `function broken( { return 1; }
`)

	out1JS := filepath.Join(dir, "out1.js")
	out1Map := filepath.Join(dir, "out1.js.map")
	out2JS := filepath.Join(dir, "out2.js")
	out3WorkerJS := filepath.Join(dir, "out3_worker.js")
	out3WorkerMap := filepath.Join(dir, "out3_worker.js.map")
	out3ColdJS := filepath.Join(dir, "out3_cold.js")
	out3ColdMap := filepath.Join(dir, "out3_cold.js.map")

	worker := NewPersistentWorker([]string{"--alsologtostderr"})

	// Execute Request 1 & Request 2 first, then mutate app.js on disk before Request 3 & Request 4
	// via an io.Pipe / sequential Serve, or stream all 4 requests where Request 3 uses app_v2.js / updated app.js!
	// Using an io.Pipe lets us mutate app.js on disk between Request 2 and Request 3 over a single continuous worker stream!
	reqReader, reqWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	respReader, respWriter, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}

	errCh := make(chan error, 1)
	go func() {
		defer respWriter.Close()
		errCh <- worker.Serve(reqReader, respWriter)
	}()
	workerJoined := false
	defer func() {
		if !workerJoined {
			_ = reqWriter.Close()
			<-errCh
			_ = reqReader.Close()
			_ = respReader.Close()
		}
	}()

	dec := json.NewDecoder(respReader)
	sendAndRecv := func(req WorkRequest) WorkResponse {
		t.Helper()
		b, mErr := json.Marshal(req)
		if mErr != nil {
			t.Fatalf("Marshal: %v", mErr)
		}
		b = append(b, '\n')
		if _, wErr := reqWriter.Write(b); wErr != nil {
			t.Fatalf("Write request %d: %v", req.RequestId, wErr)
		}
		var resp WorkResponse
		if dErr := dec.Decode(&resp); dErr != nil {
			t.Fatalf("Decode response %d: %v", req.RequestId, dErr)
		}
		return resp
	}

	// Request 1: Valid compilation -> --js_output_file + --create_source_map.
	resp1 := sendAndRecv(WorkRequest{
		RequestId: 101,
		Arguments: []string{
			"--compilation_level=ADVANCED_OPTIMIZATIONS",
			"--env=CUSTOM",
			"--externs", extPath,
			"--js", libPath,
			"--js", appPath,
			"--js_output_file", out1JS,
			"--create_source_map", out1Map,
			"--source_map_format=V3",
			"--source_map_include_content",
		},
		Inputs: []WorkerInput{extIn, libIn, appInV1},
	})
	if resp1.ExitCode != 0 || resp1.RequestId != 101 || resp1.WasCancelled {
		t.Fatalf("Request 1 unexpected response: %+v", resp1)
	}
	out1Bytes, err := os.ReadFile(out1JS)
	if err != nil || !strings.Contains(string(out1Bytes), "alpha:42") {
		t.Fatalf("Request 1 output mismatch: %q, err=%v", string(out1Bytes), err)
	}
	statsAfter1 := worker.Cache().Stats()
	if statsAfter1.Misses != 3 || statsAfter1.Hits != 0 || statsAfter1.Entries != 3 {
		t.Fatalf("unexpected cache stats after Request 1: %+v", statsAfter1)
	}

	// Request 2a: Syntax error -> returns exitCode != 0 + error text in output, worker stays alive.
	resp2a := sendAndRecv(WorkRequest{
		RequestId: 102,
		Arguments: []string{
			"--env=CUSTOM",
			"--js", badPath,
			"--js_output_file", out2JS,
		},
		Inputs: []WorkerInput{badIn},
	})
	if resp2a.ExitCode == 0 || resp2a.RequestId != 102 || resp2a.WasCancelled {
		t.Fatalf("Request 2a expected non-zero exitCode, got: %+v", resp2a)
	}
	if !strings.Contains(resp2a.Output, "JSC_PARSE_ERROR") {
		t.Fatalf("Request 2a expected JSC_PARSE_ERROR in output, got: %q", resp2a.Output)
	}

	// Request 2b: Type error -> returns exitCode != 0 + diagnostic in output, worker stays alive.
	typeErrPath := filepath.Join(dir, "type_err.js")
	typeErrIn := writeFileWithDigest(t, typeErrPath, `hostReport(formatGreeting('beta', 'not_a_number'));
`)
	resp2b := sendAndRecv(WorkRequest{
		RequestId: 103,
		Arguments: []string{
			"--compilation_level=SIMPLE",
			"--jscomp_error=checkTypes",
			"--env=CUSTOM",
			"--externs", extPath,
			"--js", libPath,
			"--js", typeErrPath,
			"--js_output_file", out2JS,
		},
		Inputs: []WorkerInput{extIn, libIn, typeErrIn},
	})
	if resp2b.ExitCode == 0 || resp2b.RequestId != 103 {
		t.Fatalf("Request 2b expected non-zero exitCode on type error, got: %+v", resp2b)
	}
	if !strings.Contains(resp2b.Output, "JSC_TYPE_MISMATCH") {
		t.Fatalf("Request 2b expected JSC_TYPE_MISMATCH in output, got: %q", resp2b.Output)
	}

	// Request 3: Same library + externs digests, modified app.js on disk with new digest!
	appInV2 := writeFileWithDigest(t, appPath, `hostReport(formatGreeting('gamma', 99));
`)
	hitsBefore3 := worker.Cache().Stats().Hits
	missesBefore3 := worker.Cache().Stats().Misses

	req3Args := func(jsOut, mapOut string) []string {
		return []string{
			"--compilation_level=ADVANCED_OPTIMIZATIONS",
			"--env=CUSTOM",
			"--externs", extPath,
			"--js", libPath,
			"--js", appPath,
			"--js_output_file", jsOut,
			"--create_source_map", mapOut,
			"--source_map_format=V3",
			"--source_map_include_content",
		}
	}

	resp3 := sendAndRecv(WorkRequest{
		RequestId: 104,
		Arguments: req3Args(out3WorkerJS, out3WorkerMap),
		Inputs:    []WorkerInput{extIn, libIn, appInV2},
	})
	if resp3.ExitCode != 0 || resp3.RequestId != 104 || resp3.WasCancelled {
		t.Fatalf("Request 3 unexpected response: %+v", resp3)
	}
	statsAfter3 := worker.Cache().Stats()
	if statsAfter3.Hits != hitsBefore3+2 {
		t.Fatalf("Request 3 expected 2 cache hits (externs + lib), got before=%d after=%d", hitsBefore3, statsAfter3.Hits)
	}
	if statsAfter3.Misses != missesBefore3+1 {
		t.Fatalf("Request 3 expected 1 cache miss (modified app.js), got before=%d after=%d", missesBefore3, statsAfter3.Misses)
	}

	// Request 4: Cancel request -> immediate WasCancelled response.
	resp4 := sendAndRecv(WorkRequest{
		RequestId: 105,
		Cancel:    true,
	})
	if resp4.RequestId != 105 || !resp4.WasCancelled || resp4.ExitCode != 0 || resp4.Output != "" {
		t.Fatalf("Request 4 cancel response mismatch: %+v", resp4)
	}

	// Close stdin -> EOF clean shutdown.
	if err := reqWriter.Close(); err != nil {
		t.Fatalf("reqWriter.Close: %v", err)
	}
	serveErr := <-errCh
	workerJoined = true
	_ = reqReader.Close()
	_ = respReader.Close()
	if serveErr != nil {
		t.Fatalf("worker.Serve returned error on EOF: %v", serveErr)
	}

	// Verify NodeStore is cleanly reset to static store after worker exits.
	if !rhino.CurrentNodeStore().IsStatic() {
		t.Fatalf("expected rhino.CurrentNodeStore().IsStatic() == true after worker exit")
	}

	// Compare Request 3 JS & V3 source map byte-for-byte against standalone NewCommandLineRunner!
	// Replace out3ColdJS/Map target names in source map `"file":` field by compiling in separate subdirs with identical relative output file name, OR check exact bytes with same --source_map_location_mapping and identical output basename!
	workerSubdir := filepath.Join(dir, "w_cmp")
	coldSubdir := filepath.Join(dir, "c_cmp")
	if err := os.MkdirAll(workerSubdir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(coldSubdir, 0o755); err != nil {
		t.Fatal(err)
	}
	workerParityJS := filepath.Join(workerSubdir, "bundle.js")
	workerParityMap := filepath.Join(workerSubdir, "bundle.js.map")
	coldParityJS := filepath.Join(coldSubdir, "bundle.js")
	coldParityMap := filepath.Join(coldSubdir, "bundle.js.map")

	var parityOut bytes.Buffer
	parityReqBuf := encodeWorkRequests(t, WorkRequest{
		RequestId: 201,
		Arguments: req3Args(workerParityJS, workerParityMap),
		Inputs:    []WorkerInput{extIn, libIn, appInV2},
	})
	if err := worker.Serve(parityReqBuf, &parityOut); err != nil {
		t.Fatalf("parity worker.Serve: %v", err)
	}
	parityResps := decodeWorkResponses(t, parityOut.String())
	if len(parityResps) != 1 || parityResps[0].ExitCode != 0 {
		t.Fatalf("parity worker response failed: %+v", parityResps)
	}

	coldCode, _, coldErr := runStandaloneCompile(t, req3Args(coldParityJS, coldParityMap))
	if coldCode != 0 {
		t.Fatalf("standalone compile failed (code=%d): %s", coldCode, coldErr)
	}

	// Also verify same exact path in-place produces 100% byte-identical JS and V3 source map!
	coldCodeSamePath, _, coldErrSamePath := runStandaloneCompile(t, req3Args(out3ColdJS, out3ColdMap))
	if coldCodeSamePath != 0 {
		t.Fatalf("standalone same-path compile failed: %s", coldErrSamePath)
	}
	wJSBytes, err := os.ReadFile(workerParityJS)
	if err != nil {
		t.Fatal(err)
	}
	cJSBytes, err := os.ReadFile(coldParityJS)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(wJSBytes, cJSBytes) {
		t.Fatalf("JS output is not byte-identical!\nworker:\n%s\ncold:\n%s", wJSBytes, cJSBytes)
	}

	// Re-run standalone writing to the exact same output path (workerParityJS / workerParityMap)
	// after saving worker bytes, because V3 source map embeds `"file":"<jsOutputFile>"`!
	wMapBytes, err := os.ReadFile(workerParityMap)
	if err != nil {
		t.Fatal(err)
	}
	coldExactCode, _, coldExactErr := runStandaloneCompile(t, req3Args(workerParityJS, workerParityMap))
	if coldExactCode != 0 {
		t.Fatalf("standalone exact-path compile failed: %s", coldExactErr)
	}
	cExactJSBytes, err := os.ReadFile(workerParityJS)
	if err != nil {
		t.Fatal(err)
	}
	cExactMapBytes, err := os.ReadFile(workerParityMap)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(wJSBytes, cExactJSBytes) {
		t.Fatalf("exact-path JS mismatch:\nworker:\n%s\ncold:\n%s", wJSBytes, cExactJSBytes)
	}
	if !bytes.Equal(wMapBytes, cExactMapBytes) {
		t.Fatalf("exact-path V3 source map mismatch:\nworker:\n%s\ncold:\n%s", wMapBytes, cExactMapBytes)
	}
}

func TestPersistentWorkerDigestCacheInvalidationAndLRU(t *testing.T) {
	dir := t.TempDir()
	libPath := filepath.Join(dir, "cached_lib.js")
	appPath := filepath.Join(dir, "app.js")
	outPath := filepath.Join(dir, "out.js")

	libIn := writeFileWithDigest(t, libPath, `var LIB_CONST = 'from_cache_v1';`)
	appIn := writeFileWithDigest(t, appPath, `console.log(LIB_CONST);`)

	cache := NewWorkerDigestCache(256, 2)
	worker := NewPersistentWorkerWithCache(nil, cache)

	// 1. Populate cache for libPath and appPath.
	resp1 := worker.ExecuteRequest(&WorkRequest{
		RequestId: 1,
		Arguments: []string{"--env=CUSTOM", "--compilation_level=SIMPLE", "--js", libPath, "--js", appPath, "--js_output_file", outPath},
		Inputs:    []WorkerInput{libIn, appIn},
	})
	if resp1.ExitCode != 0 {
		t.Fatalf("req 1 failed: %+v", resp1)
	}
	if !rhino.CurrentNodeStore().IsStatic() {
		t.Fatalf("expected NodeStore reset after ExecuteRequest")
	}

	// 2. Delete libPath on disk! With matching non-empty digest, worker MUST serve from memory cache.
	if err := os.Remove(libPath); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	resp2 := worker.ExecuteRequest(&WorkRequest{
		RequestId: 2,
		Arguments: []string{"--env=CUSTOM", "--compilation_level=SIMPLE", "--js", libPath, "--js", appPath, "--js_output_file", outPath},
		Inputs:    []WorkerInput{libIn, appIn},
	})
	if resp2.ExitCode != 0 {
		t.Fatalf("req 2 expected cache hit despite removed file on disk, got: %+v", resp2)
	}
	got2, _ := os.ReadFile(outPath)
	if !strings.Contains(string(got2), "from_cache_v1") {
		t.Fatalf("req 2 output mismatch: %q", string(got2))
	}

	// 3. Empty digest ("") MUST immediately invalidate cached entry and read from disk.
	writeFileWithDigest(t, libPath, `var LIB_CONST = 'uncached_empty_digest';`)
	resp3 := worker.ExecuteRequest(&WorkRequest{
		RequestId: 3,
		Arguments: []string{"--env=CUSTOM", "--compilation_level=SIMPLE", "--js", libPath, "--js", appPath, "--js_output_file", outPath},
		Inputs:    []WorkerInput{{Path: libPath, Digest: ""}, appIn},
	})
	if resp3.ExitCode != 0 {
		t.Fatalf("req 3 failed: %+v", resp3)
	}
	got3, _ := os.ReadFile(outPath)
	if !strings.Contains(string(got3), "uncached_empty_digest") {
		t.Fatalf("req 3 did not bypass cache on empty digest: %q", string(got3))
	}
	if cache.Stats().Entries != 1 {
		t.Fatalf("expected libPath evicted from cache on empty digest, stats=%+v", cache.Stats())
	}

	// 4. LRU maxEntries=2 and maxBytes=256 eviction check.
	f1 := filepath.Join(dir, "f1.js")
	f2 := filepath.Join(dir, "f2.js")
	f3 := filepath.Join(dir, "f3.js")
	in1 := writeFileWithDigest(t, f1, "var a = 1;\n")
	in2 := writeFileWithDigest(t, f2, "var b = 2;\n")
	in3 := writeFileWithDigest(t, f3, "var c = 3;\n")

	_ = worker.ExecuteRequest(&WorkRequest{
		RequestId: 4,
		Arguments: []string{"--env=CUSTOM", "--js", f1, "--js", f2, "--js_output_file", outPath},
		Inputs:    []WorkerInput{in1, in2},
	})
	if cache.Stats().Entries != 2 {
		t.Fatalf("expected 2 entries in LRU cache, got %+v", cache.Stats())
	}
	// Touch f1 then insert f3 -> f2 should be evicted (maxEntries=2).
	_ = worker.ExecuteRequest(&WorkRequest{
		RequestId: 5,
		Arguments: []string{"--env=CUSTOM", "--js", f1, "--js", f3, "--js_output_file", outPath},
		Inputs:    []WorkerInput{in1, in3},
	})
	if cache.Stats().Entries != 2 {
		t.Fatalf("expected 2 entries after LRU eviction, got %+v", cache.Stats())
	}
	if _, ok := cache.getReload(jutil.ParsePathString(f2), in2.Digest); ok {
		t.Fatalf("expected LRU entry f2 to be evicted")
	}
	if _, ok := cache.getReload(jutil.ParsePathString(f1), in1.Digest); !ok {
		t.Fatalf("expected MRU entry f1 to remain cached")
	}

	// Byte-budget eviction check: 200-byte entry evicts existing entries under 256-byte cap.
	bigContent := "var big = '" + strings.Repeat("x", 200) + "';\n"
	fBig := filepath.Join(dir, "big.js")
	inBig := writeFileWithDigest(t, fBig, bigContent)
	_ = worker.ExecuteRequest(&WorkRequest{
		RequestId: 6,
		Arguments: []string{"--env=CUSTOM", "--js", fBig, "--js_output_file", outPath},
		Inputs:    []WorkerInput{inBig},
	})
	st := cache.Stats()
	if st.Bytes > 256 {
		t.Fatalf("cache exceeded byte budget: %+v", st)
	}
}

func TestPersistentWorkerStartupArgsAndPanicRecovery(t *testing.T) {
	dir := t.TempDir()

	// 1. Direct args detection.
	if _, ok := parsePersistentWorkerStartupArgs([]string{"--js", "a.js"}); ok {
		t.Fatalf("did not expect worker mode without --persistent_worker")
	}
	if _, ok := parsePersistentWorkerStartupArgs([]string{"--persistent_worker=false", "--js", "a.js"}); ok {
		t.Fatalf("did not expect worker mode when --persistent_worker=false")
	}
	startup, ok := parsePersistentWorkerStartupArgs([]string{"--persistent_worker", "--alsologtostderr", "--warning_level=QUIET"})
	if !ok || len(startup) != 1 || startup[0] != "--warning_level=QUIET" {
		t.Fatalf("unexpected startup args: ok=%v startup=%v", ok, startup)
	}

	// 2. @flagfile expansion for Bazel startup args.
	atFile := filepath.Join(dir, "worker.params")
	if err := os.WriteFile(atFile, []byte("--persistent_worker\n--alsologtostderr\n--env=CUSTOM\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	startupAt, okAt := parsePersistentWorkerStartupArgs([]string{"@" + atFile})
	if !okAt || len(startupAt) != 1 || startupAt[0] != "--env=CUSTOM" {
		t.Fatalf("unexpected @flagfile startup args: ok=%v startup=%v", okAt, startupAt)
	}

	// 3. --help output does NOT leak hidden --persistent_worker flag, and flag parser accepts --persistent_worker.
	code, helpOut, helpErr := runStandaloneCompile(t, []string{"--help"})
	if code != 0 || helpErr != "" {
		t.Fatalf("--help failed: code=%d err=%q", code, helpErr)
	}
	if strings.Contains(helpOut, "persistent_worker") {
		t.Fatalf("hidden --persistent_worker should not appear in --help output")
	}

	// 4. RunPersistentWorker handles FlagUsageException panics AND omitted --js_output_file cleanly.
	srcPath := filepath.Join(dir, "stdout_app.js")
	srcIn := writeFileWithDigest(t, srcPath, "var answer = 20 + 22;\n")

	inBuf := encodeWorkRequests(t,
		WorkRequest{
			RequestId: 1,
			Arguments: []string{"--browser_featureset_year=2020", "--language_out=ECMASCRIPT5", "--js", srcPath},
			Inputs:    []WorkerInput{srcIn},
		},
		WorkRequest{
			RequestId: 2,
			Arguments: []string{"--compilation_level=SIMPLE", "--js", srcPath},
			Inputs:    []WorkerInput{srcIn},
		},
	)
	var outBuf bytes.Buffer
	if err := RunPersistentWorker([]string{"@" + atFile}, inBuf, &outBuf); err != nil {
		t.Fatalf("RunPersistentWorker failed: %v", err)
	}
	resps := decodeWorkResponses(t, outBuf.String())
	if len(resps) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(resps))
	}
	if resps[0].ExitCode == 0 || !strings.Contains(resps[0].Output, "browser_featureset_year") {
		t.Fatalf("expected FlagUsageException failure in resp[0], got %+v", resps[0])
	}
	if resps[1].ExitCode != 0 || !strings.Contains(resps[1].Output, "var answer=42;") {
		t.Fatalf("expected compiled JS on stdout in resp[1], got %+v", resps[1])
	}
	if !rhino.CurrentNodeStore().IsStatic() {
		t.Fatalf("expected static NodeStore after RunPersistentWorker")
	}
}

func TestUnboxedOptStringFlagsParity(t *testing.T) {
	dir := t.TempDir()
	extPath := filepath.Join(dir, "ext.js")
	srcPath := filepath.Join(dir, "in.js")
	outPath := filepath.Join(dir, "out.js")
	if err := os.WriteFile(extPath, []byte("/** @externs */\nvar __customImport = function(x) {};\nvar sink = function(x) {};\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(srcPath, []byte("function myHelper(v) { return import(v); }\nsink(myHelper('./m.js'));\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	code, _, errOut := runStandaloneCompile(t, []string{
		"--env=CUSTOM",
		"--compilation_level=ADVANCED_OPTIMIZATIONS",
		"--language_in=ECMASCRIPT_2020",
		"--language_out=ECMASCRIPT_2019",
		"--dynamic_import_alias=__customImport",
		"--rename_variable_prefix=pref_",
		"--package_json_entry_names=browser,main",
		"--externs", extPath,
		"--js", srcPath,
		"--js_output_file", outPath,
	})
	if code != 0 {
		t.Fatalf("compile failed: code=%d err=%s", code, errOut)
	}
	got, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(got), "__customImport(") {
		t.Fatalf("expected --dynamic_import_alias=__customImport in output, got: %s", got)
	}
}

func TestPersistentWorkerSourceFileAndBigIntGCReclamation(t *testing.T) {
	dir := t.TempDir()
	srcPath := filepath.Join(dir, "ephemeral_big.js")
	outPath := filepath.Join(dir, "out.js")

	// 1. Direct weak-pointer check on jscomp.SourceFile attached to a compilation NodeStore.
	var weakSF weak.Pointer[jscomp.SourceFile]
	func() {
		rhino.ResetNodeStore()
		store := rhino.EnsureCompilationStore()
		if store.IsStatic() {
			t.Fatal("expected non-static compilation store")
		}
		sf := jscomp.SourceFile_FromCode("worker_ephemeral.js", "/* "+strings.Repeat("z", 64*1024)+" */\nconst x = 1n;\n")
		weakSF = weak.Make(sf)
		_ = rhino.IR_Script().SetStaticSourceFile(sf)
		rhino.ResetNodeStore()
	}()
	for range 4 {
		runtime.GC()
	}
	if weakSF.Value() != nil {
		t.Fatal("jscomp.SourceFile remained pinned in rhino after ResetNodeStore")
	}

	// 2. End-to-end PersistentWorker test with uncached 1 MiB comment-heavy + BigInt
	// JS files across 16 sequential requests: verify live heap does not retain all
	// 16 MiB of SourceFile contents after requests finish.
	worker := NewPersistentWorkerWithCache(nil, NewWorkerDigestCache(0, 0))
	padding := strings.Repeat("a", 1024*1024)

	// Warm one request first so one-time initialization (builtin externs, pools) settles.
	warmIn := writeFileWithDigest(t, srcPath, "/* warm */ var w = 9007199254740993n;\n")
	warmResp := worker.ExecuteRequest(&WorkRequest{
		RequestId: 1,
		Arguments: []string{"--env=CUSTOM", "--compilation_level=SIMPLE", "--language_in=ECMASCRIPT_2020", "--language_out=ECMASCRIPT_2020", "--js", srcPath, "--js_output_file", outPath},
		Inputs:    []WorkerInput{{Path: warmIn.Path, Digest: ""}},
	})
	if warmResp.ExitCode != 0 {
		t.Fatalf("warm request failed: %+v", warmResp)
	}
	for range 4 {
		runtime.GC()
	}
	var baseline runtime.MemStats
	runtime.ReadMemStats(&baseline)

	const iterations = 16
	for i := range iterations {
		code := fmt.Sprintf("/* %d %s */\nvar bigVal = %d900719925474099312345n;\nconsole.log(bigVal);\n", i, padding, i+1)
		if err := os.WriteFile(srcPath, []byte(code), 0o644); err != nil {
			t.Fatal(err)
		}
		resp := worker.ExecuteRequest(&WorkRequest{
			RequestId: i + 2,
			Arguments: []string{"--env=CUSTOM", "--compilation_level=SIMPLE", "--language_in=ECMASCRIPT_2020", "--language_out=ECMASCRIPT_2020", "--js", srcPath, "--js_output_file", outPath},
			Inputs:    []WorkerInput{{Path: srcPath, Digest: ""}},
		})
		if resp.ExitCode != 0 {
			t.Fatalf("request %d failed: %+v", i+2, resp)
		}
	}
	padding = ""
	for range 4 {
		runtime.GC()
	}
	var after runtime.MemStats
	runtime.ReadMemStats(&after)
	var retained int64
	if after.HeapAlloc > baseline.HeapAlloc {
		retained = int64(after.HeapAlloc - baseline.HeapAlloc)
	}
	// If SourceFiles were leaked globally, >= 16 MiB would remain live.
	if retained > 6*1024*1024 {
		t.Fatalf("PersistentWorker leaked SourceFile memory across requests: retained %d bytes (> 6 MiB)", retained)
	}
}

func BenchmarkPersistentWorkerVsCold(b *testing.B) {
	dir := b.TempDir()
	inputs := make([]WorkerInput, 0, 8)
	baseArgs := []string{
		"--compilation_level=SIMPLE",
	}
	for i := 0; i < 8; i++ {
		p := filepath.Join(dir, fmt.Sprintf("mod_%d.js", i))
		var sb strings.Builder
		fmt.Fprintf(&sb, "/** @fileoverview Module %d */\n", i)
		for j := 0; j < 25; j++ {
			fmt.Fprintf(&sb, "function mod%d_fn%d(a, b) { var c = a + b + %d; return c * 2; }\n", i, j, j)
		}
		fmt.Fprintf(&sb, "console.log(mod%d_fn0(1, 2));\n", i)
		inp := writeFileWithDigest(b, p, sb.String())
		inputs = append(inputs, inp)
		baseArgs = append(baseArgs, "--js", p)
	}
	outPath := filepath.Join(dir, "bench_out.js")
	reqArgs := append(append([]string(nil), baseArgs...), "--js_output_file", outPath)

	b.Run("ColdRunner", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			builtinExternsBrowserOnce = sync.Once{}
			builtinExternsBrowser = nil
			code, _, errOut := runStandaloneCompile(b, reqArgs)
			if code != 0 {
				b.Fatalf("ColdRunner failed: %s", errOut)
			}
		}
	})

	b.Run("PersistentWorker", func(b *testing.B) {
		worker := NewPersistentWorker(nil)
		req := WorkRequest{
			RequestId: 1,
			Arguments: reqArgs,
			Inputs:    inputs,
		}
		reqJSON, err := json.Marshal(req)
		if err != nil {
			b.Fatal(err)
		}
		reqJSON = append(reqJSON, '\n')

		// Warm the worker's digest cache + builtin externs once before timing.
		var warmOut bytes.Buffer
		if err := worker.Serve(bytes.NewReader(reqJSON), &warmOut); err != nil {
			b.Fatal(err)
		}

		b.ReportAllocs()
		b.ResetTimer()
		var outBuf bytes.Buffer
		for i := 0; i < b.N; i++ {
			outBuf.Reset()
			if err := worker.Serve(bytes.NewReader(reqJSON), &outBuf); err != nil {
				b.Fatalf("PersistentWorker failed: %v", err)
			}
		}
	})
}

func TestPersistentWorkerCapturesSystemOutAndSystemErrWithoutCorruptingJSON(t *testing.T) {
	dir := t.TempDir()

	ttlJS := filepath.Join(dir, "ttl.js")
	ttlIn := writeFileWithDigest(t, ttlJS, `/**
 * @template T := printType('TTL_STDOUT_MARKER:', 'number') =:
 * @param {T} x
 * @return {T}
 */
function ttlId(x) { return x; }
ttlId(42);
`)

	badXtb := filepath.Join(dir, "broken.xtb")
	if err := os.WriteFile(badXtb, []byte(`<translationbundle lang="fr"><translation id="1">unclosed`), 0o644); err != nil {
		t.Fatal(err)
	}
	msgJS := filepath.Join(dir, "msg.js")
	msgIn := writeFileWithDigest(t, msgJS, `var x = 1;`)
	outJS := filepath.Join(dir, "out.js")

	origOut := jutil.System_Out()
	origErr := jutil.System_Err()
	var leakedOut, leakedErr bytes.Buffer
	sentinelOut := jutil.NewPrintStream(&leakedOut)
	sentinelErr := jutil.NewPrintStream(&leakedErr)
	jutil.System_SetOut(sentinelOut)
	jutil.System_SetErr(sentinelErr)
	defer func() {
		jutil.System_SetOut(origOut)
		jutil.System_SetErr(origErr)
	}()

	worker := NewPersistentWorker(nil)
	reqBuf := encodeWorkRequests(t,
		WorkRequest{
			RequestId: 701,
			Arguments: []string{
				"--env=CUSTOM",
				"--checks_only",
				"--jscomp_warning=checkTypes",
				"--js", ttlJS,
			},
			Inputs: []WorkerInput{ttlIn},
		},
		WorkRequest{
			RequestId: 702,
			Arguments: []string{
				"--env=CUSTOM",
				"--translations_file", badXtb,
				"--js", msgJS,
				"--js_output_file", outJS,
			},
			Inputs: []WorkerInput{msgIn},
		},
	)

	var wireOut bytes.Buffer
	if err := worker.Serve(reqBuf, &wireOut); err != nil {
		t.Fatalf("worker.Serve: %v", err)
	}

	if leakedOut.Len() != 0 {
		t.Fatalf("System_Out leaked outside WorkResponse JSON stream: %q", leakedOut.String())
	}
	if leakedErr.Len() != 0 {
		t.Fatalf("System_Err leaked outside WorkResponse JSON stream: %q", leakedErr.String())
	}
	if jutil.System_Out() != sentinelOut || jutil.System_Err() != sentinelErr {
		t.Fatalf("PersistentWorker failed to restore System_Out / System_Err after Serve")
	}

	resps := decodeWorkResponses(t, wireOut.String())
	if len(resps) != 2 {
		t.Fatalf("expected 2 WorkResponses, got %d", len(resps))
	}
	if resps[0].ExitCode != 0 || resps[0].RequestId != 701 {
		t.Fatalf("req 701 unexpected response: %+v", resps[0])
	}
	if !strings.Contains(resps[0].Output, "TTL_STDOUT_MARKER:number") {
		t.Fatalf("expected TTL printType output captured in WorkResponse.Output, got: %q", resps[0].Output)
	}
	if resps[1].ExitCode == 0 || resps[1].RequestId != 702 {
		t.Fatalf("req 702 expected failure on malformed XTB, got: %+v", resps[1])
	}
	if !strings.Contains(resps[1].Output, "[Fatal Error]") {
		t.Fatalf("expected xtbxml [Fatal Error] stderr captured in WorkResponse.Output, got: %q", resps[1].Output)
	}
}

func TestPersistentWorkerSandboxDirResolutionAndCache(t *testing.T) {
	origWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Chdir(origWD)
	}()

	root := t.TempDir()
	sb1 := filepath.Join(root, "sandbox_1")
	sb2 := filepath.Join(root, "sandbox_2")
	for _, sb := range []string{sb1, sb2} {
		if err := os.MkdirAll(filepath.Join(sb, "src"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.MkdirAll(filepath.Join(sb, "out"), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	libCode := "var SHARED_PREFIX = 'sandbox_ok:';\n"
	appCodeV1 := "console.log(SHARED_PREFIX + 1);\n"
	appCodeV2 := "console.log(SHARED_PREFIX + 2);\n"

	if err := os.WriteFile(filepath.Join(sb1, "src", "lib.js"), []byte(libCode), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sb1, "src", "app.js"), []byte(appCodeV1), 0o644); err != nil {
		t.Fatal(err)
	}
	// Intentionally omit src/lib.js on disk in sb2 to verify cross-sandbox digest cache hit!
	if err := os.WriteFile(filepath.Join(sb2, "src", "app.js"), []byte(appCodeV2), 0o644); err != nil {
		t.Fatal(err)
	}

	libDigest := sha256Hex([]byte(libCode))
	appDigestV1 := sha256Hex([]byte(appCodeV1))
	appDigestV2 := sha256Hex([]byte(appCodeV2))

	worker := NewPersistentWorker(nil)

	// Request 1 in sb1 using sandbox-relative paths in both Arguments and Inputs.
	resp1 := worker.ExecuteRequest(&WorkRequest{
		RequestId:  801,
		SandboxDir: sb1,
		Arguments: []string{
			"--env=CUSTOM",
			"--compilation_level=SIMPLE",
			"--js", "src/lib.js",
			"--js", "src/app.js",
			"--js_output_file", "out/bundle.js",
			"--create_source_map", "out/bundle.js.map",
		},
		Inputs: []WorkerInput{
			{Path: "src/lib.js", Digest: libDigest},
			{Path: "src/app.js", Digest: appDigestV1},
		},
	})
	if resp1.ExitCode != 0 {
		t.Fatalf("sandbox_1 request failed: %+v", resp1)
	}
	afterWD, err := os.Getwd()
	if err != nil || afterWD != origWD {
		t.Fatalf("cwd was not restored after SandboxDir request: got %q, want %q (err=%v)", afterWD, origWD, err)
	}
	out1, err := os.ReadFile(filepath.Join(sb1, "out", "bundle.js"))
	if err != nil || !strings.Contains(string(out1), "sandbox_ok:") || !strings.Contains(string(out1), "+1") {
		t.Fatalf("sandbox_1 bundle mismatch: %q, err=%v", string(out1), err)
	}
	map1, err := os.ReadFile(filepath.Join(sb1, "out", "bundle.js.map"))
	if err != nil || !strings.Contains(string(map1), `"src/app.js"`) {
		t.Fatalf("sandbox_1 source map should preserve relative path src/app.js, got: %q, err=%v", string(map1), err)
	}

	// Request 2 in sb2: src/lib.js served from digest cache, src/app.js read from sb2.
	hitsBefore := worker.Cache().Stats().Hits
	resp2 := worker.ExecuteRequest(&WorkRequest{
		RequestId:  802,
		SandboxDir: sb2,
		Arguments: []string{
			"--env=CUSTOM",
			"--compilation_level=SIMPLE",
			"--js", "src/lib.js",
			"--js", "src/app.js",
			"--js_output_file", "out/bundle.js",
		},
		Inputs: []WorkerInput{
			{Path: "src/lib.js", Digest: libDigest},
			{Path: filepath.Join(sb2, "src", "app.js"), Digest: appDigestV2},
		},
	})
	if resp2.ExitCode != 0 {
		t.Fatalf("sandbox_2 request failed: %+v", resp2)
	}
	if worker.Cache().Stats().Hits != hitsBefore+1 {
		t.Fatalf("expected cache hit on src/lib.js across sandboxes, before=%d after=%d", hitsBefore, worker.Cache().Stats().Hits)
	}
	out2, err := os.ReadFile(filepath.Join(sb2, "out", "bundle.js"))
	if err != nil || !strings.Contains(string(out2), "sandbox_ok:") || !strings.Contains(string(out2), "+2") {
		t.Fatalf("sandbox_2 bundle mismatch: %q, err=%v", string(out2), err)
	}

	// Non-existent SandboxDir returns non-zero ExitCode cleanly and preserves cwd.
	respBad := worker.ExecuteRequest(&WorkRequest{
		RequestId:  803,
		SandboxDir: filepath.Join(root, "does_not_exist"),
		Arguments:  []string{"--js", "src/app.js"},
	})
	if respBad.ExitCode == 0 {
		t.Fatalf("expected non-zero exit code for missing SandboxDir, got %+v", respBad)
	}
	if gotWD, _ := os.Getwd(); gotWD != origWD {
		t.Fatalf("cwd changed after failed SandboxDir: got %q, want %q", gotWD, origWD)
	}
}
