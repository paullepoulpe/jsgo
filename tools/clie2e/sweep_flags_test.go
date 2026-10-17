package main

import (
	"slices"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jscomp"
)

// The flag sweep's diagnostic groups are the registered ones of the pinned
// compiler, in registration order, so that moving the pin cannot leave a
// group that --jscomp_error/--jscomp_warning/--jscomp_off accept out of the
// f-jscomp-* configurations. Go-only test: the port's registry is Java's.
func TestDiagnosticGroupsAreTheRegisteredOnes(t *testing.T) {
	registered := slices.Collect(jscomp.DiagnosticGroups_GetRegisteredGroups().Keys())
	if !slices.Equal(diagnosticGroups, registered) {
		t.Errorf("diagnosticGroups in sweep_flags.go:\n%q\nDiagnosticGroups.getRegisteredGroups():\n%q", diagnosticGroups, registered)
	}
}
