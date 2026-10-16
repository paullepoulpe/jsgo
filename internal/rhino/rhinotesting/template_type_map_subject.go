// Port of com/google/javascript/rhino/testing/TemplateTypeMapSubject.java @ closure-compiler v20260415

package rhinotesting

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// TemplateTypeMapSubject is a Truth Subject for TemplateTypeMap.
type TemplateTypeMapSubject struct {
	t      testing.TB
	actual *rhino.TemplateTypeMap
}

// Java: TemplateTypeMapSubject.assertThat(TemplateTypeMap)
func TemplateTypeMapSubject_AssertThat(t testing.TB, actual *rhino.TemplateTypeMap) *TemplateTypeMapSubject {
	return &TemplateTypeMapSubject{t: t, actual: actual}
}

// Java: TemplateTypeMapSubject.hasKeysAndValues(ImmutableList<TemplateType>, ImmutableList<? extends JSType>)
func (s *TemplateTypeMapSubject) HasKeysAndValues(keys []*rhino.TemplateType, values []rhino.JSType) {
	s.t.Helper()
	tr := Truth{T: s.t}
	tr.AssertWithMessage("value of: templateTypeMap.getTemplateKeys()").That(s.actual.GetTemplateKeys()).IsEqualTo(keys)
	tr.AssertWithMessage("value of: templateTypeMap.getTemplateValues()").That(s.actual.GetTemplateValues()).IsEqualTo(values)
}

// Java: TemplateTypeMapSubject.hasSubmapsAt(int...)
func (s *TemplateTypeMapSubject) HasSubmapsAt(indices ...int) {
	s.t.Helper()
	Truth{T: s.t}.AssertWithMessage("value of: templateTypeMap.getIndicesOfSubmapsForTesting()").
		That(s.actual.GetIndicesOfSubmapsForTesting()).
		ContainsExactlyElementsIn(indices)
}
