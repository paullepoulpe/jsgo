package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const (
	// DefaultWorkerCacheMaxBytes bounds cached UTF-8 input file contents in a
	// persistent worker process (256 MiB).
	DefaultWorkerCacheMaxBytes int64 = 256 << 20

	// DefaultWorkerCacheMaxEntries bounds the number of cached input files in a
	// persistent worker process.
	DefaultWorkerCacheMaxEntries int = 4096
)

// WorkerInput represents a single file input entry in a Bazel JSON WorkRequest.
type WorkerInput struct {
	Path   string `json:"path"`
	Digest string `json:"digest"`
}

// WorkRequest represents a Bazel JSON persistent worker request.
type WorkRequest struct {
	Arguments  []string      `json:"arguments"`
	Inputs     []WorkerInput `json:"inputs"`
	RequestId  int           `json:"requestId"`
	Cancel     bool          `json:"cancel"`
	Verbosity  int           `json:"verbosity"`
	SandboxDir string        `json:"sandboxDir,omitempty"`
}

// WorkResponse represents a Bazel JSON persistent worker response.
type WorkResponse struct {
	ExitCode     int    `json:"exitCode"`
	Output       string `json:"output"`
	RequestId    int    `json:"requestId,omitempty"`
	WasCancelled bool   `json:"wasCancelled,omitempty"`
}

// WorkerCacheStats reports hit/miss and occupancy counters for a WorkerDigestCache.
type WorkerCacheStats struct {
	Entries int
	Bytes   int64
	Hits    uint64
	Misses  uint64
}

type workerCacheKey struct {
	path   string
	digest string
}

type workerCacheEntry struct {
	key    workerCacheKey
	code   string
	reload func() string
	prev   *workerCacheEntry
	next   *workerCacheEntry
}

// WorkerDigestCache is a bounded byte-budget and entry-count LRU cache storing
// immutable UTF-8 file contents keyed by (normalizedPath, cryptographicDigest).
// Live AST nodes are never stored in this cache.
type WorkerDigestCache struct {
	mu         sync.Mutex
	maxBytes   int64
	maxEntries int
	curBytes   int64
	entries    map[workerCacheKey]*workerCacheEntry
	byPath     map[string]*workerCacheEntry
	head       *workerCacheEntry // most recently used
	tail       *workerCacheEntry // least recently used
	hits       uint64
	misses     uint64
}

// NewWorkerDigestCache creates a bounded LRU digest cache.
func NewWorkerDigestCache(maxBytes int64, maxEntries int) *WorkerDigestCache {
	if maxBytes < 0 {
		maxBytes = 0
	}
	if maxEntries < 0 {
		maxEntries = 0
	}
	return &WorkerDigestCache{
		maxBytes:   maxBytes,
		maxEntries: maxEntries,
		entries:    make(map[workerCacheKey]*workerCacheEntry),
		byPath:     make(map[string]*workerCacheEntry),
	}
}

// Stats returns a snapshot of the cache's occupancy and hit/miss counters.
func (c *WorkerDigestCache) Stats() WorkerCacheStats {
	c.mu.Lock()
	defer c.mu.Unlock()
	return WorkerCacheStats{
		Entries: len(c.entries),
		Bytes:   c.curBytes,
		Hits:    c.hits,
		Misses:  c.misses,
	}
}

// InvalidatePath immediately evicts any cached entry associated with normPath.
func (c *WorkerDigestCache) InvalidatePath(normPath string) {
	c.mu.Lock()
	if existing, ok := c.byPath[normPath]; ok {
		c.removeEntryLocked(existing)
	}
	c.mu.Unlock()
}

func (c *WorkerDigestCache) evictIfDigestMismatch(normPath string, digest string) {
	c.mu.Lock()
	if existing, ok := c.byPath[normPath]; ok && existing.key.digest != digest {
		c.removeEntryLocked(existing)
	}
	c.mu.Unlock()
}

func (c *WorkerDigestCache) getReload(normPath string, digest string) (func() string, bool) {
	if digest == "" {
		c.InvalidatePath(normPath)
		return nil, false
	}
	c.mu.Lock()
	entry, ok := c.entries[workerCacheKey{path: normPath, digest: digest}]
	if !ok {
		if stale, hasStale := c.byPath[normPath]; hasStale {
			c.removeEntryLocked(stale)
		}
		c.misses++
		c.mu.Unlock()
		return nil, false
	}
	c.moveToFrontLocked(entry)
	c.hits++
	reload := entry.reload
	c.mu.Unlock()
	return reload, true
}

