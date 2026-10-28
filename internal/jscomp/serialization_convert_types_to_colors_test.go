package jscomp

// Go-only: direct-lowering equivalence and zero-allocation fast-skip tests for
// ConvertTypesToColors, JSDocSerializer, and Token <-> NodeKind translation.

import (
	"io"
	"slices"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type astColorSnapshot struct {
	token         rhino.Token
	color         *colors.Color
	colorFromCast bool
	hasJsdoc      bool
	jsdocSummary  string
}

func collectAstColorSnapshots(root rhino.Node) []astColorSnapshot {
	var out []astColorSnapshot
	var walk func(n rhino.Node)
	walk = func(n rhino.Node) {
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			walk(c)
		}
		jsdoc := n.GetJSDocInfo()
		var summary string
		if jsdoc != nil {
			summary = jsdoc.ToStringVerbose()
		}
		out = append(out, astColorSnapshot{
			token:         n.GetToken(),
			color:         n.GetColor(),
			colorFromCast: n.IsColorFromTypeCast(),
			hasJsdoc:      jsdoc != nil,
			jsdocSummary:  summary,
		})
	}
	walk(root)
	return out
}

func assertIdenticalColorGraphs(
	t *testing.T,
	directSnaps, protoSnaps []astColorSnapshot,
	directReg, protoReg *colors.ColorRegistry,
	includeDebugMismatches bool,
) {
	t.Helper()
	if len(directSnaps) != len(protoSnaps) {
		t.Fatalf("node snapshot count mismatch: direct=%d proto=%d", len(directSnaps), len(protoSnaps))
	}

	seenColors := make(map[colors.ColorId]*colors.Color)
	protoColors := make(map[colors.ColorId]*colors.Color)

	var recordColor func(cDirect, cProto *colors.Color)
	recordColor = func(cDirect, cProto *colors.Color) {
		if (cDirect == nil) != (cProto == nil) {
			t.Fatalf("nil color mismatch: direct=%v proto=%v", cDirect, cProto)
		}
		if cDirect == nil {
			return
		}
		if !cDirect.Equals(cProto) {
			t.Fatalf("Color mismatch for id %s:\ndirect: %+v\nproto:  %+v", cDirect.GetId(), cDirect, cProto)
		}
		if !slices.Equal(cDirect.GetOwnProperties(), cProto.GetOwnProperties()) {
			t.Fatalf("ownProperties order mismatch for %s: direct=%v proto=%v",
				cDirect.GetId(), cDirect.GetOwnProperties(), cProto.GetOwnProperties())
		}
		if len(cDirect.GetUnionElements()) != len(cProto.GetUnionElements()) {
			t.Fatalf("unionElements count mismatch for %s", cDirect.GetId())
		}
		for i := range cDirect.GetUnionElements() {
			if !cDirect.GetUnionElements()[i].Equals(cProto.GetUnionElements()[i]) {
				t.Fatalf("unionElements[%d] mismatch for %s: direct=%v proto=%v",
					i, cDirect.GetId(), cDirect.GetUnionElements()[i], cProto.GetUnionElements()[i])
			}
		}
		if len(cDirect.GetInstanceColors()) != len(cProto.GetInstanceColors()) {
			t.Fatalf("instanceColors count mismatch for %s", cDirect.GetId())
		}
		for i := range cDirect.GetInstanceColors() {
			if !cDirect.GetInstanceColors()[i].Equals(cProto.GetInstanceColors()[i]) {
				t.Fatalf("instanceColors[%d] mismatch for %s", i, cDirect.GetId())
			}
		}
		if len(cDirect.GetPrototypes()) != len(cProto.GetPrototypes()) {
			t.Fatalf("prototypes count mismatch for %s", cDirect.GetId())
		}
		for i := range cDirect.GetPrototypes() {
			if !cDirect.GetPrototypes()[i].Equals(cProto.GetPrototypes()[i]) {
				t.Fatalf("prototypes[%d] mismatch for %s", i, cDirect.GetId())
			}
		}

		id := cDirect.GetId()
		if _, visited := seenColors[id]; visited {
			return
		}
		seenColors[id] = cDirect
		protoColors[id] = cProto

		for i := range cDirect.GetUnionElements() {
			recordColor(cDirect.GetUnionElements()[i], cProto.GetUnionElements()[i])
		}
		for i := range cDirect.GetInstanceColors() {
			recordColor(cDirect.GetInstanceColors()[i], cProto.GetInstanceColors()[i])
		}
		for i := range cDirect.GetPrototypes() {
			recordColor(cDirect.GetPrototypes()[i], cProto.GetPrototypes()[i])
		}

		directSups := directReg.GetDisambiguationSupertypes(cDirect).ToSlice()
		protoSups := protoReg.GetDisambiguationSupertypes(cProto).ToSlice()
		if len(directSups) != len(protoSups) {
			t.Fatalf("disambiguationSupertypes count mismatch for %s: direct=%d proto=%d",
				id, len(directSups), len(protoSups))
		}
		for i := range directSups {
			recordColor(directSups[i], protoSups[i])
		}
	}

	for i := range directSnaps {
		d, p := directSnaps[i], protoSnaps[i]
		if d.token != p.token {
			t.Fatalf("node[%d] token mismatch: direct=%v proto=%v", i, d.token, p.token)
		}
		if d.colorFromCast != p.colorFromCast {
			t.Fatalf("node[%d] colorFromCast mismatch: direct=%v proto=%v", i, d.colorFromCast, p.colorFromCast)
		}
		if d.hasJsdoc != p.hasJsdoc || d.jsdocSummary != p.jsdocSummary {
			t.Fatalf("node[%d] jsdoc mismatch:\ndirect: %s\nproto:  %s", i, d.jsdocSummary, p.jsdocSummary)
		}
		recordColor(d.color, p.color)
	}

	for reqID := range colors.ColorRegistry_REQUIRED_IDS.All() {
		recordColor(directReg.Get(reqID), protoReg.Get(reqID))
	}

	if includeDebugMismatches {
		dMis := directReg.GetMismatchLocationsForDebugging()
		pMis := protoReg.GetMismatchLocationsForDebugging()
		for id := range seenColors {
			dLocs := dMis.Get(id).ToSlice()
			pLocs := pMis.Get(id).ToSlice()
			if !slices.Equal(dLocs, pLocs) {
				t.Fatalf("mismatch locations differ for %s: direct=%v proto=%v", id, dLocs, pLocs)
			}
		}
	}
}

