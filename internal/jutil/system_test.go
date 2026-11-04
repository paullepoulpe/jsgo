package jutil_test

import (
	"slices"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

func TestSetSystemExitHook(t *testing.T) {
	var events []string
	var gotStatus int

	prev := jutil.SetSystemExitHook(func(status int) {
		events = append(events, "exit_hook")
		gotStatus = status
	})
	defer jutil.SetSystemExitHook(prev)

	if jutil.GetSystemExitHook() == nil {
		t.Fatal("expected non-nil GetSystemExitHook")
	}

	jutil.System_AtExit(func() { events = append(events, "atexit_1") })
	jutil.System_AtExit(func() { events = append(events, "atexit_2") })

	jutil.System_Exit(17)

	if gotStatus != 17 {
		t.Fatalf("got status %d, want 17", gotStatus)
	}
	wantEvents := []string{"atexit_2", "atexit_1", "exit_hook"}
	if !slices.Equal(events, wantEvents) {
		t.Fatalf("events = %v, want %v", events, wantEvents)
	}
}