func (c *WorkerDigestCache) put(normPath string, digest string, code string) func() string {
	reload := func() string { return code }
	if digest == "" {
		c.InvalidatePath(normPath)
		return reload
	}
	size := int64(len(code))
	c.mu.Lock()
	defer c.mu.Unlock()

	if existing, ok := c.byPath[normPath]; ok {
		if existing.key.digest == digest {
			c.moveToFrontLocked(existing)
			return existing.reload
		}
		c.removeEntryLocked(existing)
	}

	if c.maxBytes <= 0 || c.maxEntries <= 0 || size > c.maxBytes {
		return reload
	}

	for (c.curBytes+size > c.maxBytes || len(c.entries) >= c.maxEntries) && c.tail != nil {
		c.removeEntryLocked(c.tail)
	}

	entry := &workerCacheEntry{
		key:    workerCacheKey{path: normPath, digest: digest},
		code:   code,
		reload: reload,
	}
	c.entries[entry.key] = entry
	c.byPath[normPath] = entry
	c.curBytes += size
	c.pushFrontLocked(entry)
	return reload
}

func (c *WorkerDigestCache) pushFrontLocked(e *workerCacheEntry) {
	e.prev = nil
	e.next = c.head
	if c.head != nil {
		c.head.prev = e
	} else {
		c.tail = e
	}
	c.head = e
}

func (c *WorkerDigestCache) moveToFrontLocked(e *workerCacheEntry) {
	if c.head == e {
		return
	}
	c.unlinkLocked(e)
	c.pushFrontLocked(e)
}

func (c *WorkerDigestCache) unlinkLocked(e *workerCacheEntry) {
	if e.prev != nil {
		e.prev.next = e.next
	} else {
		c.head = e.next
	}
	if e.next != nil {
		e.next.prev = e.prev
	} else {
		c.tail = e.prev
	}
	e.prev = nil
	e.next = nil
}

func (c *WorkerDigestCache) removeEntryLocked(e *workerCacheEntry) {
	c.unlinkLocked(e)
	delete(c.entries, e.key)
	if cur, ok := c.byPath[e.key.path]; ok && cur == e {
		delete(c.byPath, e.key.path)
	}
	c.curBytes -= int64(len(e.code))
}

// PersistentWorker executes Bazel JSON WorkRequests sequentially with hermetic
// per-request NodeStore resets and a bounded SHA-256 content-digest file cache.
type PersistentWorker struct {
	startupArgs []string
	cache       *WorkerDigestCache
}

// NewPersistentWorker creates a PersistentWorker with the default 256 MiB /
// 4,096-entry SHA-256 content-digest cache.
func NewPersistentWorker(startupArgs []string) *PersistentWorker {
	return NewPersistentWorkerWithCache(
		startupArgs,
		NewWorkerDigestCache(DefaultWorkerCacheMaxBytes, DefaultWorkerCacheMaxEntries),
	)
}

// NewPersistentWorkerWithCache creates a PersistentWorker using the provided cache.
func NewPersistentWorkerWithCache(startupArgs []string, cache *WorkerDigestCache) *PersistentWorker {
	if cache == nil {
		cache = NewWorkerDigestCache(DefaultWorkerCacheMaxBytes, DefaultWorkerCacheMaxEntries)
	}
	cleanedStartup, _ := extractWorkerArgs(expandWorkerStartupTokens(startupArgs))
	return &PersistentWorker{
		startupArgs: cleanedStartup,
		cache:       cache,
	}
}

// Cache returns the worker's SHA-256 digest input cache.
func (w *PersistentWorker) Cache() *WorkerDigestCache {
	return w.cache
}

// Serve reads JSON WorkRequests from in, executes them sequentially, and writes
// single-line JSON WorkResponses to out until EOF.
func (w *PersistentWorker) Serve(in io.Reader, out io.Writer) error {
	// Ensure no stale compilation store is active before serving requests,
	// and always release the compilation store when the worker loop exits.
	resetStandaloneStore := func() {
		rhino.LockStandaloneStore()
		rhino.ResetNodeStore()
		rhino.UnlockStandaloneStore()
	}
	resetStandaloneStore()
	defer resetStandaloneStore()

	dec := json.NewDecoder(in)
	for {
		var req WorkRequest
		if err := dec.Decode(&req); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		resp := w.ExecuteRequest(&req)
		line, err := json.Marshal(resp)
		if err != nil {
			return err
		}
		line = append(line, '\n')
		if _, err := out.Write(line); err != nil {
			return err
		}
	}
}

type workerSyncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *workerSyncBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	n, err := b.buf.Write(p)
	b.mu.Unlock()
	return n, err
}

func registerWorkerInputDigest(
	cache *WorkerDigestCache,
	reqDigests map[string]string,
	rawPath string,
	digest string,
	rawSandbox string,
	absSandbox string,
) {
	normPath := jutil.ParsePathString(rawPath)
	applyKey := func(k string, overwrite bool) {
		if k == "" {
			return
		}
		if digest == "" {
			cache.InvalidatePath(k)
			return
		}
		if !overwrite {
			if _, exists := reqDigests[k]; exists {
				return
			}
		}
		cache.evictIfDigestMismatch(k, digest)
		reqDigests[k] = digest
	}
	applyKey(normPath, true)
	if absSandbox == "" {
		return
	}
	sep := string(filepath.Separator)
	if filepath.IsAbs(normPath) {
		if rel, err := filepath.Rel(absSandbox, normPath); err == nil &&
			rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+sep) {
			applyKey(jutil.ParsePathString(rel), false)
		}
		return
	}
	applyKey(jutil.ParsePathString(filepath.Join(absSandbox, normPath)), false)
	if rawSbNorm := jutil.ParsePathString(rawSandbox); rawSbNorm != "" && !filepath.IsAbs(rawSbNorm) {
		if rel, err := filepath.Rel(rawSbNorm, normPath); err == nil &&
			rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+sep) {
			applyKey(jutil.ParsePathString(rel), false)
		}
	}
}

// ExecuteRequest runs a single Bazel WorkRequest hermetically and returns its WorkResponse.
func (w *PersistentWorker) ExecuteRequest(req *WorkRequest) WorkResponse {
	if req.Cancel {
		return WorkResponse{
			RequestId:    req.RequestId,
			WasCancelled: true,
		}
	}

	rhino.LockStandaloneStore()
	defer rhino.UnlockStandaloneStore()

	// Strictly reset the global NodeStore before AND after every compilation request
	// so AST blocks never leak across requests and are immediately eligible for GC.
	rhino.ResetNodeStore()
	defer rhino.ResetNodeStore()

	var absSandbox string
	if req.SandboxDir != "" {
		origDir, wdErr := os.Getwd()
		if wdErr != nil {
			return WorkResponse{
				ExitCode:  1,
				Output:    wdErr.Error() + "\n",
				RequestId: req.RequestId,
			}
		}
		var absErr error
		absSandbox, absErr = filepath.Abs(req.SandboxDir)
		if absErr != nil {
			return WorkResponse{
				ExitCode:  1,
				Output:    absErr.Error() + "\n",
				RequestId: req.RequestId,
			}
		}
		if chErr := os.Chdir(absSandbox); chErr != nil {
			return WorkResponse{
				ExitCode:  1,
				Output:    chErr.Error() + "\n",
				RequestId: req.RequestId,
			}
		}
		defer func() {
			_ = os.Chdir(origDir)
		}()
	}

	var reqDigests map[string]string
	if len(req.Inputs) > 0 {
		reqDigests = make(map[string]string, len(req.Inputs)*2)
		for i := range req.Inputs {
			inp := req.Inputs[i]
			if inp.Path == "" {
				continue
			}
			registerWorkerInputDigest(w.cache, reqDigests, inp.Path, inp.Digest, req.SandboxDir, absSandbox)
		}
	}

	reqArgs := w.buildRequestArgs(req.Arguments)

	var stdoutBuf workerSyncBuffer
	var stderrBuf workerSyncBuffer
	outPS := jutil.NewPrintStream(&stdoutBuf)
	errPS := jutil.NewPrintStream(&stderrBuf)
	emptyIn := bytes.NewReader(nil)

	exitCode := 0
	prevOut := jutil.System_Out()
	prevErr := jutil.System_Err()
	jutil.System_SetOut(errPS)
	jutil.System_SetErr(errPS)
	func() {
		defer func() {
			jutil.System_SetOut(prevOut)
			jutil.System_SetErr(prevErr)
		}()
		defer func() {
			if rec := recover(); rec != nil {
				thrown := jutil.AsThrowable(rec)
				if flagErr, ok := thrown.(jscomp.IFlagUsageException); ok {
					errPS.Println(flagErr.GetMessage())
				} else {
					jutil.Throwable_PrintStackTrace(thrown, errPS)
				}
				exitCode = 1
			}
		}()

		runner := NewCommandLineRunnerWithStreams(reqArgs, emptyIn, outPS, errPS)
		runner.SetExitCodeReceiver(func(code int) {
			exitCode = code
		})
		if !runner.GetCommandLineConfig().saveCompilationStateToFilename.IsPresent() &&
			!runner.flags.typedAstOutputFile.IsPresent() {
			runner.sourceFileCache = func(
				filename string,
				charset *jutil.Charset,
				kind rhino.StaticSourceFile_SourceKind,
			) (*jscomp.SourceFile, bool) {
				return w.lookupSourceFile(reqDigests, filename, charset, kind)
			}
		}

		if runner.ShouldRunCompiler() {
			runner.Run()
		} else if runner.HasErrors() {
			exitCode = 1
		}
	}()

	outPS.Flush()
	errPS.Flush()

	if exitCode < 0 {
		exitCode = 1
	}

	return WorkResponse{
		ExitCode:  exitCode,
		Output:    combineWorkerStreams(&stderrBuf.buf, &stdoutBuf.buf),
		RequestId: req.RequestId,
	}
}

