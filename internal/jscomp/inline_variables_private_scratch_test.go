package jscomp

import "testing"

// Fixed-only private guard test; omit from the stock public-API negative overlay.
func TestInlineVariablesPrivateScratchCallbackOwnership(t *testing.T) {
	// A callback can request another public expert while the private driver owns
	// its scratch slot. The callback result and active scratch remain separate.
	behavior, _, a, ar, b, br := inlineExpertOwnershipFixture()
	behavior.initScratchVarExpert(a, ar, false, false)
	behavior.scratchExpertInUse = true
	callback := behavior.CreateVarExpert(b, br, false).(*InlineVariables_InliningBehavior_StandardVarExpert)
	if callback == &behavior.scratchExpert || behavior.scratchExpert.v != a || callback.v != b || !behavior.scratchExpertInUse {
		t.Fatal("public callback reused or altered the private driver's active expert")
	}
}