func runReferenceProtoRoundTrip(
	compiler *Compiler,
	externsRoot, jsRoot rhino.Node,
	opts *SerializationOptions,
	initRuntimeLibs bool,
) *colors.ColorRegistry {
	externsAndJsRoot := jsRoot.GetParent()

	stringPoolBuilder := StringPool_Builder_()
	serializeJstypes := SerializeTypesToPointers_Create(compiler, stringPoolBuilder, opts)
	serializeJstypes.GatherTypesOnAst(externsAndJsRoot)

	typePool := serializeJstypes.GetTypePool()
	stringPool := stringPoolBuilder.Build()

	colorPoolBuilder := ColorPool_Builder_()
	if initRuntimeLibs {
		compiler.InitRuntimeLibraryTypedAsts(colorPoolBuilder)
	}
	shard := colorPoolBuilder.AddShard(typePool, stringPool)
	colorPool := colorPoolBuilder.Build()

	var applyWalk func(n rhino.Node)
	applyWalk = func(n rhino.Node) {
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			applyWalk(c)
		}
		oldType, hadBeforeCast, jsdoc := n.TakeTypesAndJSDocForColorConversion()
		if jsdoc != nil {
			n.SetJSDocInfo(JSDocSerializer_ConvertJSDocInfoForOptimizations(jsdoc))
		}
		if oldType != nil {
			if ptr, ok := serializeJstypes.jstypeReconserializer.typePointers[oldType]; ok {
				if col := shard.GetColor(int(ptr)); col != nil {
					n.SetColor(col)
				}
			}
		}
		if hadBeforeCast {
			n.SetColorFromTypeCast()
		}
	}
	applyWalk(externsAndJsRoot)
	convertTypesToColors_forDetachedRoots(compiler, externsAndJsRoot, applyWalk)
	compiler.ClearJSTypeRegistry()
	if initRuntimeLibs {
		compiler.SetColorRegistry(colorPool.GetRegistry())
	}
	return colorPool.GetRegistry()
}

func prepareCheckedCompilerForColorsTest(t *testing.T, externsCode, srcCode string) *Compiler {
	t.Helper()
	compiler := NewCompilerWithOutStream(jutil.NewPrintStream(io.Discard))
	errMgr := newTypeCheckDump_errorManager()
	compiler.SetErrorManager(errMgr)
	options := typeCheckDump_options()
	externsFile := SourceFile_FromCode("externs.js", externsCode)
	inputFile := SourceFile_FromCode("input.js", srcCode)
	compiler.Init([]*SourceFile{externsFile}, []*SourceFile{inputFile}, options)
	compiler.Parse()
	serializationDump_typeCheck(compiler, options)
	NewRemoveCastNodes(compiler).Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	return compiler
}

