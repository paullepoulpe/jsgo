// Port of com/google/javascript/refactoring/RequireAliasGeneratorTest.java @ closure-compiler v20260415

package jscomp

import (
	"testing"
)

func TestRequireAliasGenerator_GetShortNameForRequire(t *testing.T) {
	cases := []struct {
		namespace string
		aliases   []string
	}{
		{"ns.Name", []string{"Name", "NsName", "NsName0", "NsName1"}},
		{"ns.name", []string{"name", "nsName", "nsName0", "nsName1"}},
		// Skiplisted aliases
		{"Array", []string{"Array0"}},
		{"array", []string{"array0"}},
		{"goog.Promise", []string{"GoogPromise"}},
		{"goog.Thenable", []string{"GoogThenable"}},
		{"goog.array", []string{"googArray"}},
		{"goog.events.Event", []string{"EventsEvent"}},
		{"goog.math", []string{"googMath"}},
		{"goog.object", []string{"googObject"}},
		{"goog.string", []string{"googString"}},
		{"goog.structs.Map", []string{"StructsMap"}},
		{"goog.structs.Set", []string{"StructsSet"}},
		{"wiz.Object", []string{"WizObject"}},
		{"xid.String", []string{"XidString"}},
	}
	for _, c := range cases {
		var actual []string
		for alias := range RequireAliasGenerator_Over(c.namespace).All() {
			if len(actual) == len(c.aliases) {
				break
			}
			actual = append(actual, alias)
		}
		assertSameElements(t, c.namespace, actual, c.aliases)
	}
}
