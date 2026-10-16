// Port of com/google/javascript/jscomp/deps/SortedDependencies.java @ closure-compiler v20260415

package jscomp

import (
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// SortedDependencies is a sorted list of inputs following the ES6 module
// ordering spec.
//
// Orders such that each input always comes after its dependencies. Circular
// references are allowed by emitting the current input in the moment before
// a loop would complete.
//
// The resulting order is not the same as the user-provided order but is
// influenced by it since it may take more than one graph traversal to
// account for all provided inputs and the graph traversals start with the
// first user provided input and continue from there.
//
// Also exposes other information about the inputs, like which inputs do not
// provide symbols.
//
// Inputs are compared with equals/hashCode, as in the Java collections.
type SortedDependencies[InputT DependencyInfo] struct {
	userOrderedInputs          *jutil.ArrayList[InputT]
	importOrderedInputs        *jutil.ArrayList[InputT]
	completedInputs            *jutil.ObjectLinkedHashSet[InputT]
	nonExportingInputs         *jutil.LinkedHashMap[string, InputT]
	exportingInputBySymbolName *jutil.LinkedHashMap[string, InputT]
	// Maps an input A to the inputs it depends on, ie, inputs that provide
	// stuff that A requires.
	importedInputByImportingInput *jutil.SetMultimap[InputT, InputT]
}

// Java: SortedDependencies(List<InputT>)
func NewSortedDependencies[InputT DependencyInfo](userOrderedInputs []InputT) *SortedDependencies[InputT] {
	d := &SortedDependencies[InputT]{
		userOrderedInputs:             jutil.NewArrayList[InputT](),
		importOrderedInputs:           jutil.NewArrayList[InputT](),
		completedInputs:               jutil.NewObjectLinkedHashSet[InputT](),
		nonExportingInputs:            jutil.NewLinkedHashMap[string, InputT](),
		exportingInputBySymbolName:    jutil.NewLinkedHashMap[string, InputT](),
		importedInputByImportingInput: jutil.NewSetMultimap[InputT, InputT](),
	}
	d.userOrderedInputs.AddAll(userOrderedInputs...)
	d.processInputs()
	return d
}

// GetStrongDependenciesOf gets all the strong dependencies of the given
// roots. The inputs must be returned in a stable order. In other words, if A
// comes before B, and A does not transitively depend on B, then A must also
// come before B in the returned list.
//
// sorted: if true, get them in topologically sorted order. If false, get
// them in the original order they were passed to the compiler.
//
// Java: SortedDependencies.getStrongDependenciesOf(List<InputT>, boolean)
func (d *SortedDependencies[InputT]) GetStrongDependenciesOf(rootInputs []InputT, sorted bool) []InputT {
	includedInputs := jutil.NewObjectLinkedHashSet[InputT]()
	worklist := append([]InputT(nil), rootInputs...)
	for len(worklist) > 0 {
		input := worklist[0]
		worklist = worklist[1:]
		if includedInputs.Add(input) {
			for _, symbolName := range input.GetRequiredSymbols() {
				if importedSymbolName, ok := d.exportingInputBySymbolName.Lookup(symbolName); ok && !jutil.IsNil(importedSymbolName) {
					worklist = append(worklist, importedSymbolName)
				}
			}
		}
	}

	var builder []InputT
	inputs := d.userOrderedInputs
	if sorted {
		inputs = d.importOrderedInputs
	}
	for input := range inputs.All() {
		if includedInputs.Contains(input) {
			builder = append(builder, input)
		}
	}
	return builder
}

// GetInputProviding returns the input that gives us the given symbol. It
// panics with *SortedDependencies_MissingProvideException if there is no
// input for this symbol.
//
// Java: SortedDependencies.getInputProviding(String)
func (d *SortedDependencies[InputT]) GetInputProviding(symbolName string) InputT {
	input := d.MaybeGetInputProviding(symbolName)
	if !jutil.IsNil(input) {
		return input
	}
	panic(NewSortedDependencies_MissingProvideException(symbolName))
}

// Java: SortedDependencies.getInputsWithoutProvides()
func (d *SortedDependencies[InputT]) GetInputsWithoutProvides() []InputT {
	return slices.Collect(d.nonExportingInputs.Values())
}

// GetSortedStrongDependenciesOf gets all the strong dependencies of the
// given roots, in topologically sorted order.
//
// Java: SortedDependencies.getSortedStrongDependenciesOf(List<InputT>)
func (d *SortedDependencies[InputT]) GetSortedStrongDependenciesOf(roots []InputT) []InputT {
	return d.GetStrongDependenciesOf(roots, true)
}

// GetSortedWeakDependenciesOf gets all the weak dependencies of the given
// roots. The inputs must be returned in stable order. In other words, if A
// comes before B, and A does not transitively depend on B, then A must also
// come before B in the returned list.
//
// The weak dependencies are those that are only reachable via type requires
// from the roots. Note that if a root weakly requires another input, then
// all of its transitive dependencies (strong or weak) that are not strongly
// reachable from the roots will be included. e.g. if A weakly requires B,
// and B strongly requires C, and A is the sole root, then this will return
// B and C. However, if we add D as a root, and D strongly requires C, then
// this will only return B.
//
// Root inputs will never be in the returned list as they are all considered
// strong.
//
// Java: SortedDependencies.getSortedWeakDependenciesOf(List<InputT>)
func (d *SortedDependencies[InputT]) GetSortedWeakDependenciesOf(rootInputs []InputT) []InputT {
	strongInputs := jutil.NewObjectLinkedHashSet(d.GetSortedStrongDependenciesOf(rootInputs)...)
	weakInputs := jutil.NewObjectLinkedHashSet[InputT]()
	worklist := strongInputs.ToSlice()
	for len(worklist) > 0 {
		input := worklist[0]
		worklist = worklist[1:]
		isStrong := strongInputs.Contains(input)

		var edges []string
		if isStrong {
			edges = input.GetTypeRequires()
		} else {
			edges = append(append([]string(nil), input.GetRequiredSymbols()...), input.GetTypeRequires()...)
		}

		if !isStrong && !weakInputs.Add(input) {
			continue
		}

		for _, symbolName := range edges {
			importedSymbolName, ok := d.exportingInputBySymbolName.Lookup(symbolName)
			if ok && !jutil.IsNil(importedSymbolName) &&
				!strongInputs.Contains(importedSymbolName) &&
				!weakInputs.Contains(importedSymbolName) {
				worklist = append(worklist, importedSymbolName)
			}
		}
	}

	var builder []InputT
	for input := range d.importOrderedInputs.All() {
		if weakInputs.Contains(input) {
			builder = append(builder, input)
		}
	}
	return builder
}

// GetSortedList returns the inputs in import order (an unmodifiable view in
// Java; callers must not modify it).
//
// Java: SortedDependencies.getSortedList()
func (d *SortedDependencies[InputT]) GetSortedList() []InputT {
	return d.importOrderedInputs.Slice()
}

// MaybeGetInputProviding returns the input that gives us the given symbol,
// or null (the zero InputT).
//
// Java: SortedDependencies.maybeGetInputProviding(String)
func (d *SortedDependencies[InputT]) MaybeGetInputProviding(symbol string) InputT {
	if d.exportingInputBySymbolName.ContainsKey(symbol) {
		return d.exportingInputBySymbolName.Get(symbol)
	}
	return d.nonExportingInputs.Get(ModuleNames_FileToModuleName(symbol))
}

func (d *SortedDependencies[InputT]) orderInput(input InputT) {
	if d.completedInputs.Contains(input) {
		return
	}

	d.completedInputs.Add(input)
	for importedInput := range d.importedInputByImportingInput.Get(input).All() {
		d.orderInput(importedInput)
	}

	// Emit an input after its imports have been emitted.
	d.importOrderedInputs.Add(input)
}

func (d *SortedDependencies[InputT]) processInputs() {
	// Index.
	for userOrderedInput := range d.userOrderedInputs.All() {
		provides := userOrderedInput.GetProvides()
		var firstProvide *string
		if len(provides) > 0 {
			firstProvide = &provides[0]
		}
		if firstProvide == nil ||
			// "module$" indicates the provide is generated from the path. If
			// this is the only thing the module provides and it is not an ES6
			// module then it is just a script and doesn't export anything.
			(len(provides) == 1 &&
				strings.HasPrefix(*firstProvide, "module$") &&
				// ES6 modules should always be considered as exporting
				// something.
				!userOrderedInput.IsEs6Module()) {
			d.nonExportingInputs.Put(ModuleNames_FileToModuleName(userOrderedInput.GetName()), userOrderedInput)
		}
		for _, providedSymbolName := range userOrderedInput.GetProvides() {
			d.exportingInputBySymbolName.Put(providedSymbolName, userOrderedInput)
		}
	}
	for userOrderedInput := range d.userOrderedInputs.All() {
		for _, symbolName := range userOrderedInput.GetRequiredSymbols() {
			importedInput, ok := d.exportingInputBySymbolName.Lookup(symbolName)
			if ok && !jutil.IsNil(importedInput) {
				d.importedInputByImportingInput.Put(userOrderedInput, importedInput)
			}
		}
	}

	// Order.
	// For each input, traverse in user-provided order.
	for userOrderedInput := range d.userOrderedInputs.All() {
		// Traverse the graph starting from this input and record any
		// newly-reached inputs.
		d.orderInput(userOrderedInput)
	}

	// Free temporary indexes.
	d.completedInputs.Clear()
	d.importedInputByImportingInput.Clear()
}

// SortedDependencies_MissingProvideException is
// SortedDependencies.MissingProvideException, a checked exception.
type SortedDependencies_MissingProvideException struct{ jutil.JavaException }

func (*SortedDependencies_MissingProvideException) missingProvideException() {}

// ISortedDependencies_MissingProvideException matches
// SortedDependencies_MissingProvideException.
type ISortedDependencies_MissingProvideException interface {
	jutil.IException
	missingProvideException()
}

// Java: SortedDependencies.MissingProvideException(String)
func NewSortedDependencies_MissingProvideException(provide string) *SortedDependencies_MissingProvideException {
	e := &SortedDependencies_MissingProvideException{}
	jutil.InitThrowable(e, "com.google.javascript.jscomp.deps.SortedDependencies$MissingProvideException", &provide, nil)
	return e
}

// Java: SortedDependencies.MissingProvideException(String, Exception)
func NewSortedDependencies_MissingProvideExceptionWithCause(provide string, e jutil.Throwable) *SortedDependencies_MissingProvideException {
	ex := &SortedDependencies_MissingProvideException{}
	jutil.InitThrowable(ex, "com.google.javascript.jscomp.deps.SortedDependencies$MissingProvideException", &provide, e)
	return ex
}
