// Port of com/google/javascript/jscomp/modules/ModuleRequestResolver.java @ closure-compiler v20260415

package jscomp

// ModuleRequestResolver resolves requests for other modules.
type ModuleRequestResolver interface {
	// Resolve returns the module that this import references, if it exists
	// in the compilation (nil otherwise).
	//
	// Java: ModuleRequestResolver.resolve(Import)
	Resolve(i *Import) UnresolvedModule

	// ResolveWithExport returns the module that this export references, if
	// it exists in the compilation (nil otherwise).
	//
	// Java: ModuleRequestResolver.resolve(Export)
	ResolveWithExport(e *Export) UnresolvedModule
}