func (w *PersistentWorker) lookupSourceFile(
	reqDigests map[string]string,
	filename string,
	charset *jutil.Charset,
	kind rhino.StaticSourceFile_SourceKind,
) (*jscomp.SourceFile, bool) {
	if filename == "" || isInferredZipEntryPath(filename) {
		return nil, false
	}
	if charset != nil && charset != jutil.StandardCharsets_UTF_8 {
		return nil, false
	}
	normPath := jutil.ParsePathString(filename)
	digest := ""
	if reqDigests != nil {
		digest = reqDigests[normPath]
	}
	if digest == "" {
		w.cache.InvalidatePath(normPath)
		return nil, false
	}
	if reload, ok := w.cache.getReload(normPath, digest); ok {
		return jscomp.SourceFile_Builder_().
			WithKind(kind).
			WithPath(filename).
			WithReloadableContent(reload).
			Build(), true
	}
	b, err := os.ReadFile(normPath)
	if err != nil || !utf8.Valid(b) {
		w.cache.InvalidatePath(normPath)
		return nil, false
	}
	reload := w.cache.put(normPath, digest, string(b))
	return jscomp.SourceFile_Builder_().
		WithKind(kind).
		WithPath(filename).
		WithReloadableContent(reload).
		Build(), true
}

func isInferredZipEntryPath(path string) bool {
	if strings.IndexByte(path, '!') < 0 {
		return false
	}
	bangSlash := "!" + jscomp.Platform_GetFileSeperator()
	idx := strings.Index(path, bangSlash)
	if idx < 0 {
		return false
	}
	zipPath := path[:idx]
	entryPath := path[idx+len(bangSlash):]
	return strings.HasSuffix(zipPath, ".zip") &&
		(strings.HasSuffix(entryPath, ".js") || strings.HasSuffix(entryPath, ".js.map"))
}

func combineWorkerStreams(stderrBuf, stdoutBuf *bytes.Buffer) string {
	errLen := stderrBuf.Len()
	outLen := stdoutBuf.Len()
	if errLen == 0 && outLen == 0 {
		return ""
	}
	if outLen == 0 {
		return stderrBuf.String()
	}
	if errLen == 0 {
		return stdoutBuf.String()
	}
	errBytes := stderrBuf.Bytes()
	needNewline := errBytes[len(errBytes)-1] != '\n'
	var sb strings.Builder
	sb.Grow(errLen + outLen + 1)
	sb.Write(errBytes)
	if needNewline {
		sb.WriteByte('\n')
	}
	sb.Write(stdoutBuf.Bytes())
	return sb.String()
}

func (w *PersistentWorker) buildRequestArgs(reqArgs []string) []string {
	needsSanitize := false
	for _, a := range reqArgs {
		if strings.Contains(a, "persistent_worker") || strings.Contains(a, "alsologtostderr") {
			needsSanitize = true
			break
		}
	}
	if !needsSanitize && len(w.startupArgs) == 0 {
		return reqArgs
	}
	cleanedReq := reqArgs
	if needsSanitize {
		cleanedReq, _ = extractWorkerArgs(reqArgs)
	}
	if len(w.startupArgs) == 0 {
		return cleanedReq
	}
	combined := make([]string, 0, len(w.startupArgs)+len(cleanedReq))
	combined = append(combined, w.startupArgs...)
	combined = append(combined, cleanedReq...)
	return combined
}

