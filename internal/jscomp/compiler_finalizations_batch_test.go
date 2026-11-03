package jscomp

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type batchLocaleFixture struct {
	locale   string
	welcome  string
	cart     string
	cta      string
	currency string
}

type capturedLocaleOutput struct {
	locale      string
	chunks      []RenderedChunk
	variableMap []byte
	propertyMap []byte
	errors      []string
	warnings    []string
}

const batchTestProjectID = "ShopCheckoutApp"

var fiveLocaleFixtures = []batchLocaleFixture{
	{
		locale:   "en",
		welcome:  "Welcome back, <ph name=\"USER_NAME\"/>!",
		cart:     "Items in cart: <ph name=\"ITEM_COUNT\"/>",
		cta:      "Complete purchase",
		currency: "USD",
	},
	{
		locale:   "fr",
		welcome:  "Bon retour, <ph name=\"USER_NAME\"/> !",
		cart:     "Articles dans le panier : <ph name=\"ITEM_COUNT\"/>",
		cta:      "Finaliser la commande",
		currency: "EUR",
	},
	{
		locale:   "de",
		welcome:  "Willkommen zurueck, <ph name=\"USER_NAME\"/>!",
		cart:     "Artikel im Warenkorb: <ph name=\"ITEM_COUNT\"/>",
		cta:      "Kauf abschliessen",
		currency: "EUR",
	},
	{
		locale:   "es",
		welcome:  "Bienvenido de nuevo, <ph name=\"USER_NAME\"/>!",
		cart:     "Articulos en el carrito: <ph name=\"ITEM_COUNT\"/>",
		cta:      "Completar compra",
		currency: "EUR",
	},
	{
		locale:   "ja",
		welcome:  "Okaerinasai, <ph name=\"USER_NAME\"/>-san!",
		cart:     "Kaato nai no shohin: <ph name=\"ITEM_COUNT\"/>",
		cta:      "Konyu o kakutei",
		currency: "JPY",
	},
}

var eightLocaleFixtures = append(slices.Clone(fiveLocaleFixtures),
	batchLocaleFixture{
		locale:   "it",
		welcome:  "Bentornato, <ph name=\"USER_NAME\"/>!",
		cart:     "Articoli nel carrello: <ph name=\"ITEM_COUNT\"/>",
		cta:      "Completa acquisto",
		currency: "EUR",
	},
	batchLocaleFixture{
		locale:   "pt-BR",
		welcome:  "Bem-vindo de volta, <ph name=\"USER_NAME\"/>!",
		cart:     "Itens no carrinho: <ph name=\"ITEM_COUNT\"/>",
		cta:      "Concluir compra",
		currency: "BRL",
	},
	batchLocaleFixture{
		locale:   "ko",
		welcome:  "Dasi osin geoseul hwanyeonghamnida, <ph name=\"USER_NAME\"/>!",
		cart:     "Jangbaguni sangpum: <ph name=\"ITEM_COUNT\"/>",
		cta:      "Gumae wanryo",
		currency: "KRW",
	},
)

