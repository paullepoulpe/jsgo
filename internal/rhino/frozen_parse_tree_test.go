package rhino

import (
	"math/big"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestFrozenParseTreePayloadsAndFallback(t *testing.T) {
	cooked := "cooked"
	for _, n := range []*Node{NewNode(Token_ROOT), Node_NewStringWithToken(Token_NAME, "x"), Node_NewNumber(3.5), Node_NewBigInt(big.NewInt(123)), Node_NewTemplateLitString(&cooked, "raw")} {
		n.SetOriginalName(jutil.StringPtr("original"))
		n.SetLinenoCharno(4, 7)
		n.SetLength(9)
		file := NewSimpleSourceFile("externs.zip//payload.js", StaticSourceFile_SourceKind_EXTERN)
		n.SetStaticSourceFile(file)
		frozen := Node_FreezeParseTree(n)
		if frozen == nil {
			t.Fatal("payload incorrectly rejected")
		}
		got := frozen.CloneWithSourceFile(file)
		if got.String() != n.CloneTree().String() || !got.IsEquivalentTo(n) {
			t.Fatal("payload/source/position changed")
		}
		got.SetOriginalName(jutil.StringPtr("changed"))
		if *n.GetOriginalName() != "original" {
			t.Fatal("mutable extended payload shared")
		}
	}
	unknown := IR_Name("unknown")
	unknown.PutProp(Node_Prop_DECLARED_TYPE_EXPR, IR_Name("type"))
	if Node_FreezeParseTree(unknown) != nil {
		t.Fatal("unknown mutable Node property must fall back")
	}
	shadow := IR_Name("shadow")
	shadow.SetClosureUnawareShadow(IR_Root())
	if Node_FreezeParseTree(shadow) != nil {
		t.Fatal("closure-unaware shadow must fall back")
	}
}
