package jscomp

import (
	"os"
	"testing"
)

// Test-only portability for checked-in Java corpus goldens. Per-oracle overrides
// retain their existing precedence and the historical default remains unchanged.
// An explicit global override must be usable: it must never silently reduce a
// requested corpus run to only hand-written inputs.
func oracleTestCorpusRoot(t *testing.T) string {
	t.Helper()
	root, explicit := os.LookupEnv("JSGO_ORACLE_CORPUS")
	if !explicit {
		return "/home/user/jsgo/.cache/corpus"
	}
	if root == "" {
		t.Fatal("JSGO_ORACLE_CORPUS is explicitly empty")
	}
	info, err := os.Stat(root)
	if err != nil {
		t.Fatalf("JSGO_ORACLE_CORPUS %q is unavailable: %v", root, err)
	}
	if !info.IsDir() {
		t.Fatalf("JSGO_ORACLE_CORPUS %q is not a directory", root)
	}
	return root
}
