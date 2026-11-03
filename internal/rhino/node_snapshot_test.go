package rhino

import (
	"fmt"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestNodeStoreSnapshotDisallowsNilAndStatic(t *testing.T) {
	assertPanics := func(name string, fn func()) {
		t.Helper()
		defer func() {
			if recover() == nil {
				t.Fatalf("%s: expected panic, got none", name)
			}
		}()
		fn()
	}

	var nilStore *NodeStore
	assertPanics("nil.Snapshot", func() {
		nilStore.Snapshot()
	})
	assertPanics("staticStore.Snapshot", func() {
		staticStore.Snapshot()
	})

	s1 := NewNodeStore()
	s2 := NewNodeStore()
	snap1 := s1.Snapshot()

	assertPanics("nil.RestoreSnapshot", func() {
		nilStore.RestoreSnapshot(snap1)
	})
	assertPanics("staticStore.RestoreSnapshot", func() {
		staticStore.RestoreSnapshot(snap1)
	})
	assertPanics("RestoreSnapshot(nil)", func() {
		s1.RestoreSnapshot(nil)
	})
	assertPanics("mismatched store RestoreSnapshot", func() {
		s2.RestoreSnapshot(snap1)
	})
}

func TestNodeStoreSnapshotMutationsAndMultipleRewinds(t *testing.T) {
	store := NewNodeStore()
	prev := SetCurrentNodeStore(store)
	defer SetCurrentNodeStore(prev)

	registry := NewJSTypeRegistry(nil)
	numType := registry.GetNativeType(JSTypeNative_NUMBER_TYPE)
	strType := registry.GetNativeType(JSTypeNative_STRING_TYPE)

	colA := colors.Color_SingleBuilder().SetId(colors.ColorId_FromUnsigned(1001)).Build()
	colB := colors.Color_SingleBuilder().SetId(colors.ColorId_FromUnsigned(1002)).Build()

	root := IR_Root()
	script := IR_Script()
	root.AddChildToBack(script)

	// Node with unshared + shared cold prop chain:
	// cloneSrc creates shared tail items; then we prepend unshared cold props on richNode.
	cloneSrc := IR_Name("sharedPropSeed")
	cloneSrc.PutProp(Node_Prop_NON_JSDOC_COMMENT, "shared-comment")
	cloneSrc.PutIntProp(Node_Prop_ACCESS_MODIFIER, 2)
	richNode := cloneSrc.CloneNode() // marks cloneSrc's prop chain shared
	richNode.SetString("richIdent")
	richNode.SetOriginalName(jutil.StringPtr("origRichIdent"))
	richNode.PutProp(Node_Prop_REFLECTED_OBJECT, "unshared-initial")
	richNode.PutIntProp(Node_Prop_SIDE_EFFECT_FLAGS, 100) // > maxInline (15): forces kindFlagHasExtra
	richNode.PutBooleanProp(Node_Prop_IS_CONSTANT_NAME, true)
	richNode.PutBooleanProp(Node_Prop_IS_IN_CLOSURE_UNAWARE_SUBTREE, true) // coldProps bit in extBox
	richNode.SetLength(200_000)                                            // > 0xFFFF: stored in extBox.length
	richNode.SetColor(colA)
	shadowSubtree := IR_ScriptOfNodes(IR_ExprResult(IR_Name("shadowBody")))
	richNode.PutProp(Node_Prop_CLOSURE_UNAWARE_SHADOW, shadowSubtree)
	script.AddChildToBack(IR_Var(richNode))

	// Node with inline original name (kindFlagOrigNameIsString) and JSType.
	typedNode := IR_Name("inlineOrig")
	typedNode.SetOriginalName(jutil.StringPtr("inlineOrig"))
	typedNode.SetJSType(numType)
	numLit := Node_NewNumber(3.14159)
	typedVar := IR_Var(typedNode)
	typedNode.AddChildToFront(numLit)
	script.AddChildToBack(typedVar)

	// Allocate > 1 full block of nodes (> 16,384) to span multiple nodeBlocks,
	// nodeColumns, and pagedColumns.
	const extraVars = blockLen + 512
	bulkNodes := make([]Node, extraVars)
	for i := range extraVars {
		nm := IR_Name(fmt.Sprintf("v_%d", i))
		if i%97 == 0 {
			nm.SetOriginalName(jutil.StringPtr(fmt.Sprintf("orig_v_%d", i)))
			nm.SetColor(colB)
		}
		bulkNodes[i] = nm
		script.AddChildToBack(IR_ExprResult(nm))
	}

	wantCount := store.NodeCount()
	wantBase := store.base
	wantTree := root.ToStringTree()
	richID := richNode.ID()
	typedID := typedNode.ID()
	shadowID := shadowSubtree.ID()

	snap := store.Snapshot()

	verifyState := func(round string) {
		t.Helper()
		if got := store.base; got != wantBase {
			t.Fatalf("%s: base = %d, want %d", round, got, wantBase)
		}
		if got := store.NodeCount(); got != wantCount {
			t.Fatalf("%s: NodeCount = %d, want %d", round, got, wantCount)
		}
		if richNode.ID() != richID || typedNode.ID() != typedID {
			t.Fatalf("%s: node IDs changed", round)
		}
		if gotTree := root.ToStringTree(); gotTree != wantTree {
			t.Fatalf("%s: tree mismatch", round)
		}
		if richNode.GetString() != "richIdent" {
			t.Fatalf("%s: richNode string = %q, want richIdent", round, richNode.GetString())
		}
		if orig := richNode.GetOriginalName(); orig == nil || *orig != "origRichIdent" {
			t.Fatalf("%s: richNode originalName = %v, want origRichIdent", round, orig)
		}
		if got := richNode.GetProp(Node_Prop_REFLECTED_OBJECT); got != "unshared-initial" {
			t.Fatalf("%s: REFLECTED_OBJECT = %v, want unshared-initial", round, got)
		}
		if got := richNode.GetProp(Node_Prop_NON_JSDOC_COMMENT); got != "shared-comment" {
			t.Fatalf("%s: NON_JSDOC_COMMENT = %v, want shared-comment", round, got)
		}
		if got := richNode.GetIntProp(Node_Prop_ACCESS_MODIFIER); got != 2 {
			t.Fatalf("%s: ACCESS_MODIFIER = %d, want 2", round, got)
		}
		if got := richNode.GetIntProp(Node_Prop_SIDE_EFFECT_FLAGS); got != 100 {
			t.Fatalf("%s: SIDE_EFFECT_FLAGS = %d, want 100", round, got)
		}
		if !richNode.GetBooleanProp(Node_Prop_IS_CONSTANT_NAME) {
			t.Fatalf("%s: lost IS_CONSTANT_NAME", round)
		}
		if !richNode.GetBooleanProp(Node_Prop_IS_IN_CLOSURE_UNAWARE_SUBTREE) {
			t.Fatalf("%s: lost IS_IN_CLOSURE_UNAWARE_SUBTREE", round)
		}
		if got := richNode.GetLength(); got != 200_000 {
			t.Fatalf("%s: length = %d, want 200000", round, got)
		}
		if richNode.GetColor() != colA {
			t.Fatalf("%s: richNode color = %v, want colA", round, richNode.GetColor())
		}
		if sh, ok := richNode.GetProp(Node_Prop_CLOSURE_UNAWARE_SHADOW).(Node); !ok || sh.ID() != shadowID || sh.GetFirstChild().GetFirstChild().GetString() != "shadowBody" {
			t.Fatalf("%s: shadow subtree corrupted: %v", round, richNode.GetProp(Node_Prop_CLOSURE_UNAWARE_SHADOW))
		}
		if typedNode.GetString() != "inlineOrig" {
			t.Fatalf("%s: typedNode string = %q", round, typedNode.GetString())
		}
		if orig := typedNode.GetOriginalName(); orig == nil || *orig != "inlineOrig" {
			t.Fatalf("%s: typedNode originalName = %v", round, orig)
		}
		if typedNode.GetJSType() != numType {
			t.Fatalf("%s: typedNode JSType = %v, want numType", round, typedNode.GetJSType())
		}
		if numLit.GetDouble() != 3.14159 {
			t.Fatalf("%s: numLit double = %v", round, numLit.GetDouble())
		}
		for i, nm := range bulkNodes {
			wantName := fmt.Sprintf("v_%d", i)
			if nm.GetString() != wantName {
				t.Fatalf("%s: bulkNodes[%d] string = %q, want %q", round, i, nm.GetString(), wantName)
			}
			orig := nm.GetOriginalName()
			if i%97 == 0 {
				wantOrig := fmt.Sprintf("orig_v_%d", i)
				if orig == nil || *orig != wantOrig {
					t.Fatalf("%s: bulkNodes[%d] orig = %v, want %q", round, i, orig, wantOrig)
				}
				if nm.GetColor() != colB {
					t.Fatalf("%s: bulkNodes[%d] lost colB", round, i)
				}
			} else {
				if orig != nil {
					t.Fatalf("%s: bulkNodes[%d] unexpected orig = %q", round, i, *orig)
				}
				if nm.GetColor() != nil {
					t.Fatalf("%s: bulkNodes[%d] unexpected color", round, i)
				}
			}
		}
	}

	verifyState("initial")

	for rewind := 1; rewind <= 3; rewind++ {
		round := fmt.Sprintf("rewind #%d", rewind)

		// Mutate unshared prop in-place, remove shared prop, clear cold prop, and clone richNode (markShared).
		richNode.PutProp(Node_Prop_REFLECTED_OBJECT, fmt.Sprintf("mutated-%d", rewind))
		_ = richNode.CloneNode()
		richNode.RemoveProp(Node_Prop_NON_JSDOC_COMMENT)
		richNode.PutBooleanProp(Node_Prop_IS_CONSTANT_NAME, false)
		richNode.PutBooleanProp(Node_Prop_IS_IN_CLOSURE_UNAWARE_SUBTREE, false)
		richNode.PutIntProp(Node_Prop_SIDE_EFFECT_FLAGS, 3)
		richNode.SetLength(12)
		richNode.SetString(fmt.Sprintf("renamed_%d", rewind))

		// Register brand-new colors (across handleTable block boundary on rewind #1)
		// and mutate existing node colors/types.
		extraColors := 64
		if rewind == 1 {
			extraColors = tableBlockLen + 25
		}
		var lastExtraCol *colors.Color
		for cIdx := range extraColors {
			lastExtraCol = colors.Color_SingleBuilder().
				SetId(colors.ColorId_FromUnsigned(int64(rewind*1_000_000 + cIdx + 2000))).
				Build()
			bulkNodes[cIdx].SetColor(lastExtraCol)
		}
		richNode.SetColor(lastExtraCol)
		typedNode.SetJSType(strType)
		typedNode.SetString("renamedInline")
		numLit.SetDouble(99.5)
		typedVar.Detach()

		// Populate previously empty columnPages in origNames, and clear existing ones.
		for i := range 300 {
			if i%97 == 0 {
				bulkNodes[i].SetOriginalName(nil)
			} else {
				bulkNodes[i].SetOriginalName(jutil.StringPtr("temp_orig"))
			}
		}

		// Allocate an entire extra block of new nodes and link into the AST.
		var postSnapNode Node
		for i := range blockLen + 100 {
			n := IR_Name(fmt.Sprintf("temp_alloc_%d_%d", rewind, i))
			if i == 0 {
				postSnapNode = n
			}
			script.AddChildToBack(IR_ExprResult(n))
		}

		store.RestoreSnapshot(snap)
		verifyState(round)

		// Any node allocated in blocks past snap must fail loudly if dereferenced.
		func() {
			defer func() {
				if recover() == nil {
					t.Fatalf("%s: expected panic accessing node from retired block", round)
				}
			}()
			_ = postSnapNode.GetString()
		}()

		// Verify newly allocated nodes after RestoreSnapshot start cleanly zeroed
		// and existing colors/types still resolve handles properly.
		fresh := IR_Name("after_restore")
		if !fresh.GetParent().IsNil() || !fresh.GetFirstChild().IsNil() || fresh.GetOriginalName() != nil || fresh.GetColor() != nil {
			t.Fatalf("%s: fresh node after restore was not clean", round)
		}
		fresh.SetColor(colA)
		if fresh.GetColor() != colA {
			t.Fatalf("%s: SetColor(colA) after restore failed", round)
		}
		// Restore back to exact snap before next iteration.
		store.RestoreSnapshot(snap)
		verifyState(round + " clean")
	}

	// Verify another store created afterwards starts strictly above all blocks ever used.
	otherStore := NewNodeStore()
	SetCurrentNodeStore(otherStore)
	otherNode := IR_Name("isolated")
	if otherNode.ID() <= bulkNodes[len(bulkNodes)-1].ID()+uint32(blockLen) {
		t.Fatalf("new store reused block range: otherNode ID = %d", otherNode.ID())
	}
}

func TestNodeStoreSnapshotRewindAfterCompaction(t *testing.T) {
	store := NewNodeStore()
	prev := SetCurrentNodeStore(store)
	defer SetCurrentNodeStore(prev)

	registry := NewJSTypeRegistry(nil)
	boolType := registry.GetNativeType(JSTypeNative_BOOLEAN_TYPE)
	fnColor := colors.Color_SingleBuilder().SetId(colors.ColorId_FromUnsigned(3001)).Build()

	// Create > 1 block of dead nodes up front so both Slide and Permute
	// relocate live nodes to completely different IDs and block counts.
	for range blockLen + 2000 {
		_ = IR_Name("dead_prefix")
	}

	paramTypeNode := IR_Name("string")
	retTypeNode := Node_NewStringWithToken(Token_PIPE, "|")
	retTypeNode.AddChildToBack(IR_Name("number"))
	retTypeNode.AddChildToBack(IR_Name("boolean"))
	ttlBodyNode := IR_Name("T")

	docBuilder := NewJSDocInfo_Builder()
	docBuilder.RecordParameter("arg", NewJSTypeExpression(paramTypeNode, "test.js"))
	docBuilder.RecordReturnType(NewJSTypeExpression(retTypeNode, "test.js"))
	docBuilder.RecordTypeTransformation("Out", ttlBodyNode)
	fnDoc := docBuilder.Build()

	fnName := IR_Name("targetFn")
	fnName.SetOriginalName(jutil.StringPtr("origTargetFn"))
	fnParam := IR_Name("arg")
	fnParam.SetJSType(boolType)
	fnRetExpr := IR_Name("arg")
	fnRetExpr.SetColor(fnColor)
	fn := IR_Function(fnName, IR_ParamList(fnParam), IR_BlockOfNode(IR_ReturnNodeWithExpr(fnRetExpr)))
	fn.SetJSDocInfo(fnDoc)

	shadowAst := IR_ScriptOfNodes(IR_Var(IR_Name("innerShadow")))
	fn.PutProp(Node_Prop_CLOSURE_UNAWARE_SHADOW, shadowAst)
	fn.PutProp(Node_Prop_REFLECTED_OBJECT, "keepMe")

	script := IR_ScriptOfNodes(fn)
	root := IR_Root(script)

	// Dead nodes after the live tree too.
	for range 4000 {
		_ = IR_Var(IR_Name("dead_suffix"))
	}

	origBase := store.base
	origCount := store.NodeCount()
	origRootID := root.ID()
	origFnID := fn.ID()
	origShadowID := shadowAst.ID()
	origParamTypeID := paramTypeNode.ID()
	origRetTypeID := retTypeNode.ID()
	origTTLNodeID := ttlBodyNode.ID()
	origTree := root.ToStringTree()

	snap := store.Snapshot()

	verifyRestored := func(label string) {
		t.Helper()
		if store.base != origBase {
			t.Fatalf("%s: store.base = %d, want %d", label, store.base, origBase)
		}
		if store.NodeCount() != origCount {
			t.Fatalf("%s: NodeCount = %d, want %d", label, store.NodeCount(), origCount)
		}
		if root.ID() != origRootID || fn.ID() != origFnID {
			t.Fatalf("%s: root/fn IDs mismatch", label)
		}
		if gotTree := root.ToStringTree(); gotTree != origTree {
			t.Fatalf("%s: tree mismatch:\ngot:\n%s\nwant:\n%s", label, gotTree, origTree)
		}
		if orig := fnName.GetOriginalName(); orig == nil || *orig != "origTargetFn" {
			t.Fatalf("%s: fnName originalName = %v, want origTargetFn", label, orig)
		}
		if fnParam.GetJSType() != boolType {
			t.Fatalf("%s: fnParam lost JSType", label)
		}
		if fnRetExpr.GetColor() != fnColor {
			t.Fatalf("%s: fnRetExpr lost Color", label)
		}
		if got := fn.GetProp(Node_Prop_REFLECTED_OBJECT); got != "keepMe" {
			t.Fatalf("%s: REFLECTED_OBJECT = %v, want keepMe", label, got)
		}
		sh, ok := fn.GetProp(Node_Prop_CLOSURE_UNAWARE_SHADOW).(Node)
		if !ok || sh.ID() != origShadowID || sh.GetFirstChild().GetFirstChild().GetString() != "innerShadow" {
			t.Fatalf("%s: CLOSURE_UNAWARE_SHADOW = %v (want ID %d)", label, sh, origShadowID)
		}
		doc := fn.GetJSDocInfo()
		if doc != fnDoc {
			t.Fatalf("%s: JSDocInfo pointer changed", label)
		}
		if got := doc.GetParameterType("arg").GetRoot(); got.ID() != origParamTypeID || got.GetString() != "string" {
			t.Fatalf("%s: param JSDoc type root ID = %d, want %d", label, got.ID(), origParamTypeID)
		}
		if got := doc.GetReturnType().GetRoot(); got.ID() != origRetTypeID || got.GetFirstChild().GetString() != "number" {
			t.Fatalf("%s: return JSDoc type root ID = %d, want %d", label, got.ID(), origRetTypeID)
		}
		transforms := doc.GetTypeTransformationEntries()
		if len(transforms) != 1 || transforms[0].TypeNode.ID() != origTTLNodeID || transforms[0].TypeNode.GetString() != "T" {
			t.Fatalf("%s: JSDoc type transformation node not restored: %+v", label, transforms)
		}
	}

	// Cycle 1: Compactor.SlideWithThreads -> RestoreSnapshot
	{
		k := NewCompactor()
		k.Mark(root)
		k.SlideWithThreads(4)
		compactedRoot := k.Forward(root)
		compactedFn := k.Forward(fn)
		k.Finish()

		if store.base == origBase || compactedRoot.ID() == origRootID {
			t.Fatal("Slide did not relocate store base / root ID")
		}
		if fnDoc.GetReturnType().GetRoot().ID() == origRetTypeID {
			t.Fatal("Slide did not forward JSDoc return type root")
		}

		store.RestoreSnapshot(snap)
		verifyRestored("after Slide rewind")

		// Compacted ID must now be unmapped and panic if touched.
		func() {
			defer func() {
				if recover() == nil {
					t.Fatal("expected panic reading post-Slide node ID after RestoreSnapshot")
				}
			}()
			_ = compactedFn.GetFirstChild()
		}()
	}

	// Cycle 2: Compactor.RecordPreorder + Permute -> mutate compacted AST -> RestoreSnapshot
	{
		k := NewCompactor()
		k.RecordPreorder()
		k.Mark(root)
		k.Permute(4)
		compactedRoot := k.Forward(root)
		compactedFn := k.Forward(fn)
		k.Finish()

		if store.base == origBase || compactedRoot.ID() == origRootID {
			t.Fatal("Permute did not relocate store base / root ID")
		}
		if fnDoc.GetTypeTransformationEntries()[0].TypeNode.ID() == origTTLNodeID {
			t.Fatal("Permute did not forward JSDoc type transformation node")
		}
		compactedFn.PutProp(Node_Prop_REFLECTED_OBJECT, "after-permute")
		compactedFn.GetFirstChild().SetString("renamedAfterPermute")

		store.RestoreSnapshot(snap)
		verifyRestored("after Permute rewind")

		func() {
			defer func() {
				if recover() == nil {
					t.Fatal("expected panic reading post-Permute node ID after RestoreSnapshot")
				}
			}()
			_ = compactedRoot.GetFirstChild()
		}()
	}
}
