// Port of com/google/javascript/jscomp/CommandLineRunner.java @ closure-compiler v20260415
// (covers AbstractCommandLineRunner.getBuiltinExterns + DefaultExterns.prepareExterns).
//
// Closure ships its default externs as a nested zip — `externs.zip`
// inside `closure-compiler.jar` — containing `es3.js`, `es5.js`,
// `es6.js`, `es6_collections.js` at the root and a `browser/` folder
// with `html5.js`, `w3c_dom1.js`, `w3c_event.js`, … . The Java side
// reads the zip via AbstractCommandLineRunner.getBuiltinExterns and
// orders the entries through DefaultExterns.prepareExterns: ES lang
// externs first, then a fixed BROWSER_EXTERN_DEP_ORDER, then the
// remaining browser/* entries in insertion order.
//
// jsgo's Makefile rule for $(CLOSURE_JAR) extracts the zip into
// $(CLOSURE_EXTERNS_DIR) (`.cache/closure-externs/`). This loader
// reads that directory (located via runtime.Caller so the test
// process and the CLI both find it without depending on cwd), parses
// each file with our regular parser, and returns the SCRIPT roots in
// Closure's documented load order. On any I/O / parse failure it
// returns "" so the caller can fall back to the hand-curated
// synthetic bundle in default_externs.go — keeping the type checker
// usable in environments where the externs cache hasn't been
// populated yet (fresh checkout without `make .cache/...`).
//
// The loader caches the merged source string (not the parsed AST)
// because TypedScopeCreator mutates per-AST node state during scope
// building; sharing one parsed externs tree across compiles produces
// cyclic JSType references and infinite recursion in IsSubtypeOf.
// Re-parsing per createExternScope mirrors the synthetic-bundle
// behaviour and matches Java's Compiler#parseExternFiles which also
// produces fresh extern AST per Compiler instance.

package jscomp

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
)

// builtinLangExterns mirrors DefaultExterns.BUILTIN_LANG_EXTERNS — the
// ES lang externs that ship at the zip root and apply for every
// CompilerOptions.Environment (BROWSER and CUSTOM alike).
var builtinLangExterns = []string{
	"es3.js",
	"es5.js",
	"es6.js",
	"es6_collections.js",
}

// browserExternDepOrder mirrors DefaultExterns.BROWSER_EXTERN_DEP_ORDER
// — the fixed-order BROWSER subset. Files in `browser/` not listed
// here are appended afterwards in alphabetical order (Closure's
// LinkedHashMap is populated by ZipInputStream.getNextEntry; entry
// order in the canonical jar is alphabetical).
var browserExternDepOrder = []string{
	"intl.js",
	"w3c_event.js",
	"w3c_event3.js",
	"gecko_event.js",
	"ie_event.js",
	"webkit_event.js",
	"w3c_device_sensor_event.js",
	"w3c_dom1.js",
	"w3c_dom2.js",
	"w3c_dom3.js",
	"w3c_dom4.js",
	"gecko_dom.js",
	"ie_dom.js",
	"webkit_dom.js",
	"w3c_css.js",
	"gecko_css.js",
	"ie_css.js",
	"webkit_css.js",
	"w3c_touch_event.js",
	"webstorage.js",
	"whatwg_console.js",
	"nonstandard_console.js",
	"w3c_rtc.js",
}

// externsCacheDirOnce / externsCacheDir memoise the on-disk location
// of the extracted externs. Resolution walks up from this source file
// (runtime.Caller) to find the first ancestor containing
// `.cache/closure-externs/es3.js`. Returns "" when the cache hasn't
// been populated — `make .cache/closure-compiler.jar` is the trigger.
var (
	externsCacheDirOnce sync.Once
	externsCacheDir     string
)

func resolveExternsCacheDir() string {
	externsCacheDirOnce.Do(func() {
		_, file, _, ok := runtime.Caller(0)
		if !ok {
			return
		}
		dir := filepath.Dir(file)
		for {
			candidate := filepath.Join(dir, ".cache", "closure-externs")
			if _, err := os.Stat(filepath.Join(candidate, "es3.js")); err == nil {
				externsCacheDir = candidate
				return
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				return
			}
			dir = parent
		}
	})
	return externsCacheDir
}

