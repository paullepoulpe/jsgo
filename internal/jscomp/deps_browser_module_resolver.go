// Port of com/google/javascript/jscomp/deps/BrowserModuleResolver.java @ closure-compiler v20260415

package jscomp

import "github.com/paullepoulpe/jsgo/internal/jutil"

// BrowserModuleResolver is the resolution algorithm for Browsers.
//
// Only unambiguous paths are supported. Paths must specify a file
// extension.
type BrowserModuleResolver struct {
	ModuleResolverBase
}

var _ ModuleResolver = (*BrowserModuleResolver)(nil)

// Java: BrowserModuleResolver.FACTORY
var BrowserModuleResolver_FACTORY ModuleLoader_ModuleResolverFactory = ModuleLoader_ModuleResolverFactoryFunc(
	func(modulePaths *jutil.LinkedHashSet[string], moduleRootPaths []string, errorHandler ErrorHandler, pathEscaper ModuleLoader_PathEscaper) ModuleResolver {
		return NewBrowserModuleResolver(modulePaths, moduleRootPaths, errorHandler, pathEscaper)
	})

// Java: BrowserModuleResolver(ImmutableSet<String>, ImmutableList<String>, ErrorHandler, PathEscaper)
func NewBrowserModuleResolver(
	modulePaths *jutil.LinkedHashSet[string],
	moduleRootPaths []string,
	errorHandler ErrorHandler,
	pathEscaper ModuleLoader_PathEscaper) *BrowserModuleResolver {
	r := &BrowserModuleResolver{}
	r.InitModuleResolverBase(modulePaths, moduleRootPaths, errorHandler, pathEscaper)
	return r
}

// ResolveJsModule may return nil.
//
// Java: BrowserModuleResolver.resolveJsModule(String, String, String, int, int)
func (r *BrowserModuleResolver) ResolveJsModule(
	scriptAddress string, moduleAddress string, sourcename string, lineno int, colno int) *string {

	if ModuleLoader_IsAmbiguousIdentifier(moduleAddress) {
		r.errorHandler.Report(
			CheckLevel_WARNING.Ptr(),
			JSError_MakeWithSourceLocation(
				&sourcename,
				lineno,
				colno,
				ModuleLoader_INVALID_MODULE_PATH,
				moduleAddress,
				ModuleLoader_ResolutionMode_BROWSER.String()))
		return nil
	}

	loadAddress := r.locate(scriptAddress, moduleAddress)
	if loadAddress == nil {
		r.errorHandler.Report(
			CheckLevel_WARNING.Ptr(),
			JSError_MakeWithSourceLocation(&sourcename, lineno, colno, ModuleLoader_LOAD_WARNING, moduleAddress))
	}
	return loadAddress
}

// ResolveJsModuleSilently may return nil.
//
// Java: BrowserModuleResolver.resolveJsModuleSilently(String, String)
func (r *BrowserModuleResolver) ResolveJsModuleSilently(scriptAddress string, moduleAddress string) *string {
	if ModuleLoader_IsAmbiguousIdentifier(moduleAddress) {
		return nil
	}
	return r.locate(scriptAddress, moduleAddress)
}
