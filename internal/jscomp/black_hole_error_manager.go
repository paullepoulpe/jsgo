// Port of com/google/javascript/jscomp/BlackHoleErrorManager.java @ closure-compiler v20260415

package jscomp

// BlackHoleErrorManager is an ErrorManager that silently swallows all
// messages.
type BlackHoleErrorManager struct {
	SortingErrorManagerBase
}

var _ SortingErrorManager = (*BlackHoleErrorManager)(nil)

// Java: BlackHoleErrorManager()
func NewBlackHoleErrorManager() *BlackHoleErrorManager {
	m := &BlackHoleErrorManager{}
	m.InitSortingErrorManagerBase(m, nil)
	return m
}