func TestConvertTypesToColors_DirectMatchesProtoRoundTrip(t *testing.T) {
	const externs = `
/** @constructor */ function Object() {}
/** @constructor */ function Function() {}
/** @constructor @param {*=} x @return {string} */ function String(x) {}
/** @type {function(): string} */ String.prototype.trim = function() {};
/** @constructor @param {*=} x @return {number} */ function Number(x) {}
/** @type {function(number=): string} */ Number.prototype.toFixed = function(d) {};
/** @constructor @param {*=} x @return {boolean} */ function Boolean(x) {}
/** @constructor @param {*=} x @return {symbol} */ function Symbol(x) {}
/** @constructor @param {*=} x @return {bigint} */ function BigInt(x) {}
`

	const src = `
/**
 * @closurePrimitive {asserts.truthy}
 * @param {*} cond
 * @return {*}
 */
function assertTruthy(cond) { return cond; }

/** @interface */
class Serializable {
  /** @return {string} */
  serialize() {}
}

/** @interface @extends {Serializable} */
class NamedSerializable {
  /** @return {string} */
  getName() {}
}

/** @abstract @implements {NamedSerializable} */
class BaseEntity {
  constructor() {
    /** @type {number} */
    this.id = 1;
    /** @type {?BaseEntity} */
    this.parent = null;
  }
  /** @abstract @return {string} */
  getName() {}
  /** @override @return {string} */
  serialize() { return String(this.id).trim(); }
}

class ConcreteEntity extends BaseEntity {
  /** @param {string} label */
  constructor(label) {
    super();
    /** @type {string} */
    this.label = label;
    this.p1 = 1;
    this.p2 = 2;
    this.p3 = 3;
    this.p4 = 4;
    this.p5 = 5;
    this.p6 = 6;
    this.p7 = 7;
    this.p8 = 8;
    this.p9 = 9;
    this.p10 = 10;
  }
  /** @override @noinline @return {string} */
  getName() { return this.label; }
}

/** @enum {number} */
const Kind = {
  ALPHA: 1,
  BETA: 2,
};

/** @typedef {(string|number|boolean|ConcreteEntity|null)} */
let ComplexUnion;

/**
 * @param {ComplexUnion|Kind} val
 * @param {number} n
 * @modifies {arguments}
 * @suppress {partialAlias,messageConventions}
 * @return {string}
 */
function processValue(val, n) {
  assertTruthy(n);
  const casted = /** @type {!ConcreteEntity} */ (val);
  const sum = casted.p1 + casted.p2 + casted.p3 + casted.p4 + casted.p5 +
      casted.p6 + casted.p7 + casted.p8 + casted.p9 + casted.p10;
  // Intentional type mismatch to exercise invalidating types & mismatch locations.
  /** @type {!ConcreteEntity} */
  const bad = "not-an-entity";
  void bad;
  return casted.getName() + casted.serialize() + sum.toFixed(1) + String(Kind.ALPHA);
}

const instance = new ConcreteEntity("ok");
processValue(instance, 42);
`

	for _, tc := range []struct {
		name               string
		includeDebug       bool
		withRuntimeLibPool bool
	}{
		{name: "production_with_runtime_libs", includeDebug: false, withRuntimeLibPool: true},
		{name: "debug_and_validation_with_runtime_libs", includeDebug: true, withRuntimeLibPool: true},
		{name: "single_shard_no_runtime_libs", includeDebug: true, withRuntimeLibPool: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			opts := SerializationOptions_Builder_().
				SetIncludeDebugInfo(tc.includeDebug).
				SetRunValidation(tc.includeDebug).
				Build()

			cDirect := prepareCheckedCompilerForColorsTest(t, externs, src)
			var directReg *colors.ColorRegistry
			if tc.withRuntimeLibPool {
				NewConvertTypesToColors(cDirect, opts).
					Process(cDirect.GetExternsRoot(), cDirect.GetJsRoot())
				directReg = cDirect.GetColorRegistry()
			} else {
				// Exercise single-shard direct JSType lowering (0 proto shards).
				root := cDirect.GetJsRoot().GetParent()
				ser := SerializeTypesToPointers_Create(cDirect, nil, opts)
				ser.gatherTypesForColors(root)
				cpb := ColorPool_Builder_()
				cpb.addJSTypeShard(ser.jstypeReconserializer)
				cp := cpb.Build()
				var cb ConvertTypesToColors_RemoveTypesAndApplyColors
				cb.jstypeReconserializer = ser.jstypeReconserializer
				cb.offsetToColor = ser.jstypeReconserializer.offsetToColor
				cb.walk(root)
				convertTypesToColors_forDetachedRoots(cDirect, root, cb.walk)
				cDirect.ClearJSTypeRegistry()
				directReg = cp.GetRegistry()
			}

			cProto := prepareCheckedCompilerForColorsTest(t, externs, src)
			protoReg := runReferenceProtoRoundTrip(
				cProto,
				cProto.GetExternsRoot(),
				cProto.GetJsRoot(),
				opts,
				tc.withRuntimeLibPool,
			)

			directSnaps := collectAstColorSnapshots(cDirect.GetJsRoot().GetParent())
			protoSnaps := collectAstColorSnapshots(cProto.GetJsRoot().GetParent())
			assertIdenticalColorGraphs(
				t,
				directSnaps,
				protoSnaps,
				directReg,
				protoReg,
				tc.includeDebug,
			)
		})
	}
}

