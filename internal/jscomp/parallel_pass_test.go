package jscomp

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func compileWithParallelThreads(
	t *testing.T,
	externSources []string,
	inputSources []string,
	numThreads int,
	configure func(opts *CompilerOptions),
) (out string, errs []string, warns []string, formattedReport string) {
	t.Helper()

	externs := make([]*SourceFile, len(externSources))
	for i, code := range externSources {
		externs[i] = SourceFile_FromCode(fmt.Sprintf("extern_%d.js", i), code)
	}
	inputs := make([]*SourceFile, len(inputSources))
	for i, code := range inputSources {
		inputs[i] = SourceFile_FromCode(fmt.Sprintf("input_%d.js", i), code)
	}

	opts := NewCompilerOptions()
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(opts)
	WarningLevel_VERBOSE.SetOptionsForWarningLevel(opts)
	opts.SetCheckTypes(true)
	opts.SetClosurePass(true)
	opts.SetStrictModeInput(true)
	opts.SetNumParallelThreads(numThreads)
	if configure != nil {
		configure(opts)
	}

	var reportBuf bytes.Buffer
	comp := NewCompilerWithOutStream(jutil.NewPrintStream(&reportBuf))
	comp.DisableThreads()
	res := comp.Compile(externs, inputs, opts)
	for _, e := range res.Errors {
		errs = append(errs, e.String())
	}
	for _, w := range res.Warnings {
		warns = append(warns, w.String())
	}
	if len(res.Errors) == 0 {
		out = comp.ToSource()
	}
	return out, errs, warns, reportBuf.String()
}

func assertParallelParity(
	t *testing.T,
	externs []string,
	inputs []string,
	configure func(opts *CompilerOptions),
) {
	t.Helper()
	seqOut, seqErrs, seqWarns, seqReport := compileWithParallelThreads(t, externs, inputs, 1, configure)
	parOut, parErrs, parWarns, parReport := compileWithParallelThreads(t, externs, inputs, 4, configure)

	if strings.Join(seqErrs, "\n") != strings.Join(parErrs, "\n") {
		t.Fatalf("errors mismatch:\n--- seq (threads=1) ---\n%s\n--- par (threads=4) ---\n%s",
			strings.Join(seqErrs, "\n"), strings.Join(parErrs, "\n"))
	}
	if strings.Join(seqWarns, "\n") != strings.Join(parWarns, "\n") {
		t.Fatalf("warnings mismatch:\n--- seq (threads=1) ---\n%s\n--- par (threads=4) ---\n%s",
			strings.Join(seqWarns, "\n"), strings.Join(parWarns, "\n"))
	}
	if seqReport != parReport {
		t.Fatalf("formatted report mismatch:\n--- seq (threads=1) ---\n%s\n--- par (threads=4) ---\n%s",
			seqReport, parReport)
	}
	if seqOut != parOut {
		t.Fatalf("output JS mismatch:\n--- seq (threads=1) ---\n%s\n--- par (threads=4) ---\n%s",
			seqOut, parOut)
	}
}

