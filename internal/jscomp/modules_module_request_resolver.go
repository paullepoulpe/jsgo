// Port of com/google/javascript/jscomp/modules/ModuleRequestResolver.java @ closure-compiler v20261021

/*
 * Copyright 2019 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
