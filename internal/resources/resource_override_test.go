package resources_test

import (
	"bytes"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/resources"
)

func TestRegisterAndClearResourceOverrides(t *testing.T) {
	defer resources.ClearOverrides()

	origTypedAst, err := resources.LoadBytes("runtime_libs.typedast")
	if err != nil || len(origTypedAst) == 0 {
		t.Fatalf("LoadBytes(runtime_libs.typedast): err=%v len=%d", err, len(origTypedAst))
	}

	customPayload := []byte("custom-typedast-override")
	epoch0 := resources.OverridesEpoch()
	resources.RegisterOverride("/runtime_libs.typedast", customPayload)
	if resources.OverridesEpoch() <= epoch0 {
		t.Fatal("expected OverridesEpoch to advance after RegisterOverride")
	}
	// Mutate caller slice to verify defensive copy:
	customPayload[0] = 'X'

	got, err := resources.LoadBytes("runtime_libs.typedast")
	if err != nil || !bytes.Equal(got, []byte("custom-typedast-override")) {
		t.Fatalf("got override %q err=%v", got, err)
	}

	resources.RegisterOverride("custom/new_resource.txt", []byte("hello"))
	if !resources.Exists("/custom/new_resource.txt") || resources.LoadText("custom/new_resource.txt") != "hello" {
		t.Fatal("expected custom/new_resource.txt override to exist")
	}

	resources.UnregisterOverride("runtime_libs.typedast")
	restored, err := resources.LoadBytes("runtime_libs.typedast")
	if err != nil || len(restored) != len(origTypedAst) {
		t.Fatalf("expected original runtime_libs.typedast after UnregisterOverride, got len=%d err=%v", len(restored), err)
	}

	resources.ClearOverrides()
	if resources.Exists("custom/new_resource.txt") {
		t.Fatal("expected custom/new_resource.txt to be removed after ClearOverrides")
	}
}