func TestJSDocSerializer_NonOptimizationFastSkipZeroAllocs(t *testing.T) {
	b := rhino.JSDocInfo_Builder_()
	b.RecordParameter("x", jsDocSerializer_placeholderType)
	b.RecordReturnType(jsDocSerializer_placeholderType)
	b.RecordOverride()
	b.RecordSuppressions(jutil.NewLinkedHashSet("checkTypes", "extraRequire"))
	nonOptDoc := b.Build()
	if nonOptDoc == nil {
		t.Fatal("expected non-nil JSDocInfo")
	}

	pool := StringPool_Builder_()
	if got := JSDocSerializer_ConvertJSDocInfoForOptimizations(nonOptDoc); got != nil {
		t.Fatalf("expected nil converted JSDocInfo, got %v", got)
	}
	if got := JSDocSerializer_SerializeJsdoc(nonOptDoc, pool); got != nil {
		t.Fatalf("expected nil OptimizationJsdoc, got %v", got)
	}

	var reusableBuilder rhino.JSDocInfo_Builder
	allocsConvert := testing.AllocsPerRun(100, func() {
		if JSDocSerializer_ConvertJSDocInfoForOptimizations(nonOptDoc) != nil {
			t.Fatal("unexpected result")
		}
		if jsDocSerializer_convertJSDocInfoWithBuilder(nonOptDoc, &reusableBuilder) != nil {
			t.Fatal("unexpected result")
		}
	})
	if allocsConvert != 0 {
		t.Fatalf("expected 0 allocs on non-optimization ConvertJSDocInfoForOptimizations, got %v", allocsConvert)
	}

	allocsSerialize := testing.AllocsPerRun(100, func() {
		if JSDocSerializer_SerializeJsdoc(nonOptDoc, pool) != nil {
			t.Fatal("unexpected result")
		}
	})
	if allocsSerialize != 0 {
		t.Fatalf("expected 0 allocs on non-optimization SerializeJsdoc, got %v", allocsSerialize)
	}
}

func TestTokenNodeKindTables_RoundTripAllValidKinds(t *testing.T) {
	sf := SourceFile_FromCode("test.js", "")
	tmpl := sf.getOrCreateSourceInfoTemplate()
	spBuilder := StringPool_Builder_()
	numPtr := int32(spBuilder.Put("12345678901234567890"))
	identPtr := int32(spBuilder.Put("propName"))
	sp := spBuilder.Build()

	deser := &ScriptNodeDeserializer{
		sourceFile: sf,
		stringPool: sp,
		filePool:   []*SourceFile{sf},
	}
	runner := newScriptNodeDeserializer_Runner(deser)
	ser := NewTypedAstSerializer(nil, SerializationOptions_Builder_().Build())

	for kind := NodeKind(1); kind < NodeKind_UNRECOGNIZED; kind++ {
		strPtr := identPtr
		if kind == NodeKind_BIGINT_LITERAL {
			strPtr = numPtr
		}
		n := runner.deserializeSingleNode(kind, strPtr, 3.5, identPtr, identPtr, tmpl)
		got := ser.KindTranslator(n)
		if got != kind {
			t.Fatalf("round-trip mismatch for %s: deserialized token %v serialized back to %s",
				kind, n.GetToken(), got)
		}
	}
}