func TestParallelPasses_EndToEndDeterminismAndRaceSafety(t *testing.T) {
	externs := []string{
		`/** @externs */
/** @type {!Object} */ var window;
/** @nosideeffects @param {number} x @return {number} */
function pureExternFn(x) {}
/** @constructor */ function RegExp(p, f) {}
/** @constructor @template T */ function Promise(fn) {}
/** @param {...*} args */ function sink(args) {}
`,
		`/** @externs */
/** @param {string} s @return {*} */ function eval(s) {}
`,
	}

	var inputs []string
	inputs = append(inputs, `/** @fileoverview Closure base stub */
var goog = {};
goog.provide = function(ns) {};
goog.module = function(ns) {};
goog.require = function(ns) {};
goog.exportSymbol = function(name, obj) { window[name] = obj; };
`)

	for i := 0; i < 24; i++ {
		inputs = append(inputs, fmt.Sprintf(`
/**
 * @param {number} a
 * @param {number} b
 * @return {number}
 */
function addHelper_%[1]d(a, b) {
  const casted = /** @type {number} */ (a);
  return (casted + 0) + (b * 1) + (2 + 3);
}

/**
 * @param {number} x
 * @return {number}
 */
function blockHelper_%[1]d(x) {
  let acc = x + 1;
  if (acc > 1000) {
    return acc - 1;
  }
  return acc * 2;
}

/** @param {number} x @return {number} */
var reassignedHelper_%[1]d = function(x) { return x + %[1]d; };

class Base_%[1]d {
  /** @param {number} v */
  constructor(v) {
    /** @type {number} */
    this.v = v;
  }
  /** @return {number} */
  compute() {
    return this.v + 1;
  }
}

class Sub_%[1]d extends Base_%[1]d {
  /** @param {number} v */
  constructor(v) {
    super(v + 2);
  }
  /** @override @return {number} */
  compute() {
    return super.compute() + (function(k) { return k + 7; })(%[1]d);
  }
}

/**
 * @param {number} n
 * @return {!Promise<number>}
 */
async function asyncCompute_%[1]d(n) {
  return addHelper_%[1]d(n, %[1]d);
}

/**
 * @param {number} n
 * @return {number}
 */
function runScript_%[1]d(n) {
  pureExternFn(n);
  const re = /abc%[1]d+/g;
  let total = new Sub_%[1]d(n).compute();
  switch (n) {
    case 1:
      total += blockHelper_%[1]d(addHelper_%[1]d(n, 2));
      break;
    case 2:
      total += (n > 0 ? addHelper_%[1]d(n, 4) : 0);
      break;
    default:
      total += blockHelper_%[1]d(n);
      break;
  }
  if (re.toString().length > 0) {
    total += 1;
  }
  return total;
}
`, i))
	}

	// Cross-script references, reassignments, goog.module, and eval forbidTemps with extra args.
	inputs = append(inputs, `
reassignedHelper_3 = function(x) { return x * 99; };

function callerWithEval(x) {
  eval("1");
  return addHelper_0(x, 2, sink(x));
}

var allResults = [];
for (var idx = 0; idx < 3; idx++) {
  allResults.push(
    runScript_0(idx) +
    runScript_5(idx) +
    runScript_12(idx) +
    runScript_23(idx) +
    reassignedHelper_3(idx) +
    callerWithEval(idx)
  );
  sink(asyncCompute_2(idx));
}
goog.exportSymbol('keptResults', allResults);
`)

	inputs = append(inputs, `goog.module('demo.mod.alpha');
class ModClass {
  /** @param {string} s */
  constructor(s) {
    this.s = s;
  }
  /** @return {string} */
  val() {
    return this.s + '_ok';
  }
}
exports = ModClass;
`)

	inputs = append(inputs, `goog.module('demo.mod.beta');
const ModClass = goog.require('demo.mod.alpha');
const inst = new ModClass('hello');
goog.exportSymbol('keptMod', inst.val());
`)

	assertParallelParity(t, externs, inputs, nil)
}

func TestParallelPasses_DiagnosticsAndHaltingErrorsOrder(t *testing.T) {
	externs := []string{
		`/** @externs */
/** @nosideeffects @param {number} x @return {number} */
function noSideFx(x) {}
/** @constructor */ function RegExp(p, f) {}
/** @param {...*} args */ function use(args) {}
`,
	}

	var inputs []string
	for i := 0; i < 12; i++ {
		inputs = append(inputs, fmt.Sprintf(`
/**
 * @param {number} x
 * @return {number}
 */
function missingRet_%[1]d(x) {
  if (x > 0) {
    return x + 1;
  }
}

function suspiciousAndDupCase_%[1]d(x) {
  noSideFx(x);
  var obj = {a: 1, a: 2};
  switch (x) {
    case 1:
      return obj.a;
    case 1:
      return 2;
  }
  if (x === NaN) {
    return 3;
  }
  return 4;
}
use(missingRet_%[1]d(1), suspiciousAndDupCase_%[1]d(2));
`, i))
	}

	assertParallelParity(t, externs, inputs, func(opts *CompilerOptions) {
		opts.SetWarningLevel(DiagnosticGroups_MISSING_RETURN, CheckLevel_WARNING)
		opts.SetWarningLevel(DiagnosticGroups_CHECK_USELESS_CODE, CheckLevel_WARNING)
		opts.SetWarningLevel(DiagnosticGroups_CHECK_REGEXP, CheckLevel_WARNING)
	})

	// Also verify halting error stop order when CheckMissingReturn is promoted to ERROR.
	assertParallelParity(t, externs, inputs, func(opts *CompilerOptions) {
		opts.SetWarningLevel(DiagnosticGroups_MISSING_RETURN, CheckLevel_ERROR)
	})
}

