package rhino

import (
	"math/big"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/colors"
)

type testStaticSourceFile struct {
	name string
}

func (f testStaticSourceFile) GetName() string { return f.name }
func (f testStaticSourceFile) GetKind() StaticSourceFile_SourceKind {
	return StaticSourceFile_SourceKind_STRONG
}
func (f testStaticSourceFile) IsStrong() bool             { return true }
func (f testStaticSourceFile) IsWeak() bool               { return false }
func (f testStaticSourceFile) IsExtern() bool             { return false }
func (f testStaticSourceFile) IsNonCode() bool            { return false }
func (f testStaticSourceFile) IsTypeScriptSource() bool   { return false }
func (f testStaticSourceFile) IsClosureUnawareCode() bool { return false }
func (f testStaticSourceFile) GetLineOffset(int) int      { return 0 }
func (f testStaticSourceFile) GetLineOfOffset(int) int    { return 1 }
func (f testStaticSourceFile) GetColumnOfOffset(int) int  { return 0 }

func TestCaptureAndRewindStaticBaseline(t *testing.T) {
	LockStandaloneStore()
	defer UnlockStandaloneStore()
	ResetNodeStore()

	CaptureStaticBaseline()
	if !HasStaticBaseline() {
		t.Fatal("expected HasStaticBaseline() == true after CaptureStaticBaseline")
	}

	baseNext := staticStore.next.Load()
	baseBlocks := len(staticStore.blocks)
	baseSourceFiles := staticStore.sourceFiles.count.Load()
	baseColors := staticStore.colors.count.Load()
	baseBigInts := staticStore.bigInts.count.Load()

	numberTable.mu.Lock()
	baseNumberNext := numberTable.next
	baseNumberIndexLen := len(numberTable.index)
	numberTable.mu.Unlock()

	for iter := range 3 {
		for i := range 200 {
			tmpl := NewTemplateNode(Token_SCRIPT)
			tmpl.SetStaticSourceFile(testStaticSourceFile{name: "req_file.js"})
			tmpl.SetOriginalName("origName")
			tmpl.SetColor(&colors.Color{})
			tmpl.setTC(7)
			b := new(JSDocInfo_Builder)
			b.RecordConstancy()
			tmpl.SetJSDocInfo(b.Build())

			s := NewNodeStore()
			current = s
			numNode := Node_NewNumber(float64(iter*10000+i) + 0.123456789)
			if numNode.GetDouble() == 0 {
				t.Fatal("unexpected zero double")
			}
			ResetNodeStore()
		}
		// Allocate a BigInt node directly in staticStore (current == nil):
		_ = Node_NewBigInt(big.NewInt(int64(42 + iter)))

		if staticStore.next.Load() <= baseNext {
			t.Fatalf("iter %d: expected staticStore.next to grow past %d", iter, baseNext)
		}
		numberTable.mu.Lock()
		grownNumNext := numberTable.next
		numberTable.mu.Unlock()
		if grownNumNext <= baseNumberNext {
			t.Fatalf("iter %d: expected numberTable.next to grow past %d", iter, baseNumberNext)
		}

		RewindStaticBaseline()

		if got := staticStore.next.Load(); got != baseNext {
			t.Fatalf("iter %d: staticStore.next = %d, want %d", iter, got, baseNext)
		}
		if got := len(staticStore.blocks); got != baseBlocks {
			t.Fatalf("iter %d: len(staticStore.blocks) = %d, want %d", iter, got, baseBlocks)
		}
		if got := staticStore.sourceFiles.count.Load(); got != baseSourceFiles {
			t.Fatalf("iter %d: sourceFiles.count = %d, want %d", iter, got, baseSourceFiles)
		}
		if got := staticStore.colors.count.Load(); got != baseColors {
			t.Fatalf("iter %d: colors.count = %d, want %d", iter, got, baseColors)
		}
		if got := staticStore.bigInts.count.Load(); got != baseBigInts {
			t.Fatalf("iter %d: bigInts.count = %d, want %d", iter, got, baseBigInts)
		}
		numberTable.mu.Lock()
		gotNumNext := numberTable.next
		gotNumIndexLen := len(numberTable.index)
		numberTable.mu.Unlock()
		if gotNumNext != baseNumberNext {
			t.Fatalf("iter %d: numberTable.next = %d, want %d", iter, gotNumNext, baseNumberNext)
		}
		if gotNumIndexLen != baseNumberIndexLen {
			t.Fatalf("iter %d: len(numberTable.index) = %d, want %d", iter, gotNumIndexLen, baseNumberIndexLen)
		}
	}

	// Verify RewindStaticBaseline() also resets an active per-compilation NodeStore
	// without requiring a separate ResetNodeStore() call:
	active := EnsureCompilationStore()
	_ = IR_Name("activeCompilationNode")
	if CurrentNodeStore() != active || active.IsStatic() {
		t.Fatal("expected non-static active compilation NodeStore before RewindStaticBaseline")
	}
	RewindStaticBaseline()
	if !CurrentNodeStore().IsStatic() {
		t.Fatal("expected RewindStaticBaseline() to reset current NodeStore back to nil (static fallback)")
	}

	// Verify strpool shards are also rewound by CaptureStaticBaseline / RewindStaticBaseline:
	baseAnchor := rhinoStringPoolIntern("rhino_baseline_strpool_anchor")
	CaptureStaticBaseline()
	baseStrpool := staticBaseline.strpool
	for i := range 1500 {
		_ = rhinoStringPoolIntern("rhino_ephemeral_req_ident_" + string(rune('a'+(i%26))) + "_" + string(rune('0'+(i/26))))
	}
	RewindStaticBaseline()
	if got := staticBaseline.strpool; got != baseStrpool {
		t.Fatalf("strpool baseline mutated: %v != %v", got, baseStrpool)
	}
	if rhinoStringFromID(baseAnchor) != "rhino_baseline_strpool_anchor" {
		t.Fatalf("baseline strpool anchor lost after RewindStaticBaseline")
	}
}