func buildLocaleXtbBundle(projectID string, lf batchLocaleFixture) *XtbMessageBundle {
	idGen := NewGoogleJsMessageIdGeneratorForProject(projectID)
	welcomeID := idGen.GenerateIdForMeaning("MSG_WELCOME", []JsMessage_Part{
		JsMessage_StringPart_Create("Welcome back, "),
		JsMessage_PlaceholderReference_CreateForCanonicalName("USER_NAME"),
		JsMessage_StringPart_Create("!"),
	})
	cartID := idGen.GenerateIdForMeaning("MSG_CART", []JsMessage_Part{
		JsMessage_StringPart_Create("Items in cart: "),
		JsMessage_PlaceholderReference_CreateForCanonicalName("ITEM_COUNT"),
	})
	ctaID := idGen.GenerateIdForMeaning("checkout_cta", []JsMessage_Part{
		JsMessage_StringPart_Create("Complete purchase"),
	})

	xtb := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE translationbundle SYSTEM "translationbundle.dtd">
<translationbundle lang=%q>
  <translation id=%q>%s</translation>
  <translation id=%q>%s</translation>
  <translation id=%q>%s</translation>
</translationbundle>`,
		lf.locale,
		welcomeID, lf.welcome,
		cartID, lf.cart,
		ctaID, lf.cta,
	)
	return NewXtbMessageBundleOpt(strings.NewReader(xtb), base.SomeString(projectID))
}

func buildEmptyStage2XtbBundle(projectID string) *XtbMessageBundle {
	const emptyXtb = `<?xml version="1.0" encoding="UTF-8"?>
<translationbundle lang="en"></translationbundle>`
	return NewXtbMessageBundleOpt(strings.NewReader(emptyXtb), base.SomeString(projectID))
}

func makeMultiChunkSources(extraModules int) ([]*SourceFile, []*JSChunk) {
	externs := []*SourceFile{
		SourceFile_FromCode("externs.js", `/** @externs */
/** @constructor */
function Object() {}
/**
 * @constructor
 * @param {*=} opt_str
 * @return {string}
 */
function String(opt_str) {}
var window = {};
/** @type {function(...*): void} */
window.emitChunkOutput = function(var_args) {};
/** @constructor */
function ExternalSink() {}
/** @param {string} payload */
ExternalSink.prototype.send = function(payload) {};
/** @type {!ExternalSink} */
window.sink;
`),
	}

	baseCode := `/**
 * @fileoverview
 * @provideGoog
 */
var goog = {};
/** @define {string} */
goog.LOCALE = 'en';
/**
 * @param {string} str
 * @param {!Object<string, string>=} opt_values
 * @return {string}
 */
goog.getMsg = function(str, opt_values) { return str; };

class PriceFormatter {
  /** @param {number} cents */
  constructor(cents) {
    /** @type {number} */
    this.centsAmount = cents;
  }
  /** @return {string} */
  formatAmount() {
    const unit = goog.LOCALE === 'ja' ? 'JPY' : (goog.LOCALE === 'en' ? 'USD' : 'EUR');
    return unit + ':' + (this.centsAmount / 100);
  }
}

class BannerModel {
  /** @param {string} userName */
  constructor(userName) {
    /** @type {string} */
    this.memberLabel = userName;
  }
  /** @return {string} */
  renderLabel() {
    /** @desc Welcome banner message */
    var MSG_WELCOME = goog.getMsg('Welcome back, {$userName}!', {'userName': this.memberLabel});
    return '[' + goog.LOCALE + '] ' + MSG_WELCOME;
  }
}

// Helper only called from checkout chunk -> eligible for CrossChunkCodeMotion!
/**
 * @param {string} a
 * @param {string} b
 * @return {string}
 */
function checkoutOnlyJoiner(a, b) {
  return a + ' :: ' + b;
}
`

	uiCode := `
/**
 * @param {string} user
 * @param {number} count
 * @return {string}
 */
function renderHeaderWidget(user, count) {
  const banner = new BannerModel(user);
  /** @desc Cart summary badge */
  var MSG_CART = goog.getMsg('Items in cart: {$itemCount}', {'itemCount': String(count)});
  return banner.renderLabel() + ' | ' + MSG_CART;
}
window.emitChunkOutput('ui', renderHeaderWidget('Ada', 3));
`

	checkoutCode := `
/**
 * @param {string} user
 * @param {number} count
 * @param {number} totalCents
 * @return {string}
 */
function renderCheckoutPanel(user, count, totalCents) {
  const formatter = new PriceFormatter(totalCents);
  /**
   * @desc Primary checkout button CTA
   * @meaning checkout_cta
   */
  var MSG_CTA = goog.getMsg('Complete purchase');
  const summary = checkoutOnlyJoiner(renderHeaderWidget(user, count), formatter.formatAmount());
  return summary + ' -> ' + MSG_CTA;
}
window.sink.send(renderCheckoutPanel('Grace', 5, 4999));
`

	analyticsCode := `
class TelemetryEvent {
  /** @param {string} action */
  constructor(action) {
    /** @type {string} */
    this.eventAction = action;
  }
  /** @return {string} */
  serializeEvent() {
    return goog.LOCALE + ':' + this.eventAction;
  }
}
window.emitChunkOutput('analytics', new TelemetryEvent('pageview').serializeEvent());
`

	baseChunk := NewJSChunk("base")
	baseChunk.Add(NewCompilerInput(SourceFile_FromCode("base.js", baseCode)))

	uiChunk := NewJSChunk("ui")
	uiChunk.AddDependency(baseChunk)
	uiChunk.Add(NewCompilerInput(SourceFile_FromCode("ui.js", uiCode)))

	checkoutChunk := NewJSChunk("checkout")
	checkoutChunk.AddDependency(baseChunk)
	checkoutChunk.AddDependency(uiChunk)
	checkoutChunk.Add(NewCompilerInput(SourceFile_FromCode("checkout.js", checkoutCode)))

	analyticsChunk := NewJSChunk("analytics")
	analyticsChunk.AddDependency(baseChunk)
	analyticsChunk.Add(NewCompilerInput(SourceFile_FromCode("analytics.js", analyticsCode)))

	targetChunks := []*JSChunk{uiChunk, checkoutChunk, analyticsChunk}
	for m := range extraModules {
		code := fmt.Sprintf(`
class FeatureCard_%d {
  /** @param {string} title @param {number} price */
  constructor(title, price) {
    /** @type {string} */
    this.cardTitle_%d = title;
    /** @type {!PriceFormatter} */
    this.cardPrice_%d = new PriceFormatter(price);
  }
  /** @return {string} */
  describeCard() {
    /**
     * @desc Localized CTA reused in feature card
     * @meaning checkout_cta
     */
    var MSG_CARD_CTA_%d = goog.getMsg('Complete purchase');
    return '[' + goog.LOCALE + '] ' + this.cardTitle_%d + '=' + this.cardPrice_%d.formatAmount() + ':' + MSG_CARD_CTA_%d;
  }
}
window.emitChunkOutput('card_%d', new FeatureCard_%d('item_%d', %d).describeCard());
`, m, m, m, m, m, m, m, m, m, m, (m+1)*250)
		ch := targetChunks[m%len(targetChunks)]
		ch.Add(NewCompilerInput(SourceFile_FromCode(fmt.Sprintf("feature_%02d.js", m), code)))
	}

	return externs, []*JSChunk{baseChunk, uiChunk, checkoutChunk, analyticsChunk}
}

func makeStageOptions() *CompilerOptions {
	opts := NewCompilerOptions()
	CompilationLevel_ADVANCED_OPTIMIZATIONS.SetOptionsForCompilationLevel(opts)
	opts.SetLanguageIn(CompilerOptions_LanguageMode_ECMASCRIPT_2020)
	opts.SetLanguageOut(CompilerOptions_LanguageMode_ECMASCRIPT_2015)
	opts.SetCheckTypes(true)
	opts.SetAmbiguateProperties(true)
	opts.SetDoLateLocalization(true)
	opts.SetMessageBundle(buildEmptyStage2XtbBundle(batchTestProjectID))
	opts.SetAlwaysGatherSourceMapInfo(true)
	opts.SetSourceMapDetailLevel(SourceMap_DetailLevel_ALL)
	opts.SetSourceMapFormat(SourceMap_Format_V3)
	opts.SetSourceMapIncludeSourcesContent(true)
	return opts
}

func compileStages1And2ToBytes(tb testing.TB, extraModules int) []byte {
	tb.Helper()
	c := NewCompiler()
	c.SetErrorManager(NewBlackHoleErrorManager())
	externs, chunks := makeMultiChunkSources(extraModules)
	opts := makeStageOptions()

	c.InitChunks(externs, chunks, opts)
	if c.HasErrors() {
		tb.Fatalf("InitChunks errors: %v", c.GetErrors())
	}
	c.ParseForCompilation()
	if c.HasErrors() {
		tb.Fatalf("ParseForCompilation errors: %v", c.GetErrors())
	}
	c.Stage1Passes()
	if c.HasErrors() {
		tb.Fatalf("Stage1Passes errors: %v", c.GetErrors())
	}
	c.Stage2Passes(CompilerOptions_SegmentOfCompilationToRun_OPTIMIZATIONS)
	if c.HasErrors() {
		tb.Fatalf("Stage2Passes errors: %v", c.GetErrors())
	}

	var buf bytes.Buffer
	c.SaveState(&buf)
	return buf.Bytes()
}

func collectLocaleResult(tb testing.TB, locale string, c *Compiler) capturedLocaleOutput {
	tb.Helper()
	rendered, err := c.RenderChunksAndSourceMaps(c.GetChunks(), RenderChunkOptions{
		GenerateSourceMaps: true,
		ChunkFileName: func(m *JSChunk) string {
			return fmt.Sprintf("%s.%s.js", m.GetName(), locale)
		},
		WrapperPrefix: func(m *JSChunk) string {
			return fmt.Sprintf("/* locale:%s chunk:%s */\n", locale, m.GetName())
		},
		NumWorkers: 2,
	})
	if err != nil {
		tb.Fatalf("locale %s RenderChunksAndSourceMaps failed: %v", locale, err)
	}

	var varMapBytes, propMapBytes []byte
	if vm := c.GetVariableMap(); vm != nil {
		varMapBytes = vm.ToBytes()
	}
	if pm := c.GetPropertyMap(); pm != nil {
		propMapBytes = pm.ToBytes()
	}

	errs := make([]string, len(c.GetErrors()))
	for i, e := range c.GetErrors() {
		errs[i] = e.String()
	}
	warns := make([]string, len(c.GetWarnings()))
	for i, w := range c.GetWarnings() {
		warns[i] = w.String()
	}

	return capturedLocaleOutput{
		locale:      locale,
		chunks:      rendered,
		variableMap: varMapBytes,
		propertyMap: propMapBytes,
		errors:      errs,
		warnings:    warns,
	}
}

func runSeparateFinalizationsReference(
	tb testing.TB,
	savedState []byte,
	extraModules int,
	specs []FinalizationLocaleSpec,
) []capturedLocaleOutput {
	tb.Helper()
	out := make([]capturedLocaleOutput, len(specs))
	for i, spec := range specs {
		c := NewCompiler()
		c.SetErrorManager(NewBlackHoleErrorManager())
		externs, chunks := makeMultiChunkSources(extraModules)
		opts := makeStageOptions()
		if spec.Locale != "" {
			opts.SetLocale(spec.Locale)
		}
		if spec.MessageBundle != nil {
			opts.SetMessageBundle(spec.MessageBundle)
		}
		if spec.ConfigureOptions != nil {
			spec.ConfigureOptions(opts)
		}
		c.InitChunks(externs, chunks, opts)
		c.RestoreState(bytes.NewReader(savedState))
		if c.HasErrors() {
			tb.Fatalf("separate locale %s RestoreState errors: %v", spec.Locale, c.GetErrors())
		}
		c.Stage3Passes()
		c.PerformPostCompilationTasks()
		out[i] = collectLocaleResult(tb, spec.Locale, c)
	}
	return out
}

func makeLocaleSpecs(fixtures []batchLocaleFixture) []FinalizationLocaleSpec {
	specs := make([]FinalizationLocaleSpec, len(fixtures))
	for i, lf := range fixtures {
		bundle := buildLocaleXtbBundle(batchTestProjectID, lf)
		loc := lf.locale
		specs[i] = FinalizationLocaleSpec{
			Locale:        loc,
			MessageBundle: bundle,
			ConfigureOptions: func(opts *CompilerOptions) {
				opts.SetSourceMapLocationMappings([]SourceMap_LocationMapping{
					NewSourceMap_PrefixLocationMapping("", fmt.Sprintf("/src/%s/", loc)),
				})
			},
		}
	}
	return specs
}

func assertLocaleOutputsIdentical(t *testing.T, want, got []capturedLocaleOutput) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("locale count = %d, want %d", len(got), len(want))
	}
	for i := range want {
		w, g := want[i], got[i]
		if g.locale != w.locale {
			t.Fatalf("idx %d: locale = %q, want %q", i, g.locale, w.locale)
		}
		if !slices.Equal(g.errors, w.errors) {
			t.Fatalf("locale %s: errors = %v, want %v", w.locale, g.errors, w.errors)
		}
		if !slices.Equal(g.warnings, w.warnings) {
			t.Fatalf("locale %s: warnings = %v, want %v", w.locale, g.warnings, w.warnings)
		}
		if !bytes.Equal(g.variableMap, w.variableMap) {
			t.Fatalf("locale %s: variableMap mismatch:\ngot:\n%s\nwant:\n%s", w.locale, g.variableMap, w.variableMap)
		}
		if !bytes.Equal(g.propertyMap, w.propertyMap) {
			t.Fatalf("locale %s: propertyMap mismatch:\ngot:\n%s\nwant:\n%s", w.locale, g.propertyMap, w.propertyMap)
		}
		if len(w.variableMap) == 0 || len(w.propertyMap) == 0 {
			t.Fatalf("locale %s: expected non-empty variable and property maps", w.locale)
		}
		if len(g.chunks) != len(w.chunks) {
			t.Fatalf("locale %s: chunk count = %d, want %d", w.locale, len(g.chunks), len(w.chunks))
		}
		for cIdx := range w.chunks {
			wc, gc := w.chunks[cIdx], g.chunks[cIdx]
			if gc.Chunk.GetName() != wc.Chunk.GetName() {
				t.Fatalf("locale %s chunk %d: name = %q, want %q", w.locale, cIdx, gc.Chunk.GetName(), wc.Chunk.GetName())
			}
			if gc.Code != wc.Code {
				t.Fatalf("locale %s chunk %s: JS code mismatch:\ngot:\n%s\nwant:\n%s",
					w.locale, wc.Chunk.GetName(), gc.Code, wc.Code)
			}
			if gc.SourceMap != wc.SourceMap {
				t.Fatalf("locale %s chunk %s: V3 SourceMap mismatch:\ngot:\n%s\nwant:\n%s",
					w.locale, wc.Chunk.GetName(), gc.SourceMap, wc.SourceMap)
			}
			if !strings.Contains(gc.SourceMap, "\"version\":3") {
				t.Fatalf("locale %s chunk %s: expected V3 source map, got: %s", w.locale, wc.Chunk.GetName(), gc.SourceMap)
			}
			if !strings.Contains(gc.SourceMap, fmt.Sprintf("/src/%s/", w.locale)) {
				t.Fatalf("locale %s chunk %s: source map missing per-locale prefix mapping", w.locale, wc.Chunk.GetName())
			}
		}
	}
}

func TestRunFinalizationsBatchMultiChunkParity(t *testing.T) {
	savedState := compileStages1And2ToBytes(t, 6)
	specs := makeLocaleSpecs(fiveLocaleFixtures)

	want := runSeparateFinalizationsReference(t, savedState, 6, specs)

	batchCompiler := NewCompiler()
	batchCompiler.SetErrorManager(NewBlackHoleErrorManager())
	externs, chunks := makeMultiChunkSources(6)
	batchCompiler.InitChunks(externs, chunks, makeStageOptions())

	got := make([]capturedLocaleOutput, 0, len(specs))
	err := batchCompiler.RunFinalizationsBatch(
		bytes.NewReader(savedState),
		specs,
		func(idx int, spec FinalizationLocaleSpec, c *Compiler) error {
			if idx != len(got) {
				t.Fatalf("unexpected callback idx = %d, want %d", idx, len(got))
			}
			if c.HasErrors() {
				t.Fatalf("locale %s unexpected errors: %v", spec.Locale, c.GetErrors())
			}
			got = append(got, collectLocaleResult(t, spec.Locale, c))
			return nil
		},
	)
	if err != nil {
		t.Fatalf("RunFinalizationsBatch returned error: %v", err)
	}

	assertLocaleOutputsIdentical(t, want, got)

	// Verify localized translations and goog.LOCALE substitution materialized per locale.
	for i, lf := range fiveLocaleFixtures {
		var combined strings.Builder
		for _, ch := range got[i].chunks {
			combined.WriteString(ch.Code)
		}
		allJS := combined.String()
		if !strings.Contains(allJS, lf.cta) {
			t.Fatalf("locale %s output missing CTA %q:\n%s", lf.locale, lf.cta, allJS)
		}
		if !strings.Contains(allJS, lf.currency+":") {
			t.Fatalf("locale %s output missing folded currency %q:\n%s", lf.locale, lf.currency, allJS)
		}
		if !strings.Contains(allJS, "["+lf.locale+"]") {
			t.Fatalf("locale %s output missing substituted goog.LOCALE %q:\n%s", lf.locale, lf.locale, allJS)
		}
	}
}

func TestRunFinalizationsBatchWithForcedStage3Compaction(t *testing.T) {
	savedState := compileStages1And2ToBytes(t, 4)
	specs := makeLocaleSpecs(fiveLocaleFixtures)
	want := runSeparateFinalizationsReference(t, savedState, 4, specs)

	for _, mode := range []string{"preorder_permute", "slide"} {
		t.Run(mode, func(t *testing.T) {
			prevHook := DebugAfterPassHook
			prevForce := slices.Clone(astLayoutForce)
			prevLayout := astLayoutEnv
			defer func() {
				DebugAfterPassHook = prevHook
				astLayoutForce = prevForce
				astLayoutEnv = prevLayout
			}()

			var activeCompiler *Compiler
			compactionsThisLocale := 0
			if mode == "preorder_permute" {
				astLayoutEnv = ""
				astLayoutForce = []string{PassNames_CROSS_CHUNK_CODE_MOTION, "ambiguateProperties", "renameProperties"}
				DebugAfterPassHook = func(c *Compiler, passName string) {
					if c == activeCompiler && passName == "SubstituteLocaleData" {
						c.initASTCompaction()
						c.beginASTCompaction()
					}
				}
			} else {
				astLayoutEnv = "slide"
				DebugAfterPassHook = func(c *Compiler, passName string) {
					if c == activeCompiler && passName == "SubstituteLocaleData" {
						c.initASTCompaction()
						c.beginASTCompaction()
					}
					if c == activeCompiler && (passName == PassNames_CROSS_CHUNK_CODE_MOTION || passName == "ambiguateProperties" || passName == "renameProperties") {
						c.compactAST(passName, true, &c.externsRoot, &c.jsRoot, nil)
						compactionsThisLocale++
					}
				}
			}

			batchCompiler := NewCompiler()
			batchCompiler.SetErrorManager(NewBlackHoleErrorManager())
			externs, chunks := makeMultiChunkSources(4)
			batchCompiler.InitChunks(externs, chunks, makeStageOptions())
			activeCompiler = batchCompiler

			var got []capturedLocaleOutput
			err := batchCompiler.RunFinalizationsBatch(
				bytes.NewReader(savedState),
				specs,
				func(idx int, spec FinalizationLocaleSpec, c *Compiler) error {
					count := c.astCompact.count
					if mode == "slide" {
						count = compactionsThisLocale
						compactionsThisLocale = 0
					}
					if count < 2 {
						t.Fatalf("mode=%s locale=%s: expected >= 2 forced Stage 3 AST compactions, got %d",
							mode, spec.Locale, count)
					}
					got = append(got, collectLocaleResult(t, spec.Locale, c))
					return nil
				},
			)
			if err != nil {
				t.Fatalf("mode=%s RunFinalizationsBatch failed: %v", mode, err)
			}

			assertLocaleOutputsIdentical(t, want, got)
		})
	}
}

func TestRunFinalizationsBatchEdgeCases(t *testing.T) {
	savedState := compileStages1And2ToBytes(t, 0)

	// 1. Nil reader returns error cleanly.
	{
		c := NewCompiler()
		c.SetErrorManager(NewBlackHoleErrorManager())
		externs, chunks := makeMultiChunkSources(0)
		c.InitChunks(externs, chunks, makeStageOptions())
		if err := c.RunFinalizationsBatch(nil, nil, nil); err == nil {
			t.Fatal("expected error for nil savedState reader")
		}
	}

	// 2. Empty locales slice restores once and succeeds with 0 callbacks.
	{
		c := NewCompiler()
		c.SetErrorManager(NewBlackHoleErrorManager())
		externs, chunks := makeMultiChunkSources(0)
		c.InitChunks(externs, chunks, makeStageOptions())
		called := false
		err := c.RunFinalizationsBatch(bytes.NewReader(savedState), nil, func(int, FinalizationLocaleSpec, *Compiler) error {
			called = true
			return nil
		})
		if err != nil || called {
			t.Fatalf("empty locales: err=%v, called=%v", err, called)
		}
	}

	// 3. Callback error aborts immediately and propagates exact error.
	{
		c := NewCompiler()
		c.SetErrorManager(NewBlackHoleErrorManager())
		externs, chunks := makeMultiChunkSources(0)
		c.InitChunks(externs, chunks, makeStageOptions())
		specs := makeLocaleSpecs(fiveLocaleFixtures[:3])
		stopErr := errors.New("stop after second locale")
		visited := 0
		err := c.RunFinalizationsBatch(bytes.NewReader(savedState), specs, func(idx int, _ FinalizationLocaleSpec, _ *Compiler) error {
			visited++
			if idx == 1 {
				return stopErr
			}
			return nil
		})
		if !errors.Is(err, stopErr) || visited != 2 {
			t.Fatalf("callback abort: got err=%v visited=%d, want stopErr and 2", err, visited)
		}
	}

	// 4. Per-locale strict missing-message error in locale #0 does not leak into locale #1.
	{
		c := NewCompiler()
		c.SetErrorManager(NewBlackHoleErrorManager())
		externs, chunks := makeMultiChunkSources(0)
		c.InitChunks(externs, chunks, makeStageOptions())

		specs := []FinalizationLocaleSpec{
			{
				Locale:        "fr",
				MessageBundle: buildEmptyStage2XtbBundle(batchTestProjectID), // missing messages!
				ConfigureOptions: func(opts *CompilerOptions) {
					opts.strictMessageReplacement = true
				},
			},
			makeLocaleSpecs(fiveLocaleFixtures[1:2])[0], // complete fr bundle, non-strict
		}

		var errCounts [2]int
		err := c.RunFinalizationsBatch(bytes.NewReader(savedState), specs, func(idx int, _ FinalizationLocaleSpec, comp *Compiler) error {
			errCounts[idx] = comp.GetErrorCount()
			return nil
		})
		if err != nil {
			t.Fatalf("unexpected batch error: %v", err)
		}
		if errCounts[0] == 0 {
			t.Fatal("expected strict missing-message errors in locale #0")
		}
		if errCounts[1] != 0 {
			t.Fatalf("locale #0 errors leaked into locale #1: got %d error(s)", errCounts[1])
		}
	}
}

func BenchmarkFinalizationsBatchVsSeparate(b *testing.B) {
	const benchExtraModules = 20
	savedState := compileStages1And2ToBytes(b, benchExtraModules)
	specs := makeLocaleSpecs(eightLocaleFixtures)

	b.Run("Separate8Locales", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			for _, spec := range specs {
				rhino.ResetNodeStore()
				c := NewCompiler()
				c.SetErrorManager(NewBlackHoleErrorManager())
				externs, chunks := makeMultiChunkSources(benchExtraModules)
				opts := makeStageOptions()
				opts.SetLocale(spec.Locale)
				opts.SetMessageBundle(spec.MessageBundle)
				if spec.ConfigureOptions != nil {
					spec.ConfigureOptions(opts)
				}
				c.InitChunks(externs, chunks, opts)
				c.RestoreState(bytes.NewReader(savedState))
				c.Stage3Passes()
				c.PerformPostCompilationTasks()
				if _, err := c.RenderChunksAndSourceMaps(c.GetChunks(), RenderChunkOptions{
					GenerateSourceMaps: true,
					NumWorkers:         2,
				}); err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	b.Run("Batch8Locales", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for range b.N {
			rhino.ResetNodeStore()
			c := NewCompiler()
			c.SetErrorManager(NewBlackHoleErrorManager())
			externs, chunks := makeMultiChunkSources(benchExtraModules)
			c.InitChunks(externs, chunks, makeStageOptions())
			err := c.RunFinalizationsBatch(
				bytes.NewReader(savedState),
				specs,
				func(_ int, _ FinalizationLocaleSpec, comp *Compiler) error {
					_, err := comp.RenderChunksAndSourceMaps(comp.GetChunks(), RenderChunkOptions{
						GenerateSourceMaps: true,
						NumWorkers:         2,
					})
					return err
				},
			)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
