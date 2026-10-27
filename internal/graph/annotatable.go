// Port of com/google/javascript/jscomp/graph/Annotatable.java @ closure-compiler v20261021

/*
 * Copyright 2008 The Closure Compiler Authors.
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

package graph

// Annotatable is an object that has an annotation.
type Annotatable interface {
	// SetAnnotation annotates a piece of information to the object.
	//
	// Java: Annotatable.setAnnotation(Annotation)
	SetAnnotation(data Annotation)

	// GetAnnotation retrieves a piece of information that has been
	// annotated. Returns the annotation or null if the object has not been
	// annotated.
	//
	// Go-only: language-adapter — Java's `<A extends Annotation> A` is an
	// unchecked cast at the call site; callers type-assert the result.
	//
	// Java: Annotatable.getAnnotation()
	GetAnnotation() Annotation
}
