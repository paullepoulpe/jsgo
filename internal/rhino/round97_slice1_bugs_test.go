/*
 * Copyright 2026 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package rhino_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

func TestRound97ReadableTypeNameNilNodeAndZeroChildGuards(t *testing.T) {
	base := rhinotesting.NewBaseJSTypeTestCase(t)
	reg := base.Registry

	unknown := reg.GetNativeType(rhino.JSTypeNative_UNKNOWN_TYPE)
	if got := reg.GetJSTypeOrUnknown(rhino.NilNode); got != unknown {
		t.Fatalf("GetJSTypeOrUnknown(NilNode) = %v, want %v", got, unknown)
	}
	if got := reg.GetReadableTypeName(rhino.NilNode); got != "?" {
		t.Fatalf("GetReadableTypeName(NilNode) = %q, want %q", got, "?")
	}
	if got := reg.GetReadableTypeNameNoDeref(rhino.NilNode); got != "?" {
		t.Fatalf("GetReadableTypeNameNoDeref(NilNode) = %q, want %q", got, "?")
	}

	// Malformed 0-child GETPROP carrying a non-simple type must fall back cleanly
	// instead of dereferencing NilNode in GetJSTypeOrUnknown(n.GetFirstChild()).
	gp0 := rhino.Node_NewStringWithToken(rhino.Token_GETPROP, "bar")
	gp0.SetJSType(reg.GetNativeType(rhino.JSTypeNative_ALL_TYPE))
	if got := reg.GetReadableTypeName(gp0); got != "*" {
		t.Fatalf("GetReadableTypeName(0-child GETPROP) = %q, want %q", got, "*")
	}
	if got := reg.GetReadableTypeNameNoDeref(gp0); got != "*" {
		t.Fatalf("GetReadableTypeNameNoDeref(0-child GETPROP) = %q, want %q", got, "*")
	}

	// Constructor whose source is a malformed 0-child FUNCTION or CLASS node must
	// not crash in JSTypeRegistry_SimpleReadableJSTypeName.
	for _, tok := range []rhino.Token{rhino.Token_FUNCTION, rhino.Token_CLASS} {
		src := rhino.NewNode(tok)
		closer := reg.GetResolver().OpenForDefinition()
		ctor := reg.CreateConstructorType("AnonCtor", src, nil, nil, nil, false)
		inst := ctor.GetInstanceType()
		closer.Close()

		got, ok := rhino.JSTypeRegistry_SimpleReadableJSTypeName(inst)
		if !ok || got == "" {
			t.Fatalf("SimpleReadableJSTypeName(%v) = (%q, %v), want non-empty ok=true", tok, got, ok)
		}
	}
}

type failingStringWriter struct{}

func (failingStringWriter) WriteString(string) (int, error) {
	return 0, errors.New("boom")
}

func TestRound97NodeSourcePositionAndFormattingPolish(t *testing.T) {
	prev := rhino.SetCurrentNodeStore(rhino.NewNodeStore())
	defer rhino.SetCurrentNodeStore(prev)

	n := rhino.Node_NewString("line1\n\"quoted\"")

	// Unset position returns -1 / -1.
	if n.GetLineno() != -1 || n.GetCharno() != -1 {
		t.Fatalf("initial position = (%d, %d), want (-1, -1)", n.GetLineno(), n.GetCharno())
	}

	// Normal line/column packing & extraction.
	n.SetLinenoCharno(42, 17)
	if n.GetLineno() != 42 || n.GetCharno() != 17 {
		t.Fatalf("SetLinenoCharno(42, 17) -> (%d, %d)", n.GetLineno(), n.GetCharno())
	}

	// Column > Node_MAX_COLUMN_NUMBER clamps charno to Node_MAX_COLUMN_NUMBER while preserving line.
	n.SetLinenoCharno(100, rhino.Node_MAX_COLUMN_NUMBER+5)
	if n.GetLineno() != 100 || n.GetCharno() != rhino.Node_MAX_COLUMN_NUMBER {
		t.Fatalf("clamped column -> (%d, %d), want (100, %d)", n.GetLineno(), n.GetCharno(), rhino.Node_MAX_COLUMN_NUMBER)
	}

	// Negative line resets both line and charno to -1.
	n.SetLinenoCharno(-1, 5)
	if n.GetLineno() != -1 || n.GetCharno() != -1 {
		t.Fatalf("negative line -> (%d, %d), want (-1, -1)", n.GetLineno(), n.GetCharno())
	}

	// AppendStringTree works with both *strings.Builder and custom io.StringWriter.
	parent := rhino.IR_Block()
	parent.AddChildToBack(n)
	wantTree := parent.ToStringTree()
	var sb strings.Builder
	parent.AppendStringTree(&sb)
	if sb.String() != wantTree {
		t.Fatalf("AppendStringTree(*strings.Builder) = %q, want %q", sb.String(), wantTree)
	}

	// Error from custom io.StringWriter translates to jutil.IOException panic.
	func() {
		defer func() {
			r := recover()
			if _, ok := r.(jutil.IIOException); !ok {
				t.Fatalf("AppendStringTree(failingStringWriter) panic = %T (%v), want IIOException", r, r)
			}
		}()
		parent.AppendStringTree(failingStringWriter{})
	}()

	// ToJson escapes newlines and double quotes in node property values.
	n.PutProp(rhino.Node_Prop_NON_JSDOC_COMMENT, "line1\n\"quoted\"")
	var jsonBuf strings.Builder
	parent.ToJson(&jsonBuf)
	jsonOut := jsonBuf.String()
	if !strings.Contains(jsonOut, `\n\"quoted\"`) {
		t.Fatalf("ToJson() = %s, expected escaped newline and quotes", jsonOut)
	}
}
