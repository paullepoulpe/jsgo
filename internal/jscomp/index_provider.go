// Port of com/google/javascript/jscomp/IndexProvider.java @ closure-compiler v20261015

/*
 * Copyright 2017 The Closure Compiler Authors.
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

import "reflect"

// IndexProvider builds a datastructure and incrementally updates it when
// scopes change.
type IndexProvider[T any] interface {
	// Get returns an instance of type T.
	//
	// Java: IndexProvider.get()
	Get() T

	// GetType returns a class literal specialized on T.
	//
	// Go-only: language-adapter — Java's Class<T> is a reflect.Type.
	//
	// Java: IndexProvider.getType()
	GetType() reflect.Type
}