func TestParallelPasses_Stage2FunctionLocalParity(t *testing.T) {
	externs := []string{
		`/** @externs */
/** @type {!Object} */ var window;
/** @type {!Object} */ var Object;
/** @param {!Object} o @param {string} p @param {!Object} d */
Object.defineProperty = function(o, p, d) {};
/** @param {!Object} o @param {!Object} p */
Object.defineProperties = function(o, p) {};
/** @param {...*} args @return {*} */ function extEffect(args) {}
`,
	}

	var inputs []string
	// Script 0 calls hoisted functions declared in later scripts before their declarations!
	inputs = append(inputs, `
var GLOBAL_CONST = 42;
var GLOBAL_MUT = extEffect('init');
window['earlyHoisted'] = hoistedFromLate_5(3) + hoistedFromLate_2(7);
`)

	for i := 0; i < 10; i++ {
		inputs = append(inputs, fmt.Sprintf(`
function leafBlockInl_%[1]d(x) {
  var t = x + %[1]d;
  if (t > 10) {
    extEffect(t);
    return t * 2;
  }
  return t + 1;
}

function midConflict_%[1]d(a, b = 1) {
  return leafBlockInl_%[1]d(a) + b;
}

function exerciseDecomposeAndVars_%[1]d(obj, n, unusedMiddle, unusedTrail1, unusedTrail2 = 5) {
  var aliasGlobalConst = GLOBAL_CONST;
  var aliasGlobalMut = GLOBAL_MUT;
  var localOnce = extEffect(n) + aliasGlobalConst;
  var unusedLocal = 999;
  var r1 = extEffect('lhs', n) + leafBlockInl_%[1]d(localOnce);
  var r2 = obj?.nested?.(leafBlockInl_%[1]d(n + 1));
  var r3 = (n > 0 ? midConflict_%[1]d(n) : 0) + (aliasGlobalMut ? 1 : 0);
  return r1 + (r2 || 0) + r3;
}

var accHolder_%[1]d = {
  get propA_%[1]d() { return %[1]d; },
  set propB_%[1]d(v) { extEffect(v); }
};
Object.defineProperty(accHolder_%[1]d, 'dynProp_%[1]d', {
  get: function() { return %[1]d + 1; }
});
Object.defineProperties(accHolder_%[1]d, {
  multiProp_%[1]d: {
    set: function(v) { extEffect(v); }
  }
});

class AccClass_%[1]d extends (extEffect({get extGetter_%[1]d() { return 1; }}), Object) {
  [extEffect('k_%[1]d')]() {
    return ({set compSetter_%[1]d(v) { extEffect(v); }});
  }
  get propB_%[1]d() { return %[1]d; }
}

function hoistedFromLate_%[1]d(k) {
  const localConst = GLOBAL_CONST + %[1]d;
  return leafBlockInl_%[1]d(k + localConst) + (new AccClass_%[1]d().propB_%[1]d);
}

window['entry_%[1]d'] = exerciseDecomposeAndVars_%[1]d(window['o'], %[1]d, 'mid') + accHolder_%[1]d.propA_%[1]d;
window['gdecomp_%[1]d'] = extEffect('g_%[1]d') + leafBlockInl_%[1]d(%[1]d);
`, i))
	}

	// 1. ADVANCED_OPTIMIZATIONS (InlineVariables ALL, InlineFunctions Decompose + Conflicts,
	//    GatherGetterAndSetterProperties, RemoveUnusedCode with removeGlobals + FunctionArgs).
	assertParallelParity(t, externs, inputs, nil)

	// 2. SIMPLE_OPTIMIZATIONS (InlineVariables LOCALS_ONLY + RemoveUnusedCode local-only).
	assertParallelParity(t, externs, inputs, func(opts *CompilerOptions) {
		CompilationLevel_SIMPLE_OPTIMIZATIONS.SetOptionsForCompilationLevel(opts)
		opts.SetCheckTypes(false)
	})

	// 3. InlineVariables CONSTANTS_ONLY mode.
	assertParallelParity(t, externs, inputs, func(opts *CompilerOptions) {
		opts.SetInlineVariables(false)
		opts.SetInlineLocalVariables(false)
		opts.SetInlineConstantVars(true)
	})
}

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
