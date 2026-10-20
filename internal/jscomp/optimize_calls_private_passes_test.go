package jscomp

import (
	"reflect"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/diagnostic"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type optimizeCallsPrivateClosePanicLog struct{ diagnostic.LogFile }

func (*optimizeCallsPrivateClosePanicLog) Close() { panic("private-close") }

type optimizeCallsPrivateClosePanicCompiler struct{ *Compiler }

func (c *optimizeCallsPrivateClosePanicCompiler) CreateOrReopenIndexedLog(_ reflect.Type, _ string, _ ...string) diagnostic.LogFile {
	return &optimizeCallsPrivateClosePanicLog{}
}

func TestOptimizeCallsPrivateParametersRestoresBorrowedScope(t *testing.T) {
	for _, shouldPanic := range []bool{false, true} {
		t.Run(map[bool]string{false: "return", true: "panic"}[shouldPanic], func(t *testing.T) {
			fixture := newSyntacticScopeCreatorTest(t)
			var compiler AbstractCompiler = fixture.compiler
			if shouldPanic {
				compiler = &optimizeCallsPrivateClosePanicCompiler{fixture.compiler}
			}
			adapter := newOptimizeCallsPrivateParameters(compiler).(optimizeCallsPrivateParameters)
			previousScope := Scope_CreateGlobalScope(rhino.IR_Root())
			borrowedScope := Scope_CreateGlobalScope(rhino.IR_Root())
			adapter.pass.globalScope = previousScope
			references := NewOptimizeCalls_ReferenceMap()
			references.globalScope = borrowedScope
			var caught any
			func() {
				defer func() { caught = recover() }()
				adapter.Process(rhino.IR_Root(), rhino.IR_Root(), references)
			}()
			if shouldPanic {
				if caught != "private-close" {
					t.Fatalf("panic = %v, want private-close", caught)
				}
			} else if caught != nil {
				t.Fatalf("unexpected panic: %v", caught)
			}
			if adapter.pass.globalScope != previousScope {
				t.Fatal("private adapter retained its borrowed global scope")
			}
		})
	}
}