// loadDefaultExternsSource returns the concatenated externs source in
// Closure's documented load order — BUILTIN_LANG_EXTERNS, then
// BROWSER_EXTERN_DEP_ORDER, then the remaining `browser/*.js`. Each
// file's content is wrapped in `// === <name> ===` markers so parse
// errors stay attributable. Returns "" if the cache isn't populated
// or any required file fails to read — caller falls back to the
// synthetic bundle.
//
// Order matches AbstractCommandLineRunner.getBuiltinExterns +
// DefaultExterns.prepareExterns(BROWSER, …).
func loadDefaultExternsSource() string {
	cacheDir := resolveExternsCacheDir()
	if cacheDir == "" {
		return ""
	}
	var sb strings.Builder

	// 1. BUILTIN_LANG_EXTERNS — required; missing any aborts the load.
	for _, name := range builtinLangExterns {
		content, ok := readExternFile(filepath.Join(cacheDir, name))
		if !ok {
			return ""
		}
		sb.WriteString("\n// === ")
		sb.WriteString(name)
		sb.WriteString(" ===\n")
		sb.WriteString(content)
	}

	// 2. BROWSER_EXTERN_DEP_ORDER — applies when the browser/
	// subdirectory exists; we treat BROWSER as the default environment
	// (matches CommandLineRunner.getDefaultExterns which hard-codes
	// Environment.BROWSER).
	browserDir := filepath.Join(cacheDir, "browser")
	if _, err := os.Stat(browserDir); err != nil {
		// ES lang only is still useful — return what we have.
		return sb.String()
	}
	loaded := map[string]bool{}
	for _, name := range browserExternDepOrder {
		content, ok := readExternFile(filepath.Join(browserDir, name))
		if !ok {
			return ""
		}
		sb.WriteString("\n// === browser/")
		sb.WriteString(name)
		sb.WriteString(" ===\n")
		sb.WriteString(content)
		loaded[name] = true
	}

	// 3. Remaining browser/*.js — alphabetical, mirroring zip-entry
	// order in the canonical jar.
	entries, err := os.ReadDir(browserDir)
	if err != nil {
		return sb.String()
	}
	rest := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".js" {
			continue
		}
		if loaded[e.Name()] {
			continue
		}
		rest = append(rest, e.Name())
	}
	sort.Strings(rest)
	for _, name := range rest {
		content, ok := readExternFile(filepath.Join(browserDir, name))
		if !ok {
			// Optional externs failure is non-fatal.
			continue
		}
		sb.WriteString("\n// === browser/")
		sb.WriteString(name)
		sb.WriteString(" ===\n")
		sb.WriteString(content)
	}
	return sb.String()
}

// readExternFile reads a single externs file. Returns ("", false) on
// I/O error so callers can decide whether to fall back.
func readExternFile(path string) (string, bool) {
	src, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	return string(src), true
}

// defaultExternsSourceOnce / defaultExternsSourceCache memoise the
// merged source string. Reading + concatenation happens once; the
// caller (createExternScope) re-parses on each compile so each gets a
// fresh AST — matches the synthetic-bundle behaviour and avoids
// JSType cycles caused by sharing externs nodes across compiles.
var (
	defaultExternsSourceOnce  sync.Once
	defaultExternsSourceCache string
)

// upstreamExternsSource returns the concatenated upstream externs
// source (read from .cache/closure-externs/), or "" when the cache
// isn't populated. Caller is expected to fall back to the synthetic
// defaultExternsSrc.
func upstreamExternsSource() string {
	defaultExternsSourceOnce.Do(func() {
		defaultExternsSourceCache = loadDefaultExternsSource()
	})
	return defaultExternsSourceCache
}

// useUpstreamExterns reports whether createExternScope should consume
// the loader's output rather than the synthetic bundle. Gated on
// JSGO_USE_UPSTREAM_EXTERNS=1 because the upstream externs currently
// trigger infinite recursion in IsSubtypeOf on Array<T>/IArrayLike<T>
// (the recursion bug is in the type system, not the loader). Future
// agents can fix the recursion and remove the gate to make upstream
// the default.
func useUpstreamExterns() bool {
	return os.Getenv("JSGO_USE_UPSTREAM_EXTERNS") == "1"
}