// RunPersistentWorker runs the Bazel JSON persistent worker request/response loop
// over in and out using any shared startup arguments.
func RunPersistentWorker(startupArgs []string, in io.Reader, out io.Writer) error {
	return NewPersistentWorker(startupArgs).Serve(in, out)
}

func runPersistentWorker(startupArgs []string, in io.Reader, out io.Writer) error {
	return RunPersistentWorker(startupArgs, in, out)
}

func parsePersistentWorkerStartupArgs(args []string) ([]string, bool) {
	mayHaveWorker := false
	for _, a := range args {
		if strings.Contains(a, "persistent_worker") ||
			strings.HasPrefix(a, "@") ||
			strings.HasPrefix(a, "--flagfile") ||
			strings.HasPrefix(a, "-flagfile") {
			mayHaveWorker = true
			break
		}
	}
	if !mayHaveWorker {
		return nil, false
	}
	expanded := expandWorkerStartupTokens(args)
	startupArgs, isWorker := extractWorkerArgs(expanded)
	if !isWorker {
		return nil, false
	}
	return startupArgs, true
}

func expandWorkerStartupTokens(args []string) []string {
	var out []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if len(arg) > 1 && arg[0] == '@' {
			file := jutil.NewFile(arg[1:])
			if file.Exists() {
				if lines, ex := jutil.Files_AsCharSource(file, jutil.StandardCharsets_UTF_8).TryReadLines(); ex == nil {
					for _, line := range lines {
						if line != "" {
							out = append(out, line)
						}
					}
					continue
				}
			}
			out = append(out, arg)
			continue
		}
		if (arg == "--flagfile" || arg == "-flagfile") && i+1 < len(args) {
			if raw, err := readUTF8FileStrict(jutil.Path_Of(args[i+1]).String()); err == nil {
				out = append(out, tokenizeFlagFileContent(raw)...)
				i++
				continue
			}
		} else if path, ok := cutFlagFilePrefix(arg); ok {
			if raw, err := readUTF8FileStrict(jutil.Path_Of(path).String()); err == nil {
				out = append(out, tokenizeFlagFileContent(raw)...)
				continue
			}
		}
		out = append(out, arg)
	}
	return out
}

func cutFlagFilePrefix(arg string) (string, bool) {
	if strings.HasPrefix(arg, "--flagfile=") {
		return arg[len("--flagfile="):], true
	}
	if strings.HasPrefix(arg, "-flagfile=") {
		return arg[len("-flagfile="):], true
	}
	return "", false
}

func extractWorkerArgs(tokens []string) ([]string, bool) {
	isWorker := false
	cleaned := make([]string, 0, len(tokens))
	for i := 0; i < len(tokens); i++ {
		tok := tokens[i]
		if tok == "--persistent_worker" || tok == "-persistent_worker" {
			if i+1 < len(tokens) {
				if val, isBool := parseWorkerBoolLiteral(tokens[i+1]); isBool {
					isWorker = val
					i++
					continue
				}
			}
			isWorker = true
			continue
		}
		if valStr, ok := cutNamedOptionValue(tok, "persistent_worker"); ok {
			if val, isBool := parseWorkerBoolLiteral(valStr); isBool {
				isWorker = val
				continue
			}
		}
		if tok == "--alsologtostderr" || tok == "-alsologtostderr" {
			if i+1 < len(tokens) {
				if _, isBool := parseWorkerBoolLiteral(tokens[i+1]); isBool {
					i++
				}
			}
			continue
		}
		if _, ok := cutNamedOptionValue(tok, "alsologtostderr"); ok {
			continue
		}
		cleaned = append(cleaned, tok)
	}
	return cleaned, isWorker
}

func cutNamedOptionValue(tok string, name string) (string, bool) {
	if strings.HasPrefix(tok, "--") {
		rest := tok[2:]
		if len(rest) > len(name) && strings.HasPrefix(rest, name) && rest[len(name)] == '=' {
			return rest[len(name)+1:], true
		}
	} else if strings.HasPrefix(tok, "-") {
		rest := tok[1:]
		if len(rest) > len(name) && strings.HasPrefix(rest, name) && rest[len(name)] == '=' {
			return rest[len(name)+1:], true
		}
	}
	return "", false
}

func parseWorkerBoolLiteral(s string) (val bool, ok bool) {
	if len(s) >= 2 && (s[0] == '\'' || s[0] == '"') && s[len(s)-1] == s[0] {
		s = s[1 : len(s)-1]
	}
	switch asciiLower(s) {
	case "true", "1", "yes", "on":
		return true, true
	case "false", "0", "no", "off":
		return false, true
	default:
		return false, false
	}
}
